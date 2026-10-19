package termgeneric

import (
	"context"
	"runtime"
	"sync"
)

// TaskFunc defines the worker execution signature for generic task processing.
type TaskFunc[T any, R any] func(ctx context.Context, item T) (R, error)

// TaskResult encapsulates the input, output, and error for a single processed task.
type TaskResult[T any, R any] struct {
	Input  T
	Output R
	Err    error
}

// WorkerPool manages a concurrent pool of worker goroutines processing items of type T and yielding R.
type WorkerPool[T any, R any] struct {
	workerCount int
	queueSize   int
	handler     TaskFunc[T, R]
	tasks       chan T
	results     chan TaskResult[T, R]
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	errMu       sync.Mutex
	errors      []error
	isStarted   bool
	isClosed    bool
}

// NewWorkerPool creates a new WorkerPool with bounded worker count and queue capacity.
func NewWorkerPool[T any, R any](ctx context.Context, workerCount int, queueSize int, handler TaskFunc[T, R]) *WorkerPool[T, R] {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}
	if queueSize <= 0 {
		queueSize = workerCount * 2
	}

	subCtx, cancel := context.WithCancel(ctx)

	return &WorkerPool[T, R]{
		workerCount: workerCount,
		queueSize:   queueSize,
		handler:     handler,
		tasks:       make(chan T, queueSize),
		results:     make(chan TaskResult[T, R], queueSize),
		ctx:         subCtx,
		cancel:      cancel,
		errors:      make([]error, 0),
	}
}

// Start launches the worker goroutines.
func (p *WorkerPool[T, R]) Start() {
	if p.isStarted {
		return
	}
	p.isStarted = true

	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.workerLoop()
	}
}

// Submit queues a task for execution. Returns true if queued, or false if pool is closed or context canceled.
func (p *WorkerPool[T, R]) Submit(item T) bool {
	if p.isClosed {
		return false
	}

	isCtxDone := p.ctx.Err() != nil
	if isCtxDone {
		return false
	}

	select {
	case <-p.ctx.Done():
		return false
	case p.tasks <- item:
		return true
	}
}

// Close gracefully closes the task queue and waits for all active workers to finish, then closes results channel.
func (p *WorkerPool[T, R]) Close() {
	if p.isClosed {
		return
	}
	p.isClosed = true

	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}

// Stop cancels the context and halts worker routines immediately.
func (p *WorkerPool[T, R]) Stop() {
	p.cancel()
	p.Close()
}

// Results returns the read-only results channel.
func (p *WorkerPool[T, R]) Results() <-chan TaskResult[T, R] {
	return p.results
}

// Errors returns all accumulated errors across workers.
func (p *WorkerPool[T, R]) Errors() []error {
	p.errMu.Lock()
	defer p.errMu.Unlock()

	errCopy := make([]error, len(p.errors))
	copy(errCopy, p.errors)
	return errCopy
}

// HasErrors indicates whether any task returned an error.
func (p *WorkerPool[T, R]) HasErrors() bool {
	p.errMu.Lock()
	defer p.errMu.Unlock()

	return len(p.errors) > 0
}

// ErrorCount returns the number of recorded errors.
func (p *WorkerPool[T, R]) ErrorCount() int {
	p.errMu.Lock()
	defer p.errMu.Unlock()

	return len(p.errors)
}

func (p *WorkerPool[T, R]) recordError(err error) {
	if err == nil {
		return
	}

	p.errMu.Lock()
	defer p.errMu.Unlock()

	p.errors = append(p.errors, err)
}

func (p *WorkerPool[T, R]) workerLoop() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case item, ok := <-p.tasks:
			if !ok {
				return
			}
			p.processItem(item)
		}
	}
}

func (p *WorkerPool[T, R]) processItem(item T) {
	output, err := p.handler(p.ctx, item)
	res := TaskResult[T, R]{
		Input:  item,
		Output: output,
		Err:    err,
	}

	if err != nil {
		p.recordError(err)
	}

	select {
	case <-p.ctx.Done():
		return
	case p.results <- res:
	}
}

// RunBatch executes tasks for all items in batch mode and returns aggregated results.
func RunBatch[T any, R any](ctx context.Context, items []T, workerCount int, handler TaskFunc[T, R]) ([]TaskResult[T, R], error) {
	if len(items) == 0 {
		return make([]TaskResult[T, R], 0), nil
	}

	pool := NewWorkerPool(ctx, workerCount, len(items), handler)
	pool.Start()
	feedItems(pool, items)
	pool.Close()

	results := drainResults(pool, len(items))
	if pool.HasErrors() {
		return results, pool.Errors()[0]
	}

	return results, nil
}

func feedItems[T any, R any](pool *WorkerPool[T, R], items []T) {
	for _, item := range items {
		isQueued := pool.Submit(item)
		if !isQueued {
			break
		}
	}
}

func drainResults[T any, R any](pool *WorkerPool[T, R], capacity int) []TaskResult[T, R] {
	results := make([]TaskResult[T, R], 0, capacity)
	for res := range pool.Results() {
		results = append(results, res)
	}

	return results
}
