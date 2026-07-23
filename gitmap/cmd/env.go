package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runEnv handles the "env" subcommand routing.
func runEnv(args []string) {
	checkHelp("env", args)
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, constants.ErrEnvSubcommand, "")
		os.Exit(1)
	}

	sub := args[0]
	rest := args[1:]

	routeEnvSub(sub, rest)
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

	fmt.Fprintf(os.Stderr, constants.ErrEnvSubcommand, sub)
	os.Exit(1)
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

	fmt.Fprintf(os.Stderr, constants.ErrEnvSubcommand, "path "+sub)
	os.Exit(1)
}
