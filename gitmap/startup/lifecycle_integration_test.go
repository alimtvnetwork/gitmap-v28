package startup

// End-to-end integration tests for the autostart lifecycle:
//
// 	Add → List → Add (idempotent) → Add --force → Remove →
// 	Remove (idempotent no-op) → List (empty)
//
// Why a dedicated file (vs piggy-backing on add_test.go /
// startup_test.go): the existing per-function tests verify each
// status code in isolation against hand-crafted fixtures. These
// tests verify the FULL contract — that the file Add writes is
// recognized by List and accepted by Remove, on both supported
// OSes — so a regression in any single layer (renderer, marker
// parser, name normalizer) surfaces as a lifecycle failure rather
// than slipping past one of the unit tests.
//
// Platform routing:
//   - Linux/Unix: uses withFakeAutostartDir (XDG_CONFIG_HOME
//     redirected to t.TempDir()).
//   - macOS:      uses withFakeLaunchAgentsDir ($HOME redirected).
//   - Windows:    skipped (the package's runtime guard returns
//     ErrStartupUnsupportedOS for the file-based AutostartDir; the
//     Registry/Startup-folder backends have their own integration
//     coverage in windows_test.go).
//
// Both helpers already live in the package's existing test files,
// so this file adds zero new fixture plumbing — just the lifecycle
// orchestration on top.

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// withLifecycleAutostartDir picks the right fake-dir helper for the
// current OS so the same lifecycle test body runs on Linux + macOS
// without per-test build tags. Returns the absolute autostart dir.
func withLifecycleAutostartDir(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("file-based lifecycle tests skip Windows; see windows_test.go")
	}
	if runtime.GOOS == "darwin" {
		return withFakeLaunchAgentsDir(t)
	}

	return withFakeAutostartDir(t)
}

// lifecycleNames returns the (logical, on-disk-basename) pair for
// the current OS. The on-disk basename is what List surfaces in
// Entry.Name and what Remove accepts as input — Linux prefixes with
// "gitmap-", macOS with "gitmap.".
func lifecycleNames(logical string) (entryName, fileBase string) {
	if runtime.GOOS == "darwin" {
		return "gitmap." + logical, "gitmap." + logical + ".plist"
	}

	return "gitmap-" + logical, "gitmap-" + logical + ".desktop"
}

