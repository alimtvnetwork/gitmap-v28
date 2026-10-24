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
		if isMachineAliasFlag(a) {
			return true
		}
	}

	first := constants.AliasActionType(args[0])
	if first.IsSet() && len(args) >= 3 {
		return false
	}

	return first.IsList() || first.IsSet()
}

func isMachineAliasFlag(a string) bool {
	return a == "--ssh" || a == "-s" || a == "change" || a == "revert" || isHelpFlag(a)
}

// dispatchAlias routes legacy repository alias subcommands to their handlers.
func dispatchAlias(sub string, args []string) *apperror.AppError {
	action := constants.AliasActionType(sub)
	if action.IsSet() {
		return dispatchAliasSet(args)
	}

	if action.IsRemove() {
		return dispatchAliasRemove(args)
	}

	if action.IsList() {
		return dispatchAliasList()
	}

	if action.IsShow() {
		return dispatchAliasShow(args)
	}

	if action.IsSuggest() {
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
