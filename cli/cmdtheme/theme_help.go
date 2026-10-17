// Package cmdtheme — theme_help.go: --help texts and help-example registry
// data for `theme list` / `themes` and `theme preview` / `themep`
// (program 275, spec 03 format).
package cmdtheme

// isHelpRequested reports whether the user asked for help.
func isHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}

	return false
}

// themeHelpText is the full `--help` text for `theme`.
const themeHelpText = `gitmap theme — terminal palettes: list, preview, activate

USAGE
  gitmap theme list              list palettes and mark the active one (alias: gitmap themes)
  gitmap theme preview <name>    render samples in a palette without switching (alias: gitmap themep <name>)

EXAMPLE
  gitmap theme list
  gitmap theme preview tokyo-night

NOTES
  - Read-only: nothing here changes the active theme.
  - Activate with: gitmap --theme <name> <command>, or set GITMAP_THEME.
  - Piped output never carries ANSI — colors render on a TTY only.
`

// themeListHelpText is the full `--help` text for `theme list`.
const themeListHelpText = `gitmap theme list — list design-system palettes and mark the active one

USAGE
  gitmap theme list     (alias: gitmap themes)

EXAMPLE
  gitmap theme list

NOTES
  - Slugs come from the 273 theme registry (navy-purple default, tokyo-night,
    vscode-dark-plus, monokai-pro); legacy bright/standard/monochrome/mono
    labels still work via --theme / GITMAP_THEME.
  - Read-only.
`

// themePreviewHelpText is the full `--help` text for `theme preview`.
const themePreviewHelpText = `gitmap theme preview — render sample widgets in a palette without switching

USAGE
  gitmap theme preview <name>     (alias: gitmap themep <name>)

EXAMPLE
  gitmap theme preview tokyo-night
  gitmap themep monochrome

NOTES
  - Shows a sample status table, checkmarks, and an error banner.
  - The active theme is unchanged afterwards.
  - Unknown names print the valid slug list.
  - Read-only.
`

// ThemeHelpExample is the canonical help-example shape, field-for-field
// identical to helpdoc.HelpExample (spec 03 §2). It drops verbatim into
// the cli/helpdoc registry (wired by the lead).
type ThemeHelpExample struct {
	WhatItDoes string
	Examples   []ThemeExample
}

// ThemeExample is one concrete invocation + its expected output (spec 03 §2).
type ThemeExample struct {
	Invocation     string
	ExpectedOutput string
}

const (
	themeListExampleWhat = "Lists the design-system palette slugs from the 273 theme registry plus legacy labels, marking the active theme."
	themeListExampleInv  = "gitmap theme list"
	themeListExampleOut  = "THEMES (design-system palettes):\n" +
		"  ○ monokai-pro\n" +
		"  ● navy-purple (active)\n" +
		"  ○ tokyo-night\n" +
		"  ○ vscode-dark-plus\n" +
		"LEGACY LABELS:\n" +
		"  bright → navy-purple  (legacy alias)\n" +
		"  standard → vscode-dark-plus  (legacy alias)\n" +
		"  monochrome → strip mode  (no color, plain output)\n" +
		"  mono → strip mode  (short alias for monochrome)\n" +
		"active: navy-purple (set via --theme <name> or GITMAP_THEME)"

	themePreviewExampleWhat = "Renders a sample status table, checkmarks, and an error banner in the named palette without changing the active theme."
	themePreviewExampleInv  = "gitmap theme preview tokyo-night"
	themePreviewExampleOut  = "PREVIEW: tokyo-night\n" +
		"REPO               BRANCH       STATUS\n" +
		"gitmap-v28         main         ✓ clean\n" +
		"go-email-reader    main         ✓ clean\n" +
		"story-telling-v1   main         ! 2 behind\n" +
		"✗ release failed: windows asset missing\n" +
		"(active theme unchanged)"
)

// ThemeListHelpExamples registers the `theme list` / `themes` help example,
// keyed by dispatcher verb AND alias (spec 03 §3).
var ThemeListHelpExamples = map[string]ThemeHelpExample{
	"theme list": {WhatItDoes: themeListExampleWhat, Examples: []ThemeExample{{Invocation: themeListExampleInv, ExpectedOutput: themeListExampleOut}}},
	"themes":     {WhatItDoes: themeListExampleWhat, Examples: []ThemeExample{{Invocation: "gitmap themes", ExpectedOutput: themeListExampleOut}}},
}

// ThemePreviewHelpExamples registers the `theme preview` / `themep` help example.
var ThemePreviewHelpExamples = map[string]ThemeHelpExample{
	"theme preview": {WhatItDoes: themePreviewExampleWhat, Examples: []ThemeExample{{Invocation: themePreviewExampleInv, ExpectedOutput: themePreviewExampleOut}}},
	"themep":        {WhatItDoes: themePreviewExampleWhat, Examples: []ThemeExample{{Invocation: "gitmap themep tokyo-night", ExpectedOutput: themePreviewExampleOut}}},
}
