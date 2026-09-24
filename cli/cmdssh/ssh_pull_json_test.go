package cmdssh

import (
	"encoding/json"
	"testing"
)

func TestRenderSingleNodePullJSON_Parse(t *testing.T) {
	summary := remotePullSummary{
		Total:         5,
		ActiveCount:   2,
		InactiveCount: 3,
		States: []remotePullRepoItem{
			{RepoName: "api-backend", Status: "up-to-date", Changes: "synced"},
			{RepoName: "web-client", Status: "2 files changed", Changes: "2 files changed"},
		},
		DurationMs: 450,
	}

	bytes, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	res := FleetNodeResult{
		Alias:   "node-prod",
		IP:      "10.0.0.1",
		Success: true,
		Output:  string(bytes),
	}

	// Should not panic or error
	renderSingleNodePullJSON(res, true)
}

func TestNormalizePullTarget(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "all-nodes"},
		{"ssh", "all-nodes"},
		{"--ssh", "all-nodes"},
		{"all", "all-nodes"},
		{"worker-1", "worker-1"},
	}

	for _, tc := range tests {
		got := normalizePullTarget(tc.input)
		if got != tc.want {
			t.Errorf("normalizePullTarget(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
