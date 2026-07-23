package completion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestDefaultPowerShellProfilePathsWindows(t *testing.T) {
	home := filepath.Join(os.TempDir(), "alim")
	paths := defaultPowerShellProfilePaths(home, "windows")
	expected := []string{
		filepath.Join(home, "Documents", "PowerShell", "profile.ps1"),
		filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"),
		filepath.Join(home, "Documents", "WindowsPowerShell", "profile.ps1"),
		filepath.Join(home, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"),
	}

	if len(paths) != len(expected) {
		t.Fatalf("expected %d paths, got %d", len(expected), len(paths))
	}
	for i, want := range expected {
		if paths[i] != want {
			t.Fatalf("path %d mismatch: want %s, got %s", i, want, paths[i])
		}
	}
}

func TestAddSourceLineCreatesProfileDir(t *testing.T) {
	scriptPath := filepath.Join(t.TempDir(), "gitmap", constants.CompFilePS)
	profilePath := filepath.Join(t.TempDir(), "Documents", "WindowsPowerShell", "profile.ps1")

	err := addSourceLine(scriptPath, profilePath, constants.ShellPowerShell)
	if err != nil {
		t.Fatalf("addSourceLine failed: %v", err)
	}

	data, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}

	if !strings.Contains(string(data), buildSourceLine(scriptPath, constants.ShellPowerShell)) {
		t.Fatal("expected PowerShell source line in created profile")
	}
}

func TestUniqueProfilePathsDropsDuplicatesAndEmptyValues(t *testing.T) {
	paths := uniqueProfilePaths([]string{"", "a", "a", " b ", "b", "c"})
	expected := []string{"a", "b", "c"}

	if len(paths) != len(expected) {
		t.Fatalf("expected %d unique paths, got %d", len(expected), len(paths))
	}
	for i, want := range expected {
		if paths[i] != want {
			t.Fatalf("path %d mismatch: want %s, got %s", i, want, paths[i])
		}
	}
}
