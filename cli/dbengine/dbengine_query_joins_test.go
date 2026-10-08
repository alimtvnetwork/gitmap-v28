package dbengine

import (
	"context"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestQueryBuilder_InnerJoinAndInnerWhere(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	// Create secondary table for join tests
	_, execErr := wrapper.Exec(ctx, `
CREATE TABLE TestDetail (
    DetailId INTEGER PRIMARY KEY AUTOINCREMENT,
    ItemId INTEGER NOT NULL,
    DetailText TEXT NOT NULL
);
INSERT INTO TestDetail (ItemId, DetailText) VALUES (1, 'Detail for Alpha');
INSERT INTO TestDetail (ItemId, DetailText) VALUES (2, 'Detail for Beta');
`)
	if execErr != nil {
		t.Fatalf("failed creating TestDetail: %v", execErr)
	}

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	// Test QueryBuilder with Select, scoped InnerJoin, and OnField
	qb := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName).
		InnerJoin("TestDetail").
		Select("DetailText").
		OnField(TestItemDb.ItemId, SqlOperators.Equal, "TestDetail.ItemId").
		Where(TestItemDb.ItemName, "=", "Alpha")

	sqlStr, args := qb.BuildSelect()
	expectedSql := `SELECT "TestItem"."ItemId", "TestItem"."ItemName", "TestDetail"."DetailText" FROM "TestItem" INNER JOIN "TestDetail" ON "TestItem"."ItemId" = "TestDetail"."ItemId" WHERE "TestItem"."ItemName" = ?;`
	if sqlStr != expectedSql {
		t.Errorf("BuildSelect mismatch:\ngot:  %s\nwant: %s", sqlStr, expectedSql)
	}

	if len(args) != 1 || args[0] != "Alpha" {
		t.Errorf("expected 1 arg ('Alpha'), got %v", args)
	}

	// Test joined query execution scanning projected columns
	row, queryErr := wrapper.QueryRow(ctx, sqlStr, args...)
	if queryErr != nil {
		t.Fatalf("joined query failed: %v", queryErr)
	}

	var id uint64
	var name, detail string
	if err := row.Scan(&id, &name, &detail); err != nil {
		t.Fatalf("joined scan failed: %v", err)
	}

	if name != "Alpha" || detail != "Detail for Alpha" {
		t.Errorf("unexpected scan result: %s / %s", name, detail)
	}

	// Test typed entity execution returning full model via repo.Query()
	firstRes := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName, TestItemDb.Category, TestItemDb.IsActive).
		InnerJoin("TestDetail").
		OnField(TestItemDb.ItemId, SqlOperators.Equal, "TestDetail.ItemId").
		Where(TestItemDb.ItemName, "=", "Alpha").
		First(ctx)
	if firstRes.IsFailed() {
		t.Fatalf("first execution failed: %v", firstRes.Err)
	}

	if firstRes.Value.ItemName != "Alpha" {
		t.Errorf("expected Alpha, got %s", firstRes.Value.ItemName)
	}

	// Test Dynamic Query Builder via SelectTable
	dynQb := SelectTable(wrapper, "TestItem", "ItemId", "ItemName").
		Where("Category", "=", "Tool").
		OrderBy("ItemId", "ASC")

	dynRes := dynQb.Compile()
	if dynRes.IsFailed() {
		t.Fatalf("dynamic query compile failed: %v", dynRes.Err)
	}

	dynCq := dynRes.Value
	// Verify dynSql contains expected tables and clauses
	if !strings.Contains(dynCq.SQL, "SELECT \"ItemId\", \"ItemName\" FROM \"TestItem\"") {
		t.Errorf("dynamic query missing projected columns: %s", dynCq.SQL)
	}

	if len(dynCq.Args) == 0 {
		t.Errorf("expected dynamic query args")
	}

	if len(dynCq.QueryHash) == 0 {
		t.Errorf("expected non-empty QueryHash")
	}
}

