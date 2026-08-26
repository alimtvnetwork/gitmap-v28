package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
)

func TestCGMultiRepoSuite(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Test child repo discovery
	repo1 := filepath.Join(tempDir, "repo1")
	os.MkdirAll(filepath.Join(repo1, ".git"), 0755)

	repos, errDiscover := fsutil.DiscoverChildGitRepos(tempDir)
	if errDiscover != nil || len(repos) != 1 {
		t.Fatalf("expected 1 discovered child repo, got %v (err: %v)", repos, errDiscover)
	}

	// 2. Test version.json writing and reading
	meta := cmd.CGMetadata{
		Version: "v24.2.0",
		Status:  "active",
	}
	if errWrite := cmd.WriteCGMetadata(repo1, meta); errWrite != nil {
		t.Fatalf("WriteCGMetadata failed: %v", errWrite)
	}

	readMeta, errRead := cmd.ReadCGMetadata(repo1)
	if errRead != nil || readMeta.Version != "v24.2.0" {
		t.Fatalf("ReadCGMetadata unexpected result: %+v (err: %v)", readMeta, errRead)
	}

	// 3. Test resolve target
	resolved, ok := cmd.ResolveCGTarget(repo1)
	if !ok || resolved == "" {
		t.Fatalf("expected resolved path, got %s (ok: %v)", resolved, ok)
	}
}
