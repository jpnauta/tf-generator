package tf

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// HCLExporter is an interface for exporting HCL
type HCLExporter interface {
	ExportHCL(body *hclwrite.Body) error
}

// HclAsString exports the given HCLExporter as a string
func HclAsString(e HCLExporter) (string, error) {
	f := hclwrite.NewEmptyFile()
	if err := e.ExportHCL(f.Body()); err != nil {
		return "", err
	}
	return string(f.Bytes()), nil
}
