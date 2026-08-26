# AppError Package Reference — Usage examples, service adapter unwrap pattern

> **Parent:** [AppError Package Reference](./00-overview.md)  
> **Version:** 1.3.0  
> **Updated:** 2026-03-31

---

## 9. Usage Examples

### Service Method Returning Result[T]

```go
func (s *PluginService) GetById(context stdctx.Context, id int64) apperror.Result[Plugin] {
    plugin, err := s.repo.FindById(context, id)

    if err != nil {
        return apperror.FailWrap[Plugin](err, apperror.ErrDatabaseQuery, "get plugin by id").
            WithValue("PluginId", fmt.Sprintf("%d", id))
    }

    if plugin == nil {
        return apperror.FailNew[Plugin](apperror.ErrNotFound, "plugin not found")
    }

    return apperror.Ok(*plugin)
}
```

### Handler Consuming Result[T]

```go
func (h *Handler) GetPlugin(w http.ResponseWriter, r *http.Request) {
    result := h.plugins.GetById(r.Context(), pluginId)

    if result.HasError() {
        writeError(w, result.AppError())

        return
    }

    writeJson(w, result.Value())
}
```

### Error with Values

```go
return apperror.Wrap(err, apperror.ErrFSRead, "failed to read config").
    WithValue("path", configPath).
    WithValue("format", "yaml")
```

### Using `apperrtype` Enums (Preferred)

Three escalating levels of type safety — Level 3 is the target for all new code:

```go
// ❌ Level 1 — raw strings (flagged by CODE-RED-008 lint rule)
apperror.New("E2010", "site not found")

// ✅ Level 2 — enum code, manual message
apperror.New(apperrtype.SiteNotFound.Code(), "site not found")

// ✅✅ Level 3 — enum with built-in message (best)
apperror.NewType(apperrtype.SiteNotFound)
```

> **Side note:** `FailBool` is a convenience constructor for `Result[bool]`. It creates a failed
> `Result[bool]` from an `*AppError` — saving you from writing `apperror.Fail[bool](err)` everywhere.
> The same pattern applies to `FailSettings`, `FailString`, etc. — each is a type alias shortcut.

### With Convenience Constructors (Type Aliases)

```go
// ✅ Best practice — enum + type alias + convenience constructor
return apperror.FailBool(apperror.NewType(apperrtype.SiteNotFound))
return apperror.FailSettings(apperror.NewType(apperrtype.ConfigKeyMissing))

// Equivalent long-form (what FailBool replaces):
return apperror.Fail[bool](apperror.NewType(apperrtype.SiteNotFound))
```

### Error with Diagnostics + Values + ErrorType

```go
return apperror.WrapType(err, apperrtype.WPConnectionFailed).
    WithValue("url", siteURL).
    WithValue("plugin", pluginSlug).
    WithStatusCode(resp.StatusCode).
    WithMethod("GET").
    WithEndpoint("/wp-json/wp/v2/plugins")
```

---

## 10. Service Adapter Unwrap Pattern

### 10.1 Architectural Boundary

Services return `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` to preserve rich error context and type safety within the domain layer. HTTP handlers consume **adapter interfaces** that expose standard `(T, error)` tuples. A dedicated **Service Adapter** sits between them, acting as the single unwrap boundary.

```
┌─────────────┐    Result[T]    ┌──────────────────┐   (T, error)   ┌──────────┐
│   Service    │ ─────────────► │  ServiceAdapter   │ ─────────────► │  Handler │
│  (domain)    │                │  (unwrap layer)   │                │  (HTTP)  │
└─────────────┘                └──────────────────┘                └──────────┘
```

**Rules:**
- Services **never** return raw `(T, error)` for data-fetching operations — use `Result[T]` or `ResultSlice[T]`
- Void operations (`Delete`, `MarkSynced`, etc.) may return plain `error`
- Adapters are the **only** place that calls `.Value()`, `.Items()`, or `.AppError()` to convert back to tuples
- Handlers and other transport-layer code **never** import `apperror.Result` types directly

### 10.2 Adapter Implementation

Each service gets a dedicated adapter file (e.g., `adapter_plugin.go`, `adapter_site.go`, `adapter_sync.go`) in the `handlers` package:

```go
// SiteServiceAdapter wraps *site.Service to implement SiteServiceInterface
type SiteServiceAdapter struct {
    *site.Service
}

// Result[T] → (*T, error) unwrap for single-value returns
func (a *SiteServiceAdapter) GetById(context stdctx.Context, id int64) (*models.Site, error) {
    result := a.Service.GetById(context, id)  // returns apperror.Result[models.Site]

    if result.HasError() {
        return nil, result.AppError()
    }

    v := result.Value()

    return &v, nil
}

// ResultSlice[T] → ([]T, error) unwrap for collection returns
func (a *SiteServiceAdapter) List(context stdctx.Context) ([]models.Site, error) {
    result := a.Service.List(context)  // returns apperror.ResultSlice[models.Site]

    if result.HasError() {
        return nil, result.AppError()
    }

    return result.Items(), nil
}
```

### 10.3 Compile-Time Verification

All adapters include compile-time interface checks in `adapters.go`:

```go
var _ SiteServiceInterface = (*SiteServiceAdapter)(nil)
var _ PluginServiceInterface = (*PluginServiceAdapter)(nil)
var _ SyncServiceInterface = (*SyncServiceAdapter)(nil)
```

### 10.4 Cross-Service Consumption

When **Service A** holds a direct reference to **Service B** (not through the adapter), Service A must consume Result types directly using `.HasError()` / `.Value()` / `.IsSafe()`:

