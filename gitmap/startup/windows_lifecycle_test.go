//go:build windows

package startup

// Windows full-lifecycle integration tests for `gitmap startup-add` /
// `startup-list` / `startup-remove`.
//
// Why this file (vs windows_test.go / winbackend_scoped_test.go):
//
//   - windows_test.go covers each backend's stages in isolation
//     (Add round-trip, third-party refuse, Remove no-op).
//   - winbackend_scoped_test.go covers --backend-scoped removal.
//   - Neither walks the FULL canonical sequence
//
//       Add → List → Add (idempotent) → Add --force →
//       Remove → Remove (idempotent no-op) → List (empty)
//
//     in one test, on the Windows surface. That sequence is exactly
//     what provisioning scripts depend on, and is what
//     lifecycle_integration_test.go guarantees on Linux/macOS. This
//     file gives Windows the same level of integration coverage.
//
// Isolation strategy:
//
//   - Registry backend: the package's registry constants are not
//     overrideable without invasive refactoring, so we touch real
//     HKCU under per-test unique value names (`uniqueName(t)` already
//     namespaces by t.TempDir()) and rely on t.Cleanup +
//     cleanupRegistryEntry to leave the host pristine even when a
//     test fails mid-way. The pattern is already established in
//     windows_test.go — this file reuses the helpers verbatim.
//   - .lnk backend: APPDATA is redirected to t.TempDir() via
//     withIsolatedAppData so .lnk files land in an ephemeral tree.
//     The tracking subkey under HKCU\Software\Gitmap\StartupFolder
//     IS still on real HKCU and is cleaned up the same way.
//
// All tests skip cleanly when run outside Windows (build tag) and
// the .lnk lifecycle test additionally skips when powershell.exe is
// absent (matching windows_test.go's TestAddWindowsStartupFolder_*
// posture).

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestWindowsLifecycle_RegistryFullSequence walks the canonical
// add → list → idempotent-add → forced-add → remove → idempotent-
// remove → list-empty path against the Registry backend. One big
// test (not seven tiny ones) because the value here is the SEQUENCE
// — splitting it would hide ordering regressions like Remove only
// succeeding because a previous Add wrote to the wrong value name.
func TestWindowsLifecycle_RegistryFullSequence(t *testing.T) {
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	valueName := constants.StartupWinValuePrefix + name
	originalExec := `C:\gitmap.exe watch`
	updatedExec := `C:\gitmap.exe watch --quiet`

	// 1. Add (fresh) → AddCreated, value visible in Run key.
	first, err := Add(AddOptions{Name: name, Exec: originalExec, Backend: BackendRegistry})
	if err != nil {
		t.Fatalf("first Add: %v", err)
	}
	if first.Status != AddCreated {
		t.Fatalf("first Add status = %d, want AddCreated", first.Status)
	}
	assertRegistryListContains(t, valueName, originalExec)

	// 2. Add (same name, no --force) → AddExists, value unchanged.
	idem, err := Add(AddOptions{Name: name, Exec: updatedExec, Backend: BackendRegistry})
	if err != nil {
		t.Fatalf("idempotent Add: %v", err)
	}
	if idem.Status != AddExists {
		t.Fatalf("idempotent Add status = %d, want AddExists", idem.Status)
	}
	assertRegistryListContains(t, valueName, originalExec)

	// 3. Add (--force) → AddOverwritten, value updated to new exec.
	forced, err := Add(AddOptions{Name: name, Exec: updatedExec,
		Backend: BackendRegistry, Force: true})
	if err != nil {
		t.Fatalf("forced Add: %v", err)
	}
	if forced.Status != AddOverwritten {
		t.Fatalf("forced Add status = %d, want AddOverwritten", forced.Status)
	}
	assertRegistryListContains(t, valueName, updatedExec)

	// 4. Remove (managed) → RemoveDeleted, value gone from Run key.
	rm, err := RemoveWithOptions(name, RemoveOptions{Backend: BackendRegistry})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if rm.Status != RemoveDeleted {
		t.Fatalf("Remove status = %d, want RemoveDeleted", rm.Status)
	}
	assertRegistryListMissing(t, valueName)

	// 5. Remove (same name again) → RemoveNoOp. This is the
	//    idempotency guarantee provisioning scripts depend on.
	rm2, err := RemoveWithOptions(name, RemoveOptions{Backend: BackendRegistry})
	if err != nil {
		t.Fatalf("second Remove: %v", err)
	}
	if rm2.Status != RemoveNoOp {
		t.Fatalf("second Remove status = %d, want RemoveNoOp", rm2.Status)
	}
}

