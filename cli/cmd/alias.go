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
	if first == "set" && len(args) >= 3 {
		return false
	}

	return first == "ls" || first == "set"
}

// dispatchAlias routes legacy repository alias subcommands to their handlers.
func dispatchAlias(sub string, args []string) *apperror.AppError {
	if sub == constants.SubCmdAliasSet {
		return dispatchAliasSet(args)
	}
	if sub == constants.SubCmdAliasRm {
		return dispatchAliasRemove(args)
	}
	if sub == constants.SubCmdAliasList {
		return dispatchAliasList()
	}
	if sub == constants.SubCmdAliasShow {
		return dispatchAliasShow(args)
	}
	if sub == constants.SubCmdAliasSug {
		return dispatchAliasSuggest(args)
	}

	return apperror.New(fmt.Sprintf(constants.ErrUnknownCommand, sub), "E9000", nil)
}

func dispatchAliasSet(args []string) *apperror.AppError {
	if err := runAliasSet(args); err != nil {
		return apperror.WrapSimple(err, "alias set")
	}

	return nil
}

func dispatchAliasRemove(args []string) *apperror.AppError {
	if err := runAliasRemove(args); err != nil {
		return apperror.WrapSimple(err, "alias remove")
	}

	return nil
}

func dispatchAliasList() *apperror.AppError {
	if err := runAliasList(); err != nil {
		return apperror.WrapSimple(err, "alias list")
	}

	return nil
}

func dispatchAliasShow(args []string) *apperror.AppError {
	if err := runAliasShow(args); err != nil {
		return apperror.WrapSimple(err, "alias show")
	}

	return nil
}

func dispatchAliasSuggest(args []string) *apperror.AppError {
	if err := runAliasSuggest(args); err != nil {
		return apperror.WrapSimple(err, "alias suggest")
	}

	return nil
}
