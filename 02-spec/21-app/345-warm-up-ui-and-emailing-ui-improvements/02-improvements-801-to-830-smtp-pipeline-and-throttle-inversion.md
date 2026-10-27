# Improvements 801 to 830: SMTP Pipeline Multiplexing, Throttle Inversion & Jitter Waveform TUI

This document provides the exhaustive architectural specification for Improvements 801 through 830 across Clusters 81, 82, and 83 of Sequence 345.

---

# Cluster 81: Autonomous SMTP Pipeline Multiplexing & Adaptive Connection Teardown (801–810)

---

### Improvement 801: Autonomous SMTP Pipeline Session Multiplexer
- **ID**: `IMP-801`
- **Cluster**: 81 (Autonomous SMTP Pipeline Multiplexing & Adaptive Connection Teardown)
- **Goal**: Interleaves multiple outbound email transactions across a single persistent, authenticated, and TLS-encrypted SMTP session using RFC 2920 command pipelining without blocking on per-command acknowledgments.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual connection matrix rendered via Lipgloss displaying multiplexed transaction slots per open socket. Uses cyan (`#00F0FF`) for active slots and muted gray (`#4A5568`) for idle slots.
  - *Web*: Custom Web Component `<smtp-session-multiplexer>` with SVG slot lanes showing real-time concurrent message processing without external CDNs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline multiplex --target-mx=<host> [--max-streams=<n>] [--is-dry-run] [--has-telemetry]`
  - Flag `--max-streams`: Maximum concurrent transactional streams multiplexed onto a single TCP socket (default: 5).
  - Flag `--is-dry-run`: Validates capability advertisement without dispatching payload data.
  - Flag `--has-telemetry`: Emits real-time event logs to stdout.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpMultiplexSessions (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL UNIQUE,
      targetMxHost TEXT NOT NULL,
      maxStreams INTEGER NOT NULL DEFAULT 5,
      activeStreamCount INTEGER NOT NULL DEFAULT 0,
      isPipeliningSupported INTEGER NOT NULL DEFAULT 1,
      hasTlsActive INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpMultiplexSession struct {
      SessionUUID           string `json:"sessionUuid"`
      TargetMxHost          string `json:"targetMxHost"`
      MaxStreams            int    `json:"maxStreams"`
      ActiveStreamCount     int    `json:"activeStreamCount"`
      IsPipeliningSupported bool   `json:"isPipeliningSupported"`
      HasTlsActive          bool   `json:"hasTlsActive"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/pipeline/multiplexer.go` to multiplex outbox dispatch.

---

### Improvement 802: Adaptive Poisson Connection Teardown Scheduler
- **ID**: `IMP-802`
- **Cluster**: 81
- **Goal**: Computes a non-deterministic, Poisson-distributed lifetime and transaction cap for active SMTP connections to prevent destination MTAs from identifying fixed-count batching patterns.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss progress gauge displaying remaining session lifetime and probabilistic teardown threshold with amber alert styling (`#FFB800`).
  - *Web*: `<poisson-teardown-gauge>` rendering dynamic Poisson distribution curves with current session position marked by an animated SVG indicator.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline teardown-plan --target-mx=<host> [--lambda-mean=<n>] [--has-jitter] [--is-enforced]`
  - Flag `--lambda-mean`: Mean Poisson transactions before scheduled session retirement (default: 8).
  - Flag `--has-jitter`: Injects millisecond micro-delays into teardown execution.
  - Flag `--is-enforced`: Automatically executes graceful `QUIT` upon threshold arrival.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TeardownSchedules (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL UNIQUE,
      targetMxHost TEXT NOT NULL,
      transactionQuota INTEGER NOT NULL,
      transactionsCompleted INTEGER NOT NULL DEFAULT 0,
      isTeardownScheduled INTEGER NOT NULL DEFAULT 1,
      hasQuitSent INTEGER NOT NULL DEFAULT 0,
      scheduledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TeardownSchedule struct {
      SessionUUID           string `json:"sessionUuid"`
      TargetMxHost          string `json:"targetMxHost"`
      TransactionQuota      int    `json:"transactionQuota"`
      TransactionsCompleted int    `json:"transactionsCompleted"`
      IsTeardownScheduled   bool   `json:"isTeardownScheduled"`
      HasQuitSent           bool   `json:"hasQuitSent"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/pipeline/teardown.go` before next command cycle.

---

### Improvement 803: Command Inter-Arrival Jitter Engine
- **ID**: `IMP-803`
- **Cluster**: 81
- **Goal**: Inserts pseudo-random Gaussian micro-delays (15ms to 65ms) between pipelined `MAIL FROM`, `RCPT TO`, and `DATA` commands to emulate human or non-mechanized transmission intervals.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Real-time Lipgloss table displaying command micro-delay variance across recent pipeline writes with green delta metrics (`#00E676`).
  - *Web*: `<inter-arrival-jitter-chart>` component displaying interactive timeline markers of command write events.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline jitter-sync --target-mx=<host> [--min-ms=<n>] [--max-ms=<n>] [--is-adaptive]`
  - Flag `--min-ms`: Minimum millisecond delay between pipelined commands (default: 15).
  - Flag `--max-ms`: Maximum millisecond delay between pipelined commands (default: 65).
  - Flag `--is-adaptive`: Automatically scales delay based on observed remote TCP round-trip latency.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CommandJitterProfiles (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL,
      commandType TEXT NOT NULL,
      jitterDelayMs REAL NOT NULL,
      isAdaptiveDelay INTEGER NOT NULL DEFAULT 1,
      hasExecutedSuccessfully INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CommandJitterProfile struct {
      SessionUUID             string  `json:"sessionUuid"`
      CommandType             string  `json:"commandType"`
      JitterDelayMs           float64 `json:"jitterDelayMs"`
      IsAdaptiveDelay         bool    `json:"isAdaptiveDelay"`
      HasExecutedSuccessfully bool    `json:"hasExecutedSuccessfully"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Hooked into `pkg/warmup/pipeline/jitter.go` wrapping standard TCP buffer writes.

---

