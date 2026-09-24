# Subtask 91-01: SSH Batch Common Join Engine (`ssh-join-common` / `sjc`)

Spec Reference: [02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md](../../../../02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md)
Parent Plan: [.ai-memory/plans/completed/91-ssh-join-common-scan-and-agy-rop.md](../../completed/91-ssh-join-common-scan-and-agy-rop.md)

## Objective
Implement `gitmap ssh-join-common` and short-form alias `gitmap sjc` allowing batch onboarding of multiple SSH nodes sharing common username and password.

## Functional Requirements
1. **Shorthand Octet Parser:**
   - Parse IP string like `192.168.1.3(w1),7(w2),12(w3)`.
   - First token `192.168.1.3(w1)` has full IPv4 `192.168.1.3`, alias `w1`, setting base subnet prefix `192.168.1.`.
   - Second token `7(w2)` is recognized as an octet since it contains no dots; appended to `192.168.1.` yielding `192.168.1.7`, alias `w2`.
   - Third token `12(w3)` resolves to `192.168.1.12`, alias `w3`.
   - If subsequent token contains full IP (e.g. `10.0.0.5(u1)`), updates the base prefix to `10.0.0.`.
   - If alias in `(...)` is omitted, defaults to `node-<IP>`.
2. **Batch Runner:**
   - Accepts common username as first argument (e.g. `administrator`).
   - Accepts `--pass <password>` flag. If omitted, prompts once securely for the common password.
   - Concurrently or sequentially dials each node, verifies authentication, executes deep OS detection, and enrolls host into host SQLite database (`installation.db` and `gitmap.db`).
   - Outputs clear progress indicators and an aligned summary table of joined nodes.
3. **CLI Dispatching:**
   - Registered under root commands `ssh-join-common`, `sjc`, and subcommands under `gitmap ssh join-common`.

## Files to Create/Modify
- `cli/cmdssh/sshjoin_common_types.go` [NEW]
- `cli/cmdssh/sshjoin_common_parser.go` [NEW]
- `cli/cmdssh/sshjoin_common.go` [NEW]
- `cli/cmdssh/sshjoin_common_cmd.go` [NEW]
- `cli/cmdssh/ssh_common_join_test.go` [NEW]
- `cli/cmd/root.go` & `cli/cmd/rootdispatch.go` [MODIFY]
