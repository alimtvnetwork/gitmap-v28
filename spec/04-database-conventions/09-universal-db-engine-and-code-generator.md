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

### 2.1 Database Type Enum (`DbType`)
The `DbType` enum is the canonical database engine selector and includes gold-standard receiver methods (`Name()`, `String()`, `Value()`, `IsCompare()`):
```go
type DbType string

const (
    DbSQLite     DbType = "sqlite"
    DbPostgreSQL DbType = "postgres"
    DbMySQL      DbType = "mysql"
    DbMariaDB    DbType = "mariadb"
    DbMSSQL      DbType = "mssql"
    DbOracle     DbType = "oracle"
    DbMongoDB    DbType = "mongodb"
)

func (d DbType) Name() string { return string(d) }
func (d DbType) String() string { return string(d) }
func (d DbType) Value() string { return string(d) }
func (d DbType) IsCompare(target any) bool {
    switch v := target.(type) {
    case DbType:
        return d == v
    case string:
        return string(d) == v
    case fmt.Stringer:
        return string(d) == v.String()
    default:
        return false
    }
}

// dbTypeValidMap provides O(1) map validation for database types.
var dbTypeValidMap = map[DbType]bool{
    DbSQLite:     true,
    DbPostgreSQL: true,
    DbMySQL:      true,
    DbMariaDB:    true,
    DbMSSQL:      true,
    DbOracle:     true,
    DbMongoDB:    true,
}

func (d DbType) IsCompare(target DbType) bool { return d == target }
func (d DbType) IsEnum() bool                 { return dbTypeValidMap[d] }
func (d DbType) IsSQLite() bool               { return d == DbSQLite }
func (d DbType) IsPostgreSQL() bool           { return d == DbPostgreSQL }

// JSON conversion methods with AppError
func (d DbType) MarshalJSON() ([]byte, error) { return json.Marshal(string(d)) }
func (d *DbType) UnmarshalJSON(data []byte) error { ... }
func (d DbType) ToJSON() (string, *apperror.AppError) { ... }
func (d *DbType) FromJSON(s string) *apperror.AppError { ... }

// Scoped registry for zero-magic strings
type dbTypeRegistry struct {
    SQLite     DbType
    PostgreSQL DbType
    Postgres   DbType
    MySQL      DbType
    MariaDB    DbType
    MSSQL      DbType
    Oracle     DbType
    MongoDB    DbType
}

func (r dbTypeRegistry) All() []DbType { ... }
func (r dbTypeRegistry) Names() []string { ... }
func (r dbTypeRegistry) IsEnum(target DbType) bool { return dbTypeValidMap[target] }
func (r dbTypeRegistry) IsSQLite(target DbType) bool { return target == r.SQLite }
func (r dbTypeRegistry) IsPostgreSQL(target DbType) bool { return target == r.PostgreSQL }
func (r dbTypeRegistry) ToJSON() (string, *apperror.AppError) { ... }

var DbTypes = dbTypeRegistry{
    SQLite:     DbSQLite,
    PostgreSQL: DbPostgreSQL,
    Postgres:   DbPostgreSQL,
    MySQL:      DbMySQL,
    MariaDB:    DbMariaDB,
    MSSQL:      DbMSSQL,
    Oracle:     DbOracle,
    MongoDB:    DbMongoDB,
}

type DatabaseDialectType = DbType
```

### 2.2 Dialect Compiler Contract
```go
type DialectCompiler interface {
    Dialect() DbType
    Placeholder(paramIndex int) string
    QuoteIdentifier(name string) string
    CompilePagination(limit, offset int) string
    CompileSearch(table string, fields []string, limit int) string
    CompileCreateView(name string, selectSql string) string
    CompileDropView(name string) string
    CompileFunctionCall(name string, argCount int) string
    CompileCount(table string, field string) string
    CompileDelete(table string, field string) string
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

### 2.4 Dedicated SQLite Engine Section (`gitmap/dbengine/sqlite/`)
Each database dialect implementation resides in its own package section under `gitmap/dbengine/<dialect>/`:
- `compiler.go`: SQLite-specific query compilation (`INSTR(field, ?) > 0` for `CompileLocate`, `LIMIT/OFFSET`, parameterized search, counts, deletes).
- `views.go`: View lifecycle (`CREATE VIEW IF NOT EXISTS`, `DROP VIEW IF EXISTS`, and ad-hoc CTE `CompileAdHocCTE` via `WITH "ViewName" AS (...)`).
- `functions.go`: Scalar function invocation and execution (`CompileFunctionCall`, `CompileScalarFunctionExpression`).

---

## 3. Unified DbWrapper & Result Envelope Contract

To ensure zero raw unwrapped tuples and preserve `*apperror.AppError` context without swallowing, data-layer query and execution methods return strongly typed Result envelopes:
```go
type Result[T any] struct {
    Value T
    Err   *apperror.AppError
}

