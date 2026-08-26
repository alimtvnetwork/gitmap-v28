# Consolidated: Error Management — Complete Reference

**Version:** 3.3.0  
**Updated:** 2026-04-22  
**Source Module:** [`spec/03-error-manage/`](../03-error-manage/00-overview.md)

---

## Purpose

This is the **standalone consolidated reference** for all error management. An AI reading only this file must be able to enforce every rule without consulting source specs.

> 🔴 **#1 PRIORITY** — Error management is the highest priority specification. It must be implemented from the very first line of code.

### Source-Folder Coverage Map

The source module `spec/03-error-manage/` contains 3 top-level subfolders and 13 nested ones. Every one is summarized below. Use the table to jump to the relevant section.

| Source Path | Section | Status |
|-------------|---------|--------|
| `01-error-resolution/` (cheat sheet, cross-ref diagram) | §§1–4, §17 | ✅ Full |
| `01-error-resolution/03-retrospectives/` | §18 Retrospectives | ✅ Reference |
| `01-error-resolution/04-verification-patterns/` | §19 Verification Patterns | ✅ Reference |
| `01-error-resolution/05-debugging-guides/` | §20 Per-Language Debugging | ✅ Full |
| `02-error-architecture/01-error-handling-reference.md` | §§5–10 (envelope, codes, registry) | ✅ Full |
| `02-error-architecture/03-notification-colors.md` | §21 Notification Colors | ✅ Full |
| `02-error-architecture/04-error-modal/` (6 files) | §22 Error Modal System | ✅ Full |
| `02-error-architecture/05-response-envelope/` (6 files + 6 JSON fixtures) | §23 Response Envelope Reference | ✅ Full |
| `02-error-architecture/06-apperror-package/` | §24 apperror Package Reference | ✅ Full |
| `02-error-architecture/07-logging-and-diagnostics/` | §25 Logging & Diagnostics | ✅ Full |
| `03-error-code-registry/` (registry + linter scripts + schemas) | §26 Error Code Registry | ✅ Full |

---

## Three-Tier Error Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Frontend (React/TypeScript)                   │
│  ┌─────────────────┐    ┌──────────────────┐    ┌───────────────┐   │
│  │ API Client       │───▸│ Error Store       │───▸│ Global Error  │   │
│  │ (parseEnvelope)  │    │ (captureError)    │    │ Modal (tabs)  │   │
│  └─────────────────┘    └──────────────────┘    └───────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              ▲ Universal Response Envelope
┌─────────────────────────────────────────────────────────────────────┐
│                        Backend (Go)                                  │
│  ┌─────────────────┐    ┌──────────────────┐    ┌───────────────┐   │
│  │ apperror.Wrap() │───▸│ Session Logger    │───▸│ error.log.txt │   │
│  │ + .WithContext() │    │ (per-request ID)  │    │ (deduped)     │   │
│  └─────────────────┘    └──────────────────┘    └───────────────┘   │
│         │ DelegatedRequestServer Builder                             │
│         │ • Captures: endpoint, method, statusCode                   │
│         │ • Captures: requestBody, response, stackTrace              │
│         │ • Injects into Envelope.Errors block                       │
└─────────────────────────────────────────────────────────────────────┘
                              ▲ REST API (JSON)
┌─────────────────────────────────────────────────────────────────────┐
│              Delegated Server (PHP / Chrome Extension / Other)       │
│  ┌─────────────────┐    ┌──────────────────┐    ┌───────────────┐   │
│  │ safe_execute()  │───▸│ FileLogger        │───▸│ stacktrace.txt│   │
│  │ catch Throwable │    │ (6-frame backtrace)│   │ error.txt     │   │
│  └─────────────────┘    └──────────────────┘    └───────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### Request Chain (3-Hop)

```
React Frontend → Go Backend → Delegated Server (PHP/other)
     │                │                │
     │ GET /api/v1/   │ GET /wp-json/  │
     │ sites/1/data   │ plugin/v1/data │
     │───────────────▸│───────────────▸│
     │                │     HTTP 403   │
     │                │◀──────────────│
     │   HTTP 500     │               │
     │   Envelope     │               │
     │◀──────────────│               │
     │ captureError() │               │
```

| Tier | Layer | Responsibility | Technology |
|------|-------|----------------|------------|
| 1 | Delegated Server (PHP) | Structured error responses with error codes, `TypedQuery` result envelopes | PHP `ResponseKeyType` enum, `DbResult<T>`, `DbResultSet<T>`, `DbExecResult` |
| 2 | Go Backend | Error wrapping, stack traces, typed error codes, structured logging | `apperror` package with `Result[T]`, `Wrap()`, `WithContext()` |
| 3 | Frontend | Error store, Global Error Modal, toast notifications, retry logic | React error boundary, Zustand error store, Sonner toasts |

