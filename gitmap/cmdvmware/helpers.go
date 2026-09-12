package cmdvmware

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/helptext"
)

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)

			return
		}
	}
}

func hasDryRunFlag(args []string) bool {
	for _, a := range args {
		if a == "--dry-run" {
			return true
		}
	}

	return false
}
