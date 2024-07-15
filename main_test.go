package main

import (
	"os"
	"strings"
	"testing"
	"tf-generator/initialize"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

type ValidFixture struct {
	args []string
}

type InvalidFixture struct {
	args                    []string
	expectedMessageContains string
	isDiag                  bool
}

func TestValidFixtures(t *testing.T) {
	for _, fixture := range []ValidFixture{
		{
			args: []string{"--file", "examples/basics-01-load-and-combine/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "examples/basics-02-locals/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "examples/basics-03-merge-tfvars/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "examples/basics-04-remove-tfvar-keys/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "examples/basics-05-combine-with-inject/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "fixtures/valid/empty/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "fixtures/valid/empty-tfvars/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "fixtures/valid/duplicate-includes/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "fixtures/valid/include-file-comments/tf-generator.hcl"},
		},
		{
			args: []string{"--file", "fixtures/valid/locals-referencing-locals/tf-generator.hcl"},
		},
		{
			args: []string{"--glob", "fixtures/valid/*/tf-generator.hcl"},
		},
	} {
		t.Run(strings.Join(fixture.args, " "), func(t *testing.T) {
			args := append([]string{"generate", "--check"}, fixture.args...)
			if err := run(args); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestInvalidFixtures(t *testing.T) {
	for _, fixture := range []InvalidFixture{
		{
			args:                    []string{"--file", "fixtures/invalid/unknown-block/tf-generator.hcl"},
			expectedMessageContains: "Unsupported block type; Blocks of type \"unknown\" are not expected here.",
			isDiag:                  true,
		},
		{
			args:                    []string{"--file", "fixtures/invalid/syntax-error/tf-generator.hcl"},
			expectedMessageContains: "Unclosed configuration block; There is no closing brace for this block before the end of the file.",
			isDiag:                  true,
		},
		{
			args:                    []string{"--file", "fixtures/invalid/unknown/tf-generator.hcl"},
			expectedMessageContains: "<nil>: Configuration file not found; The configuration file fixtures/invalid/unknown/tf-generator.hcl does not exist.",
			isDiag:                  true,
		},
		{
			args:                    []string{"--file", "fixtures/invalid/tfvars-does-not-exist/tf-generator.hcl"},
			expectedMessageContains: "no such file or directory.",
			isDiag:                  true,
		},
		{
			args:                    []string{"--file", "fixtures/invalid/empty-generated-tfvars/tf-generator.hcl"},
			expectedMessageContains: "the new tfvars file does not match the existing file.",
			isDiag:                  false,
		},
		{
			args:                    []string{"--file", "fixtures/invalid/invalid-tfvars/tf-generator.hcl"},
			expectedMessageContains: "parse config: [locals.tfvars:1,9-10: Unclosed configuration block; There is no closing brace for this block",
			isDiag:                  true,
		},
		{
			args:                    []string{"--file", "fixtures/invalid/injected-tfvar-does-not-exist/tf-generator.hcl"},
			expectedMessageContains: "Could not find tfvar for ref `injectvar.does-not-exist`",
			isDiag:                  true,
		},
		{
			args:                    []string{"--file", "fixtures/invalid/duplicate-local/tf-generator.hcl"},
			expectedMessageContains: `tf-generator.hcl:6,3-8: local "a" already defined`,
			isDiag:                  true,
		},
		{
			args:                    []string{"--glob", "fixtures/invalid/*/tf-generator.hcl"},
			expectedMessageContains: `local "a" already defined`,
			isDiag:                  true,
		},
		{
			args:                    []string{"--glob", "fixtures/**/tf-generator.hcl"},
			expectedMessageContains: `local "a" already defined`,
			isDiag:                  true,
		},
		{
			args:                    []string{"--file", "test", "--glob", "test"},
			expectedMessageContains: "cannot specify --file and --glob",
			isDiag:                  false,
		},
		{
			args:                    []string{},
			expectedMessageContains: "The configuration file tf-generator.hcl does not exist.",
			isDiag:                  true,
		},
	} {
		t.Run(strings.Join(fixture.args, " "), func(t *testing.T) {
			args := append([]string{"generate", "--check"}, fixture.args...)
			err := run(args)
			assert.NotNilf(t, err, "expected error")
			errMsg := err.Error()
			assert.Contains(t, errMsg, fixture.expectedMessageContains)
			_, isDiag := err.(hcl.Diagnostics)
			assert.Equal(t, fixture.isDiag, isDiag)
		})
	}
}

func cleanInitFiles() {
	for _, initFile := range initialize.InitFiles {
		filePath := initFile.FileName
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			os.Remove(filePath)
		}
	}
}

func initTestWrapper(testFunc func()) {
	cleanInitFiles()
	testFunc()
	cleanInitFiles()
}

func TestInit(t *testing.T) {
	initTestWrapper(func() {
		args := []string{"init"}
		if err := run(args); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitFilesAlreadyExist(t *testing.T) {
	initTestWrapper(func() {
		args := []string{"init"}
		if err := run(args); err != nil {
			t.Fatal(err)
		}
		err := run(args)
		assert.NotNilf(t, err, "expected error")
		assert.Equal(t, err.Error(), "file 'tf-generator.hcl' already exists")
	})
}