// TestLifecycle_AddListRemove walks the canonical happy path and
// asserts each stage's contract. One big test (not five tiny ones)
// because the value here is the SEQUENCE — splitting it up would
// hide ordering regressions (e.g., Remove succeeding only because a
// previous Add wrote to the wrong path).
func TestLifecycle_AddListRemove(t *testing.T) {
	dir := withLifecycleAutostartDir(t)
	entryName, fileBase := lifecycleNames("watcher")
	exec := "/usr/local/bin/gitmap watch"

	// 1. Add (fresh) → AddCreated, file on disk.
	first, err := Add(AddOptions{Name: "watcher", Exec: exec})
	if err != nil {
		t.Fatalf("first Add: %v", err)
	}
	if first.Status != AddCreated {
		t.Fatalf("first Add status = %d, want AddCreated", first.Status)
	}
	wantPath := filepath.Join(dir, fileBase)
	if first.Path != wantPath {
		t.Fatalf("first Add path = %s, want %s", first.Path, wantPath)
	}
	if _, statErr := os.Stat(wantPath); statErr != nil {
		t.Fatalf("file not written: %v", statErr)
	}

	// 2. List → exactly one managed entry, Exec round-trips.
	entries, err := List()
	if err != nil {
		t.Fatalf("List after Add: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("List returned %d entries, want 1: %#v", len(entries), entries)
	}
	if entries[0].Name != entryName {
		t.Errorf("entry name = %s, want %s", entries[0].Name, entryName)
	}
	if !strings.Contains(entries[0].Exec, "/usr/local/bin/gitmap") {
		t.Errorf("entry exec = %q, want it to contain the binary path", entries[0].Exec)
	}

	// 3. Add (same name, no --force) → AddExists, body unchanged.
	second, err := Add(AddOptions{Name: "watcher", Exec: "/different/binary"})
	if err != nil {
		t.Fatalf("second Add: %v", err)
	}
	if second.Status != AddExists {
		t.Fatalf("second Add status = %d, want AddExists", second.Status)
	}
	body, _ := os.ReadFile(wantPath)
	if !strings.Contains(string(body), "/usr/local/bin/gitmap") {
		t.Errorf("idempotent re-add modified body:\n%s", body)
	}
	if strings.Contains(string(body), "/different/binary") {
		t.Errorf("idempotent re-add wrote new Exec:\n%s", body)
	}

	// 4. Add (--force) → AddOverwritten, body updated.
	forced, err := Add(AddOptions{Name: "watcher", Exec: "/different/binary", Force: true})
	if err != nil {
		t.Fatalf("forced Add: %v", err)
	}
	if forced.Status != AddOverwritten {
		t.Fatalf("forced Add status = %d, want AddOverwritten", forced.Status)
	}
	body, _ = os.ReadFile(wantPath)
	if !strings.Contains(string(body), "/different/binary") {
		t.Errorf("forced Add did not replace body:\n%s", body)
	}

	// 5. Remove (managed) → RemoveDeleted, file gone.
	rm, err := Remove(entryName)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if rm.Status != RemoveDeleted {
		t.Fatalf("Remove status = %d, want RemoveDeleted", rm.Status)
	}
	if _, statErr := os.Stat(wantPath); !os.IsNotExist(statErr) {
		t.Fatalf("file still present after Remove: %v", statErr)
	}

	// 6. Remove (same name again) → RemoveNoOp. This is the
	//    idempotency guarantee provisioning scripts depend on.
	rm2, err := Remove(entryName)
	if err != nil {
		t.Fatalf("second Remove: %v", err)
	}
	if rm2.Status != RemoveNoOp {
		t.Fatalf("second Remove status = %d, want RemoveNoOp", rm2.Status)
	}

	// 7. List → empty again.
	entries, err = List()
	if err != nil {
		t.Fatalf("final List: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("final List = %d entries, want 0: %#v", len(entries), entries)
	}
}

// TestLifecycle_MultipleEntriesIndependent verifies the lifecycle
// works correctly when several managed entries coexist: removing
// one must NOT touch the others, and List must still surface the
// survivors. This catches a class of regressions where a name-
// matching bug in Remove (e.g., prefix vs exact match) would
// silently delete more than requested.
func TestLifecycle_MultipleEntriesIndependent(t *testing.T) {
	withLifecycleAutostartDir(t)
	for _, name := range []string{"alpha", "beta", "gamma"} {
		res, err := Add(AddOptions{Name: name, Exec: "/bin/" + name})
		if err != nil || res.Status != AddCreated {
			t.Fatalf("Add %s: status=%d err=%v", name, res.Status, err)
		}
	}
	if entries, err := List(); err != nil || len(entries) != 3 {
		t.Fatalf("after 3 Adds: len=%d err=%v", len(entries), err)
	}
	betaName, _ := lifecycleNames("beta")
	rm, err := Remove(betaName)
	if err != nil || rm.Status != RemoveDeleted {
		t.Fatalf("Remove beta: status=%d err=%v", rm.Status, err)
	}
	entries, err := List()
	if err != nil {
		t.Fatalf("List after Remove: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 survivors, got %d: %#v", len(entries), entries)
	}
	for _, e := range entries {
		if e.Name == betaName {
			t.Errorf("beta should be gone, still present: %#v", e)
		}
	}
}

// TestLifecycle_RefuseThirdPartyAcrossOps verifies the
// "managed-only, never escalate" guarantee survives a full
// lifecycle: a third-party file with the gitmap- prefix that
// existed BEFORE Add must be refused on Add (even with --force) and
// must NOT appear in List or be deleted by Remove.
func TestLifecycle_RefuseThirdPartyAcrossOps(t *testing.T) {
	dir := withLifecycleAutostartDir(t)
	_, fileBase := lifecycleNames("intruder")
	thirdParty := filepath.Join(dir, fileBase)
	if err := os.WriteFile(thirdParty, []byte("not-managed-by-gitmap"), 0o644); err != nil {
		t.Fatalf("seed third-party file: %v", err)
	}
	res, err := Add(AddOptions{Name: "intruder", Exec: "/x", Force: true})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddRefused {
		t.Fatalf("Add status = %d, want AddRefused", res.Status)
	}
	body, _ := os.ReadFile(thirdParty)
	if string(body) != "not-managed-by-gitmap" {
		t.Errorf("third-party file modified: %q", body)
	}
	entries, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("List surfaced un-managed file: %#v", entries)
	}
	intruderName, _ := lifecycleNames("intruder")
	rm, err := Remove(intruderName)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if rm.Status != RemoveRefused {
		t.Fatalf("Remove status = %d, want RemoveRefused", rm.Status)
	}
	if _, statErr := os.Stat(thirdParty); statErr != nil {
		t.Fatalf("third-party file deleted: %v", statErr)
	}
}
