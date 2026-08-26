package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ensureZshUbuntuStep coordinates Ubuntu ZSH setup.
func ensureZshUbuntuStep(isDryRun bool) {
	if !isUbuntuOS() {
		return
	}
	fmt.Printf("\n  %s■ ZSH Setup (Ubuntu) —%s\n", constants.ColorYellow, constants.ColorReset)
	if isDryRun {
		fmt.Printf("  %s[dry-run]%s would install ZSH\n", constants.ColorDim, constants.ColorReset)
		return
	}
	promptZshInstallation()
}

func isUbuntuOS() bool {
	contentBytes, readErr := os.ReadFile("/etc/os-release")
	if readErr != nil {
		return false
	}
	fileContentString := string(contentBytes)
	return strings.Contains(fileContentString, "Ubuntu")
}

func promptZshInstallation() {
	fmt.Print("  Install ZSH and Oh-My-Zsh? (y/N): ")
	inputReader := bufio.NewReader(os.Stdin)
	responseString, readErr := inputReader.ReadString('\n')
	if readErr != nil {
		return
	}
	cleanedResponseString := strings.TrimSpace(strings.ToLower(responseString))
	evaluateZshPrompt(cleanedResponseString)
}

func evaluateZshPrompt(userResponseString string) {
	isYesAnswer := userResponseString == "y" || userResponseString == "yes"
	if isYesAnswer {
		executeZshInstall()
	} else {
		fmt.Println("  Skipped ZSH installation.")
	}
}

func executeZshInstall() {
	installCommand := exec.Command("sudo", "apt", "install", "-y", "zsh")
	installCommand.Stdout = os.Stdout
	installCommand.Stderr = os.Stderr
	if executeErr := installCommand.Run(); executeErr != nil {
		fmt.Printf("  Failed ZSH install: %v\n", executeErr)
		return
	}
	executeOhMyZshInstall()
}

func executeOhMyZshInstall() {
	installScriptString := `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended`
	ohMyZshCommand := exec.Command("sh", "-c", installScriptString)
	ohMyZshCommand.Stdout = os.Stdout
	ohMyZshCommand.Stderr = os.Stderr
	if executeErr := ohMyZshCommand.Run(); executeErr != nil {
		fmt.Printf("  Failed Oh-My-Zsh install: %v\n", executeErr)
		return
	}
	configureZshTheme()
}

func configureZshTheme() {
	homeDirectoryPath, getHomeErr := os.UserHomeDir()
	if getHomeErr != nil {
		return
	}
	zshrcFilePath := filepath.Join(homeDirectoryPath, ".zshrc")
	zshrcContentBytes, readErr := os.ReadFile(zshrcFilePath)
	if readErr != nil {
		return
	}
	applyThemeReplacement(zshrcFilePath, zshrcContentBytes)
}

func applyThemeReplacement(zshrcFilePath string, zshrcContentBytes []byte) {
	originalContentString := string(zshrcContentBytes)
	replacedContentString := strings.Replace(originalContentString, `ZSH_THEME="robbyrussell"`, `ZSH_THEME="agnoster"`, 1)
	writeErr := os.WriteFile(zshrcFilePath, []byte(replacedContentString), 0644)
	if writeErr != nil {
		fmt.Printf("  Failed writing theme: %v\n", writeErr)
		return
	}
	fmt.Println("  ZSH theme configured successfully.")
}
