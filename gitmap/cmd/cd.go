package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runCD handles the "cd" subcommand routing.
func runCD(args []string) {
	checkHelp("cd", args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrCDUsage)
		os.Exit(1)
	}

	sub := args[0]
	rest := args[1:]

	routeCDSub(sub, rest)
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
