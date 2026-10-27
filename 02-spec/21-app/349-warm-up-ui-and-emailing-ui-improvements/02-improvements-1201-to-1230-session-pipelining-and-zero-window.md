# Improvements 1201 to 1230: Session Pipelining, Socket Health & Jitter Waveform TUI

This document details the complete specifications for Improvements 1201 through 1230 across Clusters 121, 122, and 123 of Sequence 349.

---

# Cluster 121: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (1201–1210)

---

### Improvement 1201: RFC 2920 Command Streaming Pipeline & Zero-RTT Command Burst Engine
- **ID**: `IMP-1201`
- **Cluster**: 121 (Autonomous SMTP Session Pipelining & TCP Zero-Window Detection)
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

### Improvement 1202: Real-Time Zero-Copy TCP Receive Window Monitor via Socket Reflection
- **ID**: `IMP-1202`
- **Cluster**: 121
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

### Improvement 1203: Autonomous TCP Zero-Window Detection & Preemptive Backoff Controller
- **ID**: `IMP-1203`
- **Cluster**: 121
- **Goal**: Detect when remote ISP MTA advertises a zero-length receive window (`snd_wnd == 0`), automatically halting pipelined envelope transmission and engaging microsecond backoff timers without dropping the TCP connection.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss banner rendering a pulsing magenta alert (`#E040FB`) signaling zero-window state, with an ASCII countdown timer for next drain probe.
  - *Web*: `<zero-window-banner>` Web Component providing interactive backoff status with zero CDN dependencies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline detect-zero-window --session-id=<id> [--probe-timeout-ms=<ms>] [--is-drain-pause-active] [--has-backoff-scaled]`
  - Flag `--session-id`: Active SMTP session ID.
  - Flag `--probe-timeout-ms`: Maximum duration to wait for window expansion before rerouting envelopes (default: `3000`).
  - Flag `--is-drain-pause-active`: Enforces pause on local outbound buffer flush.
  - Flag `--has-backoff-scaled`: Indicates whether exponential backoff multiplier is active.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowEvents (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      stallDurationMs REAL NOT NULL DEFAULT 0.0,
      probeCount INTEGER NOT NULL DEFAULT 0,
      isDrainPauseActive INTEGER NOT NULL DEFAULT 1,
      hasBackoffScaled INTEGER NOT NULL DEFAULT 1,
      detectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowEvent struct {
      ID                 string    `json:"id"`
      SessionID          string    `json:"sessionId"`
      StallDurationMs    float64   `json:"stallDurationMs"`
      ProbeCount         int       `json:"probeCount"`
      IsDrainPauseActive bool      `json:"isDrainPauseActive"`
      HasBackoffScaled   bool      `json:"hasBackoffScaled"`
      DetectedAt         time.Time `json:"detectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tcp/zero_window.go` inside the socket write pump.

---

### Improvement 1204: Outbound Socket Drain Telemetry & Epoll/Kqueue Edge-Triggered Poller
- **ID**: `IMP-1204`
- **Cluster**: 121
- **Goal**: Leverage edge-triggered `epoll` (Linux) / `kqueue` (BSD/macOS) / IOCP (Windows) to asynchronously wake up outbound socket write routines as soon as remote window space opens up.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss stream widget rendering microsecond poll wake events with neon cyan (`#00E5FF`) progress indicators.
  - *Web*: Shadow DOM canvas component `<socket-drain-graph>` rendering an event loop latency sparkline.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline poll-drain --socket-fd=<fd> [--wake-threshold-bytes=<bytes>] [--is-edge-triggered] [--has-iocp-registered]`
  - Flag `--socket-fd`: Socket descriptor to monitor.
  - Flag `--wake-threshold-bytes`: Minimum available buffer capacity required to wake the writer (default: `4096`).
  - Flag `--is-edge-triggered`: Configures EPOLLET / EV_CLEAR polling mode.
  - Flag `--has-iocp-registered`: Verifies IOCP port association on Windows.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketDrainPollers (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      wakeLatencyUs INTEGER NOT NULL DEFAULT 0,
      freedBytes INTEGER NOT NULL DEFAULT 0,
      isEdgeTriggered INTEGER NOT NULL DEFAULT 1,
      hasIocpRegistered INTEGER NOT NULL DEFAULT 1,
      polledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketDrainPoller struct {
      ID                string    `json:"id"`
      SocketFd          int       `json:"socketFd"`
      WakeLatencyUs     int       `json:"wakeLatencyUs"`
      FreedBytes        int       `json:"freedBytes"`
      IsEdgeTriggered   bool      `json:"isEdgeTriggered"`
      HasIocpRegistered bool      `json:"hasIocpRegistered"`
      PolledAt          time.Time `json:"polledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/transport/poller.go` governing low-level network loops.

---

### Improvement 1205: RFC 2920 Multi-Line Response Demultiplexer & Async Status Parser
- **ID**: `IMP-1205`
- **Cluster**: 121
- **Goal**: Parse and align batched multi-line SMTP server replies (e.g. `250-SIZE`, `250 OK`) matching pipelined commands in an asynchronous lock-free ring buffer without blocking execution.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss two-column table displaying outbound commands on the left and parsed matching response codes on the right with color-coded badges (`#00E676` for 2xx, `#FFD600` for 4xx, `#FF1744` for 5xx).
  - *Web*: `<smtp-response-demux>` Web Component highlighting response codes with interactive drill-down tooltips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline parse-responses --stream-id=<id> [--is-strict-sync] [--has-partial-buffer]`
  - Flag `--stream-id`: ID of the target command stream.
  - Flag `--is-strict-sync`: Enforces strict command-to-reply sequencing validation.
  - Flag `--has-partial-buffer`: Indicates if incomplete responses are currently buffered.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ResponseDemuxLogs (
      id TEXT PRIMARY KEY,
      streamId TEXT NOT NULL,
      matchedCommandCount INTEGER NOT NULL DEFAULT 0,
      unmatchedLineCount INTEGER NOT NULL DEFAULT 0,
      isStrictSync INTEGER NOT NULL DEFAULT 1,
      hasPartialBuffer INTEGER NOT NULL DEFAULT 0,
      parsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ResponseDemuxLog struct {
      ID                  string    `json:"id"`
      StreamID            string    `json:"streamId"`
      MatchedCommandCount int       `json:"matchedCommandCount"`
      UnmatchedLineCount  int       `json:"unmatchedLineCount"`
      IsStrictSync        bool      `json:"isStrictSync"`
      HasPartialBuffer    bool      `json:"hasPartialBuffer"`
      ParsedAt            time.Time `json:"parsedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/smtp/demux.go`.

