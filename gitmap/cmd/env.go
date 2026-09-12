package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

func routeEnvVariableSub(sub string, args []string) (error, bool) {
	if sub == constants.CmdEnvSet {
		return runEnvSet(args), true
	}
	if sub == constants.CmdEnvGet {
		return runEnvGet(args), true
	}
	if sub == constants.CmdEnvDelete {
		return runEnvDelete(args), true
	}
	if sub == constants.CmdEnvList {
		return runEnvList(), true
	}

	return nil, false
}

// routeEnvSub routes to the appropriate env subcommand.
func routeEnvSub(sub string, args []string) error {
	if err, isHandled := routeEnvVariableSub(sub, args); isHandled {
		return err
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
