# Specification 342: Warm-Up UI and Emailing UI Improvements (Overview & Architecture)

- **Sequence:** 342
- **Slug:** `342-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author 1

---

## 1. Executive Summary & Core Innovations (Innovations 501 to 600)

Sequence 342 advances the GitMap Warm-Up & Outbound Emailing platform into an enterprise-grade, neuromorphic execution and physical connection layer. Where Sequence 341 established the predictive Bayesian inference and geo-spatial delivery telemetry, Sequence 342 provides real-time physical connection multiplexing, synthetic nervous system velocity regulators, sub-millisecond deliverability heat-wave rendering, and verifiable socket health matrices.

```
                          ┌──────────────────────────────────────────────┐
                          │    GitMap Sequence 342 Architecture Matrix   │
                          └──────────────────────┬───────────────────────┘
                                                 │
                   ┌─────────────────────────────┼─────────────────────────────┐
                   ▼                             ▼                             ▼
       ┌───────────────────────┐     ┌───────────────────────┐     ┌───────────────────────┐
       │  Cluster 51 (501-510) │     │  Cluster 52 (511-520) │     │  Cluster 53 (521-530) │
       │  Adaptive Multi-ISP   │     │  Neuromorphic Velocity│     │  Deliverability Heat- │
       │  Connection Pool &    │     │  Regulator & Dynamic  │     │  Wave Visualizer &    │
       │  Socket Health Matrix │     │  Send-Cadence Engine  │     │  Bounded Jitter Graphs│
       └───────────┬───────────┘     └───────────┬───────────┘     └───────────┬───────────┘
                   │                             │                             │
                   ├─────────────────────────────┴─────────────────────────────┤
                   ▼                                                           ▼
       ┌───────────────────────┐                                   ┌───────────────────────┐
       │ Clusters 54-57        │                                   │ Clusters 58-60        │
       │ (531-570) Deep Auth,  │                                   │ (571-600) Autonomous  │
       │ MIME Synthetic Jitter,│                                   │ Multi-Tenant Sharding,│
       │ Micro-Burst Dampening │                                   │ Cold-to-Hot Gates     │
       └───────────┬───────────┘                                   └───────────┬───────────┘
                   │                                                           │
                   └─────────────────────────────┬─────────────────────────────┘
                                                 │
                                                 ▼
                         ┌──────────────────────────────────────────────┐
                         │      Core Architecture & Shared Invariants   │
                         │   - Split-DB SQLite (`warmup.db`, `audit.db`)│
                         │   - Error Return Sovereignty (Zero Exits)    │
                         │   - Pure HTML5/SVG/CSS3 (Zero External CDN)  │
                         │   - Coupled CLI Tools (`gitmap email warmup`)│
                         │   - Positive Booleans (`is*`, `has*` Only)   │
                         └──────────────────────────────────────────────┘
