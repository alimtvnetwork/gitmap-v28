// Package probe — background.go
//
// BackgroundRunner is a fire-and-forget version-probe pool. The scan
// command kicks it off as soon as repo records are persisted, then
// blocks on Wait() right before exiting so the DB is up-to-date by
// the time the user runs `gitmap find-next`.
//
// Design constraints (locked in via UX questions):
//
//   - Worker count defaults to 3 and is overridable with
//     `--probe-concurrency N`. Anything <1 disables the runner.
//   - The runner exposes a one-line API: Start(record), Wait(). All
//     pool / mutex / channel mechanics are encapsulated; callers do
//     not import sync, see goroutines, or learn channel directions.
//   - Wait is idempotent — calling it twice is a no-op so the scan
//     command can safely Wait in a defer AND in the explicit drain
//     line, without risking a double close.
//   - Persistence happens inside the worker via a user-supplied
//     Sink so background.go stays free of store / model imports
//     beyond the probe.Result type. This keeps the package's
//     dependency graph identical to the pre-feature shape.
//
// Concurrency invariants are documented at each method.
package probe

import (
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// Sink is the persistence callback the runner invokes for each
// completed probe. The scan command supplies a closure that wraps
// (db, repo, result) → store.RecordVersionProbe. Returning a sink
// rather than calling the store directly keeps probe/ free of any
// DB dependency and makes the runner trivially testable with an
// in-memory recorder (see background_test.go).
type Sink func(record model.ScanRecord, result Result)

// BackgroundRunner is a bounded worker pool that probes repos
// asynchronously. Use New, then Start each repo, then Wait once
// before process exit.
type BackgroundRunner struct {
	jobs    chan model.ScanRecord
	wg      sync.WaitGroup
	sink    Sink
	urlPick func(model.ScanRecord) string
	closed  bool
	closeMu sync.Mutex
	stats   runnerStats

	// onFailure, when non-nil, is invoked from the worker goroutine
	// for every probe whose Result.Error is non-empty. Receives the
	// originating ScanRecord (so the caller can read RelativePath /
	// CloneURL / etc.) plus the probe result. Implementations MUST be
	// goroutine-safe — multiple workers may call it concurrently.
	// Set via SetFailureHook BEFORE the first Start so the workers
	// observe the assignment without a data race.
	onFailure func(record model.ScanRecord, result Result)

	// cloneDepth is the `--depth N` passed to the shallow-clone
	// fallback. Defaults to constants.ProbeDefaultDepth (1). Override
	// via SetCloneDepth before the first Start so workers observe a
	// stable value without locking.
	cloneDepth int
}

// SetFailureHook installs a per-failure callback. Must be called
// BEFORE the first Start — the runner does not synchronize the read
// in workerLoop, since the typical wiring (scan command) sets the
// hook immediately after construction and the workers don't begin
// pulling jobs until Start enqueues the first record. A nil receiver
// is a silent no-op so call sites can chain unconditionally.
func (r *BackgroundRunner) SetFailureHook(hook func(model.ScanRecord, Result)) {
	if r == nil {
		return
	}
	r.onFailure = hook
}

// SetCloneDepth overrides the shallow-clone --depth value used by the
// fallback strategy. Must be called BEFORE the first Start; values < 1
// are coerced to 1. Nil receiver is a silent no-op.
func (r *BackgroundRunner) SetCloneDepth(depth int) {
	if r == nil {
		return
	}
	if depth < 1 {
		depth = 1
	}
	r.cloneDepth = depth
}

// runnerStats tracks per-bucket counters under its own mutex.
type runnerStats struct {
	mu                           sync.Mutex
	queued                       int
	available, unchanged, failed int
}

// NewBackgroundRunner spins up `workers` goroutines feeding from a
// buffered job channel sized for `expectedJobs`. `urlPick` selects
// HTTPS-or-SSH per record (the scan caller wires this to the same
// pickProbeURL helper used by `gitmap probe`). `sink` persists each
// outcome.
//
// A workers value <1 returns nil so callers can treat
// "concurrency disabled" the same as "no runner needed" with a
// single nil-check.
func NewBackgroundRunner(workers, expectedJobs int, urlPick func(model.ScanRecord) string, sink Sink) *BackgroundRunner {
	if workers < 1 {
		return nil
	}
	r := &BackgroundRunner{
		jobs:       make(chan model.ScanRecord, expectedJobs),
		sink:       sink,
		urlPick:    urlPick,
		cloneDepth: constants.ProbeDefaultDepth,
	}
	r.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go r.workerLoop()
	}

	return r
}

