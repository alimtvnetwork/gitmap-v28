# Plan 301: Error Management & Strongly-Typed Exit Codes

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 02: `01-prompts/15-cg-execute/02-error-management.md`

## Strategic Intent & System Reasoning
Standardize error handling architecture across CLI command handlers, replacing magic integer exit codes with strongly-typed `cliexit.ExitCodeType` constants and specialized helpers. Ensure errors bubble up cleanly and exit codes are typed and audited.

## Disjoint File Allocations
- **Worker 01:** `cli/cmdzip/zip.go`, `cli/cmdzip/unzipcompact.go`, `cli/cmdzip/zipgroup.go`, `cli/cmdzip/zipgroupcreate.go`, `cli/cmdzip/zipgroupshow.go`, `cli/cmdzip/zipgroupops.go`
- **Worker 02:** `cli/cmdvscode/vscode_cmd.go`, `cli/cmdvscode/vscodeworkspace.go`, `cli/helpdoc/print.go`
- **Worker 03:** `cli/cmdauditlegacy/auditlegacy.go`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022180856-301-cg-execute-02-error-management`).
2. [x] Author architecture spec `02-spec/21-app/301-error-management-and-typed-exit-codes/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`).
6. [x] Single atomic commit via `gitmap cpf "error - standardize typed exit codes and apperror sovereignty"`.
