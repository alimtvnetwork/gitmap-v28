package cmdautomation

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

	files := collectSearchFiles(opts)
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
	opts.Pattern = cleanSearchPattern(opts.Pattern)
	if len(opts.Pattern) == 0 {
		return apperror.NewValidationError("search pattern is required")
	}
	autoPromoteRegex(opts)
	if err := validateRegexOption(opts); err != nil {
		return err
	}
	normalizeSearchDefaults(opts)
	return nil
}

// cleanSearchPattern removes outer wrapping quotes and shell escapes from search patterns.
func cleanSearchPattern(pattern string) string {
	trimmed := strings.TrimSpace(pattern)
	if len(trimmed) < 2 {
		return trimmed
	}

	trimmed = stripOuterQuotePairs(trimmed)
	trimmed = unescapeInternalQuotes(trimmed)

	return trimmed
}

func stripOuterQuotePairs(pattern string) string {
	current := pattern
	for {
		stripped, hasChanged := tryStripOuterPair(current)
		if !hasChanged || len(stripped) < 2 {
			return stripped
		}
		current = stripped
	}
}

func tryStripOuterPair(val string) (string, bool) {
	if stripped, hasMatch := tryStripEscapedQuotes(val); hasMatch {
		return stripped, true
	}
	if stripped, hasMatch := tryStripRegularQuotes(val); hasMatch {
		return stripped, true
	}
	return val, false
}

func tryStripEscapedQuotes(val string) (string, bool) {
	if strings.HasPrefix(val, `\"`) && strings.HasSuffix(val, `\"`) && len(val) >= 4 {
		return strings.TrimSuffix(strings.TrimPrefix(val, `\"`), `\"`), true
	}
	if strings.HasPrefix(val, `\'`) && strings.HasSuffix(val, `\'`) && len(val) >= 4 {
		return strings.TrimSuffix(strings.TrimPrefix(val, `\'`), `\'`), true
	}
	return val, false
}

func tryStripRegularQuotes(val string) (string, bool) {
	if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) && len(val) >= 2 {
		return strings.TrimSuffix(strings.TrimPrefix(val, `"`), `"`), true
	}
	if strings.HasPrefix(val, `'`) && strings.HasSuffix(val, `'`) && len(val) >= 2 {
		return strings.TrimSuffix(strings.TrimPrefix(val, `'`), `'`), true
	}
	if strings.HasPrefix(val, "`") && strings.HasSuffix(val, "`") && len(val) >= 2 {
		return strings.TrimSuffix(strings.TrimPrefix(val, "`"), "`"), true
	}
	return val, false
}

func unescapeInternalQuotes(val string) string {
	result := strings.ReplaceAll(val, `\"`, `"`)
	return strings.ReplaceAll(result, `\'`, `'`)
}

func validateRegexOption(opts *SearchOptions) *apperror.AppError {
	if !opts.IsRegex {
		return nil
	}
	_, err := GetRegex(opts.Pattern, opts.IsCaseInsensitive)
	return err
}

func normalizeSearchDefaults(opts *SearchOptions) {
	if opts.Dir == "" {
		opts.Dir = "."
	}
	if opts.Workers <= 0 {
		opts.Workers = runtime.NumCPU()
	}
}

func collectSearchFiles(opts SearchOptions) []string {
	if explicitFiles, hasExplicit := resolveSearchCandidateFiles(opts.Dir); hasExplicit {
		return explicitFiles
	}
	var files []string
	extMap := buildExtMap(opts.Extensions)
	exclusions := loadActiveExclusions()
	maxJsonBytes := resolveMaxJsonBytes(opts.MaxJsonKb)

	walkRoot := opts.Dir
	if walkRoot == "" {
		walkRoot = "."
	}
	_ = filepath.Walk(walkRoot, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && isExcludedDir(info.Name()) {
			return filepath.SkipDir
		}
		if !info.IsDir() && isSearchCandidate(p, info, opts, extMap, exclusions, maxJsonBytes) {
			files = append(files, p)
		}
		return nil
	})
	return files
}

func resolveMaxJsonBytes(maxKb int) int64 {
	if maxKb <= 0 {
		return 500 * 1024
	}
	return int64(maxKb * 1024)
}

func isSearchCandidate(p string, info os.FileInfo, opts SearchOptions, extMap map[string]bool, exclusions []string, maxJsonBytes int64) bool {
	if !matchesExtensionFilter(p, extMap) {
		return false
	}
	ext := filepath.Ext(p)
	if !opts.IncludeBinaries && IsBinaryExtension(ext) {
		return false
	}
	if !opts.IncludeLargeJson && ext == ".json" && info.Size() > maxJsonBytes {
		return false
	}
	rel := filepath.ToSlash(p)
	if IsPathExcluded(rel, exclusions) {
		return false
	}
	return true
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

	startSearchWorkers(&wg, jobs, results, opts)
	feedSearchJobs(jobs, files)
	go waitAndCloseResults(&wg, results)

	return gatherSearchResults(results)
}

func startSearchWorkers(wg *sync.WaitGroup, jobs <-chan string, results chan<- []SearchMatch, opts SearchOptions) {
	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go searchWorkerRoutine(jobs, results, opts, wg)
	}
}

func feedSearchJobs(jobs chan<- string, files []string) {
	for _, f := range files {
		jobs <- f
	}
	close(jobs)
}

func waitAndCloseResults(wg *sync.WaitGroup, results chan<- []SearchMatch) {
	wg.Wait()
	close(results)
}
