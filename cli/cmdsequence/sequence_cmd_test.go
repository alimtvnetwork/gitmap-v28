package cmdsequence

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
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
