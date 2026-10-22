# Plan 306: React Frontend Architecture & Nested If Elimination

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 08: `01-prompts/15-cg-execute/08-react-frontend-guidelines.md`

## Strategic Intent & System Reasoning
Enforce React frontend architecture guidelines across `src/` by eliminating nested conditionals in event listeners and DOM inspectors, replacing explicit boolean equality checks (`=== false`) with implicit/affirmative evaluations, and renaming negative boolean identifiers to positive concepts.

## Disjoint File Allocations
- **Worker 01:** `src/components/docs/TabOrderMap.tsx`
- **Worker 02:** `src/components/terminal/TerminalDrawer.tsx`
- **Worker 03:** `src/pages/Commands.tsx`, `src/components/docs/CloneNextCommandBuilder.tsx`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022183459-306-cg-execute-08-react-frontend-guidelines`).
2. [x] Author architecture spec `02-spec/21-app/306-react-frontend-guidelines/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`).
6. [x] Single atomic commit via `gitmap cpf "react - enforce frontend guidelines and flatten nested conditionals"`.

> **Closed 2026-10-23 (task 190, subtask 06 drain):** plan complete per its own status record and the spec 02 audit (substantively complete with matching commits); moved from pending/ to completed/ so the pending folder reflects only genuinely open work.
