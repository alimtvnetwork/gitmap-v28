# Improvements 1301 to 1330: Session Pipelining, Socket Health & Jitter Waveform TUI

This document details the complete specifications for Improvements 1301 through 1330 across Clusters 131, 132, and 133 of Sequence 350.

---

# Cluster 131: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (1301–1310)

---

### Improvement 1301: RFC 2920 Command Streaming Pipeline & Zero-RTT Command Burst Engine
- **ID**: `IMP-1301`
- **Cluster**: 131 (Autonomous SMTP Session Pipelining & TCP Zero-Window Detection)
- **Goal**: Eliminate lockstep request-response latency by packing `RSET`, `MAIL FROM`, and `RCPT TO` commands into a single TCP payload burst, demultiplexing asynchronous multi-line server responses in zero-copy userspace memory.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss split-viewport displaying command transmission bursts in electric cyan (`#00E5FF`) alongside demuxed server return codes (`250 OK`) in spring green (`#00E676`).
  - *Web*: Custom Web Component `<smtp-pipeline-stream>` rendering real-time streaming timeline bars with CSS transitions and SVG status badges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline stream --session-id=<id> [--burst-size=<size>] [--is-pipeline-forced] [--has-early-data]`
  - Flag `--session-id`: Identifier of the target warming SMTP session.
  - Flag `--burst-size`: Maximum number of commands coalesced per outbound TCP frame (default: `16`).
  - Flag `--is-pipeline-forced`: Bypasses EHLO capability check when remote server supports RFC 2920.
  - Flag `--has-early-data`: Enables TLS 1.3 0-RTT early data command payload transmission.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpPipelineStreams (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      burstCommandCount INTEGER NOT NULL DEFAULT 0,
      burstDurationMs REAL NOT NULL DEFAULT 0.0,
      isPipelineForced INTEGER NOT NULL DEFAULT 0,
      hasEarlyData INTEGER NOT NULL DEFAULT 0,
      streamedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpPipelineStream struct {
      ID                string    `json:"id"`
      SessionID         string    `json:"sessionId"`
      BurstCommandCount int       `json:"burstCommandCount"`
      BurstDurationMs   float64   `json:"burstDurationMs"`
      IsPipelineForced  bool      `json:"isPipelineForced"`
      HasEarlyData      bool      `json:"hasEarlyData"`
      StreamedAt        time.Time `json:"streamedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/smtp/pipeline.go` in the socket command dispatcher.

---

### Improvement 1302: Real-Time Zero-Copy TCP Receive Window Monitor via Socket Reflection
- **ID**: `IMP-1302`
- **Cluster**: 131
- **Goal**: Continuously query OS kernel socket buffers using zero-copy `getsockopt` (`TCP_INFO`) to track advertised receiver window size (`snd_wnd`) in microseconds, preventing outbound buffer congestion.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status card showing real-time receive window capacity. Uses bright emerald (`#00E676`) for windows > 64KB, yellow amber (`#FFD600`) for 8–64KB, and vibrant red (`#FF1744`) for < 8KB.
  - *Web*: `<tcp-window-gauge>` Web Component displaying an SVG dial meter reflecting current remote buffer headroom.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline monitor-window --socket-fd=<fd> [--sample-rate-us=<us>] [--is-live-trace] [--has-kernel-probe]`
  - Flag `--socket-fd`: Active socket file descriptor under telemetry monitoring.
  - Flag `--sample-rate-us`: Sampling interval in microseconds (default: `500`).
  - Flag `--is-live-trace`: Streams live kernel window sizes directly to standard output.
  - Flag `--has-kernel-probe`: Hooks eBPF tracepoint for kernel TCP buffer tracking.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpWindowMetrics (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      advertisedWindowBytes INTEGER NOT NULL DEFAULT 0,
      sndBufferOccupancy INTEGER NOT NULL DEFAULT 0,
      isWindowExhausted INTEGER NOT NULL DEFAULT 0,
      hasKernelProbe INTEGER NOT NULL DEFAULT 1,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpWindowMetric struct {
      ID                    string    `json:"id"`
      SocketFd              int       `json:"socketFd"`
      AdvertisedWindowBytes int       `json:"advertisedWindowBytes"`
      SndBufferOccupancy    int       `json:"sndBufferOccupancy"`
      IsWindowExhausted     bool      `json:"isWindowExhausted"`
      HasKernelProbe        bool      `json:"hasKernelProbe"`
      MeasuredAt            time.Time `json:"measuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/tcp/window_monitor.go` prior to flushing pipelined buffers.

---

### Improvement 1303: Autonomous TCP Zero-Window Detection & Preemptive Backoff Controller
- **ID**: `IMP-1303`
- **Cluster**: 131
- **Goal**: Detect when remote ISP MTA advertises a zero-length receive window (`snd_wnd == 0`), automatically halting pipelined envelope transmission and engaging microsecond backoff timers without dropping the TCP connection.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss banner rendering a pulsing magenta alert (`#E040FB`) signaling zero-window state, with an ASCII countdown timer for next drain probe.
  - *Web*: `<zero-window-banner>` Web Component providing interactive backoff status with zero CDN dependencies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline zero-window-backoff --session-id=<id> [--probe-interval-ms=<ms>] [--is-aggressive-probe] [--has-exponential-decay]`
  - Flag `--session-id`: Active SMTP session ID experiencing zero-window condition.
  - Flag `--probe-interval-ms`: Initial window probe interval (default: `250`).
  - Flag `--is-aggressive-probe`: Sends single-byte window probe segments to prompt receiver window updates.
  - Flag `--has-exponential-decay`: Applies exponential backoff modifier to subsequent probe intervals.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowEvents (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      backoffDurationMs INTEGER NOT NULL DEFAULT 0,
      probeCount INTEGER NOT NULL DEFAULT 0,
      isResolved INTEGER NOT NULL DEFAULT 0,
      hasProbePayload INTEGER NOT NULL DEFAULT 1,
      detectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowEvent struct {
      ID                string    `json:"id"`
      SessionID         string    `json:"sessionId"`
      BackoffDurationMs int       `json:"backoffDurationMs"`
      ProbeCount        int       `json:"probeCount"`
      IsResolved        bool      `json:"isResolved"`
      HasProbePayload   bool      `json:"hasProbePayload"`
      DetectedAt        time.Time `json:"detectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/tcp/backoff_controller.go`.

---

