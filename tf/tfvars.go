package tf

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfvars_parser "github.com/musukvl/tfvars-parser"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
	"os"
	"path"
	"slices"
	"sort"
	"strings"
)

type Tfvars struct {
	FileName string
	Values   map[string]cty.Value
}

// EmptyTfvars creates an empty tfvars
func EmptyTfvars() *Tfvars {
	return &Tfvars{
		Values: map[string]cty.Value{},
	}
}

func NewTfvars(fileName string, fileContent []byte) (*Tfvars, error) {
	rawTfvars, err := tfvars_parser.ConvertFileContent(fileContent, fileName)
	if err != nil {
		return nil, err
	}

	ctyTfvars := map[string]cty.Value{}
	for k, v := range rawTfvars {
		ctyValue, err := toCty(v)
		if err != nil {
			return nil, err
		}
		ctyTfvars[k] = *ctyValue
	}

	t := Tfvars{
		FileName: fileName,
		Values:   ctyTfvars,
	}
	return &t, nil
}

// LoadTfvarsFile loads the specified `.tfvars` file
func LoadTfvarsFile(filePath string) (*Tfvars, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return NewTfvars(path.Base(filePath), file)
}

// IterateTfvarsPaths iterates over all .tfvars files in a directory
func iterateTfvarsPaths(dirPath string, fn func(string) error) error {
	fileInfos, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, f := range fileInfos {
		fileName := f.Name()
		if !f.IsDir() && strings.HasSuffix(fileName, ".tfvars") {
			if err := fn(path.Join(dirPath, fileName)); err != nil {
				return err
			}
		}
	}
	return nil
}

// LoadTfvarsListFromProject loads all `.tfvars` files from a terraform project
func LoadTfvarsListFromProject(dirPath string) ([]*Tfvars, error) {
	var allTfvars []*Tfvars
	err := iterateTfvarsPaths(dirPath, func(filePath string) error {
		tfvars, err := LoadTfvarsFile(filePath)
		if err != nil {
			return err
		}
		allTfvars = append(allTfvars, tfvars)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allTfvars, nil
}

// LoadTfvarsFromProjectExcluding loads all `.tfvars` files from a terraform project and merges them into one
func LoadTfvarsFromProjectExcluding(dirPath string, ignoredPaths []string) (*Tfvars, error) {
	var allTfvars []*Tfvars
	err := iterateTfvarsPaths(dirPath, func(filePath string) error {
		if !slices.Contains(ignoredPaths, path.Base(filePath)) {
			tfvars, err := LoadTfvarsFile(filePath)
			if err != nil {
				return err
			}
			allTfvars = append(allTfvars, tfvars)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return MergeTfvars(allTfvars), nil
}

// RemoveKeys removes any of the keys from the current tfvars
func (t *Tfvars) RemoveKeys(keys []string) *Tfvars {
	cleanedTfvars := EmptyTfvars()
	for k, v := range t.Values {
		if !slices.Contains(keys, k) {
			cleanedTfvars.Values[k] = v
		}
	}
	return cleanedTfvars
}

// MergeTfvars merge all tfvars into one - first tfvars imported should take precedence
func MergeTfvars(tfVarsList []*Tfvars) *Tfvars {
	mergedTfvars := EmptyTfvars()
	for _, tfvars := range tfVarsList {
		for k, v := range tfvars.Values {
			if _, ok := mergedTfvars.Values[k]; !ok {
				mergedTfvars.Values[k] = v
			}
		}
	}
	return mergedTfvars
}

// Keys returns all keys in the tfvars
func (t *Tfvars) Keys() []string {
	keys := make([]string, 0, len(t.Values))
	for k := range t.Values {
		keys = append(keys, k)
	}
	return keys
}

// toCty converts an arbitrary value to a cty.Value
func toCty(value interface{}) (*cty.Value, error) {
	var ctyVal cty.Value
	if l, ok := value.([]interface{}); ok { // List
		ctyList := []cty.Value{}
		for _, v := range l {
			ctyV, err := toCty(v)
			if err != nil {
				return nil, err
			}
			ctyList = append(ctyList, *ctyV)
		}
		if len(ctyList) == 0 {
			ctyVal = cty.ListValEmpty(cty.DynamicPseudoType)
		} else {
			ctyVal = cty.ListVal(ctyList)
		}
	} else if m, ok := value.(map[string]interface{}); ok { // Object
		ctyMap := map[string]cty.Value{}
		for k, v := range m {
			ctyV, err := toCty(v)
			if err != nil {
				return nil, err
			}
			ctyMap[k] = *ctyV
		}
		ctyVal = cty.ObjectVal(ctyMap)
	} else if v, ok := value.(string); ok { // String
		ctyVal = cty.StringVal(v)
	} else if v, ok := value.(json.SimpleJSONValue); ok { // Ints/bool
		ctyVal = v.Value
	} else {
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
	return &ctyVal, nil
}

// Copy creates a copy of the current tfvars
func (t *Tfvars) Copy() *Tfvars {
	newTfvars := EmptyTfvars()
	for k, v := range t.Values {
		newTfvars.Values[k] = v
	}
	return newTfvars
}

func (t *Tfvars) ExportHCL(body *hclwrite.Body) error {
	// Sort keys by name for consistency
	keys := t.Keys()
	sort.Strings(keys)

	for _, key := range keys {
		body.SetAttributeValue(key, t.Values[key])
	}
	return nil
}
