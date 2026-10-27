# Specification 341: Improvements 401 to 430 (Advanced Clustering & Analytics)

- **Sequence:** 341
- **Slug:** `341-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/02-improvements-401-to-430-advanced-clustering-and-analytics.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## Improvement Catalog Summary (401 - 430)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 41** | Items 401–410 | Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner |
| **Cluster 42** | Items 411–420 | Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map |
| **Cluster 43** | Items 421–430 | Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer |

---

# Cluster 41: Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner (Items 401–410)

---

### Improvement 401: Bayesian Prior Distribution Curve Web Component
- **ID:** `WARMUP-ANALYTICS-401`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Render interactive Beta-Binomial prior and posterior warmup volume distribution curves with sub-second parameter re-calculation in an offline Web Component.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-bayesian-curve>`
  - **Lipgloss Terminal Palette:** Base `#0b0f19`, Curve Line `#06b6d4` (Cyan), Area Fill `rgba(6,182,212,0.15)`, Credible Interval `#10b981` (Emerald).
  - **Shadow DOM Layout:** Encapsulated shadow root (`mode: 'open'`). Top toolbar displays shape parameters `Alpha` and `Beta`, current step count, and credible interval boundary badges (`.ci-badge.is-active`). Inline SVG canvas (`viewBox="0 0 600 240"`) renders the probability density function (PDF) spline, dynamically updating as drag handles adjust prior belief weights.
  - **Zero-CDN Invariant:** Native SVG math calculation using cubic Bezier approximations directly in vanilla ECMAScript. Zero external charting dependencies.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup bayesian curve --domain "outbound.gitmap.dev" --isp "google" --format json
  gitmap email warmup bayesian curve tune --domain "outbound.gitmap.dev" --alpha 4.5 --beta 1.2
  ```
  - `--domain`: FQDN of the outbound sending domain (string).
  - `--isp`: Mailbox provider class (`google`, `microsoft`, `yahoo`, `generic`).
  - `--alpha`: Alpha shape parameter representing successful deliveries (float, >0).
  - `--beta`: Beta shape parameter representing throttle/bounce feedback (float, >0).
  - `--format`: Output formatter (`json`, `text`, `table`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type BayesianCurveState struct {
      DomainName       string  `json:"domainName"`
      IspClass         string  `json:"ispClass"`
      AlphaShape       float64 `json:"alphaShape"`
      BetaShape        float64 `json:"betaShape"`
      ConfidenceRating float64 `json:"confidenceRating"`
      IsActive         bool    `json:"isActive"`
      HasConverged     bool    `json:"hasConverged"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS BayesianCurveModel (
      modelId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      ispClass TEXT NOT NULL,
      alphaShape REAL NOT NULL,
      betaShape REAL NOT NULL,
      confidenceRating REAL NOT NULL,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasConverged INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Hooked directly into `pkg/warmup/bayesian` dispatcher; emits state transitions over GitMap internal event bus topic `warmup.bayesian.curve.updated`.

---

### Improvement 402: Stochastic Gradient Descent Ramp Velocity Tuner
- **ID:** `WARMUP-ANALYTICS-402`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Dynamically tune the daily volume step multiplier via stochastic gradient descent based on real-time soft-bounce backpressure signals.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-gradient-tuner>`
  - **Lipgloss Terminal Palette:** Gradient track `#111827`, Learning Rate Dial `#8b5cf6` (Violet), Loss Indicator `#f59e0b` (Amber).
  - **Shadow DOM Layout:** Real-time mini loss curve chart tracking loss minimization across epochs. Step multiplier gauge (`1.10x` to `2.00x`) with pulsating boundary indicator (`.indicator.has-converged`).
  - **Zero-CDN Invariant:** Self-contained CSS grid layout with inline SVG coordinate plots.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup gradient status --domain "outbound.gitmap.dev"
  gitmap email warmup gradient calibrate --domain "outbound.gitmap.dev" --learning-rate 0.015 --epochs 50
  ```
  - `--domain`: Sending domain to calibrate.
  - `--learning-rate`: Step size for weight adjustments (float, default 0.01).
  - `--epochs`: Optimization iterations over historical telemetry (int).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type GradientTuningRecord struct {
      DomainName       string  `json:"domainName"`
      StepMultiplier   float64 `json:"stepMultiplier"`
      LearningRate     float64 `json:"learningRate"`
      CurrentLoss      float64 `json:"currentLoss"`
      IterationCount   int     `json:"iterationCount"`
      IsOptimizing     bool    `json:"isOptimizing"`
      HasPlateaued     bool    `json:"hasPlateaued"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS GradientTunerState (
      tunerId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL UNIQUE,
      stepMultiplier REAL NOT NULL,
      learningRate REAL NOT NULL,
      currentLoss REAL NOT NULL,
      iterationCount INTEGER NOT NULL,
      isOptimizing INTEGER NOT NULL DEFAULT 0,
      hasPlateaued INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Invoked during the daily schedule compilation stage in `pkg/warmup/scheduler`; calculates tomorrow's dispatch allocation safely without human intervention.

---

### Improvement 403: Posterior Credible Interval Boundary Slider
- **ID:** `WARMUP-ANALYTICS-403`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Provide interactive visual bounds (90%, 95%, 99% credible intervals) allowing operators to set aggressive or defensive risk margins for sending ramp schedules.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-credible-interval>`
  - **Lipgloss Terminal Palette:** Lower Bound `#ef4444` (Rose), Median Path `#06b6d4` (Cyan), Upper Bound `#10b981` (Emerald).
  - **Shadow DOM Layout:** Range slider with tick markers at 90%, 95%, and 99%. A shaded probability mass region visually expands or contracts on the SVG canvas upon thumb motion.
  - **Zero-CDN Invariant:** Native range inputs with custom shadow CSS pseudo-elements (`::-webkit-slider-thumb`) and CSS custom variables.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup interval set --domain "outbound.gitmap.dev" --confidence-level 0.95
  gitmap email warmup interval inspect --domain "outbound.gitmap.dev"
  ```
  - `--domain`: Target domain.
  - `--confidence-level`: Float threshold between 0.80 and 0.99.
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type CredibleIntervalConfig struct {
      DomainName      string  `json:"domainName"`
      ConfidenceLevel float64 `json:"confidenceLevel"`
      LowerBoundVol   int     `json:"lowerBoundVol"`
      UpperBoundVol   int     `json:"upperBoundVol"`
      IsStrictGuard   bool    `json:"isStrictGuard"`
      HasUserOverride bool    `json:"hasUserOverride"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS CredibleIntervalSetting (
      settingId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      confidenceLevel REAL NOT NULL,
      lowerBoundVol INTEGER NOT NULL,
      upperBoundVol INTEGER NOT NULL,
      isStrictGuard INTEGER NOT NULL DEFAULT 1,
      hasUserOverride INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Guards maximum ceiling calculation inside `pkg/warmup/engine`; caps daily step increments to `UpperBoundVol`.

---

### Improvement 404: Likelihood Shift Monte Carlo Simulator
- **ID:** `WARMUP-ANALYTICS-404`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Run 10,000 Monte Carlo forward simulations predicting deliverability trajectories and flagging insolvency risks across 30-day warmup windows.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-monte-carlo-sim>`
  - **Lipgloss Terminal Palette:** Simulation Fan `#1f2937` (Charcoal overlay), P50 Line `#10b981` (Emerald), P05 Failure Boundary `#ef4444` (Rose).
  - **Shadow DOM Layout:** Multi-path trajectory bundle displaying 100 sample paths overlaying the median projection. Summary badge displaying "Ramp Solvency: 99.4%".
  - **Zero-CDN Invariant:** Canvas element rendered via 2D context using pure client-side deterministic PRNG (Mersenne Twister in JS).
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup simulate monte-carlo --domain "outbound.gitmap.dev" --iterations 10000 --days 30
  ```
  - `--iterations`: Number of simulated trajectories (default: 10000).
  - `--days`: Projection horizon in days (default: 30).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type MonteCarloOutcome struct {
      DomainName       string  `json:"domainName"`
      SimulatedDays    int     `json:"simulatedDays"`
      SolvencyRate     float64 `json:"solvencyRate"`
      P05VolumeLimit   int     `json:"p05VolumeLimit"`
      P50VolumeTarget  int     `json:"p50VolumeTarget"`
      P95VolumeCeiling int     `json:"p95VolumeCeiling"`
      IsHighConfidence bool    `json:"isHighConfidence"`
      HasRiskAnomaly   bool    `json:"hasRiskAnomaly"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS MonteCarloSimulationRun (
      runId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      simulatedDays INTEGER NOT NULL,
      solvencyRate REAL NOT NULL,
      p05VolumeLimit INTEGER NOT NULL,
      p50VolumeTarget INTEGER NOT NULL,
      p95VolumeCeiling INTEGER NOT NULL,
      isHighConfidence INTEGER NOT NULL,
      hasRiskAnomaly INTEGER NOT NULL,
      completedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** CLI batch command outputs ANSI sparklines in terminal stdout; reports failure probability to the health monitor.

---

### Improvement 405: ISP Backpressure Feedback Friction Gauge
- **ID:** `WARMUP-ANALYTICS-405`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Quantify and display ISP friction index (0.0 to 1.0) derived from SMTP 4xx deferral responses, connection drops, and recipient rate limits.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-friction-gauge>`
  - **Lipgloss Terminal Palette:** Arc Background `#111827`, Low Friction `#10b981`, Moderate Friction `#f59e0b`, Critical `#ef4444`.
  - **Shadow DOM Layout:** Semi-circular SVG radial gauge with animated needle damping. Numerical display shows friction coefficient with sub-text "Optimal Delivery Velocity".
  - **Zero-CDN Invariant:** Native SVG transform rotations driven by CSS variables (`--needle-angle: 45deg;`).
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup friction inspect --domain "outbound.gitmap.dev" --isp "microsoft"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type FrictionMetrics struct {
      DomainName        string  `json:"domainName"`
      IspProvider       string  `json:"ispProvider"`
      FrictionIndex     float64 `json:"frictionIndex"`
      DeferralCount     int     `json:"deferralCount"`
      IsThrottled       bool    `json:"isThrottled"`
      HasBackpressure   bool    `json:"hasBackpressure"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS IspFrictionSnapshot (
      snapshotId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      ispProvider TEXT NOT NULL,
      frictionIndex REAL NOT NULL,
      deferralCount INTEGER NOT NULL,
      isThrottled INTEGER NOT NULL DEFAULT 0,
      hasBackpressure INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Directly informs the concurrency throttling layer in `pkg/warmup/smtp/pool`.

---

### Improvement 406: Convergence Acceleration Trigger Terminal HUD
- **ID:** `WARMUP-ANALYTICS-406`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Display an interactive terminal HUD indicating when a domain qualifies for fast-tracked warmup acceleration due to impeccable reputation telemetry.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-convergence-hud>`
  - **Lipgloss Terminal Palette:** Border `#06b6d4`, Fast-Track Badge `#10b981`, Metrics Label `#94a3b8`.
  - **Terminal Render:** Double-bordered Lipgloss box with ASCII status stars, qualifying milestones, and affirmative confirmation prompts.
  - **Zero-CDN Invariant:** Fully responsive ANSI terminal stream format with HTML Shadow DOM parity fallback.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup accelerate evaluate --domain "outbound.gitmap.dev"
  gitmap email warmup accelerate activate --domain "outbound.gitmap.dev" --has-ack
  ```
  - `--has-ack`: Affirmative boolean flag confirming operator authorization.
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type ConvergenceStatus struct {
      DomainName          string  `json:"domainName"`
      ConsecutiveCleanDays int     `json:"consecutiveCleanDays"`
      PlacementRate       float64 `json:"placementRate"`
      IsEligible          bool    `json:"isEligible"`
      HasActiveBoost      bool    `json:"hasActiveBoost"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS ConvergenceAcceleration (
      accelerateId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL UNIQUE,
      consecutiveCleanDays INTEGER NOT NULL,
      placementRate REAL NOT NULL,
      isEligible INTEGER NOT NULL DEFAULT 0,
      hasActiveBoost INTEGER NOT NULL DEFAULT 0,
      activatedAt DATETIME
  );
  ```
- **Antigravity / GitMap Integration Point:** Automatically updates the dispatch step incrementor in `pkg/warmup/engine`.

---

### Improvement 407: Historical Ramp Trajectory Overlay Comparator
- **ID:** `WARMUP-ANALYTICS-407`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Overlay current warmup trajectory curves against historical successful warmups across the fleet to benchmark progression velocity.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-trajectory-comparator>`
  - **Lipgloss Terminal Palette:** Target Trace `#06b6d4` (Bold), Historical Fleet Median `#6b7280` (Dashed), Anomaly Bands `#f59e0b`.
  - **Shadow DOM Layout:** Multi-trace SVG line chart with hoverable scrubber points showing volume on Day N for both target and cohort average.
  - **Zero-CDN Invariant:** SVG path generator embedded directly within Web Component methods.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup trajectory compare --domain "outbound.gitmap.dev" --cohort "saas-transactional"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type TrajectoryComparison struct {
      DomainName       string    `json:"domainName"`
      CohortName       string    `json:"cohortName"`
      CurrentDay       int       `json:"currentDay"`
      VelocityDeltaPct float64   `json:"velocityDeltaPct"`
      IsPacingOptimal  bool      `json:"isPacingOptimal"`
      HasCohortMatch   bool      `json:"hasCohortMatch"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS TrajectoryComparisonLog (
      comparisonId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      cohortName TEXT NOT NULL,
      currentDay INTEGER NOT NULL,
      velocityDeltaPct REAL NOT NULL,
      isPacingOptimal INTEGER NOT NULL DEFAULT 1,
      hasCohortMatch INTEGER NOT NULL DEFAULT 1,
      evaluatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Queries fleet statistics from `repodb/warmup.db` to calculate cohort baselines.

---

### Improvement 408: Reputation Loss Penalty Decay Inspector
- **ID:** `WARMUP-ANALYTICS-408`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Visualize the half-life decay curve of historical bounce and spam complaint penalties applied to the Bayesian prior.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-penalty-decay>`
  - **Lipgloss Terminal Palette:** Penalty Bars `#ef4444`, Residual Weight `#f59e0b`, Restored Velocity `#10b981`.
  - **Shadow DOM Layout:** Exponential decay curve showing how an incident from 5 days ago dampens over time with half-life $t_{1/2} = 7$ days.
  - **Zero-CDN Invariant:** Native SVG exponential curve plotting without math library dependencies.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup penalty decay --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type PenaltyDecayState struct {
      DomainName       string  `json:"domainName"`
      InitialPenalty   float64 `json:"initialPenalty"`
      CurrentWeight    float64 `json:"currentWeight"`
      HalfLifeDays     int     `json:"halfLifeDays"`
      IsFullyForgiven  bool    `json:"isFullyForgiven"`
      HasActivePenalty bool    `json:"hasActivePenalty"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS PenaltyDecayRecord (
      penaltyId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      initialPenalty REAL NOT NULL,
      currentWeight REAL NOT NULL,
      halfLifeDays INTEGER NOT NULL DEFAULT 7,
      isFullyForgiven INTEGER NOT NULL DEFAULT 0,
      hasActivePenalty INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Injected into `pkg/warmup/bayesian` when updating the Alpha/Beta shapes.

---

### Improvement 409: Seed Feedback Dispersion Entropy Visualizer
- **ID:** `WARMUP-ANALYTICS-409`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Calculate and display Shannon entropy of seed inbox delivery feedback across mailbox providers to detect inconsistent ISP categorization.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-dispersion-entropy>`
  - **Lipgloss Terminal Palette:** High Entropy `#ef4444` (Chaotic/Inconsistent), Low Entropy `#10b981` (Uniform/Predictable), Metric Label `#94a3b8`.
  - **Shadow DOM Layout:** Circular entropy bubble chart where bubble size correlates with variance and color reflects predictability score.
  - **Zero-CDN Invariant:** SVG circle packing algorithm implemented in standard JS.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup entropy dispersion --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type DispersionEntropy struct {
      DomainName        string  `json:"domainName"`
      EntropyBits       float64 `json:"entropyBits"`
      MaxPossibleBits   float64 `json:"maxPossibleBits"`
      IsUniformDelivery bool    `json:"isUniformDelivery"`
      HasHighEntropy    bool    `json:"hasHighEntropy"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS DispersionEntropyAudit (
      auditId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      entropyBits REAL NOT NULL,
      maxPossibleBits REAL NOT NULL,
      isUniformDelivery INTEGER NOT NULL DEFAULT 1,
      hasHighEntropy INTEGER NOT NULL DEFAULT 0,
      auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Consumed by deliverability audit command `gitmap email audit dispersion`.

---

### Improvement 410: Adaptive Gradient Step Quorum Synchronizer
- **ID:** `WARMUP-ANALYTICS-410`
- **Cluster:** 41 (Bayesian Predictive Ramp Curve Visualizer & ML Gradient Tuner)
- **Goal:** Synchronize Bayesian gradient updates across multi-node peer swarms via decentralized Raft-lite consensus before applying new volume limits.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-gradient-quorum>`
  - **Lipgloss Terminal Palette:** Quorum Banner `#8b5cf6`, Node Affirmative `#10b981`, Node Pending `#94a3b8`.
  - **Shadow DOM Layout:** Visual quorum ring showing connected nodes, committed gradient vectors, and consensus round sequence numbers.
  - **Zero-CDN Invariant:** Native SVG arcs and CSS ring layout.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup quorum status --domain "outbound.gitmap.dev"
  gitmap email warmup quorum sync --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type QuorumSyncState struct {
      DomainName      string   `json:"domainName"`
      QuorumRound     int64    `json:"quorumRound"`
      ConfirmedNodes  []string `json:"confirmedNodes"`
      IsQuorumReached bool     `json:"isQuorumReached"`
      HasSyncConflict bool     `json:"hasSyncConflict"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS GradientQuorumLog (
      roundId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      quorumRound INTEGER NOT NULL,
      nodeCount INTEGER NOT NULL,
      isQuorumReached INTEGER NOT NULL DEFAULT 0,
      hasSyncConflict INTEGER NOT NULL DEFAULT 0,
      committedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Integrates with P2P fleet networking in `pkg/cluster/mesh`.

---

# Cluster 42: Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map (Items 411–420)

---

### Improvement 411: Global Inbox Placement Vector Map Web Component
- **ID:** `WARMUP-ANALYTICS-411`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Display a zero-CDN interactive global vector map plotting real-time seed inbox placement across geographic MX data centers.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-geo-map>`
  - **Lipgloss Terminal Palette:** Continents `#1f2937`, Background `#0b0f19`, Placement Nodes: Inbox `#10b981`, Spam `#ef4444`, Promo `#06b6d4`.
  - **Shadow DOM Layout:** Lightweight vector SVG map projecting worldwide topology. Interactive dots pulse at major cloud data center regions (us-east-1, eu-central-1, ap-southeast-1). Clicking a node opens an inspection panel showing folder placement breakdown and ISP distribution.
  - **Zero-CDN Invariant:** Encapsulated inline SVG path strings of global landmasses (18KB simplified GeoJSON-to-SVG paths) stored inside the Web Component class.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup geo map --domain "outbound.gitmap.dev" --format ascii
  gitmap email warmup geo inspect --region "eu-west-1"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type GeoPlacementNode struct {
      RegionCode       string  `json:"regionCode"`
      Latitude         float64 `json:"latitude"`
      Longitude        float64 `json:"longitude"`
      InboxPct         float64 `json:"inboxPct"`
      SpamPct          float64 `json:"spamPct"`
      IsHealthyRegion  bool    `json:"isHealthyRegion"`
      HasActiveIncident bool   `json:"hasActiveIncident"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS GeoPlacementReport (
      reportId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      regionCode TEXT NOT NULL,
      latitude REAL NOT NULL,
      longitude REAL NOT NULL,
      inboxPct REAL NOT NULL,
      spamPct REAL NOT NULL,
      isHealthyRegion INTEGER NOT NULL DEFAULT 1,
      hasActiveIncident INTEGER NOT NULL DEFAULT 0,
      reportedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Provides the visual rendering target for seed polling results ingested via `pkg/warmup/seed`.

---

### Improvement 412: Sub-Second ISP Round-Trip Latency Heatmap
- **ID:** `WARMUP-ANALYTICS-412`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Track and visualize SMTP TCP handshake and DATA command round-trip latency to each major ISP edge to identify regional peering bottlenecks.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-isp-latency>`
  - **Lipgloss Terminal Palette:** Fast (<50ms) `#10b981`, Normal (50-150ms) `#06b6d4`, Slow (>150ms) `#f59e0b`, Degraded (>400ms) `#ef4444`.
  - **Shadow DOM Layout:** Heatmap matrix with hours of day on X-axis and ISPs (Google, O365, Yahoo, iCloud) on Y-axis. Cells shaded proportionally to median RTT.
  - **Zero-CDN Invariant:** Native CSS grid with dynamic `style="background-color: var(--cell-heat);"` bindings.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup latency matrix --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type IspLatencyRecord struct {
      IspName           string `json:"ispName"`
      Region            string `json:"region"`
      HandshakeMs       int    `json:"handshakeMs"`
      DataAckMs         int    `json:"dataAckMs"`
      IsAcceptableDelay bool   `json:"isAcceptableDelay"`
      HasSpikeAnomaly   bool   `json:"hasSpikeAnomaly"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS IspLatencyLog (
      logId TEXT PRIMARY KEY,
      ispName TEXT NOT NULL,
      region TEXT NOT NULL,
      handshakeMs INTEGER NOT NULL,
      dataAckMs INTEGER NOT NULL,
      isAcceptableDelay INTEGER NOT NULL DEFAULT 1,
      hasSpikeAnomaly INTEGER NOT NULL DEFAULT 0,
      measuredAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Fed directly from SMTP pool connection benchmarks in `pkg/warmup/smtp`.

---

### Improvement 413: MX Connection Edge TLS Handshake Inspector
- **ID:** `WARMUP-ANALYTICS-413`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Verify and visually document TLS 1.3 cipher suite negotiation, cipher strength, and certificate authority validity for each MX gateway.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-tls-inspector>`
  - **Lipgloss Terminal Palette:** Cipher TLS 1.3 `#10b981`, TLS 1.2 `#06b6d4`, Deprecated `#ef4444`, Key Length `#8b5cf6`.
  - **Shadow DOM Layout:** Security audit card listing cipher suite (e.g. `TLS_AES_256_GCM_SHA384`), curve (X25519), and certificate expiration countdown.
  - **Zero-CDN Invariant:** Pure HTML/CSS badge indicators.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup tls inspect --mx "aspmx.l.google.com"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type TlsAuditReport struct {
      MxHost          string `json:"mxHost"`
      ProtocolVersion string `json:"protocolVersion"`
      CipherSuite     string `json:"cipherSuite"`
      CertExpiryDays  int    `json:"certExpiryDays"`
      IsTlsModern     bool   `json:"isTlsModern"`
      HasValidChain   bool   `json:"hasValidChain"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS MxTlsAudit (
      auditId TEXT PRIMARY KEY,
      mxHost TEXT NOT NULL,
      protocolVersion TEXT NOT NULL,
      cipherSuite TEXT NOT NULL,
      certExpiryDays INTEGER NOT NULL,
      isTlsModern INTEGER NOT NULL DEFAULT 1,
      hasValidChain INTEGER NOT NULL DEFAULT 1,
      auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Evaluated in pre-flight delivery checks in `pkg/warmup/smtp/tls`.

---

### Improvement 414: ASN Peering BGP Route Stability Indicator
- **ID:** `WARMUP-ANALYTICS-414`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Detect autonomous system number (ASN) route flaps and BGP instability between outbound sending nodes and receiving mailbox provider clusters.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-asn-stability>`
  - **Lipgloss Terminal Palette:** Stable `#10b981`, Route Flapping `#ef4444`, Hop Count `#94a3b8`.
  - **Shadow DOM Layout:** Visual traceroute hop chain with node badges indicating ASN, hop count, and latency variance.
  - **Zero-CDN Invariant:** Horizontal flexbox flow with inline SVG interconnect arrows.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup bgp check --target-asn 15169 --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type AsnRouteStatus struct {
      TargetAsn       int     `json:"targetAsn"`
      AsnOrg          string  `json:"asnOrg"`
      HopCount        int     `json:"hopCount"`
      JitterMs        float64 `json:"jitterMs"`
      IsRouteStable   bool    `json:"isRouteStable"`
      HasBgpFlapping  bool    `json:"hasBgpFlapping"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS AsnRouteTelemetry (
      telemetryId TEXT PRIMARY KEY,
      targetAsn INTEGER NOT NULL,
      asnOrg TEXT NOT NULL,
      hopCount INTEGER NOT NULL,
      jitterMs REAL NOT NULL,
      isRouteStable INTEGER NOT NULL DEFAULT 1,
      hasBgpFlapping INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Used by route selection algorithms in `pkg/cluster/mesh`.

---

### Improvement 415: Seed Placement Polar Coordinate Radar
- **ID:** `WARMUP-ANALYTICS-415`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Display deliverability health across 6 discrete ISP axes (Google, Microsoft, Yahoo, Apple, Proton, Enterprise) simultaneously on a radar chart.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-radar-placement>`
  - **Lipgloss Terminal Palette:** Radar Grid `#1f2937`, Polygon Fill `rgba(16,185,129,0.2)`, Polygon Outline `#10b981`.
  - **Shadow DOM Layout:** Regular hexagonal SVG radar canvas with polygon points mapped to placement percentages (0-100%).
  - **Zero-CDN Invariant:** Trigonometric SVG polygon vertices calculated directly using native JavaScript `Math.cos` and `Math.sin`.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup radar status --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type RadarDeliverabilityScore struct {
      DomainName     string             `json:"domainName"`
      ProviderScores map[string]float64 `json:"providerScores"`
      OverallScore   float64            `json:"overallScore"`
      IsFullyBalanced bool              `json:"isFullyBalanced"`
      HasWeakProvider bool              `json:"hasWeakProvider"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS RadarScoreSnapshot (
      snapshotId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      providerScoresJson TEXT NOT NULL,
      overallScore REAL NOT NULL,
      isFullyBalanced INTEGER NOT NULL DEFAULT 1,
      hasWeakProvider INTEGER NOT NULL DEFAULT 0,
      evaluatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Renders high-level scorecard in `gitmap email warmup status`.

---

### Improvement 416: Real-time Seed Delivery Pulse Stream
- **ID:** `WARMUP-ANALYTICS-416`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Stream live seed inbox confirmations in real time with sub-second WebSocket or local IPC notifications.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-pulse-stream>`
  - **Lipgloss Terminal Palette:** Event Item `#111827`, Status OK `#10b981`, Timestamp `#64748b`.
  - **Shadow DOM Layout:** Scrolling live event feed with smooth CSS entry transitions. Each item highlights recipient domain, folder, and confirmation latency.
  - **Zero-CDN Invariant:** Vanilla custom element maintaining an internal ring buffer of the 50 most recent events.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup stream follow --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type SeedDeliveryEvent struct {
      EventId        string `json:"eventId"`
      DomainName     string `json:"domainName"`
      SeedRecipient  string `json:"seedRecipient"`
      TargetFolder   string `json:"targetFolder"`
      ElapsedSeconds float64 `json:"elapsedSeconds"`
      IsInbox        bool   `json:"isInbox"`
      HasSpamTag     bool   `json:"hasSpamTag"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS SeedEventJournal (
      eventId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      seedRecipient TEXT NOT NULL,
      targetFolder TEXT NOT NULL,
      elapsedSeconds REAL NOT NULL,
      isInbox INTEGER NOT NULL DEFAULT 1,
      hasSpamTag INTEGER NOT NULL DEFAULT 0,
      occurredAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Pushes events directly to terminal listeners and connected web visualizers.

---

### Improvement 417: Geographic Seed Node Availability Grid
- **ID:** `WARMUP-ANALYTICS-417`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Display operational readiness and health status of geo-distributed seed mailbox listener daemons across multi-cloud regions.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-seed-nodes>`
  - **Lipgloss Terminal Palette:** Active Node `#10b981`, Syncing Node `#06b6d4`, Offline `#ef4444`.
  - **Shadow DOM Layout:** 4x4 matrix grid showing seed nodes, IMAP connection state, and unread verification queues.
  - **Zero-CDN Invariant:** Native CSS Grid with `@container` responsive rules.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup seed nodes list
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type SeedNodeState struct {
      NodeId          string `json:"nodeId"`
      CloudRegion     string `json:"cloudRegion"`
      ImapHost        string `json:"imapHost"`
      ActiveSeeds     int    `json:"activeSeeds"`
      IsNodeConnected bool   `json:"isNodeConnected"`
      HasSyncError    bool   `json:"hasSyncError"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS SeedNodeRegistry (
      nodeId TEXT PRIMARY KEY,
      cloudRegion TEXT NOT NULL,
      imapHost TEXT NOT NULL,
      activeSeeds INTEGER NOT NULL,
      isNodeConnected INTEGER NOT NULL DEFAULT 1,
      hasSyncError INTEGER NOT NULL DEFAULT 0,
      lastHeartbeat DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Managed by the seed coordinator in `pkg/warmup/seed/manager`.

---

### Improvement 418: ISP Soft-Bounce Code Hex Bin Visualizer
- **ID:** `WARMUP-ANALYTICS-418`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Cluster and map SMTP 4xx and 5xx rejection strings into an interactive hexagonal bin chart to identify emerging spam filter patterns.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-hex-bounce>`
  - **Lipgloss Terminal Palette:** Hex Border `#1f2937`, High Density `#ef4444`, Low Density `#1e293b`.
  - **Shadow DOM Layout:** SVG honeycomb tessellation where hexagon size and intensity represent error code frequency (e.g. `4.7.1 Rate Limit Exceeded`).
  - **Zero-CDN Invariant:** Hexagon coordinate math generated via inline math formulas.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup bounce hexbin --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type BounceClusterBin struct {
      SmtpCode        string `json:"smtpCode"`
      PatternSignature string `json:"patternSignature"`
      OccurrenceCount int    `json:"occurrenceCount"`
      IsBlocking       bool   `json:"isBlocking"`
      HasRateLimitFlag bool   `json:"hasRateLimitFlag"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS BounceClusterLog (
      clusterId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      smtpCode TEXT NOT NULL,
      patternSignature TEXT NOT NULL,
      occurrenceCount INTEGER NOT NULL,
      isBlocking INTEGER NOT NULL DEFAULT 0,
      hasRateLimitFlag INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Analyzed by the error classification subsystem in `pkg/warmup/bounce`.

---

### Improvement 419: MX Lookup Latency & DNS Anycast Validator
- **ID:** `WARMUP-ANALYTICS-419`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Measure DNS resolution time for MX records across multiple anycast resolvers (1.1.1.1, 8.8.8.8, 9.9.9.9) and detect DNS poisoning or cache anomalies.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-dns-latency>`
  - **Lipgloss Terminal Palette:** Fast Resolver `#10b981`, Slow Resolver `#f59e0b`, Outlier `#ef4444`.
  - **Shadow DOM Layout:** Bar chart comparing resolution times across 5 global DNS providers with TTL badges.
  - **Zero-CDN Invariant:** Pure CSS flex-based progress bars.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup dns bench --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type DnsResolverBench struct {
      ResolverIp     string `json:"resolverIp"`
      ProviderName   string `json:"providerName"`
      LatencyMs      int    `json:"latencyMs"`
      ResolvedMx     string `json:"resolvedMx"`
      IsRecordMatch  bool   `json:"isRecordMatch"`
      HasCacheStale  bool   `json:"hasCacheStale"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS DnsBenchmarkLog (
      benchId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      resolverIp TEXT NOT NULL,
      providerName TEXT NOT NULL,
      latencyMs INTEGER NOT NULL,
      resolvedMx TEXT NOT NULL,
      isRecordMatch INTEGER NOT NULL DEFAULT 1,
      hasCacheStale INTEGER NOT NULL DEFAULT 0,
      testedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Integrated into DNS health check diagnostics in `pkg/warmup/dns`.

---

### Improvement 420: Cross-Region Seed Placement Divergence Delta
- **ID:** `WARMUP-ANALYTICS-420`
- **Cluster:** 42 (Real-time Inbox Placement Geo-Spatial Mesh & ISP Latency Map)
- **Goal:** Identify geo-divergence where emails hit the inbox in North America but land in the spam folder in Europe, flagging regional IP reputation divergence.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-region-delta>`
  - **Lipgloss Terminal Palette:** Divergence Alert `#ef4444`, Harmonized `#10b981`, Delta Label `#f59e0b`.
  - **Shadow DOM Layout:** Side-by-side comparison cards comparing NA, EU, and APAC inbox rates with delta percentage badges.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with clean grid styling.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup geo divergence --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type RegionDivergence struct {
      DomainName        string  `json:"domainName"`
      PrimaryRegion     string  `json:"primaryRegion"`
      SecondaryRegion   string  `json:"secondaryRegion"`
      PlacementDeltaPct float64 `json:"placementDeltaPct"`
      IsSignificantGap  bool    `json:"isSignificantGap"`
      HasDivergenceWarn bool    `json:"hasDivergenceWarn"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS RegionDivergenceAlert (
      alertId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      primaryRegion TEXT NOT NULL,
      secondaryRegion TEXT NOT NULL,
      placementDeltaPct REAL NOT NULL,
      isSignificantGap INTEGER NOT NULL DEFAULT 0,
      hasDivergenceWarn INTEGER NOT NULL DEFAULT 0,
      detectedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Triggers warning alerts in the central notification stream.

---

# Cluster 43: Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer (Items 421–430)

---

### Improvement 421: Dynamic DMARC Compliance Enforcement Graph
- **ID:** `WARMUP-ANALYTICS-421`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Render a visual dependency graph showing alignment between RFC 7489 DMARC policy (`p=none`, `p=quarantine`, `p=reject`), SPF authentication, and DKIM signature validity.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-dmarc-graph>`
  - **Lipgloss Terminal Palette:** Policy Strict `#10b981`, Quarantine `#f59e0b`, None/Permissive `#ef4444`, Link Connect `#3b82f6`.
  - **Shadow DOM Layout:** Three connected circular nodes (SPF, DKIM, DMARC) with animated flow pulses. If both SPF and DKIM pass alignment, the central DMARC shield glows green; otherwise, warning halos indicate failure points.
  - **Zero-CDN Invariant:** Inline SVG tree graph generated deterministically without D3.js.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup auth graph --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type DmarcAlignmentGraph struct {
      DomainName      string `json:"domainName"`
      DmarcPolicy     string `json:"dmarcPolicy"`
      SpfAlignment    string `json:"spfAlignment"`
      DkimAlignment   string `json:"dkimAlignment"`
      IsDmarcEnforced bool   `json:"isDmarcEnforced"`
      HasAlignmentPass bool  `json:"hasAlignmentPass"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS DmarcGraphState (
      stateId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL UNIQUE,
      dmarcPolicy TEXT NOT NULL,
      spfAlignment TEXT NOT NULL,
      dkimAlignment TEXT NOT NULL,
      isDmarcEnforced INTEGER NOT NULL DEFAULT 1,
      hasAlignmentPass INTEGER NOT NULL DEFAULT 1,
      verifiedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Enforces security precondition gate before dispatching warm-up batches in `pkg/warmup/auth`.

---

### Improvement 422: Autonomous DKIM Key Rotation Timeline Scheduler
- **ID:** `WARMUP-ANALYTICS-422`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Manage zero-downtime DKIM 2048-bit RSA/Ed25519 key rotation with automated selector switching, DNS verification, and scheduled retirement.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-dkim-rotation>`
  - **Lipgloss Terminal Palette:** Active Selector `#10b981`, Pending DNS `#06b6d4`, Retiring `#94a3b8`.
  - **Shadow DOM Layout:** Interactive Gantt-style timeline showing active key lifespan, overlap propagation windows (minimum 72 hours), and upcoming rotation dates.
  - **Zero-CDN Invariant:** CSS grid timeline with native tooltip hovers.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup dkim rotate --domain "outbound.gitmap.dev" --selector "gm202610" --key-type "ed25519"
  ```
  - `--selector`: DNS selector string (e.g. `gm202610`).
  - `--key-type`: Cryptographic algorithm (`rsa2048`, `ed25519`).
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type DkimKeyLifecycle struct {
      DomainName      string `json:"domainName"`
      CurrentSelector string `json:"currentSelector"`
      NextSelector    string `json:"nextSelector"`
      OverlapHours    int    `json:"overlapHours"`
      IsKeyActive     bool   `json:"isKeyActive"`
      HasDnsPublished bool   `json:"hasDnsPublished"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS DkimRotationSchedule (
      rotationId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      currentSelector TEXT NOT NULL,
      nextSelector TEXT NOT NULL,
      overlapHours INTEGER NOT NULL DEFAULT 72,
      isKeyActive INTEGER NOT NULL DEFAULT 1,
      hasDnsPublished INTEGER NOT NULL DEFAULT 0,
      scheduledAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Orchestrates key updates with DNS provider APIs and internal signer modules.

---

### Improvement 423: SPF 10-DNS-Lookup Void Traversal Analyzer
- **ID:** `WARMUP-ANALYTICS-423`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Recursively resolve and count DNS lookups in SPF records, warning when approaching or exceeding the hard RFC 7208 10-lookup limit.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-spf-lookup>`
  - **Lipgloss Terminal Palette:** Safe (<=7) `#10b981`, Warning (8-9) `#f59e0b`, Exceeded (>=10) `#ef4444`.
  - **Shadow DOM Layout:** Stepped meter showing `N / 10 Lookups Used` with tree breakdown of included domains (`include:_spf.google.com`, `include:mailgun.org`).
  - **Zero-CDN Invariant:** Native SVG progress arc and expandable HTML `<details>` elements.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup spf inspect --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type SpfLookupAudit struct {
      DomainName       string   `json:"domainName"`
      LookupCount      int      `json:"lookupCount"`
      IncludedDomains  []string `json:"includedDomains"`
      IsLimitCompliant bool     `json:"isLimitCompliant"`
      HasVoidLookup    bool     `json:"hasVoidLookup"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS SpfLookupRecord (
      auditId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      lookupCount INTEGER NOT NULL,
      includedDomainsJson TEXT NOT NULL,
      isLimitCompliant INTEGER NOT NULL DEFAULT 1,
      hasVoidLookup INTEGER NOT NULL DEFAULT 0,
      auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Pre-dispatch validator prevents delivery failures caused by `PermError`.

---

### Improvement 424: DMARC Forensic RUF/RUA Aggregate Parse Tree
- **ID:** `WARMUP-ANALYTICS-424`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Ingest and parse XML DMARC aggregate reports (RUA), rendering an interactive drilldown breakdown of failing source IP addresses.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-dmarc-rua>`
  - **Lipgloss Terminal Palette:** Clean IP `#10b981`, Spoofed/Failed IP `#ef4444`, Org Name `#06b6d4`.
  - **Shadow DOM Layout:** Sortable data table with expandable rows showing detailed XML failure samples, auth results, and PTR records.
  - **Zero-CDN Invariant:** Vanilla JS DOM table sorting and filtering without third-party datatables.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup dmarc parse --file "rua-report.xml" --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type DmarcAggregateRecord struct {
      DomainName     string `json:"domainName"`
      ReportingOrg   string `json:"reportingOrg"`
      SourceIP       string `json:"sourceIP"`
      MessageCount   int    `json:"messageCount"`
      Disposition    string `json:"disposition"`
      IsAuthAligned  bool   `json:"isAuthAligned"`
      HasSpfFailure  bool   `json:"hasSpfFailure"`
      HasDkimFailure bool   `json:"hasDkimFailure"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS DmarcAggregateLog (
      recordId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      reportingOrg TEXT NOT NULL,
      sourceIP TEXT NOT NULL,
      messageCount INTEGER NOT NULL,
      disposition TEXT NOT NULL,
      isAuthAligned INTEGER NOT NULL DEFAULT 1,
      hasSpfFailure INTEGER NOT NULL DEFAULT 0,
      hasDkimFailure INTEGER NOT NULL DEFAULT 0,
      receivedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Feeds real failure data into the reputation sentinel database.

---

### Improvement 425: Multi-RBL Real-time Blacklist Scanner Matrix
- **ID:** `WARMUP-ANALYTICS-425`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Concurrently scan outbound IPs against 60+ DNS-based real-time blacklists (Spamhaus, Barracuda, SORBS, Invaluement) with sub-second aggregate completion.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-rbl-matrix>`
  - **Lipgloss Terminal Palette:** Clean `#10b981`, Listed `#ef4444`, Query Timeout `#f59e0b`.
  - **Shadow DOM Layout:** Grid of mini badges, each representing an RBL provider. Green checkmarks indicate clean status; red alert badges trigger an immediate modal alert.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with inline SVG status icons.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup rbl scan --ip "198.51.100.25"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type RblScanSummary struct {
      IpAddress      string   `json:"ipAddress"`
      TotalScanned   int      `json:"totalScanned"`
      CleanCount     int      `json:"cleanCount"`
      ListedCount    int      `json:"listedCount"`
      ListedProviders []string `json:"listedProviders"`
      IsCleanReputation bool  `json:"isCleanReputation"`
      HasBlacklistHit   bool  `json:"hasBlacklistHit"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS RblScanResult (
      scanId TEXT PRIMARY KEY,
      ipAddress TEXT NOT NULL,
      totalScanned INTEGER NOT NULL,
      cleanCount INTEGER NOT NULL,
      listedCount INTEGER NOT NULL,
      listedProvidersJson TEXT NOT NULL,
      isCleanReputation INTEGER NOT NULL DEFAULT 1,
      hasBlacklistHit INTEGER NOT NULL DEFAULT 0,
      scannedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Immediately halts warmup dispatch if `hasBlacklistHit` evaluates to `true`.

---

### Improvement 426: Reverse DNS (PTR) Forward-Confirmed Parity Inspector
- **ID:** `WARMUP-ANALYTICS-426`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Verify that the IP's PTR hostname resolves back to the same IP (FCrDNS), a strict delivery requirement for Gmail and Yahoo.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-fcrdns-checker>`
  - **Lipgloss Terminal Palette:** Matched FCrDNS `#10b981`, Mismatch Alert `#ef4444`, Hostname `#06b6d4`.
  - **Shadow DOM Layout:** Two-way verification diagram connecting `IP -> PTR -> A -> IP` with green verification checkmarks.
  - **Zero-CDN Invariant:** Inline SVG bidirectional arrows and styled cards.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup dns ptr --ip "198.51.100.25" --expected-host "mail.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type FcrdnsReport struct {
      IpAddress       string `json:"ipAddress"`
      PtrHostname     string `json:"ptrHostname"`
      ForwardResolved string `json:"forwardResolved"`
      IsMatched       bool   `json:"isMatched"`
      HasValidPtr     bool   `json:"hasValidPtr"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS FcrdnsAuditLog (
      auditId TEXT PRIMARY KEY,
      ipAddress TEXT NOT NULL,
      ptrHostname TEXT NOT NULL,
      forwardResolved TEXT NOT NULL,
      isMatched INTEGER NOT NULL DEFAULT 1,
      hasValidPtr INTEGER NOT NULL DEFAULT 1,
      verifiedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Executed during IP onboarding and daily health auditing in `pkg/warmup/network`.

---

### Improvement 427: BIMI SVG Tiny-PS Validation & Certificate Visualizer
- **ID:** `WARMUP-ANALYTICS-427`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Audit Brand Indicators for Message Identification (BIMI) DNS records, SVG Tiny P/S compliance, and Verified Mark Certificates (VMC).
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-bimi-validator>`
  - **Lipgloss Terminal Palette:** Valid Logo `#10b981`, SVG Error `#ef4444`, VMC Certified `#8b5cf6`.
  - **Shadow DOM Layout:** Brand preview card displaying the validated SVG logo in a realistic Gmail/Apple Mail inbox row mockup.
  - **Zero-CDN Invariant:** SVG parsed locally in Shadow DOM; zero external image proxying.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup bimi check --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type BimiValidationResult struct {
      DomainName     string `json:"domainName"`
      LogoSvgUrl     string `json:"logoSvgUrl"`
      VmcCertUrl     string `json:"vmcCertUrl"`
      IsSvgTinyCompliant bool `json:"isSvgTinyCompliant"`
      HasVmcCertified    bool `json:"hasVmcCertified"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS BimiValidationLog (
      validationId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL UNIQUE,
      logoSvgUrl TEXT NOT NULL,
      vmcCertUrl TEXT NOT NULL,
      isSvgTinyCompliant INTEGER NOT NULL DEFAULT 1,
      hasVmcCertified INTEGER NOT NULL DEFAULT 0,
      validatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Informs domain trust level inside `pkg/warmup/bimi`.

---

### Improvement 428: SNDS / Google Postmaster Autonomous Scraper Sentinel
- **ID:** `WARMUP-ANALYTICS-428`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Ingest Microsoft Smart Network Data Services (SNDS) and Google Postmaster Tools API telemetry, mapping domain reputation tiers into local SQLite metrics.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-postmaster-sync>`
  - **Lipgloss Terminal Palette:** High Rep `#10b981`, Medium `#06b6d4`, Low `#f59e0b`, Bad `#ef4444`.
  - **Shadow DOM Layout:** Tiered reputation step card showing Spam Rate (<0.10%), Domain Reputation, and User-Reported Spam Volume.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with clean tabular styling.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup postmaster sync --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type PostmasterMetrics struct {
      DomainName       string  `json:"domainName"`
      ReputationTier   string  `json:"reputationTier"`
      UserSpamRatePct  float64 `json:"userSpamRatePct"`
      IsHighReputation bool    `json:"isHighReputation"`
      HasSpamRateAlert bool    `json:"hasSpamRateAlert"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS PostmasterDailyMetric (
      metricId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      reputationTier TEXT NOT NULL,
      userSpamRatePct REAL NOT NULL,
      isHighReputation INTEGER NOT NULL DEFAULT 1,
      hasSpamRateAlert INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Modulates daily ramp step sizing in `pkg/warmup/bayesian`.

---

### Improvement 429: TLS-RPT (RFC 8460) Failure Report Stream Aggregator
- **ID:** `WARMUP-ANALYTICS-429`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Ingest and visualize SMTP TLS Reporting (TLS-RPT) JSON payloads, alerting on MTA-STS negotiation downgrades and TLS intercept attempts.
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-tls-rpt>`
  - **Lipgloss Terminal Palette:** Safe TLS `#10b981`, Intercept Warn `#ef4444`, Downgrade `#f59e0b`.
  - **Shadow DOM Layout:** Event log listing reporting organizations (e.g. `google.com`), total successful sessions, and failure session counters.
  - **Zero-CDN Invariant:** Native DOM event list with CSS keyframe highlighting.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup tlsrpt ingest --file "tlsrpt-google.json"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type TlsRptReport struct {
      DomainName       string `json:"domainName"`
      Provider         string `json:"provider"`
      SuccessCount     int    `json:"successCount"`
      FailureCount     int    `json:"failureCount"`
      IsTlsEnforced    bool   `json:"isTlsEnforced"`
      HasDowngradeWarn bool   `json:"hasDowngradeWarn"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS TlsRptSummary (
      reportId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      provider TEXT NOT NULL,
      successCount INTEGER NOT NULL,
      failureCount INTEGER NOT NULL,
      isTlsEnforced INTEGER NOT NULL DEFAULT 1,
      hasDowngradeWarn INTEGER NOT NULL DEFAULT 0,
      processedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Analyzed in `pkg/warmup/security/tlsrpt`.

---

### Improvement 430: MTA-STS Dynamic Policy Engine & Validator
- **ID:** `WARMUP-ANALYTICS-430`
- **Cluster:** 43 (Autonomous Reputation Drift Sentinel & Dynamic SPF/DKIM/DMARC Visualizer)
- **Goal:** Verify and deploy RFC 8461 Mail Transfer Agent Strict Transport Security (MTA-STS) policy files and DNS TXT markers (`_mta-sts.domain.com`).
- **UI/UX Terminal & Web Specification:**
  - **Tag Name:** `<warmup-mta-sts>`
  - **Lipgloss Terminal Palette:** Enforce Mode `#10b981`, Testing Mode `#06b6d4`, Missing `#ef4444`.
  - **Shadow DOM Layout:** Policy inspector displaying `mode: enforce`, `mx` match patterns, `max_age`, and live HTTPS policy file fetch status.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with clean code-block styling.
- **CLI Command Signature & Flag Definitions:**
  ```bash
  gitmap email warmup mtasts verify --domain "outbound.gitmap.dev"
  ```
- **Data Contract & Split-DB Schema:**
  ```go
  package types

  type MtaStsPolicyStatus struct {
      DomainName     string `json:"domainName"`
      PolicyMode     string `json:"policyMode"`
      MaxAgeSeconds  int    `json:"maxAgeSeconds"`
      IsEnforced     bool   `json:"isEnforced"`
      HasValidHttps  bool   `json:"hasValidHttps"`
  }
  ```
  ```sql
  -- Storage in repodb/warmup.db
  CREATE TABLE IF NOT EXISTS MtaStsAudit (
      policyId TEXT PRIMARY KEY,
      domainName TEXT NOT NULL UNIQUE,
      policyMode TEXT NOT NULL,
      maxAgeSeconds INTEGER NOT NULL,
      isEnforced INTEGER NOT NULL DEFAULT 1,
      hasValidHttps INTEGER NOT NULL DEFAULT 1,
      auditedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **Antigravity / GitMap Integration Point:** Finalizes high-security deliverability verification before production scaling.
