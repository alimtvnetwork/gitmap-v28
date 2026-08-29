// Package cmd — glyphsflag.go: global `--glyphs` switch parser.
//
// Mirrors stripThemeFlag exactly so the global-flag inventory stays
// homogeneous. Strips `--glyphs <mode>` / `--glyphs=<mode>` (and short
// `-glyphs` form) from os.Args, validates the value, and exports
// GITMAP_GLYPHS for the glyphs package and any subprocess gitmap spawns.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/glyphs"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// stripGlyphsFlag removes every `--glyphs` / `-glyphs` occurrence and
// returns the cleaned argv slice. Aborts with status 2 on a bad value
// to surface CI typos.
func stripGlyphsFlag(args []string) []string {
	short := "-" + constants.FlagGlyphs
	long := "--" + constants.FlagGlyphs

	out := make([]string, 0, len(args))
	chosen := ""
	i := 0
	for i < len(args) {
		a := args[i]
		val, consumed, matched := matchGlyphsArg(a, args, i, short, long)
		if !matched {
			out = append(out, a)
			i++

			continue
		}
		chosen = val
		i += consumed
	}

	if chosen != "" {
		applyGlyphsChoice(chosen)
	}

	return out
}

// matchGlyphsArg recognizes the four legal flag forms.
func matchGlyphsArg(
	a string,
	args []string,
	i int,
	short,
	long string,
) (val string, consumed int, matched bool) {
	matchBare := a == short || a == long
	if matchBare && i+1 < len(args) && glyphs.IsValidLabel(args[i+1]) {
		return args[i+1], 2, true
	}
	if matchBare {
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

// applyGlyphsChoice validates and exports the choice, or aborts.
func applyGlyphsChoice(choice string) {
	if !glyphs.IsValidLabel(choice) {
		fmt.Fprintf(os.Stderr,
			"gitmap: invalid --glyphs value %q (want: %s | %s | %s)\n",
			choice,
			constants.GlyphsAuto,
			constants.GlyphsRich,
			constants.GlyphsSafe,
		)
		cliexit.HandleError(nil, 2)
	}
	os.Setenv(constants.EnvGlyphs, choice)
}
