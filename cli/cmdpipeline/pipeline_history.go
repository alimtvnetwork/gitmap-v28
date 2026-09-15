package cmdpipeline

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
	val, err := strconv.Atoi(arg)
	if err == nil && val < 0 {
		return val, true
	}

	return 0, false
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

func formatStatusBadge(conclusion, status string) string {
	switch conclusion {
	case "success":
		return constants.ColorGreen + "PASS" + constants.ColorReset
	case "failure":
		return constants.ColorRed + "FAIL" + constants.ColorReset
	}
	if status == "in_progress" || status == "queued" {
		return constants.ColorYellow + "RUNNING" + constants.ColorReset
	}

	return conclusion
}

func truncateHistoryStr(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}

	return str[:maxLen]
}

// stripANSI removes ANSI escape codes from string s.
func stripANSI(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			i = j
			continue
		}
		out.WriteByte(s[i])
	}

	return out.String()
}

// visibleLen returns the visible terminal width of s.
func visibleLen(s string) int {
	return runewidth.StringWidth(stripANSI(s))
}

// padRightVisible appends spaces to s until its visible terminal width reaches targetWidth.
func padRightVisible(s string, targetWidth int) string {
	vLen := visibleLen(s)
	if vLen >= targetWidth {
		return s
	}

	return s + strings.Repeat(" ", targetWidth-vLen)
}

// RenderHistorySummaryTable outputs a formatted summary table of recent commit pipeline runs.
func RenderHistorySummaryTable(runs []ghRunItem) {
	if len(runs) == 0 {
		return
	}

	groups := GroupRunsByCommit(runs)
	renderRecentCommitsSummaryTable(groups, 5)
}

func resolveInspectGroup(groups []CommitPipelineGroup, offset int) (*CommitPipelineGroup, bool) {
	group, isFound := ResolveCommitGroupByOffset(groups, offset)
	if isFound {
		return group, true
	}
	if offset == -1 && len(groups) == 1 {
		return &groups[0], true
	}

	return nil, false
}

// InspectPositionalRun renders positional commit details for passing, in-progress, or failing runs.
func InspectPositionalRun(repo string, runs []ghRunItem, offset int, isJSON bool) error {
	groups := GroupRunsByCommit(runs)
	group, isFound := resolveInspectGroup(groups, offset)
	if isFound {
		return dispatchInspectGroupOutput(repo, runs, groups, group, offset, isJSON)
	}

	fmt.Printf("No pipeline commit found at offset %d for %s.\n", offset, repo)

	return nil
}

func dispatchInspectGroupOutput(repo string, runs []ghRunItem, groups []CommitPipelineGroup, group *CommitPipelineGroup, offset int, isJSON bool) error {
	if isJSON {
		return printPositionalJSON(group)
	}
	renderPositionalCommitTerminal(repo, runs, groups, group, offset)

	return nil
}

func printPositionalJSON(group *CommitPipelineGroup) error {
	err := printJSON(group)
	if err != nil {
		return apperror.WrapSimple(err, "inspect_positional_run")
	}

	return nil
}

func renderPositionalCommitTerminal(repo string, runs []ghRunItem, groups []CommitPipelineGroup, group *CommitPipelineGroup, offset int) {
	var sb strings.Builder
	renderCommitInspectorContent(&sb, repo, runs, group, offset)
	renderRecentCommitsSummaryToBuilder(&sb, groups, 5)
	output := CollapseConsecutiveEmptyLines(sb.String())
	fmt.Print(output)
}

func isCommitGroupFailure(group *CommitPipelineGroup) bool {
	return group.Conclusion == "failure"
}

func isCommitGroupInProgress(group *CommitPipelineGroup) bool {
	return group.Conclusion == "in_progress" || group.Status == "in_progress"
}

