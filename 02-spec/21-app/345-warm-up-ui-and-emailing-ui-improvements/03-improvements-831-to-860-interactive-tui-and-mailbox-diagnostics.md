# Specification 345: Improvements 831 to 860 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 345
- **Slug:** `345-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/03-improvements-831-to-860-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (831 - 860)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 84** | Items 831–840 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 85** | Items 841–850 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 86** | Items 851–860 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 84: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 831–840)

---

### Improvement 831: Real-Time Seed Pool Verification Matrix & Deliverability Heatmap HUD
- **ID:** `SEED-MATRIX-831`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a live terminal matrix view of seed mailboxes partitioned by major consumer and enterprise ESPs (Google Workspace, Microsoft 365, Yahoo/AOL, ProtonMail, iCloud) tracking spam folder placement, inbox placement, and round-trip engagement loop latency.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density matrix table rendering active seed nodes with 4 distinct state indicators: `[INBOX]` (emerald `#10b981`), `[SPAM]` (crimson `#ef4444`), `[DELAY]` (amber `#f59e0b`), `[MISS]` (slate `#64748b`).
    - Interactive cursor with Lipgloss border styling (`#06b6d4`, electric cyan) for keyboard navigation (arrows/Vim keys `hjkl`) showing instantaneous delivery headers, SPF/DKIM/DMARC alignment pass rate, and seed inbox ping latency.
    - Model view snippet:
      ```text
      ┌─ [SEED POOL VERIFICATION MATRIX: 128 ACTIVE SEED INBOXES] ──────────────────────────┐
      │ ESP GROUP         │ TOTAL │ INBOX RATE │ SPAM RATE │ AVG LATENCY │ STATUS            │
      ├───────────────────┼───────┼────────────┼───────────┼─────────────┼───────────────────┤
      │ Google Workspace  │    40 │      97.5% │      2.5% │      1.42s  │ [●] OPTIMAL       │
      │ Microsoft 365     │    35 │      94.2% │      5.8% │      2.10s  │ [●] OPTIMAL       │
      │ Yahoo / AOL       │    25 │      92.0% │      8.0% │      3.45s  │ [▲] ELEVATED LAG  │
      │ ProtonMail        │    15 │     100.0% │      0.0% │      0.88s  │ [●] PRISTINE      │
      │ Apple iCloud      │    13 │      84.6% │     15.4% │      4.12s  │ [■] ATTENTION REQ │
      ├───────────────────┴───────┴────────────┴───────────┴─────────────┴───────────────────┤
      │ [CURSOR: icloud-seed-09@icloud.com]                                                  │
      │ ├─ Last Arrival: SPAM (Folder: Junk) | Subject Hash: 8f4e2a1b                        │
      │ ├─ Authentication: SPF=PASS, DKIM=PASS, DMARC=PASS (Strict Alignment)               │
      │ ├─ Engagement Agent Action: Auto-Moved to Inbox & Flagged as Starred (2s ago)        │
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
    CREATE INDEX IF NOT EXISTS idx_seed_matrix_provider ON SeedPoolVerificationMatrix(provider_type);
    ```
- **Antigravity / GitMap Integration Point:**
  - Registered as sub-command under `cli/cmd/email/warmup/seed_matrix.go`.
  - Telemetry streamed through `pkg/email/warmup/seed_evaluator.go` using GitMap memory audit bus.

---

