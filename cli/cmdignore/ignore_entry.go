package cmdignore

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/ignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runIgnore(args []string) error {
	checkHelp(constants.CmdIgnore, args)
	if err := ignore.RunIgnore(args); err != nil {
		return apperror.WrapSimple(err, "Error:")
	}

	return nil
}

func RunIgnoreRm(args []string) error {
	checkHelp(constants.CmdIgnoreRm, args)
	if err := ignore.RunIgnoreRm(args); err != nil {
		return apperror.WrapSimple(err, "Error:")
	}

	return nil
}
