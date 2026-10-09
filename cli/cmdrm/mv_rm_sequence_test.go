package cmdrm

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// NOTE: these mv/rm integration tests live in cmdrm (not cmdsequence)
// because removeRepoDB is unexported here. The updateRepoInDB helper
// below is a minimal test-local equivalent of cmd's version.

func setupTestStoreDB(t *testing.T) *store.DB {
	dbPath := filepath.Join(t.TempDir(), "gitmap.db")
	db, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("open test db failed: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("migrate test db failed: %v", err)
	}

	return db
}

// updateRepoInDB updates a repo's path and name. Minimal test-local
// equivalent of cmd.updateRepoInDB (direct SQL, no transaction wrapper).
func updateRepoInDB(db *store.DB, repoID int64, newPath, newName string) error {
	cleanPath := store.NormalizeStoragePath(newPath)
	_, err := db.Conn().Exec("UPDATE Repo SET AbsolutePath = ?, RepoName = ? WHERE RepoId = ?", cleanPath, newName, repoID)
	return err
}

func TestUpdateRepoInDBAndRemoveRepoDB(t *testing.T) {
	db := setupTestStoreDB(t)
	defer db.Close()

	rec := model.ScanRecord{
		AbsolutePath: "/test/repos/old-path",
		RepoName:     "old-path",
		Slug:         "old-path",
	}

	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		t.Fatalf("upsert repo failed: %v", err)
	}

	repos, _ := db.FindByPath("/test/repos/old-path")
	if len(repos) == 0 {
		t.Fatalf("repo not found after insert")
	}

	testUpdateAndRemoveFlow(t, db, repos[0])
}

func testUpdateAndRemoveFlow(t *testing.T, db *store.DB, rec model.ScanRecord) {
	if err := updateRepoInDB(db, rec.ID, "/test/repos/new-path", "new-path"); err != nil {
		t.Fatalf("updateRepoInDB failed: %v", err)
	}

	updated, _ := db.FindByPath("/test/repos/new-path")
	if len(updated) == 0 || updated[0].RepoName != "new-path" {
		t.Fatalf("expected updated repo with new path, got %+v", updated)
	}

	if err := removeRepoDB(db, updated[0]); err != nil {
		t.Fatalf("removeRepoDB failed: %v", err)
	}

	deleted, _ := db.FindByPath("/test/repos/new-path")
	if len(deleted) != 0 {
		t.Fatalf("expected 0 repos after removeRepoDB, got %d", len(deleted))
	}
}

func TestRemoveRepoDBWithAliasAtomicity(t *testing.T) {
	db := setupTestStoreDB(t)
	defer db.Close()

	rec := model.ScanRecord{
		AbsolutePath: "/test/repos/aliased-path",
		RepoName:     "aliased-path",
		Slug:         "aliased-path",
	}

	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		t.Fatalf("upsert repo failed: %v", err)
	}

	repos, _ := db.FindByPath("/test/repos/aliased-path")
	if len(repos) == 0 {
		t.Fatalf("repo not found")
	}

	testAliasRemovalFlow(t, db, repos[0])
}

func testAliasRemovalFlow(t *testing.T, db *store.DB, rec model.ScanRecord) {
	if _, err := db.CreateAlias("my-alias", rec.ID); err != nil {
		t.Fatalf("create alias failed: %v", err)
	}

	if err := removeRepoDB(db, rec); err != nil {
		t.Fatalf("removeRepoDB failed: %v", err)
	}

	alias, _ := db.FindAliasByName("my-alias")
	if alias.ID != 0 {
		t.Fatalf("expected alias to be deleted, got %+v", alias)
	}
}
