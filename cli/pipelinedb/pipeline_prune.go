package pipelinedb

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DiskSizeBytes calculates the total on-disk size of the SQLite DB including WAL and SHM files.
func (p *PipelineSplitDb) DiskSizeBytes() int64 {
	var total int64
	targets := []string{p.Path, p.Path + "-wal", p.Path + "-shm"}
	for _, target := range targets {
		if fi, err := os.Stat(target); err == nil && !fi.IsDir() {
			total += fi.Size()
		}
	}

	return total
}

// ResolveConfiguredPipelineMaxDbBytes returns the maximum DB size ceiling in bytes (default 10 MB).
func ResolveConfiguredPipelineMaxDbBytes() int64 {
	const defaultMaxBytes = int64(10 * 1024 * 1024)
	envVal := os.Getenv("GITMAP_PIPELINE_MAX_DB_MB")
	if len(envVal) == 0 {
		return defaultMaxBytes
	}

	mb, err := strconv.ParseInt(envVal, 10, 64)
	if err == nil && mb > 0 {
		return mb * 1024 * 1024
	}

	return defaultMaxBytes
}

// PruneIfExceedsSize prunes old runs, error logs, and vacuums if the DB exceeds the maximum byte limit.
func (p *PipelineSplitDb) PruneIfExceedsSize(maxBytes int64) (bool, int, error) {
	limit := maxBytes
	if limit <= 0 {
		limit = ResolveConfiguredPipelineMaxDbBytes()
	}
	if p.DiskSizeBytes() <= limit {
		return false, 0, nil
	}

	return p.executeIterativePrune(limit)
}

func (p *PipelineSplitDb) executeIterativePrune(limit int64) (bool, int, error) {
	totalPruned := 0
	for iteration := 0; iteration < 5; iteration++ {
		pruned, err := p.pruneOldestBatch()
		if err != nil {
			return totalPruned > 0, totalPruned, err
		}
		totalPruned += pruned
		_ = p.truncateOlderRawLogs(5)
		_ = p.runCheckpointAndVacuum()
		if p.DiskSizeBytes() <= limit || pruned == 0 {
			break
		}
	}

	return totalPruned > 0, totalPruned, nil
}

func (p *PipelineSplitDb) pruneOldestBatch() (int, error) {
	runIds := p.queryOldestRunIds(10)
	if len(runIds) == 0 {
		return 0, nil
	}

	p.purgeCacheFiles(runIds)

	return len(runIds), p.deleteRunIdsCascade(runIds)
}

func (p *PipelineSplitDb) queryOldestRunIds(limit int) []uint64 {
	query := "SELECT RunId FROM PipelineRun ORDER BY PipelineRunId ASC LIMIT ?;"
	rows, err := p.conn.Query(query, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()

	return extractRunIdsFromRows(rows)
}

func (p *PipelineSplitDb) deleteRunIdsCascade(runIds []uint64) error {
	for _, id := range runIds {
		_, _ = p.conn.Exec("DELETE FROM PipelineCompactErrorLog WHERE RunId = ?;", id)
		_, _ = p.conn.Exec("DELETE FROM PipelineDetailErrorLog WHERE RunId = ?;", id)
		_, _ = p.conn.Exec("DELETE FROM PipelineErrorLog WHERE RunId = ?;", id)
		_, _ = p.conn.Exec("DELETE FROM PipelineSegment WHERE RunId = ?;", id)
		_, _ = p.conn.Exec("DELETE FROM PipelineRun WHERE RunId = ?;", id)
	}

	return nil
}

func (p *PipelineSplitDb) truncateOlderRawLogs(keepRecent int) error {
	q1 := "UPDATE PipelineErrorLog SET RawLogs = '' WHERE PipelineErrorLogId NOT IN (SELECT PipelineErrorLogId FROM PipelineErrorLog ORDER BY PipelineErrorLogId DESC LIMIT ?);"
	q2 := "UPDATE PipelineDetailErrorLog SET RawLogs = '' WHERE PipelineDetailErrorLogId NOT IN (SELECT PipelineDetailErrorLogId FROM PipelineDetailErrorLog ORDER BY PipelineDetailErrorLogId DESC LIMIT ?);"
	_, _ = p.conn.Exec(q1, keepRecent)
	_, _ = p.conn.Exec(q2, keepRecent)

	return nil
}

func (p *PipelineSplitDb) runCheckpointAndVacuum() error {
	_, _ = p.conn.Exec("PRAGMA wal_checkpoint(TRUNCATE);")
	_, err := p.conn.Exec("VACUUM;")
	_, _ = p.conn.Exec("PRAGMA wal_checkpoint(TRUNCATE);")

	return err
}

// HasErrorLogForSha checks whether an error log is associated with a given commit SHA.
func (p *PipelineSplitDb) HasErrorLogForSha(sha string) bool {
	if len(sha) < 7 {
		return false
	}
	var exists int
	query := "SELECT 1 FROM PipelineErrorLog e JOIN PipelineRun r ON e.RunId = r.RunId WHERE r.Sha LIKE ? LIMIT 1;"
	err := p.conn.QueryRow(query, sha+"%").Scan(&exists)
	if err != nil {
		return false
	}

	return exists == 1
}

// HasCompletedRun checks if a completed run exists for a given runId or commit SHA.
func (p *PipelineSplitDb) HasCompletedRun(runId uint64, sha string) bool {
	var exists int
	query := "SELECT 1 FROM PipelineRun WHERE (RunId = ? OR (length(?) >= 7 AND Sha LIKE ?)) AND (Status = 'completed' OR Conclusion IN ('success', 'failure')) LIMIT 1;"
	err := p.conn.QueryRow(query, runId, sha, sha+"%").Scan(&exists)
	if err != nil {
		return false
	}

	return exists == 1
}

// QueryRunBySha retrieves a pipeline run record by its commit SHA prefix.
func (p *PipelineSplitDb) QueryRunBySha(sha string) (*PipelineRunRecord, error) {
	if len(sha) < 7 {
		return nil, apperror.NewValidationError("sha too short for lookup")
	}
	query := "SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha, EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt FROM PipelineRun WHERE Sha LIKE ? ORDER BY PipelineRunId DESC LIMIT 1;"
	rows, err := p.conn.Query(query, sha+"%")
	if err != nil {
		return nil, apperror.WrapSimple(err, "query run by sha")
	}
	defer rows.Close()

	return extractFirstRunResult(rows)
}

// QueryRunByRunId retrieves a pipeline run record by its GitHub run ID.
func (p *PipelineSplitDb) QueryRunByRunId(runId uint64) (*PipelineRunRecord, error) {
	query := "SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha, EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt FROM PipelineRun WHERE RunId = ? LIMIT 1;"
	rows, err := p.conn.Query(query, runId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query run by runId")
	}
	defer rows.Close()

	return extractFirstRunResult(rows)
}

func extractFirstRunResult(rows *sql.Rows) (*PipelineRunRecord, error) {
	runsRes := collectRecentRuns(rows)
	if runsRes.IsFailure() || runsRes.IsEmpty() {
		return nil, runsRes.AppError()
	}

	return &runsRes.Data[0], nil
}
