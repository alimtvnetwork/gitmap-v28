package cmdai

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// OpenAiAnalysisSplitDB opens or initializes the AI analysis split SQLite database.
func OpenAiAnalysisSplitDB(repoRoot string) (*sql.DB, error) {
	splitDB, err := store.OpenAiAnalysisSplitDB("", repoRoot)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.OpenAiAnalysisSplitDB")
	}

	return splitDB.Conn(), nil
}

// OpenAiAnalysisSplitDBAt opens or initializes the split database at an explicit file path.
func OpenAiAnalysisSplitDBAt(dbPath string) (*sql.DB, error) {
	splitDB, err := store.OpenAiAnalysisSplitDBAt(dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.OpenAiAnalysisSplitDBAt")
	}

	return splitDB.Conn(), nil
}

// CreateAiAnalysisTask inserts a new task session into the database.
func CreateAiAnalysisTask(db *sql.DB, task AiAnalysisTask) error {
	startedAt := formatTimestamp(task.StartedAt)
	createdAt := formatTimestamp(task.CreatedAt)
	var completedAt *string
	hasCompleted := task.CompletedAt != nil
	if hasCompleted {
		val := formatTimestamp(*task.CompletedAt)
		completedAt = &val
	}

	query := `INSERT INTO AiAnalysisTask (
		task_id, repo_path, task_description, model_name, status,
		started_at, completed_at, total_files_read, total_files_modified,
		total_files_deleted, reasoning_summary, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res := store.ExecWrapper(db, query,
		task.TaskId, task.RepoPath, task.TaskDescription, task.ModelName, string(task.Status),
		startedAt, completedAt, task.TotalFilesRead, task.TotalFilesModified,
		task.TotalFilesDeleted, task.ReasoningSummary, createdAt,
	)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.CreateAiAnalysisTask")
	}

	return nil
}

// RecordAiAnalysisLine logs a file observation or mutation and updates task counters.
func RecordAiAnalysisLine(db *sql.DB, line AiAnalysisLine) error {
	createdAt := formatTimestamp(line.CreatedAt)
	lineCount := line.EndLine - line.StartLine + 1
	hasPositiveLineCount := lineCount > 0
	if !hasPositiveLineCount {
		lineCount = 1
	}

	query := `INSERT INTO AiAnalysisLine (
		task_id, file_path, operation_type, start_line, end_line,
		line_count, content_snippet, reasoning, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res := store.ExecWrapper(db, query,
		line.TaskId, line.FilePath, string(line.OperationType), line.StartLine, line.EndLine,
		lineCount, line.ContentSnippet, line.Reasoning, createdAt,
	)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.RecordAiAnalysisLine")
	}

	return updateTaskFileCounters(db, line.TaskId)
}

// BatchRecordAiAnalysisLines inserts multiple line records in an atomic transaction.
func BatchRecordAiAnalysisLines(db *sql.DB, lines []AiAnalysisLine) error {
	tx, err := db.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.BatchRecordAiAnalysisLines.begin")
	}
	defer tx.Rollback()

	taskID := ""
	for _, line := range lines {
		taskID = line.TaskId
		insertErr := executeInsertLine(tx, line)
		if insertErr != nil {
			return insertErr
		}
	}

	hasTask := taskID != ""
	if !hasTask {
		return tx.Commit()
	}

	countErr := executeCounterUpdate(tx, taskID)
	if countErr != nil {
		return countErr
	}

	return tx.Commit()
}

