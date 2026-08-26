# Cross-Language Code Style — Braces, Nesting, Spacing & Function Size

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`code-style` · `braces` · `nesting` · `spacing` · `function-size` · `formatting` · `cross-language`

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

Cross-language code style rules governing control-flow formatting and function design across **PHP, TypeScript, and Go**. Previously a single 1,458-line file, now split into focused modules under 300 lines each.

These rules are the **single source of truth** — language-specific specs reference this folder.

---

## Document Inventory

| # | File | Purpose | Rules |
|---|------|---------|-------|
| 01 | [01-braces-and-nesting.md](./01-braces-and-nesting.md) | Brace enforcement, zero-nesting ban, exemptions | 1, 2, 7 |
| 02 | [02-conditions-and-extraction.md](./02-conditions-and-extraction.md) | Extract complex multi-part conditions | 3 |
| 03 | [03-blank-lines-and-spacing.md](./03-blank-lines-and-spacing.md) | Blank lines before/after blocks and control structures | 4, 5, 10 |
| 04 | [04-function-and-type-size.md](./04-function-and-type-size.md) | 15-line function limit, 120-line struct/class limit | 6, 17 |
| 05 | [05-multi-line-formatting.md](./05-multi-line-formatting.md) | Multi-line arguments, method chaining, apperror formatting | 9, 11, apperror |
| 06 | [06-comments-and-documentation.md](./06-comments-and-documentation.md) | Comment formatting, doc comments, dead code, backslash rule | 8, 14, 15, 16 |
| 07 | [07-checklist.md](./07-checklist.md) | PR checklist summary + cross-references | — |
| — | 99-consistency-report.md | — | — |

| — | 99-consistency-report.md | — | — |
---

## Cross-References

- [Parent Overview](../00-overview.md) — Cross-Language root
- [Boolean Principles](../02-boolean-principles/00-overview.md) — P1–P6 boolean naming rules
- [No Raw Negations](../12-no-negatives.md) — Positive guard functions
- [Function Naming](../10-function-naming.md) — No boolean flag parameters
- [Strict Typing](../13-strict-typing.md) — Type declarations, max 3 parameters
- [Go Enum Specification](../../03-golang/01-enum-specification/00-overview.md) — Go enum pattern
- [TypeScript Enums](../../02-typescript/00-overview.md) — TypeScript string enums
- [PHP Enum Classes](../../04-php/01-enums.md) — PHP backed enum patterns
