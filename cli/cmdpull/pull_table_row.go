package cmdpull

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func resolveRepoStatusStyle(isDirty bool) lipgloss.Style {
	if isDirty {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c"))
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
}

func (l *PullTableLayout) PrintRow(r model.PullTableRow) {
	if l.IsWide {
		l.printWideRow(r)

		return
	}

	l.printCompactRow(r)
}

func (l *PullTableLayout) printWideRow(r model.PullTableRow) {
	renderedRepo, _ := l.renderRepoCol(r.RepoName, r.IsDirty)
	formattedBranch := formatBranchName(r.Branch, l.MaxBranch)
	formattedLatestBr := formatBranchName(r.LatestBranch, l.MaxLatestBr)
	commitRange := resolveRowCommitRange(r, l.MaxRange)
	changes := resolveRowChanges(r, l.MaxChanges)
	line := l.formatWideRowLine(renderedRepo, formattedBranch, formattedLatestBr, commitRange, changes, r)

	fmt.Println(line)
}

func (l *PullTableLayout) formatWideRowLine(repo, br, latestBr, rng, ch string, r model.PullTableRow) string {
	pr := middleTruncate(r.PRStatus, l.MaxPR, 3)
	status, _ := l.renderStatusCol(r.PullStatus, r.IsDirty)
	sep := "   "

	return "  " +
		PadVisual(repo, l.MaxRepo) + sep +
		PadVisual(br, l.MaxBranch) + sep +
		PadVisual(latestBr, l.MaxLatestBr) + sep +
		PadVisual(rng, l.MaxRange) + sep +
		PadVisual(ch, l.MaxChanges) + sep +
		PadVisual(pr, l.MaxPR) + sep +
		PadVisual(status, l.MaxStatus) + sep +
		r.Duration
}

func (l *PullTableLayout) printCompactRow(r model.PullTableRow) {
	renderedRepo, _ := l.renderRepoCol(r.RepoName, r.IsDirty)
	formattedBranch := formatCombinedBranch(r.Branch, r.LatestBranch, l.MaxBranch)
	commitRange := resolveRowCommitRange(r, l.MaxRange)
	changes := resolveRowChanges(r, l.MaxChanges)
	renderedStatus, _ := l.renderStatusCol(r.PullStatus, r.IsDirty)

	sep := "  "
	line := "  " +
		PadVisual(renderedRepo, l.MaxRepo) + sep +
		PadVisual(formattedBranch, l.MaxBranch) + sep +
		PadVisual(commitRange, l.MaxRange) + sep +
		PadVisual(changes, l.MaxChanges) + sep +
		PadVisual(renderedStatus, l.MaxStatus) + sep +
		r.Duration

	fmt.Println(line)
}

func resolveRowCommitRange(r model.PullTableRow, maxLen int) string {
	val := r.CommitRange
	if val == "" {
		val = r.LastSHA
	}

	return middleTruncate(val, maxLen, 3)
}

func resolveRowChanges(r model.PullTableRow, maxLen int) string {
	val := r.Changes
	if val == "" {
		val = "-"
	}

	return middleTruncate(val, maxLen, 2)
}

func (l *PullTableLayout) renderRepoCol(name string, isDirty bool) (string, int) {
	repoStyle := resolveRepoStatusStyle(isDirty)
	formattedRepo := formatRepoName(name, l.MaxRepo)
	renderedRepo := repoStyle.Render(formattedRepo)
	padRepo := calcAnsiPadding(renderedRepo, l.MaxRepo)

	return renderedRepo, padRepo
}

func (l *PullTableLayout) renderStatusCol(status string, isDirty bool) (string, int) {
	renderedStatus := formatPullStatus(status, isDirty)
	padStatus := calcAnsiPadding(renderedStatus, l.MaxStatus)

	return renderedStatus, padStatus
}