---

## Universal Response Envelope

**All APIs MUST return this envelope.** No exceptions.

### Success Response

```json
{
  "Status": {
    "IsSuccess": true,
    "Code": 200,
    "Message": "OK"
  },
  "Attributes": {
    "RequestedAt": "/plugin/v1/status",
    "Duration": "45ms",
    "TotalRecords": 150,
    "Limit": 50,
    "Offset": 0
  },
  "Results": [{ "..." }]
}
```

### Error Response

```json
{
  "Status": {
    "IsSuccess": false,
    "Code": 500,
    "Message": "Database connection failed"
  },
  "Error": {
    "ErrorCode": 1001,
    "ErrorType": "DATABASE_ERROR",
    "Detail": "Connection timeout after 5s",
    "StackTrace": "..."
  },
  "Errors": {
    "Backend": ["apperror stack frame 1", "..."],
    "BackendMessage": "open database: file locked",
    "DelegatedRequestServer": {
      "Endpoint": "https://example.com/wp-json/plugin/v1/data",
      "Method": "GET",
      "StatusCode": 403,
      "Response": "{ ... raw JSON ... }",
      "StackTrace": "PHP stack trace...",
      "RequestBody": null
    }
  },
  "MethodsStack": ["Handler.ServeHTTP", "Router.Route", "Middleware.Auth"],
  "SessionId": "sess_abc123",
  "Attributes": { "RequestedAt": "/api/v1/sites/1/data" },
  "Results": []
}
```

### Envelope Key Rules

- **HTTP status is the primary indicator** — frontend checks 2xx, not `IsSuccess`
- **`Results` is always an array** — even for single items
- **`Attributes` always includes `RequestedAt`** — the route path
- All response keys use **PascalCase** via `ResponseKeyType` enum
- **`Errors.DelegatedRequestServer`** — present when Go proxied to a downstream server that failed
- **`MethodsStack`** — Go call chain that handled the request

---

## Go `apperror` Package

### Single Return Value — `Result[T]`

Go functions **never** return `(T, error)`. They return `apperror.Result[T]` which is a single value containing either the result or the error.

```go
// ❌ FORBIDDEN — dual return
func GetUser(id string) (*User, error) { ... }

// ✅ REQUIRED — single Result[T]
func GetUser(id string) apperror.Result[*User] { ... }
```

### Result[T] API

| Method | Purpose |
|--------|---------|
| `result.HasError()` | Check if operation failed |
| `result.Value()` | Get the success value (only after HasError check) |
| `result.Error()` | Get the error (only when HasError is true) |
| `result.PropagateError()` | Re-wrap and return error to caller |

### Error Wrapping — Always `apperror.Wrap()`

```go
// ❌ FORBIDDEN — raw fmt.Errorf
return fmt.Errorf("failed to open: %w", err)

// ✅ REQUIRED — structured wrapping with code
return apperror.Wrap(
    err,
    ErrDbOpen,
    "open database",
).WithContext("path", dbPath)
```

### Error Context Methods

| Method | Purpose | Example |
|--------|---------|---------|
| `.WithContext(key, value)` | Add key-value context | `.WithContext("path", "/db/main.db")` |
| `.WithPath(path)` | Shorthand for file path context | `.WithPath(configPath)` |
| `.WithStatusCode(code)` | Attach HTTP status | `.WithStatusCode(404)` |

### Error Propagation Pattern

```go
result := doSomething()
if result.HasError() {
    return result.PropagateError()
}
value := result.Value()
```

### 🔴 CODE RED — Swallowed Errors

Swallowing errors is the **highest severity violation**.

```go
// ❌ CODE RED — swallowed error (underscore discard)
result, _ := doSomething()

// ❌ CODE RED — empty catch
} catch (e) {}

// ❌ CODE RED — generic message without path
return errors.New("file not found")  // WHICH file?

// ✅ REQUIRED — always include context
return apperror.Wrap(err, ErrFileNotFound, "read config").WithPath(configPath)
```

---

## Tier 1: Delegated Server Error Handling (PHP)

### Safe Execution Pattern

Every REST endpoint handler is wrapped in `safeExecute`:

```php
public function handleRequest(WP_REST_Request $request): WP_REST_Response {
    return $this->safeExecute(
        fn() => $this->executeRequest($request),
        'request-name',
    );
}
```

`safeExecute` catches `Throwable`, logs via `FileLogger`, and returns a structured error envelope. Debug mode controls whether stack traces appear in the response.

### FileLogger — Two-Tier Logging

| Tier | Class | When Available | Purpose |
|------|-------|---------------|---------|
| 1 | `error_log()` / `ErrorLogHelper` | Always | Fallback before autoloader loads |
| 2 | `FileLogger` (singleton) | After autoloader | Primary structured log with rotation, dedup, stack traces |

