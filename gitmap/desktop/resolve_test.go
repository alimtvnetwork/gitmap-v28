package desktop

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestResolveCLI_FindsKnownInstall is the end-to-end guard for the silent
// `gitmap desktop-sync` failure: when `github` is not on PATH but the
// Desktop installer's per-user shim exists, ResolveCLI must still find it.
func TestResolveCLI_FindsKnownInstall(t *testing.T) {
	if runtime.GOOS != constants.OSWindows {
		t.Skip("known-install fallback test targets Windows layout")
	}
	tmp := t.TempDir()
	binDir := filepath.Join(tmp, "GitHubDesktop", "bin")
	mkdirErr := os.MkdirAll(binDir, 0o755)
	if mkdirErr != nil {
		t.Fatalf("mkdir: %v", mkdirErr)
	}
	shim := filepath.Join(binDir, "github.bat")
	writeErr := os.WriteFile(shim, []byte("@echo off\r\n"), 0o644)
	if writeErr != nil {
		t.Fatalf("write shim: %v", writeErr)
	}
	t.Setenv("LOCALAPPDATA", tmp)
	t.Setenv("PATH", "")

	got := ResolveCLI()
	if got != shim {
		t.Fatalf("ResolveCLI = %q, want %q", got, shim)
	}
}

// TestResolveCLI_MissingReturnsEmpty asserts the resolver does not panic
// or fabricate a path when GitHub Desktop is not installed anywhere.
func TestResolveCLI_MissingReturnsEmpty(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("LOCALAPPDATA", tmp)
	t.Setenv("HOME", tmp)
	t.Setenv("PATH", tmp)

	got := ResolveCLI()
	if got != "" {
		t.Fatalf("ResolveCLI = %q, want empty", got)
	}
}

// TestCollectAppDirs verifies the Squirrel `app-*` filter, which underpins
// newest-version fallback when Desktop is installed but the top-level
// `bin\github.bat` shim has not been rewritten yet by the installer.
func TestCollectAppDirs(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"app-3.4.1", "app-3.4.10", "bin", "other"} {
		mkErr := os.Mkdir(filepath.Join(tmp, name), 0o755)
		if mkErr != nil {
			t.Fatalf("mkdir %s: %v", name, mkErr)
		}
	}
	entries, readErr := os.ReadDir(tmp)
	if readErr != nil {
		t.Fatalf("readdir: %v", readErr)
	}
	got := collectAppDirs(entries)
	if len(got) != 2 {
		t.Fatalf("collectAppDirs got %v, want 2 entries", got)
	}
}
