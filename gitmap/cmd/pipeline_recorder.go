package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func recordPipelineInDB(p PipelineStatusPayload, runs []ghRunItem) {
	recordInPipelineSplitDB(p, runs)
	recordInMasterDB(p, runs)
}

func recordInPipelineSplitDB(p PipelineStatusPayload, runs []ghRunItem) {
	pipeDb, err := pipelinedb.OpenPipelineSplitDB(p.Repo)
	if err != nil {
		return
	}
	defer pipeDb.Close()

	for _, r := range runs {
		_ = pipeDb.RecordRun(pipelinedb.PipelineRunRecord{
			RunID:        r.DatabaseId,
			RepoSlug:     p.Repo,
			WorkflowName: r.Name,
			Status:       r.Status,
			Conclusion:   r.Conclusion,
			Branch:       r.HeadBranch,
			Sha:          r.HeadSha,
			EtaSeconds:   p.EtaSeconds,
			RunURL:       r.URL,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
		recordSingleFailedRun(pipeDb, p.Repo, r)
	}
}

func recordSingleFailedRun(pipeDb *pipelinedb.PipelineSplitDB, repo string, r ghRunItem) {
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
	_ = pipeDb.RecordErrorLog(pipelinedb.PipelineErrorRecord{
		RunID:        r.DatabaseId,
		RepoSlug:     repo,
		WorkflowName: r.Name,
		StepName:     "Failed Step",
		ErrorText:    clean,
		RawLogs:      raw,
	})
}

func recordInMasterDB(p PipelineStatusPayload, runs []ghRunItem) {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	for _, r := range runs {
		_ = db.InsertOrUpdatePipelineRun(store.PipelineRun{
			RunID:        r.DatabaseId,
			Repo:         p.Repo,
			WorkflowName: r.Name,
			Status:       r.Status,
			Conclusion:   r.Conclusion,
			Branch:       r.HeadBranch,
			Sha:          r.HeadSha,
			EtaSeconds:   p.EtaSeconds,
			URL:          r.URL,
		})
	}
}
