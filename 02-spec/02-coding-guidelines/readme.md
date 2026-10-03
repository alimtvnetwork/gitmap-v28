# Coding Guidelines (AI Execution Prompt)

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 02 Coding Guidelines.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Internalize and enforce `01-specification-and-coding-guideline-standard.md` as the overarching 4-part execution meta-standard across all guidelines.
- [ ] `/learn` Adhere strictly to affirmative boolean conventions (`is`, `has`), ban explicit `true` comparisons (`if isReady`), and eliminate mixed-polarity conditionals (`if isA && !isB`).
- [ ] `/goal` Comply with strict code style constraints: max 3 parameters per function signature, mandatory vertical spacing before/after control flow, and guard clauses for early returns.
- [ ] `/learn` Run automated verification linters via `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only` ensuring zero CODE-RED violations before committing.

. **CRITICAL AI INSTRUCTION:** This readme.md file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.4.0
**Status:** Active
**Updated:** 2026-10-03
**AI Confidence:** Production-Ready

---

## 📐 Primary Meta-Specification: Active AI Execution Standard

The foundational blueprint governing this entire directory is documented in [`01-specification-and-coding-guideline-standard.md`](./01-specification-and-coding-guideline-standard.md).

Every coding guideline within this module is an **active AI execution prompt** rather than passive documentation. All guideline files strictly adhere to a standardized **4-part anatomy**:

1. **Prompt Header (`/goal` & `/learn`):** Unambiguous operational targets and learning directives informing AI agents of the scope and constraints prior to code generation.
2. **Actionable CI/CD & Agent Checklist (`- [ ] /goal` & `- [ ] /learn`):** Machine-verifiable step-by-step checklist verifying non-negotiable compliance, linting rules, and repository hygiene.
3. **Core Specification & Code Contracts:** Explicit architectural rules, bounded sizing limits, structured error returns, and concrete `BAD`/`GOOD` code examples across supported languages.
4. **Acceptance Criteria & Verification Command (`AC-CG-XXX`):** Deterministic, testable acceptance criteria paired with a runnable CLI verification command to enforce continuous zero-defect integration.

---

## 🤖 MUST FOLLOW INSTRUCTIONS FOR ALL AI AGENTS

> **CRITICAL DIRECTIVE**: You are bound by this document. Before generating any code, writing any script, or modifying any architecture, you **MUST** internalize and apply these rules. Failure to apply these rules will result in an immediate rejection of your code.

### 1. No Generated Code or Artifacts (Never Commit)

Never commit generated code (e.g., ORM models, gRPC clients), test results, test reports, or compiled binaries. They belong in build artifacts or CI, never in source control.

### 2. Error Management is the #1 Priority

Error handling must be implemented from the **very first line of code**. Never write business logic without proper error handling wrapping it. Use the `AppError` / `AppException` architecture explicitly defined in the `02-spec/03-error-manage/` folder. This is non-negotiable.

### 3. Boolean Logic & Naming (Strict Affirmative Only)

All booleans **MUST** use `is` or `has` prefixes and are **positively named only** (`can`, `should`, `was` are banned).

- ❌ FORBIDDEN: `!isSuccess`, `isDisabled`
- ✅ REQUIRED: `isFail`, `isActive`
- **No Explicit True Checks (TOTAL BAN):** NEVER evaluate a boolean explicitly against `true` (e.g., `if isReady == true`). Positive booleans MUST ALWAYS be evaluated implicitly: `if isReady { ... }`.
- **No Mixed Polarity (TOTAL BAN):** NEVER combine a positive check and a negative check in the same `if` condition (e.g., `if isA && !isB`). Split into discrete conditions or extract into an affirmative variable.

### 4. Nesting and Flow Control

Zero nesting. Use early returns and guard clauses. No nested `if` blocks. If you find yourself nesting, extract the logic into a separate function immediately.

### 5. Semantic Naming (No Generics)

