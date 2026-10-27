# Specification 351: Improvements 1431 to 1460 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 351
- **Slug:** `351-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/03-improvements-1431-to-1460-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/349-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (1431 - 1460)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 144** | Items 1431–1440 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 145** | Items 1441–1450 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 146** | Items 1451–1460 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 144: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 1431–1440)

---

### Improvement 1431: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Topology Visualizer
- **ID:** `SEED-MATRIX-TOPOLOGY-1431`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a real-time, interactive Bubble Tea terminal grid visualizing distributed seed mailbox nodes, health metrics, Anycast edge routing latencies, and active deliverability cohorts across Gmail, Microsoft 365, Yahoo, and iCloud seed pools.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [BUBBLE TEA SEED POOL HEALTH MATRIX & TOPOLOGY] ────────────────────────────────────────┐
      │ POOL NODES: 720 | ACTIVE SEEDS: 694 | DEGRADED: 22 | QUARANTINED: 4 | MEAN HEALTH: 98.6%   │
      ├───────────────────────┬────────────┬─────────────┬──────────────┬─────────────┬────────────┤
      │ SEED MAILBOX ID       │ PROVIDER   │ REGION      │ INBOX %      │ RTT (ms)    │ STATUS     │
      ├───────────────────────┼────────────┼─────────────┼──────────────┼─────────────┼────────────┤
      │ seed-gsuite-us-101    │ Google     │ us-east-1   │ 99.8%        │ 17.8ms      │ [●] STABLE │
      │ seed-m365-eu-042      │ Microsoft  │ eu-west-1   │ 98.1%        │ 22.4ms      │ [●] STABLE │
      │ seed-yahoo-us-019     │ Yahoo      │ us-central1 │ 92.4%        │ 41.5ms      │ [▲] WARN   │
      │ seed-icloud-ap-011    │ Apple      │ ap-northeast│ 99.1%        │ 34.2ms      │ [●] STABLE │
      ├───────────────────────┴────────────┴─────────────┴──────────────┴─────────────┴────────────┤
      │ [SELECTED: seed-yahoo-us-019] Health: 92.4% | Target MX: mx-core.mail.yahoo.com             │
      │ Keys: [↑/↓: Navigate] [enter: Drill Down] [f: Toggle Alert Filter] [r: Force Probe] [q: Esc]│
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-matrix`
  - Flags:
    - `--provider <string>`: Filter seed pool view by ESP provider (`google`, `microsoft`, `yahoo`, `apple`, `all`; default: `all`).
    - `--refresh-rate <duration>`: Interval for real-time matrix refresh and probe streaming (default: `1s`).
    - `--is-interactive`: Enable interactive keyboard navigation cursor and detail drill-down (default: `true`).
    - `--has-alerts-only`: Filter matrix view to display only degraded or critical seed mailboxes (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedMailboxHealthMatrixRecord struct {
        MailboxID                string         `json:"mailbox_id"`
        ProviderName             string         `json:"provider_name"`
        ClusterRegion            string         `json:"cluster_region"`
        HealthScorePercent       float64        `json:"health_score_percent"`
        InboxPlacementRate       float64        `json:"inbox_placement_rate"`
        SpamFolderRate           float64        `json:"spam_folder_rate"`
        LatencyRttMs             float64        `json:"latency_rtt_ms"`
        IsHealthy                bool           `json:"is_healthy"`
        IsQuarantined            bool           `json:"is_quarantined"`
        HasAnomalousSpike        bool           `json:"has_anomalous_spike"`
        HasVerifiedAuth          bool           `json:"has_verified_auth"`
    }

    type SeedMailboxHealthMatrixModel struct {
        Records        []SeedMailboxHealthMatrixRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedMailboxHealthMatrixRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
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
  - Bound to `gitmap email warmup seed-matrix`. Directly queries `repodb/pipeline.db` using the shared read replica and delivers real-time Bubble Tea frame updates via tea.Program loop.

---

### Improvement 1432: Seed Mailbox IMAP Auto-Drain & Trash Quarantine Inspector
- **ID:** `SEED-IMAP-DRAIN-INSPECT-1432`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide high-speed IMAP queue flushing, spam-to-inbox rescue verification, and quarantine telemetry across active warmup seeds.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live IMAP drain visualizer displaying rescued message counters, false-positive flag clears, and batch deletion progress.
    - ANSI Progress Bar: `[====================>-----] 80% Drained (412/515 msgs) | Rescued: 28`
  - **Zero-CDN Web View:**
    - Inline HTML/CSS status table with pure embedded styles showing mailbox IMAP connection status, unread message backlog, and auto-rescue actions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-drain`
  - Flags:
    - `--mailbox-id <string>`: Specific mailbox identifier or `all` for cluster-wide sweep.
    - `--batch-size <int>`: Number of messages processed per IMAP chunk (default: `50`).
    - `--is-dry-run`: Evaluate inbox rescue rules without deleting or moving messages (default: `false`).
    - `--has-auto-rescue`: Automatically move spam warmup messages back to primary inbox (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedImapDrainRecord struct {
        DrainSessionID     string `json:"drain_session_id"`
        MailboxID          string `json:"mailbox_id"`
        ProcessedMessages  int    `json:"processed_messages"`
        RescuedSpamCount   int    `json:"rescued_spam_count"`
        PurgedTrashCount   int    `json:"purged_trash_count"`
        DurationMs         int64  `json:"duration_ms"`
        IsCompleted        bool   `json:"is_completed"`
        HasErrors          bool   `json:"has_errors"`
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
  - Integrates with `gitmap email warmup seed-drain` and triggers asynchronous worker pools via lemon daemon background tasks.

---

### Improvement 1433: Seed Reputation Trend Graph & Cohort Degradation Detector
- **ID:** `SEED-REPUTATION-TREND-1433`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Compute moving averages and trend derivatives for deliverability cohorts to isolate seed mailboxes showing early signs of ISP throttling.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Braille-character trend graph (`⡀⠄⠂⠁⠁⠄⡀`) displaying deliverability curves over the last 14 days.
    - Alert card highlighting cohorts dropping below the 95% inbox placement SLA.
  - **Zero-CDN Web View:**
    - SVG sparklines rendered purely using server-side XML string generators with zero client JS libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-trends`
  - Flags:
    - `--timeframe <string>`: Analysis window (`24h`, `7d`, `14d`, `30d`; default: `7d`).
    - `--threshold <float>`: SLA drop threshold percentage (default: `95.0`).
    - `--is-verbose`: Include granular per-hour probe measurements (default: `false`).
    - `--has-auto-demote`: Automatically remove degrading seeds from active warmup pools (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedReputationTrendRecord struct {
        CohortID           string  `json:"cohort_id"`
        Provider           string  `json:"provider"`
        CurrentSlaPercent  float64 `json:"current_sla_percent"`
        SlopeDerivative    float64 `json:"slope_derivative"`
        ObservationWindow  string  `json:"observation_window"`
        IsDegrading        bool    `json:"is_degrading"`
        HasBreachedSla     bool    `json:"has_breached_sla"`
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
  - Directly invoked by `gitmap email warmup seed-trends` and feeds seed telemetry into the multi-domain concurrency governor.

---

### Improvement 1434: Seed Mailbox Anycast Routing & Latency Heatmap Inspector
- **ID:** `SEED-ANYCAST-LATENCY-1434`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Measure and display regional round-trip times (RTT) from sender nodes to seed mailboxes across worldwide data centers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - ANSI heatmap terminal matrix color-coded by ping times: `<20ms` Green, `20-60ms` Yellow, `>60ms` Red.
    - Regional routing summary: `US-East: 18ms | US-West: 34ms | EU-Central: 22ms | AP-South: 68ms`.
  - **Zero-CDN Web View:**
    - Pure CSS grid layout with semantic badge tokens and zero third-party font/asset loads.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-latency`
  - Flags:
    - `--region <string>`: Target edge region (`us-east`, `us-west`, `eu-west`, `ap-east`, `all`; default: `all`).
    - `--ping-count <int>`: Probes per region (default: `5`).
    - `--is-live`: Continuously stream latency samples (default: `false`).
    - `--has-jitter-check`: Calculate standard deviation across successive probe packets (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedLatencyProbeRecord struct {
        ProbeID       string  `json:"probe_id"`
        SourceRegion  string  `json:"source_region"`
        TargetRegion  string  `json:"target_region"`
        AvgLatencyMs  float64 `json:"avg_latency_ms"`
        JitterMs      float64 `json:"jitter_ms"`
        PacketLossPct float64 `json:"packet_loss_pct"`
        IsOptimal     bool    `json:"is_optimal"`
        HasHighJitter bool    `json:"has_high_jitter"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedLatencyProbeRecord (
        probe_id TEXT PRIMARY KEY,
        source_region TEXT NOT NULL,
        target_region TEXT NOT NULL,
        avg_latency_ms REAL NOT NULL DEFAULT 0.0,
        jitter_ms REAL NOT NULL DEFAULT 0.0,
        packet_loss_pct REAL NOT NULL DEFAULT 0.0,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_high_jitter INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup seed-latency` and stores probe logs in SQLite database partitions.

---

### Improvement 1435: Seed Authentication Audit & SPF/DKIM/DMARC Pass-Rate HUD
- **ID:** `SEED-AUTH-AUDIT-HUD-1435`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Continuous verification of cryptographic sender authentication headers received across seed mailboxes to verify 100% alignment.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Lipgloss table showing SPF alignment, DKIM selector validity, and DMARC enforcement modes.
    - Status badges: `[SPF: PASS] [DKIM: PASS (2048-bit)] [DMARC: PASS (p=reject)]`.
  - **Zero-CDN Web View:**
    - Standalone security scorecard with SVG shield icons and pure inline styling.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-auth`
  - Flags:
    - `--domain <string>`: Sending domain to audit against seed incoming logs.
    - `--strict-alignment`: Require both SPF and DKIM identifier alignment (default: `true`).
    - `--is-enforced`: Verify DMARC policy is set to quarantine or reject (default: `true`).
    - `--has-detailed-headers`: Dump raw Authentication-Results headers (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedAuthAuditRecord struct {
        AuditID         string `json:"audit_id"`
        DomainName      string `json:"domain_name"`
        SpfResult       string `json:"spf_result"`
        DkimResult      string `json:"dkim_result"`
        DmarcResult     string `json:"dmarc_result"`
        IsFullyAligned  bool   `json:"is_fully_aligned"`
        HasValidDkimKey bool   `json:"has_valid_dkim_key"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedAuthAuditRecord (
        audit_id TEXT PRIMARY KEY,
        domain_name TEXT NOT NULL,
        spf_result TEXT NOT NULL,
        dkim_result TEXT NOT NULL,
        dmarc_result TEXT NOT NULL,
        is_fully_aligned INTEGER NOT NULL DEFAULT 1,
        has_valid_dkim_key INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - CLI command `gitmap email warmup seed-auth` validates DKIM/SPF headers and provides pass/fail verification for deliverability audits.

---

### Improvement 1436: Seed Pool Cohort Entropy & Threading Pattern Simulator
- **ID:** `SEED-ENTROPY-THREADING-1436`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Simulate natural conversational thread depth and human reading cadences across warmup seeds to prevent synthetic pattern detection.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live thread simulator display showing conversation branches, reply delay variance (e.g., 4m to 38m), and entropy scores.
    - ANSI gauge: `Entropy Factor: 0.89 [HIGH VARIANCE / NATURAL]`.
  - **Zero-CDN Web View:**
    - Hierarchical DOM tree rendering conversation reply threads with CSS indentation and zero external JS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-entropy`
  - Flags:
    - `--depth <int>`: Maximum conversation reply depth (default: `4`).
    - `--min-delay <duration>`: Minimum interval before simulated seed reply (default: `5m`).
    - `--is-randomized`: Apply Poisson distribution to reply timing (default: `true`).
    - `--has-body-variation`: Inject AI-generated linguistic variations into replies (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedEntropyThreadRecord struct {
        ThreadID          string  `json:"thread_id"`
        SeedSenderID      string  `json:"seed_sender_id"`
        SeedReceiverID    string  `json:"seed_receiver_id"`
        ConversationDepth int     `json:"conversation_depth"`
        EntropyScore      float64 `json:"entropy_score"`
        IsNaturalPattern  bool    `json:"is_natural_pattern"`
        HasReplyPending   bool    `json:"has_reply_pending"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedEntropyThreadRecord (
        thread_id TEXT PRIMARY KEY,
        seed_sender_id TEXT NOT NULL,
        seed_receiver_id TEXT NOT NULL,
        conversation_depth INTEGER NOT NULL DEFAULT 1,
        entropy_score REAL NOT NULL DEFAULT 0.5,
        is_natural_pattern INTEGER NOT NULL DEFAULT 1,
        has_reply_pending INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed by `gitmap email warmup seed-entropy`, updating thread state in SQLite to coordinate seed response engines.

---

### Improvement 1437: Seed Spam-Trap Honeypot Quarantine & Disqualification Matrix
- **ID:** `SEED-HONEYPOT-QUARANTINE-1437`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Continuously cross-reference warmup seed pools against known spam trap networks, MX honeytokens, and dormant address recycling databases.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Alert HUD in bold red Lipgloss styling highlighting quarantined addresses with immediate disqualify/purge action prompts.
    - Status banner: `[!] HONEYPOT DETECTED: seed-legacy-77@domain.com | ACTION: QUARANTINED`.
  - **Zero-CDN Web View:**
    - Security quarantine table with instant isolation buttons and SHA-256 evidence hashes.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-quarantine`
  - Flags:
    - `--action <string>`: Operation to perform (`audit`, `isolate`, `purge`; default: `audit`).
    - `--risk-level <string>`: Minimum severity threshold (`low`, `medium`, `high`; default: `medium`).
    - `--is-auto-quarantine`: Automatically isolate flagged addresses without manual prompt (default: `true`).
    - `--has-audit-log`: Generate signed JSON audit trail in `.ai-memory/` (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedQuarantineRecord struct {
        RecordID         string `json:"record_id"`
        SeedMailboxID    string `json:"seed_mailbox_id"`
        DetectionVector  string `json:"detection_vector"`
        ThreatLevel      string `json:"threat_level"`
        IsQuarantined    bool   `json:"is_quarantined"`
        HasPurgeApproved bool   `json:"has_purge_approved"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedQuarantineRecord (
        record_id TEXT PRIMARY KEY,
        seed_mailbox_id TEXT NOT NULL,
        detection_vector TEXT NOT NULL,
        threat_level TEXT NOT NULL,
        is_quarantined INTEGER NOT NULL DEFAULT 1,
        has_purge_approved INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed via `gitmap email warmup seed-quarantine` to safeguard sender domains from honeypot exposure.

---

### Improvement 1438: Seed Inbox Folder Placement Live Ticker & Provider SLA Matrix
- **ID:** `SEED-PLACEMENT-TICKER-1438`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Real-time terminal ticker displaying incoming warmup delivery destinations (Inbox, Promotions, Spam, Trash) categorized by ESP.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density scrolling ticker with provider breakdowns:
      `[GMAIL] 98.2% Inbox | 1.8% Promo | [M365] 100% Focused | [YAHOO] 95.0% Inbox | 5.0% Spam`.
  - **Zero-CDN Web View:**
    - Live-updating HTML table using native `text/event-stream` SSE with embedded CSS flex containers.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-ticker`
  - Flags:
    - `--window <duration>`: Rolling calculation window (default: `1h`).
    - `--min-samples <int>`: Minimum messages required for statistically valid rate (default: `20`).
    - `--is-streaming`: Continuously stream live arrival ticker events (default: `true`).
    - `--has-color-coding`: Apply dynamic ANSI terminal colors based on threshold (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedPlacementTickerRecord struct {
        TickerID          string  `json:"ticker_id"`
        Provider          string  `json:"provider"`
        InboxPlacementPct float64 `json:"inbox_placement_pct"`
        PromoPlacementPct float64 `json:"promo_placement_pct"`
        SpamPlacementPct  float64 `json:"spam_spam_pct"`
        SampleCount       int     `json:"sample_count"`
        IsMeetingSla      bool    `json:"is_meeting_sla"`
        HasRecentSpamDrop bool    `json:"has_recent_spam_drop"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedPlacementTickerRecord (
        ticker_id TEXT PRIMARY KEY,
        provider TEXT NOT NULL,
        inbox_placement_pct REAL NOT NULL DEFAULT 100.0,
        promo_placement_pct REAL NOT NULL DEFAULT 0.0,
        spam_placement_pct REAL NOT NULL DEFAULT 0.0,
        sample_count INTEGER NOT NULL DEFAULT 0,
        is_meeting_sla INTEGER NOT NULL DEFAULT 1,
        has_recent_spam_drop INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Linked to `gitmap email warmup seed-ticker` for live terminal monitoring during outbound campaigns.

---

### Improvement 1439: Seed Mailbox IMAP/SMTP Credential Health & OAuth Token Expiry Gauge
- **ID:** `SEED-CREDENTIAL-GAUGE-1439`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Inspect seed pool authentication tokens, monitor OAuth refresh lifecycles, and alert ahead of credential expiration or account disconnects.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Token health gauge displaying days/hours remaining until OAuth refresh token expiration or re-auth requirement.
    - Status indicators: `[seed-04: OAuth Valid (28d remaining)]` vs `[seed-19: Token Expiring in 2h!]`.
  - **Zero-CDN Web View:**
    - Standalone tabular status list with time-remaining countdown badges and re-authentication deep links.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-creds`
  - Flags:
    - `--expiring-within <duration>`: Filter tokens expiring within threshold (default: `48h`).
    - `--auto-refresh`: Attempt automatic OAuth token renewal if refresh token is present (default: `true`).
    - `--is-authenticated`: Filter by active authentication status (default: `true`).
    - `--has-expired-only`: Show only disconnected or failed accounts (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedCredentialGaugeRecord struct {
        AccountID         string `json:"account_id"`
        MailboxAddress    string `json:"mailbox_address"`
        AuthType          string `json:"auth_type"`
        ExpiresAtUnix     int64  `json:"expires_at_unix"`
        SecondsRemaining  int64  `json:"seconds_remaining"`
        IsValid           bool   `json:"is_valid"`
        HasAutoRenewSucceeded bool `json:"has_auto_renew_succeeded"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedCredentialGaugeRecord (
        account_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        auth_type TEXT NOT NULL,
        expires_at_unix INTEGER NOT NULL DEFAULT 0,
        seconds_remaining INTEGER NOT NULL DEFAULT 0,
        is_valid INTEGER NOT NULL DEFAULT 1,
        has_auto_renew_succeeded INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated with `gitmap email warmup seed-creds` to ensure continuous seed availability and unattended operation.

---

### Improvement 1440: Comprehensive Seed Diagnostics Self-Test & Diagnostic Export Bundle
- **ID:** `SEED-DIAGNOSTICS-EXPORT-1440`
- **Cluster:** 144 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Run comprehensive end-to-end self-tests across the entire seed cluster (IMAP, SMTP, latency, placement, DNS) and export a signed diagnostic report bundle.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-stage test progress spinner: `[✓] DNS Resolution -> [✓] IMAP Handshake -> [✓] SMTP Submission -> [✓] Auth Alignment`.
    - Final report card with green confirmation border: `ALL 12 DIAGNOSTIC CHECKS PASSED (100% Green)`.
  - **Zero-CDN Web View:**
    - Self-contained HTML report with embedded SVG icons, CSS print styles, and download button for JSON diagnostics bundle.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-selftest`
  - Flags:
    - `--output-format <string>`: Export format (`summary`, `json`, `bundle`; default: `summary`).
    - `--dest-file <string>`: Output destination relative to repo root (e.g. `02-spec/21-app/diagnostics.json`).
    - `--is-deep-probe`: Execute full end-to-end send/receive loop test (default: `true`).
    - `--has-summary-only`: Output one-line exit summary for CI/CD assertions (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedDiagnosticBundleRecord struct {
        BundleID          string  `json:"bundle_id"`
        GeneratedAtUnix   int64   `json:"generated_at_unix"`
        TotalSeedsAudited int     `json:"total_seeds_audited"`
        PassedSeedCount   int     `json:"passed_seed_count"`
        OverallHealthPct  float64 `json:"overall_health_pct"`
        IsPassed          bool    `json:"is_passed"`
        HasAnomalies      bool    `json:"has_anomalies"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedDiagnosticBundleRecord (
        bundle_id TEXT PRIMARY KEY,
        generated_at_unix INTEGER NOT NULL DEFAULT 0,
        total_seeds_audited INTEGER NOT NULL DEFAULT 0,
        passed_seed_count INTEGER NOT NULL DEFAULT 0,
        overall_health_pct REAL NOT NULL DEFAULT 100.0,
        is_passed INTEGER NOT NULL DEFAULT 1,
        has_anomalies INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked via `gitmap email warmup seed-selftest`, creating verified health snapshots in SQLite.

---

# Cluster 145: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 1441–1450)

---

### Improvement 1441: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor Controller
- **ID:** `GOV-BACKOFF-INVERSION-1441`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Dynamically govern concurrency limits per target MX domain, inverting backoff delays upon transient 421/450 SMTP error detection while maintaining overall cluster throughput.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dynamic governor dashboard with concurrency gauges per domain and backoff inversion multipliers.
    - Layout:
      ```text
      ┌─ [MULTI-DOMAIN CONCURRENCY GOVERNOR & BACKOFF INVERTER] ─────────────────────────────────┐
      │ GLOBAL SLOTS: 250 | OCCUPIED: 142 | THROTTLED TARGETS: 2 | ACTIVE BACKOFF MULTIPLIER: 1.4x│
      ├───────────────────────┬──────────┬──────────┬──────────────┬──────────────┬──────────────┤
      │ DOMAIN / MX GROUP     │ LIMIT    │ IN-FLIGHT│ BACKOFF DELAY│ 4XX ERROR %  │ GOVERNOR ST  │
      ├───────────────────────┼──────────┼──────────┼──────────────┼──────────────┼──────────────┤
      │ google.com            │ 50       │ 48       │ 100ms        │ 0.02%        │ [●] OPTIMAL  │
      │ outlook.com           │ 35       │ 32       │ 150ms        │ 0.15%        │ [●] OPTIMAL  │
      │ yahoo.com             │ 15       │ 8        │ 1250ms       │ 4.80%        │ [▲] BACKOFF  │
      │ orange.fr             │ 10       │ 3        │ 3500ms       │ 8.20%        │ [■] INVERTED │
      ├───────────────────────┴──────────┴──────────┴──────────────┴──────────────┴──────────────┤
      │ [ACTIVE POLICY: ADAPTIVE-AIMD] Additive Increase, Multiplicative Decrease active            │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Pure CSS status bars with real-time SSE stream updates showing active connection slots vs hard ceilings.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup governor-status`
  - Flags:
    - `--target-domain <string>`: Filter governor state by target domain (default: `all`).
    - `--aimd-factor <float>`: Multiplicative decrease factor on 4xx reception (default: `0.5`).
    - `--is-adaptive`: Enable automatic AIMD rate adaptation (default: `true`).
    - `--has-strict-ceilings`: Enforce hard ceiling limits regardless of queue backlog (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainConcurrencyGovernorRecord struct {
        TargetMxGroup      string  `json:"target_mx_group"`
        CurrentConcurrency int     `json:"current_concurrency"`
        MaxAllowedSlots    int     `json:"max_allowed_slots"`
        BackoffDelayMs     int64   `json:"backoff_delay_ms"`
        ErrorRatePercent   float64 `json:"error_rate_percent"`
        IsThrottled        bool    `json:"is_throttled"`
        HasInvertedBackoff bool    `json:"has_inverted_backoff"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainConcurrencyGovernorRecord (
        target_mx_group TEXT PRIMARY KEY,
        current_concurrency INTEGER NOT NULL DEFAULT 1,
        max_allowed_slots INTEGER NOT NULL DEFAULT 10,
        backoff_delay_ms INTEGER NOT NULL DEFAULT 100,
        error_rate_percent REAL NOT NULL DEFAULT 0.0,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        has_inverted_backoff INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup governor-status` and coordinates with lemon background outbound workers.

---

### Improvement 1442: AIMD Additive Increase / Multiplicative Decrease Algorithm Telemetry HUD
- **ID:** `GOV-AIMD-TELEMETRY-1442`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Visualize real-time AIMD window size adjustments per sending IP and MX domain pairing to audit rate convergence.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Sparkline visualizer of congestion window (cwnd) expansion and contraction phases.
    - ANSI notation: `[AIMD] Phase: ADDITIVE_INCREASE | Window: 18 -> 19 msgs/sec (+1.0)`.
  - **Zero-CDN Web View:**
    - Inline SVG window graph generated directly from server memory without client dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup aimd-telemetry`
  - Flags:
    - `--ip-address <string>`: Filter by sending interface IP (default: `all`).
    - `--window-history <int>`: Number of AIMD cycles to display (default: `30`).
    - `--is-streaming`: Stream live window changes to stdout (default: `false`).
    - `--has-loss-indicators`: Highlight packet drops or SMTP reject markers (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type AimdTelemetryRecord struct {
        CycleID           string  `json:"cycle_id"`
        SendingIP         string  `json:"sending_ip"`
        TargetMX          string  `json:"target_mx"`
        WindowSize        int     `json:"window_size"`
        PacingIntervalMs  int64   `json:"pacing_interval_ms"`
        IsAdditivePhase   bool    `json:"is_additive_phase"`
        HasDetectedLoss   bool    `json:"has_detected_loss"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AimdTelemetryRecord (
        cycle_id TEXT PRIMARY KEY,
        sending_ip TEXT NOT NULL,
        target_mx TEXT NOT NULL,
        window_size INTEGER NOT NULL DEFAULT 1,
        pacing_interval_ms INTEGER NOT NULL DEFAULT 1000,
        is_additive_phase INTEGER NOT NULL DEFAULT 1,
        has_detected_loss INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `gitmap email warmup aimd-telemetry` to audit congestion avoidance behavior.

---

### Improvement 1443: Leaky Bucket Outbound Dispatch Shaper & Inter-Burst Jitter Governor
- **ID:** `GOV-LEAKY-BUCKET-SHAPER-1443`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Apply leaky bucket token pacing with microsecond random jitter to eliminate predictable burst signatures across destination mail servers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Pacing meter showing bucket token fill level, leak rate, and injected inter-burst jitter milliseconds.
    - ANSI display: `Token Bucket: [████████░░] 80% Capacity | Leak Rate: 4.2 msg/s | Jitter: ±184ms`.
  - **Zero-CDN Web View:**
    - HTML5 meter elements styled with CSS custom properties and zero external stylesheets.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup shaper-config`
  - Flags:
    - `--leak-rate <float>`: Target tokens leaked per second (default: `5.0`).
    - `--max-burst <int>`: Maximum bucket capacity before queue rejection (default: `20`).
    - `--is-jitter-enabled`: Inject random Gaussian jitter between dispatches (default: `true`).
    - `--has-priority-bypass`: Allow transactional warmup bypass for high-priority seed health pings (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LeakyBucketShaperRecord struct {
        ShaperID         string  `json:"shaper_id"`
        DomainGroup      string  `json:"domain_group"`
        TokenCapacity    int     `json:"token_capacity"`
        CurrentTokens    float64 `json:"current_tokens"`
        LeakRatePerSec   float64 `json:"leak_rate_per_sec"`
        IsPacingActive   bool    `json:"is_pacing_active"`
        HasJitterApplied bool    `json:"has_jitter_applied"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS LeakyBucketShaperRecord (
        shaper_id TEXT PRIMARY KEY,
        domain_group TEXT NOT NULL,
        token_capacity INTEGER NOT NULL DEFAULT 10,
        current_tokens REAL NOT NULL DEFAULT 10.0,
        leak_rate_per_sec REAL NOT NULL DEFAULT 2.0,
        is_pacing_active INTEGER NOT NULL DEFAULT 1,
        has_jitter_applied INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup shaper-config` to tune local dispatch queues without daemon restart.

---

### Improvement 1444: ISP Greylist Trap Classifier & Automatic Deferral Circuit Breaker
- **ID:** `GOV-GREYLIST-CIRCUIT-1444`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Distinguish greylisting temporary deferrals (`451 4.7.1 Greylisted`) from reputation-based rate limits and trip protective circuit breakers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Circuit breaker status matrix: `[CLOSED: NORMAL]` -> `[OPEN: GREYLIST PAUSE (15m remaining)]` -> `[HALF-OPEN: PROBING]`.
    - Detailed regex breakdown of matched SMTP 4xx response codes.
  - **Zero-CDN Web View:**
    - Diagnostic card with status badges and manual breaker reset toggle buttons.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup greylist-breaker`
  - Flags:
    - `--circuit-state <string>`: Filter by state (`closed`, `open`, `half-open`, `all`; default: `all`).
    - `--reset-domain <string>`: Force reset circuit breaker for target domain.
    - `--is-auto-tripped`: Filter circuits tripped by automatic error thresholds (default: `true`).
    - `--has-cooldown-active`: Filter domains currently in backoff cooldown (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GreylistCircuitRecord struct {
        CircuitID        string `json:"circuit_id"`
        TargetDomain     string `json:"target_domain"`
        State            string `json:"state"`
        FailureCount     int    `json:"failure_count"`
        CooldownUntilUnix int64 `json:"cooldown_until_unix"`
        IsClosed         bool   `json:"is_closed"`
        HasTrippedRecent bool   `json:"has_tripped_recent"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GreylistCircuitRecord (
        circuit_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        state TEXT NOT NULL DEFAULT 'closed',
        failure_count INTEGER NOT NULL DEFAULT 0,
        cooldown_until_unix INTEGER NOT NULL DEFAULT 0,
        is_closed INTEGER NOT NULL DEFAULT 1,
        has_tripped_recent INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed by `gitmap email warmup greylist-breaker`, persisting breaker states across process reboots.

---

### Improvement 1445: Multi-Tenant Pool Concurrency Fair-Share Arbiter
- **ID:** `GOV-FAIR-SHARE-ARBITER-1445`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Prevent aggressive warmup sender campaigns from starving shared network bandwidth or target MX domain quotas using weighted fair queuing.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Fair-share allocation table displaying tenant weight, consumed slots, and queue wait times.
    - Lipgloss table: `Tenant-A (wt: 50%): 25/50 slots | Tenant-B (wt: 50%): 25/50 slots`.
  - **Zero-CDN Web View:**
    - Pure CSS multi-column proportional bar chart comparing allocated vs utilized quotas.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup fairshare-status`
  - Flags:
    - `--tenant-id <string>`: Filter quota display by specific tenant identifier.
    - `--rebalance`: Recalculate weights based on active pending queue depth.
    - `--is-starvation-prevented`: Verify zero-starvation guarantee across low-weight tenants (default: `true`).
    - `--has-active-work`: Show only tenants with messages currently in-flight (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FairShareTenantRecord struct {
        TenantID         string  `json:"tenant_id"`
        WeightRatio      float64 `json:"weight_ratio"`
        AllocatedSlots   int     `json:"allocated_slots"`
        ActiveInFlight   int     `json:"active_in_flight"`
        QueueWaitMs      int64   `json:"queue_wait_ms"`
        IsSlaCompliant   bool    `json:"is_sla_compliant"`
        HasStarvationRisk bool   `json:"has_starvation_risk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FairShareTenantRecord (
        tenant_id TEXT PRIMARY KEY,
        weight_ratio REAL NOT NULL DEFAULT 1.0,
        allocated_slots INTEGER NOT NULL DEFAULT 10,
        active_in_flight INTEGER NOT NULL DEFAULT 0,
        queue_wait_ms INTEGER NOT NULL DEFAULT 0,
        is_sla_compliant INTEGER NOT NULL DEFAULT 1,
        has_starvation_risk INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked via `gitmap email warmup fairshare-status` to inspect multi-tenant scheduler arbitration.

---

### Improvement 1446: Dynamic MX Pipelining Handshake Latency Optimization Engine
- **ID:** `GOV-MX-PIPELINING-1446`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Optimize SMTP pipelining negotiation and TLS session resumption parameters per remote MX host to minimize connection setup RTTs.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Pipelining efficiency scorecard: `EHLO -> STARTTLS (0-RTT TLS 1.3) -> PIPELINING (MAIL FROM + RCPT TO grouped)`.
    - Handshake latency saving badge: `Saved: 48.2ms per delivery (-62% latency)`.
  - **Zero-CDN Web View:**
    - Comparative latency timeline rendered with pure HTML/CSS flexbox tracks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup pipelining-audit`
  - Flags:
    - `--mx-host <string>`: Target mail exchange host to test for RFC 2920 pipelining support.
    - `--tls-version <string>`: Minimum TLS protocol version (`1.2`, `1.3`; default: `1.2`).
    - `--is-pipelining-supported`: Filter hosts supporting pipelined command transmission (default: `true`).
    - `--has-tls-resumption`: Verify TLS ticket caching for 0-RTT re-handshakes (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MxPipeliningAuditRecord struct {
        MxHost             string  `json:"mx_host"`
        SupportsPipelining bool    `json:"supports_pipelining"`
        SupportsStartTLS   bool    `json:"supports_starttls"`
        HandshakeDurationMs float64 `json:"handshake_duration_ms"`
        OptimizedSavingMs  float64 `json:"optimized_saving_ms"`
        IsZeroRttCapable   bool    `json:"is_zero_rtt_capable"`
        HasValidCert       bool    `json:"has_valid_cert"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxPipeliningAuditRecord (
        mx_host TEXT PRIMARY KEY,
        supports_pipelining INTEGER NOT NULL DEFAULT 1,
        supports_starttls INTEGER NOT NULL DEFAULT 1,
        handshake_duration_ms REAL NOT NULL DEFAULT 0.0,
        optimized_saving_ms REAL NOT NULL DEFAULT 0.0,
        is_zero_rtt_capable INTEGER NOT NULL DEFAULT 1,
        has_valid_cert INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email warmup pipelining-audit` to log remote MX capability flags into local SQLite cache.

---

### Improvement 1447: Outbound IP Warmup Phase Ramp-Up Governor & Ceiling Automator
- **ID:** `GOV-RAMPUP-CEILING-1447`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Automatically enforce day-by-day sending ceilings (Day 1: 50, Day 2: 100, Day 3: 200...) based on verified historical reputation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Step-ladder ramp visualizer showing current day phase, quota utilized, and next tier unlock timer.
    - ANSI Progress: `Day 4/14: [██████████░░░░░] 420/600 msgs (70%) | Next Tier: +300 msgs in 6h 12m`.
  - **Zero-CDN Web View:**
    - Interactive stepped progress bar with CSS transitions and SVG milestone medals.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup rampup-schedule`
  - Flags:
    - `--sending-ip <string>`: Target IP address under warmup schedule.
    - `--current-day <int>`: Force override warmup day step (requires confirmation).
    - `--is-auto-increment`: Automatically advance to next tier if error rate < 1% (default: `true`).
    - `--has-strict-ceiling`: Block attempts to exceed current day quota (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type WarmupRampupScheduleRecord struct {
        ScheduleID        string  `json:"schedule_id"`
        SendingIP         string  `json:"sending_ip"`
        CurrentDayNumber  int     `json:"current_day_number"`
        DailyLimit        int     `json:"daily_limit"`
        SentToday         int     `json:"sent_today"`
        ErrorRateTodayPct float64 `json:"error_rate_today_pct"`
        IsQuotaExhausted  bool    `json:"is_quota_exhausted"`
        HasAdvancePermit  bool    `json:"has_advance_permit"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupRampupScheduleRecord (
        schedule_id TEXT PRIMARY KEY,
        sending_ip TEXT NOT NULL,
        current_day_number INTEGER NOT NULL DEFAULT 1,
        daily_limit INTEGER NOT NULL DEFAULT 50,
        sent_today INTEGER NOT NULL DEFAULT 0,
        error_rate_today_pct REAL NOT NULL DEFAULT 0.0,
        is_quota_exhausted INTEGER NOT NULL DEFAULT 0,
        has_advance_permit INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `gitmap email warmup rampup-schedule`, enforcing daily outbound ceilings at dispatch time.

---

### Improvement 1448: Multi-IP Round-Robin Pool Affinity & Sticky Domain Binding
- **ID:** `GOV-IP-POOL-AFFINITY-1448`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Bind destination domains to dedicated IP pools while applying consistent hashing to prevent cross-contamination of IP reputations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Topology map showing source IP addresses bound to destination domain tiers (Tier 1: Gmail/M365, Tier 2: Yahoo/AOL).
    - Affinity table: `IP 192.0.2.14 -> outlook.com (Sticky Session: 94h remaining)`.
  - **Zero-CDN Web View:**
    - Matrix diagram rendered in HTML table format with CSS colored indicators.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ip-affinity`
  - Flags:
    - `--pool-name <string>`: Name of IP allocation pool (default: `production-warmup`).
    - `--domain <string>`: Filter binding rules by target domain.
    - `--is-sticky`: Maintain persistent IP binding for active recipient conversations (default: `true`).
    - `--has-auto-fallback`: Allow secondary IP usage if primary IP is throttled (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type IpPoolAffinityRecord struct {
        BindingID         string `json:"binding_id"`
        PoolName          string `json:"pool_name"`
        SendingIP         string `json:"sending_ip"`
        DestinationDomain string `json:"destination_domain"`
        AssignedTier      string `json:"assigned_tier"`
        IsStickyActive    bool   `json:"is_sticky_active"`
        HasHealthyRating  bool   `json:"has_healthy_rating"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS IpPoolAffinityRecord (
        binding_id TEXT PRIMARY KEY,
        pool_name TEXT NOT NULL,
        sending_ip TEXT NOT NULL,
        destination_domain TEXT NOT NULL,
        assigned_tier TEXT NOT NULL,
        is_sticky_active INTEGER NOT NULL DEFAULT 1,
        has_healthy_rating INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed via `gitmap email warmup ip-affinity`, querying routing tables stored in `repodb/pipeline.db`.

---

### Improvement 1449: Transient SMTP Error Cluster Classification & Feedback Loop
- **ID:** `GOV-ERROR-CLASSIFICATION-1449`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Parse unstructured SMTP bounce and rejection messages using regex and heuristic tokenizers to classify transient errors and trigger targeted backoff adjustments.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Error classification stream displaying parsed codes: `421 4.7.0 [RATE_LIMIT]`, `450 4.2.1 [MAILBOX_BUSY]`, `550 5.1.1 [USER_UNKNOWN]`.
    - Real-time resolution action badge: `[ACTION: BACKOFF_INVERT (delay: 2400ms)]`.
  - **Zero-CDN Web View:**
    - Tabular log viewer with regex highlight tags and search filters.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup error-classifier`
  - Flags:
    - `--log-source <string>`: Path or stream of SMTP transaction logs.
    - `--min-severity <string>`: Severity filter (`info`, `warning`, `critical`; default: `warning`).
    - `--is-transient-only`: Show only temporary 4xx failures triggering backoff (default: `true`).
    - `--has-suggested-action`: Include automated remediation suggestion (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SmtpErrorClassificationRecord struct {
        ClassificationID string `json:"classification_id"`
        RawResponseText  string `json:"raw_response_text"`
        SmtpCode         int    `json:"smtp_code"`
        EnhancedCode     string `json:"enhanced_code"`
        CategoryTag      string `json:"category_tag"`
        SuggestedAction  string `json:"suggested_action"`
        IsTransient      bool   `json:"is_transient"`
        HasTriggeredRule bool   `json:"has_triggered_rule"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpErrorClassificationRecord (
        classification_id TEXT PRIMARY KEY,
        raw_response_text TEXT NOT NULL,
        smtp_code INTEGER NOT NULL,
        enhanced_code TEXT NOT NULL,
        category_tag TEXT NOT NULL,
        suggested_action TEXT NOT NULL,
        is_transient INTEGER NOT NULL DEFAULT 1,
        has_triggered_rule INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - CLI command `gitmap email warmup error-classifier` provides real-time error taxonomy to the governor.

---

### Improvement 1450: Global Governor Configuration Snapshot & Dynamic Re-Arm Engine
- **ID:** `GOV-CONFIG-SNAPSHOT-1450`
- **Cluster:** 145 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Capture, validate, hot-reload, and roll back governor parameters (AIMD rates, leaky bucket sizes, domain ceilings) without dropping in-flight SMTP connections.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Hot-reload status card with diff display comparing active vs pending governor configuration.
    - Confirmation banner: `[✓] GOVERNOR HOT-RELOAD SUCCESSFUL (0 in-flight connections dropped)`.
  - **Zero-CDN Web View:**
    - Clean form-based configuration UI with pure HTML validation attributes and zero external polyfills.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup governor-reload`
  - Flags:
    - `--config-file <string>`: Relative path to new governor JSON configuration.
    - `--validate-only`: Perform schema validation and dry-run comparison without applying.
    - `--is-hot-reload`: Apply configuration changes immediately to running daemons (default: `true`).
    - `--has-rollback-backup`: Save prior configuration snapshot for automatic fallback (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GovernorConfigSnapshotRecord struct {
        SnapshotID        string `json:"snapshot_id"`
        AppliedAtUnix     int64  `json:"applied_at_unix"`
        ConfigChecksumSha string `json:"config_checksum_sha"`
        ActiveTenantCount int    `json:"active_tenant_count"`
        IsActive          bool   `json:"is_active"`
        HasBackupStored   bool   `json:"has_backup_stored"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorConfigSnapshotRecord (
        snapshot_id TEXT PRIMARY KEY,
        applied_at_unix INTEGER NOT NULL DEFAULT 0,
        config_checksum_sha TEXT NOT NULL,
        active_tenant_count INTEGER NOT NULL DEFAULT 0,
        is_active INTEGER NOT NULL DEFAULT 1,
        has_backup_stored INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed via `gitmap email warmup governor-reload`, orchestrating live updates across all active worker instances.

---

# Cluster 146: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 1451–1460)

---

### Improvement 1451: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD
- **ID:** `SEC-HTML-SANDBOX-HUD-1451`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Provide an isolated, zero-CDN HTML/CSS sanitizer HUD that inspects email templates, strips dangerous script/tracking tags, and validates CSS responsiveness against strict Content Security Policy (CSP) rules.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Split-screen TUI inspector showing original HTML AST vs sanitized DOM structure, highlighted with Lipgloss syntax colors.
    - Security overview table:
      ```text
      ┌─ [EMAIL CONTENT SECURITY SANDBOX & SANITIZER HUD] ───────────────────────────────────────┐
      │ AUDITED TEMPLATES: 84 | CLEAN: 79 | STRIPPED TAGS: 12 | CSP VIOLATIONS: 3 | ZERO-CDN: 100%│
      ├───────────────────────┬────────────┬──────────────┬──────────────┬──────────────┬──────────┤
      │ TEMPLATE NAME         │ TOTAL SIZE │ JS/IFRAME    │ EXTERNAL CSS │ TRACKING PIX │ STATUS   │
      ├───────────────────────┼────────────┼──────────────┼──────────────┼──────────────┼──────────┤
      │ onboarding-intro-v2   │ 14.2 KB    │ 0 detected   │ 0 external   │ 1 local SVG  │ [●] PASS │
      │ warm-daily-digest-01  │ 22.8 KB    │ 1 stripped   │ 0 external   │ 0 local      │ [▲] WARN │
      │ promo-newsletter-09   │ 48.6 KB    │ 2 stripped   │ 1 blocked    │ 3 external   │ [■] FAIL │
      ├───────────────────────┴────────────┴──────────────┴──────────────┴──────────────┴──────────┤
      │ [SELECTED: warm-daily-digest-01] Stripped: <script type="text/javascript"> in footer       │
      │ Action Keys: [s: Sanitize & Save] [v: View Diff] [p: Test Preview] [q: Quit]               │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - Hardened iframe container with `sandbox="allow-same-origin"` and `Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'; img-src data:;`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template sanitize`
  - Flags:
    - `--template-path <string>`: Path to HTML template relative to repo root (e.g. `02-spec/templates/welcome.html`).
    - `--strip-external-resources`: Remove all non-data URI images and stylesheets (default: `true`).
    - `--is-interactive-hud`: Open Bubble Tea interactive AST explorer (default: `true`).
    - `--has-strict-csp`: Enforce zero external domain requests (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TemplateSanitizerRecord struct {
        TemplateID       string `json:"template_id"`
        TemplateName     string `json:"template_name"`
        OriginalBytes    int    `json:"original_bytes"`
        SanitizedBytes   int    `json:"sanitized_bytes"`
        StrippedTags     int    `json:"stripped_tags"`
        CspViolations    int    `json:"csp_violations"`
        IsCspCompliant   bool   `json:"is_csp_compliant"`
        HasResponsiveCss bool   `json:"has_responsive_css"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TemplateSanitizerRecord (
        template_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        original_bytes INTEGER NOT NULL DEFAULT 0,
        sanitized_bytes INTEGER NOT NULL DEFAULT 0,
        stripped_tags INTEGER NOT NULL DEFAULT 0,
        csp_violations INTEGER NOT NULL DEFAULT 0,
        is_csp_compliant INTEGER NOT NULL DEFAULT 1,
        has_responsive_css INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked via `gitmap email template sanitize` to enforce zero-CDN and strict CSP before email rendering.

---

### Improvement 1452: Email CSS Media Query & Mobile Responsive Grid Auditing Engine
- **ID:** `SEC-RESPONSIVE-GRID-1452`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit HTML templates against common mobile viewport constraints (320px, 375px, 414px) and verify table/div fluid grid collapse without horizontal overflow.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Viewport test matrix with status indicators: `[320px: PASS] [375px: PASS] [414px: PASS] [768px: PASS]`.
    - Horizontal scroll detection warning: `Warning: Element <table id="hero-grid"> exceeds 320px width (actual: 340px)`.
  - **Zero-CDN Web View:**
    - Dual iframe viewport comparator running locally without external font or styling libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template responsive-check`
  - Flags:
    - `--template-file <string>`: Path to HTML email template relative to workspace.
    - `--max-width <int>`: Standard desktop target width (default: `600`).
    - `--is-fluid-layout`: Verify percentage-based width definitions (default: `true`).
    - `--has-overflow-detected`: Return non-zero status if horizontal overflow occurs (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TemplateResponsiveAuditRecord struct {
        AuditID          string `json:"audit_id"`
        TemplatePath     string `json:"template_path"`
        TargetWidthPx    int    `json:"target_width_px"`
        ComputedWidthPx  int    `json:"computed_width_px"`
        HasOverflow      bool   `json:"has_overflow"`
        IsMobileFriendly bool   `json:"is_mobile_friendly"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TemplateResponsiveAuditRecord (
        audit_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        target_width_px INTEGER NOT NULL DEFAULT 600,
        computed_width_px INTEGER NOT NULL DEFAULT 600,
        has_overflow INTEGER NOT NULL DEFAULT 0,
        is_mobile_friendly INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed by `gitmap email template responsive-check` and integrated into pre-dispatch pipeline checks.

---

### Improvement 1453: Tracking Pixel Decoupler & Local SVG Beacon Transcoder
- **ID:** `SEC-PIXEL-DECOUPLER-1453`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Replace third-party HTTP tracking pixels with privacy-preserving, zero-CDN local SVG markers and local webhook dispatchers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual pixel sanitizer summary: `Found 3 external tracking pixels -> Transcoded 3 to inline data:image/svg+xml`.
    - Decoupled domains list: `[STRIPPED] track.external-esp.com/pixel.gif -> [INJECTED] data:image/svg+xml;utf8,...`.
  - **Zero-CDN Web View:**
    - Clean DOM preview highlighting converted beacon tokens in soft teal CSS borders.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template decouple-pixels`
  - Flags:
    - `--input-file <string>`: Relative path to input template file.
    - `--output-file <string>`: Destination path for decoupled template.
    - `--is-svg-transcoding`: Convert external GIFs to zero-byte inline SVG data URIs (default: `true`).
    - `--has-privacy-first`: Strip user-identifiable URL query params from all links (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PixelDecouplerRecord struct {
        DecoupleID       string `json:"decouple_id"`
        InputTemplate    string `json:"input_template"`
        PixelsStripped   int    `json:"pixels_stripped"`
        SvgBeaconsAdded  int    `json:"svg_beacons_added"`
        IsPrivacySecured bool   `json:"is_privacy_secured"`
        HasExternalLinks bool   `json:"has_external_links"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PixelDecouplerRecord (
        decouple_id TEXT PRIMARY KEY,
        input_template TEXT NOT NULL,
        pixels_stripped INTEGER NOT NULL DEFAULT 0,
        svg_beacons_added INTEGER NOT NULL DEFAULT 0,
        is_privacy_secured INTEGER NOT NULL DEFAULT 1,
        has_external_links INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked via `gitmap email template decouple-pixels` prior to warmup outbound generation.

---

### Improvement 1454: Dark Mode CSS Variable Inverter & Color Contrast Validator
- **ID:** `SEC-DARKMODE-CONTRAST-1454`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify and automatically inject `@media (prefers-color-scheme: dark)` styling with WCAG AA compliant 4.5:1 text-to-background contrast ratios.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual contrast preview HUD comparing Light Mode vs Dark Mode colors.
    - Contrast badge: `Text: #1e293b on #ffffff [Ratio: 12.8:1 - AAA PASS] | Dark Mode: #f8fafc on #0f172a [Ratio: 15.1:1 - AAA PASS]`.
  - **Zero-CDN Web View:**
    - Live theme toggle slider switching pure local CSS custom properties without JavaScript dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template darkmode-audit`
  - Flags:
    - `--template-file <string>`: Path to HTML email template.
    - `--target-wcag-level <string>`: Contrast standard (`AA`, `AAA`; default: `AA`).
    - `--is-auto-fix`: Inject dark mode CSS media blocks if missing (default: `true`).
    - `--has-contrast-breach`: Return error code if any text element fails contrast ratio (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeContrastRecord struct {
        AuditID          string  `json:"audit_id"`
        TemplateName     string  `json:"template_name"`
        LightModeRatio   float64 `json:"light_mode_ratio"`
        DarkModeRatio    float64 `json:"dark_mode_ratio"`
        IsWcagCompliant  bool    `json:"is_wcag_compliant"`
        HasInjectedRules bool    `json:"has_injected_rules"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DarkModeContrastRecord (
        audit_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        light_mode_ratio REAL NOT NULL DEFAULT 4.5,
        dark_mode_ratio REAL NOT NULL DEFAULT 4.5,
        is_wcag_compliant INTEGER NOT NULL DEFAULT 1,
        has_injected_rules INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email template darkmode-audit` to guarantee email legibility across Apple Mail and Gmail dark themes.

---

### Improvement 1455: Font Subsetting & Inline Base64 WOFF2 Bundler
- **ID:** `SEC-FONT-WOFF2-BUNDLER-1455`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Subset typography to only characters used in the template and embed lightweight WOFF2 fonts as inline Base64 data URIs, eliminating Google Fonts CDN requests.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Compression statistics table: `Original Inter Font: 148 KB -> Subsetted (38 chars): 3.2 KB (-97.8% size)`.
    - Verification flag: `Zero External Network Font Calls: VERIFIED`.
  - **Zero-CDN Web View:**
    - Typography specimen cards rendered using the embedded Base64 `@font-face` definitions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template bundle-fonts`
  - Flags:
    - `--template-file <string>`: Path to template requiring font bundling.
    - `--font-name <string>`: Target typeface family to bundle (`inter`, `roboto`, `system`; default: `inter`).
    - `--is-subset-only`: Strip all glyphs not present in template body (default: `true`).
    - `--has-system-fallback`: Include standard sans-serif fallback stack (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FontBundleRecord struct {
        BundleID          string `json:"bundle_id"`
        TemplateName      string `json:"template_name"`
        FontFamily        string `json:"font_family"`
        OriginalFontBytes int    `json:"original_font_bytes"`
        SubsetFontBytes   int    `json:"subset_font_bytes"`
        IsFullyEmbedded   bool   `json:"is_fully_embedded"`
        HasZeroCdnCalls   bool   `json:"has_zero_cdn_calls"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FontBundleRecord (
        bundle_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        font_family TEXT NOT NULL,
        original_font_bytes INTEGER NOT NULL DEFAULT 0,
        subset_font_bytes INTEGER NOT NULL DEFAULT 0,
        is_fully_embedded INTEGER NOT NULL DEFAULT 1,
        has_zero_cdn_calls INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed via `gitmap email template bundle-fonts` to eliminate remote privacy leaks and CDN blocks.

---

### Improvement 1456: Email Payload Size & 102KB Gmail Clipping Guard
- **ID:** `SEC-GMAIL-CLIPPING-GUARD-1456`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Continuously audit raw MIME message payload sizes to prevent Gmail from clipping messages exceeding the 102KB ceiling.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Payload size gauge with visual safety threshold: `[████████████░░░░░░] 64.2 KB / 102.0 KB Max (Safe margin: 37.8 KB)`.
    - Warning banner triggered if payload exceeds 95KB: `[▲] APPROACHING 102KB CEILING (98.4 KB) - CLIPPING RISK!`.
  - **Zero-CDN Web View:**
    - HTML progress bar with color thresholds (green <80KB, yellow 80-100KB, red >102KB).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template check-clipping`
  - Flags:
    - `--template-file <string>`: Path to HTML template.
    - `--warn-threshold-kb <float>`: Warning threshold in kilobytes (default: `90.0`).
    - `--is-minify-enabled`: Automatically run HTML minifier if size breaches threshold (default: `true`).
    - `--has-clipping-risk`: Return non-zero exit code if payload exceeds 102KB (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GmailClippingGuardRecord struct {
        GuardID         string  `json:"guard_id"`
        TemplatePath    string  `json:"template_path"`
        PayloadBytes    int     `json:"payload_bytes"`
        SafetyMarginKb  float64 `json:"safety_margin_kb"`
        IsUnderCeiling  bool    `json:"is_under_ceiling"`
        HasMinifiedCopy bool    `json:"has_minified_copy"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GmailClippingGuardRecord (
        guard_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        payload_bytes INTEGER NOT NULL DEFAULT 0,
        safety_margin_kb REAL NOT NULL DEFAULT 102.0,
        is_under_ceiling INTEGER NOT NULL DEFAULT 1,
        has_minified_copy INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `gitmap email template check-clipping` as a pre-flight validator before template approval.

---

### Improvement 1457: Phishing & Malicious Anchor Link Disinfection Analyzer
- **ID:** `SEC-ANCHOR-DISINFECT-1457`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit all `<a>` tags for mismatched anchor text vs target URLs, IDN homograph attacks, and suspicious redirectors that trigger spam filters.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Anchor link audit scorecard showing inspected URLs, punycode resolutions, and target reputations.
    - Security flag: `[✓] All 6 outbound links match verified domain reputation whitelist`.
  - **Zero-CDN Web View:**
    - Tabular display of extracted URLs with visual warning chips for URL shorteners or deceptive anchor texts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template link-audit`
  - Flags:
    - `--template-file <string>`: Path to template to scan.
    - `--allowed-domains <string>`: Comma-separated list of permitted destination root domains.
    - `--is-homograph-checked`: Inspect domains for Unicode/Punycode homoglyphs (default: `true`).
    - `--has-deceptive-anchors`: Flag links where text appears to be a domain different from `href` (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type AnchorDisinfectRecord struct {
        DisinfectID        string `json:"disinfect_id"`
        TemplateName       string `json:"template_name"`
        ScannedLinkCount   int    `json:"scanned_link_count"`
        FlaggedLinkCount   int    `json:"flagged_link_count"`
        IsPunycodeSafe     bool   `json:"is_punycode_safe"`
        HasHomographRisk   bool   `json:"has_homograph_risk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AnchorDisinfectRecord (
        disinfect_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        scanned_link_count INTEGER NOT NULL DEFAULT 0,
        flagged_link_count INTEGER NOT NULL DEFAULT 0,
        is_punycode_safe INTEGER NOT NULL DEFAULT 1,
        has_homograph_risk INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed via `gitmap email template link-audit` to prevent reputation blacklisting caused by bad links.

---

### Improvement 1458: Pre-Header & Snippet Text Optimization Live Inspector
- **ID:** `SEC-PREHEADER-INSPECTOR-1458`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Manage invisible pre-header text and padding spacer characters to control the exact snippet displayed in mobile mailbox list views.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Mobile inbox preview mockup displaying Subject Line and Pre-header snippet side by side.
    - Layout:
      ```text
      ┌─ [MOBILE INBOX LIST VIEW SNIPPET PREVIEW] ───────────────────────────────────────────────┐
      │ SENDER: Team GitMap <warmup@gitmap.dev>                                                   │
      │ SUBJECT: Your daily infrastructure warmup report is ready                                 │
      │ SNIPPET: 99.4% seed deliverability achieved across 410 nodes today...                      │
      │ PADDING SPACERS: 85 &zwnj;&nbsp; entities injected (prevents body bleed-through)          │
      └───────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web View:**
    - HTML5 mobile phone viewport frame displaying simulated notification banners.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template preheader-tune`
  - Flags:
    - `--template-file <string>`: Path to HTML template.
    - `--preheader-text <string>`: New pre-header text string to inject.
    - `--is-padding-injected`: Inject zero-width non-joiners to hide email body text (default: `true`).
    - `--has-optimal-length`: Verify snippet is between 40 and 110 characters (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PreheaderInspectorRecord struct {
        InspectorID       string `json:"inspector_id"`
        TemplateName      string `json:"template_name"`
        PreheaderLength   int    `json:"preheader_length"`
        SpacerCount       int    `json:"spacer_count"`
        IsLengthOptimal   bool   `json:"is_length_optimal"`
        HasPaddingApplied bool   `json:"has_padding_applied"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PreheaderInspectorRecord (
        inspector_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        preheader_length INTEGER NOT NULL DEFAULT 0,
        spacer_count INTEGER NOT NULL DEFAULT 0,
        is_length_optimal INTEGER NOT NULL DEFAULT 1,
        has_padding_applied INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bound to `gitmap email template preheader-tune` to streamline sender presentation and open rates.

---

### Improvement 1459: Responsive Table Nesting & MSO Outlook VML Conditioner
- **ID:** `SEC-MSO-OUTLOOK-VML-1459`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically generate Microsoft Outlook Vector Markup Language (VML) and `<!--[if mso]>` conditionals for bulletproof rendering in Outlook 2016-2024 desktop clients.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - MSO compatibility checklist showing table width enforcement, ghost table injection, and VML button tags.
    - Badge: `[MSO Compatibility: 100% | 4 Ghost Tables Injected | 2 VML Buttons Generated]`.
  - **Zero-CDN Web View:**
    - Side-by-side rendering comparing standards-mode rendering with simulated Outlook Word-rendering engine.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template mso-conditioner`
  - Flags:
    - `--template-file <string>`: Path to HTML template.
    - `--is-ghost-table-enabled`: Wrap fluid divs with `<!--[if (gte mso 9)|(IE)]>` table wrappers (default: `true`).
    - `--has-vml-buttons`: Convert CSS round-corner buttons to native VML shape tags (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MsoOutlookConditionerRecord struct {
        ConditionerID    string `json:"conditioner_id"`
        TemplateName     string `json:"template_name"`
        GhostTablesAdded int    `json:"ghost_tables_added"`
        VmlButtonsAdded  int    `json:"vml_buttons_added"`
        IsMsoReady       bool   `json:"is_mso_ready"`
        HasVmlFallbacks  bool   `json:"has_vml_fallbacks"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MsoOutlookConditionerRecord (
        conditioner_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        ghost_tables_added INTEGER NOT NULL DEFAULT 0,
        vml_buttons_added INTEGER NOT NULL DEFAULT 0,
        is_mso_ready INTEGER NOT NULL DEFAULT 1,
        has_vml_fallbacks INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed by `gitmap email template mso-conditioner` to eliminate desktop Outlook rendering glitches.

---

### Improvement 1460: Zero-CDN Template Export & Cryptographic Fingerprint Registry
- **ID:** `SEC-TEMPLATE-EXPORT-FINGERPRINT-1460`
- **Cluster:** 146 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Export validated zero-CDN email templates and compute immutable SHA-256 fingerprints to register approved content versions in the repository database.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Export certificate card displaying template identity, SHA-256 fingerprint, byte size, and approval status.
    - Final confirmation: `[✓] TEMPLATE APPROVED & REGISTERED IN REPOSITORY DB: SHA-256: 7f83b165...`.
  - **Zero-CDN Web View:**
    - Downloadable certified HTML view with embedded cryptographic signature comment block in `<head>`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email template register`
  - Flags:
    - `--template-file <string>`: Path to finalized HTML template.
    - `--version-tag <string>`: Semantic version identifier for template (e.g. `v1.4.0`).
    - `--is-certified-safe`: Enforce passing all sanitizer, responsive, and clipping checks prior to registration (default: `true`).
    - `--has-immutable-hash`: Record SHA-256 hash in `repodb/pipeline.db` (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TemplateFingerprintRecord struct {
        TemplateID       string `json:"template_id"`
        VersionTag       string `json:"version_tag"`
        Sha256Checksum   string `json:"sha256_checksum"`
        RegisteredAtUnix int64  `json:"registered_at_unix"`
        IsCertifiedSafe  bool   `json:"is_certified_safe"`
        HasValidSignature bool  `json:"has_valid_signature"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TemplateFingerprintRecord (
        template_id TEXT PRIMARY KEY,
        version_tag TEXT NOT NULL,
        sha256_checksum TEXT NOT NULL,
        registered_at_unix INTEGER NOT NULL DEFAULT 0,
        is_certified_safe INTEGER NOT NULL DEFAULT 1,
        has_valid_signature INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - CLI command `gitmap email template register` stores final verified template fingerprints in `repodb/pipeline.db`.

---
