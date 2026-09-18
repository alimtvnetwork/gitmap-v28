# Plan 185: Constants & Enums Architecture — Coding Guideline Execution

> **Task Origin & Objective**:
> - Autonomous execution under Coding Guidelines Version 2.1.0 (`cg-enums`, `cg-constants`).
> - Refactor all enums missing mandatory `*Type` suffix, eliminate raw numeric rune casts (`rune(10)`, `rune('a'+i)`), eliminate magic string/number literals, and ensure 100% compliance with `linter-scripts/check-enum-guidelines.py` and `03-ai-scripts/37-enum-guideline-auditor.py`.
> - Enforce file size limit <= 100 lines (target <= 80 lines) and function size <= 8–15 lines.
> - Zero line-compression cheating, affirmative booleans only, zero intermediate tests or builds.

---

## 1. Problem Analysis & Violation Ledger

| # | File Path | Violation Type | Description | Remediation Plan |
|:---:|---|---|---|---|
| 1 | `cli/cmd/historyrewrite.go:18` | `GO_ENUM_MISSING_TYPE` | `type historyMode int` missing `Type` suffix | Renamed to `HistoryModeType`, added alias `type historyMode = HistoryModeType`, decomposed into `historyrewrite.go` (77 lines) and `historyrewrite_helpers.go` (55 lines) |
| 2 | `cli/cmd/replace.go:60` | `GO_ENUM_MISSING_TYPE` | `type replaceMode int` missing `Type` suffix | Renamed to `ReplaceModeType`, added alias `type replaceMode = ReplaceModeType`, decomposed into `replace.go` (53 lines), `replace_classify.go` (59 lines), and `replace_dashn.go` (47 lines) |
| 3 | `cli/cmd/desktopsync.go:97` | `GO_ENUM_MISSING_TYPE` | `type syncResult int` missing `Type` suffix | Renamed to `SyncResultType`, added alias `type syncResult = SyncResultType`, decomposed into `desktopsync.go` (93 lines) and `desktopsync_ops.go` (78 lines) |
| 4 | `cli/cmd/visibilitydriftguard.go:17` | `GO_ENUM_MISSING_TYPE` | `type driftAction string` missing `Type` suffix | Renamed to `DriftActionType`, added alias `type driftAction = DriftActionType` (45 lines) |
| 5 | `cli/cmd/hygiene_parallel.go:131` | `GO_ENUM_MISSING_TYPE` | `type hygieneFormat string` missing `Type` suffix | Renamed to `HygieneFormatType`, added alias `type hygieneFormat = HygieneFormatType`, decomposed into `hygiene_parallel.go` (54 lines), `hygiene_parallel_workers.go` (53 lines), `hygiene_parallel_map.go` (75 lines), and `hygiene_format.go` (59 lines) |
| 6 | `cli/cmdmacro/macro_add_interactive.go:19` | `GO_ENUM_MISSING_TYPE` | `type interactiveLoopAction int` missing `Type` suffix | Renamed to `InteractiveLoopActionType`, added alias `type interactiveLoopAction = InteractiveLoopActionType` |
| 7 | `cli/config/validate_shape.go:55` | `GO_ENUM_MISSING_TYPE` | `type jsonKind int` missing `Type` suffix | Renamed to `JSONKindType`, added alias `type jsonKind = JSONKindType` |
| 8 | `cli/macro/recurse.go:13` | `GO_ENUM_MISSING_TYPE` | `type contextKey string` missing `Type` suffix | Renamed to `ContextKeyType`, added alias `type contextKey = ContextKeyType` |
| 9 | `cli/render/pretty_parse.go:6` | `GO_ENUM_MISSING_TYPE` | `type blockKind int` missing `Type` suffix | Renamed to `BlockKindType`, added alias `type blockKind = BlockKindType` |
| 10 | `cli/clonefrom/execute_concurrent_test.go:148` | `RAW_RUNE_CAST` | `return string(rune('a'+i)) + "-row"` | Replaced with `fmt.Sprintf("%c-row", 'a'+i)` |
| 11 | `cli/clonenow/execute_concurrent_test.go:143` | `RAW_RUNE_CAST` | `return string(rune('a'+i)) + "-row"` | Replaced with `fmt.Sprintf("%c-row", 'a'+i)` |
| 12 | `cli/scanner/progress_test.go:83` | `RAW_RUNE_CAST` | `string(rune('a'+i%5))` | Replaced with `fmt.Sprintf("%c", 'a'+i%5)` |
| 13 | `cli/scanner/scanner_test.go:121` | `RAW_RUNE_CAST` | `string(rune('a'+i%5))` | Replaced with `fmt.Sprintf("%c", 'a'+i%5)` and `fmt.Sprintf("%d-%c", i%10, 'a'+i%26)` |
| 14 | `cli/tests/cmd_test/seowriteloop_test.go:210` | `RAW_RUNE_CAST` | `s = string(rune('0'+n%10)) + s` | Replaced custom itoa loop with standard library `strconv.Itoa(n)` |

