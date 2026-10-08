package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsetup"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runCD handles the "cd" subcommand routing.
func runCD(args []string) error {
	checkHelp("cd", args)
	if len(args) == 0 {
		return handleBareCD()
	}

	return result.AsError(routeCDSub(args[0], args[1:]))
}

func handleBareCD() error {
	workPath, hasDefault := resolveDefaultWorkDirPath()
	if hasDefault {
		fmt.Print(workPath)
		WriteShellHandoff(workPath)
		cmdsetup.WarnIfNoWrapper()

		return nil
	}

	fmt.Fprint(os.Stderr, constants.ErrCDUsage)

	return apperror.NewValidationError("repo name or work directory required")
}

// routeCDSub routes to the appropriate cd handler.
func routeCDSub(sub string, args []string) result.ErrorWrapper {
	if isConfigCDSub(sub) {
		return result.MatchWrapper(routeConfigCDSub(sub, args))
	}
	if isSpecialRepoCDAlias(sub) {
		return result.MatchWrapper(runCDSpecialRepo(sub, args))
	}
	return result.MatchWrapper(runCDLookup(sub, args))
}

func isConfigCDSub(sub string) bool {
	return sub == constants.CmdCDRepos || sub == constants.CmdCDSetDefault || sub == constants.CmdCDClearDefault
}

func routeConfigCDSub(sub string, args []string) error {
	if sub == constants.CmdCDRepos {
		return runCDRepos(args)
	}
	if sub == constants.CmdCDSetDefault {
		return runCDSetDefault(args)
	}
	return runCDClearDefault(args)
}

func isSpecialRepoCDAlias(sub string) bool {
	switch sub {
	case "rs", "repo-secrets", "rc", "repo-cache", "repo-storage":
		return true
	default:
		return false
	}
}
