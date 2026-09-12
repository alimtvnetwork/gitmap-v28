package visibility

import (
	"reflect"
	"testing"
)

func TestAutoFixVDigitPatterns(t *testing.T) {
	pats, _ := ParsePatternList("macro-51, gitmap")
	names := []string{"macro-v51", "gitmap-v28", "other"}
	fixed := AutoFixVDigitPatterns(pats, names)
	if len(fixed) != 2 {
		t.Fatalf("expected 2 fixed patterns, got %d", len(fixed))
	}

	if fixed[0].Raw != "macro-v51" || fixed[1].Raw != "gitmap-v28" {
		t.Fatalf("unexpected fixed patterns: %+v", fixed)
	}
}

func TestNearMisses(t *testing.T) {
	pats, _ := ParsePatternList("gitmapp")
	names := []string{"gitmap", "foo", "bar"}
	near := NearMisses(pats, names, 3, 2)
	want := []string{"gitmap"}
	if !reflect.DeepEqual(near, want) {
		t.Fatalf("got %v, want %v", near, want)
	}
}

func TestLevenshtein(t *testing.T) {
	if d := levenshtein("abc", "abc"); d != 0 {
		t.Fatalf("expected 0, got %d", d)
	}

	if d := levenshtein("", "abc"); d != 3 {
		t.Fatalf("expected 3, got %d", d)
	}

	if d := levenshtein("abc", ""); d != 3 {
		t.Fatalf("expected 3, got %d", d)
	}

	if d := levenshtein("kitten", "sitting"); d != 3 {
		t.Fatalf("expected 3, got %d", d)
	}
}
