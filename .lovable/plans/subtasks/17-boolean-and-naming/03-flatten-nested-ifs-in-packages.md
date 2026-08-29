# Subtask 03: Flatten Nested If Statements in `gitmap` Sub-Packages

> **Parent Plan:** `17-boolean-and-naming-audit.md`  
> **Scope:** `gitmap/archive`, `gitmap/clonenow`, `gitmap/cluster`, `gitmap/db`, etc.

## Objective

Eliminate all nested `if` statements in all Go sub-packages outside `cmd/`.

## Action Steps

1. Flatten nested `if` blocks in `gitmap/archive/`, `gitmap/clonefrom/`, `gitmap/clonenext/`, `gitmap/cluster/`, `gitmap/db/`.
2. Extract inner branches into small, dedicated private helpers.
3. Ensure no single-line `if` statements are introduced.
4. Verify package compilation across all packages with `go test ./...`.
