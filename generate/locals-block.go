package generate

import (
	"github.com/hashicorp/hcl/v2"
)

// LocalsBlock Parses a single `locals{}` block
type LocalsBlock struct {
	Values hcl.Body `hcl:",remain"`
}

type LocalsBlocks []*LocalsBlock

// LoadAll parses the contents of all locals and stores them as locals for later usage
func (l LocalsBlocks) LoadAll(generateContext *GenerateContext) hcl.Diagnostics {
	for _, block := range l {
		if diags := generateContext.AddLocals(block.Values); diags.HasErrors() {
			return diags
		}
	}
	return nil
}
