// Package vscodepm — optimize.go implements project deduplication and clearance.
package vscodepm

import (
	"strings"
)

// OptimizeSummary holds metrics from an optimize run.
type OptimizeSummary struct {
	Removed   int
	Remaining int
}

// OptimizeProjects removes duplicate entries sharing the same RootPath.
func OptimizeProjects(exceptList []string, dryRun bool) (OptimizeSummary, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return OptimizeSummary{}, err
	}
	return OptimizeProjectsAt(path, exceptList, dryRun)
}

// OptimizeProjectsAt deduplicates entries in a specified file.
func OptimizeProjectsAt(filePath string, exceptList []string, dryRun bool) (OptimizeSummary, error) {
	entries, err := readEntries(filePath)
	if err != nil {
		return OptimizeSummary{}, err
	}
	deduped, removed := deduplicateEntries(entries, exceptList)
	if !dryRun && removed > 0 {
		if err := writeEntriesAtomic(filePath, deduped); err != nil {
			return OptimizeSummary{}, err
		}
	}
	return OptimizeSummary{Removed: removed, Remaining: len(deduped)}, nil
}

func deduplicateEntries(entries []Entry, exceptList []string) ([]Entry, int) {
	seen := make(map[string]int)
	deduped := make([]Entry, 0, len(entries))
	removed := 0

	for _, e := range entries {
		key := normalizePath(e.RootPath)
		if idx, found := seen[key]; found && !isEntryExcepted(e, exceptList) {
			deduped[idx] = mergeExistingEntry(deduped[idx], e)
			removed++
			continue
		}
		seen[key] = len(deduped)
		deduped = append(deduped, e)
	}
	return deduped, removed
}

func mergeExistingEntry(primary Entry, dup Entry) Entry {
	primary.Paths = unionPaths(primary.Paths, dup.Paths)
	primary.Tags = unionPaths(primary.Tags, dup.Tags)
	return primary
}

// ClearProjects removes entries while preserving those in exceptList.
func ClearProjects(exceptList []string, onlyMissing, dryRun bool) (OptimizeSummary, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return OptimizeSummary{}, err
	}
	return ClearProjectsAt(path, exceptList, onlyMissing, dryRun)
}

// ClearProjectsAt cleans entries from a specific projects.json file.
func ClearProjectsAt(filePath string, exceptList []string, onlyMissing, dryRun bool) (OptimizeSummary, error) {
	entries, err := readEntries(filePath)
	if err != nil {
		return OptimizeSummary{}, err
	}
	remaining, removed := filterClearEntries(entries, exceptList, onlyMissing)
	if !dryRun && removed > 0 {
		if err := writeEntriesAtomic(filePath, remaining); err != nil {
			return OptimizeSummary{}, err
		}
	}
	return OptimizeSummary{Removed: removed, Remaining: len(remaining)}, nil
}

func filterClearEntries(entries []Entry, exceptList []string, onlyMissing bool) ([]Entry, int) {
	remaining := make([]Entry, 0, len(entries))
	removed := 0

	for _, e := range entries {
		if isEntryExcepted(e, exceptList) {
			remaining = append(remaining, e)
			continue
		}
		if onlyMissing && dirExists(e.RootPath) {
			remaining = append(remaining, e)
			continue
		}
		removed++
	}
	return remaining, removed
}

func isEntryExcepted(e Entry, exceptList []string) bool {
	for _, ex := range exceptList {
		if ex == "" {
			continue
		}
		if strings.EqualFold(e.Name, ex) || strings.EqualFold(e.RootPath, ex) {
			return true
		}
		if strings.Contains(strings.ToLower(e.RootPath), strings.ToLower(ex)) {
			return true
		}
	}
	return false
}
