package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/add"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runAdd(args []string) error {
	checkHelp(constants.CmdAdd, args)
	if err := add.Run(args); err != nil {
		return apperror.Wrap(err, "Error:", nil)
	}
	return nil
}
