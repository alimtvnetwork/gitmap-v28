# Cross-Language Coding Guidelines — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all cross-language coding guideline specifications.
> **/learn** Reference individual criteria by canonical ID (`AC-CG-[CATEGORY]-[NUM]`) and verify compliance using targeted linters.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `01-cross-language/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

---

## 1. Boolean Principles & Guards (`AC-CG-BOOL-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-BOOL-001` | Boolean Principles Index & Directory Conformance | [`02-boolean-principles/readme.md`](02-boolean-principles/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles --check-only` |
| `AC-CG-BOOL-002` | Affirmative Naming Prefixes and Negative Word Elimination | [`02-boolean-principles/02-naming-prefixes.md`](02-boolean-principles/02-naming-prefixes.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/02-naming-prefixes.md --check-only` |
| `AC-CG-BOOL-003` | Named Guard Inversion and Complex Expression Extraction | [`02-boolean-principles/03-guards-and-extraction.md`](02-boolean-principles/03-guards-and-extraction.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/03-guards-and-extraction.md --check-only` |
| `AC-CG-BOOL-004` | Explicit Parameters, Polarity Isolation, and Condition Hygiene | [`02-boolean-principles/04-parameters-and-conditions.md`](02-boolean-principles/04-parameters-and-conditions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md --check-only` |
| `AC-CG-BOOL-005` | Quick Reference and Boolean Anti-Pattern Remediation | [`02-boolean-principles/05-quick-reference.md`](02-boolean-principles/05-quick-reference.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/05-quick-reference.md --check-only` |
| `AC-CG-BOOL-006` | Static Factory Exemptions and Structured Result Wrapper API | [`02-boolean-principles/06-exemptions-and-api.md`](02-boolean-principles/06-exemptions-and-api.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/06-exemptions-and-api.md --check-only` |
| `AC-CG-BOOL-012` | Positive Guard Functions (No Raw Negations) | [`12-no-negatives.md`](12-no-negatives.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/12-no-negatives.md --check-only` |
| `AC-CG-BOOL-024` | Boolean Flag Method Splitting | [`24-boolean-flag-methods.md`](24-boolean-flag-methods.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/24-boolean-flag-methods.md --check-only` |

---

## 2. Code Style & Spacing (`AC-CG-STYLE-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-STYLE-001` | Cross-Language Code Style Index & Directory Conformance | [`04-code-style/readme.md`](04-code-style/readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style --check-only` |
| `AC-CG-STYLE-002` | Brace Enforcement and Zero-Nesting Ban | [`04-code-style/02-braces-and-nesting.md`](04-code-style/02-braces-and-nesting.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style/02-braces-and-nesting.md --check-only` |
| `AC-CG-STYLE-003` | Multi-Part Condition Extraction and Guard Inversion | [`04-code-style/03-conditions-and-extraction.md`](04-code-style/03-conditions-and-extraction.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style/03-conditions-and-extraction.md --check-only` |
| `AC-CG-STYLE-004` | Blank Lines and Vertical Spacing Hygiene | [`04-code-style/04-blank-lines-and-spacing.md`](04-code-style/04-blank-lines-and-spacing.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style/04-blank-lines-and-spacing.md --check-only` |
| `AC-CG-STYLE-005` | Function and Type Size Caps | [`04-code-style/05-function-and-type-size.md`](04-code-style/05-function-and-type-size.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md --check-only` |
| `AC-CG-STYLE-006` | Multi-Line Formatting and Trailing Commas | [`04-code-style/06-multi-line-formatting.md`](04-code-style/06-multi-line-formatting.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style/06-multi-line-formatting.md --check-only` |
| `AC-CG-STYLE-007` | Comments and Self-Documenting Code | [`04-code-style/07-comments-and-documentation.md`](04-code-style/07-comments-and-documentation.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style/07-comments-and-documentation.md --check-only` |
| `AC-CG-STYLE-008` | Comprehensive Code Style Checklist | [`04-code-style/08-checklist.md`](04-code-style/08-checklist.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-code-style/08-checklist.md --check-only` |
| `AC-CG-STYLE-009` | Vertical Newline Spacing & Whitespace Conformance | [`21-newline-styling-examples.md`](21-newline-styling-examples.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |

---

## 3. Naming Conventions & Database Standards (`AC-CG-NAME-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-NAME-007` | PascalCase Database Naming Compliance | [`07-database-naming.md`](07-database-naming.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/07-database-naming.md --check-only` |
| `AC-CG-NAME-010` | Explicit Function Naming and Boolean Flag Elimination | [`10-function-naming.md`](10-function-naming.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/10-function-naming.md --check-only` |
| `AC-CG-NAME-011` | PascalCase Key Naming Standard | [`11-key-naming-pascalcase.md`](11-key-naming-pascalcase.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md --check-only` |
| `AC-CG-NAME-022` | Cross-Language Variable and Collection Naming | [`22-variable-naming-conventions.md`](22-variable-naming-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/22-variable-naming-conventions.md --check-only` |
| `AC-CG-NAME-028` | Deterministic Lowercase Kebab-Case Slug Compliance | [`28-slug-conventions.md`](28-slug-conventions.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/28-slug-conventions.md --check-only` |

---

## 4. Type Safety, Immutability & Parameter Architecture (`AC-CG-TYPE-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-TYPE-004` | Type Casting Elimination and Centralized Accessors | [`04-casting-elimination-patterns.md`](04-casting-elimination-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/04-casting-elimination-patterns.md --check-only` |
| `AC-CG-TYPE-013` | Strict Typing and Type Safety Enforcement | [`13-strict-typing.md`](13-strict-typing.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/13-strict-typing.md --check-only` |
| `AC-CG-TYPE-018` | In-Place Mutation Avoidance and Immutability | [`18-code-mutation-avoidance.md`](18-code-mutation-avoidance.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md --check-only` |
| `AC-CG-TYPE-032` | Branch Immutability and Clean Object Construction | [`32-branch-immutability-and-clean-construction.md`](32-branch-immutability-and-clean-construction.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/32-branch-immutability-and-clean-construction.md --check-only` |
| `AC-CG-TYPE-033` | Variadic & Spread Parameter Conformance | [`33-variadic-and-spread-parameters.md`](33-variadic-and-spread-parameters.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/33-variadic-and-spread-parameters.md --check-only` |
| `AC-CG-TYPE-034` | String Normalization & EqualFoldAny Conformance | [`34-string-normalization-and-equalfoldany.md`](34-string-normalization-and-equalfoldany.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/34-string-normalization-and-equalfoldany.md --check-only` |

---

## 5. Architecture, Patterns & Complexity (`AC-CG-ARCH-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-ARCH-006` | Cyclomatic Complexity Reduction and Guard Clauses | [`06-cyclomatic-complexity.md`](06-cyclomatic-complexity.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/06-cyclomatic-complexity.md --check-only` |
| `AC-CG-ARCH-008` | DRY Principles and Code Deduplication | [`08-dry-principles.md`](08-dry-principles.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/08-dry-principles.md --check-only` |
| `AC-CG-ARCH-016` | Lazy Evaluation and Deferred Initialization | [`16-lazy-evaluation-patterns.md`](16-lazy-evaluation-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/16-lazy-evaluation-patterns.md --check-only` |
| `AC-CG-ARCH-019` | Null Pointer Safety and Nil Dereference Prevention | [`19-null-pointer-safety.md`](19-null-pointer-safety.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/19-null-pointer-safety.md --check-only` |
| `AC-CG-ARCH-020` | Nesting Resolution and Branch Flattening | [`20-nesting-resolution-patterns.md`](20-nesting-resolution-patterns.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/20-nesting-resolution-patterns.md --check-only` |
| `AC-CG-ARCH-023` | SOLID Principles Conformance | [`23-solid-principles.md`](23-solid-principles.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/23-solid-principles.md --check-only` |

---

## 6. Changelog & Version History (`AC-CG-LOG-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-LOG-CROSS` | Cross-Language Changelog Conformance | [`changelog.md`](../../../changelog.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only` |

---

## Cross-References

- [Cross-Language Overview](./readme.md)
- [Repository Changelog](../../../changelog.md)
- [Parent Coding Guidelines Acceptance Criteria](../97-acceptance-criteria.md)
- [Canonical Architecture Specification](../../21-app/07-coding-guideline-actionable-checklist-and-acceptance/01-architecture-spec.md)

---

## Verification & Acceptance Criteria

### AC-CG-STYLE-009: Vertical Newline Spacing & Whitespace Conformance

**Given** Source code files containing function bodies, control flow blocks, and return statements across Go, TypeScript, and PHP.
**When** Code guideline linters or CI autofixers scan vertical whitespace patterns.
**Then** All functions adhere strictly to vertical newline rhythm: blank line before return (in multi-line functions), blank line after closing `}`, no blank line at function start, zero double blank lines, and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-XLANG-REG-001: Cross-Language Acceptance Criteria Registry Conformance

**Given** The cross-language acceptance criteria registry in `02-spec/02-coding-guidelines/01-cross-language/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All criteria entries across Boolean, Style, Naming, Type Safety, Architecture, Testing, and Static Analysis categories correctly map to active specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
