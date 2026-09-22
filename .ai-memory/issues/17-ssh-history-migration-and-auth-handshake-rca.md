# Root Cause Analysis (RCA-17): SSH History Migration Missing Column and Join Handshake Protocol Failure

> **Raw Error Log:** [`17-ssh-history-migration-and-auth-handshake-error.txt`](./17-ssh-history-migration-and-auth-handshake-error.txt)

## 1. Symptom
1. **`SQL logic error: no such column: forward_payload (1)` on `gitmap ssh nodes rm`:**
   When removing an enrolled SSH node interactively or via CLI:
   ```text
   PS C:\Users\Alim> gitmap ssh nodes rm
   ...
   Are you sure you want to remove 1 node(s)? [y/N]: y
   [QueryWrapper Error]: exec failed: SQL logic error: no such column: forward_payload (1)
   query: INSERT OR IGNORE INTO TaskHistory
                           (TaskId, Section, Action, Target, ForwardPayload, InversePayload, RestoredAt, CreatedAt)
                           SELECT task_id, 'ssh', action, target, forward_payload, inverse_payload, COALESCE(restored_at, ''), created_at
                           FROM ssh_task_history;
   ✓ 1 machine(s) removed from SSH registry.
   ```
   The error was logged to `os.Stderr` by `store.ExecWrapper` every time `openSSHHistoryDB()` was called, and legacy snapshot entries were never migrated into `TaskHistory`.

2. **`ssh: handshake failed: ssh: unexpected message type 51 (expected 60)` on `gitmap ssh join`:**
   When enrolling a remote Windows or Linux node with password:
   ```text
   PS C:\Users\Alim> gitmap ssh join a@192.168.1.3 w1
   Enter SSH password for a@192.168.1.3:
     ℹ Saving password as encrypted representation (RSA-OAEP/AES) in local vault.
     ℹ Review saved password anytime with: gitmap ssh pass show w1
     ⚠ Failed to authenticate with remote machine: ssh: handshake failed: ssh: unexpected message type 51 (expected 60)
   gitmap ssh: execute failed: [E9000:EXECUTION] ExecuteSSHJoinEnrollment.auth: ssh: handshake failed: ssh: unexpected message type 51 (expected 60) (at=cmdssh/sshjoin_enroll.go:365)
   ```

3. **Premature Password Vault Notification:**
   In `promptAndConnectTarget()`, GitMap announced that the password was saved as an encrypted representation in the vault before authentication was even attempted or verified. When authentication failed, the node was never enrolled and no password was stored, misleading the user.

---

## 2. Root Cause
1. **Unchecked Column Projection in Legacy Table Migration:**
   In `cli/cmdssh/ssh_history_db.go`, `migrateLegacySSHHistoryTables()` executed `SELECT ..., forward_payload, inverse_payload FROM ssh_task_history;`. Legacy `ssh_task_history` tables created prior to commit `c2d192e` only contained `task_id, action, target, payload_json, created_at, restored_at`. Because `forward_payload` did not exist in the table, SQLite threw `SQL logic error: no such column: forward_payload (1)`. Furthermore, the legacy table was never dropped after migration, causing the failed query to repeat on every subsequent command.

2. **Simultaneous Password and Keyboard-Interactive Auth in Go SSH Client:**
   In `cli/cmdssh/sshjoin_enroll.go`, `dialNodeWithPassword()` provided both `ssh.Password(pass)` and `buildKeyboardInteractiveAuth(pass)` in a single `ClientConfig.Auth` slice. The target host (`192.168.1.3`) is running `OpenSSH_for_Windows_9.5`. When password authentication failed, Go's `x/crypto/ssh` client immediately fell back to `keyboard-interactive`. Windows OpenSSH does not implement interactive challenge prompts and responds immediately with `SSH_MSG_USERAUTH_FAILURE` (packet type 51). Because Go's `client_auth.go` expected `SSH_MSG_USERAUTH_INFO_REQUEST` (packet type 60), it aborted the handshake with `ssh: handshake failed: ssh: unexpected message type 51 (expected 60)`, completely masking the underlying credential rejection.

3. **Premature UI Notification Before Transaction Completion:**
   `notifyPasswordEncryptedStorage` was called inside `promptAndConnectTarget()` immediately after receiving the password string from the console prompt, prior to calling `connectWithGivenPass()` and before verifying that the password successfully authenticated.

---

## 3. Resolution
1. **Detect-Then-Act Migration in `cli/cmdssh/ssh_history_migrate.go`:**
   - Implemented `checkColumnExists(conn, tableName, columnName)` using `PRAGMA table_info`.
   - If `forward_payload` exists, columns are migrated directly.
   - If `forward_payload` is missing but `payload_json` exists, mapped `target` to `ForwardPayload` and `payload_json` to `InversePayload`.
   - Dropped legacy `ssh_task_history` and `ssh_task_queue` post-migration so the migration runs strictly once.
   - Decomposed `cli/cmdssh/ssh_history_db.go` into `ssh_history_db.go`, `ssh_history_migrate.go`, and `ssh_history_tasks.go`, keeping each file $\le 100$ lines.

2. **Sequential Password Dialing & Error Normalization in `cli/cmdssh/sshjoin_dial.go`:**
   - Refactored `dialNodeWithPassword` to attempt standard `ssh.Password(pass)` first.
   - If a network error occurs (connection refused, timeout, unreachable), returned immediately without attempting fallback.
   - Only attempted `keyboard-interactive` if `isSSHPasswordUnsupported(errPass)` is true (i.e., the server disabled password auth).
   - If credentials are rejected, prevented fallback to `keyboard-interactive` on servers that advertise it without PAM support.
   - Implemented `normalizeSSHAuthFailure` to translate cryptic protocol errors (`unexpected message type 51 (expected 60)`, `unable to authenticate`) into clear, human-readable errors: `ssh: authentication failed: invalid password or remote server rejected credentials`.
   - Enhanced `tryConnectDefaultKey` to scan all user keys from `findAllUserSSHKeys()`.

3. **Deferred Password Vault Notification:**
   - Removed `notifyPasswordEncryptedStorage` from `promptAndConnectTarget`.
   - Added `notifyPasswordEncryptedStorage` to `finalizeEnrollment`, ensuring the notification only prints after authentication succeeds and enrollment is committed to the database.

---

## 4. Prevention & Learnings
1. **Detect-Then-Act Rule for SQLite Table Migrations:**
   Whenever migrating legacy SQLite tables, NEVER assume column names or schemas remain static across versions. Always inspect schema with `PRAGMA table_info` before executing `INSERT ... SELECT`.
2. **Clean Up Migrated Legacy Tables:**
   One-time migration shims MUST drop or rename the legacy source tables once migration completes to prevent re-executing migrations and polluting error logs on every connection.
3. **Never Mix Password and Keyboard-Interactive in Single Auth Slice:**
   In Go's `crypto/ssh`, never combine `ssh.Password` and `ssh.KeyboardInteractive` in the same `ClientConfig.Auth` slice unless you intend for a failed password to immediately trigger a challenge request. On servers without interactive PAM (such as Windows OpenSSH), this triggers an unexpected packet type 51 protocol failure.
4. **Announce Vault / Storage Mutations Only After Success:**
   Never output success/saving notices prior to completing the authentication or write phase.
