# Improvements 1401 to 1430: Session Pipelining, Socket Health & Jitter Waveform TUI

This document details the complete specifications for Improvements 1401 through 1430 across Clusters 141, 142, and 143 of Sequence 351.

---

# Cluster 141: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (1401–1410)

---

### Improvement 1401: RFC 2920 Command Streaming Pipeline & Zero-RTT Command Burst Engine
- **ID**: `IMP-1401`
- **Cluster**: 141 (Autonomous SMTP Session Pipelining & TCP Zero-Window Detection)
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

### Improvement 1402: Real-Time Zero-Copy TCP Receive Window Monitor via Socket Reflection
- **ID**: `IMP-1402`
- **Cluster**: 141
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

### Improvement 1403: Autonomous TCP Zero-Window Detection & Preemptive Backoff Controller
- **ID**: `IMP-1403`
- **Cluster**: 141
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
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tcp/backoff_controller.go`.

---

### Improvement 1404: Adaptive Pipeline Chunk Sizing via Sliding Path MTU Discovery
- **ID**: `IMP-1404`
- **Cluster**: 141
- **Goal**: Dynamically calculate optimal command burst chunk sizes based on real-time PMTU discovery and TCP MSS values to eliminate packet fragmentation across cloud edge routers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss bar chart in deep slate (`#161F30`) and electric cyan (`#00E5FF`) illustrating current Path MTU versus chunk occupancy.
  - *Web*: `<pmtu-chunk-visualizer>` Web Component rendering interactive SVG payload fragmentation heatmaps.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline pmtu-tune --remote-host=<host> [--is-clamp-mss] [--has-df-bit]`
  - Flag `--remote-host`: Destination MX hostname or IP.
  - Flag `--is-clamp-mss`: Automatically clamps maximum segment size to avoid DF bit drops.
  - Flag `--has-df-bit`: Sets Dont-Fragment IP header bit during discovery probes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PmtuChunkConfigs (
      id TEXT PRIMARY KEY,
      remoteHost TEXT NOT NULL,
      pathMtuBytes INTEGER NOT NULL DEFAULT 1500,
      optimalChunkSize INTEGER NOT NULL DEFAULT 1460,
      isClampedMss INTEGER NOT NULL DEFAULT 1,
      hasDfBitSet INTEGER NOT NULL DEFAULT 1,
      calculatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PmtuChunkConfig struct {
      ID               string    `json:"id"`
      RemoteHost       string    `json:"remoteHost"`
      PathMtuBytes     int       `json:"pathMtuBytes"`
      OptimalChunkSize int       `json:"optimalChunkSize"`
      IsClampedMss     bool      `json:"isClampedMss"`
      HasDfBitSet      bool      `json:"hasDfBitSet"`
      CalculatedAt     time.Time `json:"calculatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/net/pmtu.go`.

---

### Improvement 1405: Non-Blocking Multi-Line SMTP Response Deserializer
- **ID**: `IMP-1405`
- **Cluster**: 141
- **Goal**: Ingest streaming TCP responses without allocating intermediate strings, parsing RFC 5321 multiline status codes (`250-`, `250 `) using SIMD byte scanning for maximum parsing throughput.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss log console displaying incoming token streams color-coded by RFC status class: 2xx green (`#00E676`), 3xx blue (`#2979FF`), 4xx amber (`#FFD600`), 5xx red (`#FF1744`).
  - *Web*: `<smtp-response-stream>` Web Component with zero-CDN virtualized list scrolling.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline parse-response --stream-fd=<fd> [--is-simd-accelerated] [--has-verbose-tokens]`
  - Flag `--stream-fd`: Streaming file descriptor.
  - Flag `--is-simd-accelerated`: Enables AVX2/NEON SIMD vector byte matching.
  - Flag `--has-verbose-tokens`: Emits tokenized structural AST in JSON lines format.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpParsedResponses (
      id TEXT PRIMARY KEY,
      statusCode INTEGER NOT NULL,
      isMultiline INTEGER NOT NULL DEFAULT 0,
      hasSimdScanned INTEGER NOT NULL DEFAULT 1,
      statusMessage TEXT NOT NULL,
      parsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpParsedResponse struct {
      ID             string    `json:"id"`
      StatusCode     int       `json:"statusCode"`
      IsMultiline    bool      `json:"isMultiline"`
      HasSimdScanned bool      `json:"hasSimdScanned"`
      StatusMessage  string    `json:"statusMessage"`
      ParsedAt       time.Time `json:"parsedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/smtp/parser.go`.

---

### Improvement 1406: Microsecond ACK Drift Compensator for Outbound Frame Throttling
- **ID**: `IMP-1406`
- **Cluster**: 141
- **Goal**: Continuously measure TCP ACK arrival drift relative to outbound transmission timestamps, applying microsecond pacing offsets to prevent triggering carrier-grade traffic shaping algorithms.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss micro-histogram showing ACK drift in microseconds (`0µs` to `5000µs`) with purple gradient styling (`#7C4DFF`).
  - *Web*: `<ack-drift-meter>` Web Component rendering real-time SVG scatter plot.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline ack-drift --session-id=<id> [--target-drift-us=<us>] [--is-auto-compensate] [--has-smoothing]`
  - Flag `--session-id`: Active session identifier.
  - Flag `--target-drift-us`: Acceptable ACK arrival variance threshold in microseconds.
  - Flag `--is-auto-compensate`: Dynamically adjusts frame pacing delay.
  - Flag `--has-smoothing`: Uses exponential moving average (EMA) on drift metrics.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS AckDriftMetrics (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      driftMicroseconds INTEGER NOT NULL,
      appliedDelayUs INTEGER NOT NULL DEFAULT 0,
      isCompensated INTEGER NOT NULL DEFAULT 1,
      hasSmoothedMetric INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type AckDriftMetric struct {
      ID                string    `json:"id"`
      SessionID         string    `json:"sessionId"`
      DriftMicroseconds int       `json:"driftMicroseconds"`
      AppliedDelayUs    int       `json:"appliedDelayUs"`
      IsCompensated     bool      `json:"isCompensated"`
      HasSmoothedMetric bool      `json:"hasSmoothedMetric"`
      RecordedAt        time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tcp/pacer.go`.

