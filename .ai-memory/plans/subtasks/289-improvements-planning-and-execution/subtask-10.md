# Subtask Plan 10: Add Strongly Typed Query Enums for Repo Transport Lookups

- **Subtask ID:** `subtask-10`
- **Subtask Name:** `10-constants-store-repo-transport-enums`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/constants/constants_store.go`
  - `cli/store/repo_transport.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/constants/constants_store.go`, SQL statements for querying and updating repository transport (`SQLUpdateRepoTransport`, `SQLSelectRepoPrimaryURLByPath`) are defined as raw untyped string constants.

Callers in `cli/store/` and `cli/cmdpull/` pass arbitrary string arguments to identify transport query modes (e.g. `"path"`, `"slug"`, `"ssh"`, `"https"`), creating risks of typos, inconsistent case-folding, and untyped query dispatch:
1. Under **Rule 5 (Strongly Typed Enums)**, all enum-like identifiers MUST end with the `Type` suffix.
2. We must define strongly typed enum types: `RepoTransportQueryType` (for lookup strategy) and `RepoTransportActionType` (for database mutations).
3. These types prevent magic strings across database lookup routines and enforce compile-time type safety across the repository repository pattern.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Declare Typed Enums in `cli/constants/constants_store.go`
- **File:** `cli/constants/constants_store.go`
- **Additions:**
  ```go
  // RepoTransportQueryType represents the query target strategy for repository transport lookups.
  type RepoTransportQueryType string

  const (
  	RepoTransportQueryByPath RepoTransportQueryType = "by_path"
  	RepoTransportQueryBySlug RepoTransportQueryType = "by_slug"
  	RepoTransportQueryByID   RepoTransportQueryType = "by_id"
  )

  // RepoTransportActionType represents the database operation type for transport management.
  type RepoTransportActionType string

  const (
  	RepoTransportActionUpdate RepoTransportActionType = "update"
  	RepoTransportActionClear  RepoTransportActionType = "clear"
  	RepoTransportActionVerify RepoTransportActionType = "verify"
  )
  ```

### 2.2 Wire Typed Query Dispatch in `cli/store/`
- **File:** `cli/store/repo_transport.go`
- **Enhancement:**
  - Provide typed lookup helper consuming `RepoTransportQueryType`:
    ```go
    // GetRepoPrimaryURLByStrategy queries the primary URL and transport using a typed lookup strategy.
    func (db *DB) GetRepoPrimaryURLByStrategy(strategy constants.RepoTransportQueryType, target string) (int64, string, string, error) {
    	switch strategy {
    	case constants.RepoTransportQueryByPath:
    		cleanPath := NormalizeStoragePath(target)
    		return db.queryPrimaryURLByPath(cleanPath)
    	default:
    		return 0, "", "", appfault.NewValidationError("unsupported transport query strategy: " + string(strategy))
    	}
    }
    ```

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Strongly Typed Enums:** `RepoTransportQueryType` and `RepoTransportActionType` end with `Type` suffix.
- [ ] **No Magic Strings:** Eliminates untyped string tokens in transport queries.
- [ ] **15-Line Limit:** Function bodies remain strictly <= 15 lines.
- [ ] **Zero Nested Ifs:** Switch and guard clauses without nesting.
- [ ] **Positive Booleans:** Positive strategy validation checks.
- [ ] **Return Styling:** Standard vertical whitespace before returns.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
