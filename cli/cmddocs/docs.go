package cmddocs

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddocs/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddocs/enums"
)

var startCmdFn = func(cmd *exec.Cmd) error {
	return cmd.Start()
}

// ResolveViewerMode returns the default documentation viewing mode.
func ResolveViewerMode() enums.DocsViewerModeType {
	return enums.DocsViewerModeBrowser
}

// resolveBrowserCmd constructs an exec.Cmd for the specified OS and URL.
func resolveBrowserCmd(goos string, url string) *exec.Cmd {
	switch goos {
	case constants.OSWindows:
		return exec.Command(constants.CmdWindowsShell, constants.CmdArgSlashC, constants.CmdArgStart, url)
	case constants.OSDarwin:
		return exec.Command(constants.CmdOpen, url)
	default:
		return exec.Command(constants.CmdXdgOpen, url)
	}
}

// launchBrowser executes the platform-specific browser launcher command.
func launchBrowser(url string) error {
	cmd := resolveBrowserCmd(runtime.GOOS, url)
	err := startCmdFn(cmd)

	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDocsOpen)
	}

	return nil
}

// RunDocs opens the documentation website in the default browser.
func RunDocs(args []string) error {
	checkHelp(constants.CmdDocs, args)

	url := constants.DocsURL
	err := launchBrowser(url)

	if err != nil {
		return err
	}

	fmt.Printf(constants.MsgDocsOpened, url)

	return nil
}
