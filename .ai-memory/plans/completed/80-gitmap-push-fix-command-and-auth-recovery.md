# Plan 80: GitMap Push-Fix Command & Authentication Self-Healing Engine (Completed)

> **Plan Status:** Completed  
> **Completed At:** 2026-10-04  
> **Traceability IDs:** Subtask 210-01 .. Subtask 210-05  
> **Spec Reference:** [02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md](../../../02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md)  
> **Target Subsystems:** `cli/cmdpushfix`, `cli/cmdpull`, `cli/cmdssh`, `cli/store`, `cli/constants`, `cli/helptext`  

---

## 1. Architectural Summary & Results

The `gitmap push-fix` autonomous diagnosis, recovery, and push suite has been implemented, tested, and verified end-to-end:
1. **Dual Routing & Aliases**: Supports top-level `gitmap push-fix`, `gitmap pf`, `gitmap pushfix`, `gitmap fix-push` alongside subcommand `gitmap push fix` and `gitmap ph fix`.
2. **4-Phase Recovery Pipeline**:
   - Phase 1: Pre-flight git status inspection (working tree, branch, unpushed commits count).
   - Phase 2: Remote auth probe (bounded 5s SSH handshake with `BatchMode=yes`, anti-hang env `GIT_TERMINAL_PROMPT=0`, `GCM_INTERACTIVE=never`).
   - Phase 3: Self-healing engine (candidate key promotion `id_rsa_-y` -> `id_rsa`, dashed bogus keys cleanup from `gitmap.db`, dual-profile `.ssh/config` synchronization, HTTPS to SSH remote transport conversion).
   - Phase 4: Push dispatch with auto-upstream tracking (`-u`), non-fast-forward rejection detection (`CheckIsNonFastForward`), auto-rebase (`git pull --rebase`), and rich ANSI/ASCII boxed HUD cards.
3. **Bug Fixes**:
   - `gitmap cd` path deduplication: normalized paths so Windows backslashes and forward slashes deduplicate cleanly.
   - `gitmap push` Go runtime stack trace suppression: on failure, displays a clean `[tip] To diagnose and auto-repair SSH authentication & push issues, run: gitmap push-fix` without dumping ugly Go stack traces.

---

## 2. Completed Subtasks & Deliverables

| Subtask ID | Focus Area | Deliverables | Verification |
| :--- | :--- | :--- | :--- |
| **Subtask 210-01** | `cli/cmdpushfix`, `cli/constants` | `pushfix_types.go`, `pushfix_cmd.go`, `pushfix_pipeline.go` | PASS (flags, options, anti-hang env) |
| **Subtask 210-02** | `cli/cmdpushfix`, `cli/cmdssh` | `pushfix_auth.go`, `pushfix_rebase.go` | PASS (SSH probe, key sync, auto-rebase) |
| **Subtask 210-03** | `cli/cmd`, `cli/constants` | `rootcore.go`, `clihelpers.go`, `constants_cli.go` | PASS (routing, subcommand intercept) |
| **Subtask 210-04** | `cli/cmdpushfix`, `cli/helptext` | `pushfix_ui.go`, `push_fix.md`, `print.go`, `catalog.go` | PASS (rich HUD cards, help doc) |
| **Subtask 210-05** | `cli/cmdpushfix`, `cli/constants` | `pushfix_test.go`, `cmd_constants_test.go`, `constants_profile.go` | PASS (10 unit tests, alias uniqueness) |

---

## 3. Acceptance Verification

1. **Non-Interactive Execution Guarantee:** Subprocesses inherit `BuildSafeGitEnv` (`GIT_TERMINAL_PROMPT=0`, `GCM_INTERACTIVE=never`).
2. **Zero Bogus Key Residue:** SQLite `SshKey` table purged of dashed keys (`Name LIKE '-%'`).
3. **Dual Profile Parity:** Both `~/.ssh/config` and `%SystemDrive%\Users\%USERNAME%\.ssh\config` synchronized with canonical `Host github.com`.
4. **Transport Resilience:** HTTPS remotes convert to SSH remotes when SSH keys are authenticated.
5. **Divergence Recovery:** Out-of-sync branches automatically rebase remote commits and push cleanly.
6. **Telemetry & HUD:** High-contrast ANSI/ASCII boxed HUD cards with status badges.
