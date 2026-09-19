package cmdautomation

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunSearch executes multi-core parallel file search across repository files.
func RunSearch(opts SearchOptions) (SearchResult, *apperror.AppError) {
	start := time.Now()
	valErr := validateSearchOptions(&opts)
	if valErr != nil {
		return SearchResult{}, valErr
	}

	files := collectSearchFiles(opts.Dir, opts.Extensions)
	matches := dispatchSearchWorkers(files, opts)
	duration := time.Since(start)

	return SearchResult{
		Matches:    matches,
		TotalFiles: len(files),
		Duration:   duration,
		TotalHits:  len(matches),
	}, nil
}

func validateSearchOptions(opts *SearchOptions) *apperror.AppError {
	isEmpty := len(opts.Pattern) == 0
	if isEmpty {
		return apperror.NewValidationError("search pattern is required")
	}
	if opts.Dir == "" {
		opts.Dir = "."
	}
	if opts.Workers <= 0 {
		opts.Workers = runtime.NumCPU()
	}
	return nil
}

func collectSearchFiles(dir string, exts []string) []string {
	var files []string
	extMap := buildExtMap(exts)

	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && isExcludedDir(info.Name()) {
			return filepath.SkipDir
		}
		if !info.IsDir() && matchesExtensionFilter(p, extMap) {
			files = append(files, p)
		}
		return nil
	})
	return files
}

func buildExtMap(exts []string) map[string]bool {
	m := make(map[string]bool)
	for _, e := range exts {
		m[e] = true
	}
	return m
}

func matchesExtensionFilter(path string, extMap map[string]bool) bool {
	if len(extMap) == 0 {
		return IsPolyglotTextExtension(filepath.Ext(path))
	}
	return extMap[filepath.Ext(path)]
}

func dispatchSearchWorkers(files []string, opts SearchOptions) []SearchMatch {
	jobs := make(chan string, len(files))
	results := make(chan []SearchMatch, opts.Workers)
	var wg sync.WaitGroup

	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go searchWorkerRoutine(jobs, results, opts, &wg)
	}

	for _, f := range files {
		jobs <- f
	}
	close(jobs)

	wg.Wait()
	close(results)

	return gatherSearchResults(results)
}
