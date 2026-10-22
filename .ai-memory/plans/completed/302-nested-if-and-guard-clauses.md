# Plan 302: Nested If Elimination & Guard Clauses

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 03: `01-prompts/15-cg-execute/03-nested-if-and-guard-clauses.md`

## Strategic Intent & System Reasoning
Flatten all nested conditionals across target files, enforcing guard clauses, early returns, helper extraction, and <= 8–15 line function caps to ensure clean, readable, single-level cyclomatic complexity.

## Disjoint File Allocations
- **Worker 01:** `cli/fsutil/recursive_top_level.go`, `cli/fsutil/scan_exclusions.go`, `cli/gitutil/exclude.go`
- **Worker 02:** `cli/helpdoc/examples.go`, `cli/cmd/path_cli.go`, `cli/cmd/folder/folder.go`, `cli/cmd/folder/parse_args.go`
- **Worker 03:** `cli/cmdautofix/fix.go`, `cli/cmdautofix/scan.go`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022181700-302-cg-execute-03-nested-if-elimination`).
2. [x] Author architecture spec `02-spec/21-app/302-nested-if-and-guard-clauses/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`).
6. [x] Single atomic commit via `gitmap cpf "style - eliminate nested ifs and invert guard clauses"`.

> **Closed 2026-10-23 (task 190, subtask 06 drain):** plan complete per its own status record and the spec 02 audit (substantively complete with matching commits); moved from pending/ to completed/ so the pending folder reflects only genuinely open work.
