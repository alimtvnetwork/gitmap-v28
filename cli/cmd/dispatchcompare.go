package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// dispatchCompare routes `gitmap compare` / `gitmap cmp`.
func dispatchCompare(command string) (bool, error) {
	if isCompareCommand(command) {
		err := cmdgit.RunCompare(os.Args[2:])

		return true, err
	}

	return false, nil
}

func isCompareCommand(command string) bool {
	return command == constants.CmdCompare || command == constants.CmdCompareAlias
}