// TestWindowsLifecycle_StartupFolderFullSequence is the .lnk-backend
// mirror of the registry sequence. APPDATA is redirected so the .lnk
// lands in t.TempDir(); the tracking subkey is cleaned up via
// cleanupRegistryEntry.
func TestWindowsLifecycle_StartupFolderFullSequence(t *testing.T) {
	if _, err := exec.LookPath("powershell.exe"); err != nil {
		t.Skip("powershell.exe not on PATH; .lnk backend cannot be exercised")
	}
	dir := withIsolatedAppData(t)
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	originalExec := `C:\Windows\System32\notepad.exe`
	updatedExec := `C:\Windows\System32\calc.exe`
	lnkPath := filepath.Join(dir,
		constants.StartupWinValuePrefix+name+constants.StartupLnkExt)

	// 1. Add (fresh) → AddCreated, .lnk on disk.
	first, err := Add(AddOptions{Name: name, Exec: originalExec,
		Backend: BackendStartupFolder})
	if err != nil {
		t.Fatalf("first Add: %v", err)
	}
	if first.Status != AddCreated {
		t.Fatalf("first Add status = %d, want AddCreated", first.Status)
	}
	if _, err := os.Stat(lnkPath); err != nil {
		t.Fatalf(".lnk not on disk after Add: %v", err)
	}
	assertStartupFolderListContains(t, name)

	// 2. Add (idempotent) → AddExists, .lnk still present.
	idem, err := Add(AddOptions{Name: name, Exec: updatedExec,
		Backend: BackendStartupFolder})
	if err != nil {
		t.Fatalf("idempotent Add: %v", err)
	}
	if idem.Status != AddExists {
		t.Fatalf("idempotent Add status = %d, want AddExists", idem.Status)
	}
	if _, err := os.Stat(lnkPath); err != nil {
		t.Fatalf(".lnk vanished after idempotent Add: %v", err)
	}

	// 3. Add (--force) → AddOverwritten, .lnk regenerated.
	forced, err := Add(AddOptions{Name: name, Exec: updatedExec,
		Backend: BackendStartupFolder, Force: true})
	if err != nil {
		t.Fatalf("forced Add: %v", err)
	}
	if forced.Status != AddOverwritten {
		t.Fatalf("forced Add status = %d, want AddOverwritten", forced.Status)
	}
	if _, err := os.Stat(lnkPath); err != nil {
		t.Fatalf(".lnk missing after forced Add: %v", err)
	}

	// 4. Remove → RemoveDeleted, .lnk deleted from disk.
	rm, err := RemoveWithOptions(name, RemoveOptions{Backend: BackendStartupFolder})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if rm.Status != RemoveDeleted {
		t.Fatalf("Remove status = %d, want RemoveDeleted", rm.Status)
	}
	if _, err := os.Stat(lnkPath); !os.IsNotExist(err) {
		t.Fatalf(".lnk still on disk after Remove: %v", err)
	}

	// 5. Remove again → RemoveNoOp.
	rm2, err := RemoveWithOptions(name, RemoveOptions{Backend: BackendStartupFolder})
	if err != nil {
		t.Fatalf("second Remove: %v", err)
	}
	if rm2.Status != RemoveNoOp {
		t.Fatalf("second Remove status = %d, want RemoveNoOp", rm2.Status)
	}
}

// TestWindowsLifecycle_StartupFolderRefusesThirdParty mirrors the
// existing TestAddWindowsRegistry_RefusesThirdParty for the .lnk
// backend: a hand-placed .lnk WITHOUT the tracking subkey must be
// refused on Add (even with --force) and survive a Remove call.
//
// Why this gap matters: the .lnk backend's "managed" classification
// depends on the tracking subkey under HKCU\Software\Gitmap\
// StartupFolder. A regression in classifyShortcut would let `gitmap
// startup-add --force` clobber a user's hand-made shortcut — exactly
// the data-loss scenario the marker contract exists to prevent.
func TestWindowsLifecycle_StartupFolderRefusesThirdParty(t *testing.T) {
	if _, err := exec.LookPath("powershell.exe"); err != nil {
		t.Skip("powershell.exe not on PATH; .lnk backend cannot be exercised")
	}
	dir := withIsolatedAppData(t)
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	thirdParty := filepath.Join(dir,
		constants.StartupWinValuePrefix+name+constants.StartupLnkExt)
	original := []byte("not-a-real-lnk-but-claims-the-filename")
	if err := os.WriteFile(thirdParty, original, 0o644); err != nil {
		t.Fatalf("seed third-party .lnk: %v", err)
	}

	res, err := Add(AddOptions{Name: name, Exec: `C:\safe.exe`,
		Backend: BackendStartupFolder, Force: true})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddRefused {
		t.Fatalf("Add status = %d, want AddRefused", res.Status)
	}
	body, err := os.ReadFile(thirdParty)
	if err != nil {
		t.Fatalf("re-read third-party .lnk: %v", err)
	}
	if string(body) != string(original) {
		t.Errorf("third-party .lnk modified: %q", body)
	}

	rm, err := RemoveWithOptions(name, RemoveOptions{Backend: BackendStartupFolder})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if rm.Status != RemoveRefused {
		t.Fatalf("Remove status = %d, want RemoveRefused", rm.Status)
	}
	if _, err := os.Stat(thirdParty); err != nil {
		t.Fatalf("third-party .lnk deleted by Remove: %v", err)
	}
}

