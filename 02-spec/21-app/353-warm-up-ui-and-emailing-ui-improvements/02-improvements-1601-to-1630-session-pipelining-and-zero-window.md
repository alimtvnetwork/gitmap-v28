# Improvements 1601 to 1630: Session Pipelining, Socket Health & Jitter Waveform TUI

This document details the complete technical specifications for Improvements 1601 through 1630 across Clusters 161, 162, and 163 of Sequence 353.

---

# Cluster 161: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (1601–1610)

---

### Improvement 1601: Dynamic RFC 2920 Multi-Command Pipelining Orchestrator
- **ID**: `IMP-1601`
- **Cluster**: 161 (Autonomous SMTP Session Pipelining & TCP Zero-Window Detection)
- **Goal**: Coalesce outbound SMTP commands (`RSET`, `MAIL FROM`, `RCPT TO`, `DATA`) into single TCP payload packets to eliminate sequential round-trip wait times without triggering firewall stream inspection aborts.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Dual-column Lipgloss pane displaying queued command vectors in cyan (`#00E5FF`) and demultiplexed asynchronous response vectors in emerald green (`#00E676`).
  - *Web*: Zero-CDN `<smtp-pipeline-coalescer>` custom Web Component rendering real-time command sequence diagrams via SVG vector streams.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline coalesce --session-id=<id> [--batch-size=<int>] [--is-adaptive-coalesce] [--has-early-dispatch]`
  - Flag `--session-id`: Target SMTP session UUID.
  - Flag `--batch-size`: Max commands batched per TCP segment (default: `32`).
  - Flag `--is-adaptive-coalesce`: Automatically dynamically resizes batch based on remote TCP window headroom.
  - Flag `--has-early-dispatch`: Sends next message envelope before receiving 250 OK for preceding DATA end-of-dot.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpPipelinedSessions (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      batchSize INTEGER NOT NULL DEFAULT 16,
      commandsCoalesced INTEGER NOT NULL DEFAULT 0,
      isAdaptiveCoalesce INTEGER NOT NULL DEFAULT 1,
      hasEarlyDispatch INTEGER NOT NULL DEFAULT 0,
      pipelinedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpPipelinedSession struct {
      ID                 string    `json:"id"`
      SessionID          string    `json:"sessionId"`
      BatchSize          int       `json:"batchSize"`
      CommandsCoalesced  int       `json:"commandsCoalesced"`
      IsAdaptiveCoalesce bool      `json:"isAdaptiveCoalesce"`
      HasEarlyDispatch   bool      `json:"hasEarlyDispatch"`
      PipelinedAt        time.Time `json:"pipelinedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/smtp/pipeliner.go` within outbound transmission stream dispatchers.

---

### Improvement 1602: Zero-Copy Kernel TCP Window Telemetry Reflection
- **ID**: `IMP-1602`
- **Cluster**: 161
- **Goal**: Sample OS TCP socket buffer state via `getsockopt` with `TCP_INFO` in microsecond loops, extracting `tcpi_snd_wnd` and `tcpi_rtt` without context-switching memory copies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss gauge displaying advertised receiver window with color thresholds: `#00E676` (>64KB), `#FFD600` (8KB–64KB), `#FF1744` (<8KB).
  - *Web*: `<tcp-window-telemetry>` Web Component with live SVG SVG circular dials and microsecond tick graphs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tcp-window inspect --socket-fd=<fd> [--interval-us=<us>] [--is-realtime-stream] [--has-buffer-alert]`
  - Flag `--socket-fd`: Monitored socket file descriptor.
  - Flag `--interval-us`: Sampling interval in microseconds (default: `250`).
  - Flag `--is-realtime-stream`: Streams continuous window telemetry updates to stdout.
  - Flag `--has-buffer-alert`: Emits alert log when buffer drops below 2x Maximum Segment Size (MSS).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpZeroWindowSnapshots (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      advertisedWindowBytes INTEGER NOT NULL DEFAULT 0,
      rttMicroseconds INTEGER NOT NULL DEFAULT 0,
      isRealtimeStream INTEGER NOT NULL DEFAULT 0,
      hasBufferAlert INTEGER NOT NULL DEFAULT 0,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpZeroWindowSnapshot struct {
      ID                    string    `json:"id"`
      SocketFd              int       `json:"socketFd"`
      AdvertisedWindowBytes int       `json:"advertisedWindowBytes"`
      RttMicroseconds       int       `json:"rttMicroseconds"`
      IsRealtimeStream      bool      `json:"isRealtimeStream"`
      HasBufferAlert        bool      `json:"hasBufferAlert"`
      SampledAt             time.Time `json:"sampledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tcp/socket_telemetry.go` for low-level socket state monitoring.

---

### Improvement 1603: Autonomous TCP Zero-Window Probing & Preemptive Flow Restraint
- **ID**: `IMP-1603`
- **Cluster**: 161
- **Goal**: Detect remote window collapse (`snd_wnd == 0`) and engage microsecond zero-window keep-alive probes without breaking active SMTP transactions or triggering timeout drops.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Pulsing amber banner (`#FFD600`) rendered via Lipgloss alert boxes indicating zero-window probe in progress with countdown timer.
  - *Web*: `<tcp-zero-window-probe>` Web Component rendering an animated pulse node with backoff retry progression.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tcp-window probe --socket-fd=<fd> [--probe-timeout-ms=<ms>] [--is-probe-active] [--has-exponential-backoff]`
  - Flag `--socket-fd`: Monitored socket descriptor.
  - Flag `--probe-timeout-ms`: Maximum allowed freeze duration before connection reallocation (default: `15000`).
  - Flag `--is-probe-active`: Indicates probe is currently active.
  - Flag `--has-exponential-backoff`: Doubles probe interval after each probe ACK.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpZeroWindowProbes (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      probeCount INTEGER NOT NULL DEFAULT 0,
      durationMs REAL NOT NULL DEFAULT 0.0,
      isProbeActive INTEGER NOT NULL DEFAULT 1,
      hasExponentialBackoff INTEGER NOT NULL DEFAULT 1,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpZeroWindowProbe struct {
      ID                    string    `json:"id"`
      SocketFd              int       `json:"socketFd"`
      ProbeCount            int       `json:"probeCount"`
      DurationMs            float64   `json:"durationMs"`
      IsProbeActive         bool      `json:"isProbeActive"`
      HasExponentialBackoff bool      `json:"hasExponentialBackoff"`
      ProbedAt              time.Time `json:"probedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/tcp/zero_window_prober.go` to handle socket backpressure.

