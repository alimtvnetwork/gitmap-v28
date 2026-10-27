# Specification 346: Improvements 931 to 960 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 346
- **Slug:** `346-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/346-warm-up-ui-and-emailing-ui-improvements/03-improvements-931-to-960-interactive-tui-and-seed-intelligence.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (931 - 960)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 94** | Items 931–940 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 95** | Items 941–950 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 96** | Items 951–960 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 94: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 931–940)

---

### Improvement 931: Autonomous Real-Time Seed Pool Verification Matrix & Deliverability Heatmap HUD
- **ID:** `SEED-MATRIX-931`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a live terminal matrix view of seed mailboxes partitioned across tier-1 consumer and enterprise email service providers (Google Workspace, Microsoft 365, Yahoo/AOL, ProtonMail, and iCloud) tracking spam folder placement, inbox placement rates, and round-trip engagement loop latency.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density terminal matrix table rendering active seed nodes with 4 distinct state indicators: `[INBOX]` (emerald `#10b981`), `[SPAM]` (crimson `#ef4444`), `[DELAY]` (amber `#f59e0b`), and `[MISS]` (slate `#64748b`).
    - Interactive cursor with Lipgloss border styling (`#06b6d4`, electric cyan) for keyboard navigation (`hjkl` / arrow keys) showing instantaneous delivery headers, SPF/DKIM/DMARC alignment pass rate, and seed inbox ping latency.
    - Model view snippet:
      ```text
      ┌─ [SEED POOL VERIFICATION MATRIX: 160 ACTIVE SEED INBOXES] ──────────────────────────┐
      │ ESP GROUP         │ TOTAL │ INBOX RATE │ SPAM RATE │ AVG LATENCY │ STATUS            │
      ├───────────────────┼───────┼────────────┼───────────┼─────────────┼───────────────────┤
      │ Google Workspace  │    50 │      98.2% │      1.8% │      1.35s  │ [●] OPTIMAL       │
      │ Microsoft 365     │    45 │      95.6% │      4.4% │      1.95s  │ [●] OPTIMAL       │
      │ Yahoo / AOL       │    30 │      93.3% │      6.7% │      3.10s  │ [▲] ELEVATED LAG  │
      │ ProtonMail        │    20 │     100.0% │      0.0% │      0.82s  │ [●] PRISTINE      │
      │ Apple iCloud      │    15 │      86.7% │     13.3% │      3.95s  │ [■] ATTENTION REQ │
      ├───────────────────┴───────┴────────────┴───────────┴─────────────┴───────────────────┤
      │ [CURSOR: icloud-seed-14@icloud.com]                                                  │
      │ ├─ Last Arrival: SPAM (Folder: Junk) | Subject Hash: a4f8e91c                        │
      │ ├─ Authentication: SPF=PASS, DKIM=PASS, DMARC=PASS (Strict Alignment)               │
      │ ├─ Engagement Agent Action: Auto-Moved to Inbox & Flagged as Starred (1.5s ago)      │
      │ └─ Keys: [enter: Inspect Headers] [t: Trigger Engagement] [r: Refresh] [q: Quit]     │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Responsive pure CSS grid and inline SVG status indicators. Zero third-party web font or external CDN stylesheet dependencies.
    - Embedded SVG sparklines illustrating hourly inbox rate distributions per ESP category.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-matrix`
  - Flags:
    - `--esp <string>`: Filter seed pool by target provider (`google`, `microsoft`, `yahoo`, `proton`, `icloud`, or `all`; default: `all`).
    - `--min-inbox-rate <float>`: Highlight seed accounts failing target threshold (default: `90.0`).
    - `--watch`: Enable continuous real-time live terminal streaming mode (default: `true`).
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
    CREATE INDEX IF NOT EXISTS idx_seed_provider ON SeedPoolVerificationMatrix(provider_type, is_active);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates with `gitmap email warmup` command family and invokes Lemon seed agent orchestrator via `lemon/seed_daemon.go` to dispatch background IMAP checks without user blocking.

---

