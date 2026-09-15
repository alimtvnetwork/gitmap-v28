package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runCD handles the "cd" subcommand routing.
func runCD(args []string) error {
	checkHelp("cd", args)
	if len(args) == 0 {
		return handleBareCD()
	}

	sub := args[0]
	rest := args[1:]

	return routeCDSub(sub, rest).AsError()
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
func routeCDSub(sub string, args []string) result.ErrorWrapper {
	if sub == constants.CmdCDRepos {
		return result.MatchWrapper(runCDRepos(args))
	}

	if sub == constants.CmdCDSetDefault {
		return result.MatchWrapper(runCDSetDefault(args))
	}

	if sub == constants.CmdCDClearDefault {
		return result.MatchWrapper(runCDClearDefault(args))
	}

	return result.MatchWrapper(runCDLookup(sub, args))
}
