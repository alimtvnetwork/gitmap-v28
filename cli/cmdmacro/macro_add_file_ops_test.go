package cmdmacro

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

func TestRunCpDirectoryRecursive(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "git-work")
	_ = os.MkdirAll(filepath.Join(srcDir, "sub"), 0755)
	_ = os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("data1"), 0644)
	_ = os.WriteFile(filepath.Join(srcDir, "sub", "file2.txt"), []byte("data2"), 0644)

	targetDir := filepath.Join(tempDir, "test")
	_ = os.MkdirAll(targetDir, 0755)

	// Test cp -r srcDir targetDir
	if err := executeInteractiveCopy("cp ../git-work .", []string{"-r", srcDir, targetDir}); err != nil {
		t.Fatalf("executeInteractiveCopy failed: %v", err)
	}

	destFile1 := filepath.Join(targetDir, "git-work", "file1.txt")
	destFile2 := filepath.Join(targetDir, "git-work", "sub", "file2.txt")

	if _, err := os.Stat(destFile1); err != nil {
		t.Errorf("expected %s to exist: %v", destFile1, err)
	}
	if _, err := os.Stat(destFile2); err != nil {
		t.Errorf("expected %s to exist: %v", destFile2, err)
	}
}