### Improvement 932: Interactive IMAP Folder Triage & Auto-Rescue Terminal Modal
- **ID:** `SEED-RESCUE-932`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Surface an interactive triage popup in Bubble Tea whenever seed emails land in Spam/Junk, automating IMAP `STORE +FLAGS (\\Seen)` and `COPY INBOX` operations with visual confirmation.
- **UI/UX Terminal & Web Specification:**
  - Modal overlay rendered using Lipgloss float margins, tinted background backdrop (`#0f172a`), and bright magenta border (`#ec4899`).
  - Terminal ASCII layout:
    ```text
    ┌──────────────────────── [SEED RESCUE IN ACTION] ────────────────────────┐
    │ Target Seed Account: seed-m365-corp-04@warmup-fleet.org                 │
    │ Message Subject:     "Q4 Infrastructure Review Alignment"              │
    │ Message ID:          <20261027140022.9912.seed@gitmap.internal>         │
    │ Detected Folder:     Junk Email (Classification: False Positive Spam)   │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Automated Remediation Steps:                                            │
    │  [✓] 1. Connect via IMAP TLS (port 993)                                 │
    │  [✓] 2. Issue IMAP UID MOVE -> INBOX                                    │
    │  [✓] 3. Mark message as Flagged & Starred                               │
    │  [●] 4. Generate human-simulated reply body (LLM token stream)          │
    │  [ ] 5. Dispatch reply via outbound SMTP TLS relay                      │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ [Status: RESCUING 80%]  Press [Esc: Abort] [Space: Pause] [Enter: Skip] │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-rescue`
  - Flags:
    - `--seed-id <string>`: Target seed mailbox UUID.
    - `--auto-reply`: Automatically send a contextual response to establish bidirectional conversation (default: `true`).
    - `--move-folder <string>`: Target destination folder (default: `INBOX`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedRescueTask struct {
        TaskID          string    `json:"taskId"`
        SeedID          string    `json:"seedId"`
        OriginalFolder  string    `json:"originalFolder"`
        TargetFolder    string    `json:"targetFolder"`
        RescueLatencyMs int64     `json:"rescueLatencyMs"`
        TriggeredAt     time.Time `json:"triggeredAt"`
        IsCompleted     bool      `json:"isCompleted"`
        HasAutoReplied  bool      `json:"hasAutoReplied"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedRescueLog (
        task_id TEXT PRIMARY KEY,
        seed_id TEXT NOT NULL,
        original_folder TEXT NOT NULL,
        target_folder TEXT NOT NULL,
        rescue_latency_ms INTEGER NOT NULL,
        triggered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_completed INTEGER NOT NULL DEFAULT 0,
        has_auto_replied INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_rescue_seed ON SeedRescueLog(seed_id, is_completed);
    ```
- **Antigravity / GitMap Integration Point:**
  - Hooked directly into `go-email-reader` IMAP event loop to trigger when unread junk mail headers match warmup seed tracking tokens.

---

### Improvement 933: Seed Reputation Decay Predictor & Risk Scoring HUD
- **ID:** `SEED-DECAY-933`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Predict upcoming mailbox reputation degradation before hard spam blocks occur, evaluating sudden bounce spikes, SPF/DKIM validation drift, and provider throttle events.
- **UI/UX Terminal & Web Specification:**
  - Horizontal bar gauges rendered with Lipgloss progress tokens (`░▒▓█`), color-ramping from cyan (`#06b6d4`) through amber (`#f59e0b`) to crimson (`#ef4444`).
  - Terminal HUD Preview:
    ```text
    ┌─ [SEED REPUTATION HEALTH & DECAY PREDICTION HUD] ────────────────────────┐
    │ POOL: Enterprise-M365-Group-A  |  MONITORED NODES: 64                   │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Aggregate Reputation Score:  89.4 / 100  [HEALTHY]                      │
    │ Predicted 72h Drift Vector:  -3.2 pts    [MODERATE RISK]                │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ RISK BREAKDOWN:                                                         │
    │  DKIM Key Drift Rate:    [████░░░░░░░░░░░░░░░░]  18.2% (Warning)        │
    │  Spam Trap Touch Risk:   [█░░░░░░░░░░░░░░░░░░░]   4.1% (Low)            │
    │  SMTP Greylist Lag:      [████████░░░░░░░░░░░░]  38.5% (Elevated)       │
    │  ESP Rate-Limit Halts:   [██░░░░░░░░░░░░░░░░░░]   9.0% (Normal)         │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Recommendation: Lower daily sending limit from 75 to 45 msgs/day.        │
    │ Action: Press [a] to Auto-Throttle Pool, [c] to Calibrate, [q] to Exit. │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-decay-audit`
  - Flags:
    - `--pool-id <string>`: Identifier of the warmup seed pool.
    - `--horizon-hours <int>`: Predictive lookahead window (default: `72`).
    - `--enforce-auto-throttle`: Automatically scale down send rate if risk score exceeds threshold (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedReputationMetric struct {
        PoolID           string    `json:"poolId"`
        ReputationScore  float64   `json:"reputationScore"`
        DriftVelocity    float64   `json:"driftVelocity"`
        DkimDriftRate    float64   `json:"dkimDriftRate"`
        SpamTrapTouchPct float64   `json:"spamTrapTouchPct"`
        ComputedAt       time.Time `json:"computedAt"`
        IsHighRisk       bool      `json:"isHighRisk"`
        HasActiveMitig   bool      `json:"hasActiveMitig"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedReputationHistory (
        pool_id TEXT NOT NULL,
        computed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        reputation_score REAL NOT NULL,
        drift_velocity REAL NOT NULL,
        dkim_drift_rate REAL NOT NULL,
        spam_trap_touch_pct REAL NOT NULL,
        is_high_risk INTEGER NOT NULL DEFAULT 0,
        has_active_mitig INTEGER NOT NULL DEFAULT 0,
        PRIMARY KEY (pool_id, computed_at)
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Feeds telemetry directly into GitMap Warmup Concurrency Manager to execute preemptive ramp-down schedules.

---

### Improvement 934: Multi-Provider Seed Conversation Thread Generator & Natural Sentiment Visualizer
- **ID:** `SEED-CONV-934`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Display real-time natural language sentiment classification and multi-turn message progression across simulated peer seed mailboxes to guarantee engagement authenticity.
- **UI/UX Terminal & Web Specification:**
  - Lipgloss conversation bubble cards with alternating alignment: Left-aligned for incoming seed, Right-aligned for outgoing sender.
  - Sentiment token indicators: `[POSITIVE: 0.94]` (emerald), `[NEUTRAL: 0.88]` (indigo), `[CONFIRMATORY: 0.91]` (cyan).
  - Terminal View:
    ```text
    ┌─ [CONVERSATION THREAD: SEED-CONV-88129] ────────────────────────────────┐
    │ PARTICIPANTS: alex@m365-corp.org <-> devon@gmail-consult.com            │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ [alex@m365-corp.org (10:14 AM)]                                         │
    │ "Could you review the updated API schemas for the Q4 rollout?"          │
    │  Sentiment: POSITIVE (0.87) | Tone: PROFESSIONAL COLLABORATIVE          │
    │                                                                         │
    │                     [devon@gmail-consult.com (10:19 AM)]                │
    │                     "Looks solid! Verified endpoints 12 through 16."     │
    │                      Sentiment: POSITIVE (0.95) | Tone: ENTHUSIASTIC    │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Thread Depth: 4 / 6 msgs | Avg Reply Delay: 4m 32s (Natural Variance)    │
    │ Active Status: [●] Natural Exchange in Progress                          │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup thread-visualize`
  - Flags:
    - `--thread-id <string>`: Specific conversation thread identifier.
    - `--min-turn <int>`: Minimum required conversational turns (default: `3`).
    - `--sentiment-check`: Compute and validate polarity scores on render (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedConversationMessage struct {
        MessageID       string    `json:"messageId"`
        ThreadID        string    `json:"threadId"`
        SenderAddress   string    `json:"senderAddress"`
        RecipientEmail  string    `json:"recipientEmail"`
        SentimentScore  float64   `json:"sentimentScore"`
        SentAt          time.Time `json:"sentAt"`
        IsOutbound      bool      `json:"isOutbound"`
        HasValidHeaders bool      `json:"hasValidHeaders"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedConversationThread (
        message_id TEXT PRIMARY KEY,
        thread_id TEXT NOT NULL,
        sender_address TEXT NOT NULL,
        recipient_email TEXT NOT NULL,
        sentiment_score REAL NOT NULL DEFAULT 0.0,
        sent_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_outbound INTEGER NOT NULL DEFAULT 1,
        has_valid_headers INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_thread_lookup ON SeedConversationThread(thread_id, sent_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Invoked during test warmups to demonstrate organic email patterns to external spam telemetry filters.

---

### Improvement 935: Seed Pool Geographic Distribution & ASN Routing Map
- **ID:** `SEED-GEO-935`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Map the autonomous network distribution and physical geographic location of seed mailbox IMAP/SMTP endpoints to avoid IP subnet clustering penalties.
- **UI/UX Terminal & Web Specification:**
  - High-resolution ANSI terminal world region histogram and ASN breakdown panel:
    ```text
    ┌─ [SEED AS NETWORK & GEOGRAPHIC DIVERSITY MAP] ──────────────────────────┐
    │ TOTAL NODES: 240  |  UNIQUE ASNs: 18  |  SUBNET CONCENTRATION: 4.2% (OK)│
    ├─────────────────────────────────────────────────────────────────────────┤
    │ REGION            │ COUNT │ TOP ASN                │ SUBNET DIVERSITY   │
    ├───────────────────┼───────┼────────────────────────┼────────────────────┤
    │ North America East│    85 │ AS15169 (Google Cloud) │ 98.2% [OPTIMAL]    │
    │ North America West│    60 │ AS8075 (Microsoft Corp)│ 96.5% [OPTIMAL]    │
    │ Europe West (FRA) │    45 │ AS16509 (Amazon.com)   │ 94.0% [OPTIMAL]    │
    │ Europe North (FIN)│    30 │ AS24940 (Hetzner)      │ 91.5% [SAFE]       │
    │ Asia Pacific (TYO)│    20 │ AS2516 (KDDI)          │ 88.0% [WATCH]      │
    ├───────────────────┴───────┴────────────────────────┴────────────────────┤
    │ Risk Warning: 0 ASNs exceed maximum 20% repository fleet threshold.    │
    │ Action: Press [e] to Rebalance Geographies, [p] to Ping RTTs.           │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-geo-audit`
  - Flags:
    - `--max-asn-concentration <float>`: Alert threshold percentage per single ASN (default: `20.0`).
    - `--validate-rdns`: Check matching reverse DNS PTR records for all IPs (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedNetworkNode struct {
        NodeID       string    `json:"nodeId"`
        IPAddress    string    `json:"ipAddress"`
        ASN          int       `json:"asn"`
        RegionCode   string    `json:"regionCode"`
        PingRTTMs    int64     `json:"pingRttMs"`
        RegisteredAt time.Time `json:"registeredAt"`
        IsValidPTR   bool      `json:"isValidPtr"`
        HasSubnetCol bool      `json:"hasSubnetCol"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedNetworkDistribution (
        node_id TEXT PRIMARY KEY,
        ip_address TEXT NOT NULL,
        asn INTEGER NOT NULL,
        region_code TEXT NOT NULL,
        ping_rtt_ms INTEGER NOT NULL,
        registered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_valid_ptr INTEGER NOT NULL DEFAULT 1,
        has_subnet_col INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_seed_asn ON SeedNetworkDistribution(asn, region_code);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `gitmap cluster` and `gitmap email warmup` infrastructure verification pre-flight checks.

---

### Improvement 936: Interactive Seed Mailbox Quarantine & Instant Expulsion TUI
- **ID:** `SEED-QUARANTINE-936`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a rapid keyboard-driven quarantine interface allowing operators or automated heuristics to isolate compromised, disabled, or locked seed accounts instantly.
- **UI/UX Terminal & Web Specification:**
  - Bubble Tea interactive multi-select checkbox list with search filter (`/` to search).
  - Selected nodes render with strike-through or dimmed text (`#64748b`) and red danger badge `[QUARANTINED]`.
  - Terminal Action Bar:
    ```text
    ┌─ [SEED MAILBOX QUARANTINE MANAGER] ─────────────────────────────────────┐
    │ FILTER: "auth-failed" (Found 3 candidate nodes)                         │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ [x] seed-acc-42@yahoo.com       │ Reason: IMAP 535 Auth Credentials Bad │
    │ [x] seed-m365-corp-12@corp.net │ Reason: Microsoft Captcha Challenge   │
    │ [ ] seed-google-biz-08@team.org │ Reason: Greylist Temporary Delay      │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Selected: 2 nodes targeted for immediate isolation.                    │
    │ Action: Press [Enter: Confirm Quarantine] [Space: Toggle] [Esc: Cancel] │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-quarantine`
  - Flags:
    - `--seed-id <string>`: Individual seed account identifier.
    - `--reason <string>`: Mandatory quarantine justification string.
    - `--auto-drain`: Safely reroute active inbound warmup threads (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedQuarantineRecord struct {
        RecordID         string    `json:"recordId"`
        SeedEmail        string    `json:"seedEmail"`
        QuarantineReason string    `json:"quarantineReason"`
        IsolatedAt       time.Time `json:"isolatedAt"`
        IsPermanent      bool      `json:"isPermanent"`
        HasDrainedQueue  bool      `json:"hasDrainedQueue"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedQuarantineLedger (
        record_id TEXT PRIMARY KEY,
        seed_email TEXT NOT NULL,
        quarantine_reason TEXT NOT NULL,
        isolated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_permanent INTEGER NOT NULL DEFAULT 0,
        has_drained_queue INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_quar_email ON SeedQuarantineLedger(seed_email);
    ```
- **Antigravity / GitMap Integration Point:**
  - Signals the active send dispatcher pool to exclude the quarantined mailbox from target recipient pools immediately.

---

### Improvement 937: Seed Engagement Activity Timeline & Interaction Inspector
- **ID:** `SEED-ACTIVITY-937`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Display a microsecond-accurate timeline of seed activities (read, star, unspam, click, reply, archive) to audit autonomous bot behavior against ESP anti-automation patterns.
- **UI/UX Terminal & Web Specification:**
  - Vertical timeline rendered with Lipgloss ASCII glyphs (`│`, `├─`, `└─`, `◆`), color-coded by action category:
    - `◆ OPEN` (cyan `#38bdf8`)
    - `◆ CLICK` (violet `#a855f7`)
    - `◆ RESCUE` (emerald `#10b981`)
    - `◆ REPLY` (fuchsia `#ec4899`)
  - Terminal Preview:
    ```text
    ┌─ [SEED ACTIVITY CHRONOLOGY: seed-google-31@fastmail.com] ───────────────┐
    │ 14:02:11.102  ◆ ARRIVAL     Inbound mail detected in INBOX              │
    │ 14:02:45.320  ◆ DELAY PAUSE Natural human read delay (34.2s simulated) │
    │ 14:02:45.890  ◆ OPEN        Message marked as \\Seen                    │
    │ 14:03:12.441  ◆ STAR        Flagged high priority                       │
    │ 14:04:02.118  ◆ CLICK       Followed validation tracking hyperlink      │
    │ 14:05:18.901  ◆ REPLY       Autonomous contextual response dispatched   │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Natural Behavior Score: 98.4 / 100 (Passes Google Bot Heuristic)        │
    │ Press [j/k: Scroll Timeline] [f: Filter Actions] [q: Back]              │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-activity-log`
  - Flags:
    - `--seed-id <string>`: Seed mailbox identifier.
    - `--limit <int>`: Event history limit (default: `50`).
    - `--format <string>`: Output format (`tui`, `json`, `yaml`; default: `tui`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedActivityEvent struct {
        EventID       string    `json:"eventId"`
        SeedID        string    `json:"seedId"`
        ActionType    string    `json:"actionType"`
        SimulatedSec  float64   `json:"simulatedSec"`
        OccurredAt    time.Time `json:"occurredAt"`
        IsOrganic     bool      `json:"isOrganic"`
        HasAttachment bool      `json:"hasAttachment"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedActivityLog (
        event_id TEXT PRIMARY KEY,
        seed_id TEXT NOT NULL,
        action_type TEXT NOT NULL,
        simulated_sec REAL NOT NULL DEFAULT 0.0,
        occurred_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_organic INTEGER NOT NULL DEFAULT 1,
        has_attachment INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_activity_seed ON SeedActivityLog(seed_id, occurred_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Used by CI/CD simulated warmup test suites to assert randomized delay adherence without static deterministic intervals.

---

### Improvement 938: Seed Inbox Storage Quota & Capacity Gauge HUD
- **ID:** `SEED-STORAGE-938`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Prevent unexpected warmup delivery failures caused by seed mailbox mailbox full (`552 5.2.2 Over quota`) errors through automated quota monitoring and rolling archival.
- **UI/UX Terminal & Web Specification:**
  - Lipgloss progress gauge rendering used bytes vs maximum mailbox capacity:
    ```text
    ┌─ [SEED STORAGE QUOTA & RETENTION HUD] ──────────────────────────────────┐
    │ TOTAL FLEET STORAGE: 48.2 GB / 120.0 GB (40.1% UTILIZATION)             │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ TOP CONGESTED MAILBOXES:                                                │
    │  seed-corp-01@outlook.com  [████████████████░░░░] 82.4% (Purge Queued)  │
    │  seed-biz-09@gmail.com     [████████████░░░░░░░░] 61.0% (Normal)        │
    │  seed-dev-22@icloud.com    [████████░░░░░░░░░░░░] 44.5% (Normal)        │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Rolling Retention Policy: Expire messages older than 30 days.           │
    │ Action: Press [p] to Trigger Dry-Run Purge, [r] to Refresh Quotas.      │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-quota-check`
  - Flags:
    - `--warning-pct <float>`: Mailbox warning percentage threshold (default: `80.0`).
    - `--auto-purge`: Automatically trigger IMAP expunge for expired items (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedStorageMetric struct {
        SeedID          string    `json:"seedId"`
        UsedBytes       int64     `json:"usedBytes"`
        TotalBytes      int64     `json:"totalBytes"`
        PercentageUsed  float64   `json:"percentageUsed"`
        LastPurgedAt    time.Time `json:"lastPurgedAt"`
        IsNearCapacity  bool      `json:"isNearCapacity"`
        HasPurgeEnabled bool      `json:"hasPurgeEnabled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedStorageCapacity (
        seed_id TEXT PRIMARY KEY,
        used_bytes INTEGER NOT NULL,
        total_bytes INTEGER NOT NULL,
        percentage_used REAL NOT NULL,
        last_purged_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_near_capacity INTEGER NOT NULL DEFAULT 0,
        has_purge_enabled INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Coordinates with the Lemon email daemon to issue `EXPUNGE` commands during low-traffic maintenance windows.

---

### Improvement 939: Seed DMARC/DKIM Cryptographic Key Rotation Verifier
- **ID:** `SEED-KEYROT-939`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Verify that all active sending domains and seed test accounts are validating against current 2048-bit RSA or Ed25519 DKIM keys with zero DNS propagation latency.
- **UI/UX Terminal & Web Specification:**
  - Terminal cryptographic status badges: `[DKIM: VALID 2048-BIT]` (emerald `#10b981`), `[KEY EXPIRED]` (crimson `#ef4444`), `[DNS MISMATCH]` (amber `#f59e0b`).
  - Terminal Table Layout:
    ```text
    ┌─ [CRYPTOGRAPHIC KEY ROTATION AUDITOR] ──────────────────────────────────┐
    │ DOMAIN             │ SELECTOR  │ KEY ALGORITHM │ ROTATION DUE │ STATUS  │
    ├────────────────────┼───────────┼───────────────┼──────────────┼─────────┤
    │ fleet-alpha.org    │ k1-2026   │ RSA-2048      │ in 42 days   │ [PASS]  │
    │ warmup-engine.net  │ default   │ Ed25519       │ in 88 days   │ [PASS]  │
    │ relays-fast.io     │ oct26     │ RSA-1024      │ DEPRECATED   │ [ALERT] │
    ├────────────────────┴───────────┴───────────────┴──────────────┴─────────┤
    │ Alert: relays-fast.io using weak 1024-bit RSA key. Upgrade required.   │
    │ Press [u: Generate 2048-bit Key] [v: Verify DNS Propagation] [q: Quit] │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup verify-crypto-keys`
  - Flags:
    - `--domain <string>`: Target domain to verify.
    - `--selector <string>`: DKIM DNS selector string (default: `default`).
    - `--min-bits <int>`: Minimum required RSA key length (default: `2048`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainCryptoStatus struct {
        DomainName    string    `json:"domainName"`
        Selector      string    `json:"selector"`
        KeyAlgorithm  string    `json:"keyAlgorithm"`
        KeyBitLength  int       `json:"keyBitLength"`
        LastRotatedAt time.Time `json:"lastRotatedAt"`
        IsValidLength bool      `json:"isValidLength"`
        HasDnsSynced  bool      `json:"hasDnsSynced"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainCryptoVerification (
        domain_name TEXT NOT NULL,
        selector TEXT NOT NULL,
        key_algorithm TEXT NOT NULL,
        key_bit_length INTEGER NOT NULL,
        last_rotated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_valid_length INTEGER NOT NULL DEFAULT 1,
        has_dns_synced INTEGER NOT NULL DEFAULT 1,
        PRIMARY KEY (domain_name, selector)
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executes as a mandatory pre-condition before initiating large volume warmup sequences.

---

### Improvement 940: Comprehensive Seed Fleet Diagnostic Summary & Report Exporter
- **ID:** `SEED-REPORT-940`
- **Cluster:** 94 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Generate a single, comprehensive Markdown, JSON, and terminal overview report aggregating all seed pool metrics, deliverability percentages, and remediation histories.
- **UI/UX Terminal & Web Specification:**
  - Clean Lipgloss multi-column summary card with ANSI colored sparkline histograms and one-click JSON/Markdown export trigger.
  - Terminal Report Preview:
    ```text
    ================================================================================
      GITMAP EMAIL WARMUP FLEET DIAGNOSTIC AUDIT REPORT
      Generated: 2026-10-27T14:05:00Z | Status: 96.8% INBOX HEALTH
    ================================================================================
      Active Seed Nodes:       160 Mailboxes
      Average Delivery Latency: 1.84 seconds
      Total Messages Triaged:  1,420 msgs (Last 24 Hours)
      Spam Folder Rescues:     34 Successful Moves (100% Resolved)
      ESP Health Ratios:
       - Google Workspace:     98.2% Inbox / 1.8% Spam
       - Microsoft 365:        95.6% Inbox / 4.4% Spam
       - Yahoo / AOL:          93.3% Inbox / 6.7% Spam
       - ProtonMail:           100.0% Inbox / 0.0% Spam
    --------------------------------------------------------------------------------
      Output saved to: 02-spec/21-app/reports/seed-health-2026-10-27.md
    ================================================================================
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-summary`
  - Flags:
    - `--output-file <string>`: Target relative file path for written report.
    - `--format <string>`: Output serialization format (`text`, `markdown`, `json`; default: `text`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedFleetSummaryReport struct {
        ReportID          string    `json:"reportId"`
        GeneratedAt       time.Time `json:"generatedAt"`
        TotalMailboxes    int       `json:"totalMailboxes"`
        AggregateInboxPct float64   `json:"aggregateInboxPct"`
        RescuedMessages   int       `json:"rescuedMessages"`
        IsHealthyFleet    bool      `json:"isHealthyFleet"`
        HasActionRequired bool      `json:"hasActionRequired"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedFleetReports (
        report_id TEXT PRIMARY KEY,
        generated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        total_mailboxes INTEGER NOT NULL,
        aggregate_inbox_pct REAL NOT NULL,
        rescued_messages INTEGER NOT NULL,
        is_healthy_fleet INTEGER NOT NULL DEFAULT 1,
        has_action_required INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Directly invoked by `./run.ps1` or `./run.sh` diagnostic subroutines.

---

# Cluster 95: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 941–950)

---

### Improvement 941: Interactive Concurrency Throttle & Dynamic Backoff Inversion HUD
- **ID:** `THROTTLE-GOV-941`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Provide a live Bubble Tea dashboard to monitor and manipulate sending concurrency across distinct destination domains in real time, automatically inverting backoff delays into positive throughput throttling.
- **UI/UX Terminal & Web Specification:**
  - Lipgloss interactive sliders adjusting concurrency workers per destination domain:
    ```text
    ┌─ [DYNAMIC CONCURRENCY GOVERNOR & BACKOFF INVERSION HUD] ────────────────┐
    │ DESTINATION DOMAIN  │ CONCURRENCY │ CURRENT MPS │ ERROR RATE │ BACKOFF RATE │
    ├─────────────────────┼─────────────┼─────────────┼────────────┼──────────────┤
    │ gmail.com           │  12 workers │  4.2 msg/s  │   0.2%     │   0.00s      │
    │ outlook.com         │   6 workers │  1.8 msg/s  │   1.1%     │  +0.50s (Lag)│
    │ yahoo.com           │   4 workers │  0.9 msg/s  │   4.8%     │  +2.20s (Inv)│
    │ icloud.com          │   2 workers │  0.4 msg/s  │   8.1%     │  +5.00s (Halt│
    ├─────────────────────┴─────────────┴─────────────┴────────────┴──────────────┤
    │ Inversion Engine Status: [ACTIVE] (Adaptive backoff active for 2 domains)   │
    │ Use [↑/↓] Select Domain  [+/-] Modify Concurrency  [space: Toggle Pause]    │
    └─────────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup concurrency-governor`
  - Flags:
    - `--domain <string>`: Focus monitor on single domain (default: `all`).
    - `--max-global-workers <int>`: Cap aggregate concurrency across all destinations (default: `50`).
    - `--auto-invert`: Dynamically calculate backoff curve based on SMTP 421/451 replies (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainConcurrencyState struct {
        TargetDomain     string    `json:"targetDomain"`
        ActiveWorkers    int       `json:"activeWorkers"`
        MessagesPerSec   float64   `json:"messagesPerSec"`
        RecentErrorRate  float64   `json:"recentErrorRate"`
        BackoffDelaySec  float64   `json:"backoffDelaySec"`
        LastAdjustedAt   time.Time `json:"lastAdjustedAt"`
        IsThrottled      bool      `json:"isThrottled"`
        HasInvertedQueue bool      `json:"hasInvertedQueue"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainConcurrencyGovernor (
        target_domain TEXT PRIMARY KEY,
        active_workers INTEGER NOT NULL DEFAULT 5,
        messages_per_sec REAL NOT NULL DEFAULT 1.0,
        recent_error_rate REAL NOT NULL DEFAULT 0.0,
        backoff_delay_sec REAL NOT NULL DEFAULT 0.0,
        last_adjusted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        has_inverted_queue INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Governs outbound goroutine worker pools in `lemon/smtp_dispatcher.go` preventing temporary IP reputation bans.

---

### Improvement 942: SMTP 421 / 451 Greylisting Detection & Automated Adaptive Re-Queue
- **ID:** `THROTTLE-GREYLIST-942`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Intercept transient greylisting response codes (`451 4.7.1 Greylisted - Please try again later`) and requeue messages with mathematical jitter without incrementing permanent failure counters.
- **UI/UX Terminal & Web Specification:**
  - Real-time animated bubble stream in Bubble Tea showing intercepted greylist events, retry countdown clocks, and progressive jitter allocation:
    ```text
    [GREYLIST TRAP DETECTED]
     ├─ Target: user@enterprise-target.com (MX: mail.enterprise-target.com)
     ├─ SMTP Code: 451 4.7.1 Service unavailable - try again in 300 seconds
     ├─ Adaptive Strategy: Jitter Re-Queue (+312s calculated)
     └─ Queue Placement: Priority Deferred Slot #14 (Next attempt: 14:10:22)
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup greylist-monitor`
  - Flags:
    - `--base-delay-sec <int>`: Minimum wait time before greylist retry (default: `300`).
    - `--jitter-window-sec <int>`: Random dispersion window to avoid lockstep retries (default: `60`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GreylistRetryEntry struct {
        QueueID        string    `json:"queueId"`
        TargetMX       string    `json:"targetMx"`
        OriginalCode   int       `json:"originalCode"`
        RetryScheduled time.Time `json:"retryScheduled"`
        AttemptCount   int       `json:"attemptCount"`
        IsPendingRetry bool      `json:"isPendingRetry"`
        HasPassedRetry bool      `json:"hasPassedRetry"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GreylistRetryQueue (
        queue_id TEXT PRIMARY KEY,
        target_mx TEXT NOT NULL,
        original_code INTEGER NOT NULL,
        retry_scheduled DATETIME NOT NULL,
        attempt_count INTEGER NOT NULL DEFAULT 1,
        is_pending_retry INTEGER NOT NULL DEFAULT 1,
        has_passed_retry INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_greylist_sched ON GreylistRetryQueue(retry_scheduled, is_pending_retry);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated with `go-email-sender` outbound pipeline to prevent false-positive campaign bounce drops.

---

### Improvement 943: Multi-IP Pool Outbound Round-Robin & Reputation Load Balancer
- **ID:** `THROTTLE-IPROT-943`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Distribute outbound warmup emails over multiple dedicated IP addresses using a health-weighted round-robin algorithm that decreases volume on IPs experiencing elevated bounces.
- **UI/UX Terminal & Web Specification:**
  - Dual bar charts showing outbound volume distribution alongside IP reputation scores:
    ```text
    ┌─ [MULTI-IP LOAD BALANCER & WEIGHTED ROUTER] ────────────────────────────┐
    │ IP ADDRESS     │ WEIGHT │ SENT (1H) │ REPUTATION │ HEALTH STATUS        │
    ├────────────────┼────────┼───────────┼────────────┼──────────────────────┤
    │ 198.51.100.14  │  40%   │   280 msg │ 99.1 / 100 │ [●] HEALTHY (Active) │
    │ 198.51.100.15  │  40%   │   276 msg │ 98.4 / 100 │ [●] HEALTHY (Active) │
    │ 198.51.100.16  │  15%   │   102 msg │ 86.0 / 100 │ [▲] DEGRADED (Ramped)│
    │ 198.51.100.17  │   5%   │    31 msg │ 72.3 / 100 │ [■] RESTRICTED       │
    ├────────────────┴────────┴───────────┴────────────┴──────────────────────┤
    │ Routing Mode: WEIGHTED HEALTH-PROPORTIONAL                             │
    │ Press [w: Adjust Weights] [f: Flush Queues] [q: Quit]                   │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ip-pool-status`
  - Flags:
    - `--pool-name <string>`: Name of IP routing group (default: `production-warmup`).
    - `--rebalance`: Dynamically recalculate weights based on 24h bounce rates (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type IPPoolMember struct {
        IPAddress       string    `json:"ipAddress"`
        WeightRatio     float64   `json:"weightRatio"`
        SentLastHour    int       `json:"sentLastHour"`
        ReputationScore float64   `json:"reputationScore"`
        LastUpdated     time.Time `json:"lastUpdated"`
        IsActive        bool      `json:"isActive"`
        HasRestrictions bool      `json:"hasRestrictions"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS OutboundIPPool (
        ip_address TEXT PRIMARY KEY,
        weight_ratio REAL NOT NULL DEFAULT 1.0,
        sent_last_hour INTEGER NOT NULL DEFAULT 0,
        reputation_score REAL NOT NULL DEFAULT 100.0,
        last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_active INTEGER NOT NULL DEFAULT 1,
        has_restrictions INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Configures outbound network socket binding in `lemon/transport_pool.go`.

---

### Improvement 944: Destination Domain MX Cluster Resolver & Grouped Rate Limiter
- **ID:** `THROTTLE-MXMAP-944`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Group multiple disparate recipient domains that share underlying mail servers (e.g., custom domains routed to `aspmx.l.google.com` or `mail.protection.outlook.com`) under a unified aggregate rate limit.
- **UI/UX Terminal & Web Specification:**
  - Tree-structured terminal view grouping child domains under parent MX infrastructures:
    ```text
    ┌─ [MX INFRASTRUCTURE AGGREGATE RATE LIMITER] ────────────────────────────┐
    │ ▼ GOOGLE WORKSPACE INFRASTRUCTURE (aspmx.l.google.com)                  │
    │   ├─ Capacity: 35 / 60 MPS | Health: 99.4%                              │
    │   ├─ Domains: company-a.com, tech-startup.io, global-group.org (82 total)│
    │   └─ Concurrency Throttle: 18 Active Workers                            │
    │ ▼ MICROSOFT 365 CLUSTER (mail.protection.outlook.com)                   │
    │   ├─ Capacity: 22 / 40 MPS | Health: 96.1%                              │
    │   ├─ Domains: enterprise-inc.net, consulting-firm.co (54 total)        │
    │   └─ Concurrency Throttle: 12 Active Workers                            │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mx-group-status`
  - Flags:
    - `--resolve-cache`: Refresh DNS MX lookups for all queued recipient domains (default: `true`).
    - `--max-mps-per-mx <int>`: Cap maximum messages per second per physical MX cluster (default: `50`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MXClusterRateLimit struct {
        MXClusterHost  string    `json:"mxClusterHost"`
        MaxMessagesSec int       `json:"maxMessagesSec"`
        CurrentMPS     float64   `json:"currentMps"`
        DomainCount    int       `json:"domainCount"`
        LastEvaluated  time.Time `json:"lastEvaluated"`
        IsSaturated    bool      `json:"isSaturated"`
        HasBackpressure bool     `json:"hasBackpressure"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MXClusterLimits (
        mx_cluster_host TEXT PRIMARY KEY,
        max_messages_sec INTEGER NOT NULL DEFAULT 50,
        current_mps REAL NOT NULL DEFAULT 0.0,
        domain_count INTEGER NOT NULL DEFAULT 1,
        last_evaluated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_saturated INTEGER NOT NULL DEFAULT 0,
        has_backpressure INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Directly prevents spam filtering blocks from major email providers who enforce cluster-wide limits across custom vanity domains.

---

### Improvement 945: Exponential Backoff Dynamic Inversion Visualizer
- **ID:** `THROTTLE-VIS-945`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Visualize real-time mathematical backoff curves (linear, exponential, and inverse decay) in terminal ASCII plots as errors fluctuate.
- **UI/UX Terminal & Web Specification:**
  - Lipgloss XY-plot graph rendering backoff curves over time:
    ```text
    [BACKOFF INVERSION CURVE: OUTLOOK.COM]
    Delay (s)
     10 ┤                ╭─────────── (Inverted Recovery Plateau)
      8 ┤              ╭─╯
      6 ┤          ╭───╯
      4 ┤     ╭────╯
      2 ┤ ╭───╯
      0 ┼─┴───┴───┴───┴───┴───┴───┴───┴─── Time (minutes)
        0   2   4   6   8  10  12  14
    Current Formula: D(t) = Min(D_max, D_base * 2^(err_ratio * 1.5))
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup backoff-curve`
  - Flags:
    - `--domain <string>`: Target domain to graph.
    - `--width <int>`: Terminal graph column width (default: `60`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BackoffCurvePoint struct {
        TimestampSec int64   `json:"timestampSec"`
        DelaySeconds float64 `json:"delaySeconds"`
        ErrorDensity float64 `json:"errorDensity"`
        IsProjected  bool    `json:"isProjected"`
        HasClamped   bool    `json:"hasClamped"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BackoffCurveSamples (
        sample_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        timestamp_sec INTEGER NOT NULL,
        delay_seconds REAL NOT NULL,
        error_density REAL NOT NULL,
        is_projected INTEGER NOT NULL DEFAULT 0,
        has_clamped INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_curve_domain ON BackoffCurveSamples(target_domain, timestamp_sec);
    ```
- **Antigravity / GitMap Integration Point:**
  - Used for developer tuning and visual diagnostics in GitMap TUI mode.

---

### Improvement 946: Dynamic Warmup Ramp Schedule Automator
- **ID:** `THROTTLE-RAMP-946`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Automate Day 1 through Day 30 sending ramp curves per mailbox, automatically pausing progression if daily spam or bounce rates exceed critical safety parameters.
- **UI/UX Terminal & Web Specification:**
  - Multi-stage progression timeline showing current mailbox day and projected volume ceilings:
    ```text
    ┌─ [WARMUP RAMP PROGRESSION: Day 14 of 30] ───────────────────────────────┐
    │ Current Daily Quota:  65 / 70 sent  (92.8%)                             │
    │ Health Gate Checks:   SPF [PASS]  DKIM [PASS]  Bounces: 0.0% [PASS]     │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Schedule Trajectory:                                                    │
    │  [Days 1-7]   Start Ramp (10 -> 35 msgs/day)    [COMPLETED]             │
    │  [Days 8-14]  Phase 2 Ramp (40 -> 70 msgs/day)  [ACTIVE - DAY 14]       │
    │  [Days 15-21] Enterprise Tier (75 -> 120 msgs)  [PENDING TOMORROW]      │
    │  [Days 22-30] Production Readiness (150+ msgs)  [LOCKED]                │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Gate Status: [●] ALL GATES GREEN. Tomorrow: Increment Quota to 75.      │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ramp-status`
  - Flags:
    - `--mailbox <string>`: Email address of sender mailbox.
    - `--force-advance`: Manually step forward one day in the schedule (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type WarmupRampSchedule struct {
        MailboxAddress string    `json:"mailboxAddress"`
        CurrentDay     int       `json:"currentDay"`
        MaxDailyQuota  int       `json:"maxDailyQuota"`
        SentToday      int       `json:"sentToday"`
        StartDate      time.Time `json:"startDate"`
        IsPaused       bool      `json:"isPaused"`
        HasMetGates    bool      `json:"hasMetGates"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupRampLedger (
        mailbox_address TEXT PRIMARY KEY,
        current_day INTEGER NOT NULL DEFAULT 1,
        max_daily_quota INTEGER NOT NULL DEFAULT 10,
        sent_today INTEGER NOT NULL DEFAULT 0,
        start_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_paused INTEGER NOT NULL DEFAULT 0,
        has_met_gates INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Enforced by `lemon/send_limiter.go` prior to dispatching any new outbound warmup message.

---

### Improvement 947: SMTP Connection Pool Warmup & Keep-Alive Governor
- **ID:** `THROTTLE-KEEPALIVE-947`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Optimize TLS handshake latency and TCP connection reuse by managing warm connection pools to high-volume MX endpoints with graceful idle timeouts.
- **UI/UX Terminal & Web Specification:**
  - Table of active socket connections, idle duration meters, and connection reuse metrics:
    ```text
    ┌─ [PERSISTENT SMTP CONNECTION POOL] ─────────────────────────────────────┐
    │ DESTINATION HOST               │ OPEN CONNS │ IDLE SEC │ REUSE COUNT    │
    ├────────────────────────────────┼────────────┼──────────┼────────────────┤
    │ aspmx.l.google.com:25          │     8      │   4.2s   │ 142 msgs       │
    │ mail.protection.outlook.com:25 │     5      │  12.0s   │  88 msgs       │
    │ mx-aol.mail.gm0.yahoodns.net:25│     3      │   8.5s   │  45 msgs       │
    ├────────────────────────────────┴────────────┴──────────┴────────────────┤
    │ Total Open Sockets: 16 | Efficiency Gain: -420ms TLS Overhead/Msg       │
    │ Press [f: Flush Idle] [c: Close All] [q: Return]                        │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup pool-status`
  - Flags:
    - `--max-idle-sec <int>`: Socket idle timeout before closing (default: `30`).
    - `--max-open-per-host <int>`: Upper limit on open sockets per MX host (default: `10`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ConnectionPoolHostMetric struct {
        HostAddress      string    `json:"hostAddress"`
        OpenConnections  int       `json:"openConnections"`
        IdleDurationSec  float64   `json:"idleDurationSec"`
        TotalReusedCount int       `json:"totalReusedCount"`
        LastConnected    time.Time `json:"lastConnected"`
        IsSaturated      bool      `json:"isSaturated"`
        HasActiveLeak    bool      `json:"hasActiveLeak"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SMTPConnectionPoolStats (
        host_address TEXT PRIMARY KEY,
        open_connections INTEGER NOT NULL DEFAULT 0,
        idle_duration_sec REAL NOT NULL DEFAULT 0.0,
        total_reused_count INTEGER NOT NULL DEFAULT 0,
        last_connected DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_saturated INTEGER NOT NULL DEFAULT 0,
        has_active_leak INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into core Lemon outbound connection multiplexer.

---

### Improvement 948: Outbound SMTP Circuit Breaker & Hard Failure Isolator
- **ID:** `THROTTLE-CIRCUIT-948`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Implement an autonomous circuit breaker that trips instantly to `OPEN` state upon detecting 5xx series hard bounce anomalies or Spamhaus/RBL listing responses.
- **UI/UX Terminal & Web Specification:**
  - Prominent red circuit breaker status indicator with countdown timer to test probe `HALF-OPEN` state:
    ```text
    ┌────────────────────── [CIRCUIT BREAKER: TRIPPED] ──────────────────────┐
    │ TARGET DOMAIN:    yahoo.com                                             │
    │ CURRENT STATE:    OPEN (All outbound traffic halted immediately)        │
    │ TRIP CAUSE:       553 5.3.0 High bounce threshold (RBL Listing Detected)│
    │ TRIPPED AT:       14:08:19 UTC (42 seconds ago)                         │
    │ COOLDOWN WINDOW:  15 minutes (Next probe at 14:23:19 UTC)               │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Action Required: Audit domain DNS records and IP reputation.            │
    │ Override: Press [F9: Force Half-Open Probe] [Esc: Dismiss]              │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup circuit-status`
  - Flags:
    - `--domain <string>`: Check status for destination domain.
    - `--reset`: Manually reset tripped breaker back to `CLOSED` (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CircuitBreakerState struct {
        DomainName     string    `json:"domainName"`
        StateCode      string    `json:"stateCode"` // CLOSED, OPEN, HALF-OPEN
        FailureCount   int       `json:"failureCount"`
        TrippedAt      time.Time `json:"trippedAt"`
        NextProbeAt    time.Time `json:"nextProbeAt"`
        IsOpenState    bool      `json:"isOpenState"`
        HasManualReset bool      `json:"hasManualReset"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS OutboundCircuitBreaker (
        domain_name TEXT PRIMARY KEY,
        state_code TEXT NOT NULL DEFAULT 'CLOSED',
        failure_count INTEGER NOT NULL DEFAULT 0,
        tripped_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        next_probe_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_open_state INTEGER NOT NULL DEFAULT 0,
        has_manual_reset INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Protects sender reputation by terminating doomed delivery attempts before ESP blocklists cascade.

---

### Improvement 949: Warmup Dispatch Priority Queue Re-Ordering Engine
- **ID:** `THROTTLE-PRIO-949`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Dynamically reprioritize outbound message queues, elevating high-reputation seed interactions ahead of standard warmup exchanges to maintain inbox velocity.
- **UI/UX Terminal & Web Specification:**
  - Reorderable list view showing active queue priority tiers:
    ```text
    ┌─ [WARMUP DISPATCH QUEUE RE-ORDERING] ───────────────────────────────────┐
    │ PRIORITY │ MESSAGE TYPE         │ IN QUEUE │ TARGET ESP    │ HEAD LATENCY│
    ├──────────┼──────────────────────┼──────────┼───────────────┼─────────────┤
    │ P0 [CRIT]│ Seed Rescue Reply    │    4 msgs│ Google / M365 │    0.2s     │
    │ P1 [HIGH]│ Seed Thread Continue │   18 msgs│ Multi-Provider│    1.1s     │
    │ P2 [NORM]│ Initial Warmup Ping  │   62 msgs│ Global Fleet  │    3.4s     │
    │ P3 [LOW] │ Archival Exchange    │   12 msgs│ Internal Seed │    8.0s     │
    ├──────────┴──────────────────────┴──────────┴───────────────┴─────────────┤
    │ Queue Health: 96 queued | Dispatch throughput: 3.8 msg/sec              │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup queue-reorder`
  - Flags:
    - `--elevate-type <string>`: Elevate specific warmup message class to P0.
    - `--drain-priority <string>`: Flush queue in strict priority descending order.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PriorityQueueEntry struct {
        ItemUUID      string    `json:"itemUuid"`
        PriorityLevel int       `json:"priorityLevel"`
        MessageType   string    `json:"messageType"`
        QueuedAt      time.Time `json:"queuedAt"`
        IsDispatched  bool      `json:"isDispatched"`
        HasElevated   bool      `json:"hasElevated"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupPriorityQueue (
        item_uuid TEXT PRIMARY KEY,
        priority_level INTEGER NOT NULL DEFAULT 2,
        message_type TEXT NOT NULL,
        queued_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_dispatched INTEGER NOT NULL DEFAULT 0,
        has_elevated INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_prio_disp ON WarmupPriorityQueue(priority_level, is_dispatched, queued_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Feeds dispatcher consumer channels in `lemon/worker_pool.go`.

---

### Improvement 950: Multi-Domain Throughput Telemetry & Historical Trends HUD
- **ID:** `THROTTLE-HIST-950`
- **Cluster:** 95 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Render a 24-hour historical heatmap and latency trendline in the terminal, charting throughput vs throttle rejections across every major target domain.
- **UI/UX Terminal & Web Specification:**
  - ANSI colored terminal heatmap grid (24 hourly buckets x 6 major ESP categories):
    ```text
    ┌─ [24-HOUR WARMUP THROUGHPUT & REJECTION HEATMAP] ───────────────────────┐
    │ ESP \ Hour  │ 00 02 04 06 08 10 12 14 16 18 20 22 │ AVG MPS │ TOTAL VOL │
    ├─────────────┼─────────────────────────────────────┼─────────┼───────────┤
    │ Google      │ ░░ ░░ ▒▒ ▒▒ ██ ██ ██ ██ ▓▓ ▓▓ ▒▒ ░░ │ 3.4 m/s │ 2,840 msg │
    │ Microsoft   │ ░░ ░░ ▒▒ ▒▒ ▓▓ ▓▓ ██ ▓▓ ▒▒ ▒▒ ░░ ░░ │ 2.8 m/s │ 2,120 msg │
    │ Yahoo       │ ░░ ░░ ░░ ▒▒ ▒▒ ▒▒ ▓▓ ▒▒ ▒▒ ░░ ░░ ░░ │ 1.2 m/s │ 1,050 msg │
    │ Apple       │ ░░ ░░ ░░ ░░ ▒▒ ▒▒ ▒▒ ▒▒ ░░ ░░ ░░ ░░ │ 0.8 m/s │   620 msg │
    ├─────────────┴─────────────────────────────────────┴─────────┴───────────┤
    │ Legend: [░: <0.5 MPS] [▒: 1-2 MPS] [▓: 2-3 MPS] [█: 3+ MPS]             │
    │ Press [d: Export Daily CSV] [q: Quit]                                   │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup throughput-history`
  - Flags:
    - `--days <int>`: History duration in days (default: `1`).
    - `--export-csv <string>`: Destination relative path for CSV export.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ThroughputHourlyBucket struct {
        DomainCategory string    `json:"domainCategory"`
        BucketHour     time.Time `json:"bucketHour"`
        DeliveredCount int       `json:"deliveredCount"`
        RejectedCount  int       `json:"rejectedCount"`
        AvgMPS         float64   `json:"avgMps"`
        IsPeakHour     bool      `json:"isPeakHour"`
        HasRejections  bool      `json:"hasRejections"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupHourlyThroughput (
        domain_category TEXT NOT NULL,
        bucket_hour DATETIME NOT NULL,
        delivered_count INTEGER NOT NULL DEFAULT 0,
        rejected_count INTEGER NOT NULL DEFAULT 0,
        avg_mps REAL NOT NULL DEFAULT 0.0,
        is_peak_hour INTEGER NOT NULL DEFAULT 0,
        has_rejections INTEGER NOT NULL DEFAULT 0,
        PRIMARY KEY (domain_category, bucket_hour)
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Provides long-term capacity telemetry to the Lemon planning engine.

---

# Cluster 96: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 951–960)

---

### Improvement 951: Zero-CDN Responsive HTML Sanitizer & Security Sandbox HUD
- **ID:** `SEC-SANDBOX-951`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Parse, audit, and sanitize warmup and marketing email HTML payloads inside an isolated, local security sandbox without external CDN scripts, validating iframe bans, script tag stripping, and CSS injection guards.
- **UI/UX Terminal & Web Specification:**
  - Split terminal screen with raw HTML code view on left and security analysis breakdown on right:
    ```text
    ┌─ [EMAIL HTML SECURITY SANITIZER & AUDIT HUD] ───────────────────────────┐
    │ RAW PAYLOAD (Lines 1-18)           │ SECURITY VIOLATION AUDIT           │
    ├────────────────────────────────────┼────────────────────────────────────┤
    │ <div style="font-family:sans;">    │ [✓] Script Tags: NONE (PASS)       │
    │   <h1>Welcome to the Review</h1>   │ [✓] Iframes: NONE (PASS)           │
    │   <p>Click below to verify:</p>    │ [▲] External Fonts: Google Fonts   │
    │   <link rel="stylesheet"           │     (Auto-stripped to local stack) │
    │     href="https://fonts.googleapis"│ [✓] Base64 Inline: 0 images (PASS) │
    │   <a href="https://gitmap.internal │ [✓] Tracking Pixels: Validated     │
    │ </div>                             │ [✓] Zero External CDNs Detected    │
    ├────────────────────────────────────┴────────────────────────────────────┤
    │ Sanitizer Outcome: 1 Violation Remediated | HTML Weight: 1.42 KB (Safe) │
    │ Press [s: Save Sanitized HTML] [v: View Full Diff] [q: Exit]            │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sanitize-html`
  - Flags:
    - `--input-file <string>`: Path to candidate email HTML file.
    - `--strict-zero-cdn`: Enforce complete removal of remote stylesheet/font links (default: `true`).
    - `--output-clean <string>`: Path to write sanitized output payload.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SanitizerAuditResult struct {
        PayloadHash       string    `json:"payloadHash"`
        TotalViolations   int       `json:"totalViolations"`
        RemovedTagsCount  int       `json:"removedTagsCount"`
        CleanedPayload    string    `json:"cleanedPayload"`
        EvaluatedAt       time.Time `json:"evaluatedAt"`
        IsZeroCdnPass     bool      `json:"isZeroCdnPass"`
        HasScriptsStripped bool     `json:"hasScriptsStripped"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HtmlSanitizerAudit (
        payload_hash TEXT PRIMARY KEY,
        total_violations INTEGER NOT NULL DEFAULT 0,
        removed_tags_count INTEGER NOT NULL DEFAULT 0,
        cleaned_payload TEXT NOT NULL,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_zero_cdn_pass INTEGER NOT NULL DEFAULT 1,
        has_scripts_stripped INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `go-email-sender` before queuing any outbound HTML MIME body part.

---

### Improvement 952: Zero-Dependency Pure CSS Dark Mode Email Compatibility Checker
- **ID:** `SEC-DARKMODE-952`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Analyze email templates for `@media (prefers-color-scheme: dark)` overrides and inverted contrast ratios, ensuring illegible black-text-on-black-background rendering failures are flagged.
- **UI/UX Terminal & Web Specification:**
  - Dual ANSI rendering simulator displaying standard Light Mode vs inverted Dark Mode terminal previews:
    ```text
    ┌─ [EMAIL DUAL-THEME CONTRAST SIMULATOR] ─────────────────────────────────┐
    │ LIGHT MODE PREVIEW                 │ DARK MODE PREVIEW                  │
    ├────────────────────────────────────┼────────────────────────────────────┤
    │ ┌────────────────────────────────┐ │ ┌────────────────────────────────┐ │
    │ │ Background: #ffffff (White)    │ │ │ Background: #121212 (Dark)     │ │
    │ │ Header:     #0f172a (Navy)     │ │ │ Header:     #f8fafc (Light)    │ │
    │ │ Body:       #334155 (Slate)    │ │ │ Body:       #94a3b8 (Slate)    │ │
    │ │ Contrast:   12.4:1 [AAA PASS]  │ │ │ Contrast:   8.9:1 [AAA PASS]   │ │
    │ └────────────────────────────────┘ │ └────────────────────────────────┘ │
    ├────────────────────────────────────┴────────────────────────────────────┤
    │ Dark Mode Compatibility Score: 100% (WCAG 2.1 AAA Compliant)            │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email check-darkmode`
  - Flags:
    - `--html-file <string>`: Relative path to HTML template file.
    - `--min-contrast <float>`: Minimum acceptable contrast ratio (default: `4.5`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeEvaluation struct {
        TemplateUUID     string    `json:"templateUuid"`
        LightContrast    float64   `json:"lightContrast"`
        DarkContrast     float64   `json:"darkContrast"`
        EvaluatedAt      time.Time `json:"evaluatedAt"`
        IsAaaCompliant   bool      `json:"isAaaCompliant"`
        HasMediaOverride bool      `json:"hasMediaOverride"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DarkModeAudits (
        template_uuid TEXT PRIMARY KEY,
        light_contrast REAL NOT NULL,
        dark_contrast REAL NOT NULL,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_aaa_compliant INTEGER NOT NULL DEFAULT 1,
        has_media_override INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates with template linting scripts in `03-ai-scripts/`.

---

### Improvement 953: Inline SVG vs CID Attachment Security Optimizer
- **ID:** `SEC-ATTACH-953`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Detect inline SVG tags or remote imagery, converting them automatically into localized RFC 2392 multipart/related `cid:` MIME attachments with sanitized XML attributes.
- **UI/UX Terminal & Web Specification:**
  - Interactive table showing transformed media entities:
    ```text
    ┌─ [INLINE MEDIA TO CID ATTACHMENT CONVERTER] ────────────────────────────┐
    │ SOURCE ENTITY             │ DETECTED TYPE │ ACTION TAKEN    │ CID HASH  │
    ├───────────────────────────┼───────────────┼─────────────────┼───────────┤
    │ <svg viewBox="0 0 24 24"> │ Raw XML SVG   │ Sanitized -> PNG│ cid:ico01 │
    │ <img src="data:image/png; │ Base64 Inline │ Extracted -> CID│ cid:img02 │
    ├───────────────────────────┴───────────────┴─────────────────┴───────────┤
    │ MIME Optimization: Payload weight reduced by 34.2%                      │
    │ Press [c: Apply Conversion] [q: Discard]                                │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email convert-cid`
  - Flags:
    - `--input <string>`: Source HTML template.
    - `--output <string>`: Target optimized HTML file.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CIDAttachmentMeta struct {
        ContentID       string `json:"contentId"`
        MimeType        string `json:"mimeType"`
        ByteSize        int64  `json:"byteSize"`
        IsXmlSanitized  bool   `json:"isXmlSanitized"`
        HasValidHeaders bool   `json:"hasValidHeaders"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CIDMediaAssets (
        content_id TEXT PRIMARY KEY,
        mime_type TEXT NOT NULL,
        byte_size INTEGER NOT NULL,
        is_xml_sanitized INTEGER NOT NULL DEFAULT 1,
        has_valid_headers INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Hooked into `go-email-sender` MIME builder module.

---

### Improvement 954: Tracking Pixel Local Obfuscator & Privacy Guard HUD
- **ID:** `SEC-PIXEL-954`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify that open-tracking 1x1 GIF pixels use zero external tracking networks, validating HMAC-signed self-hosted internal tokens with anti-replay timestamp verification.
- **UI/UX Terminal & Web Specification:**
  - Terminal cryptographic token verification panel:
    ```text
    ┌─ [TRACKING PIXEL PRIVACY & HMAC SECURITY AUDIT] ────────────────────────┐
    │ TRACKING URL: https://gitmap.internal/t/a8f9c1e4d0b2                     │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ HMAC-SHA256 Signature: VALID (Key: id_warmup_token_v1)                 │
    │ Nonce / Timestamp:      2026-10-27T14:10:00Z (Within 24h validity)      │
    │ External Leak Check:    ZERO THIRD-PARTY CALLS (100% On-Premise)        │
    │ Cache-Control Header:   no-store, no-cache, must-revalidate             │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Privacy Rating: [PRISTINE] Meets Enterprise GDPR & CCPA Zero-Share      │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email verify-tracking-token`
  - Flags:
    - `--token <string>`: Encrypted tracking string.
    - `--secret-env <string>`: Environment variable containing HMAC secret (default: `GITMAP_TRACKING_SECRET`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TrackingTokenAudit struct {
        TokenHash    string    `json:"tokenHash"`
        RecipientID  string    `json:"recipientId"`
        IssuedAt     time.Time `json:"issuedAt"`
        IsValidToken bool      `json:"isValidToken"`
        HasExpired   bool      `json:"hasExpired"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TrackingTokenRegistry (
        token_hash TEXT PRIMARY KEY,
        recipient_id TEXT NOT NULL,
        issued_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_valid_token INTEGER NOT NULL DEFAULT 1,
        has_expired INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Validated by GitMap email telemetry listeners.

---

### Improvement 955: Link Reputation & Destination Host Reputation Scanner
- **ID:** `SEC-LINKSCAN-955`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan every hyperlink embedded in email templates against local domain blocklists, Google Safe Browsing APIs, and redirect hops to ensure no flagged or spammy URLs exist.
- **UI/UX Terminal & Web Specification:**
  - Link inventory table with status indicators:
    ```text
    ┌─ [EMBEDDED HYPERLINK REPUTATION SCANNER] ───────────────────────────────┐
    │ DESTINATION URL               │ REDIRECTS │ RBL STATUS  │ SAFE BROWSING │
    ├───────────────────────────────┼───────────┼─────────────┼───────────────┤
    │ https://gitmap.internal/docs  │ 0 hops    │ [PASS] CLEAN│ [PASS] OK     │
    │ https://company.org/signup    │ 1 hop     │ [PASS] CLEAN│ [PASS] OK     │
    │ http://bit.ly/unknown-redirect│ 3 hops    │ [FAIL] SUSP │ [WARN] UNKNOWN│
    ├───────────────────────────────┴───────────┴─────────────┴───────────────┤
    │ Warning: Generic shortener detected (bit.ly). High spam-score risk.     │
    │ Action: Press [r: Replace with Vanity URL] [q: Exit]                    │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email scan-links`
  - Flags:
    - `--template <string>`: Path to email file.
    - `--disallow-shorteners`: Reject bit.ly, tinyurl, and similar link shorteners (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LinkScanResult struct {
        TargetURL       string    `json:"targetUrl"`
        HopCount        int       `json:"hopCount"`
        FinalDestination string   `json:"finalDestination"`
        ScannedAt       time.Time `json:"scannedAt"`
        IsCleanTarget   bool      `json:"isCleanTarget"`
        HasShortener    bool      `json:"hasShortener"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailLinkAudit (
        target_url TEXT PRIMARY KEY,
        hop_count INTEGER NOT NULL DEFAULT 0,
        final_destination TEXT NOT NULL,
        scanned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_clean_target INTEGER NOT NULL DEFAULT 1,
        has_shortener INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Runs in CI/CD pipeline prior to approving marketing or warmup campaign templates.

---

### Improvement 956: Zero-CDN CSS Inliner & Redundant Rule Purger
- **ID:** `SEC-CSSINLINE-956`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Inline all external and `<style>` block CSS declarations directly into HTML tag `style=""` attributes while safely discarding unused selectors to reduce payload size.
- **UI/UX Terminal & Web Specification:**
  - Compression summary HUD:
    ```text
    ┌─ [ZERO-CDN CSS INLINER & PAYLOAD PURGER] ───────────────────────────────┐
    │ Input HTML Size:    14.8 KB                                             │
    │ Output Inlined Size: 8.2 KB (44.6% Reduction)                           │
    │ Rules Processed:    84 selectors                                        │
    │ Unused Purged:      38 unreferenced CSS rules removed                   │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Output Status: [●] Pristine Single-File Self-Contained Payload          │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email inline-css`
  - Flags:
    - `--src <string>`: Path to uncompiled template.
    - `--dest <string>`: Destination path for inlined output.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CSSInliningMetrics struct {
        SourceHash      string  `json:"sourceHash"`
        OriginalBytes   int64   `json:"originalBytes"`
        InlinedBytes    int64   `json:"inlinedBytes"`
        PurgedRuleCount int     `json:"purgedRuleCount"`
        IsCleanOutput   bool    `json:"isCleanOutput"`
        HasUnusedRules  bool    `json:"hasUnusedRules"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CSSInliningHistory (
        source_hash TEXT PRIMARY KEY,
        original_bytes INTEGER NOT NULL,
        inlined_bytes INTEGER NOT NULL,
        purged_rule_count INTEGER NOT NULL,
        is_clean_output INTEGER NOT NULL DEFAULT 1,
        has_unused_rules INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into template build scripts in `go-email-sender`.

---

### Improvement 957: Spam Filter Heuristic Trigger Word Scanner & Visual Heatmap
- **ID:** `SEC-SPAMWORD-957`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan subject lines and body copy for high-risk spam words ("100% FREE", "ACT NOW", "RISK-FREE GUARANTEE", "WIRE TRANSFER"), highlighting density in terminal colors.
- **UI/UX Terminal & Web Specification:**
  - Terminal text preview with crimson highlighted trigger words and spam density meter:
    ```text
    ┌─ [SPAM TRIGGER WORD SCANNER] ───────────────────────────────────────────┐
    │ SUBJECT: "ACT NOW for your 100% FREE consultation!"                     │
    │           ^^^^^^^          ^^^^^^^^^                                    │
    │ Trigger Density: 4.8% (Danger Threshold: > 1.5%)                        │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ DETECTED WORDS:                                                         │
    │  - "ACT NOW"     (Weight: +2.4 pts) -> Suggested: "Review at your leisure"│
    │  - "100% FREE"   (Weight: +3.1 pts) -> Suggested: "Complimentary"       │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Predicted SpamAssassin Score: 4.2 / 5.0 [HIGH PROBABILITY JUNK]         │
    │ Press [r: Apply AI Suggestions] [q: Back]                               │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit-spam-words`
  - Flags:
    - `--subject <string>`: Candidate subject line.
    - `--body-file <string>`: Path to body text file.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamWordAudit struct {
        AuditUUID     string    `json:"auditUuid"`
        SubjectText   string    `json:"subjectText"`
        SpamScore     float64   `json:"spamScore"`
        WordHitsCount int       `json:"wordHitsCount"`
        AuditedAt     time.Time `json:"auditedAt"`
        IsHighSpam    bool      `json:"isHighSpam"`
        HasUrgentTone bool      `json:"hasUrgentTone"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamWordAudits (
        audit_uuid TEXT PRIMARY KEY,
        subject_text TEXT NOT NULL,
        spam_score REAL NOT NULL,
        word_hits_count INTEGER NOT NULL,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_high_spam INTEGER NOT NULL DEFAULT 0,
        has_urgent_tone INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Prevents accidental dispatch of high-risk warmup templates.

---

### Improvement 958: Plaintext MIME Part Sync & Parity Validator
- **ID:** `SEC-PLAINTEXT-958`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Ensure every multipart/alternative email contains a high-fidelity plaintext counterpart that accurately mirrors all HTML links and content without divergence penalties.
- **UI/UX Terminal & Web Specification:**
  - Side-by-side terminal diff showing HTML stripped text vs actual plaintext part:
    ```text
    ┌─ [MULTIPART/ALTERNATIVE PARITY CHECKER] ────────────────────────────────┐
    │ HTML EXTRACTED TEXT                │ PLAINTEXT MIME PART                │
    ├────────────────────────────────────┼────────────────────────────────────┤
    │ Welcome to the GitMap team.        │ Welcome to the GitMap team.        │
    │ Click https://gitmap.internal      │ Click https://gitmap.internal      │
    │ Unsubscribe: https://gitmap.org/u  │ Unsubscribe: https://gitmap.org/u  │
    ├────────────────────────────────────┴────────────────────────────────────┤
    │ Parity Match: 100.0% [PERFECT ALIGNMENT]                                │
    │ Spam Score Benefit: -1.2 pts (Valid multipart structure)               │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email verify-mime-parity`
  - Flags:
    - `--raw-eml <string>`: Path to raw RFC 822 `.eml` file.
    - `--min-match-pct <float>`: Minimum acceptable similarity percentage (default: `95.0`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MimeParityResult struct {
        MessageUUID     string    `json:"messageUuid"`
        ParityScore     float64   `json:"parityScore"`
        HasHtmlPart     bool      `json:"hasHtmlPart"`
        HasPlaintext    bool      `json:"hasPlaintext"`
        CheckedAt       time.Time `json:"checkedAt"`
        IsSynchronized  bool      `json:"isSynchronized"`
        HasMissingLinks bool      `json:"hasMissingLinks"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MimeParityAudits (
        message_uuid TEXT PRIMARY KEY,
        parity_score REAL NOT NULL,
        has_html_part INTEGER NOT NULL DEFAULT 1,
        has_plaintext INTEGER NOT NULL DEFAULT 1,
        checked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_synchronized INTEGER NOT NULL DEFAULT 1,
        has_missing_links INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Enforced by `go-email-sender` MIME validation pipeline.

---

### Improvement 959: Unsubscribe Header & RFC 8058 One-Click Post Compliance HUD
- **ID:** `SEC-UNSUB-959`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify presence and compliance of `List-Unsubscribe` and `List-Unsubscribe-Post: List-Unsubscribe=One-Click` headers to meet Google and Yahoo 2024 bulk sender mandates.
- **UI/UX Terminal & Web Specification:**
  - Terminal header compliance checklist:
    ```text
    ┌─ [RFC 8058 ONE-CLICK UNSUBSCRIBE COMPLIANCE AUDITOR] ───────────────────┐
    │ MANDATORY BULK SENDER HEADERS:                                          │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ [✓] List-Unsubscribe: <mailto:unsub@gitmap.org>, <https://gitmap.org/u> │
    │ [✓] List-Unsubscribe-Post: List-Unsubscribe=One-Click                   │
    │ [✓] Post Endpoint Responds HTTP 200 to POST (Verified in 120ms)         │
    ├─────────────────────────────────────────────────────────────────────────┤
    │ Compliance Status: [100% VALID] Passes Yahoo & Gmail Bulk Mandates      │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email verify-unsub-headers`
  - Flags:
    - `--eml-file <string>`: Path to candidate email file.
    - `--test-post-endpoint`: Issue dry-run HTTP POST probe (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubscribeComplianceStatus struct {
        EmailHeaderID    string    `json:"emailHeaderId"`
        HasMailtoUnsub   bool      `json:"hasMailtoUnsub"`
        HasHttpUnsub     bool      `json:"hasHttpUnsub"`
        HasOneClickPost  bool      `json:"hasOneClickPost"`
        VerifiedAt       time.Time `json:"verifiedAt"`
        IsBulkCompliant  bool      `json:"isBulkCompliant"`
        HasActiveEndpoint bool     `json:"hasActiveEndpoint"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnsubscribeCompliance (
        email_header_id TEXT PRIMARY KEY,
        has_mailto_unsub INTEGER NOT NULL DEFAULT 1,
        has_http_unsub INTEGER NOT NULL DEFAULT 1,
        has_one_click_post INTEGER NOT NULL DEFAULT 1,
        verified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_bulk_compliant INTEGER NOT NULL DEFAULT 1,
        has_active_endpoint INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Mandatory pre-flight gate in `lemon/campaign_dispatcher.go`.

---

### Improvement 960: Autonomous Template Security Sandbox Verification Matrix
- **ID:** `SEC-MATRIX-960`
- **Cluster:** 96 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Execute a comprehensive security and deliverability test suite across all email templates in the repository, generating an executive pass/fail matrix report.
- **UI/UX Terminal & Web Specification:**
  - Multi-template terminal test matrix with Lipgloss styled checkmarks:
    ```text
    ┌─ [REPOSITORY-WIDE EMAIL TEMPLATE SECURITY MATRIX] ──────────────────────┐
    │ TEMPLATE NAME             │ ZERO-CDN │ DARKMODE │ SPAMWORD │ RFC 8058 │ │
    ├───────────────────────────┼──────────┼──────────┼──────────┼──────────┼─┤
    │ warmup-welcome.html       │  [PASS]  │  [PASS]  │  [PASS]  │  [PASS]  │ │
    │ seed-exchange-v2.html     │  [PASS]  │  [PASS]  │  [PASS]  │  [PASS]  │ │
    │ monthly-digest.html       │  [WARN]  │  [PASS]  │  [PASS]  │  [PASS]  │ │
    │ password-reset.html       │  [PASS]  │  [PASS]  │  [PASS]  │  [N/A]   │ │
    ├───────────────────────────┴──────────┴──────────┴──────────┴──────────┴─┤
    │ 4 templates audited: 3 PRISTINE, 1 WARNING (monthly-digest: external img│
    │ Report saved to: 02-spec/21-app/reports/template-security-matrix.md     │
    └─────────────────────────────────────────────────────────────────────────┘
    ```
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit-template-suite`
  - Flags:
    - `--dir <string>`: Directory containing email templates (default: `templates/email`).
    - `--fail-on-warning`: Exit with code 1 if any warning is detected (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TemplateMatrixReport struct {
        SuiteID         string    `json:"suiteId"`
        TotalTemplates  int       `json:"totalTemplates"`
        PassedTemplates int       `json:"passedTemplates"`
        WarningCount    int       `json:"warningCount"`
        ExecutedAt      time.Time `json:"executedAt"`
        IsSuiteGreen    bool      `json:"isSuiteGreen"`
        HasCriticalFail bool      `json:"hasCriticalFail"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TemplateSuiteAuditReports (
        suite_id TEXT PRIMARY KEY,
        total_templates INTEGER NOT NULL,
        passed_templates INTEGER NOT NULL,
        warning_count INTEGER NOT NULL,
        executed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_suite_green INTEGER NOT NULL DEFAULT 1,
        has_critical_fail INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Configured as a Git pre-commit hook and CI/CD automated test gate in `run.ps1` and `run.sh`.

---
