# Improvements 1501 to 1530: Session Pipelining, Socket Health & Jitter Waveform TUI

This document details the complete specifications for Improvements 1501 through 1530 across Clusters 151, 152, and 153 of Sequence 352.

---

# Cluster 151: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (1501–1510)

---

### Improvement 1501: RFC 2920 Command Streaming Pipeline & Zero-RTT Command Burst Engine
- **ID**: `IMP-1501`
- **Cluster**: 151 (Autonomous SMTP Session Pipelining & TCP Zero-Window Detection)
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

### Improvement 1502: Real-Time Zero-Copy TCP Receive Window Monitor via Socket Reflection
- **ID**: `IMP-1502`
- **Cluster**: 151
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

### Improvement 1503: Autonomous TCP Zero-Window Detection & Preemptive Backoff Controller
- **ID**: `IMP-1503`
- **Cluster**: 151
- **Goal**: Detect when remote ISP MTA advertises a zero-length receive window (`snd_wnd == 0`), automatically halting pipelined envelope transmission and engaging microsecond backoff timers without dropping the TCP connection.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss banner rendering a pulsing magenta alert (`#E040FB`) signaling zero-window state, with an ASCII countdown timer for next drain probe.
  - *Web*: `<zero-window-banner>` Web Component providing interactive backoff status with zero CDN dependencies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline zero-window-control --session-id=<id> [--probe-interval-ms=<ms>] [--is-backoff-active] [--has-exponential-decay]`
  - Flag `--session-id`: Active SMTP session ID.
  - Flag `--probe-interval-ms`: Keepalive window probe interval in milliseconds (default: `100`).
  - Flag `--is-backoff-active`: Forces immediate entry into backoff state.
  - Flag `--has-exponential-decay`: Enables exponential decay multiplier on probe intervals.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowEvents (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      backoffDurationMs INTEGER NOT NULL DEFAULT 0,
      probeCount INTEGER NOT NULL DEFAULT 0,
      isBackoffActive INTEGER NOT NULL DEFAULT 1,
      hasRecovered INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowEvent struct {
      ID                string    `json:"id"`
      SessionID         string    `json:"sessionId"`
      BackoffDurationMs int       `json:"backoffDurationMs"`
      ProbeCount        int       `json:"probeCount"`
      IsBackoffActive   bool      `json:"isBackoffActive"`
      HasRecovered      bool      `json:"hasRecovered"`
      CreatedAt         time.Time `json:"createdAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated with `pkg/warmup/tcp/backoff.go`.

---

### Improvement 1504: Dynamic TCP Zero-Window Probe Injector & Drain ACK Evaluator
- **ID**: `IMP-1504`
- **Cluster**: 151
- **Goal**: Emit discrete single-byte TCP keepalive probes into closed zero-window sessions, analyzing incoming TCP ACK segment window increments to resume transmission with minimal latency.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss table row showing probe sequence numbers in violet (`#7C4DFF`) and returned ACK advertised window sizes in emerald (`#00E676`).
  - *Web*: `<window-probe-tracer>` Web Component rendering a live SVG dot-plot of probe-to-ACK response cycles.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline inject-probe --socket-fd=<fd> [--max-probes=<count>] [--is-drain-verified] [--has-payload-byte]`
  - Flag `--socket-fd`: Target socket descriptor.
  - Flag `--max-probes`: Maximum allowable keepalive zero-window probes before timeout.
  - Flag `--is-drain-verified`: Verifies remote buffer drain before burst resumption.
  - Flag `--has-payload-byte`: Sends standard RFC 793 1-byte payload probe.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowProbes (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      probeSeq INTEGER NOT NULL,
      ackWindowBytes INTEGER NOT NULL DEFAULT 0,
      isDrainVerified INTEGER NOT NULL DEFAULT 0,
      hasPayloadByte INTEGER NOT NULL DEFAULT 1,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowProbe struct {
      ID              string    `json:"id"`
      SocketFd        int       `json:"socketFd"`
      ProbeSeq        int       `json:"probeSeq"`
      AckWindowBytes  int       `json:"ackWindowBytes"`
      IsDrainVerified bool      `json:"isDrainVerified"`
      HasPayloadByte  bool      `json:"hasPayloadByte"`
      ProbedAt        time.Time `json:"probedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Placed in `pkg/warmup/tcp/probe_injector.go`.

---

### Improvement 1505: RFC 2920 Multi-Line Response Parser with SIMD Token Scanning
- **ID**: `IMP-1505`
- **Cluster**: 151
- **Goal**: Accelerate SMTP pipelining response parsing using zero-allocation SIMD byte-scanning (AVX2/NEON) to demultiplex intermediate `250-` and final `250 ` status codes across concurrent burst streams.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss inline badge displaying parsed responses/sec throughput in electric cyan (`#00E5FF`) with zero garbage collection pause markers.
  - *Web*: `<simd-parser-monitor>` Web Component showing a live sparkline of parser throughput and memory allocations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline parse-response --input-buffer=<hex> [--is-simd-enabled] [--has-multiline-continuation]`
  - Flag `--input-buffer`: Hex-encoded raw socket response payload.
  - Flag `--is-simd-enabled`: Toggles AVX2/NEON SIMD hardware acceleration.
  - Flag `--has-multiline-continuation`: Indicates presence of hyphenated continuation response.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpParsedResponses (
      id TEXT PRIMARY KEY,
      statusCode INTEGER NOT NULL,
      isSimdEnabled INTEGER NOT NULL DEFAULT 1,
      hasMultilineContinuation INTEGER NOT NULL DEFAULT 0,
      parseDurationNs INTEGER NOT NULL DEFAULT 0,
      parsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpParsedResponse struct {
      ID                       string    `json:"id"`
      StatusCode               int       `json:"statusCode"`
      IsSimdEnabled            bool      `json:"isSimdEnabled"`
      HasMultilineContinuation bool      `json:"hasMultilineContinuation"`
      ParseDurationNs          int64     `json:"parseDurationNs"`
      ParsedAt                 time.Time `json:"parsedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/smtp/simd_parser.go`.

---

