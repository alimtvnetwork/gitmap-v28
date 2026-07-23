//go:build windows

package startup

// Regression tests that lock in the direct-value Run-key model:
// after `gitmap startup-add`, the HKCU Run key MUST contain
// exactly one new value (`gitmap-<name>`) and NO sibling
// `.gitmap-managed` companion. A future refactor that re-introduces
// the sibling marker would silently break Task Manager's Startup
// tab and execute the literal string "true" at every login.

import (
	"strings"
	"testing"

	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestAddWindowsRegistry_NoSiblingMarker is the contract test for
// the direct-value model. Snapshots Run-key value names before and
// after Add; the diff MUST be exactly {gitmap-<name>}.
func TestAddWindowsRegistry_NoSiblingMarker(t *testing.T) {
	name := uniqueName(t)
	t.Cleanup(func() { cleanupRegistryEntry(t, name) })

	before := snapshotRunValueNames(t)
	res, err := Add(AddOptions{Name: name, Exec: `C:\gitmap.exe watch`,
		Backend: BackendRegistry})
	if err != nil || res.Status != AddCreated {
		t.Fatalf("Add: %v / status=%d", err, res.Status)
	}
	added := diffRunValueNames(snapshotRunValueNames(t), before)
	want := constants.StartupWinValuePrefix + name
	if len(added) != 1 || added[0] != want {
		t.Fatalf("Run-key new values = %v, want exactly [%s]", added, want)
	}
	for _, n := range added {
		if strings.HasSuffix(n, ".gitmap-managed") {
			t.Fatalf("sibling marker leaked into Run key: %s", n)
		}
	}
}

// snapshotRunValueNames reads every value name under HKCU\Run.
// Returns nil if the key does not exist (fresh test environment).
func snapshotRunValueNames(t *testing.T) []string {
	t.Helper()
	k, err := registry.OpenKey(registry.CURRENT_USER,
		constants.RegRunKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		t.Fatalf("open Run key: %v", err)
	}
	defer k.Close()
	names, err := k.ReadValueNames(-1)
	if err != nil {
		t.Fatalf("read value names: %v", err)
	}
	return names
}

// diffRunValueNames returns the set difference (after - before).
// O(n*m) is fine — Run keys typically hold <30 values.
func diffRunValueNames(after, before []string) []string {
	seen := make(map[string]struct{}, len(before))
	for _, n := range before {
		seen[n] = struct{}{}
	}
	var out []string
	for _, n := range after {
		if _, ok := seen[n]; !ok {
			out = append(out, n)
		}
	}
	return out
}
