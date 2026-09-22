package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRenameCandidates(t *testing.T) {
	tempDir := t.TempDir()
	readmePath := filepath.Join(tempDir, "README.md")
	_ = os.WriteFile(readmePath, []byte("# Hello"), 0644)
	subDir := filepath.Join(tempDir, "docs")
	_ = os.MkdirAll(subDir, 0755)
	subPath := filepath.Join(subDir, "GUIDE.MD")
	_ = os.WriteFile(subPath, []byte("# Guide"), 0644)
	lowerPath := filepath.Join(tempDir, "lower.txt")
	_ = os.WriteFile(lowerPath, []byte("ok"), 0644)

	opts := LowerCaseFixOptions{Patterns: []string{"*.md"}}
	pairs, scanned, err := findRenameCandidates(tempDir, opts)
	if err != nil {
		t.Fatalf("findRenameCandidates error: %v", err)
	}
	if scanned < 3 {
		t.Errorf("expected at least 3 files scanned, got %d", scanned)
	}
	if len(pairs) != 2 {
		t.Errorf("expected 2 rename candidates, got %d", len(pairs))
	}
}

func TestPerformTwoStepFSRename(t *testing.T) {
	tempDir := t.TempDir()
	src := filepath.Join(tempDir, "FILE.txt")
	tmp := filepath.Join(tempDir, "FILE.txt.tmp-lcf")
	dst := filepath.Join(tempDir, "file.txt")
	_ = os.WriteFile(src, []byte("content"), 0644)

	err := performTwoStepFSRename(src, tmp, dst)
	if err != nil {
		t.Fatalf("performTwoStepFSRename error: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("expected %s to exist: %v", dst, err)
	}
}
