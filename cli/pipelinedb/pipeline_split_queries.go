package pipelinedb

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func iterateRunIdRows(rows *sql.Rows) PipelineRunIdSliceResult {
	var list []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return result.FailSlice[uint64](apperror.WrapSimple(err, "scan cached run id"))
		}
		list = append(list, id)
	}

	return result.OkSlice(list)
}

func collectRunIdList(rows *sql.Rows) PipelineRunIdSliceResult {
	res := iterateRunIdRows(rows)
	if res.IsFailure() {
		return res
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return result.FailSlice[uint64](apperror.WrapSimple(rowsErr, "iterate cached run ids"))
	}

	return res
}

// QueryCachedErrorRunIds retrieves distinct RunIds cached in PipelineErrorLog.
func (p *PipelineSplitDb) QueryCachedErrorRunIds() PipelineRunIdSliceResult {
	rows, err := p.conn.Query(sqlQueryCachedErrorRunIds)
	if err != nil {
		return result.FailSlice[uint64](apperror.WrapSimple(err, "query cached error run ids"))
	}

	defer rows.Close()

	return collectRunIdList(rows)
}

// QueryCachedErrorRunIdMap returns a map set of cached error RunIds for fast lookups.
func (p *PipelineSplitDb) QueryCachedErrorRunIdMap() PipelineRunIdMapResult {
	runRes := p.QueryCachedErrorRunIds()
	if runRes.IsFailure() {
		return result.FailMap[uint64, bool](runRes.AppError())
	}

	idMap := make(map[uint64]bool, runRes.Count())
	for _, id := range runRes.Data {
		idMap[id] = true
	}

	return result.OkMap(idMap)
}

func normalizeNegativeOffset(offset int) int {
	if offset < 0 {
		offset = -offset
	}

	if offset > 0 {
		return offset - 1
	}

	return 0
}

// QueryRunByNegativeOffset retrieves a run by 1-based negative offset (-1 = latest).
func (p *PipelineSplitDb) QueryRunByNegativeOffset(offset int) (*PipelineRunRecord, error) {
	rows, err := p.conn.Query(sqlQueryRunByOffset, normalizeNegativeOffset(offset))
	if err != nil {
		return nil, apperror.WrapSimple(err, "query run by offset")
	}
	defer rows.Close()

	runsRes := collectRecentRuns(rows)
	if runsRes.IsFailure() || runsRes.IsEmpty() {
		return nil, runsRes.AppError()
	}

	return &runsRes.Data[0], nil
}

// QueryRunsByNegativeOffset is an alias for QueryRunByNegativeOffset.
func (p *PipelineSplitDb) QueryRunsByNegativeOffset(offset int) (*PipelineRunRecord, error) {
	return p.QueryRunByNegativeOffset(offset)
}

// QueryLastFailedRuns retrieves the most recent failed pipeline runs up to limit.
func (p *PipelineSplitDb) QueryLastFailedRuns(limit int) PipelineRunSliceResult {
	limitVal := resolveLimit(limit, 5)
	rows, err := p.conn.Query(sqlQueryLastFailedRuns, limitVal)
	if err != nil {
		return result.FailSlice[PipelineRunRecord](apperror.WrapSimple(err, "query last failed runs"))
	}

	defer rows.Close()

	return collectRecentRuns(rows)
}

// QueryLastNFailedRuns is an alias for QueryLastFailedRuns.
func (p *PipelineSplitDb) QueryLastNFailedRuns(limit int) PipelineRunSliceResult {
	return p.QueryLastFailedRuns(limit)
}

// QueryErrorLogsByRunId retrieves all error diagnostics recorded for a specific run ID.
func (p *PipelineSplitDb) QueryErrorLogsByRunId(runId uint64) PipelineErrorSliceResult {
	rows, err := p.conn.Query(sqlQueryErrorLogsByRunId, runId)
	if err != nil {
		return result.FailSlice[PipelineErrorRecord](apperror.WrapSimple(err, "query error logs by run id"))
	}

	defer rows.Close()

	return collectRecentErrors(rows)
}

// QueryDetailedErrorLogsByRunId retrieves uncompressed detailed errors for a run ID.
func (p *PipelineSplitDb) QueryDetailedErrorLogsByRunId(runId uint64) PipelineErrorSliceResult {
	rows, err := p.conn.Query(sqlQueryDetailErrorLogsByRunId, runId)
	if err != nil {
		return result.FailSlice[PipelineErrorRecord](apperror.WrapSimple(err, "query detail error logs by run id"))
	}

	defer rows.Close()

	return collectRecentErrors(rows)
}

// QueryCompactErrorLogsByRunId retrieves compact filtered errors for a run ID.
func (p *PipelineSplitDb) QueryCompactErrorLogsByRunId(runId uint64) PipelineCompactErrorSliceResult {
	rows, err := p.conn.Query(sqlQueryCompactErrorLogsByRunId, runId)
	if err != nil {
		return result.FailSlice[PipelineCompactErrorRecord](apperror.WrapSimple(err, "query compact error logs by run id"))
	}

	defer rows.Close()

	return collectRecentCompactErrors(rows)
}

// GetRunWorkflowName looks up the workflow name for a run ID from PipelineRun.
func (p *PipelineSplitDb) GetRunWorkflowName(runId uint64) string {
	var name string
	err := p.conn.QueryRow("SELECT WorkflowName FROM PipelineRun WHERE RunId = ? LIMIT 1;", runId).Scan(&name)
	if err != nil || len(name) == 0 {
		return "CI"
	}

	return name
}

// QuerySuccessfulRunDurations retrieves historical durations of successful runs.
func (p *PipelineSplitDb) QuerySuccessfulRunDurations(workflowName string, limit int) []int {
	if p.conn == nil || limit <= 0 {
		return nil
	}
	query := "SELECT DurationSeconds FROM PipelineRun WHERE IsSuccess = 1 AND DurationSeconds >= 10 AND (WorkflowName = ? OR ? = '') ORDER BY CreatedAt DESC, RunId DESC LIMIT ?;"
	rows, err := p.conn.Query(query, workflowName, workflowName, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()

	return scanDurations(rows)
}

func scanDurations(rows *sql.Rows) []int {
	var durs []int
	for rows.Next() {
		var dur int
		if err := rows.Scan(&dur); err == nil && dur > 0 {
			durs = append(durs, dur)
		}
	}

	return durs
}