// Start enqueues one repo for probing. Safe to call from many
// goroutines, but in practice scan calls it sequentially right
// after DB upsert. Calling Start after Wait is a silent no-op
// (the channel is closed) — this prevents the classic "send on
// closed channel" panic if a caller mis-orders the calls.
func (r *BackgroundRunner) Start(record model.ScanRecord) {
	if r == nil {
		return
	}
	r.closeMu.Lock()
	closed := r.closed
	r.closeMu.Unlock()
	if closed {
		return
	}
	r.stats.mu.Lock()
	r.stats.queued++
	r.stats.mu.Unlock()
	r.jobs <- record
}

// Wait closes the job channel (first call only) and blocks until
// every queued probe has finished and persisted. Idempotent: a
// second call returns immediately with the same final stats.
func (r *BackgroundRunner) Wait() Stats {
	if r == nil {
		return Stats{}
	}
	r.closeMu.Lock()
	if !r.closed {
		close(r.jobs)
		r.closed = true
	}
	r.closeMu.Unlock()
	r.wg.Wait()

	return r.Stats()
}

// Stats is a non-blocking snapshot of the runner's counters.
// Useful for the "(N remaining)" line printed before Wait blocks.
type Stats struct {
	Queued, Available, Unchanged, Failed int
}

// Stats returns a copy of the current counters.
func (r *BackgroundRunner) Stats() Stats {
	if r == nil {
		return Stats{}
	}
	r.stats.mu.Lock()
	defer r.stats.mu.Unlock()

	return Stats{
		Queued:    r.stats.queued,
		Available: r.stats.available,
		Unchanged: r.stats.unchanged,
		Failed:    r.stats.failed,
	}
}

// Remaining returns Queued minus completed-so-far, used for the
// "waiting for N probes" line. Reads through Stats so the math
// stays consistent under the same lock.
func (r *BackgroundRunner) Remaining() int {
	s := r.Stats()
	done := s.Available + s.Unchanged + s.Failed

	return s.Queued - done
}

// workerLoop drains jobs until the channel closes. Each job runs
// the probe, hands the result to the sink, and updates counters.
func (r *BackgroundRunner) workerLoop() {
	defer r.wg.Done()
	for record := range r.jobs {
		result := r.probeOne(record)
		if r.sink != nil {
			r.sink(record, result)
		}
		if len(result.Error) > 0 && r.onFailure != nil {
			r.onFailure(record, result)
		}
		r.tally(result)
	}
}

// probeOne runs the actual probe and tags missing-URL records with
// the standard error so the sink stores a row regardless. Mirrors
// the foreground `cmd/probe.go` behavior so background and
// foreground rows look identical in the DB.
func (r *BackgroundRunner) probeOne(record model.ScanRecord) Result {
	url := ""
	if r.urlPick != nil {
		url = r.urlPick(record)
	}
	if url == "" {
		return Result{Method: constants.ProbeMethodNone, Error: "empty clone url"}
	}

	return RunOneWithDepth(url, r.cloneDepth)
}

// tally updates per-bucket counters under the stats mutex. Split
// out to keep workerLoop short.
func (r *BackgroundRunner) tally(result Result) {
	r.stats.mu.Lock()
	defer r.stats.mu.Unlock()
	if len(result.Error) > 0 {
		r.stats.failed++

		return
	}
	if result.IsAvailable {
		r.stats.available++

		return
	}
	r.stats.unchanged++
}
