# AppError Package Reference — JSON serialization, Result guard rule

> **Parent:** [AppError Package Reference](./00-overview.md)  
> **Version:** 1.3.0  
> **Updated:** 2026-03-31

---

## 11. JSON Serialization (Invariant I-1)

> **Core Guarantee:** Every `*AppError` is fully round-trippable through JSON serialization. `json.Marshal(err)` produces a complete diagnostic payload; `json.Unmarshal(data, &err)` reconstructs it with code, message, details, values, diagnostics, stack trace, and cause message preserved. This enables error transport across HTTP APIs, subprocess protocols, database storage, and log aggregation — boundaries where raw `error` fails silently.

### 11.1 JSON Tag Convention

**Rule:** Only add `json:"..."` tags when the JSON key **differs** from the Go field name OR when `omitempty` is needed. PascalCase field names serialize as PascalCase by default — redundant tags must be removed.

| Scenario | Tag Required? | Example |
|----------|--------------|---------|
| Field name matches JSON key | ❌ No | `Version string` (serializes as `"Version"`) |
| Field needs omitempty only | ✅ Yes | `Details string \`json:",omitempty"\`` |
| Excluded from JSON | ✅ Yes | `Cause error \`json:"-"\`` |
| External API contract | ✅ Yes | `AccessToken string \`json:"access_token"\` // EXEMPTED: External API` |

### 11.2 Existing JSON Tags

All core structs use implicit PascalCase (redundant tags removed). Only functional tags remain:

- `AppError` — `Details`, `Values`, `Diagnostic` use `json:",omitempty"`; `Cause` uses `json:"-"` (EXEMPTED: I-2)
- `StackTrace` — `PreviousTrace` uses `json:",omitempty"`
- `StackFrame` — no tags (all fields always present)
- `ErrorDiagnostic` — all 15+ fields use `json:",omitempty"` only (redundant key names removed)

### 11.2 Custom MarshalJSON

The `Cause` field is excluded from default JSON marshaling (`json:"-"`) because it's a Go `error` interface. A custom `MarshalJSON` must serialize the cause message as a string:

**File:** `backend/pkg/apperror/error_json.go`

```go
func (e *AppError) MarshalJSON() ([]byte, error) {
    type alias AppError

    return json.Marshal(&struct {
        *alias
        CauseMessage string `json:"Cause,omitempty"` // Maps CauseMessage → Cause key
    }{
        alias:        (*alias)(e),
        CauseMessage: causeMessage(e),
    })
}

func causeMessage(e *AppError) string {
    if e.Cause == nil {
        return ""
    }

    return e.Cause.Error()
}
```

**Rules:**
- Uses type alias to prevent infinite recursion
- `Cause` serialized as `"Cause"` string field (not nested error object)
- Empty cause omitted via `omitempty`

### 11.3 Custom UnmarshalJSON

Reconstructs `Cause` from the serialized string:

```go
func (e *AppError) UnmarshalJSON(data []byte) error {
    var alias appErrorJson

    if err := json.Unmarshal(data, &alias); err != nil {
        return fmt.Errorf("apperror.UnmarshalJSON: failed to decode AppError (received %d bytes: %s): %w",
            len(data), truncateData(data, 200), err)
    }

    e.Code = alias.Code
    e.Message = alias.Message
    e.Details = alias.Details
    e.Values = alias.Values
    e.Diagnostic = alias.Diagnostic
    e.Stack = alias.Stack

    if alias.CauseMessage != "" {
        e.Cause = &plainError{msg: alias.CauseMessage}
    }

    return nil
}

// truncateData returns a string preview of raw JSON data, capped at maxLen bytes.
func truncateData(data []byte, maxLen int) string {
    if len(data) <= maxLen {
        return string(data)
    }

    return string(data[:maxLen]) + "..."
}
```

