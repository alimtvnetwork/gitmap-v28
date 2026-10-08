package cmdbookmark

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// DispatchFn is wired by cmd/di_hooks.go to the canonical command dispatcher.
var DispatchFn func(command string)

func dispatch(command string) {
	if DispatchFn != nil {
		DispatchFn(command)
	}
}