### Improvement 804: Half-Close Socket Detector & RST Evasion Sentinel
- **ID**: `IMP-804`
- **Cluster**: 81
- **Goal**: Proactively polls TCP socket descriptors for remote `FIN` or unexpected disconnection states prior to issuing command pipelining bursts, eliminating unhandled `SIGPIPE` errors and TCP RST packets.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Sentinel health indicator badge (`[TCP READY]` in mint green or `[PEER CLOSED: FIN]` in warning red) integrated into the outbox HUD.
  - *Web*: `<socket-sentinel-status>` badge rendering a real-time TCP socket health heartbeat icon.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline socket-health --session-id=<uuid> [--is-verbose] [--has-probe]`
  - Flag `--session-id`: Specific session UUID to probe.
  - Flag `--is-verbose`: Outputs low-level TCP socket buffer and state information.
  - Flag `--has-probe`: Transmits an RFC zero-byte TCP probe or `NOOP` check.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SocketHealthProbes (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL,
      tcpState TEXT NOT NULL,
      isReadable INTEGER NOT NULL DEFAULT 1,
      isWritable INTEGER NOT NULL DEFAULT 1,
      hasRemoteFinDetected INTEGER NOT NULL DEFAULT 0,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SocketHealthProbe struct {
      SessionUUID          string `json:"sessionUuid"`
      TcpState             string `json:"tcpState"`
      IsReadable           bool   `json:"isReadable"`
      IsWritable           bool   `json:"isWritable"`
      HasRemoteFinDetected bool   `json:"hasRemoteFinDetected"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Built into `pkg/warmup/transport/socket_sentinel.go` as a guard clause before pipeline writes.

---

