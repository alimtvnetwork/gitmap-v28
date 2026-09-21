package pipelinedb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
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
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId, StepName) DO UPDATE SET
    RepoSlug = excluded.RepoSlug,
    WorkflowName = excluded.WorkflowName,
    ErrorText = excluded.ErrorText,
    RawLogs = excluded.RawLogs,
    Notes = excluded.Notes,
    Comments = excluded.Comments,
    CreatedAt = excluded.CreatedAt;`

const sqlRecordDetailErrorLog = `
INSERT INTO PipelineDetailErrorLog (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, RawLogs, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId, StepName) DO UPDATE SET
    RepoSlug = excluded.RepoSlug,
    WorkflowName = excluded.WorkflowName,
    ErrorText = excluded.ErrorText,
    RawLogs = excluded.RawLogs,
    Notes = excluded.Notes,
    Comments = excluded.Comments,
    CreatedAt = excluded.CreatedAt;`

const sqlRecordCompactErrorLog = `
INSERT INTO PipelineCompactErrorLog (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, CompactLogs, FilteredOkCount, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId, StepName) DO UPDATE SET
    RepoSlug = excluded.RepoSlug,
    WorkflowName = excluded.WorkflowName,
    ErrorText = excluded.ErrorText,
    CompactLogs = excluded.CompactLogs,
    FilteredOkCount = excluded.FilteredOkCount,
    Notes = excluded.Notes,
    Comments = excluded.Comments,
    CreatedAt = excluded.CreatedAt;`

const sqlQueryRecentRuns = `
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun ORDER BY PipelineRunId DESC LIMIT ?;`

const sqlQueryRecentErrors = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineErrorLog ORDER BY PipelineErrorLogId DESC LIMIT ?;`

const sqlQueryRecentDetailErrors = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineDetailErrorLog ORDER BY PipelineDetailErrorLogId DESC LIMIT ?;`

const sqlQueryRecentCompactErrors = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(CompactLogs, ''), FilteredOkCount, CreatedAt
FROM PipelineCompactErrorLog ORDER BY PipelineCompactErrorLogId DESC LIMIT ?;`

const sqlQueryCachedErrorRunIds = `
SELECT DISTINCT RunId FROM (
	SELECT RunId FROM PipelineErrorLog
	UNION
	SELECT RunId FROM PipelineDetailErrorLog
	UNION
	SELECT RunId FROM PipelineCompactErrorLog
) ORDER BY RunId DESC;`

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

const sqlQueryDetailErrorLogsByRunId = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineDetailErrorLog WHERE RunId = ? ORDER BY PipelineDetailErrorLogId ASC;`

const sqlQueryCompactErrorLogsByRunId = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(CompactLogs, ''), FilteredOkCount, CreatedAt
FROM PipelineCompactErrorLog WHERE RunId = ? ORDER BY PipelineCompactErrorLogId ASC;`

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

func scanPipelineError(rows *sql.Rows) (PipelineErrorRecord, *apperror.AppError) {
	var e PipelineErrorRecord
	err := rows.Scan(&e.RunId, &e.RepoSlug, &e.WorkflowName, &e.StepName, &e.ErrorText, &e.RawLogs, &e.CreatedAt)
	if err != nil {
		return e, apperror.WrapSimple(err, "scan pipeline error log row")
	}

	return e, nil
}

func iterateErrorRows(rows *sql.Rows) PipelineErrorSliceResult {
	var list []PipelineErrorRecord
	for rows.Next() {
		e, scanErr := scanPipelineError(rows)
		if scanErr != nil {
			return result.FailSlice[PipelineErrorRecord](scanErr)
		}
		list = append(list, e)
	}

	return result.OkSlice(list)
}

func collectRecentErrors(rows *sql.Rows) PipelineErrorSliceResult {
	res := iterateErrorRows(rows)
	if res.IsFailure() {
		return res
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return result.FailSlice[PipelineErrorRecord](apperror.WrapSimple(rowsErr, "iterate pipeline error log rows"))
	}

	return res
}

// QueryRecentErrorLogs retrieves stored error diagnostics.
func (p *PipelineSplitDb) QueryRecentErrorLogs(limit int) PipelineErrorSliceResult {
	limitVal := resolveLimit(limit, 20)
	rows, err := p.conn.Query(sqlQueryRecentErrors, limitVal)
	if err != nil {
		return result.FailSlice[PipelineErrorRecord](apperror.WrapSimple(err, "query recent error logs"))
	}

	defer rows.Close()

	return collectRecentErrors(rows)
}