---

### Improvement 1604: Asynchronous Multi-Line SMTP Response Demultiplexer
- **ID**: `IMP-1604`
- **Cluster**: 161
- **Goal**: Parse and demultiplex non-blocking, multi-line pipelined SMTP reply streams (e.g. `250-SIZE`, `250-PIPELINING`, `250 OK`) via ring buffer tokenizers without thread blocking.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status card formatting parsed reply codes in magenta (`#7C4DFF`) with tokenized reply parameters in muted slate (`#94A3B8`).
  - *Web*: `<smtp-response-demux>` custom Web Component with collapsible multi-line command/reply pairing accordion.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline parse-replies --session-id=<id> [--is-strict-delim] [--has-trailing-validation]`
  - Flag `--session-id`: Monitored session UUID.
  - Flag `--is-strict-delim`: Enforces strict CRLF byte sequence validation per RFC 5321.
  - Flag `--has-trailing-validation`: Checks for trailing bytes after 250 DATA completion.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpParsedReplies (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      replyCode INTEGER NOT NULL,
      replyLineCount INTEGER NOT NULL DEFAULT 1,
      isStrictDelim INTEGER NOT NULL DEFAULT 1,
      hasTrailingValidation INTEGER NOT NULL DEFAULT 1,
      parsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpParsedReply struct {
      ID                    string    `json:"id"`
      SessionID             string    `json:"sessionId"`
      ReplyCode             int       `json:"replyCode"`
      ReplyLineCount        int       `json:"replyLineCount"`
      IsStrictDelim         bool      `json:"isStrictDelim"`
      HasTrailingValidation bool      `json:"hasTrailingValidation"`
      ParsedAt              time.Time `json:"parsedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/smtp/reply_parser.go` for low-latency line parsing.

---

### Improvement 1605: Out-of-Order Pipelined RCPT Failure Isolation Boundary
- **ID**: `IMP-1605`
- **Cluster**: 161
- **Goal**: Dynamically isolate 550 or 450 recipient rejection reply codes within a pipelined command burst without terminating the entire envelope batch for valid recipients.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss split-view rendering individual recipient outcomes; valid recipients in spring emerald (`#00E676`) and rejected recipients in crimson (`#FF1744`).
  - *Web*: `<smtp-rcpt-matrix>` Web Component with interactive recipient grid highlighting deliverable versus rejected addresses.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline isolate-rcpt --session-id=<id> [--is-fail-fast] [--has-quarantine-routing]`
  - Flag `--session-id`: Target SMTP session UUID.
  - Flag `--is-fail-fast`: Aborts remaining pipeline batch immediately if any recipient fails.
  - Flag `--has-quarantine-routing`: Automatically re-routes rejected addresses to quarantine queue.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpRcptIsolations (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      acceptedCount INTEGER NOT NULL DEFAULT 0,
      rejectedCount INTEGER NOT NULL DEFAULT 0,
      isFailFast INTEGER NOT NULL DEFAULT 0,
      hasQuarantineRouting INTEGER NOT NULL DEFAULT 1,
      isolatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpRcptIsolation struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      AcceptedCount        int       `json:"acceptedCount"`
      RejectedCount        int       `json:"rejectedCount"`
      IsFailFast           bool      `json:"isFailFast"`
      HasQuarantineRouting bool      `json:"hasQuarantineRouting"`
      IsolatedAt           time.Time `json:"isolatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/smtp/rcpt_isolator.go`.

---

### Improvement 1606: Adaptive TCP SndBuf Dynamic Scaling during Mail Ingestion
- **ID**: `IMP-1606`
- **Cluster**: 161
- **Goal**: Automatically tune `SO_SNDBUF` socket sizes via kernel syscalls during transmission bursts to match remote advertised bandwidth-delay product (BDP).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss meter showing calculated BDP versus current socket send buffer allocation in electric cyan (`#00E5FF`).
  - *Web*: `<tcp-sndbuf-scaler>` custom element showing dynamic buffer curve charts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tcp-buf scale --socket-fd=<fd> [--target-bdp-bytes=<bytes>] [--is-auto-tuning] [--has-clamping-guard]`
  - Flag `--socket-fd`: Socket descriptor.
  - Flag `--target-bdp-bytes`: Target bandwidth-delay product headroom.
  - Flag `--is-auto-tuning`: Continuously adjusts buffer size on each RTT measurement.
  - Flag `--has-clamping-guard`: Restricts buffer from exceeding 4MB to prevent memory bloat.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpBufferScalings (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      calculatedBdpBytes INTEGER NOT NULL DEFAULT 0,
      currentSndBufBytes INTEGER NOT NULL DEFAULT 0,
      isAutoTuning INTEGER NOT NULL DEFAULT 1,
      hasClampingGuard INTEGER NOT NULL DEFAULT 1,
      scaledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpBufferScaling struct {
      ID                 string    `json:"id"`
      SocketFd           int       `json:"socketFd"`
      CalculatedBdpBytes int       `json:"calculatedBdpBytes"`
      CurrentSndBufBytes int       `json:"currentSndBufBytes"`
      IsAutoTuning       bool      `json:"isAutoTuning"`
      HasClampingGuard   bool      `json:"hasClampingGuard"`
      ScaledAt           time.Time `json:"scaledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tcp/buffer_tuner.go`.

---

### Improvement 1607: Zero-Copy Sendfile Payload Chunking with TCP Corking
- **ID**: `IMP-1607`
- **Cluster**: 161
- **Goal**: Transmit RFC 5322 MIME envelope bodies using Linux `sendfile(2)` / Windows `TransmitFile` with `TCP_CORK` socket options to prevent packet fragmentation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status bar showing transmission throughput with cork status icon in spring green (`#00E676`).
  - *Web*: `<sendfile-stream-meter>` Web Component displaying real-time zero-copy byte counters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup payload transmit-corked --envelope-id=<id> [--is-cork-enabled] [--has-direct-io]`
  - Flag `--envelope-id`: Message envelope UUID.
  - Flag `--is-cork-enabled`: Engages `TCP_CORK` / `TCP_NOPUSH` during MIME stream.
  - Flag `--has-direct-io`: Bypasses kernel page cache using direct I/O buffers.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpZeroCopyTransmissions (
      id TEXT PRIMARY KEY,
      envelopeId TEXT NOT NULL,
      totalBytesSent INTEGER NOT NULL DEFAULT 0,
      isCorkEnabled INTEGER NOT NULL DEFAULT 1,
      hasDirectIo INTEGER NOT NULL DEFAULT 0,
      transmittedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpZeroCopyTransmission struct {
      ID              string    `json:"id"`
      EnvelopeID      string    `json:"envelopeId"`
      TotalBytesSent  int64     `json:"totalBytesSent"`
      IsCorkEnabled   bool      `json:"isCorkEnabled"`
      HasDirectIo     bool      `json:"hasDirectIo"`
      TransmittedAt   time.Time `json:"transmittedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/smtp/zerocopy_sender.go`.

