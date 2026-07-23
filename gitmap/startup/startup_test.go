package startup

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// withFakeAutostartDir points $XDG_CONFIG_HOME at a temp dir so the
// Linux-shape List / Remove tests below operate against an isolated
// .desktop tree. Skipped on Windows (unsupported) and macOS (whose
// LaunchAgent path is exercised by plist_test.go instead — those
// tests use a different fixture format and dir layout).
func withFakeAutostartDir(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("startup package does not support Windows")
	}
	if runtime.GOOS == "darwin" {
		t.Skip(".desktop tests are Linux-only; plist_test.go covers macOS")
	}
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	dir := filepath.Join(root, "autostart")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	return dir
}

// writeDesktop is a small helper to drop a .desktop fixture into the
// autostart dir. `managed` controls whether the marker line is
// emitted, letting tests build both gitmap-managed and third-party
// fixtures from the same call site.
func writeDesktop(t *testing.T, dir, name string, managed bool, exec string) string {
	t.Helper()
	body := "[Desktop Entry]\nType=Application\nName=" + name + "\nExec=" + exec + "\n"
	if managed {
		body += constants.StartupMarkerKey + "=" + constants.StartupMarkerVal + "\n"
	}
	full := filepath.Join(dir, name)
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}

	return full
}

// TestList_OnlyReturnsManaged covers the core safety guarantee: even
// when the autostart dir contains a mix of gitmap-managed and
// third-party files (and even when a third-party file has the
// gitmap- prefix in its name), List returns ONLY entries whose body
// carries the X-Gitmap-Managed=true marker.
func TestList_OnlyReturnsManaged(t *testing.T) {
	dir := withFakeAutostartDir(t)
	writeDesktop(t, dir, "gitmap-foo.desktop", true, "/usr/bin/foo")
	writeDesktop(t, dir, "gitmap-bar.desktop", true, "/usr/bin/bar")
	writeDesktop(t, dir, "gitmap-spoof.desktop", false, "/evil") // prefix only, no marker
	writeDesktop(t, dir, "thirdparty.desktop", true, "/x")       // marker but wrong prefix

	got, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 managed entries, got %d (%+v)", len(got), got)
	}
	names := map[string]bool{got[0].Name: true, got[1].Name: true}
	if !names["gitmap-foo"] || !names["gitmap-bar"] {
		t.Fatalf("unexpected entries: %+v", got)
	}
}

// TestList_MissingDirReturnsEmpty confirms that a fresh user account
// (no ~/.config/autostart at all) produces an empty list, NOT an
// error. Idempotent CLI behavior demands this. Linux-only because
// the env-var manipulation here targets XDG_CONFIG_HOME; the
// equivalent macOS case is covered in plist_test.go.
func TestList_MissingDirReturnsEmpty(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		t.Skip("Linux-only; macOS missing-dir behavior covered in plist_test.go")
	}
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	// Note: no autostart subdir created.
	got, err := List()
	if err != nil {
		t.Fatalf("List on missing dir must not error, got: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

// TestRemove_StatusMatrix exercises every RemoveStatus branch with
// dedicated fixtures so each row of the contract is covered:
//  1. RemoveDeleted    — managed file present, gets unlinked.
//  2. RemoveNoOp       — name doesn't exist, clean exit.
//  3. RemoveRefused    — file exists but lacks the marker.
//  4. RemoveBadName    — input contains a path separator.
func TestRemove_StatusMatrix(t *testing.T) {
	dir := withFakeAutostartDir(t)
	managedPath := writeDesktop(t, dir, "gitmap-keep.desktop", true, "/x")
	writeDesktop(t, dir, "gitmap-thirdparty.desktop", false, "/y")

	// 1. Deleted
	res, err := Remove("gitmap-keep")
	if err != nil || res.Status != RemoveDeleted {
		t.Fatalf("delete: status=%v err=%v", res.Status, err)
	}
	if _, statErr := os.Stat(managedPath); !os.IsNotExist(statErr) {
		t.Fatalf("file should be gone, stat err=%v", statErr)
	}

	// 2. NoOp
	res, _ = Remove("does-not-exist")
	if res.Status != RemoveNoOp {
		t.Fatalf("noop: got status=%v", res.Status)
	}

	// 3. Refused (file present but unmanaged)
	res, _ = Remove("gitmap-thirdparty")
	if res.Status != RemoveRefused {
		t.Fatalf("refused: got status=%v", res.Status)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "gitmap-thirdparty.desktop")); statErr != nil {
		t.Fatalf("refused file must remain on disk, stat err=%v", statErr)
	}

	// 4. BadName (path traversal attempt)
	res, _ = Remove("../../etc/passwd")
	if res.Status != RemoveBadName {
		t.Fatalf("badname: got status=%v", res.Status)
	}
}

// TestRemove_DotDesktopSuffixTolerated verifies that users who paste
// `gitmap-foo.desktop` (with the suffix) get the same result as
// users who type `gitmap-foo` — the normalizer strips it.
func TestRemove_DotDesktopSuffixTolerated(t *testing.T) {
	dir := withFakeAutostartDir(t)
	writeDesktop(t, dir, "gitmap-suffixed.desktop", true, "/x")

	res, err := Remove("gitmap-suffixed.desktop")
	if err != nil || res.Status != RemoveDeleted {
		t.Fatalf("got status=%v err=%v", res.Status, err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "gitmap-suffixed.desktop")); !os.IsNotExist(statErr) {
		t.Fatalf("file should be gone")
	}
}
