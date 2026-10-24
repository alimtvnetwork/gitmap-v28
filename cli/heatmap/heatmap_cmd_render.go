package heatmap

import (
	"fmt"
	"io"

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

	return fmt.Sprintf(" [%s:%d-%d]", filePath, startLine, endLine)
}

func renderTreeText(out io.Writer, nodes []store.HeatmapExecutionTreeNode, indent string) {
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

func renderStatusText(out io.Writer, s store.HeatmapStatusSummary) {
	fmt.Fprintf(out, "Heatmap Status Summary (Task: %s)\n--------------------------------------------------\n", s.TaskSlug)
	fmt.Fprintf(out, "Total Actions:    %d\nExpectations Met: %d (%.1f%%)\n\nReasoning States:\n",
		s.TotalActions, s.ExpectationsMetCount, s.ExpectationsMetPct)
	for state, count := range s.ReasoningStateCounts {
		fmt.Fprintf(out, "  %-18s: %d\n", state, count)
	}
	printTopTouchedFiles(out, s.TopTouchedFiles)
}

func printTopTouchedFiles(out io.Writer, files []store.FileTouchCount) {
	if len(files) == 0 {
		return
	}
	fmt.Fprintln(out, "\nTop Touched Files:")
	for _, f := range files {
		fmt.Fprintf(out, "  %-40s (%d touches)\n", f.FilePath, f.Count)
	}
}
