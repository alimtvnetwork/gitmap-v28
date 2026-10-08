package cmdselfinstall

import (
	"os/exec"
)

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// CleanCorruptedInstallDirsSilentFn is wired by cmd/di_hooks.go.
var CleanCorruptedInstallDirsSilentFn func()

func cleanCorruptedInstallDirsSilent() {
	if CleanCorruptedInstallDirsSilentFn != nil {
		CleanCorruptedInstallDirsSilentFn()
	}
}

// SetHiddenProcessAttrFn is wired by cmd/di_hooks.go.
var SetHiddenProcessAttrFn func(cmd *exec.Cmd)

func setHiddenProcessAttr(cmd *exec.Cmd) {
	if SetHiddenProcessAttrFn != nil {
		SetHiddenProcessAttrFn(cmd)
	}
}
