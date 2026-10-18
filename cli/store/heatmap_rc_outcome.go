package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RepoDbPath returns the filesystem path to the per-repo heatmap database.
func (s *RcHeatmapStore) RepoDbPath() string {
	return s.repoDbPath
}

// RootDbPath returns the filesystem path to the root-heatmap database.
func (s *RcHeatmapStore) RootDbPath() string {
	return s.rootDbPath
}

// LocalDbPath returns the filesystem path to the local tracking database.
func (s *RcHeatmapStore) LocalDbPath() string {
	return s.localDbPath
}

// AiName returns the canonical AI engine name.
func (s *RcHeatmapStore) AiName() string {
	return s.aiName
}

// RepoSlug returns the target repository slug.
func (s *RcHeatmapStore) RepoSlug() string {
	return s.repoSlug
}

// GetActionLogById retrieves a single HeatmapActionLog record by ID from the per-repo database.
func (s *RcHeatmapStore) GetActionLogById(actionId int64) (*HeatmapActionLogRecord, error) {
	const q = `SELECT HeatmapActionLogId, ParentLogId, TaskSlug, RepoSlug, ActionType, FilePath, StartLine, EndLine, COALESCE(CodeSnippet, ''), COALESCE(CommandText, ''), ReasoningState, COALESCE(ReasoningNotes, ''), COALESCE(ExpectedOutcome, ''), COALESCE(ActualOutcome, ''), IsExpectationMet, ExecutionDepth, ExecutionOrder, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt FROM HeatmapActionLog WHERE HeatmapActionLogId = ?`
	row := s.repoConn.QueryRow(q, actionId)
	return scanSingleActionLog(row, "heatmap_rc.getActionById")
}

// UpdateActionOutcome updates the expectation, actual outcome, and met status of an action.
func (s *RcHeatmapStore) UpdateActionOutcome(actionId int64, expectedOutcome, actualOutcome string, isExpectationMet bool, notes string) (*HeatmapActionLogRecord, error) {
	existing, err := s.GetActionLogById(actionId)
	if err != nil {
		return nil, err
	}
	exp, act, n := resolveOutcomeFields(existing, expectedOutcome, actualOutcome, notes)
	const q = `UPDATE HeatmapActionLog SET ExpectedOutcome = ?, ActualOutcome = ?, IsExpectationMet = ?, Notes = ? WHERE HeatmapActionLogId = ?`
	if _, err := s.repoConn.Exec(q, exp, act, boolToInt(isExpectationMet), n, actionId); err != nil {
		return nil, wrapErr(err, "heatmap_rc.updateOutcome")
	}
	_ = s.incrementLocalUncommitted()
	return s.GetActionLogById(actionId)
}

// GetActionLogById retrieves a single HeatmapActionLog record by ID from the local heatmap database.
func (h *HeatmapDB) GetActionLogById(actionId int64) (*HeatmapActionLogRecord, error) {
	const q = `SELECT HeatmapActionLogId, ParentLogId, TaskSlug, RepoName, ActionType, FilePath, StartLine, EndLine, COALESCE(CodeSnippet, ''), COALESCE(CommandText, ''), ReasoningState, COALESCE(ReasoningNotes, ''), COALESCE(ExpectedOutcome, ''), COALESCE(ActualOutcome, ''), IsExpectationMet, ExecutionDepth, ExecutionOrder, COALESCE(Notes, ''), COALESCE(Comments, ''), CreatedAt FROM HeatmapActionLog WHERE HeatmapActionLogId = ?`
	row := h.conn.QueryRow(q, actionId)
	return scanSingleActionLog(row, "heatmap.getActionById")
}

// UpdateActionOutcome updates the expectation, actual outcome, and met status in the local database.
func (h *HeatmapDB) UpdateActionOutcome(actionId int64, expectedOutcome, actualOutcome string, isExpectationMet bool, notes string) (*HeatmapActionLogRecord, error) {
	existing, err := h.GetActionLogById(actionId)
	if err != nil {
		return nil, err
	}
	exp, act, n := resolveOutcomeFields(existing, expectedOutcome, actualOutcome, notes)
	const q = `UPDATE HeatmapActionLog SET ExpectedOutcome = ?, ActualOutcome = ?, IsExpectationMet = ?, Notes = ? WHERE HeatmapActionLogId = ?`
	if _, err := h.conn.Exec(q, exp, act, boolToInt(isExpectationMet), n, actionId); err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.updateOutcome")
	}
	return h.GetActionLogById(actionId)
}

func resolveOutcomeFields(existing *HeatmapActionLogRecord, expected, actual, notes string) (string, string, string) {
	exp, act, n := expected, actual, notes
	if exp == "" {
		exp = existing.ExpectedOutcome
	}
	if act == "" {
		act = existing.ActualOutcome
	}
	if n == "" {
		n = existing.Notes
	}
	return exp, act, n
}

func scanSingleActionLog(row *sql.Row, scope string) (*HeatmapActionLogRecord, error) {
	var r HeatmapActionLogRecord
	var parentID sql.NullInt64
	var isMet int
	if err := row.Scan(
		&r.HeatmapActionLogId, &parentID, &r.TaskSlug, &r.RepoName, &r.ActionType,
		&r.FilePath, &r.StartLine, &r.EndLine, &r.CodeSnippet, &r.CommandText,
		&r.ReasoningState, &r.ReasoningNotes, &r.ExpectedOutcome,
		&r.ActualOutcome, &isMet, &r.ExecutionDepth, &r.ExecutionOrder,
		&r.Notes, &r.Comments, &r.CreatedAt,
	); err != nil {
		return nil, apperror.WrapSimple(err, scope)
	}
	if parentID.Valid {
		pid := parentID.Int64
		r.ParentLogId = &pid
	}
	r.IsExpectationMet = isMet == 1
	return &r, nil
}
