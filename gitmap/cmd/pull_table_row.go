package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func (l *PullTableLayout) PrintRow(r model.PullTableRow) {
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
	if r.PullStatus != "UP_TO_DATE" && r.PullStatus != "synced" {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c")).Bold(true)
	}

	repoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
	if r.IsDirty {
		repoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c"))
	}

	fmt.Printf("  %-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %s\n",
		l.MaxRepo, repoStyle.Render(r.RepoName),
		l.MaxBranch, r.Branch,
		l.MaxLatestBr, r.LatestBranch,
		l.MaxPR, r.PRStatus,
		l.MaxStatus, statusStyle.Render(r.PullStatus),
		l.MaxSHA, r.LastSHA,
		r.Duration,
	)
}
