# Rust Coding Standards

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
| `00-overview.md` present | ✅ |
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

See 01-naming-conventions.md for the complete reference with examples.
```

---

## Purpose

Rust-specific coding standards for the Time Log CLI and any future Rust-based projects. Extends the [Cross-Language Guidelines](../01-cross-language/00-overview.md) but **overrides the naming convention** to follow Rust community standards (RFC 430) with two explicit PascalCase exceptions at system boundaries.

This override exists because Rust's compiler actively enforces `snake_case` for functions/variables via lint warnings, making the project-wide PascalCase mandate impractical for Rust code. The two exceptions (database and enum strings) are where Rust code interacts with other systems in the stack that expect PascalCase.

---

## Document Inventory

| File | Description |
|------|-------------|
| 01-naming-conventions.md | Rust naming rules: snake_case default, PascalCase for DB + enum strings, serialization, module structure |
| 02-error-handling.md | Error types, Result patterns, thiserror/anyhow usage |
| 03-async-patterns.md | Tokio async conventions, channel patterns, cancellation |
| 04-memory-safety.md | Ownership idioms, lifetime rules, unsafe policy |
| 05-testing-standards.md | Unit/integration test structure, mocking, property testing |
| 06-ffi-platform.md | FFI safety rules, conditional compilation, platform abstractions |
| 97-acceptance-criteria.md | Compliance requirements |
| 98-changelog.md | Version history |
| 99-consistency-report.md | Structural health |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Guidelines | `../01-cross-language/00-overview.md` |
| Coding Guidelines Root | `../00-overview.md` |
| Database Conventions (PascalCase) | `../../../04-database-conventions/00-overview.md` |
| Enum Standards (Cross-Language) | `../../../../17-consolidated-guidelines/04-enum-standards.md` |
| 01-naming-conventions.md | — |
| 02-error-handling.md | — |
| 03-async-patterns.md | — |
| 04-memory-safety.md | — |
| 05-testing-standards.md | — |
| 06-ffi-platform.md | — |
| 97-acceptance-criteria.md | — |
| 98-changelog.md | — |
| 99-consistency-report.md | — |
