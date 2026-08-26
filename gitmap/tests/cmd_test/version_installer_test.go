package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestVersionManifestInheritance(t *testing.T) {
	manifest := model.RepositoryVersionManifest{
		Version: "6.108.0",
		Backend: &model.ComponentVersion{
			Version: "inherit",
		},
		Frontend: &model.ComponentVersion{
			Version: "3.2.1",
		},
	}

	if backendVer := manifest.ResolveVersion("backend"); backendVer != "6.108.0" {
		t.Fatalf("expected backend to inherit root version 6.108.0, got %s", backendVer)
	}

	if frontendVer := manifest.ResolveVersion("frontend"); frontendVer != "3.2.1" {
		t.Fatalf("expected frontend explicit version 3.2.1, got %s", frontendVer)
	}

	if cliVer := manifest.ResolveVersion("cli"); cliVer != "6.108.0" {
		t.Fatalf("expected unconfigured component to default to root version 6.108.0, got %s", cliVer)
	}
}

func TestVersionInstallerDryRun(t *testing.T) {
	tempDir := t.TempDir()
	cfg := cmd.DefaultVersionInstallConfig("1.0.0")

	if err := cmd.InstallVersionJSON(tempDir, cfg, true); err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "version.json")); !os.IsNotExist(err) {
		t.Fatal("dry-run should not create version.json on disk")
	}
}
