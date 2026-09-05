# Universal Database Engine, Dialect Compiler & Field Enum Generator

**Version:** 1.0.0  
**Updated:** 2026-09-05  
**Category:** Database Architecture / Multi-Dialect Query Engine / Code Generation  

---

## 1. Overview & Acronym Conventions

### 1.1 Strict PascalCase Acronym Standards
In accordance with repository-wide naming conventions:
- Acronyms MUST be PascalCased with only the first letter capitalized: `Id`, `Db`, `Url`, `Api`.
- All uppercase acronym clusters (`ID`, `DB`, `URL`, `API`) are strictly forbidden across method names, parameter names, struct fields, and variables.
  - ✅ `runId`, `databaseId`, `pipeDb`, `PipelineSplitDb`
  - ❌ `runID`, `databaseID`, `pipeDB`, `PipelineSplitDB`

### 1.2 Architectural Objectives
1. **Unified Error Contract**: All database operations return `(*T, *AppError)` or `([]T, *AppError)`. Zero swallowed errors.
2. **Multi-Dialect Compilation**: Dialects compile parameterized queries across SQLite, PostgreSQL, MySQL, MariaDB, MSSQL, Oracle, and MongoDB.
3. **Zero Magic Strings**: Model struct fields are automatically extracted into strongly-typed enums (`*FieldType`) via a code generator.
4. **Generic Search API**: Generic repositories (`Repository[T, F]`) handle 1-parameter, 2-parameter, single-record, and paginated searches with type safety.
5. **External Migration Workflow**: Standalone runner executes migrations outside the Go binary.

---

## 2. Multi-Dialect Engine & Compilation

### 2.1 Database Dialect Enum
```go
type DatabaseDialectType string

const (
    DatabaseDialectSQLite     DatabaseDialectType = "sqlite"
    DatabaseDialectPostgreSQL DatabaseDialectType = "postgres"
    DatabaseDialectMySQL      DatabaseDialectType = "mysql"
    DatabaseDialectMariaDB    DatabaseDialectType = "mariadb"
    DatabaseDialectMSSQL      DatabaseDialectType = "mssql"
    DatabaseDialectOracle     DatabaseDialectType = "oracle"
    DatabaseDialectMongoDB    DatabaseDialectType = "mongodb"
)
```

### 2.2 Dialect Compiler Contract
```go
type DialectCompiler interface {
    Dialect() DatabaseDialectType
    Placeholder(paramIndex int) string
    QuoteIdentifier(name string) string
    CompilePagination(query string, limit, offset int) string
    CompileSearch(table string, fields []string, limit int) (string, []string)
}
```

### 2.3 Syntax Matrix
| Dialect | Parameter Placeholder | Identifier Quoting | Limit / Pagination Syntax |
| :--- | :--- | :--- | :--- |
| **SQLite** | `?` | `"ColumnName"` | `LIMIT {limit} OFFSET {offset}` |
| **MySQL / MariaDB** | `?` | `` `ColumnName` `` | `LIMIT {limit} OFFSET {offset}` |
| **PostgreSQL** | `$1, $2, ...` | `"ColumnName"` | `LIMIT {limit} OFFSET {offset}` |
| **MSSQL Server** | `@p1, @p2, ...` | `[ColumnName]` | `OFFSET {offset} ROWS FETCH NEXT {limit} ROWS ONLY` |
| **Oracle** | `:1, :2, ...` | `"ColumnName"` | `OFFSET {offset} ROWS FETCH NEXT {limit} ROWS ONLY` |
| **MongoDB** | BSON Filter | Document Key | `{"$limit": N, "$skip": M}` |

---

## 3. Unified DbWrapper & AppError Return Contract

