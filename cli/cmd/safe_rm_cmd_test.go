package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunSafeRmCLI_NonExistentTarget(t *testing.T) {
	nonExistent := filepath.Join(t.TempDir(), "non_existent_folder_or_file")
	err := RunSafeRmCLI([]string{nonExistent})
	if err != nil {
		t.Fatalf("expected RunSafeRmCLI to succeed on non-existent path, got: %v", err)
	}
}

func TestRunSafeRmCLI_ExistingFileAndDirectory(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testSubDir := filepath.Join(tempDir, "sub_dir")

	if err := os.WriteFile(testFile, []byte("hello"), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.MkdirAll(testSubDir, 0o755); err != nil {
		t.Fatalf("failed to create test sub dir: %v", err)
	}

	err := RunSafeRmCLI([]string{"--force", testFile, testSubDir})
	if err != nil {
		t.Fatalf("expected RunSafeRmCLI to succeed on existing paths, got: %v", err)
	}

	if isSafeRmPathExists(testFile) {
		t.Errorf("expected testFile to be removed: %s", testFile)
	}
	if isSafeRmPathExists(testSubDir) {
		t.Errorf("expected testSubDir to be removed: %s", testSubDir)
	}
}

func TestRunSafeRmCLI_NoArgs(t *testing.T) {
	err := RunSafeRmCLI([]string{})
	if err != nil {
		t.Fatalf("expected RunSafeRmCLI with no args to succeed (printing usage), got: %v", err)
	}
}
