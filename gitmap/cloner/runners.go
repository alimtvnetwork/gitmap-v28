// Package cloner — runners.go
//
// Dispatcher and the sequential runner, split out of cloner.go to keep
// each file focused (cloner.go = entry points + parsing, runners.go =
// orchestration, concurrent.go = worker-pool path). The dispatcher
// (cloneAll) is the single point that decides between sequential and
// parallel execution and the only writer of the "parallel enabled"
// header line.
package cloner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// cloneAll iterates records and clones each one with progress tracking.
//
// Sequential vs parallel dispatch is decided by opts.MaxConcurrency. Both
// paths share the same Progress + CloneCache instances; the parallel
// path's thread-safety contract lives in concurrent.go.
//
// Hierarchy preservation: every repo lands at filepath.Join(targetDir,
// rec.RelativePath), so the nested folder layout captured by `gitmap
// scan` is reproduced exactly under targetDir — even at MaxConcurrency
// > 1, where ordering of progress lines is no longer sequential.
func cloneAll(records []model.ScanRecord, targetDir string, opts CloneOptions) model.CloneSummary {
	// Apply --default-branch fallback BEFORE existing-repo / safe-pull
	// detection so audit, cache, and progress all see the same patched
	// records the actual `git clone` will receive. No-op when
	// opts.DefaultBranch is empty.
	records = applyDefaultBranchFallback(records, opts.DefaultBranch)

	handleConflicts(&opts, records, targetDir)

	if !opts.SafePull && !opts.Clean && !opts.MissingOnly && hasExistingRepos(records, targetDir) {
		opts.SafePull = true
		fmt.Print(constants.MsgAutoSafePull)
	}

	cache := LoadCloneCache(targetDir)
	progress := NewProgress(len(records), opts.Quiet)

	workers := normalizeWorkers(opts.MaxConcurrency, len(records))

	var summary model.CloneSummary
	if workers > 1 {
		fmt.Fprintf(os.Stderr, constants.MsgCloneConcurrencyEnabledFmt, workers)
		summary = runConcurrent(records, targetDir, opts, workers, progress, cache)
	} else {
		summary = runSequential(records, targetDir, opts, progress, cache)
	}

	// Best-effort cache persistence — never fail the run on write errors.
	_ = cache.Save()

	progress.PrintSummary()

	return summary
}

// normalizeWorkers clamps the requested worker count to a sane range.
// Zero or negative → 1 (sequential). Larger than the work queue → queue
// length (no point spawning idle workers).
func normalizeWorkers(requested, jobs int) int {
	if requested < 1 {
		return 1
	}
	if requested > jobs && jobs > 0 {
		return jobs
	}

	return requested
}

// runSequential is the legacy in-order runner. Kept as a separate
// function so concurrent.go can stay focused on the worker-pool path.
func runSequential(records []model.ScanRecord, targetDir string, opts CloneOptions,
	progress *Progress, cache *CloneCache) model.CloneSummary {
	summary := model.CloneSummary{}
	for _, rec := range records {
		progress.Begin(repoDisplayName(rec))

		dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
		if cache.IsUpToDate(rec, dest) {
			result := model.CloneResult{Record: rec, IsSuccess: true}
			progress.Skip(result)
			summary = updateSummarySkipped(summary, result)
			continue
		}

		result := cloneOrPullOne(rec, targetDir, opts)
		trackResult(progress, result, rec, targetDir, opts.SafePull)
		summary = updateSummary(summary, result)

		if result.IsSuccess {
			cache.Record(rec, dest)
		}
	}

	return summary
}

// repoDisplayName returns a display name for progress output.
func repoDisplayName(rec model.ScanRecord) string {
	if len(rec.RepoName) > 0 {
		return rec.RepoName
	}

	return rec.RelativePath
}

// trackResult updates progress based on clone/pull outcome.
func trackResult(
	p *Progress,
	result model.CloneResult,
	rec model.ScanRecord,
	targetDir string,
	safePull bool,
) {
	if result.IsSuccess {
		pulled := safePull && isGitRepo(filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath)))
		p.Done(result, pulled)

		return
	}

	p.Fail(result)
}

func countConflicts(records []model.ScanRecord, targetDir string) int {
	conflicts := 0
	for _, rec := range records {
		dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
		_, err := os.Stat(dest)
		if err == nil && !isGitRepo(dest) {
			conflicts++
		}
	}
	return conflicts
}

func handleConflicts(opts *CloneOptions, records []model.ScanRecord, targetDir string) {
	if opts.Clean || opts.MissingOnly {
		return
	}
	conflicts := countConflicts(records, targetDir)
	if conflicts > 0 {
		promptAndSetClean(opts, conflicts)
	}
}

func promptAndSetClean(opts *CloneOptions, conflicts int) {
	fmt.Printf("Warning: Detected %d existing directories that are not git repositories.\n", conflicts)
	fmt.Printf("These will fail to clone. Do you want to forcefully clean them? [y/N]: ")
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(strings.TrimSpace(response)) == "y" {
		opts.Clean = true
	} else {
		fmt.Println("Proceeding without --clean. Conflicting directories will fail.")
	}
}
