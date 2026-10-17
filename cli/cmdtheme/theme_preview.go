// Package cmdtheme — theme_preview.go: `gitmap theme preview <name>` (alias `themep`).
//
// Renders a sample status table, checkmarks, and an error banner in the
// named palette WITHOUT changing the active theme: the GITMAP_THEME
// override is process-local and restored before return. termout.RoleSGR
// stays TTY-gated, so piped output is plain. Read-only.
package cmdtheme

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// themePreviewUsageLine is the one-line usage for argument errors.
const themePreviewUsageLine = "usage: gitmap theme preview <name>\n"

// RunThemePreview executes `gitmap theme preview <name>`.
func RunThemePreview(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(themePreviewHelpText)
		cliexit.Exit(0)
	}

	if len(args) == 0 {
		return appfault.NewValidationError("missing theme name\n" + themePreviewUsageLine +
			"valid: " + strings.Join(validThemeLabels(), ", "))
	}

	label := strings.ToLower(strings.TrimSpace(args[0]))
	if !isValidThemeLabel(label) {
		return appfault.NewValidationError("unknown theme " + strconv.Quote(label) +
			" — valid: " + strings.Join(validThemeLabels(), ", "))
	}

	renderThemePreview(label)

	return nil
}

// isValidThemeLabel reports whether the label names a registry palette or a
// legacy label (bright/standard/monochrome/mono).
func isValidThemeLabel(label string) bool {
	for _, slug := range termout.ThemeSlugs() {
		if label == slug {
			return true
		}
	}

	switch label {
	case constants.ThemeBright, constants.ThemeStandard, constants.ThemeMonochrome, constants.ThemeMono:
		return true
	}

	return false
}

// validThemeLabels returns registry slugs plus legacy labels, sorted.
func validThemeLabels() []string {
	labels := append([]string{}, termout.ThemeSlugs()...)
	labels = append(labels,
		constants.ThemeBright, constants.ThemeStandard, constants.ThemeMonochrome, constants.ThemeMono)
	sort.Strings(labels)

	return labels
}

// renderThemePreview renders the sample widgets with the named palette
// active for this process only; the previous GITMAP_THEME value is
// restored afterwards so the active theme never changes.
func renderThemePreview(label string) {
	previous := os.Getenv(constants.EnvTheme)
	os.Setenv(constants.EnvTheme, label)
	defer os.Setenv(constants.EnvTheme, previous)

	ok := termout.RoleSGR(termout.RoleOK)
	warn := termout.RoleSGR(termout.RoleWarn)
	errColor := termout.RoleSGR(termout.RoleError)
	primary := termout.RoleSGRBold(termout.RolePrimary)
	muted := termout.RoleSGR(termout.RoleFGMuted)
	reset := termout.Reset()

	fmt.Printf("%sPREVIEW: %s%s\n", primary, label, reset)
	fmt.Printf("%s%-18s %-12s %s%s\n", muted, "REPO", "BRANCH", "STATUS", reset)
	fmt.Printf("%s%-18s %-12s %s✓ clean%s\n", muted, "gitmap-v28", "main", ok, reset)
	fmt.Printf("%s%-18s %-12s %s✓ clean%s\n", muted, "go-email-reader", "main", ok, reset)
	fmt.Printf("%s%-18s %-12s %s! 2 behind%s\n", muted, "story-telling-v1", "main", warn, reset)
	fmt.Printf("%s✗ release failed: windows asset missing%s\n", errColor, reset)
	fmt.Printf("%s(active theme unchanged)%s\n", muted, reset)
}
