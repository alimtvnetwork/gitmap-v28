package sqlite

import (
	"testing"
)

func TestCompiler_Basics(t *testing.T) {
	c := NewCompiler()

	if c.Dialect() != "sqlite" {
		t.Errorf("expected sqlite dialect, got %s", c.Dialect())
	}
	if c.Placeholder(1) != "?" {
		t.Errorf("expected ?, got %s", c.Placeholder(1))
	}
	if c.QuoteIdentifier("User") != "\"User\"" {
		t.Errorf("expected \"User\", got %s", c.QuoteIdentifier("User"))
	}
}

func TestCompiler_Pagination(t *testing.T) {
	c := NewCompiler()

	if c.CompilePagination(0, 0) != "" {
		t.Errorf("expected empty pagination for 0,0, got %s", c.CompilePagination(0, 0))
	}
	if c.CompilePagination(10, 0) != "LIMIT 10" {
		t.Errorf("expected LIMIT 10, got %s", c.CompilePagination(10, 0))
	}
	if c.CompilePagination(10, 20) != "LIMIT 10 OFFSET 20" {
		t.Errorf("expected LIMIT 10 OFFSET 20, got %s", c.CompilePagination(10, 20))
	}
}

func TestCompiler_Locate(t *testing.T) {
	c := NewCompiler()

	locateSql := c.CompileLocate("WorkflowName")
	expected := `INSTR("WorkflowName", ?) > 0`
	if locateSql != expected {
		t.Errorf("expected %s, got %s", expected, locateSql)
	}
}

func TestCompiler_Search(t *testing.T) {
	c := NewCompiler()

	allSql := c.CompileSearch("TestTable", nil, 5)
	expectedAll := `SELECT * FROM "TestTable" LIMIT 5;`
	if allSql != expectedAll {
		t.Errorf("expected %s, got %s", expectedAll, allSql)
	}

	searchSql := c.CompileSearch("TestTable", []string{"ColA", "ColB"}, 10)
	expectedSearch := `SELECT * FROM "TestTable" WHERE "ColA" = ? AND "ColB" = ? LIMIT 10;`
	if searchSql != expectedSearch {
		t.Errorf("expected %s, got %s", expectedSearch, searchSql)
	}
}

func TestCompiler_Views(t *testing.T) {
	c := NewCompiler()

	createSql := c.CompileCreateView("ActiveUsers", "SELECT * FROM Users WHERE IsActive = 1")
	expectedCreate := `CREATE VIEW IF NOT EXISTS "ActiveUsers" AS SELECT * FROM Users WHERE IsActive = 1;`
	if createSql != expectedCreate {
		t.Errorf("expected %s, got %s", expectedCreate, createSql)
	}

	dropSql := c.CompileDropView("ActiveUsers")
	expectedDrop := `DROP VIEW IF EXISTS "ActiveUsers";`
	if dropSql != expectedDrop {
		t.Errorf("expected %s, got %s", expectedDrop, dropSql)
	}

	cteSql := c.CompileAdHocCTE("Filtered", "SELECT * FROM Users", "SELECT * FROM Filtered")
	expectedCte := `WITH "Filtered" AS (SELECT * FROM Users) SELECT * FROM Filtered;`
	if cteSql != expectedCte {
		t.Errorf("expected %s, got %s", expectedCte, cteSql)
	}
}

func TestCompiler_Functions(t *testing.T) {
	c := NewCompiler()

	callZero := c.CompileFunctionCall("sqlite_version", 0)
	if callZero != "SELECT sqlite_version();" {
		t.Errorf("expected SELECT sqlite_version();, got %s", callZero)
	}

	callArgs := c.CompileFunctionCall("substr", 3)
	if callArgs != "SELECT substr(?, ?, ?);" {
		t.Errorf("expected SELECT substr(?, ?, ?);, got %s", callArgs)
	}

	scalarExpr := c.CompileScalarFunctionExpression("coalesce", 2)
	if scalarExpr != "coalesce(?, ?)" {
		t.Errorf("expected coalesce(?, ?), got %s", scalarExpr)
	}

	scalarNoArgs := c.CompileScalarFunctionExpression("random", 0)
	if scalarNoArgs != "random()" {
		t.Errorf("expected random(), got %s", scalarNoArgs)
	}
}

func TestCompiler_CountAndDelete(t *testing.T) {
	c := NewCompiler()

	countAll := c.CompileCount("Items", "")
	if countAll != `SELECT COUNT(*) FROM "Items";` {
		t.Errorf("unexpected countAll: %s", countAll)
	}

	countBy := c.CompileCount("Items", "Category")
	if countBy != `SELECT COUNT(*) FROM "Items" WHERE "Category" = ?;` {
		t.Errorf("unexpected countBy: %s", countBy)
	}

	del := c.CompileDelete("Items", "ItemId")
	if del != `DELETE FROM "Items" WHERE "ItemId" = ?;` {
		t.Errorf("unexpected del: %s", del)
	}
}
