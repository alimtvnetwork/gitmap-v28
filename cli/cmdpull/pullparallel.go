package cmdpull

import (
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// runPullParallel pulls records concurrently using a worker pool.
func runPullParallel(records []model.ScanRecord, bar *PullProgressBar, parallel int) *apperror.AppError {
	limit := resolveParallelLimit(parallel, len(records))
	jobs := make(chan model.ScanRecord, len(records))
	var wg sync.WaitGroup

	startParallelWorkers(limit, jobs, bar, &wg)
	dispatchParallelJobs(records, jobs, bar)
	wg.Wait()

	return nil
}

func resolveParallelLimit(parallel, count int) int {
	if parallel < 1 {
		return 1
	}
	if parallel > count {
		return count
	}

	return parallel
}

func startParallelWorkers(limit int, jobs <-chan model.ScanRecord, bar *PullProgressBar, wg *sync.WaitGroup) {
	for workerID := 0; workerID < limit; workerID++ {
		wg.Add(1)
		go parallelPullWorker(workerID, jobs, bar, wg)
	}
}

func dispatchParallelJobs(records []model.ScanRecord, jobs chan<- model.ScanRecord, bar *PullProgressBar) {
	for _, rec := range records {
		if bar != nil && bar.IsStopped() {
			break
		}
		jobs <- rec
	}
	close(jobs)
}

func parallelPullWorker(workerID int, jobs <-chan model.ScanRecord, bar *PullProgressBar, wg *sync.WaitGroup) {
	defer wg.Done()
	for rec := range jobs {
		if bar != nil && bar.IsStopped() {
			return
		}
		runWorkerPull(workerID, rec, bar)
	}
}

func runWorkerPull(workerID int, rec model.ScanRecord, bar *PullProgressBar) {
	if bar != nil {
		bar.RegisterWorker(workerID, rec.RepoName)
	}
	ExecuteTrackedPullWithWorker(rec, bar, workerID)
	if bar != nil {
		bar.UnregisterWorker(workerID)
	}
}
