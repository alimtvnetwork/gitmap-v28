package cmdcargo

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CheckHelpFn delegates help checking.
var CheckHelpFn func(command string, args []string)

// RunCargo is the main entry point for the gitmap cargo / crg command.
func RunCargo(args []string) error {
	checkHelp("cargo", args)
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
	return sub == "status" || sub == "st" || sub == "info"
}

func isInstallSubcommand(sub string) bool {
	return sub == "install" || sub == "in"
}

func runCargoDefault() error {
	cargoBin := ResolveCargoBinary()
	if cargoBin == "" {
		PrintCargoInstallSuggestions()
		return NewMissingCargoError()
	}

	return execCargo(cargoBin, []string{"--help"})
}

func dispatchCargoExecution(args []string) error {
	cleanArgs, hasInstall := extractInstallFlag(args)
	cargoBin := ResolveCargoBinary()
	if cargoBin == "" && hasInstall {
		_ = runCargoInstall([]string{"--yes"})
		cargoBin = ResolveCargoBinary()
	}

	if cargoBin == "" {
		PrintCargoInstallSuggestions()
		return NewMissingCargoError()
	}

	return execCargo(cargoBin, cleanArgs)
}

func runCargoInstall(rest []string) error {
	return cmdinstall.RunInstall(append([]string{constants.ToolCargo}, rest...))
}

func checkHelp(cmd string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(cmd, args)
	}
}
