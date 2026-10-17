// Package cmdtheme — theme.go: `gitmap theme <list|preview>` dispatcher.
//
// Internal subcommand dispatch mirrors cmdinstall's pattern: the lead
// registers `theme` once and this file routes list/preview. `themep` is a
// top-level alias that calls RunThemePreview directly. Read-only.
package cmdtheme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// themeUsageLine is the one-line usage for subcommand errors.
const themeUsageLine = "usage: gitmap theme <list|preview>\n"

// RunTheme executes `gitmap theme list` / `gitmap theme preview <name>`.
func RunTheme(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(themeHelpText)
		cliexit.Exit(0)
	}

	if len(args) == 0 {
		return appfault.NewValidationError("missing subcommand\n" + themeUsageLine +
			"subcommands: list, preview")
	}

	switch sub := strings.ToLower(strings.TrimSpace(args[0])); sub {
	case "list":
		return runThemeList(args[1:])
	case "preview":
		return RunThemePreview(args[1:])
	default:
		return appfault.NewValidationError("unknown theme subcommand " + strconv.Quote(sub) + "\n" +
			themeUsageLine + "subcommands: list, preview")
	}
}
