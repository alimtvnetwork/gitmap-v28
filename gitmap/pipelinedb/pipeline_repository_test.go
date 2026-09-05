package pipelinedb

import (
	"context"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	_ "modernc.org/sqlite"
)

func setupTestPipelineRepo(t *testing.T) (*PipelineRepository, func()) {
	t.Helper()
	ctx := context.Background()
	wrapper, appErr := dbengine.OpenDb(dbengine.DbSQLite, ":memory:")
	if appErr != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", appErr)
	}

	repo := NewPipelineRepository(wrapper)
	initErr := repo.InitSchema(ctx)
	if initErr != nil {
		_ = wrapper.Close()
		t.Fatalf("failed to init repository schema: %v", initErr)
	}

	cleanup := func() {
		_ = wrapper.Close()
	}
	return repo, cleanup
}

func TestPipelineRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := setupTestPipelineRepo(t)
	defer cleanup()

	// 1. Insert runs
	run1 := PipelineRunRecord{
		RunId:           101,
		RepoSlug:        "owner/repo",
		WorkflowName:    "build-test",
		Status:          "completed",
		Conclusion:      "success",
		Branch:          "main",
		Sha:             "commit-sha-1",
		EtaSeconds:      120,
		DurationSeconds: 110,
		RunUrl:          "https://ci.example.com/runs/101",
		IsSuccess:       true,
		CreatedAt:       "2026-09-05T10:00:00Z",
		UpdatedAt:       "2026-09-05T10:02:00Z",
	}
	run2 := PipelineRunRecord{
		RunId:           102,
		RepoSlug:        "owner/repo",
		WorkflowName:    "build-test",
		Status:          "completed",
		Conclusion:      "failure",
		Branch:          "feature",
		Sha:             "commit-sha-2",
		EtaSeconds:      150,
		DurationSeconds: 140,
		RunUrl:          "https://ci.example.com/runs/102",
		IsSuccess:       false,
		CreatedAt:       "2026-09-05T11:00:00Z",
		UpdatedAt:       "2026-09-05T11:03:00Z",
	}

	res1 := repo.InsertRun(ctx, run1)
	if res1.IsFailed() {
		t.Fatalf("failed to insert run 101: %v", res1.Err)
	}

	res2 := repo.InsertRun(ctx, run2)
	if res2.IsFailed() {
		t.Fatalf("failed to insert run 102: %v", res2.Err)
	}

	// 2. GetRunById
	getRes := repo.GetRunById(ctx, 101)
	if getRes.IsFailed() {
		t.Fatalf("failed to get run by id 101: %v", getRes.Err)
	}
	runRecord := getRes.Value
	if runRecord.WorkflowName != "build-test" {
		t.Errorf("expected workflow name 'build-test', got '%s'", runRecord.WorkflowName)
	}
	if !runRecord.IsSuccess {
		t.Errorf("expected run 101 to be successful")
	}

	// 3. GetRecentRuns
	recentRes := repo.GetRecentRuns(ctx, "owner/repo", 10)
	if recentRes.IsFailed() {
		t.Fatalf("failed to get recent runs: %v", recentRes.Err)
	}
	runs := recentRes.Value
	if len(runs) != 2 {
		t.Fatalf("expected 2 runs, got %d", len(runs))
	}
	if runs[0].RunId != 102 {
		t.Errorf("expected first run to be 102 (descending order), got %d", runs[0].RunId)
	}
}

