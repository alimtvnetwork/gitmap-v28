# Modern C++ Standards — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all Modern C++ coding guideline specifications in `13-cpp/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-CPP-[NUM]`), modern C++20 standard baseline, concepts over enable_if, smart pointer memory safety (std::unique_ptr/std::shared_ptr), Rule of Zero/Five, PascalCase types, snake_case functions, and FFI exception containment.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `13-cpp/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-03
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Modern C++ Criteria Inventory (`AC-CG-CPP-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-CPP-001` | Modern C++ Guidelines Directory Index Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| `AC-CG-CPP-002` | Modern C++ Standards Conformance | [`02-standards.md`](02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| `AC-CG-CPP-003` | C++ Memory Safety & RAII Resource Management | [`02-standards.md`](02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| `AC-CG-CPP-004` | C++ FFI Boundary Exception Safety & Standard Types | [`02-standards.md`](02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |
| `AC-CG-CPP-REG-001` | Modern C++ Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-CPP-001: Modern C++ Guidelines Directory Index Conformance

- [ ] Root `readme.md` provides navigation, scoring table, document inventory, and cross-references.
- [ ] All files in `13-cpp/` adhere to active prompt anatomy, checklist headers, and acceptance criteria blocks.
- [ ] Internal relative markdown links resolve accurately with zero broken links.

**Given** C++ specifications and documentation within `13-cpp/`.
**When** Audited against Prompt Architect structure rules and link validity.
**Then** `13-cpp/readme.md` features active execution prompts, an actionable agent checklist, accurate relative links, and zero absolute paths.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CPP-002: Modern C++ Standards Conformance

- [ ] All C++ code is written against C++20 standard baseline or newer.
- [ ] Template constraints use C++20 concepts instead of raw SFINAE or `enable_if`.
- [ ] Naming conventions strictly enforce `PascalCase` for Structs, Classes, and Enums.
- [ ] Naming conventions strictly enforce `snake_case` for Functions and Variables.
- [ ] Macros use `SCREAMING_SNAKE_CASE` and are avoided where `constexpr` can be used.

**Given** C++ header and implementation files across repositories.
**When** Audited for language dialect baseline and template constraints.
**Then** Code conforms to the C++20 standard or newer, concepts are utilized instead of raw SFINAE / `enable_if`, and naming rules (`PascalCase` for types/enums, `snake_case` for functions/variables) are strictly followed with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CPP-003: C++ Memory Safety & RAII Resource Management

- [ ] Manual `new` and `delete` invocations are strictly prohibited across all modern C++ code.
- [ ] Memory ownership is explicitly managed via `std::unique_ptr` (exclusive) or `std::shared_ptr` (shared).
- [ ] Resource lifetimes are bound deterministically to object lifecycles under RAII principles.
- [ ] Classes adhere to the Rule of Zero; if custom lifecycle operations are required, all five special member functions (destructor, copy ctor, move ctor, copy assign, move assign) are explicitly defined or deleted (Rule of Five).

**Given** C++ classes, resource handles, and heap-allocated objects.
**When** Inspected for memory management and ownership patterns.
**Then** Manual `new` and `delete` invocations are strictly absent, ownership is expressed via `std::unique_ptr` or `std::shared_ptr`, and classes adhere to the Rule of Zero (or complete Rule of Five) with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CPP-004: C++ FFI Boundary Exception Safety & Standard Types

- [ ] Exceptions are never used for routine control flow.
- [ ] Exceptions are strictly contained within C++ code and prohibited from escaping into C/FFI layers.
- [ ] All boundary functions interfacing with foreign runtimes are marked `noexcept` and wrapped with explicit `try / catch (...)` handlers.
- [ ] Cross-boundary results are communicated via status codes, enums, or C-compatible outcome structs.

**Given** C++ modules exposing C-compatible FFI or receiving foreign function calls.
**When** Audited for exception leakage and ABI stability.
**Then** Zero exceptions escape C++ boundaries into foreign runtime contexts (enforced via `noexcept` and `try/catch` wrapping), and error codes/enums communicate status across FFI layers with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-CPP-REG-001: Modern C++ Acceptance Criteria Registry Conformance

- [ ] Master registry at `02-spec/02-coding-guidelines/97-acceptance-criteria.md` references the Modern C++ registry under section `AC-13`.
- [ ] Every criterion in the Modern C++ registry maps 1:1 to an authoritative specification in `13-cpp/`.
- [ ] All criteria follow the structured `Given / When / Then` contract.
- [ ] Verification commands execute cleanly with exit code 0.
- [ ] Strict relative paths are used throughout with zero absolute paths or `file:///` URIs.

**Given** The C++ criteria registry at `13-cpp/97-acceptance-criteria.md`.
**When** Evaluated for taxonomic integrity and traceability.
**Then** Every entry links 1:1 to authoritative specifications in `13-cpp/`, features explicit `Given/When/Then` definitions, and defines functional verification commands with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/13-cpp --check-only
```
**Expected:** exit 0. Zero violations.

---

## 3. Cross-References & Traceability

- [`readme.md`](readme.md) — Modern C++ Standards Overview
- [`02-standards.md`](02-standards.md) — Modern C++ Coding Standards
- [Coding Guidelines Master AC Registry](../97-acceptance-criteria.md) — Master acceptance criteria registry
- [Cross-Language Standards](../01-cross-language/readme.md) — Cross-language guidelines
- [Cross-Language Acceptance Criteria Registry](../01-cross-language/97-acceptance-criteria.md) — Cross-language criteria
