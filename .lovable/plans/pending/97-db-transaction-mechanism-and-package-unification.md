# Plan 97: DB Transaction Mechanism & Package Repo-Wide Unification — Master Architectural Specification

Trigger Keywords & Aliases: `cg-db-tx`, `confirm db transactions`, `audit db transactions`, `fix db transactions`, `sqlite transaction unification`

> **Prompt Version:** 2.1.0  
> **Synchronization:** Main Meta-Repo & Connected Workspaces  
> **Budget:** N = 130 (PHASE_1_STEPS = 65, PHASE_2_STEPS = 65)

---

## 1. Executive Summary & Objective

Autonomously scan, plan, refactor, and verify all database transaction sites across the codebase, eliminating raw `db.Begin()`, `db.BeginTx()`, and unconfigured `sql.Open("sqlite", ...)` calls. Migrate all database mutations requiring transactional atomicity to the canonical `dbengine.WrapDb(...).WithTransaction(...)` / `WithImmediateTransaction(...)` API with automatic rollback on error or panic and typed `*apperror.AppError` returns until 100% green.

Deep repository scanning identified 7 raw transaction sites and 3 unconfigured SQLite open sites:
1. `gitmap/cmd/sequence_cmd.go:573`: Raw `repoDB.BeginTx(ctx, nil)`.
2. `gitmap/cmd/sshjoin_cmd.go:23`: Raw `db.BeginTx(ctx, nil)`.
3. `gitmap/cmd/ssh_alias_cmd.go:43`: Raw `db.Conn().BeginTx(ctx, nil)`.
4. `gitmap/cmd/commitin/runlog/inputs.go:35`: Raw `db.Begin()`.
5. `gitmap/store/makeallvisibility.go:44` & `94`: Raw `db.conn.Begin()` and legacy `commitOrWrap`.
6. `gitmap/store/owner_repo_name_index.go:26`: Raw `db.conn.Begin()`.
7. `gitmap/store/pendingtask.go:90`: Raw `db.conn.Begin()`.
8. `gitmap/cmd/chromeprofile_export_all.go:179`: Raw `sql.Open("sqlite", outPath)` and multi-insert loop without transaction.
9. `gitmap/cmd/agy_conv_scanner.go:63`: Raw `sql.Open("sqlite", dbPath)`.
10. `gitmap/store/chromeprofile_delete.go:29`: Multi-table DELETE operations executed outside transaction.

---

## 2. Authoritative Spec Citations

- `spec/04-database-conventions/09-universal-db-engine-and-code-generator.md`: All database transaction management must route through `dbengine.DbWrapper` (`WithTransaction`, `WithImmediateTransaction`, `WithExclusiveTransaction`, `WithReadOnlyTransaction`). Manual `Begin`/`Commit`/`Rollback` calls are strictly banned.
- `spec/04-database-conventions/01-sqlite-conventions.md`: All SQLite connections must be configured with `SetMaxOpenConns(1)`, WAL journal mode, `foreign_keys=ON`, `busy_timeout=5000`, and `synchronous=NORMAL` via `store.OpenSQLiteDB` or `store.ConfigureSQLiteConn`.
- `spec/03-error-manage/01-principles.md`: Zero-swallow policy; all database errors must be wrapped with `*apperror.AppError` providing operation context.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: Functions <= 8 lines preferred, <= 15 lines hard cap.
- `spec/02-coding-guidelines/01-cross-language/04-code-style/04-blank-lines-and-spacing.md`: Single blank line before `return` and after closing `}`.

---

## 3. Domain-Specific Rules for DB Transactions & SQLite Hardening

1. **Rule 1 (Mandatory `dbengine.WithTransaction`):** Direct invocations of `db.Begin()`, `db.BeginTx()`, or manual `tx.Commit()` / `defer tx.Rollback()` are strictly forbidden in business and store logic. All transaction blocks must be executed via `dbengine.WrapDb(conn, dbengine.DbSQLite).WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError { ... })`.
2. **Rule 2 (Centralized SQLite Factory `OpenSQLiteDB`):** Every SQLite connection opened in commands, split databases, or store operations must use `store.OpenSQLiteDB(path)` to enforce `SetMaxOpenConns(1)`, foreign keys, WAL mode, and busy timeouts.
3. **Rule 3 (Atomicity Over Multi-Statement Writes):** Any function performing multiple INSERT, UPDATE, or DELETE statements (e.g. `chromeprofile_export_all.go`, `chromeprofile_delete.go`, `transaction.go:deleteTransactionRows`) must enclose all statements in a transaction block.
4. **Rule 4 (Zero Swallowed Errors & Typed AppErrors):** Callback functions passed to `WithTransaction` must return `*apperror.AppError` directly when any operation fails, ensuring automatic rollback.
5. **Rule 5 (Function Size <= 15 Lines & Affirmative Naming):** Helper decomposition must maintain strict function limits (<= 15 lines) with affirmative booleans (`is*`, `has*`).

---