Absolutely NO generic garbage names. Variables named `temp`, `data`, `obj`, `comp_100` will trigger an instant rejection. All unit tests must be behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`).

### 6. Function Metrics & Signatures

- Functions: 8-15 lines. Files: < 300 lines. React components: < 100 lines.
- **Maximum 3 Parameters:** See the strict formatting rules in `03-coding-style-checklist.md`. Signatures requiring 4+ parameters MUST be refactored into a parameter options struct.

### 7. Never Hallucinate

If a requirement is unclear or missing, **ask a clarifying question** instead of guessing. Wrong assumptions cause rewrites.

---

## 🏗 Directory Index

### Root Specifications

| Specification | Purpose / Scope |
|---|---|
| [`01-specification-and-coding-guideline-standard.md`](./01-specification-and-coding-guideline-standard.md) | Universal 4-part AI execution prompt standard and guideline anatomy |
| [`02-canonical-size-tier.md`](./02-canonical-size-tier.md) | Single source of truth for file, function, and component line sizing tiers |
| [`03-coding-style-checklist.md`](./03-coding-style-checklist.md) | Root coding style rules: 3-parameter cap, PascalCase acronyms, vertical spacing |
| [`04-consolidated-review-guide-condensed.md`](./04-consolidated-review-guide-condensed.md) | Condensed quick-reference review checklist for rapid agent audits |
| [`05-consolidated-review-guide.md`](./05-consolidated-review-guide.md) | Comprehensive cross-language code review standard and guidelines |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | Master acceptance criteria registry and verification index for all guidelines |

### Cross-Language & Core Conventions

| Module | Description |
|---|---|
| [`01-cross-language`](./01-cross-language/readme.md) | Cross-language code style, booleans, semantic naming, and strong typing |
| [`06-ai-optimization`](./06-ai-optimization/readme.md) | Anti-hallucination, citation rules, grounded checklists, and token efficiency |
| [`06-cicd-integration`](./06-cicd-integration/readme.md) | CI/CD integration, SARIF contracts, GitHub Actions linter pack, and runners |
| [`08-file-folder-naming`](./08-file-folder-naming/readme.md) | Strict lowercase filenames, hyphen separators, and repository hierarchy |
| [`11-security`](./11-security/readme.md) | Security baselines, JWT verification, cryptographic standards, secret isolation |

### Polyglot Language Guidelines

| Language / Domain | Description |
|---|---|
| [`02-typescript`](./02-typescript/readme.md) | TypeScript standards: strict types, interfaces, async/await, and error envelopes |
| [`03-golang`](./03-golang/readme.md) | Go standards: structured `*appfault.AppError`, `Result[T]`, value semantics |
| [`04-php`](./04-php/readme.md) | PHP standards: strict types, PSR-12, typed properties, and response envelopes |
| [`05-rust`](./05-rust/readme.md) | Rust standards: idiomatic error handling, memory safety, borrowing, traits |
| [`07-csharp`](./07-csharp/readme.md) | C# standards: modern language features, async/await, LINQ, nullable types |
| [`12-python`](./12-python/readme.md) | Python standards: type hinting, PEP 8, enum standards, and script architecture |
| [`13-cpp`](./13-cpp/readme.md) | C++ standards: modern C++20, RAII, memory safety, zero raw pointer bloat |
| [`09-powershell-integration`](./09-powershell-integration/readme.md) | PowerShell conventions, cross-platform scripting, and automation pipelines |
| [`10-research`](./10-research/readme.md) | Research methodologies, technology evaluations, and exploratory notes |

### Application Specifications

| Module | Description |
|---|---|
| [`21-app`](./21-app/readme.md) | Application architecture, core features, component designs, and workflows |
| [`22-app-issues`](./22-app-issues/readme.md) | App issue tracking, 4-part root cause analysis (RCA), and fix ledger |
| [`23-app-db`](./23-app-db/readme.md) | Application database models, SQLite split-DB architecture, and migrations |
| [`24-app-ui-design-system`](./24-app-ui-design-system/readme.md) | UI design system, color tokens, typography, and responsive component patterns |

---

## Verification & Acceptance Criteria

_Auto-generated section — see [`02-spec/02-coding-guidelines/97-acceptance-criteria.md`](./97-acceptance-criteria.md) for the full criteria index._

### AC-CG-001: Coding guideline conformance: Index

**Given** Run the cross-language coding-guidelines validator against `src/` and language-specific source roots.
**When** Run the verification command shown below.
**Then** Zero CODE-RED violations are reported (functions ≤ 15 lines, files ≤ 300 lines, no nested ifs, max 2 boolean operands).

**Verification command:**

```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-10-03_
