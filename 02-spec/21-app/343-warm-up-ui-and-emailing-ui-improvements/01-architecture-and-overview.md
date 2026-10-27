# Architecture & Overview: Sequence 343 - Warm Up UI and Emailing UI Improvements (601–700)

## 1. Executive Summary & Purpose

Sequence 343 establishes the next-generation autonomous warming engine, high-fidelity terminal UI, and dynamic email diagnostic suite for `gitmap email warmup` and GitMap developer telemetry fleets. Covering Innovations 601 through 700 across 10 functional clusters (Clusters 61 to 70), this architecture bridges low-level RFC-compliant mail handling (DMARC, ARC, DKIM, SPF, MIME stream parsing) with resilient terminal visualization and zero-external-CDN web components.

All operational state, reputation telemetry, circuit-breaker metrics, and audit logs are governed by GitMap's SQLite Split-DB architecture, ensuring sub-millisecond local reads, air-gapped security, zero external CDN dependencies, and strictly positive boolean properties (`is`/`has` naming conventions only).

---

## 2. Cluster Overview (Innovations 601–700)

The 100 innovations comprising Sequence 343 are partitioned into 10 cohesive architectural clusters:

| Cluster | Range | Domain & Title | Core Architectural Focus |
|---|---|---|---|
| **Cluster 61** | 601–610 | **Cross-ISP DMARC / ARC Cryptographic Verification Engine & Dynamic Header Studio** | Cryptographic validation of Authenticated Received Chain (ARC) seals, DKIM selector rotation, DMARC alignment matrices, and interactive ASCII header inspection. |
| **Cluster 62** | 611–620 | **Autonomous Reputation Feedback Loops & Micro-Bounce Dampening Circuit Breakers** | Closed-loop reputation tracking across Gmail, Microsoft 365, Yahoo/AOL, and custom MTAs; automated circuit breaker dampening on micro-bounces and SPF softfails. |
| **Cluster 63** | 621–630 | **Live MIME Composition Tree Diffing & Attachment Payload Streamer** | High-performance AST representation of nested MIME structures, zero-copy payload streaming, boundary validation, and visual delta trees. |
| **Cluster 64** | 631–640 | **ISP Inboxing Heatmap Matrix & Seed-List Deliverability Predictor** | Probabilistic seed inboxing projections, simulated MX deliverability grid, sender IP pool distribution heatmaps, and Bayesian reputation forecast models. |
| **Cluster 65** | 641–650 | **Dynamic Volume Ramp Scheduler & Cold-Domain Velocity Governors** | Non-linear warmup velocity curves (Fibonacci, sigmoid, logarithmic), per-domain throttle caps, multi-tenant throttle synchronization, and anti-burst rate limiting. |
| **Cluster 66** | 651–660 | **Zero-CDN High-Definition ANSI Canvas Renderer & Shadow DOM Studio** | Dual-mode Lipgloss terminal canvas with 24-bit TrueColor support alongside zero-CDN Shadow DOM Custom Web Components with native SVG and WebGL rendering. |
| **Cluster 67** | 661–670 | **Interactive Bubble Tea Warmup Dashboard & Real-Time Telemetry HUD** | Terminal HUD with responsive keyboard navigation, microsecond timeline scrubber, live metric dials, threadpool status bars, and active queue watchdogs. |
| **Cluster 68** | 671–680 | **Automated Spam-Trap Sentinel & Honeypot Detection Network** | Passive MX query sniffing, MX blacklist matching, honeypot syntax pattern heuristics, recipient syntax risk score, and real-time seed sanitization. |
| **Cluster 69** | 681–690 | **Multi-Tenant SMTP/IMAP Pool Orchestrator & Credential Vault Sentinel** | Secure local credential vaulting, multi-identity connection pooling, TLS cipher suite enforcement, and automated proxy/MTA egress routing. |
| **Cluster 70** | 691–700 | **Compliance, Forensics & SQLite Split-DB Audit Event Ledger** | Tamper-evident hash-chained audit trails, RFC 5322 diagnostic archiving, one-click export to sanitized JSON/Markdown, and GitMap AUM fleet telemetry bridge. |

---

## 3. Detailed Subsystem Architectures

### 3.1 Cross-ISP DMARC / ARC Cryptographic Verification Engine & Dynamic Header Studio
RFC 8617 (ARC) and RFC 7489 (DMARC) validation require multi-stage cryptographic evaluation across untrusted intermediaries. The engine parses and verifies ARC-Authentication-Results (`AAR`), ARC-Message-Signature (`AMS`), and ARC-Seal (`AS`) sets across sequential message validation instances (`i=1..N`).

