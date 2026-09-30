# RCA-092: Cross-Platform Matrix Test Failures and Nested If Linter Violation

**Date:** 2026-09-30  
**Status:** Resolved  
**Severity:** High (CI/CD Pipeline Failure on Cross-Platform Build & Boolean/Enum Linter)  
**Affected Workflows:** `Cross-Platform Build` (macOS, Ubuntu), `CI` (`Boolean & Enum Linter`)

---

## 1. Root Cause Analysis
Two distinct issues caused failures across the matrix:
1. **Linter Gate Failure in CI (`Boolean & Enum Linter`):**
   `check-enum-and-boolean.py` detected a nested `if` at depth 2 in `cli/cmdnodes/nodes_clone_table.go:134`:
   ```go
   if !isSkipLocal {
       if isLocalSuccess { ... }
   ```
   Violated the repository's zero-nesting guideline.
2. **Cross-Platform Test Failures (`cmdmacro` and `cmdnodes`):**
   - In `cli/cmdnodes/nodes_clone_remote.go`, `buildWindowsWorkDirExecString` used `%q` to format Windows paths, causing `strconv.Quote` to double-escape backslashes (`D:\\\\work`), failing test assertions expecting `Set-Location "D:\work"`.
   - In `cli/cmdmacro/macro_add_interactive.go`, `isExplicitHelperCmd` was mistakenly wrapping `processInBuilderCommand` in `processInteractiveStepLine`, which prevented standard interactive helpers like `ls` from executing as builder helpers, causing `TestProcessInteractiveStepLine` to fail.

## 2. Why Not Caught Earlier?
The prior release run used `--skip-tests` to accelerate release execution, deferring matrix testing to GitHub Actions. `python linter-scripts/check-enum-and-boolean.py` and package unit tests were not executed prior to the tag.

## 3. Remediation
1. **Nested If Refactoring (`cli/cmdnodes/nodes_clone_table.go`):**
   Decomposed `renderFleetSummaryFooter` into `countFleetResults`, `adjustForLocalResult`, and `renderSkippedLocalSummary` with early returns and guard clauses. Zero nested ifs remaining.
2. **Windows Path Quoting Fix (`cli/cmdnodes/nodes_clone_remote.go`):**
   Replaced `%q` formatting with literal quotes `\"%s\"` in `buildWindowsWorkDirExecString` to preserve clean single backslashes for PowerShell remote execution.
3. **Interactive Step Helper Handling (`cli/cmdmacro/macro_add_interactive.go`):**
   Restored direct `processInBuilderCommand` delegation in `processInteractiveStepLine`.
4. **Validation:**
   - `python linter-scripts/check-enum-and-boolean.py`: PASS (3072 files).
   - `go test -short ./cmdmacro/... ./cmdnodes/...`: PASS.
   - `go test -run=^$ ./...`: PASS across all packages.

## 4. Prevention
- Always run `python linter-scripts/check-enum-and-boolean.py` whenever modifying control flow or conditionals in Go or TypeScript files.
- Run targeted package unit tests locally before initiating automated release branching.
