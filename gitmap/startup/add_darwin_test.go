package startup

// macOS LaunchAgent test suite for Add(). Mirrors add_test.go's
// Linux/.desktop coverage so darwin gets the same five behavioral
// guarantees: managed-marker writeback, third-party refusal,
// idempotent re-runs, --force overwrite, and bad-name rejection.
//
// All tests skip on non-darwin so a Linux CI run sees them as
// "skipped" rather than failing on a missing $HOME/Library path.

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// writeRawPlist drops a third-party (un-marked) plist in place so
// the refusal test has something to NOT overwrite. The body is the
// minimum launchd accepts; absence of the XGitmapManaged key is what
// flags it as not-ours.
func writeRawPlist(t *testing.T, dir, filename, program string) string {
	t.Helper()
	body := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>` + filename + `</string>
  <key>Program</key>
  <string>` + program + `</string>
</dict>
</plist>
`
	full := filepath.Join(dir, filename)
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}

	return full
}

// TestAddDarwin_CreatesManagedFile is the happy path: a fresh
// LaunchAgents dir gets a new gitmap.<name>.plist file containing
// the XGitmapManaged <true/> marker and the Exec split into
// ProgramArguments.
func TestAddDarwin_CreatesManagedFile(t *testing.T) {
	dir := withFakeLaunchAgentsDir(t)
	res, err := Add(AddOptions{Name: "watch", Exec: "/usr/local/bin/gitmap watch"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddCreated {
		t.Fatalf("status = %d, want AddCreated", res.Status)
	}
	want := filepath.Join(dir, "gitmap.watch.plist")
	if res.Path != want {
		t.Fatalf("path = %s, want %s", res.Path, want)
	}
	body, err := os.ReadFile(res.Path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	s := string(body)
	if !strings.Contains(s, "<key>"+constants.StartupPlistMarker+"</key>") {
		t.Errorf("body missing managed marker key:\n%s", s)
	}
	if !strings.Contains(s, "<true/>") {
		t.Errorf("body missing <true/> for marker:\n%s", s)
	}
	if !strings.Contains(s, "<string>/usr/local/bin/gitmap</string>") ||
		!strings.Contains(s, "<string>watch</string>") {
		t.Errorf("body missing split ProgramArguments:\n%s", s)
	}
	if !strings.Contains(s, "<key>Label</key>") ||
		!strings.Contains(s, "<string>gitmap.watch</string>") {
		t.Errorf("body missing Label=gitmap.watch:\n%s", s)
	}
}

// TestAddDarwin_RefusesNonManagedOverwrite is the security-critical
// case: a third-party plist with the same prefixed name must NOT be
// overwritten, even with --force.
func TestAddDarwin_RefusesNonManagedOverwrite(t *testing.T) {
	dir := withFakeLaunchAgentsDir(t)
	target := writeRawPlist(t, dir, "gitmap.watch.plist", "/evil")
	res, err := Add(AddOptions{Name: "watch", Exec: "/safe", Force: true})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddRefused {
		t.Fatalf("status = %d, want AddRefused", res.Status)
	}
	body, _ := os.ReadFile(target)
	if !strings.Contains(string(body), "/evil") {
		t.Errorf("third-party plist was modified:\n%s", body)
	}
}

// TestAddDarwin_ExistsWithoutForce confirms idempotent re-runs.
func TestAddDarwin_ExistsWithoutForce(t *testing.T) {
	withFakeLaunchAgentsDir(t)
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
	if !strings.Contains(string(body), "<string>/v1</string>") {
		t.Errorf("file overwritten without --force:\n%s", body)
	}
}

// TestAddDarwin_ForceOverwritesOurOwn confirms --force lifts the
// AddExists guard for entries gitmap itself created.
func TestAddDarwin_ForceOverwritesOurOwn(t *testing.T) {
	withFakeLaunchAgentsDir(t)
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
	if !strings.Contains(string(body), "<string>/v2</string>") {
		t.Errorf("force did not replace body:\n%s", body)
	}
}

// TestAddDarwin_PrefixNotDoubled covers the cosmetic-but-important
// case where a caller passes a name that already starts with the
// gitmap. prefix; the on-disk filename must not become
// gitmap.gitmap.<name>.plist.
func TestAddDarwin_PrefixNotDoubled(t *testing.T) {
	dir := withFakeLaunchAgentsDir(t)
	res, err := Add(AddOptions{Name: "gitmap.watch", Exec: "/x"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	want := filepath.Join(dir, "gitmap.watch.plist")
	if res.Path != want {
		t.Fatalf("path = %s, want %s (no double prefix)", res.Path, want)
	}
}

// TestAddDarwin_AutoCreatesDir confirms a missing LaunchAgents dir
// is created rather than producing an error — same idempotent-first-
// run guarantee as the Linux side gives autostart/.
func TestAddDarwin_AutoCreatesDir(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("plist tests are macOS-only")
	}
	root := t.TempDir()
	t.Setenv("HOME", root)
	res, err := Add(AddOptions{Name: "watch", Exec: "/x"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Status != AddCreated {
		t.Fatalf("status = %d, want AddCreated", res.Status)
	}
	if _, err := os.Stat(filepath.Join(root, "Library", "LaunchAgents")); err != nil {
		t.Errorf("LaunchAgents dir not created: %v", err)
	}
}

// TestAddDarwin_ListSeesAddedEntry is the round-trip integration
// test: write with Add(), then List() must surface it. This is the
// guarantee that the marker rendered by addplist.go matches the
// marker parsed by plist.go — a regression in either would break
// `gitmap startup-list` for darwin users immediately after add.
func TestAddDarwin_ListSeesAddedEntry(t *testing.T) {
	withFakeLaunchAgentsDir(t)
	if _, err := Add(AddOptions{Name: "watch", Exec: "/usr/local/bin/gitmap watch"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	entries, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("List returned %d entries, want 1: %#v", len(entries), entries)
	}
	if entries[0].Name != "gitmap.watch" {
		t.Errorf("entry name = %s, want gitmap.watch", entries[0].Name)
	}
	if entries[0].Exec != "/usr/local/bin/gitmap watch" {
		t.Errorf("entry exec = %s, want /usr/local/bin/gitmap watch", entries[0].Exec)
	}
}
