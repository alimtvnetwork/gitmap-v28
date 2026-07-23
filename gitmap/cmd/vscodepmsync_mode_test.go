// Package cmd — vscodepmsync_mode_test.go: regression coverage for
// the v4.37.0 `--mode` flag added to `gitmap vscode-pm-sync`.
//
// We exercise all three modes against a fixture projects.json:
//
//  1. UNION (default, no flag)        — existing ∪ detected, dedup'd.
//  2. REPLACE (`--mode replace`)       — detector wins, user tags dropped,
//     brand survives because the detector
//     always pre-pends it.
//  3. INTERSECTION (`--mode intersection`)
//     — only tags in BOTH sources survive,
//     plus the gitmap brand is PINNED
//     (added back even when intersection
//     is empty) per the v4.37.0 contract.
//
// The "bad mode" path (unknown literal) is covered by the package-
// level ParseMergeMode unit test in vscodepm/mergemode_test.go;
// here we only assert the runner-side behavior of the three valid
// modes end-to-end against a real on-disk file.
//
// All tests skip on GOOS=windows (XDG path layout). Windows path
// resolution is covered separately in vscodepm/path_test.go.
package cmd

import (
	"runtime"
	"testing"
)

// TestVSCodePMSyncModeUnionDefaultPreservesUserTags asserts the
// no-flag default still behaves like v4.36.0: user-added tags are
// kept and the brand tag is added.
func TestVSCodePMSyncModeUnionDefaultPreservesUserTags(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	_, restore := setupVSCodePMSyncFixtureWithTags(t,
		[]string{"user-only-tag"})
	defer restore()

	runVSCodePMSync(nil)

	got := loadVSCodePMSyncProjectsJSON(t)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if !containsTag(got[0].Tags, "user-only-tag") {
		t.Errorf("union dropped user tag: %v", got[0].Tags)
	}
	if !containsTag(got[0].Tags, "gitmap") {
		t.Errorf("union missing brand tag: %v", got[0].Tags)
	}
}

// TestVSCodePMSyncModeReplaceDropsUserTags asserts `--mode replace`
// overwrites the existing tag set with the detector output. The
// "user-only-tag" must be GONE, but "gitmap" must remain because
// DetectTagsCustom always pre-pends it.
func TestVSCodePMSyncModeReplaceDropsUserTags(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	_, restore := setupVSCodePMSyncFixtureWithTags(t,
		[]string{"user-only-tag", "gitmap"})
	defer restore()

	runVSCodePMSync([]string{"--mode", "replace"})

	got := loadVSCodePMSyncProjectsJSON(t)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if containsTag(got[0].Tags, "user-only-tag") {
		t.Errorf("replace did not drop user tag: %v", got[0].Tags)
	}
	if !containsTag(got[0].Tags, "gitmap") {
		t.Errorf("replace dropped brand tag (detector should pre-pend it): %v",
			got[0].Tags)
	}
}

// TestVSCodePMSyncModeIntersectionDropsExclusiveTags asserts
// `--mode intersection` keeps only tags present in BOTH the existing
// on-disk set AND the freshly-detected set. A tag that exists only
// in one of the two sources must be GONE after the run.
//
// We seed with both "user-only-tag" (exists only in the file) and
// "gitmap" (exists in both — the detector always emits it). After
// the run only "gitmap" should survive.
func TestVSCodePMSyncModeIntersectionDropsExclusiveTags(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	_, restore := setupVSCodePMSyncFixtureWithTags(t,
		[]string{"user-only-tag", "gitmap"})
	defer restore()

	runVSCodePMSync([]string{"--mode", "intersection"})

	got := loadVSCodePMSyncProjectsJSON(t)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if containsTag(got[0].Tags, "user-only-tag") {
		t.Errorf("intersection retained tag that's only in existing set: %v",
			got[0].Tags)
	}
	if !containsTag(got[0].Tags, "gitmap") {
		t.Errorf("intersection lost brand tag (should always be pinned): %v",
			got[0].Tags)
	}
}

// TestVSCodePMSyncModeIntersectionPinsBrandWhenAbsent asserts the
// MOST IMPORTANT contract of the intersection mode: even when the
// existing on-disk set does NOT contain "gitmap", the runner must
// add it back. Otherwise a strict intersection would silently strip
// the brand from any entry the detector hadn't tagged before.
//
// Seed has only ["user-only-tag"] — no brand. After intersection
// (which would normally produce []) the brand-pin step adds gitmap.
func TestVSCodePMSyncModeIntersectionPinsBrandWhenAbsent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG path layout — windows path resolution covered elsewhere")
	}
	_, restore := setupVSCodePMSyncFixtureWithTags(t,
		[]string{"user-only-tag"})
	defer restore()

	runVSCodePMSync([]string{"--mode", "intersection"})

	got := loadVSCodePMSyncProjectsJSON(t)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if !containsTag(got[0].Tags, "gitmap") {
		t.Errorf("brand-pin failed under intersection: %v", got[0].Tags)
	}
	if countTag(got[0].Tags, "gitmap") != 1 {
		t.Errorf("brand pinned more than once under intersection: %v",
			got[0].Tags)
	}
}
