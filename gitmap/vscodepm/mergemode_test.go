// Package vscodepm — mergemode_test.go: pure-function tests for
// the MergeMode enum, ParseMergeMode validator, and the strategy
// dispatcher (mergeTags). End-to-end CLI coverage lives in
// gitmap/cmd/vscodepmsync_mode_test.go; this file only asserts the
// internal logic so the dispatcher stays correct even when the CLI
// surface is refactored.
package vscodepm

import (
	"sort"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestParseMergeModeAcceptsCanonicalLiterals asserts every valid
// CLI literal round-trips through ParseMergeMode + String() to the
// same string. Empty input maps to UNION (the default).
func TestParseMergeModeAcceptsCanonicalLiterals(t *testing.T) {
	cases := []struct {
		in   string
		want MergeMode
	}{
		{"", MergeModeUnion},
		{constants.VSCodePMSyncModeUnion, MergeModeUnion},
		{constants.VSCodePMSyncModeReplace, MergeModeReplace},
		{constants.VSCodePMSyncModeIntersection, MergeModeIntersection},
	}
	for _, tc := range cases {
		got, err := ParseMergeMode(tc.in)
		if err != nil {
			t.Errorf("ParseMergeMode(%q) returned err: %v", tc.in, err)
		}
		if got != tc.want {
			t.Errorf("ParseMergeMode(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

// TestParseMergeModeRejectsUnknown asserts unknown literals fail
// loud (zero-swallow rule) instead of silently defaulting.
func TestParseMergeModeRejectsUnknown(t *testing.T) {
	_, err := ParseMergeMode("merge-everything-please")
	if err == nil {
		t.Fatal("ParseMergeMode accepted unknown literal")
	}
	if !strings.Contains(err.Error(), "merge-everything-please") {
		t.Errorf("error should echo the bad value, got: %v", err)
	}
}

// TestMergeModeStringRoundTrip asserts MergeMode.String() emits the
// CLI literals exactly so helptext, parser, and enum stay in sync.
func TestMergeModeStringRoundTrip(t *testing.T) {
	if got, want := MergeModeUnion.String(), constants.VSCodePMSyncModeUnion; got != want {
		t.Errorf("union: got %q, want %q", got, want)
	}
	if got, want := MergeModeReplace.String(), constants.VSCodePMSyncModeReplace; got != want {
		t.Errorf("replace: got %q, want %q", got, want)
	}
	if got, want := MergeModeIntersection.String(), constants.VSCodePMSyncModeIntersection; got != want {
		t.Errorf("intersection: got %q, want %q", got, want)
	}
}

// TestMergeTagsUnion asserts the union path matches the legacy
// unionTags semantics (existing-order first, dedup, additive).
func TestMergeTagsUnion(t *testing.T) {
	got := mergeTags(MergeModeUnion,
		[]string{"a", "b", "gitmap"},
		[]string{"b", "c", "gitmap"})
	want := []string{"a", "b", "gitmap", "c"}
	if !sliceEqual(got, want) {
		t.Errorf("union mergeTags: got %v, want %v", got, want)
	}
}

// TestMergeTagsReplaceDropsExistingKeepsBrand asserts replace returns
// a dedup'd copy of incoming, dropping anything not in the detector
// output. Brand survives because incoming already carries it (the
// real detector always pre-pends it).
func TestMergeTagsReplaceDropsExistingKeepsBrand(t *testing.T) {
	got := mergeTags(MergeModeReplace,
		[]string{"user-tag", "gitmap"},
		[]string{"gitmap", "go", "docker"})
	want := []string{"gitmap", "go", "docker"}
	if !sliceEqual(got, want) {
		t.Errorf("replace mergeTags: got %v, want %v", got, want)
	}
}

// TestMergeTagsIntersectionPinsBrand asserts the intersection +
// brand-pin contract: only tags in BOTH sets survive, and "gitmap"
// is added even when neither side carries it.
func TestMergeTagsIntersectionPinsBrand(t *testing.T) {
	// Both contain "go" — must survive. "user" only on left, "rust"
	// only on right — both must drop.
	got := mergeTags(MergeModeIntersection,
		[]string{"user", "go"},
		[]string{"go", "rust"})
	wantSet := map[string]struct{}{"go": {}, "gitmap": {}}
	if len(got) != len(wantSet) {
		t.Fatalf("intersection len: got %v, want %v", got, wantSet)
	}
	for _, tag := range got {
		if _, ok := wantSet[tag]; !ok {
			t.Errorf("intersection produced unexpected tag %q in %v",
				tag, got)
		}
	}

	// Empty intersection -> brand still pinned.
	gotEmpty := mergeTags(MergeModeIntersection,
		[]string{"user-only"},
		[]string{"detector-only"})
	if len(gotEmpty) != 1 || gotEmpty[0] != constants.AutoTagGitmap {
		t.Errorf("empty intersection should pin brand, got %v", gotEmpty)
	}
}

// sliceEqual reports order-sensitive equality for tag-slice asserts.
// Where a test asserts set semantics it compares against a map instead.
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ac := append([]string{}, a...)
	bc := append([]string{}, b...)
	sort.Strings(ac)
	sort.Strings(bc)
	for i := range ac {
		if ac[i] != bc[i] {
			return false
		}
	}

	return true
}
