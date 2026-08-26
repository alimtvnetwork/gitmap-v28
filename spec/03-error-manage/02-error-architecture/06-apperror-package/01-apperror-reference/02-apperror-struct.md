# AppError Package Reference — AppError struct and constructors

> **Parent:** [AppError Package Reference](./00-overview.md)  
> **Version:** 2.0.0  
> **Updated:** 2026-04-02

---

## 2. AppError

### 2.1 Struct

```go
type AppError struct {
    Code       string
    Message    string
    Details    string            `json:",omitempty"`
    Values     map[string]string `json:",omitempty"`
    Diagnostic ErrorDiagnostic   `json:",omitempty"`
    Stack      StackTrace
    Cause      error             `json:"-"` // EXEMPTED: AppError internal cause (I-2)
}
```

**Fields:**
- `Code` — error code from constants (e.g., `ErrNotFound`, `ErrDatabaseQuery`)
- `Message` — human-readable error description
- `Details` — additional context (auto-set from cause on `Wrap`)
- `Values` — key-value map for injecting variables relevant to the error context (paths, IDs, names, etc.)
- `Diagnostic` — typed diagnostic fields for structured reporting
- `Stack` — mandatory stack trace captured at creation
- `Cause` — wrapped underlying error (implements `Unwrap()`)

### 2.2 Constructors

Every constructor captures a stack trace automatically. **Three things are always required: cause (or nil), code, and message.**

```go
// New creates a new AppError with code + message. Stack captured at caller.
func New(code, message string) *AppError

// NewWithSkip creates a new AppError with explicit skip for stack capture.
func NewWithSkip(code, message string, skip int) *AppError

// Wrap wraps an existing error with code + message. Stack captured at caller.
// If cause is an *AppError, its stack is preserved in PreviousTrace.
func Wrap(cause error, code, message string) *AppError

// WrapWithSkip wraps with explicit skip for stack capture.
func WrapWithSkip(cause error, code, message string, skip int) *AppError

// NewType creates a new AppError from an apperrtype.Variation enum.
// Code, message, and name are pulled from the global variantRegistry.
func NewType(errType apperrtype.ErrorType) *AppError

// WrapType wraps an existing error using an apperrtype.Variation enum.
// Uses the registry's default message. cause's stack is preserved.
func WrapType(cause error, errType apperrtype.ErrorType) *AppError

// WrapTypeMsg wraps an existing error with a Variation enum + custom message.
// Overrides the registry's default message with the provided one.
func WrapTypeMsg(cause error, errType apperrtype.ErrorType, message string) *AppError
```

**`NewType` / `WrapType` / `WrapTypeMsg` examples:**

```go
// NewType — new error from enum (no cause)
return apperror.NewType(apperrtype.SiteNotFound)

// WrapType — wrap with enum's default message (2 args: cause + errType)
return apperror.WrapType(err, apperrtype.WPConnectionFailed).
    WithValue("url", siteURL)

// WrapTypeMsg — wrap with custom message (3 args: cause + errType + msg)
return apperror.WrapTypeMsg(err, apperrtype.WPConnectionFailed, "failed during health check").
    WithValue("url", siteURL).
    WithStatusCode(resp.StatusCode)
```

### 2.2.1 Path Convenience Constructors

Shorthand constructors for file system errors — automatically set the `path` diagnostic and use the appropriate `Variation`:

```go
// PathError creates a path-related AppError with the given Variation.
// Automatically sets WithPath(path) diagnostic.
func PathError(errType apperrtype.ErrorType, path string) *AppError

// WrapPathError wraps a cause with a path-related Variation.
// Automatically sets WithPath(path) diagnostic.
func WrapPathError(cause error, errType apperrtype.ErrorType, path string) *AppError
```

**Usage with path variants:**

```go
// Path validation — no underlying error
if path == "" {
    return apperror.FailBool(apperror.PathError(apperrtype.EmptyFilePath, path))
}

if !isValidPath(path) {
    return apperror.FailBool(apperror.PathError(apperrtype.PathInvalid, path))
}

// Wrapping an OS error with path context
data, err := os.ReadFile(path)
if err != nil {
    return apperror.FailBytes(apperror.WrapPathError(err, apperrtype.PathFailedToRead, path))
}

// Creating a directory
if err := os.MkdirAll(dir, 0755); err != nil {
    return apperror.FailBool(apperror.WrapPathError(err, apperrtype.PathFailedToCreate, dir))
}
```

