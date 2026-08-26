# Golang Coding Standards — Concurrency, forbidden patterns, imports, common mistakes

> **Parent:** [Golang Coding Standards](./00-overview.md)  
> **Version:** 3.7.0  
> **Updated:** 2026-03-31

---

## Concurrency Patterns

### `sync.Once` for Lazy Initialization

```go
var (
    openAPISpec     []byte
    openAPISpecOnce sync.Once
)

func GetOpenAPISpec() []byte {
    openAPISpecOnce.Do(func() {
        openAPISpec, _ = pathutil.ReadFile("api/openapi.json")
    })

    return openAPISpec
}
```

### Context Propagation

All long-running operations must accept `context.Context`:

```go
func (s *PublishService) Upload(ctx context.Context, req UploadRequest) error { ... }
```

---


---

## Forbidden Patterns

| Pattern | Why | Alternative |
|---------|-----|-------------|
| `interface{}` / `any` in exported APIs | Untyped | Concrete types or generics |
| `fmt.Errorf` for service errors | No stack trace | `apperror.Wrap` |
| Raw `error` return from services | No stack trace, no code | `apperror.Result[T]` or `*apperror.AppError` |
| Panic in handlers | Crashes server | Return error |
| `init()` functions | Hidden side effects | Explicit initialization |
| Global mutable state | Race conditions | Dependency injection |
| `map[string]interface{}` in APIs | Untyped | Defined structs |
| Raw `(T, error)` from services | No semantic methods | `apperror.Result[T]` |
| Raw `(T, error)` from API call helpers | No semantic methods | `apperror.Result[T]` |
| `"POST"` / `"GET"` string literals | Magic string | `httpmethod.Post` / `httpmethod.Get` enum |
| `"snapshot cleanup"` operation strings | Magic string | Per-domain operation enum constant |
| `!fn()` raw negation | Easy to miss `!` | Positive guard function |
| `x != nil && x.IsValid()` | Compound nil+valid check | `x.IsDefinedAndValid()` |
| Nested `if` (any depth) | **Zero tolerance** | Flatten with early returns |
| Functions > 15 lines | Hard to read | Extract small helpers |
| Files > 400 lines | Hard to navigate | Split with suffix convention (target 300) |
| Magic strings/numbers | Brittle | Typed constants |
| Boolean flag parameters | Unclear intent | Separate named methods |
| `json:"FieldName"` matching field name | Redundant — Go marshals PascalCase by default | Omit tag; use `json:",omitempty"` only when needed |

---


---

## Import Organization — 3 Groups

```go
import (
    // stdlib
    "context"
    "fmt"

    // internal packages
    "project/pkg/apperror"
    "project/internal/domain"

    // third-party
    "github.com/lib/pq"
)
```

---


---

## Common Mistakes — Go

These are real violations found and fixed. Reference to avoid repeating.

### Mistake 1: snake_case in `variantLabels`

```go
// ❌ WRONG — snake_case labels
var variantLabels = [...]string{
    Invalid:  "invalid",
    PerTable: "per_table",
}

// ✅ CORRECT — PascalCase labels
var variantLabels = [...]string{
    Invalid:  "Invalid",
    PerTable: "PerTable",
}
```

### Mistake 2: `!v.IsValid()` Instead of `v.IsInvalid()`

```go
// ❌ WRONG — raw negation
func (v Variant) String() string {
    if !v.IsValid() {
        return variantLabels[Invalid]
    }

    return variantLabels[v]
}

// ✅ CORRECT — positive counterpart
func (v Variant) String() string {
    if v.IsInvalid() {
        return variantLabels[Invalid]
    }

    return variantLabels[v]
}
```

### Mistake 3: `!pathutil.IsDir()` Without Counterpart

```go
// ❌ WRONG — raw negation on utility
if !pathutil.IsDir(gitDir) {
    return apperror.FailNew[StatusResult](
        apperror.ErrGitNotRepo,
        "not a git repo",
    )
}

// ✅ CORRECT — use IsDirMissing()
if pathutil.IsDirMissing(gitDir) {
    return apperror.FailNew[StatusResult](
        apperror.ErrGitNotRepo,
        "not a git repo",
    )
}
```

### Mistake 4: `fmt.Errorf()` for Service Errors

```go
// ❌ WRONG — no stack trace, no error code
return fmt.Errorf("failed to upload: %w", err)

// ✅ CORRECT — apperror with automatic stack trace
return apperror.Wrap(
    err,
    apperror.ErrUploadFailed,
    "failed to upload plugin",
)
```

### Mistake 5: Raw `(T, error)` from Service Methods

```go
// ❌ WRONG — raw tuple, no semantic methods
func (s *PluginService) GetById(ctx context.Context, id int64) (*Plugin, error) { ... }

// ✅ CORRECT — typed result wrapper
func (s *PluginService) GetById(ctx context.Context, id int64) apperror.Result[Plugin] { ... }
```

### Mistake 5b: Raw `(T, error)` from API Call Helpers

```go
// ❌ WRONG — tuple return + magic strings
func doAPICall[T any](c *Client, input apiCallInput) (*T, error) {
    data, err := c.doAPICallRaw(input)
    if err != nil {
        return nil, err
    }

    return decodeAPIResponse[T](data, input.Operation)
}

// ✅ CORRECT — Result[T] return, blank line after closing brace
func apiCallTo[T any](c *Client, input apiCallInput) apperror.Result[T] {
    data, err := c.apiCallToRaw(input)
    if err != nil {
        return apperror.FailWrap[T](
            err,
            "E4010",
            "API call failed: "+input.Operation.Label(),
        )
    }

    return decodeApiResponse[T](data, input.Operation)
}
```

### Mistake 5c: Magic Strings for HTTP Method and Operation

```go
// ❌ WRONG — inline string literals
callInput := apiCallInput{
    Method:    "POST",
    Operation: "snapshot cleanup",
}

// ✅ CORRECT — enum constants
callInput := apiCallInput{
    Method:    httpmethod.Post,
    Operation: snapshotoperationtype.Cleanup,
}
```

### Mistake 6: `interface{}` / `any` in Business Logic

```go
// ❌ WRONG — type erasure
func ProcessData(data interface{}) interface{} { ... }

// ✅ CORRECT — concrete types
func ProcessData(data PluginDetails) apperror.Result[PluginSummary] { ... }
```

### Mistake 7: snake_case in SQL / Struct Tags After Migration

```go
// ❌ WRONG — old snake_case
const query = `SELECT plugin_slug FROM transactions`
type Tx struct {
    PluginSlug string `db:"plugin_slug"`
}

// ✅ CORRECT — PascalCase
const query = `SELECT PluginSlug FROM Transactions`
type Tx struct {
    PluginSlug string `db:"PluginSlug"`
}
```

### Mistake 8: Compound Negation Without Named Boolean

```go
// ❌ WRONG — inline negated compound
if !config.BuildEnabled || config.BuildCommand == "" {
    return apperror.FailNew[BuildResult](
        apperror.ErrBuildNotConfigured,
        "build not configured",
    )
}

// ✅ CORRECT — extract negation to positive counterpart, then compose
isBuildDisabled := !config.BuildEnabled
isBuildCommandEmpty := config.BuildCommand == ""
isBuildMissing := isBuildDisabled || isBuildCommandEmpty

if isBuildMissing {
    return apperror.FailNew[BuildResult](
        apperror.ErrBuildNotConfigured,
        "build not configured",
    )
}
```

---
