// Package theme resolves the active terminal color palette and installs
// an ANSI rewrite filter on stdout / stderr so every existing Print /
// Printf call adapts to the user's --theme choice without per-site
// changes.
//
// Three modes:
//
//   - bright     (default): passthrough — the bright + bold ANSI
//     palette baked into constants.Color* is what the user sees.
//   - standard:             downgrade bright codes to plain 3X codes
//     (the pre-v5.13 look) for users on light themes or older
//     terminals where bright-bold is too loud.
//   - monochrome / mono:    strip every SGR escape so output stays
//     readable when piped into tools that don't grok ANSI (diff,
//     less without -R, log scrapers).
//
// Resolution order (first wins):
//
//  1. `--theme=<mode>` / `--theme <mode>` on the command line
//     (stripped from os.Args by cmd.stripThemeFlag before subcommand
//     dispatch).
//  2. GITMAP_THEME env var (also exported by the flag stripper so
//     subprocesses inherit the choice).
//  3. constants.ThemeDefault.
package theme

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Mode is the resolved palette selection.
type Mode int

const (
	// ModeBright is the loud bright+bold palette (default).
	ModeBright Mode = iota
	// ModeStandard is the muted pre-v5.13 plain palette.
	ModeStandard
	// ModeMono strips all ANSI SGR codes.
	ModeMono
)

// String returns the canonical lowercase label for the mode.
func (m Mode) String() string {
	switch m {
	case ModeStandard:
		return constants.ThemeStandard
	case ModeMono:
		return constants.ThemeMonochrome
	default:
		return constants.ThemeBright
	}
}

// Parse maps a user-supplied label to a Mode. Unknown labels (and the
// empty string) fall back to ModeBright so a typo never crashes a CLI
// run — they just render in the default palette.
func Parse(label string) Mode {
	switch strings.ToLower(strings.TrimSpace(label)) {
	case constants.ThemeStandard:
		return ModeStandard
	case constants.ThemeMonochrome, constants.ThemeMono:
		return ModeMono
	default:
		return ModeBright
	}
}

// Resolve picks the active mode from env (set by the flag stripper or
// the user's shell). Called once at startup by Install.
func Resolve() Mode {
	return Parse(os.Getenv(constants.EnvTheme))
}

// IsValidLabel reports whether label names a known theme. Used by the
// flag stripper to reject typos with a clear error instead of silently
// falling back.
func IsValidLabel(label string) bool {
	switch strings.ToLower(strings.TrimSpace(label)) {
	case constants.ThemeBright,
		constants.ThemeStandard,
		constants.ThemeMonochrome,
		constants.ThemeMono:
		return true
	default:
		return false
	}
}