type EntityResult[T any] = result.Result[*T]
type ListResult[T any]   = result.Result[[]T]
type Uint64Result        = result.Result[uint64]
type Int64Result         = result.Result[int64]
type StringResult        = result.Result[string]
type BoolResult          = result.Result[bool]
type RowsAffectedResult  = result.Result[int64]
```

### DbWrapper Operations
```go
type DbWrapper struct {
    conn     *sql.DB
    dialect  DbType
    compiler DialectCompiler
}

func (w *DbWrapper) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, *apperror.AppError)
func (w *DbWrapper) Query(ctx context.Context, query string, args ...any) (*sql.Rows, *apperror.AppError)
func (w *DbWrapper) Exec(ctx context.Context, query string, args ...any) (sql.Result, *apperror.AppError)
func (w *DbWrapper) ExecRowsAffected(ctx context.Context, query string, args ...any) RowsAffectedResult
func (w *DbWrapper) WithTransaction(ctx context.Context, fn func(tx *TxWrapper) *apperror.AppError) *apperror.AppError

// View & Function Management
func (w *DbWrapper) CreateView(ctx context.Context, name string, selectSql string) BoolResult
func (w *DbWrapper) DropView(ctx context.Context, name string) BoolResult
func (w *DbWrapper) CallFunction(ctx context.Context, name string, args ...any) StringResult
```

---

## 4. Struct-to-Table Model & CLI Enum Generator

### 4.1 Model Struct Standard
Model structs represent database tables and use tags for column metadata. Non-negative database IDs use `uint64`:
```go
type PipelineErrorLog struct {
    PipelineErrorLogId uint64 `db:"PipelineErrorLogId" dbtype:"INTEGER PRIMARY KEY AUTOINCREMENT"`
    RunId              uint64 `db:"RunId" dbtype:"INTEGER NOT NULL"`
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

### 4.2 Auto-Generated Field Enums & Scoped Registries
To avoid repeating prefixes across flat constants (e.g. avoiding repetitive `PipelineErrorLogFieldPipelineErrorLogId`), field enums are grouped into scoped database registries. Developers access fields via the canonical `<Model>Db.<FieldName>` pattern (e.g. `PipelineErrorLogDb.RunId`, `PipelineRunRecordDb.RunId`):
```go
type PipelineErrorLogFieldType string

func (e PipelineErrorLogFieldType) Name() string   { return string(e) }
func (e PipelineErrorLogFieldType) String() string { return string(e) }
func (e PipelineErrorLogFieldType) Value() string  { return string(e) }
func (e PipelineErrorLogFieldType) IsCompare(target PipelineErrorLogFieldType) bool {
    return e == target
}
func (e PipelineErrorLogFieldType) IsEnum() bool { return pipelineErrorLogValidMap[e] }
func (e PipelineErrorLogFieldType) IsRunId() bool { return e == PipelineErrorLogDb.RunId }
func (e PipelineErrorLogFieldType) IsRepoSlug() bool { return e == PipelineErrorLogDb.RepoSlug }

// JSON conversion methods with AppError
func (e PipelineErrorLogFieldType) MarshalJSON() ([]byte, error) { return json.Marshal(string(e)) }
func (e *PipelineErrorLogFieldType) UnmarshalJSON(data []byte) error { ... }
func (e PipelineErrorLogFieldType) ToJSON() (string, *apperror.AppError) { ... }
func (e *PipelineErrorLogFieldType) FromJSON(s string) *apperror.AppError { ... }

type pipelineErrorLogDbRegistry struct {
    PipelineErrorLogId PipelineErrorLogFieldType
    RunId              PipelineErrorLogFieldType
    RepoSlug           PipelineErrorLogFieldType
    WorkflowName       PipelineErrorLogFieldType
    StepName           PipelineErrorLogFieldType
    ErrorText          PipelineErrorLogFieldType
    RawLogs            PipelineErrorLogFieldType
    Notes              PipelineErrorLogFieldType
    Comments           PipelineErrorLogFieldType
    CreatedAt          PipelineErrorLogFieldType
}

func (r pipelineErrorLogDbRegistry) All() []PipelineErrorLogFieldType { ... }
func (r pipelineErrorLogDbRegistry) Names() []string { ... }
func (r pipelineErrorLogDbRegistry) IsEnum(target PipelineErrorLogFieldType) bool { return pipelineErrorLogValidMap[target] }
func (r pipelineErrorLogDbRegistry) IsRunId(target PipelineErrorLogFieldType) bool { return target == r.RunId }
func (r pipelineErrorLogDbRegistry) ToJSON() (string, *apperror.AppError) { ... }

// Canonical scoped access: PipelineErrorLogDb.RunId
var PipelineErrorLogDb = pipelineErrorLogDbRegistry{
    PipelineErrorLogId: "PipelineErrorLogId",
    RunId:              "RunId",
    RepoSlug:           "RepoSlug",
    WorkflowName:       "WorkflowName",
    StepName:           "StepName",
    ErrorText:          "ErrorText",
    RawLogs:            "RawLogs",
    Notes:              "Notes",
    Comments:           "Comments",
    CreatedAt:          "CreatedAt",
}

// O(1) valid enum map
var pipelineErrorLogValidMap = map[PipelineErrorLogFieldType]bool{
    PipelineErrorLogDb.PipelineErrorLogId: true,
    PipelineErrorLogDb.RunId:              true,
    PipelineErrorLogDb.RepoSlug:           true,
    // ...
}

// PipelineErrorLogField is an alias to PipelineErrorLogDb
var PipelineErrorLogField = PipelineErrorLogDb

const PipelineErrorLogTable = "PipelineErrorLog"
```

---

## 5. Generic Repository & Fluent Query Builder

All repository and query builder methods return wrapped `result.Result[T]` types to eliminate raw unhandled tuples:

### 5.1 Repository API
```go
type Repository[T any, F ~string] struct {
    db        *DbWrapper
    tableName string
    scanner   ModelScanner[T]
}

// Fluent Query Builder entry point
func (r *Repository[T, F]) Query() *QueryBuilder[T, F]

// 1-parameter search returning limit 1
func (r *Repository[T, F]) First(ctx context.Context, field F, value any) EntityResult[T]
func (r *Repository[T, F]) FindById(ctx context.Context, idField F, idValue any) EntityResult[T]

// Multi-parameter searches
func (r *Repository[T, F]) FindBy(ctx context.Context, field F, value any, limit int) ListResult[T]
func (r *Repository[T, F]) FindBy2(ctx context.Context, field1 F, val1 any, field2 F, val2 any, limit int) ListResult[T]
func (r *Repository[T, F]) FindAll(ctx context.Context, limit int) ListResult[T]

// Aggregate and mutation
func (r *Repository[T, F]) Count(ctx context.Context, field F, value any) Int64Result
func (r *Repository[T, F]) CountAll(ctx context.Context) Int64Result
func (r *Repository[T, F]) DeleteBy(ctx context.Context, field F, value any) RowsAffectedResult
```

### 5.2 Fluent QueryBuilder
The `QueryBuilder[T, F]` provides fluent query construction with substring locating, table joining, ad-hoc CTE views, sorting, and pagination:
```go
qb := repo.Query().
    Where(PipelineErrorLogDb.WorkflowName, "ci").
    Locate(PipelineErrorLogDb.ErrorText, "timeout").
    OrderByDesc(PipelineErrorLogDb.PipelineErrorLogId).
    Limit(10).
    Offset(0)

// Terminal executions returning wrapped results
firstRes := qb.First(ctx)       // EntityResult[T]
listRes  := qb.FindAll(ctx)     // ListResult[T]
countRes := qb.Count(ctx)       // Int64Result
delRes   := qb.Delete(ctx)      // RowsAffectedResult

// Joining and Ad-Hoc CTE Views
qb.Join("OtherTable", "OtherTable.RunId = PipelineErrorLog.RunId")
qb.LeftJoin("AuditLog", "AuditLog.RunId = PipelineErrorLog.RunId")
qb.WithView("ActiveErrors", "SELECT * FROM PipelineErrorLog WHERE ErrorText != ''")
```

---

## 6. External Migration Runner Workflow

1. Standalone Python runner `03-ai-scripts/26-db-migration-runner.py` inspects model struct definitions and migration SQL files.
2. Applies migrations with `CREATE TABLE IF NOT EXISTS` and `PRAGMA table_info` schema detection.
3. Automatically triggers the enum generator and marks generated output as `// Code generated by gitmap db generate. DO NOT EDIT.`
