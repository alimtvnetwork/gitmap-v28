# Phase 2 — Enums and Coding Style

. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


> **Purpose:** Define enum patterns, coding style, and naming conventions for WordPress plugins.

---

## Index

| File | Purpose |
|------|---------|
| [02-enum-architecture.md](02-enum-architecture.md) | Core enum pattern, standard categories, comparison methods, coding style, naming |
| [03-enum-metadata-pattern.md](03-enum-metadata-pattern.md) | `match`-based metadata methods (label, icon, cssClass) and `is*()` helpers |
| [04-self-update-status-enum.md](04-self-update-status-enum.md) | `SelfUpdateStatusType` — reference impl (17 cases, deployment domain) |
| [05-action-type-enum.md](05-action-type-enum.md) | `ActionType` — reference impl (40+ cases, transaction logging domain) |

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

See [03-enum-metadata-pattern.md](03-enum-metadata-pattern.md) for the full pattern.

---

## Cross-References

- [Go Enum Specification](../../02-coding-guidelines/03-golang/01-enum-specification/01-index.md) — equivalent pattern for Go
- [Go Info-Object Pattern](../../02-coding-guidelines/03-golang/01-enum-specification/06-info-object-pattern.md) — Go version of the metadata pattern (uses info-object, not `match`)
- [Phase 10 — Deployment Patterns](../11-deployment-patterns.md) — uses `SelfUpdateStatusType`
