package initialize

import (
	"fmt"
	"os"
)

type InitFile struct {
	FileName string
	Content  string
}

var InitFiles = []InitFile{
	{
		FileName: "tf-generator.hcl",
		Content: `locals {
  a        = load("a.tf")
  b        = load("b.tf")
  a-plus-b = combine([local.a, local.b])
}

generate {
  content = local.a-plus-b
  output  = "tf-generator.tf"
}
`,
	},
	{
		FileName: "a.tf",
		Content: `resource "a" "a" {
  name = "a"
}
`,
	},
	{
		FileName: "b.tf",
		Content: `resource "b" "b" {
  name = "b"
}
`,
	},
}

func Run() error {
	// Check if files already exist
	for _, initFile := range InitFiles {
		filePath := initFile.FileName

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			return fmt.Errorf("file '%s' already exists", filePath)
		}

	}

	// Write files
	for _, initFile := range InitFiles {
		filePath := initFile.FileName
		err := os.WriteFile(initFile.FileName, []byte(initFile.Content), 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %v", filePath, err)
		}

		fmt.Printf("Created file '%s'\n", filePath)
	}

	fmt.Printf("DONE\n\n")
	fmt.Printf("Please run 'tf-generator generate' to generate the resulting 'tf-generator.tf'.\n")

	return nil
}
