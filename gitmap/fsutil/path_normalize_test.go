// Package fsutil — path_normalize_test.go tests path normalization utilities.
package fsutil

import (
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
