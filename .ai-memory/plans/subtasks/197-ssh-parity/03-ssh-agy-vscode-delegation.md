# Subtask 03: Remote AGY and VS Code Delegation via SSH

## Status: COMPLETED

## Summary of Accomplishments
1. Implemented `gitmap ssh agy <args>` in `cli/cmdssh/ssh_agy_cmd.go`:
   - Runs Antigravity CLI remotely or opens remote folders in AGY (`gitmap ssh agy open <folder>`).
   - Uses `apperror.NewExecutionError` and `apperror.WrapSimple` for full error and stack trace capture.
2. Implemented `gitmap ssh code <args>` in `cli/cmdssh/ssh_code_cmd.go`:
   - Attempts local VS Code SSH remote launch (`code --remote ssh-remote+<user@ip> <path>`).
   - Falls back to remote `code` CLI binary execution over SSH if local CLI is unavailable.
3. Routed `agy` and `code` subcommands through `dispatchPrimarySSH` in `cli/cmdssh/ssh.go`.
4. Enforced canonical <= 100 line limit on all delegation modules.
