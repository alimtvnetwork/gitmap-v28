# Go AppError — Mandatory Return Type, DisplayError & Linter Enforcement

> **Parent:** [AppError Package Reference](./01-index.md)
> **Version:** 1.0.0
> **Updated:** 2026-08-28
> **Rule IDs:** CODE-RED-026, CODE-RED-027

---

## Purpose

This document defines the **Go-specific AppError enforcement standard**. Every Go service function MUST return `*AppError` (not the bare `error` interface). Every error printed to a terminal or log output MUST first be wrapped inside an `*AppError`. The `validate-guidelines.go` linter enforces these rules automatically.

---

## Why Not Bare `error`?

The bare `error` interface is a string carrier. It loses all context the moment it crosses a function boundary:

| Property | `error` interface | `*AppError` |
|----------|------------------|-------------|
| Stack trace | None | Captured at creation |
| Error code | None | Typed constant from `apperrtype` |
| Key-value context | None | `Values` map (paths, IDs, URLs) |
| User-facing message | None | `DisplayError` field |
| Structured diagnostics | None | `ErrorDiagnostic` (endpoint, method, slug) |
| JSON serialization | Loses all fields | Full round-trip preserved |
| AI debugging (clipboard) | Useless string | `ToClipboard()` — full reproduction context |
| Cross-boundary transport | Breaks | HTTP response, subprocess protocol, DB |

**Rule I-4 (non-negotiable):** Service methods return `*AppError` or `Result[T]`, never bare `error`.

---

## The AppError Struct (Canonical — Go)

`go
// AppError is the canonical error type for all Go service methods.
// Invariant I-3: Every AppError captures a stack trace at creation.
// Invariant I-4: Service methods return *AppError, never bare error.
type AppError struct {
    Code         string            // Error code from apperrtype enum (e.g. "E2010")
    Message      string            // Developer-facing technical description
    DisplayError string            json:",omitempty" // User-facing safe message for UI/terminal output
    Details      string            json:",omitempty" // Extra context (auto-set from cause on Wrap)
    Values       map[string]string json:",omitempty" // Variable context (path, ID, URL, slug, etc.)
    Diagnostic   ErrorDiagnostic   json:",omitempty" // Typed structured fields (endpoint, method, statusCode)
    Stack        StackTrace                          // Mandatory stack trace (always captured)
    Cause        error             json:"-"          // Wrapped underlying error — EXEMPTED per I-2
}
`

### The `DisplayError` Field

`DisplayError` carries the **user-facing message** — a clean, human-readable string safe to show in a terminal, UI dialog, or REST API response body.

| Property | `Message` | `DisplayError` |
|----------|-----------|----------------|
| Audience | Developer / logs | End user / terminal / UI |
| Technical detail | Yes — stack, paths, IDs | No — no internal paths or codes |
| JSON serialization | Always present | omitempty — only when set |
| Example | `"failed to read plugin file: /var/www/mu-plugins/foo.php"` | `"Could not load the plugin. Please try again."` |

**Rule:** Never print `Message` directly to the terminal. Always set `DisplayError` and print that instead.

---

## Constructors

Every constructor captures a stack trace automatically using `runtime.Callers`.

`go
// New — new AppError, developer message only.
func New(code, message string) *AppError

// NewDisplay — new AppError with BOTH developer message and user-facing display message.
func NewDisplay(code, message, displayError string) *AppError

// Wrap — wraps an existing error, preserving the original stack in PreviousTrace.
func Wrap(cause error, code, message string) *AppError

// WrapDisplay — wraps an existing error AND sets a user-facing display message.
func WrapDisplay(cause error, code, message, displayError string) *AppError

// WithDisplayError — fluent setter to add/override the user-facing display message.
func (e *AppError) WithDisplayError(displayError string) *AppError
`

### Constructor Decision Tree

`
Is there an underlying error to wrap?
  YES -> use Wrap* or WrapType*
  NO  -> use New* or NewType*

Will this error be shown to the user (terminal / CLI / UI)?
  YES -> use *Display variant, set displayError
  NO  -> standard constructor is fine

Is there a typed apperrtype enum for this case?
  YES -> use NewType / WrapType / WrapTypeMsg (preferred)
  NO  -> use New / Wrap with a raw code string
`

---

## Reference Code: Stack Capture via codestack Package

The `codestack` package from `gitlab.com/auk-go/core` provides helpers AppError uses internally for stack trace capture.

### codestack Reference Files

`go
// Source: gitlab.com/auk-go/core/-/raw/develop/codestack/funcs.go
package codestack

type (
    FilterFunc func(trace *Trace) (isTake, isBreak bool)
    Formatter  func(trace *Trace) (output string)
)
`

`go
// Source: gitlab.com/auk-go/core/-/raw/develop/codestack/fileGetter.go
// PathLineSep returns the caller file path and line number, used internally by AppError.
func (it fileGetter) PathLineSep(skipStack int) (
    filePath string, lineNumber int,
) {
    stack := New.Create(Skip1 + skipStack)
    fileWithLine := stack.FileWithLine()
    filePath = fileWithLine.FullFilePath()
    lineNumber = fileWithLine.LineNumber()
    stack.Dispose()

    return filePath, lineNumber
}
`

`go
// Source: gitlab.com/auk-go/core/-/raw/develop/codestack/dirGetter.go
// CurDir returns the directory of the file where the call was made.
func (it dirGetter) CurDir() string {
    _, filePath, _, isOkay := runtime.Caller(defaultInternalSkip)

    if isOkay {
        return filepath.Dir(filePath)
    }

    return constants.EmptyString
}
`

