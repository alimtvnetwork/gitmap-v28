package cmdpull

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

var pullBatchJSONWriter io.Writer = os.Stdout

// SetPullBatchJSONWriter configures the output destination for batch JSON summaries.
func SetPullBatchJSONWriter(w io.Writer) {
	pullBatchJSONWriter = w
}

// ResetPullBatchJSONWriter resets the batch JSON output destination to os.Stdout.
func ResetPullBatchJSONWriter() {
	pullBatchJSONWriter = os.Stdout
}

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
	w := pullBatchJSONWriter
	if w == nil {
		w = os.Stdout
	}
	return json.NewEncoder(w).Encode(summary)
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
		item := buildSinglePullRepoSummaryItem(s)
		items = append(items, item)
	}

	return items
}

func buildSinglePullRepoSummaryItem(s *PullRepoState) PullRepoSummaryItem {
	status := ResolveRepoStatusLabel(s.Changes)
	item := PullRepoSummaryItem{
		RepoName: s.RepoName,
		Status:   status,
		Changes:  s.Changes,
	}
	if status == "failed" || status == "dirty" || s.ErrorMsg != "" {
		item.ErrorDetails = ResolvePullErrorDetails(s)
		item.RemediationHint = ResolvePullRemediationHint(s)
	}
	return item
}
