// Package store — backup_prompts_split_ops.go handles CRUD operations for running prompts Split-DB.
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
	if err != nil {
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
	if len(cAt) == 0 {
		cAt = time.Now().UTC().Format(time.RFC3339)
	}
	note := strings.Join(summary.ProjectNames, ",")
	_, err := tx.Exec(q, summary.BatchId, summary.DatabasePath, summary.TotalPrompts, summary.RunningCount, summary.EnqueuedCount, cAt, note)
	if err != nil {
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
	if err != nil {
		return apperror.WrapSimple(err, "insert batch item")
	}
	return nil
}

func resolveItemCreatedAt(createdAt string) string {
	if len(createdAt) == 0 {
		return time.Now().UTC().Format(time.RFC3339)
	}
	return createdAt
}

func marshalMediaPaths(paths []string) string {
	if len(paths) == 0 {
		return "[]"
	}
	raw, err := json.Marshal(paths)
	if err != nil {
		return "[]"
	}
	return string(raw)
}

func unmarshalMediaPaths(raw string) []string {
	var paths []string
	if len(strings.TrimSpace(raw)) == 0 {
		return nil
	}
	_ = json.Unmarshal([]byte(raw), &paths)
	return paths
}

// ExtractUniqueProjectNames extracts unique human-readable project names from items.
func ExtractUniqueProjectNames(items []RunningPromptRecord) []string {
	seen := make(map[string]bool)
	var names []string
	for _, it := range items {
		name := strings.TrimSpace(it.ProjectName)
		if len(name) > 0 && !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

// ListBackupBatches returns all recorded backup batches with their active/restored status.
func (db *BackupPromptsSplitDB) ListBackupBatches() ([]PromptBackupBatchRecord, error) {
	q := `SELECT b.BatchId, b.SourcePath, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, b.CreatedAt, b.Note,
		CASE WHEN r.RestoreId IS NOT NULL THEN 'restored' ELSE 'active' END as Status
	FROM PromptBackupBatch b
	LEFT JOIN (SELECT BatchId, RestoreId FROM PromptRestoreLedger GROUP BY BatchId) r ON b.BatchId = r.BatchId
	ORDER BY b.CreatedAt DESC`
	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, apperror.WrapSimple(err, "list backup batches")
	}
	defer rows.Close()
	return scanBatchRecords(rows)
}

func scanBatchRecords(rows *sql.Rows) ([]PromptBackupBatchRecord, error) {
	var records []PromptBackupBatchRecord
	for rows.Next() {
		var rec PromptBackupBatchRecord
		if err := rows.Scan(&rec.BatchID, &rec.SourcePath, &rec.TotalPrompts, &rec.RunningCount, &rec.EnqueuedCount, &rec.CreatedAt, &rec.Note, &rec.Status); err != nil {
			return nil, apperror.WrapSimple(err, "scan batch row")
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// GetBackupBatch returns summary and items for a specific batch.
func (db *BackupPromptsSplitDB) GetBackupBatch(batchId string) (*PromptBackupSummary, error) {
	summary, err := db.scanBatchSummaryRow(batchId)
	if err != nil {
		return nil, err
	}
	return db.populateBatchItems(summary, batchId)
}

func (db *BackupPromptsSplitDB) scanBatchSummaryRow(batchId string) (*PromptBackupSummary, error) {
	q := `SELECT BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt FROM PromptBackupBatch WHERE BatchId = ?`
	var s PromptBackupSummary
	err := db.conn.QueryRow(q, batchId).Scan(&s.BatchId, &s.DatabasePath, &s.TotalPrompts, &s.RunningCount, &s.EnqueuedCount, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NewSimple(fmt.Sprintf("batch %q not found", batchId), "E404")
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "query backup batch")
	}
	return &s, nil
}

func (db *BackupPromptsSplitDB) populateBatchItems(summary *PromptBackupSummary, batchId string) (*PromptBackupSummary, error) {
	items, err := db.getBackupBatchItems(batchId)
	if err != nil {
		return nil, err
	}
	summary.Items = items
	summary.ProjectNames = ExtractUniqueProjectNames(items)
	_, size, _ := db.GetStorageInfo()
	summary.DatabaseSize = size
	return summary, nil
}

func (db *BackupPromptsSplitDB) getBackupBatchItems(batchId string) ([]RunningPromptRecord, error) {
	q := `SELECT ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt
	FROM PromptBackupItem WHERE BatchId = ? ORDER BY CreatedAt ASC`
	rows, err := db.conn.Query(q, batchId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query batch items")
	}
	defer rows.Close()
	return scanItemRecords(rows)
}

func scanItemRecords(rows *sql.Rows) ([]RunningPromptRecord, error) {
	var items []RunningPromptRecord
	for rows.Next() {
		var item RunningPromptRecord
		var statusStr, mediaJSON string
		if err := rows.Scan(&item.ProjectName, &item.ProjectPath, &item.ProjectId, &item.ConversationId, &item.SequenceId, &item.Prompt, &statusStr, &item.WordCount, &mediaJSON, &item.CreatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "scan item row")
		}
		item.Status = PromptStatusType(statusStr)
		item.MediaPaths = unmarshalMediaPaths(mediaJSON)
		items = append(items, item)
	}
	return items, rows.Err()
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
	if err != nil {
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
	if count > 0 {
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
	if err != nil {
		return 0, err
	}
	return db.pruneBatchList(batchIDs)
}

func resolvePruneCutoff(retentionDuration time.Duration) time.Time {
	if retentionDuration <= 0 {
		retentionDuration = 24 * time.Hour
	}
	return time.Now().UTC().Add(-retentionDuration)
}

func (db *BackupPromptsSplitDB) pruneBatchList(batchIDs []string) (int, error) {
	count, delErr := db.deleteBatches(batchIDs)
	if delErr != nil {
		return 0, delErr
	}
	if count > 0 {
		fmt.Println("old data has been removed")
	}
	return count, nil
}

func (db *BackupPromptsSplitDB) queryExpiredRestoreBatches(cutoff time.Time) ([]string, error) {
	q := "SELECT BatchId, RestoredAt FROM PromptRestoreLedger WHERE IsKept = 0"
	rows, err := db.conn.Query(q)
	if err != nil {
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
		if isTimestampExpired(rAt, cutoff) {
			batchIDs = append(batchIDs, bID)
		}
	}
	return batchIDs, rows.Err()
}

func isTimestampExpired(ts string, cutoff time.Time) bool {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", ts)
	}
	if err != nil {
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
