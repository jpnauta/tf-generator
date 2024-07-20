package generate

import (
	"fmt"
	"os"
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

func RunRecursive(filename string, check bool) error {
	if hasDirs(filename) {
		return fmt.Errorf("file name %s cannot contain folders", filename)
	}

	cwd, _ := os.Getwd()
	files, err := findFilesByName(cwd, filename)
	if err != nil {
		return err
	}
	fmt.Printf("Found %d generator files matching glob\n", len(files))
	for _, file := range files {
		if err := Run(file, check); err != nil {
			return err
		}
	}
	return nil
}
