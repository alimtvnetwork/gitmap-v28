# RCA-52: SSH Join Re-Enrollment Authentication Failure, Premature Interactive Prompt, and Windows Double Shell Wrapping

## 1. Symptom

When enrolling or re-enrolling an SSH node via `gitmap ssh join a@192.168.1.20 main`:
1. On initial join, the user is prompted for credentials and the machine joins successfully.
2. However, upon re-running the join command (`gitmap ssh join a@192.168.1.20 main`), GitMap repeatedly prompts for the password again rather than reusing stored vault credentials or falling back to configuration files (`vmpass.json` / `06-vmpass.json`).
3. Furthermore, public key authentication (`○ [3] Public Key Auth: checked 1 default user keys, none accepted`) fails to authenticate on subsequent runs because the initial public key deployment via PowerShell failed silently on the remote Windows host.
4. If the user presses Enter or enters a mismatched password on the subsequent join attempt, the command aborts with:
   ```
   ▲ Failed to authenticate with remote machine: ssh: authentication failed: invalid password or remote server rejected credentials: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain
   ```

## 2. Root Cause

1. **Premature Interactive Password Prompting:** `ExecuteSSHJoinEnrollment` called `promptPasswordIfInteractive(opts)` prior to testing default user SSH keys (`~/.ssh/id_ed25519`, `~/.ssh/id_rsa`), local vault credentials (`gitmap.db`), or configuration fallbacks (`vmpass.json`). This forced unnecessary interactive prompts even when valid credentials or keys were already available.
2. **Missing Local Vault & Fallback Credential Resolution in SSH Join:** `detectAndConnectAuth` did not query `lookupVaultPassword` or `ResolveFallbackCredentials` when public key authentication did not immediately succeed, forcing manual user re-entry on every join command.
3. **PowerShell Double-Wrapping Syntax Error:** `buildWindowsAuthKeyScript` wrapped the key injection script in `powershell -NoProfile -Command "..."`, while `wrapCommandForShell` in `cli/crypto/ssh_client.go` unconditionally wrapped commands with `shellType == "ps"` in a secondary `powershell -NoProfile -Command "..."`. This nested quote invocation caused a PowerShell command parsing error (`\ : The term '\' is not recognized as the name of a cmdlet...`), preventing the public key from being written to `authorized_keys` during the initial enrollment.
4. **Single-Target Scope on Windows:** Windows OpenSSH server configurations vary between checking `__PROGRAMDATA__/ssh/administrators_authorized_keys` and the user's `$HOME\.ssh\authorized_keys`. Writing exclusively to one location caused authentication misses when non-admin users connected or when group matching was disabled.

## 3. Resolution

1. **Eliminated Premature Prompt:** Removed the eager `promptPasswordIfInteractive` call in `ExecuteSSHJoinEnrollment`, ensuring authentication attempts follow the secure precedence:
   - Primary: Installed user SSH keys (`~/.ssh/id_ed25519`, `~/.ssh/id_rsa`, etc.)
   - Secondary: Local encrypted vault credentials for target alias or IP (`gitmap.db`)
   - Tertiary: Fallback cluster credentials from `06-vmpass.json` / `vmpass.json` (`ResolveFallbackCredentials`)
   - Interactive: User prompt only when all non-interactive methods are exhausted.
2. **Auto-Recovery in `resolveTargetClient`:** When explicit or prompted passwords are provided, `resolveTargetClient` and `detectAndConnectAuth` automatically attempt fallback to stored vault credentials and configuration defaults if the primary attempt encounters an authentication rejection.
3. **Prevented Double Shell Wrapping:** Updated `wrapCommandForShell` in `cli/crypto/ssh_client.go` to inspect the command string and avoid double-wrapping commands that already begin with `powershell`, `pwsh`, `cmd.exe`, or `bash`.
4. **Dual Key Deployment for Windows OpenSSH:** Updated `winAuthKeyDeployScript` in `cli/cmdssh/ssh_auth_key_deploy.go` to deploy public keys to both `$env:USERPROFILE\.ssh\authorized_keys` (with proper user ACLs) and `$env:ProgramData\ssh\administrators_authorized_keys` (with Administrator/SYSTEM ACLs), ensuring reliable public key acceptance across all Windows OpenSSH configurations.

## 4. Prevention & Learnings

- **Precedence Rule:** Network commands that can authenticate non-interactively must ALWAYS verify keys, vault, and fallback stores before prompting for manual user input.
- **Shell Wrapper Hygiene:** Shell wrapping utilities must be idempotent; if a command already contains an explicit shell interpreter invocation, subsequent wrappers must never nest it.
- **Cross-Platform Host Key Verification:** Unit tests for remote script generation must assert that generated command strings contain only valid unnested quoting.
