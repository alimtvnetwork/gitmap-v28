package cmd

import (
	"sort"
	"sync"
)

// collectGitRepos fans out git repository validation across hygieneWorkers().
func collectGitRepos(candidates []string) []string {
	jobs := make(chan string)
	results := make(chan string, len(candidates))
	wg := spawnRepoWorkers(jobs, results, hygieneWorkers())

	for _, c := range candidates {
		jobs <- c
	}

	close(jobs)
	wg.Wait()
	close(results)

	return drainAndSortResults(results, len(candidates))
}

func spawnRepoWorkers(jobs <-chan string, results chan<- string, count int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go runRepoWorker(&wg, jobs, results)
	}

	return &wg
}

func runRepoWorker(wg *sync.WaitGroup, jobs <-chan string, results chan<- string) {
	defer wg.Done()
	for p := range jobs {
		if isGitRepo(p) {
			results <- p
		}
	}
}

func drainAndSortResults(results <-chan string, capacity int) []string {
	out := make([]string, 0, capacity)
	for r := range results {
		out = append(out, r)
	}

	sort.Strings(out)

	return out
}
