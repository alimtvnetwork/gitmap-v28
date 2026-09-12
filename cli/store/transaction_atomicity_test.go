package store

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func setupAtomicityTestDB(t *testing.T) *DB {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_atomicity.db")

	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	if err := db.Migrate(); err != nil {
		db.Close()
		t.Fatalf("Migrate failed: %v", err)
	}

	return db
}

func seedRepoForRelease(t *testing.T, db *DB) int64 {
	t.Helper()
	if err := db.UpsertRepos([]model.ScanRecord{
		{Slug: "release-repo", RepoName: "Release Repo", AbsolutePath: "/path/to/release-repo"},
	}); err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}

	repos, err := db.FindBySlug("release-repo")
	if err != nil || len(repos) == 0 {
		t.Fatalf("failed to find created repo: %v", err)
	}

	return repos[0].ID
}

func upsertTwoReleases(t *testing.T, db *DB, repoID int64) {
	t.Helper()
	r1 := model.ReleaseRecord{RepoID: repoID, Version: "v1.0.0", Tag: "v1.0.0", IsLatest: true}
	if err := db.UpsertRelease(r1); err != nil {
		t.Fatalf("UpsertRelease v1 failed: %v", err)
	}

	r2 := model.ReleaseRecord{RepoID: repoID, Version: "v2.0.0", Tag: "v2.0.0", IsLatest: true}
	if err := db.UpsertRelease(r2); err != nil {
		t.Fatalf("UpsertRelease v2 failed: %v", err)
	}
}

func verifyLatestReleases(t *testing.T, db *DB) {
	t.Helper()
	rel1, err := db.FindReleaseByTag("v1.0.0")
	if err != nil || rel1.IsLatest {
		t.Errorf("expected v1.0.0 IsLatest false, err: %v", err)
	}

	rel2, err := db.FindReleaseByTag("v2.0.0")
	if err != nil || !rel2.IsLatest {
		t.Errorf("expected v2.0.0 IsLatest true, err: %v", err)
	}
}

func TestUpsertRelease_TransactionAtomicity(t *testing.T) {
	db := setupAtomicityTestDB(t)
	defer db.Close()

	repoID := seedRepoForRelease(t, db)
	upsertTwoReleases(t, db, repoID)
	verifyLatestReleases(t, db)
}

func createExportDataSample() model.DatabaseExport {
	return model.DatabaseExport{
		Repos: []model.ScanRecord{
			{Slug: "import-repo-1", RepoName: "Import One", AbsolutePath: "/import/1"},
			{Slug: "import-repo-2", RepoName: "Import Two", AbsolutePath: "/import/2"},
		},
		Groups: []model.GroupExport{
			{Group: model.Group{Name: "import-group", Description: "Desc", Color: "blue"}, RepoSlugs: []string{"import-repo-1"}},
		},
		History: []model.CommandHistoryRecord{
			{Command: "scan", StartedAt: "2026-09-01T00:00:00Z"},
		},
		Bookmarks: []model.BookmarkRecord{
			{Name: "bm1", Command: "scan"},
		},
	}
}

func TestImportAll_TransactionAtomicity(t *testing.T) {
	db := setupAtomicityTestDB(t)
	defer db.Close()

	if err := db.ImportAll(createExportDataSample()); err != nil {
		t.Fatalf("ImportAll failed: %v", err)
	}

	repos, err := db.ListRepos()
	if err != nil || len(repos) != 2 {
		t.Errorf("expected 2 imported repos, err: %v", err)
	}
}

func TestRemoveScanFolder_TransactionAtomicity(t *testing.T) {
	db := setupAtomicityTestDB(t)
	defer db.Close()

	folder, err := db.EnsureScanFolder("/path/to/folder", "Work", "Notes")
	if err != nil {
		t.Fatalf("EnsureScanFolder failed: %v", err)
	}

	if _, _, err := db.RemoveScanFolderByID(folder.ID); err != nil {
		t.Fatalf("RemoveScanFolderByID failed: %v", err)
	}

	folders, err := db.ListScanFolders()
	if err != nil || len(folders) != 0 {
		t.Errorf("expected 0 scan folders after delete, err: %v", err)
	}
}

func TestUpsertRepos_TransactionAtomicity(t *testing.T) {
	db := setupAtomicityTestDB(t)
	defer db.Close()

	records := []model.ScanRecord{
		{Slug: "repo1", RepoName: "Repo One", AbsolutePath: "/path/to/repo1"},
		{Slug: "repo2", RepoName: "Repo Two", AbsolutePath: "/path/to/repo2"},
	}

	if err := db.UpsertRepos(records); err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}

	repos, err := db.ListRepos()
	if err != nil || len(repos) != 2 {
		t.Errorf("expected 2 repos, err: %v", err)
	}
}
