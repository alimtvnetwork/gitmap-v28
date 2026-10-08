package cmdprofiles

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// IsHelpFlagFn is wired by cmd/di_hooks.go.
var IsHelpFlagFn func(arg string) bool

func isHelpFlag(arg string) bool {
	if IsHelpFlagFn != nil {
		return IsHelpFlagFn(arg)
	}
	return false
}
