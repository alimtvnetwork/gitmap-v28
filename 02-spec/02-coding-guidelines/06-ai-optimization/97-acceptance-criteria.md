# AI Optimization — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all AI optimization specifications in `06-ai-optimization/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-AI-[NUM]`), anti-hallucination rules, pre-output validation checklists, common mistake remediations, enum standards, citation mandates, and memory lifecycle controls.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `06-ai-optimization/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths across all guideline links.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. AI Optimization Criteria Inventory (`AC-CG-AI-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-AI-001` | AI Optimization Index & Architecture Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-002` | Anti-Hallucination Rules Enforcement | [`02-anti-hallucination-rules.md`](02-anti-hallucination-rules.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-003` | AI Quick-Reference Checklist Validation | [`03-ai-quick-reference-checklist.md`](03-ai-quick-reference-checklist.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-004` | Common AI Mistakes Prevention & Remediation | [`04-common-ai-mistakes.md`](04-common-ai-mistakes.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-005` | Condensed Master Guidelines Conformance | [`05-condensed-master-guidelines.md`](05-condensed-master-guidelines.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-006` | Mandatory Citations & Relative Path Governance | [`06-citation-requirement.md`](06-citation-requirement.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-007` | Cross-Language Enum Standards Enforcement | [`07-enum-naming-quick-reference.md`](07-enum-naming-quick-reference.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-008` | Hallucination Prevention & Static Verification | [`08-hallucination-checks.md`](08-hallucination-checks.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-009` | Agent Memory Lifecycle & TTL Maintenance | [`09-agent-memory-lifecycle.md`](09-agent-memory-lifecycle.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |
| `AC-CG-AI-REG-001` | AI Optimization Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only` |

---

## 2. Core Requirements & Validation Rules

### Required

- [ ] Anti-hallucination rules cover all 5 language categories (cross-language, Go, TS, PHP, Rust)
- [ ] Each rule has a unique ID (e.g., `AH-N1`)
- [ ] Quick-reference checklist is machine-parsable with checkboxes
- [ ] Common mistakes include before/after code examples
- [ ] All rules cross-reference the canonical spec file they enforce

### Validation

- [ ] Zero overlap between anti-hallucination rules and quick-reference checklist (no duplicate content)
- [ ] All referenced spec files exist and links resolve accurately with zero broken links

---

## 3. Detailed Acceptance Criteria Specifications

### AC-CG-AI-001: AI Optimization Index & Architecture Conformance

- [ ] Root `readme.md` provides complete module navigation, execution prompts, and cross-references.
- [ ] All specification files in `06-ai-optimization/` adhere to active prompt anatomy, checklist headers, and acceptance criteria blocks.
- [ ] Internal relative markdown links resolve accurately with zero broken links.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-002: Anti-Hallucination Rules Enforcement

- [ ] All 40 anti-hallucination rules are actively enforced across naming, type safety, boolean logic, structure, error handling, enums, C#, and caching.
- [ ] Zero generated code files (`*.generated.*`, ORM models) or test artifacts are staged into Git.
- [ ] Abbreviations use first-letter-only casing (`Id`, `Url`, `Api`, `Json`).

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-003: AI Quick-Reference Checklist Validation

- [ ] Pre-output validation checklist covers naming, structure, error management, and caching boundaries.
- [ ] Machine-parsable checkboxes allow rapid verification of code correctness before emission.
- [ ] Zero nested `if` statements and maximum 15-line function bodies are enforced.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-004: Common AI Mistakes Prevention & Remediation

- [ ] Top 15 ranked AI mistake catalog with explicit before-and-after corrections is documented and enforced.
- [ ] Prohibits camelCase JSON keys, uppercase abbreviations, multi-return Go signatures, and raw `fmt.Errorf`.
- [ ] Mandates query cache stale times and immediate cache invalidation upon entity mutations.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-005: Condensed Master Guidelines Conformance

- [ ] High-density distillation of all master coding guidelines is maintained for LLM context windows.
- [ ] Rule Zero strictly prohibits checking in generated code, binaries, and build artifacts.
- [ ] Positive boolean logic, enum standards, and universal error envelopes are reinforced.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-006: Mandatory Citations & Relative Path Governance

- [ ] All code generation, planning tasks, subtasks, and memory logs cite authoritative `02-spec/` or `.ai-memory/` paths.
- [ ] Total ban on absolute filesystem paths (`/absolute/...`, `C:\...`) and `file:///` URIs across all markdown and code files.
- [ ] Portable relative paths from the git repository root are strictly enforced.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-007: Cross-Language Enum Standards Enforcement

- [ ] Universal enum rules eliminate magic strings and enforce PascalCase values with co-located string constants.
- [ ] Go: `Variant byte`, `Invalid` zero value, package-scoped constants, and full method suite (`String`, `Label`, `IsValid`, `MarshalJSON`).
- [ ] TypeScript: String enums with UPPER_SNAKE values under `src/lib/enums/`.
- [ ] PHP: String-backed enums with `Type` suffix and `isEqual()` method under `includes/Enums/`.
- [ ] Rust: `#[serde(rename_all = "PascalCase")]` and `#[sqlx(rename_all = "PascalCase")]` on all serialized and DB-mapped enums.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-008: Hallucination Prevention & Static Verification

- [ ] Strict "Read Before Write" protocol mandates inspecting router definitions, database schemas, and package manifests before emitting code.
- [ ] Strict null and type enforcement (`strict: true`, `Result[T]`, `mypy`) catches invalid property access early.
- [ ] Immediate static analysis and linter validation runs after code generation to ensure correctness.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-AI-009: Agent Memory Lifecycle & TTL Maintenance

- [ ] Tasks in `.ai-memory/plans/pending/` and `planned/` are transitioned to `done/` upon task completion.
- [ ] Deduplication scans detect and resolve redundant memory entries when rules are formalized into `02-spec/`.
- [ ] Conflicting memory entries are purged immediately, upholding canonical `02-spec/` supremacy.
- [ ] Explicit `Updated:` ISO date stamps are required on all memory records to ensure deterministic conflict resolution.

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

---

## Verification & Acceptance Criteria

### AC-CG-AI-REG-001: AI Optimization Criteria Registry Conformance

**Given** The AI optimization criteria registry in `02-spec/02-coding-guidelines/06-ai-optimization/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All AI optimization specification criteria correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.
