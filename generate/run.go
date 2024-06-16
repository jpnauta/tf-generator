package generate

import (
	"fmt"
)

// Run main entry point for the `generate` command
func Run(filePath string, check bool) error {
	fmt.Printf("Loading %s and its references...\n", filePath)
	generateFile, err := LoadGenerateFile(filePath)
	if err != nil {
		return err
	}

	results, diags := generateFile.LoadAll()
	if diags.HasErrors() {
		return diags
	}

	for _, result := range results {
		if check {
			fmt.Printf("checking the file contents of %s...\n", result.OutputFile)
			if err := result.Check(); err != nil {
				return err
			}
		} else {
			fmt.Printf("updating %s)...\n", result.OutputFile)
			if err := result.Save(); err != nil {
				return err
			}
		}
	}

	fmt.Println("DONE")
	return nil
}
