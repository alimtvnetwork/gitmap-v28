package cmd

import (
	"fmt"
	"io"
	"os"
)

func runCat(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gitmap cat <file>")
		os.Exit(1)
	}

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
}