---

### Improvement 1407: TCP Window Probe Synthesizer with Zero-Payload Keep-Alives
- **ID**: `IMP-1407`
- **Cluster**: 141
- **Goal**: Generate precise, RFC-compliant zero-window probes containing a single already-acknowledged sequence byte, ensuring remote TCP stacks refresh their receive window status without discarding the connection.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss badge showing probe sequence number in bright violet (`#7C4DFF`) and remote TCP ACK window reply in spring emerald (`#00E676`).
  - *Web*: `<window-probe-badge>` Web Component with animated pulse waves.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline send-probe --socket-fd=<fd> [--is-raw-socket] [--has-ack-validation]`
  - Flag `--socket-fd`: Target socket descriptor.
  - Flag `--is-raw-socket`: Uses raw packet injection when standard TCP socket probe fails.
  - Flag `--has-ack-validation`: Verifies sequence and ACK numbers match expected kernel state.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WindowProbes (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      probeSeqNumber INTEGER NOT NULL,
      windowReplyBytes INTEGER NOT NULL,
      isWindowRestored INTEGER NOT NULL DEFAULT 0,
      hasAckValidated INTEGER NOT NULL DEFAULT 1,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WindowProbe struct {
      ID               string    `json:"id"`
      SocketFd         int       `json:"socketFd"`
      ProbeSeqNumber   int64     `json:"probeSeqNumber"`
      WindowReplyBytes int       `json:"windowReplyBytes"`
      IsWindowRestored bool      `json:"isWindowRestored"`
      HasAckValidated  bool      `json:"hasAckValidated"`
      ProbedAt         time.Time `json:"probedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/tcp/prober.go`.

---

### Improvement 1408: Pipeline Command Coalescing with Dynamic Nagle Override
- **ID**: `IMP-1408`
- **Cluster**: 141
- **Goal**: Dynamically toggle `TCP_NODELAY` and `TCP_CORK` flags per command stage to coalesce envelope transactions into single MTU frames while avoiding artificial Nagle algorithm transmission delays.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss toggle indicator displaying `CORK: ON` (`#FFD600`) during batch construction and `NODELAY: ACTIVE` (`#00E5FF`) on frame flush.
  - *Web*: `<socket-cork-indicator>` Web Component with interactive state transitions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline nagle-tune --session-id=<id> [--is-cork-enabled] [--has-immediate-flush]`
  - Flag `--session-id`: Target SMTP session identifier.
  - Flag `--is-cork-enabled`: Enables TCP_CORK during command buffering.
  - Flag `--has-immediate-flush`: Immediately forces TCP_NODELAY after DATA phase.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS NagleOverrides (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      isCorkActive INTEGER NOT NULL DEFAULT 0,
      hasFlushedImmediately INTEGER NOT NULL DEFAULT 1,
      coalescedBytes INTEGER NOT NULL DEFAULT 0,
      executedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type NagleOverride struct {
      ID                    string    `json:"id"`
      SessionID             string    `json:"sessionId"`
      IsCorkActive          bool      `json:"isCorkActive"`
      HasFlushedImmediately bool      `json:"hasFlushedImmediately"`
      CoalescedBytes        int       `json:"coalescedBytes"`
      ExecutedAt            time.Time `json:"executedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/net/socket_options.go`.

---

### Improvement 1409: Outbound Kernel Buffer Pressure Inversion Controller
- **ID**: `IMP-1409`
- **Cluster**: 141
- **Goal**: Monitor Linux kernel write queue capacity via `SIOCOUTQ` / `TIOCOUTQ`, automatically back-pressuring upstream warming worker queues before kernel socket buffers drop outbound email frames.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss horizontal gauge representing kernel queue fill ratio (0% to 100%) in amber flare (`#FFD600`) and vibrant crimson (`#FF1744`).
  - *Web*: `<kernel-queue-gauge>` Web Component with smooth CSS gauge animations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline buffer-pressure --socket-fd=<fd> [--threshold-bytes=<bytes>] [--is-inverted-flow] [--has-worker-signal]`
  - Flag `--socket-fd`: Inspected socket descriptor.
  - Flag `--threshold-bytes`: Buffer occupancy threshold that triggers backpressure.
  - Flag `--is-inverted-flow`: Inverts upstream queue intake to trickle rate.
  - Flag `--has-worker-signal`: Emits POSIX realtime signals to warming spoolers.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS BufferPressureStates (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      queuedBytes INTEGER NOT NULL,
      maxBufferBytes INTEGER NOT NULL,
      isInvertedFlow INTEGER NOT NULL DEFAULT 0,
      hasWorkerSignaled INTEGER NOT NULL DEFAULT 1,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type BufferPressureState struct {
      ID                string    `json:"id"`
      SocketFd          int       `json:"socketFd"`
      QueuedBytes       int       `json:"queuedBytes"`
      MaxBufferBytes    int       `json:"maxBufferBytes"`
      IsInvertedFlow    bool      `json:"isInvertedFlow"`
      HasWorkerSignaled bool      `json:"hasWorkerSignaled"`
      SampledAt         time.Time `json:"sampledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tcp/pressure.go`.

---

