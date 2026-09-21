package cmdpipeline

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// BuildFailedChecksTree generates a tree view of failed workflows and jobs for a commit.
func BuildFailedChecksTree(p PipelineErrorLogsPayload, isColor bool) string {
	if len(p.FailedRuns) == 0 {
		return ""
	}

	var sb strings.Builder
	writeTreeHeader(&sb, p, isColor)
	for i, fr := range p.FailedRuns {
		isLast := i == len(p.FailedRuns)-1
		renderTreeWorkflowNode(&sb, fr, isLast, isColor)
	}

	return sb.String()
}

func writeTreeHeader(sb *strings.Builder, p PipelineErrorLogsPayload, isColor bool) {
	writeTreeHeaderSha(sb, resolveTreeCommitLabel(p), isColor)
}

func resolveTreeCommitLabel(p PipelineErrorLogsPayload) string {
	sha := p.LastHash
	if len(sha) == 0 {
		sha = p.Sha
	}
	if len(sha) > 0 {
		return sha
	}

	return "current"
}

func writeTreeHeaderSha(sb *strings.Builder, fullSha string, isColor bool) {
	sha := fullSha
	if len(sha) > 7 {
		sha = sha[:7]
	}
	if isColor {
		fmt.Fprintf(sb, "  %s● Failed Pipeline Checks Tree:%s\n", constants.ColorCyan, constants.ColorReset)
		fmt.Fprintf(sb, "    Commit %s%s%s:\n", constants.ColorWhite, sha, constants.ColorReset)

		return
	}

	fmt.Fprintf(sb, "FAILED PIPELINE CHECKS TREE:\n")
	fmt.Fprintf(sb, "  Commit %s:\n", sha)
}

// BuildCommitGroupFailureTree generates a tree view of failed workflows and jobs for a commit group.
func BuildCommitGroupFailureTree(repo string, group *CommitPipelineGroup, isColor bool) string {
	if group == nil {
		return ""
	}

	failingWorkflows := filterFailingWorkflows(group.Workflows)
	if len(failingWorkflows) == 0 {
		return ""
	}

	var sb strings.Builder
	writeTreeHeaderSha(&sb, group.HeadSha, isColor)
	for i, wf := range failingWorkflows {
		isLast := i == len(failingWorkflows)-1
		rawLogs := queryFailedRunLogs(repo, wf.DatabaseId)
		jobs := CorrelateRunFailedJobs(repo, wf.DatabaseId, rawLogs)
		renderTreeWorkflowDetails(&sb, wf.Name, wf.DatabaseId, jobs, isLast, isColor)
	}

	return sb.String()
}

func filterFailingWorkflows(workflows []CommitWorkflowItem) []CommitWorkflowItem {
	var failing []CommitWorkflowItem
	for _, wf := range workflows {
		if isWorkflowFailure(wf) {
			failing = append(failing, wf)
		}
	}

	return failing
}

func renderTreeWorkflowNode(sb *strings.Builder, fr FailedRunItem, isLast bool, isColor bool) {
	renderTreeWorkflowDetails(sb, fr.WorkflowName, fr.RunId, fr.FailedJobs, isLast, isColor)
}

func renderTreeWorkflowDetails(sb *strings.Builder, name string, runId uint64, jobs []FailedJobItem, isLast, isColor bool) {
	branch := "├── "
	indent := "│   "
	if isLast {
		branch = "└── "
		indent = "    "
	}

	sym := formatTreeSymbol(isColor)
	fmt.Fprintf(sb, "    %s%s%s (#%d)\n", branch, sym, name, runId)
	renderTreeJobs(sb, jobs, indent, isColor)
}

func renderTreeJobs(sb *strings.Builder, jobs []FailedJobItem, indent string, isColor bool) {
	for j, job := range jobs {
		isLastJob := j == len(jobs)-1
		renderTreeJobNode(sb, job, indent, isLastJob, isColor)
	}
}

func renderTreeJobNode(sb *strings.Builder, job FailedJobItem, indent string, isLast bool, isColor bool) {
	branch := "├── "
	if isLast {
		branch = "└── "
	}

	sym := formatTreeSymbol(isColor)
	stepPart := formatTreeStepPart(job.StepName)
	fmt.Fprintf(sb, "    %s%s%s%s%s\n", indent, branch, sym, job.JobName, stepPart)
}

func formatTreeStepPart(stepName string) string {
	clean := strings.TrimSpace(stepName)
	if len(clean) > 0 && clean != "Job Execution" && clean != "Failed Step" {
		return fmt.Sprintf(" [Step: %s]", clean)
	}

	return ""
}

func formatTreeSymbol(isColor bool) string {
	if isColor {
		return constants.ColorRed + "✖ " + constants.ColorReset
	}

	return "✖ "
}

func renderFailureTreeTerminal(p PipelineErrorLogsPayload) {
	tree := BuildFailedChecksTree(p, true)
	if len(tree) > 0 {
		fmt.Println(tree)
	}
}

func appendClipboardFailureTree(sb *strings.Builder, p PipelineErrorLogsPayload) {
	tree := BuildFailedChecksTree(p, false)
	if len(tree) > 0 {
		sb.WriteString(tree)
		sb.WriteString("\n")
	}
}
