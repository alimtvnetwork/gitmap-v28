package cmdpull

import (
	"encoding/json"
	"os"
	"time"
)

// PullBatchSummary represents the serialized JSON summary of a pull-all batch operation.
type PullBatchSummary struct {
	Total        int                   `json:"total"`
	PulledCount  int                   `json:"pulledCount"`
	SuccessCount int                   `json:"successCount"`
	FailedCount  int                   `json:"failedCount"`
	States       []PullRepoSummaryItem `json:"states"`
	DurationMs   int64                 `json:"durationMs"`
}

func renderPullBatchJSONSummary(total int, states []*PullRepoState, dur time.Duration) error {
	summary := buildPullBatchSummary(total, states, dur)

	return json.NewEncoder(os.Stdout).Encode(summary)
}

func buildPullBatchSummary(total int, states []*PullRepoState, dur time.Duration) PullBatchSummary {
	items := buildPullRepoSummaryItems(states)
	successCount, failedCount := countBatchStateOutcomes(states)

	return PullBatchSummary{
		Total:        total,
		PulledCount:  len(states),
		SuccessCount: successCount,
		FailedCount:  failedCount,
		States:       items,
		DurationMs:   dur.Milliseconds(),
	}
}

func buildPullRepoSummaryItems(states []*PullRepoState) []PullRepoSummaryItem {
	items := make([]PullRepoSummaryItem, 0, len(states))
	for _, s := range states {
		items = append(items, PullRepoSummaryItem{
			RepoName: s.RepoName,
			Status:   ResolveRepoStatusLabel(s.Changes),
			Changes:  s.Changes,
		})
	}

	return items
}
