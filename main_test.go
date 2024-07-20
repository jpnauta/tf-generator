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
	cwd  string
	args []string
}

type InvalidFixture struct {
	cwd                     string
	args                    []string
	expectedMessageContains string
	isDiag                  bool
}

func withCwd(t *testing.T, newDir string, fn func()) {
	// Save current working directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current working directory: %v", err)
	}

	// Change cwd temporarily
	if err := os.Chdir(newDir); err != nil {
		t.Fatalf("error changing directory to %s: %v", newDir, err)
	}

	fn()

	// Restore original cwd
	if err := os.Chdir(oldDir); err != nil {
		t.Fatalf("error restoring original directory: %v", err)
	}
}

func TestValidFixtures(t *testing.T) {
	for _, fixture := range []ValidFixture{
		{
			cwd:  ".",
			args: []string{"--file", "examples/basics-01-load-and-combine/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "examples/basics-02-locals/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "examples/basics-03-merge-tfvars/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "examples/basics-04-remove-tfvar-keys/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "examples/basics-05-combine-with-inject/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "fixtures/valid/empty/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "fixtures/valid/empty-tfvars/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "fixtures/valid/duplicate-includes/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "fixtures/valid/include-file-comments/tf-generator.hcl"},
		},
		{
			cwd:  ".",
			args: []string{"--file", "fixtures/valid/locals-referencing-locals/tf-generator.hcl"},
		},
		{
			cwd:  "fixtures/valid/",
			args: []string{"--recursive"},
		},
	} {
		t.Run(strings.Join(fixture.args, " "), func(t *testing.T) {
			withCwd(t, fixture.cwd, func() {
				args := append([]string{"generate", "--check"}, fixture.args...)
				if err := run(args); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}

func TestInvalidFixtures(t *testing.T) {
	for _, fixture := range []InvalidFixture{
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/unknown-block/tf-generator.hcl"},
			expectedMessageContains: "Unsupported block type; Blocks of type \"unknown\" are not expected here.",
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/syntax-error/tf-generator.hcl"},
			expectedMessageContains: "Unclosed configuration block; There is no closing brace for this block before the end of the file.",
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/unknown/tf-generator.hcl"},
			expectedMessageContains: "<nil>: Configuration file not found; The configuration file fixtures/invalid/unknown/tf-generator.hcl does not exist.",
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/tfvars-does-not-exist/tf-generator.hcl"},
			expectedMessageContains: "no such file or directory.",
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/empty-generated-tfvars/tf-generator.hcl"},
			expectedMessageContains: "the new tfvars file does not match the existing file.",
			isDiag:                  false,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/invalid-tfvars/tf-generator.hcl"},
			expectedMessageContains: "parse config: [locals.tfvars:1,9-10: Unclosed configuration block; There is no closing brace for this block",
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/injected-tfvar-does-not-exist/tf-generator.hcl"},
			expectedMessageContains: "Could not find tfvar for ref `injectvar.does-not-exist`",
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "fixtures/invalid/duplicate-local/tf-generator.hcl"},
			expectedMessageContains: `tf-generator.hcl:6,3-8: local "a" already defined`,
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{"--file", "a/tf-generator.hcl", "--recursive"},
			expectedMessageContains: `file name a/tf-generator.hcl cannot contain folders`,
			isDiag:                  false,
		},
		{
			cwd:                     "fixtures/invalid/",
			args:                    []string{"--recursive"},
			expectedMessageContains: `local "a" already defined`,
			isDiag:                  true,
		},
		{
			cwd:                     "fixtures/invalid/nested-folder-with-error",
			args:                    []string{"--recursive"},
			expectedMessageContains: "parse config: [locals.tfvars:1,9-10: Unclosed configuration block; There is no closing brace for this block",
			isDiag:                  true,
		},
		{
			cwd:                     "fixtures/",
			args:                    []string{"--recursive"},
			expectedMessageContains: `local "a" already defined`,
			isDiag:                  true,
		},
		{
			cwd:                     ".",
			args:                    []string{},
			expectedMessageContains: "The configuration file tf-generator.hcl does not exist.",
			isDiag:                  true,
		},
	} {
		t.Run(strings.Join(fixture.args, " "), func(t *testing.T) {
			withCwd(t, fixture.cwd, func() {
				args := append([]string{"generate", "--check"}, fixture.args...)
				err := run(args)
				assert.NotNilf(t, err, "expected error")
				errMsg := err.Error()
				assert.Contains(t, errMsg, fixture.expectedMessageContains)
				_, isDiag := err.(hcl.Diagnostics)
				assert.Equal(t, fixture.isDiag, isDiag)
			})
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
