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
func ExecuteWithHooksConcurrent(params ConcurrentExecutionParams) []Result {
	if params.Workers <= 1 {
		return ExecuteWithHooks(params.Plan, params.Cwd, params.Progress, params.BeforeRow)
	}

	params.Cwd = resolveCwd(params.Cwd)
	out := make([]Result, len(params.Plan.Rows))
	dispatchConcurrent(CloneFromDispatchParams{
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
// row's input index.
func dispatchConcurrent(params CloneFromDispatchParams) {
	jobs := make(chan concurrentJob, len(params.Plan.Rows))
	var wg sync.WaitGroup
	wg.Add(params.Workers)
	for i := 0; i < params.Workers; i++ {
		go runConcurrentWorker(CloneFromWorkerParams{
			Jobs: jobs,
			Cwd:  params.Cwd,
			Out:  params.Out,
			Wg:   &wg,
		})
	}

	enqueueConcurrentJobs(params.Plan, params.BeforeRow, jobs)
	close(jobs)
	wg.Wait()
}

// runConcurrentWorker is the per-goroutine drain loop.
func runConcurrentWorker(params CloneFromWorkerParams) {
	defer params.Wg.Done()
	for j := range params.Jobs {
		params.Out[j.idx] = executeRow(j.row, params.Cwd)
	}
}

// enqueueConcurrentJobs fires the BeforeRow hook (synchronously,
// in input order) and enqueues each row.
func enqueueConcurrentJobs(plan Plan, beforeRow BeforeRowHook,
	jobs chan<- concurrentJob) {
	total := len(plan.Rows)
	for i, r := range plan.Rows {
		if beforeRow != nil {
			invokeBeforeRow(BeforeRowInvokeParams{
				Hook:         beforeRow,
				CurrentIndex: i + 1,
				TotalCount:   total,
				Row:          r,
			})
		}

		jobs <- concurrentJob{idx: i, row: r}
	}
}

func invokeBeforeRow(params BeforeRowInvokeParams) {
	dest := params.Row.Dest
	if len(dest) == 0 {
		dest = DeriveDest(params.Row.URL)
	}

	params.Hook(params.CurrentIndex, params.TotalCount, params.Row, dest)
}

// emitProgressInOrder prints progress lines in input order AFTER
// the pool drains. Keeps `[i/total]` lines monotonic for scripts.
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
