package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
)

func TestRecursiveDiscoverySuite(t *testing.T) {
	tempDir := t.TempDir()

	repo1 := filepath.Join(tempDir, "workspace", "repo1")
	os.MkdirAll(filepath.Join(repo1, ".git"), 0755)

	repo2 := filepath.Join(tempDir, "workspace", "repo2")
	os.MkdirAll(filepath.Join(repo2, ".git"), 0755)

	records := cmd.ResolvePullDirectoryTargets(filepath.Join(tempDir, "workspace"))
	if len(records) != 2 {
		t.Fatalf("expected 2 resolved records, got %d", len(records))
	}
}
