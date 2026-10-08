package cmdfoldertree

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/folder"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunFolder(args []string) error {
	checkHelp(constants.CmdFolder, args)
	if err := folder.Run(args); err != nil {
		return apperror.WrapSimple(err, "Error:")
	}

	return nil
}