### Improvement 1506: Adaptive TCP Send Buffer Right-Sizing & Zero-Window Guard
- **ID**: `IMP-1506`
- **Cluster**: 151
- **Goal**: Dynamically tune kernel `SO_SNDBUF` socket sizes via feedback loops to align perfectly with remote advertised windows, completely eliminating buffer-bloat-induced zero-window stalls.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss dual-meter rendering current OS `sndbuf` against remote `snd_wnd`, highlighted with amber bars (`#FFD600`) when divergence exceeds 20%.
  - *Web*: `<sndbuf-balancer-card>` Web Component displaying buffer utilization histograms.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline tune-sndbuf --socket-fd=<fd> [--target-bytes=<bytes>] [--is-auto-tuned] [--has-kernel-sync]`
  - Flag `--socket-fd`: Target socket file descriptor.
  - Flag `--target-bytes`: Target send buffer size in bytes.
  - Flag `--is-auto-tuned`: Activates continuous proportional-integral (PI) auto-tuning.
  - Flag `--has-kernel-sync`: Immediately commits buffer changes via `setsockopt`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SndBufTuningLogs (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      allocatedBytes INTEGER NOT NULL DEFAULT 0,
      optimalBytes INTEGER NOT NULL DEFAULT 0,
      isAutoTuned INTEGER NOT NULL DEFAULT 1,
      hasKernelSync INTEGER NOT NULL DEFAULT 1,
      tunedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SndBufTuningLog struct {
      ID             string    `json:"id"`
      SocketFd       int       `json:"socketFd"`
      AllocatedBytes int       `json:"allocatedBytes"`
      OptimalBytes   int       `json:"optimalBytes"`
      IsAutoTuned    bool      `json:"isAutoTuned"`
      HasKernelSync  bool      `json:"hasKernelSync"`
      TunedAt        time.Time `json:"tunedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Positioned in `pkg/warmup/tcp/sndbuf_tuner.go`.

---

### Improvement 1507: Pipelined DATA Command Early Termination on RCPT Failure
- **ID**: `IMP-1507`
- **Cluster**: 151
- **Goal**: Preemptively abort transmission of large MIME payloads when pipelined `RCPT TO` returns a 5xx fatal rejection, saving outbound bandwidth and avoiding wasted TCP window capacity.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss notification tag showing intercepted 5xx codes in crimson (`#FF1744`) and bytes saved in emerald (`#00E676`).
  - *Web*: Custom Web Component `<pipeline-abort-modal>` showing early termination details.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline abort-early --session-id=<id> [--is-early-abort-active] [--has-partial-data]`
  - Flag `--session-id`: ID of target warming session.
  - Flag `--is-early-abort-active`: Enables instant socket flush abortion on fatal recipient code.
  - Flag `--has-partial-data`: Reports whether MIME body payload had partially hit the wire.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EarlyAbortEvents (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      fatalStatusCode INTEGER NOT NULL,
      bytesSaved INTEGER NOT NULL DEFAULT 0,
      isEarlyAbortActive INTEGER NOT NULL DEFAULT 1,
      hasPartialData INTEGER NOT NULL DEFAULT 0,
      abortedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EarlyAbortEvent struct {
      ID                 string    `json:"id"`
      SessionID          string    `json:"sessionId"`
      FatalStatusCode    int       `json:"fatalStatusCode"`
      BytesSaved         int64     `json:"bytesSaved"`
      IsEarlyAbortActive bool      `json:"isEarlyAbortActive"`
      HasPartialData     bool      `json:"hasPartialData"`
      AbortedAt          time.Time `json:"abortedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Placed in `pkg/warmup/smtp/early_abort.go`.

---

### Improvement 1508: TCP Zero-Window Out-of-Band Queue Failover Daemon
- **ID**: `IMP-1508`
- **Cluster**: 151
- **Goal**: Automatically divert queued warming envelopes to alternative IP interface routes when the primary route encounters persistent zero-window state exceeding 1500ms.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss flow diagram indicating route rerouting from stalled IP (`192.0.2.1`) in amber (`#FFD600`) to failover IP (`198.51.100.2`) in cyan (`#00E5FF`).
  - *Web*: `<failover-route-graph>` Web Component rendering interactive multi-path egress topology.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline failover-route --route-id=<id> [--timeout-ms=<ms>] [--is-failover-engaged] [--has-alt-ip]`
  - Flag `--route-id`: Current route identifier.
  - Flag `--timeout-ms`: Threshold duration in zero-window before failover triggers.
  - Flag `--is-failover-engaged`: Manual override to activate failover.
  - Flag `--has-alt-ip`: Validates presence of secondary IP interface.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RouteFailoverEvents (
      id TEXT PRIMARY KEY,
      originalRouteId TEXT NOT NULL,
      targetRouteId TEXT NOT NULL,
      zeroWindowDurationMs INTEGER NOT NULL DEFAULT 0,
      isFailoverEngaged INTEGER NOT NULL DEFAULT 1,
      hasAltIp INTEGER NOT NULL DEFAULT 1,
      divertedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RouteFailoverEvent struct {
      ID                   string    `json:"id"`
      OriginalRouteID      string    `json:"originalRouteId"`
      TargetRouteID        string    `json:"targetRouteId"`
      ZeroWindowDurationMs int       `json:"zeroWindowDurationMs"`
      IsFailoverEngaged    bool      `json:"isFailoverEngaged"`
      HasAltIp             bool      `json:"hasAltIp"`
      DivertedAt           time.Time `json:"divertedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tcp/route_failover.go`.

---

