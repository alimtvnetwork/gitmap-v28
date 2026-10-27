# Specification 352: Improvements 1531 to 1560 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 352
- **Slug:** `352-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/03-improvements-1531-to-1560-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (1531 - 1560)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 154** | Items 1531–1540 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 155** | Items 1541–1550 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 156** | Items 1551–1560 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 154: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 1531–1540)

---

### Improvement 1531: Interactive Bubble Tea Mailbox Health Matrix & Distributed Edge Seed Topology Visualizer
- **ID:** `SEED-MATRIX-TOPOLOGY-1531`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a real-time, interactive Bubble Tea terminal grid visualizing distributed seed mailbox nodes across global Edge points of presence (PoPs), tracking Anycast routing latencies, ISP inbox placement cohorts (Gmail, M365, Yahoo, Fastmail), and automated seed cluster drain states.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status badges (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` quarantined rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [BUBBLE TEA SEED POOL HEALTH MATRIX & TOPOLOGY] ────────────────────────────────────────┐
      │ POOL NODES: 864 | ACTIVE SEEDS: 842 | DEGRADED: 18 | QUARANTINED: 4 | MEAN HEALTH: 99.1%   │
      ├───────────────────────┬────────────┬─────────────┬──────────────┬─────────────┬────────────┤
      │ SEED MAILBOX ID       │ PROVIDER   │ REGION      │ INBOX %      │ RTT (ms)    │ STATUS     │
      ├───────────────────────┼────────────┼─────────────┼──────────────┼─────────────┼────────────┤
      │ seed-gsuite-us-204    │ Google     │ us-east-1   │ 99.8%        │ 16.4ms      │ [●] STABLE │
      │ seed-m365-eu-108      │ Microsoft  │ eu-central-1│ 98.7%        │ 21.2ms      │ [●] STABLE │
      │ seed-yahoo-us-044     │ Yahoo      │ us-west-2   │ 93.1%        │ 39.8ms      │ [▲] WARN   │
      │ seed-fastmail-ap-012  │ Fastmail   │ ap-east-1   │ 99.4%        │ 28.5ms      │ [●] STABLE │
      ├───────────────────────┴────────────┴─────────────┴──────────────┴─────────────┴────────────┤
      │ [SELECTED: seed-yahoo-us-044] Health: 93.1% | Edge MX: mx-core.mail.yahoo.com              │
      │ Keys: [↑/↓: Navigate] [enter: Drill Down] [f: Toggle Alert Filter] [r: Force Probe] [q: Esc]│
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive CSS flexbox/grid containers, embedded system fonts (`ui-monospace`, `monospace`), and semantic HTML5 cards with zero external CDN scripts or remote stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-matrix`
  - Flags:
    - `--provider <string>`: Filter seed pool view by ESP provider (`google`, `microsoft`, `yahoo`, `fastmail`, `all`; default: `all`).
    - `--refresh-rate <duration>`: Interval for real-time matrix refresh and probe streaming (default: `1s`).
    - `--is-interactive`: Enable interactive keyboard navigation cursor and detail drill-down (default: `true`).
    - `--has-alerts-only`: Filter matrix view to display only degraded or critical seed mailboxes (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedMailboxHealthMatrixRecord struct {
        MailboxID          string  `json:"mailbox_id"`
        ProviderName       string  `json:"provider_name"`
        ClusterRegion      string  `json:"cluster_region"`
        HealthScorePercent float64 `json:"health_score_percent"`
        InboxPlacementRate float64 `json:"inbox_placement_rate"`
        SpamFolderRate     float64 `json:"spam_folder_rate"`
        LatencyRttMs       float64 `json:"latency_rtt_ms"`
        IsHealthy          bool    `json:"is_healthy"`
        IsQuarantined      bool    `json:"is_quarantined"`
        HasAnomalousSpike  bool    `json:"has_anomalous_spike"`
        HasVerifiedAuth    bool    `json:"has_verified_auth"`
    }

    type SeedMailboxHealthMatrixModel struct {
        Records        []SeedMailboxHealthMatrixRecord `json:"records"`
        CursorIndex    int                             `json:"cursor_index"`
        SelectedRecord SeedMailboxHealthMatrixRecord   `json:"selected_record"`
        IsStreaming    bool                            `json:"is_streaming"`
        HasAlerts      bool                            `json:"has_alerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedMailboxHealthMatrixRecord (
        mailbox_id TEXT PRIMARY KEY,
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
  - Bound to `gitmap email warmup diagnostics seed-matrix`. Directly queries `repodb/pipeline.db` using the shared read replica and delivers real-time Bubble Tea frame updates via `tea.Program` loop without intermediate disk locks.

---

