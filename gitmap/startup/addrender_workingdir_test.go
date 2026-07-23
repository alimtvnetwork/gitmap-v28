package startup

// Focused renderer tests for the --working-dir option. We assert
// each per-OS body shape from the pure renderer (no filesystem)
// so the test runs identically on every platform: renderDesktop
// is OS-agnostic, and renderPlist is callable directly.

import (
	"strings"
	"testing"
)

// TestRenderDesktop_EmitsPathWhenWorkingDirSet confirms a non-empty
// WorkingDir produces an XDG-spec `Path=<dir>` line and that the
// line is placed before Terminal= so `desktop-file-validate`'s
// recommended field order is preserved.
func TestRenderDesktop_EmitsPathWhenWorkingDirSet(t *testing.T) {
	body := string(renderDesktop("watch", AddOptions{
		Exec: "/usr/local/bin/gitmap watch", WorkingDir: "/srv/work",
	}))
	if !strings.Contains(body, "\nPath=/srv/work\n") {
		t.Errorf("missing Path= line:\n%s", body)
	}
	pathIdx := strings.Index(body, "Path=")
	termIdx := strings.Index(body, "Terminal=")
	if pathIdx < 0 || termIdx < 0 || pathIdx > termIdx {
		t.Errorf("Path= must precede Terminal= for desktop-file-validate order:\n%s", body)
	}
}

// TestRenderDesktop_OmitsPathWhenWorkingDirEmpty confirms the
// default (no flag) produces NO Path= line — keeps the file shape
// identical to pre-feature behavior for callers that don't opt in.
func TestRenderDesktop_OmitsPathWhenWorkingDirEmpty(t *testing.T) {
	body := string(renderDesktop("watch", AddOptions{Exec: "/x"}))
	if strings.Contains(body, "Path=") {
		t.Errorf("unexpected Path= line in default body:\n%s", body)
	}
}

// TestRenderPlist_EmitsWorkingDirectoryKey confirms the macOS
// LaunchAgent renderer emits the canonical <key>WorkingDirectory
// </key><string>...</string> pair right after RunAtLoad.
func TestRenderPlist_EmitsWorkingDirectoryKey(t *testing.T) {
	body := string(renderPlist("watch", AddOptions{
		Exec: "/usr/local/bin/gitmap watch", WorkingDir: "/srv/work",
	}))
	if !strings.Contains(body, "<key>WorkingDirectory</key>") {
		t.Errorf("missing WorkingDirectory key:\n%s", body)
	}
	if !strings.Contains(body, "<string>/srv/work</string>") {
		t.Errorf("missing WorkingDirectory value:\n%s", body)
	}
	runAt := strings.Index(body, "<key>RunAtLoad</key>")
	wd := strings.Index(body, "<key>WorkingDirectory</key>")
	if runAt < 0 || wd < 0 || wd < runAt {
		t.Errorf("WorkingDirectory must follow RunAtLoad:\n%s", body)
	}
}

// TestRenderPlist_OmitsWorkingDirectoryWhenEmpty confirms the key
// is absent for callers that did not pass --working-dir, so the
// LaunchAgent file stays minimal in the common case.
func TestRenderPlist_OmitsWorkingDirectoryWhenEmpty(t *testing.T) {
	body := string(renderPlist("watch", AddOptions{Exec: "/x"}))
	if strings.Contains(body, "WorkingDirectory") {
		t.Errorf("unexpected WorkingDirectory in default body:\n%s", body)
	}
}
