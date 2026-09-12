---
name: cg-result-wrapper
description: Autonomously scan, audit, refactor, and verify Go functions returning multi-value error tuples, replacing them with strongly-typed Result, ResultMap, and ResultSlice wrappers and structured AppError returns.
---

# Skill: Result Wrapper Types, Collections & AppError Returns (`cg-result-wrapper`)

This skill governs autonomous scanning, auditing, refactoring, and verification of Go functions returning multi-value tuples pairing return values with standard library errors (such as `(map[K]V, error)`, `([]T, error)`, or `(T, error)`).

## Core Architectural Directives

1. **Single Return Object Mandate:**
   - Multi-value error returns (e.g. `(map[string]int, error)`, `([]string, error)`, `(MyType, error)`) are prohibited on domain, store, and service functions.
   - All operations must return a single, strongly-typed envelope:
     - Key-Value Maps: `result.ResultMap[K, V]`
     - Slices / Lists: `result.ResultSlice[T]`
     - Scalar Values: `result.Result[T]`
     - Pure Side-Effects: `*apperror.AppError` (zero bare `error` returns)

2. **Structured Domain Errors:**
   - Raw standard library `error` returns are strictly banned in domain, store, and service layers.
   - All errors must use structured `*apperror.AppError` created with `apperror.New()` or `apperror.Wrap()`.

3. **Standardized Outer-Layer Inspection Predicates:**
   - `res.IsSuccess()`: `bool` - operation completed without error.
   - `res.IsFailure()` / `res.IsFailed()`: `bool` - operation encountered an error.
   - `res.HasError()`: `bool` - active error attached.
   - `res.IsEmptyError()` / `res.HasNoError()`: `bool` - no active error attached.
   - `res.IsEmpty()`: `bool` - collection has 0 items or data is empty.
   - `res.Data` / `res.Value`: Direct access to underlying payload.
   - `res.AppError()` / `res.Fault()`: Retrieves structured `*apperror.AppError`.
   - `res.Get(key)`: Safely retrieves map element `(V, bool)`.
   - `res.Has(key)`: Reports whether key exists in map.
   - `res.Count()`: Returns length of collection.
   - `res.Keys()`: Returns deterministically sorted slice of keys.
   - `res.Values()`: Returns slice of values ordered by keys.

4. **Zero Inversion & Affirmative Predicates:**
   - Never write `if !res.IsSuccess()`.
   - Always write `if res.IsFailure() { ... }` or `if res.IsFailed() { ... }`.

5. **Guard Clauses & Zero Nested Ifs:**
   - Maximum nesting depth is 1 (`if` within `if` is prohibited).
   - Use early returns and guard clauses.
   - Separate consecutive guard clauses with blank lines.
   - Always insert a blank line before `return`.

6. **Targeted Verification & CI Protection:**
   - DO NOT run `06-cicd-local-runner.py` during routine refactoring turns.
   - Run targeted linters (`golangci-lint run <pkg>`, `go vet <pkg>`, `python linter-scripts/...`).
   - Append modified files to `.lovable/temp/recent-file-changes.json` under lock.
