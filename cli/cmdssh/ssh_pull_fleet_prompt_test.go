package cmdssh

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestExtractFleetIssues_NodeFailure(t *testing.T) {
	outcomes := []FleetNodePullOutcome{
		{
			NodeName: "w1",
			IP:       "192.168.1.3",
			IsLocal:  false,
			Success:  false,
			ErrorMsg: "connection refused",
		},
		{
			NodeName:  "w3",
			IP:        "192.168.1.12",
			IsLocal:   false,
			Success:   false,
			IsSkipped: true,
			ErrorMsg:  "offline",
		},
	}

	issues := ExtractFleetIssues(outcomes)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !issues[0].HasNodeError || issues[0].NodeError != "connection refused" {
		t.Errorf("expected node error issue, got %+v", issues[0])
	}
}

func TestExtractFleetIssues_DirtyAndFailedRepos(t *testing.T) {
	outcomes := []FleetNodePullOutcome{
		{
			NodeName: "AI-Main",
			IP:       "127.0.0.1",
			IsLocal:  true,
			Success:  true,
			ActiveStates: []fleetPullRepoItem{
				{RepoName: "repo-a", Status: "dirty", Changes: "dirty"},
				{RepoName: "repo-b", Status: "failed", ErrorDetails: "conflict on pull"},
				{RepoName: "repo-c", Status: "up-to-date", Changes: "synced"},
			},
		},
	}

	issues := ExtractFleetIssues(outcomes)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if len(issues[0].RepoIssues) != 2 {
		t.Fatalf("expected 2 repo issues, got %d", len(issues[0].RepoIssues))
	}
	if issues[0].RepoIssues[0].RepoName != "repo-a" || !issues[0].RepoIssues[0].IsDirty {
		t.Errorf("expected repo-a dirty, got %+v", issues[0].RepoIssues[0])
	}
	if issues[0].RepoIssues[1].RepoName != "repo-b" || !issues[0].RepoIssues[1].IsFailed {
		t.Errorf("expected repo-b failed, got %+v", issues[0].RepoIssues[1])
	}
}

func TestExtractFleetIssues_AllClean(t *testing.T) {
	outcomes := []FleetNodePullOutcome{
		{
			NodeName: "AI-Main",
			IP:       "127.0.0.1",
			IsLocal:  true,
			Success:  true,
			ActiveStates: []fleetPullRepoItem{
				{RepoName: "repo-c", Status: "up-to-date", Changes: "synced"},
			},
		},
		{
			NodeName:  "w3",
			IP:        "192.168.1.12",
			IsSkipped: true,
		},
	}

	issues := ExtractFleetIssues(outcomes)
	if len(issues) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(issues))
	}
	if HasAnyFleetIssues(issues) {
		t.Errorf("expected HasAnyFleetIssues to be false")
	}
}

func TestCountTotalFleetIssues(t *testing.T) {
	issues := []FleetNodeIssue{
		{
			NodeName:     "w1",
			HasNodeError: true,
			NodeError:    "offline",
		},
		{
			NodeName: "w2",
			RepoIssues: []FleetRepoIssueItem{
				{RepoName: "a"},
				{RepoName: "b"},
			},
		},
	}

	total := CountTotalFleetIssues(issues)
	if total != 3 {
		t.Errorf("expected total 3 issues, got %d", total)
	}
}

func TestRenderFleetPostPullUnresolvedList_Output(t *testing.T) {
	issues := []FleetNodeIssue{
		{
			NodeName:     "w1",
			IP:           "192.168.1.3",
			IsLocal:      false,
			HasNodeError: true,
			NodeError:    "invalid JSON output",
		},
		{
			NodeName: "AI-Main",
			IP:       "127.0.0.1",
			IsLocal:  true,
			RepoIssues: []FleetRepoIssueItem{
				{RepoName: "my-repo", Reason: "working tree dirty / autostash conflict"},
			},
		},
	}

	out := captureStdout(func() {
		RenderFleetPostPullUnresolvedList(issues)
	})

	if !strings.Contains(out, "require attention after pull") {
		t.Errorf("expected banner text, got: %s", out)
	}
	if !strings.Contains(out, "w1") || !strings.Contains(out, "invalid JSON output") {
		t.Errorf("expected node w1 details, got: %s", out)
	}
	if !strings.Contains(out, "AI-Main") || !strings.Contains(out, "my-repo") {
		t.Errorf("expected AI-Main details, got: %s", out)
	}
}

func TestRenderFleetPostPullMenu_Output(t *testing.T) {
	out := captureStdout(func() {
		RenderFleetPostPullMenu()
	})

	if !strings.Contains(out, "Send commits across affected nodes") {
		t.Errorf("expected Send commits choice, got: %s", out)
	}
	if !strings.Contains(out, "Fix All across affected nodes") {
		t.Errorf("expected Fix All choice, got: %s", out)
	}
	if !strings.Contains(out, "Retry Pull-All") {
		t.Errorf("expected Retry Pull-All choice, got: %s", out)
	}
}

func TestRenderFleetPostPullHelpGuide_Output(t *testing.T) {
	out := captureStdout(func() {
		RenderFleetPostPullHelpGuide()
	})

	if !strings.Contains(out, "Multi-Machine Fleet Remediation Guide") {
		t.Errorf("expected help title, got: %s", out)
	}
	if !strings.Contains(out, "Send Commits") || !strings.Contains(out, "Fix All") {
		t.Errorf("expected strategy explanations, got: %s", out)
	}
}

func TestResolveConnectionForNode(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "w1", IPAddress: "192.168.1.3"},
		{Alias: "w2", IPAddress: "192.168.1.7"},
	}

	c, hasConn := resolveConnectionForNode("w1", "192.168.1.3", conns)
	if !hasConn || c.Alias != "w1" {
		t.Errorf("expected match for w1, got %+v", c)
	}

	c2, hasConn2 := resolveConnectionForNode("unknown", "192.168.1.7", conns)
	if !hasConn2 || c2.Alias != "w2" {
		t.Errorf("expected match for IP 192.168.1.7, got %+v", c2)
	}

	_, hasConn3 := resolveConnectionForNode("w99", "10.0.0.1", conns)
	if hasConn3 {
		t.Errorf("expected no match for w99")
	}
}

func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}
