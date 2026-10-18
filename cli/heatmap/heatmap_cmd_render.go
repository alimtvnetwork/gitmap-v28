package heatmap

import (
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
		fStr := r.FilePath
		if r.StartLine > 0 {
			fStr += fmt.Sprintf(":%d", r.StartLine)
		}
		if len(fStr) > 30 {
			fStr = "..." + fStr[len(fStr)-27:]
		}
		fmt.Fprintf(out, "%-6d %-14s %-8s %-30s %s\n", r.HeatmapActionLogId, r.ReasoningState, r.ActionType, fStr, r.ReasoningNotes)
	}
}

func renderTreeText(out io.Writer, nodes []store.HeatmapExecutionTreeNode, indent string) {
	for i, node := range nodes {
		branch, nextIndent := "├── ", indent+"│   "
		if i == len(nodes)-1 {
			branch, nextIndent = "└── ", indent+"    "
		}
		if indent == "" {
			branch, nextIndent = "• ", "  "
		}

		lineInfo := ""
		if node.FilePath != "" {
			lineInfo = fmt.Sprintf(" [%s:%d-%d]", node.FilePath, node.StartLine, node.EndLine)
		}

		fmt.Fprintf(out, "%s%s#%d (%s) %s%s: %s\n", indent, branch, node.HeatmapActionLogId,
			node.ReasoningState, node.ActionType, lineInfo, node.ReasoningNotes)

		if len(node.Children) > 0 {
			renderTreeText(out, node.Children, nextIndent)
		}
	}
}

func renderRcTreeText(out io.Writer, nodes []store.RcHeatmapExecutionTreeNode, indent string) {
	for i, node := range nodes {
		branch, nextIndent := "├── ", indent+"│   "
		if i == len(nodes)-1 {
			branch, nextIndent = "└── ", indent+"    "
		}
		if indent == "" {
			branch, nextIndent = "• ", "  "
		}

		lineInfo := ""
		if node.FilePath != "" {
			lineInfo = fmt.Sprintf(" [%s:%d-%d]", node.FilePath, node.StartLine, node.EndLine)
		}

		fmt.Fprintf(out, "%s%s#%d (%s) %s%s: %s\n", indent, branch, node.HeatmapActionLogId,
			node.ReasoningState, node.ActionType, lineInfo, node.ReasoningNotes)

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

	if len(s.TopTouchedFiles) > 0 {
		fmt.Fprintln(out, "\nTop Touched Files:")
		for _, f := range s.TopTouchedFiles {
			fmt.Fprintf(out, "  %-40s (%d touches)\n", f.FilePath, f.Count)
		}
	}
}

func printHeatmapHelp(out io.Writer) {
	fmt.Fprint(out, `Usage: gitmap heatmap <command> [flags]

Commands:
  log     Log an AI cognitive decision or file action
  update  Update outcome and expectation status of an existing action
  search  Search logged actions by query, task, or file
  tree    Display hierarchical causal execution tree
  status  Summarize action volume, states, and file touches
  tables  List tables and row counts in heatmap SQLite DB
  view    Inspect table rows with --limit, --format, and filters
  query   Run safe read-only SQL SELECT queries against heatmap DB
  rc      Manage repo-cache isolated multi-AI heatmap telemetry
  help    Show this help message

Examples:
  gitmap heatmap log --task t1 --action READ --file cli/store/path.go --state EXPLORATION --expected "Inspect paths"
  gitmap heatmap update --id 1 --outcome "Inspected successfully" --met true
  gitmap heatmap tree --task t1
  gitmap heatmap status --task t1
  gitmap heatmap search "base paths" --json
`)
}
