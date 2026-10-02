package cmdpipeline

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func TestIsHistoryAISubcmd(t *testing.T) {
	cases := []struct {
		arg      string
		expected bool
	}{
		{"history-ai", true},
		{"HISTORY-AI", true},
		{"hai", true},
		{"history_ai", true},
		{"history", false},
		{"errors", false},
		{"pe", false},
	}

	for _, tc := range cases {
		actual := IsHistoryAISubcmd(tc.arg)
		if actual != tc.expected {
			t.Errorf("IsHistoryAISubcmd(%q) = %v; expected %v", tc.arg, actual, tc.expected)
		}
	}
}

func TestHasHistoryAISubcmd(t *testing.T) {
	if !HasHistoryAISubcmd([]string{"history-ai", "5"}) {
		t.Errorf("expected true for history-ai in args")
	}
	if !HasHistoryAISubcmd([]string{"pe", "hai", "--repo", "demo"}) {
		t.Errorf("expected true for hai in args")
	}
	if HasHistoryAISubcmd([]string{"pe", "-1"}) {
		t.Errorf("expected false when history-ai is not in args")
	}
}

func TestExtractHistoryAILimit(t *testing.T) {
	cases := []struct {
		args     []string
		expected int
	}{
		{[]string{"history-ai", "3"}, 3},
		{[]string{"history-ai", "10"}, 10},
		{[]string{"history-ai"}, 5},
		{[]string{"pe", "history-ai", "7"}, 7},
		{[]string{"hai", "--json"}, 5},
	}

	for _, tc := range cases {
		actual := extractHistoryAILimit(tc.args)
		if actual != tc.expected {
			t.Errorf("extractHistoryAILimit(%v) = %d; expected %d", tc.args, actual, tc.expected)
		}
	}
}

func TestGenerateHistoryAIMarkdown(t *testing.T) {
	payload := HistoryAIPayload{
		Repo:        "test/repo",
		GeneratedAt: "2026-10-02T03:00:00Z",
		CommitCount: 1,
		TotalErrors: 1,
		Commits: []HistoryAICommit{
			{
				Sha:          "abc1234567890",
				ShortSha:     "abc1234",
				Branch:       "main",
				WorkflowName: "CI",
				RunId:        123456,
				RunUrl:       "https://github.com/test/repo/actions/runs/123456",
				FailedSteps: []HistoryAIFailedStep{
					{
						JobName:        "Lint",
						StepName:       "golangci-lint",
						FailureSummary: "Lint error in file.go",
						ErrorLines:     []string{"file.go:10: linter violation"},
					},
				},
			},
		},
	}

	md := generateHistoryAIMarkdown(payload)
	if !strings.Contains(md, "Pipeline Historical Errors & AI Anti-Mistake Training Dossier") {
		t.Errorf("expected title in markdown")
	}
	if !strings.Contains(md, "abc1234") {
		t.Errorf("expected commit sha in markdown")
	}
	if !strings.Contains(md, "golangci-lint") {
		t.Errorf("expected step name in markdown")
	}
	if !strings.Contains(md, "Code Quality / Linter") {
		t.Errorf("expected failure taxonomy in markdown")
	}
}

func TestCategorizeFailure(t *testing.T) {
	if cat := categorizeFailure("golangci-lint (strict)"); cat != "Code Quality / Linter" {
		t.Errorf("expected Linter category, got %s", cat)
	}
	if cat := categorizeFailure("Policy Check Guard"); cat != "Coding Guideline Violation" {
		t.Errorf("expected Coding Guideline category, got %s", cat)
	}
	if cat := categorizeFailure("go test ./..."); cat != "Unit / Integration Test" {
		t.Errorf("expected Test category, got %s", cat)
	}
	if cat := categorizeFailure("windows-latest / go build"); cat != "Cross-Platform Build" {
		t.Errorf("expected Build category, got %s", cat)
	}
}

func TestGroupRunsByFailingCommits(t *testing.T) {
	runs := []pipelinedb.PipelineRunRecord{
		{RunId: 1, Sha: "sha1", Branch: "main", WorkflowName: "CI"},
		{RunId: 2, Sha: "sha1", Branch: "main", WorkflowName: "Build"}, // duplicate sha
		{RunId: 3, Sha: "sha2", Branch: "feature", WorkflowName: "CI"},
	}

	commits := groupRunsByFailingCommits(nil, "dummy/repo", runs, 2)
	if len(commits) != 2 {
		t.Fatalf("expected 2 distinct commits, got %d", len(commits))
	}
	if commits[0].Sha != "sha1" || commits[1].Sha != "sha2" {
		t.Errorf("unexpected commit shas: %v", commits)
	}
}
