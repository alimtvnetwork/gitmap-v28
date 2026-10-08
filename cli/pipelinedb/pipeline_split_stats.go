package pipelinedb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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
	if stats.JobCount, err = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineJob;"); err != nil {
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
