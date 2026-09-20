package graph

import (
	"strings"
	"testing"
)

func TestRenderExecutionGraph_Empty(t *testing.T) {
	out := RenderExecutionGraph(nil)
	if out != "" {
		t.Errorf("expected empty output for nil events, got %q", out)
	}
}

func TestRenderExecutionGraph_Linear(t *testing.T) {
	events := []GraphEvent{
		{
			CommitSha:  "abc123456",
			BranchName: "main",
			Message:    "feat: initial commit",
		},
		{
			CommitSha:  "def789012",
			BranchName: "main",
			ReleaseTag: "v1.0.0",
			Message:    "release: v1.0.0",
		},
	}

	out := RenderExecutionGraph(events)
	if !strings.Contains(out, "main:") {
		t.Errorf("expected output to contain 'main:', got %s", out)
	}
	if !strings.Contains(out, "abc1234") {
		t.Errorf("expected output to contain short sha 'abc1234', got %s", out)
	}
	if !strings.Contains(out, "v1.0.0") {
		t.Errorf("expected output to contain release tag 'v1.0.0', got %s", out)
	}
}

func TestRenderExecutionGraph_Branched(t *testing.T) {
	events := []GraphEvent{
		{
			CommitSha:  "abc123456",
			BranchName: "main",
			Message:    "feat: add base",
		},
		{
			CommitSha:  "bcd234567",
			BranchName: "feature/auth",
			PRNumber:   42,
			Message:    "feat: add token login",
		},
		{
			CommitSha:  "cde345678",
			BranchName: "main",
			IsMerge:    true,
			PRNumber:   42,
			Message:    "Merge pull request #42 from feature/auth",
		},
	}

	out := RenderExecutionGraph(events)
	if !strings.Contains(out, "branch:") {
		t.Errorf("expected branched graph to show 'branch:' track, got %s", out)
	}
	if !strings.Contains(out, "PR #42") {
		t.Errorf("expected output to mention PR #42, got %s", out)
	}
}
