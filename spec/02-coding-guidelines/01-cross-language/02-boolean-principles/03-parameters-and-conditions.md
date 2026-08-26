# Boolean Principles — P5: explicit params, P6: no mixed booleans, P7: no inline statements, P8: no raw system calls

> **Parent:** [Boolean Principles](./00-overview.md)  
> **Version:** 2.6.0  
> **Updated:** 2026-03-31

---

## Principle 5: Boolean Parameters Must Be Explicit

Never use bare `true`/`false` at call sites. If a function accepts a boolean parameter, either:
1. Use separate, explicitly named methods
2. Use an enum or options object

```typescript
// ❌ FORBIDDEN — What does `true` mean here?
fetchData(userId, true);

// ✅ REQUIRED — Option A: Named methods
fetchDataWithCache(userId);
fetchDataWithoutCache(userId);

// ✅ REQUIRED — Option B: Options object
fetchData(userId, { isUseCache: true });
```

```php
// ❌ FORBIDDEN
$this->log($message, true);

// ✅ REQUIRED — Separate methods
$this->logWithTrace($message);
$this->log($message);
```

See also: [function-naming.md](../10-function-naming.md)

---


---

## Principle 6: Never Mix Positive and Negative Booleans in a Single Condition

Combining a positive boolean with a negated boolean in the same `if` condition (e.g., `isX && !y`, `IsReady && !overwrite`, `cacheEnabled && !forceRefresh`) is a **code smell**. It forces the reader to mentally switch polarity mid-expression, creating cognitive load and hiding intent.

**The fix:** Extract the combined condition into a single, positively named boolean that captures the **actual intent**.

```go
// ❌ FORBIDDEN — Mixed polarity: positive + negative
if s.isCacheEnabled && !isForceRefresh {
    return cachedResult
}

// ✅ REQUIRED — Extract negation to positive counterpart, then compose
isNormalRefresh := !isForceRefresh
isCacheHit := s.isCacheEnabled && isNormalRefresh

if isCacheHit {
    return cachedResult
}
```

```go
// ❌ FORBIDDEN — Mixed polarity: positive + negative
if isProjectExists && !isOverwrite {
    return fmt.Errorf("conflict")
}

// ✅ REQUIRED — Extract negation to positive counterpart, then compose
isReadOnly := !isOverwrite
isConflict := isProjectExists && isReadOnly

if isConflict {
    return fmt.Errorf("conflict")
}
```

```php
// ❌ FORBIDDEN — Mixed polarity
if ($isAuthenticated && !$isAuthorized) {
    throw new ForbiddenException();
}

// ✅ REQUIRED — Extract negation to positive counterpart, then compose
$isUnauthorized = !$isAuthorized;
$isAccessDenied = $isAuthenticated && $isUnauthorized;

if ($isAccessDenied) {
    throw new ForbiddenException();
}
```

```typescript
// ❌ FORBIDDEN — Mixed polarity
if (isLoggedIn && !hasPermission) {
    redirect('/unauthorized');
}

// ✅ REQUIRED — Single intent
const isUnauthorized = isLoggedIn && !hasPermission;

if (isUnauthorized) {
    redirect('/unauthorized');
}
```

### Why This Matters

| Pattern | Problem | Fix |
|---|---|---|
| `isX && !y` | Reader must switch polarity mid-expression | Extract to `isConflict` / `isUnauthorized` / `isDenied` |
| `isCacheEnabled && !forceRefresh` | Mixed polarity + missing `is` prefix on `forceRefresh` | Use `isCacheHit` |
| `isReady && !isOverwrite` | `!isOverwrite` lacks semantic meaning | Use `isFreshImport` or extract full condition |
| `hasData && !isProcessed` | Two separate concerns crammed together | Extract to `isPendingProcessing` |

### Rule Summary

1. **Never combine `isX` with `!isY`** in the same `if` condition
2. **Always extract** the combined condition into a named boolean with a positive semantic name
3. The named boolean should express the **intent** (e.g., `isConflict`, `isAccessDenied`, `isPending`, `isCacheHit`) — not just restate the logic

---


---

## Principle 7: No Inline Statements in Conditions

Inline variable assignment inside `if`/`while` conditions is **prohibited** across all languages. This pattern hides computation inside control flow, couples unrelated operations, and makes debugging harder.

**Exception (Go only):** The idiomatic comma-ok pattern (`if v, ok := m[k]; ok {`) and type assertions are exempt.