```go
type DbWrapper struct {
    conn     *sql.DB
    dialect  DatabaseDialectType
    compiler DialectCompiler
}

func (w *DbWrapper) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, *apperror.AppError)
func (w *DbWrapper) Query(ctx context.Context, query string, args ...any) (*sql.Rows, *apperror.AppError)
func (w *DbWrapper) Exec(ctx context.Context, query string, args ...any) (sql.Result, *apperror.AppError)
func (w *DbWrapper) WithTransaction(ctx context.Context, fn func(tx *TxWrapper) *apperror.AppError) *apperror.AppError
```

---

## 4. Struct-to-Table Model & CLI Enum Generator

### 4.1 Model Struct Standard
Model structs represent database tables and use tags for column metadata:
```go
type PipelineErrorLog struct {
    PipelineErrorLogId int64  `db:"PipelineErrorLogId" dbtype:"INTEGER PRIMARY KEY AUTOINCREMENT"`
    RunId              int64  `db:"RunId" dbtype:"INTEGER NOT NULL"`
    RepoSlug           string `db:"RepoSlug" dbtype:"TEXT NOT NULL"`
    WorkflowName       string `db:"WorkflowName" dbtype:"TEXT NOT NULL"`
    StepName           string `db:"StepName" dbtype:"TEXT NOT NULL"`
    ErrorText          string `db:"ErrorText" dbtype:"TEXT NOT NULL"`
    RawLogs            string `db:"RawLogs" dbtype:"TEXT NULL"`
    Notes              string `db:"Notes" dbtype:"TEXT NULL"`
    Comments           string `db:"Comments" dbtype:"TEXT NULL"`
    CreatedAt          string `db:"CreatedAt" dbtype:"TEXT NOT NULL"`
}
```

### 4.2 Auto-Generated Field Enums
The code generator produces typed string enums eliminating all magic strings:
```go
type PipelineErrorLogFieldType string

const (
    PipelineErrorLogFieldPipelineErrorLogId PipelineErrorLogFieldType = "PipelineErrorLogId"
    PipelineErrorLogFieldRunId              PipelineErrorLogFieldType = "RunId"
    PipelineErrorLogFieldRepoSlug           PipelineErrorLogFieldType = "RepoSlug"
    PipelineErrorLogFieldWorkflowName       PipelineErrorLogFieldType = "WorkflowName"
    PipelineErrorLogFieldStepName           PipelineErrorLogFieldType = "StepName"
    PipelineErrorLogFieldErrorText          PipelineErrorLogFieldType = "ErrorText"
    PipelineErrorLogFieldRawLogs            PipelineErrorLogFieldType = "RawLogs"
    PipelineErrorLogFieldNotes              PipelineErrorLogFieldType = "Notes"
    PipelineErrorLogFieldComments           PipelineErrorLogFieldType = "Comments"
    PipelineErrorLogFieldCreatedAt          PipelineErrorLogFieldType = "CreatedAt"
)
```

---

## 5. Generic Query Builder & Search API

```go
type Repository[T any, F ~string] struct {
    db        *DbWrapper
    tableName string
    scanner   ModelScanner[T]
}

// 1-parameter search returning limit 1
func (r *Repository[T, F]) First(ctx context.Context, field F, value any) (*T, *apperror.AppError)

// 1-parameter search returning up to limit N
func (r *Repository[T, F]) FindBy(ctx context.Context, field F, value any, limit int) ([]T, *apperror.AppError)

// 2-parameter search returning up to limit N
func (r *Repository[T, F]) FindBy2(ctx context.Context, field1 F, val1 any, field2 F, val2 any, limit int) ([]T, *apperror.AppError)
```

---

## 6. External Migration Runner Workflow

1. Standalone Python runner `03-ai-scripts/26-db-migration-runner.py` inspects model struct definitions and migration SQL files.
2. Applies migrations with `CREATE TABLE IF NOT EXISTS` and `PRAGMA table_info` schema detection.
3. Automatically triggers the enum generator and marks generated output as `// Code generated by gitmap db generate. DO NOT EDIT.`
