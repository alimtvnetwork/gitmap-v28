---
name: go-appwriter-and-result
description: Implement enterprise Go packages using structured AppError, monadic Result wrappers, deterministic enums, and pluggable StreamWriter pipelines.
---

# Go Architecture: AppWriter, Result Wrappers, & AppError

This skill provides guidelines for writing and maintaining enterprise Go code within `04-code/golang/pkg/`.

## Core Pillars

1. **Structured Failure Standard (`*appfault.AppError`):**
   - Standard failure return type across all packages: `*appfault.AppError`.
   - Classification via `errtype.Variation uint16` enum (`None`, `Validation`, `NotFound`, `Execution`, `Database`, `Network`, `Timeout`, `Internal`).
   - Structured caller tracking: `CallerInfo` (Function, File, Line), never a raw string.

2. **Monadic Result Container (`result.Wrap[T]`):**
   - Functions returning a value along with possible failure must return `result.Wrap[T]`.
   - Eliminates bare tuples `(T, error)` in domain boundaries.
   - Provides methods: `.IsSuccess()`, `.IsFailure()`, `.Value()`, `.AppError()`, `.Match()`.

3. **Pluggable Write Pipelines (`pkg/streamwriter`):**
   - Universal `Writer[T]` interface satisfying `sync.Locker`.
   - Runtime swappable formatters (`SetFormatMethod`), handlers (`SetWriteMethod`), and destinations (`SetDestination`).
   - Supported targets: Console, File (with POSIX permissions), JSON serializer, REST API client.

4. **No Bare Void in Go:**
   - Every function must return either a value container (`result.Wrap[T]`) or an error (`*appfault.AppError`).
   - Parameter structs: Use `*Params` structs when arguments exceed 2-3 loose parameters.
