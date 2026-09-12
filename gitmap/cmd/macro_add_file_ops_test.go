package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunTouchAndCatCmd(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "subdir", "test_file.txt")

	if err := runTouchCmd([]string{filePath}); err != nil {
		t.Fatalf("runTouchCmd failed: %v", err)
	}

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		t.Fatalf("expected touched file to exist, err=%v", err)
	}

	if err := runCatCmd([]string{filePath}); err != nil {
		t.Fatalf("runCatCmd failed: %v", err)
	}
}

func TestRunMkfileCmd(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "nested", "sample.md")
	content := "Hello GitMap Macro File Ops"

	if err := runMkfileCmd([]string{filePath, content}); err != nil {
		t.Fatalf("runMkfileCmd failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil || string(data) != content {
		t.Fatalf("expected %q, got %q (err: %v)", content, string(data), err)
	}
}

func TestRunCpAndRmFileCmd(t *testing.T) {
	tempDir := t.TempDir()
	src := filepath.Join(tempDir, "source.txt")
	dst := filepath.Join(tempDir, "copied", "dest.txt")

	_ = os.WriteFile(src, []byte("copy payload"), 0644)
	if err := runCpFileCmd([]string{src, dst}); err != nil {
		t.Fatalf("runCpFileCmd failed: %v", err)
	}

	destData, err := os.ReadFile(dst)
	if err != nil || string(destData) != "copy payload" {
		t.Fatalf("copy content mismatch: %q", string(destData))
	}

	if err := runRmFileCmd([]string{src}); err != nil {
		t.Fatalf("runRmFileCmd failed: %v", err)
	}
}