**Available path variants:**

| Variant | Code | Default Message |
|---------|------|-----------------|
| `PathNotFound` | E4012 | path not found |
| `PathInvalid` | E4013 | invalid path |
| `PathStatFailed` | E4014 | failed to stat path |
| `PathMissing` | E4016 | required path is missing |
| `PathFailedToCreate` | E4017 | failed to create path |
| `PathFailedToRead` | E4018 | failed to read path |
| `PathFailedToWrite` | E4019 | failed to write to path |
| `PathFailedToDelete` | E4020 | failed to delete path |
| `EmptyFilePath` | E4011 | file path is empty |

### 2.2.2 URL Convenience Constructors

Shorthand for network/API errors — automatically sets the `url` diagnostic:

```go
// UrlError creates a URL-related AppError. Auto-sets WithUrl(url).
func UrlError(errType apperrtype.ErrorType, url string) *AppError

// WrapUrlError wraps a cause with a URL-related Variation. Auto-sets WithUrl(url).
func WrapUrlError(cause error, errType apperrtype.ErrorType, url string) *AppError
```

**Usage:**

```go
// WordPress API connection failure
return apperror.FailSettings(
    apperror.WrapUrlError(err, apperrtype.WPConnectionFailed, siteURL),
)

// Invalid URL — no underlying error
return apperror.FailBool(
    apperror.UrlError(apperrtype.RequestFailed, endpoint),
)
```

### 2.2.3 Slug Convenience Constructors

Shorthand for plugin/resource slug errors — automatically sets the `slug` diagnostic:

```go
// SlugError creates a slug-related AppError. Auto-sets WithSlug(slug).
func SlugError(errType apperrtype.ErrorType, slug string) *AppError

// WrapSlugError wraps a cause with a slug-related Variation. Auto-sets WithSlug(slug).
func WrapSlugError(cause error, errType apperrtype.ErrorType, slug string) *AppError
```

**Usage:**

```go
// Plugin not found by slug
return apperror.FailBool(
    apperror.SlugError(apperrtype.PluginNotFound, pluginSlug),
)

// Plugin activation failure with underlying error
return apperror.FailBool(
    apperror.WrapSlugError(err, apperrtype.PluginAlreadyActive, pluginSlug),
)
```

### 2.2.4 Site Convenience Constructors

Shorthand for site-scoped errors — automatically sets the `siteId` diagnostic:

```go
// SiteError creates a site-related AppError. Auto-sets WithSiteId(siteId).
func SiteError(errType apperrtype.ErrorType, siteId int64) *AppError

// WrapSiteError wraps a cause with a site-related Variation. Auto-sets WithSiteId(siteId).
func WrapSiteError(cause error, errType apperrtype.ErrorType, siteId int64) *AppError
```

**Usage:**

```go
// Site not found
return apperror.FailBool(
    apperror.SiteError(apperrtype.SiteNotFound, siteId),
)

// Sync failure on a specific site
return apperror.FailBool(
    apperror.WrapSiteError(err, apperrtype.SyncTimeout, siteId),
)
```

### 2.2.5 Endpoint Convenience Constructors

Shorthand for HTTP request errors — automatically sets `endpoint`, `method`, and `statusCode` diagnostics:

```go
// EndpointError creates an HTTP-related AppError.
// Auto-sets WithEndpoint(ep), WithMethod(method), WithStatusCode(statusCode).
func EndpointError(errType apperrtype.ErrorType, method, endpoint string, statusCode int) *AppError

// WrapEndpointError wraps a cause with HTTP diagnostics.
// Auto-sets WithEndpoint(ep), WithMethod(method), WithStatusCode(statusCode).
func WrapEndpointError(cause error, errType apperrtype.ErrorType, method, endpoint string, statusCode int) *AppError
```

**Usage:**

```go
// API returned unexpected status
return apperror.FailSettings(
    apperror.EndpointError(apperrtype.WPResponseInvalid, "GET", "/wp-json/wp/v2/plugins", resp.StatusCode),
)

// Wrapping a network error with full HTTP context
return apperror.FailSettings(
    apperror.WrapEndpointError(err, apperrtype.WPConnectionFailed, "POST", endpoint, 0).
        WithValue("payload", truncatedBody),
)
```

