# Centralized Error Handling Architecture & Anti-Pattern Elimination

> **Version:** 1.0.0  
> **Status:** Mandatory / Canonical  
> **Target:** Cross-Stack (Go, TypeScript, Python, C#, Rust)

---

## 1. Executive Summary & The Problem

In modern software systems, improper error handling is one of the leading causes of production outages, un-debuggable incidents, and degraded user experience. 

When encountering abnormal conditions, poorly designed code typically falls into one of two destructive extremes:
1. **The Chaotic Crash Anti-Pattern** (e.g. `panic("fatal error")`, `throw "error"`): Dumps raw runtime internals, frame pointers, or meaningless strings directly to the user, bypassing structured telemetry and resource cleanup.
2. **The Silent Abort Anti-Pattern** (e.g. bare `os.Exit(1)`, `sys.exit(1)`, swallowed `catch {}`): Abruptly halts process execution or swallows the error, destroying all diagnostic context, failing to execute cleanup hooks, and leaving developers with zero clue as to why the operation failed.

**Both patterns are fundamentally flawed.** Quality engineering demands a **Centralized Error Handling Architecture** where errors are structured, attributed, contextualized, and dispatched through a single, configurable handler.

---

## 2. Analysis of the Anti-Patterns (Visual Diff Breakdown)

Consider a typical code review diff encountered during refactoring:

```diff
- panic("fatal error")
+ os.Exit(1)
```

### Why the Original Code (`panic("fatal error")`) is Wrong

* **Magic String Noise**: `"fatal error"` provides zero operational intelligence. It fails to answer: *Which command failed? What parameters were provided? What precondition was violated?*
* **Ugly Runtime Stack Dump**: In production CLI tools, an unhandled panic spills runtime stack frames across the user's terminal, eroding user trust.
* **Bypasses Application Telemetry**: The application cannot record the failure into audit logs, database journals, or metrics pipelines.

### Why the AI Modification (`os.Exit(1)`) is Equally Wrong

* **Silent Termination**: Calling `os.Exit(1)` directly terminates execution without printing structured diagnostic information.
* **Loss of Domain Metadata**: The failure is never wrapped in an `AppError`. Information regarding who created the error, the error code (`E...`), the error category, and parameter context is discarded.
* **Bypasses Lifecycle Cleanups**: Calling `os.Exit` immediately aborts the process, preventing deferred cleanup hooks (e.g., closing file descriptors, flushing buffered stdout/stderr pipes, releasing database locks) from executing.
* **Fractured Error Policies**: Exit behavior is scattered across dozens of arbitrary source files rather than managed by a single central policy.

---

## 3. The Centralized Error Architecture Standard

Centralized error handling separates **Error Construction** from **Error Handling/Dispatch**.

```
┌────────────────────────────────────────────────────────┐
│                   Domain / Caller Layer                │
│                                                        │
│   Constructs structured AppError:                      │
│   - Op: "command.reinstall"                            │
│   - Code: "E2002"                                      │
│   - Type: ErrorTypeAbort                               │
│   - Severity: SeverityWarn                             │
│   - Creator: "cmd.reinstall"                           │
│   - Message: "Reinstall aborted by user"               │
│   - Ctx: {"mode": "auto", "confirmed": false}          │
│   - Cause: underlying error (if wrapping)              │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│             Central Handler / Dispatcher               │
│               (e.g. cliexit.HandleError)               │
│                                                        │
│   1. Inspects error type, code, and severity.          │
│   2. Formats uniform output to stderr with context.    │
│   3. Drains pipe buffers and flushes log sinks.        │
│   4. Records structured failure in database/audit.     │
│   5. Executes configured strategy:                     │
│      • CLI Mode: Clean process exit with exit code     │
│      • Debug Mode: Panic with rich diagnostic payload  │
│      • API Mode: Serialize to Universal Envelope JSON  │
└────────────────────────────────────────────────────────┘
```

### Core Requirements of `AppError`

Every structured domain error MUST contain:
1. **`Op` (Operation)**: Semantic action label (e.g., `git.clone`, `auth.verify`, `db.query`).
2. **`Code` (Error Code)**: Unique, registered error identifier (e.g., `E1001`, `E2004`).
3. **`Type` (Error Category)**: Enum classifying the failure mode (`VALIDATION`, `NOT_FOUND`, `PRECONDITION`, `EXECUTION`, `ABORT`, `INTERNAL`).
4. **`Severity`**: Impact level (`INFO`, `WARN`, `ERROR`, `FATAL`).
5. **`Creator` / `Attribution`**: The subsystem or component that detected the error (e.g., `cmd.rootadd`, `auth.jwt`).
6. **`Message`**: Clear, human-readable description of what went wrong and how to fix it.
7. **`Ctx` (Context Map)**: Key-value bag containing critical runtime arguments, paths, and flags.
8. **`Cause` (Underlying Error)**: The root cause error being wrapped (preserving full unwrapping capability).

---

## 4. Universal "Never Be Silent" Policy

1. **No Swallowed Errors**: Catching or receiving an error and doing nothing, or printing a generic string and exiting with `os.Exit(1)`, is strictly banned.
2. **Always Provide Actionable Diagnostics**: Every error emitted to the console or log must clearly tell the developer or end-user *what failed*, *why it failed*, and *how to resolve it*.
3. **Pluggable Exit Strategies**: The central dispatcher (`HandleError`) determines how to terminate based on the environment:
   - In production CLI mode: Writes formatted diagnostics to stderr, runs cleanup flushers, and exits with a deterministic exit code.
   - In developer / debug mode (`GITMAP_ERROR_PANIC=1`): Emits full diagnostics and raises a panic with the structured `AppError` for interactive debugger inspection.
   - In HTTP/API mode: Serializes the `AppError` into a universal response envelope (`{ "data": null, "errors": [...] }`).

---

## 5. Reference Implementations

### Go Implementation

```go
// 1. Structured AppError
type AppError struct {
    Op       string
    Code     string
    Type     ErrorType
    Severity SeverityType
    Creator  string
    Message  string
    Ctx      map[string]any
    Cause    error
}

// 2. Caller Usage
func runReleasePull(args []string) error {
    if !release.IsInsideGitRepo() {
        err := apperror.NewWithDetails(
            "release.pull",
            "E2001",
            "not inside a valid git repository",
            "cmd.releasepull",
            apperror.ErrorTypePrecondition,
            apperror.SeverityError,
            map[string]any{"cwd": cwd},
        )
        cliexit.HandleError(err, 1)
    }
    return nil
}

// 3. Central Handler
func HandleError(err error, defaultCode ...int) {
    if err == nil {
        return
    }
    code := 1
    if len(defaultCode) > 0 && defaultCode[0] > 0 {
        code = defaultCode[0]
    }
    var appErr *apperror.AppError
    if !errors.As(err, &appErr) {
        appErr = apperror.WrapSimple(err, "cli")
    }
    WriteAppErrorReport(os.Stderr, appErr)
    runFlushers()
    if os.Getenv("GITMAP_ERROR_PANIC") == "1" {
        panic(appErr)
    }
    exitFunc(code)
}
```

### TypeScript / Frontend Implementation

```typescript
export class AppError extends Error {
  public readonly code: string;
  public readonly type: ErrorType;
  public readonly creator: string;
  public readonly context: Record<string, unknown>;
  public readonly cause?: Error;

  constructor(options: {
    op: string;
    code: string;
    type: ErrorType;
    creator: string;
    message: string;
    context?: Record<string, unknown>;
    cause?: Error;
  }) {
    super(`[${options.code}:${options.type}] ${options.op}: ${options.message}`);
    this.name = "AppError";
    this.code = options.code;
    this.type = options.type;
    this.creator = options.creator;
    this.context = options.context ?? {};
    this.cause = options.cause;
  }
}

export function handleError(error: unknown): void {
  const appError = error instanceof AppError 
    ? error 
    : new AppError({
        op: "unknown",
        code: "E9000",
        type: "INTERNAL",
        creator: "globalHandler",
        message: String(error),
        cause: error instanceof Error ? error : undefined,
      });

  // Log to structured telemetry / global error modal store
  errorStore.captureError(appError);
}
```
