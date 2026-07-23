//go:build windows

package startup

// Windows-only test for --backend-scoped removal. Confirms that
// passing RemoveOptions.Backend = BackendRegistry only touches the
// Run-key entry (NOT the .lnk Startup folder entry), and vice
// versa. Without scoping, removeWindows falls back to the other
// backend on a NoOp — which would silently delete an entry the
// user did not target.

import (
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestRemoveWindows_BackendScopedRegistry creates entries in BOTH
// backends with the same name, then removes scoped to registry.
// Asserts the Run-key value is gone but the .lnk file is still
// on disk.
func TestRemoveWindows_BackendScopedRegistry(t *testing.T) {
	dir := withIsolatedAppData(t)
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	addBoth(t, name)

	res, err := RemoveWithOptions(name, RemoveOptions{Backend: BackendRegistry})
	if err != nil || res.Status != RemoveDeleted {
		t.Fatalf("scoped Remove(registry): %v / status=%d", err, res.Status)
	}
	lnk := dir + "\\" + constants.StartupWinValuePrefix + name + constants.StartupLnkExt
	if _, err := os.Stat(lnk); err != nil {
		t.Errorf(".lnk was incorrectly deleted: %v", err)
	}
}

// TestRemoveWindows_BackendScopedStartupFolder is the mirror test:
// scoped removal of the .lnk leaves the registry entry intact.
func TestRemoveWindows_BackendScopedStartupFolder(t *testing.T) {
	withIsolatedAppData(t)
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	addBoth(t, name)

	res, err := RemoveWithOptions(name, RemoveOptions{Backend: BackendStartupFolder})
	if err != nil || res.Status != RemoveDeleted {
		t.Fatalf("scoped Remove(folder): %v / status=%d", err, res.Status)
	}
	entries, _ := listWindowsRegistry()
	for _, e := range entries {
		if e.Name == constants.StartupWinValuePrefix+name {
			return // Registry entry survived as expected
		}
	}
	t.Errorf("registry entry was incorrectly deleted")
}

// addBoth is a tiny helper that creates the same logical name in
// both Windows backends so the scoped-removal tests have something
// to scope.
func addBoth(t *testing.T, name string) {
	t.Helper()
	if _, err := Add(AddOptions{Name: name, Exec: `C:\gitmap.exe watch`,
		Backend: BackendRegistry}); err != nil {
		t.Fatalf("Add registry: %v", err)
	}
	if _, err := Add(AddOptions{Name: name, Exec: `C:\gitmap.exe watch`,
		Backend: BackendStartupFolder}); err != nil {
		t.Fatalf("Add startup-folder: %v", err)
	}
}
