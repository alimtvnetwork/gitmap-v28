package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func resolvePullStatusStyle(pullStatus string) lipgloss.Style {
	isDefaultStatus := pullStatus == "UP_TO_DATE" || pullStatus == "synced"
	if isDefaultStatus {
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))

		return statusStyle
	}
	alertStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c")).Bold(true)

	return alertStyle
}

func resolveRepoStatusStyle(isDirty bool) lipgloss.Style {
	if isDirty {
		dirtyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c"))

		return dirtyStyle
	}
	cleanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))

	return cleanStyle
}

func (l *PullTableLayout) PrintRow(r model.PullTableRow) {
	statusStyle := resolvePullStatusStyle(r.PullStatus)
	repoStyle := resolveRepoStatusStyle(r.IsDirty)

	formattedRepo := formatRepoName(r.RepoName, l.MaxRepo)
	renderedRepo := repoStyle.Render(formattedRepo)
	padRepo := calcAnsiPadding(renderedRepo, l.MaxRepo)

	formattedBranch := formatBranchName(r.Branch, l.MaxBranch)
	formattedLatestBr := formatBranchName(r.LatestBranch, l.MaxLatestBr)

	renderedStatus := statusStyle.Render(r.PullStatus)
	padStatus := calcAnsiPadding(renderedStatus, l.MaxStatus)

	fmt.Printf("  %-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %s\n",
		padRepo, renderedRepo,
		l.MaxBranch, formattedBranch,
		l.MaxLatestBr, formattedLatestBr,
		l.MaxPR, r.PRStatus,
		padStatus, renderedStatus,
		l.MaxSHA, r.LastSHA,
		r.Duration,
	)
}
