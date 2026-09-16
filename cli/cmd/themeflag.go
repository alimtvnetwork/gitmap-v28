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

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/theme"
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
		nextIdx, isTheme := tryConsumeThemeArg(args, i, short, long)
		if isTheme {
			i = nextIdx
			continue
		}

		cleaned = append(cleaned, args[i])
	}

	return cleaned
}

// applyThemeChoice validates choice and exports GITMAP_THEME, or
// aborts with a friendly error listing the accepted values.
func applyThemeChoice(choice string) {
	if !theme.IsValidLabel(choice) {
		failInvalidThemeChoice(choice)
	}

	os.Setenv(constants.EnvTheme, choice)
}

func failInvalidThemeChoice(choice string) {
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
