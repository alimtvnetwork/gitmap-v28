package cmdpipeline

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func extractNegativeOffset(args []string) (int, bool) {
	for _, arg := range args {
		val, hasOffset := tryParseNegativeOffset(arg)
		if hasOffset {
			return val, true
		}
	}

	return 0, false
}

func tryParseNegativeOffset(arg string) (int, bool) {
	if !strings.HasPrefix(arg, "-") || len(arg) <= 1 {
		return 0, false
	}

	val, err := strconv.Atoi(arg)
	if err != nil || val >= 0 {
		return 0, false
	}

	return val, true
}

func extractLastFailuresFlag(args []string) (int, bool) {
	if val, hasVal := ParseLastFailuresFlag(args); hasVal {
		return val, true
	}

	if hasArgFlag(args, "--last-failures") {
		return 5, true
	}

	return 0, false
}

func capHistoryRuns(runs []ghRunItem, limit int) []ghRunItem {
	if len(runs) > limit {
		return runs[:limit]
	}

	return runs
}

func printHistoryTableHeader() {
	fmt.Printf("    %-5s %-10s %-16s %-10s %-10s %-12s %-8s\n",
		"Pos", "Run ID", "Workflow", "Status", "Duration", "Branch", "Commit")
	fmt.Printf("    %-5s %-10s %-16s %-10s %-10s %-12s %-8s\n",
		"---", "------", "--------", "------", "--------", "------", "------")
}

func formatStatusBadge(conclusion, status string) string {
	if conclusion == "success" {
		return constants.ColorGreen + "PASS" + constants.ColorReset
	}

	if conclusion == "failure" {
		return constants.ColorRed + "FAIL" + constants.ColorReset
	}

	if status == "in_progress" || status == "queued" {
		return constants.ColorYellow + "RUNNING" + constants.ColorReset
	}

	return conclusion
}

func normalizeNegativeOffset(offset int) int {
	if offset < 0 {
		offset = -offset
	}

	if offset > 0 {
		return offset - 1
	}

	return 0
}

func truncateHistoryStr(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}

	return str[:maxLen]
}

func printHistoryTableRow(r ghRunItem, index int) {
	posStr := fmt.Sprintf("-%d", index)
	runIdStr := fmt.Sprintf("#%d", r.DatabaseId)
	badge := formatStatusBadge(r.Conclusion, r.Status)
	dur := formatDurationSeconds(calculateRunDuration(r.CreatedAt, r.UpdatedAt))
	branch := truncateHistoryStr(r.HeadBranch, 11)
	sha := truncateHistoryStr(r.HeadSha, 7)
	wf := truncateHistoryStr(r.Name, 15)
	fmt.Printf("    %-5s %-10s %-16s %-10s %-10s %-12s %-8s\n",
		posStr, runIdStr, wf, badge, dur, branch, sha)
}

// RenderHistorySummaryTable outputs a formatted summary table of the last 5 pipeline runs.
func RenderHistorySummaryTable(runs []ghRunItem) {
	if len(runs) == 0 {
		return
	}

	limitRuns := capHistoryRuns(runs, 5)
	fmt.Printf("  %s● Recent Pipeline Execution History (Last %d Runs):%s\n",
		constants.ColorCyan, len(limitRuns), constants.ColorReset)
	printHistoryTableHeader()
	for i, r := range limitRuns {
		printHistoryTableRow(r, i+1)
	}

	fmt.Println()
}

func resolveRunByOffset(runs []ghRunItem, offset int) (ghRunItem, bool) {
	absIdx := normalizeNegativeOffset(offset)
	if absIdx < len(runs) {
		return runs[absIdx], true
	}

	return ghRunItem{}, false
}

