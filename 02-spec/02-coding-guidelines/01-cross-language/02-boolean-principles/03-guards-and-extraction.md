# Boolean Principles — P3: named guards, P4: extract complex expressions (AI Execution Prompt)

> **/goal** Eliminate raw call-site negations by utilizing affirmative named guards and decompose complex boolean expressions exceeding two operands or mixing operators into dedicated intermediate variables.
> **/learn** Master semantic inverse helpers (e.g., `isInvalid()`, `isFileMissing()`), the 2-operand rule (P4a), strict isolation between `&&` and `||` (P4b), and separation of negative checks from positive conditions (P4c).

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace raw negations on function calls (`!fn()`) with semantic named guard methods or utilities (`fnInvalid()`, `isMissing()`).
- [ ] `/learn` Cap chained boolean operands at a maximum of two per expression; break 3+ conditions into intermediate variables.
- [ ] `/goal` Never mix `&&` and `||` operators in a single conditional statement; decompose into separate named booleans.
- [ ] `/learn` Never combine negative checks (`err != nil`, `isNil`) and positive checks in the same expression; isolate guards first.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Boolean Principles](./readme.md)
> **Version:** 2.6.0
> **Updated:** 2026-10-02

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

See also: [code-style.md — Rule 3](../04-code-style/03-conditions-and-extraction.md#rule-3-extract-complex-conditions--no-inline-multi-part-checks)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-BOOL-003: Named Guard Inversion and Complex Expression Extraction

**Given** Complex conditionals and raw boolean negations across application source code.
**When** Codebases are analyzed for logical operator density, raw negation patterns, and operator mixing.
**Then** Raw negations are replaced by affirmative named guards, all conditional expressions contain at most two operands without mixing `&&` and `||`, and mixed-polarity expressions are eliminated.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles --check-only
```
**Expected:** exit 0. Zero violations.
