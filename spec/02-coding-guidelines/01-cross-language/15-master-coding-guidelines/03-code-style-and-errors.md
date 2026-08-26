# Master Coding Guidelines — Code style formatting, error handling

> **Parent:** [Master Coding Guidelines](./00-overview.md)  
> **Version:** 2.1.0  
> **Updated:** 2026-03-31

---

## 5. Code Style — Formatting Rules

> Full reference: [code-style.md](../04-code-style/00-overview.md)

| Rule | Description |
|------|-------------|
| R1 | Always use braces — no single-line `if` |
| R2 | Zero nested `if` — absolute ban |
| R3 | Extract complex conditions into named booleans |
| R4 | Blank line before `return`/`throw` when preceded by statements |
| R5 | Blank line after `}` when followed by more code (5a: if, 5b: loops, 5c: try/switch) |
| R6 | Max 15 lines per function body |
| R7 | Zero nested `if` reinforcement |
| R9a | Function signatures >2 params → one per line with trailing comma |
| R9b | Function calls >2 args → one per line |
| R9c | PHP array literals >2 items → one per line |
| R10 | Blank line before control structures when preceded by assignments |
| R11 | Long string concatenations → line-by-line |
| R12 | No empty line after opening brace |
| R13 | No empty line at start of file |

### Common Mistakes

```php
// ❌ R2 violation: Nested if
if ($request !== null) {
    if ($request->hasParam('file')) {
        $this->process($request);
    }
}

// ✅ CORRECT: Early return + flat
if ($request === null) {
    return;
}

if ($request->hasParam('file')) {
    $this->process($request);
}
```

```php
// ❌ R4 violation: No blank line before return
$result = $this->compute($data);
return $result;

// ✅ CORRECT
$result = $this->compute($data);

return $result;
```

```go
// ❌ R6 violation: Function too long (>15 lines)
func ProcessUpload(ctx context.Context, req Request) error {
    // 25 lines of code...
}

// ✅ CORRECT: Decompose into helpers
func ProcessUpload(ctx context.Context, req Request) error {
    if err := validateUpload(req); err != nil {
        return err
    }

    result, err := executeUpload(ctx, req)
    if err != nil {
        return apperror.Wrap(err, apperror.ErrUploadFailed, "upload failed")
    }

    return logAndRespond(ctx, result)
}
```

---


---

## 6. Error Handling

### PHP

- Use `try/catch` with `Throwable` (unqualified, imported via `use`)
- Never use leading backslash: `\Throwable` → `Throwable`

### Go

- All errors via `apperror.New()` or `apperror.Wrap()` — automatic stack traces
- Never use `fmt.Errorf()` for service errors
- Service methods return `apperror.Result[T]` — never raw `(T, error)` or multi-value tuples (§7.1)
- Zero type assertions in business logic — use concrete typed structs (§7.2)
- Error codes follow `E{category}xxx` convention

### 6.0 — Serialization Invariant & No Raw `error` in Structs

**`*AppError` is fully serializable.** Every `*AppError` round-trips through JSON (and YAML) preserving code, message, details, values, diagnostics, stack trace, and cause message. This enables error transport across HTTP APIs, subprocess protocols, database storage, and log aggregation.

**Struct fields that hold errors must use `*AppError` (Go), `Throwable` (PHP), or the framework's structured error type (TypeScript) — never raw `error`, `string`, or untyped exceptions.**

```go
// ❌ FORBIDDEN — error interface is not serializable
type JobResult struct {
    Output string
    Err    error     // json.Marshal → {} — all context lost
}

// ✅ REQUIRED — *AppError serializes with full diagnostic context
type JobResult struct {
    Output   string
    AppError *apperror.AppError `json:",omitempty"`
}
```

```php
// ❌ FORBIDDEN — bare string error
class JobResult {
    public string $error;  // no stack trace, no code, no diagnostics
}

// ✅ REQUIRED — Throwable with stack trace
class JobResult {
    public ?Throwable $error;  // preserves stack trace and context
}
```

**Rationale:** Raw `error` is an opaque interface — it cannot be serialized, queried, or transported. `*AppError` carries structured data (code, stack, diagnostics) that survives every boundary: HTTP responses, subprocess JSON protocol, error history DB, AI diagnostic clipboard, and log aggregation.