**Rules:**
- Reconstructed `Cause` is a `plainError` struct — the original type is lost (acceptable for deserialization)
- Stack trace, values, and diagnostics are fully preserved
- **Error messages include the raw received data** (truncated to 200 bytes) for debugging malformed payloads

### 11.4 Serialization Output Example

```json
{
    "Code": "E3001",
    "Message": "failed to connect to WordPress",
    "Details": "dial tcp 192.168.1.100:443: connect: connection refused",
    "Values": {
        "url": "https://example.com",
        "plugin": "my-plugin"
    },
    "Diagnostic": {
        "Url": "https://example.com/wp-json/wp/v2/plugins",
        "StatusCode": 0,
        "Method": "GET"
    },
    "Stack": {
        "Frames": [
            {"Function": "wordpress.(*Client).ListPlugins", "File": "client.go", "Line": 42},
            {"Function": "sync.(*Service).CheckSync", "File": "service.go", "Line": 128}
        ]
    },
    "Cause": "dial tcp 192.168.1.100:443: connect: connection refused"
}
```

---


---

## 12. Result Guard Rule — Mandatory Error Check Before Value Access

Every call site that receives a `Result[T]`, `ResultSlice[T]`, or `ResultMap[K, V]` (Go) or `DbResult`, `DbResultSet`, `DbExecResult` (PHP) **MUST** check `HasError()` / `hasError()` or `IsSafe()` / `isSafe()` before calling `.Value()` / `.value()`, `.Items()` / `.items()`, or `.Get()`. Accessing the contained value without a guard is a **spec violation**.

**Principle:** No error may ever be swallowed. If a result carries an error, it must be explicitly handled — logged, returned, or propagated. The framework-level accessor should log immediately when called on an errored result, reducing diagnostic steps.

> **Important:** In Go, the method to retrieve the error is named `.AppError()` (not `.Error()`) to avoid confusion with Go's native `error` interface. `.AppError()` returns `*apperror.AppError` which carries the full stack trace, error code, and diagnostic context. In PHP, the equivalent method remains `.error()` returning `Throwable`.

### Go Examples

#### ❌ WRONG — No Guard

```go
// Silent failure: if result has an error, Value() panics or returns zero
result := svc.GetById(ctx, id)
plugin := result.Value()
```

#### ✅ CORRECT — Direct Propagation (Same Type)

```go
// result is already Result[T] with the error — just return it
result := svc.GetById(ctx, id)

if result.HasError() {
    return result
}

plugin := result.Value()
```

#### ✅ CORRECT — Cross-Type Propagation (Result[T] → Result[U])

```go
// When the return type differs, Fail re-wrapping IS needed:
siteResult := siteSvc.GetById(ctx, siteId)

if siteResult.HasError() {
    return apperror.Fail[PluginList](siteResult.AppError())
}
```

#### ❌ WRONG — Redundant Re-Wrapping (Same Type)

```go
// result.AppError() is already *AppError — no need to re-wrap into same Result[T]
result := svc.GetById(ctx, id)
if result.HasError() {
    return apperror.Fail[Plugin](result.AppError()) // redundant
}
```

#### Same-Type vs Cross-Type — ResultSlice and ResultMap

The same rule applies to `FailSlice` and `FailMap`:

```go
// ✅ Same-type ResultSlice — direct return
func (s *Service) ListActive(context stdctx.Context) apperror.ResultSlice[models.Site] {
    result := s.List(context) // returns ResultSlice[models.Site]

    if result.HasError() {
        return result // same type — direct return
    }

    return apperror.OkSlice(filterActive(result.Items()))
}

// ❌ WRONG — redundant FailSlice re-wrapping (same type)
if result.HasError() {
    return apperror.FailSlice[models.Site](result.AppError()) // redundant
}

// ✅ Cross-type ResultSlice — FailSlice IS needed
func (s *Service) CheckAll(context stdctx.Context) apperror.ResultSlice[SyncResult] {
    plugins := s.pluginService.List(context) // returns ResultSlice[models.Plugin]

    if plugins.HasError() {
        return apperror.FailSlice[SyncResult](plugins.AppError()) // different T
    }
    // ...
}

// ✅ Same-type ResultMap — direct return
func (s *Service) GetCached(context stdctx.Context) apperror.ResultMap[string, Config] {
    result := s.loadAll(context) // returns ResultMap[string, Config]

    if result.HasError() {
        return result // same type — direct return
    }
    // ...
}

// ✅ Cross-type ResultMap — FailMap IS needed
if configResult.HasError() {
    return apperror.FailMap[string, Summary](configResult.AppError()) // different V
}
```

