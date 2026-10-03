# AppError Package Reference — Usage examples, service adapter unwrap pattern

> **Parent:** [AppError Package Reference](./readme.md)
> **Version:** 1.3.0
> **Updated:** 2026-03-31

---

## 9. Usage Examples (Canonical Reference: `04-code/golang/examples/`)

> Real-world implementations are maintained in [`04-code/golang/examples/database_query.go`](../../../../../04-code/golang/examples/database_query.go) and [`04-code/golang/examples/workflow_service.go`](../../../../../04-code/golang/examples/workflow_service.go).

### Service Method Returning Concrete Result Type (`types.go`)

```go
// -----------------------------------------------------------------------------
// Step 1: Declare Concrete Types in `types.go` (Mandatory Rule)
// -----------------------------------------------------------------------------
// In types.go:
// type PluginWrapResult = result.Wrap[Plugin]
// -----------------------------------------------------------------------------

func (s *PluginService) GetById(ctx context.Context, id int64) PluginWrapResult {
    if id <= 0 {
        fault := appfault.New(errtype.Validation, "plugin id must be positive").
            WithOp("PluginService.GetById").
            WithVar("id", id)

        return result.WrapFailure[Plugin](fault)
    }

    plugin, err := s.repo.FindById(ctx, id)

    if err != nil {
        fault := appfault.Wrap(errtype.Database, err, "failed to get plugin by id").
            WithOp("PluginService.GetById").
            WithVar("id", id)

        return result.WrapFailure[Plugin](fault)
    }

    return result.WrapSuccess(*plugin)
}
```

### Handler Consuming `result.Wrap[T]`

```go
func (h *Handler) GetPlugin(w http.ResponseWriter, r *http.Request) {
    res := h.plugins.GetById(r.Context(), pluginId)

    if res.IsFailed() {
        writeErrorEnvelope(w, res.Fault())

        return
    }

    writeJsonEnvelope(w, res.Value())
}
```

### Direct Error Typing (No Type Branching or Nested Ifs)

```go
// Direct error typing: select the error type reflecting this layer (errtype.IO).
// Never branch on error types or nest ifs!
fault := appfault.WrapFile(errtype.IO, err, relativePath, "failed to read config").
    WithOp("config.Load").
    WithVar("format", "yaml")

    return result.WrapFailure[Config](fault)
}
```

### Error Construction with Variation Enum & Fluent Builders

```go
// Direct constructor from Variation enum
fault := appfault.New(errtype.Validation, "site not found").
    WithOp("site.Find").
    WithVar("siteId", siteId)

return result.WrapFailure[Site](fault)
```

### Error with Diagnostics, Variables, and Fluent Properties

```go
fault := appfault.Wrap(errtype.IO, err, "failed during health check").
    WithOp("site.CheckHealth").
    WithVar("url", siteURL).
    WithVar("plugin", pluginSlug).
    WithVar("statusCode", resp.StatusCode)

return result.WrapFailure[Site](fault)
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
