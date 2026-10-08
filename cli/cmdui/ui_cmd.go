package cmdui

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

func init() {
	cmdssh.RunSSHUIFn = RunUI
}

// runUI dispatches gitmap ui and gitmap <module> ui.
func RunUICmd(args []string) error {
	checkHelp("ui", args)
	page := "settings"
	hasArg := len(args) > 0
	if hasArg {
		page = args[0]
	}

	return RunUI(page, 8080)
}
