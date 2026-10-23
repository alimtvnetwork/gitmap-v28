package cmdpull

import (
	"sync"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// PullWorkerPool coordinates dynamic CPU-aware parallel repository pull workers.
type PullWorkerPool struct {
	workerCount      int
	isSSHSession     bool
	hasActiveWorkers bool
	jobQueue         chan model.ScanRecord
	bar              *PullProgressBar
	waitGroup        sync.WaitGroup
}

// NewPullWorkerPool constructs a parallel worker pool with dynamic CPU-aware sizing.
func NewPullWorkerPool(requestedWorkers, repoCount int, isSSHSession bool, bar *PullProgressBar) *PullWorkerPool {
	workers := ResolveDynamicPullWorkers(requestedWorkers, isSSHSession)
	effectiveWorkers := ClampWorkers(workers, repoCount)
	queueCap := resolveQueueCapacity(effectiveWorkers, repoCount)

	return &PullWorkerPool{
		workerCount:  effectiveWorkers,
		isSSHSession: isSSHSession,
		jobQueue:     make(chan model.ScanRecord, queueCap),
		bar:          bar,
	}
}

func resolveQueueCapacity(workers, repoCount int) int {
	if repoCount > 0 && repoCount < workers {
		return repoCount
	}

	if workers > 0 {
		return workers
	}

	return 1
}

// WorkerCount returns the active number of dynamic workers allocated.
func (p *PullWorkerPool) WorkerCount() int {
	return p.workerCount
}

// HasActiveWorkers returns true if workers are currently running.
func (p *PullWorkerPool) HasActiveWorkers() bool {
	return p.hasActiveWorkers
}

// IsSSHSession returns true if the pool is configured for an SSH session.
func (p *PullWorkerPool) IsSSHSession() bool {
	return p.isSSHSession
}

// Start launches dynamic worker goroutines matching the pool worker count.
func (p *PullWorkerPool) Start() {
	if p.hasActiveWorkers || p.workerCount <= 0 {
		return
	}

	p.hasActiveWorkers = true
	for id := 0; id < p.workerCount; id++ {
		p.waitGroup.Add(1)
		go p.workerLoop(id)
	}
}

func (p *PullWorkerPool) workerLoop(workerID int) {
	defer p.waitGroup.Done()
	for rec := range p.jobQueue {
		if p.isStopped() {
			return
		}

		p.executeJob(workerID, rec)
	}
}

func (p *PullWorkerPool) executeJob(workerID int, rec model.ScanRecord) {
	if p.bar != nil {
		p.bar.RegisterWorker(workerID, rec.RepoName)
	}

	ExecuteTrackedPullWithWorker(rec, p.bar, workerID)
	if p.bar != nil {
		p.bar.UnregisterWorker(workerID)
	}
}

func (p *PullWorkerPool) isStopped() bool {
	if p.bar == nil {
		return false
	}

	return p.bar.IsStopped()
}

// Submit queues a single repository record for pulling.
func (p *PullWorkerPool) Submit(rec model.ScanRecord) bool {
	if p.isStopped() {
		return false
	}

	p.jobQueue <- rec

	return true
}

// CloseQueue closes the job queue signaling workers to drain and finish.
func (p *PullWorkerPool) CloseQueue() {
	close(p.jobQueue)
}

// Wait blocks until all dynamic worker goroutines complete.
func (p *PullWorkerPool) Wait() *appfault.AppError {
	p.waitGroup.Wait()
	p.hasActiveWorkers = false

	return nil
}

// Dispatch sends records to the pool's job queue and closes it.
func (p *PullWorkerPool) Dispatch(records []model.ScanRecord) *appfault.AppError {
	for _, rec := range records {
		if p.isStopped() {
			break
		}

		p.jobQueue <- rec
	}

	p.CloseQueue()

	return nil
}

// Execute runs the full parallel pull pipeline with dynamic workers.
func (p *PullWorkerPool) Execute(records []model.ScanRecord) *appfault.AppError {
	if len(records) == 0 {
		return nil
	}

	p.Start()
	p.Dispatch(records)

	return p.Wait()
}

// RunDynamicPullParallel executes parallel pulls across repositories using dynamic CPU workers.
func RunDynamicPullParallel(records []model.ScanRecord, bar *PullProgressBar, requested int, isSSH bool) *appfault.AppError {
	if len(records) == 0 {
		return nil
	}

	pool := NewPullWorkerPool(requested, len(records), isSSH, bar)

	return pool.Execute(records)
}
