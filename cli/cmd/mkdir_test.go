package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveMkdirAbsPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine user home dir")
	}

	res, err := resolveMkdirAbsPath("~/test_mkdir_subfolder")
	if err != nil {
		t.Fatalf("resolveMkdirAbsPath failed: %v", err)
	}

	expectedPrefix := filepath.Clean(home)
	if !strings.HasPrefix(filepath.Clean(res), expectedPrefix) {
		t.Errorf("expected path starting with %s, got %s", expectedPrefix, res)
	}
}

func TestRunMkdir(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "sample_dir")

	err := runMkdir([]string{testPath})
	if err != nil {
		t.Fatalf("runMkdir failed: %v", err)
	}

	info, err := os.Stat(testPath)
	if err != nil || !info.IsDir() {
		t.Errorf("expected directory to exist: %s", testPath)
	}
}
