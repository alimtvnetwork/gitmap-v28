// Package cmd — shared helpers for `stale`, `dedupe`, `size`, `orphans`:
// parallel repo scanning and uniform --format=table|json|csv output.
// v6.71.0.
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
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

// scanForReposParallel walks the immediate children of root and returns
// directories that contain a .git folder. The .git probe is fanned out
// across hygieneWorkers() goroutines so large directories stay snappy.
func scanForReposParallel(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	candidates := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			candidates = append(candidates, filepath.Join(root, e.Name()))
		}
	}
	jobs := make(chan string)
	results := make(chan string, len(candidates))
	var wg sync.WaitGroup
	for i := 0; i < hygieneWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if isGitRepo(p) {
					results <- p
				}
			}
		}()
	}
	for _, c := range candidates {
		jobs <- c
	}
	close(jobs)
	wg.Wait()
	close(results)
	out := make([]string, 0, len(candidates))
	for r := range results {
		out = append(out, r)
	}
	sort.Strings(out)

	return out
}

// mapReposParallel runs fn against every repo using hygieneWorkers().
// Results are returned in the same order as repos. nil entries are dropped.
func mapReposParallel[T any](repos []string, fn func(string) (T, bool)) []T {
	type indexed struct {
		i  int
		v  T
		ok bool
	}
	jobs := make(chan int)
	results := make(chan indexed, len(repos))
	var wg sync.WaitGroup
	for i := 0; i < hygieneWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				v, ok := fn(repos[idx])
				results <- indexed{i: idx, v: v, ok: ok}
			}
		}()
	}
	for i := range repos {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	close(results)
	buf := make([]indexed, 0, len(repos))
	for r := range results {
		if r.ok {
			buf = append(buf, r)
		}
	}
	sort.Slice(buf, func(a, b int) bool { return buf[a].i < buf[b].i })
	out := make([]T, 0, len(buf))
	for _, b := range buf {
		out = append(out, b.v)
	}

	return out
}

// hygieneFormat is the value of --format; "" means default table output.
type hygieneFormat string

const (
	hygieneFormatTable hygieneFormat = ""
	hygieneFormatJSON  hygieneFormat = "json"
	hygieneFormatCSV   hygieneFormat = "csv"
)

// parseHygieneFormat normalizes and validates a --format value.
func parseHygieneFormat(s string) (hygieneFormat, error) {
	switch s {
	case "", "table":
		return hygieneFormatTable, nil
	case "json":
		return hygieneFormatJSON, nil
	case "csv":
		return hygieneFormatCSV, nil
	}

	return "", fmt.Errorf("invalid --format %q (want table|json|csv)", s)
}

// emitJSON writes v as indented JSON to stdout.
func emitJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "  json encode: %v\n", err)
	}
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
