package cli

import (
	"flag"
	"fmt"
	"tf-generator/generate"
)

type GenerateCommand struct {
	fs *flag.FlagSet

	file  string
	check bool
}

// NewGenerateCommand sub-command to generate files
func NewGenerateCommand() *GenerateCommand {
	c := &GenerateCommand{
		fs: flag.NewFlagSet("generate", flag.ContinueOnError),
	}

	c.fs.StringVar(&c.file, "file", "tf-generator.hcl", "file used to configure file generation")
	c.fs.BoolVar(&c.check, "check", false, "only check if file is up-to-date, do not update it")

	return c
}

func (c *GenerateCommand) Name() string {
	return c.fs.Name()
}

func (c *GenerateCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	if len(c.fs.Args()) != 0 {
		return fmt.Errorf("expected no positional arguments, got %d", len(c.fs.Args()))
	}
	return nil
}

func (c *GenerateCommand) Run() error {
	return generate.Run(c.file, c.check)
}