```

The 100 innovations comprising Sequence 342 span 10 dedicated functional clusters (Clusters 51 to 60):

1. **Cluster 51 (501–510): Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix**: Persistent TCP socket state machines, dynamic keep-alive tuning, TLS session cache reuse across Google, Microsoft, and Yahoo MX clusters, and real-time socket health matrices.
2. **Cluster 52 (511–520): Neuromorphic Velocity Regulator & Dynamic Send-Cadence Engine**: Synthetic nervous system mimicking human biological circadian rhythms, dynamic Poisson inter-arrival scheduling, refractory period enforcement, and feedback-damped volume acceleration.
3. **Cluster 53 (521–530): Live Deliverability Heat-Wave Visualizer & Bounded Jitter Graphs**: Ultra-low-overhead SVG and terminal Lipgloss heat-wave waterfalls depicting sub-second delivery confirmation, jitter bounds, TLS handshake overhead, and greylisting delays.
4. **Cluster 54 (531–540): Deep Cryptographic Authenticity & Ephemeral DKIM Key Mesh**: Ephemeral rotating DKIM signatures, sub-hour selector failover, DNS EDNS0 payload optimization, and automated ARC (Authenticated Received Chain) validation.
5. **Cluster 55 (541–550): Content Polymorphism & Dynamic MIME Tree Synthesizer**: Micro-structural body permutations, multi-part alternative boundary randomization, invisible zero-width typography variation, and heuristic anti-fingerprinting.
6. **Cluster 56 (551–560): Adaptive Concurrency Throttle & Micro-Burst Damping Circuit**: Hardware-informed leaky-bucket rate limiters, backoff hysteresis loops, TCP window scale adaptation, and zero-drop buffer management.
7. **Cluster 57 (561–570): Autonomous Reputation Sentinel & MX Ban-Wave Predictor**: Predictive anti-abuse pattern recognition, RBL query decentralization, automated IP pool rotation, and preemptive quarantine triggers.
8. **Cluster 58 (571–580): Seed Account Engagement Simulator & Swarm Thread Emulation**: Autonomous IMAP/POP3 interaction loops, thread replying, star/label workflows, spam folder rescue actions, and human dwell-time emulation.
9. **Cluster 59 (581–590): Multi-Tenant IP Sharding & Hierarchical Sub-Domain Governance**: Partitioned pool routing, isolation barriers between transactional and warmup traffic, parent/child reputation weighting, and automated egress gateway balancing.
10. **Cluster 60 (591–600): Cold-to-Hot Transition Gatekeeper & Enterprise Production Handoff**: Strict graduation criteria validation, dual-homed DNS canary routing, automated cryptographic handoff tokens, and production queue safety interlocks.

---

## 2. Deep Architectural Design of Core Subsystems

### 2.1 Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix

The Adaptive Multi-ISP Connection Pool manages outbound SMTP socket lifecycles with granular per-ISP tuning. Outbound connections are multiplexed, verified, and recycled without incurring redundant TCP three-way handshake and TLS negotiation overhead:

```
  ┌─────────────────────────────────────────────────────────────────────────────┐
  │ Outbound Dispatch Router                                                    │
  └──────────────────────────────────────┬──────────────────────────────────────┘
                                         │
                   ┌─────────────────────┴─────────────────────┐
                   ▼                                           ▼
       ┌───────────────────────┐                   ┌───────────────────────┐
       │ Google Pool (Port 25) │                   │ Microsoft Pool (25)   │
       │ - Max Idle: 12        │                   │ - Max Idle: 8         │
       │ - Keep-Alive: 30s     │                   │ - Keep-Alive: 20s     │
       │ - TLS 1.3 Fast Open   │                   │ - TLS Session Resump  │
       └───────────┬───────────┘                   └───────────┬───────────┘
                   │                                           │
                   └─────────────────────┬─────────────────────┘
                                         │
                                         ▼
  ┌─────────────────────────────────────────────────────────────────────────────┐
  │ Socket Health Matrix (repodb/warmup.db -> SocketHealthRecord)               │
  │ - Active: `isActive`                                                        │
  │ - Verified: `hasPassedHealthCheck`                                          │
  │ - Reused TLS: `hasTlsSessionResumed`                                        │
  │ - Degradation: `isDegraded`                                                 │
  └─────────────────────────────────────────────────────────────────────────────┘
```

1. **Deterministic Socket State Engine**: Sockets exist in strict states: `Idle`, `InUse`, `Verifying`, and `Retiring`. Sockets exceeding the ISP-specific keep-alive deadline are closed gracefully via `QUIT` rather than abrupt RST teardown.
2. **Session Resumption & TLS Cache**: TLS session tickets are preserved in memory and cached in SQLite `warmup.db`, slashing handshakes from 2-RTT to 0-RTT/1-RTT for repeated connections.
3. **Continuous Socket Health Probing**: Health probes issue lightweight `NOOP` or `RSET` commands periodically. If an ISP drops connection or introduces latency > 1,200ms, the pool flags `isDegraded = 1` and diverts pending dispatches to healthy backup sockets.

### 2.2 Neuromorphic Feedback Loop for Inbox Velocity

The Neuromorphic Velocity Regulator abandons simplistic linear ramp schedules in favor of a biological neural firing rate model. The engine computes send cadence based on an integrate-and-fire dynamic:

```
  Feedback Inputs:
  - Seed Open Rate: `hasPositiveSignal` (Membrane Potential +ΔV)
  - Greylist / 421: `isThrottled`        (Membrane Potential -ΔV)
  - Spam Placement: `hasSpamPlacement`   (Enter Refractory State)
                       │
                       ▼
         ┌───────────────────────────┐
         │ Leaky Integrate & Fire    │
         │ - Leak Rate: λ_decay      │
         │ - Threshold: V_threshold  │
         │ - Dynamic Refractory: τ   │
         └─────────────┬─────────────┘
                       │
         ┌─────────────┴─────────────┐
         ▼                           ▼
  [Spike: Dispatch Batch]     [Refractory Hold]
