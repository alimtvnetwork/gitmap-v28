package cmddesktopsync

import "github.com/alimtvnetwork/gitmap-v28/cli/constants"

// HasGHDesktopInstallFlagFn is wired by cmd/di_hooks.go.
var HasGHDesktopInstallFlagFn func(args []string) bool

func hasGHDesktopInstallFlag(args []string) bool {
	if HasGHDesktopInstallFlagFn != nil {
		return HasGHDesktopInstallFlagFn(args)
	}
	for _, arg := range args {
		if arg == constants.FlagGHDesktopInstall || arg == constants.FlagGHDesktopInstallShort {
			return true
		}
	}
	return false
}
