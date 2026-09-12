# Boolean Principles — P1: is/has prefixes, P2: no negative words

> **Parent:** [Boolean Principles](./01-index.md)
> **Version:** 2.6.0
> **Updated:** 2026-03-31

---

## Principle 1: Always Use `is` or `has` Prefixes

Every boolean identifier — variable, property, parameter, or method — **must** start with `is` or `has`.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
$active = true;
$loaded = false;
$blocked = true;

// ✅ REQUIRED
$isActive = true;
$isLoaded = false;
$isBlocked = true;
$hasPermission = true;
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
const loading = true;
const valid = false;
const overdue = checkOverdue();

// ✅ REQUIRED
const isLoading = true;
const isValid = false;
const hasOverdue = checkOverdue();
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
blocked := true
connected := false

// ✅ REQUIRED
isBlocked := true
isConnected := false
hasItems := len(items) > 0
```

### Method Names Follow the Same Rule

```php
// ❌ FORBIDDEN
$order->overdue();
$user->admin();

// ✅ REQUIRED
$order->hasOverdue();
$user->isAdmin();
```

This mirrors industry best practices. For example, .NET's `char` type exposes `IsLetter`, `IsDigit`, `IsUpper`, `IsLower`, `IsNumber`, `IsPunctuation`, `IsSeparator`, `IsSymbol`, `IsControl`, `IsLetterOrDigit` — all boolean methods with the `Is` prefix.

### Function & Method Parameters: Total Ban on Single-Letter & Bare Names

Boolean parameters in function and method signatures (such as setters) MUST NEVER use lazy single-letter identifiers (`v bool`, `b bool`, `flag bool`) or bare unprefixed verbs (`stop bool`, `pause bool`). They must always carry an affirmative prefix describing the exact state being set:

```go
// ❌ FORBIDDEN: Single-letter parameter `v bool` or bare verb `stop bool`
func (p *BatchProgress) SetStopOnFail(v bool) {
    p.stopOnFail = v
}

func (w *Worker) SetStopped(stop bool) {
    w.stop = stop
}

// ✅ REQUIRED: Meaningful, affirmative boolean parameter and property
func (p *BatchProgress) SetStopOnFail(isStopOnFail bool) {
    p.stopOnFail = isStopOnFail
}

func (w *Worker) SetStopped(isStopped bool) {
    w.isStopped = isStopped
}
```

### Struct Fields & State Properties: Total Ban on Bare Names (e.g. `defined` -> `isDefined`)

Struct fields, class properties, and state flags representing boolean states MUST ALWAYS use affirmative `is*` or `has*` prefixes. A bare name such as `defined bool`, `ready bool`, `active bool` is strictly FORBIDDEN:

```go
// ❌ FORBIDDEN: Bare boolean field name in struct
type Result[T any] struct {
    value   T
    err     *AppError
    defined bool // VIOLATION: bare boolean without is/has prefix
}

// ✅ REQUIRED: Explicit affirmative boolean prefix
type Result[T any] struct {
    value     T
    err       *AppError
    isDefined bool // COMPLIANT: starts with affirmative 'is'
}
```

### Ban on Awkward `isExists` / `isUserExist` (Use `isDefined` / `isFound`)

"Exists" is a verb. Combining `is` with a verb (`isExists`, `IsExists`, `isUserExist`) is grammatically malformed and strictly banned. Always use `isDefined` (or `IsDefined`) for state or resource presence, and `isFound` for map/cache lookup presence.

---

---

## Principle 2: Never Use Negative Words in Boolean Names

The words **`not`**, **`no`**, and **`non`** are **absolutely banned** from boolean variable names, function names, and method names. These words create cognitive overhead — the reader must mentally invert the meaning. Instead, always use a **positive semantic synonym** that describes what the state actually **is**.

Double negatives (`!isNot...`, `!isNotBlocked`) are the worst form and must never appear.

### Naming Strategy: Describe What It IS, Not What It ISN'T

| ❌ Forbidden Name | ✅ Required Name | Semantic Meaning |
|---|---|---|
| `isExists` / `isUserExist` | `isDefined` / `isFound` | The resource or entity exists and is defined |
| `isNotReady` | `isPending` | The order is waiting |
| `isNotInList` | `isAbsentFromList` | The item is absent |
| `isNoRecentErrors` | `isErrorListClear` | The error list is clean |
| `isNotDirectory` | `isDirAbsent` | The directory doesn't exist |
| `isNotRegularFile` | `isIrregularPath` | The path is irregular |
| `isNotPHP` | `isSkippableEntry` | The entry should be skipped |
| `isNotBlocked` | `isActive` | The entity is active |
| `isClassNotLoaded` | `isClassUnregistered` | The class is unregistered |
| `hasNoPermission` | `isUnauthorized` | The user lacks access |
| `isNotHonest` | `isDishonest` | Use the inverse naming instead of 'not' |

```typescript
// ❌ FORBIDDEN — "not" in the variable name
const isNotReady = order.status !== 'ready';
if (isNotReady) {
    throw new Error('Order is not ready');
}

// ✅ REQUIRED — Positive semantic synonym
const isPending = order.status !== 'ready';
if (isPending) {
    throw new Error('Order is not ready');
}
```

```php
// ❌ FORBIDDEN — "No" in the variable name
$isNoRecentErrors = empty($errors) || !$hasUnseen;

// ✅ REQUIRED — Describes the positive state
$isErrorListClear = empty($errors) || !$hasUnseen;
```

### Rule: Name booleans for the **positive semantic state**, then negate only once in guard if needed

```typescript
// ❌ FORBIDDEN: Negative boolean names
const hasNoColors = !colorConfig.length;
if (hasNoColors) {
    return null;
}

const hasNoPayload = !payload?.length;
if (hideLabel || hasNoPayload) {
    return null;
}

// ✅ REQUIRED: Positive framing + inverted guard condition
const hasColors = colorConfig.length > 0;
if (!hasColors) {
    return null;
}

const hasPayload = Boolean(payload?.length);
if (hideLabel || !hasPayload) {
    return null;
}
```

---

## Principle 3: Total Ban on Bare `ok` Identifiers

In Go type assertions, map lookups, and channel receives, the bare identifier `ok` is **strictly forbidden**. It violates the boolean prefix rule and obscures the domain meaning. Always use a semantic affirmative boolean starting with `is` or `has`:

```go
// ❌ FORBIDDEN: Bare ok identifier in type assertions
if appErr, ok := err.(*apperror.AppError); ok {
    if appErr.Code != "E_INTERNAL_ERROR" {
        t.Errorf("expected E_INTERNAL_ERROR, got %s", appErr.Code)
    }
} else {
    t.Errorf("expected AppError, got %T", err)
}

// ✅ REQUIRED: Semantic isAppErr + inverted guard clause
appErr, isAppErr := err.(*apperror.AppError)
if !isAppErr {
    t.Fatalf("expected AppError, got %T", err)
}

if appErr.Code != "E_INTERNAL_ERROR" {
    t.Errorf("expected E_INTERNAL_ERROR, got %s", appErr.Code)
}
```

```go
// ❌ FORBIDDEN: Bare ok in map lookups
val, ok := userCache[id]
if !ok {
    return nil, ErrUserNotFound
}

// ✅ REQUIRED: Semantic affirmative boolean
val, isFound := userCache[id]
if !isFound {
    return nil, ErrUserNotFound
}
```
