# Boolean Principles

. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Keywords

`02-boolean-principles` · `coding-standards`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `01-index.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Previously a single 858-line file, now split into focused modules under 300 lines each.

---

## Document Inventory

| # | File | Purpose | Lines |
|---|------|---------|-------|
| — | [01-naming-prefixes.md](./02-naming-prefixes.md) | P1: is/has prefixes, P2: no negative words | 134 |
| — | [02-guards-and-extraction.md](./03-guards-and-extraction.md) | P3: named guards, P4: extract complex expressions | 205 |
| — | [03-parameters-and-conditions.md](./04-parameters-and-conditions.md) | P5: explicit params, P6: no mixed booleans, P7: no inline statements, P8: no raw system calls, P9: no explicit true checks | 262 |
| — | [04-quick-reference.md](./05-quick-reference.md) | Quick reference table, common mistakes | 155 |
| — | [05-exemptions-and-api.md](./06-exemptions-and-api.md) | Static factory exemption, Result wrapper API | 139 |
| — | 99-consistency-report.md | — | — |

| — | 99-consistency-report.md | — | — |
---

## Database ↔ Code Inverse Pattern (Rule 9)

When a boolean originates in the **database**, the storage layer holds the
canonical positive form (e.g. `IsActive`) and the **inverted sibling**
(e.g. `IsInactive`) is auto-generated as a computed property in code —
never as a second column. This is the database-side counterpart to the
in-memory semantic-inverse pairs documented in
[`12-no-negatives.md`](../12-no-negatives.md#object-level-semantic-inverses).

> **Authoritative spec:** [Database Naming Conventions — Rule 9: Auto-Generated Inverted (Computed) Fields](../../../04-database-conventions/02-naming-conventions.md#rule-9-auto-generated-inverted-computed-fields-in-code)
>
> **Codegen tool:** [`linters-cicd/codegen/`](../../../../linters-cicd/codegen/readme.md) — emits Go methods, PHP traits, and TypeScript getters from `Is*`/`Has*` db-tagged fields.
>
> **Linter:** `BOOL-NEG-001` rejects `Not`/`No`-prefixed column names (`IsNotActive`, `HasNoLicense`) at CI time.

---

## Cross-References

- [No Raw Negations](../12-no-negatives.md) — Full guard function inventory
- [Database Naming — Rule 9 (Inverted Fields)](../../../04-database-conventions/02-naming-conventions.md#rule-9-auto-generated-inverted-computed-fields-in-code) — DB-side inverse contract
- [Code Style — Rule 3](../04-code-style/01-index.md) — Complex condition extraction
- [Function Naming](../10-function-naming.md) — No boolean flag parameters
- [PHP Boolean Guard Inventory](../../04-php/07-php-standards-reference/01-index.md) — PHP-specific helpers
- [Go Boolean Standards](../../03-golang/02-boolean-standards.md) — Go-specific rules and exemptions (P7, P8)
- [Master Coding Guidelines](../15-master-coding-guidelines/01-index.md) — Consolidated reference
- [Issues & Fixes Log](../02-issues-and-fixes-log.md) — Historical fixes
- [apperror Package — Result Guard Rule](../../../03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/08-serialization-and-guards.md#12-result-guard-rule--mandatory-error-check-before-value-access)
