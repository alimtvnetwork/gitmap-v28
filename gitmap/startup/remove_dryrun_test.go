package startup

// Tests for RemoveWithOptions(..., DryRun:true). Verifies the
// classification logic matches the live Remove path EXACTLY (same
// status, same Path, DryRun=true) while the filesystem stays
// untouched. Linux fixtures only — same skip as startup_test.go;
// the macOS dry-run guarantee is implicitly covered because the
// dry-run branch bypasses the OS-specific delete entirely.

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRemove_DryRun_DeletedReportedButFileRemains is the headline
// guarantee: a managed file produces RemoveDeleted with DryRun=true
// and the file is STILL on disk after the call.
func TestRemove_DryRun_DeletedReportedButFileRemains(t *testing.T) {
	dir := withFakeAutostartDir(t)
	managedPath := writeDesktop(t, dir, "gitmap-keep.desktop", true, "/usr/bin/x")

	res, err := RemoveWithOptions("gitmap-keep", RemoveOptions{DryRun: true})
	if err != nil {
		t.Fatalf("dry-run remove: %v", err)
	}
	if res.Status != RemoveDeleted {
		t.Fatalf("status = %v, want RemoveDeleted (dry-run preview)", res.Status)
	}
	if !res.DryRun {
		t.Fatalf("DryRun = false, want true")
	}
	if res.Path != managedPath {
		t.Fatalf("Path = %q, want %q", res.Path, managedPath)
	}
	if _, statErr := os.Stat(managedPath); statErr != nil {
		t.Fatalf("file must remain on disk after dry-run, stat err=%v", statErr)
	}
}

// TestRemove_DryRun_RefusedDoesNotTouchThirdParty mirrors the live
// Refused branch: a non-gitmap file produces RemoveRefused with
// DryRun=true and is not deleted.
func TestRemove_DryRun_RefusedDoesNotTouchThirdParty(t *testing.T) {
	dir := withFakeAutostartDir(t)
	thirdPartyPath := writeDesktop(t, dir, "gitmap-thirdparty.desktop", false, "/y")

	res, _ := RemoveWithOptions("gitmap-thirdparty", RemoveOptions{DryRun: true})
	if res.Status != RemoveRefused || !res.DryRun {
		t.Fatalf("got status=%v dryRun=%v, want RemoveRefused/true", res.Status, res.DryRun)
	}
	if _, statErr := os.Stat(thirdPartyPath); statErr != nil {
		t.Fatalf("third-party file must remain, stat err=%v", statErr)
	}
}

// TestRemove_DryRun_NoOpAndBadNamePropagateFlag confirms the
// non-mutating outcomes (NoOp, BadName) also carry DryRun=true so
// renderers can pick the `(dry-run)` mirror message uniformly.
func TestRemove_DryRun_NoOpAndBadNamePropagateFlag(t *testing.T) {
	withFakeAutostartDir(t)

	res, _ := RemoveWithOptions("does-not-exist", RemoveOptions{DryRun: true})
	if res.Status != RemoveNoOp || !res.DryRun {
		t.Fatalf("noop: got status=%v dryRun=%v", res.Status, res.DryRun)
	}

	res, _ = RemoveWithOptions("../../etc/passwd", RemoveOptions{DryRun: true})
	if res.Status != RemoveBadName || !res.DryRun {
		t.Fatalf("badname: got status=%v dryRun=%v", res.Status, res.DryRun)
	}
}

// TestRemove_LiveCallStillDeletes is a regression guard: adding the
// dry-run branch must not have broken the live delete path. We run
// the SAME setup as the dry-run test above and assert the file IS
// gone after a non-dry-run call.
func TestRemove_LiveCallStillDeletes(t *testing.T) {
	dir := withFakeAutostartDir(t)
	managedPath := writeDesktop(t, dir, "gitmap-live.desktop", true, "/x")

	res, err := RemoveWithOptions("gitmap-live", RemoveOptions{})
	if err != nil || res.Status != RemoveDeleted || res.DryRun {
		t.Fatalf("live: status=%v dryRun=%v err=%v", res.Status, res.DryRun, err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "gitmap-live.desktop")); !os.IsNotExist(statErr) {
		t.Fatalf("file should be gone, stat err=%v", statErr)
	}
	_ = managedPath
}
