// Package store — backup_prompts_split_query.go handles query and inspection for running prompts Split-DB.
package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ExtractUniqueProjectNames extracts unique human-readable project names from items.
func ExtractUniqueProjectNames(items []RunningPromptRecord) []string {
	seen := make(map[string]bool)
	var names []string
	for _, it := range items {
		name := strings.TrimSpace(it.ProjectName)
		hasName := len(name) > 0 && !seen[name]
		if hasName {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

func unmarshalMediaPaths(raw string) []string {
	var paths []string
	isRawEmpty := len(strings.TrimSpace(raw)) == 0
	if isRawEmpty {
		return nil
	}
	_ = json.Unmarshal([]byte(raw), &paths)
	return paths
}

// ListBackupBatches returns all recorded backup batches with their active/restored status.
func (db *BackupPromptsSplitDB) ListBackupBatches() ([]PromptBackupBatchRecord, error) {
	q := `SELECT b.BatchId, b.SourcePath, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, b.CreatedAt, b.Note,
		CASE WHEN r.RestoreId IS NOT NULL THEN 'restored' ELSE 'active' END as Status
	FROM PromptBackupBatch b
	LEFT JOIN (SELECT BatchId, RestoreId FROM PromptRestoreLedger GROUP BY BatchId) r ON b.BatchId = r.BatchId
	ORDER BY b.CreatedAt DESC`
	rows, err := db.conn.Query(q)
	isQueryErr := err != nil
	if isQueryErr {
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
	isScanErr := err != nil
	if isScanErr {
		return nil, err
	}
	return db.populateBatchItems(summary, batchId)
}

func (db *BackupPromptsSplitDB) scanBatchSummaryRow(batchId string) (*PromptBackupSummary, error) {
	q := `SELECT BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt FROM PromptBackupBatch WHERE BatchId = ?`
	var s PromptBackupSummary
	err := db.conn.QueryRow(q, batchId).Scan(&s.BatchId, &s.DatabasePath, &s.TotalPrompts, &s.RunningCount, &s.EnqueuedCount, &s.CreatedAt)
	isNotFound := err == sql.ErrNoRows
	if isNotFound {
		return nil, apperror.NewSimple(fmt.Sprintf("batch %q not found", batchId), "E404")
	}
	isErr := err != nil
	if isErr {
		return nil, apperror.WrapSimple(err, "query backup batch")
	}
	return &s, nil
}

func (db *BackupPromptsSplitDB) populateBatchItems(summary *PromptBackupSummary, batchId string) (*PromptBackupSummary, error) {
	items, err := db.getBackupBatchItems(batchId)
	isItemsErr := err != nil
	if isItemsErr {
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
	isQueryErr := err != nil
	if isQueryErr {
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
