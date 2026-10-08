package cmdcreate

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}
