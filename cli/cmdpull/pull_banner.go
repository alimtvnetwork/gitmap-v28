package cmdpull

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintPullBanner prints runtime execution header with full command name, CWD, and GitMap version.
func PrintPullBanner(fullCmd, invokedAlias string, isShortForm bool) {
	cwd, _ := os.Getwd()
	arrow := resolveSubArrow()
	versionStr := formatVersionTag(constants.Version)

	fmt.Printf("\n  %s%s%s %sgitmap %s%s %s(%s)%s %s(cwd: %s)%s\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorBold, fullCmd, constants.ColorReset,
		constants.ColorCyan, versionStr, constants.ColorReset,
		constants.ColorDim, cwd, constants.ColorReset)

	if isShortForm && invokedAlias != "" {
		printAliasExpansionLine(invokedAlias, fullCmd, arrow)
	}
}

func printAliasExpansionLine(alias, fullCmd, arrow string) {
	fmt.Printf("    %s[alias: %s %s gitmap %s]%s\n",
		constants.ColorDim, alias, arrow, fullCmd, constants.ColorReset)
}

func formatVersionTag(raw string) string {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return "v0.0.0"
	}
	if !strings.HasPrefix(clean, "v") {
		return "v" + clean
	}

	return clean
}

// FormatPullBannerString returns formatted banner text for testing or logging.
func FormatPullBannerString(fullCmd, invokedAlias, cwd, version string, isShortForm bool) string {
	arrow := "→"
	versionStr := formatVersionTag(version)
	header := fmt.Sprintf("  %s gitmap %s (%s) (cwd: %s)", arrow, fullCmd, versionStr, cwd)
	if isShortForm && invokedAlias != "" {
		header += fmt.Sprintf("\n    [alias: %s %s gitmap %s]", invokedAlias, arrow, fullCmd)
	}

	return header
}
