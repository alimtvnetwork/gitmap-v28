# Specification 344: Improvements 731 to 760 (Interactive Bubble Tea Mailbox Health Matrix, Dynamic Greylisting Auto-Detection & Zero-CDN Content Security Sandbox)

- **Sequence:** 344
- **Slug:** `344-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/03-improvements-731-to-760-interactive-tui-and-mailbox-telemetry.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (731 - 760)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 74** | Items 731–740 | Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD |
| **Cluster 75** | Items 741–750 | Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler |
| **Cluster 76** | Items 751–760 | Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor |

---

# Cluster 74: Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD (Items 731–740)

---

### Improvement 731: High-Density 2D Mailbox Health Matrix Grid with Bidirectional Cursor Navigation
- **ID:** `MAILBOX-HEALTH-731`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Render a high-density, multi-mailbox 2D health matrix in Bubble Tea allowing operators to visually inspect aggregate warmup velocity, deliverability percentage, and latency across hundreds of active mailboxes with arrow-key cursor navigation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Grid layout displaying 10 mailboxes per row with 2-character glyph indicators (`[●]`, `[▲]`, `[■]`, `[✖]`) color-coded by health score tier.
    - Cursor selection box with high-contrast inverted border (`#38bdf8` bright sky blue) showing real-time telemetry drilldown in a docked footer panel.
    - Lipgloss token palette: Pristine/Warm (`#10b981`, emerald), Active Scaling (`#3b82f6`, blue), Latency Warning (`#f59e0b`, amber), Blocked/Greylisted (`#ef4444`, red), Cursor Frame (`#e0e7ff`, white-indigo).
    - Model view snippet:
      ```text
      ┌─ [MAILBOX HEALTH MATRIX GRID: 64 ACTIVE INBOXES] ──────────────────────────┐
      │ [●] mbx-01  [●] mbx-02  [▲] mbx-03  [●] mbx-04  [■] mbx-05  [●] mbx-06     │
      │ [●] mbx-07  [✖] mbx-08  [●] mbx-09  [●] mbx-10  [●] mbx-11  [▲] mbx-12     │
      │ [●] mbx-13  [●] mbx-14  [●] mbx-15  [●] mbx-16  [●] mbx-17  [●] mbx-18     │
      │ [CURSOR: mbx-08@outbound.net]                                              │
      │ ├─ Health Score: 41.2% | State: GREYLISTED_SPIKE | Outbound: 142/250        │
      │ ├─ Target ISP: Google Workspace (451 4.7.1 Rate limit exceeded)            │
      │ └─ Action: [t: Triage] [p: Pause] [r: Force Re-check] [q: Quit]            │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - CSS Grid matrix layout (`display: grid; grid-template-columns: repeat(auto-fill, minmax(80px, 1fr)); gap: 6px;`).
    - Pure inline SVG status indicators with tooltips styled with pure CSS (`:hover` pseudo-selectors, no external JS runtime).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup matrix`
  - Flags:
    - `--domain <string>`: Filter matrix accounts by sender domain (default: `all`).
    - `--min-health <float>`: Filter mailboxes with health score below threshold (default: `100.0`).
    - `--refresh-interval <duration>`: TUI polling refresh frequency (default: `2s`).
    - `--interactive`: Launch full-screen interactive Bubble Tea grid navigator (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxGridCell struct {
        MailboxID       string  `json:"mailboxId"`
        EmailAddress    string  `json:"emailAddress"`
        DomainName      string  `json:"domainName"`
        HealthScore     float64 `json:"healthScore"`
        SentCount       int     `json:"sentCount"`
        TargetQuota     int     `json:"targetQuota"`
        GridRow         int     `json:"gridRow"`
        GridCol         int     `json:"gridCol"`
        IsSelected      bool    `json:"isSelected"`
        IsHealthy       bool    `json:"isHealthy"`
        HasActiveAlert  bool    `json:"hasActiveAlert"`
    }

    type MailboxMatrixModel struct {
        TotalMailboxes  int               `json:"totalMailboxes"`
        ActiveCursorRow int               `json:"activeCursorRow"`
        ActiveCursorCol int               `json:"activeCursorCol"`
        Cells           []MailboxGridCell `json:"cells"`
        IsLoaded        bool              `json:"isLoaded"`
        HasFocus        bool              `json:"hasFocus"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxHealthMatrixGrid (
        mailbox_id TEXT PRIMARY KEY,
        domain_name TEXT NOT NULL,
        email_address TEXT NOT NULL,
        health_score REAL NOT NULL DEFAULT 100.0,
        sent_count INTEGER NOT NULL DEFAULT 0,
        target_quota INTEGER NOT NULL DEFAULT 50,
        is_healthy INTEGER NOT NULL DEFAULT 1,
        has_active_alert INTEGER NOT NULL DEFAULT 0,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_mbx_health_grid_domain ON MailboxHealthMatrixGrid(domain_name);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/ui/matrix.go` and registers CLI command under `cmd/email/warmup/matrix.go` hooked directly to GitMap's terminal tea runner.

---

### Improvement 732: Instant Seed-Pool Triage HUD with Real-Time Drift Detection
- **ID:** `SEED-TRIAGE-732`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Provide a dedicated instant triage HUD for operator intervention when seed accounts report delivery drift, unread spikes, or inbox placement drops below SLA thresholds.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Split-pane layout: Left pane displays list of drifted seed accounts with delta indicators; Right pane shows diagnostic failure log and remediation actions.
    - Lipgloss border styling with pulsing status badge (`[TRIAGE REQUIRED]` `#ef4444`).
    - Model view snippet:
      ```text
      ┌─ [SEED-POOL DRIFT TRIAGE HUD] ─────────────────────────────────────────────┐
      │ DRIFTED SEEDS (3 DETECTED)     │ DIAGNOSTIC LOG: seed-04@yahoo-test.com      │
      │ ▶ seed-04@yahoo.com [-14.2%]   │ ├─ Last Placed: SPAM (Foldering Drift)     │
      │   seed-12@outlook.com [-8.1%]  │ ├─ Seed Engagement: Missed Star Action     │
      │   seed-19@gmail.com   [-6.5%]  │ ├─ Probable Cause: Greylist Retry Expired  │
      │                                │ └─ Action: [a: Auto-Recover] [r: Rotate]   │
      │ [↑/↓: Navigate] [Space: Inspect] [a: Remediate All] [Esc: Back]             │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline flexbox split layout styled via `<style>` tag without framework dependencies.
    - HTML5 native `<details>` and `<summary>` tags for expanding triage diagnostic logs.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup triage`
  - Flags:
    - `--threshold-drift <float>`: Drift percentage threshold triggering triage (default: `5.0`).
    - `--auto-remediate`: Automatically apply recovery scripts to drifted seed accounts (default: `false`).
    - `--pool-id <string>`: Filter triage to specific seed pool identifier (default: `all`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedTriageRecord struct {
        SeedAddress      string    `json:"seedAddress"`
        PoolID           string    `json:"poolId"`
        ExpectedRate     float64   `json:"expectedRate"`
        ActualRate       float64   `json:"actualRate"`
        DriftPercentage  float64   `json:"driftPercentage"`
        LastPlacement    string    `json:"lastPlacement"`
        RemediationStep  string    `json:"remediationStep"`
        IsUnderTriage    bool      `json:"isUnderTriage"`
        HasAutoRecovered bool      `json:"hasAutoRecovered"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedTriageIncidents (
        incident_id TEXT PRIMARY KEY,
        seed_address TEXT NOT NULL,
        pool_id TEXT NOT NULL,
        drift_percentage REAL NOT NULL,
        last_placement TEXT NOT NULL,
        remediation_action TEXT NOT NULL,
        is_under_triage INTEGER NOT NULL DEFAULT 1,
        has_auto_recovered INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Dispatches diagnostic actions through `pkg/warmup/triage/engine.go` and integrates with `gitmap email warmup triage`.

---

### Improvement 733: One-Key Seed Account Quarantine and Instant Hot-Swap Orchestrator
- **ID:** `SEED-HOTSWAP-733`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Allow instant isolation of burned or flagged seed mailboxes with automated allocation and synchronization of pre-warmed backup seed accounts via a single keystroke.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Modal confirmation dialogue overlaying the main HUD with hot-swap telemetry (`Replacing seed-04 with standby seed-88`).
    - Lipgloss token palette: Modal Frame (`#6366f1`, indigo), Confirmation Key (`#22c55e`, green), Warning Badge (`#f97316`, orange).
    - ANSI Modal:
      ```text
      ┌───────────────── [CONFIRM SEED HOT-SWAP] ─────────────────┐
      │ Flagged Seed:  seed-04@yahoo-test.com (Placement: 12.0%)  │
      │ Standby Seed:  seed-88@yahoo-test.com (Health: 99.4%)     │
      │ Action: Quarantine target and migrate warmup thread quota │
      │                                                           │
      │ Press [Enter] to Execute Swap | Press [Esc] to Cancel     │
      └───────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure CSS modal dialog utilizing `:target` selector without JavaScript dialog polyfills.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-swap`
  - Flags:
    - `--target <email>`: Target seed address to quarantine.
    - `--replacement <email>`: Specific standby seed address (optional; auto-allocates if omitted).
    - `--reason <string>`: Diagnostic reason for quarantine (default: `manual_triage`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedHotSwapEvent struct {
        EventID          string    `json:"eventId"`
        QuarantinedSeed  string    `json:"quarantinedSeed"`
        ReplacementSeed  string    `json:"replacementSeed"`
        SwapReason       string    `json:"swapReason"`
        CompletedAt      time.Time `json:"completedAt"`
        IsSuccessful     bool      `json:"isSuccessful"`
        HasStateMigrated bool      `json:"hasStateMigrated"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedHotSwapAudit (
        event_id TEXT PRIMARY KEY,
        quarantined_seed TEXT NOT NULL,
        replacement_seed TEXT NOT NULL,
        swap_reason TEXT NOT NULL,
        is_successful INTEGER NOT NULL DEFAULT 1,
        has_state_migrated INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Triggered via keybinding `x` in `pkg/warmup/ui/triage_hud.go` and executed by `pkg/warmup/pool/rotator.go`.

---

### Improvement 734: Multi-Target ISP Latency Sparkline Ticker inside TUI Footer
- **ID:** `ISP-SPARKLINE-734`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Embed multi-ISP live round-trip latency sparklines (`Google`, `Microsoft`, `Yahoo`, `Apple`) directly in the TUI status bar using Braille Unicode glyphs.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Footer bar featuring compact sparklines rendered via Unicode Braille characters (`⡀`, `⣀`, `⣄`, `⣤`, `⣦`, `⣶`, `⣷`, `⣿`).
    - Lipgloss token palette: Normal latency (`#38bdf8`), Latency climb (`#fbbf24`), Degradation (`#f43f5e`).
    - Model display snippet:
      ```text
      ISP LATENCY TICKER:
      GGL [⡀⣀⣄⣄⣦⣦⣄ 42ms] | MSFT [⡀⡀⣀⣄⣦⣶⣷ 112ms ⚠] | YHO [⡀⣀⣀⡀⡀⣀⡀ 38ms] | APL [⡀⡀⡀⡀⡀⡀⡀ 24ms]
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline SVG sparklines dynamically embedded using `<polyline>` with `stroke-width="1.5"` and transparent backgrounds.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup latency-ticker`
  - Flags:
    - `--history-points <int>`: Number of sliding telemetry samples to render (default: `16`).
    - `--isp <string>`: Filter ticker to target ISP (default: `all`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ISPLatencySample struct {
        ISPName       string    `json:"ispName"`
        CurrentRTTMs  int64     `json:"currentRttMs"`
        HistoryPoints []int64   `json:"historyPoints"`
        SparklineStr  string    `json:"sparklineStr"`
        IsDegraded    bool      `json:"isDegraded"`
        HasJitterAlert bool     `json:"hasJitterAlert"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ISPLatencyTelemetry (
        sample_id INTEGER PRIMARY KEY AUTOINCREMENT,
        isp_name TEXT NOT NULL,
        rtt_ms INTEGER NOT NULL,
        is_degraded INTEGER NOT NULL DEFAULT 0,
        has_jitter_alert INTEGER NOT NULL DEFAULT 0,
        sampled_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_isp_latency_time ON ISPLatencyTelemetry(isp_name, sampled_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Consumed by `pkg/warmup/ui/components/sparkline.go` and integrated into the global status ticker.

---

### Improvement 735: Interactive Filter and Search Palette with Fuzzy Mailbox Matching
- **ID:** `MAILBOX-SEARCH-735`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Provide an on-the-fly interactive fuzzy search modal (`/` key) allowing operators to jump directly to any mailbox across thousands of entries by address, relay IP, or domain.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Input overlay field at the top of the grid with instant highlight matching and fuzzy score ranking.
    - Lipgloss token palette: Search prompt (`#a855f7`, purple), Matched characters (`#facc15`, yellow bold), Active row (`#1e293b`, dark blue bg).
    - Model view snippet:
      ```text
      ┌─ [FUZZY MAILBOX FILTER] ───────────────────────────────────────────────────┐
      │ Search: outb-seed/                                                        │
      │ > outb-seed-01@domain-a.com   [IP: 198.51.100.10] [Health: 99.1%]          │
      │   outb-seed-02@domain-a.com   [IP: 198.51.100.11] [Health: 98.4%]          │
      │   outb-seed-05@domain-b.com   [IP: 198.51.100.18] [Health: 88.0% ⚠]        │
      │ [Enter: Select] [Ctrl+N/Ctrl+P: Cycle] [Esc: Dismiss]                      │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native client-side search using vanilla HTML input with vanilla inline script utilizing Trie index for sub-millisecond filtering.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup find`
  - Flags:
    - `--query <string>`: Immediate query string for non-interactive output.
    - `--limit <int>`: Max matched results to return (default: `10`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxSearchResult struct {
        MailboxID       string  `json:"mailboxId"`
        EmailAddress    string  `json:"emailAddress"`
        DomainName      string  `json:"domainName"`
        AssignedIP      string  `json:"assignedIp"`
        MatchScore      int     `json:"matchScore"`
        IsExactMatch    bool    `json:"isExactMatch"`
        HasAlertStatus  bool    `json:"hasAlertStatus"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE VIRTUAL TABLE IF NOT EXISTS MailboxFTS USING fts5(
        mailbox_id,
        email_address,
        domain_name,
        assigned_ip
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/ui/search_palette.go` hooking directly into SQLite FTS5 backend.

---

### Improvement 736: Mailbox Lifecycle Stage Progress Bar Visualizer
- **ID:** `WARMUP-STAGE-736`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Display individual mailbox progression across 5 standardized warmup stages (`Stage 0: Cold`, `Stage 1: Ramp`, `Stage 2: Acceleration`, `Stage 3: Mature`, `Stage 4: Autonomous Production`) with milestone markers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Segmented progress bar with UTF-8 block elements (`█`, `▓`, `▒`, `░`) and stage checkpoint tags (`[C]`, `[R]`, `[A]`, `[M]`, `[P]`).
    - Lipgloss token palette: Completed segments (`#10b981`), Current segment (`#38bdf8`), Remaining (`#334155`).
    - Model snippet:
      ```text
      STAGE PROGRESS: [████████████▓▓▓▓▓░░░░░░░░░░░░] 54% (Stage 2: Acceleration)
      Day: 14/30 | Current Target: 350 emails/day | Peak Target: 2,500 emails/day
      ```
  - **Zero-CDN Web Dashboard:**
    - CSS-only multi-segment bar using CSS Flexbox and gradient fills (`linear-gradient(...)`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup stage-status`
  - Flags:
    - `--mailbox <string>`: Target mailbox address.
    - `--detail`: Show milestone audit checklist.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxStageProgress struct {
        MailboxAddress    string `json:"mailboxAddress"`
        CurrentStage      int    `json:"currentStage"`
        StageName         string `json:"stageName"`
        ProgressPercent   int    `json:"progressPercent"`
        DaysActive        int    `json:"daysActive"`
        IsMilestoneMet    bool   `json:"isMilestoneMet"`
        HasPassedSafetyGate bool `json:"hasPassedSafetyGate"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxStageProgression (
        mailbox_address TEXT PRIMARY KEY,
        current_stage INTEGER NOT NULL DEFAULT 0,
        stage_name TEXT NOT NULL DEFAULT 'Cold',
        progress_percent INTEGER NOT NULL DEFAULT 0,
        days_active INTEGER NOT NULL DEFAULT 0,
        is_milestone_met INTEGER NOT NULL DEFAULT 0,
        has_passed_safety_gate INTEGER NOT NULL DEFAULT 0,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Evaluated by `pkg/warmup/stage/evaluator.go` and visualized in `pkg/warmup/ui/stage_view.go`.

---

### Improvement 737: Bulk Pause, Resume, and Throttling Action Bar
- **ID:** `BULK-ACTION-737`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Enable bulk operations across flagged or selected subsets of mailboxes directly from the TUI with visual selection counter and confirmation lock.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Selection mode (`v` to enter visual mode, `space` to toggle selection, `a` for select all in current filter).
    - Bottom command bar showing active selection count and action keys (`[P: Pause]`, `[R: Resume]`, `[T: Throttle 50%]`).
    - Model view snippet:
      ```text
      ┌─ [SELECTED: 14 MAILBOXES] ─────────────────────────────────────────────────┐
      │ Action: [P] Pause Warmup  [R] Resume  [T] Apply 50% Throttle  [C] Clear    │
      │ Warning: Pausing will temporarily park active seed-pool delivery threads.  │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Multi-select checkbox table with pure CSS floating action toolbar appearing when checkboxes are checked.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup bulk-control`
  - Flags:
    - `--action <string>`: Action type (`pause`, `resume`, `throttle`).
    - `--filter-status <string>`: Target mailboxes matching status.
    - `--throttle-percent <int>`: Quota reduction percentage (default: `50`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BulkControlRequest struct {
        TargetAddresses []string `json:"targetAddresses"`
        ActionType      string   `json:"actionType"`
        ThrottlePercent int      `json:"throttlePercent"`
        ExecutedBy      string   `json:"executedBy"`
        IsDryRun        bool     `json:"isDryRun"`
        HasConfirmed    bool     `json:"hasConfirmed"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BulkOperationAudit (
        operation_id TEXT PRIMARY KEY,
        action_type TEXT NOT NULL,
        affected_count INTEGER NOT NULL,
        executed_by TEXT NOT NULL,
        is_dry_run INTEGER NOT NULL DEFAULT 0,
        has_confirmed INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/ui/bulk_action.go` and executed through `pkg/warmup/dispatch/controller.go`.

---

### Improvement 738: Real-Time Seed Reply Rate Breakdown & Gauge Meter
- **ID:** `REPLY-GAUGE-738`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Render a high-precision semi-circular or horizontal gauge displaying real-time seed-to-sender reply rates against configured warmup targets.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-resolution ASCII horizontal needle gauge with threshold bounds (`[Min: 20%] [Target: 35%] [Current: 38.4%]`).
    - Lipgloss token palette: Under-target (`#ef4444`), On-target (`#10b981`), Over-target warning (`#f59e0b`).
    - Model view snippet:
      ```text
      SEED REPLY ENGAGEMENT GAUGE:
      0%  [──────────|──────────◆──────────|──────────] 100%
                    20%        38.4%      60%
      Status: OPTIMAL ENGAGEMENT (Simulated Human Pattern Verified)
      ```
  - **Zero-CDN Web Dashboard:**
    - Standalone SVG semi-circle gauge element with stroke-dasharray transitions driven by CSS variables.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup reply-gauge`
  - Flags:
    - `--domain <string>`: Target domain to inspect.
    - `--window <duration>`: Time window for calculation (default: `24h`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ReplyEngagementMetrics struct {
        DomainName     string  `json:"domainName"`
        TotalSent      int     `json:"totalSent"`
        RepliesReceived int    `json:"repliesReceived"`
        ReplyRate      float64 `json:"replyRate"`
        TargetRate     float64 `json:"targetRate"`
        IsOptimal      bool    `json:"isOptimal"`
        HasReplyDeficit bool   `json:"hasReplyDeficit"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedEngagementMetrics (
        domain_name TEXT PRIMARY KEY,
        total_sent INTEGER NOT NULL DEFAULT 0,
        replies_received INTEGER NOT NULL DEFAULT 0,
        reply_rate REAL NOT NULL DEFAULT 0.0,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_reply_deficit INTEGER NOT NULL DEFAULT 0,
        calculated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Rendered in `pkg/warmup/ui/reply_meter.go` utilizing metrics aggregated from `pkg/warmup/telemetry/reply_aggregator.go`.

---

### Improvement 739: Visual Seed Distribution Pie/Donut Chart Component
- **ID:** `SEED-DONUT-739`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Present an ANSI/Unicode text-based distribution donut chart showing proportional allocation of seed accounts across major providers (Gmail, Outlook, Yahoo, Custom MX).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Terminal donut block rendering using shaded block glyphs with legend badges showing account counts and target quotas.
    - Lipgloss token palette: Google (`#ea4335`, red), Microsoft (`#00a4ef`, blue), Yahoo (`#6001d2`, purple), Other (`#64748b`, slate).
    - Model view snippet:
      ```text
      ┌─ [SEED ALLOCATION BREAKDOWN] ──────────────────────────────────────────────┐
      │   ▄████████▄     ■ Gmail:     45% (90 seeds) [Allocated]                     │
      │  ███▛      ▜███  ■ Outlook:   30% (60 seeds) [Allocated]                     │
      │  ███       ███   ■ Yahoo:     15% (30 seeds) [Allocated]                     │
      │   ▀████████▀     ■ Custom MX: 10% (20 seeds) [Allocated]                     │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline SVG `<circle>` donut chart with CSS `stroke-dasharray` and `stroke-dashoffset` for zero-JS rendering.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-distribution`
  - Flags:
    - `--chart-type <string>`: Presentation format (`donut`, `bar`, `table`) (default: `donut`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedDistributionSlice struct {
        ProviderName string  `json:"providerName"`
        SeedCount    int     `json:"seedCount"`
        Percentage   float64 `json:"percentage"`
        ColorHex     string  `json:"colorHex"`
        IsBalanced   bool    `json:"isBalanced"`
        HasUnderflow bool    `json:"hasUnderflow"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedPoolDistribution (
        provider_name TEXT PRIMARY KEY,
        seed_count INTEGER NOT NULL,
        percentage REAL NOT NULL,
        is_balanced INTEGER NOT NULL DEFAULT 1,
        has_underflow INTEGER NOT NULL DEFAULT 0,
        synced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/warmup/ui/charts/donut.go` and displayed in the main seed pool overview screen.

---

### Improvement 740: Exportable TUI Snapshot Report Generator (Markdown/ANSI/SVG)
- **ID:** `TUI-SNAPSHOT-740`
- **Cluster:** 74 (Interactive Bubble Tea Mailbox Health Matrix & Instant Seed-Pool Triage HUD)
- **Goal:** Enable one-key instant export of the current TUI view buffer to formatted Markdown, raw ANSI text, or standalone vector SVG for sharing in pull requests and incident tickets.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Quick-action notification banner (`Snapshot captured to .ai-memory/reports/warmup-snap-20261027.svg`).
    - Model view snippet:
      ```text
      [INFO: Snapshot saved: 02-spec/21-app/reports/warmup-triage-snap-731.md]
      ```
  - **Zero-CDN Web Dashboard:**
    - Download button injecting standalone SVG embedding all inline styles and system font definitions (`ui-monospace, monospace`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup snapshot`
  - Flags:
    - `--format <string>`: Output format (`svg`, `ansi`, `markdown`, `json`) (default: `svg`).
    - `--output-path <string>`: Relative filesystem destination path.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SnapshotExportMeta struct {
        ExportID     string    `json:"exportId"`
        FormatType   string    `json:"formatType"`
        FilePath     string    `json:"filePath"`
        ByteSize     int64     `json:"byteSize"`
        CapturedAt   time.Time `json:"capturedAt"`
        IsCompressed bool      `json:"isCompressed"`
        HasExportSucceeded bool `json:"hasExportSucceeded"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SnapshotExportHistory (
        export_id TEXT PRIMARY KEY,
        format_type TEXT NOT NULL,
        file_path TEXT NOT NULL,
        byte_size INTEGER NOT NULL,
        is_compressed INTEGER NOT NULL DEFAULT 0,
        has_export_succeeded INTEGER NOT NULL DEFAULT 1,
        captured_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/ui/snapshot.go` with keybinding `Ctrl+S` registered in Bubble Tea event loop.

---

# Cluster 75: Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler (Items 741–750)

---

### Improvement 741: SMTP 451/421 Greylisting Signature Pattern Matcher
- **ID:** `GREYLIST-DETECT-741`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Implement an autonomous regular expression and token parser for SMTP 4xx response strings that accurately categorizes temporary rejections as greylisting vs rate limiting vs socket timeout.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Real-time diagnostic stream showing parsed SMTP responses with regex match classification badges (`[GREYLIST]`, `[RATE_LIMIT]`, `[TEMP_FAIL]`).
    - Lipgloss token palette: Greylist badge (`#a855f7`, purple), Rate limit badge (`#f59e0b`, amber), Match highlight (`#38bdf8`, cyan).
    - Model view snippet:
      ```text
      ┌─ [SMTP DIAGNOSTIC RESPONSE PARSER] ────────────────────────────────────────┐
      │ 2026-10-27 21:52:01 | mx.google.com -> 451 4.7.1 Service unavailable:      │
      │   Pattern: `4\.7\.1.*greylist|try again later` -> [CLASSIFIED: GREYLIST]     │
      │   Calculated Hold Window: 15m 00s (Auto-Deferral Activated)                  │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live-updating log table with color-coded classification tags and searchable regex filter.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup parse-smtp-error`
  - Flags:
    - `--raw-error <string>`: Test arbitrary SMTP response string against greylist engine.
    - `--strict-mode`: Require exact ISP code matching instead of generic fallback.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GreylistPatternMatch struct {
        RawResponse      string `json:"rawResponse"`
        Classification   string `json:"classification"`
        MatchedRegex     string `json:"matchedRegex"`
        TargetISP        string `json:"targetIsp"`
        RecommendedHoldS int    `json:"recommendedHoldS"`
        IsGreylisted     bool   `json:"isGreylisted"`
        HasKnownISPCode  bool   `json:"hasKnownIspCode"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GreylistSignatures (
        signature_id TEXT PRIMARY KEY,
        isp_identifier TEXT NOT NULL,
        regex_pattern TEXT NOT NULL,
        recommended_hold_seconds INTEGER NOT NULL DEFAULT 900,
        is_greylisted INTEGER NOT NULL DEFAULT 1,
        has_known_isp_code INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/warmup/greylist/parser.go` and is called by the delivery dispatcher prior to scheduling retries.

---

### Improvement 742: ISP-Specific Adaptive Exponential Backoff with Jitter Injector
- **ID:** `BACKOFF-JITTER-742`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Calculate adaptive retry intervals using exponential backoff with full randomized jitter tailored specifically to destination MX server behavioral profiles.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Mathematical calculation display showing base interval, multiplier, randomized jitter window, and scheduled retry timestamp.
    - Lipgloss token palette: Formula labels (`#94a3b8`), Calculated duration (`#22c55e`), Jitter range (`#f43f5e`).
    - Model view snippet:
      ```text
      ADAPTIVE RETRY SCHEDULER:
      Target: mx1.outlook.com | Attempt: 2/5 | Base: 300s | Multiplier: 2.0
      Jitter Range: [-45s, +45s] | Final Calculated Wait: 9m 14s (ETA: 22:01:15)
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline mathematical plot rendering retry curves across multiple attempts using SVG line paths.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup calculate-backoff`
  - Flags:
    - `--attempt <int>`: Current retry attempt index (1-based).
    - `--isp <string>`: Target ISP name (`google`, `microsoft`, `yahoo`, `generic`).
    - `--base-seconds <int>`: Override base backoff interval (default: `300`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BackoffCalculation struct {
        TargetMX        string        `json:"targetMx"`
        AttemptNumber   int           `json:"attemptNumber"`
        BaseDuration    time.Duration `json:"baseDuration"`
        JitterOffset    time.Duration `json:"jitterOffset"`
        TotalWaitTime   time.Duration `json:"totalWaitTime"`
        ScheduledTime   time.Time     `json:"scheduledTime"`
        IsMaxCapReached bool          `json:"isMaxCapReached"`
        HasJitterApplied bool         `json:"hasJitterApplied"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AdaptiveRetrySchedule (
        schedule_id TEXT PRIMARY KEY,
        message_id TEXT NOT NULL,
        target_mx TEXT NOT NULL,
        attempt_number INTEGER NOT NULL,
        total_wait_seconds INTEGER NOT NULL,
        scheduled_for DATETIME NOT NULL,
        is_max_cap_reached INTEGER NOT NULL DEFAULT 0,
        has_jitter_applied INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_retry_scheduled_for ON AdaptiveRetrySchedule(scheduled_for);
    ```
- **Antigravity / GitMap Integration Point:**
  - Logic resides in `pkg/warmup/greylist/backoff.go` interacting with GitMap's task queue runner.

---

### Improvement 743: Triplet Identification & Greylist Cache Registry (Sender, Recipient, Client IP)
- **ID:** `TRIPLET-REGISTRY-743`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Maintain an in-memory and SQLite-backed registry of active greylisting triplets `(ClientIP, SenderAddress, RecipientAddress)` to avoid redundant immediate attempts.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Triplet cache table showing remaining embargo countdown timers and verification flags.
    - Lipgloss token palette: Active Embargo (`#f59e0b`), Cleared/Passed (`#10b981`), Expired (`#64748b`).
    - Model view snippet:
      ```text
      ┌─ [GREYLIST TRIPLET REGISTRY] ──────────────────────────────────────────────┐
      │ CLIENT IP     │ SENDER                │ RECIPIENT            │ EMBARGO REM │
      │ 198.51.100.12 │ out-01@domain.io      │ seed-04@yahoo.com    │ 04m 22s     │
      │ 198.51.100.15 │ out-02@domain.io      │ seed-18@outlook.com  │ 11m 40s     │
      │ [Active Triplets: 2] [Cleared Today: 48] [Auto-Purge Expired: Enabled]     │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Live searchable HTML table with client-side countdown clocks rendering without external dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup triplets`
  - Flags:
    - `--active-only`: Filter triplets with remaining embargo time (default: `true`).
    - `--purge-cleared`: Remove expired cleared records from SQLite registry.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GreylistTriplet struct {
        TripletHash      string    `json:"tripletHash"`
        ClientIP         string    `json:"clientIp"`
        SenderAddress    string    `json:"senderAddress"`
        RecipientAddress string    `json:"recipientAddress"`
        FirstAttemptAt   time.Time `json:"firstAttemptAt"`
        EmbargoExpiresAt time.Time `json:"embargoExpiresAt"`
        IsEmbargoActive  bool      `json:"isEmbargoActive"`
        HasPassedCheck   bool      `json:"hasPassedCheck"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GreylistTripletCache (
        triplet_hash TEXT PRIMARY KEY,
        client_ip TEXT NOT NULL,
        sender_address TEXT NOT NULL,
        recipient_address TEXT NOT NULL,
        first_attempt_at DATETIME NOT NULL,
        embargo_expires_at DATETIME NOT NULL,
        is_embargo_active INTEGER NOT NULL DEFAULT 1,
        has_passed_check INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_triplet_embargo ON GreylistTripletCache(embargo_expires_at);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/warmup/greylist/triplet_store.go` with automatic LRU memory caching.

---

### Improvement 744: Bubble Tea Interactive Retry Timeline and Queue Visualizer
- **ID:** `RETRY-TIMELINE-744`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Display a graphical timeline in Bubble Tea visualizing pending retries scheduled across future minutes and hours with density heat indicators.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Horizontal timeline bar representing the next 60 minutes in 5-minute buckets, with bar heights denoting scheduled dispatch load.
    - Lipgloss token palette: Normal dispatch (`#38bdf8`), High queue density (`#f59e0b`), Throttling threshold line (`#ef4444`).
    - Model view snippet:
      ```text
      RETRY QUEUE DISPATCH TIMELINE (NEXT 60 MINUTES):
      NOW  +5m   +10m  +15m  +20m  +25m  +30m  +35m  +40m  +45m  +50m  +55m  +60m
      [█]  [▄]   [ ]   [▆]   [█]   [▃]   [ ]   [ ]   [▂]   [ ]   [ ]   [ ]   [ ]
      Items: 12   4     0     18    22    6     0     0     3     0     0     0     0
      Peak Capacity: 22/50 msgs/bucket | Next Immediate Dispatch in 01m 14s
      ```
  - **Zero-CDN Web Dashboard:**
    - Responsive SVG bar histogram rendered directly with semantic HTML tags.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup retry-queue`
  - Flags:
    - `--buckets <int>`: Number of timeline slices to compute (default: `12`).
    - `--bucket-duration <duration>`: Duration of each slice (default: `5m`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type RetryTimelineBucket struct {
        BucketIndex       int       `json:"bucketIndex"`
        StartTime         time.Time `json:"startTime"`
        EndTime           time.Time `json:"endTime"`
        ScheduledCount    int       `json:"scheduledCount"`
        IsPeakBucket      bool      `json:"isPeakBucket"`
        HasExceededSafety bool      `json:"hasExceededSafety"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RetryTimelineBuckets (
        bucket_id INTEGER PRIMARY KEY AUTOINCREMENT,
        start_time DATETIME NOT NULL,
        end_time DATETIME NOT NULL,
        scheduled_count INTEGER NOT NULL DEFAULT 0,
        is_peak_bucket INTEGER NOT NULL DEFAULT 0,
        has_exceeded_safety INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Component coded in `pkg/warmup/ui/timeline_view.go` fed by the retry scheduler.

---

### Improvement 745: MX Fallback and Alternative Route Balancer for Persistent Greylisting
- **ID:** `MX-FALLBACK-745`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Automatically rotate to secondary or tertiary MX records when the primary MX host consistently greylists connection attempts past the maximum backoff threshold.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Route switching log showing MX priority fallbacks (`mx1.domain.com [FAIL] -> mx2.domain.com [ACTIVE]`).
    - Lipgloss token palette: Failed MX (`#ef4444`), Active fallback (`#10b981`), Priority tag (`#64748b`).
    - Model view snippet:
      ```text
      ROUTE BALANCER: Destination Domain: enterprise-seed.com
      ├─ [Pri 10] mx1.enterprise-seed.com: FAILED (Greylisted 3x consecutive)
      └─ [Pri 20] mx2.enterprise-seed.com: CONNECTED (TLS Handshake OK, Accepted)
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline SVG routing tree showing primary vs fallback MX nodes with status pills.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mx-route`
  - Flags:
    - `--domain <string>`: Target destination domain to resolve and test.
    - `--force-fallback`: Force routing to secondary MX regardless of primary state.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MXRouteOption struct {
        HostName      string `json:"hostName"`
        Preference    int    `json:"preference"`
        IPAddress     string `json:"ipAddress"`
        ConsecutiveFails int `json:"consecutiveFails"`
        IsActiveRoute bool   `json:"isActiveRoute"`
        HasPassedTLS  bool   `json:"hasPassedTls"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MXRoutePool (
        domain_name TEXT NOT NULL,
        host_name TEXT NOT NULL,
        preference INTEGER NOT NULL,
        ip_address TEXT NOT NULL,
        is_active_route INTEGER NOT NULL DEFAULT 0,
        has_passed_tls INTEGER NOT NULL DEFAULT 1,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY(domain_name, host_name)
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Encapsulated in `pkg/warmup/greylist/mx_fallback.go` integrated with DNS resolver cache.


---

### Improvement 746: Automated Recipient Domain Throttling Governor
- **ID:** `THROTTLE-GOV-746`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Dynamically enforce concurrency limits and messages-per-minute caps on a per-destination ISP basis upon detecting repetitive 4xx throttling spikes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Live governor gauge showing token bucket depth, active concurrency, and enforced dispatch limit.
    - Lipgloss token palette: Normal flow (`#10b981`), Governor active (`#f59e0b`), Saturated/Clamped (`#ef4444`).
    - Model view snippet:
      ```text
      ISP CONCURRENCY GOVERNOR: mx.yahoo.com
      Token Bucket: [████████░░] 8/10 tokens | Concurrency Limit: 2 connections (CLAMPED)
      State: ACTIVE_GOVERNANCE (Rate lowered from 60 msg/min to 15 msg/min)
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure HTML/CSS progress bar reflecting token bucket levels without JavaScript charting libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup governor`
  - Flags:
    - `--isp <string>`: Target ISP domain identifier.
    - `--max-concurrency <int>`: Set maximum concurrent SMTP sockets (default: `2`).
    - `--rate-limit <int>`: Messages dispatched per minute cap (default: `20`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainThrottleRule struct {
        TargetDomain       string `json:"targetDomain"`
        MaxConcurrency     int    `json:"maxConcurrency"`
        MessagesPerMinute  int    `json:"messagesPerMinute"`
        IsGovernanceActive bool   `json:"isGovernanceActive"`
        HasRecentSpike     bool   `json:"hasRecentSpike"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainThrottleRules (
        target_domain TEXT PRIMARY KEY,
        max_concurrency INTEGER NOT NULL DEFAULT 5,
        messages_per_minute INTEGER NOT NULL DEFAULT 60,
        is_governance_active INTEGER NOT NULL DEFAULT 0,
        has_recent_spike INTEGER NOT NULL DEFAULT 0,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Governed by `pkg/warmup/greylist/governor.go` connected directly to the socket connection pool.

---

### Improvement 747: Greylisting Embargo Countdown HUD Bar
- **ID:** `EMBARGO-BAR-747`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Provide a high-precision countdown bar in the TUI tracking the precise embargo expiration moment across all deferred messages.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Active countdown meter updating at 1Hz frequency with minute and second precision.
    - Lipgloss token palette: Clock icon (`#38bdf8`), Expiring soon (`#10b981`), Long wait (`#64748b`).
    - Model snippet:
      ```text
      ACTIVE GREYLIST EMBARGOES:
      [ID: msg-9921] -> seed-14@comcast.net | Wait: 06m 12s [████████████░░░░] 68%
      [ID: msg-9924] -> seed-02@verizon.net | Wait: 12m 44s [████░░░░░░░░░░░░] 22%
      ```
  - **Zero-CDN Web Dashboard:**
    - Lightweight `<time>` tags formatted with pure CSS keyframe progress bar animation.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup embargo-hud`
  - Flags:
    - `--watch`: Persistent interactive terminal countdown monitor.
    - `--min-wait <duration>`: Filter embargoes exceeding duration.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmbargoTracker struct {
        MessageID       string        `json:"messageId"`
        TargetRecipient string        `json:"targetRecipient"`
        TotalEmbargo    time.Duration `json:"totalEmbargo"`
        RemainingWait   time.Duration `json:"remainingWait"`
        IsImminent      bool          `json:"isImminent"`
        HasAutoFired    bool          `json:"hasAutoFired"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ActiveEmbargoes (
        message_id TEXT PRIMARY KEY,
        target_recipient TEXT NOT NULL,
        total_embargo_seconds INTEGER NOT NULL,
        expires_at DATETIME NOT NULL,
        is_imminent INTEGER NOT NULL DEFAULT 0,
        has_auto_fired INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Component in `pkg/warmup/ui/embargo_hud.go` refreshing on Bubble Tea tick message `tea.Tick`.

---

### Improvement 748: Dead-Letter Queue (DLQ) Inspector and Re-Drive Modal
- **ID:** `DLQ-REDRIVE-748`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Display exhausted messages that failed all adaptive retry attempts in a dead-letter quarantine inspector with root cause summaries and one-key redrive capability.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Quarantined message inspection table with error stack preview and keybindings for manual redrive or permanent purge.
    - Lipgloss token palette: DLQ Warning (`#ef4444`), Redrive action (`#38bdf8`), Purge action (`#94a3b8`).
    - Model view snippet:
      ```text
      ┌─ [DEAD-LETTER QUARANTINE: 4 EXHAUSTED MESSAGES] ───────────────────────────┐
      │ MSG ID   │ SENDER            │ DESTINATION       │ FAILURE REASON          │
      │ ▶ msg-81 │ out-01@domain.io  │ seed-09@aol.com   │ 550 5.1.1 User unknown  │
      │   msg-94 │ out-03@domain.io  │ seed-44@sbc.net   │ Max retries exceeded    │
      │ [Space: Inspect Payload] [r: Re-Drive Selected] [d: Discard] [q: Quit]     │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Accessible HTML table with `<details>` collapsible raw SMTP dialogue transcript.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup dlq`
  - Flags:
    - `--redrive-all`: Force immediate requeuing of all qualifying DLQ items.
    - `--purge-older-than <duration>`: Discard stale DLQ messages (default: `72h`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DeadLetterRecord struct {
        MessageID        string    `json:"messageId"`
        SenderAddress    string    `json:"senderAddress"`
        RecipientAddress string    `json:"recipientAddress"`
        FinalError       string    `json:"finalError"`
        RetryCount       int       `json:"retryCount"`
        QuarantinedAt    time.Time `json:"quarantinedAt"`
        IsRedrivable     bool      `json:"isRedrivable"`
        HasPermanentFail bool      `json:"hasPermanentFail"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DeadLetterMessages (
        message_id TEXT PRIMARY KEY,
        sender_address TEXT NOT NULL,
        recipient_address TEXT NOT NULL,
        final_error TEXT NOT NULL,
        retry_count INTEGER NOT NULL,
        is_redrivable INTEGER NOT NULL DEFAULT 1,
        has_permanent_fail INTEGER NOT NULL DEFAULT 0,
        quarantined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/warmup/greylist/dlq.go` and surfaced through `gitmap email warmup dlq`.

---

### Improvement 749: Greylisting Recovery Success Rate Analytics Heatmap
- **ID:** `RECOVERY-HEATMAP-749`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Visualize percentage recovery rates following initial greylisting events across destination domains over 24-hour periods using a calendar-style heat grid.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - 24-hour horizontal heatmap grid using UTF-8 shades (`░` 0%, `▒` 25%, `▓` 75%, `█` 100% recovery).
    - Lipgloss token palette: Cold/Zero (`#1e293b`), Medium recovery (`#38bdf8`), Full recovery (`#10b981`).
    - Model view snippet:
      ```text
      GREYLIST RECOVERY SUCCESS (HOURLY PAST 24H):
      00:00 [████████████████████████] 23:00
      Yahoo:   [████▓▓████████████████] Avg: 96.2% Recovery
      Outlook: [▓▓▒▒▓▓▓▓████████▓▓▓▓██] Avg: 84.1% Recovery
      Comcast: [██████████████████████] Avg: 99.8% Recovery
      ```
  - **Zero-CDN Web Dashboard:**
    - SVG heat matrix with native HTML tooltip styling via CSS `attr(data-title)`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup recovery-heatmap`
  - Flags:
    - `--days <int>`: History span in days (default: `1`).
    - `--isp <string>`: Target ISP filter (default: `all`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HourlyRecoveryMetric struct {
        HourIndex      int     `json:"hourIndex"`
        ISPName        string  `json:"ispName"`
        TotalGreylisted int    `json:"totalGreylisted"`
        TotalRecovered int     `json:"totalRecovered"`
        RecoveryRate   float64 `json:"recoveryRate"`
        IsOptimal      bool    `json:"isOptimal"`
        HasSevereDrop  bool    `json:"hasSevereDrop"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HourlyRecoveryStats (
        stat_id INTEGER PRIMARY KEY AUTOINCREMENT,
        isp_name TEXT NOT NULL,
        hour_slot DATETIME NOT NULL,
        total_greylisted INTEGER NOT NULL DEFAULT 0,
        total_recovered INTEGER NOT NULL DEFAULT 0,
        recovery_rate REAL NOT NULL DEFAULT 0.0,
        is_optimal INTEGER NOT NULL DEFAULT 1,
        has_severe_drop INTEGER NOT NULL DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_recovery_slot ON HourlyRecoveryStats(isp_name, hour_slot);
    ```
- **Antigravity / GitMap Integration Point:**
  - Computed by `pkg/warmup/telemetry/recovery.go` and rendered by `pkg/warmup/ui/heat_view.go`.

---

### Improvement 750: Autonomous Dynamic Hold-Time Tuner with Reinforcement Feedback
- **ID:** `HOLD-TUNER-750`
- **Cluster:** 75 (Dynamic Greylisting Auto-Detection & Adaptive Exponential Retry Scheduler)
- **Goal:** Autonomously adjust the default first-retry hold duration per destination MX based on statistical historical clearing windows to minimize delivery deferral time.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Parameter convergence display showing previous hold duration vs statistically optimized new hold setting.
    - Lipgloss token palette: Tuned parameter (`#a855f7`), Delta indicator (`#22c55e`), Baseline (`#64748b`).
    - Model view snippet:
      ```text
      DYNAMIC HOLD-TIME TUNER (REINFORCEMENT ACTIVE):
      Target: mx.aol.com | Previous Hold: 15m 00s -> New Tuned Hold: 09m 20s (-37.7%)
      Confidence: 98.4% (Based on 142 cleared deliveries, min clearance window: 08m 10s)
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline HTML telemetry card with mathematical convergence graph rendered in pure SVG.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup tune-hold`
  - Flags:
    - `--auto-apply`: Commit tuned hold parameters to SQLite database (default: `true`).
    - `--min-samples <int>`: Minimum sample size required to tune (default: `30`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TunedHoldParam struct {
        TargetMX         string        `json:"targetMx"`
        BaselineHold     time.Duration `json:"baselineHold"`
        OptimizedHold    time.Duration `json:"optimizedHold"`
        SampleCount      int           `json:"sampleCount"`
        ConfidenceScore  float64       `json:"confidenceScore"`
        IsApplied        bool          `json:"isApplied"`
        HasHighVariance  bool          `json:"hasHighVariance"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TunedHoldParameters (
        target_mx TEXT PRIMARY KEY,
        baseline_hold_seconds INTEGER NOT NULL,
        optimized_hold_seconds INTEGER NOT NULL,
        sample_count INTEGER NOT NULL,
        confidence_score REAL NOT NULL,
        is_applied INTEGER NOT NULL DEFAULT 1,
        has_high_variance INTEGER NOT NULL DEFAULT 0,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/warmup/greylist/tuner.go` scheduled via GitMap background daemon.

---

# Cluster 76: Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor (Items 751–760)

---

### Improvement 751: Local Offline HTML Content Security Policy (CSP) Strict Auditor
- **ID:** `CSP-AUDITOR-751`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Scan outbound warmup email templates locally to detect external scripts, dynamic stylesheets, tracking beacons, and unauthorized external assets, enforcing strict Zero-CDN isolation.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Template audit checklist with violation tags (`[FAIL: External Script]`, `[FAIL: Remote Font]`, `[PASS: Clean Inlined CSS]`).
    - Lipgloss token palette: Passed check (`#10b981`), Security violation (`#ef4444`), Audit summary (`#38bdf8`).
    - Model view snippet:
      ```text
      ┌─ [EMAIL CONTENT SECURITY AUDIT (ZERO-CDN ENFORCEMENT)] ────────────────────┐
      │ Template: welcome_warmup_v2.html                                           │
      │ [✓] Script Tags: NONE DETECTED (Strict Static DOM)                         │
      │ [✗] External CSS: VIOLATION DETECTED (`fonts.googleapis.com` referenced)   │
      │ [✓] Tracking Pixels: Valid Internal Seed Signature                        │
      │ Result: 1 VIOLATION FOUND | Security Score: 82/100 [BLOCKED FROM DISPATCH]│
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - HTML summary card with line-numbered code view highlighting forbidden tags in red.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup audit-csp`
  - Flags:
    - `--template <path>`: Relative filesystem path to HTML template file.
    - `--strict`: Disallow all external assets without exception (default: `true`).
    - `--fix-inlining`: Automatically strip or inline offending resources.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CSPViolation struct {
        TagName      string `json:"tagName"`
        Attribute    string `json:"attribute"`
        TargetURI    string `json:"targetUri"`
        LineNumber   int    `json:"lineNumber"`
        Severity     string `json:"severity"`
        IsFatal      bool   `json:"isFatal"`
        HasAutoFix   bool   `json:"hasAutoFix"`
    }

    type CSPAuditReport struct {
        TemplatePath string         `json:"templatePath"`
        SecurityScore int           `json:"securityScore"`
        Violations   []CSPViolation `json:"violations"`
        IsPassed     bool           `json:"isPassed"`
        HasWarning   bool           `json:"hasWarning"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailCSPAuditLog (
        audit_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        security_score INTEGER NOT NULL,
        is_passed INTEGER NOT NULL DEFAULT 1,
        has_warning INTEGER NOT NULL DEFAULT 0,
        violation_count INTEGER NOT NULL DEFAULT 0,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded in `pkg/email/security/csp_auditor.go` as a mandatory pre-dispatch check.

---

### Improvement 752: Inline CSS Optimizer and Zero-CDN Embedded Font Stripper
- **ID:** `CSS-INLINER-752`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Automatically parse HTML `<style>` tags and remote `@import` directives, strip external webfonts, and inline CSS properties directly into HTML elements with native CSS fallback stacks.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Optimization transformation log displaying original size vs inlined size and stripped resources.
    - Lipgloss token palette: Optimization delta (`#10b981`), Stripped tag (`#f59e0b`), Size badge (`#64748b`).
    - Model view snippet:
      ```text
      CSS INLINING ENGINE:
      Original: 48.2 KB (External CSS & Fonts) -> Inlined: 18.4 KB (-61.8%)
      Stripped: 2 external font links | Replaced with `system-ui, -apple-system, sans-serif`
      ```
  - **Zero-CDN Web Dashboard:**
    - Side-by-side diff view of raw HTML vs inlined HTML using semantic `<pre><code>` blocks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup inline-css`
  - Flags:
    - `--input <path>`: Source HTML file path.
    - `--output <path>`: Target output file path.
    - `--minify`: Strip excessive whitespace and comments (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type InliningSummary struct {
        OriginalBytes int  `json:"originalBytes"`
        FinalBytes    int  `json:"finalBytes"`
        StrippedRules int  `json:"strippedRules"`
        IsInlined     bool `json:"isInlined"`
        HasFontStack  bool `json:"hasFontStack"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InliningOperationLog (
        operation_id TEXT PRIMARY KEY,
        input_path TEXT NOT NULL,
        original_bytes INTEGER NOT NULL,
        final_bytes INTEGER NOT NULL,
        is_inlined INTEGER NOT NULL DEFAULT 1,
        has_font_stack INTEGER NOT NULL DEFAULT 1,
        processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/security/inliner.go` utilizing standard Go `golang.org/x/net/html`.

---

### Improvement 753: Responsive Mobile/Desktop Dual-Viewport Terminal Previewer
- **ID:** `VIEWPORT-PREVIEW-753`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Render a dual-pane ASCII/ANSI responsive email preview in Bubble Tea simulating 375px mobile screen width alongside 650px desktop email client layouts.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Split terminal screen displaying mobile portrait viewport on the left and full desktop viewport on the right, rendering converted plain text layout and ASCII frame borders.
    - Lipgloss token palette: Mobile frame (`#38bdf8`), Desktop frame (`#a855f7`), Content text (`#f8fafc`).
    - Model view snippet:
      ```text
      ┌─ [MOBILE PREVIEW: 375px] ───────┐ ┌─ [DESKTOP PREVIEW: 650px] ─────────────────┐
      │ ┌─────────────────────────────┐ │ │ ┌─────────────────────────────────────────┐ │
      │ │ Warmup Verification Update  │ │ │ │ Warmup Verification Update              │ │
      │ │ Your warmup thread is active│ │ │ │ Your warmup thread is active. Verified! │ │
      │ │ [Verify Seed Delivery]      │ │ │ │ [ Verify Seed Delivery ]                │ │
      │ └─────────────────────────────┘ │ │ └─────────────────────────────────────────┘ │
      │ [m: Toggle Mobile] [d: Toggle Desktop] [Tab: Switch Pane] [q: Exit Preview]   │
      └─────────────────────────────────┘ └───────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Native CSS container queries (`@container`) simulating mobile and desktop viewports side-by-side.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup preview`
  - Flags:
    - `--template <path>`: Path to email template.
    - `--mode <string>`: Initial view mode (`split`, `mobile`, `desktop`) (default: `split`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ViewportConfig struct {
        MobileWidth   int  `json:"mobileWidth"`
        DesktopWidth  int  `json:"desktopWidth"`
        IsSplitActive bool `json:"isSplitActive"`
        HasReflowed   bool `json:"hasReflowed"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ViewportPreferences (
        user_id TEXT PRIMARY KEY,
        mobile_width INTEGER NOT NULL DEFAULT 375,
        desktop_width INTEGER NOT NULL DEFAULT 650,
        is_split_active INTEGER NOT NULL DEFAULT 1,
        has_reflowed INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Surfaced via `pkg/warmup/ui/previewer.go` integrated with Bubble Tea view matrix.

---

### Improvement 754: Dark Mode / Light Mode Contrast Safety Validator
- **ID:** `CONTRAST-VALIDATE-754`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Calculate WCAG 2.1 AA/AAA color contrast ratios across background and text colors under simulated email client auto-inverted dark modes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Contrast analysis table showing foreground hex, background hex, calculated ratio, and Pass/Fail badge.
    - Lipgloss token palette: Contrast Pass (`#10b981`), Contrast Fail (`#ef4444`), Ratio tag (`#38bdf8`).
    - Model view snippet:
      ```text
      ┌─ [DARK/LIGHT CONTRAST AUDIT] ──────────────────────────────────────────────┐
      │ MODE         │ FG COLOR │ BG COLOR │ RATIO   │ WCAG AA │ WCAG AAA          │
      │ Light Normal │ #1e293b  │ #ffffff  │ 12.6:1  │ [PASS]  │ [PASS]            │
      │ Inverted Dark│ #94a3b8  │ #0f172a  │ 6.8:1   │ [PASS]  │ [FAIL]            │
      │ Inverted Btn │ #ffffff  │ #38bdf8  │ 2.9:1   │ [FAIL]  │ [FAIL] (⚠ Review) │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Inline HTML swatch table with CSS color-scheme simulation (`color-scheme: light dark`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup audit-contrast`
  - Flags:
    - `--template <path>`: Target HTML template file.
    - `--min-ratio <float>`: Minimum acceptable contrast ratio (default: `4.5`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ColorContrastResult struct {
        Selector        string  `json:"selector"`
        ForegroundHex   string  `json:"foregroundHex"`
        BackgroundHex   string  `json:"backgroundHex"`
        CalculatedRatio float64 `json:"calculatedRatio"`
        IsPassedAA      bool    `json:"isPassedAa"`
        IsPassedAAA     bool    `json:"isPassedAaa"`
        HasInvertHazard bool    `json:"hasInvertHazard"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ContrastAuditHistory (
        audit_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        lowest_ratio REAL NOT NULL,
        is_passed_aa INTEGER NOT NULL DEFAULT 1,
        is_passed_aaa INTEGER NOT NULL DEFAULT 0,
        has_invert_hazard INTEGER NOT NULL DEFAULT 0,
        audited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Executed by `pkg/email/security/contrast.go` and surfaced in terminal preview.

---

### Improvement 755: Embedded Base64 Image Quarantine and Size Guard
- **ID:** `BASE64-GUARD-755`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Detect inline `data:image/...;base64` payloads in email bodies, verifying size constraints (<50KB) and converting oversized payloads to multipart CID attachments.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Inline image scan progress list with byte sizes, encoding format, and conversion decision flags.
    - Lipgloss token palette: Inline OK (`#10b981`), Quarantine/Oversized (`#f59e0b`), Conversion tag (`#38bdf8`).
    - Model view snippet:
      ```text
      IMAGE PAYLOAD SCANNER:
      [IMG 1] data:image/png;base64 (12.4 KB) -> [STATUS: OK - INLINE EMBEDDED]
      [IMG 2] data:image/jpeg;base64 (148.2 KB) -> [QUARANTINE: Auto-Converted to CID Attachment]
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure HTML list showing image metadata, dimensions, and computed SHA-256 digests.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup guard-images`
  - Flags:
    - `--max-inline-kb <int>`: Max allowed base64 size before auto-converting to CID (default: `50`).
    - `--template <path>`: Source HTML template.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ImagePayloadMeta struct {
        ImageID         string `json:"imageId"`
        MimeType        string `json:"mimeType"`
        ByteSize        int64  `json:"byteSize"`
        ContentID       string `json:"contentId"`
        IsInlineAllowed bool   `json:"isInlineAllowed"`
        HasCIDConverted bool   `json:"hasCidConverted"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ImagePayloadAudit (
        image_id TEXT PRIMARY KEY,
        template_path TEXT NOT NULL,
        byte_size INTEGER NOT NULL,
        content_id TEXT NOT NULL,
        is_inline_allowed INTEGER NOT NULL DEFAULT 1,
        has_cid_converted INTEGER NOT NULL DEFAULT 0,
        scanned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/email/security/image_guard.go` invoked during email MIME assembly.

---

### Improvement 756: Zero-External-Font Fallback Font Metric Matcher
- **ID:** `FONT-MATCHER-756`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Ensure font-family stacks use metrically-matched standard local system fonts (`system-ui`, `Segoe UI`, `Roboto`, `Helvetica Neue`, `sans-serif`) to prevent layout shifts across clients without external CDN font requests.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Font family evaluation card showing original font stack, detected missing local fallbacks, and synthesized safe fallback stack.
    - Lipgloss token palette: Verified stack (`#10b981`), Substituted font (`#38bdf8`), Shift warning (`#f59e0b`).
    - Model view snippet:
      ```text
      FONT METRIC MATCHER:
      Input Stack:  "Inter", sans-serif (Inter not available on local Windows/Outlook)
      Safe Stack:   "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif
      Layout Shift Risk: LOW (< 1.2% bounding box deviation) [STATUS: OPTIMIZED]
      ```
  - **Zero-CDN Web Dashboard:**
    - Comparative typography specimen rendered with native system fonts in pure HTML.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup sanitize-fonts`
  - Flags:
    - `--template <path>`: Target HTML template file.
    - `--apply-safe-stack`: Automatically replace missing webfont stacks (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FontStackProfile struct {
        OriginalStack    string  `json:"originalStack"`
        NormalizedStack  string  `json:"normalizedStack"`
        LayoutShiftRatio float64 `json:"layoutShiftRatio"`
        IsZeroCDN        bool    `json:"isZeroCdn"`
        HasFallbackMatch bool    `json:"hasFallbackMatch"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FontStackProfiles (
        profile_id TEXT PRIMARY KEY,
        original_stack TEXT NOT NULL,
        normalized_stack TEXT NOT NULL,
        layout_shift_ratio REAL NOT NULL,
        is_zero_cdn INTEGER NOT NULL DEFAULT 1,
        has_fallback_match INTEGER NOT NULL DEFAULT 1,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Processed in `pkg/email/security/fonts.go` during pre-render normalization.

---

### Improvement 757: Spam Assassin & Rspamd Offline Rules Evaluation Engine
- **ID:** `SPAM-RULES-757`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Run a local, offline evaluation of common SpamAssassin and Rspamd scoring rules (`MISSING_MIMEOLE`, `HTML_MESSAGE`, `MIME_HTML_ONLY`, `BAYES_00`, `URIBL_BLOCKED`) against generated warmup message payloads.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Detailed scoring scorecard showing matched rules, individual point weights, and aggregate spam probability score.
    - Lipgloss token palette: Negative score/Clean (`#10b981`), Spam penalty (`#ef4444`), Neutral (`#94a3b8`).
    - Model view snippet:
      ```text
      ┌─ [OFFLINE SPAM RULES SCORECARD: SCORE -1.8 / 5.0 (SAFE)] ──────────────────┐
      │ RULE NAME              │ POINTS │ DESCRIPTION                              │
      │ DKIM_SIGNED            │ -0.10  │ Message has a valid DKIM signature       │
      │ SPF_PASS               │ -0.01  │ SPF record matches client IP             │
      │ HTML_FONT_LOW_CONTRAST │ +0.50  │ Slight contrast penalty detected in footer│
      │ MIME_HTML_ONLY         │ +1.10  │ ⚠ Missing multipart/alternative plain text│
      │ Net Score: -1.8 (PASS - Message safe for autonomous delivery)             │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Color-coded HTML table with expandable descriptions for all triggered filter rules.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup score-spam`
  - Flags:
    - `--raw-eml <path>`: Path to generated RFC822 .eml file.
    - `--threshold <float>`: Max allowable spam score (default: `2.0`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamRuleHit struct {
        RuleName    string  `json:"ruleName"`
        ScoreDelta  float64 `json:"scoreDelta"`
        Description string  `json:"description"`
        IsPenalty   bool    `json:"isPenalty"`
        HasPassed   bool    `json:"hasPassed"`
    }

    type SpamScoreReport struct {
        MessageID    string        `json:"messageId"`
        AggregateScore float64     `json:"aggregateScore"`
        RuleHits     []SpamRuleHit `json:"ruleHits"`
        IsUnderCap   bool          `json:"isUnderCap"`
        HasFatalRule bool          `json:"hasFatalRule"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamScoreAudit (
        audit_id TEXT PRIMARY KEY,
        message_id TEXT NOT NULL,
        aggregate_score REAL NOT NULL,
        rule_hit_count INTEGER NOT NULL,
        is_under_cap INTEGER NOT NULL DEFAULT 1,
        has_fatal_rule INTEGER NOT NULL DEFAULT 0,
        evaluated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Logic in `pkg/email/security/spam_evaluator.go` preventing dispatch if `aggregateScore > threshold`.

---

### Improvement 758: Disallow List and URL Keyword Sanitizer HUD
- **ID:** `URL-SANITIZER-758`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Scan all embedded URLs and anchor text for suspicious domains, open redirects, shortened links (`bit.ly`, `tinyurl`), and spam-associated keywords.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - URL discovery list with destination hostname classification, redirect status, and security rating badges.
    - Lipgloss token palette: Safe domain (`#10b981`), Blocked/Shortener (`#ef4444`), Warning (`#f59e0b`).
    - Model view snippet:
      ```text
      URL SECURITY SCANNER: 3 Links Found
      [1] https://warmup.gitmap.local/confirm  -> [PASS: Verified Internal Domain]
      [2] https://bit.ly/3xY71K                -> [FAIL: URL Shorteners Disallowed]
      [3] http://insecure-http.com/seed        -> [FAIL: Unencrypted HTTP Scheme]
      Result: 2 DANGEROUS LINKS DETECTED (Dispatch Halted)
      ```
  - **Zero-CDN Web Dashboard:**
    - Interactive link table highlighting destination URL protocols and link reputation pills.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup audit-urls`
  - Flags:
    - `--template <path>`: Source HTML template.
    - `--allow-shorteners`: Allow shortened URLs (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ExtractedURLAudit struct {
        RawURL          string `json:"rawUrl"`
        ParsedHost      string `json:"parsedHost"`
        IsHTTPS         bool   `json:"isHttps"`
        IsShortener     bool   `json:"isShortener"`
        IsDomainAllowed bool   `json:"isDomainAllowed"`
        HasThreatMatch  bool   `json:"hasThreatMatch"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS URLSecurityAudit (
        audit_id TEXT PRIMARY KEY,
        raw_url TEXT NOT NULL,
        parsed_host TEXT NOT NULL,
        is_https INTEGER NOT NULL DEFAULT 1,
        is_shortener INTEGER NOT NULL DEFAULT 0,
        is_domain_allowed INTEGER NOT NULL DEFAULT 1,
        has_threat_match INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `pkg/email/security/url_sanitizer.go` and executed in template validation pipeline.

---

### Improvement 759: Multipart Alternative Plaintext Auto-Synthesizer
- **ID:** `PLAINTEXT-SYNTH-759`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Autonomously generate a clean, perfectly-formatted RFC-compliant plaintext (`text/plain; charset=UTF-8`) alternative from the sanitized HTML DOM, including markdown link conversion and table formatting.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual text preview comparing HTML render to synthesized plaintext formatting with line wrapping indicators.
    - Lipgloss token palette: Plaintext header (`#94a3b8`), Link footnote (`#38bdf8`), Border (`#334155`).
    - Model view snippet:
      ```text
      SYNTHESIZED PLAINTEXT VIEW (RFC822 text/plain):
      ========================================================================
      WARMUP VERIFICATION UPDATE
      ========================================================================
      Your warmup thread is currently active and healthy.

      To inspect your seed placement, visit:
      [1] https://warmup.gitmap.local/confirm
      ------------------------------------------------------------------------
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure `<pre>` block styled with monospace font and accessible color scheme.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup synth-plaintext`
  - Flags:
    - `--html-path <path>`: Source HTML email file.
    - `--output-path <path>`: Destination path for .txt plaintext file.
    - `--wrap-cols <int>`: Column width for word wrapping (default: `72`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PlaintextSynthesisResult struct {
        SourceHTMLPath   string `json:"sourceHtmlPath"`
        WordCount        int    `json:"wordCount"`
        ExtractedLinks   int    `json:"extractedLinks"`
        IsWrappedProperly bool  `json:"isWrappedProperly"`
        HasValidEncoding bool   `json:"hasValidEncoding"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PlaintextSynthesisLog (
        synthesis_id TEXT PRIMARY KEY,
        source_html_path TEXT NOT NULL,
        word_count INTEGER NOT NULL,
        extracted_links INTEGER NOT NULL,
        is_wrapped_properly INTEGER NOT NULL DEFAULT 1,
        has_valid_encoding INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `pkg/email/security/plaintext_synth.go` automatically assembling RFC2046 multipart bodies.

---

### Improvement 760: Standalone Zero-CDN Offline HTML Test Harness & Report Exporter
- **ID:** `SANDBOX-REPORT-760`
- **Cluster:** 76 (Zero-CDN Responsive Email Content Security Sandbox & CSP Auditor)
- **Goal:** Compile all CSP, security, font, contrast, spam rules, and viewport checks into a self-contained, offline HTML report with embedded CSS and zero external dependencies for audit signing.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Overall security badge with one-key report export notification (`Report saved to 02-spec/21-app/reports/csp-audit-report.html`).
    - Lipgloss token palette: Certification Badge (`#10b981`), Summary text (`#f8fafc`), File link (`#38bdf8`).
    - Model view snippet:
      ```text
      ┌─ [CONTENT SECURITY CERTIFICATION: PASSED] ─────────────────────────────────┐
      │ All 6 Security Filters Passed: CSP [✓] | CSS Inlining [✓] | Contrast [✓]    │
      │                                Images [✓] | Fonts [✓] | Spam Score [✓]     │
      │ Export: Standalone HTML report generated with 0 external network calls.    │
      │ Path: 02-spec/21-app/reports/email-security-cert-760.html                  │
      └────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Fully self-contained HTML report with CSS Grid layout, SVG badges, and local system fonts.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup generate-audit-report`
  - Flags:
    - `--template <path>`: Path to HTML template.
    - `--output-file <path>`: Destination path for HTML audit report.
    - `--open`: Open generated report in default system browser.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailSecurityCert struct {
        CertID           string    `json:"certId"`
        TemplateName     string    `json:"templateName"`
        OverallScore     int       `json:"overallScore"`
        GeneratedAt      time.Time `json:"generatedAt"`
        IsCertifiedSafe  bool      `json:"isCertifiedSafe"`
        HasZeroRemoteRef bool      `json:"hasZeroRemoteRef"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailSecurityCertificates (
        cert_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        overall_score INTEGER NOT NULL,
        is_certified_safe INTEGER NOT NULL DEFAULT 1,
        has_zero_remote_ref INTEGER NOT NULL DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/security/cert_generator.go` registering the root command `gitmap email warmup generate-audit-report`.

---

## Verification & Architecture Adherence Checklist

- [x] **Strict Relative Paths**: All spec references, schema destinations, and report paths use relative repository paths (`02-spec/...`, `.ai-memory/...`).
- [x] **Positive Boolean Invariant**: 100% of struct fields and schema columns utilize positive boolean nomenclature (`is...`, `has...`). No negative booleans.
- [x] **Zero CDN Guarantee**: All web visualizers and previews strictly utilize inline SVGs, CSS variables, system font stacks, and zero external script tags or CDNs.
- [x] **Split-DB Conformance**: Schemas partitioned across `gitmap.db` (governance/audit/rules), `installation.db` (local cache/state), and `repodb/pipeline.db` (telemetry/incidents).
- [x] **High Technical Depth**: Complete ANSI terminal UI mockups, Lipgloss token palettes, CLI commands, Go structs, and SQLite DDL provided for every item from 731 through 760 without truncation.
