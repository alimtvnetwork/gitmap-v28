package cmdssh

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

func installFreshGitmap(client *ssh.Client, header, osType string) {
	fmt.Printf("  %s %sgitmap missing, installing fresh...%s\n", header, constants.ColorCyan, constants.ColorReset)
	cmd := BuildGitmapInstallOneLiner(osType, "latest")
	out, err := crypto.RunCommand(client, cmd, resolveRemoteShell(osType))
	reportRemoteExecution(header, "Installed", out, err)
}

func updateExistingGitmap(client *ssh.Client, header, osType string) {
	fmt.Printf("  %s %sgitmap found, updating to latest...%s\n", header, constants.ColorYellow, constants.ColorReset)
	out, err := crypto.RunCommand(client, "gitmap update", resolveRemoteShell(osType))
	if err != nil {
		fmt.Printf("  %s %sUpdate failed, attempting reinstall:%s %v\n", header, constants.ColorYellow, constants.ColorReset, err)
		installFreshGitmap(client, header, osType)
		return
	}
	fmt.Printf("  %s %sUpdated successfully!%s\n%s\n", header, constants.ColorGreen, constants.ColorReset, strings.TrimSpace(out))
}

func installRemotePackage(client *ssh.Client, header, osType, pkg string) {
	fmt.Printf("  %s %sInstalling package '%s' via gitmap...%s\n", header, constants.ColorCyan, pkg, constants.ColorReset)
	out, err := crypto.RunCommand(client, "gitmap install "+pkg, resolveRemoteShell(osType))
	reportRemoteExecution(header, "Installed "+pkg, out, err)
}

func isNodeAvailable(ip, header string) bool {
	isOnline, reason := CheckConnLiveness(context.Background(), ip, 22, 0)
	if !isOnline {
		fmt.Printf("  %s %sOFFLINE (skipped: %s)%s\n", header, constants.ColorYellow, reason, constants.ColorReset)
		return false
	}
	return true
}

func resolveRemoteShell(osType string) string {
	if isWindowsOS(osType) {
		return "ps"
	}
	return "bash"
}

func reportRemoteExecution(header, action, out string, err error) {
	if err != nil {
		fmt.Printf("  %s %s%s failed:%s %v\n", header, constants.ColorRed, action, constants.ColorReset, err)
		return
	}
	fmt.Printf("  %s %s%s successfully!%s\n%s\n", header, constants.ColorGreen, action, constants.ColorReset, strings.TrimSpace(out))
}