### Improvement 1304: Multi-Stage Pipelined EHLO/STARTTLS Capability Handshake Negotiator
- **ID**: `IMP-1304`
- **Cluster**: 131
- **Goal**: Pipeline initial SMTP greeting, capability discovery, and TLS negotiation into concurrent execution phases, reducing connection setup overhead by 60%.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss table presenting handshake stages (TCP Syn, Banner, EHLO, STARTTLS, TLS Ready) with nanosecond latency badges in violet (`#7C4DFF`).
  - *Web*: `<handshake-timeline>` Web Component with interactive step progression and SVG connection graph.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline negotiate-tls --target-host=<host> [--tls-version=<ver>] [--is-strict-validation] [--has-session-ticket]`
  - Flag `--target-host`: Hostname or IP of the destination MTA.
  - Flag `--tls-version`: Minimum required TLS version (e.g., `1.3`).
  - Flag `--is-strict-validation`: Enforces strict X.509 certificate and DANE-TLSA validation.
  - Flag `--has-session-ticket`: Leverages TLS session resumption tickets for instantaneous 0-RTT handshakes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HandshakeNegotiations (
      id TEXT PRIMARY KEY,
      targetHost TEXT NOT NULL,
      tlsVersion TEXT NOT NULL,
      handshakeLatencyMs REAL NOT NULL DEFAULT 0.0,
      isStrictValidation INTEGER NOT NULL DEFAULT 1,
      hasSessionTicket INTEGER NOT NULL DEFAULT 0,
      negotiatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HandshakeNegotiation struct {
      ID                 string    `json:"id"`
      TargetHost         string    `json:"targetHost"`
      TlsVersion         string    `json:"tlsVersion"`
      HandshakeLatencyMs float64   `json:"handshakeLatencyMs"`
      IsStrictValidation bool      `json:"isStrictValidation"`
      HasSessionTicket   bool      `json:"hasSessionTicket"`
      NegotiatedAt       time.Time `json:"negotiatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/smtp/handshake.go`.

---

### Improvement 1305: Sliding Send Buffer Throttling for Outbound Socket Pressure Normalization
- **ID**: `IMP-1305`
- **Cluster**: 131
- **Goal**: Regulate userland send buffer writes against kernel socket queue pressure, dynamically pacing batch pipelined chunks to avoid triggering receiver buffer clampdowns.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss progress gauge rendering buffer saturation percentage in gradient from teal (`#00B0FF`) to amber (`#FF9100`).
  - *Web*: `<socket-pressure-gauge>` Web Component rendering real-time queue pressure bars with SVG markers.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline pace-send --socket-fd=<fd> [--max-pressure-pct=<pct>] [--is-dynamic-clamp] [--has-drain-notifier]`
  - Flag `--socket-fd`: File descriptor for the active outbound socket.
  - Flag `--max-pressure-pct`: Maximum buffer saturation threshold before pacing triggers (default: `80`).
  - Flag `--is-dynamic-clamp`: Adapts maximum threshold dynamically based on recent RTT oscillations.
  - Flag `--has-drain-notifier`: Registers kernel `POLLOUT` notifications for instant resumption upon queue drain.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketPressureStats (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      currentPressurePct INTEGER NOT NULL DEFAULT 0,
      pacingDelayUs INTEGER NOT NULL DEFAULT 0,
      isDynamicClamp INTEGER NOT NULL DEFAULT 1,
      hasDrainNotifier INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketPressureStat struct {
      ID                 string    `json:"id"`
      SocketFd           int       `json:"socketFd"`
      CurrentPressurePct int       `json:"currentPressurePct"`
      PacingDelayUs      int       `json:"pacingDelayUs"`
      IsDynamicClamp     bool      `json:"isDynamicClamp"`
      HasDrainNotifier   bool      `json:"hasDrainNotifier"`
      RecordedAt         time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tcp/pacer.go`.

---

### Improvement 1306: Asynchronous Response Demultiplexer & RFC 2920 Multi-Line Parser
- **ID**: `IMP-1306`
- **Cluster**: 131
- **Goal**: Efficiently parse streamed multi-line SMTP responses (e.g., `250-SIZE`, `250-PIPELINING`, `250 OK`) matching them to previously batched pipelined requests with zero allocations using slice windows.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss dual-pane view showing batched outbound commands on the left and mapped server response status codes on the right in muted gray (`#9E9E9E`) and neon cyan (`#00E5FF`).
  - *Web*: `<response-demux-inspector>` Web Component providing expandable command-response pairs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline parse-responses --session-id=<id> [--buffer-size-bytes=<bytes>] [--is-zero-alloc] [--has-multiline-folding]`
  - Flag `--session-id`: SMTP warming session identifier.
  - Flag `--buffer-size-bytes`: Internal memory ring buffer size for streaming tokenization (default: `8192`).
  - Flag `--is-zero-alloc`: Enforces zero-allocation string parsing utilizing byte slice windows.
  - Flag `--has-multiline-folding`: Folds continuation lines into unified structured tokens.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DemuxedResponses (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      commandVerb TEXT NOT NULL,
      statusCode INTEGER NOT NULL,
      isZeroAlloc INTEGER NOT NULL DEFAULT 1,
      hasMultilineFolding INTEGER NOT NULL DEFAULT 1,
      parsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DemuxedResponse struct {
      ID                  string    `json:"id"`
      SessionID           string    `json:"sessionId"`
      CommandVerb         string    `json:"commandVerb"`
      StatusCode          int       `json:"statusCode"`
      IsZeroAlloc         bool      `json:"isZeroAlloc"`
      HasMultilineFolding bool      `json:"hasMultilineFolding"`
      ParsedAt            time.Time `json:"parsedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/smtp/demux.go`.

---

### Improvement 1307: Dynamic Window Size Prediction via Exponential Moving Average (EMA)
- **ID**: `IMP-1307`
- **Cluster**: 131
- **Goal**: Predict impending zero-window clampdowns by tracking the first and second derivatives of receiver advertised window decay across consecutive pipelined bursts using an exponential moving average.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss line chart visualizing EMA projected window headroom with warning threshold markers in safety orange (`#FF6D00`).
  - *Web*: `<window-ema-chart>` Web Component rendering real-time SVG predictive decay curves.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline predict-window --session-id=<id> [--smoothing-alpha=<val>] [--is-predictive-drain] [--has-decay-alert]`
  - Flag `--session-id`: Target SMTP warming session identifier.
  - Flag `--smoothing-alpha`: EMA weight factor between 0.01 and 0.99 (default: `0.25`).
  - Flag `--is-predictive-drain`: Automatically throttles burst dispatch before window reaches 0.
  - Flag `--has-decay-alert`: Fires terminal event trigger when window decay velocity exceeds threshold.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WindowEmaPredictions (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      predictedWindowBytes INTEGER NOT NULL DEFAULT 0,
      decayVelocity REAL NOT NULL DEFAULT 0.0,
      isPredictiveDrain INTEGER NOT NULL DEFAULT 1,
      hasDecayAlert INTEGER NOT NULL DEFAULT 0,
      calculatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WindowEmaPrediction struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      PredictedWindowBytes int       `json:"predictedWindowBytes"`
      DecayVelocity        float64   `json:"decayVelocity"`
      IsPredictiveDrain    bool      `json:"isPredictiveDrain"`
      HasDecayAlert        bool      `json:"hasDecayAlert"`
      CalculatedAt         time.Time `json:"calculatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/tcp/window_predictor.go`.

