package cmd

import (
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
)

// TestParsePrettyFlagDefaultsToAuto guards the no-flag path — every
// command that adopts the flag must keep its existing TTY-auto behavior
// unchanged for users who never type --pretty.
func TestParsePrettyFlagDefaultsToAuto(t *testing.T) {
	rest, mode := ParsePrettyFlag([]string{"foo", "bar"})
	if mode != render.PrettyAuto {
		t.Errorf("mode = %v, want PrettyAuto", mode)
	}
	if !reflect.DeepEqual(rest, []string{"foo", "bar"}) {
		t.Errorf("rest = %v, want unchanged passthrough", rest)
	}
}

// TestParsePrettyFlagAcceptsAllPositiveForms locks the synonym set so a
// future cleanup can't silently drop one (e.g. shell aliases that pass
// --pretty=on rely on this).
func TestParsePrettyFlagAcceptsAllPositiveForms(t *testing.T) {
	for _, arg := range []string{"--pretty", "--pretty=true", "--pretty=on", "--pretty=1", "--pretty=YES"} {
		_, mode := ParsePrettyFlag([]string{arg})
		if mode != render.PrettyOn {
			t.Errorf("%s → %v, want PrettyOn", arg, mode)
		}
	}
}

// TestParsePrettyFlagAcceptsAllNegativeForms locks the negation set,
// including the dedicated --no-pretty alias for callers who prefer it
// to the =false suffix style.
func TestParsePrettyFlagAcceptsAllNegativeForms(t *testing.T) {
	for _, arg := range []string{"--pretty=false", "--pretty=off", "--pretty=0", "--no-pretty", "--pretty=NO"} {
		_, mode := ParsePrettyFlag([]string{arg})
		if mode != render.PrettyOff {
			t.Errorf("%s → %v, want PrettyOff", arg, mode)
		}
	}
}

// TestParsePrettyFlagExplicitAutoResets covers the niche but useful
// `--pretty=auto` case — handy when scripting wants to override an
// upstream alias that hard-codes --pretty=false.
func TestParsePrettyFlagExplicitAutoResets(t *testing.T) {
	_, mode := ParsePrettyFlag([]string{"--pretty=false", "--pretty=auto"})
	if mode != render.PrettyAuto {
		t.Errorf("explicit auto reset → %v, want PrettyAuto", mode)
	}
}

// TestParsePrettyFlagLastWriterWins matches stdlib flag.Parse semantics
// so users can stack flags from multiple sources (env-prefilled +
// per-invocation override) without surprises.
func TestParsePrettyFlagLastWriterWins(t *testing.T) {
	_, mode := ParsePrettyFlag([]string{"--pretty=true", "--no-pretty"})
	if mode != render.PrettyOff {
		t.Fatalf("last flag should win: got %v, want PrettyOff", mode)
	}
}

// TestParsePrettyFlagStripsFromArgs is the load-bearing test for
// integration with downstream flag.FlagSet parsers — they must never
// see --pretty in their args, otherwise they'd choke with "flag
// provided but not defined".
func TestParsePrettyFlagStripsFromArgs(t *testing.T) {
	rest, _ := ParsePrettyFlag([]string{"--latest", "--pretty=false", "--limit", "3"})
	if !reflect.DeepEqual(rest, []string{"--latest", "--limit", "3"}) {
		t.Errorf("rest = %v, want pretty-flag stripped", rest)
	}
}

// TestParsePrettyFlagPassesThroughUnknownValues protects the escape
// hatch: an unrecognized value (typo / future extension) must surface
// through the downstream parser instead of being silently swallowed.
func TestParsePrettyFlagPassesThroughUnknownValues(t *testing.T) {
	rest, mode := ParsePrettyFlag([]string{"--pretty=maybe"})
	if mode != render.PrettyAuto {
		t.Errorf("unknown value → %v, want unchanged PrettyAuto", mode)
	}
	if !reflect.DeepEqual(rest, []string{"--pretty=maybe"}) {
		t.Errorf("rest = %v, want passthrough so downstream errors clearly", rest)
	}
}

// TestParsePrettyFlagAcceptsColorSynonyms locks the `--no-color` /
// `--color[=value]` aliases. `--no-color` is the conventional spelling
// users reach for first (mirrors NO_COLOR env), and the CLI's only
// ANSI surface is the pretty markdown pipeline, so the two flags
// share a single resolver.
func TestParsePrettyFlagAcceptsColorSynonyms(t *testing.T) {
	cases := []struct {
		arg  string
		want render.PrettyMode
	}{
		{"--no-color", render.PrettyOff},
		{"--color=false", render.PrettyOff},
		{"--color=off", render.PrettyOff},
		{"--color", render.PrettyOn},
		{"--color=true", render.PrettyOn},
		{"--color=on", render.PrettyOn},
		{"--color=auto", render.PrettyAuto},
	}
	for _, tc := range cases {
		rest, mode := ParsePrettyFlag([]string{tc.arg})
		if mode != tc.want {
			t.Errorf("%s → %v, want %v", tc.arg, mode, tc.want)
		}
		if len(rest) != 0 {
			t.Errorf("%s left residue in args: %v", tc.arg, rest)
		}
	}
}

// TestParsePrettyFlagColorAndPrettyInteract guards last-writer-wins
// across the synonym families — a script may compose flags from
// multiple sources and the final mode must reflect the last token
// regardless of which spelling it used.
func TestParsePrettyFlagColorAndPrettyInteract(t *testing.T) {
	_, mode := ParsePrettyFlag([]string{"--pretty", "--no-color"})
	if mode != render.PrettyOff {
		t.Fatalf("--pretty then --no-color → %v, want PrettyOff", mode)
	}
	_, mode = ParsePrettyFlag([]string{"--no-color", "--pretty=on"})
	if mode != render.PrettyOn {
		t.Fatalf("--no-color then --pretty=on → %v, want PrettyOn", mode)
	}
}

// TestParsePrettyFlagDoesNotEatColorPrefixedFlags protects unrelated
// flags that happen to share a prefix substring with `--color`
// (e.g. a hypothetical `--colorblind`) from being silently swallowed.
func TestParsePrettyFlagDoesNotEatColorPrefixedFlags(t *testing.T) {
	rest, mode := ParsePrettyFlag([]string{"--colorblind", "--color=blue"})
	if mode != render.PrettyAuto {
		t.Errorf("unknown color value → %v, want PrettyAuto", mode)
	}
	if !reflect.DeepEqual(rest, []string{"--colorblind", "--color=blue"}) {
		t.Errorf("rest = %v, want both flags passed through", rest)
	}
}
