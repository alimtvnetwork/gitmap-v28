package store

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// QueryActionLogs executes a read query and returns scanned HeatmapActionLogRecord records.
func (h *HeatmapDB) QueryActionLogs(query string, args ...any) ([]HeatmapActionLogRecord, error) {
	rows, err := h.conn.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.queryActionLogs")
	}

	defer rows.Close()

	return scanActionLogs(rows)
}

// QueryHeatmapActionLogs executes a read query on any sql.DB and returns records.
func QueryHeatmapActionLogs(db *sql.DB, query string, args ...any) ([]HeatmapActionLogRecord, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.queryActionLogs")
	}

	defer rows.Close()

	return scanActionLogs(rows)
}

type actionLogScanTargets struct {
	record   HeatmapActionLogRecord
	parentID sql.NullInt64
	isMet    int
	code     sql.NullString
	cmd      sql.NullString
	reason   sql.NullString
	exp      sql.NullString
	act      sql.NullString
	model    sql.NullString
	tool     sql.NullString
	notes    sql.NullString
	comments sql.NullString
}

func (t *actionLogScanTargets) bindCoreTarget(col string) any {
	switch strings.ToLower(col) {
	case "heatmapactionlogid":
		return &t.record.HeatmapActionLogId
	case "parentlogid":
		return &t.parentID
	case "taskslug":
		return &t.record.TaskSlug
	case "reponame", "reposlug":
		return &t.record.RepoName
	case "actiontype":
		return &t.record.ActionType
	case "filepath":
		return &t.record.FilePath
	case "startline":
		return &t.record.StartLine
	case "endline":
		return &t.record.EndLine
	default:
		return nil
	}
}

func (t *actionLogScanTargets) bindDetailsTarget(col string) any {
	switch strings.ToLower(col) {
	case "codesnippet":
		return &t.code
	case "commandtext":
		return &t.cmd
	case "reasoningstate":
		return &t.record.ReasoningState
	case "reasoningnotes":
		return &t.reason
	case "expectedoutcome":
		return &t.exp
	case "actualoutcome":
		return &t.act
	case "isexpectationmet":
		return &t.isMet
	default:
		return nil
	}
}

func (t *actionLogScanTargets) bindMetaTarget(col string) any {
	switch strings.ToLower(col) {
	case "executiondepth":
		return &t.record.ExecutionDepth
	case "executionorder":
		return &t.record.ExecutionOrder
	case "modelname":
		return &t.model
	case "toolname":
		return &t.tool
	case "notes":
		return &t.notes
	case "comments":
		return &t.comments
	case "createdat":
		return &t.record.CreatedAt
	default:
		return nil
	}
}

func (t *actionLogScanTargets) bindDetailsOrMeta(col string) any {
	if p := t.bindDetailsTarget(col); p != nil {
		return p
	}

	return t.bindMetaTarget(col)
}

func (t *actionLogScanTargets) bindTarget(col string) any {
	if p := t.bindCoreTarget(col); p != nil {
		return p
	}

	if p := t.bindDetailsOrMeta(col); p != nil {
		return p
	}

	var discard any

	return &discard
}

func (t *actionLogScanTargets) finalizeStrings() {
	t.finalizeCommandStrings()
	t.finalizeOutcomeStrings()
}

func (t *actionLogScanTargets) finalizeCommandStrings() {
	if t.code.Valid {
		t.record.CodeSnippet = t.code.String
	}
	if t.cmd.Valid {
		t.record.CommandText = t.cmd.String
	}
	if t.reason.Valid {
		t.record.ReasoningNotes = t.reason.String
	}
}

func (t *actionLogScanTargets) finalizeOutcomeStrings() {
	if t.exp.Valid {
		t.record.ExpectedOutcome = t.exp.String
	}
	if t.act.Valid {
		t.record.ActualOutcome = t.act.String
	}
}

func (t *actionLogScanTargets) finalizeMeta() {
	t.finalizeToolMeta()
	t.finalizeNotesMeta()
}

func (t *actionLogScanTargets) finalizeToolMeta() {
	if t.model.Valid {
		t.record.ModelName = t.model.String
	}
	if t.tool.Valid {
		t.record.ToolName = t.tool.String
	}
}

func (t *actionLogScanTargets) finalizeNotesMeta() {
	if t.notes.Valid {
		t.record.Notes = t.notes.String
	}
	if t.comments.Valid {
		t.record.Comments = t.comments.String
	}
}

func (t *actionLogScanTargets) finalize() HeatmapActionLogRecord {
	if t.parentID.Valid {
		pid := t.parentID.Int64
		t.record.ParentLogId = &pid
	}

	t.record.IsExpectationMet = t.isMet == 1
	t.finalizeStrings()
	t.finalizeMeta()

	return t.record
}

func scanActionLogRow(rows *sql.Rows, cols []string) (HeatmapActionLogRecord, error) {
	var targets actionLogScanTargets
	dest := make([]any, len(cols))
	for i, col := range cols {
		dest[i] = targets.bindTarget(col)
	}

	if err := rows.Scan(dest...); err != nil {
		return targets.record, apperror.WrapSimple(err, "heatmap.scanRow")
	}

	return targets.finalize(), nil
}

func scanActionLogsRows(rows *sql.Rows, cols []string) ([]HeatmapActionLogRecord, error) {
	var records []HeatmapActionLogRecord
	for rows.Next() {
		rec, err := scanActionLogRow(rows, cols)
		if err != nil {
			return nil, err
		}

		records = append(records, rec)
	}

	return records, nil
}

func scanActionLogs(rows *sql.Rows) ([]HeatmapActionLogRecord, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.scanColumns")
	}

	return scanActionLogsRows(rows, cols)
}

func tryAttachNodeToParent(node *HeatmapExecutionTreeNode, parentID *int64, nodeMap map[int64]*HeatmapExecutionTreeNode) bool {
	hasParent := parentID != nil && *parentID > 0
	if !hasParent {
		return false
	}

	parent, hasNode := nodeMap[*parentID]
	if !hasNode {
		return false
	}

	parent.Children = append(parent.Children, *node)

	return true
}

func buildTreeNodeMap(records []HeatmapActionLogRecord) map[int64]*HeatmapExecutionTreeNode {
	nodeMap := make(map[int64]*HeatmapExecutionTreeNode, len(records))
	for i := range records {
		rec := records[i]
		nodeMap[rec.HeatmapActionLogId] = &HeatmapExecutionTreeNode{
			HeatmapActionLogRecord: rec,
			Children:               make([]HeatmapExecutionTreeNode, 0),
		}
	}

	return nodeMap
}

func collectRootTreeNodes(records []HeatmapActionLogRecord, nodeMap map[int64]*HeatmapExecutionTreeNode) []HeatmapExecutionTreeNode {
	var roots []HeatmapExecutionTreeNode
	for _, rec := range records {
		node := nodeMap[rec.HeatmapActionLogId]
		hasAttached := tryAttachNodeToParent(node, rec.ParentLogId, nodeMap)
		if hasAttached {
			continue
		}

		roots = append(roots, *node)
	}

	return roots
}

func assembleTreeNodes(records []HeatmapActionLogRecord) []HeatmapExecutionTreeNode {
	nodeMap := buildTreeNodeMap(records)

	return collectRootTreeNodes(records, nodeMap)
}
