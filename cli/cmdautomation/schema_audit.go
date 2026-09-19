package cmdautomation

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	reCreateTable = regexp.MustCompile(`(?is)CREATE\s+TABLE(?:\s+IF\s+NOT\s+EXISTS)?\s+["` + "`" + `]?([A-Za-z0-9_]+)["` + "`" + `]?\s*\((.*?)\)(?:;|\s*` + "`" + `)`)
)

// RunSchemaAudit inspects database schemas and verifies naming conventions.
func RunSchemaAudit(opts SchemaAuditOptions) SchemaAuditResultMonad {
	start := time.Now()
	var res SchemaAuditResult
	if opts.DbPath != "" {
		auditDbFile(opts.DbPath, opts.IsStrict, &res)
	}
	scanDir := opts.Dir
	if opts.DbPath == "" && scanDir == "" {
		scanDir = "."
	}
	if scanDir != "" {
		auditSourceFiles(scanDir, opts.IsStrict, &res)
	}
	res.Duration = time.Since(start)
	res.IsClean = len(res.Violations) == 0
	return result.Ok(res)
}

func auditDbFile(dbPath string, isStrict bool, res *SchemaAuditResult) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return
	}
	defer db.Close()
	tables, errApp := inspectDbTables(dbPath)
	if errApp != nil {
		return
	}
	res.ScannedFiles++
	for _, t := range tables {
		res.TableCount++
		inspectTableConventions(dbPath, t, isStrict, res)
	}
}

func auditSourceFiles(dir string, isStrict bool, res *SchemaAuditResult) {
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return handleSourceDirWalk(info, err)
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != ".go" && ext != ".sql" {
			return nil
		}
		res.ScannedFiles++
		scanFileSchema(filepath.ToSlash(path), isStrict, res)
		return nil
	})
}

func handleSourceDirWalk(info os.FileInfo, err error) error {
	if err != nil || info == nil {
		return nil
	}
	if isTopologyExcludedDir(info.Name()) {
		return filepath.SkipDir
	}
	return nil
}

func scanFileSchema(filePath string, isStrict bool, res *SchemaAuditResult) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	matches := reCreateTable.FindAllStringSubmatch(string(content), -1)
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		tableName := m[1]
		body := m[2]
		res.TableCount++
		auditParsedTable(filePath, tableName, body, isStrict, res)
	}
}

func auditParsedTable(file, table, body string, isStrict bool, res *SchemaAuditResult) {
	if !isPascalCase(table) {
		res.Violations = append(res.Violations, SchemaViolation{
			File: file, Table: table, Issue: "Table name is not PascalCase", Severity: "error",
		})
	}
	checkTablePrimaryKey(file, table, body, res)
}

func checkTablePrimaryKey(file, table, body string, res *SchemaAuditResult) {
	if !strings.Contains(strings.ToUpper(body), "PRIMARY KEY") {
		return
	}
	if strings.HasSuffix(table, "Metadata") {
		return
	}
	expectedPk := strings.ToLower(table + "Id")
	expectedId := strings.ToLower(table + "_id")
	bodyLower := strings.ToLower(body)
	if !strings.Contains(bodyLower, expectedPk) && !strings.Contains(bodyLower, expectedId) {
		res.Violations = append(res.Violations, SchemaViolation{
			File:     file,
			Table:    table,
			Issue:    fmt.Sprintf("Primary key does not match %sId or %s_id", table, table),
			Severity: "warning",
		})
	}
}

func inspectTableConventions(dbFile string, t TableSchemaInfo, isStrict bool, res *SchemaAuditResult) {
	if !isPascalCase(t.Name) {
		res.Violations = append(res.Violations, SchemaViolation{
			File: dbFile, Table: t.Name, Issue: "Table name is not PascalCase", Severity: "error",
		})
	}
	for _, col := range t.Columns {
		checkColumnConventions(dbFile, t.Name, col, isStrict, res)
	}
}

func checkColumnConventions(file, table string, col TableColumnInfo, isStrict bool, res *SchemaAuditResult) {
	colUpper := strings.ToUpper(col.Type)
	if !strings.Contains(colUpper, "BOOL") {
		return
	}
	colLower := strings.ToLower(col.Name)
	hasAffirmative := strings.HasPrefix(colLower, "is") || strings.HasPrefix(colLower, "has") || strings.HasPrefix(colLower, "can")
	if !hasAffirmative {
		res.Violations = append(res.Violations, SchemaViolation{
			File:     file,
			Table:    table,
			Issue:    fmt.Sprintf("Boolean column '%s' lacks affirmative prefix (is*, has*, can*)", col.Name),
			Severity: "warning",
		})
	}
}

func isPascalCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	runes := []rune(s)
	if !unicode.IsUpper(runes[0]) {
		return false
	}
	for _, r := range runes {
		if r == '_' || r == '-' || (!unicode.IsLetter(r) && !unicode.IsDigit(r)) {
			return false
		}
	}
	return true
}
