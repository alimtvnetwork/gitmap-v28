package cmdpipeline

import (
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
	isInProgress    bool
	dataStatus      string
}

func executeAllPipelineErrorLogs(flags PipelineErrorFlags, args []string) error {
	workers := 0
	if flags.HasWorkers {
		workers = flags.Workers
	}
	summary := collectAllPipelineSummaryWithOptions(workers, flags.HasForce)

	if flags.IsJSON && flags.FilePath == "" {
		return emitAllPipelineSummaryJSON(summary)
	}

	if err := maybeWriteReportFile(flags, summary); err != nil {
		return err
	}

	if flags.IsJSON {
		return emitAllPipelineSummaryJSON(summary)
	}

	renderAllPipelineSummaryTerminal(summary)

	return nil
}

func maybeWriteReportFile(flags PipelineErrorFlags, summary AllPipelineSummary) error {
	if flags.FilePath == "" {
		return nil
	}

	return writeAllPipelineErrorsReportFile(flags, summary)
}

func emitAllPipelineSummaryJSON(summary AllPipelineSummary) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(summary)
}

func collectAllPipelineSummary() AllPipelineSummary {
	return collectAllPipelineSummaryWithOptions(0, false)
}

func collectAllPipelineSummaryWithOptions(workers int, bypassCache bool) AllPipelineSummary {
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
	inspections := runPeAllWorkerPool(slugs, absPaths, poolSize, bypassCache)

	for i, slug := range slugs {
		inspection := inspections[i]
		summary.repoDataStatus[slug] = inspection.dataStatus

		if inspection.isCached {
			summary.CachedCount++
			summary.CachedRepos = append(summary.CachedRepos, slug)
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
	if requested > 0 && repoCount > 0 && requested > repoCount {
		return repoCount
	}
	if requested > 0 {
		return requested
	}
	n := runtime.NumCPU()
	if n < 3 {
		n = 3
	}
	if repoCount > 0 && n > repoCount {
		return repoCount
	}

	return n
}

// runPeAllWorkerPool fans the per-repo inspections out to a semaphore-
// bounded worker pool. Each worker processes its repos sequentially.
func runPeAllWorkerPool(slugs []string, absPaths map[string]string, poolSize int, bypassCache bool) []singleRepoInspection {
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
			inspections[idx] = inspectSingleRepoCached(s, absPaths[s], cacheConn, &cacheMu, bypassCache)
		}(i, slug)
	}
	wg.Wait()

	return inspections
}

func renderAllPipelineSummaryTerminal(s AllPipelineSummary) {
	fmt.Printf("\n  %s● Pipeline Error Inspector: ALL REPOSITORIES%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    Repos Scanned:              %d\n", s.ReposScanned)
	fmt.Printf("    Passing / Clean:            %d\n", s.CleanCount)
	fmt.Printf("    Failing:                    %d\n", s.FailedCount)
	fmt.Printf("    No pipeline data:           %d\n\n", s.NoDataCount)

	switch {
	case s.FailedCount == 0 && s.NoDataCount == 0:
		fmt.Printf("  %s✔ All tracked CI/CD pipelines across all repositories are PASSING (100%% green).%s\n\n",
			constants.ColorGreen, constants.ColorReset)
	case s.FailedCount == 0:
		fmt.Printf("  %s⚠ No failing pipelines detected, but %d repositories had no pipeline data (status unknown, NOT verified green).%s\n\n",
			constants.ColorYellow, s.NoDataCount, constants.ColorReset)
	default:
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
