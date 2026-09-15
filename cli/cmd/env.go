package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runEnv handles the "env" subcommand routing.
func runEnv(args []string) error {
	checkHelp("env", args)
	if len(args) < 1 {
		return apperror.NewSimple("constants.ErrEnvSubcommand "+"", "E9000")
	}

	sub := args[0]
	rest := args[1:]

	return routeEnvSub(sub, rest)
}

func routeEnvVariableSub(sub string, args []string) result.Result[bool] {
	if sub == constants.CmdEnvSet {
		return result.RouteMatched(runEnvSet(args))
	}

	if sub == constants.CmdEnvGet {
		return result.RouteMatched(runEnvGet(args))
	}

	if sub == constants.CmdEnvDelete {
		return result.RouteMatched(runEnvDelete(args))
	}

	if sub == constants.CmdEnvList {
		return result.RouteMatched(runEnvList())
	}

	return result.RouteUnmatched()
}

// routeEnvSub routes to the appropriate env subcommand.
func routeEnvSub(sub string, args []string) error {
	resVar := routeEnvVariableSub(sub, args)
	if resVar.Data {
		return resVar.AppError()
	}

	if sub == constants.CmdEnvPathAdd {
		return routeEnvPath(args)
	}

	return apperror.NewSimple(constants.ErrEnvSubcommand, "E9000")
}

// routeEnvPath routes path subcommands (path add, path remove, path list).
func routeEnvPath(args []string) error {
	if len(args) < 1 {
		return runEnvPathList()
	}

	sub := args[0]
	rest := args[1:]

	return dispatchEnvPath(sub, rest)
}

func dispatchEnvPath(sub string, rest []string) error {
	if sub == constants.CmdEnvPathSub {
		return runEnvPathAdd(rest)
	}

	if sub == constants.CmdEnvPathRemove {
		return runEnvPathRemove(rest)
	}

	if sub == constants.CmdEnvPathList {
		return runEnvPathList()
	}

	return apperror.NewSimple("constants.ErrEnvSubcommand "+"path "+sub, "E9000")
}
