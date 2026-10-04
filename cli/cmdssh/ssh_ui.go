package cmdssh

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSSHUIFn is the pluggable callback provided by cmdui package.
var RunSSHUIFn func(page string, preferredPort int) error

func runSSHUI(args []string) error {
	if RunSSHUIFn != nil {
		return RunSSHUIFn("ssh", 8080)
	}

	targetURL := "http://localhost:8080/ssh"
	fmt.Printf("%s⚡ Launching GitMap SSH Fleet Web UI at %s%s\n", constants.ColorGreen, targetURL, constants.ColorReset)

	return openBrowserFallbackURL(targetURL)
}

func openBrowserFallbackURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
