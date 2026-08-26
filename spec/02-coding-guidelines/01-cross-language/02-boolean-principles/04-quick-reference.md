# Boolean Principles — Quick reference table, common mistakes

> **Parent:** [Boolean Principles](./00-overview.md)  
> **Version:** 2.6.0  
> **Updated:** 2026-03-31

---

## Quick Reference

| ❌ Forbidden | ✅ Required | Principle |
|-------------|------------|-----------|
| `$active` | `$isActive` | P1: `is`/`has` prefix |
| `$loaded` | `$isLoaded` | P1: `is`/`has` prefix |
| `!isNotBlocked` | `isBlocked` | P2: No negative words |
| `isNotBlocked` | `isActive` (synonym) | P2: No negative words |
| `isNotReady` | `isPending` (synonym) | P2: No negative words |
| `!$obj->isValid()` | `$obj->isInvalid()` | P3: Named guards |
| `if (a && b \|\| c)` | `if (isValid(x))` | P4: Extract expressions |
| `fn(true)` | `fnWithOption()` | P5: Explicit params |
| `isX && !isY` | `isConflict` (extracted) | P6: No mixed polarity |
| `if x := fn(); x > 0` | Separate assignment | P7: No inline statements |
| `os.Stat(path)` | `pathutil.IsDir(path)` | P8: No raw filesystem |

---


---

## Common Mistakes — Boolean Logic

### Mistake 1: Missing `is`/`has` Prefix (P1)

```php
// ❌ WRONG
$active = true;        // What is active? No semantic meaning.
$loaded = false;       // Ambiguous.

// ✅ CORRECT
$isActive = true;
$isLoaded = false;
```

### Mistake 2: Negative Word in Name (P2)

```go
// ❌ WRONG — "not" in the name
isNotReady := order.Status != "ready"
hasNoPermission := !user.HasPermission("admin")

// ✅ CORRECT — positive semantic synonym
isPending := order.Status != "ready"
isUnauthorized := !user.HasPermission("admin")
```

### Mistake 3: Raw `!` on Function Call (P3)

```php
// ❌ WRONG
if (!$order->isValid()) { return; }
if (!file_exists($path)) { return; }

// ✅ CORRECT
if ($order->isInvalid()) { return; }
if (PathHelper::isFileMissing($path)) { return; }
```

```go
// ❌ WRONG
if !v.IsValid() { return variantLabels[Invalid] }
if !pathutil.IsDir(gitDir) { return err }

// ✅ CORRECT
if v.IsInvalid() { return variantLabels[Invalid] }
if pathutil.IsDirMissing(gitDir) { return err }
```

### Mistake 4: Mixed Polarity in Condition (P6)

```typescript
// ❌ WRONG — positive + negative in same if
if (isLoggedIn && !hasPermission) {
    redirect('/unauthorized');
}

// ✅ CORRECT — extract to single-intent name
const isUnauthorized = isLoggedIn && !hasPermission;
if (isUnauthorized) {
    redirect('/unauthorized');
}
```

### Mistake 5: Inline Statement in `if` (P7)

```go
// ❌ WRONG — inline os.Stat hides computation
if _, err := os.Stat(dir); err == nil {
    fmt.Println("exists")
}

// ✅ CORRECT — separate computation from condition
isProjectExists := pathutil.IsDir(dir)
if isProjectExists {
    fmt.Println("exists")
}
```

### Mistake 6: Raw Filesystem Call (P8)

```go
// ❌ WRONG — raw os.MkdirAll with raw error
if err := os.MkdirAll(dir, 0755); err != nil {
    return fmt.Errorf("mkdir failed: %w", err)
}

// ✅ CORRECT — pathutil wrapper with apperror
if err := pathutil.EnsureDir(dir); err != nil {
    return apperror.Fail[ProjectResult](err)
}
```

### Mistake 7: All Rules Violated Together (P6 + P7 + P8)

```go
// ❌ WRONG — P6 (mixed polarity) + P7 (inline stmt) + P8 (raw os.Stat) + raw error
isProjectConflict := err == nil && !isOverwrite
if _, err := os.Stat(projectDir); isProjectConflict {
    return fmt.Errorf("project exists, use isOverwrite=true to replace")
}

// ✅ CORRECT — all rules applied
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

### Go-Specific Exemptions

These patterns are **exempt** from the no-negation rule in Go:
- `if !ok` — idiomatic comma-ok pattern
- `if !requireService(w, svc, "name")` — handler guard returns
- `if err != nil` — idiomatic error check
- `if !strings.HasPrefix(...)` — stdlib calls (extract if repeated 3+ times)
- `if v, ok := m[k]; ok {` — inline comma-ok (exempt from P7)

---
