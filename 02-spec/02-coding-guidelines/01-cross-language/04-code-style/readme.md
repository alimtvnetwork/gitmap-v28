# Cross-Language Code Style — Braces, Nesting, Spacing & Function Size (AI Execution Prompt)

> **/goal** Standardize cross-language control-flow formatting, mandatory brace enforcement, zero-nesting early return guards, vertical line spacing, and strict function and type size caps across PHP, TypeScript, and Go.
> **/learn** Master the code style rules (Rules 1-17), guard clause flattening patterns, discrete condition extraction, mandatory vertical blank line spacing rules, and small-function decomposition architectures.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce mandatory braces and ban all nested if statements across control flow via early guard returns.
- [ ] `/learn` Extract complex multi-part conditions into discrete, positively named intermediate boolean variables.
- [ ] `/goal` Maintain strict vertical blank line spacing: before if, after }, before return, and around multiline structures.
- [ ] `/learn` Enforce function length limits (max 15 lines) and type/struct/class caps (max 120 lines) through clean decomposition.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

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
| `readme.md` present | ✅ |
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
| 01 | [01-braces-and-nesting.md](./02-braces-and-nesting.md) | Brace enforcement, zero-nesting ban, exemptions | 1, 2, 7 |
| 02 | [02-conditions-and-extraction.md](./03-conditions-and-extraction.md) | Extract complex multi-part conditions | 3 |
| 03 | [03-blank-lines-and-spacing.md](./04-blank-lines-and-spacing.md) | Blank lines before/after blocks and control structures | 4, 5, 10 |
| 04 | [04-function-and-type-size.md](./05-function-and-type-size.md) | 15-line function limit, 120-line struct/class limit | 6, 17 |
| 05 | [05-multi-line-formatting.md](./06-multi-line-formatting.md) | Multi-line arguments, method chaining, apperror formatting | 9, 11, apperror |
| 06 | [06-comments-and-documentation.md](./07-comments-and-documentation.md) | Comment formatting, doc comments, dead code, backslash rule | 8, 14, 15, 16 |
| 07 | [08-checklist.md](./08-checklist.md) | PR checklist summary + cross-references | — |
| — | 99-consistency-report.md | — | — |

---

## Core Rules Summary & Code Examples

### Rule 1: Mandatory Braces & No Single-Line Statements

- ❌ FORBIDDEN:
  ```typescript
  if (isValid) return true;
  ```
- ✅ REQUIRED:
  ```typescript
  if (isValid) {
    return true;
  }
  ```

### Rule 2: Vertical Spacing Hygiene

- ❌ FORBIDDEN:
  ```go
  result := compute()
  return result
  ```
- ✅ REQUIRED:
  ```go
  result := compute()

  return result
  ```

---

## Cross-References

- [Parent Overview](../readme.md) — Cross-Language root
- [Boolean Principles](../02-boolean-principles/readme.md) — P1–P6 boolean naming rules
- [No Raw Negations](../12-no-negatives.md) — Positive guard functions
- [Function Naming](../10-function-naming.md) — No boolean flag parameters
- [Strict Typing](../13-strict-typing.md) — Type declarations, max 3 parameters
- [Go Enum Specification](../../03-golang/01-enum-specification/readme.md) — Go enum pattern
- [TypeScript Enums](../../02-typescript/readme.md) — TypeScript string enums
- [PHP Enum Classes](../../04-php/02-enums.md) — PHP backed enum patterns

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-STYLE-001: Cross-Language Code Style Index & Directory Conformance

**Given** The code style guideline directory `02-spec/02-coding-guidelines/01-cross-language/04-code-style/`.
**When** Codebases and guideline specifications are audited for control-flow formatting, brace standards, and size caps.
**Then** All code style rules are deterministically verifiable with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style --check-only
```
**Expected:** exit 0. Zero violations.
