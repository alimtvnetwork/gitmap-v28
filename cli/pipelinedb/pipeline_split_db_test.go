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
	logRes := db.QueryRecentErrorLogs(5)
	if logRes.IsFailure() || logRes.Count() != 1 {
		t.Fatalf("expected 1 error log, got %d (err: %v)", logRes.Count(), logRes.AppError())
	}
	logs := logRes.Data

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
	// 1. Registry field checks (object-based)
	if !PipelineRunRecordDb.IsRunId(PipelineRunRecordDb.RunId) {
		t.Errorf("expected PipelineRunRecordDb.IsRunId(PipelineRunRecordDb.RunId) to be true")
	}

	if PipelineRunRecordDb.IsRunId(PipelineRunRecordDb.RepoSlug) {
		t.Errorf("expected PipelineRunRecordDb.IsRunId(PipelineRunRecordDb.RepoSlug) to be false")
	}

	// 2. Registry IsEnum (map-based object lookup)
	if !PipelineRunRecordDb.IsEnum(PipelineRunRecordDb.RunId) {
		t.Errorf("expected PipelineRunRecordDb.IsEnum(PipelineRunRecordDb.RunId) to be true")
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

	// 4. Registry JSON serialization with AppError
	regJson, appErr := PipelineRunRecordDb.ToJSON()
	if appErr != nil || len(regJson) == 0 {
		t.Errorf("failed serializing registry to JSON: %v", appErr)
	}

	// 5. FieldType enum methods
	field := PipelineRunRecordDb.RunId
	if field.Name() != "RunId" || field.String() != "RunId" || field.Value() != "RunId" {
		t.Errorf("unexpected field values: %s, %s, %s", field.Name(), field.String(), field.Value())
	}

	if !field.IsCompare(PipelineRunRecordDb.RunId) {
		t.Errorf("expected IsCompare to be true for RunId")
	}

	if field.IsCompare(PipelineRunRecordDb.RepoSlug) {
		t.Errorf("expected IsCompare to be false for RepoSlug")
	}

	if !field.IsEnum() {
		t.Errorf("expected field.IsEnum() to be true")
	}

	if !field.IsRunId() {
		t.Errorf("expected field.IsRunId() to be true")
	}

	if field.IsRepoSlug() {
		t.Errorf("expected field.IsRepoSlug() to be false")
	}

	// 6. FieldType JSON serialization and deserialization with AppError
	jsonStr, jsonErr := field.ToJSON()
	if jsonErr != nil || jsonStr != `"RunId"` {
		t.Errorf("expected `\"RunId\"`, got %s (err: %v)", jsonStr, jsonErr)
	}

	var parsedField PipelineRunRecordFieldType
	if err := parsedField.FromJSON(`"WorkflowName"`); err != nil || parsedField != PipelineRunRecordDb.WorkflowName {
		t.Errorf("expected parsed WorkflowName, got %v (err: %v)", parsedField, err)
	}

	// 7. ErrorRecord field checks
	if !PipelineErrorRecordDb.IsRunId(PipelineErrorRecordDb.RunId) {
		t.Errorf("expected PipelineErrorRecordDb.IsRunId(PipelineErrorRecordDb.RunId) to be true")
	}

	if !PipelineErrorRecordDb.IsStepName(PipelineErrorRecordDb.StepName) {
		t.Errorf("expected PipelineErrorRecordDb.IsStepName(PipelineErrorRecordDb.StepName) to be true")
	}

	if PipelineErrorRecordDb.IsRunId(PipelineErrorRecordDb.StepName) {
		t.Errorf("expected PipelineErrorRecordDb.IsRunId(PipelineErrorRecordDb.StepName) to be false")
	}
}

func TestPipelineDbScanError(t *testing.T) {
	db, err := OpenPipelineSplitDB("test-owner/test-scan-err")
	if err != nil {
		t.Fatalf("failed to open pipeline split db: %v", err)
	}

	defer db.Close()

	query := `INSERT INTO PipelineRun (
		RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
		EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
	) VALUES ('not_an_int', 's', 'w', 'st', 'c', 'b', 'sh', 'invalid', 'invalid', 'u', 'bad', 'c', 'u');`
	if _, execErr := db.conn.Exec(query); execErr != nil {
		return
	}

	if queryRes := db.QueryRecentRuns(5); queryRes.IsSuccess() {
		t.Errorf("expected QueryRecentRuns to propagate error on corrupt row")
	}
}

func TestPipelineSplitDB_3ConnectedTables(t *testing.T) {
	db, err := OpenPipelineSplitDB("test-owner/test-3tables")
	if err != nil {
		t.Fatalf("failed to open split db: %v", err)
	}
	defer db.Close()
	_ = db.Reset()

	verifyMasterAndDetailStorage(t, db)
	verifyCompactStorageAndQuery(t, db)
	verifyCleanAndReset3Tables(t, db)
}

func verifyMasterAndDetailStorage(t *testing.T, db *PipelineSplitDb) {
	run := PipelineRunRecord{
		RunId: 99001, RepoSlug: "test-owner/test-3tables", WorkflowName: "CI",
		Status: "completed", Conclusion: "failure", Branch: "main", Sha: "sha99001",
		CreatedAt: "2026-09-12T10:00:00Z", UpdatedAt: "2026-09-12T10:05:00Z",
	}
	if err := db.RecordRun(run); err != nil {
		t.Fatalf("failed to record master run: %v", err)
	}
	detail := PipelineErrorRecord{
		RunId: 99001, RepoSlug: "test-owner/test-3tables", WorkflowName: "CI",
		StepName: "Step 1", ErrorText: "fail", RawLogs: "ok 1\nFAIL 2\nok 3",
	}
	if err := db.RecordDetailErrorLog(detail); err != nil {
		t.Fatalf("failed to record detail error log: %v", err)
	}
}

func verifyCompactStorageAndQuery(t *testing.T, db *PipelineSplitDb) {
	compact := PipelineCompactErrorRecord{
		RunId: 99001, RepoSlug: "test-owner/test-3tables", WorkflowName: "CI",
		StepName: "Step 1", ErrorText: "fail", CompactLogs: "FAIL 2", FilteredOkCount: 2,
	}
	if err := db.RecordCompactErrorLog(compact); err != nil {
		t.Fatalf("failed to record compact error log: %v", err)
	}
	details := db.QueryDetailedErrorLogsByRunId(99001)
	compacts := db.QueryCompactErrorLogsByRunId(99001)
	if details.Count() != 1 || compacts.Count() != 1 {
		t.Fatalf("expected 1 detail and 1 compact log, got %d and %d", details.Count(), compacts.Count())
	}
	if compacts.Data[0].FilteredOkCount != 2 {
		t.Errorf("expected FilteredOkCount 2, got %d", compacts.Data[0].FilteredOkCount)
	}
}

func verifyCleanAndReset3Tables(t *testing.T, db *PipelineSplitDb) {
	if !db.HasDetailErrorLog(99001) || !db.HasCompactErrorLog(99001) {
		t.Errorf("expected HasDetailErrorLog and HasCompactErrorLog to be true")
	}
	if err := db.Clear(); err != nil {
		t.Fatalf("failed to clear db: %v", err)
	}
	details := db.QueryDetailedErrorLogsByRunId(99001)
	if !details.IsEmpty() {
		t.Errorf("expected 0 detail logs after clear, got %d", details.Count())
	}
}
