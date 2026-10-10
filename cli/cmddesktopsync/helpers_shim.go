package cmddesktopsync

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
)

// ensureGHDesktopInstalled auto-installs GitHub Desktop if --install/-i is present and CLI is missing.
func ensureGHDesktopInstalled(args []string) error {
	hasInstall := hasGHDesktopInstallFlag(args)
	if !hasInstall {
		return nil
	}

	cli := desktop.ResolveCLI()
	if cli != "" {
		return nil
	}

	return cmdinstall.RunInstall([]string{constants.ToolGitHubDesktop, "--yes"})
}
