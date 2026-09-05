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

func (c *PostgresCompiler) CompileCreateView(name string, selectSql string) string {
	cleanSql := strings.TrimRight(strings.TrimSpace(selectSql), ";")
	return fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s;", c.QuoteIdentifier(name), cleanSql)
}

func (c *PostgresCompiler) CompileDropView(name string) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s CASCADE;", c.QuoteIdentifier(name))
}

func (c *PostgresCompiler) CompileFunctionCall(name string, argCount int) string {
	if argCount <= 0 {
		return fmt.Sprintf("SELECT %s();", name)
	}
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf("SELECT %s(%s);", name, strings.Join(placeholders, ", "))
}

func (c *PostgresCompiler) CompileCount(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(field) == 0 {
		return fmt.Sprintf("SELECT COUNT(*) FROM %s;", quotedTable)
	}
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1;", quotedTable, c.QuoteIdentifier(field))
}

func (c *PostgresCompiler) CompileDelete(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	return fmt.Sprintf("DELETE FROM %s WHERE %s = $1;", quotedTable, c.QuoteIdentifier(field))
}

func (c *PostgresCompiler) CompileInspectColumns(tableOrView string) string {
	return "SELECT column_name FROM information_schema.columns WHERE table_name = $1;"
}

func (c *PostgresCompiler) CompileInspectViewExists(viewName string) string {
	return "SELECT 1 FROM information_schema.views WHERE table_name = $1;"
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

func (c *MySQLCompiler) CompileCreateView(name string, selectSql string) string {
	cleanSql := strings.TrimRight(strings.TrimSpace(selectSql), ";")
	return fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s;", c.QuoteIdentifier(name), cleanSql)
}

func (c *MySQLCompiler) CompileDropView(name string) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s;", c.QuoteIdentifier(name))
}

func (c *MySQLCompiler) CompileFunctionCall(name string, argCount int) string {
	if argCount <= 0 {
		return fmt.Sprintf("SELECT %s();", name)
	}
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}
	return fmt.Sprintf("SELECT %s(%s);", name, strings.Join(placeholders, ", "))
}

func (c *MySQLCompiler) CompileCount(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(field) == 0 {
		return fmt.Sprintf("SELECT COUNT(*) FROM %s;", quotedTable)
	}
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?;", quotedTable, c.QuoteIdentifier(field))
}

func (c *MySQLCompiler) CompileDelete(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	return fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", quotedTable, c.QuoteIdentifier(field))
}

func (c *MySQLCompiler) CompileInspectColumns(tableOrView string) string {
	return "SELECT column_name FROM information_schema.columns WHERE table_name = ?;"
}

func (c *MySQLCompiler) CompileInspectViewExists(viewName string) string {
	return "SELECT 1 FROM information_schema.views WHERE table_name = ?;"
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
	return "[" + name + "]"
}

func (c *MSSQLCompiler) CompilePagination(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	cleanOffset := offset
	if cleanOffset < 0 {
		cleanOffset = 0
	}
	cleanLimit := limit
	if cleanLimit <= 0 {
		cleanLimit = 1000
	}
	return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", cleanOffset, cleanLimit)
}

func (c *MSSQLCompiler) CompileSearch(table string, fields []string, limit int) string {
	quotedTable := c.QuoteIdentifier(table)
	pagination := c.CompilePagination(limit, 0)
	if len(fields) == 0 {
		return buildMssqlSelect(quotedTable, "", pagination)
	}

	whereClauses := make([]string, 0, len(fields))
	for i, f := range fields {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", c.QuoteIdentifier(f), c.Placeholder(i+1)))
	}

	whereSql := strings.Join(whereClauses, " AND ")
	return buildMssqlSelect(quotedTable, whereSql, pagination)
}

func buildMssqlSelect(quotedTable, whereSql, pagination string) string {
	base := "SELECT * FROM " + quotedTable
	if len(whereSql) > 0 {
		base += " WHERE " + whereSql
	}
	if len(pagination) == 0 {
		return base + ";"
	}
	return base + " ORDER BY (SELECT NULL) " + pagination + ";"
}

func (c *MSSQLCompiler) CompileCreateView(name string, selectSql string) string {
	cleanSql := strings.TrimRight(strings.TrimSpace(selectSql), ";")
	return fmt.Sprintf("CREATE OR ALTER VIEW %s AS %s;", c.QuoteIdentifier(name), cleanSql)
}

func (c *MSSQLCompiler) CompileDropView(name string) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s;", c.QuoteIdentifier(name))
}

