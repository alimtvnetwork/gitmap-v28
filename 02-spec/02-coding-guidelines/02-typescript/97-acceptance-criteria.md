# TypeScript Standards — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all TypeScript coding guideline specifications in `02-typescript/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-TS-[NUM]`), typed PascalCase enums, complete elimination of `any` and `unknown`, strict generic constraints, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `02-typescript/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. TypeScript Criteria Inventory (`AC-CG-TS-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-TS-000` | TypeScript Guidelines Index Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-002` | ConnectionStatusEnum Definition and Validation | [`02-connection-status-enum.md`](02-connection-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-003` | EntityStatusEnum Definition and Validation | [`03-entity-status-enum.md`](03-entity-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-004` | ExecutionStatusEnum Definition and Validation | [`04-execution-status-enum.md`](04-execution-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-005` | ExportStatusEnum Definition and Validation | [`05-export-status-enum.md`](05-export-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-006` | HttpMethodEnum Definition and Validation | [`06-http-method-enum.md`](06-http-method-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-007` | MessageStatusEnum Definition and Validation | [`07-message-status-enum.md`](07-message-status-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-008` | TypeScript Type Safety Remediation Plan and Any Elimination | [`08-type-safety-remediation-plan.md`](08-type-safety-remediation-plan.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-009` | TypeScript Standards Reference and Architectural Patterns | [`09-typescript-standards-reference.md`](09-typescript-standards-reference.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-010` | TypeScript Promise, Async/Await, and Concurrency Patterns | [`10-promise-await-patterns.md`](10-promise-await-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-011` | LogLevelEnum Definition and Validation | [`11-log-level-enum.md`](11-log-level-enum.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-012` | TypeScript ESLint Rules, Type Checking, and Lint Automation | [`12-eslint-enforcement.md`](12-eslint-enforcement.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-013` | TypeScript Discriminated Unions and Exhaustive Type Narrowing | [`13-discriminated-union-patterns.md`](13-discriminated-union-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-014` | TypeScript Enum Runtime Validation and Parse Guard Utilities | [`14-enum-checking-and-validation.md`](14-enum-checking-and-validation.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |
| `AC-CG-TS-015` | TypeScript State Management, Stores, and Reactive Architecture | [`15-state-management.md`](15-state-management.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-TS-000: TypeScript Guidelines Index Conformance

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All TypeScript guidelines and directory index standards are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-002: ConnectionStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All connection status handling strictly utilizes `ConnectionStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-003: EntityStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All general entity lifecycle handling strictly utilizes `EntityStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-004: ExecutionStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All task and pipeline execution lifecycle states strictly utilize `ExecutionStatus` enum constants, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-005: ExportStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All import and export operation states strictly utilize `ExportStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-006: HttpMethodEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All HTTP request method assignments and configuration tables strictly utilize `HttpMethod` enum constants, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-007: MessageStatusEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All chat message lifecycle states strictly utilize `MessageStatus` enum constants instead of magic strings, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-008: TypeScript Type Safety Remediation Plan and Any Elimination

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** All occurrences of `any`, unsafe type assertions, and untyped method signatures are systematically remediated with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-009: TypeScript Standards Reference and Architectural Patterns

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** Generics-first typing, strict enum conventions, function length bounds, and clean boolean principles are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-010: TypeScript Promise, Async/Await, and Concurrency Patterns

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** All asynchronous flows execute independent promises in parallel and eliminate redundant Promise wrappers with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-011: LogLevelEnum Definition and Validation

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All log level comparisons, interface models, and theme mappings strictly utilize `LogLevel` enum constants, with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-012: TypeScript ESLint Rules, Type Checking, and Lint Automation

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** All static analysis and ESLint enforcement mappings pass with zero warnings, zero errors, and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-013: TypeScript Discriminated Unions and Exhaustive Type Narrowing

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** Discriminated unions strictly employ named variant interfaces and PascalCase enum discriminants with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-014: TypeScript Enum Runtime Validation and Parse Guard Utilities

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** Runtime inputs are strictly verified through type guards with zero blind assertions and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-TS-015: TypeScript State Management, Stores, and Reactive Architecture

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** Global state architectures enforce domain-bounded stores and strict immutability with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

## Verification & Acceptance Criteria

### AC-CG-TS-REG-001: TypeScript Acceptance Criteria Registry Conformance

**Given** The TypeScript acceptance criteria registry file `97-acceptance-criteria.md`.
**When** Audited against master acceptance criteria specifications and CI/CD validation.
**Then** The registry file provides a complete traceable index of all criteria (`AC-CG-TS-000` through `AC-CG-TS-015`), validates structured Given / When / Then specifications, and passes guideline verification with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.

---

## Cross-References

- [TypeScript Coding Standards Overview](./readme.md)
- [Master Coding Guidelines Acceptance Criteria](../97-acceptance-criteria.md)
- [Cross-Language Standards](../01-cross-language/readme.md)
- [Cross-Language Acceptance Criteria](../01-cross-language/97-acceptance-criteria.md)
- [Golang Acceptance Criteria Registry](../03-golang/97-acceptance-criteria.md)
