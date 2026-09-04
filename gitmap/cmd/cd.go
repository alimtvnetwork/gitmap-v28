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

	routeCDSub(sub, rest)

	return nil
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
func routeCDSub(sub string, args []string) {
	if sub == constants.CmdCDRepos {
		runCDRepos(args)

		return
	}
	if sub == constants.CmdCDSetDefault {
		runCDSetDefault(args)

		return
	}
	if sub == constants.CmdCDClearDefault {
		runCDClearDefault(args)

		return
	}

	runCDLookup(sub, args)
}