---

### Improvement 1206: TCP Window Probing Keep-Alive Synthesizer & Silent Drop Detector
- **ID**: `IMP-1206`
- **Cluster**: 121
- **Goal**: Synthesize single-byte TCP keep-alive ACK probes during extended zero-window states to distinguish between busy MTA queue ingestion and dead/dropped carrier middlebox connections.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss probe ping widget displaying microsecond RTT responses and probe retry counts in electric violet (`#7C4DFF`).
  - *Web*: `<tcp-probe-card>` Web Component displaying probe timeline bars and drop warnings.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline send-probe --socket-fd=<fd> [--max-probes=<count>] [--is-ack-validated] [--has-drop-detected]`
  - Flag `--socket-fd`: Monitored socket descriptor.
  - Flag `--max-probes`: Number of probes before terminating stall state (default: `5`).
  - Flag `--is-ack-validated`: Checks validity of TCP ACK sequence numbers.
  - Flag `--has-drop-detected`: Set when socket silently stopped responding to probes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpWindowProbes (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      probeSeqNumber INTEGER NOT NULL DEFAULT 0,
      roundTripLatencyUs INTEGER NOT NULL DEFAULT 0,
      isAckValidated INTEGER NOT NULL DEFAULT 1,
      hasDropDetected INTEGER NOT NULL DEFAULT 0,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpWindowProbe struct {
      ID                 string    `json:"id"`
      SocketFd           int       `json:"socketFd"`
      ProbeSeqNumber     int       `json:"probeSeqNumber"`
      RoundTripLatencyUs int       `json:"roundTripLatencyUs"`
      IsAckValidated     bool      `json:"isAckValidated"`
      HasDropDetected    bool      `json:"hasDropDetected"`
      ProbedAt           time.Time `json:"probedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Bound in `pkg/warmup/tcp/probe_synthesizer.go`.

---

### Improvement 1207: Kernel Socket Send-Buffer Auto-Tuning via Congestion Window Feedback
- **ID**: `IMP-1207`
- **Cluster**: 121
- **Goal**: Dynamically tune `SO_SNDBUF` and `TCP_NOTSENT_LOWAT` socket options in response to real-time TCP congestion window (`snd_cwnd`) feedback to prevent local buffer bloat during warming bursts.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss metric bar rendering buffer size dynamic adjustments in spring green (`#00E676`) and current CWND in cyan (`#00E5FF`).
  - *Web*: `<socket-buffer-tuner>` Web Component providing real-time SVG area charts of socket buffer allocation.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline tune-buffer --socket-fd=<fd> [--buffer-size-bytes=<bytes>] [--is-auto-tuned] [--has-buffer-bloat]`
  - Flag `--socket-fd`: Target network socket.
  - Flag `--buffer-size-bytes`: Calculated ideal buffer capacity.
  - Flag `--is-auto-tuned`: Enables automated kernel option mutation.
  - Flag `--has-buffer-bloat`: Flagged when unacknowledged bytes exceed bandwidth-delay product.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketBufferTunings (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      allocatedBufferBytes INTEGER NOT NULL DEFAULT 0,
      cwndSegments INTEGER NOT NULL DEFAULT 0,
      isAutoTuned INTEGER NOT NULL DEFAULT 1,
      hasBufferBloat INTEGER NOT NULL DEFAULT 0,
      tunedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketBufferTuning struct {
      ID                   string    `json:"id"`
      SocketFd             int       `json:"socketFd"`
      AllocatedBufferBytes int       `json:"allocatedBufferBytes"`
      CwndSegments         int       `json:"cwndSegments"`
      IsAutoTuned          bool      `json:"isAutoTuned"`
      HasBufferBloat       bool      `json:"hasBufferBloat"`
      TunedAt              time.Time `json:"tunedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tcp/buffer_tuner.go`.

---

### Improvement 1208: ISP-Specific Pipelined Command Burst Throttling Limiter
- **ID**: `IMP-1208`
- **Cluster**: 121
- **Goal**: Cap the maximum number of pipelined commands per burst according to historical provider thresholds (e.g., Google: 15, Microsoft: 10, Yahoo: 8) to prevent triggering anti-abuse burst filters.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss card showing ISP target domain, active burst ceiling, and current utilization in amber (`#FFD600`).
  - *Web*: `<isp-burst-governor>` Web Component rendering interactive slider limits and provider compliance status.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline set-burst-limit --isp-domain=<domain> [--max-burst=<count>] [--is-adaptive-limit] [--has-enforced-cap]`
  - Flag `--isp-domain`: Target provider domain (e.g., `google.com`, `microsoft.com`).
  - Flag `--max-burst`: Ceiling for command pipelining batch.
  - Flag `--is-adaptive-limit`: Scales down limit automatically on 4xx responses.
  - Flag `--has-enforced-cap`: Enforces hard rejection if queue exceeds burst limit.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspBurstLimits (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL,
      maxBurstCount INTEGER NOT NULL DEFAULT 10,
      effectiveBurstCount INTEGER NOT NULL DEFAULT 10,
      isAdaptiveLimit INTEGER NOT NULL DEFAULT 1,
      hasEnforcedCap INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspBurstLimit struct {
      ID                  string    `json:"id"`
      IspDomain           string    `json:"ispDomain"`
      MaxBurstCount       int       `json:"maxBurstCount"`
      EffectiveBurstCount int       `json:"effectiveBurstCount"`
      IsAdaptiveLimit     bool      `json:"isAdaptiveLimit"`
      HasEnforcedCap      bool      `json:"hasEnforcedCap"`
      UpdatedAt           time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/limiter/burst_governor.go`.

---