func renderCommitInspectorContent(sb *strings.Builder, repo string, runs []ghRunItem, group *CommitPipelineGroup, offset int) {
	if isCommitGroupFailure(group) {
		renderFailingPositionalCommit(sb, repo, group, offset)
		return
	}
	if isCommitGroupInProgress(group) {
		renderInProgressPositionalCommit(sb, runs, group, offset)
		return
	}
	renderPassingPositionalCommit(sb, group, offset)
}

func renderPassingPositionalCommit(sb *strings.Builder, group *CommitPipelineGroup, offset int) {
	shortSha := truncateHistoryStr(group.HeadSha, 7)
	dur := formatDurationSeconds(group.TotalDuration)
	fmt.Fprintf(sb, "\n  %s● Positional Commit Inspector [Offset %d, Commit %s]: PASSING (clean)%s\n",
		constants.ColorGreen, offset, shortSha, constants.ColorReset)
	fmt.Fprintf(sb, "    • Branch:          %s\n", group.HeadBranch)
	fmt.Fprintf(sb, "    • Status:          %s (conclusion: %s)\n", group.Status, group.Conclusion)
	fmt.Fprintf(sb, "    • Total Duration:  %s\n", dur)
	renderPassingWorkflowsList(sb, group.Workflows)
	fmt.Fprintf(sb, "    %s✓ All pipeline workflows for this commit succeeded with zero errors.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func renderPassingWorkflowsList(sb *strings.Builder, workflows []CommitWorkflowItem) {
	fmt.Fprintf(sb, "    • Workflows (%d):\n", len(workflows))
	for _, wf := range workflows {
		dur := formatDurationSeconds(wf.Duration)
		fmt.Fprintf(sb, "      - %s (#%d): %s (%s)\n", wf.Name, wf.DatabaseId,
			formatStatusBadge(wf.Conclusion, wf.Status), dur)
	}
}

func renderInProgressPositionalCommit(sb *strings.Builder, runs []ghRunItem, group *CommitPipelineGroup, offset int) {
	shortSha := truncateHistoryStr(group.HeadSha, 7)
	dur := formatDurationSeconds(group.TotalDuration)
	fmt.Fprintf(sb, "\n  %s● Positional Commit Inspector [Offset %d, Commit %s]: IN_PROGRESS%s\n",
		constants.ColorYellow, offset, shortSha, constants.ColorReset)
	fmt.Fprintf(sb, "    • Branch:          %s\n", group.HeadBranch)
	fmt.Fprintf(sb, "    • Status:          in_progress\n")
	fmt.Fprintf(sb, "    • Total Duration:  %s\n", dur)
	fmt.Fprintf(sb, "    • Active Workflows:\n")
	renderActiveWorkflowsList(sb, runs, group.Workflows)
	fmt.Fprintln(sb)
}

func renderActiveWorkflowsList(sb *strings.Builder, runs []ghRunItem, workflows []CommitWorkflowItem) {
	for _, wf := range workflows {
		renderSingleWorkflowStatus(sb, runs, wf)
	}
}

func renderSingleWorkflowStatus(sb *strings.Builder, runs []ghRunItem, wf CommitWorkflowItem) {
	badge := formatStatusBadge(wf.Conclusion, wf.Status)
	eta := calculateAverageDuration(runs, wf.Name)
	etaStr := formatWorkflowETA(eta, wf.Status)
	fmt.Fprintf(sb, "      - %s (#%d): %s %s\n", wf.Name, wf.DatabaseId, badge, etaStr)
}

func formatWorkflowETA(eta int, status string) string {
	if (status == "in_progress" || status == "queued") && eta > 0 {
		return fmt.Sprintf("(ETA: ~%ds)", eta)
	}

	return ""
}

func renderFailingPositionalCommit(sb *strings.Builder, repo string, group *CommitPipelineGroup, offset int) {
	shortSha := truncateHistoryStr(group.HeadSha, 7)
	dur := formatDurationSeconds(group.TotalDuration)
	fmt.Fprintf(sb, "\n  %s● Positional Commit Inspector [Offset %d, Commit %s]: FAILING%s\n",
		constants.ColorRed, offset, shortSha, constants.ColorReset)
	fmt.Fprintf(sb, "    • Branch:          %s\n", group.HeadBranch)
	fmt.Fprintf(sb, "    • Status:          %s (conclusion: %s)\n", group.Status, group.Conclusion)
	fmt.Fprintf(sb, "    • Total Duration:  %s\n", dur)
	renderFailingWorkflowsDiagnostics(sb, repo, group.Workflows)
	renderPassingWorkflowsSummary(sb, group.Workflows)
}

func renderFailingWorkflowsDiagnostics(sb *strings.Builder, repo string, workflows []CommitWorkflowItem) {
	for _, wf := range workflows {
		if isWorkflowFailure(wf) {
			renderSingleFailingWorkflowLogs(sb, repo, wf)
		}
	}
}

func renderSingleFailingWorkflowLogs(sb *strings.Builder, repo string, wf CommitWorkflowItem) {
	fmt.Fprintf(sb, "    • Failed Workflow: %s (#%d) - %s\n", wf.Name, wf.DatabaseId, wf.Url)
	rawLogs := queryFailedRunLogs(repo, wf.DatabaseId)
	renderFailedJobItemsToBuilder(sb, ParseFailedLogLines(rawLogs))
}

func renderPassingWorkflowsSummary(sb *strings.Builder, workflows []CommitWorkflowItem) {
	var passingNames []string
	for _, wf := range workflows {
		if wf.Conclusion == "success" {
			passingNames = append(passingNames, fmt.Sprintf("%s (#%d)", wf.Name, wf.DatabaseId))
		}
	}
	if len(passingNames) > 0 {
		fmt.Fprintf(sb, "    • Passing Workflows: %s\n\n", strings.Join(passingNames, ", "))
	}
}

func renderFailedJobItemsToBuilder(sb *strings.Builder, jobs []FailedJobItem) {
	if len(jobs) == 0 {
		fmt.Fprintln(sb, "    • No detailed failing steps could be parsed.")

		return
	}

	fmt.Fprintln(sb, "    • Failing Steps & Diagnostics:")
	for _, j := range jobs {
		printFailedJobItemToBuilder(sb, j)
	}
	fmt.Fprintln(sb)
}

func printFailedJobItemToBuilder(sb *strings.Builder, j FailedJobItem) {
	fmt.Fprintf(sb, "      %s[Job: %s | Step: %s]%s\n",
		constants.ColorCyan, j.JobName, j.StepName, constants.ColorReset)
	if len(j.FailureSummary) > 0 {
		fmt.Fprintf(sb, "        Error: %s%s%s\n",
			constants.ColorRed, j.FailureSummary, constants.ColorReset)
	}
}

func renderRecentCommitsSummaryTable(groups []CommitPipelineGroup, limit int) {
	if len(groups) == 0 {
		return
	}

	var sb strings.Builder
	renderRecentCommitsSummaryToBuilder(&sb, groups, limit)
	fmt.Print(CollapseConsecutiveEmptyLines(sb.String()))
}

func renderRecentCommitsSummaryToBuilder(sb *strings.Builder, groups []CommitPipelineGroup, limit int) {
	if len(groups) == 0 {
		return
	}

	displayGroups := capCommitGroups(groups, limit)
	printRecentCommitsHeader(sb, len(displayGroups))
	for i, g := range displayGroups {
		printRecentCommitRow(sb, g, i)
	}
	fmt.Fprintln(sb)
}

func capCommitGroups(groups []CommitPipelineGroup, limit int) []CommitPipelineGroup {
	if len(groups) > limit {
		return groups[:limit]
	}

	return groups
}

func printRecentCommitsHeader(sb *strings.Builder, count int) {
	fmt.Fprintf(sb, "  %s● Recent Commits Pipeline Summary (Last %d Commits):%s\n",
		constants.ColorCyan, count, constants.ColorReset)
	fmt.Fprintf(sb, "    %-8s %-9s %-14s %-10s %-32s %-8s\n",
		"Offset", "Commit", "Branch", "Status", "Workflows", "Failures")
	fmt.Fprintf(sb, "    %-8s %-9s %-14s %-10s %-32s %-8s\n",
		"------", "------", "------", "------", "---------", "--------")
}

func printRecentCommitRow(sb *strings.Builder, g CommitPipelineGroup, index int) {
	offsetStr := formatCommitOffsetLabel(index)
	sha := truncateHistoryStr(g.HeadSha, 7)
	branch := truncateHistoryStr(g.HeadBranch, 13)
	badge := formatStatusBadge(g.Conclusion, g.Status)
	paddedBadge := padRightVisible(badge, 10)
	wfSummary := formatGroupWorkflowsSummary(g.Workflows, 32)
	failuresStr := strconv.Itoa(g.FailedWorkflows)
	fmt.Fprintf(sb, "    %-8s %-9s %-14s %s %-32s %-8s\n",
		offsetStr, sha, branch, paddedBadge, wfSummary, failuresStr)
}

func formatGroupWorkflowsSummary(workflows []CommitWorkflowItem, maxWidth int) string {
	if len(workflows) == 0 {
		return "-"
	}
	full := summarizeGroupWorkflows(workflows)
	if visibleLen(full) <= maxWidth {
		return full
	}

	return formatWorkflowsWithRemaining(workflows, maxWidth)
}

func formatWorkflowsWithRemaining(workflows []CommitWorkflowItem, maxWidth int) string {
	var parts []string
	for i, wf := range workflows {
		item := fmt.Sprintf("%s [%s]", wf.Name, formatWorkflowShortStatus(wf))
		remaining := len(workflows) - (i + 1)
		candidate := buildWorkflowCandidate(parts, item, remaining)
		if visibleLen(candidate) > maxWidth {
			return finalizeTruncatedSummary(workflows, parts, remaining, maxWidth)
		}
		parts = append(parts, item)
	}

	return strings.Join(parts, ", ")
}

func buildWorkflowCandidate(parts []string, nextItem string, remaining int) string {
	suffix := ""
	if remaining > 0 {
		suffix = fmt.Sprintf(" (+%d)", remaining)
	}
	if len(parts) == 0 {
		return nextItem + suffix
	}

	return strings.Join(parts, ", ") + ", " + nextItem + suffix
}

func finalizeTruncatedSummary(workflows []CommitWorkflowItem, parts []string, remaining int, maxWidth int) string {
	if len(parts) > 0 {
		return strings.Join(parts, ", ") + fmt.Sprintf(" (+%d)", remaining)
	}
	if len(workflows) == 0 {
		return "-"
	}
	first := fmt.Sprintf("%s [%s]", workflows[0].Name, formatWorkflowShortStatus(workflows[0]))
	if len(workflows) > 1 {
		first += fmt.Sprintf(" (+%d)", len(workflows)-1)
	}
	if visibleLen(first) > maxWidth && maxWidth > 3 {
		return first[:maxWidth-3] + "..."
	}

	return first
}

func formatCommitOffsetLabel(index int) string {
	if index == 0 {
		return "latest"
	}

	return fmt.Sprintf("-%d", index)
}

func summarizeGroupWorkflows(workflows []CommitWorkflowItem) string {
	var parts []string
	for _, wf := range workflows {
		statusShort := formatWorkflowShortStatus(wf)
		parts = append(parts, fmt.Sprintf("%s [%s]", wf.Name, statusShort))
	}

	return strings.Join(parts, ", ")
}

func formatWorkflowShortStatus(wf CommitWorkflowItem) string {
	switch wf.Conclusion {
	case "success":
		return "PASS"
	case "failure":
		return "FAIL"
	}
	if wf.Status == "in_progress" || wf.Status == "queued" {
		return "RUN"
	}

	return wf.Conclusion
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
		printEmptyCachedFailures(repo, db.Path)

		return
	}

	printCachedFailuresList(db, relDb, runs, isDetailed)
}

