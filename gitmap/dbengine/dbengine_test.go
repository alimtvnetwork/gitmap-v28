package dbengine

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	_ "modernc.org/sqlite"
)

type TestItem struct {
	ItemId   uint64
	ItemName string
	Category string
	IsActive bool
}

type TestItemFieldType string

func (e TestItemFieldType) Name() string   { return string(e) }
func (e TestItemFieldType) String() string { return string(e) }
func (e TestItemFieldType) Value() string  { return string(e) }
func (e TestItemFieldType) IsCompare(target TestItemFieldType) bool {
	return e == target
}

func (e TestItemFieldType) IsEnum() bool {
	return testItemValidMap[e]
}

func (e TestItemFieldType) IsItemId() bool   { return e == TestItemDb.ItemId }
func (e TestItemFieldType) IsItemName() bool { return e == TestItemDb.ItemName }
func (e TestItemFieldType) IsCategory() bool { return e == TestItemDb.Category }
func (e TestItemFieldType) IsIsActive() bool { return e == TestItemDb.IsActive }

func (e TestItemFieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

func (e *TestItemFieldType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	target := TestItemFieldType(s)
	if !testItemValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid test item enum: %s", s), "unmarshal test item field")
	}
	*e = target
	return nil
}

func (e TestItemFieldType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(e))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize field to json")
	}
	return string(b), nil
}

func (e *TestItemFieldType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize field from json")
	}
	target := TestItemFieldType(str)
	if !testItemValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid test item enum: %s", str), "validate field from json")
	}
	*e = target
	return nil
}

type testItemDbRegistry struct {
	ItemId   TestItemFieldType
	ItemName TestItemFieldType
	Category TestItemFieldType
	IsActive TestItemFieldType
}

func (r testItemDbRegistry) All() []TestItemFieldType {
	return []TestItemFieldType{r.ItemId, r.ItemName, r.Category, r.IsActive}
}

func (r testItemDbRegistry) Names() []string {
	return []string{"ItemId", "ItemName", "Category", "IsActive"}
}

func (r testItemDbRegistry) IsEnum(target TestItemFieldType) bool {
	return testItemValidMap[target]
}

func (r testItemDbRegistry) IsItemId(target TestItemFieldType) bool {
	return target == r.ItemId
}

func (r testItemDbRegistry) IsItemName(target TestItemFieldType) bool {
	return target == r.ItemName
}

func (r testItemDbRegistry) IsCategory(target TestItemFieldType) bool {
	return target == r.Category
}

func (r testItemDbRegistry) IsIsActive(target TestItemFieldType) bool {
	return target == r.IsActive
}

func (r testItemDbRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize test item db registry to json")
	}
	return string(b), nil
}

var TestItemDb = testItemDbRegistry{
	ItemId:   "ItemId",
	ItemName: "ItemName",
	Category: "Category",
	IsActive: "IsActive",
}

var testItemValidMap = map[TestItemFieldType]bool{
	TestItemDb.ItemId:   true,
	TestItemDb.ItemName: true,
	TestItemDb.Category: true,
	TestItemDb.IsActive: true,
}

var TestItemField = TestItemDb

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
	itemsRes := repo.FindAll(ctx, 10)
	if itemsRes.IsFailed() {
		t.Fatalf("FindAll failed: %v", itemsRes.Err)
	}
	if len(itemsRes.Value) != 4 {
		t.Errorf("expected 4 items (Delta committed, Echo rolled back), got %d", len(itemsRes.Value))
	}
}