func (c *MSSQLCompiler) CompileFunctionCall(name string, argCount int) string {
	if argCount <= 0 {
		return fmt.Sprintf("SELECT %s();", name)
	}
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = fmt.Sprintf("@p%d", i+1)
	}
	return fmt.Sprintf("SELECT %s(%s);", name, strings.Join(placeholders, ", "))
}

func (c *MSSQLCompiler) CompileCount(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(field) == 0 {
		return fmt.Sprintf("SELECT COUNT(*) FROM %s;", quotedTable)
	}
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = @p1;", quotedTable, c.QuoteIdentifier(field))
}

func (c *MSSQLCompiler) CompileDelete(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	return fmt.Sprintf("DELETE FROM %s WHERE %s = @p1;", quotedTable, c.QuoteIdentifier(field))
}

func (c *MSSQLCompiler) CompileInspectColumns(tableOrView string) string {
	return "SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = @p1;"
}

func (c *MSSQLCompiler) CompileInspectViewExists(viewName string) string {
	return "SELECT 1 FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_NAME = @p1;"
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
	cleanOffset := offset
	if cleanOffset < 0 {
		cleanOffset = 0
	}
	cleanLimit := limit
	if cleanLimit <= 0 {
		cleanLimit = 1000
	}
	return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", cleanOffset, cleanLimit)
}

func (c *OracleCompiler) CompileSearch(table string, fields []string, limit int) string {
	quotedTable := c.QuoteIdentifier(table)
	pagination := c.CompilePagination(limit, 0)
	if len(fields) == 0 {
		return buildSelectWithoutWhere(quotedTable, pagination)
	}

	whereClauses := make([]string, 0, len(fields))
	for i, f := range fields {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", c.QuoteIdentifier(f), c.Placeholder(i+1)))
	}

	whereSql := strings.Join(whereClauses, " AND ")
	return buildSelectWithWhere(quotedTable, whereSql, pagination)
}

func (c *OracleCompiler) CompileCreateView(name string, selectSql string) string {
	cleanSql := strings.TrimRight(strings.TrimSpace(selectSql), ";")
	return fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s;", c.QuoteIdentifier(name), cleanSql)
}

func (c *OracleCompiler) CompileDropView(name string) string {
	return fmt.Sprintf("DROP VIEW %s;", c.QuoteIdentifier(name))
}

func (c *OracleCompiler) CompileFunctionCall(name string, argCount int) string {
	if argCount <= 0 {
		return fmt.Sprintf("SELECT %s() FROM DUAL;", name)
	}
	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = fmt.Sprintf(":%d", i+1)
	}
	return fmt.Sprintf("SELECT %s(%s) FROM DUAL;", name, strings.Join(placeholders, ", "))
}

func (c *OracleCompiler) CompileCount(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	if len(field) == 0 {
		return fmt.Sprintf("SELECT COUNT(*) FROM %s;", quotedTable)
	}
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = :1;", quotedTable, c.QuoteIdentifier(field))
}

func (c *OracleCompiler) CompileDelete(table string, field string) string {
	quotedTable := c.QuoteIdentifier(table)
	return fmt.Sprintf("DELETE FROM %s WHERE %s = :1;", quotedTable, c.QuoteIdentifier(field))
}

func (c *OracleCompiler) CompileInspectColumns(tableOrView string) string {
	return "SELECT COLUMN_NAME FROM ALL_TAB_COLUMNS WHERE TABLE_NAME = :1;"
}

func (c *OracleCompiler) CompileInspectViewExists(viewName string) string {
	return "SELECT 1 FROM ALL_VIEWS WHERE VIEW_NAME = :1;"
}

// MongoDBCompiler compiles queries into MongoDB representation.
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

func (c *MongoDBCompiler) CompileCreateView(name string, selectSql string) string {
	return fmt.Sprintf("db.createView('%s', '%s', [])", name, name)
}

func (c *MongoDBCompiler) CompileDropView(name string) string {
	return fmt.Sprintf("db.%s.drop()", name)
}

func (c *MongoDBCompiler) CompileFunctionCall(name string, argCount int) string {
	return fmt.Sprintf("%s()", name)
}

func (c *MongoDBCompiler) CompileCount(table string, field string) string {
	return fmt.Sprintf("db.%s.countDocuments()", table)
}

func (c *MongoDBCompiler) CompileDelete(table string, field string) string {
	return fmt.Sprintf("db.%s.deleteMany({})", table)
}

func (c *MongoDBCompiler) CompileInspectColumns(tableOrView string) string {
	return ""
}

func (c *MongoDBCompiler) CompileInspectViewExists(viewName string) string {
	return ""
}
