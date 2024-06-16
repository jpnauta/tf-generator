package cli

// Run main entry point for the application
func Run(args []string) error {
	return RunSubcommand(args, []Command{
		NewGenerateCommand(),
	})
}