### Improvement 1209: Early Data (0-RTT) TLS 1.3 Command Pre-Warming Engine
- **ID**: `IMP-1209`
- **Cluster**: 121
- **Goal**: Transmit initial SMTP `EHLO` and pre-auth handshakes inside TLS 1.3 early data (`0-RTT`) payloads on resumed connections to eliminate an entire handshake round-trip.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss badge in neon green (`#00E676`) denoting active 0-RTT TLS session resumption with handshake latency stats.
  - *Web*: `<tls-early-data-badge>` Web Component rendering cryptographic session ticket indicators.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline resume-tls --host=<host> [--ticket-id=<id>] [--is-early-data-accepted] [--has-resumption-cached]`
  - Flag `--host`: Remote SMTP server host.
  - Flag `--ticket-id`: Serialized TLS 1.3 session ticket identifier.
  - Flag `--is-early-data-accepted`: Indicates remote server acknowledged early data payload.
  - Flag `--has-resumption-cached`: Verifies valid TLS ticket cached in local storage.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TlsEarlyDataResumptions (
      id TEXT PRIMARY KEY,
      remoteHost TEXT NOT NULL,
      ticketId TEXT NOT NULL,
      handshakeSavedMs REAL NOT NULL DEFAULT 0.0,
      isEarlyDataAccepted INTEGER NOT NULL DEFAULT 1,
      hasResumptionCached INTEGER NOT NULL DEFAULT 1,
      resumedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TlsEarlyDataResumption struct {
      ID                  string    `json:"id"`
      RemoteHost          string    `json:"remoteHost"`
      TicketID            string    `json:"ticketId"`
      HandshakeSavedMs    float64   `json:"handshakeSavedMs"`
      IsEarlyDataAccepted bool      `json:"isEarlyDataAccepted"`
      HasResumptionCached bool      `json:"hasResumptionCached"`
      ResumedAt           time.Time `json:"resumedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/tls/early_data.go`.

---

### Improvement 1210: TCP Zero-Window Forensic Audit Capture & Merkle Proof Generator
- **ID**: `IMP-1210`
- **Cluster**: 121
- **Goal**: Cryptographically hash and record TCP zero-window events and socket diagnostics into `repodb/audit.db` with Merkle tree proofs for delivery compliance verification.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss audit block printing SHA-256 Merkle root hashes in muted gray (`#9E9E9E`) and verification badges in gold (`#FFD700`).
  - *Web*: `<audit-merkle-proof>` Web Component allowing one-click cryptographic hash verification.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup audit verify-zero-window --event-id=<id> [--is-merkle-verified] [--has-signature-valid]`
  - Flag `--event-id`: Zero-window occurrence identifier.
  - Flag `--is-merkle-verified`: Verifies proof against local Merkle tree leaf nodes.
  - Flag `--has-signature-valid`: Validates ECDSA signature on audit ledger entry.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowAuditProofs (
      id TEXT PRIMARY KEY,
      eventId TEXT NOT NULL,
      merkleRootHash TEXT NOT NULL,
      leafHash TEXT NOT NULL,
      isMerkleVerified INTEGER NOT NULL DEFAULT 1,
      hasSignatureValid INTEGER NOT NULL DEFAULT 1,
      auditedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowAuditProof struct {
      ID                string    `json:"id"`
      EventID           string    `json:"eventId"`
      MerkleRootHash    string    `json:"merkleRootHash"`
      LeafHash          string    `json:"leafHash"`
      IsMerkleVerified  bool      `json:"isMerkleVerified"`
      HasSignatureValid bool      `json:"hasSignatureValid"`
      AuditedAt         time.Time `json:"auditedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Bound into `pkg/audit/merkle/zero_window.go`.

---

# Cluster 122: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (1211–1220)

---

### Improvement 1211: Multi-ISP Persistent Socket Pool Orchestrator & Reuse Balancer
- **ID**: `IMP-1211`
- **Cluster**: 122 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal**: Maintain an intelligent pool of warm, authenticated TLS SMTP connections per ISP provider, maximizing connection reuse while observing ISP concurrency boundaries.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss dashboard showing connection pool tables with active, idle, and reserved sockets partitioned by ISP domain in cyan (`#00E5FF`).
  - *Web*: `<connection-pool-matrix>` Web Component rendering real-time socket allocation grids.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool list --isp-domain=<domain> [--is-pool-healthy] [--has-available-slot]`
  - Flag `--isp-domain`: Filter by target provider domain.
  - Flag `--is-pool-healthy`: Filter pools by positive health status.
  - Flag `--has-available-slot`: Filter pools with open socket capacity.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspConnectionPools (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL,
      activeConnections INTEGER NOT NULL DEFAULT 0,
      idleConnections INTEGER NOT NULL DEFAULT 0,
      maxAllowedConnections INTEGER NOT NULL DEFAULT 10,
      isPoolHealthy INTEGER NOT NULL DEFAULT 1,
      hasAvailableSlot INTEGER NOT NULL DEFAULT 1,
      lastRecycledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspConnectionPool struct {
      ID                    string    `json:"id"`
      IspDomain             string    `json:"ispDomain"`
      ActiveConnections     int       `json:"activeConnections"`
      IdleConnections       int       `json:"idleConnections"`
      MaxAllowedConnections int       `json:"maxAllowedConnections"`
      IsPoolHealthy         bool      `json:"isPoolHealthy"`
      HasAvailableSlot      bool      `json:"hasAvailableSlot"`
      LastRecycledAt        time.Time `json:"lastRecycledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/pool/orchestrator.go`.

---

### Improvement 1212: Microsecond Socket Heartbeat Probe & Keep-Alive Calibrator
- **ID**: `IMP-1212`
- **Cluster**: 122
- **Goal**: Send non-disruptive `NOOP` or `RSET` heartbeats at dynamically calibrated intervals to prevent stateful NAT firewalls from severing idle pooled connections.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss table row showing heartbeat intervals, last ACK latency in microseconds, and status badge in spring green (`#00E676`).
  - *Web*: `<heartbeat-telemetry>` Web Component displaying microsecond latency charts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool ping-heartbeat --connection-id=<id> [--is-keep-alive-active] [--has-nat-timeout-calibrated]`
  - Flag `--connection-id`: Target pooled socket ID.
  - Flag `--is-keep-alive-active`: Confirms active background ping worker.
  - Flag `--has-nat-timeout-calibrated`: Indicates dynamic calibration against NAT timeout.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketHeartbeats (
      id TEXT PRIMARY KEY,
      connectionId TEXT NOT NULL,
      heartbeatIntervalSec INTEGER NOT NULL DEFAULT 30,
      responseLatencyUs INTEGER NOT NULL DEFAULT 0,
      isKeepAliveActive INTEGER NOT NULL DEFAULT 1,
      hasNatTimeoutCalibrated INTEGER NOT NULL DEFAULT 1,
      pingedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketHeartbeat struct {
      ID                      string    `json:"id"`
      ConnectionID            string    `json:"connectionId"`
      HeartbeatIntervalSec    int       `json:"heartbeatIntervalSec"`
      ResponseLatencyUs       int       `json:"responseLatencyUs"`
      IsKeepAliveActive       bool      `json:"isKeepAliveActive"`
      HasNatTimeoutCalibrated bool      `json:"hasNatTimeoutCalibrated"`
      PingedAt                time.Time `json:"pingedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/pool/heartbeat.go`.

