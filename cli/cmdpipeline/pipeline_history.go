package cmdpipeline

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	return ParseNegativeIndex(arg)
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
	if isFailingConclusion(conclusion) {
		return constants.ColorRed + "FAIL" + constants.ColorReset
	}
	switch conclusion {
	case "success":
		return constants.ColorGreen + "PASS" + constants.ColorReset
	case "skipped":
		return constants.ColorYellow + "SKIPPED" + constants.ColorReset
	}
	if isWorkflowStatusActive(status) {
		return constants.ColorYellow + "RUNNING" + constants.ColorReset
	}
	if len(conclusion) > 0 {
		return conclusion
	}
	if len(status) > 0 {
		return strings.ToUpper(status)
	}

	return "-"
}

func isWorkflowStatusActive(status string) bool {
	return status == "in_progress" || status == "queued" || status == "waiting" || status == "pending" || status == "requested"
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

	return renderPositionalCommitTerminal(repo, runs, groups, group, offset)
}

func printPositionalJSON(group *CommitPipelineGroup) error {
	err := printJSON(group)
	if err != nil {
		return apperror.WrapSimple(err, "inspect_positional_run")
	}

	return nil
}

func renderPositionalCommitTerminal(repo string, runs []ghRunItem, groups []CommitPipelineGroup, group *CommitPipelineGroup, offset int) error {
	var sb strings.Builder
	renderCommitInspectorContent(&sb, repo, runs, group, offset)
	renderRecentCommitsSummaryToBuilder(&sb, groups, 5)
	output := CollapseConsecutiveEmptyLines(sb.String())
	fmt.Print(output)

	return copyPositionalReportToClipboard(output, isCommitGroupFailure(group))
}

func copyPositionalReportToClipboard(terminalOutput string, isFailure bool) error {
	cleanOutput := stripANSI(terminalOutput)
	if err := writeClipboard(cleanOutput); err != nil {
		fmt.Printf("\n  %s⚠ Could not copy to clipboard: %v%s\n", constants.ColorYellow, err, constants.ColorReset)

		return err
	}

	label := "positional pipeline report"
	if isFailure {
		label = "positional pipeline error report"
	}
	fmt.Printf("\n  %s📋 Copied %s to clipboard%s\n", constants.ColorCyan, label, constants.ColorReset)

	return nil
}

func isCommitGroupFailure(group *CommitPipelineGroup) bool {
	return isFailingConclusion(group.Conclusion) || group.FailedWorkflows > 0
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
	eta := resolveCommitWorkflowETA(runs, wf)
	etaStr := formatWorkflowETA(eta, wf.Status)
	fmt.Fprintf(sb, "      - %s (#%d): %s %s\n", wf.Name, wf.DatabaseId, badge, etaStr)
}

func resolveCommitWorkflowETA(runs []ghRunItem, wf CommitWorkflowItem) int {
	if wf.Status != "in_progress" && wf.Status != "queued" {
		return 0
	}
	run := ghRunItem{
		DatabaseId: wf.DatabaseId,
		Name:       wf.Name,
		Status:     wf.Status,
		CreatedAt:  wf.CreatedAt,
		UpdatedAt:  wf.UpdatedAt,
	}

	return CalculateRunETA(run, runs)
}

func formatWorkflowETA(eta int, status string) string {
	if (status == "in_progress" || status == "queued") && eta != 0 {
		return fmt.Sprintf("(ETA: %s)", formatEtaDisplay(eta))
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
	appendCommitGroupFailureTree(sb, repo, group)
	renderFailingWorkflowsDiagnostics(sb, repo, group.Workflows)
}

func appendCommitGroupFailureTree(sb *strings.Builder, repo string, group *CommitPipelineGroup) {
	tree := BuildCommitGroupFailureTree(repo, group, true)
	if len(tree) > 0 {
		sb.WriteString("\n")
		sb.WriteString(tree)
		sb.WriteString("\n")
	}
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
	jobs := CorrelateRunFailedJobs(repo, wf.DatabaseId, rawLogs)
	renderFailedJobItemsToBuilder(sb, jobs)
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
	printFailedJobErrorLinesToBuilder(sb, j.ErrorLines, j.FailureSummary)
	printFailedJobStackTraceToBuilder(sb, j.StackTrace)
}

func printFailedJobErrorLinesToBuilder(sb *strings.Builder, lines []string, summary string) {
	dedupLines := filterOutSummaryLine(lines, summary)
	for _, line := range dedupLines {
		fmt.Fprintf(sb, "        %s\n", line)
	}
}

func printFailedJobStackTraceToBuilder(sb *strings.Builder, stack string) {
	if len(stack) == 0 {
		return
	}

	fmt.Fprintf(sb, "        Stack Trace:\n")
	for _, line := range strings.Split(strings.TrimSpace(stack), "\n") {
		fmt.Fprintf(sb, "          %s\n", line)
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
	fmt.Fprintf(sb, "\n  %s● Recent Commits Pipeline Summary (Last %d Commits):%s\n",
		constants.ColorCyan, count, constants.ColorReset)
	fmt.Fprintf(sb, "    %-8s %-9s %-20s %-11s %-10s %-32s %-8s\n",
		"Offset", "Commit", "Branch", "Release", "Status", "Workflows", "Failures")
	fmt.Fprintf(sb, "    %-8s %-9s %-20s %-11s %-10s %-32s %-8s\n",
		"------", "------", "------", "-------", "------", "---------", "--------")
}

func printRecentCommitRow(sb *strings.Builder, g CommitPipelineGroup, index int) {
	offsetStr := formatCommitOffsetLabel(index)
	sha := truncateHistoryStr(g.HeadSha, 7)
	branch := truncateHistoryStr(g.HeadBranch, 19)
	releaseVer := resolveCommitRelease(g)
	badge := formatStatusBadge(g.Conclusion, g.Status)
	paddedBadge := padRightVisible(badge, 10)
	wfSummary := formatGroupWorkflowsSummary(g.Workflows, 32)
	failuresStr := strconv.Itoa(g.FailedWorkflows)
	fmt.Fprintf(sb, "    %-8s %-9s %-20s %-11s %s %-32s %-8s\n",
		offsetStr, sha, branch, releaseVer, paddedBadge, wfSummary, failuresStr)
}

func resolveCommitRelease(g CommitPipelineGroup) string {
	if len(g.Release) > 0 {
		return g.Release
	}

	relFromBranch := extractReleaseFromBranch(g.HeadBranch)
	if len(relFromBranch) > 0 {
		return relFromBranch
	}

	relFromSha := extractReleaseFromSha(g.HeadSha)
	if len(relFromSha) > 0 {
		return relFromSha
	}

	return "-"
}

func extractReleaseFromBranch(branch string) string {
	lower := strings.ToLower(branch)
	isReleasePrefix := strings.HasPrefix(lower, "release/")
	trimmed := strings.TrimPrefix(branch, "release/")
	hasRelease := isReleasePrefix && len(trimmed) > 0
	if hasRelease {
		return trimmed
	}

	isVPrefix := strings.HasPrefix(lower, "v")
	if isVPrefix && isSemverLike(branch) {
		return branch
	}

	return ""
}

func isSemverLike(s string) bool {
	raw := strings.TrimPrefix(s, "v")
	parts := strings.Split(raw, ".")
	hasEnoughParts := len(parts) >= 2
	if !hasEnoughParts {
		return false
	}

	_, err := strconv.Atoi(parts[0])

	return err == nil
}
