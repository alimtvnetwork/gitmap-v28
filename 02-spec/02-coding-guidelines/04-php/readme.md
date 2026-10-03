# PHP Standards (AI Execution Prompt)

> **/goal** Master and enforce PHP-specific coding standards, native backed enums with Type suffix, zero magic strings, affirmative boolean patterns, PascalCase response keys, and clean architecture across all PHP modules.
> **/learn** Internalize PHP 8.1+ backed enums with `isEqual()`, PSR-12 conventions, eliminating raw negations with positive guards, file-level trait namespace imports, ResultHelper responses, and ResponseKeyType usage.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce PHP 8.1+ native backed enums with mandatory `Type` suffix and universal `isEqual()` method.
- [ ] `/learn` Eliminate all forbidden patterns: catch `Throwable` instead of `Exception`, prohibit `wp_die()` in REST handlers, and ban raw negations.
- [ ] `/goal` Standardize service result arrays with `ResultHelper` factory methods and `ResponseKeyType` enum keys in PascalCase.
- [ ] `/learn` Verify that all PHP files, trait imports, and specifications comply with `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Status:** Active
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Keywords

`coding`, `guidelines`, `php`, `enums`, `naming`, `spacing`, `response-key`

---

## Purpose

PHP-specific coding standards and patterns for the RiseupAsia namespace.

---

## Document Inventory

| File | Description |
|------|-------------|
| [`02-enums.md`](./02-enums.md) | PHP 8.1+ native backed enums, Type suffix, and isEqual() comparison standard |
| [`03-forbidden-patterns.md`](./03-forbidden-patterns.md) | Forbidden PHP patterns and required replacements checklist |
| [`04-naming-conventions.md`](./04-naming-conventions.md) | PHP naming conventions, PSR-12 baseline, and case rules |
| [`05-response-array-standard.md`](./05-response-array-standard.md) | Response array standards, ResultHelper, and ResponseKeyType |
| [`06-php-standards-reference.md`](./06-php-standards-reference.md) | Comprehensive PHP standards reference |
| [`07-spacing-and-imports.md`](./07-spacing-and-imports.md) | Spacing, vertical whitespace, and file-level import rules |
| [`08-response-key-type-inventory.md`](./08-response-key-type-inventory.md) | ResponseKeyType case inventory (176 cases) |
| [`09-php-go-consistency-audit.md`](./09-php-go-consistency-audit.md) | PHP–Go cross-language consistency audit |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | PHP acceptance criteria registry |
| [`99-consistency-report.md`](../../99-consistency-report.md) | Guideline consistency and audit report |

**Total:** 8 spec files + acceptance criteria, consistency report

---

## Cross-References

- [Cross-Language Guidelines](../01-cross-language/readme.md)
- [Go Standards](../03-golang/readme.md) — for PHP–Go parity
- [Parent Overview](../readme.md)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PHP-001: PHP Guidelines Index & Module Conformance

**Given** PHP source code and specifications across modules and plugins.
**When** Guideline linters audit the codebase.
**Then** All PHP specifications and source files adhere to native backed enum standards, forbidden patterns elimination, PSR-12 naming conventions, and PascalCase response array keys with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.
