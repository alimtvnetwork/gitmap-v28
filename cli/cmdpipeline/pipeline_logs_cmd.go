package cmdpipeline

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/atotto/clipboard"
)

// PipelineLogsOptions encapsulates flags and arguments for gitmap pipeline logs.
type PipelineLogsOptions struct {
	Target       string
	Workflow     string
	FilePath     string
	TempFileName string
	FailedOnly   bool
	IsDetailed   bool
	IsJSON       bool
	DoCopy       bool
}

func handlePipelineLogs(args []string) error {
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") {
		printPipelineLogsHelp()

		return nil
	}

	opts := parsePipelineLogsOptions(args)
	repo := resolveCurrentRepoSlug()
	runs := queryWorkflowRuns(repo)
	if len(runs) == 0 {
		return reportNoRunsFound(repo)
	}

	return executePipelineLogsForTarget(repo, runs, opts)
}

func reportNoRunsFound(repo string) error {
	fmt.Printf("No pipeline runs found for %s.\n", repo)

	return nil
}

func parsePipelineLogsOptions(args []string) PipelineLogsOptions {
	return PipelineLogsOptions{
		Target:       extractLogsTargetArg(args),
		Workflow:     extractLogsWorkflowFlag(args),
		FilePath:     extractFlagVal(args, "--file"),
		TempFileName: extractFlagVal(args, "--tempfile"),
		FailedOnly:   hasArgFlag(args, "--failed-only") || hasArgFlag(args, "--failed"),
		IsDetailed:   hasDetailedArg(args),
		IsJSON:       hasArgFlag(args, "--json"),
		DoCopy:       hasArgFlag(args, "--copy") || hasArgFlag(args, "-c") || hasArgFlag(args, "--clip"),
	}
}

func extractLogsWorkflowFlag(args []string) string {
	val := extractFlagVal(args, "--workflow")
	if len(val) > 0 {
		return val
	}

	return extractFlagVal(args, "-w")
}

func extractLogsTargetArg(args []string) string {
	for _, a := range args {
		if isLogsTargetCandidate(a) {
			return a
		}
	}

	return "latest"
}

func isLogsTargetCandidate(a string) bool {
	if tryIsNumericNegativeOffset(a) {
		return true
	}
	if strings.HasPrefix(a, "-") {
		return false
	}

	return true
}

func tryIsNumericNegativeOffset(arg string) bool {
	hasDash := strings.HasPrefix(arg, "-")
	hasChars := len(arg) > 1
	if hasDash && hasChars {
		return isAllDigits(arg[1:])
	}

	return false
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func executePipelineLogsForTarget(repo string, runs []ghRunItem, opts PipelineLogsOptions) error {
	groups := GroupRunsByCommit(runs)
	targetGroup, hasGroup := ResolveCommitGroupByTarget(groups, opts.Target)
	if !hasGroup {
		return reportTargetCommitNotFound(repo, opts.Target)
	}

	combinedLogs, logItems := collectCommitWorkflowLogs(repo, targetGroup, opts)
	recordLogsInSplitDb(repo, targetGroup, combinedLogs)

	return dispatchLogsOutput(combinedLogs, logItems, opts)
}

func reportTargetCommitNotFound(repo, target string) error {
	msg := fmt.Sprintf("no pipeline commit found matching '%s' for %s", target, repo)

	return apperror.NewNotFoundError(msg)
}

func collectCommitWorkflowLogs(repo string, group *CommitPipelineGroup, opts PipelineLogsOptions) (string, []CommitWorkflowLogItem) {
	var sb strings.Builder
	var items []CommitWorkflowLogItem
	for _, wf := range group.Workflows {
		if isWorkflowLogSkipped(wf, opts) {
			continue
		}
		raw := queryAllRunLogs(repo, wf.DatabaseId)
		processed := processWorkflowLogText(raw, opts.IsDetailed)
		appendWorkflowLogSection(&sb, wf, processed)
		items = append(items, buildWorkflowLogItem(wf, processed))
	}

	return CollapseConsecutiveEmptyLines(sb.String()), items
}

func isWorkflowLogSkipped(wf CommitWorkflowItem, opts PipelineLogsOptions) bool {
	if opts.FailedOnly && wf.Conclusion != "failure" {
		return true
	}
	if len(opts.Workflow) > 0 && !strings.EqualFold(wf.Name, opts.Workflow) {
		return true
	}

	return false
}

func processWorkflowLogText(raw string, isDetailed bool) string {
	if isDetailed {
		return raw
	}

	return extractCleanErrorLines(raw)
}

func appendWorkflowLogSection(sb *strings.Builder, wf CommitWorkflowItem, logs string) {
	header := fmt.Sprintf("=== Workflow: %s (#%d) | Status: %s | Conclusion: %s ===\n",
		wf.Name, wf.DatabaseId, wf.Status, wf.Conclusion)
	sb.WriteString(header)
	if len(strings.TrimSpace(logs)) > 0 {
		sb.WriteString(logs)
		sb.WriteString("\n\n")
	} else {
		sb.WriteString("(No log entries or errors recorded)\n\n")
	}
}

// CommitWorkflowLogItem captures structured logs for a workflow in a commit.
type CommitWorkflowLogItem struct {
	DatabaseId uint64 `json:"databaseId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	Logs       string `json:"logs"`
}

func buildWorkflowLogItem(wf CommitWorkflowItem, logs string) CommitWorkflowLogItem {
	return CommitWorkflowLogItem{
		DatabaseId: wf.DatabaseId,
		Name:       wf.Name,
		Status:     wf.Status,
		Conclusion: wf.Conclusion,
		Logs:       logs,
	}
}

func recordLogsInSplitDb(repo string, group *CommitPipelineGroup, logs string) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return
	}
	defer db.Close()
	for _, wf := range group.Workflows {
		_ = db.RecordDetailErrorLog(pipelinedb.PipelineErrorRecord{
			RunId:        wf.DatabaseId,
			RepoSlug:     repo,
			WorkflowName: wf.Name,
			StepName:     "pipeline-logs",
			ErrorText:    logs,
			RawLogs:      logs,
		})
	}
}

func dispatchLogsOutput(combined string, items []CommitWorkflowLogItem, opts PipelineLogsOptions) error {
	if opts.IsJSON {
		return printJSON(items)
	}
	if len(opts.FilePath) > 0 {
		return writeLogsToFile(opts.FilePath, combined)
	}
	if len(opts.TempFileName) > 0 {
		return writeLogsToTempFile(opts.TempFileName, combined)
	}
	if opts.DoCopy {
		_ = clipboard.WriteAll(combined)
		fmt.Println("✓ Copied pipeline logs to system clipboard.")
	}

	fmt.Println(combined)

	return nil
}

func writeLogsToFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "write logs to file")
	}
	fmt.Printf("✓ Saved pipeline logs to %s\n", path)

	return nil
}

func writeLogsToTempFile(name, content string) error {
	tempPath := buildTempLogPath(name)
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "write logs to tempfile")
	}
	fmt.Printf("✓ Saved pipeline logs to %s\n", tempPath)

	return nil
}

func buildTempLogPath(name string) string {
	_ = os.MkdirAll(".lovable/temp", 0755)

	return fmt.Sprintf(".lovable/temp/%s", name)
}
