package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runEnv handles the "env" subcommand routing.
func runEnv(args []string) error {
	checkHelp("env", args)
	if len(args) < 1 || isHelpSubcommand(args[0]) {
		printEnvHelp()

		return nil
	}

	return result.AsError(routeEnvSub(args[0], args[1:]))
}

func printEnvHelp() {
	_, mode := ParsePrettyFlag(nil)
	helptext.PrintWithMode("env", mode)
	printUsageFooterShort()
}

func isHelpSubcommand(sub string) bool {
	return sub == "help" || sub == "-h" || sub == "--help"
}

func routeEnvVariableSub(sub string, args []string) result.ErrorWrapper {
	if isEnvSetSub(sub) {
		return result.MatchWrapper(runEnvSet(args))
	}

	if sub == constants.CmdEnvGet {
		return result.MatchWrapper(runEnvGet(args))
	}

	if isEnvDeleteSub(sub) {
		return result.MatchWrapper(runEnvDelete(args))
	}

	if isEnvListSub(sub) {
		return result.MatchWrapper(runEnvList())
	}

	return result.UnmatchedWrapper()
}

func isEnvSetSub(sub string) bool {
	return sub == constants.CmdEnvSet || sub == "add"
}

func isEnvDeleteSub(sub string) bool {
	return sub == constants.CmdEnvDelete || sub == "rm" || sub == "remove" || sub == "del"
}

func isEnvListSub(sub string) bool {
	return sub == constants.CmdEnvList || sub == "ls"
}

// routeEnvSub routes to the appropriate env subcommand.
func routeEnvSub(sub string, args []string) result.ErrorWrapper {
	resVar := routeEnvVariableSub(sub, args)
	if resVar.IsMatched() {
		return resVar
	}

	if sub == constants.CmdEnvPathAdd {
		return routeEnvPath(args)
	}

	return result.FailureWrapper(apperror.NewSimple(constants.ErrEnvSubcommand, "E9000"))
}

// routeEnvPath routes path subcommands (path add, path remove, path list).
func routeEnvPath(args []string) result.ErrorWrapper {
	if len(args) < 1 {
		return result.MatchWrapper(runEnvPathList())
	}

	sub := args[0]
	rest := args[1:]

	return dispatchEnvPath(sub, rest)
}

func dispatchEnvPath(sub string, rest []string) result.ErrorWrapper {
	if sub == constants.CmdEnvPathSub {
		return result.MatchWrapper(runEnvPathAdd(rest))
	}

	if sub == constants.CmdEnvPathRemove {
		return result.MatchWrapper(runEnvPathRemove(rest))
	}

	if sub == constants.CmdEnvPathList {
		return result.MatchWrapper(runEnvPathList())
	}

	return result.FailureWrapper(apperror.NewSimple("constants.ErrEnvSubcommand "+"path "+sub, "E9000"))
}
