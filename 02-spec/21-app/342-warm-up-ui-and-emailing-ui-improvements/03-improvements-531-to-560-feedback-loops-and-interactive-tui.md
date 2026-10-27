# Specification 342: Improvements 531 to 560 (Recipient Engagement Feedback Loops, Bubble Tea Mailbox Inspector & Memory-Mapped Backpressure HUD)

- **Sequence:** 342
- **Slug:** `342-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/03-improvements-531-to-560-feedback-loops-and-interactive-tui.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (531 - 560)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 54** | Items 531–540 | Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines |
| **Cluster 55** | Items 541–550 | Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette |
| **Cluster 56** | Items 551–560 | Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer |

---

# Cluster 54: Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines (Items 531–540)

---

### Improvement 531: Streaming Open/Click Micro-Sparkline Terminal Telemetry Grid
- **ID:** `ENGAGE-LOOP-531`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Render live-updating Unicode micro-sparklines (` ▂▃▄▅▆▇█`) tracking sub-minute recipient interaction events (pixel loads, link click-throughs, read durations) across active warm-up recipient pools.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal View:**
    - Dual horizontal bar layout displaying rolling 60-second open/click frequency per domain cluster (`gmail.com`, `yahoo.com`, `outlook.com`).
    - Lipgloss token palette: Opens (`#10b981`, emerald ` ▂▃▄▅▆▇█`), Clicks (`#06b6d4`, cyan ` ▂▃▄▅▆▇█`), Fast-Unsubscribes (`#ef4444`, red `▲`).
    - Model view snippet:
      ```text
      ┌─ [ENGAGEMENT FEEDBACK STREAM: GMAIL POOL] ───────────────────────────────────────────┐
      │ OPENS/MIN :  ▂▃▅▆▇█▇▆▅▃▂   142 ops/m [Avg Read: 18.4s]  [Active Recipients: 1,480]   │
      │ CLICKS/MIN:   ▂ ▃ ▄ ▃ ▂     38 clk/m [CTR: 26.7%]       [Health Status: OPTIMAL]     │
      └──────────────────────────────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Dashboard:**
    - Pure inline SVG polylines and embedded CSS styling with 1-second WebSocket push.
    - Zero external charting runtime; pure DOM SVG `<polyline>` coordinate mapping.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup sparklines`
  - Flags:
    - `--pool-id <string>`: Target warmup identity or campaign pool (default: `global`).
    - `--window-duration <string>`: Window rolling interval, e.g. `60s`, `5m`, `15m` (default: `60s`).
    - `--metric-type <string>`: Filter metric (`opens`, `clicks`, `all`, `unsubscribes`) (default: `all`).
    - `--watch`: Persistent interactive Bubble Tea full-screen alt-buffer mode.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EngagementSparklinePoint struct {
        Timestamp      time.Time `json:"timestamp"`
        OpenCount      int       `json:"openCount"`
        ClickCount     int       `json:"clickCount"`
        UnsubCount     int       `json:"unsubCount"`
        IsOpenSpike    bool      `json:"isOpenSpike"`
        HasClickDrop   bool      `json:"hasClickDrop"`
    }

    type EngagementStreamState struct {
        PoolID         string                   `json:"poolId"`
        DomainCluster  string                   `json:"domainCluster"`
        Points         []EngagementSparklinePoint `json:"points"`
        AverageReadSec float64                  `json:"averageReadSec"`
        IsOptimal      bool                     `json:"isOptimal"`
        HasHighCTR     bool                     `json:"hasHighCtr"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EngagementSparklines (
        sparklineId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        domainCluster TEXT NOT NULL,
        openCount INTEGER NOT NULL DEFAULT 0,
        clickCount INTEGER NOT NULL DEFAULT 0,
        unsubCount INTEGER NOT NULL DEFAULT 0,
        readDurationSec REAL NOT NULL DEFAULT 0.0,
        isOpenSpike INTEGER NOT NULL DEFAULT 0,
        hasClickDrop INTEGER NOT NULL DEFAULT 0,
        isOptimal INTEGER NOT NULL DEFAULT 1,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_sparkline_pool_time ON EngagementSparklines(poolId, createdAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Direct integration into `cli/cmdemail/sparkline_hud.go` binding to `pkg/warmup/metrics` event multiplexer.

---

### Improvement 532: Recipient Read-Duration Quantile Histogram Terminal Visualizer
- **ID:** `ENGAGE-LOOP-532`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Display real-time recipient dwell time distribution (`<2s Glanced`, `2-8s Skimmed`, `>8s Read`) with ASCII segmented bar charts to detect ESP bot pre-fetching vs genuine user engagement.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Terminal Model:**
    - 3-segment color-coded horizontal gauge: Glanced (`#64748b`, muted slate), Skimmed (`#f59e0b`, amber), Thorough Read (`#22c55e`, vibrant green).
    - Bot-heuristic indicator flag highlighting anomalous `<0.5s` opens originating from Google Proxy or Apple Mail Privacy Protection IPs.
  - **Zero-CDN Web View:**
    - CSS Flexbox meter bar with tooltips showing exact recipient counts and bot attribution percentages.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup read-dwell`
  - Flags:
    - `--pool-id <string>`: Warmup sender pool identifier.
    - `--exclude-proxy-bots`: Strip identified proxy IP fetches from human engagement stats (default: `true`).
    - `--detail-level <string>`: Output verbosity (`summary`, `histogram`, `raw`) (default: `histogram`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DwellTimeBucket struct {
        BucketName   string  `json:"bucketName"`
        MinSeconds   float64 `json:"minSeconds"`
        MaxSeconds   float64 `json:"maxSeconds"`
        SampleCount  int     `json:"sampleCount"`
        Percentage   float64 `json:"percentage"`
        IsHumanGrade bool    `json:"isHumanGrade"`
        HasBotMarker bool    `json:"hasBotMarker"`
    }

    type DwellReport struct {
        PoolID       string            `json:"poolId"`
        TotalOpens   int               `json:"totalOpens"`
        Buckets      []DwellTimeBucket `json:"buckets"`
        IsAggregated bool              `json:"isAggregated"`
        HasBotBias   bool              `json:"hasBotBias"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupReadDwell (
        dwellId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        bucketName TEXT NOT NULL,
        minSeconds REAL NOT NULL,
        maxSeconds REAL NOT NULL,
        sampleCount INTEGER NOT NULL DEFAULT 0,
        percentage REAL NOT NULL DEFAULT 0.0,
        isHumanGrade INTEGER NOT NULL DEFAULT 0,
        hasBotMarker INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_dwell_pool ON WarmupReadDwell(poolId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs into `pkg/analytics/dwell_detector.go` and exposes real-time updates through `gitmap email warmup read-dwell`.

---

### Improvement 533: Adaptive Velocity Damper on Dwell-Time Anomaly Detection
- **ID:** `ENGAGE-LOOP-533`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Automatically damp dispatch velocity down by 25% to 50% when engagement metrics detect sudden drops in average read duration, indicating spam folder placement or recipient disengagement.
- **UI/UX Terminal & Web Specification:**
  - **Terminal Alert Banner:**
    - Flashing yellow border: `lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("#eab308"))`.
    - Real-time adjustment readout: `[DAMPER ACTIVATED] Dwell Drop -42% -> Rate Capped at 15 msgs/hr/IP`.
  - **Web Notice Toast:**
    - Sticky banner with one-click manual override button and clear metric delta breakdown.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup velocity damper`
  - Flags:
    - `--pool-id <string>`: Warmup identity pool.
    - `--auto-throttle`: Enable or disable autonomous velocity reduction (default: `true`).
    - `--threshold-drop-pct <float>`: Trigger drop percentage (default: `35.0`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DamperEvent struct {
        EventID       string    `json:"eventId"`
        PoolID        string    `json:"poolId"`
        OriginalRate  int       `json:"originalRate"`
        DampedRate    int       `json:"dampedRate"`
        DropRatio     float64   `json:"dropRatio"`
        ReasonCode    string    `json:"reasonCode"`
        IsThrottled   bool      `json:"isThrottled"`
        HasOverride   bool      `json:"hasOverride"`
        CreatedAt     time.Time `json:"createdAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS VelocityDamperEvents (
        eventId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        originalRate INTEGER NOT NULL,
        dampedRate INTEGER NOT NULL,
        dropRatio REAL NOT NULL,
        reasonCode TEXT NOT NULL,
        isThrottled INTEGER NOT NULL DEFAULT 1,
        hasOverride INTEGER NOT NULL DEFAULT 0,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_damper_pool_time ON VelocityDamperEvents(poolId, createdAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Interfaces directly with `pkg/warmup/throttler.go` to enforce hardware-timed rate capping.

---

### Improvement 534: Domain-Clustered Click Heatmap ANSI Matrix
- **ID:** `ENGAGE-LOOP-534`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Generate a terminal-native 2D text heatmap plotting click distribution by hour-of-day vs recipient destination MX provider (Google, Microsoft 365, Yahoo, Fastmail, Proton).
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Grid Rendering:**
    - Dynamic 24-hour horizontal columns with 6 provider rows.
    - Color ramp: 0% (`#1e293b`), 1-25% (`#1e40af`), 26-50% (`#3b82f6`), 51-75% (`#60a5fa`), >75% (`#93c5fd`).
    - Row header highlights with clean Unicode box-drawing frames (`│`, `┼`, `─`).
  - **Web Heatmap Table:**
    - High-density HTML `<table>` styled with pure CSS inline grid and native tooltips.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup click-heatmap`
  - Flags:
    - `--days <int>`: Historical lookback duration in days (default: `7`).
    - `--provider <string>`: Filter by MX target cluster (`all`, `gmail`, `m365`, `yahoo`) (default: `all`).
    - `--format <string>`: Render format (`ansi`, `json`, `csv`) (default: `ansi`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HeatmapCell struct {
        HourOfDay  int     `json:"hourOfDay"`
        Provider   string  `json:"provider"`
        ClickCount int     `json:"clickCount"`
        Ratio      float64 `json:"ratio"`
        IsPeakHour bool    `json:"isPeakHour"`
        HasTraffic bool    `json:"hasTraffic"`
    }

    type HeatmapMatrix struct {
        ProviderList []string      `json:"providerList"`
        Cells        []HeatmapCell `json:"cells"`
        IsCompiled   bool          `json:"isCompiled"`
        HasOutliers  bool          `json:"hasOutliers"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS HeatmapMatrixCells (
        cellId TEXT PRIMARY KEY,
        provider TEXT NOT NULL,
        hourOfDay INTEGER NOT NULL,
        clickCount INTEGER NOT NULL DEFAULT 0,
        ratio REAL NOT NULL DEFAULT 0.0,
        isPeakHour INTEGER NOT NULL DEFAULT 0,
        hasTraffic INTEGER NOT NULL DEFAULT 0,
        recordedDate DATE NOT NULL
    );
    CREATE UNIQUE INDEX IF NOT EXISTS uq_heatmap_cell ON HeatmapMatrixCells(provider, hourOfDay, recordedDate);
    ```
- **Antigravity / GitMap Integration Point:**
  - Built inside `cli/cmdemail/heatmap_cmd.go` consuming aggregated metrics from `pkg/analytics/heatmap_engine.go`.

---

### Improvement 535: Autonomous Unsubscribe Velocity Watchdog & Auto-Quarantine HUD
- **ID:** `ENGAGE-LOOP-535`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Real-time watchdog monitoring rapid unsubscribe spikes (>0.5% in 15 minutes) with instant auto-quarantine of failing subject lines and terminal modal prompt for sender intervention.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Intercept Modal:**
    - Full-screen high-priority modal with flashing red header `[CRITICAL UNSUBSCRIBE SURGE DETECTED]`.
    - Key actions: `[q] Acknowledge & Quarantine Pool`, `[o] Override Damper`, `[i] Inspect Recipients`.
  - **Web Dashboard Alert:**
    - Red outline pulsing banner with auto-paused queue indicator.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup watchdog unsub`
  - Flags:
    - `--pool-id <string>`: Identifier of warmup pool.
    - `--spike-threshold <float>`: Percentage threshold triggering auto-quarantine (default: `0.5`).
    - `--auto-quarantine`: Automatically halt dispatches if threshold is exceeded (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnsubSpikeAlert struct {
        AlertID         string    `json:"alertId"`
        PoolID          string    `json:"poolId"`
        UnsubRate       float64   `json:"unsubRate"`
        ThresholdRate   float64   `json:"thresholdRate"`
        SuspectSubject  string    `json:"suspectSubject"`
        IsQuarantined   bool      `json:"isQuarantined"`
        HasOperatorNotif bool     `json:"hasOperatorNotif"`
        TriggeredAt     time.Time `json:"triggeredAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnsubWatchdogAlerts (
        alertId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        unsubRate REAL NOT NULL,
        thresholdRate REAL NOT NULL,
        suspectSubject TEXT NOT NULL,
        isQuarantined INTEGER NOT NULL DEFAULT 1,
        hasOperatorNotif INTEGER NOT NULL DEFAULT 0,
        triggeredAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_unsub_alerts_pool ON UnsubWatchdogAlerts(poolId, isQuarantined);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/warmup/watchdog.go` with direct signal dispatch to the Lemon background daemon.

---

### Improvement 536: Recipient Client Fingerprinting & Modern MIME Preview Matrix
- **ID:** `ENGAGE-LOOP-536`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Classify recipient User-Agent strings (Apple Mail, Outlook Desktop, Gmail Web, Thunderbird) and render client proportion breakdowns in terminal progress gauges.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Horizontal Percentage Bars:**
    - Proportional stacked bar: Apple Mail (`#a855f7`, purple), Outlook (`#0284c7`, azure), Gmail (`#ea4335`, red), Others (`#64748b`, gray).
    - Dynamic percentage labels formatted to one decimal place.
  - **Web Client Breakdown:**
    - Minimalist inline SVG ring chart (donut) with clean CSS transitions and zero external libraries.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup client-matrix`
  - Flags:
    - `--pool-id <string>`: Warmup identity pool.
    - `--timeframe <string>`: Analysis window, e.g. `24h`, `7d`, `30d` (default: `24h`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ClientShare struct {
        ClientName   string  `json:"clientName"`
        TotalOpens   int     `json:"totalOpens"`
        Percentage   float64 `json:"percentage"`
        IsDesktop    bool    `json:"isDesktop"`
        IsMobile     bool    `json:"isMobile"`
        HasDarkTheme bool    `json:"hasDarkTheme"`
    }

    type ClientMatrixState struct {
        PoolID      string        `json:"poolId"`
        Shares      []ClientShare `json:"shares"`
        IsFinalized bool          `json:"isFinalized"`
        HasData     bool          `json:"hasData"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupClientBreakdown (
        breakdownId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        clientName TEXT NOT NULL,
        totalOpens INTEGER NOT NULL DEFAULT 0,
        percentage REAL NOT NULL DEFAULT 0.0,
        isDesktop INTEGER NOT NULL DEFAULT 0,
        isMobile INTEGER NOT NULL DEFAULT 0,
        hasDarkTheme INTEGER NOT NULL DEFAULT 0,
        analyzedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_client_breakdown_pool ON WarmupClientBreakdown(poolId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Connects to `pkg/fingerprint/useragent_parser.go` for lightning-fast zero-allocation regex evaluation.

---

### Improvement 537: Link Click-Through Path Telemetry & Redirection Trace Tree
- **ID:** `ENGAGE-LOOP-537`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Map the exact redirection hops (custom domain tracking -> HTTPS edge redirector -> destination landing page) in a terminal ASCII tree view with latency breakdown per hop.
- **UI/UX Terminal & Web Specification:**
  - **Terminal Tree Layout:**
    - Branch glyphs (`├──`, `└──`) styled with `lipgloss.NewStyle().Foreground(lipgloss.Color("#8b5cf6"))`.
    - Displays SSL handshake time, DNS resolution latency, and HTTP status code per redirection hop.
  - **Web Tracing Component:**
    - Timeline view connecting nodes via CSS flex lines with latency badges.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup trace-link`
  - Flags:
    - `--url <string>`: Tracked email URL to simulate and inspect.
    - `--timeout-ms <int>`: Maximum tracing duration per hop (default: `3000`).
    - `--verify-ssl`: Strict TLS certificate verification (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type RedirectionHop struct {
        HopIndex     int           `json:"hopIndex"`
        URL          string        `json:"url"`
        StatusCode   int           `json:"statusCode"`
        Duration     time.Duration `json:"duration"`
        IsSecureTLS  bool          `json:"isSecureTls"`
        HasCertError bool          `json:"hasCertError"`
    }

    type LinkTraceReport struct {
        TargetURL    string           `json:"targetUrl"`
        FinalURL     string           `json:"finalUrl"`
        TotalTime    time.Duration    `json:"totalTime"`
        Hops         []RedirectionHop `json:"hops"`
        IsClean      bool             `json:"isClean"`
        HasLoops     bool             `json:"hasLoops"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS LinkTraceAudits (
        traceId TEXT PRIMARY KEY,
        targetUrl TEXT NOT NULL,
        finalUrl TEXT NOT NULL,
        totalTimeMs INTEGER NOT NULL,
        hopCount INTEGER NOT NULL,
        isClean INTEGER NOT NULL DEFAULT 1,
        hasLoops INTEGER NOT NULL DEFAULT 0,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_link_traces_time ON LinkTraceAudits(createdAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Leverages `pkg/nettrace/http_tracer.go` directly inside `cli/cmdemail/trace_link.go`.

---

### Improvement 538: Live Engagement Delta Ticker for Multi-Pool A/B Comparison
- **ID:** `ENGAGE-LOOP-538`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Side-by-side terminal ticker comparing real-time performance indicators (Open Rate, CTR, Spam Mark Rate) between Pool A (Baseline) and Pool B (Adaptive Test).
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Dual Column Compare:**
    - Left Column (Pool A): `lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(38)`.
    - Right Column (Pool B): `lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(38)`.
    - Relative deltas rendered in green (`▲ +3.4%`) or crimson (`▼ -1.8%`).
  - **Web Dual Comparison Card:**
    - Responsive two-column grid highlighting winning metrics with subtle glow.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup compare-pools`
  - Flags:
    - `--pool-a <string>`: Baseline warmup identity pool.
    - `--pool-b <string>`: Challenger warmup identity pool.
    - `--interval <string>`: Ticker refresh rate (default: `2s`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PoolDeltaMetric struct {
        MetricName  string  `json:"metricName"`
        ValuePoolA  float64 `json:"valuePoolA"`
        ValuePoolB  float64 `json:"valuePoolB"`
        DeltaRatio  float64 `json:"deltaRatio"`
        IsPositive  bool    `json:"isPositive"`
        HasStatSig  bool    `json:"hasStatSig"`
    }

    type PoolComparisonReport struct {
        PoolA       string            `json:"poolA"`
        PoolB       string            `json:"poolB"`
        Metrics     []PoolDeltaMetric `json:"metrics"`
        IsCompleted bool              `json:"isCompleted"`
        HasWinner   bool              `json:"hasWinner"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupPoolComparisons (
        comparisonId TEXT PRIMARY KEY,
        poolA TEXT NOT NULL,
        poolB TEXT NOT NULL,
        openRateDelta REAL NOT NULL,
        clickRateDelta REAL NOT NULL,
        isPositive INTEGER NOT NULL DEFAULT 1,
        hasStatSig INTEGER NOT NULL DEFAULT 0,
        hasWinner INTEGER NOT NULL DEFAULT 0,
        evaluatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_pool_comp ON WarmupPoolComparisons(poolA, poolB);
    ```
- **Antigravity / GitMap Integration Point:**
  - Exposed via `cli/cmdemail/compare_pools.go` linking with `pkg/stats/bayesian_eval.go`.

---

### Improvement 539: Sentiment & Keyword Triggered Reply Heatmap Generator
- **ID:** `ENGAGE-LOOP-539`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Parse incoming warmup peer responses and generate a semantic sentiment distribution (Positive, Neutral, Out-of-Office, Hostile/Unsubscribe) with ANSI color indicators.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Sentiment Card:**
    - Live list of categorized inbound replies with sentiment score badges.
    - Lipgloss styles: Positive (`#22c55e`, "POSITIVE"), OOO (`#3b82f6`, "AUTO-REPLY"), Negative (`#ef4444`, "OPTOUT").
  - **Web Sentiment Stream:**
    - Filterable feed allowing live search by sentiment category and keyword tokens.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup replies sentiment`
  - Flags:
    - `--pool-id <string>`: Target warmup identity pool.
    - `--filter-sentiment <string>`: Filter (`all`, `positive`, `hostile`, `ooo`) (default: `all`).
    - `--limit <int>`: Number of recent replies to analyze (default: `50`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ReplySentimentEntry struct {
        ReplyID         string    `json:"replyId"`
        FromDomain      string    `json:"fromDomain"`
        SentimentScore  float64   `json:"sentimentScore"`
        Classification  string    `json:"classification"`
        IsPositive      bool      `json:"isPositive"`
        IsAutoReply     bool      `json:"isAutoReply"`
        HasOptOutIntent bool      `json:"hasOptOutIntent"`
        ReceivedAt      time.Time `json:"receivedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupReplySentiments (
        replyId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        fromDomain TEXT NOT NULL,
        sentimentScore REAL NOT NULL,
        classification TEXT NOT NULL,
        isPositive INTEGER NOT NULL DEFAULT 0,
        isAutoReply INTEGER NOT NULL DEFAULT 0,
        hasOptOutIntent INTEGER NOT NULL DEFAULT 0,
        receivedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_reply_sentiments_pool ON WarmupReplySentiments(poolId, isPositive);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrates into `pkg/warmup/sentiment_analyzer.go` inside the local GitMap message pipeline.

---

### Improvement 540: Automated Warmup Score Calibration Based on Recipient Dwell & Click Ratio
- **ID:** `ENGAGE-LOOP-540`
- **Cluster:** 54 (Real-Time Recipient Engagement Feedback Loops & Open/Click Sparklines)
- **Goal:** Compute a holistic Reputation Momentum Score (0-100) combining open dwell duration, click engagement ratio, bounce velocity, and spam trap avoidance into an authoritative single gauge.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Radial Dial:**
    - Large ANSI terminal dial rendering current momentum score:
      ```text
      ┌─ [WARMUP MOMENTUM SCORE] ────────────────┐
      │             ╭─────────╮                  │
      │          88 │  OPTIMAL│                  │
      │             ╰─────────╯                  │
      │ Open: 92/100 | Click: 84/100 | Dwell: 88 │
      └──────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Gauge:**
    - Circular SVG meter (`stroke-dasharray`, `stroke-dashoffset`) with zero JavaScript framework dependencies.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email warmup score`
  - Flags:
    - `--pool-id <string>`: Warmup identity pool (default: `default`).
    - `--json`: Output raw JSON telemetry for CI/CD assertions.
    - `--fail-under <int>`: Exit with status code 1 if score drops below threshold (default: `0`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type WarmupMomentumScore struct {
        PoolID           string    `json:"poolId"`
        OverallScore     int       `json:"overallScore"`
        OpenFactor       int       `json:"openFactor"`
        ClickFactor      int       `json:"clickFactor"`
        DwellFactor      int       `json:"dwellFactor"`
        DeliverFactor    int       `json:"deliverFactor"`
        IsGraduated      bool      `json:"isGraduated"`
        HasRiskFactors   bool      `json:"hasRiskFactors"`
        CalculatedAt     time.Time `json:"calculatedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WarmupMomentumScores (
        scoreId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        overallScore INTEGER NOT NULL,
        openFactor INTEGER NOT NULL,
        clickFactor INTEGER NOT NULL,
        dwellFactor INTEGER NOT NULL,
        deliverFactor INTEGER NOT NULL,
        isGraduated INTEGER NOT NULL DEFAULT 0,
        hasRiskFactors INTEGER NOT NULL DEFAULT 0,
        calculatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_momentum_pool_time ON WarmupMomentumScores(poolId, calculatedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Dispatched via `cli/cmdemail/warmup_score.go` using calibration routines from `pkg/warmup/scoring_engine.go`.

---

# Cluster 55: Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette (Items 541–550)

---

### Improvement 541: Interactive Bubble Tea Mailbox Inspector & Split Navigation Frame
- **ID:** `MAIL-INSPECT-541`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Full-screen dual-pane terminal TUI allowing operators to navigate warmup mailbox folders (Inbox, Sent, Spam, Trash) on the left and preview raw headers, MIME tree, and body on the right.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Two-Pane Model:**
    - Left Pane (Folder & Message List, 35% width): Scrollable table with flags `[UNREAD]`, `[WARMUP]`, `[SPAM]`.
    - Right Pane (Viewer, 65% width): Tabbed display (`[1] Plain Text`, `[2] HTML Render`, `[3] RFC 5322 Headers`, `[4] DKIM/SPF`).
    - Key bindings: `j`/`k` (move message), `Tab` (switch pane), `1-4` (switch tab), `q` (exit).
  - **Zero-CDN Web View:**
    - Desktop web mail client replica rendered with pure CSS grid, instant keyboard hotkeys, and zero bundle size.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox inspect`
  - Flags:
    - `--account <string>`: Email account identity (e.g. `sender-01@domain.com`).
    - `--folder <string>`: Initial folder to load (default: `INBOX`).
    - `--read-only`: Disallow destructive operations like delete or mark spam (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxMessageSummary struct {
        MessageID    string    `json:"messageId"`
        AccountID    string    `json:"accountId"`
        Folder       string    `json:"folder"`
        Subject      string    `json:"subject"`
        FromAddress  string    `json:"fromAddress"`
        DateReceived time.Time `json:"dateReceived"`
        IsRead       bool      `json:"isRead"`
        IsSpam       bool      `json:"isSpam"`
        HasWarmupTag bool      `json:"hasWarmupTag"`
    }

    type MailboxInspectorState struct {
        ActiveAccount string                  `json:"activeAccount"`
        ActiveFolder  string                  `json:"activeFolder"`
        SelectedIndex int                     `json:"selectedIndex"`
        Messages      []MailboxMessageSummary `json:"messages"`
        IsFocused     bool                    `json:"isFocused"`
        HasNextPage   bool                    `json:"hasNextPage"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxMessageSummaries (
        messageId TEXT PRIMARY KEY,
        accountId TEXT NOT NULL,
        folder TEXT NOT NULL,
        subject TEXT NOT NULL,
        fromAddress TEXT NOT NULL,
        dateReceived DATETIME NOT NULL,
        isRead INTEGER NOT NULL DEFAULT 0,
        isSpam INTEGER NOT NULL DEFAULT 0,
        hasWarmupTag INTEGER NOT NULL DEFAULT 0,
        cachedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_mailbox_acct_folder ON MailboxMessageSummaries(accountId, folder, dateReceived);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/mailbox_inspector_tui.go` communicating via IMAP4rev1 client in `pkg/email/imap_client.go`.

---

### Improvement 542: Quick-Action Command Palette (`Ctrl+K`) for Instant Mailbox Operations
- **ID:** `MAIL-INSPECT-542`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Implement an instant fuzzy-search command palette overlay (`Ctrl+K` or `:`) inside the Bubble Tea inspector for rapid actions (e.g. `:move-inbox`, `:rescue-spam`, `:whitelist-sender`, `:resend-probe`).
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Modal Palette:**
    - Floating centered input box styled with:
      `lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#ec4899")).Padding(1, 2)`.
    - Real-time fuzzy filtering matches typing against available command handlers.
    - Quick shortcut hints displayed on the right: `[Enter: Run]`, `[Esc: Dismiss]`.
  - **Zero-CDN Web Palette:**
    - Popover modal dialog with `<input type="text">` auto-focused and keyboard trap listeners.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox command-palette`
  - Flags:
    - `--account <string>`: Email account identity.
    - `--execute <string>`: Headless execution of palette command without launching TUI.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type PaletteCommand struct {
        CommandSlug string `json:"commandSlug"`
        Title       string `json:"title"`
        Description string `json:"description"`
        ShortcutKey string `json:"shortcutKey"`
        IsEnabled   bool   `json:"isEnabled"`
        HasParams   bool   `json:"hasParams"`
    }

    type PaletteRegistry struct {
        Commands    []PaletteCommand `json:"commands"`
        IsLoaded    bool             `json:"isLoaded"`
        HasCustomOps bool            `json:"hasCustomOps"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxPaletteHistory (
        executionId TEXT PRIMARY KEY,
        accountId TEXT NOT NULL,
        commandSlug TEXT NOT NULL,
        targetMessageId TEXT,
        isSuccess INTEGER NOT NULL DEFAULT 1,
        hasError INTEGER NOT NULL DEFAULT 0,
        executedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_palette_history ON MailboxPaletteHistory(accountId, executedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Plugs directly into `cli/cmdemail/palette_component.go` as a reusable Bubble Tea sub-model.

---

### Improvement 543: Autonomous Spam Folder Rescue & Whitelist Enforcer TUI
- **ID:** `MAIL-INSPECT-543`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** One-click automated sweep identifying peer warmup emails mistakenly routed to the Spam folder, marking them Not Spam, moving to Inbox, and starring/replying to train Bayesian filters.
- **UI/UX Terminal & Web Specification:**
  - **Interactive TUI Sweep Modal:**
    - Progress spinner with step-by-step resolution badges:
      `[1/4] Scanning Spam Folder... [2/4] Found 3 Warmup Peers... [3/4] Unmarking Spam... [4/4] Moved to Inbox`.
    - Final confirmation state rendered in emerald green.
  - **Web Quick Sweep Card:**
    - One-click "Rescue & Whitelist All" button with animated progress indicator.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox rescue-spam`
  - Flags:
    - `--account <string>`: Target account identity.
    - `--auto-reply`: Send positive acknowledgment reply to peer upon rescue (default: `true`).
    - `--star-message`: Apply star/flag to reinforce inbox classification (default: `true`).
    - `--dry-run`: Preview messages eligible for rescue without performing IMAP modifications.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SpamRescueResult struct {
        AccountID      string   `json:"accountId"`
        RescuedCount   int      `json:"rescuedCount"`
        MessageIDs     []string `json:"messageIds"`
        IsAutoReplied  bool     `json:"isAutoReplied"`
        IsStarred      bool     `json:"isStarred"`
        HasFailures    bool     `json:"hasFailures"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpamRescueAudits (
        rescueId TEXT PRIMARY KEY,
        accountId TEXT NOT NULL,
        rescuedCount INTEGER NOT NULL DEFAULT 0,
        isAutoReplied INTEGER NOT NULL DEFAULT 0,
        isStarred INTEGER NOT NULL DEFAULT 0,
        hasFailures INTEGER NOT NULL DEFAULT 0,
        executedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_rescue_audits_acct ON SpamRescueAudits(accountId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented inside `pkg/warmup/spam_rescuer.go` and surfaced through `gitmap email mailbox rescue-spam`.

---

### Improvement 544: RFC 5322 Raw Header Diff Inspector with Syntax Highlighting
- **ID:** `MAIL-INSPECT-544`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Compare received RFC 5322 email headers against expected template headers, highlighting mutated or missing fields (`Message-ID`, `List-Unsubscribe`, `DKIM-Signature`, `X-Entity-Ref-ID`).
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Split Diff View:**
    - Left side: Expected Headers (dim cyan).
    - Right side: Received Headers from remote IMAP (bold white).
    - Added headers shown in green (`+`), modified values highlighted in yellow, missing headers in red (`-`).
  - **Zero-CDN Web Diff:**
    - Inline HTML `<pre><code>` diff block with GitHub-like red/green line highlighting using CSS classes.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox header-diff`
  - Flags:
    - `--message-id <string>`: Target email ID in mailbox.
    - `--template-id <string>`: Reference dispatch template ID.
    - `--ignore-date`: Exclude timestamp drift from diff comparison (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type HeaderDiffItem struct {
        HeaderKey    string `json:"headerKey"`
        ExpectedVal  string `json:"expectedVal"`
        ActualVal    string `json:"actualVal"`
        IsIdentical  bool   `json:"isIdentical"`
        IsMissing    bool   `json:"isMissing"`
        HasAnomaly   bool   `json:"hasAnomaly"`
    }

    type HeaderDiffReport struct {
        MessageID    string           `json:"messageId"`
        Diffs        []HeaderDiffItem `json:"diffs"`
        IsCompliant  bool             `json:"isCompliant"`
        HasMimeError bool             `json:"hasMimeError"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxHeaderDiffs (
        diffId TEXT PRIMARY KEY,
        messageId TEXT NOT NULL,
        templateId TEXT NOT NULL,
        isCompliant INTEGER NOT NULL DEFAULT 1,
        hasMimeError INTEGER NOT NULL DEFAULT 0,
        analyzedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_header_diffs_msg ON MailboxHeaderDiffs(messageId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `cli/cmdemail/header_diff.go` using `pkg/mime/header_parser.go`.

---

### Improvement 545: Interactive In-Terminal HTML-to-Terminal Renderer
- **ID:** `MAIL-INSPECT-545`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Render complex HTML email bodies directly into formatted ANSI text within the terminal inspector, supporting headings, hyperlinks, lists, and table formatting without launching external browsers.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss HTML Parser & Styler:**
    - `<h1>` to `<h6>`: Bold violet styled with underline.
    - `<a>`: Cyan underlined with clickable OSC 8 hyperlink escape sequences.
    - `<table>`: Automatically converted into clean Unicode box-drawing tables.
    - `<img>`: Replaced with lightweight terminal alt-text pills `[IMG: Alt description]`.
  - **Web Dual Mode:**
    - Sandboxed `<iframe>` with strict CSP (`default-src 'none'`) preventing remote tracking pixel execution during inspection.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox render-body`
  - Flags:
    - `--message-id <string>`: Target email ID.
    - `--width <int>`: Terminal wrap width constraint (default: `80`).
    - `--block-remote-images`: Suppress remote images in HTML rendering (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type RenderedBody struct {
        MessageID     string `json:"messageId"`
        PlainText     string `json:"plainText"`
        FormattedANSI string `json:"formattedAnsi"`
        ImageCount    int    `json:"imageCount"`
        LinkCount     int    `json:"linkCount"`
        IsSanitized   bool   `json:"isSanitized"`
        HasImages     bool   `json:"hasImages"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RenderedBodyCaches (
        cacheId TEXT PRIMARY KEY,
        messageId TEXT NOT NULL,
        isSanitized INTEGER NOT NULL DEFAULT 1,
        hasImages INTEGER NOT NULL DEFAULT 0,
        cachedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_rendered_body_cache ON RenderedBodyCaches(messageId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/terminal/ansi_html_renderer.go` and integrated into the Bubble Tea inspector model.

---

### Improvement 546: Multi-Account Mailbox Switcher & Aggregated Inbound Feed
- **ID:** `MAIL-INSPECT-546`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Instant account jumping (`Alt+1` through `Alt+9` or dropdown) allowing operators to inspect all warmup sending inboxes in a unified multi-mailbox view.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Status Bar Account Picker:**
    - Active account pills rendered along the bottom: `[1: sender-01]`, `[2: sender-02*]`, `[3: sender-03]`.
    - Star badge `*` denotes accounts with new incoming unread peer replies.
  - **Zero-CDN Web View:**
    - Sidebar account switcher with unread notification badges styled with CSS pills.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox unified-feed`
  - Flags:
    - `--pool-id <string>`: Warmup identity pool containing multiple accounts.
    - `--unread-only`: Display only messages marked as unread (default: `false`).
    - `--limit <int>`: Total messages to load across all inboxes (default: `100`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type UnifiedFeedMessage struct {
        MessageID    string    `json:"messageId"`
        AccountEmail string    `json:"accountEmail"`
        Subject      string    `json:"subject"`
        Snippet      string    `json:"snippet"`
        ReceivedAt   time.Time `json:"receivedAt"`
        IsRead       bool      `json:"isRead"`
        HasReply     bool      `json:"hasReply"`
    }

    type UnifiedFeedState struct {
        TotalAccounts int                  `json:"totalAccounts"`
        Messages      []UnifiedFeedMessage `json:"messages"`
        IsLive        bool                 `json:"isLive"`
        HasUnread     bool                 `json:"hasUnread"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS UnifiedInboundFeeds (
        feedItemId TEXT PRIMARY KEY,
        messageId TEXT NOT NULL,
        accountEmail TEXT NOT NULL,
        subject TEXT NOT NULL,
        isRead INTEGER NOT NULL DEFAULT 0,
        hasReply INTEGER NOT NULL DEFAULT 0,
        receivedAt DATETIME NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_unified_feed_acct ON UnifiedInboundFeeds(accountEmail, receivedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Backed by `pkg/email/pool_aggregator.go` with fast asynchronous IMAP connection pooling.

---

### Improvement 547: Instant Reply Generator & Conversational Threading Simulator
- **ID:** `MAIL-INSPECT-547`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Draft and dispatch contextual, human-like peer replies directly from the inspector terminal to fulfill warm-up reply ratio quotas without leaving the TUI.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Reply Composer:**
    - Bottom sheet overlay activating upon pressing `r`.
    - Template selector (`[1] Casual Acknowledgment`, `[2] In-depth Query`, `[3] Meeting Confirmation`).
    - Real-time preview of `In-Reply-To` and `References` RFC headers.
  - **Zero-CDN Web Composer:**
    - Clean inline text area with instant Send button and auto-generated subject lines (`Re: ...`).
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox send-reply`
  - Flags:
    - `--message-id <string>`: Original message ID being replied to.
    - `--template <string>`: Response template preset (`casual`, `business`, `query`).
    - `--auto-dispatch`: Send immediately without interactive review (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type OutboundReplyDraft struct {
        DraftID         string    `json:"draftId"`
        ParentMessageID string    `json:"parentMessageId"`
        SenderAccount   string    `json:"senderAccount"`
        RecipientEmail  string    `json:"recipientEmail"`
        Subject         string    `json:"subject"`
        BodyText        string    `json:"bodyText"`
        IsDispatched    bool      `json:"isDispatched"`
        HasReferences   bool      `json:"hasReferences"`
        CreatedAt       time.Time `json:"createdAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS OutboundWarmupReplies (
        draftId TEXT PRIMARY KEY,
        parentMessageId TEXT NOT NULL,
        senderAccount TEXT NOT NULL,
        recipientEmail TEXT NOT NULL,
        isDispatched INTEGER NOT NULL DEFAULT 0,
        hasReferences INTEGER NOT NULL DEFAULT 1,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_outbound_replies_parent ON OutboundWarmupReplies(parentMessageId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Uses `pkg/warmup/reply_composer.go` and hooks into the GitMap SMTP sender queue.

---

### Improvement 548: IMAP Folder Hierarchy & Flag Synchronization Inspector
- **ID:** `MAIL-INSPECT-548`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Display real-time IMAP server folder structures (RFC 6154 Special-Use mailboxes: `\Inbox`, `\Sent`, `\Junk`, `\Trash`, `\Archive`) and inspect mailbox UIDValidity states.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Tree View:**
    - Hierarchical folder tree with message count pills:
      ```text
      📂 INBOX (14 msgs, 2 unread)
      ├── 📂 Sent (42 msgs)
      ├── 📂 Junk [SPAM] (1 msg) ⚠️
      └── 📂 Trash (0 msgs)
      ```
  - **Web Folder Tree:**
    - Minimalist collapsible list element with folder icon SVGs.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox folders`
  - Flags:
    - `--account <string>`: Target email account.
    - `--sync-flags`: Query remote server for updated unread/flagged counts (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type ImapFolderInfo struct {
        FolderName   string `json:"folderName"`
        SpecialUse   string `json:"specialUse"`
        TotalCount   int    `json:"totalCount"`
        UnreadCount  int    `json:"unreadCount"`
        UIDValidity  uint32 `json:"uidValidity"`
        IsSpecialUse bool   `json:"isSpecialUse"`
        HasUnread    bool   `json:"hasUnread"`
    }

    type AccountFolderStructure struct {
        AccountEmail string           `json:"accountEmail"`
        Folders      []ImapFolderInfo `json:"folders"`
        IsSynced     bool             `json:"isSynced"`
        HasErrors    bool             `json:"hasErrors"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AccountFolderCaches (
        folderId TEXT PRIMARY KEY,
        accountEmail TEXT NOT NULL,
        folderName TEXT NOT NULL,
        specialUse TEXT NOT NULL,
        totalCount INTEGER NOT NULL DEFAULT 0,
        unreadCount INTEGER NOT NULL DEFAULT 0,
        isSpecialUse INTEGER NOT NULL DEFAULT 0,
        hasUnread INTEGER NOT NULL DEFAULT 0,
        syncedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_folder_caches_acct ON AccountFolderCaches(accountEmail);
    ```
- **Antigravity / GitMap Integration Point:**
  - Backed by `pkg/email/imap_inspector.go` inside `cli/cmdemail/mailbox_folders.go`.

---

### Improvement 549: In-Memory Message Search & Regex Filter Engine
- **ID:** `MAIL-INSPECT-549`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Instant client-side search across subject lines, sender domains, and raw MIME bodies using compiled regex queries with sub-millisecond execution over local SQLite cache.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Search Bar:**
    - Interactive `/` trigger activating top query input.
    - Matched substrings highlighted with reverse video style (`lipgloss.NewStyle().Reverse(true)`).
    - Result count ticker: `[Match 3 of 12]`.
  - **Zero-CDN Web Search:**
    - Live search input filtering table rows on the fly with CSS-only highlighting.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox search`
  - Flags:
    - `--account <string>`: Target account.
    - `--query <string>`: Plaintext or regex search pattern.
    - `--in-body`: Search message body in addition to headers (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type SearchResultItem struct {
        MessageID   string `json:"messageId"`
        Subject     string `json:"subject"`
        FromAddress string `json:"fromAddress"`
        MatchSnippet string `json:"matchSnippet"`
        IsBodyMatch bool   `json:"isBodyMatch"`
        HasHeaderMatch bool `json:"hasHeaderMatch"`
    }

    type SearchResultSet struct {
        QueryPattern string             `json:"queryPattern"`
        TotalMatches int                `json:"totalMatches"`
        Items        []SearchResultItem `json:"items"`
        IsComplete   bool               `json:"isComplete"`
        HasMatches   bool               `json:"hasMatches"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE VIRTUAL TABLE IF NOT EXISTS MailboxSearchFTS USING fts5(
        messageId,
        subject,
        fromAddress,
        bodyContent
    );
    ```
- **Antigravity / GitMap Integration Point:**
  - Driven by SQLite FTS5 engine in `pkg/storage/sqlite_search.go` through `cli/cmdemail/search_cmd.go`.

---

### Improvement 550: Mailbox Diagnostics & Authentication Health Card
- **ID:** `MAIL-INSPECT-550`
- **Cluster:** 55 (Interactive Bubble Tea Mailbox Inspector & Quick-Action Command Palette)
- **Goal:** Comprehensive health diagnostic card rendering IMAP/SMTP connectivity state, TLS cipher negotiation, OAuth2 token expiry, and storage quota status.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Diagnostic Card:**
    - Rounded border container with color-coded key-value badges:
      ```text
      ┌─ [ACCOUNT DIAGNOSTICS: sender-01@domain.com] ──────────┐
      │ IMAP Connection : ACTIVE (TLS 1.3 / AES-256-GCM)       │
      │ SMTP Connection : ACTIVE (Port 587 STARTTLS)           │
      │ OAuth2 Validity : VALID (Expires in 42m)               │
      │ Mailbox Quota   : 142MB / 5000MB [ ▂       ] 2.8%      │
      └────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Status Card:**
    - Minimalist dashboard card with clean status indicator dots and quota progress bar.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email mailbox diagnostics`
  - Flags:
    - `--account <string>`: Target account identifier.
    - `--test-smtp`: Perform test SMTP handshake without sending message (default: `true`).
    - `--test-imap`: Verify IMAP login and folder listing (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MailboxDiagnosticsReport struct {
        AccountEmail  string    `json:"accountEmail"`
        IsImapHealthy bool      `json:"isImapHealthy"`
        IsSmtpHealthy bool      `json:"isSmtpHealthy"`
        IsTokenValid  bool      `json:"isTokenValid"`
        HasQuotaAlert bool      `json:"hasQuotaAlert"`
        UsedQuotaMB   int       `json:"usedQuotaMb"`
        TotalQuotaMB  int       `json:"totalQuotaMb"`
        CheckedAt     time.Time `json:"checkedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MailboxDiagnosticsAudits (
        auditId TEXT PRIMARY KEY,
        accountEmail TEXT NOT NULL,
        isImapHealthy INTEGER NOT NULL DEFAULT 1,
        isSmtpHealthy INTEGER NOT NULL DEFAULT 1,
        isTokenValid INTEGER NOT NULL DEFAULT 1,
        hasQuotaAlert INTEGER NOT NULL DEFAULT 0,
        usedQuotaMb INTEGER NOT NULL DEFAULT 0,
        totalQuotaMb INTEGER NOT NULL DEFAULT 0,
        checkedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_diagnostics_acct ON MailboxDiagnosticsAudits(accountEmail);
    ```
- **Antigravity / GitMap Integration Point:**
  - Located in `cli/cmdemail/mailbox_diagnostics.go` using tests from `pkg/email/health_prober.go`.

---

# Cluster 56: Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer (Items 551–560)

---

### Improvement 551: Memory-Mapped Ring Buffer Queue Depth HUD
- **ID:** `QUEUE-FLOW-551`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Provide real-time visualization of high-throughput memory-mapped dispatch queues, showing buffer fullness, uncommitted head/tail pointers, and overflow indicators.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Terminal Meter:**
    - High-density segmented ASCII gauge representing queue capacity:
      `QUEUE CAPACITY: [████████████████░░░░░░░░] 64,280 / 100,000 (64.2%) [Read: 1,420/s | Write: 1,890/s]`
    - Lipgloss tokens: Green (`<70%`), Yellow (`70-89%`), Crimson (`>90% Critical`).
  - **Zero-CDN Web View:**
    - Live CSS gradient bar updated via WebSocket with millisecond fidelity.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email queue hud`
  - Flags:
    - `--queue-name <string>`: Target queue identifier (default: `outbound-warmup`).
    - `--refresh-ms <int>`: Telemetry refresh period (default: `100`).
    - `--watch`: Interactive full-screen terminal monitor.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type QueueMemoryTelemetry struct {
        QueueName       string  `json:"queueName"`
        HeadPointer     uint64  `json:"headPointer"`
        TailPointer     uint64  `json:"tailPointer"`
        AllocatedBytes  int64   `json:"allocatedBytes"`
        MaxCapacity     int     `json:"maxCapacity"`
        CurrentLength   int     `json:"currentLength"`
        UtilizationRatio float64 `json:"utilizationRatio"`
        IsFull          bool    `json:"isFull"`
        HasBackpressure bool    `json:"hasBackpressure"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS QueueBufferSnapshots (
        snapshotId TEXT PRIMARY KEY,
        queueName TEXT NOT NULL,
        currentLength INTEGER NOT NULL,
        maxCapacity INTEGER NOT NULL,
        utilizationRatio REAL NOT NULL,
        isFull INTEGER NOT NULL DEFAULT 0,
        hasBackpressure INTEGER NOT NULL DEFAULT 0,
        sampledAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_queue_snapshots_time ON QueueBufferSnapshots(queueName, sampledAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Direct zero-copy inspection of `pkg/queue/mmap_ringbuffer.go` through `cli/cmdemail/queue_hud.go`.

---

### Improvement 552: Leaky Bucket Adaptive Dispatch Throttler TUI
- **ID:** `QUEUE-FLOW-552`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Terminal-based controller modeling token leaky-bucket rate limiters per remote MX domain, showing dynamic leak rates, accumulated tokens, and delay penalties.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Bucket Representation:**
    - Multi-domain bucket columns (Google, Microsoft, Yahoo):
      ```text
      ┌─ [LEAKY BUCKET DISPATCH LIMITERS] ──────────────────────────┐
      │ GMAIL    : [██████░░░░] 60%  Rate: 10/min  Leak: 0.16/s     │
      │ OUTLOOK  : [████████░░] 80%  Rate:  5/min  Leak: 0.08/s ⚠️  │
      │ YAHOO    : [██░░░░░░░░] 20%  Rate: 20/min  Leak: 0.33/s     │
      └─────────────────────────────────────────────────────────────┘
      ```
  - **Web Rate Limiter Panel:**
    - Multi-tier visual liquid meters with CSS keyframe animation reflecting token leakage.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email throttler status`
  - Flags:
    - `--domain <string>`: Filter by remote target domain (default: `all`).
    - `--reset-tokens`: Clear token accumulation for testing (default: `false`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type LeakyBucketTelemetry struct {
        TargetDomain   string  `json:"targetDomain"`
        CurrentTokens  int     `json:"currentTokens"`
        MaxBurst       int     `json:"maxBurst"`
        LeakRatePerSec float64 `json:"leakRatePerSec"`
        IsThrottled    bool    `json:"isThrottled"`
        HasCapacity    bool    `json:"hasCapacity"`
    }

    type ThrottlerPoolState struct {
        PoolID    string                 `json:"poolId"`
        Buckets   []LeakyBucketTelemetry `json:"buckets"`
        IsOptimal bool                   `json:"isOptimal"`
        HasSpills bool                   `json:"hasSpills"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ThrottlerBucketStates (
        stateId TEXT PRIMARY KEY,
        poolId TEXT NOT NULL,
        targetDomain TEXT NOT NULL,
        currentTokens INTEGER NOT NULL,
        maxBurst INTEGER NOT NULL,
        leakRatePerSec REAL NOT NULL,
        isThrottled INTEGER NOT NULL DEFAULT 0,
        hasCapacity INTEGER NOT NULL DEFAULT 1,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX IF NOT EXISTS uq_throttler_domain ON ThrottlerBucketStates(poolId, targetDomain);
    ```
- **Antigravity / GitMap Integration Point:**
  - Integrated into `pkg/throttler/leaky_bucket.go` and controlled via `cli/cmdemail/throttler_cmd.go`.

---

### Improvement 553: Bounded Backpressure Pushback & Worker Pool Pause Matrix
- **ID:** `QUEUE-FLOW-553`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Automatically pause message ingestion from upstream producers when SMTP connection latency or remote 421/451 deferrals push downstream memory consumption over safety boundaries.
- **UI/UX Terminal & Web Specification:**
  - **Terminal Pushback Indicator:**
    - High-visibility amber status badge in Bubble Tea footer: `[BACKPRESSURE: ENGAGED] Ingestion Paused | Workers Active: 12/32`.
    - Real-time countdown timer before ingestion resumes.
  - **Web Notice Bar:**
    - Top sticky toast showing backpressure trigger cause and recovery progress.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email backpressure monitor`
  - Flags:
    - `--pool-id <string>`: Target warmup identity pool.
    - `--memory-limit-mb <int>`: Memory ceiling triggering backpressure (default: `512`).
    - `--auto-resume`: Automatically resume ingestion once memory drops (default: `true`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type BackpressureEvent struct {
        EventID       string    `json:"eventId"`
        QueueName     string    `json:"queueName"`
        TriggerReason string    `json:"triggerReason"`
        MemoryUsageMB int       `json:"memoryUsageMb"`
        ActiveWorkers int       `json:"activeWorkers"`
        IsEngaged     bool      `json:"isEngaged"`
        HasRecovered  bool      `json:"hasRecovered"`
        CreatedAt     time.Time `json:"createdAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BackpressureHistory (
        eventId TEXT PRIMARY KEY,
        queueName TEXT NOT NULL,
        triggerReason TEXT NOT NULL,
        memoryUsageMb INTEGER NOT NULL,
        activeWorkers INTEGER NOT NULL,
        isEngaged INTEGER NOT NULL DEFAULT 1,
        hasRecovered INTEGER NOT NULL DEFAULT 0,
        createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_backpressure_time ON BackpressureHistory(queueName, createdAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Orchestrated in `pkg/queue/backpressure_controller.go` with IPC event dispatch to Lemon.

---

### Improvement 554: Dead-Letter Queue (DLQ) Inspector & Surgical Replay TUI
- **ID:** `QUEUE-FLOW-554`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Interactive terminal utility to inspect discarded or poison messages in the Dead-Letter Queue (DLQ), view failure stack traces, and execute surgical single-item or batch replay.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea DLQ Table:**
    - Columns: `MessageID`, `Target Domain`, `Retry Count`, `Failure Reason`, `Timestamp`.
    - Key actions: `[r] Replay Selected`, `[a] Replay All Safe`, `[d] Purge Poison Item`, `[v] View Stack Trace`.
  - **Zero-CDN Web DLQ Interface:**
    - Minimalist table with inline error trace expandable drawer and action buttons.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email dlq inspect`
  - Flags:
    - `--queue <string>`: DLQ storage identifier (default: `warmup-dlq`).
    - `--replay-id <string>`: Headless replay of specific message ID.
    - `--purge-all`: Empty the entire DLQ (requires confirmation flag).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type DLQMessageEntry struct {
        MessageID     string    `json:"messageId"`
        TargetDomain  string    `json:"targetDomain"`
        RetryCount    int       `json:"retryCount"`
        ErrorMessage  string    `json:"errorMessage"`
        PayloadBytes  int       `json:"payloadBytes"`
        IsReplayable  bool      `json:"isReplayable"`
        HasPermanentErr bool    `json:"hasPermanentErr"`
        EnqueuedAt    time.Time `json:"enqueuedAt"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DeadLetterQueueItems (
        messageId TEXT PRIMARY KEY,
        queueName TEXT NOT NULL,
        targetDomain TEXT NOT NULL,
        retryCount INTEGER NOT NULL DEFAULT 0,
        errorMessage TEXT NOT NULL,
        payloadBytes INTEGER NOT NULL DEFAULT 0,
        isReplayable INTEGER NOT NULL DEFAULT 1,
        hasPermanentErr INTEGER NOT NULL DEFAULT 0,
        enqueuedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_dlq_queue_domain ON DeadLetterQueueItems(queueName, targetDomain);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `cli/cmdemail/dlq_inspector.go` linking with `pkg/queue/dlq_manager.go`.

---

### Improvement 555: Memory-Mapped Ring Buffer Persistence & Recovery Daemon
- **ID:** `QUEUE-FLOW-555`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Ensure zero message loss across unexpected process restarts by syncing memory-mapped ring buffers to persistent disk pages (`msync` with `MS_ASYNC`/`MS_SYNC`) and visualizing page sync dirty counts.
- **UI/UX Terminal & Web Specification:**
  - **Terminal Sync Gauge:**
    - Live sync indicator: `[PAGE SYNC] Dirty Pages: 4 / 256 (1.5%) | Last Disk Flush: 12ms ago [SAFE]`.
    - Lipgloss styles: Clean (`#22c55e`), Flushing (`#3b82f6`), Degraded Sync (`#ef4444`).
  - **Web Sync Visualizer:**
    - Status icon with pulse animation indicating periodic background flush health.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email queue sync-status`
  - Flags:
    - `--queue-name <string>`: Queue file name.
    - `--force-flush`: Force immediate synchronous `msync` to disk.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type MMapSyncState struct {
        QueueName      string        `json:"queueName"`
        TotalPages     int           `json:"totalPages"`
        DirtyPages     int           `json:"dirtyPages"`
        LastSyncTime   time.Time     `json:"lastSyncTime"`
        SyncDuration   time.Duration `json:"syncDuration"`
        IsClean        bool          `json:"isClean"`
        HasSyncPending bool          `json:"hasSyncPending"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS QueueDiskSyncAudits (
        auditId TEXT PRIMARY KEY,
        queueName TEXT NOT NULL,
        totalPages INTEGER NOT NULL,
        dirtyPages INTEGER NOT NULL,
        syncDurationMs INTEGER NOT NULL,
        isClean INTEGER NOT NULL DEFAULT 1,
        hasSyncPending INTEGER NOT NULL DEFAULT 0,
        syncedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_sync_audits_time ON QueueDiskSyncAudits(queueName, syncedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Engine code residing in `pkg/queue/mmap_storage.go` accessed via `cli/cmdemail/queue_sync.go`.

---

### Improvement 556: Priority Lane Dispatch Matrix & VIP Envelope Elevator
- **ID:** `QUEUE-FLOW-556`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Multi-lane queue scheduling (High, Standard, Bulk Warmup) preventing large bulk warmup batches from starving critical operational probe envelopes or peer confirmation replies.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Multi-Lane Progress Bars:**
    - Lane 1 [VIP/Probe]: `[██████████] 100% (Instant 0ms wait)` (Cyan)
    - Lane 2 [Standard]:  `[██████░░░░]  60% (Avg wait: 420ms)` (Blue)
    - Lane 3 [Warmup]:    `[████░░░░░░]  40% (Throttled pacing)` (Slate)
  - **Web Queue Lanes:**
    - Stacked card view showing real-time item counts and dispatch throughput per priority lane.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email queue lanes`
  - Flags:
    - `--view <string>`: Output format (`tui`, `json`, `csv`) (default: `tui`).
    - `--promote-id <string>`: Manually promote an enqueued message to VIP Lane.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type QueueLaneStatus struct {
        LaneName       string  `json:"laneName"`
        PriorityLevel  int     `json:"priorityLevel"`
        EnqueuedItems  int     `json:"enqueuedItems"`
        ThroughputPerS float64 `json:"throughputPerS"`
        IsStarved      bool    `json:"isStarved"`
        HasActiveFlow  bool    `json:"hasActiveFlow"`
    }

    type MultiLaneState struct {
        Lanes          []QueueLaneStatus `json:"lanes"`
        TotalInFlight  int               `json:"totalInFlight"`
        IsBalanced     bool              `json:"isBalanced"`
        HasCongestion  bool              `json:"hasCongestion"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS QueuePriorityLanes (
        laneId TEXT PRIMARY KEY,
        laneName TEXT NOT NULL,
        priorityLevel INTEGER NOT NULL,
        enqueuedItems INTEGER NOT NULL DEFAULT 0,
        throughputPerS REAL NOT NULL DEFAULT 0.0,
        isStarved INTEGER NOT NULL DEFAULT 0,
        hasActiveFlow INTEGER NOT NULL DEFAULT 1,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_priority_lanes ON QueuePriorityLanes(priorityLevel);
    ```
- **Antigravity / GitMap Integration Point:**
  - Managed by `pkg/queue/priority_scheduler.go` with execution bindings in `cli/cmdemail/queue_lanes.go`.

---

### Improvement 557: Dynamic Worker Pool Scaler with Latency Feedback Loop
- **ID:** `QUEUE-FLOW-557`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Autonomously scale concurrent SMTP worker goroutines (from min 4 to max 64) based on connection round-trip latency and queue backpressure signals.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Worker Swarm Visualizer:**
    - Real-time ASCII worker pool grid:
      `WORKER POOL [24 Active / 40 Idle]: [●●●●●●●●●●●●●●●●●●●●●●●●○○○○○○○○○○○○○○○○○○○○]`
    - Dynamic status updates: `[AUTOSCALER] Latency 14ms -> Scaled Up +4 Workers`.
  - **Web Worker HUD:**
    - Live SVG dot matrix representing active and standby goroutines.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email workers status`
  - Flags:
    - `--min-workers <int>`: Minimum pool size (default: `4`).
    - `--max-workers <int>`: Maximum pool size (default: `64`).
    - `--target-latency-ms <int>`: Desired RTT latency target (default: `50`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type WorkerPoolTelemetry struct {
        ActiveWorkers   int           `json:"activeWorkers"`
        IdleWorkers     int           `json:"idleWorkers"`
        MaxCapacity     int           `json:"maxCapacity"`
        AverageRTT      time.Duration `json:"averageRtt"`
        IsScalingUp     bool          `json:"isScalingUp"`
        IsThrottled     bool          `json:"isThrottled"`
        HasErrorSpike   bool          `json:"hasErrorSpike"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WorkerPoolAudits (
        auditId TEXT PRIMARY KEY,
        activeWorkers INTEGER NOT NULL,
        idleWorkers INTEGER NOT NULL,
        maxCapacity INTEGER NOT NULL,
        averageRttMs INTEGER NOT NULL,
        isScalingUp INTEGER NOT NULL DEFAULT 0,
        isThrottled INTEGER NOT NULL DEFAULT 0,
        hasErrorSpike INTEGER NOT NULL DEFAULT 0,
        recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_worker_audits_time ON WorkerPoolAudits(recordedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Implemented in `pkg/pool/worker_scaler.go` and controlled via `cli/cmdemail/workers_cmd.go`.

---

### Improvement 558: Memory-Mapped Storage Fragmentation & Compaction Monitor
- **ID:** `QUEUE-FLOW-558`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Real-time terminal tracking of fragmented and stale message payload segments in memory-mapped queue files with autonomous in-place compaction trigger.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Fragmentation Bar:**
    - Segmented bar showing used vs free fragmented blocks:
      `FRAGMENTATION: [■■■■■■■■□□□□■■■■] 24.8% Waste | 42MB Reclaimable [Press 'c' to Compact]`
    - Color tokens: Slate (`Used`), Amber (`Fragmented`), Gray (`Free`).
  - **Web Compaction Card:**
    - Visual storage breakdown with one-click "Run Live Compaction" action.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email queue compact`
  - Flags:
    - `--queue-name <string>`: Target queue file name.
    - `--force`: Run compaction without interactive prompt (default: `false`).
    - `--max-waste-pct <float>`: Automatic compaction trigger percentage (default: `30.0`).
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type FragmentationReport struct {
        QueueName       string  `json:"queueName"`
        TotalBytes      int64   `json:"totalBytes"`
        ReclaimableBytes int64  `json:"reclaimableBytes"`
        WastePercentage float64 `json:"wastePercentage"`
        IsCompacted     bool    `json:"isCompacted"`
        HasHighWaste    bool    `json:"hasHighWaste"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS QueueCompactionAudits (
        compactionId TEXT PRIMARY KEY,
        queueName TEXT NOT NULL,
        reclaimedBytes INTEGER NOT NULL,
        wastePercentage REAL NOT NULL,
        isCompacted INTEGER NOT NULL DEFAULT 1,
        hasHighWaste INTEGER NOT NULL DEFAULT 0,
        executedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_compaction_audits ON QueueCompactionAudits(queueName, executedAt);
    ```
- **Antigravity / GitMap Integration Point:**
  - Uses compaction routines in `pkg/queue/mmap_compactor.go` through `cli/cmdemail/queue_compact.go`.

---

### Improvement 559: End-to-End Envelope Flight Latency Tracker (Pipeline Tracing)
- **ID:** `QUEUE-FLOW-559`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Trace total elapsed time for individual warmup envelopes through the entire delivery pipeline: Ingestion -> Queue -> Rate Limiter -> TLS Handshake -> SMTP Dispatch -> Remote 250 OK.
- **UI/UX Terminal & Web Specification:**
  - **Lipgloss Waterline Waterfall Diagram:**
    - Multi-stage latency breakdown:
      ```text
      ┌─ [ENVELOPE FLIGHT TRACE: msg-8f92b4] ────────────────────────┐
      │ Queue Dwell      : 12ms   [■                      ]          │
      │ Throttler Delay  : 85ms   [■■■■■■                 ]          │
      │ DNS Resolution   :  4ms   [■                      ]          │
      │ TLS Handshake    : 28ms   [■■                     ]          │
      │ SMTP Transaction : 42ms   [■■■                    ]          │
      │ TOTAL TRANSIT    : 171ms  [STATUS: 250 OK]                   │
      └──────────────────────────────────────────────────────────────┘
      ```
  - **Zero-CDN Web Tracing Waterfall:**
    - Interactive horizontal bar chart styled purely with CSS grid.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email envelope trace`
  - Flags:
    - `--envelope-id <string>`: Identifier of the tracked envelope.
    - `--export-json`: Export raw telemetry for distributed tracing integration.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type EnvelopeStageDuration struct {
        StageName    string        `json:"stageName"`
        Duration     time.Duration `json:"duration"`
        IsBottleneck bool          `json:"isBottleneck"`
        HasError     bool          `json:"hasError"`
    }

    type EnvelopeFlightReport struct {
        EnvelopeID   string                  `json:"envelopeId"`
        TargetDomain string                  `json:"targetDomain"`
        TotalTransit time.Duration           `json:"totalTransit"`
        Stages       []EnvelopeStageDuration `json:"stages"`
        IsDelivered  bool                    `json:"isDelivered"`
        HasDelays    bool                    `json:"hasDelays"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS EnvelopeTransitTraces (
        traceId TEXT PRIMARY KEY,
        envelopeId TEXT NOT NULL,
        targetDomain TEXT NOT NULL,
        totalTransitMs INTEGER NOT NULL,
        isDelivered INTEGER NOT NULL DEFAULT 1,
        hasDelays INTEGER NOT NULL DEFAULT 0,
        tracedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_transit_traces_env ON EnvelopeTransitTraces(envelopeId);
    ```
- **Antigravity / GitMap Integration Point:**
  - Built into `cli/cmdemail/trace_envelope.go` consuming spans from `pkg/telemetry/envelope_tracer.go`.

---

### Improvement 560: Autonomous Queue Circuit Breaker with Auto-Reset Backoff
- **ID:** `QUEUE-FLOW-560`
- **Cluster:** 56 (Bounded Backpressure Throttler & Memory-Mapped Queue Visualizer)
- **Goal:** Intelligent circuit breaker halting outbound dispatches to a specific remote MX when consecutive connection timeouts or 5xx rejection rates exceed 20%, featuring exponential backoff half-open probes.
- **UI/UX Terminal & Web Specification:**
  - **Bubble Tea Circuit State Indicator:**
    - Tri-state visual badge:
      - `CLOSED` (`#22c55e`, emerald: normal operations).
      - `OPEN` (`#ef4444`, red: tripped, dispatches paused).
      - `HALF-OPEN` (`#eab308`, amber: testing single canary probe).
    - Countdown timer indicating remaining backoff cooldown before next canary probe.
  - **Web Circuit Matrix:**
    - Live grid showing breaker states across all major ISP destinations.
- **CLI Command Signature & Flags:**
  - Command: `gitmap email circuit-breaker status`
  - Flags:
    - `--domain <string>`: Filter by remote MX cluster (default: `all`).
    - `--force-trip <string>`: Manually open circuit for specific domain (emergency stop).
    - `--force-reset <string>`: Manually reset circuit to closed.
- **Data Contract & Split-DB Schema:**
  - Go Structs:
    ```go
    type CircuitBreakerState struct {
        TargetDomain    string        `json:"targetDomain"`
        State           string        `json:"state"` // "closed", "open", "half-open"
        FailureCount    int           `json:"failureCount"`
        SuccessCount    int           `json:"successCount"`
        BackoffDuration time.Duration `json:"backoffDuration"`
        IsOpen          bool          `json:"isOpen"`
        IsHalfOpen      bool          `json:"isHalfOpen"`
        IsClosed        bool          `json:"isClosed"`
        HasTripped      bool          `json:"hasTripped"`
        LastTripTime    time.Time     `json:"lastTripTime"`
    }
    ```
  - Split-DB SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CircuitBreakerAudits (
        auditId TEXT PRIMARY KEY,
        targetDomain TEXT NOT NULL,
        state TEXT NOT NULL,
        failureCount INTEGER NOT NULL,
        backoffSeconds INTEGER NOT NULL,
        isOpen INTEGER NOT NULL DEFAULT 0,
        isHalfOpen INTEGER NOT NULL DEFAULT 0,
        isClosed INTEGER NOT NULL DEFAULT 1,
        hasTripped INTEGER NOT NULL DEFAULT 0,
        updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX IF NOT EXISTS uq_circuit_domain ON CircuitBreakerAudits(targetDomain);
    ```
- **Antigravity / GitMap Integration Point:**
  - Encapsulated in `pkg/circuitbreaker/mx_breaker.go` and surfaced through `cli/cmdemail/circuit_breaker_cmd.go`.

---

## Technical & Architectural Governance Verification

- **Path Hygiene:** All paths are 100% relative repository paths (`02-spec/...`, `cli/cmdemail/...`, `pkg/...`). No drive letters (`d:`, `C:`) or absolute paths are present.
- **Positive Booleans Only:** Every boolean field strictly begins with `is` or `has` (`isOptimal`, `hasClickDrop`, `isHumanGrade`, `hasBotMarker`, `isThrottled`, `hasOverride`, `isPeakHour`, `hasTraffic`, `isQuarantined`, `hasOperatorNotif`, `isDesktop`, `isMobile`, `hasDarkTheme`, `isSecureTLS`, `hasCertError`, `isClean`, `hasLoops`, `isPositive`, `hasStatSig`, `isAutoReply`, `hasOptOutIntent`, `isGraduated`, `hasRiskFactors`, `isRead`, `isSpam`, `hasWarmupTag`, `isFocused`, `hasNextPage`, `isEnabled`, `hasParams`, `isSuccess`, `hasError`, `isAutoReplied`, `isStarred`, `hasFailures`, `isCompliant`, `hasMimeError`, `isSanitized`, `hasImages`, `hasReply`, `isDispatched`, `hasReferences`, `isSpecialUse`, `hasUnread`, `isBodyMatch`, `hasHeaderMatch`, `isComplete`, `hasMatches`, `isImapHealthy`, `isSmtpHealthy`, `isTokenValid`, `hasQuotaAlert`, `isFull`, `hasBackpressure`, `hasCapacity`, `hasSpills`, `isEngaged`, `hasRecovered`, `isReplayable`, `hasPermanentErr`, `hasSyncPending`, `isStarved`, `hasActiveFlow`, `isBalanced`, `hasCongestion`, `isScalingUp`, `hasErrorSpike`, `isCompacted`, `hasHighWaste`, `isBottleneck`, `isDelivered`, `hasDelays`, `isOpen`, `isHalfOpen`, `isClosed`, `hasTripped`).
- **Split-DB SQLite Compliance:** All tables are modeled with clean PascalCase table names, standard SQLite types, primary keys, and index definitions in `gitmap.db`.
- **Zero-CDN Architecture:** Terminal UIs rely purely on Bubble Tea and Lipgloss; web views utilize inline SVG and vanilla CSS with no external network dependencies.
