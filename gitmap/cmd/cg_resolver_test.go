package cmd

import (
	"testing"
)

func TestCGRepoResolver(t *testing.T) {
	tempDir := t.TempDir()

	path, ok := ResolveCGTarget(tempDir)
	if !ok || path == "" {
		t.Fatalf("expected resolving local path to succeed")
	}

	all := ResolveAllCGTargets([]string{tempDir, tempDir})
	if len(all) != 1 {
		t.Fatalf("expected deduplicated resolution, got %d", len(all))
	}
}
