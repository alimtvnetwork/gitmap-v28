package cmdpipeline

import (
	"path/filepath"
	"sync"
)

func populateFailedRunsPayload(repo string, failedRuns []ghRunItem, p *PipelineErrorLogsPayload) {
	initFailedRunTopLevel(p, failedRuns[0])
	p.FailedRuns = fetchAllFailedRunsParallel(repo, failedRuns)
	p.FailedTree = BuildFailedChecksTree(*p, false)
	p.SectionFailures = extractAllSectionFailures(p.FailedRuns)
	p.CombinedErrors = formatCombinedSectionFailures(p.SectionFailures)
	p.ErrorLogs = formatAggregatedErrorLogs(p.FailedRuns)
}

func fetchAllFailedRunsParallel(repo string, failedRuns []ghRunItem) []FailedRunItem {
	total := len(failedRuns)
	if total == 0 {
		return nil
	}
	if total == 1 {
		return []FailedRunItem{fetchAndBuildFailedRunItem(repo, failedRuns[0])}
	}

	return executeParallelFetchWorkers(repo, failedRuns)
}

func executeParallelFetchWorkers(repo string, failedRuns []ghRunItem) []FailedRunItem {
	results := make([]FailedRunItem, len(failedRuns))
	sem := make(chan struct{}, resolveFetchConcurrency(len(failedRuns)))
	var wg sync.WaitGroup

	for i, fr := range failedRuns {
		wg.Add(1)
		go dispatchFetchRunWorker(&wg, sem, results, repo, fr, i)
	}
	wg.Wait()

	return results
}

func dispatchFetchRunWorker(wg *sync.WaitGroup, sem chan struct{}, results []FailedRunItem, repo string, fr ghRunItem, idx int) {
	defer wg.Done()
	sem <- struct{}{}
	results[idx] = fetchAndBuildFailedRunItem(repo, fr)
	<-sem
}

func resolveFetchConcurrency(total int) int {
	if total <= 8 {
		return total
	}

	return 8
}

func initFailedRunTopLevel(p *PipelineErrorLogsPayload, fr ghRunItem) {
	p.Conclusion = fr.Conclusion
	if len(p.Conclusion) == 0 {
		p.Conclusion = "failure"
	}
	p.WorkflowName = fr.Name
	p.RunId = fr.DatabaseId
	p.Url = fr.Url
	p.Branch = fr.HeadBranch
	p.Sha = fr.HeadSha
	p.CreatedAt = fr.CreatedAt
	p.UpdatedAt = fr.UpdatedAt
	p.DurationSeconds = calculateRunDuration(fr.CreatedAt, fr.UpdatedAt)
	p.SavedLogFile = filepath.ToSlash(getCachedPipelineLogPathForRepo(p.Repo, fr.DatabaseId))
}

func fetchAndBuildFailedRunItem(repo string, fr ghRunItem) FailedRunItem {
	rawLogs, ghJobs := fetchRunLogsAndJobsConcurrently(repo, fr.DatabaseId)
	jobs := CorrelateFailedJobsWithRunJobs(rawLogs, ghJobs)
	item := buildBaseFailedRunItem(repo, fr, rawLogs)
	item.FailedJobs = jobs
	item.StackTrace = extractStackTraceFromLog(rawLogs)

	return item
}

func fetchRunLogsAndJobsConcurrently(repo string, runId uint64) (string, []ghJobItem) {
	var rawLogs string
	var ghJobs []ghJobItem
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		rawLogs = queryFailedRunLogs(repo, runId)
	}()
	go func() {
		defer wg.Done()
		ghJobs = queryRunJobs(repo, runId)
	}()
	wg.Wait()

	return rawLogs, ghJobs
}

func buildBaseFailedRunItem(repo string, fr ghRunItem, rawLogs string) FailedRunItem {
	item := FailedRunItem{
		WorkflowName: fr.Name, RunId: fr.DatabaseId,
		Conclusion: fr.Conclusion, Status: fr.Status,
		Branch: fr.HeadBranch, Sha: fr.HeadSha,
		CreatedAt: fr.CreatedAt, UpdatedAt: fr.UpdatedAt,
		DurationSeconds: calculateRunDuration(fr.CreatedAt, fr.UpdatedAt),
		SavedLogFile:    filepath.ToSlash(getCachedPipelineLogPathForRepo(repo, fr.DatabaseId)),
		SavedMetaFile:   filepath.ToSlash(getCachedPipelineJSONPathForRepo(repo, fr.DatabaseId)),
		Url:             fr.Url, RawErrors: rawLogs,
	}

	return item
}
