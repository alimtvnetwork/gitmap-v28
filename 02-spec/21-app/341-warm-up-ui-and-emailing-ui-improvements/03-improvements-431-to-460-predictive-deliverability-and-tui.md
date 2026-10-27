# Specification 341: Improvements 431 to 460 (Predictive Deliverability, Bubble Tea Hex-Grid HUD & Distributed SMTP Tracer)

- **Sequence:** 341
- **Slug:** `341-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/03-improvements-431-to-460-predictive-deliverability-and-tui.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (431 - 460)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 44** | Items 431–440 | Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD |
| **Cluster 45** | Items 441–450 | Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing |
| **Cluster 46** | Items 451–460 | Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector |

---

# Cluster 44: Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD (Items 431–440)

---

### Improvement 431: Interactive 2D Hex-Grid Terminal Swarm Topography Matrix
- **ID:** `WARMUP-TUI-431`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Provide operators with an interactive, terminal-native hexagonal node grid rendering up to 64 sender nodes in real-time with health-coded ANSI glyphs, live hover inspectors, and directional arrow traversal.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea & Lipgloss Terminal HUD:**
    - Hexagonal coordinate projection mapping sending nodes into staggered rows (`/ \_ / \`).
    - Color tokens: Healthy (`#10b981`, bold green `⬡`), Warming (`#3b82f6`, blue `⬢`), Rate-limited (`#f59e0b`, amber `⬡`), Tripped/Degraded (`#ef4444`, flashing crimson `⬢`).
    - Real-time cursor focus with border styling: `lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#a855f7")).Padding(0, 1)`.
    - Live contextual footer rendering focused node IP, PTR record, current hourly velocity, and active warm-up stage.
  - **Zero-CDN Web UI:**
    - Inline SVG hex-tile canvas (`viewBox="0 0 960 480"`) using pure CSS variable tokens (`var(--status-healthy)`).
    - Client-side click handlers for immediate drawer slide-out displaying telemetry timeline without external JS bundles.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup hud hex-grid`
  - Flags:
    - `--pool <string>`: Target warmup pool identifier (default: `default`).
    - `--refresh-ms <int>`: Terminal rendering refresh rate in milliseconds (default: `250`).
    - `--filter-status <string>`: Filter grid by state (`all`, `healthy`, `degraded`, `tripped`).
    - `--watch`: Persistent interactive Bubble Tea full-screen alt-buffer mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HexNodeTelemetry struct {
        NodeID           string    `json:"nodeId"`
        IPv4Address      string    `json:"ipv4Address"`
        WarmupStage      int       `json:"warmupStage"`
        HourlyVelocity   int       `json:"hourlyVelocity"`
        HealthScore      float64   `json:"healthScore"`
        IsHealthy        bool      `json:"isHealthy"`
        IsWarming        bool      `json:"isWarming"`
        HasRateLimitFlag bool      `json:"hasRateLimitFlag"`
        IsTripped        bool      `json:"isTripped"`
        LastHeartbeat    time.Time `json:"lastHeartbeat"`
    }

    type HexGridState struct {
        PoolName    string             `json:"poolName"`
        TotalNodes  int                `json:"totalNodes"`
        ActiveNodes []HexNodeTelemetry `json:"activeNodes"`
        IsLive      bool               `json:"isLive"`
        HasErrors   bool               `json:"hasErrors"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupHexNodes (
        nodeId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        ipv4Address TEXT NOT NULL,
        warmupStage INTEGER NOT NULL DEFAULT 1,
        hourlyVelocity INTEGER NOT NULL DEFAULT 0,
        healthScore REAL NOT NULL DEFAULT 100.0,
        isHealthy INTEGER NOT NULL DEFAULT 1,
        isWarming INTEGER NOT NULL DEFAULT 1,
        hasRateLimitFlag INTEGER NOT NULL DEFAULT 0,
        isTripped INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_warmup_hex_pool ON WarmupHexNodes(poolId, isHealthy);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements `cli/cmdwarmup/hex_hud.go` utilizing GitMap's custom Bubble Tea event loop, streaming telemetry from Lemon background daemon socket.

---

### Improvement 432: Multi-Window Terminal Splitter with Real-Time Flow Velocity Sparklines
- **ID:** `WARMUP-TUI-432`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Display simultaneous dual-pane or quad-pane terminal windows showing real-time sender dispatch velocity sparklines alongside incoming bounce telemetry streams.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dynamic window layout using `bubbles/viewport` partitioned vertically or horizontally based on terminal resize signals (`tea.WindowSizeMsg`).
    - Braille UTF-8 sparklines (` ▂▃▄▅▆▇█`) rendering rolling 60-second delivery throughput in pane 1, and connection latency histograms in pane 2.
    - Lipgloss border tabs (`Active Relays [1]`, `Bounces [2]`, `Inbound DSN [3]`) switchable via hotkeys `Tab` and `Shift+Tab`.
  - **Zero-CDN Web UI:**
    - CSS Grid dashboard (`grid-template-columns: repeat(auto-fit, minmax(320px, 1fr))`) containing pure SVG path sparkline elements.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup monitor split`
  - Flags:
    - `--panes <int>`: Number of panes to instantiate (2 or 4, default: `2`).
    - `--sparkline-window <duration>`: Rolling window for velocity points (default: `60s`).
    - `--focus-pane <int>`: Default focused sub-window index (default: `1`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SparklineStreamSample struct {
        TimestampRFC3339 string  `json:"timestamp"`
        ThroughputEPS    float64 `json:"throughputEps"`
        LatencyMS        float64 `json:"latencyMs"`
        BounceCount      int     `json:"bounceCount"`
        IsThrottled      bool    `json:"isThrottled"`
        HasWarning       bool    `json:"hasWarning"`
    }
    ```
  - Split-DB SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TerminalSplitSessions (
        sessionId TEXT PRIMARY KEY,
        paneCount INTEGER NOT NULL DEFAULT 2,
        windowSeconds INTEGER NOT NULL DEFAULT 60,
        isActive INTEGER NOT NULL DEFAULT 1,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Hooks into `cli/term/splitter.go` to cleanly detach and isolate UI render loops without blocking outbound SMTP worker goroutines.

---

### Improvement 433: Interactive Hotkey Ramp Velocity Throttle Dial & Override Switch
- **ID:** `WARMUP-TUI-433`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Enable operators to live-tune warmup ramp speeds (+/- 5%, 10%, 25%) directly inside the terminal via keyboard hotkeys (`+`, `-`, `space`, `p`) with instant Redis/SQLite state synchronization.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea TUI Dial:**
    - Semicircular progress bar dial displaying current multiplier (e.g., `1.25x [Ramp Turbo]`).
    - Emergency killswitch toggle hotkey (`Shift+K` or `Ctrl+C` safe halt) instantly setting velocity to `0` and triggering graceful connection draining.
    - Status message alert bar flashing yellow for manual override: `[OVERRIDE ACTIVE: +15% Manual Boost]`.
  - **Zero-CDN Web UI:**
    - Pure HTML `<input type="range">` styled with CSS custom properties, updating local storage and posting JSON commands via Fetch API without third-party frameworks.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup throttle tune`
  - Flags:
    - `--pool <string>`: Name of warmup pool (default: `production`).
    - `--multiplier <float>`: Target velocity multiplier (e.g., `0.5`, `1.0`, `1.5`).
    - `--emergency-drain`: If true, immediately halts injection and drains open TLS handshakes.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ThrottleOverrideEvent struct {
        PoolID             string    `json:"poolId"`
        BaseVelocity       int       `json:"baseVelocity"`
        AdjustedVelocity   int       `json:"adjustedVelocity"`
        Multiplier         float64   `json:"multiplier"`
        IsManualOverride   bool      `json:"isManualOverride"`
        IsEmergencyDrained bool      `json:"isEmergencyDrained"`
        HasOperatorAuth    bool      `json:"hasOperatorAuth"`
        ModifiedAt         time.Time `json:"modifiedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupThrottleOverrides (
        overrideId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        multiplier REAL NOT NULL DEFAULT 1.0,
        isManualOverride INTEGER NOT NULL DEFAULT 0,
        isEmergencyDrained INTEGER NOT NULL DEFAULT 0,
        hasOperatorAuth INTEGER NOT NULL DEFAULT 1,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates with `pkg/warmup/throttler.go` via thread-safe atomic pointer swaps for lock-free velocity re-evaluation.

---

### Improvement 434: ISP Quota Boundary Visualizer with Color-Coded Radial Guages
- **ID:** `WARMUP-TUI-434`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Display real-time consumption of hourly and daily recipient quotas broken down by receiving ISP (Google, Microsoft 365, Yahoo, Proton, iCloud) in terminal radial gauges.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Component:**
    - Custom radial quadrant gauge rendered using Unicode quadrant and block elements (`▰▰▰▱▱▱ 60%`).
    - Color ramp thresholds: <70% Green (`#10b981`), 70-89% Yellow (`#eab308`), >=90% Orange (`#f97316`), 100% Red (`#ef4444`).
    - Real-time countdown timer to rolling hour reset (`Quota reset in: 14m 32s`).
  - **Zero-CDN Web UI:**
    - SVG donut charts (`stroke-dasharray`, `stroke-dashoffset`) embedded directly in HTML template with zero external CDN dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup quota inspect`
  - Flags:
    - `--isp <string>`: Target ISP filter (`google`, `microsoft`, `yahoo`, `all`).
    - `--pool <string>`: Warmup sending pool (default: `default`).
    - `--json`: Output raw quota metrics in structured JSON format.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ISPQuotaMetric struct {
        ISPName          string    `json:"ispName"`
        CurrentHourSent  int       `json:"currentHourSent"`
        HourlyMaxQuota   int       `json:"hourlyMaxQuota"`
        CurrentDaySent   int       `json:"currentDaySent"`
        DailyMaxQuota    int       `json:"dailyMaxQuota"`
        IsSoftCapped     bool      `json:"isSoftCapped"`
        IsHardCapped     bool      `json:"isHardCapped"`
        HasReachedTarget bool      `json:"hasReachedTarget"`
        WindowResetTime  time.Time `json:"windowResetTime"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupISPQuotas (
        quotaId TEXT PRIMARY KEY,
        ispName TEXT NOT NULL,
        currentHourSent INTEGER NOT NULL DEFAULT 0,
        hourlyMaxQuota INTEGER NOT NULL DEFAULT 100,
        currentDaySent INTEGER NOT NULL DEFAULT 0,
        dailyMaxQuota INTEGER NOT NULL DEFAULT 1000,
        isSoftCapped INTEGER NOT NULL DEFAULT 0,
        isHardCapped INTEGER NOT NULL DEFAULT 0,
        hasReachedTarget INTEGER NOT NULL DEFAULT 0,
        resetAt DATETIME NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_warmup_isp ON WarmupISPQuotas(ispName, resetAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Connects to `pkg/warmup/quota_limiter.go` to preemptively backpressure sender queues before ISP throttling triggers 421/451 temporary deferrals.

---

### Improvement 435: Live Interactive Terminal Filter Search for In-Flight Envelopes
- **ID:** `WARMUP-TUI-435`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Provide a high-performance Fuzzy-Search input field (`bubbles/textinput`) allowing operators to filter active in-flight warmup envelopes by recipient domain, relay IP, or current SMTP phase.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal Component:**
    - Pinned header text input with prefix `Search [ / ]: ` styled with cyan foreground (`#06b6d4`).
    - Real-time tabular result list powered by fuzzy matching algorithm with matching characters underlined.
    - Live summary counter: `Showing 14 of 420 active envelopes (Filtered by: "yahoo.com")`.
  - **Zero-CDN Web UI:**
    - Standard HTML5 `<input type="search">` tied to vanilla JavaScript table row filter evaluating `dataset.domain`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup search inflight`
  - Flags:
    - `--query <string>`: Initial search query string.
    - `--limit <int>`: Max results to render simultaneously (default: `25`).
    - `--sort <string>`: Sort field (`time`, `recipient`, `ip`, `stage`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type InflightEnvelopeSummary struct {
        EnvelopeID       string    `json:"envelopeId"`
        RecipientDomain  string    `json:"recipientDomain"`
        SenderIP         string    `json:"senderIp"`
        CurrentPhase     string    `json:"currentPhase"`
        ElapsedTimeMS    int64     `json:"elapsedTimeMs"`
        IsTLSActive      bool      `json:"isTlsActive"`
        HasResponseError bool      `json:"hasResponseError"`
        CreatedAt        time.Time `json:"createdAt"`
    }
    ```
  - Split-DB SQLite Schema (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InflightWarmupEnvelopes (
        envelopeId TEXT PRIMARY KEY,
        recipientDomain TEXT NOT NULL,
        senderIp TEXT NOT NULL,
        currentPhase TEXT NOT NULL,
        elapsedTimeMs INTEGER NOT NULL DEFAULT 0,
        isTlsActive INTEGER NOT NULL DEFAULT 1,
        hasResponseError INTEGER NOT NULL DEFAULT 0,
        dispatchedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_inflight_domain ON InflightWarmupEnvelopes(recipientDomain);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `cli/term/fuzzy_envelope.go` to provide instant feedback during heavy warm-up load testing sessions.

---

### Improvement 436: Bubble Tea Progress Timeline for Tiered Warmup Schedules
- **ID:** `WARMUP-TUI-436`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Render a structured 30-day tiered warmup schedule in terminal tables with step milestones, percentage completions, and projected target dates.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Timeline Display:**
    - Vertical timeline graph with tree branching glyphs (`├── [Day 05]`, `└── [Day 30]`).
    - Completed days marked with green checkmark (`✔`), current day highlighted with blue pulsing badge (`● Day 12 - 450 emails/day`), and future days dimmed (`○`).
    - Lipgloss table borders with soft purple styling (`#8b5cf6`).
  - **Zero-CDN Web UI:**
    - Clean semantic HTML `<ol class="warmup-timeline">` with CSS flexbox styling and responsive mobile layout.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup schedule show`
  - Flags:
    - `--pool <string>`: Warmup schedule pool ID (default: `default`).
    - `--tier-plan <string>`: Plan archetype (`conservative`, `standard`, `aggressive`).
    - `--start-date <date>`: Custom origin start date (`YYYY-MM-DD`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type WarmupScheduleStep struct {
        DayNumber        int       `json:"dayNumber"`
        TargetVolume     int       `json:"targetVolume"`
        ActualSent       int       `json:"actualSent"`
        CompletionRate   float64   `json:"completionRate"`
        IsCurrentDay     bool      `json:"isCurrentDay"`
        IsCompleted      bool      `json:"isCompleted"`
        HasVolumeSurpassed bool    `json:"hasVolumeSurpassed"`
        ScheduledDate    time.Time `json:"scheduledDate"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupTierSchedules (
        scheduleId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        dayNumber INTEGER NOT NULL,
        targetVolume INTEGER NOT NULL,
        actualSent INTEGER NOT NULL DEFAULT 0,
        isCurrentDay INTEGER NOT NULL DEFAULT 0,
        isCompleted INTEGER NOT NULL DEFAULT 0,
        hasVolumeSurpassed INTEGER NOT NULL DEFAULT 0,
        scheduledDate DATE NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_warmup_sched ON WarmupTierSchedules(poolId, dayNumber);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements `pkg/warmup/schedule_viewer.go` ensuring alignment with GitMap automated plan sequencing standards.

---

### Improvement 437: Real-Time Bounced Envelope Hex Highlight & Auto-Isolation HUD
- **ID:** `WARMUP-TUI-437`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Visually isolate and flash degraded or hard-bouncing nodes on the Bubble Tea hex HUD, presenting interactive keybinds (`i` to isolate, `r` to recycle socket).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Node tile dynamically changes glyph to inverted crimson warning symbol (`▲ ERR`) when bounce count > 3 within 60 seconds.
    - Sound bell trigger (`\a`) optional flag for terminal alerts upon critical isolation.
    - Quick-action modal popover displaying exact SMTP error code (`554 5.7.1 Relay Access Denied`).
  - **Zero-CDN Web UI:**
    - Live CSS keyframe pulsing border animation on affected node cards (`@keyframes alert-pulse`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup node isolate`
  - Flags:
    - `--node-id <string>`: Target node identifier.
    - `--reason <string>`: Isolation rationale (`hard-bounce-spike`, `spamhaus-flag`, `manual`).
    - `--auto-drain`: Automatically transfer existing queued jobs to healthy peers.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type NodeIsolationRecord struct {
        NodeID             string    `json:"nodeId"`
        IsolationReason    string    `json:"isolationReason"`
        RecentBounceCount  int       `json:"recentBounceCount"`
        IsIsolated         bool      `json:"isIsolated"`
        HasDrainedJobs     bool      `json:"hasDrainedJobs"`
        IsolatedAt         time.Time `json:"isolatedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupNodeIsolations (
        isolationId TEXT PRIMARY KEY,
        nodeId TEXT NOT NULL,
        isolationReason TEXT NOT NULL,
        recentBounceCount INTEGER NOT NULL DEFAULT 0,
        isIsolated INTEGER NOT NULL DEFAULT 1,
        hasDrainedJobs INTEGER NOT NULL DEFAULT 0,
        isolatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Directly updates `pkg/warmup/circuit_guard.go` and notifies active swarm supervisors.

---

### Improvement 438: Interactive Latency Waterfall Inspector in Bubble Tea Viewport
- **ID:** `WARMUP-TUI-438`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Render a detailed step-by-step latency waterfall for individual SMTP warmup transactions (DNS -> TCP Connect -> TLS Handshake -> EHLO -> MAIL FROM -> RCPT TO -> DATA -> 250 OK) inside a scrollable viewport.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Viewport:**
    - ASCII Gantt-style horizontal bars visualizing milliseconds spent in each phase:
      `DNS      [██░░░░░░░░]  12ms`
      `TCP      [████░░░░░░]  24ms`
      `TLS      [██████████]  68ms`
      `SMTP OK  [███░░░░░░░]  18ms`
    - Color gradients shifting from cyan (`#06b6d4`) to violet (`#8b5cf6`).
  - **Zero-CDN Web UI:**
    - Pure HTML/CSS flexbox waterfall bars with hover tooltips displaying socket connection state and TCP window sizes.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup trace waterfall`
  - Flags:
    - `--envelope-id <string>`: Target envelope tracking UUID.
    - `--verbose`: Include raw socket TCP ACK timing and cipher payload sizes.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SMTPWaterfallMetrics struct {
        EnvelopeID       string  `json:"envelopeId"`
        DNSLookupMS      float64 `json:"dnsLookupMs"`
        TCPConnectMS     float64 `json:"tcpConnectMs"`
        TLSHandshakeMS   float64 `json:"tlsHandshakeMs"`
        EHLOExchangeMS   float64 `json:"ehloExchangeMs"`
        DataTransferMS   float64 `json:"dataTransferMs"`
        TotalDurationMS  float64 `json:"totalDurationMs"`
        IsSuccess        bool    `json:"isSuccess"`
        HasLagAnomaly    bool    `json:"hasLagAnomaly"`
    }
    ```
  - Split-DB SQLite Schema (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupWaterfallLogs (
        traceId TEXT PRIMARY KEY,
        envelopeId TEXT NOT NULL,
        dnsLookupMs REAL NOT NULL,
        tcpConnectMs REAL NOT NULL,
        tlsHandshakeMs REAL NOT NULL,
        ehloExchangeMs REAL NOT NULL,
        dataTransferMs REAL NOT NULL,
        totalDurationMs REAL NOT NULL,
        isSuccess INTEGER NOT NULL DEFAULT 1,
        hasLagAnomaly INTEGER NOT NULL DEFAULT 0,
        recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_waterfall_env ON WarmupWaterfallLogs(envelopeId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `cli/term/waterfall.go` using high-resolution monotonic clocks (`time.Now()`).

---

### Improvement 439: Terminal Heatmap Distribution of Send Velocity Across 24h Windows
- **ID:** `WARMUP-TUI-439`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Display an interactive 7-day x 24-hour terminal ANSI heatmap showing hourly warmup density, optimal ISP deliverability windows, and detected rate-limiting hours.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - 7 rows (Mon-Sun) and 24 columns (00h-23h) grid using block intensity glyphs:
      `░` (Low volume), `▒` (Medium volume), `▓` (High volume), `█` (Peak volume).
    - Color ramp using green and emerald tones, with red markers on hours containing throttle deferrals.
    - Arrow key navigation across cells with detail popover showing exact send count and average open rate.
  - **Zero-CDN Web UI:**
    - Native CSS grid heatmap with zero JS libraries, using CSS native variables `color-mix` for color interpolation.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup heatmap show`
  - Flags:
    - `--days <int>`: Number of days to render (default: `7`).
    - `--pool <string>`: Sender warmup pool (default: `default`).
    - `--ansi-ramp <string>`: Color palette selection (`emerald`, `amber`, `monochrome`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HourlyHeatmapCell struct {
        DayOfWeek         int     `json:"dayOfWeek"`
        HourOfDay         int     `json:"hourOfDay"`
        TotalSent         int     `json:"totalSent"`
        BounceRate        float64 `json:"bounceRate"`
        IsOptimalWindow   bool    `json:"isOptimalWindow"`
        HasThrottleEvents bool    `json:"hasThrottleEvents"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupHourlyHeatmap (
        cellId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        dayOfWeek INTEGER NOT NULL,
        hourOfDay INTEGER NOT NULL,
        totalSent INTEGER NOT NULL DEFAULT 0,
        bounceRate REAL NOT NULL DEFAULT 0.0,
        isOptimalWindow INTEGER NOT NULL DEFAULT 0,
        hasThrottleEvents INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX IF NOT EXISTS idx_warmup_heatmap ON WarmupHourlyHeatmap(poolId, dayOfWeek, hourOfDay);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements `cli/term/heatmap_table.go` leveraging GitMap's existing `termtable` package.

---

### Improvement 440: Swarm Peer-to-Peer Consensus Health Summary Bar
- **ID:** `WARMUP-TUI-440`
- **Cluster:** 44 (Interactive Terminal Warmup Flow Matrix & Live Bubble Tea Hex-Grid HUD)
- **Goal:** Provide a pinned top HUD status bar displaying real-time Lemon mesh peer consensus, active Raft leader, pool health ratio, and sync latency across warmup worker nodes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Status Bar:**
    - Pinned top header bar using high-contrast slate background:
      `LEMON MESH: LEADER [node-us-west] | PEERS: 12/12 ONLINE | HEALTH: 98.4% | DRIFT: 1.2ms`
    - Micro-indicators: Blue dot (`●`) for syncing, Green dot (`●`) for synced, Orange (`▲`) for quorum warning.
  - **Zero-CDN Web UI:**
    - Sticky top navbar with CSS flexbox layout and offline indicator via `navigator.onLine`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap warmup swarm status`
  - Flags:
    - `--timeout <duration>`: Consensus discovery timeout (default: `5s`).
    - `--compact`: Output single-line condensed format for integration into external scripts.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SwarmConsensusSummary struct {
        LeaderNodeID     string    `json:"leaderNodeId"`
        ActivePeersCount int       `json:"activePeersCount"`
        TotalPeersCount  int       `json:"totalPeersCount"`
        QuorumHealthPct  float64   `json:"quorumHealthPct"`
        SyncDriftMS      float64   `json:"syncDriftMs"`
        IsQuorumReached  bool      `json:"isQuorumReached"`
        HasNetworkSplit  bool      `json:"hasNetworkSplit"`
        LastSyncedAt     time.Time `json:"lastSyncedAt"`
    }
    ```
  - Split-DB SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupSwarmConsensus (
        consensusId TEXT PRIMARY KEY,
        leaderNodeId TEXT NOT NULL,
        activePeersCount INTEGER NOT NULL,
        totalPeersCount INTEGER NOT NULL,
        quorumHealthPct REAL NOT NULL,
        syncDriftMs REAL NOT NULL,
        isQuorumReached INTEGER NOT NULL DEFAULT 1,
        hasNetworkSplit INTEGER NOT NULL DEFAULT 0,
        syncedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Ingests telemetry directly from GitMap's Lemon agent communication bus (`pkg/lemon/consensus.go`).

---

# Cluster 45: Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing (Items 441–450)

---

### Improvement 441: Zero-CDN Client Email Canvas with Live Resizing & Responsive Rulers
- **ID:** `WARMUP-CANVAS-441`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Provide a local, zero-external-CDN HTML canvas previewing warmup email templates with live draggable viewport handles (Desktop 600px, Tablet 480px, Mobile 320px) and pixel-accurate dimension rulers.
- **UI/UX Terminal & Web Specification:**
  - **Zero-CDN Web UI:**
    - Sandboxed `<iframe>` wrapper with `sandbox="allow-same-origin"` blocking external network tracking pixels.
    - SVG-based draggable slider controls to dynamically adjust iframe width in real-time without JavaScript libraries.
    - High-precision CSS rulers along top and left margins with 10px tick marks.
  - **Lipgloss Terminal UX:**
    - ASCII terminal dimension badge: `[PREVIEW CANVAS: 600x800px | Mobile-Responsive: OK]`.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email canvas preview`
  - Flags:
    - `--template <file>`: Path to local HTML email template file.
    - `--port <int>`: Local preview HTTP server port (default: `8999`).
    - `--open-browser`: Automatically launch system default browser.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CanvasPreviewConfig struct {
        TemplatePath     string `json:"templatePath"`
        TargetWidthPX    int    `json:"targetWidthPx"`
        TargetHeightPX   int    `json:"targetHeightPx"`
        IsResponsiveMode bool   `json:"isResponsiveMode"`
        HasExternalFonts bool   `json:"hasExternalFonts"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailCanvasTemplates (
        templateId TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        htmlContent TEXT NOT NULL,
        renderedWidthPx INTEGER NOT NULL DEFAULT 600,
        isResponsiveMode INTEGER NOT NULL DEFAULT 1,
        hasExternalFonts INTEGER NOT NULL DEFAULT 0,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements lightweight zero-dependency HTTP preview handler in `pkg/email/canvas_server.go`.

---

### Improvement 442: WCAG 2.1 AA/AAA Terminal Contrast Matrix Analyzer for Email HTML
- **ID:** `WARMUP-CANVAS-442`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Scan all text and background CSS color declarations across warmup email HTML templates and calculate WCAG 2.1 contrast ratios in a terminal matrix, flagging unreadable body text.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Tabular contrast report displaying Element Selector, Foreground Hex, Background Hex, Calculated Ratio, and Compliance Badges (`AA PASS [4.8:1]`, `AAA FAIL [4.8:1 < 7.0:1]`).
    - Compliant elements styled in bold emerald green; failing elements highlighted in amber/crimson with suggested replacement color values.
  - **Zero-CDN Web UI:**
    - Interactive DOM node inspector highlighting low-contrast paragraphs directly on the preview iframe with an overlay badge.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit contrast`
  - Flags:
    - `--template <file>`: HTML template path to analyze.
    - `--level <string>`: Target WCAG compliance level (`AA` or `AAA`, default: `AA`).
    - `--strict`: Return non-zero CLI exit code on any AAA failure.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ContrastAuditEntry struct {
        Selector        string  `json:"selector"`
        ForegroundHex   string  `json:"foregroundHex"`
        BackgroundHex   string  `json:"backgroundHex"`
        CalculatedRatio float64 `json:"calculatedRatio"`
        IsAAPass        bool    `json:"isAaPass"`
        IsAAAPass       bool    `json:"isAaaPass"`
        HasLargeText    bool    `json:"hasLargeText"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailContrastAudits (
        auditId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        selector TEXT NOT NULL,
        calculatedRatio REAL NOT NULL,
        isAaPass INTEGER NOT NULL DEFAULT 1,
        isAaaPass INTEGER NOT NULL DEFAULT 0,
        hasLargeText INTEGER NOT NULL DEFAULT 0,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `pkg/email/contrast_checker.go` using standard sRGB luminance formula per W3C specification.

---

### Improvement 443: Dark Mode Inversion Simulator (Gmail, Apple Mail & Outlook Schemes)
- **ID:** `WARMUP-CANVAS-443`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Autonomously simulate client-specific CSS dark mode mutations (Full Inversion, Partial Color Inversion, and Prefers-Color-Scheme media queries) to detect invisible logos and illegible typography.
- **UI/UX Terminal & Web Specification:**
  - **Zero-CDN Web UI:**
    - Segmented client toggle buttons: `Standard Dark`, `Gmail Android Inversion`, `Outlook Windows Inversion`, `Apple Mail Clean`.
    - Real-time CSS filter application (`filter: invert(100%) hue-rotate(180deg)`) simulating native mail client render bugs.
  - **Lipgloss Terminal UX:**
    - High-level warning summary:
      `[!] Gmail Inversion Warning: Logo PNG lacks transparent padding border`
      `[✔] Apple Mail: @media (prefers-color-scheme: dark) correctly specified`
- **CLI Command Signature & Flags:**
  - Command: `gitmap email simulate dark-mode`
  - Flags:
    - `--template <file>`: Target HTML template file.
    - `--client <string>`: Mail client engine (`gmail`, `outlook`, `apple`, `all`).
    - `--dump-transformed-html`: Write converted dark-mode HTML to disk for diffing.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DarkModeSimulationResult struct {
        ClientEngine       string   `json:"clientEngine"`
        InversionType      string   `json:"inversionType"`
        BrokenElementCount int      `json:"brokenElementCount"`
        IsLegible          bool     `json:"isLegible"`
        HasDarkModeCSS     bool     `json:"hasDarkModeCss"`
        DetectedIssues     []string `json:"detectedIssues"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailDarkModeSimulations (
        simulationId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        clientEngine TEXT NOT NULL,
        brokenElementCount INTEGER NOT NULL DEFAULT 0,
        isLegible INTEGER NOT NULL DEFAULT 1,
        hasDarkModeCss INTEGER NOT NULL DEFAULT 0,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements AST HTML parser in `pkg/email/darkmode_sim.go` inspecting inline styles and CSS media blocks.

---

### Improvement 444: Embedded Asset Size & 102KB Gmail Clipping Warning Inspector
- **ID:** `WARMUP-CANVAS-444`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Calculate the exact total byte count of email HTML envelopes (including inline CSS, base64 images, and tracking tags), raising strict terminal alerts if the size nears or exceeds Gmail's 102KB clipping cutoff.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Progress meter: `Total HTML Size: 84.2 KB / 102.0 KB [████████░░] 82%`.
    - Warning banner if >100KB: `[WARNING: High Risk of "[Message clipped] View entire message" link in Gmail]`.
    - Breakdown table itemizing raw markup bytes, inline style bytes, and embedded base64 bytes.
  - **Zero-CDN Web UI:**
    - Gauge indicator with red threshold line at 102KB and one-click "Minify HTML Markup" action button.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit size`
  - Flags:
    - `--template <file>`: Input template file.
    - `--max-bytes <int>`: Custom ceiling limit in bytes (default: `104448` - 102KB).
    - `--minify`: Provide projected compressed size after stripping redundant whitespaces and comments.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EmailSizeAudit struct {
        RawByteSize         int  `json:"rawByteSize"`
        MinifiedByteSize    int  `json:"minifiedByteSize"`
        InlineCSSBytes      int  `json:"inlineCssBytes"`
        EmbeddedBase64Bytes int  `json:"embeddedBase64Bytes"`
        IsUnderGmailLimit   bool `json:"isUnderGmailLimit"`
        IsClipped           bool `json:"isClipped"`
        HasBase64Assets     bool `json:"hasBase64Assets"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailSizeAudits (
        auditId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        rawByteSize INTEGER NOT NULL,
        minifiedByteSize INTEGER NOT NULL,
        isUnderGmailLimit INTEGER NOT NULL DEFAULT 1,
        isClipped INTEGER NOT NULL DEFAULT 0,
        hasBase64Assets INTEGER NOT NULL DEFAULT 0,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Extends `pkg/email/size_guard.go` invoked in automated pre-send validation hooks.

---

### Improvement 445: Broken Link & Unsub Header RFC 8058 Compliance Verifier
- **ID:** `WARMUP-CANVAS-445`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Crawl all hyperlinks inside the email body and verify required `List-Unsubscribe` and `List-Unsubscribe-Post` headers to ensure 100% compliance with Google/Yahoo 2024 sender requirements.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Header validation summary:
      `[✔] List-Unsubscribe: <mailto:unsub@domain.com>, <https://domain.com/unsub>`
      `[✔] List-Unsubscribe-Post: List-Unsubscribe=One-Click`
    - Status table of all parsed hyperlinks with status codes (`200 OK`, `301 REDIR`, `404 NOT FOUND`).
  - **Zero-CDN Web UI:**
    - Link verification badges next to rendered anchor tags with one-click test dispatch.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit links`
  - Flags:
    - `--template <file>`: HTML template or raw `.eml` file.
    - `--check-remote-http`: Actively ping HTTP endpoints to confirm 200 OK responses.
    - `--timeout <duration>`: HTTP ping timeout per link (default: `3s`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubscribeComplianceReport struct {
        HasListUnsubscribeHeader bool     `json:"hasListUnsubscribeHeader"`
        HasOneClickPostHeader    bool     `json:"hasOneClickPostHeader"`
        IsRFC8058Compliant       bool     `json:"isRfc8058Compliant"`
        TotalLinksChecked        int      `json:"totalLinksChecked"`
        BrokenLinksCount         int      `json:"brokenLinksCount"`
        DeadURLList              []string `json:"deadUrlList"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailLinkAudits (
        auditId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        hasListUnsubscribeHeader INTEGER NOT NULL DEFAULT 1,
        hasOneClickPostHeader INTEGER NOT NULL DEFAULT 1,
        isRfc8058Compliant INTEGER NOT NULL DEFAULT 1,
        brokenLinksCount INTEGER NOT NULL DEFAULT 0,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/email/rfc8058_validator.go` with zero external dependencies using Go's `net/http` client.

---

### Improvement 446: Dynamic Merge Tag Syntax Linter & Missing Variable Safeguard
- **ID:** `WARMUP-CANVAS-446`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Parse template strings for merge tags (e.g., `{{FirstName}}`, `{{Company}}`, `{{UnsubscribeURL}}`), validating against provided JSON test fixtures and preventing fatal dispatch of raw curly brackets to recipient inboxes.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Code snippet view highlighting unmapped variables with crimson underline:
      `Hello {{FirstName}}, welcome to {{Organization}}`
      `      ^^^^^^^^^^^^^ Missing in fixture`
    - Summary card indicating variable resolution completeness (`4/5 Variables Bound`).
  - **Zero-CDN Web UI:**
    - Interactive fixture editor side-by-side with live template updating immediately as JSON keys are modified.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email lint tags`
  - Flags:
    - `--template <file>`: Template file containing merge tags.
    - `--fixture <file>`: JSON file providing key-value pairs.
    - `--fail-on-missing`: Cause CLI process to return exit code 1 if any unmapped tag is detected.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MergeTagLintReport struct {
        TotalTagsFound      int      `json:"totalTagsFound"`
        ResolvedTagsCount   int      `json:"resolvedTagsCount"`
        UnresolvedTagsList  []string `json:"unresolvedTagsList"`
        IsFullyResolved     bool     `json:"isFullyResolved"`
        HasMalformedSyntax  bool     `json:"hasMalformedSyntax"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailMergeTagLints (
        lintId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        totalTagsFound INTEGER NOT NULL,
        resolvedTagsCount INTEGER NOT NULL,
        isFullyResolved INTEGER NOT NULL DEFAULT 1,
        hasMalformedSyntax INTEGER NOT NULL DEFAULT 0,
        lintedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements AST scanner in `pkg/email/tag_linter.go` preventing embarrassing deliverability accidents.

---

### Improvement 447: Client-Specific CSS Fallback Generator for Outlook MSO Conditional Comments
- **ID:** `WARMUP-CANVAS-447`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Automatically analyze HTML button and table elements, injecting required Outlook MSO conditional comments (`<!--[if mso]> ... <![endif]-->`) and VML vector shapes for bulletproof rendering in legacy Microsoft Outlook desktop clients.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Diff viewer displaying before/after code blocks with green highlighted `<!--[if mso]>` injections.
    - Stats indicator: `Injected 3 VML button wrappers, 2 table width guards`.
  - **Zero-CDN Web UI:**
    - Side-by-side Outlook 2016 vs Webkit preview canvas toggles.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email inject mso-fallbacks`
  - Flags:
    - `--template <file>`: Input template path.
    - `--out-file <file>`: Destination path for normalized template.
    - `--wrap-buttons`: Specifically convert CSS border-radius buttons to VML shape tags.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MSOFallbackInjection struct {
        InjectedButtonCount int  `json:"injectedButtonCount"`
        InjectedTableGuards int  `json:"injectedTableGuards"`
        HasVMLDefinitions   bool `json:"hasVmlDefinitions"`
        IsOutlookCompatible bool `json:"isOutlookCompatible"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailMSOTransformations (
        transformId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        injectedButtonCount INTEGER NOT NULL DEFAULT 0,
        injectedTableGuards INTEGER NOT NULL DEFAULT 0,
        isOutlookCompatible INTEGER NOT NULL DEFAULT 1,
        transformedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/email/mso_transformer.go` with unit test fixtures covering Outlook 2013-2021 quirks.

---

### Improvement 448: Spam Trigger Phrase & Lexical Sentiment Terminal Scanner
- **ID:** `WARMUP-CANVAS-448`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Inspect template copy, subject lines, and preheaders against an updated database of over 800 spam trigger keywords (e.g., "100% free", "urgent transfer", excessive dollar signs `$$$`, ALL-CAPS words) and output a weighted Spam Score.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Score badge with gradient color scale: `Spam Risk Score: 1.4 / 10.0 [LOW RISK]`.
    - Lexical findings table itemizing trigger words, occurrence count, location (Subject, Body, Preheader), and severity weight.
  - **Zero-CDN Web UI:**
    - Visual text highlight overlay marking trigger words with yellow/red underlines and tooltips offering alternative synonyms.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit spam-score`
  - Flags:
    - `--template <file>`: HTML or plain text file.
    - `--subject <string>`: Email subject line to evaluate.
    - `--threshold <float>`: Maximum allowed spam score before exiting with failure (default: `4.5`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamScoreFinding struct {
        WordOrPhrase string  `json:"wordOrPhrase"`
        Location     string  `json:"location"`
        Weight       float64 `json:"weight"`
    }

    type SpamScoreAuditReport struct {
        CalculatedScore float64            `json:"calculatedScore"`
        MaxAllowedScore float64            `json:"maxAllowedScore"`
        IsUnderCeiling  bool               `json:"isUnderCeiling"`
        HasAllCapsTitle bool               `json:"hasAllCapsTitle"`
        Findings        []SpamScoreFinding `json:"findings"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailSpamScoreAudits (
        auditId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        calculatedScore REAL NOT NULL,
        isUnderCeiling INTEGER NOT NULL DEFAULT 1,
        hasAllCapsTitle INTEGER NOT NULL DEFAULT 0,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Powers `pkg/email/spam_scorer.go` acting as a pre-flight deliverability gate in warmup pipelines.

---

### Improvement 449: Plain-Text Alternative (multipart/alternative) Auto-Generator
- **ID:** `WARMUP-CANVAS-449`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Autonomously generate a clean, readable RFC 2046 `text/plain` alternative part from rich HTML templates, preserving layout hierarchy, list bulleting, and converting hyperlink anchors to bracketed plain-text URLs.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Side-by-side terminal comparison view: Left pane shows HTML structure preview, right pane shows generated plain-text output.
    - Character encoding confirmation badge (`UTF-8 Validated`, `No Raw CRLF Orphans`).
  - **Zero-CDN Web UI:**
    - Tab switcher on the preview interface between "Rich HTML" and "Generated Plain Text".
- **CLI Command Signature & Flags:**
  - Command: `gitmap email generate text-part`
  - Flags:
    - `--html-file <file>`: Source HTML email file.
    - `--out-file <file>`: Destination file for generated text part.
    - `--wrap-column <int>`: Word wrap column limit (default: `78` per RFC 5322).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PlainTextGenerationResult struct {
        SourceHTMLPath     string `json:"sourceHtmlPath"`
        GeneratedTextSize  int    `json:"generatedTextSize"`
        ExtractedLinkCount int    `json:"extractedLinkCount"`
        IsUTF8Compliant    bool   `json:"isUtf8Compliant"`
        HasUnwrappedLines  bool   `json:"hasUnwrappedLines"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailTextAlternatives (
        textId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        textContent TEXT NOT NULL,
        extractedLinkCount INTEGER NOT NULL DEFAULT 0,
        isUtf8Compliant INTEGER NOT NULL DEFAULT 1,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Embedded into `pkg/email/plain_text_builder.go` using HTML tokenizers without external markdown dependencies.

---

### Improvement 450: Zero-CDN Responsive Email Screenshot Capture & Local Diff Engine
- **ID:** `WARMUP-CANVAS-450`
- **Cluster:** 45 (Zero-CDN Client Email Canvas Live Preview & Dark Mode Contrast Auditing)
- **Goal:** Capture pixel-perfect headless browser screenshots across mobile and desktop breakpoints locally via Chrome DevTools protocol, calculating perceptual image diff hashes (pHash) against reference designs.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Status reporting progress bar: `Capturing Local Renders: [Desktop ✔] [Tablet ✔] [Mobile ✔]`.
    - Perceptual diff metric: `Visual Drift Score: 0.02% (Below 0.5% tolerance threshold - PASS)`.
  - **Zero-CDN Web UI:**
    - Interactive split-slider image comparison view (swipe left/right to reveal diff between base and modified version).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email canvas diff`
  - Flags:
    - `--template <file>`: Path to current template.
    - `--reference-image <file>`: Baseline image for perceptual comparison.
    - `--tolerance <float>`: Max acceptable pixel mismatch percentage (default: `0.05`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type VisualDiffAudit struct {
        TemplateID       string  `json:"templateId"`
        DiffPercentage   float64 `json:"diffPercentage"`
        IsMatchWithinTol bool    `json:"isMatchWithinTol"`
        HasVisualDrift   bool    `json:"hasVisualDrift"`
        SnapshotPath     string  `json:"snapshotPath"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EmailVisualDiffAudits (
        diffId TEXT PRIMARY KEY,
        templateId TEXT NOT NULL,
        diffPercentage REAL NOT NULL,
        isMatchWithinTol INTEGER NOT NULL DEFAULT 1,
        hasVisualDrift INTEGER NOT NULL DEFAULT 0,
        snapshotPath TEXT NOT NULL,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates with local Chrome headless runners in `pkg/email/visual_diff.go`.

---

# Cluster 46: Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector (Items 451–460)

---

### Improvement 451: Distributed Step-by-Step SMTP Handshake Packet & Command Tracer
- **ID:** `WARMUP-SEC-451`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Intercept and log every raw command/response packet exchanged during the SMTP conversation across distributed sender nodes, providing instant millisecond latency timestamps and reply code parsing.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Colorized interactive conversational trace:
      `>>> 220 mx.google.com ESMTP (0.0ms)` (Green)
      `<<< EHLO mail01.sender.org (1.2ms)` (Cyan)
      `>>> 250-mx.google.com at your service (14.5ms)` (Green)
      `<<< STARTTLS (15.1ms)` (Cyan)
      `>>> 220 2.0.0 Ready to start TLS (18.2ms)` (Green)
    - Bordered summary box highlighting total handshake duration and negotiate cipher.
  - **Zero-CDN Web UI:**
    - Live-scrolling terminal console widget with replay and copy-to-clipboard functionality.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email test handshake`
  - Flags:
    - `--mx-host <string>`: Destination mail exchanger hostname (e.g., `gmail-smtp-in.l.google.com`).
    - `--port <int>`: Remote SMTP port (default: `25`, supports `587`, `465`).
    - `--helo-domain <string>`: FQDN hostname to announce in EHLO greeting.
    - `--trace-raw`: Print full unmasked network byte stream.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SMTPTraceStep struct {
        Direction     string  `json:"direction"` // INBOUND or OUTBOUND
        Payload       string  `json:"payload"`
        ElapsedMS     float64 `json:"elapsedMs"`
        StatusCode    int     `json:"statusCode"`
        IsSuccessCode bool    `json:"isSuccessCode"`
        HasWarning    bool    `json:"hasWarning"`
    }

    type SMTPHandshakeReport struct {
        TargetMXHost    string          `json:"targetMxHost"`
        Port            int             `json:"port"`
        TotalDurationMS float64         `json:"totalDurationMs"`
        IsCompleted     bool            `json:"isCompleted"`
        HasTLSSupport   bool            `json:"hasTlsSupport"`
        Steps           []SMTPTraceStep `json:"steps"`
    }
    ```
  - Split-DB SQLite Schema (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpHandshakeTraces (
        traceId TEXT PRIMARY KEY,
        targetMxHost TEXT NOT NULL,
        port INTEGER NOT NULL DEFAULT 25,
        totalDurationMs REAL NOT NULL,
        isCompleted INTEGER NOT NULL DEFAULT 1,
        hasTlsSupport INTEGER NOT NULL DEFAULT 1,
        recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_smtp_trace_host ON SmtpHandshakeTraces(targetMxHost);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements custom `net.Conn` proxy wrapper in `pkg/smtp/tracer.go` instrumenting read/write buffers without degrading socket throughput.

---

### Improvement 452: TLS 1.3 Cipher Suite Negotiator & Forward Secrecy Inspector
- **ID:** `WARMUP-SEC-452`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Inspect negotiated TLS versions (TLS 1.2 vs TLS 1.3), active cipher suites, and Elliptic Curve Diffie-Hellman (ECDHE) parameters to verify Perfect Forward Secrecy (PFS) and reject obsolete/weak algorithms.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Security scorecard badge: `TLS Security Grade: A+ [TLS_AES_256_GCM_SHA384 | X25519 | PFS OK]`.
    - Warning indicator in bold amber if legacy TLS 1.0/1.1 is allowed or if RC4/3DES/CBC ciphers are encountered.
  - **Zero-CDN Web UI:**
    - Key exchange visual diagram displaying client hello offered suites vs server selected suite.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email test tls-security`
  - Flags:
    - `--mx-host <string>`: Target mail exchanger host.
    - `--min-version <string>`: Minimum acceptable TLS version (`1.2` or `1.3`, default: `1.2`).
    - `--enforce-pfs`: Require Perfect Forward Secrecy for successful exit status.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type TLSCipherSecurityReport struct {
        Host                string `json:"host"`
        TLSVersionName      string `json:"tlsVersionName"`
        CipherSuiteName     string `json:"cipherSuiteName"`
        KeyExchangeGroup    string `json:"keyExchangeGroup"`
        IsTLS13             bool   `json:"isTls13"`
        HasForwardSecrecy   bool   `json:"hasForwardSecrecy"`
        HasInsecureCiphers  bool   `json:"hasInsecureCiphers"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpTlsSecurityAudits (
        auditId TEXT PRIMARY KEY,
        host TEXT NOT NULL,
        tlsVersionName TEXT NOT NULL,
        cipherSuiteName TEXT NOT NULL,
        isTls13 INTEGER NOT NULL DEFAULT 1,
        hasForwardSecrecy INTEGER NOT NULL DEFAULT 1,
        hasInsecureCiphers INTEGER NOT NULL DEFAULT 0,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Connects to `pkg/smtp/tls_inspector.go` utilizing standard `crypto/tls` connection state metadata.

---

### Improvement 453: MTA-STS (RFC 8461) Policy Validation Engine & TLSRPT Listener
- **ID:** `WARMUP-SEC-453`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Autonomously query and validate Mail Transfer Agent Strict Transport Security (MTA-STS) TXT records and HTTPS policy endpoints (`mta-sts.<domain>/.well-known/mta-sts.txt`), ensuring strict enforcement against downgrade attacks.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Checklist output:
      `[✔] DNS TXT: _mta-sts.example.com -> v=STSv1; id=20261027T01;`
      `[✔] HTTPS Policy: https://mta-sts.example.com/.well-known/mta-sts.txt (200 OK)`
      `[✔] Mode: enforce | Max-Age: 604800s (7 days)`
    - Green checkmarks for valid MX match declarations.
  - **Zero-CDN Web UI:**
    - Status dashboard widget highlighting policy expiration countdown and testing MX alignment.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit mta-sts`
  - Flags:
    - `--domain <string>`: Domain to audit.
    - `--expected-mode <string>`: Expected mode policy (`testing` or `enforce`, default: `enforce`).
    - `--check-tlsrpt`: Also verify `_smtp._tls.<domain>` TLSRPT reporting record.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MTASTSPolicyReport struct {
        Domain         string   `json:"domain"`
        PolicyMode     string   `json:"policyMode"`
        MaxAgeSeconds  int      `json:"maxAgeSeconds"`
        MXPatterns     []string `json:"mxPatterns"`
        IsEnforced     bool     `json:"isEnforced"`
        HasValidHTTPS  bool     `json:"hasValidHttps"`
        HasTLSRPT      bool     `json:"hasTlsrpt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpMTASTSPolicies (
        policyId TEXT PRIMARY KEY,
        domain TEXT NOT NULL,
        policyMode TEXT NOT NULL,
        maxAgeSeconds INTEGER NOT NULL,
        isEnforced INTEGER NOT NULL DEFAULT 1,
        hasValidHttps INTEGER NOT NULL DEFAULT 1,
        hasTlsrpt INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX IF NOT EXISTS idx_mta_sts_domain ON SmtpMTASTSPolicies(domain);
    ```
- **Antigravity / GitMap Integration Point:**
  - Handled by `pkg/smtp/mta_sts.go` caching policies per RFC specifications to minimize external HTTPS requests.

---

### Improvement 454: DANE TLSA (RFC 6698) DNSSEC Authentication Validator
- **ID:** `WARMUP-SEC-454`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Validate DNS-based Authentication of Named Entities (DANE) by querying TLSA records (`_25._tcp.mx.<domain>`) over DNSSEC-validated resolvers, ensuring certificate pinning integrity.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Terminal security badge: `DANE/TLSA Status: SECURE [DNSSEC Validated | Match Type: SPKI SHA-256]`.
    - Warning alert if DNSSEC validation fails or AD (Authenticated Data) flag is absent in DNS header.
  - **Zero-CDN Web UI:**
    - Certificate pin comparison tree illustrating public key hash matches.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit dane`
  - Flags:
    - `--mx-host <string>`: Target mail exchanger host.
    - `--dnssec-resolver <ip>`: DNSSEC-validating DNS server (default: `1.1.1.1:53` or system resolver).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DANETLSAReport struct {
        MXHost             string `json:"mxHost"`
        UsageField         int    `json:"usageField"`
        SelectorField      int    `json:"selectorField"`
        MatchingType       int    `json:"matchingType"`
        CertificateData    string `json:"certificateData"`
        IsDNSSECValidated  bool   `json:"isDnssecValidated"`
        IsCertMatchingTLSA bool   `json:"isCertMatchingTlsa"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpDANEAudits (
        auditId TEXT PRIMARY KEY,
        mxHost TEXT NOT NULL,
        isDnssecValidated INTEGER NOT NULL DEFAULT 1,
        isCertMatchingTlsa INTEGER NOT NULL DEFAULT 1,
        auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements DNSSEC resolver calls in `pkg/smtp/dane.go`.

---

### Improvement 455: Certificate Chain Expiration, SAN & Revocation (OCSP/CRL) Guard
- **ID:** `WARMUP-SEC-455`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Verify the full X.509 certificate chain of remote MX servers, asserting Common Name (CN) and Subject Alternative Names (SAN) match the connected host, verifying expiration countdowns, and querying OCSP stapling status.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Cert status card:
      `ISSUER: Let's Encrypt Authority X3`
      `SUBJECT: mx1.cleanmail.org`
      `SAN: *.cleanmail.org, cleanmail.org, mx1.cleanmail.org`
      `VALIDITY: 64 days remaining (Expires: 2026-12-30)`
      `OCSP STATUS: Stapled (Good)`
  - **Zero-CDN Web UI:**
    - Visual certificate expiration timeline with colored urgency badges.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email audit cert-chain`
  - Flags:
    - `--host <string>`: Remote server hostname.
    - `--port <int>`: Port (default: `25`).
    - `--warn-days <int>`: Days threshold to trigger warning (default: `30`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CertChainAuditResult struct {
        CommonName          string    `json:"commonName"`
        IssuerOrganization  string    `json:"issuerOrganization"`
        DaysUntilExpiration int       `json:"daysUntilExpiration"`
        IsSANMatching       bool      `json:"isSanMatching"`
        IsOCSPStapled       bool      `json:"isOcspStapled"`
        HasRevocationAlert  bool      `json:"hasRevocationAlert"`
        ExpiresAt           time.Time `json:"expiresAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpCertAudits (
        certId TEXT PRIMARY KEY,
        host TEXT NOT NULL,
        commonName TEXT NOT NULL,
        daysUntilExpiration INTEGER NOT NULL,
        isSanMatching INTEGER NOT NULL DEFAULT 1,
        isOcspStapled INTEGER NOT NULL DEFAULT 1,
        hasRevocationAlert INTEGER NOT NULL DEFAULT 0,
        expiresAt DATETIME NOT NULL,
        checkedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Part of `pkg/smtp/cert_guard.go` used during proactive relay health checks.

---

### Improvement 456: Real-Time SMTP Banner Greeter & Capabilities (EHLO Keywords) Analyzer
- **ID:** `WARMUP-SEC-456`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Parse 220 greeting banners and 250 EHLO responses to catalog all supported SMTP extensions (`SIZE`, `8BITMIME`, `STARTTLS`, `ENHANCEDSTATUSCODES`, `CHUNKING`, `PIPELINING`, `SMTPUTF8`, `DSN`).
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Extensions grid with status badges:
      `[✔] PIPELINING (RFC 2920)   | [✔] 8BITMIME (RFC 6152)`
      `[✔] SIZE 104857600 (100MB) | [✔] CHUNKING (RFC 3030)`
      `[✔] STARTTLS (RFC 3207)    | [✖] SMTPUTF8 (Not Supported)`
  - **Zero-CDN Web UI:**
    - Interactive matrix displaying capabilities of target mail providers side-by-side.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email test capabilities`
  - Flags:
    - `--mx-host <string>`: Destination server hostname.
    - `--require <string>`: Comma-separated list of required extensions (e.g., `STARTTLS,PIPELINING`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SMTPCapabilitiesSummary struct {
        BannerGreeting   string   `json:"bannerGreeting"`
        MaxMessageBytes  int64    `json:"maxMessageBytes"`
        SupportsTLS      bool     `json:"supportsTls"`
        SupportsPipeline bool     `json:"supportsPipeline"`
        SupportsChunking bool     `json:"supportsChunking"`
        SupportsUTF8     bool     `json:"supportsUtf8"`
        HasAuthExtension bool     `json:"hasAuthExtension"`
        RawExtensions    []string `json:"rawExtensions"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpCapabilitiesCache (
        host TEXT PRIMARY KEY,
        bannerGreeting TEXT NOT NULL,
        maxMessageBytes INTEGER NOT NULL DEFAULT 0,
        supportsTls INTEGER NOT NULL DEFAULT 1,
        supportsPipeline INTEGER NOT NULL DEFAULT 1,
        supportsChunking INTEGER NOT NULL DEFAULT 0,
        supportsUtf8 INTEGER NOT NULL DEFAULT 1,
        hasAuthExtension INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Informs `pkg/smtp/envelope_sender.go` whether pipelining can be leveraged to accelerate message injection.

---

### Improvement 457: SMTP Smuggling (RFC 5321 End-of-Data Desync) Vulnerability Probe
- **ID:** `WARMUP-SEC-457`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Safely probe outbound relays and inbound gateways for SMTP smuggling vulnerabilities caused by desynchronized `<CR><LF>.<CR><LF>` vs `<LF>.<LF>` handling between transit relays.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Security audit alert:
      `[SEC-PROBE] Testing SMTP Smuggling Ingestion Patterns...`
      `[✔] Non-standard LF.LF: REJECTED (500 5.5.2 Syntax error)`
      `[✔] Pipelined Carriage Return: STRIPPED`
      `VULNERABILITY VERDICT: IMMUNE [SAFE]`
  - **Zero-CDN Web UI:**
    - Explanatory packet sequence animation showing how intermediary proxies sanitize malformed line feeds.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email test smuggling-probe`
  - Flags:
    - `--relay-host <string>`: Outbound relay host to test.
    - `--test-mode <string>`: Probe vector (`bare-lf`, `cr-only`, `pipelined-eod`, `all`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SMTPSmugglingAuditReport struct {
        TargetRelay       string `json:"targetRelay"`
        IsImmune          bool   `json:"isImmune"`
        HasBareLFAccepted bool   `json:"hasBareLfAccepted"`
        HasCROnlyAccepted bool   `json:"hasCrOnlyAccepted"`
        RiskLevel         string `json:"riskLevel"` // SAFE, ELEVATED, CRITICAL
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpSmugglingAudits (
        probeId TEXT PRIMARY KEY,
        targetRelay TEXT NOT NULL,
        isImmune INTEGER NOT NULL DEFAULT 1,
        hasBareLfAccepted INTEGER NOT NULL DEFAULT 0,
        riskLevel TEXT NOT NULL DEFAULT 'SAFE',
        probedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Security inspection tool located in `pkg/smtp/smuggling_probe.go` adhering to CERT advisory standards.

---

### Improvement 458: Distributed Latency Anomaly & Geo-Route BGP Route Inspector
- **ID:** `WARMUP-SEC-458`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Trace network latency and route hops to receiving MX servers across distributed geographic swarm nodes, identifying sub-optimal routing paths, packet jitter, and ASN transit anomalies.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Terminal route hop table:
      `HOP  IP               ASN       RTT     GEO`
      `1    192.0.2.1        AS13335   0.4ms   US-ORD`
      `2    198.51.100.5     AS15169   8.2ms   US-IAD`
      `3    142.250.100.26   AS15169   12.1ms  US-IAD (Google MX)`
  - **Zero-CDN Web UI:**
    - World map SVG projection marking node locations and connection hops with latency labels.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email test route-tracer`
  - Flags:
    - `--destination-mx <string>`: Target MX IP or hostname.
    - `--from-node <string>`: Specific warm-up swarm node ID.
    - `--max-hops <int>`: Max traceroute hops (default: `15`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type RouteHopRecord struct {
        HopIndex  int     `json:"hopIndex"`
        IPAddress string  `json:"ipAddress"`
        ASNNumber int     `json:"asnNumber"`
        RTTMS     float64 `json:"rttMs"`
        GeoRegion string  `json:"geoRegion"`
    }

    type GeoRouteTraceReport struct {
        OriginNodeID  string           `json:"originNodeId"`
        TargetMX      string           `json:"targetMx"`
        TotalHops     int              `json:"totalHops"`
        AverageRTTMS  float64          `json:"averageRttMs"`
        IsOptimalPath bool             `json:"isOptimalPath"`
        HasHighJitter bool             `json:"hasHighJitter"`
        Hops          []RouteHopRecord `json:"hops"`
    }
    ```
  - Split-DB SQLite Schema (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpGeoRouteTraces (
        traceId TEXT PRIMARY KEY,
        originNodeId TEXT NOT NULL,
        targetMx TEXT NOT NULL,
        totalHops INTEGER NOT NULL,
        averageRttMs REAL NOT NULL,
        isOptimalPath INTEGER NOT NULL DEFAULT 1,
        hasHighJitter INTEGER NOT NULL DEFAULT 0,
        tracedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Implements lightweight UDP/ICMP echo tracer in `pkg/net/route_tracer.go`.

---

### Improvement 459: Automated DSN Bounce Classification & Enhanced Status Code (RFC 3463) Parser
- **ID:** `WARMUP-SEC-459`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Ingest Delivery Status Notifications (DSNs) and parse Enhanced Mail System Status Codes (e.g., `5.7.1` Security/Policy, `4.4.1` No Answer from Host, `5.2.2` Mailbox Full), categorizing bounces into hard, soft, or transient rate limits.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - High-visibility diagnostic card:
      `BOUNCE TYPE: Hard Policy Rejection [5.7.26]`
      `SUMMARY: SPF/DKIM Authentication failed for sender domain`
      `ACTION: Automatic circuit suppression applied to recipient`
  - **Zero-CDN Web UI:**
    - Donut chart showing bounce distribution categorized by RFC sub-codes.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email parse dsn`
  - Flags:
    - `--raw-dsn-file <file>`: File containing raw bounced MIME message.
    - `--json`: Output parsed classification structure.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EnhancedDSNClassification struct {
        RawCode         string `json:"rawCode"` // e.g., "5.7.26"
        ClassCode       int    `json:"classCode"` // 2, 4, 5
        SubjectCategory string `json:"subjectCategory"`
        DetailReason    string `json:"detailReason"`
        IsHardBounce    bool   `json:"isHardBounce"`
        IsPolicyReject  bool   `json:"isPolicyReject"`
        HasSuppressed   bool   `json:"hasSuppressed"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpDSNBounceEvents (
        bounceId TEXT PRIMARY KEY,
        rawCode TEXT NOT NULL,
        classCode INTEGER NOT NULL,
        subjectCategory TEXT NOT NULL,
        isHardBounce INTEGER NOT NULL DEFAULT 1,
        isPolicyReject INTEGER NOT NULL DEFAULT 0,
        hasSuppressed INTEGER NOT NULL DEFAULT 1,
        recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_smtp_dsn_code ON SmtpDSNBounceEvents(rawCode);
    ```
- **Antigravity / GitMap Integration Point:**
  - Core parser located in `pkg/smtp/dsn_parser.go` adhering strictly to RFC 3463 and RFC 1893 specifications.

---

### Improvement 460: Dynamic MX Load Balancing & Fallback Circuit Health Prober
- **ID:** `WARMUP-SEC-460`
- **Cluster:** 46 (Distributed SMTP Handshake Tracer & TLS Cipher Security Inspector)
- **Goal:** Concurrently probe all published MX records for a target recipient domain, measuring handshake latencies and connection availability across primary and secondary MX servers, automatically falling back to lower-priority MX hosts when primary hosts reject connections.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal UX:**
    - Priority table:
      `PREF  MX HOST                       STATUS       RTT     HEALTH`
      `5     alt1.aspmx.l.google.com       ONLINE       14ms    [█████] 100%`
      `10    alt2.aspmx.l.google.com       ONLINE       19ms    [█████] 100%`
      `20    aspmx2.googlemail.com         ONLINE       28ms    [████░]  80%`
      `30    aspmx3.googlemail.com         DEGRADED     120ms   [██░░░]  40%`
  - **Zero-CDN Web UI:**
    - Visual connection topology tree highlighting chosen path and fallback alternates.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email test mx-pool`
  - Flags:
    - `--domain <string>`: Recipient domain to probe (e.g., `google.com`).
    - `--parallel-probes`: Probe all resolved MX records simultaneously.
    - `--timeout <duration>`: Probe timeout per MX target (default: `2s`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MXHostHealthSummary struct {
        Host            string  `json:"host"`
        Preference      int     `json:"preference"`
        LatencyMS       float64 `json:"latencyMs"`
        IsOnline        bool    `json:"isOnline"`
        IsPrimary       bool    `json:"isPrimary"`
        HasConnectionErr bool   `json:"hasConnectionErr"`
    }

    type DomainMXPoolReport struct {
        Domain          string                `json:"domain"`
        SelectedMX      string                `json:"selectedMx"`
        TotalMXCount    int                   `json:"totalMxCount"`
        IsFullyHealthy  bool                  `json:"isFullyHealthy"`
        HasFallbackUsed bool                  `json:"hasFallbackUsed"`
        Hosts           []MXHostHealthSummary `json:"hosts"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SmtpDomainMXHealth (
        healthId TEXT PRIMARY KEY,
        domain TEXT NOT NULL,
        selectedMx TEXT NOT NULL,
        totalMxCount INTEGER NOT NULL,
        isFullyHealthy INTEGER NOT NULL DEFAULT 1,
        hasFallbackUsed INTEGER NOT NULL DEFAULT 0,
        probedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_domain_mx ON SmtpDomainMXHealth(domain);
    ```
- **Antigravity / GitMap Integration Point:**
  - Powers `pkg/smtp/mx_pool.go` to provide intelligent routing and fault isolation across multi-homed mail delivery targets.
