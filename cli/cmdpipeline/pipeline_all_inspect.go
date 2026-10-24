package cmdpipeline

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func queryInspectionFromCache(cacheConn *sql.DB, cacheMu *sync.Mutex, slug, headSha string) (singleRepoInspection, bool) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	return getCachedInspection(cacheConn, slug, headSha)
}

func saveInspectionToCache(cacheConn *sql.DB, cacheMu *sync.Mutex, slug, headSha string, inspection singleRepoInspection) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cacheErr := recordCachedInspection(cacheConn, slug, headSha, inspection); cacheErr != nil {
		fmt.Fprintf(os.Stderr, "warning: pe-all cache write failed for %s: %v\n", slug, cacheErr)
	}
}

func tryGetCachedInspection(cacheConn *sql.DB, cacheMu *sync.Mutex, slug, headSha string, bypassCache bool) (singleRepoInspection, bool) {
	if bypassCache || headSha == "" || cacheConn == nil {
		return singleRepoInspection{}, false
	}

	return queryInspectionFromCache(cacheConn, cacheMu, slug, headSha)
}

// inspectSingleRepoCached checks the commit-hash cache before inspecting.
// Cache hit (HEAD SHA already recorded): skip the inspection, restore cached result.
// Cache miss: run the full inspection, then record the SHA.
func inspectSingleRepoCached(slug, absPath string, cacheConn *sql.DB, cacheMu *sync.Mutex, bypassCache bool) singleRepoInspection {
	headSha := getRepoHeadSha(absPath)
	if cached, isHit := tryGetCachedInspection(cacheConn, cacheMu, slug, headSha, bypassCache); isHit {
		return cached
	}

	inspection := inspectSinglePipelineRepo(slug)

	if !inspection.isInProgress && headSha != "" && cacheConn != nil {
		saveInspectionToCache(cacheConn, cacheMu, slug, headSha, inspection)
	}

	return inspection
}

func resolveCatalogRepoSlug(rec model.ScanRecord) string {
	slug := pipelinedb.SanitizeRepoSlug(rec.Slug)
	urlCandidate := rec.HTTPSUrl
	if len(urlCandidate) == 0 {
		urlCandidate = rec.SSHUrl
	}
	if len(urlCandidate) == 0 {
		return slug
	}

	parsed := parseSlugFromGitURL(urlCandidate)
	if !strings.Contains(parsed, "/") {
		return slug
	}

	candSlug := pipelinedb.SanitizeRepoSlug(parsed)
	candDb := pipelinedb.PipelineDbPath(candSlug)
	origDb := pipelinedb.PipelineDbPath(slug)
	if getFileSizeBytes(candDb) > getFileSizeBytes(origDb) {
		return candSlug
	}

	return slug
}

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
			slug:    resolveCatalogRepoSlug(rec),
			absPath: rec.AbsolutePath,
		})
	}

	return entries
}

func getFileSizeBytes(p string) int64 {
	info, err := os.Stat(p)
	if err != nil {
		return 0
	}

	return info.Size()
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

	if run.IsSuccess || strings.EqualFold(run.Conclusion, "success") || strings.EqualFold(run.Conclusion, "skipped") || strings.EqualFold(run.Conclusion, "neutral") {
		return singleRepoInspection{hasPipelineData: true, isClean: true, dataStatus: allRepoStatusClean}
	}

	if isWorkflowStatusActive(run.Status) && !isFailingConclusion(run.Conclusion) && !db.HasErrorLogForSha(run.Sha) {
		return singleRepoInspection{
			hasPipelineData: true,
			isClean:         true,
			isInProgress:    true,
			dataStatus:      allRepoStatusClean,
		}
	}

	errRecords := detailedErrorRecordsForRun(db, run.RunId)
	if !isFailingConclusion(run.Conclusion) && !db.HasErrorLogForSha(run.Sha) && len(errRecords) == 0 {
		return singleRepoInspection{
			hasPipelineData: true,
			isClean:         true,
			isInProgress:    isWorkflowStatusActive(run.Status),
			dataStatus:      allRepoStatusClean,
		}
	}

	errSummary := extractRunErrorSummary(db, run.RunId)
	if errSummary == "" && isFailingConclusion(run.Conclusion) {
		errSummary = fmt.Sprintf("Workflow %s concluded with %s", run.WorkflowName, run.Conclusion)
	}

	return singleRepoInspection{
		hasPipelineData: true,
		isInProgress:    isWorkflowStatusActive(run.Status),
		dataStatus:      allRepoStatusFailed,
		errorRecords:    errRecords,
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

	recordsRes := pullDb.QueryLatestPullErrors("all", 10)

	return recordsRes.Items()
}