### Improvement 832: Interactive Seed Auto-Rescue & Quarantine Flow TUI
- **ID:** `SEED-RESCUE-832`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Render a focused Bubble Tea modal workflow to execute automated IMAP rescue operations on emails detected in Spam/Junk folders (marking unread, dragging to inbox, whitelisting sender, simulating human reading scroll time).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Modal popup window overlaid on matrix grid with lipgloss shadow styling.
    - Step-by-step progress checklist displaying active IMAP rescue steps: `[1] IMAP Login`, `[2] Search Junk UID`, `[3] Flag NOT-SPAM`, `[4] Move INBOX`, `[5] Mark Starred`, `[6] Synthetic Reply Loop`.
    - Keybindings: `[y: Execute Rescue]`, `[q: Quarantine Mailbox]`, `[esc: Dismiss Modal]`.
  - **Zero-CDN Web View:**
    - Pure CSS dialog tag (`<dialog>`) styled with native backdrop filter and system font stacks. Zero JavaScript libraries loaded over external networks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-rescue`
  - Flags:
    - `--seed-id <string>`: Unique identifier of seed mailbox requiring rescue triage.
    - `--auto-reply`: Synthesize and send contextual reply back to outbound sender (default: `true`).
    - `--quarantine-threshold <int>`: Consecutively flagged spam count to trigger quarantine (default: `3`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedRescueTask struct {
        TaskID          string    `json:"taskId"`
        SeedID          string    `json:"seedId"`
        OutboundEmailID string    `json:"outboundEmailId"`
        CurrentStep     string    `json:"currentStep"`
        ExecutionStatus string    `json:"executionStatus"`
        RetryCount      int       `json:"retryCount"`
        CreatedAt       time.Time `json:"createdAt"`
        IsCompleted     bool      `json:"isCompleted"`
        IsSuccessful    bool      `json:"isSuccessful"`
        HasReplied      bool      `json:"hasReplied"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedRescueOperations (
        task_id TEXT PRIMARY KEY,
        seed_id TEXT NOT NULL,
        outbound_email_id TEXT NOT NULL,
        current_step TEXT NOT NULL,
        execution_status TEXT NOT NULL,
        retry_count INTEGER NOT NULL DEFAULT 0,
        is_completed INTEGER NOT NULL DEFAULT 0,
        is_successful INTEGER NOT NULL DEFAULT 0,
        has_replied INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Dispatched from `pkg/email/warmup/imap_rescue.go` with audit logging logged into `installation.db` session audit history.

---

### Improvement 833: Mailbox Cluster Entropy & Synthetic Subject Diversification Meter
- **ID:** `ENTROPY-METER-833`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Calculate and visualize lexical Shannon entropy, token variance, and semantic uniqueness of outbound warmup message subject lines and body previews across seed threads to prevent spam-filter template fingerprinting.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live horizontal ASCII/Unicode entropy gauge: `[██████████░░░░] 7.21 / 8.00 bits (High Diversity)`.
    - Real-time n-gram collision warning table highlighting over-utilized phrases (e.g. "quick question regarding", "following up on previous").
    - Lipgloss gradient tokens: Low entropy (`#ef4444`, dangerous cluster), Medium (`#f59e0b`), High (`#10b981`, optimal randomness).
  - **Zero-CDN Web View:**
    - CSS-only radial progress wheel using SVG stroke-dashoffset formulas and localized system font styling.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup entropy-meter`
  - Flags:
    - `--sample-size <int>`: Number of recent generated messages to evaluate (default: `100`).
    - `--min-bits <float>`: Minimum acceptable Shannon entropy threshold in bits (default: `6.5`).
    - `--strict`: Fail CLI exit code if entropy falls below threshold (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EntropyMetrics struct {
        ClusterID       string    `json:"clusterId"`
        EvaluatedCount  int       `json:"evaluatedCount"`
        ShannonEntropy  float64   `json:"shannonEntropy"`
        NgramCollision  float64   `json:"ngramCollision"`
        DominantKeyword string    `json:"dominantKeyword"`
        EvaluatedAt     time.Time `json:"evaluatedAt"`
        IsOptimal       bool      `json:"isOptimal"`
        HasCollisions   bool      `json:"hasCollisions"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxEntropyAudit (
        cluster_id TEXT PRIMARY KEY,
        evaluated_count INTEGER NOT NULL DEFAULT 0,
        shannon_entropy REAL NOT NULL DEFAULT 0.0,
        ngram_collision REAL NOT NULL DEFAULT 0.0,
        dominant_keyword TEXT NOT NULL,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_collisions INTEGER NOT NULL DEFAULT 0,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Ingested by `pkg/email/warmup/entropy_analyzer.go` before queuing batches to `gitmap email dispatch`.

---

### Improvement 834: Seed Inbox Response Engagement Simulator & Typing Cadence HUD
- **ID:** `ENGAGEMENT-HUD-834`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a live terminal HUD monitoring synthetic reply generation, simulating human dwell time, typing velocity (words per minute), and jittered delay curves prior to dispatching responses from seed mailboxes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Real-time status ticker visualizing synthetic human interaction: `[READING: 42s] -> [TYPING: 65 WPM] -> [REPLY SENT]`.
    - Dynamic sparkline indicating simulated reading pause duration relative to incoming message length.
  - **Zero-CDN Web View:**
    - Pure HTML/CSS timeline visualization with CSS keyframe pulse animation.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup engagement-hud`
  - Flags:
    - `--wpm-target <int>`: Average simulated typing velocity in WPM (default: `68`).
    - `--jitter-pct <int>`: Percentage variance applied to timing curves (default: `25`).
    - `--simulate-read`: Include randomized read-time delays prior to composition (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EngagementCadence struct {
        ActionID        string    `json:"actionId"`
        SeedAddress     string    `json:"seedAddress"`
        SimulatedWPM    int       `json:"simulatedWPM"`
        DwellTimeSec    float64   `json:"dwellTimeSec"`
        ResponseDelaySec float64  `json:"responseDelaySec"`
        IsSynthetic     bool      `json:"isSynthetic"`
        HasDispatched   bool      `json:"hasDispatched"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EngagementCadenceLog (
        action_id TEXT PRIMARY KEY,
        seed_address TEXT NOT NULL,
        simulated_wpm INTEGER NOT NULL,
        dwell_time_sec REAL NOT NULL,
        response_delay_sec REAL NOT NULL,
        is_synthetic INTEGER NOT NULL DEFAULT 1,
        has_dispatched INTEGER NOT NULL DEFAULT 0,
        recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/email/warmup/human_emulator.go` and surfaces metrics to interactive dashboard.

---

### Improvement 835: Multi-Tenancy Seed Pool Balance & IP Shuffling Inspector
- **ID:** `SEED-SHUFFLE-835`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Display an interactive distribution map illustrating the dynamic rotation and pairing between outbound sending IP pools and destination seed inboxes to prevent circular routing signatures.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual-column interactive selector connecting Outbound Node IDs to Inbound Seed Node IDs with arrow paths (`[Node A] ──> [Seed G]`).
    - Colorized node tags indicating IP subnet diversity and ASN separation.
  - **Zero-CDN Web View:**
    - Pure CSS SVG node link graph rendered server-side with zero external JavaScript.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-shuffle`
  - Flags:
    - `--max-pairings <int>`: Max allowable interactions between any distinct pair within 7 days (default: `3`).
    - `--subnet-isolation`: Enforce strict /24 subnet separation between sender and seed receiver (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedPairingPolicy struct {
        PairingID       string    `json:"pairingId"`
        SenderIP        string    `json:"senderIP"`
        SeedEmail       string    `json:"seedEmail"`
        InteractionCount int      `json:"interactionCount"`
        LastPairedAt    time.Time `json:"lastPairedAt"`
        IsSubnetIsolated bool     `json:"isSubnetIsolated"`
        HasAvailableSlot bool     `json:"hasAvailableSlot"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedPairingHistory (
        pairing_id TEXT PRIMARY KEY,
        sender_ip TEXT NOT NULL,
        seed_email TEXT NOT NULL,
        interaction_count INTEGER NOT NULL DEFAULT 1,
        is_subnet_isolated INTEGER NOT NULL DEFAULT 1,
        has_available_slot INTEGER NOT NULL DEFAULT 1,
        last_paired_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_seed_pair_sender ON SeedPairingHistory(sender_ip);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/email/warmup/rotator.go` algorithm.

---

### Improvement 836: Real-Time MX Record & SMTP Banner Divergence Probe
- **ID:** `MX-PROBE-836`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Continuously probe destination seed domain MX records, TLS cipher suites, and 220 greeting banners, surfacing unexpected route shifts, security downgrades, or spam gateway insertions (e.g. sudden Barracuda/Mimecast insertion).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live probe table highlighting MX hostname, IP, Banner String, TLS 1.3 confirmation, and divergence flag `[DIVERGENT: Barracuda Detected]`.
    - Instant hotkey `[f: Flush Cache]` to re-resolve authoritative DNS over DNS-over-HTTPS (DoH).
  - **Zero-CDN Web View:**
    - Static HTML report with monospace styled banner diff blocks and color-coded TLS status badges.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mx-probe`
  - Flags:
    - `--domain <string>`: Target seed domain to probe (default: `all`).
    - `--timeout <duration>`: SMTP handshake connection timeout (default: `5s`).
    - `--enforce-dane`: Verify TLSA DNS records for DANE compliance (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MXProbeResult struct {
        DomainName      string    `json:"domainName"`
        MXHost          string    `json:"mxHost"`
        BannerGreeting  string    `json:"bannerGreeting"`
        TLSVersion      string    `json:"tlsVersion"`
        IsDivergent     bool      `json:"isDivergent"`
        HasTLSSupport   bool      `json:"hasTLSSupport"`
        HasDANEValid    bool      `json:"hasDANEValid"`
        CheckedAt       time.Time `json:"checkedAt"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedMXProbeCache (
        domain_name TEXT PRIMARY KEY,
        mx_host TEXT NOT NULL,
        banner_greeting TEXT NOT NULL,
        tls_version TEXT NOT NULL,
        is_divergent INTEGER NOT NULL DEFAULT 0,
        has_tls_support INTEGER NOT NULL DEFAULT 1,
        has_dane_valid INTEGER NOT NULL DEFAULT 0,
        checked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/email/warmup/dns_probe.go` as pre-flight validation hook.

---

### Improvement 837: Seed Account IMAP Quota & Storage Gauge HUD
- **ID:** `IMAP-QUOTA-837`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Track mailbox storage capacity across all active seed inboxes, rendering interactive horizontal gauges to warn operators before storage exhaustion triggers `552 5.2.2 Mailbox full` delivery bounces.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Compact storage meters: `[seed-42@yahoo.com] [████████░░] 78% (1.17 GB / 1.50 GB)`.
    - Automatic purge action suggestion `[p: Run Safe Purge]` to strip attachment storage while retaining headers.
  - **Zero-CDN Web View:**
    - HTML5 native `<progress>` elements styled with zero external CSS frameworks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup imap-quota`
  - Flags:
    - `--warn-threshold <float>`: Storage utilization percentage to flag warning (default: `80.0`).
    - `--auto-purge`: Automatically purge messages older than 30 days if quota exceeds warning threshold (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxQuotaMetric struct {
        SeedID          string  `json:"seedId"`
        EmailAddress    string  `json:"emailAddress"`
        UsedBytes       int64   `json:"usedBytes"`
        MaxBytes        int64   `json:"maxBytes"`
        UsagePercentage float64 `json:"usagePercentage"`
        IsWarning       bool    `json:"isWarning"`
        HasPurgeJob     bool    `json:"hasPurgeJob"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedMailboxQuota (
        seed_id TEXT PRIMARY KEY,
        email_address TEXT NOT NULL,
        used_bytes INTEGER NOT NULL DEFAULT 0,
        max_bytes INTEGER NOT NULL DEFAULT 0,
        usage_percentage REAL NOT NULL DEFAULT 0.0,
        is_warning INTEGER NOT NULL DEFAULT 0,
        has_purge_job INTEGER NOT NULL DEFAULT 0,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Automated audit scheduled via `pkg/email/warmup/quota_monitor.go`.

---

### Improvement 838: Instant Mailbox Health Score Card & Factor Decomposition Modal
- **ID:** `SCORE-DECOMP-838`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Present an interactive multi-factor scorecard decomposing an individual mailbox's composite health rating (0–100%) into weighted sub-scores: Deliverability (40%), Spam Complaint Delta (25%), Dwell Engagement (20%), and Sending Cadence Regularity (15%).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Layered overlay window rendered using Lipgloss rounded borders (`border.Rounded`) and adaptive terminal color backgrounds (`#1e1e2e` base, `#89b4fa` accent).
    - Keyboard-expandable factor rows revealing historical trends and active penalties.
  - **Zero-CDN Web View:**
    - Pure CSS card layout with CSS variables for dynamic color threshold binding (`--score-color: #10b981;`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup score-inspect`
  - Flags:
    - `--mailbox-id <string>`: Identifier of mailbox to inspect.
    - `--detailed`: Show individual penalty events from past 14 days (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HealthScoreDecomposition struct {
        MailboxID       string  `json:"mailboxId"`
        CompositeScore  float64 `json:"compositeScore"`
        DeliverySubscore float64 `json:"deliverySubscore"`
        SpamSubscore    float64 `json:"spamSubscore"`
        EngagementSubscore float64 `json:"engagementSubscore"`
        CadenceSubscore float64 `json:"cadenceSubscore"`
        IsPassing       bool    `json:"isPassing"`
        HasPenalty      bool    `json:"hasPenalty"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxScoreDecomposition (
        mailbox_id TEXT PRIMARY KEY,
        composite_score REAL NOT NULL,
        delivery_subscore REAL NOT NULL,
        spam_subscore REAL NOT NULL,
        engagement_subscore REAL NOT NULL,
        cadence_subscore REAL NOT NULL,
        is_passing INTEGER NOT NULL DEFAULT 1,
        has_penalty INTEGER NOT NULL DEFAULT 0,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Core scoring logic implemented in `pkg/email/warmup/score_calculator.go`.

---

### Improvement 839: Interactive Mailbox Tagging & Dynamic Cohort Slicer
- **ID:** `COHORT-SLICER-839`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Enable live interactive filtering and batch tagging of warmup mailboxes by dynamic cohorts (e.g. `cohort:infant`, `cohort:scaling`, `cohort:veteran`, `cohort:cold-quarantine`) directly within the Bubble Tea matrix interface.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Filter bar with quick-tag toggles: `[/]` to search, `[tab]` to toggle cohort pills (`[All]`, `[Infant: <7d]`, `[Scaling: 7-21d]`, `[Veteran: >21d]`).
    - Bulk selection mode (`spacebar` to toggle, `a` to select all, `t` to apply batch tag).
  - **Zero-CDN Web View:**
    - Monolithic static HTML page with vanilla CSS `:target` and form filter controls without client-side React/Vue payloads.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup cohort-filter`
  - Flags:
    - `--cohort <string>`: Filter mailboxes by assigned cohort name.
    - `--assign-tag <string>`: Apply tag to all mailboxes matching current filter.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxCohortAssignment struct {
        MailboxID    string    `json:"mailboxId"`
        CohortName   string    `json:"cohortName"`
        TagList      string    `json:"tagList"`
        AssignedAt   time.Time `json:"assignedAt"`
        IsInfant     bool      `json:"isInfant"`
        IsScaling    bool      `json:"isScaling"`
        IsVeteran    bool      `json:"isVeteran"`
        HasCustomTag bool      `json:"hasCustomTag"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxCohortTags (
        mailbox_id TEXT PRIMARY KEY,
        cohort_name TEXT NOT NULL,
        tag_list TEXT NOT NULL DEFAULT '',
        is_infant INTEGER NOT NULL DEFAULT 0,
        is_scaling INTEGER NOT NULL DEFAULT 0,
        is_veteran INTEGER NOT NULL DEFAULT 0,
        has_custom_tag INTEGER NOT NULL DEFAULT 0,
        assigned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_mbx_cohort_name ON MailboxCohortTags(cohort_name);
    ```
- **Antigravity / GitMap Integration Point:**
  - Connected to `pkg/email/warmup/cohort_manager.go`.

---

### Improvement 840: Headless Head-of-Line Diagnostic Dumper & Seed Anomaly Ledger
- **ID:** `DIAG-DUMP-840`
- **Cluster:** 84 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a deterministic non-interactive CLI diagnostic dumper that captures head-of-line seed anomalies, authentication failures, and latency outliers for automated CI/CD assertion and audit logs.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Streamlined tabular ASCII report with ANSI colors disabled when non-TTY output is detected.
    - Summary status line: `RESULT: 0 FATAL ANOMALIES, 2 WARNINGS (EXIT 0)`.
  - **Zero-CDN Web View:**
    - Plaintext preformatted `<pre>` view optimized for fast loading and zero script execution.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-diagnostics`
  - Flags:
    - `--json`: Output full anomaly ledger as formatted JSON.
    - `--strict-exit`: Return exit code 1 if any fatal seed anomalies exist (default: `true`).
    - `--output-file <path>`: Write structured report to relative disk path.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedDiagnosticSummary struct {
        Timestamp       time.Time `json:"timestamp"`
        TotalSeeds      int       `json:"totalSeeds"`
        HealthyCount    int       `json:"healthyCount"`
        WarningCount    int       `json:"warningCount"`
        CriticalCount   int       `json:"criticalCount"`
        HasFatalErrors  bool      `json:"hasFatalErrors"`
        IsAuditPassed   bool      `json:"isAuditPassed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedDiagnosticSnapshots (
        snapshot_id TEXT PRIMARY KEY,
        total_seeds INTEGER NOT NULL,
        healthy_count INTEGER NOT NULL,
        warning_count INTEGER NOT NULL,
        critical_count INTEGER NOT NULL,
        has_fatal_errors INTEGER NOT NULL DEFAULT 0,
        is_audit_passed INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmd/email/warmup/diagnostics.go` with exit code governance managed via `pkg/cliexit`.

---

# Cluster 85: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 841–850)

---

### Improvement 841: Multi-Domain Concurrency Governor HUD & Live Token-Bucket Ticker
- **ID:** `CONCURRENCY-GOV-841`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Manage parallel SMTP worker allocations across target destination domains (e.g. `google.com`, `outlook.com`, `yahoo.com`) using token-bucket rate limiters with real-time Bubble Tea bucket visualization to prevent provider-level throttling.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live multi-progress bar rendered via Lipgloss showing token availability per domain:
      `google.com  [██████████░░░░] 24/36 tokens (12 req/min)`
      `outlook.com [██████░░░░░░░░] 12/30 tokens ( 8 req/min)`
    - Instant hotkey `[+]` or `[-]` to dynamically raise or depress global concurrency ceilings on the fly.
  - **Zero-CDN Web View:**
    - Monolithic pure CSS dashboard with SVG token capacity bars, local CSS keyframe animations, and zero external JS frameworks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup concurrency-gov`
  - Flags:
    - `--domain-cap <string>`: Key-value string pairs specifying domain concurrency caps (e.g. `google.com=10,outlook.com=5`).
    - `--default-rate <int>`: Default token refill rate per minute (default: `15`).
    - `--burst-size <int>`: Max allowable token accumulation for bursts (default: `25`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainConcurrencyState struct {
        TargetDomain    string    `json:"targetDomain"`
        CurrentTokens   int       `json:"currentTokens"`
        MaxCapacity     int       `json:"maxCapacity"`
        RefillRateSec   float64   `json:"refillRateSec"`
        ActiveWorkers   int       `json:"activeWorkers"`
        LastRefillAt    time.Time `json:"lastRefillAt"`
        IsThrottled     bool      `json:"isThrottled"`
        HasAvailableCap bool      `json:"hasAvailableCap"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainConcurrencyLimits (
        target_domain TEXT PRIMARY KEY,
        current_tokens INTEGER NOT NULL,
        max_capacity INTEGER NOT NULL,
        refill_rate_sec REAL NOT NULL,
        active_workers INTEGER NOT NULL DEFAULT 0,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        has_available_cap INTEGER NOT NULL DEFAULT 1,
        last_refill_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/email/warmup/concurrency_manager.go` and governs worker pool distribution in `cli/cmd/email/warmup/worker.go`.

---

### Improvement 842: Dynamic Backoff Inversion Engine & Greylist Escalation Breaker
- **ID:** `BACKOFF-INVERT-842`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Invert exponential backoff timers into intelligent retry cadences when SMTP servers respond with 4xx greylisting codes (e.g. 451 4.7.1), shifting re-attempt windows into high-probability acceptance slots (15m, 45m, 2h) instead of naive rapid hammering.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Visual timeline of queued delayed retry slots: `[Attempt 1: +15m] -> [Attempt 2: +45m] -> [Attempt 3: +120m]`.
    - Warning indicator badge `[INVERSION ACTIVE: Greylist Avoidance Pattern]` with countdown timers.
  - **Zero-CDN Web View:**
    - Timeline grid styled using CSS flexbox and inline SVGs.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup backoff-status`
  - Flags:
    - `--auto-invert`: Automatically engage inverted scheduling upon encountering 421 or 451 responses (default: `true`).
    - `--max-delay <duration>`: Ceiling on inverted retry delay (default: `4h`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BackoffInversionPlan struct {
        RetryID         string        `json:"retryId"`
        TargetDomain    string        `json:"targetDomain"`
        OriginalError   string        `json:"originalError"`
        RetryInterval   time.Duration `json:"retryInterval"`
        ScheduledTime   time.Time     `json:"scheduledTime"`
        AttemptCount    int           `json:"attemptCount"`
        IsInverted      bool          `json:"isInverted"`
        HasExecuted     bool          `json:"hasExecuted"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BackoffInversionQueue (
        retry_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        original_error TEXT NOT NULL,
        scheduled_time DATETIME NOT NULL,
        attempt_count INTEGER NOT NULL DEFAULT 0,
        is_inverted INTEGER NOT NULL DEFAULT 1,
        has_executed INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_backoff_scheduled ON BackoffInversionQueue(scheduled_time, has_executed);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/warmup/backoff_governor.go`.

---

### Improvement 843: ISP Response Telemetry Parser & Enhanced Error Code Matrix
- **ID:** `RESP-PARSER-843`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Parse granular SMTP diagnostic strings across diverse mail transfer agents (Postfix, Exim, Sendmail, Exchange) and classify raw error strings into structured failure taxonomies (Rate Limiting, IP Reputation, Content Spam, Authentication Drop).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live feed of intercepted SMTP response status messages with classified tags:
      `[2026-10-27 21:00:12] 550 5.7.1 [Google] -> TAG: REPUTATION_PENALTY (IP Blocklist)`
      `[2026-10-27 21:00:15] 452 4.2.2 [Outlook] -> TAG: MAILBOX_FULL (Remote Inactive)`
    - High-visibility Lipgloss styles separating transient 4xx warnings from fatal 5xx rejections.
  - **Zero-CDN Web View:**
    - Filterable table with native CSS dropdowns for instantaneous category sorting without external script assets.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup parse-errors`
  - Flags:
    - `--raw-error <string>`: Parse a specific raw SMTP response string.
    - `--tail <int>`: Stream the last N categorized SMTP interactions (default: `50`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SMTPResponseClassification struct {
        EventID         string    `json:"eventId"`
        RawResponse     string    `json:"rawResponse"`
        StatusCode      int       `json:"statusCode"`
        EnhancedCode    string    `json:"enhancedCode"`
        CategoryTag     string    `json:"categoryTag"`
        Timestamp       time.Time `json:"timestamp"`
        IsFatal         bool      `json:"isFatal"`
        HasReputationRisk bool    `json:"hasReputationRisk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SMTPResponseAudit (
        event_id TEXT PRIMARY KEY,
        raw_response TEXT NOT NULL,
        status_code INTEGER NOT NULL,
        enhanced_code TEXT NOT NULL,
        category_tag TEXT NOT NULL,
        is_fatal INTEGER NOT NULL DEFAULT 0,
        has_reputation_risk INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_smtp_audit_tag ON SMTPResponseAudit(category_tag);
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded inside `pkg/email/warmup/smtp_classifier.go`.

---

### Improvement 844: Domain-Level Circuit Breaker & Automatic Cooldown State Machine
- **ID:** `CIRCUIT-BREAKER-844`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Implement a three-state circuit breaker pattern (`CLOSED`, `HALF-OPEN`, `OPEN`) for individual destination email domains; trips to `OPEN` when error rates exceed 5% over a 5-minute rolling window, pausing outbound traffic to that domain.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Circuit status pills: `[CLOSED: Normal]` (green), `[HALF-OPEN: Probing 10%]` (yellow), `[OPEN: Cooldown 30m]` (pulsing red).
    - Countdown timer indicating remaining duration until automatic transition to `HALF-OPEN`.
  - **Zero-CDN Web View:**
    - Status dashboard cards with SVG circuit state icons and pure CSS color variables.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup circuit-status`
  - Flags:
    - `--reset <domain>`: Manually force-reset a tripped domain circuit to `CLOSED`.
    - `--trip-threshold <float>`: Failure rate percentage required to trip breaker (default: `5.0`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainCircuitBreaker struct {
        DomainName      string    `json:"domainName"`
        CurrentState    string    `json:"currentState"`
        FailureCount    int       `json:"failureCount"`
        TripThreshold   float64   `json:"tripThreshold"`
        CooldownUntil   time.Time `json:"cooldownUntil"`
        IsOpen          bool      `json:"isOpen"`
        IsHalfOpen      bool      `json:"isHalfOpen"`
        HasTripped      bool      `json:"hasTripped"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainCircuitState (
        domain_name TEXT PRIMARY KEY,
        current_state TEXT NOT NULL DEFAULT 'CLOSED',
        failure_count INTEGER NOT NULL DEFAULT 0,
        trip_threshold REAL NOT NULL DEFAULT 5.0,
        cooldown_until DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_open INTEGER NOT NULL DEFAULT 0,
        is_half_open INTEGER NOT NULL DEFAULT 0,
        has_tripped INTEGER NOT NULL DEFAULT 0,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Maintained by `pkg/email/warmup/circuit_breaker.go`.

---

### Improvement 845: Jittered Micro-Staggering Scheduler & Burst Prevention Engine
- **ID:** `MICRO-STAGGER-845`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Enforce microsecond and millisecond Poisson-distributed jitter between concurrent outbound SMTP connections to eliminate synchronization anomalies that trigger ISP pattern-recognition alarms.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live micro-histogram visualizing dispatch gap intervals (0–1000ms buckets).
    - Visual indicator confirming jitter distribution compliance (`Gaussian Mean: 450ms, StDev: 120ms`).
  - **Zero-CDN Web View:**
    - Server-rendered SVG histogram graph styled with minimal inline CSS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup micro-stagger`
  - Flags:
    - `--mean-delay-ms <int>`: Desired mean dispatch gap in milliseconds (default: `450`).
    - `--stdev-ms <int>`: Standard deviation for jitter calculations (default: `120`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type StaggerConfiguration struct {
        ConfigID        string  `json:"configId"`
        TargetDomain    string  `json:"targetDomain"`
        MeanDelayMs     int     `json:"meanDelayMs"`
        StdevMs         int     `json:"stdevMs"`
        IsPoissonJitter bool    `json:"isPoissonJitter"`
        HasActiveFlow   bool    `json:"hasActiveFlow"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MicroStaggerConfig (
        config_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        mean_delay_ms INTEGER NOT NULL DEFAULT 450,
        stdev_ms INTEGER NOT NULL DEFAULT 120,
        is_poisson_jitter INTEGER NOT NULL DEFAULT 1,
        has_active_flow INTEGER NOT NULL DEFAULT 1,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Core scheduling loop executed in `pkg/email/warmup/stagger_scheduler.go`.

---

### Improvement 846: Multi-Domain Throughput Balancer & Fair-Share Queue Arbitrator
- **ID:** `QUEUE-ARBITER-846`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Arbitrate outbound message queues using weighted fair queuing (WFQ), preventing large bulk batches targeting a single destination domain from starving outbound delivery queues for other destinations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Queued message distribution tree showing pending mail counts segmented by destination MX:
      `[Pending Queues: 1,420 total]`
      `├─ google.com:  240 queued [Weight: 30% | Active: 4]`
      `├─ outlook.com: 180 queued [Weight: 25% | Active: 3]`
      `└─ others:      1,000 queued [Weight: 45% | Active: 6]`
  - **Zero-CDN Web View:**
    - CSS flexbox queue depth bars showing real-time allocation percentages.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup queue-arbitrate`
  - Flags:
    - `--enforce-fairness`: Ensure no single domain occupies >40% of active concurrency slots (default: `true`).
    - `--max-queue-depth <int>`: Alert threshold for pending messages per domain (default: `500`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainQueueWeight struct {
        DomainName      string  `json:"domainName"`
        PendingCount    int     `json:"pendingCount"`
        WeightPercent   float64 `json:"weightPercent"`
        AllocatedSlots  int     `json:"allocatedSlots"`
        IsStarved       bool    `json:"isStarved"`
        HasActiveDrain  bool    `json:"hasActiveDrain"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainQueueArbitration (
        domain_name TEXT PRIMARY KEY,
        pending_count INTEGER NOT NULL DEFAULT 0,
        weight_percent REAL NOT NULL DEFAULT 10.0,
        allocated_slots INTEGER NOT NULL DEFAULT 1,
        is_starved INTEGER NOT NULL DEFAULT 0,
        has_active_drain INTEGER NOT NULL DEFAULT 1,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/email/warmup/queue_arbitrator.go`.

---

### Improvement 847: Automated IP Warmup Ramp-Up Curve Adjuster & Step Governor
- **ID:** `RAMP-GOVERNOR-847`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Automatically calculate and advance daily volume sending thresholds for new IP addresses and domains following smooth logarithmic expansion curves, throttling back down if negative feedback spikes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Ramp progress HUD displaying Day N of 30, Target Quota vs. Actual Sent, and Next Stage Projection.
    - Color indicator: Green (`[RAMP STEADY]`), Yellow (`[HOLD AT CURRENT STEP]`), Red (`[DEGRADE STEP -1]`).
  - **Zero-CDN Web View:**
    - Pure SVG logarithmic curve chart with current sending point plotted as an SVG `<circle>`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ramp-status`
  - Flags:
    - `--account-id <string>`: Target account or domain to inspect.
    - `--override-step <int>`: Force manual adjustment of warmup day index.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type WarmupRampSchedule struct {
        AccountID       string    `json:"accountId"`
        CurrentDayIndex int       `json:"currentDayIndex"`
        TargetQuota     int       `json:"targetQuota"`
        ActualSentToday int       `json:"actualSentToday"`
        CurveModel      string    `json:"curveModel"`
        IsAdvancementAllowed bool `json:"isAdvancementAllowed"`
        HasStepDegraded bool      `json:"hasStepDegraded"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupRampProgress (
        account_id TEXT PRIMARY KEY,
        current_day_index INTEGER NOT NULL DEFAULT 1,
        target_quota INTEGER NOT NULL DEFAULT 10,
        actual_sent_today INTEGER NOT NULL DEFAULT 0,
        curve_model TEXT NOT NULL DEFAULT 'LOGARITHMIC',
        is_advancement_allowed INTEGER NOT NULL DEFAULT 1,
        has_step_degraded INTEGER NOT NULL DEFAULT 0,
        last_evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Evaluated nightly by `pkg/email/warmup/ramp_calculator.go`.

---

### Improvement 848: Real-Time Connection Pool Reuse & TCP Keep-Alive Monitor
- **ID:** `CONN-POOL-848`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Monitor and optimize persistent TLS/SMTP socket reuse across worker threads to avoid repeated three-way TCP handshakes and TLS renegotiation overhead while obeying remote server max-connection directives.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Connection pool breakdown table: `[Domain: google.com] [Active: 6] [Idle: 2] [Reuse Ratio: 88.4%]`.
    - Alert flag for excessive socket thrashing or premature connection drops.
  - **Zero-CDN Web View:**
    - Tabular HTML view with CSS status dot markers.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup conn-pool`
  - Flags:
    - `--max-idle-conns <int>`: Maximum idle socket retention per target host (default: `5`).
    - `--idle-timeout <duration>`: Duration before terminating idle socket (default: `30s`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ConnectionPoolTelemetry struct {
        HostName        string  `json:"hostName"`
        ActiveSockets   int     `json:"activeSockets"`
        IdleSockets     int     `json:"idleSockets"`
        TotalReuses     int64   `json:"totalReuses"`
        ReuseRatio      float64 `json:"reuseRatio"`
        IsOptimal       bool    `json:"isOptimal"`
        HasSocketErrors bool    `json:"hasSocketErrors"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SMTPConnectionPoolStats (
        host_name TEXT PRIMARY KEY,
        active_sockets INTEGER NOT NULL DEFAULT 0,
        idle_sockets INTEGER NOT NULL DEFAULT 0,
        total_reuses INTEGER NOT NULL DEFAULT 0,
        reuse_ratio REAL NOT NULL DEFAULT 0.0,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_socket_errors INTEGER NOT NULL DEFAULT 0,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Handled in `pkg/email/warmup/conn_pool.go`.

---

### Improvement 849: Domain Reputation Cross-Correlation & Early-Warning Canary
- **ID:** `CANARY-CORRELATE-849`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Correlate minor latency increases, sudden greeting banner delays, and TLS handshake latency across canary sending mailboxes to predict imminent IP blocks 15–30 minutes before hard 550 bounces emerge.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Warning banner: `[!] CANARY EARLY WARNING: Latency spike (+240ms) detected on mx.google.com`.
    - Automatic recommended preemptive action: `[Action: Soft Throttle -50% for 20m]`.
  - **Zero-CDN Web View:**
    - Pure CSS alert banner with high-contrast text and zero external telemetry scripts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup canary-check`
  - Flags:
    - `--sensitivity <string>`: Detection sensitivity tier (`low`, `medium`, `high`; default: `medium`).
    - `--auto-throttle`: Automatically throttle sending pool by 50% if canary triggers (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CanaryReputationSignal struct {
        SignalID        string    `json:"signalId"`
        TargetDomain    string    `json:"targetDomain"`
        LatencyDeltaMs  int       `json:"latencyDeltaMs"`
        GreetingLagMs   int       `json:"greetingLagMs"`
        RiskScore       float64   `json:"riskScore"`
        TriggeredAt     time.Time `json:"triggeredAt"`
        IsCanaryAlert   bool      `json:"isCanaryAlert"`
        HasThrottled    bool      `json:"hasThrottled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CanaryTelemetrySignals (
        signal_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        latency_delta_ms INTEGER NOT NULL,
        greeting_lag_ms INTEGER NOT NULL,
        risk_score REAL NOT NULL,
        is_canary_alert INTEGER NOT NULL DEFAULT 0,
        has_throttled INTEGER NOT NULL DEFAULT 0,
        triggered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/email/warmup/canary_watcher.go`.

---

### Improvement 850: Concurrency Governor Benchmark & Stress Test Harness
- **ID:** `GOV-BENCHMARK-850`
- **Cluster:** 85 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Run an isolated local stress benchmark verifying that multi-domain token bucket governors, fair queue arbitrators, and backoff inversion modules accurately cap and schedule 10,000 synthetic requests without CPU pegging or lock contention.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live benchmark progress HUD showing operations/sec, p99 scheduling latency, lock contention microseconds, and memory allocations.
    - Final report card: `P99 LATENCY: 0.12ms | LOCK CONTENTION: 0.00% | PASS`.
  - **Zero-CDN Web View:**
    - Static performance benchmark table rendered with pure HTML/CSS.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup benchmark-governor`
  - Flags:
    - `--requests <int>`: Number of simulated concurrent dispatch evaluations (default: `10000`).
    - `--domains <int>`: Number of distinct destination domains to simulate (default: `50`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BenchmarkResult struct {
        BenchmarkID     string    `json:"benchmarkId"`
        TotalRequests   int       `json:"totalRequests"`
        DurationMs      int64     `json:"durationMs"`
        P99LatencyMicro int64     `json:"p99LatencyMicro"`
        ThroughputRPS   float64   `json:"throughputRPS"`
        IsSuccess       bool      `json:"isSuccess"`
        HasContention   bool      `json:"hasContention"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GovernorBenchmarkRuns (
        benchmark_id TEXT PRIMARY KEY,
        total_requests INTEGER NOT NULL,
        duration_ms INTEGER NOT NULL,
        p99_latency_micro INTEGER NOT NULL,
        throughput_rps REAL NOT NULL,
        is_success INTEGER NOT NULL DEFAULT 1,
        has_contention INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `cli/cmd/email/warmup/benchmark.go` and verified via standard test suites.

---

# Cluster 86: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 851–860)

---

### Improvement 851: Zero-CDN Local HTML Email Sandbox Renderer & CSS Inliner Preview
- **ID:** `SANDBOX-RENDER-851`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Render responsive HTML email templates locally inside an isolated, scriptless sandbox previewing layout behavior across simulated mobile (320px/375px) and desktop (600px/800px) viewports with zero external CDN dependencies.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-fidelity terminal HTML renderer converting tags to formatted Lipgloss text, preserving table column layouts, link anchor annotations, and inline image placeholders.
    - Viewport toggle hotkey: `[m: Mobile 375px]` / `[d: Desktop 600px]`.
  - **Zero-CDN Web View:**
    - Dual iframe sandboxed view (`sandbox="allow-same-origin"`, script execution blocked) with CSS media query toggles. Local static font declarations.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup preview-sandbox`
  - Flags:
    - `--template <path>`: Path to relative HTML email template file.
    - `--inline-css`: Execute CSS inlining engine before previewing (default: `true`).
    - `--port <int>`: Local port for zero-CDN browser preview (default: `8911`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailSandboxView struct {
        TemplateID      string    `json:"templateId"`
        ViewportMode    string    `json:"viewportMode"`
        RawHTMLSize     int       `json:"rawHTMLSize"`
        InlinedHTMLSize int       `json:"inlinedHTMLSize"`
        IsResponsive    bool      `json:"isResponsive"`
        HasExternalRefs bool      `json:"hasExternalRefs"`
        RenderedAt      time.Time `json:"renderedAt"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailTemplateSandbox (
        template_id TEXT PRIMARY KEY,
        viewport_mode TEXT NOT NULL DEFAULT 'DESKTOP',
        raw_html_size INTEGER NOT NULL,
        inlined_html_size INTEGER NOT NULL,
        is_responsive INTEGER NOT NULL DEFAULT 1,
        has_external_refs INTEGER NOT NULL DEFAULT 0,
        rendered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Core parsing performed in `pkg/email/warmup/template_renderer.go`.

---

### Improvement 852: Content Security Policy (CSP) & Zero-Leakage Static Resource Auditor
- **ID:** `CSP-AUDIT-852`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Audit HTML warmup message templates for privacy leaks, third-party tracking beacons, unverified external font calls, and insecure `http://` image references that trigger recipient mail client security warnings.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Diagnostic tree listing detected asset origins:
      `[TEMPLATE ASSET AUDIT: PASS]`
      `├─ External Scripts: 0 (Strictly Banned)`
      `├─ Insecure URIs:     0 (100% HTTPS/CID)`
      `└─ Tracking Pixels:  0 Detected`
    - High-contrast alert badge if remote unverified domains are referenced.
  - **Zero-CDN Web View:**
    - Audit scorecard with pure CSS warning callouts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup audit-csp`
  - Flags:
    - `--strict-local`: Ban all remote asset URLs, requiring inline CID or Base64 images (default: `true`).
    - `--allow-hosts <string>`: Comma-separated list of whitelisted external media hosts.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CSPAuditReport struct {
        ReportID        string    `json:"reportId"`
        TemplateID      string    `json:"templateId"`
        ExternalCallCount int     `json:"externalCallCount"`
        InsecureLinkCount int     `json:"insecureLinkCount"`
        IsCompliant     bool      `json:"isCompliant"`
        HasTrackingLeak bool      `json:"hasTrackingLeak"`
        AuditedAt       time.Time `json:"auditedAt"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailCSPAuditHistory (
        report_id TEXT PRIMARY KEY,
        template_id TEXT NOT NULL,
        external_call_count INTEGER NOT NULL DEFAULT 0,
        insecure_link_count INTEGER NOT NULL DEFAULT 0,
        is_compliant INTEGER NOT NULL DEFAULT 1,
        has_tracking_leak INTEGER NOT NULL DEFAULT 0,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Linting rule managed in `pkg/email/warmup/csp_linter.go`.

---

### Improvement 853: HTML AST Sanitizer & Dangerous Tag Stripper HUD
- **ID:** `HTML-SANITIZE-853`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Parse outbound email HTML into an Abstract Syntax Tree (AST), mechanically removing prohibited tags (`<script>`, `<applet>`, `<embed>`, `<object>`, `<iframe>`, `on*` event handlers) while displaying a side-by-side terminal diff of stripped elements.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Side-by-side terminal diff showing Raw Source vs Sanitized Source with removed constructs highlighted in red strikethrough.
    - Status line: `STRIPPED 3 ILLEGAL TOKENS: [onclick, <script>, onerror]`.
  - **Zero-CDN Web View:**
    - Pure HTML/CSS diff table with line numbers and green/red block styling.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup sanitize-html`
  - Flags:
    - `--file <path>`: HTML file to sanitize.
    - `--inplace`: Overwrite file directly with sanitized payload (default: `false`).
    - `--strip-comments`: Remove all HTML comments to save payload bytes (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SanitizerResult struct {
        SanitizeID      string    `json:"sanitizeId"`
        SourceLength    int       `json:"sourceLength"`
        SanitizedLength int       `json:"sanitizedLength"`
        StrippedTagCount int      `json:"strippedTagCount"`
        IsSanitary      bool      `json:"isSanitary"`
        HasMutations    bool      `json:"hasMutations"`
        ProcessedAt     time.Time `json:"processedAt"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HTMLSanitizerAudit (
        sanitize_id TEXT PRIMARY KEY,
        source_length INTEGER NOT NULL,
        sanitized_length INTEGER NOT NULL,
        stripped_tag_count INTEGER NOT NULL,
        is_sanitary INTEGER NOT NULL DEFAULT 1,
        has_mutations INTEGER NOT NULL DEFAULT 0,
        processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/warmup/html_sanitizer.go`.

---

### Improvement 854: Spam Trigger Keyword Scanner & Lexical Density Inspector
- **ID:** `SPAM-SCANNER-854`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Analyze text body and subject line copy against an embedded offline dictionary of high-risk spam trigger tokens (financial gimmicks, urgency hooks, spam phrases), highlighting offending tokens in terminal view.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Text preview panel with flagged tokens highlighted in bold red background:
      `Subject: [ACT NOW!] Exclusive [100% FREE] opportunity for your team`
    - Score deduction summary showing individual token weights and recommendations.
  - **Zero-CDN Web View:**
    - Clean text view with SVG warning pills and inline suggestions.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup scan-keywords`
  - Flags:
    - `--body <string>`: Email body string to inspect.
    - `--max-spam-score <float>`: Max allowable spam token weight (default: `2.0`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamKeywordScan struct {
        ScanID          string    `json:"scanId"`
        TotalWords      int       `json:"totalWords"`
        SpamTokenCount  int       `json:"spamTokenCount"`
        CalculatedRisk  float64   `json:"calculatedRisk"`
        FlaggedTokens   string    `json:"flaggedTokens"`
        IsAcceptable    bool      `json:"isAcceptable"`
        HasUrgencyHooks bool      `json:"hasUrgencyHooks"`
        ScannedAt       time.Time `json:"scannedAt"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamKeywordAudits (
        scan_id TEXT PRIMARY KEY,
        total_words INTEGER NOT NULL,
        spam_token_count INTEGER NOT NULL,
        calculated_risk REAL NOT NULL,
        flagged_tokens TEXT NOT NULL DEFAULT '',
        is_acceptable INTEGER NOT NULL DEFAULT 1,
        has_urgency_hooks INTEGER NOT NULL DEFAULT 0,
        scanned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Evaluated in `pkg/email/warmup/keyword_scanner.go`.

---

### Improvement 855: Responsive Table & Mobile Viewport Geometry Validator
- **ID:** `TABLE-GEOM-855`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Validate HTML email layout geometry, ensuring all nested `<table>` elements enforce responsive widths (`max-width: 600px`, `width: 100%`) and prevent rigid horizontal blowouts on mobile screens (<375px).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Table tree layout visualizing DOM hierarchy with table dimensions:
      `<table> [max-width: 600px | width: 100%] -> OK`
      `  └── <table> [width: 850px (FIXED)] -> [!] BLOWOUT RISK ON MOBILE`
  - **Zero-CDN Web View:**
    - Interactive ruler overlay in CSS preview displaying pixel widths across device breakpoints.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup check-tables`
  - Flags:
    - `--file <path>`: Template file to inspect.
    - `--fix-widths`: Automatically convert fixed pixel table widths to responsive percentage rules (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TableGeometryReport struct {
        ReportID        string  `json:"reportId"`
        TotalTables     int     `json:"totalTables"`
        FixedTableCount int     `json:"fixedTableCount"`
        MaxCanvasWidth  int     `json:"maxCanvasWidth"`
        IsResponsive    bool    `json:"isResponsive"`
        HasFixedBlowout bool    `json:"hasFixedBlowout"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailTableGeometry (
        report_id TEXT PRIMARY KEY,
        total_tables INTEGER NOT NULL,
        fixed_table_count INTEGER NOT NULL,
        max_canvas_width INTEGER NOT NULL,
        is_responsive INTEGER NOT NULL DEFAULT 1,
        has_fixed_blowout INTEGER NOT NULL DEFAULT 0,
        analyzed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `pkg/email/warmup/geometry_validator.go`.

---

### Improvement 856: Dark-Mode Contrast & Inverted Color Palette Simulator
- **ID:** `DARK-MODE-856`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Simulate aggressive dark mode color inversions implemented by Outlook iOS, Apple Mail, and Gmail, identifying illegible low-contrast text (e.g. dark grey text on auto-inverted black backgrounds).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual rendering showing Light Mode vs Inverted Dark Mode representation in ANSI colors.
    - Contrast ratio warning flag: `[!] Text contrast falls to 2.1:1 on dark mode (WCAG AA Fail)`.
  - **Zero-CDN Web View:**
    - CSS `@media (prefers-color-scheme: dark)` toggle container with local slider control.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup dark-mode-check`
  - Flags:
    - `--min-contrast <float>`: Minimum acceptable WCAG contrast ratio (default: `4.5`).
    - `--inject-dark-css`: Automatically inject `@media (prefers-color-scheme: dark)` fallback overrides (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeAudit struct {
        AuditID         string  `json:"auditId"`
        LightContrast   float64 `json:"lightContrast"`
        DarkContrast    float64 `json:"darkContrast"`
        IsWCAGCompliant bool    `json:"isWCAGCompliant"`
        HasInversionBug bool    `json:"hasInversionBug"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailDarkModeAudit (
        audit_id TEXT PRIMARY KEY,
        light_contrast REAL NOT NULL,
        dark_contrast REAL NOT NULL,
        is_wcag_compliant INTEGER NOT NULL DEFAULT 1,
        has_inversion_bug INTEGER NOT NULL DEFAULT 0,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Handled by `pkg/email/warmup/contrast_evaluator.go`.

---

### Improvement 857: RFC 5322 MIME Structure & Header Conformance Linter
- **ID:** `MIME-LINT-857`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify that the generated raw MIME structure adheres strictly to RFC 5322 and RFC 2045/2046 standards, ensuring valid boundary tokens, matching `multipart/alternative` parts, UTF-8 header encoding, and correct `Date` formatting.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Tree representation of MIME body parts:
      `multipart/alternative; boundary="==_NextPart_001"`
      `├─ text/plain; charset=UTF-8 [Encoding: 7bit | Size: 1.2 KB]`
      `└─ text/html; charset=UTF-8 [Encoding: quoted-printable | Size: 4.8 KB]`
    - Linter validation checks: `[PASS: 8/8 RFC Checks]`.
  - **Zero-CDN Web View:**
    - Static HTML report presenting MIME tree nodes and raw header dumps.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup lint-mime`
  - Flags:
    - `--raw-eml <path>`: Path to raw `.eml` or RFC 822 file.
    - `--strict`: Fail on minor whitespace or header line-length warnings (>78 chars).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MIMELintSummary struct {
        FileID          string  `json:"fileId"`
        PartCount       int     `json:"partCount"`
        MaxLineLength   int     `json:"maxLineLength"`
        IsRFC5322Valid  bool    `json:"isRFC5322Valid"`
        HasPlainFallback bool   `json:"hasPlainFallback"`
        HasValidHeaders bool    `json:"hasValidHeaders"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailMIMELintHistory (
        file_id TEXT PRIMARY KEY,
        part_count INTEGER NOT NULL,
        max_line_length INTEGER NOT NULL,
        is_rfc5322_valid INTEGER NOT NULL DEFAULT 1,
        has_plain_fallback INTEGER NOT NULL DEFAULT 1,
        has_valid_headers INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Dispatched via `pkg/email/warmup/mime_linter.go`.

---

### Improvement 858: Inline Image CID Packager & Base64 Payload Sizing Guard
- **ID:** `CID-PACKAGER-858`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Package local brand images into MIME Content-ID (CID) inline attachments, enforcing strict payload size limits (<100 KB total message weight) to prevent Gmail clipping (`[Message clipped] View entire message`).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Asset packager progress bar: `Packing logo.png [CID: logo-img-01] -> 8.4 KB (Total: 42.1 KB / 102 KB Limit)`.
    - Warning alert if payload exceeds 102 KB Gmail clipping threshold.
  - **Zero-CDN Web View:**
    - Monolithic static table summarizing embedded asset sizes and CID mappings.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup package-cid`
  - Flags:
    - `--html-file <path>`: HTML template with relative local image sources.
    - `--max-total-kb <int>`: Max allowable message size before triggering error (default: `98`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CIDPackageManifest struct {
        PackageID       string  `json:"packageId"`
        AssetCount      int     `json:"assetCount"`
        TotalBytes      int64   `json:"totalBytes"`
        IsUnderGmailCap bool    `json:"isUnderGmailCap"`
        HasInlineAssets bool    `json:"hasInlineAssets"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CIDPackageAudit (
        package_id TEXT PRIMARY KEY,
        asset_count INTEGER NOT NULL,
        total_bytes INTEGER NOT NULL,
        is_under_gmail_cap INTEGER NOT NULL DEFAULT 1,
        has_inline_assets INTEGER NOT NULL DEFAULT 1,
        packaged_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/warmup/cid_packager.go`.

---

### Improvement 859: Unsubscribe Link & Header Compliance Verification Modal
- **ID:** `UNSUB-VERIFY-859`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Verify that outbound email templates and generated MIME headers satisfy 2024–2026 Google & Yahoo sender mandates for one-click unsubscribe headers (`List-Unsubscribe: <mailto:...>, <https:...>` and `List-Unsubscribe-Post: List-Unsubscribe=One-Click`).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Compliance checklist badge:
      `[MANDATE CHECK: RFC 8058 ONE-CLICK UNSUBSCRIBE]`
      `├─ List-Unsubscribe Header:      [PRESENT & VALID]`
      `├─ List-Unsubscribe-Post Header: [PRESENT: List-Unsubscribe=One-Click]`
      `└─ Body Unsubscribe Link:        [MATCHED FOOTER ANCHOR]`
    - Status: `100% SENDER MANDATE COMPLIANT`.
  - **Zero-CDN Web View:**
    - Pure HTML5 badge layout with CSS status checks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup check-unsubscribe`
  - Flags:
    - `--template <path>`: Template file to inspect.
    - `--enforce-mandate`: Fail build if RFC 8058 headers are absent (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubscribeCompliance struct {
        TemplateID      string  `json:"templateId"`
        HasMailtoTarget bool    `json:"hasMailtoTarget"`
        HasHTTPSURL     bool    `json:"hasHTTPSURL"`
        HasOneClickPost bool    `json:"hasOneClickPost"`
        IsMandatePassed bool    `json:"isMandatePassed"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnsubscribeComplianceLog (
        template_id TEXT PRIMARY KEY,
        has_mailto_target INTEGER NOT NULL DEFAULT 1,
        has_https_url INTEGER NOT NULL DEFAULT 1,
        has_one_click_post INTEGER NOT NULL DEFAULT 1,
        is_mandate_passed INTEGER NOT NULL DEFAULT 1,
        verified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Validated by `pkg/email/warmup/unsub_validator.go`.

---

### Improvement 860: Comprehensive Pre-Flight Email Sandbox Audit Suite & Certification Dumper
- **ID:** `SANDBOX-SUITE-860`
- **Cluster:** 86 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Aggregate all Cluster 86 sandbox checks (CSP audit, AST sanitization, spam keyword density, table geometry, dark mode contrast, MIME RFC conformance, CID sizing, and unsubscribe verification) into a single deterministic pre-flight certification suite that issues a cryptographic verification token.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Full-screen summary dashboard aggregating 8 diagnostic categories into a tabular score card:
      `┌─ [PRE-FLIGHT EMAIL CERTIFICATION SUITE: TEMPLATE-WELCOME-V2] ────────────┐`
      `│ CATEGORY                 │ STATUS     │ SEVERITY   │ DETAILS              │`
      `├──────────────────────────┼────────────┼────────────┼──────────────────────┤`
      `│ 1. CSP / Resource Leak   │ [PASS]     │ None       │ 0 External Calls     │`
      `│ 2. AST HTML Sanitizer    │ [PASS]     │ None       │ 0 Disallowed Tags    │`
      `│ 3. Spam Keyword Density  │ [PASS]     │ Low        │ Score 0.4 / 2.0      │`
      `│ 4. Table Geometry        │ [PASS]     │ None       │ Max 600px Responsive │`
      `│ 5. Dark Mode Contrast    │ [PASS]     │ None       │ Min 5.2:1 Ratio      │`
      `│ 6. MIME RFC 5322 Lint    │ [PASS]     │ None       │ Valid Part Tree      │`
      `│ 7. CID / Gmail Sizing    │ [PASS]     │ None       │ Total 46.2 KB (<98K) │`
      `│ 8. One-Click Unsubscribe │ [PASS]     │ None       │ RFC 8058 Compliant   │`
      `├──────────────────────────┴────────────┴────────────┴──────────────────────┤`
      `│ OVERALL RESULT: CERTIFIED READY FOR WARMUP (TOKEN: a7f89c20e14d)         │`
      `└──────────────────────────────────────────────────────────────────────────┘`
  - **Zero-CDN Web View:**
    - High-authority HTML certificate page with inline CSS and printable stylesheet rules (`@media print`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup certify-template`
  - Flags:
    - `--template <path>`: Template file or directory to certify.
    - `--output-cert <path>`: Path to write signed JSON certificate token.
    - `--strict`: Terminate with non-zero exit code on any single category warning (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PreFlightCertificate struct {
        CertificateID   string    `json:"certificateId"`
        TemplateID      string    `json:"templateId"`
        VerificationToken string  `json:"verificationToken"`
        TotalChecks     int       `json:"totalChecks"`
        PassedChecks    int       `json:"passedChecks"`
        IsCertified     bool      `json:"isCertified"`
        HasWarnings     bool      `json:"hasWarnings"`
        CertifiedAt     time.Time `json:"certifiedAt"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PreFlightCertificates (
        certificate_id TEXT PRIMARY KEY,
        template_id TEXT NOT NULL,
        verification_token TEXT NOT NULL,
        total_checks INTEGER NOT NULL,
        passed_checks INTEGER NOT NULL,
        is_certified INTEGER NOT NULL DEFAULT 1,
        has_warnings INTEGER NOT NULL DEFAULT 0,
        certified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_cert_token ON PreFlightCertificates(verification_token);
    ```
- **Antigravity / GitMap Integration Point:**
  - Orchestrated in `cli/cmd/email/warmup/certify.go` and logged to `gitmap.db` audit ledger.

---
