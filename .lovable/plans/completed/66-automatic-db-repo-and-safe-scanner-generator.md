# 66 — Automatic Typed DbRepo & Safe Row Scanners Generator

**Status:** COMPLETED  
**Priority:** High  
**Category:** Database Architecture / Code Generation / Type-Safe Repositories  
**Completed:** 2026-09-05  

---

## 1. Problem Statement & Root Cause

### Problem Statement
Previously, database model structs generated column enums (`*FieldType`) and registries (`*DbRegistry`), but developers still had to write boilerplate row scanners and typed repository structs manually. Specifically:
1. Handwritten scanners (`ScanPipelineRunRecord`, `ScanPipelineErrorRecord`) duplicated struct fields across files.
2. Standard SQLite driver scanning caused errors or panics when handling `NULL` values or converting SQLite `int64` columns into Go `uint64` / `bool` types.
3. Repositories had to be instantiated manually via generic `dbengine.NewRepository[T, F]` and wrapped with boilerplate query methods (`FindAll`, `First`, `Count`, `Query`).

### Root Cause
The database code generator (`03-ai-scripts/30-db-struct-enum-generator.py`) only generated field enums and JSON serializers/deserializers; it lacked AST-driven repository code generation and null-safe scanner helpers.

---

## 2. Architecture & Design Decisions

### 2.1 Safe Row Scan Helpers (`gitmap/dbengine/scan_helpers.go`)
Implemented type-converting, null-safe row scanners:
- `ScanString(v any) string`: Handles `nil`, `string`, `[]byte`, `fmt.Sprint`.
- `ScanInt(v any) int`: Coerces `int64`, `int32`, `uint64` into `int`; returns `0` for `nil`.
- `ScanInt64(v any) int64`: Coerces numeric interfaces into `int64`; returns `0` for `nil`.
- `ScanUint64(v any) uint64`: Coerces `int64`, `int`, `uint` into `uint64`; returns `0` for `nil`.
- `ScanUint(v any) uint`: Converts uint64 into `uint`.
- `ScanBool(v any) bool`: Handles `bool`, integer `0`/`1` boolean representations; returns `false` for `nil`.
- `ScanFloat64(v any) float64`: Coerces float and integer values into `float64`; returns `0.0` for `nil`.

### 2.2 Automated Scanner Generation (`Scan{StructName}`)
The generator analyzes struct fields and generates:
```go
func ScanPipelineRunRecord(row dbengine.RowScanner) (*PipelineRunRecord, error)
```
- Declares temporary `any` variables for each field.
- Scans row safely.
- Maps raw values into struct fields via `dbengine.Scan*` helpers.

### 2.3 Dedicated Single Data Type QueryBuilders & Generic Repositories
The generator automatically outputs concrete type aliases for each struct:
- `type {s_name}QueryBuilder = dbengine.QueryBuilder[{s_name}, {enum_type}]`
- `type {short_name}QueryBuilder = {s_name}QueryBuilder`
- `type {s_name}Repository = dbengine.Repository[{s_name}, {enum_type}]`
- `type {short_name}Repository = {s_name}Repository`
This completely converts verbose generic signatures (`*dbengine.QueryBuilder[PipelineSplitDb, PipelineSplitDbFieldType]`) into concise, single named types (`*PipelineSplitDbQueryBuilder`, `*PipelineRunQueryBuilder`).

### 2.4 Value Semantics & Pointer Reduction
- Repository structs (`{StructName}DbRepo`) are lightweight wrappers around immutable handles (`*DbWrapper`, `*Repository`).
- Constructors return values: `NewPipelineRunDbRepo(db) PipelineRunDbRepo` (0 heap allocations).
- Methods use value receivers: `func (r PipelineRunDbRepo) Query() *PipelineRunQueryBuilder`.
- Standard query methods:
  - `FindAll(ctx context.Context) dbengine.ListResult[PipelineRunRecord]`
  - `First(ctx context.Context) dbengine.EntityResult[PipelineRunRecord]`
  - `Count(ctx context.Context) dbengine.Int64Result`
  - `Query() *PipelineRunRecordQueryBuilder`
  - `QueryBare() *PipelineRunRecordQueryBuilder`
  - `Db() *dbengine.DbWrapper`
  - `Repo() *PipelineRunRecordRepository`

### 2.5 Domain Repository Embedding Pattern
Domain repositories (such as `pipelinedb.PipelineRepository`) embed `PipelineRunRecordDbRepo` by value:
```go
type PipelineRepository struct {
    PipelineRunRecordDbRepo
}

func NewPipelineRepository(db *dbengine.DbWrapper) PipelineRepository {
    return PipelineRepository{
        PipelineRunRecordDbRepo: NewPipelineRunRecordDbRepo(db),
    }
}
```
Eliminates all boilerplate scanner and generic repository delegations while allowing domain repositories to cleanly define domain-specific queries (`GetRunById`, `GetRecentRuns`, `EnsureActiveErrorsView`) using value semantics.

---

## 3. Verification & Quality Gates

| Check / Suite | Status | Results |
| :--- | :--- | :--- |
| `TestScanString`, `TestScanInt`, `TestScanInt64`, `TestScanUint64`, `TestScanUint`, `TestScanBool`, `TestScanFloat64` | PASS | Full coverage of null safety and type coercions |
| `TestPipelineRunDbRepo_GeneratedRepo` | PASS | Verified `NewPipelineRunDbRepo`, `FindAll`, `First`, `Count`, and fluent queries |
| All PipelineDB Tests | PASS | 6 test functions green |
| All DBEngine Tests | PASS | 21 test functions green |
| Nested If Linter (`check-nested-ifs.py`) | PASS | 3,122 files scanned, 0 violations |
| Boolean & Enum Linter (`check-enum-and-boolean.py`) | PASS | 2,308 files scanned, 0 violations |
| CI/CD Local Runner (`06-cicd-local-runner.py`) | PASS | 16 of 16 quality gates green |
