package pipelinedb

import (
	"database/sql"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"strings"
)

func clearTableQueries() []string {
	return []string{
		"DELETE FROM PipelineCompactErrorLog;",
		"DELETE FROM PipelineDetailErrorLog;",
		"DELETE FROM PipelineErrorLog;",
		"DELETE FROM PipelineSegment;",
		"DELETE FROM PipelineJob;",
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
	query := "DELETE FROM sqlite_sequence WHERE name IN ('PipelineCompactErrorLog', 'PipelineDetailErrorLog', 'PipelineErrorLog', 'PipelineSegment', 'PipelineJob', 'PipelineRun');"
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
