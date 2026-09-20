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
		return "where.exe gitmap 2>nul || gitmap version 2>nul || gitmap --version 2>nul", ""
	}

	return wrapUnixPath("for bin in gitmap \"$HOME/.local/bin/gitmap\" \"$HOME/.local/bin/gitmap-cli/gitmap\" \"$HOME/bin/gitmap\" /usr/local/bin/gitmap /usr/bin/gitmap /snap/bin/gitmap; do if command -v \"$bin\" >/dev/null 2>&1 || [ -x \"$bin\" ]; then if \"$bin\" version >/dev/null 2>&1 || \"$bin\" --version >/dev/null 2>&1; then exit 0; fi; fi; done; exit 1"), "bash"
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

	return performGitmapInstall(client, osType, header, checkCmd, shell)
}

func performGitmapInstall(client *ssh.Client, osType, header, checkCmd, shell string) error {
	fmt.Printf("%s gitmap not found, installing...\n", header)
	installCmd, installShell := getGitmapInstallCmd(osType)
	if _, err := crypto.RunCommand(client, installCmd, installShell); err != nil {
		return fmt.Errorf("auto-install failed: %w", err)
	}

	if _, err := crypto.RunCommand(client, checkCmd, shell); err != nil {
		return fmt.Errorf("gitmap install verification failed: %w", err)
	}
	fmt.Printf("%s ✓ gitmap installed successfully\n", header)

	return nil
}

func wrapUnixPath(cmd string) string {
	hasPathPrefix := strings.HasPrefix(cmd, "export PATH=")
	if hasPathPrefix {
		return cmd
	}

	return "export PATH=\"$HOME/.local/bin:$HOME/.local/bin/gitmap-cli:$HOME/bin:/usr/local/bin:/usr/local/go/bin:/snap/bin:$PATH\"; " + cmd
}
