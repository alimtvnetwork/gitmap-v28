package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
)

// runUI dispatches gitmap ui and gitmap <module> ui.
func runUI(args []string) error {
	page := "settings"
	hasArg := len(args) > 0
	if hasArg {
		page = args[0]
	}

	return cmdui.RunUI(page, 8080)
}
