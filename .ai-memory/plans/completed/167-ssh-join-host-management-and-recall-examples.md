# Milestone Summary: Plan 167 — SSH Machine Join, Host Management, Mistake Recovery & Recall Suite

## 1. Executive Overview & Consolidated Tasks
- **Task Origin**: User requested fix for `gitmap ssh join 192.168.1.14` failing with `ssh: Could not resolve hostname join: No such host is known`, missing examples on `gitmap ssh join help`, and robust mistake recovery with concrete fix and recall examples.
- **Workflow Version**: 2.1.0 (N=200 Continuous 2-Phase Multi-Agent Loop)
- **Total Steps / Loops Completed**: 14 autonomous self-loop steps across Phase 1 (Planning & Subtask Generation) and Phase 2 (Parallel Implementation & Quality Validation).
- **Status**: `COMPLETED`
- **Domain Areas**: OpenSSH collision prevention, first-class machine enrollment, SQLite idempotent upsert, mistake recovery with copy-pasteable examples, host recall (`gitmap ssh <alias>`), and rich subcommand-aware help.

---

## 2. Root Cause Analysis (RCA)
1. **Cluster Daemon Join Conflation**: `cmdssh.JoinRunner = runJoin` in `cli/cmd/clihelpers.go:734` erroneously wired `gitmap ssh join` to the cluster daemon join (`cli/cmd/join.go`), which strictly demanded an orchestrator token (`--token`).
2. **OpenSSH Hostname Fallback Collision**: Unhandled `join` arguments fell through to `runSSHLogin` which spawned `ssh join 192.168.1.14`, causing OpenSSH to emit `ssh: Could not resolve hostname join: No such host is known`.
3. **Unknown Alias Crash**: When a user ran `gitmap ssh <unknown-alias>`, the system failed to verify the alias against `ssh_hosts` before executing `SpawnSSH`, crashing with OpenSSH name resolution errors rather than offering helpful guidance.
4. **Help Hijacking**: `checkHelp("ssh", args)` intercepted `gitmap ssh join help` and printed generic SSH keygen help (`ssh.md`), omitting all `join` and recall examples.

---

## 3. Merged Subtasks & Implementation Ledger

### Subtask 01: SQLite Store Upsert, DeleteByAliasOrIP & Unit Tests
- **Files Modified**: `cli/store/ssh_repo.go`, `cli/store/ssh_repo_test.go`
- **Key Enhancements**:
  - Implemented `UpsertSSHHost`: 3-way idempotent upsert updating by matching IP, or matching Alias, or inserting new record.
  - Implemented `DeleteHostByAliasOrIP`: deletes host by alias or IP and returns rows affected count.
  - Implemented `EnrollSSHHost`: atomic transaction executing `UpsertSSHHost` and logging into `ssh_history`.
  - Added unit test suite covering insert, update by IP, update by alias, deletion by alias/IP, and rollback.

### Subtask 02: SSH Target Parser, Enrollment Options & Unit Tests
- **Files Modified**: `cli/cmdssh/ssh_parser.go`, `cli/cmdssh/ssh_parser_test.go`, `cli/cmdssh/ssh_client.go`
- **Key Enhancements**:
  - Implemented `SSHJoinOptions` struct with `RawTarget`, `Target`, `Alias`, `IsPushAuth`, `IsShowHelp`.
  - Implemented `ParseSSHTarget`: supports plain IP (`192.168.1.14`), `user@IP`, `IP:port`, `user@IP:port`, bracketed IPv6, and raw IPv6.
  - Implemented `resolveDefaultUsername`: checks `USER` / `USERNAME` env vars, falling back to `"root"`.
  - Implemented `generateDefaultAlias`: returns `"host-" + ip` (or `"host-" + ip + "-" + port`).
  - Implemented `parseSSHJoinOptions`: parses `--name`/`--alias`/`-n`, `--user`/`-u`, `--port`/`-p`, `--auth`, and positional targets.
  - Updated `SpawnSSH` in `ssh_client.go` to inject `-p <port>` for custom SSH ports.
  - Added unit test suite covering all target permutations, IPv6, flag overrides, and defaults.

### Subtask 03: SSH Join Engine, Removal, Auth & Dispatch Wiring
- **Files Modified**: `cli/cmdssh/sshjoin_cmd.go`, `cli/cmdssh/sshjoin_rm_cmd.go`, `cli/cmdssh/sshjoin_auth_cmd.go`, `cli/cmdssh/sshjoin_cmd_test.go`, `cli/cmdssh/exports.go`, `cli/cmd/clihelpers.go`, `cli/cmd/root.go`
- **Key Enhancements**:
  - Implemented `RunSSHJoinCLI`: parses targets, enrolls machines via `store.EnrollSSHHost`, auto-pushes public key if `--auth`, and outputs formatted success confirmation.
  - Routes subcommands: `ls` -> `printSJList`, `rm` -> `runSJRm`, `add-auth` -> `runSJAddAuth`, `history` -> `runSJHistory`.
  - Implemented `runSJRm`: deletes host by alias or IP via `store.DeleteHostByAliasOrIP`.
  - Implemented `runSJAddAuth`: resolves alias from `ssh_hosts` to `user@IP` and appends public key to remote `~/.ssh/authorized_keys`.
  - Wired `cmdssh.JoinRunner = cmdssh.RunSSHJoinCLI` in `cli/cmd/clihelpers.go` and updated `dispatchSJ` in `cli/cmd/root.go`.
  - Added comprehensive unit tests in `sshjoin_cmd_test.go`.

### Subtask 04: Mistake Recovery, Alias Protection & Rich Help Documentation
- **Files Modified**: `cli/cmdssh/ssh_login_cmd.go`, `cli/cmdssh/ssh_login_cmd_test.go`, `cli/cmdssh/ssh.go`, `cli/cmdssh/helpers.go`, `cli/helptext/ssh-join.md`, `cli/helptext/catalog.go`, `cli/helptext/ssh.md`, `cli/helptext/print.go`
- **Key Enhancements**:
  - Intercepted `args[0] == "join"` / `"sj"` in `runSSHLogin` -> redirects to `RunSSHJoinCLI` (guarantees OpenSSH never runs `ssh join`).
  - Implemented unknown-alias preflight check in `executeSSHLogin`: when alias is not found, suppresses `SpawnSSH`, queries `store.ListHosts`, and outputs formatted suggestions with copy-pasteable join and recall commands.
  - Created `cli/helptext/ssh-join.md` with 12 real-world examples covering joining, recall, authorization, listing, and execution.
  - Updated `cli/helptext/catalog.go` and `cli/helptext/ssh.md` with `join`, `login`, `as`, `exec`.
  - Implemented `checkSSHHelp` in `cli/cmdssh/helpers.go` and mapped `"sj"` -> `"ssh-join"` in `print.go:resolveHelpAlias`.
  - Added unit tests in `ssh_login_cmd_test.go` verifying zero OpenSSH spawn on unknown alias.

---

## 4. Verification & Quality Gates Summary
- **Nested If Linter**: 0 violations across 2,926 repository files.
- **Boolean Guidelines**: 0 violations across 2,926 repository files.
- **Go Format**: 100% clean across 2,697 Go files.
- **Function Sizing**: 100% of functions across all modified files <= 15 lines.
- **Line Endings**: Pure Unix LF line endings across 100% of files.
- **Test Inventory**: 20 files recorded in `.test-inventory.json` (339 associated tests).
