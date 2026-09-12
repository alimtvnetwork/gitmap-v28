package downloaderconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateInFlightStageDir(t *testing.T) {
	stageDir, err := CreateInFlightStageDir()
	if err != nil {
		t.Fatalf("CreateInFlightStageDir failed: %v", err)
	}

	defer os.RemoveAll(stageDir)

	info, err := os.Stat(stageDir)
	if err != nil || !info.IsDir() {
		t.Errorf("Expected valid directory at %s", stageDir)
	}
}

func TestResolvePersistentKeepDir(t *testing.T) {
	keepDir, err := ResolvePersistentKeepDir()
	if err != nil {
		t.Fatalf("ResolvePersistentKeepDir failed: %v", err)
	}

	if keepDir == "" {
		t.Errorf("Expected non-empty keep dir")
	}
}

func TestEnsureKeepDirectoryStructure(t *testing.T) {
	tempDir := t.TempDir()
	keepDir := filepath.Join(tempDir, ".gitmap-installation")

	if err := EnsureKeepDirectoryStructure(keepDir); err != nil {
		t.Fatalf("EnsureKeepDirectoryStructure failed: %v", err)
	}

	subdirs := []string{"downloads", "scripts", "logs"}
	for _, sub := range subdirs {
		path := filepath.Join(keepDir, sub)
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			t.Errorf("Expected subdir %s to exist as a directory", sub)
		}
	}
}

func TestPromoteStagedFile(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "stage.part")
	dstFile := filepath.Join(tempDir, "keep", "downloads", "final.deb")

	expectedContent := "gitmap test package content"
	if err := os.WriteFile(srcFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	if err := PromoteStagedFile(srcFile, dstFile); err != nil {
		t.Fatalf("PromoteStagedFile failed: %v", err)
	}

	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Errorf("Expected source file to be removed")
	}

	readBytes, err := os.ReadFile(dstFile)
	if err != nil || string(readBytes) != expectedContent {
		t.Errorf("Expected destination file content %q, got %q (err: %v)", expectedContent, string(readBytes), err)
	}
}
