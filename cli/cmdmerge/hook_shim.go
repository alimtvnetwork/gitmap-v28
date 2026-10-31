package cmdmerge

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

// CommitInRunner is wired by cmd/di_hooks_wave2.go to cmdcommitin.RunCommitIn.
var CommitInRunner func(args []string) error

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}
