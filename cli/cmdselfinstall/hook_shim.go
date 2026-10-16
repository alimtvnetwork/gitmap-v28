package cmdselfinstall

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
