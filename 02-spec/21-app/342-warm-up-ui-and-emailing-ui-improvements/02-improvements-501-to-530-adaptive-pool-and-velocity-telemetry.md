# Specification 342: Improvements 501 to 530 (Adaptive Pool & Velocity Telemetry)

- **Sequence:** 342
- **Slug:** `342-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/02-improvements-501-to-530-adaptive-pool-and-velocity-telemetry.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author 1

---

## Improvement Catalog Summary (501 - 530)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 51** | Items 501–510 | Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix |
| **Cluster 52** | Items 511–520 | Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine |
| **Cluster 53** | Items 521–530 | Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs |

---

# Cluster 51: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (Items 501–510)

---

### Improvement 501: Multi-ISP Persistent Socket Matrix Visualizer
- **ID:** `WARMUP-POOL-501`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Render a high-density, real-time matrix of outbound TCP sockets grouped by ISP (Google, Microsoft, Yahoo, Apple, Generic) with sub-second health status indicators.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-socket-matrix>`
  - **Lipgloss Terminal Palette:** Background `#0b0f19`, Border `#1e293b`, Header `#06b6d4` (Cyan), Socket Active `#10b981` (Emerald), Socket Idle `#0ea5e9` (Sky), Socket Degraded `#f59e0b` (Amber), Socket Closed `#f43f5e` (Rose).
  - **Shadow DOM Layout:** Encapsulated shadow root (`mode: 'open'`). Top toolbar displays aggregate socket counters (`Active: 24`, `Idle: 18`, `Max: 60`). The body renders a responsive grid of ISP pool cards. Each card displays an ASCII/Unicode connection slot array: `[●●●○··]`, accompanied by socket latency metrics in milliseconds.
  - **Zero-CDN Invariant:** Native SVG and CSS Grid rendering. Zero external libraries. Fully interactive keyboard navigation with `aria-keyshortcuts="ArrowRight ArrowLeft"`.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool matrix --domain "outbound.gitmap.dev" --watch --interval 1000ms
  gitmap email warmup pool matrix inspect --pool-id "pool-google-01" --format json
  ```
  - `--domain`: FQDN of the outbound sending domain (string).
  - `--watch`: Continuous terminal refresh flag (bool).
  - `--interval`: Refresh polling interval (duration, default `1000ms`).
  - `--pool-id`: Specific ISP pool identifier (string).
  - `--format`: Output format (`text`, `json`, `table`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type SocketMatrixState struct {
      PoolID            string  `json:"poolId"`
      DomainName        string  `json:"domainName"`
      IspClass          string  `json:"ispClass"`
      ActiveSockets     int     `json:"activeSockets"`
      IdleSockets       int     `json:"idleSockets"`
      MaxAllowedSockets int     `json:"maxAllowedSockets"`
      AverageLatencyMs  float64 `json:"averageLatencyMs"`
      IsActive          bool    `json:"isActive"`
      HasDegradedSocket bool    `json:"hasDegradedSocket"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS SocketMatrixState (
      poolId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      ispClass TEXT NOT NULL,
      activeSockets INTEGER NOT NULL DEFAULT 0,
      idleSockets INTEGER NOT NULL DEFAULT 0,
      maxAllowedSockets INTEGER NOT NULL DEFAULT 10,
      averageLatencyMs REAL NOT NULL DEFAULT 0.0,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasDegradedSocket INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Subscribed to `pkg/warmup/pool` event stream; dispatches updates over bus topic `warmup.pool.matrix.state`.

---

### Improvement 502: Automated SMTP Keep-Alive & Graceful QUIT Pipeline
- **ID:** `WARMUP-POOL-502`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Prevent abrupt TCP RST disconnects by scheduling automated NOOP probes and graceful QUIT terminations based on ISP-specific idle deadlines.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-keepalive-monitor>`
  - **Lipgloss Terminal Palette:** Border `#1e293b`, Text `#94a3b8`, Probe Success `#10b981`, Graceful Quit `#8b5cf6` (Purple).
  - **Shadow DOM Layout:** A real-time timeline stream depicting upcoming socket expiration counts with countdown progress bars. Badges indicate graceful terminations (`QUIT SENT`) vs proactive health probes (`NOOP 250`).
  - **Zero-CDN Invariant:** Self-contained CSS keyframe animations for countdown bars.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool keepalive status --domain "outbound.gitmap.dev"
  gitmap email warmup pool keepalive tune --domain "outbound.gitmap.dev" --isp "google" --timeout 25s
  ```
  - `--domain`: Outbound domain (string).
  - `--isp`: Target mailbox provider class (string).
  - `--timeout`: Maximum idle socket retention time before graceful QUIT (duration).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type KeepAliveConfig struct {
      ConfigID            string `json:"configId"`
      DomainName          string `json:"domainName"`
      IspClass            string `json:"ispClass"`
      KeepAliveTimeoutSec int    `json:"keepAliveTimeoutSec"`
      IsProbingEnabled    bool   `json:"isProbingEnabled"`
      HasGracefulShutdown bool   `json:"hasGracefulShutdown"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS KeepAliveConfig (
      configId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      ispClass TEXT NOT NULL,
      keepAliveTimeoutSec INTEGER NOT NULL DEFAULT 30,
      isProbingEnabled INTEGER NOT NULL DEFAULT 1,
      hasGracefulShutdown INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Integrates with `pkg/warmup/smtp/pool_worker.go`; logs termination events into `audit.db -> ConnectionAuditTrail`.

---

### Improvement 503: TLS 1.3 Session Ticket Reuse & 0-RTT Cache Manager
- **ID:** `WARMUP-POOL-503`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Cache and reuse TLS 1.3 session tickets for outbound SMTP connections to reduce handshake roundtrips and eliminate cryptographic overhead.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-tls-session-cache>`
  - **Lipgloss Terminal Palette:** Background `#0b0f19`, Hit Accent `#10b981` (Emerald), Miss Accent `#f59e0b` (Amber), Cipher Label `#06b6d4`.
  - **Shadow DOM Layout:** Metric cards highlighting `Ticket Hit Ratio: 92.4%`, `Handshake Time Saved: 48ms`, and an active cipher suite indicator pill (`TLS_AES_256_GCM_SHA384`).
  - **Zero-CDN Invariant:** Native SVG donut chart showing hit vs miss ratio.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool tls-cache stats --domain "outbound.gitmap.dev"
  gitmap email warmup pool tls-cache purge --domain "outbound.gitmap.dev" --expired-only
  ```
  - `--domain`: Sending domain (string).
  - `--expired-only`: Purge only invalid/expired session tickets (bool).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type TlsSessionTicketRecord struct {
      TicketID         string `json:"ticketId"`
      DomainName       string `json:"domainName"`
      RemoteMX         string `json:"remoteMx"`
      CipherSuite      string `json:"cipherSuite"`
      TicketExpiresAt  int64  `json:"ticketExpiresAt"`
      IsActive         bool   `json:"isActive"`
      HasEarlyDataUsed bool   `json:"hasEarlyDataUsed"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS TlsSessionTicketRecord (
      ticketId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      remoteMx TEXT NOT NULL,
      cipherSuite TEXT NOT NULL,
      ticketExpiresAt INTEGER NOT NULL,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasEarlyDataUsed INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Hooked to `crypto/tls.ClientSessionCache` implementation inside `pkg/warmup/smtp/tls_session.go`.

---

### Improvement 504: Socket Degradation & High-RTT Auto-Drain Circuit
- **ID:** `WARMUP-POOL-504`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Autonomously identify sockets with RTT latency > 1,200ms or packet retransmission spikes, drain them from the active queue, and replace them transparently.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-socket-drain-circuit>`
  - **Lipgloss Terminal Palette:** Degraded Alert `#f59e0b`, Draining `#6366f1` (Indigo), Purged `#94a3b8`.
  - **Shadow DOM Layout:** Table listing socket IDs, current RTT latency, packet loss percentage, and drain state (`DRAINING`, `DRAINED`, `HEALTHY`).
  - **Zero-CDN Invariant:** Pure vanilla DOM elements with CSS transitions.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool drain-circuit status --domain "outbound.gitmap.dev"
  gitmap email warmup pool drain-circuit trip --socket-id "sock-9912" --reason "high-rtt"
  ```
  - `--domain`: Sending domain (string).
  - `--socket-id`: Individual socket identifier (string).
  - `--reason`: Trigger reason string (`high-rtt`, `packet-loss`, `manual`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type SocketDrainEvent struct {
      EventID          string  `json:"eventId"`
      SocketID         string  `json:"socketId"`
      ObservedLatency  float64 `json:"observedLatency"`
      ThresholdLatency float64 `json:"thresholdLatency"`
      IsDrained        bool    `json:"isDrained"`
      HasReplaced      bool    `json:"hasReplaced"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS SocketDrainEvent (
      eventId TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      observedLatency REAL NOT NULL,
      thresholdLatency REAL NOT NULL DEFAULT 1200.0,
      isDrained INTEGER NOT NULL DEFAULT 1,
      hasReplaced INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Invoked by `pkg/warmup/pool/circuit_breaker.go`; emits incident to `audit.db -> ConnectionAuditTrail`.

---

### Improvement 505: Dynamic MX Priority & Failover Route Mapper
- **ID:** `WARMUP-POOL-505`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Dynamically map and failover outbound connections across MX records according to DNS preference weight and observed delivery health.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-mx-route-map>`
  - **Lipgloss Terminal Palette:** Primary Route `#06b6d4`, Secondary Route `#3b82f6`, Degraded Route `#ef4444`.
  - **Shadow DOM Layout:** Visual topological graph rendered in inline SVG showing the primary MX route and fallback paths with connection weight indicators.
  - **Zero-CDN Invariant:** Handcrafted SVG topology layout; zero external graph engines.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool mx routes --domain "outbound.gitmap.dev" --target-domain "gmail.com"
  gitmap email warmup pool mx failover --target-domain "gmail.com" --preferred-mx "alt1.gmail-smtp-in.l.google.com"
  ```
  - `--domain`: Outbound source domain (string).
  - `--target-domain`: Recipient domain (string).
  - `--preferred-mx`: MX host override (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type MxRoutePreference struct {
      RouteID        string  `json:"routeId"`
      TargetDomain   string  `json:"targetDomain"`
      HostName       string  `json:"hostName"`
      Priority       int     `json:"priority"`
      SuccessRate    float64 `json:"successRate"`
      IsActive       bool    `json:"isActive"`
      HasPrimaryRole bool    `json:"hasPrimaryRole"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS MxRoutePreference (
      routeId TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      hostName TEXT NOT NULL,
      priority INTEGER NOT NULL,
      successRate REAL NOT NULL DEFAULT 1.0,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasPrimaryRole INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Connected to `pkg/warmup/dns/mx_router.go`.

---

### Improvement 506: Connection Pool Concurrency Governor
- **ID:** `WARMUP-POOL-506`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Throttle maximum concurrent TCP connections per destination ISP to prevent temporary IP reputation penalties caused by aggressive connection bursts.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-concurrency-governor>`
  - **Lipgloss Terminal Palette:** Limit Gauge Background `#1e293b`, Usage Fill `#06b6d4`, Warning `#f59e0b`.
  - **Shadow DOM Layout:** Radial meter displaying current concurrent connections vs ISP limit threshold (`Current: 8 / Max: 12`).
  - **Zero-CDN Invariant:** SVG stroke-dasharray radial meter calculation.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool governor status --domain "outbound.gitmap.dev"
  gitmap email warmup pool governor set --domain "outbound.gitmap.dev" --isp "microsoft" --max-concurrency 8
  ```
  - `--domain`: Sending domain (string).
  - `--isp`: Mailbox provider identifier (string).
  - `--max-concurrency`: Maximum simultaneous active sockets (int).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type ConcurrencyLimitConfig struct {
      GovernorID        string `json:"governorId"`
      DomainName        string `json:"domainName"`
      IspClass          string `json:"ispClass"`
      MaxConcurrency    int    `json:"maxConcurrency"`
      CurrentActive     int    `json:"currentActive"`
      IsThrottled       bool   `json:"isThrottled"`
      HasCapacitySpill  bool   `json:"hasCapacitySpill"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS ConcurrencyLimitConfig (
      governorId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      ispClass TEXT NOT NULL,
      maxConcurrency INTEGER NOT NULL DEFAULT 10,
      currentActive INTEGER NOT NULL DEFAULT 0,
      isThrottled INTEGER NOT NULL DEFAULT 0,
      hasCapacitySpill INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Enforces channel capacity locks inside `pkg/warmup/pool/semaphore.go`.

---

### Improvement 507: Asynchronous Socket Health Verification Daemon
- **ID:** `WARMUP-POOL-507`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Run non-blocking background probes validating that idle sockets in the pool respond to EHLO and RSET commands without hanging.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-health-daemon-hud>`
  - **Lipgloss Terminal Palette:** Daemon Active `#10b981`, Pulse Indicator `#06b6d4`, Failure Red `#ef4444`.
  - **Shadow DOM Layout:** Status pill indicating `Daemon: RUNNING`, alongside an inspection log feed updating with sub-second probe responses.
  - **Zero-CDN Invariant:** Self-rendering DOM log buffer (capped at 50 records) with automated scroll pinning.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool health-daemon run --domain "outbound.gitmap.dev" --interval 5s
  gitmap email warmup pool health-daemon probe --socket-id "sock-001"
  ```
  - `--domain`: Outbound domain (string).
  - `--interval`: Inspection frequency (duration).
  - `--socket-id`: Individual socket to probe immediately (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type HealthDaemonProbeResult struct {
      ProbeID          string  `json:"probeId"`
      SocketID         string  `json:"socketId"`
      CommandIssued    string  `json:"commandIssued"`
      ResponseCode     int     `json:"responseCode"`
      DurationMs       float64 `json:"durationMs"`
      IsHealthy        bool    `json:"isHealthy"`
      HasTimedOut      bool    `json:"hasTimedOut"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS HealthDaemonProbeResult (
      probeId TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      commandIssued TEXT NOT NULL,
      responseCode INTEGER NOT NULL,
      durationMs REAL NOT NULL,
      isHealthy INTEGER NOT NULL DEFAULT 1,
      hasTimedOut INTEGER NOT NULL DEFAULT 0,
      testedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Daemon runner in `pkg/warmup/pool/health_daemon.go`; logs probe telemetry into `repodb/warmup.db`.

---

### Improvement 508: TCP Window Scale & Buffer Autotuner
- **ID:** `WARMUP-POOL-508`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Dynamically calibrate socket receive and send buffer sizes (`SO_RCVBUF`, `SO_SNDBUF`) based on outbound email payload sizes and bandwidth-delay product.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-buffer-tuner>`
  - **Lipgloss Terminal Palette:** Primary `#06b6d4`, Secondary `#94a3b8`, Highlight `#10b981`.
  - **Shadow DOM Layout:** Real-time buffer allocation gauges displaying active buffer size, BDP estimation, and packet overhead stats.
  - **Zero-CDN Invariant:** CSS linear gradients and flex layout for bar visualizations.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool buffer-tune --domain "outbound.gitmap.dev" --isp "google"
  gitmap email warmup pool buffer-tune set --send-buffer 65536 --recv-buffer 65536
  ```
  - `--domain`: Target domain (string).
  - `--isp`: Mailbox provider class (string).
  - `--send-buffer`: Outbound buffer size in bytes (int).
  - `--recv-buffer`: Inbound buffer size in bytes (int).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type SocketBufferConfig struct {
      BufferID         string `json:"bufferId"`
      DomainName       string `json:"domainName"`
      SendBufferSize   int    `json:"sendBufferSize"`
      RecvBufferSize   int    `json:"recvBufferSize"`
      IsAutotuned      bool   `json:"isAutotuned"`
      HasWindowScaled  bool   `json:"hasWindowScaled"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS SocketBufferConfig (
      bufferId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      sendBufferSize INTEGER NOT NULL DEFAULT 65536,
      recvBufferSize INTEGER NOT NULL DEFAULT 65536,
      isAutotuned INTEGER NOT NULL DEFAULT 1,
      hasWindowScaled INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Low-level socket options applied via `syscall` in `pkg/warmup/smtp/socket_dialer.go`.

---

### Improvement 509: Multi-Homed Outbound Egress IP Balancer
- **ID:** `WARMUP-POOL-509`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Bind outbound TCP connections across multiple local network interfaces and IP addresses according to assigned warmup quotas.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-egress-balancer>`
  - **Lipgloss Terminal Palette:** IP Badge Active `#10b981`, IP Badge Standby `#64748b`, Usage Accent `#06b6d4`.
  - **Shadow DOM Layout:** Multi-card layout presenting available network interfaces (`eth0`, `eth1`), their bound public IPs, and real-time dispatch share percentages.
  - **Zero-CDN Invariant:** Native SVG distribution bar chart with accessible tooltips.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool egress list --domain "outbound.gitmap.dev"
  gitmap email warmup pool egress add --ip "198.51.100.24" --interface "eth0" --weight 50
  ```
  - `--domain`: Sending domain (string).
  - `--ip`: Outbound egress IP address (string).
  - `--interface`: System network interface name (string).
  - `--weight`: Traffic distribution weight (int, 1-100).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type EgressIpAssignment struct {
      AssignmentID    string `json:"assignmentId"`
      DomainName      string `json:"domainName"`
      PublicIP        string `json:"publicIp"`
      InterfaceName   string `json:"interfaceName"`
      TrafficWeight   int    `json:"trafficWeight"`
      IsActive        bool   `json:"isActive"`
      HasDirectRoute  bool   `json:"hasDirectRoute"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS EgressIpAssignment (
      assignmentId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      publicIp TEXT NOT NULL,
      interfaceName TEXT NOT NULL,
      trafficWeight INTEGER NOT NULL DEFAULT 50,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasDirectRoute INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Dialers use `net.Dialer.LocalAddr` configured in `pkg/warmup/smtp/dialer_pool.go`.

---

### Improvement 510: Connection Pool Historical Audit Logger
- **ID:** `WARMUP-POOL-510`
- **Cluster:** 51 (Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix)
- **Goal:** Immutably persist every connection lifecycle event (creation, handshake, verification, eviction, error) to `repodb/audit.db` for forensic compliance.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-connection-audit-log>`
  - **Lipgloss Terminal Palette:** Monospace Font, Success Green `#10b981`, Warning Yellow `#f59e0b`, Error Red `#f43f5e`, Timestamp Gray `#64748b`.
  - **Shadow DOM Layout:** Terminal log scroll pane displaying structured audit records: `[2026-10-27T14:02:11Z] [CONNECT] pool=google-01 ip=142.250.150.26 latency=14.2ms status=OK`.
  - **Zero-CDN Invariant:** HTML5 `<pre>` and `<code>` container with hardware-accelerated smooth scrolling.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup pool audit tail --domain "outbound.gitmap.dev" --limit 50
  gitmap email warmup pool audit export --format csv --output ./audit-export.csv
  ```
  - `--domain`: Sending domain (string).
  - `--limit`: Number of historical entries to display (int).
  - `--format`: Export format (`json`, `csv`, `text`).
  - `--output`: File destination path (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type ConnectionAuditRecord struct {
      AuditID              string  `json:"auditId"`
      SocketID             string  `json:"socketId"`
      PoolID               string  `json:"poolId"`
      EventType            string  `json:"eventType"`
      LatencyMs            float64 `json:"latencyMs"`
      IsSuccess            bool    `json:"isSuccess"`
      HasSecurityViolation bool    `json:"hasSecurityViolation"`
  }
  ```
  ```sql
  -- Storage in repodb/audit.db
  CREATE TABLE IF NOT EXISTS ConnectionAuditRecord (
      auditId TEXT PRIMARY KEY,
      socketId TEXT NOT NULL,
      poolId TEXT NOT NULL,
      eventType TEXT NOT NULL,
      latencyMs REAL NOT NULL,
      isSuccess INTEGER NOT NULL DEFAULT 1,
      hasSecurityViolation INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Asynchronous write worker in `pkg/warmup/audit/writer.go` appending directly to `repodb/audit.db`.

---

# Cluster 52: Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine (Items 511–520)

---

### Improvement 511: Neuromorphic Membrane Potential Visualizer
- **ID:** `WARMUP-VELOCITY-511`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Visualize the continuous membrane potential of the leaky integrate-and-fire send-cadence engine, showing how positive delivery signals drive volume spikes.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-membrane-potential>`
  - **Lipgloss Terminal Palette:** Background `#0b0f19`, Trace Cyan `#06b6d4`, Threshold Amber `#f59e0b`, Spike Emerald `#10b981`.
  - **Shadow DOM Layout:** Real-time oscilloscope-style SVG line graph displaying current potential `V(t)` alongside threshold `V_th` and dynamic leak decay.
  - **Zero-CDN Invariant:** Native SVG vector path constructed via cubic spline interpolation in pure vanilla JS.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity potential --domain "outbound.gitmap.dev" --watch
  gitmap email warmup velocity potential reset --domain "outbound.gitmap.dev"
  ```
  - `--domain`: Sending domain (string).
  - `--watch`: Live update stream in terminal (bool).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type MembranePotentialTelemetry struct {
      RegulatorID       string  `json:"regulatorId"`
      DomainName        string  `json:"domainName"`
      CurrentPotential  float64 `json:"currentPotential"`
      FiringThreshold   float64 `json:"firingThreshold"`
      LeakCoefficient   float64 `json:"leakCoefficient"`
      IsRefractory      bool    `json:"isRefractory"`
      HasActionFired    bool    `json:"hasActionFired"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS MembranePotentialTelemetry (
      regulatorId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      currentPotential REAL NOT NULL DEFAULT 0.0,
      firingThreshold REAL NOT NULL DEFAULT 1.0,
      leakCoefficient REAL NOT NULL DEFAULT 0.05,
      isRefractory INTEGER NOT NULL DEFAULT 0,
      hasActionFired INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Emits from `pkg/warmup/neuromorphic/membrane.go` to the main event bus.

---

### Improvement 512: Poisson Inter-Arrival Dispatch Scheduler
- **ID:** `WARMUP-VELOCITY-512`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Schedule outbound warmup emails with non-uniform, Poisson-distributed intervals mimicking human communication bursts.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-poisson-scheduler>`
  - **Lipgloss Terminal Palette:** Primary `#06b6d4`, Interval Bubble `#3b82f6`, Mean Cadence Pill `#10b981`.
  - **Shadow DOM Layout:** Distribution histogram showing inter-arrival time frequency bins, highlighting current lambda parameter (`λ = 4.2 emails/hour`).
  - **Zero-CDN Invariant:** Pure CSS flex-box histogram bars rendered without external graphing tools.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity poisson status --domain "outbound.gitmap.dev"
  gitmap email warmup velocity poisson set-lambda --domain "outbound.gitmap.dev" --lambda 5.0
  ```
  - `--domain`: Sending domain (string).
  - `--lambda`: Mean arrival rate parameter (float).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type PoissonScheduleConfig struct {
      ScheduleID       string  `json:"scheduleId"`
      DomainName       string  `json:"domainName"`
      LambdaRate       float64 `json:"lambdaRate"`
      MinIntervalSec   int     `json:"minIntervalSec"`
      MaxIntervalSec   int     `json:"maxIntervalSec"`
      IsActive         bool    `json:"isActive"`
      HasJitterApplied bool    `json:"hasJitterApplied"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS PoissonScheduleConfig (
      scheduleId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      lambdaRate REAL NOT NULL DEFAULT 4.0,
      minIntervalSec INTEGER NOT NULL DEFAULT 30,
      maxIntervalSec INTEGER NOT NULL DEFAULT 900,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasJitterApplied INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Driver loop in `pkg/warmup/scheduler/poisson_timer.go`.

---

### Improvement 513: Dynamic Refractory Period Circuit
- **ID:** `WARMUP-VELOCITY-513`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Enforce an unalterable refractory pause following any volume surge or minor ISP deferral, preventing catastrophic reputation cascades.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-refractory-circuit>`
  - **Lipgloss Terminal Palette:** Idle State `#10b981`, Active Refractory `#f59e0b`, Critical Lockout `#ef4444`.
  - **Shadow DOM Layout:** Countdown circular meter showing remaining refractory duration (`Hold Remaining: 14m 22s`), with disabled dispatch controls.
  - **Zero-CDN Invariant:** SVG stroke-dashoffset countdown animation.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity refractory status --domain "outbound.gitmap.dev"
  gitmap email warmup velocity refractory trip --domain "outbound.gitmap.dev" --duration 15m
  ```
  - `--domain`: Sending domain (string).
  - `--duration`: Refractory hold time (duration, e.g., `15m`, `30m`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type RefractoryState struct {
      DomainName       string `json:"domainName"`
      HoldExpiresAt    int64  `json:"holdExpiresAt"`
      TriggerReason    string `json:"triggerReason"`
      IsHolding        bool   `json:"isHolding"`
      HasForcedRelease bool   `json:"hasForcedRelease"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS RefractoryState (
      domainName TEXT PRIMARY KEY,
      holdExpiresAt INTEGER NOT NULL,
      triggerReason TEXT NOT NULL,
      isHolding INTEGER NOT NULL DEFAULT 0,
      hasForcedRelease INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Gatekeeper in `pkg/warmup/neuromorphic/refractory_gate.go`.

---

### Improvement 514: Circadian Rhythm Volume Modulator
- **ID:** `WARMUP-VELOCITY-514`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Modulate outbound sending volume according to recipient timezone circadian curves, suppressing night-time dispatches and concentrating natural daytime flows.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-circadian-modulator>`
  - **Lipgloss Terminal Palette:** Daytime Sun Yellow `#f59e0b`, Twilight Indigo `#6366f1`, Night Sky Deep Blue `#1e1b4b`.
  - **Shadow DOM Layout:** 24-hour visual dial indicating recipient daytime peaks (09:00 - 17:00 local time) and current sending phase factor (`0.85x`).
  - **Zero-CDN Invariant:** Native SVG sine wave rendering representing solar progression.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity circadian status --domain "outbound.gitmap.dev" --timezone "America/New_York"
  gitmap email warmup velocity circadian configure --domain "outbound.gitmap.dev" --night-suppression
  ```
  - `--domain`: Sending domain (string).
  - `--timezone`: Target recipient timezone (string).
  - `--night-suppression`: Flag to strictly suspend non-essential warmup between 22:00 and 06:00 (bool).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type CircadianModulationConfig struct {
      ConfigID             string  `json:"configId"`
      DomainName           string  `json:"domainName"`
      TargetTimezone       string  `json:"targetTimezone"`
      NightReductionFactor float64 `json:"nightReductionFactor"`
      IsSuppressedAtNight  bool    `json:"isSuppressedAtNight"`
      HasPeakWindowEnabled bool    `json:"hasPeakWindowEnabled"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS CircadianModulationConfig (
      configId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      targetTimezone TEXT NOT NULL DEFAULT 'UTC',
      nightReductionFactor REAL NOT NULL DEFAULT 0.1,
      isSuppressedAtNight INTEGER NOT NULL DEFAULT 1,
      hasPeakWindowEnabled INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Applied as dynamic scaling multiplier in `pkg/warmup/neuromorphic/circadian.go`.

---

### Improvement 515: ISP Feedback Damping Attenuator
- **ID:** `WARMUP-VELOCITY-515`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Automatically depress send velocity by an exponential damping factor whenever SMTP 421 or 451 temporary deferral codes are returned.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-damping-attenuator>`
  - **Lipgloss Terminal Palette:** Attenuation High `#ef4444`, Moderate `#f59e0b`, Normal `#10b981`.
  - **Shadow DOM Layout:** Damping coefficient visual slider showing current multiplier (`γ = 0.42`), recent deferral event counts, and recovery projection.
  - **Zero-CDN Invariant:** Plain HTML5 slider bar rendered with CSS custom properties.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity damping status --domain "outbound.gitmap.dev"
  gitmap email warmup velocity damping reset --domain "outbound.gitmap.dev" --isp "microsoft"
  ```
  - `--domain`: Sending domain (string).
  - `--isp`: Mailbox provider identifier (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type FeedbackDampingState struct {
      DomainName          string  `json:"domainName"`
      IspClass            string  `json:"ispClass"`
      DampingFactor       float64 `json:"dampingFactor"`
      ConsecutiveFailures int     `json:"consecutiveFailures"`
      IsDamped            bool    `json:"isDamped"`
      HasRecovered        bool    `json:"hasRecovered"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS FeedbackDampingState (
      domainName TEXT NOT NULL,
      ispClass TEXT NOT NULL,
      dampingFactor REAL NOT NULL DEFAULT 1.0,
      consecutiveFailures INTEGER NOT NULL DEFAULT 0,
      isDamped INTEGER NOT NULL DEFAULT 0,
      hasRecovered INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY(domainName, ispClass)
  );
  ```
- **Antigravity / GitMap Integration Point:** Hooked to SMTP response parser in `pkg/warmup/neuromorphic/attenuator.go`.

---

### Improvement 516: Seed Engagement Velocity Catalyst
- **ID:** `WARMUP-VELOCITY-516`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Accelerate allowable warmup velocity when positive seed signals (inbox read, reply, star, mark important) are verified across seed networks.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-velocity-catalyst>`
  - **Lipgloss Terminal Palette:** Catalyst Green `#10b981`, Acceleration Blue `#0ea5e9`, Baseline Gray `#64748b`.
  - **Shadow DOM Layout:** Momentum tachometer gauge displaying current boost factor (`+1.35x Boost`), driven by verified positive engagement ratios.
  - **Zero-CDN Invariant:** Native SVG gauge needle with CSS transform rotation.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity catalyst metrics --domain "outbound.gitmap.dev"
  gitmap email warmup velocity catalyst record-signal --type "thread-reply" --domain "outbound.gitmap.dev"
  ```
  - `--domain`: Sending domain (string).
  - `--type`: Signal type (`seed-read`, `thread-reply`, `star`, `rescue-from-spam`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type CatalystBoostState struct {
      DomainName       string  `json:"domainName"`
      BoostFactor      float64 `json:"boostFactor"`
      RecentEngagement int     `json:"recentEngagement"`
      IsBoosted        bool    `json:"isBoosted"`
      HasMaximumCap    bool    `json:"hasMaximumCap"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS CatalystBoostState (
      domainName TEXT PRIMARY KEY,
      boostFactor REAL NOT NULL DEFAULT 1.0,
      recentEngagement INTEGER NOT NULL DEFAULT 0,
      isBoosted INTEGER NOT NULL DEFAULT 0,
      hasMaximumCap INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Ingests events from `pkg/warmup/seed/engine.go`; recalculates membrane potential.

---

### Improvement 517: Backoff Hysteresis Smoothing Governor
- **ID:** `WARMUP-VELOCITY-517`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Prevent abrupt oscillating transitions between fast dispatch and throttling by implementing a dual-threshold hysteresis loop.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-hysteresis-loop>`
  - **Lipgloss Terminal Palette:** Upper Threshold `#10b981`, Lower Threshold `#f59e0b`, Current Cadence `#06b6d4`.
  - **Shadow DOM Layout:** Dual-band graph displaying upper turn-on threshold vs lower turn-off threshold, avoiding rapid state chattering.
  - **Zero-CDN Invariant:** SVG path rendering of hysteresis loop.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity hysteresis status --domain "outbound.gitmap.dev"
  gitmap email warmup velocity hysteresis configure --high-threshold 0.85 --low-threshold 0.60
  ```
  - `--domain`: Sending domain (string).
  - `--high-threshold`: Upper engagement boundary for step-up (float).
  - `--low-threshold`: Lower engagement boundary for step-down (float).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type HysteresisConfig struct {
      DomainName     string  `json:"domainName"`
      HighThreshold  float64 `json:"highThreshold"`
      LowThreshold   float64 `json:"lowThreshold"`
      IsUpperState   bool    `json:"isUpperState"`
      HasStabilized  bool    `json:"hasStabilized"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS HysteresisConfig (
      domainName TEXT PRIMARY KEY,
      highThreshold REAL NOT NULL DEFAULT 0.80,
      lowThreshold REAL NOT NULL DEFAULT 0.50,
      isUpperState INTEGER NOT NULL DEFAULT 1,
      hasStabilized INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Smooths velocity adjustments in `pkg/warmup/neuromorphic/hysteresis.go`.

---

### Improvement 518: Micro-Burst Emission Limiter
- **ID:** `WARMUP-VELOCITY-518`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Restrict dispatch clustering to a strict maximum of `N` emails within any rolling 60-second window to prevent trigger-happy ISP spam heuristics.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-burst-limiter>`
  - **Lipgloss Terminal Palette:** Token Available `#10b981`, Token Depleted `#ef4444`, Border `#1e293b`.
  - **Shadow DOM Layout:** Token bucket visualization showing available dispatch tokens refilling in real time.
  - **Zero-CDN Invariant:** Pure CSS animated token pill indicators.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity burst status --domain "outbound.gitmap.dev"
  gitmap email warmup velocity burst set --domain "outbound.gitmap.dev" --max-burst 3 --window 60s
  ```
  - `--domain`: Sending domain (string).
  - `--max-burst`: Maximum emails in window (int).
  - `--window`: Rolling time window duration (duration).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type MicroBurstState struct {
      DomainName      string `json:"domainName"`
      MaxBurstTokens  int    `json:"maxBurstTokens"`
      AvailableTokens int    `json:"availableTokens"`
      WindowSeconds   int    `json:"windowSeconds"`
      IsRateLimited   bool   `json:"isRateLimited"`
      HasTokens       bool   `json:"hasTokens"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS MicroBurstState (
      domainName TEXT PRIMARY KEY,
      maxBurstTokens INTEGER NOT NULL DEFAULT 3,
      availableTokens INTEGER NOT NULL DEFAULT 3,
      windowSeconds INTEGER NOT NULL DEFAULT 60,
      isRateLimited INTEGER NOT NULL DEFAULT 0,
      hasTokens INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Token bucket governor in `pkg/warmup/neuromorphic/burst_limiter.go`.

---

### Improvement 519: Adaptive Hourly Ramp Curve Synthesizer
- **ID:** `WARMUP-VELOCITY-519`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Synthesize an adaptive 24-hour target volume ramp schedule that continuously recalculates based on daily milestone achievements and error ratios.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-hourly-ramp>`
  - **Lipgloss Terminal Palette:** Target Curve `#06b6d4`, Actual Volume `#10b981`, Deficit Red `#f43f5e`.
  - **Shadow DOM Layout:** 24-bar stepped histogram comparing projected hourly warmup target vs realized dispatch count.
  - **Zero-CDN Invariant:** Native SVG stepped bar chart with zero dependencies.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity ramp plan --domain "outbound.gitmap.dev"
  gitmap email warmup velocity ramp recalibrate --domain "outbound.gitmap.dev" --day 5
  ```
  - `--domain`: Sending domain (string).
  - `--day`: Current day index of warmup plan (int).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type HourlyRampPlan struct {
      PlanID          string `json:"planId"`
      DomainName      string `json:"domainName"`
      DayIndex        int    `json:"dayIndex"`
      TargetHourly    int    `json:"targetHourly"`
      RealizedHourly  int    `json:"realizedHourly"`
      IsMilestoneMet  bool   `json:"isMilestoneMet"`
      HasDeviation    bool   `json:"hasDeviation"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS HourlyRampPlan (
      planId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      dayIndex INTEGER NOT NULL,
      targetHourly INTEGER NOT NULL,
      realizedHourly INTEGER NOT NULL DEFAULT 0,
      isMilestoneMet INTEGER NOT NULL DEFAULT 0,
      hasDeviation INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Orchestrated in `pkg/warmup/neuromorphic/ramp_planner.go`.

---

### Improvement 520: Velocity Audit & Regulation Telemetry Store
- **ID:** `WARMUP-VELOCITY-520`
- **Cluster:** 52 (Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine)
- **Goal:** Persist every velocity regulation decision, damping trigger, and refractory state transition to `repodb/audit.db` for post-campaign deliverability reviews.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-velocity-audit-table>`
  - **Lipgloss Terminal Palette:** Monospace, Timestamp `#64748b`, Potential Delta `#06b6d4`, Trigger `#f59e0b`.
  - **Shadow DOM Layout:** High-density virtualized data table displaying past 100 velocity modulation events.
  - **Zero-CDN Invariant:** HTML5 table with CSS sticky headers and CSS scrollbars.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup velocity audit list --domain "outbound.gitmap.dev" --limit 20
  gitmap email warmup velocity audit export --domain "outbound.gitmap.dev" --format json
  ```
  - `--domain`: Sending domain (string).
  - `--limit`: Entry count limit (int).
  - `--format`: Output format (`json`, `csv`, `table`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type VelocityRegulationAuditEvent struct {
      EventID          string  `json:"eventId"`
      RegulatorID      string  `json:"regulatorId"`
      DomainName       string  `json:"domainName"`
      EventTrigger     string  `json:"eventTrigger"`
      PotentialDelta   float64 `json:"potentialDelta"`
      IsRefractoryHeld bool    `json:"isRefractoryHeld"`
      HasAlertTripped  bool    `json:"hasAlertTripped"`
  }
  ```
  ```sql
  -- Storage in repodb/audit.db
  CREATE TABLE IF NOT EXISTS VelocityRegulationAuditEvent (
      eventId TEXT PRIMARY KEY,
      regulatorId TEXT NOT NULL,
      domainName TEXT NOT NULL,
      eventTrigger TEXT NOT NULL,
      potentialDelta REAL NOT NULL,
      isRefractoryHeld INTEGER NOT NULL DEFAULT 0,
      hasAlertTripped INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Asynchronous persistence pipeline in `pkg/warmup/audit/velocity_recorder.go`.

---

# Cluster 53: Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs (Items 521–530)

---

### Improvement 521: Sub-Second Deliverability Heat-Wave Canvas
- **ID:** `WARMUP-HEATWAVE-521`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Render a continuous 2D chromatic heat-wave waterfall showing sub-second delivery confirmations and latency variance across all destination MX servers.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-heatwave-canvas>`
  - **Lipgloss Terminal Palette:** Background `#0b0f19`, Latency < 200ms `#10b981`, 200-500ms `#0ea5e9`, 500-1200ms `#f59e0b`, > 1200ms `#f43f5e`.
  - **Shadow DOM Layout:** Full-width canvas/SVG waterfall streaming delivery events from top to bottom. Columns represent ISP clusters; cells are color-coded by delivery duration.
  - **Zero-CDN Invariant:** HTML5 2D Canvas rendering with vanilla JavaScript pixel plotting. Zero WebGL or external charting dependencies.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave stream --domain "outbound.gitmap.dev" --speed fast
  gitmap email warmup heatwave snapshot --domain "outbound.gitmap.dev" --output-svg ./heatwave.svg
  ```
  - `--domain`: Sending domain (string).
  - `--speed`: Waterfall scroll speed (`slow`, `normal`, `fast`).
  - `--output-svg`: Path to export static SVG snapshot (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type HeatWavePixel struct {
      PixelID          string  `json:"pixelId"`
      DomainName       string  `json:"domainName"`
      IspClass         string  `json:"ispClass"`
      DurationMs       float64 `json:"durationMs"`
      TimestampEpoch   int64   `json:"timestampEpoch"`
      IsDelivered      bool    `json:"isDelivered"`
      HasAnomalyColor  bool    `json:"hasAnomalyColor"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS HeatWavePixel (
      pixelId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      ispClass TEXT NOT NULL,
      durationMs REAL NOT NULL,
      timestampEpoch INTEGER NOT NULL,
      isDelivered INTEGER NOT NULL DEFAULT 1,
      hasAnomalyColor INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Streams events over WebSockets/SSE from `pkg/warmup/heatwave/streamer.go`.

---

### Improvement 522: Bounded Jitter Distribution Graph
- **ID:** `WARMUP-HEATWAVE-522`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Plot inter-packet and dispatch timing jitter distributions within strict statistical bounds (p50, p95, p99) to detect unnatural robotic periodicity.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-jitter-graph>`
  - **Lipgloss Terminal Palette:** P50 Line `#10b981`, P95 Line `#06b6d4`, P99 Line `#f59e0b`, Outlier Red `#ef4444`.
  - **Shadow DOM Layout:** Box-and-whisker and violin distribution plot indicating variance spread. Outlier points exceeding 3-sigma bounds are flagged with warning dots.
  - **Zero-CDN Invariant:** Native SVG polygon and circle rendering computed in vanilla JS.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave jitter --domain "outbound.gitmap.dev" --interval 1h
  gitmap email warmup heatwave jitter bounds --domain "outbound.gitmap.dev"
  ```
  - `--domain`: Sending domain (string).
  - `--interval`: Aggregation window (duration).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type JitterDistributionStats struct {
      StatID       string  `json:"statId"`
      DomainName   string  `json:"domainName"`
      P50JitterMs  float64 `json:"p50JitterMs"`
      P95JitterMs  float64 `json:"p95JitterMs"`
      P99JitterMs  float64 `json:"p99JitterMs"`
      IsOrganic    bool    `json:"isOrganic"`
      HasOutliers  bool    `json:"hasOutliers"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS JitterDistributionStats (
      statID TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      p50JitterMs REAL NOT NULL,
      p95JitterMs REAL NOT NULL,
      p99JitterMs REAL NOT NULL,
      isOrganic INTEGER NOT NULL DEFAULT 1,
      hasOutliers INTEGER NOT NULL DEFAULT 0,
      calculatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Jitter calculation routines in `pkg/warmup/heatwave/jitter_calc.go`.

---

### Improvement 523: TLS Handshake Latency Breakdown Strip
- **ID:** `WARMUP-HEATWAVE-523`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Display granular timing decomposition for outbound SMTP sessions: DNS resolution, TCP SYN/ACK, STARTTLS negotiation, and certificate validation.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-tls-latency-strip>`
  - **Lipgloss Terminal Palette:** DNS `#8b5cf6`, TCP `#3b82f6`, TLS `#06b6d4`, App `#10b981`.
  - **Shadow DOM Layout:** Segmented horizontal progress bar breaking down total connection setup time into component millisecond durations.
  - **Zero-CDN Invariant:** CSS flexbox bar segments with dynamic percentage widths.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave tls-latency --domain "outbound.gitmap.dev" --sample 20
  gitmap email warmup heatwave tls-latency inspect --session-id "sess-102"
  ```
  - `--domain`: Sending domain (string).
  - `--sample`: Number of recent connections to aggregate (int).
  - `--session-id`: Specific connection identifier (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type TlsLatencyDecomposition struct {
      SessionID      string  `json:"sessionId"`
      DomainName     string  `json:"domainName"`
      DnsLookupMs    float64 `json:"dnsLookupMs"`
      TcpHandshakeMs float64 `json:"tcpHandshakeMs"`
      TlsHandshakeMs float64 `json:"tlsHandshakeMs"`
      TotalSetupMs   float64 `json:"totalSetupMs"`
      IsFastTrack    bool    `json:"isFastTrack"`
      HasTlsResumed  bool    `json:"hasTlsResumed"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS TlsLatencyDecomposition (
      sessionId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      dnsLookupMs REAL NOT NULL,
      tcpHandshakeMs REAL NOT NULL,
      tlsHandshakeMs REAL NOT NULL,
      totalSetupMs REAL NOT NULL,
      isFastTrack INTEGER NOT NULL DEFAULT 0,
      hasTlsResumed INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Trace collector in `pkg/warmup/smtp/tls_tracer.go`.

---

### Improvement 524: Greylisting Delay & Retry Latency Waterfall
- **ID:** `WARMUP-HEATWAVE-524`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Map and visualize the delay duration imposed by destination ISP greylisting rules from first 451 deferral to final 250 OK acceptance.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-greylist-waterfall>`
  - **Lipgloss Terminal Palette:** Deferral Amber `#f59e0b`, Retry Indigo `#6366f1`, Success Emerald `#10b981`.
  - **Shadow DOM Layout:** Waterfall step diagram depicting attempt 1 (deferred), retry delay interval (e.g. 15 min), and attempt 2 (delivered).
  - **Zero-CDN Invariant:** Native SVG stepped chart.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave greylist status --domain "outbound.gitmap.dev"
  gitmap email warmup heatwave greylist history --domain "outbound.gitmap.dev" --isp "yahoo"
  ```
  - `--domain`: Outbound domain (string).
  - `--isp`: Mailbox provider identifier (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type GreylistEvent struct {
      EventID          string  `json:"eventId"`
      DomainName       string  `json:"domainName"`
      RecipientDomain  string  `json:"recipientDomain"`
      FirstDeferralAt  int64   `json:"firstDeferralAt"`
      AcceptedAt       int64   `json:"acceptedAt"`
      DelayMinutes     float64 `json:"delayMinutes"`
      IsResolved       bool    `json:"isResolved"`
      HasPermanentDrop bool    `json:"hasPermanentDrop"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS GreylistEvent (
      eventId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      recipientDomain TEXT NOT NULL,
      firstDeferralAt INTEGER NOT NULL,
      acceptedAt INTEGER NOT NULL DEFAULT 0,
      delayMinutes REAL NOT NULL DEFAULT 0.0,
      isResolved INTEGER NOT NULL DEFAULT 0,
      hasPermanentDrop INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Handled by retry coordinator in `pkg/warmup/queue/greylist_tracker.go`.

---

### Improvement 525: Real-Time Seed Inbox Placement Matrix
- **ID:** `WARMUP-HEATWAVE-525`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Render a high-contrast placement grid showing seed account inbox, promotions, and spam folder categorization in real time.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-placement-matrix>`
  - **Lipgloss Terminal Palette:** Inbox Green `#10b981`, Promo Yellow `#f59e0b`, Spam Red `#ef4444`, Missing Gray `#64748b`.
  - **Shadow DOM Layout:** Interactive tabular grid grouping seed addresses by provider. Cells show folder placement badge and arrival latency.
  - **Zero-CDN Invariant:** CSS grid table with accessible ARIA grid roles.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave placement --domain "outbound.gitmap.dev"
  gitmap email warmup heatwave placement summary --domain "outbound.gitmap.dev" --format table
  ```
  - `--domain`: Sending domain (string).
  - `--format`: Display format (`table`, `json`, `text`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type SeedPlacementReport struct {
      ReportID         string  `json:"reportId"`
      DomainName       string  `json:"domainName"`
      SeedAddress      string  `json:"seedAddress"`
      FolderPlacement  string  `json:"folderPlacement"`
      ArrivalLatencyMs float64 `json:"arrivalLatencyMs"`
      IsInbox          bool    `json:"isInbox"`
      HasAlertTriggered bool   `json:"hasAlertTriggered"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS SeedPlacementReport (
      reportId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      seedAddress TEXT NOT NULL,
      folderPlacement TEXT NOT NULL,
      arrivalLatencyMs REAL NOT NULL,
      isInbox INTEGER NOT NULL DEFAULT 1,
      hasAlertTriggered INTEGER NOT NULL DEFAULT 0,
      verifiedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Feed from seed verification engine `pkg/warmup/seed/verifier.go`.

---

### Improvement 526: Delivery Latency Percentile Sparklines
- **ID:** `WARMUP-HEATWAVE-526`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Compact terminal and web sparklines displaying 60-minute rolling latency percentiles directly inside CLI status headers.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-latency-sparkline>`
  - **Lipgloss Terminal Palette:** Unicode Block Sparkline ` ▂▃▄▅▆▇█`, Accent Cyan `#06b6d4`, Danger `#f43f5e`.
  - **Shadow DOM Layout:** Inline sparkline element (`width: 120px; height: 20px`) with mini tooltip revealing min, median, and max latency.
  - **Zero-CDN Invariant:** Native SVG polyline path element.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave sparkline --domain "outbound.gitmap.dev"
  ```
  - `--domain`: Target domain (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type LatencySparklinePoint struct {
      PointID       string  `json:"pointId"`
      DomainName    string  `json:"domainName"`
      BucketMinute  int64   `json:"bucketMinute"`
      AverageMs     float64 `json:"averageMs"`
      IsPeakAnomaly bool    `json:"isPeakAnomaly"`
      HasSlowdown   bool    `json:"hasSlowdown"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS LatencySparklinePoint (
      pointId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      bucketMinute INTEGER NOT NULL,
      averageMs REAL NOT NULL,
      isPeakAnomaly INTEGER NOT NULL DEFAULT 0,
      hasSlowdown INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Rendered via Lipgloss termpad in `pkg/warmup/ui/sparkline.go`.

---

### Improvement 527: Envelope Size vs Delivery Latency Scatter Plot
- **ID:** `WARMUP-HEATWAVE-527`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Correlate MIME message byte size against transmission duration to identify payload thresholds triggering ISP content-filtering latency.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-size-latency-scatter>`
  - **Lipgloss Terminal Palette:** Normal Dot `#3b82f6`, Heavy Payload `#8b5cf6`, Anomalous Slow `#ef4444`.
  - **Shadow DOM Layout:** Cartesian SVG coordinate space plotting MIME size (KB) on X-axis and delivery latency (ms) on Y-axis.
  - **Zero-CDN Invariant:** SVG scatter plot calculation in pure ES script.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave scatter --domain "outbound.gitmap.dev" --sample-size 100
  ```
  - `--domain`: Sending domain (string).
  - `--sample-size`: Number of data points to plot (int).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type PayloadLatencyPoint struct {
      PointID          string  `json:"pointId"`
      DomainName       string  `json:"domainName"`
      PayloadBytes     int     `json:"payloadBytes"`
      TransmissionMs   float64 `json:"transmissionMs"`
      IsOversized      bool    `json:"isOversized"`
      HasContentFilter bool    `json:"hasContentFilter"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS PayloadLatencyPoint (
      pointId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      payloadBytes INTEGER NOT NULL,
      transmissionMs REAL NOT NULL,
      isOversized INTEGER NOT NULL DEFAULT 0,
      hasContentFilter INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Telemetry recorded by `pkg/warmup/heatwave/scatter_collector.go`.

---

### Improvement 528: Real-Time Bounced & Deferred Stream HUD
- **ID:** `WARMUP-HEATWAVE-528`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Terminal HUD presenting real-time streaming notifications of SMTP 4xx and 5xx responses with parsed diagnostic codes.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-bounce-stream-hud>`
  - **Lipgloss Terminal Palette:** Alert Badge Red `#f43f5e`, Warning Badge Yellow `#f59e0b`, Border `#1e293b`.
  - **Shadow DOM Layout:** Real-time event toast list highlighting error codes (`550 5.7.1 Blocked`, `421 4.7.0 Try later`) with direct recommendations.
  - **Zero-CDN Invariant:** Pure CSS slide-in animations.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave bounce-stream --domain "outbound.gitmap.dev"
  gitmap email warmup heatwave bounce-stream export --format json --output ./bounces.json
  ```
  - `--domain`: Sending domain (string).
  - `--format`: Export format (`json`, `csv`).
  - `--output`: File output path (string).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type BounceStreamEvent struct {
      EventID          string `json:"eventId"`
      DomainName       string `json:"domainName"`
      RecipientDomain  string `json:"recipientDomain"`
      SmtpCode         int    `json:"smtpCode"`
      EnhancedCode     string `json:"enhancedCode"`
      DiagnosticReason string `json:"diagnosticReason"`
      IsPermanent      bool   `json:"isPermanent"`
      HasSuppressed    bool   `json:"hasSuppressed"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS BounceStreamEvent (
      eventId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      recipientDomain TEXT NOT NULL,
      smtpCode INTEGER NOT NULL,
      enhancedCode TEXT NOT NULL,
      diagnosticReason TEXT NOT NULL,
      isPermanent INTEGER NOT NULL DEFAULT 0,
      hasSuppressed INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Connected to `pkg/warmup/bounce/parser.go`.

---

### Improvement 529: Heat-Wave Jitter Anomaly Detector & Sound Beacon
- **ID:** `WARMUP-HEATWAVE-529`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Autonomously detect statistical timing clustering that could alert spam filters, optionally triggering a terminal bell or visual beacon.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-anomaly-beacon>`
  - **Lipgloss Terminal Palette:** Beacon Flashing `#ef4444`, Normal Muted `#64748b`.
  - **Shadow DOM Layout:** Prominent pulsating alert indicator flashing on variance collapse (e.g. standard deviation < 2ms).
  - **Zero-CDN Invariant:** CSS `@keyframes pulse` animation without JS timers.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave beacon status --domain "outbound.gitmap.dev"
  gitmap email warmup heatwave beacon enable-sound --domain "outbound.gitmap.dev"
  ```
  - `--domain`: Sending domain (string).
  - `--enable-sound`: Enable terminal ASCII `\a` bell on anomaly detection (bool).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type JitterAnomalyBeaconState struct {
      DomainName        string  `json:"domainName"`
      VarianceStdDev    float64 `json:"varianceStdDev"`
      ThresholdMinDev   float64 `json:"thresholdMinDev"`
      IsAnomalyDetected bool    `json:"isAnomalyDetected"`
      HasSoundEnabled   bool    `json:"hasSoundEnabled"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS JitterAnomalyBeaconState (
      domainName TEXT PRIMARY KEY,
      varianceStdDev REAL NOT NULL,
      thresholdMinDev REAL NOT NULL DEFAULT 10.0,
      isAnomalyDetected INTEGER NOT NULL DEFAULT 0,
      hasSoundEnabled INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Evaluated in `pkg/warmup/heatwave/anomaly_detector.go`.

---

### Improvement 530: Heat-Wave Audit & Historical Telemetry Exporter
- **ID:** `WARMUP-HEATWAVE-530`
- **Cluster:** 53 (Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs)
- **Goal:** Export high-resolution heat-wave deliverability metrics and socket jitter logs to `repodb/audit.db` and external compliance bundles.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-heatwave-audit-exporter>`
  - **Lipgloss Terminal Palette:** Button Cyan `#06b6d4`, Text White `#f8fafc`, Border Gray `#334155`.
  - **Shadow DOM Layout:** Export modal with configurable date range pickers, format selector chips (`JSON`, `CSV`, `PARQUET`), and export progress indicator.
  - **Zero-CDN Invariant:** Native browser `Blob` creation and `URL.createObjectURL` download triggers.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup heatwave audit export --domain "outbound.gitmap.dev" --output-dir "./telemetry" --format json
  gitmap email warmup heatwave audit verify --domain "outbound.gitmap.dev"
  ```
  - `--domain`: Sending domain (string).
  - `--output-dir`: Export target directory (string).
  - `--format`: Format type (`json`, `csv`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type HeatWaveAuditExportRecord struct {
      ExportID         string `json:"exportId"`
      DomainName       string `json:"domainName"`
      TotalRecords     int    `json:"totalRecords"`
      ExportFilePath   string `json:"exportFilePath"`
      IsVerified       bool   `json:"isVerified"`
      HasSignedChecksum bool  `json:"hasSignedChecksum"`
  }
  ```
  ```sql
  -- Storage in repodb/audit.db
  CREATE TABLE IF NOT EXISTS HeatWaveAuditExportRecord (
      exportId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      totalRecords INTEGER NOT NULL,
      exportFilePath TEXT NOT NULL,
      isVerified INTEGER NOT NULL DEFAULT 1,
      hasSignedChecksum INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Export coordinator in `pkg/warmup/audit/heatwave_exporter.go`.
