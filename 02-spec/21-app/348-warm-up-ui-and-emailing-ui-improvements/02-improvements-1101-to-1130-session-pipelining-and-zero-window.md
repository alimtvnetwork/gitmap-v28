# Improvements 1101 to 1130: Session Pipelining, Socket Health & Jitter Waveform TUI

This document details the complete specifications for Improvements 1101 through 1130 across Clusters 111, 112, and 113 of Sequence 348.

---

# Cluster 111: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (1101–1110)

---

### Improvement 1101: RFC 2920 Command Streaming Pipeline & Zero-RTT Command Burst Engine
- **ID**: `IMP-1101`
- **Cluster**: 111 (Autonomous SMTP Session Pipelining & TCP Zero-Window Detection)
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

### Improvement 1102: Real-Time Zero-Copy TCP Receive Window Monitor via Socket Reflection
- **ID**: `IMP-1102`
- **Cluster**: 111
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

### Improvement 1103: Autonomous TCP Zero-Window Detection & Preemptive Backoff Controller
- **ID**: `IMP-1103`
- **Cluster**: 111
- **Goal**: Detect when remote ISP MTA advertises a zero-length receive window (`snd_wnd == 0`), automatically halting pipelined envelope transmission and engaging microsecond backoff timers without dropping the TCP connection.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss banner rendering a pulsing magenta alert (`#E040FB`) signaling zero-window state, with an ASCII countdown timer for next drain probe.
  - *Web*: `<zero-window-banner>` Web Component providing interactive backoff status with zero CDN dependencies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline zero-window-backoff --socket-fd=<fd> [--timeout-ms=<ms>] [--is-auto-drain] [--has-exponential-scale]`
  - Flag `--socket-fd`: The saturated socket file descriptor.
  - Flag `--timeout-ms`: Maximum backoff window before aborting session (default: `10000`).
  - Flag `--is-auto-drain`: Automatically resumes streaming upon receive window expansion.
  - Flag `--has-exponential-scale`: Applies exponential multiplier to retry probing intervals.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZeroWindowBackoffEvents (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      initialWindowBytes INTEGER NOT NULL DEFAULT 0,
      backoffDurationMs REAL NOT NULL DEFAULT 0.0,
      isAutoDrain INTEGER NOT NULL DEFAULT 1,
      hasExponentialScale INTEGER NOT NULL DEFAULT 1,
      occurredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZeroWindowBackoffEvent struct {
      ID                  string    `json:"id"`
      SocketFd            int       `json:"socketFd"`
      InitialWindowBytes  int       `json:"initialWindowBytes"`
      BackoffDurationMs   float64   `json:"backoffDurationMs"`
      IsAutoDrain         bool      `json:"isAutoDrain"`
      HasExponentialScale bool      `json:"hasExponentialScale"`
      OccurredAt          time.Time `json:"occurredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/tcp/backoff_controller.go`.

---

### Improvement 1104: Non-Intrusive TCP Zero-Window ACK Probe Dispatcher
- **ID**: `IMP-1104`
- **Cluster**: 111
- **Goal**: Dispatch 1-byte non-intrusive TCP zero-window probe segments during receive window exhaustion to prompt destination MTA TCP stacks for window update ACK packets.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss table row showing probe sequence numbers, RTT response latencies in cyan (`#00E5FF`), and updated window sizes.
  - *Web*: `<window-probe-table>` Web Component displaying probe timing heatmaps using pure HTML/CSS.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline send-probe --socket-fd=<fd> [--probe-interval-ms=<ms>] [--is-stealth-probe] [--has-sequence-sync]`
  - Flag `--socket-fd`: Saturated connection socket.
  - Flag `--probe-interval-ms`: Frequency of probe transmission (default: `250`).
  - Flag `--is-stealth-probe`: Uses standard keep-alive payload rather than raw segment injection.
  - Flag `--has-sequence-sync`: Resynchronizes TCP sequence numbers upon receiving window update.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WindowProbeDispatches (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      probeSequenceNumber INTEGER NOT NULL,
      updatedWindowSize INTEGER NOT NULL DEFAULT 0,
      rttLatencyUs INTEGER NOT NULL DEFAULT 0,
      isStealthProbe INTEGER NOT NULL DEFAULT 1,
      hasSequenceSync INTEGER NOT NULL DEFAULT 1,
      dispatchedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WindowProbeDispatch struct {
      ID                  string    `json:"id"`
      SocketFd            int       `json:"socketFd"`
      ProbeSequenceNumber int64     `json:"probeSequenceNumber"`
      UpdatedWindowSize   int       `json:"updatedWindowSize"`
      RttLatencyUs        int64     `json:"rttLatencyUs"`
      IsStealthProbe      bool      `json:"isStealthProbe"`
      HasSequenceSync     bool      `json:"hasSequenceSync"`
      DispatchedAt        time.Time `json:"dispatchedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tcp/probe_dispatcher.go`.

---

### Improvement 1105: Outbound Socket Buffer Self-Healing & Ring-Drain Optimizer
- **ID**: `IMP-1105`
- **Cluster**: 111
- **Goal**: Dynamically resize userspace memory buffers and flush outbound circular queues when downstream TCP windows collapse, avoiding application thread blocking.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Bubble Tea v2 viewport showing ring buffer fill levels with color-coded ANSI progress bars (green < 50%, amber 50–85%, red > 85%).
  - *Web*: `<ring-buffer-meter>` Web Component rendering interactive SVG circular progress meters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline drain-buffer --socket-fd=<fd> [--target-headroom-kb=<kb>] [--is-aggressive-drain] [--has-drop-prevention]`
  - Flag `--socket-fd`: Socket descriptor to drain.
  - Flag `--target-headroom-kb`: Target minimum free buffer capacity in kilobytes.
  - Flag `--is-aggressive-drain`: Offloads stalled outbound payloads to disk-backed spillover queues.
  - Flag `--has-drop-prevention`: Guarantees zero dropped envelopes during queue reallocation.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketBufferDrains (
      id TEXT PRIMARY KEY,
      socketFd INTEGER NOT NULL,
      flushedBytes INTEGER NOT NULL DEFAULT 0,
      bufferCapacityBytes INTEGER NOT NULL DEFAULT 0,
      isAggressiveDrain INTEGER NOT NULL DEFAULT 0,
      hasDropPrevention INTEGER NOT NULL DEFAULT 1,
      drainedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketBufferDrain struct {
      ID                  string    `json:"id"`
      SocketFd            int       `json:"socketFd"`
      FlushedBytes        int       `json:"flushedBytes"`
      BufferCapacityBytes int       `json:"bufferCapacityBytes"`
      IsAggressiveDrain   bool      `json:"isAggressiveDrain"`
      HasDropPrevention   bool      `json:"hasDropPrevention"`
      DrainedAt           time.Time `json:"drainedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed inside `pkg/warmup/smtp/ring_buffer.go`.

