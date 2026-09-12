package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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
}

func recordRunsToSplitDb(pipeDb *pipelinedb.PipelineSplitDb, p PipelineStatusPayload, runs []ghRunItem) {
	for _, r := range runs {
		recordSingleSplitRun(pipeDb, p, r)
		recordSingleFailedRun(pipeDb, p.Repo, r)
	}
}

func recordSingleSplitRun(pipeDb *pipelinedb.PipelineSplitDb, p PipelineStatusPayload, r ghRunItem) {
	rec := buildPipelineRunRecord(p, r)
	if err := pipeDb.RecordRun(rec); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record pipeline run %d: %v\n", r.DatabaseId, err)
	}
}

func buildPipelineRunRecord(p PipelineStatusPayload, r ghRunItem) pipelinedb.PipelineRunRecord {
	return pipelinedb.PipelineRunRecord{
		RunId:        r.DatabaseId,
		RepoSlug:     p.Repo,
		WorkflowName: r.Name,
		Status:       r.Status,
		Conclusion:   r.Conclusion,
		Branch:       r.HeadBranch,
		Sha:          r.HeadSha,
		EtaSeconds:   p.EtaSeconds,
		RunUrl:       r.Url,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func recordSingleFailedRun(pipeDb *pipelinedb.PipelineSplitDb, repo string, r ghRunItem) {
	if r.Conclusion != "failure" {
		return
	}
	if pipeDb.HasErrorLog(r.DatabaseId) {
		return
	}
	if isSkipDelayRequested() {
		return
	}
	raw := queryFailedRunLogs(repo, r.DatabaseId)
	clean := extractCleanErrorLines(raw)
	if clean == "" {
		return
	}
	persistSingleFailedRunLog(pipeDb, repo, r, clean, raw)
}

func persistSingleFailedRunLog(pipeDb *pipelinedb.PipelineSplitDb, repo string, r ghRunItem, clean, raw string) {
	err := pipeDb.RecordErrorLog(pipelinedb.PipelineErrorRecord{
		RunId:        r.DatabaseId,
		RepoSlug:     repo,
		WorkflowName: r.Name,
		StepName:     "Failed Step",
		ErrorText:    clean,
		RawLogs:      raw,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record pipeline error log for run %d: %v\n", r.DatabaseId, err)
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
		EtaSeconds:   p.EtaSeconds,
		URL:          r.Url,
	}
	if err := db.InsertOrUpdatePipelineRun(run); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record master pipeline run %d: %v\n", r.DatabaseId, err)
	}
}
