package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdir_DeepDirectory(t *testing.T) {
	tempDir := t.TempDir()
	deepDir := filepath.Join(tempDir, "a", "b", "c", "d")
	if err := runMkdir([]string{"-p", deepDir}); err != nil {
		t.Fatalf("runMkdir failed: %v", err)
	}
	if info, err := os.Stat(deepDir); err != nil || !info.IsDir() {
		t.Fatalf("expected dir %s to exist", deepDir)
	}
}

func TestMkdir_FlagPositionIndependence(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "x", "y")
	if err := runMkdir([]string{targetDir, "-p"}); err != nil {
		t.Fatalf("runMkdir with flag after target failed: %v", err)
	}
	if info, err := os.Stat(targetDir); err != nil || !info.IsDir() {
		t.Fatalf("expected dir %s to exist", targetDir)
	}
}

func TestMkdir_FileCreation(t *testing.T) {
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "nested", "folder", "sample.txt")
	if err := runMkdir([]string{"-f", targetFile}); err != nil {
		t.Fatalf("runMkdir -f failed: %v", err)
	}
	if info, err := os.Stat(targetFile); err != nil || info.IsDir() {
		t.Fatalf("expected file %s to exist", targetFile)
	}
}

func TestMkdir_Idempotency(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "idempotent_dir")
	if err := runMkdir([]string{"-p", targetDir}); err != nil {
		t.Fatalf("first runMkdir failed: %v", err)
	}
	if err := runMkdir([]string{"-p", targetDir}); err != nil {
		t.Fatalf("second runMkdir failed: %v", err)
	}
}

func TestMkdir_EmptyArgs(t *testing.T) {
	if err := runMkdir([]string{}); err == nil {
		t.Fatal("expected error on empty args, got nil")
	}
}

func TestMkdir_MixedSlashes(t *testing.T) {
	tempDir := t.TempDir()
	rawPath := tempDir + "/mixed//slash\\sub//dir"
	if err := runMkdir([]string{"-p", rawPath}); err != nil {
		t.Fatalf("runMkdir mixed slashes failed: %v", err)
	}
}
