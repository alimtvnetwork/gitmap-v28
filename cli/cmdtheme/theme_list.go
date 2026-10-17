// Package cmdtheme — theme_list.go: `gitmap theme list` (alias `themes`).
//
// Lists the design-system palette slugs straight from the 273 theme
// registry (termout.ThemeSlugs — never a hardcoded copy) plus the legacy
// labels, and marks the active theme. Read-only.
package cmdtheme

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// themeLegacyAlias is one legacy label and what it resolves to.
type themeLegacyAlias struct {
	label  string
	target string
	note   string
}

// themeLegacyAliases documents the pre-273 labels.
var themeLegacyAliases = []themeLegacyAlias{
	{label: constants.ThemeBright, target: termout.ThemeNavyPurple, note: "legacy alias"},
	{label: constants.ThemeStandard, target: termout.ThemeVSCodeDarkPlus, note: "legacy alias"},
	{label: constants.ThemeMonochrome, target: "strip mode", note: "no color, plain output"},
	{label: constants.ThemeMono, target: "strip mode", note: "short alias for monochrome"},
}

// runThemeList executes `gitmap theme list`.
func runThemeList(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(themeListHelpText)
		cliexit.Exit(0)
	}

	activeSlug, isStrip := activeThemeSlug()
	fmt.Println("THEMES (design-system palettes):")
	for _, slug := range termout.ThemeSlugs() {
		mark := "○"
		state := ""
		if !isStrip && slug == activeSlug {
			mark = "●"
			state = " (active)"
		}
		fmt.Printf("  %s %s%s\n", mark, slug, state)
	}

	fmt.Println("LEGACY LABELS:")
	for _, alias := range themeLegacyAliases {
		fmt.Printf("  %s → %s  (%s)\n", alias.label, alias.target, alias.note)
	}

	fmt.Printf("active: %s (set via --theme <name> or %s)\n", activeThemeLabel(activeSlug, isStrip), constants.EnvTheme)

	return nil
}

// activeThemeSlug resolves the effective design-system slug for this
// process. isStrip is true for monochrome/mono mode (no palette).
func activeThemeSlug() (string, bool) {
	palette, isStrip := termout.ResolveTheme()
	if isStrip || palette == nil {
		return "", true
	}

	return palette.Slug, false
}

// activeThemeLabel renders the effective theme for the summary line.
func activeThemeLabel(activeSlug string, isStrip bool) string {
	if isStrip {
		label := strings.ToLower(strings.TrimSpace(os.Getenv(constants.EnvTheme)))
		if label == constants.ThemeMono {
			return constants.ThemeMono + " (strip mode)"
		}

		return constants.ThemeMonochrome + " (strip mode)"
	}

	return activeSlug
}
