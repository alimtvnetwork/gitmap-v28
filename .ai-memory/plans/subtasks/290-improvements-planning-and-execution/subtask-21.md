# Subtask Plan 21: Transport-Based Repository Listing (`cli/store/repo.go`)

> **Subtask ID:** `subtask-21`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/repo.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `ListReposByTransport(transport string) ([]model.ScanRecord, error)` on `*DB` in `cli/store/repo.go`.
2. Allow consumers to directly query repositories matching a specified transport protocol (`SSH`, `HTTPS`) from SQLite.
3. Fall back gracefully to `db.ListRepos()` when `transport` is empty or whitespace.
4. Adhere strictly to Error Return Sovereignty (Rule 1), Function Length Limits (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
Currently, commands and subsystems that filter repositories by transport protocol (e.g. `cmdpull`, `cmdfixauth`, fleet synchronization) call `db.ListRepos()` to load all records into memory, then iterate through the entire slice checking `strings.EqualFold(r.IdentifiedTransport, transport)`.

This causes:
1. Redundant memory allocation of full repository lists when only a small subset is needed.
2. Repetitive in-memory filtering boilerplate across multiple CLI commands.
3. Lack of database index utilization for transport filtering.

---

## 3. Concrete Implementation Details

### 3.1 Method Addition in `cli/store/repo.go`
Add `ListReposByTransport` directly after `ListRepos`:

```go
// ListReposByTransport returns all repositories configured for the specified transport protocol.
func (db *DB) ListReposByTransport(transport string) ([]model.ScanRecord, error) {
	clean := strings.TrimSpace(transport)
	if clean == "" {
		return db.ListRepos()
	}

	rows, err := QueryWrapper(db.conn, constants.SQLSelectReposByTransport, clean).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}

	defer rows.Close()

	return scanRows(rows)
}
```

### 3.2 SQL Constant in `cli/constants/constants_store.go`
Ensure `constants.SQLSelectReposByTransport` is defined:
```go
const SQLSelectReposByTransport = `SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE IdentifiedTransport = ? COLLATE NOCASE ORDER BY Slug ASC`
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Database query errors bubble directly to the caller via `return nil, fmt.Errorf(...)`. No exit handlers (`cliexit.HandleError`) or `os.Exit` are invoked.
- **Rule 2 (Function Limit):** 14 lines of code, strictly complying with the <= 15 line rule.
- **Rule 3 (Zero Nested Ifs):** Uses early guard return for empty transport string. Zero nested conditionals.
- **Rule 4 (Positive Booleans):** No negative booleans used.
- **Rule 6 (Vertical Whitespace):** Exactly one blank line after `}` and before `return`.

---

## 5. Verification Criteria
- [x] Method `ListReposByTransport(transport string) ([]model.ScanRecord, error)` is declared on `*DB`.
- [x] Returns filtered repository records matching `IdentifiedTransport` case-insensitively.
- [x] Returns `db.ListRepos()` when transport is empty string.
- [x] Function is 14 lines or fewer with zero nested ifs.
- [x] All file paths are strictly relative lowercase Git paths.