### Improvement 1410: Autonomous Zero-Window Predictive Exhaustion Forecaster
- **ID**: `IMP-1410`
- **Cluster**: 141
- **Goal**: Predict imminent remote TCP receive window collapses 50–150 milliseconds prior to zero-window events by fitting second-order polynomial regressions over window delta derivatives.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss forecast panel showing predicted time-to-exhaustion (`TTE`) in milliseconds, highlighted in vibrant crimson (`#FF1744`) when TTE < 100ms.
  - *Web*: `<zero-window-forecast>` Web Component rendering interactive forecast trajectories.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline forecast-exhaustion --session-id=<id> [--lookahead-ms=<ms>] [--is-preemptive-pause] [--has-drift-model]`
  - Flag `--session-id`: Target SMTP session.
  - Flag `--lookahead-ms`: Prediction horizon window (default: `200`).
  - Flag `--is-preemptive-pause`: Pauses command streaming before window hits zero.
  - Flag `--has-drift-model`: Incorporates historical ISP gateway drift coefficients.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowForecasts (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      predictedTteMs REAL NOT NULL,
      regressionConfidence REAL NOT NULL,
      isPreemptivelyPaused INTEGER NOT NULL DEFAULT 0,
      hasDriftModelApplied INTEGER NOT NULL DEFAULT 1,
      forecastedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowForecast struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      PredictedTteMs       float64   `json:"predictedTteMs"`
      RegressionConfidence float64   `json:"regressionConfidence"`
      IsPreemptivelyPaused bool      `json:"isPreemptivelyPaused"`
      HasDriftModelApplied bool      `json:"hasDriftModelApplied"`
      ForecastedAt         time.Time `json:"forecastedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Executed in `pkg/warmup/ml/exhaustion_forecaster.go`.

---

# Cluster 142: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (1411–1420)

---

### Improvement 1411: Multi-ISP Persistent Socket Pool Orchestrator & Keep-Alive Tuner
- **ID**: `IMP-1411`
- **Cluster**: 142 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal**: Maintain warmed, pre-authenticated, and TLS-session-cached persistent SMTP connection pools partitioned across major mailbox providers (Google, Microsoft, Yahoo, iCloud), minimizing repeated TCP and TLS handshakes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss multi-column table displaying active, idle, and reserved sockets per ISP provider with status pills in electric cyan (`#00E5FF`).
  - *Web*: `<isp-connection-pool-matrix>` Web Component with interactive connection count meters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool status [--isp=<isp-name>] [--is-healthy-only] [--has-idle-reclaim]`
  - Flag `--isp`: Filters pool status by specific ISP (e.g. `google`, `microsoft`).
  - Flag `--is-healthy-only`: Excludes degraded or draining sockets.
  - Flag `--has-idle-reclaim`: Triggers garbage collection of idle sockets exceeding max lifetime.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspConnectionPools (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      activeConnections INTEGER NOT NULL DEFAULT 0,
      idleConnections INTEGER NOT NULL DEFAULT 0,
      isHealthyOnly INTEGER NOT NULL DEFAULT 1,
      hasIdleReclaim INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspConnectionPool struct {
      ID                string    `json:"id"`
      IspName           string    `json:"ispName"`
      ActiveConnections int       `json:"activeConnections"`
      IdleConnections   int       `json:"idleConnections"`
      IsHealthyOnly     bool      `json:"isHealthyOnly"`
      HasIdleReclaim    bool      `json:"hasIdleReclaim"`
      UpdatedAt         time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/pool/manager.go`.

---

### Improvement 1412: Continuous Microsecond RTT Heatmap Matrix per Destination MX
- **ID**: `IMP-1412`
- **Cluster**: 142
- **Goal**: Track round-trip handshake and response latencies per MX destination in a continuous circular buffer, computing percentiles (p50, p95, p99) to detect ISP transit route degradation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss color-coded latency grid with p50 (<20ms: emerald `#00E676`), p95 (<60ms: amber `#FFD600`), and p99 (>100ms: crimson `#FF1744`).
  - *Web*: `<mx-rtt-heatmap>` Web Component with dynamic SVG heatmap rendering.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool rtt-matrix --mx-domain=<domain> [--is-percentile-view] [--has-anomaly-alert]`
  - Flag `--mx-domain`: Target MX gateway domain name.
  - Flag `--is-percentile-view`: Outputs p50, p90, p95, and p99 statistical percentiles.
  - Flag `--has-anomaly-alert`: Emits terminal warning when latency exceeds 3 standard deviations.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MxRttMatrices (
      id TEXT PRIMARY KEY,
      mxDomain TEXT NOT NULL,
      p50RttMs REAL NOT NULL DEFAULT 0.0,
      p95RttMs REAL NOT NULL DEFAULT 0.0,
      p99RttMs REAL NOT NULL DEFAULT 0.0,
      isDegraded INTEGER NOT NULL DEFAULT 0,
      hasAnomalyAlert INTEGER NOT NULL DEFAULT 1,
      evaluatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MxRttMatrix struct {
      ID              string    `json:"id"`
      MxDomain        string    `json:"mxDomain"`
      P50RttMs        float64   `json:"p50RttMs"`
      P95RttMs        float64   `json:"p95RttMs"`
      P99RttMs        float64   `json:"p99RttMs"`
      IsDegraded      bool      `json:"isDegraded"`
      HasAnomalyAlert bool      `json:"hasAnomalyAlert"`
      EvaluatedAt     time.Time `json:"evaluatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/telemetry/rtt_tracker.go`.

---

### Improvement 1413: Socket Health Degradation Matrix & Auto-Prune Engine
- **ID**: `IMP-1413`
- **Cluster**: 142
- **Goal**: Score individual pooled TCP sockets across multiple dimensions (retransmission count, ACK jitter, TLS alert history), automatically retiring degraded sockets before they drop live email deliveries.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status card showing socket degradation score (0–100) with color transition from spring green (`#00E676`) to crimson (`#FF1744`).
  - *Web*: `<socket-health-card>` Web Component with animated degradation progress bar.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool prune-degraded [--min-score=<score>] [--is-force-drain] [--has-audit-log]`
  - Flag `--min-score`: Minimum acceptable health score (default: `75`).
  - Flag `--is-force-drain`: Gracefully completes in-flight commands then closes socket.
  - Flag `--has-audit-log`: Appends socket retirement event to `repodb/audit.db`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketHealthMatrices (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      healthScore REAL NOT NULL DEFAULT 100.0,
      retransmits INTEGER NOT NULL DEFAULT 0,
      isDrained INTEGER NOT NULL DEFAULT 0,
      hasAuditLogged INTEGER NOT NULL DEFAULT 1,
      inspectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketHealthMatrix struct {
      ID             string    `json:"id"`
      SocketFd       int       `json:"socketFd"`
      HealthScore    float64   `json:"healthScore"`
      Retransmits    int       `json:"retransmits"`
      IsDrained      bool      `json:"isDrained"`
      HasAuditLogged bool      `json:"hasAuditLogged"`
      InspectedAt    time.Time `json:"inspectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/health_evaluator.go`.

