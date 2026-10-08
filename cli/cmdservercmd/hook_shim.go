package cmdservercmd

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// CheckHelpOrEmptyFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpOrEmptyFn func(command string, args []string)

func CheckHelpOrEmpty(command string, args []string) {
	if CheckHelpOrEmptyFn != nil {
		CheckHelpOrEmptyFn(command, args)
	}
}
