package utils

import (
	"sync"
)

// ProcessAsync processes items concurrently with a max concurrency limit.
func ProcessAsync(concurrency int, total int, processor func(index int)) {
	if total == 0 {
		return
	}
	sem := buildSemaphore(concurrency)
	var wg sync.WaitGroup
	executeWorkers(&wg, sem, total, processor)
	wg.Wait()
}

func buildSemaphore(concurrency int) chan struct{} {
	if concurrency < 1 {
		return make(chan struct{}, 1)
	}
	return make(chan struct{}, concurrency)
}

func executeWorkers(wg *sync.WaitGroup, sem chan struct{}, total int, processor func(int)) {
	for i := 0; i < total; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go executeSingleWorker(wg, sem, i, processor)
	}
}

func executeSingleWorker(wg *sync.WaitGroup, sem chan struct{}, idx int, processor func(int)) {
	defer wg.Done()
	processor(idx)
	<-sem
}
