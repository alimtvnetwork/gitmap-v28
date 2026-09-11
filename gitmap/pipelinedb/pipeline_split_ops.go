package pipelinedb

import (
	"database/sql"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const sqlRecordRun = `
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

const sqlRecordErrorLog = `
INSERT INTO PipelineErrorLog (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, RawLogs, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

const sqlQueryRecentRuns = `
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun ORDER BY PipelineRunId DESC LIMIT ?;`

const sqlQueryRecentErrors = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineErrorLog ORDER BY PipelineErrorLogId DESC LIMIT ?;`

const sqlQueryCachedErrorRunIds = `
SELECT DISTINCT RunId FROM PipelineErrorLog ORDER BY RunId DESC;`

const sqlQueryRunByOffset = `
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun ORDER BY PipelineRunId DESC LIMIT 1 OFFSET ?;`

const sqlQueryLastFailedRuns = `
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun WHERE IsSuccess = 0 ORDER BY PipelineRunId DESC LIMIT ?;`

const sqlQueryErrorLogsByRunId = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineErrorLog WHERE RunId = ? ORDER BY PipelineErrorLogId ASC;`

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
	err := p.conn.QueryRow("SELECT 1 FROM PipelineErrorLog WHERE RunId = ? LIMIT 1;", runId).Scan(&exists)
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

func collectRecentRuns(rows *sql.Rows) ([]PipelineRunRecord, error) {
	var list []PipelineRunRecord
	for rows.Next() {
		r, scanErr := scanPipelineRun(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		list = append(list, r)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "iterate pipeline run rows")
	}
	return list, nil
}

// QueryRecentRuns retrieves recent pipeline executions.
func (p *PipelineSplitDb) QueryRecentRuns(limit int) ([]PipelineRunRecord, error) {
	limitVal := resolveLimit(limit, 10)
	rows, err := p.conn.Query(sqlQueryRecentRuns, limitVal)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query recent runs")
	}
	defer rows.Close()
	return collectRecentRuns(rows)
}

func scanPipelineError(rows *sql.Rows) (PipelineErrorRecord, *apperror.AppError) {
	var e PipelineErrorRecord
	err := rows.Scan(&e.RunId, &e.RepoSlug, &e.WorkflowName, &e.StepName, &e.ErrorText, &e.RawLogs, &e.CreatedAt)
	if err != nil {
		return e, apperror.WrapSimple(err, "scan pipeline error log row")
	}
	return e, nil
}

func collectRecentErrors(rows *sql.Rows) ([]PipelineErrorRecord, error) {
	var list []PipelineErrorRecord
	for rows.Next() {
		e, scanErr := scanPipelineError(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		list = append(list, e)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "iterate pipeline error log rows")
	}
	return list, nil
}

// QueryRecentErrorLogs retrieves stored error diagnostics.
func (p *PipelineSplitDb) QueryRecentErrorLogs(limit int) ([]PipelineErrorRecord, error) {
	limitVal := resolveLimit(limit, 20)
	rows, err := p.conn.Query(sqlQueryRecentErrors, limitVal)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query recent error logs")
	}
	defer rows.Close()
	return collectRecentErrors(rows)
}

// Clear truncates all recorded runs, error logs, and segments.
func (p *PipelineSplitDb) Clear() error {
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
func (p *PipelineSplitDb) Reset() error {
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

func (p *PipelineSplitDb) optimizePragmas() *apperror.AppError {
	if _, err := p.conn.Exec("PRAGMA wal_checkpoint(TRUNCATE);"); err != nil {
		return apperror.WrapSimple(err, "wal checkpoint pipeline db")
	}
	if _, err := p.conn.Exec("VACUUM;"); err != nil {
		return apperror.WrapSimple(err, "vacuum pipeline db")
	}
	if _, err := p.conn.Exec("PRAGMA optimize;"); err != nil {
		return apperror.WrapSimple(err, "optimize pipeline db")
	}
	return nil
}

// Optimize executes WAL checkpoint and VACUUM, returning reclaimed bytes.
func (p *PipelineSplitDb) Optimize() (int64, error) {
	sizeBefore := getFileSize(p.Path)
	if err := p.optimizePragmas(); err != nil {
		return 0, err
	}

	sizeAfter := getFileSize(p.Path)
	if sizeBefore <= sizeAfter {
		return 0, nil
	}

	return sizeBefore - sizeAfter, nil
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	if info.Size() < 0 {
		return 0
	}

	return info.Size()
}

func safeInt64ToUint64(val int64) uint64 {
	if val < 0 {
		return 0
	}

	return uint64(val)
}

func countQuery(conn *sql.DB, query string) (int, *apperror.AppError) {
	var count int
	if err := conn.QueryRow(query).Scan(&count); err != nil {
		return 0, apperror.WrapSimple(err, "count query: "+query)
	}

	return count, nil
}

func queryLastUpdated(conn *sql.DB) (string, *apperror.AppError) {
	var lastUpdated string
	query := "SELECT COALESCE(MAX(UpdatedAt), '') FROM PipelineRun;"
	if err := conn.QueryRow(query).Scan(&lastUpdated); err != nil {
		return "", apperror.WrapSimple(err, "query last updated")
	}

	return lastUpdated, nil
}

func (p *PipelineSplitDb) loadRunStatsCounts(stats *PipelineDbStats) *apperror.AppError {
	var err *apperror.AppError
	if stats.TotalRuns, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun;"); err != nil {
		return err
	}
	if stats.SuccessRuns, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 1;"); err != nil {
		return err
	}
	if stats.FailedRuns, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 0;"); err != nil {
		return err
	}

	return nil
}

