package cmd

import (
	"sync"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runPullParallel pulls every record concurrently using a worker pool of
// the given width. BatchProgress is not goroutine-safe by itself, so all
// progress mutations happen under progMu.
//
// stopOnFail is honored: once any worker reports a failure, the dispatcher
// drains the queue without spawning more work and returning workers exit
// after their in-flight task finishes.
func runPullParallel(records []model.ScanRecord, prog *cloner.BatchProgress, parallel int, stopOnFail bool) {
	if parallel < 1 {
		parallel = 1
	}
	if parallel > len(records) {
		parallel = len(records)
	}

	jobs := make(chan model.ScanRecord, len(records))
	var (
		wg      sync.WaitGroup
		progMu  sync.Mutex
		stopped bool
	)

	startPullWorkers(parallel, jobs, prog, &progMu, &wg, stopOnFail, &stopped)
	dispatchPullJobs(records, jobs, &progMu, &stopped)
	wg.Wait()
}

// startPullWorkers spins up `count` workers, each draining the jobs channel
// until it closes.
func startPullWorkers(count int, jobs <-chan model.ScanRecord, prog *cloner.BatchProgress,
	progMu *sync.Mutex, wg *sync.WaitGroup, stopOnFail bool, stopped *bool) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go pullWorker(jobs, prog, progMu, wg, stopOnFail, stopped)
	}
}

// dispatchPullJobs feeds records into the jobs channel respecting stopOnFail.
// Closes the channel when done so workers exit.
func dispatchPullJobs(records []model.ScanRecord, jobs chan<- model.ScanRecord,
	progMu *sync.Mutex, stopped *bool) {
	for _, rec := range records {
		progMu.Lock()
		halted := *stopped
		progMu.Unlock()
		if halted {
			break
		}
		jobs <- rec
	}
	close(jobs)
}

// pullWorker drains the channel and runs SafePullOne on each record.
// All BatchProgress mutations are guarded by progMu.
func pullWorker(jobs <-chan model.ScanRecord, prog *cloner.BatchProgress,
	progMu *sync.Mutex, wg *sync.WaitGroup, stopOnFail bool, stopped *bool) {
	defer wg.Done()

	for rec := range jobs {
		runOnePullJob(rec, prog, progMu, stopOnFail, stopped)
	}
}

// runOnePullJob handles a single record under the progress mutex. Sets
// *stopped when a failure occurs and stopOnFail is enabled.
func runOnePullJob(rec model.ScanRecord, prog *cloner.BatchProgress,
	progMu *sync.Mutex, stopOnFail bool, stopped *bool) {
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		progMu.Lock()
		prog.BeginItem(rec.RepoName)
		prog.Skip()
		progMu.Unlock()
		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)

	progMu.Lock()
	prog.BeginItem(rec.RepoName)
	if result.Success {
		prog.Succeed()
	} else {
		prog.FailWithError(rec.RepoName, result.Error)
		if stopOnFail {
			*stopped = true
		}
	}
	progMu.Unlock()
}