See [apperror §10.6 — No Raw `error` in Struct Fields](../../../03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/05-usage-and-adapters.md#106-no-raw-error-in-struct-fields-invariant-i-2) for the full rule.

### 6.1 — Result Guard Rule (Zero Silent Failures)

Every Result/DbResult wrapper **MUST** have its error state checked before accessing the contained value. Accessing `.value()` / `.Value()` without a prior `hasError()` or `isSafe()` guard is a **spec violation**.

**Principle:** No error may ever be swallowed. If a result carries an error, it must be explicitly handled — logged, returned, or propagated. The framework-level `.value()` / `.Value()` accessor should log immediately when called on an errored result, reducing diagnostic steps. If an error exists, the accessor returns empty/zero and the framework logs the error automatically.

```php
// PHP — DbResult / DbResultSet / DbExecResult
$result = $query->queryOne(...);

// ❌ WRONG: No guard — error silently swallowed
$result->value();

// ✅ CORRECT: Guard before access
if ($result->hasError()) {
    $this->logger->logException($result->error(), 'context');

    return null;
}

return $result->value();
```

```php
// PHP — DbResultSet (collection access)
$results = $query->queryAll(...);

// ❌ WRONG: No guard — iterating potentially empty/errored set
foreach ($results->items() as $row) { ... }

// ✅ CORRECT: Guard before iteration
if ($results->hasError()) {
    $this->logger->logException($results->error(), 'query failed');

    return [];
}

return $results->items();
```

```php
// PHP — DbExecResult (write operations)
$execResult = $query->execute(...);

// ❌ WRONG: No guard — assuming success
$execResult->affectedRows();

// ✅ CORRECT: Guard before access
if ($execResult->hasError()) {
    $this->logger->logException($execResult->error(), 'execute failed');

    return false;
}

return $execResult->affectedRows() > 0;
```

```go
// Go — Propagation Rules (Result[T], ResultSlice[T], ResultMap[K,V])
// .AppError() returns *AppError — always preserves stack trace and context.
// Named AppError() (not Error()) to avoid confusion with Go's native error interface.

// ✅ Same-type → direct return (applies to Result, ResultSlice, ResultMap)
result := svc.GetById(ctx, id)           // Result[Plugin]

if result.HasError() { return result }    // no re-wrapping needed
plugin := result.Value()

// ✅ Cross-type → Fail/FailSlice/FailMap IS needed
plugins := s.pluginService.List(ctx)      // ResultSlice[Plugin]

if plugins.HasError() {
    return apperror.FailSlice[SyncResult](plugins.AppError())
}

// ❌ WRONG — redundant (same type re-wrapped)
if result.HasError() { return apperror.Fail[Plugin](result.AppError()) }

// ✅ Collection access — early-return guard, then iterate
if result.HasError() { return result }

for _, item := range result.Items() { process(item) }

// ✅ Adapter unwrap — Result[T] → (*T, error) // EXEMPTED: framework boundary adapter
func (a *Adapter) GetById(ctx context.Context, id int64) (*models.Plugin, error) {
    result := a.Service.GetById(ctx, id)

    if result.HasError() { return nil, result.AppError() }
    v := result.Value()

    return &v, nil
}
```

> **Full examples with PHP/Go/TypeScript:** see [apperror § Result Guard Rule](../../../03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/06-serialization-and-guards.md#12-result-guard-rule--mandatory-error-check-before-value-access)

#### Enforcement Checklist

- [ ] Every `result.Value()` / `$result->value()` call is preceded by `HasError()` / `hasError()` or `IsSafe()` / `isSafe()`
- [ ] Every `result.Items()` / `$results->items()` call is preceded by a guard
- [ ] Every `result.Get(key)` on `ResultMap` is preceded by a guard
- [ ] Every `$execResult->affectedRows()` on `DbExecResult` is preceded by a guard
- [ ] No error is silently discarded — all errors are logged, returned, or propagated
- [ ] Cross-service callers (direct `*service.Service` refs) guard results the same way

### Common Mistakes

```go
// ❌ MISTAKE: Raw error without stack trace
return fmt.Errorf("failed to upload: %w", err)

// ✅ CORRECT: apperror with automatic stack trace
return apperror.Wrap(err, apperror.ErrUploadFailed, "failed to upload plugin")
```

```go
// ❌ MISTAKE: Compound guard — redundant double-check after HasError()
if publishResult.HasError() || (publishResult.IsSafe() && !publishResult.Value().IsSuccess) {
    return publishResult
}

// ✅ CORRECT: HasError() is the only guard needed
if publishResult.HasError() {
    return publishResult
}
```

> **Rule:** `HasError()` is the single source of truth for failure. Never combine it with `IsSafe() && !Value().Field` conditions — that conflates transport-level errors with business-level validation, which must be handled separately at the domain layer if needed.

```php
// ❌ MISTAKE: Leading backslash on global types
catch (\Throwable $e) { ... }

// ✅ CORRECT: Use import
use Throwable;
// ...
catch (Throwable $e) { ... }
```

---
