package generate

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"os"
)

// GenerateResult tracks the result of a `generate{}` block and performs the resulting actions
type GenerateResult struct {
	Content    []byte
	OutputFile string
}

type GenerateResults []*GenerateResult

func NewGenerateResult(content []byte, name string) *GenerateResult {
	return &GenerateResult{
		Content:    content,
		OutputFile: name,
	}
}

// Check checks if the output file matches the result content
func (r *GenerateResult) Check() error {
	// Read the file into a byte slice
	actualFileContent, err := os.ReadFile(r.OutputFile)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(actualFileContent), string(r.Content), true)

	if !(len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual) {
		return fmt.Errorf(
			"the new tfvars file does not match the existing file.\n%s\n%s",
			dmp.DiffToDelta(diffs),
			dmp.DiffPrettyText(diffs),
		)
	}
	return nil
}

// Save replaces the output file with the result content
func (r *GenerateResult) Save() error {
	return os.WriteFile(r.OutputFile, r.Content, 0644)
}
