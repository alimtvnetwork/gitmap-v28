# Subtask 02: Database Mutations Cache Invalidation Hooks

## Status: READY
- **Subtask Slug**: `02-db-mutations-cache-invalidation`
- **Parent Plan**: `.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md`
- **Architecture Spec**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/01-architecture-spec.md`
- **Owner**: Worker 01 (Database Operations & Cache Engineer)
- **Target Subsystem**: `cli/cmddb`
- **Exclusive Owned Files**:
  - `cli/cmddb/cmd_add.go`
  - `cli/cmddb/cmd_edit.go`
  - `cli/cmddb/cmd_record_del.go`
  - `cli/cmddb/cmd_rm.go`
  - `cli/cmddb/cmddb_reset.go`
  - `cli/cmddb/rest_handlers_crud.go`
  - `cli/cmddb/rest_handlers_ddl.go`

---

## 1. Context & Objectives

To guarantee absolute data consistency across GitMap and prevent stale reads following write operations:
1. Every CLI command that modifies database records or schemas must proactively evict the database entry from the in-memory metadata cache.
2. Every REST API mutation endpoint (CRUD record operations and DDL table/view migrations) must trigger immediate cache eviction.
3. Database reset operations must clear all cached entries across all databases (`InvalidateAllDbTableCaches()`).
4. Invalidation functions:
   - `InvalidateTableMetadata(dbPath string)`: Evicts cached metadata for a specific database path.
   - `InvalidateAllDbTableCaches()`: Clears the entire in-memory metadata cache.
5. All modifications must adhere to:
   - Maximum 15 lines per function.
   - Zero nested `if` statements.
   - Positive boolean naming only.
   - Direct error returns (no `os.Exit` or exit handlers in leaf functions).

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Wire Invalidation in Record Insertion (`cli/cmddb/cmd_add.go`)

1. **Locate Target Function**:
   - `executeInsert` or `RunDbAdd` after the SQL `INSERT` statement executes successfully.
2. **Hook Invalidation**:
   - Resolve target database path: `dbPath, _ := ResolveDbPath(flags.targetDb)`.
   - Call `InvalidateTableMetadata(dbPath)`.
3. **Verify Error Isolation**:
   - Ensure cache invalidation is invoked only when the insert succeeds.
   - Cache invalidation itself is infallible and must not mask or alter transaction return errors.

---

### Step 2: Wire Invalidation in Record Update (`cli/cmddb/cmd_edit.go`)

1. **Locate Target Function**:
   - `executeUpdate` or `RunDbEdit` following successful transaction commit.
2. **Hook Invalidation**:
   - Resolve target database path: `dbPath, _ := ResolveDbPath(flags.targetDb)`.
   - Call `InvalidateTableMetadata(dbPath)`.
3. **Guard Logic**:
   - Ensure early error returns bypass cache invalidation.

---

### Step 3: Wire Invalidation in Record Deletion (`cli/cmddb/cmd_record_del.go`)

1. **Locate Target Function**:
   - `executeRecordDeletion` or `RunDbRecordDel` after the `DELETE FROM ...` statement completes.
2. **Hook Invalidation**:
   - Resolve target database path: `dbPath, _ := ResolveDbPath(flags.targetDb)`.
   - Call `InvalidateTableMetadata(dbPath)`.
3. **Unit Test Verification**:
   - Verify `cmd_record_del_test.go` still passes cleanly.

---

### Step 4: Wire Invalidation in Database Removal & Unregistration (`cli/cmddb/cmd_rm.go`)

1. **Locate Target Function**:
   - `executeAliasRemoval` or `RunDbRm` after deleting database files or unregistering aliases.
2. **Hook Invalidation**:
   - Retrieve resolved path for the target database alias.
   - Call `InvalidateTableMetadata(dbPath)`.
3. **Alias Cleanup**:
   - If an entire database file was deleted, also evict any alias pointing to it.

---

### Step 5: Wire Invalidation in Database Reset (`cli/cmddb/cmddb_reset.go`)

1. **Locate Target Function**:
   - `runDbResetAction` after database deletion/re-seeding completes.
2. **Hook Invalidation**:
   - Call `InvalidateAllDbTableCaches()`.
   - This ensures all cached schema versions, row counts, and table definitions for `gitmap.db`, `installation.db`, and all split databases are immediately purged.

---

### Step 6: Wire Invalidation in REST API Mutation Handlers

1. **In `cli/cmddb/rest_handlers_crud.go`**:
   - `handleInsertRecord`: After `res.LastInsertId()` succeeds, call `InvalidateTableMetadata(dbPath)`.
   - `handleUpdateRecord`: After `res.RowsAffected()` succeeds, call `InvalidateTableMetadata(dbPath)`.
   - `handleDeleteRecord`: After deletion succeeds, call `InvalidateTableMetadata(dbPath)`.
2. **In `cli/cmddb/rest_handlers_ddl.go`**:
   - `handleCreateTable`: Call `InvalidateTableMetadata(dbPath)` after table creation.
   - `handleAlter`: Call `InvalidateTableMetadata(dbPath)` after column addition.
   - `handleMigrate`: Call `InvalidateTableMetadata(dbPath)` after migration scripts execute.
   - `handleClone`: Call `InvalidateTableMetadata(destDbPath)` after database cloning.
   - `handleCreateView` & `handleDeleteView`: Call `InvalidateTableMetadata(dbPath)` after view alterations.

---

### Step 7: Regression and Consistency Testing

1. Create integration test verifying the complete read-write-read lifecycle:
   - Step A: `gitmap db read <testdb>` (populates cache with row count N).
   - Step B: `gitmap db add <testdb> <table> ...` (inserts record and evicts cache).
   - Step C: `gitmap db read <testdb>` (queries database, verifies row count N+1, repopulates cache).
2. Repeat for update, delete, reset, and REST API endpoints.

---

## 3. Verification & Quality Gates

- **Function Line Cap**: <= 15 lines per function.
- **Positive Booleans**: Use positive flags only (`hasSucceeded`, `isCommitted`, `isEvicted`).
- **Error Handling**: Preserve AppError return contracts; never suppress underlying SQLite errors.
- **Relative Path Standards**: All references use strictly relative paths.