---

### Improvement 1106: Dynamic Pipelined Command Window Scaling Engine
- **ID**: `IMP-1106`
- **Cluster**: 111
- **Goal**: Automatically tune the number of pipelined commands dispatched in a single batch based on TCP round-trip time and historical recipient gateway acknowledgment speed.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss line chart visualizing dynamic command window expansion and contraction over time.
  - *Web*: `<dynamic-window-chart>` Web Component plotting latency vs. pipeline depth using inline SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline tune-window --host=<hostname> [--min-depth=<min>] [--max-depth=<max>] [--is-additive-increase] [--has-rtt-weight]`
  - Flag `--host`: Remote SMTP server domain or IP address.
  - Flag `--min-depth`: Lower bound for pipelined command burst depth (default: `1`).
  - Flag `--max-depth`: Upper limit for pipelined command burst depth (default: `32`).
  - Flag `--is-additive-increase`: Applies AIMD (Additive Increase, Multiplicative Decrease) scaling.
  - Flag `--has-rtt-weight`: Weights window adjustments proportionally to network RTT.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PipelineWindowTunings (
      id TEXT PRIMARY KEY,
      remoteHost TEXT NOT NULL,
      currentWindowDepth INTEGER NOT NULL DEFAULT 4,
      measuredRttMs REAL NOT NULL DEFAULT 0.0,
      isAdditiveIncrease INTEGER NOT NULL DEFAULT 1,
      hasRttWeight INTEGER NOT NULL DEFAULT 1,
      tunedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PipelineWindowTuning struct {
      ID                 string    `json:"id"`
      RemoteHost         string    `json:"remoteHost"`
      CurrentWindowDepth int       `json:"currentWindowDepth"`
      MeasuredRttMs      float64   `json:"measuredRttMs"`
      IsAdditiveIncrease bool      `json:"isAdditiveIncrease"`
      HasRttWeight       bool      `json:"hasRttWeight"`
      TunedAt            time.Time `json:"tunedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Governed in `pkg/warmup/smtp/window_tuner.go`.

---

### Improvement 1107: Asynchronous Multi-Line SMTP Response Demultiplexer
- **ID**: `IMP-1107`
- **Cluster**: 111
- **Goal**: Parse and correlate out-of-order and chunked multi-line SMTP response status codes (e.g., `250-`, `250 `) from pipelined command batches without blocking thread execution.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss formatted diagnostic pane showing pipelined command stream mapped to respective multi-line response blocks with syntax highlights.
  - *Web*: `<smtp-response-demuxer>` Web Component rendering interactive collapsible tree nodes for multi-line replies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline parse-response --input-buffer=<hex> [--is-strict-rfc] [--has-continuation-check]`
  - Flag `--input-buffer`: Raw hex/ASCII buffer received from socket stream.
  - Flag `--is-strict-rfc`: Validates RFC 5321 compliant line terminators (`CRLF`).
  - Flag `--has-continuation-check`: Asserts proper hyphen syntax on multi-line continuation codes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpResponseDemuxLogs (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      commandVerb TEXT NOT NULL,
      responseCode INTEGER NOT NULL,
      responsePayload TEXT NOT NULL,
      isStrictRfc INTEGER NOT NULL DEFAULT 1,
      hasContinuationCheck INTEGER NOT NULL DEFAULT 1,
      parsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpResponseDemuxLog struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      CommandVerb          string    `json:"commandVerb"`
      ResponseCode         int       `json:"responseCode"`
      ResponsePayload      string    `json:"responsePayload"`
      IsStrictRfc          bool      `json:"isStrictRfc"`
      HasContinuationCheck bool      `json:"hasContinuationCheck"`
      ParsedAt             time.Time `json:"parsedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Incorporated into `pkg/warmup/smtp/response_demux.go`.

---

### Improvement 1108: Smtp Pipelining Capability Cache with Negative Expiry
- **ID**: `IMP-1108`
- **Cluster**: 111
- **Goal**: Cache remote MTA `250 PIPELINING` capability discoveries in SQLite with negative TTL caching to eliminate redundant EHLO interrogations on recurring connections.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss lookup card displaying cached MX host capabilities with TTL indicators in cyan (`#00E5FF`).
  - *Web*: `<capability-cache-view>` Web Component rendering interactive host tables with instant search filters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline get-cached-capability --host=<hostname> [--is-fresh-check] [--has-negative-ttl]`
  - Flag `--host`: Remote MX host domain.
  - Flag `--is-fresh-check`: Forces live EHLO probe if cache entry is within 10% of expiration.
  - Flag `--has-negative-ttl`: Stores temporary failure states to avoid hammering non-pipelining hosts.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpCapabilityCaches (
      id TEXT PRIMARY KEY,
      remoteHost TEXT NOT NULL UNIQUE,
      isPipeliningSupported INTEGER NOT NULL DEFAULT 1,
      hasEnhancedStatusCodes INTEGER NOT NULL DEFAULT 1,
      hasStartTls INTEGER NOT NULL DEFAULT 1,
      ttlSeconds INTEGER NOT NULL DEFAULT 86400,
      cachedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpCapabilityCache struct {
      ID                     string    `json:"id"`
      RemoteHost             string    `json:"remoteHost"`
      IsPipeliningSupported  bool      `json:"isPipeliningSupported"`
      HasEnhancedStatusCodes bool      `json:"hasEnhancedStatusCodes"`
      HasStartTls            bool      `json:"hasStartTls"`
      TtlSeconds             int       `json:"ttlSeconds"`
      CachedAt               time.Time `json:"cachedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Persistent cache lookup in `pkg/warmup/smtp/capability_cache.go`.

---

