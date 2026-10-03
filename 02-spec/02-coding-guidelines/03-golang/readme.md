# Golang Standards (AI Execution Prompt)

> **/goal** Master and enforce Go-specific coding standards, affirmative boolean patterns, error handling conventions, resource defer rules, and clean architecture across all Go packages and CLI modules.
> **/learn** Internalize positive boolean naming (`Is*`, `Has*`), elimination of inline negations, `*appfault.AppError` and monadic `Result[T]` returns, single-defer lifecycle rules, and typed enum usage.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce affirmative `is` and `has` boolean prefixes; ban explicit `== true` evaluations and mixed polarity conditionals.
- [ ] `/learn` Never use multiple `defer` calls in a single function; extract subroutines or use explicit resource management.
- [ ] `/goal` Use `httpmethodtype.Variant` enum values and methods instead of raw HTTP method string literals (`"GET"`, `"POST"`).
- [ ] `/learn` Verify that all Go packages and specifications comply with `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Status:** Active
**Updated:** 2026-04-16
**AI Confidence:** High
**Ambiguity:** None

---

## Keywords

`coding`, `golang`, `guidelines`

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

Go-specific coding standards and patterns.

---

## Document Inventory

| File | Description |
|------|-------------|
| [`02-boolean-standards.md`](./02-boolean-standards.md) | Positive boolean naming, negation elimination, and guard clauses |
| [`03-httpmethod-enum.md`](./03-httpmethod-enum.md) | HttpMethod enum standard and magic string elimination |
| [`04-golang-standards-reference.md`](./04-golang-standards-reference.md) | Comprehensive Golang standards and conventions reference |
| [`05-defer-rules.md`](./05-defer-rules.md) | Resource defer rules and loop safety standards |
| [`06-string-slice-internals.md`](./06-string-slice-internals.md) | Go string and slice memory representation and preallocation |
| [`07-code-severity-taxonomy.md`](./07-code-severity-taxonomy.md) | Code severity taxonomy, error classification, and fault logging |
| [`08-pathutil-fileutil-spec.md`](./08-pathutil-fileutil-spec.md) | Unified PathUtil & FileUtil cross-platform filesystem specification |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | Golang acceptance criteria registry |
| [`98-changelog.md`](./98-changelog.md) | Version history and specification changelog |
| [`99-consistency-report.md`](./99-consistency-report.md) | Guideline consistency and audit report |

---

## Cross-References

_See parent folder's [`readme.md`](../readme.md) for broader context._

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-000: Golang Coding Standards Index & Specification Conformance

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** All Go specifications and source files adhere to affirmative boolean conventions, single-defer rules, typed enums, and structured error management with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.
