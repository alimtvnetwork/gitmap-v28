package pipelinedb

import (
	"database/sql"
	"math"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// QueryStageSummary computes stage duration aggregations and speedup factor.
func (p *PipelineSplitDb) QueryStageSummary(runId uint64) (*PipelineStageSummary, error) {
	run, err := p.QueryRunById(runId)
	if err != nil {
		return nil, err
	}
	jobs, jobsErr := p.QueryPipelineJobs(runId)
	if jobsErr != nil {
		return nil, jobsErr
	}
	stageSum := calculateJobsTotalDuration(jobs)
	speedup := calculateSpeedupMultiplier(stageSum, run.DurationSeconds)

	return &PipelineStageSummary{
		RunId:             runId,
		RepoSlug:          run.RepoSlug,
		WorkflowName:      run.WorkflowName,
		Status:            run.Status,
		Conclusion:        run.Conclusion,
		WallClockSeconds:  run.DurationSeconds,
		StageSumSeconds:   stageSum,
		SpeedupMultiplier: speedup,
		Jobs:              jobs,
	}, nil
}

func calculateJobsTotalDuration(jobs []PipelineJobRecord) int {
	total := 0
	for _, j := range jobs {
		total += j.DurationSeconds
	}
	return total
}

func calculateSpeedupMultiplier(stageSum, wallClock int) float64 {
	hasZeroWallClock := wallClock <= 0 || stageSum <= 0
	if hasZeroWallClock {
		return 1.0
	}
	val := float64(stageSum) / float64(wallClock)
	return math.Round(val*100) / 100
}

// QueryPipelineSegments returns all steps/segments recorded for a workflow run ID.
func (p *PipelineSplitDb) QueryPipelineSegments(runId uint64) ([]PipelineSegmentRecord, error) {
	query := `SELECT PipelineSegmentId, RunId, JobName, StepName, StepNumber,
		Status, Conclusion, DurationSeconds, CreatedAt
		FROM PipelineSegment WHERE RunId = ? ORDER BY StepNumber ASC, PipelineSegmentId ASC;`
	rows, err := p.conn.Query(query, runId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query pipeline segments")
	}
	defer rows.Close()
	return scanPipelineSegments(rows)
}

func scanPipelineSegments(rows *sql.Rows) ([]PipelineSegmentRecord, error) {
	var list []PipelineSegmentRecord
	for rows.Next() {
		var s PipelineSegmentRecord
		if err := rows.Scan(
			&s.SegmentId, &s.RunId, &s.JobName, &s.StepName, &s.StepNumber,
			&s.Status, &s.Conclusion, &s.DurationSeconds, &s.CreatedAt,
		); err != nil {
			return nil, apperror.WrapSimple(err, "scan pipeline segment")
		}
		list = append(list, s)
	}
	return list, nil
}
