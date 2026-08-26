# Phase 2 — Enums and Coding Style

> **Purpose:** Define enum patterns, coding style, and naming conventions for WordPress plugins.

---

## Index

| File | Purpose |
|------|---------|
| [01-enum-architecture.md](01-enum-architecture.md) | Core enum pattern, standard categories, comparison methods, coding style, naming |
| [02-enum-metadata-pattern.md](02-enum-metadata-pattern.md) | `match`-based metadata methods (label, icon, cssClass) and `is*()` helpers |
| [03-self-update-status-enum.md](03-self-update-status-enum.md) | `SelfUpdateStatusType` — reference impl (17 cases, deployment domain) |
| [04-action-type-enum.md](04-action-type-enum.md) | `ActionType` — reference impl (40+ cases, transaction logging domain) |

---

## Quick Reference

### Standard Enum Template

```php
enum ExampleType: string
{
    case SomeName  = 'some_value';
    case OtherName = 'other_value';

    // Per-case helpers
    public function isSomeName(): bool  { return $this->isEqual(self::SomeName); }
    public function isOtherName(): bool { return $this->isEqual(self::OtherName); }

    // Standard comparison methods
    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Metadata via `match` (PHP)

```php
public function label(): string
{
    return match ($this) {
        self::SomeName  => 'Some Label',
        self::OtherName => 'Other Label',
    };
}
```

See [02-enum-metadata-pattern.md](02-enum-metadata-pattern.md) for the full pattern.

---

## Cross-References

- [Go Enum Specification](../../02-coding-guidelines/03-golang/01-enum-specification/00-overview.md) — equivalent pattern for Go
- [Go Info-Object Pattern](../../02-coding-guidelines/03-golang/01-enum-specification/05-info-object-pattern.md) — Go version of the metadata pattern (uses info-object, not `match`)
- [Phase 10 — Deployment Patterns](../10-deployment-patterns.md) — uses `SelfUpdateStatusType`
