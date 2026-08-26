// Package cmd — pull_table_row.go renders an individual rich pull table row.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func (l *PullTableLayout) PrintRow(r model.PullTableRow) {
	statusFmt := formatPullStatus(r.PullStatus, r.IsDirty)
	fmt.Printf("  %-*s   %-*s   %-*s   %-*s   %-*s   %s\n",
		l.MaxRepo, r.RepoName,
		l.MaxBranch, r.Branch,
		l.MaxSHA, r.LastSHA,
		l.MaxPR, r.PRStatus,
		l.MaxStatus, statusFmt,
		r.Duration,
	)
}
