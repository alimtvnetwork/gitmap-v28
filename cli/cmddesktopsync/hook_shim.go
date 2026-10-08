package cmddesktopsync

// HasGHDesktopInstallFlagFn is wired by cmd/di_hooks.go.
var HasGHDesktopInstallFlagFn func(args []string) bool

func hasGHDesktopInstallFlag(args []string) bool {
	if HasGHDesktopInstallFlagFn != nil {
		return HasGHDesktopInstallFlagFn(args)
	}
	return false
}
