# Specification 350: Improvements 1331 to 1360 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics, Dynamic Backoff Inversion & Multi-Domain Concurrency Governor, Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)

- **Sequence:** 350
- **Slug:** `350-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/03-improvements-1331-to-1360-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/349-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/348-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (1331 - 1360)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 134** | Items 1331–1340 | Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics |
| **Cluster 135** | Items 1341–1350 | Dynamic Backoff Inversion & Multi-Domain Concurrency Governor |
| **Cluster 136** | Items 1351–1360 | Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD |

---

# Cluster 134: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Items 1331–1340)

---

### Improvement 1331: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Topology Visualizer
- **ID:** `SEED-MATRIX-TOPOLOGY-1331`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide a real-time, interactive Bubble Tea terminal grid visualizing distributed seed mailbox nodes, health metrics, Anycast edge routing latencies, and active deliverability cohorts across Gmail, Microsoft 365, Yahoo, and iCloud seed pools.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [BUBBLE TEA SEED POOL HEALTH MATRIX & TOPOLOGY] ────────────────────────────────────────┐
      │ POOL NODES: 640 | ACTIVE SEEDS: 618 | DEGRADED: 18 | QUARANTINED: 4 | MEAN HEALTH: 98.4%   │
      ├───────────────────────┬────────────┬─────────────┬──────────────┬─────────────┬────────────┤
      │ SEED MAILBOX ID       │ PROVIDER   │ REGION      │ INBOX %      │ RTT (ms)    │ STATUS     │
      ├───────────────────────┼────────────┼─────────────┼──────────────┼─────────────┼────────────┤
      │ seed-gsuite-us-041    │ Google     │ us-east-1   │ 99.8%        │ 18.2ms      │ [●] STABLE │
      │ seed-m365-eu-019      │ Microsoft  │ eu-west-1   │ 97.4%        │ 24.6ms      │ [●] STABLE │
      │ seed-yahoo-us-008     │ Yahoo      │ us-central1 │ 91.2%        │ 42.1ms      │ [▲] WARN   │
      │ seed-icloud-ap-003    │ Apple      │ ap-northeast│ 98.9%        │ 36.4ms      │ [●] STABLE │
      ├───────────────────────┴────────────┴─────────────┴──────────────┴─────────────┴────────────┤
      │ [SELECTED: seed-yahoo-us-008] Health: 91.2% | Target MX: mx-core.mail.yahoo.com             │
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
  - Integrates into `gitmap email warmup seed-matrix` using Bubble Tea viewport and GitMap SQLite telemetry storage at `repodb/pipeline.db` for instant seed degradation root-cause lookup.

---

### Improvement 1332: Autonomous Seed Mailbox SPF, DKIM & DMARC Alignment Diagnostic Inspector
- **ID:** `SEED-AUTH-PROBE-1332`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Autonomously verify strict alignment of SPF (Sender Policy Framework), DKIM (DomainKeys Identified Mail), and DMARC policies for each participating seed mailbox in the warmup swarm to prevent spoofing classification.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED MAILBOX AUTHENTICATION & DNS ALIGNMENT AUDIT] ──────────────────────────────────────┐
      │ AUDITED DOMAINS: 32 | VERIFIED NODES: 640 | SPF PASS: 100% | DKIM 2048: 100% | DMARC: PASS │
      ├──────────────────────────┬───────────┬─────────────┬──────────────┬───────────┬────────────┤
      │ SEED MAILBOX             │ SPF ALIGN │ DKIM SELECT │ KEY LENGTH   │ DMARC POL │ STATUS     │
      ├──────────────────────────┼───────────┼─────────────┼──────────────┼───────────┼────────────┤
      │ warm-alpha@node.gitmap.io│ PASS (v=1)│ gm2026      │ 2048-bit RSA │ p=reject  │ [✓] ALIGNED│
      │ warm-beta@node.gitmap.io │ PASS (v=1)│ gm2026      │ 2048-bit RSA │ p=reject  │ [✓] ALIGNED│
      │ seed-edge@partner.org    │ PASS (v=1)│ k1          │ 1024-bit RSA │ p=none    │ [▲] WEAK   │
      ├──────────────────────────┴───────────┴─────────────┴──────────────┴───────────┴────────────┤
      │ [WARNING: seed-edge@partner.org] 1024-bit RSA deprecated by Gmail/Yahoo 2026 guidelines.   │
      │ Recommendation: Rotate to 2048-bit RSA or Ed25519 selector.                                │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-auth-probe`
  - Flags:
    - `--domain <string>`: Filter diagnostic audit by sender or seed domain.
    - `--is-strict-alignment`: Enforce strict (s=s) rather than relaxed DMARC evaluation (default: `true`).
    - `--has-dnssec-validation`: Require authoritative DNSSEC resolution verification (default: `true`).
    - `--timeout <duration>`: DNS and socket query timeout per seed record (default: `3s`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedAuthAlignmentRecord struct {
        ProbeID                  string         `json:"probe_id"`
        MailboxAddress           string         `json:"mailbox_address"`
        SpfResultStatus          string         `json:"spf_result_status"`
        DkimKeyBitLength         int            `json:"dkim_key_bit_length"`
        DmarcPolicyEnforced      string         `json:"dmarc_policy_enforced"`
        IsSpfAligned             bool           `json:"is_spf_aligned"`
        IsDkimAligned            bool           `json:"is_dkim_aligned"`
        IsDmarcPassing           bool           `json:"is_dmarc_passing"`
        HasDnssecVerified        bool           `json:"has_dnssec_verified"`
    }

    type SeedAuthAlignmentModel struct {
        Records        []SeedAuthAlignmentRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedAuthAlignmentRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedAuthAlignmentRecord (
        probe_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        spf_result_status TEXT NOT NULL,
        dkim_key_bit_length INTEGER NOT NULL DEFAULT 2048,
        dmarc_policy_enforced TEXT NOT NULL,
        is_spf_aligned INTEGER NOT NULL DEFAULT 1,
        is_dkim_aligned INTEGER NOT NULL DEFAULT 1,
        is_dmarc_passing INTEGER NOT NULL DEFAULT 1,
        has_dnssec_verified INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Registers into `gitmap email warmup seed-auth-probe` directly auditing seed sender identities against RFC 7208, RFC 6376, and RFC 7489 rules.

---

### Improvement 1333: Interactive Seed User Engagement Behavior Simulation & Feedback TUI
- **ID:** `SEED-ENGAGE-SIMULATOR-1333`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Simulate and visualize natural human email engagement patterns (mark as important, move from spam to primary inbox, threaded replies, simulated reading dwell time) across the automated seed swarm.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED INTERACTIVE ENGAGEMENT BEHAVIOR MONITOR] ───────────────────────────────────────────┐
      │ SIMULATED AGENTS: 48 | DWELL MEDIAN: 26.4s | RESCUE ACTIONS: 14 | REPLIES DISPATCHED: 84   │
      ├────────────────────────────┬──────────────────┬──────────────┬─────────────┬───────────────┤
      │ SEED MAILBOX               │ THREAD SUBJECT   │ DWELL TIME   │ RESCUE OP   │ ACTION DISP   │
      ├────────────────────────────┼──────────────────┼──────────────┼─────────────┼───────────────┤
      │ seed-09@warmup.gitmap.dev  │ Q3 Architecture  │ 34.2s        │ N/A (Inbox) │ REPLIED       │
      │ seed-14@warmup.gitmap.dev  │ Sync Pipeline    │ 18.9s        │ RESCUED (SP)│ STARRED & MOVE│
      │ seed-22@warmup.gitmap.dev  │ Weekly Digest    │ 41.5s        │ N/A (Inbox) │ MARKED IMP    │
      ├────────────────────────────┴──────────────────┴──────────────┴─────────────┴───────────────┤
      │ [RESCUE HOOK: seed-14] Removed from spam folder, marked not spam, added sender to contacts. │
      │ Keys: [p: Pause Simulation] [m: Force Rescue All] [s: Adjust Speeds] [q: Back]             │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-engagement`
  - Flags:
    - `--cohort-id <string>`: Specify seed mailbox cohort group identifier.
    - `--dwell-time-min <duration>`: Minimum simulated dwell reading time (default: `15s`).
    - `--reply-rate-target <float>`: Target auto-reply percentage within cohort (default: `0.35`).
    - `--is-rescued-from-spam`: Enable rescue protocol moving detected spam warmup mail to inbox (default: `true`).
    - `--has-star-engagement`: Enable random starring and marking as important (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedEngagementEventRecord struct {
        EventID                  string         `json:"event_id"`
        MailboxAddress           string         `json:"mailbox_address"`
        ThreadID                 string         `json:"thread_id"`
        ActionType               string         `json:"action_type"`
        DwellDurationMs          int64          `json:"dwell_duration_ms"`
        IsInboxRescued           bool           `json:"is_inbox_rescued"`
        IsThreadReplied          bool           `json:"is_thread_replied"`
        HasMarkedImportant       bool           `json:"has_marked_important"`
        HasUserInteraction       bool           `json:"has_user_interaction"`
    }

    type SeedEngagementModel struct {
        Records        []SeedEngagementEventRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedEngagementEventRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedEngagementEventRecord (
        event_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        thread_id TEXT NOT NULL,
        action_type TEXT NOT NULL,
        dwell_duration_ms INTEGER NOT NULL DEFAULT 0,
        is_inbox_rescued INTEGER NOT NULL DEFAULT 0,
        is_thread_replied INTEGER NOT NULL DEFAULT 0,
        has_marked_important INTEGER NOT NULL DEFAULT 1,
        has_user_interaction INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Hooks into GitMap email warm-up scheduler to continuously stimulate receiving ESP algorithms into recognizing incoming messages as high-priority user mail.

---

### Improvement 1334: Multi-Provider Seed Reputation Radar & IP Blacklist Real-Time Monitor
- **ID:** `SEED-REPUTATION-RADAR-1334`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Provide an aggregated radar dashboard actively cross-referencing seed sending/receiving IPs and domains against 50+ Real-Time Blackhole Lists (Spamhaus ZEN, Barracuda, SORBS, Invaluement) and Postmaster tools.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED REPUTATION RADAR & RBL AGGREGATED STATUS] ──────────────────────────────────────────┐
      │ MONITORED IPS: 16 | DOMAINS: 8 | ACTIVE RBLS CHECKED: 52 | CLEAN RATE: 100.0%             │
      ├────────────────────┬────────────────────────┬──────────────┬───────────────┬───────────────┤
      │ HOST / IP          │ RBL FEED               │ RETURN CODE  │ REPUTATION    │ HEALTH STATUS │
      ├────────────────────┼────────────────────────┼──────────────┼───────────────┼───────────────┤
      │ 198.51.100.44      │ zen.spamhaus.org       │ NXDOMAIN     │ 100/100 (HIGH)│ [●] CLEAN     │
      │ 198.51.100.44      │ b.barracudacentral.org │ NXDOMAIN     │ 100/100 (HIGH)│ [●] CLEAN     │
      │ warmup.gitmap.dev  │ multi.surbl.org        │ NXDOMAIN     │ 99/100 (HIGH) │ [●] CLEAN     │
      ├────────────────────┴────────────────────────┴──────────────┴───────────────┴───────────────┤
      │ Google Postmaster Tier: HIGH (0.01% user spam complaint rate across rolling 7 days).       │
      │ Microsoft SNDS Status: GREEN (Zero trap hits detected across Outlook/Hotmail corpus).       │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup blacklist-radar`
  - Flags:
    - `--ip <string>`: Specific IP address or CIDR range to query across RBLs.
    - `--domain <string>`: Domain or subdomain to cross-reference against URIBL/SURBL.
    - `--is-live-watch`: Continuously monitor blacklist state and notify on state changes (default: `true`).
    - `--has-postmaster-api`: Fetch Google Postmaster / Microsoft SNDS API reputation metrics (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedRblReputationRecord struct {
        ListingID                string         `json:"listing_id"`
        TargetHost               string         `json:"target_host"`
        RblProvider              string         `json:"rbl_provider"`
        ResponseCode             string         `json:"response_code"`
        ReputationScore          float64        `json:"reputation_score"`
        IsListed                 bool           `json:"is_listed"`
        IsClean                  bool           `json:"is_clean"`
        HasDelistingUrl          bool           `json:"has_delisting_url"`
        HasTelemetryAlert        bool           `json:"has_telemetry_alert"`
    }

    type SeedRblReputationModel struct {
        Records        []SeedRblReputationRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedRblReputationRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedRblReputationRecord (
        listing_id TEXT PRIMARY KEY,
        target_host TEXT NOT NULL,
        rbl_provider TEXT NOT NULL,
        response_code TEXT NOT NULL,
        reputation_score REAL NOT NULL DEFAULT 100.0,
        is_listed INTEGER NOT NULL DEFAULT 0,
        is_clean INTEGER NOT NULL DEFAULT 1,
        has_delisting_url INTEGER NOT NULL DEFAULT 0,
        has_telemetry_alert INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates into `gitmap email warmup blacklist-radar` feeding into GitMap quarantine pipeline to automatically suspend compromised IP addresses.

---

### Improvement 1335: Interactive Visual Heatmap of Per-Folder Seed Mailbox Placement
- **ID:** `SEED-INBOX-PLACEMENT-MAP-1335`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Render an interactive ANSI/ASCII heat map displaying exact mail folder destinations (Primary, Promotions, Social, Updates, Spam, Trash) across ESP targets during warmup iterations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED FOLDER PLACEMENT HEATMAP DISTRIBUTION] ─────────────────────────────────────────────┐
      │ SAMPLED DELIVERIES: 1,420 | PRIMARY: 93.8% | PROMOTIONS: 4.6% | SPAM: 1.4% | TRASH: 0.2%   │
      ├────────────────────┬───────────┬──────────────┬─────────────┬────────────┬─────────────────┤
      │ ESP PROVIDER       │ PRIMARY   │ PROMOTIONS   │ SPAM        │ TRASH      │ HEAT GRADIENT   │
      ├────────────────────┼───────────┼──────────────┼─────────────┼────────────┼─────────────────┤
      │ Gmail / Workspace  │ 96.2%     │  3.1%        │  0.7%       │ 0.0%       │ [██████████░░░] │
      │ Microsoft 365 / OL │ 94.1%     │  4.2%        │  1.5%       │ 0.2%       │ [█████████░░░░] │
      │ Yahoo / AOL Mail   │ 89.0%     │  8.4%        │  2.4%       │ 0.2%       │ [████████░░░░░] │
      │ Apple iCloud Mail  │ 98.4%     │  1.2%        │  0.4%       │ 0.0%       │ [██████████░░░] │
      ├────────────────────┴───────────┴──────────────┴─────────────┴────────────┴─────────────────┤
      │ Yahoo Promotions threshold active. Adjusting payload tone to reduce promotional markers.   │
      │ Controls: [←/→: Shift Time Window] [v: Switch View Provider/Template] [q: Exit]            │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup placement-heatmap`
  - Flags:
    - `--window <duration>`: Time range for placement aggregation (default: `24h`).
    - `--group-by <string>`: Aggregation dimension (`provider`, `subject-template`, `ip`; default: `provider`).
    - `--is-color-coded`: Enable rich 256-color heat gradient rendering (default: `true`).
    - `--has-tabulated-export`: Export heat tabular results to structured CSV/JSON (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedFolderPlacementRecord struct {
        PlacementID              string         `json:"placement_id"`
        ProviderName             string         `json:"provider_name"`
        FolderPrimaryCount       int            `json:"folder_primary_count"`
        FolderPromoCount         int            `json:"folder_promo_count"`
        FolderSpamCount          int            `json:"folder_spam_count"`
        FolderTrashCount         int            `json:"folder_trash_count"`
        PrimaryRatio             float64        `json:"primary_ratio"`
        IsPrimaryDominant        bool           `json:"is_primary_dominant"`
        HasSpamLeakage           bool           `json:"has_spam_leakage"`
        HasPromoDiversion        bool           `json:"has_promo_diversion"`
    }

    type SeedFolderPlacementModel struct {
        Records        []SeedFolderPlacementRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedFolderPlacementRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedFolderPlacementRecord (
        placement_id TEXT PRIMARY KEY,
        provider_name TEXT NOT NULL,
        folder_primary_count INTEGER NOT NULL DEFAULT 0,
        folder_promo_count INTEGER NOT NULL DEFAULT 0,
        folder_spam_count INTEGER NOT NULL DEFAULT 0,
        folder_trash_count INTEGER NOT NULL DEFAULT 0,
        primary_ratio REAL NOT NULL DEFAULT 1.0,
        is_primary_dominant INTEGER NOT NULL DEFAULT 1,
        has_spam_leakage INTEGER NOT NULL DEFAULT 0,
        has_promo_diversion INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Visualizes deliverability trends directly in GitMap CLI without external tracking webhooks or third-party SaaS platforms.

---

### Improvement 1336: Dynamic Seed Pool Auto-Rotation & Mailbox Exhaustion Governor
- **ID:** `SEED-POOL-ROTATION-GOV-1336`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Enforce automatic rotation and rest cycles for seed mailboxes exceeding maximum daily interaction quotas, preserving account longevity and preventing algorithmic burn.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED POOL EXHAUSTION & AUTO-ROTATION GOVERNOR] ───────────────────────────────────────────┐
      │ TOTAL SEEDS: 640 | ACTIVE: 580 | RESTING (COOLDOWN): 48 | EMERGENCY RESERVE: 12            │
      ├──────────────────────────┬──────────────┬──────────────┬───────────────┬───────────────────┤
      │ SEED MAILBOX             │ DAILY USAGE  │ QUOTA REMAIN │ STATUS        │ COOLDOWN REMAIN   │
      ├──────────────────────────┼──────────────┼──────────────┼───────────────┼───────────────────┤
      │ seed-01@warmup.gitmap.dev│ 30 / 30      │  0 msgs      │ [▲] RESTING   │ 07h 42m remaining │
      │ seed-02@warmup.gitmap.dev│ 28 / 30      │  2 msgs      │ [●] WARNING   │ Active (Throttled)│
      │ seed-03@warmup.gitmap.dev│ 14 / 30      │ 16 msgs      │ [●] ACTIVE    │ Active (Optimal)  │
      │ seed-res-01@gitmap.dev   │  0 / 30      │ 30 msgs      │ [★] STANDBY   │ Ready for swap    │
      ├──────────────────────────┴──────────────┴──────────────┴───────────────┴───────────────────┤
      │ Swapped seed-01 with seed-res-01. Warmup pipeline flow preserved without interruption.     │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-rotation`
  - Flags:
    - `--max-daily-inbound <int>`: Maximum received warmup emails per seed per 24 hours (default: `30`).
    - `--rest-period <duration>`: Mandatory cooldown duration when quota reached (default: `12h`).
    - `--is-auto-rebalance`: Automatically spin up cold reserve seeds upon rotation (default: `true`).
    - `--has-circuit-breaker`: Halt entire cohort if consecutive authentication failures exceed 3 (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedRotationQuotaRecord struct {
        RotationID               string         `json:"rotation_id"`
        MailboxAddress           string         `json:"mailbox_address"`
        InboundCountToday        int            `json:"inbound_count_today"`
        DailyMaxQuota            int            `json:"daily_max_quota"`
        CooldownExpiresAt        string         `json:"cooldown_expires_at"`
        IsActive                 bool           `json:"is_active"`
        IsCoolingDown            bool           `json:"is_cooling_down"`
        HasQuotaHeadroom         bool           `json:"has_quota_headroom"`
        HasReserveAssigned       bool           `json:"has_reserve_assigned"`
    }

    type SeedRotationQuotaModel struct {
        Records        []SeedRotationQuotaRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedRotationQuotaRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedRotationQuotaRecord (
        rotation_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        inbound_count_today INTEGER NOT NULL DEFAULT 0,
        daily_max_quota INTEGER NOT NULL DEFAULT 30,
        cooldown_expires_at TEXT NOT NULL,
        is_active INTEGER NOT NULL DEFAULT 1,
        is_cooling_down INTEGER NOT NULL DEFAULT 0,
        has_quota_headroom INTEGER NOT NULL DEFAULT 1,
        has_reserve_assigned INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Governs the active seed lifecycle inside GitMap warmup runner, reading quotas directly from `gitmap.db`.

---

### Improvement 1337: Warmup Email Conversational Thread Integrity & Header Stitching Inspector
- **ID:** `SEED-THREAD-INTEGRITY-1337`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Audit RFC 5322 In-Reply-To and References message header chaining across simulated multi-turn warmup conversations to guarantee genuine conversational topology to ESP anti-abuse parsers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [CONVERSATIONAL THREAD INTEGRITY & HEADER TREE] ──────────────────────────────────────────┐
      │ THREAD: <msg-root-91823@gitmap.dev> | SUBJECT: Re: Architecture Blueprint v4               │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ ├─ [Turn 1] Originator: author@gitmap.dev (Message-ID: <msg-root-91823@gitmap.dev>)        │
      │ │   └─ SPF: Pass | DKIM: Valid | Date: 2026-10-27 10:14:02 UTC                             │
      │ ├─ [Turn 2] Responder: seed-04@warmup.gitmap.dev                                           │
      │ │   ├─ In-Reply-To: <msg-root-91823@gitmap.dev> [VERIFIED]                                │
      │ │   ├─ References: <msg-root-91823@gitmap.dev> [CHAIN INTACT]                             │
      │ │   └─ Dwell Time: 42s | Simulated Typo Fixes: Yes                                        │
      │ └─ [Turn 3] Originator Reply: author@gitmap.dev                                            │
      │     ├─ In-Reply-To: <msg-reply-seed04-129@warmup.gitmap.dev> [VERIFIED]                   │
      │     └─ References: <msg-root-91823@gitmap.dev> <msg-reply-seed04-129@...> [VALID CHAIN]  │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Audit Status: [✓] PERFECT CONVERSATIONAL CONTINUITY (ESP Trust Factor: 99.4%)             │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup thread-audit`
  - Flags:
    - `--thread-depth-min <int>`: Minimum expected conversation replies in thread (default: `2`).
    - `--is-validate-headers`: Inspect In-Reply-To and References RFC syntax and checksums (default: `true`).
    - `--has-subject-re-prefix`: Ensure strict localization of Re: / Fwd: prefixes (default: `true`).
    - `--output-format <string>`: Output terminal render or JSON tree (`terminal`, `tree`, `json`; default: `tree`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedThreadIntegrityRecord struct {
        ThreadAuditID            string         `json:"thread_audit_id"`
        RootMessageID            string         `json:"root_message_id"`
        SubjectLine              string         `json:"subject_line"`
        ReplyCount               int            `json:"reply_count"`
        HeaderChainStatus        string         `json:"header_chain_status"`
        IsChainValid             bool           `json:"is_chain_valid"`
        IsReferencesComplete     bool           `json:"is_references_complete"`
        HasMatchingSubject       bool           `json:"has_matching_subject"`
        HasHeaderAnomaly         bool           `json:"has_header_anomaly"`
    }

    type SeedThreadIntegrityModel struct {
        Records        []SeedThreadIntegrityRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedThreadIntegrityRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedThreadIntegrityRecord (
        thread_audit_id TEXT PRIMARY KEY,
        root_message_id TEXT NOT NULL,
        subject_line TEXT NOT NULL,
        reply_count INTEGER NOT NULL DEFAULT 0,
        header_chain_status TEXT NOT NULL,
        is_chain_valid INTEGER NOT NULL DEFAULT 1,
        is_references_complete INTEGER NOT NULL DEFAULT 1,
        has_matching_subject INTEGER NOT NULL DEFAULT 1,
        has_header_anomaly INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded in GitMap mail generation engine, validating outgoing replies prior to transmission over SMTP.

---

### Improvement 1338: Global Edge Seed POP Latency Topology & Anycast Route Degradation Monitor
- **ID:** `SEED-GEO-LATENCY-DIAG-1338`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Trace and isolate edge network latency spikes, AS peering bottlenecks, and BGP route flapping between warmup sender nodes and destination MX clusters.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [GLOBAL EDGE SEED POP LATENCY & BGP ROUTE TRACER] ────────────────────────────────────────┐
      │ SOURCE: us-east-1 (Ashburn) -> DESTINATION: mx1.mail.icloud.com (Anycast AS6185)           │
      ├─────┬──────────────────────────┬──────────────┬────────────┬───────────────────────────────┤
      │ HOP │ NODE / ROUTER IP         │ ASN & ENTITY │ RTT (ms)   │ STATUS                        │
      ├─────┼──────────────────────────┼──────────────┼────────────┼───────────────────────────────┤
      │ 1   │ 10.0.0.1                 │ Local GW     │ 0.32ms     │ [●] Clean                     │
      │ 2   │ 198.51.100.254           │ AS16509 AWS  │ 1.12ms     │ [●] Clean                     │
      │ 3   │ 64.233.174.45            │ AS15169 Equi │ 4.88ms     │ [●] Direct Peering            │
      │ 4   │ 17.58.23.1               │ AS6185 Apple │ 11.45ms    │ [●] Optimal Edge POP          │
      ├─────┴──────────────────────────┴──────────────┴────────────┴───────────────────────────────┤
      │ Route Summary: 4 Hops | Total Round-Trip Time: 11.45ms | Jitter: 0.2ms | Loss: 0.0%        │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup geo-latency`
  - Flags:
    - `--target-mx <string>`: Destination MX host to ping and route-trace.
    - `--max-rtt-alert <duration>`: Threshold to flag network degradation (default: `60ms`).
    - `--is-traceroute-enabled`: Execute ICMP/TCP traceroute on latency breach (default: `true`).
    - `--has-peering-lookup`: Query ASN peering info from PeeringDB (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedGeoLatencyRecord struct {
        LatencyHopID             string         `json:"latency_hop_id"`
        SourceRegion             string         `json:"source_region"`
        TargetMxHost             string         `json:"target_mx_host"`
        MeanRttMs                float64        `json:"mean_rtt_ms"`
        JitterMs                 float64        `json:"jitter_ms"`
        HopCount                 int            `json:"hop_count"`
        IsDegraded               bool           `json:"is_degraded"`
        IsOptimalRoute           bool           `json:"is_optimal_route"`
        HasFlappingRoute         bool           `json:"has_flapping_route"`
        HasDirectPeer            bool           `json:"has_direct_peer"`
    }

    type SeedGeoLatencyModel struct {
        Records        []SeedGeoLatencyRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedGeoLatencyRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedGeoLatencyRecord (
        latency_hop_id TEXT PRIMARY KEY,
        source_region TEXT NOT NULL,
        target_mx_host TEXT NOT NULL,
        mean_rtt_ms REAL NOT NULL DEFAULT 0.0,
        jitter_ms REAL NOT NULL DEFAULT 0.0,
        hop_count INTEGER NOT NULL DEFAULT 0,
        is_degraded INTEGER NOT NULL DEFAULT 0,
        is_optimal_route INTEGER NOT NULL DEFAULT 1,
        has_flapping_route INTEGER NOT NULL DEFAULT 0,
        has_direct_peer INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Diagnostics tool exposed via `gitmap email warmup geo-latency` to verify zero carrier-level drops before sending high-volume batches.

---

### Improvement 1339: Autonomous Asynchronous Bounce & DSN Semantic Classifier TUI
- **ID:** `SEED-BOUNCE-CLASSIFIER-1339`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Parse asynchronous Delivery Status Notifications (DSNs) and SMTP bounce bodies using regex and AST classifiers to distinguish hard bounces (5.1.1 unknown recipient) from soft throttle pauses (4.2.1 rate limit).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [BOUNCE & DSN ASYNCHRONOUS SEMANTIC CLASSIFIER] ──────────────────────────────────────────┐
      │ ANALYZED DSNS: 18 | HARD BOUNCES: 1 | SOFT THROTTLES: 17 | AUTO-SUPPRESSED: 1               │
      ├──────────────────────────┬───────────┬──────────────┬───────────────┬──────────────────────┤
      │ RECIPIENT                │ CODE      │ CLASSIF      │ ROOT CAUSE    │ ACTION TAKEN         │
      ├──────────────────────────┼───────────┼──────────────┼───────────────┼──────────────────────┤
      │ typo-user@domain.invalid │ 5.1.1     │ HARD BOUNCE  │ Bad Mailbox   │ Added to Suppression │
      │ test-inbox@fastmail.com  │ 4.2.1     │ SOFT RATE    │ Conn Throttled│ Backoff 300s (Queue) │
      │ test-seed@free.fr        │ 4.7.1     │ GREYLIST     │ Greylisting   │ Retry in 900s        │
      ├──────────────────────────┴───────────┴──────────────┴───────────────┴──────────────────────┤
      │ [SELECTION: typo-user@domain.invalid] Diagnostic: 550 5.1.1 Recipient address rejected      │
      │ Actions: [s: View Full SMTP DSN] [u: Undo Suppression] [c: Clear Resolved] [q: Quit]      │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup bounce-audit`
  - Flags:
    - `--since <duration>`: Lookback time for bounce events (default: `48h`).
    - `--filter-code <string>`: Filter by enhanced status code (e.g. `5.1.1`, `4.7.1`).
    - `--is-auto-suppress`: Add hard-bounced addresses immediately to suppression table (default: `true`).
    - `--has-raw-dsn-dump`: Display full raw MIME multipart DSN payload in drill-down (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedBounceAuditRecord struct {
        BounceID                 string         `json:"bounce_id"`
        RecipientAddress         string         `json:"recipient_address"`
        EnhancedStatusCode       string         `json:"enhanced_status_code"`
        BounceCategory           string         `json:"bounce_category"`
        SmtpDiagnosticText       string         `json:"smtp_diagnostic_text"`
        IsHardBounce             bool           `json:"is_hard_bounce"`
        IsSoftThrottle           bool           `json:"is_soft_throttle"`
        HasSuppressionApplied    bool           `json:"has_suppression_applied"`
        HasRetryScheduled        bool           `json:"has_retry_scheduled"`
    }

    type SeedBounceAuditModel struct {
        Records        []SeedBounceAuditRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedBounceAuditRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedBounceAuditRecord (
        bounce_id TEXT PRIMARY KEY,
        recipient_address TEXT NOT NULL,
        enhanced_status_code TEXT NOT NULL,
        bounce_category TEXT NOT NULL,
        smtp_diagnostic_text TEXT NOT NULL,
        is_hard_bounce INTEGER NOT NULL DEFAULT 0,
        is_soft_throttle INTEGER NOT NULL DEFAULT 0,
        has_suppression_applied INTEGER NOT NULL DEFAULT 0,
        has_retry_scheduled INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Writes hard bounces directly to `repodb/pipeline.db` suppression tables, protecting warm domain reputations.

---

### Improvement 1340: Seed Swarm Behavioral Anomaly & Sybil Attack Intrusion Guard
- **ID:** `SEED-ANOMALY-ALERT-1340`
- **Cluster:** 134 (Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics)
- **Goal:** Detect compromised seed mailboxes displaying unauthorized outbound relaying, password changes, or abnormal webhook traffic to protect swarm hygiene from malicious hijacking.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SEED SWARM BEHAVIORAL ANOMALY & INTRUSION GUARD] ────────────────────────────────────────┐
      │ MONITORED ACCOUNTS: 640 | INTRUSION EVENTS: 0 | QUARANTINED: 0 | INTEGRITY SCORE: 100%      │
      ├──────────────────────────┬────────────────┬───────────────┬───────────────┬────────────────┤
      │ SEED MAILBOX             │ ANOMALY TYPE   │ CONFIDENCE    │ STATE         │ MITIGATION     │
      ├──────────────────────────┼────────────────┼───────────────┼───────────────┼────────────────┤
      │ seed-099@warmup.gitmap   │ NORMAL PATTERN │ 0.00 (CLEAN)  │ [●] PROTECTED │ None Required  │
      │ seed-102@warmup.gitmap   │ HIGH VOLUME TX │ 0.88 (ALERT)  │ [▲] ISOLATED  │ Token Revoked  │
      ├──────────────────────────┴────────────────┴───────────────┴───────────────┴────────────────┤
      │ [SECURITY INCIDENT ISOLATED] seed-102 attempted 150 SMTP auths from unapproved IP ASN.    │
      │ Token severed automatically. Zero impact on wider warmup fleet.                            │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup seed-guard`
  - Flags:
    - `--threat-level <string>`: Alert sensitivity threshold (`low`, `medium`, `high`, `paranoid`; default: `medium`).
    - `--is-auto-quarantine`: Instantly isolate compromised seed upon threat trigger (default: `true`).
    - `--has-slack-webhook`: Forward security incidents to configured alert webhook (default: `false`).
    - `--evict-token`: Invalidate and rotate IMAP/OAuth2 session tokens immediately.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SeedSecurityGuardRecord struct {
        IncidentID               string         `json:"incident_id"`
        MailboxAddress           string         `json:"mailbox_address"`
        ThreatVector             string         `json:"threat_vector"`
        AnomalyScore             float64        `json:"anomaly_score"`
        MitigationAction         string         `json:"mitigation_action"`
        IsCompromised            bool           `json:"is_compromised"`
        IsQuarantined            bool           `json:"is_quarantined"`
        HasTokenRevoked          bool           `json:"has_token_revoked"`
        HasAuditAlertSent        bool           `json:"has_audit_alert_sent"`
    }

    type SeedSecurityGuardModel struct {
        Records        []SeedSecurityGuardRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SeedSecurityGuardRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedSecurityGuardRecord (
        incident_id TEXT PRIMARY KEY,
        mailbox_address TEXT NOT NULL,
        threat_vector TEXT NOT NULL,
        anomaly_score REAL NOT NULL DEFAULT 0.0,
        mitigation_action TEXT NOT NULL,
        is_compromised INTEGER NOT NULL DEFAULT 0,
        is_quarantined INTEGER NOT NULL DEFAULT 0,
        has_token_revoked INTEGER NOT NULL DEFAULT 0,
        has_audit_alert_sent INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded security agent in GitMap daemon guaranteeing zero unauthorized relay abuse of warmup seed accounts.

---

# Cluster 135: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (Items 1341–1350)

---

### Improvement 1341: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor
- **ID:** `BACKOFF-INVERSION-ORCH-1341`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Dynamically govern outbound SMTP concurrency per destination domain, inverting connection pools upon receiving 4xx throttling replies and dynamically adjusting sliding windows.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [MULTI-DOMAIN CONCURRENCY GOVERNOR & BACKOFF INVERSION] ──────────────────────────────────┐
      │ ACTIVE DOMAINS: 142 | GLOBAL CONCURRENCY: 380 | ADAPTIVE BACKOFF EVENTS: 3                 │
      ├────────────────────┬─────────────┬─────────────┬──────────────┬──────────────┬─────────────┤
      │ DESTINATION DOMAIN │ CONCURRENCY │ MAX LIMIT   │ 4XX COUNTER  │ BACKOFF COEF │ STATE       │
      ├────────────────────┼─────────────┼─────────────┼──────────────┼──────────────┼─────────────┤
      │ gmail.com          │ 24 workers  │ 40 workers  │ 0 errors     │ 1.00x (Norm) │ [●] RAMPING │
      │ outlook.com        │ 12 workers  │ 20 workers  │ 1 error      │ 1.50x (Slow) │ [▲] BACKOFF │
      │ yahoo.com          │  8 workers  │ 15 workers  │ 0 errors     │ 1.00x (Norm) │ [●] OPTIMAL │
      │ orange.fr          │  2 workers  │  5 workers  │ 2 errors     │ 2.25x (Hold) │ [■] HALTED  │
      ├────────────────────┴─────────────┴─────────────┴──────────────┴──────────────┴─────────────┤
      │ [INVERSION TRIGGER: orange.fr] Received 4.2.1 Service unavailable. Dropping concurrency.    │
      │ Concurrency halved: 4 -> 2 workers. Inter-message delay shifted: 250ms -> 1200ms.          │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup concurrency-governor`
  - Flags:
    - `--domain <string>`: Target destination domain to tune (default: `all`).
    - `--initial-concurrency <int>`: Base parallel worker count per domain (default: `5`).
    - `--max-concurrency <int>`: Ceiling limit for aggressive warmup phases (default: `50`).
    - `--is-backoff-inversion`: Enable inverse rate-limiting upon 4xx SMTP responses (default: `true`).
    - `--has-jitter-smoothing`: Apply randomized exponential jitter to transmission timers (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DomainConcurrencyGovernorRecord struct {
        DomainID                 string         `json:"domain_id"`
        TargetDomain             string         `json:"target_domain"`
        CurrentConcurrency       int            `json:"current_concurrency"`
        MaxSafeConcurrency       int            `json:"max_safe_concurrency"`
        BackoffMultiplier        float64        `json:"backoff_multiplier"`
        ConsecutiveThrottles     int            `json:"consecutive_throttles"`
        IsThrottled              bool           `json:"is_throttled"`
        IsRampingUp              bool           `json:"is_ramping_up"`
        HasAdaptiveJitter        bool           `json:"has_adaptive_jitter"`
        HasSmtpLock              bool           `json:"has_smtp_lock"`
    }

    type DomainConcurrencyGovernorModel struct {
        Records        []DomainConcurrencyGovernorRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord DomainConcurrencyGovernorRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DomainConcurrencyGovernorRecord (
        domain_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        current_concurrency INTEGER NOT NULL DEFAULT 5,
        max_safe_concurrency INTEGER NOT NULL DEFAULT 20,
        backoff_multiplier REAL NOT NULL DEFAULT 1.0,
        consecutive_throttles INTEGER NOT NULL DEFAULT 0,
        is_throttled INTEGER NOT NULL DEFAULT 0,
        is_ramping_up INTEGER NOT NULL DEFAULT 1,
        has_adaptive_jitter INTEGER NOT NULL DEFAULT 1,
        has_smtp_lock INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Core execution governor in GitMap SMTP dispatcher preventing domain-level blocking via dynamic backoff.

---

### Improvement 1342: Hierarchical Token-Bucket Rate Limiter with Sliding Window Quantization
- **ID:** `BACKOFF-TOKEN-BUCKET-1342`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Enforce microsecond-accurate token-bucket rate limiting across per-domain, per-IP, and per-account boundaries with sliding window interval quantization.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [HIERARCHICAL TOKEN-BUCKET SLIDING WINDOW MONITOR] ───────────────────────────────────────┐
      │ ACTIVE BUCKETS: 512 | BURST CONSUMPTION: 84 tok/s | REFILL VELOCITY: 120 tok/s             │
      ├──────────────────────────┬──────────────┬─────────────┬──────────────┬─────────────────────┤
      │ BUCKET IDENTIFIER        │ TOKENS / CAP │ FILL RATE/s │ DRAIN DELAY  │ PERMIT STATUS       │
      ├──────────────────────────┼──────────────┼─────────────┼──────────────┼─────────────────────┤
      │ domain:gmail.com         │ 18.4 / 20    │ 15.0 /s     │ 0.0ms        │ [●] GRANTED         │
      │ ip:198.51.100.44         │  4.2 / 20    │  8.0 /s     │ 12.5ms       │ [●] FLOW RESTRICTED │
      │ account:seed-01@gitmap   │  0.0 / 5     │  1.0 /s     │ 480.0ms      │ [▲] BUCKET EXHAUSTED│
      ├──────────────────────────┴──────────────┴─────────────┴──────────────┴─────────────────────┤
      │ Token refill calculation using atomic CAS nano clocks. Zero OS thread locks.               │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup token-bucket`
  - Flags:
    - `--fill-rate <float>`: Tokens added per second per target bucket (default: `10.0`).
    - `--capacity <int>`: Maximum burst token capacity in bucket (default: `20`).
    - `--is-hierarchical`: Evaluate parent domain and IP bucket limits concurrently (default: `true`).
    - `--has-zero-allocation`: Use zero-alloc atomic bitfield state tracking in Go (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TokenBucketStateRecord struct {
        BucketKey                string         `json:"bucket_key"`
        RemainingTokens          float64        `json:"remaining_tokens"`
        RefillRatePerSec         float64        `json:"refill_rate_per_sec"`
        Capacity                 int            `json:"capacity"`
        LastRefillNano           int64          `json:"last_refill_nano"`
        IsAllowed                bool           `json:"is_allowed"`
        IsExhausted              bool           `json:"is_exhausted"`
        HasBurstAllowance        bool           `json:"has_burst_allowance"`
        HasParentThrottle        bool           `json:"has_parent_throttle"`
    }

    type TokenBucketStateModel struct {
        Records        []TokenBucketStateRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord TokenBucketStateRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TokenBucketStateRecord (
        bucket_key TEXT PRIMARY KEY,
        remaining_tokens REAL NOT NULL DEFAULT 20.0,
        refill_rate_per_sec REAL NOT NULL DEFAULT 10.0,
        capacity INTEGER NOT NULL DEFAULT 20,
        last_refill_nano INTEGER NOT NULL DEFAULT 0,
        is_allowed INTEGER NOT NULL DEFAULT 1,
        is_exhausted INTEGER NOT NULL DEFAULT 0,
        has_burst_allowance INTEGER NOT NULL DEFAULT 1,
        has_parent_throttle INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Zero-lock atomic token-bucket engine in `gitmap` managing exact transmission velocity per ESP.

---

### Improvement 1343: Three-State Circuit Breaker (Closed, Open, Half-Open) per MX Route
- **ID:** `BACKOFF-CIRCUIT-BREAKER-1343`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Provide a hardened three-state circuit breaker per destination MX server to isolate dead, greylisting, or failing mail hosts without stalling the entire sending queue.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [MX ROUTE THREE-STATE CIRCUIT BREAKER TELEMETRY] ─────────────────────────────────────────┐
      │ CIRCUITS TRACKED: 84 | CLOSED (HEALTHY): 82 | OPEN (TRIPPED): 1 | HALF-OPEN (PROBING): 1   │
      ├──────────────────────────┬──────────────┬──────────────┬──────────────┬────────────────────┤
      │ TARGET MX ROUTE          │ STATE        │ FAIL COUNTER │ PROBE RATIO  │ COOLDOWN REMAIN    │
      ├──────────────────────────┼──────────────┼──────────────┼──────────────┼────────────────────┤
      │ alt1.aspmx.l.google.com  │ CLOSED       │ 0 / 5 fails  │ N/A (Online) │ 0s (Streaming)     │
      │ mail.protection.outlook  │ HALF-OPEN    │ 1 fail       │ 2 / 3 probes │ 12s remaining      │
      │ mx.rogers.com            │ OPEN         │ 5 / 5 fails  │ Tripped      │ 44s cooldown       │
      ├──────────────────────────┴──────────────┴──────────────┴──────────────┴────────────────────┤
      │ [ALERT: mx.rogers.com] Circuit OPEN: 554 Relay Access Denied x5. Traffic re-queued.        │
      │ Probing scheduled in 44s with single probe handshake.                                      │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup circuit-breaker`
  - Flags:
    - `--failure-threshold <int>`: Consecutive failures to trip circuit from Closed to Open (default: `5`).
    - `--open-timeout <duration>`: Sleep duration before transitioning Open to Half-Open probe (default: `60s`).
    - `--half-open-successes <int>`: Successful probes required to restore Closed state (default: `3`).
    - `--is-auto-tripped`: Trip immediately upon 554 Transaction Failed (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CircuitBreakerStateRecord struct {
        MxHostKey                string         `json:"mx_host_key"`
        StateName                string         `json:"state_name"`
        FailureCount             int            `json:"failure_count"`
        SuccessProbeCount        int            `json:"success_probe_count"`
        NextProbeAllowedAt       string         `json:"next_probe_allowed_at"`
        IsClosed                 bool           `json:"is_closed"`
        IsOpen                   bool           `json:"is_open"`
        IsHalfOpen               bool           `json:"is_half_open"`
        HasTrippedAlert          bool           `json:"has_tripped_alert"`
    }

    type CircuitBreakerStateModel struct {
        Records        []CircuitBreakerStateRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord CircuitBreakerStateRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CircuitBreakerStateRecord (
        mx_host_key TEXT PRIMARY KEY,
        state_name TEXT NOT NULL,
        failure_count INTEGER NOT NULL DEFAULT 0,
        success_probe_count INTEGER NOT NULL DEFAULT 0,
        next_probe_allowed_at TEXT NOT NULL,
        is_closed INTEGER NOT NULL DEFAULT 1,
        is_open INTEGER NOT NULL DEFAULT 0,
        is_half_open INTEGER NOT NULL DEFAULT 0,
        has_tripped_alert INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Protects worker queue threads from hung sockets or connection timeouts when remote MX hosts undergo outages.

---

### Improvement 1344: TCP Fast Open (TFO) & TLS 1.3 0-RTT Connection Pipelining Governor
- **ID:** `BACKOFF-TCP-FASTOPEN-1344`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Optimize connection handshake overhead to remote MX servers using TCP Fast Open (RFC 7413) and TLS 1.3 session resumption with fallback protection against replay vulnerabilities.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [TCP FAST OPEN & TLS 1.3 SESSION RESUMPTION METRICS] ─────────────────────────────────────┐
      │ ACTIVE POOLED CONNS: 94 | TFO CACHE HITS: 78% | TLS 1.3 RESUMPTION: 84% | MEAN HANDSHAKE: 9ms │
      ├──────────────────────────┬──────────────┬──────────────┬──────────────┬────────────────────┤
      │ MX HOST                  │ TLS VERSION  │ TFO COOKIE   │ 0-RTT STATUS │ HANDSHAKE TIME     │
      ├──────────────────────────┼──────────────┼──────────────┼──────────────┼────────────────────┤
      │ aspmx.l.google.com       │ TLSv1.3      │ VALID (Cached│ ACCEPTED     │ 3.2ms  [●] OPTIMAL │
      │ mail.protection.outlook  │ TLSv1.3      │ VALID (Cached│ REJECTED(FB) │ 14.8ms [●] OPTIMAL │
      │ mx1.mail.yahoo.com       │ TLSv1.2      │ NOT SUPPORTED│ N/A (Full HS)│ 42.1ms [▲] ELEVATED│
      ├──────────────────────────┴──────────────┴──────────────┴──────────────┴────────────────────┤
      │ Handshake acceleration reduces SMTP latency from 45ms to 3.2ms on Google MX clusters.      │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup connection-pipeline`
  - Flags:
    - `--is-tfo-enabled`: Enable TCP Fast Open cookies for outbound SMTP dialer (default: `true`).
    - `--is-tls-resumption`: Store and reuse TLS 1.3 session tickets (default: `true`).
    - `--has-anti-replay-guard`: Reject 0-RTT early data for non-idempotent operations (default: `true`).
    - `--max-idle-conns <int>`: Maximum idle persistent SMTP connections per MX host (default: `10`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TlsHandshakeOptimizationRecord struct {
        ConnectionID             string         `json:"connection_id"`
        TargetMxHost             string         `json:"target_mx_host"`
        HandshakeDurationMs      float64        `json:"handshake_duration_ms"`
        TlsVersion               string         `json:"tls_version"`
        HasSessionReused         bool           `json:"has_session_reused"`
        IsTfoActive              bool           `json:"is_tfo_active"`
        IsZeroRttAccepted        bool           `json:"is_zero_rtt_accepted"`
        HasValidSessionTicket    bool           `json:"has_valid_session_ticket"`
        HasHandshakeAnomaly      bool           `json:"has_handshake_anomaly"`
    }

    type TlsHandshakeOptimizationModel struct {
        Records        []TlsHandshakeOptimizationRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord TlsHandshakeOptimizationRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TlsHandshakeOptimizationRecord (
        connection_id TEXT PRIMARY KEY,
        target_mx_host TEXT NOT NULL,
        handshake_duration_ms REAL NOT NULL DEFAULT 0.0,
        tls_version TEXT NOT NULL,
        has_session_reused INTEGER NOT NULL DEFAULT 0,
        is_tfo_active INTEGER NOT NULL DEFAULT 1,
        is_zero_rtt_accepted INTEGER NOT NULL DEFAULT 0,
        has_valid_session_ticket INTEGER NOT NULL DEFAULT 1,
        has_handshake_anomaly INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Low-level network layer inside GitMap Go dialer optimizing TCP/TLS handshakes.

---

### Improvement 1345: Predictive Greylisting Detection & Automated Rescheduling Engine
- **ID:** `BACKOFF-GREYLIST-RESCHED-1345`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Identify 450/451 Greylisting responses, accurately calculate remote MX retry expectations (e.g. 5m, 15m, 30m), and reschedule messages into an isolated low-priority retry ring without clogging live queues.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [PREDICTIVE GREYLISTING DETECTION & RETRY SCHEDULER] ─────────────────────────────────────┐
      │ DELAYED QUEUE: 23 | PROBE SUCCESS RATE: 95.6% | LEARNED GREYLIST SIGNATURES: 12 MX DOMAINS │
      ├──────────────────────────┬──────────────┬──────────────┬──────────────┬────────────────────┤
      │ DESTINATION DOMAIN       │ TARGET MX    │ DETECTED WAIT│ RETRY DUE IN │ STATUS             │
      ├──────────────────────────┼──────────────┼──────────────┼──────────────┼────────────────────┤
      │ posteo.de                │ mx1.posteo.de│ 300s (5 min) │ 01m 42s      │ [▲] WAITING RETRY  │
      │ gmx.net                  │ mx00.gmx.net │ 900s (15 min)│ 08m 10s      │ [▲] WAITING RETRY  │
      │ custom-corp.com          │ mx.mail.org  │ 600s (10 min)│ 00m 05s      │ [●] DISPATCHING NOW│
      ├──────────────────────────┴──────────────┴──────────────┴──────────────┴────────────────────┤
      │ Learned Posteo pattern: Re-try after 5 minutes succeeds 99.8% on second attempt.           │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup greylist-scheduler`
  - Flags:
    - `--min-greylist-delay <duration>`: Minimum initial wait before retrying greylisted MX (default: `5m`).
    - `--max-greylist-attempts <int>`: Maximum retry attempts before flagging hard error (default: `4`).
    - `--is-auto-isolated`: Shift greylisted domains to isolated delay queue (default: `true`).
    - `--has-mx-fingerprint`: Learn destination greylisting timeout patterns per MX provider (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type GreylistScheduleRecord struct {
        ScheduleID               string         `json:"schedule_id"`
        RecipientDomain          string         `json:"recipient_domain"`
        TargetMxHost             string         `json:"target_mx_host"`
        GreylistDelaySeconds     int            `json:"greylist_delay_seconds"`
        RetryAttempt             int            `json:"retry_attempt"`
        IsDelayed                bool           `json:"is_delayed"`
        IsFingerprintLearned     bool           `json:"is_fingerprint_learned"`
        HasExceededAttempts      bool           `json:"has_exceeded_attempts"`
        HasRetryDispatched       bool           `json:"has_retry_dispatched"`
    }

    type GreylistScheduleModel struct {
        Records        []GreylistScheduleRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord GreylistScheduleRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GreylistScheduleRecord (
        schedule_id TEXT PRIMARY KEY,
        recipient_domain TEXT NOT NULL,
        target_mx_host TEXT NOT NULL,
        greylist_delay_seconds INTEGER NOT NULL DEFAULT 300,
        retry_attempt INTEGER NOT NULL DEFAULT 1,
        is_delayed INTEGER NOT NULL DEFAULT 1,
        is_fingerprint_learned INTEGER NOT NULL DEFAULT 1,
        has_exceeded_attempts INTEGER NOT NULL DEFAULT 0,
        has_retry_dispatched INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Decouples greylist retries from the hot outbound queue into a dedicated scheduler ring in SQLite.

---

### Improvement 1346: Algorithmic IP Warmup Ramp Curve Governor & Volume Ceiling Enforcer
- **ID:** `BACKOFF-IP-WARM-RAMP-1346`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Enforce strict mathematical ramp curves (Linear, Exponential, Sigmoidal) for newly commissioned sending IPs, preventing premature high-volume spikes that trigger spam traps.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [ALGORITHMIC IP WARMUP RAMP CURVE GOVERNOR] ──────────────────────────────────────────────┐
      │ MONITORED IPS: 4 | RAMP TYPE: SIGMOIDAL | ACTIVE PROGRESSION: DAY 8 / 30                   │
      ├────────────────────┬──────────┬──────────────┬──────────────┬──────────────┬───────────────┤
      │ SENDING IP         │ DAY NO.  │ CEILING CAP  │ DISPATCHED   │ USAGE %      │ STATUS        │
      ├────────────────────┼──────────┼──────────────┼──────────────┼──────────────┼───────────────┤
      │ 198.51.100.41      │ Day 8    │ 500 msgs/day │ 412 msgs     │ 82.4%        │ [●] IN-FLIGHT │
      │ 198.51.100.42      │ Day 8    │ 500 msgs/day │ 500 msgs     │ 100.0%       │ [■] CAPPED    │
      │ 198.51.100.43      │ Day 1    │  50 msgs/day │  38 msgs     │ 76.0%        │ [●] IN-FLIGHT │
      ├────────────────────┴──────────┴──────────────┴──────────────┴──────────────┴───────────────┤
      │ [198.51.100.42] Daily volume ceiling met with 0 bounces. Next expansion: +250 at 00:00 UTC.│
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup ramp-curve`
  - Flags:
    - `--ip <string>`: Target sending IP to inspect or adjust.
    - `--curve-type <string>`: Ramp curve formula (`linear`, `exponential`, `sigmoid`; default: `sigmoid`).
    - `--current-day <int>`: Current progression day of the warmup schedule (1-30).
    - `--is-enforcing-ceiling`: Block messages exceeding current day ceiling (default: `true`).
    - `--has-auto-advance`: Automatically increment schedule day upon 24h clean run (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type IpWarmupScheduleRecord struct {
        ScheduleID               string         `json:"schedule_id"`
        SendingIP                string         `json:"sending_ip"`
        WarmupDay                int            `json:"warmup_day"`
        DailyCeilingVolume       int            `json:"daily_ceiling_volume"`
        DispatchedToday          int            `json:"dispatched_today"`
        CurveType                string         `json:"curve_type"`
        IsCeilingReached         bool           `json:"is_ceiling_reached"`
        IsCompliant              bool           `json:"is_compliant"`
        HasCleanHistory          bool           `json:"has_clean_history"`
        HasAdvancedToday         bool           `json:"has_advanced_today"`
    }

    type IpWarmupScheduleModel struct {
        Records        []IpWarmupScheduleRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord IpWarmupScheduleRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS IpWarmupScheduleRecord (
        schedule_id TEXT PRIMARY KEY,
        sending_ip TEXT NOT NULL,
        warmup_day INTEGER NOT NULL DEFAULT 1,
        daily_ceiling_volume INTEGER NOT NULL DEFAULT 50,
        dispatched_today INTEGER NOT NULL DEFAULT 0,
        curve_type TEXT NOT NULL,
        is_ceiling_reached INTEGER NOT NULL DEFAULT 0,
        is_compliant INTEGER NOT NULL DEFAULT 1,
        has_clean_history INTEGER NOT NULL DEFAULT 1,
        has_advanced_today INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Governs sending limits in GitMap dispatcher, persisting schedule state in `repodb/pipeline.db`.

---

### Improvement 1347: Smart SMTP Connection Multiplexing & Keep-Alive Session Pool
- **ID:** `BACKOFF-SMTP-SESSION-POOL-1347`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Maintain pooled, authenticated SMTP client sessions per MX host using RSET command re-use to transmit multiple messages per connection without repeated TCP/TLS handshakes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [SMART SMTP CONNECTION POOL & MULTIPLEXING HUD] ──────────────────────────────────────────┐
      │ ACTIVE POOLS: 18 | TOTAL OPEN SOCKETS: 42 | IDLE: 12 | RE-USE EFFICIENCY: 94.2%             │
      ├──────────────────────────┬────────────┬───────────┬──────────────┬─────────────────────────┤
      │ TARGET MX HOST           │ ACTIVE     │ IDLE      │ MSG / CONN   │ STATUS                  │
      ├──────────────────────────┼────────────┼───────────┼──────────────┼─────────────────────────┤
      │ aspmx.l.google.com       │ 14 sockets │ 4 sockets │ 84.2 msgs    │ [●] PIPELINING (RSET)   │
      │ mail.protection.outlook  │  8 sockets │ 2 sockets │ 42.1 msgs    │ [●] PIPELINING (RSET)   │
      │ mx1.mail.icloud.com      │  4 sockets │ 2 sockets │ 25.0 msgs    │ [●] PIPELINING (RSET)   │
      ├──────────────────────────┴────────────┴───────────┴──────────────┴─────────────────────────┤
      │ RSET re-use avoided 1,420 TCP 3-way handshakes and TLS negotiations in past 60 minutes.   │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup connection-pool`
  - Flags:
    - `--max-messages-per-conn <int>`: Maximum messages sent through single connection before close (default: `100`).
    - `--idle-timeout <duration>`: Keep-alive duration before closing idle session (default: `30s`).
    - `--is-rset-pipelined`: Issue RSET immediately following message acknowledgement (default: `true`).
    - `--has-leak-detector`: Verify closed socket descriptors to prevent file descriptor leaks (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SmtpConnectionPoolRecord struct {
        PoolID                   string         `json:"pool_id"`
        MxHostKey                string         `json:"mx_host_key"`
        ActiveConns              int            `json:"active_conns"`
        IdleConns                int            `json:"idle_conns"`
        TotalMessagesPiped       int            `json:"total_messages_piped"`
        IsPoolHealthy            bool           `json:"is_pool_healthy"`
        IsSaturated              bool           `json:"is_saturated"`
        HasFdLeak                bool           `json:"has_fd_leak"`
        HasKeepAliveSupported    bool           `json:"has_keep_alive_supported"`
    }

    type SmtpConnectionPoolModel struct {
        Records        []SmtpConnectionPoolRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SmtpConnectionPoolRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpConnectionPoolRecord (
        pool_id TEXT PRIMARY KEY,
        mx_host_key TEXT NOT NULL,
        active_conns INTEGER NOT NULL DEFAULT 0,
        idle_conns INTEGER NOT NULL DEFAULT 0,
        total_messages_piped INTEGER NOT NULL DEFAULT 0,
        is_pool_healthy INTEGER NOT NULL DEFAULT 1,
        is_saturated INTEGER NOT NULL DEFAULT 0,
        has_fd_leak INTEGER NOT NULL DEFAULT 0,
        has_keep_alive_supported INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Slashes connection overhead and CPU usage in high-speed GitMap email transmission operations.

---

### Improvement 1348: Machine Learning Heuristic Classifier for Provider Throttle Signatures
- **ID:** `BACKOFF-THROTTLE-LEARNER-1348`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Automatically categorize proprietary, non-standard 4xx error strings from Yahoo, Microsoft, and ProtonMail into structured backoff delay intervals using lightweight token classification.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [PROVIDER THROTTLE SIGNATURE PATTERN RECOGNITION] ────────────────────────────────────────┐
      │ LEARNED SIGNATURES: 38 | CLASSIFICATION ACCURACY: 99.4% | ADAPTIVE ADJUSTMENTS: 14         │
      ├────────────────────┬─────────────────────────────────┬──────────────┬──────────────────────┤
      │ PROVIDER           │ MATCHED ERROR PATTERN           │ ADAPTIVE DEL │ CONFIDENCE           │
      ├────────────────────┼─────────────────────────────────┼──────────────┼──────────────────────┤
      │ Yahoo Mail         │ 421 4.7.0 [TSS04] Message temp  │ 1,200s (20m) │ 99.8%  [●] VERIFIED  │
      │ Microsoft 365      │ 451 4.7.500 Server busy         │   300s (5m)  │ 98.9%  [●] VERIFIED  │
      │ Comcast Mail       │ 421 resend too quickly          │   600s (10m) │ 97.4%  [●] VERIFIED  │
      ├────────────────────┴─────────────────────────────────┴──────────────┴──────────────────────┤
      │ Rule applied: Yahoo TSS04 automatically triggers 20-minute quiet period on sending IP.     │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup throttle-classifier`
  - Flags:
    - `--sample-size <int>`: Number of recent error events to analyze (default: `100`).
    - `--is-realtime-update`: Update dynamic backoff rules immediately upon new signature (default: `true`).
    - `--has-feedback-loop`: Feed outcome of retry attempts back into classifier weights (default: `true`).
    - `--export-rules <string>`: Export learned rules table to YAML/JSON.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ThrottleSignatureRuleRecord struct {
        SignatureID              string         `json:"signature_id"`
        ProviderName             string         `json:"provider_name"`
        ErrorPatternRegex        string         `json:"error_pattern_regex"`
        RecommendedWaitSec       int            `json:"recommended_wait_sec"`
        ConfidenceRate           float64        `json:"confidence_rate"`
        IsRecognized             bool           `json:"is_recognized"`
        IsActiveRule             bool           `json:"is_active_rule"`
        HasSuccessfulRetry       bool           `json:"has_successful_retry"`
        HasAdaptiveTuning        bool           `json:"has_adaptive_tuning"`
    }

    type ThrottleSignatureRuleModel struct {
        Records        []ThrottleSignatureRuleRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord ThrottleSignatureRuleRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ThrottleSignatureRuleRecord (
        signature_id TEXT PRIMARY KEY,
        provider_name TEXT NOT NULL,
        error_pattern_regex TEXT NOT NULL,
        recommended_wait_sec INTEGER NOT NULL DEFAULT 300,
        confidence_rate REAL NOT NULL DEFAULT 0.95,
        is_recognized INTEGER NOT NULL DEFAULT 1,
        is_active_rule INTEGER NOT NULL DEFAULT 1,
        has_successful_retry INTEGER NOT NULL DEFAULT 1,
        has_adaptive_tuning INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Learns from incoming SMTP diagnostics to continuously refine the concurrency and backoff engine.

---

### Improvement 1349: Dynamic DNS MX Record TTL Observer & Route Caching Invalidator
- **ID:** `BACKOFF-DNS-TTL-OBSERVER-1349`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Track TTL values on destination MX DNS queries, actively invalidating cached IP mappings upon TTL expiration to prevent routing emails to decommissioned mail servers during provider migrations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [DNS MX RECORD TTL OBSERVER & CACHE CONTROLLER] ──────────────────────────────────────────┐
      │ CACHED MX DOMAINS: 1,240 | PREFETCH QUEUED: 18 | CACHE HIT RATIO: 98.7%                    │
      ├────────────────────┬────────────────────────┬──────────────┬──────────────┬────────────────┤
      │ DOMAIN             │ PRIMARY MX HOST        │ RESOLVED IP  │ TTL REMAIN   │ STATUS         │
      ├────────────────────┼────────────────────────┼──────────────┼──────────────┼────────────────┤
      │ proton.me          │ mail.protonmail.ch     │ 185.70.42.18 │ 184s         │ [●] VALID      │
      │ fastmail.com       │ in1-smtp.messagingengine│ 66.111.4.74  │ 42s (Prefetch│ [●] PREFETCHING│
      │ corp-client.com    │ mx1.corp-client.com    │ 203.0.113.88 │ 0s (Expired) │ [▲] REFRESHING │
      ├────────────────────┴────────────────────────┴──────────────┴──────────────┴────────────────┤
      │ Prefetch executed in background for fastmail.com. Zero query latency added to dispatch.   │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup dns-cache-observer`
  - Flags:
    - `--flush-domain <string>`: Manually flush and re-query MX records for target domain.
    - `--is-ttl-enforced`: Strictly respect upstream DNS authoritative TTLs (default: `true`).
    - `--has-prefetch`: Prefetch MX records 60 seconds before TTL expiration (default: `true`).
    - `--dns-resolver <string>`: Custom upstream DNS resolver (`system`, `1.1.1.1`, `8.8.8.8`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DnsMxCacheRecord struct {
        DnsQueryID               string         `json:"dns_query_id"`
        TargetDomain             string         `json:"target_domain"`
        ResolvedMxHost           string         `json:"resolved_mx_host"`
        ResolvedIP               string         `json:"resolved_ip"`
        TtlSecondsRemaining      int            `json:"ttl_seconds_remaining"`
        IsExpired                bool           `json:"is_expired"`
        IsValid                  bool           `json:"is_valid"`
        HasPrefetchScheduled     bool           `json:"has_prefetch_scheduled"`
        HasDnssecVerified        bool           `json:"has_dnssec_verified"`
    }

    type DnsMxCacheModel struct {
        Records        []DnsMxCacheRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord DnsMxCacheRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DnsMxCacheRecord (
        dns_query_id TEXT PRIMARY KEY,
        target_domain TEXT NOT NULL,
        resolved_mx_host TEXT NOT NULL,
        resolved_ip TEXT NOT NULL,
        ttl_seconds_remaining INTEGER NOT NULL DEFAULT 300,
        is_expired INTEGER NOT NULL DEFAULT 0,
        is_valid INTEGER NOT NULL DEFAULT 1,
        has_prefetch_scheduled INTEGER NOT NULL DEFAULT 1,
        has_dnssec_verified INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Protects against stale DNS cache routing in long-running GitMap email sending daemons.

---

### Improvement 1350: Automated Multi-Homed MX Priority Failover & Fallback Coordinator
- **ID:** `BACKOFF-FAILOVER-FALLBACK-1350`
- **Cluster:** 135 (Dynamic Backoff Inversion & Multi-Domain Concurrency Governor)
- **Goal:** Handle MX record preference weighting (e.g. 10, 20, 30), automatically routing traffic to secondary or tertiary backup MX endpoints when primary MX hosts reject connections or timeout.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [MULTI-HOMED MX PRIORITY FAILOVER & FALLBACK] ────────────────────────────────────────────┐
      │ MONITORED DOMAINS: 320 | ACTIVE FAILOVER ROUTES: 2 | PRIMARY RESTORATION PROBES: 2/2        │
      ├────────────────────┬───────────┬────────────────────────────┬─────────────┬────────────────┤
      │ DOMAIN             │ ACTIVE MX │ HOSTNAME                   │ PREFERENCE  │ STATUS         │
      ├────────────────────┼───────────┼────────────────────────────┼─────────────┼────────────────┤
      │ enterprise.com     │ PRIMARY   │ mx1.enterprise.com         │ 10          │ [●] ONLINE     │
      │ disaster-test.org  │ SECONDARY │ mx2.disaster-test.org      │ 20          │ [▲] FAILOVER   │
      │ global-corp.net    │ PRIMARY   │ mx-a.global-corp.net       │ 10          │ [●] ONLINE     │
      ├────────────────────┴───────────┴────────────────────────────┴─────────────┴────────────────┤
      │ disaster-test.org primary MX timed out after 3 retries. Traffic shifted to preference 20.  │
      │ Background restoration probe running every 60s.                                            │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup mx-failover`
  - Flags:
    - `--domain <string>`: Inspect priority failover chain for specific domain.
    - `--is-fallback-auto`: Fail back to primary MX immediately when health check recovers (default: `true`).
    - `--probe-interval <duration>`: Interval to verify primary MX restoration (default: `60s`).
    - `--has-weight-balancing`: Distribute load evenly across MX hosts of identical preference (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MxFailoverRouteRecord struct {
        RouteID                  string         `json:"route_id"`
        DomainName               string         `json:"domain_name"`
        PrimaryMx                string         `json:"primary_mx"`
        ActiveFailoverMx         string         `json:"active_failover_mx"`
        CurrentPreference        int            `json:"current_preference"`
        IsPrimaryActive          bool           `json:"is_primary_active"`
        IsFailoverEngaged        bool           `json:"is_failover_engaged"`
        HasHealthyBackup         bool           `json:"has_healthy_backup"`
        HasRestorationProbe      bool           `json:"has_restoration_probe"`
    }

    type MxFailoverRouteModel struct {
        Records        []MxFailoverRouteRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord MxFailoverRouteRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxFailoverRouteRecord (
        route_id TEXT PRIMARY KEY,
        domain_name TEXT NOT NULL,
        primary_mx TEXT NOT NULL,
        active_failover_mx TEXT NOT NULL,
        current_preference INTEGER NOT NULL DEFAULT 10,
        is_primary_active INTEGER NOT NULL DEFAULT 1,
        is_failover_engaged INTEGER NOT NULL DEFAULT 0,
        has_healthy_backup INTEGER NOT NULL DEFAULT 1,
        has_restoration_probe INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Ensures uninterrupted delivery by dynamically traversing RFC 5321 MX preference trees.

---

# Cluster 136: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (Items 1351–1360)

---

### Improvement 1351: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD
- **ID:** `EMAIL-SANDBOX-HUD-1351`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Provide a terminal-based HTML AST parser and zero-CDN interactive sandbox auditing email payloads for security threats (embedded scripts, external tracker beacons, malicious iframes, CSS injection).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [ZERO-CDN EMAIL CONTENT SECURITY SANDBOX & SANITIZER] ────────────────────────────────────┐
      │ PAYLOAD: welcome-warmup-v2.html | SIZE: 14.2 KB | SECURITY SCORE: 98.5% | CSP: STRICT       │
      ├────────────────────┬───────────┬──────────────────────────────────────────┬────────────────┤
      │ THREAT CATEGORY    │ SEVERITY  │ ELEMENT / ATTRIBUTE MATCHED              │ ACTION TAKEN   │
      ├────────────────────┼───────────┼──────────────────────────────────────────┼────────────────┤
      │ Script Tag         │ CRITICAL  │ <script src="http://untrusted...">       │ STRIPPED       │
      │ Tracking Pixel     │ MEDIUM    │ <img src="https://trk.pixel/1x1.png"...> │ STRIPPED       │
      │ Unsafe CSS Style   │ LOW       │ style="behavior: url(...)"               │ PURGED         │
      ├────────────────────┴───────────┴──────────────────────────────────────────┴────────────────┤
      │ Sanitized payload rendered in terminal preview. Zero external network calls executed.       │
      │ Interactive: [v: Render Raw Tree] [e: Export Clean HTML] [c: View CSP Rules] [q: Quit]     │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox inspect`
  - Flags:
    - `--file <string>`: Path to HTML email payload file to sanitize and inspect.
    - `--is-strict-csp`: Enforce strict zero-external-network Content Security Policy (default: `true`).
    - `--has-tracker-detection`: Scan for zero-pixel tracking GIFs and CSS background beacons (default: `true`).
    - `--output-sanitized <string>`: Export sanitized secure HTML to designated path.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailContentSanitizerRecord struct {
        AuditID                  string         `json:"audit_id"`
        PayloadChecksum          string         `json:"payload_checksum"`
        StrippedTagsCount        int            `json:"stripped_tags_count"`
        DetectedTrackersCount    int            `json:"detected_trackers_count"`
        CspComplianceScore       float64        `json:"csp_compliance_score"`
        IsClean                  bool           `json:"is_clean"`
        IsModified               bool           `json:"is_modified"`
        HasScriptInjection       bool           `json:"has_script_injection"`
        HasExternalBeacon        bool           `json:"has_external_beacon"`
    }

    type EmailContentSanitizerModel struct {
        Records        []EmailContentSanitizerRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord EmailContentSanitizerRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailContentSanitizerRecord (
        audit_id TEXT PRIMARY KEY,
        payload_checksum TEXT NOT NULL,
        stripped_tags_count INTEGER NOT NULL DEFAULT 0,
        detected_trackers_count INTEGER NOT NULL DEFAULT 0,
        csp_compliance_score REAL NOT NULL DEFAULT 100.0,
        is_clean INTEGER NOT NULL DEFAULT 1,
        is_modified INTEGER NOT NULL DEFAULT 0,
        has_script_injection INTEGER NOT NULL DEFAULT 0,
        has_external_beacon INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Protects warmup recipients and inspection TUIs from executing untrusted dynamic HTML/JS content.

---

### Improvement 1352: Interactive Email Dark Mode & Color Contrast TUI Simulator
- **ID:** `EMAIL-DARKMODE-SIM-1352`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Simulate and verify email rendering across Light Mode, Partial Dark Mode (Outlook/Gmail auto-inversion), and Full Dark Mode (Apple Mail @media prefers-color-scheme) with WCAG 2.1 AA contrast auditing.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [EMAIL DARK MODE RENDERING & WCAG AA CONTRAST SIM] ──────────────────────────────────────┐
      │ TEMPLATE: newsletter-digest.html | MODE: Apple Mail Full Dark Mode | CONTRAST RATIO: 7.2:1 │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ ┌─ [DARK MODE PREVIEW] ──────────────────────────────────────────────────────────────────┐ │
      │ │ #121212 BACKGROUND | #FFFFFF HEADING 1 (Contrast: 18.2:1 [✓] AAA)                      │ │
      │ │ #E0E0E0 BODY TEXT | #38BDF8 ACCENT LINK (Contrast: 6.8:1 [✓] AA)                       │ │
      │ │                                                                                        │ │
      │ │ [!] DETECTED INVERSION HAZARD: <img src="dark-logo.png"> has black transparent bounds. │ │
      │ │ Recommendation: Add style="filter: drop-shadow(0 0 1px #fff)" or provide SVG icon.     │ │
      │ └────────────────────────────────────────────────────────────────────────────────────────┘ │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Status: [✓] WCAG 2.1 AA PASSED | 0 Critical Contrast Violations                            │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox dark-mode`
  - Flags:
    - `--html <string>`: Path to input email HTML template.
    - `--sim-mode <string>`: Simulation mode (`light`, `dark-partial`, `dark-full`; default: `dark-full`).
    - `--is-wcag-audit`: Verify text/background luminance ratio >= 4.5:1 (default: `true`).
    - `--has-inversion-guard`: Inject vendor CSS rules preventing broken logo inversions (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailDarkModeSimRecord struct {
        SimID                    string         `json:"sim_id"`
        TemplateName             string         `json:"template_name"`
        SimulatedMode            string         `json:"simulated_mode"`
        LowestContrastRatio      float64        `json:"lowest_contrast_ratio"`
        ViolationsCount          int            `json:"violations_count"`
        IsWcagCompliant          bool           `json:"is_wcag_compliant"`
        IsLegible                bool           `json:"is_legible"`
        HasAutoInversionBug      bool           `json:"has_auto_inversion_bug"`
        HasFixInjected           bool           `json:"has_fix_injected"`
    }

    type EmailDarkModeSimModel struct {
        Records        []EmailDarkModeSimRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord EmailDarkModeSimRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailDarkModeSimRecord (
        sim_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        simulated_mode TEXT NOT NULL,
        lowest_contrast_ratio REAL NOT NULL DEFAULT 4.5,
        violations_count INTEGER NOT NULL DEFAULT 0,
        is_wcag_compliant INTEGER NOT NULL DEFAULT 1,
        is_legible INTEGER NOT NULL DEFAULT 1,
        has_auto_inversion_bug INTEGER NOT NULL DEFAULT 0,
        has_fix_injected INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Ensures warmup and marketing templates render legibly across all major mobile/desktop dark themes.

---

### Improvement 1353: Responsive Viewport Breakpoint & Fluid Hybrid Layout Tester
- **ID:** `EMAIL-RESPONSIVE-BREAKPOINT-1353`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Test and validate fluid-hybrid email layout behavior across Mobile (320px-480px), Tablet (600px-768px), and Desktop (1024px+) viewports in pure terminal ANSI blocks.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [RESPONSIVE VIEWPORT BREAKPOINT & FLUID-HYBRID TEST] ─────────────────────────────────────┐
      │ VIEWPORT: 375px (Mobile Portrait) | GHOST TABLES: 4 PRESENT | HORIZONTAL SCROLL: NONE      │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ ┌─ [375px WIREFRAME BOUNDS] ─────────────────────────────────────────────────────────────┐ │
      │ │ ┌────────────────────────────────────────────────────────────────────────────────────┐ │ │
      │ │ │ HEADER (375px fluid) -> Logo centered (120px)                                      │ │ │
      │ │ ├────────────────────────────────────────────────────────────────────────────────────┤ │ │
      │ │ │ 2-COLUMN COLLAPSE:                                                                 │ │ │
      │ │ │ Column 1 (was 300px desktop) -> 100% (345px padded) [STACKED]                      │ │ │
      │ │ │ Column 2 (was 300px desktop) -> 100% (345px padded) [STACKED]                      │ │ │
      │ │ ├────────────────────────────────────────────────────────────────────────────────────┤ │ │
      │ │ │ CTA BUTTON: 100% width button (min-height: 48px, accessible tap target)            │ │ │
      │ │ └────────────────────────────────────────────────────────────────────────────────────┘ │ │
      │ └────────────────────────────────────────────────────────────────────────────────────────┘ │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Layout Validation: [✓] FLUID-HYBRID PASS (Zero horizontal overflow detected).              │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox responsive-check`
  - Flags:
    - `--viewport <string>`: Test target viewport width (`mobile-375`, `tablet-768`, `desktop-1024`; default: `mobile-375`).
    - `--is-fluid-hybrid`: Validate ghost table compatibility for Outlook MSO desktop engines (default: `true`).
    - `--has-horizontal-overflow`: Detect elements exceeding viewport bounding box (default: `true`).
    - `--output-terminal-grid`: Render layout wireframe directly to terminal standard output.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailResponsiveCheckRecord struct {
        CheckID                  string         `json:"check_id"`
        TemplateName             string         `json:"template_name"`
        TargetViewportWidth      int            `json:"target_viewport_width"`
        MaxContentWidth          int            `json:"max_content_width"`
        GhostTablesDetected      int            `json:"ghost_tables_detected"`
        IsResponsive             bool           `json:"is_responsive"`
        IsFittingViewport        bool           `json:"is_fitting_viewport"`
        HasHorizontalOverflow    bool           `json:"has_horizontal_overflow"`
        HasMsoConditionals       bool           `json:"has_mso_conditionals"`
    }

    type EmailResponsiveCheckModel struct {
        Records        []EmailResponsiveCheckRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord EmailResponsiveCheckRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailResponsiveCheckRecord (
        check_id TEXT PRIMARY KEY,
        template_name TEXT NOT NULL,
        target_viewport_width INTEGER NOT NULL DEFAULT 375,
        max_content_width INTEGER NOT NULL DEFAULT 375,
        ghost_tables_detected INTEGER NOT NULL DEFAULT 4,
        is_responsive INTEGER NOT NULL DEFAULT 1,
        is_fitting_viewport INTEGER NOT NULL DEFAULT 1,
        has_horizontal_overflow INTEGER NOT NULL DEFAULT 0,
        has_mso_conditionals INTEGER NOT NULL DEFAULT 1
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Validates complex HTML email layout stacking before transmission without requiring cloud rendering engines.

---

### Improvement 1354: Lexical AST Spam Trigger Word & Misleading Subject Analyzer
- **ID:** `EMAIL-SPAM-TRIGGER-AST-1354`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan subject lines and HTML text nodes using a tokenizer AST to detect spam trigger keywords (e.g. '100% free', 'make money fast', all-caps strings, deceptive Re: / Fwd: markers).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [LEXICAL AST SPAM TRIGGER & SUBJECT REPUTATION SCORE] ────────────────────────────────────┐
      │ SUBJECT: "Quick sync regarding project roadmap" | ESTIMATED SPAM SCORE: 0.2 / 5.0 (SAFE)  │
      ├────────────────────┬───────────┬──────────────┬────────────────────────────────────────────┤
      │ TOKEN / PHRASE     │ WEIGHT    │ STATUS       │ REASON & RECOMMENDATION                    │
      ├────────────────────┼───────────┼──────────────┼────────────────────────────────────────────┤
      │ "Quick sync"       │ +0.10     │ LOW RISK     │ Natural conversational phrasing            │
      │ "project roadmap"  │ +0.05     │ LOW RISK     │ Contextual business token                  │
      │ Capitalization     │ 0.00      │ CLEAN        │ Sentence case compliant                    │
      ├────────────────────┴───────────┴──────────────┴────────────────────────────────────────────┤
      │ Overall Status: [✓] EXCELLENT (Clean natural language; zero commercial spam markers).     │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox spam-score`
  - Flags:
    - `--subject <string>`: Subject line string to audit.
    - `--body-file <string>`: Path to email plain text or HTML body file.
    - `--is-strict-scoring`: Enforce strict SpamAssassin / Rspamd heuristic thresholds (default: `true`).
    - `--has-remediation-suggestions`: Generate alternative lexical suggestions for flagged words (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailSpamLexicalRecord struct {
        AuditID                  string         `json:"audit_id"`
        SubjectText              string         `json:"subject_text"`
        HeuristicSpamScore       float64        `json:"heuristic_spam_score"`
        FlaggedKeywordsCount     int            `json:"flagged_keywords_count"`
        AllCapsRatio             float64        `json:"all_caps_ratio"`
        IsSpamThresholdClean     bool           `json:"is_spam_threshold_clean"`
        IsSubjectHonest          bool           `json:"is_subject_honest"`
        HasDeceptivePrefix       bool           `json:"has_deceptive_prefix"`
        HasUrgencyTriggers       bool           `json:"has_urgency_triggers"`
    }

    type EmailSpamLexicalModel struct {
        Records        []EmailSpamLexicalRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord EmailSpamLexicalRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailSpamLexicalRecord (
        audit_id TEXT PRIMARY KEY,
        subject_text TEXT NOT NULL,
        heuristic_spam_score REAL NOT NULL DEFAULT 0.0,
        flagged_keywords_count INTEGER NOT NULL DEFAULT 0,
        all_caps_ratio REAL NOT NULL DEFAULT 0.0,
        is_spam_threshold_clean INTEGER NOT NULL DEFAULT 1,
        is_subject_honest INTEGER NOT NULL DEFAULT 1,
        has_deceptive_prefix INTEGER NOT NULL DEFAULT 0,
        has_urgency_triggers INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - CLI command in `gitmap email sandbox spam-score` calculating offline spam heuristics.

---

### Improvement 1355: Zero-Script SVG Sanitizer & Safe Vector Inlining Engine
- **ID:** `EMAIL-SVG-SANITIZER-INLINE-1355`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Sanitize embedded or linked SVG vector graphics in email bodies by stripping XML entity expansions (XXE), `<foreignObject>`, `<script>`, and dangerous event handlers, converting safe vectors to inline MIME parts.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [ZERO-SCRIPT SVG SANITIZER & XXE SECURITY AUDIT] ─────────────────────────────────────────┐
      │ TARGET SVG: assets/logo-vector.svg | ORIGINAL: 8.4 KB | SANITIZED: 3.1 KB (-63% bytes)     │
      ├──────────────────────────┬──────────────┬───────────────────────────┬──────────────────────┤
      │ PARSER INSPECTION CHECK  │ OCCURRENCES  │ DETAIL                    │ ACTION TAKEN         │
      ├──────────────────────────┼──────────────┼───────────────────────────┼──────────────────────┤
      │ XML External Entities    │ 0 detected   │ DTD processing disabled   │ [✓] IMMUNE           │
      │ <foreignObject> Tags     │ 1 detected   │ Embedded HTML block       │ [✗] STRIPPED SAFELY  │
      │ onload / onerror hooks   │ 0 detected   │ No inline JS attributes   │ [✓] CLEAN            │
      │ Inkscape Metadata        │ 4 tags       │ Redundant authoring data  │ [✓] PURGED           │
      ├──────────────────────────┴──────────────┴───────────────────────────┴──────────────────────┤
      │ Safe vector generated. Ready for CID inline attachment or base64 CSS embedding.           │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox sanitize-svg`
  - Flags:
    - `--input-svg <string>`: Path to source SVG graphic.
    - `--is-xxe-immune`: Disable all DTD and external XML entity resolution (default: `true`).
    - `--has-lossless-opt`: Strip redundant metadata, editor comments, and namespaces (default: `true`).
    - `--output-safe <string>`: Path to save sanitized safe SVG vector.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SvgSanitizationRecord struct {
        SvgID                    string         `json:"svg_id"`
        SourceFilePath           string         `json:"source_file_path"`
        OriginalBytes            int            `json:"original_bytes"`
        SanitizedBytes           int            `json:"sanitized_bytes"`
        RemovedThreatsCount      int            `json:"removed_threats_count"`
        IsSafe                   bool           `json:"is_safe"`
        IsOptimized              bool           `json:"is_optimized"`
        HasXmlEntityExploit      bool           `json:"has_xml_entity_exploit"`
        HasForeignObject         bool           `json:"has_foreign_object"`
    }

    type SvgSanitizationModel struct {
        Records        []SvgSanitizationRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord SvgSanitizationRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SvgSanitizationRecord (
        svg_id TEXT PRIMARY KEY,
        source_file_path TEXT NOT NULL,
        original_bytes INTEGER NOT NULL DEFAULT 0,
        sanitized_bytes INTEGER NOT NULL DEFAULT 0,
        removed_threats_count INTEGER NOT NULL DEFAULT 0,
        is_safe INTEGER NOT NULL DEFAULT 1,
        is_optimized INTEGER NOT NULL DEFAULT 1,
        has_xml_entity_exploit INTEGER NOT NULL DEFAULT 0,
        has_foreign_object INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Secures all outbound vector assets from XML and HTML injection vulnerabilities.

---

### Improvement 1356: Autonomous Offline Asset Bundler & RFC 2392 CID Attachment Packager
- **ID:** `EMAIL-ZERO-CDN-ASSET-BUNDLER-1356`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Bundle external image URLs into self-contained MIME multipart/related structures with unique Content-ID (CID) headers, enabling 100% offline rendering without CDN dependencies or external image blocking.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [OFFLINE ASSET BUNDLER & RFC 2392 CID PACKAGER] ──────────────────────────────────────────┐
      │ TEMPLATE: onboarding-email.html | BUNDLED ASSETS: 3 IMAGES | TOTAL MIME SIZE: 142 KB        │
      ├────────────────────┬───────────┬──────────────┬─────────────────────────────┬──────────────┤
      │ ASSET NAME         │ TYPE      │ DISK SIZE    │ CONTENT-ID (CID)            │ INLINE STATE │
      ├────────────────────┼───────────┼──────────────┼─────────────────────────────┼──────────────┤
      │ header-logo.png    │ image/png │ 24.2 KB      │ <logo@gitmap.internal>      │ [●] INLINED  │
      │ hero-diagram.png   │ image/png │ 86.4 KB      │ <hero@gitmap.internal>      │ [●] INLINED  │
      │ check-icon.png     │ image/png │  4.1 KB      │ <check@gitmap.internal>     │ [●] INLINED  │
      ├────────────────────┴───────────┴──────────────┴─────────────────────────────┴──────────────┤
      │ Zero external CDN dependencies. Recipient clients display images without clicking         │
      │ "Load Remote Content". Complete privacy protection guaranteed.                             │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox bundle-assets`
  - Flags:
    - `--html-in <string>`: Input HTML template containing image references.
    - `--is-cid-inlined`: Convert `src="http..."` to `src="cid:..."` MIME pointers (default: `true`).
    - `--max-bundle-bytes <int>`: Maximum allowable MIME multipart payload size (default: `1048576`).
    - `--has-compression`: Losslessly compress PNG/JPEG attachments prior to encoding (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailAssetBundleRecord struct {
        BundleID                 string         `json:"bundle_id"`
        HtmlTemplatePath         string         `json:"html_template_path"`
        BundledImagesCount       int            `json:"bundled_images_count"`
        TotalMimeBytes           int            `json:"total_mime_bytes"`
        IsZeroCdnCompliant       bool           `json:"is_zero_cdn_compliant"`
        IsWithinSizeCap          bool           `json:"is_within_size_cap"`
        HasLosslessOptimized     bool           `json:"has_lossless_optimized"`
        HasBrokenReferences      bool           `json:"has_broken_references"`
    }

    type EmailAssetBundleModel struct {
        Records        []EmailAssetBundleRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord EmailAssetBundleRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailAssetBundleRecord (
        bundle_id TEXT PRIMARY KEY,
        html_template_path TEXT NOT NULL,
        bundled_images_count INTEGER NOT NULL DEFAULT 0,
        total_mime_bytes INTEGER NOT NULL DEFAULT 0,
        is_zero_cdn_compliant INTEGER NOT NULL DEFAULT 1,
        is_within_size_cap INTEGER NOT NULL DEFAULT 1,
        has_lossless_optimized INTEGER NOT NULL DEFAULT 1,
        has_broken_references INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Bundles emails directly inside GitMap sending orchestrator for offline, privacy-first delivery.

---

### Improvement 1357: Automated CAN-SPAM, GDPR & List-Unsubscribe One-Click Header Guard
- **ID:** `EMAIL-CAN-SPAM-FOOTER-1357`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Validate and inject compliant postal address blocks, privacy policy links, and RFC 8058 List-Unsubscribe-Post one-click unsubscribe headers into email headers and body footers.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [CAN-SPAM, GDPR & RFC 8058 ONE-CLICK UNSUBSCRIBE AUDIT] ──────────────────────────────────┐
      │ STATUS: [✓] FULLY COMPLIANT | YAHOO / GMAIL 2026 BULK SENDER CRITERIA: SATISFIED            │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [MANDATORY HEADERS GENERATED]                                                              │
      │ List-Unsubscribe: <https://gitmap.io/unsub?id=c891>, <mailto:unsub@gitmap.io?subject=unsub>│
      │ List-Unsubscribe-Post: List-Unsubscribe=One-Click                                          │
      │                                                                                            │
      │ [FOOTER AUDIT INSPECTION]                                                                  │
      │ ├─ Physical Postal Address: "100 Innovation Way, Suite 400, Austin, TX 78701" [FOUND]     │
      │ ├─ Direct One-Click Unsubscribe Link: Present and functioning                              │
      │ └─ Legible Font Size: 12px (>= 10px minimum requirement)                                   │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox compliance-footer`
  - Flags:
    - `--physical-address <string>`: Registered company physical postal mailing address.
    - `--unsubscribe-url <string>`: HTTPS endpoint for RFC 8058 one-click unsubscription.
    - `--is-rfc8058-enforced`: Inject List-Unsubscribe and List-Unsubscribe-Post headers (default: `true`).
    - `--has-gdpr-consent-text`: Include explicit opt-in consent and data processing disclosure (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailComplianceFooterRecord struct {
        ComplianceID             string         `json:"compliance_id"`
        MailID                   string         `json:"mail_id"`
        UnsubscribeHttpsUrl      string         `json:"unsubscribe_https_url"`
        PhysicalAddressString    string         `json:"physical_address_string"`
        IsCanSpamValid           bool           `json:"is_can_spam_valid"`
        IsGdprCompliant          bool           `json:"is_gdpr_compliant"`
        IsRfc8058Active          bool           `json:"is_rfc8058_active"`
        HasOneClickHeader        bool           `json:"has_one_click_header"`
        HasMissingPostalAddress  bool           `json:"has_missing_postal_address"`
    }

    type EmailComplianceFooterModel struct {
        Records        []EmailComplianceFooterRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord EmailComplianceFooterRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailComplianceFooterRecord (
        compliance_id TEXT PRIMARY KEY,
        mail_id TEXT NOT NULL,
        unsubscribe_https_url TEXT NOT NULL,
        physical_address_string TEXT NOT NULL,
        is_can_spam_valid INTEGER NOT NULL DEFAULT 1,
        is_gdpr_compliant INTEGER NOT NULL DEFAULT 1,
        is_rfc8058_active INTEGER NOT NULL DEFAULT 1,
        has_one_click_header INTEGER NOT NULL DEFAULT 1,
        has_missing_postal_address INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Validates every outbound message against legal compliance regulations before queue entry.

---

### Improvement 1358: Bilingual HTML-to-Plaintext AST Synchronizer & MIME Multiparter
- **ID:** `EMAIL-PLAIN-TEXT-AST-SYNC-1358`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Automatically generate an accessible, beautifully formatted `text/plain` alternative MIME part from `text/html` email templates, preserving links, lists, and hierarchy to prevent spam filter penalties.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [HTML-TO-PLAINTEXT AST SYNCHRONIZER & MULTIPART AUDIT] ───────────────────────────────────┐
      │ HTML SIZE: 14.2 KB | PLAINTEXT SIZE: 1.8 KB | WORD RATIO: 1.00 (EXACT MATCH)               │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ ┌─ [GENERATED TEXT/PLAIN PREVIEW (Wrapped at 78 cols)] ──────────────────────────────────┐ │
      │ │ Welcome to GitMap Architecture!                                                        │ │
      │ │                                                                                        │ │
      │ │ We have compiled your weekly engineering sync summary:                                 │ │
      │ │ * Completed 12 High-Priority Specs [1]                                                 │ │
      │ │ * Resolved Greylisting Delays on European ESPs                                         │ │
      │ │                                                                                        │ │
      │ │ ────────────────────────────────────────────────────────────────────────────────────── │ │
      │ │ [1] https://gitmap.io/docs/specs                                                       │ │
      │ └────────────────────────────────────────────────────────────────────────────────────────┘ │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ Validation: [✓] MULTIPART/ALTERNATIVE GENERATED (Prevents HTML-only spam penalties).      │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox sync-plaintext`
  - Flags:
    - `--html-input <string>`: Path to source HTML email template.
    - `--wrap-width <int>`: Maximum line wrap column for plain text (default: `78`).
    - `--is-link-footnotes`: Convert inline links to numbered reference footnotes (default: `true`).
    - `--has-semantic-bulleting`: Render bullet points as clean ASCII asterisks or dashes (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PlaintextSyncRecord struct {
        SyncID                   string         `json:"sync_id"`
        HtmlBytes                int            `json:"html_bytes"`
        PlaintextBytes           int            `json:"plaintext_bytes"`
        WordCountRatio           float64        `json:"word_count_ratio"`
        FootnotesCount           int            `json:"footnotes_count"`
        IsSynced                 bool           `json:"is_synced"`
        IsMimeCompliant          bool           `json:"is_mime_compliant"`
        HasTextPartMissing       bool           `json:"has_text_part_missing"`
        HasRatioPenalty          bool           `json:"has_ratio_penalty"`
    }

    type PlaintextSyncModel struct {
        Records        []PlaintextSyncRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord PlaintextSyncRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PlaintextSyncRecord (
        sync_id TEXT PRIMARY KEY,
        html_bytes INTEGER NOT NULL DEFAULT 0,
        plaintext_bytes INTEGER NOT NULL DEFAULT 0,
        word_count_ratio REAL NOT NULL DEFAULT 1.0,
        footnotes_count INTEGER NOT NULL DEFAULT 0,
        is_synced INTEGER NOT NULL DEFAULT 1,
        is_mime_compliant INTEGER NOT NULL DEFAULT 1,
        has_text_part_missing INTEGER NOT NULL DEFAULT 0,
        has_ratio_penalty INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Generates synchronized multi-part payloads in GitMap, eliminating missing plaintext spam scoring penalties.

---

### Improvement 1359: Embedded Hyperlink Security & Reputation Destination Gatekeeper
- **ID:** `EMAIL-LINK-REPUTATION-GATE-1359`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Scan all `href` URLs in email bodies against Google Safe Browsing, PhishTank, and URLhaus to prevent sending links that lead to compromised or blacklisted web destinations.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [EMBEDDED HYPERLINK REPUTATION & REDIRECT GATEKEEPER] ────────────────────────────────────┐
      │ EXTRACTED URLS: 4 | REDIRECT CHAINS CHECKED: 4 | CLEAN: 4 | PHISHING MATCHES: 0            │
      ├──────────────────────────┬──────┬────────────────────────────┬──────────────┬──────────────┤
      │ EMBEDDED URL             │ HOPS │ FINAL DESTINATION          │ REPUTATION   │ SAFETY STATE │
      ├──────────────────────────┼──────┼────────────────────────────┼──────────────┼──────────────┤
      │ https://gitmap.io/login  │ 0    │ https://gitmap.io/login    │ 100/100 HIGH │ [●] SAFE     │
      │ https://gitmap.io/docs   │ 1    │ https://gitmap.io/docs/v6  │ 100/100 HIGH │ [●] SAFE     │
      │ https://go.dev           │ 0    │ https://go.dev             │ 100/100 HIGH │ [●] SAFE     │
      ├──────────────────────────┴──────┴────────────────────────────┴──────────────┴──────────────┤
      │ Zero link shorteners (bit.ly, t.co) or suspicious raw IPs found. Domain trust verified.    │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox audit-links`
  - Flags:
    - `--html-file <string>`: Path to email HTML payload to scan for embedded hyperlinks.
    - `--is-deep-redirect-resolve`: Trace HTTP 301/302 redirect chains to final destination (default: `true`).
    - `--has-ip-url-block`: Block raw IP addresses inside href attributes (default: `true`).
    - `--timeout <duration>`: Timeout per link verification probe (default: `2s`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailLinkReputationRecord struct {
        LinkAuditID              string         `json:"link_audit_id"`
        TargetUrl                string         `json:"target_url"`
        ResolvedFinalUrl         string         `json:"resolved_final_url"`
        SafetyReputationScore    float64        `json:"safety_reputation_score"`
        RedirectHops             int            `json:"redirect_hops"`
        IsSafe                   bool           `json:"is_safe"`
        IsSuspicious             bool           `json:"is_suspicious"`
        HasRedirectChain         bool           `json:"has_redirect_chain"`
        HasBlacklistMatch        bool           `json:"has_blacklist_match"`
    }

    type EmailLinkReputationModel struct {
        Records        []EmailLinkReputationRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord EmailLinkReputationRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailLinkReputationRecord (
        link_audit_id TEXT PRIMARY KEY,
        target_url TEXT NOT NULL,
        resolved_final_url TEXT NOT NULL,
        safety_reputation_score REAL NOT NULL DEFAULT 100.0,
        redirect_hops INTEGER NOT NULL DEFAULT 0,
        is_safe INTEGER NOT NULL DEFAULT 1,
        is_suspicious INTEGER NOT NULL DEFAULT 0,
        has_redirect_chain INTEGER NOT NULL DEFAULT 0,
        has_blacklist_match INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Gatekeeper tool ensuring zero outbound email leads to malicious redirectors or blacklisted domains.

---

### Improvement 1360: Interactive Bubble Tea MIME Structure & Boundary Tree Visualizer
- **ID:** `EMAIL-MIME-TREE-EXPLORER-1360`
- **Cluster:** 136 (Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD)
- **Goal:** Render an interactive Bubble Tea terminal tree of raw MIME multipart headers, nested boundaries, character encodings, and attachment transfer encodings (quoted-printable, base64) with inline hex dumps.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - High-density Lipgloss TUI interface with border styling (`#38bdf8`, sky blue) and adaptive status indicators (`[●]` healthy emerald `#10b981`, `[▲]` warning amber `#f59e0b`, `[■]` critical rose `#f43f5e`).
    - Model view layout rendering:
      ```text
      ┌─ [INTERACTIVE BUBBLE TEA MIME STRUCTURE & BOUNDARY DISSECTOR] ─────────────────────────────┐
      │ MESSAGE-ID: <20261027.9184@gitmap.internal> | SIZE: 148.4 KB | PARTS: 4 | SYNTAX: RFC 2046 │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ ▼ multipart/mixed; boundary="====_MIXED_BOUNDARY_01_===="                                  │
      │   ├─ ▼ multipart/alternative; boundary="====_ALT_BOUNDARY_02_===="                       │
      │   │   ├─ [Part 1] text/plain; charset="utf-8"; transfer-encoding="quoted-printable" (1.8K)  │
      │   │   └─ [Part 2] text/html; charset="utf-8"; transfer-encoding="quoted-printable" (14.2K)  │
      │   ├─ [Part 3] image/png; name="logo.png"; Content-ID="<logo@gitmap>"; base64 (24.2K)       │
      │   └─ [Part 4] application/pdf; name="summary.pdf"; base64 (108.2K)                        │
      ├────────────────────────────────────────────────────────────────────────────────────────────┤
      │ [SELECTED PART 3: image/png] SHA-256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca49 │
      │ Controls: [↑/↓: Select Part] [enter: Hex Viewer] [x: Export Part] [q: Quit]              │
      └────────────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure SVG coordinate plots, responsive flexbox CSS containers, and semantic HTML5 cards with zero external CDN scripts or stylesheet dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email sandbox mime-tree`
  - Flags:
    - `--raw-eml <string>`: Path to raw .eml RFC 5322 message file to dissect.
    - `--is-interactive`: Enable Bubble Tea interactive tree expansion and folding (default: `true`).
    - `--has-hex-preview`: Show hex/ASCII byte viewer for selected MIME payload (default: `true`).
    - `--verify-checksums`: Verify SHA-256 integrity of all attached payload parts (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MimeStructureAuditRecord struct {
        MimeTreeID               string         `json:"mime_tree_id"`
        MessageID                string         `json:"message_id"`
        RootContentType          string         `json:"root_content_type"`
        TotalPartsCount          int            `json:"total_parts_count"`
        TotalMimeSize            int            `json:"total_mime_size"`
        IsStructureValid         bool           `json:"is_structure_valid"`
        IsCleanBoundaries        bool           `json:"is_clean_boundaries"`
        HasNestedMultiparts      bool           `json:"has_nested_multiparts"`
        HasEncodingError         bool           `json:"has_encoding_error"`
    }

    type MimeStructureAuditModel struct {
        Records        []MimeStructureAuditRecord `json:"records"`
        CursorIndex    int                    `json:"cursorIndex"`
        SelectedRecord MimeStructureAuditRecord   `json:"selectedRecord"`
        IsStreaming    bool                   `json:"isStreaming"`
        HasAlerts      bool                   `json:"hasAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MimeStructureAuditRecord (
        mime_tree_id TEXT PRIMARY KEY,
        message_id TEXT NOT NULL,
        root_content_type TEXT NOT NULL,
        total_parts_count INTEGER NOT NULL DEFAULT 1,
        total_mime_size INTEGER NOT NULL DEFAULT 0,
        is_structure_valid INTEGER NOT NULL DEFAULT 1,
        is_clean_boundaries INTEGER NOT NULL DEFAULT 1,
        has_nested_multiparts INTEGER NOT NULL DEFAULT 1,
        has_encoding_error INTEGER NOT NULL DEFAULT 0
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Interactive TUI inspector in GitMap for deep diagnosis of corrupted MIME headers and broken multipart boundaries.

---
