# Subtask 06: AGM Fleet Deployment, Auto-Removal & Transit Security

> **Subtask ID:** Subtask-06  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/cmdnodes/`, `cli/cmdagm/`  
> **Owned Files:**  
> - `cli/cmdnodes/nodes_deploy_agm.go`  
> - `cli/cmdnodes/nodes_deploy_transit.go`  
> - `cli/cmdnodes/nodes_shred.go`  
> - `cli/cmdagm/agm_cleanup.go`  
> - `cli/cmdagm/agm_routes.go`  
> - `cli/cmdnodes/nodes_deploy_agm_test.go`  

---

## 1. Concrete Objectives

1. **Zero-Disk In-Memory Streaming:**
   - Eliminate legacy staging files `C:\Windows\Temp\agm_accounts_sync.tar.gz` and `/tmp/agm_accounts_sync.tar.gz` in `cli/cmdnodes/nodes_deploy_agm.go`.
   - Implement `StreamZeroDiskAGM` piping in-memory tarball directly to remote process standard input:
     - On Windows target nodes: `powershell.exe -NoProfile -Command "$dest = Join-Path $env:USERPROFILE '.antigravity_tools'; if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Force -Path $dest | Out-Null }; tar.exe -xzf - -C $dest"`
     - On Linux/macOS target nodes: `sh -c "mkdir -p ~/.antigravity_tools && tar -xzf - -C ~/.antigravity_tools && chmod 700 ~/.antigravity_tools && chmod 600 ~/.antigravity_tools/accounts.json ~/.antigravity_tools/user_tokens.db 2>/dev/null || true"`
2. **AES-256-GCM & RSA Transit Cryptography:**
   - Implement `EncryptTransitEnvelope` in `cli/cmdnodes/nodes_deploy_transit.go`.
   - Symmetrically encrypt packaged credentials using an ephemeral 256-bit AES-GCM session key and unique 12-byte nonce.
   - Encrypt the AES session key using the destination node's RSA public key (stored in Split-DB vault).
   - Bundle into an authenticated transit envelope with SHA-256 integrity verification and a 5-minute replay prevention timestamp.
3. **Dual REST API Deployment Endpoints:**
   - In `cli/cmdagm/agm_routes.go`, register `POST /api/v1/agm/deploy` on the GitMap local daemon.
   - Validate bearer JWT authorization tokens before unpacking and extracting in-memory payloads.
4. **Temporary Backup JSON Auto-Removal (`D:\agm_accounts_backup_*.json`):**
   - Implement `SecureShredFile` in `cli/cmdnodes/nodes_shred.go` executing a 3-pass sanitization overwrite (Pass 1: `0x55`, Pass 2: `0xAA`, Pass 3: cryptographically secure random bytes) before unlinking.
   - Implement `DiscoverTemporaryBackupFiles()` scanning root drive patterns (`D:\agm_accounts_backup_*.json`, `%USERPROFILE%\agm_accounts_backup_*.json`, `%TEMP%\agm_*.json`).
   - Hook automatic cleanup into `executeFleetAGMDeployment`: upon 100% successful node deployment, trigger secure shredding of discovered temporary backup JSON files.
5. **Split-DB Audit Persistence:**
   - Persist deployment outcomes to `AgmDeploymentAudit` and shredded files to `AgmBackupCleanupAudit` in `.gitmap/data/installation/agm_fleet.db`.
6. **Dedicated CLI Command:**
   - Expose `gitmap agm cleanup-backups [--force] [--dry-run]` and support `--cleanup-backups` on `gitmap nodes deploy agm-accounts`.

---

## 2. Core Domain Types & Structs

```go
package cmdnodes

import (
	"time"
)

// AgmTransitEnvelope encapsulates an encrypted credential deployment payload.
type AgmTransitEnvelope struct {
	Version          string    `json:"version"`           // "v1-aes256gcm"
	TargetNodeAlias  string    `json:"target_node_alias"`
	WrappedKeyHex    string    `json:"wrapped_key_hex"`   // RSA-OAEP encrypted AES key
	NonceHex         string    `json:"nonce_hex"`         // 12-byte GCM nonce
	CiphertextBase64 string    `json:"ciphertext_base64"` // AES-256-GCM encrypted tarball
	TotalAccounts    int       `json:"total_accounts"`
	PayloadSha256    string    `json:"payload_sha256"`
	ExpiresAt        time.Time `json:"expires_at"`
}

// AgmBackupCleanupRecord documents shredded temporary backup files.
type AgmBackupCleanupRecord struct {
	Id            string    `json:"id"`
	FilePath      string    `json:"file_path"`
	FileSizeBytes int64     `json:"file_size_bytes"`
	PassCount     int       `json:"pass_count"`
	ShreddedAt    time.Time `json:"shredded_at"`
	Status        string    `json:"status"` // "shredded", "failed"
	ErrorMessage  string    `json:"error_message,omitempty"`
}

