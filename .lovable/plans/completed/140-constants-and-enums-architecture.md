# Plan 140: Constants & Enums Architecture — Master Architectural Specification

Trigger Keywords & Aliases: `cg-enums`, `cg-constants`, `cg-execute enums`, `audit constants`, `fix enums`, `eliminate magic strings`, `eliminate magic numbers`, `enforce enum suffix`, `constants and enums audit`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)

---

## 1. Executive Summary & Objective

Autonomously scan, plan, refactor, and fix all constants, enums, magic string literals, raw character/rune literals (`rune(10)`), and magic number violations across the codebase, modifying source files directly to enforce the `*Type` enum suffix, extract magic numbers/strings into dedicated constant files, and use typed enums and traits until 100% green without stopping.

Automated scan across Go, TypeScript, and Python identified:
- **Zero raw rune numeric casts** (`rune(10)`).
- **Zero TypeScript status/state string unions** (all properly using `*Type` enums).
- **11 Go enum types missing mandatory `*Type` suffix** across `cli/cmd/commitin/`, `cli/dbengine_old/`, and `scripts/changelog/`.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md`: Universal enum and constant conventions and mandatory `*Type` suffix.
- `spec/02-coding-guidelines/06-ai-optimization/05-enum-naming-quick-reference.md`: Strict naming rules for enum types and constants.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100 coding lines).
- `spec/02-coding-guidelines/06-ai-optimization/01-index.md`: Zero truncation, zero placeholders, zero ghost diffs.
- `.lovable/coding-guidelines.md`: Master consolidated coding guidelines.

---

## 3. Domain-Specific Rules for Constants & Enums

1. **Rule 1 (Mandatory `*Type` Suffix):** Every Enum type definition MUST end with `Type` (e.g. `ConflictModeType`, `InputKindType`, `SourceKindType`, `ModeType`). Backward-compatible type aliases (`type ConflictMode = ConflictModeType`) must be provided where needed to prevent API drift.
2. **Rule 2 (Zero Raw Rune Numeric Conversions):** Never use raw integer character casts like `rune(10)` or `string(rune(10))`. Use centralized character/string constants like `NewLineUnix`.
3. **Rule 3 (Centralized Dedicated Packages):** All constants and enums must reside in dedicated directories (`constants/`, `enums/`, `src/types/`, `src/enums/`), never scattered inline.
4. **Rule 4 (No String Unions for Enums):** TypeScript status and state collections must use typed enums or `as const` object enums ending with `*Type`.

---

## 4. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Violation Description | Planned Fix | Status |
|:---|:---|:---:|:---|:---|:---|:---:|
| V-01 | `cli/cmd/commitin/enums.go` | 15 | `type ConflictMode uint8` | Enum missing mandatory `*Type` suffix | Rename to `ConflictModeType`, add `type ConflictMode = ConflictModeType` | Fixed |
| V-02 | `cli/cmd/commitin/enums.go` | 38 | `type InputKind uint8` | Enum missing mandatory `*Type` suffix | Rename to `InputKindType`, add `type InputKind = InputKindType` | Fixed |
| V-03 | `cli/cmd/commitin/enums.go` | 64 | `type RunStatus uint8` | Enum missing mandatory `*Type` suffix | Rename to `RunStatusType`, add `type RunStatus = RunStatusType` | Fixed |
| V-04 | `cli/cmd/commitin/enums.go` | 96 | `type CommitOutcome uint8` | Enum missing mandatory `*Type` suffix | Rename to `CommitOutcomeType`, add `type CommitOutcome = CommitOutcomeType` | Fixed |
| V-05 | `cli/cmd/commitin/enums.go` | 122 | `type SkipReason uint8` | Enum missing mandatory `*Type` suffix | Rename to `SkipReasonType`, add `type SkipReason = SkipReasonType` | Fixed |
| V-06 | `cli/cmd/commitin/enums.go` | 151 | `type ExclusionKind uint8` | Enum missing mandatory `*Type` suffix | Rename to `ExclusionKindType`, add `type ExclusionKind = ExclusionKindType` | Fixed |
| V-07 | `cli/cmd/commitin/enums.go` | 174 | `type MessageRuleKind uint8` | Enum missing mandatory `*Type` suffix | Rename to `MessageRuleKindType`, add `type MessageRuleKind = MessageRuleKindType` | Fixed |
| V-08 | `cli/cmd/commitin/enums.go` | 200 | `type FunctionIntelLanguage uint8` | Enum missing mandatory `*Type` suffix | Rename to `FunctionIntelLanguageType`, add `type FunctionIntelLanguage = FunctionIntelLanguageType` | Fixed |
| V-09 | `cli/cmd/commitin/workspace/source.go` | 25 | `type SourceKind uint8` | Enum missing mandatory `*Type` suffix | Rename to `SourceKindType`, add `type SourceKind = SourceKindType` | Fixed |
| V-10 | `cli/dbengine_old/operators.go` | 11 | `type SqlOperator string` | Enum missing mandatory `*Type` suffix | Rename to `SqlOperatorType`, add `type SqlOperator = SqlOperatorType` | Fixed |
| V-11 | `scripts/changelog/internal/runner/args.go` | 11 | `type Mode string` | Enum missing mandatory `*Type` suffix | Rename to `ModeType`, add `type Mode = ModeType` | Fixed |

---

## 5. Subtask Decomposition

- **Subtask 01:** `.lovable/plans/subtasks/140-enums/01-audit-and-scan-enums.md`
  - Targets: AST scan results, `linter-scripts/check-enum-guidelines.py`, `.agents/skills/cg-enums/skill.md`
  - Scope: Create `cg-enums` skill, enhance `check-enum-guidelines.py` with Go enum validation.
- **Subtask 02:** `.lovable/plans/subtasks/140-enums/02-refactor-go-enums-with-type-suffix.md`
  - Targets: `cli/cmd/commitin/enums.go`, `cli/cmd/commitin/workspace/source.go`, `cli/dbengine_old/operators.go`, `scripts/changelog/internal/runner/args.go`
  - Scope: Refactor all 11 enums with `*Type` suffix and backward-compatible aliases.
- **Subtask 03:** `.lovable/plans/subtasks/140-enums/03-verify-linters-and-milestone-consolidation.md`
  - Targets: Entire codebase verification, `check-enum-guidelines.py`, `check-enum-and-boolean.mjs`, `go vet ./...`.
  - Scope: Verify 0 enum violations repo-wide, record test inventory, consolidate plan milestone.

---

## 6. Verification Plan

1. **Enum Guidelines Linter:** `python linter-scripts/check-enum-guidelines.py` must pass with exit code 0.
2. **Enum & Boolean MJS Linter:** `node linter-scripts/check-enum-and-boolean.mjs` must pass with exit code 0.
3. **Go Compilation & Vet:** `go vet ./...` in `cli/` and `scripts/changelog/` must pass with 0 errors.
4. **Boolean Guidelines Linter:** `python linter-scripts/check-boolean-guidelines.py` must pass.
5. **Relative Path Linter:** `python linter-scripts/check-relative-paths.py` must pass.

---

## 7. Execution & Verification Summary

All 11 Go enum type definitions have been refactored to use the mandatory `*Type` suffix along with backward-compatible type aliases to guarantee zero breakage for callers.
- Tooling: Added `cg-enums` canonical skill (`.agents/skills/cg-enums/skill.md`).
- Linter: Enhanced `linter-scripts/check-enum-guidelines.py` with Go enum AST/regex checks.
- Verification:
  - `python linter-scripts/check-enum-guidelines.py`: PASS (0 violations)
  - `node linter-scripts/check-enum-and-boolean.mjs`: PASS (0 violations)
  - `go vet ./...` in `cli/`: PASS (0 errors)
  - `go vet ./...` in `scripts/changelog/`: PASS (0 errors)
  - `python linter-scripts/check-boolean-guidelines.py`: PASS (0 violations across 2,863 files)
  - `python linter-scripts/check-relative-paths.py`: PASS (0 violations across 6,827 files)
  - `python linter-scripts/check-nested-ifs.py`: PASS (0 violations across 2,863 files)