### Improvement 805: RFC 2920 Advertised Capability Cache & Dynamic Fallback
- **ID**: `IMP-805`
- **Cluster**: 81
- **Goal**: Caches remote MX server `EHLO` extension keywords (`PIPELINING`, `CHUNKING`, `SIZE`, `8BITMIME`, `SMTPUTF8`) with automatic fallback to lock-step synchronous command-response mode if pipelining is omitted or rejected.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual capability checklist rendered with Lipgloss glyphs (`✓ PIPELINING`, `✓ 8BITMIME`, `✗ CHUNKING`).
  - *Web*: `<mx-capability-matrix>` displaying remote MTA features in an interactive grid.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline pipelining-cap --target-mx=<host> [--is-cached] [--has-fallback-check]`
  - Flag `--is-cached`: Reads existing capability profile from `repodb/warmup.db`.
  - Flag `--has-fallback-check`: Simulates fallback to lock-step mode.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TargetCapabilityCaches (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL UNIQUE,
      isPipeliningSupported INTEGER NOT NULL DEFAULT 1,
      isChunkingSupported INTEGER NOT NULL DEFAULT 0,
      hasSmtpUtf8Supported INTEGER NOT NULL DEFAULT 1,
      hasSizeLimitExceeded INTEGER NOT NULL DEFAULT 0,
      maxSizeBytes INTEGER NOT NULL DEFAULT 35651584,
      cachedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TargetCapabilityCache struct {
      TargetMxHost         string `json:"targetMxHost"`
      IsPipeliningSupported bool  `json:"isPipeliningSupported"`
      IsChunkingSupported   bool  `json:"isChunkingSupported"`
      HasSmtpUtf8Supported  bool  `json:"hasSmtpUtf8Supported"`
      HasSizeLimitExceeded  bool  `json:"hasSizeLimitExceeded"`
      MaxSizeBytes         int64  `json:"maxSizeBytes"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/pipeline/capability.go` guiding command serialization.

---

### Improvement 806: Dynamic TCP Write Buffer Tuner
- **ID**: `IMP-806`
- **Cluster**: 81
- **Goal**: Dynamically sizes the underlying TCP write buffer based on destination MTU and bandwidth-delay product to ensure pipelined SMTP command clusters fit within single IP packets, minimizing packet segmentation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Compact Lipgloss metric card showing active TCP buffer size (`4096 bytes`) and packet utilization percentage.
  - *Web*: `<tcp-buffer-gauge>` showing real-time MTU alignment visualization.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline buffer-tune --target-mx=<host> [--buffer-size=<bytes>] [--is-auto-tuned]`
  - Flag `--buffer-size`: Explicit socket write buffer size.
  - Flag `--is-auto-tuned`: Computes optimal buffer size based on path MTU discovery.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpBufferTuningProfiles (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      configuredBufferSize INTEGER NOT NULL DEFAULT 4096,
      pathMtuSize INTEGER NOT NULL DEFAULT 1500,
      isAutoTuned INTEGER NOT NULL DEFAULT 1,
      hasOptimalPacketFit INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpBufferTuningProfile struct {
      TargetMxHost         string `json:"targetMxHost"`
      ConfiguredBufferSize int    `json:"configuredBufferSize"`
      PathMtuSize          int    `json:"pathMtuSize"`
      IsAutoTuned          bool   `json:"isAutoTuned"`
      HasOptimalPacketFit  bool   `json:"hasOptimalPacketFit"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/transport/socket_opts.go` via OS socket control.

---

### Improvement 807: Transaction Reset Governor & Clean Slate Validator
- **ID**: `IMP-807`
- **Cluster**: 81
- **Goal**: Enforces an explicit RFC 5321 `RSET` command between pipelined transactions on reused connections, validating that the remote MTA state returns to clean envelope addressing before new recipient parameters are written.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual state transition stepper in Lipgloss (`[DATA END] -> [250 OK] -> [RSET] -> [250 OK] -> [READY]`).
  - *Web*: `<smtp-state-stepper>` showing real-time transaction lifecycle checkpoints.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline rset-governor --session-id=<uuid> [--is-enforced] [--has-audit-log]`
  - Flag `--is-enforced`: Rejects subsequent commands if `RSET` returns non-250 status code.
  - Flag `--has-audit-log`: Writes state validation events to `repodb/audit.db`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TransactionResetLogs (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL,
      isCleanStateVerified INTEGER NOT NULL DEFAULT 1,
      hasRsetAcknowledged INTEGER NOT NULL DEFAULT 1,
      resetLatencyMs REAL NOT NULL,
      executedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TransactionResetLog struct {
      SessionUUID          string  `json:"sessionUuid"`
      IsCleanStateVerified bool    `json:"isCleanStateVerified"`
      HasRsetAcknowledged  bool    `json:"hasRsetAcknowledged"`
      ResetLatencyMs       float64 `json:"resetLatencyMs"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/pipeline/reset.go` between outbox dispatches.

---

### Improvement 808: Recycled Session Egress Tokenizer
- **ID**: `IMP-808`
- **Cluster**: 81
- **Goal**: Issues high-priority token allowances to outbox queues for destinations with already established, verified-warm SMTP sessions, avoiding cold TLS setup overhead and maximizing throughput.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Token allocation table styled with cyan borders, displaying warm-socket tokens vs cold-dial tokens.
  - *Web*: `<session-token-visualizer>` rendering active token pools allocated across destination MX nodes.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline session-recycle --target-mx=<host> [--token-boost=<n>] [--is-active]`
  - Flag `--token-boost`: Additional token capacity granted for warm sockets (default: 3).
  - Flag `--is-active`: Enables opportunistic warm session routing.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RecycledSessionTokens (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      warmTokenCount INTEGER NOT NULL DEFAULT 10,
      coldTokenCount INTEGER NOT NULL DEFAULT 2,
      isOpportunisticRoutingActive INTEGER NOT NULL DEFAULT 1,
      hasActiveSessionAvailable INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RecycledSessionToken struct {
      TargetMxHost                 string `json:"targetMxHost"`
      WarmTokenCount               int    `json:"warmTokenCount"`
      ColdTokenCount               int    `json:"coldTokenCount"`
      IsOpportunisticRoutingActive bool   `json:"isOpportunisticRoutingActive"`
      HasActiveSessionAvailable    bool   `json:"hasActiveSessionAvailable"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Injected into `pkg/warmup/ratelimit/session_tokenizer.go` before dispatcher dequeue.

---

### Improvement 809: Async Pipelined TLS Handshake Pre-Warmer
- **ID**: `IMP-809`
- **Cluster**: 81
- **Goal**: Pre-establishes TLS 1.3 handshakes to impending destination MX hosts asynchronously in the background based on forecasted queue arrivals, eliminating handshake latency when messages are dequeued.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual pre-warm queue monitor displaying background connection establishment in dark cyan Lipgloss (`#00838F`).
  - *Web*: `<tls-prewarm-monitor>` showing background TLS handshake timelines and pre-warmed socket availability.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline tls-handshake --target-mx=<host> [--prewarm-count=<n>] [--has-tls13-only] [--is-background]`
  - Flag `--prewarm-count`: Number of background pre-warmed TLS sockets to maintain (default: 2).
  - Flag `--has-tls13-only`: Restricts pre-warming strictly to TLS 1.3 handshakes.
  - Flag `--is-background`: Detaches runner into background goroutine.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PrewarmedTlsSockets (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      tlsVersion TEXT NOT NULL,
      isPrewarmedReady INTEGER NOT NULL DEFAULT 1,
      hasBeenClaimed INTEGER NOT NULL DEFAULT 0,
      establishedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      expiresAt DATETIME NOT NULL
  );
  ```
  ```go
  type PrewarmedTlsSocket struct {
      TargetMxHost     string `json:"targetMxHost"`
      TlsVersion       string `json:"tlsVersion"`
      IsPrewarmedReady bool   `json:"isPrewarmedReady"`
      HasBeenClaimed   bool   `json:"hasBeenClaimed"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Executed in `pkg/warmup/transport/tls_prewarmer.go` during idle CPU cycles.

---

### Improvement 810: Graceful Pool Drain Orchestrator
- **ID**: `IMP-810`
- **Cluster**: 81
- **Goal**: Coordinates orderly drainage and closure of persistent SMTP sessions during graceful daemon shutdown or maintenance mode, ensuring all inflight transactions finish and `QUIT` is transmitted without abruptly severing sockets.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Animated Bubble Tea v2 shutdown spinner showing draining connections count down to zero with green completion checkmarks.
  - *Web*: `<pool-drain-banner>` overlay displaying real-time graceful teardown progress.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup pipeline drain-queue [--timeout-sec=<n>] [--is-force-after-timeout] [--has-audit-flush]`
  - Flag `--timeout-sec`: Maximum graceful wait period before forceful socket teardown (default: 15).
  - Flag `--is-force-after-timeout`: Sends TCP RST if remote peer fails to acknowledge `QUIT` within timeout.
  - Flag `--has-audit-flush`: Flushes all pending audit logs to `repodb/audit.db` prior to exit.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PoolDrainEvents (
      id TEXT PRIMARY KEY,
      totalSessionsToDrain INTEGER NOT NULL,
      sessionsDrainedSuccessfully INTEGER NOT NULL DEFAULT 0,
      isCleanDrainCompleted INTEGER NOT NULL DEFAULT 1,
      hasTimedOut INTEGER NOT NULL DEFAULT 0,
      drainedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PoolDrainEvent struct {
      TotalSessionsToDrain        int  `json:"totalSessionsToDrain"`
      SessionsDrainedSuccessfully int  `json:"sessionsDrainedSuccessfully"`
      IsCleanDrainCompleted       bool `json:"isCleanDrainCompleted"`
      HasTimedOut                 bool `json:"hasTimedOut"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Orchestrated in `pkg/warmup/pipeline/drain.go` handling OS termination signals.

---

# Cluster 82: Neuromorphic Feedback Calibration & Real-time ISP Throttle Inversion (811–820)

---

### Improvement 811: Neuromorphic ISP Feedback Weight Visualizer
- **ID**: `IMP-811`
- **Cluster**: 82 (Neuromorphic Feedback Calibration & Real-time ISP Throttle Inversion)
- **Goal**: Computes and displays dynamic synaptic weight vectors (0.00 to 1.00) assigned to remote Mailbox Providers based on real-time feedback signals (latency, 4xx deferrals, TLS RTT).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss colored bar chart showing synaptic weight distribution across major providers (Google, Microsoft, Yahoo, Apple) with color ramps from Mint (`#00E676`) to Red (`#FF1744`).
  - *Web*: `<neuro-weights-radar>` Web Component rendering a multi-axis SVG radar diagram of provider friction factors.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle neuro-weights [--domain=<domain>] [--is-realtime] [--has-decay]`
  - Flag `--domain`: Specific target domain to inspect (or all if omitted).
  - Flag `--is-realtime`: Launches live terminal update feed via Bubble Tea v2.
  - Flag `--has-decay`: Applies mathematical half-life decay to historic weights.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS NeuromorphicWeights (
      id TEXT PRIMARY KEY,
      providerName TEXT NOT NULL UNIQUE,
      latencyWeight REAL NOT NULL DEFAULT 0.50,
      deferralWeight REAL NOT NULL DEFAULT 0.50,
      inversionFactor REAL NOT NULL DEFAULT 1.00,
      isAdaptiveActive INTEGER NOT NULL DEFAULT 1,
      hasRecentPenalty INTEGER NOT NULL DEFAULT 0,
      calibratedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type NeuromorphicWeight struct {
      ProviderName     string  `json:"providerName"`
      LatencyWeight    float64 `json:"latencyWeight"`
      DeferralWeight   float64 `json:"deferralWeight"`
      InversionFactor  float64 `json:"inversionFactor"`
      IsAdaptiveActive bool    `json:"isAdaptiveActive"`
      HasRecentPenalty bool    `json:"hasRecentPenalty"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Calculated in `pkg/warmup/throttle/neuro_weights.go` before schedule evaluation.

---

### Improvement 812: Throttle Gradient Inverter
- **ID**: `IMP-812`
- **Cluster**: 82
- **Goal**: Automatically inverts defensive MTA throttling curves by decreasing egress rate proportionally to the gradient of 4xx deferral velocity, preventing hard blocklist triggers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual gradient slope indicator in Lipgloss (`[SLOPE: -0.42/min -> RATE INVERTED: 60% CAP]`) with amber warning accents.
  - *Web*: `<throttle-gradient-curve>` rendering interactive Bézier curves showing raw pressure vs inverted dispatch cap.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle invert --provider=<name> [--sensitivity=<float>] [--is-active] [--has-auto-recovery]`
  - Flag `--sensitivity`: Gradient responsiveness multiplier (default: 1.25).
  - Flag `--is-active`: Activates real-time rate throttling inversion.
  - Flag `--has-auto-recovery`: Gradually ramps throughput back up when deferral velocity flattens.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ThrottleInversionStates (
      id TEXT PRIMARY KEY,
      providerName TEXT NOT NULL UNIQUE,
      deferralVelocity REAL NOT NULL DEFAULT 0.0,
      appliedRateCap REAL NOT NULL DEFAULT 1.0,
      isThrottleInverted INTEGER NOT NULL DEFAULT 0,
      hasAutoRecoveryEnabled INTEGER NOT NULL DEFAULT 1,
      lastEvaluatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ThrottleInversionState struct {
      ProviderName           string  `json:"providerName"`
      DeferralVelocity       float64 `json:"deferralVelocity"`
      AppliedRateCap         float64 `json:"appliedRateCap"`
      IsThrottleInverted     bool    `json:"isThrottleInverted"`
      HasAutoRecoveryEnabled bool    `json:"hasAutoRecoveryEnabled"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Engine component in `pkg/warmup/throttle/inverter.go` moderating outbox rate limits.

---

### Improvement 813: Predictive 4xx Deferral Shifter
- **ID**: `IMP-813`
- **Cluster**: 82
- **Goal**: Detects early diagnostic text patterns in transient 421 and 451 responses (such as "rate exceeded" or "too many connections") and instantly reschedules deferred messages to alternate outbound intervals.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss table of matched deferral signatures, extracted wait instructions, and shifted dispatch timestamps.
  - *Web*: `<deferral-shifter-timeline>` component displaying deferred items reassigned across future hourly buckets.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle backpressure --target-mx=<host> [--has-pattern-match] [--is-shifted]`
  - Flag `--has-pattern-match`: Enables regex decomposition of MTA error message bodies.
  - Flag `--is-shifted`: Confirms automatic modification of outbox dispatch queue timestamps.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DeferralShiftEvents (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      rawSmtpResponse TEXT NOT NULL,
      detectedSubcode TEXT NOT NULL,
      shiftDelaySeconds INTEGER NOT NULL,
      isShiftApplied INTEGER NOT NULL DEFAULT 1,
      hasSpecificBackoffHeader INTEGER NOT NULL DEFAULT 0,
      shiftedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DeferralShiftEvent struct {
      TargetMxHost             string `json:"targetMxHost"`
      RawSmtpResponse          string `json:"rawSmtpResponse"`
      DetectedSubcode          string `json:"detectedSubcode"`
      ShiftDelaySeconds        int    `json:"shiftDelaySeconds"`
      IsShiftApplied           bool   `json:"isShiftApplied"`
      HasSpecificBackoffHeader bool   `json:"hasSpecificBackoffHeader"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/throttle/deferral_shifter.go` upon receiving 4xx responses.

---

### Improvement 814: MX Cluster Affinity Calibrator
- **ID**: `IMP-814`
- **Cluster**: 82
- **Goal**: Identifies and groups disparate MX hostnames belonging to identical autonomous systems (e.g., `smtp.google.com` and `aspmx.l.google.com`) to enforce unified neuromorphic throttle state across related endpoints.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual tree hierarchy showing cluster root nodes with aggregated throttle coefficients and member endpoints.
  - *Web*: `<mx-affinity-tree>` component presenting interactive node cluster maps.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle mx-affinity --group=<name> [--has-asn-lookup] [--is-clustered]`
  - Flag `--has-asn-lookup`: Resolves BGP ASN numbers to cluster related hostnames.
  - Flag `--is-clustered`: Groups endpoints into a single shared throttle token bucket.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MxClusterAffinities (
      id TEXT PRIMARY KEY,
      clusterGroupName TEXT NOT NULL,
      mxHostname TEXT NOT NULL UNIQUE,
      asnNumber INTEGER NOT NULL,
      isClusteredPolicyActive INTEGER NOT NULL DEFAULT 1,
      hasSharedTokens INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MxClusterAffinity struct {
      ClusterGroupName        string `json:"clusterGroupName"`
      MxHostname              string `json:"mxHostname"`
      AsnNumber               int    `json:"asnNumber"`
      IsClusteredPolicyActive bool   `json:"isClusteredPolicyActive"`
      HasSharedTokens         bool   `json:"hasSharedTokens"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Injected in `pkg/warmup/dns/affinity.go` for cross-host rate limiting.

---

### Improvement 815: Round-Trip Time Gradient Tracker
- **ID**: `IMP-815`
- **Cluster**: 82
- **Goal**: Continuously monitors the round-trip latency gradient of TCP handshakes and SMTP command ACKs to detect early-stage network or MTA congestion before explicit deferrals occur.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-density ANSI sparkline showing RTT trends over the last 100 transactions with latency delta percentages.
  - *Web*: `<rtt-gradient-sparkline>` rendering micro-charts in zero-CDN SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle rtt-gradient --target-mx=<host> [--window=<n>] [--is-streaming]`
  - Flag `--window`: Number of historic samples to include in the gradient calculation (default: 50).
  - Flag `--is-streaming`: Continuously stream RTT deltas to the console.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RttGradientSamples (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      sampleLatencyMs REAL NOT NULL,
      calculatedGradient REAL NOT NULL,
      isCongestionSuspected INTEGER NOT NULL DEFAULT 0,
      hasSustainedIncrease INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RttGradientSample struct {
      TargetMxHost          string  `json:"targetMxHost"`
      SampleLatencyMs       float64 `json:"sampleLatencyMs"`
      CalculatedGradient    float64 `json:"calculatedGradient"`
      IsCongestionSuspected bool    `json:"isCongestionSuspected"`
      HasSustainedIncrease  bool    `json:"hasSustainedIncrease"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Measured in `pkg/warmup/metrics/rtt_tracker.go` on each network frame.

---

### Improvement 816: TCP Receive Window Shrink Detector
- **ID**: `IMP-816`
- **Cluster**: 82
- **Goal**: Inspects TCP ACK packets returned by destination MTAs for shrinking TCP receive window sizes (zero-window probes or low window advertisements), indicating MTA server buffer exhaustion.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Warning badge (`[TCP WINDOW STRAIN: 2048 BYTES]`) displayed in Lipgloss amber when window size drops below 4096 bytes.
  - *Web*: `<tcp-window-meter>` showing live socket window capacity meters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle window-dampen --session-id=<uuid> [--threshold-bytes=<n>] [--is-dampened]`
  - Flag `--threshold-bytes`: Window size threshold below which dispatch is throttled (default: 4096).
  - Flag `--is-dampened`: Automatically pauses writes until window expansion is advertised.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TcpWindowEvents (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL,
      advertisedWindowBytes INTEGER NOT NULL,
      isWindowExhausted INTEGER NOT NULL DEFAULT 0,
      hasDampeningTriggered INTEGER NOT NULL DEFAULT 0,
      detectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TcpWindowEvent struct {
      SessionUUID           string `json:"sessionUuid"`
      AdvertisedWindowBytes int    `json:"advertisedWindowBytes"`
      IsWindowExhausted     bool   `json:"isWindowExhausted"`
      HasDampeningTriggered bool   `json:"hasDampeningTriggered"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/transport/window_watcher.go` wrapping raw socket reads.

---

### Improvement 817: Greeting Banner Delay Profiler
- **ID**: `IMP-817`
- **Cluster**: 82
- **Goal**: Measures the latency between initial TCP socket connection and the remote MTA's initial `220` SMTP greeting banner to profile remote greylisting defenses (teergrubing/tarpitting).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss bar chart profiling greeting banner delay across destination providers with threshold zones (<500ms normal, >2000ms tarpit warning).
  - *Web*: `<banner-delay-profiler>` displaying interactive latency distribution histograms.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle banner-pause --target-mx=<host> [--is-tarpit-detected] [--has-historical-avg]`
  - Flag `--is-tarpit-detected`: Evaluates whether current delay exceeds statistical tarpitting limits.
  - Flag `--has-historical-avg`: Compares measurement against rolling 7-day average.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS GreetingBannerDelays (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      bannerLatencyMs REAL NOT NULL,
      isTarpitDetected INTEGER NOT NULL DEFAULT 0,
      hasBannerCompliant INTEGER NOT NULL DEFAULT 1,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type GreetingBannerDelay struct {
      TargetMxHost       string  `json:"targetMxHost"`
      BannerLatencyMs    float64 `json:"bannerLatencyMs"`
      IsTarpitDetected   bool    `json:"isTarpitDetected"`
      HasBannerCompliant bool    `json:"hasBannerCompliant"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Recorded in `pkg/warmup/transport/banner.go` immediately post-connect.

---

### Improvement 818: Asymptotic Recovery Rate Tuner
- **ID**: `IMP-818`
- **Cluster**: 82
- **Goal**: Manages the gradual, non-linear restoration of outbound sending quotas following a throttle event, enforcing asymptotic curves that prevent instant re-triggering of ISP rate limits.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual ASCII recovery curve showing step-by-step restoration milestones over time.
  - *Web*: `<asymptotic-recovery-graph>` showing mathematical recovery trajectories.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle adaptive-burst --provider=<name> [--recovery-half-life-min=<n>] [--is-enforced]`
  - Flag `--recovery-half-life-min`: Time in minutes required to recover 50% of previous quota (default: 30).
  - Flag `--is-enforced`: Enforces rate limiter ceiling to asymptotic curve.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS AsymptoticRecoverySchedules (
      id TEXT PRIMARY KEY,
      providerName TEXT NOT NULL UNIQUE,
      initialSuppressionRate REAL NOT NULL,
      currentRestoredRate REAL NOT NULL,
      targetFullRate REAL NOT NULL,
      isRecoveryInProgress INTEGER NOT NULL DEFAULT 1,
      hasFullRateRestored INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type AsymptoticRecoverySchedule struct {
      ProviderName           string  `json:"providerName"`
      InitialSuppressionRate float64 `json:"initialSuppressionRate"`
      CurrentRestoredRate    float64 `json:"currentRestoredRate"`
      TargetFullRate         float64 `json:"targetFullRate"`
      IsRecoveryInProgress   bool    `json:"isRecoveryInProgress"`
      HasFullRateRestored    bool    `json:"hasFullRateRestored"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Calculated in `pkg/warmup/throttle/recovery.go` adjusting token replenishment.

---

### Improvement 819: Domain Heuristic Penalty Matrix
- **ID**: `IMP-819`
- **Cluster**: 82
- **Goal**: Maintains a cross-referenced matrix of ISP-specific penalty heuristics (e.g., maximum daily recipient caps for fresh IPs on Yahoo vs Microsoft) to prevent unforced reputation violations.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Structured Lipgloss table displaying provider rules, daily capacity utilization, and remaining headroom.
  - *Web*: `<penalty-matrix-dashboard>` rendering a status grid of compliance margins per provider.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle provider-matrix [--is-strict] [--has-overrides]`
  - Flag `--is-strict`: Enforces conservative caps over standard guidelines.
  - Flag `--has-overrides`: Displays custom operator rule overrides.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ProviderPenaltyMatrices (
      id TEXT PRIMARY KEY,
      providerName TEXT NOT NULL UNIQUE,
      maxHourlyVolume INTEGER NOT NULL,
      maxConnectionsPerIp INTEGER NOT NULL,
      isStrictEnforcement INTEGER NOT NULL DEFAULT 1,
      hasPenaltyActive INTEGER NOT NULL DEFAULT 0,
      lastRuleSyncAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ProviderPenaltyMatrix struct {
      ProviderName        string `json:"providerName"`
      MaxHourlyVolume     int    `json:"maxHourlyVolume"`
      MaxConnectionsPerIP int    `json:"maxConnectionsPerIp"`
      IsStrictEnforcement bool   `json:"isStrictEnforcement"`
      HasPenaltyActive    bool   `json:"hasPenaltyActive"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Consulted in `pkg/warmup/governance/matrix.go` during warming plan generation.

---

### Improvement 820: Dynamic Calibration Weight Reset
- **ID**: `IMP-820`
- **Cluster**: 82
- **Goal**: Provides an administrative override mechanism to safely reset accumulated neuromorphic throttle weights and recovery penalties for a provider after underlying network or IP configurations have been resolved.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Interactive confirmation dialog styled with Lipgloss warning borders (`#FF1744`), prompting operator before purging calibration states.
  - *Web*: `<calibration-reset-modal>` component requiring explicit confirmation before clearing weights.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup throttle feedback-reset --provider=<name> [--is-confirmed] [--has-audit-trail]`
  - Flag `--is-confirmed`: Skips interactive confirmation prompt.
  - Flag `--has-audit-trail`: Records reset reason and operator ID in `repodb/audit.db`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CalibrationResetEvents (
      id TEXT PRIMARY KEY,
      providerName TEXT NOT NULL,
      operatorIdentity TEXT NOT NULL,
      isResetSuccessful INTEGER NOT NULL DEFAULT 1,
      hasPreviousWeightsArchived INTEGER NOT NULL DEFAULT 1,
      resetTimestamp DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CalibrationResetEvent struct {
      ProviderName               string `json:"providerName"`
      OperatorIdentity           string `json:"operatorIdentity"`
      IsResetSuccessful          bool   `json:"isResetSuccessful"`
      HasPreviousWeightsArchived bool   `json:"hasPreviousWeightsArchived"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/throttle/reset.go` clearing runtime cache.

---

# Cluster 83: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (821–830)

---

### Improvement 821: Microsecond SMTP Flamegraph SVG Generator
- **ID**: `IMP-821`
- **Cluster**: 83 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Generates zero-external-dependency, microsecond-accurate SVG flamegraphs illustrating the exact time breakdown across DNS lookup, TCP dial, TLS handshake, EHLO, MAIL, RCPT, and DATA transmission phases.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual block breakdown table in Lipgloss showing phase durations and percentage share of total latency.
  - *Web*: `<smtp-flamegraph-hud>` custom Shadow DOM element rendering self-contained vector SVG flamegraphs with zoom and pan capabilities.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform flamegraph --transaction-id=<id> [--out-file=<path>] [--is-standalone] [--has-css]`
  - Flag `--out-file`: Destination path for generated `.svg` file.
  - Flag `--is-standalone`: Embeds inline styles for direct browser viewing.
  - Flag `--has-css`: Enables interactive SVG hover animations.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TransactionFlamegraphs (
      id TEXT PRIMARY KEY,
      transactionId TEXT NOT NULL UNIQUE,
      dnsDurationUs INTEGER NOT NULL,
      tcpConnectDurationUs INTEGER NOT NULL,
      tlsHandshakeDurationUs INTEGER NOT NULL,
      ehloDurationUs INTEGER NOT NULL,
      dataTransferDurationUs INTEGER NOT NULL,
      isFlamegraphRendered INTEGER NOT NULL DEFAULT 1,
      hasOutlierPhases INTEGER NOT NULL DEFAULT 0,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TransactionFlamegraph struct {
      TransactionID          string `json:"transactionId"`
      DnsDurationUs          int64  `json:"dnsDurationUs"`
      TcpConnectDurationUs   int64  `json:"tcpConnectDurationUs"`
      TlsHandshakeDurationUs int64  `json:"tlsHandshakeDurationUs"`
      EhloDurationUs         int64  `json:"ehloDurationUs"`
      DataTransferDurationUs int64  `json:"dataTransferDurationUs"`
      IsFlamegraphRendered   bool   `json:"isFlamegraphRendered"`
      HasOutlierPhases       bool   `json:"hasOutlierPhases"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Synthesized in `pkg/warmup/ui/flamegraph.go` after transaction completion.

---

### Improvement 822: Lipgloss TrueColor Jitter Waveform HUD
- **ID**: `IMP-822`
- **Cluster**: 83
- **Goal**: Renders a live terminal oscilloscope-style jitter waveform using Bubble Tea v2 and Lipgloss TrueColor palettes, depicting millisecond transmission variations across outgoing emails in real-time.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Continuous horizontal waveform graph drawn using block runes (` ▂▃▄▅▆▇█`) colorized from Cyan (`#00F0FF`) through Mint (`#00E676`) to Red (`#FF1744`).
  - *Web*: `<jitter-waveform-viewer>` element replicating terminal waveforms in responsive canvas/SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform live-jitter [--refresh-rate-ms=<n>] [--is-fullscreen] [--has-ansi-colors]`
  - Flag `--refresh-rate-ms`: UI update frequency in milliseconds (default: 100).
  - Flag `--is-fullscreen`: Launches full alternate terminal screen buffer.
  - Flag `--has-ansi-colors`: Emits TrueColor ANSI escape sequences.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterWaveformSamples (
      id TEXT PRIMARY KEY,
      sampleTimestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
      measuredJitterMs REAL NOT NULL,
      isWarningThresholdExceeded INTEGER NOT NULL DEFAULT 0,
      hasAnomalyDetected INTEGER NOT NULL DEFAULT 0
  );
  ```
  ```go
  type JitterWaveformSample struct {
      SampleTimestamp            string  `json:"sampleTimestamp"`
      MeasuredJitterMs           float64 `json:"measuredJitterMs"`
      IsWarningThresholdExceeded bool    `json:"isWarningThresholdExceeded"`
      HasAnomalyDetected         bool    `json:"hasAnomalyDetected"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Executed in `pkg/warmup/ui/waveform_hud.go` via Bubble Tea v2 model loops.

---

### Improvement 823: Phase-Wise Micro-Latency Breakdown
- **ID**: `IMP-823`
- **Cluster**: 83
- **Goal**: Decomposes total message delivery time into granular protocol micro-phases (Socket Dial, TLS Negotiation, EHLO, MAIL FROM, RCPT TO, DATA Transfer, QUIT) with statistical p50, p95, and p99 metrics.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss summary table displaying phase percentiles with highlighted slowest hops.
  - *Web*: `<latency-phase-breakdown>` displaying interactive phase distribution box plots.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform phase-breakdown --target-mx=<host> [--is-percentile] [--has-p99]`
  - Flag `--is-percentile`: Renders p50/p90/p95/p99 percentiles instead of simple averages.
  - Flag `--has-p99`: Flags phases where p99 diverges by >300% from p50.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PhaseLatencyStats (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      phaseName TEXT NOT NULL,
      p50LatencyUs INTEGER NOT NULL,
      p95LatencyUs INTEGER NOT NULL,
      p99LatencyUs INTEGER NOT NULL,
      isPhaseDegraded INTEGER NOT NULL DEFAULT 0,
      hasOutlierSpike INTEGER NOT NULL DEFAULT 0,
      calculatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PhaseLatencyStat struct {
      TargetMxHost    string `json:"targetMxHost"`
      PhaseName       string `json:"phaseName"`
      P50LatencyUs    int64  `json:"p50LatencyUs"`
      P95LatencyUs    int64  `json:"p95LatencyUs"`
      P99LatencyUs    int64  `json:"p99LatencyUs"`
      IsPhaseDegraded bool   `json:"isPhaseDegraded"`
      HasOutlierSpike bool   `json:"hasOutlierSpike"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Aggregated in `pkg/warmup/metrics/phase_stats.go` for developer diagnostics.

---

### Improvement 824: Zero-CDN Standalone Visualizer Packager
- **ID**: `IMP-824`
- **Cluster**: 83
- **Goal**: Bundles telemetry data and interactive SVG/CSS visualization components into a completely self-contained, air-gapped HTML file requiring zero internet connectivity or third-party CDNs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Completion message with file path and byte size confirmation in Lipgloss cyan.
  - *Web*: Generated autonomous HTML file with embedded Shadow DOM components, inline CSS, and SVG flamegraphs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform shadow-dom-view --out-file=<path> [--is-embedded] [--has-compressed-data]`
  - Flag `--out-file`: Output file destination (e.g., `dist/report.html`).
  - Flag `--is-embedded`: Embeds telemetry JSON payload directly within script tags.
  - Flag `--has-compressed-data`: Applies base64 compression to embedded datasets.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS StandaloneReportBundles (
      id TEXT PRIMARY KEY,
      bundlePath TEXT NOT NULL UNIQUE,
      bundleSizeBytes INTEGER NOT NULL,
      isAirGappedVerified INTEGER NOT NULL DEFAULT 1,
      hasZeroExternalCdn INTEGER NOT NULL DEFAULT 1,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type StandaloneReportBundle struct {
      BundlePath          string `json:"bundlePath"`
      BundleSizeBytes     int64  `json:"bundleSizeBytes"`
      IsAirGappedVerified bool   `json:"isAirGappedVerified"`
      HasZeroExternalCdn  bool   `json:"hasZeroExternalCdn"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Packaged via `pkg/warmup/ui/packager.go` for offline inspection.

---

### Improvement 825: Terminal UTF-8 Braille Sparkline Histogram
- **ID**: `IMP-825`
- **Cluster**: 83
- **Goal**: Utilizes Unicode Braille patterns (U+2800 to U+28FF) to achieve 2x horizontal and 4x vertical resolution density for in-terminal jitter and delivery latency histograms.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-resolution Braille sparklines (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`) rendered smoothly across narrow terminal windows without layout breaks.
  - *Web*: `<braille-histogram-component>` presenting terminal-identical Braille visual representations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform braille-spark --samples=<n> [--is-compact] [--has-color-ramp]`
  - Flag `--samples`: Number of recent latency samples to plot (default: 60).
  - Flag `--is-compact`: Renders single-line inline sparkline.
  - Flag `--has-color-ramp`: Colorizes Braille glyphs based on latency severity.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS BrailleSparklineSnapshots (
      id TEXT PRIMARY KEY,
      sparklineGlyphs TEXT NOT NULL,
      sampleCount INTEGER NOT NULL,
      isCompactFormat INTEGER NOT NULL DEFAULT 1,
      hasColorApplied INTEGER NOT NULL DEFAULT 1,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type BrailleSparklineSnapshot struct {
      SparklineGlyphs string `json:"sparklineGlyphs"`
      SampleCount     int    `json:"sampleCount"`
      IsCompactFormat bool   `json:"isCompactFormat"`
      HasColorApplied bool   `json:"hasColorApplied"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Formatted in `pkg/warmup/ui/braille.go` for terminal CLI command summaries.

---

### Improvement 826: Latency Tail Outlier Highlighting
- **ID**: `IMP-826`
- **Cluster**: 83
- **Goal**: Automatically identifies and visually highlights latency tail anomalies (events exceeding 3 standard deviations above the moving mean) on live TUI waveforms and flamegraphs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Distinct bright crimson highlight (`#FF1744`) on outlier graph spikes with tooltip annotations showing exact delta.
  - *Web*: `<outlier-annotation-layer>` rendering pulsing SVG alert markers over anomalous latency peaks.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform outlier-detect [--sigma=<float>] [--is-realtime] [--has-notification]`
  - Flag `--sigma`: Standard deviation multiplier for anomaly classification (default: 3.0).
  - Flag `--is-realtime`: Listens to live stream and highlights anomalies as they arrive.
  - Flag `--has-notification`: Emits terminal audio bell or alert event upon detection.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS LatencyOutlierEvents (
      id TEXT PRIMARY KEY,
      transactionId TEXT NOT NULL,
      observedLatencyMs REAL NOT NULL,
      meanLatencyMs REAL NOT NULL,
      sigmaDeviation REAL NOT NULL,
      isOutlierConfirmed INTEGER NOT NULL DEFAULT 1,
      hasAlertDispatched INTEGER NOT NULL DEFAULT 0,
      detectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type LatencyOutlierEvent struct {
      TransactionID      string  `json:"transactionId"`
      ObservedLatencyMs  float64 `json:"observedLatencyMs"`
      MeanLatencyMs      float64 `json:"meanLatencyMs"`
      SigmaDeviation     float64 `json:"sigmaDeviation"`
      IsOutlierConfirmed bool    `json:"isOutlierConfirmed"`
      HasAlertDispatched bool    `json:"hasAlertDispatched"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Evaluated in `pkg/warmup/metrics/outlier.go` during trace recording.

---

### Improvement 827: High-Definition Vector SVG Exporter
- **ID**: `IMP-827`
- **Cluster**: 83
- **Goal**: Exports cryptographic-grade, pixel-perfect vector SVG files of session flamegraphs and jitter waveforms suitable for embedding in post-mortem reports and compliance documentation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Status message displaying export path, dimensions (`1920x1080`), and file hash in Lipgloss mint green.
  - *Web*: `<svg-export-dialog>` component offering export options (dark mode, light mode, resolution scale).
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform export-svg --session-id=<uuid> --out=<path> [--width=<px>] [--height=<px>] [--is-dark-theme]`
  - Flag `--width`: SVG viewport width in pixels (default: 1920).
  - Flag `--height`: SVG viewport height in pixels (default: 1080).
  - Flag `--is-dark-theme`: Generates dark theme optimized color palette.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SvgExportRecords (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      filePath TEXT NOT NULL,
      viewBoxWidth INTEGER NOT NULL DEFAULT 1920,
      viewBoxHeight INTEGER NOT NULL DEFAULT 1080,
      isDarkTheme INTEGER NOT NULL DEFAULT 1,
      hasValidXmlHeader INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SvgExportRecord struct {
      SessionID         string `json:"sessionId"`
      FilePath          string `json:"filePath"`
      ViewBoxWidth      int    `json:"viewBoxWidth"`
      ViewBoxHeight     int    `json:"viewBoxHeight"`
      IsDarkTheme       bool   `json:"isDarkTheme"`
      HasValidXmlHeader bool   `json:"hasValidXmlHeader"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/ui/svg_export.go` for report generation.

---

### Improvement 828: Real-Time Stream Latency Differ
- **ID**: `IMP-828`
- **Cluster**: 83
- **Goal**: Compares live latency telemetry streams against historic baseline reference profiles, computing and displaying millisecond delta differentials across all protocol phases.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Two-column Lipgloss diff viewer showing baseline vs current latency with colored diff indicators (`+14ms [DEGRADED]` in red, `-8ms [FASTER]` in green).
  - *Web*: `<latency-stream-diff>` component rendering interactive comparative bar charts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform stream-diff --baseline=<id> --current=<id> [--is-normalized] [--has-phase-breakdown]`
  - Flag `--is-normalized`: Normalizes values to percentage variance.
  - Flag `--has-phase-breakdown`: Breaks down comparison by protocol sub-phases.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS StreamDiffResults (
      id TEXT PRIMARY KEY,
      baselineProfileId TEXT NOT NULL,
      currentProfileId TEXT NOT NULL,
      totalDeltaMs REAL NOT NULL,
      isPerformanceImproved INTEGER NOT NULL DEFAULT 0,
      hasSignificantDegradation INTEGER NOT NULL DEFAULT 0,
      comparedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type StreamDiffResult struct {
      BaselineProfileID         string  `json:"baselineProfileId"`
      CurrentProfileID          string  `json:"currentProfileId"`
      TotalDeltaMs              float64 `json:"totalDeltaMs"`
      IsPerformanceImproved     bool    `json:"isPerformanceImproved"`
      HasSignificantDegradation bool    `json:"hasSignificantDegradation"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Executed in `pkg/warmup/metrics/differ.go` during performance regressions checks.

---

### Improvement 829: Fullscreen Bubble Tea v2 HUD
- **ID**: `IMP-829`
- **Cluster**: 83
- **Goal**: Delivers an immersive, responsive, full-terminal dashboard utilizing Bubble Tea v2 and Lipgloss layout containers, combining active connection pools, throttle gauges, and live waveforms into a unified console.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Multi-pane terminal cockpit (Outbox Pool, Throttle Gauges, Live Jitter Waveform, System Health) with dynamic keyboard navigation (`tab` switch, `q` quit).
  - *Web*: `<terminal-hud-mirror>` web view replicating full terminal layout in modern responsive CSS grid.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform ansi-hud [--is-interactive] [--has-keybinds] [--theme=<name>]`
  - Flag `--is-interactive`: Enables mouse and keyboard interaction within terminal UI.
  - Flag `--has-keybinds`: Displays shortcut hints footer (`[Tab] Next [P] Pause [Q] Quit`).
  - Flag `--theme`: Color palette selection (`cyberpunk`, `monokai`, `matrix`).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TerminalHudSessions (
      id TEXT PRIMARY KEY,
      terminalWidth INTEGER NOT NULL,
      terminalHeight INTEGER NOT NULL,
      activeTheme TEXT NOT NULL DEFAULT 'cyberpunk',
      isInteractive INTEGER NOT NULL DEFAULT 1,
      hasMouseTracking INTEGER NOT NULL DEFAULT 1,
      startedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TerminalHudSession struct {
      TerminalWidth    int    `json:"terminalWidth"`
      TerminalHeight   int    `json:"terminalHeight"`
      ActiveTheme      string `json:"activeTheme"`
      IsInteractive    bool   `json:"isInteractive"`
      HasMouseTracking bool   `json:"hasMouseTracking"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Primary operator interface in `pkg/warmup/ui/hud/app.go`.

---

### Improvement 830: Zero-Overhead Render Benchmark
- **ID**: `IMP-830`
- **Cluster**: 83
- **Goal**: Benchmarks and profiles the rendering pipeline of both the Lipgloss ANSI engine and the Shadow DOM SVG generator to verify sub-millisecond render overhead and zero allocation bloat on high-throughput nodes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Benchmark execution scorecard showing frames-per-second, memory allocation per render, and p99 render durations.
  - *Web*: `<render-benchmark-scorecard>` component displaying performance profiling results.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup waveform render-benchmark [--iterations=<n>] [--is-profiling-cpu] [--has-alloc-stats]`
  - Flag `--iterations`: Number of render loops to execute (default: 1000).
  - Flag `--is-profiling-cpu`: Outputs CPU profile to `repodb/benchmark.prof`.
  - Flag `--has-alloc-stats`: Measures heap allocations per frame.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RenderBenchmarkResults (
      id TEXT PRIMARY KEY,
      iterationCount INTEGER NOT NULL,
      meanRenderDurationUs INTEGER NOT NULL,
      allocationsPerOp INTEGER NOT NULL,
      bytesPerOp INTEGER NOT NULL,
      isSubMillisecondVerified INTEGER NOT NULL DEFAULT 1,
      hasZeroAllocationVerified INTEGER NOT NULL DEFAULT 0,
      testedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RenderBenchmarkResult struct {
      IterationCount            int   `json:"iterationCount"`
      MeanRenderDurationUs      int64 `json:"meanRenderDurationUs"`
      AllocationsPerOp          int   `json:"allocationsPerOp"`
      BytesPerOp                int   `json:"bytesPerOp"`
      IsSubMillisecondVerified  bool  `json:"isSubMillisecondVerified"`
      HasZeroAllocationVerified bool  `json:"hasZeroAllocationVerified"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Verified in `pkg/warmup/ui/benchmark_test.go` protecting core runtime performance.