**Log files** (under `wp-content/uploads/{plugin-slug}/logs/`):

| File | Contents |
|------|----------|
| `info.log` | All successful operations + debug entries (when debug mode ON) |
| `error.log` | All errors + exceptions with stack traces |
| `stacktrace.log` | Dedicated stack trace file (6-frame limit) |

### TypedQuery — Database Result Envelopes

| Class | Purpose | Key Methods |
|-------|---------|-------------|
| `DbResult<T>` | Single-row query | `isDefined()`, `isEmpty()`, `hasError()`, `isSafe()`, `value()`, `error()`, `stackTrace()` |
| `DbResultSet<T>` | Multi-row query | `hasAny()`, `isEmpty()`, `count()`, `hasError()`, `isSafe()`, `items()`, `first()` |
| `DbExecResult` | INSERT/UPDATE/DELETE | `isEmpty()`, `hasError()`, `isSafe()`, `affectedRows()`, `lastInsertId()` |

```php
$tq = new TypedQuery($pdo);

// Single row
$result = $tq->queryOne(
    'SELECT * FROM plugins WHERE id = :id',
    [':id' => $id],
    fn(array $row): PluginInfo => PluginInfo::fromRow($row),
);

if ($result->hasError()) { /* handle */ }
if ($result->isEmpty()) { /* not found */ }
$plugin = $result->value();
```

---

## DelegatedRequestServer — Proxy Error Capture

When the Go backend proxies to a downstream server (PHP, Node.js, etc.) and the downstream server fails, the Go backend captures full diagnostics:

| Field | Type | Description |
|-------|------|-------------|
| `Endpoint` | string | Full URL of the downstream request |
| `Method` | string | HTTP method used (GET, POST, etc.) |
| `StatusCode` | int | HTTP status from downstream |
| `Response` | string | Raw response body from downstream |
| `StackTrace` | string | Downstream stack trace (if provided) |
| `RequestBody` | string/null | Request body sent (POST/PUT/PATCH only) |

This block appears in `Envelope.Errors.DelegatedRequestServer` and is essential for debugging proxy-chain errors.

---

## Error Code Registry

### Structure

- Master registry: `error-codes-master.json`
- Ranges assigned per module to prevent collisions
- Overlap validator script enforces no duplicate ranges

### Module Ranges

| Range | Prefix | Module |
|-------|--------|--------|
| 1000–1999 | `GEN` | General/Shared (Core, Database, Auth, File System) |
| 2000–2999 | `SM` | Spec Management Software |
| 7000–7099 | `GS` | GSearch CLI Core |
| 7100–7599 | `BR` | BRun CLI |
| 7600–7919 | `GS` | GSearch sub-modules (Movie, BI, Multi-Source, etc.) |
| 8000–8399 | `NF` | Nexus Flow |
| 9000–9499 | `AB` | AI Bridge Core |
| 9500–9599 | `PS` | PowerShell Integration |
| 9600–9999 | `AB` | AI Bridge Extended modules |
| 10000–10499 | `WPB` | WP Plugin Builder |
| 11000–11999 | `SRC` | Spec Reverse CLI |
| 12000–12599 | `WSP` | WP SEO Publish |
| 13000–13999 | `WPP` | WP Plugin Publish |
| 14000–14499 | `AIT` | AI Transcribe CLI |
| 14500–14999 | `EQM` | Exam Manager |
| 15000–15999 | `LM` | Link Manager |
| 16000–16799 | `SM-CG` | SM Code Generation |
| 17000–17999 | `SM-PE` | SM Project Editor |
| 19000–19019 | `AB` | AI Bridge Lovable Reasoning |
| 20000–20999 | `AB-TR` | AI Bridge Non-Vector RAG |

### Code Registration Format

```json
{
  "Code": 1001,
  "Name": "ErrDbConnectionTimeout",
  "Module": "core/database",
  "Severity": "error",
  "Message": "Database connection timed out",
  "Resolution": "Check database path exists and is not locked"
}
```

### Overlap Prevention Rules

- Each module owns a contiguous integer range — no gaps, no overlaps
- Collision validator script runs in CI
- When ranges collide, the newer module gets reassigned (13 resolutions documented)
- PS/AB SEO share 9500-9540 intentionally — distinguishable by format (prefixed vs flat integer)

---

## Frontend Error Handling

### Error Store (Zustand)

```typescript
interface ErrorState {
  errors: AppError[];
  addError: (error: AppError) => void;
  clearErrors: () => void;
  dismissError: (id: string) => void;
}
```

### Global Error Modal

