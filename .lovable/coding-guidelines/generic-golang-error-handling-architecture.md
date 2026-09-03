# Centralized Error Handling Architecture & AI Prompt for Go Applications

> **Target:** Universal Go Applications (CLI, Services, Daemons, Microservices)
> **Status:** Production Standard / Canonical Guide
> **Compatibility:** Go 1.20+ (Standard Library Compliant)

---

## Part 1: The AI System Prompt (Drop-In Ready)

Copy and paste this section directly into your AI configuration (`.cursorrules`, system prompts, or agent instructions) to permanently prevent anti-patterns in error handling.

```markdown

# Non-Negotiable Directive: Centralized Error Handling Architecture

You are an expert Go software engineer. You NEVER write careless, unstructured, or silent error handling code.

### The Double Anti-Pattern Rule

1. **NEVER use bare `panic("...")` or `panic(err)` in business logic or command handlers**: Crashing with raw strings or untyped runtime panics dumps unformatted stack frames to end-users and bypasses application telemetry and resource cleanup.
2. **NEVER use bare `os.Exit(...)` or silent error returns**: Terminating abruptly or failing silently destroys context, loses caller attribution, skips cleanup flushers, and makes post-incident debugging impossible.

**Both patterns represent bad engineering. All errors MUST be constructed as structured domain types (`AppError`) and dispatched through a centralized error handler (`HandleError`).**

### Core Requirements

1. **Always Structure Errors**: Every error must specify:
   - `Op`: Operation label (e.g. `user.create`, `config.load`, `db.migrate`).
   - `Code`: A registered, unique error code (e.g. `E1001`, `E2004`).
   - `Type`: A classified category enum (`VALIDATION`, `NOT_FOUND`, `PRECONDITION`, `EXECUTION`, `ABORT`, `INTERNAL`).
   - `Severity`: Impact level (`INFO`, `WARN`, `ERROR`, `FATAL`).
   - `Creator`: The subsystem/package raising the error (e.g. `cmd.setup`, `auth.jwt`).
   - `Message`: Human-readable explanation with clear remediation instructions.
   - `Ctx`: Key-value map capturing runtime parameters, variables, paths, and flags.
   - `Cause`: The underlying root cause error (for wrapping).
2. **Always Dispatch Through the Central Handler**: Pass errors to `errorhandler.HandleError(err, exitCode)` or equivalent. The central handler writes formatted diagnostics to `os.Stderr`, runs registered pipe/buffer flushers, and executes the configured exit strategy.
3. **Never Be Silent**: Swallowed errors or silent terminations are auto-reject violations. Every failure must produce actionable diagnostic breadcrumbs.

### Mandatory Pre-Commit Checklist

- [ ] Zero instances of `panic("string")` or raw `panic(err)` in production code.
- [ ] Zero instances of bare `os.Exit(...)` bypassing the central error handler.
- [ ] All error sites wrap or construct `AppError` with full metadata.
- [ ] Central error dispatcher drains buffers before exiting.
- [ ] Unit tests verify diagnostic formatting and unwrapping.
```

---

## Part 2: Detailed Analysis of the Anti-Patterns

When developers or AI assistants refactor error handling code, they often commit the mistake of replacing one anti-pattern with another:

```diff
- panic("fatal error")
+ os.Exit(1)
```

### 1. Why `panic("fatal error")` is an Anti-Pattern

* **Zero Operational Context**: The magic string `"fatal error"` tells neither the user nor the developer what actually went wrong, which inputs were invalid, or how to resolve the issue.
* **Polluted Terminal Output**: An unhandled panic prints raw goroutine stack traces and memory addresses, which confuses end-users and looks unprofessional in production software.
* **Bypasses Application Telemetry**: Telemetry pipelines, audit tables, and structured logging sinks never get invoked.

### 2. Why Bare `os.Exit(1)` is Equally Flawed

* **Silent Failure**: Calling `os.Exit(1)` abruptly halts execution without explaining the root cause.
* **Destroys Deferred Cleanups**: In Go, `os.Exit` **does not run deferred functions** (`defer`). Open file handles, temporary scratch directories, database connections, and buffered I/O streams are terminated immediately without flushing.
* **Loss of Domain Metadata**: The failure is never captured into a domain model. The error code, caller attribution, and contextual variables are dropped.
* **Fractured Control Flow**: Exit policies are hardcoded across dozens of arbitrary source files rather than managed by a single central policy.

