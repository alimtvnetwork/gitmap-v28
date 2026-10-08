package dbengine

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRepository_Queries(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	// First: limit 1
	firstRes := repo.First(ctx, TestItemDb.ItemName, "Alpha")
	if firstRes.IsFailed() {
		t.Fatalf("First failed: %v", firstRes.Err)
	}

	item := firstRes.Value
	if item.ItemName != "Alpha" || item.ItemId != 1 {
		t.Errorf("unexpected item: %+v", item)
	}

	// FindById: primary key lookup
	idRes := repo.FindById(ctx, TestItemDb.ItemId, 1)
	if idRes.IsFailed() || idRes.Value.ItemName != "Alpha" {
		t.Errorf("FindById failed: %+v (err: %v)", idRes.Value, idRes.Err)
	}

	// FindBy: 1-parameter
	toolsRes := repo.FindBy(ctx, TestItemDb.Category, "Tool", 10)
	if toolsRes.IsFailed() {
		t.Fatalf("FindBy failed: %v", toolsRes.Err)
	}

	if len(toolsRes.Value) != 2 {
		t.Errorf("expected 2 tools, got %d", len(toolsRes.Value))
	}

	// FindBy2: 2-parameters
	activeRes := repo.FindBy2(ctx, TestItemDb.Category, "Tool", TestItemDb.IsActive, 1, 10)
	if activeRes.IsFailed() {
		t.Fatalf("FindBy2 failed: %v", activeRes.Err)
	}

	if len(activeRes.Value) != 2 {
		t.Errorf("expected 2 active tools, got %d", len(activeRes.Value))
	}

	// FindAll
	allRes := repo.FindAll(ctx, 10)
	if allRes.IsFailed() {
		t.Fatalf("FindAll failed: %v", allRes.Err)
	}

	if len(allRes.Value) != 3 {
		t.Errorf("expected 3 total items, got %d", len(allRes.Value))
	}

	// Count & CountAll
	countRes := repo.Count(ctx, TestItemDb.Category, "Tool")
	if countRes.IsFailed() || countRes.Value != 2 {
		t.Errorf("expected 2 tools count, got %d (err: %v)", countRes.Value, countRes.Err)
	}

	totalCountRes := repo.CountAll(ctx)
	if totalCountRes.IsFailed() || totalCountRes.Value != 3 {
		t.Errorf("expected 3 total count, got %d (err: %v)", totalCountRes.Value, totalCountRes.Err)
	}
}

func TestDatabaseViewsAndFunctions(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	// 1. Create View
	createRes := wrapper.CreateView(ctx, "ActiveToolsView", "SELECT * FROM TestItem WHERE Category = 'Tool' AND IsActive = 1")
	if createRes.IsFailed() {
		t.Fatalf("CreateView failed: %v", createRes.Err)
	}

	// 2. Query the created view using Repository
	viewRepo := NewRepository[TestItem, TestItemFieldType](wrapper, "ActiveToolsView", scanTestItem)
	viewItemsRes := viewRepo.FindAll(ctx, 10)
	if viewItemsRes.IsFailed() || len(viewItemsRes.Value) != 2 {
		t.Errorf("expected 2 items from view, got %d (err: %v)", len(viewItemsRes.Value), viewItemsRes.Err)
	}

	// 3. Call database function
	funcRes := wrapper.CallFunction(ctx, "UPPER", "test string")
	if funcRes.IsFailed() || funcRes.Value != "TEST STRING" {
		t.Errorf("expected 'TEST STRING', got %s (err: %v)", funcRes.Value, funcRes.Err)
	}

	// 4. Drop View
	dropRes := wrapper.DropView(ctx, "ActiveToolsView")
	if dropRes.IsFailed() {
		t.Errorf("DropView failed: %v", dropRes.Err)
	}
}

func TestCompiledQueryCache(t *testing.T) {
	cache := NewCompiledQueryCache()
	if cache.Size() != 0 {
		t.Errorf("expected initial size 0, got %d", cache.Size())
	}

	cache.Put("key1", "SELECT 1;")
	val, found := cache.Get("key1")
	if !found || val != "SELECT 1;" {
		t.Errorf("expected cache hit with SELECT 1;, got %s (found: %v)", val, found)
	}

	_, found2 := cache.Get("nonexistent")
	if found2 {
		t.Errorf("expected cache miss for nonexistent key")
	}

	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", cache.Size())
	}
}