### Improvement 1109: Zero-Window Exhaustion Circuit Breaker & Envelope Re-Router
- **ID**: `IMP-1109`
- **Cluster**: 111
- **Goal**: Automatically trip a circuit breaker when a destination MX sustains zero-window lock for > 5 seconds, diverting queued warming envelopes to secondary MX endpoints without message loss.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss alert modal notifying operators of circuit trip with fallback MX target highlighted in vivid orange (`#FF6D00`).
  - *Web*: `<circuit-breaker-modal>` Web Component showing circuit state transitions (Closed -> Open -> Half-Open).
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline trip-circuit --host=<hostname> [--reason=<reason>] [--is-forced-trip] [--has-reroute-active]`
  - Flag `--host`: Target host to isolate.
  - Flag `--reason`: Textual justification for circuit state mutation.
  - Flag `--is-forced-trip`: Manually forces circuit breaker into Open state for testing.
  - Flag `--has-reroute-active`: Automatically activates standby backup egress tunnels.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CircuitBreakerEvents (
      id TEXT PRIMARY KEY,
      targetHost TEXT NOT NULL,
      circuitState TEXT NOT NULL,
      failureThreshold INTEGER NOT NULL DEFAULT 5,
      isForcedTrip INTEGER NOT NULL DEFAULT 0,
      hasRerouteActive INTEGER NOT NULL DEFAULT 1,
      trippedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CircuitBreakerEvent struct {
      ID               string    `json:"id"`
      TargetHost       string    `json:"targetHost"`
      CircuitState     string    `json:"circuitState"`
      FailureThreshold int       `json:"failureThreshold"`
      IsForcedTrip     bool      `json:"isForcedTrip"`
      HasRerouteActive bool      `json:"hasRerouteActive"`
      TrippedAt        time.Time `json:"trippedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Registered in `pkg/warmup/circuit/breaker.go`.

---

### Improvement 1110: Kernel TCP Zero-Copy eBPF Telemetry Hook
- **ID**: `IMP-1110`
- **Cluster**: 111
- **Goal**: Inject lightweight eBPF tracepoint probes into `tcp_probe` and `sock:inet_sock_set_state` to stream microsecond-accurate TCP receive window metrics directly into GitMap userspace buffers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss real-time HUD streaming kernel probe events with microsecond timestamps and CPU core indices.
  - *Web*: `<ebpf-telemetry-stream>` Web Component showing real-time event logs with zero JS library overhead.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline ebpf-attach --interface=<iface> [--is-ring-buffer] [--has-bpf-verifier-pass]`
  - Flag `--interface`: Network interface to monitor.
  - Flag `--is-ring-buffer`: Allocates modern BPF perf ring buffer for kernel-to-user transmission.
  - Flag `--has-bpf-verifier-pass`: Performs pre-flight verification on eBPF bytecode integrity.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EbpfProbeHooks (
      id TEXT PRIMARY KEY,
      interfaceName TEXT NOT NULL,
      probeName TEXT NOT NULL,
      eventCounter INTEGER NOT NULL DEFAULT 0,
      isRingBuffer INTEGER NOT NULL DEFAULT 1,
      hasBpfVerifierPass INTEGER NOT NULL DEFAULT 1,
      attachedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EbpfProbeHook struct {
      ID                 string    `json:"id"`
      InterfaceName      string    `json:"interfaceName"`
      ProbeName          string    `json:"probeName"`
      EventCounter       int64     `json:"eventCounter"`
      IsRingBuffer       bool      `json:"isRingBuffer"`
      HasBpfVerifierPass bool      `json:"hasBpfVerifierPass"`
      AttachedAt         time.Time `json:"attachedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/ebpf/kernel_hook.go`.

---

# Cluster 112: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (1111–1120)

---

### Improvement 1111: Adaptive Multi-Tenant SMTP Socket Reuse Pool
- **ID**: `IMP-1111`
- **Cluster**: 112 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal**: Maintain warmed, pre-authenticated, and idle TLS 1.3 socket pools per target ISP domain to eliminate TCP 3-way handshake and TLS negotiation overhead for scheduled envelope dispatches.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss grid table displaying active, idle, and max pool capacities per destination ISP with electric cyan (`#00E5FF`) progress bars.
  - *Web*: `<connection-pool-matrix>` Web Component rendering visual connection slots with real-time status pulses.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool status --isp=<domain> [--max-idle=<count>] [--is-tls-warmed] [--has-keepalive]`
  - Flag `--isp`: Destination ISP domain (e.g. `google.com`, `yahoo.com`).
  - Flag `--max-idle`: Upper limit of retained idle sockets in pool (default: `10`).
  - Flag `--is-tls-warmed`: Asserts that idle sockets maintain completed TLS handshakes.
  - Flag `--has-keepalive`: Enables periodic TCP keep-alive probe transmission.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ConnectionPoolStates (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL UNIQUE,
      activeSockets INTEGER NOT NULL DEFAULT 0,
      idleSockets INTEGER NOT NULL DEFAULT 0,
      maxPoolCapacity INTEGER NOT NULL DEFAULT 20,
      isTlsWarmed INTEGER NOT NULL DEFAULT 1,
      hasKeepalive INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ConnectionPoolState struct {
      ID              string    `json:"id"`
      IspDomain       string    `json:"ispDomain"`
      ActiveSockets   int       `json:"activeSockets"`
      IdleSockets     int       `json:"idleSockets"`
      MaxPoolCapacity int       `json:"maxPoolCapacity"`
      IsTlsWarmed     bool      `json:"isTlsWarmed"`
      HasKeepalive    bool      `json:"hasKeepalive"`
      UpdatedAt       time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Controlled via `pkg/warmup/pool/manager.go`.

---

### Improvement 1112: Sub-Millisecond Socket Degradation & TCP Half-Open Detector
- **ID**: `IMP-1112`
- **Cluster**: 112
- **Goal**: Actively detect TCP half-open states and stealth connection drops caused by stateful firewall timeout silent drops before attempting envelope dispatch.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss alert badge highlighting degraded sockets with warning amber (`#FFD600`) and half-closed icons.
  - *Web*: `<socket-health-card>` Web Component displaying real-time TCP state transitions via SVG icons.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool check-degradation --socket-id=<id> [--is-active-probe] [--has-epoll-hup-check]`
  - Flag `--socket-id`: Unique identifier of the connection in the pool.
  - Flag `--is-active-probe`: Emits non-destructive RSET command to confirm remote socket readiness.
  - Flag `--has-epoll-hup-check`: Inspects kernel epoll `EPOLLRDHUP` event status.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketDegradationLogs (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      tcpState TEXT NOT NULL,
      isDegraded INTEGER NOT NULL DEFAULT 0,
      hasEpollHupCheck INTEGER NOT NULL DEFAULT 1,
      isActiveProbe INTEGER NOT NULL DEFAULT 0,
      checkedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketDegradationLog struct {
      ID               string    `json:"id"`
      SocketID         string    `json:"socketId"`
      TcpState         string    `json:"tcpState"`
      IsDegraded       bool      `json:"isDegraded"`
      HasEpollHupCheck bool      `json:"hasEpollHupCheck"`
      IsActiveProbe    bool      `json:"isActiveProbe"`
      CheckedAt        time.Time `json:"checkedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Placed in `pkg/warmup/pool/health_checker.go`.

