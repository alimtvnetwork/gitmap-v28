package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// dispatchDiff routes `gitmap diff` / `gitmap df`.
//
// Spec: companion to spec/01-app/97-move-and-merge.md
func dispatchDiff(command string) bool {
	if command == constants.CmdDiff || command == constants.CmdDiffAlias {
		runDiff(os.Args[2:])

		return true
	}

	return false
}