---

### Improvement 1308: Out-of-Order Pipelined RCPT Failure Isolation & Auto-RSET Recovery
- **ID**: `IMP-1308`
- **Cluster**: 131
- **Goal**: Gracefully isolate individual recipient rejections (`550 Mailbox Not Found`) within a pipelined batch without corrupting subsequent valid recipients, auto-dispatching `RSET` only when fatal session errors occur.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss checklist showing individual recipient validation badges; rejected addresses tagged in crimson (`#FF1744`), accepted in lime green (`#76FF03`).
  - *Web*: `<recipient-batch-table>` Web Component displaying isolation status with SVG state indicators.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline isolate-rcpt --session-id=<id> [--is-auto-rset] [--has-quarantine-redirect]`
  - Flag `--session-id`: Active SMTP session identifier.
  - Flag `--is-auto-rset`: Emits atomic `RSET` command when all recipients within batch fail.
  - Flag `--has-quarantine-redirect`: Automatically shifts invalid recipient seeds to isolation queue.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RcptPipelineResults (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      recipientEmail TEXT NOT NULL,
      smtpStatusCode INTEGER NOT NULL,
      isAutoRset INTEGER NOT NULL DEFAULT 0,
      hasQuarantineRedirect INTEGER NOT NULL DEFAULT 1,
      processedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RcptPipelineResult struct {
      ID                    string    `json:"id"`
      SessionID             string    `json:"sessionId"`
      RecipientEmail        string    `json:"recipientEmail"`
      SmtpStatusCode        int       `json:"smtpStatusCode"`
      IsAutoRset            bool      `json:"isAutoRset"`
      HasQuarantineRedirect bool      `json:"hasQuarantineRedirect"`
      ProcessedAt           time.Time `json:"processedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/smtp/rcpt_guard.go`.

---

### Improvement 1309: High-Throughput Zero-Copy Chunked DATA Streaming Engine
- **ID**: `IMP-1309`
- **Cluster**: 131
- **Goal**: Stream large email bodies and headers directly from memory-mapped files to kernel TCP buffers via `sendfile`/`splice` syscalls after pipelined `DATA` response confirmation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss transfer monitor showing real-time streaming throughput in MB/s in bright sapphire (`#2979FF`).
  - *Web*: `<zero-copy-stream-card>` Web Component displaying byte offset progression with animated SVG bars.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline stream-data --session-id=<id> --payload-path=<path> [--is-splice-enabled] [--has-dot-stuffing]`
  - Flag `--session-id`: Target SMTP session identifier.
  - Flag `--payload-path`: Relative path to raw RFC 5322 MIME envelope payload.
  - Flag `--is-splice-enabled`: Utilizes zero-copy Linux `splice` or Windows `TransmitFile` primitives.
  - Flag `--has-dot-stuffing`: Pre-calculates transparent dot-stuffing without runtime string reallocations.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroCopyDataStreams (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      totalBytesStreamed INTEGER NOT NULL DEFAULT 0,
      throughputMbps REAL NOT NULL DEFAULT 0.0,
      isSpliceEnabled INTEGER NOT NULL DEFAULT 1,
      hasDotStuffing INTEGER NOT NULL DEFAULT 1,
      completedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroCopyDataStream struct {
      ID                 string    `json:"id"`
      SessionID          string    `json:"sessionId"`
      TotalBytesStreamed int64     `json:"totalBytesStreamed"`
      ThroughputMbps     float64   `json:"throughputMbps"`
      IsSpliceEnabled    bool      `json:"isSpliceEnabled"`
      HasDotStuffing     bool      `json:"hasDotStuffing"`
      CompletedAt        time.Time `json:"completedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/smtp/data_stream.go`.

---

### Improvement 1310: Autonomous Session Pipeline Diagnostics & TCP Metrics Export
- **ID**: `IMP-1310`
- **Cluster**: 131
- **Goal**: Aggregate end-to-end pipelining metrics, socket buffer occupancy histograms, and zero-window incident logs into a unified structured JSON/SQLite telemetry report.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss formatted diagnostic summary card with tabular pipeline execution efficiency scores.
  - *Web*: `<pipeline-diagnostic-report>` Web Component with downloadable JSON summary and interactive metrics dashboard.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline export-metrics [--session-id=<id>] [--format=<json|text>] [--is-detailed] [--has-window-histogram]`
  - Flag `--session-id`: Specific session filter (omit for all active sessions).
  - Flag `--format`: Output format choice (`json` or `text`).
  - Flag `--is-detailed`: Includes nanosecond socket state transitions in export payload.
  - Flag `--has-window-histogram`: Appends 64-bin TCP receive window distribution histogram.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PipelineAuditReports (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      efficiencyScore REAL NOT NULL DEFAULT 0.0,
      zeroWindowIncidents INTEGER NOT NULL DEFAULT 0,
      isDetailed INTEGER NOT NULL DEFAULT 0,
      hasWindowHistogram INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PipelineAuditReport struct {
      ID                  string    `json:"id"`
      SessionID           string    `json:"sessionId"`
      EfficiencyScore     float64   `json:"efficiencyScore"`
      ZeroWindowIncidents int       `json:"zeroWindowIncidents"`
      IsDetailed          bool      `json:"isDetailed"`
      HasWindowHistogram  bool      `json:"hasWindowHistogram"`
      ExportedAt          time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/telemetry/exporter.go`.

---

# Cluster 132: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (1311–1320)

---

