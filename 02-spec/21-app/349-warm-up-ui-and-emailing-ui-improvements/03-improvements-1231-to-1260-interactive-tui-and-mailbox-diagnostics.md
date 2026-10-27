# Specification 349: Improvements 1231 to 1260 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 349
- **Slug:** `349-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/349-warm-up-ui-and-emailing-ui-improvements/03-improvements-1231-to-1260-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/348-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/347-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (1231 - 1260)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 124** | Items 1231–1240 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 125** | Items 1241–1250 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 126** | Items 1251–1260 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 124: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 1231–1240)

---

### Improvement 1231: Distributed Seed Mailbox Geo-Latency Topology & Anycast Routing Matrix
- **ID:** `SEED-TOPO-1231`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Visualize and diagnose round-trip edge network transit latencies across globally distributed warmup seed nodes (US-East, US-West, EU-Central, AP-Southeast) targeting consumer ESP edge endpoints with dynamic BGP Anycast ping tracking.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-column topological view showing regional POP latency, route jitter, and packet degradation across active seed cohorts.
    - Lipgloss border styling (`#38bdf8`, sky blue) with status badges: `[OPTIMAL]` (emerald `#10b981`), `[ELEVATED]` (amber `#f59e0b`), `[CRITICAL]` (rose `#f43f5e`).
    - Model view snippet:
      ```text
      ┌─ [DISTRIBUTED SEED GEO-LATENCY TOPOLOGY & ANYCAST ROUTING] ──────────────────────────┐
      │ GLOBAL SEED CLUSTERS: 4 REGIONS | MONITORED NODES: 480 | MEAN ANYCAST RTT: 28.4ms     │
      ├───────────────────────┬────────────┬───────────┬─────────────┬─────────────────────────┤
      │ REGION / EDGE POP     │ ACTIVE IPs │ PKT LOSS  │ JITTER (ms) │ MEAN TRANSIT RTT        │
      ├───────────────────────┼────────────┼───────────┼─────────────┼─────────────────────────┤
      │ us-east-1 (N. Virginia)│        160 │     0.00% │      0.82ms │ 14.2ms [●] OPTIMAL      │
      │ us-west-2 (Oregon)    │        120 │     0.02% │      1.14ms │ 22.8ms [●] OPTIMAL      │
      │ eu-central-1 (Frankfurt)│       120 │     0.00% │      0.95ms │ 34.6ms [●] OPTIMAL      │
      │ ap-southeast-1 (Singapore)│      80 │     0.11% │      3.40ms │ 89.2ms [▲] ELEVATED     │
      ├───────────────────────┴────────────┴───────────┴─────────────┴─────────────────────────┤
      │ [CURSOR: ap-southeast-1 :: node-sg-04.seed.gitmap.internal]                            │
      │ ├─ Target MX: aspmx.l.google.com (Anycast IPv6: 2a00:1450:400c:c0b::1b)                │
      │ ├─ TCP SYN Handshake: 42.1ms | TLS 1.3 ClientHello -> Finished: 47.1ms (Total: 89.2ms) │
      │ ├─ Path Deviation: +18.4ms above regional median (Peering Hop: AS15169 SingTel Node)   │
      │ └─ Keys: [r: Retest Ping] [m: Force Alternate POP] [c: View Trace] [q: Quit]          │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plot mapping regional seed cluster nodes without Google Maps or Mapbox dependencies. Inline SVG tooltips.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-topology`
  - Flags:
    - `--region <string>`: Filter seed cohort by target cloud region (`us-east-1`, `us-west-2`, `eu-central-1`, `ap-southeast-1`, `all`; default: `all`).
    - `--max-acceptable-rtt <duration>`: Threshold to trigger route re-evaluation alerts (default: `75ms`).
    - `--is-streaming`: Continuous Bubble Tea telemetry stream (default: `true`).
    - `--has-traceroute`: Execute automatic ICMP/TCP path tracing upon elevated latency (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedGeoTopologyNode struct {
        NodeID             string        `json:"nodeId"`
        RegionCode         string        `json:"regionCode"`
        EdgePopEndpoint    string        `json:"edgePopEndpoint"`
        TargetMxHost       string        `json:"targetMxHost"`
        MeanTransitRtt     time.Duration `json:"meanTransitRtt"`
        PacketLossRate     float64       `json:"packetLossRate"`
        JitterVarianceMs   float64       `json:"jitterVarianceMs"`
        LastProbedAt       time.Time     `json:"lastProbedAt"`
        IsOptimal          bool          `json:"isOptimal"`
        HasDirectPeering   bool          `json:"hasDirectPeering"`
        HasTraceAlert      bool          `json:"hasTraceAlert"`
    }

    type SeedGeoTopologyModel struct {
        Nodes          []SeedGeoTopologyNode `json:"nodes"`
        CursorIndex    int                   `json:"cursorIndex"`
        SelectedNode   SeedGeoTopologyNode   `json:"selectedNode"`
        GlobalMeanRtt  time.Duration         `json:"globalMeanRtt"`
        IsStreaming    bool                  `json:"isStreaming"`
        HasAlerts      bool                  `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedGeoTopologyRecord (
        node_id TEXT PRIMARY KEY,
        region_code TEXT NOT NULL,
        edge_pop_endpoint TEXT NOT NULL,
        target_mx_host TEXT NOT NULL,
        mean_transit_rtt_ms REAL NOT NULL DEFAULT 0.0,
        packet_loss_rate REAL NOT NULL DEFAULT 0.0,
        jitter_variance_ms REAL NOT NULL DEFAULT 0.0,
        last_probed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_direct_peering INTEGER NOT NULL DEFAULT 1,
        has_trace_alert INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_seed_geo_region ON SeedGeoTopologyRecord(region_code, is_optimal);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/seedpool/topology.go` and surfaces interactive Bubble Tea matrix via `gitmap email warmup seed-topology`.

---

### Improvement 1232: Seed Mailbox Auto-Authentication & Automated TLS Keypin Rotator
- **ID:** `SEED-KEYPIN-1232`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Manage autonomous TLS certificate pinning, private key rotation, and OAuth2 refresh synchronization across high-security seed mailboxes, proactively warning before key expiration and ensuring zero unauthenticated drops.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Cryptographic key vault dashboard with real-time countdown clocks, SHA-256 fingerprint verifications, and one-key rollover triggering:
      ```text
      ┌─ [SEED MAILBOX TLS KEYPIN & AUTHENTICATION VAULT] ────────────────────────────────────┐
      │ MONITORED VAULT ITEMS: 340 SEEDS | EXPIRED: 0 | ROTATING: 3 | STRICT PINNING: ENFORCED│
      ├───────────────────────┬────────────┬─────────────────────────────┬─────────────────────┤
      │ SEED ACCOUNT          │ AUTH PROTO │ CERT PIN / PUBKEY HASH      │ EXPIRATION / STATUS │
      ├───────────────────────┼────────────┼─────────────────────────────┼─────────────────────┤
      │ vault-01@corp-mail.io │ mTLS X.509 │ sha256:d8a1...9b42 (RSA4096)│ 42d left [●] ACTIVE │
      │ seed-39@proton.custom │ PGP Subkey │ openpgp:4f71...a108 (ED25519)│ 180d left[●] ACTIVE │
      │ oauth-admin@gsuite.cc │ OAuth2 OIDC│ token_exp: 34m (Auto-Ref)   │ 34m left [▲] ROTATE │
      │ sec-agent@o365-ent.org│ Azure Cert │ sha256:c119...ee01 (ECDSA)  │ 3d left  [■] URGENT │
      ├───────────────────────┴────────────┴─────────────────────────────┴─────────────────────┤
      │ [ACTIONS: Press 'r' to trigger instant token refresh, 'g' to generate mTLS keypair]    │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Active green `#22c55e`, Rotating amber `#f59e0b`, Urgent rose `#f43f5e`, Key violet `#c084fc`.
  - **Zero-CDN Web Dashboard:**
    - Local vault management panel with SHA-256 fingerprint visual checksum block and zero external fonts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-auth-rotator`
  - Flags:
    - `--days-before-expiry <int>`: Alert threshold for expiring certificates (default: `14`).
    - `--is-auto-renew`: Automatically request OAuth2/mTLS token renewal when remaining life < 20% (default: `true`).
    - `--auth-provider <string>`: Filter by auth mechanism (`mtls`, `oauth2`, `pgp`, `all`; default: `all`).
    - `--has-strict-pinning`: Reject TLS sessions if public key pinning mismatch is detected (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedKeypinItem struct {
        VaultID           string    `json:"vaultId"`
        SeedEmailAddress  string    `json:"seedEmailAddress"`
        AuthProtocol      string    `json:"authProtocol"`
        KeyFingerprint    string    `json:"keyFingerprint"`
        ExpiresAt         time.Time `json:"expiresAt"`
        LastRefreshedAt   time.Time `json:"lastRefreshedAt"`
        IsValid           bool      `json:"isValid"`
        IsAutoRenewable   bool      `json:"isAutoRenewable"`
        HasStrictPinning  bool      `json:"hasStrictPinning"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedKeypinVaultRecord (
        vault_id TEXT PRIMARY KEY,
        seed_email_address TEXT NOT NULL UNIQUE,
        auth_protocol TEXT NOT NULL,
        key_fingerprint TEXT NOT NULL,
        expires_at DATETIME NOT NULL,
        last_refreshed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_valid INTEGER NOT NULL DEFAULT 1,
        is_auto_renewable INTEGER NOT NULL DEFAULT 1,
        has_strict_pinning INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_seed_keypin_expiry ON SeedKeypinVaultRecord(expires_at, is_valid);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/seedpool/auth_vault.go` and integrates with GitMap secrets manager.

---

### Improvement 1233: Seed Inbound Header Forensic Timeline & Authentication Tracer HUD
- **ID:** `SEED-FORENSIC-1233`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide an interactive step-by-step forensic timeline decomposing received message headers (`Received: from ... by ...`, `Authentication-Results`, `ARC-Seal`) in seed mailboxes to verify ESP hop authenticity and trace transport security degradations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Timeline view detailing each relay hop with TLS cipher negotiation, hop duration, and validation badges.
    - Model view snippet:
      ```text
      ┌─ [SEED INBOUND HEADER FORENSIC TRACER: MSG-ID <wup-8921@out.internal>] ──────────────┐
      │ HOP 1/3 (0.12s) │ out-01.relay.internal -> smtp-in.google.com [TLS 1.3 / X25519]       │
      │ HOP 2/3 (0.45s) │ mail-sorter-04.google.com -> mx-cluster-b.google.com [INTERNAL]      │
      │ HOP 3/3 (0.08s) │ Delivery to Mailbox: seed-gsuite-19@enterprise.net [INBOX]           │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [AUTHENTICATION SUMMARY MATRIX]                                                      │
      │ ├─ SPF:   PASS (ip=198.51.100.12, client-ip=198.51.100.12, helo=mail.domain.io)      │
      │ ├─ DKIM:  PASS (header.i=@domain.io header.s=s1 header.b=8f29ea10 RSA-2048)          │
      │ ├─ DMARC: PASS (p=reject sp=reject pct=100 dis=none)                                 │
      │ ├─ ARC:   PASS (i=1 as.b=0b31 spf=pass dkim=pass dmarc=pass)                         │
      │ └─ TLS:   STRICT (TLS_AES_256_GCM_SHA384 / ECDHE / Perfect Forward Secrecy)          │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Pass green `#22c55e`, Hop cyan `#06b6d4`, Timeline line slate `#475569`.
  - **Zero-CDN Web Dashboard:**
    - Pure CSS timeline diagram with collapsible header inspection accordions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup header-forensics`
  - Flags:
    - `--msg-id <string>`: Specific warmup Message-ID to trace.
    - `--is-deep-arc`: Validate full ARC header chain integrity across all hops (default: `true`).
    - `--has-raw-headers`: Display untouched RFC 5322 header payload (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HeaderHopTrace struct {
        HopIndex       int           `json:"hopIndex"`
        OriginRelay    string        `json:"originRelay"`
        RecipientRelay string        `json:"recipientRelay"`
        TlsCipherSuite string        `json:"tlsCipherSuite"`
        HopDuration    time.Duration `json:"hopDuration"`
        IsTlsSecured   bool          `json:"isTlsSecured"`
        HasArcVerified bool          `json:"hasArcVerified"`
    }

    type HeaderForensicResult struct {
        MessageID       string           `json:"messageId"`
        SeedAddress     string           `json:"seedAddress"`
        TotalTransit    time.Duration    `json:"totalTransit"`
        Hops            []HeaderHopTrace `json:"hops"`
        IsSpfPass       bool             `json:"isSpfPass"`
        IsDkimPass      bool             `json:"isDkimPass"`
        IsDmarcPass     bool             `json:"isDmarcPass"`
        HasStrictTls    bool             `json:"hasStrictTls"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HeaderForensicRecord (
        message_id TEXT PRIMARY KEY,
        seed_address TEXT NOT NULL,
        total_transit_ms REAL NOT NULL,
        is_spf_pass INTEGER NOT NULL DEFAULT 1,
        is_dkim_pass INTEGER NOT NULL DEFAULT 1,
        is_dmarc_pass INTEGER NOT NULL DEFAULT 1,
        has_strict_tls INTEGER NOT NULL DEFAULT 1,
        recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_forensic_seed ON HeaderForensicRecord(seed_address, recorded_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/seedpool/forensics.go` for real-time validation of delivery authentication.

---

### Improvement 1234: Automated Multi-ESP Seed Mailbox Persona Emulation Engine
- **ID:** `SEED-PERSONA-1234`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Emulate diverse human email client behaviors (open pauses, reading scroll delays, folder movements, marking as important, replying with generated context) across distinct persona profiles (executive, developer, consumer) to defeat ISP heuristic spam filters.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live persona execution monitor displaying assigned behavior archetypes, pending actions, and realism confidence scores:
      ```text
      ┌─ [SEED PERSONA EMULATION ENGINE: 120 CONCURRENT HUMAN REPLICAS] ──────────────────────┐
      │ ACTIVE PERSONAS: 120 | MEAN ENGAGEMENT REALISM SCORE: 98.4% | EMULATING MUA: THUNDERBIRD│
      ├───────────────────────┬────────────┬─────────────────────────────┬─────────────────────┤
      │ SEED MAILBOX          │ PERSONA    │ CURRENT EMULATION ACTION    │ READ DELAY / PACE   │
      ├───────────────────────┼────────────┼─────────────────────────────┼─────────────────────┤
      │ cto-seed@alpha.tech   │ Executive  │ Moving to 'Projects' Folder │ 4.8s reading dwell  │
      │ dev-ops@gamma.io      │ Developer  │ Replying to Thread #481     │ 14.2s draft pause   │
      │ consumer-91@yahoo.com │ Shopper    │ Marking as 'Starred / VIP'  │ 2.1s dwell pause    │
      │ analyst-02@corp.fin   │ Financial  │ Downloading Attached Report │ 7.4s inspection     │
      ├───────────────────────┴────────────┴─────────────────────────────┴─────────────────────┤
      │ [STATUS: All personas operating with randomized Gaussian jitter (mean=5.2s, σ=1.4s)]   │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Persona magenta `#d946ef`, Action cyan `#06b6d4`, Dwell amber `#fbbf24`.
  - **Zero-CDN Web Dashboard:**
    - Persona behavior profile editor and action dispatch log viewer.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup persona-emulate`
  - Flags:
    - `--persona-type <string>`: Filter or assign archetype (`executive`, `developer`, `consumer`, `random`; default: `random`).
    - `--min-dwell-sec <float>`: Minimum simulated message reading dwell time (default: `2.5`).
    - `--is-reply-enabled`: Allow automated contextual reply generation (default: `true`).
    - `--has-gaussian-jitter`: Apply statistical Gaussian variance to all action timings (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedPersonaState struct {
        PersonaID         string        `json:"personaId"`
        SeedAddress       string        `json:"seedAddress"`
        Archetype         string        `json:"archetype"`
        LastAction        string        `json:"lastAction"`
        DwellDuration     time.Duration `json:"dwellDuration"`
        RealismScore      float64       `json:"realismScore"`
        LastEngagedAt     time.Time     `json:"lastEngagedAt"`
        IsReplyEnabled    bool          `json:"isReplyEnabled"`
        HasGaussianJitter bool          `json:"hasGaussianJitter"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedPersonaStateRecord (
        persona_id TEXT PRIMARY KEY,
        seed_address TEXT NOT NULL UNIQUE,
        archetype TEXT NOT NULL,
        last_action TEXT NOT NULL,
        dwell_duration_ms INTEGER NOT NULL,
        realism_score REAL NOT NULL DEFAULT 1.0,
        last_engaged_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_reply_enabled INTEGER NOT NULL DEFAULT 1,
        has_gaussian_jitter INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_persona_archetype ON SeedPersonaStateRecord(archetype, is_reply_enabled);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/seedpool/persona.go` providing engagement behaviors.

---

### Improvement 1235: Seed Mailbox IMAP Auto-Folder Sanitizer & Archival Vacuum HUD
- **ID:** `SEED-VACUUM-1235`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Routinely prune, sanitize, and vacuum seed mailbox folders (Trash, Junk, Archive) over IMAP/POP3, maintaining mailbox responsiveness, preventing quota starvation, and keeping mail accounts indistinguishable from active users.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual storage reclaim progress bar and folder cleanup log:
      ```text
      ┌─ [SEED MAILBOX IMAP FOLDER SANITIZER & ARCHIVAL VACUUM] ─────────────────────────────┐
      │ VACUUM RUN: ACTIVE | TARGETED MAILBOXES: 320 | RECLAIMED STORAGE: 14.8 GB / 22.0 GB   │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [CURRENT PROGRESS]                                                                   │
      │ Overall Cleanup:  [████████████████████░░░░░░░░░░] 67.2% (215/320 Accounts Cleaned)   │
      │ Junk Purged:      12,490 messages | Expired Trash Expunged: 8,910 messages           │
      │ Folder Balancing: Inbox maintain target: 45 messages (Active Thread Preservation)    │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [ACTIVE THREAD: seed-yahoo-44@yahoo.com]                                             │
      │ ├─ IMAP Command: EXPUNGE (Junk folder: 142 items cleared, 18.2 MB reclaimed)        │
      │ └─ Response:     OK EXPUNGE completed (0.42s)                                       │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Storage blue `#3b82f6`, Progress cyan `#06b6d4`, Success green `#10b981`.
  - **Zero-CDN Web Dashboard:**
    - Storage breakdown charts rendered via inline HTML5 canvas-free flex-bars.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-vacuum`
  - Flags:
    - `--retention-days <int>`: Max age for non-thread warmup messages in seed boxes (default: `14`).
    - `--is-dry-run`: Preview expunge candidates without deleting (default: `false`).
    - `--has-trash-purge`: Empty Trash/Bin folders after inbox rotation (default: `true`).
    - `--target-inbox-size <int>`: Retain top N most recent messages in inbox (default: `50`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedVacuumReport struct {
        ReportID          string    `json:"reportId"`
        SeedAddress       string    `json:"seedAddress"`
        MessagesExpunged  int       `json:"messagesExpunged"`
        BytesReclaimed    int64     `json:"bytesReclaimed"`
        RemainingInboxMsg int       `json:"remainingInboxMsg"`
        VacuumedAt        time.Time `json:"vacuumedAt"`
        IsClean           bool      `json:"isClean"`
        HasTrashPurged    bool      `json:"hasTrashPurged"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedVacuumExecutionRecord (
        report_id TEXT PRIMARY KEY,
        seed_address TEXT NOT NULL,
        messages_expunged INTEGER NOT NULL DEFAULT 0,
        bytes_reclaimed INTEGER NOT NULL DEFAULT 0,
        remaining_inbox_msg INTEGER NOT NULL DEFAULT 0,
        vacuumed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_clean INTEGER NOT NULL DEFAULT 1,
        has_trash_purged INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_vacuum_seed ON SeedVacuumExecutionRecord(seed_address, vacuumed_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/seedpool/vacuum.go` and executed as a GitMap scheduled task.

---

### Improvement 1236: Seed Cluster Cross-Provider Inbox Affinity & Scoring HUD
- **ID:** `SEED-AFFINITY-1236`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Compute multi-dimensional affinity and trust scores between specific outbound IP/domain combinations and seed mailbox providers, dynamically directing test volumes toward lower-scoring providers to build missing trust signals.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Affinity matrix table mapping sending domains against major recipient providers with dynamic trust scoring (0–100):
      ```text
      ┌─ [CROSS-PROVIDER INBOX AFFINITY & REPUTATION SCOREBOARD] ────────────────────────────┐
      │ DOMAIN: dispatch.marketflow.io | AGGREGATE TRUST: 88.4/100 | REPUTATION TIER: TIER-1  │
      ├───────────────────────┬────────────┬─────────────┬─────────────┬─────────────────────┤
      │ DESTINATION ESP       │ TRUST (0-100)│ INBOX SCORE │ REPLY RATIO │ RECOMMENDED ACTION  │
      ├───────────────────────┼────────────┼─────────────┼─────────────┼─────────────────────┤
      │ Google Workspace      │       94.2 │       99.4% │       34.1% │ Maintain Cruise Pacing│
      │ Microsoft 365         │       88.1 │       96.2% │       28.4% │ Maintain Cruise Pacing│
      │ Yahoo / AOL           │       72.4 │       88.0% │       19.0% │ Boost Seed Affinity │
      │ Apple iCloud Mail     │       68.9 │       84.5% │       12.5% │ Throttle & Seed Resc│
      ├───────────────────────┴────────────┴─────────────┴─────────────┴─────────────────────┤
      │ [DISPATCH POLICY: Weight increased +25% on Yahoo/Apple seeds to rebuild reputation]  │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Trust emerald `#10b981`, Boost amber `#f59e0b`, Throttle rose `#f43f5e`.
  - **Zero-CDN Web Dashboard:**
    - Hexagonal trust radar chart rendered via inline SVG paths.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup provider-affinity`
  - Flags:
    - `--domain <string>`: Filter affinity by sending domain.
    - `--min-trust-threshold <float>`: Flag providers falling below trust score (default: `75.0`).
    - `--is-auto-rebalance`: Automatically divert higher seed volume to low-affinity providers (default: `true`).
    - `--has-detailed-factors`: Display individual scoring factor breakdowns (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ProviderAffinityScore struct {
        ScoreID           string    `json:"scoreId"`
        SendingDomain     string    `json:"sendingDomain"`
        ProviderType      string    `json:"providerType"`
        TrustScore        float64   `json:"trustScore"`
        InboxPlacement    float64   `json:"inboxPlacement"`
        EngagementRatio   float64   `json:"engagementRatio"`
        LastEvaluatedAt   time.Time `json:"lastEvaluatedAt"`
        IsOptimal         bool      `json:"isOptimal"`
        HasVolumeBoost    bool      `json:"hasVolumeBoost"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ProviderAffinityScoreRecord (
        score_id TEXT PRIMARY KEY,
        sending_domain TEXT NOT NULL,
        provider_type TEXT NOT NULL,
        trust_score REAL NOT NULL,
        inbox_placement REAL NOT NULL,
        engagement_ratio REAL NOT NULL,
        last_evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_volume_boost INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_affinity_domain ON ProviderAffinityScoreRecord(sending_domain, provider_type);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/seedpool/affinity.go`, driving intelligent seed target selection.

---

### Improvement 1237: Real-Time Seed Mailbox Quarantine Alert & Instant Webhook Dispatcher
- **ID:** `SEED-ALERT-1237`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Detect sudden mailbox quarantine triggers (e.g. 3 consecutive spam placements or unexpected bounce rejections) and fire sub-millisecond local Unix-socket or HTTP webhooks to suspend sending instantly.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-visibility emergency breaker banner with audio-visual terminal flash and instant pause button:
      ```text
      ┌─ [EMERGENCY QUARANTINE CIRCUIT BREAKER ALERT] ───────────────────────────────────────┐
      │ STATUS: QUARANTINE TRIPPED [CRITICAL] | AFFECTED EGRESS: 198.51.100.19 (Cluster B)   │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [INCIDENT TELEMETRY]                                                                 │
      │ ├─ Trigger: 4 consecutive Spam folder arrivals detected on Microsoft 365 seeds       │
      │ ├─ Impact:  Domain: reachout-mail.net | Egress IP: 198.51.100.19                     │
      │ ├─ Action:  CIRCUIT OPENED -> All outbound queues for reachout-mail.net PAUSED       │
      │ └─ Webhook: Dispatched to unix:///var/run/gitmap-alert.sock (0.84ms latency)         │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [Press 'o' to override and resume, 'd' to run diagnostic scan, 'q' to acknowledge]   │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Alert red `#dc2626`, Circuit yellow `#facc15`, Text white `#f8fafc`.
  - **Zero-CDN Web Dashboard:**
    - Live toast notification system powered by native Server-Sent Events (SSE).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup quarantine-monitor`
  - Flags:
    - `--consecutive-spam-threshold <int>`: Consecutive spam landings before trip (default: `3`).
    - `--webhook-endpoint <string>`: Target webhook URL or Unix socket path.
    - `--is-auto-freeze`: Automatically freeze sending queue on trip (default: `true`).
    - `--has-sound-alert`: Emit terminal ASCII bell on trip (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type QuarantineIncident struct {
        IncidentID       string    `json:"incidentId"`
        SendingDomain    string    `json:"sendingDomain"`
        EgressIP         string    `json:"egressIp"`
        ConsecutiveSpam  int       `json:"consecutiveSpam"`
        TrippedAt        time.Time `json:"trippedAt"`
        IsCircuitOpen    bool      `json:"isCircuitOpen"`
        IsWebhookSent    bool      `json:"isWebhookSent"`
        HasHumanOverride bool      `json:"hasHumanOverride"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS QuarantineIncidentRecord (
        incident_id TEXT PRIMARY KEY,
        sending_domain TEXT NOT NULL,
        egress_ip TEXT NOT NULL,
        consecutive_spam INTEGER NOT NULL,
        tripped_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_circuit_open INTEGER NOT NULL DEFAULT 1,
        is_webhook_sent INTEGER NOT NULL DEFAULT 1,
        has_human_override INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_quarantine_domain ON QuarantineIncidentRecord(sending_domain, tripped_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/warmup/seedpool/quarantine.go` hooked into the core sending pipeline.

---

### Improvement 1238: Seed Mailbox IMAP Search UID Synchronization & Cache Accel HUD
- **ID:** `SEED-UIDSYNC-1238`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Maintain persistent local SQLite UID maps of seed inboxes, using IMAP `UID SEARCH SINCE` incremental synchronization to minimize bandwidth usage and achieve sub-10ms arrival scans across thousands of mailboxes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Fast UID synchronization ticker displaying cached vs fetched message headers:
      ```text
      ┌─ [SEED IMAP UID CACHE & INCREMENTAL SYNC ENGINE] ────────────────────────────────────┐
      │ SYNCHRONIZED SEEDS: 420 | LOCAL UID MAPS: CACHED | BANDWIDTH SAVED TODAY: 94.2%      │
      ├───────────────────────┬────────────┬─────────────┬─────────────┬─────────────────────┤
      │ SEED MAILBOX          │ CACHED UID │ REMOTE UID  │ DELTA MSGS  │ SYNC DURATION       │
      ├───────────────────────┼────────────┼─────────────┼─────────────┼─────────────────────┤
      │ seed-01@o365-corp.org │      14920 │       14922 │      +2 new │ 8.4ms (Fast Sync)   │
      │ seed-18@fastmail.com  │       4101 │        4101 │           0 │ 3.2ms (Zero Delta)  │
      │ seed-82@gmail-acc.net │       9820 │        9824 │      +4 new │ 12.1ms (Fast Sync)  │
      ├───────────────────────┴────────────┴─────────────┴─────────────┴─────────────────────┤
      │ [CACHE HIT RATE: 99.1% | Zero full-folder downloads required]                        │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Sync cyan `#06b6d4`, Delta green `#10b981`, Zero slate `#64748b`.
  - **Zero-CDN Web Dashboard:**
    - Cache hit/miss ratio gauge and incremental bandwidth metrics table.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup uid-sync-status`
  - Flags:
    - `--sync-interval <duration>`: Periodic sync cycle (default: `30s`).
    - `--is-fast-cache`: Enforce local SQLite UID cache queries before network requests (default: `true`).
    - `--has-full-resync`: Force complete UID recalculation and validation (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedUidSyncState struct {
        SyncID         string        `json:"syncId"`
        SeedAddress    string        `json:"seedAddress"`
        CachedMaxUID   uint32        `json:"cachedMaxUid"`
        RemoteMaxUID   uint32        `json:"remoteMaxUid"`
        NewMessageDelta int          `json:"newMessageDelta"`
        SyncDuration   time.Duration `json:"syncDuration"`
        LastSyncedAt   time.Time     `json:"lastSyncedAt"`
        IsSynchronized bool          `json:"isSynchronized"`
        HasCacheHit    bool          `json:"hasCacheHit"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedUidSyncRecord (
        sync_id TEXT PRIMARY KEY,
        seed_address TEXT NOT NULL UNIQUE,
        cached_max_uid INTEGER NOT NULL,
        remote_max_uid INTEGER NOT NULL,
        new_message_delta INTEGER NOT NULL DEFAULT 0,
        sync_duration_ms REAL NOT NULL,
        last_synced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_synchronized INTEGER NOT NULL DEFAULT 1,
        has_cache_hit INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_uidsync_seed ON SeedUidSyncRecord(seed_address, is_synchronized);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/seedpool/uidsync.go`, accelerating IMAP polling routines.

---

### Improvement 1239: Seed Mailbox Multi-Factor Authentication (MFA) App-Password Health HUD
- **ID:** `SEED-MFA-1239`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Monitor and audit app-specific password authorizations and session validity across 2FA-enforced seed accounts (Google, Yahoo, iCloud), alerting before security policy revocations cause batch disconnects.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - App password status dashboard with real-time auth verification and last-success markers:
      ```text
      ┌─ [SEED MFA APP-PASSWORD HEALTH & SESSION VALIDATION] ────────────────────────────────┐
      │ TOTAL 2FA ACCOUNTS: 280 | HEALTHY: 278 | REVOKED/EXPIRED: 2 | LAST VERIFIED: 1m ago  │
      ├───────────────────────┬────────────┬─────────────────────────────┬─────────────────────┤
      │ SEED ACCOUNT          │ PROVIDER   │ LAST AUTH SUCCESS           │ HEALTH STATUS       │
      ├───────────────────────┼────────────┼─────────────────────────────┼─────────────────────┤
      │ dev-seed-01@gmail.com │ Google     │ 14:32:10 UTC (2m ago)       │ [●] VALID / ACTIVE  │
      │ app-seed-14@yahoo.com │ Yahoo      │ 14:30:00 UTC (4m ago)       │ [●] VALID / ACTIVE  │
      │ icloud-seed@icloud.com│ Apple      │ 12:10:00 UTC (2.5h ago)     │ [▲] AUTH_INVALID    │
      ├───────────────────────┴────────────┴─────────────────────────────┴─────────────────────┤
      │ [WARNING: icloud-seed@icloud.com rejected APP-PASSWORD. Re-issuance required]        │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Active green `#22c55e`, Invalid red `#ef4444`, Provider cyan `#06b6d4`.
  - **Zero-CDN Web Dashboard:**
    - App-password health matrix with one-click credential re-verification triggers.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mfa-health`
  - Flags:
    - `--provider <string>`: Filter by provider (`google`, `yahoo`, `apple`, `all`; default: `all`).
    - `--is-auto-test`: Run non-intrusive IMAP login check during scan (default: `true`).
    - `--has-revocation-alert`: Send alert notification on auth rejection (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedMfaHealthStatus struct {
        StatusID         string    `json:"statusId"`
        SeedAddress      string    `json:"seedAddress"`
        ProviderType     string    `json:"providerType"`
        LastAuthSuccess  time.Time `json:"lastAuthSuccess"`
        FailedAttempts   int       `json:"failedAttempts"`
        IsValid          bool      `json:"isValid"`
        IsRevoked        bool      `json:"isRevoked"`
        HasAlertTriggered bool     `json:"hasAlertTriggered"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedMfaHealthRecord (
        status_id TEXT PRIMARY KEY,
        seed_address TEXT NOT NULL UNIQUE,
        provider_type TEXT NOT NULL,
        last_auth_success DATETIME NOT NULL,
        failed_attempts INTEGER NOT NULL DEFAULT 0,
        is_valid INTEGER NOT NULL DEFAULT 1,
        is_revoked INTEGER NOT NULL DEFAULT 0,
        has_alert_triggered INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_mfa_provider ON SeedMfaHealthRecord(provider_type, is_valid);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/warmup/seedpool/mfa.go` with automated status notifications.

---

### Improvement 1240: Seed Inbound Spam-Score Reverse Engineering & Heuristic Inspector HUD
- **ID:** `SEED-SPAMREVERSE-1240`
- **Cluster:** 124 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Extract and parse vendor-specific spam score headers (`X-Spam-Status`, `X-Microsoft-Antispam`, `X-Forefront-Antispam-Report`, `X-Google-DKIM-Signature`) from arriving warmup test messages to reverse-engineer spam engine decisions in real time.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-section spam breakdown HUD showing extracted BCL (Bulk Complaint Level), SCL (Spam Confidence Level), and SpamAssassin rule hits:
      ```text
      ┌─ [SEED SPAM-SCORE REVERSE ENGINEERING & HEURISTIC INSPECTOR] ────────────────────────┐
      │ TARGET MSG-ID: <wup-9041@mail.org> | SEED: m365-seed@enterprise.cc | RESULT: INBOX    │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [MICROSOFT 365 ANTISPAM REPORT]                                                      │
      │ ├─ SCL (Spam Confidence Level): 1 (Safe, Range: 0-9)                                  │
      │ ├─ BCL (Bulk Complaint Level):  2 (Low Bulk Risk, Range: 0-9)                         │
      │ ├─ CAT (Category):              NONE (Not marked as bulk, spam, or phishing)         │
      │ └─ PCL (Phish Confidence Level):0 (Clean)                                             │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [SPAMASSASSIN HEURISTIC RULE MATCHES (Total Score: -1.2 / 5.0)]                      │
      │ ├─ [-1.0] DKIM_SIGNED         Message has a valid DKIM signature                     │
      │ ├─ [-0.5] SPF_PASS            Sender host matches SPF record                         │
      │ └─ [+0.3] HTML_MESSAGE        HTML included in payload                               │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Safe green `#22c55e`, Rule cyan `#06b6d4`, Score yellow `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Interactive spam factor breakdown tree with score sensitivity simulator.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup spam-reverse-inspect`
  - Flags:
    - `--msg-id <string>`: Inspect headers for specific message.
    - `--is-strict-scl`: Flag messages with SCL > 3 as high risk (default: `true`).
    - `--has-rule-explanations`: Output comprehensive text descriptions for every rule code (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamReverseReport struct {
        ReportID       string    `json:"reportId"`
        MessageID      string    `json:"messageId"`
        SeedAddress    string    `json:"seedAddress"`
        SclScore       int       `json:"sclScore"`
        BclScore       int       `json:"bclScore"`
        SpamAssassinTotal float64 `json:"spamAssassinTotal"`
        EvaluatedAt    time.Time `json:"evaluatedAt"`
        IsInboxSafe    bool      `json:"isInboxSafe"`
        HasHeuristicFlag bool    `json:"hasHeuristicFlag"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamReverseReportRecord (
        report_id TEXT PRIMARY KEY,
        message_id TEXT NOT NULL,
        seed_address TEXT NOT NULL,
        scl_score INTEGER NOT NULL DEFAULT 0,
        bcl_score INTEGER NOT NULL DEFAULT 0,
        spam_assassin_total REAL NOT NULL DEFAULT 0.0,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_inbox_safe INTEGER NOT NULL DEFAULT 1,
        has_heuristic_flag INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_spam_reverse_msg ON SpamReverseReportRecord(message_id, is_inbox_safe);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/seedpool/spam_reverse.go` supplying deliverability diagnostics.


# Cluster 125: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 1241–1250)

---

### Improvement 1241: Token-Bucket Dynamic Egress Pacer & Burst Dispersion HUD
- **ID:** `THROTTLE-BURST-1241`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Govern outbound SMTP email delivery bursts using a dual token-bucket rate limiter that enforces strict micro-cadence dispersion (spacing individual messages by calculated millisecond offsets) to prevent triggering ISP burst defenses.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live token bucket fluid level meter and microsecond burst dispersion gauge:
      ```text
      ┌─ [TOKEN-BUCKET DYNAMIC EGRESS PACER & BURST DISPERSION] ─────────────────────────────┐
      │ DESTINATION: mx.google.com | CAPACITY: 60 TOKENS | REFILL RATE: 2.5 token/sec        │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [TOKEN BUCKET FLUID GAUGE]                                                           │
      │ Current Level:    [████████████████████░░░░░░░░░░] 42 / 60 Tokens (Equilibrium: OK)  │
      │ Inter-Message Gap: 412ms programmed offset | Jitter: ±35ms (Uniform Dispersion)      │
      │ Dropped/Throttled: 0 messages | Smoothing Efficiency: 99.8%                         │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [MICRO-CADENCE DISPERSION TRAIN]                                                     │
      │ [MSG #1041] ──── 415ms ────► [MSG #1042] ──── 398ms ────► [MSG #1043] ──── 421ms ──► │
      │ Egress Interface: eth0 (198.51.100.10) | Socket TCP Buffer: OPTIMAL (32KB write)    │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Token green `#10b981`, Cadence sky `#0284c7`, Warning amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Animated SVG token bucket tank level visualizer and cadence distribution histogram.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup burst-pacer`
  - Flags:
    - `--capacity <int>`: Token bucket maximum burst capacity (default: `60`).
    - `--refill-rate <float>`: Tokens replenished per second (default: `2.0`).
    - `--min-gap-ms <int>`: Minimum millisecond gap between consecutive SMTP handshakes (default: `350`).
    - `--is-jitter-enabled`: Add randomized timing jitter to defeat periodic heuristic filters (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TokenBucketState struct {
        BucketID        string        `json:"bucketId"`
        TargetDomain    string        `json:"targetDomain"`
        CurrentTokens   float64       `json:"currentTokens"`
        MaxCapacity     int           `json:"maxCapacity"`
        RefillRatePerSec float64      `json:"refillRatePerSec"`
        InterMessageGap time.Duration `json:"interMessageGap"`
        LastRefilledAt  time.Time     `json:"lastRefilledAt"`
        IsBalanced      bool          `json:"isBalanced"`
        HasJitter       bool          `json:"hasJitter"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TokenBucketStateRecord (
        bucket_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL UNIQUE,
        current_tokens REAL NOT NULL,
        max_capacity INTEGER NOT NULL,
        refill_rate_per_sec REAL NOT NULL,
        inter_message_gap_ms INTEGER NOT NULL,
        last_refilled_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_balanced INTEGER NOT NULL DEFAULT 1,
        has_jitter INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_token_bucket_domain ON TokenBucketStateRecord(target_domain, is_balanced);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/transport/token_bucket.go`, throttling raw socket dispatch routines.

---

### Improvement 1242: Multi-Egress IP Affinity Balancer & Reputation Equalizer HUD
- **ID:** `THROTTLE-EQUALIZE-1242`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Dynamically balance outbound warmup traffic across multiple dedicated egress IP interfaces based on real-time ISP reputation scores, distributing volume away from cold or recovering IPs toward established interfaces.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-IP egress pool matrix with load allocation bars, current stage, and reputation grades:
      ```text
      ┌─ [MULTI-EGRESS IP AFFINITY BALANCER & REPUTATION EQUALIZER] ─────────────────────────┐
      │ POOL: Production Outbound Pool A | TOTAL IPs: 4 | ACTIVE TRAFFIC ALLOCATION: 100%    │
      ├─────────────────┬──────────┬──────────┬─────────────────────────────┬──────────────────┤
      │ EGRESS IP       │ REPUTATION│ STAGE    │ TRAFFIC WEIGHT ALLOCATION   │ CONCURRENT CONNS │
      ├─────────────────┼──────────┼──────────┼─────────────────────────────┼──────────────────┤
      │ 198.51.100.10   │ 99/100   │ Stage 4  │ [████████████████] 40.0%    │  8 / 12 slots    │
      │ 198.51.100.11   │ 96/100   │ Stage 3  │ [████████████░░░░] 30.0%    │  6 / 10 slots    │
      │ 198.51.100.12   │ 84/100   │ Stage 2  │ [████████░░░░░░░░] 20.0%    │  4 /  8 slots    │
      │ 198.51.100.13   │ 71/100   │ Recovery │ [████░░░░░░░░░░░░] 10.0%    │  2 /  5 slots    │
      ├─────────────────┴──────────┴──────────┴─────────────────────────────┴──────────────────┤
      │ [REBALANCING ENGINE: IP .13 placed on 50% reduced volume due to recent 421 throttle]   │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: HighRep emerald `#10b981`, Recovery amber `#f59e0b`, Progress cyan `#06b6d4`.
  - **Zero-CDN Web Dashboard:**
    - Visual egress IP traffic split diagram with drag-and-drop weight adjustment override.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ip-equalizer`
  - Flags:
    - `--pool-id <string>`: Identifier for IP egress pool (default: `default`).
    - `--is-auto-rebalance`: Dynamically adjust traffic percentage based on bounce telemetry (default: `true`).
    - `--min-recovery-weight <float>`: Minimum traffic percentage allocated to recovering IPs (default: `0.05`).
    - `--has-strict-stage-caps`: Enforce maximum hourly sends dictated by warmup stage (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EgressIpAllocation struct {
        AllocationID       string    `json:"allocationId"`
        EgressIP           string    `json:"egressIp"`
        PoolID             string    `json:"poolId"`
        WarmupStage        string    `json:"warmupStage"`
        ReputationScore    int       `json:"reputationScore"`
        AllocatedWeight    float64   `json:"allocatedWeight"`
        ActiveConnections  int       `json:"activeConnections"`
        MaxConnections     int       `json:"maxConnections"`
        LastRebalancedAt   time.Time `json:"lastRebalancedAt"`
        IsActive           bool      `json:"isActive"`
        HasRecoveryBackoff bool      `json:"hasRecoveryBackoff"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EgressIpAllocationRecord (
        allocation_id TEXT PRIMARY KEY,
        egress_ip TEXT NOT NULL UNIQUE,
        pool_id TEXT NOT NULL,
        warmup_stage TEXT NOT NULL,
        reputation_score INTEGER NOT NULL DEFAULT 100,
        allocated_weight REAL NOT NULL DEFAULT 0.25,
        active_connections INTEGER NOT NULL DEFAULT 0,
        max_connections INTEGER NOT NULL DEFAULT 10,
        last_rebalanced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_active INTEGER NOT NULL DEFAULT 1,
        has_recovery_backoff INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_egress_pool ON EgressIpAllocationRecord(pool_id, is_active);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/transport/equalizer.go`, controlling source IP binding on SMTP dialers.

---

### Improvement 1243: SMTP Destination MX Health Circuit Breaker & Fallback Matrix HUD
- **ID:** `THROTTLE-CIRCUIT-1243`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Monitor destination MX cluster health (detecting consecutive SMTP 5xx permanent failures, persistent connection timeouts, or TLS negotiation failures), immediately tripping local circuit breakers and routing to secondary MX records.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual circuit state matrix (`CLOSED`, `HALF-OPEN`, `OPEN`) per destination MX host:
      ```text
      ┌─ [SMTP DESTINATION MX HEALTH CIRCUIT BREAKER MATRIX] ────────────────────────────────┐
      │ DESTINATION DOMAIN: corp-enterprise.com | MX HOSTS MONITORED: 2 | STATUS: STABLE     │
      ├───────────────────────┬────────────┬─────────────┬─────────────┬─────────────────────┤
      │ MX HOSTNAME           │ PREFERENCE │ STATE       │ ERROR RATE  │ CONSECUTIVE FAILS   │
      ├───────────────────────┼────────────┼─────────────┼─────────────┼─────────────────────┤
      │ mx1.corp-enterprise.cc│         10 │ CLOSED [●]  │       0.12% │ 0 / 5 failures      │
      │ mx2.corp-enterprise.cc│         20 │ OPEN   [▲]  │      42.10% │ 5 / 5 [CIRCUIT TRIP]│
      ├───────────────────────┴────────────┴─────────────┴─────────────┴─────────────────────┤
      │ [BREAKER TELEMETRY: mx2 tripped due to TLS handshake timeout. Retry in 4m 12s]       │
      │ [ROUTING POLICY: 100% of outbound connections directed to primary mx1]              │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Closed green `#10b981`, HalfOpen amber `#f59e0b`, Open red `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Circuit state transition diagram with manual trip/reset controls.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup circuit-status`
  - Flags:
    - `--domain <string>`: Filter circuit breaker status by destination domain.
    - `--failure-threshold <int>`: Consecutive failures before circuit opens (default: `5`).
    - `--cooldown-duration <duration>`: Cooldown before half-open probe trial (default: `5m`).
    - `--is-auto-reset`: Automatically attempt half-open recovery probes (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MxCircuitBreakerState struct {
        CircuitID          string        `json:"circuitId"`
        MxHostname         string        `json:"mxHostname"`
        DestinationDomain  string        `json:"destinationDomain"`
        State              string        `json:"state"`
        ConsecutiveFails   int           `json:"consecutiveFails"`
        FailureThreshold   int           `json:"failureThreshold"`
        CooldownDuration   time.Duration `json:"cooldownDuration"`
        LastStateChange    time.Time     `json:"lastStateChange"`
        IsOpen             bool          `json:"isOpen"`
        IsHalfOpen         bool          `json:"isHalfOpen"`
        HasFallbackRouted  bool          `json:"hasFallbackRouted"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxCircuitBreakerRecord (
        circuit_id TEXT PRIMARY KEY,
        mx_hostname TEXT NOT NULL UNIQUE,
        destination_domain TEXT NOT NULL,
        state TEXT NOT NULL DEFAULT 'CLOSED',
        consecutive_fails INTEGER NOT NULL DEFAULT 0,
        failure_threshold INTEGER NOT NULL DEFAULT 5,
        cooldown_duration_sec INTEGER NOT NULL DEFAULT 300,
        last_state_change DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_open INTEGER NOT NULL DEFAULT 0,
        is_half_open INTEGER NOT NULL DEFAULT 0,
        has_fallback_routed INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_circuit_domain ON MxCircuitBreakerRecord(destination_domain, state);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/warmup/transport/circuit.go`, intercepting outbound connection establishment.

---

### Improvement 1244: Cross-Tenant Fair-Share Priority Queue & Starvation Guard HUD
- **ID:** `THROTTLE-FAIRNESS-1244`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Prevent high-volume warmup campaigns from starving low-volume transactional or test email queues using Deficit Round-Robin (DRR) fair-share scheduling across tenants sharing outbound IP pools.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live multi-tenant queue monitor with quantum deficit balances and round-robin cycle counters:
      ```text
      ┌─ [DEFICIT ROUND-ROBIN FAIR-SHARE MULTI-TENANT QUEUE] ────────────────────────────────┐
      │ ACTIVE TENANTS: 3 | DRR QUANTUM: 1,500 BYTES | SYSTEM DISPATCH CAPACITY: 20 msg/sec   │
      ├───────────────────────┬────────────┬─────────────┬─────────────┬─────────────────────┤
      │ TENANT IDENTIFIER     │ PRIORITY   │ DEFICIT (B) │ QUEUE DEPTH │ CURRENT DISPATCH    │
      ├───────────────────────┼────────────┼─────────────┼─────────────┼─────────────────────┤
      │ Tenant-A (Marketing)  │ Normal     │        420  │ 1,420 msgs  │ Serving Packet...   │
      │ Tenant-B (Warmup Seed)│ High       │      1,250  │   140 msgs  │ Quantum Available   │
      │ Tenant-C (Transact)   │ Urgent     │      3,000  │     2 msgs  │ High Priority Serve │
      ├───────────────────────┴────────────┴─────────────┴─────────────┴─────────────────────┤
      │ [FAIRNESS CHECK: Tenant-A prevented from hogging egress socket; Starvation Guard OK] │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Urgent magenta `#ec4899`, High cyan `#06b6d4`, Normal slate `#64748b`.
  - **Zero-CDN Web Dashboard:**
    - Tenant throughput allocation stream and latency distribution percentiles (p50, p95, p99).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup fairness-queue`
  - Flags:
    - `--quantum-bytes <int>`: DRR quantum slice size in bytes (default: `1500`).
    - `--max-tenant-starvation <duration>`: Maximum acceptable wait time before priority boost (default: `15s`).
    - `--is-starvation-guard`: Automatically promote stalled tenant queues (default: `true`).
    - `--has-urgent-bypass`: Allow transactional messages to preempt bulk warmup queues (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TenantQueueState struct {
        TenantID         string    `json:"tenantId"`
        PriorityTier     string    `json:"priorityTier"`
        DeficitBytes     int       `json:"deficitBytes"`
        QueueDepth       int       `json:"queueDepth"`
        MeanWaitTimeMs   float64   `json:"meanWaitTimeMs"`
        LastDispatchedAt time.Time `json:"lastDispatchedAt"`
        IsStarved        bool      `json:"isStarved"`
        HasUrgentBypass  bool      `json:"hasUrgentBypass"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TenantQueueStateRecord (
        tenant_id TEXT PRIMARY KEY,
        priority_tier TEXT NOT NULL,
        deficit_bytes INTEGER NOT NULL DEFAULT 0,
        queue_depth INTEGER NOT NULL DEFAULT 0,
        mean_wait_time_ms REAL NOT NULL DEFAULT 0.0,
        last_dispatched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_starved INTEGER NOT NULL DEFAULT 0,
        has_urgent_bypass INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_tenant_priority ON TenantQueueStateRecord(priority_tier, is_starved);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/transport/fairness.go`, organizing internal dispatch channels.

---

### Improvement 1245: SMTP 421/451 Transient Greylist Backoff & Exponential Jitter Matrix
- **ID:** `THROTTLE-GREYLIST-1245`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Distinguish destination ISP greylisting delays (e.g. `451 4.7.1 Service unavailable - try again later`) from harsh rate limits, applying an adaptive exponential retry schedule with full jitter to avoid spam traps.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Greylist retry schedule timeline with ISP retry window indicators:
      ```text
      ┌─ [DESTINATION ISP GREYLIST ADAPTIVE RETRY SCHEDULE] ─────────────────────────────────┐
      │ ACTIVE GREYLIST TRACKERS: 18 DESTINATIONS | ESTIMATED GREYLIST WINDOW: 10 - 15 mins   │
      ├───────────────────────┬────────────┬─────────────┬─────────────┬─────────────────────┤
      │ RECIPIENT DOMAIN      │ ATTEMPT    │ FIRST 451   │ RETRY DUE   │ CALCULATED JITTER   │
      ├───────────────────────┼────────────┼─────────────┼─────────────┼─────────────────────┤
      │ corporate-buyer.de    │ 2nd of 4   │ 14:15:10 UTC│ in 4m 12s   │ ±42s (Full Jitter)  │
      │ secure-bank.co.uk     │ 1st of 4   │ 14:28:40 UTC│ in 11m 20s  │ ±1m 15s (Full Jitter│
      │ medical-health.org    │ 3rd of 4   │ 13:58:00 UTC│ in 1m 05s   │ ±18s (Full Jitter)  │
      ├───────────────────────┴────────────┴─────────────┴─────────────┴─────────────────────┤
      │ [STRATEGY: Greylist detected via triplet key (ClientIP, Sender, Recipient)]          │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Greylist amber `#f59e0b`, Due green `#10b981`, Jitter cyan `#06b6d4`.
  - **Zero-CDN Web Dashboard:**
    - Greylist survival probability curve and retry queue scheduler inspector.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup greylist-status`
  - Flags:
    - `--domain <string>`: Check greylist state for recipient domain.
    - `--base-retry-delay <duration>`: Initial retry sleep after first 451 (default: `10m`).
    - `--max-retry-attempts <int>`: Max retry iterations before flagging as hard rejection (default: `5`).
    - `--is-full-jitter`: Apply statistical uniform full jitter to retry delay (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GreylistTracker struct {
        TrackerID        string    `json:"trackerId"`
        RecipientDomain  string    `json:"recipientDomain"`
        SenderAddress    string    `json:"senderAddress"`
        AttemptCount     int       `json:"attemptCount"`
        FirstAttemptAt   time.Time `json:"firstAttemptAt"`
        NextRetryDueAt   time.Time `json:"nextRetryDueAt"`
        IsGreylisted     bool      `json:"isGreylisted"`
        HasFullJitter    bool      `json:"hasFullJitter"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GreylistTrackerRecord (
        tracker_id TEXT PRIMARY KEY,
        recipient_domain TEXT NOT NULL,
        sender_address TEXT NOT NULL,
        attempt_count INTEGER NOT NULL DEFAULT 1,
        first_attempt_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        next_retry_due_at DATETIME NOT NULL,
        is_greylisted INTEGER NOT NULL DEFAULT 1,
        has_full_jitter INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_greylist_domain ON GreylistTrackerRecord(recipient_domain, next_retry_due_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/transport/greylist.go`, governing message reschedule timers.

---

### Improvement 1246: Asynchronous Egress MAC/VLAN Interface Cycling & Pacer HUD
- **ID:** `THROTTLE-VLAN-1246`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Cycle outbound SMTP traffic across virtual network interfaces (VLANs / sub-interfaces with unique MAC and IP bindings) on Linux/BSD host engines to prevent local ARP table bottlenecking and distribute NIC interrupt loading.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live virtual network adapter grid displaying packet rate, CPU core affinity, and NIC RX/TX queues:
      ```text
      ┌─ [ASYNCHRONOUS EGRESS VLAN / MAC INTERFACE CYCLER] ──────────────────────────────────┐
      │ INTERFACES ACTIVE: 4 VLANs | LINK SPEED: 10 Gbps | NIC INTERRUPT SPREAD: BALANCED    │
      ├─────────────────┬──────────┬─────────────────────────────┬─────────────┬───────────────┤
      │ INTERFACE NAME  │ VLAN ID  │ VIRTUAL MAC ADDRESS         │ TX MSG RATE │ CPU PINNING   │
      ├─────────────────┼──────────┼─────────────────────────────┼─────────────┼───────────────┤
      │ eth0.100        │      100 │ 02:42:c0:a8:01:10           │ 18.2 msg/s  │ Core #0, #1   │
      │ eth0.101        │      101 │ 02:42:c0:a8:01:11           │ 19.4 msg/s  │ Core #2, #3   │
      │ eth0.102        │      102 │ 02:42:c0:a8:01:12           │ 17.8 msg/s  │ Core #4, #5   │
      │ eth0.103        │      103 │ 02:42:c0:a8:01:13           │ 18.9 msg/s  │ Core #6, #7   │
      ├─────────────────┴──────────┴─────────────────────────────┴─────────────┴───────────────┤
      │ [NIC STATUS: Zero dropped frames | Hardware TCP Offload (TSO/GSO): ENABLED]          │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Interface sky `#0284c7`, Hardware green `#10b981`, MAC slate `#94a3b8`.
  - **Zero-CDN Web Dashboard:**
    - Hardware NIC packet flow visualizer with throughput heatmaps.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup vlan-cycler`
  - Flags:
    - `--base-interface <string>`: Primary network adapter (default: `eth0`).
    - `--vlan-range <string>`: Range of VLAN IDs to cycle (e.g. `100-103`).
    - `--is-mac-randomized`: Ensure each sub-interface maintains unique locally administered MAC (default: `true`).
    - `--has-cpu-pinning`: Pin worker threads to dedicated NUMA CPU cores (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type VlanInterfaceState struct {
        InterfaceID     string    `json:"interfaceId"`
        InterfaceName   string    `json:"interfaceName"`
        VlanID          int       `json:"vlanId"`
        VirtualMac      string    `json:"virtualMac"`
        TxMsgRate       float64   `json:"txMsgRate"`
        BoundCpuCore    int       `json:"boundCpuCore"`
        LastRotatedAt   time.Time `json:"lastRotatedAt"`
        IsUp            bool      `json:"isUp"`
        HasHardwareTso  bool      `json:"hasHardwareTso"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS VlanInterfaceStateRecord (
        interface_id TEXT PRIMARY KEY,
        interface_name TEXT NOT NULL UNIQUE,
        vlan_id INTEGER NOT NULL,
        virtual_mac TEXT NOT NULL,
        tx_msg_rate REAL NOT NULL DEFAULT 0.0,
        bound_cpu_core INTEGER NOT NULL,
        last_rotated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_up INTEGER NOT NULL DEFAULT 1,
        has_hardware_tso INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_vlan_status ON VlanInterfaceStateRecord(is_up, bound_cpu_core);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/warmup/transport/vlan.go` providing OS network binding hooks.

---

### Improvement 1247: Outbound SMTP Socket Micro-Pacer & Inter-Packet Gap (IPG) HUD
- **ID:** `THROTTLE-MICROPACE-1247`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Control the exact nanosecond/microsecond spacing of outbound SMTP TCP packets (SYN, TLS handshake chunks, DATA blocks), emulating biological user delays and smoothing traffic through Deep Packet Inspection (DPI) edge firewalls.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Oscilloscope-style packet timeline detailing inter-packet gaps:
      ```text
      ┌─ [OUTBOUND SMTP SOCKET MICRO-PACER & IPG CONTROLLER] ────────────────────────────────┐
      │ SOCKET: fd=42 (198.51.100.12 -> 142.250.102.26:25) | PACING ENGINE: NANOSECOND LEAKY  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [TCP PACKET TRAIN TIMELINE]                                                          │
      │ [SYN]   ── 1.2ms ──► [TLS CH] ── 840μs ──► [AUTH] ── 420μs ──► [DATA] ── 2.1ms ──►[QUIT]│
      │ Mean Inter-Packet Gap: 812μs | Target Dispersion: 750μs - 1,200μs (Gaussian)         │
      │ Packet Fragmentation:  Prevented (MTU=1500 / MSS=1460 Optimal)                       │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [HARDWARE INTERFACE STATUS: SO_MAX_PACING_RATE enforced via Linux FQ qdisc]          │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Packet cyan `#06b6d4`, Timing green `#10b981`, Gap amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Live microsecond timing scatter plot rendered with inline SVG dots.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup micro-pacer`
  - Flags:
    - `--target-ipg-micros <int>`: Target inter-packet gap in microseconds (default: `800`).
    - `--is-fq-qdisc`: Enforce Linux Fair Queueing pacing socket options (default: `true`).
    - `--has-gaussian-variance`: Introduce biological microsecond jitter (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MicroPacerTelemetry struct {
        SocketID         string    `json:"socketId"`
        FileDescriptor   int       `json:"fileDescriptor"`
        TargetIpgMicros  int       `json:"targetIpgMicros"`
        ActualMeanMicros float64   `json:"actualMeanMicros"`
        PacketsPaced     int64     `json:"packetsPaced"`
        LastPacedAt      time.Time `json:"lastPacedAt"`
        IsActive         bool      `json:"isActive"`
        HasGaussianJitter bool     `json:"hasGaussianJitter"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MicroPacerTelemetryRecord (
        socket_id TEXT PRIMARY KEY,
        file_descriptor INTEGER NOT NULL,
        target_ipg_micros INTEGER NOT NULL DEFAULT 800,
        actual_mean_micros REAL NOT NULL DEFAULT 0.0,
        packets_paced INTEGER NOT NULL DEFAULT 0,
        last_paced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_active INTEGER NOT NULL DEFAULT 1,
        has_gaussian_jitter INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_pacer_socket ON MicroPacerTelemetryRecord(socket_id, is_active);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/transport/pacer.go` wrapping standard TCP connections.

---

### Improvement 1248: Destination ISP TLS 1.3 0-RTT Early Data Governor HUD
- **ID:** `THROTTLE-ZERORTT-1248`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Safely govern TLS 1.3 0-RTT Early Data session resumptions when reconnecting to modern MX relays (e.g. ProtonMail, Fastmail, Postfix 3.4+), eliminating a full round-trip handshake while enforcing replay attack defenses.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Handshake round-trip speed comparison HUD:
      ```text
      ┌─ [TLS 1.3 0-RTT EARLY DATA SESSION RESUMPTION GOVERNOR] ─────────────────────────────┐
      │ PROTOCOL: TLS 1.3 / X25519 / CHACHA20-POLY1305 | 0-RTT SAVINGS: 28.4ms per connection │
      ├───────────────────────┬────────────┬─────────────────────────────┬─────────────────────┤
      │ DESTINATION MX        │ TLS TICKET │ 0-RTT EARLY DATA STATUS     │ LATENCY BENEFIT     │
      ├───────────────────────┼────────────┼─────────────────────────────┼─────────────────────┤
      │ mail.protonmail.ch    │ RESUMED    │ [●] ACCEPTED (EHLO sent 0ms)│ -32.1ms (1 RTT)     │
      │ in1-smtp.messaging.cc │ RESUMED    │ [●] ACCEPTED (EHLO sent 0ms)│ -24.8ms (1 RTT)     │
      │ aspmx.l.google.com    │ NEW TICKET │ [▲] 1-RTT REQUIRED (Initial)│ Base Transit RTT    │
      ├───────────────────────┴────────────┴─────────────────────────────┴─────────────────────┤
      │ [ANTI-REPLAY DEFENSE: Monotonic ticket nonce tracking prevents duplicate commands]   │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Accepted green `#10b981`, Resumed cyan `#06b6d4`, Required amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Handshake waterfall chart comparing 0-RTT, 1-RTT, and legacy 2-RTT connections.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup tls-early-data`
  - Flags:
    - `--is-0rtt-enabled`: Allow TLS 1.3 Early Data on resumed sessions (default: `true`).
    - `--ticket-cache-size <int>`: Max cached session tickets (default: `1000`).
    - `--has-anti-replay`: Enforce single-use ticket nonce registry (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TlsEarlyDataSession struct {
        SessionID       string    `json:"sessionId"`
        DestinationHost string    `json:"destinationHost"`
        TicketHash      string    `json:"ticketHash"`
        HandshakeRttMs  float64   `json:"handshakeRttMs"`
        RttSavedMs      float64   `json:"rttSavedMs"`
        ExpiresAt       time.Time `json:"expiresAt"`
        IsResumed       bool      `json:"isResumed"`
        HasAntiReplay   bool      `json:"hasAntiReplay"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TlsEarlyDataSessionRecord (
        session_id TEXT PRIMARY KEY,
        destination_host TEXT NOT NULL,
        ticket_hash TEXT NOT NULL UNIQUE,
        handshake_rtt_ms REAL NOT NULL,
        rtt_saved_ms REAL NOT NULL DEFAULT 0.0,
        expires_at DATETIME NOT NULL,
        is_resumed INTEGER NOT NULL DEFAULT 1,
        has_anti_replay INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_tls_ticket ON TlsEarlyDataSessionRecord(destination_host, expires_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/transport/tls_session.go` optimizing cryptographic dials.

---

### Improvement 1249: SMTP Response Code Semantics Parser & ISP Anomaly Detector HUD
- **ID:** `THROTTLE-ANOMALY-1249`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Parse non-standard, cryptic, or vendor-specific SMTP response strings (e.g. Proofpoint block messages, Barracuda rate limits, Mimecast greylisting codes), normalizing them into actionable decision enums for the concurrency governor.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live anomaly stream mapping raw SMTP text to normalized semantic codes:
      ```text
      ┌─ [SMTP RESPONSE CODE SEMANTICS & ISP ANOMALY PARSER] ────────────────────────────────┐
      │ PARSED RESPONSES: 12,410 | ANOMALIES RESOLVED: 18 | NORMALIZATION ACCURACY: 99.9%     │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [RAW RESPONSE]                                                                       │
      │ 452 4.3.1 Insufficient system storage (Proofpoint PPS-49219b: Egress IP Throttled)   │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [SEMANTIC RESOLUTION]                                                                │
      │ ├─ Normalized Enum:    THROTTLE_VENDOR_STORAGE_QUOTA                                 │
      │ ├─ Recognized Vendor:  Proofpoint Enterprise Protection                              │
      │ ├─ Recommended Action: Apply 15-minute egress IP cooldown; do NOT retry immediately  │
      │ └─ Confidence Score:   99.4% (Matched Pattern Rule #PP-452)                          │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Raw yellow `#facc15`, Semantic violet `#a855f7`, Action green `#10b981`.
  - **Zero-CDN Web Dashboard:**
    - Searchable database of parsed response strings with regex rule tester.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup parse-response`
  - Flags:
    - `--raw-string <string>`: Raw SMTP response line to parse and classify.
    - `--is-auto-update-rules`: Enable learning heuristics for unrecognized codes (default: `true`).
    - `--has-vendor-attribution`: Tag recognized security vendor in output (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SmtpSemanticResolution struct {
        ResolutionID     string    `json:"resolutionId"`
        RawResponseText  string    `json:"rawResponseText"`
        NormalizedAction string    `json:"normalizedAction"`
        DetectedVendor   string    `json:"detectedVendor"`
        ConfidenceScore  float64   `json:"confidenceScore"`
        ParsedAt         time.Time `json:"parsedAt"`
        IsThrottlingCode bool      `json:"isThrottlingCode"`
        HasPermanentFail bool      `json:"hasPermanentFail"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpSemanticResolutionRecord (
        resolution_id TEXT PRIMARY KEY,
        raw_response_text TEXT NOT NULL,
        normalized_action TEXT NOT NULL,
        detected_vendor TEXT NOT NULL,
        confidence_score REAL NOT NULL,
        parsed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_throttling_code INTEGER NOT NULL DEFAULT 1,
        has_permanent_fail INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_semantic_vendor ON SmtpSemanticResolutionRecord(detected_vendor, is_throttling_code);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/transport/response_parser.go`, driving routing decisions.

---

### Improvement 1250: Multi-Tenant Outbound Warmup Concurrency Semaphore HUD
- **ID:** `THROTTLE-SEMAPHORE-1250`
- **Cluster:** 125 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Enforce hard, non-blocking asynchronous semaphore limits across simultaneous outbound connections per destination ISP cluster (Google, Microsoft, Yahoo), preventing connection starvation on the host while obeying recipient limits.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live multi-cluster semaphore acquisition status bars with waiting queues and grant timers:
      ```text
      ┌─ [MULTI-TENANT OUTBOUND SEMAPHORE & CONCURRENCY ARBITER] ────────────────────────────┐
      │ ACTIVE CLUSTER SEMAPHORES: 5 | TOTAL GLOBAL ACQUIRED: 34 / 80 SLOTS                   │
      ├───────────────────────┬────────────┬─────────────────────────────┬─────────────────────┤
      │ DESTINATION CLUSTER   │ MAX SLOTS  │ ACQUIRED / IN-USE           │ WAITING QUEUE       │
      ├───────────────────────┼────────────┼─────────────────────────────┼─────────────────────┤
      │ Google MX Cluster     │   20 slots │ [████████████████░░░░] 16/20│ 2 requests waiting  │
      │ Microsoft 365 Cluster │   20 slots │ [██████████░░░░░░░░░░] 10/20│ 0 requests waiting  │
      │ Yahoo / AOL Cluster   │   15 slots │ [████░░░░░░░░░░░░░░░░]  4/15│ 0 requests waiting  │
      │ ProtonMail Relay      │   10 slots │ [████░░░░░░░░░░░░░░░░]  4/10│ 0 requests waiting  │
      ├───────────────────────┴────────────┴─────────────────────────────┴─────────────────────┤
      │ [ACQUISITION LATENCY: Mean Semaphore Wait: 1.4ms | Lock contention: 0.02%]           │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Acquired green `#10b981`, Warning amber `#f59e0b`, Wait cyan `#06b6d4`.
  - **Zero-CDN Web Dashboard:**
    - Real-time gauge cluster showing semaphore saturation per destination.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup semaphore-status`
  - Flags:
    - `--cluster <string>`: Filter status by destination cluster name.
    - `--max-wait-timeout <duration>`: Maximum time to wait for slot acquisition (default: `5s`).
    - `--is-non-blocking`: Reject immediately if semaphore capacity is full (default: `false`).
    - `--has-priority-queue`: Prioritize transactional messages over bulk warmup (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ClusterSemaphoreState struct {
        ClusterName     string        `json:"clusterName"`
        MaxCapacity     int           `json:"maxCapacity"`
        AcquiredSlots   int           `json:"acquiredSlots"`
        WaitingRequests int           `json:"waitingRequests"`
        MeanWaitTime    time.Duration `json:"meanWaitTime"`
        LastUpdated     time.Time     `json:"lastUpdated"`
        IsSaturated     bool          `json:"isSaturated"`
        HasWaitingQueue bool          `json:"hasWaitingQueue"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ClusterSemaphoreStateRecord (
        cluster_name TEXT PRIMARY KEY,
        max_capacity INTEGER NOT NULL,
        acquired_slots INTEGER NOT NULL DEFAULT 0,
        waiting_requests INTEGER NOT NULL DEFAULT 0,
        mean_wait_time_ms REAL NOT NULL DEFAULT 0.0,
        last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_saturated INTEGER NOT NULL DEFAULT 0,
        has_waiting_queue INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_semaphore_sat ON ClusterSemaphoreStateRecord(is_saturated, acquired_slots);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/warmup/transport/semaphore.go`, governing dial concurrency limits.


# Cluster 126: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 1251–1260)

---

### Improvement 1251: Sandboxed Local Shadow DOM Email Frame & Pixel Quarantine HUD
- **ID:** `CONTENT-SHADOW-1251`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Render raw HTML email messages inside a localized Web Component Shadow DOM boundary with full CSS style encapsulation and automated quarantine of external 1x1 tracking pixels to protect developer privacy.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Terminal Shadow DOM boundary tree showing quarantined attributes and isolated elements:
      ```text
      ┌─ [SANDBOXED SHADOW DOM EMAIL FRAME & PIXEL QUARANTINE] ──────────────────────────────┐
      │ SHADOW ROOT: #shadow-root (closed) | DOM NODES: 86 | TRACKING BEACONS NEUTRALIZED: 4  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [ENCAPSULATION INTEGRITY MATRIX]                                                     │
      │ ├─ Outer Styles Bleed: 0% (Scoped css isolation active via Shadow DOM boundary)      │
      │ ├─ Inline Javascript:  TOTAL PURGE (0 <script> / 0 event handler attributes)         │
      │ ├─ Discovered Pixels:  4x remote 1x1 GIF/PNG tracking beacons identified             │
      │ └─ Neutralization:     Replaced with inline base64 transparent SVG badges            │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [NEUTRALIZED TRACKING BEACONS]                                                       │
      │ ├─ [PIXEL #1] https://track.sendgrid.net/wf/open?upn=...   -> [QUARANTINED]          │
      │ ├─ [PIXEL #2] https://ci3.googleusercontent.com/proxy/... -> [QUARANTINED]          │
      │ └─ [PIXEL #3] https://pixel.hubspot.com/v1/track.gif       -> [QUARANTINED]          │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Shadow purple `#a855f7`, Quarantine rose `#f43f5e`, Scoped green `#10b981`.
  - **Zero-CDN Web Dashboard:**
    - Self-contained custom element `<email-sandbox-frame>` with zero external script imports.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content shadow-frame`
  - Flags:
    - `--html-file <path>`: Path to local HTML template or raw RFC 5322 MIME file.
    - `--is-pixel-quarantine`: Automatically replace tracking pixels with local placeholders (default: `true`).
    - `--has-closed-shadow`: Enforce strict `{mode: 'closed'}` on the Shadow Root (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ShadowDomSandboxReport struct {
        ReportID          string    `json:"reportId"`
        TemplateChecksum  string    `json:"templateChecksum"`
        TotalDomNodes     int       `json:"totalDomNodes"`
        QuarantinedPixels int       `json:"quarantinedPixels"`
        StrippedScripts   int       `json:"strippedScripts"`
        RenderedAt        time.Time `json:"renderedAt"`
        IsEncapsulated    bool      `json:"isEncapsulated"`
        HasPixelsBlocked  bool      `json:"hasPixelsBlocked"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ShadowDomSandboxRecord (
        report_id TEXT PRIMARY KEY,
        template_checksum TEXT NOT NULL,
        total_dom_nodes INTEGER NOT NULL,
        quarantined_pixels INTEGER NOT NULL DEFAULT 0,
        stripped_scripts INTEGER NOT NULL DEFAULT 0,
        rendered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_encapsulated INTEGER NOT NULL DEFAULT 1,
        has_pixels_blocked INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_shadow_checksum ON ShadowDomSandboxRecord(template_checksum, is_encapsulated);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/shadow_dom.go` providing safe preview isolation.

---

### Improvement 1252: Responsive Viewport Grid & Multi-Screen Media-Query Simulator HUD
- **ID:** `CONTENT-VIEWPORT-1252`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Emulate and validate email rendering across standard device viewports (Desktop 1200px, Tablet 768px, Mobile 375px / 320px), verifying responsive CSS breakpoints (`@media only screen and (max-width: 600px)`) in the terminal.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-viewport split view depicting character widths, overflow columns, and fluid table wraps:
      ```text
      ┌─ [RESPONSIVE EMAIL VIEWPORT GRID & MEDIA-QUERY SIMULATOR] ───────────────────────────┐
      │ TEMPLATE: newsletter-v4.html | FLUID WRAP: TESTED | HORIZONTAL OVERFLOW: 0 DETECTED   │
      ├───────────────────────┬────────────┬─────────────┬─────────────┬─────────────────────┤
      │ VIEWPORT PROFILE      │ WIDTH (PX) │ STATUS      │ WRAP BEHAVIOR│ MEDIA-QUERY TRIGGER │
      ├───────────────────────┼────────────┼─────────────┼─────────────┼─────────────────────┤
      │ Desktop Large         │ 1200px     │ [●] PASS    │ Fixed 600px │ Inactive (Base CSS) │
      │ Tablet Landscape      │  768px     │ [●] PASS    │ Fluid 100%  │ Inactive (Base CSS) │
      │ Mobile Large (iPhone) │  375px     │ [●] PASS    │ Stacked Cols│ ACTIVE (@media 600px│
      │ Mobile Small (SE)     │  320px     │ [●] PASS    │ Font Scaled │ ACTIVE (@media 480px│
      ├───────────────────────┴────────────┴─────────────┴─────────────┴─────────────────────┤
      │ [STACKING DIAGNOSTIC: 2-column tables successfully collapsed to single column at 375px│
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Pass green `#10b981`, Active cyan `#06b6d4`, Warning amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Live responsive multi-frame viewer with synchronized scroll and zero external dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content test-viewports`
  - Flags:
    - `--template <path>`: Source HTML email file to test.
    - `--viewports <string>`: Comma-separated list of widths (default: `1200,768,375,320`).
    - `--is-strict-overflow`: Flag any element exceeding viewport bounding box as error (default: `true`).
    - `--has-table-collapser`: Verify automated single-column stacking for Outlook/MUA tables (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ViewportTestProfile struct {
        ProfileID       string    `json:"profileId"`
        TemplatePath    string    `json:"templatePath"`
        ViewportWidth   int       `json:"viewportWidth"`
        HasOverflow     bool      `json:"hasOverflow"`
        IsResponsive    bool      `json:"isResponsive"`
        IsStackedColumn bool      `json:"isStackedColumn"`
        TestedAt        time.Time `json:"testedAt"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ViewportTestProfileRecord (
        profile_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        viewport_width INTEGER NOT NULL,
        has_overflow INTEGER NOT NULL DEFAULT 0,
        is_responsive INTEGER NOT NULL DEFAULT 1,
        is_stacked_column INTEGER NOT NULL DEFAULT 1,
        tested_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_viewport_test ON ViewportTestProfileRecord(template_path, viewport_width);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/email/content/viewports.go` validating responsive designs.

---

### Improvement 1253: Email CSS Inliner & Redundant Specificity Optimizer HUD
- **ID:** `CONTENT-INLINER-1253`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically inline global `<style>` tags into direct element `style="..."` attributes with CSS specificity calculation, stripping unused selectors and preserving `@media` queries in a preserved head tag.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - CSS compression and inlining metric summary HUD:
      ```text
      ┌─ [EMAIL CSS INLINER & REDUNDANT SPECIFICITY OPTIMIZER] ──────────────────────────────┐
      │ SOURCE SIZE: 48.2 KB | INLINED SIZE: 34.1 KB | COMPRESSION RATIO: -29.2%             │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [INLINING OPTIMIZATION METRICS]                                                      │
      │ ├─ Total CSS Rules Parsed:  142 rules (Resolved via W3C specificity cascade)         │
      │ ├─ Inlined Element Styles:  310 direct DOM element node modifications                 │
      │ ├─ Redundant Rules Stripped: 48 unused classes (.hidden, .btn-secondary, etc.)       │
      │ ├─ Preserved in <style>:    8 @media responsive queries (Safe for modern MUAs)       │
      │ └─ Outlook VML Tags:        Preserved <!-[if mso]> conditional wrappers intact       │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [ACTIONS: Press 'w' to write inlined HTML, 'd' to view inline diff, 'q' to quit]     │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Size cyan `#06b6d4`, Optimized green `#10b981`, Purged slate `#94a3b8`.
  - **Zero-CDN Web Dashboard:**
    - Side-by-side visual HTML diff viewer showing stripped selectors and resulting markup.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content inline-css`
  - Flags:
    - `--input <path>`: Source HTML email file with `<style>` blocks.
    - `--output <path>`: Destination path for inlined HTML output.
    - `--is-purge-unused`: Remove unreferenced CSS classes and IDs (default: `true`).
    - `--has-mso-preservation`: Keep Microsoft Outlook MML/VML conditional tags intact (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CssInlinerReport struct {
        ReportID          string    `json:"reportId"`
        SourceFilePath    string    `json:"sourceFilePath"`
        OriginalSizeBytes int64     `json:"originalSizeBytes"`
        InlinedSizeBytes  int64     `json:"inlinedSizeBytes"`
        RulesInlined      int       `json:"rulesInlined"`
        RulesStripped     int       `json:"rulesStripped"`
        ProcessedAt       time.Time `json:"processedAt"`
        IsSuccessful      bool      `json:"isSuccessful"`
        HasMsoPreserved   bool      `json:"hasMsoPreserved"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CssInlinerExecutionRecord (
        report_id TEXT PRIMARY KEY,
        source_file_path TEXT NOT NULL,
        original_size_bytes INTEGER NOT NULL,
        inlined_size_bytes INTEGER NOT NULL,
        rules_inlined INTEGER NOT NULL DEFAULT 0,
        rules_stripped INTEGER NOT NULL DEFAULT 0,
        processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_successful INTEGER NOT NULL DEFAULT 1,
        has_mso_preserved INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_inliner_file ON CssInlinerExecutionRecord(source_file_path, is_successful);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/email/content/inliner.go` preparing raw markup for transmission.

---

### Improvement 1254: Dark Mode Palette Inversion & Contrast Accessibility Auditor HUD
- **ID:** `CONTENT-DARKMODE-1254`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit and simulate how major email clients (Apple Mail, Gmail, Outlook iOS) automatically invert HTML email colors in dark mode, verifying that text maintains WCAG AAA (7:1) contrast against dynamically inverted backgrounds.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual light/dark contrast comparison table with WCAG compliance indicators:
      ```text
      ┌─ [DARK MODE PALETTE INVERSION & WCAG CONTRAST AUDITOR] ──────────────────────────────┐
      │ TESTED CLIENTS: Apple Mail, Gmail Android, Outlook iOS | OVERALL ACCESSIBILITY: WCAG AAA│
      ├───────────────────────┬────────────┬─────────────┬─────────────┬─────────────────────┤
      │ ELEMENT / SELECTOR    │ LIGHT MODE │ INVERTED DARK│ CONTRAST    │ ACCESSIBILITY PASS  │
      ├───────────────────────┼────────────┼─────────────┼─────────────┼─────────────────────┤
      │ Body Text (p, td)     │ #1e293b    │ #f1f5f9     │ 14.2:1      │ [●] PASS (AAA)      │
      │ Heading (h1, h2)      │ #0f172a    │ #ffffff     │ 18.5:1      │ [●] PASS (AAA)      │
      │ Primary CTA Button    │ #2563eb    │ #3b82f6     │  8.1:1      │ [●] PASS (AAA)      │
      │ Muted Footer Legal    │ #64748b    │ #94a3b8     │  4.8:1      │ [▲] PASS (AA Only)  │
      ├───────────────────────┴────────────┴─────────────┴─────────────┴─────────────────────┤
      │ [DARK MODE FIX: Injected [data-ogsc] / @media (prefers-color-scheme: dark) overrides] │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: PassAAA green `#10b981`, PassAA amber `#f59e0b`, Fail rose `#f43f5e`.
  - **Zero-CDN Web Dashboard:**
    - Interactive split-slider widget simulating light and dark mode side-by-side.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content audit-darkmode`
  - Flags:
    - `--html-file <path>`: Source HTML template.
    - `--target-wcag-level <string>`: Target compliance (`AA` or `AAA`; default: `AAA`).
    - `--is-auto-patch`: Inject modern CSS dark mode media queries for broken contrasts (default: `true`).
    - `--has-image-inversion-guard`: Add CSS filters to prevent logo color distortions in dark mode (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeContrastCheck struct {
        CheckID          string    `json:"checkId"`
        ElementSelector  string    `json:"elementSelector"`
        LightHexColor    string    `json:"lightHexColor"`
        DarkHexColor     string    `json:"darkHexColor"`
        ContrastRatio    float64   `json:"contrastRatio"`
        WcagLevelPassed  string    `json:"wcagLevelPassed"`
        AuditedAt        time.Time `json:"auditedAt"`
        IsCompliant      bool      `json:"isCompliant"`
        HasInjectedPatch bool      `json:"hasInjectedPatch"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DarkModeContrastCheckRecord (
        check_id TEXT PRIMARY KEY,
        element_selector TEXT NOT NULL,
        light_hex_color TEXT NOT NULL,
        dark_hex_color TEXT NOT NULL,
        contrast_ratio REAL NOT NULL,
        wcag_level_passed TEXT NOT NULL,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_compliant INTEGER NOT NULL DEFAULT 1,
        has_injected_patch INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_darkmode_comp ON DarkModeContrastCheckRecord(wcag_level_passed, is_compliant);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/darkmode.go` guaranteeing cross-client legibility.

---

### Improvement 1255: MIME Multipart/Alternative RFC 2046 Synchronizer & Plaintext Generator
- **ID:** `CONTENT-MIME-1255`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically parse rich HTML emails and synthesize clean, formatted `text/plain` alternatives (with Markdown headers, ASCII table borders, and URL footnote lists), ensuring strict RFC 2046 boundary compliance.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Side-by-side terminal viewer comparing rendered HTML markup with generated plain text:
      ```text
      ┌─ [MIME MULTIPART/ALTERNATIVE SYNCHRONIZER: RFC 2046] ────────────────────────────────┐
      │ BOUNDARY: =_gitmap_boundary_4819 | HTML PART: 24.2 KB | PLAINTEXT PART: 3.8 KB        │
      ├──────────────────────────────────────┬───────────────────────────────────────────────┤
      │ RENDERED HTML STRUCTURE              │ GENERATED text/plain EQUIVALENT               │
      ├──────────────────────────────────────┼───────────────────────────────────────────────┤
      │ <h1>Infrastructure Alert</h1>        │ # Infrastructure Alert                        │
      │ <p>Your cluster CPU reached 92%.</p> │ Your cluster CPU reached 92%.                 │
      │ <a href="https://app.io/d">View</a>  │ [View Dashboard] [1]                          │
      │ <table><tr><td>Node</td>...</tr>...  │ --------------------------------------------- │
      │                                      │ [1] https://app.io/d                          │
      ├──────────────────────────────────────┴───────────────────────────────────────────────┤
      │ [MIME HEALTH: Boundaries aligned; text/plain appears first as per RFC 2046: PASS]    │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Boundary cyan `#06b6d4`, Plaintext yellow `#facc15`, Pass green `#10b981`.
  - **Zero-CDN Web Dashboard:**
    - Live MIME structure tree and plaintext generation playground.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content generate-plaintext`
  - Flags:
    - `--html-input <path>`: Source HTML template.
    - `--output-mime <path>`: Target synthesized RFC 5322 MIME message file.
    - `--is-footnote-links`: Place hyperlinked URLs in numbered footnotes (default: `true`).
    - `--has-ascii-tables`: Render HTML tables as clean ASCII text tables (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MimeMultipartSyncReport struct {
        ReportID          string    `json:"reportId"`
        BoundaryKey       string    `json:"boundaryKey"`
        HtmlSizeBytes     int64     `json:"htmlSizeBytes"`
        PlaintextSizeBytes int64    `json:"plaintextSizeBytes"`
        LinksExtracted    int       `json:"linksExtracted"`
        GeneratedAt       time.Time `json:"generatedAt"`
        IsRfcCompliant    bool      `json:"isRfcCompliant"`
        HasAsciiTables    bool      `json:"hasAsciiTables"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MimeMultipartSyncRecord (
        report_id TEXT PRIMARY KEY,
        boundary_key TEXT NOT NULL,
        html_size_bytes INTEGER NOT NULL,
        plaintext_size_bytes INTEGER NOT NULL,
        links_extracted INTEGER NOT NULL DEFAULT 0,
        generated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_rfc_compliant INTEGER NOT NULL DEFAULT 1,
        has_ascii_tables INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_mime_boundary ON MimeMultipartSyncRecord(boundary_key, is_rfc_compliant);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/email/content/mime_synth.go` powering transactional message assembly.

---

### Improvement 1256: RFC 2387 Multipart/Related CID Inline Asset Packager HUD
- **ID:** `CONTENT-CID-1256`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Convert remote or local email assets (logos, icons) into Content-ID (`cid:`) attachments packaged in `multipart/related` MIME containers, guaranteeing instant image display without triggering ESP remote image blocking banners.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Embedded CID asset packing matrix with byte size budgets and Content-Type inspection:
      ```text
      ┌─ [RFC 2387 MULTIPART/RELATED CID ASSET PACKAGER] ────────────────────────────────────┐
      │ ATTACHED CID IMAGES: 3 | TOTAL ASSET PAYLOAD: 42.1 KB (SAFE BUDGET: < 100 KB)        │
      ├───────────────────────┬────────────┬─────────────────────────────┬─────────────────────┤
      │ ASSET FILENAME        │ CONTENT-ID │ CONTENT-TYPE / CHECKSUM     │ SIZE / COMPRESSION  │
      ├───────────────────────┼────────────┼─────────────────────────────┼─────────────────────┤
      │ company-logo.png      │ <logo@cid> │ image/png (sha256:4f1a...)  │ 24.8 KB [●] OPTIMAL │
      │ icon-check.png        │ <chk@cid>  │ image/png (sha256:9b20...)  │  8.4 KB [●] OPTIMAL │
      │ social-github.png     │ <gh@cid>   │ image/png (sha256:1a8e...)  │  8.9 KB [●] OPTIMAL │
      ├───────────────────────┴────────────┴─────────────────────────────┴─────────────────────┤
      │ [HTML REPLACEMENT: <img src="logo.png"> -> <img src="cid:logo@cid"> completed]        │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Asset cyan `#06b6d4`, CID purple `#c084fc`, Budget green `#10b981`.
  - **Zero-CDN Web Dashboard:**
    - Interactive asset dropzone with real-time PNG/JPEG micro-compression and CID preview.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content package-cid`
  - Flags:
    - `--html-file <path>`: Source HTML email file.
    - `--assets-dir <path>`: Directory containing local images to embed.
    - `--max-budget-kb <int>`: Max allowable aggregate CID size (default: `100`).
    - `--is-auto-compress`: Losslessly optimize images before packaging (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CidEmbeddedAsset struct {
        AssetID        string    `json:"assetId"`
        Filename       string    `json:"filename"`
        ContentID      string    `json:"contentId"`
        ContentType    string    `json:"contentType"`
        SizeBytes      int64     `json:"sizeBytes"`
        Sha256Checksum string    `json:"sha256Checksum"`
        PackagedAt     time.Time `json:"packagedAt"`
        IsCompressed   bool      `json:"isCompressed"`
        HasValidCid    bool      `json:"hasValidCid"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CidEmbeddedAssetRecord (
        asset_id TEXT PRIMARY KEY,
        filename TEXT NOT NULL,
        content_id TEXT NOT NULL UNIQUE,
        content_type TEXT NOT NULL,
        size_bytes INTEGER NOT NULL,
        sha256_checksum TEXT NOT NULL,
        packaged_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_compressed INTEGER NOT NULL DEFAULT 1,
        has_valid_cid INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_cid_checksum ON CidEmbeddedAssetRecord(sha256_checksum, is_compressed);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/cid_packager.go` building rich multipart envelopes.

---

### Improvement 1257: Outbound Email Link Reputation & Redirect Red-Flag Analyzer HUD
- **ID:** `CONTENT-LINKS-1257`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan all hyperlink destinations embedded in email templates against known threat intelligence blocklists, checking for URL shorteners, excessive redirect hops, unencrypted HTTP targets, or suspicious TLDs that harm ESP reputation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Hyperlink safety audit ledger showing redirect trace hops and domain age ratings:
      ```text
      ┌─ [OUTBOUND EMAIL LINK REPUTATION & REDIRECT ANALYZER] ───────────────────────────────┐
      │ AUDITED HYPERLINKS: 6 | SUSPICIOUS: 0 | SHORTENERS: 0 | REDIRECT HOPS: ALL < 2       │
      ├───────────────────────────────────┬────────────┬─────────────┬─────────────────────────┤
      │ TARGET URL                        │ PROTOCOL   │ HOPS        │ REPUTATION / STATUS     │
      ├───────────────────────────────────┼────────────┼─────────────┼─────────────────────────┤
      │ https://acme.org/dashboard        │ HTTPS Strict│ 0 (Direct)  │ [●] SAFE (Domain age 8y)│
      │ https://docs.acme.org/v2/api      │ HTTPS Strict│ 1 (301 Perm)│ [●] SAFE (Canonical)    │
      │ http://insecure-cdn.biz/img       │ HTTP Clear │ 0 (Direct)  │ [▲] WARN (Cleartext)    │
      ├───────────────────────────────────┴────────────┴─────────────┴─────────────────────────┤
      │ [SECURITY POLICY: Cleartext HTTP target detected. Upgraded to HTTPS automatically]   │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Safe green `#10b981`, Warn amber `#f59e0b`, Danger rose `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Interactive link graph visualizer mapping outbound anchor tags to resolved canonical endpoints.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content scan-links`
  - Flags:
    - `--html-file <path>`: Source HTML template to audit.
    - `--max-allowed-hops <int>`: Max allowable HTTP redirects before warning (default: `2`).
    - `--is-https-enforced`: Flag and rewrite non-TLS HTTP links to HTTPS (default: `true`).
    - `--has-shortener-ban`: Reject generic shorteners (bit.ly, tinyurl) known to trigger spam filters (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LinkAuditItem struct {
        LinkID          string    `json:"linkId"`
        OriginalHref    string    `json:"originalHref"`
        FinalResolvedUrl string   `json:"finalResolvedUrl"`
        RedirectHops    int       `json:"redirectHops"`
        IsHttpsStrict   bool      `json:"isHttpsStrict"`
        IsSafe          bool      `json:"isSafe"`
        HasShortener    bool      `json:"hasShortener"`
        LastAuditedAt   time.Time `json:"lastAuditedAt"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS LinkAuditItemRecord (
        link_id TEXT PRIMARY KEY,
        original_href TEXT NOT NULL,
        final_resolved_url TEXT NOT NULL,
        redirect_hops INTEGER NOT NULL DEFAULT 0,
        is_https_strict INTEGER NOT NULL DEFAULT 1,
        is_safe INTEGER NOT NULL DEFAULT 1,
        has_shortener INTEGER NOT NULL DEFAULT 0,
        last_audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_link_safe ON LinkAuditItemRecord(is_safe, is_https_strict);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/email/content/link_audit.go` auditing marketing and transactional bodies.

---

### Improvement 1258: RFC 8058 One-Click Unsubscribe Header & HTTPS Receiver HUD
- **ID:** `CONTENT-UNSUB-1258`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify and inject mandatory RFC 8058 `List-Unsubscribe` and `List-Unsubscribe-Post: List-Unsubscribe=One-Click` headers, providing a built-in HTTPS handler to process one-click POST unsubscribe requests seamlessly.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - RFC 8058 validation status and live incoming unsubscribe request monitor:
      ```text
      ┌─ [RFC 8058 ONE-CLICK LIST-UNSUBSCRIBE GOVERNOR] ─────────────────────────────────────┐
      │ POLICY STATUS: 100% COMPLIANT (Google/Yahoo 2024 Bulk Sender Requirement Met)        │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [INJECTED HEADER PAYLOAD]                                                            │
      │ List-Unsubscribe: <https://mail.domain.io/unsub?id=c819>, <mailto:unsub@domain.io>   │
      │ List-Unsubscribe-Post: List-Unsubscribe=One-Click                                    │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [RECENT ONE-CLICK POST DISPATCHES]                                                   │
      │ [14:41:02] POST /unsub?id=c819 from Google Mail User-Agent -> PROCESSED (0.42ms)     │
      │ Recipient suppressed from active warmup cohort instantaneously without user friction │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Compliant green `#10b981`, Post cyan `#06b6d4`, Header violet `#a855f7`.
  - **Zero-CDN Web Dashboard:**
    - Live suppression list manager with one-click unsubscribe event audit trail.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content verify-unsub`
  - Flags:
    - `--domain <string>`: Sending domain to check header templates for.
    - `--is-auto-inject`: Inject RFC 8058 headers automatically during MIME serialization (default: `true`).
    - `--has-mailto-fallback`: Include mailto URI alongside HTTPS endpoint (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubscribeConfiguration struct {
        ConfigID          string    `json:"configId"`
        DomainName        string    `json:"domainName"`
        HttpsEndpoint     string    `json:"httpsEndpoint"`
        MailtoAddress     string    `json:"mailtoAddress"`
        LastVerifiedAt    time.Time `json:"lastVerifiedAt"`
        IsOneClickEnabled bool      `json:"isOneClickEnabled"`
        HasMailtoFallback bool      `json:"hasMailtoFallback"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnsubscribeConfigurationRecord (
        config_id TEXT PRIMARY KEY,
        domain_name TEXT NOT NULL UNIQUE,
        https_endpoint TEXT NOT NULL,
        mailto_address TEXT NOT NULL,
        last_verified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_one_click_enabled INTEGER NOT NULL DEFAULT 1,
        has_mailto_fallback INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_unsub_domain ON UnsubscribeConfigurationRecord(domain_name, is_one_click_enabled);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/unsubscribe.go` ensuring deliverability standards compliance.

---

### Improvement 1259: Email Spam Trigger Heuristic Scanner & Keyword Weighting HUD
- **ID:** `CONTENT-HEURISTIC-1259`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit subject lines and email body content for spam trigger words ("100% free", "guaranteed cash", "act now", "urgent wire transfer"), calculating compound heuristic risk scores before mail is sent.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Interactive keyword risk heatmap with line-number citations and alternative copy suggestions:
      ```text
      ┌─ [EMAIL SPAM TRIGGER HEURISTIC SCANNER & KEYWORD WEIGHTING] ─────────────────────────┐
      │ AUDITED CONTENT: "Quarterly Growth Briefing" | RISK SCORE: 0.8 / 10.0 (LOW RISK)      │
      ├───────────────────────┬────────────┬─────────────────────────────┬─────────────────────┤
      │ TRIGGER KEYWORD       │ LINE NO.   │ RISK WEIGHT / CATEGORY      │ SUGGESTED REWRITE   │
      ├───────────────────────┼────────────┼─────────────────────────────┼─────────────────────┤
      │ "guaranteed ROI"      │ Line #14   │ +0.5 (Financial Claim)      │ "projected returns" │
      │ "act now"             │ Line #42   │ +0.3 (High Pressure Urgency)│ "review options"    │
      ├───────────────────────┴────────────┴─────────────────────────────┴─────────────────────┤
      │ [HEURISTIC EVALUATION: Subject Line: PRISTINE | Body Text: OPTIMAL (Safe for Inbox)]  │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: LowRisk green `#10b981`, Suggestion cyan `#06b6d4`, Weight amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Live content editor with real-time heuristic keyword highlighting and synonym suggestions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content scan-heuristics`
  - Flags:
    - `--file <path>`: Source HTML or text email file to scan.
    - `--max-acceptable-risk <float>`: Risk threshold to fail CI/CD build (default: `2.5`).
    - `--is-subject-scanned`: Include email Subject header in heuristic audit (default: `true`).
    - `--has-auto-suggestions`: Output suggested contextual rewrites for flagged phrases (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamHeuristicResult struct {
        ResultID        string    `json:"resultId"`
        FileChecksum    string    `json:"fileChecksum"`
        TotalRiskScore  float64   `json:"totalRiskScore"`
        FlaggedKeywords int       `json:"flaggedKeywords"`
        EvaluatedAt     time.Time `json:"evaluatedAt"`
        IsSafeToSend    bool      `json:"isSafeToSend"`
        HasSuggestions  bool      `json:"hasSuggestions"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamHeuristicResultRecord (
        result_id TEXT PRIMARY KEY,
        file_checksum TEXT NOT NULL,
        total_risk_score REAL NOT NULL DEFAULT 0.0,
        flagged_keywords INTEGER NOT NULL DEFAULT 0,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_safe_to_send INTEGER NOT NULL DEFAULT 1,
        has_suggestions INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_heuristic_safe ON SpamHeuristicResultRecord(is_safe_to_send, total_risk_score);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/email/content/heuristics.go` serving as a pre-flight deliverability gate.

---

### Improvement 1260: Cryptographic DKIM Canonicalization (relaxed/relaxed) Visualizer HUD
- **ID:** `CONTENT-DKIMVIS-1260`
- **Cluster:** 126 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Visualize exact byte-by-byte string transformations during RFC 6376 DKIM canonicalization (`relaxed/relaxed` vs `simple/simple`), showing how whitespace collapsing, lowercase header conversion, and body hashes are computed prior to signing.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Terminal diff inspector showing raw bytes transformed into canonicalized hashing inputs:
      ```text
      ┌─ [CRYPTOGRAPHIC DKIM CANONICALIZATION (relaxed/relaxed) VISUALIZER] ─────────────────┐
      │ SELECTOR: s1._domainkey.gitmap.io | ALGORITHM: rsa-sha256 | HASH ACCURACY: 100% PASS   │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [HEADER CANONICALIZATION (relaxed: lowercase key, collapse WSP, unfold CRLF)]        │
      │ RAW:       Subject:  Quarterly   Engineering   Briefing 
                          │
      │ CANONICAL: subject:Quarterly Engineering Briefing
                                │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [BODY HASH COMPUTATION (relaxed: collapse trailing empty lines, strip trailing WSP)] │
      │ Raw Body Size: 4,120 bytes | Canonicalized Body Size: 3,980 bytes                    │
      │ Computed Body Hash (bh=): 8vF9d1pA2xKl09+M4Qe8sY2bU1vOqE5=                           │
      │ Signature Status: [●] VALID RSA-2048 CRYPTOGRAPHIC SIGNATURE GENERATED               │
      └────────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Canonical teal `#14b8a6`, Hash purple `#c084fc`, Valid green `#22c55e`.
  - **Zero-CDN Web Dashboard:**
    - Interactive ASCII/Hex canonicalization byte difference explorer.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content visualize-dkim`
  - Flags:
    - `--raw-email <path>`: Source MIME message file to canonicalize and test.
    - `--selector <string>`: DKIM DNS selector identifier (default: `s1`).
    - `--canonicalization <string>`: Scheme (`relaxed/relaxed` or `simple/simple`; default: `relaxed/relaxed`).
    - `--is-signature-verified`: Cryptographically verify generated signature against public key (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DkimVisualizationState struct {
        VisualID         string    `json:"visualId"`
        SelectorKey      string    `json:"selectorKey"`
        Canonicalization string    `json:"canonicalization"`
        BodyHashComputed string    `json:"bodyHashComputed"`
        SignaturePayload string    `json:"signaturePayload"`
        ComputedAt       time.Time `json:"computedAt"`
        IsValidSignature bool      `json:"isValidSignature"`
        HasRelaxedHeader bool      `json:"hasRelaxedHeader"`
        HasRelaxedBody   bool      `json:"hasRelaxedBody"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DkimVisualizationStateRecord (
        visual_id TEXT PRIMARY KEY,
        selector_key TEXT NOT NULL,
        canonicalization TEXT NOT NULL,
        body_hash_computed TEXT NOT NULL,
        signature_payload TEXT NOT NULL,
        computed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_valid_signature INTEGER NOT NULL DEFAULT 1,
        has_relaxed_header INTEGER NOT NULL DEFAULT 1,
        has_relaxed_body INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_dkim_vis_selector ON DkimVisualizationStateRecord(selector_key, is_valid_signature);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/dkim_visualizer.go` providing cryptographic diagnostics.