---

### Improvement 1213: Real-Time Socket Degradation Scorer & Preemptive Socket Evictor
- **ID**: `IMP-1213`
- **Cluster**: 122
- **Goal**: Compute a continuous health degradation score (0–100) based on packet retransmissions, TCP RTT variance, and slow response codes, gracefully evicting decaying sockets before failures occur.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss progress bar displaying socket health score: green (`#00E676`) for > 80, amber (`#FFD600`) for 50–80, red (`#FF1744`) for < 50.
  - *Web*: `<socket-health-card>` Web Component displaying degradation factor breakdown.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool score-socket --connection-id=<id> [--is-eviction-recommended] [--has-high-retransmits]`
  - Flag `--connection-id`: Socket identifier.
  - Flag `--is-eviction-recommended`: True if health score falls below eviction threshold.
  - Flag `--has-high-retransmits`: Flagged when retransmit percentage exceeds 2%.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketDegradationScores (
      id TEXT PRIMARY KEY,
      connectionId TEXT NOT NULL,
      healthScore INTEGER NOT NULL DEFAULT 100,
      retransmissionRate REAL NOT NULL DEFAULT 0.0,
      isEvictionRecommended INTEGER NOT NULL DEFAULT 0,
      hasHighRetransmits INTEGER NOT NULL DEFAULT 0,
      scoredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketDegradationScore struct {
      ID                    string    `json:"id"`
      ConnectionID          string    `json:"connectionId"`
      HealthScore           int       `json:"healthScore"`
      RetransmissionRate    float64   `json:"retransmissionRate"`
      IsEvictionRecommended bool      `json:"isEvictionRecommended"`
      HasHighRetransmits    bool      `json:"hasHighRetransmits"`
      ScoredAt              time.Time `json:"scoredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Built into `pkg/warmup/pool/scorer.go`.

---

### Improvement 1214: Multi-MTA Concurrent Handshake Shaper & Staggered Syn Dispatcher
- **ID**: `IMP-1214`
- **Cluster**: 122
- **Goal**: Stagger concurrent TCP `SYN` packets across destination mail exchangers using jittered Poisson intervals, preventing synchronized spikes that trigger ISP edge rate limiters.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss timeline displaying staggered connection attempts with millisecond delta offsets in purple (`#7C4DFF`).
  - *Web*: `<staggered-syn-visualizer>` Web Component rendering interactive scatter plots of connection initiations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool stagger-syn --isp-domain=<domain> [--stagger-interval-ms=<ms>] [--is-poisson-distributed] [--has-syn-flood-prevented]`
  - Flag `--isp-domain`: Target provider domain.
  - Flag `--stagger-interval-ms`: Base stagger delay (default: `120`).
  - Flag `--is-poisson-distributed`: Uses Poisson distribution for jittered intervals.
  - Flag `--has-syn-flood-prevented`: Confirms outbound SYN rate limiter active.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SynStaggerEvents (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL,
      staggerDelayMs INTEGER NOT NULL DEFAULT 0,
      activeHandshakeCount INTEGER NOT NULL DEFAULT 0,
      isPoissonDistributed INTEGER NOT NULL DEFAULT 1,
      hasSynFloodPrevented INTEGER NOT NULL DEFAULT 1,
      dispatchedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SynStaggerEvent struct {
      ID                   string    `json:"id"`
      IspDomain            string    `json:"ispDomain"`
      StaggerDelayMs       int       `json:"staggerDelayMs"`
      ActiveHandshakeCount int       `json:"activeHandshakeCount"`
      IsPoissonDistributed bool      `json:"isPoissonDistributed"`
      HasSynFloodPrevented bool      `json:"hasSynFloodPrevented"`
      DispatchedAt         time.Time `json:"dispatchedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/syn_shaper.go`.

---

### Improvement 1215: Cross-ISP Handshake Latency Matrix & Edge Routing Benchmark
- **ID**: `IMP-1215`
- **Cluster**: 122
- **Goal**: Maintain an in-memory matrix of handshake latencies (DNS, TCP, TLS) across all major providers, dynamically picking the lowest-latency edge egress node.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss matrix table rendering providers on rows and latency phases on columns with gradient heat map styling.
  - *Web*: `<latency-matrix-heat-map>` Web Component rendering interactive SVG color ramps.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool latency-matrix [--is-optimal-selected] [--has-stale-benchmarks]`
  - Flag `--is-optimal-selected`: Highlights lowest latency egress route.
  - Flag `--has-stale-benchmarks`: Indicates benchmarks older than TTL.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HandshakeLatencyMatrices (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL,
      dnsLatencyMs REAL NOT NULL DEFAULT 0.0,
      tcpLatencyMs REAL NOT NULL DEFAULT 0.0,
      tlsLatencyMs REAL NOT NULL DEFAULT 0.0,
      isOptimalSelected INTEGER NOT NULL DEFAULT 1,
      hasStaleBenchmarks INTEGER NOT NULL DEFAULT 0,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HandshakeLatencyMatrix struct {
      ID                 string    `json:"id"`
      IspDomain          string    `json:"ispDomain"`
      DnsLatencyMs       float64   `json:"dnsLatencyMs"`
      TcpLatencyMs       float64   `json:"tcpLatencyMs"`
      TlsLatencyMs       float64   `json:"tlsLatencyMs"`
      IsOptimalSelected  bool      `json:"isOptimalSelected"`
      HasStaleBenchmarks bool      `json:"hasStaleBenchmarks"`
      MeasuredAt         time.Time `json:"measuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/pool/latency_matrix.go`.

---

