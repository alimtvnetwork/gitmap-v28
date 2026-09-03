# Casting Elimination Patterns

> **Version:** 2.0.0
> **Updated:** 2026-03-09
> **Status:** Complete
> **Priority:** Critical

---

## 1. Purpose

Documents the canonical patterns for eliminating raw type assertions (`.(Type)`) from business logic per §7.2. All casts must be centralized into typed accessor functions, with `// EXEMPTED:` annotations at the cast boundary.

---

## 2. Pattern Categories

| Category | Cast Location | Annotation |
|----------|--------------|------------|
| Context values | Typed accessor function (e.g., `GetUserId`) | `// EXEMPTED: typed context accessor internal (§7.2)` |
| Error types | `errors.As()` or `apperror.Extract()` | `// EXEMPTED: apperror package internal (§7.2)` |
| Cache (`sync.Map`) | Typed accessor wrapper | `// EXEMPTED: typed accessor internal — cache stores known T values (§7.2)` |
| External JSON/WebSocket | Immediate deserialization boundary | `// EXEMPTED: external API — dynamic JSON (§7.2)` |
| `sql.Scanner` | `Scan(src any)` implementation | `// EXEMPTED: sql.Scanner stdlib interface (§7.2)` |
| `.(type)` switch | JSON decoder boundary only | `// EXEMPTED: generic settings accessor internal (§7.2)` |

---

## 3. Context Values

### ❌ Prohibited — Raw `ctx.Value` Cast in Business Logic

```go
func HandleRequest(ctx context.Context) {
    userId := ctx.Value("user_id").(string)  // §7.2 violation
    // ...
}
```

### ✅ Required — Typed Accessor Function

```go
// In pkg/ctxutil/context.go

type contextKey string

const userIdKey contextKey = "UserId"

func WithUserId(ctx context.Context, userId string) context.Context {
    return context.WithValue(ctx, userIdKey, userId)
}

// EXEMPTED: typed context accessor internal — this IS the centralized cast location (§7.2)
func GetUserId(ctx context.Context) string {
    if v := ctx.Value(userIdKey); v != nil {
        return v.(string)
    }

    return ""
}
```

### ✅ Business Logic Consumption

```go
func HandleRequest(ctx context.Context) {
    userId := ctxutil.GetUserId(ctx)  // No cast visible
    // ...
}
```

---

## 4. Error Types

### ❌ Prohibited — Raw Type Assertion

```go
if appErr, ok := err.(*apperror.AppError); ok {
    log.Error("failed", "Code", appErr.Code)
}
```

### ✅ Required — `errors.As()` Pattern

```go
var appErr *apperror.AppError

if errors.As(err, &appErr) {
    log.Error("failed", "Code", appErr.Code)
}
```

### ✅ Required — Helper Functions for Common Checks

```go
// In pkg/apperror/helpers.go

func IsRetryable(err error) bool {
    var appErr *AppError

    if !errors.As(err, &appErr) {
        return false
    }

    return appErr.Retryable
}

func GetExitCode(err error) int {
    var appErr *AppError

    if !errors.As(err, &appErr) {
        return 1
    }

    return appErr.ExitCode
}
```

### ✅ Stdlib Boundary — Exempt with Annotation

```go
// EXEMPTED: stdlib boundary — exec.ExitError (§7.2)
var exitErr *exec.ExitError

if errors.As(err, &exitErr) {
    return exitErr.ExitCode()
}
```

---

## 5. Cache / `sync.Map`

### ❌ Prohibited — Bare Cast in Business Logic

```go
if cached, ok := s.cache.Load(key); ok {
    return cached.([]string), nil  // §7.2 violation
}
```

### ✅ Required — Annotated Typed Accessor

```go
func (s *Service) GetCachedStrings(key string) ([]string, bool) {
    // EXEMPTED: typed accessor internal — cache stores known []string values (§7.2)
    if cached, ok := s.cache.Load(key); ok {
        return cached.([]string), true
    }

    return nil, false
}
```

### ✅ Preferred — Generic Typed Cache Wrapper

```go
type TypedCache[T any] struct {
    inner sync.Map
}

func (c *TypedCache[T]) Load(key string) (T, bool) {
    // EXEMPTED: generic cache internal — single centralized cast point (§7.2)
    if v, ok := c.inner.Load(key); ok {
        return v.(T), true
    }

    var zero T

    return zero, false
}

func (c *TypedCache[T]) Store(key string, value T) {
    c.inner.Store(key, value)
}
```

---

## 6. External JSON / WebSocket

### ❌ Prohibited — Bare Casts Without Annotation

```go
var msg map[string]any
conn.ReadJSON(&msg)
text := msg["text"].(string)  // panic-prone + §7.2 violation
```

### ✅ Required — Comma-Ok with Exemption Annotation