func scanPipelineCompactError(rows *sql.Rows) (PipelineCompactErrorRecord, *apperror.AppError) {
	var c PipelineCompactErrorRecord
	err := rows.Scan(&c.RunId, &c.RepoSlug, &c.WorkflowName, &c.StepName, &c.ErrorText, &c.CompactLogs, &c.FilteredOkCount, &c.CreatedAt)
	if err != nil {
		return c, apperror.WrapSimple(err, "scan pipeline compact error row")
	}

	return c, nil
}

func iterateCompactErrorRows(rows *sql.Rows) PipelineCompactErrorSliceResult {
	var list []PipelineCompactErrorRecord
	for rows.Next() {
		c, scanErr := scanPipelineCompactError(rows)
		if scanErr != nil {
			return result.FailSlice[PipelineCompactErrorRecord](scanErr)
		}
		list = append(list, c)
	}

	return result.OkSlice(list)
}

func collectRecentCompactErrors(rows *sql.Rows) PipelineCompactErrorSliceResult {
	res := iterateCompactErrorRows(rows)
	if res.IsFailure() {
		return res
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return result.FailSlice[PipelineCompactErrorRecord](apperror.WrapSimple(rowsErr, "iterate pipeline compact error rows"))
	}

	return res
}

// QueryDetailedErrors retrieves stored uncompressed detailed error diagnostics.
func (p *PipelineSplitDb) QueryDetailedErrors(limit int) PipelineErrorSliceResult {
	limitVal := resolveLimit(limit, 20)
	rows, err := p.conn.Query(sqlQueryRecentDetailErrors, limitVal)
	if err != nil {
		return result.FailSlice[PipelineErrorRecord](apperror.WrapSimple(err, "query recent detail error logs"))
	}

	defer rows.Close()

	return collectRecentErrors(rows)
}

// QueryCompactErrors retrieves stored noise-filtered compact error diagnostics.
func (p *PipelineSplitDb) QueryCompactErrors(limit int) PipelineCompactErrorSliceResult {
	limitVal := resolveLimit(limit, 20)
	rows, err := p.conn.Query(sqlQueryRecentCompactErrors, limitVal)
	if err != nil {
		return result.FailSlice[PipelineCompactErrorRecord](apperror.WrapSimple(err, "query recent compact error logs"))
	}

	defer rows.Close()

	return collectRecentCompactErrors(rows)
}

func clearTableQueries() []string {
	return []string{
		"DELETE FROM PipelineCompactErrorLog;",
		"DELETE FROM PipelineDetailErrorLog;",
		"DELETE FROM PipelineErrorLog;",
		"DELETE FROM PipelineSegment;",
		"DELETE FROM PipelineRun;",
	}
}

// Clear truncates all recorded runs, error logs, and segments, resets sequences, vacuums, and purges cache files.
func (p *PipelineSplitDb) Clear() error {
	runIds := p.collectAllRunIds()
	if err := p.truncateAllTables(); err != nil {
		return err
	}
	_ = p.resetSqliteSequence()
	_ = p.runVacuum()
	p.purgeCacheFiles(runIds)

	return nil
}

func (p *PipelineSplitDb) collectAllRunIds() []uint64 {
	rows, err := p.conn.Query("SELECT RunId FROM PipelineRun;")
	if err != nil {
		return nil
	}
	defer rows.Close()

	return extractRunIdsFromRows(rows)
}

func extractRunIdsFromRows(rows *sql.Rows) []uint64 {
	var ids []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}

func (p *PipelineSplitDb) truncateAllTables() error {
	for _, q := range clearTableQueries() {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "clear pipeline split db")
		}
	}

	return nil
}

func (p *PipelineSplitDb) resetSqliteSequence() error {
	query := "DELETE FROM sqlite_sequence WHERE name IN ('PipelineCompactErrorLog', 'PipelineDetailErrorLog', 'PipelineErrorLog', 'PipelineSegment', 'PipelineRun');"
	_, err := p.conn.Exec(query)
	if err != nil && !strings.Contains(err.Error(), "no such table") {
		return apperror.WrapSimple(err, "reset sqlite sequence")
	}

	return nil
}

func (p *PipelineSplitDb) runVacuum() error {
	if _, err := p.conn.Exec("VACUUM;"); err != nil {
		return apperror.WrapSimple(err, "vacuum pipeline db")
	}

	return nil
}

