package cmd

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TestUpdateRepoInDB_Success verifies the mv database update.
// NOTE: this test lives in cmd (not cmdrm) because updateRepoInDB is
// defined in cmd/mv_db.go. The rm transaction tests live in
// cmdrm/mv_rm_tx_test.go. Helpers are duplicated from there.

func TestUpdateRepoInDB_Success(t *testing.T) {
	db := setupMvTestDB(t)
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

func setupMvTestDB(t *testing.T) *store.DB {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_mv.db")

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

func assertRepoUpdated(t *testing.T, db *store.DB, repoID int64, expectedPath, expectedName string) {
	t.Helper()
	var actualPath, actualName string
	err := db.Conn().QueryRow("SELECT AbsolutePath, RepoName FROM Repo WHERE RepoId = ?", repoID).Scan(&actualPath, &actualName)
	if err != nil {
		t.Fatalf("failed to query updated repo: %v", err)
	}

	normalizedExpected := store.NormalizeStoragePath(expectedPath)
	isPathMatch := actualPath == expectedPath || actualPath == normalizedExpected
	if !isPathMatch || actualName != expectedName {
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
