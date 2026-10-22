# Subtask 10: Repository Query Extraction, Value Semantics & Function Sizing

> **Subtask ID:** Subtask-10  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-10.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope

1. Decompose the oversized `cli/store/repo.go` (currently 350 lines) by extracting all repository read/query methods into `cli/store/repo_query.go`.
2. Enforce Golang value semantics and pointer reduction across query parameters and record returns.
3. Enforce the hard limit of $\le 15$ lines per function ($\le 8$ lines preferred) across all extracted query methods.
4. Eliminate nested conditional statements with guard clauses and early returns.
5. Standardize boolean parameters and local checks with positive `is*` or `has*` prefixes (e.g., `hasFilter`, `hasExactMatch`).
6. Apply strict vertical blank line rules (blank line before `return` and after `}`).

---

## 2. Detailed Technical Plan

### 2.1 File Decomposition Strategy
- In `cli/store/repo.go`, retain only mutations, lifecycle hooks, and DDL execution:
  - `NormalizeStoragePath`
  - `OnRepoAddedHook`
  - `UpsertRepos`, `upsertReposTx`, `upsertOneRepo`
  - `DeleteByPath`, `DeleteBySlug`
- Extract all repository querying operations into `cli/store/repo_query.go`:
  - `ListRepos() ([]model.ScanRecord, error)`
  - `ListReposByTransport(transport string) ([]model.ScanRecord, error)`
  - `CountRepos() (int, error)`
  - `FindByID(id int64) ([]model.ScanRecord, error)`
  - `FindBySlug(slug string) ([]model.ScanRecord, error)`
  - `FindByPath(absPath string) (*model.ScanRecord, error)`
  - `SearchBySlug(term string) ([]model.ScanRecord, error)`
  - `SearchRepos(query string) ([]model.ScanRecord, error)`

### 2.2 Function Sizing & Guard Clause Inversion ($\le 15$ Lines Limit)
- Ensure every extracted function in `repo_query.go` adheres to the 15-line limit:
  ```go
  // Compliant ListReposByTransport implementation:
  func (db *DB) ListReposByTransport(transport string) ([]model.ScanRecord, error) {
      clean := strings.TrimSpace(transport)
      if clean == "" {
          return db.ListRepos()
      }

      rows, err := QueryWrapper(db.conn, sqlSelectReposByTransport, clean).Destruct()
      if err != nil {
          return nil, fmt.Errorf(constants.ErrDBQuery, err)
      }

      defer rows.Close()

      return scanRows(rows)
  }
  ```
- Guarantee every closing brace `}` is followed by a blank line, and every `return` is preceded by a blank line.

### 2.3 Value Semantics & Positive Boolean Naming
- Return `model.ScanRecord` or `[]model.ScanRecord` slices by value rather than pointer slices `[]*model.ScanRecord` where mutation is not required, avoiding heap allocation pressure.
- Replace any negative boolean flags with positive naming (`hasFilter`, `isExactMatch`).

---

## 3. Verification Criteria

- [ ] All query routines extracted from `cli/store/repo.go` to `cli/store/repo_query.go`.
- [ ] `cli/store/repo.go` reduced from 350 lines to $\le 120$ lines.
- [ ] Every function in `cli/store/repo_query.go` is $\le 15$ lines.
- [ ] Zero nested `if` statements in `cli/store/repo_query.go`.
- [ ] Value semantics applied across query returns.
- [ ] Positive boolean naming (`is*`, `has*`) applied consistently.
- [ ] Blank line formatting rules strictly respected.
