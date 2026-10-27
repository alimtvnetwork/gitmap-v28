package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// HeatmapDB wraps an isolated SQLite database connection for decision heatmap tracking.
type HeatmapDB struct {
	conn *sql.DB
	Path string
}

// Conn returns the underlying database connection.
func (h *HeatmapDB) Conn() *sql.DB {
	return h.conn
}

// Close closes the underlying database connection.
func (h *HeatmapDB) Close() error {
	if h.conn == nil {
		return nil
	}
	return h.conn.Close()
}

// OpenHeatmapRootDB opens or creates the root global heatmap database.
func OpenHeatmapRootDB(repoRoot string) (*HeatmapDB, error) {
	dbPath := ResolveHeatmapRootDbPath(repoRoot)
	return OpenHeatmapDBAtPath(dbPath)
}

// OpenHeatmapRepoDB opens or creates a repository-specific split heatmap database.
func OpenHeatmapRepoDB(repoSlug, repoRoot string) (*HeatmapDB, error) {
	dbPath := ResolveHeatmapRepoSplitDbPath(repoSlug, repoRoot)
	return OpenHeatmapDBAtPath(dbPath)
}

// OpenHeatmapDBAtPath opens or creates a heatmap database at the specified filesystem path.
func OpenHeatmapDBAtPath(dbPath string) (*HeatmapDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.open")
	}

	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "heatmap.config")
	}

	if err := EnsureHeatmapTables(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &HeatmapDB{conn: conn, Path: dbPath}, nil
}

