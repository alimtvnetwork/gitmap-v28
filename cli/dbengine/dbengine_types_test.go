package dbengine

import (
	"testing"

	_ "modernc.org/sqlite"
)

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

func TestSqlOperator_Methods(t *testing.T) {
	op := SqlOperators.Equal
	if op.Name() != "=" || op.String() != "=" || op.Value() != "=" {
		t.Errorf("unexpected op string: %s", op.String())
	}

	if !op.IsCompare(SqlOpEqual) {
		t.Errorf("expected IsCompare true for Equal")
	}

	if !op.IsEnum() {
		t.Errorf("expected IsEnum true for Equal")
	}

	if !op.IsEqual() {
		t.Errorf("expected IsEqual true for Equal")
	}

	if op.IsNotEqual() {
		t.Errorf("expected IsNotEqual false for Equal")
	}

	if !SqlOperators.NotEqual.IsNotEqual() {
		t.Errorf("expected IsNotEqual true for NotEqual")
	}

	if !SqlOperators.LessThan.IsLessThan() {
		t.Errorf("expected IsLessThan true for LessThan")
	}

	if !SqlOperators.GreaterThan.IsGreaterThan() {
		t.Errorf("expected IsGreaterThan true for GreaterThan")
	}

	if !SqlOperators.Like.IsLike() {
		t.Errorf("expected IsLike true for Like")
	}

	if !SqlOperators.In.IsIn() {
		t.Errorf("expected IsIn true for In")
	}

	invalidOp := SqlOperator("INVALID_OP")
	if invalidOp.IsEnum() {
		t.Errorf("expected invalidOp to not be enum")
	}

	jsonStr, appErr := op.ToJSON()
	if appErr != nil || jsonStr != `"="` {
		t.Errorf("unexpected json: %s (err: %v)", jsonStr, appErr)
	}

	var parsedOp SqlOperator
	if err := parsedOp.FromJSON(`"!="`); err != nil || parsedOp != SqlOperators.NotEqual {
		t.Errorf("unexpected parsed op: %v (err: %v)", parsedOp, err)
	}

	// Registry methods
	if len(SqlOperators.All()) != 13 {
		t.Errorf("expected 13 operators, got %d", len(SqlOperators.All()))
	}

	if len(SqlOperators.Names()) != 13 {
		t.Errorf("expected 13 operator names, got %d", len(SqlOperators.Names()))
	}

	if !SqlOperators.IsEnum(SqlOpLike) {
		t.Errorf("expected registry IsEnum true for Like")
	}

	if !SqlOperators.IsEqual(SqlOpEqual) {
		t.Errorf("expected registry IsEqual true for Equal")
	}
}
