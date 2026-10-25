package cmdui

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogviewer"
)

// RunShowLog launches the interactive Show Log UI.
func RunShowLog(args []string) error {
	return cmdlogviewer.RunShowLog(args)
}