---

### Improvement 1414: Cross-ISP TLS Handshake Latency Profiler & Session Ticket Resumption
- **ID**: `IMP-1414`
- **Cluster**: 142
- **Goal**: Measure TLS 1.3 ClientHello to Finished latency per ISP endpoint, optimizing session ticket resumption (`PSK`) to bypass full asymmetric key exchange on pooled connections.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss timeline chart displaying asymmetric exchange (`~45ms`) versus PSK resumption (`~4ms`) in bright cyan (`#00E5FF`).
  - *Web*: `<tls-resumption-profiler>` Web Component with interactive handshake comparison.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool tls-resume --isp=<name> [--is-psk-enforced] [--has-ticket-cache]`
  - Flag `--isp`: Target ISP name.
  - Flag `--is-psk-enforced`: Rejects non-PSK handshakes on warm pool connections.
  - Flag `--has-ticket-cache`: Persists encrypted session tickets across process restarts.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TlsHandshakeProfiles (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      fullHandshakeMs REAL NOT NULL,
      pskHandshakeMs REAL NOT NULL,
      isPskEnforced INTEGER NOT NULL DEFAULT 1,
      hasTicketCached INTEGER NOT NULL DEFAULT 1,
      profiledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TlsHandshakeProfile struct {
      ID              string    `json:"id"`
      IspName         string    `json:"ispName"`
      FullHandshakeMs float64   `json:"fullHandshakeMs"`
      PskHandshakeMs  float64   `json:"pskHandshakeMs"`
      IsPskEnforced   bool      `json:"isPskEnforced"`
      HasTicketCached bool      `json:"hasTicketCached"`
      ProfiledAt      time.Time `json:"profiledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tls/session_ticket.go`.

---

### Improvement 1415: ISP Concurrency Governor & Artificial Delay Inversion Sentinel
- **ID**: `IMP-1415`
- **Cluster**: 142
- **Goal**: Dynamically throttle simultaneous open connections to single destination MX clusters when remote gateways introduce artificial delays, preventing tarpits from exhausting local socket handles.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss throttle meter indicating active vs maximum permitted concurrency per MX cluster with amber flare styling (`#FFD600`).
  - *Web*: `<isp-concurrency-meter>` Web Component rendering real-time capacity dials.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool govern-concurrency --mx-cluster=<cluster> [--is-tarpit-detected] [--has-auto-scaling]`
  - Flag `--mx-cluster`: Destination MX cluster identifier.
  - Flag `--is-tarpit-detected`: Signals active tarpit behavior detected by response delay.
  - Flag `--has-auto-scaling`: Scales maximum concurrency dynamically based on completion rate.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspConcurrencyGovernors (
      id TEXT PRIMARY KEY,
      mxCluster TEXT NOT NULL,
      maxAllowedSockets INTEGER NOT NULL DEFAULT 10,
      activeSockets INTEGER NOT NULL DEFAULT 0,
      isTarpitDetected INTEGER NOT NULL DEFAULT 0,
      hasAutoScaled INTEGER NOT NULL DEFAULT 1,
      governedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspConcurrencyGovernor struct {
      ID                string    `json:"id"`
      MxCluster         string    `json:"mxCluster"`
      MaxAllowedSockets int       `json:"maxAllowedSockets"`
      ActiveSockets     int       `json:"activeSockets"`
      IsTarpitDetected  bool      `json:"isTarpitDetected"`
      HasAutoScaled     bool      `json:"hasAutoScaled"`
      GovernedAt        time.Time `json:"governedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/pool/concurrency_governor.go`.

---

### Improvement 1416: Zero-Copy Local Interface IP Shuffling & Binding Rotator
- **ID**: `IMP-1416`
- **Cluster**: 142
- **Goal**: Rotate outbound socket bindings across multiple local public IP interfaces using `SO_BINDTODEVICE` without memory allocations, balancing warm reputation across CIDR subnets.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss IP rotation roster displaying active binding interface and outbound message count per IP in electric cyan (`#00E5FF`).
  - *Web*: `<interface-rotator-roster>` Web Component with live packet count badges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool rotate-ip [--interface-group=<group>] [--is-round-robin] [--has-reputation-weighting]`
  - Flag `--interface-group`: Group of outbound network interfaces.
  - Flag `--is-round-robin`: Rotates interfaces in strict round-robin sequence.
  - Flag `--has-reputation-weighting`: Prioritizes interfaces with higher ISP trust scores.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS InterfaceIpRotations (
      id TEXT PRIMARY KEY,
      interfaceName TEXT NOT NULL,
      boundIpAddress TEXT NOT NULL,
      outboundCount INTEGER NOT NULL DEFAULT 0,
      isRoundRobinActive INTEGER NOT NULL DEFAULT 1,
      hasReputationWeight INTEGER NOT NULL DEFAULT 1,
      rotatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type InterfaceIpRotation struct {
      ID                  string    `json:"id"`
      InterfaceName       string    `json:"interfaceName"`
      BoundIpAddress      string    `json:"boundIpAddress"`
      OutboundCount       int       `json:"outboundCount"`
      IsRoundRobinActive  bool      `json:"isRoundRobinActive"`
      HasReputationWeight bool      `json:"hasReputationWeight"`
      RotatedAt           time.Time `json:"rotatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/net/interface_rotator.go`.

---

