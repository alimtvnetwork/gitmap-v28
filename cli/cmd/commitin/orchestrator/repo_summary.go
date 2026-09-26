package orchestrator

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/workspace"
)

// RepoSummaryItem represents the completed summary for one replayed repository.
type RepoSummaryItem struct {
	OrderIndex      int    `json:"orderIndex"`
	Name            string `json:"name"`
	OriginalInput   string `json:"originalInput"`
	StagedPath      string `json:"stagedPath"`
	CommitsReplayed int    `json:"commitsReplayed"`
	CommitsSkipped  int    `json:"commitsSkipped"`
	CommitsFailed   int    `json:"commitsFailed"`
	PRsProcessed    int    `json:"prsProcessed"`
	HeadSha         string `json:"headSha"`
	IsPushed        bool   `json:"isPushed"`
	Status          string `json:"status"`
	CompletedAt     string `json:"completedAt"`
	SummaryFile     string `json:"summaryFile,omitempty"`
}

// MigrationIndexSummary aggregates all completed repositories for user inspection.
type MigrationIndexSummary struct {
	Target         string            `json:"target"`
	TotalInputs    int               `json:"totalInputs"`
	CompletedCount int               `json:"completedCount"`
	TotalCommits   int               `json:"totalCommits"`
	TotalPRs       int               `json:"totalPRs"`
	CurrentHead    string            `json:"currentHead"`
	Status         string            `json:"status"`
	UpdatedAt      string            `json:"updatedAt"`
	SummaryDir     string            `json:"summaryDir,omitempty"`
	Repositories   []RepoSummaryItem `json:"repositories"`
}

func finalizeOneInput(
	ctx *runContext,
	staged workspace.StagedInput,
	created, skipped, failed, prs int,
	stdout io.Writer,
) {
	checkoutTargetHead(ctx.Source.Path)
	isPushed := maybePushImmediate(ctx, staged, created, stdout)
	recordRepoSummary(ctx, staged, created, skipped, failed, prs, isPushed, stdout)
}

func checkoutTargetHead(targetPath string) {
	_ = exec.Command("git", "-C", targetPath, "checkout", "-f", "HEAD").Run()
}

func maybePushImmediate(ctx *runContext, staged workspace.StagedInput, created int, stdout io.Writer) bool {
	if !ctx.Raw.IsPushImmediate || ctx.Raw.IsDryRun {
		return false
	}
	ensureSSHRemote(ctx.Source.Path)
	out, err := exec.Command("git", "-C", ctx.Source.Path, "push", "-u", "origin", "main", "--force").CombinedOutput()
	if err != nil {
		if stdout != nil {
			fmt.Fprintf(stdout, "\n  ⚠️ Push warning for %s: %s\n\n", staged.Input.Original, strings.TrimSpace(string(out)))
		}
		return false
	}
	if stdout != nil {
		fmt.Fprintf(stdout, "\n  ✔ Staged repo %d: %s pushed to origin/main (%d commits)\n\n", staged.Input.OrderIndex, staged.Input.Original, created)
	}
	return true
}

func ensureSSHRemote(repoPath string) {
	out, err := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin").CombinedOutput()
	if err != nil {
		return
	}
	urlStr := strings.TrimSpace(string(out))
	if strings.HasPrefix(urlStr, "https://github.com/") {
		trimmed := strings.TrimPrefix(urlStr, "https://github.com/")
		if !strings.HasSuffix(trimmed, ".git") {
			trimmed += ".git"
		}
		sshURL := "git@github.com:" + trimmed
		_ = exec.Command("git", "-C", repoPath, "remote", "set-url", "origin", sshURL).Run()
	}
}