func TestPipelineRepository_ActiveErrorsViewAndJoin(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := setupTestPipelineRepo(t)
	defer cleanup()

	// 1. Insert run and corresponding error record
	run := PipelineRunRecord{
		RunId:           201,
		RepoSlug:        "owner/repo",
		WorkflowName:    "deploy",
		Status:          "completed",
		Conclusion:      "failure",
		Branch:          "main",
		Sha:             "sha-deploy-fail",
		EtaSeconds:      300,
		DurationSeconds: 280,
		RunUrl:          "https://ci.example.com/runs/201",
		IsSuccess:       false,
		CreatedAt:       "2026-09-05T12:00:00Z",
		UpdatedAt:       "2026-09-05T12:05:00Z",
	}
	_ = repo.InsertRun(ctx, run)

	errRec := PipelineErrorRecord{
		RunId:        201,
		RepoSlug:     "owner/repo",
		WorkflowName: "deploy",
		StepName:     "terraform-apply",
		ErrorText:    "Access denied on cloud resource",
		CreatedAt:    "2026-09-05T12:04:00Z",
	}
	_ = repo.InsertErrorRecord(ctx, errRec)

	// 2. EnsureActiveErrorsView creation
	viewRes := repo.EnsureActiveErrorsView(ctx)
	if viewRes.IsFailed() {
		t.Fatalf("failed creating ActiveCiErrors view: %v", viewRes.Err)
	}

	// 3. Verify hash is recorded in __gitmap_view_meta
	hash1, hashErr := repo.Db().GetViewHash(ctx, "ActiveCiErrors")
	if hashErr != nil {
		t.Fatalf("failed retrieving view hash: %v", hashErr)
	}
	if len(hash1) == 0 {
		t.Fatalf("expected recorded query hash for ActiveCiErrors")
	}

	// 4. Second invocation should match hash and reuse view with 0 DDL
	reuseRes := repo.EnsureActiveErrorsView(ctx)
	if reuseRes.IsFailed() {
		t.Fatalf("failed re-ensuring ActiveCiErrors view: %v", reuseRes.Err)
	}
	hash2, _ := repo.Db().GetViewHash(ctx, "ActiveCiErrors")
	if hash1 != hash2 {
		t.Errorf("expected hash to be unchanged on reuse: %s vs %s", hash1, hash2)
	}

	// 5. Query ActiveCiErrors view directly
	rows, queryErr := repo.Db().Query(ctx, "SELECT RunId, RepoSlug, StepName, ErrorText FROM ActiveCiErrors;")
	if queryErr != nil {
		t.Fatalf("failed querying ActiveCiErrors view: %v", queryErr)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var runId uint64
		var repoSlug, stepName, errorText string
		scanErr := rows.Scan(&runId, &repoSlug, &stepName, &errorText)
		if scanErr != nil {
			t.Fatalf("failed scanning view row: %v", scanErr)
		}
		count++
		if runId != 201 || stepName != "terraform-apply" {
			t.Errorf("unexpected view row content: runId=%d, step=%s", runId, stepName)
		}
	}
	if count != 1 {
		t.Errorf("expected 1 joined error record in view, got %d", count)
	}
}

func TestPipelineRepository_FluentQueryWithEnums(t *testing.T) {
	repo, cleanup := setupTestPipelineRepo(t)
	defer cleanup()

	// Compile query with typed field enums and operator enums
	qb := repo.QueryBare().
		Select(PipelineRunRecordDb.RunId, PipelineRunRecordDb.WorkflowName).
		WhereOp(PipelineRunRecordDb.Status, dbengine.SqlOperators.Equal, "completed").
		WhereOp(PipelineRunRecordDb.Conclusion, dbengine.SqlOperators.Equal, "failure").
		OrderByDesc(PipelineRunRecordDb.RunId).
		Limit(5)

	compRes := qb.Compile()
	if compRes.IsFailed() {
		t.Fatalf("failed compiling query: %v", compRes.Err)
	}

	cq := compRes.Value
	if !strings.Contains(cq.SQL, "SELECT \"PipelineRunRecord\".\"RunId\", \"PipelineRunRecord\".\"WorkflowName\" FROM \"PipelineRunRecord\"") &&
		!strings.Contains(cq.SQL, "SELECT \"RunId\", \"WorkflowName\" FROM \"PipelineRunRecord\"") {
		t.Errorf("SQL missing expected projections: %s", cq.SQL)
	}
	if len(cq.Args) != 2 {
		t.Fatalf("expected 2 bound arguments, got %d", len(cq.Args))
	}
	if len(cq.QueryHash) == 0 {
		t.Errorf("expected non-empty QueryHash")
	}
}

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

