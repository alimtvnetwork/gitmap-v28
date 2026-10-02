package vscodepm

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ProjectOptimizationAdvice provides relocation and deduplication advice for duplicate repository checkouts.
type ProjectOptimizationAdvice struct {
	ProjectName   string
	DuplicatePath string
	CanonicalPath string
	Advice        string
}

// OptimizeSummary contains counts of removed and remaining entries, plus relocation advice.
type OptimizeSummary struct {
	RemovedDuplicates int
	RemovedMissing    int
	Removed           int
	Remaining         int
	Advices           []ProjectOptimizationAdvice
}

// OptimizeProjects cleans up duplicate and missing entries in projects.json and advises on relocations.
func OptimizeProjects(exceptList []string, dryRun bool) (OptimizeSummary, error) {
	path, err := ProjectsJSONPath()
	if err != nil {
		return OptimizeSummary{}, err
	}

	return OptimizeProjectsAt(path, exceptList, dryRun)
}

// OptimizeProjectsAt deduplicates entries and prunes missing entries in a specified file.
func OptimizeProjectsAt(filePath string, exceptList []string, dryRun bool) (OptimizeSummary, error) {
	entries, err := readEntries(filePath)
	if err != nil {
		return OptimizeSummary{}, err
	}

	validEntries, missingCount := pruneMissingEntries(entries, exceptList)
	deduped, dupCount := deduplicateEntries(validEntries, exceptList)
	advices := generateRelocationAdvice(deduped)

	totalRemoved := missingCount + dupCount
	summary := OptimizeSummary{
		RemovedDuplicates: dupCount,
		RemovedMissing:    missingCount,
		Removed:           totalRemoved,
		Remaining:         len(deduped),
		Advices:           advices,
	}

	if isWriteEnabled(dryRun, totalRemoved) {
		return commitOptimizedEntries(filePath, deduped, summary)
	}

	return summary, nil
}

func isWriteEnabled(dryRun bool, totalRemoved int) bool {
	return !dryRun && totalRemoved > 0
}

func pruneMissingEntries(entries []Entry, exceptList []string) ([]Entry, int) {
	valid := make([]Entry, 0, len(entries))
	missing := 0
	for i, e := range entries {
		if isEntryExcepted(e, exceptList, i+1) {
			valid = append(valid, e)
			continue
		}
		if !dirExists(e.RootPath) {
			missing++
			continue
		}
		valid = append(valid, e)
	}
	return valid, missing
}

func generateRelocationAdvice(entries []Entry) []ProjectOptimizationAdvice {
	byBaseName := make(map[string][]Entry)
	for _, e := range entries {
		base := strings.ToLower(filepath.Base(e.RootPath))
		byBaseName[base] = append(byBaseName[base], e)
	}

	var advices []ProjectOptimizationAdvice
	for base, group := range byBaseName {
		if len(group) <= 1 {
			continue
		}
		canonical := pickCanonicalEntry(group)
		for _, e := range group {
			if normalizePath(e.RootPath) == normalizePath(canonical.RootPath) {
				continue
			}
			advices = append(advices, ProjectOptimizationAdvice{
				ProjectName:   base,
				DuplicatePath: e.RootPath,
				CanonicalPath: canonical.RootPath,
				Advice:        fmt.Sprintf("Relocate work from duplicate directory %q to canonical %q", e.RootPath, canonical.RootPath),
			})
		}
	}
	return advices
}

func pickCanonicalEntry(group []Entry) Entry {
	for _, e := range group {
		low := strings.ToLower(e.RootPath)
		if strings.Contains(low, "\\work\\") || strings.Contains(low, "/work/") {
			return e
		}
	}
	return group[0]
}

func commitOptimizedEntries(filePath string, deduped []Entry, summary OptimizeSummary) (OptimizeSummary, error) {
	if err := writeEntriesAtomic(filePath, deduped); err != nil {
		return summary, err
	}

	return summary, nil
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
	if err := maybeWriteRemaining(filePath, remaining, len(targets), dryRun); err != nil {
		return OptimizeSummary{}, nil, err
	}

	return OptimizeSummary{Removed: len(targets), Remaining: len(remaining)}, targets, nil
}

func maybeWriteRemaining(filePath string, remaining []Entry, targetCount int, dryRun bool) error {
	if dryRun || targetCount == 0 {
		return nil
	}

	return writeEntriesAtomic(filePath, remaining)
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

		if ex == idStr || ex == idPad || strings.EqualFold(ex, e.Name) || strings.EqualFold(ex, slug) || strings.EqualFold(ex, e.RootPath) {
			return true
		}

		if strings.HasPrefix(lowName, ex) || strings.HasPrefix(lowSlug, ex) {
			return true
		}

		if matchesPathException(lowPath, ex) {
			return true
		}
	}

	return false
}

func matchesPathException(lowPath, ex string) bool {
	if !strings.Contains(ex, "/") && !strings.Contains(ex, "\\") {
		return false
	}

	cleanEx := strings.ToLower(filepath.Clean(ex))

	return lowPath == cleanEx || strings.HasSuffix(lowPath, cleanEx)
}
