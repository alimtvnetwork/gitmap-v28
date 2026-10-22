# Plan 307: Code Hygiene, Line Endings & Parameter Reduction

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 09: `01-prompts/15-cg-execute/09-code-hygiene.md`

## Strategic Intent & System Reasoning
Standardize repository code hygiene by normalizing line endings to Unix LF (`\n`), resolving vertical newline gaps before returns and after braces, eliminating explicit boolean equality checks (`== false`), and reducing parameter bloat via structured parameter objects.

## Disjoint File Allocations
- **Worker 01:** `cli/cmdpipeline/pipeline_history.go`
- **Worker 02:** `src/components/terminal/TerminalDrawer.tsx`
- **Worker 03:** `cli/cmdpipeline/types.go`, `cli/cmdpipeline/pipeline_sync_cache.go`, `cli/cmdpipeline/exports.go`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022183938-307-cg-execute-09-code-hygiene`).
2. [x] Author architecture spec `02-spec/21-app/307-code-hygiene/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`) and newline linter (`check-newline-styling.py`).
6. [x] Single atomic commit via `gitmap cpf "hygiene - normalize line endings, newline spacing, and parameter structs"`.