### Improvement 1509: Sliding-Window Outbound Command Coalescing Dispatcher
- **ID**: `IMP-1509`
- **Cluster**: 151
- **Goal**: Implement a dynamic credit-based sliding window algorithm for pipelining SMTP commands, throttling in-flight commands according to remote server processing capacity.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss credit ledger showing available in-flight command slots (`Credits: 8/16`) in bright cyan (`#00E5FF`).
  - *Web*: `<pipeline-credits-meter>` Web Component displaying sliding window capacity.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline credit-window --session-id=<id> [--max-credits=<n>] [--is-credit-throttled] [--has-auto-scale]`
  - Flag `--session-id`: Active SMTP session ID.
  - Flag `--max-credits`: Maximum concurrent in-flight unacknowledged commands.
  - Flag `--is-credit-throttled`: Reports whether credit queue is currently blocked.
  - Flag `--has-auto-scale`: Dynamically scales credit window based on RTT variance.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PipelineCreditWindows (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      availableCredits INTEGER NOT NULL DEFAULT 0,
      maxCredits INTEGER NOT NULL DEFAULT 16,
      isCreditThrottled INTEGER NOT NULL DEFAULT 0,
      hasAutoScale INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PipelineCreditWindow struct {
      ID                string    `json:"id"`
      SessionID         string    `json:"sessionId"`
      AvailableCredits  int       `json:"availableCredits"`
      MaxCredits        int       `json:"maxCredits"`
      IsCreditThrottled bool      `json:"isCreditThrottled"`
      HasAutoScale      bool      `json:"hasAutoScale"`
      UpdatedAt         time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Added to `pkg/warmup/smtp/sliding_credit.go`.

---

### Improvement 1510: Forensics-Grade TCP Window Collapse Audit Stream
- **ID**: `IMP-1510`
- **Cluster**: 151
- **Goal**: Record microsecond-accurate TCP window collapse forensic events into the append-only audit database with cryptographic hash chains to prove ISP throttle compliance during disputes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss audit log table displaying collapsed window events with SHA-256 hash digests in deep violet (`#7C4DFF`).
  - *Web*: `<audit-collapse-table>` Web Component providing verifiable cryptographic audit trail.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup audit verify-collapse --event-id=<id> [--is-verified] [--has-valid-hash]`
  - Flag `--event-id`: Audit event UUID.
  - Flag `--is-verified`: Returns exit code 0 if cryptographic hash chain is valid.
  - Flag `--has-valid-hash`: Verifies SHA-256 hash against parent block.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpWindowCollapseAudits (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      windowSize INTEGER NOT NULL DEFAULT 0,
      parentHash TEXT NOT NULL,
      eventHash TEXT NOT NULL,
      isVerified INTEGER NOT NULL DEFAULT 1,
      hasValidHash INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpWindowCollapseAudit struct {
      ID           string    `json:"id"`
      SessionID    string    `json:"sessionId"`
      WindowSize   int       `json:"windowSize"`
      ParentHash   string    `json:"parentHash"`
      EventHash    string    `json:"eventHash"`
      IsVerified   bool      `json:"isVerified"`
      HasValidHash bool      `json:"hasValidHash"`
      RecordedAt   time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Placed in `pkg/audit/tcp_forensics.go`.

---

# Cluster 152: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (1511–1520)

---

### Improvement 1511: Multi-ISP Connection Reuse Pool with Pre-Warmed TLS Keepalive
- **ID**: `IMP-1511`
- **Cluster**: 152 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal**: Maintain an intelligent pool of persistent, pre-authenticated, and TLS-encrypted connections across Google, Microsoft, and Yahoo gateways to reduce handshake overhead during warming bursts.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss pool summary card showing active idle sockets per ISP in deep slate (`#161F30`) with live keep-alive pulses in emerald (`#00E676`).
  - *Web*: `<isp-pool-matrix>` Web Component displaying socket connection grids grouped by MX domain.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool status [--isp-name=<name>] [--is-idle-only] [--has-active-tls]`
  - Flag `--isp-name`: Filter pool metrics by target ISP (e.g. `google`, `microsoft`).
  - Flag `--is-idle-only`: Lists only warm idle connections waiting for dispatch.
  - Flag `--has-active-tls`: Confirms TLS session resumption status.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspConnectionPools (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      maxConnections INTEGER NOT NULL DEFAULT 20,
      activeConnections INTEGER NOT NULL DEFAULT 0,
      idleConnections INTEGER NOT NULL DEFAULT 0,
      isIdleOnly INTEGER NOT NULL DEFAULT 0,
      hasActiveTls INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspConnectionPool struct {
      ID                string    `json:"id"`
      IspName           string    `json:"ispName"`
      MaxConnections    int       `json:"maxConnections"`
      ActiveConnections int       `json:"activeConnections"`
      IdleConnections   int       `json:"idleConnections"`
      IsIdleOnly        bool      `json:"isIdleOnly"`
      HasActiveTls      bool      `json:"hasActiveTls"`
      UpdatedAt         time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/manager.go`.

---

### Improvement 1512: Cross-ISP Real-Time Socket Health Matrix Evaluator
- **ID**: `IMP-1512`
- **Cluster**: 152
- **Goal**: Compute a dynamic 0.0–1.0 health score for every active socket connection factoring in RTT drift, TCP retransmissions, TLS session tickets, and SMTP greeting delays.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss colored status pills showing socket score: green (`#00E676`) for > 0.85, amber (`#FFD600`) for 0.50–0.84, red (`#FF1744`) for < 0.50.
  - *Web*: `<socket-health-radar>` Web Component rendering interactive multi-metric radar charts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool health-matrix --socket-id=<id> [--is-degraded] [--has-retransmission-spike]`
  - Flag `--socket-id`: UUID of the tracked socket.
  - Flag `--is-degraded`: Filters for sockets with health score < 0.60.
  - Flag `--has-retransmission-spike`: Highlights sockets with > 3% TCP packet retransmission.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketHealthMatrices (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      ispName TEXT NOT NULL,
      healthScore REAL NOT NULL DEFAULT 1.0,
      rttDriftUs INTEGER NOT NULL DEFAULT 0,
      isDegraded INTEGER NOT NULL DEFAULT 0,
      hasRetransmissionSpike INTEGER NOT NULL DEFAULT 0,
      evaluatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketHealthMatrix struct {
      ID                     string    `json:"id"`
      SocketID               string    `json:"socketId"`
      IspName                string    `json:"ispName"`
      HealthScore            float64   `json:"healthScore"`
      RttDriftUs             int       `json:"rttDriftUs"`
      IsDegraded             bool      `json:"isDegraded"`
      HasRetransmissionSpike bool      `json:"hasRetransmissionSpike"`
      EvaluatedAt            time.Time `json:"evaluatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/pool/health.go`.

---

