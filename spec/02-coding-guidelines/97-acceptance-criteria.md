# Coding Guidelines — Acceptance Criteria

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

---

## Overview

30 testable criteria across 5 guideline categories, consolidating subfolder-level acceptance criteria.

---

## AC-01: Cross-Language Standards

| # | Criterion | Source |
|---|-----------|--------|
| AC-001 | Boolean principles define naming (`isX`, `hasX`, `canX`) and evaluation patterns | `01-cross-language/02-boolean-principles/00-overview.md` |
| AC-002 | Casting elimination patterns cover type-safe alternatives to type assertions | `01-cross-language/03-casting-elimination-patterns.md` |
| AC-003 | Code style defines formatting, naming, and structural conventions | `01-cross-language/04-code-style/00-overview.md` |
| AC-004 | All guidelines include ❌ (forbidden) and ✅ (compliant) code examples | `01-cross-language/15-master-coding-guidelines/00-overview.md` |
| AC-005 | DRY principles documented with refactoring patterns | `01-cross-language/08-dry-principles.md` |
| AC-006 | Cyclomatic complexity limits defined with enforcement rules | `01-cross-language/06-cyclomatic-complexity.md` |

---

## AC-02: TypeScript Standards

| # | Criterion | Source |
|---|-----------|--------|
| AC-007 | Connection status enums define all valid states with TypeScript string literals | `02-typescript/` |
| AC-008 | Type definitions avoid `any` and use proper generic constraints | `02-typescript/` |
| AC-009 | React component patterns follow functional component with hooks style | `02-typescript/` |
| AC-010 | State management patterns use Zustand stores with typed selectors | `02-typescript/` |

---

## AC-03: Golang Standards

| # | Criterion | Source |
|---|-----------|--------|
| AC-011 | Boolean standards define naming and evaluation patterns per Go idioms | `03-golang/02-boolean-standards.md` |
| AC-012 | Error handling uses `apperror.Result[T]` pattern consistently | `03-golang/` |
| AC-013 | HTTP method enum defines typed constants | `03-golang/03-httpmethod-enum.md` |
| AC-014 | Service layer follows interface-based dependency injection | `03-golang/` |

---

## AC-04: PHP Standards

| # | Criterion | Source |
|---|-----------|--------|
| AC-015 | Class naming follows WordPress PSR-4 autoloading conventions | `04-php/` |
| AC-016 | Database queries use $wpdb prepared statements exclusively | `04-php/` |
| AC-017 | Type declarations (parameter + return types) required on all functions | `04-php/` |
| AC-018 | Input sanitization and output escaping follow WordPress security standards | `04-php/` |

---

## AC-05: Rust Standards

| # | Criterion | Source |
|---|-----------|--------|
| AC-019 | Naming conventions follow Rust idioms (snake_case for functions, PascalCase for types) | `05-rust/` |
| AC-020 | Error handling uses `Result<T, E>` pattern with custom error types | `05-rust/` |
| AC-021 | Async patterns use tokio runtime with proper cancellation handling | `05-rust/` |
| AC-022 | Memory safety patterns documented for FFI boundaries | `05-rust/` |

---

## Cross-References

- [Overview](./00-overview.md)
- [Cross-Language Standards](./01-cross-language/00-overview.md)
- [TypeScript Standards](./02-typescript/00-overview.md)
- [Golang Standards](./03-golang/00-overview.md)
- [PHP Standards](./04-php/00-overview.md)
- [Rust Standards](./05-rust/00-overview.md)