### Improvement 1311: Persistent TLS Connection Reuse Pool with ISP Route Affinity
- **ID**: `IMP-1311`
- **Cluster**: 132 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal**: Maintain warmed, pre-authenticated TLS SMTP sockets grouped by destination MX provider, reducing handshakes per batch while respecting per-IP concurrency caps.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss multi-column table displaying ISP connection pools (Google, Microsoft, Yahoo, Fastmail) with active/idle connection counts in neon cyan (`#00E5FF`).
  - *Web*: `<isp-pool-matrix>` Web Component with live connection counters and responsive SVG status badges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool status [--isp=<name>] [--is-idle-visible] [--has-health-metrics]`
  - Flag `--isp`: Filter by specific target provider (e.g., `google`, `microsoft`).
  - Flag `--is-idle-visible`: Displays idle persistent connections waiting in pool.
  - Flag `--has-health-metrics`: Appends round-trip latency and packet loss rates to status table.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspConnectionPools (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      activeConnections INTEGER NOT NULL DEFAULT 0,
      idleConnections INTEGER NOT NULL DEFAULT 0,
      isIdleVisible INTEGER NOT NULL DEFAULT 1,
      hasHealthMetrics INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspConnectionPool struct {
      ID                string    `json:"id"`
      IspName           string    `json:"ispName"`
      ActiveConnections int       `json:"activeConnections"`
      IdleConnections   int       `json:"idleConnections"`
      IsIdleVisible     bool      `json:"isIdleVisible"`
      HasHealthMetrics  bool      `json:"hasHealthMetrics"`
      UpdatedAt         time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Core pool manager in `pkg/warmup/pool/manager.go`.

---

### Improvement 1312: Calibrated SMTP Keep-Alive Heartbeat Engine via NOOP/RSET
- **ID**: `IMP-1312`
- **Cluster**: 132
- **Goal**: Prevent stateful firewall idle connection termination by scheduling dynamic `NOOP` and `RSET` heartbeats calibrated precisely below ISP TCP idle timeout thresholds.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss pulse indicator showing heartbeat cadence countdown for each idle socket in soft green (`#69F0AE`).
  - *Web*: `<heartbeat-cadence-meter>` Web Component rendering circular animated SVG pulse rings.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool tune-heartbeat --isp=<name> [--interval-sec=<sec>] [--is-noop-preferred] [--has-jitter]`
  - Flag `--isp`: Target ISP identifier.
  - Flag `--interval-sec`: Keep-alive interval in seconds (default: `180`).
  - Flag `--is-noop-preferred`: Sends `NOOP` instead of `RSET` as primary keep-alive ping.
  - Flag `--has-jitter`: Introduces random ±15% timing jitter to avoid synchronized fleet heartbeat spikes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ConnectionHeartbeats (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      heartbeatIntervalSec INTEGER NOT NULL DEFAULT 180,
      heartbeatCount INTEGER NOT NULL DEFAULT 0,
      isNoopPreferred INTEGER NOT NULL DEFAULT 1,
      hasJitter INTEGER NOT NULL DEFAULT 1,
      dispatchedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ConnectionHeartbeat struct {
      ID                   string    `json:"id"`
      IspName              string    `json:"ispName"`
      HeartbeatIntervalSec int       `json:"heartbeatIntervalSec"`
      HeartbeatCount       int       `json:"heartbeatCount"`
      IsNoopPreferred      bool      `json:"isNoopPreferred"`
      HasJitter            bool      `json:"hasJitter"`
      DispatchedAt         time.Time `json:"dispatchedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/pool/heartbeat.go`.

---

### Improvement 1313: Dynamic Socket Health Matrix & Multi-Factor Degradation Scoring
- **ID**: `IMP-1313`
- **Cluster**: 132
- **Goal**: Score individual socket health continuously based on SYN-ACK latency, TLS renegotiation frequency, and TCP segment retransmissions, culling degraded sockets before dispatch errors occur.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss matrix grid rendering socket health scores (0–100) color-coded from emerald (`#00E676`) down to scarlet (`#D50000`).
  - *Web*: `<socket-health-matrix>` Web Component with interactive cell tooltips and live SVG score heatmaps.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool score-health [--min-score=<score>] [--is-auto-prune] [--has-retransmit-penalty]`
  - Flag `--min-score`: Minimum threshold score (0–100) below which sockets are retired (default: `65`).
  - Flag `--is-auto-prune`: Closes and evicts unhealthy sockets automatically.
  - Flag `--has-retransmit-penalty`: Weights kernel TCP retransmission counters heavily in score computation.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketHealthMatrices (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      healthScore REAL NOT NULL DEFAULT 100.0,
      retransmitCount INTEGER NOT NULL DEFAULT 0,
      isAutoPrune INTEGER NOT NULL DEFAULT 1,
      hasRetransmitPenalty INTEGER NOT NULL DEFAULT 1,
      evaluatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketHealthMatrix struct {
      ID                   string    `json:"id"`
      SocketFd             int       `json:"socketFd"`
      HealthScore          float64   `json:"healthScore"`
      RetransmitCount      int       `json:"retransmitCount"`
      IsAutoPrune          bool      `json:"isAutoPrune"`
      HasRetransmitPenalty bool      `json:"hasRetransmitPenalty"`
      EvaluatedAt          time.Time `json:"evaluatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/pool/health.go`.

---

### Improvement 1314: Multi-MTA Concurrent Handshake Pre-Warming & Speculative Dialing
- **ID**: `IMP-1314`
- **Cluster**: 132
- **Goal**: Speculatively dial and pre-warm backup TCP/TLS sessions across secondary MX targets before primary connection saturation, eliminating tail latency spikes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss banner showing speculative dialing queue with live handshake completion progress in violet (`#7C4DFF`).
  - *Web*: `<speculative-dial-card>` Web Component displaying pre-warmed socket ready states.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool prewarm --domain=<domain> [--concurrency=<num>] [--is-tls-prewarmed] [--has-speculative-bound]`
  - Flag `--domain`: Target recipient domain (e.g. `gmail.com`).
  - Flag `--concurrency`: Number of speculative connections to initiate (default: `2`).
  - Flag `--is-tls-prewarmed`: Completes full TLS handshake during speculative pre-warm.
  - Flag `--has-speculative-bound`: Caps total speculative sockets to 20% of pool size.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PrewarmedConnections (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      prewarmCount INTEGER NOT NULL DEFAULT 0,
      isTlsPrewarmed INTEGER NOT NULL DEFAULT 1,
      hasSpeculativeBound INTEGER NOT NULL DEFAULT 1,
      initiatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PrewarmedConnection struct {
      ID                  string    `json:"id"`
      TargetDomain        string    `json:"targetDomain"`
      PrewarmCount        int       `json:"prewarmCount"`
      IsTlsPrewarmed      bool      `json:"isTlsPrewarmed"`
      HasSpeculativeBound bool      `json:"hasSpeculativeBound"`
      InitiatedAt         time.Time `json:"initiatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/prewarm.go`.

---

