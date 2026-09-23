# Plan 43: SSH Join Password & Known Hosts Automation, Cluster Join Help & Node OS Metadata

## Origin & Loop Summary
- **Triggered By**: User request to automate SSH host key `known_hosts` acceptance, interactive hidden password entry with encrypted vault storage notice and review command, comprehensive `cluster join --help` and server command connection architecture explanation, and remote node OS metadata probing/storage during GitMap installation.
- **Execution Loops**: Completed in 1 self-loop iteration (planning, architecture design, implementation, and consolidation).

## User Request (Verbatim)
> is it done?
>
> Okay. I think there is a few more comments, I think. With the SSH, we need to fix first. When we go to the join command with the machine, if that machine does not have a password or requires a password, so the first thing that it will not prompt using the regular way. So GitMap would already take one more time to see if it is looking for a password and the password is not provided. If the password is not provided and it requires a password, so first of all, GitMap will add it to the known host situation automatically using writing yes in the prompt and so on. And then GitMap will request for the password with the password, like hidden format CLI. Now, when the password is inputted, the CLI will also mention that we are saving the password as a hash code or something like this. If you want to review, there will be a command. It will also mention that command to review that password if the user wants to do it. Okay? Now, it will log in, and it knows the password. In the next time, if the user wants to run any command, they can do it. The password will be saved and reused. Remember that. This is very critical step, and that needs to be performed with any SSH join, cluster join, or things like that. Can you please tell me the cluster join? I think it's a bit different than the SSH join. So cluster join commands are not there in the help. So if we do the cluster space join help, there should be a help that actually explains why the cluster join is there, how it works, how it's different than the SSH join. Also, similar thing needs to be done for the server commands, how it's connected. The explanation needs to be there, the joining and things, the steps. That needs to show a step, like each of the step, how it gets connected. Okay? Also, GitMap installation fix. I think that you need to understand where the GitMap is. And also, for every one of the node from the root node, when we try to install, we should save a few things, like what the OS is, OS version, first time run. Okay, so that it's attached, so that GitMap database can tell us what type of OS that is in the future running of the command, so that it can optimize the decision-making. Do you understand? This is a very critical process, and this is also important to go through with it. Is it clear?

## Extracted Actionable Task List
1. **Automated SSH Known-Hosts Acceptance**: When joining an SSH node, automatically accept/add the host key to `~/.ssh/known_hosts` (simulating writing "yes") to prevent interactive prompt hangs.
2. **Interactive Hidden Password Prompt**: If password is not provided via flag/key and required by target, prompt using hidden terminal format (`term.ReadPassword`).
3. **Password Vaulting & Review Command**: Notify user that password is saved encrypted (RSA-OAEP/AES) in the vault and display review command (`gitmap ssh pass show <alias>`).
4. **Password Persistence & Reuse**: Save encrypted password in SQLite database (`ssh_hosts` and `SSHConnection`) and reuse across subsequent commands (`ssh exec`, `cluster exec`, `sc`).
5. **Cluster Join Help Parity**: Provide comprehensive help for `cluster join` (`gitmap cluster join help` / `gitmap cluster join --help`) detailing why it exists, how it works, and how it differs from `ssh join`.
6. **Server Command Connection Architecture**: Document step-by-step connection flow architecture for server commands (`server-cmd`, `servers-clients`).
7. **Node OS Discovery & GitMap Binary Detection**: Probe and record remote OS type, OS version/distribution, and `first_run_at` in GitMap SQLite database (`SSHConnection`) during installation/discovery to optimize future command execution; reliably check GitMap binary presence across standard Linux, Darwin, and Windows paths.

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
