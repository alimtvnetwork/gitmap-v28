# Master Coding Guidelines — Type safety, single return value, no casting

> **Parent:** [Master Coding Guidelines](./01-index.md)
> **Version:** 2.1.0
> **Updated:** 2026-03-31

---

## 7. Type Safety

### PHP

- Native type declarations on all parameters, return values, properties
- Remove redundant PHPDoc when native types are present
- Max 3 parameters per function

### Go

- Zero `any`/`interface{}`/`map[string]any` in business logic
- `json.RawMessage` only at architectural boundaries
- Concrete domain models for all handler decoding

### Common Mistakes

```go
// ❌ MISTAKE: Type erasure
func ProcessData(data interface{}) interface{} { ... }

// ✅ CORRECT: Concrete types
func ProcessData(data PluginDetails) apperror.Result[PluginSummary] { ... }
```

---


---

## 7.1 Single Return Value Rule (Go)

**Every Go function must return exactly one value.** Multiple return values like `(bool, *Result, bool, *AppError)` are **prohibited**. Use a typed result struct instead.

### Why

Multiple return values:

- Create ambiguous call sites — callers must remember positional meaning
- Break serialization — you cannot JSON-encode a multi-return tuple
- Defeat the Result pattern — error handling becomes scattered across return positions

### Rule

| ❌ Prohibited | ✅ Required |
|--------------|------------|
| `func F() (T, error)` | `func F() apperror.Result[T]` |
| `func F() (T, bool)` | `func F() apperror.Result[T]` (use `.IsEmpty()` for "not found") |
| `func F() (bool, *R, bool, *AppError)` | `func F() apperror.Result[OnboardResult]` with a typed payload struct |
| `func F() (int, string, error)` | `func F() apperror.Result[MyOutput]` with `MyOutput` struct |

**One exception:** Go's idiomatic comma-ok pattern for map lookups (`v, ok := m[k]`) and type assertions within the `apperror` package internals are exempt. All **custom function signatures** follow this rule.

### Example: Before and After

```go
// ❌ PROHIBITED: 4 return values — positional, ambiguous, unserializable
func (s *Service) OnboardUpload(
    ctx context.Context,
    req OnboardRequest,
) (bool, *OnboardUploadResult, bool, *apperror.AppError) {
    // ...
    return true, result, false, nil
}

// Caller — positional guessing, fragile
isNew, result, hasConflict, appErr := svc.OnboardUpload(ctx, req)
if appErr != nil { ... }
if hasConflict { ... }
```

```go
// ✅ REQUIRED: Single typed result
type OnboardOutcome struct {
    IsNew       bool
    HasConflict bool
    Upload      OnboardUploadResult
}

func (s *Service) OnboardUpload(
    ctx context.Context,
    req OnboardRequest,
) apperror.Result[OnboardOutcome] {
    // ...
    return apperror.Ok(OnboardOutcome{
        IsNew:       true,
        HasConflict: false,
        Upload:      uploadResult,
    })
}

// Caller — clear, self-documenting
result := svc.OnboardUpload(ctx, req)

if result.HasError() { return result }
outcome := result.Value()

if outcome.HasConflict { ... }
```

### Payload Struct Naming

Name the return struct after the operation + `Outcome` or `Output`:

| Operation | Payload Struct |
|-----------|---------------|
| `OnboardUpload` | `OnboardOutcome` |
| `CheckSync` | `SyncCheckOutput` |
| `PublishPlugin` | `PublishOutcome` |
| `ValidateConfig` | `ValidationOutput` |

---


---

## 7.2 No Type Assertions / Casting (Go)

**Type assertions (`.(*Type)`, `.(string)`, `.(float64)`) are prohibited** in business logic and service code. Their presence indicates missing concrete types upstream.

**Canonical utility:** All unavoidable casts must go through `typecast.CastOrFail[T]()` from `pkg/typecast/`, which returns `apperror.Result[T]` on failure. See [Casting Elimination Patterns](../04-casting-elimination-patterns.md) for the full specification, including `CastSliceOrFail[T]`, stack-skip requirements, and test patterns.

### Why

Type assertions:

- Panic at runtime if the unchecked form is used (`x.(string)` without comma-ok)
- Signal that the data model is untyped (`interface{}`, `map[string]any`)
- Defeat the generics-first and strict-typing policies

### Rule

| ❌ Prohibited | ✅ Required |
|--------------|------------|
| `msg["text"].(string)` | Deserialize into a concrete struct |
| `cached.(float64)` | Use a typed cache: `TypedCache[float64]` |
| `payload.(*AIRequest)` | Use generics: `Envelope[AIRequest]` |
| `int(normMap["max"].(float64))` | Struct field: `NormConfig.Max int` |
| `x, _ := val.(T)` (swallowed error) | `typecast.CastOrFail[T](val)` — never discard |
| `x.(T)` bare assertion | `typecast.CastOrFail[T](x)` — returns `AppError` |

**Cast errors must NEVER be swallowed.** The blank-identifier pattern (`_, _ :=`) and bare assertions are both prohibited. All casts produce an `*apperror.AppError` on failure that must be returned or logged.

**Exemptions:**

- Go standard library interfaces: `net.Listener.Addr().(*net.TCPAddr)` — framework boundary
- `apperror` package internals — error unwrapping
- Test helpers using `require` / `assert` (test-only code) — but prefer `typecast.CastOrFail[T]` + `t.Fatalf`
- Annotate any exemption with `// EXEMPTED: <reason>` and use `.WithSkip(1)` in wrapper functions

### Example: Before and After

```go
// ❌ PROHIBITED: Type assertion on map values — runtime panic risk
func (a *Analyzer) loadWeights(cfg map[string]any) {
    a.weights = Weights{
        Stars:    cfg["github_stars"].(float64),
        Jobs:     cfg["job_postings"].(float64),
        Packages: cfg["package_downloads"].(float64),
    }
}

// ✅ REQUIRED: Concrete config struct — compile-time safe
type WeightsConfig struct {
    Stars    float64
    Jobs     float64
    Packages float64
}

func (a *Analyzer) loadWeights(cfg WeightsConfig) {
    a.weights = Weights{
        Stars:    cfg.Stars,
        Jobs:     cfg.Jobs,
        Packages: cfg.Packages,
    }
}
```

```go
// ❌ PROHIBITED: Swallowed cast error
results, _ := resp.Results.([]interface{})

// ✅ REQUIRED: Safe cast via utility
result := typecast.CastOrFail[[]interface{}](resp.Results)
if result.HasError() {
    return apperror.Fail[MyOutput](result.AppError())
}
results := result.Value()
```

```go
// ❌ PROHIBITED: Casting from cache
if cached, ok := s.cache.Load(key); ok {
    return cached.(float64), nil
}

// ✅ REQUIRED: Typed cache wrapper
type TypedCache[T any] struct { sync.Map }

func (c *TypedCache[T]) Load(key string) (T, bool) {
    v, ok := c.Map.Load(key)
    if !ok {
        var zero T

        return zero, false
    }

    return v.(T), true  // EXEMPTED: Generic wrapper internal cast
}

// Caller — no assertion needed
if val, ok := s.cache.Load(key); ok {
    return val, nil
}
```

---
