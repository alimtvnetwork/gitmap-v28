package cmdcluster

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// HasHelpFlagFn is wired by cmd/di_hooks.go.
var HasHelpFlagFn func(args []string) bool

func hasHelpFlag(args []string) bool {
	if HasHelpFlagFn != nil {
		return HasHelpFlagFn(args)
	}
	return false
}

// CheckHelpOrEmptyFn is wired by cmd/di_hooks.go.
var CheckHelpOrEmptyFn func(command string, args []string)

func CheckHelpOrEmpty(command string, args []string) {
	if CheckHelpOrEmptyFn != nil {
		CheckHelpOrEmptyFn(command, args)
	}
}
