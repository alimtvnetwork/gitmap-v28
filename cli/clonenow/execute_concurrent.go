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
func ExecuteWithHooksConcurrent(params ConcurrentExecutionParams) []Result {
	if params.Workers <= 1 {
		return ExecuteWithHooks(params.Plan, params.Cwd, params.Progress, params.BeforeRow)
	}

	wd, err := os.Getwd()
	if len(params.Cwd) == 0 && err == nil {
		params.Cwd = wd
	}

	out := make([]Result, len(params.Plan.Rows))
	dispatchConcurrent(ConcurrentDispatchParams{
		Plan:      params.Plan,
		Cwd:       params.Cwd,
		BeforeRow: params.BeforeRow,
		Workers:   params.Workers,
		Out:       out,
	})
	emitProgressInOrder(params.Progress, out)

	return out
}

// dispatchConcurrent runs the worker pool and fills `out` at each
// row's input index. Split out so ExecuteWithHooksConcurrent stays
// under the 15-line function cap.
func dispatchConcurrent(params ConcurrentDispatchParams) {
	jobs := make(chan concurrentJob, len(params.Plan.Rows))
	var wg sync.WaitGroup
	wg.Add(params.Workers)
	for i := 0; i < params.Workers; i++ {
		go runConcurrentWorker(ConcurrentWorkerParams{
			Jobs: jobs,
			Plan: params.Plan,
			Cwd:  params.Cwd,
			Out:  params.Out,
			Wg:   &wg,
		})
	}

	enqueueConcurrentJobs(params.Plan, params.BeforeRow, jobs)
	close(jobs)
	wg.Wait()
}

// runConcurrentWorker is the per-goroutine drain loop. Pulled out
// so dispatchConcurrent stays under the function-length cap.
func runConcurrentWorker(params ConcurrentWorkerParams) {
	defer params.Wg.Done()
	for j := range params.Jobs {
		params.Out[j.idx] = executeRow(j.row, params.Plan, params.Cwd)
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
		writeProgress(ProgressWriteParams{
			Writer:       w,
			CurrentIndex: i + 1,
			TotalCount:   total,
			Result:       res,
		})
	}
}
