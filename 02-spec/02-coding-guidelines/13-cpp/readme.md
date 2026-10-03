# C++ Standards Overview (AI Execution Prompt)

> **/goal** Provide the primary entry point and authoritative index for modern C++ (C++20 baseline) standards, memory safety rules, and RAII architecture.
> **/learn** Master C++20 concepts, smart pointer ownership models, Rule of Zero/Five, PascalCase/snake_case naming conventions, and FFI boundary safety.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Consult and enforce the modern C++ standards module across all C++ codebase components.
- [ ] `/learn` Ban raw `new` and `delete`; mandate `std::unique_ptr` and `std::shared_ptr` for resource management.
- [ ] `/goal` Enforce C++20 baseline with concepts over `enable_if` and strict RAII resource lifetimes.
- [ ] `/learn` Validate C++ guidelines using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 1.0.0
**Updated:** 2026-08-08
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose

This module outlines modern C++ conventions (C++20 baseline) to ensure memory safety, readable structures, and consistent application of RAII principles.

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

## Keywords

`cpp` · `cpp20` · `smart-pointers` · `memory-safety` · `raii`

## Files

| # | File | Category | Description |
|---|------|----------|-------------|
| 01 | [01-standards.md](./02-standards.md) | Logic / Rules | Modern C++ coding standards |
| 97 | [97-acceptance-criteria.md](./97-acceptance-criteria.md) | Acceptance Criteria Registry | Complete testable Gherkin criteria inventory |
| 99 | [99-consistency-report.md](./99-consistency-report.md) | Meta | Consistency and compliance report |

## Cross-References

- [Cross-Language Standards](../01-cross-language/readme.md)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CPP-001: Modern C++ Standards Index Conformance

**Given** Polyglot development guidelines and language standards.
**When** Audited against this language specification.
**Then** Zero non-compliant conventions or syntax patterns are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
