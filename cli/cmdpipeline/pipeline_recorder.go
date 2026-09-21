package cmdpipeline

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func recordPipelineInDB(p PipelineStatusPayload, runs []ghRunItem) {
	recordInPipelineSplitDb(p, runs)
	recordInPipelineTasksDb(p, runs)
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

func recordInPipelineTasksDb(p PipelineStatusPayload, runs []ghRunItem) {
	pdb, err := store.OpenSectionTasksDB("pipeline", "")
	hasErr := err != nil
	if hasErr {
		return
	}
	defer pdb.Close()

	for _, r := range runs {
		insertSinglePipelineTask(pdb.Conn(), p.Repo, r)
	}
}

func insertSinglePipelineTask(conn *sql.DB, repo string, r ghRunItem) {
	taskId := fmt.Sprintf("pipe_%s_%d", repo, r.DatabaseId)
	action := "pipeline_run_" + r.Status
	target := r.HeadBranch
	fwd := fmt.Sprintf(`{"runId":%d,"workflow":"%s","conclusion":"%s","sha":"%s"}`, r.DatabaseId, r.Name, r.Conclusion, r.HeadSha)
	sqlQuery := `INSERT OR IGNORE INTO TaskHistory 
		(TaskId, Section, Action, Target, ForwardPayload, InversePayload, Status) 
		VALUES (?, 'pipeline', ?, ?, ?, '', 'completed')`
	_ = store.ExecWrapper(conn, sqlQuery, taskId, action, target, fwd)
}

func safeUint64ToInt64(val uint64) int64 {
	if val > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(val)
}
