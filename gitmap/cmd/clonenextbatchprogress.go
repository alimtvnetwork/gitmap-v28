package cmd

// Real-time per-repo progress reporting for `gitmap cn --all` /
// `--csv` (v3.124.0+).
//
// The collector in clonenextbatchconcurrent.go fires onResult once per
// finished repo. This file owns the counter + printing logic so the
// concurrent file stays free of fmt imports and the dispatcher in
// clonenextbatch.go reads as a one-liner.
//
// Concurrency contract: the collector is a single goroutine calling
// OnResult sequentially as it dequeues results — no locking is needed
// here even though the underlying workers run in parallel. Documented
// at OnResult so future maintainers don't add a redundant mutex.

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// batchProgressReporter tracks completed/failed/skipped counters and
// prints a one-line update each time a repo finishes.
type batchProgressReporter struct {
	total               int
	done                int
	ok, failed, skipped int
	silent              bool
}

// newBatchProgressReporter returns a reporter sized for `total` jobs.
// When `silent` is true, OnResult still bumps counters but suppresses
// the per-repo print (the final summary is owned by printBatchSummary
// and prints regardless).
func newBatchProgressReporter(total int, silent bool) *batchProgressReporter {
	return &batchProgressReporter{total: total, silent: silent}
}

// OnResult is the callback handed to the batch collector. Called
// exactly once per finished repo, in completion order. Safe to call
// from a single collector goroutine; not safe to share across
// callers without external locking (intentional — the collector
// already serializes results so a mutex here would just add cost).
func (r *batchProgressReporter) OnResult(row batchRowResult) {
	r.done++
	r.tally(row.Status)
	if r.silent {
		return
	}
	fmt.Printf(constants.MsgCloneNextBatchProgressFmt,
		r.done, r.total,
		filepath.Base(row.RepoPath),
		row.Status,
		r.ok, r.failed, r.skipped)
}

// tally increments the bucket matching this row's status. Unknown
// statuses are intentionally ignored — the final summary's
// tallyBatch helper is the source of truth for the report rows;
// this counter is a live progress indicator only.
func (r *batchProgressReporter) tally(status string) {
	switch status {
	case constants.BatchStatusOK:
		r.ok++
	case constants.BatchStatusFailed:
		r.failed++
	case constants.BatchStatusSkipped:
		r.skipped++
	}
}
