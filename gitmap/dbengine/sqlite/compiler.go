package sqlite

import (
	"fmt"
	"strings"
)

// Compiler compiles queries for SQLite databases.
type Compiler struct{}

// NewCompiler returns a new SQLite Compiler instance.
func NewCompiler() *Compiler {
	return &Compiler{}
}

// Dialect returns "sqlite".
func (c *Compiler) Dialect() string {
	return "sqlite"
}

// Placeholder returns ? for SQLite parameters.
func (c *Compiler) Placeholder(paramIndex int) string {
	return "?"
}

// QuoteIdentifier wraps column or table names in double quotes.
func (c *Compiler) QuoteIdentifier(name string) string {
	return "\"" + name + "\""
}

// CompilePagination produces LIMIT/OFFSET syntax for SQLite.
func (c *Compiler) CompilePagination(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	if offset <= 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

// CompileLocate builds an INSTR(field, ?) > 0 clause for substring location.
func (c *Compiler) CompileLocate(field string) string {
	return fmt.Sprintf("INSTR(%s, ?) > 0", c.QuoteIdentifier(field))
}

// CompileSearch builds a parameterized SELECT query.
func (c *Compiler) CompileSearch(table string, fields []string, limit int) string {
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

// CompileCount builds a SELECT COUNT(*) query.
func (c *Compiler) CompileCount(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(field) == 0 {
		return fmt.Sprintf("SELECT COUNT(*) FROM %s;", quotedTable)
	}
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?;", quotedTable, c.QuoteIdentifier(field))
}

// CompileDelete builds a DELETE query.
func (c *Compiler) CompileDelete(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	return fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", quotedTable, c.QuoteIdentifier(field))
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