func recordRepoSummary(
	ctx *runContext,
	staged workspace.StagedInput,
	created, skipped, failed, prs int,
	isPushed bool,
	stdout io.Writer,
) {
	dir := resolveSummaryDir(ctx)
	if dir == "" {
		return
	}
	_ = os.MkdirAll(dir, 0o755)
	headSha := readTargetHeadSha(ctx.Source.Path)
	fname := cleanRepoFilename(staged.Input.OrderIndex, staged.Input.Original)
	status := "completed"
	if isPushed {
		status = "pushed and completed"
	}
	item := RepoSummaryItem{
		OrderIndex:      staged.Input.OrderIndex,
		Name:            filepath.Base(staged.WorkPath),
		OriginalInput:   staged.Input.Original,
		StagedPath:      staged.WorkPath,
		CommitsReplayed: created,
		CommitsSkipped:  skipped,
		CommitsFailed:   failed,
		PRsProcessed:    prs,
		HeadSha:         headSha,
		IsPushed:        isPushed,
		Status:          status,
		CompletedAt:     time.Now().UTC().Format(time.RFC3339),
		SummaryFile:     fname,
	}
	writeRepoSummaryFile(dir, fname, item)
	updateSummaryIndex(dir, ctx, item)
	if stdout != nil {
		fmt.Fprintf(stdout, "  ✔ Repo %d summary: %s (Status: %s, Commits: %d, PRs: %d)\n", staged.Input.OrderIndex, filepath.Join(dir, fname), status, created, prs)
	}
}

func resolveSummaryDir(ctx *runContext) string {
	if ctx.Raw.SummaryDir != "" {
		return ctx.Raw.SummaryDir
	}
	if ctx.Paths != nil && ctx.Paths.CommitInRoot != "" {
		return filepath.Join(ctx.Paths.CommitInRoot, "summaries")
	}

	return filepath.Join(os.TempDir(), "gitmap-summaries", fmt.Sprintf("%d", ctx.RunID))
}

func printFinalSummaryLocation(ctx *runContext, stdout io.Writer) {
	dir := resolveSummaryDir(ctx)
	if dir == "" || stdout == nil {
		return
	}
	fmt.Fprintf(stdout, "\n  📋 Migration Summary Index: %s\n\n", filepath.Join(dir, "index.json"))
}

func cleanRepoFilename(orderIndex int, original string) string {
	base := filepath.Base(original)
	base = strings.TrimSuffix(base, ".git")
	base = strings.ReplaceAll(base, "https://github.com/", "")
	base = strings.ReplaceAll(base, "http://", "")
	base = strings.ReplaceAll(base, "https://", "")
	base = strings.ReplaceAll(base, "/", "-")

	return fmt.Sprintf("%02d-%s.json", orderIndex, base)
}

func writeRepoSummaryFile(dir, fname string, item RepoSummaryItem) {
	data, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, fname), data, 0o644)
}

func updateSummaryIndex(dir string, ctx *runContext, item RepoSummaryItem) {
	idxPath := filepath.Join(dir, "index.json")
	idx := loadExistingIndex(idxPath, ctx)
	idx.CurrentHead = item.HeadSha
	idx.UpdatedAt = item.CompletedAt
	idx.SummaryDir = dir
	idx.Repositories = appendOrUpdateRepo(idx.Repositories, item)
	idx.CompletedCount = len(idx.Repositories)
	idx.TotalCommits = 0
	idx.TotalPRs = 0
	allPushed := true
	for _, r := range idx.Repositories {
		idx.TotalCommits += r.CommitsReplayed
		idx.TotalPRs += r.PRsProcessed
		if !r.IsPushed {
			allPushed = false
		}
	}
	if allPushed && len(idx.Repositories) > 0 {
		idx.Status = "pushed and completed"
	} else {
		idx.Status = "completed"
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err == nil {
		_ = os.WriteFile(idxPath, data, 0o644)
	}
}

func loadExistingIndex(path string, ctx *runContext) MigrationIndexSummary {
	var idx MigrationIndexSummary
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &idx)
	}
	if idx.Target == "" {
		idx.Target = ctx.Source.Path
		idx.TotalInputs = len(ctx.Raw.Inputs)
	}

	return idx
}

func appendOrUpdateRepo(list []RepoSummaryItem, item RepoSummaryItem) []RepoSummaryItem {
	for i, r := range list {
		if r.OrderIndex == item.OrderIndex {
			list[i] = item
			return list
		}
	}

	return append(list, item)
}

func readTargetHeadSha(targetPath string) string {
	out, err := exec.Command("git", "-C", targetPath, "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