Key verification components:
1. **ARC Seal Chain Traversal**: Recursively evaluates cryptographic signatures using Ed25519 or RSA public keys resolved via DNS TXT records.
2. **DMARC Strict vs Relaxed Alignment**: Computes Header From (`RFC 5322`) alignment against Envelope From (`RFC 5321 Return-Path`) and DKIM `d=` parameters across organizational domains.
3. **Dynamic Header Studio**: Renders tree structures of parsed RFC 5322 headers, highlighting anomalous hops, DKIM canonicalization mutations (`relaxed/relaxed` vs `simple/simple`), and SPF routing mismatches.

```
+--------------------------------------------------------------------------+
|                 Dynamic Header Studio & ARC Parser                       |
|                                                                          |
|  [RFC 5322 Raw Header Stream]                                            |
|              |                                                           |
|              v                                                           |
|  +------------------------+      +------------------------------------+  |
|  | ARC-Seal Validator (AS)| ---> | DNS Public Key Resolver (Cached)   |  |
|  +------------------------+      +------------------------------------+  |
|              |                                                           |
|              v                                                           |
|  +------------------------+      +------------------------------------+  |
|  | ARC-Message-Sig (AMS)  | ---> | Body Canonicalization Hash Check   |  |
|  +------------------------+      +------------------------------------+  |
|              |                                                           |
|              v                                                           |
|  +------------------------+      +------------------------------------+  |
|  | ARC-Auth-Results (AAR) | ---> | Alignment Matrix (DKIM / SPF)      |  |
|  +------------------------+      +------------------------------------+  |
|              |                                                           |
|              +-------------------+                                       |
|                                  v                                       |
|                     [Verified Cryptographic Chain]                       |
+--------------------------------------------------------------------------+
```

### 3.2 Autonomous Reputation Feedback Loops & Micro-Bounce Dampening Circuit Breakers
Warmup velocity must dynamically adapt when receiving soft bounces (`421`, `450`, `451`, `452`) or ISP rate-limit deferrals. The dampening circuit breaker operates as a state machine:
- **State `Closed` (Normal Operation)**: Target warmup velocity active; throughput matches target schedule.
- **State `Half-Open` (Probabilistic Probing)**: Triggered after dampening timer expires; sends low-velocity synthetic probes to assess ISP queue latency and rejection codes.
- **State `Open` (Trip & Throttle)**: Tripped when micro-bounce rate exceeds threshold ($> 1.8\%$ within a rolling 50-send window). Instantly dials back dispatch rate by $50\%\dots 90\%$ and halts high-risk recipient segments.

```
           +-----------------------------+
           |        Normal Sending       |
           |       (State: Closed)       |
           +-----------------------------+
                          |
              Micro-bounce rate > 1.8%
                          |
                          v
           +-----------------------------+
           |       Trip & Throttle       | <--------+
           |        (State: Open)        |          | Probe fails
           +-----------------------------+          |
                          |                         |
               Cool-off window expires              |
                          |                         |
                          v                         |
           +-----------------------------+          |
           |      Adaptive Probing       | ---------+
           |      (State: Half-Open)     |
           +-----------------------------+
                          |
              Probe pass (100% green)
                          |
                          v
                 [Resume Back to Closed]
```

### 3.3 Zero-CDN High-Definition ANSI Canvas Renderer
Visual components run in two environments without relying on external network requests or CDN dependencies:
1. **Terminal Canvas (Bubble Tea + Lipgloss)**: Emits 24-bit TrueColor ANSI escape sequences. Utilizes cell buffers with double-buffering to eliminate terminal flicker during 60 FPS real-time warmup streaming.
2. **Web Component Studio (Shadow DOM)**: Custom elements (`<warmup-arc-inspector>`, `<warmup-bounce-circuit>`, `<warmup-mime-tree>`) encapsulate local CSS, inline SVG icons, and DOM logic. Zero CDN script tags, external web fonts, or tracker scripts are used.

```
   CLI / Terminal Mode                        Web / GUI Mode
+-------------------------+             +-------------------------+
|  Lipgloss Token Palette |             | Shadow DOM Web Element  |
|  ANSI TrueColor 24-bit  |             | Encapsulated CSS3 / SVG |
|  Double-Buffered Canvas |             | Pure Offline Vanilla JS |
+-------------------------+             +-------------------------+
              ^                                       ^
              |                                       |
              +-------------------+-------------------+
                                  |
                      [Unified ViewModel Provider]
                                  |
                      [Split-DB SQLite Storage]
```

---

## 4. Split-DB SQLite Storage Architecture

GitMap enforces clear architectural separation across database boundaries. Warmup domain operational data resides in `repodb/warmup.db`, while immutable security audit events and diagnostic logs reside in `repodb/audit.db`.

### 4.1 `repodb/warmup.db` (Operational Warmup & Deliverability Engine)

