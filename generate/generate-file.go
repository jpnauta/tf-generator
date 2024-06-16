package generate

import (
	"bytes"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"os"
	"path"
)

// GenerateFile parses the contents of a codegen .hcl file
type GenerateFile struct {
	Locals          LocalsBlocks   `hcl:"locals,block"`
	GenerateBlocks  GenerateBlocks `hcl:"generate,block"`
	FilePath        string
	GenerateContext *GenerateContext
}

// LoadGenerateFile loads a GenerateFile from an `.hcl` file
func LoadGenerateFile(filePath string) (*GenerateFile, error) {
	g := &GenerateFile{
		GenerateBlocks:  []*GenerateBlock{},
		FilePath:        filePath,
		GenerateContext: NewGenerateContext(path.Dir(filePath)),
	}

	ctx := g.GenerateContext.EvalContext
	if err := hclsimple.DecodeFile(filePath, ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

// LoadAll Invokes all HCL functions in the generate file, and returns the contents of
// the files to be checked/written.
func (g *GenerateFile) LoadAll() (GenerateResults, hcl.Diagnostics) {
	// Parse locals if needed
	if diags := g.Locals.LoadAll(g.GenerateContext); diags.HasErrors() {
		return nil, diags
	}

	// Determine generate results
	return g.GenerateBlocks.LoadAll(g.GenerateContext)
}

// Save writes the GenerateFile to an `.hcl` file.
func (g *GenerateFile) Save() error {
	// Render file
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(g, f.Body())
	fileContent := f.Bytes()

	// Remove the leading newline
	if bytes.HasPrefix(fileContent, []byte("\n")) {
		fileContent = fileContent[1:]
	}

	// Add empty newline between import-tfvars blocks
	fileContent = bytes.ReplaceAll(fileContent, []byte("}\nimport-tfvars {"), []byte("}\n\nimport-tfvars {"))

	// Write file
	return os.WriteFile(g.FilePath, fileContent, 0644)
}
