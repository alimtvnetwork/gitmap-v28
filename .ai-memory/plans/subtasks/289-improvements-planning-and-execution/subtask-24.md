# Subtask 24: File Touch Type Enumeration

> **Subtask ID:** Subtask-24  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/heatmap_schema.go`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Define a strongly typed enum `FileTouchType` in `cli/store/heatmap_schema.go`.
2. Define exported constants for supported touch categories: `FileTouchRead`, `FileTouchEdit`, and `FileTouchCreated`.
3. Eliminate magic string comparisons for file operations across heatmap logging, agent servers, and reporting routines.
4. Adhere strictly to the project convention requiring all enum types to end with the `Type` suffix.

---

## 2. Detailed Technical Plan

### 2.1 Enum & Constant Declarations in `cli/store/heatmap_schema.go`
- Add the following type and constant block to `cli/store/heatmap_schema.go`:
  ```go
  // FileTouchType represents the categorical operation applied to a touched file.
  type FileTouchType string

  const (
  	FileTouchRead    FileTouchType = "READ"
  	FileTouchEdit    FileTouchType = "EDIT"
  	FileTouchCreated FileTouchType = "CREATED"
  )
  ```

### 2.2 Validation Helper Method
- Provide a clean helper to validate whether a string corresponds to a known `FileTouchType`:
  ```go
  // IsValid reports whether the touch type is a recognized enumeration value.
  func (t FileTouchType) IsValid() bool {
  	switch t {
  	case FileTouchRead, FileTouchEdit, FileTouchCreated:
  		return true
  	default:
  		return false
  	}
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Rule 4 (Positive Boolean):** Helper method `IsValid()` uses positive polarity.
- **Rule 5 (Enum Suffix):** Type name `FileTouchType` strictly ends with `Type`.
- **Rule 2 (Function Limit):** `IsValid` is 8 lines (<= 15 line limit).
- **Whitespace Rules:** Blank line before `return` and after `}`.

---

## 4. Verification Criteria
- [ ] `FileTouchType` compiles cleanly under Go 1.22+.
- [ ] Constants `FileTouchRead`, `FileTouchEdit`, `FileTouchCreated` exported and accessible across packages.
- [ ] `IsValid()` returns true for all three defined enum values and false for unrecognized strings.
- [ ] Zero compiler builds and zero tests executed during planning.
