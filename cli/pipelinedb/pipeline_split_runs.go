package pipelinedb

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func isRunSuccess(r PipelineRunRecord) int {
	if r.IsSuccess {
		return 1
	}

	if r.Conclusion == "success" {
		return 1
	}

	return 0
}

// RecordRun inserts or updates a pipeline run execution record.
func (p *PipelineSplitDb) RecordRun(r PipelineRunRecord) error {
	isSuccessInt := isRunSuccess(r)
	_, err := p.conn.Exec(sqlRecordRun,
		r.RunId, r.RepoSlug, r.WorkflowName, r.Status, r.Conclusion, r.Branch, r.Sha,
		r.EtaSeconds, r.DurationSeconds, r.RunUrl, isSuccessInt, r.Notes, r.Comments, r.CreatedAt, r.UpdatedAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record pipeline run")
	}

	return nil
}

func resolveCreatedAt(createdAt string) string {
	if createdAt != "" {
		return createdAt
	}

	return time.Now().UTC().Format(time.RFC3339)
}

// RecordErrorLog inserts an error diagnostic entry for a failing run.
func (p *PipelineSplitDb) RecordErrorLog(e PipelineErrorRecord) error {
	createdAt := resolveCreatedAt(e.CreatedAt)
	_, err := p.conn.Exec(sqlRecordErrorLog,
		e.RunId, e.RepoSlug, e.WorkflowName, e.StepName, e.ErrorText, e.RawLogs, e.Notes, e.Comments, createdAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record pipeline error log")
	}

	return nil
}

// HasErrorLog checks if an error diagnostic entry for a run has already been recorded.
func (p *PipelineSplitDb) HasErrorLog(runId uint64) bool {
	var exists int
	const query = `SELECT 1 FROM (
		SELECT RunId FROM PipelineErrorLog
		UNION
		SELECT RunId FROM PipelineDetailErrorLog
		UNION
		SELECT RunId FROM PipelineCompactErrorLog
	) WHERE RunId = ? LIMIT 1;`
	err := p.conn.QueryRow(query, runId).Scan(&exists)
	if err != nil {
		return false
	}

	return exists == 1
}

// RecordDetailErrorLog inserts an uncompressed diagnostic error log into PipelineDetailErrorLog.
func (p *PipelineSplitDb) RecordDetailErrorLog(e PipelineErrorRecord) error {
	createdAt := resolveCreatedAt(e.CreatedAt)
	_, err := p.conn.Exec(sqlRecordDetailErrorLog,
		e.RunId, e.RepoSlug, e.WorkflowName, e.StepName, e.ErrorText, e.RawLogs, e.Notes, e.Comments, createdAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record pipeline detail error log")
	}

	return nil
}

// RecordCompactErrorLog inserts a noise-filtered compact error log into PipelineCompactErrorLog.
func (p *PipelineSplitDb) RecordCompactErrorLog(c PipelineCompactErrorRecord) error {
	createdAt := resolveCreatedAt(c.CreatedAt)
	_, err := p.conn.Exec(sqlRecordCompactErrorLog,
		c.RunId, c.RepoSlug, c.WorkflowName, c.StepName, c.ErrorText, c.CompactLogs, c.FilteredOkCount, c.Notes, c.Comments, createdAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record pipeline compact error log")
	}

	return nil
}

// RecordDualErrorLog inserts both detail and compact error diagnostics for a run.
func (p *PipelineSplitDb) RecordDualErrorLog(detail PipelineErrorRecord, compact PipelineCompactErrorRecord) error {
	_ = p.RecordErrorLog(detail)
	if err := p.RecordDetailErrorLog(detail); err != nil {
		return err
	}

	return p.RecordCompactErrorLog(compact)
}

// HasDetailErrorLog checks if a detail error log exists for the run.
func (p *PipelineSplitDb) HasDetailErrorLog(runId uint64) bool {
	var exists int
	err := p.conn.QueryRow("SELECT 1 FROM PipelineDetailErrorLog WHERE RunId = ? LIMIT 1;", runId).Scan(&exists)
	if err != nil {
		return false
	}

	return exists == 1
}

// HasCompactErrorLog checks if a compact error log exists for the run.
func (p *PipelineSplitDb) HasCompactErrorLog(runId uint64) bool {
	var exists int
	err := p.conn.QueryRow("SELECT 1 FROM PipelineCompactErrorLog WHERE RunId = ? LIMIT 1;", runId).Scan(&exists)
	if err != nil {
		return false
	}

	return exists == 1
}

func resolveLimit(limit int, fallback int) int {
	if limit <= 0 {
		return fallback
	}

	return limit
}

func scanPipelineRun(rows *sql.Rows) (PipelineRunRecord, *apperror.AppError) {
	var r PipelineRunRecord
	var isSuccessInt int
	err := rows.Scan(
		&r.RunId, &r.RepoSlug, &r.WorkflowName, &r.Status, &r.Conclusion, &r.Branch, &r.Sha,
		&r.EtaSeconds, &r.DurationSeconds, &r.RunUrl, &isSuccessInt, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return r, apperror.WrapSimple(err, "scan pipeline run row")
	}

	r.IsSuccess = isSuccessInt == 1

	return r, nil
}

func iterateRunRows(rows *sql.Rows) PipelineRunSliceResult {
	var list []PipelineRunRecord
	for rows.Next() {
		r, scanErr := scanPipelineRun(rows)
		if scanErr != nil {
			return result.FailSlice[PipelineRunRecord](scanErr)
		}
		list = append(list, r)
	}

	return result.OkSlice(list)
}

func collectRecentRuns(rows *sql.Rows) PipelineRunSliceResult {
	res := iterateRunRows(rows)
	if res.IsFailure() {
		return res
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return result.FailSlice[PipelineRunRecord](apperror.WrapSimple(rowsErr, "iterate pipeline run rows"))
	}

	return res
}

// QueryRecentRuns retrieves recent pipeline executions.
func (p *PipelineSplitDb) QueryRecentRuns(limit int) PipelineRunSliceResult {
	limitVal := resolveLimit(limit, 10)
	rows, err := p.conn.Query(sqlQueryRecentRuns, limitVal)
	if err != nil {
		return result.FailSlice[PipelineRunRecord](apperror.WrapSimple(err, "query recent runs"))
	}

	defer rows.Close()

	return collectRecentRuns(rows)
}