### Improvement 1417: TCP Keep-Alive Dynamic Interval Shaper for ISP Middleboxes
- **ID**: `IMP-1417`
- **Cluster**: 142
- **Goal**: Tune TCP keep-alive probe parameters (`TCP_KEEPIDLE`, `TCP_KEEPINTVL`, `TCP_KEEPCNT`) per destination ISP to keep idle sockets alive through NAT firewalls without triggering rate-limiting heuristics.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss parameter table showing socket keep-alive timers tuned per ISP in deep violet (`#7C4DFF`).
  - *Web*: `<keepalive-shaper-table>` Web Component with interactive countdown counters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool keepalive-tune --isp=<name> [--idle-secs=<secs>] [--is-kernel-tuned] [--has-probe-feedback]`
  - Flag `--isp`: Target ISP name.
  - Flag `--idle-secs`: Idle duration before first probe (default: `60`).
  - Flag `--is-kernel-tuned`: Directly tunes socket via setsockopt.
  - Flag `--has-probe-feedback`: Logs TCP ACK response timing from keep-alive probes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS KeepAliveProfiles (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      idleSecs INTEGER NOT NULL DEFAULT 60,
      intervalSecs INTEGER NOT NULL DEFAULT 15,
      isKernelTuned INTEGER NOT NULL DEFAULT 1,
      hasProbeFeedback INTEGER NOT NULL DEFAULT 1,
      configuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type KeepAliveProfile struct {
      ID               string    `json:"id"`
      IspName          string    `json:"ispName"`
      IdleSecs         int       `json:"idleSecs"`
      IntervalSecs     int       `json:"intervalSecs"`
      IsKernelTuned    bool      `json:"isKernelTuned"`
      HasProbeFeedback bool      `json:"hasProbeFeedback"`
      ConfiguredAt     time.Time `json:"configuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tcp/keepalive.go`.

---

### Improvement 1418: Socket Connection Leak Detector & Garbage Collection Sentinel
- **ID**: `IMP-1418`
- **Cluster**: 142
- **Goal**: Detect orphaned, hung, or leaked sockets across warm pool workers by scanning kernel file descriptor tables, automatically terminating unreferenced sockets after timeout thresholds.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss warning banner rendering orphaned socket count in vibrant red (`#FF1744`) with immediate reclaim action triggers.
  - *Web*: `<socket-leak-banner>` Web Component with zero-CDN interactive reclaim button.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool detect-leaks [--timeout-secs=<secs>] [--is-auto-reclaim] [--has-leak-trace]`
  - Flag `--timeout-secs`: Socket unreferenced idle duration before classification as leak (default: `120`).
  - Flag `--is-auto-reclaim`: Automatically closes leaked sockets and frees file descriptors.
  - Flag `--has-leak-trace`: Captures stack trace of thread that opened the leaked socket.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketLeakEvents (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      idleDurationSecs INTEGER NOT NULL,
      isReclaimed INTEGER NOT NULL DEFAULT 0,
      hasTraceCaptured INTEGER NOT NULL DEFAULT 1,
      detectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketLeakEvent struct {
      ID               string    `json:"id"`
      SocketFd         int       `json:"socketFd"`
      IdleDurationSecs int       `json:"idleDurationSecs"`
      IsReclaimed      bool      `json:"isReclaimed"`
      HasTraceCaptured bool      `json:"hasTraceCaptured"`
      DetectedAt       time.Time `json:"detectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/pool/leak_detector.go`.

---

### Improvement 1419: MX Failure Fast-Reroute & DNS Failover Matrix
- **ID**: `IMP-1419`
- **Cluster**: 142
- **Goal**: Instantly fail over egress email streams to secondary and tertiary MX records when primary gateways return connection refusal or sustained zero-window deferrals, maintaining warming continuity.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss routing switch diagram illustrating dynamic redirect from primary MX (`status: FAILED` in `#FF1744`) to secondary MX (`status: ACTIVE` in `#00E676`).
  - *Web*: `<mx-failover-flow>` Web Component rendering animated SVG switching paths.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool mx-failover --domain=<domain> [--is-failover-active] [--has-ttl-override]`
  - Flag `--domain`: Target recipient email domain.
  - Flag `--is-failover-active`: Forces traffic onto backup MX target.
  - Flag `--has-ttl-override`: Overrides local DNS cache TTL for fast propagation.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MxFailoverRoutes (
      id TEXT PRIMARY KEY,
      domain TEXT NOT NULL,
      primaryMx TEXT NOT NULL,
      activeBackupMx TEXT NOT NULL,
      isFailoverActive INTEGER NOT NULL DEFAULT 0,
      hasTtlOverride INTEGER NOT NULL DEFAULT 1,
      switchedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MxFailoverRoute struct {
      ID               string    `json:"id"`
      Domain           string    `json:"domain"`
      PrimaryMx        string    `json:"primaryMx"`
      ActiveBackupMx   string    `json:"activeBackupMx"`
      IsFailoverActive bool      `json:"isFailoverActive"`
      HasTtlOverride   bool      `json:"hasTtlOverride"`
      SwitchedAt       time.Time `json:"switchedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/dns/mx_failover.go`.

---

### Improvement 1420: Cross-Session Pre-Warmed TLS Context Multiplexer
- **ID**: `IMP-1420`
- **Cluster**: 142
- **Goal**: Multiplex pre-initialized cryptographic TLS client contexts across parallel worker routines, reducing memory footprint and eliminating repeated CA bundle loading during pool expansion.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss context status card showing shared TLS configuration reference count and memory savings in spring emerald (`#00E676`).
  - *Web*: `<tls-context-card>` Web Component displaying live memory footprint reduction stats.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool tls-mux status [--is-shared-context] [--has-crypto-bench]`
  - Flag `--is-shared-context`: Confirms single TLS config reference shared across workers.
  - Flag `--has-crypto-bench`: Executes micro-benchmarks comparing isolated vs multiplexed contexts.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TlsContextStats (
      id TEXT PRIMARY KEY,
      contextId TEXT NOT NULL,
      sharedWorkerCount INTEGER NOT NULL DEFAULT 1,
      savedMemoryBytes INTEGER NOT NULL DEFAULT 0,
      isSharedContext INTEGER NOT NULL DEFAULT 1,
      hasCryptoBenchRun INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TlsContextStat struct {
      ID                string    `json:"id"`
      ContextID         string    `json:"contextId"`
      SharedWorkerCount int       `json:"sharedWorkerCount"`
      SavedMemoryBytes  int64     `json:"savedMemoryBytes"`
      IsSharedContext   bool      `json:"isSharedContext"`
      HasCryptoBenchRun bool      `json:"hasCryptoBenchRun"`
      RecordedAt        time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/tls/context_mux.go`.

