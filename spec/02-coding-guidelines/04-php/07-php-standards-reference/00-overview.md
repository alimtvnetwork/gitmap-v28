# PHP Coding Standards

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`07-php-standards-reference` · `coding-standards`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Previously a single 841-line file, now split into focused modules under 300 lines each.

---

## Document Inventory

| # | File | Purpose | Lines |
|---|------|---------|-------|
| — | [01-naming-and-errors.md](./01-naming-and-errors.md) | Naming conventions, error handling, structured responses | 158 |
| — | [02-constants-and-deps.md](./02-constants-and-deps.md) | Constants, enums, dependency checks, file paths | 146 |
| — | [03-initialization-and-booleans.md](./03-initialization-and-booleans.md) | Constructor rules, boolean logic, isDefined guards | 252 |
| — | [04-code-style.md](./04-code-style.md) | Braces, nesting, spacing, function size | 235 |
| — | [05-forbidden-and-database.md](./05-forbidden-and-database.md) | Forbidden patterns, database wrapper | 94 |
| — | 99-consistency-report.md | — | — |

| — | 99-consistency-report.md | — | — |
---

## Cross-References

- WordPress Plugin Development Spec — Full 10-document guide *(Phase 4 target)*
- [Error Handling Spec](../../../03-error-manage/02-error-architecture/01-error-handling-reference.md) — Cross-language error strategy
- Generic Enforce Spec — Type safety rules *(Phase 5/6 target)*
- [DRY Principles](../../01-cross-language/08-dry-principles.md) — Cross-language DRY rules
- [Cross-Language Code Style](../../01-cross-language/04-code-style/00-overview.md) — Braces, nesting & spacing rules (canonical)
- [Function Naming](../../01-cross-language/10-function-naming.md) — No boolean flag parameters (all languages)
- [Strict Typing](../../01-cross-language/13-strict-typing.md) — Type declarations & docblock rules (all languages)
