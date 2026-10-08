package dbengine

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

func TestFluentQueryBuilder(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	// 1. WhereEq and OrderByDesc
	queryRes := repo.Query().
		WhereEq(TestItemDb.Category, "Tool").
		OrderByDesc(TestItemDb.ItemId).
		FindAll(ctx)

	if queryRes.IsFailed() {
		t.Fatalf("Fluent query failed: %v", queryRes.Err)
	}

	if len(queryRes.Value) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(queryRes.Value))
	}

	if queryRes.Value[0].ItemId != 2 || queryRes.Value[1].ItemId != 1 {
		t.Errorf("expected descending order (2, 1), got (%d, %d)", queryRes.Value[0].ItemId, queryRes.Value[1].ItemId)
	}

	// 2. Locate substring test (INSTR in SQLite)
	locateRes := repo.Query().
		Locate(TestItemDb.ItemName, "lph").
		First(ctx)

	if locateRes.IsFailed() || locateRes.Value.ItemName != "Alpha" {
		t.Errorf("Locate failed: %+v (err: %v)", locateRes.Value, locateRes.Err)
	}

	// 3. Ad-hoc CTE view (WithView)
	cteRes := repo.Query().
		WithView("ActiveTools", "SELECT * FROM TestItem WHERE Category = 'Tool' AND IsActive = 1").
		WhereEq(TestItemDb.ItemName, "Beta").
		First(ctx)

	if cteRes.IsFailed() || cteRes.Value.ItemName != "Beta" {
		t.Errorf("CTE query failed: %+v (err: %v)", cteRes.Value, cteRes.Err)
	}
}

func TestQueryBuilder_CompileAndCache(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	GlobalQueryCache.Clear()

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)
	qb := repo.Query().
		Select(TestItemDb.ItemName).
		Where(TestItemDb.Category, "=", "Tool")

	res1 := qb.Compile()
	if res1.IsFailed() {
		t.Fatalf("first compile failed: %v", res1.Err)
	}

	cq1 := res1.Value
	if GlobalQueryCache.Size() != 1 {
		t.Errorf("expected cache size 1 after compile, got %d", GlobalQueryCache.Size())
	}

	res2 := qb.Compile()
	if res2.IsFailed() {
		t.Fatalf("second compile failed: %v", res2.Err)
	}

	cq2 := res2.Value
	if cq1.SQL != cq2.SQL {
		t.Errorf("expected cached SQL to match: %s vs %s", cq1.SQL, cq2.SQL)
	}

	if len(cq1.Args) != len(cq2.Args) {
		t.Errorf("expected args count to match")
	}

	if cq1.QueryHash != cq2.QueryHash {
		t.Errorf("expected query hashes to match: %s vs %s", cq1.QueryHash, cq2.QueryHash)
	}
}

func TestQueryBuilder_ErrorGuards(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)
	qb := repo.Query()

	simulatedErr := apperror.NewWithDetails(
		"TestQueryBuilder_ErrorGuards",
		"ERR_SIMULATED",
		"Simulated builder error",
		"system",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		nil,
	)

	qb.SetError(simulatedErr)

	// Chaining methods should preserve error and not panic
	qb.Select(TestItemDb.ItemId).
		Where(TestItemDb.IsActive, "=", true).
		OrderBy(TestItemDb.ItemId, "ASC").
		Limit(5)

	if qb.Err() != simulatedErr {
		t.Errorf("expected builder to retain error")
	}

	// Terminal methods must return failure results wrapping simulatedErr
	compileRes := qb.Compile()
	if compileRes.IsSuccess() {
		t.Errorf("expected Compile to fail when builder has error")
	}

	if compileRes.Err != simulatedErr {
		t.Errorf("expected exact simulatedErr in compile result")
	}

	firstRes := qb.First(ctx)
	if firstRes.IsSuccess() {
		t.Errorf("expected First to fail when builder has error")
	}

	findRes := qb.FindAll(ctx)
	if findRes.IsSuccess() {
		t.Errorf("expected FindAll to fail when builder has error")
	}

	countRes := qb.Count(ctx)
	if countRes.IsSuccess() {
		t.Errorf("expected Count to fail when builder has error")
	}

	delRes := qb.Delete(ctx)
	if delRes.IsSuccess() {
		t.Errorf("expected Delete to fail when builder has error")
	}

	viewRes := qb.CreateViewOrUseView(ctx, "ErrorView")
	if viewRes.IsSuccess() {
		t.Errorf("expected CreateViewOrUseView to fail when builder has error")
	}
}
