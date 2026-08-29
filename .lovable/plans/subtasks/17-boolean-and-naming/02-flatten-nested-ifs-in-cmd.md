# Subtask 02: Flatten Nested If Statements in `gitmap/cmd`

> **Parent Plan:** `17-boolean-and-naming-audit.md`  
> **Scope:** `gitmap/cmd/*.go`

## Objective

Eliminate all nested `if` statements (nesting depth > 1) in `gitmap/cmd` by applying guard clauses, early returns, and extracting helper functions.

## Action Steps

1. Flatten nested `if` conditions in `cmd/cg*.go`, `cmd/chromeprofile*.go`, `cmd/cluster*.go`, etc.
2. Invert conditions into guard clauses with early returns.
3. Extract validation logic into helper functions (< 8 lines).
4. Verify package compilation with `go test ./cmd`.
