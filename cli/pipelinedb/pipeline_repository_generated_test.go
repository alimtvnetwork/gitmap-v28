package pipelinedb

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb/enums"
	_ "modernc.org/sqlite"
)

func TestPipelineRunDbRepo_GeneratedRepo(t *testing.T) {
	ctx := context.Background()
	domainRepo, cleanup := setupTestPipelineRepo(t)
	defer cleanup()

	// Insert test runs via domain repo
	run := PipelineRunRecord{
		RunId:           301,
		RepoSlug:        "owner/auto-repo",
		WorkflowName:    "ci-cd",
		Status:          "completed",
		Conclusion:      "success",
		Branch:          "main",
		Sha:             "sha-301",
		EtaSeconds:      60,
		DurationSeconds: 55,
		RunUrl:          "https://ci.example.com/runs/301",
		IsSuccess:       true,
		CreatedAt:       "2026-09-05T14:00:00Z",
		UpdatedAt:       "2026-09-05T14:01:00Z",
	}

	insRes := domainRepo.InsertRun(ctx, run)
	if insRes.IsFailed() {
		t.Fatalf("failed inserting run 301: %v", insRes.Err)
	}

	// 1. Direct usage of auto-generated PipelineRunDbRepo alias constructor
	runDbRepo := NewPipelineRunDbRepo(domainRepo.Db())
	if runDbRepo == nil {
		t.Fatalf("expected non-nil PipelineRunDbRepo")
	}

	// Verify dedicated single data type query builder
	var qb *PipelineRunQueryBuilder = runDbRepo.Query()
	if qb == nil {
		t.Fatalf("expected non-nil PipelineRunQueryBuilder")
	}

	// 2. FindAll returns dbengine.ListResult[PipelineRunRecord]
	listRes := runDbRepo.FindAll(ctx)
	if listRes.IsFailed() {
		t.Fatalf("FindAll failed: %v", listRes.Err)
	}

	if len(listRes.Value) != 1 {
		t.Fatalf("expected 1 record from FindAll, got %d", len(listRes.Value))
	}

	firstFound := listRes.Value[0]
	if firstFound.RunId != 301 {
		t.Errorf("expected RunId 301, got %d", firstFound.RunId)
	}

	// 3. First returns dbengine.EntityResult[PipelineRunRecord]
	entRes := runDbRepo.First(ctx)
	if entRes.IsFailed() {
		t.Fatalf("First failed: %v", entRes.Err)
	}

	if entRes.Value.WorkflowName != "ci-cd" {
		t.Errorf("expected workflow 'ci-cd', got '%s'", entRes.Value.WorkflowName)
	}

	// 4. Count returns dbengine.Int64Result
	countRes := runDbRepo.Count(ctx)
	if countRes.IsFailed() {
		t.Fatalf("Count failed: %v", countRes.Err)
	}

	if countRes.Value != 1 {
		t.Errorf("expected count 1, got %d", countRes.Value)
	}

	// 5. Query() fluent filtering with generated enums
	filteredRes := runDbRepo.Query().
		WhereOp(PipelineRunRecordDb.RunId, dbengine.SqlOperators.Equal, 301).
		FindAll(ctx)
	if filteredRes.IsFailed() {
		t.Fatalf("Query().WhereOp() failed: %v", filteredRes.Err)
	}

	if len(filteredRes.Value) != 1 {
		t.Errorf("expected 1 filtered record, got %d", len(filteredRes.Value))
	}

	// 6. Direct usage of auto-generated PipelineErrorDbRepo
	errDbRepo := NewPipelineErrorDbRepo(domainRepo.Db())
	if errDbRepo == nil {
		t.Fatalf("expected non-nil PipelineErrorDbRepo")
	}

	errCountRes := errDbRepo.Count(ctx)
	if errCountRes.IsFailed() {
		t.Fatalf("errDbRepo.Count failed: %v", errCountRes.Err)
	}

	if errCountRes.Value != 0 {
		t.Errorf("expected error count 0, got %d", errCountRes.Value)
	}
}