## 4. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
|:---|:---|:---:|:---|:---|:---:|
| V-TX01 | `gitmap/cmd/sequence_cmd.go` | 573 | `tx, err := repoDB.BeginTx(ctx, nil)` | Migrate to `dbengine.WrapDb(repoDB, dbengine.DbSQLite).WithTransaction(ctx, ...)` | Pending |
| V-TX02 | `gitmap/cmd/sshjoin_cmd.go` | 23 | `tx, err := db.BeginTx(ctx, nil)` | Migrate to `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(ctx, ...)` | Pending |
| V-TX03 | `gitmap/cmd/ssh_alias_cmd.go` | 43 | `tx, err := db.Conn().BeginTx(ctx, nil)` | Migrate to `dbengine.WrapDb(db.Conn(), dbengine.DbSQLite).WithTransaction(ctx, ...)` | Pending |
| V-TX04 | `gitmap/cmd/commitin/runlog/inputs.go` | 35 | `tx, err := db.Begin()` | Migrate to `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(ctx, ...)` | Pending |
| V-TX05 | `gitmap/store/makeallvisibility.go` | 44 | `tx, err := db.conn.Begin()` | Migrate `InsertMakeAllVisibilityPendingResults` to `dbengine.WrapDb(db.conn, ...).WithTransaction` | Pending |
| V-TX06 | `gitmap/store/makeallvisibility.go` | 94 | `tx, err := db.conn.Begin()` | Migrate `MarkMakeAllVisibilityResultsExcluded` to `dbengine.WrapDb(db.conn, ...).WithTransaction` | Pending |
| V-TX07 | `gitmap/store/owner_repo_name_index.go` | 26 | `tx, err := db.conn.Begin()` | Migrate `UpsertOwnerRepoNameIndex` to `dbengine.WrapDb(db.conn, ...).WithTransaction` | Pending |
| V-TX08 | `gitmap/store/pendingtask.go` | 90 | `tx, err := db.conn.Begin()` | Migrate `CompleteTask` to `dbengine.WrapDb(db.conn, ...).WithTransaction` | Pending |
| V-TX09 | `gitmap/cmd/chromeprofile_export_all.go` | 179 | `db, err := sql.Open("sqlite", outPath)` | Use `store.OpenSQLiteDB(outPath)` and wrap profile inserts in a transaction | Pending |
| V-TX10 | `gitmap/cmd/agy_conv_scanner.go` | 63 | `conn, err := sql.Open("sqlite", dbPath)` | Use `store.OpenSQLiteDB(dbPath)` | Pending |
| V-TX11 | `gitmap/store/chromeprofile_delete.go` | 29 | Two bare `ExecWrapper` delete calls | Enclose both deletions inside `dbengine.WrapDb(db.conn, ...).WithTransaction` | Pending |
| V-TX12 | `gitmap/store/transaction.go` | 131 | `deleteTransactionRows` loop bare Exec | Enclose pruning deletions in a transaction | Pending |
| V-TX13 | `gitmap/store/txhelpers.go` | 24 | `commitOrWrap` legacy helper | Clean up / deprecate legacy manual commit helper | Pending |

---

## 5. Subtask Decomposition

- [Subtask 01: CMD Package Transaction & SQLite Open Unification](plans/subtasks/97-db-transaction-mechanism-and-package-unification/01-task-cmd-transaction-migration-and-sqlite-open.md)
  - Target files: `gitmap/cmd/sequence_cmd.go`, `gitmap/cmd/sshjoin_cmd.go`, `gitmap/cmd/ssh_alias_cmd.go`, `gitmap/cmd/commitin/runlog/inputs.go`, `gitmap/cmd/chromeprofile_export_all.go`, `gitmap/cmd/agy_conv_scanner.go`.
- [Subtask 02: STORE Package Transaction & Atomicity Unification](plans/subtasks/97-db-transaction-mechanism-and-package-unification/02-task-store-transaction-migration-and-atomicity.md)
  - Target files: `gitmap/store/makeallvisibility.go`, `gitmap/store/owner_repo_name_index.go`, `gitmap/store/pendingtask.go`, `gitmap/store/chromeprofile_delete.go`, `gitmap/store/transaction.go`, `gitmap/store/txhelpers.go`, `gitmap/store/ssh_repo.go`.
- [Subtask 03: Unit Testing & CI Verification](plans/subtasks/97-db-transaction-mechanism-and-package-unification/03-task-transaction-tests-and-ci-verification.md)
  - Target files: `gitmap/cmd/sequence_cmd_test.go`, `gitmap/cmd/sshjoin_cmd_test.go`, `gitmap/store/ssh_repo_test.go`, `gitmap/store/transaction_atomicity_test.go`.
  - Quality Gate: `python 03-ai-scripts/06-cicd-local-runner.py --force --no-tests`.

---

## 6. Acceptance Criteria

```text
1. Zero occurrences of raw `Begin()` or `BeginTx()` in `gitmap/cmd/` and `gitmap/store/`.
2. All database mutations requiring atomicity route through `dbengine.WrapDb(...).WithTransaction(...)` or `WithImmediateTransaction(...)`.
3. All SQLite connections opened in utility commands use `store.OpenSQLiteDB(...)`.
4. All functions <= 15 lines with affirmative booleans and strict blank-line spacing.
5. All 29 gates in `06-cicd-local-runner.py` pass 100% green (`exit 0`).
```
