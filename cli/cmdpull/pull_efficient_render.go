package cmdpull

import (
	"fmt"
	"io"
	"os"
)

const (
	minConciseRepoColWidth = 26
	conciseReservedColGap  = 30
)

// renderConciseActiveResults prints bullet list with aligned columns to stdout.
func renderConciseActiveResults(states []*PullRepoState) {
	RenderConciseActiveResultsTo(os.Stdout, states)
}

// RenderConciseActiveResultsTo renders concise repo states to the provided writer.
func RenderConciseActiveResultsTo(w io.Writer, states []*PullRepoState) {
	fmt.Fprintln(w)
	colWidth := ResolveConciseRepoColWidth(states)
	for _, s := range states {
		statusLabel := ResolveRepoStatusLabel(s.Changes)
		fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, s.RepoName, statusLabel))
	}
}

// ResolveConciseRepoColWidth dynamically calculates the repo column width to prevent overflow.
func ResolveConciseRepoColWidth(states []*PullRepoState) int {
	colWidth := minConciseRepoColWidth
	for _, s := range states {
		if len(s.RepoName) > colWidth {
			colWidth = len(s.RepoName)
		}
	}

	termWidth := detectTerminalWidth()
	maxAllowed := termWidth - conciseReservedColGap
	if maxAllowed > minConciseRepoColWidth && colWidth > maxAllowed {
		colWidth = maxAllowed
	}

	return colWidth
}

// ResolveRepoStatusLabel normalizes empty or synced changes to up-to-date.
func ResolveRepoStatusLabel(changes string) string {
	if changes == "" || changes == "synced" {
		return "up-to-date"
	}
	return changes
}

// FormatConciseActiveResultLine formats a single row with dynamic padding and clean spacing.
func FormatConciseActiveResultLine(colWidth int, repoName, statusLabel string) string {
	displayName := repoName
	if len(displayName) > colWidth && colWidth > 5 {
		displayName = displayName[:colWidth-3] + "..."
	}
	return fmt.Sprintf("    • %-*s  %s", colWidth, displayName, statusLabel)
}

