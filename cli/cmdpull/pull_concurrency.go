package cmdpull

import (
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
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
		issues = appendColdIssueIfAny(rec, i, seenPath, issues)
	}

	return sortIgnoreIssues(issues)
}

func appendColdIssueIfAny(rec model.ScanRecord, idx int, seen map[string]bool, list []IgnoreRepoIssue) []IgnoreRepoIssue {
	pathKey := CanonicalRepoPathKey(rec.AbsolutePath)
	if isPathAlreadySeen(pathKey, seen) {
		return list
	}
	yieldBetweenIgnoreAudits(idx)
	issue := inspectAndCacheRepoIgnore(rec)
	cleaned := sanitizeIgnoreIssue(issue)
	if cleaned.HasIssues() {
		return append(list, cleaned)
	}

	return list
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

// IsSystemCPUIdle evaluates if the host CPU is under 50% load to safely scale concurrency.
func IsSystemCPUIdle() bool {
	metrics, err := cmdos.GetHardwareMetrics()
	if err != nil {
		return true
	}

	return metrics.CPU.LoadPercent < 50
}

// HasHighCPUPressure evaluates if host CPU utilization exceeds 80% to throttle workers.
func HasHighCPUPressure() bool {
	metrics, err := cmdos.GetHardwareMetrics()
	if err != nil {
		return false
	}

	return metrics.CPU.LoadPercent > 80
}

// ResolveDynamicPullWorkers calculates dynamic concurrency based on CPU cores and session type.
func ResolveDynamicPullWorkers(requestedWorkers int, isSSHSession bool) int {
	if isSSHSession {
		return resolveSSHWorkers(requestedWorkers)
	}
	if requestedWorkers > 0 {
		return requestedWorkers
	}

	return resolveLocalDynamicWorkers()
}

func resolveSSHWorkers(requestedWorkers int) int {
	if requestedWorkers > 0 && requestedWorkers < 2 {
		return requestedWorkers
	}

	return 2
}

func resolveLocalDynamicWorkers() int {
	if HasHighCPUPressure() {
		return 2
	}
	cores := runtime.NumCPU()
	if IsSystemCPUIdle() {
		return resolveIdleCoreWorkers(cores)
	}

	return resolveModerateCoreWorkers(cores)
}

func resolveIdleCoreWorkers(cores int) int {
	if cores >= 16 {
		return 16
	}
	if cores >= 8 {
		return 12
	}
	if cores >= 4 {
		return 4
	}

	return resolveLowCoreWorkers(cores)
}

func resolveModerateCoreWorkers(cores int) int {
	if cores >= 8 {
		return 8
	}
	if cores >= 4 {
		return 4
	}

	return resolveLowCoreWorkers(cores)
}

func resolveLowCoreWorkers(cores int) int {
	if cores > 1 {
		return 2
	}

	return 1
}

// ClampWorkers bounds worker count to at least 1 and at most repoCount (if repoCount > 0).
func ClampWorkers(workers, repoCount int) int {
	if workers < 1 {
		return 1
	}
	if repoCount > 0 && workers > repoCount {
		return repoCount
	}

	return workers
}

// ResolveDynamicPullWorkersForRepos calculates dynamic pull concurrency bounded by repository count.
func ResolveDynamicPullWorkersForRepos(requestedWorkers, repoCount int, isSSHSession bool) int {
	workers := ResolveDynamicPullWorkers(requestedWorkers, isSSHSession)

	return ClampWorkers(workers, repoCount)
}
