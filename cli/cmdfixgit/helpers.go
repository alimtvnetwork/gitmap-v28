package cmdfixgit

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func checkHelp(command string, args []string) {
	if !hasHelpFlag(args) {
		return
	}
	RenderFixGitHelp()
	cliexit.Exit(0)
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}

	return false
}
