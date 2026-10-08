package cmdgitrm

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/gitrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunGitRm(args []string) error {
	checkHelp(constants.CmdGitRm, args)
	if err := gitrm.Run(args); err != nil {
		return apperror.WrapSimple(err, "Error:")
	}

	return nil
}
