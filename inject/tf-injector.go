package inject

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"tf-generator/tf"
	"slices"
	"sort"
)

// TfFileInjector Performs injections on the given `.tf` files
type TfFileInjector struct {
	Injections []*TfFileInjection
	Tfvars     *tf.Tfvars
}

func NewTfInjector(injections []*TfFileInjection, tfvars *tf.Tfvars) *TfFileInjector {
	return &TfFileInjector{
		Injections: injections,
		Tfvars:     tfvars,
	}
}

// ExportHCL perform all injections and exports the combined `.tf` file to HCL
func (i *TfFileInjector) ExportHCL(body *hclwrite.Body) error {
	injectedTfFile, err := i.injectAndCombine()
	if err != nil {
		return err
	}
	return injectedTfFile.ExportHCL(body)
}

// injectAndCombine performs all injections and combines them into a single `.tf` file
func (i *TfFileInjector) injectAndCombine() (*tf.TfFile, error) {
	// Modify references to injected variables and keep track of injected locals
	allUsedTfvarNames := []string{}
	for _, injection := range i.Injections {
		usedTfvarNames, err := injection.injectAll(i.Tfvars)
		if err != nil {
			return nil, err
		}
		allUsedTfvarNames = append(allUsedTfvarNames, usedTfvarNames...)
	}

	// Process used tfvar names
	slices.Compact[[]string, string](allUsedTfvarNames) // Remove duplicates
	sort.Slice(allUsedTfvarNames, func(i, j int) bool { // sort alphabetically for consistency
		return allUsedTfvarNames[i] < allUsedTfvarNames[j]
	})

	// Add the injected `locals` block to the top of the file
	localsTfFile := tf.EmptyTfFile()
	if len(allUsedTfvarNames) > 0 {
		localsBody := localsTfFile.File.Body().AppendNewBlock("locals", nil).Body()
		for _, tfvarName := range allUsedTfvarNames {
			localsBody.SetAttributeValue(injectLocalPrefix+tfvarName, i.Tfvars.Values[tfvarName])
		}
		localsTfFile.File.Body().AppendNewline()
	}

	// Combine all locals and tf files
	tfilesToCombine := []*tf.TfFile{localsTfFile}
	for _, injection := range i.Injections {
		tfilesToCombine = append(tfilesToCombine, injection.TfFile)
	}
	return tf.CombineTfFiles(tfilesToCombine), nil
}
