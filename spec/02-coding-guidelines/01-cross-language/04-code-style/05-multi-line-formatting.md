# Multi-Line Formatting

> **Version:** 4.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** PHP, TypeScript, Go  
> **Rules covered:** 9, 11, apperror multi-line

---

## Rule 9: Multi-Line Arguments — Signatures, Calls, and Arrays

When a function/method **signature or call** has **more than two arguments**, each argument must be on its own line with consistent indentation and a **trailing comma** after the last argument (where syntax permits).

This applies equally to:
- **Function/method signatures** (parameter declarations)
- **Function/method calls** (argument expressions)
- **Constructor calls** (`new Foo(...)`)

### 9a: Function Signatures (>2 Parameters)

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN (>2 params on one line)
function buildRecord(string $label, string $path, bool $success, ?string $error): void {

// ✅ REQUIRED
function buildRecord(
    string $label,
    string $path,
    bool $success,
    ?string $error,
): void {

// ✅ OK: 2 params — single line is fine
function loadFile(string $label, string $path): bool {
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN (>2 params on one line)
function buildRecord(label: string, path: string, success: boolean, error?: string): void {

// ✅ REQUIRED
function buildRecord(
    label: string,
    path: string,
    success: boolean,
    error?: string,
): void {
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN (>2 params on one line)
func BuildRecord(label string, path string, success bool, errMsg string) {

// ✅ REQUIRED
func BuildRecord(
	label string,
	path string,
	success bool,
	errMsg string,
) {
```

### 9b: Function Calls (>2 Arguments)

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN (>2 args on one line)
$this->logAction($agentId, ActionType::AgentTest->value, null, StatusType::Failed->value, null, $error->get_error_message());

// ✅ REQUIRED
$this->logAction(
    $agentId,
    ActionType::AgentTest->value,
    null,
    StatusType::Failed->value,
    null,
    $error->get_error_message(),
);

// ✅ OK: 2 args — single line is fine
$this->updateAgent($agentId, $data);
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN (>2 args on one line)
const result = buildRecord(label, path, true, errorMessage);

// ✅ REQUIRED
const result = buildRecord(
    label,
    path,
    true,
    errorMessage,
);

// ✅ OK: 2 args — single line is fine
const result = fetchData(url, options);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN (>2 args on one line)
result := buildRecord(label, path, true, errMsg)

// ✅ REQUIRED
result := buildRecord(
	label,
	path,
	true,
	errMsg,
)
```

### 9c: PHP Arrays — Each Item on Its Own Line

In PHP, `array(...)` and `[...]` literals with **more than two items** must place each item on its own line with a trailing comma.

```php
// ❌ FORBIDDEN (>2 items on one line)
$statuses = array(301, 302, 303, 307, 308);
$data = ['agentId' => $agentId, 'action' => $action, 'slug' => $slug];

// ✅ REQUIRED
$statuses = array(
    301,
    302,
    303,
    307,
    308,
);

$data = [
    'agentId' => $agentId,
    'action'  => $action,
    'slug'    => $slug,
];

// ✅ OK: 2 items — single line is fine
$pair = array('key' => $value, 'name' => $name);
```

```typescript
// ── TypeScript / Go ─────────────────────────────────────────
// Same principle applies to array/slice literals with >2 items.
// Each item on its own line with trailing comma.

// ❌ FORBIDDEN
const codes = [301, 302, 303, 307, 308];

// ✅ REQUIRED
const codes = [
    301,
    302,
    303,
    307,
    308,
];
```

---

## Rule 11: Method Chaining — Each Call on Its Own Line

When using **method chaining** (fluent API), every chained `.Method()` call **must** be on its own line, indented one level. Never chain multiple methods on a single line.

This applies to all fluent APIs including error wrapping (`apperror.Wrap`), query builders, HTTP clients, and any builder pattern.

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN — Multiple chained calls on one line
return apperror.Wrap(err, apperror.ErrWPConnection, "failed to stream snapshot ZIP").WithSiteId(siteId).WithSnapshotId(snapshotId).WithUrl(meta.Url)

// ✅ REQUIRED — Each chained call on its own line
return apperror.Wrap(
    err, apperror.ErrWPConnection, "failed to stream snapshot ZIP",
).
    WithSiteId(siteId).
    WithSnapshotId(snapshotId).
    WithUrl(meta.Url)
```

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN — Chained calls on one line
$query->select('*')->from('users')->where('isActive', true)->orderBy('name')->limit(10);

// ✅ REQUIRED — Each call on its own line
$query
    ->select('*')
    ->from('users')
    ->where('isActive', true)
    ->orderBy('name')
    ->limit(10);
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN — Chained calls on one line
const result = queryBuilder.select('*').from('users').where({ isActive: true }).orderBy('name').execute();

// ✅ REQUIRED — Each call on its own line
const result = queryBuilder
    .select('*')
    .from('users')
    .where({ isActive: true })
    .orderBy('name')
    .execute();
```

### Exception

Chains with **2 or fewer** method calls may remain on a single line if the total line length stays under 100 characters:

```go
// ✅ OK — single chained method, short line
return apperror.Wrap(
    err,
    code,
    msg,
).WithPath(path)
```

---

## Rule: Universal Multi-Line Formatting for `apperror` Calls

All `apperror.Wrap`, `apperror.New`, `apperror.FailWrap`, `apperror.FailNew`, and `apperror.WrapNew` calls with **2 or more arguments** must use multi-line formatting — **regardless of line length**. This is a universal rule with no exceptions.

### `apperror.New` — Always Multi-Line

```go
// ❌ FORBIDDEN — inline multi-arg
return apperror.New(apperror.ErrValidationEmpty, "site name is required")

// ✅ REQUIRED — each arg on its own line
return apperror.New(
    apperror.ErrValidationEmpty,
    "site name is required",
)
```

### `apperror.Wrap` — Always Multi-Line

```go
// ❌ FORBIDDEN — inline multi-arg
return apperror.Wrap(err, apperror.ErrDatabaseExec, "failed to update")

// ✅ REQUIRED — each arg on its own line
return apperror.Wrap(
    err,
    apperror.ErrDatabaseExec,
    "failed to update",
)
```

### Nested `apperror.Fail[T](apperror.Wrap(...))` — Always Multi-Line

```go
// ❌ FORBIDDEN — nested call on single line
return apperror.Fail[*http.Response](apperror.Wrap(errors.ErrHTTPConnection, "create request", err))

// ✅ REQUIRED — each level expanded
return apperror.Fail[*http.Response](
    apperror.Wrap(
        errors.ErrHTTPConnection,
        "create request",
        err,
    ),
)
```

### With Method Chaining

When combined with `.WithContext()` or `.WithPath()` chaining, the base call is multi-line and each chain method gets its own line:

```go
// ✅ REQUIRED — multi-line base + chained methods
return apperror.Wrap(
    err,
    apperror.ErrWPConnect,
    "failed to stream snapshot",
).
    WithSiteId(siteId).
    WithSnapshotId(snapshotId).
    WithUrl(meta.Url)
```

### Single-Arg Calls

Calls with exactly **one argument** may remain inline:

```go
// ✅ OK — single arg
return apperror.Fail[Plugin](appErr)
```

---

*Part of [Code Style](./00-overview.md) — Rules 9, 11, apperror multi-line*