```go
// sync service calls plugin service directly (not through adapter)
plugResult := s.pluginService.GetById(ctx, pluginId)

if plugResult.HasError() {
    return apperror.FailWrap[PushSyncResult](plugResult.AppError(), apperror.ErrDatabaseQuery, "failed to get plugin")
}

plug := plugResult.Value()
```

**Cross-service audit checklist** — when migrating a service to Result types, verify:
1. All cross-service callers that hold a direct `*service.Service` reference
2. All `main.go` initialization code that calls service methods
3. All adapter methods are updated to unwrap the new return types

### 10.5 Zero Raw Error Rule

**No service method may return a bare `error` from the standard library.** Every error returned from a service function must be an `*apperror.AppError` (created via `apperror.New`, `apperror.Wrap`, or contained within a `Result[T]`). This guarantees every error carries a stack trace for diagnostics.

**Forbidden patterns:**
```go
// ❌ NEVER — no stack trace captured
return err
return fmt.Errorf("something failed: %w", err)
return errors.New("something failed")
```

**Required patterns:**
```go
// ✅ Wraps with stack trace + error code
return apperror.Wrap(
    err, apperror.ErrDatabaseExec, "failed to update config",
)

// ✅ New error with stack trace + error code
return apperror.New(
    apperror.ErrNotFound, "entry not found",
)

// ✅✅ BEST: error type enum — code + message from enum, zero duplication
return apperror.NewType(apperrtype.EntryNotFound)
```

**Exemptions:**
- `filepath.Walk` callbacks (framework requires `error` interface)
- E2E test harness (`e2e/` package) — test assertion errors, not production
- Enum `UnmarshalJSON` / `MarshalJSON` methods (`internal/enums/*/variant.go`) — circular import risk with `apperror` package; these are standard library interface implementations

### 10.6 No Raw `error` in Struct Fields (Invariant I-2)

**Struct fields that represent errors must use `*apperror.AppError`, never Go's `error` interface.** The `error` interface is not serializable — `json.Marshal` cannot introspect its internals, producing only `{}` or requiring custom marshaling at every usage site. `*AppError` is fully serializable by design (see §11), carrying code, message, stack trace, diagnostics, and cause through any transport boundary.

```go
// ❌ FORBIDDEN — error interface is not serializable
type JobResult struct {
    Output string
    Err    error     // json.Marshal produces {} — all diagnostic context lost
}

// ❌ FORBIDDEN — Inner/Cause as raw error in custom structs
type OperationLog struct {
    Action    string
    Inner     error    // not serializable, not queryable
}

// ✅ REQUIRED — *AppError carries full context and serializes cleanly
type JobResult struct {
    Output   string
    AppError *apperror.AppError `json:",omitempty"`  // fully serializable
}

// ✅ REQUIRED — structured error in operation logs
type OperationLog struct {
    Action   string
    AppError *apperror.AppError `json:",omitempty"`
}
```

**Why this matters:**
- `*AppError` serializes to a complete JSON object with code, message, stack, values, and diagnostics
- Raw `error` serializes to `{}` or requires per-struct custom `MarshalJSON` (DRY violation)
- Error history DB stores `*AppError` as structured JSON — raw `error` cannot be queried
- Subprocess JSON protocol transmits `*AppError` — raw `error` loses all context

**Exemptions:**
- The `Cause` field on `AppError` itself uses `error` (handled by custom `MarshalJSON/UnmarshalJSON` — see §11.2/§11.3)
- Constructor parameters (`Wrap(cause error, ...)`) accept `error` at the wrapping boundary — this is where raw errors enter the system and get wrapped

**Cross-language equivalent:**
- **PHP:** Struct/class error fields use the framework's `Throwable` type with `stackTrace()` method, never bare `string` or `null`
- **TypeScript:** Error fields use the framework's structured error type with `code`, `message`, and `stack` properties, never bare `Error` or `string`

### 10.6 Migrated Services

| Service | Result Types | Adapter File |
|---------|-------------|--------------|
| Plugin | `List`, `GetById`, `Create`, `Update`, `ScanDirectory`, `GetMappings`, `GetMappingsBySite`, `CreateMapping` | `adapter_plugin.go` |
| Site | `List`, `GetById`, `GetByUrl`, `Create`, `Update` | `adapter_site.go` |
| Sync | `CheckSync`, `CheckAllSites`, `CheckAllPlugins`, `PushSync`, `GetFileChanges` | `adapter_sync.go` |
| Publish | `Publish`, `PublishFiles`, `PreviewPublish`, `GetFileDiff` | `adapter_publish.go` |
| Git | `Pull`, `PullAll`, `Build`, `PullAndBuild`, `GetConfig`, `Status`, `Commit`, `Push` | `adapter_git.go` |
| Watcher | `TriggerScan`, `ScanAfterGitPull`, `ScanAll` | `adapter_sync.go` |
| Backup | `Create`, `List`, `GetById`, `Restore`, `ExportToZip`, `ImportFromZip` | `adapter_publish.go` |
| Session | `GetSession`, `GetSessionLogs`, `GetSessionDiagnostics`, `ListSessions` | `adapter_session.go` |
| ErrorHistory | `Save`, `List`, `GetById`, `GetByErrorId`, `Clear`, `BulkExport`, `GetStats` | `adapter_session.go` |
| SiteHealth | `CheckSite`, `CheckAllSites`, `GetHistory`, `GetSummaries`, `GetStats`, `ClearHistory` | `adapter_history.go` |
| PublishHistory | `Record`, `List`, `GetById`, `GetStats`, `Clear` | `adapter_history.go` |

---
