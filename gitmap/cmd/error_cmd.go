package cmd

import (
	"fmt"
	"os"
)

func runErrorCmd(args []string) error {
	checkHelp("error", args)

	if len(args) == 0 {
		printErrorUsage()
		return nil
	}

	sub := args[0]
	subArgs := args[1:]

	if sub == "export" {
		return runErrorExport(subArgs)
	}

	printErrorUsage()
	return nil
}

func runErrorExport(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("error export: missing destination file")
	}
	dest := args[0]

	lastErrFile := ".gitmap/last_error.log"
	data, err := os.ReadFile(lastErrFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("error export: no recent error found to export")
		}
		return fmt.Errorf("error export: could not read last error: %w", err)
	}

	if err := os.WriteFile(dest, data, 0644); err != nil {
		return fmt.Errorf("error export: could not write to %s: %w", dest, err)
	}

	fmt.Printf("gitmap error: exported last error to %s\n", dest)
	return nil
}

func printErrorUsage() {
	fmt.Println(`Usage: gitmap error <command> [arguments]

Commands:
  export <file>    Export the last failed command error to a file

Examples:
  gitmap error export error.txt`)
}
