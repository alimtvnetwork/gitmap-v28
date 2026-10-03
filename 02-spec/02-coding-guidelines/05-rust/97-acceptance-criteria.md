# Rust Standards — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all Rust coding guideline specifications in `05-rust/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-RUST-[NUM]`), RFC 430 naming conventions with PascalCase boundary exceptions, structured `thiserror`/`anyhow` errors, Tokio async cancellation safety, zero unannotated `unsafe`, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `05-rust/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Rust Criteria Inventory (`AC-CG-RUST-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-RUST-001` | Rust Coding Guidelines Index & Architecture Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| `AC-CG-RUST-002` | Rust Naming Conventions and Positive Booleans | [`02-naming-conventions.md`](02-naming-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| `AC-CG-RUST-003` | Rust Result Handling, Error Trait and No Panics | [`03-error-handling.md`](03-error-handling.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| `AC-CG-RUST-004` | Rust Async Tokio Cancellation and Send Bounds | [`04-async-patterns.md`](04-async-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| `AC-CG-RUST-005` | Rust Zero Unsafe and RAII Lifetime Safety | [`05-memory-safety.md`](05-memory-safety.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| `AC-CG-RUST-006` | Rust Test Suite Organization and Mocking Safety | [`06-testing-standards.md`](06-testing-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| `AC-CG-RUST-007` | Rust FFI Safety Boundaries and Platform Isolation | [`07-ffi-platform.md`](07-ffi-platform.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |
| `AC-CG-RUST-REG-001` | Rust Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-RUST-001: Rust Coding Guidelines Index & Architecture Conformance

**Given** Rust source code, crate manifests, and specifications across modules and subsystems.
**When** Guideline linters audit the codebase against Rust coding standards and architecture rules.
**Then** All Rust specifications and source files adhere to RFC 430 conventions, PascalCase boundaries, structured error handling, async patterns, and memory safety rules with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-RUST-002: Rust Naming Conventions and Positive Booleans

**Given** Rust source code across crate modules, database structs, serde models, and functions.
**When** Guideline linters audit the codebase for naming convention and boolean compliance.
**Then** All functions, methods, variables, modules, and crates adhere to RFC 430 snake_case, database identifiers and serialized enum values use PascalCase, and booleans strictly employ affirmative `is_*`/`has_*` naming with implicit evaluations and zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-RUST-003: Rust Result Handling, Error Trait and No Panics

**Given** Rust source code across library crates, domain modules, and application binaries.
**When** Audited against this specification using guideline linters and static checks.
**Then** All domain errors implement `std::error::Error` (via `thiserror`), application orchestration uses `anyhow::Result` with context, production code avoids `.unwrap()` and `.expect()` panics, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-RUST-004: Rust Async Tokio Cancellation and Send Bounds

**Given** Rust asynchronous tasks, channel endpoints, and Tokio event loops.
**When** Audited against async runtime discipline and cancellation safety standards.
**Then** All spawned futures satisfy `Send + Sync + 'static`, channels use bounded buffers, cancel-safe idioms protect I/O from lost writes, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-RUST-005: Rust Zero Unsafe and RAII Lifetime Safety

**Given** Rust codebases handling memory allocations, lifetimes, references, and unsafe blocks.
**When** Audited against memory safety standards and lifetime hygiene rules.
**Then** Borrowing is preferred over cloning, every `unsafe` block includes a mandatory `// SAFETY:` comment justifying pointer validity and sound invariants, all unsafe calls are encapsulated in safe RAII wrappers, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-RUST-006: Rust Test Suite Organization and Mocking Safety

**Given** Rust unit and integration test suites across crates and services.
**When** Audited against test architecture, naming standards, and isolation guidelines.
**Then** All test functions follow standard naming, async tests use isolated Tokio runtimes, OS dependencies are mocked via traits, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-RUST-007: Rust FFI Safety Boundaries and Platform Isolation

**Given** Rust codebases integrating with external C libraries, Win32 APIs, or Unix system calls.
**When** Audited against FFI safety guidelines and platform isolation architecture.
**Then** All external C interfaces are isolated behind safe `PlatformApi` trait boundaries, every OS handle implements deterministic RAII drop cleanup, conditional compilation is confined to module boundaries, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Rust Standards Overview | [`readme.md`](readme.md) |
| Rust Naming Conventions | [`02-naming-conventions.md`](02-naming-conventions.md) |
| Rust Error Handling | [`03-error-handling.md`](03-error-handling.md) |
| Rust Async Patterns | [`04-async-patterns.md`](04-async-patterns.md) |
| Rust Memory Safety | [`05-memory-safety.md`](05-memory-safety.md) |
| Rust Testing Standards | [`06-testing-standards.md`](06-testing-standards.md) |
| Rust FFI & Platform | [`07-ffi-platform.md`](07-ffi-platform.md) |
| Coding Guidelines Root AC Registry | [`../97-acceptance-criteria.md`](../97-acceptance-criteria.md) |
| Cross-Language Standards | [`../01-cross-language/readme.md`](../01-cross-language/readme.md) |

---

## Verification & Acceptance Criteria

### AC-CG-RUST-REG-001: Rust Acceptance Criteria Registry Conformance

**Given** The Rust acceptance criteria registry in `02-spec/02-coding-guidelines/05-rust/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All Rust specification criteria correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.
