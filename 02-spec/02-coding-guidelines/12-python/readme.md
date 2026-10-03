# Python Standards Overview (AI Execution Prompt)

> **/goal** Provide the primary entry point and authoritative index for Python language standards, enforcing PEP-8 compliance, strict static typing, and schema validation.
> **/learn** Master Pydantic model validation at system boundaries, Black code formatting (100 char limit), elimination of bare `except`, and explicit type annotations.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Consult and enforce the Python standards module across all internal automation and tooling scripts.
- [ ] `/learn` Avoid untyped signatures, `Any` types, and untyped dictionary bags at domain interfaces.
- [ ] `/goal` Require Black formatting, 100 character line length ceiling, and granular exception hierarchy handling.
- [ ] `/learn` Validate Python guidelines using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 1.0.0
**Updated:** 2026-08-08
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose

This module defines the standard practices for writing Python code within the repository, focusing on PEP-8 compliance, strict typing, and data validation.

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

## Keywords

`python` · `pep-8` · `pydantic` · `typing` · `data-validation`

## Files

| # | File | Category | Description |
|---|------|----------|-------------|
| 01 | [01-standards.md](./02-standards.md) | Logic / Rules | Python-specific coding standards |
| 97 | [97-acceptance-criteria.md](./97-acceptance-criteria.md) | Acceptance Criteria Registry | Complete testable Gherkin criteria inventory |
| 99 | [99-consistency-report.md](./99-consistency-report.md) | Meta | Consistency and compliance report |

## Cross-References

- [Cross-Language Standards](../01-cross-language/readme.md)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PY-001: Python Standards Index Conformance

**Given** Polyglot development guidelines and language standards.
**When** Audited against this language specification.
**Then** Zero non-compliant conventions or syntax patterns are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