### Improvement 1216: Graceful Idle Socket Retirement & TCP FIN Handshake Finalizer
- **ID**: `IMP-1216`
- **Cluster**: 122
- **Goal**: Gracefully retire aged pooled connections using standard `QUIT` and bidirectional TCP `FIN` exchanges instead of abrupt `RST` teardowns, preserving remote MTA trust scores.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status notification in electric cyan (`#00E5FF`) confirming clean TCP FIN handshake completion.
  - *Web*: `<socket-teardown-badge>` Web Component displaying termination state.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool retire-socket --connection-id=<id> [--is-clean-fin-received] [--has-quit-acknowledged]`
  - Flag `--connection-id`: Target socket ID.
  - Flag `--is-clean-fin-received`: Confirms receipt of remote FIN ACK.
  - Flag `--has-quit-acknowledged`: Confirms 221 response received to QUIT command.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketRetirements (
      id TEXT PRIMARY KEY,
      connectionId TEXT NOT NULL,
      socketLifetimeSec INTEGER NOT NULL DEFAULT 0,
      messagesProcessed INTEGER NOT NULL DEFAULT 0,
      isCleanFinReceived INTEGER NOT NULL DEFAULT 1,
      hasQuitAcknowledged INTEGER NOT NULL DEFAULT 1,
      retiredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketRetirement struct {
      ID                  string    `json:"id"`
      ConnectionID        string    `json:"connectionId"`
      SocketLifetimeSec   int       `json:"socketLifetimeSec"`
      MessagesProcessed   int       `json:"messagesProcessed"`
      IsCleanFinReceived  bool      `json:"isCleanFinReceived"`
      HasQuitAcknowledged bool      `json:"hasQuitAcknowledged"`
      RetiredAt           time.Time `json:"retiredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/pool/finalizer.go`.

---

### Improvement 1217: Dynamic IP-Binding Interface Affinity & Local Address Rotator
- **ID**: `IMP-1217`
- **Cluster**: 122
- **Goal**: Bind outbound sockets across multiple local network interfaces and IP addresses based on destination reputation cohorts, balancing load evenly across multi-homed NICs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss card displaying interface name, bound local IP, egress throughput, and active socket tally.
  - *Web*: `<interface-affinity-map>` Web Component displaying network interface topology.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool bind-interface --local-ip=<ip> [--interface-name=<name>] [--is-binding-active] [--has-ip-healthy]`
  - Flag `--local-ip`: Outbound source IP address.
  - Flag `--interface-name`: Target network adapter.
  - Flag `--is-binding-active`: Confirms socket bound to specific interface.
  - Flag `--has-ip-healthy`: Verifies IP not blacklisted.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS InterfaceBindings (
      id TEXT PRIMARY KEY,
      interfaceName TEXT NOT NULL,
      localIp TEXT NOT NULL,
      boundSocketCount INTEGER NOT NULL DEFAULT 0,
      isBindingActive INTEGER NOT NULL DEFAULT 1,
      hasIpHealthy INTEGER NOT NULL DEFAULT 1,
      boundAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type InterfaceBinding struct {
      ID               string    `json:"id"`
      InterfaceName    string    `json:"interfaceName"`
      LocalIP          string    `json:"localIp"`
      BoundSocketCount int       `json:"boundSocketCount"`
      IsBindingActive  bool      `json:"isBindingActive"`
      HasIPHealthy     bool      `json:"hasIpHealthy"`
      BoundAt          time.Time `json:"boundAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/pool/interface_binder.go`.

---

### Improvement 1218: Transport Security Cipher Negotiation Matrix & Downgrade Shield
- **ID**: `IMP-1218`
- **Cluster**: 122
- **Goal**: Enforce strict TLS 1.3/1.2 cipher prioritization, shielding connections against cryptographic downgrade attacks while maintaining compatibility across legacy MTA gateways.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss security card rendering cipher suite strings in bright purple (`#7C4DFF`) with forward-secrecy verification marks in spring green (`#00E676`).
  - *Web*: `<tls-cipher-matrix>` Web Component displaying negotiated protocol breakdown.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool check-ciphers --isp-domain=<domain> [--is-tls13-preferred] [--has-forward-secrecy]`
  - Flag `--isp-domain`: Target provider domain.
  - Flag `--is-tls13-preferred`: Verifies TLS 1.3 negotiated when available.
  - Flag `--has-forward-secrecy`: Ensures ECDHE cipher suites are enforced.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CipherNegotiations (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL,
      tlsVersion TEXT NOT NULL,
      negotiatedCipherSuite TEXT NOT NULL,
      isTls13Preferred INTEGER NOT NULL DEFAULT 1,
      hasForwardSecrecy INTEGER NOT NULL DEFAULT 1,
      inspectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CipherNegotiation struct {
      ID                    string    `json:"id"`
      IspDomain             string    `json:"ispDomain"`
      TlsVersion            string    `json:"tlsVersion"`
      NegotiatedCipherSuite string    `json:"negotiatedCipherSuite"`
      IsTls13Preferred      bool      `json:"isTls13Preferred"`
      HasForwardSecrecy     bool      `json:"hasForwardSecrecy"`
      InspectedAt           time.Time `json:"inspectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Bound into `pkg/warmup/tls/cipher_matrix.go`.

---

### Improvement 1219: ISP Connection Limit Governor & Global Socket Ceiling Enforcer
- **ID**: `IMP-1219`
- **Cluster**: 122
- **Goal**: Monitor and enforce global concurrent connection limits across distributed sender nodes per ISP destination, preventing simultaneous connections from tripping provider firewall thresholds.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss gauge displaying global active connections vs. ISP ceiling in yellow amber (`#FFD600`).
  - *Web*: `<isp-connection-ceiling>` Web Component rendering real-time capacity thermometers.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool enforce-ceiling --isp-domain=<domain> [--ceiling=<limit>] [--is-ceiling-exceeded] [--has-global-consensus]`
  - Flag `--isp-domain`: Target provider domain.
  - Flag `--ceiling`: Maximum concurrent open sockets.
  - Flag `--is-ceiling-exceeded`: True if current socket count equals or exceeds ceiling.
  - Flag `--has-global-consensus`: Validated against distributed cluster state.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspConnectionCeilings (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL,
      maxAllowedSockets INTEGER NOT NULL DEFAULT 20,
      currentActiveSockets INTEGER NOT NULL DEFAULT 0,
      isCeilingExceeded INTEGER NOT NULL DEFAULT 0,
      hasGlobalConsensus INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspConnectionCeiling struct {
      ID                   string    `json:"id"`
      IspDomain            string    `json:"ispDomain"`
      MaxAllowedSockets    int       `json:"maxAllowedSockets"`
      CurrentActiveSockets int       `json:"currentActiveSockets"`
      IsCeilingExceeded    bool      `json:"isCeilingExceeded"`
      HasGlobalConsensus   bool      `json:"hasGlobalConsensus"`
      RecordedAt           time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/pool/governor.go`.

---