---

### Improvement 1608: Multi-Stage Pipelined DATA Command Stream Interleaving
- **ID**: `IMP-1608`
- **Cluster**: 161
- **Goal**: Pipeline the `DATA` command immediately following `RCPT TO` sequences and stream MIME chunks concurrently with server 354 Start Mail Input responses.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual ASCII pipeline stream showing interleaved commands with status color tags.
  - *Web*: `<data-interleave-view>` custom element visualizing interleaved command packets.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline interleave-data --session-id=<id> [--is-interleaving-active] [--has-pre-flight-354]`
  - Flag `--session-id`: SMTP session identifier.
  - Flag `--is-interleaving-active`: Enables opportunistic payload streaming.
  - Flag `--has-pre-flight-354`: Validates 354 response before releasing terminal dot.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpDataInterleaves (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      interleavedBytes INTEGER NOT NULL DEFAULT 0,
      isInterleavingActive INTEGER NOT NULL DEFAULT 1,
      hasPreFlight354 INTEGER NOT NULL DEFAULT 1,
      interleavedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpDataInterleave struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      InterleavedBytes     int64     `json:"interleavedBytes"`
      IsInterleavingActive bool      `json:"isInterleavingActive"`
      HasPreFlight354      bool      `json:"hasPreFlight354"`
      InterleavedAt        time.Time `json:"interleavedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/smtp/interleave.go`.

---

### Improvement 1609: Remote TCP Window Stall Early-Warning Forensics
- **ID**: `IMP-1609`
- **Cluster**: 161
- **Goal**: Capture microsecond socket stall events where remote receive window decreases below 3x MSS, writing forensic diagnostic records to `repodb/audit.db`.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Warning badge in amber flare (`#FFD600`) displaying socket stall duration and remote peer IP.
  - *Web*: `<window-stall-alert>` Web Component presenting historical stall incident cards.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup forensics audit-stall --socket-fd=<fd> [--is-stall-persisted] [--has-stack-trace]`
  - Flag `--socket-fd`: Monitored socket descriptor.
  - Flag `--is-stall-persisted`: Commits diagnostic snapshot to `repodb/audit.db`.
  - Flag `--has-stack-trace`: Captures current execution goroutine stack trace.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowForensicProofs (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      remotePeerIp TEXT NOT NULL,
      windowBytes INTEGER NOT NULL DEFAULT 0,
      stallDurationMs REAL NOT NULL DEFAULT 0.0,
      isStallPersisted INTEGER NOT NULL DEFAULT 1,
      hasStackTrace INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowForensicProof struct {
      ID               string    `json:"id"`
      SocketFd         int       `json:"socketFd"`
      RemotePeerIP     string    `json:"remotePeerIp"`
      WindowBytes      int       `json:"windowBytes"`
      StallDurationMs  float64   `json:"stallDurationMs"`
      IsStallPersisted bool      `json:"isStallPersisted"`
      HasStackTrace    bool      `json:"hasStackTrace"`
      RecordedAt       time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/audit/stall_forensics.go`.

---

### Improvement 1610: TCP Zero-Window Recovery Handshake & Auto-Resume Trigger
- **ID**: `IMP-1610`
- **Cluster**: 161
- **Goal**: Seamlessly detect window reopening ACK notifications (`snd_wnd > MSS`), triggering instant resumption of pending pipelined batches without connection teardown.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status flash in spring emerald (`#00E676`) indicating zero-window recovery with latency metrics.
  - *Web*: `<window-recovery-badge>` Web Component rendering real-time resume animations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tcp-window resume --socket-fd=<fd> [--is-resume-immediate] [--has-rate-taper]`
  - Flag `--socket-fd`: Socket descriptor.
  - Flag `--is-resume-immediate`: Immediately fires queued command frames.
  - Flag `--has-rate-taper`: Tapers transmission rate for first 5 envelopes post-recovery.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpWindowRecoveries (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      resumedWindowBytes INTEGER NOT NULL DEFAULT 0,
      recoveryLatencyMs REAL NOT NULL DEFAULT 0.0,
      isResumeImmediate INTEGER NOT NULL DEFAULT 1,
      hasRateTaper INTEGER NOT NULL DEFAULT 1,
      resumedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpWindowRecovery struct {
      ID                 string    `json:"id"`
      SocketFd           int       `json:"socketFd"`
      ResumedWindowBytes int       `json:"resumedWindowBytes"`
      RecoveryLatencyMs  float64   `json:"recoveryLatencyMs"`
      IsResumeImmediate  bool      `json:"isResumeImmediate"`
      HasRateTaper       bool      `json:"hasRateTaper"`
      ResumedAt          time.Time `json:"resumedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/tcp/recovery_trigger.go`.

---

# Cluster 162: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (1611–1620)

---

