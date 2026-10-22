# Subtask Plan 01: Error Return Sovereignty in Repository Matching (`cli/cmdcd/cdops.go`)

- **Subtask ID:** `subtask-01`
- **Subtask Name:** `01-cdops-repo-matching-error-sovereignty`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdcd/cdops.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdcd/cdops.go`, repository lookup and directory resolution functions must strictly adhere to **Rule 1 (Error Return Sovereignty)**.

Currently, `findCDRecords(db *store.DB, name string) []model.ScanRecord` ignores database retrieval errors:
1. When `db.FindBySlug(...)` fails or returns no records, it falls back to `all, listErr := db.ListRepos()`.
2. If `listErr != nil`, it merely prints a warning to `os.Stderr` (`fmt.Fprintf(os.Stderr, "  ⚠ Could not list repos: %v\n", listErr)`) and continues executing with an empty or corrupted `all` slice, returning an empty list to `lookupCDRecords`.
3. In `lookupCDRecords(name string) ([]model.ScanRecord, error)`, the database call is wrapped as `return findCDRecords(db, name), nil`, unconditionally discarding any database listing failure.
4. This violates Error Return Sovereignty: leaf functions that encounter I/O or DB failures MUST propagate them rather than swallowing them and returning `nil`.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `findCDRecords`
- **File:** `cli/cmdcd/cdops.go`
- **Signature:** `func findCDRecords(db *store.DB, name string) ([]model.ScanRecord, error)`
- **Requirements:**
  - Attempt lookup via `db.FindBySlug(strings.ToLower(cleanName))`.
  - If a fatal database query error occurs (other than no rows found), return `nil, err`.
  - If no records match slug, execute `all, listErr := db.ListRepos()`.
  - If `listErr != nil`, return `nil, listErr`.
  - Return `deduplicateCDRecords(cmdpull.FindBySlug(all, cleanName)), nil`.
  - Keep function length <= 15 lines.

```go
func findCDRecords(db *store.DB, name string) ([]model.ScanRecord, error) {
	cleanName := strings.TrimRight(name, "/\\")
	repos, err := db.FindBySlug(strings.ToLower(cleanName))
	if err != nil {
		return nil, err
	}

	hasValidRepos := len(repos) > 0
	if hasValidRepos {
		return deduplicateCDRecords(repos), nil
	}

	all, listErr := db.ListRepos()
	if listErr != nil {
		return nil, listErr
	}

	return deduplicateCDRecords(cmdpull.FindBySlug(all, cleanName)), nil
}
```

### 2.2 Refactor `lookupCDRecords`
- **File:** `cli/cmdcd/cdops.go`
- **Signature:** `func lookupCDRecords(name string) ([]model.ScanRecord, error)`
- **Requirements:**
  - Check `openDB()` error and return immediately.
  - Defer `db.Close()`.
  - Directly propagate the result of `findCDRecords(db, name)`.
  - Function length <= 12 lines.

```go
func lookupCDRecords(name string) ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)

		return nil, fmt.Errorf(constants.ErrListDBFailed, err)
	}

	defer db.Close()

	return findCDRecords(db, name)
}
```

### 2.3 Verify Upstream Callers in `cli/cmdcd/cdops.go`
- In `runCDLookup(name string, args []string) error`:
  - `records, err := lookupCDRecords(name)` is already checking `if err != nil { return err }`.
  - With `findCDRecords` propagating errors, database corruption or connection failures will now properly bubble up through `runCDLookup` to the root dispatcher and `finishCommandAudit`.

---

## 3. Implementation Verification Checklist

- [x] `findCDRecords` updated to return `([]model.ScanRecord, error)`.
- [x] Database errors in `db.FindBySlug` and `db.ListRepos` are returned directly instead of swallowed.
- [x] `lookupCDRecords` forwards error directly from `findCDRecords(db, name)`.
- [x] Function lengths <= 15 lines each.
- [x] Zero nested `if` statements.
- [x] Blank lines preserved before `return` and after `}`.
