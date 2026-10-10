package cmdsize

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhygiene"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
)

// dirSize returns the total byte count under path (best-effort).
func dirSize(path string) int64 {
	var total int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		total += info.Size()

		return nil
	})

	return total
}

// emitCSV writes a header + rows to stdout via encoding/csv.
func emitCSV(header []string, rows [][]string) {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write(header)
	for _, r := range rows {
		_ = w.Write(r)
	}

	w.Flush()
}

// emitJSON writes v as indented JSON to stdout.
func emitJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "  json encode: %v\n", err)
	}
}

// parseHygieneFormat normalizes and validates a --format value.
func parseHygieneFormat(s string) (cmdhygiene.HygieneFormatType, error) {
	switch s {
	case "", "table":
		return cmdhygiene.HygieneFormatTypeTable, nil
	case "json":
		return cmdhygiene.HygieneFormatTypeJSON, nil
	case "csv":
		return cmdhygiene.HygieneFormatTypeCSV, nil
	}

	return "", fmt.Errorf("invalid --format %q (want table|json|csv)", s)
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
func humanBytes(n int64) string {
	const k = 1024.0
	if n < int64(k) {
		return fmt.Sprintf("%d B", n)
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	f := float64(n) / k
	for _, u := range units {
		if f < k {
			return fmt.Sprintf("%.1f %s", f, u)
		}

		f /= k
	}

	return fmt.Sprintf("%.1f EB", f)
}

func runRepoWorker(wg *sync.WaitGroup, jobs <-chan string, results chan<- string) {
	defer wg.Done()
	for p := range jobs {
		if isGitRepo(p) {
			results <- p
		}
	}
}

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

func spawnRepoWorkers(jobs <-chan string, results chan<- string, count int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go runRepoWorker(&wg, jobs, results)
	}
	return &wg
}

func isGitRepo(path string) bool {
	return cmdhygiene.IsGitRepo(path)
}

func drainAndSortResults(results <-chan string, capacity int) []string {
	out := make([]string, 0, capacity)
	for r := range results {
		out = append(out, r)
	}

	sort.Strings(out)

	return out
}