func renderPassingPositionalRun(runs []ghRunItem, run ghRunItem, offset int) {
	eta := calculateAverageDuration(runs, run.Name)
	dur := calculateRunDuration(run.CreatedAt, run.UpdatedAt)
	fmt.Printf("\n  %s● Positional Run Inspector [Offset %d, Run #%d]: PASSING (clean)%s\n",
		constants.ColorGreen, offset, run.DatabaseId, constants.ColorReset)
	fmt.Printf("    • Workflow:        %s\n", run.Name)
	fmt.Printf("    • Status:          %s (conclusion: %s)\n", run.Status, run.Conclusion)
	fmt.Printf("    • Duration:        %s\n", formatDurationSeconds(dur))
	fmt.Printf("    • Historical ETA:  ~%ds\n", eta)
	fmt.Printf("    • Branch / Commit: %s (%s)\n", run.HeadBranch, run.HeadSha)
	fmt.Printf("    • Web Run URL:     %s\n", run.Url)
	fmt.Printf("    %s✓ This historical pipeline run succeeded with zero errors.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func printFailedJobItem(j FailedJobItem) {
	fmt.Printf("      %s[Job: %s | Step: %s]%s\n",
		constants.ColorCyan, j.JobName, j.StepName, constants.ColorReset)
	if len(j.FailureSummary) > 0 {
		fmt.Printf("        Error: %s%s%s\n",
			constants.ColorRed, j.FailureSummary, constants.ColorReset)
	}
}

func renderFailedJobItems(jobs []FailedJobItem) {
	if len(jobs) == 0 {
		fmt.Println("    • No detailed failing steps could be parsed.")

		return
	}

	fmt.Println("    • Failing Steps & Diagnostics:")
	for _, j := range jobs {
		printFailedJobItem(j)
	}

	fmt.Println()
}

func renderFailingPositionalRun(repo string, run ghRunItem, offset int) {
	dur := calculateRunDuration(run.CreatedAt, run.UpdatedAt)
	fmt.Printf("\n  %s● Positional Run Inspector [Offset %d, Run #%d]: FAILING%s\n",
		constants.ColorRed, offset, run.DatabaseId, constants.ColorReset)
	fmt.Printf("    • Workflow:        %s\n", run.Name)
	fmt.Printf("    • Status:          %s (conclusion: %s)\n", run.Status, run.Conclusion)
	fmt.Printf("    • Duration:        %s\n", formatDurationSeconds(dur))
	fmt.Printf("    • Branch / Commit: %s (%s)\n", run.HeadBranch, run.HeadSha)
	fmt.Printf("    • Web Run URL:     %s\n", run.Url)
	rawLogs := queryFailedRunLogs(repo, run.DatabaseId)
	renderFailedJobItems(ParseFailedLogLines(rawLogs))
}

func renderPositionalRunTerminal(repo string, runs []ghRunItem, run ghRunItem, offset int) {
	isFailure := run.Conclusion == "failure"
	if isFailure {
		renderFailingPositionalRun(repo, run, offset)

		return
	}

	renderPassingPositionalRun(runs, run, offset)
}

// InspectPositionalRun renders positional run details for passing or failing runs.
func InspectPositionalRun(repo string, runs []ghRunItem, offset int, isJSON bool) error {
	run, hasRun := resolveRunByOffset(runs, offset)
	if !hasRun {
		fmt.Printf("No pipeline run found at offset %d for %s.\n", offset, repo)

		return nil
	}

	if isJSON {
		return printJSON(run)
	}

	renderPositionalRunTerminal(repo, runs, run, offset)

	return nil
}

func renderCachedErrorsList(errors []pipelinedb.PipelineErrorRecord) {
	if len(errors) == 0 {
		fmt.Println("  │ Error: (No error diagnostic logs cached for this run)")

		return
	}

	for _, e := range errors {
		fmt.Printf("  │ Step:  %s\n", e.StepName)
		fmt.Printf("  │ Error: %s%s%s\n", constants.ColorRed, e.ErrorText, constants.ColorReset)
	}
}

func renderSingleCachedFailureCard(db *pipelinedb.PipelineSplitDb, r pipelinedb.PipelineRunRecord, idx, total int, isDetailed bool) {
	fmt.Printf("  %s┌─ [%d/%d] %s #%d ──────────────────────────%s\n",
		constants.ColorRed, idx, total, r.WorkflowName, r.RunId, constants.ColorReset)
	fmt.Printf("  │ Branch / Commit: %s (%s)\n", r.Branch, r.Sha)
	fmt.Printf("  │ When Run:        %s | Duration: %s\n",
		formatRunTimestamp(r.CreatedAt), formatDurationSeconds(r.DurationSeconds))
	renderRunErrorsByDetailMode(db, r.RunId, isDetailed)
	fmt.Printf("  %s└──────────────────────────────────────────────────────────%s\n\n",
		constants.ColorRed, constants.ColorReset)
}

func renderRunErrorsByDetailMode(db *pipelinedb.PipelineSplitDb, runId uint64, isDetailed bool) {
	if isDetailed {
		renderDetailedOrLegacyErrors(db, runId)

		return
	}

	renderCompactOrLegacyErrors(db, runId)
}

func renderDetailedOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	detailRes := db.QueryDetailedErrorLogsByRunId(runId)
	if detailRes.HasRecord() {
		renderCachedErrorsList(detailRes.Data)

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Data)
}

func renderCompactOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	compactRes := db.QueryCompactErrorLogsByRunId(runId)
	if compactRes.HasRecord() {
		renderCachedCompactErrorsList(compactRes.Data)

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Data)
}

func renderCachedCompactErrorsList(errors []pipelinedb.PipelineCompactErrorRecord) {
	if len(errors) == 0 {
		fmt.Println("  │ Error: (No error diagnostic logs cached for this run)")

		return
	}

	for _, e := range errors {
		fmt.Printf("  │ Step:  %s\n", e.StepName)
		fmt.Printf("  │ Error: %s%s%s\n", constants.ColorRed, e.ErrorText, constants.ColorReset)
	}
}

func renderCachedFailuresTerminal(db *pipelinedb.PipelineSplitDb, repo string, runs []pipelinedb.PipelineRunRecord, isDetailed bool) {
	relDb := FormatRelativeDbPath(db.Path)
	if len(runs) == 0 {
		fmt.Printf("\n  No cached pipeline failures found in SQLite for %s.\n", repo)
		fmt.Printf("  Pipeline DB: %s\n\n", relDb)

		return
	}

	fmt.Printf("\n  %s● Last %d Cached Pipeline Failure(s) from SQLite (%s):%s\n",
		constants.ColorRed, len(runs), relDb, constants.ColorReset)
	for i, r := range runs {
		renderSingleCachedFailureCard(db, r, i+1, len(runs), isDetailed)
	}
}

func fetchLastCachedFailures(repo string, count int) (*pipelinedb.PipelineSplitDb, []pipelinedb.PipelineRunRecord, error) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, nil, err
	}

	limit := resolveMaxSyncLimit(count)
	runRes := db.QueryLastFailedRuns(limit)
	if runRes.IsFailure() {
		_ = db.Close()

		return nil, nil, runRes.AppError()
	}

	return db, runRes.Data, nil
}

// RenderLastCachedFailures displays cached failure logs from SQLite up to count.
func RenderLastCachedFailures(repo string, count int, isJSON bool, isDetailed ...bool) error {
	detailed := len(isDetailed) > 0 && isDetailed[0]
	db, runs, err := fetchLastCachedFailures(repo, count)
	if err != nil {
		return err
	}

	defer db.Close()
	if isJSON {
		return printJSON(runs)
	}

	renderCachedFailuresTerminal(db, repo, runs, detailed)

	return nil
}

// HandlePipelineHistoryErrors inspects args for negative offset or --last-failures flag and executes them.
func HandlePipelineHistoryErrors(args []string) (bool, error) {
	repo := resolveCurrentRepoSlug()
	isJSON := hasArgFlag(args, "--json")
	offset, hasOffset := extractNegativeOffset(args)
	if hasOffset {
		runs := queryWorkflowRuns(repo)

		return true, InspectPositionalRun(repo, runs, offset, isJSON)
	}

	count, hasLastFailures := extractLastFailuresFlag(args)
	if hasLastFailures {
		isDetailed := hasDetailedArg(args)

		return true, RenderLastCachedFailures(repo, count, isJSON, isDetailed)
	}

	return false, nil
}
