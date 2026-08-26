# Golang Coding Standards — Type safety, error handling, Result types

> **Parent:** [Golang Coding Standards](./00-overview.md)  
> **Version:** 3.7.0  
> **Updated:** 2026-03-31

---

## Type Safety — No `interface{}` or `any`

### Rule: Never use `interface{}` or `any` in exported APIs

```go
// ❌ FORBIDDEN
func ProcessData(data interface{}) interface{} { ... }
func FetchResults() (any, error) { ... }

// ✅ REQUIRED: Use concrete types or generics
func ProcessData(data PluginDetails) apperror.Result[PluginSummary] { ... }
func FetchResults[T any]() apperror.Result[T] { ... }
```

### Acceptable `any` Usage

1. **SQL query arguments** — `args ...any` in `dbutil` (framework boundary)
2. **Logger variadic parameters** — `map[string]any` for structured log fields (internal only)
3. **Third-party library interfaces** — When a library requires `interface{}`

---


---

## Error Handling — `apperror` Package

### Rule 1: Every error carries a mandatory stack trace

All errors created via `apperror.New()` or `apperror.Wrap()` automatically capture a full `StackTrace` at creation — no opt-in needed.

```go
// ❌ FORBIDDEN: loses stack trace
return fmt.Errorf("failed to upload: %w", err)

// ✅ REQUIRED: full stack trace captured automatically
return apperror.Wrap(
    err,
    "E5001",
    "failed to upload plugin",
)
```

### Rule 2: AppError-Only — No Raw `error` in Application Code

**All application code MUST use `*apperror.AppError` instead of Go's raw `error` type.** The standard `error` interface is only acceptable at framework boundaries — the moment you receive a raw `error` from a framework, stdlib call, or third-party library, you **must immediately wrap it** into an `AppError`.

```go
// ❌ FORBIDDEN: Returning raw error from application code
func (s *PluginService) Upload(ctx context.Context, req UploadRequest) error {
    file, err := os.Open(req.Path)
    if err != nil {
        return err  // ❌ Raw error escaping application boundary
    }
    // ...
}

// ✅ REQUIRED: Use fileutil helpers that already return apperror.Result
func (s *PluginService) Upload(ctx context.Context, req UploadRequest) apperror.Result[UploadResult] {
    result := fileutil.Open(req.Path) // returns apperror.Result[*os.File]
    if result.IsErr() {
        return apperror.Fail[UploadResult](result.AppError())
    }
    file := result.Value
    // ...use file...
}

// fileutil package — thin wrappers returning apperror.Result instead of raw error
// See 08-pathutil-fileutil-spec.md for full signatures, error codes, and usage examples
// Key functions: Open, ReadAll, WriteFile, MkdirAll, EnsureDir, Stat
// Boolean guards (pathutil): IsFileExists, IsFileMissing, IsDir, IsDirMissing
```

#### Where Raw `error` Is Acceptable

| Context | Why | Example |
|---------|-----|---------|
| Receiving from stdlib/framework | You don't control the return type | `os.Open()`, `json.Unmarshal()`, `sql.Query()` |
| Implementing stdlib interfaces | Go requires `error` return | `io.Reader`, `http.Handler`, `encoding.Marshaler` |
| `main()` or top-level bootstrap | Before apperror is initialized | `log.Fatal(err)` |

#### Where Raw `error` Is FORBIDDEN

| Context | Required Instead |
|---------|-----------------|
| Service method returns | `apperror.Result[T]` / `*apperror.AppError` |
| Repository/store returns | `apperror.Result[T]` / `apperror.ResultSlice[T]` |
| Error propagation between services | `apperror.Fail[T](src.AppError())` |
| Error creation for business logic | `apperror.New("E1001", "message")` |
| Wrapping framework errors | `apperror.Wrap(err, "E1002", "context")` |

#### The Wrap-Immediately Pattern

```go
// Pattern: Wrap at first contact, propagate as AppError thereafter
func (s *SyncService) Pull(ctx context.Context, pluginID int64) apperror.Result[PullResult] {
    // Framework boundary — raw error received and IMMEDIATELY wrapped
    output, err := exec.CommandContext(ctx, "git", "pull").Output()

    if err != nil {
        return apperror.FailWrap[PullResult](
            err,
            "E7001",
            "git pull failed",
        )
    }

    // Application boundary — already AppError, propagate directly
    plugin := s.pluginService.GetById(ctx, pluginID)

    if plugin.HasError() {
        return apperror.Fail[PullResult](plugin.AppError())
    }

    return apperror.Ok(PullResult{Output: string(output)})
}
```

### StackTrace Type

```go
// Captured automatically — structured frames, not raw strings
type StackFrame struct {
    Function string
    File     string
    Line     int
}

type StackTrace []StackFrame

// Display methods
trace.String()      // full formatted multi-line trace
trace.CallerLine()  // "file.go:42" — compact single line
trace.IsEmpty()     // no frames captured
trace.Depth()       // number of frames
```

### AppError Display Methods

```go
err.Error()       // "[E5001] upload failed" — implements error interface
err.FullString()  // code + message + diagnostics + stack + cause chain
err.ToClipboard() // markdown-formatted error report for AI paste
```

### Context Enrichment — Typed Diagnostic Setters

```go
// ✅ Enriched error with diagnostic context
return apperror.Wrap(
    err,
    "E5002",
    "remote site request failed",
).
    WithUrl(requestUrl).
    WithSlug(pluginSlug).
    WithStatusCode(resp.StatusCode).
    WithSiteId(siteId)
```

### Error Code Convention

