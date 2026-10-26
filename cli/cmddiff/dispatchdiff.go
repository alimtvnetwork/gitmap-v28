package cmddiff

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff/constants"
)

// DispatchDiff routes `gitmap diff` / `gitmap df`.
//
// Spec: companion to 02-spec/01-app/97-move-and-merge.md
func DispatchDiff(command string) (bool, error) {
	if command == constants.CmdDiff || command == constants.CmdDiffAlias {
		err := runDiff(os.Args[constants.ArgsSubcommandOffset:])

		return true, err
	}

	return false, nil
}
