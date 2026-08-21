package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathNormalizeAndSlashes(t *testing.T) {
	p := `D:\work\project\sub\`
	trimmed := TrimTrailingSlashes(p)
	if trimmed != `D:\work\project\sub` {
		t.Errorf("TrimTrailingSlashes got %s", trimmed)
	}

	norm := NormalizeSlashes(p)
	if norm != `D:/work/project/sub` {
		t.Errorf("NormalizeSlashes got %s", norm)
	}

	if !EqualPaths(`D:\work\foo`, `d:/work/foo/`) {
		t.Errorf("EqualPaths failed for case and slash variations")
	}

	if !IsSubdirectory(`D:/work/project`, `D:/work/project/nested/deep`) {
		t.Errorf("IsSubdirectory failed for nested child")
	}

	if IsSubdirectory(`D:/work/project`, `D:/work/other`) {
		t.Errorf("IsSubdirectory false positive for unrelated folder")
	}
}

func TestEnsureAndStripLongPath(t *testing.T) {
	raw := `D:\work\myrepo`
	long := EnsureLongPath(raw)
	if len(long) == 0 {
		t.Fatalf("EnsureLongPath returned empty string")
	}
	stripped := StripLongPathPrefix(long)
	if stripped != raw && filepath.Clean(stripped) != filepath.Clean(raw) {
		t.Errorf("StripLongPathPrefix mismatch: got %s, want %s", stripped, raw)
	}
}

func TestSafeFsOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-fsutil-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	src := filepath.Join(tempDir, "src")
	dst := filepath.Join(tempDir, "dst")
	_ = os.MkdirAll(src, 0755)
	_ = os.WriteFile(filepath.Join(src, "file.txt"), []byte("hello"), 0644)

	if err := SafeRename(src, dst); err != nil {
		t.Fatalf("SafeRename failed: %v", err)
	}

	if !FileExists(filepath.Join(dst, "file.txt")) {
		t.Errorf("Destination file missing after SafeRename")
	}

	if err := SafeRemoveAll(dst); err != nil {
		t.Fatalf("SafeRemoveAll failed: %v", err)
	}
}
