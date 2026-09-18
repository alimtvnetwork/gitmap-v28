# Subtask 02: SSH Multi-Machine Join & Open Machine Discovery Verification

## Scope
- Inspect `cli/cmdssh/sshjoin.go`, `cli/cmdssh/sshscan.go`, `cli/cmdssh/ssh_liveness.go`, and `cli/cmdssh/sshrecall.go`.
- Verify registering and joining multiple machines, alias routing, and credential caching.
- Verify liveness scanning (`gitmap ssh scan`, `gitmap ssh check`) to detect open machines (TCP port 22, latency probing, cache invalidation/refresh).
- Confirm host recall and interactive table display for discovered nodes.

## Acceptance Criteria
- [x] Joining multiple machines properly stores host configuration in SQLite SSH store (`ssh_hosts`, `ssh_history`).
- [x] Direct `gitmap ssh check` / `health` / `ping` command added to `dispatchPrimarySSH` in `cli/cmdssh/ssh.go`.
- [x] Scan and check accurately test port 22 liveness across registered nodes with latency and colored status.
- [x] Terminal help text details `join`, `scan`, `check`, and `recall` workflows with concrete examples.

## Completed Changes
- Connected `check`, `health`, `ping` subcommands under `gitmap ssh` directly to `RunSJStatus` (`ExecuteHealthCheck`).
- Updated `MsgSSHAvailableCommands` in `cli/constants/constants_ssh.go` to include `ssh check [target]`.
- Updated `cli/helptext/ssh.md` with table row and examples for checking open port 22 connectivity.
- Added `ssh join` and `ssh check` definitions to `src/data/commands.ts`.
