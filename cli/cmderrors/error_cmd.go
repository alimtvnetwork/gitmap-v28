package cmderrors

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func RunErrorCmd(args []string) error {
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
		return apperror.NewValidationError("error export: missing destination file")
	}

	dest := args[0]

	lastErrFile := ".gitmap/last_error.log"
	data, err := os.ReadFile(lastErrFile)
	if os.IsNotExist(err) {
		return apperror.NewNotFoundError("error export: no recent error found to export")
	}

	if err != nil {
		return apperror.WrapSimple(err, "error export: could not read last error")
	}

	if err := os.WriteFile(dest, data, 0644); err != nil {
		return apperror.WrapSimple(err, "error export: could not write to "+dest)
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
