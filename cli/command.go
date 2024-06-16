package cli

import "fmt"

// Command interface for all subcommands
type Command interface {
	Init([]string) error
	Run() error
	Name() string
}

func RunSubcommand(args []string, cmds []Command) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand")
	}

	// Run specified subcommand
	subcommand := args[0]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			args := args[1:]
			if err := cmd.Init(args); err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}
