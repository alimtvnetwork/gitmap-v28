package cmdpull

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const (
	minConciseRepoColWidth = 26
	conciseReservedColGap  = 30
)

// renderConciseActiveResults prints bullet list with aligned columns to stdout.
func renderConciseActiveResults(states []*PullRepoState, allRecords ...[]model.ScanRecord) {
	RenderConciseActiveResultsTo(os.Stdout, states, allRecords...)
}

// RenderConciseActiveResultsTo renders concise repo states to the provided writer.
func RenderConciseActiveResultsTo(w io.Writer, states []*PullRepoState, allRecords ...[]model.ScanRecord) {
	fmt.Fprintln(w)
	colWidth := ResolveConciseRepoColWidth(states, allRecords...)
	for _, s := range states {
		statusLabel := ResolveRepoStatusLabel(s.Changes)
		fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, s.RepoName, statusLabel))
	}
}

// ResolveConciseRepoColWidth dynamically calculates the repo column width to prevent overflow.
func ResolveConciseRepoColWidth(states []*PullRepoState, allRecords ...[]model.ScanRecord) int {
	colWidth := minConciseRepoColWidth
	if len(allRecords) > 0 {
		for _, r := range allRecords[0] {
			if len(r.RepoName) > colWidth {
				colWidth = len(r.RepoName)
			}
		}
	}
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

// StyleRepoStatusLabel decorates status labels with ANSI colors for fast visual triage.
func StyleRepoStatusLabel(statusLabel string) string {
	switch statusLabel {
	case "up-to-date":
		return constants.ColorGreen + "up-to-date" + constants.ColorReset
	case "dirty":
		return constants.ColorYellow + "dirty" + constants.ColorReset
	case "failed":
		return constants.ColorRed + "failed" + constants.ColorReset
	}

	if strings.HasPrefix(statusLabel, "+") {
		return formatDiffStatsColor(statusLabel)
	}

	return statusLabel
}

func formatDiffStatsColor(statusLabel string) string {
	parts := strings.SplitN(statusLabel, " ", 2)
	diffPart := parts[0]
	slashIdx := strings.Index(diffPart, "/")
	if slashIdx == -1 {
		return constants.ColorGreen + statusLabel + constants.ColorReset
	}

	ins := diffPart[:slashIdx]
	del := diffPart[slashIdx:]
	coloredDiff := constants.ColorGreen + ins + constants.ColorReset + constants.ColorRed + del + constants.ColorReset
	if len(parts) > 1 {
		return coloredDiff + " " + constants.ColorDim + parts[1] + constants.ColorReset
	}

	return coloredDiff
}

// FormatConciseActiveResultLine formats a single row with dynamic padding and clean spacing.
func FormatConciseActiveResultLine(colWidth int, repoName, statusLabel string) string {
	displayName := repoName
	if len(displayName) > colWidth && colWidth > 5 {
		displayName = displayName[:colWidth-3] + "..."
	}
	styledStatus := StyleRepoStatusLabel(statusLabel)
	return fmt.Sprintf("    • %-*s  %s", colWidth, displayName, styledStatus)
}

// FormatWrappedInactiveList formats repository names wrapped cleanly at commas without splitting tokens.
func FormatWrappedInactiveList(names []string, indent string, maxLineLen int) string {
	if len(names) == 0 {
		return ""
	}
	if maxLineLen <= len(indent)+10 {
		maxLineLen = 120
	}

	var sb strings.Builder
	currentLineLen := len(indent)
	sb.WriteString(indent)

	for i, name := range names {
		token := name
		if i < len(names)-1 {
			token += ", "
		}

		if currentLineLen+len(token) > maxLineLen && currentLineLen > len(indent) {
			sb.WriteString("\n")
			sb.WriteString(indent)
			currentLineLen = len(indent)
		}

		sb.WriteString(token)
		currentLineLen += len(token)
	}

	return sb.String()
}