### 2.2.6 Convenience Constructor Summary

| Constructor | Auto-sets | Typical Variants |
|-------------|-----------|------------------|
| `PathError` / `WrapPathError` | `WithPath(path)` | `PathInvalid`, `PathFailedToRead`, `PathMissing` |
| `UrlError` / `WrapUrlError` | `WithUrl(url)` | `WPConnectionFailed`, `RequestFailed`, `ConnectionFailed` |
| `SlugError` / `WrapSlugError` | `WithSlug(slug)` | `PluginNotFound`, `PluginAlreadyActive`, `PluginSlugMissing` |
| `SiteError` / `WrapSiteError` | `WithSiteId(id)` | `SiteNotFound`, `SiteBlocked`, `SyncTimeout` |
| `EndpointError` / `WrapEndpointError` | `WithEndpoint` + `WithMethod` + `WithStatusCode` | `WPEndpointNotFound`, `WPResponseInvalid`, `WPRateLimited` |

> **Rule:** If a diagnostic field is relevant, always use the convenience constructor instead of manual `.WithXxx()` chaining — it ensures no context is accidentally omitted.

### 2.2.7 Error Merge

Combine multiple `AppError` instances into a single error when an operation collects several failures (e.g., batch validation, multi-step processing):

```go
// Merge combines multiple AppErrors into a single AppError.
// The first error's code is used as the merged error's code.
// All errors are preserved in the Values map and cause chain.
func Merge(errors []*AppError) *AppError

// MergeWithCode combines multiple AppErrors under a specific error code.
func MergeWithCode(code string, message string, errors []*AppError) *AppError
```

**Usage:**

```go
// Batch validation — collect all failures, then merge
var errs []*apperror.AppError

if site == nil {
    errs = append(errs, apperror.NewType(apperrtype.SiteNotFound))
}
if pluginSlug == "" {
    errs = append(errs, apperror.NewType(apperrtype.PluginSlugMissing))
}
if configPath == "" {
    errs = append(errs, apperror.PathError(apperrtype.EmptyFilePath, configPath))
}

if len(errs) > 0 {
    return apperror.FailBool(
        apperror.MergeWithCode(apperrtype.ValidationFailed.Code(), "multiple validation errors", errs),
    )
}

// Multi-step processing — accumulate errors
var errs []*apperror.AppError
for _, site := range sites {
    if err := syncSite(site); err != nil {
        errs = append(errs, apperror.WrapSiteError(err, apperrtype.SyncConflict, site.Id))
    }
}
if len(errs) > 0 {
    merged := apperror.Merge(errs)
    log.Error(merged.FullString())
    return apperror.FailBool(merged)
}
```

**Merge output format (in `FullString()`):**

```
[E9002] multiple validation errors (3 errors merged)
  1. [E2010] site not found
  2. [E2012] plugin slug required
  3. [E4011] file path is empty — path: ""
```

### 2.3 Display Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `Error()` | `string` | `"[CODE] message"` — implements `error` interface |
| `FullString()` | `string` | Code + message + details + values + diagnostics + stack + cause chain |
| `String()` | `string` | Alias for `FullString()` — complete error representation |
| `ToClipboard()` | `string` | Markdown-formatted error report for AI paste |

### 2.3.1 Variation Display Methods

The `Variation` enum itself carries display and introspection methods — no `AppError` instance required:

```go
v := apperrtype.SiteNotFound

v.String()               // "SiteNotFound (Code - 10) : site not found"
v.CodeTypeName()         // "(#10 - SiteNotFound)"
v.CodeTypeNameWithReferences("example.com", "wp-admin")
                         // "(#10 - SiteNotFound) example.com, wp-admin"
v.Name()                 // "SiteNotFound"
v.Code()                 // "E2010"
v.Message()              // "site not found"
v.IsValid()              // true
```

| Method | Returns | Description |
|--------|---------|-------------|
| `String()` | `string` | Delegates to `TypeNameCodeMessage()` — `"Name (Code - N) : Message"` |
| `CodeTypeName()` | `string` | `"(#N - Name)"` — compact numeric identifier |
| `CodeTypeNameWithReferences(refs...)` | `string` | `"(#N - Name) ref1, ref2"` — with context |
| `Name()` | `string` | PascalCase variant name |
| `Code()` | `string` | String error code (e.g. `"E2010"`) |
| `Message()` | `string` | Human-readable default message |
| `IsValid()` | `bool` | True if between `NoError` and `MaxError` |