---

# Cluster 143: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1421–1430)

---

### Improvement 1421: Bubble Tea v2 Sub-Millisecond Jitter Waveform HUD Viewport
- **ID**: `IMP-1421`
- **Cluster**: 143 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Render a 60 FPS terminal HUD viewport using Bubble Tea v2 and Lipgloss depicting real-time microsecond socket packet arrival jitter across warming connections.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Real-time ASCII waveform chart using block runes (` ▂▃▄▅▆▇█`) with electric cyan (`#00E5FF`) baseline and crimson (`#FF1744`) peak spikes.
  - *Web*: `<jitter-waveform-hud>` Web Component using canvas-free inline SVG waveform vectors.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui jitter [--fps=<fps>] [--is-truecolor] [--has-sparkline]`
  - Flag `--fps`: Rendering frame rate (default: `60`).
  - Flag `--is-truecolor`: Forces 24-bit TrueColor ANSI output mode.
  - Flag `--has-sparkline`: Includes inline micro-sparkline per active socket thread.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterWaveformSamples (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      jitterUs INTEGER NOT NULL,
      peakJitterUs INTEGER NOT NULL,
      isTrueColorRendered INTEGER NOT NULL DEFAULT 1,
      hasSparklineIncluded INTEGER NOT NULL DEFAULT 1,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterWaveformSample struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      JitterUs             int       `json:"jitterUs"`
      PeakJitterUs         int       `json:"peakJitterUs"`
      IsTrueColorRendered  bool      `json:"isTrueColorRendered"`
      HasSparklineIncluded bool      `json:"hasSparklineIncluded"`
      SampledAt            time.Time `json:"sampledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/tui/jitter_view.go`.

---

### Improvement 1422: Interactive SMTP Phase Latency Flamegraph Renderer
- **ID**: `IMP-1422`
- **Cluster**: 143
- **Goal**: Generate high-definition flamegraphs decomposing SMTP transactions into granular phases (`DNS`, `TCP SYN`, `TLS Handshake`, `EHLO`, `MAIL`, `RCPT`, `DATA`, `QUIT`), enabling visual diagnosis of bottleneck phases.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss nested horizontal bar flamegraph color-coded by phase (DNS: violet `#7C4DFF`, TLS: cyan `#00E5FF`, DATA: emerald `#00E676`).
  - *Web*: `<smtp-flamegraph-view>` Web Component with zero-CDN interactive phase zoom and tooltips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui flamegraph --session-id=<id> [--is-interactive] [--has-phase-breakdown]`
  - Flag `--session-id`: Target SMTP transaction session.
  - Flag `--is-interactive`: Enables terminal keyboard cursor navigation across phase spans.
  - Flag `--has-phase-breakdown`: Prints detailed nanosecond duration table below flamegraph.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpPhaseFlamegraphs (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      totalDurationMs REAL NOT NULL,
      dnsDurationMs REAL NOT NULL,
      tlsDurationMs REAL NOT NULL,
      dataDurationMs REAL NOT NULL,
      isInteractiveRendered INTEGER NOT NULL DEFAULT 1,
      hasPhaseBreakdown INTEGER NOT NULL DEFAULT 1,
      renderedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpPhaseFlamegraph struct {
      ID                    string    `json:"id"`
      SessionID             string    `json:"sessionId"`
      TotalDurationMs       float64   `json:"totalDurationMs"`
      DnsDurationMs         float64   `json:"dnsDurationMs"`
      TlsDurationMs         float64   `json:"tlsDurationMs"`
      DataDurationMs        float64   `json:"dataDurationMs"`
      IsInteractiveRendered bool      `json:"isInteractiveRendered"`
      HasPhaseBreakdown     bool      `json:"hasPhaseBreakdown"`
      RenderedAt            time.Time `json:"renderedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tui/flamegraph_view.go`.

---

### Improvement 1423: Zero-CDN Shadow DOM SVG Waveform Exporter
- **ID**: `IMP-1423`
- **Cluster**: 143
- **Goal**: Export standalone HTML/SVG reports with embedded Shadow DOM styles and interactive SVG vector curves for email warming telemetry, requiring zero external script or font CDNs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss confirmation box displaying generated HTML file size and relative path in spring emerald (`#00E676`).
  - *Web*: Autonomous `<warmup-report>` Web Component with encapsulated CSS variables.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup export report --output=<relpath> [--is-self-contained] [--has-shadow-dom]`
  - Flag `--output`: Relative file path destination for the generated HTML/SVG report.
  - Flag `--is-self-contained`: Inlines all CSS, SVG definitions, and metadata into one file.
  - Flag `--has-shadow-dom`: Encapsulates markup within Shadow DOM boundary to avoid styling conflicts.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SvgExportLogs (
      id TEXT PRIMARY KEY,
      outputPath TEXT NOT NULL,
      fileSizeBytes INTEGER NOT NULL,
      isSelfContained INTEGER NOT NULL DEFAULT 1,
      hasShadowDom INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SvgExportLog struct {
      ID              string    `json:"id"`
      OutputPath      string    `json:"outputPath"`
      FileSizeBytes   int64     `json:"fileSizeBytes"`
      IsSelfContained bool      `json:"isSelfContained"`
      HasShadowDom    bool      `json:"hasShadowDom"`
      ExportedAt      time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/export/svg_exporter.go`.

---

