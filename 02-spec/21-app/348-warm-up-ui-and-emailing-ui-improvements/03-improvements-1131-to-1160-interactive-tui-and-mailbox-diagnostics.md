# Specification 348: Improvements 1131 to 1160 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 348
- **Slug:** `348-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/348-warm-up-ui-and-emailing-ui-improvements/03-improvements-1131-to-1160-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/347-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/346-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (1131 - 1160)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 114** | Items 1131–1140 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 115** | Items 1141–1150 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 116** | Items 1151–1160 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 114: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 1131–1140)

---

### Improvement 1131: Autonomous Multi-Tier Seed Pool Health Matrix & Placement Heatmap HUD
- **ID:** `SEED-MATRIX-1131`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Render a high-density, real-time terminal matrix table classifying active warmup seed mailboxes across primary consumer and enterprise ESP networks (Google Workspace, Microsoft 365, Yahoo/AOL, ProtonMail, and Apple iCloud) with granular placement rate telemetry, round-trip latency scoring, and automated spam rescue dispatch.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density terminal matrix table rendering active seed nodes with 4 distinct state indicators: `[INBOX]` (emerald `#10b981`), `[SPAM]` (crimson `#ef4444`), `[DELAY]` (amber `#f59e0b`), and `[MISS]` (slate `#64748b`).
    - Interactive cursor with Lipgloss border styling (`#06b6d4`, electric cyan) for keyboard navigation (`hjkl` / arrow keys) displaying instantaneous delivery headers, SPF/DKIM/DMARC alignment pass rate, and seed inbox ping latency.
    - Model view snippet:
      ```text
      ┌─ [SEED POOL VERIFICATION MATRIX: 320 ACTIVE SEED INBOXES] ──────────────────────────┐
      │ ESP CLUSTER       │ TOTAL │ INBOX RATE │ SPAM RATE │ AVG LATENCY │ STATUS            │
      ├───────────────────┼───────┼────────────┼───────────┼─────────────┼───────────────────┤
      │ Google Workspace  │   110 │      99.1% │      0.9% │      0.98s  │ [●] OPTIMAL       │
      │ Microsoft 365     │    90 │      97.4% │      2.6% │      1.45s  │ [●] OPTIMAL       │
      │ Yahoo / AOL       │    50 │      95.2% │      4.8% │      2.30s  │ [▲] ELEVATED LAG  │
      │ ProtonMail        │    40 │     100.0% │      0.0% │      0.65s  │ [●] PRISTINE      │
      │ Apple iCloud      │    30 │      92.3% │      7.7% │      2.95s  │ [■] ATTENTION REQ │
      ├───────────────────┴───────┴────────────┴───────────┴─────────────┴───────────────────┤
      │ [CURSOR: icloud-seed-27@icloud.com]                                                  │
      │ ├─ Last Arrival: SPAM (Folder: Junk) | Subject Hash: a9b3e108                        │
      │ ├─ Authentication: SPF=PASS, DKIM=PASS, DMARC=PASS (Strict Alignment)               │
      │ ├─ Engagement Agent Action: Auto-Moved to Inbox & Starred (0.4s ago)                 │
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
    CREATE INDEX IF NOT EXISTS idx_seed_provider ON SeedPoolVerificationMatrix(provider_type, is_healthy);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/seedpool/matrix.go` and surfaces interactive Bubble Tea matrix through GitMap CLI subcommand `gitmap email warmup seed-matrix`.

---