### Improvement 1315: Multi-Interface Egress IP Binding & Local Route Failover
- **ID**: `IMP-1315`
- **Cluster**: 132
- **Goal**: Bind outbound sockets across multiple local network interfaces and dedicated IP addresses, rotating egress paths and executing instant failover on local interface saturation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss interface card displaying local IP bindings with egress bandwidth meters in cyan (`#00E5FF`).
  - *Web*: `<interface-egress-table>` Web Component showing real-time interface status and route health.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool bind-interfaces --ips=<ip1,ip2> [--is-round-robin] [--has-interface-failover]`
  - Flag `--ips`: Comma-separated list of local egress IP addresses.
  - Flag `--is-round-robin`: Rotates outgoing socket bindings evenly across provided IPs.
  - Flag `--has-interface-failover`: Automatically reroutes connections away from failed local interfaces.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS InterfaceBindings (
      id TEXT PRIMARY KEY,
      ipAddress TEXT NOT NULL,
      egressPacketCount INTEGER NOT NULL DEFAULT 0,
      isRoundRobin INTEGER NOT NULL DEFAULT 1,
      hasInterfaceFailover INTEGER NOT NULL DEFAULT 1,
      boundAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type InterfaceBinding struct {
      ID                   string    `json:"id"`
      IPAddress            string    `json:"ipAddress"`
      EgressPacketCount    int64     `json:"egressPacketCount"`
      IsRoundRobin         bool      `json:"isRoundRobin"`
      HasInterfaceFailover bool      `json:"hasInterfaceFailover"`
      BoundAt              time.Time `json:"boundAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Injected into `pkg/warmup/pool/dialer.go`.

---

### Improvement 1316: Cross-ISP Latency Baseline Tracker & Moving Median Calibrator
- **ID**: `IMP-1316`
- **Cluster**: 132
- **Goal**: Continuously record TCP handshake and command response RTTs across ISP endpoints, maintaining a rolling median baseline to identify ISP-wide network anomalies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss sparkline visualizer rendering rolling RTT medians for top 5 ISPs with baseline threshold lines.
  - *Web*: `<isp-latency-sparkline>` Web Component providing interactive SVG trendlines.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool track-latency [--window-size=<samples>] [--is-median-smoothed] [--has-anomaly-alarm]`
  - Flag `--window-size`: Number of recent latency samples included in moving window (default: `100`).
  - Flag `--is-median-smoothed`: Uses median filtering instead of mean to eliminate outlier bias.
  - Flag `--has-anomaly-alarm`: Generates alert triggers when median RTT exceeds 2.5x baseline.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspLatencyBaselines (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      medianRttMs REAL NOT NULL DEFAULT 0.0,
      sampleWindowCount INTEGER NOT NULL DEFAULT 100,
      isMedianSmoothed INTEGER NOT NULL DEFAULT 1,
      hasAnomalyAlarm INTEGER NOT NULL DEFAULT 0,
      calculatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspLatencyBaseline struct {
      ID                string    `json:"id"`
      IspName           string    `json:"ispName"`
      MedianRttMs       float64   `json:"medianRttMs"`
      SampleWindowCount int       `json:"sampleWindowCount"`
      IsMedianSmoothed  bool      `json:"isMedianSmoothed"`
      HasAnomalyAlarm   bool      `json:"hasAnomalyAlarm"`
      CalculatedAt      time.Time `json:"calculatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/latency_tracker.go`.

---

### Improvement 1317: Autonomous Connection Rebalancing & Egress Drain Orchestrator
- **ID**: `IMP-1317`
- **Cluster**: 132
- **Goal**: Smoothly drain and rebalance connection pools when ISP rate limiters dynamically contract allowed concurrency, closing idle sockets gracefully with `QUIT`.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status badge showing drain state (`DRAINING`, `BALANCED`) with active socket drain counter in amber (`#FFD600`).
  - *Web*: `<connection-drain-panel>` Web Component with animated progress indicators and socket state table.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool rebalance --isp=<name> --target-concurrency=<num> [--is-graceful-quit] [--has-spool-pause]`
  - Flag `--isp`: Target ISP identifier undergoing rebalance.
  - Flag `--target-concurrency`: Desired maximum concurrent socket threshold.
  - Flag `--is-graceful-quit`: Emits standard SMTP `QUIT` command on drained sockets.
  - Flag `--has-spool-pause`: Temporarily pauses outbound spool dispatch during rebalance transition.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ConnectionRebalances (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      targetConcurrency INTEGER NOT NULL,
      drainedSocketCount INTEGER NOT NULL DEFAULT 0,
      isGracefulQuit INTEGER NOT NULL DEFAULT 1,
      hasSpoolPause INTEGER NOT NULL DEFAULT 0,
      rebalancedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ConnectionRebalance struct {
      ID                 string    `json:"id"`
      IspName            string    `json:"ispName"`
      TargetConcurrency  int       `json:"targetConcurrency"`
      DrainedSocketCount int       `json:"drainedSocketCount"`
      IsGracefulQuit     bool      `json:"isGracefulQuit"`
      HasSpoolPause      bool      `json:"hasSpoolPause"`
      RebalancedAt       time.Time `json:"rebalancedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/pool/rebalancer.go`.

---

### Improvement 1318: TCP Fast Open (TFO) Connection Booster for SMTP Sockets
- **ID**: `IMP-1318`
- **Cluster**: 132
- **Goal**: Enable TCP Fast Open (`MSG_FASTOPEN`) on supported Linux and Windows kernels to bundle initial data in the SYN packet, saving one full RTT during SMTP reconnects.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status tag showing `TFO: ENABLED` in neon green (`#00E676`) with cookie caching status.
  - *Web*: `<tfo-status-badge>` Web Component rendering interactive SVG hardware acceleration badges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool enable-tfo [--is-system-validated] [--has-cookie-cache]`
  - Flag `--is-system-validated`: Checks kernel `/proc/sys/net/ipv4/tcp_fastopen` or Windows TFO support flags.
  - Flag `--has-cookie-cache`: Maintains local fast-open binary cookie cache per target MX IP.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpFastOpenStats (
      id TEXT PRIMARY KEY,
      targetIp TEXT NOT NULL,
      tfoCookieHex TEXT,
      synDataBytes INTEGER NOT NULL DEFAULT 0,
      isSystemValidated INTEGER NOT NULL DEFAULT 1,
      hasCookieCache INTEGER NOT NULL DEFAULT 1,
      testedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpFastOpenStat struct {
      ID                string    `json:"id"`
      TargetIP          string    `json:"targetIp"`
      TfoCookieHex      string    `json:"tfoCookieHex"`
      SynDataBytes      int       `json:"synDataBytes"`
      IsSystemValidated bool      `json:"isSystemValidated"`
      HasCookieCache    bool      `json:"hasCookieCache"`
      TestedAt          time.Time `json:"testedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/tcp/tfo.go`.

---

