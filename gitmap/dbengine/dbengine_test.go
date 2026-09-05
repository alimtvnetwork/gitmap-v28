package dbengine

import (
	"context"
	"database/sql"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	_ "modernc.org/sqlite"
)

type TestItem struct {
	ItemId   int64
	ItemName string
	Category string
	IsActive bool
}

type TestItemFieldType string

const (
	TestItemFieldItemId   TestItemFieldType = "ItemId"
	TestItemFieldItemName TestItemFieldType = "ItemName"
	TestItemFieldCategory TestItemFieldType = "Category"
	TestItemFieldIsActive TestItemFieldType = "IsActive"
)

func scanTestItem(s RowScanner) (*TestItem, error) {
	var item TestItem
	var activeInt int
	err := s.Scan(&item.ItemId, &item.ItemName, &item.Category, &activeInt)
	if err != nil {
		return nil, err
	}
	item.IsActive = activeInt == 1
	return &item, nil
}

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
	expectedMssql := `SELECT * FROM [User] WHERE [UserId] = @p1 OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY;`
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

func TestRepository_Queries(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)

	// First: limit 1
	item, appErr := repo.First(ctx, TestItemFieldItemName, "Alpha")
	if appErr != nil {
		t.Fatalf("First failed: %v", appErr)
	}
	if item.ItemName != "Alpha" || item.ItemId != 1 {
		t.Errorf("unexpected item: %+v", item)
	}

	// FindBy: 1-parameter
	tools, appErr := repo.FindBy(ctx, TestItemFieldCategory, "Tool", 10)
	if appErr != nil {
		t.Fatalf("FindBy failed: %v", appErr)
	}
	if len(tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(tools))
	}

	// FindBy2: 2-parameters
	activeTools, appErr := repo.FindBy2(ctx, TestItemFieldCategory, "Tool", TestItemFieldIsActive, 1, 10)
	if appErr != nil {
		t.Fatalf("FindBy2 failed: %v", appErr)
	}
	if len(activeTools) != 2 {
		t.Errorf("expected 2 active tools, got %d", len(activeTools))
	}

	// FindAll
	allItems, appErr := repo.FindAll(ctx, 10)
	if appErr != nil {
		t.Fatalf("FindAll failed: %v", appErr)
	}
	if len(allItems) != 3 {
		t.Errorf("expected 3 total items, got %d", len(allItems))
	}
}

func TestTransaction_CommitAndRollback(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	// Successful transaction
	err := wrapper.WithTransaction(ctx, func(tx *TxWrapper) *apperror.AppError {
		_, execErr := tx.tx.Exec("INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Service', 1)")
		if execErr != nil {
			return apperror.WrapSimple(execErr, "insert delta")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected transaction success, got %v", err)
	}

	// Rollback transaction
	_ = wrapper.WithTransaction(ctx, func(tx *TxWrapper) *apperror.AppError {
		_, _ = tx.tx.Exec("INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Echo', 'Service', 1)")
		return apperror.WrapSimple(sql.ErrTxDone, "simulated tx failure")
	})

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)
	items, _ := repo.FindAll(ctx, 10)
	if len(items) != 4 {
		t.Errorf("expected 4 items (Delta committed, Echo rolled back), got %d", len(items))
	}
}
