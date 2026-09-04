package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

const (
	defaultPullTableMaxRepo     = 20
	defaultPullTableMaxBranch   = 16
	defaultPullTableMaxLatestBr = 18
	defaultPullTableMaxPR       = 10
	defaultPullTableMaxStatus   = 10
	defaultPullTableMaxSHA      = 7
)

type PullTableLayout struct {
	MaxRepo     int
	MaxBranch   int
	MaxLatestBr int
	MaxSHA      int
	MaxPR       int
	MaxStatus   int
	Rows        []model.PullTableRow
}

func NewPullTableLayout(rows []model.PullTableRow) *PullTableLayout {
	layout := &PullTableLayout{
		MaxRepo:     defaultPullTableMaxRepo,
		MaxBranch:   defaultPullTableMaxBranch,
		MaxLatestBr: defaultPullTableMaxLatestBr,
		MaxSHA:      defaultPullTableMaxSHA,
		MaxPR:       defaultPullTableMaxPR,
		MaxStatus:   defaultPullTableMaxStatus,
		Rows:        rows,
	}

	return layout
}

func (l *PullTableLayout) PrintHeader() {
	// Reorder: REPO | BRANCH | LATEST BRANCH | PR/TRACK | STATUS | SHA | TIME
	fmt.Printf("  %-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %s\n",
		l.MaxRepo, "REPO",
		l.MaxBranch, "BRANCH",
		l.MaxLatestBr, "LATEST BRANCH",
		l.MaxPR, "PR/TRACK",
		l.MaxStatus, "STATUS",
		l.MaxSHA, "SHA",
		"TIME",
	)
	fmt.Println("  -------------------------------------------------------------------------------------------------------")
}
