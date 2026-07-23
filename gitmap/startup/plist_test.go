package startup

// macOS LaunchAgent tests for List / Remove. Mirrors the structure
// of startup_test.go (the Linux .desktop tests) so the two suites
// can be diffed line-for-line and either platform's coverage gaps
// are immediately visible. All tests skip on non-darwin so a Linux
// CI run still goes green.

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// writePlist drops a LaunchAgent fixture into the dir. `managed`
// controls whether the XGitmapManaged marker is emitted; `argv`
// becomes the ProgramArguments array (empty argv → no Program /
// ProgramArguments key, exercising the "(no Exec line)" path).
func writePlist(t *testing.T, dir, name string, managed bool, argv ...string) string {
	t.Helper()
	body := plistBody(name, managed, argv)
	full := filepath.Join(dir, name)
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}

	return full
}

// plistBody renders a minimal LaunchAgent plist. Kept inline (no
// template engine) so the fixture is self-evident at the test call
// site and the marker placement is obvious to reviewers.
func plistBody(name string, managed bool, argv []string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<plist version="1.0"><dict>` + "\n")
	b.WriteString("<key>Label</key><string>" + name + "</string>\n")
	if managed {
		b.WriteString("<key>" + constants.StartupPlistMarker + "</key><true/>\n")
	}
	if len(argv) > 0 {
		b.WriteString("<key>ProgramArguments</key><array>\n")
		for _, a := range argv {
			b.WriteString("<string>" + a + "</string>\n")
		}
		b.WriteString("</array>\n")
	}
	b.WriteString("</dict></plist>\n")

	return b.String()
}

// TestList_Plist_OnlyReturnsManaged is the macOS analog of
// TestList_OnlyReturnsManaged: mixed dir contents must produce a
// list filtered down to entries with BOTH the gitmap. prefix AND the
// XGitmapManaged marker.
func TestList_Plist_OnlyReturnsManaged(t *testing.T) {
	dir := withFakeLaunchAgentsDir(t)
	writePlist(t, dir, "gitmap.foo.plist", true, "/usr/local/bin/gitmap", "watch")
	writePlist(t, dir, "gitmap.bar.plist", true, "/usr/local/bin/bar")
	writePlist(t, dir, "gitmap.spoof.plist", false, "/evil")    // prefix only, no marker
	writePlist(t, dir, "com.apple.something.plist", true, "/x") // marker but wrong prefix

	got, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 managed entries, got %d (%+v)", len(got), got)
	}
	names := map[string]bool{got[0].Name: true, got[1].Name: true}
	if !names["gitmap.foo"] || !names["gitmap.bar"] {
		t.Fatalf("unexpected entries: %+v", got)
	}
}

// TestList_Plist_ExecJoinsProgramArguments confirms the renderer-
// facing Exec field is the space-joined ProgramArguments array.
// Without this, the table view would show a useless "(no Exec line)"
// for every macOS entry whose plist uses the array form (the norm).
func TestList_Plist_ExecJoinsProgramArguments(t *testing.T) {
	dir := withFakeLaunchAgentsDir(t)
	writePlist(t, dir, "gitmap.watcher.plist", true, "/usr/local/bin/gitmap", "watch", "~/projects")

	got, err := List()
	if err != nil || len(got) != 1 {
		t.Fatalf("List: err=%v entries=%+v", err, got)
	}
	want := "/usr/local/bin/gitmap watch ~/projects"
	if got[0].Exec != want {
		t.Fatalf("Exec = %q, want %q", got[0].Exec, want)
	}
}

// TestList_Plist_MissingDirReturnsEmpty confirms a fresh account
// (no ~/Library/LaunchAgents at all) produces an empty list, not an
// error — same idempotency guarantee as Linux.
func TestList_Plist_MissingDirReturnsEmpty(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("plist tests are macOS-only")
	}
	root := t.TempDir()
	t.Setenv("HOME", root)
	got, err := List()
	if err != nil {
		t.Fatalf("List on missing dir must not error, got: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

// TestRemove_Plist_StatusMatrix exercises every RemoveStatus branch
// against macOS fixtures. Same four rows as the Linux version so a
// regression on either platform is caught by mirror coverage.
func TestRemove_Plist_StatusMatrix(t *testing.T) {
	dir := withFakeLaunchAgentsDir(t)
	managedPath := writePlist(t, dir, "gitmap.keep.plist", true, "/x")
	writePlist(t, dir, "gitmap.thirdparty.plist", false, "/y")

	res, err := Remove("gitmap.keep")
	if err != nil || res.Status != RemoveDeleted {
		t.Fatalf("delete: status=%v err=%v", res.Status, err)
	}
	if _, statErr := os.Stat(managedPath); !os.IsNotExist(statErr) {
		t.Fatalf("file should be gone, stat err=%v", statErr)
	}

	res, _ = Remove("does-not-exist")
	if res.Status != RemoveNoOp {
		t.Fatalf("noop: got status=%v", res.Status)
	}

	res, _ = Remove("gitmap.thirdparty")
	if res.Status != RemoveRefused {
		t.Fatalf("refused: got status=%v", res.Status)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "gitmap.thirdparty.plist")); statErr != nil {
		t.Fatalf("refused file must remain on disk, stat err=%v", statErr)
	}

	res, _ = Remove("../../etc/passwd")
	if res.Status != RemoveBadName {
		t.Fatalf("badname: got status=%v", res.Status)
	}
}

// TestRemove_Plist_DotPlistSuffixTolerated mirrors the Linux
// .desktop-suffix-tolerated test: pasting the full filename works
// the same as the bare name.
func TestRemove_Plist_DotPlistSuffixTolerated(t *testing.T) {
	dir := withFakeLaunchAgentsDir(t)
	writePlist(t, dir, "gitmap.suffixed.plist", true, "/x")

	res, err := Remove("gitmap.suffixed.plist")
	if err != nil || res.Status != RemoveDeleted {
		t.Fatalf("got status=%v err=%v", res.Status, err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "gitmap.suffixed.plist")); !os.IsNotExist(statErr) {
		t.Fatalf("file should be gone")
	}
}
