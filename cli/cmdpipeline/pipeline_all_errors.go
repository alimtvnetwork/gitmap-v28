package cmdpipeline

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// Pipeline repo data states for the trailing report table.
const (
	allRepoStatusClean         = "clean"
	allRepoStatusFailed        = "failed"
	allRepoStatusNoData        = "no-pipeline-data"
	allRepoStatusDbUnreachable = "db-unavailable"
)

// repoCatalogEntry is one catalog repository: sanitized pipeline slug plus
// its absolute path. Kept as an ordered slice (not a slug-keyed map) so
// two catalog entries sharing a slug are never silently dropped from reports.
type repoCatalogEntry struct {
	slug    string
	absPath string
}

// AllPipelineSummary aggregates pipeline error statuses across all repositories.
type AllPipelineSummary struct {
	TotalPipelines    int                     `json:"totalPipelines"`
	ReposScanned      int                     `json:"reposScanned"`
	CleanCount        int                     `json:"cleanCount"`
	FailedCount       int                     `json:"failedCount"`
	NoDataCount       int                     `json:"noDataCount"`
	CachedCount       int                     `json:"cachedCount"`
	TotalErrorEntries int                     `json:"totalErrorEntries"`
	FailedItems       []AllPipelineFailedItem `json:"failedItems,omitempty"`
	CleanRepos        []string                `json:"cleanRepos,omitempty"`
	NoDataRepos       []string                `json:"noDataRepos,omitempty"`
	CachedRepos       []string                `json:"cachedRepos,omitempty"`
	PullErrors        []store.PullErrorRecord `json:"pullErrors,omitempty"`

	repoAbsPaths   map[string]string
	repoDataStatus map[string]string
	errorDetails   map[string][]pipelinedb.PipelineErrorRecord
	stepJobs       map[string]map[string]string
	// repoList preserves every catalog entry in order (duplicates kept);
	// the slug-keyed maps above collapse duplicate slugs.
	repoList []repoCatalogEntry
}

// AllPipelineFailedItem represents a failed pipeline run discovered across repositories.
type AllPipelineFailedItem struct {
	RepoSlug     string `json:"repoSlug"`
	WorkflowName string `json:"workflowName"`
	RunId        uint64 `json:"runId"`
	Branch       string `json:"branch"`
	Sha          string `json:"sha"`
	ErrorSummary string `json:"errorSummary"`
	UpdatedAt    string `json:"updatedAt"`
}

// singleRepoInspection carries the tri-state inspection result for one catalog repo.
type singleRepoInspection struct {
	failedItem      *AllPipelineFailedItem
	errorRecords    []pipelinedb.PipelineErrorRecord
	stepJobs        map[string]string
	isClean         bool
	hasPipelineData bool
	isCached        bool
	dataStatus      string
}

func executeAllPipelineErrorLogs(flags PipelineErrorFlags, args []string) error {
	workers := 0
	if flags.HasWorkers {
		workers = flags.Workers
	}
	summary := collectAllPipelineSummaryWithWorkers(workers)

	if flags.IsJSON && flags.FilePath == "" {
		return emitAllPipelineSummaryJSON(summary)
	}

	if flags.FilePath != "" {
		if err := writeAllPipelineErrorsReportFile(flags, summary); err != nil {
			return err
		}
	}

	if flags.IsJSON {
		return emitAllPipelineSummaryJSON(summary)
	}

	renderAllPipelineSummaryTerminal(summary)

	return nil
}

func emitAllPipelineSummaryJSON(summary AllPipelineSummary) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(summary)
}

func collectAllPipelineSummary() AllPipelineSummary {
	return collectAllPipelineSummaryWithWorkers(0)
}

