package repodb

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	_ "modernc.org/sqlite"
)

func initSchemas(ctx context.Context, db *sql.DB) error {
	if err := InitRepoSchema(ctx, db); err != nil {
		return err
	}
	return InitRootSchema(ctx, db)
}

func setupTestRepoWrapper(t *testing.T) (*dbengine.DbWrapper, func()) {
	t.Helper()
	wrapper, appErr := dbengine.OpenDb(dbengine.DbSQLite, ":memory:")
	if appErr != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", appErr)
	}
	if err := initSchemas(context.Background(), wrapper.Conn()); err != nil {
		_ = wrapper.Close()
		t.Fatalf("failed to init schemas: %v", err)
	}
	return wrapper, func() { _ = wrapper.Close() }
}

func testInsertRepoFile(t *testing.T, ctx context.Context, repo *RepoFileDbRepo) *RepoFile {
	t.Helper()
	item := &RepoFile{
		RelativePath: "src/main.go",
		AbsolutePath: "/repo/src/main.go",
		Content:      "package main",
		IsBig:        true,
		WriteTime:    1000,
		CreatedAt:    2000,
		UpdatedAt:    3000,
	}
	if res := repo.Insert(ctx, item); res.IsFailed() {
		t.Fatalf("Insert RepoFile failed: %v", res.Err)
	}
	return item
}

func testUpdateRepoFile(t *testing.T, ctx context.Context, repo *RepoFileDbRepo, item *RepoFile) {
	t.Helper()
	item.Content = "package main // updated"
	item.UpdatedAt = 4000
	if res := repo.Update(ctx, item); res.IsFailed() {
		t.Fatalf("Update RepoFile failed: %v", res.Err)
	}
	firstRes := repo.First(ctx)
	if firstRes.IsFailed() || firstRes.Value.Content != item.Content {
		t.Fatalf("expected updated content: %v", firstRes.Value)
	}
}

func testDeleteRepoFile(t *testing.T, ctx context.Context, repo *RepoFileDbRepo, id uint64) {
	t.Helper()
	if res := repo.DeleteById(ctx, id); res.IsFailed() {
		t.Fatalf("DeleteById RepoFile failed: %v", res.Err)
	}
	countRes := repo.Count(ctx)
	if countRes.IsFailed() || countRes.Value != 0 {
		t.Fatalf("expected count 0 after delete, got %d", countRes.Value)
	}
}

func TestRepoFileDbRepo_Mutations(t *testing.T) {
	ctx := context.Background()
	wrapper, cleanup := setupTestRepoWrapper(t)
	defer cleanup()

	repo := NewRepoFileDbRepo(wrapper)
	item := testInsertRepoFile(t, ctx, repo)
	firstRes := repo.First(ctx)
	if firstRes.IsFailed() || firstRes.Value.RepoFileId == 0 {
		t.Fatalf("expected valid RepoFileId: %v", firstRes.Value)
	}
	item.RepoFileId = firstRes.Value.RepoFileId
	testUpdateRepoFile(t, ctx, repo, item)
	testDeleteRepoFile(t, ctx, repo, item.RepoFileId)
}

func testSearchCacheOps(t *testing.T, ctx context.Context, repo *SearchCacheDbRepo) {
	t.Helper()
	item := &SearchCache{Query: "test", Hits: 1, ResultJson: "[]", CreatedAt: 1, UpdatedAt: 1}
	if res := repo.Insert(ctx, item); res.IsFailed() {
		t.Fatalf("Insert SearchCache failed: %v", res.Err)
	}
	firstRes := repo.First(ctx)
	if firstRes.IsFailed() || firstRes.Value.SearchCacheId == 0 {
		t.Fatalf("expected SearchCacheId: %v", firstRes.Value)
	}
	item.SearchCacheId = firstRes.Value.SearchCacheId
	item.Hits = 5
	if res := repo.Update(ctx, item); res.IsFailed() {
		t.Fatalf("Update SearchCache failed: %v", res.Err)
	}
	if res := repo.DeleteById(ctx, item.SearchCacheId); res.IsFailed() {
		t.Fatalf("DeleteById SearchCache failed: %v", res.Err)
	}
}

func TestSearchCacheDbRepo_Mutations(t *testing.T) {
	ctx := context.Background()
	wrapper, cleanup := setupTestRepoWrapper(t)
	defer cleanup()

	repo := NewSearchCacheDbRepo(wrapper)
	testSearchCacheOps(t, ctx, repo)
}

func testFileSequenceOps(t *testing.T, ctx context.Context, repo *FileSequenceDbRepo) {
	t.Helper()
	item := &FileSequence{Directory: "dir", Filename: "a.go", SequenceNumber: 1, BaseName: "a", UpdatedAt: 1}
	if res := repo.Insert(ctx, item); res.IsFailed() {
		t.Fatalf("Insert FileSequence failed: %v", res.Err)
	}
	firstRes := repo.First(ctx)
	if firstRes.IsFailed() || firstRes.Value.FileSequenceId == 0 {
		t.Fatalf("expected FileSequenceId: %v", firstRes.Value)
	}
	item.FileSequenceId = firstRes.Value.FileSequenceId
	item.SequenceNumber = 2
	if res := repo.Update(ctx, item); res.IsFailed() {
		t.Fatalf("Update FileSequence failed: %v", res.Err)
	}
	if res := repo.DeleteById(ctx, item.FileSequenceId); res.IsFailed() {
		t.Fatalf("DeleteById FileSequence failed: %v", res.Err)
	}
}

