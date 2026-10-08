package cmdadd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/add"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunAdd(args []string) error {
	checkHelp(constants.CmdAdd, args)
	if err := add.Run(args); err != nil {
		return apperror.WrapSimple(err, "Error:")
	}

	return nil
}
