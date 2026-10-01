package cmdpull

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
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

// RenderConciseActiveResultsTo renders concise repo states grouped into categories.
func RenderConciseActiveResultsTo(w io.Writer, states []*PullRepoState, allRecords ...[]model.ScanRecord) {
	colWidth := ResolveConciseRepoColWidth(states, allRecords...)
	var updated, dirty, failed []*PullRepoState
	for _, s := range states {
		if s.IsDirty || s.Changes == "dirty" {
			dirty = append(dirty, s)
		} else if s.ErrorMsg != "" || s.Step == PullStepTypeError || s.Step == PullStepTypeConflict || s.Changes == "failed" {
			failed = append(failed, s)
		} else if isUpdatedRepoState(s) {
			updated = append(updated, s)
		}
	}

	if len(updated) == 0 && len(dirty) == 0 && len(failed) == 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  (all repositories are up-to-date)")
		return
	}

	if len(updated) > 0 {
		renderUpdatedGroup(w, colWidth, updated)
	}
	if len(dirty) > 0 {
		renderDirtyGroup(w, colWidth, dirty)
	}
	if len(failed) > 0 {
		renderFailedGroup(w, colWidth, failed)
	}
}

func isUpdatedRepoState(s *PullRepoState) bool {
	if s.Step == PullStepTypeFastForward || s.Step == PullStepTypeMerging {
		return true
	}
	status := ResolveRepoStatusLabel(s.Changes)
	return status != "up-to-date" && status != "dirty" && status != "failed"
}

func renderUpdatedGroup(w io.Writer, colWidth int, updated []*PullRepoState) {
	fmt.Fprintf(w, "\n  %s%s Updated Repositories (%d):%s\n", constants.ColorGreen, constants.ColorBold, len(updated), constants.ColorReset)
	for _, s := range updated {
		statusLabel := resolveUpdatedStatusLabel(s)
		fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, s.RepoName, statusLabel))
		if s.CommitRange != "" {
			fmt.Fprintf(w, "        %s↳ Commits: %s%s\n", constants.ColorDim, s.CommitRange, constants.ColorReset)
		}
	}
}

func resolveUpdatedStatusLabel(s *PullRepoState) string {
	statusLabel := ResolveRepoStatusLabel(s.Changes)
	if statusLabel == "up-to-date" {
		if s.Step == PullStepTypeFastForward {
			return "fast-forward"
		}
		if s.Step == PullStepTypeMerging {
			return "merged"
		}
	}
	return statusLabel
}

func renderDirtyGroup(w io.Writer, colWidth int, dirty []*PullRepoState) {
	fmt.Fprintf(w, "\n  %s%s Dirty Repositories (%d):%s\n", constants.ColorYellow, constants.ColorBold, len(dirty), constants.ColorReset)
	for _, s := range dirty {
		fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, s.RepoName, "dirty"))
		if s.RepoPath != "" {
			diag := gitutil.InspectDirtyState(s.RepoPath)
			renderItemizedDirtyFiles(w, diag)
		}
		errDetails := ResolvePullErrorDetails(s)
		if errDetails != "" {
			fmt.Fprintf(w, "        %s↳ Reason: %s%s\n", constants.ColorDim, errDetails, constants.ColorReset)
		}
		remHint := ResolvePullRemediationHint(s)
		if remHint != "" {
			fmt.Fprintf(w, "        %s↳ Next Step: %s%s\n", constants.ColorCyan, remHint, constants.ColorReset)
		}
	}
}

func renderItemizedDirtyFiles(w io.Writer, diag gitutil.DirtyDiagnosis) {
	for i, f := range diag.ModifiedFiles {
		if i >= 5 {
			fmt.Fprintf(w, "        %s... and %d more modified files%s\n", constants.ColorDim, len(diag.ModifiedFiles)-5, constants.ColorReset)
			break
		}
		fmt.Fprintf(w, "        %smodified: %s%s\n", constants.ColorDim, f, constants.ColorReset)
	}
	for i, f := range diag.UntrackedFiles {
		if i >= 5 {
			fmt.Fprintf(w, "        %s... and %d more untracked files%s\n", constants.ColorDim, len(diag.UntrackedFiles)-5, constants.ColorReset)
			break
		}
		fmt.Fprintf(w, "        %suntracked: %s%s\n", constants.ColorDim, f, constants.ColorReset)
	}
}

func renderFailedGroup(w io.Writer, colWidth int, failed []*PullRepoState) {
	fmt.Fprintf(w, "\n  %s%s Failed Repositories (%d):%s\n", constants.ColorRed, constants.ColorBold, len(failed), constants.ColorReset)
	for _, s := range failed {
		fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, s.RepoName, "failed"))
		errDetails := ResolvePullErrorDetails(s)
		if errDetails == "" {
			errDetails = "pull execution failed"
		}
		fmt.Fprintf(w, "        %s↳ Reason: %s%s\n", constants.ColorDim, errDetails, constants.ColorReset)
		remHint := ResolvePullRemediationHint(s)
		if remHint == "" {
			remHint = fmt.Sprintf("gitmap status %s or gitmap fix %s", s.RepoName, s.RepoName)
		}
		fmt.Fprintf(w, "        %s↳ Next Step: %s%s\n", constants.ColorCyan, remHint, constants.ColorReset)
	}
}

// ResolveConciseRepoColWidth dynamically calculates the repo column width to prevent overflow.
func ResolveConciseRepoColWidth(states []*PullRepoState, allRecords ...[]model.ScanRecord) int {
	colWidth := minConciseRepoColWidth
	var records []model.ScanRecord
	if len(allRecords) > 0 {
		records = allRecords[0]
	}
	for _, r := range records {
		if len(r.RepoName) > colWidth {
			colWidth = len(r.RepoName)
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
