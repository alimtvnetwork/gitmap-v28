# Plan 308: Coding Style, Vertical Newline Rhythm & Return New Line Concept

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 10: `01-prompts/15-cg-execute/10-style-guidelines.md`

## Strategic Intent & System Reasoning
Enforce repository coding style and Return New Line rules (R13-R16) repository-wide, eliminating all 14 newline spacing violations across `src/pages/AGYPrompts.tsx`, `src/pages/Database.tsx`, and `src/pages/MergeAIPage.tsx` to achieve 100% full-scan pass on `check-newline-styling.py --all`.

## Disjoint File Allocations
- **Worker 01:** `src/pages/AGYPrompts.tsx`
- **Worker 02:** `src/pages/Database.tsx`
- **Worker 03:** `src/pages/MergeAIPage.tsx`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022184239-308-cg-execute-10-style-guidelines`).
2. [x] Author architecture spec `02-spec/21-app/308-style-guidelines/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`) and full newline linter (`check-newline-styling.py --all`).
6. [x] Single atomic commit via `gitmap cpf "style - enforce return new line and closing brace vertical spacing"`.
