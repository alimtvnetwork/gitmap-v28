// Package cmd — shared helpers for `stale`, `dedupe`, `size`, `orphans`:
// parallel repo scanning and uniform --format=table|json|csv output.
// v6.71.0.
package cmd

import (
	"os"
	"path/filepath"
	"runtime"
)

// hygieneWorkers caps fan-out for parallel git invocations.
func hygieneWorkers() int {
	n := runtime.NumCPU()
	if n < 2 {
		return 2
	}

	if n > 8 {
		return 8
	}

	return n
}

// ScanForReposParallel discovers git repository directories concurrently.
func ScanForReposParallel(root string) []string {
	return scanForReposParallel(root)
}

// scanForReposParallel walks the immediate children of root and returns
// directories that contain a .git folder. The .git probe is fanned out
// across hygieneWorkers() goroutines so large directories stay snappy.
func scanForReposParallel(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}

	candidates := filterCandidateDirs(root, entries)

	return collectGitRepos(candidates)
}

func filterCandidateDirs(root string, entries []os.DirEntry) []string {
	candidates := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			candidates = append(candidates, filepath.Join(root, e.Name()))
		}
	}

	return candidates
}
