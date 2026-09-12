# Plan 99: Constants & Enums Architecture Audit

## Executive Summary
This master architectural plan establishes repo-wide compliance with authoritative enum and constant conventions defined in `spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md`, `spec/02-coding-guidelines/06-ai-optimization/05-enum-naming-quick-reference.md`, and `.lovable/coding-guidelines.md`.

## Core Objectives
1. **Mandatory `*Type` Suffix on All Enums:** Ensure all enum definitions across Go, TypeScript, and Python end with `Type` (e.g. `MergeModeType`, `SqlOperatorType`, `ConfigKeyType`), providing type aliases where needed for clean backward compatibility.
2. **Elimination of Magic Strings & Numbers:** Centralize error codes, statuses, and protocol values into dedicated `constants/` and `enums/` packages.
3. **Total Ban on Raw Rune Casts:** Zero numeric rune conversions (`rune(10)`) or hardcoded control characters.
4. **Dedicated Definition Files:** Ensure all enums live in designated modules and types directories.

## Violation Ledger & Resolutions

| ID | File Path | Violation Type | Original Identifier | Resolved Architecture | Status |
|---|---|---|---|---|---|
| V-01 | `gitmap/vscodepm/mergemode.go` | Enum Missing `Type` Suffix | `type MergeMode uint8` | `type MergeModeType uint8` (with `type MergeMode = MergeModeType`) | RESOLVED |
| V-02 | `gitmap/dbengine/operators.go` | Enum Missing `Type` Suffix | `type SqlOperator string` | `type SqlOperatorType string` (with `type SqlOperator = SqlOperatorType`) | RESOLVED |
| V-03 | `gitmap/downloaderconfig/downloaderconfig.go` | Enum Missing `Type` Suffix | `type ConfigKey string` | `type ConfigKeyType string` (with `type ConfigKey = ConfigKeyType`) | RESOLVED |

## Subtasks Executed
- `01-task-enforce-enum-type-suffix-go.md`: Add `*Type` suffix to `MergeMode`, `SqlOperator`, and `ConfigKey` with backward-compatible aliases. (COMPLETE)
- `02-task-verify-linter-and-ast-constants.md`: Run `check-enum-guidelines.py`, `check-enum-and-boolean.mjs`, and AST constants checks. (COMPLETE)
- `03-task-verify-ci-local-runner.md`: Verify all 28 quality gates via `06-cicd-local-runner.py` and consolidate Plan 99. (COMPLETE)

## Verification Results
- `go build ./...`: PASS (Clean compile with zero errors)
- `python linter-scripts/check-enum-guidelines.py`: PASS
- `node linter-scripts/check-enum-and-boolean.mjs`: PASS
- `python 03-ai-scripts/06-cicd-local-runner.py --force --no-tests`: PASS (28/28 gates green)
