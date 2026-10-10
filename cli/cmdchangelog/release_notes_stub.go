package cmdchangelog

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
)

// RunReleaseNotes is a shim delegating to cmdrelease.RunReleaseNotes.
func RunReleaseNotes(args []string) error {
	return cmdrelease.RunReleaseNotes(args)
}