func TestDbType_Methods(t *testing.T) {
	if DbTypes.SQLite.Name() != "sqlite" {
		t.Errorf("expected sqlite, got %s", DbTypes.SQLite.Name())
	}
	if DbTypes.SQLite.String() != "sqlite" {
		t.Errorf("expected sqlite, got %s", DbTypes.SQLite.String())
	}
	if DbTypes.SQLite.Value() != "sqlite" {
		t.Errorf("expected sqlite, got %s", DbTypes.SQLite.Value())
	}
	if !DbTypes.SQLite.IsCompare("sqlite") {
		t.Errorf("expected IsCompare('sqlite') to be true")
	}
	if !DbTypes.SQLite.IsCompare(DbSQLite) {
		t.Errorf("expected IsCompare(DbSQLite) to be true")
	}
	if DbTypes.SQLite.IsCompare("postgres") {
		t.Errorf("expected IsCompare('postgres') to be false")
	}

	// Test IsEnum and Is<Dialect> on registry (object based)
	if !DbTypes.IsEnum(DbSQLite) {
		t.Errorf("expected DbTypes.IsEnum(DbSQLite) to be true")
	}
	if !DbTypes.IsEnum(DbPostgreSQL) {
		t.Errorf("expected DbTypes.IsEnum(DbPostgreSQL) to be true")
	}
	if DbTypes.IsEnum("invalid_db") {
		t.Errorf("expected DbTypes.IsEnum('invalid_db') to be false")
	}
	if !DbTypes.IsSQLite(DbSQLite) {
		t.Errorf("expected DbTypes.IsSQLite(DbSQLite) to be true")
	}
	if DbTypes.IsSQLite(DbPostgreSQL) {
		t.Errorf("expected DbTypes.IsSQLite(DbPostgreSQL) to be false")
	}

	// Test Is<Dialect>() on DbType instance
	if !DbTypes.SQLite.IsSQLite() {
		t.Errorf("expected DbTypes.SQLite.IsSQLite() to be true")
	}
	if DbTypes.SQLite.IsPostgreSQL() {
		t.Errorf("expected DbTypes.SQLite.IsPostgreSQL() to be false")
	}

	// Test JSON methods with AppError
	jsonStr, appErr := DbTypes.SQLite.ToJSON()
	if appErr != nil || jsonStr != `"sqlite"` {
		t.Errorf("expected JSON `\"sqlite\"`, got %s (err: %v)", jsonStr, appErr)
	}

	var parsed DbType
	if err := parsed.FromJSON(`"postgres"`); err != nil || parsed != DbPostgreSQL {
		t.Errorf("expected parsed DbPostgreSQL, got %v (err: %v)", parsed, err)
	}
}

func TestFieldType_Methods(t *testing.T) {
	field := TestItemDb.ItemId
	if field.Name() != "ItemId" {
		t.Errorf("expected ItemId, got %s", field.Name())
	}
	if field.String() != "ItemId" {
		t.Errorf("expected ItemId, got %s", field.String())
	}
	if field.Value() != "ItemId" {
		t.Errorf("expected ItemId, got %s", field.Value())
	}
	if !field.IsCompare(TestItemDb.ItemId) {
		t.Errorf("expected IsCompare(TestItemDb.ItemId) to be true")
	}
	if field.IsCompare(TestItemDb.ItemName) {
		t.Errorf("expected IsCompare(TestItemDb.ItemName) to be false")
	}

	// Test IsEnum on field (zero args, checks map)
	if !field.IsEnum() {
		t.Errorf("expected field.IsEnum() to be true")
	}

	// Test field-specific object methods on field instance
	if !field.IsItemId() {
		t.Errorf("expected field.IsItemId() to be true")
	}
	if field.IsItemName() {
		t.Errorf("expected field.IsItemName() to be false")
	}

	// Test Is<Field>(target) on registry
	if !TestItemDb.IsItemId(TestItemDb.ItemId) {
		t.Errorf("expected TestItemDb.IsItemId(TestItemDb.ItemId) to be true")
	}
	if TestItemDb.IsItemId(TestItemDb.Category) {
		t.Errorf("expected TestItemDb.IsItemId(TestItemDb.Category) to be false")
	}

	// Test IsEnum on registry
	if !TestItemDb.IsEnum(TestItemDb.Category) {
		t.Errorf("expected TestItemDb.IsEnum(Category) to be true")
	}
	if TestItemDb.IsEnum("NonExistentColumn") {
		t.Errorf("expected TestItemDb.IsEnum('NonExistentColumn') to be false")
	}

	// Test JSON methods with AppError
	jsonStr, appErr := field.ToJSON()
	if appErr != nil || jsonStr != `"ItemId"` {
		t.Errorf("expected JSON `\"ItemId\"`, got %s (err: %v)", jsonStr, appErr)
	}

	var parsedField TestItemFieldType
	if err := parsedField.FromJSON(`"ItemName"`); err != nil || parsedField != TestItemDb.ItemName {
		t.Errorf("expected parsed ItemName, got %v (err: %v)", parsedField, err)
	}
}
