# Boolean Principles — P3: named guards, P4: extract complex expressions

> **Parent:** [Boolean Principles](./00-overview.md)  
> **Version:** 2.6.0  
> **Updated:** 2026-03-31

---

## Principle 3: Replace Raw Negation With Named Guards

Never use raw `!` on function calls or existence checks at call sites. Instead, wrap every negative check in a **positively named utility function**.

```php
// ❌ FORBIDDEN — Raw negation on function call
if (!$order->isValid()) {
    return;
}

// ✅ REQUIRED — Semantic inverse method on the object
if ($order->isInvalid()) {
    return;
}
```

```typescript
// ❌ FORBIDDEN
if (!isDefined(value)) {
    return;
}

// ✅ REQUIRED — Use a positive guard
if (isUndefined(value)) {
    return;
}
```

```go
// ❌ FORBIDDEN
if !IsFileExists(path) {
    return apperror.New(
        "E4010", "file not found",
    )
}

// ✅ REQUIRED
if IsFileMissing(path) {
    return apperror.New(
        "E4010", "file not found",
    )
}
```

For the full guard function inventory, see [no-negatives.md](../12-no-negatives.md).

---


---

## Principle 4: Extract Complex Boolean Expressions

When a boolean expression contains **2+ operators** (`&&`, `||`, `!`), it **must** be extracted into a named boolean variable or a dedicated method. The `if` statement should read as a single intent.

### P4a — Maximum Two Conditions Per Expression

A single boolean expression may combine **at most two** operands with `&&` or `||`. Three or more operands **must** be decomposed into intermediate named booleans.

```php
// ❌ FORBIDDEN: Three conditions chained together
$hasFileParam = $request !== null
    && $request->hasParam('file')
    && $request->getParam('file') !== '';

// ✅ REQUIRED: Decompose into two-condition steps
$hasRequest = $request !== null;
$hasNonEmptyFile = $hasRequest
    && $request->hasParam('file')
    && $request->getParam('file') !== '';

// ✅ BETTER: Early return to eliminate null guard, then two-condition max
if ($request === null) {
    return;
}

$hasNonEmptyFile = $request->hasParam('file')
    && $request->getParam('file') !== '';

if ($hasNonEmptyFile) {
    $this->process($request);
}
```

```go
// ❌ FORBIDDEN: Three conditions chained
isUpstreamError := err != nil && resp != nil && resp.StatusCode >= 400

// ✅ REQUIRED: Decompose — early return for error, then check response
if err != nil {
    return apperror.Wrap(err, "E5001", "upstream call failed")
}

isUpstreamError := resp != nil && resp.StatusCode >= 400

if isUpstreamError {
    handleUpstreamError(resp)
}
```

```typescript
// ❌ FORBIDDEN: Three conditions chained
const isDelegatedError = response != null
    && response.status >= 400
    && response.data?.code?.startsWith('E8');

// ✅ REQUIRED: Decompose into two-condition steps
const isErrorResponse = response != null && response.status >= 400;
const isDelegatedError = isErrorResponse && response.data?.code?.startsWith('E8');

if (isDelegatedError) {
    showDelegatedError(response);
}
```

### P4b — Never Mix `&&` with `||` in One Expression

A single boolean expression must use **only one** logical operator type. Mixing `&&` and `||` creates ambiguity and bugs.

```csharp
// ❌ FORBIDDEN: Mixed && and || — ambiguous precedence
if (value > 0 && value % 2 == 0 || value < -10)
{
    // ...
}

// ✅ REQUIRED: Separate into named booleans, one operator each
bool isPositiveEven = value > 0 && value % 2 == 0;
bool isBelowThreshold = value < -10;
bool isValueValid = isPositiveEven || isBelowThreshold;

if (isValueValid)
{
    // ...
}
```

```typescript
// ❌ FORBIDDEN: Mixed && and ||
if (isAdmin && hasPermission || isSuperUser) {
    grantAccess();
}

// ✅ REQUIRED: Decompose
const hasAdminAccess = isAdmin && hasPermission;
const canAccess = hasAdminAccess || isSuperUser;

if (canAccess) {
    grantAccess();
}
```

### P4c — Never Mix Negative and Positive Checks in One Expression

A single boolean expression must not combine a negative check (`!`, `null`, `=== false`) with a positive check. Separate them.

```php
// ❌ FORBIDDEN: Negative (null check) mixed with positive
$hasFileParam = $request !== null && $request->hasParam('file');

// ✅ REQUIRED: Separate the null guard, then check positive
if ($request === null) {
    return;
}

$hasFileParam = $request->hasParam('file');
```

```go
// ❌ FORBIDDEN: Negative (nil check) mixed with positive (status check)
isValid := err == nil && resp.StatusCode == 200

// ✅ REQUIRED: Early return for negative, then positive check
if err != nil {
    return apperror.Wrap(err, "E5001", "request failed")
}

isValid := resp.StatusCode == 200
```

```typescript
// ❌ FORBIDDEN: Negative mixed with positive
const isReady = response != null && response.status === 200;

// ✅ REQUIRED: Guard the negative case first
if (isNullish(response)) {
    return;
}

const isReady = response.status === 200;
```

See also: [code-style.md — Rule 3](../04-code-style/02-conditions-and-extraction.md#rule-3-extract-complex-conditions--no-inline-multi-part-checks)

---