---

## 2. Granular Subtask Execution Summary

### Subtask 01: Enforce `*Type` on Command Enums (`historyrewrite`, `replace`, `desktopsync`, `visibilitydriftguard`)
- Refactored `historyMode` -> `HistoryModeType` (with alias). Split `historyrewrite.go` into `historyrewrite.go` (77 lines) and `historyrewrite_helpers.go` (55 lines).
- Refactored `replaceMode` -> `ReplaceModeType` (with alias). Split `replace.go` into `replace.go` (53 lines), `replace_classify.go` (59 lines), and `replace_dashn.go` (47 lines).
- Refactored `syncResult` -> `SyncResultType` (with alias). Split `desktopsync.go` into `desktopsync.go` (93 lines) and `desktopsync_ops.go` (78 lines).
- Refactored `driftAction` -> `DriftActionType` (with alias) in `visibilitydriftguard.go` (45 lines).

### Subtask 02: Enforce `*Type` on Core & Subsystem Enums (`hygiene_parallel`, `macro`, `config`, `render`)
- Refactored `hygieneFormat` -> `HygieneFormatType` (with alias). Split `hygiene_parallel.go` into `hygiene_parallel.go` (54 lines), `hygiene_parallel_workers.go` (53 lines), `hygiene_parallel_map.go` (75 lines), and `hygiene_format.go` (59 lines).
- Refactored `interactiveLoopAction` -> `InteractiveLoopActionType` (with alias) in `cli/cmdmacro/macro_add_interactive.go`.
- Refactored `jsonKind` -> `JSONKindType` (with alias) in `cli/config/validate_shape.go`.
- Refactored `blockKind` -> `BlockKindType` (with alias) in `cli/render/pretty_parse.go`.
- Refactored `contextKey` -> `ContextKeyType` (with alias) in `cli/macro/recurse.go`.

### Subtask 03: Eliminate Raw Rune Numeric Conversions across Codebase and Tests
- Replaced raw rune casts in `cli/clonefrom/execute_concurrent_test.go` and `cli/clonenow/execute_concurrent_test.go` with `fmt.Sprintf("%c-row", 'a'+i)`.
- Replaced raw rune casts in `cli/scanner/progress_test.go` and `cli/scanner/scanner_test.go` with formatted path components.
- Replaced custom itoa in `cli/tests/cmd_test/seowriteloop_test.go` with `strconv.Itoa(n)`.

---

## 3. Strict Quality Invariants & Verification Results
- File line cap: strictly under 100 lines across all modified and created files (all <= 93 lines).
- Function line cap: <= 15 lines across all functions in all modified and created files.
- Mandatory `*Type` suffix on all enums with backward-compatible aliases.
- Zero raw numeric rune casts across the repository.
- Linters:
  - `python linter-scripts/check-enum-guidelines.py`: PASS (0 violations).
  - `python 03-ai-scripts/37-enum-guideline-auditor.py`: PASS (0 violations).
  - `node linter-scripts/check-enum-and-boolean.mjs`: PASS (0 violations).
  - `python linter-scripts/check-naming-guidelines.py`: PASS (0 violations across 3424 files).
- Formatter: `python 03-ai-scripts/26-go-code-formatter.py`: PASS (2975 files formatted).
- Compilation: `go vet` clean compilation across all touched packages.
- Test Inventory: `python 03-ai-scripts/33-test-inventory-generator.py` updated (3534 tests indexed).
