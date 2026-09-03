package worker

import (
	"context"
	"sync"
)

// TaskFunc is a generic function signature for processing a single item.
type TaskFunc[T any, R any] func(ctx context.Context, input T) (R, error)

// Result holds the output or error of a processed task.
type Result[R any] struct {
	Value R
	Err   error
}

// Pool represents a generic worker pool.
type Pool[T any, R any] struct {
	workerCount int
	taskFunc    TaskFunc[T, R]
}

// NewPool creates a new generic worker pool.
func NewPool[T any, R any](workerCount int, taskFunc TaskFunc[T, R]) *Pool[T, R] {
	if workerCount <= 0 {
		workerCount = 1
	}

	return &Pool[T, R]{
		workerCount: workerCount,
		taskFunc:    taskFunc,
	}
}

// Run executes the worker pool over a channel of inputs.
// It returns a channel of results.

func (p *Pool[T, R]) Run(ctx context.Context, inputs <-chan T) <-chan Result[R] {
	results := make(chan Result[R])
	var wg sync.WaitGroup

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case input, ok := <-inputs:
					if !ok {
						return
					}

					res, err := p.taskFunc(ctx, input)
					select {
					case <-ctx.Done():
						return
					case results <- Result[R]{Value: res, Err: err}:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
