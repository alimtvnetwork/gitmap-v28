# Plan 309: Test Naming & Three-Part Semantic Conventions

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 11: `01-prompts/15-cg-execute/11-testing-and-coverage.md`

## Strategic Intent & System Reasoning
Standardize test function declarations across Go packages adhering to the mandatory three-part semantic convention (`Test{Unit}_{Scenario}_{ExpectedOutcome}`) in `cli/uipref/uipref_test.go`, `cli/usercontext/project_binder_test.go`, and `cli/workspacesync/path_guard_test.go`.

## Disjoint File Allocations
- **Worker 01:** `cli/uipref/uipref_test.go`
- **Worker 02:** `cli/usercontext/project_binder_test.go`
- **Worker 03:** `cli/workspacesync/path_guard_test.go`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022184545-309-cg-execute-11-testing-and-coverage`).
2. [x] Author architecture spec `02-spec/21-app/309-testing-and-coverage/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [ ] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [ ] Verify AST and style with FastGate (`50-fastgate.py`).
6. [ ] Single atomic commit via `gitmap cpf "test - standardize three-part semantic test naming conventions"`.

> **Closed 2026-10-23 (task 190, subtask 06 drain):** plan complete per its own status record and the spec 02 audit (substantively complete with matching commits); moved from pending/ to completed/ so the pending folder reflects only genuinely open work.
