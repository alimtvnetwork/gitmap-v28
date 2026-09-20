# Plan 43: SSH Join Password & Known Hosts Automation, Cluster Join Help & Node OS Metadata

## Origin & Loop Summary
- **Triggered By**: User request to automate SSH host key `known_hosts` acceptance, interactive hidden password entry with encrypted vault storage notice and review command, comprehensive `cluster join --help` and server command connection architecture explanation, and remote node OS metadata probing/storage during GitMap installation.
- **Execution Loops**: Completed in 1 self-loop iteration (planning, architecture design, implementation, and consolidation).

---

## Executive Summary

1. **SSH Known-Hosts Auto-Acceptance**:
   - In `cli/cmdssh/sshjoin_enroll.go`, integrated `autoTrustTargetHost` into `resolveTargetClient`.
   - Before attempting SSH dial, the target host's public key is captured via `TrustRemoteTarget` and automatically appended to `~/.ssh/known_hosts` and the local SQLite security ledger.
   - Eliminates terminal hangs caused by OpenSSH's `Are you sure you want to continue connecting (yes/no)?` prompt.

2. **Hidden Interactive Password Prompt & Encrypted Vault Storage**:
   - When key authentication fails and no password was supplied, GitMap prompts for the password using hidden terminal input (`PromptSSHPassword` / `term.ReadPassword`).
   - Displays clear notification to the user:
     ```text
     ℹ Saving password as encrypted representation (RSA-OAEP/AES) in local vault.
     ℹ Review saved password anytime with: gitmap ssh pass show <alias>
     ```
   - Persists the encrypted password into `ssh_hosts` and `SSHConnection` tables so subsequent commands (`gitmap ssh exec`, `gitmap cluster exec`, `gitmap sc`) reuse it seamlessly.
   - Implemented `gitmap ssh pass show <alias|ip>` and `gitmap ssh pass ls` in `cli/cmdssh/ssh_pass_cmd.go` to safely review and decrypt saved credentials using local SSH RSA private keys.

3. **Node OS Metadata Discovery & GitMap Installation**:
   - Enhanced `cli/cmdssh/ssh_exec_install.go`:
     - Updated `getGitmapCheckCmd` with a robust multi-path check that finds `gitmap` across `$HOME/.local/bin`, `/usr/local/bin`, `/usr/bin`, `$HOME/go/bin`, `/snap/bin`, and on Windows in `%LOCALAPPDATA%\gitmap\bin` and PATH.
     - Fixed `wrapCommandForShell` in `cli/crypto/ssh_client.go` to properly quote and wrap commands in `bash -c` and `sh -c`.
   - Extended `cli/db/sshconnection.go`:
     - Added `OSVersion string` and `FirstRunAt time.Time` to `SSHConnection`.
     - Implemented automatic schema migration (`ALTER TABLE SSHConnection ADD COLUMN ...`).
     - Added OS version probing (`probeRemoteOSVersion`) in `cli/cmdssh/sshjoin_enroll.go` to discover distribution and version (e.g. `Ubuntu 22.04.4 LTS`, `Microsoft Windows 10 Pro`).
     - Stored `FirstRunAt` timestamp upon initial enrollment.

4. **Cluster Join Help & Server Connection Architecture Documentation**:
   - Expanded `cli/helptext/cluster-join.md`:
     - Documented why `cluster join` exists and how it works.
     - Triad Comparison table: `ssh join` (ad-hoc SSH access) vs `cluster join` (role-based cluster admission, topology orchestration, Kubernetes lifecycle, recipes).
     - Documented the full 7-step connection and admission protocol.
   - Expanded `cli/helptext/server-cmd.md`:
     - Documented the 6-step server command connection architecture and execution flow.
   - Updated `cli/helptext/servers-clients.md` and `cli/helptext/ssh-join.md`:
     - Added password review examples and step-by-step admission details.

---

## Deliverables & Modified Files

1. `cli/cmdssh/ssh_pass_cmd.go` & `cli/cmdssh/ssh_pass_cmd_test.go`:
   - Password review command (`gitmap ssh pass show`, `gitmap ssh pass ls`).
2. `cli/cmdssh/ssh.go`:
   - Routed `pass` and `password` subcommands in `dispatchCoreSSH`.
3. `cli/cmdssh/sshjoin_cmd.go`:
   - Routed `pass` and `password` in `dispatchSJActionSubcommand`.
4. `cli/cmdssh/sshjoin_enroll.go`:
   - Integrated `autoTrustTargetHost`, `notifyPasswordEncryptedStorage`, `probeRemoteOSVersion`, and passed `osVersion` to `persistDualTables`.
5. `cli/cmdssh/ssh_exec_install.go`:
   - Robust `getGitmapCheckCmd` covering all standard paths for Linux/Unix and Windows.
6. `cli/crypto/ssh_client.go`:
   - Wrapped `bash` and `sh` commands in `wrapCommandForShell`.
7. `cli/db/sshconnection.go` & `cli/db/sshconnection_test.go`:
   - Added `OSVersion` and `FirstRunAt` to `SSHConnection`, with auto-migration and query scanners.
8. `cli/helptext/cluster-join.md`:
   - Detailed cluster join help, triad comparison, and 7-step connection flow.
9. `cli/helptext/server-cmd.md`, `cli/helptext/servers-clients.md`, `cli/helptext/ssh-join.md`:
   - Server connection flow and password review documentation.

---

## Verification & Guidelines Compliance
- Functions <= 15 lines max: Verified across all modified and newly created files.
- Affirmative booleans only: No `!` negation operators.
- Universal `*apperror.AppError` wrappers: Verified.
- Blank line before every return: Verified.
- Unix LF line endings: Enforced.
- Zero test runner or routine build executions during execution turns.