---

## Part 3: Centralized Architecture Overview

```
┌────────────────────────────────────────────────────────┐
│                   Domain / Caller Layer                │
│                                                        │
│   Constructs structured AppError:                      │
│   • Op:       "config.load"                            │
│   • Code:     "E1001"                                  │
│   • Type:     ErrorTypeNotFound                        │
│   • Severity: SeverityError                            │
│   • Creator:  "config"                                 │
│   • Message:  "Configuration file does not exist"      │
│   • Ctx:      {"path": "/etc/app.yaml"}                │
│   • Cause:    underlying os.ErrNotExist                │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│             Central Handler / Dispatcher               │
│               (cliexit.HandleError)                    │
│                                                        │
│   1. Unwraps AppError or creates default envelope.     │
│   2. Formats uniform output to stderr with context.    │
│   3. Runs registered pipe flushers (runFlushers).      │
│   4. Records structured failure in database/audit.     │
│   5. Executes configured strategy:                     │
│      • CLI Mode: Clean process exit with exit code     │
│      • Debug Mode: Panic with rich diagnostic payload  │
│      • API Mode: Serialize to Universal Envelope JSON  │
└────────────────────────────────────────────────────────┘
```

---

## Part 4: Step-by-Step Implementation Guide

To implement this architecture in any Go application, create two lightweight, reusable packages:

### Package 1: `apperror` (Domain Error Envelope)

#### File 1.1: `apperror/types.go`

Defines the standard error categorization and severity taxonomy.

```go
package apperror

// ErrorType identifies the category of an application error.
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION"
	ErrorTypePrecondition ErrorType = "PRECONDITION"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeExecution    ErrorType = "EXECUTION"
	ErrorTypeAbort        ErrorType = "ABORT"
	ErrorTypeInternal     ErrorType = "INTERNAL"
)

// SeverityType indicates the severity level of an error.
type SeverityType string

const (
	SeverityInfo  SeverityType = "INFO"
	SeverityWarn  SeverityType = "WARN"
	SeverityError SeverityType = "ERROR"
	SeverityFatal SeverityType = "FATAL"
)
```

#### File 1.2: `apperror/apperror.go`

Defines the rich `AppError` struct, standard error interface compliance, unwrapping, and constructors.

```go
package apperror

import (
	"errors"
	"fmt"
	"strings"
)

// AppError is a typed, domain-rich error that captures operation labels,
// creator attribution, contextual metadata, severity, and root cause.
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

// Error formats the diagnostic description of the AppError.
func (e *AppError) Error() string {
	parts := make([]string, 0, 4)
	if e.Code != "" || e.Type != "" {
		parts = append(parts, fmt.Sprintf("[%s:%s]", e.Code, e.Type))
	}
	if e.Op != "" {
		parts = append(parts, e.Op+":")
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	} else if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	if e.Creator != "" {
		parts = append(parts, fmt.Sprintf("(creator=%s)", e.Creator))
	}
	if len(e.Ctx) > 0 {
		parts = append(parts, fmt.Sprintf("(ctx=%v)", e.Ctx))
	}
	if e.Cause != nil && e.Message != "" {
		parts = append(parts, fmt.Sprintf("(cause=%v)", e.Cause))
	}

	return strings.Join(parts, " ")
}

// Unwrap enables compatibility with errors.Is and errors.As.
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithContext appends a key-value pair to the error's context map.
func (e *AppError) WithContext(key string, val any) *AppError {
	if e.Ctx == nil {
		e.Ctx = make(map[string]any)
	}
	e.Ctx[key] = val

	return e
}

// New creates a new AppError without an underlying cause.
func New(op string, code string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Ctx:      ctx,
	}
}

// NewSimple creates a basic AppError with only operation and code.
func NewSimple(op string, code string) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
	}
}

// NewWithDetails creates a fully specified AppError without cause.
func NewWithDetails(
	op, code, msg, creator string,
	errType ErrorType,
	sev SeverityType,
	ctx map[string]any,
) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     errType,
		Severity: sev,
		Creator:  creator,
		Message:  msg,
		Ctx:      ctx,
	}
}

// Wrap wraps an existing error with an operation label and context.
func Wrap(err error, op string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     "E9000",
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Ctx:      ctx,
		Cause:    err,
	}
}

// WrapWithDetails wraps an existing error with full metadata.
func WrapWithDetails(
	err error,
	op, code, msg, creator string,
	errType ErrorType,
	sev SeverityType,
	ctx map[string]any,
) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     errType,
		Severity: sev,
		Creator:  creator,
		Message:  msg,
		Ctx:      ctx,
		Cause:    err,
	}
}
```

