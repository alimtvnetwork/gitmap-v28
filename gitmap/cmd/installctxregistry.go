package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runRegistryCommandsCtx executes a set of `reg` commands and returns
// the success count. Distinct from runRegistryCommands in
// installctxmenu.go to keep gitmap-ctx error messages independent of
// the older vscode-ctx / pwsh-ctx wording.
func runRegistryCommandsCtx(commands [][]string) int {
	success := 0

	for _, args := range commands {
		c := exec.Command(args[0], args[1:]...)
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, constants.MsgCtxRegFail, err)

			continue
		}

		success++
	}

	return success
}