func TestPipelineDbGeneratedConsts(t *testing.T) {
	// Canonical table name constants
	if PipelineRunRecordTable != "PipelineRunRecord" {
		t.Errorf("expected PipelineRunRecord, got %s", PipelineRunRecordTable)
	}

	if PipelineRunTable != PipelineRunRecordTable {
		t.Errorf("expected PipelineRunTable == PipelineRunRecordTable")
	}

	if PipelineErrorRecordTable != "PipelineErrorRecord" {
		t.Errorf("expected PipelineErrorRecord, got %s", PipelineErrorRecordTable)
	}

	if PipelineErrorTable != PipelineErrorRecordTable {
		t.Errorf("expected PipelineErrorTable == PipelineErrorRecordTable")
	}

	if PipelineDbStatsTable != "PipelineDbStats" {
		t.Errorf("expected PipelineDbStats, got %s", PipelineDbStatsTable)
	}

	// Type aliases compile-time checks
	var _ *PipelineRunRecordQueryBuilder = nil
	var _ *PipelineRunQueryBuilder = nil
	var _ *PipelineErrorRecordQueryBuilder = nil
	var _ *PipelineErrorQueryBuilder = nil
	var _ *PipelineDbStatsQueryBuilder = nil

	var _ *PipelineRunRecordRepository = nil
	var _ *PipelineRunRepository = nil
	var _ *PipelineErrorRecordRepository = nil
	var _ *PipelineErrorRepository = nil
	var _ *PipelineDbStatsRepository = nil

	var _ *PipelineRunDbRepo = nil
	var _ *PipelineErrorDbRepo = nil
}

func TestPipelineDbEnumsPackageDirectUsage(t *testing.T) {
	// 1. Verify direct enums package table constants
	if enums.PipelineRunRecordTable != "PipelineRunRecord" {
		t.Errorf("expected enums.PipelineRunRecordTable == PipelineRunRecord")
	}

	if enums.PipelineRunTable != "PipelineRunRecord" {
		t.Errorf("expected enums.PipelineRunTable == PipelineRunRecord")
	}

	// 2. Verify direct enums package registry fields
	if enums.PipelineRunRecordDb.RunId != "RunId" {
		t.Errorf("expected enums.PipelineRunRecordDb.RunId == RunId")
	}

	if enums.PipelineRunDb.RunId != "RunId" {
		t.Errorf("expected enums.PipelineRunDb.RunId == RunId")
	}

	if enums.PipelineErrorRecordDb.ErrorText != "ErrorText" {
		t.Errorf("expected enums.PipelineErrorRecordDb.ErrorText == ErrorText")
	}

	if enums.PipelineErrorDb.ErrorText != "ErrorText" {
		t.Errorf("expected enums.PipelineErrorDb.ErrorText == ErrorText")
	}

	if enums.PipelineDbStatsDb.TotalRuns != "TotalRuns" {
		t.Errorf("expected enums.PipelineDbStatsDb.TotalRuns == TotalRuns")
	}

	// 3. Compile-time check: enums types match re-exported types
	var _ enums.PipelineRunRecordFieldType = PipelineRunRecordDb.RunId
	var _ PipelineRunRecordFieldType = enums.PipelineRunRecordDb.RunId
}

func TestPipelineRunRecordDbRepo_Mutations(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := setupTestPipelineRepo(t)
	defer cleanup()

	runRepo := repo.PipelineRunRecordDbRepo

	// Test Insert
	record := PipelineRunRecord{
		RunId:        501,
		RepoSlug:     "test/repo",
		WorkflowName: "ci",
		Status:       "in_progress",
		Conclusion:   "neutral",
		Branch:       "main",
		Sha:          "sha501",
		CreatedAt:    "2026-09-10T12:00:00Z",
		UpdatedAt:    "2026-09-10T12:00:00Z",
	}

	insRes := runRepo.Insert(ctx, &record)
	if insRes.IsFailed() {
		t.Fatalf("Insert failed: %v", insRes.Err)
	}

	// Test Update
	record.Status = "completed"
	record.Conclusion = "success"
	updRes := runRepo.Update(ctx, &record)
	if updRes.IsFailed() {
		t.Fatalf("Update failed: %v", updRes.Err)
	}

	// Verify update took effect
	findRes := runRepo.First(ctx)
	if findRes.IsFailed() || findRes.Value.Status != "completed" {
		t.Fatalf("expected updated status 'completed', got %v", findRes.Value)
	}

	// Test DeleteById
	delRes := runRepo.DeleteById(ctx, 501)
	if delRes.IsFailed() {
		t.Fatalf("DeleteById failed: %v", delRes.Err)
	}

	countRes := runRepo.Count(ctx)
	if countRes.IsFailed() || countRes.Value != 0 {
		t.Fatalf("expected 0 records after DeleteById, got %d", countRes.Value)
	}
}
