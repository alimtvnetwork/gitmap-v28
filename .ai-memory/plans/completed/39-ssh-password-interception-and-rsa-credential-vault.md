# Milestone 39: SSH Password Interception, RSA Credential Vault, and Push Self-Healing

- **Slug:** `ssh-password-interception-and-rsa-credential-vault`
- **Milestone Index:** `39`
- **Status:** `COMPLETED`
- **Source Plans Merged:** Plans 67, 68, 78, 80, 203, 204, 208, 210
- **Folded Subtask Folders:**
  - `01-ports-and-ssh-enablement` (3 files)
  - `207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca` (7 files)
  - `208-detectedproject-foreign-key-and-ssh-lifecycle-hardening` (7 files)
- **Target Subsystems:** `cli/cmdssh/`, `cli/vault/`, `cli/store/`, `cli/cmdpushfix/`

---

## 1. Domain Context & Architectural Problem

Connecting and authenticating across heterogeneous multi-node environments presented significant security and reliability hurdles:
1. **Plaintext Password Exposure:** Early SSH commands allowed password parameters in cleartext shell arguments or stored credentials unencrypted in SQLite databases, violating basic security standards and risking exposure in shell histories.
2. **Missing Consent & Key Lifecycle Protocols:** Automated SSH key distribution lacked explicit interactive consent workflows, meaning automated scripts could alter remote authorized keys without administrator awareness.
3. **Database Foreign Key Violations during DB Reset:** Resetting GitMap databases (`gitmap db reset`) triggered foreign key constraint failures across `DetectedProject` and `SSHNode` tables due to non-cascading deletion sequences.
4. **Git Remote Push Rejections:** Push routines frequently failed due to stale credentials, expired tokens, or desynchronized tracking branches without an automated self-healing mechanism.

---

## 2. Synthesized Architectural Outcomes

### 2.1 Masked Password Capture & User Consent Protocol
- **Interactive Masked Terminal Input:** Implemented in `cli/cmdssh/prompt_auth.go` using raw terminal mode (`golang.org/x/term`) to securely capture user passwords with bullet masking (`*`) and zero echo back to standard output or log streams.
- **Explicit Consent Protocol:** Prior to enrolling an SSH node or deploying authorized keys, GitMap prompts for explicit user consent, displaying target IP, SSH port, user account, and cryptographic fingerprint.
- **Non-Interactive Bypass:** For headless CI/CD execution, credentials may be provided via environment variables (`GITMAP_SSH_PASSWORD`, `GITMAP_VAULT_KEY`) or encrypted vault records.

### 2.2 RSA-OAEP Salt Credential Encryption Vault
- **Encrypted Local Vault (`credentials.vault`):** Designed and deployed a secure local credential storage subsystem in `cli/vault/` leveraging RSA-2048 with OAEP padding and PBKDF2 salt derivation (SHA-256 with 100,000 iterations).
- **Zero Plaintext at Rest:** Node passwords, SSH private keys, and API access tokens are encrypted with machine-specific entropy before being persisted to disk (`~/.gitmap/credentials.vault`).
- **Memory Zeroization:** Sensitive buffers in memory are explicitly wiped (`ZeroMemory`) following encryption and decryption routines.

### 2.3 SSH Daemon & Cross-OS Firewall Hardening
- **Windows OpenSSH Server Management:** Added `gitmap ssh service enable` and `gitmap ports inspect` to detect Windows Update (`wuauserv`) dependencies, configure OpenSSH Server Windows service (`sshd`), and create inbound Windows Defender Firewall rules for port 22.
- **Linux UFW & Systemd Hardening:** Automated SSH daemon verification on Linux systems (`systemctl enable --now ssh`, `ufw allow 22/tcp`), inspecting listen sockets and verifying key-based authentication settings in `sshd_config`.

### 2.4 DetectedProject Foreign Key Lifecycle Hardening & Safe DB Reset
- **Cascading Foreign Key Management:** Restructured database schemas in `cli/store/schema.go` with explicit `ON DELETE CASCADE` constraints across `DetectedProjects`, `ProjectBookmarks`, and `SSHNodes`.
- **Safe DB Reset Workflow (`gitmap db reset`):** Implemented pre-flight dependency analysis that unlinks dependent entities, truncates tables in strict topological order, vacuum-cleans SQLite storage, and re-initializes baseline configuration safely.

### 2.5 Git Push Self-Healing Engine (`gitmap push-fix`)
- **Automated Push Recovery:** Implemented `cli/cmdpushfix/push_fix.go` to diagnose rejected push operations.
- **Self-Healing Strategy:**
  1. Detects rejection reason (rejected non-fast-forward, remote branch diverged, missing upstream, or authentication expired).
  2. Executes safe fetch and fast-forward rebase (`git fetch && git rebase origin/<branch>`).
  3. Prompts for updated credentials or re-authenticates via credential vault if auth failure is detected.
  4. Retries push with verified validation checks.

---

## 3. Go Type Contracts & Architecture

```go
// VaultCredential represents an encrypted credential entity
type VaultCredential struct {
    CredentialID string    `json:"credentialId"`
    NodeHost     string    `json:"nodeHost"`
    Username     string    `json:"username"`
    EncryptedKey []byte    `json:"encryptedKey"`
    Salt         []byte    `json:"salt"`
    CreatedAt    time.Time `json:"createdAt"`
    LastUsedAt   time.Time `json:"lastUsedAt"`
}

// SSHServiceStatus models daemon and firewall state
type SSHServiceStatus struct {
    IsDaemonInstalled bool   `json:"isDaemonInstalled"`
    IsDaemonRunning   bool   `json:"isDaemonRunning"`
    Port              int    `json:"port"`
    IsFirewallAllowed bool   `json:"isFirewallAllowed"`
    ServiceError      string `json:"serviceError,omitempty"`
}

// PushFixResult captures push self-healing outcomes
type PushFixResult struct {
    IsRecovered   bool   `json:"isRecovered"`
    OriginalError string `json:"originalError"`
    FixActionTaken string `json:"fixActionTaken"` // "REBASE_APPLIED", "UPSTREAM_SET", "AUTH_REFRESHED"
    CommitCount   int    `json:"commitCount"`
}
```

Cryptographic errors and database foreign key violations map cleanly to structured error types with actionable remediation instructions.

---

## 4. Subtask Verification Ledger

| Folded Subtask Directory | Source Subtask Files | Verified Criteria & Delivered Artifacts |
| :--- | :--- | :--- |
| `01-ports-and-ssh-enablement` | 3 subtask files | Port inspection CLI, OpenSSH Windows service installer, firewall rule generator. |
| `207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca` | 7 subtask files | Database reset safety checks, SSH key lifecycle management, foreign key dependency solver. |
| `208-detectedproject-foreign-key-and-ssh-lifecycle-hardening` | 7 subtask files | Foreign key constraints schema update, cascading deletion triggers, SSH connection test harness. |

---

## 5. Quality & Coding Guideline Compliance

- **Positive Booleans Only:** Enforced `isDaemonInstalled`, `isDaemonRunning`, `isFirewallAllowed`, `isRecovered`.
- **Zero Plaintext Secrets:** Passwords never logged, exported to JSON, or exposed in error messages.
- **Relatively Anchored DB:** Vault and database files resolve relative to user home directory (`os.UserHomeDir`) with secure permissions (`0600`).
- **Comprehensive Linters Passed:** Zero violations across boolean naming, relative paths, and error wrapping rules.
