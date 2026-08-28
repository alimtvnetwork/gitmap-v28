package committransfer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDirtySource_RejectsRunRight ensures that commit-right will not run if the source repository
// has uncommitted changes, preventing the detached HEAD stranding bug.
func TestDirtySource_RejectsRunRight(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	source := mustInitCountRepo(t, filepath.Join(root, "src"))
	target := mustInitCountRepo(t, filepath.Join(root, "dst"))

	// Create 2 commits in source
	mustCommitCount(t, source, "file1.txt", "content1", "commit 1")
	mustCommitCount(t, source, "file2.txt", "content2", "commit 2")

	// Dirtify source working tree
	dirtyFile := filepath.Join(source, "file1.txt")
	if err := os.WriteFile(dirtyFile, []byte("uncommitted change"), 0644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	opts := Options{
		LogPrefix: "[commit-right]",
		DryRun:    false,
		Yes:       true,
		NoPush:    true,
	}

	err := RunRight(source, target, opts)
	if err == nil {
		t.Fatal("expected RunRight to fail on dirty source, but it succeeded")
	}

	if !strings.Contains(err.Error(), "uncommitted changes") {
		t.Fatalf("expected uncommitted changes error, got: %v", err)
	}
}

