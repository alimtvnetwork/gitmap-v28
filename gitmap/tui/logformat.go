package tui

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// formatLogRow renders a single command history list row.
func formatLogRow(e model.CommandHistoryRecord) string {
	return fmt.Sprintf("%-16s %-10s %-30s %-10s %-6d %s",
		e.Command, truncateStr(e.Alias, 10),
		truncateStr(e.Args, 30), formatDurationMs(e.DurationMs),
		e.ExitCode, e.StartedAt)
}

// formatDurationMs converts milliseconds to a human-readable string.
func formatDurationMs(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}

	seconds := float64(ms) / 1000.0
	if seconds < 60 {
		return fmt.Sprintf("%.1fs", seconds)
	}

	minutes := int(seconds) / 60
	secs := int(seconds) % 60

	return fmt.Sprintf("%dm%ds", minutes, secs)
}
