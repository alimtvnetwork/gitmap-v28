# Cross-Language Rule: No Raw Negations — Use Positive Guard Functions

> **Version:** 2.1.0  
> **Updated:** 2026-03-31  
> **Applies to:** PHP, TypeScript, Go

---

## Principle

**Never use raw negation operators (`!`, `not`) on function calls or existence checks in conditions.** Instead, wrap every negative check in a **positively named utility function** that reads as a single intent.

Raw negations are easy to miss during code review, cause cognitive overhead, and scatter low-level logic across call sites. A named guard function centralizes the check, is self-documenting, and eliminates the visual noise of `!`.

---

## The Rule

| ❌ Forbidden | ✅ Required | Why |
|-------------|------------|-----|
| `!file_exists($path)` | `PathHelper::isFileMissing($path)` | Positive name; no `!` to overlook |
| `!is_dir($path)` | `PathHelper::isDirMissing($path)` | Self-documenting intent |
| `!class_exists('X')` | `BooleanHelpers::isClassMissing('X')` | Centralized, testable |
| `!function_exists('f')` | `BooleanHelpers::isFuncMissing('f')` | Same principle |
| `!extension_loaded('e')` | `BooleanHelpers::isExtensionMissing('e')` | Same principle |
| `!$obj->isActive()` | `$obj->isDisabled()` | Semantic inverse on object |
| `!arr.includes(x)` | `isMissing(arr, x)` | Named guard |
| `!strings.Contains(s, x)` | `IsMissing(s, x)` | Named guard |

### Key: Every negative check becomes a **positively named function**

The function name must express the **positive assertion** of what is being checked:
- "is missing" not "is not existing"
- "is disabled" not "is not active"  
- "is empty" not "is not filled"
- "is disconnected" not "is not connected"

---

## Language-Specific Examples

### PHP (camelCase methods)

```php
// ❌ FORBIDDEN: Raw negation on function call
if (!file_exists($path)) {
    return false;
}

if (!is_dir($dir)) {
    mkdir($dir, 0755, true);
}

if (!class_exists('PDO')) {
    throw new RuntimeException('PDO not available');
}

// ✅ REQUIRED: Positive guard function from PathHelper / BooleanHelpers
use RiseupAsia\Helpers\PathHelper;
use RiseupAsia\Helpers\BooleanHelpers;

if (PathHelper::isFileMissing($path)) {
    return false;
}

if (PathHelper::isDirMissing($dir)) {
    mkdir($dir, 0755, true);
}

if (BooleanHelpers::isClassMissing('PDO')) {
    throw new RuntimeException('PDO not available');
}
```

**Utility classes:**
- `PathHelper` (`RiseupAsia\Helpers\PathHelper`) — file/directory guards
- `BooleanHelpers` (`RiseupAsia\Helpers\BooleanHelpers`) — function/class/extension/database guards

| Guard Method | Replaces | Class |
|-------------|----------|-------|
| `isFileMissing($path)` | `!file_exists($path)` | `PathHelper` |
| `isFileExists($path)` | `file_exists($path)` (with null guard) | `PathHelper` |
| `isDirMissing($path)` | `!is_dir($path)` | `PathHelper` |
| `isDirExists($path)` | `is_dir($path)` (with null guard) | `PathHelper` |
| `isDirWritable($path)` | `is_dir($path) && is_writable($path)` | `PathHelper` |
| `isDirReadonly($path)` | `!is_dir($path) \|\| !is_writable($path)` | `PathHelper` |
| `isClassMissing($name)` | `!class_exists($name)` | `BooleanHelpers` |
| `isClassExists($name)` | `class_exists($name)` | `BooleanHelpers` |
| `isFuncMissing($name)` | `!function_exists($name)` | `BooleanHelpers` |
| `isFuncExists($name)` | `function_exists($name)` | `BooleanHelpers` |
| `isExtensionMissing($name)` | `!extension_loaded($name)` | `BooleanHelpers` |
| `isExtensionLoaded($name)` | `extension_loaded($name)` | `BooleanHelpers` |
| `isDbConnected($db)` | `$db !== null && $db->isConnected()` | `BooleanHelpers` |
| `isDbDisconnected($db)` | `$db === null \|\| !$db->isConnected()` | `BooleanHelpers` |

### TypeScript (camelCase functions)

```typescript
// ❌ FORBIDDEN: Raw negation
if (!fs.existsSync(path)) {
    throw new Error('File not found');
}

if (!response.ok) {
    handleError(response);
}

if (!array.includes(item)) {
    array.push(item);
}

// ✅ REQUIRED: Positive guard function
if (isFileMissing(path)) {
    throw new Error('File not found');
}

if (isResponseFailed(response)) {
    handleError(response);
}

if (isItemMissing(array, item)) {
    array.push(item);
}
```

