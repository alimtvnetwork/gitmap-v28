package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runAlias handles the "alias" subcommand and routes to sub-handlers.
func runAlias(args []string) *apperror.AppError {
	checkHelp("alias", args)
	if len(args) == 0 {
		runAliasList()

		return nil
	}
	return dispatchAlias(args[0], args[1:])
}

// dispatchAlias routes alias subcommands to their handlers.
func dispatchAlias(sub string, args []string) *apperror.AppError {
	if sub == constants.SubCmdAliasSet {
		runAliasSet(args)
		return nil
	}
	if sub == constants.SubCmdAliasRm {
		runAliasRemove(args)
		return nil
	}
	if sub == constants.SubCmdAliasList {
		runAliasList()
		return nil
	}
	if sub == constants.SubCmdAliasShow {
		runAliasShow(args)
		return nil
	}
	if sub == constants.SubCmdAliasSug {
		runAliasSuggest(args)
		return nil
	}

	return apperror.New(fmt.Sprintf(constants.ErrUnknownCommand, sub), "E9000", nil)
}
