package pipelinedb

import (
	"testing"
)

func TestPipelineSplitDBLifecycle(t *testing.T) {
	db, err := OpenPipelineSplitDB("test-owner/test-repo")
	if err != nil {
		t.Fatalf("failed to open pipeline split db: %v", err)
	}
	defer db.Close()

	// 1. Record Run
	run := PipelineRunRecord{
		RunId:        12345,
		RepoSlug:     "test-owner/test-repo",
		WorkflowName: "CI",
		Status:       "completed",
		Conclusion:   "failure",
		Branch:       "main",
		Sha:          "abc1234",
		EtaSeconds:   60,
		RunUrl:       "https://github.com/test-owner/test-repo/actions/runs/12345",
		CreatedAt:    "2026-09-03T10:00:00Z",
		UpdatedAt:    "2026-09-03T10:02:00Z",
	}
	if err := db.RecordRun(run); err != nil {
		t.Fatalf("failed to record run: %v", err)
	}

	// 2. Record Error Log
	errLog := PipelineErrorRecord{
		RunId:        12345,
		RepoSlug:     "test-owner/test-repo",
		WorkflowName: "CI",
		StepName:     "Test Step",
		ErrorText:    "##[error]Process exited with code 1",
	}
	if err := db.RecordErrorLog(errLog); err != nil {
		t.Fatalf("failed to record error log: %v", err)
	}

	// 3. Query Stats
	stats, err := db.GetStats()
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}
	if stats.TotalRuns != 1 {
		t.Errorf("expected 1 run, got %d", stats.TotalRuns)
	}
	if stats.ErrorLogCount != 1 {
		t.Errorf("expected 1 error log, got %d", stats.ErrorLogCount)
	}

	// 4. Query Error Logs
	logs, err := db.QueryRecentErrorLogs(5)
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected 1 error log, got %d (err: %v)", len(logs), err)
	}
	if logs[0].StepName != "Test Step" {
		t.Errorf("expected step name 'Test Step', got %s", logs[0].StepName)
	}

	// 5. Optimize
	if _, err := db.Optimize(); err != nil {
		t.Errorf("failed to optimize: %v", err)
	}

	// 6. Clear
	if err := db.Clear(); err != nil {
		t.Fatalf("failed to clear: %v", err)
	}
	statsAfterClear, _ := db.GetStats()
	if statsAfterClear.TotalRuns != 0 || statsAfterClear.ErrorLogCount != 0 {
		t.Errorf("expected 0 runs and 0 error logs after clear, got %d, %d",
			statsAfterClear.TotalRuns, statsAfterClear.ErrorLogCount)
	}

	// 7. Reset
	if err := db.Reset(); err != nil {
		t.Fatalf("failed to reset: %v", err)
	}
}

func TestPipelineDbGeneratedFields(t *testing.T) {
	// 1. Registry field checks
	if !PipelineRunRecordDb.IsRunId("RunId") {
		t.Errorf("expected PipelineRunRecordDb.IsRunId('RunId') to be true")
	}
	if !PipelineRunRecordDb.IsRunId(PipelineRunRecordDb.RunId) {
		t.Errorf("expected PipelineRunRecordDb.IsRunId(PipelineRunRecordDb.RunId) to be true")
	}
	if PipelineRunRecordDb.IsRunId("RepoSlug") {
		t.Errorf("expected PipelineRunRecordDb.IsRunId('RepoSlug') to be false")
	}

	// 2. Registry IsEnum
	if !PipelineRunRecordDb.IsEnum("RunId") {
		t.Errorf("expected PipelineRunRecordDb.IsEnum('RunId') to be true")
	}
	if !PipelineRunRecordDb.IsEnum(PipelineRunRecordDb.Sha) {
		t.Errorf("expected PipelineRunRecordDb.IsEnum(PipelineRunRecordDb.Sha) to be true")
	}
	if PipelineRunRecordDb.IsEnum("NonExistentColumn") {
		t.Errorf("expected PipelineRunRecordDb.IsEnum('NonExistentColumn') to be false")
	}

	// 3. Registry All and Names
	allFields := PipelineRunRecordDb.All()
	if len(allFields) != 15 {
		t.Errorf("expected 15 fields, got %d", len(allFields))
	}
	names := PipelineRunRecordDb.Names()
	if len(names) != 15 || names[0] != "RunId" {
		t.Errorf("unexpected names: %v", names)
	}

	// 4. Registry JSON serialization
	regJson, err := PipelineRunRecordDb.ToJSON()
	if err != nil || len(regJson) == 0 {
		t.Errorf("failed serializing registry to JSON: %v", err)
	}

	// 5. FieldType enum methods
	field := PipelineRunRecordDb.RunId
	if field.Name() != "RunId" || field.String() != "RunId" || field.Value() != "RunId" {
		t.Errorf("unexpected field values: %s, %s, %s", field.Name(), field.String(), field.Value())
	}
	if !field.IsCompare("RunId") || !field.IsEnum("RunId") {
		t.Errorf("expected IsCompare and IsEnum to be true for 'RunId'")
	}
	if field.IsCompare("RepoSlug") || field.IsEnum("RepoSlug") {
		t.Errorf("expected IsCompare and IsEnum to be false for 'RepoSlug'")
	}
	if !field.IsRunId() || !field.IsRunId("RunId") {
		t.Errorf("expected field.IsRunId to be true")
	}
	if field.IsRunId("RepoSlug") {
		t.Errorf("expected field.IsRunId('RepoSlug') to be false")
	}

	// 6. FieldType JSON serialization and deserialization
	jsonStr, err := field.ToJSON()
	if err != nil || jsonStr != `"RunId"` {
		t.Errorf("expected `\"RunId\"`, got %s (err: %v)", jsonStr, err)
	}

	var parsedField PipelineRunRecordFieldType
	if err := parsedField.FromJSON(`"WorkflowName"`); err != nil || parsedField != PipelineRunRecordDb.WorkflowName {
		t.Errorf("expected parsed WorkflowName, got %v (err: %v)", parsedField, err)
	}

	// 7. ErrorRecord field checks
	if !PipelineErrorRecordDb.IsRunId("RunId") {
		t.Errorf("expected PipelineErrorRecordDb.IsRunId('RunId') to be true")
	}
	if !PipelineErrorRecordDb.IsStepName("StepName") {
		t.Errorf("expected PipelineErrorRecordDb.IsStepName('StepName') to be true")
	}
	if PipelineErrorRecordDb.IsRunId("StepName") {
		t.Errorf("expected PipelineErrorRecordDb.IsRunId('StepName') to be false")
	}
}