```go
// ❌ FORBIDDEN — inline os.Stat inside if
if _, err := os.Stat(dir); err == nil {
    // exists
}

// ✅ REQUIRED — separate computation
isProjectExists := pathutil.IsDir(dir)
if isProjectExists {
    // exists
}
```

```php
// ❌ FORBIDDEN — inline assignment in condition
if (($result = $db->query($sql)) && $result->hasRows()) {
    process($result);
}

// ✅ REQUIRED — separate assignment
$result = $db->query($sql);
$hasRows = $result !== null && $result->hasRows();

if ($hasRows) {
    process($result);
}
```

```typescript
// ❌ FORBIDDEN — assignment expression in condition
let match;
if ((match = regex.exec(input)) !== null) {
    process(match);
}

// ✅ REQUIRED — separate assignment
const match = regex.exec(input);
const isMatchFound = match !== null;

if (isMatchFound) {
    process(match);
}
```

See [Go Boolean Standards P7](../../03-golang/02-boolean-standards.md#28--no-inline-statements-in-if-conditions-rule-p7) for Go-specific details.

---


---

## Principle 8: No Raw Filesystem / System Calls in Application Code

Application code **must not** call raw filesystem functions (`os.Stat`, `file_exists`, `fs.existsSync`) directly. Instead, use wrapper utilities that:

1. Return the framework's structured error type (not raw errors)
2. Provide positive-named boolean helpers (`IsDir`, `IsDirMissing`, `isFileMissing`)
3. Handle edge cases consistently

```go
// ❌ FORBIDDEN — raw os.Stat
if _, err := os.Stat(projectDir); err == nil {
    // exists
}

// ✅ REQUIRED — pathutil wrapper
isProjectExists := pathutil.IsDir(projectDir)
```

```php
// ❌ FORBIDDEN — raw file_exists
if (file_exists($path)) {
    // exists
}

// ✅ REQUIRED — PathHelper wrapper
if (PathHelper::isFilePresent($path)) {
    // exists
}
```

```typescript
// ❌ FORBIDDEN — raw fs.existsSync
if (fs.existsSync(path)) {
    // exists
}

// ✅ REQUIRED — pathUtil wrapper
const isFilePresent = pathUtil.isFile(path);
if (isFilePresent) {
    // exists
}
```

### Comprehensive Example — All Rules Combined (P6 + P7 + P8)

This example demonstrates how P6, P7, and P8 violations compound, and the correct fix:

```go
// ❌ FORBIDDEN — 4 violations in 3 lines:
//   1. P8: raw os.Stat (must use pathutil)
//   2. P7: inline statement in if (semicolon)
//   3. P6: mixed polarity (err == nil && !isOverwrite)
//   4. Raw fmt.Errorf instead of apperror
isProjectConflict := err == nil && !isOverwrite
if _, err := os.Stat(projectDir); isProjectConflict {
    return fmt.Errorf("project exists, use isOverwrite=true to replace")
}

// ✅ REQUIRED — clean, readable, all rules applied:
//   1. P8: pathutil.IsDir() wraps os.Stat
//   2. P7: no inline statement; all variables computed before if
//   3. P6: mixed polarity extracted to single-intent boolean
//   4. apperror.FailNew returns structured *apperror.AppError
isProjectExists := pathutil.IsDir(projectDir)
isReadOnly := !isOverwrite
isProjectConflict := isProjectExists && isReadOnly

if isProjectConflict {
    return apperror.FailNew[ProjectResult](
        errors.ErrFsConflict,
        "project already exists; use isOverwrite=true to replace",
    )
}
```

---

## Principle 9: Strict Conditional Joins

Never mix logical operators (`&&` and `||`) in the same `if` condition. Furthermore, keep it to a single conditional join (max two operands).
If you need more than one join, you must extract the logic into well-named boolean variables.

```go
// ❌ FORBIDDEN - Multiple conditional joins
if isUserActive && hasSubscription && isBillingCurrent { ... }

// ❌ FORBIDDEN - Mixed operators
if isUserActive && (hasSubscription || isTrial) { ... }

// ✅ REQUIRED - Extracted named booleans
canAccessContent := hasSubscription || isTrial
isEligible := isUserActive && canAccessContent
if isEligible { ... }
```

## Principle 10: No Mixed Polarity

Never mix positive and negative conditions in a single conditional join. Both sides of the join must be positive.

```typescript
// ❌ FORBIDDEN - Mixed polarity
if (isReady && !isExpired) { ... }

// ✅ REQUIRED - All positive
if (isReady && isValid) { ... }
```
