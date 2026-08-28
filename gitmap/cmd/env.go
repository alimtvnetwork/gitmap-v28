package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runEnv handles the "env" subcommand routing.
func runEnv(args []string) error {
	checkHelp("env", args)
	if len(args) < 1 {
		return apperror.Wrap("", "constants.ErrEnvSubcommand", nil)
	}

	sub := args[0]
	rest := args[1:]

	routeEnvSub(sub, rest)
	return nil
}

// routeEnvSub routes to the appropriate env subcommand.
func routeEnvSub(sub string, args []string) {
	if sub == constants.CmdEnvSet {
		runEnvSet(args)

		return
	}
	if sub == constants.CmdEnvGet {
		runEnvGet(args)

		return
	}
	if sub == constants.CmdEnvDelete {
		runEnvDelete(args)

		return
	}
	if sub == constants.CmdEnvList {
		runEnvList()

		return
	}
	if sub == constants.CmdEnvPathAdd {
		routeEnvPath(args)

		return
	}

	return apperror.New(constants.ErrEnvSubcommand, "E9000", nil)
}

// routeEnvPath routes path subcommands (path add, path remove, path list).
func routeEnvPath(args []string) {
	if len(args) < 1 {
		runEnvPathList()

		return
	}

	sub := args[0]
	rest := args[1:]

	if sub == constants.CmdEnvPathSub {
		runEnvPathAdd(rest)

		return
	}
	if sub == constants.CmdEnvPathRemove {
		runEnvPathRemove(rest)

		return
	}
	if sub == constants.CmdEnvPathList {
		runEnvPathList()

		return
	}

	return apperror.Wrap("path "+sub, "constants.ErrEnvSubcommand", nil)
}
