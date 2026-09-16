package cmd

import (
	"sort"
	"sync"
)

type hygieneIndexedResult[T any] struct {
	i       int
	v       T
	isFound bool
}

// mapReposParallel runs fn against every repo using hygieneWorkers().
// Results are returned in the same order as repos.
func mapReposParallel[T any](repos []string, fn func(string) (T, bool)) []T {
	jobs := make(chan int)
	results := make(chan hygieneIndexedResult[T], len(repos))
	wg := spawnMapWorkers(repos, jobs, results, fn, hygieneWorkers())

	for i := range repos {
		jobs <- i
	}

	close(jobs)
	wg.Wait()
	close(results)

	return collectIndexedResults(results, len(repos))
}

func spawnMapWorkers[T any](repos []string, jobs <-chan int, results chan<- hygieneIndexedResult[T], fn func(string) (T, bool), count int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go runMapWorker(&wg, repos, jobs, results, fn)
	}

	return &wg
}

func runMapWorker[T any](wg *sync.WaitGroup, repos []string, jobs <-chan int, results chan<- hygieneIndexedResult[T], fn func(string) (T, bool)) {
	defer wg.Done()
	for idx := range jobs {
		v, isFound := fn(repos[idx])
		results <- hygieneIndexedResult[T]{i: idx, v: v, isFound: isFound}
	}
}

func collectIndexedResults[T any](results <-chan hygieneIndexedResult[T], capacity int) []T {
	buf := drainIndexedBuffer(results, capacity)
	sort.Slice(buf, func(a, b int) bool { return buf[a].i < buf[b].i })

	return unwrapIndexedValues(buf)
}

func drainIndexedBuffer[T any](results <-chan hygieneIndexedResult[T], capacity int) []hygieneIndexedResult[T] {
	buf := make([]hygieneIndexedResult[T], 0, capacity)
	for r := range results {
		if r.isFound {
			buf = append(buf, r)
		}
	}

	return buf
}

func unwrapIndexedValues[T any](buf []hygieneIndexedResult[T]) []T {
	out := make([]T, 0, len(buf))
	for _, b := range buf {
		out = append(out, b.v)
	}

	return out
}