// collectAllPipelineSummaryWithWorkers inspects every catalog repo using a
// worker pool. workers <= 0 means auto: max(runtime.NumCPU(), 3).
// Each worker runs the per-repo inspection SEQUENTIALLY — no nested
// parallelism inside a worker. Results are aggregated in catalog order,
// so output is deterministic regardless of worker scheduling.
func collectAllPipelineSummaryWithWorkers(workers int) AllPipelineSummary {
	entries := discoverCatalogRepoSlugs()
	slugs := make([]string, 0, len(entries))
	absPaths := make(map[string]string, len(entries))
	for _, e := range entries {
		slugs = append(slugs, e.slug)
		absPaths[e.slug] = e.absPath
	}
	summary := AllPipelineSummary{
		TotalPipelines: len(slugs),
		ReposScanned:   len(slugs),
		repoAbsPaths:   absPaths,
		repoList:       entries,
		repoDataStatus: make(map[string]string, len(slugs)),
		errorDetails:   make(map[string][]pipelinedb.PipelineErrorRecord),
		stepJobs:       make(map[string]map[string]string),
	}

	poolSize := resolvePeAllWorkers(workers, len(slugs))
	inspections := runPeAllWorkerPool(slugs, absPaths, poolSize)

	for i, slug := range slugs {
		inspection := inspections[i]
		summary.repoDataStatus[slug] = inspection.dataStatus

		if inspection.isCached {
			summary.CachedCount++
			summary.CachedRepos = append(summary.CachedRepos, slug)
			continue
		}

		if !inspection.hasPipelineData {
			summary.NoDataCount++
			summary.NoDataRepos = append(summary.NoDataRepos, slug)
			continue
		}

		if inspection.isClean {
			summary.CleanCount++
			summary.CleanRepos = append(summary.CleanRepos, slug)
			continue
		}

		if inspection.failedItem == nil {
			continue
		}

		summary.FailedCount++
		summary.FailedItems = append(summary.FailedItems, *inspection.failedItem)
		summary.TotalErrorEntries += len(inspection.errorRecords)
		summary.errorDetails[slug] = inspection.errorRecords
		summary.stepJobs[slug] = inspection.stepJobs
	}

	sort.Slice(summary.FailedItems, func(i, j int) bool {
		return summary.FailedItems[i].RepoSlug < summary.FailedItems[j].RepoSlug
	})

	summary.PullErrors = queryPullErrorsForAll()

	return summary
}

// resolvePeAllWorkers returns the worker count: explicit override wins,
// otherwise max(CPU threads, 3), capped at the repo count.
func resolvePeAllWorkers(requested, repoCount int) int {
	if requested > 0 {
		if requested > repoCount && repoCount > 0 {
			return repoCount
		}

		return requested
	}
	n := runtime.NumCPU()
	if n < 3 {
		n = 3
	}
	if n > repoCount && repoCount > 0 {
		return repoCount
	}

	return n
}

// runPeAllWorkerPool fans the per-repo inspections out to a semaphore-
// bounded worker pool. Each worker processes its repos sequentially.
func runPeAllWorkerPool(slugs []string, absPaths map[string]string, poolSize int) []singleRepoInspection {
	inspections := make([]singleRepoInspection, len(slugs))
	if len(slugs) == 0 {
		return inspections
	}

	cacheConn, cacheErr := openPeAllCache()
	if cacheErr != nil {
		cacheConn = nil
	}
	if cacheConn != nil {
		defer cacheConn.Close()
	}
	var cacheMu sync.Mutex

	sem := make(chan struct{}, poolSize)
	var wg sync.WaitGroup
	for i, slug := range slugs {
		wg.Add(1)
		go func(idx int, s string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			inspections[idx] = inspectSingleRepoCached(s, absPaths[s], cacheConn, &cacheMu)
		}(i, slug)
	}
	wg.Wait()

	return inspections
}

// inspectSingleRepoCached checks the commit-hash cache before inspecting.
// Cache hit (HEAD SHA already recorded): skip the inspection, mark done.
// Cache miss: run the full inspection, then record the SHA.
func inspectSingleRepoCached(slug, absPath string, cacheConn *sql.DB, cacheMu *sync.Mutex) singleRepoInspection {
	headSha := getRepoHeadSha(absPath)
	if headSha != "" && cacheConn != nil {
		cacheMu.Lock()
		cachedSha := getCachedCommitSha(cacheConn, slug)
		cacheMu.Unlock()
		if cachedSha == headSha {
			return singleRepoInspection{isCached: true, dataStatus: "cached-unchanged"}
		}
	}

	inspection := inspectSinglePipelineRepo(slug)

	if headSha != "" && cacheConn != nil {
		cacheMu.Lock()
		if cacheErr := recordCommitSha(cacheConn, slug, headSha); cacheErr != nil {
			fmt.Fprintf(os.Stderr, "warning: pe-all cache write failed for %s: %v\n", slug, cacheErr)
		}
		cacheMu.Unlock()
	}

	return inspection
}

// discoverCatalogRepoSlugs enumerates every repository in the catalog as
// discoverCatalogRepoSlugs enumerates every repository in the catalog as an
// ordered list of (sanitized slug, absolute path) entries. Duplicate slugs
// are preserved as separate entries and never collapsed.
func discoverCatalogRepoSlugs() []repoCatalogEntry {
	entries := []repoCatalogEntry{}

	db, err := store.OpenDefault()

	if err != nil {
		return entries
	}
	defer db.Close()

	records, err := db.ListRepos()

	if err != nil {
		return entries
	}

	for _, rec := range records {
		entries = append(entries, repoCatalogEntry{
			slug:    pipelinedb.SanitizeRepoSlug(rec.Slug),
			absPath: rec.AbsolutePath,
		})
	}

	return entries
}