```sql
-- Warmup Schedules and Velocity Configuration
CREATE TABLE IF NOT EXISTS WarmupCampaigns (
    id TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    senderIdentity TEXT NOT NULL,
    strategyCurve TEXT NOT NULL, -- 'fibonacci', 'sigmoid', 'linear'
    currentStep INTEGER NOT NULL DEFAULT 1,
    targetDailyVolume INTEGER NOT NULL,
    sentTodayCount INTEGER NOT NULL DEFAULT 0,
    isPaused INTEGER NOT NULL DEFAULT 0,
    hasAdaptiveThrottling INTEGER NOT NULL DEFAULT 1,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Real-Time ISP Reputation State
CREATE TABLE IF NOT EXISTS IspReputationStates (
    id TEXT PRIMARY KEY,
    campaignId TEXT NOT NULL,
    ispProvider TEXT NOT NULL, -- 'google', 'microsoft', 'yahoo', 'custom'
    reputationScore REAL NOT NULL DEFAULT 100.0,
    consecutiveBounces INTEGER NOT NULL DEFAULT 0,
    circuitBreakerState TEXT NOT NULL DEFAULT 'closed', -- 'closed', 'half_open', 'open'
    isThrottled INTEGER NOT NULL DEFAULT 0,
    hasFeedbackLoopEnabled INTEGER NOT NULL DEFAULT 1,
    lastDampenedAt DATETIME,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaignId) REFERENCES WarmupCampaigns(id) ON DELETE CASCADE
);

-- Cryptographic ARC Verification Records
CREATE TABLE IF NOT EXISTS ArcVerifications (
    id TEXT PRIMARY KEY,
    messageId TEXT NOT NULL,
    campaignId TEXT NOT NULL,
    instanceIndex INTEGER NOT NULL,
    sealStatus TEXT NOT NULL, -- 'pass', 'fail', 'none'
    amsStatus TEXT NOT NULL,
    dmarcAlignment TEXT NOT NULL, -- 'strict_pass', 'relaxed_pass', 'fail'
    isChainValid INTEGER NOT NULL DEFAULT 0,
    hasDkimAlignment INTEGER NOT NULL DEFAULT 0,
    hasSpfAlignment INTEGER NOT NULL DEFAULT 0,
    verifiedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaignId) REFERENCES WarmupCampaigns(id) ON DELETE CASCADE
);
```

### 4.2 `repodb/audit.db` (Immutable Forensic & Deliverability Events)

```sql
-- Hash-Chained Audit Ledger for Cryptographic Forensics
CREATE TABLE IF NOT EXISTS WarmupAuditLedger (
    id TEXT PRIMARY KEY,
    campaignId TEXT NOT NULL,
    eventType TEXT NOT NULL, -- 'arc_validation', 'circuit_trip', 'rate_throttle', 'seed_delivery'
    payloadHash TEXT NOT NULL,
    previousBlockHash TEXT NOT NULL,
    currentBlockHash TEXT NOT NULL,
    isCompromised INTEGER NOT NULL DEFAULT 0,
    hasSignatureVerified INTEGER NOT NULL DEFAULT 1,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Detailed Bounce & Diagnostic Raw Event Archive
CREATE TABLE IF NOT EXISTS BounceDiagnostics (
    id TEXT PRIMARY KEY,
    campaignId TEXT NOT NULL,
    recipientDomain TEXT NOT NULL,
    diagnosticCode TEXT NOT NULL,
    smtpReplyCode INTEGER NOT NULL,
    enhancedStatusCode TEXT NOT NULL,
    isTransient INTEGER NOT NULL DEFAULT 1,
    hasCircuitTripped INTEGER NOT NULL DEFAULT 0,
    rawHeaders TEXT NOT NULL,
    recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 5. Coding & Naming Guidelines Compliance

1. **Boolean Naming Rigor**: Every boolean column in SQLite schemas (`isPaused`, `hasAdaptiveThrottling`, `isThrottled`, `hasFeedbackLoopEnabled`, `isChainValid`, `hasDkimAlignment`, `hasSpfAlignment`, `isCompromised`, `hasSignatureVerified`, `isTransient`, `hasCircuitTripped`) and every Go struct field begins exclusively with `is` or `has`. Negative booleans and legacy prefixes (`can`, `should`, `was`, `will`) are strictly excluded.
2. **Relative Path Hygiene**: All referenced paths use relative repository structure (e.g. `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/`, `repodb/warmup.db`, `repodb/audit.db`). No Windows drive letters (`D:\`), absolute user paths, or localhost URLs exist in the specification.
3. **Split-DB Isolation**: No cross-database SQLite foreign keys or queries are permitted. The operational state is completely decoupled from the audit ledger.
