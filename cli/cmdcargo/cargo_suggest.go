package cmdcargo

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	cargoconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintCargoInstallSuggestions prints formatted installation suggestions.
func PrintCargoInstallSuggestions() {
	fmt.Fprint(os.Stderr, FormatCargoInstallSuggestions(runtime.GOOS))
}

// FormatCargoInstallSuggestions formats cross-platform suggestions.
func FormatCargoInstallSuggestions(targetOS string) string {
	var sb strings.Builder
	writeMissingHeader(&sb)
	writeRecommendedSuggestions(&sb)
	writeFallbackSuggestions(&sb, targetOS)
	writeOfficialInstaller(&sb)
	writeAutoInstallTip(&sb)

	return sb.String()
}

func writeMissingHeader(sb *strings.Builder) {
	sb.WriteString("\n  " + constants.ColorRed + cargoconst.SuggestMissingHeader + constants.ColorReset + "\n\n")
}

func writeRecommendedSuggestions(sb *strings.Builder) {
	sb.WriteString("  Recommended (Install via GitMap):\n")
	sb.WriteString("    " + constants.ColorCyan + cargoconst.SuggestGitMapInstall + constants.ColorReset + "\n")
	sb.WriteString("    (or: " + cargoconst.SuggestGitMapInCargo + ")\n")
	sb.WriteString("    (or: " + cargoconst.SuggestGitMapInstallRust + ")\n\n")
}

func writeFallbackSuggestions(sb *strings.Builder, targetOS string) {
	sb.WriteString("  Package Manager Fallback:\n")
	sb.WriteString(formatFallbackCmd(targetOS))
}

func writeOfficialInstaller(sb *strings.Builder) {
	sb.WriteString("  Official Rustup Installer:\n")
	sb.WriteString(cargoconst.SuggestOfficialRustup)
}

func writeAutoInstallTip(sb *strings.Builder) {
	sb.WriteString("  Tip: Automatically install and run:\n")
	sb.WriteString("    gitmap cargo " + cargoconst.FlagLongInstall + " <command...>\n\n")
}

func formatFallbackCmd(targetOS string) string {
	switch targetOS {
	case constants.PlatformWindows:
		return cargoconst.WingetInstallRust
	case constants.PlatformDarwin:
		return cargoconst.BrewInstallRust
	default:
		return cargoconst.AptInstallCargo
	}
}

// NewMissingCargoError returns an AppError with code E7100.
func NewMissingCargoError() *apperror.AppError {
	return apperror.NewSimple(cargoconst.MissingCargoErrorMsg, cargoconst.ErrCodeMissingCargo)
}
