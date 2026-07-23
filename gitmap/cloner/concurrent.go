// Package cloner — concurrent.go
//
// Bounded worker pool for parallel clone/pull execution. Wired in by
// cloneAll() when CloneOptions.MaxConcurrency > 1; the sequential path
// stays unchanged for the default invocation.
//
// Each worker pulls one record at a time from a buffered job channel,
// performs the clone-or-pull through the same cloneOrPullOne path used
// by the sequential runner, and reports outcomes back through a single
// result channel. The collector goroutine serializes Progress / cache /
// summary updates so the public Progress + CloneCache types only need
// the lightweight mutex they already carry.
//
// Concurrency invariants:
//   - Progress.Begin / .Done / .Skip / .Fail are guarded internally by
//     Progress.mu (see progress.go), so concurrent calls cannot
//     interleave a half-written stderr line.
//   - CloneCache.Record is guarded by CloneCache.mu (see cache.go).
//   - Order of progress lines matches completion order, NOT input order.
//     Manifest rows are still cloned into their recorded RelativePath
//     so the on-disk hierarchy is unaffected.
package cloner

import (
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// cloneJob is the unit of work handed to each worker.
type cloneJob struct {
	rec  model.ScanRecord
	dest string
}

// cloneOutcome is what a worker sends back to the collector.
type cloneOutcome struct {
	rec    model.ScanRecord
	dest   string
	result model.CloneResult
	cached bool
}

// runConcurrent fans the records out across `workers` goroutines and
// returns the same CloneSummary shape as the sequential runner. The
// caller is responsible for picking a sane worker count (>=1).
func runConcurrent(records []model.ScanRecord, targetDir string, safePull bool,
	workers int, progress *Progress, cache *CloneCache) model.CloneSummary {
	jobs := make(chan cloneJob, len(records))
	out := make(chan cloneOutcome, len(records))

	startWorkers(workers, jobs, out, targetDir, safePull)
	enqueueJobs(records, targetDir, cache, progress, jobs, out)
	close(jobs)

	return collectOutcomes(records, targetDir, safePull, progress, cache, out)
}

// startWorkers spins up the worker goroutines.
func startWorkers(workers int, jobs <-chan cloneJob, out chan<- cloneOutcome,
	targetDir string, safePull bool) {
	for i := 0; i < workers; i++ {
		go cloneWorker(jobs, out, targetDir, safePull)
	}
}

// cloneWorker drains the job channel until it closes.
func cloneWorker(jobs <-chan cloneJob, out chan<- cloneOutcome,
	targetDir string, safePull bool) {
	for job := range jobs {
		result := cloneOrPullOne(job.rec, targetDir, safePull)
		out <- cloneOutcome{rec: job.rec, dest: job.dest, result: result}
	}
}

// enqueueJobs short-circuits cache hits (reported synchronously so the
// progress line lands before any worker output) and dispatches the rest
// onto the job channel.
func enqueueJobs(records []model.ScanRecord, targetDir string, cache *CloneCache,
	progress *Progress, jobs chan<- cloneJob, out chan<- cloneOutcome) {
	for _, rec := range records {
		progress.Begin(repoDisplayName(rec))
		dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
		if cache.IsUpToDate(rec, dest) {
			out <- cloneOutcome{
				rec:    rec,
				dest:   dest,
				result: model.CloneResult{Record: rec, Success: true},
				cached: true,
			}

			continue
		}
		jobs <- cloneJob{rec: rec, dest: dest}
	}
}

// collectOutcomes drains the outcome channel, updating progress, cache,
// and summary in the same order outcomes complete.
func collectOutcomes(records []model.ScanRecord, targetDir string, safePull bool,
	progress *Progress, cache *CloneCache, out <-chan cloneOutcome) model.CloneSummary {
	summary := model.CloneSummary{}
	for i := 0; i < len(records); i++ {
		o := <-out
		if o.cached {
			progress.Skip(o.result)
			summary = updateSummarySkipped(summary, o.result)

			continue
		}
		trackResult(progress, o.result, o.rec, targetDir, safePull)
		summary = updateSummary(summary, o.result)
		if o.result.Success {
			cache.Record(o.rec, o.dest)
		}
	}

	return summary
}
