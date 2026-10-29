package heatmap

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderSearchText(out io.Writer, records []store.HeatmapActionLogRecord) {
	if len(records) == 0 {
		fmt.Fprintln(out, "No heatmap action logs found matching query.")
		return
	}

	fmt.Fprintf(out, "Found %d heatmap action log(s):\n", len(records))
	fmt.Fprintf(out, "%-6s %-14s %-8s %-30s %s\n", "ID", "STATE", "ACTION", "FILE", "REASON")
	fmt.Fprintf(out, "%-6s %-14s %-8s %-30s %s\n", "------", "--------------", "--------", "------------------------------", "--------------------")
	for _, r := range records {
		printSearchRecord(out, r)
	}
}

func printSearchRecord(out io.Writer, r store.HeatmapActionLogRecord) {
	fStr := r.FilePath
	if r.StartLine > 0 {
		fStr += fmt.Sprintf(":%d", r.StartLine)
	}
	if len(fStr) > 30 {
		fStr = "..." + fStr[len(fStr)-27:]
	}
	fmt.Fprintf(out, "%-6d %-14s %-8s %-30s %s\n", r.HeatmapActionLogId, r.ReasoningState, r.ActionType, fStr, r.ReasoningNotes)
}

func formatTreeBranch(branch string) string {
	return termout.RoleSGR(termout.RoleFGSubtle) + branch + termout.Reset()
}

func formatNodeLogID(id int64) string {
	return termout.RoleSGRBold(termout.RoleAccent) + fmt.Sprintf("#%d", id) + termout.Reset()
}

func formatReasoningState(state string) string {
	role := termout.RoleFG
	switch state {
	case "PASS", "SUCCESS", "OPTIMAL":
		role = termout.RoleOK
	case "FAIL", "ERROR", "CRITICAL":
		role = termout.RoleError
	case "RETRY", "WARN":
		role = termout.RoleWarn
	}

	return termout.RoleSGR(role) + state + termout.Reset()
}

func resolveBranchIndent(i, count int, indent string) (string, string) {
	if indent == "" {
		return "• ", "  "
	}
	if i == count-1 {
		return "└── ", indent + "    "
	}

	return "├── ", indent + "│   "
}

func formatFileLineInfo(filePath string, startLine, endLine int) string {
	if filePath == "" {
		return ""
	}

	cleanPath := filepath.ToSlash(filePath)
	if startLine == 0 && endLine == 0 {
		return fmt.Sprintf(" [%s]", cleanPath)
	}

	return fmt.Sprintf(" [%s:%d-%d]", cleanPath, startLine, endLine)
}

func renderTreeText(out io.Writer, nodes []store.HeatmapExecutionTreeNode, indent string) {
	if len(nodes) == 0 && indent == "" {
		fmt.Fprintln(out, "No in-flight file touches or action tree recorded.")

		return
	}

	for i, node := range nodes {
		branch, nextIndent := resolveBranchIndent(i, len(nodes), indent)
		lineInfo := formatFileLineInfo(node.FilePath, node.StartLine, node.EndLine)
		styledBranch := formatTreeBranch(branch)
		styledID := formatNodeLogID(node.HeatmapActionLogId)
		styledState := formatReasoningState(node.ReasoningState)

		fmt.Fprintf(out, "%s%s%s (%s) %s%s: %s\n", indent, styledBranch, styledID,
			styledState, node.ActionType, lineInfo, node.ReasoningNotes)

		if len(node.Children) > 0 {
			renderTreeText(out, node.Children, nextIndent)
		}
	}
}

func renderRcTreeText(out io.Writer, nodes []store.RcHeatmapExecutionTreeNode, indent string) {
	if len(nodes) == 0 && indent == "" {
		fmt.Fprintln(out, "No in-flight file touches or action tree recorded.")

		return
	}

	for i, node := range nodes {
		branch, nextIndent := resolveBranchIndent(i, len(nodes), indent)
		lineInfo := formatFileLineInfo(node.FilePath, node.StartLine, node.EndLine)
		styledBranch := formatTreeBranch(branch)
		styledID := formatNodeLogID(node.HeatmapActionLogId)
		styledState := formatReasoningState(node.ReasoningState)

		fmt.Fprintf(out, "%s%s%s (%s) %s%s: %s\n", indent, styledBranch, styledID,
			styledState, node.ActionType, lineInfo, node.ReasoningNotes)

		if len(node.Children) > 0 {
			renderRcTreeText(out, node.Children, nextIndent)
		}
	}
}

