# AppError Package Reference — Overview, invariants, StackTrace

> **Parent:** [AppError Package Reference](./01-index.md)
> **Version:** 1.3.0
> **Updated:** 2026-03-31

---

## Overview

The `apperror` package provides **structured application errors with mandatory stack traces** and **generic result wrappers** for all service-level returns. Every error created through this package automatically captures a stack trace at the point of creation. No error is ever swallowed or lost.

### Core Invariants

| # | Invariant | Description |
|---|-----------|-------------|
| I-1 | **Serialization Round-Trip** | Every `*AppError` is fully serializable to JSON/YAML and deserializable back. An error serialized with `json.Marshal` and deserialized with `json.Unmarshal` must preserve: code, message, details, values, diagnostics, stack trace, and cause message. This enables error transport across process boundaries (subprocess JSON protocol, HTTP responses, log storage, error history DB). |
| I-2 | **No Raw `error` in Structs** | Struct fields that hold errors must use `*AppError`, never Go's `error` interface. The `error` interface is not serializable — it has no guaranteed JSON representation. Only `*AppError` carries the structured data (code, stack, diagnostics) needed for serialization, logging, and debugging. |
| I-3 | **Stack Trace Always Present** | Every `*AppError` captures a stack trace at creation. No error exists without diagnostic context. |
| I-4 | **Zero Raw Error Returns** | Service methods return `*AppError` or `Result[T]`, never bare `error`. See §10.5. |

### Why Serialization Matters

`*AppError` travels across boundaries that raw `error` cannot:

| Boundary | Mechanism | Raw `error` | `*AppError` |
|----------|-----------|-------------|-------------|
| HTTP API response | JSON envelope | ❌ Loses context | ✅ Full diagnostic payload |
| Subprocess protocol | JSON stdout | ❌ String only | ✅ Code + stack + values |
| Error history DB | GORM/SQLite | ❌ Unstructured | ✅ Queryable fields |
| Log aggregation | Structured logging | ❌ Flat string | ✅ Indexed by code/stack |
| AI diagnostic | Clipboard format | ❌ Useless | ✅ Full reproduction context |
| Cross-service call | Direct method | ❌ No stack | ✅ Chained PreviousTrace |

### Package Files

| File | Purpose | Target Lines |
|------|---------|--------------|
| `stack_trace.go` | StackFrame, StackTrace capture & formatting | ≤300 |
| `error.go` | AppError struct & constructors (New, Wrap) | ≤300 |
| `error_diagnostic.go` | ErrorDiagnostic struct & typed setters | ≤400 |
| `error_values.go` | Values map & WithValue/WithValues setters | ≤150 |
| `clipboard.go` | AI-friendly error report formatting | ≤200 |
| `match.go` | Error matching utilities | ≤50 |
| `codes.go` | Error code constants | ≤200 |
| `result.go` | Result[T] — single value wrapper | ≤150 |
| `result_slice.go` | ResultSlice[T] — collection wrapper | ≤150 |
| `result_map.go` | ResultMap[K, V] — associative map wrapper | ≤150 |

---


---

## 1. StackTrace

### 1.1 StackFrame

```go
type StackFrame struct {
    Function string
    File     string
    Line     int
}
```

**Methods:**

- `String() string` — formats as `"function\n      file:line"`

### 1.2 StackTrace Type

```go
type StackTrace struct {
    Frames        []StackFrame
    PreviousTrace string       `json:",omitempty"`
}
```

**Fields:**

- `Frames` — ordered list of captured stack frames
- `PreviousTrace` — when two stack traces are merged (e.g., re-wrapping an error), the original trace is stored here as a formatted string

### 1.3 Capture Functions

```go
// CaptureStack captures a stack trace, skipping `skip` caller frames.
// Maximum 18 frames are captured by default.
func CaptureStack(skip int) StackTrace

// CaptureStackN captures a stack trace with a custom max frame depth.
func CaptureStackN(skip int, maxFrames int) StackTrace
```

**Rules:**

- Default max frames: **18** (sufficient for most call chains)
- `skip` parameter: number of frames to skip from the top
  - `New()` and `Wrap()` use `skip=2` (skip `runtime.Callers` + constructor)
  - `FailWrap()`, `FailSliceWrap()`, `FailMapWrap()` use `skip=3` (skip wrapper + `Wrap` + `runtime.Callers`)
- Runtime frame filtering uses `strings.HasPrefix(fn, "runtime.")` (NOT `Contains`) to avoid false positives with domain functions

### 1.4 StackTrace Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `String()` | `string` | Full formatted multi-line trace including `PreviousTrace` |
| `CallerLine()` | `string` | Top frame as `"file:line"` — compact single line |
| `FinalLine()` | `string` | Bottom frame as `"file:line"` — deepest origin point |
| `IsEmpty()` | `bool` | True if no frames captured |
| `Depth()` | `int` | Number of captured frames |
| `HasPrevious()` | `bool` | True if a previous trace exists from merging |

### 1.5 Merging Behavior

When an `AppError` is re-wrapped (wrapping an error that already has a `StackTrace`), the original trace is preserved:

```go
// Original error with trace at line 42
original := apperror.New(
    "E5001", "file not found",
)

// Re-wrapping preserves the original trace in PreviousTrace
wrapped := apperror.Wrap(
    original, "E5002", "upload failed",
)
// wrapped.Stack.HasPrevious() == true
// wrapped.Stack.PreviousTrace contains the original trace
```

---
