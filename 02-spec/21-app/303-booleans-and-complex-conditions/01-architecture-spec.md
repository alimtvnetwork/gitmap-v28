# Architecture Spec: 303 - Booleans, Implicit Evaluation & Complex Conditions

## 1. Executive Summary & Intent
This specification formalizes the remediation of boolean naming and explicit comparison anti-patterns across the codebase in accordance with `01-prompts/15-cg-execute/04-booleans-and-complex-conditions.md` and repository guidelines (`02-spec/02-coding-guidelines/`).
All boolean variables must use affirmative prefixes (`is` or `has`), eliminate single-letter boolean parameters (`b bool`, `v bool`), and enforce implicit boolean evaluation (eliminating all instances of `== false`).

## 2. Invariants & Rules
1. **Implicit Evaluation:** Replace all `val == false` with `!val`.
2. **Positive Prefixes:** All booleans MUST begin with `is` or `has` only (`can`, `should`, `was`, `will` are banned).
3. **Parameter Clarity:** Replace single-letter booleans (`b bool`, `v bool`) with semantic parameters (`isEnabled bool`).
4. **Zero Nested Ifs:** All conditionals remain flattened with early returns.
5. **File Size Bounds:** All touched files must stay <= 500 lines.
6. **Zero Builds & Zero Tests:** FastGate AST verification only.

## 3. Disjoint Micro-Batch Allocation
- **Worker 01 (CmdAGY & Common Boolean Helpers):**
  - Files: `cli/cmdagy/agy_agentapi.go`, `cli/cmdagy/agy_agentapi_dispatch.go`, `cli/cmdagy/agy_conv_summaries.go`, `cli/cmdagy/agy_conv_transcript.go`, `cli/cmdagy/agy_ls_discovery.go`, `cli/cmdos/os_machine_alias_render.go`, `cli/tui/relformat.go`.
  - Action: Eliminate explicit `== false` checks; replace `b bool` / `v bool` parameters with `isEnabled bool`.
- **Worker 02 (CmdAntigravity & Automation):**
  - Files: `cli/cmdantigravity/prompt_help.go`, `cli/cmdantigravity/prompt_project.go`, `cli/cmdantigravity/prompt_cmd.go`, `cli/cmdantigravity/prompt_queue.go`, `cli/cmdautomation/clean_artifacts_cmd.go`.
  - Action: Eliminate explicit `== false` checks across prompt handlers and artifact cleanup.
- **Worker 03 (CmdSSH Subsystem):**
  - Files: `cli/cmdssh/ssh_cp.go`, `cli/cmdssh/ssh_macro_sync.go`, `cli/cmdssh/ssh_pass_cmd.go`, `cli/cmdssh/sshjoin_enroll.go`, `cli/cmdssh/ssh_target_loader.go`, `cli/cmdssh/ssh_export_all.go`, `cli/cmdssh/ssh_import_all.go`.
  - Action: Eliminate explicit `== false` checks in SSH file transfer, synchronization, key enroll, and credential handling.

## 4. Verification Evidence Contract
- FastGate AST inspection: `python 03-ai-scripts/50-fastgate.py` exit 0.
- Forbidden strings check: `python linter-scripts/check-forbidden-strings.py` exit 0.
