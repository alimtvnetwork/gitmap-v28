package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/gitrm"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runGitRm(args []string) error {
	checkHelp(constants.CmdGitRm, args)
	if err := gitrm.Run(args); err != nil {
		return apperror.WrapSimple(err, "Error:")
	}
	return nil
}
