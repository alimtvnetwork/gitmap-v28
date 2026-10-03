# Python Standards — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all Python coding guideline specifications in `12-python/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-PY-[NUM]`), PEP-8 compliance, Black formatting (max 100 chars), explicit static typing, Pydantic data schemas, specific exception handling, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `12-python/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-03
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Python Criteria Inventory (`AC-CG-PY-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-PY-001` | Python Guidelines Directory Index Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| `AC-CG-PY-002` | Python Coding Standards Conformance | [`02-standards.md`](02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| `AC-CG-PY-003` | Python Dynamic Enum & Array Constants Standard | [`02-standards.md`](02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| `AC-CG-PY-004` | Python DRY Architecture & Engine Caching Conformance | [`02-standards.md`](02-standards.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |
| `AC-CG-PY-REG-001` | Python Acceptance Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-PY-001: Python Guidelines Directory Index Conformance

- [ ] Root `readme.md` provides navigation, scoring table, document inventory, and cross-references.
- [ ] All files in `12-python/` adhere to active prompt anatomy, checklist headers, and acceptance criteria blocks.
- [ ] Internal relative markdown links resolve accurately with zero broken links.

**Given** Python specifications and documentation within `12-python/`.
**When** Audited against Prompt Architect structure rules and relative path link integrity.
**Then** `12-python/readme.md` contains active `/goal` and `/learn` directives, an actionable checklist, accurate file references, and zero absolute paths.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PY-002: Python Coding Standards Conformance

- [ ] Every public function includes Python type hints for all arguments and return values.
- [ ] Raw `Any` types are eliminated; use generics, `Union`, or `Optional` when needed.
- [ ] Data validation models (`pydantic` or `@dataclass`) are enforced at system boundaries.
- [ ] Code is formatted with Black, maintaining strict PEP-8 compliance and a 100-character line length ceiling.
- [ ] Bare `except:` and `except Exception:` blocks are strictly prohibited; specific exception classes are caught.
- [ ] Low-level exceptions are wrapped in application domain-specific error structures.

**Given** Python source files and scripts across the codebase.
**When** Validated for static type hinting, data validation models, and PEP-8 compliance.
**Then** All public functions declare explicit type hints (banning raw `Any`), data schemas employ `pydantic` or `@dataclass`, line length respects the 100-character ceiling, and bare `except:` blocks are strictly absent with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PY-003: Python Dynamic Enum & Array Constants Standard

- [ ] Dynamic enums and `StrEnum` / `Enum` classes are used for state machines and categorical strings.
- [ ] Cartesian string permutations repeating host + port combinations are banned.
- [ ] Magic string literals and hardcoded array combinations are replaced by structured enum lookups and dynamic builders.
- [ ] Positive boolean names (`is_*`, `has_*`) are enforced with implicit evaluations.

**Given** Python scripts and modules defining configurable options or state machines.
**When** Inspected for hardcoded string constants and Cartesian string permutations.
**Then** Identifiers use `Enum` or `StrEnum` types, dynamic array builders eliminate repetitive string combinations, and magic literals are eliminated with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PY-004: Python DRY Architecture & Engine Caching Conformance

- [ ] Common utilities, CLI helpers, and caching logic import from `03-ai-scripts/02-shared-engine.py`.
- [ ] Duplicate utility implementations across Python scripts are eliminated and consolidated.
- [ ] Idempotent caching is implemented where applicable to guarantee sub-millisecond repeated checks.
- [ ] Unified `ExitCodeType` exit codes (`SUCCESS = 0`, `ERROR = 1`) are enforced across CLI scripts.

**Given** Python CI/CD, automation, or linting scripts in `03-ai-scripts/`.
**When** Scanned for duplicate utility implementations or unshared helpers.
**Then** Scripts import shared logic from `03-ai-scripts/02-shared-engine.py`, enforce idempotent caching where applicable, and maintain unified exit codes with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-PY-REG-001: Python Acceptance Criteria Registry Conformance

- [ ] Master registry at `02-spec/02-coding-guidelines/97-acceptance-criteria.md` references the Python registry under section `AC-12`.
- [ ] Every criterion in the Python registry maps 1:1 to an authoritative specification in `12-python/`.
- [ ] All criteria follow the structured `Given / When / Then` contract.
- [ ] Verification commands execute cleanly with exit code 0.
- [ ] Strict relative paths are used throughout with zero absolute paths or `file:///` URIs.

**Given** The Python criteria registry at `12-python/97-acceptance-criteria.md`.
**When** Evaluated for taxonomic integrity and traceability.
**Then** All listed criteria map 1:1 to specification files in `12-python/`, employ strict `Given/When/Then` structures, and specify executable verification commands with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/12-python --check-only
```
**Expected:** exit 0. Zero violations.

---

## 3. Cross-References & Traceability

- [`readme.md`](readme.md) — Python Coding Standards Overview
- [`02-standards.md`](02-standards.md) — Core Python Coding Standards
- [Coding Guidelines Master AC Registry](../97-acceptance-criteria.md) — Master acceptance criteria registry
- [Cross-Language Standards](../01-cross-language/readme.md) — Cross-language guidelines
- [Cross-Language Acceptance Criteria Registry](../01-cross-language/97-acceptance-criteria.md) — Cross-language criteria