### Improvement 1611: Multi-ISP Persistent Socket Pool Registry & Lifecycle Balancer
- **ID**: `IMP-1611`
- **Cluster**: 162 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal**: Manage persistent, pre-warmed keep-alive socket pools partitioned by destination MX provider (Google, Microsoft, Yahoo) to eliminate repeated TCP/TLS handshake latency.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss grid table displaying active, idle, and max socket slots across destination providers in electric cyan (`#00E5FF`).
  - *Web*: `<isp-pool-balancer>` Web Component displaying socket pool capacity bars.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool list [--isp-filter=<isp>] [--is-active-only] [--has-idle-headroom]`
  - Flag `--isp-filter`: Filters by target ISP (e.g. `google`, `microsoft`).
  - Flag `--is-active-only`: Limits output to currently transmitting pools.
  - Flag `--has-idle-headroom`: Includes pool warm-standby headroom counts.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspConnectionPoolNodes (
      id TEXT PRIMARY KEY,
      ispProvider TEXT NOT NULL,
      activeSockets INTEGER NOT NULL DEFAULT 0,
      idleSockets INTEGER NOT NULL DEFAULT 0,
      maxPoolCapacity INTEGER NOT NULL DEFAULT 100,
      isActiveOnly INTEGER NOT NULL DEFAULT 0,
      hasIdleHeadroom INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspConnectionPoolNode struct {
      ID              string    `json:"id"`
      IspProvider     string    `json:"ispProvider"`
      ActiveSockets   int       `json:"activeSockets"`
      IdleSockets     int       `json:"idleSockets"`
      MaxPoolCapacity int       `json:"maxPoolCapacity"`
      IsActiveOnly    bool      `json:"isActiveOnly"`
      HasIdleHeadroom bool      `json:"hasIdleHeadroom"`
      UpdatedAt       time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Centralized in `pkg/warmup/pool/registry.go`.

---

### Improvement 1612: TCP Handshake & TLS Negotiation Latency Telemetry Matrix
- **ID**: `IMP-1612`
- **Cluster**: 162
- **Goal**: Record sub-millisecond nanosecond-precision metrics for TCP SYN/ACK roundtrips, TLS 1.3 ClientHello/ServerHello timing, and cipher negotiation per destination gateway.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss bar chart with handshake breakdown in electric violet (`#7C4DFF`) and TLS completion in cyan (`#00E5FF`).
  - *Web*: `<handshake-latency-matrix>` Web Component with interactive SVG latency heatmaps.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool measure-handshake --target-mx=<mx> [--is-tls13-required] [--has-alpn-probe]`
  - Flag `--target-mx`: Destination MX domain or IP.
  - Flag `--is-tls13-required`: Requires TLS 1.3 protocol verification.
  - Flag `--has-alpn-probe`: Tests ALPN negotiation support.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketHealthMetrics (
      id TEXT PRIMARY KEY,
      targetMx TEXT NOT NULL,
      tcpHandshakeMs REAL NOT NULL DEFAULT 0.0,
      tlsHandshakeMs REAL NOT NULL DEFAULT 0.0,
      isTls13Required INTEGER NOT NULL DEFAULT 1,
      hasAlpnProbe INTEGER NOT NULL DEFAULT 0,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketHealthMetric struct {
      ID               string    `json:"id"`
      TargetMx         string    `json:"targetMx"`
      TcpHandshakeMs   float64   `json:"tcpHandshakeMs"`
      TlsHandshakeMs   float64   `json:"tlsHandshakeMs"`
      IsTls13Required  bool      `json:"isTls13Required"`
      HasAlpnProbe     bool      `json:"hasAlpnProbe"`
      MeasuredAt       time.Time `json:"measuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/health_matrix.go`.

---

### Improvement 1613: Proactive Keep-Alive Heartbeat & NOOP Socket Rejuvenator
- **ID**: `IMP-1613`
- **Cluster**: 162
- **Goal**: Maintain idle socket warm state by injecting opportunistic `NOOP` or `RSET` command bursts at dynamically jittered intervals below provider inactivity timeouts.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Mini status indicator showing keep-alive heartbeat cadence in spring emerald (`#00E676`).
  - *Web*: `<socket-heartbeat-ticker>` Web Component with ticking pulse icons.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool heartbeat --pool-id=<id> [--interval-sec=<sec>] [--is-jitter-randomized] [--has-silent-discard]`
  - Flag `--pool-id`: Target connection pool UUID.
  - Flag `--interval-sec`: Base heartbeat interval in seconds (default: `45`).
  - Flag `--is-jitter-randomized`: Applies ±15% Gaussian jitter to prevent synchronized burst spikes.
  - Flag `--has-silent-discard`: Discards heartbeat 250 response without emitting event logs.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolHeartbeatSchedules (
      id TEXT PRIMARY KEY,
      poolId TEXT NOT NULL,
      intervalSeconds INTEGER NOT NULL DEFAULT 45,
      isJitterRandomized INTEGER NOT NULL DEFAULT 1,
      hasSilentDiscard INTEGER NOT NULL DEFAULT 1,
      scheduledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolHeartbeatSchedule struct {
      ID                 string    `json:"id"`
      PoolID             string    `json:"poolId"`
      IntervalSeconds    int       `json:"intervalSeconds"`
      IsJitterRandomized bool      `json:"isJitterRandomized"`
      HasSilentDiscard   bool      `json:"hasSilentDiscard"`
      ScheduledAt        time.Time `json:"scheduledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/pool/heartbeat.go`.

---

### Improvement 1614: Composite Socket Degradation Scoring & Eviction Governor
- **ID**: `IMP-1614`
- **Cluster**: 162
- **Goal**: Compute a 0–100 real-time socket health score combining TCP RTT jitter, TLS session re-use rate, and banner latency, automatically evicting degraded sockets before sending mail.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss colored badge: green for scores > 80, yellow for 50–80, red for < 50.
  - *Web*: `<socket-score-gauge>` Web Component displaying score arc with eviction threshold marker.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool score --socket-id=<id> [--eviction-threshold=<score>] [--is-auto-evict] [--has-tarpit-penalty]`
  - Flag `--socket-id`: Monitored socket UUID.
  - Flag `--eviction-threshold`: Score cutoff for immediate eviction (default: `40`).
  - Flag `--is-auto-evict`: Closes socket immediately if score falls below threshold.
  - Flag `--has-tarpit-penalty`: Deducts 30 points if remote peer delays response > 2500ms.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketDegradationScores (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      compositeScore INTEGER NOT NULL DEFAULT 100,
      isAutoEvict INTEGER NOT NULL DEFAULT 1,
      hasTarpitPenalty INTEGER NOT NULL DEFAULT 0,
      scoredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketDegradationScore struct {
      ID               string    `json:"id"`
      SocketID         string    `json:"socketId"`
      CompositeScore   int       `json:"compositeScore"`
      IsAutoEvict      bool      `json:"isAutoEvict"`
      HasTarpitPenalty bool      `json:"hasTarpitPenalty"`
      ScoredAt         time.Time `json:"scoredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/pool/eviction_governor.go`.

---

