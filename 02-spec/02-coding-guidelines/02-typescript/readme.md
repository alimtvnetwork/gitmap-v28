# TypeScript Standards (AI Execution Prompt)

> **/goal** Master and enforce TypeScript coding standards, typed enum conventions, type safety enforcement, and zero magic strings across frontend and shared codebases.
> **/learn** Internalize PascalCase string enum rules, prohibition of string union types for domain states, complete eradication of `any` and `unknown`, and strict type safety patterns.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce typed `enum` constructs with PascalCase members for domain and status models; ban raw string union types.
- [ ] `/learn` Eliminate all usages of `any`, `unknown`, and `Record<string, unknown>` in favor of strongly-typed domain interfaces.
- [ ] `/goal` Replace all magic string literals in status checks and HTTP calls with canonical enum constants.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Status:** Active
**Updated:** 2026-10-02
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Overview

TypeScript-specific coding standards, enum definitions, and type safety enforcement rules. All enums must use proper `enum` syntax with PascalCase values and a `Type` suffix — string union types are prohibited. Generics must use concrete type parameters; `unknown`, `any`, and `Record<string, unknown>` are banned.

---

## Keywords

`typescript` · `enums` · `type-safety` · `pascalcase` · `connection-status` · `entity-status` · `execution-status` · `export-status` · `http-method` · `message-status` · `remediation-plan` · `standards-reference`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `readme.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

| # | File | Category | Description |
|---|------|----------|-------------|
| 01 | [01-connection-status-enum.md](./02-connection-status-enum.md) | Enum | Connection status enum definition |
| 02 | [02-entity-status-enum.md](./03-entity-status-enum.md) | Enum | Entity status enum definition |
| 03 | [03-execution-status-enum.md](./04-execution-status-enum.md) | Enum | Execution status enum definition |
| 04 | [04-export-status-enum.md](./05-export-status-enum.md) | Enum | Export status enum definition |
| 05 | [05-http-method-enum.md](./06-http-method-enum.md) | Enum | HTTP method enum definition |
| 06 | [06-message-status-enum.md](./07-message-status-enum.md) | Enum | Message status enum definition |
| 07 | [07-type-safety-remediation-29-plan.md](./08-type-safety-remediation-plan.md) | Plan | Type safety remediation plan (v2.0.0) — eliminates `any`, `unknown`, string unions |
| 08 | [08-typescript-standards-reference.md](./09-typescript-standards-reference.md) | Reference | Comprehensive TypeScript standards reference |
| 09 | [09-promise-await-patterns.md](./10-promise-await-patterns.md) | Patterns | Promise/await patterns and async conventions |
| 10 | [10-log-level-enum.md](./11-log-level-enum.md) | Enum | Log level enum definition (Debug, Info, Warn, Error, Fatal) |
| 11 | [11-eslint-enforcement.md](./12-eslint-enforcement.md) | Enforcement | ESLint rule mapping + SonarQube integration |
| 12 | [12-discriminated-union-patterns.md](./13-discriminated-union-patterns.md) | Patterns | Discriminated union & action type patterns — no inline types, PascalCase enums |
| 97 | [97-acceptance-criteria.md](./97-acceptance-criteria.md) | Testing | Acceptance criteria |

---

## Document Inventory

| File |
|------|
| 99-consistency-report.md |

## Cross-References

| Reference | Location |
|-----------|----------|
| Parent Overview | `../readme.md` |
| Cross-Language Rules | `../01-cross-language/readme.md` |
| Coding Guidelines Memory | `../../../.ai-memory/memories/constraints/coding-guidelines.md` |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-000: TypeScript Guidelines Index Conformance

**Given** TypeScript source code defining domain types, enums, and API models.
**When** Codebases are audited against TypeScript enum standards.
**Then** All TypeScript guidelines and directory index standards are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
