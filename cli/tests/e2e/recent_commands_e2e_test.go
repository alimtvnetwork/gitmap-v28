//go:build e2e || tempe2e

package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/heatmap"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestE2E_RepoCreateAndFlatCommit(t *testing.T) {
	requireE2EEnv(t)

	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "e2e-flat-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.name", "E2E Tester")
	runGit(t, repoDir, "config", "user.email", "e2e@gitmap.local")

	readmePath := filepath.Join(repoDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# E2E Flat Commit Test\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	runGit(t, repoDir, "add", "README.md")
	runGit(t, repoDir, "commit", "-m", "init - create initial flat commit")

	headSHA := strings.TrimSpace(runGit(t, repoDir, "rev-parse", "HEAD"))
	if len(headSHA) != 40 {
		t.Fatalf("expected 40-char commit SHA, got %q", headSHA)
	}

	logMsg := strings.TrimSpace(runGit(t, repoDir, "log", "-1", "--pretty=%B"))
	if !strings.Contains(logMsg, "init - create initial flat commit") {
		t.Fatalf("unexpected commit message: %s", logMsg)
	}
}

func TestE2E_HeatmapJSONResponseAndExpectation(t *testing.T) {
	requireE2EEnv(t)

	rec := store.HeatmapActionLogRecord{
		HeatmapActionLogId: 42,
		TaskSlug:           "e2e-heatmap-task",
		RepoName:           "gitmap-e2e",
		ActionType:         "EDIT",
		FilePath:           "cli/main.go",
		ExpectedOutcome:    "Compile cleanly and satisfy expectations",
		ActualOutcome:      "Compilation succeeded",
		IsExpectationMet:   true,
	}

	fakeDbPath := "D:/work/repo-cache/.heatmap/01-antigravity/gitmap-e2e/heatmap.db"
	env := heatmap.BuildRcLogEnvelope(rec, fakeDbPath, "antigravity", "gitmap-e2e")

	if env.Status != "success" {
		t.Fatalf("expected status 'success', got %s", env.Status)
	}

	if env.ActionId != 42 {
		t.Fatalf("expected action_id 42, got %d", env.ActionId)
	}

	if strings.Contains(env.WhereSaved.RelativeDbPath, "\\") {
		t.Fatalf("expected relativeDbPath with forward slashes, got: %s", env.WhereSaved.RelativeDbPath)
	}

	if len(env.WhereSaved.RepoSecretsUrl) == 0 {
		t.Fatal("expected non-empty repo_secrets_url")
	}

	if env.DbInfo.Driver != "sqlite" || env.DbInfo.Table != "HeatmapActionLog" {
		t.Fatalf("unexpected db_info: %+v", env.DbInfo)
	}

	if env.AiExpectation.Status != "MET" || !env.AiExpectation.IsExpectationMet {
		t.Fatalf("unexpected ai_expectation: %+v", env.AiExpectation)
	}
}

func TestE2E_PeerDelegation(t *testing.T) {
	requireE2EEnv(t)

	node := cluster.ClusterNode{
		ID:    "node-remote-1",
		Alias: "w1",
		OS:    "windows",
	}

	payload := cluster.PeerDelegationPayload{
		TaskID:  "task-e2e-delegation",
		Action:  "heatmap status --json",
		Command: "gitmap heatmap status --json",
	}

	resp, _ := cluster.DispatchPeerDelegation(context.Background(), node, payload)
	if resp == nil {
		t.Fatal("expected non-nil response from DispatchPeerDelegation")
	}

	if resp.TaskID != "task-e2e-delegation" {
		t.Fatalf("expected task-e2e-delegation, got %s", resp.TaskID)
	}

	if resp.NodeAlias != "w1" {
		t.Fatalf("expected w1, got %s", resp.NodeAlias)
	}
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v, out: %s", strings.Join(args, " "), err, string(out))
	}
	return string(out)
}