### Improvement 1513: Sub-Millisecond Heartbeat Pinger & Keepalive Idle Balancer
- **ID**: `IMP-1513`
- **Cluster**: 152
- **Goal**: Dispatch periodic lightweight `NOOP` or `RSET` commands across idle pooled connections to prevent ISP stateful firewall connection drops while distributing heartbeat jitter.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss heartbeat progress indicator rendering rhythmic cyan pulses (`#00E5FF`) along with idle timer countdowns.
  - *Web*: `<heartbeat-rhythm-strip>` Web Component visualizing ping cadence.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool ping --pool-id=<id> [--interval-ms=<ms>] [--is-jitter-applied] [--has-heartbeat-passed]`
  - Flag `--pool-id`: Target connection pool ID.
  - Flag `--interval-ms`: Base heartbeat ping interval in milliseconds.
  - Flag `--is-jitter-applied`: Adds +/- 15% pseudorandom jitter to ping schedule.
  - Flag `--has-heartbeat-passed`: Returns boolean status of last ping round.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolHeartbeatLogs (
      id TEXT PRIMARY KEY,
      poolId TEXT NOT NULL,
      pingIntervalMs INTEGER NOT NULL DEFAULT 30000,
      latencyUs INTEGER NOT NULL DEFAULT 0,
      isJitterApplied INTEGER NOT NULL DEFAULT 1,
      hasHeartbeatPassed INTEGER NOT NULL DEFAULT 1,
      pingedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolHeartbeatLog struct {
      ID                 string    `json:"id"`
      PoolID             string    `json:"poolId"`
      PingIntervalMs     int       `json:"pingIntervalMs"`
      LatencyUs          int       `json:"latencyUs"`
      IsJitterApplied    bool      `json:"isJitterApplied"`
      HasHeartbeatPassed bool      `json:"hasHeartbeatPassed"`
      PingedAt           time.Time `json:"pingedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/heartbeat.go`.

---

### Improvement 1514: Autonomous Socket Retiring & Zombie Connection Pruning
- **ID**: `IMP-1514`
- **Cluster**: 152
- **Goal**: Identify and gracefully retire stale or zombie sockets that have exceeded maximum reuse limits (e.g. 100 messages) or maximum lifetime (e.g. 300s), sending clean `QUIT` commands.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss cleanup notification tag showing retired socket IDs in slate (`#64748B`) and bytes reclaimed in emerald (`#00E676`).
  - *Web*: `<socket-reaper-card>` Web Component displaying socket lifecycle metrics.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool prune [--max-age-sec=<sec>] [--is-force-prune] [--has-graceful-quit]`
  - Flag `--max-age-sec`: Maximum allowable socket age before pruning (default: `300`).
  - Flag `--is-force-prune`: Forces immediate closure without awaiting in-flight queue drain.
  - Flag `--has-graceful-quit`: Emits standard SMTP `QUIT` command before TCP termination.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketPruningLogs (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      totalLifetimeSec INTEGER NOT NULL DEFAULT 0,
      messagesProcessed INTEGER NOT NULL DEFAULT 0,
      isForcePrune INTEGER NOT NULL DEFAULT 0,
      hasGracefulQuit INTEGER NOT NULL DEFAULT 1,
      prunedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketPruningLog struct {
      ID                string    `json:"id"`
      SocketID          string    `json:"socketId"`
      TotalLifetimeSec  int       `json:"totalLifetimeSec"`
      MessagesProcessed int       `json:"messagesProcessed"`
      IsForcePrune      bool      `json:"isForcePrune"`
      HasGracefulQuit   bool      `json:"hasGracefulQuit"`
      PrunedAt          time.Time `json:"prunedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled by `pkg/warmup/pool/pruner.go`.

---

### Improvement 1515: ISP Socket Capacity Forecaster via Linear Kalman Filtering
- **ID**: `IMP-1515`
- **Cluster**: 152
- **Goal**: Predict future socket connection acceptance limits for specific MX destinations by applying 1D Kalman filters to rolling handshake RTT and rejection frequencies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss forecast gauge indicating predicted safe concurrent socket limit per provider in bright cyan (`#00E5FF`).
  - *Web*: `<kalman-capacity-chart>` Web Component showing predicted vs actual connection limits over time.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool forecast-capacity --isp-name=<name> [--is-converged] [--has-high-confidence]`
  - Flag `--isp-name`: Target mailbox provider.
  - Flag `--is-converged`: Indicates Kalman state filter convergence.
  - Flag `--has-high-confidence`: True when estimation variance is below threshold.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspCapacityForecasts (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      predictedMaxSockets INTEGER NOT NULL DEFAULT 10,
      estimationVariance REAL NOT NULL DEFAULT 0.0,
      isConverged INTEGER NOT NULL DEFAULT 1,
      hasHighConfidence INTEGER NOT NULL DEFAULT 1,
      forecastedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspCapacityForecast struct {
      ID                  string    `json:"id"`
      IspName             string    `json:"ispName"`
      PredictedMaxSockets int       `json:"predictedMaxSockets"`
      EstimationVariance  float64   `json:"estimationVariance"`
      IsConverged         bool      `json:"isConverged"`
      HasHighConfidence   bool      `json:"hasHighConfidence"`
      ForecastedAt        time.Time `json:"forecastedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/kalman_forecast.go`.

---

### Improvement 1516: Dynamic Multi-Interface Source IP Affinity Binder
- **ID**: `IMP-1516`
- **Cluster**: 152
- **Goal**: Bind pooled outbound sockets across multiple local network interfaces and secondary IPv4/IPv6 addresses, matching source IP subnet allocation to destination reputation tiers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss multi-column IP routing table with binding status in electric cyan (`#00E5FF`) and interface name badges.
  - *Web*: `<ip-affinity-matrix>` Web Component displaying source-to-destination mapping.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool bind-affinity --source-ip=<ip> --isp-name=<name> [--is-bound] [--has-ipv6-support]`
  - Flag `--source-ip`: Egress source IP address.
  - Flag `--isp-name`: Target ISP group.
  - Flag `--is-bound`: Verifies binding lock to specified interface.
  - Flag `--has-ipv6-support`: Validates IPv6 reachability for target ISP.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IpAffinityBindings (
      id TEXT PRIMARY KEY,
      sourceIp TEXT NOT NULL,
      interfaceName TEXT NOT NULL,
      ispName TEXT NOT NULL,
      isBound INTEGER NOT NULL DEFAULT 1,
      hasIpv6Support INTEGER NOT NULL DEFAULT 0,
      boundAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IpAffinityBinding struct {
      ID             string    `json:"id"`
      SourceIP       string    `json:"sourceIp"`
      InterfaceName  string    `json:"interfaceName"`
      IspName        string    `json:"ispName"`
      IsBound        bool      `json:"isBound"`
      HasIpv6Support bool      `json:"hasIpv6Support"`
      BoundAt        time.Time `json:"boundAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/pool/ip_affinity.go`.

