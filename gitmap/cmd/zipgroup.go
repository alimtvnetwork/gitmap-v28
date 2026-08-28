package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runZipGroup handles the "zip-group" subcommand and routes to sub-handlers.
func runZipGroup(args []string) error {
	checkHelp("zip-group", args)
	if len(args) == 0 {
		runZipGroupList()

		return nil
	}
	dispatchZipGroup(args[0], args[1:])
	return nil
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

	return apperror.New(constants.ErrUnknownCommand, "E9000", nil)
}
