# Golang Coding Standards

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`04-golang-standards-reference` · `coding-standards`

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

Previously a single 1281-line file, now split into focused modules under 300 lines each.

---

## Document Inventory

| # | File | Purpose | Lines |
|---|------|---------|-------|
| — | [01-file-and-function-rules.md](./01-file-and-function-rules.md) | File naming, size, function size, nesting ban | 224 |
| — | [02-type-safety-and-errors.md](./02-type-safety-and-errors.md) | Type safety, error handling, Result types | 362 |
| — | [03-database-and-structs.md](./03-database-and-structs.md) | Database naming, dbutil wrapper, struct design | 123 |
| — | [04-naming-and-organization.md](./04-naming-and-organization.md) | File organization, naming conventions, negations, guards | 272 |
| — | [05-enums-and-dry.md](./05-enums-and-dry.md) | Typed constants, enums, DRY enforcement | 186 |
| — | [06-concurrency-and-patterns.md](./06-concurrency-and-patterns.md) | Concurrency, forbidden patterns, imports, common mistakes | 274 |
| — | 99-consistency-report.md | — | — |

| — | 99-consistency-report.md | — | — |
---

## Cross-References

- [No Raw Negations](../../01-cross-language/12-no-negatives.md) — Positive guard functions (all languages)
- [Cross-Language Code Style](../../01-cross-language/04-code-style/00-overview.md) — Braces, nesting & spacing rules
- [Function Naming](../../01-cross-language/10-function-naming.md) — No boolean flag parameters
- [Strict Typing](../../01-cross-language/13-strict-typing.md) — Type declarations & docblock rules
- [DRY Principles](../../01-cross-language/08-dry-principles.md)
- [Boolean Standards](../02-boolean-standards.md) — Go-specific positive logic rules and exemptions
- apperror Package Spec — Full StackTrace, AppError, Result types specification <!-- external: spec/03-error-manage/01-error-resolution/10-apperror-package/01-apperror-reference.md -->
- [Enum Specification](../01-enum-specification/00-overview.md) — Byte-based enum pattern, required methods, folder structure
- [Master Coding Guidelines](../../01-cross-language/15-master-coding-guidelines/00-overview.md) — Consolidated cross-language reference
- [Issues & Fixes Log](../../01-cross-language/01-issues-and-fixes-log.md) — Full historical fixes
- [golangci-lint Enforcement](../../01-cross-language/16-static-analysis/02-go-golangci-lint.md) — Linter rule mapping for Go guidelines