---

### Package 2: `cliexit` / `errorhandler` (Central Dispatcher)

#### File 2.1: `cliexit/handle.go`

Implements buffer flushing, formatted stderr diagnostics, debug panic switches, and clean termination.

```go
package cliexit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"yourmodule/apperror"
)

var (
	flushMu  sync.Mutex
	flushers []func()
	exitFunc = os.Exit
)

// RegisterFlusher records a cleanup function to invoke before process exit.
func RegisterFlusher(f func()) {
	if f == nil {
		return
	}
	flushMu.Lock()
	flushers = append(flushers, f)
	flushMu.Unlock()
}

// runFlushers executes all registered flushers safely.
func runFlushers() {
	flushMu.Lock()
	defer flushMu.Unlock()
	for _, f := range flushers {
		f()
	}
}

// SetExitFunc allows unit tests to intercept os.Exit without terminating the runner.
func SetExitFunc(fn func(int)) func(int) {
	prev := exitFunc
	exitFunc = fn

	return prev
}

// HandleError processes an error through centralized logging, flushing,
// and process exit (or panic if debug mode is active).
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
		appErr = apperror.Wrap(err, "cli", nil)
	}
	WriteAppErrorReport(os.Stderr, appErr)
	runFlushers()
	if os.Getenv("APP_ERROR_PANIC") == "1" {
		panic(appErr)
	}
	exitFunc(code)
}

// WriteAppErrorReport writes formatted structured error diagnostics to an io.Writer.
func WriteAppErrorReport(w io.Writer, e *apperror.AppError) {
	fmt.Fprintf(w, "error [%s:%s] in %s: %s\n", e.Code, e.Type, e.Op, e.Message)
	if e.Creator != "" {
		fmt.Fprintf(w, "  creator: %s\n", e.Creator)
	}
	if len(e.Ctx) > 0 {
		fmt.Fprintf(w, "  context: %v\n", e.Ctx)
	}
	if e.Cause != nil && e.Message != "" {
		fmt.Fprintf(w, "  cause: %v\n", e.Cause)
	}
}
```

---

## Part 5: Before & After Usage Examples

### Example 1: Command Precondition Check

#### ❌ Anti-Pattern (Before)

```go
func runSetup(args []string) error {
    if !isConfigured() {
        fmt.Fprint(os.Stderr, "setup is not configured")
        os.Exit(1) // Silent, unformatted, hardcoded exit
    }
    return nil
}
```

#### ✅ Centralized Architecture (After)

```go
func runSetup(args []string) error {
    if !isConfigured() {
        err := apperror.NewWithDetails(
            "cmd.setup",
            "E1002",
            "application is not initialized; please run initialize first",
            "cmd.setup",
            apperror.ErrorTypePrecondition,
            apperror.SeverityError,
            map[string]any{"args": args},
        )
        cliexit.HandleError(err, 1)
    }
    return nil
}
```

---

### Example 2: Wrapped File / Database Errors

#### ❌ Anti-Pattern (Before)

```go
func loadDatabase(path string) *DB {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        panic("fatal error") // Raw string crash, unhelpful stack dump
    }
    return db
}
```

#### ✅ Centralized Architecture (After)

```go
func loadDatabase(path string) *DB {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        appErr := apperror.WrapWithDetails(
            err,
            "db.open",
            "E3001",
            "failed to open SQLite database file",
            "store.sqlite",
            apperror.ErrorTypeExecution,
            apperror.SeverityFatal,
            map[string]any{"path": path},
        )
        cliexit.HandleError(appErr, 1)
    }
    return db
}
```

---

## Part 6: Verification & Quality Assurance Commands

To verify that your Go codebase is 100% compliant and free of anti-patterns, run these automated checks:

```bash

# 1. Ensure zero instances of bare panic("...") remain

grep -rn 'panic("' .

# 2. Check for bare os.Exit outside of the central error handler

grep -rn 'os\.Exit(' .

# 3. Verify that Go packages compile and pass vet

go vet ./...

# 4. Run unit tests for apperror and cliexit packages

go test -v ./apperror ./cliexit
```
