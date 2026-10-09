package cmdautofix

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type scanFileResult struct {
	abs        string
	violations []Violation
	fromCache  bool
	scanned    bool
	entry      cacheEntry
}

// scanOneFile stats, cache-checks, reads, probes, and checks one file.
func scanOneFile(absPath string, root string, cats []Category, catNames []string, opts Options, cache map[string]cacheEntry) scanFileResult {
	res := scanFileResult{abs: absPath}
	rel := relSlash(root, absPath)

	info, err := os.Stat(absPath)
	if err != nil {
		return res
	}
	if entry, ok := cache[absPath]; ok && cacheHit(entry, info, catNames) {
		res.violations = filterViolations(entry.Violations, catNames)
		res.fromCache = true
		res.entry = entry
		return res
	}

	src, err := os.ReadFile(absPath)
	if err != nil {
		return res
	}
	if isBinaryProbe(src) {
		return res
	}

	for _, cat := range cats {
		if !matchesScope(rel, cat, opts) {
			continue
		}
		res.violations = append(res.violations, cat.Check(rel, src, opts)...)
	}

	res.entry = cacheEntry{
		ModUnix:    info.ModTime().Unix(),
		Size:       info.Size(),
		Categories: catNames,
		Violations: res.violations,
	}
	res.scanned = true
	return res
}

// Scan walks opts.Root, runs every selected category's Check in a worker
// pool, consults the scan cache, and returns the deterministic result.
func Scan(opts Options) (*ScanResult, error) {
	start := time.Now()
	if opts.Root == "" {
		opts.Root = "."
	}
	cats, err := selectedCategories(opts)
	if err != nil {
		return nil, err
	}
	for _, cat := range cats {
		if cat.Exec {
			if terr := catToolCheck(cat.Name); terr != nil {
				return nil, terr
			}
		}
	}

	catNames := make([]string, 0, len(cats))
	for _, cat := range cats {
		catNames = append(catNames, cat.Name)
	}

	var cache map[string]cacheEntry
	if opts.NoCache {
		cache = map[string]cacheEntry{}
	} else {
		cache = loadCache()
	}

	workers := opts.Workers
	if workers < 1 {
		workers = 1
	}

	fileCh := make(chan string, workers*4)
	resultCh := make(chan scanFileResult, workers*4)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for absPath := range fileCh {
				resultCh <- scanOneFile(absPath, opts.Root, cats, catNames, opts, cache)
			}
		}()
	}

	go func() {
		defer close(fileCh)
		_ = filepath.WalkDir(opts.Root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if d.Name() == ".git" {
					return filepath.SkipDir
				}
				return nil
			}
			fileCh <- path
			return nil
		})
	}()
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	out := &ScanResult{Root: opts.Root, Categories: catNames}
	newCache := make(map[string]cacheEntry, len(cache))
	for res := range resultCh {
		out.FilesScanned++
		if res.fromCache {
			out.FilesFromCache++
		}
		out.Violations = append(out.Violations, res.violations...)
		if res.fromCache || res.scanned {
			newCache[res.abs] = res.entry
		}
	}

	// Deterministic output regardless of worker scheduling.
	sort.Slice(out.Violations, func(i, j int) bool {
		if out.Violations[i].Path != out.Violations[j].Path {
			return out.Violations[i].Path < out.Violations[j].Path
		}
		if out.Violations[i].Category != out.Violations[j].Category {
			return out.Violations[i].Category < out.Violations[j].Category
		}
		return out.Violations[i].Line < out.Violations[j].Line
	})
	out.Violations = dedupeViolations(out.Violations)
	out.Elapsed = time.Since(start)

	if !opts.NoCache {
		saveCache(newCache)
	}

	return out, nil
}

// dedupeViolations drops exact duplicate findings (same path, category,
// line, detail).
