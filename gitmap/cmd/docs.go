package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runDocs opens the documentation website in the default browser.
func runDocs(args []string) {
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
		fmt.Fprintf(os.Stderr, constants.ErrDocsOpen, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgDocsOpened, url)
}
