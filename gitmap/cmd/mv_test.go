package cmd

import (
	"path/filepath"
	"testing"
)

func TestCalculateDestPath(t *testing.T) {
	src := `D:\work\projects\my-repo`

	// 1. Double dot destination (parent)
	dest, err := calculateDestPath(src, "..")
	if err != nil {
		t.Fatalf("calculateDestPath .. failed: %v", err)
	}
	expected := `D:\work\my-repo`
	if filepath.Clean(dest) != filepath.Clean(expected) {
		t.Errorf("calculateDestPath .. got %s, want %s", dest, expected)
	}

	// 2. Preflight identical check
	if err := preflightMove(src, src); err == nil {
		t.Errorf("preflightMove expected error for identical paths")
	}

	// 3. Preflight inside self check
	nested := filepath.Join(src, "subfolder")
	if err := preflightMove(src, nested); err == nil {
		t.Errorf("preflightMove expected error for nested path")
	}
}
