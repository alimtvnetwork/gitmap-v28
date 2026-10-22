# Subtask Plan 07: Enforce Positive Boolean `hasUpdated` in `repo_transport.go`

- **Subtask ID:** `subtask-07`
- **Subtask Name:** `07-store-repo-transport-positive-boolean`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/store/repo_transport.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/store/repo_transport.go`, `UpdateRepoTransport` executes:
`UPDATE Repo SET SshUrl = ?, IdentifiedTransport = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE AbsolutePath = ?`

1. The method currently discards the `sql.Result` returned by `ExecWrapper(...).Destruct()`, ignoring whether any rows were actually affected.
2. If `repoPath` does not exist in the database, SQLite executes the statement successfully with 0 rows affected, returning `nil` error. Callers assume the record was updated, while the repository remained absent from `gitmap.db`.
3. In addition, existing codebase patterns frequently used negative boolean checks (`if rows == 0` or `isMissing`).
4. Under **Rule 4 (Positive Boolean Naming)** and **Rule 1 (Error Return Sovereignty)**:
   - Introduce `UpdateRepoTransportWithStatus(repoPath, sshURL, transport string) (bool, error)` returning the positive boolean `hasUpdated`.
   - Maintain `UpdateRepoTransport` as a backwards-compatible delegate returning `error`.
   - Ensure all conditional branches inspect `hasUpdated` directly.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Implement `UpdateRepoTransportWithStatus`
- **File:** `cli/store/repo_transport.go`
- **Signature:** `func (db *DB) UpdateRepoTransportWithStatus(repoPath, sshURL, transport string) (bool, error)`
- **Requirements:**
  - Canonicalize path, URL, and transport.
  - Call `ExecWrapper(db.conn, constants.SQLUpdateRepoTransport, cleanURL, cleanTransport, cleanPath).Destruct()`.
  - Check error; return `(false, err)` immediately if non-nil.
  - Query `res.RowsAffected()`.
  - Evaluate positive boolean: `hasUpdated := rows > 0`.
  - Return `(hasUpdated, nil)`.
  - Function length <= 15 lines.

```go
// UpdateRepoTransportWithStatus updates remote transport and returns whether a row was modified.
func (db *DB) UpdateRepoTransportWithStatus(repoPath, sshURL, transport string) (bool, error) {
	cleanPath := NormalizeStoragePath(repoPath)
	cleanURL := strings.TrimSpace(sshURL)
	cleanTransport := cleanTransportOrDefault(transport)

	res, err := ExecWrapper(db.conn, constants.SQLUpdateRepoTransport, cleanURL, cleanTransport, cleanPath).Destruct()
	if err != nil {
		return false, err
	}

	rows, errRows := res.RowsAffected()
	if errRows != nil {
		return false, errRows
	}

	hasUpdated := rows > 0

	return hasUpdated, nil
}
```

### 2.2 Preserve Backwards-Compatible Delegate `UpdateRepoTransport`
- **File:** `cli/store/repo_transport.go`
- **Signature:** `func (db *DB) UpdateRepoTransport(repoPath, sshURL, transport string) error`
- **Implementation:**
  ```go
  // UpdateRepoTransport updates the remote SSH URL and identified transport for a repository.
  func (db *DB) UpdateRepoTransport(repoPath, sshURL, transport string) error {
  	_, err := db.UpdateRepoTransportWithStatus(repoPath, sshURL, transport)

  	return err
  }
  ```

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** Handles both `ExecWrapper` error and `RowsAffected` error.
- [ ] **15-Line Limit:** Both functions strictly <= 14 lines.
- [ ] **Zero Nested Ifs:** Flattened with guard clauses.
- [ ] **Positive Booleans:** Uses `hasUpdated` (no `isUnchanged` or `hasFailed`).
- [ ] **Return Styling:** Standard vertical whitespace before returns.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