func (p *PipelineSplitDb) loadStatsCounts(stats *PipelineDbStats) *apperror.AppError {
	if err := p.loadRunStatsCounts(stats); err != nil {
		return err
	}
	var err *apperror.AppError
	if stats.ErrorLogCount, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineErrorLog;"); err != nil {
		return err
	}
	if stats.SegmentCount, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineSegment;"); err != nil {
		return err
	}

	return nil
}

// GetStats returns telemetry metrics for the split database.
func (p *PipelineSplitDb) GetStats() (PipelineDbStats, error) {
	var stats PipelineDbStats
	stats.Path = p.Path
	stats.Size = safeInt64ToUint64(getFileSize(p.Path))
	if err := p.loadStatsCounts(&stats); err != nil {
		return stats, err
	}
	lastUpdated, err := queryLastUpdated(p.conn)
	if err != nil {
		return stats, err
	}
	stats.LastUpdated = lastUpdated

	return stats, nil
}

func collectRunIdList(rows *sql.Rows) ([]uint64, error) {
	var list []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, apperror.WrapSimple(err, "scan cached run id")
		}
		list = append(list, id)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "iterate cached run ids")
	}

	return list, nil
}

// QueryCachedErrorRunIds retrieves distinct RunIds cached in PipelineErrorLog.
func (p *PipelineSplitDb) QueryCachedErrorRunIds() ([]uint64, error) {
	rows, err := p.conn.Query(sqlQueryCachedErrorRunIds)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query cached error run ids")
	}
	defer rows.Close()

	return collectRunIdList(rows)
}

// QueryCachedErrorRunIdMap returns a map set of cached error RunIds for fast lookups.
func (p *PipelineSplitDb) QueryCachedErrorRunIdMap() (map[uint64]bool, error) {
	ids, err := p.QueryCachedErrorRunIds()
	if err != nil {
		return nil, err
	}
	idMap := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		idMap[id] = true
	}

	return idMap, nil
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
	sqlOffset := normalizeNegativeOffset(offset)
	rows, err := p.conn.Query(sqlQueryRunByOffset, sqlOffset)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query run by offset")
	}
	defer rows.Close()
	runs, err := collectRecentRuns(rows)
	if err != nil || len(runs) == 0 {
		return nil, err
	}

	return &runs[0], nil
}

// QueryRunsByNegativeOffset is an alias for QueryRunByNegativeOffset.
func (p *PipelineSplitDb) QueryRunsByNegativeOffset(offset int) (*PipelineRunRecord, error) {
	return p.QueryRunByNegativeOffset(offset)
}

// QueryLastFailedRuns retrieves the most recent failed pipeline runs up to limit.
func (p *PipelineSplitDb) QueryLastFailedRuns(limit int) ([]PipelineRunRecord, error) {
	limitVal := resolveLimit(limit, 5)
	rows, err := p.conn.Query(sqlQueryLastFailedRuns, limitVal)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query last failed runs")
	}
	defer rows.Close()

	return collectRecentRuns(rows)
}

// QueryLastNFailedRuns is an alias for QueryLastFailedRuns.
func (p *PipelineSplitDb) QueryLastNFailedRuns(limit int) ([]PipelineRunRecord, error) {
	return p.QueryLastFailedRuns(limit)
}

// QueryErrorLogsByRunId retrieves all error diagnostics recorded for a specific run ID.
func (p *PipelineSplitDb) QueryErrorLogsByRunId(runId uint64) ([]PipelineErrorRecord, error) {
	rows, err := p.conn.Query(sqlQueryErrorLogsByRunId, runId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query error logs by run id")
	}
	defer rows.Close()

	return collectRecentErrors(rows)
}