- Displays when `errors.length > 0`
- **Tabbed interface:** Overview, Backend Details, Delegated Server, Session Info
- Shows error code, message, and suggested resolution
- **Copy button (split button):** Main click → Compact Report (instant), dropdown → Full Report / With Backend Logs / error.log.txt / log.txt
- **Download menu:** Full Bundle (ZIP), error.log.txt, log.txt, Report (.md)
- Dismiss button clears the error
- Auto-dismiss for non-critical errors after 5 seconds

### Compact Report (Default Copy Format)

The default copy output includes:
- Error code, message, HTTP status
- Request URL, method, timestamp
- **Delegated Server Info** (built from `CapturedError.envelopeErrors.DelegatedRequestServer`) — no API call needed
- Frontend execution chain

### Toast Notifications

| Type | Color | Auto-dismiss |
|------|-------|-------------|
| Success | Green | 3s |
| Warning | Yellow | 5s |
| Error | Red | Persists until dismissed |
| Info | Blue | 3s |

---

## Notification Colors

| Notification Type | Color | Hex |
|-------------------|-------|-----|
| Success | Green | `#22c55e` |
| Warning | Amber | `#f59e0b` |
| Error | Red | `#ef4444` |
| Info | Blue | `#3b82f6` |

---

## Debugging Cheat Sheet

### Initialization Order (ALL Languages)

```
1. Configuration    → Load env vars and config files FIRST
2. Directories      → Ensure all required directories exist
3. Database         → Initialize connections (only after dirs exist)
4. Services         → Initialize business logic components
5. Server/App       → Start ONLY after all dependencies ready
```

### Common Pitfalls

| Symptom | Likely Cause | Resolution |
|---------|--------------|------------|
| "Backend disconnected" but backend is running | Response format mismatch | Compare handler output to frontend detection logic |
| 404 on API base URL | No index route registered | Check router for `GET /api/v1` handler |
| VITE_API_URL shows wrong value | Resolved vs raw env confusion | Distinguish raw env var from resolved origin |
| HTML instead of JSON | SPA fallback serving `index.html` | Check if route exists in backend router |
| CORS errors | Missing CORS headers | Check backend CORS middleware configuration |
| 401/403 on protected routes | Token not sent or expired | Check Authorization header, token validity |
| Plugin won't activate | PDO SQLite extension missing | Check `extension_loaded('pdo_sqlite')` |
| Database connection fails | Directory permissions | Check path exists and is writable |

### Quick Debug Commands

```bash

# Go — check server

curl -s http://localhost:8080/api/v1/health | jq .
tail -f logs/app.log

# PHP — check logs

cat wp-content/uploads/{plugin-slug}/logs/error.log

# TypeScript — check API

console.log(import.meta.env.VITE_API_URL)
```

---

## Core Principles

1. **Never assume — always verify** both backend response and frontend detection
2. **HTTP status is primary indicator** — use 2xx, not response body fields
3. **Structured errors only** — no unstructured error strings
4. **Every error gets context** — path, entity ID, operation name
5. **Stack traces in development** — stripped in production responses
6. **Envelope is universal** — every API endpoint returns the same structure
7. **Delegated server errors propagate** — Go captures and forwards downstream diagnostics

---

## Forbidden Patterns

| Pattern | Why | Required Alternative |
|---------|-----|---------------------|
| `catch (Exception $e)` | Misses PHP 7+ `Error` types | `catch (Throwable $e)` |
| `error_log()` for diagnostics | No structure | `FileLogger` / `Logger` |
| Empty `catch` block | Swallowed error — CODE RED | Handle or re-throw |
| `fmt.Errorf("failed: %w", err)` | No error code | `apperror.Wrap(err, ErrCode, "context")` |
| `result, _ := fn()` | Swallowed error — CODE RED | Check `result.HasError()` |
| Generic "file not found" | No path context — CODE RED | Include exact file path |
| `(T, error)` dual return in Go | Breaks single-return pattern | `apperror.Result[T]` |
| Inline error strings | Not machine-parseable | Error code constants |

---

## Validation

Run `linter-scripts/validate-guidelines.py` — zero **CODE-RED** violations required.

---

## 18. Retrospectives

**Source:** `01-error-resolution/03-retrospectives/` (4 retrospective files)

Every production error escalation produces a retrospective file in this folder. Required sections:

1. **What happened** — symptom, blast radius, time to detect / mitigate / resolve.
2. **Root cause** — exact line of code or config drift.
3. **Why guards failed** — which spec rule should have caught it (cite section).
4. **Spec change** — which spec file gets a new rule, with the diff.
5. **Linter change** — new rule ID added to `linter-scripts/`.
6. **Memory update** — new entry in `mem://issues/` if the lesson is project-wide.

Existing retrospectives cover: health-endpoint mismatch, retry/debounce/dedup race, ZIP finalization-before-return bug, activation-endpoint casing fix.

---

## 19. Verification Patterns

**Source:** `01-error-resolution/04-verification-patterns/`

