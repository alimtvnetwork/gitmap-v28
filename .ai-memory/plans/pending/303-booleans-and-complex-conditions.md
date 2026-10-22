# Plan 303: Booleans, Implicit Evaluation & Complex Conditions

## User Request (Verbatim)
`let's run all the prompts sequencetially from @[01-prompts/15-cg-execute/readme.md]`
-> Prompt 04: `01-prompts/15-cg-execute/04-booleans-and-complex-conditions.md`

## Strategic Intent & System Reasoning
Remediate boolean naming and explicit comparison anti-patterns across CmdAGY, CmdAntigravity, and CmdSSH, enforcing implicit evaluation (`!val` instead of `val == false`), positive boolean naming, and semantic parameter clarity.

## Disjoint File Allocations
- **Worker 01:** `cli/cmdagy/agy_agentapi.go`, `cli/cmdagy/agy_agentapi_dispatch.go`, `cli/cmdagy/agy_conv_summaries.go`, `cli/cmdagy/agy_conv_transcript.go`, `cli/cmdagy/agy_ls_discovery.go`, `cli/cmdos/os_machine_alias_render.go`, `cli/tui/relformat.go`
- **Worker 02:** `cli/cmdantigravity/prompt_help.go`, `cli/cmdantigravity/prompt_project.go`, `cli/cmdantigravity/prompt_cmd.go`, `cli/cmdantigravity/prompt_queue.go`, `cli/cmdautomation/clean_artifacts_cmd.go`
- **Worker 03:** `cli/cmdssh/ssh_cp.go`, `cli/cmdssh/ssh_macro_sync.go`, `cli/cmdssh/ssh_pass_cmd.go`, `cli/cmdssh/ssh_target_loader.go`, `cli/cmdssh/ssh_export_all.go`, `cli/cmdssh/ssh_import_all.go`

## Execution Steps
1. [x] Preflight, platform handshake, and SQLite Task DB initialization (`task-20261022182103-303-cg-execute-04-booleans-and-complex-conditions`).
2. [x] Author architecture spec `02-spec/21-app/303-booleans-and-complex-conditions/01-architecture-spec.md`.
3. [x] Register subtasks in SQLite Task DB.
4. [x] Spawn 3 parallel workers via `invoke_subagent` (`Model: 'flash_lite'`).
5. [x] Verify AST and style with FastGate (`50-fastgate.py`).
6. [x] Single atomic commit via `gitmap cpf "boolean - enforce implicit evaluation and positive naming"`.