func (p *PipelineSplitDb) purgeCacheFiles(runIds []uint64) {
	dir := filepath.Dir(p.Path)
	purgeRunIdCacheFiles(dir, runIds)
	purgeRepoCacheDir(filepath.Join(dir, strings.ReplaceAll(p.RepoSlug, "/", "_")))
	purgeRepoCacheDir(filepath.Join(dir, SanitizeRepoSlug(p.RepoSlug)))
	purgeRepoMatchingJsonFiles(dir, p.RepoSlug)
	purgePipelineReports(dir)
	purgeParentLegacyFiles(dir, p.RepoSlug, runIds)
}

func purgePipelineReports(dir string) {
	_ = os.Remove(filepath.Join(dir, "pipeline_errors.log"))
	_ = os.Remove(filepath.Join(dir, "last_error.log"))
}

func purgeParentLegacyFiles(dir, repoSlug string, runIds []uint64) {
	parent := filepath.Dir(dir)
	if filepath.Base(parent) == "pipeline" {
		slug := SanitizeRepoSlug(repoSlug)
		_ = os.Remove(filepath.Join(parent, "pipeline_"+slug+".db"))
		_ = os.Remove(filepath.Join(parent, slug+".db"))
		purgeRunIdCacheFiles(parent, runIds)
	}
}

func purgeRunIdCacheFiles(dir string, runIds []uint64) {
	for _, id := range runIds {
		_ = os.Remove(filepath.Join(dir, fmt.Sprintf("%d.log", id)))
		_ = os.Remove(filepath.Join(dir, fmt.Sprintf("%d.json", id)))
		_ = os.Remove(filepath.Join(dir, fmt.Sprintf("%d.jobs.json", id)))
	}
}

func purgeRepoCacheDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() && isPurgeableCacheFile(e.Name()) {
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
	_ = os.Remove(dir)
}

func isPurgeableCacheFile(name string) bool {
	lower := strings.ToLower(name)

	return strings.HasSuffix(lower, ".log") || strings.HasSuffix(lower, ".json")
}

func purgeRepoMatchingJsonFiles(dir, repoSlug string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
			purgeIfRepoMatches(filepath.Join(dir, e.Name()), repoSlug)
		}
	}
}

func purgeIfRepoMatches(jsonPath, repoSlug string) {
	content, err := os.ReadFile(jsonPath)
	if err != nil || !strings.Contains(string(content), repoSlug) {
		return
	}
	_ = os.Remove(jsonPath)
	baseNoExt := strings.TrimSuffix(jsonPath, ".json")
	_ = os.Remove(baseNoExt + ".log")
	_ = os.Remove(strings.TrimSuffix(baseNoExt, ".jobs") + ".log")
}

func dropTableQueries() []string {
	return []string{
		"DROP TABLE IF EXISTS PipelineCompactErrorLog;",
		"DROP TABLE IF EXISTS PipelineDetailErrorLog;",
		"DROP TABLE IF EXISTS PipelineErrorLog;",
		"DROP TABLE IF EXISTS PipelineSegment;",
		"DROP TABLE IF EXISTS PipelineRun;",
	}
}

// Reset drops all tables and re-initializes the schema.
func (p *PipelineSplitDb) Reset() error {
	for _, q := range dropTableQueries() {
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

func (p *PipelineSplitDb) loadRunOutcomeCounts(stats *PipelineDbStats) *apperror.AppError {
	var err *apperror.AppError
	if stats.SuccessRuns, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 1;"); err != nil {
		return err
	}
	if stats.FailedRuns, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 0;"); err != nil {
		return err
	}

	return nil
}

func (p *PipelineSplitDb) loadRunStatsCounts(stats *PipelineDbStats) *apperror.AppError {
	var err *apperror.AppError
	if stats.TotalRuns, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun;"); err != nil {
		return err
	}

	return p.loadRunOutcomeCounts(stats)
}

func (p *PipelineSplitDb) loadDetailStatsCounts(stats *PipelineDbStats) *apperror.AppError {
	var err *apperror.AppError
	if stats.ErrorLogCount, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineErrorLog;"); err != nil {
		return err
	}
	if stats.SegmentCount, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineSegment;"); err != nil {
		return err
	}

	return nil
}

func (p *PipelineSplitDb) loadStatsCounts(stats *PipelineDbStats) *apperror.AppError {
	if err := p.loadRunStatsCounts(stats); err != nil {
		return err
	}

	return p.loadDetailStatsCounts(stats)
}