func printEmptyCachedFailures(repo, dbPath string) {
	fmt.Printf("\n  No cached pipeline failures found in SQLite for %s.\n", repo)
	fmt.Printf("  Pipeline DB: %s\n", FormatRelativeDbPath(dbPath))
	fmt.Printf("  DB Size:     %s\n\n", ResolveDbFileSize(dbPath))
}

func printCachedFailuresList(db *pipelinedb.PipelineSplitDb, relDb string, runs []pipelinedb.PipelineRunRecord, isDetailed bool) {
	fmt.Printf("\n  %s● Last %d Cached Pipeline Failure(s) from SQLite (%s):%s\n",
		constants.ColorRed, len(runs), relDb, constants.ColorReset)
	for i, r := range runs {
		renderSingleCachedFailureCard(db, r, i+1, len(runs), isDetailed)
	}
}

func fetchLastCachedFailures(repo string, count int) (*pipelinedb.PipelineSplitDb, []pipelinedb.PipelineRunRecord, error) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "fetch_last_cached_failures")
	}

	runs, qErr := queryCachedFailedRuns(db, count)
	if qErr != nil {
		_ = db.Close()

		return nil, nil, qErr
	}

	return db, runs, nil
}

func queryCachedFailedRuns(db *pipelinedb.PipelineSplitDb, count int) ([]pipelinedb.PipelineRunRecord, error) {
	limit := resolveMaxSyncLimit(count)
	runRes := db.QueryLastFailedRuns(limit)
	if runRes.IsFailure() {
		return nil, runRes.AppError()
	}

	return runRes.Data, nil
}