### 2.3.2 VariantStructure Lookup

Use `Structure()` to retrieve the full metadata from the global registry:

```go
vs := apperrtype.SiteNotFound.Structure()
// vs.Name    → "SiteNotFound"
// vs.Code    → "E2010"
// vs.Message → "site not found"
// vs.Variant → apperrtype.SiteNotFound
```

`Structure()` returns a `VariantStructure` with all fields populated from `variantRegistry`. If the variant is unregistered, it returns a zero-value struct.

`VariantStructure` has its own display methods:

| Method | Returns | Description |
|--------|---------|-------------|
| `String()` | `string` | Alias for `TypeNameCodeMessage()` |
| `TypeNameCodeMessage()` | `string` | `"Name (Code - N) : Message"` |
| `CodeTypeName()` | `string` | `"(#N - Name)"` |
| `CodeTypeNameWithMessage(msg)` | `string` | `"(#N - Name) msg"` |

### 2.3.3 Direct Error Creation from VariantStructure

`VariantStructure` (not `Variation`) can create stdlib errors and panics directly:

```go
// .Error() — returns a stdlib error with "[Code] Message: detail" format
err := apperrtype.SiteNotFound.Structure().Error("domain example.com")
// err.Error() → "[e2010] site not found: domain example.com"

// .ErrorNoRefs() — returns a stdlib error with no additional context
err := apperrtype.SiteNotFound.Structure().ErrorNoRefs()
// err.Error() → "[e2010] site not found"

// .Panic() — panics with the formatted error (startup failures only)
apperrtype.ConfigFileMissing.Structure().Panic("required key: db_host")
```

| Method | Receiver | Returns | Description |
|--------|----------|---------|-------------|
| `Error(detail)` | `VariantStructure` | `error` | `errors.New("[code] message: detail")` (lowercased) |
| `ErrorNoRefs()` | `VariantStructure` | `error` | `errors.New("[code] message")` (lowercased) |
| `Panic(detail)` | `VariantStructure` | — | Panics with `Error(detail)` |

### 2.4 Values — Variable Injection

Anytime an error occurs while working with a variable (path, ID, name, URL), that variable **must** be injected into the error's `Values` map so no context is lost.

```go
// WithValue adds a single key-value pair.
func (e *AppError) WithValue(key, value string) *AppError

// WithValues merges multiple key-value pairs.
func (e *AppError) WithValues(values map[string]string) *AppError
```

**Usage:**
```go

return apperror.Wrap(err, ErrFSRead, "failed to read plugin file").
    WithValue("path", filePath).
    WithValue("plugin", pluginSlug)
```

The `Values` map is included in `FullString()`, `String()`, and `ToClipboard()` output, compiling into a readable error message.

### 2.5 Flow Control Methods

```go
// Panic logs the full error and panics with the formatted message.
// Use ONLY for unrecoverable initialization failures.
func (e *AppError) Panic(message string)

// Throw panics with the AppError itself (recoverable via recover).
// The AppError can be extracted from the panic value.
func (e *AppError) Throw()
```

**Rules:**
- `Panic()` is reserved for startup/initialization failures only
- `Throw()` enables structured panic/recover patterns where the `AppError` is preserved
- Neither should be used in request handlers — return errors instead

### 2.6 Query Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `Unwrap()` | `error` | Returns cause for `errors.Is/As` |
| `Is(target)` | `bool` | True if error codes match |
| `HasCause()` | `bool` | True if a wrapped cause exists |
| `HasValues()` | `bool` | True if Values map is populated |
| `HasDiagnostic()` | `bool` | True if any diagnostic field is set |

### 2.7 Diagnostic Setters (Fluent)

All typed diagnostic setters return `*AppError` for chaining:

```go
err.WithPath(p string)
err.WithFile(f string)
err.WithFilePath(p string)
err.WithUrl(u string)
err.WithSlug(s string)
err.WithSiteId(id int64)
err.WithPluginId(id int64)
err.WithStatusCode(code int)
err.WithMethod(m string)
err.WithEndpoint(ep string)
err.WithUsername(u string)
// ... etc (see error_diagnostic.go for full list)
```

---
