# Plan 305: Database Conventions & PascalCase Schema Views

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 07: `01-prompts/15-cg-execute/07-data-and-schema.md`

## Strategic Intent & System Reasoning
Enforce database conventions across SQLite split-databases by generating PascalCase singular views with `{TableName}Id` primary keys and positive boolean columns for legacy snake_case tables (`schedule_config`, `schedule_logs`, `ssh_history`, `spec_numbers`).

## Disjoint File Allocations
- **Worker 01:** `cli/store/schedule_split_db.go`
- **Worker 02:** `cli/store/migrations_ssh_hist.go`
- **Worker 03:** `cli/store/spec_numbers.go`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022183017-305-cg-execute-07-data-and-schema`).
2. [x] Author architecture spec `02-spec/21-app/305-data-and-schema/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`).
6. [x] Single atomic commit via `gitmap cpf "schema - enforce pascalcase views and primary key conventions"`.
