# TypeScript Standards

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
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
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---


| # | File | Category | Description |
|---|------|----------|-------------|
| 01 | [01-connection-status-enum.md](./01-connection-status-enum.md) | Enum | Connection status enum definition |
| 02 | [02-entity-status-enum.md](./02-entity-status-enum.md) | Enum | Entity status enum definition |
| 03 | [03-execution-status-enum.md](./03-execution-status-enum.md) | Enum | Execution status enum definition |
| 04 | [04-export-status-enum.md](./04-export-status-enum.md) | Enum | Export status enum definition |
| 05 | [05-http-method-enum.md](./05-http-method-enum.md) | Enum | HTTP method enum definition |
| 06 | [06-message-status-enum.md](./06-message-status-enum.md) | Enum | Message status enum definition |
| 07 | [07-type-safety-remediation-plan.md](./07-type-safety-remediation-plan.md) | Plan | Type safety remediation plan (v2.0.0) — eliminates `any`, `unknown`, string unions |
| 08 | [08-typescript-standards-reference.md](./08-typescript-standards-reference.md) | Reference | Comprehensive TypeScript standards reference |
| 09 | [09-promise-await-patterns.md](./09-promise-await-patterns.md) | Patterns | Promise/await patterns and async conventions |
| 10 | [10-log-level-enum.md](./10-log-level-enum.md) | Enum | Log level enum definition (Debug, Info, Warn, Error, Fatal) |
| 11 | [11-eslint-enforcement.md](./11-eslint-enforcement.md) | Enforcement | ESLint rule mapping + SonarQube integration |
| 12 | [12-discriminated-union-patterns.md](./12-discriminated-union-patterns.md) | Patterns | Discriminated union & action type patterns — no inline types, PascalCase enums |
| 97 | [97-acceptance-criteria.md](./97-acceptance-criteria.md) | Testing | Acceptance criteria |
| 98 | [98-changelog.md](./98-changelog.md) | Meta | Changelog |

---

## Document Inventory

| File |
|------|
| 99-consistency-report.md |


## Cross-References

| Reference | Location |
|-----------|----------|
| Parent Overview | `../00-overview.md` |
| Cross-Language Rules | `../01-cross-language/00-overview.md` |
| Coding Guidelines Memory | `../../../.lovable/memories/constraints/coding-guidelines.md` |
