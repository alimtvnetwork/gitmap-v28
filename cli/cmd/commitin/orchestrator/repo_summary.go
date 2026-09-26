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
	HeadSha         string `json:"headSha"`
	Status          string `json:"status"`
	CompletedAt     string `json:"completedAt"`
	SummaryFile     string `json:"summaryFile,omitempty"`
}

// MigrationIndexSummary aggregates all completed repositories for user inspection.
type MigrationIndexSummary struct {
	Target         string            `json:"target"`
	TotalInputs    int               `json:"totalInputs"`
	CompletedCount int               `json:"completedCount"`
	CurrentHead    string            `json:"currentHead"`
	UpdatedAt      string            `json:"updatedAt"`
	Repositories   []RepoSummaryItem `json:"repositories"`
}

func finalizeOneInput(
	ctx *runContext,
	staged workspace.StagedInput,
	created, skipped, failed int,
	stdout io.Writer,
) {
	checkoutTargetHead(ctx.Source.Path)
	recordRepoSummary(ctx, staged, created, skipped, failed)
	maybePushImmediate(ctx, staged, created, stdout)
}

func checkoutTargetHead(targetPath string) {
	_ = exec.Command("git", "-C", targetPath, "checkout", "-f", "HEAD").Run()
}

func maybePushImmediate(ctx *runContext, staged workspace.StagedInput, created int, stdout io.Writer) {
	if !ctx.Raw.IsPushImmediate || ctx.Raw.IsDryRun {
		return
	}
	ensureSSHRemote(ctx.Source.Path)
	out, err := exec.Command("git", "-C", ctx.Source.Path, "push", "-u", "origin", "main", "--force").CombinedOutput()
	if err != nil {
		if stdout != nil {
			fmt.Fprintf(stdout, "\n  ⚠️ Push warning for %s: %s\n\n", staged.Input.Original, strings.TrimSpace(string(out)))
		}
		return
	}
	if stdout != nil {
		fmt.Fprintf(stdout, "\n  ✔ Staged repo %d: %s pushed to origin/main (%d commits)\n\n", staged.Input.OrderIndex, staged.Input.Original, created)
	}
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

func recordRepoSummary(ctx *runContext, staged workspace.StagedInput, created, skipped, failed int) {
	dir := resolveSummaryDir(ctx)
	if dir == "" {
		return
	}
	_ = os.MkdirAll(dir, 0o755)
	headSha := readTargetHeadSha(ctx.Source.Path)
	fname := cleanRepoFilename(staged.Input.OrderIndex, staged.Input.Original)
	item := RepoSummaryItem{
		OrderIndex:      staged.Input.OrderIndex,
		Name:            filepath.Base(staged.WorkPath),
		OriginalInput:   staged.Input.Original,
		StagedPath:      staged.WorkPath,
		CommitsReplayed: created,
		CommitsSkipped:  skipped,
		CommitsFailed:   failed,
		HeadSha:         headSha,
		Status:          "completed",
		CompletedAt:     time.Now().UTC().Format(time.RFC3339),
		SummaryFile:     fname,
	}
	writeRepoSummaryFile(dir, fname, item)
	updateSummaryIndex(dir, ctx, item)
}

func resolveSummaryDir(ctx *runContext) string {
	if ctx.Raw.SummaryDir != "" {
		return ctx.Raw.SummaryDir
	}
	if ctx.Paths != nil && ctx.Paths.CommitInRoot != "" {
		return filepath.Join(ctx.Paths.CommitInRoot, "summaries")
	}

	return ""
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
	idx.Repositories = appendOrUpdateRepo(idx.Repositories, item)
	idx.CompletedCount = len(idx.Repositories)
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
