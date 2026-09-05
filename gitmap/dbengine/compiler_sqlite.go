package dbengine

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine/sqlite"
)

// SQLiteCompiler compiles queries for SQLite databases.
type SQLiteCompiler struct {
	impl *sqlite.Compiler
}

func (c *SQLiteCompiler) getImpl() *sqlite.Compiler {
	if c.impl == nil {
		c.impl = sqlite.NewCompiler()
	}
	return c.impl
}

// Dialect returns DatabaseDialectSQLite.
func (c *SQLiteCompiler) Dialect() DatabaseDialectType {
	return DatabaseDialectSQLite
}

// Placeholder returns ? for SQLite parameters.
func (c *SQLiteCompiler) Placeholder(paramIndex int) string {
	return c.getImpl().Placeholder(paramIndex)
}

// QuoteIdentifier wraps column/table names in double quotes.
func (c *SQLiteCompiler) QuoteIdentifier(name string) string {
	return c.getImpl().QuoteIdentifier(name)
}

// CompilePagination produces LIMIT/OFFSET syntax.
func (c *SQLiteCompiler) CompilePagination(limit, offset int) string {
	return c.getImpl().CompilePagination(limit, offset)
}

// CompileSearch builds a parameterized SELECT query.
func (c *SQLiteCompiler) CompileSearch(table string, fields []string, limit int) string {
	return c.getImpl().CompileSearch(table, fields, limit)
}

// CompileCreateView builds a CREATE VIEW IF NOT EXISTS statement.
func (c *SQLiteCompiler) CompileCreateView(name string, selectSql string) string {
	return c.getImpl().CompileCreateView(name, selectSql)
}

// CompileDropView builds a DROP VIEW IF EXISTS statement.
func (c *SQLiteCompiler) CompileDropView(name string) string {
	return c.getImpl().CompileDropView(name)
}

// CompileFunctionCall builds a SELECT <func>(args...) query.
func (c *SQLiteCompiler) CompileFunctionCall(name string, argCount int) string {
	return c.getImpl().CompileFunctionCall(name, argCount)
}

// CompileCount builds a SELECT COUNT(*) query.
func (c *SQLiteCompiler) CompileCount(table string, field string) string {
	return c.getImpl().CompileCount(table, field)
}

// CompileDelete builds a DELETE query.
func (c *SQLiteCompiler) CompileDelete(table string, field string) string {
	return c.getImpl().CompileDelete(table, field)
}

// CompileInspectColumns returns the PRAGMA query to inspect columns of a table or view.
func (c *SQLiteCompiler) CompileInspectColumns(tableOrView string) string {
	return c.getImpl().CompileInspectColumns(tableOrView)
}

// CompileInspectViewExists returns a query to check if a view exists in sqlite_master.
func (c *SQLiteCompiler) CompileInspectViewExists(viewName string) string {
	return c.getImpl().CompileInspectViewExists(viewName)
}

func buildSelectWithoutWhere(quotedTable, pagination string) string {
	if len(pagination) == 0 {
		return "SELECT * FROM " + quotedTable + ";"
	}
	return "SELECT * FROM " + quotedTable + " " + pagination + ";"
}

func buildSelectWithWhere(quotedTable, whereSql, pagination string) string {
	if len(pagination) == 0 {
		return "SELECT * FROM " + quotedTable + " WHERE " + whereSql + ";"
	}
	return "SELECT * FROM " + quotedTable + " WHERE " + whereSql + " " + pagination + ";"
}
