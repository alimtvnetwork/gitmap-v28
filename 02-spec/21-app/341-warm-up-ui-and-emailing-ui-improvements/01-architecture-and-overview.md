# Specification 341: Warm-Up UI and Emailing UI Improvements (Overview & Architecture)

- **Sequence:** 341
- **Slug:** `341-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## 1. Executive Summary & Core Innovations (Innovations 401 to 500)

Sequence 341 establishes the predictive intelligence, geo-spatial delivery telemetry, and autonomous reputation resilience layer for the GitMap Warm-Up & Outbound Emailing platform. Building upon the multi-agent consensus swarm, cognitive sending rhythms, and peer-to-peer fleet mesh introduced in Sequence 340, Sequence 341 shifts warmup and sending execution from reactive adjustment to proactive Bayesian and neural-mesh forecasting:

```
                          ┌──────────────────────────────────────────────┐
                          │    GitMap Sequence 341 Architecture Matrix   │
                          └──────────────────────┬───────────────────────┘
                                                 │
                   ┌─────────────────────────────┼─────────────────────────────┐
                   ▼                             ▼                             ▼
       ┌───────────────────────┐     ┌───────────────────────┐     ┌───────────────────────┐
       │  Cluster 41 (401-410) │     │  Cluster 42 (411-420) │     │  Cluster 43 (421-430) │
       │  Bayesian Ramp Curve  │     │  Inbox Placement Mesh │     │  Reputation Sentinel  │
       │  & ML Gradient Tuner  │     │  & ISP Latency Map    │     │  & Dynamic Auth Graph │
       └───────────┬───────────┘     └───────────┬───────────┘     └───────────┬───────────┘
                   │                             │                             │
                   ├─────────────────────────────┴─────────────────────────────┤
                   ▼                                                           ▼
       ┌───────────────────────┐                                   ┌───────────────────────┐
       │  Cluster 44 (431-440) │                                   │  Cluster 45 (441-450) │
       │  Predictive Deliver-  │                                   │  Seed Inbox Swarm     │
       │  ability Neural Mesh  │                                   │  Engagement Simulator │
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

The 100 innovations comprising Sequence 341 span 10 dedicated functional clusters:
1. **Cluster 41 (401–410): Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner**: Continuous statistical inference calculating optimal volume step-up curves conditioned on seed inbox engagement rates, bounce feedback intervals, and reputation trajectory vectors.
2. **Cluster 42 (411–420): Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map**: Zero-dependency vector cartography rendering sub-second seed delivery confirmation, MX server round-trip latency, TLS handshake overhead, and inbox/spam/promotions placement breakdown across global autonomous regions.
3. **Cluster 43 (421–430): Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer**: Continuous DNS record hygiene validation, cryptographic key rotation scheduling, DMARC forensic aggregate parse tree inspection, and automated DNS propagation health graphs.
4. **Cluster 44 (431–440): Predictive Deliverability Neural Mesh**: Multi-layer perceptron and recurrent deliverability inference engine running locally via lightweight weight matrices in Go to detect ISP rate-limiting thresholds before soft bounces occur.
5. **Cluster 45 (441–450): Seed Inbox Autonomous Engagement Simulator**: Simulates organic user behaviors (read duration, mark as important, star thread, reply loop, folder moves) across seed accounts with humanized timing jitter.
6. **Cluster 46 (451–460): SMTP Connection Multiplexing & Adaptive Concurrency Pipeline**: Dynamic TCP connection pooling, pipelined STARTTLS sessions, and automated IP pool failover under ISP grey-listing.
7. **Cluster 47 (461–470): Content Entropy, Fingerprint Obfuscation & Micro-Variance Engine**: Real-time MIME structure fuzzing, spintax tree mutation, CSS class isolation, and invisible watermark synthesis preventing template fingerprinting by heuristic filters.
8. **Cluster 48 (471–480): Feedback Loop (FBL) Stream Consumer & Instant Unsubscribe Orchestrator**: Sub-second ARF (Abuse Reporting Format) parsing, automated suppression list replication, and instant one-click unsubscribe header enforcement.
9. **Cluster 49 (481–490): Multi-Tenant IP Pool Warmup Sharding & Sub-Domain Hierarchy Visualizer**: Cross-tenant isolation boundaries, shared IP warmup credits, parent-child DNS hierarchy graphs, and reverse-DNS pointer validation.
10. **Cluster 50 (491–500): Cold-to-Hot Transition Gatekeeper & Enterprise Production Handoff**: Automated milestone gating, compliance audit certificate generator, and production queue traffic cross-fading without volume spikes.

---

## 2. Deep Architectural Design of Core Subsystems

### 2.1 Predictive Deliverability Neural Mesh

The Predictive Deliverability Neural Mesh operates as a sub-millisecond local inference engine evaluating outbound dispatch vectors before SMTP transmission. It transforms historical ISP bounce profiles, current queue concurrency, recent complaint ratios, and time-of-day circadian factors into a continuous risk score:

```
  ┌──────────────────┐    ┌──────────────────────┐    ┌─────────────────────────┐
  │ Dispatch Request │───▶│ Feature Vectorizer   │───▶│ Forward Neural Pipeline │
  │ - ISP Target     │    │ - `hasRecentBounce`  │    │ - Hidden Layer 1 (ReLU) │
  │ - Payload Size   │    │ - `isPeakCircadian`  │    │ - Hidden Layer 2 (Tanh) │
  │ - IP Pool Age    │    │ - `hasTls13Support`  │    │ - Softmax Risk Tensor   │
  └──────────────────┘    └──────────────────────┘    └────────────┬────────────┘
                                                                   │
                                                                   ▼
  ┌─────────────────────────────────────────────────────────────────────────────┐
  │ Decision & Routing Hub                                                      │
  │ - Risk Score <= 0.25: Direct Instant Dispatch                               │
  │ - Risk Score 0.26 - 0.70: Adaptive Delay Injection & Concurrency Damping   │
  │ - Risk Score > 0.70: Autonomous IP Pool Reroute & Failsafe Hold             │
  └─────────────────────────────────────────────────────────────────────────────┘
```

1. **Feature Ingestion & Vector Normalization:** Telemetry signals are normalized into a fixed-length float vector `x in [0, 1]^16`. Categorical states are encoded using positive boolean indicators (`hasRecentBounce`, `isPeakCircadian`, `hasTls13Support`, `isSeedVerified`).
2. **Deterministic Go Inference Kernel:** To ensure zero external Python runtime dependencies, model weights are stored in embedded read-only byte arrays and calculated via optimized Go matrix multiplication routines with SIMD loop unrolling.
3. **Continuous Local Calibration:** Whenever an SMTP server returns a `250 OK`, `421 Try Again Later`, `451 Greylisted`, or `550 Spam Rejected` response, the local inference weights receive a stochastic gradient descent update written to `repodb/warmup.db`.

### 2.2 Bayesian Warmup Ramp Simulator

The Bayesian Warmup Ramp Simulator models volume progression not as a static linear schedule, but as a posterior probability distribution conditioned on real-world deliverability feedback:

```
  Prior Distribution P(Volume|Day)
                │
                ▼
  [Daily Dispatch Phase] ──▶ Observed Feedback (Bounces, Seed Placement, Spam Flags)
                │
                ▼
  Likelihood Function L(Feedback|Volume)
                │
                ▼
  Posterior Distribution P(Volume|Feedback) ──▶ Maximum A Posteriori (MAP) Next Target
```

- **Prior Specification:** Beta-Binomial conjugacy is maintained for each ISP destination (Google, Microsoft 365, Yahoo/AOL, Proton, Private MX). The prior represents the historical baseline ramp velocity for that provider class.
- **Likelihood Computation:** When daily seed tests demonstrate high inbox placement (>95%) and zero soft bounces, the likelihood distribution shifts rightward, enabling accelerated volume increases (+25% to +40% per step).
- **Conservative Guard Rails:** If the posterior distribution indicates a variance escalation or a non-zero probability of reputation downgrade, the simulator triggers an automatic dampening coefficient, clamping daily growth to 0% or initiating an emergency 20% rollback.

### 2.3 Real-Time Inbox Placement Heatmap

The Real-Time Inbox Placement Heatmap provides sub-second visibility into global deliverability performance without relying on third-party SaaS tracking APIs:

- **Geo-Spatial Grid Rendering:** Native SVG canvas projecting standard Equirectangular or Mercator topologies with GPU-accelerated styling. Points represent physical MX edge locations (e.g., Ashburn, Frankfurt, Singapore, Dublin).
- **Placement Triangulation:** Seed accounts deployed across distinct autonomous system numbers (ASNs) report delivery status via authenticated webhook pings.
- **Micro-State Classification:** Each coordinate node visualizes one of four discrete placement classifications:
  - `Inbox Placement` (`#10b981` Emerald glow)
  - `Promotions / Tabbed` (`#06b6d4` Cyan)
  - `Spam / Junk Folder` (`#f59e0b` Amber warn)
  - `Blocked / Dropped` (`#ef4444` Rose alert)
- **Zero Heavy Graphing Libraries:** All chart lines, nodes, connections, and pulse vectors are calculated dynamically in native Go CLI terminal renderers (via Lipgloss and ANSI matrix blocks) and browser Web Components (via native SVG `<path>` and `<circle>` tags).

### 2.4 Zero-CDN Reactive Web Visualizer

In strict accordance with GitMap UI guidelines, all visual interfaces run completely offline with zero external CDNs, Google Fonts, or runtime dependencies:

- **Encapsulated Shadow DOM Components:** Every visualization is isolated within an Autonomous Custom Element exposing an open shadow root (`mode: 'open'`). This prevents host page style pollution and ensures 100% theme isolation.
- **CSS Custom Property Design System:** Standard design tokens:
  - `--warmup-bg-deep`: `#0b0f19` (Obsidian baseline)
  - `--warmup-bg-card`: `#111827` (Card surface)
  - `--warmup-border-hairline`: `rgba(255, 255, 255, 0.08)`
  - `--warmup-accent-cyan`: `#06b6d4` (Telemetry accent)
  - `--warmup-accent-emerald`: `#10b981` (Health confirmation)
  - `--warmup-accent-amber`: `#f59e0b` (Degradation warning)
  - `--warmup-accent-rose`: `#ef4444` (Reputation critical)
- **Inline SVG Asset Inlining:** Icons, charts, status glyphs, and geo-paths are embedded directly in the component template markup. No external image fetches are permitted.

---

## 3. Strict Coding Conventions & Architecture Invariants

All code examples, contracts, structs, and implementations across Sequence 341 strictly follow these core invariants:

1. **Error Return Sovereignty (No Exits in Leaf Functions):**
   - Leaf functions returning `error` or `result.Result[T]` must return errors up the stack directly.
   - Calling `os.Exit` or `cliexit.HandleError` inside business logic or storage layers is strictly forbidden.
   - Root command dispatchers alone intercept returned errors and determine exit codes.
2. **Positive Boolean Naming Invariant:**
   - All boolean variables, struct fields, JSON keys, and CLI flags MUST start with `is` or `has` ONLY.
   - Prefixes such as `can`, `should`, `was`, `will`, `allow`, or `enable` are completely banned.
   - Negative booleans (e.g., `isNotReady`, `hasNoErrors`, `isDisabled`) are strictly forbidden.
3. **100% Relative Path Discipline:**
   - All path references in specifications, documents, and code MUST be repository-relative (e.g., `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`).
   - Absolute filesystem paths, Windows drive letters (such as drive letter prefixes), and home-directory tildes (`~`) are strictly prohibited.
4. **Three-Tier Split-DB Architecture:**
   - System state, warmup schedules, and Bayesian priors reside in `repodb/warmup.db`.
   - Command dispatch logs and security compliance events reside in `repodb/audit.db`.
   - Configuration files reside in `.gitmap/config.json`.
5. **Function Size & Spacing Standards:**
   - Functions are limited to <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses and early returns are mandatory).
   - Blank line after closing braces `}` and blank line before `return` statements.

---

## 4. Split-DB Storage Schema Foundations

The underlying data model for Sequence 341 leverages SQLite with WAL (Write-Ahead Logging) mode enabled for high concurrency:

```sql
-- Schema migration for warmup.db
CREATE TABLE IF NOT EXISTS BayesianRampPrior (
    priorId TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    ispClass TEXT NOT NULL,
    alphaShape REAL NOT NULL DEFAULT 1.0,
    betaShape REAL NOT NULL DEFAULT 1.0,
    currentStep INTEGER NOT NULL DEFAULT 1,
    maxDailyVolume INTEGER NOT NULL DEFAULT 500,
    isActive INTEGER NOT NULL DEFAULT 1,
    hasConverged INTEGER NOT NULL DEFAULT 0,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS InboxPlacementTest (
    testId TEXT PRIMARY KEY,
    campaignId TEXT NOT NULL,
    targetDomain TEXT NOT NULL,
    seedAddress TEXT NOT NULL,
    ispProvider TEXT NOT NULL,
    placementFolder TEXT NOT NULL,
    latencyMilliseconds INTEGER NOT NULL,
    isDelivered INTEGER NOT NULL DEFAULT 0,
    hasAuthPassed INTEGER NOT NULL DEFAULT 0,
    hasSpamFlag INTEGER NOT NULL DEFAULT 0,
    testedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ReputationAuditRecord (
    recordId TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    spfRecord TEXT NOT NULL,
    dkimSelector TEXT NOT NULL,
    dmarcPolicy TEXT NOT NULL,
    sndsScore REAL NOT NULL DEFAULT 100.0,
    isDnsPropagated INTEGER NOT NULL DEFAULT 1,
    hasDkimValid INTEGER NOT NULL DEFAULT 1,
    hasSpfAligned INTEGER NOT NULL DEFAULT 1,
    hasDmarcEnforced INTEGER NOT NULL DEFAULT 1,
    auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## 5. Verification Matrix & Quality Gates

| Subsystem | Metric Target | Verification Method |
| :--- | :--- | :--- |
| **Bayesian Ramp Curve** | Inference latency < 2ms | Local bench tests with synthetic 30-day feedback vectors |
| **Geo-Spatial Placement Mesh** | First-render < 16ms (60 FPS) | Shadow DOM paint benchmark on 4K display profile |
| **Reputation Sentinel** | Zero unhandled DNS timeout errors | Mock DNS server injecting 100% packet loss scenarios |
| **CLI Command Suite** | Zero panic exits, full error wrapping | End-to-end testing with `gitmap email warmup ...` commands |
| **Boolean Compliance** | 100% compliance with `is*`/`has*` | Static AST analysis scanner over all structs and schemas |