### Improvement 1319: Multi-ISP Connection Pool Quarantine & Circuit Breaker Engine
- **ID**: `IMP-1319`
- **Cluster**: 132
- **Goal**: Automatically trip circuit breaker when connection failures to a specific ISP exceed error thresholds, placing the pool into quarantine and diverting warm-up cohorts to non-impacted providers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss circuit breaker HUD showing `TRIPPED: QUARANTINE` in bold scarlet (`#FF1744`) with countdown reset timer.
  - *Web*: `<circuit-breaker-card>` Web Component displaying interactive breaker states and trip logs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool circuit-breaker --isp=<name> [--trip-threshold=<num>] [--is-auto-reset] [--has-alert-dispatch]`
  - Flag `--isp`: Target ISP identifier.
  - Flag `--trip-threshold`: Consecutive error threshold before circuit trips (default: `5`).
  - Flag `--is-auto-reset`: Automatically attempts canary reconnection after cooldown window expires.
  - Flag `--has-alert-dispatch`: Dispatches event notification to central audit ledger.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CircuitBreakerEvents (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      consecutiveErrors INTEGER NOT NULL DEFAULT 0,
      quarantineSeconds INTEGER NOT NULL DEFAULT 300,
      isAutoReset INTEGER NOT NULL DEFAULT 1,
      hasAlertDispatch INTEGER NOT NULL DEFAULT 1,
      trippedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CircuitBreakerEvent struct {
      ID                string    `json:"id"`
      IspName           string    `json:"ispName"`
      ConsecutiveErrors int       `json:"consecutiveErrors"`
      QuarantineSeconds int       `json:"quarantineSeconds"`
      IsAutoReset       bool      `json:"isAutoReset"`
      HasAlertDispatch  bool      `json:"hasAlertDispatch"`
      TrippedAt         time.Time `json:"trippedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Core resilience guard in `pkg/warmup/pool/circuit_breaker.go`.

---

### Improvement 1320: Connection Pool Historical Audit & Socket Churn Telemetry
- **ID**: `IMP-1320`
- **Cluster**: 132
- **Goal**: Export detailed historical socket churn metrics, connection reuse efficiencies, and handshake failure rates to append-only forensic audit storage.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss audit log table rendering historical churn events with timestamp markers in deep violet (`#7C4DFF`).
  - *Web*: `<pool-audit-table>` Web Component with searchable historical log records.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool export-audit [--since=<duration>] [--is-compressed] [--has-error-stack]`
  - Flag `--since`: Time window for historical query (e.g. `24h`).
  - Flag `--is-compressed`: Compresses export stream utilizing Gzip compression.
  - Flag `--has-error-stack`: Appends detailed error stack traces for failed socket creations.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolAuditLedgers (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      socketLifecycleEvent TEXT NOT NULL,
      reusedConnectionCount INTEGER NOT NULL DEFAULT 0,
      isCompressed INTEGER NOT NULL DEFAULT 0,
      hasErrorStack INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolAuditLedger struct {
      ID                    string    `json:"id"`
      IspName               string    `json:"ispName"`
      SocketLifecycleEvent  string    `json:"socketLifecycleEvent"`
      ReusedConnectionCount int       `json:"reusedConnectionCount"`
      IsCompressed          bool      `json:"isCompressed"`
      HasErrorStack         bool      `json:"hasErrorStack"`
      RecordedAt            time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/telemetry/pool_audit.go`.

---

# Cluster 133: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1321–1330)

---

### Improvement 1321: High-Definition Sub-Millisecond Latency Flamegraph TUI Engine
- **ID**: `IMP-1321`
- **Cluster**: 133 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Render real-time sub-millisecond flamegraphs of SMTP connection phases (DNS, TCP, TLS, EHLO, MAIL, RCPT, DATA) in terminal viewports using Lipgloss TrueColor formatting.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss hierarchical flamegraph bars in cascading colors (DNS: `#00E5FF`, TCP: `#7C4DFF`, TLS: `#E040FB`, SMTP: `#00E676`).
  - *Web*: `<latency-flamegraph-svg>` Web Component rendering interactive SVG vector flamegraphs with zero CDN dependencies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui flamegraph --session-id=<id> [--refresh-ms=<ms>] [--is-interactive] [--has-sub-millisecond]`
  - Flag `--session-id`: Target SMTP warming session identifier.
  - Flag `--refresh-ms`: Screen refresh interval in milliseconds (default: `100`).
  - Flag `--is-interactive`: Enables arrow-key drilldown navigation through flamegraph nodes.
  - Flag `--has-sub-millisecond`: Formats timing resolution in microseconds (`µs`).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS FlamegraphSamples (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      phaseName TEXT NOT NULL,
      durationUs INTEGER NOT NULL DEFAULT 0,
      isInteractive INTEGER NOT NULL DEFAULT 1,
      hasSubMillisecond INTEGER NOT NULL DEFAULT 1,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type FlamegraphSample struct {
      ID                string    `json:"id"`
      SessionID         string    `json:"sessionId"`
      PhaseName         string    `json:"phaseName"`
      DurationUs        int64     `json:"durationUs"`
      IsInteractive     bool      `json:"isInteractive"`
      HasSubMillisecond bool      `json:"hasSubMillisecond"`
      CapturedAt        time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/flamegraph.go`.

---

### Improvement 1322: Real-Time Bounded Jitter Waveform Oscilloscope Visualizer
- **ID**: `IMP-1322`
- **Cluster**: 133
- **Goal**: Render real-time TCP packet and ACK jitter as an ASCII/Unicode waveform oscilloscope in terminal viewports, highlighting microsecond phase jitter.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Unicode Braille/block-based waveform rendering oscillating signals in vivid cyan (`#00E5FF`) with trigger lines.
  - *Web*: `<jitter-waveform-scope>` Web Component providing responsive HTML5 Canvas oscilloscope with zero external scripts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui jitter-scope [--window-width=<cols>] [--is-auto-scale] [--has-phosphor-glow]`
  - Flag `--window-width`: Viewport column width for waveform rendering (default: `80`).
  - Flag `--is-auto-scale`: Dynamically scales vertical amplitude to fit peak jitter signals.
  - Flag `--has-phosphor-glow`: Renders visual decay trail simulating analog CRT phosphor screens.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterWaveformSamples (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      jitterValueUs REAL NOT NULL DEFAULT 0.0,
      amplitudeRatio REAL NOT NULL DEFAULT 1.0,
      isAutoScale INTEGER NOT NULL DEFAULT 1,
      hasPhosphorGlow INTEGER NOT NULL DEFAULT 1,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterWaveformSample struct {
      ID              string    `json:"id"`
      SocketFd        int       `json:"socketFd"`
      JitterValueUs   float64   `json:"jitterValueUs"`
      AmplitudeRatio  float64   `json:"amplitudeRatio"`
      IsAutoScale     bool      `json:"isAutoScale"`
      HasPhosphorGlow bool      `json:"hasPhosphorGlow"`
      SampledAt       time.Time `json:"sampledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/oscilloscope.go`.