### Improvement 1220: Connection Pool Forensic Audit & State Transition Ledger
- **ID**: `IMP-1220`
- **Cluster**: 122
- **Goal**: Record all socket state transitions (allocation, reuse, degradation, eviction) into `repodb/audit.db` with cryptographic sequence hashing for forensic performance review.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss audit log format displaying microsecond timestamps and transition events in muted cyan (`#80DEEA`).
  - *Web*: `<pool-audit-stream>` Web Component displaying interactive event timelines.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup audit verify-pool-transitions --connection-id=<id> [--is-audit-valid] [--has-sequence-intact]`
  - Flag `--connection-id`: Target connection ID.
  - Flag `--is-audit-valid`: Confirms audit hash verification pass.
  - Flag `--has-sequence-intact`: Verifies monotonic sequence numbers.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolStateAuditRecords (
      id TEXT PRIMARY KEY,
      connectionId TEXT NOT NULL,
      previousState TEXT NOT NULL,
      newState TEXT NOT NULL,
      stateHash TEXT NOT NULL,
      isAuditValid INTEGER NOT NULL DEFAULT 1,
      hasSequenceIntact INTEGER NOT NULL DEFAULT 1,
      transitionedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolStateAuditRecord struct {
      ID                string    `json:"id"`
      ConnectionID      string    `json:"connectionId"`
      PreviousState     string    `json:"previousState"`
      NewState          string    `json:"newState"`
      StateHash         string    `json:"stateHash"`
      IsAuditValid      bool      `json:"isAuditValid"`
      HasSequenceIntact bool      `json:"hasSequenceIntact"`
      TransitionedAt    time.Time `json:"transitionedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/audit/ledger/pool_audit.go`.

---

# Cluster 123: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1221–1230)

---

### Improvement 1221: Bubble Tea v2 HUD Viewport & Multi-Panel Egress Terminal Workspace
- **ID**: `IMP-1221`
- **Cluster**: 123 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Architect a high-frame-rate terminal HUD using Bubble Tea v2 and Lipgloss layout primitives, rendering synchronized sub-viewports for socket status, transmission throughput, and jitter waveforms.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Multi-pane Lipgloss TUI featuring 3 bounded horizontal panes, bordered in deep slate (`#37474F`) with active focus highlighted in electric cyan (`#00E5FF`).
  - *Web*: `<egress-hud-workspace>` Web Component rendering flexbox panes replicating terminal geometry.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui hud [--fps=<rate>] [--is-hud-fullscreen] [--has-mouse-enabled]`
  - Flag `--fps`: Rendering frames per second (default: `30`).
  - Flag `--is-hud-fullscreen`: Runs HUD in alternate screen buffer.
  - Flag `--has-mouse-enabled`: Captures mouse scroll and click events.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HudViewportStates (
      id TEXT PRIMARY KEY,
      sessionToken TEXT NOT NULL,
      activePaneIndex INTEGER NOT NULL DEFAULT 0,
      renderFps INTEGER NOT NULL DEFAULT 30,
      isHudFullscreen INTEGER NOT NULL DEFAULT 1,
      hasMouseEnabled INTEGER NOT NULL DEFAULT 1,
      initializedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HudViewportState struct {
      ID              string    `json:"id"`
      SessionToken    string    `json:"sessionToken"`
      ActivePaneIndex int       `json:"activePaneIndex"`
      RenderFps       int       `json:"renderFps"`
      IsHudFullscreen bool      `json:"isHudFullscreen"`
      HasMouseEnabled bool      `json:"hasMouseEnabled"`
      InitializedAt   time.Time `json:"initializedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/ui/hud/viewport.go`.

---

### Improvement 1222: Sub-Millisecond SMTP Phase Latency Flamegraph Renderer
- **ID**: `IMP-1222`
- **Cluster**: 123
- **Goal**: Measure and visualize DNS, TCP, TLS, Pipelining, and DATA payload phases as a horizontal color-coded flamegraph in both terminal ANSI TrueColor and inline SVG formats.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Horizontal flamegraph bar: DNS in cyan (`#00E5FF`), TCP in spring green (`#00E676`), TLS in purple (`#7C4DFF`), PIPE in amber (`#FFD600`), and DATA in magenta (`#E040FB`).
  - *Web*: Native `<smtp-flamegraph>` Web Component rendering interactive SVG rect elements with hover duration tooltips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui flamegraph --transaction-id=<id> [--is-flamegraph-detailed] [--has-svg-export]`
  - Flag `--transaction-id`: Target SMTP transaction ID.
  - Flag `--is-flamegraph-detailed`: Renders microsecond labels on each phase block.
  - Flag `--has-svg-export`: Writes inline SVG string to output buffer.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PhaseFlamegraphLogs (
      id TEXT PRIMARY KEY,
      transactionId TEXT NOT NULL,
      totalDurationMs REAL NOT NULL DEFAULT 0.0,
      dnsDurationMs REAL NOT NULL DEFAULT 0.0,
      tcpDurationMs REAL NOT NULL DEFAULT 0.0,
      tlsDurationMs REAL NOT NULL DEFAULT 0.0,
      pipeDurationMs REAL NOT NULL DEFAULT 0.0,
      dataDurationMs REAL NOT NULL DEFAULT 0.0,
      isFlamegraphDetailed INTEGER NOT NULL DEFAULT 1,
      hasSvgExport INTEGER NOT NULL DEFAULT 0,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PhaseFlamegraphLog struct {
      ID                   string    `json:"id"`
      TransactionID        string    `json:"transactionId"`
      TotalDurationMs      float64   `json:"totalDurationMs"`
      DnsDurationMs        float64   `json:"dnsDurationMs"`
      TcpDurationMs        float64   `json:"tcpDurationMs"`
      TlsDurationMs        float64   `json:"tlsDurationMs"`
      PipeDurationMs       float64   `json:"pipeDurationMs"`
      DataDurationMs       float64   `json:"dataDurationMs"`
      IsFlamegraphDetailed bool      `json:"isFlamegraphDetailed"`
      HasSvgExport         bool      `json:"hasSvgExport"`
      CapturedAt           time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Built into `pkg/ui/flamegraph/renderer.go`.

---

### Improvement 1223: Braille & Block-Based ASCII Jitter Waveform Oscilloscope
- **ID**: `IMP-1223`
- **Cluster**: 123
- **Goal**: Render real-time TCP ACK jitter waveforms using Unicode Braille (`⠋⠙⠹...`) and quarter-block characters in the terminal, visualizing millisecond timing jitter with oscilloscope precision.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss oscilloscope screen in dark blue (`#0D47A1`) with green Braille dots (`#00E676`) charting ACK jitter across a sliding 60-second window.
  - *Web*: `<jitter-oscilloscope>` Web Component rendering HTML5 canvas equivalent waveform.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui waveform --window-sec=<sec> [--is-braille-mode] [--has-grid-lines]`
  - Flag `--window-sec`: Duration of sliding window in seconds (default: `60`).
  - Flag `--is-braille-mode`: Uses high-density Unicode Braille characters.
  - Flag `--has-grid-lines`: Draws horizontal latency threshold guide marks.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterWaveformConfigs (
      id TEXT PRIMARY KEY,
      viewportWidth INTEGER NOT NULL DEFAULT 80,
      windowDurationSec INTEGER NOT NULL DEFAULT 60,
      isBrailleMode INTEGER NOT NULL DEFAULT 1,
      hasGridLines INTEGER NOT NULL DEFAULT 1,
      configuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterWaveformConfig struct {
      ID                string    `json:"id"`
      ViewportWidth     int       `json:"viewportWidth"`
      WindowDurationSec int       `json:"windowDurationSec"`
      IsBrailleMode     bool      `json:"isBrailleMode"`
      HasGridLines      bool      `json:"hasGridLines"`
      ConfiguredAt      time.Time `json:"configuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/ui/waveform/oscilloscope.go`.