```

1. **Membrane Potential Calculation**: Positive feedback (verified inbox placement, seed thread replies, stars) increases velocity potential. Negative feedback (421 greylisting, SMTP temporary deferrals, spam folder placement) immediately depresses potential.
2. **Refractory Period Invariant**: Following any ISP deferral or rapid step-up, a dynamic refractory pause (`τ_refractory`) is enforced. During this interval, only background seed probes are permitted; warmup dispatches are suspended to avoid triggering automated spam heuristics.
3. **Stochastic Jitter & Poisson Distribution**: Inter-dispatch delays are sampled from a Poisson process parameterized by the current neuromorphic firing rate, guaranteeing organic dispatch spacing with mathematical randomness.

### 2.3 Zero-CDN Visual Web Component Architecture

All visual widgets are implemented as native Web Components with full Shadow DOM encapsulation, adhering to strict zero-CDN principles:

```
  ┌─────────────────────────────────────────────────────────────────┐
  │ <warmup-connection-matrix> / <warmup-velocity-gauge>            │
  │                                                                 │
  │  #shadow-root (open)                                            │
  │  ┌───────────────────────────────────────────────────────────┐  │
  │  │ <style>                                                   │  │
  │  │   :host { display: block; font-family: ui-monospace; }   │  │
  │  │   /* WCAG AAA High Contrast Lipgloss-Mapped Tokens */     │  │
  │  │   --color-bg-base: #0b0f19;                               │  │
  │  │   --color-border: #1e293b;                                │  │
  │  │   --color-cyan: #06b6d4;                                  │  │
  │  │   --color-emerald: #10b981;                               │  │
  │  │   --color-amber: #f59e0b;                                 │  │
  │  │   --color-rose: #f43f5e;                                  │  │
  │  │ </style>                                                  │  │
  │  │ <svg viewBox="0 0 800 320" preserveAspectRatio="none">   │  │
  │  │   <!-- Native SVG Paths & Polylines Rendered in Go/ES --> │  │
  │  │ </svg>                                                    │  │
  │  │ <div class="matrix-hud" aria-live="polite">               │  │
  │  │   <span class="badge is-healthy">HEALTHY</span>           │  │
  │  │ </div>                                                    │  │
  │  └───────────────────────────────────────────────────────────┘  │
  └─────────────────────────────────────────────────────────────────┘
```

- **Zero External Dependencies**: No Tailwind CDN, no Chart.js, no D3, and no external font loading.
- **Lipgloss Palette Alignment**: All styles strictly match GitMap's terminal Lipgloss design system tokens (`#0b0f19`, `#1e293b`, `#06b6d4`, `#10b981`, `#f59e0b`, `#f43f5e`).
- **WCAG AAA Compliance**: Minimum 7:1 contrast ratio for all textual elements and UI indicators. Accessible aria labels and live regions for screen reader compatibility.

---

## 3. Split-DB SQLite Storage Architecture

GitMap strictly separates warmup operational telemetry from compliance audit logging across two SQLite databases:

```
  repodb/
  ├── warmup.db   <── Operational runtime tables (connection pools, sockets, velocity)
  └── audit.db    <── Immutably appended compliance & event records (security, audit trail)
```

### 3.1 `repodb/warmup.db` Tables