---

### Improvement 1323: Interactive Bubble Tea v2 HUD Viewport & Multi-Window Layout
- **ID**: `IMP-1323`
- **Cluster**: 133
- **Goal**: Provide a responsive terminal dashboard using Bubble Tea v2 architecture, dividing the screen into collapsible tiles for connection pools, queue pressure, and active streams.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss framed windows with active tile highlighted in bright sapphire (`#2979FF`) and inactive tiles in muted slate (`#37474F`).
  - *Web*: `<tui-hud-layout>` Web Component mirroring terminal tiled grid with responsive CSS Grid layout.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui hud [--layout=<quad|split|focus>] [--is-mouse-enabled] [--has-fullscreen-toggle]`
  - Flag `--layout`: Dashboard panel arrangement mode (`quad`, `split`, `focus`).
  - Flag `--is-mouse-enabled`: Captures terminal mouse clicks and scroll events for panel selection.
  - Flag `--has-fullscreen-toggle`: Permits single-key expansion of focused panel (`F` key).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HudLayoutConfigs (
      id TEXT PRIMARY KEY,
      layoutMode TEXT NOT NULL DEFAULT 'quad',
      activePaneIndex INTEGER NOT NULL DEFAULT 0,
      isMouseEnabled INTEGER NOT NULL DEFAULT 1,
      hasFullscreenToggle INTEGER NOT NULL DEFAULT 1,
      savedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HudLayoutConfig struct {
      ID                  string    `json:"id"`
      LayoutMode          string    `json:"layoutMode"`
      ActivePaneIndex     int       `json:"activePaneIndex"`
      IsMouseEnabled      bool      `json:"isMouseEnabled"`
      HasFullscreenToggle bool      `json:"hasFullscreenToggle"`
      SavedAt             time.Time `json:"savedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Application runtime in `pkg/warmup/tui/hud_model.go`.

---

### Improvement 1324: Zero-CDN Autonomous SVG Exporter & Static Telemetry Bundler
- **ID**: `IMP-1324`
- **Cluster**: 133
- **Goal**: Generate fully self-contained SVG files containing embedded vector assets, CSS animations, and inline data payloads for offline reporting with zero CDN references.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss success card confirming export path and file size in kilobytes.
  - *Web*: Native standalone SVG document viewable in any modern browser without network access.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui export-svg --output-path=<path> [--is-animated] [--has-embedded-css]`
  - Flag `--output-path`: Target filesystem destination path for generated SVG document.
  - Flag `--is-animated`: Embeds CSS keyframe animations for live waveform motion.
  - Flag `--has-embedded-css`: Inlines CSS styling directly inside SVG `<style>` tags.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SvgExportRecords (
      id TEXT PRIMARY KEY,
      destinationPath TEXT NOT NULL,
      fileSizeBytes INTEGER NOT NULL DEFAULT 0,
      isAnimated INTEGER NOT NULL DEFAULT 1,
      hasEmbeddedCss INTEGER NOT NULL DEFAULT 1,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SvgExportRecord struct {
      ID              string    `json:"id"`
      DestinationPath string    `json:"destinationPath"`
      FileSizeBytes   int64     `json:"fileSizeBytes"`
      IsAnimated      bool      `json:"isAnimated"`
      HasEmbeddedCss  bool      `json:"hasEmbeddedCss"`
      GeneratedAt     time.Time `json:"generatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/exporter/svg_bundle.go`.

---

### Improvement 1325: Adaptive Lipgloss Color Theme Negotiator & ANSI Fallback Driver
- **ID**: `IMP-1325`
- **Cluster**: 133
- **Goal**: Automatically negotiate terminal color capabilities (TrueColor 24-bit, ANSI-256, ANSI-16, or monochrome), falling back gracefully without visual clipping or syntax errors.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual palette calibration bar displaying 16 test swatches with capability detection status.
  - *Web*: `<theme-negotiator-preview>` Web Component demonstrating rendering under simulated color-depth restrictions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui test-colors [--force-depth=<24bit|256|16|mono>] [--is-truecolor-checked] [--has-high-contrast]`
  - Flag `--force-depth`: Manually overrides terminal color depth detection.
  - Flag `--is-truecolor-checked`: Probes `$COLORTERM` and terminal query sequences for 24-bit capability.
  - Flag `--has-high-contrast`: Activates WCAG-compliant high-contrast color mappings.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ColorThemeSettings (
      id TEXT PRIMARY KEY,
      colorDepth TEXT NOT NULL DEFAULT '24bit',
      isTruecolorChecked INTEGER NOT NULL DEFAULT 1,
      hasHighContrast INTEGER NOT NULL DEFAULT 0,
      configuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ColorThemeSetting struct {
      ID                 string    `json:"id"`
      ColorDepth         string    `json:"colorDepth"`
      IsTruecolorChecked bool      `json:"isTruecolorChecked"`
      HasHighContrast    bool      `json:"hasHighContrast"`
      ConfiguredAt       time.Time `json:"configuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/tui/theme.go`.

---

### Improvement 1326: Real-Time Egress Queue Micro-Depth Sparkline Indicator
- **ID**: `IMP-1326`
- **Cluster**: 133
- **Goal**: Render inline terminal sparklines representing outbound spool queue depth in real-time, displaying backlog fluctuations per ISP.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Compact single-line Unicode sparkline (` ▂▃▄▅▆▇█`) in electric cyan (`#00E5FF`).
  - *Web*: `<queue-depth-sparkline>` Web Component rendering real-time SVG vector bar sparklines.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui queue-sparkline [--isp=<name>] [--is-realtime-stream] [--has-peak-indicator]`
  - Flag `--isp`: Filter sparkline telemetry to specific ISP.
  - Flag `--is-realtime-stream`: Continuously stream sparkline updates in terminal stdout.
  - Flag `--has-peak-indicator`: Highlights queue depth peaks with marker glyphs (`▲`).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS QueueSparklineSamples (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      queueDepth INTEGER NOT NULL DEFAULT 0,
      peakDepth INTEGER NOT NULL DEFAULT 0,
      isRealtimeStream INTEGER NOT NULL DEFAULT 1,
      hasPeakIndicator INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type QueueSparklineSample struct {
      ID               string    `json:"id"`
      IspName          string    `json:"ispName"`
      QueueDepth       int       `json:"queueDepth"`
      PeakDepth        int       `json:"peakDepth"`
      IsRealtimeStream bool      `json:"isRealtimeStream"`
      HasPeakIndicator bool      `json:"hasPeakIndicator"`
      RecordedAt       time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/tui/sparkline.go`.

---