func executeInsertLine(tx *sql.Tx, line AiAnalysisLine) error {
	createdAt := formatTimestamp(line.CreatedAt)
	lineCount := line.EndLine - line.StartLine + 1
	hasPositiveLineCount := lineCount > 0
	if !hasPositiveLineCount {
		lineCount = 1
	}

	query := `INSERT INTO AiAnalysisLine (
		task_id, file_path, operation_type, start_line, end_line,
		line_count, content_snippet, reasoning, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res := store.ExecWrapper(tx, query,
		line.TaskId, line.FilePath, string(line.OperationType), line.StartLine, line.EndLine,
		lineCount, line.ContentSnippet, line.Reasoning, createdAt,
	)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.executeInsertLine")
	}

	return nil
}

func updateTaskFileCounters(db *sql.DB, taskID string) error {
	query := `UPDATE AiAnalysisTask SET
		total_files_read = (SELECT COUNT(DISTINCT file_path) FROM AiAnalysisLine WHERE task_id = ? AND operation_type = 'read'),
		total_files_modified = (SELECT COUNT(DISTINCT file_path) FROM AiAnalysisLine WHERE task_id = ? AND operation_type = 'edit'),
		total_files_deleted = (SELECT COUNT(DISTINCT file_path) FROM AiAnalysisLine WHERE task_id = ? AND operation_type = 'delete')
	WHERE task_id = ?`
	res := store.ExecWrapper(db, query, taskID, taskID, taskID, taskID)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.updateTaskFileCounters")
	}

	return nil
}

func executeCounterUpdate(tx *sql.Tx, taskID string) error {
	query := `UPDATE AiAnalysisTask SET
		total_files_read = (SELECT COUNT(DISTINCT file_path) FROM AiAnalysisLine WHERE task_id = ? AND operation_type = 'read'),
		total_files_modified = (SELECT COUNT(DISTINCT file_path) FROM AiAnalysisLine WHERE task_id = ? AND operation_type = 'edit'),
		total_files_deleted = (SELECT COUNT(DISTINCT file_path) FROM AiAnalysisLine WHERE task_id = ? AND operation_type = 'delete')
	WHERE task_id = ?`
	res := store.ExecWrapper(tx, query, taskID, taskID, taskID, taskID)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.executeCounterUpdate")
	}

	return nil
}

// GetAiAnalysisTask fetches a task session by its unique task ID.
func GetAiAnalysisTask(db *sql.DB, taskID string) (*AiAnalysisTask, error) {
	query := `SELECT
		task_id, repo_path, task_description, model_name, status,
		started_at, completed_at, total_files_read, total_files_modified,
		total_files_deleted, reasoning_summary, created_at
	FROM AiAnalysisTask WHERE task_id = ? LIMIT 1`
	row := store.QueryRowWrapper(db, query, taskID)

	return scanAiAnalysisTask(row, taskID)
}

func scanAiAnalysisTask(row *sql.Row, taskID string) (*AiAnalysisTask, error) {
	var task AiAnalysisTask
	var startedStr, createdStr string
	var completedStr sql.NullString
	var statusStr string

	err := row.Scan(
		&task.TaskId, &task.RepoPath, &task.TaskDescription, &task.ModelName, &statusStr,
		&startedStr, &completedStr, &task.TotalFilesRead, &task.TotalFilesModified,
		&task.TotalFilesDeleted, &task.ReasoningSummary, &createdStr,
	)
	if err != nil {
		return nil, apperror.NewNotFound("cmd.ai.analysis", "E3012", fmt.Sprintf("task %q not found", taskID))
	}

	task.Status = AiTaskStatusType(statusStr)
	task.StartedAt = parseTimestamp(startedStr)
	task.CreatedAt = parseTimestamp(createdStr)
	hasCompleted := completedStr.Valid
	if hasCompleted {
		parsed := parseTimestamp(completedStr.String)
		task.CompletedAt = &parsed
	}

	return &task, nil
}

// ListAiAnalysisTasks retrieves a paginated slice of analysis tasks.
func ListAiAnalysisTasks(db *sql.DB, limit, offset int) ([]AiAnalysisTask, error) {
	hasInvalidLimit := limit <= 0
	if hasInvalidLimit {
		limit = 50
	}
	hasNegativeOffset := offset < 0
	if hasNegativeOffset {
		offset = 0
	}

	query := `SELECT
		task_id, repo_path, task_description, model_name, status,
		started_at, completed_at, total_files_read, total_files_modified,
		total_files_deleted, reasoning_summary, created_at
	FROM AiAnalysisTask ORDER BY created_at DESC LIMIT ? OFFSET ?`
	res := store.QueryWrapper(db, query, limit, offset)
	if res.IsFailure {
		return nil, apperror.WrapSimple(res.Error, "cmdai.ListAiAnalysisTasks")
	}
	defer res.Data.Close()

	return scanAiAnalysisTasks(res.Data)
}

func scanAiAnalysisTasks(rows *sql.Rows) ([]AiAnalysisTask, error) {
	var tasks []AiAnalysisTask
	for rows.Next() {
		task, err := scanSingleTaskRow(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func scanSingleTaskRow(rows *sql.Rows) (AiAnalysisTask, error) {
	var task AiAnalysisTask
	var startedStr, createdStr string
	var completedStr sql.NullString
	var statusStr string

	err := rows.Scan(
		&task.TaskId, &task.RepoPath, &task.TaskDescription, &task.ModelName, &statusStr,
		&startedStr, &completedStr, &task.TotalFilesRead, &task.TotalFilesModified,
		&task.TotalFilesDeleted, &task.ReasoningSummary, &createdStr,
	)
	if err != nil {
		return task, apperror.WrapSimple(err, "cmdai.scanSingleTaskRow")
	}

	task.Status = AiTaskStatusType(statusStr)
	task.StartedAt = parseTimestamp(startedStr)
	task.CreatedAt = parseTimestamp(createdStr)
	hasCompleted := completedStr.Valid
	if hasCompleted {
		parsed := parseTimestamp(completedStr.String)
		task.CompletedAt = &parsed
	}

	return task, nil
}

// GetAiAnalysisLinesForTask retrieves all logged lines for a specific task.
func GetAiAnalysisLinesForTask(db *sql.DB, taskID string) ([]AiAnalysisLine, error) {
	query := `SELECT
		id, task_id, file_path, operation_type, start_line, end_line,
		line_count, content_snippet, reasoning, created_at
	FROM AiAnalysisLine WHERE task_id = ? ORDER BY id ASC`
	res := store.QueryWrapper(db, query, taskID)
	if res.IsFailure {
		return nil, apperror.WrapSimple(res.Error, "cmdai.GetAiAnalysisLinesForTask")
	}
	defer res.Data.Close()

	return scanAiAnalysisLines(res.Data)
}

func scanAiAnalysisLines(rows *sql.Rows) ([]AiAnalysisLine, error) {
	var lines []AiAnalysisLine
	for rows.Next() {
		line, err := scanSingleLineRow(rows)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func scanSingleLineRow(rows *sql.Rows) (AiAnalysisLine, error) {
	var line AiAnalysisLine
	var opStr, createdStr string

	err := rows.Scan(
		&line.Id, &line.TaskId, &line.FilePath, &opStr, &line.StartLine, &line.EndLine,
		&line.LineCount, &line.ContentSnippet, &line.Reasoning, &createdStr,
	)
	if err != nil {
		return line, apperror.WrapSimple(err, "cmdai.scanSingleLineRow")
	}

	line.OperationType = AiOperationType(opStr)
	line.CreatedAt = parseTimestamp(createdStr)

	return line, nil
}

// UpdateAiAnalysisTaskStatus modifies status, summary, and completed timestamp.
func UpdateAiAnalysisTaskStatus(db *sql.DB, taskID string, status AiTaskStatusType, summary string) error {
	var completedAt *string
	isFinished := status == AiTaskStatusCompleted || status == AiTaskStatusFailed || status == AiTaskStatusCancelled
	if isFinished {
		now := formatTimestamp(time.Now())
		completedAt = &now
	}

	query := `UPDATE AiAnalysisTask SET
		status = ?,
		reasoning_summary = CASE WHEN ? != '' THEN ? ELSE reasoning_summary END,
		completed_at = COALESCE(?, completed_at)
	WHERE task_id = ?`
	res := store.ExecWrapper(db, query, string(status), summary, summary, completedAt, taskID)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.UpdateAiAnalysisTaskStatus")
	}

	return checkRowsAffected(res.Data, taskID)
}

func checkRowsAffected(res sql.Result, taskID string) error {
	affected, err := res.RowsAffected()
	hasZeroAffected := err == nil && affected == 0
	if hasZeroAffected {
		return apperror.NewNotFound("cmd.ai.analysis", "E3012", fmt.Sprintf("task %q not found", taskID))
	}

	return nil
}

func formatTimestamp(t time.Time) string {
	isZero := t.IsZero()
	if isZero {
		return time.Now().UTC().Format(time.RFC3339)
	}

	return t.UTC().Format(time.RFC3339)
}

func parseTimestamp(raw string) time.Time {
	clean := raw
	isEmpty := clean == ""
	if isEmpty {
		return time.Time{}
	}

	parsed, err := time.Parse(time.RFC3339, clean)
	hasParsed := err == nil
	if hasParsed {
		return parsed
	}

	fallback, err2 := time.Parse("2006-01-02 15:04:05", clean)
	hasFallback := err2 == nil
	if hasFallback {
		return fallback
	}

	return time.Time{}
}
