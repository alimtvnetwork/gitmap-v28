package cmd

import (
	"strings"
	"testing"
)

func TestParseTagFromGitLine(t *testing.T) {
	line := "67fcff9d23d4d88ce81f282dbebce7789450a65b\trefs/tags/v4.6.9"
	ver, ok := parseTagFromGitLine(line)
	if !ok {
		t.Fatalf("expected valid parse, got ok=false")
	}
	if ver.CoreString() != "4.6.9" {
		t.Errorf("expected 4.6.9, got %s", ver.CoreString())
	}
}

func TestParseTagFromGitLine_Invalid(t *testing.T) {
	line := "invalid_line_without_tag"
	_, ok := parseTagFromGitLine(line)
	if ok {
		t.Errorf("expected parse failure for invalid line")
	}
}

func TestFindHighestSemverInOutput(t *testing.T) {
	raw := `8e04b45c78e8d2ca0dadd97731251c93329489ff\trefs/tags/v3.3.18
c6ab5e6fa689271a6453ed9eefead7caa99285fc\trefs/tags/v3.3.31
82f61414e5a314c95b3be67a350f024960d47235\trefs/tags/v4.0.1
67fcff9d23d4d88ce81f282dbebce7789450a65b\trefs/tags/v4.6.9
1df0461084a795971dbc44fcc5895e73d1e9fb67\trefs/tags/v4.6.8`

	highest, isFound := findHighestSemverInOutput(raw)
	if !isFound {
		t.Fatalf("expected semver to be found")
	}
	if highest.CoreString() != "4.6.9" {
		t.Errorf("expected highest 4.6.9, got %s", highest.CoreString())
	}
}

func TestResolveAgManagerAssetURL_SpecifiedVersion(t *testing.T) {
	url, ver, err := resolveAgManagerAssetURL("4.6.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ver != "4.6.0" {
		t.Errorf("expected ver 4.6.0, got %s", ver)
	}
	if !strings.Contains(url, "v4.6.0") {
		t.Errorf("expected url to contain v4.6.0, got %s", url)
	}
}

func TestResolveAgManagerAssetURL_WithVPrefix(t *testing.T) {
	url, ver, err := resolveAgManagerAssetURL("v4.5.8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ver != "4.5.8" {
		t.Errorf("expected ver 4.5.8, got %s", ver)
	}
	if !strings.Contains(url, "v4.5.8") {
		t.Errorf("expected url to contain v4.5.8, got %s", url)
	}
}

func TestAgyInstallFlags(t *testing.T) {
	dryFlag := agyInstallCmd.Flags().Lookup("dry-run")
	if dryFlag == nil {
		t.Errorf("expected --dry-run flag on agyInstallCmd")
	}
	verFlag := agyInstallCmd.Flags().Lookup("version")
	if verFlag == nil {
		t.Errorf("expected --version flag on agyInstallCmd")
	}
	yesFlag := agyInstallCmd.Flags().Lookup("yes")
	if yesFlag == nil {
		t.Errorf("expected --yes flag on agyInstallCmd")
	}
}

func TestAgyInstallDryRunDispatch(t *testing.T) {
	opts := installOptions{DryRun: true}
	err := dispatchAgyInstallTarget("manager", opts)
	if err != nil {
		t.Errorf("expected nil error for dry-run dispatch, got %v", err)
	}
	errCli := dispatchAgyInstallTarget("cli", opts)
	if errCli != nil {
		t.Errorf("expected nil error for cli dry-run dispatch, got %v", errCli)
	}
}