```go
// EXEMPTED: external WebSocket API — protocol returns dynamic JSON (§7.2)
var msg map[string]any

if err := conn.ReadJSON(&msg); err != nil {
    return
}

msgType, _ := msg["type"].(string)
switch msgType {
case "transcript":
    text, _ := msg["text"].(string)
    // ...
}
```

### ✅ Preferred — Deserialize into Concrete Struct

```go
type TranscriptMessage struct {
    Type       string `json:"type"`
    Text       string `json:"text"`
    IsFinal    bool   `json:"is_final"`
}

var msg TranscriptMessage

if err := conn.ReadJSON(&msg); err != nil {
    return
}
// No casts needed — all fields are typed
```

---

## 7. `sql.Scanner` Interface

### ✅ Required — Annotated `.(type)` Switch

```go
// EXEMPTED: sql.Scanner stdlib interface — src is `any` by contract (§7.2)
func (id *ProjectId) Scan(src any) error {
    switch v := src.(type) {
    case string:
        *id = ProjectId(v)
    case []byte:
        *id = ProjectId(string(v))
    default:
        return fmt.Errorf("unsupported type for ProjectId: %T", src)
    }

    return nil
}
```

---

## 8. JSON-LD / Third-Party Schema

### ✅ Required — Annotated Dynamic Schema Parsing

```go
// EXEMPTED: external JSON-LD schema — dynamic structure from third-party HTML (§7.2)
if sameAs, ok := schema["sameAs"]; ok {
    switch v := sameAs.(type) {
    case string:
        urls = append(urls, v)
    case []interface{}:

        for _, item := range v {
            if str, ok := item.(string); ok {
                urls = append(urls, str)
            }
        }
    }
}
```

---

## 9. Centralized Safe-Cast Utility (`CastOrFail[T]`)

All casting/conversion operations **must** go through a centralized utility function that returns `*apperror.AppError` on failure. This ensures every failed cast is traceable, loggable, and propagatable — never silently swallowed.

**Error Code Mapping** (see Error Code Registry — GEN-600 <!-- external: 02-spec/03-error-manage/03-error-code-registry/01-registry.md -->):

| Go Constant | Registry Code | Name | Emitted By |
|-------------|---------------|------|------------|
| `ECast001` | `GEN-600-01` | `CAST_TYPE_ASSERTION_FAILED` | `typecast.CastOrFail[T]()` |
| `ECast002` | `GEN-600-02` | `CAST_SLICE_ELEMENT_FAILED` | `typecast.CastSliceOrFail[T]()` |

### 9.1 Canonical Implementation

```go
// In pkg/typecast/cast.go

package typecast

import (
    "fmt"
    "runtime"

    "myapp/pkg/apperror"
)

// CastOrFail attempts a type assertion and returns an AppError on failure.
// The extra stack skip ensures the error points to the CALLER's call site,
// not to this utility function itself.
func CastOrFail[T any](value any) apperror.Result[T] {
    result, ok := value.(T)

    if !ok {
        var zero T

        return apperror.Fail[T](
            apperror.New(
                "ECast001",  // Registry: GEN-600-01 (CAST_TYPE_ASSERTION_FAILED)
                fmt.Sprintf(
                    "type assertion failed: expected %T, got %T",
                    zero,
                    value,
                ),
            ).WithSkip(1),  // Skip this frame — error points to caller
        )
    }

    return apperror.Ok(result)
}

// CastSliceOrFail casts []interface{} elements to []T.
func CastSliceOrFail[T any](slice []interface{}) apperror.Result[[]T] {
    out := make([]T, 0, len(slice))

    for i, item := range slice {
        result := CastOrFail[T](item)

        if result.HasError() {
            return apperror.Fail[[]T](
                result.AppError().
                    WithContext("SliceIndex", fmt.Sprintf("%d", i)).
                    WithSkip(1),
            )
        }

        out = append(out, result.Value())
    }

    return apperror.Ok(out)
}
```

### 9.2 Usage in Business Logic

```go
// ✅ CORRECT: Safe cast via utility — error is never swallowed
result := typecast.CastOrFail[[]SearchResult](resp.Results)

if result.HasError() {
    return apperror.Fail[MyOutput](result.AppError())
}

results := result.Value()
```

```go
// ✅ CORRECT: Slice casting with full error context
sliceResult := typecast.CastSliceOrFail[SearchResult](rawSlice)

if sliceResult.HasError() {
    return apperror.Fail[MyOutput](sliceResult.AppError())
}
```

### 9.3 Usage in Tests

```go
// ✅ CORRECT: Tests also use CastOrFail — never bare assertions
result := typecast.CastOrFail[[]interface{}](resp.Results)

if result.HasError() {
    t.Fatalf("cast failed: %s", result.AppError().Message)
}

results := result.Value()

if len(results) != 1 {
    t.Errorf("expected 1 result, got %d", len(results))
}
```