func TestQueryBuilder_JoinsAndGroupBy(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	_, execErr := wrapper.Exec(ctx, `
CREATE TABLE TestDetail (
    DetailId INTEGER PRIMARY KEY AUTOINCREMENT,
    ItemId INTEGER NOT NULL,
    DetailText TEXT NOT NULL
);
INSERT INTO TestDetail (ItemId, DetailText) VALUES (1, 'Detail for Alpha');
INSERT INTO TestDetail (ItemId, DetailText) VALUES (1, 'Second detail for Alpha');
INSERT INTO TestDetail (ItemId, DetailText) VALUES (2, 'Detail for Beta');
`)
	if execErr != nil {
		t.Fatalf("failed creating TestDetail: %v", execErr)
	}

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	// Test LeftJoin with And filter on JoinBuilder
	leftQb := repo.Query().
		Select(TestItemDb.ItemName).
		LeftJoin("TestDetail").
		Select("DetailText").
		And("DetailText", SqlOperators.Like, "%Alpha%").
		OnField(TestItemDb.ItemId, SqlOperators.Equal, "TestDetail.ItemId").
		Where(TestItemDb.Category, "=", "Tool")

	leftSql, leftArgs := leftQb.BuildSelect()
	if !strings.Contains(leftSql, "LEFT JOIN \"TestDetail\" ON \"TestItem\".\"ItemId\" = \"TestDetail\".\"ItemId\" AND \"TestDetail\".\"DetailText\" LIKE ?") {
		t.Errorf("LeftJoin SQL mismatch: %s", leftSql)
	}

	if len(leftArgs) != 2 {
		t.Errorf("expected 2 args for left join query, got %d", len(leftArgs))
	}

	// Test GroupBy and HavingCount
	groupQb := repo.Query().
		Select(TestItemDb.Category).
		GroupBy(TestItemDb.Category).
		HavingCount(SqlOperators.GreaterThan, 1)

	groupSql, groupArgs := groupQb.BuildSelect()
	expectedGroup := `SELECT "Category" FROM "TestItem" GROUP BY "TestItem"."Category" HAVING COUNT(*) > ?;`
	if groupSql != expectedGroup {
		t.Errorf("GroupBy SQL mismatch:\ngot:  %s\nwant: %s", groupSql, expectedGroup)
	}

	if len(groupArgs) != 1 || groupArgs[0] != int64(1) {
		t.Errorf("expected having count arg 1, got %v", groupArgs)
	}

	// Test RightJoin and OuterJoin constructors
	rightQb := repo.Query().RightJoin("TestDetail").On("1=1")
	rightSql, _ := rightQb.BuildSelect()
	if !strings.Contains(rightSql, "RIGHT JOIN \"TestDetail\" ON 1=1") {
		t.Errorf("RightJoin SQL mismatch: %s", rightSql)
	}

	outerQb := repo.Query().OuterJoin("TestDetail").On("1=1")
	outerSql, _ := outerQb.BuildSelect()
	if !strings.Contains(outerSql, "FULL OUTER JOIN \"TestDetail\" ON 1=1") {
		t.Errorf("OuterJoin SQL mismatch: %s", outerSql)
	}
}

func TestQueryBuilder_TypedJoinsWithoutMagicStrings(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	const TestDetailTable = "TestDetail"
	type TestDetailFieldType string
	const (
		DetailItemId TestDetailFieldType = "ItemId"
		DetailText   TestDetailFieldType = "DetailText"
	)

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	// Zero string literals: pass table constant and typed field enums to JoinBuilder
	qb := repo.Query().
		Select(TestItemDb.ItemId, TestItemDb.ItemName).
		InnerJoin(TestDetailTable).
		Select(DetailText).
		OnField(TestItemDb.ItemId, SqlOperators.Equal, DetailItemId).
		InnerWhere(TestItemDb.ItemId, SqlOperators.Equal, DetailItemId)

	compRes := qb.Compile()
	if compRes.IsFailed() {
		t.Fatalf("compile failed: %v", compRes.Err)
	}

	cq := compRes.Value
	expectedSub := "INNER JOIN \"TestDetail\" ON \"TestItem\".\"ItemId\" = \"TestDetail\".\"ItemId\""
	if !strings.Contains(cq.SQL, expectedSub) {
		t.Errorf("expected SQL to contain '%s', got: %s", expectedSub, cq.SQL)
	}
}