### Improvement 1327: High-Resolution Packet Timing Scatterplot Terminal Matrix
- **ID**: `IMP-1327`
- **Cluster**: 133
- **Goal**: Map individual packet arrival times against ACK delays on an ANSI terminal scatterplot grid, revealing clustered packet loss and micro-bursting anomalies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: 2D scatterplot matrix using braille and dot characters; green dots (`#00E676`) for normal packets, orange dots (`#FF9100`) for delayed packets.
  - *Web*: `<packet-scatterplot-matrix>` Web Component with interactive point inspection and zoom.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui scatterplot --session-id=<id> [--is-burst-highlighted] [--has-outlier-rejection]`
  - Flag `--session-id`: Target SMTP warming session identifier.
  - Flag `--is-burst-highlighted`: Renders packet bursts with distinct color hues.
  - Flag `--has-outlier-rejection`: Scales plot bounds by rejecting statistical 3-sigma outliers.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ScatterplotDataPoints (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      relativeTimeMs REAL NOT NULL DEFAULT 0.0,
      ackDelayUs INTEGER NOT NULL DEFAULT 0,
      isBurstHighlighted INTEGER NOT NULL DEFAULT 1,
      hasOutlierRejection INTEGER NOT NULL DEFAULT 1,
      loggedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ScatterplotDataPoint struct {
      ID                  string    `json:"id"`
      SessionID           string    `json:"sessionId"`
      RelativeTimeMs      float64   `json:"relativeTimeMs"`
      AckDelayUs          int       `json:"ackDelayUs"`
      IsBurstHighlighted  bool      `json:"isBurstHighlighted"`
      HasOutlierRejection bool      `json:"hasOutlierRejection"`
      LoggedAt            time.Time `json:"loggedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/scatterplot.go`.

---

### Improvement 1328: Zero-CDN Web Component Telemetry Bridge & Shadow DOM Exporter
- **ID**: `IMP-1328`
- **Cluster**: 133
- **Goal**: Export live telemetry streams over local WebSocket connections into standard zero-dependency Shadow DOM Web Components (`<warmup-telemetry-hud>`), eliminating external framework bundles.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss card showing active local WebSocket clients and stream bitrate in KB/s.
  - *Web*: Pure native HTML5/ES6 Custom Elements rendered inside Shadow DOM with encapsulated CSS tokens.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui serve-web [--port=<port>] [--is-shadow-dom-strict] [--has-live-websocket]`
  - Flag `--port`: Local HTTP/WebSocket listening port (default: `8970`).
  - Flag `--is-shadow-dom-strict`: Enforces `mode: 'closed'` Shadow DOM encapsulation.
  - Flag `--has-live-websocket`: Streams microsecond socket events directly to connected web clients.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WebBridgeSessions (
      id TEXT PRIMARY KEY,
      clientIp TEXT NOT NULL,
      connectedPort INTEGER NOT NULL DEFAULT 8970,
      isShadowDomStrict INTEGER NOT NULL DEFAULT 1,
      hasLiveWebsocket INTEGER NOT NULL DEFAULT 1,
      startedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WebBridgeSession struct {
      ID                string    `json:"id"`
      ClientIP          string    `json:"clientIp"`
      ConnectedPort     int       `json:"connectedPort"`
      IsShadowDomStrict bool      `json:"isShadowDomStrict"`
      HasLiveWebsocket  bool      `json:"hasLiveWebsocket"`
      StartedAt         time.Time `json:"startedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Server bridge located in `pkg/warmup/web/bridge.go`.

---

### Improvement 1329: Keyboard-Driven Viewport Pan/Zoom & Terminal Inspector Controller
- **ID**: `IMP-1329`
- **Cluster**: 133
- **Goal**: Enable keyboard-driven pan and zoom across latency flamegraphs and waveforms using `+`, `-`, `h`, `j`, `k`, `l` navigation in Bubble Tea v2 viewports.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Floating Lipgloss help footer listing available navigation hotkeys in gold (`#FFD700`).
  - *Web*: `<pan-zoom-controls>` Web Component with clickable on-screen zoom controls and keyboard event listeners.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui inspect [--is-vim-keys] [--has-zoom-indicator]`
  - Flag `--is-vim-keys`: Enables Vim navigation hotkeys (`h`, `j`, `k`, `l`) alongside arrow keys.
  - Flag `--has-zoom-indicator`: Displays real-time zoom magnification level (e.g. `2.5x`) in viewport header.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ViewportNavigationStates (
      id TEXT PRIMARY KEY,
      zoomLevel REAL NOT NULL DEFAULT 1.0,
      panOffsetCols INTEGER NOT NULL DEFAULT 0,
      isVimKeys INTEGER NOT NULL DEFAULT 1,
      hasZoomIndicator INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ViewportNavigationState struct {
      ID               string    `json:"id"`
      ZoomLevel        float64   `json:"zoomLevel"`
      PanOffsetCols    int       `json:"panOffsetCols"`
      IsVimKeys        bool      `json:"isVimKeys"`
      HasZoomIndicator bool      `json:"hasZoomIndicator"`
      UpdatedAt        time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Viewport navigation handler in `pkg/warmup/tui/navigation.go`.

---

### Improvement 1330: Terminal Session Snapshot Archiver & Golden Trace Generator
- **ID**: `IMP-1330`
- **Cluster**: 133
- **Goal**: Snapshot the entire terminal viewport into ANSI-escaped text files and JSON-encoded trace logs for automated test verification and retrospective debugging.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss modal dialog confirming snapshot generation with checksum hash in violet (`#7C4DFF`).
  - *Web*: `<snapshot-viewer>` Web Component rendering ANSI color codes natively into browser viewports.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui snapshot --output-dir=<dir> [--is-ansi-stripped] [--has-json-metadata]`
  - Flag `--output-dir`: Destination directory for snapshot trace file.
  - Flag `--is-ansi-stripped`: Also writes a clean plaintext version without ANSI escape sequences.
  - Flag `--has-json-metadata`: Generates accompanying `.meta.json` containing timestamp and window metrics.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TerminalSnapshots (
      id TEXT PRIMARY KEY,
      snapshotPath TEXT NOT NULL,
      sha256Hash TEXT NOT NULL,
      isAnsiStripped INTEGER NOT NULL DEFAULT 0,
      hasJsonMetadata INTEGER NOT NULL DEFAULT 1,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TerminalSnapshot struct {
      ID              string    `json:"id"`
      SnapshotPath    string    `json:"snapshotPath"`
      Sha256Hash      string    `json:"sha256Hash"`
      IsAnsiStripped  bool      `json:"isAnsiStripped"`
      HasJsonMetadata bool      `json:"hasJsonMetadata"`
      CapturedAt      time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/snapshot.go`.