---

### Improvement 1113: Cross-ISP Handshake Latency Profiler & Scoring Engine
- **ID**: `IMP-1113`
- **Cluster**: 112
- **Goal**: Score and compare TCP and TLS handshake latency across major destination ISPs, continuously optimizing pool allocation toward lower-latency ingress gateways.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss bar chart ranking ISP handshake performance in ascending latency order.
  - *Web*: `<isp-latency-ranker>` Web Component rendering responsive SVG horizontal bar charts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool benchmark-handshakes [--sample-count=<count>] [--is-detailed-tls] [--has-ranking-update]`
  - Flag `--sample-count`: Number of handshakes initiated per target domain (default: `5`).
  - Flag `--is-detailed-tls`: Breaks down TLS 1.3 key exchange vs. certificate transmission phases.
  - Flag `--has-ranking-update`: Automatically updates internal routing priority matrix.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspHandshakeProfiles (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL,
      averageTcpRttMs REAL NOT NULL DEFAULT 0.0,
      averageTlsDurationMs REAL NOT NULL DEFAULT 0.0,
      compositeScore INTEGER NOT NULL DEFAULT 100,
      isDetailedTls INTEGER NOT NULL DEFAULT 1,
      hasRankingUpdate INTEGER NOT NULL DEFAULT 1,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspHandshakeProfile struct {
      ID                   string    `json:"id"`
      IspDomain            string    `json:"ispDomain"`
      AverageTcpRttMs      float64   `json:"averageTcpRttMs"`
      AverageTlsDurationMs float64   `json:"averageTlsDurationMs"`
      CompositeScore       int       `json:"compositeScore"`
      IsDetailedTls        bool      `json:"isDetailedTls"`
      HasRankingUpdate     bool      `json:"hasRankingUpdate"`
      MeasuredAt           time.Time `json:"measuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/pool/profiler.go`.

---

### Improvement 1114: Dynamic TCP Keep-Alive Cadence Calibrator
- **ID**: `IMP-1114`
- **Cluster**: 112
- **Goal**: Compute the optimal TCP keep-alive transmission interval dynamically per ISP to prevent firewall NAT session drops while avoiding ISP anti-ping heuristics.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss numeric readout of calibrated keep-alive parameters (idle time, probe count, interval).
  - *Web*: `<keepalive-calibrator-card>` Web Component displaying interactive slider gauges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool calibrate-keepalive --isp=<domain> [--is-adaptive] [--has-nat-drift-defense]`
  - Flag `--isp`: Destination provider to evaluate.
  - Flag `--is-adaptive`: Adjusts keep-alive frequency based on observed NAT drop timeouts.
  - Flag `--has-nat-drift-defense`: Employs jittered timings to prevent rhythmic traffic detection.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS KeepAliveCalibrations (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL UNIQUE,
      optimalIntervalSeconds INTEGER NOT NULL DEFAULT 45,
      detectedNatTimeoutSeconds INTEGER NOT NULL DEFAULT 60,
      isAdaptive INTEGER NOT NULL DEFAULT 1,
      hasNatDriftDefense INTEGER NOT NULL DEFAULT 1,
      calibratedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type KeepAliveCalibration struct {
      ID                        string    `json:"id"`
      IspDomain                 string    `json:"ispDomain"`
      OptimalIntervalSeconds    int       `json:"optimalIntervalSeconds"`
      DetectedNatTimeoutSeconds int       `json:"detectedNatTimeoutSeconds"`
      IsAdaptive                bool      `json:"isAdaptive"`
      HasNatDriftDefense        bool      `json:"hasNatDriftDefense"`
      CalibratedAt              time.Time `json:"calibratedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/pool/keepalive.go`.

---

### Improvement 1115: Socket Eviction & Graceful TLS Teardown Orchestrator
- **ID**: `IMP-1115`
- **Cluster**: 112
- **Goal**: Retire and teardown idle sockets exceeding max lifetime using clean TLS close_notify and SMTP QUIT commands rather than abrupt TCP RST resets.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss log viewer displaying graceful teardown handshakes in subtle lavender (`#B388FF`).
  - *Web*: `<socket-teardown-log>` Web Component detailing socket lifecycle history in searchable tables.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool evict-expired [--max-age-seconds=<sec>] [--is-graceful] [--has-quit-command]`
  - Flag `--max-age-seconds`: Maximum retention time before retiring idle connection (default: `300`).
  - Flag `--is-graceful`: Transmits TLS close_notify alert before closing TCP socket.
  - Flag `--has-quit-command`: Dispatches SMTP `QUIT` command to finalize remote server session.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketEvictionLogs (
      id TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      ageSeconds INTEGER NOT NULL,
      teardownDurationMs REAL NOT NULL DEFAULT 0.0,
      isGraceful INTEGER NOT NULL DEFAULT 1,
      hasQuitCommand INTEGER NOT NULL DEFAULT 1,
      evictedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketEvictionLog struct {
      ID                 string    `json:"id"`
      SocketID           string    `json:"socketId"`
      AgeSeconds         int       `json:"ageSeconds"`
      TeardownDurationMs float64   `json:"teardownDurationMs"`
      IsGraceful         bool      `json:"isGraceful"`
      HasQuitCommand     bool      `json:"hasQuitCommand"`
      EvictedAt          time.Time `json:"evictedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked inside `pkg/warmup/pool/eviction.go`.

---

### Improvement 1116: Multipath TCP Subflow State Machine for Connection Pools
- **ID**: `IMP-1116`
- **Cluster**: 112
- **Goal**: Maintain and orchestrate MPTCP subflows across bonded local network adapters inside the connection pool, dynamically switching egress subflows if interface degradation occurs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss subflow topology map indicating active physical interfaces, subflow IDs, and individual RTTs.
  - *Web*: `<mptcp-subflow-map>` Web Component rendering network interface topologies in dynamic SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool mptcp-subflows --pool-id=<id> [--is-multipath-active] [--has-backup-path]`
  - Flag `--pool-id`: Connection pool instance under inspection.
  - Flag `--is-multipath-active`: Validates that kernel MPTCP socket options are enabled.
  - Flag `--has-backup-path`: Configures secondary uplink in hot-standby backup mode.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MptcpSubflowStates (
      id TEXT PRIMARY KEY,
      poolId TEXT NOT NULL,
      subflowId INTEGER NOT NULL,
      interfaceName TEXT NOT NULL,
      rttMs REAL NOT NULL DEFAULT 0.0,
      isMultipathActive INTEGER NOT NULL DEFAULT 1,
      hasBackupPath INTEGER NOT NULL DEFAULT 0,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MptcpSubflowState struct {
      ID                string    `json:"id"`
      PoolID            string    `json:"poolId"`
      SubflowID         int       `json:"subflowId"`
      InterfaceName     string    `json:"interfaceName"`
      RttMs             float64   `json:"rttMs"`
      IsMultipathActive bool      `json:"isMultipathActive"`
      HasBackupPath     bool      `json:"hasBackupPath"`
      MeasuredAt        time.Time `json:"measuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/pool/mptcp_subflow.go`.

