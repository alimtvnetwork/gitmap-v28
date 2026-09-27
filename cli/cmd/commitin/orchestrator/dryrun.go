package orchestrator

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/workspace"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// DryRunSlot captures per-slot metadata for the Array Async Pool dry-run engine.
type DryRunSlot struct {
	Input       workspace.ResolvedInput
	WorkPath    string
	IsReady     bool
	IsFound     bool
	CommitCount int
	Status      string
	Err         error
}

func executeDryRunArrayAsyncPool(ctx *runContext, inputs []workspace.ResolvedInput, stdout io.Writer) int {
	total := len(inputs)
	if total == 0 {
		return constants.CommitInExitOk
	}
	fmt.Fprintf(stdout, "\n  %s=== Dry-Run Array Async Pool Probing (%d Repositories) ===%s\n\n", constants.ColorCyan, total, constants.ColorReset)
	slots := make([]DryRunSlot, total)
	var mu sync.Mutex
	dispatchDryRunProbers(ctx, inputs, slots, &mu)
	streamOrderedDryRunSlots(ctx, stdout, inputs, slots, &mu)
	return constants.CommitInExitOk
}

func dispatchDryRunProbers(ctx *runContext, inputs []workspace.ResolvedInput, slots []DryRunSlot, mu *sync.Mutex) {
	sem := make(chan struct{}, 16)
	for i, in := range inputs {
		go probeDryRunWorker(ctx, in, i, slots, sem, mu)
	}
}

func probeDryRunWorker(ctx *runContext, in workspace.ResolvedInput, idx int, slots []DryRunSlot, sem chan struct{}, mu *sync.Mutex) {
	sem <- struct{}{}
	defer func() { <-sem }()
	res := probeSingleDryRunSlot(ctx, in)
	mu.Lock()
	slots[idx] = res
	mu.Unlock()
}

func probeSingleDryRunSlot(ctx *runContext, in workspace.ResolvedInput) DryRunSlot {
	if in.Kind == constants.CommitInInputKindLocalFolder || in.Kind == constants.CommitInInputKindVersionedSibling {
		return probeLocalFolderSlot(in)
	}
	cached := workspace.FindCachedCloneDir(ctx.TempDir, in)
	if workspace.HasGitMetadata(cached) {
		count := countCommitsAt(cached)
		return DryRunSlot{Input: in, WorkPath: cached, IsReady: true, IsFound: true, CommitCount: count, Status: "cached"}
	}
	return probeRemoteRepoSlot(in)
}

func probeLocalFolderSlot(in workspace.ResolvedInput) DryRunSlot {
	if !workspace.HasGitMetadata(in.AbsPath) {
		return DryRunSlot{Input: in, IsReady: true, IsFound: false, Status: "not found", Err: fmt.Errorf("no git repository")}
	}
	count := countCommitsAt(in.AbsPath)
	return DryRunSlot{Input: in, WorkPath: in.AbsPath, IsReady: true, IsFound: true, CommitCount: count, Status: "local"}
}

func probeRemoteRepoSlot(in workspace.ResolvedInput) DryRunSlot {
	cmd := exec.Command("git", "ls-remote", "--heads", in.URL)
	out, err := cmd.CombinedOutput()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		return DryRunSlot{Input: in, IsReady: true, IsFound: false, Status: "not found", Err: err}
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return DryRunSlot{Input: in, IsReady: true, IsFound: true, CommitCount: len(lines), Status: "remote"}
}

func countCommitsAt(dir string) int {
	cmd := exec.Command("git", "-C", dir, "rev-list", "--count", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0
	}
	count, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return count
}

func streamOrderedDryRunSlots(ctx *runContext, stdout io.Writer, inputs []workspace.ResolvedInput, slots []DryRunSlot, mu *sync.Mutex) {
	total := len(inputs)
	readyCount, totalCommits := 0, 0
	for cursor := 0; cursor < total; cursor++ {
		for {
			mu.Lock()
			slot := slots[cursor]
			mu.Unlock()
			if slot.IsReady {
				readyCount, totalCommits = renderDryRunSlotOutput(stdout, cursor, total, slot, readyCount, totalCommits)
				break
			}
			time.Sleep(15 * time.Millisecond)
		}
	}
	ctx.Counters.Skipped = totalCommits
	printDryRunSummary(stdout, total, readyCount, totalCommits)
}

func renderDryRunSlotOutput(stdout io.Writer, cursor, total int, slot DryRunSlot, readyCount, totalCommits int) (int, int) {
	if slot.IsFound {
		fmt.Fprintf(stdout, "  [DRY-RUN] [%d/%d] %s: %d commits (%s) [OK]\n", cursor+1, total, slot.Input.Original, slot.CommitCount, slot.Status)
		return readyCount + 1, totalCommits + slot.CommitCount
	}
	fmt.Fprintf(stdout, "  [DRY-RUN] [%d/%d] %s: NOT FOUND (skipped)\n", cursor+1, total, slot.Input.Original)
	return readyCount, totalCommits
}

func printDryRunSummary(stdout io.Writer, total, readyCount, totalCommits int) {
	fmt.Fprintf(stdout, "\n  %s✓%s [DRY-RUN COMPLETE] %d/%d repositories verified, %d total commits discovered.\n\n",
		constants.ColorGreen, constants.ColorReset, readyCount, total, totalCommits)
}
