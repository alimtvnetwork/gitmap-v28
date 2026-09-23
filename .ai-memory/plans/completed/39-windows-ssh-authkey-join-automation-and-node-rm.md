# Plan 39: Windows SSH Authorized Key Architecture, Interactive SSH Join Password Automation, GitMap Host-Node Bootstrapping, and Comprehensive Node Removal

> **Origin:** User request incorporating Windows OpenSSH PowerShell elevation, `$env:ProgramData\ssh\administrators_authorized_keys`, `icacls` ACL inheritance and permissions, `sshd` service health/startup, password prompt auto-detection during `gitmap ssh join`, remote GitMap installation/key exchange, and dual-table node removal (`gitmap ssh rm <ip|alias|--all>`).
> **Execution Loops:** 1 continuous multi-agent orchestration loop with 2 parallel subagents across 4 subtasks.
> **Status:** COMPLETED

---

## Executive Summary & Outcomes

1. **Windows OpenSSH Authorized Keys Specification (`02-spec/23-ssh-fleet-management/01-windows-ssh-authorized-keys-spec.md`)**:
   - Reverse-engineered and formalized the dual authorized keys architecture on Windows.
   - Documented elevation and token checks, localgroup administrators detection, Administrator path (`$env:ProgramData\ssh\administrators_authorized_keys`) with ACL hardening (`icacls $keysFile /inheritance:r /grant "Administrators:F" /grant "SYSTEM:F"`), and standard user profile path (`$profile\.ssh\authorized_keys`).
   - Documented key format validation, whitespace token split deduplication, and automated `sshd` service verification and startup.

2. **Native Go Windows Injection Script (`cli/cmdssh/ssh_auth_key_deploy.go`)**:
   - Upgraded `buildWindowsAuthKeyScript` to generate the complete, self-contained PowerShell one-liner embodying the exact privilege detection, path routing, `icacls` permission hardening, deduplication, and service health activation without external script dependencies.

3. **Dual Table Node Removal Command (`gitmap ssh rm` / `gitmap ssh remove`)**:
   - Added `DeleteAllSSHHosts(ctx, db)` to `cli/store/ssh_repo.go`.
   - Added `DeleteSSHConnectionByTarget(ctx, db, target)` and `DeleteAllSSHConnections(ctx, db)` to `cli/db/sshconnection.go`.
   - Enhanced `cli/cmdssh/sshjoin_rm_cmd.go` with `resolveRmTarget`, `purgeAllSSHRecords`, `purgeTargetSSHRecords`, and `printRmSuccess` to cleanly delete by IP address, alias/ADSC name, or `--all` / `all` across both `ssh_hosts` and `SSHConnection` SQLite tables.
   - Routed `gitmap ssh rm` and `gitmap ssh remove` in `cli/cmdssh/ssh.go` directly to `runSJRm`.

4. **Interactive SSH Join Password Automation & Auto-Accept Host Key (`cli/cmdssh/sshjoin_enroll.go`, `cli/cmdssh/ssh_client.go`)**:
   - Implemented `ExecuteSSHJoinEnrollment` in `cli/cmdssh/sshjoin_enroll.go`.
   - When joining an SSH node without password or pre-authorized key, detects authentication requirement and prompts user interactively using `PromptSSHPassword` if a TTY is available.
   - Automatically configures host key validation with auto-accept on first connection (`NewAutoAcceptHostKeyConfig` using `ssh.InsecureIgnoreHostKey()`).
   - Securely encrypts passwords via `crypto.EncryptSSHPassword` and persists records to both `ssh_hosts` and `SSHConnection` SQLite tables.

5. **Remote GitMap Bootstrapping & Mutual Key Exchange (`cli/cmdssh/sshjoin_enroll.go`)**:
   - During `ssh join` enrollment, verifies whether GitMap is installed on the remote machine (`gitmap --version`).
   - If missing, auto-bootstraps GitMap using platform-specific canonical one-liners (`irm ... | iex` for Windows, `curl ... | bash` for Unix).
   - Discovers host machine's public SSH key (`~/.ssh/id_ed25519.pub`, `id_rsa.pub`, `id_ecdsa.pub`) and deploys it to the node's authorized keys using the upgraded `buildInjectAuthKeyScript`.
   - Confirms bidirectional communication between host GitMap and node GitMap.

---

## Consolidated Subtasks

### Subtask 01: Windows SSH Authorized Keys Spec & Remote Injection Script
- **Target Files:** `02-spec/23-ssh-fleet-management/01-windows-ssh-authorized-keys-spec.md`, `cli/cmdssh/ssh_auth_key_deploy.go`
- **Action:** Authored architectural specification for Windows OpenSSH authorized keys. Upgraded `buildWindowsAuthKeyScript` in `ssh_auth_key_deploy.go` to handle Administrator vs Standard users, `administrators_authorized_keys` location, ACL inheritance reset with `icacls`, key body deduplication, and `sshd` service verification.

### Subtask 02: SSH Node Removal Command & Dual Table Purge
- **Target Files:** `cli/cmdssh/sshjoin_rm_cmd.go`, `cli/store/ssh_repo.go`, `cli/db/sshconnection.go`, `cli/cmdssh/ssh.go`
- **Action:** Implemented `gitmap ssh rm` and `gitmap ssh remove` supporting deletion by IP, alias/ADSC name, or `--all` / `all`. Added `DeleteAllSSHHosts` and dual-table cleanup removing records from both `ssh_hosts` and `SSHConnection` SQLite tables.

### Subtask 03: Interactive SSH Join Password Prompt & Host Key Automation
- **Target Files:** `cli/cmdssh/sshjoin_enroll.go`, `cli/cmdssh/ssh_client.go`, `cli/cmdssh/sshjoin_cmd.go`
- **Action:** Implemented interactive password prompt using `PromptSSHPassword`, auto-accept host key callback, and dual-table encrypted password persistence in `sshjoin_enroll.go`.

### Subtask 04: Remote GitMap Bootstrap & Mutual Key Exchange
- **Target Files:** `cli/cmdssh/sshjoin_enroll.go`, `cli/cmdssh/cluster_install_cmd.go`, `cli/cmdssh/ssh_auth_key_deploy.go`
- **Action:** Added remote GitMap version check, auto-bootstrap installation via one-liners, host public key deployment, and bidirectional communication confirmation.
