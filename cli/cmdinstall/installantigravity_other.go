//go:build !linux && !windows && !darwin

package cmdinstall

import "github.com/alimtvnetwork/gitmap-v28/cli/apperror"

func findInstalledAntigravityDesktopPath() (string, bool) {
	return "", false
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	return apperror.NewSimple("unsupported platform for Antigravity desktop", "E_UNSUPPORTED_PLATFORM")
}
