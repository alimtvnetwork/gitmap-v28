package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runAlias handles the "alias" subcommand for machine identity aliasing and repo aliases.
func runAlias(args []string) *apperror.AppError {
	if !isMachineAliasInvocation(args) {
		return dispatchAlias(args[0], args[1:])
	}
	if err := cmdos.RunAliasCLI(args); err != nil {
		return apperror.New(err.Error(), "E9000", nil)
	}
	return nil
}

func isMachineAliasInvocation(args []string) bool {
	if len(args) == 0 {
		return true
	}
	for _, a := range args {
		if a == "--ssh" || a == "-s" || a == "change" || a == "revert" || a == "help" || a == "-h" || a == "--help" {
			return true
		}
	}
	first := args[0]
	return first == "ls" || first == "set"
}

// dispatchAlias routes legacy repository alias subcommands to their handlers.
func dispatchAlias(sub string, args []string) *apperror.AppError {
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
