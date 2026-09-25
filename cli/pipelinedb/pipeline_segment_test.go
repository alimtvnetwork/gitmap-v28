package pipelinedb

import (
	"path/filepath"
	"testing"
)

func TestPipelineSegmentsAndStats(t *testing.T) {
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

	segs := []PipelineSegmentRecord{
		{
			RunId:           202,
			JobName:         "Build Frontend",
			StepName:        "Setup Node",
			StepNumber:      1,
			Status:          "completed",
			Conclusion:      "success",
			DurationSeconds: 15,
		},
		{
			RunId:           202,
			JobName:         "Build Frontend",
			StepName:        "Build Vite App",
			StepNumber:      2,
			Status:          "completed",
			Conclusion:      "success",
			DurationSeconds: 65,
		},
	}
	if sErr := p.RecordPipelineSegments(202, segs); sErr != nil {
		t.Fatalf("RecordPipelineSegments failed: %v", sErr)
	}

	queriedSegs, qsErr := p.QueryPipelineSegments(202)
	if qsErr != nil {
		t.Fatalf("QueryPipelineSegments failed: %v", qsErr)
	}
	if len(queriedSegs) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(queriedSegs))
	}

	info := p.GetDatabaseInfo()
	if info.SegmentCount != 2 {
		t.Errorf("expected SegmentCount 2, got %d", info.SegmentCount)
	}
}
