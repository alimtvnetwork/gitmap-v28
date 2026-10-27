# Specification 340: Improvements 301 to 330 (Swarm, Rhythms & Fleet Mesh)

- **Sequence:** 340
- **Slug:** `340-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/02-improvements-301-to-330-swarm-rhythms-fleet-mesh.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## Improvement Catalog Summary (301 - 330)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 31** | Items 301–310 | Multi-Agent Collaborative Deliverability Swarm |
| **Cluster 32** | Items 311–320 | Cognitive Sending Rhythms & ISP Fatigue Visualizer |
| **Cluster 33** | Items 321–330 | Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing |

---

# Cluster 31: Multi-Agent Collaborative Deliverability Swarm (Items 301–310)

---

### Improvement 301: Swarm Consensus Topology Canvas Web Component
- **ID:** `WARMUP-SWARM-301`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** Operators have no cohesive visual interface to observe multi-agent warm-up swarms, their communication vectors, vote states, and consensus weights in real time.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-swarm-canvas>`
  - **Shadow DOM Structure:**
    - Host styled with `display: block; position: relative; background: #0b0f19; border-radius: 8px; border: 1px solid rgba(255,255,255,0.08); overflow: hidden;`.
    - SVG interactive network layer (`<svg class="swarm-topology" viewBox="0 0 800 600">`) containing force-directed agent nodes (`<g class="agent-node">`) and consensus links (`<line class="agent-link">`).
    - Animated particle pulses along consensus links visualizing active message passing between agents.
    - Floating agent HUD tooltip showing active role, current vote weight, and consensus status (`.agent-tooltip`).
    - Quick-action footer with toggle controls for swarm auto-pilot (`.btn-autopilot`), force re-vote (`.btn-revote`), and consensus quorum threshold selector.
  - **Zero-CDN Invariant:** Implemented strictly via native Shadow DOM, inline SVGs, and pure CSS keyframe animations (`@keyframes pulseGlow`). Zero external graph libraries (e.g. D3 or Cytoscape).
