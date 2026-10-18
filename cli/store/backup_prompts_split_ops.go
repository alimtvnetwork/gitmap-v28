// Package store — backup_prompts_split_ops.go handles mutation operations for running prompts Split-DB.
package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// InsertBackupBatch records a new backup batch and its associated prompt items.
func (db *BackupPromptsSplitDB) InsertBackupBatch(summary PromptBackupSummary) error {
	tx, err := db.conn.Begin()
	isBeginErr := err != nil
	if isBeginErr {
		return apperror.WrapSimple(err, "begin insert batch tx")
	}
	defer tx.Rollback()
	if err := insertBatchSummaryRow(tx, summary); err != nil {
		return err
	}
	if err := insertBatchItemRows(tx, summary.BatchId, summary.Items); err != nil {
		return err
	}
	return commitTx(tx, "insert backup batch")
}

func commitTx(tx *sql.Tx, action string) error {
	if err := tx.Commit(); err != nil {
		return apperror.WrapSimple(err, "commit "+action)
	}
	return nil
}

func insertBatchSummaryRow(tx *sql.Tx, summary PromptBackupSummary) error {
	q := `INSERT INTO PromptBackupBatch (BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt, Note)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	cAt := summary.CreatedAt
	isAtEmpty := len(cAt) == 0
	if isAtEmpty {
		cAt = time.Now().UTC().Format(time.RFC3339)
	}
	note := strings.Join(summary.ProjectNames, ",")
	_, err := tx.Exec(q, summary.BatchId, summary.DatabasePath, summary.TotalPrompts, summary.RunningCount, summary.EnqueuedCount, cAt, note)
	isExecErr := err != nil
	if isExecErr {
		return apperror.WrapSimple(err, "insert batch summary")
	}
	return nil
}

func insertBatchItemRows(tx *sql.Tx, batchID string, items []RunningPromptRecord) error {
	for idx, item := range items {
		if err := insertSingleBatchItemRow(tx, batchID, idx+1, item); err != nil {
			return err
		}
	}
	return nil
}

func insertSingleBatchItemRow(tx *sql.Tx, batchID string, seq int, item RunningPromptRecord) error {
	q := `INSERT INTO PromptBackupItem (ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	itemID := fmt.Sprintf("item-%s-%d", batchID, seq)
	cAt := resolveItemCreatedAt(item.CreatedAt)
	mediaJSON := marshalMediaPaths(item.MediaPaths)
	_, err := tx.Exec(q, itemID, batchID, item.ProjectName, item.ProjectPath, item.ProjectId, item.ConversationId, item.SequenceId, item.Prompt, string(item.Status), item.WordCount, mediaJSON, cAt)
	isExecErr := err != nil
	if isExecErr {
		return apperror.WrapSimple(err, "insert batch item")
	}
	return nil
}

func resolveItemCreatedAt(createdAt string) string {
	isCreatedEmpty := len(createdAt) == 0
	if isCreatedEmpty {
		return time.Now().UTC().Format(time.RFC3339)
	}
	return createdAt
}

func marshalMediaPaths(paths []string) string {
	isPathsEmpty := len(paths) == 0
	if isPathsEmpty {
		return "[]"
	}
	raw, err := json.Marshal(paths)
	isMarshalErr := err != nil
	if isMarshalErr {
		return "[]"
	}
	return string(raw)
}

// MarkBatchRestored registers a restore event in PromptRestoreLedger.
func (db *BackupPromptsSplitDB) MarkBatchRestored(batchId string, isKeep bool, count int) error {
	restoreID := fmt.Sprintf("r-%d", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC3339)
	isKeptInt := 0
	if isKeep {
		isKeptInt = 1
	}
	q := `INSERT INTO PromptRestoreLedger (RestoreId, BatchId, RestoredAt, RestoredCount, IsKept, PrunedAt)
	VALUES (?, ?, ?, ?, ?, '')`
	_, err := db.conn.Exec(q, restoreID, batchId, now, count, isKeptInt)
	isExecErr := err != nil
	if isExecErr {
		return apperror.WrapSimple(err, "mark batch restored")
	}
	return nil
}

// PruneOldRestoredEntries removes expired restore records or cleans all data when isForce is set.
func (db *BackupPromptsSplitDB) PruneOldRestoredEntries(retentionDuration time.Duration, isForce bool) (int, error) {
	if isForce {
		return db.pruneForce()
	}
	return db.pruneByRetention(retentionDuration)
}

func (db *BackupPromptsSplitDB) pruneForce() (int, error) {
	var count int
	_ = db.conn.QueryRow("SELECT COUNT(*) FROM PromptBackupBatch").Scan(&count)
	if err := db.execPruneForceStmts(); err != nil {
		return 0, err
	}
	hasPruned := count > 0
	if hasPruned {
		fmt.Println("old data has been removed")
	}
	return count, nil
}

func (db *BackupPromptsSplitDB) execPruneForceStmts() error {
	stmts := []string{
		"DELETE FROM PromptBackupItem",
		"DELETE FROM PromptRestoreLedger",
		"DELETE FROM PromptBackupBatch",
	}
	for _, stmt := range stmts {
		if _, err := db.conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "prune force exec")
		}
	}
	return nil
}