| Range | Category |
|-------|----------|
| E1xxx | Configuration errors |
| E2xxx | Database errors |
| E3xxx | WordPress API errors |
| E4xxx | File system errors |
| E5xxx | Sync errors |
| E6xxx | Backup errors |
| E7xxx | Git errors |
| E8xxx | Build errors |
| E9xxx | General errors |
| E10xxx | E2E test errors |
| E11xxx | Publish errors |
| E12xxx | Version errors |
| E13xxx | Session errors |
| E14xxx | Crypto errors |

---


---

## Generic Result Types — `apperror` Package

Three generic result types for all service returns. Replaces raw `(T, error)` tuples.

### `Result[T]` — Single Value

For operations that return one item or nothing.

```go
// Construction
result := apperror.Ok(plugin)             // success
result := apperror.Fail[Plugin](appErr)   // from AppError
result := apperror.FailWrap[Plugin](
    err,
    "E5001",
    "load failed",
)
result := apperror.FailNew[Plugin](
    "E4004",
    "not found",
)

// Query methods
result.HasError()    // true if operation failed
result.IsSafe()      // true if value exists AND no error
result.IsDefined()   // true if value was set
result.IsEmpty()     // true if no value was set

// Access methods
result.Value()             // returns T; panics if HasError
result.ValueOr(fallback)   // returns T or fallback if empty
result.AppError()          // returns *AppError or nil (named AppError to avoid confusion with Go's error)
result.Unwrap()            // bridges to (T, error) pattern
```

### `ResultSlice[T]` — Collection (Array)

For operations that return lists of items.

```go
// Construction
set := apperror.OkSlice(plugins)
set := apperror.FailSlice[Plugin](appErr)
set := apperror.FailSliceWrap[Plugin](err, "E5011", "query failed")

// Query methods
set.HasError()     // true if operation failed
set.IsSafe()       // true if no error (items may be empty)
set.HasItems()     // true if at least one item
set.IsEmpty()      // true if zero items
set.Count()        // number of items

// Access methods
set.Items()        // returns []T (nil if error)
set.First()        // Result[T] for first item
set.Last()         // Result[T] for last item
set.GetAt(index)   // Result[T] at index; empty if out of bounds
set.AppError()     // returns *AppError or nil

// Mutation methods
set.Append(items...)  // adds items; no-op if in error state
```

### `ResultMap[K, V]` — Associative Map

For operations that return key-value data.

```go
// Construction
m := apperror.OkMap(pluginsBySlug)
m := apperror.FailMap[string, Plugin](appErr)
m := apperror.FailMapWrap[string, Plugin](err, "E5012", "index failed")

// Query methods
m.HasError()     // true if operation failed
m.IsSafe()       // true if no error (map may be empty)
m.HasItems()     // true if at least one entry
m.IsEmpty()      // true if zero entries
m.Count()        // number of entries
m.Has(key)       // true if key exists

// Access methods
m.Items()        // returns map[K]V (nil if error)
m.Get(key)       // Result[V] for key; empty if not found
m.Keys()         // returns []K
m.Values()       // returns []V
m.AppError()     // returns *AppError or nil

// Mutation methods
m.Set(key, value)   // adds/updates; no-op if error state
m.Remove(key)       // deletes key; no-op if error state
```

### Service Usage Pattern

```go
// ✅ Same-type propagation — use bridge method (no unwrap+rewrap)
func (s *SiteService) ListAll(ctx context.Context) apperror.ResultSlice[Site] {
    set := dbutil.QueryMany[Site](ctx, s.db, query, scanSite)

    if set.HasError() {
        return set.ToAppResultSlice()
    }

    return apperror.OkSlice(set.Items())
}

// ✅ Single-row with post-processing — direct propagation via AppError()
func (s *PluginService) GetById(ctx context.Context, id int64) apperror.Result[Plugin] {
    dbResult := dbutil.QueryOne[Plugin](ctx, s.db, query, scanPlugin, id)

    if dbResult.HasError() {
        return apperror.Fail[Plugin](dbResult.AppError())
    }

    if dbResult.IsEmpty() {
        return apperror.FailNew[Plugin](
            ErrNotFound,
            "plugin not found",
        )
    }

    return apperror.Ok(dbResult.Value())
}

// ✅ Cross-type propagation — different T, Fail[NewT] is correct
func (s *GitService) Pull(ctx context.Context, pluginID int64) apperror.Result[PullResult] {
    pResult := s.pluginService.GetById(ctx, pluginID)

    if pResult.HasError() {
        return apperror.Fail[PullResult](pResult.AppError())
    }
    // ...
}

// ✅ Handler consuming Result[T]
func (h *Handler) GetPlugin(w http.ResponseWriter, r *http.Request) {
    result := h.plugins.GetById(r.Context(), pluginId)

    if result.HasError() {
        writeError(w, result.AppError())

        return
    }

    writeJSON(w, result.Value())
}
```

### Error Propagation Rules

| Scenario | Pattern | Example |
|----------|---------|---------|
| Same T, dbutil→apperror slice | Bridge method | `set.ToAppResultSlice()` |
| Same T, dbutil→apperror single | Bridge method | `result.ToAppResult()` |
| Different T (cross-type) | `Fail[NewT](src.AppError())` | `apperror.Fail[BuildResult](pluginResult.AppError())` |
| Same wrapper, same T | Direct return | `return existingResult` |

> **Anti-pattern:** Never unwrap an error just to re-wrap it into the same type parameter:
> ```go
> // ❌ FORBIDDEN: redundant unwrap+rewrap (same T)
> return apperror.FailSlice[Plugin](set.AppError())
>
> // ✅ REQUIRED: use bridge method
> return set.ToAppResultSlice()
> ```

---
