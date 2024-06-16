package main

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

type ValidFixture struct {
	dirPath string
}

type InvalidFixture struct {
	dirPath                 string
	expectedMessageContains string
	isDiag                  bool
}

func TestValidFixtures(t *testing.T) {
	for _, fixture := range []ValidFixture{
		{
			dirPath: "examples/basics-01-load-and-combine/",
		},
		{
			dirPath: "examples/basics-02-locals/",
		},
		{
			dirPath: "examples/basics-03-merge-tfvars/",
		},
		{
			dirPath: "examples/basics-04-remove-tfvar-keys/",
		},
		{
			dirPath: "examples/basics-05-combine-with-inject/",
		},
		{
			dirPath: "examples/usecases-01-copy-tfvars/dev/eastus/project1/",
		},
		{
			dirPath: "examples/usecases-01-copy-tfvars/dev/eastus/project2/",
		},
		{
			dirPath: "examples/usecases-02-reference-remote-state/dev/eastus/project1/",
		},
		{
			dirPath: "examples/usecases-02-reference-remote-state/dev/eastus/project2/",
		},
		{
			dirPath: "examples/usecases-03-import-local-module/dev/eastus/project1/",
		},
		{
			dirPath: "examples/usecases-03-import-local-module/dev/eastus/project2/",
		},
		{
			dirPath: "examples/usecases-04-import-remote-module/dev/eastus/project1/",
		},
		{
			dirPath: "examples/usecases-04-import-remote-module/dev/eastus/project2/",
		},
		{
			dirPath: "examples/usecases-05-terraform-versions/dev/eastus/project1/",
		},
		{
			dirPath: "examples/usecases-05-terraform-versions/dev/eastus/project2/",
		},
		{
			dirPath: "fixtures/valid/empty/",
		},
		{
			dirPath: "fixtures/valid/empty-tfvars/",
		},
		{
			dirPath: "fixtures/valid/empty-tfvars/",
		},
		{
			dirPath: "fixtures/valid/duplicate-includes/",
		},
		{
			dirPath: "fixtures/valid/include-file-comments/",
		},
		{
			dirPath: "fixtures/valid/locals-referencing-locals/",
		},
		// TODO test variable reference hell 2
	} {
		t.Run(fixture.dirPath, func(t *testing.T) {
			args := []string{"generate", "--file", path.Join(fixture.dirPath, "tf-generator.hcl"), "--check"}
			if err := run(args); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestInvalidFixtures(t *testing.T) {
	for _, fixture := range []InvalidFixture{
		{
			dirPath:                 "fixtures/invalid/unknown-block/",
			expectedMessageContains: "Unsupported block type; Blocks of type \"unknown\" are not expected here.",
			isDiag:                  true,
		},
		{
			dirPath:                 "fixtures/invalid/syntax-error/",
			expectedMessageContains: "Unclosed configuration block; There is no closing brace for this block before the end of the file.",
			isDiag:                  true,
		},
		{
			dirPath:                 "fixtures/invalid/unknown/",
			expectedMessageContains: "<nil>: Configuration file not found; The configuration file fixtures/invalid/unknown/tf-generator.hcl does not exist.",
			isDiag:                  true,
		},
		{
			dirPath:                 "fixtures/invalid/tfvars-does-not-exist/",
			expectedMessageContains: "no such file or directory.",
			isDiag:                  true,
		},
		{
			dirPath:                 "fixtures/invalid/empty-generated-tfvars/",
			expectedMessageContains: "the new tfvars file does not match the existing file.",
			isDiag:                  false,
		},
		{
			dirPath:                 "fixtures/invalid/invalid-tfvars/",
			expectedMessageContains: "parse config: [locals.tfvars:1,9-10: Unclosed configuration block; There is no closing brace for this block",
			isDiag:                  true,
		},
		{
			dirPath:                 "fixtures/invalid/injected-tfvar-does-not-exist/",
			expectedMessageContains: "Could not find tfvar for ref `injectvar.does-not-exist`",
			isDiag:                  true,
		},
		{
			dirPath:                 "fixtures/invalid/duplicate-local/",
			expectedMessageContains: `tf-generator.hcl:6,3-8: local "a" already defined`,
			isDiag:                  true,
		},
		// TODO test injected tfvar does not exist 2
	} {
		t.Run(fixture.dirPath, func(t *testing.T) {
			filePath := path.Join(fixture.dirPath, "tf-generator.hcl")
			args := []string{"generate", "--file", filePath, "--check"}
			err := run(args)
			assert.NotNilf(t, err, "expected error")
			assert.Contains(t, err.Error(), fixture.expectedMessageContains)
			_, isDiag := err.(hcl.Diagnostics)
			assert.Equal(t, fixture.isDiag, isDiag)
		})
	}
}
