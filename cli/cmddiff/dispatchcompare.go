package cmddiff

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgit"
)

// DispatchCompare routes `gitmap compare` / `gitmap cmp`.
func DispatchCompare(command string) (bool, error) {
	if isCompareCommand(command) {
		err := cmdgit.RunCompare(os.Args[constants.ArgsSubcommandOffset:])

		return true, err
	}

	return false, nil
}

func isCompareCommand(command string) bool {
	return command == constants.CmdCompare || command == constants.CmdCompareAlias
}
