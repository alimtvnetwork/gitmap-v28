package cmdpipeline

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func syncRunJobsAndSegments(repo string, runId uint64, jobs []ghJobItem) {
	hasNoJobs := runId == 0 || len(jobs) == 0
	if hasNoJobs {
		return
	}

	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return
	}
	defer db.Close()

	jobRecs, segRecs := convertGhJobsToRecords(repo, runId, jobs)
	_ = db.RecordPipelineJobs(jobRecs)
	_ = db.RecordPipelineSegments(runId, segRecs)
}

func convertGhJobsToRecords(repo string, runId uint64, jobs []ghJobItem) ([]pipelinedb.PipelineJobRecord, []pipelinedb.PipelineSegmentRecord) {
	var jobRecs []pipelinedb.PipelineJobRecord
	var segRecs []pipelinedb.PipelineSegmentRecord

	for _, j := range jobs {
		dur := computeDurationSeconds(j.StartedAt, j.CompletedAt, j.Status)
		jobRecs = append(jobRecs, pipelinedb.PipelineJobRecord{
			JobId:           j.DatabaseId,
			RunId:           runId,
			RepoSlug:        repo,
			JobName:         j.Name,
			Status:          j.Status,
			Conclusion:      j.Conclusion,
			StartedAt:       j.StartedAt,
			CompletedAt:     j.CompletedAt,
			DurationSeconds: dur,
			JobUrl:          j.Url,
		})
		segRecs = append(segRecs, extractStepSegments(runId, j)...)
	}

	return jobRecs, segRecs
}

func extractStepSegments(runId uint64, j ghJobItem) []pipelinedb.PipelineSegmentRecord {
	var segs []pipelinedb.PipelineSegmentRecord
	for _, s := range j.Steps {
		stepDur := computeDurationSeconds(s.StartedAt, s.CompletedAt, s.Status)
		segs = append(segs, pipelinedb.PipelineSegmentRecord{
			RunId:           runId,
			JobName:         j.Name,
			StepName:        s.Name,
			StepNumber:      s.Number,
			Status:          s.Status,
			Conclusion:      s.Conclusion,
			DurationSeconds: stepDur,
		})
	}
	return segs
}

func computeDurationSeconds(startedStr, completedStr, status string) int {
	if status == "in_progress" {
		return computeInProgressDurationSec(startedStr)
	}
	hasEmptyTime := startedStr == "" || completedStr == ""
	if hasEmptyTime {
		return 0
	}
	t1, err1 := time.Parse(time.RFC3339, startedStr)
	t2, err2 := time.Parse(time.RFC3339, completedStr)
	isInvalidTime := err1 != nil || err2 != nil || t2.Before(t1)
	if isInvalidTime {
		return 0
	}
	return int(t2.Sub(t1).Seconds())
}

func computeInProgressDurationSec(startedStr string) int {
	t, err := time.Parse(time.RFC3339, startedStr)
	if err != nil {
		return 0
	}
	elapsed := int(time.Since(t).Seconds())
	if elapsed < 0 {
		return 0
	}
	return elapsed
}
