# Specification 347: Improvements 1031 to 1060 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 347
- **Slug:** `347-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/347-warm-up-ui-and-emailing-ui-improvements/03-improvements-1031-to-1060-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/346-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (1031 - 1060)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 104** | Items 1031–1040 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 105** | Items 1041–1050 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 106** | Items 1051–1060 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 104: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 1031–1040)

---

### Improvement 1031: Autonomous Multi-Tier Seed Pool Health Matrix & Placement Heatmap HUD
- **ID:** `SEED-MATRIX-1031`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Render a high-density, real-time terminal matrix table classifying active warmup seed mailboxes across primary consumer and enterprise ESP networks (Google Workspace, Microsoft 365, Yahoo/AOL, ProtonMail, and Apple iCloud) with granular placement rate telemetry, round-trip latency scoring, and spam rescue dispatch.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density terminal matrix table rendering active seed nodes with 4 distinct state indicators: `[INBOX]` (emerald `#10b981`), `[SPAM]` (crimson `#ef4444`), `[DELAY]` (amber `#f59e0b`), and `[MISS]` (slate `#64748b`).
    - Interactive cursor with Lipgloss border styling (`#06b6d4`, electric cyan) for keyboard navigation (`hjkl` / arrow keys) displaying instantaneous delivery headers, SPF/DKIM/DMARC alignment pass rate, and seed inbox ping latency.
    - Model view snippet:
      ```text
      ┌─ [SEED POOL VERIFICATION MATRIX: 240 ACTIVE SEED INBOXES] ──────────────────────────┐
      │ ESP CLUSTER       │ TOTAL │ INBOX RATE │ SPAM RATE │ AVG LATENCY │ STATUS            │
      ├───────────────────┼───────┼────────────┼───────────┼─────────────┼───────────────────┤
      │ Google Workspace  │    80 │      98.7% │      1.3% │      1.12s  │ [●] OPTIMAL       │
      │ Microsoft 365     │    70 │      96.2% │      3.8% │      1.74s  │ [●] OPTIMAL       │
      │ Yahoo / AOL       │    40 │      94.1% │      5.9% │      2.85s  │ [▲] ELEVATED LAG  │
      │ ProtonMail        │    30 │     100.0% │      0.0% │      0.79s  │ [●] PRISTINE      │
      │ Apple iCloud      │    20 │      89.5% │     10.5% │      3.40s  │ [■] ATTENTION REQ │
      ├───────────────────┴───────┴────────────┴───────────┴─────────────┴───────────────────┤
      │ [CURSOR: icloud-seed-19@icloud.com]                                                  │
      │ ├─ Last Arrival: SPAM (Folder: Junk) | Subject Hash: f4c901a8                        │
      │ ├─ Authentication: SPF=PASS, DKIM=PASS, DMARC=PASS (Strict Alignment)               │
      │ ├─ Engagement Agent Action: Auto-Moved to Inbox & Starred (0.8s ago)                 │
      │ └─ Keys: [enter: Inspect Headers] [t: Trigger Engagement] [r: Refresh] [q: Quit]     │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure HTML5 grid and inline SVG status indicators. Zero third-party web font or external CDN stylesheet dependencies.
    - Inline SVG sparklines depicting hourly inbox rate distributions per ESP category.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-matrix`
  - Flags:
    - `--esp <string>`: Filter seed pool by target provider (`google`, `microsoft`, `yahoo`, `proton`, `icloud`, or `all`; default: `all`).
    - `--min-inbox-rate <float>`: Highlight seed accounts failing target threshold (default: `90.0`).
    - `--is-streaming`: Enable continuous real-time live terminal streaming mode (default: `true`).
    - `--export-json <path>`: Dump instantaneous snapshot to relative JSON file.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedPoolNode struct {
        SeedID          string    `json:"seedId"`
        EmailAddress    string    `json:"emailAddress"`
        ProviderType    string    `json:"providerType"`
        InboxPlacement  int       `json:"inboxPlacement"`
        SpamPlacement   int       `json:"spamPlacement"`
        AverageLagSec   float64   `json:"averageLagSec"`
        LastCheckedAt   time.Time `json:"lastCheckedAt"`
        IsActive        bool      `json:"isActive"`
        IsHealthy       bool      `json:"isHealthy"`
        HasRecentRescue bool      `json:"hasRecentRescue"`
    }

    type SeedMatrixModel struct {
        Nodes           []SeedPoolNode `json:"nodes"`
        CursorIndex     int            `json:"cursorIndex"`
        SelectedNode    SeedPoolNode   `json:"selectedNode"`
        TotalInboxRatio float64        `json:"totalInboxRatio"`
        IsStreaming     bool           `json:"isStreaming"`
        HasAlerts       bool           `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedPoolVerificationMatrix (
        seed_id TEXT PRIMARY KEY,
        email_address TEXT NOT NULL,
        provider_type TEXT NOT NULL,
        inbox_placement INTEGER NOT NULL DEFAULT 0,
        spam_placement INTEGER NOT NULL DEFAULT 0,
        average_lag_sec REAL NOT NULL DEFAULT 0.0,
        last_checked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_active INTEGER NOT NULL DEFAULT 1,
        is_healthy INTEGER NOT NULL DEFAULT 1,
        has_recent_rescue INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_seedpool_provider_health ON SeedPoolVerificationMatrix(provider_type, is_healthy);
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded into `cli/email/warmup/seed_matrix.go` utilizing Bubble Tea and Lipgloss terminal styling. Connects to `repodb/pipeline.db` using the repository pattern.

---

### Improvement 1032: Interactive Seed Engagement Trajectory Timeline & Action Dispatcher
- **ID:** `SEED-TIMELINE-1032`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Visualize granular chronological interactions between seed mailboxes and sender domains, illustrating open events, thread replies, star actions, mark-not-spam corrections, and categorization changes in a navigable Bubble Tea timeline.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Vertical timeline stream using unicode connectors (`│`, `├─`, `└─`, `◆`, `◇`).
    - Event tags color-coded with Lipgloss tokens: `[OPEN]` (sky blue `#38bdf8`), `[REPLY]` (violet `#a855f7`), `[UNSPAM]` (emerald `#10b981`), `[STAR]` (amber `#fbbf24`), and `[LABEL]` (slate `#94a3b8`).
    - ANSI Preview:
      ```text
      ┌─ [SEED INTERACTION TRAJECTORY: seed-42@m365.warmup.io] ─────────────────────────────┐
      │ 14:02:18  ◆ [REPLY]   Replied with sentiment positive (Token count: 48)            │
      │           │           Thread ID: th-90812 | Response Lag: 42s                       │
      │ 14:01:36  ◆ [STAR]    Marked conversation as starred / important                   │
      │ 14:01:22  ◆ [UNSPAM]  Rescued message from Junk folder -> Moved to Primary Inbox    │
      │ 14:01:05  ◆ [OPEN]    Message opened, images loaded, duration: 18.4s               │
      │ 13:58:10  ◇ [RECEIVE] Message arrived at MTA gateway (SPF=PASS, DKIM=PASS)          │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [STATISTICS] Interaction Depth: 4 steps | Human Mimicry Score: 98.4%                │
      │ [ACTIONS]    (a) Manual Rescue  (r) Force Thread Reply  (d) Detail View  (q) Back   │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline HTML SVG rendering of interaction steps with hoverable DOM tooltips showing RFC headers and payload hashes.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-timeline`
  - Flags:
    - `--seed-id <string>`: Target seed mailbox identifier (required).
    - `--limit <int>`: Maximum number of trajectory events to render (default: `50`).
    - `--is-auto-refresh`: Poll timeline dynamically every 2 seconds (default: `true`).
    - `--has-replies-only`: Filter timeline to reply events only (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedEventTrajectory struct {
        EventID       string    `json:"eventId"`
        SeedID        string    `json:"seedId"`
        SenderDomain  string    `json:"senderDomain"`
        ThreadID      string    `json:"threadId"`
        EventType     string    `json:"eventType"`
        DwellDuration float64   `json:"dwellDuration"`
        OccurredAt    time.Time `json:"occurredAt"`
        IsAutomated   bool      `json:"isAutomated"`
        HasReplied    bool      `json:"hasReplied"`
    }

    type SeedTimelineModel struct {
        Events         []SeedEventTrajectory `json:"events"`
        SelectedSeedID string                `json:"selectedSeedId"`
        CursorPosition int                   `json:"cursorPosition"`
        IsAutoRefresh  bool                  `json:"isAutoRefresh"`
        HasError       bool                  `json:"hasError"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedEventTrajectory (
        event_id TEXT PRIMARY KEY,
        seed_id TEXT NOT NULL,
        sender_domain TEXT NOT NULL,
        thread_id TEXT NOT NULL,
        event_type TEXT NOT NULL,
        dwell_duration REAL NOT NULL DEFAULT 0.0,
        occurred_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_automated INTEGER NOT NULL DEFAULT 1,
        has_replied INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_seed_event_time ON SeedEventTrajectory(seed_id, occurred_at DESC);
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed under `cli/email/warmup/seed_timeline.go`. Streams real-time pipeline events from SQLite database.

---

### Improvement 1033: Mailbox Warmup Ramp Gradient & Daily Quota Governor Visualizer
- **ID:** `WARMUP-RAMP-1033`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide an interactive ramp curve visualizer indicating linear, quadratic, and logarithmic daily volume ramp schedules per mailbox, warning operators against spam filter velocity tripwires.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-precision ASCII bar chart graphing historical sends against projected safety curves.
    - Status cards with dynamic Lipgloss borders: green (`#22c55e`) for safe pacing, yellow (`#eab308`) for approaching limit, and red (`#ef4444`) for velocity lock.
    - Terminal UI preview:
      ```text
      ┌─ [MAILBOX RAMP GRADIENT: outreach-node-04@domain.com] ──────────────────────────────┐
      │ Ramp Stage: DAY 14 / 30 | Curve: QUADRATIC_SMOOTH | Target Cap: 80 msgs/day        │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ MSG/DAY                                                                             │
      │  80 │                                                            ╭──── [MAX CAP]    │
      │  60 │                                               ╭────────────╯                  │
      │  40 │                                  ╭────────────╯                               │
      │  20 │                     ╭────────────╯                                            │
      │   0 ┼───┬────┬────┬────┬──┴─┬────┬────┬────┬────┬────┬────┬────┬────┬───            │
      │     D1  D3   D5   D7   D9  D11  D13  D15  D17  D19  D21  D23  D25  D27             │
      │                                  ▲                                                  │
      │                           CURRENT POSITION                                          │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ Today Dispatched: 38 / 42 (90.4%) | Inbox Rate: 97.6% | Health Score: 96/100        │
      │ [FLAGS] IsThrottleActive: FALSE | HasVelocityTripped: FALSE | IsRampPaused: FALSE    │
      │ [KEYS]  (+) Bump Daily Target  (-) Lower Target  (p) Pause Ramp  (q) Exit           │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Embedded SVG dynamic line charts plotting 30-day velocity profiles with safe threshold fill zones.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ramp-status`
  - Flags:
    - `--mailbox <string>`: Email address of the warmup mailbox (required).
    - `--curve <string>`: Pacing curve type (`linear`, `quadratic`, `logarithmic`; default: `quadratic`).
    - `--is-verbose`: Display granular hourly delivery distribution (default: `false`).
    - `--has-warnings-only`: Filter display to mailboxes with pacing warnings (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxRampStatus struct {
        MailboxAddress    string    `json:"mailboxAddress"`
        CurrentDay        int       `json:"currentDay"`
        MaxDays           int       `json:"maxDays"`
        TodaySentCount    int       `json:"todaySentCount"`
        DailyTargetCount  int       `json:"dailyTargetCount"`
        MaximumCap        int       `json:"maximumCap"`
        CurveType         string    `json:"curveType"`
        LastDispatchedAt  time.Time `json:"lastDispatchedAt"`
        IsThrottleActive  bool      `json:"isThrottleActive"`
        HasVelocityAlert  bool      `json:"hasVelocityAlert"`
        IsRampPaused      bool      `json:"isRampPaused"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxRampSchedule (
        mailbox_address TEXT PRIMARY KEY,
        current_day INTEGER NOT NULL DEFAULT 1,
        max_days INTEGER NOT NULL DEFAULT 30,
        today_sent_count INTEGER NOT NULL DEFAULT 0,
        daily_target_count INTEGER NOT NULL DEFAULT 5,
        maximum_cap INTEGER NOT NULL DEFAULT 100,
        curve_type TEXT NOT NULL DEFAULT 'quadratic',
        last_dispatched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_throttle_active INTEGER NOT NULL DEFAULT 0,
        has_velocity_alert INTEGER NOT NULL DEFAULT 0,
        is_ramp_paused INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/email/warmup/ramp_visualizer.go`. Interacts with scheduler algorithms to guard against sudden spam-triggering send spikes.

---

### Improvement 1034: Deliverability Diagnostic Inspector & Seed Header Parsing HUD
- **ID:** `DIAG-HEADER-1034`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Enable deep inspection of raw RFC 5322 email headers captured upon seed inbox arrival, isolating `Authentication-Results`, `Received-SPF`, `DKIM-Signature`, `X-Forefront-Antispam-Report`, and `ARC-Authentication-Results` in an interactive viewer.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual-pane layout: Left pane displays message arrivals; right pane parses security headers with syntax-highlighted status flags.
    - Color highlighting: Green (`#4ade80`) for passing signatures, Red (`#f87171`) for failed/softfail/permerror verdicts, Blue (`#60a5fa`) for routing IP metadata.
    - ANSI Preview:
      ```text
      ┌─ [DIAGNOSTIC HEADER INSPECTOR: MSG-20261027-8912] ──────────────────────────────────┐
      │ ARRIVAL FEED                │ RFC AUTHENTICATION RESULTS                             │
      ├─────────────────────────────┼────────────────────────────────────────────────────────┤
      │ ▶ 14:15:22 m365-seed-01     │ Authentication-Results: mx.google.com;                 │
      │   14:12:05 gmail-seed-14    │   dkim=pass header.i=@outreach.domain.com;             │
      │   14:09:44 yahoo-seed-02    │   spf=pass (google.com: domain of sender designates    │
      │   14:05:12 proton-seed-05   │     198.51.100.24 as permitted sender);               │
      │   13:59:18 icloud-seed-08   │   dmarc=pass (p=REJECT sp=REJECT dis=NONE)             │
      │                             │ ARC-Authentication-Results: i=1; mx.google.com;        │
      │                             │   arc=pass (as=pass)                                   │
      │                             │ X-Spam-Score: -1.2 (Clean)                             │
      │                             │ X-Forefront-Antispam-Report: BCL:0;PCL:0;SFV:NSP;      │
      ├─────────────────────────────┴────────────────────────────────────────────────────────┤
      │ [STATUS] SPF=PASS | DKIM=PASS | DMARC=PASS | ARC=PASS | SCL=1 (Clean)                │
      │ [KEYS]   (j/k) Scroll Headers  (tab) Switch Pane  (c) Copy Raw  (q) Back             │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Monospace clean HTML view with collapsible header blocks and one-click raw copy functionality.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup header-inspect`
  - Flags:
    - `--message-id <string>`: Unique message identifier or hash (required).
    - `--is-raw-mode`: Display unmodified RFC headers without parsing (default: `false`).
    - `--has-auth-only`: Suppress standard transport headers and only display auth lines (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailHeaderReport struct {
        MessageID          string    `json:"messageId"`
        SeedRecipient      string    `json:"seedRecipient"`
        SenderIP           string    `json:"senderIp"`
        SPFVerdict         string    `json:"spfVerdict"`
        DKIMVerdict        string    `json:"dkimVerdict"`
        DMARCVerdict       string    `json:"dmarcVerdict"`
        SpamConfidenceLvl  int       `json:"spamConfidenceLvl"`
        RawHeaderPayload   string    `json:"rawHeaderPayload"`
        CapturedAt         time.Time `json:"capturedAt"`
        IsStrictAligned    bool      `json:"isStrictAligned"`
        HasArcValidation   bool      `json:"hasArcValidation"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailHeaderReports (
        message_id TEXT PRIMARY KEY,
        seed_recipient TEXT NOT NULL,
        sender_ip TEXT NOT NULL,
        spf_verdict TEXT NOT NULL,
        dkim_verdict TEXT NOT NULL,
        dmarc_verdict TEXT NOT NULL,
        spam_confidence_lvl INTEGER NOT NULL DEFAULT 0,
        raw_header_payload TEXT NOT NULL,
        captured_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_strict_aligned INTEGER NOT NULL DEFAULT 1,
        has_arc_validation INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_header_recipient ON EmailHeaderReports(seed_recipient);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/warmup/header_inspect.go`. Provides instant forensic insights into deliverability and authentication drops.

---

### Improvement 1035: Autonomous Reputation Trap & Honeypot Detection Sentinel
- **ID:** `SEED-HONEYPOT-1035`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Continuously evaluate seed mailbox pools and outbound contact lists for pristine spam traps, recycled honeypots, and typo domains, immediately isolating high-risk targets before dispatch.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-visibility warning alert modal when honeypot patterns are detected.
    - Security risk grade table: Red badge for `[TRAP DETECTED]`, Amber badge for `[RECYCLED RISKY]`, Green badge for `[VERIFIED SEED]`.
    - ANSI UI:
      ```text
      ┌─ [REPUTATION TRAP & HONEYPOT SENTINEL] ─────────────────────────────────────────────┐
      │ AUDITED POOL: 1,450 ADDRESSES | THREATS NEUTRALIZED: 3                                │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ RISK │ TARGET ADDRESS                │ REASON                                       │
      ├──────┼───────────────────────────────┼──────────────────────────────────────────────┤
      │ CRIT │ abuse-collector@spambait.org  │ Matched known Project HoneyPot MX footprint  │
      │ HIGH │ info@micros0ft-support.com    │ Typo-squatted domain target with MX trap     │
      │ WARN │ dormant-user-91@yahoo.com     │ Inactive > 36 months, potential recycled MX  │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [AUTO-ACTION] Targets flagged and quarantined to Prevent Reputation Degradation     │
      │ [FLAGS]       HasTrapFound: TRUE | IsQuarantineEnforced: TRUE | IsAutoPurge: TRUE   │
      │ [COMMANDS]    (q) Quarantine  (w) Whitelist Override  (e) Export List  (x) Dismiss  │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Quarantined table with danger banners and threat intelligence breakdown.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup honeypot-scan`
  - Flags:
    - `--pool-id <string>`: Identifier of the mailbox list or seed pool (default: `default`).
    - `--is-auto-quarantine`: Automatically isolate suspected traps from dispatch queues (default: `true`).
    - `--has-strict-heuristics`: Enable aggressive MX and age heuristics (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HoneypotAuditResult struct {
        AuditID              string    `json:"auditId"`
        TargetAddress        string    `json:"targetAddress"`
        ThreatLevel          string    `json:"threatLevel"`
        DetectionReason      string    `json:"detectionReason"`
        MXRecordHost         string    `json:"mxRecordHost"`
        DetectedAt           time.Time `json:"detectedAt"`
        IsQuarantined        bool      `json:"isQuarantined"`
        HasFalsePositiveFlag bool      `json:"hasFalsePositiveFlag"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HoneypotAuditResults (
        audit_id TEXT PRIMARY KEY,
        target_address TEXT NOT NULL,
        threat_level TEXT NOT NULL,
        detection_reason TEXT NOT NULL,
        mx_record_host TEXT NOT NULL,
        detected_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_quarantined INTEGER NOT NULL DEFAULT 1,
        has_false_positive_flag INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_honeypot_target ON HoneypotAuditResults(target_address);
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/email/warmup/honeypot_sentinel.go`. Blocks risky email addresses during warmup and outbound operations.

---

### Improvement 1036: Seed Pool Geo-Distribution & ISP Diversity Analyzer
- **ID:** `SEED-GEODIST-1036`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Analyze and visualize the geographical and autonomous system (ASN) diversity of the seed mailbox network, preventing artificial single-datacenter warmup clustering.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Unicode ASCII world/regional map visualizer highlighting seed distribution concentrations across North America, Europe, Asia-Pacific, and Latin America.
    - Diversity score index bar (0–100%) indicating ASN and ISP dispersion quality.
    - ANSI Layout:
      ```text
      ┌─ [SEED POOL GEO & ISP DIVERSITY ANALYZER] ──────────────────────────────────────────┐
      │ TOTAL SEEDS: 250 | DIVERSITY INDEX: 94.2% [EXEMPLARY]                               │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ REGION          │ SEED COUNT │ PRIMARY ASNS                    │ STATUS             │
      ├─────────────────┼────────────┼─────────────────────────────────┼────────────────────┤
      │ North America   │        115 │ AS15169 (Google), AS8075 (MSFT) │ [●] BALANCED       │
      │ Europe (EU/UK)  │         75 │ AS3356 (Lumen), AS1273 (Vodaf) │ [●] BALANCED       │
      │ Asia-Pacific    │         42 │ AS4755 (TATA), AS4637 (Telstra) │ [●] BALANCED       │
      │ Latin America   │         18 │ AS27699 (Claro), AS7738 (Telef) │ [▲] NEED +10 SEEDS │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [WARNING] ASN 15169 concentration is 38.0% (Threshold: < 40.0% required)           │
      │ [KEYS]    (r) Rebalance Seed Pool  (f) Fetch ASN Details  (q) Exit                  │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG choropleth map highlighting seed locations with no external mapping CDN dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup geo-diversity`
  - Flags:
    - `--min-diversity-score <float>`: Alert if overall score falls below threshold (default: `85.0`).
    - `--is-cluster-view`: Group distribution by ASN rather than continent (default: `false`).
    - `--has-rebalance-recommendations`: Output suggested seed acquisitions (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GeoDistributionEntry struct {
        RegionCode        string    `json:"regionCode"`
        SeedCount         int       `json:"seedCount"`
        PercentageShare   float64   `json:"percentageShare"`
        TopASNs           []string  `json:"topAsns"`
        LastAuditedAt     time.Time `json:"lastAuditedAt"`
        IsWithinTolerance bool      `json:"isWithinTolerance"`
        HasUnderCapacity  bool      `json:"hasUnderCapacity"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedGeoDistribution (
        region_code TEXT PRIMARY KEY,
        seed_count INTEGER NOT NULL DEFAULT 0,
        percentage_share REAL NOT NULL DEFAULT 0.0,
        top_asns TEXT NOT NULL,
        last_audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_within_tolerance INTEGER NOT NULL DEFAULT 1,
        has_under_capacity INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded into `cli/email/warmup/geo_diversity.go`. Ensures warmup interactions closely simulate organic, globally distributed consumer behavior.

---

### Improvement 1037: Thread Reply Sentiment & LLM Persona Tuning Console
- **ID:** `SEED-PERSONA-1037`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide an interactive configuration and test console for LLM-driven organic conversation replies between seed mailboxes and warmup senders, supporting custom personas and sentiment variance.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Split-screen model: Left pane lists seed persona presets (`Executive B2B`, `Curious Consumer`, `Technical Lead`, `Procurement Specialist`); right pane provides a live response generator preview.
    - Lipgloss border framing with live token cost and reply readability metrics (Flesch-Kincaid score).
    - Terminal UI:
      ```text
      ┌─ [SEED REPLY PERSONA & SENTIMENT TUNING CONSOLE] ───────────────────────────────────┐
      │ ACTIVE PERSONA: B2B_EXECUTIVE_BUSY | SENTIMENT: POSITIVE_INTERESTED                 │
      ├────────────────────────────────────┬────────────────────────────────────────────────┤
      │ AVAILABLE PERSONAS                 │ GENERATED REPLY PREVIEW                        │
      ├────────────────────────────────────┼────────────────────────────────────────────────┤
      │ ▶ Executive B2B (Concise)          │ "Thanks for reaching out, Alex. The architecture│
      │   Curious Consumer (Inquisitive)   │ breakdown looks promising. Could you share the │
      │   Engineering Lead (Deep Dive)     │ benchmark comparisons against Redis? Let's talk│
      │   Procurement Lead (Budget Focus)  │ next Thursday at 2 PM EST."                    │
      │                                    ├────────────────────────────────────────────────┤
      │                                    │ Readability: Grade 8.2 | Sentiment Score: +0.82│
      │                                    │ Word Count: 31 | Tokens: 44 | Latency: 420ms   │
      ├────────────────────────────────────┴────────────────────────────────────────────────┤
      │ [FLAGS] IsContextAware: TRUE | HasHumanTypoSimulation: TRUE | IsSpamProof: TRUE     │
      │ [KEYS]  (p) Change Persona  (g) Re-Generate  (s) Save Preset  (q) Exit              │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live text editor with immediate sentiment gauge meter and simulated typist latency curves.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup persona-console`
  - Flags:
    - `--persona <string>`: Active persona preset name (default: `b2b_executive`).
    - `--has-typo-simulation`: Inject organic micro-typos and corrections (default: `true`).
    - `--is-sentiment-forced`: Require positive reply sentiment (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedReplyPersona struct {
        PersonaID              string    `json:"personaId"`
        PersonaName            string    `json:"personaName"`
        ToneProfile            string    `json:"toneProfile"`
        AverageLengthWords     int       `json:"averageLengthWords"`
        SentimentBias          float64   `json:"sentimentBias"`
        CreatedAt              time.Time `json:"createdAt"`
        IsContextAware         bool      `json:"isContextAware"`
        HasTypoSimulation      bool      `json:"hasTypoSimulation"`
        IsActive               bool      `json:"isActive"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedReplyPersonas (
        persona_id TEXT PRIMARY KEY,
        persona_name TEXT NOT NULL,
        tone_profile TEXT NOT NULL,
        average_length_words INTEGER NOT NULL DEFAULT 35,
        sentiment_bias REAL NOT NULL DEFAULT 0.75,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_context_aware INTEGER NOT NULL DEFAULT 1,
        has_typo_simulation INTEGER NOT NULL DEFAULT 1,
        is_active INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/warmup/persona_console.go`. Tunes the organic conversation generation engine to bypass modern AI spam detection heuristics.

---

### Improvement 1038: Mailbox Engagement Rescue & Auto-Categorization Automation
- **ID:** `SEED-RESCUE-1038`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Execute automated recovery routines on seed mailboxes when warmup messages land in spam or promotions folders, moving them back to the primary inbox, marking as important, and establishing positive sender affinity.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live feed of autonomous rescue operations executing via IMAP/Graph/Gmail APIs.
    - Status pills: `[RESCUED FROM SPAM]` in bold emerald, `[PROMOTIONS -> PRIMARY]` in electric blue, `[STARRED]` in gold.
    - Terminal UI:
      ```text
      ┌─ [AUTONOMOUS SEED RESCUE & CATEGORIZATION DISPATCHER] ──────────────────────────────┐
      │ ACTIVE SESSIONS: 12 IMAP / 8 GMAIL API / 6 GRAPH API                                │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ TIMESTAMP │ SEED MAILBOX              │ SOURCE FOLDER │ ACTION APPLIED              │
      ├───────────┼───────────────────────────┼───────────────┼─────────────────────────────┤
      │ 14:22:01  │ test-seed-81@outlook.com  │ Junk / Spam   │ Moved to Inbox + Starred    │
      │ 14:21:49  │ marketing-seed@gmail.com  │ Promotions    │ Moved to Primary + Starred  │
      │ 14:20:15  │ dev-seed-04@yahoo.com     │ Spam          │ Unmarked Spam + Replied     │
      │ 14:18:33  │ seed-vip-12@icloud.com    │ Junk          │ Moved to Inbox              │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ Total Rescues Today: 47 | Sender Reputation Uplift Factor: +12.4%                   │
      │ [FLAGS] IsRescueEnabled: TRUE | HasAutoStarring: TRUE | HasReplyDispatch: TRUE      │
      │ [KEYS]  (p) Pause Automation  (r) Run Manual Sweep  (l) View Logs  (q) Exit         │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live event log table updating via lightweight HTML server-sent events (SSE).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup auto-rescue`
  - Flags:
    - `--is-continuous`: Run rescue daemon continuously in the background (default: `true`).
    - `--has-auto-starring`: Automatically apply VIP/starred flag upon rescue (default: `true`).
    - `--rescue-delay-sec <int>`: Humanization delay before executing rescue (default: `15`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedRescueEvent struct {
        RescueID         string    `json:"rescueId"`
        SeedMailbox      string    `json:"seedMailbox"`
        SenderAddress    string    `json:"senderAddress"`
        OriginalFolder   string    `json:"originalFolder"`
        TargetFolder     string    `json:"targetFolder"`
        RescuedAt        time.Time `json:"rescuedAt"`
        IsRescueEnabled  bool      `json:"isRescueEnabled"`
        HasAutoStarring  bool      `json:"hasAutoStarring"`
        HasReplied       bool      `json:"hasReplied"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedRescueEvents (
        rescue_id TEXT PRIMARY KEY,
        seed_mailbox TEXT NOT NULL,
        sender_address TEXT NOT NULL,
        original_folder TEXT NOT NULL,
        target_folder TEXT NOT NULL,
        rescued_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_rescue_enabled INTEGER NOT NULL DEFAULT 1,
        has_auto_starring INTEGER NOT NULL DEFAULT 1,
        has_replied INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_rescue_sender ON SeedRescueEvents(sender_address);
    ```
- **Antigravity / GitMap Integration Point:**
  - Under `cli/email/warmup/auto_rescue.go`. Directs mailbox client APIs to simulate natural recipient rescue behavior.

---

### Improvement 1039: Warmup Pool Quarantine & Degraded Node Eviction Manager
- **ID:** `SEED-QUARANTINE-1039`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Automatically detect dead, locked, or blacklisted seed mailboxes, isolating them from the active warmup network to preserve sender domain reputation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Split list: Active Healthy Pool vs. Quarantined Eviction Staging.
    - Quarantine reason tags: `[IMAP AUTH FAIL]` (red), `[CHALLENGE PROMPT]` (orange), `[INBOX 0% CAP]` (yellow).
    - Terminal UI:
      ```text
      ┌─ [SEED POOL QUARANTINE & DEGRADED NODE EVICTION] ───────────────────────────────────┐
      │ ACTIVE NODES: 232 | QUARANTINED NODES: 8 | AUTOMATIC CLEANUP: ACTIVE                 │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ SEED ADDRESS              │ STATUS       │ REASON                │ DOWNTIME         │
      ├───────────────────────────┼──────────────┼───────────────────────┼──────────────────┤
      │ seed-99@outlook.com       │ QUARANTINED  │ Account locked by MS365│ 4h 12m          │
      │ tester-04@yahoo.com       │ QUARANTINED  │ IMAP Authentication 535│ 1h 05m          │
      │ worker-19@proton.me       │ EVALUATING   │ Delivery failure > 3x │ 0h 25m           │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [ACTIONS] (r) Re-Test Auth  (e) Evict Node Permanently  (u) Un-Quarantine  (q) Back  │
      │ [FLAGS]   IsEvictionActive: TRUE | HasAutoReplace: TRUE | IsStrictQuarantine: TRUE   │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Clean table with immediate action triggers to re-verify credentials or swap seed nodes.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup quarantine-manager`
  - Flags:
    - `--is-auto-evict`: Permanently remove nodes unrecovered after 24h (default: `true`).
    - `--max-error-threshold <int>`: Number of consecutive errors before quarantine (default: `3`).
    - `--has-auto-replace`: Spin up standby replacement seeds on eviction (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type QuarantinedSeedNode struct {
        SeedAddress       string    `json:"seedAddress"`
        QuarantineReason  string    `json:"quarantineReason"`
        ErrorCount        int       `json:"errorCount"`
        QuarantinedAt     time.Time `json:"quarantinedAt"`
        IsEvictionActive  bool      `json:"isEvictionActive"`
        HasAutoReplace    bool      `json:"hasAutoReplace"`
        IsPermanentBanned bool      `json:"isPermanentBanned"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS QuarantinedSeedNodes (
        seed_address TEXT PRIMARY KEY,
        quarantine_reason TEXT NOT NULL,
        error_count INTEGER NOT NULL DEFAULT 1,
        quarantined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_eviction_active INTEGER NOT NULL DEFAULT 1,
        has_auto_replace INTEGER NOT NULL DEFAULT 1,
        is_permanent_banned INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/email/warmup/quarantine_manager.go`. Protects the warmup cluster from sending into poisoned seed accounts.

---

### Improvement 1040: Seed Inbound Velocity & IMAP Polling Performance Dashboard
- **ID:** `SEED-POLL-1040`
- **Cluster:** 104 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Monitor and optimize IMAP IDLE connections and scheduled polling intervals across hundreds of seed accounts, preventing provider rate limiting while ensuring sub-second notification of message arrivals.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Real-time connection pool status monitor showing IDLE worker threads, connection latency, socket errors, and active keep-alive pings.
    - Terminal UI:
      ```text
      ┌─ [SEED INBOUND VELOCITY & IMAP POLLING PERFORMANCE] ────────────────────────────────┐
      │ ACTIVE SOCKETS: 128 | IDLE THREADS: 120 | AVG ROUNDTRIP: 142ms                      │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ PROTOCOL   │ CONNECTED │ ACTIVE POLLING │ RATE LIMIT THROTTLES │ SOCKET ERRORS/HR   │
      ├────────────┼───────────┼────────────────┼──────────────────────┼────────────────────┤
      │ IMAP IDLE  │        95 │            N/A │                    0 │                  0 │
      │ IMAP POLL  │        33 │         30 sec │                    1 │                  2 │
      │ GMAIL API  │        80 │          PUSH  │                    0 │                  0 │
      │ MS GRAPH   │        42 │          PUSH  │                    0 │                  0 │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [METRICS] Inbound Velocity: 14.2 msgs/min | P99 Latency: 420ms | Status: HEALTHY    │
      │ [KEYS]    (t) Tune Polling Frequency  (r) Reconnect Sockets  (q) Exit                │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live latency gauges and connection state telemetry rendered in pure CSS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup poll-metrics`
  - Flags:
    - `--poll-interval-sec <int>`: Set base IMAP polling cadence (default: `30`).
    - `--is-idle-preferred`: Enforce IMAP IDLE push where supported (default: `true`).
    - `--has-rate-limit-alerts`: Emit warning sound or notification on 429 errors (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type InboundPollMetric struct {
        ProtocolType      string    `json:"protocolType"`
        ActiveConnections int       `json:"activeConnections"`
        AveragePingMs     float64   `json:"averagePingMs"`
        ThrottleCount     int       `json:"throttleCount"`
        SocketErrorCount  int       `json:"socketErrorCount"`
        RecordedAt        time.Time `json:"recordedAt"`
        IsIdlePreferred   bool      `json:"isIdlePreferred"`
        HasRecentThrottle bool      `json:"hasRecentThrottle"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InboundPollMetrics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        protocol_type TEXT NOT NULL,
        active_connections INTEGER NOT NULL DEFAULT 0,
        average_ping_ms REAL NOT NULL DEFAULT 0.0,
        throttle_count INTEGER NOT NULL DEFAULT 0,
        socket_error_count INTEGER NOT NULL DEFAULT 0,
        recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_idle_preferred INTEGER NOT NULL DEFAULT 1,
        has_recent_throttle INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `cli/email/warmup/poll_metrics.go`. Coordinates with IMAP worker pools to maintain persistent connections without tripping security firewalls.

---

# Cluster 105: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 1041–1050)

---

### Improvement 1041: Multi-Domain Concurrency Governor & Per-ESP Rate Limiter HUD
- **ID:** `CONCUR-GOV-1041`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Dynamically govern outbound message dispatch velocity across multiple sender domains, enforcing distinct concurrency envelopes for target MX providers (e.g., maximum 5 simultaneous connections to Google, 3 to Microsoft, 2 to Yahoo).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-bar gauge matrix depicting active concurrency slots against provider thresholds.
    - Color coding: Cyan (`#06b6d4`) for normal concurrency, Orange (`#f97316`) for near-cap saturation, and Violet (`#a855f7`) for dynamic queuing.
    - Terminal UI:
      ```text
      ┌─ [MULTI-DOMAIN CONCURRENCY GOVERNOR & ESP RATE LIMITS] ─────────────────────────────┐
      │ GLOBAL SLOTS: 42 / 100 IN-FLIGHT | DISPATCH ENGINE: ADAPTIVE LEAKY BUCKET          │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ DESTINATION ESP   │ ACTIVE / MAX SLOTS │ SATURATION │ CURRENT TPS │ QUEUE DEPTH     │
      ├───────────────────┼────────────────────┼────────────┼─────────────┼─────────────────┤
      │ Google (AS15169)  │        12 / 20     │ [████░░░░] │  4.2 msg/s  │ 18 pending      │
      │ Microsoft (AS8075)│         8 / 15     │ [█████░░░] │  2.8 msg/s  │ 12 pending      │
      │ Yahoo/AOL (AS3471)│         5 /  8     │ [██████░░] │  1.4 msg/s  │  7 pending      │
      │ iCloud (AS714)    │         2 /  5     │ [████░░░░] │  0.8 msg/s  │  2 pending      │
      │ Other / Corporate │        15 / 50     │ [███░░░░░] │  5.1 msg/s  │ 35 pending      │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [STATUS] No 421/450 backoff active. Throughput optimal.                             │
      │ [FLAGS]  IsGovernorActive: TRUE | HasAutoThrottle: TRUE | IsStrictConcur: TRUE      │
      │ [KEYS]   (+) Increase Global Slots  (-) Reduce Slots  (p) Pause Governor  (q) Exit  │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Real-time SVG bar charts and slot utilization dials with pure CSS animations.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor status`
  - Flags:
    - `--max-global-slots <int>`: Total global concurrent outbound connections (default: `100`).
    - `--is-auto-throttle`: Dynamically reduce concurrency on remote 4xx deferrals (default: `true`).
    - `--has-strict-concurrency`: Hard cap limits without burst allowance (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ESPConcurrencyStatus struct {
        ESPName             string    `json:"espName"`
        ActiveSlots         int       `json:"activeSlots"`
        MaximumSlots        int       `json:"maximumSlots"`
        CurrentTPS          float64   `json:"currentTps"`
        PendingQueueDepth   int       `json:"pendingQueueDepth"`
        LastUpdated         time.Time `json:"lastUpdated"`
        IsGovernorActive    bool      `json:"isGovernorActive"`
        HasAutoThrottle     bool      `json:"hasAutoThrottle"`
        IsThrottled         bool      `json:"isThrottled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ESPConcurrencyStatus (
        esp_name TEXT PRIMARY KEY,
        active_slots INTEGER NOT NULL DEFAULT 0,
        maximum_slots INTEGER NOT NULL DEFAULT 10,
        current_tps REAL NOT NULL DEFAULT 0.0,
        pending_queue_depth INTEGER NOT NULL DEFAULT 0,
        last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_governor_active INTEGER NOT NULL DEFAULT 1,
        has_auto_throttle INTEGER NOT NULL DEFAULT 1,
        is_throttled INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/email/governor/concurrency_hud.go`. Manages connection token buckets in the outbound dispatch queue.

---

### Improvement 1042: Dynamic Backoff Inversion Engine & Exponential Jitter Modulator
- **ID:** `BACKOFF-INV-1042`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Invert retry schedules upon detecting greylisting or 4xx temporary deferrals, replacing rigid intervals with decorrelated full jitter and MX-specific backoff curves to prevent herd reconnection stamps.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Interactive graph detailing backoff delay distribution across retry attempts.
    - Live breakdown of SMTP deferred codes (`421 4.7.0 Try again later`, `451 4.4.0 DNS query failed`).
    - Terminal UI:
      ```text
      ┌─ [DYNAMIC BACKOFF INVERSION & RETRY JITTER MODULATOR] ──────────────────────────────┐
      │ ACTIVE DEFERRALS: 14 | RESOLUTION RATE: 92.8% | CURVE: FULL_DECORRELATED_JITTER    │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ RETRY # │ BASE DELAY │ MIN JITTER │ MAX JITTER │ ACTUAL APPLIED DELAY │ INVERT STATE │
      ├─────────┼────────────┼────────────┼────────────┼──────────────────────┼──────────────┤
      │ Try 1   │     30 sec │     12 sec │     45 sec │     28 sec           │ [●] INVERTED │
      │ Try 2   │    120 sec │     60 sec │    180 sec │    142 sec           │ [●] INVERTED │
      │ Try 3   │    600 sec │    300 sec │    900 sec │    715 sec           │ [●] INVERTED │
      │ Try 4   │   3600 sec │   1800 sec │   5400 sec │   4102 sec           │ [▲] BACKOFF  │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [LATEST 4XX CODE] 421 4.7.28 [Google] Rate limit exceeded for IP. Backoff engaged.   │
      │ [FLAGS]           IsJitterEnabled: TRUE | HasAutoInvert: TRUE | IsStrictDecorrel: TRUE│
      │ [KEYS]            (r) Reset Backoff State  (j) Test Jitter Formula  (q) Exit         │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Mathematical curve visualizer showing probability distribution density for next retry window.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor backoff-tune`
  - Flags:
    - `--curve <string>`: Jitter formula (`decorrelated`, `full_jitter`, `equal_jitter`; default: `decorrelated`).
    - `--base-delay-sec <int>`: Initial retry backoff seconds (default: `30`).
    - `--max-backoff-sec <int>`: Upper boundary ceiling for backoff (default: `7200`).
    - `--is-jitter-enabled`: Enable random jitter to break resonance (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BackoffInversionPlan struct {
        DestinationMX      string    `json:"destinationMx"`
        RetryAttempt       int       `json:"retryAttempt"`
        BaseDelaySeconds   int       `json:"baseDelaySeconds"`
        ActualDelaySeconds int       `json:"actualDelaySeconds"`
        DeferredCode       string    `json:"deferredCode"`
        NextScheduledAt    time.Time `json:"nextScheduledAt"`
        IsJitterEnabled    bool      `json:"isJitterEnabled"`
        HasAutoInvert      bool      `json:"hasAutoInvert"`
        IsResolved         bool      `json:"isResolved"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BackoffInversionPlans (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        destination_mx TEXT NOT NULL,
        retry_attempt INTEGER NOT NULL DEFAULT 1,
        base_delay_seconds INTEGER NOT NULL DEFAULT 30,
        actual_delay_seconds INTEGER NOT NULL DEFAULT 30,
        deferred_code TEXT NOT NULL,
        next_scheduled_at DATETIME NOT NULL,
        is_jitter_enabled INTEGER NOT NULL DEFAULT 1,
        has_auto_invert INTEGER NOT NULL DEFAULT 1,
        is_resolved INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_backoff_mx_sched ON BackoffInversionPlans(destination_mx, next_scheduled_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/governor/backoff_engine.go`. Recalculates deferral timings within the delivery scheduler.

---

### Improvement 1043: Multi-IP Warmup Pool Circuit Breaker & Automatic Failover
- **ID:** `CONCUR-BREAKER-1043`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Implement an autonomous circuit breaker across dedicated sending IP addresses, tripping an IP into quarantine upon detecting hard blocklists (Spamhaus, Barracuda) or spike in 5xx rejection rates.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live matrix of sending IPs with state indicator badges: `[CLOSED / HEALTHY]` (green), `[HALF-OPEN / PROBING]` (yellow), `[OPEN / TRIPPED]` (red).
    - Trip counter and automatic failover target routing table.
    - Terminal UI:
      ```text
      ┌─ [IP POOL CIRCUIT BREAKER & FAILOVER MONITOR] ─────────────────────────────────────┐
      │ ACTIVE IP POOL: 6 DEDICATED IPS | TRIPPED BREAKERS: 1 | FAILOVER TRAFFIC: 16%       │
      ├────────────────────────────────────────────────────────────────────────────────────┤
      │ IP ADDRESS     │ REPUTATION │ STATE     │ 5XX RATE (5m) │ FAILOVER ROUTE           │
      ├────────────────┼────────────┼───────────┼───────────────┼──────────────────────────┤
      │ 198.51.100.20  │ 99/100     │ CLOSED    │ 0.02%         │ Normal Routing           │
      │ 198.51.100.21  │ 98/100     │ CLOSED    │ 0.05%         │ Normal Routing           │
      │ 198.51.100.22  │ 61/100     │ OPEN      │ 8.40% [TRIP]  │ -> Failover: .20 & .21   │
      │ 198.51.100.23  │ 84/100     │ HALF-OPEN │ 1.10%         │ Probe Traffic (10%)      │
      ├────────────────────────────────────────────────────────────────────────────────────┤
      │ [BREAKER EVENT] IP 198.51.100.22 tripped: Spike in 550 5.7.1 blocklist responses.   │
      │ [FLAGS]         IsBreakerArmed: TRUE | HasAutoReroute: TRUE | IsStrictQuarantine:TRUE│
      │ [KEYS]          (r) Reset Breaker  (t) Trigger Manual Trip  (p) Probe IP  (q) Exit   │
      └────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Graphical circuit breaker state machine with immediate toggle controls.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor circuit-breaker`
  - Flags:
    - `--error-rate-threshold <float>`: 5xx rejection percentage to trip breaker (default: `5.0`).
    - `--half-open-timeout-sec <int>`: Duration to remain open before half-open probe (default: `1800`).
    - `--is-auto-reroute`: Automatically re-allocate pending queue to clean IPs (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CircuitBreakerNode struct {
        IPAddress         string    `json:"ipAddress"`
        ReputationScore   int       `json:"reputationScore"`
        CircuitState      string    `json:"circuitState"`
        RecentErrorRate   float64   `json:"recentErrorRate"`
        FailoverTargetIP  string    `json:"failoverTargetIp"`
        LastTrippedAt     time.Time `json:"lastTrippedAt"`
        IsBreakerArmed    bool      `json:"isBreakerArmed"`
        HasAutoReroute    bool      `json:"hasAutoReroute"`
        IsQuarantined     bool      `json:"isQuarantined"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CircuitBreakerNodes (
        ip_address TEXT PRIMARY KEY,
        reputation_score INTEGER NOT NULL DEFAULT 100,
        circuit_state TEXT NOT NULL DEFAULT 'CLOSED',
        recent_error_rate REAL NOT NULL DEFAULT 0.0,
        failover_target_ip TEXT NOT NULL DEFAULT '',
        last_tripped_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_breaker_armed INTEGER NOT NULL DEFAULT 1,
        has_auto_reroute INTEGER NOT NULL DEFAULT 1,
        is_quarantined INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `cli/email/governor/circuit_breaker.go`. Prevents systemic domain damage when single IP reputation is compromised.

---

### Improvement 1044: Sender Domain Reputation Drift Tracker & Velocity Dampener
- **ID:** `CONCUR-DRIFT-1044`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Continuously track reputation drift across Postmaster Tools and SNDS metrics, automatically applying a proportional velocity dampening factor if a domain's trust score slips.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual chart tracking sender domain reputation alongside current throttle coefficient ($0.1\times$ to $1.0\times$).
    - Color indicator: High reputation in lime (`#84cc16`), medium in yellow (`#eab308`), decaying in red (`#ef4444`).
    - Terminal UI:
      ```text
      ┌─ [DOMAIN REPUTATION DRIFT & VELOCITY DAMPENER] ─────────────────────────────────────┐
      │ MONITORED DOMAINS: 8 | DRIFT DETECTED: 1 | ENGINE: AUTO-DAMPENING                   │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ DOMAIN NAME          │ REPUTATION │ 7D DRIFT │ DAMPENING FACTOR │ EFFECTIVE RATE    │
      ├──────────────────────┼────────────┼──────────┼──────────────────┼───────────────────┤
      │ outreach.alpha.io    │ HIGH (98)  │ +1.2%    │ 1.00x (Full)     │ 50 msg/hr         │
      │ connect.beta.io      │ MED  (76)  │ -4.8%    │ 0.65x (Dampened) │ 32 msg/hr         │
      │ notify.gamma.io      │ HIGH (94)  │ -0.4%    │ 0.95x (Nominal)  │ 47 msg/hr         │
      │ cold.delta.io        │ LOW  (48)  │ -12.1%   │ 0.20x (Restricted│ 10 msg/hr         │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [ALERT] cold.delta.io dropped below threshold. Daily volume reduced by 80%.         │
      │ [FLAGS] IsDampenerActive: TRUE | HasAutoBrake: TRUE | IsReputationGated: TRUE       │
      │ [KEYS]  (d) Detail Audit  (o) Override Dampener  (r) Refresh Reputation  (q) Exit   │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live radar charts displaying multi-dimensional reputation drift metrics.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor reputation-drift`
  - Flags:
    - `--domain <string>`: Target sending domain name (optional; defaults to all).
    - `--min-reputation <int>`: Minimum score before triggering dampener (default: `80`).
    - `--is-auto-brake`: Automatically slow send cadence when drift is negative (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainReputationDrift struct {
        DomainName         string    `json:"domainName"`
        CurrentScore       int       `json:"currentScore"`
        SevenDayDriftPct   float64   `json:"sevenDayDriftPct"`
        DampeningFactor    float64   `json:"dampeningFactor"`
        EffectiveHourlyTPS float64   `json:"effectiveHourlyTps"`
        AuditedAt          time.Time `json:"auditedAt"`
        IsDampenerActive   bool      `json:"isDampenerActive"`
        HasAutoBrake       bool      `json:"hasAutoBrake"`
        IsCriticalWarning  bool      `json:"isCriticalWarning"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainReputationDrifts (
        domain_name TEXT PRIMARY KEY,
        current_score INTEGER NOT NULL DEFAULT 100,
        seven_day_drift_pct REAL NOT NULL DEFAULT 0.0,
        dampening_factor REAL NOT NULL DEFAULT 1.0,
        effective_hourly_tps REAL NOT NULL DEFAULT 50.0,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_dampener_active INTEGER NOT NULL DEFAULT 1,
        has_auto_brake INTEGER NOT NULL DEFAULT 1,
        is_critical_warning INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed under `cli/email/governor/reputation_drift.go`. Integrates telemetry feeds to throttle sender domains before spam blacklists trigger.

---

### Improvement 1045: Token Bucket Leaky Valve Dispatch Governor & Burst Limiter
- **ID:** `CONCUR-TOKEN-1045`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Regulate individual SMTP socket connection handshakes using a mathematical leaky-bucket token governor, guaranteeing strict spacing between commands to defeat ESP behavioral fingerprinting.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Animated ASCII bucket fill level with drip interval indicator in milliseconds.
    - Token generation rate ($r$) vs capacity ($b$) gauges.
    - Terminal UI:
      ```text
      ┌─ [TOKEN BUCKET LEAKY VALVE DISPATCH GOVERNOR] ─────────────────────────────────────┐
      │ BUCKET CAPACITY: 60 TOKENS | FILL RATE: 2.0 TOKENS/SEC | LEAK INTERVAL: 500ms       │
      ├────────────────────────────────────────────────────────────────────────────────────┤
      │ TOKEN RESERVOIR LEVEL                                                              │
      │ [████████████████████████████████████░░░░░░░░░░░░░░░░░░] 40 / 60 Tokens (66.7%)    │
      ├────────────────────────────────────────────────────────────────────────────────────┤
      │ PARAMETER           │ CONFIGURED VALUE │ REAL-TIME OBSERVED                        │
      ├─────────────────────┼──────────────────┼───────────────────────────────────────────┤
      │ Burst Allowance     │        10 tokens │ Maximum 8 tokens released in 1s           │
      │ Inter-Message Delay │           500 ms │ Mean: 512ms | Jitter: ±42ms               │
      │ Token Refill Speed  │       2.0 / sec  │ Nominal 120 msg/min                       │
      │ Dropped (Throttled) │              N/A │ 0 messages dropped (0.0%)                 │
      ├────────────────────────────────────────────────────────────────────────────────────┤
      │ [FLAGS] IsLeakyBucketActive: TRUE | HasBurstLimiter: TRUE | IsJitterSpaced: TRUE    │
      │ [KEYS]  (u) Increase Refill  (d) Decrease Refill  (b) Toggle Burst  (q) Exit       │
      └────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Canvas-free dynamic SVG bucket animation illustrating token drip physics.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor token-bucket`
  - Flags:
    - `--capacity <int>`: Maximum burst token capacity (default: `60`).
    - `--refill-rate <float>`: Tokens added per second (default: `2.0`).
    - `--is-jitter-spaced`: Inject micro-delay jitter between tokens (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TokenBucketConfig struct {
        BucketID           string    `json:"bucketId"`
        Capacity           int       `json:"capacity"`
        RefillRatePerSec   float64   `json:"refillRatePerSec"`
        CurrentTokens      float64   `json:"currentTokens"`
        LastRefilledAt     time.Time `json:"lastRefilledAt"`
        IsLeakyValveActive bool      `json:"isLeakyValveActive"`
        HasBurstLimiter    bool      `json:"hasBurstLimiter"`
        IsThrottled        bool      `json:"isThrottled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TokenBucketConfigs (
        bucket_id TEXT PRIMARY KEY,
        capacity INTEGER NOT NULL DEFAULT 60,
        refill_rate_per_sec REAL NOT NULL DEFAULT 2.0,
        current_tokens REAL NOT NULL DEFAULT 60.0,
        last_refilled_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_leaky_valve_active INTEGER NOT NULL DEFAULT 1,
        has_burst_limiter INTEGER NOT NULL DEFAULT 1,
        is_throttled INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/governor/token_bucket.go`. Regulates message release rate to match human typing and clicking profiles.

---

### Improvement 1046: Cross-Domain Warmup Affinity Pairing & Graph Topology Visualizer
- **ID:** `CONCUR-AFFINITY-1046`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Generate and enforce optimal cross-domain sending affinity pairings, ensuring sender domains exchange traffic across varied recipient domains without forming isolated reciprocal sending rings.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Interactive directed graph representation using UTF-8 arrows (`→`, `⇄`, `↷`).
    - Ring detection warning alerts: highlights illicit circular sending patterns (`A → B → C → A`) that trigger spam heuristics.
    - Terminal UI:
      ```text
      ┌─ [CROSS-DOMAIN AFFINITY GRAPH & RING PREVENTION TOPOLOGY] ──────────────────────────┐
      │ ACTIVE SENDER NODES: 24 | GRAPH DENSITY: 0.18 | RING TOPOLOGY WARNINGS: 0           │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ SOURCE DOMAIN       │ TARGET DOMAIN PAIRS                   │ RECIPROCITY │ HEALTH  │
      ├─────────────────────┼───────────────────────────────────────┼─────────────┼─────────┤
      │ outreach.domain-a   │ → seed-group-1, → ext-corp-b, → m365-c │ 8.4% (Low)  │ [●] SAFE│
      │ outbound.domain-b   │ → seed-group-4, → proton-d, → yahoo-e │ 5.1% (Low)  │ [●] SAFE│
      │ growth.domain-c     │ → seed-group-2, → gsuite-f, → corp-a  │ 6.2% (Low)  │ [●] SAFE│
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [TOPOLOGY INTEGRITY] Perfectly balanced acyclic sending flow. Zero self-loops.      │
      │ [FLAGS]              IsAcyclicEnforced: TRUE | HasRingDetection: TRUE | IsValid:TRUE│
      │ [KEYS]               (r) Recompute Graph  (t) Test Topology  (q) Exit               │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Dynamic SVG node-link graph with interactive node dragging and edge weight highlighting.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor affinity-graph`
  - Flags:
    - `--max-reciprocity-pct <float>`: Maximum allowed reciprocal volume between two domains (default: `10.0`).
    - `--is-acyclic-enforced`: Block dispatch plans containing closed topological cycles (default: `true`).
    - `--has-auto-rebalance`: Auto-shuffle target seeds to break emerging rings (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainAffinityEdge struct {
        SourceDomain      string    `json:"sourceDomain"`
        TargetDomain      string    `json:"targetDomain"`
        VolumeWeight      int       `json:"volumeWeight"`
        ReciprocityScore  float64   `json:"reciprocityScore"`
        LastDispatchedAt  time.Time `json:"lastDispatchedAt"`
        IsAcyclicEnforced bool      `json:"isAcyclicEnforced"`
        HasRingWarning    bool      `json:"hasRingWarning"`
        IsValidTopology   bool      `json:"isValidTopology"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainAffinityEdges (
        source_domain TEXT NOT NULL,
        target_domain TEXT NOT NULL,
        volume_weight INTEGER NOT NULL DEFAULT 0,
        reciprocity_score REAL NOT NULL DEFAULT 0.0,
        last_dispatched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_acyclic_enforced INTEGER NOT NULL DEFAULT 1,
        has_ring_warning INTEGER NOT NULL DEFAULT 0,
        is_valid_topology INTEGER NOT NULL DEFAULT 1,
        PRIMARY KEY (source_domain, target_domain)
    );
    CREATE INDEX IF NOT EXISTS idx_affinity_recip ON DomainAffinityEdges(reciprocity_score);
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `cli/email/governor/affinity_graph.go`. Prevents the formation of detectable synthetic peer networks among managed domains.

---

### Improvement 1047: MX Provider Graylist Detection & Early Retry Accommodator
- **ID:** `CONCUR-GRAY-1047`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Distinguish transient graylisting delays from true deliverability blocks, automatically caching graylist challenge intervals per receiving MX server to optimize subsequent dispatch timings.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live feed of MX graylisting challenges and cached delay windows.
    - Status badges: `[GRAYLIST_CHALLENGE]` in sky blue, `[WINDOW_PASSED]` in emerald, `[CHALLENGE_ACCEPTED]` in violet.
    - Terminal UI:
      ```text
      ┌─ [MX GRAYLIST DETECTION & RETRY ACCOMMODATOR] ──────────────────────────────────────┐
      │ CACHED MX WINDOWS: 84 | ACTIVE CHALLENGES: 5 | SUCCESSFUL REPASSES: 99.4%           │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ TARGET MX HOST          │ INITIAL REJECTION CODE  │ DETECTED WINDOW │ NEXT RETRY    │
      ├─────────────────────────┼─────────────────────────┼─────────────────┼───────────────┤
      │ mx1.puremessage.com     │ 451 4.7.1 Greylisted    │ 15 minutes      │ in 4m 12s     │
      │ mail.mimecast.com       │ 451 Internal deferral   │ 10 minutes      │ in 1m 05s     │
      │ mx.spamexperts.com      │ 450 4.3.2 Please retry  │ 5 minutes       │ READY         │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [INTELLIGENCE] MX host mx.spamexperts.com verified. Delivery will succeed on retry. │
      │ [FLAGS]        IsGraylistLearned: TRUE | HasAutoSchedule: TRUE | IsStrictCache: TRUE│
      │ [KEYS]         (f) Flush Expired  (s) Force Send Ready  (q) Exit                     │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Clean tabular view of active graylist windows with live countdown timers.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor graylist-status`
  - Flags:
    - `--mx-host <string>`: Target mail exchange host (optional).
    - `--is-auto-reschedule`: Move graylisted items to precise release window (default: `true`).
    - `--has-learning-cache`: Persist learned delays to `pipeline.db` (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GraylistWindowEntry struct {
        MXHost             string    `json:"mxHost"`
        ObservedDelaySec   int       `json:"observedDelaySec"`
        SuccessCount       int       `json:"successCount"`
        LastChallengeAt    time.Time `json:"lastChallengeAt"`
        IsGraylistLearned  bool      `json:"isGraylistLearned"`
        HasAutoSchedule    bool      `json:"hasAutoSchedule"`
        IsActive           bool      `json:"isActive"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GraylistWindowEntries (
        mx_host TEXT PRIMARY KEY,
        observed_delay_sec INTEGER NOT NULL DEFAULT 600,
        success_count INTEGER NOT NULL DEFAULT 1,
        last_challenge_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_graylist_learned INTEGER NOT NULL DEFAULT 1,
        has_auto_schedule INTEGER NOT NULL DEFAULT 1,
        is_active INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/email/governor/graylist_accommodator.go`. Schedules queue release precisely when remote graylisting timeouts expire.

---

### Improvement 1048: SMTP Session Multiplexing & TLS Keep-Alive Connection Pool HUD
- **ID:** `CONCUR-SMTP-1048`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Maintain persistent, multiplexed SMTP/STARTTLS connections to major recipient MTAs, minimizing TLS handshake CPU overhead and DNS lookup thrashing while adhering to max-messages-per-connection limits.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Connection pool table detailing socket status, cipher suites, TLS version (TLS 1.3), message counter per socket, and idle expiration timers.
    - Terminal UI:
      ```text
      ┌─ [SMTP SESSION MULTIPLEXING & TLS CONNECTION POOL] ─────────────────────────────────┐
      │ POOL STATUS: 32 ACTIVE SOCKETS | 18 IDLE (KEEP-ALIVE) | HANDSHAKES SAVED: 1,420     │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ DESTINATION HOST        │ PROTOCOL │ CIPHER             │ MSGS/CONN │ IDLE REMAIN   │
      ├─────────────────────────┼──────────┼────────────────────┼───────────┼───────────────┤
      │ aspmx.l.google.com:25   │ TLS 1.3  │ TLS_AES_256_GCM    │ 4 / 10    │ 24s           │
      │ mail-protection.ms:25   │ TLS 1.3  │ TLS_CHACHA20_POLY  │ 8 / 10    │ 12s           │
      │ mta5.am0.yahoodns.net:25│ TLS 1.2  │ ECDHE_RSA_AES_128  │ 2 / 5     │ 45s           │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [PERFORMANCE] Mean Handshake Latency: 42ms | Handshake Overhead Reduction: 78.4%   │
      │ [FLAGS]       IsMultiplexActive: TRUE | HasTlsSessionReuse: TRUE | IsPristine: TRUE │
      │ [KEYS]        (k) Kill Idle Sockets  (t) Tune Max Per Conn  (q) Exit                 │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live socket lifetime visualizer rendered via plain SVG timelines.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor session-pool`
  - Flags:
    - `--max-idle-sec <int>`: Socket idle timeout before closing (default: `60`).
    - `--max-msgs-per-conn <int>`: Safety ceiling for messages per single socket (default: `10`).
    - `--is-tls-reuse`: Enable TLS Session Ticket resumption (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SMTPSessionSocket struct {
        SocketID           string    `json:"socketId"`
        RemoteHost         string    `json:"remoteHost"`
        TLSVersion         string    `json:"tlsVersion"`
        CipherSuite        string    `json:"cipherSuite"`
        MessagesDelivered  int       `json:"messagesDelivered"`
        MaxAllowedMessages int       `json:"maxAllowedMessages"`
        ConnectedAt        time.Time `json:"connectedAt"`
        IsMultiplexActive  bool      `json:"isMultiplexActive"`
        HasTlsReuse        bool      `json:"hasTlsReuse"`
        IsIdle             bool      `json:"isIdle"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SMTPSessionSockets (
        socket_id TEXT PRIMARY KEY,
        remote_host TEXT NOT NULL,
        tls_version TEXT NOT NULL,
        cipher_suite TEXT NOT NULL,
        messages_delivered INTEGER NOT NULL DEFAULT 0,
        max_allowed_messages INTEGER NOT NULL DEFAULT 10,
        connected_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_multiplex_active INTEGER NOT NULL DEFAULT 1,
        has_tls_reuse INTEGER NOT NULL DEFAULT 1,
        is_idle INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/governor/session_pool.go`. Accelerates dispatch pipeline throughput and eliminates TLS handshake latency penalties.

---

### Improvement 1049: Real-Time MX Deferral Rate Trend & Spike Alerting HUD
- **ID:** `CONCUR-DEFER-1049`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Track rolling 5-minute, 1-hour, and 24-hour deferral rate trends across recipient domains, firing real-time terminal alarms when deferrals spike above statistical standard deviation bounds.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual ASCII sparklines showing historical deferral vs success counts.
    - Audio/visual flashing banner in Lipgloss red (`#dc2626`) upon threshold breach.
    - Terminal UI:
      ```text
      ┌─ [REAL-TIME MX DEFERRAL TREND & ANOMALY SPIKE ALERT] ───────────────────────────────┐
      │ ROLLING DEFERRAL RATE: 1.8% (BASELINE: 1.2%) | SPIKE ALERTS ACTIVE: 0                │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ ROLLING 60-MINUTE DEFERRAL TREND (SPARKLINE)                                        │
      │ 5% │                                                ╭──╮                            │
      │ 2% │ ───────────────╮                 ╭─────────────╯  ╰──────────                  │
      │ 0% ┼────────────────┴─────────────────┴───────────────────────────                  │
      │    0m               15m              30m              45m        60m                │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ PROVIDER          │ 5M RATE │ 1H RATE │ DEVIATION (σ) │ THREAT ASSESSMENT           │
      ├───────────────────┼─────────┼─────────┼───────────────┼─────────────────────────────┤
      │ Google Workspace  │    0.8% │    0.9% │        +0.2σ  │ [●] NOMINAL                 │
      │ Microsoft 365     │    2.4% │    2.1% │        +0.8σ  │ [●] NOMINAL                 │
      │ Yahoo Mail        │    6.8% │    3.2% │        +2.6σ  │ [▲] ELEVATED SPIKE          │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [FLAGS] IsAnomalyDetectorActive: TRUE | HasSmsAlert: FALSE | IsAutoMuted: FALSE     │
      │ [KEYS]  (a) Acknowledge Warning  (t) Adjust Sensitivity  (q) Exit                   │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live canvas-free SVG area chart with anomaly band shading.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor deferral-trend`
  - Flags:
    - `--sigma-threshold <float>`: Standard deviation threshold to trigger alert (default: `2.5`).
    - `--is-auto-mute`: Automatically silence alerts during known maintenance windows (default: `false`).
    - `--has-sound-alert`: Play terminal bell on critical spike (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DeferralTrendEntry struct {
        ProviderName            string    `json:"providerName"`
        FiveMinRate             float64   `json:"fiveMinRate"`
        OneHourRate             float64   `json:"oneHourRate"`
        SigmaDeviation          float64   `json:"sigmaDeviation"`
        ObservedAt              time.Time `json:"observedAt"`
        IsAnomalyDetectorActive bool      `json:"isAnomalyDetectorActive"`
        HasSpikeAlert           bool      `json:"hasSpikeAlert"`
        IsUnderInvestigation    bool      `json:"isUnderInvestigation"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DeferralTrendEntries (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        provider_name TEXT NOT NULL,
        five_min_rate REAL NOT NULL DEFAULT 0.0,
        one_hour_rate REAL NOT NULL DEFAULT 0.0,
        sigma_deviation REAL NOT NULL DEFAULT 0.0,
        observed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_anomaly_detector_active INTEGER NOT NULL DEFAULT 1,
        has_spike_alert INTEGER NOT NULL DEFAULT 0,
        is_under_investigation INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_deferral_prov ON DeferralTrendEntries(provider_name, observed_at DESC);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `cli/email/governor/deferral_trend.go`. Notifies operators early before temporary throttling escalates into hard IP blocks.

---

### Improvement 1050: Multi-Tenant Quota Distribution & Priority Queue Orchestrator
- **ID:** `CONCUR-PRIO-1050`
- **Cluster:** 105 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Allocate sender mailbox bandwidth dynamically between high-priority transactional emails, scheduled marketing sequences, and warmup traffic, preventing bulk campaigns from starving critical validation messages.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-queue priority lane monitor with real-time throughput metrics per tier (`P0: Transactional`, `P1: Warmup Seed`, `P2: Scheduled Campaign`).
    - Starvation prevention timer ensuring lower-priority queues make steady progress.
    - Terminal UI:
      ```text
      ┌─ [MULTI-TENANT QUOTA ALLOCATION & PRIORITY QUEUE ORCHESTRATOR] ─────────────────────┐
      │ GLOBAL THROUGHPUT: 12.4 MSG/S | DISPATCH STRATEGY: WEIGHTED FAIR QUEUING (WFQ)      │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ PRIORITY LANE   │ WEIGHT │ QUEUE DEPTH │ ALLOCATED TPS │ STARVATION STATUS          │
      ├─────────────────┼────────┼─────────────┼───────────────┼────────────────────────────┤
      │ P0 (Transact)   │    70% │           0 │    8.7 msg/s  │ [●] ZERO DELAY (0ms wait)  │
      │ P1 (Warmup)     │    20% │          42 │    2.5 msg/s  │ [●] NOMINAL (1.2s wait)    │
      │ P2 (Campaign)   │    10% │       1,240 │    1.2 msg/s  │ [●] CONTROLLED (42s wait)  │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [WFQ FAIRNESS] Anti-starvation active: Guaranteeing min 1 msg/s to Campaign tier.   │
      │ [FLAGS]        IsPriorityActive: TRUE | HasAntiStarve: TRUE | IsStrictWeighted: TRUE│
      │ [KEYS]         (w) Adjust Weights  (p) Pause Campaign Tier  (q) Exit                │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live lane throughput distribution visualization using CSS flex layouts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email governor queue-prio`
  - Flags:
    - `--transact-weight <int>`: Weight percentage for transactional tier (default: `70`).
    - `--warmup-weight <int>`: Weight percentage for warmup tier (default: `20`).
    - `--campaign-weight <int>`: Weight percentage for campaign tier (default: `10`).
    - `--is-anti-starve`: Force execution of lower priority items after timeout (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PriorityQueueMetric struct {
        PriorityTier       string    `json:"priorityTier"`
        WeightPercent      int       `json:"weightPercent"`
        CurrentQueueDepth  int       `json:"currentQueueDepth"`
        AllocatedTPS       float64   `json:"allocatedTps"`
        MaxLatencySeconds  float64   `json:"maxLatencySeconds"`
        LastProcessedAt    time.Time `json:"lastProcessedAt"`
        IsPriorityActive   bool      `json:"isPriorityActive"`
        HasAntiStarve      bool      `json:"hasAntiStarve"`
        IsStarvationSafe   bool      `json:"isStarvationSafe"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PriorityQueueMetrics (
        priority_tier TEXT PRIMARY KEY,
        weight_percent INTEGER NOT NULL DEFAULT 33,
        current_queue_depth INTEGER NOT NULL DEFAULT 0,
        allocated_tps REAL NOT NULL DEFAULT 0.0,
        max_latency_seconds REAL NOT NULL DEFAULT 0.0,
        last_processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_priority_active INTEGER NOT NULL DEFAULT 1,
        has_anti_starve INTEGER NOT NULL DEFAULT 1,
        is_starvation_safe INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/governor/priority_orchestrator.go`. Guarantees mission-critical emails route with zero latency regardless of bulk background load.

---

# Cluster 106: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 1051–1060)

---

### Improvement 1051: Zero-CDN Responsive Email HTML Sandbox & Render Preview HUD
- **ID:** `EMAIL-SANDBOX-1051`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Provide a local, zero-external-CDN HTML sanitizer and responsive email sandbox that previews HTML/plain-text emails across desktop, tablet, and mobile viewport breakpoints without leaking external tracker assets.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Side-by-side terminal rendering: Plain-text fallback vs. stripped HTML DOM tree preview with CSS property extraction.
    - Viewport switcher: Toggle between 320px (Mobile), 600px (Tablet/Standard Email Width), and 1024px (Desktop).
    - Terminal UI:
      ```text
      ┌─ [ZERO-CDN RESPONSIVE EMAIL SANDBOX: VIEWPORT 600px] ───────────────────────────────┐
      │ TEMPLATE: template-intro-v3.html | SECURITY GRADE: PRISTINE (NO EXTERNAL LEAKS)      │
      ├───────────────────────────────────┬─────────────────────────────────────────────────┤
      │ DOM STRUCTURE & STYLES            │ RENDERED PLAIN-TEXT FALLBACK PREVIEW            │
      ├───────────────────────────────────┼─────────────────────────────────────────────────┤
      │ <table> [width: 100%, max: 600px] │ Alex,                                           │
      │  └── <tr>                         │                                                 │
      │       └── <td> [padding: 24px]    │ Following up on our benchmark discussion from   │
      │            ├── <h1> [font: 22px]  │ yesterday. Attached is the GitMap architecture  │
      │            ├── <p>  [font: 14px]  │ specification for review.                       │
      │            └── <a>  [btn: #2563eb]│                                                 │
      │                                   │ [View Full Architecture Spec]                   │
      │                                   │ https://gitmap.local/spec/347                   │
      ├───────────────────────────────────┴─────────────────────────────────────────────────┤
      │ [COMPATIBILITY] Outlook Desktop: 100% | Gmail Mobile: 100% | Apple Mail: 100%       │
      │ [KEYS] (1) Mobile (320px)  (2) Tablet (600px)  (3) Desktop  (s) Security Scan (q)Back │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Sandboxed iframe (`sandbox="allow-same-origin"`) with zero remote script execution, stripping all tracking pixels and remote web fonts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox preview`
  - Flags:
    - `--file <path>`: Relative path to HTML email template (required).
    - `--width <int>`: Virtual viewport width in pixels (default: `600`).
    - `--is-plain-text-preview`: Display plain-text extraction side-by-side (default: `true`).
    - `--has-sandbox-isolation`: Strictly forbid remote asset fetching (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailSandboxReport struct {
        TemplatePath        string    `json:"templatePath"`
        HTMLPayloadSize     int       `json:"htmlPayloadSize"`
        PlainTextLength     int       `json:"plainTextLength"`
        ResponsiveBreakpts  []int     `json:"responsiveBreakpts"`
        ParsedAt            time.Time `json:"parsedAt"`
        IsSandboxIsolated   bool      `json:"isSandboxIsolated"`
        HasPlainTextValid   bool      `json:"hasPlainTextValid"`
        IsCompatibleOutlook bool      `json:"isCompatibleOutlook"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailSandboxReports (
        template_path TEXT PRIMARY KEY,
        html_payload_size INTEGER NOT NULL DEFAULT 0,
        plain_text_length INTEGER NOT NULL DEFAULT 0,
        responsive_breakpts TEXT NOT NULL DEFAULT '320,600',
        parsed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_sandbox_isolated INTEGER NOT NULL DEFAULT 1,
        has_plain_text_valid INTEGER NOT NULL DEFAULT 1,
        is_compatible_outlook INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed under `cli/email/sandbox/html_preview.go`. Provides developers with an offline preview environment free from tracking leaks.

---

### Improvement 1052: Autonomous Dark Mode Inversion & Color Contrast Auditor
- **ID:** `EMAIL-DARK-1052`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Simulate and evaluate how email clients (Apple Mail, Outlook Dark Mode, Gmail App) invert custom CSS colors, alerting if inverted contrast ratios violate WCAG AAA accessibility or produce unreadable black-on-dark text.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual color swatch comparison: Original Light Mode vs. Simulated Dark Mode Inversion.
    - WCAG Contrast ratio indicator (`4.5:1` AA minimum, `7.0:1` AAA target).
    - Terminal UI:
      ```text
      ┌─ [DARK MODE INVERSION & CONTRAST AUDITOR] ──────────────────────────────────────────┐
      │ TEMPLATE: template-onboarding.html | INVERSION PROFILE: GMAIL_ANDROID_FULL_INVERT   │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ ELEMENT    │ LIGHT MODE VALUE │ INVERTED COLOR   │ CONTRAST RATIO │ STATUS          │
      ├────────────┼──────────────────┼──────────────────┼────────────────┼─────────────────┤
      │ Body BG    │ #ffffff (White)  │ #121212 (Dark)   │ 16.8:1 (AAA)   │ [●] PASS        │
      │ Card BG    │ #f8fafc (Slate)  │ #1e293b (Slate9) │ 14.2:1 (AAA)   │ [●] PASS        │
      │ Primary Txt│ #0f172a (Dark)   │ #f1f5f9 (Light)  │ 15.1:1 (AAA)   │ [●] PASS        │
      │ Link Anchor│ #2563eb (Blue6)  │ #1d4ed8 (Blue7)  │  3.8:1 (FAIL)  │ [▲] LOW CONTRAST│
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [REMEDIATION] Add `@media (prefers-color-scheme: dark)` override for link colors.   │
      │ [FLAGS]       IsDarkCompliant: FALSE | HasMetaSchemeTag: TRUE | IsFixSuggested: TRUE│
      │ [KEYS]        (f) Auto-Fix Stylesheet  (c) Change Inversion Profile  (q) Exit       │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Split-screen comparison slider displaying real-time dark mode simulation.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox dark-mode`
  - Flags:
    - `--file <path>`: Email template path (required).
    - `--min-contrast-ratio <float>`: Target WCAG contrast threshold (default: `4.5`).
    - `--is-auto-fix`: Inject dark-mode CSS overrides into `<style>` block (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeAuditResult struct {
        TemplateFile        string    `json:"templateFile"`
        InversionEngine     string    `json:"inversionEngine"`
        LowestContrastRatio float64   `json:"lowestContrastRatio"`
        FailedElementsCount int       `json:"failedElementsCount"`
        AuditedAt           time.Time `json:"auditedAt"`
        IsDarkCompliant     bool      `json:"isDarkCompliant"`
        HasMetaSchemeTag    bool      `json:"hasMetaSchemeTag"`
        IsAutoRemediated    bool      `json:"isAutoRemediated"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DarkModeAuditResults (
        template_file TEXT PRIMARY KEY,
        inversion_engine TEXT NOT NULL,
        lowest_contrast_ratio REAL NOT NULL DEFAULT 7.0,
        failed_elements_count INTEGER NOT NULL DEFAULT 0,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_dark_compliant INTEGER NOT NULL DEFAULT 1,
        has_meta_scheme_tag INTEGER NOT NULL DEFAULT 1,
        is_auto_remediated INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/sandbox/dark_mode_audit.go`. Prevents illegible layouts in dark-themed email apps.

---

### Improvement 1053: Spam Trigger Lexical & Semantic Scanner HUD
- **ID:** `EMAIL-SPAMLEX-1053`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan subject lines, preview text, and email body content for known spam trigger words, high-risk financial tokens, deceptive urgency patterns, and punctuation anomalies before dispatch.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Interactive text buffer with identified trigger terms underlined in red Lipgloss styling.
    - Risk category breakdown: `Financial/Free`, `Deceptive Urgency`, `Punctuation/Caps Excess`.
    - Terminal UI:
      ```text
      ┌─ [SPAM TRIGGER LEXICAL & SEMANTIC SCANNER] ─────────────────────────────────────────┐
      │ SCANNED FILE: campaign-v1.html | TOTAL WORDS: 142 | RISK SCORE: 28/100 [MILD RISK] │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ LINE │ TRIGGER PHRASE          │ CATEGORY           │ RECOMMENDED ALTERNATIVE       │
      ├──────┼─────────────────────────┼────────────────────┼───────────────────────────────┤
      │   1  │ "ACT NOW: Exclusive..." │ Deceptive Urgency  │ "Update: New capabilities..." │
      │   8  │ "100% Risk Free Trial"  │ Financial Guarantee│ "Explore without commitment"  │
      │  14  │ "$$$ Save Big $$$"      │ Symbol / Token Bait│ "Optimize operational budget" │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [SEMANTIC DENSITY] Trigger token density: 2.1% (Safe ceiling is < 1.0%).            │
      │ [FLAGS]            IsSpamClean: FALSE | HasDeceptiveSubject: TRUE | IsFixReady: TRUE│
      │ [KEYS]             (r) Auto-Replace with Safe Synonyms  (e) Edit Source  (q) Exit   │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline HTML editor with real-time word squiggles and replacement chips.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox spam-scanner`
  - Flags:
    - `--file <path>`: Template file to inspect (required).
    - `--max-risk-score <int>`: Rejection threshold score (default: `25`).
    - `--is-auto-replace`: Automatically swap identified words for safe synonyms (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamLexicalScanResult struct {
        DocumentPath         string    `json:"documentPath"`
        TotalWordCount       int       `json:"totalWordCount"`
        TriggerWordCount     int       `json:"triggerWordCount"`
        CalculatedRiskScore  int       `json:"calculatedRiskScore"`
        ScannedAt            time.Time `json:"scannedAt"`
        IsSpamClean          bool      `json:"isSpamClean"`
        HasDeceptiveSubject  bool      `json:"hasDeceptiveSubject"`
        IsReplacementApplied bool      `json:"isReplacementApplied"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamLexicalScanResults (
        document_path TEXT PRIMARY KEY,
        total_word_count INTEGER NOT NULL DEFAULT 0,
        trigger_word_count INTEGER NOT NULL DEFAULT 0,
        calculated_risk_score INTEGER NOT NULL DEFAULT 0,
        scanned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_spam_clean INTEGER NOT NULL DEFAULT 1,
        has_deceptive_subject INTEGER NOT NULL DEFAULT 0,
        is_replacement_applied INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed under `cli/email/sandbox/spam_scanner.go`. Guards copywriters against accidental usage of words flagged by enterprise Bayes filters.

---

### Improvement 1054: Zero-Storage Tracking Pixel Sanitizer & Image Proxy Pre-Fetcher
- **ID:** `EMAIL-PIXEL-1054`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Strip external tracking web bugs ($1\times1$ GIFs, hidden spy pixels) from inbound messages and route necessary remote content through a privacy-preserving zero-storage caching proxy before preview.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Table displaying all remote images found in the HTML payload, flagging tracking pixels and extracting dimensions.
    - Terminal UI:
      ```text
      ┌─ [TRACKING PIXEL SANITIZER & REMOTE ASSET AUDITOR] ─────────────────────────────────┐
      │ ASSETS DETECTED: 4 | TRACKERS STRIPPED: 2 | EMBEDDED CONTENT: CLEAN                  │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ ASSET URL                          │ DIMENSIONS │ TYPE       │ ACTION TAKEN         │
      ├────────────────────────────────────┼────────────┼────────────┼──────────────────────┤
      │ https://spy.track.io/open?uid=9081 │ 1x1 px     │ Web Bug    │ [STRIPPED / BLOCKED] │
      │ https://analytics.corp.com/px.gif  │ 0x0 px     │ Hidden Gif │ [STRIPPED / BLOCKED] │
      │ https://assets.cdn.com/logo.png    │ 240x60 px  │ Valid Hero │ [PROXIED / SAFE]     │
      │ cid:attachment-header.jpg          │ 600x200 px │ Inline CID │ [PASSED DIRECT]      │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [PRIVACY GUARANTEE] Sender IP and User-Agent leak eliminated. Zero CDN tracking.    │
      │ [FLAGS]             IsSanitized: TRUE | HasTrackersStripped: TRUE | IsSecure: TRUE  │
      │ [KEYS]              (d) Dump Sanitized HTML  (v) View Assets  (q) Exit              │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Side-by-side sanitization diff highlighting removed tracking code.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox sanitize-pixels`
  - Flags:
    - `--input <path>`: Source email HTML file (required).
    - `--output <path>`: Sanitized destination HTML file (optional).
    - `--is-aggressive-filter`: Strip all 1x1, 0x0, and hidden image tags (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PixelSanitizeReport struct {
        TemplateFile        string    `json:"templateFile"`
        TotalImagesFound    int       `json:"totalImagesFound"`
        TrackersStripped    int       `json:"trackersStripped"`
        SanitizedAt         time.Time `json:"sanitizedAt"`
        IsSanitized         bool      `json:"isSanitized"`
        HasTrackersStripped bool      `json:"hasTrackersStripped"`
        IsPrivacyClean      bool      `json:"isPrivacyClean"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PixelSanitizeReports (
        template_file TEXT PRIMARY KEY,
        total_images_found INTEGER NOT NULL DEFAULT 0,
        trackers_stripped INTEGER NOT NULL DEFAULT 0,
        sanitized_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_sanitized INTEGER NOT NULL DEFAULT 1,
        has_trackers_stripped INTEGER NOT NULL DEFAULT 0,
        is_privacy_clean INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/email/sandbox/pixel_sanitizer.go`. Protects user privacy while validating incoming and outgoing email content.

---

### Improvement 1055: CSS Inliner & Email Client Compatibility Linter HUD
- **ID:** `EMAIL-INLINER-1055`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically inline external and `<style>` block CSS rules directly into HTML `style=""` attributes while warning of unsupported CSS properties across legacy Outlook, Gmail, and Yahoo webmail engines.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Real-time compatibility checklist showing support percentage across major client engines (`Outlook 2016-2021`, `Gmail App iOS/Android`, `Apple Mail`, `Yahoo Mail Web`).
    - Unsupported CSS feature alerts (`display: flex`, `position: absolute`, `calc()`).
    - Terminal UI:
      ```text
      ┌─ [CSS INLINER & CLIENT COMPATIBILITY LINTER] ───────────────────────────────────────┐
      │ SOURCE: newsletter.html | INLINED NODES: 84 | OVERALL COMPATIBILITY: 98.2%          │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ CSS PROPERTY        │ TARGET COMPONENT     │ INLINE STATUS │ CLIENT WARNING         │
      ├─────────────────────┼──────────────────────┼───────────────┼────────────────────────┤
      │ display: flex       │ <div class="col-2">  │ REMAPPED      │ Converted to <table>   │
      │ border-radius: 8px  │ <a class="button">   │ INLINED       │ Ignored by Outlook Win │
      │ max-width: 600px    │ <table class="main"> │ INLINED       │ Requires VML wrapper   │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [AUTO-CORRECTION] Injected MSO condition comments for Outlook desktop fidelity.     │
      │ [FLAGS]           IsInlined: TRUE | HasMsoWrappers: TRUE | IsStrictCompliant: TRUE  │
      │ [KEYS]            (e) Export Inlined HTML  (v) View Diff  (q) Exit                   │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Tabular CSS rule matrix with green checkmarks and yellow warning chips.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox inline-css`
  - Flags:
    - `--file <path>`: Input HTML template file (required).
    - `--has-mso-fallbacks`: Generate conditional VML tags for Microsoft Outlook (default: `true`).
    - `--is-minify-enabled`: Strip extraneous whitespace after inlining (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CSSInlinerReport struct {
        SourcePath          string    `json:"sourcePath"`
        InlinedRuleCount    int       `json:"inlinedRuleCount"`
        CompatibilityRating float64   `json:"compatibilityRating"`
        InlinedAt           time.Time `json:"inlinedAt"`
        IsInlined           bool      `json:"isInlined"`
        HasMsoFallbacks     bool      `json:"hasMsoFallbacks"`
        IsOutlookTested     bool      `json:"isOutlookTested"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CSSInlinerReports (
        source_path TEXT PRIMARY KEY,
        inlined_rule_count INTEGER NOT NULL DEFAULT 0,
        compatibility_rating REAL NOT NULL DEFAULT 100.0,
        inlined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_inlined INTEGER NOT NULL DEFAULT 1,
        has_mso_fallbacks INTEGER NOT NULL DEFAULT 1,
        is_outlook_tested INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `cli/email/sandbox/css_inliner.go`. Guarantees pixel-perfect rendering across fragmented desktop and web email clients.

---

### Improvement 1056: Link Reputation & Blacklist Target Verification Scanner
- **ID:** `EMAIL-LINK-1056`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Extract every anchor link, redirect domain, and image host inside outbound templates, checking each destination against Google Safe Browsing, SURBL, and URIBL blacklists.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Interactive link list with destination reputation flags and redirect resolution paths.
    - Danger highlighting: Red `[BLACKLISTED]` tag with immediate terminal alert beep.
    - Terminal UI:
      ```text
      ┌─ [LINK REPUTATION & DOMAIN BLACKLIST VERIFIER] ─────────────────────────────────────┐
      │ AUDITED HYPERLINKS: 6 | THREATS DETECTED: 0 | RESOLUTION SPEED: 180ms                │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ ANCHOR TEXT         │ FINAL DESTINATION URL        │ SURBL │ URIBL │ STATUS         │
      ├─────────────────────┼──────────────────────────────┼───────┼───────┼────────────────┤
      │ "Documentation"     │ https://docs.gitmap.local/   │ CLEAN │ CLEAN │ [●] VERIFIED   │
      │ "Schedule Meeting"  │ https://cal.partner.io/alex  │ CLEAN │ CLEAN │ [●] VERIFIED   │
      │ "Unsubscribe"       │ https://outreach.domain/opt  │ CLEAN │ CLEAN │ [●] VERIFIED   │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [VERDICT] All outbound destination domains verified clean. Zero DNS traps found.    │
      │ [FLAGS]   IsScanClean: TRUE | HasShortenersFound: FALSE | IsVerified: TRUE          │
      │ [KEYS]    (r) Re-Verify Links  (s) Expand Redirects  (q) Exit                       │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live URL tree view showing intermediate HTTP 301/302 hops and final target DNS status.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox check-links`
  - Flags:
    - `--file <path>`: Template file to analyze (required).
    - `--is-expand-redirects`: Follow HTTP 301/302 redirects to destination (default: `true`).
    - `--has-shortener-ban`: Reject generic shorteners (bit.ly, t.co) as high spam risks (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LinkScanEntry struct {
        AnchorText        string    `json:"anchorText"`
        TargetURL         string    `json:"targetUrl"`
        FinalResolvedURL  string    `json:"finalResolvedUrl"`
        SurblStatus       string    `json:"surblStatus"`
        UriblStatus       string    `json:"uriblStatus"`
        ScannedAt         time.Time `json:"scannedAt"`
        IsScanClean       bool      `json:"isScanClean"`
        HasShortenerFound bool      `json:"hasShortenerFound"`
        IsValidLink       bool      `json:"isValidLink"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS LinkScanEntries (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        anchor_text TEXT NOT NULL,
        target_url TEXT NOT NULL,
        final_resolved_url TEXT NOT NULL,
        surbl_status TEXT NOT NULL DEFAULT 'CLEAN',
        uribl_status TEXT NOT NULL DEFAULT 'CLEAN',
        scanned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_scan_clean INTEGER NOT NULL DEFAULT 1,
        has_shortener_found INTEGER NOT NULL DEFAULT 0,
        is_valid_link INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_link_url ON LinkScanEntries(target_url);
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed under `cli/email/sandbox/link_verifier.go`. Prevents emails from tripping spam filters due to compromised destination websites.

---

### Improvement 1057: Email Attachment Payload Sandbox & Anti-Malware Validator
- **ID:** `EMAIL-ATTACH-1057`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify that email attachments (PDFs, ICS calendars, CSV data) comply with RFC sizing rules, do not contain embedded JavaScript or macros, and stay within strict size limits (< 10MB total).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Attachment inspection view showing file MIME types, byte sizes, magic byte validation, and macro detection results.
    - Terminal UI:
      ```text
      ┌─ [EMAIL ATTACHMENT SECURITY SANDBOX & VALIDATOR] ───────────────────────────────────┐
      │ INSPECTED FILES: 2 | TOTAL PAYLOAD: 2.4 MB / 10.0 MB | THREAT SCAN: 100% CLEAN      │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ FILE NAME       │ MIME TYPE       │ FILE SIZE │ MACROS / JS   │ VALIDATION VERDICT  │
      ├─────────────────┼─────────────────┼───────────┼───────────────┼─────────────────────┤
      │ proposal.pdf    │ application/pdf │  2.1 MB   │ NONE (Pure)   │ [●] PASSED (Safe)   │
      │ calendar.ics    │ text/calendar   │   12 KB   │ NONE (RFC5545)│ [●] PASSED (Safe)   │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [SECURITY SUMMARY] Magic numbers verified. No embedded executable streams found.    │
      │ [FLAGS]            IsPayloadSafe: TRUE | HasExecutableByte: FALSE | IsWithinCap:TRUE│
      │ [KEYS]             (s) Deep Binary Scan  (d) Detach File  (q) Exit                   │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Interactive file card list with security badges and hex-view preview.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox audit-attachment`
  - Flags:
    - `--file <path>`: Attachment file path (required).
    - `--max-size-mb <int>`: Upper size boundary threshold (default: `10`).
    - `--is-strict-magic-check`: Verify MIME type against internal byte header (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type AttachmentAuditEntry struct {
        AttachmentName    string    `json:"attachmentName"`
        MimeType          string    `json:"mimeType"`
        SizeBytes         int64     `json:"sizeBytes"`
        HasEmbeddedScript bool      `json:"hasEmbeddedScript"`
        AuditedAt         time.Time `json:"auditedAt"`
        IsPayloadSafe     bool      `json:"isPayloadSafe"`
        HasValidMagicByte bool      `json:"hasValidMagicByte"`
        IsWithinCap       bool      `json:"isWithinCap"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AttachmentAuditEntries (
        attachment_name TEXT PRIMARY KEY,
        mime_type TEXT NOT NULL,
        size_bytes INTEGER NOT NULL DEFAULT 0,
        has_embedded_script INTEGER NOT NULL DEFAULT 0,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_payload_safe INTEGER NOT NULL DEFAULT 1,
        has_valid_magic_byte INTEGER NOT NULL DEFAULT 1,
        is_within_cap INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `cli/email/sandbox/attachment_auditor.go`. Guarantees outbound attachments do not trip corporate gateway anti-malware traps.

---

### Improvement 1058: One-Click List-Unsubscribe Header & RFC 8058 Compliance HUD
- **ID:** `EMAIL-UNSUB-1058`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify and inject mandatory RFC 8058 `List-Unsubscribe` and `List-Unsubscribe-Post` headers into outbound email configurations, ensuring compliance with Google and Yahoo 2024+ sender requirements.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Compliance checklist indicating mailto URI formatting, HTTPS POST endpoint validity, and one-click token verification.
    - Terminal UI:
      ```text
      ┌─ [ONE-CLICK LIST-UNSUBSCRIBE & RFC 8058 COMPLIANCE AUDITOR] ────────────────────────┐
      │ REQUIREMENT: GOOGLE & YAHOO 2024 SENDER DIRECTIVE | COMPLIANCE: 100% COMPLETE       │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ HEADER FIELD           │ VALUE CONFIGURED                     │ COMPLIANCE STATUS   │
      ├────────────────────────┼──────────────────────────────────────┼─────────────────────┤
      │ List-Unsubscribe (URL) │ <https://api.domain.io/unsub?token=…>│ [●] VALID HTTPS     │
      │ List-Unsubscribe (Mail)│ <mailto:unsub@domain.io?subject=…>   │ [●] VALID MAILTO    │
      │ List-Unsubscribe-Post  │ List-Unsubscribe=One-Click           │ [●] STRICT RFC 8058 │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [VERIFICATION] Remote endpoint tested with HTTP POST: Returned 200 OK in 82ms.       │
      │ [FLAGS]        IsRFC8058Compliant: TRUE | HasHttpsEndpoint: TRUE | IsActive: TRUE   │
      │ [KEYS]         (t) Test Live Webhook  (g) Generate Header Tokens  (q) Exit          │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - RFC compliance validation card with test-dispatch button.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox unsub-audit`
  - Flags:
    - `--domain <string>`: Sender domain name (required).
    - `--is-strict-rfc8058`: Ensure both mailto and https one-click post exist (default: `true`).
    - `--has-live-test`: Send simulated probe POST request to verify 200 OK (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubscribeCompliance struct {
        DomainName         string    `json:"domainName"`
        HttpsURL           string    `json:"httpsUrl"`
        MailtoURI          string    `json:"mailtoUri"`
        HasOneClickPost    bool      `json:"hasOneClickPost"`
        LastVerifiedAt     time.Time `json:"lastVerifiedAt"`
        IsRFC8058Compliant bool      `json:"isRfc8058Compliant"`
        HasHttpsEndpoint   bool      `json:"hasHttpsEndpoint"`
        IsValidEndpoint    bool      `json:"isValidEndpoint"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnsubscribeCompliance (
        domain_name TEXT PRIMARY KEY,
        https_url TEXT NOT NULL,
        mailto_uri TEXT NOT NULL,
        has_one_click_post INTEGER NOT NULL DEFAULT 1,
        last_verified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_rfc8058_compliant INTEGER NOT NULL DEFAULT 1,
        has_https_endpoint INTEGER NOT NULL DEFAULT 1,
        is_valid_endpoint INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/email/sandbox/unsub_compliance.go`. Prevents immediate bulk spam blocking by enforcing mandatory unsubscribe architecture.

---

### Improvement 1059: Dynamic HTML Email Size Minifier & 102KB Clipping Preventer
- **ID:** `EMAIL-CLIP-1059`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Monitor and compress email HTML payloads to ensure total file size remains strictly under 102KB, preventing Gmail from clipping the message with `[Message clipped] View entire message` (which hides unsubscribe links and breaks layout).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Size gauge bar showing current raw bytes, compressed bytes, and delta to the 102KB clipping cliff.
    - Automated optimizations list: Comment stripping, attribute quote normalization, CSS dead-code elimination.
    - Terminal UI:
      ```text
      ┌─ [EMAIL HTML MINIFIER & 102KB CLIPPING PREVENTER] ──────────────────────────────────┐
      │ PAYLOAD: 64.2 KB / 102.0 KB SAFETY CEILING | HEADROOM: 37.8 KB (37.0% SAFE)         │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ METRIC                  │ RAW TEMPLATE │ OPTIMIZED TEMPLATE │ REDUCTION RATIO       │
      ├─────────────────────────┼──────────────┼────────────────────┼───────────────────────┤
      │ Total File Size         │     118.4 KB │            64.2 KB │ -45.8% (Saved 54.2 KB)│
      │ HTML Node Count         │    842 nodes │          842 nodes │ 0% (DOM Preserved)    │
      │ Extraneous Whitespace   │      24.2 KB │             0.0 KB │ -100% (Purged)        │
      │ HTML Comments           │       8.4 KB │             0.0 KB │ -100% (Purged)        │
      ├─────────────────────────────────────────────────────────────────────────────────────┤
      │ [GMAIL SAFETY STATUS] Message will NOT be clipped. Footer & Unsub fully visible.    │
      │ [FLAGS]               IsUnderClipLimit: TRUE | HasMinified: TRUE | IsPreserved: TRUE│
      │ [KEYS]                (m) Run Deep Minifier  (e) Export Minified File  (q) Exit     │
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Interactive file size cliff indicator with real-time payload breakdown.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox minify-html`
  - Flags:
    - `--file <path>`: Source HTML email file (required).
    - `--max-size-kb <int>`: Maximum threshold boundary (default: `100`).
    - `--is-comments-stripped`: Remove all non-IE conditional comments (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HTMLMinifierResult struct {
        FilePath         string    `json:"filePath"`
        RawSizeBytes     int       `json:"rawSizeBytes"`
        MinifiedSizeBytes int      `json:"minifiedSizeBytes"`
        SavingsPercent   float64   `json:"savingsPercent"`
        OptimizedAt      time.Time `json:"optimizedAt"`
        IsUnderClipLimit bool      `json:"isUnderClipLimit"`
        HasMinified      bool      `json:"hasMinified"`
        IsSafeForGmail   bool      `json:"isSafeForGmail"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HTMLMinifierResults (
        file_path TEXT PRIMARY KEY,
        raw_size_bytes INTEGER NOT NULL DEFAULT 0,
        minified_size_bytes INTEGER NOT NULL DEFAULT 0,
        savings_percent REAL NOT NULL DEFAULT 0.0,
        optimized_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_under_clip_limit INTEGER NOT NULL DEFAULT 1,
        has_minified INTEGER NOT NULL DEFAULT 1,
        is_safe_for_gmail INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed under `cli/email/sandbox/html_minifier.go`. Prevents Gmail mobile and web clients from severing tracking and legal footers.

---

### Improvement 1060: Zero-CDN Responsive Email Template Registry & Component Sandbox HUD
- **ID:** `EMAIL-REGISTRY-1060`
- **Cluster:** 106 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Maintain an offline, zero-external-CDN component registry of pre-tested, responsive email building blocks (Hero Banners, Data Tables, CTA Buttons, Footers) guaranteed to render consistently across all major email clients.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Interactive component library browser with keyboard selection (`hjkl` / tab).
    - Instant ASCII wireframe preview and code snippet export for chosen components.
    - Terminal UI:
      ```text
      ┌─ [ZERO-CDN EMAIL COMPONENT REGISTRY & BUILDER HUD] ─────────────────────────────────┐
      │ REGISTERED COMPONENTS: 24 | TESTING COVERAGE: 18 EMAIL CLIENTS | CDN ASSETS: 0      │
      ├─────────────────────────────────────┬───────────────────────────────────────────────┤
      │ COMPONENT CATEGORIES                │ COMPONENT PREVIEW (ASCII WIREFRAME)           │
      ├─────────────────────────────────────┼───────────────────────────────────────────────┤
      │ ▶ CTA Button (Bulletproof VML)      │ ┌───────────────────────────────────────────┐ │
      │   Two-Column Responsive Card        │ │       [ SCHEDULE ARCHITECTURE DEMO ]      │ │
      │   Data Comparison Table             │ └───────────────────────────────────────────┘ │
      │   Standard Compliance Footer        │ VML Fallback: YES | Border-Radius: 6px        │
      │   Notification Alert Banner         │ Mobile Stacking: Full Width (100%)            │
      ├─────────────────────────────────────┴───────────────────────────────────────────────┤
      │ [COMPONENT HEALTH] Tested across Outlook 2019, Apple Mail 17, Gmail App, Android.   │
      │ [FLAGS]            IsComponentValid: TRUE | HasZeroCdn: TRUE | IsReadyToCopy: TRUE  │
      │ [KEYS]             (enter) Copy HTML  (v) View Full Markup  (t) Test Render (q) Exit│
      └─────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure HTML/CSS component gallery showcasing visual layout variations without external dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox component-registry`
  - Flags:
    - `--category <string>`: Filter by component category (`buttons`, `cards`, `tables`, `footers`; default: `all`).
    - `--is-export-vml`: Include Outlook VML vector code in snippet (default: `true`).
    - `--has-zero-cdn-check`: Enforce zero external fonts or stylesheets (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailComponentSnippet struct {
        ComponentID      string    `json:"componentId"`
        ComponentName    string    `json:"componentName"`
        CategoryName     string    `json:"categoryName"`
        HTMLSnippet      string    `json:"htmlSnippet"`
        CompatibilityPct float64   `json:"compatibilityPct"`
        CreatedAt        time.Time `json:"createdAt"`
        IsComponentValid bool      `json:"isComponentValid"`
        HasZeroCdn       bool      `json:"hasZeroCdn"`
        IsBulletproofVML bool      `json:"isBulletproofVml"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailComponentSnippets (
        component_id TEXT PRIMARY KEY,
        component_name TEXT NOT NULL,
        category_name TEXT NOT NULL,
        html_snippet TEXT NOT NULL,
        compatibility_pct REAL NOT NULL DEFAULT 100.0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_component_valid INTEGER NOT NULL DEFAULT 1,
        has_zero_cdn INTEGER NOT NULL DEFAULT 1,
        is_bulletproof_vml INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_comp_category ON EmailComponentSnippets(category_name);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/email/sandbox/component_registry.go`. Standardizes email UI construction across all Antigravity and GitMap communication pipelines.

---
