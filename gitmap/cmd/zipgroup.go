package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runZipGroup handles the "zip-group" subcommand and routes to sub-handlers.
func runZipGroup(args []string) {
	checkHelp("zip-group", args)
	if len(args) == 0 {
		runZipGroupList()

		return
	}
	dispatchZipGroup(args[0], args[1:])
}

// dispatchZipGroup routes zip-group subcommands to their handlers.
func dispatchZipGroup(sub string, args []string) {
	if sub == constants.SubCmdZGCreate {
		runZipGroupCreate(args)

		return
	}
	if sub == constants.SubCmdZGAdd {
		runZipGroupAdd(args)

		return
	}
	if sub == constants.SubCmdZGRemove {
		runZipGroupRemove(args)

		return
	}
	if sub == constants.SubCmdZGList {
		runZipGroupList()

		return
	}
	if sub == constants.SubCmdZGShow {
		runZipGroupShow(args)

		return
	}
	if sub == constants.SubCmdZGDelete {
		runZipGroupDelete(args)

		return
	}
	if sub == constants.SubCmdZGRename {
		runZipGroupRename(args)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrUnknownCommand, sub)
	os.Exit(1)
}