// inspectSinglePipelineRepo inspects one repo's pipeline DB without creating it:
// the DB path is stat'ed BEFORE OpenPipelineSplitDb, which initializes missing
// DBs as a side effect. Repos without a DB file report no pipeline data.
func inspectSinglePipelineRepo(slug string) singleRepoInspection {
	dbPath := pipelinedb.ResolvePipelineDbPath(slug)

	if _, statErr := os.Stat(dbPath); statErr != nil {
		return singleRepoInspection{hasPipelineData: false, dataStatus: allRepoStatusNoData}
	}

	db, err := pipelinedb.OpenPipelineSplitDb(slug)

	if err != nil {
		return singleRepoInspection{hasPipelineData: false, dataStatus: allRepoStatusDbUnreachable}
	}
	defer db.Close()

	run, runErr := db.QueryRunByNegativeOffset(-1)

	if runErr != nil || run == nil {
		return singleRepoInspection{hasPipelineData: false, dataStatus: allRepoStatusNoData}
	}

	if run.IsSuccess || strings.EqualFold(run.Conclusion, "success") {
		return singleRepoInspection{hasPipelineData: true, isClean: true, dataStatus: allRepoStatusClean}
	}

	errSummary := extractRunErrorSummary(db, run.RunId)

	return singleRepoInspection{
		hasPipelineData: true,
		dataStatus:      allRepoStatusFailed,
		errorRecords:    detailedErrorRecordsForRun(db, run.RunId),
		stepJobs:        stepJobNamesForRun(db, run.RunId),
		failedItem: &AllPipelineFailedItem{
			RepoSlug:     slug,
			WorkflowName: run.WorkflowName,
			RunId:        run.RunId,
			Branch:       run.Branch,
			Sha:          run.Sha,
			ErrorSummary: errSummary,
			UpdatedAt:    run.UpdatedAt,
		},
	}
}

func extractRunErrorSummary(db *pipelinedb.PipelineSplitDb, runId uint64) string {
	compactRes := db.QueryCompactErrorLogsByRunId(runId)

	if compactRes.IsSuccess() && len(compactRes.Data) > 0 {
		return compactRes.Data[0].ErrorText
	}

	return "Pipeline workflow run failed"
}

func queryPullErrorsForAll() []store.PullErrorRecord {
	pullDb, err := store.OpenPullSplitDB()

	if err != nil {
		return nil
	}
	defer pullDb.Close()

	records, _ := pullDb.QueryLatestPullErrors("all", 10)

	return records
}

func renderAllPipelineSummaryTerminal(s AllPipelineSummary) {
	fmt.Printf("\n  %s● Pipeline Error Inspector: ALL REPOSITORIES%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    Repos Scanned:              %d\n", s.ReposScanned)
	fmt.Printf("    Passing / Clean:            %d\n", s.CleanCount)
	fmt.Printf("    Failing:                    %d\n", s.FailedCount)
	fmt.Printf("    No pipeline data:           %d\n\n", s.NoDataCount)

	if s.FailedCount == 0 {
		fmt.Printf("  %s✔ All tracked CI/CD pipelines across all repositories are PASSING (100%% green).%s\n\n",
			constants.ColorGreen, constants.ColorReset)
	}

	if s.FailedCount > 0 {
		renderFailingPipelinesList(s)
	}

	renderNoPipelineDataList(s.NoDataRepos)

	if len(s.PullErrors) > 0 {
		renderPullErrorsNotice(len(s.PullErrors))
	}
}

func renderFailingPipelinesList(s AllPipelineSummary) {
	fmt.Printf("  %sFailing Repositories (%d):%s\n", constants.ColorRed, len(s.FailedItems), constants.ColorReset)

	for _, item := range s.FailedItems {
		errorCount := len(s.errorDetails[item.RepoSlug])
		fmt.Printf("    Repo: %s%s%s · Errors: %d\n", constants.ColorYellow, item.RepoSlug, constants.ColorReset, errorCount)
		fmt.Printf("      └── ✖ %s (Run #%d, %s, %s, %s)\n", item.WorkflowName, item.RunId, item.Branch, item.Sha, item.UpdatedAt)
		fmt.Printf("          Error: %s\n", item.ErrorSummary)
		fmt.Printf("          ↳ Inspect: gitmap pe %s\n\n", item.RepoSlug)
	}
}

func renderNoPipelineDataList(repos []string) {
	if len(repos) == 0 {
		return
	}

	fmt.Printf("  No pipeline data (%d): %s\n\n", len(repos), truncateSlugDisplay(repos))
}

func truncateSlugDisplay(repos []string) string {
	const maxShown = 12

	if len(repos) <= maxShown {
		return strings.Join(repos, ", ")
	}

	return strings.Join(repos[:maxShown], ", ") + fmt.Sprintf(", … and %d more", len(repos)-maxShown)
}

func renderPullErrorsNotice(count int) {
	fmt.Printf("  %s● Recorded Pull Errors (%d failure(s)):%s\n", constants.ColorYellow, count, constants.ColorReset)
	fmt.Printf("    ↳ To inspect pull stack traces: gitmap pull-error all (or: gitmap pulle all)\n\n")
}
