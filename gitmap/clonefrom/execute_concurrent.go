package clonefrom

// execute_concurrent.go — bounded worker-pool variant of
// ExecuteWithHooks. Used by the cmd layer when --max-concurrency
// resolves to >1.
//
// Design contract (mirrors clonenow.ExecuteWithHooksConcurrent):
//
//   - On-disk layout is unchanged: every worker resolves its dest
//     via the row's Dest / DeriveDest verbatim, so increasing the
//     worker count NEVER reshuffles where repos land.
//   - Result ORDER matches input order. Per-row PROGRESS LINES are
//     emitted in input order AFTER the pool drains, keeping
//     `[i/total]` numbering monotonic for downstream scripts.
//   - The BeforeRow hook fires synchronously on the dispatcher
//     goroutine in input order, BEFORE the row enters the work
//     queue.
//   - workers <= 1 falls back to ExecuteWithHooks so there is
//     exactly one code path per regime.

import (
	"io"
	"os"
	"sync"
)

// concurrentJob is the named pool work unit (anonymous struct
// would be assignable but harder to refactor / test).
type concurrentJob struct {
	idx int
	row Row
}

// ExecuteWithHooksConcurrent is the parallel sibling of
// ExecuteWithHooks. See file header for the contract.
func ExecuteWithHooksConcurrent(plan Plan, cwd string, progress io.Writer,
	beforeRow BeforeRowHook, workers int) []Result {
	if workers <= 1 {
		return ExecuteWithHooks(plan, cwd, progress, beforeRow)
	}
	if len(cwd) == 0 {
		if wd, err := os.Getwd(); err == nil {
			cwd = wd
		}
	}
	out := make([]Result, len(plan.Rows))
	dispatchConcurrent(plan, cwd, beforeRow, workers, out)
	emitProgressInOrder(progress, out)

	return out
}

// dispatchConcurrent runs the worker pool and fills `out` at each
// row's input index.
func dispatchConcurrent(plan Plan, cwd string, beforeRow BeforeRowHook,
	workers int, out []Result) {
	jobs := make(chan concurrentJob, len(plan.Rows))
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go runConcurrentWorker(jobs, cwd, out, &wg)
	}
	enqueueConcurrentJobs(plan, beforeRow, jobs)
	close(jobs)
	wg.Wait()
}

// runConcurrentWorker is the per-goroutine drain loop.
func runConcurrentWorker(jobs <-chan concurrentJob, cwd string,
	out []Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		out[j.idx] = executeRow(j.row, cwd)
	}
}

// enqueueConcurrentJobs fires the BeforeRow hook (synchronously,
// in input order) and enqueues each row.
func enqueueConcurrentJobs(plan Plan, beforeRow BeforeRowHook,
	jobs chan<- concurrentJob) {
	total := len(plan.Rows)
	for i, r := range plan.Rows {
		if beforeRow != nil {
			dest := r.Dest
			if len(dest) == 0 {
				dest = DeriveDest(r.URL)
			}
			beforeRow(i+1, total, r, dest)
		}
		jobs <- concurrentJob{idx: i, row: r}
	}
}

// emitProgressInOrder prints progress lines in input order AFTER
// the pool drains. Keeps `[i/total]` lines monotonic for scripts.
func emitProgressInOrder(w io.Writer, out []Result) {
	if w == nil {
		return
	}
	total := len(out)
	for i, res := range out {
		writeProgress(w, i+1, total, res)
	}
}
