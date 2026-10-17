package parallel

import (
	"runtime"
	"sync"
)

// DefaultMinConcurrency is the minimum fallback concurrency worker count.
const DefaultMinConcurrency = 4

// DefaultMaxConcurrency is the maximum worker count for parallel execution.
const DefaultMaxConcurrency = 20

// ResolveConcurrency calculates an optimal worker count based on runtime.NumCPU().
// If requested > 0, it clamps requested between 1 and DefaultMaxConcurrency.
// If requested <= 0, it uses runtime.NumCPU(), defaulting to at least DefaultMinConcurrency
// (4) and capped at DefaultMaxConcurrency (20).
func ResolveConcurrency(requested int) int {
	if requested > DefaultMaxConcurrency {
		return DefaultMaxConcurrency
	}
	if requested > 0 {
		return requested
	}

	cpus := runtime.NumCPU()
	if cpus < DefaultMinConcurrency {
		return DefaultMinConcurrency
	}
	if cpus > DefaultMaxConcurrency {
		return DefaultMaxConcurrency
	}

	return cpus
}

// ResolveDefaultConcurrency resolves the default concurrency level based on runtime.NumCPU().
func ResolveDefaultConcurrency() int {
	return ResolveConcurrency(0)
}

// ExecuteParallel runs worker across items concurrently using up to concurrency workers.
// workerID is passed to worker (0-indexed). Any non-nil errors returned by worker are collected
// and returned as a slice. If all items succeed, nil is returned.
func ExecuteParallel[T any](items []T, concurrency int, worker func(workerID int, item T) error) []error {
	if len(items) == 0 {
		return nil
	}

	effectiveConcurrency := concurrency
	if effectiveConcurrency <= 0 {
		effectiveConcurrency = ResolveDefaultConcurrency()
	}
	if effectiveConcurrency > len(items) {
		effectiveConcurrency = len(items)
	}

	type job struct {
		index int
		item  T
	}

	jobs := make(chan job, len(items))
	for i, item := range items {
		jobs <- job{index: i, item: item}
	}
	close(jobs)

	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	wg.Add(effectiveConcurrency)
	for w := 0; w < effectiveConcurrency; w++ {
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				if err := worker(workerID, j.item); err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
				}
			}
		}(w)
	}

	wg.Wait()

	if len(errs) == 0 {
		return nil
	}

	return errs
}
