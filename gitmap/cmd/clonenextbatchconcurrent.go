package cmd

// Bounded worker-pool path for `gitmap cn --all` / `--csv` batch mode.
//
// Split out of clonenextbatch.go to keep that file focused on the
// dispatcher + sequential path and to honor the project's 200-line
// per-file budget. The pattern mirrors gitmap/cloner/concurrent.go:
// workers drain a buffered job channel, write outcomes to a buffered
// result channel, and the collector preserves input ordering so the
// CSV report rows are deterministic regardless of worker count.

// normalizeBatchWorkers clamps the requested worker count to the work
// queue size. Mirrors cloner.normalizeWorkers; duplicated rather than
// imported to keep the cmd package free of an internal cloner dep
// (this helper is 5 lines — the cost is trivial).
func normalizeBatchWorkers(requested, jobs int) int {
	if requested < 1 {
		return 1
	}
	if requested > jobs && jobs > 0 {
		return jobs
	}

	return requested
}

// indexedBatchJob carries the input-list position alongside the repo
// path so the collector can restore input order.
type indexedBatchJob struct {
	idx  int
	path string
}

// indexedBatchResult pairs a per-repo outcome with its original index.
type indexedBatchResult struct {
	idx int
	row batchRowResult
}

// processBatchReposConcurrent fans per-repo work across a bounded pool.
// Each worker pulls one repo from the job channel and writes its
// outcome to the result channel; the collector re-orders by input
// index so the CSV report rows still match the caller's repo list.
//
// `onResult` fires once per dequeued result (i.e. as workers finish,
// not at the end) so callers can print real-time progress without
// reaching into the pool internals. Safe to pass nil.
//
// Concurrency contract: processOneBatchRepo owns its own per-repo
// state (probe, version read, row construction) and shares no mutable
// data with peers, so no extra mutex is required here.
func processBatchReposConcurrent(repos []string, workers int, onResult func(batchRowResult)) []batchRowResult {
	jobs := make(chan indexedBatchJob, len(repos))
	results := make(chan indexedBatchResult, len(repos))

	startBatchWorkers(workers, jobs, results)
	enqueueBatchJobs(repos, jobs)

	return collectBatchResults(repos, results, onResult)
}

// processOneBatchRepoFn is the worker entrypoint, exposed as a
// package-level var so E2E tests can swap in a deterministic stub
// without spawning real git processes. Default points at the real
// implementation in clonenextbatch.go; production code never
// reassigns it. Tests restore the original in t.Cleanup.
var processOneBatchRepoFn = processOneBatchRepo

// startBatchWorkers spins up the worker goroutines. Each one drains
// the job channel until it closes.
func startBatchWorkers(workers int, jobs <-chan indexedBatchJob, results chan<- indexedBatchResult) {
	for w := 0; w < workers; w++ {
		go func() {
			for j := range jobs {
				results <- indexedBatchResult{idx: j.idx, row: processOneBatchRepoFn(j.path)}
			}
		}()
	}
}

// enqueueBatchJobs pushes every repo onto the job channel and closes
// it so workers exit cleanly once drained.
func enqueueBatchJobs(repos []string, jobs chan<- indexedBatchJob) {
	for i, r := range repos {
		jobs <- indexedBatchJob{idx: i, path: r}
	}
	close(jobs)
}

// collectBatchResults drains the result channel and slots each row
// back into its input position. `onResult` fires synchronously for
// each dequeued row so the caller sees live "X done" feedback as
// workers complete. The collector remains the single goroutine
// touching the output slice — no locking required.
func collectBatchResults(repos []string, results <-chan indexedBatchResult, onResult func(batchRowResult)) []batchRowResult {
	out := make([]batchRowResult, len(repos))
	for i := 0; i < len(repos); i++ {
		r := <-results
		out[r.idx] = r.row
		if onResult != nil {
			onResult(r.row)
		}
	}

	return out
}
