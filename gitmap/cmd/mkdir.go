package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func runMkdir(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gitmap mkdir [-p] <path>")
		os.Exit(1)
	}

	createParents := false
	pathArg := ""

	if args[0] == "-p" {
		createParents = true
		if len(args) > 1 {
			pathArg = args[1]
		}
	} else {
		pathArg = args[0]
	}

	if pathArg == "" {
		fmt.Fprintln(os.Stderr, "Error: missing path argument")
		os.Exit(1)
	}

	absPath, err := filepath.Abs(pathArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	if createParents {
		err = os.MkdirAll(absPath, 0755)
	} else {
		err = os.Mkdir(absPath, 0755)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created directory: %s\n", absPath)
	return nil
}
