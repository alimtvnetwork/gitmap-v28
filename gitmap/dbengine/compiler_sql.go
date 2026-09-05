package dbengine

import (
	"fmt"
	"strings"
)

// PostgresCompiler compiles queries for PostgreSQL databases.
type PostgresCompiler struct{}

func (c *PostgresCompiler) Dialect() DatabaseDialectType {
	return DatabaseDialectPostgreSQL
}

func (c *PostgresCompiler) Placeholder(paramIndex int) string {
	return fmt.Sprintf("$%d", paramIndex)
}

func (c *PostgresCompiler) QuoteIdentifier(name string) string {
	return "\"" + name + "\""
}

func (c *PostgresCompiler) CompilePagination(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	if offset <= 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func (c *PostgresCompiler) CompileSearch(table string, fields []string, limit int) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(fields) == 0 {
		return buildSelectWithoutWhere(quotedTable, c.CompilePagination(limit, 0))
	}

	whereClauses := make([]string, 0, len(fields))
	for i, f := range fields {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", c.QuoteIdentifier(f), c.Placeholder(i+1)))
	}

	whereSql := strings.Join(whereClauses, " AND ")
	pagination := c.CompilePagination(limit, 0)
	return buildSelectWithWhere(quotedTable, whereSql, pagination)
}

// MySQLCompiler compiles queries for MySQL and MariaDB databases.
type MySQLCompiler struct {
	dialect DatabaseDialectType
}

func (c *MySQLCompiler) Dialect() DatabaseDialectType {
	return c.dialect
}

func (c *MySQLCompiler) Placeholder(paramIndex int) string {
	return "?"
}

func (c *MySQLCompiler) QuoteIdentifier(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (c *MySQLCompiler) CompilePagination(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	if offset <= 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func (c *MySQLCompiler) CompileSearch(table string, fields []string, limit int) string {
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

// MSSQLCompiler compiles queries for Microsoft SQL Server databases.
type MSSQLCompiler struct{}

func (c *MSSQLCompiler) Dialect() DatabaseDialectType {
	return DatabaseDialectMSSQL
}

func (c *MSSQLCompiler) Placeholder(paramIndex int) string {
	return fmt.Sprintf("@p%d", paramIndex)
}

func (c *MSSQLCompiler) QuoteIdentifier(name string) string {
	return fmt.Sprintf("[%s]", name)
}

func (c *MSSQLCompiler) CompilePagination(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)
}

func (c *MSSQLCompiler) CompileSearch(table string, fields []string, limit int) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(fields) == 0 {
		return buildSelectWithoutWhere(quotedTable, c.CompilePagination(limit, 0))
	}

	whereClauses := make([]string, 0, len(fields))
	for i, f := range fields {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", c.QuoteIdentifier(f), c.Placeholder(i+1)))
	}

	whereSql := strings.Join(whereClauses, " AND ")
	pagination := c.CompilePagination(limit, 0)
	return buildSelectWithWhere(quotedTable, whereSql, pagination)
}

// OracleCompiler compiles queries for Oracle databases.
type OracleCompiler struct{}

func (c *OracleCompiler) Dialect() DatabaseDialectType {
	return DatabaseDialectOracle
}

func (c *OracleCompiler) Placeholder(paramIndex int) string {
	return fmt.Sprintf(":%d", paramIndex)
}

func (c *OracleCompiler) QuoteIdentifier(name string) string {
	return "\"" + name + "\""
}

func (c *OracleCompiler) CompilePagination(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)
}

func (c *OracleCompiler) CompileSearch(table string, fields []string, limit int) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(fields) == 0 {
		return buildSelectWithoutWhere(quotedTable, c.CompilePagination(limit, 0))
	}

	whereClauses := make([]string, 0, len(fields))
	for i, f := range fields {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", c.QuoteIdentifier(f), c.Placeholder(i+1)))
	}

	whereSql := strings.Join(whereClauses, " AND ")
	pagination := c.CompilePagination(limit, 0)
	return buildSelectWithWhere(quotedTable, whereSql, pagination)
}

// MongoDBCompiler compiles queries into MongoDB BSON representation.
type MongoDBCompiler struct{}

func (c *MongoDBCompiler) Dialect() DatabaseDialectType {
	return DatabaseDialectMongoDB
}

func (c *MongoDBCompiler) Placeholder(paramIndex int) string {
	return "$param"
}

func (c *MongoDBCompiler) QuoteIdentifier(name string) string {
	return name
}

func (c *MongoDBCompiler) CompilePagination(limit, offset int) string {
	return fmt.Sprintf(`{"$limit": %d, "$skip": %d}`, limit, offset)
}

func (c *MongoDBCompiler) CompileSearch(table string, fields []string, limit int) string {
	if len(fields) == 0 {
		return fmt.Sprintf("db.%s.find({}).limit(%d)", table, limit)
	}
	quotedFields := make([]string, 0, len(fields))
	for _, f := range fields {
		quotedFields = append(quotedFields, `"`+f+`": $param`)
	}
	filter := strings.Join(quotedFields, ", ")
	return fmt.Sprintf("db.%s.find({%s}).limit(%d)", table, filter, limit)
}
