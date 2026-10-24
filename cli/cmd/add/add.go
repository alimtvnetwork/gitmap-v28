package add

import (
	"os"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdadd/constants"
)

// Run executes the add command logic.
func Run(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple(constants.OpAddRun, constants.ErrCodeAddMissingArg)
	}

	target := args[0]
	switch target {
	case constants.SubCmdCommonAttr:
		return writeCommonAttr()
	case constants.SubCmdCommonIgnore:
		return writeCommonIgnore()
	default:
		return apperror.New(constants.OpAddRun, constants.ErrCodeAddUnknownTarget, map[string]any{constants.KeyTarget: target})
	}
}

func writeCommonAttr() error {
	if err := os.WriteFile(constants.FileGitAttributes, []byte(constants.TemplateCommonAttr), constants.DefaultFileMode); err != nil {
		return apperror.WrapSimple(err, constants.ErrMsgWriteAttrFailed)
	}

	pterm.Success.Println(constants.MsgCreatedGitAttributes)

	return nil
}

func writeCommonIgnore() error {
	if err := os.WriteFile(constants.FileGitIgnore, []byte(constants.TemplateCommonIgnore), constants.DefaultFileMode); err != nil {
		return apperror.WrapSimple(err, constants.ErrMsgWriteIgnoreFailed)
	}

	pterm.Success.Println(constants.MsgCreatedGitIgnore)

	return nil
}
