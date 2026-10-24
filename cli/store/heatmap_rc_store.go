package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// RcHeatmapStore manages multi-tier database connections for RC heatmap telemetry.
type RcHeatmapStore struct {
	aiName, repoSlug, rootDbPath, repoDbPath, localDbPath string
	rootConn, repoConn, localConn                         *sql.DB
}

// OpenRcHeatmapStore opens connections to root, repo, and local tracking databases.
func OpenRcHeatmapStore(aiName, repoSlug string) (*RcHeatmapStore, error) {
	normAI, cleanSlug := NormalizeAIName(aiName), SanitizeSlug(repoSlug)
	rootPath, repoPath, localPath := ResolveRcHeatmapRootDbPath(normAI), ResolveRcHeatmapRepoDbPath(normAI, cleanSlug), ResolveRcHeatmapLocalTrackingDbPath()
	rootConn, err := openConfiguredDB(rootPath, EnsureRcHeatmapRootTables)
	if err != nil {
		return nil, err
	}
	repoConn, err := openConfiguredDB(repoPath, EnsureRcHeatmapRepoTables)
	if err != nil {
		_ = rootConn.Close()
		return nil, err
	}
	localConn, err := openConfiguredDB(localPath, EnsureRcHeatmapLocalTrackingTables)
	if err != nil {
		_ = rootConn.Close()
		_ = repoConn.Close()
		return nil, err
	}
	return &RcHeatmapStore{aiName: normAI, repoSlug: cleanSlug, rootConn: rootConn, repoConn: repoConn, localConn: localConn, rootDbPath: rootPath, repoDbPath: repoPath, localDbPath: localPath}, nil
}

