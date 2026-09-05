package sqlite

import (
	"fmt"
	"strings"
)

// CompileCreateView builds a CREATE VIEW IF NOT EXISTS statement.
func (c *Compiler) CompileCreateView(name string, selectSql string) string {
	cleanSql := strings.TrimRight(strings.TrimSpace(selectSql), ";")
	return fmt.Sprintf("CREATE VIEW IF NOT EXISTS %s AS %s;", c.QuoteIdentifier(name), cleanSql)
}

// CompileDropView builds a DROP VIEW IF EXISTS statement.
func (c *Compiler) CompileDropView(name string) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s;", c.QuoteIdentifier(name))
}

// CompileAdHocCTE wraps a main query with a Common Table Expression (CTE) view.
func (c *Compiler) CompileAdHocCTE(viewName string, subQuery string, mainQuery string) string {
	cleanSub := strings.TrimRight(strings.TrimSpace(subQuery), ";")
	cleanMain := strings.TrimRight(strings.TrimSpace(mainQuery), ";")
	return fmt.Sprintf("WITH %s AS (%s) %s;", c.QuoteIdentifier(viewName), cleanSub, cleanMain)
}

// CompileInspectColumns returns the PRAGMA query to inspect columns of a table or view.
func (c *Compiler) CompileInspectColumns(tableOrView string) string {
	return fmt.Sprintf("PRAGMA table_info(%s);", c.QuoteIdentifier(tableOrView))
}

// CompileInspectViewExists returns a query to check if a view exists in sqlite_master.
func (c *Compiler) CompileInspectViewExists(viewName string) string {
	return "SELECT 1 FROM sqlite_master WHERE type = 'view' AND name = ?;"
}