Documents the **frontend↔backend sync verification** pattern: every fixed bug must include a paired test (one in the producing tier, one in the consuming tier) that would have failed before the fix. CI gate: PR description must include both test paths or be rejected.

---

## 20. Per-Language Debugging Guides

**Source:** `01-error-resolution/05-debugging-guides/` (PHP, Go, TypeScript guides)

| Language | Required Tools | Mandatory First-Step |
|----------|---------------|---------------------|
| PHP | Xdebug, `FileLogger` tail, WP debug.log | `tail -f wp-content/debug.log` while reproducing |
| Go | `dlv`, session log file, `pprof` for perf bugs | Inspect `error.log.txt` for the request ID |
| TypeScript | Browser devtools, React DevTools, `errorStore` snapshot | Open the Global Error Modal → Copy Full Report |

Every debugging session **must** capture: timestamp, request ID, full stack trace, file paths involved. Verbal "it's broken" reports are insufficient.

---

## 21. Notification Colors

**Source:** `02-error-architecture/03-notification-colors.md`

Toast and modal severity colors are defined as CSS variables and **must not be hard-coded** in components.

| Severity | Token | Default HSL |
|----------|-------|-------------|
| Success | `--toast-success` | `142 71% 45%` |
| Info | `--toast-info` | `217 91% 60%` |
| Warning | `--toast-warning` | `38 92% 50%` |
| Error | `--toast-error` | `0 84% 60%` |
| Critical (CODE-RED) | `--toast-critical` | `0 100% 35%` (deeper red, plus pulse animation) |

The Global Error Modal uses these same tokens for the severity badge.

---

## 22. Error Modal System

**Source:** `02-error-architecture/04-error-modal/` (overview + 6 sub-areas: copy formats, React components, full reference, color themes, history persistence, suppression)

### 22.1 Modal Tabs

The Global Error Modal has **5 tabs**: Summary, Stack, Request, Response, Click Path. Each tab renders from the same `CapturedError` shape — never duplicate state.

### 22.2 Copy Formats (3)

| Format | Use Case | Function |
|--------|---------|----------|
| Full Report | Bug report to engineer | `generateErrorReport(error, app)` |
| Compact Report | Slack / chat paste | `generateCompactReport(error, app)` |
| Raw JSON | Programmatic re-ingest | `JSON.stringify(error, null, 2)` |

Source files: `src/components/errors/errorReportGenerator.ts` (pure functions, no React deps), `errorLogAdapter.ts` (maps backend `ErrorLog` → `CapturedError`).

### 22.3 History Persistence

Errors are persisted to IndexedDB with a 50-entry rolling window keyed by capture timestamp. The history pane is keyboard-navigable (`↑`/`↓`, `Enter` to inspect).

### 22.4 Suppression

`suppressGlobalError(predicate)` allows hiding the modal for known-recoverable errors (e.g., 401 on a probe request). Suppression is **logged** even when the modal is hidden — never silently swallow.

---

## 23. Response Envelope Reference

**Source:** `02-error-architecture/05-response-envelope/` (6 docs + 6 JSON fixtures + JSON schema)

### 23.1 Schema Fixtures

Every envelope variant has a fixture in the source folder for snapshot testing:

| Fixture | Scenario |
|---------|----------|
| `envelope-minimal.json` | Success with no Meta |
| `envelope-single.json` | Single typed result |
| `envelope-multiple.json` | Multiple result blocks |
| `envelope-error.json` | Error envelope with `Errors` array |
| `envelope-debug.json` | Debug-mode envelope including `MethodsStack` |
| `envelope.schema.json` | JSON Schema 2020-12 for the envelope |

CI runs every backend response through schema validation; mismatches fail the build.

### 23.2 Configurability

- `Meta` block is freely extensible (only key-naming rule applies).
- `Errors` block has a **fixed shape**: `Code`, `Message`, `Detail`, `StackTrace`, `Context`, `RequestId`. Additional keys are forbidden.
- `MethodsStack` and `Attributes` are debug-only — must be stripped in production builds via `STRIP_DEBUG_ENVELOPE=1`.

---

## 24. apperror Package Reference

**Source:** `02-error-architecture/06-apperror-package/01-apperror-reference/` (8 modules: overview, AppError struct, Result types, codes/policy, AppErrType enums, usage/adapters, serialization/guards)

### 24.1 Core Types

```go
type StackTrace struct {
    Frames []StackFrame  // file, line, function, package
}

type AppError struct {
    Code      AppErrType        // E1xxx–E14xxx domain enum
    Message   string            // human-readable
    Cause     error             // wrapped underlying
    Stack     StackTrace        // captured at Wrap()
    Context   map[string]any    // PascalCase keys
}

type Result[T any] struct {
    value *T
    err   *AppError
}

type ResultSlice[T any] = Result[[]T]
type ResultMap[K comparable, V any] = Result[map[K]V]
```

