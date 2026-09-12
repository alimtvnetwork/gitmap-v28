package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runCD handles the "cd" subcommand routing.
func runCD(args []string) error {
	checkHelp("cd", args)
	if len(args) == 0 {
		return handleBareCD()
	}

	sub := args[0]
	rest := args[1:]

	return routeCDSub(sub, rest)
}

func handleBareCD() error {
	workPath, hasDefault := resolveDefaultWorkDirPath()
	if hasDefault {
		fmt.Print(workPath)
		WriteShellHandoff(workPath)
		warnIfNoWrapper()

		return nil
	}

	fmt.Fprint(os.Stderr, constants.ErrCDUsage)

	return apperror.NewValidationError("repo name or work directory required")
}

// routeCDSub routes to the appropriate cd handler.
func routeCDSub(sub string, args []string) error {
	if sub == constants.CmdCDRepos {
		return runCDRepos(args)
	}

	if sub == constants.CmdCDSetDefault {
		return runCDSetDefault(args)
	}

	if sub == constants.CmdCDClearDefault {
		return runCDClearDefault(args)
	}

	return runCDLookup(sub, args)
}