### Improvement 1424: Terminal TrueColor Palette Tokenizer & Color Blindness Adaptive Switcher
- **ID**: `IMP-1424`
- **Cluster**: 143
- **Goal**: Provide standard semantic Lipgloss styling tokens across all TUI views with automated switching for protanopia, deuteranopia, and high-contrast terminal environments.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Live palette preview grid showing token mappings across standard TrueColor and accessible color-blind mode.
  - *Web*: `<palette-preview-grid>` Web Component rendering theme tokens.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui theme --mode=<default|protanopia|high-contrast> [--is-persistent] [--has-truecolor-fallback]`
  - Flag `--mode`: Active accessibility color mapping profile.
  - Flag `--is-persistent`: Saves selected theme preference to local configuration.
  - Flag `--has-truecolor-fallback`: Gracefully degrades to 16-color ANSI when 24-bit is unsupported.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiThemeConfigs (
      id TEXT PRIMARY KEY,
      themeMode TEXT NOT NULL DEFAULT 'default',
      isPersistent INTEGER NOT NULL DEFAULT 1,
      hasTruecolorFallback INTEGER NOT NULL DEFAULT 1,
      appliedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiThemeConfig struct {
      ID                   string    `json:"id"`
      ThemeMode            string    `json:"themeMode"`
      IsPersistent         bool      `json:"isPersistent"`
      HasTruecolorFallback bool      `json:"hasTruecolorFallback"`
      AppliedAt            time.Time `json:"appliedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/tui/theme.go`.

---

### Improvement 1425: Lipgloss Bounded Queue Depth Gauge & Buffer Pressure HUD
- **ID**: `IMP-1425`
- **Cluster**: 143
- **Goal**: Display outbound mail spool queue occupancy and kernel buffer pressure in a compact, real-time Lipgloss terminal gauge with sub-second refresh rates.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss progress bar with dynamic color gradients: 0–60% emerald (`#00E676`), 60–85% amber (`#FFD600`), 85–100% crimson (`#FF1744`).
  - *Web*: `<queue-depth-gauge>` Web Component rendering SVG circular dial.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui queue-gauge [--refresh-ms=<ms>] [--is-live-watch] [--has-buffer-details]`
  - Flag `--refresh-ms`: HUD refresh cadence in milliseconds (default: `100`).
  - Flag `--is-live-watch`: Runs persistent interactive watch loop.
  - Flag `--has-buffer-details`: Expands view to display per-socket TCP send buffer occupancy.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS QueueGaugeSnapshots (
      id TEXT PRIMARY KEY,
      queueOccupancyPct REAL NOT NULL,
      totalPendingMessages INTEGER NOT NULL,
      isLiveWatched INTEGER NOT NULL DEFAULT 1,
      hasBufferDetailsIncluded INTEGER NOT NULL DEFAULT 0,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type QueueGaugeSnapshot struct {
      ID                       string    `json:"id"`
      QueueOccupancyPct        float64   `json:"queueOccupancyPct"`
      TotalPendingMessages     int       `json:"totalPendingMessages"`
      IsLiveWatched            bool      `json:"isLiveWatched"`
      HasBufferDetailsIncluded bool      `json:"hasBufferDetailsIncluded"`
      CapturedAt               time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tui/queue_gauge.go`.

---

### Improvement 1426: Sub-Millisecond Event Spool Timeline & Live Scrollback Buffer
- **ID**: `IMP-1426`
- **Cluster**: 143
- **Goal**: Maintain an in-memory ring buffer of the last 10,000 warming events, permitting smooth keyboard-navigated scrollback, searching, and filtering in the Bubble Tea v2 TUI without disk I/O freezes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss scrollable viewport with sticky header and footer, keyboard page navigation (`PageUp`, `PageDown`, `j`, `k`), and fuzzy search overlay.
  - *Web*: `<event-scrollback-view>` Web Component with virtualized DOM rows.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui scrollback [--capacity=<num>] [--is-ring-buffered] [--has-fuzzy-search]`
  - Flag `--capacity`: Maximum events retained in memory (default: `10000`).
  - Flag `--is-ring-buffered`: Uses circular memory buffer with zero reallocations.
  - Flag `--has-fuzzy-search`: Enables interactive regex and fuzzy text matching.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ScrollbackBuffers (
      id TEXT PRIMARY KEY,
      bufferCapacity INTEGER NOT NULL DEFAULT 10000,
      retainedEventCount INTEGER NOT NULL DEFAULT 0,
      isRingBuffered INTEGER NOT NULL DEFAULT 1,
      hasFuzzySearchEnabled INTEGER NOT NULL DEFAULT 1,
      initializedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ScrollbackBuffer struct {
      ID                    string    `json:"id"`
      BufferCapacity        int       `json:"bufferCapacity"`
      RetainedEventCount    int       `json:"retainedEventCount"`
      IsRingBuffered        bool      `json:"isRingBuffered"`
      HasFuzzySearchEnabled bool      `json:"hasFuzzySearchEnabled"`
      InitializedAt         time.Time `json:"initializedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/scrollback.go`.

---

### Improvement 1427: Split-Screen Comparative ISP Latency Analyzer TUI
- **ID**: `IMP-1427`
- **Cluster**: 143
- **Goal**: Render a dual-pane side-by-side comparative TUI view contrasting latency distributions between two destination ISPs (e.g. Google vs Microsoft) to detect provider-specific network anomalies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss dual-pane layout split vertically with ISP 1 on the left in cyan (`#00E5FF`) and ISP 2 on the right in violet (`#7C4DFF`).
  - *Web*: `<isp-dual-compare-view>` Web Component with synchronized zoom.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui compare --isp-a=<name> --isp-b=<name> [--is-synced-scale] [--has-diff-metrics]`
  - Flag `--isp-a`: First ISP provider name.
  - Flag `--isp-b`: Second ISP provider name.
  - Flag `--is-synced-scale`: Locks vertical axis scaling between panes.
  - Flag `--has-diff-metrics`: Renders statistical delta and p-value significance.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspComparisonSessions (
      id TEXT PRIMARY KEY,
      ispA TEXT NOT NULL,
      ispB TEXT NOT NULL,
      avgLatencyDeltaMs REAL NOT NULL DEFAULT 0.0,
      isSyncedScale INTEGER NOT NULL DEFAULT 1,
      hasDiffMetricsIncluded INTEGER NOT NULL DEFAULT 1,
      comparedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspComparisonSession struct {
      ID                     string    `json:"id"`
      IspA                   string    `json:"ispA"`
      IspB                   string    `json:"ispB"`
      AvgLatencyDeltaMs      float64   `json:"avgLatencyDeltaMs"`
      IsSyncedScale          bool      `json:"isSyncedScale"`
      HasDiffMetricsIncluded bool      `json:"hasDiffMetricsIncluded"`
      ComparedAt             time.Time `json:"comparedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tui/compare_view.go`.

