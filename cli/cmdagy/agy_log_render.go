package cmdagy

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RenderAgyLogsTable prints a structured terminal table of AGY decision logs.
func RenderAgyLogsTable(logs []store.AgyDecisionLogRecord) {
	fmt.Println()
	fmt.Printf("  %s%s ANTIGRAVITY DECISION AUDIT LOGS (%d entries) %s%s\n",
		constants.ColorCyan, "╔════", len(logs), "════╗", constants.ColorReset)
	fmt.Printf("  %s%-19s  %-10s  %-8s  %-24s  %-9s  %s%s\n",
		constants.ColorWhite, "TIMESTAMP", "NODE", "CMD", "TARGET PROJECT", "STATUS", "DECISION REASON", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 100), constants.ColorReset)
	if len(logs) == 0 {
		fmt.Printf("  %sNo decision logs found matching criteria.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, l := range logs {
		printAgyLogRow(l)
	}
	fmt.Println()
}

func printAgyLogRow(l store.AgyDecisionLogRecord) {
	ts := formatLogTimestamp(l.CreatedAt)
	node := truncateLogStr(l.Node, 10)
	cmd := truncateLogStr(l.Command, 8)
	target := truncateLogStr(l.TargetProject, 24)
	statusColor := resolveLogStatusColor(l.Status)
	status := truncateLogStr(l.Status, 9)
	reason := truncateLogStr(l.DecisionReason, 38)
	fmt.Printf("  %-19s  %-10s  %-8s  %-24s  %s%-9s%s  %s\n",
		ts, node, cmd, target, statusColor, status, constants.ColorReset, reason)
}

func formatLogTimestamp(raw string) string {
	if len(raw) >= 19 {
		t := strings.Replace(raw[:19], "T", " ", 1)
		return t
	}
	return raw
}

func resolveLogStatusColor(status string) string {
	low := strings.ToLower(status)
	if strings.Contains(low, "success") || strings.Contains(low, "ok") {
		return constants.ColorGreen
	}
	if strings.Contains(low, "err") || strings.Contains(low, "fail") {
		return constants.ColorRed
	}
	if strings.Contains(low, "dry") || strings.Contains(low, "skip") {
		return constants.ColorYellow
	}
	return constants.ColorCyan
}

func truncateLogStr(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
