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
	coldRecords, err := filterColdIgnoreRecords(records, ttl)
	if err != nil || len(coldRecords) == 0 {
		handle.done <- nil
		return
	}
	issues := scanColdRecordsSequentially(coldRecords)
	handle.done <- issues
}

func filterColdIgnoreRecords(records []model.ScanRecord, ttl time.Duration) ([]model.ScanRecord, error) {
	cold, err := store.FilterReposNeedingIgnoreCheck(records, ttl)
	if err != nil {
		return records, nil
	}
	return cold, nil
}

func scanColdRecordsSequentially(records []model.ScanRecord) []IgnoreRepoIssue {
	var issues []IgnoreRepoIssue
	for i, rec := range records {
		yieldBetweenIgnoreAudits(i)
		issue := inspectAndCacheRepoIgnore(rec)
		if issue.HasIssues() {
			issues = append(issues, issue)
		}
	}
	return sortIgnoreIssues(issues)
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
