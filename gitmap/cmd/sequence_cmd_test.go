package cmd

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	_ "modernc.org/sqlite"
)

func TestScanDirectorySequence(t *testing.T) {
	tempDir := t.TempDir()

	_ = os.WriteFile(filepath.Join(tempDir, "01-first.txt"), []byte("one"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "02-second.txt"), []byte("two"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "unsequenced.txt"), []byte("three"), 0644)

	payload, err := scanDirectorySequence(tempDir)
	if err != nil {
		t.Fatalf("scanDirectorySequence failed: %v", err)
	}

	assertScanPayloadCounts(t, payload, 3, 2)
}

func assertScanPayloadCounts(t *testing.T, payload *SequencePayload, wantTotal, wantSeq int) {
	if payload.TotalFiles != wantTotal {
		t.Errorf("expected %d total files, got %d", wantTotal, payload.TotalFiles)
	}
	if payload.SequencedFiles != wantSeq {
		t.Errorf("expected %d sequenced files, got %d", wantSeq, payload.SequencedFiles)
	}
}

func TestApplySequenceOrderingWithPin(t *testing.T) {
	tempDir := t.TempDir()
	createPinTestFiles(tempDir)

	entries, _ := os.ReadDir(tempDir)
	parsedFiles := parseSeqFiles(entries, tempDir)

	flags := SequenceFlags{
		StartNum: 1,
		PinMap:   map[string]int{"index": 1, "shared-engine": 2},
	}
	applySequenceOrdering(parsedFiles, flags)

	report := executeSequenceRenames(parsedFiles, tempDir, true)
	assertPinRenameOperations(t, report)
}

func createPinTestFiles(tempDir string) {
	_ = os.WriteFile(filepath.Join(tempDir, "index.md"), []byte("index"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "shared-engine.py"), []byte("shared"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "file-manipulator.py"), []byte("manipulator"), 0644)
}

func assertPinRenameOperations(t *testing.T, report SequenceFixReport) {
	if len(report.Operations) != 3 {
		t.Fatalf("expected 3 rename operations, got %d", len(report.Operations))
	}
	for _, op := range report.Operations {
		assertSinglePinOp(t, op)
	}
}

func assertSinglePinOp(t *testing.T, op SequenceRenameOp) {
	if op.From == "index.md" && op.To != "01-index.md" {
		t.Errorf("expected index.md -> 01-index.md, got %s", op.To)
	}
	if op.From == "shared-engine.py" && op.To != "02-shared-engine.py" {
		t.Errorf("expected shared-engine.py -> 02-shared-engine.py, got %s", op.To)
	}
	if op.From == "file-manipulator.py" && op.To != "03-file-manipulator.py" {
		t.Errorf("expected file-manipulator.py -> 03-file-manipulator.py, got %s", op.To)
	}
}

func setupTestFileSequenceDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	schema := `CREATE TABLE IF NOT EXISTS FileSequence (
		Directory TEXT, Filename TEXT, SequenceNumber INTEGER, BaseName TEXT, UpdatedAt INTEGER
	);`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create FileSequence failed: %v", err)
	}

	return db
}

func TestExecuteSaveSequenceTx(t *testing.T) {
	db := setupTestFileSequenceDB(t)
	defer db.Close()

	ctx := context.Background()
	wrapper, wrapErr := dbengine.WrapDb(db, dbengine.DbSQLite)
	if wrapErr != nil {
		t.Fatalf("wrap db failed: %v", wrapErr)
	}

	payload := buildTestSequencePayload()
	err := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeSaveSequenceTx(ctx, tx, payload)
	})
	if err != nil {
		t.Fatalf("executeSaveSequenceTx failed: %v", err)
	}

	assertSavedSequenceRows(t, db, 2)
}

func buildTestSequencePayload() *SequencePayload {
	return &SequencePayload{
		Directory:      "testdir",
		TotalFiles:     2,
		SequencedFiles: 2,
		Files: []SequenceItem{
			{Sequence: 1, Filename: "01-a.txt", BaseName: "a.txt"},
			{Sequence: 2, Filename: "02-b.txt", BaseName: "b.txt"},
		},
	}
}

func assertSavedSequenceRows(t *testing.T, db *sql.DB, wantCount int) {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM FileSequence WHERE Directory = 'testdir'").Scan(&count); err != nil {
		t.Fatalf("query FileSequence count failed: %v", err)
	}
	if count != wantCount {
		t.Errorf("expected %d rows in FileSequence, got %d", wantCount, count)
	}
}

func TestExecuteSaveSequenceTxErrorPropagation(t *testing.T) {
	db := setupTestFileSequenceDB(t)
	defer db.Close()

	ctx := context.Background()
	wrapper, wrapErr := dbengine.WrapDb(db, dbengine.DbSQLite)
	if wrapErr != nil {
		t.Fatalf("wrap db failed: %v", wrapErr)
	}

	payload := buildTestSequencePayload()
	err := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		_ = tx.Tx().Rollback()

		return executeSaveSequenceTx(ctx, tx, payload)
	})
	if err == nil {
		t.Fatalf("expected executeSaveSequenceTx to fail on closed tx, got nil")
	}
}

func TestSaveSequenceToRepoDBOutsideRepo(t *testing.T) {
	payload := buildTestSequencePayload()
	err := saveSequenceToRepoDB(payload)
	if err == nil {
		t.Fatalf("expected error outside repository, got nil")
	}
}

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