func renderLiveTouches(out io.Writer, records []store.HeatmapActionLogRecord, title string) {
	if len(records) == 0 {
		fmt.Fprintln(out, "No in-flight file touches or active live actions recorded.")

		return
	}

	fmt.Fprintf(out, "In-Flight File Touches & Live Telemetry [%s]:\n", title)
	for _, r := range records {
		printLiveRecord(out, r)
	}
}

func printLiveRecord(out io.Writer, r store.HeatmapActionLogRecord) {
	cleanPath := filepath.ToSlash(r.FilePath)
	lineInfo := formatFileLineInfo(cleanPath, r.StartLine, r.EndLine)
	styledID := formatNodeLogID(r.HeatmapActionLogId)
	styledState := formatReasoningState(r.ReasoningState)

	fmt.Fprintf(out, " %s (%s) %s%s - %s\n", styledID, styledState, r.ActionType, lineInfo, r.ReasoningNotes)
}

func renderStatusDashboard(out io.Writer, s store.HeatmapStatusSummary, inv TaskInventory, alerts []HeatmapAiConfusion, cmds SuggestedCommandsPayload) {
	wsSlug := ResolveWorkspaceSlug(".")
	activeSlug := ResolveActiveTaskSlug(s.TaskSlug, inv)

	renderStatusHeaderCard(out, wsSlug, activeSlug)
	renderTelemetrySummarySection(out, s)
	renderTopTouchedFilesSection(out, s.TopTouchedFiles)
	renderCompletedTasksSection(out, inv.Completed)
	renderPendingTasksSection(out, inv.Pending)
	renderAlertsSection(out, alerts)
	renderDrillDownSection(out, cmds)
}

func renderStatusHeaderCard(out io.Writer, wsSlug, activeSlug string) {
	title := termout.TokyoBannerStyle("GITMAP HEATMAP TELEMETRY STATUS")
	activeDisplay := activeSlug
	if activeDisplay == "" {
		activeDisplay = "(none)"
	}

	meta := fmt.Sprintf("  Workspace: %s  •  AI: 01-antigravity  •  Active Task: %s", wsSlug, activeDisplay)
	card := termout.TokyoBoxStyle(title + "\n" + meta)

	fmt.Fprintln(out, card)
	fmt.Fprintln(out)
}

func renderTelemetrySummarySection(out io.Writer, s store.HeatmapStatusSummary) {
	fmt.Fprintln(out, "  "+termout.TokyoHeader("Telemetry Summary"))
	fmt.Fprintln(out, "  "+renderSectionDivider())
	fmt.Fprintf(out, "  Total Actions Recorded : %s\n", formatActionsCount(s.TotalActions))
	fmt.Fprintf(out, "  Expectation Pass Rate  : %.1f%% [%d / %d]\n", s.ExpectationsMetPct, s.ExpectationsMetCount, s.TotalActions)
	fmt.Fprintf(out, "  Active Reasoning State : %s\n\n", formatReasoningStatesLine(s.ReasoningStateCounts))
}

func renderSectionDivider() string {
	return termout.SafeTokyoStyle(termout.TokyoMuted).Render(strings.Repeat("─", 75))
}

func formatActionsCount(total int) string {
	return fmt.Sprintf("%d actions", total)
}

