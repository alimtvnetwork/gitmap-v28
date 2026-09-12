# Plan 95: Nested If Elimination & Guard Clauses — Master Architectural Specification

Trigger Keywords & Aliases: `cg-nested-if`, `cg-execute nested-if`, `audit nested if`, `fix nested if`, `flatten conditionals`, `enforce guard clauses`

> **Prompt Version:** 2.1.0  
> **Synchronization:** Main Meta-Repo & Connected Workspaces  
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)

---

## 1. Executive Summary & Objective

Autonomously scan, plan, refactor, and fix all nested `if` statements across the codebase (nesting depth > 1) and single-line collapsed `if` statements using guard clauses, early returns, inverted conditions, and function decomposition (<= 8–15 lines) until 100% green without stopping.

Automated AST scan using `linter-scripts/check-nested-ifs.py` identified exactly 5 nested `if` violations across 5 files in `gitmap/cmd/` and `gitmap/store/`.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/02-canonical-size-tier.md`: Universal size limits and zero nested `if` mandate (functions <= 8 lines preferred, hard cap 15 lines; files <= 100 coding lines).
- `spec/02-coding-guidelines/01-cross-language/04-code-style/02-braces-and-nesting.md`: Elimination of nested conditional pyramids and mandatory braces.
- `spec/02-coding-guidelines/01-cross-language/04-code-style/04-blank-lines-and-spacing.md`: Blank lines before `return` and after closing `}` (R13-R16).
- `spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md`: Function size caps and extraction of single-purpose helpers.
- `spec/02-coding-guidelines/06-ai-optimization/01-index.md`: Zero truncation, zero placeholders, zero ghost diffs.

---

## 3. Domain-Specific Rules for Control Flow Flattening

1. **Rule 1 (Inverted Guard Returns):** When validating preconditions or error states, invert the conditional expression to return early (`if err != nil { return err }` or `if len(aliasName) == 0 { return args }`), keeping the main business logic un-indented at depth 0.
2. **Rule 2 (Helper Decomposition):** When handling sub-operations within loops or transactions, extract discrete <= 8-line helper functions (`saveBlobEntry`, `syncLatestReleaseFlag`) to isolate complexity and prevent indentation cascades.
3. **Rule 3 (Formatting Integrity):** Never collapse `if` blocks onto a single line to cheat line caps. Always include explicit newlines and braces. Ensure exactly one blank line before `return` and after `}` blocks.
4. **Rule 4 (Zero Swallow Policy):** Never discard errors during guard extraction; wrap and return all internal execution failures.

---

## 4. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Violation Description | Planned Fix | Status |
|:---|:---|:---:|:---|:---|:---|:---:|
| V-01 | `gitmap/cmd/chromeprofile_export_all.go` | 317 | `if _, err := db.Exec(..., name, blobFile, bytes); err != nil {` inside `if bytes, readErr := os.ReadFile(filePath); readErr == nil && len(bytes) > 0 {` | Nested `if` (depth 2) inside file read loop | Extract `saveBlobEntry` helper with inverted guard return `if err != nil || len(bytes) == 0 { return nil }` | Completed |
| V-02 | `gitmap/cmd/rm.go` | 230 | `if err := row.Scan(&id); err != nil {` inside `if appErr == nil {` in `scanFallbackRepoID` | Nested `if` (depth 2) inside nil error check | Invert condition to `if appErr != nil { return 0 }` and return early, flattening `row.Scan` check | Completed |
| V-03 | `gitmap/cmd/root.go` | 90 | `if err := resolveAliasContext(aliasName); err != nil {` inside `if len(aliasName) > 0 {` in `Execute()` | Nested `if` (depth 2) inside alias resolution | Extract `applyAliasContextIfPresent(command, args)` helper with early return if `len(aliasName) == 0` | Completed |
| V-04 | `gitmap/store/installer_delete.go` | 56 | `if appErr.Code != "E_INSTALLER_NOT_FOUND" {` inside `if appErr != nil {` in `runDeleteInstallerTx` | Nested `if` (depth 2) inside tx error handler | Invert guard `if appErr == nil { return nil }` and return early, flattening error code override | Completed |
| V-05 | `gitmap/store/release.go` | 51 | `if err := clearLatestRunner(runner, r.RepoID); err != nil {` inside `if r.IsLatest {` in `upsertReleaseTx` | Nested `if` (depth 2) inside `r.IsLatest` check | Extract `syncLatestReleaseFlag` helper with inverted guard `if !r.IsLatest { return nil }` | Completed |

---

## 5. Subtask Decomposition

- **Subtask 01:** `.lovable/plans/subtasks/95-nested-if-elimination-and-guard-clauses/01-task-flatten-cmd-nested-ifs.md`
  - Targets: `gitmap/cmd/chromeprofile_export_all.go`, `gitmap/cmd/rm.go`, `gitmap/cmd/root.go`
  - Scope: Flatten all nested `if` statements in `cmd/`, verify with `go test -v -short ./cmd/...`.
- **Subtask 02:** `.lovable/plans/subtasks/95-nested-if-elimination-and-guard-clauses/02-task-flatten-store-nested-ifs-and-ci-register.md`
  - Targets: `gitmap/store/installer_delete.go`, `gitmap/store/release.go`, `03-ai-scripts/01-index.md`
  - Scope: Flatten all nested `if` statements in `store/`, document `check-nested-ifs.py` in `03-ai-scripts/01-index.md`, verify with `go test -v -short ./store/...` and `python 03-ai-scripts/06-cicd-local-runner.py`.

---

## 6. Verification Plan

1. **Nested-If Linter:** `python linter-scripts/check-nested-ifs.py` must return 0 violations and exit 0.
2. **Go Test Suite:** `go test -v -short ./cmd/... ./store/...` must pass with 0 failures.
3. **CI/CD Quality Gate:** `python 03-ai-scripts/06-cicd-local-runner.py` Segment 1 (Linters & AST Checks) must pass 100% green.
