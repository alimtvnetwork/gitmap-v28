package cmdchromeprofile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChromeUserDataDirOverride(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)

	got := chromeUserDataDir()

	if got != tempDir {
		t.Fatalf("expected override %q, got %q", tempDir, got)
	}
}

func TestResolveWindowsChromeUserDataDirCandidate(t *testing.T) {
	tempLocal := t.TempDir()
	expectedDir := filepath.Join(tempLocal, "Google", "Chrome", "User Data")
	_ = os.MkdirAll(expectedDir, 0o755)

	t.Setenv("LOCALAPPDATA", tempLocal)
	t.Setenv("GITMAP_CHROME_USER_DATA", "")

	got := resolveWindowsChromeUserDataDir(t.TempDir())

	if got != expectedDir {
		t.Fatalf("expected candidate %q, got %q", expectedDir, got)
	}
}

func TestResolveWindowsChromeUserDataDirDefaultWhenEmpty(t *testing.T) {
	tempHome := t.TempDir()
	expectedDir := filepath.Join(tempHome, "AppData", "Local", "Google", "Chrome", "User Data")

	t.Setenv("LOCALAPPDATA", "")
	t.Setenv("USERNAME", "")
	t.Setenv("GITMAP_CHROME_USER_DATA", "")

	got := resolveWindowsChromeUserDataDir(tempHome)

	if got != expectedDir {
		t.Fatalf("expected default %q, got %q", expectedDir, got)
	}
}

func TestResolveWindowsSystemDriveFallback(t *testing.T) {
	t.Setenv("SystemDrive", "")
	got := resolveWindowsSystemDrive()
	expected := "C:" + string(filepath.Separator)

	if got != expected {
		t.Fatalf("expected default %q, got %q", expected, got)
	}
}

func TestResolveWindowsSystemDriveExisting(t *testing.T) {
	t.Setenv("SystemDrive", "D:")
	got := resolveWindowsSystemDrive()
	expected := "D:" + string(filepath.Separator)

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestFallbackHostChromeUserDataDirEmptyUser(t *testing.T) {
	t.Setenv("USERNAME", "")
	got := fallbackHostChromeUserDataDir()

	if got != "" {
		t.Fatalf("expected empty fallback for empty username, got %q", got)
	}
}
