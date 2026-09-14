package desktop

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// URLGitHubDesktopOfficial points to the official GitHub Desktop installer.
const URLGitHubDesktopOfficial = "https://desktop.github.com"

// PrintInstallSuggestions prints formatted remediation guidance to os.Stderr.
func PrintInstallSuggestions() {
	fmt.Fprint(os.Stderr, FormatInstallSuggestions(runtime.GOOS))
}

// FormatInstallSuggestions returns clean multi-platform installation options.
func FormatInstallSuggestions(targetOS string) string {
	var sb strings.Builder
	sb.WriteString("\nTo install GitHub Desktop, choose one of the options below:\n\n")
	sb.WriteString(formatOptionOne())
	sb.WriteString(formatOptionTwo(targetOS))
	sb.WriteString(formatOptionThree())

	return sb.String()
}

func formatOptionOne() string {
	return "  Option 1 (Recommended - GitMap):\n" +
		"    gitmap install github-desktop\n" +
		"    (or: gitmap in gd)\n\n"
}

func formatOptionTwo(targetOS string) string {
	cmds := resolveNativeCommands(targetOS)
	if len(cmds) == 0 {
		return ""
	}

	return "  Option 2 (Native Package Managers):\n    " +
		strings.Join(cmds, "\n    ") + "\n\n"
}

func formatOptionThree() string {
	return "  Option 3 (Official GUI Installer):\n" +
		"    " + URLGitHubDesktopOfficial + "\n\n"
}

// resolveNativeCommands resolves native package manager commands for target OS.
func resolveNativeCommands(targetOS string) []string {
	switch targetOS {
	case constants.PlatformWindows:
		return windowsNativeCommands()
	case constants.PlatformDarwin:
		return darwinNativeCommands()
	case constants.PlatformLinux:
		return linuxNativeCommands()
	default:
		return nil
	}
}

func windowsNativeCommands() []string {
	return []string{
		"winget install --id GitHub.GitHubDesktop",
		"choco install github-desktop",
	}
}

func darwinNativeCommands() []string {
	return []string{
		"brew install --cask github-desktop",
	}
}

func linuxNativeCommands() []string {
	return []string{
		"wget -qO - https://mirror.mwt.me/ghd/gpgkey | sudo tee /etc/apt/keyrings/shiftkey-packages.asc > /dev/null && echo \"deb [arch=amd64 signed-by=/etc/apt/keyrings/shiftkey-packages.asc] https://mirror.mwt.me/ghd/deb/ any main\" | sudo tee /etc/apt/sources.list.d/shiftkey-packages.list && sudo apt update && sudo apt install -y github-desktop",
		"sudo snap install github-desktop --beta",
	}
}

// NewMissingCLIError returns a validation AppError for missing GitHub Desktop CLI.
func NewMissingCLIError() *apperror.AppError {
	return apperror.NewValidationError(constants.MsgDesktopNotFound)
}
