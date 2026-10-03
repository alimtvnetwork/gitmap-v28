# PHP Standards — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all PHP coding guideline specifications in `04-php/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-PHP-[NUM]`), native backed enums with `Type` suffix, complete elimination of forbidden patterns, PSR-12 code style, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `04-php/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. PHP Criteria Inventory (`AC-CG-PHP-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-PHP-001` | PHP Standards Index & Overview Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-002` | PHP Enums and Backed Enum Standards | [`02-enums.md`](02-enums.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-003` | PHP Forbidden Anti-Patterns and Quality Gates | [`03-forbidden-patterns.md`](03-forbidden-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-004` | PHP Naming Conventions, Identifiers, and Casings | [`04-naming-conventions.md`](04-naming-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-005` | PHP Response Array Standards and Result Envelopes | [`05-response-array-standard.md`](05-response-array-standard.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-006` | PHP Standards Reference Conformance | [`06-php-standards-reference.md`](06-php-standards-reference.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-007` | PHP Spacing, Blank Lines and Import Hygiene | [`07-spacing-and-imports.md`](07-spacing-and-imports.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-008` | PHP Response Key Type Inventory & Taxonomy | [`08-response-key-type-inventory.md`](08-response-key-type-inventory.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-009` | PHP-Go Cross-Language Architecture Consistency | [`09-php-go-consistency-audit.md`](09-php-go-consistency-audit.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-REG-001` | PHP Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-REF-000` | PHP Standards Reference Index Conformance | [`07-php-standards-reference/readme.md`](07-php-standards-reference/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-REF-002` | PHP Modular Naming and Structured Error Envelopes | [`07-php-standards-reference/02-naming-and-errors.md`](07-php-standards-reference/02-naming-and-errors.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-REF-003` | PHP Centralized Constants and Backed Enums | [`07-php-standards-reference/03-constants-and-deps.md`](07-php-standards-reference/03-constants-and-deps.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-REF-004` | PHP Constructor Initialization, Positive Booleans, and isDefined Guards | [`07-php-standards-reference/04-initialization-and-booleans.md`](07-php-standards-reference/04-initialization-and-booleans.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-REF-005` | PHP PSR-12 Code Style, Braces, Vertical Spacing, and Function Size Limits | [`07-php-standards-reference/05-code-style.md`](07-php-standards-reference/05-code-style.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |
| `AC-CG-PHP-REF-006` | PHP Forbidden Globals, Sanitization, and Prepared Database Queries | [`07-php-standards-reference/06-forbidden-and-database.md`](07-php-standards-reference/06-forbidden-and-database.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-PHP-001: PHP Standards Index & Overview Conformance

**Given** PHP development documentation and guidelines.
**When** Specifications in `04-php/` are audited against directory conventions and structural guidelines.
**Then** All PHP guideline documents, index inventories, and cross-references are complete with zero orphaned files and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-002: PHP Enums and Backed Enum Standards

**Given** PHP source code defining enums, constants, and typed values.
**When** Codebases are audited against native backed enum standards.
**Then** All enum classes use PHP 8.1+ native string-backed enums with `Type` suffix, PascalCase cases, `isEqual()` methods, and zero magic strings with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-003: PHP Forbidden Anti-Patterns and Quality Gates

**Given** PHP source code implementing business logic, controllers, and helpers.
**When** Source files are scanned for prohibited language features and anti-patterns.
**Then** Zero instances of `eval()`, direct superglobal access (`$_POST`, `$_GET`), `extract()`, `goto`, raw SQL string concatenation, or unescaped queries exist, with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-004: PHP Naming Conventions, Identifiers, and Casings

**Given** PHP declarations including classes, interfaces, traits, methods, functions, and variables.
**When** Identifiers are audited against repository naming standards.
**Then** Classes, traits, interfaces, and enums strictly use PascalCase, methods and variables use camelCase, log context keys use camelCase, and zero unapproved snake_case identifiers exist, with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-005: PHP Response Array Standards and Result Envelopes

**Given** Internal PHP service methods returning operation outcomes and error payloads.
**When** Service results are constructed across traits and providers.
**Then** All service methods strictly utilize `ResultHelper` factory methods (`ok`, `failed`, `error`, `errorWithCode`, `errorFromException`) with `ResponseKeyType` enum keys, multiline array formatting, and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-006: PHP Standards Reference Conformance

**Given** PHP implementation files and architectural references.
**When** Audited against the modularized standards reference suite.
**Then** All code conforms to the decomposed standards for naming, centralized constants, explicit property initialization, positive booleans, PSR-12 code style, and database query wrapping with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-007: PHP Spacing, Blank Lines and Import Hygiene

**Given** PHP files in the `RiseupAsia` namespace.
**When** Source files are audited for layout and statement spacing.
**Then** All `if` and `throw` statements preceded by code have a mandatory blank line, all global classes and exceptions use file-level `use` imports without leading backslashes, and reusable log keys use enums with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-008: PHP Response Key Type Inventory & Taxonomy

**Given** PHP response payloads, dictionary keys, and result maps.
**When** Codebases are audited against the 176-case `ResponseKeyType` enum catalog.
**Then** All response keys strictly map to approved PascalCase `ResponseKeyType` enum cases, with zero magic string keys, correct helper method usage (`isEqual`, `isOtherThan`, `isAnyOf`), and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-009: PHP-Go Cross-Language Architecture Consistency

**Given** Shared domain logic, database tables, and API contracts between PHP and Go.
**When** Cross-language architecture is audited for schema and enum parity.
**Then** SQLite database schemas use identical PascalCase tables and columns, enum variants and label values mirror each other symmetrically, and API response keys remain synchronized with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-REF-000: PHP Standards Reference Index Conformance

**Given** PHP standards reference files and companion plugin implementations.
**When** Audited against `07-php-standards-reference/readme.md`.
**Then** All reference guidelines comply with repository architecture and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php/07-php-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-REF-002: PHP Modular Naming and Structured Error Envelopes

**Given** PHP classes, enums, functions, and error response builders.
**When** Audited against `07-php-standards-reference/02-naming-and-errors.md`.
**Then** Class names use PascalCase, methods use camelCase, error payloads use standard envelopes (`Success`, `Error`, `Message`), and raw error strings are prohibited with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-REF-003: PHP Centralized Constants and Backed Enums

**Given** PHP constants, hooks, capabilities, HTTP methods, and file paths.
**When** Audited against `07-php-standards-reference/03-constants-and-deps.md`.
**Then** All identifiers are defined centrally in `constants.php` or backed enums in `includes/Enums/`, with zero scattered hardcoded values and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-REF-004: PHP Constructor Initialization, Positive Booleans, and isDefined Guards

**Given** PHP class constructors, property definitions, and conditional branching.
**When** Audited against `07-php-standards-reference/04-initialization-and-booleans.md`.
**Then** All properties are explicitly initialized, booleans use positive `is`/`has` prefixes without explicit `=== true` comparisons, and guard clauses use `isDefined()` methods with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-REF-005: PHP PSR-12 Code Style, Braces, Vertical Spacing, and Function Size Limits

**Given** PHP statements, braces, conditionals, returns, and function declarations.
**When** Audited against `07-php-standards-reference/05-code-style.md`.
**Then** Braces follow PSR-12, blank lines precede `if` and `return`, functions remain strictly within the 8–15 line limit, and deeply nested code is flattened with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PHP-REF-006: PHP Forbidden Globals, Sanitization, and Prepared Database Queries

**Given** PHP database access, user input handling, and runtime execution.
**When** Audited against `07-php-standards-reference/06-forbidden-and-database.md`.
**Then** Superglobals are accessed through sanitization wrappers, all SQL queries use `$wpdb->prepare` or dedicated database helpers, and execution operators are banned with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.

---

## Cross-References

- [PHP Standards Overview](./readme.md)
- [PHP Enums Specification](./02-enums.md)
- [PHP Forbidden Patterns](./03-forbidden-patterns.md)
- [PHP Naming Conventions](./04-naming-conventions.md)
- [PHP Response Array Standard](./05-response-array-standard.md)
- [PHP Standards Reference](./06-php-standards-reference.md)
- [PHP Spacing & Imports](./07-spacing-and-imports.md)
- [Response Key Type Inventory](./08-response-key-type-inventory.md)
- [PHP-Go Consistency Audit](./09-php-go-consistency-audit.md)
- [PHP Standards Reference Suite](./07-php-standards-reference/readme.md)

---

## Verification & Acceptance Criteria

### AC-CG-PHP-REG-001: PHP Acceptance Criteria Registry Conformance

**Given** The PHP acceptance criteria registry in `02-spec/02-coding-guidelines/04-php/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All PHP specification criteria correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.
