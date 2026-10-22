# Subtask 13: HelpSection Badge Data Model Enhancement

> **Subtask ID:** Subtask-13  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/termout/termhelp_types.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Extend `HelpSection` struct in `cli/termout/termhelp_types.go`.
2. Add `BadgeText` (string) for section categorization labels (e.g., `CORE`, `AI`, `SSOT`, `EXPERIMENTAL`).
3. Add `BadgeColor` (string) for custom styling of the badge label.
4. Maintain backwards compatibility with all existing callers constructing `HelpSection` literals.

---

## 2. Detailed Technical Plan

### 2.1 Struct Modification: `cli/termout/termhelp_types.go`
- Update `HelpSection` definition:
  ```go
  // HelpSection groups related command entries under a styled heading.
  type HelpSection struct {
  	Title      string
  	Color      string
  	BadgeText  string
  	BadgeColor string
  	Entries    []CommandEntry
  }
  ```

### 2.2 Compatibility Assessment
- Go struct field additions with zero-value defaults (`""`) are 100% backwards compatible with keyed struct literals across all CLI commands.
- Callers omitting `BadgeText` or `BadgeColor` remain unchanged and render standard unbadged headers.

---

## 3. Strict Coding Guidelines Compliance
- **File Sizing:** `termhelp_types.go` remains well under the 100-line limit (~30 lines).
- **Naming Conventions:** Fields use PascalCase (`BadgeText`, `BadgeColor`).
- **No Side Effects:** Pure type declaration change without logic branches.

---

## 4. Verification Criteria
- [x] `HelpSection` contains `BadgeText` and `BadgeColor` fields.
- [x] Existing callers constructing `HelpSection{Title: ..., Color: ..., Entries: ...}` continue to compile without syntax errors.
- [x] Zero linting or formatting regressions.
