package startup

// withFakeAutostartDir + writeDesktop are defined in startup_test.go
// (same package); reused here without redeclaration.

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestAdd_CreatesManagedFile verifies the happy path: a fresh
// autostart dir gets a new gitmap-<name>.desktop file containing the
// X-Gitmap-Managed=true marker and the requested Exec line.
func TestAdd_CreatesManagedFile(t *testing.T) {
	dir := withFakeAutostartDir(t)
	res, err := Add(AddOptions{Name: "watch", Exec: "/usr/local/bin/gitmap watch"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddCreated {
		t.Fatalf("status = %d, want AddCreated", res.Status)
	}
	want := filepath.Join(dir, "gitmap-watch.desktop")
	if res.Path != want {
		t.Fatalf("path = %s, want %s", res.Path, want)
	}
	body, err := os.ReadFile(res.Path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	s := string(body)
	if !strings.Contains(s, constants.StartupMarkerKey+"="+constants.StartupMarkerVal) {
		t.Errorf("body missing managed marker:\n%s", s)
	}
	if !strings.Contains(s, "Exec=/usr/local/bin/gitmap watch") {
		t.Errorf("body missing Exec line:\n%s", s)
	}
}

// TestAdd_RefusesNonManagedOverwrite is the security-critical case:
// a third-party file with the same prefixed name must NOT be
// overwritten, even when --force is passed.
func TestAdd_RefusesNonManagedOverwrite(t *testing.T) {
	dir := withFakeAutostartDir(t)
	target := writeDesktop(t, dir, "gitmap-watch.desktop", false, "/evil")
	res, err := Add(AddOptions{Name: "watch", Exec: "/safe", Force: true})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddRefused {
		t.Fatalf("status = %d, want AddRefused", res.Status)
	}
	body, _ := os.ReadFile(target)
	if !strings.Contains(string(body), "/evil") {
		t.Errorf("third-party file body was modified:\n%s", body)
	}
}

// TestAdd_ExistsWithoutForce confirms idempotent re-runs.
func TestAdd_ExistsWithoutForce(t *testing.T) {
	withFakeAutostartDir(t)
	first, err := Add(AddOptions{Name: "watch", Exec: "/v1"})
	if err != nil || first.Status != AddCreated {
		t.Fatalf("first add: %v / %d", err, first.Status)
	}
	second, err := Add(AddOptions{Name: "watch", Exec: "/v2"})
	if err != nil {
		t.Fatalf("second add: %v", err)
	}
	if second.Status != AddExists {
		t.Fatalf("status = %d, want AddExists", second.Status)
	}
	body, _ := os.ReadFile(second.Path)
	if !strings.Contains(string(body), "Exec=/v1") {
		t.Errorf("file overwritten without --force:\n%s", body)
	}
}

// TestAdd_ForceOverwritesOurOwn confirms --force lifts the
// AddExists guard for entries gitmap itself created.
func TestAdd_ForceOverwritesOurOwn(t *testing.T) {
	withFakeAutostartDir(t)
	if _, err := Add(AddOptions{Name: "watch", Exec: "/v1"}); err != nil {
		t.Fatalf("first add: %v", err)
	}
	res, err := Add(AddOptions{Name: "watch", Exec: "/v2", Force: true})
	if err != nil {
		t.Fatalf("force add: %v", err)
	}
	if res.Status != AddOverwritten {
		t.Fatalf("status = %d, want AddOverwritten", res.Status)
	}
	body, _ := os.ReadFile(res.Path)
	if !strings.Contains(string(body), "Exec=/v2") {
		t.Errorf("force did not replace body:\n%s", body)
	}
}

// TestAdd_BadName covers the validation gate: empty + slash + NUL
// all return AddBadName without writing anything.
func TestAdd_BadName(t *testing.T) {
	withFakeAutostartDir(t)
	cases := []string{"", "../escape", "with/slash", "nul\x00name"}
	for _, name := range cases {
		res, err := Add(AddOptions{Name: name, Exec: "/x"})
		if err != nil {
			t.Fatalf("Add(%q): %v", name, err)
		}
		if res.Status != AddBadName {
			t.Errorf("name=%q: status = %d, want AddBadName", name, res.Status)
		}
	}
}

// TestAdd_AutoCreatesDir confirms a missing autostart dir is created
// rather than producing an error.
func TestAdd_AutoCreatesDir(t *testing.T) {
	// Windows uses the registry / .lnk Startup-folder backends, NOT
	// a file-based AutostartDir() — the EnvStartupAutostartDir
	// override is honored only on Linux/macOS where Add() routes
	// through AutostartDir(). Skipping here keeps the cross-platform
	// CI green without weakening the Linux/macOS coverage.
	if runtime.GOOS == "windows" {
		t.Skip("autostart dir is file-based on Linux/macOS only; Windows uses registry/.lnk backends")
	}
	// Use the cross-platform $GITMAP_AUTOSTART_DIR override so the
	// test exercises the same Add() code path on Linux and macOS
	// without per-OS skips. The override is the explicit, documented
	// test-injection point baked into AutostartDir().
	root := filepath.Join(t.TempDir(), "autostart")
	t.Setenv(constants.EnvStartupAutostartDir, root)
	res, err := Add(AddOptions{Name: "watch", Exec: "/x"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddCreated {
		t.Fatalf("status = %d, want AddCreated", res.Status)
	}
	if _, err := os.Stat(root); err != nil {
		t.Errorf("autostart dir not created: %v", err)
	}
}

// TestAdd_PrefixNotDoubled covers the cosmetic but important case
// where a caller passes a name that already starts with the gitmap-
// prefix; the on-disk filename must not become gitmap-gitmap-...
func TestAdd_PrefixNotDoubled(t *testing.T) {
	dir := withFakeAutostartDir(t)
	res, err := Add(AddOptions{Name: "gitmap-watch", Exec: "/x"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	want := filepath.Join(dir, "gitmap-watch.desktop")
	if res.Path != want {
		t.Fatalf("path = %s, want %s (no double prefix)", res.Path, want)
	}
}
