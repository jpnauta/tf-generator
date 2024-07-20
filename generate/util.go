package generate

import (
	"os"
	"path/filepath"
)

// findFilesByName finds all files with a specific name in a directory and its subdirectories
func findFilesByName(directory, filename string) ([]string, error) {
	var matches []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == filename {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return matches, nil
}

// hasDirs check if the file path specifies any directories
func hasDirs(filePath string) bool {
	dir, _ := filepath.Split(filePath)
	return dir != ""
}