---

### Improvement 1517: Cross-ISP Connection Throttling & Burst Limiter Guard
- **ID**: `IMP-1517`
- **Cluster**: 152
- **Goal**: Enforce per-second and per-minute new connection creation rate limits per destination ISP to prevent rapid connection establishment triggers that fire anti-DDoS filters.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss rate limiter gauge rendering token bucket capacity in spring green (`#00E676`) and wait time in amber (`#FFD600`).
  - *Web*: `<connection-rate-meter>` Web Component displaying token bucket utilization.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool rate-limit --isp-name=<name> [--max-cps=<cps>] [--is-rate-limited] [--has-token-available]`
  - Flag `--isp-name`: Mailbox provider name.
  - Flag `--max-cps`: Maximum connections per second allowed.
  - Flag `--is-rate-limited`: Indicates active throttle engagement.
  - Flag `--has-token-available`: Checks if token bucket contains capacity.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ConnectionRateLimits (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      connectionsPerSecond INTEGER NOT NULL DEFAULT 5,
      tokensRemaining INTEGER NOT NULL DEFAULT 5,
      isRateLimited INTEGER NOT NULL DEFAULT 0,
      hasTokenAvailable INTEGER NOT NULL DEFAULT 1,
      checkedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ConnectionRateLimit struct {
      ID                   string    `json:"id"`
      IspName              string    `json:"ispName"`
      ConnectionsPerSecond int       `json:"connectionsPerSecond"`
      TokensRemaining      int       `json:"tokensRemaining"`
      IsRateLimited        bool      `json:"isRateLimited"`
      HasTokenAvailable    bool      `json:"hasTokenAvailable"`
      CheckedAt            time.Time `json:"checkedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Placed in `pkg/warmup/pool/rate_limiter.go`.

---

### Improvement 1518: Pooled Socket TLS Session Ticket Resumption Cache
- **ID**: `IMP-1518`
- **Cluster**: 152
- **Goal**: Cache and reuse RFC 5077 / RFC 8446 TLS session resumption tickets across pooled connections to avoid full asymmetric cryptographic handshakes on connection churn.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss cache statistic banner showing ticket reuse hit rate (`94.2%`) in electric cyan (`#00E5FF`) and CPU saving estimates.
  - *Web*: `<tls-ticket-cache-view>` Web Component rendering cached session ticket lifetimes.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool tls-cache [--isp-name=<name>] [--is-cache-hit] [--has-valid-ticket]`
  - Flag `--isp-name`: Target ISP name.
  - Flag `--is-cache-hit`: Checks if the last handshake reused a cached TLS ticket.
  - Flag `--has-valid-ticket`: Validates if unexpired tickets exist for the provider.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TlsSessionTickets (
      id TEXT PRIMARY KEY,
      serverName TEXT NOT NULL,
      ticketLifetimeSec INTEGER NOT NULL DEFAULT 7200,
      isCacheHit INTEGER NOT NULL DEFAULT 1,
      hasValidTicket INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TlsSessionTicket struct {
      ID                string    `json:"id"`
      ServerName        string    `json:"serverName"`
      TicketLifetimeSec int       `json:"ticketLifetimeSec"`
      IsCacheHit        bool      `json:"isCacheHit"`
      HasValidTicket    bool      `json:"hasValidTicket"`
      CreatedAt         time.Time `json:"createdAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/pool/tls_cache.go`.

---

