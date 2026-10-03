# Master Coding Guidelines (AI Execution Prompt)

> **/goal** Provide a consolidated, authoritative cross-language reference of non-negotiable coding standards, naming conventions, boolean principles, and structural rules across all supported stacks.
> **/learn** Master the cross-language architectural invariants, PascalCase naming conventions, affirmative boolean logic, strict typing, and zero-error-suppression standards.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify all code changes across Go, PHP, and TypeScript strictly adhere to universal naming and structure conventions.
- [ ] `/learn` Enforce affirmative boolean naming (`is*`, `has*`), positive null guards (`isDefined`), and zero raw negations on function returns.
- [ ] `/goal` Eliminate magic strings, nested conditionals, loose parameter signatures, and unhandled Result errors across all packages.
- [ ] `/learn` Validate compliance across master guidelines using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Updated:** 2026-10-02
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Keywords

`15-master-coding-guidelines` · `coding-standards`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `readme.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Previously a single 1122-line file, now split into focused modules under 300 lines each.

---

## Document Inventory

| # | File | Purpose | Lines |
|---|------|---------|-------|
| — | [01-naming-and-database.md](./02-naming-and-database.md) | Naming conventions, database naming, file naming | 174 |
| — | [02-boolean-and-enum.md](./03-boolean-and-enum.md) | Boolean standards, isDefined guards, enum standards | 213 |
| — | [03-code-style-and-errors.md](./04-code-style-and-errors.md) | Code style formatting, error handling | 277 |
| — | [04-type-safety.md](./05-type-safety.md) | Type safety, single return value, no casting | 221 |
| — | [05-magic-strings-and-organization.md](./06-magic-strings-and-organization.md) | Magic strings, file organization, array keys | 127 |
| — | [06-advanced-patterns.md](./07-advanced-patterns.md) | Lint, enum sync, tests, lazy eval, regex, mutation, null safety, nesting, newlines, defer | 177 |
| — | [07-checklist.md](./08-checklist.md) | Quick checklist for any code change | 41 |
| — | 99-consistency-report.md | — | — |

---

## Cross-References

- [Cross-Language Standards](../readme.md)
- [Boolean Principles](../02-boolean-principles/readme.md)
- [Code Style](../04-code-style/readme.md)
- [Parent Acceptance Criteria](../../97-acceptance-criteria.md)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-MASTER-INDEX: Master Coding Guidelines Directory Index & Conformance

**Given** Master coding guidelines directory `02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/`.
**When** Linters and CI/CD pipelines audit the master specification documents for actionable checklists, valid links, and structural compliance.
**Then** All documents pass automated checks with zero dead links and zero guideline violations.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
