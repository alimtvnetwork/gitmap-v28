// Package vscodepm — optimize.go implements project deduplication and clearance.
package vscodepm

import (
	"fmt"
	"path/filepath"
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
		if idx, found := seen[key]; found && !isEntryExcepted(e, exceptList, len(deduped)+1) {
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
	summary, _, err := ClearProjectsWithTargets(exceptList, onlyMissing, dryRun)
	return summary, err
}

// ClearProjectsAt cleans entries from a specific projects.json file.
func ClearProjectsAt(filePath string, exceptList []string, onlyMissing, dryRun bool) (OptimizeSummary, error) {
	summary, _, err := ClearProjectsWithTargetsAt(filePath, exceptList, onlyMissing, dryRun)
	return summary, err
}

// ClearProjectsWithTargets returns both the summary and targeted entries.
func ClearProjectsWithTargets(exceptList []string, onlyMissing, dryRun bool) (OptimizeSummary, []Entry, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return OptimizeSummary{}, nil, err
	}
	return ClearProjectsWithTargetsAt(path, exceptList, onlyMissing, dryRun)
}

// ClearProjectsWithTargetsAt cleans entries and returns targeted entries.
func ClearProjectsWithTargetsAt(filePath string, exceptList []string, onlyMissing, dryRun bool) (OptimizeSummary, []Entry, error) {
	entries, err := readEntries(filePath)
	if err != nil {
		return OptimizeSummary{}, nil, err
	}
	targets, remaining := GetClearTargets(entries, exceptList, onlyMissing)
	if !dryRun && len(targets) > 0 {
		if err := writeEntriesAtomic(filePath, remaining); err != nil {
			return OptimizeSummary{}, nil, err
		}
	}
	return OptimizeSummary{Removed: len(targets), Remaining: len(remaining)}, targets, nil
}

// GetClearTargets partitions entries into targets to clear and remaining entries.
func GetClearTargets(entries []Entry, exceptList []string, onlyMissing bool) ([]Entry, []Entry) {
	targets := make([]Entry, 0)
	remaining := make([]Entry, 0)
	for i, e := range entries {
		if isEntryExcepted(e, exceptList, i+1) {
			remaining = append(remaining, e)
			continue
		}
		if onlyMissing && dirExists(e.RootPath) {
			remaining = append(remaining, e)
			continue
		}
		targets = append(targets, e)
	}
	return targets, remaining
}

func filterClearEntries(entries []Entry, exceptList []string, onlyMissing bool) ([]Entry, int) {
	targets, remaining := GetClearTargets(entries, exceptList, onlyMissing)
	return remaining, len(targets)
}

func isEntryExcepted(e Entry, exceptList []string, index int) bool {
	idStr := fmt.Sprintf("%d", index)
	idPad := fmt.Sprintf("%02d", index)
	slug := filepath.Base(e.RootPath)
	lowName := strings.ToLower(e.Name)
	lowSlug := strings.ToLower(slug)
	lowPath := strings.ToLower(e.RootPath)

	for _, rawEx := range exceptList {
		ex := strings.ToLower(strings.TrimSpace(rawEx))
		if ex == "" {
			continue
		}
		if ex == idStr || ex == idPad || ex == lowName || ex == lowSlug || ex == lowPath {
			return true
		}
		if strings.HasPrefix(lowName, ex) || strings.HasPrefix(lowSlug, ex) {
			return true
		}
		if strings.Contains(ex, "/") || strings.Contains(ex, "\\") {
			cleanEx := strings.ToLower(filepath.Clean(ex))
			if lowPath == cleanEx || strings.HasSuffix(lowPath, cleanEx) {
				return true
			}
		}
	}
	return false
}
