package dbengine

import (
	"context"
	"database/sql"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
	"testing"
)

func TestResolveCompiler(t *testing.T) {
	dialects := []DatabaseDialectType{
		DatabaseDialectSQLite,
		DatabaseDialectPostgreSQL,
		DatabaseDialectMySQL,
		DatabaseDialectMariaDB,
		DatabaseDialectMSSQL,
		DatabaseDialectOracle,
		DatabaseDialectMongoDB,
	}

	for _, d := range dialects {
		compiler, err := ResolveCompiler(d)
		if err != nil {
			t.Fatalf("expected compiler for %s, got err: %v", d, err)
		}

		if compiler == nil {
			t.Fatalf("compiler for %s is nil", d)
		}
	}

	_, err := ResolveCompiler("unsupported")
	if err == nil {
		t.Fatal("expected error for unsupported dialect, got nil")
	}
}

func TestCompilerSyntaxes(t *testing.T) {
	sqliteComp := &SQLiteCompiler{}
	if sqliteComp.Placeholder(1) != "?" {
		t.Errorf("expected ?, got %s", sqliteComp.Placeholder(1))
	}

	searchSqlite := sqliteComp.CompileSearch("User", []string{"UserId", "Email"}, 5)
	expectedSqlite := `SELECT * FROM "User" WHERE "UserId" = ? AND "Email" = ? LIMIT 5;`
	if searchSqlite != expectedSqlite {
		t.Errorf("sqlite search mismatch:\ngot:  %s\nwant: %s", searchSqlite, expectedSqlite)
	}

	pgComp := &PostgresCompiler{}
	if pgComp.Placeholder(1) != "$1" || pgComp.Placeholder(2) != "$2" {
		t.Errorf("pg placeholder mismatch: %s, %s", pgComp.Placeholder(1), pgComp.Placeholder(2))
	}

	searchPg := pgComp.CompileSearch("User", []string{"UserId"}, 1)
	expectedPg := `SELECT * FROM "User" WHERE "UserId" = $1 LIMIT 1;`
	if searchPg != expectedPg {
		t.Errorf("pg search mismatch:\ngot:  %s\nwant: %s", searchPg, expectedPg)
	}

	mssqlComp := &MSSQLCompiler{}
	if mssqlComp.Placeholder(1) != "@p1" {
		t.Errorf("mssql placeholder mismatch: %s", mssqlComp.Placeholder(1))
	}

	searchMssql := mssqlComp.CompileSearch("User", []string{"UserId"}, 1)
	expectedMssql := `SELECT * FROM [User] WHERE [UserId] = @p1 ORDER BY (SELECT NULL) OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY;`
	if searchMssql != expectedMssql {
		t.Errorf("mssql search mismatch:\ngot:  %s\nwant: %s", searchMssql, expectedMssql)
	}
}

func setupInMemoryDb(t *testing.T) *DbWrapper {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	createSql := `
CREATE TABLE TestItem (
    ItemId INTEGER PRIMARY KEY AUTOINCREMENT,
    ItemName TEXT NOT NULL,
    Category TEXT NOT NULL,
    IsActive INTEGER NOT NULL DEFAULT 1
);
INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Alpha', 'Tool', 1);
INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Beta', 'Tool', 1);
INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Gamma', 'Service', 0);
`
	_, execErr := conn.Exec(createSql)
	if execErr != nil {
		t.Fatalf("failed creating test table: %v", execErr)
	}

	wrapper, appErr := WrapDb(conn, DatabaseDialectSQLite)
	if appErr != nil {
		t.Fatalf("WrapDb failed: %v", appErr)
	}

	return wrapper
}

func TestFluentQueryBuilderViaCompiler(t *testing.T) {
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

func TestQueryBuilder_CompileAndCacheViaCompiler(t *testing.T) {
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

func TestQueryBuilder_ErrorGuardsViaCompiler(t *testing.T) {
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
