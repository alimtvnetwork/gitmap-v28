package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupLegacyDeployDirNoOpWhenNotInGitmapCli(t *testing.T) {
	t.Parallel()

	ctx := updateCleanupContext{selfPath: "C:/bin/gitmap.exe"}
	removed := cleanupLegacyDeployDir(ctx)
	if removed != 0 {
		t.Errorf("expected 0 removals for non-gitmap-cli path, got %d", removed)
	}
}

func TestCleanupLegacyDeployDirRemovesLegacyBinary(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cliDir := filepath.Join(tmpDir, "gitmap-cli")
	legacyDir := filepath.Join(tmpDir, "gitmap")
	setupTestLegacyDirs(t, cliDir, legacyDir)

	selfPath := filepath.Join(cliDir, "gitmap.exe")
	ctx := updateCleanupContext{selfPath: selfPath}
	removed := cleanupLegacyDeployDir(ctx)
	assertLegacyCleanupResult(t, removed, legacyDir, selfPath)
}

func setupTestLegacyDirs(t *testing.T, cliDir, legacyDir string) {
	t.Helper()

	if err := os.MkdirAll(cliDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	createTestFile(t, filepath.Join(cliDir, "gitmap.exe"))
	createTestFile(t, filepath.Join(legacyDir, "gitmap.exe"))
	createTestFile(t, filepath.Join(legacyDir, "gitmap.ps1"))
}

func createTestFile(t *testing.T, path string) {
	t.Helper()

	if err := os.WriteFile(path, []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertLegacyCleanupResult(t *testing.T, removed int, legacyDir, selfPath string) {
	t.Helper()

	if removed != 2 {
		t.Errorf("expected 2 removed files, got %d", removed)
	}
	if _, err := os.Stat(filepath.Join(legacyDir, "gitmap.exe")); !os.IsNotExist(err) {
		t.Errorf("expected legacy binary to be deleted")
	}
	if _, err := os.Stat(selfPath); err != nil {
		t.Errorf("expected active binary to remain, err: %v", err)
	}
}
