package store

import (
	"path/filepath"
	"testing"
)

func TestPullSplitDB_SearchPullTraces(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-pull-search.db")

	db, err := OpenPullSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenPullSplitDBAt failed: %v", err)
	}
	defer db.Close()

	runRecord := &PullRunRecord{
		CommandType:   "pull all-efficient",
		WorkingDir:    tempDir,
		TotalRepos:    1,
		PulledRepos:   1,
		IsEfficient:   true,
		GitMapVersion: "v6.342.0",
	}
	runID, err := db.InsertPullRun(runRecord)
	if err != nil {
		t.Fatalf("InsertPullRun failed: %v", err)
	}

	repoRecords := []PullRepoRunRecord{
		{
			RepoPath:      filepath.Join(tempDir, "repo-alpha"),
			RepoName:      "repo-alpha",
			PullStatus:    "success",
			LastCommitSha: "abc1234",
			CommitMessage: "feat: add superfast commit trace search",
			CommitAuthor:  "ALIM",
			HasChanges:    true,
			Notes:         "abc1234 feat: add superfast commit trace search",
		},
	}
	if err := db.InsertPullRepoRuns(runID, repoRecords); err != nil {
		t.Fatalf("InsertPullRepoRuns failed: %v", err)
	}

	results, err := db.SearchPullTraces("superfast", 10)
	if err != nil {
		t.Fatalf("SearchPullTraces failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].RepoName != "repo-alpha" {
		t.Fatalf("expected repo-alpha, got %s", results[0].RepoName)
	}

	shaResults, err := db.SearchPullTraces("abc1234", 10)
	if err != nil {
		t.Fatalf("SearchPullTraces by SHA failed: %v", err)
	}
	if len(shaResults) != 1 {
		t.Fatalf("expected 1 result for SHA, got %d", len(shaResults))
	}
}