- **CLI Counterpart Command:**
  - `gitmap warmup swarm topology --format json`
  - `gitmap warmup swarm inspect --node-id "agent-consensus-master"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["canvasId", "swarmId", "isActive", "hasQuorum", "agentCount", "agents"],
    "properties": {
      "canvasId": { "type": "string" },
      "swarmId": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasQuorum": { "type": "boolean" },
      "agentCount": { "type": "integer", "minimum": 1 },
      "agents": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["agentId", "agentRole", "isActive", "hasVoted", "voteWeight"],
          "properties": {
            "agentId": { "type": "string" },
            "agentRole": { "type": "string" },
            "isActive": { "type": "boolean" },
            "hasVoted": { "type": "boolean" },
            "voteWeight": { "type": "number", "minimum": 0, "maximum": 1 }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Renders 10 distinct swarm agent nodes on canvas in <= 15ms.
  - Updates node positions and connection colors when receiving consensus state changes without re-rendering the full DOM.
  - Passes 100% positive boolean validation for `isActive`, `hasQuorum`, and `hasVoted`.

---

### Improvement 302: Autonomous Reputation Guard Agent HUD
- **ID:** `WARMUP-SWARM-302`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** IP and domain reputation drops are detected reactively after blacklisting rather than preemptively by an autonomous sentinel agent capable of intervening immediately.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-reputation>`
  - **Shadow DOM Structure:**
    - Agent status badge with glowing pulsing dot (`.pulse-dot.is-active`).
    - Reputation telemetry matrix displaying SNDS score, Google Postmaster reputation tier (High/Medium/Low/Bad), and Spamhaus status pills.
    - Autonomous action log (`.action-log-stream`) listing recent guard actions: "Auto-throttled pool 2 by 40%", "Rerouted Gmail traffic through backup IP 192.168.1.10".
    - Manual override safety toggle (`<input type="checkbox" class="toggle-override">`) with affirmative label "Autonomous Intervention Active".
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with inline SVG status glyphs and native CSS scrollable log container.
- **CLI Counterpart Command:**
  - `gitmap warmup agent reputation status --domain "outbound.gitmap.dev"`
  - `gitmap warmup agent reputation trigger-check --immediate`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "domainName", "isActive", "hasIntervened", "reputationScore", "postmasterTier"],
    "properties": {
      "agentId": { "type": "string" },
      "domainName": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasIntervened": { "type": "boolean" },
      "reputationScore": { "type": "number", "minimum": 0, "maximum": 100 },
      "postmasterTier": { "type": "string", "enum": ["high", "medium", "low", "bad"] },
      "hasBlacklistHit": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Immediately signals visual warning (rose border `#ef4444`) if `hasBlacklistHit` evaluates to true.
  - Dispatches custom DOM event `agent-intervened` when autonomous mitigation is triggered.
  - CLI command returns typed AppError if domain name is invalid.

---

### Improvement 303: Content Entropy Synthesizer Agent View
- **ID:** `WARMUP-SWARM-303`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** Warm-up emails with repetitive or low-entropy bodies trigger modern statistical spam filters (Bayesian and neural spam classifiers) due to lack of lexical variety.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-entropy>`
  - **Shadow DOM Structure:**
    - Real-time Shannon entropy meter (`.entropy-meter`) with score display (0.0 to 8.0 bits/byte) and dynamic color gradient.
    - Lexical diversity radar chart (`<svg class="lexical-radar">`) breaking down vocabulary uniqueness, n-gram dispersion, and structural variation.
    - Side-by-side generative variant inspector showing seed template and synthetic variant generated by the entropy agent.
    - Synthesizer controls: "Regenerate Variant", "Increase Temperature", "Inject Dynamic Context".
  - **Zero-CDN Invariant:** Native SVG radar graph calculations rendered inline via Vanilla JS. No chart libraries loaded.
- **CLI Counterpart Command:**
  - `gitmap warmup agent entropy evaluate --template-id "tpl-welcome-01"`
  - `gitmap warmup agent entropy synthesize --count 5 --min-entropy 6.5`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "templateId", "isActive", "hasOptimalEntropy", "shannonEntropy", "uniqueRatio"],
    "properties": {
      "agentId": { "type": "string" },
      "templateId": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasOptimalEntropy": { "type": "boolean" },
      "shannonEntropy": { "type": "number", "minimum": 0, "maximum": 8.0 },
      "uniqueRatio": { "type": "number", "minimum": 0, "maximum": 1.0 },
      "hasBayesianFlag": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately renders Shannon entropy metric with 2 decimal precision.
  - Automatically highlights low-entropy sentences in warning amber (`#f59e0b`).
  - Generates valid JSON payloads compatible with `WarmupSwarmAgents` SQLite table.

---

### Improvement 304: Persona Behavior Orchestrator Swarm View
- **ID:** `WARMUP-SWARM-304`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** Artificial warm-up conversations feel robotic and fail to simulate natural human engagement rhythms such as reply delays, forward actions, and thread branching.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-persona>`
  - **Shadow DOM Structure:**
    - Active persona carousel (`.persona-carousel`) showcasing conversational avatars with metadata: role (e.g., "DevOps Lead", "Marketing Director"), engagement propensity, and typing speed profile.
    - Thread simulation graph (`<svg class="thread-graph">`) detailing simulated back-and-forth email interactions across persona pairs.
    - Interactive parameter dials for `replyDelayMin`, `replyDelayMax`, and `isForwardingEnabled`.
    - Live conversational event ticker displaying simulated events: "Persona Alice replied to Bob after 4m 12s", "Bob marked email as Important".
  - **Zero-CDN Invariant:** Pure CSS flex-wrap carousel with hardware-accelerated transitions and self-contained SVG avatars.
- **CLI Counterpart Command:**
  - `gitmap warmup agent persona list --status active`
  - `gitmap warmup agent persona simulate-thread --persona-a "alice" --persona-b "bob"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "personaId", "personaRole", "isActive", "hasReplied", "hasThreadBranching"],
    "properties": {
      "agentId": { "type": "string" },
      "personaId": { "type": "string" },
      "personaRole": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasReplied": { "type": "boolean" },
      "hasThreadBranching": { "type": "boolean" },
      "simulatedReadDurationSeconds": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Renders persona interaction flows with realistic non-uniform delays.
  - Exposes `persona-updated` event upon manual dial adjustment.
  - Enforces positive boolean properties across all persona data records.

---

### Improvement 305: ISP Feedback Loop Negotiator Agent
- **ID:** `WARMUP-SWARM-305`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** When an ISP registers user spam complaints through Feedback Loops (FBLs), traditional systems do not dynamically coordinate with sending agents to instantly suppress target personas.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-fbl>`
  - **Shadow DOM Structure:**
    - FBL incident dashboard showing live complaint rates grouped by ISP (Yahoo, Microsoft, Comcast).
    - Auto-suppression ledger table (`.suppression-ledger`) with instantaneous timestamp, ISP origin, complaining recipient hash, and agent mitigation status.
    - Negotiator response action badge: "Instant Throttle (-50%)", "Domain Shifted", "Cool-off Initiated".
    - Manual release button for suppressed test seeds with safety confirmation modal.
  - **Zero-CDN Invariant:** Fully responsive table with native CSS sticky headers and SVG warning pills.
- **CLI Counterpart Command:**
  - `gitmap warmup agent fbl inspect --isp "yahoo"`
  - `gitmap warmup agent fbl suppress --email "seed-042@yahoo.com" --reason "complaint"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "ispName", "complaintRate", "isActive", "hasAutoSuppressed", "hasExceededThreshold"],
    "properties": {
      "agentId": { "type": "string" },
      "ispName": { "type": "string" },
      "complaintRate": { "type": "number", "minimum": 0, "maximum": 1.0 },
      "isActive": { "type": "boolean" },
      "hasAutoSuppressed": { "type": "boolean" },
      "hasExceededThreshold": { "type": "boolean" },
      "suppressedCount": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Auto-flags ISP complaint rates above 0.001 (0.1%) with red alert border.
  - Triggers instant suppression record creation in `WarmupSwarmAgents` store.
  - Zero external CDN dependencies in styling or icon rendering.

---

### Improvement 306: Bounce Sentinel & SMTP Error Triager Agent
- **ID:** `WARMUP-SWARM-306`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** Transient 4xx SMTP greylisting errors and permanent 5xx hard bounces are often conflated, resulting in unnecessary domain warming halts or catastrophic blacklisting.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-bounce>`
  - **Shadow DOM Structure:**
    - Split error visualizer categorizing live SMTP responses: Hard Bounces (5xx, red), Soft Bounces (4xx, yellow), Protocol Timeouts (orange).
    - Enhanced error code decoder breaking down RFC 3463 enhanced status codes (e.g. `5.7.1 Service unavailable`, `4.4.2 Bad connection`).
    - Smart Retry Countdown timer showing autonomous retry schedule adjusted by exponential jitter backoff.
    - Seed quarantine panel allowing immediate removal of dead or corrupted seed mailboxes.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with CSS grid and inline SVG status indicators.
- **CLI Counterpart Command:**
  - `gitmap warmup agent bounce triage --code "5.7.1"`
  - `gitmap warmup agent bounce quarantine-list`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "smtpCode", "enhancedCode", "isActive", "hasHardBounce", "isQuarantined"],
    "properties": {
      "agentId": { "type": "string" },
      "smtpCode": { "type": "integer", "minimum": 400, "maximum": 599 },
      "enhancedCode": { "type": "string" },
      "errorDiagnostic": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasHardBounce": { "type": "boolean" },
      "isQuarantined": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Correctly categorizes 4xx vs 5xx codes into distinct visual streams.
  - Automatically flags mailbox for quarantine if `hasHardBounce` is true.
  - CLI command returns clean JSON representation of triaged errors.

---

### Improvement 307: DMARC Forensic Forensic Investigator Agent
- **ID:** `WARMUP-SWARM-307`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** DMARC forensic failure reports (RUF) arrive as raw XML/MIME bodies that operators rarely parse, missing spoofing attempts and alignment failures during warm-up.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-dmarc>`
  - **Shadow DOM Structure:**
    - DMARC alignment inspector with SPF and DKIM pass/fail status bars (`.alignment-bar`).
    - Forensic XML/JSON report parser highlighting source IP, disposition (none, quarantine, reject), and selector discrepancies.
    - Threat origin geo-map snippet (`<svg class="geo-points">`) indicating geolocation of unauthorized sending sources.
    - One-click recommendation card: "Update SPF record to include `ip4:198.51.100.0/24`".
  - **Zero-CDN Invariant:** Encapsulated Shadow DOM parser with pure JavaScript XML-to-JSON transformer.
- **CLI Counterpart Command:**
  - `gitmap warmup agent dmarc audit --domain "gitmap.dev"`
  - `gitmap warmup agent dmarc parse-ruf --file "ruf-report.xml"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "domainName", "policyMode", "isActive", "hasSpfAligned", "hasDkimAligned"],
    "properties": {
      "agentId": { "type": "string" },
      "domainName": { "type": "string" },
      "policyMode": { "type": "string", "enum": ["none", "quarantine", "reject"] },
      "isActive": { "type": "boolean" },
      "hasSpfAligned": { "type": "boolean" },
      "hasDkimAligned": { "type": "boolean" },
      "hasSpoofAttempt": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately reports strict vs relaxed alignment status for SPF and DKIM.
  - Visually flags spoofing alerts in high-contrast red badge (`#ef4444`).
  - Strict adherence to positive booleans across all forensic data fields.

---

### Improvement 308: DNS Health & SPF Alignment Auditor Agent
- **ID:** `WARMUP-SWARM-308`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** DNS lookup limits (such as the RFC 7208 10-lookup SPF limit) and propagation delays can silently invalidate domain warm-up sessions across different regional resolvers.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-dns>`
  - **Shadow DOM Structure:**
    - SPF lookup counter meter (`<div class="spf-counter">`) showing current count vs 10-lookup ceiling (e.g. "7 / 10 lookups").
    - Multi-resolver propagation matrix comparing DNS responses across Cloudflare (1.1.1.1), Google (8.8.8.8), Quad9 (9.9.9.9), and OpenDNS.
    - DKIM key length validator displaying key size (1024-bit vs 2048-bit) and syntax validity check.
    - Auto-flattening suggestion preview that generates an optimized single-record SPF string.
  - **Zero-CDN Invariant:** Self-contained SVG dial for lookup counter and pure CSS status grid.
- **CLI Counterpart Command:**
  - `gitmap warmup agent dns verify --domain "outbound.gitmap.dev"`
  - `gitmap warmup agent dns spf-flatten --domain "outbound.gitmap.dev"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "domainName", "spfLookupCount", "isActive", "hasValidSyntax", "hasExceededLimit"],
    "properties": {
      "agentId": { "type": "string" },
      "domainName": { "type": "string" },
      "spfLookupCount": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasValidSyntax": { "type": "boolean" },
      "hasExceededLimit": { "type": "boolean" },
      "hasDkim2048Bit": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Highlights lookup count in red when `hasExceededLimit` is true (>= 10 lookups).
  - Multi-resolver grid renders propagation status across 4 resolvers within 100ms.
  - Zero external CDN dependencies.

---

### Improvement 309: Seed Inbox Coordination Dispatcher Agent
- **ID:** `WARMUP-SWARM-309`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** Seed inboxes become saturated or fatigued when receiving too many test emails from a single domain within a short window, distorting deliverability stats.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-agent-seed>`
  - **Shadow DOM Structure:**
    - Seed pool balance meter showing active vs resting seed mailboxes across major providers (Google, Microsoft, Yahoo, Apple Mail).
    - Intelligent seed allocation slider displaying target distribution percentages per provider.
    - Rest period countdown timer per seed pool indicating remaining time before cooldown completes.
    - "Seed Health Audit" one-click action triggering simulated IMAP fetch and spam folder rescue verification.
  - **Zero-CDN Invariant:** Encapsulated custom element using CSS variable-driven progress bars and inline SVG provider icons.
- **CLI Counterpart Command:**
  - `gitmap warmup agent seed pool-status`
  - `gitmap warmup agent seed dispatch --count 20 --balance-providers`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["agentId", "poolId", "activeSeedCount", "restingSeedCount", "isActive", "hasBalancedSpread"],
    "properties": {
      "agentId": { "type": "string" },
      "poolId": { "type": "string" },
      "activeSeedCount": { "type": "integer", "minimum": 0 },
      "restingSeedCount": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasBalancedSpread": { "type": "boolean" },
      "hasFatiguedSeeds": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Correctly displays resting seed count and flags fatigued pools.
  - Seed dispatch command triggers evenly balanced provider allocation.
  - Emits custom event `seed-dispatched` with target envelope count.

---

### Improvement 310: Swarm Quorum & Consensus Leader Dashboard
- **ID:** `WARMUP-SWARM-310`
- **Cluster:** 31 (Multi-Agent Collaborative Deliverability Swarm)
- **Problem Statement:** Decentralized agents need a unified consensus coordinator to reconcile conflicting advice (e.g. Reputation Guard urging backoff vs Volume Escalator urging speedup).
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-swarm-quorum>`
  - **Shadow DOM Structure:**
    - Quorum state badge (`.quorum-badge`) displaying current leader agent and election epoch number.
    - Consensus agreement bar (`<div class="consensus-progress">`) displaying real-time vote percentage (e.g. 85% agreement across 10 agents).
    - Vote breakdown table listing each subagent's proposed sending rate and reason code.
    - "Manual Leader Override" control switch allowing admin to designate an authoritative agent.
  - **Zero-CDN Invariant:** Standalone custom element with CSS flexbox layout and hardware-accelerated vote meter transitions.
- **CLI Counterpart Command:**
  - `gitmap warmup swarm quorum status`
  - `gitmap warmup swarm elect-leader --agent-id "reputation-guard"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["dashboardId", "leaderAgentId", "electionEpoch", "isActive", "hasQuorum", "consensusPercentage"],
    "properties": {
      "dashboardId": { "type": "string" },
      "leaderAgentId": { "type": "string" },
      "electionEpoch": { "type": "integer", "minimum": 1 },
      "isActive": { "type": "boolean" },
      "hasQuorum": { "type": "boolean" },
      "consensusPercentage": { "type": "number", "minimum": 0, "maximum": 100 },
      "hasSplitVote": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately renders consensus percentage and marks `hasQuorum` as true when consensus >= 66.7%.
  - Renders split-vote warning banner if consensus drops below quorum threshold.
  - Complies with zero-exit and error return sovereignty standards.

---

# Cluster 32: Cognitive Sending Rhythms & ISP Fatigue Visualizer (Items 311–320)

---

### Improvement 311: Circadian Sending Cadence Curve Web Component
- **ID:** `WARMUP-RHYTHM-311`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Robotic flat-rate email delivery (e.g. exactly 1 email every 60 seconds) triggers ISP heuristic pattern-matching algorithms, identifying non-human traffic instantly.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-cadence-curve>`
  - **Shadow DOM Structure:**
    - 24-hour sinusoidal curve graph (`<svg class="circadian-chart" viewBox="0 0 1200 400">`) showing business hour peaks (9 AM–11 AM, 2 PM–4 PM) and nighttime troughs.
    - Current time marker vertical line sliding across curve according to recipient timezone.
    - Adjustable circadian modulation dial controlling amplitude variation from 10% to 90%.
    - Visual overlay showing expected vs actual dispatched email timestamps.
  - **Zero-CDN Invariant:** Self-contained SVG path generator calculating cubic bezier curves from 24 hourly dispatch targets.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm curve --domain "outbound.gitmap.dev" --timezone "America/New_York"`
  - `gitmap warmup rhythm preview-cadence --day 14`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["curveId", "timezone", "isActive", "hasCircadianApplied", "hourlyTargets"],
    "properties": {
      "curveId": { "type": "string" },
      "timezone": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasCircadianApplied": { "type": "boolean" },
      "hourlyTargets": {
        "type": "array",
        "items": { "type": "integer", "minimum": 0 },
        "minItems": 24,
        "maxItems": 24
      },
      "hasPeakHour": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Generates smooth SVG path without jerky line segments across all 24 points.
  - Updates current time indicator accurately to the minute based on provided timezone.
  - Zero external graphing libraries.

---

### Improvement 312: ISP Queue Saturation & Fatigue Gauge Element
- **ID:** `WARMUP-RHYTHM-312`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Mailbox providers maintain connection queues that get saturated before deferrals occur, but operators cannot see queue pressure until errors happen.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-fatigue-gauge>`
  - **Shadow DOM Structure:**
    - Provider fatigue meter (`.fatigue-meter`) showing percentage queue fill (0%–100%) for Google, Microsoft, and Yahoo.
    - Color-coded thermal gradient: Emerald (<50%), Amber (50%–75%), Crimson (>75% fatigue).
    - Deferral risk indicator score showing probability of imminent greylisting.
    - Auto-drain button that triggers momentary sending pause to allow ISP queues to clear.
  - **Zero-CDN Invariant:** Enclosed Shadow DOM with CSS linear gradient tracks and SVG alert icons.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm fatigue --isp "gmail"`
  - `gitmap warmup rhythm drain --isp "outlook" --duration "15m"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["gaugeId", "ispName", "fatiguePercentage", "isActive", "hasHighFatigue", "deferralRiskScore"],
    "properties": {
      "gaugeId": { "type": "string" },
      "ispName": { "type": "string" },
      "fatiguePercentage": { "type": "number", "minimum": 0, "maximum": 100 },
      "isActive": { "type": "boolean" },
      "hasHighFatigue": { "type": "boolean" },
      "deferralRiskScore": { "type": "number", "minimum": 0, "maximum": 1.0 },
      "hasQueueDrained": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Gauge transitions into crimson alert when `fatiguePercentage` exceeds 75%.
  - Dispatches `isp-fatigued` custom event with provider name when threshold crossed.
  - Enforces positive booleans across all state properties.

---

### Improvement 313: Poisson Distributed Micro-Burst Timeline Element
- **ID:** `WARMUP-RHYTHM-313`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Real human outbound email sending follows clustered Poisson arrival times rather than fixed intervals; failing to simulate this makes warmup traffic easily detectable.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-burst-timeline>`
  - **Shadow DOM Structure:**
    - Timeline strip (`<div class="poisson-strip">`) plotting individual dispatch pulses as vertical tick marks along a 60-minute window.
    - Inter-arrival time density histogram (`<svg class="inter-arrival-hist">`) showing exponential distribution of gaps between sends.
    - Lambda ($\lambda$) parameter control slider (adjusting average events per minute from 0.5 to 10.0).
    - "Sample Next 10 Dispatches" simulator button generating probabilistic future send times.
  - **Zero-CDN Invariant:** Built-in Poisson distribution calculation engine inside Vanilla JS Web Component. Zero dependencies.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm poisson-sample --lambda 2.5 --count 20`
  - `gitmap warmup rhythm timeline-inspect --window "1h"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["timelineId", "lambdaValue", "isActive", "hasPoissonDistribution", "eventsDispatched"],
    "properties": {
      "timelineId": { "type": "string" },
      "lambdaValue": { "type": "number", "minimum": 0.1, "maximum": 50.0 },
      "isActive": { "type": "boolean" },
      "hasPoissonDistribution": { "type": "boolean" },
      "eventsDispatched": { "type": "integer", "minimum": 0 },
      "hasClusteredBurst": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Successfully generates tick marks distributed via inverse transform sampling of Poisson process.
  - Histogram accurately represents exponential distribution of inter-arrival delays.
  - Completely responsive within container down to 320px width.

---

### Improvement 314: ISP Greylisting & Backoff Penalty Matrix Element
- **ID:** `WARMUP-RHYTHM-314`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** When greylisted by an ISP (e.g. receiving `451 4.7.1 Try again later`), retrying too fast causes hard blacklisting, while waiting too long stalls the warmup schedule.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-backoff-matrix>`
  - **Shadow DOM Structure:**
    - Backoff stage progression cards (`.backoff-stage`): Stage 1 (5m), Stage 2 (15m), Stage 3 (45m), Stage 4 (2h).
    - Penalty cool-down radial rings indicating active penalty timer per destination domain.
    - Historical greylist encounter graph showing frequency of 4xx responses over time.
    - Immediate "Bypass & Reroute via Secondary Pool" action trigger.
  - **Zero-CDN Invariant:** Encapsulated Shadow DOM styling with CSS animation on countdown rings.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm backoff-status --isp "mimecast"`
  - `gitmap warmup rhythm reset-backoff --domain "outbound.gitmap.dev"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["matrixId", "ispName", "currentBackoffStage", "isActive", "hasActivePenalty", "retryWaitSeconds"],
    "properties": {
      "matrixId": { "type": "string" },
      "ispName": { "type": "string" },
      "currentBackoffStage": { "type": "integer", "minimum": 0, "maximum": 5 },
      "isActive": { "type": "boolean" },
      "hasActivePenalty": { "type": "boolean" },
      "retryWaitSeconds": { "type": "integer", "minimum": 0 },
      "hasClearedPenalty": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Real-time countdown timer decrements accurately without drifting.
  - Marks `hasActivePenalty` as false when countdown reaches zero and dispatches event `penalty-cleared`.
  - Zero external UI libraries.

---

### Improvement 315: Dynamic Throttling Friction Slider Control
- **ID:** `WARMUP-RHYTHM-315`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Operators lack an intuitive physical analogy control to fine-tune transmission friction and enforce gentle drag on outbound queues during sensitive warm-up phases.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-throttle-friction>`
  - **Shadow DOM Structure:**
    - Mechanical friction slider control (`.friction-slider`) styled with grooved track and metal thumb.
    - Numerical friction coefficient readout (e.g., "Friction: 1.45x Delay").
    - Visual spring tension gauge showing artificial latency injected per SMTP socket connection.
    - Preset snap buttons: "Ultra Smooth" (1.0x), "Conservative Drag" (1.5x), "Severe Brake" (3.0x).
  - **Zero-CDN Invariant:** Implemented via standard HTML5 `<input type="range">` customized inside Shadow DOM with pure CSS.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm friction set --coefficient 1.5`
  - `gitmap warmup rhythm friction get`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["controlId", "frictionCoefficient", "isActive", "hasBrakeApplied", "injectedLatencyMs"],
    "properties": {
      "controlId": { "type": "string" },
      "frictionCoefficient": { "type": "number", "minimum": 1.0, "maximum": 5.0 },
      "isActive": { "type": "boolean" },
      "hasBrakeApplied": { "type": "boolean" },
      "injectedLatencyMs": { "type": "integer", "minimum": 0 },
      "hasFrictionOverride": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Slider adjustments trigger immediate DOM update of injected latency preview.
  - Fires custom event `friction-changed` with exact coefficient payload.
  - Strict compliance with positive boolean naming.

---

### Improvement 316: Recipient Timezone Synchronizer Calendar Element
- **ID:** `WARMUP-RHYTHM-316`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Dispatching warmup emails during local recipient middle-of-the-night hours (e.g. 3 AM) lowers open rates and raises spam flags from ISP machine learning models.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-timezone-sync>`
  - **Shadow DOM Structure:**
    - World clock strip displaying 4 core timezone bands (UTC, EST, PST, CET) with local daylight indicators (sun/moon SVGs).
    - Window eligibility matrix highlighting optimal sending slots per timezone (e.g. 8:30 AM to 5:30 PM local).
    - Send prohibition guard banner alerting if active queue contains destinations in night-time windows.
    - Automatic timezone shifter toggle that postpones out-of-window envelopes until target sunrise.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with native JS `Intl.DateTimeFormat` calculations and SVG icons.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm timezone-matrix`
  - `gitmap warmup rhythm queue-shift --auto-daylight`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["calendarId", "targetTimezone", "localHour", "isActive", "hasDaylightWindow", "isShiftEnabled"],
    "properties": {
      "calendarId": { "type": "string" },
      "targetTimezone": { "type": "string" },
      "localHour": { "type": "integer", "minimum": 0, "maximum": 23 },
      "isActive": { "type": "boolean" },
      "hasDaylightWindow": { "type": "boolean" },
      "isShiftEnabled": { "type": "boolean" },
      "hasDelayedEnvelopes": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Highlights timezone column in amber if local hour is between 22:00 and 06:00.
  - Accurately computes local hour offsets using browser internationalization API.
  - Zero third-party date/time libraries (e.g. Moment.js or date-fns).

---

### Improvement 317: Cognitive Jitter & Human Typing Simulator View
- **ID:** `WARMUP-RHYTHM-317`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Outbound SMTP sessions that connect, transmit DATA, and disconnect with millisecond uniformity indicate automated bots to ISP deep packet inspection.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-jitter-simulator>`
  - **Shadow DOM Structure:**
    - SMTP phase duration breakdown showing simulated human composition delay, SMTP greeting latency, and socket hold times.
    - Gaussian jitter curve visualization displaying bell curve of injected delays (mean 120ms, stddev 35ms).
    - Interactive parameter inputs for `meanLatencyMs` and `jitterVariance`.
    - Live socket simulation stream rendering realistic step-by-step SMTP conversation timing.
  - **Zero-CDN Invariant:** Inline SVG normal distribution curve rendered via standard mathematical formulas.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm jitter-config --mean-ms 150 --variance 40`
  - `gitmap warmup rhythm jitter-test --count 10`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["simulatorId", "meanLatencyMs", "jitterVariance", "isActive", "hasJitterEnabled", "isGaussianDistribution"],
    "properties": {
      "simulatorId": { "type": "string" },
      "meanLatencyMs": { "type": "integer", "minimum": 10, "maximum": 5000 },
      "jitterVariance": { "type": "integer", "minimum": 1, "maximum": 1000 },
      "isActive": { "type": "boolean" },
      "hasJitterEnabled": { "type": "boolean" },
      "isGaussianDistribution": { "type": "boolean" },
      "hasPacketInspectionDefense": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Visual Gaussian curve dynamically redraws when mean or variance inputs change.
  - Emits custom event `jitter-adjusted` with updated distribution parameters.
  - Enforces positive booleans across all schema attributes.

---

### Improvement 318: Provider Throttling Threshold Radar Element
- **ID:** `WARMUP-RHYTHM-318`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Each mailbox provider enforces unique, undisclosed hourly connection and volume ceilings; operators cannot see how close they are to provider limits.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-throttle-radar>`
  - **Shadow DOM Structure:**
    - Multi-axis spider/radar chart (`<svg class="radar-chart">`) displaying current usage against estimated capacity across Gmail, Outlook, Yahoo, Proton, and Zoho.
    - Threshold warning perimeter ring marking the 80% caution zone.
    - ISP specific telemetry cards with hourly send counts, connection error rates, and estimated capacity.
    - One-click auto-balance recommendation: "Shift 15% outbound volume from Microsoft to Yahoo".
  - **Zero-CDN Invariant:** Native SVG polygon plotting based on 5-axis polar coordinate calculations. Zero external charting dependencies.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm radar-status`
  - `gitmap warmup rhythm capacity-check --isp "google"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["radarId", "isActive", "hasThresholdExceeded", "providerPoints"],
    "properties": {
      "radarId": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasThresholdExceeded": { "type": "boolean" },
      "providerPoints": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["ispName", "currentVolume", "estimatedCap", "utilizationRatio", "hasOptimalMargin"],
          "properties": {
            "ispName": { "type": "string" },
            "currentVolume": { "type": "integer", "minimum": 0 },
            "estimatedCap": { "type": "integer", "minimum": 1 },
            "utilizationRatio": { "type": "number", "minimum": 0, "maximum": 1.0 },
            "hasOptimalMargin": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Radar chart updates polygon vertices cleanly with animated SVG morphing.
  - Alerts operator with pulsing border if any provider's `utilizationRatio` exceeds 0.85.
  - Passes all positive boolean checks (`isActive`, `hasOptimalMargin`).

---

### Improvement 319: Warming Volume Escalation Escalator View
- **ID:** `WARMUP-RHYTHM-319`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Manual day-by-day volume increases are error-prone and fail to respond when negative reputation signals dictate pausing or stepping down the escalation ramp.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-volume-escalator>`
  - **Shadow DOM Structure:**
    - Multi-step escalator staircase graphic showing 30-day progression steps with daily quota targets.
    - Active day indicator with current progress fill bar showing emails sent today vs daily ceiling.
    - Escalation rules matrix: "Advance condition: Inbox rate > 95%", "Stepback condition: Bounce rate > 2%".
    - Manual step override buttons: "Hold Step", "Step Down", "Advance Early".
  - **Zero-CDN Invariant:** Built with native CSS grid and SVG icons within encapsulated Shadow DOM.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm escalator status --domain "outbound.gitmap.dev"`
  - `gitmap warmup rhythm escalator advance --force`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["escalatorId", "domainName", "currentDay", "maxDays", "isActive", "hasStepAdvanced", "dailyTarget"],
    "properties": {
      "escalatorId": { "type": "string" },
      "domainName": { "type": "string" },
      "currentDay": { "type": "integer", "minimum": 1, "maximum": 90 },
      "maxDays": { "type": "integer", "minimum": 14, "maximum": 90 },
      "isActive": { "type": "boolean" },
      "hasStepAdvanced": { "type": "boolean" },
      "dailyTarget": { "type": "integer", "minimum": 1 },
      "hasStepdownTriggered": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately renders staircase step corresponding to `currentDay`.
  - Dispatches `escalator-step-changed` event when day changes.
  - CLI command safely returns error on invalid domain name.

---

### Improvement 320: Real-Time Rhythm Anomaly Alert Banner Element
- **ID:** `WARMUP-RHYTHM-320`
- **Cluster:** 32 (Cognitive Sending Rhythms & ISP Fatigue Visualizer)
- **Problem Statement:** Sudden anomalies (e.g. sharp increase in delivery latency, spike in 421 deferrals) require instant prominent notification to prevent automated reputation damage.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-rhythm-alert>`
  - **Shadow DOM Structure:**
    - Dismissible glassmorphic top banner (`.alert-banner`) with severity indicator (Info: blue, Warning: amber, Critical: rose).
    - Anomaly detail typography highlighting offending metric (e.g. "Google deferrals surged 300% in last 10m").
    - Action shortcut buttons: "Engage Emergency Brake", "Mute Alert (15m)", "Inspect Logs".
    - Audio alert toggle (`<button class="btn-audio-toggle">`) with embedded SVG speaker glyph for critical alerts.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with CSS keyframe flash effect and pure SVG icons.
- **CLI Counterpart Command:**
  - `gitmap warmup rhythm alert list --unresolved`
  - `gitmap warmup rhythm alert resolve --alert-id "alert-882"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["alertId", "severity", "anomalyDescription", "isActive", "hasResolved", "isEmergencyBrakeEngaged"],
    "properties": {
      "alertId": { "type": "string" },
      "severity": { "type": "string", "enum": ["info", "warning", "critical"] },
      "anomalyDescription": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasResolved": { "type": "boolean" },
      "isEmergencyBrakeEngaged": { "type": "boolean" },
      "timestampEpoch": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Renders critical alert with high-visibility pulsing rose border (`#ef4444`).
  - Emergency brake action immediately emits `emergency-brake-engaged` event.
  - Clean error return handling in Go CLI counterpart.

---

# Cluster 33: Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing (Items 321–330)

---

### Improvement 321: P2P Fleet Mesh Geo-Node Map Web Component
- **ID:** `WARMUP-MESH-321`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Operating warm-up infrastructure across multiple distributed VPS/servers lacks a centralized visual map showing active peer nodes, health, and latency vectors.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-fleet-map>`
  - **Shadow DOM Structure:**
    - World projection vector canvas (`<svg class="world-mesh-map" viewBox="0 0 1000 500">`) with country landmass silhouettes and geo-pinned peer nodes.
    - Curved bezier connection lines between active peer nodes with animated dot pulses representing gossip heartbeats.
    - Interactive node hover card displaying node hostname, public IP, region (e.g. `us-east-1`, `eu-central-1`), and active seed capacity.
    - Fleet health summary overlay displaying total nodes online, mesh connectivity percentage, and average round-trip ping.
  - **Zero-CDN Invariant:** Self-contained low-poly SVG world map embedded directly within the component template. Zero external mapping tiles or APIs.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh nodes list`
  - `gitmap warmup mesh ping-all`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["mapId", "totalNodes", "activeNodes", "isActive", "hasFullMeshConnectivity", "nodes"],
    "properties": {
      "mapId": { "type": "string" },
      "totalNodes": { "type": "integer", "minimum": 1 },
      "activeNodes": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasFullMeshConnectivity": { "type": "boolean" },
      "nodes": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["nodeId", "hostname", "regionCode", "latitude", "longitude", "isActive", "hasHealthyGossip"],
          "properties": {
            "nodeId": { "type": "string" },
            "hostname": { "type": "string" },
            "regionCode": { "type": "string" },
            "latitude": { "type": "number" },
            "longitude": { "type": "number" },
            "isActive": { "type": "boolean" },
            "hasHealthyGossip": { "type": "boolean" },
            "pingLatencyMs": { "type": "integer", "minimum": 0 }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Correctly plots latitude/longitude coordinates to SVG canvas space via Mercator projection.
  - Highlights disconnected or degraded nodes in amber/red with flashing connection lines.
  - Renders completely within 25ms on modern browser hardware.

---

### Improvement 322: Cross-Region Seed Routing Matrix Element
- **ID:** `WARMUP-MESH-322`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Routing seed emails from an EU worker node to a US-based seed inbox introduces unnatural network latency and foreign IP flags from regional ISP filters.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-seed-matrix>`
  - **Shadow DOM Structure:**
    - Origin-to-Destination routing matrix table (`.routing-matrix`) mapping sender node regions to target seed inbox regions.
    - Regional affinity tag pills: "Same-Region Preferred", "Geo-Proximity Locked", "Multi-Hop Allowed".
    - Latency penalty readout showing estimated network hop delay per route.
    - Route optimizer button (`.btn-optimize-routes`) recalculating routing table to minimize cross-continental transit.
  - **Zero-CDN Invariant:** Encapsulated CSS grid layout with native SVG routing arrows.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh routes list`
  - `gitmap warmup mesh route set --origin "eu-west-1" --target "eu-west-1" --affinity "strict"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["matrixId", "totalRoutes", "isActive", "hasGeoAffinityEnforced", "routes"],
    "properties": {
      "matrixId": { "type": "string" },
      "totalRoutes": { "type": "integer", "minimum": 1 },
      "isActive": { "type": "boolean" },
      "hasGeoAffinityEnforced": { "type": "boolean" },
      "routes": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["routeId", "originRegion", "targetRegion", "isActive", "hasLowLatency"],
          "properties": {
            "routeId": { "type": "string" },
            "originRegion": { "type": "string" },
            "targetRegion": { "type": "string" },
            "isActive": { "type": "boolean" },
            "hasLowLatency": { "type": "boolean" },
            "averageLatencyMs": { "type": "integer", "minimum": 0 }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Flags cross-continental routes with warning indicator when `hasLowLatency` is false.
  - Dispatches `route-optimized` event upon optimizer completion.
  - Zero external styles or JS libraries.

---

### Improvement 323: Node Latency & Gossip Ping Monitor Element
- **ID:** `WARMUP-MESH-323`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Fleet nodes must maintain synchronization via gossip protocol; silent network partitions or UDP drops go unnoticed until warm-up states diverge.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-gossip-monitor>`
  - **Shadow DOM Structure:**
    - Live sparkline array showing recent 60-second ping round-trips for each peer in the fleet.
    - Gossip protocol message counter (heartbeats sent/received, delta syncs executed).
    - Partition warning alert banner showing any nodes currently uncontactable.
    - "Force Gossip Sync" button triggering immediate state broadcast across all peers.
  - **Zero-CDN Invariant:** Pure SVG micro-sparklines rendered dynamically from internal state arrays.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh gossip status`
  - `gitmap warmup mesh gossip sync-now`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["monitorId", "clusterEpoch", "isActive", "hasHealthyGossip", "gossipDropRatio"],
    "properties": {
      "monitorId": { "type": "string" },
      "clusterEpoch": { "type": "integer", "minimum": 1 },
      "isActive": { "type": "boolean" },
      "hasHealthyGossip": { "type": "boolean" },
      "gossipDropRatio": { "type": "number", "minimum": 0, "maximum": 1.0 },
      "hasPartitionDetected": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Flags network partitions immediately with flashing amber status dot.
  - Drops sparkline points smoothly off left edge as new telemetry arrives.
  - Strictly follows positive boolean rules (`hasHealthyGossip`, `isActive`).

---

### Improvement 324: WireGuard Tunnel & Peer Encryption Card Element
- **ID:** `WARMUP-MESH-324`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Node-to-node warm-up traffic and seed inbox credentials must be transmitted over encrypted point-to-point tunnels; operators need instant visibility into tunnel handshakes.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-tunnel-card>`
  - **Shadow DOM Structure:**
    - Card container (`.tunnel-card`) with glassmorphic obsidian background and green connection lock icon.
    - WireGuard endpoint status: Interface name (`wg0`), public key fingerprint, peer IP, and allowed IPs.
    - Handshake freshness counter showing seconds since last cryptographic handshake (e.g. "Last handshake: 14s ago").
    - Data transfer meter displaying total Rx/Tx bytes encrypted across the tunnel.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with native CSS flex layout and SVG lock badge.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh tunnel status --peer "node-eu-01"`
  - `gitmap warmup mesh tunnel rekey --peer "node-eu-01"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["tunnelId", "peerNodeId", "interfaceName", "isActive", "hasRecentHandshake", "lastHandshakeSeconds"],
    "properties": {
      "tunnelId": { "type": "string" },
      "peerNodeId": { "type": "string" },
      "interfaceName": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasRecentHandshake": { "type": "boolean" },
      "lastHandshakeSeconds": { "type": "integer", "minimum": 0 },
      "hasEncryptedTraffic": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Automatically marks `hasRecentHandshake` as false if `lastHandshakeSeconds` exceeds 180 seconds.
  - Displays encrypted transfer volume formatted with clean unit suffixes (KB, MB, GB).
  - Pure zero-CDN implementation.

---

### Improvement 325: Seed Inbox Rotation Pool Inspector Element
- **ID:** `WARMUP-MESH-325`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Distributed fleets sharing a seed pool can easily create collision where multiple nodes email the same seed simultaneously, causing artificial spam flags.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-seed-pool>`
  - **Shadow DOM Structure:**
    - Distributed lock status table showing which fleet node currently leases which seed inbox.
    - Lease TTL countdown timer per seed (e.g. "Lease expires in 45s").
    - Collision prevention radar badge confirming zero concurrent sends across the fleet.
    - Seed mailbox recycling filter allowing bulk release of expired or completed test seeds.
  - **Zero-CDN Invariant:** Native CSS grid table with responsive scroll container and SVG lock icons.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh pool leases`
  - `gitmap warmup mesh pool release-lock --seed-id "seed-104"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["poolId", "totalSeeds", "leasedSeeds", "isActive", "hasCollisionAvoidanceActive"],
    "properties": {
      "poolId": { "type": "string" },
      "totalSeeds": { "type": "integer", "minimum": 0 },
      "leasedSeeds": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasCollisionAvoidanceActive": { "type": "boolean" },
      "hasAvailableSeeds": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Displays warning badge if leased seeds equal total seeds (`hasAvailableSeeds` is false).
  - Emits custom event `seed-leased` when a node secures exclusive seed lock.
  - Enforces positive booleans throughout.

---

### Improvement 326: Distributed Rate-Limiting Barrier Element
- **ID:** `WARMUP-MESH-326`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Individual nodes each sending within their local limit can collectively exceed an ISP's global domain rate limit if there is no cross-node distributed token bucket.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-mesh-ratelimit>`
  - **Shadow DOM Structure:**
    - Global token bucket visualizer (`.token-bucket`) showing aggregate tokens available across the entire fleet for each target provider.
    - Local vs Global rate consumption bar comparing node's slice to total fleet burn rate.
    - Token replenishment flow animation indicating rate of refill per second.
    - Manual quota allocation dial allowing operator to reserve tokens for high-priority warmup nodes.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with animated SVG wave fill representing token bucket depth.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh ratelimit status --domain "outbound.gitmap.dev"`
  - `gitmap warmup mesh ratelimit set-limit --isp "yahoo" --global-per-minute 60`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["barrierId", "globalLimitPerMinute", "availableTokens", "isActive", "hasTokensAvailable"],
    "properties": {
      "barrierId": { "type": "string" },
      "globalLimitPerMinute": { "type": "integer", "minimum": 1 },
      "availableTokens": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasTokensAvailable": { "type": "boolean" },
      "hasRateLimitExceeded": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Token bucket graphic updates fill level smoothly on each token consume event.
  - Automatically blocks outbound dispatch when `hasTokensAvailable` becomes false.
  - Complies with zero-exit error architecture in Go CLI.

---

### Improvement 327: Multi-Cloud Seed Transit Latency Gauge
- **ID:** `WARMUP-MESH-327`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Fleets deployed across heterogeneous cloud providers (AWS, GCP, Hetzner, DigitalOcean) encounter variable inter-cloud transit latencies that distort delivery metrics.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-transit-gauge>`
  - **Shadow DOM Structure:**
    - Multi-cloud transit arc gauge comparing inter-cloud latency benchmarks across AWS EC2, GCP Compute Engine, and Hetzner bare-metal nodes.
    - Latency tier classification pills: "Optimal (<25ms)", "Acceptable (25-80ms)", "Elevated (>80ms)".
    - Inter-cloud route selection toggle prioritizing low-latency transit paths for time-critical IMAP verification.
    - Diagnostic trace-route button displaying hop-by-hop latency breakdowns.
  - **Zero-CDN Invariant:** Pure SVG semicircular gauge with dynamic gradient track and needle indicator.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh transit benchmark`
  - `gitmap warmup mesh transit trace --from "aws-us-east" --to "gcp-eu-west"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["gaugeId", "originCloud", "destinationCloud", "latencyMs", "isActive", "hasOptimalTransit"],
    "properties": {
      "gaugeId": { "type": "string" },
      "originCloud": { "type": "string" },
      "destinationCloud": { "type": "string" },
      "latencyMs": { "type": "number", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasOptimalTransit": { "type": "boolean" },
      "hasElevatedLatency": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Moves needle smoothly to target millisecond position using CSS transition.
  - Correctly marks `hasOptimalTransit` as true when latency < 25ms.
  - Zero external CDN scripts or stylesheet links.

---

### Improvement 328: Peer Failover & Split-Brain Guardian View
- **ID:** `WARMUP-MESH-328`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** In network partitions, multiple isolated sub-clusters could each elect a leader and duplicate warmup schedules, severely damaging sender reputation through double-sending.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-splitbrain-guard>`
  - **Shadow DOM Structure:**
    - Quorum health indicator showing active majority calculation (e.g. "5 of 7 nodes visible = Majority Quorum Secured").
    - Split-brain defensive lock banner: "Safe Mode: Minority cluster isolated from sending".
    - Peer witness node status list showing third-party arbiter nodes (e.g. DNS canary, S3 witness).
    - Emergency cluster rejoin button forcing reconciliation of parted nodes once connectivity resumes.
  - **Zero-CDN Invariant:** Self-contained Shadow DOM with status color badges and SVG shield icon.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh splitbrain status`
  - `gitmap warmup mesh splitbrain reconcile`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["guardId", "clusterSize", "visibleNodes", "isActive", "hasMajorityQuorum", "isSafeModeEngaged"],
    "properties": {
      "guardId": { "type": "string" },
      "clusterSize": { "type": "integer", "minimum": 1 },
      "visibleNodes": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasMajorityQuorum": { "type": "boolean" },
      "isSafeModeEngaged": { "type": "boolean" },
      "hasSplitBrainRisk": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Engages safe mode (`isSafeModeEngaged` = true) and displays red shield banner when `hasMajorityQuorum` is false.
  - Dispatches custom event `split-brain-defended` when minority cluster halts sending.
  - Pure positive boolean data representation.

---

### Improvement 329: Cross-Region Envelope Dispatch Tracker Element
- **ID:** `WARMUP-MESH-329`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Operators cannot trace individual warm-up email envelopes as they are queued on node A, routed through peer node B for IP warmth, and verified by seed node C.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-dispatch-tracker>`
  - **Shadow DOM Structure:**
    - Multi-hop envelope journey breadcrumb (`.journey-breadcrumb`): Origin -> Route Transit -> Outbound IP Dispatch -> Seed Inbox Arrival -> IMAP Rescue.
    - Timestamped transit hop log with individual duration per step.
    - Envelope metadata inspector displaying message ID, envelope sender, DKIM signature status, and recipient seed.
    - Live delivery verification badge showing "Delivered to Inbox" vs "Rescued from Spam".
  - **Zero-CDN Invariant:** Enclosed Shadow DOM with SVG breadcrumb connectors and status badges.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh trace-envelope --msg-id "msg-2026-10-8812"`
  - `gitmap warmup mesh recent-dispatches --limit 10`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["trackerId", "messageId", "senderDomain", "isActive", "hasDeliveredToInbox", "isRescuedFromSpam"],
    "properties": {
      "trackerId": { "type": "string" },
      "messageId": { "type": "string" },
      "senderDomain": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasDeliveredToInbox": { "type": "boolean" },
      "isRescuedFromSpam": { "type": "boolean" },
      "totalTransitTimeMs": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Correctly renders each step of multi-hop delivery path with accurate duration.
  - Clearly displays rescue badge if `isRescuedFromSpam` is true.
  - Strict compliance with zero-CDN and positive boolean requirements.

---

### Improvement 330: Fleet Consensus Audit Log Scroller Element
- **ID:** `WARMUP-MESH-330`
- **Cluster:** 33 (Multi-Node P2P Fleet Mesh & Cross-Region Seed Routing)
- **Problem Statement:** Fleet decisions (node joins, seed lease reallocations, rate adjustments) are distributed across multiple nodes; operators lack an immutable, aggregated audit log viewer.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-mesh-auditlog>`
  - **Shadow DOM Structure:**
    - Infinite virtualized log scroller (`.audit-log-scroller`) streaming cryptographic consensus log events from `warmup.db`.
    - Filter bar with checkboxes for event categories: "Node Membership", "Seed Leasing", "Quorum Votes", "Rate Limits".
    - Cryptographic hash verification chip on each log entry showing hash chain integrity status (green lock).
    - Quick-export button (`.btn-export-log`) generating JSON/CSV audit ledger directly in-browser.
  - **Zero-CDN Invariant:** Built with native virtual scrolling technique in Vanilla JS and pure CSS scrollbars. Zero external table or log libraries.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh audit-log --limit 50 --format json`
  - `gitmap warmup mesh audit-verify --verify-chain`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["scrollerId", "totalLogEntries", "isActive", "hasValidHashChain", "entries"],
    "properties": {
      "scrollerId": { "type": "string" },
      "totalLogEntries": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasValidHashChain": { "type": "boolean" },
      "entries": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["logId", "eventType", "originNodeId", "entryHash", "isActive", "hasVerifiedSignature"],
          "properties": {
            "logId": { "type": "string" },
            "eventType": { "type": "string" },
            "originNodeId": { "type": "string" },
            "entryHash": { "type": "string" },
            "isActive": { "type": "boolean" },
            "hasVerifiedSignature": { "type": "boolean" },
            "timestampEpoch": { "type": "integer", "minimum": 0 }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Smoothly scrolls through 500+ log entries using virtualized viewport rendering without frame drops.
  - Accurately renders green cryptographic verification lock when `hasValidHashChain` is true.
  - Returns structured AppError in Go CLI counterpart if hash chain verification fails.
