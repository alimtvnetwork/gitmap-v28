// Package fsutil — path_normalize_test.go tests path normalization utilities.
package fsutil

import (
	"runtime"
	"strings"
	"testing"
)

func TestPathNormalize(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"a\\b\\c", "a/b/c"},
		{"a/b/c", "a/b/c"},
		{"a/b/../c", "a/c"},
	}

	for _, c := range cases {
		out := NormalizeToForwardSlashes(c.input)
		if out != c.expected {
			t.Errorf("NormalizeToForwardSlashes(%q) = %q, expected %q", c.input, out, c.expected)
		}
	}

	rel, err := MakeRelativeToRoot("/root/work", "/root/work/dir/file.txt")
	if err != nil {
		t.Fatalf("MakeRelativeToRoot failed: %v", err)
	}

	if rel != "dir/file.txt" {
		t.Errorf("expected dir/file.txt, got %s", rel)
	}
}

func TestEqualPaths_OSAware(t *testing.T) {
	// Windows drive letter and mixed slash comparisons
	p1 := "c:/work/repo"
	p2 := "C:\\work\\repo"
	if !EqualPaths(p1, p2) {
		t.Errorf("EqualPaths(%q, %q) expected true, got false", p1, p2)
	}

	p3 := "D:/Projects/App"
	p4 := "d:/projects/app"
	if !EqualPaths(p3, p4) {
		t.Errorf("EqualPaths(%q, %q) expected true for drive paths, got false", p3, p4)
	}

	// Test IsPathCaseInsensitive explicitly
	if !IsPathCaseInsensitive("C:/Users/Administrator") {
		t.Errorf("IsPathCaseInsensitive(C:/Users/Administrator) expected true")
	}
	if !IsPathCaseInsensitive("d:\\work\\gitmap") {
		t.Errorf("IsPathCaseInsensitive(d:\\work\\gitmap) expected true")
	}

	if runtime.GOOS != "windows" {
		// On non-Windows platforms, pure relative or Unix absolute paths must be case-sensitive
		u1 := "/home/user/Repo"
		u2 := "/home/user/repo"
		if EqualPaths(u1, u2) {
			t.Errorf("EqualPaths(%q, %q) on Unix expected false, got true", u1, u2)
		}
	}
}

func TestCanonicalPathKey(t *testing.T) {
	if CanonicalPathKey("") != "" {
		t.Errorf("CanonicalPathKey(\"\") expected \"\", got %q", CanonicalPathKey(""))
	}

	// Windows-style drive path should be lowercased
	winPath := "C:\\Work\\MyProject/SubDir"
	key := CanonicalPathKey(winPath)
	expectedWin := "c:/work/myproject/subdir"
	if key != expectedWin {
		t.Errorf("CanonicalPathKey(%q) = %q, expected %q", winPath, key, expectedWin)
	}

	if runtime.GOOS != "windows" {
		unixPath := "/home/User/Code/MyRepo"
		unixKey := CanonicalPathKey(unixPath)
		if unixKey != "/home/User/Code/MyRepo" {
			t.Errorf("CanonicalPathKey(%q) on Unix = %q, expected %q", unixPath, unixKey, unixPath)
		}
	}
}

func TestIsSubdirectory(t *testing.T) {
	parent := "C:/Work"
	child := "c:/work/repo/sub"
	if !IsSubdirectory(parent, child) {
		t.Errorf("IsSubdirectory(%q, %q) expected true", parent, child)
	}

	parent2 := "C:/Work/Repo"
	child2 := "C:/Work/Other"
	if IsSubdirectory(parent2, child2) {
		t.Errorf("IsSubdirectory(%q, %q) expected false", parent2, child2)
	}
}

func BenchmarkEqualPaths(b *testing.B) {
	p1 := "C:/Work/GitMap/cli/fsutil/path_normalize.go"
	p2 := "c:/work/gitmap/cli/fsutil/path_normalize.go"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !strings.EqualFold(p1, p2) {
			b.Fail()
		}
	}
}
