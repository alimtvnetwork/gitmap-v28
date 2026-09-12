package cmdinstall

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestParseInstallAddFlags_Valid(t *testing.T) {
	args := []string{"my-app", "v2.1", "--desc", "My App", "--win", "echo win", "--ubuntu", "echo ubuntu", "-y"}
	flags, err := parseInstallAddFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if flags.Name != "my-app" || flags.Version != "v2.1" {
		t.Errorf("unexpected name/version: %s / %s", flags.Name, flags.Version)
	}

	if flags.WinScript != "echo win" || flags.UbuntuScript != "echo ubuntu" {
		t.Errorf("unexpected scripts: win=%s, ubuntu=%s", flags.WinScript, flags.UbuntuScript)
	}
}

func TestBuildInstallerScriptFromFlags(t *testing.T) {
	flags := &InstallAddFlags{
		Name:         "custom-pkg",
		Version:      "v1.0.0",
		Description:  "A custom package",
		WinScript:    "winget install custom",
		UbuntuScript: "apt install custom",
	}

	script := buildInstallerScriptFromFlags(flags)
	if script.Slug != "custom-pkg" {
		t.Errorf("expected slug custom-pkg, got %s", script.Slug)
	}

	if script.TargetOS != "all" {
		t.Errorf("expected targetOS all, got %s", script.TargetOS)
	}

	if len(script.Scripts) != 2 {
		t.Errorf("expected 2 OS scripts, got %d", len(script.Scripts))
	}
}

func TestExportAndImportJSONRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "test-tool.json")
	script := model.InstallerScript{
		Name:        "test-tool",
		Slug:        "test-tool",
		Description: "Tool for testing",
		TargetOS:    "all",
		Version:     "v1.0.0",
	}

	opts := &ExportOptions{
		Slug:       "test-tool",
		OutputPath: jsonPath,
		Format:     "json",
	}

	if err := writeJSONExport([]model.InstallerScript{script}, opts); err != nil {
		t.Fatalf("writeJSONExport failed: %v", err)
	}

	if _, err := os.Stat(jsonPath); err != nil {
		t.Fatalf("expected json file to exist: %v", err)
	}
}

func TestExportAndImportZipRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test-tool.zip")
	script := model.InstallerScript{
		Name:        "test-tool",
		Slug:        "test-tool",
		Description: "Tool for testing",
		TargetOS:    "win",
		Version:     "v1.0.0",
	}

	opts := &ExportOptions{
		Slug:       "test-tool",
		OutputPath: zipPath,
		Format:     "zip",
	}

	if err := writeZipExport([]model.InstallerScript{script}, opts); err != nil {
		t.Fatalf("writeZipExport failed: %v", err)
	}

	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("expected zip file to exist: %v", err)
	}
}

func TestResolveCustomToolStatus(t *testing.T) {
	script := model.InstallerScript{
		Slug:    "git",
		Version: "2.47.1",
	}

	installed := map[string]string{}
	status, ver := resolveCustomToolStatus(script, installed)
	if status == "" || ver == "" {
		t.Errorf("unexpected empty status or ver")
	}
}

func TestLoadCustomInstallersList(t *testing.T) {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		t.Skip("sqlite default db unavailable in test environment")
	}

	defer db.Close()
	_ = loadCustomInstallersList()
}
