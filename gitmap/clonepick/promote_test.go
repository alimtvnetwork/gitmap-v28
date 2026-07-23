package clonepick

// promote_test.go: tests the rename + cross-fs copy fallback path
// used by the --ask clone-once optimisation. We can't easily force
// a real EXDEV in the test sandbox, so the cross-fs copy is exercised
// via copyTreeThenRemove directly.

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPromotePreClonedSrcRenameFastPath(t *testing.T) {
	src, dest := setupPromoteDirs(t, "rename-src", "rename-dest")
	mustWrite(t, filepath.Join(src, "a.txt"), "alpha")
	mustWrite(t, filepath.Join(src, "sub", "b.txt"), "beta")

	if err := promotePreClonedSrc(src, dest); err != nil {
		t.Fatalf("promotePreClonedSrc: %v", err)
	}
	assertFile(t, filepath.Join(dest, "a.txt"), "alpha")
	assertFile(t, filepath.Join(dest, "sub", "b.txt"), "beta")
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("src should be gone after promote, got err=%v", err)
	}
}

func TestCopyTreeThenRemoveFallback(t *testing.T) {
	src, dest := setupPromoteDirs(t, "copy-src", "copy-dest")
	mustWrite(t, filepath.Join(src, "x.txt"), "xray")
	mustWrite(t, filepath.Join(src, "deep", "nested", "y.txt"), "yankee")

	if err := copyTreeThenRemove(src, dest); err != nil {
		t.Fatalf("copyTreeThenRemove: %v", err)
	}
	assertFile(t, filepath.Join(dest, "x.txt"), "xray")
	assertFile(t, filepath.Join(dest, "deep", "nested", "y.txt"), "yankee")
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("src should be removed after copy, got err=%v", err)
	}
}

// setupPromoteDirs returns (src, dest). src exists with content
// owned by the caller; dest exists and is empty (mirrors prepareDest
// contract).
func setupPromoteDirs(t *testing.T, srcName, destName string) (string, string) {
	t.Helper()
	root := t.TempDir()
	src := filepath.Join(root, srcName)
	dest := filepath.Join(root, destName)
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}

	return src, dest
}

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func assertFile(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("file %s body = %q, want %q", path, got, want)
	}
}