// TestWindowsLifecycle_DualBackendListAggregation confirms the
// cross-backend List enumerates entries from BOTH backends in one
// call. Without this, a regression that drops one backend from
// listWindows() would silently break the unscoped `gitmap
// startup-list` output that users rely on as the single source of
// truth for "what will run at login".
func TestWindowsLifecycle_DualBackendListAggregation(t *testing.T) {
	if _, err := exec.LookPath("powershell.exe"); err != nil {
		t.Skip("powershell.exe not on PATH; .lnk backend cannot be exercised")
	}
	withIsolatedAppData(t)
	regName := uniqueName(t) + "-reg"
	folderName := uniqueName(t) + "-folder"
	t.Cleanup(func() { cleanupRegistryEntry(t, regName) })
	t.Cleanup(func() { cleanupRegistryEntry(t, folderName) })

	if _, err := Add(AddOptions{Name: regName, Exec: `C:\gitmap.exe scan`,
		Backend: BackendRegistry}); err != nil {
		t.Fatalf("Add registry: %v", err)
	}
	if _, err := Add(AddOptions{Name: folderName,
		Exec:    `C:\Windows\System32\notepad.exe`,
		Backend: BackendStartupFolder}); err != nil {
		t.Fatalf("Add startup-folder: %v", err)
	}

	entries, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	regWant := constants.StartupWinValuePrefix + regName
	folderWant := constants.StartupWinValuePrefix + folderName
	sawReg, sawFolder := false, false
	for _, e := range entries {
		if e.Name == regWant {
			sawReg = true
		}
		if e.Name == folderWant {
			sawFolder = true
		}
	}
	if !sawReg {
		t.Errorf("registry entry %q missing from cross-backend List: %#v",
			regWant, entries)
	}
	if !sawFolder {
		t.Errorf("startup-folder entry %q missing from cross-backend List: %#v",
			folderWant, entries)
	}
}

// assertRegistryListContains scans the registry-backend list for an
// entry with the exact value name and exec string. Fails the test
// (not a Fatal) so the calling lifecycle stage can surface BOTH a
// list-mismatch and a follow-up status mismatch in one run.
func assertRegistryListContains(t *testing.T, valueName, exec string) {
	t.Helper()
	entries, err := listWindowsRegistry()
	if err != nil {
		t.Fatalf("listWindowsRegistry: %v", err)
	}
	for _, e := range entries {
		if e.Name == valueName {
			if e.Exec != exec {
				t.Errorf("entry %q exec = %q, want %q", valueName, e.Exec, exec)
			}

			return
		}
	}
	t.Errorf("entry %q not in registry list: %#v", valueName, entries)
}

// assertRegistryListMissing is the negative counterpart used by the
// post-Remove stage. Fails when the value is still surfaced by List.
func assertRegistryListMissing(t *testing.T, valueName string) {
	t.Helper()
	entries, err := listWindowsRegistry()
	if err != nil {
		t.Fatalf("listWindowsRegistry: %v", err)
	}
	for _, e := range entries {
		if e.Name == valueName {
			t.Errorf("entry %q still in registry list after Remove", valueName)
		}
	}
}

// assertStartupFolderListContains confirms the .lnk backend's List
// surfaces the entry. Asserts on Name only (not Path) because the
// .lnk path includes APPDATA which the test redirected — keeping
// the check on the stable Name field avoids re-deriving the
// expected path here.
func assertStartupFolderListContains(t *testing.T, name string) {
	t.Helper()
	entries, err := listWindowsStartupFolder()
	if err != nil {
		t.Fatalf("listWindowsStartupFolder: %v", err)
	}
	want := constants.StartupWinValuePrefix + name
	for _, e := range entries {
		if e.Name == want {
			return
		}
	}
	t.Errorf("entry %q not in startup-folder list: %#v", want, entries)
}
