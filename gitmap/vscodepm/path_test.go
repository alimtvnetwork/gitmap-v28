package vscodepm

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestWindowsCandidate_AppDataPresent asserts that APPDATA wins on
// Windows when set.
func TestWindowsCandidate_AppDataPresent(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only test")
	}
	t.Setenv("APPDATA", filepath.FromSlash("C:/Users/jane/AppData/Roaming"))

	got := windowsUserDataCandidate()
	want := filepath.Join(filepath.FromSlash("C:/Users/jane/AppData/Roaming"), "Code")
	if got != want {
		t.Errorf("windowsUserDataCandidate() = %q, want %q", got, want)
	}
}

// TestWindowsCandidate_NoEnvReturnsEmpty asserts the helper returns ""
// when neither APPDATA nor USERPROFILE is set.
func TestWindowsCandidate_NoEnvReturnsEmpty(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only test")
	}
	t.Setenv("APPDATA", "")
	t.Setenv("USERPROFILE", "")

	if got := windowsUserDataCandidate(); got != "" {
		t.Errorf("windowsUserDataCandidate() = %q, want empty string", got)
	}
}

// TestLinuxCandidate_XDGPrimary asserts $XDG_CONFIG_HOME wins.
func TestLinuxCandidate_XDGPrimary(t *testing.T) {
	xdg := filepath.FromSlash("/custom/xdg")
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Setenv("HOME", filepath.FromSlash("/home/jane"))

	got := linuxUserDataCandidate()
	want := filepath.Join(xdg, "Code")
	if got != want {
		t.Errorf("linuxUserDataCandidate() = %q, want %q", got, want)
	}
}

// TestLinuxCandidate_HomeFallback asserts $HOME/.config/Code fallback.
func TestLinuxCandidate_HomeFallback(t *testing.T) {
	home := filepath.FromSlash("/home/jane")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)

	got := linuxUserDataCandidate()
	want := filepath.Join(home, filepath.FromSlash(".config/Code"))
	if got != want {
		t.Errorf("linuxUserDataCandidate() = %q, want %q", got, want)
	}
}

// TestLinuxCandidate_NoEnvReturnsEmpty asserts the helper returns ""
// when neither env var is set.
func TestLinuxCandidate_NoEnvReturnsEmpty(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")

	if got := linuxUserDataCandidate(); got != "" {
		t.Errorf("linuxUserDataCandidate() = %q, want empty string", got)
	}
}

// TestDarwinCandidate_HomePresent asserts the macOS Library path is built
// off $HOME.
func TestDarwinCandidate_HomePresent(t *testing.T) {
	home := filepath.FromSlash("/Users/jane")
	t.Setenv("HOME", home)

	got := darwinUserDataCandidate()
	want := filepath.Join(home, filepath.FromSlash("Library/Application Support/Code"))
	if got != want {
		t.Errorf("darwinUserDataCandidate() = %q, want %q", got, want)
	}
}

// TestDarwinCandidate_NoHomeReturnsEmpty asserts empty $HOME yields "".
func TestDarwinCandidate_NoHomeReturnsEmpty(t *testing.T) {
	t.Setenv("HOME", "")

	if got := darwinUserDataCandidate(); got != "" {
		t.Errorf("darwinUserDataCandidate() = %q, want empty string", got)
	}
}

// TestUserDataRoot_MissingEnvReturnsSentinel asserts the public resolver
// returns ErrUserDataMissing (not a string error) when no env is set, so
// callers can use errors.Is for soft-fail handling.
func TestUserDataRoot_MissingEnvReturnsSentinel(t *testing.T) {
	clearAllVSCodeEnv(t)

	_, err := UserDataRoot()
	if !errors.Is(err, ErrUserDataMissing) {
		t.Fatalf("UserDataRoot() err = %v, want %v", err, ErrUserDataMissing)
	}
}

// TestProjectsJSONPath_MissingRootReturnsSentinel asserts the higher-level
// path builder also surfaces ErrUserDataMissing when the root cannot be
// located, without ever attempting to touch the filesystem.
func TestProjectsJSONPath_MissingRootReturnsSentinel(t *testing.T) {
	clearAllVSCodeEnv(t)

	_, err := ProjectsJSONPath()
	if !errors.Is(err, ErrUserDataMissing) {
		t.Fatalf("ProjectsJSONPath() err = %v, want %v", err, ErrUserDataMissing)
	}
}

// TestProjectsJSONPath_RootExistsExtMissingReturnsSentinel points the
// resolver at a real temp dir (so UserDataRoot succeeds) but does NOT
// create the alefragnani.project-manager subtree, asserting the second
// sentinel fires.
func TestProjectsJSONPath_RootExistsExtMissingReturnsSentinel(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("darwin user-data path is too deep to fake cleanly in a temp dir")
	}

	parent := t.TempDir()
	codeDir := filepath.Join(parent, "Code")
	if err := os.MkdirAll(codeDir, 0o755); err != nil {
		t.Fatalf("mkdir Code: %v", err)
	}

	clearAllVSCodeEnv(t)
	if runtime.GOOS == "windows" {
		t.Setenv("APPDATA", parent)
	} else {
		t.Setenv("XDG_CONFIG_HOME", parent)
	}

	_, err := ProjectsJSONPath()
	if !errors.Is(err, ErrExtensionMissing) {
		t.Fatalf("ProjectsJSONPath() err = %v, want %v", err, ErrExtensionMissing)
	}
}

// clearAllVSCodeEnv unsets every env var the resolver consults, on every
// supported platform. Uses t.Setenv so values restore at test end.
func clearAllVSCodeEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"APPDATA", "USERPROFILE", "HOME", "XDG_CONFIG_HOME"} {
		t.Setenv(k, "")
	}
}
