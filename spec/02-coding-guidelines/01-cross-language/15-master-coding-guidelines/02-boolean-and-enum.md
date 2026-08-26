# Master Coding Guidelines — Boolean standards, isDefined guards, enum standards

> **Parent:** [Master Coding Guidelines](./00-overview.md)  
> **Version:** 2.1.0  
> **Updated:** 2026-03-31

---

## 3. Boolean Standards — Positive Logic

> Full reference: [boolean-principles.md](../02-boolean-principles/00-overview.md) and [no-negatives.md](../12-no-negatives.md)

### 6 Non-Negotiable Principles

| # | Principle | Rule |
|---|-----------|------|
| P1 | `is`/`has` prefix (99%), `should` rare | Every boolean must start with `is` or `has`. Use `should` only for recommendations. Never use `can`, `was`, `will`. |
| P2 | No negative words | `not`, `no`, `non` are banned from boolean names |
| P3 | Named guards | Never use `!` on function calls — use semantic inverse |
| P4 | Extract complex expressions | 2+ operators → extract to named boolean |
| P5 | No boolean parameters | Use separate named methods or options objects |
| P6 | No mixed polarity | `isX && !isY` → extract to single-intent name |

### Common Mistakes

```php
// ❌ P1 violation: Missing prefix
$active = true;
$loaded = false;

// ✅ CORRECT
$isActive = true;
$isLoaded = false;
```

```php
// ❌ P2 violation: Negative word in name
$isNotReady = true;
$hasNoPermission = true;

// ✅ CORRECT: Positive semantic synonym
$isPending = true;
$isUnauthorized = true;
```

```php
// ❌ P3 violation: Raw negation on function call
if (!$order->isValid()) { return; }
if (!file_exists($path)) { return; }

// ✅ CORRECT: Semantic inverse / guard function
if ($order->isInvalid()) { return; }
if (PathHelper::isFileMissing($path)) { return; }
```

```go
// ❌ P3 violation in Go: Raw negation
if !v.IsValid() {
    return variantLabels[Invalid]
}
if !pathutil.IsDir(gitDir) {
    return err
}

// ✅ CORRECT: Positive counterpart
if v.IsInvalid() {
    return variantLabels[Invalid]
}
if pathutil.IsDirMissing(gitDir) {
    return err
}
```

### Go-Specific Exemptions

These Go patterns are **exempt** from the no-negation rule:
- `if !ok` — idiomatic comma-ok pattern
- `if !requireService(...)` / `if !decodeJSON(...)` — handler guard returns
- `if err != nil` — idiomatic error check
- `if !strings.HasPrefix(...)` — stdlib calls (extract if repeated 3+ times)

---


---

## 3.1 `isDefined` / `isDefinedAndValid` — Positive Null/Existence Guards

> Language-specific details: [Go Standards](../../03-golang/04-golang-standards-reference/00-overview.md) · [PHP Standards](../../04-php/07-php-standards-reference/00-overview.md) · [TypeScript Standards](../../02-typescript/08-typescript-standards-reference.md)

Raw `!== null` / `!= nil` combined with validity checks creates cognitive overhead. Use positive guard methods/functions that express intent as a single word.

### Guard Table (Cross-Language)

| Guard | Replaces | PHP | Go | TypeScript |
|-------|----------|-----|-----|------------|
| **isDefined** | `!= null`, `!= nil` | `$x->isDefined()` | `x.IsDefined()` | `isDefined(x)` |
| **isDefinedAndValid** | `!= null && isValid()` | `$x->isDefinedAndValid()` | `x.IsDefinedAndValid()` | `isDefinedAndValid(x)` |
| **isEmpty** | `== null`, `== nil` | `$x->isEmpty()` | `x.IsEmpty()` | `isEmpty(x)` |
| **isInvalid** | `!isValid()` | `$x->isInvalid()` | `x.IsInvalid()` | — |

### Common Mistakes

```php
// ❌ FORBIDDEN: Nested null + validity check
if ($config !== null) {
    if ($config->isValid()) {
        $this->applyConfig($config);
    }
}

// ❌ FORBIDDEN: Compound with null check
if ($config !== null && $config->isValid()) {
    $this->applyConfig($config);
}

// ✅ REQUIRED: Single positive guard
if ($config->isDefinedAndValid()) {
    $this->applyConfig($config);
}
```

```go
// ❌ FORBIDDEN: Raw nil check
if config != nil {
    applyConfig(config)
}

// ✅ REQUIRED: Positive existence check
if config.IsDefined() {
    applyConfig(config)
}
```

```typescript
// ❌ FORBIDDEN: Raw null/undefined checks
if (config !== null && config !== undefined) {
  applyConfig(config);
}

// ✅ REQUIRED: Type-narrowing guard function
if (isDefined(config)) {
  applyConfig(config); // TypeScript narrows to Config
}
```

### On Result Wrappers (Already Built-in)

| Method | Go | PHP | Meaning |
|--------|-----|-----|---------|
| isDefined | `result.IsDefined()` | `$result->isDefined()` | Value was set |
| isSafe | `result.IsSafe()` | `$result->isSafe()` | Value exists AND no error (≈ isDefinedAndValid) |
| isEmpty | `result.IsEmpty()` | `$result->isEmpty()` | No value set |
| hasError | `result.HasError()` | `$result->hasError()` | Operation failed |

---


---

## 4. Enum Standards

> **Canonical cross-language reference:** [Enum Naming Quick Reference](../../06-ai-optimization/05-enum-naming-quick-reference.md) — declaration, naming, methods, folder structure for Go, TypeScript, and PHP  
> **Go-specific:** [Go Enum Specification](../../03-golang/01-enum-specification/00-overview.md)  
> **PHP-specific:** [PHP Enums](../../04-php/01-enums.md)  
> **TypeScript-specific:** [TypeScript Enums](../../02-typescript/00-overview.md)

**Key rules (all languages):**
- PascalCase members/cases for TypeScript, GoLang, and C# (no magic strings). `camelCase` and `_camelCase` are strictly forbidden. For Rust or other languages, follow the ecosystem's dominant convention.
- Exhaustive `default` branch in every switch/match
- Use comparison methods (`isEqual()` in PHP, `Is{Value}()` in Go, `=== Enum.Member` in TS)

See the canonical sources above for full declaration patterns, required methods, and folder structure.

---
