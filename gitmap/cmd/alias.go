package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runAlias handles the "alias" subcommand and routes to sub-handlers.
func runAlias(args []string) {
	checkHelp("alias", args)
	if len(args) == 0 {
		runAliasList()

		return
	}
	dispatchAlias(args[0], args[1:])
}

// dispatchAlias routes alias subcommands to their handlers.
func dispatchAlias(sub string, args []string) {
	if sub == constants.SubCmdAliasSet {
		runAliasSet(args)

		return
	}
	if sub == constants.SubCmdAliasRm {
		runAliasRemove(args)

		return
	}
	if sub == constants.SubCmdAliasList {
		runAliasList()

		return
	}
	if sub == constants.SubCmdAliasShow {
		runAliasShow(args)

		return
	}
	if sub == constants.SubCmdAliasSug {
		runAliasSuggest(args)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrUnknownCommand, sub)
	os.Exit(1)
}
