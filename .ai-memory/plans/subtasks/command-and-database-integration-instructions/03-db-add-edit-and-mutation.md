# Subtask 03: `db-add` & `db-edit` Mutation Engine with Single-Writer Locks

- **Subtask ID**: `subtask-03-db-add-edit-and-mutation`
- **Parent Task**: `command-and-database-integration-instructions`
- **Specification Reference**:
  - `02-spec/21-app/command-and-database-integration-instructions/02-cli-command-contracts.md`
- **Owned Files**:
  - `cli/cmddb/cmddb_add.go`
  - `cli/cmddb/cmddb_edit.go`
  - `cli/cmddb/cmddb_mutation_helpers.go`
  - `cli/cmddb/cmddb_mutation_test.go`
- **Assigned Worker**: Worker 03
- **Status**: PENDING

---

## 1. Objectives

1. Implement `gitmap db-add` (alias `gitmap dba`) in `cli/cmddb/cmddb_add.go`:
   - **Interactive Prompt Mode:** When invoked as `gitmap db-add <alias> <table>`, introspect columns via PRAGMA and prompt the user in terminal for each non-autoincrement field.
   - **CSV Values Ingestion Mode:** When invoked with `csv "col1", "col2", ...`, parse comma-separated literals and insert directly into the target table.
2. Implement `gitmap db-edit` (alias `gitmap dbe`) in `cli/cmddb/cmddb_edit.go`:
   - Update existing record by ID: `gitmap db-edit <alias> -table "table" -id "<id>" csv "col1", "col2", ...`.
   - Introspect target table primary key dynamically to construct deterministic `UPDATE ... WHERE <pk> = ?` statements.
3. Concurrency & Transactional Safety:
   - Use `BEGIN IMMEDIATE TRANSACTION` ensuring exclusive write lock acquisition.
   - Strictly prohibit string concatenation in queries; use parameterized placeholders (`?`).
   - Validate column counts against writable schema columns prior to executing mutations.
4. Provide unit tests in `cli/cmddb/cmddb_mutation_test.go` verifying valid insertions, updates, foreign key checks, and rollbacks on constraint violations.

---

## 2. CLI Command Syntax

```bash
# Interactive insertion prompting for each column
gitmap db-add master Repositories

# Direct CSV insertion into table
gitmap db-add master Repositories csv "my-new-repo", "d:/work/my-new-repo", "ACTIVE", 0

# Modify existing record identified by primary key
gitmap db-edit master -table "Repositories" -id "105" csv "updated-name", "d:/work/my-new-repo", "ACTIVE", 1
```

---

## 3. Verification Criteria

1. `go test ./cli/cmddb -run TestDbMutation -v` passes cleanly.
2. Inserting records via `db-add` increases table row count by 1 and persists data reliably.
3. Updating records via `db-edit` modifies only the targeted row identified by primary key.
4. Calling `db-add` or `db-edit` with mismatched column counts returns an explicit error without executing partial writes.
