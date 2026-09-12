package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// PipelineSyncResult summarizes an incremental sync of workflow runs and error logs.
type PipelineSyncResult struct {
	Repo           string   `json:"repo"`
	DbPath         string   `json:"dbPath"`
	RelativeDb     string   `json:"relativeDb"`
	ScannedRuns    int      `json:"scannedRuns"`
	NewRuns        int      `json:"newRuns"`
	DownloadedLogs int      `json:"downloadedLogs"`
	CachedErrors   []uint64 `json:"cachedErrors"`
	FailedRuns     []uint64 `json:"failedRuns"`
}

// FormatRelativeDbPath converts an absolute DB path to a repo-relative path starting with ./.
func FormatRelativeDbPath(fullPath string) string {
	if len(fullPath) == 0 {
		return "./.gitmap/pipeline_db/pipeline-default.db"
	}

	rel, err := filepath.Rel(resolveRepoRootDir(), fullPath)
	if err != nil || len(rel) == 0 {
		return filepath.ToSlash(fullPath)
	}

	slashRel := filepath.ToSlash(rel)
	if !strings.HasPrefix(slashRel, ".") {
		return "./" + slashRel
	}

	return slashRel
}

func resolveMaxSyncLimit(maxRuns int) int {
	if maxRuns <= 0 {
		return 20
	}

	if maxRuns > 20 {
		return 20
	}

	return maxRuns
}

func initSyncResult(repo, dbPath string, scanned int) *PipelineSyncResult {
	return &PipelineSyncResult{
		Repo:        repo,
		DbPath:      dbPath,
		RelativeDb:  FormatRelativeDbPath(dbPath),
		ScannedRuns: scanned,
	}
}

func queryWorkflowRunsLimit(repo string, limit int) []ghRunItem {
	runs := queryWorkflowRuns(repo)
	if len(runs) > limit {
		return runs[:limit]
	}

	return runs
}

func buildRunRecord(repo string, run ghRunItem) pipelinedb.PipelineRunRecord {
	isSuccess := run.Conclusion == "success"
	dur := calculateRunDuration(run.CreatedAt, run.UpdatedAt)

	return pipelinedb.PipelineRunRecord{
		RunId: run.DatabaseId, RepoSlug: repo, WorkflowName: run.Name,
		Status: run.Status, Conclusion: run.Conclusion, Branch: run.HeadBranch,
		Sha: run.HeadSha, DurationSeconds: dur, RunUrl: run.Url,
		IsSuccess: isSuccess, CreatedAt: run.CreatedAt, UpdatedAt: run.UpdatedAt,
	}
}

func recordRunInSplitDb(db *pipelinedb.PipelineSplitDb, repo string, run ghRunItem) error {
	record := buildRunRecord(repo, run)
	if err := db.RecordRun(record); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record pipeline run %d: %v\n", run.DatabaseId, err)

		return err
	}

	return nil
}

func saveParsedFailedJobs(db *pipelinedb.PipelineSplitDb, repo string, run ghRunItem, jobs []FailedJobItem, rawLogs string) {
	for _, job := range jobs {
		rec := buildParsedFailedJobRecord(repo, run, job, rawLogs)
		if err := db.RecordErrorLog(rec); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not record pipeline error log for run %d: %v\n", run.DatabaseId, err)
		}
	}
}

func buildParsedFailedJobRecord(repo string, run ghRunItem, job FailedJobItem, rawLogs string) pipelinedb.PipelineErrorRecord {
	return pipelinedb.PipelineErrorRecord{
		RunId:        run.DatabaseId,
		RepoSlug:     repo,
		WorkflowName: run.Name,
		StepName:     job.StepName,
		ErrorText:    job.FailureSummary,
		RawLogs:      rawLogs,
		CreatedAt:    run.UpdatedAt,
	}
}

func saveFallbackErrorLog(db *pipelinedb.PipelineSplitDb, repo string, run ghRunItem, rawLogs string) {
	rec := pipelinedb.PipelineErrorRecord{
		RunId:        run.DatabaseId,
		RepoSlug:     repo,
		WorkflowName: run.Name,
		StepName:     "Execution Failure",
		ErrorText:    rawLogs,
		RawLogs:      rawLogs,
		CreatedAt:    run.UpdatedAt,
	}

	if err := db.RecordErrorLog(rec); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record fallback error log for run %d: %v\n", run.DatabaseId, err)
	}
}