// LogAction inserts an action record, computing execution depth and order from ParentLogId.
func (h *HeatmapDB) LogAction(log HeatmapActionLogRecord) (int64, error) {
	depth, order, err := h.calculateDepthAndOrder(log.ParentLogId, log.TaskSlug)
	if err != nil {
		return 0, err
	}

	isMetInt := 0
	if log.IsExpectationMet {
		isMetInt = 1
	}

	const sqlInsert = `INSERT INTO HeatmapActionLog (
		ParentLogId, TaskSlug, RepoName, ActionType, FilePath, StartLine, EndLine,
		CodeSnippet, CommandText, ReasoningState, ReasoningNotes, ExpectedOutcome,
		ActualOutcome, IsExpectationMet, ExecutionDepth, ExecutionOrder, Notes, Comments
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := h.conn.Exec(sqlInsert,
		log.ParentLogId, log.TaskSlug, log.RepoName, log.ActionType, log.FilePath,
		log.StartLine, log.EndLine, log.CodeSnippet, log.CommandText, log.ReasoningState,
		log.ReasoningNotes, log.ExpectedOutcome, log.ActualOutcome, isMetInt,
		depth, order, log.Notes, log.Comments,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "heatmap.logAction")
	}

	return res.LastInsertId()
}

func (h *HeatmapDB) calculateDepthAndOrder(parentLogID *int64, taskSlug string) (int, int, error) {
	if parentLogID == nil || *parentLogID <= 0 {
		return h.calculateRootDepthAndOrder(taskSlug)
	}

	return h.calculateChildDepthAndOrder(*parentLogID)
}

func (h *HeatmapDB) calculateRootDepthAndOrder(taskSlug string) (int, int, error) {
	var maxOrder sql.NullInt64
	query := `SELECT MAX(ExecutionOrder) FROM HeatmapActionLog WHERE TaskSlug = ? AND ParentLogId IS NULL`
	err := h.conn.QueryRow(query, taskSlug).Scan(&maxOrder)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, apperror.WrapSimple(err, "heatmap.calcRootOrder")
	}

	return 0, int(maxOrder.Int64) + 1, nil
}

func (h *HeatmapDB) calculateChildDepthAndOrder(parentLogID int64) (int, int, error) {
	var parentDepth int
	queryParent := `SELECT ExecutionDepth FROM HeatmapActionLog WHERE HeatmapActionLogId = ?`
	if err := h.conn.QueryRow(queryParent, parentLogID).Scan(&parentDepth); err != nil {
		return 0, 0, apperror.WrapSimple(err, "heatmap.parentLookup")
	}

	var maxChildOrder sql.NullInt64
	queryOrder := `SELECT MAX(ExecutionOrder) FROM HeatmapActionLog WHERE ParentLogId = ?`
	err := h.conn.QueryRow(queryOrder, parentLogID).Scan(&maxChildOrder)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, apperror.WrapSimple(err, "heatmap.calcChildOrder")
	}

	return parentDepth + 1, int(maxChildOrder.Int64) + 1, nil
}

// SearchActionLogs searches action logs by task, query string, and file path.
func (h *HeatmapDB) SearchActionLogs(taskSlug, query, filePath string, limit int) ([]HeatmapActionLogRecord, error) {
	if limit <= 0 {
		limit = 50
	}

	var whereClauses []string
	var args []any

	if strings.TrimSpace(taskSlug) != "" {
		whereClauses = append(whereClauses, "TaskSlug = ?")
		args = append(args, taskSlug)
	}
	if strings.TrimSpace(filePath) != "" {
		whereClauses = append(whereClauses, "FilePath LIKE ?")
		args = append(args, "%"+filePath+"%")
	}
	if strings.TrimSpace(query) != "" {
		q := "%" + query + "%"
		whereClauses = append(whereClauses, "(CodeSnippet LIKE ? OR CommandText LIKE ? OR ReasoningNotes LIKE ? OR ActionType LIKE ?)")
		args = append(args, q, q, q, q)
	}

	sqlText := `SELECT HeatmapActionLogId, ParentLogId, TaskSlug, RepoName, ActionType,
		FilePath, StartLine, EndLine, COALESCE(CodeSnippet, ''), COALESCE(CommandText, ''),
		ReasoningState, COALESCE(ReasoningNotes, ''), COALESCE(ExpectedOutcome, ''),
		COALESCE(ActualOutcome, ''), IsExpectationMet, ExecutionDepth, ExecutionOrder,
		COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt
		FROM HeatmapActionLog`

	if len(whereClauses) > 0 {
		sqlText += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	sqlText += " ORDER BY CreatedAt DESC LIMIT ?"
	args = append(args, limit)

	rows, err := h.conn.Query(sqlText, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.search")
	}
	defer rows.Close()

	return scanActionLogs(rows)
}

// GetExecutionTree returns hierarchical execution steps organized as a tree using recursive CTE.
func (h *HeatmapDB) GetExecutionTree(taskSlug string) ([]HeatmapExecutionTreeNode, error) {
	const sqlCTE = `WITH RECURSIVE ActionTree AS (
		SELECT HeatmapActionLogId, ParentLogId, TaskSlug, RepoName, ActionType, FilePath,
			StartLine, EndLine, CodeSnippet, CommandText, ReasoningState, ReasoningNotes,
			ExpectedOutcome, ActualOutcome, IsExpectationMet, ExecutionDepth, ExecutionOrder,
			Notes, Comments, CreatedAt,
			PRINTF('%06d', ExecutionOrder) AS TreePath
		FROM HeatmapActionLog
		WHERE (? = '' OR TaskSlug = ?) AND ParentLogId IS NULL
		UNION ALL
		SELECT c.HeatmapActionLogId, c.ParentLogId, c.TaskSlug, c.RepoName, c.ActionType, c.FilePath,
			c.StartLine, c.EndLine, c.CodeSnippet, c.CommandText, c.ReasoningState, c.ReasoningNotes,
			c.ExpectedOutcome, c.ActualOutcome, c.IsExpectationMet, c.ExecutionDepth, c.ExecutionOrder,
			c.Notes, c.Comments, c.CreatedAt,
			p.TreePath || '/' || PRINTF('%06d', c.ExecutionOrder) AS TreePath
		FROM HeatmapActionLog c
		JOIN ActionTree p ON c.ParentLogId = p.HeatmapActionLogId
	)
	SELECT HeatmapActionLogId, ParentLogId, TaskSlug, RepoName, ActionType, FilePath,
		StartLine, EndLine, COALESCE(CodeSnippet, ''), COALESCE(CommandText, ''),
		ReasoningState, COALESCE(ReasoningNotes, ''), COALESCE(ExpectedOutcome, ''),
		COALESCE(ActualOutcome, ''), IsExpectationMet, ExecutionDepth, ExecutionOrder,
		COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt
	FROM ActionTree
	ORDER BY TreePath ASC`

	rows, err := h.conn.Query(sqlCTE, taskSlug, taskSlug)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.executionTree")
	}
	defer rows.Close()

	records, err := scanActionLogs(rows)
	if err != nil {
		return nil, err
	}

	return assembleTreeNodes(records), nil
}

// GetHeatmapStatus computes statistical and touch summaries for a task.
func (h *HeatmapDB) GetHeatmapStatus(taskSlug string) (HeatmapStatusSummary, error) {
	summary := HeatmapStatusSummary{
		TaskSlug:             taskSlug,
		ReasoningStateCounts: make(map[string]int),
		TopTouchedFiles:      make([]FileTouchCount, 0),
	}

	if err := h.populateStatusCounts(&summary, taskSlug); err != nil {
		return summary, err
	}
	if err := h.populateReasoningStates(&summary, taskSlug); err != nil {
		return summary, err
	}
	if err := h.populateTopTouchedFiles(&summary, taskSlug); err != nil {
		return summary, err
	}

	return summary, nil
}

func (h *HeatmapDB) populateStatusCounts(s *HeatmapStatusSummary, taskSlug string) error {
	const sqlCounts = `SELECT COUNT(*), COALESCE(SUM(CASE WHEN IsExpectationMet = 1 THEN 1 ELSE 0 END), 0)
		FROM HeatmapActionLog WHERE (? = '' OR TaskSlug = ?)`
	err := h.conn.QueryRow(sqlCounts, taskSlug, taskSlug).Scan(&s.TotalActions, &s.ExpectationsMetCount)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.statusCounts")
	}
	if s.TotalActions > 0 {
		s.ExpectationsMetPct = (float64(s.ExpectationsMetCount) / float64(s.TotalActions)) * 100.0
	}
	return nil
}

func (h *HeatmapDB) populateReasoningStates(s *HeatmapStatusSummary, taskSlug string) error {
	const sqlStates = `SELECT ReasoningState, COUNT(*) FROM HeatmapActionLog
		WHERE (? = '' OR TaskSlug = ?) GROUP BY ReasoningState`
	rows, err := h.conn.Query(sqlStates, taskSlug, taskSlug)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.statusStates")
	}
	defer rows.Close()

	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err == nil {
			s.ReasoningStateCounts[state] = count
		}
	}
	return nil
}

const sqlTopTouchedFiles = `SELECT FilePath, COUNT(*) AS c,
	COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('READ', 'READ_FILE', 'INSPECT') THEN 1 ELSE 0 END), 0) AS reads,
	COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('EDIT', 'EDIT_FILE', 'MODIFY') THEN 1 ELSE 0 END), 0) AS edits,
	COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('CREATE', 'WRITE_FILE', 'NEW') THEN 1 ELSE 0 END), 0) AS creates
	FROM HeatmapActionLog WHERE (? = '' OR TaskSlug = ?) AND FilePath != ''
	GROUP BY FilePath ORDER BY c DESC LIMIT 10`

func (h *HeatmapDB) populateTopTouchedFiles(s *HeatmapStatusSummary, taskSlug string) error {
	rows, err := h.conn.Query(sqlTopTouchedFiles, taskSlug, taskSlug)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.statusTopFiles")
	}

	defer rows.Close()

	return scanTopTouchedRows(rows, s)
}

func scanTopTouchedRows(rows *sql.Rows, s *HeatmapStatusSummary) error {
	for rows.Next() {
		var item FileTouchCount
		err := rows.Scan(&item.FilePath, &item.Count, &item.ReadCount, &item.EditCount, &item.CreatedCount)
		if err != nil {
			return apperror.WrapSimple(err, "heatmap.scanTopFiles")
		}

		s.TopTouchedFiles = append(s.TopTouchedFiles, item)
	}

	return nil
}

const sqlPurgeStaleActionLogs = `DELETE FROM HeatmapActionLog 
WHERE CreatedAt < datetime('now', printf('-%d days', ?))
  AND TaskSlug NOT IN (SELECT TaskSlug FROM HeatmapParentTask WHERE UPPER(Status) = 'IN_PROGRESS')`

// PurgeStaleActionLogs removes action logs older than the specified retention window,
// strictly preserving logs associated with tasks currently in progress.
func (h *HeatmapDB) PurgeStaleActionLogs(olderThanDays int) (int64, error) {
	if olderThanDays <= 0 {
		return 0, nil
	}

	res, err := h.conn.Exec(sqlPurgeStaleActionLogs, olderThanDays)
	if err != nil {
		return 0, apperror.WrapSimple(err, "heatmap.purgeStale")
	}

	return res.RowsAffected()
}

const sqlHeatmapTaskSummary = `SELECT
    COUNT(*),
    COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('READ', 'READ_FILE', 'INSPECT') THEN 1 ELSE 0 END), 0),
    COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('EDIT', 'EDIT_FILE', 'MODIFY') THEN 1 ELSE 0 END), 0),
    COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('CREATE', 'WRITE_FILE', 'NEW') THEN 1 ELSE 0 END), 0)
