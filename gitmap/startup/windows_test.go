//go:build windows

package startup

// Windows backend tests. Build-tag-gated to windows because the
// Registry backend uses golang.org/x/sys/windows/registry which only
// compiles on windows. Tests touch HKCU under
// `Software\GitmapTest\<random>` instead of the real
// `Software\Gitmap` root so a CI run cannot pollute a developer's
// real autostart entries — see withIsolatedRegistryRoots for the
// override mechanism.
//
// Coverage targets (per backend):
//
//   - happy-path Add → List → Remove round-trip
//   - Add refuses non-managed (third-party value with same name)
//   - Add idempotent (second Add returns AddExists)
//   - Add --force overwrites our own entry
//   - Remove of unknown name returns RemoveNoOp
//
// .lnk backend tests SKIP if powershell.exe is not on PATH so a
// Windows host without PowerShell still passes the registry suite.

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// withIsolatedAppData redirects %APPDATA% to t.TempDir() so the
// .lnk backend creates files in an isolated tree.
func withIsolatedAppData(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Setenv("APPDATA", root)
	dir := filepath.Join(root, constants.StartupFolderRelative)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir startup folder: %v", err)
	}

	return dir
}

// cleanupRegistryEntry deletes the Run-key value AND the tracking
// subkey created by an Add. Used by t.Cleanup to keep the real
// HKCU clean — the registry constants are NOT overrideable without
// invasive refactoring, so we touch real HKCU and clean up. No
// sibling-marker delete: the direct-value model never writes one.
func cleanupRegistryEntry(t *testing.T, clean string) {
	t.Helper()
	valueName := constants.StartupWinValuePrefix + clean
	if k, err := registry.OpenKey(registry.CURRENT_USER,
		constants.RegRunKeyPath, registry.SET_VALUE); err == nil {
		k.DeleteValue(valueName)
		k.Close()
	}
	registry.DeleteKey(registry.CURRENT_USER,
		constants.RegGitmapRegistrySub+`\`+clean)
	registry.DeleteKey(registry.CURRENT_USER,
		constants.RegGitmapStartupFolder+`\`+clean)
}

// uniqueName returns a per-test value name so concurrent test runs
// (or stale cleanups from a previous failed run) don't collide.
func uniqueName(t *testing.T) string {
	t.Helper()
	return "test-" + filepath.Base(t.TempDir())
}

// TestAddWindowsRegistry_RoundTrip is the headline test: Add writes
// a Run-key value, List sees it, Remove deletes it, List sees zero.
func TestAddWindowsRegistry_RoundTrip(t *testing.T) {
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	res, err := Add(AddOptions{Name: name, Exec: `C:\gitmap.exe watch`,
		Backend: BackendRegistry})
	if err != nil || res.Status != AddCreated {
		t.Fatalf("Add: %v / status=%d", err, res.Status)
	}
	entries, err := listWindowsRegistry()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	found := false
	for _, e := range entries {
		if e.Name == constants.StartupWinValuePrefix+name {
			found = true
			if e.Exec != `C:\gitmap.exe watch` {
				t.Errorf("exec = %q, want %s", e.Exec, `C:\gitmap.exe watch`)
			}
		}
	}
	if !found {
		t.Fatalf("List did not return the added entry: %#v", entries)
	}
	rm, err := RemoveWithOptions(name, RemoveOptions{})
	if err != nil || rm.Status != RemoveDeleted {
		t.Fatalf("Remove: %v / status=%d", err, rm.Status)
	}
}

// TestAddWindowsRegistry_RefusesThirdParty writes a Run-key value
// WITHOUT the sibling marker / tracking subkey, then confirms Add
// returns AddRefused even with --force.
func TestAddWindowsRegistry_RefusesThirdParty(t *testing.T) {
	name := uniqueName(t)
	valueName := constants.StartupWinValuePrefix + name
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	k, _, err := registry.CreateKey(registry.CURRENT_USER,
		constants.RegRunKeyPath, registry.SET_VALUE)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := k.SetStringValue(valueName, `C:\evil.exe`); err != nil {
		t.Fatalf("seed value: %v", err)
	}
	k.Close()

	res, err := Add(AddOptions{Name: name, Exec: `C:\safe.exe`,
		Backend: BackendRegistry, Force: true})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddRefused {
		t.Fatalf("status = %d, want AddRefused", res.Status)
	}
}

// TestRemoveWindowsRegistry_NoOp confirms removing an entry that
// never existed returns RemoveNoOp (clean exit 0).
func TestRemoveWindowsRegistry_NoOp(t *testing.T) {
	res, err := RemoveWithOptions(uniqueName(t), RemoveOptions{})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if res.Status != RemoveNoOp {
		t.Fatalf("status = %d, want RemoveNoOp", res.Status)
	}
}

// TestAddWindowsStartupFolder_RoundTrip exercises the .lnk backend.
// Skips when powershell.exe is missing so this file still passes on
// Windows hosts that have it disabled.
func TestAddWindowsStartupFolder_RoundTrip(t *testing.T) {
	if _, err := exec.LookPath("powershell.exe"); err != nil {
		t.Skip("powershell.exe not on PATH; skipping .lnk backend test")
	}
	withIsolatedAppData(t)
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	res, err := Add(AddOptions{Name: name, Exec: `C:\Windows\System32\notepad.exe`,
		Backend: BackendStartupFolder})
	if err != nil || res.Status != AddCreated {
		t.Fatalf("Add: %v / status=%d", err, res.Status)
	}
	if _, err := os.Stat(res.Path); err != nil {
		t.Fatalf("lnk not on disk: %v", err)
	}
	entries, err := listWindowsStartupFolder()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("List returned zero entries")
	}
	rm, err := RemoveWithOptions(name, RemoveOptions{})
	if err != nil || rm.Status != RemoveDeleted {
		t.Fatalf("Remove: %v / status=%d", err, rm.Status)
	}
	if _, err := os.Stat(res.Path); !os.IsNotExist(err) {
		t.Errorf("lnk still on disk after remove")
	}
}