### 24.2 Constructors

| Function | When to use |
|----------|------------|
| `apperror.Ok(v)` | Success path, return value |
| `apperror.Err` (generic on T) | Failure path, propagate typed `*AppError` (signature: `apperror.Err` accepts `*AppError` and returns `Result` of T) |
| `apperror.Wrap(err)` | Convert any `error` to `*AppError`, capturing stack at this point |
| `apperror.New(code, msg)` | Construct from scratch with no underlying cause |

### 24.3 Stack Trace Skip Rules

`Wrap()` skips the wrap call itself + `runtime.Callers`. `New()` skips one extra frame. Service-layer adapters that immediately re-wrap should use `WrapSkip(2)` to keep the trace pointing at the originating site.

### 24.4 AppErrType Enum Ranges

| Range | Domain |
|-------|--------|
| E1xxx | Configuration / startup |
| E2xxx | Validation |
| E3xxx | Authentication / authorization |
| E4xxx | Database |
| E5xxx | Filesystem / IO |
| E6xxx | Network / HTTP client |
| E7xxx | External service (third-party APIs) |
| E8xxx | Business logic |
| E9xxx | Concurrency / queue |
| E10xxx | Cache |
| E11xxx | Cryptography |
| E12xxx | Self-update |
| E13xxx | CI/CD pipeline |
| E14xxx | Reserved for app-specific extensions |

Full enum definitions live in `05-apperrtype-enums.md` (340 lines). Adding a new code requires registry update (§26).

### 24.5 Service Adapter Unwrap Pattern

```go
func (s *Service) ProcessOrder(id string) apperror.Result[Order] {
    rawResult := s.repo.LoadOrder(id)
    if rawResult.HasError() {
        return apperror.Err[Order](rawResult.Error().
            WithCode(apperror.E8001OrderLoadFailed).
            WithContext("OrderId", id))
    }
    return apperror.Ok(rawResult.Unwrap())
}
```

The adapter pattern preserves the original stack trace while adding domain-specific context. **Never** call `.Unwrap()` without a preceding `.HasError()` guard — this is a CODE-RED.

### 24.6 JSON Serialization

`AppError.MarshalJSON()` emits: `{"Code", "Message", "Detail", "StackTrace": [...], "Context": {...}}`. Stack frames serialize as `{"File", "Line", "Function", "Package"}`. The `Cause` chain is flattened into `Detail` to keep the wire format flat.

---

## 25. Logging & Diagnostics

**Source:** `02-error-architecture/07-logging-and-diagnostics/` (3 files: overview, React execution logger, session-based logging)

### 25.1 Two-Tier Logging (Backend)

| Tier | When | Destination | Rotation |
|------|------|-------------|----------|
| Bootstrap (`ErrorLogHelper`) | Before session is established | `error.bootstrap.log` | Manual / weekly |
| Runtime (`FileLogger`) | After request enters handler | `sessions/<RequestId>/error.log.txt` | Per-request directory |

The bootstrap tier exists because session context isn't available before middleware runs. Runtime tier is mandatory for everything else.

### 25.2 Session-Based Logging

- Every request gets a UUID `RequestId` injected as the first middleware.
- All logs for that request go to `sessions/<RequestId>/`.
- Directory is preserved for **24 hours** then garbage-collected by a scheduled task.
- The `RequestId` is echoed in the response envelope's `Meta.RequestId` for client correlation.

### 25.3 React Execution Logger (Debug Mode)

When `localStorage.debugExecutionLog === '1'`, the frontend captures every component render, hook call, and event handler invocation into a circular buffer. Buffer is dumped into the Global Error Modal's Click Path tab on capture. Disabled in production builds via tree-shaking.

---

## 26. Error Code Registry

**Source:** `03-error-code-registry/` (registry, integration guide, collision resolution, utilization report, overlap validator, schemas, linter scripts, templates)

### 26.1 Registry File

`error-codes-master.json` is the **single source of truth** for every error code across all languages (Go, PHP, TypeScript). Schema:

```json
{
  "Code": "E2001",
  "Name": "ValidationFailed",
  "Domain": "Validation",
  "Severity": "Error",
  "MessageTemplate": "Validation failed for field '{Field}'",
  "OwningSpec": "spec/03-error-manage/...",
  "FirstSeenVersion": "1.4.0"
}
```

### 26.2 Adding a New Code (4-Step Process)

1. Append entry to `error-codes-master.json` (validated by `07-schemas/error-codes.schema.json`).
2. Run `python3 spec/03-error-manage/03-error-code-registry/08-linter-scripts/check-code-collisions.py` — must exit 0.
3. Run code generators to emit `apperror_codes_generated.go`, `ErrorCode.php`, `errorCodes.generated.ts` from the master.
4. Update the relevant retrospective or feature spec to reference the new code.

