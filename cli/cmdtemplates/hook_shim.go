package cmdtemplates

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical checkHelp implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// IsHelpFlagFn is wired by cmd/di_hooks.go to the canonical implementation.
var IsHelpFlagFn func(token string) bool

// IsHelpFlag reports whether token is a help request indicator.
func IsHelpFlag(token string) bool {
	if IsHelpFlagFn != nil {
		return IsHelpFlagFn(token)
	}

	return false
}

// PrintHelpAndExitFn is wired by cmd/di_hooks.go to the canonical implementation.
var PrintHelpAndExitFn func(command string, args []string)

func printHelpAndExit(command string, args []string) {
	if PrintHelpAndExitFn != nil {
		PrintHelpAndExitFn(command, args)
	}
}