---

### Improvement 1224: Zero-CDN Shadow DOM Web Component Telemetry Dashboard
- **ID**: `IMP-1224`
- **Cluster**: 123
- **Goal**: Package self-contained, standalone HTML/JS/CSS Web Components encapsulating SVG sparklines and latency gauges with zero external CDN dependencies, operable completely offline.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Command output generating self-contained HTML payloads directly to standard output or local file.
  - *Web*: Custom Element `<gitmap-warmup-dashboard>` encapsulating encapsulated shadow DOM tree with zero external script or font requests.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui export-dashboard --output-path=<path> [--is-zero-cdn-verified] [--has-embedded-assets]`
  - Flag `--output-path`: Target HTML file path.
  - Flag `--is-zero-cdn-verified`: Validates that zero external network resources are referenced.
  - Flag `--has-embedded-assets`: Embeds all styling and SVG templates inline.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WebDashboardExports (
      id TEXT PRIMARY KEY,
      outputPath TEXT NOT NULL,
      bundleSizeBytes INTEGER NOT NULL DEFAULT 0,
      isZeroCdnVerified INTEGER NOT NULL DEFAULT 1,
      hasEmbeddedAssets INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WebDashboardExport struct {
      ID                string    `json:"id"`
      OutputPath        string    `json:"outputPath"`
      BundleSizeBytes   int       `json:"bundleSizeBytes"`
      IsZeroCdnVerified bool      `json:"isZeroCdnVerified"`
      HasEmbeddedAssets bool      `json:"hasEmbeddedAssets"`
      ExportedAt        time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/ui/export/web_component.go`.

---

### Improvement 1225: Lipgloss Color Token Management & Adaptive Terminal Palette Synchronizer
- **ID**: `IMP-1225`
- **Cluster**: 123
- **Goal**: Dynamically detect terminal color capabilities (TrueColor, 256-color, 16-color) and synchronize Lipgloss styling tokens across dark and light terminal backgrounds.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Color swatch matrix showing resolved palette tokens mapped to local terminal profile.
  - *Web*: Synchronized CSS root variables adhering to same hexadecimal values.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui sync-palette [--is-true-color-detected] [--has-dark-mode-active]`
  - Flag `--is-true-color-detected`: Indicates 24-bit TrueColor support verified.
  - Flag `--has-dark-mode-active`: Resolves high-contrast tokens for dark terminals.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TerminalPalettes (
      id TEXT PRIMARY KEY,
      colorDepthBits INTEGER NOT NULL DEFAULT 24,
      backgroundProfile TEXT NOT NULL DEFAULT 'dark',
      isTrueColorDetected INTEGER NOT NULL DEFAULT 1,
      hasDarkModeActive INTEGER NOT NULL DEFAULT 1,
      detectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TerminalPalette struct {
      ID                  string    `json:"id"`
      ColorDepthBits      int       `json:"colorDepthBits"`
      BackgroundProfile   string    `json:"backgroundProfile"`
      IsTrueColorDetected bool      `json:"isTrueColorDetected"`
      HasDarkModeActive   bool      `json:"hasDarkModeActive"`
      DetectedAt          time.Time `json:"detectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/ui/theme/palette.go`.

---

### Improvement 1226: Interactive Bubble Tea Keybinding Controller & Filter Navigation
- **ID**: `IMP-1226`
- **Cluster**: 123
- **Goal**: Implement intuitive, responsive keyboard navigation (`j`/`k` for scrolling, `/` for ISP search, `tab` for pane switching, `q` for quit) within the Bubble Tea v2 HUD.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Persistent footer bar displaying available hotkeys in muted slate (`#90A4AE`) with key badges highlighted in cyan (`#00E5FF`).
  - *Web*: Virtual keyboard legend rendered in footer status bar.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui config-keys [--is-vim-mode-active] [--has-help-overlay-visible]`
  - Flag `--is-vim-mode-active`: Enables hjkl navigation binding.
  - Flag `--has-help-overlay-visible`: Toggles interactive hotkey overlay.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS KeybindingPreferences (
      id TEXT PRIMARY KEY,
      userProfile TEXT NOT NULL DEFAULT 'default',
      isVimModeActive INTEGER NOT NULL DEFAULT 1,
      hasHelpOverlayVisible INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type KeybindingPreference struct {
      ID                    string    `json:"id"`
      UserProfile           string    `json:"userProfile"`
      IsVimModeActive       bool      `json:"isVimModeActive"`
      HasHelpOverlayVisible bool      `json:"hasHelpOverlayVisible"`
      UpdatedAt             time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Bound into `pkg/ui/hud/keybindings.go`.

---

### Improvement 1227: High-Frequency Telemetry Stream Ring-Buffer & Viewport Throttler
- **ID**: `IMP-1227`
- **Cluster**: 123
- **Goal**: Buffer up to 50,000 telemetry events per second in a lock-free memory ring buffer, downsampling frame emissions to 30/60 FPS to prevent terminal emulator rendering lockup.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Buffer capacity meter in bottom HUD status line displaying throughput (events/sec) and frame drop count in green (`#00E676`).
  - *Web*: `<ring-buffer-meter>` Web Component rendering buffer utilization sparkline.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui stream-buffer [--buffer-size=<size>] [--is-downsampling-active] [--has-frame-dropped]`
  - Flag `--buffer-size`: Ring buffer event capacity (default: `65536`).
  - Flag `--is-downsampling-active`: True if event downsampler active for UI thread.
  - Flag `--has-frame-dropped`: Flagged when UI frame skips to maintain latency.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS UiBufferMetrics (
      id TEXT PRIMARY KEY,
      incomingEventRate INTEGER NOT NULL DEFAULT 0,
      renderedFrameRate INTEGER NOT NULL DEFAULT 60,
      isDownsamplingActive INTEGER NOT NULL DEFAULT 1,
      hasFrameDropped INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type UiBufferMetric struct {
      ID                   string    `json:"id"`
      IncomingEventRate    int       `json:"incomingEventRate"`
      RenderedFrameRate    int       `json:"renderedFrameRate"`
      IsDownsamplingActive bool      `json:"isDownsamplingActive"`
      HasFrameDropped      bool      `json:"hasFrameDropped"`
      RecordedAt           time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/ui/hud/ring_buffer.go`.

