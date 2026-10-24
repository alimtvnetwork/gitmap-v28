// Package cmdagy — agy_running_prompts_export_sql.go generates idempotent SQL reconstitution scripts and manages exports.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const sqlExportPromptDDL = `CREATE TABLE IF NOT EXISTS PromptBackupBatch (BatchId TEXT PRIMARY KEY, SourcePath TEXT NOT NULL DEFAULT '', TotalPrompts INTEGER NOT NULL DEFAULT 0, RunningCount INTEGER NOT NULL DEFAULT 0, EnqueuedCount INTEGER NOT NULL DEFAULT 0, CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, Note TEXT NOT NULL DEFAULT '');
CREATE TABLE IF NOT EXISTS PromptBackupItem (ItemId TEXT PRIMARY KEY, BatchId TEXT NOT NULL, ProjectName TEXT NOT NULL, ProjectPath TEXT NOT NULL, ProjectId TEXT NOT NULL, ConversationId TEXT NOT NULL, SequenceId TEXT NOT NULL DEFAULT '', PromptText TEXT NOT NULL, PromptStatus TEXT NOT NULL, WordCount INTEGER NOT NULL DEFAULT 0, MediaPathsJson TEXT NOT NULL DEFAULT '[]', CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE);
CREATE INDEX IF NOT EXISTS idx_prompt_backup_item_batch ON PromptBackupItem(BatchId);
CREATE INDEX IF NOT EXISTS idx_prompt_backup_item_project ON PromptBackupItem(ProjectId);
CREATE TABLE IF NOT EXISTS PromptRestoreLedger (RestoreId TEXT PRIMARY KEY, BatchId TEXT NOT NULL, RestoredAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, RestoredCount INTEGER NOT NULL DEFAULT 0, IsKept INTEGER NOT NULL DEFAULT 0, PrunedAt TEXT NOT NULL DEFAULT '', FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE);
CREATE INDEX IF NOT EXISTS idx_prompt_restore_batch ON PromptRestoreLedger(BatchId);`

// ValidateExportPath checks that the target file path does not attempt path traversal.
func ValidateExportPath(targetPath string) *appfault.AppError {
	trimmed := strings.TrimSpace(targetPath)
	isPathEmpty := len(trimmed) == 0
	if isPathEmpty {
		return appfault.NewValidation("ValidateExportPath", "E1000", "export target path cannot be empty")
	}

	slashPath := filepath.ToSlash(filepath.Clean(trimmed))
	isTraversal := strings.HasPrefix(slashPath, "../") || slashPath == ".." || strings.Contains(slashPath, "/../")
	if isTraversal {
		return appfault.NewValidation("ValidateExportPath", "E_PATH_TRAVERSAL", "path traversal detected: "+targetPath)
	}

	return nil
}

// GeneratePromptBackupSQL generates an idempotent SQL script to reconstitute prompt backup data in SQLite.
func GeneratePromptBackupSQL(summary store.PromptBackupSummary) string {
	return GenerateMultiBatchSQL([]store.PromptBackupSummary{summary})
}

// GenerateMultiBatchSQL generates idempotent SQL statements for a collection of backup batches.
func GenerateMultiBatchSQL(summaries []store.PromptBackupSummary) string {
	var sb strings.Builder
	sb.WriteString("-- GitMap AGY Running Prompts Multi-Batch SQL Backup Export\nPRAGMA foreign_keys = ON;\nBEGIN TRANSACTION;\n\n")
	sb.WriteString(sqlExportPromptDDL + "\n")

	for _, summary := range summaries {
		sb.WriteString("\n-- Batch: " + summary.BatchId + "\n" + formatBatchInsertSQL(summary))
		for idx, item := range summary.Items {
			sb.WriteString(formatItemInsertSQL(summary.BatchId, idx+1, item))
		}
	}

	sb.WriteString("\nCOMMIT;\n")

	return sb.String()
}

func formatBatchInsertSQL(summary store.PromptBackupSummary) string {
	cAt := resolveExportCreatedAt(summary.CreatedAt)
	note := escapeSQLString(strings.Join(summary.ProjectNames, ","))
	sourcePath := escapeSQLString(summary.DatabasePath)
	batchID := escapeSQLString(summary.BatchId)
	q := "INSERT OR REPLACE INTO PromptBackupBatch (BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt, Note) VALUES ('%s', '%s', %d, %d, %d, '%s', '%s');\n"

	return fmt.Sprintf(q, batchID, sourcePath, summary.TotalPrompts, summary.RunningCount, summary.EnqueuedCount, cAt, note)
}

func formatItemInsertSQL(batchID string, seq int, item store.RunningPromptRecord) string {
	itemID := fmt.Sprintf("item-%s-%d", batchID, seq)
	cAt := resolveExportCreatedAt(item.CreatedAt)
	mediaJSON := marshalItemMediaJSON(item.MediaPaths)
	q := "INSERT OR REPLACE INTO PromptBackupItem (ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, '%s', '%s');\n"

	return fmt.Sprintf(q, escapeSQLString(itemID), escapeSQLString(batchID), escapeSQLString(item.ProjectName), escapeSQLString(item.ProjectPath), escapeSQLString(item.ProjectId), escapeSQLString(item.ConversationId), escapeSQLString(item.SequenceId), escapeSQLString(item.Prompt), escapeSQLString(string(item.Status)), item.WordCount, escapeSQLString(mediaJSON), escapeSQLString(cAt))
}

func escapeSQLString(val string) string {
	return strings.ReplaceAll(val, "'", "''")
}

func marshalItemMediaJSON(paths []string) string {
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

func resolveExportCreatedAt(raw string) string {
	isRawEmpty := len(strings.TrimSpace(raw)) == 0
	if isRawEmpty {
		return time.Now().UTC().Format(time.RFC3339)
	}

	return raw
}

// ExportPromptBackupToSQL generates SQL and writes it to the designated target file or stdout.
func ExportPromptBackupToSQL(req PromptExportRequest) (*PromptExportResult, *appfault.AppError) {
	pathErr := validateExportFileTarget(req.TargetFile, req.IsStdout)
	isPathErr := pathErr != nil
	if isPathErr {
		return nil, pathErr
	}

	summary, err := loadExportSummary(req.BatchID, req.MaxWords, req.IsFull)
	isLoadErr := err != nil
	if isLoadErr {
		return nil, err
	}

	sqlContent := GeneratePromptBackupSQL(*summary)

	isStdout := req.IsStdout
	if isStdout {
		fmt.Print(sqlContent)
		return buildExportResult(req.TargetFile, enums.PromptExportFormatTypeSQL, summary.BatchId, summary.TotalPrompts, sqlContent), nil
	}

	writeErr := writeExportFile(req.TargetFile, []byte(sqlContent))
	isWriteErr := writeErr != nil
	if isWriteErr {
		return nil, writeErr
	}

	return buildExportResult(req.TargetFile, enums.PromptExportFormatTypeSQL, summary.BatchId, summary.TotalPrompts, sqlContent), nil
}

func validateExportFileTarget(targetFile string, isStdout bool) *appfault.AppError {
	if isStdout {
		return nil
	}

	return ValidateExportPath(targetFile)
}

func writeExportFile(filePath string, data []byte) *appfault.AppError {
	_ = os.MkdirAll(filepath.Dir(filePath), 0755)

	wErr := os.WriteFile(filePath, data, 0644)
	isWErr := wErr != nil
	if isWErr {
		return appfault.WrapSimple(wErr, "write export file")
	}

	return nil
}

// ExportPromptsByFormat exports prompts in JSON, SQL, or DB format.
func ExportPromptsByFormat(req PromptExportRequest) (*PromptExportResult, *appfault.AppError) {
	pathErr := validateExportFileTarget(req.TargetFile, req.IsStdout)
	isPathErr := pathErr != nil
	if isPathErr {
		return nil, pathErr
	}

	isSQL := req.Format.IsSQL()
	if isSQL {
		return ExportPromptBackupToSQL(req)
	}

	isJSON := req.Format.IsJSON()
	if isJSON {
		return exportPromptsToJSON(req)
	}

	isDB := req.Format.IsDB()
	if isDB {
		return exportPromptsToDB(req)
	}

	return nil, appfault.NewValidation("ExportPromptsByFormat", "E1000", fmt.Sprintf("unsupported export format: %s", req.Format))
}

func exportPromptsToJSON(req PromptExportRequest) (*PromptExportResult, *appfault.AppError) {
	summary, err := loadExportSummary(req.BatchID, req.MaxWords, req.IsFull)
	isLoadErr := err != nil
	if isLoadErr {
		return nil, err
	}

	data, mErr := json.MarshalIndent(summary, "", "  ")
	isMarshalErr := mErr != nil
	if isMarshalErr {
		return nil, appfault.WrapSimple(mErr, "marshal export json")
	}

	isStdout := req.IsStdout
	if isStdout {
		fmt.Println(string(data))
		return buildExportResult(req.TargetFile, enums.PromptExportFormatTypeJSON, summary.BatchId, summary.TotalPrompts, string(data)), nil
	}

	writeErr := writeExportFile(req.TargetFile, data)
	isWriteErr := writeErr != nil
	if isWriteErr {
		return nil, writeErr
	}

	return buildExportResult(req.TargetFile, enums.PromptExportFormatTypeJSON, summary.BatchId, summary.TotalPrompts, string(data)), nil
}

func exportPromptsToDB(req PromptExportRequest) (*PromptExportResult, *appfault.AppError) {
	summary, err := loadExportSummary(req.BatchID, req.MaxWords, req.IsFull)
	isLoadErr := err != nil
	if isLoadErr {
		return nil, err
	}

	dbErr := exportToSQLiteFile(req.TargetFile, summary.Items)
	isDbErr := dbErr != nil
	if isDbErr {
		return nil, appfault.WrapSimple(dbErr, "export to sqlite db")
	}

	return buildExportResult(req.TargetFile, enums.PromptExportFormatTypeDB, summary.BatchId, summary.TotalPrompts, ""), nil
}

func loadExportSummary(batchID string, maxWords int, isFull bool) (*store.PromptBackupSummary, *appfault.AppError) {
	cleanBatchID := strings.TrimSpace(batchID)
	hasBatchID := len(cleanBatchID) > 0

	db, err := store.OpenBackupPromptsSplitDB("")
	isDbOpenOk := err == nil
	if isDbOpenOk {
		defer db.Close()
	}

	summary := resolveDBSummary(db, cleanBatchID, hasBatchID, isDbOpenOk)
	isSummaryFound := summary != nil
	if isSummaryFound {
		return summary, nil
	}

	return buildAdHocExportSummary(maxWords, isFull)
}

func resolveDBSummary(db *store.BackupPromptsSplitDB, batchID string, hasBatchID, isDbOpenOk bool) *store.PromptBackupSummary {
	if !isDbOpenOk {
		return nil
	}

	if hasBatchID {
		summary, _ := db.GetBackupBatch(batchID)
		return summary
	}

	batches, err := db.ListBackupBatches()
	hasBatches := err == nil && len(batches) > 0
	if !hasBatches {
		return nil
	}

	summary, _ := db.GetBackupBatch(batches[0].BatchID)

	return summary
}

func buildAdHocExportSummary(maxWords int, isFull bool) (*store.PromptBackupSummary, *appfault.AppError) {
	items, cErr := CollectActiveAndQueuedPrompts(maxWords, isFull)
	isCollectErr := cErr != nil
	if isCollectErr {
		return nil, appfault.WrapSimple(cErr, "collect active prompts for export")
	}

	rCount, eCount := countPromptStatuses(items)
	adhocBatchID := fmt.Sprintf("adhoc-%x", time.Now().UnixNano()%0xffffffff)

	return &store.PromptBackupSummary{
		BatchId: adhocBatchID, TotalPrompts: len(items),
		RunningCount: rCount, EnqueuedCount: eCount,
		ProjectNames: extractCleanProjectNames(items),
		CreatedAt:    time.Now().UTC().Format(time.RFC3339), Items: items,
	}, nil
}

func buildExportResult(targetFile string, format enums.PromptExportFormatType, batchID string, total int, content string) *PromptExportResult {
	return &PromptExportResult{
		TargetFile: targetFile, Format: format, BatchID: batchID,
		TotalPrompts: total, TotalBytes: int64(len(content)),
		IsSuccess: true, SqlContent: content,
	}
}
