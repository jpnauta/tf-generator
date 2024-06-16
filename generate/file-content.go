package generate

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"tf-generator/inject"
	"tf-generator/tf"
	"path"
	"strings"
)

// FileContent Used to serialize and deserialize file contents between HCL function calls,
// as well as perform the key actions in HCL functions.
type FileContent struct {
	SourcePath string
	Content    string
}

type FileContents []*FileContent

// CtyFileContentType Expected format when a FileContent is converted to/from a cty.Value
var CtyFileContentType = cty.Object(map[string]cty.Type{
	"source-path": cty.String,
	"content":     cty.String,
})
var CtyFileContentsType = cty.List(CtyFileContentType)

func NewFileContent(sourcePath string, content string) *FileContent {
	return &FileContent{
		SourcePath: sourcePath,
		Content:    content,
	}
}

func LoadFileContent(value cty.Value) *FileContent {
	fc := FileContent{}
	it := value.ElementIterator()
	for it.Next() {
		k, v := it.Element()
		if k.AsString() == "source-path" {
			fc.SourcePath = v.AsString()
		} else if k.AsString() == "content" {
			fc.Content = v.AsString()
		}
	}
	return &fc
}

func LoadFileContents(value cty.Value) FileContents {
	fcs := FileContents{}
	it := value.ElementIterator()
	for it.Next() {
		_, v := it.Element()
		fc := LoadFileContent(v)
		fcs = append(fcs, fc)
	}
	return fcs
}

func (fc *FileContent) ToCty() cty.Value {
	values := map[string]cty.Value{
		"source-path": cty.StringVal(fc.SourcePath),
		"content":     cty.StringVal(fc.Content),
	}
	return cty.ObjectVal(values)
}

func (fcs FileContents) Combine() *FileContent {
	contents := []string{}
	for _, fc := range fcs {
		content := fc.Content
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		contents = append(contents, content)
	}
	return NewFileContent("", strings.Join(contents, "\n"))
}

func (fc *FileContent) loadTfvars() (*tf.Tfvars, error) {
	return tf.NewTfvars(fc.SourcePath, []byte(fc.Content))
}

func (fcs FileContents) MergeTfvars() (*FileContent, error) {
	// Parse all tfvars content
	tfvarsList := []*tf.Tfvars{}
	for _, fc := range fcs {
		tfvars, err := fc.loadTfvars()
		if err != nil {
			return nil, err
		}
		tfvarsList = append(tfvarsList, tfvars)
	}

	// Merge tfvars
	merged := tf.MergeTfvars(tfvarsList)

	// Export as content
	content, err := tf.HclAsString(merged)
	if err != nil {
		return nil, err
	}
	return NewFileContent("", content), nil
}

func (fc *FileContent) RemoveKeys(removeFc *FileContent) (*FileContent, error) {
	tfvars, err := fc.loadTfvars()
	if err != nil {
		return nil, err
	}
	removeTfvars, err := removeFc.loadTfvars()
	if err != nil {
		return nil, err
	}

	removedTfvars := tfvars.RemoveKeys(removeTfvars.Keys())

	content, err := tf.HclAsString(removedTfvars)
	if err != nil {
		return nil, err
	}

	return NewFileContent("", content), nil
}

func (fcs FileContents) CombineWithInject(tfvarsContent *FileContent) (*FileContent, error) {
	injections := []*inject.TfFileInjection{}
	for _, fc := range fcs {
		tfFile, err := tf.NewTfFile(fc.SourcePath, []byte(fc.Content))
		if err != nil {
			return nil, err
		}
		injection := inject.NewTfFileInjection(path.Dir(fc.SourcePath), tfFile)
		injections = append(injections, injection)
	}

	// Load tfvars content
	tfvars, err := tfvarsContent.loadTfvars()
	if err != nil {
		return nil, err
	}

	injector := inject.NewTfInjector(injections, tfvars)

	// Export as content
	content, err := tf.HclAsString(injector)
	if err != nil {
		return nil, err
	}
	return NewFileContent("", content), nil
}

func (fc *FileContent) GetTfvar(key string) (*FileContent, error) {
	tfvars, err := fc.loadTfvars()
	if err != nil {
		return nil, err
	}
	value, ok := tfvars.Values[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in tfvars", key)
	}
	return NewFileContent("", value.AsString()), nil
}
