package cmdcode

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)
