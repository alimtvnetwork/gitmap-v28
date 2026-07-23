package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// TestIsExcludedDirAndPrefix verifies the spec's two-tier exclusion
// rule: directory base names AND repo-relative path prefixes both
// short-circuit the walker.
func TestIsExcludedDirAndPrefix(t *testing.T) {
	for _, name := range []string{".git", ".gitmap", ".release", "node_modules", "vendor"} {
		if !isExcludedDir(name) {
			t.Errorf("isExcludedDir(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"src", "docs", "gitmap", "release"} {
		if isExcludedDir(name) {
			t.Errorf("isExcludedDir(%q) = true, want false", name)
		}
	}

	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, ".gitmap", "release", "v1"))
	mustMkdirAll(t, filepath.Join(root, ".gitmap", "release-assets"))
	mustMkdirAll(t, filepath.Join(root, "src"))

	if !isExcludedPrefix(root, filepath.Join(root, ".gitmap", "release")) {
		t.Error(".gitmap/release should be excluded")
	}
	if !isExcludedPrefix(root, filepath.Join(root, ".gitmap", "release", "v1")) {
		t.Error(".gitmap/release/v1 should be excluded (prefix match)")
	}
	if !isExcludedPrefix(root, filepath.Join(root, ".gitmap", "release-assets")) {
		t.Error(".gitmap/release-assets should be excluded")
	}
	if isExcludedPrefix(root, filepath.Join(root, "src")) {
		t.Error("src must not be excluded by prefix rule")
	}
}

// TestWalkRepoFilesSkipsExclusionsAndBinaries seeds a fake repo that
// covers every code path the walker promises to honor:
//   - excluded directory (.git) is skipped wholesale
//   - excluded prefix (.gitmap/release-assets) is skipped wholesale
//   - binary files (null byte in first 8 KiB) are skipped per file
//   - regular text files are included
func TestWalkRepoFilesSkipsExclusionsAndBinaries(t *testing.T) {
	root := t.TempDir()

	mustWriteFile(t, filepath.Join(root, "README.md"), []byte("hello\n"))
	mustWriteFile(t, filepath.Join(root, "src", "app.go"), []byte("package app\n"))
	mustWriteFile(t, filepath.Join(root, ".git", "HEAD"), []byte("ref: refs/heads/main\n"))
	mustWriteFile(t,
		filepath.Join(root, ".gitmap", "release-assets", "v1.zip"),
		[]byte("ignored text"),
	)
	mustWriteFile(t, filepath.Join(root, "image.png"), []byte("PNG\x00\x01\x02binary"))

	got, err := walkRepoFiles(root, nil, true)
	if err != nil {
		t.Fatalf("walkRepoFiles: %v", err)
	}

	rels := relativizeAll(t, root, got)
	sort.Strings(rels)

	want := []string{"README.md", "src/app.go"}
	if !equalStringSlice(rels, want) {
		t.Fatalf("walkRepoFiles returned %v, want %v", rels, want)
	}
}

// TestIsBinaryFile pins the null-byte sniff: a file with a NUL in the
// first 8 KiB is binary; otherwise it's text.
func TestIsBinaryFile(t *testing.T) {
	dir := t.TempDir()
	textPath := filepath.Join(dir, "a.txt")
	binPath := filepath.Join(dir, "b.bin")

	mustWriteFile(t, textPath, []byte("plain ascii content\n"))
	mustWriteFile(t, binPath, []byte{'h', 'i', 0x00, 'x'})

	if isBinaryFile(textPath) {
		t.Error("plain ascii file misclassified as binary")
	}
	if !isBinaryFile(binPath) {
		t.Error("file with null byte not classified as binary")
	}
}

// --- helpers ---

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	mustMkdirAll(t, filepath.Dir(path))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("WriteFile %s: %v", path, err)
	}
}

func relativizeAll(t *testing.T, root string, paths []string) []string {
	t.Helper()
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		rel, err := filepath.Rel(root, p)
		if err != nil {
			t.Fatalf("Rel: %v", err)
		}
		out = append(out, filepath.ToSlash(rel))
	}
	return out
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
