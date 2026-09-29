package completion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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

func TestAddSourceLinePowerShellIdempotent(t *testing.T) {
	scriptPath := filepath.Join(t.TempDir(), "gitmap", constants.CompFilePS)
	profilePath := filepath.Join(t.TempDir(), "Documents", "PowerShell", "profile.ps1")

	if err := addSourceLine(scriptPath, profilePath, constants.ShellPowerShell); err != nil {
		t.Fatalf("first addSourceLine failed: %v", err)
	}

	if err := addSourceLine(scriptPath, profilePath, constants.ShellPowerShell); err != nil {
		t.Fatalf("second addSourceLine failed: %v", err)
	}

	data, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}

	occurrences := strings.Count(string(data), compMarkerStart)
	if occurrences != 1 {
		t.Fatalf("expected exactly 1 marker block, got %d in:\n%s", occurrences, string(data))
	}
}

func TestAddSourceLinePowerShellStripsLegacyBlocks(t *testing.T) {
	scriptPath := filepath.Join(t.TempDir(), "gitmap", constants.CompFilePS)
	profilePath := filepath.Join(t.TempDir(), "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")

	initial := `# custom user tool
function mytool { Write-Host "hi" }

# gitmap shell completion
. 'C:\Users\Old\AppData\Roaming\gitmap\completions.ps1'
if ((Get-Module -ListAvailable -Name PSReadLine) -and -not [Console]::IsOutputRedirected) {
    try {
        Set-PSReadLineOption -PredictionSource History -ErrorAction SilentlyContinue
        Set-PSReadLineOption -PredictionViewStyle ListView -ErrorAction SilentlyContinue
    } catch {}
}

$env:CUSTOM_VAR = "1"
`
	if err := os.MkdirAll(filepath.Dir(profilePath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(profilePath, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial profile failed: %v", err)
	}

	if err := addSourceLine(scriptPath, profilePath, constants.ShellPowerShell); err != nil {
		t.Fatalf("addSourceLine failed: %v", err)
	}

	data, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "Old\\AppData") {
		t.Fatalf("legacy completions.ps1 path not removed:\n%s", content)
	}
	if !strings.Contains(content, "function mytool") {
		t.Fatalf("user custom function was removed:\n%s", content)
	}
	if !strings.Contains(content, "$env:CUSTOM_VAR") {
		t.Fatalf("user custom var was removed:\n%s", content)
	}
	if strings.Count(content, compMarkerStart) != 1 {
		t.Fatalf("expected exactly 1 modern marker, got %d", strings.Count(content, compMarkerStart))
	}
}