func TestFileSequenceDbRepo_Mutations(t *testing.T) {
	ctx := context.Background()
	wrapper, cleanup := setupTestRepoWrapper(t)
	defer cleanup()

	repo := NewFileSequenceDbRepo(wrapper)
	testFileSequenceOps(t, ctx, repo)
}

func testSequenceHistoryOps(t *testing.T, ctx context.Context, repo *SequenceHistoryDbRepo) {
	t.Helper()
	item := &SequenceHistory{Directory: "dir", OperationsJson: "{}", CreatedAt: 100}
	if res := repo.Insert(ctx, item); res.IsFailed() {
		t.Fatalf("Insert SequenceHistory failed: %v", res.Err)
	}
	firstRes := repo.First(ctx)
	if firstRes.IsFailed() || firstRes.Value.SequenceHistoryId == 0 {
		t.Fatalf("expected SequenceHistoryId: %v", firstRes.Value)
	}
	item.SequenceHistoryId = firstRes.Value.SequenceHistoryId
	item.OperationsJson = `{"updated": true}`
	if res := repo.Update(ctx, item); res.IsFailed() {
		t.Fatalf("Update SequenceHistory failed: %v", res.Err)
	}
	if res := repo.DeleteById(ctx, item.SequenceHistoryId); res.IsFailed() {
		t.Fatalf("DeleteById SequenceHistory failed: %v", res.Err)
	}
}

func TestSequenceHistoryDbRepo_Mutations(t *testing.T) {
	ctx := context.Background()
	wrapper, cleanup := setupTestRepoWrapper(t)
	defer cleanup()

	repo := NewSequenceHistoryDbRepo(wrapper)
	testSequenceHistoryOps(t, ctx, repo)
}

func testRepoScanLogOps(t *testing.T, ctx context.Context, repo *RepoScanLogDbRepo) {
	t.Helper()
	item := &RepoScanLog{RepoId: 10, RepoSlug: "slug", Action: "scan", Status: "success", CreatedAt: "2026-09-10"}
	if res := repo.Insert(ctx, item); res.IsFailed() {
		t.Fatalf("Insert RepoScanLog failed: %v", res.Err)
	}
	firstRes := repo.First(ctx)
	if firstRes.IsFailed() || firstRes.Value.RepoScanLogId == 0 {
		t.Fatalf("expected RepoScanLogId: %v", firstRes.Value)
	}
	item.RepoScanLogId = firstRes.Value.RepoScanLogId
	item.Status = "completed"
	if res := repo.Update(ctx, item); res.IsFailed() {
		t.Fatalf("Update RepoScanLog failed: %v", res.Err)
	}
	if res := repo.DeleteById(ctx, item.RepoScanLogId); res.IsFailed() {
		t.Fatalf("DeleteById RepoScanLog failed: %v", res.Err)
	}
}

func TestRepoScanLogDbRepo_Mutations(t *testing.T) {
	ctx := context.Background()
	wrapper, cleanup := setupTestRepoWrapper(t)
	defer cleanup()

	repo := NewRepoScanLogDbRepo(wrapper)
	testRepoScanLogOps(t, ctx, repo)
}

func testIndexedRepoOps(t *testing.T, ctx context.Context, repo *IndexedRepoDbRepo) {
	t.Helper()
	item := &IndexedRepo{Path: "/path", Slug: "path", MigratedVersion: 1, CreatedAt: 10, UpdatedAt: 10}
	if res := repo.Insert(ctx, item); res.IsFailed() {
		t.Fatalf("Insert IndexedRepo failed: %v", res.Err)
	}
	firstRes := repo.First(ctx)
	if firstRes.IsFailed() || firstRes.Value.IndexedRepoId == 0 {
		t.Fatalf("expected IndexedRepoId: %v", firstRes.Value)
	}
	item.IndexedRepoId = firstRes.Value.IndexedRepoId
	item.MigratedVersion = 2
	if res := repo.Update(ctx, item); res.IsFailed() {
		t.Fatalf("Update IndexedRepo failed: %v", res.Err)
	}
	if res := repo.DeleteById(ctx, item.IndexedRepoId); res.IsFailed() {
		t.Fatalf("DeleteById IndexedRepo failed: %v", res.Err)
	}
}

func TestIndexedRepoDbRepo_Mutations(t *testing.T) {
	ctx := context.Background()
	wrapper, cleanup := setupTestRepoWrapper(t)
	defer cleanup()

	repo := NewIndexedRepoDbRepo(wrapper)
	testIndexedRepoOps(t, ctx, repo)
}

func testClearAndReset(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	if err := ClearRepoDB(ctx, db); err != nil {
		t.Fatalf("ClearRepoDB failed: %v", err)
	}
	if err := ResetRepoDB(ctx, db); err != nil {
		t.Fatalf("ResetRepoDB failed: %v", err)
	}
}

func TestRepoDB_Lifecycle(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "repodb_test_*")
	if err != nil {
		t.Fatalf("mkdir temp failed: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	db, err := OpenRepoDB(ctx, tmpDir, filepath.Join(tmpDir, "sample-repo"), 42)
	if err != nil {
		t.Fatalf("OpenRepoDB failed: %v", err)
	}
	defer func() { _ = db.Close() }()

	testClearAndReset(t, ctx, db)
	wrap, err := OpenRepoDbWrapper(ctx, tmpDir, filepath.Join(tmpDir, "sample-repo"), 42)
	if err != nil || wrap == nil {
		t.Fatalf("OpenRepoDbWrapper failed: %v", err)
	}
	defer func() { _ = wrap.Close() }()
}
