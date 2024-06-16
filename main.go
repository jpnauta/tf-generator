package main

import (
	"fmt"
	"os"
	"tf-generator/cli"
)

func run(args []string) error {
	return cli.Run(args)
}

// main entry point of the application
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
