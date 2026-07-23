// Package cmd — themeflag.go: global `--theme` palette selector.
//
// Strips `--theme <mode>` / `--theme=<mode>` (and the short `-theme`
// form) from os.Args before subcommand dispatch and exports
// GITMAP_THEME so gitmap/theme.Install — and any subprocess gitmap
// spawns — picks up the choice. Mirrors stripVSCodeSyncDisabledFlag's
// pattern so the global-flag inventory stays homogeneous.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/theme"
)

// stripThemeFlag removes every `--theme` / `-theme` occurrence from
// args, validates its value, and sets GITMAP_THEME. Returns the
// cleaned argv slice. On an invalid value it writes a clear error to
// stderr and exits with status 2 — silently falling back would hide
// typos in CI scripts.
func stripThemeFlag(args []string) []string {
	short := "-" + constants.FlagTheme
	long := "--" + constants.FlagTheme

	out := make([]string, 0, len(args))
	chosen := ""
	i := 0
	for i < len(args) {
		a := args[i]
		val, consumed, matched := matchThemeArg(a, args, i, short, long)
		if !matched {
			out = append(out, a)
			i++

			continue
		}
		chosen = val
		i += consumed
	}

	if chosen != "" {
		applyThemeChoice(chosen)
	}

	return out
}

// matchThemeArg recognizes the four legal flag forms and returns the
// extracted value plus how many argv slots it consumed. The split
// `--theme <mode>` form only consumes the next slot when it looks
// like a theme label — otherwise a bare `--theme` would silently
// steal the subcommand name.
func matchThemeArg(a string, args []string, i int, short, long string) (val string, consumed int, matched bool) {
	if a == short || a == long {
		if i+1 < len(args) && theme.IsValidLabel(args[i+1]) {
			return args[i+1], 2, true
		}

		return "", 1, true
	}
	if v, ok := stripEqPrefix(a, long+"="); ok {
		return v, 1, true
	}
	if v, ok := stripEqPrefix(a, short+"="); ok {
		return v, 1, true
	}

	return "", 0, false
}

// stripEqPrefix returns the value after prefix when a starts with it.
func stripEqPrefix(a, prefix string) (string, bool) {
	if len(a) <= len(prefix) || a[:len(prefix)] != prefix {
		return "", false
	}

	return a[len(prefix):], true
}

// applyThemeChoice validates choice and exports GITMAP_THEME, or
// aborts with a friendly error listing the accepted values.
func applyThemeChoice(choice string) {
	if !theme.IsValidLabel(choice) {
		fmt.Fprintf(os.Stderr,
			"gitmap: invalid --theme value %q (want: %s | %s | %s)\n",
			choice,
			constants.ThemeBright,
			constants.ThemeStandard,
			constants.ThemeMonochrome,
		)
		os.Exit(2)
	}
	os.Setenv(constants.EnvTheme, choice)
}