---

### Improvement 1428: Terminal ANSI Sparkline Density Scaler for Dynamic Terminal Resizing
- **ID**: `IMP-1428`
- **Cluster**: 143
- **Goal**: Automatically rescale sparkline metric bins and character resolutions upon receiving `tea.WindowSizeMsg`, preventing line wrapping and layout distortion on small or resizing terminal windows.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Smoothly reflowing sparkline graphs that adapt from 40 columns up to 240 columns without breaking borders or truncating text.
  - *Web*: `<responsive-sparkline>` Web Component using vector viewBox auto-scaling.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui sparkline-scale [--is-auto-resize] [--has-min-width-guard]`
  - Flag `--is-auto-resize`: Listens to window resize signals and re-renders layout.
  - Flag `--has-min-width-guard`: Displays clean fallback banner if terminal width < 80 columns.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SparklineScaleEvents (
      id TEXT PRIMARY KEY,
      terminalWidth INTEGER NOT NULL,
      terminalHeight INTEGER NOT NULL,
      calculatedBins INTEGER NOT NULL,
      isAutoResized INTEGER NOT NULL DEFAULT 1,
      hasMinWidthGuarded INTEGER NOT NULL DEFAULT 1,
      resizedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SparklineScaleEvent struct {
      ID                 string    `json:"id"`
      TerminalWidth      int       `json:"terminalWidth"`
      TerminalHeight     int       `json:"terminalHeight"`
      CalculatedBins     int       `json:"calculatedBins"`
      IsAutoResized      bool      `json:"isAutoResized"`
      HasMinWidthGuarded bool      `json:"hasMinWidthGuarded"`
      ResizedAt          time.Time `json:"resizedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/sparkline.go`.

---

### Improvement 1429: Outbound Envelope Spool State Matrix & Delivery Inspector HUD
- **ID**: `IMP-1429`
- **Cluster**: 143
- **Goal**: Provide an interactive drill-down inspector in Bubble Tea v2 for exploring queued email envelopes, inspecting headers, TLS cipher suites, and retry schedules in real time.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss modal dialog presenting envelope metadata in deep slate (`#161F30`) with syntax-highlighted SMTP headers and spring emerald (`#00E676`) status tags.
  - *Web*: `<envelope-inspector-modal>` Web Component with collapsible JSON payload viewer.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui inspect --envelope-id=<id> [--is-raw-view] [--has-header-mask]`
  - Flag `--envelope-id`: Identifier of envelope in warming spool.
  - Flag `--is-raw-view`: Renders unformatted raw RFC 5322 MIME envelope payload.
  - Flag `--has-header-mask`: Masks sensitive recipient data (PII) for compliance.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EnvelopeInspections (
      id TEXT PRIMARY KEY,
      envelopeId TEXT NOT NULL,
      recipientDomain TEXT NOT NULL,
      isRawViewed INTEGER NOT NULL DEFAULT 0,
      hasHeaderMasked INTEGER NOT NULL DEFAULT 1,
      inspectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EnvelopeInspection struct {
      ID              string    `json:"id"`
      EnvelopeID      string    `json:"envelopeId"`
      RecipientDomain string    `json:"recipientDomain"`
      IsRawViewed     bool      `json:"isRawViewed"`
      HasHeaderMasked bool      `json:"hasHeaderMasked"`
      InspectedAt     time.Time `json:"inspectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/envelope_inspector.go`.

---

### Improvement 1430: Zero-Dependency Interactive SVG Export Daemon & Live Browser Bridge
- **ID**: `IMP-1430`
- **Cluster**: 143
- **Goal**: Serve live warming telemetry flamegraphs and waveforms directly to local web browsers over an embedded HTTP/WebSocket server without requiring npm packages, bundlers, or remote CDN links.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss server startup banner displaying local listening URL (`http://127.0.0.1:9095`) and connected browser client count in bright cyan (`#00E5FF`).
  - *Web*: Full-screen real-time telemetry dashboard using native HTML5 Web Components and SVG animations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup serve [--port=<port>] [--is-open-browser] [--has-websocket-live]`
  - Flag `--port`: Local HTTP server port (default: `9095`).
  - Flag `--is-open-browser`: Automatically launches default system web browser.
  - Flag `--has-websocket-live`: Streams microsecond socket events via raw WebSocket connection.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS LiveServeConfigs (
      id TEXT PRIMARY KEY,
      portNumber INTEGER NOT NULL DEFAULT 9095,
      connectedClients INTEGER NOT NULL DEFAULT 0,
      isOpenBrowserActive INTEGER NOT NULL DEFAULT 1,
      hasWebSocketLiveActive INTEGER NOT NULL DEFAULT 1,
      startedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type LiveServeConfig struct {
      ID                     string    `json:"id"`
      PortNumber             int       `json:"portNumber"`
      ConnectedClients       int       `json:"connectedClients"`
      IsOpenBrowserActive    bool      `json:"isOpenBrowserActive"`
      HasWebSocketLiveActive bool      `json:"hasWebSocketLiveActive"`
      StartedAt              time.Time `json:"startedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/serve/live_bridge.go`.
