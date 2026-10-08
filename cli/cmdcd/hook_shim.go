package cmdcd

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// WriteShellHandoffFn is wired by cmd/di_hooks.go.
var WriteShellHandoffFn func(targetPath string)

func WriteShellHandoff(targetPath string) {
	if WriteShellHandoffFn != nil {
		WriteShellHandoffFn(targetPath)
	}
}

// RunCDSpecialRepoFn is wired by cmd/di_hooks.go.
var RunCDSpecialRepoFn func(sub string, args []string) error

func runCDSpecialRepo(sub string, args []string) error {
	if RunCDSpecialRepoFn != nil {
		return RunCDSpecialRepoFn(sub, args)
	}
	return nil
}

// HasAliasFn is wired by cmd/di_hooks.go.
var HasAliasFn func() bool

func HasAlias() bool {
	if HasAliasFn != nil {
		return HasAliasFn()
	}
	return false
}

// GetAliasPathFn is wired by cmd/di_hooks.go.
var GetAliasPathFn func() string

func GetAliasPath() string {
	if GetAliasPathFn != nil {
		return GetAliasPathFn()
	}
	return ""
}

// DispatchFn is wired by cmd/di_hooks.go.
var DispatchFn func(command string)

func dispatch(command string) {
	if DispatchFn != nil {
		DispatchFn(command)
	}
}
