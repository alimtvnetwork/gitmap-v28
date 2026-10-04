# Plan 80: GitMap Push-Fix Command & Authentication Self-Healing Engine

> **Plan Status:** Active  
> **Traceability IDs:** Subtask 210-01 .. Subtask 210-05  
> **Spec Reference:** [02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md](../../02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md)  
> **Target Subsystems:** `cli/cmdpull`, `cli/cmdpush`, `cli/cmdssh`, `cli/store`, `cli/constants`, `cli/theme`  

---

## 1. Architectural Context & Subtask Mapping

| Subtask ID | Focus Area | Target Deliverable | Status |
| :--- | :--- | :--- | :--- |
| **Subtask 210-01** | `cli/cmd`, `cli/cmdpull`, `cli/constants` | Core push-fix command routing (`gitmap push-fix`, `gitmap push fix`, `pf`), pre-flight git inspection, unpushed commit counter, anti-hang environment guards (`GIT_TERMINAL_PROMPT=0`, `GCM_INTERACTIVE=never`) | In-Progress |
| **Subtask 210-02** | `cli/cmdssh`, `cli/store`, `cli/cmdpull` | SSH handshake probe with timeout, dual-profile SSH config writing (`~/.ssh/config` + Windows User profile), dashed bogus keys cleanup (`DELETE FROM SshKey WHERE Name LIKE '-%'`), and transport auto-conversion via `SetRepoIdentifiedTransport` | In-Progress |
| **Subtask 210-03** | `cli/cmdpull`, `cli/cloner` | Diverged branch auto-rebase recovery (`git pull --rebase --autostash`), conflict detection and rollback, upstream branch tracking auto-binding (`git push -u origin <branch>`) | Pending |
| **Subtask 210-04** | `cli/cmdpull`, `cli/clicolors`, `cli/theme` | High-contrast ANSI terminal HUD card formatting, commit metadata display, probe and push latency benchmarking metrics | Pending |
| **Subtask 210-05** | `cli/cmd`, `cli/helptext`, documentation | Command registration validation, CLI help text integration, README command table synchronization, and regression safety | Pending |

---

## 2. 5-Subtask Detailed Breakdown

### Subtask 210-01: Architecture & Push-Fix Pipeline Core
- **Target Files:**
  - `cli/constants/constants_cli.go`: Register `CmdPushFix`, `CmdPushFixAlias`, `CmdPushFixShort` (`pf`).
  - `cli/cmd/dispatchcommittransfer.go` / `cli/cmd/rootcore.go`: Dispatch `push-fix` and subcommand `push fix`.
  - `cli/cmdpull/push_fix.go`: Implement `RunPushFix(args []string) error` with CWD repository verification, active branch detection, unpushed commits count calculation via `git rev-list`.
  - Inject anti-hang environment variables (`GIT_TERMINAL_PROMPT=0`, `GCM_INTERACTIVE=never`, `GIT_SSH_COMMAND=ssh -o BatchMode=yes -o ConnectTimeout=5`).

### Subtask 210-02: SSH & HTTPS Authentication Remediation Engine
- **Target Files:**
  - `cli/cmdssh/auth_probe.go`: Implement non-interactive SSH handshake probe (`ssh -T git@github.com`) and HTTPS remote probe (`git ls-remote`).
  - `cli/store/ssh_keys.go`: Add `PurgeBogusDashedSSHKeys()` to remove records where `Name LIKE '-%'` or `Name = ''`.
  - `cli/cmdssh/sshconfig.go`: Update managed block generator to always include canonical `Host github.com` default block, and synchronize across dual profiles on Windows (`~/.ssh/config` and `%SystemDrive%\Users\%USERNAME%\.ssh\config`).
  - `cli/cmdpull/push_fix.go`: Trigger transport auto-conversion from HTTPS to SSH when HTTPS credentials fail but SSH keys succeed, updating `db.SetRepoIdentifiedTransport(url, "ssh")`.

### Subtask 210-03: Diverged Branch Auto-Rebase & Upstream Tracking Engine
- **Target Files:**
  - `cli/cmdpull/push_fix.go`: Check if active branch has an upstream tracking reference (`git rev-parse --abbrev-ref @{u}`).
  - If upstream tracking is missing, pass `-u` / `--set-upstream origin <branch>` to `git push`.
  - On non-fast-forward push rejection (`fetch first`, `non-fast-forward`, `[rejected]`), run `git pull --rebase --autostash origin <branch>`.
  - If rebase encounters unresolvable conflicts, abort rebase cleanly (`git rebase --abort`) and emit structured conflict guidance.

