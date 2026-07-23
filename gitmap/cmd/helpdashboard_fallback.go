package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// openHostedDocsFallback opens the hosted docs URL when the local docs site
// is unavailable (release didn't bundle docs-site.zip and download failed).
// Best-effort: prints the URL even if launching the browser fails so the user
// can copy it manually. See spec/02-app-issues/34-hd-hosted-docs-fallback.md.
func openHostedDocsFallback() {
	fmt.Fprintf(os.Stderr, constants.MsgHDHostedFallback, constants.DocsURL)
	openURL(constants.DocsURL)
}

// openURL launches the OS default browser for the given URL. Launch
// failures are logged to stderr (per zero-swallow policy) but never fatal —
// the caller always prints the URL first so the user can copy it manually.
func openURL(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case constants.OSWindows:
		cmd = exec.Command(constants.CmdWindowsShell, constants.CmdArgSlashC, constants.CmdArgStart, url)
	case constants.OSDarwin:
		cmd = exec.Command(constants.CmdOpen, url)
	default:
		cmd = exec.Command(constants.CmdXdgOpen, url)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not launch browser (%v); open the URL above manually.\n", err)
	}
}
