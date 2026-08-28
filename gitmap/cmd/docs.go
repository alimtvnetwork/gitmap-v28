package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDocs opens the documentation website in the default browser.
func runDocs(args []string) error {
	checkHelp("docs", args)

	url := constants.DocsURL

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case constants.OSWindows:
		cmd = exec.Command(constants.CmdWindowsShell, constants.CmdArgSlashC, constants.CmdArgStart, url)
	case constants.OSDarwin:
		cmd = exec.Command(constants.CmdOpen, url)
	default:
		cmd = exec.Command(constants.CmdXdgOpen, url)
	}

	err := cmd.Start()
	if err != nil {
		return apperror.Wrap(err, constants.ErrDocsOpen, nil)
	}

	fmt.Printf(constants.MsgDocsOpened, url)
	return nil
}
