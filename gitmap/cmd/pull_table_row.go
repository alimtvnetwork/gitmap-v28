// Package cmd — pull_table_row.go renders an individual rich pull table row.
package cmd

import (
	"fmt"
	"sync/atomic"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

var pullRowCounter uint64

func (l *PullTableLayout) PrintRow(r model.PullTableRow) {
	statusFmt := formatPullStatus(r.PullStatus, r.IsDirty)
	idx := atomic.AddUint64(&pullRowCounter, 1) - 1
	pastelColor := constants.ColorCycle[idx%uint64(len(constants.ColorCycle))]

	fmt.Printf("  %s%-*s%s   %-*s   %-*s   %-*s   %-*s   %s\n",
		pastelColor, l.MaxRepo, r.RepoName, constants.ColorReset,
		l.MaxBranch, r.Branch,
		l.MaxSHA, r.LastSHA,
		l.MaxPR, r.PRStatus,
		l.MaxStatus, statusFmt,
		r.Duration,
	)
}
