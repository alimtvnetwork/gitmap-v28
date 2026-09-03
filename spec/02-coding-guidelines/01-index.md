# Coding Guidelines (AI Execution Prompt)

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 02 Coding Guidelines.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `02-coding-guidelines/`.
- [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.3.0
**Status:** Active
**Updated:** 2026-08-28
**AI Confidence:** Production-Ready

## 🤖 MUST FOLLOW INSTRUCTIONS FOR ALL AI AGENTS

> **CRITICAL DIRECTIVE**: You are bound by this document. Before generating any code, writing any script, or modifying any architecture, you **MUST** internalize and apply these rules. Failure to apply these rules will result in an immediate rejection of your code.

### 1. No Generated Code or Artifacts (Never Commit)

Never commit generated code (e.g., ORM models, gRPC clients), test results, test reports, or compiled binaries. They belong in build artifacts or CI, never in source control.

### 2. Error Management is the #1 Priority

Error handling must be implemented from the **very first line of code**. Never write business logic without proper error handling wrapping it. Use the \AppError\ / \AppException\ architecture explicitly defined in the \02-spec/03-error-manage/\ folder. This is non-negotiable.

### 3. Boolean Naming (Strict Positive Assertion)

All booleans **MUST** use \is\, \has\, \can\, or \should\ prefixes and are **positively named only**.

- ❌ BAD: \!isSuccess\, \isDisabled\
- ✅ GOOD: \isFail\, \isActive\
Extract multi-part conditions into well-named boolean variables.

### 4. Nesting and Flow Control

Zero nesting. Use early returns and guard clauses. No nested \if\ blocks. If you find yourself nesting, extract the logic into a separate function immediately.

### 5. Semantic Naming (No Generics)

Absolutely NO generic garbage names. Variables named \	emp\, \data\, \obj\, \comp_100\ will trigger an instant rejection. All unit tests must be behavior-driven (e.g., \TestUpdateUser_RejectsInvalidEmail\).

### 6. Function Metrics & Signatures

- Functions: 8-15 lines. Files: < 300 lines. React components: < 100 lines.
- **Maximum 3 Parameters:** See the strict formatting rules in \03-coding-style-checklist.md\.

### 7. Never Hallucinate

If a requirement is unclear or missing, **ask a clarifying question** instead of guessing. Wrong assumptions cause rewrites.

---

## 🏗 Directory Index

- [Cross-Language Standards](./01-cross-language/01-index.md)
- [TypeScript Guidelines](./02-typescript/01-index.md)
- [Go Guidelines](./03-golang/01-index.md)
- [C# Guidelines](./07-csharp/01-index.md)
- [AI Optimization](./06-ai-optimization/01-index.md)
- **[Coding Style Checklist](./03-coding-style-checklist.md)**

---

## Verification

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-001: Coding guideline conformance: Index

**Given** Run the cross-language coding-guidelines validator against `src/` and language-specific source roots.
**When** Run the verification command shown below.
**Then** Zero CODE-RED violations are reported (functions ≤ 15 lines, files ≤ 300 lines, no nested ifs, max 2 boolean operands).

**Verification command:**

```bash
go run linter-scripts/validate-guidelines.go --path spec --max-lines 15 && python3 linter-scripts/validate-guidelines.py spec
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