```sql
-- Connection Pool State
CREATE TABLE IF NOT EXISTS ConnectionPoolState (
    poolId TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    ispClass TEXT NOT NULL,
    maxIdleSockets INTEGER NOT NULL DEFAULT 8,
    activeSocketCount INTEGER NOT NULL DEFAULT 0,
    idleSocketCount INTEGER NOT NULL DEFAULT 0,
    keepAliveSeconds INTEGER NOT NULL DEFAULT 30,
    isActive INTEGER NOT NULL DEFAULT 1,
    hasTlsSessionCache INTEGER NOT NULL DEFAULT 1,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Individual Socket Health Record
CREATE TABLE IF NOT EXISTS SocketHealthRecord (
    socketId TEXT PRIMARY KEY,
    poolId TEXT NOT NULL,
    remoteAddress TEXT NOT NULL,
    localPort INTEGER NOT NULL,
    handshakeDurationMs REAL NOT NULL,
    latencyRttMs REAL NOT NULL,
    isActive INTEGER NOT NULL DEFAULT 1,
    isDegraded INTEGER NOT NULL DEFAULT 0,
    hasPassedHealthCheck INTEGER NOT NULL DEFAULT 1,
    hasTlsSessionResumed INTEGER NOT NULL DEFAULT 0,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(poolId) REFERENCES ConnectionPoolState(poolId) ON DELETE CASCADE
);

-- Neuromorphic Velocity State
CREATE TABLE IF NOT EXISTS NeuromorphicVelocityState (
    regulatorId TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    ispClass TEXT NOT NULL,
    membranePotential REAL NOT NULL DEFAULT 0.0,
    firingThreshold REAL NOT NULL DEFAULT 1.0,
    leakRate REAL NOT NULL DEFAULT 0.05,
    refractorySeconds INTEGER NOT NULL DEFAULT 0,
    currentCadenceHourly INTEGER NOT NULL DEFAULT 5,
    maxCadenceHourly INTEGER NOT NULL DEFAULT 100,
    isActive INTEGER NOT NULL DEFAULT 1,
    isInRefractoryPeriod INTEGER NOT NULL DEFAULT 0,
    hasFiringTripped INTEGER NOT NULL DEFAULT 0,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Deliverability Heat-Wave Point
CREATE TABLE IF NOT EXISTS DeliverabilityHeatPoint (
    pointId TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    ispClass TEXT NOT NULL,
    timestampEpoch INTEGER NOT NULL,
    deliveryLatencyMs REAL NOT NULL,
    jitterMs REAL NOT NULL,
    placementStatus TEXT NOT NULL, -- 'inbox', 'spam', 'promotions', 'deferred'
    isSuccessful INTEGER NOT NULL DEFAULT 1,
    hasJitterAnomaly INTEGER NOT NULL DEFAULT 0,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 3.2 `repodb/audit.db` Tables

```sql
-- Connection Audit Trail
CREATE TABLE IF NOT EXISTS ConnectionAuditTrail (
    auditId TEXT PRIMARY KEY,
    socketId TEXT NOT NULL,
    poolId TEXT NOT NULL,
    eventType TEXT NOT NULL, -- 'CONNECT', 'TLS_HANDSHAKE', 'NOOP_PROBE', 'DISCONNECT', 'EVICT'
    responseStatus TEXT NOT NULL,
    latencyMs REAL NOT NULL,
    isSuccess INTEGER NOT NULL DEFAULT 1,
    hasSecurityViolation INTEGER NOT NULL DEFAULT 0,
    recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Velocity Regulation Audit Event
CREATE TABLE IF NOT EXISTS VelocityAuditEvent (
    eventId TEXT PRIMARY KEY,
    regulatorId TEXT NOT NULL,
    domainName TEXT NOT NULL,
    eventTrigger TEXT NOT NULL, -- 'SEED_OPEN', 'BOUNCE_421', 'SPAM_COMPLAINT', 'RECOVERY'
    potentialDelta REAL NOT NULL,
    cadenceBefore INTEGER NOT NULL,
    cadenceAfter INTEGER NOT NULL,
    isRefractoryEngaged INTEGER NOT NULL DEFAULT 0,
    hasAlertFired INTEGER NOT NULL DEFAULT 0,
    recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## 4. Coding & Architecture Invariants

1. **Error Return Sovereignty**: All leaf Go functions must return `error` or `*appfault.AppError` directly. Never invoke `cliexit.HandleError` or `os.Exit` inside leaf packages. Exits are handled solely by command orchestrators.
2. **Boolean Naming**: All booleans in Go structs, JSON contracts, database columns, and frontend classes must start strictly with `is` or `has`. Zero negative booleans (`isNotValid` is banned; use `isValid` or `hasFailed`).
3. **Zero Magic Exit Codes**: Only strongly-typed enums (`constants.ExitSuccess`, `constants.ExitGeneralError`, `constants.ExitUsageError`).
4. **Function Line Limits**: Max 15 lines per function, preferred <= 8 lines. Early returns and guard clauses; zero nested `if` blocks.
5. **Path Hygiene**: 100% relative repository paths (`02-spec/21-app/...`). Never include drive letters like `D:\` or absolute workstation paths.
