package cli

import (
	"flag"
	"fmt"
	"tf-generator/initialize"
)

type InitCommand struct {
	fs *flag.FlagSet
}

// NewInitCommand sub-command to init project
func NewInitCommand() *InitCommand {
	c := &InitCommand{
		fs: flag.NewFlagSet("init", flag.ContinueOnError),
	}

	return c
}

func (c *InitCommand) Name() string {
	return c.fs.Name()
}

func (c *InitCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	if len(c.fs.Args()) != 0 {
		return fmt.Errorf("expected no positional arguments, got %d", len(c.fs.Args()))
	}

	return nil
}

func (c *InitCommand) Run() error {
	return initialize.Run()
}