FROM HeatmapActionLog
WHERE (? = '' OR TaskSlug = ?)`

// GetHeatmapTaskSummary returns aggregated read, edit, and create counts for a specified task.
func (h *HeatmapDB) GetHeatmapTaskSummary(taskSlug string) (*HeatmapTaskSummary, error) {
	cleanSlug := strings.TrimSpace(taskSlug)
	var total, reads, edits, creates int
	err := h.conn.QueryRow(sqlHeatmapTaskSummary, cleanSlug, cleanSlug).Scan(&total, &reads, &edits, &creates)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.taskSummary")
	}

	other := max(0, total-(reads+edits+creates))

	return &HeatmapTaskSummary{
		TaskSlug:     cleanSlug,
		TotalActions: total,
		ReadCount:    reads,
		EditCount:    edits,
		CreateCount:  creates,
		OtherCount:   other,
	}, nil
}

const sqlUpsertHeatmapParentTask = `INSERT INTO HeatmapParentTask (
    TaskSlug, TaskName, RunDirectory, RootDbPath, RootSpecFilePath, TaskFolder, Status, Notes, Comments, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT(TaskSlug) DO UPDATE SET
    TaskName = CASE WHEN excluded.TaskName != '' THEN excluded.TaskName ELSE TaskName END,
    RunDirectory = CASE WHEN excluded.RunDirectory != '' THEN excluded.RunDirectory ELSE RunDirectory END,
    RootDbPath = CASE WHEN excluded.RootDbPath != '' THEN excluded.RootDbPath ELSE RootDbPath END,
    RootSpecFilePath = CASE WHEN excluded.RootSpecFilePath != '' THEN excluded.RootSpecFilePath ELSE RootSpecFilePath END,
    TaskFolder = CASE WHEN excluded.TaskFolder != '' THEN excluded.TaskFolder ELSE TaskFolder END,
    Status = excluded.Status,
    Notes = CASE WHEN excluded.Notes IS NOT NULL THEN excluded.Notes ELSE Notes END,
    Comments = CASE WHEN excluded.Comments IS NOT NULL THEN excluded.Comments ELSE Comments END,
    UpdatedAt = CURRENT_TIMESTAMP`

// RecordParentTask registers or updates a parent task record in the heatmap.
func (h *HeatmapDB) RecordParentTask(task HeatmapParentTaskRecord) error {
	status := task.Status
	if strings.TrimSpace(status) == "" {
		status = "IN_PROGRESS"
	}

	_, err := h.conn.Exec(sqlUpsertHeatmapParentTask,
		task.TaskSlug, task.TaskName, task.RunDirectory, task.RootDbPath,
		task.RootSpecFilePath, task.TaskFolder, status, task.Notes, task.Comments,
	)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.recordParentTask")
	}

	return nil
}
