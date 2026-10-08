package cmdhaschange

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// ResolveReleaseAliasPathFn is wired by cmd/di_hooks.go to the canonical implementation.
var ResolveReleaseAliasPathFn func(alias string) string

func resolveReleaseAliasPath(alias string) string {
	if ResolveReleaseAliasPathFn != nil {
		return ResolveReleaseAliasPathFn(alias)
	}
	return alias
}