**Utility location:** `src/utils/guards.ts` or domain-specific guard files

| Guard Function | Replaces |
|---------------|----------|
| `isFileMissing(path)` | `!fs.existsSync(path)` |
| `isFileExists(path)` | `fs.existsSync(path)` |
| `isResponseFailed(res)` | `!res.ok` |
| `isResponseSuccess(res)` | `res.ok` |
| `isArrayEmpty(arr)` | `!arr.length` or `arr.length === 0` |
| `hasItems(arr)` | `arr.length > 0` |
| `isNullish(val)` | `val == null` |
| `isPresent(val)` | `val != null` |
| `isStringEmpty(str)` | `!str` or `str === ''` |
| `hasContent(str)` | `!!str` or `str.length > 0` |

### Go (PascalCase exported functions)

```go
// ❌ FORBIDDEN: Raw negation
if !fileExists(path) {
    return fmt.Errorf("file not found: %s", path)
}

if !strings.Contains(s, substr) {
    return apperror.New(
        "E4010", "missing required field",
    )
}

// ✅ REQUIRED: Positive guard function
if IsFileMissing(path) {
    return apperror.New(
        "E4010", "file not found: "+path,
    )
}

if IsMissingSubstring(s, substr) {
    return apperror.New(
        "E4010", "missing required field",
    )
}
```

**Utility package:** `pkg/guards/` or `internal/guards/`

| Guard Function | Replaces |
|---------------|----------|
| `IsFileMissing(path)` | `!fileExists(path)` |
| `IsFileExists(path)` | `fileExists(path)` |
| `IsDirMissing(path)` | `!dirExists(path)` |
| `IsDirExists(path)` | `dirExists(path)` |
| `IsStringEmpty(s)` | `s == ""` or `len(s) == 0` |
| `HasContent(s)` | `s != ""` or `len(s) > 0` |
| `IsSliceEmpty(s)` | `len(s) == 0` |
| `HasItems(s)` | `len(s) > 0` |
| `IsMissingSubstring(s, sub)` | `!strings.Contains(s, sub)` |
| `ContainsSubstring(s, sub)` | `strings.Contains(s, sub)` |
| `IsDefined(x)` | `x != nil` (positive existence check) |
| `IsDefinedAndValid(x)` | `x != nil && x.IsValid()` (existence + validation) |
| `IsEmpty(x)` | `x == nil` (absence check) |

---

## Object-Level Semantic Inverses

Every boolean method on an object **must have a semantic inverse** — never negate a method call with `!`.

### Core Pairs

```php
// ❌ FORBIDDEN              → ✅ REQUIRED
if (!$plugin->isActive())    → if ($plugin->isInactive())
if (!$user->hasPermission()) → if ($user->lacksPermission())
if (!$account->isBlocked())  → if ($account->isUnblocked())
if (!$task->isComplete())    → if ($task->isIncomplete())
```

### Common Semantic Inverse Pairs

| Positive Method | Semantic Inverse | ❌ Never Use |
|----------------|-----------------|-------------|
| `isActive()` | `isInactive()` | `!isActive()` |
| `isEnabled()` | `isDisabled()` | `!isEnabled()` |
| `isBlocked()` | `isUnblocked()` | `!isBlocked()` |
| `isLocked()` | `isUnlocked()` | `!isLocked()` |
| `isConnected()` | `isDisconnected()` | `!isConnected()` |
| `isVisible()` | `isHidden()` | `!isVisible()` |
| `isOpen()` | `isClosed()` | `!isOpen()` |
| `isComplete()` | `isIncomplete()` | `!isComplete()` |
| `isValid()` | `isInvalid()` | `!isValid()` |
| `isAuthorized()` | `isUnauthorized()` | `!isAuthorized()` |
| `hasPermission()` | `lacksPermission()` | `!hasPermission()` |
| `hasAccess()` | `lacksAccess()` | `!hasAccess()` |
| `isSubscribed()` | `isUnsubscribed()` | `!isSubscribed()` |

### State Transition Pairs (Action Methods)

These are **action methods** (not boolean checks) but follow the same naming philosophy:

| Action | Inverse Action |
|--------|---------------|
| `block()` | `unblock()` |
| `activate()` | `deactivate()` |
| `lock()` | `unlock()` |
| `enable()` | `disable()` |
| `subscribe()` | `unsubscribe()` |
| `connect()` | `disconnect()` |
| `show()` | `hide()` |

