# Plan 181: Servers-Clients (SC) Bash, Shell, Join, Nodes/List Commands, Rich Examples & SSH Connection Unification

> **Task Origin & Objective**:
> - User prompt highlighted missing `bash` and `shell` (`sh`) commands in `servers-clients` / `sc`.
> - Missing `join` (`koin`) and `list` / `nodes` / `joined` commands in `servers-clients` / `sc`.
> - Incomplete help text and zero practical examples for `servers-clients` / `sc`.
> - SQL error in `gitmap ssh exec`: `SQL logic error: no such table: SSHConnection (1) (at=db/sshconnection.go:58)`.
> - Minor version bump and release required.

---

## 1. Problem Analysis & RCA
- Missing bash/shell in `cli/cmd/clustersubcmd.go` and `cli/db/enums.go`.
- Missing join/add and nodes/list/joined routing in `cli/cmd/rootcore.go`.
- Missing `SSHConnection` table in `cli/store/migrations_ssh.go` and disconnect between `ssh_hosts` and `SSHConnection`.
- Sparse help documentation in `cli/helptext/servers-clients.md` and `sc.md`.

---

## 2. Subtasks
- **Subtask 01**: `cli/store/migrations_ssh.go`, `cli/db/sshconnection.go` - Database Schema & SSH Connection Unification.
- **Subtask 02**: `cli/db/enums.go`, `cli/cmd/clustersubcmd.go`, `cli/cluster/exec_bash.go`, `cli/cmd/rootcore.go`, `cli/cmd/clustercommand.go` - SC / Servers-Clients Subcommands & Dispatch.
- **Subtask 03**: `cli/helptext/servers-clients.md`, `cli/helptext/sc.md` - Comprehensive Help Documentation & Real-World Use-Case Examples.

---

## 3. Target Release
- Minor version bump: `v6.247.0` -> `v6.248.0`.
