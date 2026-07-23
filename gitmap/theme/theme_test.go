package theme

import (
	"strings"
	"testing"
)

// TestParseRecognizesAllLabels locks the user-facing label table —
// renaming any of these strings breaks scripts in the wild that
// pass `--theme=...`.
func TestParseRecognizesAllLabels(t *testing.T) {
	cases := map[string]Mode{
		"bright":     ModeBright,
		"BRIGHT":     ModeBright,
		"standard":   ModeStandard,
		"  standard": ModeStandard,
		"monochrome": ModeMono,
		"mono":       ModeMono,
		"":           ModeBright, // fallback
		"weird":      ModeBright, // unknown → fallback
	}
	for label, want := range cases {
		if got := Parse(label); got != want {
			t.Errorf("Parse(%q) = %v, want %v", label, got, want)
		}
	}
}

// TestFilterBrightIsPassthrough guards the zero-cost default: bright
// mode MUST return the input slice unchanged so we don't burn CPU on
// the hot path.
func TestFilterBrightIsPassthrough(t *testing.T) {
	in := []byte("\033[1;92m✓ ok\033[0m\n")
	got := Filter(in, ModeBright)
	if &got[0] != &in[0] {
		t.Fatal("ModeBright must return the input slice unchanged (no copy)")
	}
}

// TestFilterMonoStripsAllSGR is the contract for piping / log
// scraping: every SGR escape disappears, payload bytes survive.
func TestFilterMonoStripsAllSGR(t *testing.T) {
	in := []byte("\033[1;92m✓\033[0m hello \033[38;5;208m\033[1mworld\033[0m")
	got := string(Filter(in, ModeMono))
	if strings.Contains(got, "\033[") {
		t.Fatalf("mono mode leaked an escape: %q", got)
	}
	if got != "✓ hello world" {
		t.Fatalf("mono payload corrupted: %q", got)
	}
}

// TestFilterStandardDowngrades verifies the bright→plain mapping the
// pre-v5.13 look depended on.
func TestFilterStandardDowngrades(t *testing.T) {
	in := []byte("\033[1;92m✓\033[0m \033[1;96mcyan\033[0m")
	got := string(Filter(in, ModeStandard))
	want := "\033[32m✓\033[0m \033[36mcyan\033[0m"
	if got != want {
		t.Fatalf("standard mode mismatch:\n got: %q\nwant: %q", got, want)
	}
}

// TestFilterPreservesNonSGRCSI lets cursor / clear escapes
// (e.g. \033[2J, \033[H used by gitmap watch) pass through untouched
// in every mode — the rewriter only owns SGR.
func TestFilterPreservesNonSGRCSI(t *testing.T) {
	in := []byte("\033[2J\033[Hredraw")
	for _, m := range []Mode{ModeStandard, ModeMono} {
		got := string(Filter(in, m))
		if got != "\033[2J\033[Hredraw" {
			t.Fatalf("mode %v mangled non-SGR CSI: %q", m, got)
		}
	}
}
