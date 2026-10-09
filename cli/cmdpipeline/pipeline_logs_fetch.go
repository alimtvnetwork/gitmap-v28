package cmdpipeline

import (
	"fmt"
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

func printPipelineErrorLogsHelp() {
	printPipelineErrorLogsUsage()
	printPipelineErrorLogsFlags()
}

func printPipelineErrorLogsUsage() {
	fmt.Println("Usage: gitmap pipeline error-logs [repo] [commit|-N|-Nn|HEAD~N] [flags]")
	fmt.Println("       gitmap pipeline errors [repo] [commit|-N] [clear [-y]] [flags]")
	fmt.Println("       gitmap pe [repo] [commit|-N|-Nn|HEAD~N] [clear [-y]] [flags]")
	fmt.Println("       gitmap pe all [--json] [--file <path>] [flags]   (alias: gitmap te all)")
	fmt.Println("       gitmap pe -f <format.json|alias> [flags]")
	fmt.Println("       gitmap pe -f <alias> -test <filepath>")
	fmt.Println("       gitmap pe -f <alias> -test-commit <commit-sha> [-repo <path>]")
	fmt.Println("       gitmap pe add-format <file.json> [alias]")
	fmt.Println("       gitmap pe remove-format (rm-format) <name|alias>")
	fmt.Println("       gitmap pe add-all <folder-path>")
	fmt.Println("       gitmap pe list-formats")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  clear [-y]                     Purge error logs, reports, and reset pipeline DB for current repo")
	fmt.Println("  all [--json] [--file <path>]   Aggregate pipeline error status across ALL catalog repos")
	fmt.Println("                                 (alias: te all; same as pe all; --all also works)")
	fmt.Println("  add-format <file.json> [alias] Register custom JSON format profile for error/warning filtering")
	fmt.Println("  remove-format <name|alias>     Remove registered format profile (alias: rm-format)")
	fmt.Println("  add-all <folder-path>          Batch register all format profiles (*.json) from folder")
	fmt.Println("  list-formats                   List all registered error log format profiles")
	fmt.Println("  preview-format <name|file>     Preview format configuration rules in JSON")
	fmt.Println()
	fmt.Println("Targeting:")
	fmt.Println("  [repo]                         Target repository by name, prefix, alias, or path")
	fmt.Println("  <commit-sha>                   Filter errors for specific commit (e.g. gitmap pe ee4a694)")
	fmt.Println("  -1, -2, -3, -1n, HEAD~1        Inspect errors for previous commits by relative offset")
	fmt.Println()
}

func printPipelineErrorLogsFlags() {
	fmt.Println("Flags:")
	printPEExecutionFlags()
	printPEOutputFlags()
}

func printPEExecutionFlags() {
	fmt.Println("  -f <format|file.json>          Apply custom format profile to filter and format errors (e.g. -f tauri)")
	fmt.Println("  -f <format> -test <file>       Test format profile against a local log file")
	fmt.Println("  -f <format> -test-commit <sha> Test format profile against remote GitHub Actions run for commit")
	fmt.Println("  -1, -2, -3, HEAD~N             Target past workflow run by relative commit offset or SHA")
	fmt.Println("  -t, --timeline                 Watch pipeline dynamic timeline until completion")
	fmt.Println("  -f, --fix                      Execute internal CI/CD diagnostic & auto-repair suite (when standalone)")
	fmt.Println("  -c, --check                    Run internal CI/CD checks without modifying files")
	fmt.Println("  -v, --detailed, --verbose      Show full raw error logs including passing ok lines")
	fmt.Println("  -y, --yes                      Auto-confirm prompts non-interactively")
	fmt.Println("  --force, --no-cache            Bypass local SQLite DB cache and pull fresh from GitHub")
}

func printPEOutputFlags() {
	fmt.Println("  --json                         Output data in structured JSON format")
	fmt.Println("  --file <path>                  Write error logs to specified file path")
	fmt.Println("  --tempfile <filename>          Write error logs to .ai-memory/temp/<filename>")
	fmt.Println("  -l, --limit, -n, --lines <N>   Limit displayed error log lines and failure run records (e.g. -l 10)")
	fmt.Println("  -n, --no-output-log            Stage error logs to disk without displaying in terminal")
}

func printPipelineLogsHelp() {
	fmt.Println("Usage: gitmap pipeline logs [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --json                  Output workflow status and URL in JSON format")
}
