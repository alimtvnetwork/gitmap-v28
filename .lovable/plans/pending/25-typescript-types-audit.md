# Plan 25: TypeScript Strict Typing, Discriminated Unions & Architecture Audit

## Overview

Comprehensive audit of TypeScript type safety across the repository, verifying zero `any` usage, strict discriminated unions, `Result<T>` return envelopes, and compile-time type validation.

---

## Phase 1: TypeScript Type Safety Ledger

| Target File | Line | Symbol / Function | Current Type / Pattern | Violation | Target Refactoring | Status |
|---|:---:|---|---|---|---|:---:|
| `src/types/result.ts` | 1 | `Result<T>` | Discriminated Union | Missing Result envelope | Implemented `Result<T>` with `AppError` | DONE |
| `src/types/helpJson.ts` | 29 | `isHelpJsonPayload` | Type Guard | Untyped JSON boundary | Strongly-typed narrowing guard | DONE |
| `src/data/commands.ts` | 1 | `CommandItem` | Typed Interface | Strict options types | Verified 100% strict typed interfaces | DONE |

---

## Subtasks Breakdown

- [x] [01-result-envelope-and-guards.md](.lovable/plans/subtasks/25-typescript-types/01-result-envelope-and-guards.md) — Implement strongly-typed `Result<T>` envelope and exhaustive pattern matching helpers.
- [x] [02-type-safety-audit.md](.lovable/plans/subtasks/25-typescript-types/02-type-safety-audit.md) — Audit codebase for `any` types and verify zero unsafe type assertions.
- [x] [03-verification.md](.lovable/plans/subtasks/25-typescript-types/03-verification.md) — Verify `tsc --noEmit`, vitest test suites, and all quality gates.