#### ❌ WRONG — Compound Guard (Over-Checking)

```go
// Never combine HasError() with additional IsSafe()/Value() conditions
if publishResult.HasError() || (publishResult.IsSafe() && publishResult.Value().IsFailed) {
    return publishResult // over-engineered — HasError() already covers failure
}

// ✅ CORRECT — HasError() is the single source of truth
if publishResult.HasError() {
    return publishResult
}
```

#### ❌ WRONG — Negating a Positive Boolean

```go
// PROHIBITED — negation operator on boolean field
if !result.IsSuccess {
    return loglevel.Error
}

// ✅ CORRECT — use the positive failure field
if result.IsFailed {
    return loglevel.Error
}
```

> **Rule:** Any struct with an `IsSuccess` field MUST also expose an `IsFailed` field (or method). Code must use `IsFailed` instead of `!IsSuccess`. This follows the universal no-negation-operator rule.

> **Rule:** `HasError()` is the only guard needed for Result error checking. Business-level success/failure flags on the contained value (e.g. `.IsSuccess`) are a separate concern and must not be combined into the Result guard expression. If you need to check domain-level outcomes, do so **after** the `HasError()` guard, on the unwrapped value.

#### ✅ CORRECT — Using IsSafe()

```go
result := svc.List(ctx)

if result.IsSafe() {
    for _, item := range result.Items() {
        process(item)
    }
}
```

#### ✅ CORRECT — Adapter Unwrap Pattern

```go
func (a *PluginServiceAdapter) GetById(context stdctx.Context, id int64) (*models.Plugin, error) {
    result := a.Service.GetById(context, id)

    if result.HasError() {
        return nil, result.AppError()
    }

    v := result.Value()

    return &v, nil
}
```

### PHP Examples

#### ❌ WRONG — No Guard (DbResult)

```php
$result = $query->queryOne(...);
$result->value(); // error silently swallowed
```

#### ✅ CORRECT — Guard Before Access (DbResult)

```php
$result = $query->queryOne(...);

if ($result->hasError()) {
    $this->logger->logException($result->error(), 'context');

    return null;
}

return $result->value();
```

#### ✅ CORRECT — Guard Before Iteration (DbResultSet)

```php
$results = $query->queryAll(...);

if ($results->hasError()) {
    $this->logger->logException($results->error(), 'query failed');

    return [];
}

return $results->items();
```

#### ✅ CORRECT — Guard Before Write Result (DbExecResult)

```php
$execResult = $query->execute(...);

if ($execResult->hasError()) {
    $this->logger->logException($execResult->error(), 'execute failed');

    return false;
}

return $execResult->affectedRows() > 0;
```

### Enforcement Checklist

- [ ] Every `result.Value()` / `$result->value()` call is preceded by `HasError()` / `hasError()` or `IsSafe()` / `isSafe()`. In Go, use `.AppError()` (not `.Error()`) to retrieve the structured error.
- [ ] Every `result.Items()` / `$results->items()` call is preceded by a guard
- [ ] Every `result.Get(key)` on `ResultMap` is preceded by a guard
- [ ] Every `$execResult->affectedRows()` on `DbExecResult` is preceded by a guard
- [ ] No error is silently discarded — all errors are logged, returned, or propagated
- [ ] Cross-service callers (direct `*service.Service` refs) guard results the same way

---
