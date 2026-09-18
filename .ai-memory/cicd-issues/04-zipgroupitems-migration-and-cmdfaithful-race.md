# CI/CD Issue 04: Legacy SQLite Migration Failures & `cmdFaithful` Concurrency Race

- **Stage**: CI `go test -race` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/store`, `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

During CI test runs with `-race` (`go test -race -count=1 -timeout=15m ./cmd/... ./cloneconcurrency/... ./visibility/... ./store/... ./uipref/...`), failures occurred in the `cmd` package accompanied by spurious SQLite migration errors to `os.Stderr`:
```
[QueryWrapper Error]: exec failed: SQL logic error: no such table: ZipGroupItems (1)
query: ALTER TABLE ZipGroupItems ADD COLUMN RepoPath TEXT NOT NULL DEFAULT ''
```

## 2. Root Cause Analysis

Two independent issues interacted during concurrent test execution:

1. **Unchecked Legacy Table Migrations & Spurious `ExecWrapper` Logging**:
   - `migrateZipGroupItemPaths()` and `migratePendingTaskColumns()` in `gitmap/store/store.go` attempted `ALTER TABLE` on legacy tables (`ZipGroupItems`, `PendingTask`, `CompletedTask`) without guarding whether the tables existed in the database (tables were renamed to singular in v15).
   - In addition, `addColumnIfNotExists` called `ExecWrapper(db.conn, stmt)` which logged every execution failure to `os.Stderr` immediately, defeating `isBenignAlterError(err)` filtering for expected duplicate column / missing table errors.

2. **Unsynchronized Request-Scoped Globals in `cmdFaithful`**:
   - `cmdFaithfulVerify` and `cmdFaithfulExitOnMismatch` in `gitmap/cmd/clonetermverifystate.go` were declared as bare `bool` variables rather than `atomic.Bool`.
   - `cmdFaithfulExiter` in `gitmap/cmd/clonetermverifyexit.go` lacked mutex synchronization, allowing concurrent tests to race when reading or mocking the exiter function under `-race`.

## 3. Solution

1. **Store Migration Hardening**:
   - Guarded `migrateZipGroupItemPaths()` with `if !db.tableExists("ZipGroupItems") { return }`.
   - Guarded `migratePendingTaskColumns()` with `db.tableExists("PendingTask")` and `db.tableExists("CompletedTask")`.
   - Switched `addColumnIfNotExists` and benign data backfill statements to execute via `db.conn.Exec(stmt)` directly, preserving benign error suppression (`isBenignAlterError`).

2. **Concurrency Safety for `cmdFaithful` Verification State**:
   - Upgraded `cmdFaithfulVerify` and `cmdFaithfulExitOnMismatch` to `atomic.Bool`.
   - Wrapped `cmdFaithfulExiter` with `sync.RWMutex` via thread-safe getter `getCmdFaithfulExiter()` and setter `setCmdFaithfulExiter(fn)`.
   - Updated `withRecordingExiter` in `clonetermverifyexit_test.go` to use thread-safe accessors.

## 4. What NOT to Repeat

- **Never use `ExecWrapper` inside benign-tolerant migration helpers**: `ExecWrapper` unconditionally writes failures to `os.Stderr`. Use `db.conn.Exec` when errors like "duplicate column" or "no such table" are expected and filtered.
- **Never rely on bare `bool` for package-level state read across test goroutines**: Always use `atomic.Bool` or mutexes for variables toggled during CLI execution or test setup.
- **Always verify legacy table existence before ALTER/UPDATE**: In post-v15 schema designs, legacy plural tables (`ZipGroupItems`, `Releases`, `Repos`) may not exist in fresh test databases.