`go
// Source: gitlab.com/auk-go/core/-/raw/develop/codestack/currentNameOf.go
// AllStackSkip resolves full method name, package name, and method name at a given stack depth.
func (it currentNameOf) AllStackSkip(stackSkipIndex int) (
    fullMethodName, packageName, methodName string,
) {
    pc, _, _, _ := runtime.Caller(stackSkipIndex + defaultInternalSkip)
    funcInfo := runtime.FuncForPC(pc)
    fullFuncName := funcInfo.Name()

    return it.All(fullFuncName)
}
`

**The skip parameter:** AppError constructors use `skip=2` internally (`runtime.Callers` + constructor itself). You never set skip manually — the constructors handle it.

---

## Enforced Linter Rules

### CODE-RED-026: No bare `error` return in Go service functions

**Trigger:** Any Go function (excluding test files, interface stubs, `main`) whose return list contains the bare `error` interface.

**Violation patterns:**

`go
// BANNED — returns bare error
func LoadPlugin(slug string) error { }

// BANNED — tuple return with error
func FetchSite(id int64) (*Site, error) { }

// BANNED — multiple returns ending in error
func ParseConfig(path string) (Config, bool, error) { }
`

**Correct patterns:**

`go
// CORRECT — returns *AppError
func LoadPlugin(slug string) *apperror.AppError { }

// CORRECT — uses Result[T] wrapper
func FetchSite(id int64) apperror.Result[*Site] { }

// CORRECT — *AppError alongside other values
func ParseConfig(path string) (Config, bool, *apperror.AppError) { }

// EXEMPT — interface implementation
func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error { }

// EXEMPT — test files (*_test.go)
func TestLoadPlugin_RejectsEmptySlug(t *testing.T) { }
`

---

### CODE-RED-027: Raw terminal/log output without AppError wrapping

**Trigger:** Any call to a terminal output or logging function passing a raw `error` or `err` variable directly.

**Violation patterns:**

`go
// BANNED — raw err printed to terminal
fmt.Println(err)
fmt.Printf("error: %v\n", err)
fmt.Fprintf(os.Stderr, "failed: %s\n", err.Error())
log.Println(err)
log.Printf("error: %v", err)
log.Fatal(err)
`

**Correct pattern:**

`go
// CORRECT — craft AppError first, then display DisplayError to terminal
appErr := apperror.WrapDisplay(
    err,
    apperrtype.PluginNotFound,
    "failed to load plugin: " + slug,
    "Could not load the plugin. Please check the slug and try again.",
).WithSlug(slug)

log.Error(appErr.FullString())        // Internal structured log
fmt.Fprintln(os.Stderr, appErr.DisplayError) // Terminal: user-safe message only
`

---

## Terminal Output Flow (Mandatory Pattern)

`
Error occurs in service layer
        |
        v
apperror.NewDisplay() or apperror.WrapDisplay()
        |
        |-- .Message        stored in AppError.Message (developer detail)
        |-- .DisplayError   stored in AppError.DisplayError (user-safe)
        |-- .Stack          stack trace captured automatically
        |-- .Values / .Diagnostic  contextual key-value pairs
        |
        v
AppError propagated UP the call chain without modification
        |
        v
CLI / HTTP handler (the BOUNDARY layer — only place output occurs)
        |
        |-- log.Error(appErr.FullString())              internal log
        |-- fmt.Fprintln(os.Stderr, appErr.DisplayError) terminal output
`

Service layer NEVER prints to terminal. Only the boundary layer outputs.

---

## Exemptions

| Exemption | Reason |
|-----------|--------|
| `*_test.go` files | Test helpers may use bare error for assertions |
| `main()` entry points | May call `log.Fatal(err)` exactly once at top level |
| Interface implementations | Must match declared interface signature |
| `init()` functions | Startup panics only — use `appErr.Panic()` |
| `internal/adapters/` | stdlib adapter wrappers may return `error` |
| `errors.Is()` / `errors.As()` | Error inspection — not service returns |

---

## Integration with validate-guidelines.go

`bash

# Full validation including AppError return type enforcement

go run linter-scripts/validate-guidelines.go --path . --max-lines 15

# Go files only, JSON output

go run linter-scripts/validate-guidelines.go --path ./services --json
`

The linter checks:

1. **CODE-RED-026:** Scans every Go function signature for `error` in the return list (excluding exemptions).
2. **CODE-RED-027:** Scans for `fmt.Println`, `fmt.Printf`, `fmt.Fprintf`, `log.Println`, `log.Printf`, `log.Fatal`, `log.Fatalf` calls where the argument contains a raw error variable.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| AppError Struct | [02-apperror-struct.md](./02-apperror-reference.md) |
| StackTrace | [01-overview-and-stack.md](./01-index.md) |
| Result types | [03-result-types.md](./01-index.md) |
| AppErrType Enums | [05-apperrtype-enums.md](./02-apperror-reference.md) |
| Linter Script | [linter-scripts/validate-guidelines.go](../../../../linter-scripts/validate-guidelines.go) |
| golangci-lint Config | [linters/golangci-lint/.golangci.yml](../../../../linters/golangci-lint/.golangci.yml) |
| codestack package | `gitlab.com/auk-go/core/codestack` |
| Error Management Overview | [02-spec/03-error-manage/01-index.md](../../../01-index.md) |
