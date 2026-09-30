package cmdpull

import (
	"encoding/json"
	"os"
	"time"
)

func renderJSONEfficientResults(total int, states []*PullRepoState, inactive []InactiveRepoDetail, dur time.Duration) error {
	summary := buildPullEfficientSummary(total, states, inactive, dur)
	w := pullBatchJSONWriter
	if w == nil {
		w = os.Stdout
	}
	return json.NewEncoder(w).Encode(summary)
}

func renderJSONInactiveResults(total int, inactive []InactiveRepoDetail) error {
	summary := PullEfficientSummary{
		Total:         total,
		ActiveCount:   0,
		InactiveCount: len(inactive),
		States:        []PullRepoSummaryItem{},
		Inactive:      inactive,
		DurationMs:    0,
	}
	w := pullBatchJSONWriter
	if w == nil {
		w = os.Stdout
	}
	return json.NewEncoder(w).Encode(summary)
}

func buildPullEfficientSummary(total int, states []*PullRepoState, inactive []InactiveRepoDetail, dur time.Duration) PullEfficientSummary {
	items := make([]PullRepoSummaryItem, 0, len(states))
	for _, s := range states {
		status := s.Changes
		if status == "" || status == "synced" {
			status = "up-to-date"
		}
		items = append(items, PullRepoSummaryItem{
			RepoName: s.RepoName,
			Status:   status,
			Changes:  s.Changes,
		})
	}
	return PullEfficientSummary{
		Total:         total,
		ActiveCount:   len(states),
		InactiveCount: len(inactive),
		States:        items,
		Inactive:      inactive,
		DurationMs:    dur.Milliseconds(),
	}
}
