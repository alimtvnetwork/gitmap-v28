package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// watchSnapshot holds a single repo's watch status.
type watchSnapshot struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Branch string `json:"branch"`
	Status string `json:"status"`
	Ahead  int    `json:"ahead"`
	Behind int    `json:"behind"`
	Stash  int    `json:"stash"`
}

// watchSummary holds aggregate watch counts.
type watchSummary struct {
	Total  int `json:"total"`
	Dirty  int `json:"dirty"`
	Behind int `json:"behind"`
	Stash  int `json:"stash"`
}

// collectAllStatuses gathers status for all repos.
func collectAllStatuses(records []model.ScanRecord, noFetch bool) []watchSnapshot {
	if !noFetch {
		fetchAllRemotes(records)
	}

	snapshots := make([]watchSnapshot, 0, len(records))

	for _, rec := range records {
		snap := collectOneStatus(rec)
		snapshots = append(snapshots, snap)
	}

	return snapshots
}

// collectOneStatus gathers status for a single repo.
func collectOneStatus(rec model.ScanRecord) watchSnapshot {
	snap := watchSnapshot{
		Name: rec.RepoName,
		Path: rec.AbsolutePath,
	}

	_, err := os.Stat(rec.AbsolutePath)
	if err != nil {
		snap.Status = "error"

		return snap
	}

	rs := gitutil.Status(rec.AbsolutePath)
	snap.Branch = rs.Branch
	snap.Ahead = rs.Ahead
	snap.Behind = rs.Behind
	snap.Stash = rs.StashCount
	snap.Status = watchStatusLabel(rs.Dirty)

	return snap
}

// watchStatusLabel returns "dirty" or "clean".
func watchStatusLabel(dirty bool) string {
	if dirty {
		return "dirty"
	}

	return "clean"
}

// fetchAllRemotes runs git fetch for each repo (best effort).
func fetchAllRemotes(records []model.ScanRecord) {
	for _, rec := range records {
		_, err := os.Stat(rec.AbsolutePath)
		if err != nil {
			continue
		}

		gitutil.FetchAll(rec.AbsolutePath)
	}
}

// buildWatchSummary aggregates counts from snapshots.
func buildWatchSummary(snapshots []watchSnapshot) watchSummary {
	s := watchSummary{Total: len(snapshots)}

	for _, snap := range snapshots {
		if snap.Status == "dirty" {
			s.Dirty++
		}
		if snap.Behind > 0 {
			s.Behind++
		}
		if snap.Stash > 0 {
			s.Stash++
		}
	}

	return s
}
