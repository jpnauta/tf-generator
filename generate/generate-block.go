package generate

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"path"
)

const generatedFileHeader = "#DO NOT EDIT! This file was generated by tf-generator.\n"

// GenerateBlock represents a single `generate{}` block
type GenerateBlock struct {
	Content       hcl.Expression `hcl:"content"`
	Output        string         `hcl:"output"`
	ExcludeHeader bool           `hcl:"exclude-header,optional"`
}

type GenerateBlocks []*GenerateBlock

// Load parses the content of the `generate{}` block and returns a GenerateResult
func (g *GenerateBlock) Load(generateContext *GenerateContext) (*GenerateResult, hcl.Diagnostics) {
	var fc map[string]string
	diags := gohcl.DecodeExpression(g.Content, generateContext.EvalContext, &fc)
	if diags.HasErrors() {
		return nil, diags
	}
	content := fc["content"]
	if !g.ExcludeHeader {
		content = generatedFileHeader + content
	}

	fileName := path.Join(generateContext.RootDir, g.Output)
	return NewGenerateResult([]byte(content), fileName), nil
}

// LoadAll loads all `generate{}` blocks and returns all GenerateResults
func (l GenerateBlocks) LoadAll(generateContext *GenerateContext) (GenerateResults, hcl.Diagnostics) {
	results := GenerateResults{}
	for _, block := range l {
		result, diags := block.Load(generateContext)
		if diags.HasErrors() {
			return nil, diags
		}
		results = append(results, result)
	}
	return results, nil
}