func saveJobsOrFallback(db *pipelinedb.PipelineSplitDb, repo string, run ghRunItem, jobs []FailedJobItem, rawLogs string) {
	hasJobs := len(jobs) > 0
	if hasJobs {
		saveParsedFailedJobs(db, repo, run, jobs, rawLogs)

		return
	}

	saveFallbackErrorLog(db, repo, run, rawLogs)
}

func fetchAndStoreRunErrorLog(db *pipelinedb.PipelineSplitDb, repo string, run ghRunItem, res *PipelineSyncResult) {
	rawLogs := queryFailedRunLogs(repo, run.DatabaseId)
	if len(rawLogs) == 0 {
		return
	}

	jobs := ParseFailedLogLines(rawLogs)
	saveJobsOrFallback(db, repo, run, jobs, rawLogs)
	res.DownloadedLogs++
	res.CachedErrors = append(res.CachedErrors, run.DatabaseId)
}

func handleSyncFailureLog(db *pipelinedb.PipelineSplitDb, repo string, run ghRunItem, cachedMap map[uint64]bool, res *PipelineSyncResult) {
	res.FailedRuns = append(res.FailedRuns, run.DatabaseId)
	hasCached := cachedMap[run.DatabaseId] || db.HasErrorLog(run.DatabaseId)
	if hasCached {
		res.CachedErrors = append(res.CachedErrors, run.DatabaseId)

		return
	}

	fetchAndStoreRunErrorLog(db, repo, run, res)
}

func processSyncRun(db *pipelinedb.PipelineSplitDb, repo string, run ghRunItem, cachedMap map[uint64]bool, res *PipelineSyncResult) {
	if err := recordRunInSplitDb(db, repo, run); err != nil {
		return
	}

	res.NewRuns++
	isFailure := run.Conclusion == "failure"
	if isFailure {
		handleSyncFailureLog(db, repo, run, cachedMap, res)
	}
}

func syncAllRunsIntoDb(db *pipelinedb.PipelineSplitDb, repo string, runs []ghRunItem, res *PipelineSyncResult) {
	cachedMap, err := db.QueryCachedErrorRunIdMap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not query cached error run ID map for %s: %v\n", repo, err)
	}

	for _, run := range runs {
		processSyncRun(db, repo, run, cachedMap, res)
	}
}

// SyncPipelineCache incrementally syncs up to maxRuns (capped at 20) pipeline runs and logs into SQLite.
func SyncPipelineCache(repo string, maxRuns int) (*PipelineSyncResult, error) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	limit := resolveMaxSyncLimit(maxRuns)
	runs := queryWorkflowRunsLimit(repo, limit)
	result := initSyncResult(repo, db.Path, len(runs))
	syncAllRunsIntoDb(db, repo, runs, result)

	return result, nil
}

func calculateAlreadyCachedCount(res *PipelineSyncResult) int {
	cachedCount := len(res.CachedErrors) - res.DownloadedLogs
	if cachedCount < 0 {
		return 0
	}

	return cachedCount
}

func renderSyncResultTerminal(res *PipelineSyncResult) {
	fmt.Printf("\n  %s● Incremental Pipeline Cache Sync (%s):%s\n",
		constants.ColorCyan, res.Repo, constants.ColorReset)
	fmt.Printf("    • Pipeline Database: %s\n", res.RelativeDb)
	fmt.Printf("    • Runs Scanned:      %d (max 20)\n", res.ScannedRuns)
	fmt.Printf("    • Failed Runs Found: %d\n", len(res.FailedRuns))
	fmt.Printf("    • New Logs Cached:   %s%d%s\n",
		constants.ColorGreen, res.DownloadedLogs, constants.ColorReset)
	cached := calculateAlreadyCachedCount(res)
	fmt.Printf("    • Already Cached:    %d (skipped download)\n\n", cached)
}

// HandlePipelineLastFailedLogs executes an incremental cache sync of the last 20 pipeline runs.
func HandlePipelineLastFailedLogs(args []string) error {
	repo := resolveCurrentRepoSlug()
	isJSON := hasArgFlag(args, "--json")
	res, err := SyncPipelineCache(repo, 20)
	if err != nil {
		return err
	}

	if isJSON {
		return printJSON(res)
	}

	renderSyncResultTerminal(res)

	return nil
}

// RunPipelineSyncCache is an alias for HandlePipelineLastFailedLogs.
func RunPipelineSyncCache(args []string) error {
	return HandlePipelineLastFailedLogs(args)
}
