# Specification 343: Improvements 661 to 700 - Vault, Header Trace, Sharding & AGY Commander

- **Sequence:** 343
- **Slug:** `343-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/04-improvements-661-to-700-sharding-audit-and-agy-fleet.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/04-improvements-561-to-600-security-audit-and-swarm-telemetry.md`
  - `02-spec/21-app/06-database-and-split-db/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification defines Improvements 661 through 700 for Sequence 343 of the GitMap Warm-Up and Emailing UI Engine. Sequence 343 expands enterprise deliverability, zero-trust security, distributed data sharding, and multi-agent coordination across four mission-critical clusters:

1. **Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD (Items 661–670):** Hardware-backed envelope encryption (AES-256-GCM / ChaCha20-Poly1305), automated key rotation, zero-swap in-memory secret caches, Casbin RBAC tenant isolation, cryptographic audit ledgers, and terminal enclave status HUDs.
2. **Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer (Items 671–680):** High-speed parsing of RFC 5322 and RFC 5321 `Received:` transit headers, multi-hop latency waterfall timelines, offline ASN and GeoIP routing enrichment, intermediate TLS cipher suite inspection, and gateway tampering detection.
3. **Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver (Items 681–690):** Physical shard routing via consistent hashing, streaming WAL frame replication, cryptographic Merkle-tree page diffing, two-phase distributed epoch commits, online non-blocking vacuum compaction, and split-db multi-cluster management cockpits.
4. **Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander (Items 691–700):** Real-time multi-agent swarm mission cockpit, peer discovery gossip mesh, subtask leasing arbitration to eliminate concurrency conflicts, token and LLM expenditure guards, circular deadlock remediation, and SQLite FTS5 conversation transcript indexing.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│               Sequence 343 Improvements 661 to 700 Cluster Distribution               │
├────────────────────────────┬────────────────────────────┬──────────────────────────────┤
│ Cluster 67 (661-670)       │ Cluster 68 (671-680)       │ Cluster 69 (681-690)         │
│ Multi-Tenant Zero-Trust    │ Forensic Header Trace      │ Distributed SQLite Sharded   │
│ Mailbox Credential Vault   │ & Cross-Gateway Visualizer │ Replicator & Merkle Resolver │
├────────────────────────────┴────────────────────────────┼──────────────────────────────┤
│ Cluster 70 (691-700)                                    │ Cross-Cutting Standards      │
│ Autonomous Antigravity Swarm Orchestration Engine       │ Positive Booleans (is/has),  │
│ & Real-Time Mission Commander                           │ Split-DB, Zero CDN, No Exits │
└─────────────────────────────────────────────────────────┴──────────────────────────────┘
```

---

## 2. Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD (Items 661–670)

