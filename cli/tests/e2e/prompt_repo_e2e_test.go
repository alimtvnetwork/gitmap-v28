// Package e2e — prompt_repo_e2e_test.go validates PromptBackupRepository and prompt management.
package e2e

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestPromptBackupRepository_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_backup.db")

	repo, err := store.NewPromptBackupRepositoryWithPath(dbPath)
	isInitErr := err != nil
	if isInitErr {
		t.Fatalf("failed to init repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()
	now := time.Now().UTC().Format(time.RFC3339)

	testBatch := &store.PromptBackupBatch{
		BatchID:       "b-test-01",
		SourcePath:    dbPath,
		TotalPrompts:  2,
		RunningCount:  1,
		EnqueuedCount: 1,
		CreatedAt:     now,
		Note:          "e2e test batch",
		Items: []store.PromptBackupItem{
			{
				ItemID:         "item-01",
				BatchID:        "b-test-01",
				ProjectName:    "gitmap",
				ProjectPath:    tempDir,
				ProjectID:      "p-01",
				ConversationID: "conv-01",
				SequenceID:     "seq-01",
				PromptText:     "Analyze repository patterns in Go",
				PromptStatus:   "running",
				WordCount:      6,
				CreatedAt:      now,
			},
			{
				ItemID:         "item-02",
				BatchID:        "b-test-01",
				ProjectName:    "gitmap",
				ProjectPath:    tempDir,
				ProjectID:      "p-01",
				ConversationID: "conv-02",
				SequenceID:     "seq-02",
				PromptText:     "Verify SQL parameterization",
				PromptStatus:   "queued",
				WordCount:      3,
				CreatedAt:      now,
			},
		},
	}

	created, cErr := repo.CreateBatch(ctx, testBatch)
	isCreateErr := cErr != nil
	if isCreateErr {
		t.Fatalf("CreateBatch failed: %v", cErr)
	}
	if created.BatchID != "b-test-01" {
		t.Fatalf("expected batchID b-test-01, got %s", created.BatchID)
	}

	batches, lErr := repo.ListBatches(ctx)
	isListErr := lErr != nil
	if isListErr {
		t.Fatalf("ListBatches failed: %v", lErr)
	}
	isBatchCountValid := len(batches) == 1
	if !isBatchCountValid {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}

	fetched, gErr := repo.GetBatch(ctx, "b-test-01")
	isGetErr := gErr != nil
	if isGetErr {
		t.Fatalf("GetBatch failed: %v", gErr)
	}
	if fetched.TotalPrompts != 2 {
		t.Fatalf("expected 2 prompts, got %d", fetched.TotalPrompts)
	}

	items, iErr := repo.GetBatchItems(ctx, "b-test-01")
	isItemsErr := iErr != nil
	if isItemsErr {
		t.Fatalf("GetBatchItems failed: %v", iErr)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	_, notFoundErr := repo.GetBatch(ctx, "non-existent-batch")
	isNotFound := notFoundErr != nil && notFoundErr.HasError()
	if !isNotFound {
		t.Fatalf("expected not found error on missing batch")
	}

	ledger := &store.PromptRestoreLedger{
		RestoreID:     "rest-01",
		BatchID:       "b-test-01",
		RestoredAt:    now,
		RestoredCount: 2,
		IsKept:        true,
	}
	rErr := repo.RecordRestore(ctx, ledger)
	isRestoreErr := rErr != nil
	if isRestoreErr {
		t.Fatalf("RecordRestore failed: %v", rErr)
	}
}

func TestPromptBackup_CopyAndSQLExport(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_copy.db")

	db, err := store.OpenBackupPromptsSplitDB(dbPath)
	isDbErr := err != nil
	if isDbErr {
		t.Fatalf("failed to open split db: %v", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	initialSummary := store.PromptBackupSummary{
		BatchId:       "b-copy-source",
		TotalPrompts:  1,
		RunningCount:  1,
		EnqueuedCount: 0,
		ProjectNames:  []string{"gitmap"},
		DatabasePath:  dbPath,
		CreatedAt:     now,
		Items: []store.RunningPromptRecord{
			{
				ProjectName: "gitmap",
				ProjectPath: tempDir,
				ProjectId:   "proj-copy",
				Prompt:      "Refactor SQL statements to repository pattern",
				Status:      store.PromptStatusRunning,
				WordCount:   7,
				CreatedAt:   now,
			},
		},
	}
	saveErr := db.InsertBackupBatch(initialSummary)
	if saveErr != nil {
		t.Fatalf("failed to insert initial batch: %v", saveErr)
	}
	db.Close()

	copyOpts := cmdagy.PromptCopyOptions{
		SourceID:        "b-copy-source",
		DestinationID:   "b-copy-dest",
		DestinationName: "Cloned Batch",
		Namespace:       "test-ns",
		TargetFile:      dbPath,
		IsDryRun:        true,
	}
	dryRes, dryErr := cmdagy.ExecutePromptCopy(copyOpts)
	isDryErr := dryErr != nil
	if isDryErr {
		t.Fatalf("dry-run copy failed: %v", dryErr)
	}
	if dryRes.CopiedCount != 1 {
		t.Fatalf("expected 1 copied item in dry-run, got %d", dryRes.CopiedCount)
	}

	copyOpts.IsDryRun = false
	realRes, realErr := cmdagy.ExecutePromptCopy(copyOpts)
	isRealErr := realErr != nil
	if isRealErr {
		t.Fatalf("real copy failed: %v", realErr)
	}
	if realRes.DestinationID != "b-copy-dest" {
		t.Fatalf("expected b-copy-dest, got %s", realRes.DestinationID)
	}

	sqlText := cmdagy.GeneratePromptBackupSQL(initialSummary)
	hasDDL := strings.Contains(sqlText, "CREATE TABLE IF NOT EXISTS PromptBackupBatch")
	hasInsert := strings.Contains(sqlText, "INSERT OR REPLACE INTO PromptBackupBatch")
	if !hasDDL || !hasInsert {
		t.Fatalf("expected generated SQL to contain DDL and INSERT statements")
	}

	sqlExportPath := filepath.Join(tempDir, "exported.sql")
	wErr := os.WriteFile(sqlExportPath, []byte(sqlText), 0644)
	if wErr != nil {
		t.Fatalf("failed to write export sql file: %v", wErr)
	}

	vErr := cmdagy.ValidateExportPath(sqlExportPath)
	isValErr := vErr != nil
	if isValErr {
		t.Fatalf("valid export path was rejected: %v", vErr)
	}

	tErr := cmdagy.ValidateExportPath("../../../etc/shadow")
	isTraversalRejected := tErr != nil && tErr.HasError()
	if !isTraversalRejected {
		t.Fatalf("expected path traversal to be rejected")
	}
}