### 26.3 Collision Rules

- Code numbers are **immutable** once shipped — never re-use.
- Names within a Domain must be unique.
- Deprecated codes get `Deprecated: true` + `ReplacedBy: "EXXXX"` instead of deletion.

### 26.4 Utilization Report

`04-error-code-utilization-report.md` is regenerated weekly: lists every code, the files that emit it, and the count of test cases. Codes with zero emitters for 30+ days trigger a "candidate for removal" flag.

---

## Cross-References

| Topic | Source Location |
|-------|----------------|
| Error Handling Reference | `03-error-manage/02-error-architecture/01-error-handling-reference.md` |
| Error Code Registry | `03-error-manage/03-error-code-registry/01-registry.md` |
| Debugging Cheat Sheet | `03-error-manage/01-error-resolution/02-debugging-cheat-sheet.md` |
| Error Modal Copy Formats | `03-error-manage/02-error-architecture/04-error-modal/01-copy-formats/` |
| Notification Colors | `03-error-manage/02-error-architecture/03-notification-colors.md` |
| Go Delegation Fix | `03-error-manage/02-error-architecture/02-go-delegation-fix.md` |
| Collision Resolution | `03-error-manage/03-error-code-registry/03-collision-resolution-summary.md` |
| Overlap Validator | `03-error-manage/03-error-code-registry/05-overlap-validator.md` |
| Session-Based Logging | `03-error-manage/07-logging-and-diagnostics/02-session-based-logging.md` |

---

---

## §27 Live Error Code Registry Snapshot

This is a **point-in-time snapshot** of the canonical error code allocations. Source of truth: `spec/03-error-manage/03-error-code-registry/error-codes-master.json`. A blind AI must consult this section before allocating new error codes to avoid range collisions.

### 27.1 Allocated Ranges by Subsystem (16 modules · 933 codes · 159 retryable)

| Project | Name | Range | Codes | Retryable | Status |
|---------|------|-------|-------|-----------|--------|
| `GEN` | General Error Codes | 0–999 | — | — | Special — cross-cutting |
| `SM` | Spec Management | 2000–2999 | 0 | 0 | ⏳ Pending |
| `SM` (misc) | Spec Editor + Error Recovery | 6010–6013, 6020–6027 | — | — | Special |
| `BR` | BRun CLI | 7100–7599 | 0 | 0 | ⏳ Pending |
| `GS` | GSearch CLI | 7000–7919 | 102 | 37 | ✅ Active |
| `NF` | Nexus Flow CLI | 8000–8349 | 15 | 1 | ✅ Active |
| `AB` | AI Bridge CLI | 9000–9999 | 168 | 8 | ✅ Active |
| `PS/AB` | PowerShell/AB SEO Overlap | 9500–9540 | — | — | Special — format-separated |
| `WPB` | WP Plugin Builder | 10000–10499 | 70 | 9 | ✅ Active |
| `SRC` | Spec Reverse CLI | 11000–11999 | 0 | 0 | ⏳ Pending |
| `WSP` | WP SEO Publish CLI | 12000–12599 | 75 | 4 | ✅ Active |
| `WPP` | WP Plugin Publish | 13000–13499 | 46 | 1 | ✅ Active |
| `AIT` | AI Transcribe CLI | 14000–14499 | 141 | 27 | ✅ Active |
| `EQM` | Exam Manager | 14500–14999 | 57 | 5 | ✅ Active |
| `LM` | License Manager | 15000–15999 | 0 | 0 | ⏳ Pending |
| `SM-CG` | Spec Mgmt — Code Generation | 16000–16799 | 79 | 32 | ✅ Active |
| `SM-PE` | Spec Mgmt — Project Editor | 17000–17999 | 82 | 3 | ✅ Active |
| `SM-GS` | Spec Mgmt — GSearch (remap) | 18000–18249 | 92 | 32 | ✅ Active — ecosystem remap of GS |
| `AB-LR` | AI Bridge — Lovable Reasoning | 19000–19049 | 6 | 0 | ✅ Active |

### 27.2 Free Ranges Available for Allocation

| Range | Note |
|-------|------|
| 1000–1999 | Reserved (post-GEN buffer — do not allocate without RFC) |
| 3000–5999 | Reserved (post-SM buffer) |
| 8350–8999 | Free — between NF and AB |
| 10500–10999 | Free — between WPB and SRC |
| 12600–12999 | Free — between WSP and WPP |
| 13500–13999 | Free — freed from WPP compression |
| 18250–18999 | Free — after SM-GS |
| 19050–19999 | Free — after AB-LR |

**Rule:** A new module gets a 500- or 1000-code range from the **Free** list above. Choose the smallest range that fits projected growth × 3.

