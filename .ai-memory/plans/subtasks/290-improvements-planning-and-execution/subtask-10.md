# Subtask Plan 10: Add Strongly Typed Table Constants (`cli/constants/constants_store.go`)

- **Subtask ID:** `subtask-10`
- **Subtask Name:** `10-constants-store-table-names`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/constants/constants_store.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/constants/constants_store.go`, database table names follow PascalCase singular naming conventions, but are declared as untyped string constants:

```go
// Table names (v15: PascalCase + singular + {Table}Id PK).
const (
	TableRepo      = "Repo"
	TableGroup     = "Group"
	TableGroupRepo = "GroupRepo"
	TableRelease   = "Release"
	TableAlias     = "Alias"
)
```

1. Newer database tables introduced in recent Split-DB schema migrations (`RepoVersionHistory` and `DetectedProject`) lack centralized table constants, leading to raw string literal queries across store packages.
2. Under **Rule 5 (Strongly Typed Enums & Constants)**, database entities and table names should be governed by a typed enum ending with the `Type` suffix (`TableNameType`).
3. Adding `TableRepoVersionHistory` and `TableDetectedProject` alongside a strongly typed `TableNameType` enum eliminates raw string SQL errors and standardizes schema access across `cli/store/` and `cli/cmdscan/`.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Add `TableNameType` Enum & Constants
- **File:** `cli/constants/constants_store.go`
- **Definition:**
```go
// TableNameType represents strongly-typed database table names across the Split-DB tiers.
type TableNameType string

const (
	TableRepo               TableNameType = "Repo"
	TableGroup              TableNameType = "Group"
	TableGroupRepo          TableNameType = "GroupRepo"
	TableRelease            TableNameType = "Release"
	TableAlias              TableNameType = "Alias"
	TableRepoVersionHistory TableNameType = "RepoVersionHistory"
	TableDetectedProject    TableNameType = "DetectedProject"
)
```

### 2.2 Preserve Backwards-Compatible Untyped Constants
- **File:** `cli/constants/constants_store.go`
- **Requirements:**
  - String constants `TableRepoVersionHistory = "RepoVersionHistory"` and `TableDetectedProject = "DetectedProject"` can be directly used as `TableNameType` or string via untyped constant aliases, ensuring 100% compatibility with existing query builders.

---

## 3. Implementation Verification Checklist

- [x] `TableNameType` enum declared ending in `Type`.
- [x] `TableRepoVersionHistory = "RepoVersionHistory"` declared.
- [x] `TableDetectedProject = "DetectedProject"` declared.
- [x] Existing table constants (`TableRepo`, `TableGroup`, `TableRelease`, `TableAlias`) remain compatible.
- [x] Zero compiler builds or test runs during spec phase.
- [x] Formatting follows repo vertical whitespace standards.
