# Subtask 10: Repository Query Filter Modularization & Value Semantics

> **Subtask ID:** Subtask-10  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-10.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Architect and author `cli/store/repo_query_filter.go` to provide specialized, filtered repository query operations (including tag patterns, unpushed tracking queries, and transport/branch filters) separated from base CRUD persistence.
2. Enforce strict Golang value semantics across filter structures, query arguments, and slice returns, avoiding unnecessary pointer indirection (`*T`) for lightweight configuration structs.
3. Prevent `cli/store/repo.go` (currently ~383 lines) from breaching the 500-line ratchet limit by establishing a dedicated, single-responsibility query filtering module.
4. Enforce strict coding guidelines across `cli/store/repo_query_filter.go`:
   - Function limit: $\le 15$ lines per function ($\le 8$ lines preferred).
   - Conditionals: Zero nested `if` statements with guard clauses and early returns.
   - Positive booleans: Variables and fields prefixed strictly with `is` or `has` (`hasFilter`, `isUnpushedOnly`, `hasTag`).
   - Vertical spacing: Blank line after closing brace `}` and blank line before every `return`.
   - File length: The new module must remain compact and well under 250 lines.

---

## 2. Detailed Technical Plan

### 2.1 Define Value-Semantics Filter Structs in `cli/store/repo_query_filter.go`
- Define `RepoFilterOptions` struct using value fields and positive booleans:
  ```go
  package store

  import (
      "fmt"
      "strings"

      "github.com/alimtvnetwork/gitmap-v28/cli/constants"
      "github.com/alimtvnetwork/gitmap-v28/cli/model"
  )

  // RepoFilterOptions specifies criteria for filtering repository queries.
  type RepoFilterOptions struct {
      Branch          string
      Transport       string
      TagPattern      string
      HasUnpushedOnly bool
      IsExactBranch   bool
  }
  ```

### 2.2 Implement Specialized Filter Query Methods on `*DB`
- Implement `ListReposByFilter(opts RepoFilterOptions) ([]model.ScanRecord, error)`:
  - Decompose into query builder helper and execution helper to respect the $\le 15$-line ceiling.
  - Value semantics: `opts` passed by value.
  - Guard clauses for empty criteria:
    ```go
    func (db *DB) ListReposByFilter(opts RepoFilterOptions) ([]model.ScanRecord, error) {
        if !hasActiveFilters(opts) {
            return db.ListRepos()
        }

        query, args := buildFilterSQL(opts)

        return db.queryFilterRecords(query, args...)
    }
    ```
- Implement `hasActiveFilters(opts RepoFilterOptions) bool`:
  ```go
  func hasActiveFilters(opts RepoFilterOptions) bool {
      hasBranch := len(strings.TrimSpace(opts.Branch)) > 0
      hasTransport := len(strings.TrimSpace(opts.Transport)) > 0
      hasTag := len(strings.TrimSpace(opts.TagPattern)) > 0

      return hasBranch || hasTransport || hasTag || opts.HasUnpushedOnly
  }
  ```
- Implement `ListUnpushedRepos() ([]model.ScanRecord, error)`:
  - Query repositories with pending local commits tracked in SQLite metadata:
    ```go
    const sqlSelectUnpushedRepos = `SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE RepoName IN (SELECT RepoSlug FROM HeatmapLocalTracking WHERE AheadCount > 0) ORDER BY Slug ASC`

    func (db *DB) ListUnpushedRepos() ([]model.ScanRecord, error) {
        rows, err := QueryWrapper(db.conn, sqlSelectUnpushedRepos).Destruct()
        if err != nil {
            return nil, fmt.Errorf(constants.ErrDBQuery, err)
        }

        defer rows.Close()

        return scanRows(rows)
    }
    ```

### 2.3 Query Execution Helper
- Implement `queryFilterRecords(query string, args ...any) ([]model.ScanRecord, error)`:
  ```go
  func (db *DB) queryFilterRecords(query string, args ...any) ([]model.ScanRecord, error) {
      rows, err := QueryWrapper(db.conn, query, args...).Destruct()
      if err != nil {
          return nil, fmt.Errorf(constants.ErrDBQuery, err)
      }

      defer rows.Close()

      return scanRows(rows)
  }
  ```

---

## 3. Verification Criteria

- [ ] `cli/store/repo_query_filter.go` is designed with clean value semantics (passing `RepoFilterOptions` by value).
- [ ] `ListReposByFilter` and `ListUnpushedRepos` declare explicit `([]model.ScanRecord, error)` returns.
- [ ] No pointers to basic types or lightweight option structs are used.
- [ ] All functions in `cli/store/repo_query_filter.go` strictly comply with $\le 15$ lines limit.
- [ ] Zero nested `if` statements throughout; all conditionals flattened.
- [ ] All boolean fields and variables use strictly positive `is*` or `has*` prefixes.
- [ ] Both `cli/store/repo.go` and `cli/store/repo_query_filter.go` stay well under the 500-line ceiling.
- [ ] Strictly relative Git paths used throughout all references.