### Improvement 1615: Multi-Egress IP Affinity & Connection Pinning Router
- **ID**: `IMP-1615`
- **Cluster**: 162
- **Goal**: Pin warming domains to dedicated local egress IPs while balancing socket load across available interface aliases to maximize reputation accumulation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss mapping table showing sender domain, pinned local IP, target MX provider, and active connection count.
  - *Web*: `<egress-ip-affinity-matrix>` Web Component with route topology diagrams.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool route-pin --domain=<domain> --egress-ip=<ip> [--is-sticky] [--has-failover-alias]`
  - Flag `--domain`: Sending domain identity.
  - Flag `--egress-ip`: Dedicated local interface IP.
  - Flag `--is-sticky`: Enforces strict session pinning to the designated egress IP.
  - Flag `--has-failover-alias`: Allows fallback to secondary IP alias if primary interface stalls.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EgressIpAffinities (
      id TEXT PRIMARY KEY,
      domain TEXT NOT NULL,
      egressIp TEXT NOT NULL,
      isSticky INTEGER NOT NULL DEFAULT 1,
      hasFailoverAlias INTEGER NOT NULL DEFAULT 1,
      pinnedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EgressIpAffinity struct {
      ID               string    `json:"id"`
      Domain           string    `json:"domain"`
      EgressIP         string    `json:"egressIp"`
      IsSticky         bool      `json:"isSticky"`
      HasFailoverAlias bool      `json:"hasFailoverAlias"`
      PinnedAt         time.Time `json:"pinnedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/pool/affinity_router.go`.

---

### Improvement 1616: Socket Pool State Transition Audit Ledger
- **ID**: `IMP-1616`
- **Cluster**: 162
- **Goal**: Record all socket state changes (`ALLOCATED`, `CONNECTED`, `IDLE`, `DRAINING`, `EVICTED`, `TERMINATED`) in `repodb/audit.db` for full compliance traceability.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss audit log stream with state transition tags in electric cyan (`#00E5FF`) and timestamps in muted slate (`#94A3B8`).
  - *Web*: `<pool-audit-stream>` Web Component presenting an append-only timeline table.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup audit pool-transitions [--socket-id=<id>] [--limit=<int>] [--is-chronological] [--has-hash-verify]`
  - Flag `--socket-id`: Filters transitions by socket UUID.
  - Flag `--limit`: Number of records returned (default: `50`).
  - Flag `--is-chronological`: Sorts ascending by timestamp.
  - Flag `--has-hash-verify`: Validates cryptographic chain hashes for returned rows.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolStateTransitionAudits (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      previousState TEXT NOT NULL,
      newState TEXT NOT NULL,
      reasonCode TEXT NOT NULL,
      isChronological INTEGER NOT NULL DEFAULT 1,
      hasHashVerify INTEGER NOT NULL DEFAULT 1,
      transitionedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolStateTransitionAudit struct {
      ID              string    `json:"id"`
      SocketID        string    `json:"socketId"`
      PreviousState   string    `json:"previousState"`
      NewState        string    `json:"newState"`
      ReasonCode      string    `json:"reasonCode"`
      IsChronological bool      `json:"isChronological"`
      HasHashVerify   bool      `json:"hasHashVerify"`
      TransitionedAt  time.Time `json:"transitionedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/audit/pool_transition_logger.go`.

---

### Improvement 1617: Anti-Tarpit Rapid Disconnect & Quarantine Switch
- **ID**: `IMP-1617`
- **Cluster**: 162
- **Goal**: Detect defensive ISP connection tarpitting (banner delay > 5000ms or byte-drip rate < 10 bytes/sec), immediately issuing `RST` teardown and placing the IP route in cooling quarantine.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-contrast red warning box (`#FF1744`) alerting tarpit detection with quarantine cooldown countdown.
  - *Web*: `<tarpit-quarantine-card>` Web Component displaying quarantined destination MX and remaining cooling period.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool tarpit-detect --socket-id=<id> [--timeout-ms=<ms>] [--is-immediate-rst] [--has-route-quarantine]`
  - Flag `--socket-id`: Monitored socket UUID.
  - Flag `--timeout-ms`: Maximum banner read delay before tarpit trigger (default: `4000`).
  - Flag `--is-immediate-rst`: Sends TCP RST instead of graceful FIN to terminate tarpit instantly.
  - Flag `--has-route-quarantine`: Puts route into 30-minute quarantine.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TarpitQuarantineRecords (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      targetMx TEXT NOT NULL,
      quarantineDurationMin INTEGER NOT NULL DEFAULT 30,
      isImmediateRst INTEGER NOT NULL DEFAULT 1,
      hasRouteQuarantine INTEGER NOT NULL DEFAULT 1,
      quarantinedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TarpitQuarantineRecord struct {
      ID                    string    `json:"id"`
      SocketID              string    `json:"socketId"`
      TargetMx              string    `json:"targetMx"`
      QuarantineDurationMin int       `json:"quarantineDurationMin"`
      IsImmediateRst        bool      `json:"isImmediateRst"`
      HasRouteQuarantine    bool      `json:"hasRouteQuarantine"`
      QuarantinedAt         time.Time `json:"quarantinedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/pool/tarpit_detector.go`.

---

### Improvement 1618: Dynamic Sched_Yield Worker Thread Scheduling for Socket Pools
- **ID**: `IMP-1618`
- **Cluster**: 162
- **Goal**: Optimize Go runtime scheduler execution across socket pool worker goroutines via selective `runtime.Gosched()` and OS thread affinity pinning to minimize context switch jitter.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss CPU core map displaying thread allocation and goroutine occupancy.
  - *Web*: `<scheduler-affinity-view>` Web Component rendering CPU core utilization heatmaps.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool tune-runtime [--worker-count=<int>] [--is-thread-pinned] [--has-yield-optimization]`
  - Flag `--worker-count`: Number of dedicated worker goroutines (default: `8`).
  - Flag `--is-thread-pinned`: Locks worker threads to OS CPU cores using runtime.LockOSThread.
  - Flag `--has-yield-optimization`: Yields processor slice when socket wait queue is empty.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WorkerRuntimeConfigs (
      id TEXT PRIMARY KEY,
      workerCount INTEGER NOT NULL DEFAULT 8,
      isThreadPinned INTEGER NOT NULL DEFAULT 0,
      hasYieldOptimization INTEGER NOT NULL DEFAULT 1,
      configuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WorkerRuntimeConfig struct {
      ID                   string    `json:"id"`
      WorkerCount          int       `json:"workerCount"`
      IsThreadPinned       bool      `json:"isThreadPinned"`
      HasYieldOptimization bool      `json:"hasYieldOptimization"`
      ConfiguredAt         time.Time `json:"configuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/pool/worker_runtime.go`.

---