### Improvement 1532: Seed Mailbox IMAP Auto-Drain & Spam Rescue Inspector HUD
- **ID:** `SEED-IMAP-DRAIN-INSPECT-1532`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide high-speed IMAP queue flushing, spam-to-inbox rescue verification, and quarantine telemetry across active warmup seeds to prevent storage exhaustion and verify positive algorithmic feedback loops.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live IMAP drain visualizer displaying rescued message counters, false-positive flag clears, and batch deletion progress.
    - ANSI Progress Bar:
      ```text
      [IMAP AUTO-DRAIN SESSION: ssn-drain-9042]
      Queue Status: [====================>-----] 80% Drained (412/515 msgs) | Rescued: 28 | Purged: 384
      Current Rate: 42.6 msgs/sec | Worker Goroutines: 16 | Errors: 0
      ```
  - **Zero-CDN Web View:**
    - Inline HTML/CSS status table with pure embedded styles showing mailbox IMAP connection status, unread message backlog, and auto-rescue actions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-drain`
  - Flags:
    - `--mailbox-id <string>`: Specific seed mailbox ID to drain or `all` (default: `all`).
    - `--batch-size <int>`: Batch size of messages per IMAP fetch cycle (default: `50`).
    - `--is-dry-run`: Simulate drain actions without deleting or moving messages (default: `false`).
    - `--has-auto-rescue`: Automatically move spam folder messages to inbox (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedImapDrainRecord struct {
        DrainSessionID    string `json:"drain_session_id"`
        MailboxID         string `json:"mailbox_id"`
        ProcessedMessages int    `json:"processed_messages"`
        RescuedSpamCount  int    `json:"rescued_spam_count"`
        PurgedTrashCount  int    `json:"purged_trash_count"`
        DurationMs        int64  `json:"duration_ms"`
        IsCompleted       bool   `json:"is_completed"`
        HasErrors         bool   `json:"has_errors"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedImapDrainRecord (
        drain_session_id TEXT PRIMARY KEY,
        mailbox_id TEXT NOT NULL,
        processed_messages INTEGER NOT NULL DEFAULT 0,
        rescued_spam_count INTEGER NOT NULL DEFAULT 0,
        purged_trash_count INTEGER NOT NULL DEFAULT 0,
        duration_ms INTEGER NOT NULL DEFAULT 0,
        is_completed INTEGER NOT NULL DEFAULT 0,
        has_errors INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates with `gitmap email warmup diagnostics seed-drain` and triggers asynchronous worker pools via lemon daemon background tasks.

---

### Improvement 1533: Seed Reputation Trend Graph & Cohort Degradation Detector
- **ID:** `SEED-REPUTATION-TREND-1533`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Compute moving averages and trend derivatives for deliverability cohorts to isolate seed mailboxes showing early signs of ISP throttling, algorithmic demotion, or greylisting.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Braille-character trend graph (`⡀⠄⠂⠁⠁⠄⡀`) displaying deliverability curves over the last 14 days.
    - Alert card highlighting cohorts dropping below the 95% inbox placement SLA.
      ```text
      ┌─ [SEED COHORT REPUTATION DEGRADATION DETECTOR] ───────────────────────────────────────────┐
      │ COHORT: m365-enterprise-de | SLA: 91.2% (WARN) | SLO SLOPE: -1.8%/day (DEMOTION RISK)     │
      │ 14-DAY CURVE: 100% ⡿⣟⣯⣷⣾⣿⣷⣯⣟⡿⣟⣯ 91.2% | TRIGGER: Microsoft EOP Rate Throttling       │
      │ Action: Auto-cooling period active (24h) | Traffic divert to: m365-enterprise-nl         │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - SVG sparklines rendered purely using server-side XML string generators with zero client JS libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-trends`
  - Flags:
    - `--timeframe <string>`: Analysis window (`24h`, `7d`, `14d`, `30d`; default: `7d`).
    - `--threshold <float>`: SLA drop threshold percentage (default: `95.0`).
    - `--is-verbose`: Include granular per-hour probe measurements (default: `false`).
    - `--has-auto-demote`: Automatically remove degrading seeds from active warmup pools (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedReputationTrendRecord struct {
        CohortID          string  `json:"cohort_id"`
        Provider          string  `json:"provider"`
        CurrentSlaPercent float64 `json:"current_sla_percent"`
        SlopeDerivative   float64 `json:"slope_derivative"`
        ObservationWindow string  `json:"observation_window"`
        IsDegrading       bool    `json:"is_degrading"`
        HasBreachedSla    bool    `json:"has_breached_sla"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedReputationTrendRecord (
        cohort_id TEXT PRIMARY KEY,
        provider TEXT NOT NULL,
        current_sla_percent REAL NOT NULL DEFAULT 100.0,
        slope_derivative REAL NOT NULL DEFAULT 0.0,
        observation_window TEXT NOT NULL,
        is_degrading INTEGER NOT NULL DEFAULT 0,
        has_breached_sla INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Directly invoked by `gitmap email warmup diagnostics seed-trends` and feeds seed telemetry into the multi-domain concurrency governor.

---

### Improvement 1534: Seed Mailbox Anycast Routing & Latency Heatmap Inspector
- **ID:** `SEED-ANYCAST-LATENCY-1534`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Measure and display regional round-trip times (RTT) from sender nodes to seed mailboxes across worldwide data centers to detect Anycast routing anomalies.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - ANSI heatmap terminal matrix color-coded by ping times: `<20ms` Green, `20-60ms` Yellow, `>60ms` Red.
    - Regional routing summary: `US-East: 17ms | US-West: 32ms | EU-Central: 21ms | AP-South: 64ms`.
  - **Zero-CDN Web View:**
    - Pure CSS grid layout with semantic badge tokens and zero third-party font or asset loads.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-latency`
  - Flags:
    - `--region <string>`: Target edge region (`us-east`, `us-west`, `eu-west`, `ap-east`, `all`; default: `all`).
    - `--ping-count <int>`: Probes per region (default: `5`).
    - `--is-live`: Continuously stream latency samples (default: `false`).
    - `--has-jitter-check`: Calculate standard deviation across successive probe packets (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedAnycastLatencyRecord struct {
        ProbeID       string  `json:"probe_id"`
        SourceRegion  string  `json:"source_region"`
        TargetRegion  string  `json:"target_region"`
        RttMedianMs   float64 `json:"rtt_median_ms"`
        JitterMs      float64 `json:"jitter_ms"`
        PacketLossPct float64 `json:"packet_loss_pct"`
        IsOptimal     bool    `json:"is_optimal"`
        HasHighJitter bool    `json:"has_high_jitter"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedAnycastLatencyRecord (
        probe_id TEXT PRIMARY KEY,
        source_region TEXT NOT NULL,
        target_region TEXT NOT NULL,
        rtt_median_ms REAL NOT NULL DEFAULT 0.0,
        jitter_ms REAL NOT NULL DEFAULT 0.0,
        packet_loss_pct REAL NOT NULL DEFAULT 0.0,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_high_jitter INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked via `gitmap email warmup diagnostics seed-latency` with results written to `repodb/pipeline.db`.

---

### Improvement 1535: Seed Authentication Audit & SPF/DKIM/DMARC Pass-Rate HUD
- **ID:** `SEED-AUTH-AUDIT-HUD-1535`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Audit inbound authentication headers (`Authentication-Results`, `ARC-Authentication-Results`) received by seed mailboxes to verify 100% cryptographic alignment.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Cryptographic compliance badges: `[SPF: PASS (100%)]`, `[DKIM: PASS 2048-bit (100%)]`, `[DMARC: PASS reject (100%)]`, `[ARC: PASS (100%)]`.
    - Real-time diagnostic split table showing DNS selector lookup status and alignment modes (`strict` vs `relaxed`).
  - **Zero-CDN Web View:**
    - High-contrast HTML audit table with embedded SVG shield badges.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-auth`
  - Flags:
    - `--domain <string>`: Sending domain to verify against seed pool headers.
    - `--sample-size <int>`: Number of recent seed header samples to parse (default: `100`).
    - `--is-strict-dmarc`: Require strict `adkim=s` and `aspf=s` alignment (default: `true`).
    - `--has-arc-verification`: Validate Authenticated Received Chain sealed headers (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedAuthAuditRecord struct {
        AuditID       string  `json:"audit_id"`
        DomainName    string  `json:"domain_name"`
        SpfPassRate   float64 `json:"spf_pass_rate"`
        DkimPassRate  float64 `json:"dkim_pass_rate"`
        DmarcPassRate float64 `json:"dmarc_pass_rate"`
        ArcPassRate   float64 `json:"arc_pass_rate"`
        IsAligned     bool    `json:"is_aligned"`
        HasFailures   bool    `json:"has_failures"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedAuthAuditRecord (
        audit_id TEXT PRIMARY KEY,
        domain_name TEXT NOT NULL,
        spf_pass_rate REAL NOT NULL DEFAULT 1.0,
        dkim_pass_rate REAL NOT NULL DEFAULT 1.0,
        dmarc_pass_rate REAL NOT NULL DEFAULT 1.0,
        arc_pass_rate REAL NOT NULL DEFAULT 1.0,
        is_aligned INTEGER NOT NULL DEFAULT 1,
        has_failures INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics seed-auth`. Outputs JSON summary logs to `.ai-memory/telemetry/` for audit persistence.

---

### Improvement 1536: Seed Pool Cohort Entropy & Conversational Threading Pattern Simulator
- **ID:** `SEED-COHORT-ENTROPY-1536`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Model and simulate organic peer-to-peer conversation dynamics (multi-turn replies, quote nesting, forward chains, randomized read delays) across seed pools to prevent heuristic spam flagging.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Tree-graph visualization of conversation threads:
      ```text
      ┌─ [CONVERSATIONAL THREADING SIMULATOR: ENTROPY INDEX 0.94] ─────────────────────────────────┐
      │ THREAD: thrd-alpha-8812 (4 Turns)                                                          │
      │ ├── [T0] sender@brand.io -> seed-01@gmail.com (Initial Outreach)                           │
      │ ├── [T1] seed-01@gmail.com -> sender@brand.io (Inquiry Reply +38m)                         │
      │ ├── [T2] sender@brand.io -> seed-01@gmail.com (Clarification +14m)                         │
      │ └── [T3] seed-01@gmail.com -> sender@brand.io (Engagement Close +1h 12m)                   │
      │ Threading Naturalness Score: 98.4% | Entropy Level: HIGH [●]                               │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Interactive nested timeline view constructed with pure HTML `<details>` and CSS borders.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-entropy`
  - Flags:
    - `--min-turns <int>`: Minimum conversational exchange turns (default: `2`).
    - `--max-turns <int>`: Maximum conversational exchange turns (default: `6`).
    - `--is-natural-timing`: Apply Poisson distribution delay model (default: `true`).
    - `--has-forwarding-simulation`: Simulate internal company forwarding (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedEntropyRecord struct {
        ThreadID           string  `json:"thread_id"`
        TurnCount          int     `json:"turn_count"`
        EntropyScore       float64 `json:"entropy_score"`
        MeanDelayMinutes   float64 `json:"mean_delay_minutes"`
        IsNaturalPattern   bool    `json:"is_natural_pattern"`
        HasForwardedAction bool    `json:"has_forwarded_action"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedEntropyRecord (
        thread_id TEXT PRIMARY KEY,
        turn_count INTEGER NOT NULL DEFAULT 1,
        entropy_score REAL NOT NULL DEFAULT 1.0,
        mean_delay_minutes REAL NOT NULL DEFAULT 30.0,
        is_natural_pattern INTEGER NOT NULL DEFAULT 1,
        has_forwarded_action INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics seed-entropy` and coordinates with lemon daemon for scheduled reply triggering.

---

### Improvement 1537: Seed Spam-Trap Honeypot Quarantine & Disqualification Matrix
- **ID:** `SEED-HONEYPOT-QUARANTINE-1537`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Automatically identify and quarantine compromised seed addresses, reciprocal spam traps, or recycled honeypots detected in test lists.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Red alert isolation banner with single-keystroke quarantine confirmation:
      ```text
      [!] HONEYPOT ANOMALY DETECTED: seed-oldbox-33@yahoo.com
      Matched Known Recycled Trap Pattern: MX Bounce Signature 554-5.7.1
      [ENTER: Instant Purge & Quarantine] | [ESC: Dismiss False Positive]
      ```
  - **Zero-CDN Web View:**
    - Security quarantine table with isolation timestamps, trigger rules, and audit logs.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-quarantine`
  - Flags:
    - `--auto-isolate`: Instantly isolate flagged seeds without prompt (default: `true`).
    - `--reason-filter <string>`: Filter by trigger reason (`honeypot`, `hard-bounce`, `rbl-listed`; default: `all`).
    - `--is-export-csv`: Export quarantine list to local CSV (default: `false`).
    - `--has-notify`: Emit OS terminal alert beep on detection (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedQuarantineRecord struct {
        QuarantineID      string `json:"quarantine_id"`
        MailboxAddress    string `json:"mailbox_address"`
        DetectionReason   string `json:"detection_reason"`
        IsolatedTimestamp int64  `json:"isolated_timestamp"`
        IsActive          bool   `json:"is_active"`
        HasPermanentBan   bool   `json:"has_permanent_ban"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedQuarantineRecord (
        quarantine_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL UNIQUE,
        detection_reason TEXT NOT NULL,
        isolated_timestamp INTEGER NOT NULL,
        is_active INTEGER NOT NULL DEFAULT 1,
        has_permanent_ban INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked by `gitmap email warmup diagnostics seed-quarantine`. Syncs instantly with dispatch engine blocklists.

---

### Improvement 1538: Seed Inbox Folder Placement Live Ticker & Provider SLA Matrix
- **ID:** `SEED-INBOX-TICKER-SLA-1538`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Deliver a high-frequency live ticker streaming placement outcomes across Primary Inbox, Promotions, Updates, and Spam folders across major webmail providers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Continuous scrolling Lipgloss ticker pane:
      ```text
      [14:22:01] Gmail -> PRIMARY (100%) | [14:22:04] Outlook -> INBOX (100%) | [14:22:08] Yahoo -> SPAM (ALERT)
      OVERALL PROVIDER SLA: Google: 99.4% | Microsoft: 98.8% | Yahoo: 92.1% | Apple: 99.9%
      ```
  - **Zero-CDN Web View:**
    - Pure CSS marquee/ticker element using CSS `@keyframes` with zero external JS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-ticker`
  - Flags:
    - `--poll-frequency <duration>`: Ticker polling interval (default: `500ms`).
    - `--provider-focus <string>`: Focus on specific provider (`google`, `microsoft`, `all`; default: `all`).
    - `--is-compact`: Render single-line compact terminal mode (default: `false`).
    - `--has-sound-alert`: Audio chime on spam placement (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedTickerEventRecord struct {
        EventID           string `json:"event_id"`
        ProviderName      string `json:"provider_name"`
        FolderPlacement   string `json:"folder_placement"`
        DeliveryTimestamp int64  `json:"delivery_timestamp"`
        IsPrimaryInbox    bool   `json:"is_primary_inbox"`
        HasAlertTriggered bool   `json:"has_alert_triggered"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedTickerEventRecord (
        event_id TEXT PRIMARY KEY,
        provider_name TEXT NOT NULL,
        folder_placement TEXT NOT NULL,
        delivery_timestamp INTEGER NOT NULL,
        is_primary_inbox INTEGER NOT NULL DEFAULT 1,
        has_alert_triggered INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics seed-ticker`. Feeds live events directly to TUI subscribers.

---

### Improvement 1539: Seed Mailbox IMAP/SMTP Credential Health & OAuth Token Expiry Gauge
- **ID:** `SEED-CREDENTIAL-GAUGE-1539`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Monitor authentication tokens, OAuth2 refresh lifecycles, and app-password validity across all seed accounts to eliminate silent credential failures during warmup runs.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Expiry gauges rendered as horizontal bars:
      ```text
      seed-google-01: [████████████████████] Valid (28d remaining)
      seed-m365-14:   [████                ] Expiring Soon (14h remaining) [WARN]
      seed-yahoo-09:  [                    ] Revoked / Bad Password [CRITICAL]
      ```
  - **Zero-CDN Web View:**
    - HTML `<progress>` elements styled via native CSS variables with zero font or script CDNs.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-credentials`
  - Flags:
    - `--warn-threshold <duration>`: Expiry threshold to trigger warning (default: `48h`).
    - `--auto-refresh`: Attempt OAuth2 token refresh automatically (default: `true`).
    - `--is-json`: Output results as JSON payload (default: `false`).
    - `--has-check-smtp`: Perform lightweight SMTP AUTH handshake probe (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedCredentialHealthRecord struct {
        AccountID         string `json:"account_id"`
        AuthType          string `json:"auth_type"`
        ExpiryTimestamp   int64  `json:"expiry_timestamp"`
        SecondsRemaining  int64  `json:"seconds_remaining"`
        IsValid           bool   `json:"is_valid"`
        HasRefreshFailure bool   `json:"has_refresh_failure"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedCredentialHealthRecord (
        account_id TEXT PRIMARY KEY,
        auth_type TEXT NOT NULL,
        expiry_timestamp INTEGER NOT NULL,
        seconds_remaining INTEGER NOT NULL,
        is_valid INTEGER NOT NULL DEFAULT 1,
        has_refresh_failure INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated with `gitmap email warmup diagnostics seed-credentials`. Uses local OS keyring or encrypted SQLite vault.

---

### Improvement 1540: Comprehensive Seed Diagnostics Self-Test & Diagnostic Export Bundle
- **ID:** `SEED-DIAGNOSTICS-BUNDLE-1540`
- **Cluster:** 154 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Run an end-to-end synthetic seed pool validation suite (DNS, IMAP, SMTP, latency, inbox placement, headers) and generate a compressed diagnostic bundle for root cause analysis.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Step-by-step diagnostic checklist with real-time pass/fail tick marks:
      ```text
      [*] Running Comprehensive Seed Diagnostics:
        [✓] DNS Resolver Alignment ......................... OK (12ms)
        [✓] IMAP Multi-Host Connectivity ................... OK (48ms)
        [✓] SMTP AUTH & Pipelining ......................... OK (29ms)
        [✓] SPF/DKIM Cryptographic Integrity .............. OK (100%)
        [!] Spam Folder Rescue Heuristics .................. 1 WARN (yahoo)
      >> Diagnostic Bundle Created: .ai-memory/diagnostics/seed-bundle-20261028.zip
      ```
  - **Zero-CDN Web View:**
    - Printable HTML diagnostic report with pure native CSS styles and SVG checkmarks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics seed-bundle`
  - Flags:
    - `--output-dir <string>`: Target path for output bundle (default: `.ai-memory/diagnostics/`).
    - `--include-raw-headers`: Include raw RFC 822 email headers in diagnostic dump (default: `false`).
    - `--is-full-suite`: Execute intensive stress tests across entire seed fleet (default: `true`).
    - `--has-anonymization`: Scrub user email prefixes before bundling (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedDiagnosticBundleRecord struct {
        BundleID          string `json:"bundle_id"`
        BundlePath        string `json:"bundle_path"`
        TestCount         int    `json:"test_count"`
        PassedTests       int    `json:"passed_tests"`
        FailedTests       int    `json:"failed_tests"`
        GeneratedAt       int64  `json:"generated_at"`
        IsComplete        bool   `json:"is_complete"`
        HasSevereFailures bool   `json:"has_severe_failures"`
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
        is_complete INTEGER NOT NULL DEFAULT 1,
        has_severe_failures INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics seed-bundle`. Stores output artifacts under relative repository path `.ai-memory/diagnostics/`.

---

# Cluster 155: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 1541–1550)

---

### Improvement 1541: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor Controller
- **ID:** `GOV-CONCURRENCY-CONTROLLER-1541`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Regulate outbound SMTP connection pools dynamically using inverted backoff models, adjusting simultaneous worker slots per destination MX domain based on real-time 421/451 deferral feedback.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live multi-column concurrency visualizer:
      ```text
      ┌─ [MULTI-DOMAIN CONCURRENCY GOVERNOR] ─────────────────────────────────────────────────────┐
      │ GLOBAL WORKERS: 128/256 | ACTIVE MX POOLS: 42 | TOTAL DEFERRALS (1m): 2 | HEALTH: 99.8%    │
      ├───────────────────────┬────────────┬─────────────┬──────────────┬─────────────┬────────────┤
      │ DESTINATION DOMAIN    │ MAX SLOTS  │ ACTIVE      │ BACKOFF MULT │ DEFER %     │ STATE      │
      ├───────────────────────┼────────────┼─────────────┼──────────────┼─────────────┼────────────┤
      │ gmail.com             │ 64         │ 48          │ 1.0x (NONE)  │ 0.0%        │ [●] OPTIMAL│
      │ outlook.com           │ 32         │ 12          │ 1.2x (LIGHT) │ 0.4%        │ [●] STABLE │
      │ yahoo.com             │ 16         │ 4           │ 2.5x (MED)   │ 3.2%        │ [▲] THROTTLE
      │ icloud.com            │ 8          │ 6           │ 1.0x (NONE)  │ 0.0%        │ [●] OPTIMAL│
      ├───────────────────────┴────────────┴─────────────┴──────────────┴─────────────┴────────────┤
      │ [SELECTED: yahoo.com] Slots throttled down from 32 to 16 due to 421 TS01 deferral spikes    │
      │ Keys: [↑/↓: Select Domain] [+: Manual Slot Boost] [-: Force Throttle] [r: Reset] [q: Quit] │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Server-rendered HTML5 dashboard with pure SVG dynamic capacity bars and CSS status indicators.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics governor-status`
  - Flags:
    - `--domain <string>`: Filter governor status by destination domain (default: `all`).
    - `--max-global-workers <int>`: Cap total concurrency across all domains (default: `256`).
    - `--is-auto-scale`: Allow dynamic allocation based on latency feedback (default: `true`).
    - `--has-strict-fairness`: Prevent high-volume domains from starving niche domains (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainConcurrencyStateRecord struct {
        DomainName         string  `json:"domain_name"`
        MaxAllowedSlots    int     `json:"max_allowed_slots"`
        ActiveActiveSlots  int     `json:"active_active_slots"`
        BackoffMultiplier  float64 `json:"backoff_multiplier"`
        DeferralPercentage float64 `json:"deferral_percentage"`
        IsThrottled        bool    `json:"is_throttled"`
        HasRecentDeferral  bool    `json:"has_recent_deferral"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainConcurrencyStateRecord (
        domain_name TEXT PRIMARY KEY,
        max_allowed_slots INTEGER NOT NULL DEFAULT 16,
        active_active_slots INTEGER NOT NULL DEFAULT 0,
        backoff_multiplier REAL NOT NULL DEFAULT 1.0,
        deferral_percentage REAL NOT NULL DEFAULT 0.0,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        has_recent_deferral INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `gitmap email warmup diagnostics governor-status`. Governs dispatch rates in real-time across workers managed by the lemon daemon.

---

### Improvement 1542: Additive Increase / Multiplicative Decrease (AIMD) Algorithm Telemetry HUD
- **ID:** `AIMD-TELEMETRY-HUD-1542`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Visualize TCP-style AIMD window adjustments for sending rates per MX domain, guaranteeing safe delivery throughput without triggering ISP volume filters.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-precision telemetry window showing saw-tooth congestion curve and window size ($W_{cwnd}$):
      ```text
      AIMD CONGESTION WINDOW: gmail.com
      W_cwnd = 48 pkts/sec | State: CONGESTION_AVOIDANCE (Additive +1/sec)
      Saw-Tooth: ⡠⠊⠑⠊⠑⠊⡠⠊⠑⡠⠊  Last Multiplicative Cut: -50% at 14:10:02 (Trigger: 452 4.2.2 Mailbox Full)
      ```
  - **Zero-CDN Web View:**
    - Inline SVG coordinate line chart with native CSS grid styling and zero JS dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics aimd-telemetry`
  - Flags:
    - `--domain <string>`: Target domain for telemetry focus (default: `gmail.com`).
    - `--history-window <duration>`: Historical curve duration to display (default: `1h`).
    - `--is-live-plot`: Render real-time streaming plot in terminal (default: `true`).
    - `--has-cut-highlight`: Emphasize multiplicative backoff cut events in red (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type AimdTelemetryRecord struct {
        WindowID          string `json:"window_id"`
        Domain            string `json:"domain"`
        CurrentCwnd       int    `json:"current_cwnd"`
        Ssthresh          int    `json:"ssthresh"`
        LastCutTimestamp  int64  `json:"last_cut_timestamp"`
        IsInSlowStart     bool   `json:"is_in_slow_start"`
        HasCongestionLoss bool   `json:"has_congestion_loss"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AimdTelemetryRecord (
        window_id TEXT PRIMARY KEY,
        domain TEXT NOT NULL,
        current_cwnd INTEGER NOT NULL DEFAULT 10,
        ssthresh INTEGER NOT NULL DEFAULT 64,
        last_cut_timestamp INTEGER NOT NULL,
        is_in_slow_start INTEGER NOT NULL DEFAULT 0,
        has_congestion_loss INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics aimd-telemetry`. Uses SQLite event streaming for ultra-low overhead state monitoring.

---

### Improvement 1543: Leaky Bucket Outbound Dispatch Shaper & Inter-Burst Jitter Governor
- **ID:** `LEAKY-BUCKET-DISPATCH-1543`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Smooth out bursty email dispatches using a leaky bucket algorithm with randomized inter-message micro-jitter (120ms–850ms) to emulate human sending distributions.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Graphical bucket fill indicator with leak rate readout:
      ```text
      ┌─ [LEAKY BUCKET DISPATCH SHAPER] ────────────────────────────────────────────────────────┐
      │ BUCKET CAPACITY: 100 msgs | CURRENT FILL: [████████░░░░░░░░░░░░] 40 msgs                 │
      │ LEAK RATE: 5.0 msgs/sec | JITTER RANGE: 150ms - 420ms (Gaussian: μ=285ms, σ=45ms)      │
      │ Output State: STEADY FLOW [●] | Queue Health: OPTIMAL | Dropped: 0                      │
      └─────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Pure CSS liquid fill animation for bucket state, embedded with vanilla HTML5 progress indicators.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics dispatch-shaper`
  - Flags:
    - `--capacity <int>`: Leaky bucket max capacity (default: `100`).
    - `--leak-rate <float>`: Dispatch leak rate in messages per second (default: `5.0`).
    - `--is-gaussian-jitter`: Apply normal Gaussian jitter distribution instead of uniform (default: `true`).
    - `--has-drop-prevention`: Dynamically backpressure upstream queue if bucket overflows (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LeakyBucketShaperRecord struct {
        ShaperID        string  `json:"shaper_id"`
        CurrentFill     int     `json:"current_fill"`
        Capacity        int     `json:"capacity"`
        CurrentLeakRate float64 `json:"current_leak_rate"`
        LastJitterMs    int     `json:"last_jitter_ms"`
        IsUnderPressure bool    `json:"is_under_pressure"`
        HasZeroLoss     bool    `json:"has_zero_loss"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS LeakyBucketShaperRecord (
        shaper_id TEXT PRIMARY KEY,
        current_fill INTEGER NOT NULL DEFAULT 0,
        capacity INTEGER NOT NULL DEFAULT 100,
        current_leak_rate REAL NOT NULL DEFAULT 5.0,
        last_jitter_ms INTEGER NOT NULL DEFAULT 200,
        is_under_pressure INTEGER NOT NULL DEFAULT 0,
        has_zero_loss INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates with `gitmap email warmup diagnostics dispatch-shaper`. Directly throttles outbound channels in `go-email-sender`.

---

### Improvement 1544: ISP Greylist Trap Classifier & Automatic Deferral Circuit Breaker
- **ID:** `GREYLIST-CIRCUIT-BREAKER-1544`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Distinguish temporary greylisting deferrals (451 4.7.1) from permanent policy blocks, automatically triggering non-destructive circuit breaker delays before retrying.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Circuit breaker tri-state display (`CLOSED [GREEN]`, `HALF-OPEN [YELLOW]`, `OPEN [RED]`):
      ```text
      CIRCUIT BREAKER: mx.zoho.com
      State: OPEN (Greylisting Trap Active) -> Hold Period: 180s (112s remaining)
      Auto-Probe Status: Scheduled at half-open threshold
      Classification: GREYLIST_TRANSIENT (Confidence: 99.4%)
      ```
  - **Zero-CDN Web View:**
    - High-contrast alert badge card with real-time countdown timer rendered via CSS styling.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics circuit-breaker`
  - Flags:
    - `--domain <string>`: Target MX domain to inspect (default: `all`).
    - `--hold-time <duration>`: Base hold time for greylisting deferrals (default: `300s`).
    - `--is-auto-recover`: Automatically probe in half-open state (default: `true`).
    - `--has-strict-classification`: Require exact regex match on SMTP response tokens (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CircuitBreakerRecord struct {
        BreakerID     string `json:"breaker_id"`
        TargetMxHost  string `json:"target_mx_host"`
        CurrentState  string `json:"current_state"`
        HoldSeconds   int    `json:"hold_seconds"`
        NextProbeTime int64  `json:"next_probe_time"`
        IsOpen        bool   `json:"is_open"`
        HasTripped    bool   `json:"has_tripped"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CircuitBreakerRecord (
        breaker_id TEXT PRIMARY KEY,
        target_mx_host TEXT NOT NULL UNIQUE,
        current_state TEXT NOT NULL DEFAULT 'CLOSED',
        hold_seconds INTEGER NOT NULL DEFAULT 300,
        next_probe_time INTEGER NOT NULL DEFAULT 0,
        is_open INTEGER NOT NULL DEFAULT 0,
        has_tripped INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked by `gitmap email warmup diagnostics circuit-breaker`. Stores circuit states in `repodb/pipeline.db`.

---

### Improvement 1545: Multi-Tenant Pool Concurrency Fair-Share Arbiter
- **ID:** `MULTI-TENANT-ARBITER-1545`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Allocate dispatch bandwidth proportionally across competing warmup campaigns and sub-accounts using weighted fair-queuing (WFQ) to prevent tenant starvation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Tenant allocation breakdown table:
      ```text
      ┌─ [MULTI-TENANT FAIR-SHARE ARBITER] ───────────────────────────────────────────────────────┐
      │ TENANT ID        │ WEIGHT │ QUOTA / SEC │ CURRENT USED │ DEFICIT BAL │ STATUS             │
      ├──────────────────┼────────┼─────────────┼──────────────┼─────────────┼────────────────────┤
      │ tenant-prod-core │ 50%    │ 25.0 msgs/s │ 24.8 msgs/s  │ 0           │ [●] SATISFIED      │
      │ tenant-warmup-a  │ 30%    │ 15.0 msgs/s │ 12.1 msgs/s  │ +2          │ [●] SATISFIED      │
      │ tenant-warmup-b  │ 20%    │ 10.0 msgs/s │ 10.0 msgs/s  │ -4          │ [▲] THROTTLED      │
      └──────────────────┴────────┴─────────────┴──────────────┴─────────────┴────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Pure CSS horizontal stacked percentage bars with embedded theme colors.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics fair-share`
  - Flags:
    - `--rebalance-interval <duration>`: Frequency of WFQ weight recalculation (default: `5s`).
    - `--enforce-hard-cap`: Strictly drop bursts exceeding calculated share (default: `false`).
    - `--is-audit-enabled`: Persist allocation logs to SQLite (default: `true`).
    - `--has-priority-boost`: Permit temporary deficit borrowing for urgent seeds (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FairShareTenantRecord struct {
        TenantID         string  `json:"tenant_id"`
        AssignedWeight   float64 `json:"assigned_weight"`
        AllocatedRate    float64 `json:"allocated_rate"`
        ActualRate       float64 `json:"actual_rate"`
        DeficitBalance   int     `json:"deficit_balance"`
        IsThrottled      bool    `json:"is_throttled"`
        HasPriorityBoost bool    `json:"has_priority_boost"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FairShareTenantRecord (
        tenant_id TEXT PRIMARY KEY,
        assigned_weight REAL NOT NULL DEFAULT 1.0,
        allocated_rate REAL NOT NULL DEFAULT 10.0,
        actual_rate REAL NOT NULL DEFAULT 0.0,
        deficit_balance INTEGER NOT NULL DEFAULT 0,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        has_priority_boost INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates with `gitmap email warmup diagnostics fair-share` to govern internal channel priority dispatch.

---

### Improvement 1546: Dynamic MX Pipelining Handshake Latency Optimization Engine
- **ID:** `MX-PIPELINING-LATENCY-1546`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Accelerate SMTP delivery cycles by negotiating RFC 2920 command pipelining, measuring TLS handshake overhead, and reusing persistent session sockets.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Command pipeline timeline:
      ```text
      [SMTP PIPELINING: mx1.google.com]
      Handshake Breakdown: TCP SYN: 12ms | TLS 1.3: 18ms | EHLO+MAIL+RCPT (Pipelined): 14ms
      Total Round-Trip Optimization: -54% Latency Reduction (Old: 98ms -> New: 44ms) [●]
      ```
  - **Zero-CDN Web View:**
    - Waterfall latency diagram rendered via embedded SVG syntax.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics pipelining-status`
  - Flags:
    - `--mx-target <string>`: Hostname of MX server to probe (default: `gmail-smtp-in.l.google.com`).
    - `--keep-alive-timeout <duration>`: Socket reuse timeout (default: `30s`).
    - `--is-tls-early-data`: Enable TLS 0-RTT if supported (default: `false`).
    - `--has-pipelining-fallback`: Fall back to sequential mode on legacy MTA (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MxPipeliningRecord struct {
        MxHostname       string  `json:"mx_hostname"`
        TcpConnectMs     float64 `json:"tcp_connect_ms"`
        TlsHandshakeMs   float64 `json:"tls_handshake_ms"`
        PipelineSavePct  float64 `json:"pipeline_save_pct"`
        SupportsPipeline bool    `json:"supports_pipeline"`
        IsSessionReused  bool    `json:"is_session_reused"`
        HasTls13         bool    `json:"has_tls13"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxPipeliningRecord (
        mx_hostname TEXT PRIMARY KEY,
        tcp_connect_ms REAL NOT NULL DEFAULT 0.0,
        tls_handshake_ms REAL NOT NULL DEFAULT 0.0,
        pipeline_save_pct REAL NOT NULL DEFAULT 0.0,
        supports_pipeline INTEGER NOT NULL DEFAULT 1,
        is_session_reused INTEGER NOT NULL DEFAULT 0,
        has_tls13 INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics pipelining-status`. Connects directly to network metrics telemetry in `repodb/pipeline.db`.

---

### Improvement 1547: Outbound IP Warmup Phase Ramp-Up Governor & Ceiling Automator
- **ID:** `IP-WARMUP-RAMP-CEILING-1547`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Automate progressive daily volume increases for new outbound dedicated IPs across a 30-day warmup curve, clamping daily ceilings if error spikes occur.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Daily warmup stage progress:
      ```text
      ┌─ [OUTBOUND IP WARMUP CEILING: 198.51.100.24] ─────────────────────────────────────────────┐
      │ CURRENT PHASE: Day 12 / 30 | TODAY'S CEILING: 1,200 msgs | SENT: 840 msgs (70% Used)      │
      │ 30-DAY RAMP: [████████████░░░░░░░░░░░░░░░░░░] Day 12                                      │
      │ Tomorrow Projection: 1,500 msgs | Error Gating: 0.1% (<0.5% threshold) -> APPROVED [●]   │
      └───────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Step-wise progression chart with interactive milestone badges.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics ip-ramp`
  - Flags:
    - `--ip-address <string>`: Target IP address to inspect or govern.
    - `--ramp-curve <string>`: Strategy curve (`linear`, `exponential`, `conservative`; default: `conservative`).
    - `--is-auto-advance`: Automatically increase ceiling at midnight if healthy (default: `true`).
    - `--has-halt-on-bounce`: Freeze ramp if bounce rate exceeds 1% (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type IpWarmupPhaseRecord struct {
        IpAddress         string  `json:"ip_address"`
        CurrentDay        int     `json:"current_day"`
        DailyCeilingLimit int     `json:"daily_ceiling_limit"`
        SentTodayCount    int     `json:"sent_today_count"`
        ErrorRatePercent  float64 `json:"error_rate_percent"`
        IsCeilingClamped  bool    `json:"is_ceiling_clamped"`
        HasAdvancedPhase  bool    `json:"has_advanced_phase"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS IpWarmupPhaseRecord (
        ip_address TEXT PRIMARY KEY,
        current_day INTEGER NOT NULL DEFAULT 1,
        daily_ceiling_limit INTEGER NOT NULL DEFAULT 50,
        sent_today_count INTEGER NOT NULL DEFAULT 0,
        error_rate_percent REAL NOT NULL DEFAULT 0.0,
        is_ceiling_clamped INTEGER NOT NULL DEFAULT 0,
        has_advanced_phase INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics ip-ramp`. Synchronizes dispatch caps in `repodb/pipeline.db`.

---

### Improvement 1548: Multi-IP Round-Robin Pool Affinity & Sticky Destination Binding
- **ID:** `MULTI-IP-AFFINITY-BINDING-1548`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Maintain persistent source-IP-to-destination-domain affinity bindings while load balancing across IP pools to prevent ISP suspicion caused by sudden IP switching.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Affinity routing map:
      ```text
      BINDING TABLE:
      gmail.com     -> Pool [198.51.100.21, 198.51.100.22] (Sticky Hash: 0x9A4F) [ACTIVE]
      microsoft.com -> Pool [198.51.100.23]                (Sticky Hash: 0x2C1B) [ACTIVE]
      yahoo.com     -> Pool [198.51.100.24]                (Sticky Hash: 0xE801) [WARMING]
      ```
  - **Zero-CDN Web View:**
    - Visual node topology graph using embedded CSS flex-flow and SVG link connectors.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics ip-affinity`
  - Flags:
    - `--affinity-mode <string>`: Strategy (`sticky-domain`, `round-robin`, `least-busy`; default: `sticky-domain`).
    - `--ttl <duration>`: Binding persistence timeout (default: `24h`).
    - `--is-pool-rebalanced`: Trigger auto-rebalancing across healthy IPs (default: `true`).
    - `--has-fallback-ip`: Fallback to secondary IP if primary IP suffers high deferrals (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type IpAffinityBindingRecord struct {
        BindingKey        string `json:"binding_key"`
        DestinationDomain string `json:"destination_domain"`
        AssignedIp        string `json:"assigned_ip"`
        BindingExpiresAt  int64  `json:"binding_expires_at"`
        IsActive          bool   `json:"is_active"`
        HasFallback       bool   `json:"has_fallback"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS IpAffinityBindingRecord (
        binding_key TEXT PRIMARY KEY,
        destination_domain TEXT NOT NULL,
        assigned_ip TEXT NOT NULL,
        binding_expires_at INTEGER NOT NULL,
        is_active INTEGER NOT NULL DEFAULT 1,
        has_fallback INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked by `gitmap email warmup diagnostics ip-affinity`. Directly read by outbound socket dialers.

---

### Improvement 1549: Transient SMTP Error Cluster Classification & Feedback Loop
- **ID:** `SMTP-ERROR-CLUSTER-1549`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Group and classify transient 4xx and 5xx SMTP server response strings into actionable categories (rate limit, reputation drop, content block, DNS failure) using deterministic regex classifiers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Error taxonomy breakdown pane:
      ```text
      ┌─ [SMTP ERROR CLUSTER CLASSIFIER] ─────────────────────────────────────────────────────────┐
      │ RECENT ERRORS (15m): 42 | RATE_LIMIT: 28 (66%) | POLICY_BLOCK: 12 (28%) | DNS_FAIL: 2 (4%) │
      ├───────────────────────┬────────┬──────────────────────────────────────────┬───────────────┤
      │ CLUSTER CATEGORY      │ COUNT  │ SAMPLE RESPONSE                          │ ACTION        │
      ├───────────────────────┼────────┼──────────────────────────────────────────┼───────────────┤
      │ RATE_LIMIT_421        │ 28     │ 421 4.7.0 Try again later, closing conn. │ Invert Backoff│
      │ POLICY_REPUTATION_550 │ 12     │ 550 5.7.1 Message rejected due to policy │ Isolate Seed  │
      │ DNS_RESOLVER_451      │ 2      │ 451 4.4.3 Temporary server error         │ Retry 60s     │
      └────────────────────────────────┴──────────────────────────────────────────┴───────────────┘
      ```
  - **Zero-CDN Web View:**
    - Severity badge matrix with code-block previews of raw SMTP lines.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics error-clusters`
  - Flags:
    - `--window <duration>`: Time window to cluster errors (default: `1h`).
    - `--min-cluster-size <int>`: Minimum errors to form a cluster (default: `3`).
    - `--is-auto-remedy`: Trigger automated governor adjustment based on cluster type (default: `true`).
    - `--has-raw-samples`: Include full SMTP transcript lines in output (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SmtpErrorClusterRecord struct {
        ClusterID           string `json:"cluster_id"`
        ClusterCategory     string `json:"cluster_category"`
        OccurrenceCount     int    `json:"occurrence_count"`
        PatternRegex        string `json:"pattern_regex"`
        RecommendedAction   string `json:"recommended_action"`
        IsActionApplied     bool   `json:"is_action_applied"`
        HasCriticalSeverity bool   `json:"has_critical_severity"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpErrorClusterRecord (
        cluster_id TEXT PRIMARY KEY,
        cluster_category TEXT NOT NULL,
        occurrence_count INTEGER NOT NULL DEFAULT 1,
        pattern_regex TEXT NOT NULL,
        recommended_action TEXT NOT NULL,
        is_action_applied INTEGER NOT NULL DEFAULT 0,
        has_critical_severity INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics error-clusters`. Directly triggers backoff updates.

---

### Improvement 1550: Global Governor Configuration Snapshot & Dynamic Re-Arm Engine
- **ID:** `GOV-SNAPSHOT-REARM-1550`
- **Cluster:** 155 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Export, import, and dynamically re-arm the complete state of the concurrency governor and backoff models without dropping active TCP connections or restarting daemon processes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Configuration state snapshot confirmation:
      ```text
      [GOVERNOR SNAPSHOT SAVED: .ai-memory/snapshots/gov-state-20261028-1425.json]
      Captured State: 64 Domain Rules, 8 IP Pools, 128 Active Leases
      [RE-ARM STATUS]: Zero-downtime hot-reload completed in 1.4ms [●]
      ```
  - **Zero-CDN Web View:**
    - Interactive JSON tree viewer built with native HTML `<details>` and zero JavaScript libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics governor-rearm`
  - Flags:
    - `--snapshot-file <string>`: Path to snapshot file for restore or export.
    - `--action <string>`: Action to perform (`save`, `restore`, `dry-run`; default: `save`).
    - `--is-hot-reload`: Apply changes live without restarting the background lemon service (default: `true`).
    - `--has-checksum-verify`: Verify SHA-256 integrity hash of snapshot before loading (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorSnapshotRecord struct {
        SnapshotID         string `json:"snapshot_id"`
        SnapshotChecksum   string `json:"snapshot_checksum"`
        ActiveDomainCount  int    `json:"active_domain_count"`
        ActiveLeaseCount   int    `json:"active_lease_count"`
        CreatedAtTimestamp int64  `json:"created_at_timestamp"`
        IsVerified         bool   `json:"is_verified"`
        HasAppliedSuccess  bool   `json:"has_applied_success"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorSnapshotRecord (
        snapshot_id TEXT PRIMARY KEY,
        snapshot_checksum TEXT NOT NULL,
        active_domain_count INTEGER NOT NULL DEFAULT 0,
        active_lease_count INTEGER NOT NULL DEFAULT 0,
        created_at_timestamp INTEGER NOT NULL,
        is_verified INTEGER NOT NULL DEFAULT 1,
        has_applied_success INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics governor-rearm`. Persists snapshots to `.ai-memory/snapshots/` adhering to strictly relative paths.

---

# Cluster 156: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 1551–1560)

---

### Improvement 1551: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD
- **ID:** `HTML-SANDBOX-SANITIZER-1551`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Provide an isolated, zero-CDN HTML/CSS parsing and sanitization engine that verifies outbound warmup templates against strict email client CSS subsets (Gmail, Outlook, Apple Mail), stripping scripts, insecure fonts, and external assets.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Sanitizer terminal HUD with syntax compliance split-pane:
      ```text
      ┌─ [ZERO-CDN EMAIL HTML SANITIZER & SECURITY SANDBOX] ──────────────────────────────────────┐
      │ TEMPLATE: welcome-warmup-v2.html | SIZE: 44.2 KB | CSS COMPLIANCE: 99.4% | SCRIPTS: 0     │
      ├───────────────────────────────────┬───────────────────────────────────────────────────────┤
      │ SECURITY SANITIZATION CHECKS      │ COMPATIBILITY MATRIX                                  │
      ├───────────────────────────────────┼───────────────────────────────────────────────────────┤
      │ [✓] Zero Remote Script Tags       │ Apple Mail (iOS/macOS): 100% [●]                      │
      │ [✓] Zero External CSS CDNs        │ Gmail (Web/Android):    98%  [●]                      │
      │ [✓] Strict Image Source Hostnames │ Outlook 365 (Desktop):  94%  [▲] (MSO VML Active)     │
      │ [✓] SVG Secure Elements Only      │ Yahoo Mail:             99%  [●]                      │
      ├───────────────────────────────────┴───────────────────────────────────────────────────────┤
      │ [OK] Template passes all zero-CDN security boundaries and is ready for dispatch.           │
      │ Keys: [v: View Diff] [s: Save Sanitized] [p: Preview Live] [q: Quit]                      │
      └───────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Sandboxed `<iframe>` with `sandbox="allow-same-origin"` and zero external stylesheet or JavaScript references, rendering pure HTML templates locally.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics html-sandbox`
  - Flags:
    - `--template-path <string>`: Relative path to HTML email template.
    - `--target-client <string>`: Target client profile (`gmail`, `outlook`, `apple`, `all`; default: `all`).
    - `--is-auto-sanitize`: Strip invalid tags automatically and output clean HTML (default: `true`).
    - `--has-strict-csp`: Enforce zero external domain connections in CSS (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HtmlSanitizerRecord struct {
        SanitizerID       string  `json:"sanitizer_id"`
        TemplateName      string  `json:"template_name"`
        OriginalSizeBytes int     `json:"original_size_bytes"`
        CleanedSizeBytes  int     `json:"cleaned_size_bytes"`
        ComplianceScore   float64 `json:"compliance_score"`
        RemovedTagCount   int     `json:"removed_tag_count"`
        IsSecure          bool    `json:"is_secure"`
        HasExternalAssets bool    `json:"has_external_assets"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlSanitizerRecord (
        sanitizer_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        original_size_bytes INTEGER NOT NULL DEFAULT 0,
        cleaned_size_bytes INTEGER NOT NULL DEFAULT 0,
        compliance_score REAL NOT NULL DEFAULT 100.0,
        removed_tag_count INTEGER NOT NULL DEFAULT 0,
        is_secure INTEGER NOT NULL DEFAULT 1,
        has_external_assets INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics html-sandbox`. Intercepts all templates prior to signing in `go-email-sender`.

---

### Improvement 1552: Email CSS Media Query & Mobile Responsive Grid Auditing Engine
- **ID:** `EMAIL-RESPONSIVE-GRID-1552`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit inline styles and `@media` queries in email templates to guarantee seamless rendering across 320px mobile viewports without horizontal scrolling or distorted column stacks.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Viewport column stack simulation:
      ```text
      VIEWPORT AUDIT:
      [320px Mobile]: Stacked Single-Column OK [✓] | Font Min: 14px OK [✓]
      [600px Tablet]: 2-Column Responsive Flow OK [✓] | Image Fluidity: 100% OK [✓]
      [960px Desktop]: Centered 600px Max-Width OK [✓]
      Result: 0 Horizontal Overflow Breaches Detected
      ```
  - **Zero-CDN Web View:**
    - Responsive viewport toggle toolbar (320px, 480px, 600px) styled purely with CSS radio inputs.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics responsive-audit`
  - Flags:
    - `--template-file <string>`: Path to HTML template to audit.
    - `--min-font-size <int>`: Minimum mobile readable font size in pixels (default: `14`).
    - `--is-fluid-images`: Require all `<img>` tags to include `max-width: 100%` (default: `true`).
    - `--has-mso-conditional`: Verify Microsoft Outlook `<!--[if mso]>` tables (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ResponsiveAuditRecord struct {
        AuditID               string  `json:"audit_id"`
        TemplatePath          string  `json:"template_path"`
        MinViewportWidth      int     `json:"min_viewport_width"`
        MaxHorizontalWidth    int     `json:"max_horizontal_width"`
        MobileFontScore       float64 `json:"mobile_font_score"`
        IsMobileCompliant     bool    `json:"is_mobile_compliant"`
        HasHorizontalOverflow bool    `json:"has_horizontal_overflow"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ResponsiveAuditRecord (
        audit_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        min_viewport_width INTEGER NOT NULL DEFAULT 320,
        max_horizontal_width INTEGER NOT NULL DEFAULT 600,
        mobile_font_score REAL NOT NULL DEFAULT 100.0,
        is_mobile_compliant INTEGER NOT NULL DEFAULT 1,
        has_horizontal_overflow INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked by `gitmap email warmup diagnostics responsive-audit` and verifies template hygiene.

---

### Improvement 1553: Tracking Pixel Decoupler & Local SVG Beacon Transcoder
- **ID:** `TRACKING-PIXEL-DECOUPLER-1553`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Transcode remote 1x1 tracking GIFs into privacy-preserving zero-CDN inline SVG beacons with cryptographically signed tokens to prevent external CDN leaks.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Beacon transcoder status:
      ```text
      TRACKING BEACON TRANSCODER:
      Remote HTTP Pixel Replaced: https://track.external.com/pixel.gif -> INLINE_SVG_BEACON
      Beacon Hash: sha256:d8a4e... | External DNS Leak: ZERO [●]
      ```
  - **Zero-CDN Web View:**
    - Embedded SVG beacon preview inspection box showing base64 representation.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics pixel-decouple`
  - Flags:
    - `--template <string>`: Target template to process.
    - `--token-salt <string>`: Cryptographic salt for beacon token generation.
    - `--is-inline-svg`: Embed as raw data URI SVG instead of external URL (default: `true`).
    - `--has-privacy-preserving`: Strip client user-agent from beacon tracking payload (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PixelDecouplerRecord struct {
        BeaconID       string `json:"beacon_id"`
        TemplateName   string `json:"template_name"`
        OriginalUrl    string `json:"original_url"`
        InlineSvgToken string `json:"inline_svg_token"`
        IsDecoupled    bool   `json:"is_decoupled"`
        HasRemoteLeak  bool   `json:"has_remote_leak"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PixelDecouplerRecord (
        beacon_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        original_url TEXT NOT NULL,
        inline_svg_token TEXT NOT NULL,
        is_decoupled INTEGER NOT NULL DEFAULT 1,
        has_remote_leak INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated via `gitmap email warmup diagnostics pixel-decouple`.

---

### Improvement 1554: Dark Mode CSS Variable Inverter & Color Contrast Validator
- **ID:** `DARK-MODE-CONTRAST-1554`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit and inject `@media (prefers-color-scheme: dark)` styles into email HTML templates to prevent Apple Mail and Outlook from arbitrarily inverting text colors into unreadable low-contrast states.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual contrast preview table:
      ```text
      CONTRAST VALIDATION (WCAG AAA >= 7.0:1):
      Light Mode: #1e293b on #ffffff -> Contrast Ratio: 12.8:1 [PASS]
      Dark Mode:  #f8fafc on #0f172a -> Contrast Ratio: 14.1:1 [PASS]
      Auto-Injected Dark Meta Tags: [✓] color-scheme: light dark [✓] supported-color-schemes: light dark
      ```
  - **Zero-CDN Web View:**
    - Side-by-side Light/Dark preview split rendered using pure CSS `@media (prefers-color-scheme: dark)`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics dark-mode`
  - Flags:
    - `--template <string>`: Path to HTML template.
    - `--min-contrast-ratio <float>`: Target WCAG contrast ratio (default: `4.5`).
    - `--is-auto-inject`: Automatically inject dark-mode CSS overrides (default: `true`).
    - `--has-image-transparency`: Check PNG/SVG logos for dark backgrounds (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeContrastRecord struct {
        AuditID             string  `json:"audit_id"`
        TemplateName        string  `json:"template_name"`
        LightModeContrast   float64 `json:"light_mode_contrast"`
        DarkModeContrast    float64 `json:"dark_mode_contrast"`
        IsWcagCompliant     bool    `json:"is_wcag_compliant"`
        HasInjectedOverride bool    `json:"has_injected_override"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DarkModeContrastRecord (
        audit_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        light_mode_contrast REAL NOT NULL DEFAULT 12.0,
        dark_mode_contrast REAL NOT NULL DEFAULT 12.0,
        is_wcag_compliant INTEGER NOT NULL DEFAULT 1,
        has_injected_override INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics dark-mode`.

---

### Improvement 1555: Font Subsetting & Inline Base64 WOFF2 Bundler
- **ID:** `FONT-SUBSET-BUNDLER-1555`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Subset custom web typography to only the characters used in the email template, encode them into base64 WOFF2 strings, and embed them directly to eliminate external Google Fonts CDN requests.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Subsetter efficiency summary:
      ```text
      FONT SUBSETTER & INLINE BUNDLER:
      Original Font: Inter-Regular.woff2 (142 KB)
      Glyphs Retained: 78 unique characters -> Subset Font: 6.4 KB (-95.5%)
      Embedded as data:font/woff2;base64,... | Remote Font Calls: 0 [●]
      ```
  - **Zero-CDN Web View:**
    - Typography specimen card rendered with embedded font subset.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics font-bundler`
  - Flags:
    - `--template <string>`: Email template containing custom typography.
    - `--font-file <string>`: Source TTF/WOFF2 font path to subset.
    - `--is-strip-ligatures`: Omit advanced OpenType ligatures for minimal size (default: `true`).
    - `--has-fallback-stack`: Guarantee standard web-safe system font fallbacks (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FontSubsetRecord struct {
        SubsetID       string `json:"subset_id"`
        FontFamilyName string `json:"font_family_name"`
        OriginalBytes  int    `json:"original_bytes"`
        SubsetBytes    int    `json:"subset_bytes"`
        RetainedGlyphs int    `json:"retained_glyphs"`
        IsEmbedded     bool   `json:"is_embedded"`
        HasRemoteCalls bool   `json:"has_remote_calls"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FontSubsetRecord (
        subset_id TEXT PRIMARY KEY,
        font_family_name TEXT NOT NULL,
        original_bytes INTEGER NOT NULL DEFAULT 0,
        subset_bytes INTEGER NOT NULL DEFAULT 0,
        retained_glyphs INTEGER NOT NULL DEFAULT 0,
        is_embedded INTEGER NOT NULL DEFAULT 1,
        has_remote_calls INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated via `gitmap email warmup diagnostics font-bundler`.

---

### Improvement 1556: Email Payload Size & 102KB Gmail Clipping Guard
- **ID:** `PAYLOAD-GMAIL-CLIPPING-1556`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Measure total MIME payload bytes (HTML, inline CSS, multipart boundaries, headers) and enforce a strict safety threshold below 102KB to prevent Gmail from clipping messages with "[Message clipped] View entire message".
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Size meter gauge:
      ```text
      ┌─ [GMAIL 102KB CLIPPING GUARD] ────────────────────────────────────────────────────────────┐
      │ TOTAL MIME SIZE: 68.4 KB / 102.0 KB [███████████████░░░░░░░] 67% (SAFE: 33.6 KB Margin)   │
      │ Breakdown: HTML Body: 42.1 KB | Inline Styles: 14.8 KB | Headers: 2.2 KB | Images: 9.3 KB  │
      │ Clipping Risk: NONE [●] | Recommended Minification Action: None required                  │
      └───────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Dynamic payload budget meter styled via CSS gradient fills.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics clipping-guard`
  - Flags:
    - `--template <string>`: Target HTML template to evaluate.
    - `--max-bytes <int>`: Clipping ceiling in bytes (default: `102400`).
    - `--is-auto-minify`: Automatically minify HTML whitespace if within 5% of cap (default: `true`).
    - `--has-warn-alert`: Return exit code 1 if payload breaches threshold (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ClippingGuardRecord struct {
        TemplateID     string `json:"template_id"`
        TotalSizeBytes int    `json:"total_size_bytes"`
        RemainingBytes int    `json:"remaining_bytes"`
        IsUnderCeiling bool   `json:"is_under_ceiling"`
        HasMinified    bool   `json:"has_minified"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ClippingGuardRecord (
        template_id TEXT PRIMARY KEY,
        total_size_bytes INTEGER NOT NULL DEFAULT 0,
        remaining_bytes INTEGER NOT NULL DEFAULT 0,
        is_under_ceiling INTEGER NOT NULL DEFAULT 1,
        has_minified INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics clipping-guard`. Blocks templates that exceed limits prior to dispatch.

---

### Improvement 1557: Phishing & Malicious Anchor Link Disinfection Analyzer
- **ID:** `PHISHING-LINK-DISINFECT-1557`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan all hyperlink destinations (`<a href="...">`) in templates against reputation blocklists, verify destination domain alignment with sender DKIM domains, and flag misleading anchor text.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Anchor link audit list:
      ```text
      ANCHOR AUDIT: 4 Links Inspected
      [✓] https://brand.com/app          (Anchor: "Go to Dashboard") -> MATCH [●]
      [✓] https://brand.com/unsubscribe  (Anchor: "Unsubscribe")    -> RFC 8058 Compliant [●]
      [!] http://insecure-cdn.com/asset  (Anchor: "View Report")    -> INSECURE PROTOCOL (HTTP) [▲]
      ```
  - **Zero-CDN Web View:**
    - Link verification scorecard with status badges.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics link-disinfect`
  - Flags:
    - `--template <string>`: Template file to inspect.
    - `--enforce-https`: Flag any non-HTTPS URLs as critical failures (default: `true`).
    - `--is-check-rbl`: Query local SURBL/URIBL mirrors for link domains (default: `true`).
    - `--has-disinfect`: Rewrite insecure links to sanitized placeholder paths (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LinkDisinfectionRecord struct {
        LinkID            string `json:"link_id"`
        TargetUrl         string `json:"target_url"`
        AnchorText        string `json:"anchor_text"`
        ProtocolType      string `json:"protocol_type"`
        IsHttpsSecure     bool   `json:"is_https_secure"`
        HasReputationRisk bool   `json:"has_reputation_risk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS LinkDisinfectionRecord (
        link_id TEXT PRIMARY KEY,
        target_url TEXT NOT NULL,
        anchor_text TEXT NOT NULL,
        protocol_type TEXT NOT NULL DEFAULT 'https',
        is_https_secure INTEGER NOT NULL DEFAULT 1,
        has_reputation_risk INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked via `gitmap email warmup diagnostics link-disinfect`.

---

### Improvement 1558: Pre-Header & Snippet Text Optimization Live Inspector
- **ID:** `PREHEADER-SNIPPET-INSPECT-1558`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Inspect and optimize hidden pre-header text and zero-width non-breaking space padding (`&zwnj;&nbsp;`) to ensure clean inbox preview snippets without bleeding into body copy.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live inbox preview simulation:
      ```text
      INBOX PREVIEW SNIPPET:
      Sender: Acme Delivery <outreach@acme.com>
      Subject: Your weekly performance summary
      Snippet: [Ready for review: Q3 deliverability metrics are up 14%...] (Padding: 120 zwnj chars)
      Snippet Bleed Risk: ZERO [●]
      ```
  - **Zero-CDN Web View:**
    - Webmail inbox list item simulator replicating desktop and mobile snippet widths.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics preheader-inspect`
  - Flags:
    - `--template <string>`: Path to HTML template.
    - `--min-padding <int>`: Minimum zero-width space characters to append (default: `90`).
    - `--is-auto-inject`: Inject missing pre-header container automatically (default: `true`).
    - `--has-subject-preview`: Render combined Subject + Snippet layout (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PreheaderInspectRecord struct {
        TemplateID       string `json:"template_id"`
        ExtractedSnippet string `json:"extracted_snippet"`
        PaddingLength    int    `json:"padding_length"`
        IsOptimized      bool   `json:"is_optimized"`
        HasSnippetBleed  bool   `json:"has_snippet_bleed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PreheaderInspectRecord (
        template_id TEXT PRIMARY KEY,
        extracted_snippet TEXT NOT NULL,
        padding_length INTEGER NOT NULL DEFAULT 0,
        is_optimized INTEGER NOT NULL DEFAULT 1,
        has_snippet_bleed INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics preheader-inspect`.

---

### Improvement 1559: Responsive Table Nesting & MSO Outlook VML Conditioner
- **ID:** `MSO-OUTLOOK-VML-CONDITIONER-1559`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically validate and condition nested HTML table structures with Microsoft Office Vector Markup Language (`<!--[if mso]> <v:roundrect>`) wrappers for pixel-perfect Outlook rendering.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - MSO validation breakdown:
      ```text
      MSO OUTLOOK CONDITIONER:
      [✓] Outlook Word Engine Compatibility Verified
      [✓] Conditional Button VML Wrapper Injected (<v:roundrect arcsize="10%">)
      [✓] Ghost Table Columns Aligned: 600px Max-Width Guaranteed [●]
      ```
  - **Zero-CDN Web View:**
    - DOM tree visualizer highlighting conditional comments and nested table borders.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics mso-conditioner`
  - Flags:
    - `--template <string>`: Path to email template.
    - `--outlook-dpi <int>`: DPI scaling target (default: `96`).
    - `--is-wrap-buttons`: Wrap bulletproof buttons in VML rect containers (default: `true`).
    - `--has-ghost-tables`: Insert conditional table scaffolding around divs (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MsoConditionerRecord struct {
        ConditionerID    string `json:"conditioner_id"`
        TemplateName     string `json:"template_name"`
        VmlButtonCount   int    `json:"vml_button_count"`
        GhostTableCount  int    `json:"ghost_table_count"`
        IsOutlookReady   bool   `json:"is_outlook_ready"`
        HasMsoValidation bool   `json:"has_mso_validation"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MsoConditionerRecord (
        conditioner_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        vml_button_count INTEGER NOT NULL DEFAULT 0,
        ghost_table_count INTEGER NOT NULL DEFAULT 0,
        is_outlook_ready INTEGER NOT NULL DEFAULT 1,
        has_mso_validation INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics mso-conditioner`.

---

### Improvement 1560: Zero-CDN Template Export & Cryptographic Fingerprint Registry
- **ID:** `TEMPLATE-EXPORT-FINGERPRINT-1560`
- **Cluster:** 156 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Generate deterministic cryptographic hashes (SHA-256) of verified zero-CDN templates, registering them into a tamper-proof local SQLite catalog before deployment to production warmup queues.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Template registry certification card:
      ```text
      ┌─ [TEMPLATE REGISTRY & CRYPTOGRAPHIC CERTIFICATION] ────────────────────────────────────────┐
      │ TEMPLATE: marketing-onboarding-v3.html                                                    │
      │ SHA-256: 7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069                │
      │ Content Status: ZERO-CDN VERIFIED [●] | Security Sign-Off: PASS | Version: 3.1.0          │
      │ Registered to: repodb/pipeline.db -> TemplateFingerprintRegistry                           │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Cryptographic catalog table with export buttons, copyable hash badges, and verification status.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup diagnostics template-register`
  - Flags:
    - `--template-file <string>`: Path to finalized HTML template.
    - `--export-dir <string>`: Directory to export certified bundle (default: `02-spec/21-app/templates/`).
    - `--is-sign`: Cryptographically sign template with local Ed25519 sender key (default: `true`).
    - `--has-overwrite`: Overwrite existing registry record if hash matches (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TemplateFingerprintRecord struct {
        TemplateHashID    string `json:"template_hash_id"`
        TemplateName      string `json:"template_name"`
        Sha256Checksum    string `json:"sha256_checksum"`
        ByteSize          int    `json:"byte_size"`
        RegisteredAtUnix  int64  `json:"registered_at_unix"`
        IsVerified        bool   `json:"is_verified"`
        HasValidSignature bool   `json:"has_valid_signature"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TemplateFingerprintRecord (
        template_hash_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        sha256_checksum TEXT NOT NULL UNIQUE,
        byte_size INTEGER NOT NULL DEFAULT 0,
        registered_at_unix INTEGER NOT NULL,
        is_verified INTEGER NOT NULL DEFAULT 1,
        has_valid_signature INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup diagnostics template-register`. Writes certified records directly to `repodb/pipeline.db`.
