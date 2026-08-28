package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/ag"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runAg(args []string) error {
	checkHelp(constants.CmdAg, args)
	if err := ag.Run(args); err != nil {
		return apperror.Wrap(err, "Error:", nil)
	}
	return nil
}