```typescript
// ❌ FORBIDDEN
if (!plugin.isActive()) { ... }
if (!user.hasPermission('admin')) { ... }
if (!account.isBlocked()) { ... }

// ✅ REQUIRED
if (plugin.isInactive()) { ... }
if (user.lacksPermission('admin')) { ... }
if (account.isUnblocked()) { ... }
```

```go
// ❌ FORBIDDEN
if !plugin.IsActive() { ... }
if !user.HasPermission("admin") { ... }
if !session.IsLocked() { ... }

// ✅ REQUIRED
if plugin.IsInactive() { ... }
if user.LacksPermission("admin") { ... }
if session.IsUnlocked() { ... }
```

### Database-Backed Inverses (Rule 9)

When the underlying boolean lives in a **database column**, the same
inverse-pair philosophy applies — but with a hard storage rule: the DB
stores **only the positive (canonical) form**, and the inverted accessor
is **auto-generated as a computed method/getter** in the model, never as
a second column.

| Stored DB column | Auto-derived code accessor | Derivation |
|------------------|----------------------------|-----------|
| `IsActive`       | `IsInactive()` / `isInactive()` | `!IsActive` |
| `IsEnabled`      | `IsDisabled()`              | `!IsEnabled` |
| `HasLicense`     | `HasNoLicense()`            | `!HasLicense` |

> **Spec:** [Database Naming Conventions — Rule 9: Auto-Generated Inverted (Computed) Fields](../../04-database-conventions/01-naming-conventions.md#rule-9-auto-generated-inverted-computed-fields-in-code)
>
> **Codegen:** [`linters-cicd/codegen/`](../../../linters-cicd/codegen/readme.md) emits the inverse methods/traits/getters for Go, PHP, and TypeScript automatically.
>
> **Linter:** `BOOL-NEG-001` blocks `Not`/`No`-prefixed column names from being introduced via migration.

---

## When Raw `!` Is Still Acceptable


Raw negation is **only** acceptable for:

1. **Simple boolean variable checks** where the variable is already a positively named `is_*`/`has_*` boolean:
   ```php

   if (!$isInitialized) { ... }  // ✅ OK — variable is already semantic
   ```

2. **Logical operators in extracted named booleans** (inside the variable/method definition, not at the call site):
   ```php
   $isInvalid = !$isValid && !$hasOverride;  // ✅ OK — inside named boolean

   if ($isInvalid) { ... }                    // ✅ Call site is clean
   ```

3. **Native type coercion** where no function exists:
   ```php

   if (!$value) { ... }  // ✅ OK — simple falsy check on primitive
   ```

---

## Checklist Summary (Copy for PRs)

```
[ ] No `!file_exists()` — use `PathHelper::isFileMissing()`
[ ] No `!is_dir()` — use `PathHelper::isDirMissing()`
[ ] No `!class_exists()` — use `BooleanHelpers::isClassMissing()` / guard function
[ ] No `!function_exists()` — use `BooleanHelpers::isFuncMissing()` / guard function
[ ] No `!extension_loaded()` — use `BooleanHelpers::isExtensionMissing()` / guard function
[ ] No `!$obj->isActive()` — use `$obj->isDisabled()` / semantic inverse
[ ] No `!array.includes()` — use `isItemMissing()` / guard function
[ ] No `!strings.Contains()` — use `IsMissingSubstring()` / guard function
[ ] Guard functions live in dedicated utility classes/packages
[ ] Every boolean method on objects has a semantic inverse
[ ] No mixed-polarity conditions (`isX && !y`) — extract to a named boolean (see P6 in boolean-principles.md)
```

---

## Cross-References

- [Boolean Principles Overview](./02-boolean-principles/00-overview.md) — Is/Has prefix rules and parent index
- [Database Naming — Rule 9 (Inverted Fields)](../../04-database-conventions/01-naming-conventions.md#rule-9-auto-generated-inverted-computed-fields-in-code) — DB-side inverse contract + codegen
- [PHP Boolean Logic](../04-php/07-php-standards-reference/03-initialization-and-booleans.md#boolean-logic) — PHP-specific helpers
- [PHP Forbidden Patterns](../04-php/02-forbidden-patterns.md) — Pattern 4.x
- [Cross-Language Code Style](./04-code-style/00-overview.md) — Braces, nesting, spacing
- [TypeScript Standards](../02-typescript/08-typescript-standards-reference.md)
- [Golang Standards](../03-golang/04-golang-standards-reference/00-overview.md)

---

*No-negatives specification v2.2.0 — 2026-04-19*
