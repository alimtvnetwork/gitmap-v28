# Plan 304: Constants & Strongly-Typed Enums

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 06: `01-prompts/15-cg-execute/06-constants-and-enums.md`

## Strategic Intent & System Reasoning
Enforce mandatory `*Type` suffixes on all custom enum types in Go and eliminate raw numeric rune casts, preserving full backward compatibility via type aliases while ensuring strict adherence to cross-language enum guidelines.

## Disjoint File Allocations
- **Worker 01:** `cli/cmdapps/types.go`, `cli/cloneconcurrency/pull_concurrency.go`, `cli/cmddownload/types.go`, `cli/usercontext/types.go`
- **Worker 02:** `cli/cmdfixreleasetags/types.go`, `cli/cmdfoldertree/treesearch_match.go`, `cli/cmdos/os_dock_types.go`, `cli/cmdupdate/update_ui_types.go`
- **Worker 03:** `cli/diag/types.go`, `cli/termgeneric/header.go`, `cli/termout/themes.go`, `cli/heatmap/heatmap_query.go`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022182444-304-cg-execute-06-constants-and-enums`).
2. [x] Author architecture spec `02-spec/21-app/304-constants-and-enums/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`) and `check-enum-guidelines.py`.
6. [x] Single atomic commit via `gitmap cpf "enums - enforce type suffix on domain enums and eliminate rune casts"`.

> **Closed 2026-10-23 (task 190, subtask 06 drain):** plan complete per its own status record and the spec 02 audit (substantively complete with matching commits); moved from pending/ to completed/ so the pending folder reflects only genuinely open work.
