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

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
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

// ConcurrentRunParams encapsulates all arguments required for concurrent clone execution.
type ConcurrentRunParams struct {
	Records   []model.ScanRecord
	TargetDir string
	Options   CloneOptions
	Workers   int
	Progress  *Progress
	Cache     *CloneCache
}

// WorkerParams encapsulates dependencies for concurrent worker goroutines.
type WorkerParams struct {
	Workers   int
	Jobs      <-chan cloneJob
	Out       chan<- cloneOutcome
	TargetDir string
	Options   CloneOptions
	Progress  *Progress
}

// EnqueueJobsParams encapsulates parameters for enqueuing clone jobs.
type EnqueueJobsParams struct {
	Records   []model.ScanRecord
	TargetDir string
	Cache     *CloneCache
	Jobs      chan<- cloneJob
	Out       chan<- cloneOutcome
}

// CollectOutcomesParams encapsulates parameters for collecting worker outcomes.
type CollectOutcomesParams struct {
	Records    []model.ScanRecord
	TargetDir  string
	IsSafePull bool
	Progress   *Progress
	Cache      *CloneCache
	Out        <-chan cloneOutcome
}

// runConcurrent fans the records out across `workers` goroutines and
// returns the same CloneSummary shape as the sequential runner. The
// caller is responsible for picking a sane worker count (>=1).
func runConcurrent(params ConcurrentRunParams) model.CloneSummary {
	jobs := make(chan cloneJob, len(params.Records))
	out := make(chan cloneOutcome, len(params.Records))

	startWorkers(WorkerParams{
		Workers:   params.Workers,
		Jobs:      jobs,
		Out:       out,
		TargetDir: params.TargetDir,
		Options:   params.Options,
		Progress:  params.Progress,
	})
	enqueueJobs(EnqueueJobsParams{
		Records:   params.Records,
		TargetDir: params.TargetDir,
		Cache:     params.Cache,
		Jobs:      jobs,
		Out:       out,
	})
	close(jobs)

	return collectOutcomes(CollectOutcomesParams{
		Records:    params.Records,
		TargetDir:  params.TargetDir,
		IsSafePull: params.Options.IsSafePull,
		Progress:   params.Progress,
		Cache:      params.Cache,
		Out:        out,
	})
}

// startWorkers spins up the worker goroutines.
func startWorkers(params WorkerParams) {
	for i := 0; i < params.Workers; i++ {
		go cloneWorker(params)
	}
}

// cloneWorker drains the job channel until it closes.
func cloneWorker(params WorkerParams) {
	for job := range params.Jobs {
		params.Progress.Begin(repoDisplayName(job.rec))
		result := cloneOrPullOne(job.rec, params.TargetDir, params.Options)
		params.Out <- cloneOutcome{rec: job.rec, dest: job.dest, result: result}
	}
}

// enqueueJobs short-circuits cache hits (reported synchronously so the
// progress line lands before any worker output) and dispatches the rest
// onto the job channel.
func enqueueJobs(params EnqueueJobsParams) {
	for _, rec := range params.Records {
		dest := filepath.Join(params.TargetDir, model.CleanRelativePath(rec.RelativePath))
		if params.Cache.IsUpToDate(rec, dest) {
			params.Out <- cloneOutcome{
				rec:    rec,
				dest:   dest,
				result: model.CloneResult{Record: rec, IsSuccess: true},
				cached: true,
			}

			continue
		}

		params.Jobs <- cloneJob{rec: rec, dest: dest}
	}
}

// collectOutcomes drains the outcome channel, updating progress, cache,
// and summary in the same order outcomes complete.
func collectOutcomes(params CollectOutcomesParams) model.CloneSummary {
	summary := model.CloneSummary{}
	for i := 0; i < len(params.Records); i++ {
		o := <-params.Out
		if o.cached {
			params.Progress.Skip(o.result)
			summary = updateSummarySkipped(summary, o.result)

			continue
		}

		_ = TrackResult(TrackResultParams{
			Progress:   params.Progress,
			Result:     o.result,
			ScanRecord: o.rec,
			TargetDir:  params.TargetDir,
			IsSafePull: params.IsSafePull,
		})
		summary = updateSummary(summary, o.result)
		if o.result.IsSuccess {
			params.Cache.Record(o.rec, o.dest)
		}
	}

	return summary
}
