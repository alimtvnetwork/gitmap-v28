package cmd

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func setupMvRmTestDB(t *testing.T) *store.DB {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_mv_rm.db")

	db, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.Migrate(); err != nil {
		db.Close()
		t.Fatalf("failed to migrate test db: %v", err)
	}

	return db
}

func insertInitialTestRepo(t *testing.T, db *store.DB, absPath, name string) int64 {
	t.Helper()
	rec := model.ScanRecord{
		AbsolutePath: absPath,
		RepoName:     name,
		Slug:         name,
	}

	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		t.Fatalf("failed to insert initial repo: %v", err)
	}

	repos, err := db.FindByPath(absPath)
	if err != nil || len(repos) == 0 {
		t.Fatalf("failed to find inserted repo: %v", err)
	}

	return repos[0].ID
}

func insertInitialTestAlias(t *testing.T, db *store.DB, repoID int64, alias string) {
	t.Helper()
	if _, err := db.CreateAlias(alias, repoID); err != nil {
		t.Fatalf("failed to insert alias: %v", err)
	}
}

func TestUpdateRepoInDB_Success(t *testing.T) {
	db := setupMvRmTestDB(t)
	defer db.Close()

	repoID := insertInitialTestRepo(t, db, "/path/old", "old-repo")
	insertInitialTestAlias(t, db, repoID, "my-alias")

	err := updateRepoInDB(db, repoID, "/path/new", "new-repo")
	if err != nil {
		t.Fatalf("updateRepoInDB failed: %v", err)
	}

	assertRepoUpdated(t, db, repoID, "/path/new", "new-repo")
	assertAliasPresent(t, db, repoID, "my-alias")
}

func assertRepoUpdated(t *testing.T, db *store.DB, repoID int64, expectedPath, expectedName string) {
	t.Helper()
	var actualPath, actualName string
	err := db.Conn().QueryRow("SELECT AbsolutePath, RepoName FROM Repo WHERE RepoId = ?", repoID).Scan(&actualPath, &actualName)
	if err != nil {
		t.Fatalf("failed to query updated repo: %v", err)
	}

	if actualPath != expectedPath || actualName != expectedName {
		t.Errorf("repo mismatch: got (%s, %s), want (%s, %s)", actualPath, actualName, expectedPath, expectedName)
	}
}

func assertAliasPresent(t *testing.T, db *store.DB, repoID int64, aliasName string) {
	t.Helper()
	alias, err := db.FindAliasByName(aliasName)
	if err != nil || alias.RepoID != repoID {
		t.Fatalf("alias not found or mismatched: %+v, err: %v", alias, err)
	}
}

func TestRemoveRepoDB_Success(t *testing.T) {
	db := setupMvRmTestDB(t)
	defer db.Close()

	repoID := insertInitialTestRepo(t, db, "/path/to/rm", "rm-repo")
	insertInitialTestAlias(t, db, repoID, "rm-alias")

	scanRec := model.ScanRecord{
		ID:           repoID,
		AbsolutePath: "/path/to/rm",
		Slug:         "rm-repo",
	}

	if err := removeRepoDB(db, scanRec); err != nil {
		t.Fatalf("removeRepoDB failed: %v", err)
	}

	assertRepoAndAliasDeleted(t, db, repoID)
	assertCommandHistoryRecorded(t, db, "/path/to/rm")
}

func assertRepoAndAliasDeleted(t *testing.T, db *store.DB, repoID int64) {
	t.Helper()
	var count int
	_ = db.Conn().QueryRow("SELECT COUNT(*) FROM Repo WHERE RepoId = ?", repoID).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 repos, got %d", count)
	}

	_ = db.Conn().QueryRow("SELECT COUNT(*) FROM Alias WHERE RepoId = ?", repoID).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 aliases, got %d", count)
	}
}

func assertCommandHistoryRecorded(t *testing.T, db *store.DB, absPath string) {
	t.Helper()
	var histCount int
	_ = db.Conn().QueryRow("SELECT COUNT(*) FROM CommandHistory WHERE Command = 'rm' AND Args = ?", absPath).Scan(&histCount)
	if histCount == 0 {
		t.Errorf("expected command history record for rm %s, got 0", absPath)
	}
}

func TestRemoveRepoDB_RollbackOnError(t *testing.T) {
	db := setupMvRmTestDB(t)
	defer db.Close()

	repoID := insertInitialTestRepo(t, db, "/path/to/rollback", "rb-repo")

	// Drop CommandHistory table to force recordRmHistory to fail.
	if _, err := db.Conn().Exec("DROP TABLE CommandHistory"); err != nil {
		t.Fatalf("failed to drop CommandHistory: %v", err)
	}

	scanRec := model.ScanRecord{
		ID:           repoID,
		AbsolutePath: "/path/to/rollback",
		Slug:         "rb-repo",
	}

	err := removeRepoDB(db, scanRec)
	if err == nil {
		t.Fatalf("expected error from removeRepoDB when history fails, got nil")
	}

	// Verify rollback preserved the repo row.
	var count int
	_ = db.Conn().QueryRow("SELECT COUNT(*) FROM Repo WHERE RepoId = ?", repoID).Scan(&count)
	if count != 1 {
		t.Errorf("expected repo row to be preserved after rollback, count=%d", count)
	}
}

func TestPipelineRecorderAndSyncCacheHelpers(t *testing.T) {
	slug := "test-owner/test-pipe-tx"
	pipeDb, err := pipelinedb.OpenPipelineSplitDB(slug)
	if err != nil {
		t.Fatalf("failed to open split db: %v", err)
	}

	defer pipeDb.Close()

	testRecordRunAndErrorInSplitDb(t, pipeDb, slug)
}

func testRecordRunAndErrorInSplitDb(t *testing.T, pipeDb *pipelinedb.PipelineSplitDb, slug string) {
	t.Helper()
	runItem := ghRunItem{
		DatabaseId: 991001,
		Name:       "CI",
		Status:     "completed",
		Conclusion: "failure",
		HeadBranch: "main",
		HeadSha:    "abc1234",
		Url:        "https://github.com/test/run/991001",
	}

	if err := recordRunInSplitDb(pipeDb, slug, runItem); err != nil {
		t.Fatalf("recordRunInSplitDb failed: %v", err)
	}

	jobs := []FailedJobItem{{StepName: "test-step", FailureSummary: "lint failed"}}
	saveParsedFailedJobs(pipeDb, slug, runItem, jobs, "raw log output")

	if !pipeDb.HasErrorLog(991001) {
		t.Errorf("expected error log for run 991001 to exist")
	}
}