---

### Improvement 1117: Autonomous ISP Greylist Telemetry & Backoff Queue Balancer
- **ID**: `IMP-1117`
- **Cluster**: 112
- **Goal**: Recognize 451/450 greylisting deferral response signatures from newly warmed IP ranges and automatically queue deferred envelopes for re-dispatch at the exact greylist expiration interval.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss table presenting detected greylisting hosts, retry countdown clocks, and target mailbox counts.
  - *Web*: `<greylist-queue-tracker>` Web Component showing real-time queue drain forecasts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool inspect-greylists [--isp=<domain>] [--is-auto-reschedule] [--has-signature-match]`
  - Flag `--isp`: Destination provider filter.
  - Flag `--is-auto-reschedule`: Moves deferred envelopes into targeted greylist replay buffers.
  - Flag `--has-signature-match`: Employs regex matching on 4xx textual response banners.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS GreylistTelemetryLogs (
      id TEXT PRIMARY KEY,
      targetIsp TEXT NOT NULL,
      responseBanner TEXT NOT NULL,
      retryDelaySeconds INTEGER NOT NULL DEFAULT 900,
      isAutoReschedule INTEGER NOT NULL DEFAULT 1,
      hasSignatureMatch INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type GreylistTelemetryLog struct {
      ID                string    `json:"id"`
      TargetIsp         string    `json:"targetIsp"`
      ResponseBanner    string    `json:"responseBanner"`
      RetryDelaySeconds int       `json:"retryDelaySeconds"`
      IsAutoReschedule  bool      `json:"isAutoReschedule"`
      HasSignatureMatch bool      `json:"hasSignatureMatch"`
      RecordedAt        time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/pool/greylist.go`.

---

### Improvement 1118: Predictive Connection Pre-Warming Scheduler
- **ID**: `IMP-1118`
- **Cluster**: 112
- **Goal**: Anticipate outbound warming schedule batches by initiating TCP handshakes and TLS sessions 2.5 seconds before scheduled dispatch timestamps, guaranteeing zero-latency launches.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss timeline bar indicating pre-warming milestones ahead of execution schedule.
  - *Web*: `<prewarm-scheduler-view>` Web Component visualizing pre-warming horizons with CSS animations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool schedule-prewarm --lead-time-ms=<ms> [--is-predictive] [--has-capacity-check]`
  - Flag `--lead-time-ms`: Milliseconds prior to dispatch to trigger socket creation (default: `2500`).
  - Flag `--is-predictive`: Evaluates warming calendar queues to size pre-warm socket counts.
  - Flag `--has-capacity-check`: Asserts local resource headroom before opening pre-warm connections.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PrewarmScheduleEvents (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      prewarmSocketsCount INTEGER NOT NULL DEFAULT 2,
      leadTimeMs INTEGER NOT NULL DEFAULT 2500,
      isPredictive INTEGER NOT NULL DEFAULT 1,
      hasCapacityCheck INTEGER NOT NULL DEFAULT 1,
      scheduledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PrewarmScheduleEvent struct {
      ID                  string    `json:"id"`
      TargetDomain        string    `json:"targetDomain"`
      PrewarmSocketsCount int       `json:"prewarmSocketsCount"`
      LeadTimeMs          int       `json:"leadTimeMs"`
      IsPredictive        bool      `json:"isPredictive"`
      HasCapacityCheck    bool      `json:"hasCapacityCheck"`
      ScheduledAt         time.Time `json:"scheduledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/pool/prewarm_scheduler.go`.

---

### Improvement 1119: Connection Pool Throttling Governor & Concurrency Clamp
- **ID**: `IMP-1119`
- **Cluster**: 112
- **Goal**: Enforce strict per-IP and per-destination concurrency caps, preventing accidental surges that trip ISP anti-abuse filters and damage reputation scores.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss gauge displaying current concurrent active connections against ISP-specific ceiling caps.
  - *Web*: `<concurrency-clamp-meter>` Web Component rendering interactive SVG limit indicators.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool set-concurrency-limit --isp=<domain> --max-concurrent=<limit> [--is-enforced] [--has-queue-spill]`
  - Flag `--isp`: Destination provider.
  - Flag `--max-concurrent`: Maximum concurrent open sockets allowed.
  - Flag `--is-enforced`: Blocks requests attempting to exceed cap.
  - Flag `--has-queue-spill`: Spills excess requests into prioritized local queue.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ConcurrencyLimitRules (
      id TEXT PRIMARY KEY,
      ispDomain TEXT NOT NULL UNIQUE,
      maxConcurrentSockets INTEGER NOT NULL DEFAULT 5,
      currentActiveSockets INTEGER NOT NULL DEFAULT 0,
      isEnforced INTEGER NOT NULL DEFAULT 1,
      hasQueueSpill INTEGER NOT NULL DEFAULT 1,
      configuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ConcurrencyLimitRule struct {
      ID                   string    `json:"id"`
      IspDomain            string    `json:"ispDomain"`
      MaxConcurrentSockets int       `json:"maxConcurrentSockets"`
      CurrentActiveSockets int       `json:"currentActiveSockets"`
      IsEnforced           bool      `json:"isEnforced"`
      HasQueueSpill        bool      `json:"hasQueueSpill"`
      ConfiguredAt         time.Time `json:"configuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Enforced in `pkg/warmup/pool/governor.go`.

---