### Subtask 210-04: Terminal HUD Telemetry & ANSI Metric Cards
- **Target Files:**
  - `cli/cmdpull/push_fix_hud.go`: Construct ANSI HUD metric card with box borders, showing repository slug, branch, head commit SHA, pushed commit count, transport protocol, and millisecond latency timers.
  - Apply high-contrast theme styling (Neon Green `#00FF7F`, Bright Yellow `#FFD700`, Cyan `#8BE9FD`).

### Subtask 210-05: End-to-End CLI Verification & Documentation
- **Target Files:**
  - `cli/helptext/helptext.go`: Add comprehensive help page for `gitmap push-fix` explaining autonomous remediation modes, auto-rebase, and transport switching.
  - Verify zero-regression against existing `gitmap push`, `gitmap pull`, and `gitmap fix-auth` commands.
  - Keep documentation updated across repository indices.

---

## 3. Traceability Matrix

| Requirement | Architectural Component | Subtask | Verification Method |
| :--- | :--- | :--- | :--- |
| **REQ-PF-01**: Command routing & aliases | `cli/cmd/`, `cli/constants` | Subtask 210-01 | CLI dispatch resolves `push-fix`, `push fix`, `pf` |
| **REQ-PF-02**: Anti-hang subprocess guards | `cli/cmdpull/push_fix.go` | Subtask 210-01 | Subprocess env inherits `GIT_TERMINAL_PROMPT=0`, `GCM_INTERACTIVE=never` |
| **REQ-PF-03**: Unpushed commits detection | `cli/cmdpull/push_fix.go` | Subtask 210-01 | Calculates commit delta between HEAD and remote tracking branch |
| **REQ-PF-04**: SSH handshake probe | `cli/cmdssh/auth_probe.go` | Subtask 210-02 | Evaluates exit code 1 with GitHub welcome string within 5s |
| **REQ-PF-05**: Dashed bogus keys purge | `cli/store/ssh_keys.go` | Subtask 210-02 | Cleanses `SshKey` table of hyphenated flags |
| **REQ-PF-06**: Dual-profile SSH config write | `cli/cmdssh/sshconfig.go` | Subtask 210-02 | Writes identical block to user home and system profile |
| **REQ-PF-07**: HTTPS -> SSH transport switch | `cli/cmdpull/push_fix.go` | Subtask 210-02 | Rewrites remote URL and updates SQLite `IdentifiedTransport` |
| **REQ-PF-08**: Upstream auto-configuration | `cli/cmdpull/push_fix.go` | Subtask 210-03 | Executes `git push -u` when `@{u}` tracking is missing |
| **REQ-PF-09**: Non-fast-forward auto-rebase | `cli/cmdpull/push_fix.go` | Subtask 210-03 | Runs `git pull --rebase` upon rejected push and retries |
| **REQ-PF-10**: ANSI HUD telemetry card | `cli/cmdpull/push_fix_hud.go`| Subtask 210-04 | Emits formatted terminal summary card with latency benchmarks |

---

## 4. Acceptance Criteria

1. **Non-Interactive Execution Guarantee:** Under no condition does `gitmap push-fix` hang awaiting terminal password or browser OAuth dialogs.
2. **Zero Bogus Key Residue:** Any record in SQLite `SshKey` where `Name` starts with `-` is removed automatically.
3. **Dual Profile Parity:** Both `~/.ssh/config` and `%SystemDrive%\Users\%USERNAME%\.ssh\config` (on Windows) contain the canonical `Host github.com` default configuration block.
4. **Transport Resilience:** Blocked HTTPS remotes seamlessly convert to valid SSH remotes if SSH keys are authenticated.
5. **Divergence Recovery:** Out-of-sync branches automatically rebase remote commits and push without requiring manual terminal rebase commands.
6. **New Branch Push:** Newly created local branches without upstream tracking are pushed using `-u origin <branch>` without error.
7. **Telemetry Visibility:** The command renders a clean, high-contrast HUD card containing commit counts, branch, SHA, remote URL, transport, and latency.
8. **Strict Quality Gates:** Strict adherence to positive boolean naming conventions, zero-build, and zero-test rules.
