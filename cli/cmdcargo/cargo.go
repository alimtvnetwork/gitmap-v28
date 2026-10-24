package cmdcargo

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	cargoconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
)

// CheckHelpFn delegates help checking.
var CheckHelpFn func(command string, args []string)

// RunCargo is the main entry point for the gitmap cargo / crg command.
func RunCargo(args []string) *apperror.AppError {
	checkHelp(cargoconst.CmdCargo, args)
	if len(args) == 0 {
		return runCargoDefault()
	}

	sub := strings.ToLower(args[0])
	if isStatusSubcommand(sub) {
		return runCargoStatus()
	}
	if isInstallSubcommand(sub) {
		return runCargoInstall(args[1:])
	}

	return dispatchCargoExecution(args)
}

func isStatusSubcommand(sub string) bool {
	sc, err := enums.ParseCargoSubcommandType(sub)
	if err != nil {
		return false
	}

	return sc.IsStatus()
}

func isInstallSubcommand(sub string) bool {
	sc, err := enums.ParseCargoSubcommandType(sub)
	if err != nil {
		return false
	}

	return sc.IsInstall()
}

func runCargoDefault() *apperror.AppError {
	cargoBin := ResolveCargoBinary()
	if cargoBin == "" {
		PrintCargoInstallSuggestions()

		return NewMissingCargoError()
	}

	return execCargo(cargoBin, []string{cargoconst.FlagLongHelp})
}

func dispatchCargoExecution(args []string) *apperror.AppError {
	cleanArgs, hasInstall := extractInstallFlag(args)
	cargoBin := ResolveCargoBinary()
	if cargoBin == "" && hasInstall {
		_ = runCargoInstall([]string{cargoconst.FlagLongYes})
		cargoBin = ResolveCargoBinary()
	}

	if cargoBin == "" {
		PrintCargoInstallSuggestions()

		return NewMissingCargoError()
	}

	return execCargo(cargoBin, cleanArgs)
}

func runCargoInstall(rest []string) *apperror.AppError {
	if err := cmdinstall.RunInstall(append([]string{cargoconst.BinCargo}, rest...)); err != nil {
		return apperror.Wrap(err, "runCargoInstall", map[string]any{"tool": cargoconst.BinCargo})
	}

	return nil
}

func checkHelp(cmd string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(cmd, args)
	}
}
