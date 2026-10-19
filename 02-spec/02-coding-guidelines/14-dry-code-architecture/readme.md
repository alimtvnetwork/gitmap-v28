# Golang DRY Code Architecture & Generic Package Standards (AI Execution Prompt)

> **/goal** Master and enforce DRY (Don't Repeat Yourself) code architecture principles across Golang packages, centralizing terminal rendering, ANSI styling, progress indicators, and concurrency worker pools into reusable generic packages.
> **/learn** Internalize the generic terminal box header framework (`termgeneric/header.go`), zero-dependency ANSI coloring engine (`termgeneric/color.go`), thread-safe progress and spinner components (`termgeneric/progress.go`), and generic concurrency worker pools (`termgeneric/worker_pool.go`).

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Centralize repetitive terminal UI formatting, box headers, and spinners into `cli/termgeneric/` instead of re-implementing ad-hoc ANSI logic across CLI commands.
- [ ] `/learn` Never use raw string length (`len(s)`) or byte counts for terminal column calculations; always use ANSI-stripped `runewidth.StringWidth(clean)` to prevent line-wrapping staircases.
- [ ] `/goal` Use `termgeneric.WorkerPool[T, R]` or `termgeneric.RunBatch` for parallel task processing rather than spawning unmanaged goroutines without queue bounds or cancellation contexts.
- [ ] `/learn` Enforce clean terminal line-clearing mechanics using `\r\033[2K` and thread-safe `sync.Mutex` state protection in all progress and spinner implementations.
- [ ] `/goal` Verify package test suite compliance via `go test ./termgeneric -v` (Cwd `cli/`) ensuring 100% pass rate.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 1.0.0  
**Updated:** 2026-10-19  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Ambiguity:** None  

---

## Purpose

This specification defines the DRY (Don't Repeat Yourself) code architecture standards for Golang packages in GitMap, establishing reusable patterns for terminal UI presentation, ANSI escape handling, and generic concurrent worker pools.

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

## Keywords

`golang` · `dry-architecture` · `generic-packages` · `termgeneric` · `ansi-color` · `worker-pool` · `progress-bar`

---

## 1. Principles of Golang Package DRY Architecture

Repetitive boilerplate across CLI commands compromises maintainability, introduces subtle formatting bugs, and complicates performance tuning. The DRY architecture enforces the following non-negotiable guidelines:

1. **Zero External Dependencies for Foundation Modules:** Foundational utilities (such as ANSI coloring, string stripping, and basic terminal borders) must rely exclusively on the Go standard library or repository-standard packages (`mattn/go-runewidth`).
2. **Generic Type Safety (`[T any, R any]`):** Concurrency pools and data pipelines must leverage Go generics rather than `interface{}` or `any` type assertions, ensuring compile-time type safety.
3. **Discrete Single Responsibility:** Separate terminal layout (`header.go`), color formatting (`color.go`), progress animation (`progress.go`), and concurrency execution (`worker_pool.go`) into dedicated, modular files under 300 lines.
4. **Context & Resource Propagation:** All concurrent worker pools must accept and propagate `context.Context` to guarantee instantaneous cancellation without goroutine leaks.

---

## 2. Declarative Box Header Framework (`header.go`)

### 2.1 Border Styles and Declarative Configuration

All terminal header blocks must use declarative structs rather than manual multi-line `fmt.Println` sequences:

```go
type BorderStyle struct {
    TopLeft     string
    TopRight    string
    BottomLeft  string
    BottomRight string
    Horizontal  string
    Vertical    string
}
```

Predefined styles include:
- `BorderRounded` (`╭`, `╮`, `╰`, `╯`, `─`, `│`): Default standard for GitMap UI headers.
- `BorderSharp` (`┌`, `┐`, `└`, `┘`, `─`, `│`): High-density technical summaries.
- `BorderDouble` (`╔`, `╗`, `╚`, `╝`, `═`, `║`): Critical alert banners.
- `BorderAscii` (`+`, `+`, `+`, `+`, `-`, `|`): Dumb terminal fallback.

### 2.2 ANSI-Aware Width & Title Truncation

Never calculate terminal widths using `len(s)` or `utf8.RuneCountInString(s)` because:
1. ANSI escape codes add non-visible bytes.
2. Wide Unicode runes and emojis take 2 terminal columns.

The `termgeneric` package computes visible width accurately:
```go
func VisibleWidth(s string) int {
    clean := StripAnsi(s)
    return runewidth.StringWidth(clean)
}
```

Titles exceeding available inner width are truncated with `...` to prevent ConHost/Terminal auto-wrapping to the next line.

---

## 3. Clean ANSI Coloring Engine (`color.go`)

### 3.1 Zero-Dependency Wrappers

Color functions wrap text cleanly without pulling in heavyweight terminal libraries:
- `Red(s string) string`
- `Green(s string) string`
- `Cyan(s string) string`
- `Yellow(s string) string`
- `Magenta(s string) string`
- `Blue(s string) string`
- `Dim(s string) string`
- `Bold(s string) string`
- `Reset(s ...string) string`

### 3.2 Global Color Toggling & Pipeline Stripping

For automated CI environments and non-TTY pipes, color output can be toggled via `SetColorEnabled(false)`. Furthermore, `StripAnsi(s)` strips both CSI (`\x1b[...m`, `\x1b[...K`) and OSC sequences without allocations when no escapes exist.

---

## 4. Pure Go Terminal Progress Bar & Spinner (`progress.go`)

### 4.1 Clean Line-Clearing Mechanics

To prevent staircase artifacts where progress updates print on successive rows, `progress.go` strictly uses `LineClearSequence = "\r\033[2K"`:
- `\r`: Returns the cursor to column 0.
- `\033[2K`: Erases the entire current line.

### 4.2 Thread Safety and Width Safety

All state mutations on `Spinner` and `ProgressBar` are guarded with `sync.Mutex`. Progress bar columns are clamped to `termWidth - 4` right-margin to avoid edge auto-wrap.

---

## 5. Generic Worker Pool Concurrency Pattern (`worker_pool.go`)

### 5.1 Architecture

The `WorkerPool[T, R]` provides bounded concurrency with zero thread leaks:

```go
type WorkerPool[T any, R any] struct {
    workerCount int
    queueSize   int
    handler     TaskFunc[T, R]
    tasks       chan T
    results     chan TaskResult[T, R]
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
    errMu       sync.Mutex
    errors      []error
}
```

### 5.2 Key Capabilities
- **Bounded Buffer:** Limits queued items to prevent out-of-memory errors on massive inputs.
- **Graceful Shutdown:** `Close()` flushes pending tasks and closes channels safely.
- **Immediate Cancellation:** `Stop()` cancels the context and halts worker routines promptly.
- **Error Aggregation:** Thread-safe collection via `Errors()`, `HasErrors()`, and `ErrorCount()`.
- **Batch Convenience:** `RunBatch[T, R](ctx, items, workerCount, handler)` for one-line concurrent processing.

---

## 6. Code Examples: Bad vs. Good

### ❌ BAD: Ad-Hoc Unbounded Goroutines with Manual ANSI
```go
// ❌ Manual goroutine leak, no context, manual unpadded ANSI
for _, repo := range repos {
    go func(r string) {
        fmt.Printf("\033[32mSyncing %s\033[0m\n", r)
        syncRepo(r)
    }(repo)
}
```

### ✅ GOOD: Declarative DRY Pattern with `termgeneric`
```go
// ✅ Bounded generic worker pool with context and clean UI
results, err := termgeneric.RunBatch(ctx, repos, 4, func(ctx context.Context, repo string) (string, error) {
    return syncRepo(ctx, repo)
})
if err != nil {
    fmt.Println(termgeneric.Red("Sync failed: " + err.Error()))
    return
}
fmt.Println(termgeneric.RenderBox("Sync Complete", 50))
```

---

## 7. Verification & Acceptance Criteria

### AC-CG-DRY-001: Unit Test Verification
**Given** The `cli/termgeneric` package implementation.  
**When** Running targeted package tests via `go test ./termgeneric -v` inside `cli/`.  
**Then** All test cases pass with exit code 0 and zero race conditions.

### AC-CG-DRY-002: Zero External Dependencies
**Given** The `cli/termgeneric/color.go` implementation.  
**When** Inspected for imports.  
**Then** Zero third-party dependencies are imported outside standard library and `go-runewidth`.

### AC-CG-DRY-003: Generic Concurrency Type Safety
**Given** `WorkerPool[T, R]` and `RunBatch[T, R]`.  
**When** Parameterized with domain types (e.g., `[string, Result]`).  
**Then** Handlers, inputs, and outputs are strictly typed without raw interface assertions.
