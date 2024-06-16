package generate

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"os"
	"path"
	"slices"
	"sort"
)

// GenerateContext implements all HCL functions available in generate files, and tracks
// the values of and locals.
type GenerateContext struct {
	RootDir     string
	EvalContext *hcl.EvalContext
	locals      map[string]cty.Value
}

func NewGenerateContext(rootDir string) *GenerateContext {
	return &GenerateContext{
		RootDir: rootDir,
		EvalContext: &hcl.EvalContext{
			Functions: map[string]function.Function{
				"load":                loadFunc(rootDir),
				"combine":             combineFunc(),
				"merge-tfvars":        mergeTfvarsFunc(),
				"remove-tfvar-keys":   removeTfvarKeys(),
				"combine-with-inject": combineWithInjectFunc(),
				"get-tfvar":           getTfvarFunc(),
			},
			Variables: map[string]cty.Value{
				"local": cty.ObjectVal(map[string]cty.Value{}),
			},
		},
		locals: map[string]cty.Value{},
	}
}

func (gc *GenerateContext) addLocal(name string, val cty.Value) bool {
	if _, ok := gc.locals[name]; ok {
		return false
	}
	gc.locals[name] = val
	gc.EvalContext.Variables["local"] = cty.ObjectVal(gc.locals)
	return true
}

func nestedRefs(varReferences map[string][]string, varName string) []string {
	nested := []string{}
	if _, ok := varReferences[varName]; ok {
		for _, refName := range varReferences[varName] {
			nested = append(nested, refName)
			if _, ok := varReferences[refName]; ok {
				nested = append(nested, nestedRefs(varReferences, refName)...)
			}
		}
	}
	return nested
}

// AddLocals loads the values of all locals defined in the body
func (gc *GenerateContext) AddLocals(body hcl.Body) hcl.Diagnostics {
	attributesMap, diags := body.JustAttributes()
	if diags.HasErrors() {
		return diags
	}
	attributes := []*hcl.Attribute{}
	for _, v := range attributesMap {
		attributes = append(attributes, v)
	}

	// Ensure locals referencing other locals are loaded last
	varReferences := map[string][]string{}
	for _, attribute := range attributes { // Find immediate references to other locals
		for _, variable := range attribute.Expr.Variables() {
			if variable.RootName() == "local" && len(variable) == 2 {
				if ref, ok := variable[1].(hcl.TraverseAttr); ok {
					varReferences[attribute.Name] = append(varReferences[attribute.Name], ref.Name)
				}
			}
		}
	}
	for varName, refNames := range varReferences { // Find nested references to other locals
		for _, refName := range nestedRefs(varReferences, varName) {
			if !slices.Contains(refNames, refName) {
				varReferences[varName] = append(varReferences[varName], refName)
			}
		}
	}
	sort.Slice(attributes, func(i, j int) bool {
		return !slices.Contains(varReferences[attributes[j].Name], attributes[i].Name)
	})
	slices.Reverse(attributes)

	for _, attr := range attributes {
		val, diags := attr.Expr.Value(gc.EvalContext)
		if diags.HasErrors() {
			return diags
		}

		if !gc.addLocal(attr.Name, val) {
			return hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("local %q already defined", attr.Name),
					Subject:  &attr.Range,
				},
			}
		}
	}
	return nil
}

func loadFunc(rootDir string) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{Name: "source", Type: cty.String},
		},
		Type: function.StaticReturnType(CtyFileContentType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			sourcePath := args[0].AsString()
			filePath := path.Join(rootDir, sourcePath)
			content, err := os.ReadFile(filePath)
			if err != nil {
				return cty.NilVal, err
			}
			fileContent := NewFileContent(sourcePath, string(content))
			return fileContent.ToCty(), nil
		},
	})
}

func combineFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{Name: "tfvars", Type: CtyFileContentsType},
		},
		Type: function.StaticReturnType(CtyFileContentType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			fileContents := LoadFileContents(args[0])
			return fileContents.Combine().ToCty(), nil
		},
	})
}

func mergeTfvarsFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{Name: "tfvars", Type: CtyFileContentsType},
		},
		Type: function.StaticReturnType(CtyFileContentType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			fileContents := LoadFileContents(args[0])
			merged, err := fileContents.MergeTfvars()
			if err != nil {
				return cty.NilVal, err
			}
			return merged.ToCty(), nil
		},
	})
}

func removeTfvarKeys() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{Name: "tfvars", Type: CtyFileContentType},
			{Name: "removeTfvars", Type: CtyFileContentType},
		},
		Type: function.StaticReturnType(CtyFileContentType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			tfvarsContent := LoadFileContent(args[0])
			removeTfvarsContent := LoadFileContent(args[1])
			removedTfvars, err := tfvarsContent.RemoveKeys(removeTfvarsContent)
			if err != nil {
				return cty.NilVal, err
			}
			return removedTfvars.ToCty(), nil
		},
	})
}

func combineWithInjectFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{Name: "tf", Type: CtyFileContentsType},
			{Name: "tfvars", Type: CtyFileContentType},
		},
		Type: function.StaticReturnType(CtyFileContentType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			tfContents := LoadFileContents(args[0])
			tfvarsContent := LoadFileContent(args[1])
			combinedContent, err := tfContents.CombineWithInject(tfvarsContent)
			if err != nil {
				return cty.NilVal, err
			}
			return combinedContent.ToCty(), nil
		},
	})
}

func getTfvarFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{Name: "tfvars", Type: CtyFileContentType},
			{Name: "key", Type: cty.String},
		},
		Type: function.StaticReturnType(CtyFileContentType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			tfvarsContent := LoadFileContent(args[0])
			key := args[1].AsString()
			content, err := tfvarsContent.GetTfvar(key)
			if err != nil {
				return cty.NilVal, err
			}
			return content.ToCty(), nil
		},
	})
}