### Improvement 1120: Cross-Interface Egress IP Rotation Matrix
- **ID**: `IMP-1120`
- **Cluster**: 112
- **Goal**: Rotate outbound connection binds across multiple local IPv4/IPv6 addresses, balancing dispatch volume across warming IP allocations according to configured maturation curves.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss multi-column table displaying local IP allocation pool, current volume allocations, and rotation health.
  - *Web*: `<ip-rotation-matrix>` Web Component rendering dynamic heatmaps of IP utilization.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pool rotate-ip [--pool-group=<group>] [--is-round-robin] [--has-ipv6-dual-stack]`
  - Flag `--pool-group`: Group of IP addresses to rotate.
  - Flag `--is-round-robin`: Rotates addresses sequentially; otherwise uses reputation-weighted selection.
  - Flag `--has-ipv6-dual-stack`: Includes verified IPv6 addresses in outbound rotation.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EgressIpRotations (
      id TEXT PRIMARY KEY,
      ipAddress TEXT NOT NULL UNIQUE,
      poolGroup TEXT NOT NULL,
      dispatchedCount INTEGER NOT NULL DEFAULT 0,
      isRoundRobin INTEGER NOT NULL DEFAULT 0,
      hasIpv6DualStack INTEGER NOT NULL DEFAULT 1,
      lastUsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EgressIpRotation struct {
      ID               string    `json:"id"`
      IPAddress        string    `json:"ipAddress"`
      PoolGroup        string    `json:"poolGroup"`
      DispatchedCount  int64     `json:"dispatchedCount"`
      IsRoundRobin     bool      `json:"isRoundRobin"`
      HasIpv6DualStack bool      `json:"hasIpv6DualStack"`
      LastUsedAt       time.Time `json:"lastUsedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Governed in `pkg/warmup/pool/ip_rotator.go`.

---

# Cluster 113: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1121–1130)

---

### Improvement 1121: Sub-Millisecond Five-Phase SMTP Latency Flamegraph Renderer
- **ID**: `IMP-1121`
- **Cluster**: 113 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Render a detailed 5-phase latency breakdown (DNS, TCP, TLS, PIPE, DATA) for every outbound transaction directly in the terminal via Lipgloss TrueColor horizontal blocks and in the browser via zero-CDN inline SVG.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-density ANSI flamegraph with Electric Cyan (`#00E5FF`, DNS), Spring Green (`#00E676`, TCP), Bright Purple (`#7C4DFF`, TLS), Amber (`#FFD600`, PIPE), and Neon Magenta (`#E040FB`, DATA).
  - *Web*: Custom Web Component `<smtp-flamegraph-card>` rendering pure SVG vector bars with interactive hover tooltips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui render-flamegraph --transaction-id=<id> [--is-ascii-fallback] [--has-percentage-labels]`
  - Flag `--transaction-id`: Outbound warming transaction record ID.
  - Flag `--is-ascii-fallback`: Renders simplified 8-color ASCII blocks on legacy terminal emulators.
  - Flag `--has-percentage-labels`: Annotates flamegraph blocks with percentage of total transaction duration.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpPhaseFlamegraphs (
      id TEXT PRIMARY KEY,
      transactionId TEXT NOT NULL UNIQUE,
      dnsDurationMs REAL NOT NULL DEFAULT 0.0,
      tcpDurationMs REAL NOT NULL DEFAULT 0.0,
      tlsDurationMs REAL NOT NULL DEFAULT 0.0,
      pipeDurationMs REAL NOT NULL DEFAULT 0.0,
      dataDurationMs REAL NOT NULL DEFAULT 0.0,
      totalDurationMs REAL NOT NULL DEFAULT 0.0,
      isAsciiFallback INTEGER NOT NULL DEFAULT 0,
      hasPercentageLabels INTEGER NOT NULL DEFAULT 1,
      renderedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpPhaseFlamegraph struct {
      ID                  string    `json:"id"`
      TransactionID       string    `json:"transactionId"`
      DnsDurationMs       float64   `json:"dnsDurationMs"`
      TcpDurationMs       float64   `json:"tcpDurationMs"`
      TlsDurationMs       float64   `json:"tlsDurationMs"`
      PipeDurationMs      float64   `json:"pipeDurationMs"`
      DataDurationMs      float64   `json:"dataDurationMs"`
      TotalDurationMs     float64   `json:"totalDurationMs"`
      IsAsciiFallback     bool      `json:"isAsciiFallback"`
      HasPercentageLabels bool      `json:"hasPercentageLabels"`
      RenderedAt          time.Time `json:"renderedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/flamegraph.go`.

---

