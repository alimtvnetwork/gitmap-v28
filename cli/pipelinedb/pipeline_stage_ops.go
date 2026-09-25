package pipelinedb

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const sqlUpsertJob = `
INSERT INTO PipelineJob (
    JobId, RunId, RepoSlug, JobName, Status, Conclusion,
    StartedAt, CompletedAt, DurationSeconds, JobUrl, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(JobId) DO UPDATE SET
    Status = excluded.Status,
    Conclusion = excluded.Conclusion,
    StartedAt = excluded.StartedAt,
    CompletedAt = excluded.CompletedAt,
    DurationSeconds = excluded.DurationSeconds,
    JobUrl = excluded.JobUrl;`

const sqlInsertSegment = `
INSERT INTO PipelineSegment (
    RunId, JobName, StepName, StepNumber, Status, Conclusion, DurationSeconds, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);`

// RecordPipelineJobs inserts or updates a list of pipeline jobs.
func (p *PipelineSplitDb) RecordPipelineJobs(jobs []PipelineJobRecord) error {
	for _, j := range jobs {
		if err := p.recordSingleJob(j); err != nil {
			return err
		}
	}
	return nil
}

func (p *PipelineSplitDb) recordSingleJob(j PipelineJobRecord) error {
	_, err := p.conn.Exec(sqlUpsertJob,
		j.JobId, j.RunId, j.RepoSlug, j.JobName, j.Status, j.Conclusion,
		j.StartedAt, j.CompletedAt, j.DurationSeconds, j.JobUrl,
	)
	if err != nil {
		return apperror.WrapSimple(err, "upsert pipeline job")
	}
	return nil
}

// RecordPipelineSegments deletes prior segments for the run and inserts the given slice.
func (p *PipelineSplitDb) RecordPipelineSegments(runId uint64, segs []PipelineSegmentRecord) error {
	if _, err := p.conn.Exec("DELETE FROM PipelineSegment WHERE RunId = ?;", runId); err != nil {
		return apperror.WrapSimple(err, "clear old segments")
	}
	for _, s := range segs {
		if err := p.recordSingleSegment(s); err != nil {
			return err
		}
	}
	return nil
}

func (p *PipelineSplitDb) recordSingleSegment(s PipelineSegmentRecord) error {
	_, err := p.conn.Exec(sqlInsertSegment,
		s.RunId, s.JobName, s.StepName, s.StepNumber, s.Status, s.Conclusion, s.DurationSeconds,
	)
	if err != nil {
		return apperror.WrapSimple(err, "insert pipeline segment")
	}
	return nil
}

// QueryPipelineJobs returns all jobs recorded for a given workflow run ID.
func (p *PipelineSplitDb) QueryPipelineJobs(runId uint64) ([]PipelineJobRecord, error) {
	query := `SELECT JobId, RunId, RepoSlug, JobName, Status, Conclusion,
		StartedAt, CompletedAt, DurationSeconds, JobUrl, CreatedAt
		FROM PipelineJob WHERE RunId = ? ORDER BY StartedAt ASC, JobId ASC;`
	rows, err := p.conn.Query(query, runId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query pipeline jobs")
	}
	defer rows.Close()
	return scanPipelineJobs(rows)
}

func scanPipelineJobs(rows *sql.Rows) ([]PipelineJobRecord, error) {
	var list []PipelineJobRecord
	for rows.Next() {
		var j PipelineJobRecord
		if err := rows.Scan(
			&j.JobId, &j.RunId, &j.RepoSlug, &j.JobName, &j.Status, &j.Conclusion,
			&j.StartedAt, &j.CompletedAt, &j.DurationSeconds, &j.JobUrl, &j.CreatedAt,
		); err != nil {
			return nil, apperror.WrapSimple(err, "scan pipeline job")
		}
		list = append(list, j)
	}
	return list, nil
}
