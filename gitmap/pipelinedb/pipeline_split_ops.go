package pipelinedb

import (
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RecordRun inserts or updates a pipeline run execution record.
func (p *PipelineSplitDB) RecordRun(r PipelineRunRecord) error {
	isSuccessInt := 0
	if r.IsSuccess || r.Conclusion == "success" {
		isSuccessInt = 1
	}
	query := `
INSERT INTO PipelineRun (
    RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
    EtaSeconds, DurationSeconds, RunUrl, IsSuccess, Notes, Comments, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId) DO UPDATE SET
    Status = excluded.Status,
    Conclusion = excluded.Conclusion,
    EtaSeconds = excluded.EtaSeconds,
    DurationSeconds = excluded.DurationSeconds,
    IsSuccess = excluded.IsSuccess,
    UpdatedAt = excluded.UpdatedAt;`

	_, err := p.conn.Exec(query,
		r.RunID, r.RepoSlug, r.WorkflowName, r.Status, r.Conclusion, r.Branch, r.Sha,
		r.EtaSeconds, r.DurationSeconds, r.RunURL, isSuccessInt, r.Notes, r.Comments, r.CreatedAt, r.UpdatedAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record pipeline run")
	}
	return nil
}

// RecordErrorLog inserts an error diagnostic entry for a failing run.
func (p *PipelineSplitDB) RecordErrorLog(e PipelineErrorRecord) error {
	query := `
INSERT INTO PipelineErrorLog (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, RawLogs, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	createdAt := e.CreatedAt
	if createdAt == "" {
		createdAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := p.conn.Exec(query,
		e.RunID, e.RepoSlug, e.WorkflowName, e.StepName, e.ErrorText, e.RawLogs, e.Notes, e.Comments, createdAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "record pipeline error log")
	}
	return nil
}

// QueryRecentRuns retrieves recent pipeline executions.
func (p *PipelineSplitDB) QueryRecentRuns(limit int) ([]PipelineRunRecord, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := p.conn.Query(`
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun ORDER BY PipelineRunId DESC LIMIT ?;`, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query recent runs")
	}
	defer rows.Close()

	var list []PipelineRunRecord
	for rows.Next() {
		var r PipelineRunRecord
		var isSuccessInt int
		if scanErr := rows.Scan(
			&r.RunID, &r.RepoSlug, &r.WorkflowName, &r.Status, &r.Conclusion, &r.Branch, &r.Sha,
			&r.EtaSeconds, &r.DurationSeconds, &r.RunURL, &isSuccessInt, &r.CreatedAt, &r.UpdatedAt,
		); scanErr == nil {
			r.IsSuccess = isSuccessInt == 1
			list = append(list, r)
		}
	}
	return list, nil
}

// QueryRecentErrorLogs retrieves stored error diagnostics.
func (p *PipelineSplitDB) QueryRecentErrorLogs(limit int) ([]PipelineErrorRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := p.conn.Query(`
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineErrorLog ORDER BY PipelineErrorLogId DESC LIMIT ?;`, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query recent error logs")
	}
	defer rows.Close()

	var list []PipelineErrorRecord
	for rows.Next() {
		var e PipelineErrorRecord
		if scanErr := rows.Scan(
			&e.RunID, &e.RepoSlug, &e.WorkflowName, &e.StepName, &e.ErrorText, &e.RawLogs, &e.CreatedAt,
		); scanErr == nil {
			list = append(list, e)
		}
	}
	return list, nil
}

// Clear truncates all recorded runs, error logs, and segments.
func (p *PipelineSplitDB) Clear() error {
	queries := []string{
		"DELETE FROM PipelineRun;",
		"DELETE FROM PipelineErrorLog;",
		"DELETE FROM PipelineSegment;",
	}
	for _, q := range queries {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "clear pipeline split db")
		}
	}
	return nil
}

// Reset drops all tables and re-initializes the schema.
func (p *PipelineSplitDB) Reset() error {
	queries := []string{
		"DROP TABLE IF EXISTS PipelineRun;",
		"DROP TABLE IF EXISTS PipelineErrorLog;",
		"DROP TABLE IF EXISTS PipelineSegment;",
	}
	for _, q := range queries {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "reset pipeline split db")
		}
	}
	return p.InitSchema()
}

// Optimize executes WAL checkpoint and VACUUM, returning reclaimed bytes.
func (p *PipelineSplitDB) Optimize() (int64, error) {
	var sizeBefore int64
	if info, err := os.Stat(p.Path); err == nil {
		sizeBefore = info.Size()
	}
	_, _ = p.conn.Exec("PRAGMA wal_checkpoint(TRUNCATE);")
	if _, err := p.conn.Exec("VACUUM;"); err != nil {
		return 0, apperror.WrapSimple(err, "vacuum pipeline db")
	}
	_, _ = p.conn.Exec("PRAGMA optimize;")
	var sizeAfter int64
	if info, err := os.Stat(p.Path); err == nil {
		sizeAfter = info.Size()
	}
	reclaimed := sizeBefore - sizeAfter
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}

// GetStats returns telemetry metrics for the split database.
func (p *PipelineSplitDB) GetStats() (PipelineDBStats, error) {
	var stats PipelineDBStats
	stats.Path = p.Path
	if info, err := os.Stat(p.Path); err == nil {
		stats.Size = info.Size()
	}
	_ = p.conn.QueryRow("SELECT COUNT(*) FROM PipelineRun;").Scan(&stats.TotalRuns)
	_ = p.conn.QueryRow("SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 1;").Scan(&stats.SuccessRuns)
	_ = p.conn.QueryRow("SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 0;").Scan(&stats.FailedRuns)
	_ = p.conn.QueryRow("SELECT COUNT(*) FROM PipelineErrorLog;").Scan(&stats.ErrorLogCount)
	_ = p.conn.QueryRow("SELECT COUNT(*) FROM PipelineSegment;").Scan(&stats.SegmentCount)
	_ = p.conn.QueryRow("SELECT COALESCE(MAX(UpdatedAt), '') FROM PipelineRun;").Scan(&stats.LastUpdated)
	return stats, nil
}
