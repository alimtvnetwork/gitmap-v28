package cmdcargo

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintCargoInstallSuggestions prints formatted installation suggestions.
func PrintCargoInstallSuggestions() {
	fmt.Fprint(os.Stderr, FormatCargoInstallSuggestions(runtime.GOOS))
}

// FormatCargoInstallSuggestions formats cross-platform suggestions.
func FormatCargoInstallSuggestions(targetOS string) string {
	var sb strings.Builder
	sb.WriteString("\n  " + constants.ColorRed + "✗ Cargo (Rust toolchain) is not installed or not found in PATH." + constants.ColorReset + "\n\n")
	sb.WriteString("  Recommended (Install via GitMap):\n")
	sb.WriteString("    " + constants.ColorCyan + "gitmap install cargo" + constants.ColorReset + "\n")
	sb.WriteString("    (or: gitmap in cargo)\n")
	sb.WriteString("    (or: gitmap install rust)\n\n")
	sb.WriteString("  Package Manager Fallback:\n")
	sb.WriteString(formatFallbackCmd(targetOS))
	sb.WriteString("  Official Rustup Installer:\n")
	sb.WriteString("    https://rustup.rs/\n")
	sb.WriteString("    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh\n\n")
	sb.WriteString("  Tip: Automatically install and run:\n")
	sb.WriteString("    gitmap cargo --install <command...>\n\n")

	return sb.String()
}

func formatFallbackCmd(targetOS string) string {
	switch targetOS {
	case constants.PlatformWindows:
		return "    winget install Rustlang.Rustup\n    (or: choco install rust)\n\n"
	case constants.PlatformDarwin:
		return "    brew install rust\n\n"
	default:
		return "    sudo apt install -y cargo\n\n"
	}
}

// NewMissingCargoError returns an AppError with code E7100.
func NewMissingCargoError() *apperror.AppError {
	msg := "cargo is not installed or not in PATH. Run 'gitmap install cargo' to install it."

	return apperror.NewSimple(msg, "E7100")
}
