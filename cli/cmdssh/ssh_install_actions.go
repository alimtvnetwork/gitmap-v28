package cmdssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

func installFreshGitmap(client *ssh.Client, header, osType string) {
	fmt.Printf("  %s %sgitmap missing, installing fresh...%s\n", header, constants.ColorCyan, constants.ColorReset)
	installCmd := BuildGitmapInstallOneLiner(osType, "latest")
	shell := "bash"
	if isWindowsOS(osType) {
		shell = "ps"
	}

	out, err := crypto.RunCommand(client, installCmd, shell)
	if err != nil {
		fmt.Printf("  %s %sInstall failed:%s %v\n", header, constants.ColorRed, constants.ColorReset, err)
		return
	}

	fmt.Printf("  %s %sInstalled successfully!%s\n%s\n", header, constants.ColorGreen, constants.ColorReset, strings.TrimSpace(out))
}

func updateExistingGitmap(client *ssh.Client, header, osType string) {
	fmt.Printf("  %s %sgitmap found, updating to latest...%s\n", header, constants.ColorYellow, constants.ColorReset)
	updateCmd := "gitmap update"
	shell := "bash"
	if isWindowsOS(osType) {
		shell = "ps"
	}

	out, err := crypto.RunCommand(client, updateCmd, shell)
	if err != nil {
		fmt.Printf("  %s %sUpdate failed, attempting reinstall:%s %v\n", header, constants.ColorYellow, constants.ColorReset, err)
		installFreshGitmap(client, header, osType)
		return
	}

	fmt.Printf("  %s %sUpdated successfully!%s\n%s\n", header, constants.ColorGreen, constants.ColorReset, strings.TrimSpace(out))
}
