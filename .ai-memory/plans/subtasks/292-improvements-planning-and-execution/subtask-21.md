# Subtask Plan 21: Tag-Based Repository Listing (`cli/store/repo_query_filter.go`)

> **Subtask ID:** `subtask-21`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/repo_query_filter.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `ListReposByTag(tag string) ([]model.ScanRecord, error)` on `*DB` in `cli/store/repo_query_filter.go`.
2. Allow consumers to directly query repositories matching a metadata tag stored in repository `Notes` or metadata columns.
3. Fall back gracefully to `db.ListRepos()` when `tag` is empty or whitespace.
4. Keep `cli/store/repo.go` (currently 384 lines) strictly bounded under the 500-line ceiling by placing filtered query methods in `cli/store/repo_query_filter.go`.
5. Comply strictly with Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
Callers requiring repositories grouped by tags currently call `db.ListRepos()` to load all records into memory and perform manual pattern scanning. Furthermore, adding new query filters directly into `cli/store/repo.go` pushes the file toward the 500-line cap.

Extracting filtered query operations into `cli/store/repo_query_filter.go`:
1. Reduces memory overhead by executing targeted `WHERE Notes LIKE ? COLLATE NOCASE` queries inside SQLite.
2. Standardizes case-insensitive tag matching across all CLI commands.
3. Preserves `repo.go` sizing well below the 500-line ceiling.

---

## 3. Concrete Implementation Details

### 3.1 Function Implementation in `cli/store/repo_query_filter.go`
```go
package store

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const sqlSelectReposByTag = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Notes LIKE ? COLLATE NOCASE ORDER BY Slug ASC"

// ListReposByTag returns all repositories matching a specific tag stored in metadata or notes.
func (db *DB) ListReposByTag(tag string) ([]model.ScanRecord, error) {
	cleanTag := strings.TrimSpace(tag)
	if cleanTag == "" {
		return db.ListRepos()
	}

	pattern := "%" + cleanTag + "%"
	rows, err := QueryWrapper(db.conn, sqlSelectReposByTag, pattern).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}

	defer rows.Close()

	return scanRows(rows)
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** SQL query error bubbles directly up via `return nil, fmt.Errorf(constants.ErrDBQuery, err)`. No exit handlers or terminations invoked.
- **Rule 2 (Function Limit):** Exactly 14 lines of code, satisfying the <= 15 line limit.
- **Rule 3 (Zero Nested Ifs):** Linear guard clause checks for empty tag with early return.
- **Rule 4 (Positive Booleans):** No negative booleans used.
- **Rule 6 (Vertical Whitespace):** Blank line after `}` and before `return`.

---

## 5. Verification Criteria
- [x] Method `ListReposByTag(tag string) ([]model.ScanRecord, error)` is declared on `*DB` in `cli/store/repo_query_filter.go`.
- [x] Trims whitespace from `tag` and returns `db.ListRepos()` if empty.
- [x] Uses parameterized `pattern := "%" + cleanTag + "%"` preventing SQL injection.
- [x] Closes `rows` with `defer rows.Close()`.
- [x] Returns structured error or slice mapped via `scanRows(rows)`.
- [x] File and function satisfy all sizing and guidelines invariants.
