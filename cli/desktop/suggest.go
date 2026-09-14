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
	sb.WriteString("\nGitHub Desktop CLI not found.\n\n")
	sb.WriteString(formatOptionGitMap())
	sb.WriteString(formatOptionNative(targetOS))
	sb.WriteString(formatOptionOfficial())

	return sb.String()
}

func formatOptionGitMap() string {
	return "  Recommended (Install via GitMap):\n" +
		"    gitmap install github-desktop\n" +
		"    (or: gitmap in gd)\n" +
		"    (or auto-install: gitmap github-desktop --install)\n\n"
}

func formatOptionNative(targetOS string) string {
	cmd := resolveNativeCommand(targetOS)
	if cmd == "" {
		return ""
	}

	return "  Package Manager Fallback:\n" +
		"    " + cmd + "\n\n"
}

func formatOptionOfficial() string {
	return "  Official GUI Installer:\n" +
		"    " + URLGitHubDesktopOfficial + "\n\n"
}

// resolveNativeCommand resolves the single, standard package manager command for target OS.
func resolveNativeCommand(targetOS string) string {
	switch targetOS {
	case constants.PlatformWindows:
		return "winget install --id GitHub.GitHubDesktop"
	case constants.PlatformDarwin:
		return "brew install --cask github-desktop"
	case constants.PlatformLinux:
		return "sudo snap install github-desktop --beta"
	default:
		return ""
	}
}

// NewMissingCLIError returns a validation AppError for missing GitHub Desktop CLI.
func NewMissingCLIError() *apperror.AppError {
	return apperror.NewValidationError(constants.MsgDesktopNotFound)
}