// AgmDeploymentAuditRecord records cluster deployment telemetry in Split-DB.
type AgmDeploymentAuditRecord struct {
	Id               string    `json:"id"`
	Initiator        string    `json:"initiator"`
	TotalTargets     int       `json:"total_targets"`
	SuccessTargets   int       `json:"success_targets"`
	AccountsCount    int       `json:"accounts_count"`
	ZeroDiskVerified bool      `json:"zero_disk_verified"`
	BackupsCleaned   int       `json:"backups_cleaned"`
	DeployedAt       time.Time `json:"deployed_at"`
	DurationMs       int64     `json:"duration_ms"`
}
```

---

## 3. Implementation Checklist

- [ ] **1. Multi-Pass Secure Shredding Engine (`cli/cmdnodes/nodes_shred.go`):**
  - Implement `SecureShredFile(filePath string) error`.
  - Validate file existence and file size via `os.Stat`.
  - Open file with `os.O_WRONLY`.
  - Write Pass 1 buffer filled with `0x55`, sync to disk via `file.Sync()`.
  - Write Pass 2 buffer filled with `0xAA`, sync to disk via `file.Sync()`.
  - Write Pass 3 buffer with `crypto/rand.Read` random bytes, sync to disk.
  - Zero out the file buffer and close the handle.
  - Remove file via `os.Remove(filePath)`.
- [ ] **2. Backup Discovery Scanner (`cli/cmdnodes/nodes_shred.go`):**
  - Implement `DiscoverTemporaryBackupFiles() ([]string, error)`.
  - Check `D:\` drive for `agm_accounts_backup_*.json` matches.
  - Check user profile directory and OS temp folder for stale backup patterns.
  - Return unique deduplicated list of discovered paths.
- [ ] **3. Transit Encryption Pipeline (`cli/cmdnodes/nodes_deploy_transit.go`):**
  - Implement `EncryptTransitEnvelope(nodeConn db.SSHConnection, tarData []byte, count int) (*AgmTransitEnvelope, error)`.
  - Generate 32-byte AES key and 12-byte nonce using `crypto/rand`.
  - Encrypt tar bytes using `cipher.NewGCM`.
  - Encrypt AES key using node's RSA public key with OAEP (SHA-256).
  - Compute SHA-256 digest of plaintext payload.
  - Assemble `AgmTransitEnvelope` with 5-minute TTL.
- [ ] **4. Zero-Disk Remote Stream Extraction (`cli/cmdnodes/nodes_deploy_agm.go`):**
  - Refactor `streamAndExtractAGM` to use `StreamZeroDiskAGM`.
  - Open SSH session, set `stdinPipe`, and pipe decrypted archive bytes directly into stdin of remote `tar -xzf -`.
  - Assert that no staging files are created in `C:\Windows\Temp\` or `/tmp/`.
- [ ] **5. Deployment Hook & Split-DB Audit (`cli/cmdnodes/nodes_deploy_agm.go`):**
  - After completing node deployments, check if `opts.CleanupBackups` is true (default).
  - If all target nodes succeeded, invoke `DiscoverTemporaryBackupFiles()` and shred matching files.
  - Insert records into `AgmDeploymentAudit` and `AgmBackupCleanupAudit` in `.gitmap/data/installation/agm_fleet.db`.
- [ ] **6. REST API Endpoints (`cli/cmdagm/agm_routes.go`):**
  - Implement `POST /api/v1/agm/deploy` handling incoming encrypted envelopes.
  - Verify JWT bearer signature and replay expiration.
  - Extract payload in memory directly to `~/.antigravity_tools`.
- [ ] **7. Comprehensive Testing (`cli/cmdnodes/nodes_deploy_agm_test.go`):**
  - Unit test `SecureShredFile`: verify file is overwritten with random bytes and removed.
  - Test transit encryption and decryption loop.
  - Test discovery of mock backup files matching `agm_accounts_backup_*.json`.

---

## 4. Acceptance Criteria

- [x] Running `gitmap nodes deploy agm-accounts` leaves zero intermediate staging files in OS temp storage.
- [x] All credential data transmitted across network streams is encrypted with AES-256-GCM.
- [x] Upon successful fleet deployment, any temporary backup JSON files matching `agm_accounts_backup_*.json` are securely shredded with 3 overwrite passes and deleted.
- [x] If a deployment fails or is executed with `--dry-run`, temporary backup files are strictly preserved.
- [x] Every deployment and cleanup event creates an immutable audit row in Split-DB SQLite (`AgmDeploymentAudit` and `AgmBackupCleanupAudit`).

---

## 5. Verification Commands

```bash
# Run unit tests
go test -v ./cli/cmdnodes/...

# Verify build
go build -v ./cli/...

# Test backup cleanup simulation
gitmap agm cleanup-backups --dry-run

# Test fleet deployment with zero-disk enforcement
gitmap nodes deploy agm-accounts --dry-run --json
```
