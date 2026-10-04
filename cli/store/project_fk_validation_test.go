package store

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestSelectRepoIDByPath_SlashAgnostic(t *testing.T) {
	db := openTempDB(t)

	rec := model.ScanRecord{
		Slug:         "gitmap-v28",
		RepoName:     "gitmap",
		AbsolutePath: `D:\work\gitmap`,
	}
	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}

	id, err := db.SelectRepoIDByPath("D:/work/gitmap")
	if err != nil || id <= 0 {
		t.Fatalf("expected to find repo by forward slash path, got id=%d, err=%v", id, err)
	}

	nestedID, nestedErr := db.SelectRepoIDByPath("D:/work/gitmap/nested/project")
	if nestedErr != nil || nestedID != id {
		t.Fatalf("expected nested prefix match to find parent repo id %d, got %d, err=%v", id, nestedID, nestedErr)
	}
}

func TestRepoExists(t *testing.T) {
	db := openTempDB(t)

	if db.RepoExists(0) || db.RepoExists(-1) {
		t.Fatalf("expected RepoExists to return false for non-positive IDs")
	}

	if db.RepoExists(99999) {
		t.Fatalf("expected RepoExists to return false for missing ID")
	}

	rec := model.ScanRecord{
		Slug:         "demo-repo",
		RepoName:     "demo",
		AbsolutePath: "/work/demo",
	}
	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}

	id, err := db.SelectRepoIDByPath("/work/demo")
	if err != nil || id <= 0 {
		t.Fatalf("failed to get repo id: %v", err)
	}

	if !db.RepoExists(id) {
		t.Fatalf("expected RepoExists to return true for existing repo id %d", id)
	}
}

func TestUpsertDetectedProject_ForeignKeyValidation(t *testing.T) {
	db := openTempDB(t)

	proj := model.DetectedProject{
		RepoID:           99999,
		ProjectTypeID:    1,
		ProjectName:      "test-proj",
		AbsolutePath:     "/non/existent/path",
		RepoPath:         "/non/existent/repo",
		RelativePath:     ".",
		PrimaryIndicator: "go.mod",
	}

	err := db.UpsertDetectedProject(proj)
	if err == nil {
		t.Fatalf("expected UpsertDetectedProject to fail when RepoID and paths do not exist")
	}
}

func TestPurgeOrphanRepoReferences(t *testing.T) {
	db := openTempDB(t)

	_, _ = db.conn.Exec("PRAGMA foreign_keys = OFF")
	_, err := db.conn.Exec(`INSERT INTO DetectedProject
		(RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
		VALUES (9999, 1, 'orphan', '/orphan/abs', '/orphan/repo', '.', 'go.mod')`)
	if err != nil {
		t.Fatalf("failed to insert orphan: %v", err)
	}

	_, _ = db.conn.Exec("PRAGMA foreign_keys = ON")

	db.purgeOrphanRepoReferences()

	var count int
	_ = db.conn.QueryRow("SELECT count(*) FROM DetectedProject WHERE RepoId = 9999").Scan(&count)
	if count != 0 {
		t.Fatalf("expected orphan count to be 0 after purge, got %d", count)
	}
}