func openConfiguredDB(dbPath string, schemaInit func(*sql.DB) error) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "heatmap_rc.mkdir")
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap_rc.open")
	}
	if err := ConfigureRcHeatmapPragmas(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := schemaInit(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

// Close closes all underlying database connections.
func (s *RcHeatmapStore) Close() error {
	var firstErr error
	for _, c := range []*sql.DB{s.rootConn, s.repoConn, s.localConn} {
		if c != nil && firstErr == nil {
			firstErr = c.Close()
		}
	}
	return firstErr
}

// LogAction records an action in the per-repo database and updates tracking counters.
func (s *RcHeatmapStore) LogAction(log HeatmapActionLogRecord) (int64, error) {
	depth, order, err := s.calculateDepthAndOrder(log.ParentLogId, log.TaskSlug)
	if err != nil {
		return 0, err
	}
	const q = `INSERT INTO HeatmapActionLog (ParentLogId, TaskSlug, RepoSlug, ActionType, FilePath, StartLine, EndLine, CodeSnippet, CommandText, ReasoningState, ReasoningNotes, ExpectedOutcome, ActualOutcome, IsExpectationMet, ExecutionDepth, ExecutionOrder, DurationMs, ModelName, ToolName, Notes, Comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`
	res, err := s.repoConn.Exec(q, log.ParentLogId, log.TaskSlug, s.repoSlug, log.ActionType, log.FilePath, log.StartLine, log.EndLine, log.CodeSnippet, log.CommandText, log.ReasoningState, log.ReasoningNotes, log.ExpectedOutcome, log.ActualOutcome, boolToInt(log.IsExpectationMet), depth, order, log.ModelName, log.ToolName, log.Notes, log.Comments)
	if err != nil {
		return 0, apperror.WrapSimple(err, "heatmap_rc.logAction")
	}
	actionId, err := res.LastInsertId()
	if err != nil {
		return 0, apperror.WrapSimple(err, "heatmap_rc.lastInsertId")
	}
	s.insertExecutionTreeNode(actionId, log.ParentLogId, depth, order)
	_ = s.incrementLocalUncommitted()
	return actionId, nil
}

func (s *RcHeatmapStore) calculateDepthAndOrder(parentID *int64, taskSlug string) (int, int, error) {
	if parentID == nil || *parentID <= 0 {
		var maxOrder sql.NullInt64
		_ = s.repoConn.QueryRow(`SELECT MAX(ExecutionOrder) FROM HeatmapActionLog WHERE TaskSlug = ? AND ParentLogId IS NULL`, taskSlug).Scan(&maxOrder)
		return 0, int(maxOrder.Int64) + 1, nil
	}
	var parentDepth int
	if err := s.repoConn.QueryRow(`SELECT ExecutionDepth FROM HeatmapActionLog WHERE HeatmapActionLogId = ?`, *parentID).Scan(&parentDepth); err != nil {
		return 0, 0, apperror.WrapSimple(err, "heatmap_rc.parentLookup")
	}
	var maxChildOrder sql.NullInt64
	_ = s.repoConn.QueryRow(`SELECT MAX(ExecutionOrder) FROM HeatmapActionLog WHERE ParentLogId = ?`, *parentID).Scan(&maxChildOrder)
	return parentDepth + 1, int(maxChildOrder.Int64) + 1, nil
}

func (s *RcHeatmapStore) resolveParentTreeVector(nodeId int64, parentId int64) (int64, string) {
	rootId, pathVector := nodeId, fmt.Sprintf("/%d", nodeId)
	var pRoot int64
	var pPath string
	if err := s.repoConn.QueryRow(`SELECT RootLogId, PathVector FROM HeatmapExecutionTree WHERE NodeLogId = ?`, parentId).Scan(&pRoot, &pPath); err == nil {
		rootId = pRoot
		pathVector = fmt.Sprintf("%s/%d", pPath, nodeId)
	}

	_, _ = s.repoConn.Exec(`UPDATE HeatmapExecutionTree SET IsLeaf = 0, HasChildren = 1 WHERE NodeLogId = ?`, parentId)

	return rootId, pathVector
}

func (s *RcHeatmapStore) insertExecutionTreeNode(nodeId int64, parentId *int64, depth, order int) {
	rootId, pathVector := nodeId, fmt.Sprintf("/%d", nodeId)
	if parentId != nil && *parentId > 0 {
		rootId, pathVector = s.resolveParentTreeVector(nodeId, *parentId)
	}

	const q = `INSERT INTO HeatmapExecutionTree (RootLogId, NodeLogId, ParentLogId, NodeDepth, SiblingOrder, PathVector, IsLeaf, HasChildren) VALUES (?, ?, ?, ?, ?, ?, 1, 0)`
	_, _ = s.repoConn.Exec(q, rootId, nodeId, parentId, depth, order, pathVector)
}

// RecordParentTask inserts or updates a parent task record in root-heatmap.db.
func (s *RcHeatmapStore) RecordParentTask(t RcHeatmapParentTaskRecord) error {
	const q = `INSERT INTO HeatmapParentTask (TaskSlug, TaskName, RepoSlug, Status, StepBudget, CompletedSteps, HasCompleted, IsActive, StartedAt, Notes, Comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?, ?) ON CONFLICT(TaskSlug) DO UPDATE SET Status = excluded.Status, CompletedSteps = excluded.CompletedSteps, HasCompleted = excluded.HasCompleted, IsActive = excluded.IsActive, CompletedAt = CASE WHEN excluded.HasCompleted = 1 THEN CURRENT_TIMESTAMP ELSE CompletedAt END, UpdatedAt = CURRENT_TIMESTAMP`
	_, err := s.rootConn.Exec(q, t.TaskSlug, t.TaskName, s.repoSlug, t.Status, t.StepBudget, t.CompletedSteps, boolToInt(t.HasCompleted), boolToInt(t.IsActive), t.Notes, t.Comments)
	return wrapErr(err, "heatmap_rc.recordParentTask")
}

// RecordAiSession inserts or updates an AI execution session record.
func (s *RcHeatmapStore) RecordAiSession(sess HeatmapAiSessionRecord) error {
	const q = `INSERT INTO HeatmapAiSession (SessionUuid, AiName, SequencePrefix, MachineHostname, OsPlatform, AgentRole, IsActive, HasCompleted, TotalTasksExecuted, TotalActionsLogged, Notes, Comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(SessionUuid) DO UPDATE SET IsActive = excluded.IsActive, HasCompleted = excluded.HasCompleted, TotalTasksExecuted = excluded.TotalTasksExecuted, TotalActionsLogged = excluded.TotalActionsLogged, CompletedAt = CASE WHEN excluded.HasCompleted = 1 THEN CURRENT_TIMESTAMP ELSE CompletedAt END, UpdatedAt = CURRENT_TIMESTAMP`
	_, err := s.rootConn.Exec(q, sess.SessionUuid, sess.AiName, sess.SequencePrefix, sess.MachineHostname, sess.OsPlatform, sess.AgentRole, boolToInt(sess.IsActive), boolToInt(sess.HasCompleted), sess.TotalTasksExecuted, sess.TotalActionsLogged, sess.Notes, sess.Comments)
	return wrapErr(err, "heatmap_rc.recordAiSession")
}

// RecordRepoRegistry inserts or updates a repo entry in root-heatmap.db.
func (s *RcHeatmapStore) RecordRepoRegistry(r HeatmapRepoRegistryRecord) error {
	const q = `INSERT INTO HeatmapRepoRegistry (RepoSlug, RepoUrl, DefaultBranch, IsActive, TotalActionsRecorded, Notes, Comments) VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT(RepoSlug) DO UPDATE SET RepoUrl = CASE WHEN excluded.RepoUrl != '' THEN excluded.RepoUrl ELSE RepoUrl END, LastActiveAt = CURRENT_TIMESTAMP, TotalActionsRecorded = excluded.TotalActionsRecorded, UpdatedAt = CURRENT_TIMESTAMP`
	_, err := s.rootConn.Exec(q, r.RepoSlug, r.RepoUrl, r.DefaultBranch, boolToInt(r.IsActive), r.TotalActionsRecorded, r.Notes, r.Comments)
	return wrapErr(err, "heatmap_rc.recordRepoRegistry")
}

// GetOrCreateLocalTracking fetches or initializes the tracking record for repo and AI.
func (s *RcHeatmapStore) GetOrCreateLocalTracking(repoSlug, aiName string) (*HeatmapLocalTrackingRecord, error) {
	normAI, cleanSlug := NormalizeAIName(aiName), SanitizeSlug(repoSlug)
	const qSel = `SELECT HeatmapLocalTrackingId, RepoSlug, AiName, SequencePrefix, RepoCachePath, LastCommittedSha, LastCommitTimestamp, UncommittedChangeCount, SyncStatus, IsDirty, LastSyncAttemptAt, LastSuccessfulSyncAt, ErrorMessage, CreatedAt, UpdatedAt FROM HeatmapLocalTracking WHERE RepoSlug = ? AND AiName = ?`
	var r HeatmapLocalTrackingRecord
	var isDirty int
	err := s.localConn.QueryRow(qSel, cleanSlug, normAI).Scan(&r.HeatmapLocalTrackingId, &r.RepoSlug, &r.AiName, &r.SequencePrefix, &r.RepoCachePath, &r.LastCommittedSha, &r.LastCommitTimestamp, &r.UncommittedChangeCount, &r.SyncStatus, &isDirty, &r.LastSyncAttemptAt, &r.LastSuccessfulSyncAt, &r.ErrorMessage, &r.CreatedAt, &r.UpdatedAt)
	if err == nil {
		r.IsDirty = isDirty == 1
		return &r, nil
	}
	if err != sql.ErrNoRows {
		return nil, wrapErr(err, "heatmap_rc.getTracking")
	}
	prefix, rcRoot := AISequencePrefix(normAI), ResolveRepoCacheRoot()
	const qIns = `INSERT INTO HeatmapLocalTracking (RepoSlug, AiName, SequencePrefix, RepoCachePath, LastCommittedSha, LastCommitTimestamp, UncommittedChangeCount, SyncStatus, IsDirty) VALUES (?, ?, ?, ?, '', '1970-01-01T00:00:00Z', 0, 'IDLE', 0)`
	res, err := s.localConn.Exec(qIns, cleanSlug, normAI, prefix, rcRoot)
	if err != nil {
		return nil, wrapErr(err, "heatmap_rc.initTracking")
	}
	id, _ := res.LastInsertId()
	return &HeatmapLocalTrackingRecord{HeatmapLocalTrackingId: id, RepoSlug: cleanSlug, AiName: normAI, SequencePrefix: prefix, RepoCachePath: rcRoot, LastCommitTimestamp: "1970-01-01T00:00:00Z", SyncStatus: "IDLE"}, nil
}

// UpdateLocalTracking updates the tracking state in the local database.
func (s *RcHeatmapStore) UpdateLocalTracking(t HeatmapLocalTrackingRecord) error {
	const sqlUpdate = `UPDATE HeatmapLocalTracking SET LastCommittedSha = ?, LastCommitTimestamp = ?, UncommittedChangeCount = ?, SyncStatus = ?, IsDirty = ?, ErrorMessage = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE RepoSlug = ? AND AiName = ?`
	_, err := s.localConn.Exec(sqlUpdate, t.LastCommittedSha, t.LastCommitTimestamp, t.UncommittedChangeCount, t.SyncStatus, boolToInt(t.IsDirty), t.ErrorMessage, t.RepoSlug, t.AiName)
	return wrapErr(err, "heatmap_rc.updateTracking")
}

func (s *RcHeatmapStore) incrementLocalUncommitted() error {
	if _, err := s.GetOrCreateLocalTracking(s.repoSlug, s.aiName); err != nil {
		return err
	}
	const sqlInc = `UPDATE HeatmapLocalTracking SET UncommittedChangeCount = UncommittedChangeCount + 1, IsDirty = 1, SyncStatus = 'DIRTY', UpdatedAt = CURRENT_TIMESTAMP WHERE RepoSlug = ? AND AiName = ?`
	_, err := s.localConn.Exec(sqlInc, s.repoSlug, s.aiName)
	return err
}

// GetUncommittedActionCount returns uncommitted change count from local tracking.
func (s *RcHeatmapStore) GetUncommittedActionCount(repoSlug, aiName string) (int, error) {
	rec, err := s.GetOrCreateLocalTracking(repoSlug, aiName)
	if err != nil {
		return 0, err
	}
	return rec.UncommittedChangeCount, nil
}

// SearchActions searches the per-repo action log by task, query substring, and filepath.
func (s *RcHeatmapStore) SearchActions(query, taskSlug, filePath string, limit int) ([]HeatmapActionLogRecord, error) {
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
	sqlText := `SELECT HeatmapActionLogId, ParentLogId, TaskSlug, RepoSlug, ActionType, FilePath, StartLine, EndLine, COALESCE(CodeSnippet, ''), COALESCE(CommandText, ''), ReasoningState, COALESCE(ReasoningNotes, ''), COALESCE(ExpectedOutcome, ''), COALESCE(ActualOutcome, ''), IsExpectationMet, ExecutionDepth, ExecutionOrder, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt FROM HeatmapActionLog`
	if len(whereClauses) > 0 {
		sqlText += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	sqlText += " ORDER BY CreatedAt DESC LIMIT ?"
	args = append(args, limit)
	rows, err := s.repoConn.Query(sqlText, args...)
	if err != nil {
		return nil, wrapErr(err, "heatmap_rc.searchActions")
	}
	defer rows.Close()
	return scanActionLogs(rows)
}

// GetExecutionTree constructs a hierarchical execution tree for the given task.
func (s *RcHeatmapStore) GetExecutionTree(taskSlug string) ([]RcHeatmapExecutionTreeNode, error) {
	records, err := s.SearchActions("", taskSlug, "", 1000)
	if err != nil {
		return nil, err
	}
	return AssembleRcTreeNodes(records), nil
}

// GetStatusSummary computes volume, expectation percentage, and top touched files.
func (s *RcHeatmapStore) GetStatusSummary(taskSlug string) (*HeatmapStatusSummary, error) {
	records, err := s.SearchActions("", taskSlug, "", 5000)
	if err != nil {
		return nil, err
	}
	summary := &HeatmapStatusSummary{TaskSlug: taskSlug, TotalActions: len(records), ReasoningStateCounts: make(map[string]int)}
	fileTouches, metCount := make(map[string]int), 0
	for _, r := range records {
		if r.IsExpectationMet {
			metCount++
		}
		summary.ReasoningStateCounts[r.ReasoningState]++
		if r.FilePath != "" {
			fileTouches[r.FilePath]++
		}
	}
	if len(records) > 0 {
		summary.ExpectationsMetCount, summary.ExpectationsMetPct = metCount, (float64(metCount)/float64(len(records)))*100.0
	}
	for f, c := range fileTouches {
		summary.TopTouchedFiles = append(summary.TopTouchedFiles, FileTouchCount{FilePath: f, Count: c})
	}
	return summary, nil
}

// QuerySessionData retrieves AI sessions from root-heatmap.db.
func (s *RcHeatmapStore) QuerySessionData(aiName string) ([]HeatmapAiSessionRecord, error) {
	const query = `SELECT HeatmapAiSessionId, SessionUuid, AiName, SequencePrefix, MachineHostname, OsPlatform, AgentRole, StartedAt, CompletedAt, IsActive, HasCompleted, TotalTasksExecuted, TotalActionsLogged, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt, UpdatedAt FROM HeatmapAiSession WHERE AiName = ? ORDER BY CreatedAt DESC LIMIT 50`
	rows, err := s.rootConn.Query(query, NormalizeAIName(aiName))
	if err != nil {
		return nil, wrapErr(err, "heatmap_rc.querySession")
	}
	defer rows.Close()
	return ScanSessionRows(rows)
}

// QueryRepoRegistry retrieves registered repos from root-heatmap.db.
func (s *RcHeatmapStore) QueryRepoRegistry(aiName string) ([]HeatmapRepoRegistryRecord, error) {
	const query = `SELECT HeatmapRepoRegistryId, RepoSlug, RepoUrl, DefaultBranch, FirstSeenAt, LastActiveAt, IsActive, TotalActionsRecorded, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt, UpdatedAt FROM HeatmapRepoRegistry ORDER BY LastActiveAt DESC LIMIT 100`
	rows, err := s.rootConn.Query(query)
	if err != nil {
		return nil, wrapErr(err, "heatmap_rc.queryRepos")
	}
	defer rows.Close()
	return ScanRepoRegistryRows(rows)
}

func wrapErr(err error, scope string) error {
	if err == nil {
		return nil
	}
	return apperror.WrapSimple(err, scope)
}