// RenderLastCachedFailures displays cached failure logs from SQLite up to count.
func RenderLastCachedFailures(repo string, count int, isJSON bool, isDetailed ...bool) error {
	detailed := isDetailedRequested(isDetailed)
	db, runs, err := fetchLastCachedFailures(repo, count)
	if err != nil {
		return err
	}
	defer db.Close()

	if isJSON {
		return printCachedRunsJSON(runs)
	}

	renderCachedFailuresTerminal(db, repo, runs, detailed)

	return nil
}

func isDetailedRequested(isDetailed []bool) bool {
	return len(isDetailed) > 0 && isDetailed[0]
}

func printCachedRunsJSON(runs []pipelinedb.PipelineRunRecord) error {
	err := printJSON(runs)
	if err != nil {
		return apperror.WrapSimple(err, "render_last_cached_failures")
	}

	return nil
}

// HandlePipelineHistoryErrors inspects args for negative offset or --last-failures flag and executes them.
func HandlePipelineHistoryErrors(args []string) (bool, error) {
	repo := resolveCurrentRepoSlug()
	isJSON := hasArgFlag(args, "--json")

	if isHandled, err := tryHandleOffsetInspection(repo, args, isJSON); isHandled {
		return true, err
	}

	return tryHandleLastFailuresInspection(repo, args, isJSON)
}

func tryHandleOffsetInspection(repo string, args []string, isJSON bool) (bool, error) {
	offset, hasOffset := extractNegativeOffset(args)
	if hasOffset {
		runs := queryWorkflowRuns(repo)

		return true, InspectPositionalRun(repo, runs, offset, isJSON)
	}

	return false, nil
}

func tryHandleLastFailuresInspection(repo string, args []string, isJSON bool) (bool, error) {
	count, hasLastFailures := extractLastFailuresFlag(args)
	if hasLastFailures {
		isDetailed := hasDetailedArg(args)

		return true, RenderLastCachedFailures(repo, count, isJSON, isDetailed)
	}

	return false, nil
}
