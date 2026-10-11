package cmdscan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestTryIncludeRepoInRepoCacheJSON(t *testing.T) {
	tempRC := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempRC, ".git"), 0755); err != nil {
		t.Fatalf("failed to create temp git dir: %v", err)
	}

	oldHook := ResolveRepoCacheRootFn
	ResolveRepoCacheRootFn = func() string { return tempRC }
	defer func() { ResolveRepoCacheRootFn = oldHook }()

	rec := model.ScanRecord{
		RepoName:     "test-added-repo",
		Slug:         "test-added-repo",
		HTTPSUrl:     "https://github.com/org/test-added-repo.git",
		RelativePath: "test-added-repo",
		Branch:       "main",
	}

	TryIncludeRepoInRepoCacheJSON([]model.ScanRecord{rec})

	manifestPath := filepath.Join(tempRC, "01-gitmap", "gitmap.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Fatalf("expected manifest at %s, but file not found", manifestPath)
	}

	recs, _, err := readExportRecords(manifestPath)
	if err != nil {
		t.Fatalf("failed to read export records: %v", err)
	}

	if len(recs) != 1 || recs[0].RepoName != "test-added-repo" {
		t.Fatalf("unexpected records in manifest: %+v", recs)
	}
}
