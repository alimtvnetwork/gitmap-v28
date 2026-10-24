// Package termout resolves the active terminal color palette and
// provides the ANSI SGR rewrite filter. Filtering runs synchronously
// inside the cli/output FilterWriter that cmd.Run builds from the
// resolved mode — os.Stdout/os.Stderr are never reassigned; UI code
// takes the explicit dispatch-context writer instead.
//
// Named themes now live in the design-system registry (themes.go):
// navy-purple (default), tokyo-night, vscode-dark-plus, monokai-pro.
// The filter modes below stay byte-identical — new themes resolve to
// the passthrough mode and render role colors via Role()/RoleBold(),
// which emit no SGR at all when stdout is not a TTY (dumb-terminal
// safety) and are stripped by the mono mode on a TTY.
//
// Three filter modes:
//
//   - bright     (default): passthrough — the bright + bold ANSI
//     palette baked into constants.Color* is what the user sees.
//     The navy-purple, tokyo-night, and monokai-pro themes also
//     resolve here; their role colors are emitted per-theme by
//     Role()/RoleBold().
//   - standard:             downgrade bright codes to plain 3X codes
//     (the pre-v5.13 look) for users on light themes or older
//     terminals where bright-bold is too loud. The vscode-dark-plus
//     theme alias resolves here; its role SGRs are not in the
//     downgrade table so they pass through untouched.
//   - monochrome / mono:    strip every SGR escape so output stays
//     readable when piped into tools that don't grok ANSI (diff,
//     less without -R, log scrapers).
//
// Resolution order (first wins):
//
//  1. `--theme=<slug>` / `--theme <slug>` on the command line
//     (stripped from os.Args by cmd.stripThemeFlag before subcommand
//     dispatch).
//  2. GITMAP_THEME env var (also exported by the flag stripper so
//     subprocesses inherit the choice).
//  3. constants.ThemeDefault.
package termout

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/termout/enums"

// ModeType is the resolved palette selection.
// ModeType aliases the strongly-typed enum from the enums package.
type ModeType = enums.ModeType

const (
	ModeBright   = enums.ModeBright
	ModeStandard = enums.ModeStandard
	ModeMono     = enums.ModeMono
)

// themeModeByLabel maps a lowercased theme label to its render mode;
// labels absent from the map fall back to ModeBright in Parse.
var themeModeByLabel = map[string]ModeType{
	constants.ThemeStandard:   ModeStandard,
	ThemeVSCodeDarkPlus:       ModeStandard,
	constants.ThemeMonochrome: ModeMono,
	constants.ThemeMono:       ModeMono,
}

// Parse maps a user-supplied label to a ModeType. Unknown labels (and the
// empty string) fall back to ModeBright so a typo never crashes a CLI
// run — they just render in the default palette. The named design-system
// themes resolve to passthrough/downgrade modes here; their actual
// palettes come from the registry in themes.go.
func Parse(label string) ModeType {
	if mode, ok := themeModeByLabel[strings.ToLower(strings.TrimSpace(label))]; ok {
		return mode
	}

	return ModeBright
}

// Resolve picks the active mode from env (set by the flag stripper or
// the user's shell). Called once at startup by Install.
func Resolve() ModeType {
	return Parse(os.Getenv(constants.EnvTheme))
}

// IsValidLabel reports whether label names a known theme — a registry
// slug or a legacy alias. Used by the flag stripper to reject typos
// with a clear error instead of silently falling back.
func IsValidLabel(label string) bool {
	return strutil.EqualFoldAnyTrim(label,
		constants.ThemeBright,
		constants.ThemeStandard,
		constants.ThemeMonochrome,
		constants.ThemeMono,
		ThemeNavyPurple,
		ThemeTokyoNight,
		ThemeVSCodeDarkPlus,
		ThemeMonokaiPro)
}
