package clonenow

// execute_concurrent.go — bounded worker-pool variant of
// ExecuteWithHooks. Used by the cmd layer when --max-concurrency
// resolves to >1.
//
// Design contract:
//
//   - The on-disk layout is unchanged: every worker resolves its
//     destination via the row's RelativePath verbatim (same as the
//     sequential path), so increasing the worker count NEVER
//     reshuffles where repos land.
//   - Result ORDER matches input order. Scripts that grep stderr
//     for "[i/total]" expect monotonic numbering, so per-row
//     progress lines are emitted in order AFTER the pool drains.
//   - The BeforeRow hook fires synchronously on the dispatcher
//     goroutine in input order, BEFORE the row enters the work
//     queue. Mirrors the sequential hook timing contract.
//   - workers <= 1 falls back to the sequential ExecuteWithHooks so
//     there is exactly one code path per regime.

import (
	"io"
	"os"
	"sync"
)

// concurrentJob is the named pool work unit (an anonymous struct
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
// row's input index. Split out so ExecuteWithHooksConcurrent stays
// under the 15-line function cap.
func dispatchConcurrent(plan Plan, cwd string, beforeRow BeforeRowHook,
	workers int, out []Result) {
	jobs := make(chan concurrentJob, len(plan.Rows))
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go runConcurrentWorker(jobs, plan, cwd, out, &wg)
	}
	enqueueConcurrentJobs(plan, beforeRow, jobs)
	close(jobs)
	wg.Wait()
}

// runConcurrentWorker is the per-goroutine drain loop. Pulled out
// so dispatchConcurrent stays under the function-length cap.
func runConcurrentWorker(jobs <-chan concurrentJob, plan Plan, cwd string,
	out []Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		out[j.idx] = executeRow(j.row, plan, cwd)
	}
}

// enqueueConcurrentJobs fires the BeforeRow hook (synchronously, in
// input order) and enqueues each row. The channel's buffer is
// sized to the row count so this never blocks the dispatcher.
func enqueueConcurrentJobs(plan Plan, beforeRow BeforeRowHook,
	jobs chan<- concurrentJob) {
	total := len(plan.Rows)
	for i, r := range plan.Rows {
		if beforeRow != nil {
			url := r.PickURL(plan.Mode)
			beforeRow(i+1, total, r, url, r.RelativePath)
		}
		jobs <- concurrentJob{idx: i, row: r}
	}
}

// emitProgressInOrder prints progress lines in input order AFTER
// the pool drains. Trade-off: progress is post-hoc rather than
// real-time, but ordering matches the sequential runner's contract
// — keeping `[i/total]` lines monotonic for scripts.
func emitProgressInOrder(w io.Writer, out []Result) {
	if w == nil {
		return
	}
	total := len(out)
	for i, res := range out {
		writeProgress(w, i+1, total, res)
	}
}
