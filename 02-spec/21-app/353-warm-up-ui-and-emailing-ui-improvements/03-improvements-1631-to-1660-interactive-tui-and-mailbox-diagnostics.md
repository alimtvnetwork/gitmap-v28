# Specification 353: Improvements 1631 to 1660 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 353
- **Slug:** `353-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/353-warm-up-ui-and-emailing-ui-improvements/03-improvements-1631-to-1660-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (1631 - 1660)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 164** | Items 1631–1640 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 165** | Items 1641–1650 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 166** | Items 1651–1660 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 164: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 1631–1640)

---

### Improvement 1631: Interactive Bubble Tea Mailbox Health Matrix & Multi-Tenant Seed Shard Visualizer
- **ID:** `SEED-MATRIX-TOPOLOGY-1631`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a real-time, interactive Bubble Tea terminal grid visualizing multi-tenant seed mailbox shards across global edge points of presence (PoPs), tracking Anycast routing latencies, tenant isolation boundaries, ISP inbox placement cohorts (Gmail, M365, Yahoo, Fastmail), and automated seed cluster drain states.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [BUBBLE TEA MULTI-TENANT SEED HEALTH MATRIX & TOPOLOGY] ───────────────────────────────┐
      │ SHARDS: 32 | ACTIVE SEEDS: 1,280 | DEGRADED: 14 | QUARANTINED: 2 | MEAN HEALTH: 99.4%       │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ SEED SHARD ID | TENANT GROUP | EDGE REGION | INBOX % | RTT (ms) | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ shard-gsuite-us-301 | Enterprise A | us-east-1 | 99.8% | 14.2ms | [●] STABLE
      │ shard-m365-eu-204 | Fintech B | eu-central-1 | 99.1% | 18.6ms | [●] STABLE
      │ shard-yahoo-us-109 | SaaS Core | us-west-2 | 94.2% | 36.1ms | [▲] WARN
      │ shard-fastm-ap-055 | Global High | ap-east-1 | 99.7% | 26.4ms | [●] STABLE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: shard-yahoo-us-109] Health: 94.2% | Target MX: mx-core.mail.yahoo.com | Shard Drain: IDLE
      │ Keys: [↑/↓: Navigate] [enter: Drill Down Shard] [f: Filter Alerts] [r: Force Health Probe] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Zero-CDN responsive HTML5/SVG topology canvas rendering interactive Voronoi regional cell partitions, local CSS variable gradients (`--color-emerald: #10b981; --color-amber: #f59e0b`), and zero-external asset offline inspection cards.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-matrix`
  - Flags:
    - `--tenant <string>`: Filter seed pool view by tenant workspace ID (default: `all`).
    - `--provider <string>`: Filter seed pool view by ESP provider (`google`, `microsoft`, `yahoo`, `fastmail`, `all`; default: `all`).
    - `--refresh-rate <duration>`: Interval for real-time matrix refresh and probe streaming (default: `1s`).
    - `--is-interactive`: Enable interactive keyboard navigation cursor and detail drill-down (default: `true`).
    - `--has-alerts-only`: Filter matrix view to display only degraded or quarantined seed shards (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedShardHealthMatrixRecord struct {
        ShardID string           `json:"shard_id"`
        TenantID string          `json:"tenant_id"`
        ProviderName string      `json:"provider_name"`
        ClusterRegion string     `json:"cluster_region"`
        HealthScorePercent float64 `json:"health_score_percent"`
        InboxPlacementRate float64 `json:"inbox_placement_rate"`
        SpamFolderRate float64   `json:"spam_folder_rate"`
        LatencyRttMs float64     `json:"latency_rtt_ms"`
        IsHealthy bool           `json:"is_healthy"`
        IsQuarantined bool       `json:"is_quarantined"`
        HasAnomalousSpike bool   `json:"has_anomalous_spike"`
        HasVerifiedAuth bool     `json:"has_verified_auth"`
    }

    type SeedShardHealthMatrixModel struct {
        Records []SeedShardHealthMatrixRecord `json:"records"`
        CursorIndex int          `json:"cursor_index"`
        SelectedRecord SeedShardHealthMatrixRecord `json:"selected_record"`
        IsStreaming bool         `json:"is_streaming"`
        HasAlerts bool           `json:"has_alerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedShardHealthMatrixRecord (
        shard_id TEXT PRIMARY KEY,
        tenant_id TEXT NOT NULL,
        provider_name TEXT NOT NULL,
        cluster_region TEXT NOT NULL,
        health_score_percent REAL NOT NULL DEFAULT 100.0,
        inbox_placement_rate REAL NOT NULL DEFAULT 1.0,
        spam_folder_rate REAL NOT NULL DEFAULT 0.0,
        latency_rtt_ms REAL NOT NULL DEFAULT 0.0,
        is_healthy INTEGER NOT NULL DEFAULT 1,
        is_quarantined INTEGER NOT NULL DEFAULT 0,
        has_anomalous_spike INTEGER NOT NULL DEFAULT 0,
        has_verified_auth INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics seed-matrix`. Reads real-time multi-tenant health records from `repodb/pipeline.db` using read-replica transactions and dispatches non-blocking tea.Msg events to Bubble Tea.

---

### Improvement 1632: Seed Mailbox IMAP Auto-Drain & Spam Rescue Inspector HUD
- **ID:** `SEED-IMAP-DRAIN-AUTO-1632`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide high-speed IMAP queue flushing, automated spam-to-inbox rescue telemetry, and quarantine inspection across active warmup seed clusters to eliminate mailbox quota saturation and reinforce sender reputation feedback.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [IMAP AUTO-DRAIN & SPAM RESCUE HUD] ───────────────────────────────────────────────────┐
      │ ACTIVE SESSIONS: 18 | DRAINED (1h): 4,892 | RESCUED: 142 | STORAGE RECLAIMED: 84.2 MB       │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ SESSION ID | SEED MAILBOX | INBOX RESCUE | PURGED | THROUGHPUT | STATE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ drain-ssn-8812 | seed-01@warm.dev | 14 msgs | 320 msgs | 48.2 msg/s | [●] ACTIVE
      │ drain-ssn-8813 | seed-02@corp.io | 8 msgs | 195 msgs | 39.5 msg/s | [●] ACTIVE
      │ drain-ssn-8814 | seed-03@inbox.net | 29 msgs | 512 msgs | 52.0 msg/s | [●] ACTIVE
      │ drain-ssn-8815 | seed-04@outpost.org | 0 msgs | 42 msgs | 12.1 msg/s | [●] FLUSH
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: drain-ssn-8814] Rescued 29 false-positives from Junk to Inbox | IMAP EXPUNGE: OK
      │ Keys: [↑/↓: Select Session] [p: Pause/Resume] [f: Force Instant Rescue] [c: Clear Expunged] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Client-rendered SVG progress gauges with system CSS monospace tables, embedded tabular progress meters, and pure local SQLite state binding.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics imap-drain`
  - Flags:
    - `--mailbox-id <string>`: Specific seed mailbox to drain and rescue (default: `all`).
    - `--rescue-threshold <int>`: Max false-positive spam messages to rescue per cycle (default: `50`).
    - `--batch-limit <int>`: Max messages purged per IMAP EXPUNGE batch (default: `250`).
    - `--is-rescue-active`: Automatically move warmup messages from Spam/Junk to Inbox (default: `true`).
    - `--has-expunge-verify`: Verify message deletion via UID search before closing session (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedImapDrainRecord struct {
        SessionID string         `json:"session_id"`
        MailboxAddress string    `json:"mailbox_address"`
        RescuedCount int         `json:"rescued_count"`
        PurgedCount int          `json:"purged_count"`
        ReclaimedBytes int64     `json:"reclaimed_bytes"`
        ThroughputRate float64   `json:"throughput_rate"`
        IsActive bool            `json:"is_active"`
        IsExpunged bool          `json:"is_expunged"`
        HasRescueError bool      `json:"has_rescue_error"`
        HasQuotaAlert bool       `json:"has_quota_alert"`
    }

    type SeedImapDrainModel struct {
        Sessions []SeedImapDrainRecord `json:"sessions"`
        ActiveCursor int         `json:"active_cursor"`
        TotalRescued int         `json:"total_rescued"`
        IsAutoDraining bool      `json:"is_auto_draining"`
        HasAlerts bool           `json:"has_alerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedImapDrainRecord (
        session_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        rescued_count INTEGER NOT NULL DEFAULT 0,
        purged_count INTEGER NOT NULL DEFAULT 0,
        reclaimed_bytes INTEGER NOT NULL DEFAULT 0,
        throughput_rate REAL NOT NULL DEFAULT 0.0,
        is_active INTEGER NOT NULL DEFAULT 1,
        is_expunged INTEGER NOT NULL DEFAULT 1,
        has_rescue_error INTEGER NOT NULL DEFAULT 0,
        has_quota_alert INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `gitmap email warmup diagnostics imap-drain`. Connects to target IMAP endpoints using TLS 1.3 socket workers and logs drain audits into `repodb/pipeline.db`.

---

### Improvement 1633: Seed Domain DNS Propagation & Resolver Consistency Matrix
- **ID:** `SEED-DNS-PROPAGATION-1633`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Verify global Anycast DNS propagation (SPF, DKIM, DMARC, MX) across authoritative edge nameservers, surfacing stale TTLs, DNSSEC cryptographic validation errors, and regional NXDOMAIN anomalies in real time.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED DOMAIN DNS RESOLVER CONSISTENCY MATRIX] ─────────────────────────────────────────┐
      │ RESOLVERS CHECKED: 16 | SEED DOMAINS: 48 | SPF SYNC: 100% | DKIM VALID: 100% | DMARC P=REJECT│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ RESOLVER NODE | LOCATION | MX RECORD | SPF MATCH | DKIM SIG | DNSSEC
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Cloudflare (1.1.1.1) | Global Anycast | 10 mx1.dev | [✓] PASS | [✓] VALID | [●] SECURE
      │ Google (8.8.8.8) | Global Anycast | 10 mx1.dev | [✓] PASS | [✓] VALID | [●] SECURE
      │ Quad9 (9.9.9.9) | Zurich, CH | 10 mx1.dev | [✓] PASS | [✓] VALID | [●] SECURE
      │ OpenDNS (208.67...) | Ashburn, US | 10 mx1.dev | [✓] PASS | [✓] VALID | [●] SECURE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: Quad9] Query Latency: 11.4ms | EDNS Client Subnet: ACTIVE | Record TTL: 300s
      │ Keys: [↑/↓: Change Resolver] [d: Toggle Domain Selector] [r: Flush Resolver Cache] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive diagnostic grid with zero-CDN CSS flexbox wrappers, SVG shield badges for DNSSEC status, and local JSON export modal.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics dns-matrix`
  - Flags:
    - `--domain <string>`: Seed domain name to probe across resolvers (default: `all`).
    - `--resolvers <string>`: Comma-separated list of target DNS resolver IPs (default: `1.1.1.1,8.8.8.8,9.9.9.9`).
    - `--timeout <duration>`: Individual UDP/DoH query timeout (default: `2s`).
    - `--is-dnssec-enforced`: Require DNSSEC validation signature verification (default: `true`).
    - `--has-doh-fallback`: Fallback to DNS-over-HTTPS if standard UDP port 53 is filtered (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedDnsResolverRecord struct {
        ResolverIP string        `json:"resolver_ip"`
        DomainName string        `json:"domain_name"`
        LocationName string      `json:"location_name"`
        MxTarget string          `json:"mx_target"`
        QueryDurationMs float64  `json:"query_duration_ms"`
        TtlSeconds int           `json:"ttl_seconds"`
        IsSpfValid bool          `json:"is_spf_valid"`
        IsDkimValid bool         `json:"is_dkim_valid"`
        IsDmarcAligned bool      `json:"is_dmarc_aligned"`
        HasDnssecValid bool      `json:"has_dnssec_valid"`
    }

    type SeedDnsResolverModel struct {
        Records []SeedDnsResolverRecord `json:"records"`
        ActiveIndex int          `json:"active_index"`
        TargetDomain string      `json:"target_domain"`
        IsProbing bool           `json:"is_probing"`
        HasFailures bool         `json:"has_failures"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedDnsResolverRecord (
        resolver_ip TEXT NOT NULL,
        domain_name TEXT NOT NULL,
        location_name TEXT NOT NULL,
        mx_target TEXT NOT NULL,
        query_duration_ms REAL NOT NULL DEFAULT 0.0,
        ttl_seconds INTEGER NOT NULL DEFAULT 300,
        is_spf_valid INTEGER NOT NULL DEFAULT 1,
        is_dkim_valid INTEGER NOT NULL DEFAULT 1,
        is_dmarc_aligned INTEGER NOT NULL DEFAULT 1,
        has_dnssec_valid INTEGER NOT NULL DEFAULT 1,
        PRIMARY KEY (resolver_ip, domain_name)
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics dns-matrix`. Executes parallel non-blocking Go goroutines querying UDP/DoH endpoints and writes consistency records to `repodb/pipeline.db`.

---

### Improvement 1634: Seed Mailbox SMTP Auth & STARTTLS Handshake Latency Probe
- **ID:** `SEED-SMTP-AUTH-PROBER-1634`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Benchmark end-to-end SMTP handshakes (EHLO, STARTTLS, AUTH PLAIN/LOGIN, MAIL FROM) against major destination MX nodes, recording TLS 1.3 cipher suite negotiation times, certificate chain validity, and handshake jitter.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SMTP AUTH & STARTTLS HANDSHAKE LATENCY PROBE] ────────────────────────────────────────┐
      │ MX TARGETS: 12 | CONCURRENT PROBES: 4 | TLS 1.3: 100% | CIPHER: AES-GCM-256 | MEAN HANDSHAKE: 34ms│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ TARGET MX HOST | PROVIDER | STARTTLS | AUTH TYPE | CERT EXPIRY | HANDSHAKE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ gmail-smtp-in.l.goo | Google | [✓] TLS1.3 | XOAUTH2 | 68 days | 22.4ms [●]
      │ outlook-com.olc.pro | Microsoft | [✓] TLS1.3 | AUTH LOGIN | 112 days | 38.1ms [●]
      │ mta5.am0.yahoodns.n | Yahoo | [✓] TLS1.2 | AUTH PLAIN | 45 days | 49.6ms [▲]
      │ mx.zoho.com | Zoho | [✓] TLS1.3 | AUTH PLAIN | 84 days | 28.9ms [●]
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: gmail-smtp-in.l.google.com] Cipher: TLS_AES_256_GCM_SHA384 | ALPN: smtp | Cert CN: *.google.com
      │ Keys: [↑/↓: Select MX] [t: Retest Handshake] [c: Inspect X.509 Certificate Chain] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML dashboard with native CSS flex grid, zero script dependencies, and pure SVG latency distribution sparks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics smtp-handshake`
  - Flags:
    - `--mx-host <string>`: Specific MX hostname to test (default: `all-primary`).
    - `--auth-method <string>`: Authentication mechanism (`login`, `plain`, `xoauth2`; default: `login`).
    - `--cert-warn-days <int>`: Threshold in days to flag upcoming TLS certificate expiration (default: `14`).
    - `--is-tls-strict`: Fail handshake if TLS version is below 1.2 (default: `true`).
    - `--has-alpn-enabled`: Negotiate ALPN SMTP token during TLS ClientHello (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedSmtpHandshakeRecord struct {
        MxHost string            `json:"mx_host"`
        ProviderName string      `json:"provider_name"`
        TlsVersion string        `json:"tls_version"`
        CipherSuite string       `json:"cipher_suite"`
        CertExpiryDays int       `json:"cert_expiry_days"`
        HandshakeDurationMs float64 `json:"handshake_duration_ms"`
        IsHandshakeSuccess bool  `json:"is_handshake_success"`
        IsCertValid bool         `json:"is_cert_valid"`
        HasAlpnSupport bool      `json:"has_alpn_support"`
        HasTlsDowngrade bool     `json:"has_tls_downgrade"`
    }

    type SeedSmtpHandshakeModel struct {
        Records []SeedSmtpHandshakeRecord `json:"records"`
        SelectedIdx int          `json:"selected_idx"`
        IsProbing bool           `json:"is_probing"`
        HasCertAlerts bool       `json:"has_cert_alerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedSmtpHandshakeRecord (
        mx_host TEXT PRIMARY KEY,
        provider_name TEXT NOT NULL,
        tls_version TEXT NOT NULL,
        cipher_suite TEXT NOT NULL,
        cert_expiry_days INTEGER NOT NULL DEFAULT 90,
        handshake_duration_ms REAL NOT NULL DEFAULT 0.0,
        is_handshake_success INTEGER NOT NULL DEFAULT 1,
        is_cert_valid INTEGER NOT NULL DEFAULT 1,
        has_alpn_support INTEGER NOT NULL DEFAULT 1,
        has_tls_downgrade INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated in `gitmap email warmup diagnostics smtp-handshake`. Connects directly over Go `crypto/tls` stack with non-blocking timeouts, saving results to `repodb/pipeline.db`.

---

### Improvement 1635: Seed Interaction & Conversational Engagement Pattern Simulator HUD
- **ID:** `SEED-ENGAGEMENT-BEHAVIOR-1635`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Audit and visualize AI-driven realistic conversation patterns between warmup seeds (read dwell times, reply latency jitter, folder moves, star/mark-as-important actions), ensuring human-like engagement profiles that satisfy ISP reputation heuristics.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED ENGAGEMENT PATTERN & INTERACTION SIMULATOR] ─────────────────────────────────────┐
      │ CONVERSATION THREADS: 84 | REPLIES GENERATED: 312 | MEAN DWELL: 46s | STAR RATE: 18.5%      │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ THREAD ID | SEED A (SENDER) | SEED B (RECEIVER) | DWELL TIME | ACTION | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ th-warm-09121 | alice@seed-alpha.io | bob@seed-beta.com | 52.4s | Reply+Star | [●] SYNCED
      │ th-warm-09122 | charlie@seed-03.net | david@seed-04.org | 38.1s | Read+Move | [●] SYNCED
      │ th-warm-09123 | elena@seed-gamma.eu | frank@seed-05.co | 14.2s | Read Only | [●] SYNCED
      │ th-warm-09124 | grace@seed-omega.us | heidi@seed-06.dev | 64.0s | Reply Sent | [●] SYNCED
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: th-warm-09121] Dynamic Persona: Enterprise DevOps | Context: Post-Mortem Review
      │ Keys: [↑/↓: Browse Threads] [e: Inject Interaction] [s: Adjust Jitter Curve] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native CSS thread timeline with SVG message flow nodes, responsive bubble dialogue containers, zero external font downloads.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics engagement-simulator`
  - Flags:
    - `--thread-limit <int>`: Number of concurrent engagement threads to visualize (default: `50`).
    - `--dwell-min <duration>`: Minimum simulated read dwell duration (default: `15s`).
    - `--dwell-max <duration>`: Maximum simulated read dwell duration (default: `90s`).
    - `--is-reply-enabled`: Simulate bidirectional thread replies between seed mailboxes (default: `true`).
    - `--has-star-action`: Simulate starring and marking messages as important (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedEngagementPatternRecord struct {
        ThreadID string          `json:"thread_id"`
        SenderMailbox string     `json:"sender_mailbox"`
        ReceiverMailbox string   `json:"receiver_mailbox"`
        DwellTimeSeconds float64 `json:"dwell_time_seconds"`
        EngagementAction string  `json:"engagement_action"`
        ContextTopic string      `json:"context_topic"`
        IsReplyGenerated bool    `json:"is_reply_generated"`
        IsMarkedImportant bool   `json:"is_marked_important"`
        HasStarApplied bool      `json:"has_star_applied"`
        HasFolderMoved bool      `json:"has_folder_moved"`
    }

    type SeedEngagementPatternModel struct {
        Threads []SeedEngagementPatternRecord `json:"threads"`
        ActiveThreadIdx int      `json:"active_thread_idx"`
        MeanDwellSeconds float64 `json:"mean_dwell_seconds"`
        IsSimulating bool        `json:"is_simulating"`
        HasActivitySpike bool    `json:"has_activity_spike"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedEngagementPatternRecord (
        thread_id TEXT PRIMARY KEY,
        sender_mailbox TEXT NOT NULL,
        receiver_mailbox TEXT NOT NULL,
        dwell_time_seconds REAL NOT NULL DEFAULT 30.0,
        engagement_action TEXT NOT NULL,
        context_topic TEXT NOT NULL,
        is_reply_generated INTEGER NOT NULL DEFAULT 1,
        is_marked_important INTEGER NOT NULL DEFAULT 1,
        has_star_applied INTEGER NOT NULL DEFAULT 1,
        has_folder_moved INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics engagement-simulator`. Interacts with `repodb/pipeline.db` to stream interaction telemetry and verify natural behavior patterns.

---

### Improvement 1636: Seed Mailbox Quarantine Automated Recovery & Re-Admission Workflow
- **ID:** `SEED-QUARANTINE-RECOVERY-1636`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Manage quarantined seed mailboxes flagged for reputation drop or spam listing, guiding them through automated multi-stage probationary recovery cycles with strict rate-limiting before re-admission into active rotation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED MAILBOX QUARANTINE & PROBATIONARY RECOVERY MANAGER] ─────────────────────────────┐
      │ QUARANTINED TOTAL: 6 | IN RECOVERY: 4 | PROBATION CLEARED (24h): 12 | FAILURE RATE: 0.0%    │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ MAILBOX ID | FLAG REASON | STAGE | PROBATION HR | HEALTH SCORE | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ seed-yahoo-091 | Spam Spike | Stage 2: Rest | 18/48h | 76.4% | [▲] PROBATION
      │ seed-gsuite-114 | Auth Flap | Stage 1: Auth | 4/24h | 88.2% | [●] RESTORING
      │ seed-m365-042 | High Deferral | Stage 3: Pacing | 36/72h | 92.0% | [●] TESTING
      │ seed-fastm-018 | Manual Flag | Stage 1: Audit | 2/12h | 81.5% | [▲] AUDIT
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: seed-yahoo-091] Reason: 3 spam placements detected | Current Hourly Pacing: 2 emails/hr
      │ Keys: [↑/↓: Select Mailbox] [a: Approve Re-Admission] [r: Reset Probation] [x: Purge Mailbox] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Interactive HTML5 quarantine dashboard with native system font styling, inline SVG progress steps, zero external scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-quarantine`
  - Flags:
    - `--mailbox-id <string>`: Filter by quarantined mailbox identifier (default: `all`).
    - `--stage <int>`: Filter by recovery probation stage (1-4; default: `0` for all).
    - `--force-admit`: Force immediate re-admission without completing probationary period (default: `false`).
    - `--is-dry-run`: Simulate stage advancement without updating live status (default: `false`).
    - `--has-auto-advance`: Automatically advance stage when health score surpasses 95% (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedQuarantineRecoveryRecord struct {
        MailboxID string         `json:"mailbox_id"`
        FlagReason string        `json:"flag_reason"`
        CurrentStage int         `json:"current_stage"`
        ProbationHoursElapsed int `json:"probation_hours_elapsed"`
        ProbationHoursTotal int  `json:"probation_hours_total"`
        HealthScore float64      `json:"health_score"`
        IsProbationActive bool   `json:"is_probation_active"`
        IsReAdmitted bool        `json:"is_re_admitted"`
        HasPassedAudit bool      `json:"has_passed_audit"`
        HasRelapseAlert bool     `json:"has_relapse_alert"`
    }

    type SeedQuarantineRecoveryModel struct {
        Records []SeedQuarantineRecoveryRecord `json:"records"`
        Cursor int               `json:"cursor"`
        QuarantinedCount int     `json:"quarantined_count"`
        IsProcessing bool        `json:"is_processing"`
        HasPendingAdmissions bool `json:"has_pending_admissions"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedQuarantineRecoveryRecord (
        mailbox_id TEXT PRIMARY KEY,
        flag_reason TEXT NOT NULL,
        current_stage INTEGER NOT NULL DEFAULT 1,
        probation_hours_elapsed INTEGER NOT NULL DEFAULT 0,
        probation_hours_total INTEGER NOT NULL DEFAULT 48,
        health_score REAL NOT NULL DEFAULT 70.0,
        is_probation_active INTEGER NOT NULL DEFAULT 1,
        is_re_admitted INTEGER NOT NULL DEFAULT 0,
        has_passed_audit INTEGER NOT NULL DEFAULT 0,
        has_relapse_alert INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `gitmap email warmup diagnostics seed-quarantine`. Updates quarantine states in `repodb/pipeline.db` and dispatches audit events to GitMap telemetry.

---

### Improvement 1637: Seed Egress IP Multi-Blocklist Radar & RBL Query Engine
- **ID:** `SEED-IP-REPUTATION-RADAR-1637`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Perform high-speed parallel checks of seed egress IP addresses against top Realtime Blackhole Lists (Spamhaus ZEN, Barracuda, Spamcop, SORBS), alerting on listings and providing instant delisting assistance links.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED EGRESS IP RBL BLOCKLIST RADAR] ──────────────────────────────────────────────────┐
      │ EGRESS IPS: 8 | RBLS MONITORED: 32 | TOTAL LISTINGS: 0 | CLEAN REPUTATION: 100%             │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ EGRESS IP ADDRESS | PROVIDER | SPAMHAUS | BARRACUDA | SPAMCOP | RBL STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ 198.51.100.24 | AWS us-east-1 | [✓] CLEAN | [✓] CLEAN | [✓] CLEAN | [●] PERFECT
      │ 198.51.100.25 | AWS us-east-1 | [✓] CLEAN | [✓] CLEAN | [✓] CLEAN | [●] PERFECT
      │ 203.0.113.88 | OVH eu-west-1 | [✓] CLEAN | [✓] CLEAN | [✓] CLEAN | [●] PERFECT
      │ 192.0.2.142 | Hetzner DE | [✓] CLEAN | [✓] CLEAN | [✓] CLEAN | [●] PERFECT
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: 198.51.100.24] PTR: mail-out-01.warmup.net | FCrDNS: VERIFIED | ASN: AS16509
      │ Keys: [↑/↓: Select IP] [r: Force RBL Rescan] [d: Delist Guidance] [l: Export CSV] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure zero-CDN CSS grid with SVG green checkmark badges and instant local printable compliance report.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics rbl-radar`
  - Flags:
    - `--ip-address <string>`: Specific egress IP address to check (default: `all-configured`).
    - `--rbl-tier <string>`: RBL checking tier (`core`, `extended`, `all`; default: `core`).
    - `--timeout <duration>`: Timeout per individual DNSBL lookup (default: `1500ms`).
    - `--is-reverse-dns-checked`: Verify Forward-Confirmed Reverse DNS (FCrDNS) mapping (default: `true`).
    - `--has-delist-links`: Generate direct one-click delisting portal URLs for any listings (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedRblRadarRecord struct {
        EgressIP string          `json:"egress_ip"`
        ProviderName string      `json:"provider_name"`
        PtrRecord string         `json:"ptr_record"`
        SpamhausStatus string    `json:"spamhaus_status"`
        BarracudaStatus string   `json:"barracuda_status"`
        SpamcopStatus string     `json:"spamcop_status"`
        TotalListings int        `json:"total_listings"`
        IsClean bool             `json:"is_clean"`
        IsPtrAligned bool        `json:"is_ptr_aligned"`
        HasDelistUrl bool        `json:"has_delist_url"`
    }

    type SeedRblRadarModel struct {
        Records []SeedRblRadarRecord `json:"records"`
        CurrentIndex int         `json:"current_index"`
        IsScanning bool          `json:"is_scanning"`
        HasListedIPs bool        `json:"has_listed_ips"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedRblRadarRecord (
        egress_ip TEXT PRIMARY KEY,
        provider_name TEXT NOT NULL,
        ptr_record TEXT NOT NULL,
        spamhaus_status TEXT NOT NULL DEFAULT 'CLEAN',
        barracuda_status TEXT NOT NULL DEFAULT 'CLEAN',
        spamcop_status TEXT NOT NULL DEFAULT 'CLEAN',
        total_listings INTEGER NOT NULL DEFAULT 0,
        is_clean INTEGER NOT NULL DEFAULT 1,
        is_ptr_aligned INTEGER NOT NULL DEFAULT 1,
        has_delist_url INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics rbl-radar`. Executes concurrent DNSBL A-record queries and caches verified IP reputation records in `repodb/pipeline.db`.

---

### Improvement 1638: Seed Account OAuth2 Token Refresh & Credential Vault Sentinel
- **ID:** `SEED-CREDENTIAL-VAULT-1638`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Monitor seed mailbox authentication credential health, tracking OAuth2 refresh token expirations, app password rotations, and invalid grant errors before warmup sessions encounter auth rejections.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED OAUTH2 TOKEN REFRESH & CREDENTIAL VAULT SENTINEL] ───────────────────────────────┐
      │ VAULT ACCOUNTS: 64 | OAUTH2: 48 | APP PASSWORDS: 16 | EXPIRING SOON (<24h): 1 | INVALID: 0  │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ ACCOUNT EMAIL | AUTH TYPE | TOKEN EXPIRY | LAST REFRESH | SCOPES OK | VAULT STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ seed-g01@warm.dev | OAuth2 | 45 mins | 15m ago | [✓] FULL | [●] HEALTHY
      │ seed-m02@corp.io | OAuth2 | 18 hours | 6h ago | [✓] FULL | [●] HEALTHY
      │ seed-y03@yahoo.com | App Pass | Never (Perm) | 30d ago | [✓] IMAP | [●] HEALTHY
      │ seed-z04@zoho.com | OAuth2 | 2 hours | 58m ago | [!] REAUTH | [▲] EXPIRING
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: seed-z04@zoho.com] Token Expires in 2h 14m | Auto-Refresh: TRIGGERED | Scopes: Mail.ReadWrite, Mail.Send
      │ Keys: [↑/↓: Select Account] [r: Force Token Refresh] [v: Rotate App Password] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Zero-external asset HTML credential audit report with masked secret fields, native CSS badge pills, and local encrypted vault export.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-credentials`
  - Flags:
    - `--account-email <string>`: Filter by specific seed account email (default: `all`).
    - `--warn-hours <int>`: Threshold in hours to warn of upcoming token expiration (default: `12`).
    - `--auto-refresh`: Automatically invoke OAuth2 token refresh flow (default: `true`).
    - `--is-dry-run`: Check credential validity without mutating refresh tokens (default: `false`).
    - `--has-keyring-backend`: Store decrypted tokens strictly in native OS keyring (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedCredentialSentinelRecord struct {
        AccountEmail string      `json:"account_email"`
        AuthType string          `json:"auth_type"`
        ExpiresInSeconds int64   `json:"expires_in_seconds"`
        LastRefreshedAt int64    `json:"last_refreshed_at"`
        GrantedScopes string     `json:"granted_scopes"`
        IsValid bool             `json:"is_valid"`
        IsExpiringSoon bool      `json:"is_expiring_soon"`
        HasRefreshError bool     `json:"has_refresh_error"`
        HasSecureStorage bool    `json:"has_secure_storage"`
    }

    type SeedCredentialSentinelModel struct {
        Records []SeedCredentialSentinelRecord `json:"records"`
        SelectedIndex int        `json:"selected_index"`
        IsRefreshing bool        `json:"is_refreshing"`
        HasExpiringTokens bool   `json:"has_expiring_tokens"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedCredentialSentinelRecord (
        account_email TEXT PRIMARY KEY,
        auth_type TEXT NOT NULL,
        expires_in_seconds INTEGER NOT NULL DEFAULT 3600,
        last_refreshed_at INTEGER NOT NULL,
        granted_scopes TEXT NOT NULL,
        is_valid INTEGER NOT NULL DEFAULT 1,
        is_expiring_soon INTEGER NOT NULL DEFAULT 0,
        has_refresh_error INTEGER NOT NULL DEFAULT 0,
        has_secure_storage INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics seed-credentials`. Bridges with encrypted local SQLite vault or OS keyring and logs credential health into `repodb/pipeline.db`.

---

### Improvement 1639: Seed Pool Placement Outlier & Heuristic Anomaly Detector
- **ID:** `SEED-ANOMALY-DETECTION-1639`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Utilize statistical Z-score algorithms across inbox placement rates and reply latencies to isolate misbehaving seed mailboxes, detecting provider shadow-bans or localized ISP filter updates before pool-wide degradation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED POOL ANOMALY & OUTLIER HEURISTIC DETECTOR] ──────────────────────────────────────┐
      │ POOLS EVALUATED: 12 | OUTLIERS DETECTED: 2 | CONFIDENCE: 98.6% | Z-SCORE THRESHOLD: > 2.5   │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ SEED IDENTIFIER | PROVIDER | METRIC TESTED | VALUE | Z-SCORE | ANOMALY
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ seed-yahoo-044 | Yahoo Mail | Spam Folder % | 14.2% | +3.18 | [!] OUTLIER
      │ seed-m365-019 | Microsoft 365 | Latency RTT | 340ms | +2.84 | [!] OUTLIER
      │ seed-gsuite-102 | Google Suite | Inbox Rate | 99.1% | -0.12 | [●] NORMAL
      │ seed-fastm-003 | Fastmail | Inbox Rate | 98.9% | -0.24 | [●] NORMAL
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: seed-yahoo-044] High Spam Placements: Z-Score +3.18 exceeds 2.5 limit | Action: Auto-Quarantine Enqueued
      │ Keys: [↑/↓: Select Anomaly] [q: Quarantine Outlier] [i: Ignore/Baseline] [s: Adjust Z-Score] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Zero-external asset HTML report featuring embedded SVG scatter plots with standard deviation bounding bands, pure local CSS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-anomalies`
  - Flags:
    - `--z-threshold <float>`: Z-score limit to trigger statistical outlier alert (default: `2.5`).
    - `--sample-window <duration>`: Time window for historical metric baselining (default: `24h`).
    - `--metric <string>`: Target metric to analyze (`inbox`, `spam`, `latency`, `all`; default: `all`).
    - `--is-auto-quarantine`: Automatically quarantine outliers exceeding threshold (default: `true`).
    - `--has-baseline-learning`: Incorporate validated healthy runs into rolling baseline (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedAnomalyDetectorRecord struct {
        SeedIdentifier string    `json:"seed_identifier"`
        ProviderName string      `json:"provider_name"`
        MetricName string        `json:"metric_name"`
        ObservedValue float64    `json:"observed_value"`
        BaselineMean float64     `json:"baseline_mean"`
        ZScore float64           `json:"z_score"`
        IsAnomaly bool           `json:"is_anomaly"`
        IsQuarantined bool       `json:"is_quarantined"`
        HasConfidenceHigh bool   `json:"has_confidence_high"`
        HasAlertDispatched bool  `json:"has_alert_dispatched"`
    }

    type SeedAnomalyDetectorModel struct {
        Records []SeedAnomalyDetectorRecord `json:"records"`
        SelectedIdx int          `json:"selected_idx"`
        ZThreshold float64       `json:"z_threshold"`
        IsAnalyzing bool         `json:"is_analyzing"`
        HasSevereOutliers bool   `json:"has_severe_outliers"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedAnomalyDetectorRecord (
        seed_identifier TEXT NOT NULL,
        provider_name TEXT NOT NULL,
        metric_name TEXT NOT NULL,
        observed_value REAL NOT NULL,
        baseline_mean REAL NOT NULL,
        z_score REAL NOT NULL,
        is_anomaly INTEGER NOT NULL DEFAULT 1,
        is_quarantined INTEGER NOT NULL DEFAULT 0,
        has_confidence_high INTEGER NOT NULL DEFAULT 1,
        has_alert_dispatched INTEGER NOT NULL DEFAULT 1,
        PRIMARY KEY (seed_identifier, metric_name)
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics seed-anomalies`. Calculates rolling variance in-memory, writes anomaly triggers to `repodb/pipeline.db`, and alerts GitMap supervisor.

---

### Improvement 1640: Comprehensive Seed Diagnostics Self-Test & Diagnostic Export Bundle
- **ID:** `SEED-DIAGNOSTIC-BUNDLE-1640`
- **Cluster:** 164 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Execute an end-to-end synthetic seed pool validation suite (DNS, IMAP, SMTP, latency, inbox placement, headers) and generate a compressed, cryptographically signed diagnostic bundle for root cause analysis.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [COMPREHENSIVE SEED DIAGNOSTICS SUITE & BUNDLE EXPORTER] ──────────────────────────────┐
      │ TEST SUITES: 6 | TOTAL CHECKS: 48 | PASSED: 47 | WARN: 1 | FAIL: 0 | DURATION: 1.84s        │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ CHECK SUITE | TARGET POOL | CHECKS RUN | PASSED | LATENCY | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ DNS & Nameservers | Global Edges | 8 checks | 8 pass | 14ms | [✓] PASS
      │ IMAP Drain & Rescue | Primary Seeds | 12 checks | 12 pass | 48ms | [✓] PASS
      │ SMTP Handshakes | Core ESPs | 10 checks | 10 pass | 32ms | [✓] PASS
      │ Placement Heuristic | Yahoo Cohort | 8 checks | 7 pass | 110ms | [!] WARN
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [BUNDLE READY] Relative Path: .ai-memory/diagnostics/seed-bundle-20261028-1640.zip | Size: 1.2 MB
      │ Keys: [↑/↓: Browse Check Suites] [b: Generate Diagnostic Bundle] [v: Verify Signature] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Zero-external asset HTML report with native CSS test run matrix, SVG check badges, and browser download trigger for the local zip bundle.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-bundle`
  - Flags:
    - `--output-dir <string>`: Target relative path for output bundle (default: `.ai-memory/diagnostics/`).
    - `--include-raw-headers`: Include anonymized raw RFC 822 email headers in diagnostic dump (default: `false`).
    - `--is-full-suite`: Execute intensive stress checks across entire seed fleet (default: `true`).
    - `--has-anonymization`: Scrub user email prefixes before bundling (default: `true`).
    - `--has-sha256-signature`: Generate SHA256 checksum file alongside the bundle (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedDiagnosticBundleRecord struct {
        BundleID string          `json:"bundle_id"`
        BundlePath string        `json:"bundle_path"`
        TestCount int            `json:"test_count"`
        PassedTests int          `json:"passed_tests"`
        FailedTests int          `json:"failed_tests"`
        GeneratedAt int64        `json:"generated_at"`
        Sha256Checksum string    `json:"sha256_checksum"`
        IsComplete bool          `json:"is_complete"`
        HasSevereFailures bool   `json:"has_severe_failures"`
        HasSignatureValid bool   `json:"has_signature_valid"`
    }

    type SeedDiagnosticBundleModel struct {
        Bundle SeedDiagnosticBundleRecord `json:"bundle"`
        ActiveTab int            `json:"active_tab"`
        IsExporting bool         `json:"is_exporting"`
        HasGenerated bool        `json:"has_generated"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedDiagnosticBundleRecord (
        bundle_id TEXT PRIMARY KEY,
        bundle_path TEXT NOT NULL,
        test_count INTEGER NOT NULL DEFAULT 0,
        passed_tests INTEGER NOT NULL DEFAULT 0,
        failed_tests INTEGER NOT NULL DEFAULT 0,
        generated_at INTEGER NOT NULL,
        sha256_checksum TEXT NOT NULL,
        is_complete INTEGER NOT NULL DEFAULT 1,
        has_severe_failures INTEGER NOT NULL DEFAULT 0,
        has_signature_valid INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics seed-bundle`. Stores output artifacts under relative repository path `.ai-memory/diagnostics/` and records run summary in `repodb/pipeline.db`.

---

# Cluster 165: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 1641–1650)

---

### Improvement 1641: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor Controller
- **ID:** `GOV-CONCURRENCY-CONTROLLER-1641`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Regulate outbound SMTP connection pools dynamically using inverted backoff models, adjusting simultaneous worker slots per destination MX domain based on real-time 421/451 deferral feedback and server load signals.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [MULTI-DOMAIN CONCURRENCY GOVERNOR & BACKOFF CONTROLLER] ──────────────────────────────┐
      │ GLOBAL SLOTS: 128/256 | ACTIVE POOLS: 36 | MEAN DEFERRAL: 0.12% | GOVERNOR: STEADY          │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ DESTINATION DOMAIN | MAX SLOTS | ACTIVE | BACKOFF MULT | DEFERRAL % | STATE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ gmail.com | 64 | 42 | 1.0x (NONE) | 0.0% | [●] OPTIMAL
      │ outlook.com | 32 | 18 | 1.1x (LIGHT) | 0.2% | [●] OPTIMAL
      │ yahoo.com | 16 | 6 | 2.4x (MED) | 2.8% | [▲] THROTTLE
      │ icloud.com | 12 | 4 | 1.0x (NONE) | 0.0% | [●] OPTIMAL
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: yahoo.com] Slots throttled down (32 -> 16) due to 421 TS01 temporary rate limiting
      │ Keys: [↑/↓: Select Domain] [+: Manual Slot Boost] [-: Force Throttle] [r: Reset Governor] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML dashboard with native CSS flex grid, zero script dependencies, and pure SVG latency distribution sparks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics governor-status`
  - Flags:
    - `--domain <string>`: Filter governor metrics by destination domain (default: `all`).
    - `--global-max-slots <int>`: Upper bound for total concurrent outbound connections (default: `256`).
    - `--deferral-window <duration>`: Rolling window for deferral calculation (default: `5m`).
    - `--is-auto-throttle`: Automatically reduce slot capacity upon receiving 4xx responses (default: `true`).
    - `--has-inverted-recovery`: Accelerate capacity restoration exponentially as error rates reach zero (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorDomainConcurrencyRecord struct {
        DestinationDomain string `json:"destination_domain"`
        MaxSlots int             `json:"max_slots"`
        ActiveSlots int          `json:"active_slots"`
        BackoffMultiplier float64 `json:"backoff_multiplier"`
        DeferralRatePercent float64 `json:"deferral_rate_percent"`
        ConsecutiveSuccesses int `json:"consecutive_successes"`
        IsThrottled bool         `json:"is_throttled"`
        IsOptimal bool           `json:"is_optimal"`
        HasInversionActive bool  `json:"has_inversion_active"`
        HasRecentDeferral bool   `json:"has_recent_deferral"`
    }

    type GovernorDomainConcurrencyModel struct {
        Domains []GovernorDomainConcurrencyRecord `json:"domains"`
        Cursor int               `json:"cursor"`
        TotalActiveSlots int     `json:"total_active_slots"`
        IsGoverning bool         `json:"is_governing"`
        HasThrottledDomains bool `json:"has_throttled_domains"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorDomainConcurrencyRecord (
        destination_domain TEXT PRIMARY KEY,
        max_slots INTEGER NOT NULL DEFAULT 32,
        active_slots INTEGER NOT NULL DEFAULT 0,
        backoff_multiplier REAL NOT NULL DEFAULT 1.0,
        deferral_rate_percent REAL NOT NULL DEFAULT 0.0,
        consecutive_successes INTEGER NOT NULL DEFAULT 0,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_inversion_active INTEGER NOT NULL DEFAULT 1,
        has_recent_deferral INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics governor-status`. Queries state from `repodb/pipeline.db` and modulates outbound queue dispatch in real time.

---

### Improvement 1642: SMTP Deferral Code Heuristic Classifier & Pattern Matcher
- **ID:** `GOV-DEFERRAL-CLASSIFIER-1642`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Parse and categorize granular SMTP temporary error responses (421, 450, 451, 452) into structured classification buckets (rate limit, IP reputation, greylisting, mailbox full, TLS requirement) to trigger customized recovery actions.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SMTP 4XX DEFERRAL HEURISTIC CLASSIFIER] ──────────────────────────────────────────────┐
      │ PARSED RESPONSES (1h): 248 | CLASSIFIED: 100% | GREYLIST: 42% | RATE LIMIT: 38% | UNKNOWN: 0%│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ RAW CODE/SUBCODE | CATEGORY | SEVERITY | SAMPLE REASON | ACTION
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ 421 4.7.0 | IP Reputation | HIGH | Our system has detected an unusual... | Pause 15m
      │ 451 4.7.1 | Greylisting | LOW | Recipient temporarily unavailable... | Retry in 5m
      │ 452 4.2.2 | Mailbox Full | MEDIUM | The user mailbox has exceeded quo... | Drop Batch
      │ 421 4.4.5 | Rate Limit | MEDIUM | Server busy, too many connections... | Cut Slots
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: 421 4.7.0 IP Reputation] Matched Pattern: RFC 3463 / RFC 5321 | Automatic Action: Governor Throttle Domain
      │ Keys: [↑/↓: Browse Classes] [t: Test Custom Error String] [a: Edit Action Map] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML deferral taxonomy tree using native CSS grid, zero scripts, and system monospace typography.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics deferral-classifier`
  - Flags:
    - `--sample-text <string>`: Raw SMTP banner string to test against classification heuristics.
    - `--category <string>`: Filter classifications by category (`reputation`, `greylist`, `ratelimit`, `all`; default: `all`).
    - `--is-regex-optimized`: Compile deterministic regex DFA tables for sub-millisecond classification (default: `true`).
    - `--has-rfc3463-matching`: Strictly parse enhanced mail system status codes (RFC 3463; default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorDeferralClassifierRecord struct {
        StatusCode string        `json:"status_code"`
        Category string          `json:"category"`
        SeverityLevel string     `json:"severity_level"`
        PatternRegex string      `json:"pattern_regex"`
        RecommendedAction string `json:"recommended_action"`
        MatchCount int           `json:"match_count"`
        IsHighSeverity bool      `json:"is_high_severity"`
        IsActionAutomated bool   `json:"is_action_automated"`
        HasEnhancedCode bool     `json:"has_enhanced_code"`
    }

    type GovernorDeferralClassifierModel struct {
        Records []GovernorDeferralClassifierRecord `json:"records"`
        Cursor int               `json:"cursor"`
        IsMatching bool          `json:"is_matching"`
        HasHighSeverityRules bool `json:"has_high_severity_rules"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorDeferralClassifierRecord (
        status_code TEXT NOT NULL,
        category TEXT NOT NULL,
        severity_level TEXT NOT NULL,
        pattern_regex TEXT NOT NULL,
        recommended_action TEXT NOT NULL,
        match_count INTEGER NOT NULL DEFAULT 0,
        is_high_severity INTEGER NOT NULL DEFAULT 0,
        is_action_automated INTEGER NOT NULL DEFAULT 1,
        has_enhanced_code INTEGER NOT NULL DEFAULT 1,
        PRIMARY KEY (status_code, category)
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated in `gitmap email warmup diagnostics deferral-classifier`. Parses server replies directly from Go SMTP client hooks and persists match frequency in `repodb/pipeline.db`.

---

### Improvement 1643: Decorrelated Full Jitter Backoff Curve Generator & Visualizer
- **ID:** `GOV-EXPONENTIAL-JITTER-1643`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Calculate and visually plot decorrelated full-jitter retry curves (`Sleep = min(Cap, random_between(Base, Sleep * 3))`) for deferred warmup deliveries, preventing synchronization thundering herds against recipient MX clusters.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [DECORRELATED JITTER RETRY CURVE GENERATOR & SIMULATOR] ───────────────────────────────┐
      │ BASE SLEEP: 2.0s | CAP SLEEP: 300.0s | JITTER FORMULA: DECORRELATED | SAMPLES: 1,000        │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ RETRY ATTEMPT | MIN SLEEP | MAX SLEEP | MEAN JITTER | SPREAD (σ) | DISTRIBUTION
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Attempt #1 | 2.0s | 6.0s | 3.8s | 1.15s | [===>------]
      │ Attempt #2 | 2.0s | 18.0s | 9.4s | 4.21s | [=====>----]
      │ Attempt #3 | 2.0s | 54.0s | 26.1s | 12.80s | [=======>--]
      │ Attempt #4 (Capped) | 2.0s | 162.0s | 78.5s | 38.40s | [=========>]
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SIMULATION] 1000 Simulated Retries: Peak concurrency spread over 284 seconds | 0 Thundering Herds
      │ Keys: [↑/↓: Change Attempt] [b: Edit Base Delay] [c: Edit Cap] [s: Resimulate 1000 Runs] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native HTML5 canvas/SVG probability density curve with zero CDN libraries, embedded CSS styling.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics jitter-curve`
  - Flags:
    - `--base-delay <duration>`: Initial base backoff sleep duration (default: `2s`).
    - `--cap-delay <duration>`: Upper maximum backoff sleep ceiling (default: `300s`).
    - `--algorithm <string>`: Jitter algorithm (`decorrelated`, `full-jitter`, `equal-jitter`; default: `decorrelated`).
    - `--is-simulation-mode`: Run statistical Monte Carlo simulation without dispatching live emails (default: `true`).
    - `--has-entropy-source`: Use crypto/rand hardware entropy for jitter generation (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorJitterCurveRecord struct {
        AttemptNumber int        `json:"attempt_number"`
        MinSleepSeconds float64  `json:"min_sleep_seconds"`
        MaxSleepSeconds float64  `json:"max_sleep_seconds"`
        MeanSleepSeconds float64 `json:"mean_sleep_seconds"`
        StandardDeviation float64 `json:"standard_deviation"`
        AlgorithmName string     `json:"algorithm_name"`
        IsCapped bool            `json:"is_capped"`
        IsOptimalSpread bool     `json:"is_optimal_spread"`
        HasHardwareEntropy bool  `json:"has_hardware_entropy"`
    }

    type GovernorJitterCurveModel struct {
        Attempts []GovernorJitterCurveRecord `json:"attempts"`
        SelectedAttempt int      `json:"selected_attempt"`
        BaseSeconds float64      `json:"base_seconds"`
        CapSeconds float64       `json:"cap_seconds"`
        IsSimulating bool        `json:"is_simulating"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorJitterCurveRecord (
        attempt_number INTEGER PRIMARY KEY,
        min_sleep_seconds REAL NOT NULL,
        max_sleep_seconds REAL NOT NULL,
        mean_sleep_seconds REAL NOT NULL,
        standard_deviation REAL NOT NULL,
        algorithm_name TEXT NOT NULL,
        is_capped INTEGER NOT NULL DEFAULT 0,
        is_optimal_spread INTEGER NOT NULL DEFAULT 1,
        has_hardware_entropy INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics jitter-curve`. Provides delay values to the dispatch queue and logs jitter evaluations in `repodb/pipeline.db`.

---

### Improvement 1644: Hierarchical Token Bucket Rate Limiter Matrix (Domain & IP Level)
- **ID:** `GOV-RATE-LIMIT-TOKEN-BUCKET-1644`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Enforce nested dual-layer token bucket rate limiting (global per-IP egress rate + granular destination domain quota), visualizing token refilling rates, burst allowances, and queue starvation metrics in real time.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [HIERARCHICAL TOKEN BUCKET RATE LIMITER MATRIX] ───────────────────────────────────────┐
      │ IP BUCKETS: 4 | DOMAIN BUCKETS: 42 | TOTAL TOKENS: 1,840/2,000 | BURST HEADROOM: 16%        │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ BUCKET TARGET | CAPACITY | CURRENT | REFILL RATE | BURST MAX | SATURATION
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ IP: 198.51.100.24 | 1000 tok | 920 tok | 50.0 tok/s | 150 tok | [=>--------]  8%
      │ Domain: gmail.com | 500 tok | 440 tok | 25.0 tok/s | 50 tok | [==>-------] 12%
      │ Domain: outlook.com | 250 tok | 210 tok | 10.0 tok/s | 25 tok | [==>-------] 16%
      │ Domain: yahoo.com | 100 tok | 34 tok | 2.0 tok/s | 10 tok | [=======>--] 66%
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: Domain: yahoo.com] Tokens: 34/100 | Refill: 2.0/s | High Saturation (66%) triggers governor pacing
      │ Keys: [↑/↓: Select Bucket] [+: Add Emergency Tokens] [r: Set Refill Rate] [c: Clear Queue] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure zero-CDN CSS progress bars with live SVG bucket visualization, zero external scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics token-buckets`
  - Flags:
    - `--target <string>`: Specific IP or domain bucket to monitor (default: `all`).
    - `--refill-override <float>`: Temporarily override token refill rate per second.
    - `--burst-allowance <int>`: Max allowable burst token capacity over standard quota (default: `20`).
    - `--is-hierarchical`: Evaluate both IP and Domain buckets before releasing dispatch token (default: `true`).
    - `--has-atomic-sync`: Use sync/atomic CPU instructions for lock-free token subtraction (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorTokenBucketRecord struct {
        BucketKey string         `json:"bucket_key"`
        BucketType string        `json:"bucket_type"`
        MaxCapacity int          `json:"max_capacity"`
        CurrentTokens int        `json:"current_tokens"`
        RefillRatePerSec float64 `json:"refill_rate_per_sec"`
        BurstLimit int           `json:"burst_limit"`
        SaturationPercent float64 `json:"saturation_percent"`
        IsSaturated bool         `json:"is_saturated"`
        IsRefillActive bool      `json:"is_refill_active"`
        HasBurstAvailable bool   `json:"has_burst_available"`
    }

    type GovernorTokenBucketModel struct {
        Buckets []GovernorTokenBucketRecord `json:"buckets"`
        SelectedIndex int        `json:"selected_index"`
        IsMonitoring bool        `json:"is_monitoring"`
        HasSaturatedBuckets bool `json:"has_saturated_buckets"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorTokenBucketRecord (
        bucket_key TEXT PRIMARY KEY,
        bucket_type TEXT NOT NULL,
        max_capacity INTEGER NOT NULL DEFAULT 500,
        current_tokens INTEGER NOT NULL DEFAULT 500,
        refill_rate_per_sec REAL NOT NULL DEFAULT 20.0,
        burst_limit INTEGER NOT NULL DEFAULT 50,
        saturation_percent REAL NOT NULL DEFAULT 0.0,
        is_saturated INTEGER NOT NULL DEFAULT 0,
        is_refill_active INTEGER NOT NULL DEFAULT 1,
        has_burst_available INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics token-buckets`. Embedded directly in the outbound delivery worker loop with state cached in memory and logged to `repodb/pipeline.db`.

---

### Improvement 1645: Per-Domain Delivery Circuit Breaker & Trip Sentinel HUD
- **ID:** `GOV-CIRCUIT-BREAKER-1645`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Implement finite-state circuit breakers (CLOSED, OPEN, HALF-OPEN) per destination domain, cutting off delivery immediately upon cascading 5xx rejections, and safely evaluating canary deliveries before re-closing the circuit.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [PER-DOMAIN DELIVERY CIRCUIT BREAKER & TRIP SENTINEL] ─────────────────────────────────┐
      │ CIRCUITS CLOSED (HEALTHY): 40 | HALF-OPEN: 1 | OPEN (TRIPPED): 1 | CANARY SUCCESS: 100%     │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ DESTINATION DOMAIN | CIRCUIT STATE | CONSEC FAIL | TRIP THRESHOLD | RESET TIMER | REPUTATION
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ gmail.com | [●] CLOSED | 0 | 5 failures | N/A | 100% [●]
      │ outlook.com | [●] CLOSED | 0 | 5 failures | N/A | 99%  [●]
      │ yahoo.com | [▲] HALF-OPEN | 0 (Canary) | 3 failures | 12s remain | 84%  [▲]
      │ corporate-mx.ru | [■] OPEN | 5 | 5 failures | 4m 18s | 42%  [■]
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: corporate-mx.ru] Tripped: 5 consecutive 550 Blocked errors | Status: OPEN | Deliveries Quenched
      │ Keys: [↑/↓: Select Domain] [f: Force Trip Circuit] [r: Manual Circuit Reset] [c: Dispatch Canary] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML circuit dashboard using native CSS states, SVG breaker toggle graphics, and zero external JS/CSS dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics circuit-breaker`
  - Flags:
    - `--domain <string>`: Destination domain to inspect or manipulate (default: `all`).
    - `--trip-threshold <int>`: Number of consecutive failures required to trip circuit open (default: `5`).
    - `--reset-timeout <duration>`: Duration circuit stays OPEN before entering HALF-OPEN canary mode (default: `5m`).
    - `--is-canary-active`: Dispatch small single canary probes when in HALF-OPEN state (default: `true`).
    - `--has-auto-recovery`: Automatically restore CLOSED status upon 3 consecutive canary successes (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorCircuitBreakerRecord struct {
        DestinationDomain string `json:"destination_domain"`
        CircuitState string      `json:"circuit_state"`
        ConsecutiveFailures int  `json:"consecutive_failures"`
        TripThreshold int        `json:"trip_threshold"`
        ResetTimeoutSeconds int  `json:"reset_timeout_seconds"`
        RemainingCooldownSeconds int `json:"remaining_cooldown_seconds"`
        CanarySuccessCount int   `json:"canary_success_count"`
        IsClosed bool            `json:"is_closed"`
        IsOpen bool              `json:"is_open"`
        HasCanaryPassed bool     `json:"has_canary_passed"`
    }

    type GovernorCircuitBreakerModel struct {
        Circuits []GovernorCircuitBreakerRecord `json:"circuits"`
        SelectedIdx int          `json:"selected_idx"`
        IsMonitoring bool        `json:"is_monitoring"`
        HasTrippedCircuits bool  `json:"has_tripped_circuits"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorCircuitBreakerRecord (
        destination_domain TEXT PRIMARY KEY,
        circuit_state TEXT NOT NULL DEFAULT 'CLOSED',
        consecutive_failures INTEGER NOT NULL DEFAULT 0,
        trip_threshold INTEGER NOT NULL DEFAULT 5,
        reset_timeout_seconds INTEGER NOT NULL DEFAULT 300,
        remaining_cooldown_seconds INTEGER NOT NULL DEFAULT 0,
        canary_success_count INTEGER NOT NULL DEFAULT 0,
        is_closed INTEGER NOT NULL DEFAULT 1,
        is_open INTEGER NOT NULL DEFAULT 0,
        has_canary_passed INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics circuit-breaker`. Guards dispatch goroutines and persists circuit transition audit logs into `repodb/pipeline.db`.

---

### Improvement 1646: MX Host Resolution In-Memory Cache & Ephemeral Pool Manager
- **ID:** `GOV-DNS-MX-CACHE-1646`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Provide high-performance, concurrent lock-free MX DNS resolution caching with automatic TTL eviction, handling multi-MX priority failover and Anycast IP rotation without triggering external DNS lookups for each email.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [MX HOST RESOLUTION CACHE & FAILOVER POOL] ────────────────────────────────────────────┐
      │ CACHED DOMAINS: 124 | CACHE HIT RATE: 99.4% | MEAN RESOLVE: 0.12ms | ACTIVE FAILOVERS: 0    │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ DOMAIN NAME | PRIMARY MX HOST | PREF | RESOLVED IP | TTL REMAIN | HEALTH
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ gmail.com | gmail-smtp-in.l.goo | 5 | 108.177.14.26 | 240s | [●] ACTIVE
      │ outlook.com | outlook-com.olc.pro | 10 | 104.47.58.33 | 180s | [●] ACTIVE
      │ yahoo.com | mta5.am0.yahoodns.n | 10 | 67.195.228.106 | 90s | [●] ACTIVE
      │ apple.com | mx1.mail.icloud.com | 10 | 17.42.251.41 | 210s | [●] ACTIVE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: gmail.com] Primary: gmail-smtp-in.l.google.com (pref 5) | Backup: alt1.gmail-smtp-in.l.google.com (pref 10)
      │ Keys: [↑/↓: Select Domain] [f: Force MX Invalidation] [p: Test MX Failover] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native CSS monospace table with embedded SVG latency indicators, local JSON export, zero CDN dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics mx-cache`
  - Flags:
    - `--domain <string>`: Domain to query in MX cache (default: `all`).
    - `--flush-domain <string>`: Evict specified domain from active resolution cache.
    - `--is-failover-tested`: Perform dry-run DNS resolution to secondary backup MX targets (default: `false`).
    - `--has-prefetch-enabled`: Asynchronously refresh MX records 30s before TTL expiration (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorMxCacheRecord struct {
        DomainName string        `json:"domain_name"`
        PrimaryMxHost string     `json:"primary_mx_host"`
        Preference int           `json:"preference"`
        ResolvedIP string        `json:"resolved_ip"`
        TtlSecondsRemaining int  `json:"ttl_seconds_remaining"`
        HitCount int64           `json:"hit_count"`
        IsPrimaryActive bool     `json:"is_primary_active"`
        IsCached bool            `json:"is_cached"`
        HasBackupAvailable bool  `json:"has_backup_available"`
        HasResolutionError bool  `json:"has_resolution_error"`
    }

    type GovernorMxCacheModel struct {
        Records []GovernorMxCacheRecord `json:"records"`
        Cursor int               `json:"cursor"`
        HitRatePercent float64   `json:"hit_rate_percent"`
        IsPrefetching bool       `json:"is_prefetching"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorMxCacheRecord (
        domain_name TEXT PRIMARY KEY,
        primary_mx_host TEXT NOT NULL,
        preference INTEGER NOT NULL DEFAULT 10,
        resolved_ip TEXT NOT NULL,
        ttl_seconds_remaining INTEGER NOT NULL DEFAULT 300,
        hit_count INTEGER NOT NULL DEFAULT 0,
        is_primary_active INTEGER NOT NULL DEFAULT 1,
        is_cached INTEGER NOT NULL DEFAULT 1,
        has_backup_available INTEGER NOT NULL DEFAULT 1,
        has_resolution_error INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics mx-cache`. Integrated with the in-process DNS resolver cache and synchronized to `repodb/pipeline.db`.

---

### Improvement 1647: Persistent SMTP Connection Pool & Keep-Alive Lifecycle Sentinel
- **ID:** `GOV-CONNECTION-POOL-REUSE-1647`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Maintain pooled, authenticated SMTP client connections for high-volume warmup sessions with RFC 5321 RSET pipeline reuse, idle timeout eviction, and maximum message dispatch thresholds per socket.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [PERSISTENT SMTP CONNECTION POOL & LIFECYCLE SENTINEL] ────────────────────────────────┐
      │ POOLED SOCKETS: 48 | IDLE: 14 | IN-FLIGHT: 34 | REUSE RATIO: 8.4 msgs/conn | RSET PASS: 100%│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ SOCKET ID | TARGET MX HOST | MSGS SENT | IDLE TIME | AGE (s) | STATE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ sock-mta-0041 | gmail-smtp-in.l.goo | 14/25 | 2.1s | 48s | [●] IN-FLIGHT
      │ sock-mta-0042 | gmail-smtp-in.l.goo | 8/25 | 4.8s | 22s | [●] IDLE
      │ sock-mta-0043 | outlook-com.olc.pro | 19/25 | 1.2s | 74s | [●] IN-FLIGHT
      │ sock-mta-0044 | mta5.am0.yahoodns.n | 25/25 | 0.0s | 92s | [▲] DRAINING
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: sock-mta-0044] Max msgs reached (25/25) | Issuing QUIT and closing socket gracefully
      │ Keys: [↑/↓: Select Socket] [d: Drain Connection] [k: Kill Socket] [r: Send RSET Probe] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML dashboard with native CSS flex grid, zero script dependencies, and pure SVG latency distribution sparks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics connection-pool`
  - Flags:
    - `--max-idle <duration>`: Max duration an idle SMTP socket remains pooled (default: `30s`).
    - `--max-messages-per-conn <int>`: Max warmup messages dispatched through a single persistent connection (default: `25`).
    - `--is-pipelining-verified`: Confirm remote MX advertises SMTP PIPELINING capability (RFC 2920; default: `true`).
    - `--has-keepalive-probes`: Send background NOOP/RSET keep-alive ping every 10 seconds (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorConnectionPoolRecord struct {
        SocketID string          `json:"socket_id"`
        TargetMxHost string      `json:"target_mx_host"`
        MessagesDispatched int   `json:"messages_dispatched"`
        MaxMessagesPerConn int   `json:"max_messages_per_conn"`
        IdleTimeSeconds float64  `json:"idle_time_seconds"`
        SocketAgeSeconds int     `json:"socket_age_seconds"`
        SocketState string       `json:"socket_state"`
        IsPooled bool            `json:"is_pooled"`
        IsDraining bool          `json:"is_draining"`
        HasPipelining bool       `json:"has_pipelining"`
    }

    type GovernorConnectionPoolModel struct {
        Sockets []GovernorConnectionPoolRecord `json:"sockets"`
        ActiveCursor int         `json:"active_cursor"`
        TotalPooledCount int     `json:"total_pooled_count"`
        IsPoolingActive bool     `json:"is_pooling_active"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorConnectionPoolRecord (
        socket_id TEXT PRIMARY KEY,
        target_mx_host TEXT NOT NULL,
        messages_dispatched INTEGER NOT NULL DEFAULT 0,
        max_messages_per_conn INTEGER NOT NULL DEFAULT 25,
        idle_time_seconds REAL NOT NULL DEFAULT 0.0,
        socket_age_seconds INTEGER NOT NULL DEFAULT 0,
        socket_state TEXT NOT NULL DEFAULT 'IDLE',
        is_pooled INTEGER NOT NULL DEFAULT 1,
        is_draining INTEGER NOT NULL DEFAULT 0,
        has_pipelining INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics connection-pool`. Interacts with internal TCP connection pool and registers socket events in `repodb/pipeline.db`.

---

### Improvement 1648: Outbound Delivery Queue Pressure Load Shedder & Priority Sorter
- **ID:** `GOV-QUEUE-PRESSURE-SHEDDING-1648`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Protect warm-up orchestrators from memory exhaustion under heavy 4xx deferral spikes by implementing priority-driven load shedding, evicting low-tier probe traffic while preserving primary domain warmup schedules.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [OUTBOUND QUEUE PRESSURE & LOAD SHEDDER SENTINEL] ─────────────────────────────────────┐
      │ QUEUE OCCUPANCY: 420/5000 (8.4%) | PRIORITY 1: 100% | PRESSURE TIER: NOMINAL | SHED: 0      │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ PRIORITY TIER | ENQUEUED | DISPATCH/s | PRESSURE | EVICTION POLICY | HEALTH
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Tier 1: High Warmup | 180 msgs | 14.2 msg/s | NOMINAL | Never Shed | [●] SECURE
      │ Tier 2: Probe Sync | 140 msgs | 8.5 msg/s | NOMINAL | Shed if > 80% | [●] SECURE
      │ Tier 3: Engagement | 80 msgs | 4.0 msg/s | NOMINAL | Shed if > 60% | [●] SECURE
      │ Tier 4: Bulk Testing | 20 msgs | 1.1 msg/s | NOMINAL | Shed if > 50% | [●] SECURE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [PRESSURE STATUS: NOMINAL] Queue memory footprint: 18.4 MB | Headroom: 91.6% | No load shedding needed
      │ Keys: [↑/↓: Select Tier] [s: Trigger Dry-Run Shedding] [t: Edit Capacity Threshold] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML queue pressure gauge using pure SVG arcs, native CSS status classes, zero scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics load-shedder`
  - Flags:
    - `--max-queue-size <int>`: Maximum queue capacity before triggering shedding (default: `5000`).
    - `--shed-tier <int>`: Minimum priority tier eligible for emergency eviction (default: `3`).
    - `--is-shedding-active`: Automatically discard non-critical probe messages under memory pressure (default: `true`).
    - `--has-shed-logging`: Persist metadata of shed messages into SQLite for post-mortem audit (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorLoadShedderRecord struct {
        TierLevel int            `json:"tier_level"`
        TierName string          `json:"tier_name"`
        EnqueuedCount int        `json:"enqueued_count"`
        DispatchRatePerSec float64 `json:"dispatch_rate_per_sec"`
        PressureState string     `json:"pressure_state"`
        ShedCountTotal int       `json:"shed_count_total"`
        IsSheddingEligible bool  `json:"is_shedding_eligible"`
        IsQueueHealthy bool      `json:"is_queue_healthy"`
        HasCapacityAvailable bool `json:"has_capacity_available"`
    }

    type GovernorLoadShedderModel struct {
        Tiers []GovernorLoadShedderRecord `json:"tiers"`
        SelectedIndex int        `json:"selected_index"`
        GlobalOccupancyPercent float64 `json:"global_occupancy_percent"`
        IsMonitoring bool        `json:"is_monitoring"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorLoadShedderRecord (
        tier_level INTEGER PRIMARY KEY,
        tier_name TEXT NOT NULL,
        enqueued_count INTEGER NOT NULL DEFAULT 0,
        dispatch_rate_per_sec REAL NOT NULL DEFAULT 0.0,
        pressure_state TEXT NOT NULL DEFAULT 'NOMINAL',
        shed_count_total INTEGER NOT NULL DEFAULT 0,
        is_shedding_eligible INTEGER NOT NULL DEFAULT 0,
        is_queue_healthy INTEGER NOT NULL DEFAULT 1,
        has_capacity_available INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics load-shedder`. Controls channel buffer thresholds in the queue worker pool and records drops in `repodb/pipeline.db`.

---

### Improvement 1649: Governor Dispatch Velocity & ISP Quota Telemetry Analytics HUD
- **ID:** `GOV-THROUGHPUT-ANALYTICS-1649`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Aggregate real-time outbound delivery throughput (emails/second, bytes/second), tracking adherence to ISP hourly ramp-up volume ceilings and flagging inadvertent quota breaches.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [DISPATCH VELOCITY & ISP QUOTA TELEMETRY ANALYTICS] ───────────────────────────────────┐
      │ VELOCITY: 18.4 msg/s | HOURLY DISPATCHED: 482/1000 | QUOTA COMPLIANCE: 100% | MEAN RTT: 28ms│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ PROVIDER COHORT | CURR RATE | HOURLY CAP | CONSUMED | REMAINING | COMPLIANCE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Google Mail | 8.2 msg/s | 400 msg/h | 180 msgs | 220 msgs | [●] 100% OK
      │ Microsoft 365 | 5.1 msg/s | 300 msg/h | 142 msgs | 158 msgs | [●] 100% OK
      │ Yahoo / AOL | 2.8 msg/s | 150 msg/h | 84 msgs | 66 msgs | [●] 100% OK
      │ Fastmail / Other | 2.3 msg/s | 150 msg/h | 76 msgs | 74 msgs | [●] 100% OK
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: Google Mail] Velocity 8.2 msg/s within optimal range | Hourly quota headroom: 55%
      │ Keys: [↑/↓: Select Cohort] [r: Reset Hourly Counter] [c: Adjust Quota Cap] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native CSS analytics charts with embedded SVG line histograms, pure local storage bindings, zero CDN scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics velocity-analytics`
  - Flags:
    - `--cohort <string>`: Provider cohort to inspect (`google`, `microsoft`, `yahoo`, `all`; default: `all`).
    - `--window <duration>`: Rolling metrics window for velocity calculation (default: `1h`).
    - `--is-quota-guarded`: Block deliveries immediately upon reaching 100% of hourly cohort cap (default: `true`).
    - `--has-pace-smoothing`: Distribute outbound messages evenly across the hour window (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorVelocityAnalyticsRecord struct {
        ProviderCohort string    `json:"provider_cohort"`
        CurrentVelocityMsgPerSec float64 `json:"current_velocity_msg_per_sec"`
        HourlyCapLimit int       `json:"hourly_cap_limit"`
        ConsumedCount int        `json:"consumed_count"`
        RemainingQuota int       `json:"remaining_quota"`
        CompliancePercent float64 `json:"compliance_percent"`
        IsCompliant bool         `json:"is_compliant"`
        IsCapApproaching bool    `json:"is_cap_approaching"`
        HasPacingActive bool     `json:"has_pacing_active"`
    }

    type GovernorVelocityAnalyticsModel struct {
        Cohorts []GovernorVelocityAnalyticsRecord `json:"cohorts"`
        Cursor int               `json:"cursor"`
        TotalVelocityMsgPerSec float64 `json:"total_velocity_msg_per_sec"`
        IsMonitoring bool        `json:"is_monitoring"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorVelocityAnalyticsRecord (
        provider_cohort TEXT PRIMARY KEY,
        current_velocity_msg_per_sec REAL NOT NULL DEFAULT 0.0,
        hourly_cap_limit INTEGER NOT NULL DEFAULT 500,
        consumed_count INTEGER NOT NULL DEFAULT 0,
        remaining_quota INTEGER NOT NULL DEFAULT 500,
        compliance_percent REAL NOT NULL DEFAULT 100.0,
        is_compliant INTEGER NOT NULL DEFAULT 1,
        is_cap_approaching INTEGER NOT NULL DEFAULT 0,
        has_pacing_active INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics velocity-analytics`. Real-time telemetry reader backed by `repodb/pipeline.db` and rendered in Lipgloss TUI.

---

### Improvement 1650: Governor Automated Health Re-Arm & Self-Healing Watchdog
- **ID:** `GOV-SUPERVISOR-REARM-1650`
- **Cluster:** 165 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Supervise the concurrency governor background threads, automatically restarting dead worker routines, resolving SQLite connection deadlocks, and clearing stale lockfiles without manual intervention.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [GOVERNOR SUPERVISOR WATCHDOG & AUTO RE-ARM SENTINEL] ─────────────────────────────────┐
      │ WATCHDOG: ARMED | HEARTBEAT: 100ms | THREADS MONITORED: 8 | FAULTS RESOLVED (24h): 0        │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ THREAD SUBSYSTEM | HEARTBEAT | UPTIME | RESTARTS | MEMORY MB | HEALTH
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Worker Pool Manager | [✓] 12ms | 18h 42m | 0 | 24.2 MB | [●] NORMAL
      │ Token Bucket Refill | [✓] 8ms | 18h 42m | 0 | 4.1 MB | [●] NORMAL
      │ Circuit Breaker Poller | [✓] 15ms | 18h 42m | 0 | 6.8 MB | [●] NORMAL
      │ DNS MX Cache Cleaner | [✓] 22ms | 18h 42m | 0 | 8.5 MB | [●] NORMAL
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SUPERVISOR STATUS: ALL NORMAL] Watchdog ping validated | Lockfile status: UNLOCKED | No leaks
      │ Keys: [↑/↓: Select Subsystem] [r: Force Re-Arm Watchdog] [c: Clear Deadlocks] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML watchdog status card using system CSS monospace, zero external libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics governor-rearm`
  - Flags:
    - `--heartbeat-timeout <duration>`: Threshold duration before declaring a background worker dead (default: `5s`).
    - `--max-restarts <int>`: Maximum automatic restart attempts before alerting human operator (default: `3`).
    - `--is-auto-rearm`: Automatically restart stalled governor threads (default: `true`).
    - `--has-panic-recovery`: Wrap all worker goroutines in structured deferred panic handlers (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorSupervisorWatchdogRecord struct {
        SubsystemName string     `json:"subsystem_name"`
        LastHeartbeatMs int64    `json:"last_heartbeat_ms"`
        UptimeSeconds int64      `json:"uptime_seconds"`
        RestartCount int         `json:"restart_count"`
        MemoryBytes int64        `json:"memory_bytes"`
        IsAlive bool             `json:"is_alive"`
        IsHung bool              `json:"is_hung"`
        HasPanicRecorded bool    `json:"has_panic_recorded"`
        HasLockCleaned bool      `json:"has_lock_cleaned"`
    }

    type GovernorSupervisorWatchdogModel struct {
        Subsystems []GovernorSupervisorWatchdogRecord `json:"subsystems"`
        ActiveIndex int          `json:"active_index"`
        IsArmed bool             `json:"is_armed"`
        HasFaults bool           `json:"has_faults"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorSupervisorWatchdogRecord (
        subsystem_name TEXT PRIMARY KEY,
        last_heartbeat_ms INTEGER NOT NULL,
        uptime_seconds INTEGER NOT NULL DEFAULT 0,
        restart_count INTEGER NOT NULL DEFAULT 0,
        memory_bytes INTEGER NOT NULL DEFAULT 0,
        is_alive INTEGER NOT NULL DEFAULT 1,
        is_hung INTEGER NOT NULL DEFAULT 0,
        has_panic_recorded INTEGER NOT NULL DEFAULT 0,
        has_lock_cleaned INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics governor-rearm`. Runs as an independent supervision loop and logs recovery actions in `repodb/pipeline.db`.

---

# Cluster 166: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 1651–1660)

---

### Improvement 1651: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD
- **ID:** `HTML-SANDBOX-SANITIZER-1651`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Provide an isolated, zero-CDN HTML/CSS parsing and sanitization engine that validates outbound warmup templates against strict email client CSS subsets (Gmail, Outlook, Apple Mail), stripping scripts, insecure fonts, and external assets.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [ZERO-CDN EMAIL HTML SANITIZER & SECURITY SANDBOX] ────────────────────────────────────┐
      │ TEMPLATE: warmup-onboard-v1.html | SIZE: 34.2 KB | CSS COMPLIANCE: 99.8% | SCRIPTS: 0       │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ SECURITY/COMPATIBILITY CHECK | TARGET SCOPE | RESULT | COMPLIANCE | REMEDY | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Zero Remote Script Tags | Global Email | 0 Found | 100% | N/A | [✓] PASS
      │ Zero External CDN CSS | Global Email | 0 Found | 100% | N/A | [✓] PASS
      │ Strict Image Host Validation | CID / HTTPS | All Match | 100% | N/A | [✓] PASS
      │ Outlook MSO Conditional VML | Desktop 365 | Aligned | 98% | MSO Wraps | [✓] PASS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [VERIFIED CLEAN] Template conforms to strict Zero-CDN email client boundaries | Disallowed tags: NONE
      │ Keys: [v: View Diff] [s: Save Sanitized] [p: Preview Live] [t: Test Client Matrix] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Sandboxed `<iframe>` with `sandbox="allow-same-origin"` and zero external stylesheet or JavaScript references, rendering pure HTML templates locally.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics html-sandbox`
  - Flags:
    - `--template-path <string>`: Relative path to HTML email template file.
    - `--target-client <string>`: Target client profile (`gmail`, `outlook`, `apple`, `all`; default: `all`).
    - `--is-auto-sanitize`: Strip invalid tags automatically and output clean HTML (default: `true`).
    - `--has-strict-csp`: Enforce zero external domain connections in CSS (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlSanitizerRecord struct {
        SanitizerID string       `json:"sanitizer_id"`
        TemplateName string      `json:"template_name"`
        OriginalSizeBytes int    `json:"original_size_bytes"`
        SanitizedSizeBytes int   `json:"sanitized_size_bytes"`
        StrippedTagsCount int    `json:"stripped_tags_count"`
        CssComplianceScore float64 `json:"css_compliance_score"`
        IsZeroCdnCompliant bool  `json:"is_zero_cdn_compliant"`
        IsClean bool             `json:"is_clean"`
        HasExternalFonts bool    `json:"has_external_fonts"`
        HasJavascriptStripped bool `json:"has_javascript_stripped"`
    }

    type HtmlSanitizerModel struct {
        Sanitizer HtmlSanitizerRecord `json:"sanitizer"`
        ActiveTab int            `json:"active_tab"`
        IsProcessing bool        `json:"is_processing"`
        HasViolations bool       `json:"has_violations"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlSanitizerRecord (
        sanitizer_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        original_size_bytes INTEGER NOT NULL,
        sanitized_size_bytes INTEGER NOT NULL,
        stripped_tags_count INTEGER NOT NULL DEFAULT 0,
        css_compliance_score REAL NOT NULL DEFAULT 100.0,
        is_zero_cdn_compliant INTEGER NOT NULL DEFAULT 1,
        is_clean INTEGER NOT NULL DEFAULT 1,
        has_external_fonts INTEGER NOT NULL DEFAULT 0,
        has_javascript_stripped INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics html-sandbox`. Validates templates locally using native Go HTML tokenizer and persists audit logs in `repodb/pipeline.db`.

---

### Improvement 1652: Deterministic CSS-to-Style Attribute Inliner & Minification Engine
- **ID:** `HTML-CSS-INLINER-1652`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Transform embedded `<style>` blocks into deterministic inline `style="..."` attributes across HTML tags, ensuring bulletproof rendering in clients that strip header `<style>` elements (Gmail Web/Mobile).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [DETERMINISTIC CSS INLINER & MINIFIER ENGINE] ─────────────────────────────────────────┐
      │ RULES INLINED: 84 | SIZE BEFORE: 42.1 KB | AFTER: 38.6 KB | REDUCTION: 8.3%                 │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ SELECTOR | PROPERTIES INLINED | ELEMENTS MATCHED | SPECIFICITY | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ table.container | width, margin, pad | 4 elements | (0, 1, 1) | [✓] INLINED
      │ td.button-cell | background, border | 2 elements | (0, 1, 1) | [✓] INLINED
      │ p.body-text | font-family, color | 18 elements | (0, 1, 1) | [✓] INLINED
      │ a.footer-link | color, text-decor | 6 elements | (0, 1, 1) | [✓] INLINED
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [TRANSFORMATION COMPLETE] All CSS rules successfully inlined | Style tag purged safely
      │ Keys: [↑/↓: Browse Selectors] [i: Inline Template] [m: Toggle Minifier] [v: Diff View] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native HTML side-by-side diff view with zero-external script syntax highlighting, pure system CSS typography.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics css-inliner`
  - Flags:
    - `--input-file <string>`: Relative path to source HTML file.
    - `--output-file <string>`: Relative path to target destination for inlined HTML.
    - `--is-minification-active`: Minify HTML whitespace and attributes after inlining (default: `true`).
    - `--has-mso-preservation`: Preserve `<!--[if mso]>` conditional blocks untouched (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlCssInlinerRecord struct {
        InlinerID string         `json:"inliner_id"`
        InputPath string         `json:"input_path"`
        OutputPath string        `json:"output_path"`
        RulesInlinedCount int    `json:"rules_inlined_count"`
        ElementsMatchedCount int `json:"elements_matched_count"`
        SizeReductionBytes int   `json:"size_reduction_bytes"`
        IsInlined bool           `json:"is_inlined"`
        IsMinified bool          `json:"is_minified"`
        HasMsoPreserved bool     `json:"has_mso_preserved"`
        HasSpecificityConflict bool `json:"has_specificity_conflict"`
    }

    type HtmlCssInlinerModel struct {
        Records []HtmlCssInlinerRecord `json:"records"`
        CursorIndex int          `json:"cursor_index"`
        IsProcessing bool        `json:"is_processing"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlCssInlinerRecord (
        inliner_id TEXT PRIMARY KEY,
        input_path TEXT NOT NULL,
        output_path TEXT NOT NULL,
        rules_inlined_count INTEGER NOT NULL DEFAULT 0,
        elements_matched_count INTEGER NOT NULL DEFAULT 0,
        size_reduction_bytes INTEGER NOT NULL DEFAULT 0,
        is_inlined INTEGER NOT NULL DEFAULT 1,
        is_minified INTEGER NOT NULL DEFAULT 1,
        has_mso_preserved INTEGER NOT NULL DEFAULT 1,
        has_specificity_conflict INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics css-inliner`. Executes in-memory DOM manipulation via Go net/html and logs execution metrics to `repodb/pipeline.db`.

---

### Improvement 1653: Email Dark Mode Adaptation & WCAG Contrast Ratio Checker
- **ID:** `HTML-DARK-MODE-CONTRAST-1653`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify email template legibility under client dark mode color inversion engines (Apple Mail, Outlook, Gmail), auditing `@media (prefers-color-scheme: dark)` overrides and verifying WCAG AAA (7:1) / AA (4.5:1) contrast ratios.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [EMAIL DARK MODE ADAPTATION & CONTRAST RATIO AUDITOR] ─────────────────────────────────┐
      │ ELEMENTS TESTED: 42 | WCAG AA PASS: 100% | WCAG AAA PASS: 92% | DARK META: DETECTED         │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ TEMPLATE SECTION | LIGHT COLOR | DARK COLOR | CONTRAST (L) | CONTRAST (D) | COMPLIANCE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Header H1 Title | #1e293b on # | #f8fafc on # | 12.4:1 (AAA) | 14.2:1 (AAA) | [●] PASS
      │ Body Paragraph Text | #334155 on # | #e2e8f0 on # | 9.8:1 (AAA) | 11.1:1 (AAA) | [●] PASS
      │ Primary CTA Button | #ffffff on # | #ffffff on # | 5.2:1 (AA) | 5.2:1 (AA) | [●] PASS
      │ Footer Legal Note | #64748b on # | #94a3b8 on # | 4.8:1 (AA) | 5.1:1 (AA) | [●] PASS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: Primary CTA Button] Light #ffffff on #2563eb (5.2:1) | Dark #ffffff on #3b82f6 (5.2:1) | WCAG AA Met
      │ Keys: [↑/↓: Browse Elements] [t: Toggle Light/Dark Preview] [s: Adjust Theme Hex] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML dual-pane light/dark preview using native CSS custom properties and zero external scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics dark-mode`
  - Flags:
    - `--template-path <string>`: Relative path to target HTML email template.
    - `--wcag-level <string>`: Required accessibility conformance level (`AA`, `AAA`; default: `AA`).
    - `--is-inverted-simulated`: Simulate aggressive Outlook Android full color inversion (default: `true`).
    - `--has-dark-meta-tags`: Verify presence of `color-scheme` and `supported-color-schemes` meta tags (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlDarkModeRecord struct {
        ElementSelector string   `json:"element_selector"`
        LightTextColor string    `json:"light_text_color"`
        LightBgColor string      `json:"light_bg_color"`
        DarkTextColor string     `json:"dark_text_color"`
        DarkBgColor string       `json:"dark_bg_color"`
        LightContrastRatio float64 `json:"light_contrast_ratio"`
        DarkContrastRatio float64 `json:"dark_contrast_ratio"`
        IsWcagAaCompliant bool   `json:"is_wcag_aa_compliant"`
        IsWcagAaaCompliant bool  `json:"is_wcag_aaa_compliant"`
        HasDarkModeMeta bool     `json:"has_dark_mode_meta"`
    }

    type HtmlDarkModeModel struct {
        Elements []HtmlDarkModeRecord `json:"elements"`
        Cursor int               `json:"cursor"`
        IsAuditing bool          `json:"is_auditing"`
        HasContrastFailures bool `json:"has_contrast_failures"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlDarkModeRecord (
        element_selector TEXT PRIMARY KEY,
        light_text_color TEXT NOT NULL,
        light_bg_color TEXT NOT NULL,
        dark_text_color TEXT NOT NULL,
        dark_bg_color TEXT NOT NULL,
        light_contrast_ratio REAL NOT NULL DEFAULT 4.5,
        dark_contrast_ratio REAL NOT NULL DEFAULT 4.5,
        is_wcag_aa_compliant INTEGER NOT NULL DEFAULT 1,
        is_wcag_aaa_compliant INTEGER NOT NULL DEFAULT 0,
        has_dark_mode_meta INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics dark-mode`. Computes relative luminance mathematically using sRGB formula and stores audit in `repodb/pipeline.db`.

---

### Improvement 1654: MIME Multipart/Related CID Image Embedding & Zero-HTTP Asset Bundler
- **ID:** `HTML-CID-IMAGE-EMBEDDER-1654`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Package local brand logos and icons as embedded MIME `Content-ID` (CID) multipart parts, rewriting `src="cid:..."` attributes to guarantee instant image display without external HTTP requests or remote image tracking blocks.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [MIME MULTIPART/RELATED CID IMAGE EMBEDDER] ───────────────────────────────────────────┐
      │ IMAGES EMBEDDED: 3 | TOTAL PAYLOAD: 28.4 KB | EXTERNAL HTTP: 0 | ENCODING: BASE64           │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ IMAGE FILENAME | CONTENT-ID (CID) | FORMAT | RAW BYTES | BASE64 | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ logo-brand.png | cid:logo-brand-01 | image/png | 12.4 KB | 16.5 KB | [✓] EMBEDDED
      │ icon-check.png | cid:icon-check-02 | image/png | 4.2 KB | 5.6 KB | [✓] EMBEDDED
      │ avatar-team.jpg | cid:avatar-team-03 | image/jp | 6.1 KB | 8.1 KB | [✓] EMBEDDED
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [EMBEDDING COMPLETE] All image sources rewritten to local CID tokens | Remote network leaks: 0
      │ Keys: [↑/↓: Select Image] [e: Extract Asset] [r: Regenerate CIDs] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML email asset inspector with pure local Base64 previews, zero CDN scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics cid-embedder`
  - Flags:
    - `--html-file <string>`: Relative path to target HTML template.
    - `--assets-dir <string>`: Relative path to directory containing image assets.
    - `--max-image-bytes <int>`: Upper byte threshold per embedded image (default: `102400`).
    - `--is-auto-rewrite`: Automatically rewrite `<img src="...">` to `cid:` format (default: `true`).
    - `--has-lossless-compression`: Run lossless PNG/JPEG optimization before embedding (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlCidImageRecord struct {
        ImageFilename string     `json:"image_filename"`
        ContentID string         `json:"content_id"`
        MimeType string          `json:"mime_type"`
        SizeBytes int64          `json:"size_bytes"`
        Base64SizeBytes int64    `json:"base64_size_bytes"`
        IsEmbedded bool          `json:"is_embedded"`
        IsSizeAllowed bool       `json:"is_size_allowed"`
        HasLosslessOptimized bool `json:"has_lossless_optimized"`
    }

    type HtmlCidImageModel struct {
        Images []HtmlCidImageRecord `json:"images"`
        ActiveCursor int         `json:"active_cursor"`
        TotalPayloadBytes int64  `json:"total_payload_bytes"`
        IsBundling bool          `json:"is_bundling"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlCidImageRecord (
        image_filename TEXT PRIMARY KEY,
        content_id TEXT NOT NULL,
        mime_type TEXT NOT NULL,
        size_bytes INTEGER NOT NULL,
        base64_size_bytes INTEGER NOT NULL,
        is_embedded INTEGER NOT NULL DEFAULT 1,
        is_size_allowed INTEGER NOT NULL DEFAULT 1,
        has_lossless_optimized INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics cid-embedder`. Integrates with the Go MIME multipart message writer and records assets in `repodb/pipeline.db`.

---

### Improvement 1655: Microsoft Outlook Vector Markup Language (VML) Generator & Fallback Injector
- **ID:** `HTML-FALLBACK-VML-OUTLOOK-1655`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Generate and inject compliant Microsoft Word rendering engine conditional comments (`<!--[if mso]> <v:roundrect ...> <![endif]-->`) for rounded buttons and background images, fixing Outlook Windows rendering glitches without external hacks.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [OUTLOOK VML COMPATIBILITY GENERATOR & FALLBACK INJECTOR] ─────────────────────────────┐
      │ MSO BLOCKS GENERATED: 4 | BUTTONS WRAPPED: 2 | BG IMAGES: 1 | PARSER: RFC 822 VALID         │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ VML WRAPPER ELEMENT | TYPE | MSO TAG | ARCSIZE | FALLBACK | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Primary CTA Button | Round Button | v:roundrect | 15% | table/td | [✓] WRAPPED
      │ Secondary Action | Round Button | v:roundrect | 15% | table/td | [✓] WRAPPED
      │ Header Hero Section | Hero Bg Image | v:rect | 0% | td bg | [✓] WRAPPED
      │ Spacer Padding Block | DPI Guard | o:pixelspert | N/A | table 0 | [✓] WRAPPED
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [VML INJECTION VERIFIED] Clean Word HTML wrapper created | Non-MSO clients receive standard CSS
      │ Keys: [↑/↓: Browse VML Elements] [g: Generate MSO Wrappers] [v: Preview Word Render] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML preview displaying both WebKit rendering and Outlook simulated Word layout side-by-side using pure system CSS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics vml-generator`
  - Flags:
    - `--html-template <string>`: Relative path to source HTML file.
    - `--button-selector <string>`: CSS selector matching buttons needing VML wrappers (default: `.btn-primary`).
    - `--is-auto-inject`: Inject MSO wrappers automatically around matched elements (default: `true`).
    - `--has-dpi-scaling`: Include 96/120 DPI XML scaling guards for high-resolution displays (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlVmlGeneratorRecord struct {
        ElementName string       `json:"element_name"`
        ElementType string       `json:"element_type"`
        MsoTagType string        `json:"mso_tag_type"`
        ArcsizePercent int       `json:"arcsize_percent"`
        FallbackType string      `json:"fallback_type"`
        IsWrapped bool           `json:"is_wrapped"`
        IsDpiProtected bool      `json:"is_dpi_protected"`
        HasMsoCommentValid bool  `json:"has_mso_comment_valid"`
    }

    type HtmlVmlGeneratorModel struct {
        Elements []HtmlVmlGeneratorRecord `json:"elements"`
        CursorIndex int          `json:"cursor_index"`
        IsInjecting bool         `json:"is_injecting"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlVmlGeneratorRecord (
        element_name TEXT PRIMARY KEY,
        element_type TEXT NOT NULL,
        mso_tag_type TEXT NOT NULL,
        arcsize_percent INTEGER NOT NULL DEFAULT 10,
        fallback_type TEXT NOT NULL,
        is_wrapped INTEGER NOT NULL DEFAULT 1,
        is_dpi_protected INTEGER NOT NULL DEFAULT 1,
        has_mso_comment_valid INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics vml-generator`. Transforms AST nodes in memory and stores generated VML metadata in `repodb/pipeline.db`.

---

### Improvement 1656: Zero-External Web Font System Font Stack Enforcer & Fallback Auditor
- **ID:** `HTML-TYPOGRAPHY-SYSTEM-STACK-1656`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan HTML templates to eliminate external web font downloads (`@import url(fonts.googleapis.com...)`), replacing them with cross-platform native system font stacks (`-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif`) to ensure zero-latency rendering and privacy compliance.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SYSTEM FONT STACK ENFORCER & ZERO-EXTERNAL FONT AUDITOR] ─────────────────────────────┐
      │ FONTS AUDITED: 8 | REMOTE DOWNLOADS: 0 | NATIVE STACKS: 100% | PRIVACY COMPLIANT: YES       │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ TARGET SELECTOR | DECLARED FONT FAMILY | REMOTE URL | SYSTEM STACK | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ body, td, p | -apple-system, BlinkMacSystemFont... | NONE | [✓] ENFORCED | [●] SAFE
      │ h1, h2, h3 | ui-sans-serif, system-ui, sans-serif | NONE | [✓] ENFORCED | [●] SAFE
      │ code, pre | ui-monospace, SFMono-Regular, Consola | NONE | [✓] ENFORCED | [●] SAFE
      │ .serif-quote | Georgia, Cambria, 'Times New Roman' | NONE | [✓] ENFORCED | [●] SAFE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [PRIVACY VALIDATED] All external font linkages stripped | 100% pure system typography enforced
      │ Keys: [↑/↓: Browse Selectors] [s: Apply Recommended System Stack] [r: Strip Font URLs] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML typography specimen sheet rendered using pure local system fonts and native CSS variables.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics system-fonts`
  - Flags:
    - `--html-template <string>`: Relative path to target HTML email template.
    - `--font-preset <string>`: System typography preset (`sans`, `serif`, `monospace`; default: `sans`).
    - `--is-remote-stripped`: Strip all remote `@font-face` and Google/Adobe Font links (default: `true`).
    - `--has-windows-fallbacks`: Include legacy Windows Arial/Tahoma font fallbacks in stack (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlTypographyStackRecord struct {
        Selector string          `json:"selector"`
        DeclaredFontFamily string `json:"declared_font_family"`
        RemoteFontUrl string     `json:"remote_font_url"`
        SystemStackName string   `json:"system_stack_name"`
        IsRemoteFree bool        `json:"is_remote_free"`
        IsSystemEnforced bool    `json:"is_system_enforced"`
        HasPlatformParity bool   `json:"has_platform_parity"`
    }

    type HtmlTypographyStackModel struct {
        Records []HtmlTypographyStackRecord `json:"records"`
        Cursor int               `json:"cursor"`
        IsAuditing bool          `json:"is_auditing"`
        HasRemoteFontViolations bool `json:"has_remote_font_violations"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlTypographyStackRecord (
        selector TEXT PRIMARY KEY,
        declared_font_family TEXT NOT NULL,
        remote_font_url TEXT NOT NULL DEFAULT 'NONE',
        system_stack_name TEXT NOT NULL,
        is_remote_free INTEGER NOT NULL DEFAULT 1,
        is_system_enforced INTEGER NOT NULL DEFAULT 1,
        has_platform_parity INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics system-fonts`. Audits CSS property values and logs typography compliance in `repodb/pipeline.db`.

---

### Improvement 1657: Gmail 102 KB Clipping Guard & HTML Payload Budget Monitor
- **ID:** `HTML-PAYLOAD-GZIP-BUDGET-1657`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Enforce strict byte size caps (< 102,400 bytes uncompressed) on outbound HTML email bodies, preventing Gmail from displaying the destructive '[Message clipped] View entire message' truncated footer.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [GMAIL 102 KB CLIPPING GUARD & PAYLOAD BUDGET MONITOR] ────────────────────────────────┐
      │ BODY SIZE: 64.2 KB (62.7% OF LIMIT) | CEILING: 102.4 KB | HEADROOM: 38.2 KB | CLIPPING: SAFE│
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ PAYLOAD SECTION | RAW BYTES | PERCENT | BUDGET ALLOC | MINIFIED | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Inlined CSS Styles | 18.4 KB | 28.6% | 25.0 KB max | [✓] YES | [●] OPTIMAL
      │ DOM Elements & Text | 32.1 KB | 50.0% | 50.0 KB max | [✓] YES | [●] OPTIMAL
      │ Base64 Small Assets | 12.8 KB | 19.9% | 20.0 KB max | [✓] YES | [●] OPTIMAL
      │ HTML Comment Garbage | 0.9 KB | 1.4% | 2.0 KB max | [✓] STRIP | [●] OPTIMAL
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [GMAIL SAFE] Total uncompressed body size: 64.2 KB | Headroom: 38.2 KB before 102 KB clip boundary
      │ Keys: [↑/↓: Browse Sections] [c: Strip HTML Comments] [o: Compress Assets] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML gauge showing payload distribution bar with warning threshold marker at 102 KB using native CSS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics payload-budget`
  - Flags:
    - `--html-file <string>`: Relative path to HTML email template.
    - `--budget-kb <float>`: Maximum payload ceiling in kilobytes (default: `100.0`).
    - `--is-comment-stripped`: Strip all non-conditional HTML comments to reduce size (default: `true`).
    - `--has-clipping-warning`: Return exit code 1 if payload exceeds 95 KB (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlPayloadBudgetRecord struct {
        TemplateName string      `json:"template_name"`
        TotalBytes int64         `json:"total_bytes"`
        BudgetCapBytes int64     `json:"budget_cap_bytes"`
        HeadroomBytes int64      `json:"headroom_bytes"`
        UsagePercent float64     `json:"usage_percent"`
        IsWithinBudget bool      `json:"is_within_budget"`
        IsClippingRisk bool      `json:"is_clipping_risk"`
        HasCommentsStripped bool `json:"has_comments_stripped"`
    }

    type HtmlPayloadBudgetModel struct {
        Budget HtmlPayloadBudgetRecord `json:"budget"`
        IsEvaluating bool        `json:"is_evaluating"`
        HasAlerts bool           `json:"has_alerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlPayloadBudgetRecord (
        template_name TEXT PRIMARY KEY,
        total_bytes INTEGER NOT NULL,
        budget_cap_bytes INTEGER NOT NULL DEFAULT 102400,
        headroom_bytes INTEGER NOT NULL,
        usage_percent REAL NOT NULL DEFAULT 0.0,
        is_within_budget INTEGER NOT NULL DEFAULT 1,
        is_clipping_risk INTEGER NOT NULL DEFAULT 0,
        has_comments_stripped INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics payload-budget`. Measures exact byte count of serialized RFC 822 body and records metrics in `repodb/pipeline.db`.

---

### Improvement 1658: RFC 5322 MIME Header Structure & Cryptographic Alignment Linter
- **ID:** `HTML-RFC5322-HEADER-LINT-1658`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit RFC 5322 message headers (Message-ID, Date, From, To, List-Unsubscribe, Feedback-ID, Precedence: bulk) to guarantee zero header syntax violations, deterministic unique ID generation, and SPF/DKIM/DMARC domain alignment.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [RFC 5322 MIME HEADER STRUCTURE & ALIGNMENT LINTER] ───────────────────────────────────┐
      │ HEADERS AUDITED: 14 | RFC 5322 COMPLIANT: 100% | CRYPTO ALIGNMENT: STRICT | ERRORS: 0       │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ HEADER FIELD | VALUE PREVIEW | RFC RULE | ALIGNMENT | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Message-ID | <20261028.9102@warmup.domain.com> | RFC 5322 | Domain OK | [✓] PASS
      │ From | "Sender Name" <alice@domain.com> | RFC 5322 | DMARC OK | [✓] PASS
      │ List-Unsubscribe | <mailto:unsub@domain.com>, <https... | RFC 8058 | One-Click | [✓] PASS
      │ Feedback-ID | warmup-pool:client-01:campaign-2 | ESP Std | Aligned | [✓] PASS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED: List-Unsubscribe] RFC 8058 One-Click Post header verified | Target: mailto:unsub@domain.com
      │ Keys: [↑/↓: Browse Headers] [v: Inspect Raw Header Block] [c: Check DMARC Alignment] [q: Return]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML raw header inspector with pure system CSS syntax coloring, zero external libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics header-linter`
  - Flags:
    - `--raw-eml-file <string>`: Relative path to raw .eml file or header dump.
    - `--is-strict-rfc8058`: Verify List-Unsubscribe and List-Unsubscribe-Post one-click compliance (default: `true`).
    - `--has-dmarc-alignment`: Verify From domain matches DKIM d= and SPF Return-Path domains (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlHeaderLinterRecord struct {
        HeaderName string        `json:"header_name"`
        HeaderValue string       `json:"header_value"`
        RfcStandard string       `json:"rfc_standard"`
        IsAligned bool           `json:"is_aligned"`
        IsValid bool             `json:"is_valid"`
        HasOneClickUnsub bool    `json:"has_one_click_unsub"`
        HasDkimMatch bool        `json:"has_dkim_match"`
    }

    type HtmlHeaderLinterModel struct {
        Headers []HtmlHeaderLinterRecord `json:"headers"`
        ActiveCursor int         `json:"active_cursor"`
        IsAuditing bool          `json:"is_auditing"`
        HasSyntaxErrors bool     `json:"has_syntax_errors"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlHeaderLinterRecord (
        header_name TEXT PRIMARY KEY,
        header_value TEXT NOT NULL,
        rfc_standard TEXT NOT NULL,
        is_aligned INTEGER NOT NULL DEFAULT 1,
        is_valid INTEGER NOT NULL DEFAULT 1,
        has_one_click_unsub INTEGER NOT NULL DEFAULT 1,
        has_dkim_match INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics header-linter`. Parses message headers with Go `net/mail` and registers linter outcomes in `repodb/pipeline.db`.

---

### Improvement 1659: Spam Trigger Keyword & Algorithmic Lexicon Scanner
- **ID:** `HTML-SPAM-TRIGGER-LEXICON-1659`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan email subject lines and body copy against high-risk spam phrase dictionaries (e.g., deceptive urgency, financial claims, excessive capitalization), calculating composite SpamAssassin-style Bayesian risk scores before sending.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SPAM TRIGGER KEYWORD & LEXICON BAYESIAN SCANNER] ─────────────────────────────────────┐
      │ WORDS SCANNED: 342 | HIGH RISK PHRASES: 0 | BAYESIAN SCORE: 0.12 (LOW) | RATING: CLEAN      │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ TRIGGER CATEGORY | TERMS FOUND | WEIGHT | IMPACT SCORE | SUGGESTION | STATUS
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Urgency & Scarcity | 0 terms | 0.20 | 0.00 | None | [●] CLEAN
      │ Financial & Guarante | 0 terms | 0.35 | 0.00 | None | [●] CLEAN
      │ Excessive Caps/Excl | 1 term | 0.10 | 0.10 | Lowercase | [✓] ACCEPT
      │ Deceptive Headers | 0 terms | 0.50 | 0.00 | None | [●] CLEAN
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [CONTENT RATING: EXCELLENT] Composite Bayesian score: 0.10 (Threshold: < 2.50) | In-box probability: 99.4%
      │ Keys: [↑/↓: Browse Categories] [e: View Highlighted Words] [s: Rescan Copy] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone responsive HTML copy editor with native CSS inline amber highlights on detected phrases, zero external JS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics spam-lexicon`
  - Flags:
    - `--subject <string>`: Subject line text to evaluate for spam keywords.
    - `--body-file <string>`: Relative path to email body text or HTML file.
    - `--score-threshold <float>`: Maximum allowable risk score before blocking delivery (default: `2.5`).
    - `--is-bayesian-active`: Apply statistical Bayesian term probability weighting (default: `true`).
    - `--has-synonym-suggestions`: Provide safe conversational synonym alternatives for flagged terms (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlSpamLexiconRecord struct {
        CategoryName string      `json:"category_name"`
        TermsFoundCount int      `json:"terms_found_count"`
        CategoryWeight float64   `json:"category_weight"`
        ImpactScore float64      `json:"impact_score"`
        SuggestedReplacement string `json:"suggested_replacement"`
        IsClean bool             `json:"is_clean"`
        IsHighRisk bool          `json:"is_high_risk"`
        HasReplacements bool     `json:"has_replacements"`
    }

    type HtmlSpamLexiconModel struct {
        Categories []HtmlSpamLexiconRecord `json:"categories"`
        Cursor int               `json:"cursor"`
        CompositeScore float64   `json:"composite_score"`
        IsScanning bool          `json:"is_scanning"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlSpamLexiconRecord (
        category_name TEXT PRIMARY KEY,
        terms_found_count INTEGER NOT NULL DEFAULT 0,
        category_weight REAL NOT NULL DEFAULT 0.2,
        impact_score REAL NOT NULL DEFAULT 0.0,
        suggested_replacement TEXT NOT NULL DEFAULT 'None',
        is_clean INTEGER NOT NULL DEFAULT 1,
        is_high_risk INTEGER NOT NULL DEFAULT 0,
        has_replacements INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics spam-lexicon`. Evaluates tokens using an embedded Trie dictionary and logs scores in `repodb/pipeline.db`.

---

### Improvement 1660: Interactive Bubble Tea & Zero-CDN Responsive Viewport Simulator HUD
- **ID:** `HTML-RESPONSIVE-PREVIEW-1660`
- **Cluster:** 166 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Render dual-mode responsive previews of email templates in both terminal ASCII/half-block text rendering and a zero-CDN local HTTP server supporting mobile (375px), tablet (768px), and desktop (1200px) viewport switching.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [RESPONSIVE EMAIL VIEWPORT SIMULATOR HUD] ─────────────────────────────────────────────┐
      │ TEMPLATE: newsletter-dev-01.html | VIEWPORT: MOBILE (375px) | SCROLL: 0% | STATUS: READY    │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ DEVICE VIEWPORT | WIDTH px | HEIGHT px | ORIENTATION | BREAKPOINT | STATE
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Mobile iPhone SE | 375px | 667px | Portrait | < 480px | [●] SELECTED
      │ Mobile iPhone Pro | 393px | 852px | Portrait | < 480px | [ ] READY
      │ Tablet iPad Mini | 768px | 1024px | Portrait | < 800px | [ ] READY
      │ Desktop Standard | 1200px | 900px | Landscape | > 800px | [ ] READY
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SIMULATING: Mobile iPhone SE] Local Zero-CDN Preview URL: http://127.0.0.1:48291/preview?view=mobile
      │ Keys: [↑/↓: Change Device] [m: Switch to Mobile] [d: Switch to Desktop] [o: Open Browser] [q: Quit]
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Zero-CDN responsive HTML wrapper with embedded iframe viewports, pure CSS toggle buttons (`375px`, `768px`, `1200px`), and zero remote JavaScript.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics responsive-preview`
  - Flags:
    - `--html-template <string>`: Relative path to target HTML email file.
    - `--port <int>`: Local ephemeral HTTP port for zero-CDN browser preview (default: `48291`).
    - `--initial-viewport <string>`: Starting device simulation (`mobile`, `tablet`, `desktop`; default: `mobile`).
    - `--is-terminal-ascii`: Render text-mode ASCII email layout directly in terminal window (default: `true`).
    - `--has-auto-reload`: Watch template file and reload browser viewport on edit (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlResponsivePreviewRecord struct {
        DeviceName string        `json:"device_name"`
        WidthPx int              `json:"width_px"`
        HeightPx int             `json:"height_px"`
        Orientation string       `json:"orientation"`
        BreakpointRule string    `json:"breakpoint_rule"`
        PreviewUrl string        `json:"preview_url"`
        IsSelected bool          `json:"is_selected"`
        IsMobile bool            `json:"is_mobile"`
        HasLiveReload bool       `json:"has_live_reload"`
    }

    type HtmlResponsivePreviewModel struct {
        Viewports []HtmlResponsivePreviewRecord `json:"viewports"`
        ActiveIndex int          `json:"active_index"`
        LocalServerPort int      `json:"local_server_port"`
        IsServerRunning bool     `json:"is_server_running"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlResponsivePreviewRecord (
        device_name TEXT PRIMARY KEY,
        width_px INTEGER NOT NULL,
        height_px INTEGER NOT NULL,
        orientation TEXT NOT NULL,
        breakpoint_rule TEXT NOT NULL,
        preview_url TEXT NOT NULL,
        is_selected INTEGER NOT NULL DEFAULT 0,
        is_mobile INTEGER NOT NULL DEFAULT 1,
        has_live_reload INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Command handler in `gitmap email warmup diagnostics responsive-preview`. Starts local Go `http.Server` bound to localhost with zero CDN dependencies and logs sessions in `repodb/pipeline.db`.

---