### Improvement 1122: Bounded Real-Time Jitter Waveform TUI Visualizer
- **ID**: `IMP-1122`
- **Cluster**: 113
- **Goal**: Stream continuous microsecond-level packet jitter waveforms across active SMTP sessions using Unicode Braille patterns (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`) in Bubble Tea v2 viewports.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Multi-row Braille waveform oscillogram with dynamic amplitude coloring (green for < 1ms jitter, amber for 1–5ms, red for > 5ms).
  - *Web*: `<jitter-waveform-canvas>` Web Component streaming HTML5 Canvas or SVG path lines with zero external dependencies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui jitter-stream --session-id=<id> [--window-seconds=<sec>] [--is-braille-mode] [--has-anomaly-highlights]`
  - Flag `--session-id`: Target SMTP session.
  - Flag `--window-seconds`: Rolling time window in seconds (default: `30`).
  - Flag `--is-braille-mode`: Uses 2x4 dot Unicode Braille glyphs for high-resolution plotting.
  - Flag `--has-anomaly-highlights`: Inverts terminal colors on anomalous jitter spikes (> 3 sigma).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterStreamMetrics (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      windowDurationSeconds INTEGER NOT NULL DEFAULT 30,
      peakJitterMicroseconds INTEGER NOT NULL DEFAULT 0,
      isBrailleMode INTEGER NOT NULL DEFAULT 1,
      hasAnomalyHighlights INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterStreamMetric struct {
      ID                     string    `json:"id"`
      SessionID              string    `json:"sessionId"`
      WindowDurationSeconds  int       `json:"windowDurationSeconds"`
      PeakJitterMicroseconds int       `json:"peakJitterMicroseconds"`
      IsBrailleMode          bool      `json:"isBrailleMode"`
      HasAnomalyHighlights   bool      `json:"hasAnomalyHighlights"`
      RecordedAt             time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/tui/jitter_waveform.go`.

---

### Improvement 1123: Zero-CDN Shadow DOM Web Flamegraph Component
- **ID**: `IMP-1123`
- **Cluster**: 113
- **Goal**: Provide a self-contained Web Component encapsulation `<smtp-flamegraph>` using isolated Shadow DOM styles and native SVG elements, requiring zero external scripts or CDN fonts.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: N/A (Web Component runtime output verification CLI).
  - *Web*: Custom Element `<smtp-flamegraph>` with open Shadow Root, zero CSS leakage, and responsive SVG coordinate scaling.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui export-flamegraph-component --out-file=<path> [--is-minified] [--has-inline-css]`
  - Flag `--out-file`: Target destination path for compiled standalone HTML/JS bundle.
  - Flag `--is-minified`: Strips whitespace and comments from generated JS/SVG markup.
  - Flag `--has-inline-css`: Embeds all CSS styling within Shadow DOM `<style>` tags.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WebComponentArtifacts (
      id TEXT PRIMARY KEY,
      componentTag TEXT NOT NULL UNIQUE,
      bundleSizeBytes INTEGER NOT NULL DEFAULT 0,
      isMinified INTEGER NOT NULL DEFAULT 1,
      hasInlineCss INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WebComponentArtifact struct {
      ID              string    `json:"id"`
      ComponentTag    string    `json:"componentTag"`
      BundleSizeBytes int       `json:"bundleSizeBytes"`
      IsMinified      bool      `json:"isMinified"`
      HasInlineCss    bool      `json:"hasInlineCss"`
      ExportedAt      time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/web/component_exporter.go`.

---

### Improvement 1124: Bubble Tea v2 Interactive Warmup HUD Dashboard
- **ID**: `IMP-1124`
- **Cluster**: 113
- **Goal**: Full-screen interactive terminal HUD built on Bubble Tea v2 with keyboard navigation (`tab`, `h/j/k/l`, `q`), multi-pane layout, real-time metrics streaming, and zero flicker.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss styled 3-pane layout: Left pane for active session tree, top right for flamegraph, bottom right for jitter waveform.
  - *Web*: `<warmup-hud-replica>` Web Component providing 1:1 mirroring of the terminal HUD layout in browser windows.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui hud [--refresh-fps=<fps>] [--is-mouse-enabled] [--has-alt-screen]`
  - Flag `--refresh-fps`: Terminal render frequency (default: `30`).
  - Flag `--is-mouse-enabled`: Captures mouse scroll and click events inside terminal viewports.
  - Flag `--has-alt-screen`: Switches terminal buffer to alternate screen mode, restoring prior buffer on exit.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TerminalHudSessions (
      id TEXT PRIMARY KEY,
      terminalCols INTEGER NOT NULL DEFAULT 80,
      terminalRows INTEGER NOT NULL DEFAULT 24,
      refreshFps INTEGER NOT NULL DEFAULT 30,
      isMouseEnabled INTEGER NOT NULL DEFAULT 1,
      hasAltScreen INTEGER NOT NULL DEFAULT 1,
      launchedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TerminalHudSession struct {
      ID             string    `json:"id"`
      TerminalCols   int       `json:"terminalCols"`
      TerminalRows   int       `json:"terminalRows"`
      RefreshFps     int       `json:"refreshFps"`
      IsMouseEnabled bool      `json:"isMouseEnabled"`
      HasAltScreen   bool      `json:"hasAltScreen"`
      LaunchedAt     time.Time `json:"launchedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Core command entry point in `pkg/warmup/tui/hud_model.go`.

---

### Improvement 1125: ANSI TrueColor Palette Normalizer & Contrast Checker
- **ID**: `IMP-1125`
- **Cluster**: 113
- **Goal**: Audit terminal color capabilities (`COLORTERM=truecolor`), automatically mapping 24-bit RGB values to 256-color or 16-color ANSI equivalents while guaranteeing WCAG 2.1 AA text contrast.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Color swatch test grid displaying normalized tokens with contrast ratio overlays (e.g. `4.5:1 PASS`).
  - *Web*: `<palette-contrast-card>` Web Component demonstrating color adaptability across light and dark system themes.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui check-palette [--is-strict-wcag] [--has-auto-fallback]`
  - Flag `--is-strict-wcag`: Enforces minimum contrast ratio of 4.5:1 against current terminal background.
  - Flag `--has-auto-fallback`: Automatically downsamples palette to 256 or 16 colors when TrueColor is unavailable.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TerminalPaletteConfigs (
      id TEXT PRIMARY KEY,
      colorDepthBits INTEGER NOT NULL DEFAULT 24,
      detectedColorTerm TEXT NOT NULL,
      isStrictWcag INTEGER NOT NULL DEFAULT 1,
      hasAutoFallback INTEGER NOT NULL DEFAULT 1,
      testedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TerminalPaletteConfig struct {
      ID                string    `json:"id"`
      ColorDepthBits    int       `json:"colorDepthBits"`
      DetectedColorTerm string    `json:"detectedColorTerm"`
      IsStrictWcag      bool      `json:"isStrictWcag"`
      HasAutoFallback   bool      `json:"hasAutoFallback"`
      TestedAt          time.Time `json:"testedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/palette_normalizer.go`.

---

### Improvement 1126: High-Frequency Sub-Millisecond Jitter Histogram TUI
- **ID**: `IMP-1126`
- **Cluster**: 113
- **Goal**: Aggregate microsecond latency variance across all pipelined commands and render real-time logarithmic bin distribution histograms in terminal and web dashboards.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Vertical bar histogram using Lipgloss gradient shading from Spring Green (low jitter) to Crimson Red (high jitter).
  - *Web*: `<jitter-histogram-view>` Web Component rendering interactive SVG bar charts with logarithmic scaling.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui jitter-histogram --session-id=<id> [--bin-count=<bins>] [--is-log-scale] [--has-p99-marker]`
  - Flag `--session-id`: Target session identifier.
  - Flag `--bin-count`: Number of histogram distribution buckets (default: `16`).
  - Flag `--is-log-scale`: Applies logarithmic y-axis scaling for wide frequency ranges.
  - Flag `--has-p99-marker`: Draws a prominent visual marker at the 99th percentile boundary.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterHistogramSnapshots (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      p50LatencyUs INTEGER NOT NULL DEFAULT 0,
      p95LatencyUs INTEGER NOT NULL DEFAULT 0,
      p99LatencyUs INTEGER NOT NULL DEFAULT 0,
      isLogScale INTEGER NOT NULL DEFAULT 1,
      hasP99Marker INTEGER NOT NULL DEFAULT 1,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterHistogramSnapshot struct {
      ID           string    `json:"id"`
      SessionID    string    `json:"sessionId"`
      P50LatencyUs int       `json:"p50LatencyUs"`
      P95LatencyUs int       `json:"p95LatencyUs"`
      P99LatencyUs int       `json:"p99LatencyUs"`
      IsLogScale   bool      `json:"isLogScale"`
      HasP99Marker bool      `json:"hasP99Marker"`
      GeneratedAt  time.Time `json:"generatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/jitter_histogram.go`.

---

### Improvement 1127: Terminal Viewport Window Resizing & Dynamic Reflow Engine
- **ID**: `IMP-1127`
- **Cluster**: 113
- **Goal**: Listen for OS `SIGWINCH` window resize signals and dynamically recompute flamegraph column widths, text truncation boundaries, and waveform aspect ratios without tearing or layout distortion.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Smooth, artifact-free re-layout upon terminal resizing; gracefully collapses panes into vertical stack when width < 80 columns.
  - *Web*: Responsive container queries `@container` adjusting Shadow DOM SVG layout to parent DOM sizing.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui test-reflow --width=<cols> --height=<rows> [--is-dry-run] [--has-overflow-protection]`
  - Flag `--width`: Virtual column width to simulate.
  - Flag `--height`: Virtual row count to simulate.
  - Flag `--is-dry-run`: Tests reflow layout calculations without emitting ANSI sequences.
  - Flag `--has-overflow-protection`: Asserts that content never spills beyond specified dimensions.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TerminalResizeEvents (
      id TEXT PRIMARY KEY,
      previousCols INTEGER NOT NULL,
      newCols INTEGER NOT NULL,
      previousRows INTEGER NOT NULL,
      newRows INTEGER NOT NULL,
      isDryRun INTEGER NOT NULL DEFAULT 0,
      hasOverflowProtection INTEGER NOT NULL DEFAULT 1,
      resizedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TerminalResizeEvent struct {
      ID                    string    `json:"id"`
      PreviousCols          int       `json:"previousCols"`
      NewCols               int       `json:"newCols"`
      PreviousRows          int       `json:"previousRows"`
      NewRows               int       `json:"newRows"`
      IsDryRun              bool      `json:"isDryRun"`
      HasOverflowProtection bool      `json:"hasOverflowProtection"`
      ResizedAt             time.Time `json:"resizedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tui/reflow_manager.go`.

---

### Improvement 1128: SVG Flamegraph Snapshot Exporter for CI/CD Reports
- **ID**: `IMP-1128`
- **Cluster**: 113
- **Goal**: Export standalone vector SVG flamegraph snapshots of completed warming batches directly to disk for embedding into CI/CD build summaries and Git pull request comments.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Status message confirming file export path and file size in bytes with green checkmark.
  - *Web*: Pure SVG vector file ready for browser preview or direct image rendering.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui export-svg --batch-id=<id> --out-file=<path> [--is-standalone] [--has-dark-theme]`
  - Flag `--batch-id`: Identifier of completed warming batch.
  - Flag `--out-file`: Target filesystem path for the output `.svg` file.
  - Flag `--is-standalone`: Includes XML declaration and embedded font definitions.
  - Flag `--has-dark-theme`: Styles background and text elements for dark mode embedding.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SvgExportRecords (
      id TEXT PRIMARY KEY,
      batchId TEXT NOT NULL,
      filePath TEXT NOT NULL,
      fileSizeBytes INTEGER NOT NULL DEFAULT 0,
      isStandalone INTEGER NOT NULL DEFAULT 1,
      hasDarkTheme INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SvgExportRecord struct {
      ID            string    `json:"id"`
      BatchID       string    `json:"batchId"`
      FilePath      string    `json:"filePath"`
      FileSizeBytes int       `json:"fileSizeBytes"`
      IsStandalone  bool      `json:"isStandalone"`
      HasDarkTheme  bool      `json:"hasDarkTheme"`
      ExportedAt    time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/tui/svg_exporter.go`.

---

### Improvement 1129: Multi-Session Concurrent Latency Comparison TUI
- **ID**: `IMP-1129`
- **Cluster**: 113
- **Goal**: Compare latency flamegraphs and jitter metrics across up to 4 concurrent warming sessions simultaneously in a side-by-side split terminal view.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss 4-quadrant layout displaying synchronous telemetry streams across distinct ISP targets.
  - *Web*: `<multi-session-comparator>` Web Component providing responsive CSS grid comparisons.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui compare-sessions --session-ids=<id1,id2,...> [--is-sync-scroll] [--has-differential-delta]`
  - Flag `--session-ids`: Comma-separated list of session IDs to compare (max 4).
  - Flag `--is-sync-scroll`: Synchronizes timeline scrolling across all quadrants.
  - Flag `--has-differential-delta`: Displays difference deltas relative to baseline session.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MultiSessionComparisons (
      id TEXT PRIMARY KEY,
      sessionIds TEXT NOT NULL,
      baselineSessionId TEXT NOT NULL,
      isSyncScroll INTEGER NOT NULL DEFAULT 1,
      hasDifferentialDelta INTEGER NOT NULL DEFAULT 1,
      comparedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MultiSessionComparison struct {
      ID                   string    `json:"id"`
      SessionIDs           string    `json:"sessionIds"`
      BaselineSessionID    string    `json:"baselineSessionId"`
      IsSyncScroll         bool      `json:"isSyncScroll"`
      HasDifferentialDelta bool      `json:"hasDifferentialDelta"`
      ComparedAt           time.Time `json:"comparedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/tui/session_comparator.go`.

---

### Improvement 1130: Zero-Flicker Double-Buffered Terminal Render Engine
- **ID**: `IMP-1130`
- **Cluster**: 113
- **Goal**: Eliminate visual tearing, blinking, and partial screen repaints by maintaining an in-memory double-buffer ANSI character grid, calculating diff deltas and emitting only changed screen coordinates.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Rock-solid 60 FPS terminal UI updates with zero flicker during high-speed streaming data bursts.
  - *Web*: Equivalent to browser `requestAnimationFrame` zero-jank screen updates.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui test-render-buffer [--fps=<fps>] [--is-double-buffer] [--has-diff-optimization]`
  - Flag `--fps`: Render frame rate target (default: `60`).
  - Flag `--is-double-buffer`: Enables front/back virtual buffer swapping.
  - Flag `--has-diff-optimization`: Emits only modified character cells via cursor positioning escape codes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RenderBufferMetrics (
      id TEXT PRIMARY KEY,
      frameNumber INTEGER NOT NULL,
      cellsUpdated INTEGER NOT NULL DEFAULT 0,
      renderDurationUs INTEGER NOT NULL DEFAULT 0,
      isDoubleBuffer INTEGER NOT NULL DEFAULT 1,
      hasDiffOptimization INTEGER NOT NULL DEFAULT 1,
      renderedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RenderBufferMetric struct {
      ID                  string    `json:"id"`
      FrameNumber         int64     `json:"frameNumber"`
      CellsUpdated        int       `json:"cellsUpdated"`
      RenderDurationUs    int       `json:"renderDurationUs"`
      IsDoubleBuffer      bool      `json:"isDoubleBuffer"`
      HasDiffOptimization bool      `json:"hasDiffOptimization"`
      RenderedAt          time.Time `json:"renderedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/tui/render_buffer.go`.
