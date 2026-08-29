# Plan 26: Multi-Language Enums, Traits & Pattern Matching Audit

## Overview

Comprehensive audit and enforcement of multi-language Enum definitions, `*Type` suffixes, string-backed enums, helper traits, and exhaustive pattern matching across Go, TypeScript, PHP, Rust, and Python.

---

## Phase 1: Enum & Trait Violation Ledger

| Target File | Line | Identifier | Current Pattern | Language | Planned Refactoring | Status |
|---|:---:|---|---|---|---|:---:|
| `gitmap/cmd/commitin/enums.go` | 10 | `ConflictModeType` | Type Alias | Go | Added explicit `*Type` type aliases | DONE |
| `gitmap/cmd/commitin/enums.go` | 24 | `ActionKindType` | Type Alias | Go | Added explicit `*Type` type aliases | DONE |
| `gitmap/cmd/commitin/enums.go` | 38 | `ValidationVerdictType` | Type Alias | Go | Added explicit `*Type` type aliases | DONE |
| `src/constants/index.ts` | 1 | `TaskStatusType` | `as const` Object | TypeScript | Strongly-typed `*Type` union export | DONE |
| `gitmap/gitutil/dirty_inspect.go` | 45 | Character code | `rune('0' + count)` | Go | Replaced with `strconv.Itoa` | DONE |

---

## Subtasks Breakdown

- [x] [01-enum-scanner.md](.lovable/plans/subtasks/26-enums-and-traits/01-enum-scanner.md) — Create `check-enum-guidelines.py` and verify enum suffixes across languages.
- [x] [02-type-suffix-enforcement.md](.lovable/plans/subtasks/26-enums-and-traits/02-type-suffix-enforcement.md) — Enforce `*Type` naming across Go, TypeScript, PHP, Rust, and Python enums.
- [x] [03-pattern-matching-helpers.md](.lovable/plans/subtasks/26-enums-and-traits/03-pattern-matching-helpers.md) — Implement exhaustive pattern matching and helper methods (`IsValid`, `IsTerminal`, `assertNever`).
- [x] [04-verification.md](.lovable/plans/subtasks/26-enums-and-traits/04-verification.md) — Run all enum linters and quality gates.
