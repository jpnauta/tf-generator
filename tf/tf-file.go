package tf

import (
	"bytes"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"os"
	"path"
)

type TfFile struct {
	FileName string
	File     *hclwrite.File
}

// NewTfFile creates a new `.tf` file from the contents of a file
func NewTfFile(filePath string, content []byte) (*TfFile, error) {
	if !bytes.HasSuffix(content, []byte("\n")) { // Ensures consistent spacing when merged
		content = append(content, '\n')
	}

	file, diags := hclwrite.ParseConfig(content, filePath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	return &TfFile{
		FileName: path.Base(filePath),
		File:     file,
	}, nil
}

// LoadTfFile loads a .tf file and parses its contents
func LoadTfFile(filePath string) (*TfFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return NewTfFile(path.Base(filePath), content)
}

func EmptyTfFile() *TfFile {
	return &TfFile{
		File: hclwrite.NewEmptyFile(),
	}
}

// CombineTfFiles merges the contents of multiple .tf files into one
func CombineTfFiles(tfFiles []*TfFile) *TfFile {
	mergedFile := hclwrite.NewEmptyFile()
	mergedFileBody := mergedFile.Body()
	for _, tfFile := range tfFiles {
		for _, block := range tfFile.File.Body().Blocks() {
			mergedFileBody.AppendBlock(block)
		}
	}
	return &TfFile{
		File: mergedFile,
	}
}

func (t *TfFile) ExportHCL(body *hclwrite.Body) error {
	blocks := t.File.Body().Blocks()
	for i, b := range blocks {
		body.AppendBlock(b)
		// Add newline after each block except the last one
		if i < len(blocks)-1 {
			body.AppendNewline()
		}
	}
	return nil
}

// VariableRefsFor finds all variables with the prefix of the specified typeName
// E.g. typeName = "var" will return all variable references like `var.my_var`
func (t *TfFile) VariableRefsFor(typeName string) []string {
	references := []string{}
	for _, block := range t.File.Body().Blocks() {
		for _, attribute := range block.Body().Attributes() {
			variables := attribute.Expr().Variables()
			for _, variable := range variables {
				tokens := variable.BuildTokens(hclwrite.Tokens{})
				if string(tokens[0].Bytes) == typeName {
					references = append(references, string(tokens[2].Bytes))
				}

			}
		}
	}

	return references
}
