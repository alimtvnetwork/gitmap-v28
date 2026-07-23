package cmd

// E2E-style ordering + count tests for `gitmap cn --all` / `--csv`
// batch mode under concurrency (v3.126.0+). These tests stub
// processOneBatchRepoFn (see clonenextbatchconcurrent_e2e_helpers_test.go)
// so they run in milliseconds without real git repos.
//
// CSV byte-equivalence + full-write-report tests live in the
// sibling clonenextbatchconcurrent_e2e_csv_test.go file.

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestE2E_BatchConcurrency_DeterministicOrdering proves that under
// --max-concurrency=8 over 50 repos, the result slice is still in
// input order even though the stub deliberately makes workers
// finish out-of-order.
func TestE2E_BatchConcurrency_DeterministicOrdering(t *testing.T) {
	peak := installStubProcessor(t)
	repos := makeRepoPaths(50)

	results := processBatchReposConcurrent(repos, 8, nil)

	if len(results) != len(repos) {
		t.Fatalf("results length %d != input %d", len(results), len(repos))
	}
	for i, r := range results {
		if r.RepoPath != repos[i] {
			t.Fatalf("ordering drift at index %d: got %s, want %s",
				i, r.RepoPath, repos[i])
		}
	}
	if atomic.LoadInt64(peak) < 2 {
		t.Fatalf("pool never went parallel (peak inflight = %d) — test is invalid",
			atomic.LoadInt64(peak))
	}
}

// TestE2E_BatchConcurrency_StatusCountsExact verifies the per-bucket
// totals exactly match what the stub emitted: 50 repos with last
// digits 0-9 cycling 5 times → 20 ok, 20 failed, 10 skipped.
func TestE2E_BatchConcurrency_StatusCountsExact(t *testing.T) {
	installStubProcessor(t)
	repos := makeRepoPaths(50)

	results := processBatchReposConcurrent(repos, 4, nil)
	ok, failed, skipped := tallyBatch(results)

	const wantOK, wantFailed, wantSkipped = 20, 20, 10
	if ok != wantOK || failed != wantFailed || skipped != wantSkipped {
		t.Fatalf("counts: ok=%d failed=%d skipped=%d, want %d/%d/%d",
			ok, failed, skipped, wantOK, wantFailed, wantSkipped)
	}
}

// TestE2E_BatchConcurrency_ProgressCallbackFires verifies the
// onResult callback handed to the collector receives exactly one
// invocation per repo, with the same total count as the returned
// slice — proving the live progress reporter sees a faithful view
// even under heavy concurrency.
func TestE2E_BatchConcurrency_ProgressCallbackFires(t *testing.T) {
	installStubProcessor(t)
	repos := makeRepoPaths(30)

	var seenCount int64
	cb := func(_ batchRowResult) {
		atomic.AddInt64(&seenCount, 1)
	}

	results := processBatchReposConcurrent(repos, 6, cb)

	if got := atomic.LoadInt64(&seenCount); int(got) != len(repos) {
		t.Fatalf("callback fired %d times, want %d", got, len(repos))
	}
	if len(results) != len(repos) {
		t.Fatalf("results length: got %d, want %d", len(results), len(repos))
	}
}

// TestE2E_BatchConcurrency_CollectorReordersByInputIndex makes the
// reorder-by-input-index contract explicit. The previous ordering
// test relied on incidental sleep variance; this one *forces*
// completion order to be the exact reverse of input order via
// per-index sleeps, then asserts:
//
//  1. Completion order recorded by the worker stub is non-monotonic
//     (proves the randomization actually happened — guards against a
//     future change that accidentally serializes workers and makes
//     the input-order check trivially pass).
//  2. results[i].RepoPath == repos[i] for every i (proves the
//     collector slots each result back into its input position
//     regardless of which worker finished first).
//
// This is the canonical regression guard for collectBatchResults.
func TestE2E_BatchConcurrency_CollectorReordersByInputIndex(t *testing.T) {
	const n = 20
	repos := makeRepoPaths(n)

	var (
		mu              sync.Mutex
		completionOrder []string
	)
	original := processOneBatchRepoFn
	processOneBatchRepoFn = func(path string) batchRowResult {
		// Sleep proportional to (n - trailing-index) so repos at the
		// end of the input list finish FIRST. With workers >= n every
		// job starts immediately and completion order is deterministic
		// reverse of input order.
		idx := indexFromRepoPath(path)
		time.Sleep(time.Duration(n-idx) * time.Millisecond)
		mu.Lock()
		completionOrder = append(completionOrder, path)
		mu.Unlock()
		return batchRowResult{RepoPath: path, FromVersion: "v1", ToVersion: "v2"}
	}
	t.Cleanup(func() { processOneBatchRepoFn = original })

	// workers == n guarantees every job starts immediately, so the
	// per-index sleep dictates completion order with no scheduling
	// noise from queue depth.
	results := processBatchReposConcurrent(repos, n, nil)

	assertCompletionOrderRandomized(t, completionOrder, repos)
	assertResultsMatchInputOrder(t, results, repos)
}

// indexFromRepoPath parses the trailing integer from the synthetic
// "/tmp/repo-N" paths produced by makeRepoPaths.
func indexFromRepoPath(path string) int {
	var idx int
	if _, err := fmt.Sscanf(path, "/tmp/repo-%d", &idx); err != nil {
		return 0
	}
	return idx
}

// assertCompletionOrderRandomized fails the test when the recorded
// completion order matches the input order — that would mean workers
// finished in input order and the collector's reorder logic is never
// actually exercised, making the input-order assertion meaningless.
func assertCompletionOrderRandomized(t *testing.T, completionOrder, repos []string) {
	t.Helper()
	if len(completionOrder) != len(repos) {
		t.Fatalf("completion order length %d != input %d",
			len(completionOrder), len(repos))
	}
	matchesInputOrder := true
	for i, p := range completionOrder {
		if p != repos[i] {
			matchesInputOrder = false
			break
		}
	}
	if matchesInputOrder {
		t.Fatalf("completion order matched input order — randomization stub failed, " +
			"reorder-by-index assertion would be trivial")
	}
}

// assertResultsMatchInputOrder is the core contract check: even
// though workers finished in scrambled order, results[i] must be
// the row produced for repos[i].
func assertResultsMatchInputOrder(t *testing.T, results []batchRowResult, repos []string) {
	t.Helper()
	if len(results) != len(repos) {
		t.Fatalf("results length %d != input %d", len(results), len(repos))
	}
	for i, r := range results {
		if r.RepoPath != repos[i] {
			t.Fatalf("collector failed to reorder: results[%d].RepoPath = %q, want %q",
				i, r.RepoPath, repos[i])
		}
	}
}
