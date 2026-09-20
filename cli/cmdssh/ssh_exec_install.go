package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"golang.org/x/crypto/ssh"
)

func getGitmapCheckCmd(osType string) (string, string) {
	if isWindowsOS(osType) {
		return "gitmap --version", ""
	}

	return wrapUnixPath("gitmap --version"), "bash"
}

func getGitmapInstallCmd(osType string) (string, string) {
	if isWindowsOS(osType) {
		return fmt.Sprintf("irm %s | iex", constants.SelfInstallRemotePwsh), "ps"
	}

	return fmt.Sprintf("curl -fsSL %s | bash", constants.SelfInstallRemoteBash), "bash"
}

func ensureGitmapInstalled(client *ssh.Client, osType, header string) error {
	checkCmd, shell := getGitmapCheckCmd(osType)
	_, err := crypto.RunCommand(client, checkCmd, shell)
	if err == nil {
		return nil
	}

	fmt.Printf("%s gitmap not found, installing...\n", header)
	installCmd, installShell := getGitmapInstallCmd(osType)
	_, err = crypto.RunCommand(client, installCmd, installShell)
	if err != nil {
		return fmt.Errorf("auto-install failed: %w", err)
	}

	return nil
}

func wrapUnixPath(cmd string) string {
	hasPathPrefix := strings.HasPrefix(cmd, "export PATH=")
	if hasPathPrefix {
		return cmd
	}

	return "export PATH=\"$HOME/.local/bin:$HOME/.local/bin/gitmap-cli:$HOME/bin:/usr/local/bin:/usr/local/go/bin:/snap/bin:$PATH\"; " + cmd
}
