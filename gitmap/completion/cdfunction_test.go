package completion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestAppendCDFunctionWritesManagedWrappers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile")

	err := appendCDFunction(constants.CDFuncBash, path)
	if err != nil {
		t.Fatalf("appendCDFunction failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}

	text := string(data)
	if !strings.Contains(text, constants.CDFuncMarker) {
		t.Fatal("expected managed wrapper marker to be written")
	}
	if !strings.Contains(text, "gitmap() {") {
		t.Fatal("expected gitmap shell wrapper to be written")
	}
	if !strings.Contains(text, "gcd() {") {
		t.Fatal("expected gcd shell wrapper to be written")
	}
}

func TestAppendCDFunctionSkipsManagedWrapperWhenPresent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile")
	block := "\n" + constants.CDFuncMarker + "\n" + constants.CDFuncBash + "\n"

	err := os.WriteFile(path, []byte(block), 0o644)
	if err != nil {
		t.Fatalf("seed profile failed: %v", err)
	}

	err = appendCDFunction(constants.CDFuncBash, path)
	if err != nil {
		t.Fatalf("appendCDFunction failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}

	if string(data) != block {
		t.Fatal("expected managed wrapper block to remain unchanged")
	}
}

func TestAppendCDFunctionDoesNotSkipPathSnippetMarker(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile")
	pathSnippet := "\n" + constants.CDFuncMarkerLegacy + " - managed by installer. Do not edit manually.\n"

	if err := os.WriteFile(path, []byte(pathSnippet), 0o644); err != nil {
		t.Fatalf("seed profile failed: %v", err)
	}
	if err := appendCDFunction(constants.CDFuncPowerShell, path); err != nil {
		t.Fatalf("appendCDFunction failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}
	if countCDFunctionStartMarkers(string(data)) != 1 {
		t.Fatal("expected command wrapper after legacy PATH snippet marker")
	}
}

func TestAppendCDFunctionAppendsManagedWrapperAfterLegacyMarker(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile")
	legacy := "\n# gitmap cd wrapper\ngcd() {\n  cd \"$(gitmap cd \"$@\")\"\n}\n"

	err := os.WriteFile(path, []byte(legacy), 0o644)
	if err != nil {
		t.Fatalf("seed profile failed: %v", err)
	}

	err = appendCDFunction(constants.CDFuncBash, path)
	if err != nil {
		t.Fatalf("appendCDFunction failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}

	text := string(data)
	if !strings.Contains(text, legacy) {
		t.Fatal("expected legacy wrapper to remain for migration safety")
	}
	if !strings.Contains(text, constants.CDFuncMarker) {
		t.Fatal("expected managed wrapper marker to be appended")
	}
	if countCDFunctionStartMarkers(text) != 1 {
		t.Fatal("expected exactly one managed wrapper marker")
	}
}

func TestAppendCDFunctionRewritesLegacyEndMarker(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile")
	block := constants.CDFuncMarker + "\nfunction gitmap { 'stale' }\n" +
		constants.CDFuncMarkerEndLegacy + "\n"

	if err := os.WriteFile(path, []byte(block), 0o644); err != nil {
		t.Fatalf("seed profile failed: %v", err)
	}
	if err := appendCDFunction(constants.CDFuncPowerShell, path); err != nil {
		t.Fatalf("appendCDFunction failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}
	text := string(data)
	if strings.Contains(text, "'stale'") {
		t.Fatal("expected stale wrapper to be replaced")
	}
	if !strings.Contains(text, constants.CDFuncMarkerEnd) {
		t.Fatal("expected canonical installer-compatible end marker")
	}
}

func TestAppendCDFunctionMovesManagedWrapperAfterStaleSnippet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile")
	seed := constants.CDFuncMarker + "\nfunction gitmap { 'old' }\n" +
		constants.CDFuncMarkerEnd + "\n# later stale gitmap shell wrapper v2\nfunction gitmap { 'stale' }\n"

	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatalf("seed profile failed: %v", err)
	}
	if err := appendCDFunction(constants.CDFuncPowerShell, path); err != nil {
		t.Fatalf("appendCDFunction failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}
	text := string(data)
	if strings.LastIndex(text, constants.CDFuncMarker) < strings.LastIndex(text, "function gitmap { 'stale' }") {
		t.Fatal("expected current command wrapper to load after stale wrapper")
	}
}

func TestAppendCDFunctionCreatesProfileDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Documents", "WindowsPowerShell", "profile.ps1")

	err := appendCDFunction(constants.CDFuncPowerShell, path)
	if err != nil {
		t.Fatalf("appendCDFunction failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read profile failed: %v", err)
	}

	if !strings.Contains(string(data), constants.CDFuncMarker) {
		t.Fatal("expected managed wrapper marker in created profile")
	}
}

func TestAppendCDFunctionsWritesToMultipleProfiles(t *testing.T) {
	base := t.TempDir()
	paths := []string{
		filepath.Join(base, "Documents", "PowerShell", "profile.ps1"),
		filepath.Join(base, "Documents", "WindowsPowerShell", "profile.ps1"),
	}

	err := appendCDFunctions(constants.CDFuncPowerShell, paths)
	if err != nil {
		t.Fatalf("appendCDFunctions failed: %v", err)
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read profile failed for %s: %v", path, err)
		}
		if !strings.Contains(string(data), constants.CDFuncMarker) {
			t.Fatalf("expected managed wrapper marker in %s", path)
		}
	}
}

func TestRenderPowerShellCommandShimPinsInstalledExe(t *testing.T) {
	got := renderPowerShellCommandShim(`C:\Tools\git'map`)
	wants := []string{
		`Join-Path -Path 'C:\Tools\git''map' -ChildPath 'gitmap.exe'`,
		"Set-Location -LiteralPath ([string]$dest)",
		constants.EnvGitmapCommandWrapper,
		constants.EnvGitmapHandoffFile,
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Fatalf("shim missing %q\n%s", want, got)
		}
	}
}

// countCDFunctionStartMarkers counts lines that exactly start with the
// managed start marker, excluding the end marker (which has the start
// marker as a prefix and would otherwise inflate the count).
func countCDFunctionStartMarkers(text string) int {
	count := 0
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == constants.CDFuncMarker {
			count++
		}
	}
	return count
}
