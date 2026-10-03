# Rust Coding Standards (AI Execution Prompt)

> **/goal** Master and enforce Rust-specific coding standards, RFC 430 idiomatic conventions, explicit PascalCase boundaries for databases and enum serialization, structured error hierarchies, and memory-safe systems programming.
> **/learn** Internalize Rust naming rules (snake_case default with PascalCase DB/enum exceptions), dual `thiserror`/`anyhow` error patterns, Tokio async runtime discipline, strict `unsafe` justification comments, and trait-based platform abstractions.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce Rust community naming conventions (RFC 430) with mandatory PascalCase for database identifiers and enum string values.
- [ ] `/learn` Implement structured error handling using `thiserror` for domain crates and `anyhow` for application orchestration.
- [ ] `/goal` Ensure all `unsafe` blocks contain explicit `// SAFETY:` justifications and Tokio async channels use bounded capacities.
- [ ] `/learn` Verify that all Rust specifications and acceptance criteria pass validation with `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Status:** Active
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Keywords

`coding`, `guidelines`, `rust`, `snake-case`, `naming`, `database-pascalcase`, `enum-string-values`

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

## ⚠️ AI Critical — Rust Naming Override

```
Rust is the ONLY language in this project that follows its own community conventions
for identifier naming (snake_case functions, snake_case variables, SCREAMING_SNAKE_CASE
constants) instead of the project-wide PascalCase mandate.

PascalCase is MANDATORY in Rust for exactly two things:
  1. Database identifiers (table names, column names, view names, primary keys)
  2. Enum string values (when an enum variant serializes to a string)

Everything else → standard Rust community conventions (RFC 430).

See 02-naming-conventions.md for the complete reference with examples.
```

---

## Purpose

Rust-specific coding standards for the Time Log CLI and any future Rust-based projects. Extends the [Cross-Language Guidelines](../01-cross-language/readme.md) but **overrides the naming convention** to follow Rust community standards (RFC 430) with two explicit PascalCase exceptions at system boundaries.

This override exists because Rust's compiler actively enforces `snake_case` for functions/variables via lint warnings, making the project-wide PascalCase mandate impractical for Rust code. The two exceptions (database and enum strings) are where Rust code interacts with other systems in the stack that expect PascalCase.

---

## Document Inventory

| File | Description |
|------|-------------|
| [`02-naming-conventions.md`](./02-naming-conventions.md) | Rust naming rules: snake_case default, PascalCase for DB + enum strings, serialization, module structure |
| [`03-error-handling.md`](./03-error-handling.md) | Error types, Result patterns, thiserror/anyhow usage |
| [`04-async-patterns.md`](./04-async-patterns.md) | Tokio async conventions, channel patterns, cancellation |
| [`05-memory-safety.md`](./05-memory-safety.md) | Ownership idioms, lifetime rules, unsafe policy |
| [`06-testing-standards.md`](./06-testing-standards.md) | Unit/integration test structure, mocking, property testing |
| [`07-ffi-platform.md`](./07-ffi-platform.md) | FFI safety rules, conditional compilation, platform abstractions |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | Compliance requirements |
| [`99-consistency-report.md`](../../99-consistency-report.md) | Structural health |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Guidelines | `../01-cross-language/readme.md` |
| Coding Guidelines Root | `../readme.md` |
| Database Conventions (PascalCase) | `../../../04-database-conventions/readme.md` |
| Enum Standards (Cross-Language) | `../../../../17-consolidated-guidelines/07-enum-standards.md` |
| [`02-naming-conventions.md`](./02-naming-conventions.md) | Rust naming and casing specifications |
| [`03-error-handling.md`](./03-error-handling.md) | Error hierarchy and Result specifications |
| [`04-async-patterns.md`](./04-async-patterns.md) | Tokio async patterns and channel guidelines |
| [`05-memory-safety.md`](./05-memory-safety.md) | Ownership, lifetimes, and unsafe guidelines |
| [`06-testing-standards.md`](./06-testing-standards.md) | Unit/integration test standards and mocking |
| [`07-ffi-platform.md`](./07-ffi-platform.md) | FFI boundaries and platform abstraction |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | Acceptance criteria registry |
| [`99-consistency-report.md`](../../99-consistency-report.md) | Consistency report |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/05-rust/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-RUST-001: Rust Coding Guidelines Index & Architecture Conformance

**Given** Rust source code, crate manifests, and specifications across modules and subsystems.
**When** Guideline linters audit the codebase against Rust coding standards and architecture rules.
**Then** All Rust specifications and source files adhere to RFC 430 conventions, PascalCase boundaries, structured error handling, async patterns, and memory safety rules with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.