### Improvement 1132: Real-Time IMAP IDLE Synchronizer & Seed Arrival Notification HUD
- **ID:** `SEED-IMAP-1132`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Establish zero-polling IMAP IDLE streaming connections across distributed seed mailboxes, capturing incoming warmup messages instantaneously, logging exact transit latency, and displaying real-time push telemetry in a compact Bubble Tea status ticker.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Split-pane terminal layout: Upper pane displays active TCP/TLS socket connection health across seed accounts; lower pane streams a ticker of arriving test payloads with microsecond transit stamps.
    - Lipgloss tokens: Connection green `#22c55e`, IDLE standby violet `#a855f7`, Alert red `#ef4444`.
    - Model view snippet:
      ```text
      ┌─ [IMAP IDLE REAL-TIME STREAM: 184 ACTIVE SOCKET CONNECTIONS] ────────────────────────┐
      │ ACTIVE SESSIONS: 184 | IDLE LISTENING: 179 | RECONNECTING: 5 | DROPPED: 0             │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [14:02:19.412] ARRIVAL  seed-gsuite-04@acme.org   <- out-node-01.relay (1.14s transit)│
      │ [14:02:20.089] ARRIVAL  seed-o365-12@corp.io      <- out-node-03.relay (0.87s transit)│
      │ [14:02:21.341] IDLE_ACK seed-yahoo-07@yahoo.com   KeepAlive PING OK (RTT 42ms)        │
      │ [14:02:22.502] ARRIVAL  seed-proton-02@pm.me      <- out-node-02.relay (1.42s transit)│
      │ [14:02:23.110] WARN     seed-icloud-09@icloud.com TLS Session Timeout -> Re-authent... │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native WebSocket stream pipe rendering incoming arrival pings via pure Vanilla JS DOM inserts with zero external bundle assets.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup imap-idle-stream`
  - Flags:
    - `--timeout-keepalive <duration>`: Keepalive NOOP/IDLE renegotiation period (default: `24m`).
    - `--max-concurrent-sockets <int>`: Max open TLS connections per host (default: `64`).
    - `--is-auto-reconnect`: Re-establish dropped IMAP sessions with exponential backoff (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ImapIdleSession struct {
        SessionID       string        `json:"sessionId"`
        MailboxAddress  string        `json:"mailboxAddress"`
        HostEndpoint    string        `json:"hostEndpoint"`
        ConnectionAge   time.Duration `json:"connectionAge"`
        LastPayloadHash string        `json:"lastPayloadHash"`
        IsConnected     bool          `json:"isConnected"`
        IsIdleActive    bool          `json:"isIdleActive"`
        HasTlsStrict    bool          `json:"hasTlsStrict"`
    }

    type ImapArrivalEvent struct {
        EventID        string    `json:"eventId"`
        SessionID      string    `json:"sessionId"`
        TransitLatency float64   `json:"transitLatency"`
        FolderDetected string    `json:"folderDetected"`
        ReceivedAt     time.Time `json:"receivedAt"`
        IsInbox        bool      `json:"isInbox"`
        HasValidDkim   bool      `json:"hasValidDkim"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ImapIdleSessionRecord (
        session_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        host_endpoint TEXT NOT NULL,
        connection_age_sec INTEGER NOT NULL DEFAULT 0,
        last_payload_hash TEXT,
        is_connected INTEGER NOT NULL DEFAULT 1,
        is_idle_active INTEGER NOT NULL DEFAULT 1,
        has_tls_strict INTEGER NOT NULL DEFAULT 1,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TABLE IF NOT EXISTS ImapArrivalEventRecord (
        event_id TEXT PRIMARY KEY,
        session_id TEXT NOT NULL REFERENCES ImapIdleSessionRecord(session_id),
        transit_latency REAL NOT NULL,
        folder_detected TEXT NOT NULL,
        received_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_inbox INTEGER NOT NULL DEFAULT 1,
        has_valid_dkim INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented under `pkg/warmup/imapidle/listener.go`, wired into background daemons via GitMap background task harness.

---

### Improvement 1133: Synthetic AI Peer Engagement Workflow & Thread Conversation Synthesizer
- **ID:** `SEED-ENGAGE-1133`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Autonomously generate realistic multi-turn conversational email replies between warmup mailboxes and seed peers, synthesizing contextually coherent thread branches, simulated human reading delays, and star/label interactions to build pristine sender authority.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual dialog simulator tree showing threaded exchanges between seed actors:
      ```text
      ┌─ [WARMUP PEER CONVERSATION SYNTHESIZER: THREAD #tr-89412] ───────────────────────────┐
      │ Root Node: sales-lead@target.com -> "Q3 Operations Review & Roadmap Planning"        │
      │ ├── [Turn 1: +14m] seed-o365-04: "Thanks for sharing, could you clarify slide 4?"    │
      │ │   └── [Turn 2: +28m] out-warmup-02: "Slide 4 incorporates recent latency audits."  │
      │ │       └── [Turn 3: +41m] seed-gsuite-08: "Understood, marked as high priority."   │
      │ ├────────────────────────────────────────────────────────────────────────────────────┤
      │ │ Simulation Parameters:                                                             │
      │ │ ├─ Reading Delay: Gaussian distribution (μ = 4.2m, σ = 1.1m)                       │
      │ │ ├─ User Actions: Mark Read [✓] | Star [✓] | Archive [✓] | Move To Folder [✓]        │
      │ │ └─ Thread Coherence Index: 0.984 (High Natural Tone Score)                         │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss tokens: Root blue `#3b82f6`, Node green `#10b981`, Accent amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - SVG-based interactive tree chart rendering multi-party conversation branches with zero external library CDN dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup synthesize-thread`
  - Flags:
    - `--thread-depth <int>`: Number of turns to simulate across peers (default: `4`).
    - `--mean-delay-min <float>`: Mean human reading delay before reply (default: `5.0`).
    - `--is-star-enabled`: Automatically flag messages as important/starred (default: `true`).
    - `--is-auto-archive`: Move threads out of inbox after completion (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EngagementThreadNode struct {
        NodeID         string    `json:"nodeId"`
        ThreadID       string    `json:"threadId"`
        SenderAddress  string    `json:"senderAddress"`
        RecipientAddr  string    `json:"recipientAddr"`
        Subject        string    `json:"subject"`
        BodySnippet    string    `json:"bodySnippet"`
        TurnIndex      int       `json:"turnIndex"`
        ScheduledAt    time.Time `json:"scheduledAt"`
        IsReplied      bool      `json:"isReplied"`
        IsStarred      bool      `json:"isStarred"`
        HasHumanDelay  bool      `json:"hasHumanDelay"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EngagementThreadRecord (
        node_id TEXT PRIMARY KEY,
        thread_id TEXT NOT NULL,
        sender_address TEXT NOT NULL,
        recipient_addr TEXT NOT NULL,
        subject TEXT NOT NULL,
        body_snippet TEXT NOT NULL,
        turn_index INTEGER NOT NULL DEFAULT 0,
        scheduled_at DATETIME NOT NULL,
        is_replied INTEGER NOT NULL DEFAULT 0,
        is_starred INTEGER NOT NULL DEFAULT 0,
        has_human_delay INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_thread_turn ON EngagementThreadRecord(thread_id, turn_index);
    ```
- **Antigravity / GitMap Integration Point:**
  - Coordinated through `pkg/warmup/engagement/synthesizer.go` with state tracked in SQLite split database `repodb/pipeline.db`.

---

### Improvement 1134: Deep Envelope Return-Path & DMARC Forensic Report (RUF) Analyzer
- **ID:** `SEED-FORENSIC-1134`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Ingest and decode RFC 6591 DMARC Forensic Reporting (RUF) XML and MIME failure payloads delivered to seed error traps, highlighting exact SPF identity misalignment and DKIM cryptographic signature mismatches in an interactive inspection viewer.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Detailed diff-split window comparing RFC 5322 header `From:` against RFC 5321 envelope `Return-Path:` and DKIM `d=` selector tags.
    - Lipgloss styles: Pass green `#22c55e`, Fail red `#ef4444`, Header cyan `#38bdf8`.
    - Model view snippet:
      ```text
      ┌─ [DMARC RUF FORENSIC REPORT VIEWER: REPORT #ruf-9104] ───────────────────────────────┐
      │ SOURCE IP: 198.51.100.44 | REPORTING MTA: dkim-report.google.com | TIME: 14:15:02    │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [AUTHENTICATION SUMMARY]                                                             │
      │ ├─ SPF Result:   PASS (Client IP matched ip4:198.51.100.0/24)                       │
      │ ├─ DKIM Result:  FAIL (Body hash verify mismatch bh=47DEQpj8... != computed hash)   │
      │ └─ DMARC Policy: REJECT (Alignment: Relaxed | Eval Domain: company.com)              │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [FORENSIC PAYLOAD HEADERS]                                                           │
      │ Return-Path: <bounce-relay@out.company.com>                                          │
      │ From: "Corporate Support" <support@company.com>                                      │
      │ DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed; d=company.com; s=k1; ...      │
      │ Diagnostic-Code: smtp; 550 5.7.26 DKIM validation failed for altered body content    │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Side-by-side XML-to-HTML tree visualizer with syntax-highlighted header tokens rendered without client-side remote scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup forensic-inspect`
  - Flags:
    - `--report-id <string>`: Specific RUF report identifier to inspect.
    - `--min-severity <string>`: Filter by severity (`info`, `warn`, `fail`; default: `warn`).
    - `--is-raw-xml`: Output raw decrypted forensic XML payload.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ForensicRufReport struct {
        ReportID       string    `json:"reportId"`
        SourceIP       string    `json:"sourceIp"`
        TargetDomain   string    `json:"targetDomain"`
        ReportingMta   string    `json:"reportingMta"`
        FailureType    string    `json:"failureType"`
        EnvelopeSender string    `json:"envelopeSender"`
        HeaderFrom     string    `json:"headerFrom"`
        ReceivedTime   time.Time `json:"receivedTime"`
        HasSpfPass     bool      `json:"hasSpfPass"`
        HasDkimPass    bool      `json:"hasDkimPass"`
        IsAligned      bool      `json:"isAligned"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ForensicRufReportRecord (
        report_id TEXT PRIMARY KEY,
        source_ip TEXT NOT NULL,
        target_domain TEXT NOT NULL,
        reporting_mta TEXT NOT NULL,
        failure_type TEXT NOT NULL,
        envelope_sender TEXT NOT NULL,
        header_from TEXT NOT NULL,
        received_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        has_spf_pass INTEGER NOT NULL DEFAULT 0,
        has_dkim_pass INTEGER NOT NULL DEFAULT 0,
        is_aligned INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_ruf_domain ON ForensicRufReportRecord(target_domain, received_time);
    ```
- **Antigravity / GitMap Integration Point:**
  - Diagnostic engine residing in `pkg/warmup/forensic/analyzer.go`, integrated with GitMap email diagnostics command suite.

---

### Improvement 1135: Automated Spam Folder Rescue & Sentiment Weighting Daemon
- **ID:** `SEED-RESCUE-1135`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Periodically scan Spam/Junk mail folders across managed seed mailboxes, identify falsely flagged warmup messages, unflag spam flags, reposition messages to the Primary Inbox, apply custom priority labels, and log heuristic rescue friction telemetry.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live terminal rescue progress bar with step-by-step IMAP state transition transitions:
      ```text
      ┌─ [AUTOMATED SPAM RESCUE DAEMON: RUNNING CYCLE #412] ─────────────────────────────────┐
      │ SCANNED: 320 INBOXES | SPAM DETECTIONS: 14 | RESCUED: 14 | RECOVERY RATE: 100%       │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [RESCUE EVENT #1092]                                                                 │
      │ ├─ Mailbox: seed-outlook-09@outlook.com                                              │
      │ ├─ Detection: Message in 'Junk' folder (Subject: "Weekly Architecture Sync")        │
      │ ├─ Actions Executed:                                                                 │
      │ │   1. MOVE 'Junk' -> 'Inbox'                        [DONE in 410ms]                │
      │ │   2. UNSET IMAP Flag \Junk                         [DONE in 120ms]                │
      │ │   3. SET Custom Header 'X-Warmup-Rescued: true'    [APPLIED]                      │
      │ │   4. ADD Seed Engagement Star                      [CONFIRMED]                    │
      │ └─ Status: RESOLVED -> Positive ESP Sender Signal Broadcasted                        │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss tokens: Processing cyan `#06b6d4`, Success green `#22c55e`, Warning amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Timeline view of rescued messages with interactive ESP classification badges.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup spam-rescue`
  - Flags:
    - `--poll-interval <duration>`: Rescue daemon sweep frequency (default: `10m`).
    - `--max-batch-rescues <int>`: Maximum rescues per sweep cycle (default: `50`).
    - `--is-dry-run`: Inspect without moving or modifying IMAP flags (default: `false`).
    - `--has-auto-star`: Star rescued messages to maximize ISP trust score (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamRescueEvent struct {
        RescueID       string    `json:"rescueId"`
        MailboxAddress string    `json:"mailboxAddress"`
        MessageUID     string    `json:"messageUid"`
        SubjectHash    string    `json:"subjectHash"`
        EspProvider    string    `json:"espProvider"`
        SourceFolder   string    `json:"sourceFolder"`
        ExecutionLagMs int64     `json:"executionLagMs"`
        RescuedAt      time.Time `json:"rescuedAt"`
        IsCompleted    bool      `json:"isCompleted"`
        HasStarred     bool      `json:"hasStarred"`
        HasUnflagged   bool      `json:"hasUnflagged"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamRescueEventRecord (
        rescue_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        message_uid TEXT NOT NULL,
        subject_hash TEXT NOT NULL,
        esp_provider TEXT NOT NULL,
        source_folder TEXT NOT NULL,
        execution_lag_ms INTEGER NOT NULL DEFAULT 0,
        rescued_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_completed INTEGER NOT NULL DEFAULT 1,
        has_starred INTEGER NOT NULL DEFAULT 1,
        has_unflagged INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_rescue_mailbox ON SpamRescueEventRecord(mailbox_address, rescued_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Packaged in `pkg/warmup/rescue/daemon.go` and controlled through GitMap daemon worker processes.

---

### Improvement 1136: Seed Mailbox IP-Reputation & RBL Blacklist Micro-Monitor
- **ID:** `SEED-RBL-1136`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Continuously query authoritative Real-Time Blackhole Lists (Spamhaus ZEN, Barracuda BRBL, Spamcop, Sorbs, and Invaluement) for all sending egress IPs and seed relay interfaces, providing immediate visual terminal warnings upon listing detection.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live multi-column grid displaying IP nodes against major DNSBL providers:
      ```text
      ┌─ [EGRESS IP & SEED RELAY BLACKLIST TELEMETRY (RBL/DNSBL)] ───────────────────────────┐
      │ TOTAL MONITORED IPS: 16 | CLEAN NODES: 15 | LISTED NODES: 1 | QUERY LATENCY: 28ms     │
      ├─────────────────┬──────────┬──────────┬──────────┬──────────┬──────────┬──────────┤
      │ IP ADDRESS      │ SPAMHAUS │ BARRACUDA│ SPAMCOP  │ SORBS    │ INVALU   │ STATUS   │
      ├─────────────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
      │ 198.51.100.12   │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [●] PASS │
      │ 198.51.100.14   │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [●] PASS │
      │ 198.51.100.18   │ [LISTED] │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [▲] WARN │
      │ 203.0.113.88    │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [CLEAR]  │ [●] PASS │
      ├─────────────────┴──────────┴──────────┴──────────┴──────────┴──────────┴──────────┤
      │ ALERT: 198.51.100.18 Listed on Spamhaus SBL (Return: 127.0.0.2). Autocooling in effect│
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Clear green `#22c55e`, Listed red `#ef4444`, Border subtle `#475569`.
  - **Zero-CDN Web Dashboard:**
    - High-contrast HTML status table with color-coded CSS status pills without external fonts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup rbl-check`
  - Flags:
    - `--ip-address <string>`: Check specific IP address or CIDR range.
    - `--nameserver <string>`: Custom authoritative DNS resolver (default: `1.1.1.1`).
    - `--is-auto-quarantine`: Automatically isolate listed IPs from warmup pool (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type RblCheckResult struct {
        CheckID        string    `json:"checkId"`
        TargetIP       string    `json:"targetIp"`
        BlacklistHost  string    `json:"blacklistHost"`
        ResponseCode   string    `json:"responseCode"`
        CheckedAt      time.Time `json:"checkedAt"`
        IsListed       bool      `json:"isListed"`
        HasDelistLink  bool      `json:"hasDelistLink"`
        DelistURL      string    `json:"delistUrl"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RblCheckRecord (
        check_id TEXT PRIMARY KEY,
        target_ip TEXT NOT NULL,
        blacklist_host TEXT NOT NULL,
        response_code TEXT NOT NULL,
        checked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_listed INTEGER NOT NULL DEFAULT 0,
        has_delist_link INTEGER NOT NULL DEFAULT 0,
        delist_url TEXT
    );
    CREATE INDEX IF NOT EXISTS idx_rbl_target ON RblCheckRecord(target_ip, is_listed);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/warmup/reputation/rbl.go`, triggering emergency traffic rerouting in GitMap email relay balancer.

---

### Improvement 1137: SMTP MX Latency & TLS Handshake Cipher Inspector HUD
- **ID:** `SEED-CIPHER-1137`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Execute deep diagnostic probes against remote destination MX servers, measuring TCP SYN handshake latency, STARTTLS protocol version negotiation (TLS 1.2 vs 1.3), and cipher suite security strength with zero external utility dependencies.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Cipher inspector terminal card rendering certificate expiration warnings and protocol negotiation:
      ```text
      ┌─ [DESTINATION MX TLS & PROTOCOL DIAGNOSTIC HUD] ─────────────────────────────────────┐
      │ TARGET DOMAIN: gmail-smtp-in.l.google.com | PORT: 25 | RESOLVED: 142.250.150.26       │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [CONNECTION TELEMETRY]                                                               │
      │ ├─ DNS MX Resolution:     8.2ms   [OK]                                               │
      │ ├─ TCP SYN Latency:       18.4ms  [EXCELLENT]                                        │
      │ ├─ SMTP Banner Read:      24.1ms  (220 mx.google.com ESMTP ...)                      │
      │ └─ STARTTLS Negotiation:  36.8ms  [TLS 1.3 | TLS_AES_256_GCM_SHA384]                │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [TLS CERTIFICATE METADATA]                                                           │
      │ ├─ Common Name: *.google.com (SAN: 42 domains)                                      │
      │ ├─ Issuer:      GTS CA 1C3                                                           │
      │ ├─ Validity:    Expires in 74 days (2026-11-10)                                      │
      │ └─ Flags:       Strict TLS [✓] | Valid Chain [✓] | Downgrade Resilient [✓]           │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Cipher purple `#c084fc`, Metric teal `#14b8a6`, Valid green `#22c55e`.
  - **Zero-CDN Web Dashboard:**
    - Clean DOM table displaying TLS cipher breakdown and handshake timeline.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup probe-tls`
  - Flags:
    - `--domain <string>`: Target destination domain to probe (e.g. `google.com`).
    - `--min-tls-version <string>`: Required minimum TLS version (`1.2` or `1.3`; default: `1.2`).
    - `--is-strict-cert`: Reject self-signed or untrusted certificate chains (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MxTlsProbeRecord struct {
        ProbeID         string    `json:"probeId"`
        TargetDomain    string    `json:"targetDomain"`
        MxHostname      string    `json:"mxHostname"`
        TlsVersion      string    `json:"tlsVersion"`
        CipherSuite     string    `json:"cipherSuite"`
        HandshakeTimeMs int64     `json:"handshakeTimeMs"`
        CertExpiresAt   time.Time `json:"certExpiresAt"`
        IsValidCert     bool      `json:"isValidCert"`
        HasTls13        bool      `json:"hasTls13"`
        IsDowngradeSafe bool      `json:"isDowngradeSafe"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxTlsProbeRecord (
        probe_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        mx_hostname TEXT NOT NULL,
        tls_version TEXT NOT NULL,
        cipher_suite TEXT NOT NULL,
        handshake_time_ms INTEGER NOT NULL,
        cert_expires_at DATETIME NOT NULL,
        is_valid_cert INTEGER NOT NULL DEFAULT 1,
        has_tls13 INTEGER NOT NULL DEFAULT 1,
        is_downgrade_safe INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_probe_domain ON MxTlsProbeRecord(target_domain);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located at `pkg/warmup/probe/tls.go`, callable directly as an Antigravity health check utility.

---

### Improvement 1138: Interactive Seed Mailbox Credential Rotator & OAuth2 Token HUD
- **ID:** `SEED-OAUTH-1138`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Manage OAuth2 refresh tokens and app passwords for Google Workspace and Microsoft Graph seed mailboxes, rendering upcoming expiration countdowns and triggering zero-downtime background re-authentication.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual token lifetime gauge table sorted by time-to-expiration:
      ```text
      ┌─ [SEED MAILBOX OAUTH2 TOKEN & APP CREDENTIAL ROTATOR] ───────────────────────────────┐
      │ REGISTERED SEED CREDENTIALS: 96 | HEALTHY: 94 | REFRESHING: 2 | EXPIRED: 0            │
      ├────────────────────────┬─────────────┬──────────────┬───────────────┬────────────────┤
      │ MAILBOX ADDRESS        │ AUTH TYPE   │ EXPIRES IN   │ TOKEN REFRESH │ STATUS         │
      ├────────────────────────┼─────────────┼──────────────┼───────────────┼────────────────┤
      │ gsuite-seed-01@co.com  │ OAuth2.0    │ 42m 18s      │ Auto-Queued   │ [●] VALID      │
      │ o365-seed-14@org.net   │ MS Graph    │ 12h 05m      │ Standby       │ [●] VALID      │
      │ yahoo-seed-03@ya.com   │ AppPass     │ 89 days      │ Fixed Secret  │ [●] VALID      │
      │ gsuite-seed-22@co.com  │ OAuth2.0    │ 02m 11s      │ Exchanging... │ [▲] REFRESHING │
      └────────────────────────┴─────────────┴──────────────┴───────────────┴────────────────┘
      ```
    - Lipgloss styles: Active green `#22c55e`, Refresh amber `#f59e0b`, Expired red `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Embedded SVG progress meters showing credential validity decay.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup rotate-credentials`
  - Flags:
    - `--provider <string>`: Filter by provider (`google`, `microsoft`, or `all`; default: `all`).
    - `--force-refresh`: Trigger immediate renewal of OAuth2 tokens regardless of age.
    - `--is-dry-run`: Validate refresh endpoints without committing new tokens (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedCredentialStatus struct {
        CredentialID    string    `json:"credentialId"`
        MailboxAddress  string    `json:"mailboxAddress"`
        AuthType        string    `json:"authType"`
        TokenExpiresAt  time.Time `json:"tokenExpiresAt"`
        LastRefreshedAt time.Time `json:"lastRefreshedAt"`
        IsValid         bool      `json:"isValid"`
        IsAutoRenewable bool      `json:"isAutoRenewable"`
        HasSecretVault  bool      `json:"hasSecretVault"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedCredentialVault (
        credential_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL UNIQUE,
        auth_type TEXT NOT NULL,
        token_expires_at DATETIME NOT NULL,
        last_refreshed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_valid INTEGER NOT NULL DEFAULT 1,
        is_auto_renewable INTEGER NOT NULL DEFAULT 1,
        has_secret_vault INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Secure token manager at `pkg/warmup/auth/rotator.go` integrated with installation key storage in `installation.db`.

---

### Improvement 1139: Seed Mailbox Storage Quota & Purge Automation Pipeline
- **ID:** `SEED-QUOTA-1139`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Monitor seed mailbox storage utilization via IMAP QUOTA extension (RFC 2087), triggering automated pruning of non-essential warmup historical threads when mailboxes exceed 85% disk quota threshold.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Horizontal storage gauge bars rendering current megabyte consumption per mailbox:
      ```text
      ┌─ [SEED MAILBOX STORAGE QUOTA & RETENTION MONITOR] ───────────────────────────────────┐
      │ TRACKED INBOXES: 140 | OVER CAPACITY (>85%): 2 | SCHEDULED FOR PURGE: 2               │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ seed-dev-01@domain.com  [████████████████████░░░░░] 78.4% (11.7 GB / 15.0 GB)  [OK]   │
      │ seed-dev-02@domain.com  [█████████████████████████] 94.2% (14.1 GB / 15.0 GB)  [PURGE]│
      │ seed-dev-03@domain.com  [████████████░░░░░░░░░░░░░] 48.0% ( 7.2 GB / 15.0 GB)  [OK]   │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ ACTION: Purging 4,200 archived threads (>60 days old) on seed-dev-02@domain.com...  │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Safe teal `#0d9488`, Warning amber `#f59e0b`, Purge red `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Visual HTML5 meter bars indicating percentage capacity without client-side charting libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup quota-cleanup`
  - Flags:
    - `--quota-threshold <float>`: Storage percentage triggering purge (default: `85.0`).
    - `--retention-days <int>`: Delete warmup threads older than N days (default: `60`).
    - `--is-dry-run`: Calculate reclaimable storage without issuing IMAP EXPUNGE (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxQuotaUsage struct {
        MailboxID       string    `json:"mailboxId"`
        EmailAddress    string    `json:"emailAddress"`
        UsedBytes       int64     `json:"usedBytes"`
        TotalBytes      int64     `json:"totalBytes"`
        PercentUsed     float64   `json:"percentUsed"`
        LastPurgedAt    time.Time `json:"lastPurgedAt"`
        IsOverQuota     bool      `json:"isOverQuota"`
        HasPurgeEnabled bool      `json:"hasPurgeEnabled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxQuotaUsageRecord (
        mailbox_id TEXT PRIMARY KEY,
        email_address TEXT NOT NULL UNIQUE,
        used_bytes INTEGER NOT NULL,
        total_bytes INTEGER NOT NULL,
        percent_used REAL NOT NULL,
        last_purged_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_over_quota INTEGER NOT NULL DEFAULT 0,
        has_purge_enabled INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_quota_percent ON MailboxQuotaUsageRecord(percent_used);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/storage/quota.go` with automated maintenance scheduled in GitMap tasks.

---

### Improvement 1140: Warmup Feedback Loop (FBL) Abuse Complaint Interceptor HUD
- **ID:** `SEED-FBL-1140`
- **Cluster:** 114 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Ingest and parse Abuse Reporting Format (ARF / RFC 5965) complaints received from major ISP Feedback Loops (Yahoo, Comcast, Microsoft JMRP), automatically matching originating warmup campaigns and initiating instant cool-down backoff.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live complaint alert banner with originating warmup campaign attribution and immediate backoff trigger:
      ```text
      ┌─ [ARF FEEDBACK LOOP (FBL) ABUSE COMPLAINT INTERCEPTOR] ──────────────────────────────┐
      │ TOTAL COMPLAINTS: 1 | OVERALL COMPLAINT RATE: 0.008% (SAFETY THRESHOLD: < 0.100%)    │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [INCIDENT #fbl-8819]                                                                 │
      │ ├─ Provider:      Yahoo / AOL Feedback Loop (ARF v1.0)                               │
      │ ├─ Feedback Type: abuse (User clicked 'This is Spam')                                │
      │ ├─ Reported At:   14:28:10 UTC (2.4m ago)                                            │
      │ ├─ Origin Egress: 198.51.100.14 | Header Msg-ID: <wup-89104@mail.biz>              │
      │ ├─ Campaign:      "Phase 3 Enterprise Warmup Flow"                                   │
      │ └─ ACTION TAKEN:  Automated Cool-down Applied: Egress 198.51.100.14 Rate Reduced 50% │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Incident red `#ef4444`, Accent cyan `#06b6d4`, Text white `#f8fafc`.
  - **Zero-CDN Web Dashboard:**
    - Incident log table with ISP ARF raw payload viewer modal.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup fbl-telemetry`
  - Flags:
    - `--max-complaint-rate <float>`: Max acceptable complaint percentage (default: `0.05`).
    - `--is-auto-throttle`: Automatically throttle sending rate upon ARF event (default: `true`).
    - `--export-csv <path>`: Export historical complaint log to CSV.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FblComplaintIncident struct {
        IncidentID     string    `json:"incidentId"`
        ProviderName   string    `json:"providerName"`
        OriginEgressIP string    `json:"originEgressIp"`
        WarmupCampaign string    `json:"warmupCampaign"`
        OriginalMsgID  string    `json:"originalMsgId"`
        ReceivedAt     time.Time `json:"receivedAt"`
        IsThrottled    bool      `json:"isThrottled"`
        HasValidArf    bool      `json:"hasValidArf"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FblComplaintIncidentRecord (
        incident_id TEXT PRIMARY KEY,
        provider_name TEXT NOT NULL,
        origin_egress_ip TEXT NOT NULL,
        warmup_campaign TEXT NOT NULL,
        original_msg_id TEXT NOT NULL,
        received_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_throttled INTEGER NOT NULL DEFAULT 1,
        has_valid_arf INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_fbl_egress ON FblComplaintIncidentRecord(origin_egress_ip, received_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located at `pkg/warmup/abuse/fbl.go`, wired into the dynamic rate limiter to enforce real-time deliverability protection.


# Cluster 115: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 1141–1150)

---

### Improvement 1141: Adaptive SMTP Backoff Inversion & ISP Throttle Decompressor HUD
- **ID:** `THROTTLE-BACKOFF-1141`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Dynamically invert SMTP 421 and 451 rate-limiting backoff curves using feedback-driven cubic delay dampening (TCP CUBIC inspired) rather than naive exponential binary backoff, avoiding aggressive ISP lockouts and maximizing throughput recovery.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live concurrency and backoff throttle visualization showing dynamic pacing curves:
      ```text
      ┌─ [ADAPTIVE SMTP BACKOFF INVERSION & DELAY DAMPENER] ─────────────────────────────────┐
      │ DOMAIN: mx.google.com | CURRENT RATE: 14 msg/min | BACKOFF FACTOR: 0.72x (CUBIC)     │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [CONCURRENCY FLOW CONTROL]                                                           │
      │ Allocated Slots:   [████████████░░░░░░░░] 6 / 10 active connections                  │
      │ Window Size (Wmax): 12 concurrency units | Epoch Duration: 45.2s                      │
      │ Throttle Events:   1x 421 4.7.0 (Too many concurrent connections) -> Smooth Recovery│
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [PACE DECOMPRESSION CURVE]                                                           │
      │ 24 msg/m ┤         ┌───────┐                                                         │
      │ 16 msg/m ┤   ┌─────┘       └─────┐                                                   │
      │  8 msg/m ┼───┘                   └────────── [Target Equilibrium: 18 msg/m]          │
      │          └────────────────────────────────────────────────────────► Time             │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Curve cyan `#06b6d4`, Warning amber `#f59e0b`, Equilibrium green `#10b981`.
  - **Zero-CDN Web Dashboard:**
    - SVG stream rendering cubic recovery trajectory without external charting libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup backoff-status`
  - Flags:
    - `--domain <string>`: Filter status by remote destination domain (default: `all`).
    - `--is-cubic-enabled`: Enforce CUBIC delay progression over binary exponential (default: `true`).
    - `--min-backoff-sec <float>`: Floor backoff delay after 4xx response (default: `1.5`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BackoffInversionState struct {
        DomainName      string    `json:"domainName"`
        CurrentWindow   int       `json:"currentWindow"`
        MaxWindow       int       `json:"maxWindow"`
        BackoffFactor   float64   `json:"backoffFactor"`
        LastThrottleAt  time.Time `json:"lastThrottleAt"`
        IsThrottled     bool      `json:"isThrottled"`
        IsRecovering    bool      `json:"isRecovering"`
        HasCubicEnabled bool      `json:"hasCubicEnabled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpBackoffInversionRecord (
        domain_name TEXT PRIMARY KEY,
        current_window INTEGER NOT NULL DEFAULT 1,
        max_window INTEGER NOT NULL DEFAULT 10,
        backoff_factor REAL NOT NULL DEFAULT 1.0,
        last_throttle_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        is_recovering INTEGER NOT NULL DEFAULT 1,
        has_cubic_enabled INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_backoff_status ON SmtpBackoffInversionRecord(is_throttled);
    ```
- **Antigravity / GitMap Integration Point:**
  - Engine in `pkg/warmup/throttler/backoff.go`, dynamically adjusting sending scheduler in GitMap delivery engine.

---

### Improvement 1142: Multi-Tenant MX Leaky-Bucket Rate Governor & Jitter Controller HUD
- **ID:** `THROTTLE-LEAKY-1142`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Implement a distributed leaky-bucket traffic shaper with cryptographic micro-jitter across multiple sending tenants sharing egress IP pools, preventing burst synchronization and ISP rate-limit triggers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual visual bucket capacity meters depicting leak rate versus incoming token velocity:
      ```text
      ┌─ [MULTI-TENANT MX LEAKY-BUCKET FLOW GOVERNOR] ──────────────────────────────────────┐
      │ TARGET MX POOL: Microsoft 365 (protection.outlook.com) | CAPACITY: 120 msg/hr        │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [BUCKET LEVEL: 42 / 120 TOKENS]                                                      │
      │ Fill State:  [████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░] 35% Full (Safe Zone)       │
      │ Leak Rate:   1 token / 30.0s | Micro-Jitter: ±2.4s (Cryptographic Uniform Random)    │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [TENANT ALLOCATION BREAKDOWN]                                                        │
      │ ├─ Tenant A (Alpha Outreach):     18 tokens/hr [15%]                                 │
      │ ├─ Tenant B (Transactional Rel):  48 tokens/hr [40%]                                 │
      │ ├─ Tenant C (Warmup Seed Agent):  36 tokens/hr [30%]                                 │
      │ └─ Unallocated Buffer:            18 tokens/hr [15%]                                 │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Safe green `#22c55e`, Warning orange `#f97316`, Metric blue `#3b82f6`.
  - **Zero-CDN Web Dashboard:**
    - Real-time SVG bucket diagram depicting water level and leak valve velocity.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup bucket-governor`
  - Flags:
    - `--mx-host <string>`: Target MX hostname or provider group (default: `all`).
    - `--jitter-range-sec <float>`: Max pseudo-random delay variance (default: `3.0`).
    - `--is-strict-quota`: Deny queue entries if bucket capacity exceeds 95% (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LeakyBucketRecord struct {
        BucketID        string    `json:"bucketId"`
        MxTargetHost    string    `json:"mxTargetHost"`
        Capacity        int       `json:"capacity"`
        CurrentTokens   int       `json:"currentTokens"`
        LeakRatePerMin  float64   `json:"leakRatePerMin"`
        LastLeakedAt    time.Time `json:"lastLeakedAt"`
        IsOverflown     bool      `json:"isOverflown"`
        HasJitterActive bool      `json:"hasJitterActive"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxLeakyBucketRecord (
        bucket_id TEXT PRIMARY KEY,
        mx_target_host TEXT NOT NULL UNIQUE,
        capacity INTEGER NOT NULL DEFAULT 100,
        current_tokens INTEGER NOT NULL DEFAULT 0,
        leak_rate_per_min REAL NOT NULL DEFAULT 2.0,
        last_leaked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_overflown INTEGER NOT NULL DEFAULT 0,
        has_jitter_active INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Core governor at `pkg/warmup/governor/leaky.go`, enforced in the send pipeline before socket acquisition.

---

### Improvement 1143: Destination MX MX-Priority Balancing & Fallback Topology HUD
- **ID:** `THROTTLE-TOPOLOGY-1143`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Analyze multi-homed destination MX DNS records, actively balancing traffic across equal-preference MX hosts, detecting regional route degradations, and routing around flapping mail exchangers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Multi-priority MX routing tree showing health and route weights:
      ```text
      ┌─ [DESTINATION MX MULTI-HOMED ROUTING TOPOLOGY: domain.tld] ──────────────────────────┐
      │ DISCOVERED MX EXCHANGERS: 4 | ACTIVE PATHS: 3 | DEGRADED: 1                          │
      ├─────┬────────────────────────────────┬──────┬──────────┬───────────┬─────────────────┤
      │ PREF│ MX HOSTNAME                    │ IP   │ RTT (ms) │ FLOW WGT  │ ROUTE STATUS    │
      ├─────┼────────────────────────────────┼──────┼──────────┼───────────┼─────────────────┤
      │  5  │ aspmx.l.google.com             │ .26  │  14.2ms  │    50%    │ [●] HEALTHY     │
      │  5  │ alt1.aspmx.l.google.com        │ .27  │  18.9ms  │    50%    │ [●] HEALTHY     │
      │ 10  │ alt2.aspmx.l.google.com        │ .27  │  42.1ms  │   STANDBY │ [○] IDLE        │
      │ 10  │ alt3.aspmx.l.google.com        │ .28  │ 290.4ms  │  DEGRADED │ [▲] LATENCY LAG │
      └─────┴────────────────────────────────┴──────┴──────────┴───────────┴─────────────────┘
      ```
    - Lipgloss styles: Healthy green `#22c55e`, Standby slate `#94a3b8`, Degraded amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Pure CSS flowchart depicting primary versus fallback MX routing nodes.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mx-topology`
  - Flags:
    - `--domain <string>`: Destination domain to resolve and map.
    - `--max-rtt-ms <int>`: Max round-trip time threshold before demoting route (default: `150`).
    - `--is-auto-failover`: Automatically switch to secondary MX upon TCP timeout (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MxRouteNode struct {
        RouteID         string    `json:"routeId"`
        TargetDomain    string    `json:"targetDomain"`
        MxHost          string    `json:"mxHost"`
        Preference      int       `json:"preference"`
        LatencyMs       float64   `json:"latencyMs"`
        AssignedWeight  int       `json:"assignedWeight"`
        LastEvaluatedAt time.Time `json:"lastEvaluatedAt"`
        IsHealthy       bool      `json:"isHealthy"`
        IsPrimaryRoute  bool      `json:"isPrimaryRoute"`
        HasFallbackLock bool      `json:"hasFallbackLock"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxRouteTopologyRecord (
        route_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        mx_host TEXT NOT NULL,
        preference INTEGER NOT NULL,
        latency_ms REAL NOT NULL DEFAULT 0.0,
        assigned_weight INTEGER NOT NULL DEFAULT 100,
        last_evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_healthy INTEGER NOT NULL DEFAULT 1,
        is_primary_route INTEGER NOT NULL DEFAULT 1,
        has_fallback_lock INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_mx_domain ON MxRouteTopologyRecord(target_domain, preference);
    ```
- **Antigravity / GitMap Integration Point:**
  - Module in `pkg/warmup/dns/mxtopology.go`, integrated into the transport dialing layer.

---

### Improvement 1144: Dynamic Sending IP Warmup Stage Transitioner & Ramp Curve HUD
- **ID:** `THROTTLE-RAMP-1144`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Govern the multi-week warmup ramping curve for fresh IP addresses and domains, enforcing strict daily volume caps, step-function progression (Day 1: 20 msg -> Day 30: 10,000 msg), and automated holdback when placement drops below 95%.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Interactive 30-day stage stepper with daily sending quota progress bars:
      ```text
      ┌─ [IP & DOMAIN WARMUP RAMP SCHEDULE: out-node-04 (198.51.100.22)] ────────────────────┐
      │ WARMUP PHASE: Day 12 / 30 | TARGET ESP: Google Workspace | REPUTATION SCORE: 98/100  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [DAILY VOLUME ENVELOPE: DAY 12]                                                      │
      │ Sent Today:   [████████████████████████░░░░░░] 380 / 480 messages (79.1%)            │
      │ Remaining:    100 messages | Hourly Pacing: ~18 msg/hour                             │
      │ Stage State:  PROGRESSING NORMALLY (Inbox Rate: 98.4% | Spam Rate: 1.6%)             │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [RAMP PREVIEW]                                                                       │
      │ Day 10: 320 msg [✓] ── Day 11: 400 msg [✓] ── Day 12: 480 msg [CURRENT]             │
      │ ── Day 13: 600 msg [NEXT] ── Day 14: 750 msg [PLANNED]                              │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Active blue `#3b82f6`, Done green `#10b981`, Future slate `#64748b`.
  - **Zero-CDN Web Dashboard:**
    - Line chart of planned vs actual daily sending volumes rendered with pure inline SVG.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ramp-status`
  - Flags:
    - `--ip <string>`: Outbound IP address to inspect.
    - `--phase-day <int>`: Manually advance or freeze warmup day schedule.
    - `--is-auto-freeze`: Halt ramp progression automatically if bounce rate exceeds 2.0% (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type WarmupRampSchedule struct {
        ScheduleID      string    `json:"scheduleId"`
        EgressIP        string    `json:"egressIp"`
        CurrentDayIndex int       `json:"currentDayIndex"`
        DailyLimit      int       `json:"dailyLimit"`
        SentToday       int       `json:"sentToday"`
        PlacementRatio  float64   `json:"placementRatio"`
        LastUpdated     time.Time `json:"lastUpdated"`
        IsStageFrozen   bool      `json:"isStageFrozen"`
        IsComplete      bool      `json:"isComplete"`
        HasRampApproved bool      `json:"hasRampApproved"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupRampScheduleRecord (
        schedule_id TEXT PRIMARY KEY,
        egress_ip TEXT NOT NULL UNIQUE,
        current_day_index INTEGER NOT NULL DEFAULT 1,
        daily_limit INTEGER NOT NULL DEFAULT 50,
        sent_today INTEGER NOT NULL DEFAULT 0,
        placement_ratio REAL NOT NULL DEFAULT 1.0,
        last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_stage_frozen INTEGER NOT NULL DEFAULT 0,
        is_complete INTEGER NOT NULL DEFAULT 0,
        has_ramp_approved INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_ramp_ip ON WarmupRampScheduleRecord(egress_ip);
    ```
- **Antigravity / GitMap Integration Point:**
  - Logic residing in `pkg/warmup/schedule/ramp.go`, gating all outbound delivery jobs.

---

### Improvement 1145: Concurrent SMTP Socket Connection Pooler & Connection Reuse HUD
- **ID:** `THROTTLE-POOL-1145`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Maximize outbound SMTP efficiency via transparent TCP/TLS connection pooling (RFC 5321 section 4.1.4 pipelining and reuse), maintaining warm keep-alive sockets per destination MX and eliminating redundant TLS handshakes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live socket connection pool HUD showing active, idle, and recycling channels:
      ```text
      ┌─ [SMTP TCP/TLS CONNECTION POOL MANAGER] ─────────────────────────────────────────────┐
      │ ACTIVE POOLS: 12 | TOTAL OPEN SOCKETS: 48 | REUSE EFFICIENCY: 84.2%                  │
      ├──────────────────────────────┬────────┬────────┬────────┬─────────────┬──────────────┤
      │ DESTINATION MX               │ ACTIVE │ IDLE   │ MAX    │ HANDSHAKE SAV│ STATUS      │
      ├──────────────────────────────┼────────┼────────┼────────┼─────────────┼──────────────┤
      │ gmail-smtp-in.l.google.com   │     12 │      4 │     16 │     420 ms  │ [●] SATURATED│
      │ mail.protection.outlook.com  │      8 │      6 │     14 │     380 ms  │ [●] OPTIMAL  │
      │ mx-aol.mail.gm0.yahoodns.net │      4 │      2 │      8 │     510 ms  │ [●] OPTIMAL  │
      │ mail.protonmail.ch           │      2 │      2 │      4 │     290 ms  │ [●] OPTIMAL  │
      └──────────────────────────────┴────────┴────────┴────────┴─────────────┴──────────────┘
      ```
    - Lipgloss styles: Pool green `#22c55e`, Saturated amber `#f59e0b`, Savings cyan `#06b6d4`.
  - **Zero-CDN Web Dashboard:**
    - Visual connection slots rendered as grid cells with CSS state transitions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup socket-pool`
  - Flags:
    - `--max-idle-conns <int>`: Max idle sockets per destination host (default: `8`).
    - `--idle-timeout <duration>`: Time before recycling idle connection (default: `45s`).
    - `--is-pipelining-enabled`: Use SMTP pipelining to batch commands (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SmtpSocketPoolNode struct {
        PoolID          string    `json:"poolId"`
        DestinationMx   string    `json:"destinationMx"`
        ActiveSockets   int       `json:"activeSockets"`
        IdleSockets     int       `json:"idleSockets"`
        MaxCapacity     int       `json:"maxCapacity"`
        TotalReuses     int64     `json:"totalReuses"`
        LastRecycledAt  time.Time `json:"lastRecycledAt"`
        IsPipelined     bool      `json:"isPipelined"`
        IsFull          bool      `json:"isFull"`
        HasKeepAlive    bool      `json:"hasKeepAlive"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpSocketPoolRecord (
        pool_id TEXT PRIMARY KEY,
        destination_mx TEXT NOT NULL UNIQUE,
        active_sockets INTEGER NOT NULL DEFAULT 0,
        idle_sockets INTEGER NOT NULL DEFAULT 0,
        max_capacity INTEGER NOT NULL DEFAULT 16,
        total_reuses INTEGER NOT NULL DEFAULT 0,
        last_recycled_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_pipelined INTEGER NOT NULL DEFAULT 1,
        is_full INTEGER NOT NULL DEFAULT 0,
        has_keep_alive INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Packaged at `pkg/warmup/transport/pool.go` serving as the base dialer for email operations.

---

### Improvement 1146: Asynchronous Egress Interface Cycling & MAC/IP Binding HUD
- **ID:** `THROTTLE-EGRESS-1146`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Round-robin or least-loaded bind local outbound SMTP sockets to designated secondary IP interfaces on multi-homed Linux/Windows servers, dynamically rotating egress routes to distribute connection volume across physical interfaces.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live hardware interface table rendering throughput and active bindings:
      ```text
      ┌─ [EGRESS NETWORK INTERFACE BINDING & TRAFFIC DISTRIBUTION] ──────────────────────────┐
      │ SYSTEM INTERFACES: 4 | ACTIVE OUTBOUND IPs: 8 | CURRENT SCHEDULER: LEAST-CONNECTIONS │
      ├──────────┬─────────────────┬──────────┬─────────────┬─────────────┬──────────────────┤
      │ IFACE    │ IP ADDRESS      │ ACTIVE   │ TOTAL SENT  │ ERROR RATE  │ HEALTH           │
      ├──────────┼─────────────────┼──────────┼─────────────┼─────────────┼──────────────────┤
      │ eth0:1   │ 198.51.100.10   │        4 │       1,420 │       0.0%  │ [●] PRISTINE     │
      │ eth0:2   │ 198.51.100.11   │        5 │       1,380 │       0.1%  │ [●] PRISTINE     │
      │ eth0:3   │ 198.51.100.12   │        3 │       1,410 │       0.0%  │ [●] PRISTINE     │
      │ eth1:0   │ 203.0.113.50    │        2 │         890 │       0.0%  │ [●] PRISTINE     │
      └──────────┴─────────────────┴──────────┴─────────────┴─────────────┴──────────────────┘
      ```
    - Lipgloss styles: Pristine green `#22c55e`, High-load amber `#f59e0b`, Border slate `#475569`.
  - **Zero-CDN Web Dashboard:**
    - Egress interface distribution pie/donut diagram rendered with pure SVG geometry.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup egress-interfaces`
  - Flags:
    - `--scheduler <string>`: Algorithm (`round-robin`, `least-conns`, `hash`; default: `least-conns`).
    - `--exclude-interface <string>`: Interface name to omit from binding pool.
    - `--is-auto-drain`: Drain traffic from failing interfaces before complete removal (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EgressInterfaceNode struct {
        InterfaceID     string    `json:"interfaceId"`
        InterfaceName   string    `json:"interfaceName"`
        IPAddress       string    `json:"ipAddress"`
        ActiveSockets   int       `json:"activeSockets"`
        TotalDelivered  int64     `json:"totalDelivered"`
        FailureCount    int       `json:"failureCount"`
        LastBoundAt     time.Time `json:"lastBoundAt"`
        IsOnline        bool      `json:"isOnline"`
        IsDraining      bool      `json:"isDraining"`
        HasIpv6Enabled  bool      `json:"hasIpv6Enabled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EgressInterfaceRecord (
        interface_id TEXT PRIMARY KEY,
        interface_name TEXT NOT NULL,
        ip_address TEXT NOT NULL UNIQUE,
        active_sockets INTEGER NOT NULL DEFAULT 0,
        total_delivered INTEGER NOT NULL DEFAULT 0,
        failure_count INTEGER NOT NULL DEFAULT 0,
        last_bound_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_online INTEGER NOT NULL DEFAULT 1,
        is_draining INTEGER NOT NULL DEFAULT 0,
        has_ipv6_enabled INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Dialer binding logic at `pkg/warmup/transport/egress.go`, wrapping `net.Dialer.LocalAddr`.

---

### Improvement 1147: Transient DNS Resolver Latency Dampener & TTL Shield HUD
- **ID:** `THROTTLE-DNS-1147`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** In-memory DNS caching layer with pre-emptive asynchronous background refreshing before TTL expiration (TTL Shielding), eliminating DNS lookup stalls during high-volume warmup bursts.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Cache hit ratio HUD and TTL decay countdown:
      ```text
      ┌─ [TRANSIENT DNS TTL SHIELD & ASYNC PRE-RESOLVER] ────────────────────────────────────┐
      │ CACHE ENTRIES: 64 | HIT RATIO: 99.4% | AVG RESOLUTION TIME: 0.12ms (CACHE HIT)       │
      ├──────────────────────────────┬──────────┬───────────┬──────────────┬─────────────────┤
      │ DOMAIN / QUERY               │ RECORD   │ TTL LEFT  │ PRE-FETCH    │ STATUS          │
      ├──────────────────────────────┼──────────┼───────────┼──────────────┼─────────────────┤
      │ gmail-smtp-in.l.google.com   │ A        │ 4m 12s    │ Queued @ 30s │ [●] WARM        │
      │ protection.outlook.com       │ MX       │ 22m 04s   │ Standby      │ [●] WARM        │
      │ yahoo.com                    │ TXT(SPF) │ 52m 19s   │ Standby      │ [●] WARM        │
      │ protonmail.ch                │ MX       │ 0m 18s    │ Refreshing...│ [▲] PREFETCHING │
      └──────────────────────────────┴──────────┴───────────┴──────────────┴─────────────────┘
      ```
    - Lipgloss styles: Warm green `#22c55e`, Prefetch cyan `#06b6d4`, Cache violet `#8b5cf6`.
  - **Zero-CDN Web Dashboard:**
    - Live gauge chart of DNS query latency comparing local cache hits against upstream authoritatives.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup dns-cache`
  - Flags:
    - `--cache-size <int>`: Max DNS entries in memory (default: `1024`).
    - `--prefetch-window <duration>`: Refresh records when TTL drops below window (default: `45s`).
    - `--is-stale-permitted`: Serve slightly stale record during upstream resolver outages (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DnsCacheEntry struct {
        EntryID         string    `json:"entryId"`
        DomainName      string    `json:"domainName"`
        RecordType      string    `json:"recordType"`
        ResolvedData    string    `json:"resolvedData"`
        TtlSeconds      int       `json:"ttlSeconds"`
        CachedAt        time.Time `json:"cachedAt"`
        ExpiresAt       time.Time `json:"expiresAt"`
        IsPrefetching   bool      `json:"isPrefetching"`
        IsStaleFallback bool      `json:"isStaleFallback"`
        HasUpstreamErr  bool      `json:"hasUpstreamErr"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DnsCacheEntryRecord (
        entry_id TEXT PRIMARY KEY,
        domain_name TEXT NOT NULL,
        record_type TEXT NOT NULL,
        resolved_data TEXT NOT NULL,
        ttl_seconds INTEGER NOT NULL,
        cached_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at DATETIME NOT NULL,
        is_prefetching INTEGER NOT NULL DEFAULT 0,
        is_stale_fallback INTEGER NOT NULL DEFAULT 0,
        has_upstream_err INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_dns_cache_query ON DnsCacheEntryRecord(domain_name, record_type);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/dns/cache.go`, integrated transparently into custom Go resolver.

---

### Improvement 1148: Destination ISP Graylist Trapper & Retry Scheduling Matrix
- **ID:** `THROTTLE-GRAYLIST-1148`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Heuristically identify SMTP 450/451 Greylisting responses ("Recipient address rejected: Greylisted"), automatically profiling the remote MTA's minimum delay requirement (typically 5 to 15 minutes), and scheduling targeted resends to bypass greylist filters.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual timeline of greylisted delivery retries:
      ```text
      ┌─ [DESTINATION ISP GREYLIST ADAPTIVE RETRY MATRIX] ───────────────────────────────────┐
      │ ACTIVE GREYLIST INCIDENTS: 3 | LEARNED PROFILES: 18 | AUTO-RESEND SUCCESS: 100%      │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [INCIDENT: seed-cust-18@corporate-firewall.net]                                      │
      │ ├─ Initial Attempt: 14:10:00 -> 451 4.7.1 Greylisted, please try again in 300s       │
      │ ├─ Profile Learned:  Minimum Greylist Penalty = 300s (5.0 minutes)                   │
      │ ├─ Reschedule Time:  14:15:05 (Buffered with +5s safety margin)                      │
      │ └─ Final Outcome:    DELIVERED ON RETRY [250 2.0.0 Ok: queued as 8B49F1]            │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Incident amber `#f59e0b`, Delivered green `#22c55e`, Text white `#f8fafc`.
  - **Zero-CDN Web Dashboard:**
    - Clean tabular schedule of deferred deliveries awaiting greylist window expiration.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup greylist-status`
  - Flags:
    - `--destination-domain <string>`: Target domain to query learned greylist profiles.
    - `--default-retry-delay <duration>`: Fallback delay when unstated by MTA (default: `8m`).
    - `--is-auto-reschedule`: Auto-requeue greylisted payloads with learned delay (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GreylistProfileRecord struct {
        ProfileID       string    `json:"profileId"`
        TargetDomain    string    `json:"targetDomain"`
        MinDelaySeconds int       `json:"minDelaySeconds"`
        SuccessCount    int       `json:"successCount"`
        LastUpdated     time.Time `json:"lastUpdated"`
        IsConfirmed     bool      `json:"isConfirmed"`
        IsBypassed      bool      `json:"isBypassed"`
        HasCustomMargin bool      `json:"hasCustomMargin"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GreylistProfileRecord (
        profile_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL UNIQUE,
        min_delay_seconds INTEGER NOT NULL DEFAULT 300,
        success_count INTEGER NOT NULL DEFAULT 0,
        last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_confirmed INTEGER NOT NULL DEFAULT 1,
        is_bypassed INTEGER NOT NULL DEFAULT 0,
        has_custom_margin INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Packaged in `pkg/warmup/throttler/greylist.go`, hooked directly into SMTP delivery error classification.

---

### Improvement 1149: Cross-Tenant Warmup Concurrency Semaphore & Fairness Arbiter HUD
- **ID:** `THROTTLE-ARBITER-1149`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Enforce fair-share scheduling across multiple concurrent warmup campaigns and mailbox tiers using a weighted round-robin concurrency semaphore, preventing high-volume campaigns from starving critical low-volume identity warming threads.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live multi-queue bandwidth distribution dashboard:
      ```text
      ┌─ [WEIGHTED CONCURRENCY SEMAPHORE & FAIRNESS ARBITER] ────────────────────────────────┐
      │ GLOBAL CONCURRENCY CAP: 32 WORKERS | ACQUIRED: 24 | IDLE: 8 | STARVATION: 0.00%     │
      ├──────────────────────┬────────┬────────┬────────┬──────────────┬─────────────────────┤
      │ CAMPAIGN TIER        │ WEIGHT │ ACTIVE │ QUEUED │ FAIR-SHARE % │ STATUS              │
      ├──────────────────────┼────────┼────────┼────────┼──────────────┼─────────────────────┤
      │ Executive VIP Warmup │    40  │     12 │      0 │        50.0% │ [●] RUNNING         │
      │ Standard Sales Dev   │    30  │      8 │      4 │        33.3% │ [●] RUNNING         │
      │ Cold Seed Placement  │    20  │      4 │     12 │        16.7% │ [●] GOVERNED        │
      │ Bulk Verification    │    10  │      0 │     40 │         0.0% │ [▲] PAUSED (BACKOFF)│
      └──────────────────────┴────────┴────────┴────────┴──────────────┴─────────────────────┘
      ```
    - Lipgloss styles: Active green `#22c55e`, Governed cyan `#06b6d4`, Paused amber `#f59e0b`.
  - **Zero-CDN Web Dashboard:**
    - Stacked horizontal bar chart showing instantaneous worker allocation per campaign tier.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup fairness-status`
  - Flags:
    - `--total-concurrency <int>`: Global max worker pool ceiling (default: `32`).
    - `--starvation-timeout <duration>`: Max queue wait time before dynamic weight boost (default: `2m`).
    - `--is-strict-fifo`: Strictly order sends within same priority tier (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FairShareQueueTier struct {
        TierID          string    `json:"tierId"`
        TierName        string    `json:"tierName"`
        PriorityWeight  int       `json:"priorityWeight"`
        ActiveWorkers   int       `json:"activeWorkers"`
        QueuedJobs      int       `json:"queuedJobs"`
        AllocatedShare  float64   `json:"allocatedShare"`
        LastScheduledAt time.Time `json:"lastScheduledAt"`
        IsActive        bool      `json:"isActive"`
        IsStarved       bool      `json:"isStarved"`
        HasPriorityLock bool      `json:"hasPriorityLock"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FairShareQueueTierRecord (
        tier_id TEXT PRIMARY KEY,
        tier_name TEXT NOT NULL UNIQUE,
        priority_weight INTEGER NOT NULL DEFAULT 10,
        active_workers INTEGER NOT NULL DEFAULT 0,
        queued_jobs INTEGER NOT NULL DEFAULT 0,
        allocated_share REAL NOT NULL DEFAULT 0.0,
        last_scheduled_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_active INTEGER NOT NULL DEFAULT 1,
        is_starved INTEGER NOT NULL DEFAULT 0,
        has_priority_lock INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located at `pkg/warmup/governor/arbiter.go`, coordinating dispatch across worker subroutines.

---

### Improvement 1150: Outbound SMTP Socket Micro-Pacing & Inter-Packet Gap Controller
- **ID:** `THROTTLE-PACING-1150`
- **Cluster:** 115 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Regulate raw TCP byte streaming pacing over SMTP sockets, introducing microsecond inter-command pauses (between `MAIL FROM`, `RCPT TO`, and `DATA` chunks) to mimic human MUA transmission cadences and evade deep packet inspection spam traps.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Microsecond oscilloscope timeline rendering command transaction delays:
      ```text
      ┌─ [SMTP RAW PROTOCOL MICRO-PACING CONTROLLER] ────────────────────────────────────────┐
      │ ACTIVE TRANSMISSIONS: 6 | CADENCE PROFILE: NATURAL MUA | INTER-COMMAND JITTER: ACTIVE│
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [TCP PACKET STREAM TELEMETRY: conn-7819 -> mx.protection.outlook.com]               │
      │ ├─ CONNECT -> 220 Banner:               [  0.0ms]                                    │
      │ ├─ [PAUSE: 142ms] -> EHLO mail.co:     [142.1ms] (Natural Handshake Delay)           │
      │ ├─ [PAUSE:  88ms] -> STARTTLS:         [230.2ms]                                     │
      │ ├─ [PAUSE: 210ms] -> MAIL FROM:        [440.4ms]                                     │
      │ ├─ [PAUSE: 175ms] -> RCPT TO:          [615.5ms]                                     │
      │ └─ [PAUSE: 310ms] -> DATA Stream:      [925.7ms] -> 250 OK Delivered                 │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Packet green `#22c55e`, Pause cyan `#38bdf8`, Delim slate `#64748b`.
  - **Zero-CDN Web Dashboard:**
    - High-precision waterfall chart rendering micro-delay phase breakdowns.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup micro-pacing`
  - Flags:
    - `--profile <string>`: Transmission profile (`natural`, `fast`, `conservative`; default: `natural`).
    - `--min-command-delay-ms <int>`: Minimum pause between SMTP verbs (default: `50`).
    - `--is-jitter-randomized`: Apply Gaussian noise to command pacing delays (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MicroPacingConfig struct {
        ProfileName       string    `json:"profileName"`
        MinDelayMs        int       `json:"minDelayMs"`
        MaxDelayMs        int       `json:"maxDelayMs"`
        GaussianMeanMs    float64   `json:"gaussianMeanMs"`
        LastAppliedAt     time.Time `json:"lastAppliedAt"`
        IsJitterActive    bool      `json:"isJitterActive"`
        IsMuaEmulated     bool      `json:"isMuaEmulated"`
        HasAdaptiveTiming bool      `json:"hasAdaptiveTiming"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MicroPacingConfigRecord (
        profile_name TEXT PRIMARY KEY,
        min_delay_ms INTEGER NOT NULL DEFAULT 50,
        max_delay_ms INTEGER NOT NULL DEFAULT 250,
        gaussian_mean_ms REAL NOT NULL DEFAULT 120.0,
        last_applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_jitter_active INTEGER NOT NULL DEFAULT 1,
        is_mua_emulated INTEGER NOT NULL DEFAULT 1,
        has_adaptive_timing INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented at `pkg/warmup/transport/pacer.go`, wrapping `net.Conn` with custom rate-limiting read/write buffers.


# Cluster 116: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 1151–1160)

---

### Improvement 1151: Self-Contained Zero-CDN Email DOM Renderer & CSP Sandbox HUD
- **ID:** `CONTENT-RENDER-1151`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Render complex responsive HTML email previews inside an ultra-isolated local sandbox with a strict `Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'` header, preventing accidental external tracking pixel fires, malicious JavaScript execution, or remote asset leaks.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Terminal HTML structure inspector and sanitizer report:
      ```text
      ┌─ [ZERO-CDN EMAIL CONTENT SECURITY SANDBOX & DOM INSPECTOR] ──────────────────────────┐
      │ PAYLOAD ID: #tpl-8941 | DOM NODES: 142 | EXTERNAL ASSETS BLOCKED: 6 | CSP: ENFORCED  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [SECURITY POLICY VIOLATION INTERCEPTIONS]                                            │
      │ ├─ BLOCKED SCRIPT: <script src="https://cdn.malicious.io/track.js">   [STRIPPED]     │
      │ ├─ INLINE STYLES:  Passed CSS Sanitizer (Safe properties only)        [ALLOWED]      │
      │ ├─ REMOTE IMAGES:  Replaced with local SVG placeholder data URIs      [ISOLATED]     │
      │ └─ TRACKING PIXEL: 1x1 GIF to http://spy.tracker.biz/pixel.gif        [QUARANTINED]  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [TERMINAL PREVIEW SUMMARY]                                                           │
      │ Subject:  "Exclusive Q3 Product Engineering Architecture Briefing"                  │
      │ Snippet:  "Welcome back! In this edition we cover zero-latency SMTP streaming..."   │
      │ Actions:  [p: Preview in Local Browser (CSP Sandboxed)] [s: View Clean HTML] [q: Quit]│
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Sandbox purple `#a855f7`, Stripped red `#ef4444`, Allowed green `#22c55e`.
  - **Zero-CDN Web Dashboard:**
    - Hardened `iframe` utilizing HTML5 `sandbox="allow-same-origin"` with zero external resources.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content render-sandbox`
  - Flags:
    - `--file <path>`: Path to local HTML template or raw MIME message.
    - `--strict-csp`: Enforce total ban on external resource loading (default: `true`).
    - `--is-interactive`: Launch interactive terminal previewer (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ContentSandboxResult struct {
        ResultID        string    `json:"resultId"`
        TemplateHash    string    `json:"templateHash"`
        TotalNodes      int       `json:"totalNodes"`
        BlockedScripts  int       `json:"blockedScripts"`
        BlockedPixels   int       `json:"blockedPixels"`
        SanitizedSize   int64     `json:"sanitizedSize"`
        ProcessedAt     time.Time `json:"processedAt"`
        IsCspCompliant  bool      `json:"isCspCompliant"`
        IsHtmlClean     bool      `json:"isHtmlClean"`
        HasRemoteAssets bool      `json:"hasRemoteAssets"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ContentSandboxRecord (
        result_id TEXT PRIMARY KEY,
        template_hash TEXT NOT NULL,
        total_nodes INTEGER NOT NULL,
        blocked_scripts INTEGER NOT NULL DEFAULT 0,
        blocked_pixels INTEGER NOT NULL DEFAULT 0,
        sanitized_size INTEGER NOT NULL,
        processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_csp_compliant INTEGER NOT NULL DEFAULT 1,
        is_html_clean INTEGER NOT NULL DEFAULT 1,
        has_remote_assets INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_sandbox_hash ON ContentSandboxRecord(template_hash);
    ```
- **Antigravity / GitMap Integration Point:**
  - Placed in `pkg/email/content/sandbox.go`, integrated into template compiler pipeline.

---

### Improvement 1152: Responsive HTML Email Viewport Matrix & Media-Query Simulator
- **ID:** `CONTENT-VIEWPORT-1152`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Autonomously test email HTML across standard viewport widths (Mobile 320px/375px/414px, Tablet 768px, Desktop 1200px), verifying CSS media-query breakpoint triggering, table fluidity, and font scaling without third-party browser emulators.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Side-by-side terminal rendering of multi-device layout geometry:
      ```text
      ┌─ [EMAIL VIEWPORT BREAKPOINT MATRIX & GEOMETRY VALIDATOR] ────────────────────────────┐
      │ TEMPLATE: newsletter-fluid-v2.html | VIEWPORTS TESTED: 4 | OVERFLOW DETECTIONS: 0    │
      ├────────────────┬──────────┬─────────────┬──────────────┬──────────────┬─────────────┤
      │ TARGET DEVICE  │ WIDTH    │ CONTAINER W │ HORIZ SCROLL │ WRAP BREAKS  │ RENDER EVAL │
      ├────────────────┼──────────┼─────────────┼──────────────┼──────────────┼─────────────┤
      │ iPhone SE      │ 320 px   │  100% fluid │ NONE [0px]   │ 2 lines OK   │ [●] PRISTINE│
      │ iPhone 15 Pro  │ 393 px   │  100% fluid │ NONE [0px]   │ 1 line OK    │ [●] PRISTINE│
      │ iPad Mini      │ 768 px   │  600 px fix │ NONE [0px]   │ Column Stack │ [●] PRISTINE│
      │ Desktop Large  │ 1200 px  │  600 px fix │ NONE [0px]   │ 2-Col Grid   │ [●] PRISTINE│
      └────────────────┴──────────┴─────────────┴──────────────┴──────────────┴─────────────┘
      ```
    - Lipgloss styles: Device teal `#0d9488`, Pristine green `#22c55e`, Overflow red `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Inline split-view viewport switcher toggling CSS container queries dynamically.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content test-viewports`
  - Flags:
    - `--template <path>`: Relative path to HTML email file.
    - `--min-width <int>`: Minimum mobile breakpoint to validate (default: `320`).
    - `--is-auto-fix-overflow`: Automatically inject `max-width: 100%` on broad tables (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ViewportAuditResult struct {
        AuditID         string    `json:"auditId"`
        TemplatePath    string    `json:"templatePath"`
        TargetDevice    string    `json:"targetDevice"`
        ViewportWidth   int       `json:"viewportWidth"`
        HasOverflow     bool      `json:"hasOverflow"`
        IsFluid         bool      `json:"isFluid"`
        IsMediaQueried  bool      `json:"isMediaQueried"`
        LastEvaluatedAt time.Time `json:"lastEvaluatedAt"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ViewportAuditRecord (
        audit_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        target_device TEXT NOT NULL,
        viewport_width INTEGER NOT NULL,
        has_overflow INTEGER NOT NULL DEFAULT 0,
        is_fluid INTEGER NOT NULL DEFAULT 1,
        is_media_queried INTEGER NOT NULL DEFAULT 1,
        last_evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located at `pkg/email/content/viewport.go`, executed during pre-send content validation.

---

### Improvement 1153: Spam Trigger Heuristic Scanner & Keyword Weighting Analyzer
- **ID:** `CONTENT-SPAM-1153`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan outbound email subject lines and body copy against an updated heuristic database of 850+ high-friction spam keywords, aggressive capitalization patterns, excessive punctuation, and spam-trap linguistic cues, generating a composite SpamAssassin/Barracuda risk score.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual risk gauge card highlighting offending tokens and recommended rephrasings:
      ```text
      ┌─ [HEURISTIC SPAM TRIGGER & PHRASE INTENSITY ANALYZER] ───────────────────────────────┐
      │ COMPOSITE SPAM SCORE: 1.8 / 10.0 [LOW RISK] | PASS THRESHOLD: < 4.0                  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [FLAGGED PHRASES & SUGGESTED REPLACEMENTS]                                           │
      │ ├─ Line 14: "ACT NOW FOR FREE ACCESS" (High Urgency + Free Bait)  [WEIGHT: +1.2]     │
      │ │   └─ Recommendation -> "Start your complimentary evaluation"                       │
      │ ├─ Subject: Multiple exclamation points detected: "Update!!!"     [WEIGHT: +0.6]     │
      │ │   └─ Recommendation -> "Update: Key Milestones"                                    │
      │ └─ Text-to-HTML Ratio: 68% Plaintext (Pristine Balance)           [BONUS:  -0.5]     │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Safe green `#22c55e`, Flagged amber `#f59e0b`, High-risk red `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Inline highlighted text preview with tooltips suggesting neutral rephrasings.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content spam-score`
  - Flags:
    - `--body-file <path>`: Path to email text or HTML file.
    - `--max-score <float>`: Maximum allowable spam risk score (default: `3.5`).
    - `--is-strict-mode`: Fail pipeline if any high-severity spam phrase is detected (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamKeywordMatch struct {
        MatchID         string    `json:"matchId"`
        OffendingPhrase string    `json:"offendingPhrase"`
        WeightPenalty   float64   `json:"weightPenalty"`
        LineNumber      int       `json:"lineNumber"`
        ReplacementHint string    `json:"replacementHint"`
        IsCritical      bool      `json:"isCritical"`
        HasBeenResolved bool      `json:"hasBeenResolved"`
    }

    type SpamScoreReport struct {
        ReportID        string              `json:"reportId"`
        TotalScore      float64             `json:"totalScore"`
        Matches         []SpamKeywordMatch  `json:"matches"`
        EvaluatedAt     time.Time           `json:"evaluatedAt"`
        IsUnderCeiling  bool                `json:"isUnderCeiling"`
        HasUrgencyWords bool                `json:"hasUrgencyWords"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamScoreReportRecord (
        report_id TEXT PRIMARY KEY,
        total_score REAL NOT NULL,
        offending_count INTEGER NOT NULL DEFAULT 0,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_under_ceiling INTEGER NOT NULL DEFAULT 1,
        has_urgency_words INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Core analyzer in `pkg/email/content/spamscore.go`, directly integrated into GitMap lint checks.

---

### Improvement 1154: RFC 2387 Multipart/Related Inline CID Asset Packager
- **ID:** `CONTENT-CID-1154`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically parse remote `<img src="http...">` links in email templates, download and verify image integrity, convert references into Content-ID (`cid:`) attachments inside an RFC 2387 `multipart/related` MIME envelope, eliminating external HTTP requests and remote asset blocking in email clients.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live MIME packaging ledger showing converted inline assets:
      ```text
      ┌─ [RFC 2387 MULTIPART/RELATED CID ASSET PACKAGER] ────────────────────────────────────┐
      │ CONVERTED ASSETS: 3 | PAYLOAD SIZE: +38.4 KB | EXTERNAL LINKS ELIMINATED: 100%       │
      ├────────────────────────┬─────────────┬──────────┬──────────────┬─────────────────────┤
      │ ASSET FILENAME         │ MIME TYPE   │ SIZE     │ CONTENT-ID   │ PACKAGING STATUS    │
      ├────────────────────────┼─────────────┼──────────┼──────────────┼─────────────────────┤
      │ logo-brand.png         │ image/png   │ 14.2 KB  │ <cid-logo-01>│ [●] EMBEDDED INLINE │
      │ icon-check.png         │ image/png   │  4.8 KB  │ <cid-ico-02> │ [●] EMBEDDED INLINE │
      │ signature-avatar.jpg   │ image/jpeg  │ 19.4 KB  │ <cid-ava-03> │ [●] EMBEDDED INLINE │
      └────────────────────────┴─────────────┴──────────┴──────────────┴─────────────────────┘
      ```
    - Lipgloss styles: Asset cyan `#06b6d4`, Embedded green `#22c55e`, Header slate `#64748b`.
  - **Zero-CDN Web Dashboard:**
    - Base64 data-URL / CID inspection table with image dimensions preview.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content package-cid`
  - Flags:
    - `--input-html <path>`: Source HTML email file with remote image URLs.
    - `--max-image-bytes <int>`: Max allowable bytes per image (default: `2097152` = 2MB).
    - `--is-lossless-compress`: Optimize PNG/JPEG bytes before packaging (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CidAssetRecord struct {
        AssetID         string    `json:"assetId"`
        ContentID       string    `json:"contentId"`
        MimeType        string    `json:"mimeType"`
        OriginalURL     string    `json:"originalUrl"`
        SizeBytes       int64     `json:"sizeBytes"`
        EmbeddedAt      time.Time `json:"embeddedAt"`
        IsCompressed    bool      `json:"isCompressed"`
        IsValidImage    bool      `json:"isValidImage"`
        HasMimeHeader   bool      `json:"hasMimeHeader"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CidAssetRecord (
        asset_id TEXT PRIMARY KEY,
        content_id TEXT NOT NULL UNIQUE,
        mime_type TEXT NOT NULL,
        original_url TEXT NOT NULL,
        size_bytes INTEGER NOT NULL,
        embedded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_compressed INTEGER NOT NULL DEFAULT 1,
        is_valid_image INTEGER NOT NULL DEFAULT 1,
        has_mime_header INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/cid.go`, hooked into the outbound MIME encoder.

---

### Improvement 1155: Zero-Dependency Pure CSS Dark Mode Compatibility Analyzer
- **ID:** `CONTENT-DARKMODE-1155`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit email CSS for dark mode rendering resilience across Apple Mail, Gmail, and Outlook, detecting missing `@media (prefers-color-scheme: dark)` overrides, low-contrast text inversions, and missing `[data-ogsc]` tags.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-contrast terminal split view simulating Light Mode vs Dark Mode colors:
      ```text
      ┌─ [ZERO-DEPENDENCY DARK MODE COLOR CONTRAST & WCAG AUDITOR] ──────────────────────────┐
      │ TEMPLATE: corporate-memo.html | DARK MODE RESILIENCE SCORE: 96/100 [AAA PASSED]     │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [LIGHT THEME]                      │ [INVERTED DARK THEME]                           │
      │ Background: #ffffff (White)        │ Background: #18181b (Zinc-900)                  │
      │ Text Color: #0f172a (Slate-900)    │ Text Color: #f8fafc (Slate-50)                  │
      │ Contrast:   16.4:1 [WCAG AAA ✓]    │ Contrast:   14.8:1 [WCAG AAA ✓]                 │
      ├────────────────────────────────────┴─────────────────────────────────────────────────┤
      │ [MICROSOFT OUTLOOK WEB OGSC FALLBACKS]                                               │
      │ ├─ [data-ogsc] .dark-text { color: #f8fafc !important; }              [PRESENT ✓]    │
      │ └─ [data-ogsb] .dark-bg   { background-color: #18181b !important; }    [PRESENT ✓]    │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Light white `#ffffff`, Dark zinc `#18181b`, Pass green `#22c55e`.
  - **Zero-CDN Web Dashboard:**
    - Side-by-side live iframe toggle simulating `prefers-color-scheme: dark`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content test-darkmode`
  - Flags:
    - `--template <path>`: Source HTML email file.
    - `--min-contrast-ratio <float>`: Minimum WCAG contrast threshold (default: `4.5`).
    - `--is-auto-inject-ogsc`: Automatically inject missing Outlook dark mode attributes (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeAuditReport struct {
        ReportID        string    `json:"reportId"`
        TemplatePath    string    `json:"templatePath"`
        LightContrast   float64   `json:"lightContrast"`
        DarkContrast    float64   `json:"darkContrast"`
        EvaluatedAt     time.Time `json:"evaluatedAt"`
        IsWcagCompliant bool      `json:"isWcagCompliant"`
        HasMediaDark    bool      `json:"hasMediaDark"`
        HasOgscFallback bool      `json:"hasOgscFallback"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DarkModeAuditRecord (
        report_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        light_contrast REAL NOT NULL,
        dark_contrast REAL NOT NULL,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_wcag_compliant INTEGER NOT NULL DEFAULT 1,
        has_media_dark INTEGER NOT NULL DEFAULT 1,
        has_ogsc_fallback INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Engine at `pkg/email/content/darkmode.go`, included in the Antigravity email pre-flight suite.

---

### Improvement 1156: Dynamic Plaintext Alternative (text/plain) Synthesizer HUD
- **ID:** `CONTENT-PLAINTEXT-1156`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Autonomously generate clean, well-formatted RFC 2046 `text/plain` alternatives from responsive HTML email sources, preserving ASCII link references, table column alignments, and semantic text hierarchy without third-party HTML parsing packages.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Terminal diff inspector comparing original HTML elements with synthesized Markdown-style ASCII output:
      ```text
      ┌─ [DYNAMIC TEXT/PLAIN ALTERNATIVE SYNTHESIZER] ───────────────────────────────────────┐
      │ SOURCE: newsletter.html (24.2 KB) -> SYNTHESIZED TEXT/PLAIN: (3.1 KB) [CLEAN ASCII]  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [GENERATED PLAINTEXT OUTPUT SAMPLE]                                                  │
      │ ==================================================================================== │
      │ EXECUTIVE ARCHITECTURE BRIEFING: Q3 UPDATE                                           │
      │ ==================================================================================== │
      │                                                                                      │
      │ Hello engineering team,                                                              │
      │                                                                                      │
      │ Here are the primary deliverability metrics for the current warmup run:              │
      │   * Inbox Delivery Rate: 99.2%                                                       │
      │   * Average Transit Lag: 1.14s                                                       │
      │                                                                                      │
      │ Review full audit log: [1] https://ops.domain.com/warmup/audit                       │
      │                                                                                      │
      │ [1] https://ops.domain.com/warmup/audit                                              │
      │ Unsubscribe: https://ops.domain.com/opt-out?id=89412                                 │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Text white `#f8fafc`, Ascii cyan `#06b6d4`, Meta slate `#64748b`.
  - **Zero-CDN Web Dashboard:**
    - Raw monospaced ASCII rendering preview pane with syntax highlighting.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content synthesize-plaintext`
  - Flags:
    - `--html-file <path>`: Source HTML template to parse.
    - `--line-width <int>`: Soft wrap line width for plain text (default: `78`).
    - `--is-numbered-links`: Append numbered link index footnotes at bottom (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PlaintextSynthesisResult struct {
        SynthesisID     string    `json:"synthesisId"`
        SourceHtmlPath  string    `json:"sourceHtmlPath"`
        PlainTextLength int       `json:"plainTextLength"`
        TotalLinksFound int       `json:"totalLinksFound"`
        GeneratedAt     time.Time `json:"generatedAt"`
        IsCompliant     bool      `json:"isCompliant"`
        HasFooterLinks  bool      `json:"hasFooterLinks"`
        HasAsciiBorders bool      `json:"hasAsciiBorders"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PlaintextSynthesisRecord (
        synthesis_id TEXT PRIMARY KEY,
        source_html_path TEXT NOT NULL,
        plain_text_length INTEGER NOT NULL,
        total_links_found INTEGER NOT NULL DEFAULT 0,
        generated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_compliant INTEGER NOT NULL DEFAULT 1,
        has_footer_links INTEGER NOT NULL DEFAULT 1,
        has_ascii_borders INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Packaged at `pkg/email/content/plaintext.go`, generating multipart alternates during compilation.

---

### Improvement 1157: Outbound Email Link Reputation & Redirect Red-Flag Analyzer
- **ID:** `CONTENT-LINKS-1157`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Inspect every hyperlink enclosed within outbound warmup emails, verifying destination HTTPS certificate validity, detecting link shorteners (bit.ly, tinyurl) known to trigger spam filters, and tracing redirect hops to their terminal destinations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Hyperlink safety ledger displaying hop counts and domain reputation status:
      ```text
      ┌─ [EMAIL HYPERLINK REPUTATION & REDIRECT CHAIN INSPECTOR] ────────────────────────────┐
      │ AUDITED LINKS: 5 | ALL HTTPS: YES | SUSPICIOUS SHORTENERS: 0 | REDIRECT HOPS: 1      │
      ├───────────────────────────────┬──────┬──────────┬──────────┬───────────┬─────────────┤
      │ URL TARGET                    │ HOPS │ FINAL DOM│ SSL/TLS  │ SHORTENER │ STATUS      │
      ├───────────────────────────────┼──────┼──────────┼──────────┼───────────┼─────────────┤
      │ https://company.com/login     │  0   │ company  │ TLS 1.3  │ NO        │ [●] SECURE  │
      │ https://go.company.com/meet   │  1   │ zoom.us  │ TLS 1.3  │ NO        │ [●] SECURE  │
      │ https://tinyurl.com/biz-offer │  2   │ spam.biz │ TLS 1.2  │ YES       │ [▲] WARN    │
      │ https://company.com/privacy   │  0   │ company  │ TLS 1.3  │ NO        │ [●] SECURE  │
      └───────────────────────────────┴──────┴──────────┴──────────┴───────────┴─────────────┘
      ```
    - Lipgloss styles: Secure green `#22c55e`, Warn amber `#f59e0b`, Fatal red `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Redirect chain waterfall table detailing intermediate status codes (`301`, `302`, `308`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content audit-links`
  - Flags:
    - `--html-file <path>`: Source HTML file containing links.
    - `--max-redirect-hops <int>`: Maximum allowable redirection hops (default: `2`).
    - `--is-shortener-banned`: Reject templates containing generic link shorteners (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LinkAuditEntry struct {
        LinkID          string    `json:"linkId"`
        SourceURL       string    `json:"sourceUrl"`
        ResolvedURL     string    `json:"resolvedUrl"`
        HopCount        int       `json:"hopCount"`
        TlsVersion      string    `json:"tlsVersion"`
        AuditedAt       time.Time `json:"auditedAt"`
        IsHttps         bool      `json:"isHttps"`
        IsShortener     bool      `json:"isShortener"`
        HasReputationOk bool      `json:"hasReputationOk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS LinkAuditRecord (
        link_id TEXT PRIMARY KEY,
        source_url TEXT NOT NULL,
        resolved_url TEXT NOT NULL,
        hop_count INTEGER NOT NULL DEFAULT 0,
        tls_version TEXT NOT NULL,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_https INTEGER NOT NULL DEFAULT 1,
        is_shortener INTEGER NOT NULL DEFAULT 0,
        has_reputation_ok INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_link_resolved ON LinkAuditRecord(resolved_url);
    ```
- **Antigravity / GitMap Integration Point:**
  - Module in `pkg/email/content/linkaudit.go`, invoked before scheduling outbound sends.

---

### Improvement 1158: RFC 8058 One-Click Unsubscribe Header & Payload Injector
- **ID:** `CONTENT-UNSUB-1158`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically generate and inject compliant RFC 8058 `List-Unsubscribe` and `List-Unsubscribe-Post: List-Unsubscribe=One-Click` headers and visible footer links required by Google Workspace and Yahoo sender guidelines, preventing bulk delivery rejections.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Unsubscribe compliance header inspector card:
      ```text
      ┌─ [RFC 8058 ONE-CLICK UNSUBSCRIBE HEADER & FOOTER AUDITOR] ───────────────────────────┐
      │ COMPLIANCE STATUS: 100% (GOVERNMENT & ISP MANDATE PASSED)                            │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [INJECTED MIME HEADERS]                                                              │
      │ List-Unsubscribe: <https://mail.co/unsub?id=894>, <mailto:unsub@mail.co?subject=unsub>│
      │ List-Unsubscribe-Post: List-Unsubscribe=One-Click                                    │
      │ Feedback-ID: warmup-phase-2:campaign-04:corp:gitmap                                  │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [VISIBLE HTML FOOTER VERIFICATION]                                                   │
      │ ├─ Text Detected: "To cease receiving warmup communications, click here."           │
      │ ├─ Font Size:     12px (Meets accessibility minimums)                                │
      │ └─ Target HREF:   https://mail.co/unsub?id=894                                        │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Valid green `#22c55e`, Header cyan `#38bdf8`, Meta slate `#64748b`.
  - **Zero-CDN Web Dashboard:**
    - Header inspection modal showing rendered email footer preview.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content verify-unsubscribe`
  - Flags:
    - `--html-file <path>`: Source email HTML template.
    - `--tenant-id <string>`: Tenant account identifier to build unique unsubscribe hash.
    - `--is-auto-inject`: Inject missing RFC 8058 headers automatically (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubscribeComplianceRecord struct {
        RecordID        string    `json:"recordId"`
        TenantID        string    `json:"tenantId"`
        UnsubHttpURL    string    `json:"unsubHttpUrl"`
        UnsubMailtoURI  string    `json:"unsubMailtoUri"`
        VerifiedAt      time.Time `json:"verifiedAt"`
        IsOneClickPost  bool      `json:"isOneClickPost"`
        IsVisibleFooter bool      `json:"isVisibleFooter"`
        HasValidHeaders bool      `json:"hasValidHeaders"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnsubscribeComplianceRecord (
        record_id TEXT PRIMARY KEY,
        tenant_id TEXT NOT NULL,
        unsub_http_url TEXT NOT NULL,
        unsub_mailto_uri TEXT NOT NULL,
        verified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_one_click_post INTEGER NOT NULL DEFAULT 1,
        is_visible_footer INTEGER NOT NULL DEFAULT 1,
        has_valid_headers INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_unsub_tenant ON UnsubscribeComplianceRecord(tenant_id);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/unsubscribe.go`, enforced as a mandatory gateway check.

---

### Improvement 1159: Email HTML CSS Inliner & Redundant Rule Purger HUD
- **ID:** `CONTENT-INLINE-1159`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Invert and inline external/internal `<style>` CSS rules directly into HTML `style=""` tag attributes for archaic email client compatibility (Outlook 2016-2021, Gmail Web), while automatically purging unused CSS declarations to minimize payload weight under 102 KB (Gmail clipping limit).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Payload size optimization ledger and Gmail clipping indicator:
      ```text
      ┌─ [EMAIL CSS INLINER & PAYLOAD PRUNING ENGINE] ────────────────────────────────────────┐
      │ SOURCE SIZE: 114.8 KB -> INLINED SIZE: 68.2 KB (-40.6%) | GMAIL CLIP THRESHOLD: 102 KB│
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [INLINING EFFICIENCY METRICS]                                                        │
      │ ├─ Selectors Processed: 184 CSS Rules                                                │
      │ ├─ Inlined Directly:    142 DOM Elements                                             │
      │ ├─ Dead Rules Purged:   42 Unused Declarations (e.g. .btn-unused, .card-hover)       │
      │ └─ GMAIL CLIPPING RISK: SAFE (68.2 KB is 33.8 KB below 102 KB clipping threshold)    │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Safe green `#22c55e`, Pruned cyan `#06b6d4`, Danger red `#ef4444`.
  - **Zero-CDN Web Dashboard:**
    - Split-screen code diff highlighting inlined attributes with byte savings meters.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content inline-css`
  - Flags:
    - `--input <path>`: Source HTML file with `<style>` blocks.
    - `--output <path>`: Output destination for inlined HTML.
    - `--is-purge-unused`: Strip unreferenced CSS classes and IDs (default: `true`).
    - `--has-clip-guard`: Throw error if output exceeds 100 KB (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CssInliningResult struct {
        InliningID      string    `json:"inliningId"`
        OriginalBytes   int64     `json:"originalBytes"`
        InlinedBytes    int64     `json:"inlinedBytes"`
        PurgedRules     int       `json:"purgedRules"`
        ProcessedAt     time.Time `json:"processedAt"`
        IsUnderGmailCap bool      `json:"isUnderGmailCap"`
        IsMediaPreserve bool      `json:"isMediaPreserve"`
        HasInliningDone bool      `json:"hasInliningDone"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CssInliningRecord (
        inlining_id TEXT PRIMARY KEY,
        original_bytes INTEGER NOT NULL,
        inlined_bytes INTEGER NOT NULL,
        purged_rules INTEGER NOT NULL DEFAULT 0,
        processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_under_gmail_cap INTEGER NOT NULL DEFAULT 1,
        is_media_preserve INTEGER NOT NULL DEFAULT 1,
        has_inlining_done INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located at `pkg/email/content/inliner.go`, invoked during template build before delivery.

---

### Improvement 1160: Cryptographic DKIM Canonicalization (relaxed/relaxed) Visualizer HUD
- **ID:** `CONTENT-DKIM-1160`
- **Cluster:** 116 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Visualize and verify RFC 6376 DKIM Canonicalization transformations (`relaxed/relaxed` and `simple/simple`), showing exactly how whitespace normalization, trailing line feeds, and header folding alter the raw message bytes prior to RSA-SHA256 signature hashing.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Byte-level canonicalization inspector rendering whitespace token transforms and SHA256 hashes:
      ```text
      ┌─ [DKIM CANONICALIZATION & RSA-SHA256 HASH VISUALIZER] ───────────────────────────────┐
      │ ALGORITHM: rsa-sha256 | CANONICALIZATION: relaxed/relaxed | SELECTOR: k1-2026        │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [HEADER CANONICALIZATION (relaxed)]                                                  │
      │ RAW:       From:   "System Alert"  <alerts@company.com> \r\n                         │
      │ CANONICAL: from:"System Alert" <alerts@company.com>                                  │
      │ RAW:       Subject:  Weekly    Update   Sync \r\n                                    │
      │ CANONICAL: subject:Weekly Update Sync                                                │
      ├──────────────────────────────────────────────────────────────────────────────────────┤
      │ [BODY HASH COMPUTATION (relaxed)]                                                    │
      │ Normalized Body Size: 2,140 bytes | Trailing CRLF collapsed: YES                     │
      │ Body Hash (bh=):       47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=                  │
      │ DKIM Signature Header: Generated & Verified [VALID RSA-2048 SIGNATURE]               │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
    - Lipgloss styles: Hash purple `#c084fc`, Canonical teal `#14b8a6`, Valid green `#22c55e`.
  - **Zero-CDN Web Dashboard:**
    - Live interactive hex/ASCII canonicalization viewer with byte diff highlights.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email content visualize-dkim`
  - Flags:
    - `--raw-email <path>`: Source MIME email file to canonicalize and sign.
    - `--selector <string>`: DKIM DNS selector to test (default: `default`).
    - `--canonicalization <string>`: Algorithm (`relaxed/relaxed` or `simple/simple`; default: `relaxed/relaxed`).
    - `--is-signature-verified`: Verify computed signature against public key DNS record (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DkimCanonicalRecord struct {
        RecordID         string    `json:"recordId"`
        Selector         string    `json:"selector"`
        Canonicalization string    `json:"canonicalization"`
        BodyHashComputed string    `json:"bodyHashComputed"`
        SignatureHeader  string    `json:"signatureHeader"`
        SignedAt         time.Time `json:"signedAt"`
        IsValidSignature bool      `json:"isValidSignature"`
        HasRelaxedHeader bool      `json:"hasRelaxedHeader"`
        HasRelaxedBody   bool      `json:"hasRelaxedBody"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DkimCanonicalRecord (
        record_id TEXT PRIMARY KEY,
        selector TEXT NOT NULL,
        canonicalization TEXT NOT NULL,
        body_hash_computed TEXT NOT NULL,
        signature_header TEXT NOT NULL,
        signed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_valid_signature INTEGER NOT NULL DEFAULT 1,
        has_relaxed_header INTEGER NOT NULL DEFAULT 1,
        has_relaxed_body INTEGER NOT NULL DEFAULT 1
    );
    CREATE INDEX IF NOT EXISTS idx_dkim_selector ON DkimCanonicalRecord(selector, signed_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/content/dkim.go`, forming the core of the cryptographic mail signer.

