package cmdjoin
// CheckHelpOrEmptyFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpOrEmptyFn func(command string, args []string)

func CheckHelpOrEmpty(command string, args []string) {
	if CheckHelpOrEmptyFn != nil {
		CheckHelpOrEmptyFn(command, args)
	}
}