### Improvement 1619: Transparent TLS Session Resumption & 0-RTT PSK Cache
- **ID**: `IMP-1619`
- **Cluster**: 162
- **Goal**: Cache TLS 1.3 Pre-Shared Key (PSK) session tickets across SMTP reconnection attempts to eliminate TLS handshake round trips with remote MX hosts.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss session cache card showing PSK hit rate percentage in electric cyan (`#00E5FF`).
  - *Web*: `<tls-psk-cache-stats>` Web Component displaying ticket expiration countdowns and hit/miss ratios.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool psk-cache inspect [--target-mx=<mx>] [--is-psk-valid] [--has-ticket-refresh]`
  - Flag `--target-mx`: Target MX domain.
  - Flag `--is-psk-valid`: Validates cached ticket signature without establishing socket.
  - Flag `--has-ticket-refresh`: Requests updated ticket on next handshake.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TlsPskCacheEntries (
      id TEXT PRIMARY KEY,
      targetMx TEXT NOT NULL,
      ticketLifetimeSec INTEGER NOT NULL DEFAULT 7200,
      isPskValid INTEGER NOT NULL DEFAULT 1,
      hasTicketRefresh INTEGER NOT NULL DEFAULT 1,
      cachedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TlsPskCacheEntry struct {
      ID                string    `json:"id"`
      TargetMx          string    `json:"targetMx"`
      TicketLifetimeSec int       `json:"ticketLifetimeSec"`
      IsPskValid        bool      `json:"isPskValid"`
      HasTicketRefresh  bool      `json:"hasTicketRefresh"`
      CachedAt          time.Time `json:"cachedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/pool/tls_cache.go`.

---

### Improvement 1620: Cross-ISP Concurrent MX Connection Race Resolver (Happy Eyeballs SMTP)
- **ID**: `IMP-1620`
- **Cluster**: 162
- **Goal**: Race IPv4 and IPv6 MX target host connections concurrently with a 250ms offset per RFC 8305 to select the lowest-latency, healthiest transport route automatically.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss split-route comparison view showing IPv4 vs IPv6 race outcome and win margins in green (`#00E676`).
  - *Web*: `<happy-eyeballs-race-view>` Web Component rendering side-by-side connection race animations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool race-mx --target-domain=<domain> [--offset-ms=<ms>] [--is-ipv6-preferred] [--has-winner-pinning]`
  - Flag `--target-domain`: Recipient domain identity.
  - Flag `--offset-ms`: Delay before launching fallback family connection (default: `250`).
  - Flag `--is-ipv6-preferred`: Prioritizes IPv6 probe launch.
  - Flag `--has-winner-pinning`: Pins winner IP family for subsequent pool connections for 60 minutes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HappyEyeballsRaces (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      winnerIpFamily TEXT NOT NULL,
      winMarginMs REAL NOT NULL DEFAULT 0.0,
      isIpv6Preferred INTEGER NOT NULL DEFAULT 1,
      hasWinnerPinning INTEGER NOT NULL DEFAULT 1,
      racedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HappyEyeballsRace struct {
      ID               string    `json:"id"`
      TargetDomain     string    `json:"targetDomain"`
      WinnerIPFamily   string    `json:"winnerIpFamily"`
      WinMarginMs      float64   `json:"winMarginMs"`
      IsIpv6Preferred  bool      `json:"isIpv6Preferred"`
      HasWinnerPinning bool      `json:"hasWinnerPinning"`
      RacedAt          time.Time `json:"racedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/pool/happy_eyeballs.go`.

---

# Cluster 163: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1621–1630)

---

### Improvement 1621: Sub-Millisecond Socket Lifecycle Flamegraph Renderer
- **ID**: `IMP-1621`
- **Cluster**: 163 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Render hierarchical execution flamegraphs of SMTP session lifecycles (DNS -> SYN -> TLS -> EHLO -> MAIL -> RCPT -> DATA) in terminal consoles via Lipgloss and in browsers via zero-CDN SVG.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss nested horizontal bar blocks styled with TrueColor gradients from navy (`#0B0F19`) to cyan (`#00E5FF`) and violet (`#7C4DFF`).
  - *Web*: `<smtp-flamegraph-renderer>` Web Component with interactive SVG tooltips showing exact microsecond phase breakdowns.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui flamegraph --session-id=<id> [--format=<ansi|svg>] [--is-sub-ms-scale] [--has-phase-drilldown]`
  - Flag `--session-id`: Target session UUID.
  - Flag `--format`: Output rendering format (`ansi` or `svg`, default: `ansi`).
  - Flag `--is-sub-ms-scale`: Formats timeline axis with microsecond precision.
  - Flag `--has-phase-drilldown`: Outputs child call breakdown for TLS and DATA phases.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS FlamegraphExecutionPhases (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      phaseName TEXT NOT NULL,
      durationMicroseconds INTEGER NOT NULL DEFAULT 0,
      isSubMsScale INTEGER NOT NULL DEFAULT 1,
      hasPhaseDrilldown INTEGER NOT NULL DEFAULT 1,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type FlamegraphExecutionPhase struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      PhaseName            string    `json:"phaseName"`
      DurationMicroseconds int64     `json:"durationMicroseconds"`
      IsSubMsScale         bool      `json:"isSubMsScale"`
      HasPhaseDrilldown    bool      `json:"hasPhaseDrilldown"`
      CapturedAt           time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/flamegraph.go`.

---

### Improvement 1622: ASCII Jitter Waveform Sparkline Generator for Bubble Tea v2
- **ID**: `IMP-1622`
- **Cluster**: 163
- **Goal**: Generate high-density ASCII jitter sparklines (` ▂▃▄▅▆▇█`) in Bubble Tea v2 HUD viewports showing microsecond ACK timing variations across the last 128 packets.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss sparkline widget with baseline in muted slate (`#94A3B8`) and peak jitter spikes in magenta flare (`#E040FB`).
  - *Web*: `<ascii-jitter-sparkline>` Web Component displaying animated SVG waveform bars.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui sparkline --stream-id=<id> [--sample-count=<count>] [--is-compact-hud] [--has-threshold-marker]`
  - Flag `--stream-id`: Monitored stream identifier.
  - Flag `--sample-count`: Number of sliding samples rendered (default: `128`).
  - Flag `--is-compact-hud`: Emits single-line compact representation suitable for status bars.
  - Flag `--has-threshold-marker`: Highlights values exceeding 2-sigma variance in red.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterWaveformTraces (
      id TEXT PRIMARY KEY,
      streamId TEXT NOT NULL,
      sampleIndex INTEGER NOT NULL DEFAULT 0,
      jitterMicroseconds INTEGER NOT NULL DEFAULT 0,
      isCompactHud INTEGER NOT NULL DEFAULT 0,
      hasThresholdMarker INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterWaveformTrace struct {
      ID                 string    `json:"id"`
      StreamID           string    `json:"streamId"`
      SampleIndex        int       `json:"sampleIndex"`
      JitterMicroseconds int       `json:"jitterMicroseconds"`
      IsCompactHud       bool      `json:"isCompactHud"`
      HasThresholdMarker bool      `json:"hasThresholdMarker"`
      RecordedAt         time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/sparkline.go`.