---

### Improvement 1228: SVG Flamegraph Snapshot Generator & File Exporter
- **ID**: `IMP-1228`
- **Cluster**: 123
- **Goal**: Export point-in-time SVG flamegraph snapshots to disk with embedded CSS styling and metadata, allowing developers to share diagnostic traces without dependencies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Notification card displaying exported SVG file path and file size in bytes with checkmark in spring green (`#00E676`).
  - *Web*: One-click "Download SVG Snapshot" button triggering local blob download.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui snapshot-svg --transaction-id=<id> --file-path=<path> [--is-compressed] [--has-embedded-metadata]`
  - Flag `--transaction-id`: Transaction to snapshot.
  - Flag `--file-path`: Destination SVG path.
  - Flag `--is-compressed`: Enables SVGZ gzip compression.
  - Flag `--has-embedded-metadata`: Encodes transaction UUID in SVG XML comment.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SvgSnapshotFiles (
      id TEXT PRIMARY KEY,
      transactionId TEXT NOT NULL,
      filePath TEXT NOT NULL,
      fileSizeBytes INTEGER NOT NULL DEFAULT 0,
      isCompressed INTEGER NOT NULL DEFAULT 0,
      hasEmbeddedMetadata INTEGER NOT NULL DEFAULT 1,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SvgSnapshotFile struct {
      ID                  string    `json:"id"`
      TransactionID       string    `json:"transactionId"`
      FilePath            string    `json:"filePath"`
      FileSizeBytes       int       `json:"fileSizeBytes"`
      IsCompressed        bool      `json:"isCompressed"`
      HasEmbeddedMetadata bool      `json:"hasEmbeddedMetadata"`
      GeneratedAt         time.Time `json:"generatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/ui/flamegraph/snapshot.go`.

---

### Improvement 1229: Multi-Provider Jitter Delta Comparator & Anomaly Highlighting
- **ID**: `IMP-1229`
- **Cluster**: 123
- **Goal**: Compare real-time jitter waveforms across multiple ISP targets simultaneously, highlighting anomalous statistical divergence (> 3 standard deviations) with visual alert tokens.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Side-by-side jitter differential graph with anomalous divergence points flashing in vibrant red (`#FF1744`).
  - *Web*: `<jitter-comparator>` Web Component displaying multi-line overlay charts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui compare-jitter --isp-primary=<p> --isp-secondary=<s> [--is-divergence-alert] [--has-baseline-exceeded]`
  - Flag `--isp-primary`: Primary comparison provider.
  - Flag `--isp-secondary`: Secondary comparison provider.
  - Flag `--is-divergence-alert`: True if divergence exceeds 3-sigma.
  - Flag `--has-baseline-exceeded`: Indicates jitter exceeds historical baseline.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterComparisons (
      id TEXT PRIMARY KEY,
      primaryIsp TEXT NOT NULL,
      secondaryIsp TEXT NOT NULL,
      deltaMeanUs INTEGER NOT NULL DEFAULT 0,
      isDivergenceAlert INTEGER NOT NULL DEFAULT 0,
      hasBaselineExceeded INTEGER NOT NULL DEFAULT 0,
      comparedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterComparison struct {
      ID                  string    `json:"id"`
      PrimaryIsp          string    `json:"primaryIsp"`
      SecondaryIsp        string    `json:"secondaryIsp"`
      DeltaMeanUs         int       `json:"deltaMeanUs"`
      IsDivergenceAlert   bool      `json:"isDivergenceAlert"`
      HasBaselineExceeded bool      `json:"hasBaselineExceeded"`
      ComparedAt          time.Time `json:"comparedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/ui/waveform/comparator.go`.

---

### Improvement 1230: TUI Session Diagnostic Forensic Audit & Trace Hasher
- **ID**: `IMP-1230`
- **Cluster**: 123
- **Goal**: Record interactive terminal UI diagnostic sessions, flamegraph summaries, and rendered anomalies into `repodb/audit.db` with cryptographic session hashes for reproducible debugging.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Audit verification badge printed upon exiting TUI HUD with session UUID and SHA-256 trace hash.
  - *Web*: `<tui-audit-card>` Web Component displaying past interactive diagnostic sessions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup audit verify-tui-trace --session-id=<id> [--is-trace-sealed] [--has-audit-hash-matched]`
  - Flag `--session-id`: UI diagnostic session ID.
  - Flag `--is-trace-sealed`: Confirms session trace cryptographically closed.
  - Flag `--has-audit-hash-matched`: Validates hash integrity against ledger record.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiSessionAuditTraces (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      terminalGeometry TEXT NOT NULL,
      eventCount INTEGER NOT NULL DEFAULT 0,
      traceHash TEXT NOT NULL,
      isTraceSealed INTEGER NOT NULL DEFAULT 1,
      hasAuditHashMatched INTEGER NOT NULL DEFAULT 1,
      loggedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiSessionAuditTrace struct {
      ID                  string    `json:"id"`
      SessionID           string    `json:"sessionId"`
      TerminalGeometry    string    `json:"terminalGeometry"`
      EventCount          int       `json:"eventCount"`
      TraceHash           string    `json:"traceHash"`
      IsTraceSealed       bool      `json:"isTraceSealed"`
      HasAuditHashMatched bool      `json:"hasAuditHashMatched"`
      LoggedAt            time.Time `json:"loggedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Bound into `pkg/audit/ledger/tui_audit.go`.