func extractSortedStateKeys(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func formatSingleStateBadge(state string, count int) string {
	styledState := formatReasoningState(state)

	return fmt.Sprintf("%s (%d)", styledState, count)
}

func formatReasoningStatesLine(counts map[string]int) string {
	if len(counts) == 0 {
		return termout.SafeTokyoStyle(termout.TokyoMuted).Render("None recorded")
	}

	keys := extractSortedStateKeys(counts)
	var parts []string
	for _, k := range keys {
		parts = append(parts, formatSingleStateBadge(k, counts[k]))
	}

	return strings.Join(parts, "  •  ")
}

func renderTopTouchedFilesSection(out io.Writer, files []store.FileTouchCount) {
	if len(files) == 0 {
		return
	}

	fmt.Fprintln(out, "  "+termout.TokyoHeader("Top Touched Files"))
	fmt.Fprintln(out, "  "+renderSectionDivider())
	for _, f := range files {
		touchStr := termout.SafeTokyoStyle(termout.TokyoMuted).Render(fmt.Sprintf("(%d touches)", f.Count))
		fmt.Fprintf(out, "    %-45s %s\n", f.FilePath, touchStr)
	}

	fmt.Fprintln(out)
}

func renderCompletedTasksSection(out io.Writer, tasks []TaskItem) {
	if len(tasks) == 0 {
		return
	}

	header := fmt.Sprintf("Completed Tasks (Last %d)", len(tasks))
	fmt.Fprintln(out, "  "+termout.TokyoHeader(header))
	fmt.Fprintln(out, "  "+renderSectionDivider())
	for _, t := range tasks {
		badge := termout.SafeTokyoStyle(termout.TokyoOK).Bold(true).Render("[✓]")
		actionStr := termout.SafeTokyoStyle(termout.TokyoMuted).Render(fmt.Sprintf("(%d actions)", t.TotalActions))
		fmt.Fprintf(out, "  %s %-56s %s\n", badge, t.TaskSlug, actionStr)
	}

	fmt.Fprintln(out)
}

func formatPendingTaskLine(t TaskItem) string {
	if t.IsActive {
		badge := termout.SafeTokyoStyle(termout.TokyoWarn).Bold(true).Render("[⏳]")
		tag := termout.SafeTokyoStyle(termout.TokyoWarn).Render("(Active)")

		return fmt.Sprintf("  %s %-56s %s", badge, t.TaskSlug, tag)
	}

	badge := termout.SafeTokyoStyle(termout.TokyoMuted).Render("[ ]")
	tag := termout.SafeTokyoStyle(termout.TokyoMuted).Render("(Queued)")

	return fmt.Sprintf("  %s %-56s %s", badge, t.TaskSlug, tag)
}

func renderPendingTasksSection(out io.Writer, tasks []TaskItem) {
	if len(tasks) == 0 {
		return
	}

	header := fmt.Sprintf("Pending Tasks Queue (Top %d)", len(tasks))
	fmt.Fprintln(out, "  "+termout.TokyoHeader(header))
	fmt.Fprintln(out, "  "+renderSectionDivider())
	for _, t := range tasks {
		fmt.Fprintln(out, formatPendingTaskLine(t))
	}

	fmt.Fprintln(out)
}

func formatAlertLine(a HeatmapAiConfusion) string {
	badge := termout.SafeTokyoStyle(termout.TokyoError).Bold(true).Render("[!]")
	loc := a.FilePath
	if a.StartLine > 0 {
		loc = fmt.Sprintf("%s:%d", a.FilePath, a.StartLine)
	}

	return fmt.Sprintf("  %s #%d %s: %s", badge, a.ConfusionId, loc, a.ReasoningNote)
}

func renderAlertsSection(out io.Writer, alerts []HeatmapAiConfusion) {
	if len(alerts) == 0 {
		return
	}

	fmt.Fprintln(out, "  "+termout.TokyoHeader("Active AI Confusions & Unmet Expectations"))
	fmt.Fprintln(out, "  "+renderSectionDivider())
	for _, a := range alerts {
		fmt.Fprintln(out, formatAlertLine(a))
	}

	fmt.Fprintln(out)
}

func renderDrillDownSection(out io.Writer, cmds SuggestedCommandsPayload) {
	fmt.Fprintln(out, "  "+termout.TokyoHeader("Drill-Down Inspection Commands"))
	fmt.Fprintln(out, "  "+renderSectionDivider())

	lblTree := termout.SafeTokyoStyle(termout.TokyoAccent).Render("• Tree View   : ")
	lblSearch := termout.SafeTokyoStyle(termout.TokyoAccent).Render("• Search Logs : ")
	lblConf := termout.SafeTokyoStyle(termout.TokyoAccent).Render("• Confusions  : ")

	fmt.Fprintf(out, "  %s%s\n", lblTree, cmds.TreeViewCommand)
	fmt.Fprintf(out, "  %s%s\n", lblSearch, cmds.SearchErrorCommand)
	fmt.Fprintf(out, "  %s%s\n", lblConf, cmds.AiConfusionSyncCommand)
}

func renderStatusText(out io.Writer, s store.HeatmapStatusSummary) {
	inv, _ := DiscoverHeatmapTasks(".", 5)
	activeSlug := ResolveActiveTaskSlug(s.TaskSlug, inv)
	cmds := BuildSuggestedCommands(activeSlug)
	alerts := CollectConfusionAlerts(nil, activeSlug)

	renderStatusDashboard(out, s, inv, alerts, cmds)
}
