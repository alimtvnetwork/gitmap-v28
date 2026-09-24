package cmdpull

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBuildPullEfficientSummary(t *testing.T) {
	states := []*PullRepoState{
		{RepoName: "repo-alpha", Changes: "synced"},
		{RepoName: "repo-beta", Changes: "3 files changed"},
	}
	inactive := []InactiveRepoDetail{
		{RepoName: "repo-gamma", Reason: "0 changes in 24h"},
	}

	summary := buildPullEfficientSummary(3, states, inactive, 1500*time.Millisecond)

	if summary.Total != 3 {
		t.Errorf("expected total 3, got %d", summary.Total)
	}
	if summary.ActiveCount != 2 {
		t.Errorf("expected active 2, got %d", summary.ActiveCount)
	}
	if summary.InactiveCount != 1 {
		t.Errorf("expected inactive 1, got %d", summary.InactiveCount)
	}
	if summary.States[0].Status != "up-to-date" {
		t.Errorf("expected status 'up-to-date', got %s", summary.States[0].Status)
	}
	if summary.States[1].Status != "3 files changed" {
		t.Errorf("expected status '3 files changed', got %s", summary.States[1].Status)
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal summary: %v", err)
	}

	var parsed PullEfficientSummary
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal summary: %v", err)
	}
	if parsed.ActiveCount != 2 || len(parsed.States) != 2 {
		t.Errorf("parsed JSON mismatch: %+v", parsed)
	}
}
