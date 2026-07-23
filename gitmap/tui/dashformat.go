package tui

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// formatDashRow renders a single dashboard row with file counts.
func formatDashRow(e statusEntry) string {
	styledStatus := formatStatus(e.Status)
	files := formatFileSummary(e.Untracked, e.Modified, e.Staged)

	return fmt.Sprintf("%-20s %-12s %-8s %-6s %-6s %-6s %s",
		e.Slug, e.Branch, styledStatus,
		formatCount(e.Ahead), formatCount(e.Behind),
		formatCount(e.Stash), files)
}

// formatStatus applies color styling to a status label.
func formatStatus(status string) string {
	switch status {
	case "dirty":
		return styleDirty.Render("dirty")
	case "error":
		return styleDirty.Render("error")
	default:
		return styleClean.Render("clean")
	}
}

// formatCount renders a number or dash for zero.
func formatCount(n int) string {
	if n == 0 {
		return "-"
	}

	return fmt.Sprintf("%d", n)
}

// formatFileSummary builds a compact file count string like "3U 2M 1S".
func formatFileSummary(untracked, modified, staged int) string {
	if untracked+modified+staged == 0 {
		return ""
	}

	parts := ""
	if staged > 0 {
		parts += fmt.Sprintf("%dS ", staged)
	}
	if modified > 0 {
		parts += fmt.Sprintf("%dM ", modified)
	}
	if untracked > 0 {
		parts += fmt.Sprintf("%dU ", untracked)
	}

	return styleHint.Render(parts)
}

// dashHeader returns the formatted column header row.
func dashHeader() string {
	return fmt.Sprintf("  %-4s %-20s %-12s %-8s %-6s %-6s %-6s %s",
		"", constants.TUIColSlug, constants.TUIColBranch,
		constants.TUIColStatus, constants.TUIColAhead,
		constants.TUIColBehind, constants.TUIColStash, "Files")
}

// dashSummary builds the bottom summary line.
func dashSummary(entries []statusEntry) string {
	dirty, behind, stash, errCount := 0, 0, 0, 0
	for _, e := range entries {
		switch e.Status {
		case "dirty":
			dirty++
		case "error":
			errCount++
		}
		if e.Behind > 0 {
			behind++
		}
		if e.Stash > 0 {
			stash++
		}
	}

	ts := time.Now().UTC().Format("15:04:05")
	base := fmt.Sprintf("  %d repos  •  %d dirty  •  %d behind  •  %d stash",
		len(entries), dirty, behind, stash)

	if errCount > 0 {
		base += fmt.Sprintf("  •  %d unreachable", errCount)
	}

	return fmt.Sprintf("%s  •  %s UTC  •  r: refresh", base, ts)
}
