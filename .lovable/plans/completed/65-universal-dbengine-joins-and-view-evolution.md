# 65 — Universal DBEngine Joins, Error-Guarded Query Builder & Automated View Evolution

**Status:** COMPLETED  
**Priority:** High  
**Category:** Database Architecture / Multi-Dialect Query Engine / Code Generation  
**Completed:** 2026-09-05  

---

## 1. Problem Statement & Root Cause

### Problem Statement
The generic repository and query builder abstractions previously had several architectural limitations:
1. Join construction lacked scoped method isolation, exposing global methods in an ambiguous chaining context.
2. Join fields, table targets, and match conditions relied on raw string literals (e.g. `"PipelineErrorRecord"`, `"ErrorText"`), bypassing compile-time safety and refactoring tooling.
3. Query compilation and terminal operations lacked internal error state tracking, returning raw tuples `(string, []any)` rather than monadic Result envelopes (`result.Result[T]`).
4. View lifecycle management (`CreateViewOrUseView`) required manual column name arguments from callers and had no database-backed mechanism to detect structural query modifications (such as changed `WHERE` filters or `JOIN` conditions).
5. There was no pre-execution SQL validation mechanism to detect syntax errors before executing DDL.
6. Domain-specific queries and business logic were mixed with generic database operations rather than residing in dedicated business logic repositories.

### Root Cause
Initial builder design focused on basic single-table CRUD operations and lacked sub-struct scoping, deterministic query hashing, and database-backed metadata tracking required for safe multi-table join and view workflows.

---

## 2. Architecture & Design Decisions

### 2.1 Scoped JoinBuilder Sub-Struct (`JoinBuilder[T, F]`)
Calling `.Join(table)`, `.InnerJoin(table)`, `.LeftJoin(table)`, `.RightJoin(table)`, or `.OuterJoin(table)` returns a dedicated `*JoinBuilder[T, F]`.
- Developer is strictly exposed to join-scoped methods:
  - `.Select(fields ...any)`: Projected fields from the joined table.
  - `.And(column any, op SqlOperator, val any)`: Extra filter conditions on the `ON` clause.
  - `.AndRaw(condition string)`: Raw SQL condition on the `ON` clause.
  - `.On(condition string)` / `.OnField(mainCol any, op SqlOperator, joinCol any)`: Completes the join and transitions back to `*QueryBuilder[T, F]`.

### 2.2 Result Envelopes & Comprehensive Error Propagation
- Both `QueryBuilder[T, F]` and `JoinBuilder[T, F]` carry an `err *apperror.AppError` field.
- If an error occurs, chained calls preserve the error state.
- Terminal methods (`First`, `FindAll`, `Count`, `Delete`, `Compile`, `CreateViewOrUseView`) guard against `b.err != nil` and return typed Result envelopes:
  - `First(ctx)` $\rightarrow$ `EntityResult[T]`
  - `FindAll(ctx)` $\rightarrow$ `ListResult[T]`
  - `Count(ctx)` $\rightarrow$ `Int64Result`
  - `Delete(ctx)` $\rightarrow$ `RowsAffectedResult`
  - `CreateViewOrUseView(...)` $\rightarrow$ `BoolResult`
  - `Compile()` $\rightarrow$ `CompiledQueryResult` (`result.Result[CompiledQuery]`)

### 2.3 Zero Magic Strings Across Joins
- All join methods accept `any` and convert through `toColumnName(col any) string`, which supports `string`, `fmt.Stringer`, and enum types.
- Field enums (e.g. `PipelineRunRecordDb.RunId`, `PipelineErrorRecordDb.ErrorText`) and table constants (`PipelineErrorRecordTable`) are used everywhere.

### 2.4 Automated View Evolution via Query Hash (`__gitmap_view_meta`)
- View metadata is stored in a dedicated database table:
  ```sql
  CREATE TABLE IF NOT EXISTS __gitmap_view_meta (
      ViewName TEXT PRIMARY KEY,
      QueryHash TEXT NOT NULL,
      ViewSql TEXT NOT NULL,
      UpdatedAt TEXT NOT NULL
  );
  ```
- If the deterministic SHA-256 hash of the view SQL matches $\rightarrow$ instant view reuse with **0 DDL**.
- If the hash differs or view does not exist $\rightarrow$ validates SQL via `EXPLAIN`, drops the stale view, creates the updated view, and records the new hash in metadata.
- Manual column parameter lists in `CreateViewOrUseView` are completely optional.

### 2.5 Pre-Execution SQL Syntax Validation (`ValidateSql`)
- `DbWrapper.ValidateSql(ctx, sqlStr)` executes `EXPLAIN <sql>` against the database engine.
- Bytecode compilation errors and missing table/column errors are caught immediately before any DDL or data changes occur.

### 2.6 Domain Business Logic Repository (`PipelineRepository`)
- Generic CRUD remains isolated in `dbengine.Repository`.
- Domain business logic lives in `pipelinedb.PipelineRepository`:
  - `GetRunById(ctx, runId uint64) EntityResult[PipelineRunRecord]`
  - `GetRecentRuns(ctx, repoSlug string, limit int) ListResult[PipelineRunRecord]`
  - `EnsureActiveErrorsView(ctx context.Context) BoolResult`
  - `Query()` and `QueryBare()` for pre-configured and custom projections.

---

## 3. Implementation Verification & Quality Gates

| Check / Suite | Status | Results |
| :--- | :--- | :--- |
| `TestDbWrapper_ValidateSql` | PASS | Verified `EXPLAIN` syntax validation |
| `TestDbWrapper_ViewHashMetaAndEvolution` | PASS | Verified metadata tracking, 0 DDL reuse, and evolution |
| `TestQueryBuilder_ErrorGuards` | PASS | Verified `*apperror.AppError` safety across all terminals |
| `TestQueryBuilder_TypedJoinsWithoutMagicStrings` | PASS | Verified join compilation with typed constants and enums |
| `TestPipelineRepository_CRUD` | PASS | Verified domain record insertions, ID retrieval, recent runs |
| `TestPipelineRepository_ActiveErrorsViewAndJoin` | PASS | Verified typed joins, view creation, hash reuse, direct querying |
| `TestPipelineRepository_FluentQueryWithEnums` | PASS | Verified bare query compilation with enum projections |
| All DBEngine Unit Tests | PASS | 14 test functions, 100% pass rate |
| All PipelineDB Unit Tests | PASS | 5 test functions, 100% pass rate |
| Nested If Linter (`check-nested-ifs.py`) | PASS | 3,120 files audited, max depth 1 enforced |
| Boolean & Enum Linter (`check-enum-and-boolean.py`) | PASS | 2,307 files audited, 0 violations |
| CI/CD Local Runner (`06-cicd-local-runner.py`) | PASS | 16 of 16 quality gates green |
| CLI Smoke Suite (`e2e-cli-smoke.py`) | PASS | 118 of 118 commands verified |
