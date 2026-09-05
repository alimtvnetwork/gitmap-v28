package dbengine

import (
	"fmt"
	"strings"
)

// SQLiteCompiler compiles queries for SQLite databases.
type SQLiteCompiler struct{}

// Dialect returns DatabaseDialectSQLite.
func (c *SQLiteCompiler) Dialect() DatabaseDialectType {
	return DatabaseDialectSQLite
}

// Placeholder returns ? for SQLite parameters.
func (c *SQLiteCompiler) Placeholder(paramIndex int) string {
	return "?"
}

// QuoteIdentifier wraps column/table names in double quotes.
func (c *SQLiteCompiler) QuoteIdentifier(name string) string {
	return "\"" + name + "\""
}

// CompilePagination produces LIMIT/OFFSET syntax.
func (c *SQLiteCompiler) CompilePagination(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	if offset <= 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

// CompileSearch builds a parameterized SELECT query.
func (c *SQLiteCompiler) CompileSearch(table string, fields []string, limit int) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(fields) == 0 {
		return buildSelectWithoutWhere(quotedTable, c.CompilePagination(limit, 0))
	}

	whereClauses := make([]string, 0, len(fields))
	for _, f := range fields {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", c.QuoteIdentifier(f)))
	}

	whereSql := strings.Join(whereClauses, " AND ")
	pagination := c.CompilePagination(limit, 0)
	return buildSelectWithWhere(quotedTable, whereSql, pagination)
}

func buildSelectWithoutWhere(quotedTable, pagination string) string {
	if len(pagination) == 0 {
		return fmt.Sprintf("SELECT * FROM %s;", quotedTable)
	}
	return fmt.Sprintf("SELECT * FROM %s %s;", quotedTable, pagination)
}

func buildSelectWithWhere(quotedTable, whereSql, pagination string) string {
	if len(pagination) == 0 {
		return fmt.Sprintf("SELECT * FROM %s WHERE %s;", quotedTable, whereSql)
	}
	return fmt.Sprintf("SELECT * FROM %s WHERE %s %s;", quotedTable, whereSql, pagination)
}
