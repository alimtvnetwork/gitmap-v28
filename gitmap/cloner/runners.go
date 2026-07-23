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

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
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

	safePull := opts.SafePull
	if !safePull && hasExistingRepos(records, targetDir) {
		safePull = true
		fmt.Print(constants.MsgAutoSafePull)
	}

	cache := LoadCloneCache(targetDir)
	progress := NewProgress(len(records), opts.Quiet)

	workers := normalizeWorkers(opts.MaxConcurrency, len(records))

	var summary model.CloneSummary
	if workers > 1 {
		fmt.Fprintf(os.Stderr, constants.MsgCloneConcurrencyEnabledFmt, workers)
		summary = runConcurrent(records, targetDir, safePull, workers, progress, cache)
	} else {
		summary = runSequential(records, targetDir, safePull, progress, cache)
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
func runSequential(records []model.ScanRecord, targetDir string, safePull bool,
	progress *Progress, cache *CloneCache) model.CloneSummary {
	summary := model.CloneSummary{}
	for _, rec := range records {
		progress.Begin(repoDisplayName(rec))

		dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
		if cache.IsUpToDate(rec, dest) {
			result := model.CloneResult{Record: rec, Success: true}
			progress.Skip(result)
			summary = updateSummarySkipped(summary, result)
			continue
		}

		result := cloneOrPullOne(rec, targetDir, safePull)
		trackResult(progress, result, rec, targetDir, safePull)
		summary = updateSummary(summary, result)

		if result.Success {
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
func trackResult(p *Progress, result model.CloneResult, rec model.ScanRecord, targetDir string, safePull bool) {
	if result.Success {
		pulled := safePull && isGitRepo(filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath)))
		p.Done(result, pulled)

		return
	}

	p.Fail(result)
}
