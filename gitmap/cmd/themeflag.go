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

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
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
	cleaned := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		a := args[i]
		if val, ok := parseThemeEqual(a, short, long); ok {
			applyThemeChoice(val)
			continue
		}
		if (a == short || a == long) && i+1 < len(args) {
			applyThemeChoice(args[i+1])
			i++
			continue
		}
		if a == short || a == long {
			applyThemeChoice("")
			continue
		}
		cleaned = append(cleaned, a)
	}

	return cleaned
}

// parseThemeEqual checks for `-theme=val` or `--theme=val`.
func parseThemeEqual(a, short, long string) (string, bool) {
	if val, ok := stripThemePrefix(a, short+"="); ok {
		return val, true
	}
	return stripThemePrefix(a, long+"=")
}

// stripThemePrefix helper returns the remainder if a begins with prefix.
func stripThemePrefix(a, prefix string) (string, bool) {
	if len(a) < len(prefix) || a[:len(prefix)] != prefix {
		return "", false
	}

	return a[len(prefix):], true
}

// applyThemeChoice validates choice and exports GITMAP_THEME, or
// aborts with a friendly error listing the accepted values.
func applyThemeChoice(choice string) {
	if !theme.IsValidLabel(choice) {
		err := apperror.NewWithDetails(
			"cmd.themeflag.apply",
			"E1152",
			fmt.Sprintf("invalid --theme value %q (want: %s | %s | %s)",
				choice,
				constants.ThemeBright,
				constants.ThemeStandard,
				constants.ThemeMonochrome),
			"cmd.themeflag",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"choice": choice},
		)
		cliexit.HandleError(err, 2)
	}
	os.Setenv(constants.EnvTheme, choice)
}
