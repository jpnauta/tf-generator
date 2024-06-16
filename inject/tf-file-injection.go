package inject

import (
	"bytes"
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"tf-generator/tf"
)

const (
	injectTypeName    = "injectvar"
	injectLocalPrefix = "INJECTED-"
	contextTypeName   = "context"
)

type TfFileInjection struct {
	SourceDir string
	TfFile    *tf.TfFile
}

// NewTfFileInjection indicates a `.tf` file to be injected by TfFileInjector
func NewTfFileInjection(sourceDir string, tfFile *tf.TfFile) *TfFileInjection {
	return &TfFileInjection{
		SourceDir: sourceDir,
		TfFile:    tfFile,
	}
}

// injectAll performs all injections on the tf file
// and returns the list of variable names that were injected
func (i *TfFileInjection) injectAll(tfvars *tf.Tfvars) ([]string, error) {
	if err := i.injectContext(); err != nil {
		return nil, err
	}
	return i.injectTfvars(tfvars)
}

// injectContext injects context variables like `context.SOURCE_DIR`
func (i *TfFileInjection) injectContext() error {
	contextVars := map[string]string{
		"SOURCE_DIR": i.SourceDir,
	}

	newFileContents := i.TfFile.File.Bytes()
	for key, value := range contextVars {
		newFileContents = bytes.ReplaceAll(
			newFileContents,
			[]byte(fmt.Sprintf("${%s.%s}", contextTypeName, key)),
			[]byte(value),
		)
		newFileContents = bytes.ReplaceAll( // TODO test 2
			newFileContents,
			[]byte(fmt.Sprintf("%s.%s", contextTypeName, key)),
			[]byte(fmt.Sprintf("\"%s\"", value)),
		)
	}

	// Save changes
	tfFile, err := tf.NewTfFile("", newFileContents)
	if err != nil {
		return err
	}
	i.TfFile = tfFile
	return nil
}

// injectTfvars replaces references to `injectvar.my-var` with `local.injected-my-var`
// and returns the list of variable names that were injected
func (i *TfFileInjection) injectTfvars(tfvars *tf.Tfvars) ([]string, error) {
	variableRefs := []string{}
	newFileContents := i.TfFile.File.Bytes()
	for _, ref := range i.TfFile.VariableRefsFor(injectTypeName) {
		fullKey := fmt.Sprintf("%s.%s", injectTypeName, ref)

		// Find corresponding tfvar for local
		value, ok := tfvars.Values[ref]
		if !ok {
			return nil, fmt.Errorf("Could not find tfvar for ref `%s`", fullKey)
		}

		if value.Type() == cty.String {
			// Inject interpolated strings inline without a local
			newFileContents = bytes.ReplaceAll(
				newFileContents,
				[]byte(fmt.Sprintf("${%s}", fullKey)),
				[]byte(value.AsString()),
			)
		}
		if bytes.Contains(newFileContents, []byte(fullKey)) {
			// Replace `injectvar.my-var` with `local.INJECTED-my-var`
			newFileContents = bytes.ReplaceAll(
				newFileContents,
				[]byte(fullKey),
				[]byte(fmt.Sprintf("local.%s%s", injectLocalPrefix, ref)),
			)
			variableRefs = append(variableRefs, ref)
		}

	}
	tfFile, err := tf.NewTfFile("", newFileContents)
	if err != nil {
		return nil, err
	}
	i.TfFile = tfFile

	return variableRefs, nil
}