func (db *BackupPromptsSplitDB) pruneByRetention(retentionDuration time.Duration) (int, error) {
	cutoff := resolvePruneCutoff(retentionDuration)
	batchIDs, err := db.queryExpiredRestoreBatches(cutoff)
	isQueryErr := err != nil
	if isQueryErr {
		return 0, err
	}
	return db.pruneBatchList(batchIDs)
}

func resolvePruneCutoff(retentionDuration time.Duration) time.Time {
	isDefaultDuration := retentionDuration <= 0
	if isDefaultDuration {
		retentionDuration = 24 * time.Hour
	}
	return time.Now().UTC().Add(-retentionDuration)
}

func (db *BackupPromptsSplitDB) pruneBatchList(batchIDs []string) (int, error) {
	count, delErr := db.deleteBatches(batchIDs)
	isDelErr := delErr != nil
	if isDelErr {
		return 0, delErr
	}
	hasDeleted := count > 0
	if hasDeleted {
		fmt.Println("old data has been removed")
	}
	return count, nil
}

func (db *BackupPromptsSplitDB) queryExpiredRestoreBatches(cutoff time.Time) ([]string, error) {
	q := "SELECT BatchId, RestoredAt FROM PromptRestoreLedger WHERE IsKept = 0"
	rows, err := db.conn.Query(q)
	isQueryErr := err != nil
	if isQueryErr {
		return nil, apperror.WrapSimple(err, "query expired restore batches")
	}
	defer rows.Close()
	return filterExpiredBatches(rows, cutoff)
}

func filterExpiredBatches(rows *sql.Rows, cutoff time.Time) ([]string, error) {
	var batchIDs []string
	for rows.Next() {
		var bID, rAt string
		if err := rows.Scan(&bID, &rAt); err != nil {
			return nil, apperror.WrapSimple(err, "scan expired batch row")
		}
		isExpired := isTimestampExpired(rAt, cutoff)
		if isExpired {
			batchIDs = append(batchIDs, bID)
		}
	}
	return batchIDs, rows.Err()
}

func isTimestampExpired(ts string, cutoff time.Time) bool {
	t, err := time.Parse(time.RFC3339, ts)
	isParseErr := err != nil
	if isParseErr {
		t, err = time.Parse("2006-01-02 15:04:05", ts)
	}
	isFinalErr := err != nil
	if isFinalErr {
		return false
	}
	return t.Before(cutoff)
}

func (db *BackupPromptsSplitDB) deleteBatches(batchIDs []string) (int, error) {
	delCount := 0
	for _, bID := range batchIDs {
		if err := db.deleteSingleBatch(bID); err != nil {
			return delCount, err
		}
		delCount++
	}
	return delCount, nil
}

func (db *BackupPromptsSplitDB) deleteSingleBatch(bID string) error {
	stmts := []string{
		"DELETE FROM PromptBackupItem WHERE BatchId = ?",
		"DELETE FROM PromptRestoreLedger WHERE BatchId = ?",
		"DELETE FROM PromptBackupBatch WHERE BatchId = ?",
	}
	for _, stmt := range stmts {
		if _, err := db.conn.Exec(stmt, bID); err != nil {
			return apperror.WrapSimple(err, "delete single batch")
		}
	}
	return nil
}
