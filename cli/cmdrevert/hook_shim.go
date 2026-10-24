package cmdrevert

// CheckHelpFn is wired by cmd/di_hooks_wave2.go.
var CheckHelpFn func(command string, args []string)

// RepoPathFn provides the active repository root path.
var RepoPathFn func() string

func checkHelp(command string, args []string) {
	if CheckHelpFn == nil {
		return
	}

	CheckHelpFn(command, args)
}
