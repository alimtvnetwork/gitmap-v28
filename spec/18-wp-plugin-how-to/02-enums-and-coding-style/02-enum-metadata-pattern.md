# Enum Metadata Pattern

**Version:** 2.0.0
**Status:** Complete
**Updated:** 2026-04-09

---

## Purpose

Define how enums expose metadata (labels, icons, CSS classes, etc.) via `match` expressions in PHP. Each metadata field is a separate method using PHP's `match` expression.

> **Note:** Go and TypeScript use the info-object / lookup-map pattern instead. See [Go Info-Object Pattern](../../02-coding-guidelines/03-golang/01-enum-specification/05-info-object-pattern.md) for details. This document covers **PHP only**.

---

## Pattern — One `match` Method per Metadata Field

Each metadata field (label, icon, cssClass, etc.) gets its own method using a `match` expression:

```php
enum StatusType: string
{
    case Success = 'Success';
    case Failed  = 'Failed';
    case Pending = 'Pending';

    public function label(): string
    {
        return match ($this) {
            self::Success => 'Operation succeeded',
            self::Failed  => 'Operation failed',
            self::Pending => 'Awaiting processing',
        };
    }

    public function icon(): string
    {
        return match ($this) {
            self::Success => '✅',
            self::Failed  => '❌',
            self::Pending => '⏳',
        };
    }

    public function cssClass(): string
    {
        return match ($this) {
            self::Success => 'text-green-600',
            self::Failed  => 'text-red-600',
            self::Pending => 'text-yellow-600',
        };
    }
}
```

### Why This Works for PHP

- PHP's `match` expression is **compile-time optimized** — equivalent to a hash lookup
- Each method is self-contained and easy to read
- Adding a new metadata field = adding one new method
- Adding a new enum case triggers compile-time errors in all `match` blocks (exhaustiveness)
- No extra classes or objects needed — zero overhead

---

## Individual `is*()` Methods — Per-Case Helpers

Every enum should provide an `is*()` method for **each case**. This makes calling code more readable:

```php
// ✅ Readable
if ($action->isUpload()) { ... }

// ❌ Verbose — avoid in calling code
if ($action->isEqual(ActionType::Upload)) { ... }
```

Implementation — each `is*()` delegates to `isEqual()`:

```php
public function isUpload(): bool { return $this->isEqual(self::Upload); }
public function isEnable(): bool { return $this->isEqual(self::Enable); }
```

### Group Helper Methods

When multiple cases share a domain concept, add a group helper:

| Pattern | Example |
|---------|---------|
| Prefix-based groups | `isSnapshot()` → `str_starts_with($this->value, 'Snapshot')` |
| Explicit groups | `isLifecycle()` → `isAnyOf(Enable, Disable, Delete)` |

Group helpers should be added when:
- Multiple cases share a **domain concept**
- The combination appears in **2+ call sites**
- The prefix-based shortcut (`str_starts_with`) is applicable

---

## Complete PHP Example

```php
namespace PluginName\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum SomeStatusType: string
{
    case Active   = 'Active';
    case Inactive = 'Inactive';
    case Archived = 'Archived';

    // ── Metadata Methods ───────────────────────────────────

    public function label(): string
    {
        return match ($this) {
            self::Active   => 'Currently active',
            self::Inactive => 'Inactive',
            self::Archived => 'Archived',
        };
    }

    public function icon(): string
    {
        return match ($this) {
            self::Active   => '🟢',
            self::Inactive => '⚪',
            self::Archived => '📦',
        };
    }

    public function cssClass(): string
    {
        return match ($this) {
            self::Active   => 'status-active',
            self::Inactive => 'status-inactive',
            self::Archived => 'status-archived',
        };
    }

    // ── Per-Case Helpers ────────────────────────────────────

    public function isActive(): bool   { return $this->isEqual(self::Active); }
    public function isInactive(): bool { return $this->isEqual(self::Inactive); }
    public function isArchived(): bool { return $this->isEqual(self::Archived); }

    // ── Group Helpers ───────────────────────────────────────

    public function isDisabled(): bool
    {
        return $this->isAnyOf(self::Inactive, self::Archived);
    }

    // ── Standard Comparison Methods ─────────────────────────

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

---

## When to Add Metadata Methods

| Scenario | Add `label()`? | Add `icon()` / `cssClass()`? |
|----------|---------------|------------------------------|
| Enum values shown in UI | ✅ Yes | If UI needs icons/styling |
| Enum used only in logic | ❌ No | ❌ No |
| Config enums (int-backed values) | ❌ No | ❌ No |

Not every enum needs metadata. Add `label()` only when values are displayed to users. Add `icon()` / `cssClass()` only when the UI requires them.

---

## Rules

### R1: Use `match` for PHP Metadata

PHP enums use `match` expressions for metadata. No `EnumInfo` class, no lookup maps.

### R2: One Method per Metadata Field

Each metadata type (label, icon, cssClass) is a separate method. Do not combine multiple metadata fields into a single method or data structure.

### R3: Every Case Must Be Covered

PHP's `match` is exhaustive — every case must have an entry. Missing cases cause compile-time errors.

### R4: Per-Case `is*()` Methods Are Mandatory

Every enum must provide individual `is*()` methods for each case. Group helpers should be added when a domain concept spans multiple cases.

### R5: Standard Comparison Methods Are Mandatory

Every enum must include `isEqual()`, `isOtherThan()`, and `isAnyOf()`.

---

## Cross-Language Notes

| Language | Metadata Pattern |
|----------|-----------------|
| PHP 8.1+ | `match` expression per method (this document) |
| Go | Info-object with `map[Variant]EnumInfo` — see [Go Info-Object Pattern](../../02-coding-guidelines/03-golang/01-enum-specification/05-info-object-pattern.md) |
| TypeScript | Info-object with `Record<EnumValue, EnumInfo>` |

PHP uses `match` because it is idiomatic, compile-time optimized, and requires no extra classes. Go and TypeScript use the info-object pattern because their type systems benefit from centralised metadata maps.

---

## Cross-References

- [01-enum-architecture.md](01-enum-architecture.md) — core enum structure and comparison methods
- [03-self-update-status-enum.md](03-self-update-status-enum.md) — reference impl (17 cases, deployment domain)
- [04-action-type-enum.md](04-action-type-enum.md) — reference impl (40+ cases, transaction logging)
- [Go Info-Object Pattern](../../02-coding-guidelines/03-golang/01-enum-specification/05-info-object-pattern.md)

---

*PHP enum metadata via `match` expressions — simple, performant, idiomatic.*
