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
		return "where.exe gitmap 2>nul || gitmap.exe version 2>nul || if exist \"%LOCALAPPDATA%\\gitmap-cli\\gitmap.exe\" (exit 0) else if exist \"%LOCALAPPDATA%\\gitmap\\bin\\gitmap.exe\" (exit 0) else (exit 1)", "cmd"
	}

	return wrapUnixPath("which gitmap >/dev/null 2>&1 || command -v gitmap >/dev/null 2>&1 || [ -x \"$HOME/.local/bin/gitmap\" ] || [ -x \"/usr/local/bin/gitmap\" ] || [ -x \"/usr/bin/gitmap\" ] || [ -x \"$HOME/go/bin/gitmap\" ] || [ -x \"/snap/bin/gitmap\" ] || [ -x \"$HOME/.local/bin/gitmap-cli/gitmap\" ]"), "sh"
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