---

### Improvement 1623: Zero-CDN Pure SVG Vector Export for Socket Latency Traces
- **ID**: `IMP-1623`
- **Cluster**: 163
- **Goal**: Export standalone vector SVG files representing full socket handshake and transmission latency traces with zero external stylesheet or JavaScript dependencies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Confirmation card with SVG file path and size metrics formatted in emerald (`#00E676`).
  - *Web*: Inlines SVG vector element directly into host DOM without CDN scripts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui export-svg --session-id=<id> --out=<path> [--is-dark-theme] [--has-embedded-styles]`
  - Flag `--session-id`: Target session UUID.
  - Flag `--out`: Relative output path for generated SVG file.
  - Flag `--is-dark-theme`: Applies dark mode styling (`#0B0F19` background).
  - Flag `--has-embedded-styles`: Inlines all CSS definitions inside `<style>` tags.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SvgExportRecords (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      outputPath TEXT NOT NULL,
      fileSizeBytes INTEGER NOT NULL DEFAULT 0,
      isDarkTheme INTEGER NOT NULL DEFAULT 1,
      hasEmbeddedStyles INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SvgExportRecord struct {
      ID                string    `json:"id"`
      SessionID         string    `json:"sessionId"`
      OutputPath        string    `json:"outputPath"`
      FileSizeBytes     int64     `json:"fileSizeBytes"`
      IsDarkTheme       bool      `json:"isDarkTheme"`
      HasEmbeddedStyles bool      `json:"hasEmbeddedStyles"`
      ExportedAt        time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/tui/svg_exporter.go`.

---

### Improvement 1624: Interactive Bubble Tea v2 HUD Viewport & Key Navigation
- **ID**: `IMP-1624`
- **Cluster**: 163
- **Goal**: Provide an interactive keyboard-navigable TUI viewport (Arrow keys, `Tab`, `q`, `Enter`) allowing operators to freeze, inspect, and drill down into active warming session buffers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Full-screen Bubble Tea v2 application with header breadcrumbs, active viewport pane, and bottom hotkey cheatsheet in Lipgloss Deep Slate (`#161F30`).
  - *Web*: Mirrors terminal keybind shortcuts inside standalone browser UI.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui live-hud [--fps=<fps>] [--is-fullscreen] [--has-mouse-support]`
  - Flag `--fps`: Viewport refresh rate (default: `30`).
  - Flag `--is-fullscreen`: Engages alternate terminal buffer mode.
  - Flag `--has-mouse-support`: Enables mouse click selection on session items.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS LiveHudSessionStates (
      id TEXT PRIMARY KEY,
      currentFps INTEGER NOT NULL DEFAULT 30,
      isFullscreen INTEGER NOT NULL DEFAULT 1,
      hasMouseSupport INTEGER NOT NULL DEFAULT 0,
      startedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type LiveHudSessionState struct {
      ID              string    `json:"id"`
      CurrentFps      int       `json:"currentFps"`
      IsFullscreen    bool      `json:"isFullscreen"`
      HasMouseSupport bool      `json:"hasMouseSupport"`
      StartedAt       time.Time `json:"startedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/tui/live_hud.go`.

---

### Improvement 1625: Real-Time Jitter Deviation Threshold Alerting
- **ID**: `IMP-1625`
- **Cluster**: 163
- **Goal**: Emit real-time audible/visual alerts in TUI when TCP acknowledgment timing jitter deviates >3 standard deviations from running baseline, persisting events to `repodb/audit.db`.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Flashing terminal bell and vibrating red highlight border (`#FF1744`) on affected session card.
  - *Web*: `<jitter-alert-banner>` Web Component displaying deviation multiplier badge.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui watch-jitter --session-id=<id> [--sigma-threshold=<float>] [--is-bell-enabled] [--has-audit-persistence]`
  - Flag `--session-id`: Target session UUID.
  - Flag `--sigma-threshold`: Standard deviation alert threshold (default: `3.0`).
  - Flag `--is-bell-enabled`: Triggers terminal ASCII bell `\a` on alert trigger.
  - Flag `--has-audit-persistence`: Writes event snapshot to `repodb/audit.db`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterDeviationAlertLogs (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      deviationSigma REAL NOT NULL DEFAULT 0.0,
      observedJitterMs REAL NOT NULL DEFAULT 0.0,
      isBellEnabled INTEGER NOT NULL DEFAULT 0,
      hasAuditPersistence INTEGER NOT NULL DEFAULT 1,
      alertedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterDeviationAlertLog struct {
      ID                  string    `json:"id"`
      SessionID           string    `json:"sessionId"`
      DeviationSigma      float64   `json:"deviationSigma"`
      ObservedJitterMs    float64   `json:"observedJitterMs"`
      IsBellEnabled       bool      `json:"isBellEnabled"`
      HasAuditPersistence bool      `json:"hasAuditPersistence"`
      AlertedAt           time.Time `json:"alertedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/jitter_watcher.go`.

---

### Improvement 1626: Shadow DOM Encapsulated HTML5 Live Telemetry Exporter
- **ID**: `IMP-1626`
- **Cluster**: 163
- **Goal**: Export self-contained single-file HTML reports bundling CSS tokens, SVG components, and embedded telemetry data within Shadow DOM scopes without network calls.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Green summary card reporting total exported HTML byte size and file hash.
  - *Web*: Completely self-contained `<gitmap-telemetry-root>` element isolating styles from embedding parent pages.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui export-html --report-id=<id> --out=<path> [--is-standalone] [--has-shadow-dom]`
  - Flag `--report-id`: Source telemetry report UUID.
  - Flag `--out`: Relative output path for generated HTML file.
  - Flag `--is-standalone`: Embeds all fonts and SVG vectors inline.
  - Flag `--has-shadow-dom`: Wraps report interface inside Shadow Root.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HtmlTelemetryExports (
      id TEXT PRIMARY KEY,
      reportId TEXT NOT NULL,
      outputPath TEXT NOT NULL,
      isStandalone INTEGER NOT NULL DEFAULT 1,
      hasShadowDom INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HtmlTelemetryExport struct {
      ID           string    `json:"id"`
      ReportID     string    `json:"reportId"`
      OutputPath   string    `json:"outputPath"`
      IsStandalone bool      `json:"isStandalone"`
      HasShadowDom bool      `json:"hasShadowDom"`
      ExportedAt   time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/html_exporter.go`.

---

### Improvement 1627: Multi-Session Concurrent Latency Comparison Splitter
- **ID**: `IMP-1627`
- **Cluster**: 163
- **Goal**: Render multi-pane comparative views in terminal screens contrasting real-time latency curves across simultaneous warming streams to different destination providers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss split-screen layout displaying up to 4 parallel session graphs side-by-side.
  - *Web*: `<multi-session-splitter>` Web Component with synchronized scrub timelines.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui compare-sessions --session-ids=<ids> [--columns=<int>] [--is-synchronized-scrub] [--has-delta-display]`
  - Flag `--session-ids`: Comma-separated list of session UUIDs.
  - Flag `--columns`: Display columns in grid (default: `2`).
  - Flag `--is-synchronized-scrub`: Locks scrub timeline across all active panes.
  - Flag `--has-delta-display`: Shows percentage latency variance relative to reference session.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SessionComparisonViews (
      id TEXT PRIMARY KEY,
      sessionIds TEXT NOT NULL,
      displayColumns INTEGER NOT NULL DEFAULT 2,
      isSynchronizedScrub INTEGER NOT NULL DEFAULT 1,
      hasDeltaDisplay INTEGER NOT NULL DEFAULT 1,
      createdDat DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SessionComparisonView struct {
      ID                  string    `json:"id"`
      SessionIDs          string    `json:"sessionIds"`
      DisplayColumns      int       `json:"displayColumns"`
      IsSynchronizedScrub bool      `json:"isSynchronizedScrub"`
      HasDeltaDisplay     bool      `json:"hasDeltaDisplay"`
      CreatedDat          time.Time `json:"createdDat"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/tui/session_comparator.go`.

---

### Improvement 1628: Nanosecond Phase Tick Precision Counter
- **ID**: `IMP-1628`
- **Cluster**: 163
- **Goal**: Provide hardware monotonic clock sampling (`runtime_nanotime`) across every socket event, storing high-precision event deltas in memory buffers for flamegraph rendering.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Nanosecond duration chips formatted with engineering unit notation (`ns`, `μs`, `ms`).
  - *Web*: `<nano-tick-counter>` Web Component presenting exact timestamp counters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui nano-timer --event-id=<id> [--is-monotonic] [--has-hardware-tsc]`
  - Flag `--event-id`: Identifier of lifecycle event.
  - Flag `--is-monotonic`: Enforces monotonic clock delta calculations.
  - Flag `--has-hardware-tsc`: Reads CPU Time Stamp Counter (TSC) when available.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS NanosecondPhaseTicks (
      id TEXT PRIMARY KEY,
      eventId TEXT NOT NULL,
      durationNanos INTEGER NOT NULL DEFAULT 0,
      isMonotonic INTEGER NOT NULL DEFAULT 1,
      hasHardwareTsc INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type NanosecondPhaseTick struct {
      ID             string    `json:"id"`
      EventID        string    `json:"eventId"`
      DurationNanos  int64     `json:"durationNanos"`
      IsMonotonic    bool      `json:"isMonotonic"`
      HasHardwareTsc bool      `json:"hasHardwareTsc"`
      RecordedAt     time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/tui/nanotimer.go`.

---

### Improvement 1629: TrueColor Lipgloss Theme Swapper for Terminal Accessibility
- **ID**: `IMP-1629`
- **Cluster**: 163
- **Goal**: Support high-contrast, protanopia, deuteranopia, and monochrome palette presets for terminal TUI views adhering strictly to WCAG 2.1 AAA contrast guidelines.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Real-time palette switcher toggling terminal colors without clearing screen state.
  - *Web*: Dynamically swaps CSS custom properties across all Web Components.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui set-theme --theme-name=<name> [--is-high-contrast] [--has-accessible-ansi]`
  - Flag `--theme-name`: Palette preset (`default`, `protanopia`, `deuteranopia`, `monochrome`).
  - Flag `--is-high-contrast`: Boosts contrast ratios above 7:1.
  - Flag `--has-accessible-ansi`: Restricts color output to 16 ANSI color standards.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiThemePreferences (
      id TEXT PRIMARY KEY,
      themeName TEXT NOT NULL DEFAULT 'default',
      isHighContrast INTEGER NOT NULL DEFAULT 0,
      hasAccessibleAnsi INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiThemePreference struct {
      ID                string    `json:"id"`
      ThemeName         string    `json:"themeName"`
      IsHighContrast    bool      `json:"isHighContrast"`
      HasAccessibleAnsi bool      `json:"hasAccessibleAnsi"`
      UpdatedAt         time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/tui/theme_swapper.go`.

---

### Improvement 1630: Headless TUI Snapshot Exporter for CI/CD Pipeline Verification
- **ID**: `IMP-1630`
- **Cluster**: 163
- **Goal**: Capture golden ANSI terminal snapshots of warming TUI screens in headless CI/CD runner environments, validating layout integrity without interactive tty devices.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Headless test runner output reporting pixel and ANSI escape code parity with reference golden snapshots.
  - *Web*: `<headless-snapshot-diff>` Web Component displaying visual before/after diff overlays.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui snapshot --view=<name> --out=<path> [--is-golden-compare] [--has-strip-ansi]`
  - Flag `--view`: Target TUI view to render (`flamegraph`, `sparkline`, `hud`).
  - Flag `--out`: Output path for generated snapshot file.
  - Flag `--is-golden-compare`: Compares output against checked-in golden file.
  - Flag `--has-strip-ansi`: Strips escape codes to produce plain text representation.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiGoldenSnapshots (
      id TEXT PRIMARY KEY,
      viewName TEXT NOT NULL,
      snapshotPath TEXT NOT NULL,
      isGoldenCompare INTEGER NOT NULL DEFAULT 1,
      hasStripAnsi INTEGER NOT NULL DEFAULT 0,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiGoldenSnapshot struct {
      ID              string    `json:"id"`
      ViewName        string    `json:"viewName"`
      SnapshotPath    string    `json:"snapshotPath"`
      IsGoldenCompare bool      `json:"isGoldenCompare"`
      HasStripAnsi    bool      `json:"hasStripAnsi"`
      CapturedAt      time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/snapshot_tester.go`.
