package cmdbranch

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// PrintHelpAndExitFn is wired by cmd/di_hooks.go to the canonical implementation.
var PrintHelpAndExitFn func(command string, args []string)

func printHelpAndExit(command string, args []string) {
	if PrintHelpAndExitFn != nil {
		PrintHelpAndExitFn(command, args)
	}
}
