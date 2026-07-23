// Package cmd — vscodepmsync_dedupe_test.go: regression coverage for
// `gitmap vscode-pm-sync` data-shape edge cases that the original
// vscodepmsync_test.go didn't cover:
//
//  1. Tag DEDUPE — when the seed file already contains the "gitmap"
//     brand tag (or any other tag DetectTagsCustom would re-emit),
//     the on-disk array must keep ONE copy after sync, not two.
//  2. Multi-source UNION — pre-existing user-only tag PLUS detected
//     tags PLUS a duplicate of one of the detected tags collapses to
//     a single dedup'd set with stable order.
//  3. MISSING projects.json — runner must not crash; it should write
//     a brand-new (possibly empty) file rather than corrupting state
//     or refusing to run.
//  4. MALFORMED projects.json — runner must SOFT-FAIL: a single
//     stderr diagnostic, NO mutation of the bytes on disk, exit 0
//     so CI never breaks because someone hand-edited the file.
//
// All tests skip on GOOS=windows (XDG path layout). Windows path
// resolution is covered separately in vscodepm/path_test.go.
package cmd

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// TestVSCodePMSyncDedupesGitmapBrandTag asserts that when the seed
// projects.json ALREADY carries "gitmap" in the tags array (e.g. a
// previous sync wrote it, or the user added it by hand), running the
// command again does NOT append a second "gitmap" — the merge layer
// must dedupe. Regression guard for the obvious "every sync grows
// the tag list" failure mode.
func TestVSCodePMSyncDedupesGitmapBrandTag(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	repoDir, restore := setupVSCodePMSyncFixtureWithTags(t,
		[]string{"gitmap", "user"})
	defer restore()
	_ = repoDir

	runVSCodePMSync(nil)
	runVSCodePMSync(nil) // run twice — dedupe must hold across re-runs

	got := loadVSCodePMSyncProjectsJSON(t)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}

	if n := countTag(got[0].Tags, "gitmap"); n != 1 {
		t.Errorf("brand tag duplicated: got %d copies in %v, want exactly 1",
			n, got[0].Tags)
	}
	if n := countTag(got[0].Tags, "user"); n != 1 {
		t.Errorf("user tag duplicated: got %d copies in %v, want exactly 1",
			n, got[0].Tags)
	}
}

// TestVSCodePMSyncDedupesAcrossSources asserts that a tag arriving
// from BOTH the existing on-disk entry AND DetectTagsCustom collapses
// to a single occurrence. We seed with ["user", "gitmap", "user"] —
// note the in-array duplicate AND the brand tag — and assert the
// final array contains each label exactly once.
func TestVSCodePMSyncDedupesAcrossSources(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	_, restore := setupVSCodePMSyncFixtureWithTags(t,
		[]string{"user", "gitmap", "user"})
	defer restore()

	runVSCodePMSync(nil)

	got := loadVSCodePMSyncProjectsJSON(t)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}

	for _, tag := range got[0].Tags {
		if c := countTag(got[0].Tags, tag); c != 1 {
			t.Errorf("tag %q appears %d times in %v, want 1",
				tag, c, got[0].Tags)
		}
	}
	if !containsTag(got[0].Tags, "gitmap") {
		t.Errorf("brand tag missing from dedup'd set: %v", got[0].Tags)
	}
	if !containsTag(got[0].Tags, "user") {
		t.Errorf("user tag missing from dedup'd set: %v", got[0].Tags)
	}
}

// TestVSCodePMSyncMissingFileIsNotAnError asserts that running vpm
// against a HOME with no projects.json on disk does not crash and
// does not produce a malformed file. There are zero entries to walk,
// so the post-condition is just "ListEntries returns []Entry{} cleanly".
//
// This guards against a regression where readEntries would propagate
// os.ErrNotExist instead of treating it as an empty input set.
func TestVSCodePMSyncMissingFileIsNotAnError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	restore := setupVSCodePMSyncEmptyHome(t)
	defer restore()

	// Must not panic / os.Exit. Soft-handles missing-file case.
	runVSCodePMSync(nil)

	got, err := vscodepm.ListEntries()
	if err != nil {
		t.Fatalf("ListEntries after missing-file run: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 entries after missing-file run, got %d: %v",
			len(got), got)
	}
}

// TestVSCodePMSyncMalformedFileIsLeftUntouched asserts the spec'd
// soft-fail contract: when projects.json is not valid JSON, the
// runner reports the parse error to stderr (via reportVSCodePMSoftError)
// and EXITS WITHOUT MUTATING the file. The bytes on disk must be
// byte-identical before and after the run — never corrupt a hand-
// edited file the user is mid-way through fixing.
//
// This is the most important regression in this file: a previous
// implementation that "recovered" by writing []Entry{} would silently
// destroy the user's work.
func TestVSCodePMSyncMalformedFileIsLeftUntouched(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	path, restore := setupVSCodePMSyncMalformedFile(t)
	defer restore()

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read malformed file before run: %v", err)
	}

	runVSCodePMSync(nil)

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read malformed file after run: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Errorf("runner mutated malformed projects.json\nbefore: %q\nafter:  %q",
			before, after)
	}
}

// countTag returns how many times tag occurs in tags (string equality).
// Pulled out so the dedupe assertions can express "exactly N copies"
// without sort-and-compact gymnastics.
func countTag(tags []string, tag string) int {
	n := 0
	for _, t := range tags {
		if t == tag {
			n++
		}
	}

	return n
}
