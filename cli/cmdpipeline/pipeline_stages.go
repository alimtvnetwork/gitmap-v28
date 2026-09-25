package cmdpipeline

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func handlePipelineStages(args []string) error {
	repo, runId := resolveTargetRepoAndRunId(args)
	hasNoRun := runId == 0
	if hasNoRun {
		return fmt.Errorf("no workflow runs found for repository: %s", repo)
	}

	summary, err := loadStageSummary(repo, runId)
	if err != nil {
		return err
	}
	if hasArgFlag(args, "--json") {
		return renderStagesJSON(summary)
	}
	renderStagesTerminal(summary)
	return nil
}

func loadStageSummary(repo string, runId uint64) (*pipelinedb.PipelineStageSummary, error) {
	jobs := queryRunJobs(repo, runId)
	syncRunJobsAndSegments(repo, runId, jobs)

	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return buildFallbackStageSummary(repo, runId, jobs), nil
	}
	defer db.Close()

	summary, summaryErr := db.QueryStageSummary(runId)
	if summaryErr != nil || len(summary.Jobs) == 0 {
		return buildFallbackStageSummary(repo, runId, jobs), nil
	}
	return summary, nil
}

func buildFallbackStageSummary(repo string, runId uint64, jobs []ghJobItem) *pipelinedb.PipelineStageSummary {
	jobRecs, _ := convertGhJobsToRecords(repo, runId, jobs)
	sumSecs := 0
	status := "completed"
	conclusion := "success"
	for _, j := range jobRecs {
		sumSecs += j.DurationSeconds
		if j.Status == "in_progress" {
			status = "in_progress"
		}
		if j.Conclusion == "failure" || j.Conclusion == "cancelled" {
			conclusion = j.Conclusion
		}
	}
	return &pipelinedb.PipelineStageSummary{
		RunId:             runId,
		RepoSlug:          repo,
		WorkflowName:      "Workflow",
		Status:            status,
		Conclusion:        conclusion,
		StageSumSeconds:   sumSecs,
		WallClockSeconds:  sumSecs,
		SpeedupMultiplier: 1.0,
		Jobs:              jobRecs,
	}
}
