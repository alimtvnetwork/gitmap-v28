package cmdpipeline

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// FindPreviousFailingRunInHistory locates the most recent failing run in history.
func FindPreviousFailingRunInHistory(runs []ghRunItem) (ghRunItem, bool) {
	for _, r := range runs {
		if r.Conclusion == "failure" {
			return r, true
		}
	}

	return ghRunItem{}, false
}

// TryRetrievePreviousDbRun retrieves the most recent failed run from pipelinedb.
func TryRetrievePreviousDbRun(repo string) (*pipelinedb.PipelineRunRecord, string, bool) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, "", false
	}
	defer db.Close()

	runRes := db.QueryLastFailedRuns(1)
	if runRes.IsFailed() || len(runRes.Data) == 0 {
		return nil, "", false
	}

	rec := runRes.Data[0]
	errLogs := queryDbRunErrorText(db, rec.RunId)

	return &rec, errLogs, true
}

func queryDbRunErrorText(db *pipelinedb.PipelineSplitDb, runId uint64) string {
	errRes := db.QueryDetailedErrorLogsByRunId(runId)
	if !errRes.IsFailed() && len(errRes.Data) > 0 {
		return formatDbErrorRecords(errRes.Data)
	}

	compactRes := db.QueryCompactErrors(10)
	if !compactRes.IsFailed() && len(compactRes.Data) > 0 {
		return formatDbCompactRecords(compactRes.Data)
	}

	return ""
}

func formatDbErrorRecords(records []pipelinedb.PipelineErrorRecord) string {
	var sb strings.Builder
	for _, r := range records {
		sb.WriteString(fmt.Sprintf("[%s > %s]\n", r.WorkflowName, r.StepName))
		if len(r.ErrorText) > 0 {
			sb.WriteString(r.ErrorText)
			sb.WriteString("\n")
		}
	}

	return strings.TrimSpace(sb.String())
}

func formatDbCompactRecords(records []pipelinedb.PipelineCompactErrorRecord) string {
	var sb strings.Builder
	for _, r := range records {
		sb.WriteString(fmt.Sprintf("[%s > %s (compact)]\n", r.WorkflowName, r.StepName))
		if len(r.ErrorText) > 0 {
			sb.WriteString(r.ErrorText)
			sb.WriteString("\n")
		}
	}

	return strings.TrimSpace(sb.String())
}

// ApplyPreviousRunFallbackToPayload applies previous run from history or pipeline DB.
func ApplyPreviousRunFallbackToPayload(
	p *PipelineErrorLogsPayload,
	repo string,
	runs []ghRunItem,
) bool {
	if prevRun, isHistoryFound := FindPreviousFailingRunInHistory(runs); isHistoryFound {
		populateFailedRunsPayload(repo, []ghRunItem{prevRun}, p)
		p.Notes = fmt.Sprintf("Retrieved from previous failed run #%d", prevRun.DatabaseId)

		return true
	}

	return ApplyPreviousDbFallbackToPayload(p, repo)
}

// ApplyPreviousDbFallbackToPayload populates payload from the previous DB run.
func ApplyPreviousDbFallbackToPayload(p *PipelineErrorLogsPayload, repo string) bool {
	rec, errLogs, isDbFound := TryRetrievePreviousDbRun(repo)
	if !isDbFound || rec == nil {
		return false
	}

	p.RunId = rec.RunId
	p.Conclusion = "failure"
	p.WorkflowName = rec.WorkflowName
	p.Branch = rec.Branch
	p.Sha = rec.Sha
	p.Url = rec.RunUrl
	p.ErrorLogs = errLogs
	p.Notes = fmt.Sprintf("Retrieved from previous pipeline DB run #%d", rec.RunId)

	return true
}

// ResolveFallbackTargetGroup finds the previous commit group when latest is empty.
func ResolveFallbackTargetGroup(groups []CommitPipelineGroup) (*CommitPipelineGroup, bool) {
	if len(groups) > 1 {
		return &groups[1], true
	}

	return nil, false
}
