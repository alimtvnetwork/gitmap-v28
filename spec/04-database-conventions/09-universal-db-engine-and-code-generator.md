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

### 5.2 SQL Operator Enum (`SqlOperator`)
The `SqlOperator` enum governs comparison operators across queries with gold-standard receiver methods:
```go
type SqlOperator string

const (
    SqlOpEqual              SqlOperator = "="
    SqlOpNotEqual           SqlOperator = "!="
    SqlOpLessThan           SqlOperator = "<"
    SqlOpLessThanOrEqual    SqlOperator = "<="
    SqlOpGreaterThan        SqlOperator = ">"
    SqlOpGreaterThanOrEqual SqlOperator = ">="
    SqlOpLike               SqlOperator = "LIKE"
    SqlOpIn                 SqlOperator = "IN"
)

// Scoped registry
var SqlOperators = sqlOperatorRegistry{
    Equal:       SqlOpEqual,
    NotEqual:    SqlOpNotEqual,
    LessThan:    SqlOpLessThan,
    GreaterThan: SqlOpGreaterThan,
    Like:        SqlOpLike,
    In:          SqlOpIn,
}
```

### 5.3 Fluent QueryBuilder, Joins, Error Guards & View Evolution
The `QueryBuilder[T, F]` provides fluent, type-safe query construction with scoped join sub-structs, error tracking, substring locating, ad-hoc CTE views, sorting, pagination, and automated database-backed view evolution:
```go
// 1. Scoped JoinBuilder with zero magic strings and error tracking
recordsRes := repo.Query().
    Select(PipelineRunRecordDb.RunId, PipelineRunRecordDb.WorkflowName).
    InnerJoin(PipelineErrorRecordTable).
        Select(PipelineErrorRecordDb.ErrorText, PipelineErrorRecordDb.StepName).
        OnField(PipelineRunRecordDb.RunId, SqlOperators.Equal, PipelineErrorRecordDb.RunId).
    WhereOp(PipelineRunRecordDb.Status, SqlOperators.Equal, "failed").
    GroupBy(PipelineRunRecordDb.WorkflowName).
    HavingCount(SqlOperators.GreaterThan, 3).
    OrderByDesc(PipelineRunRecordDb.WorkflowName).
    Limit(10).
    FindAll(ctx) // ListResult[T]

// 2. SQL compilation with thread-safe cache & QueryHash
compRes := repo.Query().
    InnerJoin(PipelineErrorRecordTable).
        OnField(PipelineRunRecordDb.RunId, SqlOperators.Equal, PipelineErrorRecordDb.RunId).
    Compile() // CompiledQueryResult (result.Result[CompiledQuery])

if compRes.IsSuccess() {
    cq := compRes.Value
    fmt.Printf("SQL: %s\nHash: %s\n", cq.SQL, cq.QueryHash)
}

// 3. Automated database-backed view evolution via QueryHash (__gitmap_view_meta)
// Zero manual column parameters required; 0 DDL on identical hash; pre-validates via EXPLAIN
viewRes := repo.Query().
    Select(PipelineRunRecordDb.WorkflowName).
    InnerJoin(PipelineErrorRecordTable).
        Select(PipelineErrorRecordDb.ErrorText).
        OnField(PipelineRunRecordDb.RunId, SqlOperators.Equal, PipelineErrorRecordDb.RunId).
    CreateViewOrUseView(ctx, "ActiveErrorsView") // BoolResult

// 4. Domain business repository pattern
pipelineRepo := pipelinedb.NewPipelineRepository(db)
runRes := pipelineRepo.GetRunById(ctx, 101)                 // EntityResult[PipelineRunRecord]
recentRuns := pipelineRepo.GetRecentRuns(ctx, "owner/repo", 20) // ListResult[PipelineRunRecord]
activeView := pipelineRepo.EnsureActiveErrorsView(ctx)      // BoolResult
```

### 5.4 Automatic Typed DbRepo & Safe Row Scanners Generation Pattern
The database code generator (`03-ai-scripts/30-db-struct-enum-generator.py`) automatically generates typed repository structs and null-safe row scanners for any Go model struct:
1. **Safe Row Scanners (`Scan{StructName}`)**:
   - Uses `dbengine.Scan*` helpers (`ScanString`, `ScanInt`, `ScanInt64`, `ScanUint64`, `ScanUint`, `ScanBool`, `ScanFloat64`).
   - Completely null-safe: handles SQLite `nil` / NULL values without panics or driver conversion errors.
   - Converts integer column types (`int64`, `int`) into struct field types (such as `uint64` or `bool`).
2. **Typed Repository Struct (`{StructName}DbRepo` & alias `{ShortName}DbRepo`)**:
   - Generates typed repository struct wrapping generic `Repository[T, F]`.
   - Provides alias for concise business usage (e.g., `PipelineRunDbRepo = PipelineRunRecordDbRepo`).
   - Generates constructor `New{ShortName}DbRepo(db *dbengine.DbWrapper)`.
   - Exposes standard typed methods returning Result envelopes:
     - `FindAll(ctx context.Context) dbengine.ListResult[T]`
     - `First(ctx context.Context) dbengine.EntityResult[T]`
     - `Count(ctx context.Context) dbengine.Int64Result`
     - `Query() *dbengine.QueryBuilder[T, F]`
     - `QueryBare() *dbengine.QueryBuilder[T, F]`
     - `Db() *dbengine.DbWrapper`
     - `Repo() *dbengine.Repository[T, F]`
3. **Domain Business Repository Integration**:
   - Domain repositories (e.g. `PipelineRepository`) embed `*{StructName}DbRepo` directly:
     ```go
     type PipelineRepository struct {
         *PipelineRunRecordDbRepo
     }
     ```
   - Automatically inherits all standard CRUD operations while cleanly adding domain business methods (`GetRecentRuns`, `EnsureActiveErrorsView`).

---

## 6. External Migration Runner Workflow

1. Standalone Python runner `03-ai-scripts/26-db-migration-runner.py` inspects model struct definitions and migration SQL files.
2. Applies migrations with `CREATE TABLE IF NOT EXISTS` and `PRAGMA table_info` schema detection.
3. Automatically triggers the enum generator and marks generated output as `// Code generated by gitmap db generate. DO NOT EDIT.`