func TestDbWrapper_CreateViewOrUseView(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	viewName := "ActiveToolsView"
	// Initial creation: view does not exist yet
	res1 := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName, TestItemDb.Category).
		Where(TestItemDb.Category, "=", "Tool").
		CreateViewOrUseView(ctx, viewName, "ItemId", "ItemName", "Category")

	if res1.IsFailed() {
		t.Fatalf("first CreateViewOrUseView failed: %v", res1.Err)
	}

	// Verify view exists in SQLite
	exists, appErr := wrapper.ViewExists(ctx, viewName)
	if appErr != nil || !exists {
		t.Fatalf("expected view to exist: %v (err: %v)", exists, appErr)
	}

	// Verify columns exist
	cols, colErr := wrapper.GetTableColumns(ctx, viewName)
	if colErr != nil || len(cols) != 3 {
		t.Fatalf("expected 3 columns in view, got %v (err: %v)", cols, colErr)
	}

	// Second run: view exists with matching columns -> reuses existing view
	res2 := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName, TestItemDb.Category).
		Where(TestItemDb.Category, "=", "Tool").
		CreateViewOrUseView(ctx, viewName, "ItemId", "ItemName", "Category")

	if res2.IsFailed() {
		t.Fatalf("second CreateViewOrUseView (reuse) failed: %v", res2.Err)
	}

	// Third run: new required column requested (e.g. IsActive) that does not exist in old view
	// Should drop old view and recreate with the new column definition
	res3 := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName, TestItemDb.Category, TestItemDb.IsActive).
		Where(TestItemDb.Category, "=", "Tool").
		CreateViewOrUseView(ctx, viewName, "ItemId", "ItemName", "Category", "IsActive")

	if res3.IsFailed() {
		t.Fatalf("third CreateViewOrUseView (recreate with new column) failed: %v", res3.Err)
	}

	updatedCols, _ := wrapper.GetTableColumns(ctx, viewName)
	if len(updatedCols) != 4 {
		t.Errorf("expected 4 columns after view recreation, got %v", updatedCols)
	}
}

func TestDbWrapper_ValidateSql(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	// Valid SQL should pass validation
	valErr := wrapper.ValidateSql(ctx, "SELECT ItemId, ItemName FROM TestItem WHERE IsActive = 1")
	if valErr != nil {
		t.Fatalf("expected valid SQL to pass validation, got: %v", valErr)
	}

	// Invalid SQL syntax should fail validation
	invalidErr := wrapper.ValidateSql(ctx, "SELECT FROM WHERE INVALID SYNTAX")
	if invalidErr == nil {
		t.Fatalf("expected invalid SQL syntax to fail validation, but it passed")
	}

	// Non-existent table should fail validation
	noTableErr := wrapper.ValidateSql(ctx, "SELECT * FROM NonExistentTableXYZ")
	if noTableErr == nil {
		t.Fatalf("expected non-existent table to fail validation, but it passed")
	}
}

func TestDbWrapper_ViewHashMetaAndEvolution(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	viewName := "DynamicActiveView"

	// 1. Initial creation via QueryBuilder with 0 manual column parameters
	qb1 := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName).
		Where(TestItemDb.IsActive, "=", true)

	res1 := qb1.CreateViewOrUseView(ctx, viewName)
	if res1.IsFailed() {
		t.Fatalf("initial CreateViewOrUseView failed: %v", res1.Err)
	}

	hash1, hashErr1 := wrapper.GetViewHash(ctx, viewName)
	if hashErr1 != nil {
		t.Fatalf("failed retrieving view hash: %v", hashErr1)
	}

	if len(hash1) == 0 {
		t.Fatalf("expected recorded query hash in __gitmap_view_meta")
	}

	// 2. Reuse view when query hash is identical
	res2 := qb1.CreateViewOrUseView(ctx, viewName)
	if res2.IsFailed() {
		t.Fatalf("subsequent CreateViewOrUseView failed: %v", res2.Err)
	}

	hash2, _ := wrapper.GetViewHash(ctx, viewName)
	if hash1 != hash2 {
		t.Errorf("expected view hash to remain identical on reuse: %s vs %s", hash1, hash2)
	}

	// 3. Automated evolution when query changes (new column added)
	qb2 := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName, TestItemDb.Category).
		Where(TestItemDb.IsActive, "=", true)

	res3 := qb2.CreateViewOrUseView(ctx, viewName)
	if res3.IsFailed() {
		t.Fatalf("evolved CreateViewOrUseView failed: %v", res3.Err)
	}

	hash3, _ := wrapper.GetViewHash(ctx, viewName)
	if hash1 == hash3 {
		t.Errorf("expected view hash to update after query change, but was identical: %s", hash3)
	}

	// 4. Verify columns in evolved view
	cols, colErr := wrapper.GetTableColumns(ctx, viewName)
	if colErr != nil {
		t.Fatalf("failed getting columns of evolved view: %v", colErr)
	}

	if len(cols) != 3 {
		t.Fatalf("expected 3 columns in evolved view, got %d: %v", len(cols), cols)
	}
}
