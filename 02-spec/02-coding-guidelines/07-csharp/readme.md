# C# Coding Standards (AI Execution Prompt)

> **/goal** Master and enforce C#-specific coding standards, PascalCase naming conventions, boolean flag splitting, async patterns, structured exception handling, pattern matching, and nullable reference types across all .NET/C# modules.
> **/learn** Internalize .NET conventions, single responsibility method design (max 15 lines, max 3 parameters), eliminating blocking async calls (.Result), specific exception catching, and compiler-enforced null safety.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce PascalCase for types/methods/properties and camelCase for locals/parameters across all C# code.
- [ ] `/learn` Split boolean flag parameters into distinct, intention-revealing methods and enforce ≤15 lines per method.
- [ ] `/goal` Ban blocking `.Result` and `.GetAwaiter().GetResult()` in favor of pure `async`/`await` and `Task.WhenAll`.
- [ ] `/learn` Prohibit bare `catch (Exception)`, enforce `ArgumentNullException` with `nameof()`, and enable nullable reference types.
- [ ] `/goal` Eliminate type casting in favor of pattern matching (`is` and `switch` expressions) and replace magic strings with enums.
- [ ] `/learn` Verify that all C# specifications adhere to `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Status:** Active
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Keywords

`coding` · `guidelines` · `csharp` · `dotnet` · `.net` · `naming` · `async` · `linq`

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

C#-specific coding standards that extend the [cross-language guidelines](../01-cross-language/readme.md). These rules apply to all .NET/C# code and align with the project's naming conventions, boolean principles, and method design patterns.

---

## File Index

| # | File | Description |
|---|------|-------------|
| 01 | [Naming and Conventions](./02-naming-and-conventions.md) | PascalCase methods, property naming, abbreviation casing |
| 02 | [Method Design](./03-method-design.md) | Boolean flag splitting, async patterns, LINQ usage |
| 03 | [Error Handling](./04-error-handling.md) | Exception patterns, Result types, guard clauses |
| 04 | [Type Safety](./05-type-safety.md) | Generics, nullable reference types, pattern matching |

---

## Document Inventory

| File | Description |
|------|-------------|
| [`02-naming-and-conventions.md`](./02-naming-and-conventions.md) | PascalCase methods, property naming, abbreviation casing |
| [`03-method-design.md`](./03-method-design.md) | Boolean flag splitting, async patterns, LINQ usage |
| [`04-error-handling.md`](./04-error-handling.md) | Exception patterns, Result types, guard clauses |
| [`05-type-safety.md`](./05-type-safety.md) | Generics, nullable reference types, pattern matching |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | C# acceptance criteria registry |
| [`99-consistency-report.md`](../../99-consistency-report.md) | Guideline consistency and audit report |

**Total:** 4 spec files + acceptance criteria, consistency report

---

## Cross-References

- [Cross-Language Guidelines](../01-cross-language/readme.md) — universal rules applied to C#
- [Boolean Flag Methods](../01-cross-language/24-boolean-flag-methods.md) — method splitting pattern with C# examples
- [Boolean Principles](../01-cross-language/02-boolean-principles/readme.md) — boolean naming rules
- [Function Naming](../01-cross-language/10-function-naming.md) — cross-language function naming
- [SOLID Principles](../01-cross-language/23-solid-principles.md) — architecture patterns

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CS-001: C# Coding Guidelines Index & Module Conformance

**Given** C# source code and specifications across modules and projects.
**When** Guideline linters audit the codebase for C# coding standards compliance.
**Then** All C# specifications and source files adhere to PascalCase naming, method sizing, async patterns, structured error handling, and type safety standards with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only
```
**Expected:** exit 0. Zero violations.
