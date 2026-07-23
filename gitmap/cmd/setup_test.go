package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestResolveSetupConfigPathPreservesCustomPath(t *testing.T) {
	want := filepath.Join(t.TempDir(), "custom-git-setup.json")
	got := resolveSetupConfigPath(want, true)
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestResolveDefaultSetupConfigPathPrefersBinaryDataDir(t *testing.T) {
	binaryDataDir := t.TempDir()
	want := writeSetupConfig(t, binaryDataDir)
	got := resolveDefaultSetupConfigPath(constants.DefaultSetupConfigPath, binaryDataDir, "")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestResolveDefaultSetupConfigPathFallsBackToRepoPath(t *testing.T) {
	root := t.TempDir()
	repoPath := filepath.Join(root, "source")
	wantDir := filepath.Join(repoPath, constants.GitMapSubdir, constants.DBDir)
	want := writeSetupConfig(t, wantDir)
	got := resolveDefaultSetupConfigPath(constants.DefaultSetupConfigPath, filepath.Join(root, "missing"), repoPath)
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestResolveDefaultSetupConfigPathUsesBinaryFallbackWhenMissing(t *testing.T) {
	binaryDataDir := t.TempDir()
	want := filepath.Join(binaryDataDir, filepath.Base(constants.DefaultSetupConfigPath))
	got := resolveDefaultSetupConfigPath(constants.DefaultSetupConfigPath, binaryDataDir, "")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func writeSetupConfig(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, filepath.Base(constants.DefaultSetupConfigPath))
	if err := os.MkdirAll(dir, constants.DirPermission); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}"), constants.FilePermission); err != nil {
		t.Fatal(err)
	}
	return path
}
