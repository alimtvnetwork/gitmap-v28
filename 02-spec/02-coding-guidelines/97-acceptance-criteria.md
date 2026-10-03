# Coding Guidelines — Master Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Establish a unified, master registry of testable acceptance criteria across all coding guideline domains in `02-spec/02-coding-guidelines/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-[CATEGORY]-[NUM]`), ensure every criterion maps to a concrete specification file, and verify compliance via automated linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `02-spec/02-coding-guidelines/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths across all guideline files.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-03
**Status:** Active
**AI Confidence:** Production-Ready

---

## Overview

Consolidated master registry index of testable criteria across all guideline categories, referencing category-specific acceptance criteria registries.

---

## AC-00: Root Style and Sizing Guidelines

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-ROOT-001` | Coding Guideline Authoring Standard Conformance | [`01-specification-and-coding-guideline-standard.md`](./01-specification-and-coding-guideline-standard.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |
| `AC-CG-ROOT-002` | Canonical Size Tier Enforcement (Function <= 15 lines, File <= 300 lines) | [`02-canonical-size-tier.md`](./02-canonical-size-tier.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |
| `AC-CG-ROOT-003` | Core Coding Style and Parameter Rules (Max 3 params, option structs, UTF-8 LF) | [`03-coding-style-checklist.md`](./03-coding-style-checklist.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |
| `AC-CG-ROOT-004` | Condensed Review Guide Rules (Zero nesting, no positive/negative mix) | [`04-consolidated-review-guide-condensed.md`](./04-consolidated-review-guide-condensed.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |
| `AC-CG-ROOT-005` | Consolidated Master Review Rules (Complete cross-language rule synthesis) | [`05-consolidated-review-guide.md`](./05-consolidated-review-guide.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |

---

## AC-01: Cross-Language Standards Registry

See [`01-cross-language/97-acceptance-criteria.md`](./01-cross-language/97-acceptance-criteria.md) for the complete 40+ testable criteria inventory across Boolean principles, code style, naming conventions, type safety, and architecture.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-BOOL-001` | Boolean principles define naming (`isX`, `hasX`) and affirmative implicit evaluation patterns | [`01-cross-language/02-boolean-principles/readme.md`](./01-cross-language/02-boolean-principles/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-TYPE-004` | Casting elimination patterns cover type-safe alternatives to type assertions | [`01-cross-language/04-casting-elimination-patterns.md`](./01-cross-language/04-casting-elimination-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-STYLE-001` | Code style defines formatting, naming, vertical spacing, and structural conventions | [`01-cross-language/04-code-style/readme.md`](./01-cross-language/04-code-style/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-STYLE-002` | Zero nesting, guard clauses, and early return patterns | [`01-cross-language/04-code-style/02-braces-and-nesting.md`](./01-cross-language/04-code-style/02-braces-and-nesting.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-STYLE-008` | Comprehensive Code Style Checklist | [`01-cross-language/04-code-style/08-checklist.md`](./01-cross-language/04-code-style/08-checklist.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-STYLE-009` | Vertical Newline Spacing & Whitespace Conformance | [`01-cross-language/21-newline-styling-examples.md`](./01-cross-language/21-newline-styling-examples.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-ARCH-008` | DRY principles documented with refactoring patterns and modular extraction | [`01-cross-language/08-dry-principles.md`](./01-cross-language/08-dry-principles.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-ARCH-006` | Cyclomatic complexity limits defined with enforcement rules (<= 10) | [`01-cross-language/06-cyclomatic-complexity.md`](./01-cross-language/06-cyclomatic-complexity.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-NAME-022` | Variable and collection naming conventions (camelCase, plural arrays, map prefixes) | [`01-cross-language/22-variable-naming-conventions.md`](./01-cross-language/22-variable-naming-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |
| `AC-CG-REG-001` | Detailed Cross-Language Criteria Registry (BOOL, STYLE, NAME, TYPE, ARCH, TEST, STATIC) | [`01-cross-language/97-acceptance-criteria.md`](./01-cross-language/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |

---

## AC-02: TypeScript Standards Registry

See [`02-typescript/97-acceptance-criteria.md`](./02-typescript/97-acceptance-criteria.md) for the complete 16 testable criteria inventory across TypeScript status enums, type-safety plans, discriminated unions, and ESLint enforcement.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-TS-000 | TypeScript Guidelines Index Conformance | [`02-typescript/readme.md`](./02-typescript/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-002 | ConnectionStatusEnum Definition and Validation | [`02-typescript/02-connection-status-enum.md`](./02-typescript/02-connection-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-003 | EntityStatusEnum Definition and Validation | [`02-typescript/03-entity-status-enum.md`](./02-typescript/03-entity-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-004 | ExecutionStatusEnum Definition and Validation | [`02-typescript/04-execution-status-enum.md`](./02-typescript/04-execution-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-005 | ExportStatusEnum Definition and Validation | [`02-typescript/05-export-status-enum.md`](./02-typescript/05-export-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-006 | HttpMethodEnum Definition and Validation | [`02-typescript/06-http-method-enum.md`](./02-typescript/06-http-method-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-007 | MessageStatusEnum Definition and Validation | [`02-typescript/07-message-status-enum.md`](./02-typescript/07-message-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-008 | TypeScript Type Safety Remediation Plan and Any Elimination | [`02-typescript/08-type-safety-remediation-plan.md`](./02-typescript/08-type-safety-remediation-plan.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-009 | TypeScript Standards Reference and Architectural Patterns | [`02-typescript/09-typescript-standards-reference.md`](./02-typescript/09-typescript-standards-reference.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-010 | TypeScript Promise, Async/Await, and Concurrency Patterns | [`02-typescript/10-promise-await-patterns.md`](./02-typescript/10-promise-await-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-011 | LogLevelEnum Definition and Validation | [`02-typescript/11-log-level-enum.md`](./02-typescript/11-log-level-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-012 | TypeScript ESLint Rules, Type Checking, and Lint Automation | [`02-typescript/12-eslint-enforcement.md`](./02-typescript/12-eslint-enforcement.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-013 | TypeScript Discriminated Unions and Exhaustive Type Narrowing | [`02-typescript/13-discriminated-union-patterns.md`](./02-typescript/13-discriminated-union-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-014 | TypeScript Enum Runtime Validation and Parse Guard Utilities | [`02-typescript/14-enum-checking-and-validation.md`](./02-typescript/14-enum-checking-and-validation.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-015 | TypeScript State Management, Stores, and Reactive Architecture | [`02-typescript/15-state-management.md`](./02-typescript/15-state-management.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| AC-CG-TS-REG-001 | TypeScript Acceptance Criteria Registry Conformance | [`02-typescript/97-acceptance-criteria.md`](./02-typescript/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |

---

## AC-03: Golang Standards Registry

See [`03-golang/97-acceptance-criteria.md`](./03-golang/97-acceptance-criteria.md) for the complete criteria inventory across Go standards, enums, and architecture.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-GO-000 | Golang Standards Index Conformance | [`03-golang/readme.md`](./03-golang/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-001 | Positive boolean naming and evaluation patterns | [`03-golang/02-boolean-standards.md`](./03-golang/02-boolean-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-003 | HttpMethod enum standard and string literal bans | [`03-golang/03-httpmethod-enum.md`](./03-golang/03-httpmethod-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-005 | Resource defer rules and loop safety | [`03-golang/05-defer-rules.md`](./03-golang/05-defer-rules.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-006 | Go string and slice internals efficiency | [`03-golang/06-string-slice-internals.md`](./03-golang/06-string-slice-internals.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-007 | Code severity taxonomy and fault logging | [`03-golang/07-code-severity-taxonomy.md`](./03-golang/07-code-severity-taxonomy.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-008 | Unified PathUtil and FileUtil specification | [`03-golang/08-pathutil-fileutil-spec.md`](./03-golang/08-pathutil-fileutil-spec.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-009 | Wrapped boolean results and monadic error returns | [`03-golang/09-wrapped-boolean-results.md`](./03-golang/09-wrapped-boolean-results.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-ENUM-000 | Go Enum Specification Registry & Patterns | [`03-golang/01-enum-specification/readme.md`](./03-golang/01-enum-specification/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-REF-000 | Go Standards Reference (Sizing, Types, DB, Naming, Concurrency) | [`03-golang/04-golang-standards-reference/readme.md`](./03-golang/04-golang-standards-reference/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| AC-CG-GO-REG-001 | Go Acceptance Criteria Registry Conformance | [`03-golang/97-acceptance-criteria.md`](./03-golang/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |

---

## AC-04: PHP Standards Registry

See [`04-php/97-acceptance-criteria.md`](./04-php/97-acceptance-criteria.md) for the complete criteria inventory across PHP backed enums, response arrays, PSR-4 naming, and standards reference.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-PHP-001 | PHP Standards Index & Overview Conformance | [`04-php/readme.md`](./04-php/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-002 | PHP Enums and Backed Enum Standards | [`04-php/02-enums.md`](./04-php/02-enums.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-003 | PHP Forbidden Anti-Patterns and Quality Gates | [`04-php/03-forbidden-patterns.md`](./04-php/03-forbidden-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-004 | PHP Naming Conventions, Identifiers, and Casings | [`04-php/04-naming-conventions.md`](./04-php/04-naming-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-005 | PHP Response Array Standards and Result Envelopes | [`04-php/05-response-array-standard.md`](./04-php/05-response-array-standard.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-006 | PHP Standards Reference Conformance | [`04-php/06-php-standards-reference.md`](./04-php/06-php-standards-reference.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-007 | PHP Spacing, Blank Lines and Import Hygiene | [`04-php/07-spacing-and-imports.md`](./04-php/07-spacing-and-imports.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-008 | PHP Response Key Type Inventory & Taxonomy | [`04-php/08-response-key-type-inventory.md`](./04-php/08-response-key-type-inventory.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-009 | PHP-Go Cross-Language Architecture Consistency | [`04-php/09-php-go-consistency-audit.md`](./04-php/09-php-go-consistency-audit.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-REF-000 | PHP Standards Reference Index Conformance | [`04-php/07-php-standards-reference/readme.md`](./04-php/07-php-standards-reference/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-REF-002 | PHP Modular Naming and Structured Error Envelopes | [`04-php/07-php-standards-reference/02-naming-and-errors.md`](./04-php/07-php-standards-reference/02-naming-and-errors.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-REF-003 | PHP Centralized Constants and Backed Enums | [`04-php/07-php-standards-reference/03-constants-and-deps.md`](./04-php/07-php-standards-reference/03-constants-and-deps.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-REF-004 | PHP Constructor Initialization, Positive Booleans, and isDefined Guards | [`04-php/07-php-standards-reference/04-initialization-and-booleans.md`](./04-php/07-php-standards-reference/04-initialization-and-booleans.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-REF-005 | PHP PSR-12 Code Style, Braces, Vertical Spacing, and Function Size Limits | [`04-php/07-php-standards-reference/05-code-style.md`](./04-php/07-php-standards-reference/05-code-style.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-REF-006 | PHP Forbidden Globals, Sanitization, and Prepared Database Queries | [`04-php/07-php-standards-reference/06-forbidden-and-database.md`](./04-php/07-php-standards-reference/06-forbidden-and-database.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| AC-CG-PHP-REG-001 | PHP Acceptance Criteria Registry Conformance | [`04-php/97-acceptance-criteria.md`](./04-php/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |

---

## AC-05: Rust Standards Registry

See [`05-rust/97-acceptance-criteria.md`](./05-rust/97-acceptance-criteria.md) for the complete criteria inventory across Rust naming, error handling, async Tokio, memory safety, and FFI boundaries.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-RUST-001 | Rust Coding Guidelines Index & Architecture Conformance | [`05-rust/readme.md`](./05-rust/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| AC-CG-RUST-002 | Rust Naming Conventions and Positive Booleans | [`05-rust/02-naming-conventions.md`](./05-rust/02-naming-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| AC-CG-RUST-003 | Rust Result Handling, Error Trait and No Panics | [`05-rust/03-error-handling.md`](./05-rust/03-error-handling.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| AC-CG-RUST-004 | Rust Async Tokio Cancellation and Send Bounds | [`05-rust/04-async-patterns.md`](./05-rust/04-async-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| AC-CG-RUST-005 | Rust Zero Unsafe and RAII Lifetime Safety | [`05-rust/05-memory-safety.md`](./05-rust/05-memory-safety.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| AC-CG-RUST-006 | Rust Test Suite Organization and Mocking Safety | [`05-rust/06-testing-standards.md`](./05-rust/06-testing-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| AC-CG-RUST-007 | Rust FFI Safety Boundaries and Platform Isolation | [`05-rust/07-ffi-platform.md`](./05-rust/07-ffi-platform.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| AC-CG-RUST-REG-001 | Rust Acceptance Criteria Registry Conformance | [`05-rust/97-acceptance-criteria.md`](./05-rust/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |

---

## AC-06: AI Optimization Registry

See [`06-ai-optimization/97-acceptance-criteria.md`](./06-ai-optimization/97-acceptance-criteria.md) for criteria covering anti-hallucination, citation requirements, quick reference checklists, and agent memory lifecycles.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-AI-000 | AI Optimization Guidelines Index Conformance | [`06-ai-optimization/readme.md`](./06-ai-optimization/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| AC-CG-AI-002 | Anti-hallucination verification and source attribution | [`06-ai-optimization/02-anti-hallucination-rules.md`](./06-ai-optimization/02-anti-hallucination-rules.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| AC-CG-AI-006 | Strict relative path citation and markdown link verification | [`06-ai-optimization/06-citation-requirement.md`](./06-ai-optimization/06-citation-requirement.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| AC-CG-AI-009 | Agent memory lifecycle, persistent knowledge, and index updates | [`06-ai-optimization/09-agent-memory-lifecycle.md`](./06-ai-optimization/09-agent-memory-lifecycle.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| AC-CG-AI-REG-001 | AI Optimization Acceptance Criteria Registry Conformance | [`06-ai-optimization/97-acceptance-criteria.md`](./06-ai-optimization/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |

---

## AC-06B: CI/CD Integration Registry

See [`06-cicd-integration/97-acceptance-criteria.md`](./06-cicd-integration/97-acceptance-criteria.md) and [`06-cicd-integration/08-fix-repo-and-installers/97-acceptance-criteria.md`](./06-cicd-integration/08-fix-repo-and-installers/97-acceptance-criteria.md) for SARIF, plugin contracts, fix repo automation, and quality gates.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-CICD-001 | CI/CD Integration Architecture Conformance (alias `AC-CG-CI-001`) | [`06-cicd-integration/readme.md`](./06-cicd-integration/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-002 | SARIF output contract and schema compliance (alias `AC-CG-CI-002`) | [`06-cicd-integration/02-sarif-contract.md`](./06-cicd-integration/02-sarif-contract.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-003 | Language Plugin Architecture and Registry Standards (alias `AC-CG-CI-003`) | [`06-cicd-integration/03-plugin-model.md`](./06-cicd-integration/03-plugin-model.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-004 | Language Rollout Roadmap and Promotion Standards (alias `AC-CG-CI-004`) | [`06-cicd-integration/04-language-roadmap.md`](./06-cicd-integration/04-language-roadmap.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-005 | Cross-Platform CI Templates and Invocation Standards (alias `AC-CG-CI-005`) | [`06-cicd-integration/05-ci-templates.md`](./06-cicd-integration/05-ci-templates.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-006 | Distribution Packaging and Release Asset Governance (alias `AC-CG-CI-006`) | [`06-cicd-integration/06-distribution.md`](./06-cicd-integration/06-distribution.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-007 | Rule Taxonomy, Severity Mapping, and Tier Coordination (alias `AC-CG-CI-007`) | [`06-cicd-integration/07-rules-mapping.md`](./06-cicd-integration/07-rules-mapping.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-008 | Probe Ordering, Parallelism, and Timeout Budgets (alias `AC-CG-CI-008`) | [`06-cicd-integration/08-performance.md`](./06-cicd-integration/08-performance.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-FIX-000 | Fix Repo & Installer Specifications Registry | [`06-cicd-integration/08-fix-repo-and-installers/readme.md`](./06-cicd-integration/08-fix-repo-and-installers/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-FAQ-001 | Linter Pack FAQ & Consumer Operations Conformance | [`06-cicd-integration/98-faq.md`](./06-cicd-integration/98-faq.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-TROUBLE-001 | Linter Pack Operations Troubleshooting Conformance | [`06-cicd-integration/99-troubleshooting.md`](./06-cicd-integration/99-troubleshooting.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |
| AC-CG-CICD-REG-001 | CI/CD Acceptance Criteria Registry Conformance | [`06-cicd-integration/97-acceptance-criteria.md`](./06-cicd-integration/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only` |

---

## AC-07: C# Standards Registry

See [`07-csharp/97-acceptance-criteria.md`](./07-csharp/97-acceptance-criteria.md) for C# naming conventions, method design, error handling, and type safety standards.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-CS-001 | C# Coding Guidelines Index & Module Conformance | [`07-csharp/readme.md`](./07-csharp/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| AC-CG-CS-002 | C# PascalCase/camelCase and Positive Booleans | [`07-csharp/02-naming-and-conventions.md`](./07-csharp/02-naming-and-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| AC-CG-CS-003 | C# Method Signatures, Parameter Limits and Guard Clauses | [`07-csharp/03-method-design.md`](./07-csharp/03-method-design.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| AC-CG-CS-004 | C# Structured Exceptions, Specific Catches and Result Types | [`07-csharp/04-error-handling.md`](./07-csharp/04-error-handling.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| AC-CG-CS-005 | C# Nullable Reference Types and Pattern Matching | [`07-csharp/05-type-safety.md`](./07-csharp/05-type-safety.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |
| AC-CG-CS-REG-001 | C# Acceptance Criteria Registry Conformance | [`07-csharp/97-acceptance-criteria.md`](./07-csharp/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/07-csharp --check-only` |

---

## AC-08: File & Folder Naming Registry

See [`08-file-folder-naming/97-acceptance-criteria.md`](./08-file-folder-naming/97-acceptance-criteria.md) for cross-language kebab-case and lowercase directory hygiene.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-FILE-000 | File & Folder Naming Directory Index Conformance | [`08-file-folder-naming/readme.md`](./08-file-folder-naming/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| AC-CG-FILE-002 | Cross-language lowercase kebab-case naming standard | [`08-file-folder-naming/02-cross-language.md`](./08-file-folder-naming/02-cross-language.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| AC-CG-FILE-REG-001 | File Naming Acceptance Criteria Registry Conformance | [`08-file-folder-naming/97-acceptance-criteria.md`](./08-file-folder-naming/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |

---

## AC-09: PowerShell Integration Registry

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-PWSH-000 | PowerShell Integration Specification Conformance | [`09-powershell-integration/readme.md`](./09-powershell-integration/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/09-powershell-integration --check-only` |

---

## AC-10: Research Standards Registry

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-RES-000 | Research Directory Specification Conformance | [`10-research/readme.md`](./10-research/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/10-research --check-only` |

---

## AC-11: Security Guidelines Registry

See [`11-security/97-acceptance-criteria.md`](./11-security/97-acceptance-criteria.md) for JWT lifecycles, encryption standards, OWASP mitigation, secret vaulting, and dependency pinning.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-SEC-000 | Security Guidelines Directory Index & Policy Conformance | [`11-security/readme.md`](./11-security/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-001 | JWT Token Lifecycle & HttpOnly Cookie Storage | [`11-security/02-jwt-standards.md`](./11-security/02-jwt-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-002 | Encryption at Rest & Key Derivation Standards | [`11-security/03-encryption-standards.md`](./11-security/03-encryption-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-003 | OWASP Top 10 Mitigation Controls & Input Validation | [`11-security/04-owasp-top-10.md`](./11-security/04-owasp-top-10.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-004 | Zero Secrets in Source Control & Environment Vaulting | [`11-security/05-secret-management.md`](./11-security/05-secret-management.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-AXIOS-000 | Axios Client Security Overview | [`11-security/01-axios-version-control/readme.md`](./11-security/01-axios-version-control/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-AXIOS-001 | Strict Axios Version Pinning & Dependency Locking | [`11-security/01-axios-version-control/02-implementation-rules.md`](./11-security/01-axios-version-control/02-implementation-rules.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-AXIOS-002 | CVE Remediation & Supply Chain Security Verification | [`11-security/01-axios-version-control/03-security-notes.md`](./11-security/01-axios-version-control/03-security-notes.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |
| AC-CG-SEC-REG-001 | Security Acceptance Criteria Registry Conformance | [`11-security/97-acceptance-criteria.md`](./11-security/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security --check-only` |

---

## AC-12: Python Standards Registry

See [`12-python/97-acceptance-criteria.md`](./12-python/97-acceptance-criteria.md) for the complete criteria inventory across Python static typing, Pydantic data validation, PEP-8 compliance, and specific exceptions.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-PY-001 | Python Guidelines Directory Index Conformance | [`12-python/readme.md`](./12-python/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| AC-CG-PY-002 | Python Coding Standards Conformance | [`12-python/02-standards.md`](./12-python/02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| AC-CG-PY-003 | Python Dynamic Enum & Array Constants Standard | [`12-python/02-standards.md`](./12-python/02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| AC-CG-PY-004 | Python DRY Architecture & Engine Caching Conformance | [`12-python/02-standards.md`](./12-python/02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| AC-CG-PY-REG-001 | Python Acceptance Criteria Registry Conformance | [`12-python/97-acceptance-criteria.md`](./12-python/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |

---

## AC-13: Modern C++ Standards Registry

See [`13-cpp/97-acceptance-criteria.md`](./13-cpp/97-acceptance-criteria.md) for the complete criteria inventory across C++20 baseline, concepts, RAII smart pointers, Rule of Zero/Five, and FFI exception boundaries.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-CPP-001 | Modern C++ Guidelines Directory Index Conformance | [`13-cpp/readme.md`](./13-cpp/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| AC-CG-CPP-002 | Modern C++ Standards Conformance | [`13-cpp/02-standards.md`](./13-cpp/02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| AC-CG-CPP-003 | C++ Memory Safety & RAII Resource Management | [`13-cpp/02-standards.md`](./13-cpp/02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| AC-CG-CPP-004 | C++ FFI Boundary Exception Safety & Standard Types | [`13-cpp/02-standards.md`](./13-cpp/02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| AC-CG-CPP-REG-001 | Modern C++ Acceptance Criteria Registry Conformance | [`13-cpp/97-acceptance-criteria.md`](./13-cpp/97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |

---

## AC-14: Application Architecture & Module Readmes

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-APP-001 | Application Specifications Index & Navigation | [`21-app/readme.md`](./21-app/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/21-app --check-only` |
| AC-CG-APP-002 | Non-CI/CD Application Issues & Bug Catalog Index | [`22-app-issues/readme.md`](./22-app-issues/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/22-app-issues --check-only` |
| AC-CG-APP-003 | Application Database Standards & Split SQLite Architecture | [`23-app-db/readme.md`](./23-app-db/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/23-app-db --check-only` |
| AC-CG-APP-004 | Application UI/UX Design System Specifications | [`24-app-ui-design-system/readme.md`](./24-app-ui-design-system/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/24-app-ui-design-system --check-only` |
| AC-CG-CONSISTENCY-001 | Guideline Consistency & Structure Conformance | [`../99-consistency-report.md`](../99-consistency-report.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |

---

## 16. Application Specifications (Coding Guidelines Architecture) (AC-15)

See [`../21-app/07-coding-guideline-actionable-checklist-and-acceptance/readme.md`](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/readme.md) for the complete coding guidelines actionable checklist, 4-part anatomy architecture, and automated verification engine.

| ID | Title | Source / Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| AC-CG-SPEC-000 | Architecture Specification Module Conformance | [`../21-app/07-coding-guideline-actionable-checklist-and-acceptance/readme.md`](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/21-app/07-coding-guideline-actionable-checklist-and-acceptance --check-only` |
| AC-CG-SPEC-001 | Coding Guidelines Architecture Specification Conformance | [`../21-app/07-coding-guideline-actionable-checklist-and-acceptance/01-architecture-spec.md`](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/01-architecture-spec.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |
| AC-CG-REG-MASTER-001 | Master Registry Completeness & Polyglot Parity | [`../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md`](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |
| AC-CG-REG-PY-001 | Python Acceptance Criteria Registry Standardization | [`../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md`](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| AC-CG-REG-CPP-001 | Modern C++ Acceptance Criteria Registry Standardization | [`../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md`](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| AC-CG-REG-TOOL-001 | Automated Guideline Autofixer & CI/CD Runner Verification | [`../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md`](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/02-component-spec.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` |

### AC-CG-SPEC-000: Architecture Specification Module Conformance

**Given** The specification files in `02-spec/21-app/07-coding-guideline-actionable-checklist-and-acceptance/`.
**When** Audited against Prompt Architect specification rules.
**Then** All files feature valid AI execution headers, actionable checklists, and testable acceptance criteria with 0 absolute paths.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/21-app/07-coding-guideline-actionable-checklist-and-acceptance --check-only
```
**Expected:** exit 0. Zero violations.

### AC-CG-SPEC-001: Coding Guidelines Architecture Specification Conformance

**Given** The coding guidelines standard specification in `02-spec/21-app/07-coding-guideline-actionable-checklist-and-acceptance/01-architecture-spec.md`.
**When** Audited by repository linters and guideline autofixers.
**Then** The file strictly includes the AI Execution Prompt header (`> **/goal**` and `> **/learn**`), the actionable checklist (`## 🎯 Actionable CI/CD & Agent Checklist`), the verbatim user prompt, complete 4-part anatomy, detailed style, boolean, and naming rules, and concludes with automated testable criteria.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations detected.

### AC-CG-REG-MASTER-001: Master Registry Completeness & Polyglot Parity

**Given** The master acceptance criteria registry at `02-spec/02-coding-guidelines/97-acceptance-criteria.md`.
**When** Audited for coverage across all guideline subdirectories in `02-spec/02-coding-guidelines/`.
**Then** All twelve domain registries (Cross-Language, TypeScript, Golang, PHP, Rust, AI Optimization, CI/CD Integration, C#, File Naming, Security, Python, and C++) are indexed with accurate relative links and summary tables.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations detected.

### AC-CG-REG-PY-001: Python Acceptance Criteria Registry Standardization

**Given** The Python guideline directory at `02-spec/02-coding-guidelines/12-python/`.
**When** Checked for registry presence and criterion completeness.
**Then** `12-python/97-acceptance-criteria.md` exists, catalogs criteria `AC-CG-PY-001` through `AC-CG-PY-REG-001` with explicit `Given/When/Then` blocks, and defines runnable verification commands.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations detected.

### AC-CG-REG-CPP-001: Modern C++ Acceptance Criteria Registry Standardization

**Given** The C++ guideline directory at `02-spec/02-coding-guidelines/13-cpp/`.
**When** Checked for registry presence and criterion completeness.
**Then** `13-cpp/97-acceptance-criteria.md` exists, catalogs criteria `AC-CG-CPP-001` through `AC-CG-CPP-REG-001` with explicit `Given/When/Then` blocks, and defines runnable verification commands.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only
```
**Expected:** exit 0. Zero violations detected.

### AC-CG-REG-TOOL-001: Automated Guideline Autofixer & CI/CD Runner Verification

**Given** The complete set of specification files across `02-spec/02-coding-guidelines/`.
**When** The automated verification engine is executed in audit mode.
**Then** The guideline autofixer passes with zero errors, zero formatting violations, and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations detected.

---

## Cross-References

- [Coding Guidelines Overview](./readme.md)
- [Coding Guideline Authoring Standard](./01-specification-and-coding-guideline-standard.md)
- [Cross-Language Standards](./01-cross-language/readme.md)
- [Cross-Language Acceptance Criteria Registry](./01-cross-language/97-acceptance-criteria.md)
- [TypeScript Standards](./02-typescript/readme.md)
- [TypeScript Acceptance Criteria Registry](./02-typescript/97-acceptance-criteria.md)
- [Golang Standards](./03-golang/readme.md)
- [Golang Acceptance Criteria Registry](./03-golang/97-acceptance-criteria.md)
- [PHP Standards](./04-php/readme.md)
- [PHP Acceptance Criteria Registry](./04-php/97-acceptance-criteria.md)
- [Rust Standards](./05-rust/readme.md)
- [Rust Acceptance Criteria Registry](./05-rust/97-acceptance-criteria.md)
- [AI Optimization Standards](./06-ai-optimization/readme.md)
- [AI Optimization Acceptance Criteria Registry](./06-ai-optimization/97-acceptance-criteria.md)
- [CI/CD Integration Standards](./06-cicd-integration/readme.md)
- [CI/CD Acceptance Criteria Registry](./06-cicd-integration/97-acceptance-criteria.md)
- [C# Standards](./07-csharp/readme.md)
- [C# Acceptance Criteria Registry](./07-csharp/97-acceptance-criteria.md)
- [File & Folder Naming](./08-file-folder-naming/readme.md)
- [File & Folder Naming Acceptance Criteria Registry](./08-file-folder-naming/97-acceptance-criteria.md)
- [PowerShell Integration Standards](./09-powershell-integration/readme.md)
- [Research Standards](./10-research/readme.md)
- [Security Guidelines](./11-security/readme.md)
- [Security Acceptance Criteria Registry](./11-security/97-acceptance-criteria.md)
- [Python Standards](./12-python/readme.md)
- [Python Acceptance Criteria Registry](./12-python/97-acceptance-criteria.md)
- [Modern C++ Standards](./13-cpp/readme.md)
- [Modern C++ Acceptance Criteria Registry](./13-cpp/97-acceptance-criteria.md)
- [Coding Guidelines Actionable Checklist & Acceptance Criteria Architecture](../21-app/07-coding-guideline-actionable-checklist-and-acceptance/readme.md)

---

## Verification & Acceptance Criteria

### AC-CG-MASTER-REG-001: Master Acceptance Criteria Registry Conformance

**Given** The consolidated master acceptance criteria registry in `02-spec/02-coding-guidelines/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All criteria entries correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
