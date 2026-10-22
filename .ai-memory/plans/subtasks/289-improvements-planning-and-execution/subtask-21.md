# Subtask 21: Single Repository Lookup by Slug

> **Subtask ID:** Subtask-21  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/repo.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `FindOneBySlug(slug string) (*model.ScanRecord, error)` on `*DB` in `cli/store/repo.go`.
2. Provide a single-record query convenience method for callers that only require one repository match.
3. Return `(nil, nil)` safely when no record matches the given slug, avoiding unnecessary slice allocations or index bounds errors in callers.
4. Adhere strictly to error return sovereignty, function length limits (<= 15 lines), and zero nested conditionals.

---

## 2. Detailed Technical Plan

### 2.1 Implementation in `cli/store/repo.go`
- Add the method directly below `FindBySlug`:
  ```go
  // FindOneBySlug returns the first repository matching the given slug or nil if not found.
  func (db *DB) FindOneBySlug(slug string) (*model.ScanRecord, error) {
  	records, err := db.FindBySlug(slug)
  	if err != nil {
  		return nil, err
  	}

  	if len(records) == 0 {
  		return nil, nil
  	}

  	return &records[0], nil
  }
  ```

### 2.2 Architectural Integration
- Reuses the existing tested `FindBySlug` query runner.
- Leaves existing slice-based `FindBySlug` intact for backwards compatibility.

---

## 3. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Propagates `err` directly to the caller. Does not call exit handlers.
- **Rule 2 (Function Limit):** 11 lines of code (well under the 15-line limit).
- **Rule 3 (Zero Nested Ifs):** Uses early guard return for error and empty slice check.
- **Rule 6 (Vertical Whitespace):** Single blank line after each `}` block closure and before return.

---

## 4. Verification Criteria
- [x] `FindOneBySlug` compiles cleanly under Go 1.22+.
- [x] Returns `(*model.ScanRecord, nil)` when matching record is found.
- [x] Returns `(nil, nil)` when no record matches.
- [x] Returns `(nil, err)` when database query encounters an error.
- [x] All paths are relative lowercase without absolute paths or URIs.