### 27.3 General Categories (GEN range 0–999)

Cross-cutting categories that every module inherits — never reallocate:

| Category | Example Codes |
|----------|---------------|
| Initialization | 0–99 |
| Authentication | 100–199 |
| Authorization | 200–299 |
| Validation | 300–399 |
| Business Logic | 400–499 |
| Database | 500–599 |
| Type Casting | 600–699 |
| File System | 700–799 |
| Network | 800–899 |
| Reserved | 900–999 |

### 27.4 Code Allocation Workflow

```
1. Pick a free range from §27.2 (or use existing module range if extending)
2. Edit spec/<your-module>/error-codes.json — add new code(s) following the schema
3. Run: node spec/03-error-manage/03-error-code-registry/08-linter-scripts/detect-collisions.mjs
4. Run: node spec/03-error-manage/03-error-code-registry/08-linter-scripts/validate-master-stats.mjs
5. Run code generators (§27.5) to emit per-language artifacts
6. Run: node spec/03-error-manage/03-error-code-registry/08-linter-scripts/check-utilization-threshold.mjs
7. Update error-codes-master.json TotalCodes / RetryableCodes (auto via generate-utilization-report.mjs)
8. Commit all artifacts together
```

### 27.5 Code Generators — Per-Language Emitters

| Generator | Reads | Emits | Language |
|-----------|-------|-------|----------|
| `gen-go-errcodes` | `spec/<module>/error-codes.json` | `internal/apperror/apperror_codes_generated.go` | Go |
| `gen-php-errcodes` | `spec/<module>/error-codes.json` | `src/Errors/ErrorCode.php` | PHP |
| `gen-ts-errcodes` | `spec/<module>/error-codes.json` | `src/lib/errors/errorCodes.generated.ts` | TypeScript |
| `gen-rust-errcodes` | `spec/<module>/error-codes.json` | `src/errors/error_code_generated.rs` | Rust |

**Invocation pattern** (run from repo root for each module):

```bash
node scripts/codegen/gen-go-errcodes.mjs   --module GS
node scripts/codegen/gen-php-errcodes.mjs  --module GS
node scripts/codegen/gen-ts-errcodes.mjs   --module GS
node scripts/codegen/gen-rust-errcodes.mjs --module GS
```

**Drift detection:** CI runs all generators and `git diff --exit-code` on the output paths. Any uncommitted generator output is a CI failure (`Error: generated error code drift in <file>`).

### 27.6 Registry Linter Scripts

| Script | Purpose | Exit Codes |
|--------|---------|------------|
| `detect-collisions.mjs` | Verifies no two codes share the same numeric value across modules | `0` clean · `1` collision |
| `validate-master-stats.mjs` | Checks `error-codes-master.json` `TotalCodes`/`RetryableCodes` match per-module files | `0` clean · `1` drift |
| `generate-utilization-report.mjs` | Produces `04-error-code-utilization-report.md` | `0` always |
| `check-utilization-threshold.mjs` | Warns if a module is < 10% utilized (encourages range compression) | `0` clean · `1` under-threshold |

All four live at `spec/03-error-manage/03-error-code-registry/08-linter-scripts/`.

### 27.7 Schema Reference

The per-module `error-codes.json` schema is documented at:
- `spec/03-error-manage/03-error-code-registry/07-schemas/error-code.schema.json`
- `spec/03-error-manage/03-error-code-registry/07-schemas/error-codes-index.schema.json`

**Minimum fields per code:**
```json
{
  "Code": 7042,
  "Name": "RagChunkNotFound",
  "Category": "RAG",
  "Retryable": false,
  "HttpStatus": 404,
  "MessageTemplate": "RAG chunk %s not found in collection %s",
  "Severity": "error"
}
```

### 27.8 Ecosystem Remap Pattern

A module that uses **local codes** internally (e.g., GS uses 1xxx–12xxx) but participates in the ecosystem **remaps** to a sub-range of `SM-*`. Example:

| Module | Local Range | Ecosystem Remap |
|--------|-------------|-----------------|
| GS | 1000–12999 (internal) | SM-GS 18000–18249 |

Remap mapping lives at the module's `EcosystemRemapIndex` path. The remap is **bidirectional** — both indexes must agree, enforced by `detect-collisions.mjs`.

### 27.9 Collision Resolution

13 historical collisions have been resolved (full log: `03-collision-resolution-summary.md`). The two intentional overlaps are:
- **PS/AB SEO 9500–9540** — format-separated (`PS-9500-00` vs flat `9500`).
- **GS local vs SM-GS remap** — by design via the ecosystem remap pattern.

**Never** introduce a new intentional collision without RFC.

---

*Live Error Code Registry Snapshot added — v3.3.0 — 2026-04-22 — sourced from `error-codes-master.json` Generated: 2026-02-28*
