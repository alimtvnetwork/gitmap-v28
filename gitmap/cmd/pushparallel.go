package cmd

import (
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runPushParallel pushes every record concurrently using a worker pool of
// the given width. BatchProgress is not goroutine-safe by itself, so all
// progress mutations happen under progMu.
//
// stopOnFail is honored: once any worker reports a failure, the dispatcher
// drains the queue without spawning more work and returning workers exit
// after their in-flight task finishes.
func runPushParallel(records []model.ScanRecord, prog *cloner.BatchProgress, parallel int, stopOnFail bool) error {
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

	startPushWorkers(parallel, jobs, prog, &progMu, &wg, stopOnFail, &stopped)
	dispatchPushJobs(records, jobs, &progMu, &stopped)
	wg.Wait()
	return nil
}

// startPushWorkers spins up `count` workers, each draining the jobs channel
// until it closes.
func startPushWorkers(count int, jobs <-chan model.ScanRecord, prog *cloner.BatchProgress,
	progMu *sync.Mutex, wg *sync.WaitGroup, stopOnFail bool, stopped *bool) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go pushWorker(jobs, prog, progMu, wg, stopOnFail, stopped)
	}
}

// dispatchPushJobs feeds records into the jobs channel respecting stopOnFail.
// Closes the channel when done so workers exit.
func dispatchPushJobs(records []model.ScanRecord, jobs chan<- model.ScanRecord,
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

// pushWorker drains the channel and runs SafePushOne on each record.
// All BatchProgress mutations are guarded by progMu.
func pushWorker(jobs <-chan model.ScanRecord, prog *cloner.BatchProgress,
	progMu *sync.Mutex, wg *sync.WaitGroup, stopOnFail bool, stopped *bool) {
	defer wg.Done()

	for rec := range jobs {
		runOnePushJob(rec, prog, progMu, stopOnFail, stopped)
	}
}

// runOnePushJob handles a single record under the progress mutex. Sets
// *stopped when a failure occurs and stopOnFail is enabled.
func runOnePushJob(rec model.ScanRecord, prog *cloner.BatchProgress,
	progMu *sync.Mutex, stopOnFail bool, stopped *bool) error {
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		progMu.Lock()
		prog.BeginItem(rec.RepoName)
		prog.Skip(rec.RepoName)
		progMu.Unlock()
		return nil
	}

	progMu.Lock()
	prog.BeginItem(rec.RepoName)
	progMu.Unlock()

	result := cloner.SafePushOne(rec, rec.AbsolutePath)

	progMu.Lock()
	if result.IsSuccess == false {
		prog.FailWithError(rec.RepoName, result.Error)
	}
	if result.IsSuccess == false && stopOnFail == true {
		*stopped = true
	}
	if result.IsSuccess == false {
		progMu.Unlock()
		return nil
	}
	if result.Notes == "up-to-date" {
		prog.UpToDate(rec.RepoName)
	} else {
		prog.Succeed(rec.RepoName)
	}
	progMu.Unlock()
	return nil
}