### Improvement 661: Hardware-Backed Envelope Encryption & Key Derivation Service
- **Title:** Hardware-Backed Envelope Encryption & Key Derivation Service
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Safeguard mailbox SMTP/IMAP passwords and OAuth tokens using AES-256-GCM / ChaCha20-Poly1305 envelope encryption backed by TPM 2.0, Apple Secure Enclave, or PKCS#11 HSM providers.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Security enclave status panel displaying hardware provider (`TPM 2.0`, `PKCS#11`, or `Software fallback`), master key fingerprint, derivation iterations (`Argon2id` parameters), and cryptographic health badges in Cyan and Emerald tokens.
  - *Web UI:* Zero-CDN SVG hardware lock indicator with visual key hierarchy breakdown (Master Key -> Key Encryption Key -> Data Encryption Key).
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault init [--provider=tpm|pkcs11|software] [--key-id=<id>] [--iterations=4] [--json]`
  - Flags: `--provider` (string, hardware enclave provider), `--key-id` (string, master key identifier), `--iterations` (int, Argon2id time cost), `--json` (bool, JSON formatted output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/enclave_types.go`):
    ```go
    type EnclaveKeyDerivation struct {
        KeyID              string `json:"keyId"`
        ProviderType       string `json:"providerType"`
        DerivationAlgorithm string `json:"derivationAlgorithm"`
        MemoryCostKiB      int    `json:"memoryCostKiB"`
        TimeCostIterations int    `json:"timeCostIterations"`
        KeyLengthBytes     int    `json:"keyLengthBytes"`
        FingerprintSHA256  string `json:"fingerprintSha256"`
        CreatedAt          int64  `json:"createdAt"`
        isHardwareBacked   bool   `json:"isHardwareBacked"`
        hasActiveSession   bool   `json:"hasActiveSession"`
        isDerivationValid  bool   `json:"isDerivationValid"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationVaultEnclaves (
        KeyID TEXT PRIMARY KEY,
        ProviderType TEXT NOT NULL,
        DerivationAlgorithm TEXT NOT NULL,
        MemoryCostKiB INTEGER NOT NULL,
        TimeCostIterations INTEGER NOT NULL,
        KeyLengthBytes INTEGER NOT NULL,
        FingerprintSHA256 TEXT NOT NULL,
        CreatedAt INTEGER NOT NULL,
        isHardwareBacked INTEGER NOT NULL CHECK (isHardwareBacked IN (0, 1)),
        hasActiveSession INTEGER NOT NULL CHECK (hasActiveSession IN (0, 1)),
        isDerivationValid INTEGER NOT NULL CHECK (isDerivationValid IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/mail/vault/enclave_manager.go` and integrates with `gitmap security vault` commands for credentials lifecycle management.

### Improvement 662: Automated Key Rotation Scheduler & Epoch Horizon Monitor
- **Title:** Automated Key Rotation Scheduler & Epoch Horizon Monitor
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Orchestrate zero-downtime re-encryption of stored mailbox secrets under rotating cryptographic epochs with configurable expiration horizons.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Epoch timeline bar showing active epoch age, days until mandatory rotation, and pending re-encryption queue depth with Amber countdown indicators.
  - *Web UI:* Zero-CDN timeline visualization of cryptographic key generations and re-encryption batch progress.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault rotate [--force] [--grace-period=72h] [--batch-size=50] [--json]`
  - Flags: `--force` (bool, trigger immediate epoch rotation), `--grace-period` (duration, backward compatibility window), `--batch-size` (int, re-encryption chunk size), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/rotation_types.go`):
    ```go
    type VaultRotationEpoch struct {
        EpochID            string `json:"epochId"`
        KeyID              string `json:"keyId"`
        ReEncryptedCount   int    `json:"reEncryptedCount"`
        TotalSecretCount   int    `json:"totalSecretCount"`
        GracePeriodSeconds int64  `json:"gracePeriodSeconds"`
        RotatedAt          int64  `json:"rotatedAt"`
        ExpiresAt          int64  `json:"expiresAt"`
        isCurrentEpoch     bool   `json:"isCurrentEpoch"`
        hasPendingSecrets  bool   `json:"hasPendingSecrets"`
        isRotationComplete bool   `json:"isRotationComplete"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationVaultEpochs (
        EpochID TEXT PRIMARY KEY,
        KeyID TEXT NOT NULL,
        ReEncryptedCount INTEGER NOT NULL,
        TotalSecretCount INTEGER NOT NULL,
        GracePeriodSeconds INTEGER NOT NULL,
        RotatedAt INTEGER NOT NULL,
        ExpiresAt INTEGER NOT NULL,
        isCurrentEpoch INTEGER NOT NULL CHECK (isCurrentEpoch IN (0, 1)),
        hasPendingSecrets INTEGER NOT NULL CHECK (hasPendingSecrets IN (0, 1)),
        isRotationComplete INTEGER NOT NULL CHECK (isRotationComplete IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/mail/vault/key_rotator.go` during scheduled background maintenance loops.

### Improvement 663: Zero-Swap In-Memory Secret Cache & Ephemeral Wiping
- **Title:** Zero-Swap In-Memory Secret Cache & Ephemeral Wiping
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Store decrypted mailbox credentials exclusively in `mlock()` memory buffers protected from OS disk paging with deterministic byte zeroing on eviction.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Memory buffer metrics display showing allocated mlock pages, cache hit ratio, byte zeroing latency, and memory page security flags in Violet and Cyan.
  - *Web UI:* Real-time SVG gauge indicating cache utilization, memory lock status, and active buffer count.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault cache-status [--purge] [--verify-mlock] [--json]`
  - Flags: `--purge` (bool, immediately overwrite all cached memory pages with zeros), `--verify-mlock` (bool, assert kernel page lock integrity), `--json` (bool, output status in JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/cache_types.go`):
    ```go
    type SecretCacheState struct {
        MailboxID         string `json:"mailboxId"`
        TenantID          string `json:"tenantId"`
        AllocatedBytes    int    `json:"allocatedBytes"`
        AccessCount       int    `json:"accessCount"`
        CachedAt          int64  `json:"cachedAt"`
        LastAccessedAt    int64  `json:"lastAccessedAt"`
        isMemoryLocked    bool   `json:"isMemoryLocked"`
        hasZeroOnEvict    bool   `json:"hasZeroOnEvict"`
        isCacheActive     bool   `json:"isCacheActive"`
    }
    ```
  - SQLite DDL (`gitmap.db` - ephemeral cache metrics log):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSecretCacheMetrics (
        MailboxID TEXT PRIMARY KEY,
        TenantID TEXT NOT NULL,
        AllocatedBytes INTEGER NOT NULL,
        AccessCount INTEGER NOT NULL,
        CachedAt INTEGER NOT NULL,
        LastAccessedAt INTEGER NOT NULL,
        isMemoryLocked INTEGER NOT NULL CHECK (isMemoryLocked IN (0, 1)),
        hasZeroOnEvict INTEGER NOT NULL CHECK (hasZeroOnEvict IN (0, 1)),
        isCacheActive INTEGER NOT NULL CHECK (isCacheActive IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/mail/vault/secret_cache.go` and utilized by SMTP/IMAP dialing workers.

### Improvement 664: Fine-Grained Casbin RBAC Tenant Isolation Gate
- **Title:** Fine-Grained Casbin RBAC Tenant Isolation Gate
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Enforce strict mathematical access separation preventing cross-tenant mailbox credential leakage across multi-tenant sender organizations.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Policy matrix renderer with tabular permission rows (`read`, `write`, `rotate`, `execute`), subject identities, and access decision audit logs in Emerald and Crimson.
  - *Web UI:* Zero-CDN permission grid showing role hierarchies, tenant scopes, and policy evaluation trace graphs.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault auth-check --tenant=<id> --user=<id> --action=<act> [--json]`
  - Flags: `--tenant` (string, tenant identifier), `--user` (string, requesting subject), `--action` (string, requested operation), `--json` (bool, JSON formatted output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/rbac_types.go`):
    ```go
    type VaultAccessPolicy struct {
        PolicyID          string `json:"policyId"`
        TenantID          string `json:"tenantId"`
        SubjectRole       string `json:"subjectRole"`
        ResourcePattern   string `json:"resourcePattern"`
        Action            string `json:"action"`
        CreatedAt         int64  `json:"createdAt"`
        isEnforced        bool   `json:"isEnforced"`
        hasStrictTenantIsolation bool `json:"hasStrictTenantIsolation"`
        isAuditRequired   bool   `json:"isAuditRequired"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapVaultAccessPolicies (
        PolicyID TEXT PRIMARY KEY,
        TenantID TEXT NOT NULL,
        SubjectRole TEXT NOT NULL,
        ResourcePattern TEXT NOT NULL,
        Action TEXT NOT NULL,
        CreatedAt INTEGER NOT NULL,
        isEnforced INTEGER NOT NULL CHECK (isEnforced IN (0, 1)),
        hasStrictTenantIsolation INTEGER NOT NULL CHECK (hasStrictTenantIsolation IN (0, 1)),
        isAuditRequired INTEGER NOT NULL CHECK (isAuditRequired IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Enforced in `pkg/mail/vault/casbin_gate.go` on every credential retrieval invocation.

### Improvement 665: Cryptographic Merkle-Linked Vault Audit Trail
- **Title:** Cryptographic Merkle-Linked Vault Audit Trail
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Log every credential access, rotation, and modification event into a tamper-evident append-only Merkle-linked audit log.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Cryptographic audit stream with block hashes, prior-block links, operator fingerprints, and tamper verification badges (`INTEGRITY VERIFIED` / `TAMPER DETECTED`).
  - *Web UI:* Interactive zero-CDN Merkle chain explorer displaying cryptographic branch verification and block receipts.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault audit-verify [--from-epoch=<id>] [--deep] [--json]`
  - Flags: `--from-epoch` (string, starting epoch hash), `--deep` (bool, verify entire historical chain back to genesis), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/audit_types.go`):
    ```go
    type VaultAuditEntry struct {
        EntryID           string `json:"entryId"`
        TenantID          string `json:"tenantId"`
        MailboxID         string `json:"mailboxId"`
        Action            string `json:"action"`
        OperatorID        string `json:"operatorId"`
        PreviousHash      string `json:"previousHash"`
        CurrentHash       string `json:"currentHash"`
        TimestampMillis   int64  `json:"timestampMillis"`
        isVerified        bool   `json:"isVerified"`
        hasValidSignature bool   `json:"hasValidSignature"`
        isTamperFree      bool   `json:"isTamperFree"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapVaultAuditEntries (
        EntryID TEXT PRIMARY KEY,
        TenantID TEXT NOT NULL,
        MailboxID TEXT NOT NULL,
        Action TEXT NOT NULL,
        OperatorID TEXT NOT NULL,
        PreviousHash TEXT NOT NULL,
        CurrentHash TEXT NOT NULL,
        TimestampMillis INTEGER NOT NULL,
        isVerified INTEGER NOT NULL CHECK (isVerified IN (0, 1)),
        hasValidSignature INTEGER NOT NULL CHECK (hasValidSignature IN (0, 1)),
        isTamperFree INTEGER NOT NULL CHECK (isTamperFree IN (0, 1))
    );
    CREATE INDEX IF NOT EXISTS idx_vault_audit_tenant ON GitmapVaultAuditEntries(TenantID);
    ```
- **Antigravity / GitMap Integration Point:**
  - Maintained by `pkg/mail/vault/audit_logger.go` with automated validation on GitMap audit runs.

### Improvement 666: OAuth2 Refresh Token Secret Lease & Scope Minimizer
- **Title:** OAuth2 Refresh Token Secret Lease & Scope Minimizer
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Manage short-lived access token renewal for Gmail and Microsoft Graph mailboxes with strict permission scope boundaries and automatic revocation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* OAuth token lifecycle meter showing lease time-to-live (TTL), refresh error counts, permitted scopes, and revocation status in Yellow and Slate tokens.
  - *Web UI:* Visual scope perimeter badge matrix and countdown timer for active OAuth2 leases.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault oauth-renew --mailbox=<id> [--force] [--json]`
  - Flags: `--mailbox` (string, target mailbox ID), `--force` (bool, trigger renewal regardless of TTL), `--json` (bool, output renewal payload in JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/oauth_types.go`):
    ```go
    type OAuthTokenLease struct {
        MailboxID         string `json:"mailboxId"`
        Provider          string `json:"provider"`
        GrantedScopes     string `json:"grantedScopes"`
        ExpiresAt         int64  `json:"expiresAt"`
        RefreshCount      int    `json:"refreshCount"`
        LastRefreshedAt   int64  `json:"lastRefreshedAt"`
        isLeaseActive     bool   `json:"isLeaseActive"`
        hasMinimalScopes  bool   `json:"hasMinimalScopes"`
        isRevoked         bool   `json:"isRevoked"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationOAuthLeases (
        MailboxID TEXT PRIMARY KEY,
        Provider TEXT NOT NULL,
        GrantedScopes TEXT NOT NULL,
        ExpiresAt INTEGER NOT NULL,
        RefreshCount INTEGER NOT NULL,
        LastRefreshedAt INTEGER NOT NULL,
        isLeaseActive INTEGER NOT NULL CHECK (isLeaseActive IN (0, 1)),
        hasMinimalScopes INTEGER NOT NULL CHECK (hasMinimalScopes IN (0, 1)),
        isRevoked INTEGER NOT NULL CHECK (isRevoked IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Used in `pkg/mail/vault/oauth_manager.go` during warmup authentication workflows.

### Improvement 667: Compromised Credential Canary & Leak Radar
- **Title:** Compromised Credential Canary & Leak Radar
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Deploy synthetic honeypot mailbox credentials and monitor external leak repositories or unauthorized vault read spikes to trigger automated lockouts.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Canary tripwire status indicator with Crimson flashing alert headers, honeytoken hit count, and auto-lockout circuit status.
  - *Web UI:* Threat map displaying simulated vs unauthorized access attempts and triggered defensive countermeasures.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault canary-status [--deploy-canary] [--reset-lockout] [--json]`
  - Flags: `--deploy-canary` (bool, inject decoy honeytoken into database), `--reset-lockout` (bool, clear administrative security lockdown), `--json` (bool, output status in JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/canary_types.go`):
    ```go
    type VaultCanaryTripwire struct {
        CanaryID          string `json:"canaryId"`
        DecoyMailboxName  string `json:"decoyMailboxName"`
        TriggerCount      int    `json:"triggerCount"`
        LastTriggeredAt   int64  `json:"lastTriggeredAt"`
        TriggeredIP       string `json:"triggeredIp"`
        CreatedAt         int64  `json:"createdAt"`
        isTripped         bool   `json:"isTripped"`
        hasLockoutFired   bool   `json:"hasLockoutFired"`
        isDecoyActive     bool   `json:"isDecoyActive"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapVaultCanaries (
        CanaryID TEXT PRIMARY KEY,
        DecoyMailboxName TEXT NOT NULL,
        TriggerCount INTEGER NOT NULL,
        LastTriggeredAt INTEGER NOT NULL,
        TriggeredIP TEXT NOT NULL,
        CreatedAt INTEGER NOT NULL,
        isTripped INTEGER NOT NULL CHECK (isTripped IN (0, 1)),
        hasLockoutFired INTEGER NOT NULL CHECK (hasLockoutFired IN (0, 1)),
        isDecoyActive INTEGER NOT NULL CHECK (isDecoyActive IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Monitored by `pkg/mail/vault/canary_sentinel.go` and integrated with Antigravity anomaly alerts.

### Improvement 668: Multi-Party Threshold Computation (MPC) Key Shard Assembler
- **Title:** Multi-Party Threshold Computation (MPC) Key Shard Assembler
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Split vault master encryption keys across Shamir secret shares requiring quorum authorization before decrypting high-security enterprise sender credentials.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Quorum status board with visual share tally (`3/5 shares present`), contributing node IDs, and threshold unlock authorization prompts in Fuchsia and Cyan.
  - *Web UI:* Zero-CDN cryptographic key reconstruction visualizer depicting Shamir polynomial evaluation.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault mpc-unlock --share=<share-payload> [--threshold=3] [--json]`
  - Flags: `--share` (string, armored Shamir share string), `--threshold` (int, minimum required quorum count), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/mpc_types.go`):
    ```go
    type MPCQuorumStatus struct {
        SessionID         string `json:"sessionId"`
        RequiredThreshold int    `json:"requiredThreshold"`
        TotalSharesCount  int    `json:"totalSharesCount"`
        SubmittedCount    int    `json:"submittedCount"`
        CreatedAt         int64  `json:"createdAt"`
        ExpiresAt         int64  `json:"expiresAt"`
        isQuorumMet       bool   `json:"isQuorumMet"`
        hasValidShares    bool   `json:"hasValidShares"`
        isSessionExpired  bool   `json:"isSessionExpired"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationMPCQuorums (
        SessionID TEXT PRIMARY KEY,
        RequiredThreshold INTEGER NOT NULL,
        TotalSharesCount INTEGER NOT NULL,
        SubmittedCount INTEGER NOT NULL,
        CreatedAt INTEGER NOT NULL,
        ExpiresAt INTEGER NOT NULL,
        isQuorumMet INTEGER NOT NULL CHECK (isQuorumMet IN (0, 1)),
        hasValidShares INTEGER NOT NULL CHECK (hasValidShares IN (0, 1)),
        isSessionExpired INTEGER NOT NULL CHECK (isSessionExpired IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/mail/vault/mpc_assembler.go` for air-gapped or high-security deployments.

### Improvement 669: Ephemeral SMTP Session Credential Injector & Proxy Wire
- **Title:** Ephemeral SMTP Session Credential Injector & Proxy Wire
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Inject credentials directly into outgoing SMTP TLS streams via loopback proxy pipe without ever exposing plaintext secrets to client code or debug logs.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Stream injection monitor reporting active proxy pipe sessions, zero-copy buffer handoffs, and injection latency stats in Emerald (`< 2ms`).
  - *Web UI:* Real-time dataflow diagram of the internal socket proxy bridge and encrypted stream boundaries.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault proxy-stats [--stream-id=<id>] [--json]`
  - Flags: `--stream-id` (string, specific active socket ID), `--json` (bool, output proxy metrics in JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/proxy_types.go`):
    ```go
    type SMTPProxyInjection struct {
        SessionID         string `json:"sessionId"`
        MailboxID         string `json:"mailboxId"`
        InjectedAuthType  string `json:"injectedAuthType"`
        BytesTransferred  int64  `json:"bytesTransferred"`
        LatencyMicros     int64  `json:"latencyMicros"`
        TimestampMillis   int64  `json:"timestampMillis"`
        isEncryptedStream bool   `json:"isEncryptedStream"`
        hasCleanBuffer    bool   `json:"hasCleanBuffer"`
        isSessionClosed   bool   `json:"isSessionClosed"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSMTPProxyInjections (
        SessionID TEXT PRIMARY KEY,
        MailboxID TEXT NOT NULL,
        InjectedAuthType TEXT NOT NULL,
        BytesTransferred INTEGER NOT NULL,
        LatencyMicros INTEGER NOT NULL,
        TimestampMillis INTEGER NOT NULL,
        isEncryptedStream INTEGER NOT NULL CHECK (isEncryptedStream IN (0, 1)),
        hasCleanBuffer INTEGER NOT NULL CHECK (hasCleanBuffer IN (0, 1)),
        isSessionClosed INTEGER NOT NULL CHECK (isSessionClosed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `pkg/mail/vault/proxy_wire.go` and consumed by the SMTP connection pool.

### Improvement 670: Real-Time Hardware Enclave HUD & Terminal Dashboard
- **Title:** Real-Time Hardware Enclave HUD & Terminal Dashboard
- **Cluster:** Cluster 67: Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD
- **Goal:** Render a high-density Bubble Tea / Lipgloss terminal HUD summarizing vault health, key epoch age, active secret leases, and audit verification status.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Unified security cockpit featuring three styled border panels: Enclave Hardware Health (Top), Active Quorums & Secret Leases (Center), and Live Cryptographic Audit Stream (Bottom) with zero ANSI flickering.
  - *Web UI:* Responsive HTML/SVG dashboard with embedded security token summaries and real-time event log viewer.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup vault hud [--refresh=2s] [--export-html=<path>]`
  - Flags: `--refresh` (duration, HUD redraw frequency), `--export-html` (string, optional static HTML diagnostic report output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/vault/hud_types.go`):
    ```go
    type EnclaveHUDModel struct {
        ActiveKeyEpoch    string `json:"activeKeyEpoch"`
        EnclaveType       string `json:"enclaveType"`
        CachedSecretCount int    `json:"cachedSecretCount"`
        TotalTenantsCount int    `json:"totalTenantsCount"`
        AuditEntriesCount int64  `json:"auditEntriesCount"`
        UptimeSeconds     int64  `json:"uptimeSeconds"`
        isEnclaveHealthy  bool   `json:"isEnclaveHealthy"`
        hasAuditTamper    bool   `json:"hasAuditTamper"`
        isHardwareSecured bool   `json:"isHardwareSecured"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapVaultHUDSnapshots (
        SnapshotID TEXT PRIMARY KEY,
        ActiveKeyEpoch TEXT NOT NULL,
        EnclaveType TEXT NOT NULL,
        CachedSecretCount INTEGER NOT NULL,
        TotalTenantsCount INTEGER NOT NULL,
        AuditEntriesCount INTEGER NOT NULL,
        UptimeSeconds INTEGER NOT NULL,
        RecordedAt INTEGER NOT NULL,
        isEnclaveHealthy INTEGER NOT NULL CHECK (isEnclaveHealthy IN (0, 1)),
        hasAuditTamper INTEGER NOT NULL CHECK (hasAuditTamper IN (0, 1)),
        isHardwareSecured INTEGER NOT NULL CHECK (isHardwareSecured IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Primary UI entry point in `cli/cmdmail/vault_cmd.go` executing `pkg/mail/vault/enclave_hud.go`.

---

## 3. Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer (Items 671–680)

### Improvement 671: RFC 5322 / RFC 5321 Received Header Deep Tokenizer
- **Title:** RFC 5322 / RFC 5321 Received Header Deep Tokenizer
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Dissect complex multi-line SMTP `Received:` headers into structured components (from, by, via, with, id, for, timestamp, and TLS parameters) with zero regex backtrack vulnerabilities.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Colorized header breakdown showing parsed tokens in syntax-highlighted blocks: Host (Cyan), Protocol (Violet), Cipher (Emerald), and Timestamp (Amber).
  - *Web UI:* Zero-CDN expandable token tree visualizing hierarchical BNF grammar components of email headers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace parse-header --file=<raw-eml> [--validate-rfc] [--json]`
  - Flags: `--file` (string, path to raw .eml or header file), `--validate-rfc` (bool, check strict RFC grammar compliance), `--json` (bool, output parsed AST in JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/tokenizer_types.go`):
    ```go
    type ReceivedHopToken struct {
        HopIndex         int    `json:"hopIndex"`
        FromHost         string `json:"fromHost"`
        FromIP           string `json:"fromIp"`
        ByHost           string `json:"byHost"`
        ProtocolWith     string `json:"protocolWith"`
        MessageID        string `json:"messageId"`
        TimestampUTC     int64  `json:"timestampUtc"`
        isRFCCompliant   bool   `json:"isRfcCompliant"`
        hasTlsProtocol   bool   `json:"hasTlsProtocol"`
        isForgedHeader   bool   `json:"isForgedHeader"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineReceivedHopTokens (
        TokenID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        HopIndex INTEGER NOT NULL,
        FromHost TEXT NOT NULL,
        FromIP TEXT NOT NULL,
        ByHost TEXT NOT NULL,
        ProtocolWith TEXT NOT NULL,
        MessageID TEXT NOT NULL,
        TimestampUTC INTEGER NOT NULL,
        isRFCCompliant INTEGER NOT NULL CHECK (isRFCCompliant IN (0, 1)),
        hasTlsProtocol INTEGER NOT NULL CHECK (hasTlsProtocol IN (0, 1)),
        isForgedHeader INTEGER NOT NULL CHECK (isForgedHeader IN (0, 1))
    );
    CREATE INDEX IF NOT EXISTS idx_hop_tokens_email ON PipelineReceivedHopTokens(EmailID);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/mail/forensics/header_tokenizer.go` as the core parsing engine for forensic traces.

### Improvement 672: Multi-Hop Latency Waterfall Timeline Visualizer
- **Title:** Multi-Hop Latency Waterfall Timeline Visualizer
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Calculate inter-gateway transit durations from sequential hop timestamps and visualize delivery bottlenecks in a terminal waterfall timeline.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Waterfall Gantt chart with proportional horizontal bar characters (`█▓░`), delta elapsed times (`+120ms`, `+4.2s`), and color thresholds (Green `< 500ms`, Yellow `500ms–3s`, Red `> 3s`).
  - *Web UI:* Responsive SVG waterfall timeline diagram displaying precise millisecond transit spans between SMTP relays.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace waterfall --msg-id=<id> [--unit=ms] [--threshold=2000] [--json]`
  - Flags: `--msg-id` (string, target email message ID), `--unit` (string, display unit `ms` or `s`), `--threshold` (int, alert latency boundary in ms), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/waterfall_types.go`):
    ```go
    type HopLatencySpan struct {
        FromHopIndex     int    `json:"fromHopIndex"`
        ToHopIndex       int    `json:"toHopIndex"`
        RelayHost        string `json:"relayHost"`
        DurationMillis   int64  `json:"durationMillis"`
        CumulativeMillis int64  `json:"cumulativeMillis"`
        isBottleneck     bool   `json:"isBottleneck"`
        hasClockSkew     bool   `json:"hasClockSkew"`
        isNormalTransit  bool   `json:"isNormalTransit"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineHopLatencySpans (
        SpanID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        FromHopIndex INTEGER NOT NULL,
        ToHopIndex INTEGER NOT NULL,
        RelayHost TEXT NOT NULL,
        DurationMillis INTEGER NOT NULL,
        CumulativeMillis INTEGER NOT NULL,
        isBottleneck INTEGER NOT NULL CHECK (isBottleneck IN (0, 1)),
        hasClockSkew INTEGER NOT NULL CHECK (hasClockSkew IN (0, 1)),
        isNormalTransit INTEGER NOT NULL CHECK (isNormalTransit IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/mail/forensics/hop_timeline.go` and surfaced in warmup delivery diagnostic reports.

### Improvement 673: Offline Autonomous ASN & BGP Route Enrichment Sentinel
- **Title:** Offline Autonomous ASN & BGP Route Enrichment Sentinel
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Enrich relay IP addresses with Autonomous System Numbers (ASN), organization owners, and geographic coordinates without external network queries using an embedded MaxMind / IP2Location database.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Network route summary table featuring ASN badges (`AS15169 Google LLC`), country codes, and BGP prefix allocations in Blue and Slate tokens.
  - *Web UI:* Interactive zero-CDN SVG world map tracing email geographical routing across intermediate continents and data centers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace asn-lookup --ip=<ip> [--db-path=<path>] [--json]`
  - Flags: `--ip` (string, IP address to inspect), `--db-path` (string, override local offline MMDB path), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/asn_types.go`):
    ```go
    type RelayASNEnrichment struct {
        IPAddress        string `json:"ipAddress"`
        ASN              int    `json:"asn"`
        OrgName          string `json:"orgName"`
        CountryCode      string `json:"countryCode"`
        CityName         string `json:"cityName"`
        isKnownDatacenter bool  `json:"isKnownDatacenter"`
        hasValidBgpRoute  bool  `json:"hasValidBgpRoute"`
        isResidentialIP   bool  `json:"isResidentialIp"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationRelayASNCache (
        IPAddress TEXT PRIMARY KEY,
        ASN INTEGER NOT NULL,
        OrgName TEXT NOT NULL,
        CountryCode TEXT NOT NULL,
        CityName TEXT NOT NULL,
        CachedAt INTEGER NOT NULL,
        isKnownDatacenter INTEGER NOT NULL CHECK (isKnownDatacenter IN (0, 1)),
        hasValidBgpRoute INTEGER NOT NULL CHECK (hasValidBgpRoute IN (0, 1)),
        isResidentialIP INTEGER NOT NULL CHECK (isResidentialIP IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Provided by `pkg/mail/forensics/asn_lookup.go` for air-gapped forensic investigations.

### Improvement 674: Intermediate TLS Cipher Suite & Forward Secrecy Inspector
- **Title:** Intermediate TLS Cipher Suite & Forward Secrecy Inspector
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Analyze TLS protocol versions (TLS 1.2, TLS 1.3), key exchange algorithms (ECDHE), and cipher suites reported across intermediate relay hops to identify downgraded or insecure transmission links.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Security grading matrix assigning letter grades (`A+`, `A`, `B`, `F`) based on cipher strength, PFS (Perfect Forward Secrecy) presence, and key exchange bit length in Green and Amber.
  - *Web UI:* Visual cipher suite analyzer highlighting weak ciphers (e.g., CBC mode, 3DES) and protocol downgrade risks.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace tls-audit --file=<raw-eml> [--min-grade=A] [--json]`
  - Flags: `--file` (string, email raw content file), `--min-grade` (string, minimum acceptable security grade), `--json` (bool, output TLS audit in JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/tls_types.go`):
    ```go
    type RelayTLSSecurity struct {
        HopIndex         int    `json:"hopIndex"`
        TLSVersion       string `json:"tlsVersion"`
        CipherSuite      string `json:"cipherSuite"`
        KeyExchange      string `json:"keyExchange"`
        SecurityGrade    string `json:"securityGrade"`
        isPfsSupported   bool   `json:"isPfsSupported"`
        hasStrongCipher  bool   `json:"hasStrongCipher"`
        isDowngradeRisk  bool   `json:"isDowngradeRisk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineRelayTLSSecurity (
        SecurityID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        HopIndex INTEGER NOT NULL,
        TLSVersion TEXT NOT NULL,
        CipherSuite TEXT NOT NULL,
        KeyExchange TEXT NOT NULL,
        SecurityGrade TEXT NOT NULL,
        isPfsSupported INTEGER NOT NULL CHECK (isPfsSupported IN (0, 1)),
        hasStrongCipher INTEGER NOT NULL CHECK (hasStrongCipher IN (0, 1)),
        isDowngradeRisk INTEGER NOT NULL CHECK (isDowngradeRisk IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/mail/forensics/tls_cipher_inspector.go` to flag plaintext or insecure intermediate hops.

### Improvement 675: Synthetic Relay Injection & Forged Hop Heuristic Detector
- **Title:** Synthetic Relay Injection & Forged Hop Heuristic Detector
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Detect spoofed `Received:` headers prepended by malicious actors or misconfigured proxies by correlating timestamps, reverse DNS PTR records, and connecting socket IPs.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Fraud diagnostic box highlighting suspected forged headers in Bright Crimson with detailed discrepancy reasons (`Timestamp inversion`, `PTR mismatch`, `Unreachable upstream gateway`).
  - *Web UI:* Visual step-by-step trust chain diagram with red broken links at detected points of header forgery.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace detect-forgery --file=<raw-eml> [--strict] [--json]`
  - Flags: `--file` (string, raw email file), `--strict` (bool, enforce microsecond monotonic clock requirements), `--json` (bool, JSON formatted output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/forgery_types.go`):
    ```go
    type HeaderForgeryAnalysis struct {
        SuspectedHopIndex int    `json:"suspectedHopIndex"`
        DiscrepancyType   string `json:"discrepancyType"`
        ReportedHost      string `json:"reportedHost"`
        ActualSocketIP    string `json:"actualSocketIp"`
        ConfidenceScore   float64 `json:"confidenceScore"`
        isForgedHop       bool   `json:"isForgedHop"`
        hasClockInversion bool   `json:"hasClockInversion"`
        isChainBroken     bool   `json:"isChainBroken"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineHeaderForgeryAnalyses (
        AnalysisID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        SuspectedHopIndex INTEGER NOT NULL,
        DiscrepancyType TEXT NOT NULL,
        ReportedHost TEXT NOT NULL,
        ActualSocketIP TEXT NOT NULL,
        ConfidenceScore REAL NOT NULL,
        isForgedHop INTEGER NOT NULL CHECK (isForgedHop IN (0, 1)),
        hasClockInversion INTEGER NOT NULL CHECK (hasClockInversion IN (0, 1)),
        isChainBroken INTEGER NOT NULL CHECK (isChainBroken IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Maintained by `pkg/mail/forensics/forgery_detector.go` to safeguard warmup seed authenticity.

### Improvement 676: Gateway Boundary Transit & Network Border Demarcator
- **Title:** Gateway Boundary Transit & Network Border Demarcator
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Classify relay hops into functional network zones (Internal Sender, Outbound Smarthost, Intermediate Transit, Recipient Edge MX, Internal Mailbox Delivery).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Zone segmentation badges displaying distinct color bands: Green (Origin Internal), Blue (Edge Smarthost), Violet (External Internet Transit), and Cyan (Target Ingress).
  - *Web UI:* Interactive zone diagram showing packet boundaries across corporate firewalls and cloud edge routers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace classify-zones --file=<raw-eml> [--json]`
  - Flags: `--file` (string, target email file), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/zone_types.go`):
    ```go
    type GatewayZoneHop struct {
        HopIndex         int    `json:"hopIndex"`
        ZoneClassification string `json:"zoneClassification"`
        GatewayRole      string `json:"gatewayRole"`
        Hostname         string `json:"hostname"`
        isInternalNetwork bool   `json:"isInternalNetwork"`
        isBorderGateway  bool   `json:"isBorderGateway"`
        isPublicTransit  bool   `json:"isPublicTransit"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineGatewayZoneHops (
        ZoneHopID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        HopIndex INTEGER NOT NULL,
        ZoneClassification TEXT NOT NULL,
        GatewayRole TEXT NOT NULL,
        Hostname TEXT NOT NULL,
        isInternalNetwork INTEGER NOT NULL CHECK (isInternalNetwork IN (0, 1)),
        isBorderGateway INTEGER NOT NULL CHECK (isBorderGateway IN (0, 1)),
        isPublicTransit INTEGER NOT NULL CHECK (isPublicTransit IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Handled by `pkg/mail/forensics/zone_demarcator.go` for network topology categorization.

### Improvement 677: DSN & MTA Diagnostic Code Correlator
- **Title:** DSN & MTA Diagnostic Code Correlator
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Link `Received:` hop delays with subsequent RFC 3463 / 3464 DSN status codes and MTA error logs to locate exactly which hop generated delivery failures or greylisting.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Correlated failure card showing target relay hostname, DSN status (`5.7.1 Service unavailable`), raw SMTP response code (`554`), and exact hop failure index.
  - *Web UI:* Visual trace graph pointing an error marker directly at the specific node responsible for the bounce.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace correlate-dsn --dsn-file=<path> --header-file=<path> [--json]`
  - Flags: `--dsn-file` (string, path to DSN bounce message), `--header-file` (string, path to original header trace), `--json` (bool, JSON formatted output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/dsn_types.go`):
    ```go
    type HopDSNCorrelation struct {
        EmailID          string `json:"emailId"`
        FailingHopIndex  int    `json:"failingHopIndex"`
        FailingRelayHost string `json:"failingRelayHost"`
        DSNStatus        string `json:"dsnStatus"`
        SMTPCode         int    `json:"smtpCode"`
        DiagnosticText   string `json:"diagnosticText"`
        isTerminalError  bool   `json:"isTerminalError"`
        hasGreylisting   bool   `json:"hasGreylisting"`
        isPolicyRejection bool  `json:"isPolicyRejection"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineHopDSNCorrelations (
        CorrelationID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        FailingHopIndex INTEGER NOT NULL,
        FailingRelayHost TEXT NOT NULL,
        DSNStatus TEXT NOT NULL,
        SMTPCode INTEGER NOT NULL,
        DiagnosticText TEXT NOT NULL,
        isTerminalError INTEGER NOT NULL CHECK (isTerminalError IN (0, 1)),
        hasGreylisting INTEGER NOT NULL CHECK (hasGreylisting IN (0, 1)),
        isPolicyRejection INTEGER NOT NULL CHECK (isPolicyRejection IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/mail/forensics/dsn_correlator.go` and cross-references bounce classifier events.

### Improvement 678: DKIM / ARC Verification Pass-Through Validator
- **Title:** DKIM / ARC Verification Pass-Through Validator
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Trace cryptographic verification results recorded by intermediate gateways in `Authentication-Results` headers across each hop to detect signature invalidation mid-flight.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Cryptographic status badge row across hops: Hop 1 `DKIM: pass`, Hop 2 `ARC: pass`, Hop 3 `DKIM: fail (body altered)` in Green and Red.
  - *Web UI:* Interactive integrity pass-through visualizer showing message modification points and canonicalization failures.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace auth-passthrough --file=<raw-eml> [--json]`
  - Flags: `--file` (string, raw email file), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/auth_passthrough_types.go`):
    ```go
    type HopAuthResultsPassThrough struct {
        HopIndex         int    `json:"hopIndex"`
        AuthServer       string `json:"authServer"`
        DKIMResult       string `json:"dkimResult"`
        SPFResult        string `json:"spfResult"`
        ARCResult        string `json:"arcResult"`
        DMARCResult      string `json:"dmarcResult"`
        isDkimValid      bool   `json:"isDkimValid"`
        hasArcSeal       bool   `json:"hasArcSeal"`
        isBodyUnaltered  bool   `json:"isBodyUnaltered"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineHopAuthPassThrough (
        PassThroughID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        HopIndex INTEGER NOT NULL,
        AuthServer TEXT NOT NULL,
        DKIMResult TEXT NOT NULL,
        SPFResult TEXT NOT NULL,
        ARCResult TEXT NOT NULL,
        DMARCResult TEXT NOT NULL,
        isDkimValid INTEGER NOT NULL CHECK (isDkimValid IN (0, 1)),
        hasArcSeal INTEGER NOT NULL CHECK (hasArcSeal IN (0, 1)),
        isBodyUnaltered INTEGER NOT NULL CHECK (isBodyUnaltered IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/mail/forensics/auth_passthrough.go` to pinpoint intermediate mail munging.

### Improvement 679: Header Anonymization & Sensitive Redaction Engine
- **Title:** Header Anonymization & Sensitive Redaction Engine
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Redact internal IP addresses, operator usernames, message IDs, and confidential domains from raw email headers before exporting traces for external troubleshooting or bug reports.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Redaction review panel showing original text alongside masked output (`10.0.4.15` -> `[REDACTED_IPV4_1]`, `user@corp.internal` -> `[REDACTED_USER_1]`) in Muted Yellow.
  - *Web UI:* Split-pane differential viewer highlighting redacted PII (Personally Identifiable Information) tokens.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace redact --input=<file> --output=<out-file> [--mask-ips] [--mask-domains] [--json]`
  - Flags: `--input` (string, source header file), `--output` (string, sanitized output destination), `--mask-ips` (bool, redact all internal/private IP blocks), `--mask-domains` (bool, mask custom domain names), `--json` (bool, output redaction mapping).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/forensics/redact_types.go`):
    ```go
    type RedactionSummary struct {
        TotalTokensRedacted int    `json:"totalTokensRedacted"`
        RedactedIPCount     int    `json:"redactedIpCount"`
        RedactedDomainCount int    `json:"redactedDomainCount"`
        SaltHash            string `json:"saltHash"`
        ProcessedAt         int64  `json:"processedAt"`
        isZeroPiiPreserved  bool   `json:"isZeroPiiPreserved"`
        hasDeterministicSalt bool  `json:"hasDeterministicSalt"`
        isSanitizedClean    bool   `json:"isSanitizedClean"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapRedactionLogs (
        LogID TEXT PRIMARY KEY,
        FileHash TEXT NOT NULL,
        TotalTokensRedacted INTEGER NOT NULL,
        RedactedIPCount INTEGER NOT NULL,
        RedactedDomainCount INTEGER NOT NULL,
        ProcessedAt INTEGER NOT NULL,
        isZeroPiiPreserved INTEGER NOT NULL CHECK (isZeroPiiPreserved IN (0, 1)),
        hasDeterministicSalt INTEGER NOT NULL CHECK (hasDeterministicSalt IN (0, 1)),
        isSanitizedClean INTEGER NOT NULL CHECK (isSanitizedClean IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `pkg/mail/forensics/redactor.go` and enforced prior to uploading diagnostics to shared repositories.

### Improvement 680: Interactive Full-Screen Forensic Hop Visualizer & Inspector
- **Title:** Interactive Full-Screen Forensic Hop Visualizer & Inspector
- **Cluster:** Cluster 68: Forensic Header Trace & Cross-Gateway Hop Visualizer
- **Goal:** Deliver a comprehensive Bubble Tea full-screen interactive explorer allowing engineers to step through email transit hops, inspect TLS ciphers, view latency bars, and examine raw BNF tokens in real time.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Three-column interactive TUI layout: Hop Hierarchy Tree (Left), Forensic Node Detail Inspector (Center), and Hex/Raw Header Viewport (Right) with keyboard navigation (`j`/`k` select hop, `t` toggle TLS view, `r` redact).
  - *Web UI:* Zero-CDN responsive HTML/SVG interactive explorer with pan-and-zoom transit route graph and click-to-inspect modal inspectors.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trace visualizer --file=<raw-eml> [--export-svg=<path>]`
  - Flags: `--file` (string, raw email file path), `--export-svg` (string, optional SVG export path for documentation).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/trace/visualizer_types.go`):
    ```go
    type ForensicVisualizerState struct {
        EmailID            string `json:"emailId"`
        TotalHopsCount     int    `json:"totalHopsCount"`
        SelectedHopIndex   int    `json:"selectedHopIndex"`
        ActiveTab          string `json:"activeTab"`
        RenderWidth        int    `json:"renderWidth"`
        RenderHeight       int    `json:"renderHeight"`
        isInteractive      bool   `json:"isInteractive"`
        hasDetectedAnomaly bool   `json:"hasDetectedAnomaly"`
        isTuiRendered      bool   `json:"isTuiRendered"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PipelineVisualizerSessions (
        SessionID TEXT PRIMARY KEY,
        EmailID TEXT NOT NULL,
        TotalHopsCount INTEGER NOT NULL,
        SelectedHopIndex INTEGER NOT NULL,
        LaunchedAt INTEGER NOT NULL,
        isInteractive INTEGER NOT NULL CHECK (isInteractive IN (0, 1)),
        hasDetectedAnomaly INTEGER NOT NULL CHECK (hasDetectedAnomaly IN (0, 1)),
        isTuiRendered INTEGER NOT NULL CHECK (isTuiRendered IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Primary TUI command in `cli/cmdmail/trace_cmd.go` executing `pkg/mail/trace/visualizer.go`.

---

## 4. Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver (Items 681–690)

### Improvement 681: Consistent Hashing Shard Topology & Virtual Ring Coordinator
- **Title:** Consistent Hashing Shard Topology & Virtual Ring Coordinator
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Distribute email warm-up accounts, logs, and telemetry across dynamic SQLite physical shard partitions using a consistent hash ring with virtual nodes (vnodes) to minimize data movement during cluster scaling.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Shard partition ring visualization displaying virtual node tokens, shard ownership percentages, key allocation distributions, and imbalance metrics in Violet, Cyan, and Slate.
  - *Web UI:* Zero-CDN SVG chord diagram depicting virtual ring partitions, node health, and token assignments.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard topology [--nodes=4] [--vnodes=128] [--rebalance] [--json]`
  - Flags: `--nodes` (int, target physical node count), `--vnodes` (int, virtual tokens per node), `--rebalance` (bool, execute key reassignment simulation), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/topology_types.go`):
    ```go
    type ShardRingNode struct {
        NodeID             string `json:"nodeId"`
        PhysicalAddress    string `json:"physicalAddress"`
        VirtualNodeCount   int    `json:"virtualNodeCount"`
        AssignedKeyRange   string `json:"assignedKeyRange"`
        ShardSizeBytes     int64  `json:"shardSizeBytes"`
        isMasterLeader     bool   `json:"isMasterLeader"`
        isShardOnline      bool   `json:"isShardOnline"`
        hasActiveReplicas  bool   `json:"hasActiveReplicas"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationShardTopology (
        NodeID TEXT PRIMARY KEY,
        PhysicalAddress TEXT NOT NULL,
        VirtualNodeCount INTEGER NOT NULL,
        AssignedKeyRange TEXT NOT NULL,
        ShardSizeBytes INTEGER NOT NULL,
        UpdatedAt INTEGER NOT NULL,
        isMasterLeader INTEGER NOT NULL CHECK (isMasterLeader IN (0, 1)),
        isShardOnline INTEGER NOT NULL CHECK (isShardOnline IN (0, 1)),
        hasActiveReplicas INTEGER NOT NULL CHECK (hasActiveReplicas IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/db/sharding/shard_topology.go` for cluster-wide database sharding.

### Improvement 682: Streaming SQLite WAL Frame Replicator & Frame Log Buffer
- **Title:** Streaming SQLite WAL Frame Replicator & Frame Log Buffer
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Replicate SQLite Write-Ahead Log (WAL) frames in real-time from leader nodes to read replicas via zero-copy Unix sockets or TCP streams with sub-millisecond lag.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Live replication lag table showing frame sequence numbers (`LSN`), streaming throughput (`frames/sec`), buffer backpressure meters, and replica sync states in Emerald and Amber.
  - *Web UI:* Real-time SVG throughput gauge indicating frame lag and network buffer occupancy.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard stream-status [--node=<id>] [--tail-frames=20] [--json]`
  - Flags: `--node` (string, target cluster node), `--tail-frames` (int, number of recent WAL frames to inspect), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/wal_stream_types.go`):
    ```go
    type WALReplicationStream struct {
        StreamID          string `json:"streamId"`
        NodeID            string `json:"nodeId"`
        CurrentLSN        int64  `json:"currentLsn"`
        ReplicatedLSN     int64  `json:"replicatedLsn"`
        LagFrameCount     int64  `json:"lagFrameCount"`
        TransitDurationMs int64  `json:"transitDurationMs"`
        isStreamHealthy   bool   `json:"isStreamHealthy"`
        hasBackpressure   bool   `json:"hasBackpressure"`
        isZeroLag         bool   `json:"isZeroLag"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationWALReplicationStreams (
        StreamID TEXT PRIMARY KEY,
        NodeID TEXT NOT NULL,
        CurrentLSN INTEGER NOT NULL,
        ReplicatedLSN INTEGER NOT NULL,
        LagFrameCount INTEGER NOT NULL,
        TransitDurationMs INTEGER NOT NULL,
        LastHeartbeatAt INTEGER NOT NULL,
        isStreamHealthy INTEGER NOT NULL CHECK (isStreamHealthy IN (0, 1)),
        hasBackpressure INTEGER NOT NULL CHECK (hasBackpressure IN (0, 1)),
        isZeroLag INTEGER NOT NULL CHECK (isZeroLag IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Handled by `pkg/db/sharding/wal_streamer.go` to provide high-availability reads across GitMap instances.

### Improvement 683: Page-Level Cryptographic Merkle-Tree Conflict Detector
- **Title:** Page-Level Cryptographic Merkle-Tree Conflict Detector
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Compute incremental cryptographic Merkle trees over SQLite 4096-byte database pages to pinpoint diverging database states across nodes without transmitting entire database files.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Merkle tree diff viewer displaying branch root hashes, differing leaf page IDs, and byte variance summaries in Amber and Cyan.
  - *Web UI:* Zero-CDN interactive SVG Merkle tree visualizing diverging cryptographic tree paths and divergent leaf nodes.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard merkle-diff --primary=<id> --secondary=<id> [--json]`
  - Flags: `--primary` (string, primary node ID), `--secondary` (string, replica node ID), `--json` (bool, JSON formatted diff payload).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/merkle/tree_types.go`):
    ```go
    type ShardMerkleTree struct {
        TreeID            string `json:"treeId"`
        NodeID            string `json:"nodeId"`
        RootHash          string `json:"rootHash"`
        TotalPagesCount   int    `json:"totalPagesCount"`
        DivergentPagesCount int  `json:"divergentPagesCount"`
        GeneratedAt       int64  `json:"generatedAt"`
        isTreeConsistent  bool   `json:"isTreeConsistent"`
        hasDivergentPages bool   `json:"hasDivergentPages"`
        isVerifiedValid   bool   `json:"isVerifiedValid"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationShardMerkleTrees (
        TreeID TEXT PRIMARY KEY,
        NodeID TEXT NOT NULL,
        RootHash TEXT NOT NULL,
        TotalPagesCount INTEGER NOT NULL,
        DivergentPagesCount INTEGER NOT NULL,
        GeneratedAt INTEGER NOT NULL,
        isTreeConsistent INTEGER NOT NULL CHECK (isTreeConsistent IN (0, 1)),
        hasDivergentPages INTEGER NOT NULL CHECK (hasDivergentPages IN (0, 1)),
        isVerifiedValid INTEGER NOT NULL CHECK (isVerifiedValid IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Maintained by `pkg/db/merkle/tree_builder.go` for continuous data integrity reconciliation.

### Improvement 684: Deterministic Last-Write-Wins (LWW) & Vector Clock Conflict Resolver
- **Title:** Deterministic Last-Write-Wins (LWW) & Vector Clock Conflict Resolver
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Reconcile divergent multi-node SQLite records using Lamport timestamps, vector clocks, and deterministic last-write-wins rules while preserving complete audit histories.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Conflict resolution resolution docket summarizing reconciled entity keys, resolved field values, vector clock vectors, and timestamp deltas in Emerald and Slate.
  - *Web UI:* Side-by-side three-way merge visualizer showing conflicting record mutations and automatic merge decisions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard resolve-conflicts [--table=<table>] [--dry-run] [--json]`
  - Flags: `--table` (string, database table to reconcile), `--dry-run` (bool, simulate conflict resolution without writes), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/merkle/conflict_types.go`):
    ```go
    type RecordConflictResolution struct {
        ResolutionID     string `json:"resolutionId"`
        TableName        string `json:"tableName"`
        RecordPrimaryKey string `json:"recordPrimaryKey"`
        WinningNodeID    string `json:"winningNodeId"`
        WinningClock     int64  `json:"winningClock"`
        ResolvedAt       int64  `json:"resolvedAt"`
        isConflictResolved bool `json:"isConflictResolved"`
        hasVectorClockTie  bool `json:"hasVectorClockTie"`
        isDeterministicLww bool `json:"isDeterministicLww"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapConflictResolutions (
        ResolutionID TEXT PRIMARY KEY,
        TableName TEXT NOT NULL,
        RecordPrimaryKey TEXT NOT NULL,
        WinningNodeID TEXT NOT NULL,
        WinningClock INTEGER NOT NULL,
        ResolvedAt INTEGER NOT NULL,
        isConflictResolved INTEGER NOT NULL CHECK (isConflictResolved IN (0, 1)),
        hasVectorClockTie INTEGER NOT NULL CHECK (hasVectorClockTie IN (0, 1)),
        isDeterministicLww INTEGER NOT NULL CHECK (isDeterministicLww IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/db/merkle/conflict_resolver.go` during partition recovery cycles.

### Improvement 685: Distributed Two-Phase Epoch Leader Lease Arbiter
- **Title:** Distributed Two-Phase Epoch Leader Lease Arbiter
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Manage physical shard write leases using non-blocking two-phase commit heartbeats with epoch numbers to prevent split-brain dual-leader write scenarios.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Leader lease status board displaying active epoch ID, leader node hostname, lease expiration countdown (`TTL 4.8s`), and voting member acknowledgments in Green and Yellow.
  - *Web UI:* Visual distributed consensus state machine diagram with real-time election status updates.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard lease-status [--epoch=<id>] [--step-down] [--json]`
  - Flags: `--epoch` (string, specific epoch query), `--step-down` (bool, trigger voluntary leadership handover), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/epoch_types.go`):
    ```go
    type ShardLeaderEpochLease struct {
        EpochNumber      int64  `json:"epochNumber"`
        ShardID          string `json:"shardId"`
        LeaderNodeID     string `json:"leaderNodeId"`
        LeaseDurationMs  int64  `json:"leaseDurationMs"`
        AcquiredAt       int64  `json:"acquiredAt"`
        ExpiresAt        int64  `json:"expiresAt"`
        isLeaseHeld      bool   `json:"isLeaseHeld"`
        hasQuorumVotes   bool   `json:"hasQuorumVotes"`
        isSplitBrainSafe bool   `json:"isSplitBrainSafe"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationShardLeaderEpochs (
        EpochNumber INTEGER PRIMARY KEY,
        ShardID TEXT NOT NULL,
        LeaderNodeID TEXT NOT NULL,
        LeaseDurationMs INTEGER NOT NULL,
        AcquiredAt INTEGER NOT NULL,
        ExpiresAt INTEGER NOT NULL,
        isLeaseHeld INTEGER NOT NULL CHECK (isLeaseHeld IN (0, 1)),
        hasQuorumVotes INTEGER NOT NULL CHECK (hasQuorumVotes IN (0, 1)),
        isSplitBrainSafe INTEGER NOT NULL CHECK (isSplitBrainSafe IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Enforced in `pkg/db/sharding/epoch_lease.go` prior to allowing write transactions.

### Improvement 686: Non-Blocking Incremental Online Vacuum Governor
- **Title:** Non-Blocking Incremental Online Vacuum Governor
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Execute SQLite `PRAGMA incremental_vacuum` in bounded micro-steps across physical shards during idle windows without locking write transactions or spiking CPU usage.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Vacuum telemetry bar showing freed freelist pages, reclaimed disk megabytes, lock contention wait times (`< 1ms`), and governor throttle state in Blue and Green.
  - *Web UI:* Zero-CDN disk savings progress chart and freelist compaction timeline.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard vacuum --shard=<id> [--pages=500] [--max-duration=5s] [--json]`
  - Flags: `--shard` (string, target shard ID), `--pages` (int, max freelist pages to reclaim per cycle), `--max-duration` (duration, cycle timeout), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/vacuum_types.go`):
    ```go
    type IncrementalVacuumRecord struct {
        RunID            string `json:"runId"`
        ShardID          string `json:"shardId"`
        PagesReclaimed   int    `json:"pagesReclaimed"`
        BytesFreed       int64  `json:"bytesFreed"`
        DurationMicros   int64  `json:"durationMicros"`
        TimestampUTC     int64  `json:"timestampUtc"`
        isThrottleActive bool   `json:"isThrottleActive"`
        hasLockContention bool  `json:"hasLockContention"`
        isRunComplete    bool   `json:"isRunComplete"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationVacuumRecords (
        RunID TEXT PRIMARY KEY,
        ShardID TEXT NOT NULL,
        PagesReclaimed INTEGER NOT NULL,
        BytesFreed INTEGER NOT NULL,
        DurationMicros INTEGER NOT NULL,
        TimestampUTC INTEGER NOT NULL,
        isThrottleActive INTEGER NOT NULL CHECK (isThrottleActive IN (0, 1)),
        hasLockContention INTEGER NOT NULL CHECK (hasLockContention IN (0, 1)),
        isRunComplete INTEGER NOT NULL CHECK (isRunComplete IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Scheduled by `pkg/db/sharding/vacuum_compact.go` during background maintenance cycles.

### Improvement 687: Sharded Cross-Node Query Router & Scatter-Gather Engine
- **Title:** Sharded Cross-Node Query Router & Scatter-Gather Engine
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Route analytical SQL queries across multiple physical SQLite database shards concurrently, aggregating results via scatter-gather pipelines with strict memory budgets.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Query execution plan tree showing fanout shard targets, per-shard execution latencies, scatter-gather merge timings, and memory footprint in Cyan and Emerald.
  - *Web UI:* Real-time SVG query orchestration timeline displaying parallel shard query execution.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard query --sql=<query> [--shards=all] [--concurrency=4] [--json]`
  - Flags: `--sql` (string, read-only SQL statement), `--shards` (string, target shards `all` or comma-separated list), `--concurrency` (int, parallel workers), `--json` (bool, output aggregated rows in JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/router_types.go`):
    ```go
    type ScatterGatherQueryPlan struct {
        QueryID          string `json:"queryId"`
        RawSQL           string `json:"rawSql"`
        TargetShardCount int    `json:"targetShardCount"`
        TotalRowsMerged  int64  `json:"totalRowsMerged"`
        ExecutionTimeMs  int64  `json:"executionTimeMs"`
        isScatterSuccess bool   `json:"isScatterSuccess"`
        hasMemoryCapHit  bool   `json:"hasMemoryCapHit"`
        isReadOnlyQuery  bool   `json:"isReadOnlyQuery"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapQueryRouterPlans (
        QueryID TEXT PRIMARY KEY,
        RawSQL TEXT NOT NULL,
        TargetShardCount INTEGER NOT NULL,
        TotalRowsMerged INTEGER NOT NULL,
        ExecutionTimeMs INTEGER NOT NULL,
        CreatedAt INTEGER NOT NULL,
        isScatterSuccess INTEGER NOT NULL CHECK (isScatterSuccess IN (0, 1)),
        hasMemoryCapHit INTEGER NOT NULL CHECK (hasMemoryCapHit IN (0, 1)),
        isReadOnlyQuery INTEGER NOT NULL CHECK (isReadOnlyQuery IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/db/sharding/query_router.go` for analytical warmup aggregations.

### Improvement 688: Online Shard Split, Merge & Dynamic Key Migrator
- **Title:** Online Shard Split, Merge & Dynamic Key Migrator
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Split oversized SQLite shard databases into twin partitions or merge underutilized shards with zero application downtime using copy-on-write page streaming.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Migration status table displaying source shard, target split shards, key migration percentage (`74.2%`), ETA countdown, and dual-write state in Violet and Emerald.
  - *Web UI:* Animated data flow diagram showing shard boundary resizing and live key movement.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard split --source=<id> --threshold-mb=500 [--dry-run] [--json]`
  - Flags: `--source` (string, shard ID to split), `--threshold-mb` (int, size trigger for splitting), `--dry-run` (bool, simulate split boundaries), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/migration_types.go`):
    ```go
    type ShardMigrationTask struct {
        TaskID           string `json:"taskId"`
        SourceShardID    string `json:"sourceShardId"`
        TargetShardID    string `json:"targetShardId"`
        KeysMigrated     int64  `json:"keysMigrated"`
        TotalKeys        int64  `json:"totalKeys"`
        DualWriteEpoch   int64  `json:"dualWriteEpoch"`
        isMigrationActive bool  `json:"isMigrationActive"`
        hasDualWriteLock  bool  `json:"hasDualWriteLock"`
        isCompletedClean  bool  `json:"isCompletedClean"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationShardMigrations (
        TaskID TEXT PRIMARY KEY,
        SourceShardID TEXT NOT NULL,
        TargetShardID TEXT NOT NULL,
        KeysMigrated INTEGER NOT NULL,
        TotalKeys INTEGER NOT NULL,
        DualWriteEpoch INTEGER NOT NULL,
        StartedAt INTEGER NOT NULL,
        isMigrationActive INTEGER NOT NULL CHECK (isMigrationActive IN (0, 1)),
        hasDualWriteLock INTEGER NOT NULL CHECK (hasDualWriteLock IN (0, 1)),
        isCompletedClean INTEGER NOT NULL CHECK (isCompletedClean IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Coordinated by `pkg/db/sharding/shard_rebalancer.go` during cluster maintenance operations.

### Improvement 689: Split-DB Cross-Tier Snapshot & Point-In-Time Backup Engine
- **Title:** Split-DB Cross-Tier Snapshot & Point-In-Time Backup Engine
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Create atomic point-in-time snapshots across the three database tiers (`gitmap.db`, `installation.db`, and `repodb/pipeline.db`) simultaneously using SQLite backup API hooks.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Backup docket table showing snapshot ID, tier verification hashes, compressed archive byte size, and backup duration in Emerald tokens.
  - *Web UI:* Snapshot catalog view with one-click point-in-time restoration preview and cryptographic verification hashes.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard snapshot --out-dir=<dir> [--compress=zstd] [--verify-hash] [--json]`
  - Flags: `--out-dir` (string, snapshot target directory), `--compress` (string, compression algorithm `zstd` or `none`), `--verify-hash` (bool, assert SHA-256 integrity), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/snapshot_types.go`):
    ```go
    type SplitDBSnapshotManifest struct {
        SnapshotID       string `json:"snapshotId"`
        GitmapHash       string `json:"gitmapHash"`
        InstallationHash string `json:"installationHash"`
        PipelineHash     string `json:"pipelineHash"`
        ArchiveSizeBytes int64  `json:"archiveSizeBytes"`
        DurationMillis   int64  `json:"durationMillis"`
        CreatedAt        int64  `json:"createdAt"`
        isVerifiedValid  bool   `json:"isVerifiedValid"`
        hasZeroErrors    bool   `json:"hasZeroErrors"`
        isAtomicSnapshot bool   `json:"isAtomicSnapshot"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationSplitDBSnapshots (
        SnapshotID TEXT PRIMARY KEY,
        GitmapHash TEXT NOT NULL,
        InstallationHash TEXT NOT NULL,
        PipelineHash TEXT NOT NULL,
        ArchiveSizeBytes INTEGER NOT NULL,
        DurationMillis INTEGER NOT NULL,
        CreatedAt INTEGER NOT NULL,
        isVerifiedValid INTEGER NOT NULL CHECK (isVerifiedValid IN (0, 1)),
        hasZeroErrors INTEGER NOT NULL CHECK (hasZeroErrors IN (0, 1)),
        isAtomicSnapshot INTEGER NOT NULL CHECK (isAtomicSnapshot IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/db/sharding/snapshot_engine.go` and integrated into GitMap automated disaster recovery.

### Improvement 690: Interactive Sharded Cluster Cockpit & Terminal Commander
- **Title:** Interactive Sharded Cluster Cockpit & Terminal Commander
- **Cluster:** Cluster 69: Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver
- **Goal:** Deliver an interactive terminal control deck for cluster operators to inspect physical shard nodes, watch live WAL replication streams, trigger rebalancing, and inspect Merkle trees.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Full-screen Bubble Tea dashboard divided into Shard Topology Ring (Top Left), WAL Stream Lag Meters (Top Right), and Consensus Lease & Conflict Resolver Stream (Bottom) with responsive terminal resizing.
  - *Web UI:* Zero-CDN responsive cluster console with live WebSocket replication status and node topology map.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard cockpit [--refresh=1s] [--export-html=<path>]`
  - Flags: `--refresh` (duration, UI refresh cycle), `--export-html` (string, optional static HTML snapshot report destination).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/sharding/cockpit_types.go`):
    ```go
    type ShardClusterCockpitModel struct {
        ActiveShardCount  int    `json:"activeShardCount"`
        TotalReplicas     int    `json:"totalReplicas"`
        MaxLagFrames      int64  `json:"maxLagFrames"`
        LeaderNodeAddress string `json:"leaderNodeAddress"`
        ClusterStatus     string `json:"clusterStatus"`
        isClusterHealthy  bool   `json:"isClusterHealthy"`
        hasPendingVacuum  bool   `json:"hasPendingVacuum"`
        isQuorumAvailable bool   `json:"isQuorumAvailable"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationCockpitSnapshots (
        SnapshotID TEXT PRIMARY KEY,
        ActiveShardCount INTEGER NOT NULL,
        TotalReplicas INTEGER NOT NULL,
        MaxLagFrames INTEGER NOT NULL,
        LeaderNodeAddress TEXT NOT NULL,
        ClusterStatus TEXT NOT NULL,
        RecordedAt INTEGER NOT NULL,
        isClusterHealthy INTEGER NOT NULL CHECK (isClusterHealthy IN (0, 1)),
        hasPendingVacuum INTEGER NOT NULL CHECK (hasPendingVacuum IN (0, 1)),
        isQuorumAvailable INTEGER NOT NULL CHECK (isQuorumAvailable IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Primary cluster management entry point in `cli/cmdcluster/shard_manage_cmd.go` executing `pkg/db/sharding/cluster_cockpit.go`.

---

## 5. Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander (Items 691–700)

### Improvement 691: Dynamic Multi-Agent Swarm Gossip Mesh & Peer Heartbeat
- **Title:** Dynamic Multi-Agent Swarm Gossip Mesh & Peer Heartbeat
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Enable autonomous Antigravity worker agents to discover peers across processes, broadcast heartbeats, exchange capability manifests, and detect ungraceful agent drops via epidemic gossip.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Agent roster table with role badges (`Architect`, `SpecAuthor`, `CodeRefactor`), process IDs, heartbeat round-trip latencies, and status lights in Emerald (`Active`), Amber (`Busy`), and Red (`Dropped`).
  - *Web UI:* Real-time SVG network topology graph visualizing active agent node connections and gossip packet broadcasts.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm peer list [--timeout=2s] [--active-only] [--json]`
  - Flags: `--timeout` (duration, peer discovery probe timeout), `--active-only` (bool, omit disconnected agents), `--json` (bool, JSON formatted peer list).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/peer_types.go`):
    ```go
    type AgentPeerHeartbeat struct {
        AgentID          string `json:"agentId"`
        RoleName         string `json:"roleName"`
        ProcessID        int    `json:"processId"`
        EndpointAddress  string `json:"endpointAddress"`
        LastHeartbeatUTC int64  `json:"lastHeartbeatUtc"`
        AssignedCluster  string `json:"assignedCluster"`
        isAgentAlive     bool   `json:"isAgentAlive"`
        hasActiveTask    bool   `json:"hasActiveTask"`
        isMeshConnected  bool   `json:"isMeshConnected"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapAgentPeerHeartbeats (
        AgentID TEXT PRIMARY KEY,
        RoleName TEXT NOT NULL,
        ProcessID INTEGER NOT NULL,
        EndpointAddress TEXT NOT NULL,
        LastHeartbeatUTC INTEGER NOT NULL,
        AssignedCluster TEXT NOT NULL,
        isAgentAlive INTEGER NOT NULL CHECK (isAgentAlive IN (0, 1)),
        hasActiveTask INTEGER NOT NULL CHECK (hasActiveTask IN (0, 1)),
        isMeshConnected INTEGER NOT NULL CHECK (isMeshConnected IN (0, 1))
    );
    CREATE INDEX IF NOT EXISTS idx_agent_peer_alive ON GitmapAgentPeerHeartbeats(isAgentAlive);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/agent/swarm/peer_mesh.go` and integrated with AGY inter-process messaging.

### Improvement 692: Subtask Exclusivity Lease Arbiter & File Mutex Registry
- **Title:** Subtask Exclusivity Lease Arbiter & File Mutex Registry
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Prevent concurrent agent write collisions by granting atomic, time-bounded leases on specific subtask markdown files and owned Go source files with auto-expiration on agent failure.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* File lock matrix displaying locked file paths, leaseholder agent IDs, lease expiration timers, and queue wait times in Amber and Cyan.
  - *Web UI:* Interactive zero-CDN resource tree showing live file locks, waiting queue depth, and lock release countdowns.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm lease acquire --file=<path> --agent=<id> [--duration=15m] [--json]`
  - Flags: `--file` (string, target file path relative to repo), `--agent` (string, acquiring agent ID), `--duration` (duration, lease TTL), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/lease_types.go`):
    ```go
    type SubtaskFileLease struct {
        FilePath         string `json:"filePath"`
        HolderAgentID    string `json:"holderAgentId"`
        LeaseEpoch       int64  `json:"leaseEpoch"`
        AcquiredAt       int64  `json:"acquiredAt"`
        ExpiresAt        int64  `json:"expiresAt"`
        isLeaseActive    bool   `json:"isLeaseActive"`
        hasExclusiveLock bool   `json:"hasExclusiveLock"`
        isStealable      bool   `json:"isStealable"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSubtaskFileLeases (
        FilePath TEXT PRIMARY KEY,
        HolderAgentID TEXT NOT NULL,
        LeaseEpoch INTEGER NOT NULL,
        AcquiredAt INTEGER NOT NULL,
        ExpiresAt INTEGER NOT NULL,
        isLeaseActive INTEGER NOT NULL CHECK (isLeaseActive IN (0, 1)),
        hasExclusiveLock INTEGER NOT NULL CHECK (hasExclusiveLock IN (0, 1)),
        isStealable INTEGER NOT NULL CHECK (isStealable IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Handled by `pkg/agent/swarm/task_arbiter.go` to enforce single-writer isolation during multi-agent workflows.

### Improvement 693: Real-Time Token Expenditure & LLM Cost Governor
- **Title:** Real-Time Token Expenditure & LLM Cost Governor
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Track input/output token counts and calculated dollar costs across all subagent LLM calls, enforcing strict hard-cap budgets per session and preventing runaway model loops.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Cost budget meter showing prompt tokens, completion tokens, dollar spend (`$0.42 / $5.00 limit`), cost velocity (`$/min`), and burn rate sparklines in Green, Yellow, and Red.
  - *Web UI:* Real-time SVG token burn chart visualizing cumulative expenditure against assigned hard-caps.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm budget status [--session-id=<id>] [--set-limit=5.00] [--json]`
  - Flags: `--session-id` (string, specific swarm session), `--set-limit` (float, update maximum budget in USD), `--json` (bool, JSON formatted output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/cost_types.go`):
    ```go
    type SwarmTokenBudget struct {
        SessionID        string  `json:"sessionId"`
        PromptTokens     int64   `json:"promptTokens"`
        CompletionTokens int64   `json:"completionTokens"`
        TotalCostUSD     float64 `json:"totalCostUsd"`
        BudgetLimitUSD   float64 `json:"budgetLimitUsd"`
        LastUpdatedUTC   int64   `json:"lastUpdatedUtc"`
        isWithinBudget   bool    `json:"isWithinBudget"`
        hasBudgetCapFired bool   `json:"hasBudgetCapFired"`
        isCostRateNormal bool    `json:"isCostRateNormal"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSwarmTokenBudgets (
        SessionID TEXT PRIMARY KEY,
        PromptTokens INTEGER NOT NULL,
        CompletionTokens INTEGER NOT NULL,
        TotalCostUSD REAL NOT NULL,
        BudgetLimitUSD REAL NOT NULL,
        LastUpdatedUTC INTEGER NOT NULL,
        isWithinBudget INTEGER NOT NULL CHECK (isWithinBudget IN (0, 1)),
        hasBudgetCapFired INTEGER NOT NULL CHECK (hasBudgetCapFired IN (0, 1)),
        isCostRateNormal INTEGER NOT NULL CHECK (isCostRateNormal IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Governed by `pkg/agent/swarm/budget_guard.go` with automatic subagent pause when spending limits are reached.

### Improvement 694: Circular Dependency & Swarm Deadlock Detector
- **Title:** Circular Dependency & Swarm Deadlock Detector
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Construct a directed wait-for graph of agent task dependencies and message queues to automatically identify cyclic deadlocks or idle agent starvation and trigger automated unblocking.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Deadlock alert panel displaying the circular dependency cycle path (`Agent A -> Task 2 -> Agent B -> Task 1 -> Agent A`) in Bright Red, followed by automated resolution recommendations.
  - *Web UI:* Interactive SVG directed acyclic graph (DAG) highlighting cyclic deadlock edges with flashing red links.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm deadlock check [--auto-break] [--json]`
  - Flags: `--auto-break` (bool, automatically revoke longest-held lease to break cycles), `--json` (bool, JSON formatted report).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/deadlock_types.go`):
    ```go
    type SwarmDeadlockCycle struct {
        CycleID          string `json:"cycleId"`
        InvolvedAgents   string `json:"involvedAgents"`
        WaitGraphCycle   string `json:"waitGraphCycle"`
        DetectedAt       int64  `json:"detectedAt"`
        ResolutionAction string `json:"resolutionAction"`
        isDeadlockActive bool   `json:"isDeadlockActive"`
        hasCycleResolved bool   `json:"hasCycleResolved"`
        isAutoBreakSafe  bool   `json:"isAutoBreakSafe"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSwarmDeadlocks (
        CycleID TEXT PRIMARY KEY,
        InvolvedAgents TEXT NOT NULL,
        WaitGraphCycle TEXT NOT NULL,
        DetectedAt INTEGER NOT NULL,
        ResolutionAction TEXT NOT NULL,
        isDeadlockActive INTEGER NOT NULL CHECK (isDeadlockActive IN (0, 1)),
        hasCycleResolved INTEGER NOT NULL CHECK (hasCycleResolved IN (0, 1)),
        isAutoBreakSafe INTEGER NOT NULL CHECK (isAutoBreakSafe IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Maintained by `pkg/agent/swarm/deadlock_watcher.go` run as a continuous background daemon during swarm runs.

### Improvement 695: Mechanical Code Quality & Guideline AST Verifier
- **Title:** Mechanical Code Quality & Guideline AST Verifier
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Verify that all Go source files modified by subagents adhere strictly to architectural rules (max 15 lines per function, zero nested `if` statements, positive booleans only, and zero `os.Exit` calls in leaf packages) using native Go AST analysis.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* AST verification scorecard displaying checked files, violation line numbers, rule categories, and green pass checkmarks (`AST VERIFIED: 100% PASS`).
  - *Web UI:* Zero-CDN side-by-side AST linting matrix with interactive AST subtree visualization.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm ast-verify --dir=<path> [--strict] [--json]`
  - Flags: `--dir` (string, directory to analyze), `--strict` (bool, fail on any warning), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/verifier_types.go`):
    ```go
    type AstVerificationReport struct {
        ReportID         string `json:"reportId"`
        TargetFilePath   string `json:"targetFilePath"`
        ViolationsCount  int    `json:"violationsCount"`
        MaxFunctionLines int    `json:"maxFunctionLines"`
        NestedIfCount    int    `json:"nestedIfCount"`
        NegativeBoolCount int   `json:"negativeBoolCount"`
        isAstCompliant   bool   `json:"isAstCompliant"`
        hasZeroLeafExits bool   `json:"hasZeroLeafExits"`
        isReadyForCommit bool   `json:"isReadyForCommit"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapAstVerifications (
        ReportID TEXT PRIMARY KEY,
        TargetFilePath TEXT NOT NULL,
        ViolationsCount INTEGER NOT NULL,
        MaxFunctionLines INTEGER NOT NULL,
        NestedIfCount INTEGER NOT NULL,
        NegativeBoolCount INTEGER NOT NULL,
        VerifiedAt INTEGER NOT NULL,
        isAstCompliant INTEGER NOT NULL CHECK (isAstCompliant IN (0, 1)),
        hasZeroLeafExits INTEGER NOT NULL CHECK (hasZeroLeafExits IN (0, 1)),
        isReadyForCommit INTEGER NOT NULL CHECK (isReadyForCommit IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/agent/swarm/ast_verifier.go` and triggered prior to committing any subagent work.

### Improvement 696: SQLite FTS5 Conversation Transcript AST Indexer & Semantic Search
- **Title:** SQLite FTS5 Conversation Transcript AST Indexer & Semantic Search
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Ingest and index subagent `transcript.jsonl` files into an embedded SQLite FTS5 full-text search table with BM25 ranking to allow instant retrieval of past technical decisions, tool calls, and error traces.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Search results list with BM25 relevance score badges, highlighted text snippets in Bright Yellow, and clickable file/line anchors.
  - *Web UI:* Fast interactive search console with query-as-you-type filtering, tool invocation faceted filters, and step timeline links.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm transcript search --query=<terms> [--limit=10] [--json]`
  - Flags: `--query` (string, full-text search phrase), `--limit` (int, maximum results count), `--json` (bool, JSON formatted search hits).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/transcript_types.go`):
    ```go
    type TranscriptSearchHit struct {
        HitID            string  `json:"hitId"`
        StepIndex        int     `json:"stepIndex"`
        AgentRole        string  `json:"agentRole"`
        SnippetText      string  `json:"snippetText"`
        BM25Score        float64 `json:"bm25Score"`
        TimestampUTC     int64   `json:"timestampUtc"`
        isToolCall       bool    `json:"isToolCall"`
        hasErrorTrace    bool    `json:"hasErrorTrace"`
        isUserVisible    bool    `json:"isUserVisible"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE VIRTUAL TABLE IF NOT EXISTS GitmapTranscriptFTS USING fts5(
        HitID UNINDEXED,
        StepIndex UNINDEXED,
        AgentRole,
        ContentText,
        ToolCallPayload,
        tokenize = 'porter unicode61'
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `pkg/agent/swarm/transcript_indexer.go` for lightning-fast memory lookups.

### Improvement 697: Inter-Agent Event Bus & Streaming Broadcast Hub
- **Title:** Inter-Agent Event Bus & Streaming Broadcast Hub
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Provide a high-throughput, low-latency publish-subscribe event bus over Unix domain sockets or named pipes enabling subagents to broadcast task progress, spec completions, and file modifications.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Live event stream monitor with colorful topic badges (`spec:created`, `ast:passed`, `lease:acquired`), publisher agent names, and message payloads in Slate and Blue.
  - *Web UI:* Streaming event log viewer with real-time WebSocket updates and topic-based filtering.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm events tail [--topic=<pattern>] [--max-events=50] [--json]`
  - Flags: `--topic` (string, event topic wildcard pattern), `--max-events` (int, stream termination count), `--json` (bool, JSON event stream).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/event_bus_types.go`):
    ```go
    type SwarmBusEvent struct {
        EventID          string `json:"eventId"`
        TopicName        string `json:"topicName"`
        PublisherAgentID string `json:"publisherAgentId"`
        PayloadBytes     int    `json:"payloadBytes"`
        PublishedAt      int64  `json:"publishedAt"`
        isBroadcast      bool   `json:"isBroadcast"`
        hasSubscribers   bool   `json:"hasSubscribers"`
        isDelivered      bool   `json:"isDelivered"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSwarmBusEvents (
        EventID TEXT PRIMARY KEY,
        TopicName TEXT NOT NULL,
        PublisherAgentID TEXT NOT NULL,
        PayloadBytes INTEGER NOT NULL,
        PublishedAt INTEGER NOT NULL,
        isBroadcast INTEGER NOT NULL CHECK (isBroadcast IN (0, 1)),
        hasSubscribers INTEGER NOT NULL CHECK (hasSubscribers IN (0, 1)),
        isDelivered INTEGER NOT NULL CHECK (isDelivered IN (0, 1))
    );
    CREATE INDEX IF NOT EXISTS idx_bus_events_topic ON GitmapSwarmBusEvents(TopicName);
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/agent/swarm/event_bus.go` as the messaging backbone of the swarm.

### Improvement 698: Autonomous Subagent Respawn & Self-Healing Sentinel
- **Title:** Autonomous Subagent Respawn & Self-Healing Sentinel
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Detect failed or stalled subagents, preserve their uncommitted work in a quarantine scratchpad, release abandoned file leases, and spawn replacement workers to resume execution automatically.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Self-healing recovery notification card showing failed worker PID, recovered file leases, new replacement agent conversation ID, and recovery latency in Emerald.
  - *Web UI:* Swarm vitality timeline indicating worker failures and automated self-healing events.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm heal [--agent=<id>] [--force-respawn] [--json]`
  - Flags: `--agent` (string, target agent ID to inspect/heal), `--force-respawn` (bool, terminate and re-launch worker), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/healing_types.go`):
    ```go
    type SwarmSelfHealingEvent struct {
        HealingID        string `json:"healingId"`
        FailedAgentID    string `json:"failedAgentId"`
        ReplacementID    string `json:"replacementId"`
        ReleasedLeases   int    `json:"releasedLeases"`
        RecoveredAt      int64  `json:"recoveredAt"`
        isRespawnSuccess bool   `json:"isRespawnSuccess"`
        hasStatePreserved bool  `json:"hasStatePreserved"`
        isHealedCleanly  bool   `json:"isHealedCleanly"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSwarmHealingEvents (
        HealingID TEXT PRIMARY KEY,
        FailedAgentID TEXT NOT NULL,
        ReplacementID TEXT NOT NULL,
        ReleasedLeases INTEGER NOT NULL,
        RecoveredAt INTEGER NOT NULL,
        isRespawnSuccess INTEGER NOT NULL CHECK (isRespawnSuccess IN (0, 1)),
        hasStatePreserved INTEGER NOT NULL CHECK (hasStatePreserved IN (0, 1)),
        isHealedCleanly INTEGER NOT NULL CHECK (isHealedCleanly IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed by `pkg/agent/swarm/self_healer.go` to provide zero-human-intervention resilience.

### Improvement 699: Cross-Subagent Task Ledger & Dependency Graph Orchestrator
- **Title:** Cross-Subagent Task Ledger & Dependency Graph Orchestrator
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Coordinate macro plans across multiple subagents using an acyclic dependency ledger, ensuring prerequisites (specs, schemas) are fully satisfied and verified before dependent code implementation begins.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Multi-tier task roadmap showing prerequisite dependencies, task readiness states (`BLOCKED`, `READY`, `IN_PROGRESS`, `DONE`), and assignee agents in Cyan, Amber, and Green.
  - *Web UI:* Interactive DAG diagram visualizing prerequisite task flows, real-time completion percentages, and critical path bottlenecks.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm ledger status [--plan=<file>] [--show-blocked] [--json]`
  - Flags: `--plan` (string, path to master plan markdown), `--show-blocked` (bool, display unready tasks with unmet dependencies), `--json` (bool, JSON formatted output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/ledger_types.go`):
    ```go
    type SwarmLedgerTask struct {
        TaskSlug         string `json:"taskSlug"`
        SequenceNumber   int    `json:"sequenceNumber"`
        RequiredPrereq   string `json:"requiredPrereq"`
        AssigneeAgentID  string `json:"assigneeAgentId"`
        TaskStatus       string `json:"taskStatus"`
        UpdatedAt        int64  `json:"updatedAt"`
        isPrereqMet      bool   `json:"isPrereqMet"`
        hasActiveWorker  bool   `json:"hasActiveWorker"`
        isTaskComplete   bool   `json:"isTaskComplete"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSwarmLedgerTasks (
        TaskSlug TEXT PRIMARY KEY,
        SequenceNumber INTEGER NOT NULL,
        RequiredPrereq TEXT NOT NULL,
        AssigneeAgentID TEXT NOT NULL,
        TaskStatus TEXT NOT NULL,
        UpdatedAt INTEGER NOT NULL,
        isPrereqMet INTEGER NOT NULL CHECK (isPrereqMet IN (0, 1)),
        hasActiveWorker INTEGER NOT NULL CHECK (hasActiveWorker IN (0, 1)),
        isTaskComplete INTEGER NOT NULL CHECK (isTaskComplete IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Maintained by `pkg/agent/swarm/ledger_orchestrator.go` to govern multi-phase sprint workflows.

### Improvement 700: Real-Time Antigravity Swarm Mission Commander HUD
- **Title:** Real-Time Antigravity Swarm Mission Commander HUD
- **Cluster:** Cluster 70: Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander
- **Goal:** Deliver an all-encompassing, high-refresh Bubble Tea terminal cockpit giving the human engineer complete operational awareness over active subagents, token expenditures, file leases, AST compliance, and live transcript events.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Multi-pane master mission deck:
    - *Header:* Sequence title, overall swarm progress bar (`100% DONE`), active subagent count, total token spend.
    - *Upper Left:* Active Agent Peer Roster & Heartbeat Health (Cyan borders).
    - *Upper Right:* Token Burn Velocity & Dollar Hard-Cap Meter (Emerald/Amber borders).
    - *Lower Left:* File Lease Mutex Matrix & Deadlock Monitor (Slate/Violet borders).
    - *Lower Right:* Live Streaming Event Bus & AST Quality Scorecard (Blue/Green borders).
  - *Web UI:* Zero-CDN responsive mission dashboard with interactive widgets, live WebSocket telemetry, and emergency pause/resume controls.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm commander [--refresh=500ms] [--export-html=<path>]`
  - Flags: `--refresh` (duration, HUD redraw cycle), `--export-html` (string, optional static HTML mission briefing export path).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/commander_hud_types.go`):
    ```go
    type SwarmMissionCommanderSnapshot struct {
        SnapshotID       string  `json:"snapshotId"`
        ActiveAgents     int     `json:"activeAgents"`
        TotalTokensSpent int64   `json:"totalTokensSpent"`
        TotalCostUSD     float64 `json:"totalCostUsd"`
        LockedFilesCount int     `json:"lockedFilesCount"`
        AstPassRate      float64 `json:"astPassRate"`
        RecordedAt       int64   `json:"recordedAt"`
        isSwarmHealthy   bool    `json:"isSwarmHealthy"`
        hasDeadlockAlert bool    `json:"hasDeadlockAlert"`
        isMissionActive  bool    `json:"isMissionActive"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSwarmCommanderSnapshots (
        SnapshotID TEXT PRIMARY KEY,
        ActiveAgents INTEGER NOT NULL,
        TotalTokensSpent INTEGER NOT NULL,
        TotalCostUSD REAL NOT NULL,
        LockedFilesCount INTEGER NOT NULL,
        AstPassRate REAL NOT NULL,
        RecordedAt INTEGER NOT NULL,
        isSwarmHealthy INTEGER NOT NULL CHECK (isSwarmHealthy IN (0, 1)),
        hasDeadlockAlert INTEGER NOT NULL CHECK (hasDeadlockAlert IN (0, 1)),
        isMissionActive INTEGER NOT NULL CHECK (isMissionActive IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Flagship CLI command registered in `cli/cmdswarm/commander_cmd.go` executing `pkg/agent/swarm/commander_hud.go`.

---
