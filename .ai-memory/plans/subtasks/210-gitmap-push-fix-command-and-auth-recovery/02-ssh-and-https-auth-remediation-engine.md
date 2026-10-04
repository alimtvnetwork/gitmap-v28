# Subtask 210.2: SSH & HTTPS Authentication Remediation Engine

- **Parent Plan:** [80-gitmap-push-fix-command-and-auth-recovery.md](../../80-gitmap-push-fix-command-and-auth-recovery.md)
- **Spec Reference:** [02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md](../../../../02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmdssh`, `cli/store`, `cli/cmdpull`

---

## Objective
Implement autonomous authentication probing and remediation for remote pushes: execute timeout-bounded SSH handshakes, purge dashed bogus records from GitMap's SQLite `SshKey` table, synchronize managed SSH configurations across dual Windows profiles, and convert blocked HTTPS repositories to authenticated SSH remotes using `SetRepoIdentifiedTransport`.

---

## Detailed Implementation Breakdown

### 1. SSH Handshake Probe (`cli/cmdssh/auth_probe.go`)
- Create `cli/cmdssh/auth_probe.go` providing timeout-bounded remote connectivity probes.
- Implement `ProbeSSHHandshake(host string, timeout time.Duration) (bool, string, error)`:
  - Executes: `ssh -T -o BatchMode=yes -o ConnectTimeout=<seconds> -o StrictHostKeyChecking=accept-new git@<host>`.
  - Default timeout: 5 seconds.
  - Inspect output and exit code:
    - GitHub returns exit status `1` with message `Hi <user>! You've successfully authenticated...` on successful authentication.
    - Parse authenticated username and return `isAuthSuccess = true`.
    - If output contains `Permission denied (publickey)`, return `isAuthSuccess = false`.

### 2. HTTPS Remote Probing (`cli/cmdssh/auth_probe.go`)
- Implement `ProbeHTTPSAuth(remoteURL string, timeout time.Duration) (bool, error)`:
  - Executes: `git ls-remote --exit-code <remoteURL> HEAD`.
  - Injects anti-hang environment variables (`GIT_TERMINAL_PROMPT=0`, `GCM_INTERACTIVE=never`).
  - If exit code is `0`, returns `isAuthSuccess = true`.
  - If exit code is `128` (authentication failure / terminal prompt rejected), returns `isAuthSuccess = false`.

### 3. Bogus Key Purging in SQLite `SshKey` Table (`cli/store/ssh_keys.go`)
- Add method `(db *DB) PurgeBogusDashedSSHKeys() (int64, error)`:
  - Executes SQL query:
    ```sql
    DELETE FROM SshKey WHERE Name LIKE '-%' OR TRIM(Name) = ''
    ```
  - Eliminates corrupted records where CLI flags (e.g., `-y`, `-f`, `--force`) were inadvertently inserted as key names.
  - Returns count of purged rows and logs sanitization notice.

### 4. Dual-Profile SSH Config Synchronization (`cli/cmdssh/sshconfig.go`)
- Refactor `updateSSHConfig(db *store.DB)`:
  - Canonical default key emission: ensure `Host github.com` stanza is always emitted for the default key (`constants.DefaultSSHKeyName` or first available key).
  - Dual-profile writing on Windows:
    1. Resolve user home config path: `filepath.Join(homeDir, ".ssh", "config")`.
    2. Resolve system profile config path: `filepath.Join(os.Getenv("SystemDrive"), "Users", os.Getenv("USERNAME"), ".ssh", "config")`.
    3. Ensure parent directory `.ssh` exists with secure permissions (`0700`).
    4. Write sanitized managed configuration block to both paths if distinct.

### 5. Remote Transport Conversion (`cli/cmdpull/push_fix.go`)
- When the remote URL uses HTTPS and HTTPS probing fails (e.g. expired credentials, OAuth prompt required):
  1. Test SSH handshake via `ProbeSSHHandshake("github.com", 5*time.Second)`.
  2. If SSH handshake succeeds, convert repository remote:
     - Derive SSH URL from HTTPS URL (e.g., `https://github.com/foo/bar.git` -> `git@github.com:foo/bar.git`).
     - Update local Git configuration: `git remote set-url origin git@github.com:foo/bar.git`.
  3. Persist the updated transport in SQLite:
     ```go
     db.SetRepoIdentifiedTransport(remoteURL, "ssh")
     ```
  4. Output notification:
     ```
     ↻ Transport converted from HTTPS to SSH (valid SSH key confirmed).
     ```
