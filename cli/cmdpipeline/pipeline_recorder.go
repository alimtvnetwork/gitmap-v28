package cmdpipeline

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func recordPipelineInDB(p PipelineStatusPayload, runs []ghRunItem) {
	recordInPipelineSplitDb(p, runs)
	recordInMasterDB(p, runs)
}

func recordInPipelineSplitDb(p PipelineStatusPayload, runs []ghRunItem) {
	pipeDb, err := pipelinedb.OpenPipelineSplitDb(p.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not open pipeline split DB for %s: %v\n", p.Repo, err)

		return
	}

	defer pipeDb.Close()

	recordRunsToSplitDb(pipeDb, p, runs)
	_, _, _ = pipeDb.PruneIfExceedsSize(0)
	touchDbTimestamp(pipeDb.Path)
	writeCacheSyncMeta(pipeDb.Path, runs)
}

// RecordFetchedRunsToSplitDb persists fresh workflow runs and failures into the repository SQLite split DB.
func RecordFetchedRunsToSplitDb(repo string, runs []ghRunItem) {
	if len(runs) == 0 {
		return
	}

	pipeDb, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return
	}
	defer pipeDb.Close()

	recordRunsToSplitDb(pipeDb, PipelineStatusPayload{Repo: repo}, runs)
	_, _, _ = pipeDb.PruneIfExceedsSize(0)
	touchDbTimestamp(pipeDb.Path)
	writeCacheSyncMeta(pipeDb.Path, runs)
}

func touchDbTimestamp(dbPath string) {
	now := time.Now()
	_ = os.Chtimes(dbPath, now, now)
}

func recordRunsToSplitDb(pipeDb *pipelinedb.PipelineSplitDb, p PipelineStatusPayload, runs []ghRunItem) {
	for _, r := range runs {
		recordSingleSplitRun(pipeDb, p, r)
	}
}

func recordSingleSplitRun(pipeDb *pipelinedb.PipelineSplitDb, p PipelineStatusPayload, r ghRunItem) {
	rec := buildPipelineRunRecord(p, r)
	if err := pipeDb.RecordRun(rec); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record pipeline run %d: %v\n", r.DatabaseId, err)
	}
}

func resolveRecordEta(status string, eta int) int {
	if status == "in_progress" || status == "queued" {
		return eta
	}

	return 0
}

func initBasicRunRecord(repo string, r ghRunItem) pipelinedb.PipelineRunRecord {
	return pipelinedb.PipelineRunRecord{
		RunId:        r.DatabaseId,
		RepoSlug:     repo,
		WorkflowName: r.Name,
		Status:       r.Status,
		Conclusion:   r.Conclusion,
		Branch:       r.HeadBranch,
		Sha:          r.HeadSha,
		RunUrl:       r.Url,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func buildPipelineRunRecord(p PipelineStatusPayload, r ghRunItem) pipelinedb.PipelineRunRecord {
	rec := initBasicRunRecord(p.Repo, r)
	rec.EtaSeconds = resolveRecordEta(r.Status, p.EtaSeconds)
	rec.DurationSeconds = calculateRunDuration(r.CreatedAt, r.UpdatedAt)
	rec.IsSuccess = r.Conclusion == "success"

	return rec
}

func persistSingleFailedRunLog(pipeDb *pipelinedb.PipelineSplitDb, repo string, r ghRunItem, clean, raw string) {
	compactClean, filteredCount := FilterCompactLogText(clean)
	compactRaw, _ := FilterCompactLogText(raw)
	detailRec := buildDetailRecord(repo, r, clean, raw)
	compactRec := buildCompactRecord(repo, r, compactClean, compactRaw, filteredCount)

	if err := pipeDb.RecordDualErrorLog(detailRec, compactRec); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record pipeline error log for run %d: %v\n", r.DatabaseId, err)
	}
}

func buildDetailRecord(repo string, r ghRunItem, clean, raw string) pipelinedb.PipelineErrorRecord {
	return pipelinedb.PipelineErrorRecord{
		RunId:        r.DatabaseId,
		RepoSlug:     repo,
		WorkflowName: r.Name,
		StepName:     "Failed Step",
		ErrorText:    clean,
		RawLogs:      raw,
	}
}

func buildCompactRecord(repo string, r ghRunItem, clean, raw string, filtered int) pipelinedb.PipelineCompactErrorRecord {
	return pipelinedb.PipelineCompactErrorRecord{
		RunId:           r.DatabaseId,
		RepoSlug:        repo,
		WorkflowName:    r.Name,
		StepName:        "Failed Step",
		ErrorText:       clean,
		CompactLogs:     raw,
		FilteredOkCount: filtered,
	}
}

func recordInMasterDB(p PipelineStatusPayload, runs []ghRunItem) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not open master DB: %v\n", err)

		return
	}

	defer db.Close()

	insertMasterRuns(db, p, runs)
}

func insertMasterRuns(db *store.DB, p PipelineStatusPayload, runs []ghRunItem) {
	for _, r := range runs {
		insertSingleMasterRun(db, p, r)
	}
}

func insertSingleMasterRun(db *store.DB, p PipelineStatusPayload, r ghRunItem) {
	run := store.PipelineRun{
		RunID:        safeUint64ToInt64(r.DatabaseId),
		Repo:         p.Repo,
		WorkflowName: r.Name,
		Status:       r.Status,
		Conclusion:   r.Conclusion,
		Branch:       r.HeadBranch,
		Sha:          r.HeadSha,
		EtaSeconds:   resolveRecordEta(r.Status, p.EtaSeconds),
		URL:          r.Url,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}

	if err := db.InsertOrUpdatePipelineRun(run); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record master pipeline run %d: %v\n", r.DatabaseId, err)
	}
}

func safeUint64ToInt64(val uint64) int64 {
	if val > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(val)
}