### 9.4 Stack Skip Requirement

When writing wrapper functions that call `CastOrFail` internally, each wrapper layer **must** add `.WithSkip(1)` so the reported error location points to the **original caller**, not the intermediate wrapper.

```go
// Wrapper adds +1 skip
func CastResponseField[T any](resp *Response, field string) apperror.Result[T] {
    raw, exists := resp.Fields[field]

    if !exists {
        return apperror.Fail[T](
            apperror.New("ECast002", fmt.Sprintf("field %q not found", field)).  // Registry: GEN-600-02 (CAST_SLICE_ELEMENT_FAILED)
                WithSkip(1),
        )
    }

    result := CastOrFail[T](raw)

    if result.HasError() {
        return apperror.Fail[T](
            result.AppError().
                WithContext("Field", field).
                WithSkip(1),  // +1 for this wrapper layer
        )
    }

    return result
}
```

---

## 10. Rule: Never Swallow Cast Errors

**Cast/conversion failures must NEVER be silently discarded.** Every cast operation must produce an `*apperror.AppError` on failure that is either:

1. **Returned** — propagated to the caller via `apperror.Result[T]`
2. **Logged** — recorded via the structured logger before any fallback

### ❌ Prohibited — Swallowed Cast Error

```go
// FORBIDDEN: comma-ok with blank identifier discards the failure signal
results, _ := resp.Results.([]interface{})

// FORBIDDEN: unchecked assertion — panics at runtime
results := resp.Results.([]interface{})
```

### ✅ Required — All Casts Through Utility

```go
result := typecast.CastOrFail[[]interface{}](resp.Results)

if result.HasError() {
    return apperror.Fail[MyOutput](result.AppError())
}

results := result.Value()
```

### Rule Summary

| Pattern | Verdict |
|---------|---------|
| `x.(T)` bare assertion | ❌ Panic risk — prohibited |
| `x, _ := val.(T)` blank discard | ❌ Swallowed error — prohibited |
| `x, ok := val.(T); if !ok { ... }` in business logic | ❌ §7.2 violation — use `CastOrFail` |
| `x, ok := val.(T)` inside `EXEMPTED` accessor | ✅ Allowed at boundary only |
| `typecast.CastOrFail[T](val)` | ✅ Canonical pattern |

---

## 11. Decision Matrix

| Scenario | Action |
|----------|--------|
| Cast in business logic | **Use `typecast.CastOrFail[T]`** |
| Cast in accessor/wrapper function | **Add `// EXEMPTED:` annotation + `.WithSkip(1)`** |
| Cast on external API response | **Prefer concrete struct; annotate if dynamic** |
| Cast in `sql.Scanner` / stdlib interface | **Annotate as stdlib boundary** |
| Cast on `error` type | **Use `errors.As()`** |
| `.(type)` switch on JSON-decoded data | **Annotate as decoder boundary** |
| Cast in tests | **Use `typecast.CastOrFail[T]` + `t.Fatalf` on error** |

---

## 12. Verification

Search for unannotated violations:

```bash

# Find all type assertions missing EXEMPTED annotation

grep -rn '\.\(string\)\|\.\(float64\)\|\.\(int\)\|\.\(\[\]' spec/ --include="*.md" \
  | grep -v 'EXEMPTED\|FORBIDDEN\|WRONG\|❌\|PROHIBITED'

# Find all .(type) switches missing annotation

grep -rn '\.\(type\)' spec/ --include="*.md" \
  | grep -v 'EXEMPTED\|FORBIDDEN\|WRONG\|❌'

# Find raw ctx.Value casts outside accessor functions

grep -rn 'ctx\.Value\|context\.Value' spec/ --include="*.md" \
  | grep '\.\(string\)' \
  | grep -v 'EXEMPTED\|ctxutil\|accessor'
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| §7.2 Type Assertion Rule | `02-spec/02-coding-guidelines/01-cross-language/13-strict-typing.md` |
| Error Code Registry (GEN-600) | `02-spec/03-error-manage/03-error-code-registry/01-registry.md` |
| Strong Typing Mandate | `.lovable/memories/standards/strong-typing-mandate.md` |
| Error Handling Standards | `02-spec/03-error-manage/01-error-resolution/10-apperror-package/01-apperror-reference.md` |
| Control Flow Rules | `.lovable/memories/architecture/coding-standards/control-flow.md` |
| Contradiction Checks | `02-spec/02-coding-guidelines/01-cross-language/05-cross-spec-contradiction-checks.md` |
| Context Accessors | `02-spec/02-spec-management-software/13-shared-packages/04-pkg-logging.md` |

---

*Updated 2026-02-28 — v2.0.0. Added centralized `CastOrFail[T]` utility pattern (§9), "never swallow cast errors" rule (§10), stack skip requirement (§9.4), and test casting patterns (§9.3).*