### Improvement 1519: Socket Degradation Early Warning Circuit Breaker
- **ID**: `IMP-1519`
- **Cluster**: 152
- **Goal**: Instantly trip a circuit breaker on an entire ISP connection pool if 3 consecutive sockets experience sudden TCP RSTs or TLS handshake timeouts, pausing dispatch for 60 seconds.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss flashing alert card in vibrant red (`#FF1744`) displaying tripped breaker status and time remaining before half-open probe.
  - *Web*: `<pool-circuit-breaker-card>` Web Component displaying tripped status and cooldown ring.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool breaker-status --isp-name=<name> [--is-breaker-tripped] [--has-half-open-probe]`
  - Flag `--isp-name`: Target mailbox provider.
  - Flag `--is-breaker-tripped`: Boolean status of the pool circuit breaker.
  - Flag `--has-half-open-probe`: Indicates if pool is testing single probe connection.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolCircuitBreakers (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      failureCount INTEGER NOT NULL DEFAULT 0,
      cooldownRemainingSec INTEGER NOT NULL DEFAULT 0,
      isBreakerTripped INTEGER NOT NULL DEFAULT 0,
      hasHalfOpenProbe INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolCircuitBreaker struct {
      ID                   string    `json:"id"`
      IspName              string    `json:"ispName"`
      FailureCount         int       `json:"failureCount"`
      CooldownRemainingSec int       `json:"cooldownRemainingSec"`
      IsBreakerTripped     bool      `json:"isBreakerTripped"`
      HasHalfOpenProbe     bool      `json:"hasHalfOpenProbe"`
      UpdatedAt            time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Positioned in `pkg/warmup/pool/circuit_breaker.go`.

---

### Improvement 1520: Cryptographic Pool Audit Logger & Handshake Compliance Ledger
- **ID**: `IMP-1520`
- **Cluster**: 152
- **Goal**: Write immutable audit logs into `repodb/audit.db` recording socket pool reuse ratios, cipher suite negotiations, and connection termination causes with cryptographic Merkle verification.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss audit tree displaying cryptographic block verification in deep violet (`#7C4DFF`) with green verification checkmarks (`#00E676`).
  - *Web*: `<pool-audit-ledger>` Web Component presenting searchable cryptographic audit log.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup audit log-pool --socket-id=<id> [--is-audited] [--has-chain-valid]`
  - Flag `--socket-id`: Socked ID being logged.
  - Flag `--is-audited`: Confirms commit into audit database.
  - Flag `--has-chain-valid`: Verifies parent block Merkle hash.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolAuditLedgers (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      ispName TEXT NOT NULL,
      cipherSuite TEXT NOT NULL,
      parentHash TEXT NOT NULL,
      isAudited INTEGER NOT NULL DEFAULT 1,
      hasChainValid INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolAuditLedger struct {
      ID            string    `json:"id"`
      SocketID      string    `json:"socketId"`
      IspName       string    `json:"ispName"`
      CipherSuite   string    `json:"cipherSuite"`
      ParentHash    string    `json:"parentHash"`
      IsAudited     bool      `json:"isAudited"`
      HasChainValid bool      `json:"hasChainValid"`
      RecordedAt    time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Maintained in `pkg/audit/pool_ledger.go`.

---

# Cluster 153: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1521–1530)

---

### Improvement 1521: Bubble Tea v2 Real-Time Egress HUD Viewport & Stream Renderer
- **ID**: `IMP-1521`
- **Cluster**: 153 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Deliver a terminal HUD powered by Bubble Tea v2 that renders high-frequency socket events, in-flight envelopes, and queue pressure in real-time with zero flickering.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss multi-pane HUD with bordered viewports, displaying active pipeline streams in cyan (`#00E5FF`) and system memory in slate (`#64748B`).
  - *Web*: `<bubbletea-hud-mirror>` Web Component streaming terminal ANSI sequences via WebSockets.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui hud [--refresh-fps=<fps>] [--is-fullscreen] [--has-mouse-support]`
  - Flag `--refresh-fps`: Target HUD refresh rate in frames per second (default: `30`).
  - Flag `--is-fullscreen`: Launches Bubble Tea in alternate screen buffer.
  - Flag `--has-mouse-support`: Enables scroll wheel and click interaction in terminal.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiHudSessions (
      id TEXT PRIMARY KEY,
      viewportWidth INTEGER NOT NULL DEFAULT 120,
      viewportHeight INTEGER NOT NULL DEFAULT 40,
      fps INTEGER NOT NULL DEFAULT 30,
      isFullscreen INTEGER NOT NULL DEFAULT 1,
      hasMouseSupport INTEGER NOT NULL DEFAULT 1,
      startedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiHudSession struct {
      ID              string    `json:"id"`
      ViewportWidth   int       `json:"viewportWidth"`
      ViewportHeight  int       `json:"viewportHeight"`
      FPS             int       `json:"fps"`
      IsFullscreen    bool      `json:"isFullscreen"`
      HasMouseSupport bool      `json:"hasMouseSupport"`
      StartedAt       time.Time `json:"startedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/hud.go`.

---

### Improvement 1522: Sub-Millisecond Latency Phase Flamegraph Terminal Renderer
- **ID**: `IMP-1522`
- **Cluster**: 153
- **Goal**: Render a horizontal ASCII/ANSI flamegraph in Lipgloss detailing individual connection phases (DNS, TCP Syn, TLS Handshake, EHLO, MAIL, RCPT, DATA, QUIT) with sub-millisecond precision.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Color-segmented horizontal flamegraph bars: DNS (Violet `#7C4DFF`), TCP (Cyan `#00E5FF`), TLS (Amber `#FFD600`), SMTP (Emerald `#00E676`).
  - *Web*: `<latency-flamegraph-svg>` Web Component rendering interactive SVG flamegraph with hover tooltips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui flamegraph --session-id=<id> [--is-normalized] [--has-sub-ms-ticks]`
  - Flag `--session-id`: Target SMTP session identifier.
  - Flag `--is-normalized`: Scales flamegraph width to 100% of terminal column width.
  - Flag `--has-sub-ms-ticks`: Renders microsecond tick rulers across the top axis.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS LatencyPhaseMetrics (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      dnsDurationMs REAL NOT NULL DEFAULT 0.0,
      tcpDurationMs REAL NOT NULL DEFAULT 0.0,
      tlsDurationMs REAL NOT NULL DEFAULT 0.0,
      smtpDurationMs REAL NOT NULL DEFAULT 0.0,
      isNormalized INTEGER NOT NULL DEFAULT 1,
      hasSubMsTicks INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type LatencyPhaseMetric struct {
      ID             string    `json:"id"`
      SessionID      string    `json:"sessionId"`
      DnsDurationMs  float64   `json:"dnsDurationMs"`
      TcpDurationMs  float64   `json:"tcpDurationMs"`
      TlsDurationMs  float64   `json:"tlsDurationMs"`
      SmtpDurationMs float64   `json:"smtpDurationMs"`
      IsNormalized   bool      `json:"isNormalized"`
      HasSubMsTicks  bool      `json:"hasSubMsTicks"`
      RecordedAt     time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/flamegraph.go`.

---

### Improvement 1523: Bounded RTT Jitter Waveform ASCII Oscilloscope
- **ID**: `IMP-1523`
- **Cluster**: 153
- **Goal**: Plot a continuous rolling ASCII/Unicode oscilloscope (` ▂▃▄▅▆▇█`) displaying RTT jitter across consecutive packets, highlighting abnormal jitter spikes that indicate ISP traffic shaping.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: 20-row high Lipgloss oscilloscope canvas rendering jitter waves in electric cyan (`#00E5FF`) with trigger lines in amber (`#FFD600`).
  - *Web*: `<jitter-oscilloscope-card>` Web Component displaying vector SVG canvas oscilloscope.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui oscilloscope --socket-id=<id> [--window-size=<n>] [--is-triggered] [--has-peak-hold]`
  - Flag `--socket-id`: Monitored socket ID.
  - Flag `--window-size`: Number of samples displayed on screen (default: `60`).
  - Flag `--is-triggered`: Locks waveform display when threshold spike is detected.
  - Flag `--has-peak-hold`: Renders ghost line indicating maximum observed jitter peak.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterWaveformSamples (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      rttJitterUs INTEGER NOT NULL DEFAULT 0,
      sampleIndex INTEGER NOT NULL DEFAULT 0,
      isTriggered INTEGER NOT NULL DEFAULT 0,
      hasPeakHold INTEGER NOT NULL DEFAULT 1,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterWaveformSample struct {
      ID          string    `json:"id"`
      SocketID    string    `json:"socketId"`
      RttJitterUs int       `json:"rttJitterUs"`
      SampleIndex int       `json:"sampleIndex"`
      IsTriggered bool      `json:"isTriggered"`
      HasPeakHold bool      `json:"hasPeakHold"`
      SampledAt   time.Time `json:"sampledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Added to `pkg/warmup/tui/oscilloscope.go`.

