package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	osReadFileHook  = os.ReadFile
	osWriteFileHook = os.WriteFile
)

func evaluateZshPrompt(ans string) {
	if ans == "y" || ans == "yes" {
		executeZshInstall()

		return
	}

	fmt.Println("  Skipped ZSH installation.")
}

func executeZshInstall() {
	installMissingZsh()
	executeOhMyZshInstall()
}

func installMissingZsh() {
	_, isZshInstalled := findZshBinary()
	if isZshInstalled {
		return
	}

	runAptInstallZsh()
}

func runAptInstallZsh() {
	installCommand := execCommandFunc("sudo", "apt", "install", "-y", "zsh")
	installCommand.Stdout, installCommand.Stderr = os.Stdout, os.Stderr
	if executeErr := installCommand.Run(); executeErr != nil {
		fmt.Printf("  Failed ZSH install: %v\n", executeErr)

		return
	}
}

func executeOhMyZshInstall() {
	if isOhMyZshInstalled() {
		return
	}

	runOhMyZshInstallScript()
}

func runOhMyZshInstallScript() {
	installScriptString := `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended`
	ohMyZshCommand := execCommandFunc("sh", "-c", installScriptString)
	ohMyZshCommand.Stdout, ohMyZshCommand.Stderr = os.Stdout, os.Stderr
	if executeErr := ohMyZshCommand.Run(); executeErr != nil {
		fmt.Printf("  Failed Oh-My-Zsh install: %v\n", executeErr)

		return
	}

	configureZshTheme()
}

func configureZshTheme() {
	homeDirectoryPath, getHomeErr := userHomeDirFunc()
	if getHomeErr != nil {
		return
	}

	zshrcFilePath := filepath.Join(homeDirectoryPath, ".zshrc")
	zshrcContentBytes, readErr := osReadFileHook(zshrcFilePath)
	if readErr != nil {
		return
	}

	applyThemeReplacement(zshrcFilePath, zshrcContentBytes)
}

func applyThemeReplacement(zshrcFilePath string, zshrcContentBytes []byte) {
	originalContentString := string(zshrcContentBytes)
	if strings.Contains(originalContentString, `ZSH_THEME="agnoster"`) {
		return
	}

	replacedContentString := strings.Replace(originalContentString, `ZSH_THEME="robbyrussell"`, `ZSH_THEME="agnoster"`, 1)
	writeErr := osWriteFileHook(zshrcFilePath, []byte(replacedContentString), 0644)
	if writeErr != nil {
		fmt.Printf("  Failed writing theme: %v\n", writeErr)

		return
	}

	fmt.Println("  ZSH theme configured successfully.")
}
