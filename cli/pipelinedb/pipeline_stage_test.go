package pipelinedb

import (
	"path/filepath"
	"testing"
)

func TestPipelineJobsAndStages(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "sql.db")
	conn, err := openCleanConn(dbPath, "test-repo")
	if err != nil {
		t.Fatalf("failed to open test conn: %v", err)
	}
	defer conn.Close()

	p := &PipelineSplitDb{conn: conn, RepoSlug: "test-repo", Path: dbPath}
	if schemaErr := p.setupSchema(); schemaErr != nil {
		t.Fatalf("setupSchema failed: %v", schemaErr)
	}

	recordTestRun(t, p, 101)
	recordTestJobs(t, p, 101)

	queriedJobs, qErr := p.QueryPipelineJobs(101)
	if qErr != nil {
		t.Fatalf("QueryPipelineJobs failed: %v", qErr)
	}
	if len(queriedJobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(queriedJobs))
	}

	summary, sumErr := p.QueryStageSummary(101)
	if sumErr != nil {
		t.Fatalf("QueryStageSummary failed: %v", sumErr)
	}
	if summary.StageSumSeconds != 190 {
		t.Errorf("expected StageSumSeconds 190, got %d", summary.StageSumSeconds)
	}
	if summary.WallClockSeconds != 120 {
		t.Errorf("expected WallClockSeconds 120, got %d", summary.WallClockSeconds)
	}
}

func recordTestRun(t *testing.T, p *PipelineSplitDb, runId uint64) {
	runRec := PipelineRunRecord{
		RunId:           runId,
		RepoSlug:        "test-repo",
		WorkflowName:    "CI",
		Status:          "completed",
		Conclusion:      "success",
		Branch:          "main",
		Sha:             "abcdef123456",
		DurationSeconds: 120,
		RunUrl:          "https://github.com/test-repo/runs/101",
		IsSuccess:       true,
		CreatedAt:       "2026-09-25T00:00:00Z",
		UpdatedAt:       "2026-09-25T00:02:00Z",
	}
	if recErr := p.RecordRun(runRec); recErr != nil {
		t.Fatalf("RecordRun failed: %v", recErr)
	}
}

func recordTestJobs(t *testing.T, p *PipelineSplitDb, runId uint64) {
	jobs := []PipelineJobRecord{
		{
			JobId:           1,
			RunId:           runId,
			RepoSlug:        "test-repo",
			JobName:         "Build Frontend",
			Status:          "completed",
			Conclusion:      "success",
			StartedAt:       "2026-09-25T00:00:10Z",
			CompletedAt:     "2026-09-25T00:01:30Z",
			DurationSeconds: 80,
			JobUrl:          "https://github.com/test-repo/jobs/1",
		},
		{
			JobId:           2,
			RunId:           runId,
			RepoSlug:        "test-repo",
			JobName:         "Check Rust Code",
			Status:          "completed",
			Conclusion:      "success",
			StartedAt:       "2026-09-25T00:00:10Z",
			CompletedAt:     "2026-09-25T00:02:00Z",
			DurationSeconds: 110,
			JobUrl:          "https://github.com/test-repo/jobs/2",
		},
	}
	if jErr := p.RecordPipelineJobs(jobs); jErr != nil {
		t.Fatalf("RecordPipelineJobs failed: %v", jErr)
	}
}