func (p *PipelineSplitDb) populateLastUpdated(stats *PipelineDbStats) error {
	lastUpdated, err := queryLastUpdated(p.conn)
	if err != nil {
		return err
	}
	stats.LastUpdated = lastUpdated

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

	return stats, p.populateLastUpdated(&stats)
}

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
	query := "SELECT DurationSeconds FROM PipelineRun WHERE IsSuccess = 1 AND DurationSeconds >= 10 AND (WorkflowName = ? OR ? = '') ORDER BY PipelineRunId DESC LIMIT ?;"
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

func computeRelativeDbPath(absPath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.ToSlash(absPath)
	}
	rel, relErr := filepath.Rel(cwd, absPath)
	if relErr != nil || strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(absPath)
	}

	return formatRelativePrefix(filepath.ToSlash(rel))
}

func formatRelativePrefix(slashRel string) string {
	if strings.HasPrefix(slashRel, "./") || strings.HasPrefix(slashRel, "../") {
		return slashRel
	}

	return "./" + slashRel
}

func formatUnitSize(val float64, unit string) string {
	if val == float64(int64(val)) {
		return fmt.Sprintf("%d %s", int64(val), unit)
	}

	return fmt.Sprintf("%.1f %s", val, unit)
}

// FormatHumanSize formats byte counts into human-readable strings (B, KB, MB, GB).
func FormatHumanSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return formatUnitSize(float64(bytes)/1024.0, "KB")
	}
	if bytes < 1024*1024*1024 {
		return formatUnitSize(float64(bytes)/(1024.0*1024.0), "MB")
	}

	return formatUnitSize(float64(bytes)/(1024.0*1024.0*1024.0), "GB")
}

func formatHumanSize(bytes int64) string {
	return FormatHumanSize(bytes)
}

func resolveSummaryPath(info PipelineDatabaseInfo) string {
	if info.RelativePath != "" {
		return info.RelativePath
	}

	return info.Path
}

func resolveRunCountLabel(totalRuns int) string {
	if totalRuns == 1 {
		return "run"
	}

	return "runs"
}

// FormatPipelineDbSummary formats a concise summary string for pipeline database telemetry.
func FormatPipelineDbSummary(info PipelineDatabaseInfo) string {
	pathStr := resolveSummaryPath(info)
	lbl := resolveRunCountLabel(info.TotalRuns)
	sizeStr := info.HumanSize
	if sizeStr == "" {
		sizeStr = formatHumanSize(info.SizeBytes)
	}

	return fmt.Sprintf("%s (%s, %d %s)", pathStr, sizeStr, info.TotalRuns, lbl)
}

func isRegularDbFile(fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}

	return true
}

func (p *PipelineSplitDb) populateFileMetadata(info *PipelineDatabaseInfo) {
	fi, err := os.Stat(p.Path)
	if err != nil {
		return
	}
	info.SizeBytes = fi.Size()
	info.HumanSize = formatHumanSize(fi.Size())
	info.IsExisting = isRegularDbFile(fi)
}

func (p *PipelineSplitDb) populateCounts(info *PipelineDatabaseInfo) {
	if p.conn == nil {
		return
	}
	info.TotalRuns, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun;")
	info.FailedRuns, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 0;")
	info.ErrorCount, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineErrorLog;")
	last, _ := queryLastUpdated(p.conn)
	info.LastUpdated = last
}

// GetDatabaseInfo returns diagnostic telemetry and metadata for the database.
func (p *PipelineSplitDb) GetDatabaseInfo() PipelineDatabaseInfo {
	var info PipelineDatabaseInfo
	info.Path = filepath.ToSlash(p.Path)
	info.RelativePath = filepath.ToSlash(computeRelativeDbPath(p.Path))
	p.populateFileMetadata(&info)
	p.populateCounts(&info)

	return info
}

// Vacuum executes SQLite VACUUM and returns reclaimed bytes.
func (p *PipelineSplitDb) Vacuum() (int64, error) {
	if p.conn == nil {
		return 0, nil
	}
	beforeSize := resolveFileSize(p.Path)
	if _, err := p.conn.Exec("VACUUM;"); err != nil {
		return 0, apperror.WrapSimple(err, "vacuum pipeline db")
	}

	return calculateFreedBytes(beforeSize, resolveFileSize(p.Path)), nil
}

func resolveFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return fi.Size()
}

func calculateFreedBytes(before, after int64) int64 {
	if before > after {
		return before - after
	}

	return 0
}
