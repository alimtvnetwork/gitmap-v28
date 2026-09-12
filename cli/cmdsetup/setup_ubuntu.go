package cmdsetup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ensureZshUbuntuStep coordinates Ubuntu ZSH setup.
func ensureZshUbuntuStep(isDryRun, isSkipZsh bool) {
	if isZshBypassed(isSkipZsh) || checkExistingZsh() {
		return
	}

	if isDryRun {
		fmt.Printf("  %s[dry-run]%s would install ZSH\n", constants.ColorDim, constants.ColorReset)

		return
	}

	promptZshIfInteractive()
}

func isZshBypassed(isSkipZsh bool) bool {
	if !isUbuntuOS() || isSkipZsh {
		return true
	}

	return isSkipZshEnv()
}

func isUbuntuOS() bool {
	b, readErr := os.ReadFile("/etc/os-release")
	if readErr != nil {
		return false
	}

	return strings.Contains(string(b), "Ubuntu")
}

func checkExistingZsh() bool {
	zshPath, isInstalled := findZshBinary()
	if !isInstalled {
		return false
	}

	printZshInstalled(zshPath)
	_ = checkExistingOhMyZsh()

	return true
}

func printZshInstalled(zshPath string) {
	ver := getZshVersion(zshPath)
	msg := "  " + constants.ColorGreen + "✓" + constants.ColorReset + " ZSH is already installed"
	if len(ver) > 0 {
		msg += " (" + ver + ")"
	}

	fmt.Println(msg)
}

func checkExistingOhMyZsh() bool {
	if isOhMyZshInstalled() {
		fmt.Printf("  %s✓%s Oh-My-Zsh is already installed\n", constants.ColorGreen, constants.ColorReset)

		return true
	}

	return false
}

func promptZshIfInteractive() {
	if isStdinTerminal() {
		promptZshInstallation()

		return
	}

	fmt.Println("  Non-interactive terminal: skipping ZSH installation prompt.")
}

func promptZshInstallation() {
	fmt.Print("  Install ZSH and Oh-My-Zsh? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	resp, readErr := reader.ReadString('\n')
	if readErr != nil {
		return
	}

	evaluateZshPrompt(strings.TrimSpace(strings.ToLower(resp)))
}
