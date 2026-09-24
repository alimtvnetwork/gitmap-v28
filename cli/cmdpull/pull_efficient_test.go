package cmdpull

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestExtractEfficientFlags(t *testing.T) {
	args := []string{"--status", "--ssh", "-t", "ssh", "--json", "extra-arg"}
	useSSH, useHTTPS, isTable, isJSON, targetSSH, rest := extractEfficientFlags(args)

	if !isTable {
		t.Fatalf("expected isTable to be true")
	}
	if !useSSH {
		t.Fatalf("expected useSSH to be true")
	}
	if useHTTPS {
		t.Fatalf("expected useHTTPS to be false")
	}
	if !isJSON {
		t.Fatalf("expected isJSON to be true")
	}
	if targetSSH != "ssh" {
		t.Fatalf("expected targetSSH to be 'ssh', got %s", targetSSH)
	}
	if len(rest) != 1 || rest[0] != "extra-arg" {
		t.Fatalf("expected rest args to have 1 element, got %v", rest)
	}
}

func TestResolveEfficientFullCmdName(t *testing.T) {
	if got := resolveEfficientFullCmdName(true); got != "pull all-efficient-table" {
		t.Fatalf("expected pull all-efficient-table, got %s", got)
	}
	if got := resolveEfficientFullCmdName(false); got != "pull all-efficient" {
		t.Fatalf("expected pull all-efficient, got %s", got)
	}
}

func TestPartitionRecordsByActivity_WithMockData(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "gitmap-pull.db")

	db, err := store.OpenPullSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenPullSplitDBAt failed: %v", err)
	}
	defer db.Close()

	repoActive := filepath.Join(tempDir, "active-repo")
	repoInactive := filepath.Join(tempDir, "inactive-repo")

	runID, err := db.InsertPullRun(&store.PullRunRecord{
		CommandType: "pull-all",
		WorkingDir:  tempDir,
	})
	if err != nil {
		t.Fatalf("InsertPullRun failed: %v", err)
	}

	// Insert 20 zero-change records for inactive repo
	var batch []store.PullRepoRunRecord
	for i := 0; i < 20; i++ {
		batch = append(batch, store.PullRepoRunRecord{
			RepoPath:     repoInactive,
			RepoName:     "inactive-repo",
			PullStatus:   "up-to-date",
			FilesChanged: 0,
			IsActive:     true,
			HasChanges:   false,
		})
	}
	// Insert 1 record with changes for active repo
	batch = append(batch, store.PullRepoRunRecord{
		RepoPath:     repoActive,
		RepoName:     "active-repo",
		PullStatus:   "success",
		FilesChanged: 5,
		IsActive:     true,
		HasChanges:   true,
	})

	if err := db.InsertPullRepoRuns(runID, batch); err != nil {
		t.Fatalf("InsertPullRepoRuns failed: %v", err)
	}

	if _, err := db.Conn().Exec("UPDATE PullRepoRun SET CreatedAt = datetime('now', '-10 minutes')"); err != nil {
		t.Fatalf("failed to update CreatedAt: %v", err)
	}

	records := []model.ScanRecord{
		{RepoName: "active-repo", AbsolutePath: repoActive},
		{RepoName: "inactive-repo", AbsolutePath: repoInactive},
	}

	var part EfficientPullPartition
	for _, rec := range records {
		classifyRepoActivity(db, rec, &part)
	}

	if len(part.ActiveRecords) != 1 || part.ActiveRecords[0].RepoName != "active-repo" {
		t.Fatalf("expected 1 active record (active-repo), got %v", part.ActiveRecords)
	}
	if len(part.InactiveRepos) != 1 || part.InactiveRepos[0].RepoName != "inactive-repo" {
		t.Fatalf("expected 1 inactive repo (inactive-repo), got %v", part.InactiveRepos)
	}
}
