package pipelinedb

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
