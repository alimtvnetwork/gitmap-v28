package cmdpull

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CalculateIgnoreWorkersForPull returns the throttled worker count for background ignore scanning.
// Always returns 1 during pull operations to eliminate process contention and CPU starvation.
func CalculateIgnoreWorkersForPull(total int) int {
	return 1
}

// StartThrottledAsyncIgnoreScan initiates a low-priority, cached background ignore inspection.
func StartThrottledAsyncIgnoreScan(records []model.ScanRecord, ttl time.Duration) *IgnoreScanHandle {
	handle := &IgnoreScanHandle{
		done: make(chan []IgnoreRepoIssue, 1),
	}
	if len(records) == 0 {
		handle.done <- nil
		return handle
	}
	go runThrottledIgnoreScanWorkflow(records, ttl, handle)
	return handle
}

func runThrottledIgnoreScanWorkflow(records []model.ScanRecord, ttl time.Duration, handle *IgnoreScanHandle) {
	defer func() {
		if r := recover(); r != nil {
			handle.done <- nil
		}
	}()
	coldRecords, err := filterColdIgnoreRecords(records, ttl)
	if err != nil || len(coldRecords) == 0 {
		handle.done <- nil
		return
	}
	issues := scanColdRecordsSequentially(coldRecords)
	handle.done <- issues
}

func filterColdIgnoreRecords(records []model.ScanRecord, ttl time.Duration) ([]model.ScanRecord, error) {
	uniqueRecords := deduplicateRecordsByPath(records)
	cold, err := store.FilterReposNeedingIgnoreCheck(uniqueRecords, ttl)
	if err != nil {
		return uniqueRecords, nil
	}
	return deduplicateRecordsByPath(cold), nil
}

func deduplicateRecordsByPath(records []model.ScanRecord) []model.ScanRecord {
	unique := make([]model.ScanRecord, 0, len(records))
	seen := make(map[string]bool, len(records))
	for _, rec := range records {
		key := CanonicalRepoPathKey(rec.AbsolutePath)
		if isPathAlreadySeen(key, seen) {
			continue
		}
		unique = append(unique, rec)
	}
	return unique
}

func scanColdRecordsSequentially(records []model.ScanRecord) []IgnoreRepoIssue {
	var issues []IgnoreRepoIssue
	seenPath := make(map[string]bool, len(records))

	for i, rec := range records {
		pathKey := CanonicalRepoPathKey(rec.AbsolutePath)
		if isPathAlreadySeen(pathKey, seenPath) {
			continue
		}
		yieldBetweenIgnoreAudits(i)
		issue := inspectAndCacheRepoIgnore(rec)
		cleaned := sanitizeIgnoreIssue(issue)
		if cleaned.HasIssues() {
			issues = append(issues, cleaned)
		}
	}

	return sortIgnoreIssues(issues)
}

func sanitizeIgnoreIssue(issue IgnoreRepoIssue) IgnoreRepoIssue {
	issue.TrackedPaths = deduplicateStringList(issue.TrackedPaths)
	issue.DuplicatePatterns = deduplicateStringList(issue.DuplicatePatterns)
	return issue
}

func deduplicateStringList(items []string) []string {
	if len(items) <= 1 {
		return items
	}
	seen := make(map[string]bool, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		result = append(result, item)
	}
	return result
}

func isPathAlreadySeen(pathKey string, seenPath map[string]bool) bool {
	if pathKey == "" {
		return false
	}
	if seenPath[pathKey] {
		return true
	}
	seenPath[pathKey] = true

	return false
}

func yieldBetweenIgnoreAudits(index int) {
	if index > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

func inspectAndCacheRepoIgnore(rec model.ScanRecord) IgnoreRepoIssue {
	start := time.Now()
	issue := inspectRepoForIgnoreIssues(rec.AbsolutePath, rec.RepoName)
	status := "clean"
	if issue.HasIssues() {
		status = "has_issues"
	}
	_ = store.RecordIgnoreCheckResult(rec.AbsolutePath, rec.RepoName, status, 0, time.Since(start))
	return issue
}
