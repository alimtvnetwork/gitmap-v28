package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/folder"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runFolder(args []string) error {
	checkHelp(constants.CmdFolder, args)
	if err := folder.Run(args); err != nil {
		return apperror.Wrap(err, "Error:", nil)
	}
	return nil
}
