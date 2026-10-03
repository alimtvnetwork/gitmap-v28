# Golang Standards — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all Go coding guideline specifications in `03-golang/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-GO-[NUM]`), positive boolean rules (`is`/`has` only), `*appfault.AppError` return types, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `03-golang/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready

---

## 1. Golang Criteria Inventory (`AC-CG-GO-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-GO-000` | Golang Coding Standards Index & Specification Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |
| `AC-CG-GO-001` | Go Positive Boolean Naming, Negation Elimination, and Guard Clauses | [`02-boolean-standards.md`](02-boolean-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/02-boolean-standards.md --check-only` |
| `AC-CG-GO-003` | HttpMethod Enum Standard and Magic String Elimination | [`03-httpmethod-enum.md`](03-httpmethod-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/03-httpmethod-enum.md --check-only` |
| `AC-CG-GO-005` | Resource Defer Rules and Loop Safety Standards | [`05-defer-rules.md`](05-defer-rules.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/05-defer-rules.md --check-only` |
| `AC-CG-GO-006` | Go String and Slice Preallocation & Memory Efficiency Standards | [`06-string-slice-internals.md`](06-string-slice-internals.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/06-string-slice-internals.md --check-only` |
| `AC-CG-GO-007` | Code Severity Taxonomy, Error Classification, and Fault Logging | [`07-code-severity-taxonomy.md`](07-code-severity-taxonomy.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/07-code-severity-taxonomy.md --check-only` |
| `AC-CG-GO-008` | Unified PathUtil & FileUtil Cross-Platform Specification | [`08-pathutil-fileutil-spec.md`](08-pathutil-fileutil-spec.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/08-pathutil-fileutil-spec.md --check-only` |
| `AC-CG-GO-009` | Wrapped Result Monad and Single Return Parameter Standards | [`09-wrapped-boolean-results.md`](09-wrapped-boolean-results.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/09-wrapped-boolean-results.md --check-only` |
| `AC-CG-GO-REG-001` | Go Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only` |

---

## 2. Go Enum Specification Criteria (`AC-CG-GO-ENUM-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-GO-ENUM-000` | Go Enum Specification Index Conformance | [`01-enum-specification/readme.md`](01-enum-specification/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only` |
| `AC-CG-GO-ENUM-002` | Go Enum Type Definition and Constant Declarations | [`01-enum-specification/02-enum-pattern.md`](01-enum-specification/02-enum-pattern.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification/02-enum-pattern.md --check-only` |
| `AC-CG-GO-ENUM-003` | Go Enum Required Methods Implementation | [`01-enum-specification/03-required-methods.md`](01-enum-specification/03-required-methods.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification/03-required-methods.md --check-only` |
| `AC-CG-GO-ENUM-004` | Go Enum Architecture and Folder Layout | [`01-enum-specification/04-folder-structure.md`](01-enum-specification/04-folder-structure.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification/04-folder-structure.md --check-only` |
| `AC-CG-GO-ENUM-005` | Go Enum Validation and Compile-Time Verification | [`01-enum-specification/05-validation-checklist.md`](01-enum-specification/05-validation-checklist.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification/05-validation-checklist.md --check-only` |
| `AC-CG-GO-ENUM-006` | Go Enum Info Object Metadata Pattern | [`01-enum-specification/06-info-object-pattern.md`](01-enum-specification/06-info-object-pattern.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification/06-info-object-pattern.md --check-only` |

---

## 3. Go Standards Reference Criteria (`AC-CG-GO-REF-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-GO-REF-000` | Go Standards Reference Index Conformance | [`04-golang-standards-reference/readme.md`](04-golang-standards-reference/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only` |
| `AC-CG-GO-REF-002` | Go File Sizing, Function Limits, and Complexity Rules | [`04-golang-standards-reference/02-file-and-function-rules.md`](04-golang-standards-reference/02-file-and-function-rules.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/02-file-and-function-rules.md --check-only` |
| `AC-CG-GO-REF-003` | Go Type Safety, Generic Result Wrappers, and AppError Wrapping | [`04-golang-standards-reference/03-type-safety-and-errors.md`](04-golang-standards-reference/03-type-safety-and-errors.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/03-type-safety-and-errors.md --check-only` |
| `AC-CG-GO-REF-004` | Go Database Schema Mapping and Parameter Struct Conventions | [`04-golang-standards-reference/04-database-and-structs.md`](04-golang-standards-reference/04-database-and-structs.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/04-database-and-structs.md --check-only` |
| `AC-CG-GO-REF-005` | Go Semantic Naming and Package Organization | [`04-golang-standards-reference/05-naming-and-organization.md`](04-golang-standards-reference/05-naming-and-organization.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/05-naming-and-organization.md --check-only` |
| `AC-CG-GO-REF-006` | Go Enum Architecture, Constants, and Code Reuse | [`04-golang-standards-reference/06-enums-and-dry.md`](04-golang-standards-reference/06-enums-and-dry.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/06-enums-and-dry.md --check-only` |
| `AC-CG-GO-REF-007` | Go Concurrency Patterns, Channels, and Errgroup Discipline | [`04-golang-standards-reference/07-concurrency-and-patterns.md`](04-golang-standards-reference/07-concurrency-and-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference/07-concurrency-and-patterns.md --check-only` |

---

---

## 4. Detailed Acceptance Criteria Specifications

### AC-CG-GO-000: Golang Coding Standards Index & Specification Conformance

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** All Go specifications and source files adhere to affirmative boolean conventions, single-defer rules, typed enums, and structured error management with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-001: Go Positive Boolean Naming, Negation Elimination, and Guard Clauses

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** All boolean identifiers use affirmative prefixes (`is`/`has`), negations are eliminated via positive counterpart variables, mixed-polarity conditions are prohibited, and error checks remain strictly isolated with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-003: HttpMethod Enum Standard and Magic String Elimination

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** All HTTP request construction, routing, and method comparisons utilize `pkg/enums/httpmethodtype` rather than magic string literals, with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-005: Resource Defer Rules and Loop Safety Standards

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** Go functions contain at most one defer statement, zero defers execute inside loop bodies, and resources are closed deterministically with zero leaks and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-006: Go String and Slice Preallocation & Memory Efficiency Standards

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** String and slice parameters adhere to value semantics without redundant pointers, slice allocations are preallocated where capacities are known, and zero memory corruption occurs.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-007: Code Severity Taxonomy, Error Classification, and Fault Logging

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** Code review issues are classified deterministically as Code Red or Dangerous, zero unchecked errors or ungrounded paths exist, and all faults log complete contextual metadata.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-008: Unified PathUtil & FileUtil Cross-Platform Specification

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** Zero raw `os`/`io` calls exist in application packages, boolean path checks utilize positive `pathutil` functions, and all file operations return wrapped `Result[T]` types.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-009: Wrapped Result Monad and Single Return Parameter Standards

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** Functions adhere to single return parameter signatures, raw boolean returns are replaced by monadic `Result[T]` containers, and error states bind to `*appfault.AppError`.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.

---

### AC-CG-GO-ENUM-000: Go Enum Specification Index Conformance

**Given** Go source code defining domain enums.
**When** Enum implementations are audited against Go enum specifications.
**Then** Universal enum architecture, byte-based typing, required methods, and directory conventions are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-ENUM-002: Go Enum Type Definition and Constant Declarations

**Given** Go source code defining domain enums.
**When** Enum implementations are audited against Go enum specifications.
**Then** Byte-based underlying type, Invalid zero-value declaration, iota-based sequence, and PascalCase labels are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-ENUM-003: Go Enum Required Methods Implementation

**Given** Go source code defining domain enums.
**When** Enum implementations are audited against Go enum specifications.
**Then** All mandatory methods including serialization, deserialization, parsing, validation, and variant predicates are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-ENUM-004: Go Enum Architecture and Folder Layout

**Given** Go source code defining domain enums.
**When** Enum implementations are audited against Go enum specifications.
**Then** Standard directory layout, `type` package suffix, `variant.go` naming, and isolated package structure are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-ENUM-005: Go Enum Validation and Compile-Time Verification

**Given** Go source code defining domain enums.
**When** Enum implementations are audited against Go enum specifications.
**Then** Full compliance across the 50-point audit checklist, zero string comparisons in switch cases, and exhaustive verification are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-ENUM-006: Go Enum Info Object Metadata Pattern

**Given** Go source code defining domain enums.
**When** Enum implementations are audited against Go enum specifications.
**Then** Info-object pattern, `VariantInfo` struct, single lookup map, and delegating accessor methods are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-REF-000: Go Standards Reference Index Conformance

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** Go standards reference index and modular specifications are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-REF-002: Go File Sizing, Function Limits, and Complexity Rules

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** File naming matches primary types, file sizes respect the 300-line target, function lengths remain under 15 lines, and nested `if` statements are completely absent with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-REF-003: Go Type Safety, Generic Result Wrappers, and AppError Wrapping

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** Exported APIs use concrete or generic types, all application errors return `*appfault.AppError` or `Result[T]` with stack traces, and raw error returns are prohibited with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-REF-004: Go Database Schema Mapping and Parameter Struct Conventions

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** Database entities use PascalCase and `dbutil` wrappers, redundant json tags are omitted, and functions adhere to the 2-3 parameter limit with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-REF-005: Go Semantic Naming and Package Organization

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** Package declarations, file organizations, and guard functions strictly adhere to semantic naming and positive boolean rules with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-REF-006: Go Enum Architecture, Constants, and Code Reuse

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** All magic literals are replaced with typed enums/constants, API helpers return Result types, and repeated patterns are consolidated with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-GO-REF-007: Go Concurrency Patterns, Channels, and Errgroup Discipline

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** Concurrency patterns propagate contexts safely, forbidden anti-patterns are zeroed, and imports conform to the 3-group standard with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.

---

## Cross-References

- [Golang Coding Standards Overview](./readme.md)
- [Master Coding Guidelines Acceptance Criteria](../97-acceptance-criteria.md)
- [Cross-Language Standards](../01-cross-language/readme.md)
- [Cross-Language Acceptance Criteria](../01-cross-language/97-acceptance-criteria.md)

---

## Verification & Acceptance Criteria

### AC-CG-GO-REG-001: Go Acceptance Criteria Registry Conformance

**Given** The Go acceptance criteria registry in `02-spec/02-coding-guidelines/03-golang/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All Go specification criteria correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations.