---

### Improvement 1524: Zero-CDN Shadow DOM SVG Waveform Exporter
- **ID**: `IMP-1524`
- **Cluster**: 153
- **Goal**: Export standalone HTML/SVG telemetry reports containing interactive embedded Shadow DOM Web Components that visualize socket jitter and pipelining timelines without requiring any external JavaScript or CSS CDNs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss confirmation badge displaying generated report file path and total byte size.
  - *Web*: Self-contained `<warmup-telemetry-report>` Web Component with zero external script tags.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup export report --out-path=<path> [--is-self-contained] [--has-dark-theme]`
  - Flag `--out-path`: Target output file path (e.g. `reports/warmup-session.html`).
  - Flag `--is-self-contained`: Inlines all CSS and SVG definitions directly.
  - Flag `--has-dark-theme`: Applies dark `#0B0F19` void palette by default.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ExportedHtmlReports (
      id TEXT PRIMARY KEY,
      outPath TEXT NOT NULL,
      byteSize INTEGER NOT NULL DEFAULT 0,
      isSelfContained INTEGER NOT NULL DEFAULT 1,
      hasDarkTheme INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ExportedHtmlReport struct {
      ID              string    `json:"id"`
      OutPath         string    `json:"outPath"`
      ByteSize        int64     `json:"byteSize"`
      IsSelfContained bool      `json:"isSelfContained"`
      HasDarkTheme    bool      `json:"hasDarkTheme"`
      ExportedAt      time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/export/html_exporter.go`.

---

### Improvement 1525: Lipgloss TrueColor Multi-ISP Throughput Matrix
- **ID**: `IMP-1525`
- **Cluster**: 153
- **Goal**: Present a high-density TrueColor terminal table displaying messages/sec, delivery success %, 4xx deferral %, and average latency grouped by ISP.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss styled table with header borders in deep slate (`#161F30`), values dynamically colored from green (`#00E676`) to red (`#FF1744`).
  - *Web*: `<isp-throughput-matrix>` Web Component mirroring the terminal layout.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui isp-matrix [--is-compact] [--has-color-coding]`
  - Flag `--is-compact`: Omits secondary headers for small terminal widths (< 80 columns).
  - Flag `--has-color-coding`: Toggles TrueColor ANSI foreground formatting.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspThroughputSnapshots (
      id TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      messagesPerSec REAL NOT NULL DEFAULT 0.0,
      successRatePct REAL NOT NULL DEFAULT 100.0,
      deferralRatePct REAL NOT NULL DEFAULT 0.0,
      isCompact INTEGER NOT NULL DEFAULT 0,
      hasColorCoding INTEGER NOT NULL DEFAULT 1,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspThroughputSnapshot struct {
      ID              string    `json:"id"`
      IspName         string    `json:"ispName"`
      MessagesPerSec  float64   `json:"messagesPerSec"`
      SuccessRatePct  float64   `json:"successRatePct"`
      DeferralRatePct float64   `json:"deferralRatePct"`
      IsCompact       bool      `json:"isCompact"`
      HasColorCoding  bool      `json:"hasColorCoding"`
      CapturedAt      time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: In `pkg/warmup/tui/matrix.go`.

---

### Improvement 1526: Interactive Bubble Tea Command Palette & Session Filter
- **ID**: `IMP-1526`
- **Cluster**: 153
- **Goal**: Provide an interactive fuzzy-search command palette overlay (activated by `:`) inside the TUI to filter live session streams by recipient domain, IP address, or status code.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss modal popover with fuzzy match highlighting in electric cyan (`#00E5FF`) and result count indicators.
  - *Web*: `<terminal-command-palette>` Web Component offering interactive query input.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui palette [--query=<query>] [--is-modal-active] [--has-fuzzy-match]`
  - Flag `--query`: Pre-filled search query.
  - Flag `--is-modal-active`: Launches TUI directly with command modal focused.
  - Flag `--has-fuzzy-match`: Enables Levenshtein distance fuzzy filtering.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiPaletteFilters (
      id TEXT PRIMARY KEY,
      query TEXT NOT NULL,
      matchCount INTEGER NOT NULL DEFAULT 0,
      isModalActive INTEGER NOT NULL DEFAULT 1,
      hasFuzzyMatch INTEGER NOT NULL DEFAULT 1,
      appliedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiPaletteFilter struct {
      ID            string    `json:"id"`
      Query         string    `json:"query"`
      MatchCount    int       `json:"matchCount"`
      IsModalActive bool      `json:"isModalActive"`
      HasFuzzyMatch bool      `json:"hasFuzzyMatch"`
      AppliedAt     time.Time `json:"appliedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/palette.go`.

---

### Improvement 1527: Sub-Millisecond Socket Event Timeline Sparkline Strip
- **ID**: `IMP-1527`
- **Cluster**: 153
- **Goal**: Render a 1-line dense ASCII sparkline strip across the top of the terminal HUD showing microsecond-level RTT fluctuations over the past 120 seconds.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Unicode sparkline (` ▂▃▄▅▆▇█`) in bright cyan (`#00E5FF`) with current rolling mean and standard deviation displayed alongside.
  - *Web*: `<sparkline-timeline-strip>` Web Component with live SVG updates.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui sparkline [--history-sec=<sec>] [--is-sparkline-active] [--has-std-dev]`
  - Flag `--history-sec`: Rolling history duration in seconds (default: `120`).
  - Flag `--is-sparkline-active`: Toggles rendering of the top sparkline bar.
  - Flag `--has-std-dev`: Displays standard deviation next to sparkline.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SparklineSnapshots (
      id TEXT PRIMARY KEY,
      meanRttMs REAL NOT NULL DEFAULT 0.0,
      stdDevMs REAL NOT NULL DEFAULT 0.0,
      isSparklineActive INTEGER NOT NULL DEFAULT 1,
      hasStdDev INTEGER NOT NULL DEFAULT 1,
      snapshottedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SparklineSnapshot struct {
      ID                string    `json:"id"`
      MeanRttMs         float64   `json:"meanRttMs"`
      StdDevMs          float64   `json:"stdDevMs"`
      IsSparklineActive bool      `json:"isSparklineActive"`
      HasStdDev         bool      `json:"hasStdDev"`
      SnapshottedAt     time.Time `json:"snapshottedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tui/sparkline.go`.

---

### Improvement 1528: Real-Time Audio-Visual Acoustic Modem Debug Chime
- **ID**: `IMP-1528`
- **Cluster**: 153
- **Goal**: Generate optional subtle terminal terminal-bell frequency pitch shifts (or PCM audio pulses) representing transmission bursts and zero-window stalls for acoustic observability.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss sound toggle badge showing sound status `[MUTED]` or `[CHIME]` in violet (`#7C4DFF`).
  - *Web*: `<audio-acoustic-debug>` Web Component with Web Audio API synthesis.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui audio-debug [--pitch-hz=<hz>] [--is-audio-enabled] [--has-stall-alarm]`
  - Flag `--pitch-hz`: Base frequency in Hz for success chimes (default: `440`).
  - Flag `--is-audio-enabled`: Enables audio feedback pulses.
  - Flag `--has-stall-alarm`: Emits distinct descending tone on TCP zero-window events.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS AudioDebugConfigs (
      id TEXT PRIMARY KEY,
      pitchHz INTEGER NOT NULL DEFAULT 440,
      volumePct INTEGER NOT NULL DEFAULT 50,
      isAudioEnabled INTEGER NOT NULL DEFAULT 0,
      hasStallAlarm INTEGER NOT NULL DEFAULT 1,
      configuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type AudioDebugConfig struct {
      ID             string    `json:"id"`
      PitchHz        int       `json:"pitchHz"`
      VolumePct      int       `json:"volumePct"`
      IsAudioEnabled bool      `json:"isAudioEnabled"`
      HasStallAlarm  bool      `json:"hasStallAlarm"`
      ConfiguredAt   time.Time `json:"configuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/tui/audio.go`.

---

### Improvement 1529: Zero-Dependency Pure CSS/SVG Heatmap Web Visualizer
- **ID**: `IMP-1529`
- **Cluster**: 153
- **Goal**: Render hourly recipient engagement and delivery throughput heatmaps using pure CSS Grid and inline SVG, avoiding external charting libraries (Chart.js/D3).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: ASCII 24x7 matrix rendering intensity characters (`░▒▓█`) with Lipgloss green styling (`#00E676`).
  - *Web*: `<pure-svg-heatmap>` Web Component with hover cell tooltips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup visualizer heatmap [--is-hourly-view] [--has-gradient-fill]`
  - Flag `--is-hourly-view`: Groups matrix into 24 one-hour columns across 7 days.
  - Flag `--has-gradient-fill`: Applies TrueColor gradient shading.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS VisualizerHeatmapCells (
      id TEXT PRIMARY KEY,
      dayOfWeek INTEGER NOT NULL,
      hourOfDay INTEGER NOT NULL,
      deliveryVolume INTEGER NOT NULL DEFAULT 0,
      isHourlyView INTEGER NOT NULL DEFAULT 1,
      hasGradientFill INTEGER NOT NULL DEFAULT 1,
      calculatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type VisualizerHeatmapCell struct {
      ID              string    `json:"id"`
      DayOfWeek       int       `json:"dayOfWeek"`
      HourOfDay       int       `json:"hourOfDay"`
      DeliveryVolume  int       `json:"deliveryVolume"`
      IsHourlyView    bool      `json:"isHourlyView"`
      HasGradientFill bool      `json:"hasGradientFill"`
      CalculatedAt    time.Time `json:"calculatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Built into `pkg/warmup/visualizer/heatmap.go`.

---

### Improvement 1530: Forensic TUI Snapshot Exporter & Audit Verification Packager
- **ID**: `IMP-1530`
- **Cluster**: 153
- **Goal**: Bundle the current Bubble Tea viewport, raw terminal ANSI buffers, and corresponding cryptographic audit ledger hashes into a compressed forensic archive for compliance audits.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss modal displaying export progress bar and final SHA-256 bundle hash in bright emerald (`#00E676`).
  - *Web*: `<audit-packager-modal>` Web Component allowing one-click download of the cryptographic bundle.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup audit export-tui-snapshot --bundle-id=<id> [--is-compressed] [--has-cryptographic-hash]`
  - Flag `--bundle-id`: Target forensic package identifier.
  - Flag `--is-compressed`: Compresses snapshot with gzip.
  - Flag `--has-cryptographic-hash`: Appends Merkle hash verification certificate.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiAuditSnapshots (
      id TEXT PRIMARY KEY,
      bundleId TEXT NOT NULL,
      byteSize INTEGER NOT NULL DEFAULT 0,
      sha256Digest TEXT NOT NULL,
      isCompressed INTEGER NOT NULL DEFAULT 1,
      hasCryptographicHash INTEGER NOT NULL DEFAULT 1,
      sealedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiAuditSnapshot struct {
      ID                   string    `json:"id"`
      BundleID             string    `json:"bundleId"`
      ByteSize             int64     `json:"byteSize"`
      Sha256Digest         string    `json:"sha256Digest"`
      IsCompressed         bool      `json:"isCompressed"`
      HasCryptographicHash bool      `json:"hasCryptographicHash"`
      SealedAt             time.Time `json:"sealedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/audit/tui_snapshot.go`.
