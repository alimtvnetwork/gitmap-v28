---
name: go-appwriter-and-result
description: Implement enterprise Go packages using structured AppError, monadic Result wrappers (result.Wrap, WrapSlice, WrapMap), deterministic enums, and pluggable StreamWriter pipelines.
---

# Go Architecture: AppWriter, Result Wrappers, & AppError

> **Invocation:** `@[go-appwriter-and-result](file;.cursor/skills/go-appwriter-and-result)`

This skill provides guidelines for writing and maintaining enterprise Go code within `04-code/golang/pkg/` and `cli/`.

## Core Pillars

1. **Structured Failure Standard (`*apperror.AppError` / `*appfault.AppError`):**
   - Standard failure return type across all packages: `*apperror.AppError` (or `*appfault.AppError`).
   - Classification via `errtype.Variation uint16` or `ErrorType` enum (`None`, `Validation`, `NotFound`, `Execution`, `Database`, `Network`, `Timeout`, `Internal`).
   - Structured caller tracking: `CallerInfo` (Function, File, Line), never a raw string.

2. **Monadic Result Containers with DRY `types.go` Reuse:**
   - Single entities: Functions returning a value along with possible failure must return `result.Wrap[T]`.
   - Collections / Slices: Functions returning slices must return `result.WrapSlice[T]`.
   - Keyed lookups: Functions returning maps must return `result.WrapMap[K, V]`.
   - Eliminates bare tuples `(T, error)`, `([]T, error)`, and `(map[K]V, error)` across all domain boundaries.
   - **DRY `types.go` Mandate:** Declare a single named type alias in the package's `types.go` first (e.g. `type ItemWrap = result.Wrap[Item]`, `type ItemSliceWrap = result.WrapSlice[Item]`, `type ItemMapWrap = result.WrapMap[string, Item]`). Implementation functions and callers MUST reuse this single type alias rather than writing raw generic brackets inline.
   - Provides fluent methods: `.IsSuccess()`, `.IsFailure()`, `.HasError()`, `.IsEmpty()`, `.HasRecord()`, `.IsCountOtherThan(N)`, `.Items()`, `.Value`, `.AppError()`.

3. **Pluggable Write Pipelines (`pkg/streamwriter`):**
   - Universal `Writer[T]` interface satisfying `sync.Locker`.
   - Runtime swappable formatters (`SetFormatMethod`), handlers (`SetWriteMethod`), and destinations (`SetDestination`).
   - Supported targets: Console, File (with POSIX permissions), JSON serializer, REST API client.

4. **No Bare Void in Go:**
   - Every function must return either a value container (`result.Wrap[T]`, `result.WrapSlice[T]`, `result.WrapMap[K, V]`) or an error (`*apperror.AppError`).
   - Parameter structs: Use `*Params` structs when arguments exceed 2-3 loose parameters.
