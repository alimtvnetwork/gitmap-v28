# Plan 197: Terminal SSH Execution, Liveness Scan, Remote Install & Cluster Parity

## Status: COMPLETED

## 1. Overview & Problem Statement
The user reported that `gitmap ssh exec <cmd>` was not executing properly across remote hosts (showing no output or failing with `No password or key configured` on machines enrolled with `gitmap ssh join user@ip alias`). Additionally:
1. Offline machines should not cause slow hangs during execution; running machines should execute immediately while offline machines are skipped with clear notices.
2. An in-memory reachability cache with TTL (45s) was requested so repeated executions skip slow timeouts, alongside an active `gitmap ssh scan` command.
3. Remote GitMap installation was missing from the SSH section: `gitmap ssh install` (installing if missing, updating if present) and `gitmap ssh update` were required.
4. Delegation for Antigravity (`gitmap ssh agy`) and VS Code (`gitmap ssh code`) was needed to open remote folders or run tools with `AppError` diagnostics on failure.
5. The relationship between `ssh`, `cluster`, and `sc` (`servers-clients`) was under-documented, requiring an interactive architecture comparison table (`gitmap ssh compare`) and enriched UI help parity.

---

## 2. Changes Made

### A. Liveness Probing & Caching
- **`cli/cmdssh/ssh_liveness.go`** (89 lines):
  - In-memory thread-safe `sync.RWMutex` cache (`map[string]livenessEntry`) with 45-second TTL.
  - `CheckNodeLiveness`, `CheckConnLiveness`, and `InvalidateLivenessCache`.
- **`cli/cmdssh/sshexec.go`**:
  - Checks `CheckConnLiveness` before dialing SSH client.
  - Skips offline nodes cleanly: `[alias|ip] OFFLINE (skipped: reason)`.
- **`cli/cmdssh/ssh_scan_cmd.go`** (44 lines) & **`cli/cmdssh/ssh_scan_render.go`** (69 lines):
  - `gitmap ssh scan`: probes registered fleet hosts in parallel, updates liveness cache, and prints formatted summary table via `termtable.PrintTable`.

### B. Smart Key Discovery & Advice
- **`cli/cmdssh/ssh_keys.go`** (62 lines):
  - `findDefaultUserSSHKey`: automatically detects standard SSH keys (`~/.ssh/id_ed25519`, `~/.ssh/id_rsa`, `~/.ssh/id_ecdsa`).
  - `connectWithDefaultKey`: auto-authenticates without requiring explicit key flags.
  - `formatMissingAuthAdvice`: renders actionable steps if credentials are completely absent.

### C. Remote GitMap Install & Update
- **`cli/cmdssh/ssh_install_remote.go`** (70 lines) & **`cli/cmdssh/ssh_install_actions.go`** (46 lines):
  - `gitmap ssh install [target]`: detects whether `gitmap` is installed via `isGitmapPresent` (`gitmap --version`).
  - If missing: executes official platform installation one-liner (curl bash or PowerShell).
  - If present: triggers `gitmap update`.
- **`cli/cmdssh/ssh_update_remote.go`** (64 lines):
  - `gitmap ssh update [target]`: updates GitMap across target nodes or all fleet machines.
- **`cli/cmdssh/ssh_target_loader.go`** (39 lines):
  - Shared connection loader filtering by target alias, IP, or `all`.

### D. AGY & VS Code Remote Delegation
- **`cli/cmdssh/ssh_agy_cmd.go`** (84 lines):
  - `gitmap ssh agy <args>`: runs Antigravity CLI or opens remote folders in AGY with `apperror.NewExecutionError` diagnostics.
- **`cli/cmdssh/ssh_code_cmd.go`** (94 lines):
  - `gitmap ssh code <args>`: opens remote folder in VS Code via SSH Remote (`code --remote ssh-remote+<user@ip> <path>`).

### E. Architecture Comparison & UI Help Parity
- **`cli/cmdssh/ssh_compare_table.go`** (82 lines):
  - Implements `gitmap ssh compare` (alias `matrix`) rendering a detailed table comparing `ssh` vs `cluster` vs `sc`.
- **`cli/cmdssh/ssh_fallback.go`** (76 lines):
  - Extracted legacy fallback handlers from `ssh.go` to keep `ssh.go` under 100 lines (82 lines).
- **`cli/constants/constants_ssh.go`**:
  - Updated `MsgSSHAvailableCommands` to list `scan`, `install`, `update`, `agy`, `code`, and `compare`.
- **`cli/helptext/ssh.md`**, **`cli/helptext/cluster.md`**, and **`cli/helptext/sc.md`**:
  - Added comparative architecture explanation, workflow differences, and command syntax examples.

### F. Unit Tests
- **`cli/cmdssh/ssh_liveness_test.go`** (45 lines)
- **`cli/cmdssh/ssh_keys_test.go`** (39 lines)
- **`cli/cmdssh/ssh_compare_table_test.go`** (31 lines)

---

## 3. Verification & Compliance
- **File Size**: All 14 new and modified Go files are strictly <= 100 lines.
- **Nested Ifs**: Zero nested if violations (depth <= 1) verified via `check-nested-ifs.py` across 3093 files.
- **Enums & Booleans**: Zero violations verified via `check-enum-and-boolean.py` across 2330 files.
- **Compiler**: Go Compile Gate passed 100% via `06-cicd-local-runner.py --filter "Compile"`.
- **Zero Test Run Rule**: No full routine test runs were executed during loops.
