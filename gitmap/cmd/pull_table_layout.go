package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
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
		MaxRepo:     10,
		MaxBranch:   8,
		MaxLatestBr: 13,
		MaxSHA:      7,
		MaxPR:       10,
		MaxStatus:   8,
		Rows:        rows,
	}

	for _, r := range rows {
		if len(r.RepoName) > layout.MaxRepo {
			layout.MaxRepo = len(r.RepoName)
		}
		if len(r.Branch) > layout.MaxBranch {
			layout.MaxBranch = len(r.Branch)
		}
		if len(r.LatestBranch) > layout.MaxLatestBr {
			layout.MaxLatestBr = len(r.LatestBranch)
		}
		if len(r.LastSHA) > layout.MaxSHA {
			layout.MaxSHA = len(r.LastSHA)
		}
		if len(r.PRStatus) > layout.MaxPR {
			layout.MaxPR = len(r.PRStatus)
		}
		if len(r.PullStatus) > layout.MaxStatus {
			layout.MaxStatus = len(r.PullStatus)
		}
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
	fmt.Println("  ---------------------------------------------------------------------------------------------------------")
}
