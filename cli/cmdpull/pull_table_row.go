package cmdpull

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func resolveRepoStatusStyle(isDirty bool) lipgloss.Style {
	if isDirty {
		dirtyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c"))

		return dirtyStyle
	}

	cleanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))

	return cleanStyle
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
	renderedStatus, _ := l.renderStatusCol(r.PullStatus, r.IsDirty)
	formattedPR := middleTruncate(r.PRStatus, l.MaxPR, 3)

	sep := "   "
	line := "  " +
		PadVisual(renderedRepo, l.MaxRepo) + sep +
		PadVisual(formattedBranch, l.MaxBranch) + sep +
		PadVisual(formattedLatestBr, l.MaxLatestBr) + sep +
		PadVisual(formattedPR, l.MaxPR) + sep +
		PadVisual(renderedStatus, l.MaxStatus) + sep +
		PadVisual(r.LastSHA, l.MaxSHA) + sep +
		r.Duration

	fmt.Println(line)
}

func (l *PullTableLayout) printCompactRow(r model.PullTableRow) {
	renderedRepo, _ := l.renderRepoCol(r.RepoName, r.IsDirty)
	formattedBranch := formatCombinedBranch(r.Branch, r.LatestBranch, l.MaxBranch)
	renderedStatus, _ := l.renderStatusCol(r.PullStatus, r.IsDirty)
	formattedPR := middleTruncate(r.PRStatus, l.MaxPR, 3)

	sep := "  "
	line := "  " +
		PadVisual(renderedRepo, l.MaxRepo) + sep +
		PadVisual(formattedBranch, l.MaxBranch) + sep +
		PadVisual(formattedPR, l.MaxPR) + sep +
		PadVisual(renderedStatus, l.MaxStatus) + sep +
		PadVisual(r.LastSHA, l.MaxSHA) + sep +
		r.Duration

	fmt.Println(line)
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
