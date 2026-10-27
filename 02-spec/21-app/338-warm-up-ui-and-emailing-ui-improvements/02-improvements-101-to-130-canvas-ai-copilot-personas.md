# Specification 338: Improvements 101 to 130 (Canvas, AI Co-Pilot & Personas)

- **Sequence:** 338
- **Slug:** `338-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/02-improvements-101-to-130-canvas-ai-copilot-personas.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## Improvement Catalog Summary (101 - 130)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 11** | Items 101–110 | Interactive Node Canvas & Delivery Path Routing |
| **Cluster 12** | Items 111–120 | AI Deliverability Co-Pilot & Subject Line Entropy |
| **Cluster 13** | Items 121–130 | Multi-Persona Seed Simulation & Conversational Dynamics |

---

# Cluster 11: Interactive Node Canvas & Delivery Path Routing (Items 101–110)

---

### Improvement 101: Force-Directed Topology Node Canvas
- **ID:** `WARMUP-CANVAS-101`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Operators manage warmup across dozens of sending IP pools and SMTP relays blind to routing bottlenecks, unable to observe topological clustering or physical server saturation in real time.
- **Concrete UI Specification:**
  - **Layout:** An infinite-pan SVG/HTML5 canvas element (`.warmup-canvas-viewport`) with hardware-accelerated pan-zoom transform matrices.
  - **Visual Elements:**
    - Force-directed interactive nodes representing sender domains (`circle.node-domain`), SMTP gateways (`rect.node-gateway`), and ISP destinations (`polygon.node-isp`).
    - Dynamic connection edges rendered via SVG path cubic splines with stroke dash animations representing active throughput (`stroke-dashoffset` keyed to send rate).
    - Node charge physics computed via local pure JavaScript Verlet integration loop with zero external libraries.
    - Hover card revealing instantaneous throughput, queue backlog depth, and connection latency with backdrop blur (`backdrop-filter: blur(8px)`).
  - **Zero-CDN Invariant:** Implemented strictly using vanilla HTML5 Canvas or SVG with CSS variables. Zero dependencies on D3.js or Vis.js.
- **CLI Counterpart Command:**
  - `gitmap warmup canvas render --topology-id "core-mesh" --output ascii`
  - `gitmap warmup canvas nodes list --filter "isActive=true"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["nodeId", "nodeType", "isActive", "hasOverload", "latencyMs", "coordinates"],
    "properties": {
      "nodeId": { "type": "string" },
      "nodeType": { "type": "string", "enum": ["domain", "gateway", "isp", "seed"] },
      "isActive": { "type": "boolean" },
      "hasOverload": { "type": "boolean" },
      "latencyMs": { "type": "number", "minimum": 0 },
      "coordinates": {
        "type": "object",
        "required": ["x", "y"],
        "properties": {
          "x": { "type": "number" },
          "y": { "type": "number" }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Canvas renders 150 concurrent nodes at 60 FPS without DOM frame drops.
  - Double clicking a gateway node opens the detailed connection inspector drawer.
  - `gitmap warmup canvas render` displays compliant ANSI ASCII graph representation.

---

### Improvement 102: Real-Time Bezier Path Packet Animator
- **ID:** `WARMUP-CANVAS-102`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Outbound email batches appear static in status tables; operators cannot visually verify packet flow velocity, queuing delays, or in-flight handoffs between MTA hops.
- **Concrete UI Specification:**
  - **Layout:** Micro-particle overlay layer embedded on top of the node canvas (`svg.packet-animation-layer`).
  - **Visual Elements:**
    - Glowing photon particles (`circle.packet-photon`, `r=3`, `fill: #06b6d4`, glow filter `drop-shadow(0 0 6px #06b6d4)`).
    - Particles traverse SVG paths via native CSS `offset-path: path(...)` and `animation: packetFlow 1.2s cubic-bezier(0.4, 0, 0.2, 1) infinite`.
    - Packet speed scales dynamically in proportion to real-time transactional velocity (messages per second).
    - Failed or bounced packets flash red (`#ef4444`) and splinter into dissipated particle bursts on target collision.
  - **Zero-CDN Invariant:** Driven by native CSS keyframes and SVG stroke animations without external tweening libraries.
- **CLI Counterpart Command:**
  - `gitmap warmup canvas trace-packets --source "mta-01.local" --destination "gmail.com" --follow`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["packetId", "routePath", "velocityRate", "isActive", "hasDelivered"],
    "properties": {
      "packetId": { "type": "string" },
      "routePath": { "type": "string" },
      "velocityRate": { "type": "number" },
      "isActive": { "type": "boolean" },
      "hasDelivered": { "type": "boolean" },
      "failureCode": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Packet particles follow exact Bezier spline geometry between connected nodes.
  - Particle spawn rate smoothly adjusts between 1 to 500 packets/sec based on telemetry state.
  - Memory consumption remains under 45MB during continuous 2-hour stress execution.

---

### Improvement 103: Dynamic SMTP Relay Weight Sliders
- **ID:** `WARMUP-CANVAS-103`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Adjusting traffic distribution across secondary SMTP relays requires editing raw YAML configs and restarting daemons, risking configuration drift and uneven domain warming.
- **Concrete UI Specification:**
  - **Layout:** Floating glassmorphic dock panel (`.relay-balancer-dock`) situated at the base of the topology canvas.
  - **Visual Elements:**
    - High-precision percentage distribution sliders (`input[type="range"].relay-slider`) with cyan fill track and magnetic snapping.
    - Proportional lock toggles (`button.lock-ratio-btn`) to maintain ratios when dragging master volume.
    - Real-time pie preview donut gauge reflecting relative load balancing shares.
    - Visual warning indicator pulsing amber when an uncertified relay receives >25% volume.
  - **Zero-CDN Invariant:** Custom styled pure CSS range inputs with hardware-accelerated linear gradients.
- **CLI Counterpart Command:**
  - `gitmap warmup relay set-weights --relay-id "smtp-out-2" --weight 35 --lock-group "tier-1"`
  - `gitmap warmup relay get-weights --format json`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["relayId", "weightPercentage", "isActive", "hasLockGroup"],
    "properties": {
      "relayId": { "type": "string" },
      "weightPercentage": { "type": "number", "minimum": 0, "maximum": 100 },
      "isActive": { "type": "boolean" },
      "hasLockGroup": { "type": "boolean" },
      "lockGroupId": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Cumulative slider values enforce a strict 100% sum invariant with auto-normalization.
  - Sliders update local SQLite configuration within 12ms of change event.
  - CLI command reflects modified weights immediately without service restart.

---

### Improvement 104: Multi-Hop MTA Connection Inspector
- **ID:** `WARMUP-CANVAS-104`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Delivery path failures often stem from intermediate TLS handshake timeouts or proxy delays that are hidden behind simple delivery failure messages.
- **Concrete UI Specification:**
  - **Layout:** Slide-over drawer component (`.mta-inspector-drawer`) anchored to the right edge with 1px obsidian border.
  - **Visual Elements:**
    - Step-by-step vertical hop ladder (`.hop-ladder-track`) showing Sender -> Milter -> Local MTA -> Outbound Proxy -> Destination MX.
    - TLS cipher badge (`span.badge-cipher`) displaying cipher suite and protocol version (`TLS 1.3 - AES_256_GCM`).
    - Hop duration micro-timelines with sub-millisecond precision color-coded by performance thresholds (<50ms green, 50-200ms yellow, >200ms red).
    - Expandable raw SMTP transcript log viewer with syntax highlighted commands (`EHLO`, `STARTTLS`, `MAIL FROM`, `RCPT TO`).
  - **Zero-CDN Invariant:** Native CSS flexbox layout with custom dark scrollbar rules.
- **CLI Counterpart Command:**
  - `gitmap warmup mta inspect-hop --message-id "msg_984210" --include-transcript`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["messageId", "hopCount", "isSuccessful", "hasTlsHandshake", "hops"],
    "properties": {
      "messageId": { "type": "string" },
      "hopCount": { "type": "integer", "minimum": 1 },
      "isSuccessful": { "type": "boolean" },
      "hasTlsHandshake": { "type": "boolean" },
      "hops": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["hopIndex", "serverHost", "durationMs", "isTlsVerified"],
          "properties": {
            "hopIndex": { "type": "integer" },
            "serverHost": { "type": "string" },
            "durationMs": { "type": "number" },
            "isTlsVerified": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Clicking any active canvas edge opens the inspector drawer populated with the latest hop telemetry.
  - Transcripts safely escape all HTML control characters.
  - Latency bar renders accurately according to the specified 50ms/200ms color thresholds.

---

### Improvement 105: Autonomous SPF/DKIM/DMARC Visual Verifier
- **ID:** `WARMUP-CANVAS-105`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** DNS authentication status is typically diagnosed through external third-party web tools, introducing latency and exposing proprietary sender domain lists.
- **Concrete UI Specification:**
  - **Layout:** Floating security radar badge card (`.dns-auth-matrix-card`) embedded inside each sending domain canvas node.
  - **Visual Elements:**
    - Triple status LED indicators (`.led-auth-pill`) for SPF, DKIM, and DMARC with pulse animation when active.
    - Interactive popover revealing parsed DNS TXT records, key lengths (e.g. 2048-bit RSA vs Ed25519), and DMARC alignment mode (`relaxed` vs `strict`).
    - Alignment failure banner highlighting subdomain mismatches in red callout boxes.
    - One-click "Re-probe DNS" trigger button with spinning sync SVG glyph.
  - **Zero-CDN Invariant:** Pure SVG icons and CSS keyframe animations.
- **CLI Counterpart Command:**
  - `gitmap warmup dns verify-auth --domain "acme-corp.com" --strict-alignment`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["domain", "isSpfValid", "isDkimValid", "isDmarcValid", "hasAlignment"],
    "properties": {
      "domain": { "type": "string" },
      "isSpfValid": { "type": "boolean" },
      "isDkimValid": { "type": "boolean" },
      "isDmarcValid": { "type": "boolean" },
      "hasAlignment": { "type": "boolean" },
      "dmarcPolicy": { "type": "string", "enum": ["none", "quarantine", "reject"] }
    }
  }
  ```
- **Acceptance Criteria:**
  - DNS probe resolves against local authoritative nameservers without external API calls.
  - LED status updates in real-time when DNS updates propagate.
  - Alignment checker correctly distinguishes between organizational root domains and RFC 5322 header subdomains.

---

### Improvement 106: Canvas Geo-Distribution Overlay
- **ID:** `WARMUP-CANVAS-106`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Outbound warm-up bursts originating exclusively from a single cloud availability zone trigger geographic anomalies in recipient ISP anti-abuse systems.
- **Concrete UI Specification:**
  - **Layout:** Dark-vector equirectangular world map projection (`.geo-canvas-layer`) underneath node topology connections.
  - **Visual Elements:**
    - High-contrast SVG geopolitical coastlines rendered in slate tone (`#1e293b`).
    - Origin server pings represented as expanding radar ripple rings (`circle.geo-radar-pulse`).
    - Arced geodesic great-circle flight paths connecting sender data centers (e.g. `us-east-1`, `eu-central-1`) to regional ISP MX clusters.
    - Geographic density heat clouds computed via client-side radial gradient stamping.
  - **Zero-CDN Invariant:** Embedded lightweight geo-path SVG vector data directly in local HTML/CSS bundle.
- **CLI Counterpart Command:**
  - `gitmap warmup geo distribution-map --cluster-id "fleet-alpha" --format terminal-ascii`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["originRegion", "targetRegion", "isActive", "hasGeoBalancing", "trafficWeight"],
    "properties": {
      "originRegion": { "type": "string" },
      "targetRegion": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasGeoBalancing": { "type": "boolean" },
      "trafficWeight": { "type": "number", "minimum": 0, "maximum": 1 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Geodesic arcs compute correct spherical trigonometry curves between coordinate pairs.
  - Map pans and zooms seamlessly synchronized with node canvas transformations.
  - CLI emits clean 80-column ASCII terminal world map with region pin points.

---

### Improvement 107: Throttling & Rate-Limit Backoff Gauge
- **ID:** `WARMUP-CANVAS-107`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** When an ISP issues `421 4.7.0 Try again later` rate-limit responses, operators lack immediate visibility into active backoff cooldown timers and queue choke states.
- **Concrete UI Specification:**
  - **Layout:** Radial cooldown countdown meter (`.backoff-gauge-container`) attached to constrained node edges.
  - **Visual Elements:**
    - Circumferential SVG countdown ring decreasing clockwise (`stroke: #f59e0b; stroke-dasharray: 283;`).
    - Centered monospace digital timer displaying remaining wait period (`04:32 remaining`).
    - Exponential backoff tier indicator (`Tier 1 (2m)` -> `Tier 2 (8m)` -> `Tier 3 (32m)`).
    - Emergency manual override button (`button.override-cooldown-btn`) protected by double-click confirmation.
  - **Zero-CDN Invariant:** Self-contained SVG ring and CSS timer transitions.
- **CLI Counterpart Command:**
  - `gitmap warmup throttle inspect --isp "yahoo.com" --include-cooldown`
  - `gitmap warmup throttle reset-cooldown --isp "yahoo.com" --force`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["ispDomain", "isThrottled", "hasCooldownActive", "remainingSeconds", "backoffTier"],
    "properties": {
      "ispDomain": { "type": "string" },
      "isThrottled": { "type": "boolean" },
      "hasCooldownActive": { "type": "boolean" },
      "remainingSeconds": { "type": "integer", "minimum": 0 },
      "backoffTier": { "type": "integer", "minimum": 1, "maximum": 5 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Countdown timer decrements synchronously with system clock without drift.
  - Edge transitions to amber pulse state whenever throttling status is activated.
  - CLI resets cooldown immediately upon executing override with audit log record.

---

### Improvement 108: Canvas Node Clustering & Domain Aggregator
- **ID:** `WARMUP-CANVAS-108`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Viewing more than 50 individual sending domains simultaneously causes canvas visual clutter, obscuring organizational trends across parent corporate brands.
- **Concrete UI Specification:**
  - **Layout:** Hierarchical bounding convex hulls (`path.cluster-hull`) grouping affiliated child domains.
  - **Visual Elements:**
    - Semi-transparent tinted polygon hulls (`fill: rgba(59, 130, 246, 0.06); stroke: #3b82f6; stroke-dasharray: 4;`).
    - Group header pill badge displaying cluster name and aggregate delivery health score (e.g. `Acme EU Fleet: 98.4%`).
    - Expand/collapse toggle button on cluster header with smooth force-directed re-layout transition.
    - Aggregated bubble view when zoomed out below 40% scale factor.
  - **Zero-CDN Invariant:** Convex hull polygon calculation executed locally via Monotone Chain algorithm in vanilla JS.
- **CLI Counterpart Command:**
  - `gitmap warmup cluster list-groups`
  - `gitmap warmup cluster create --name "Acme EU Fleet" --domain-pattern "*.eu.acme.com"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["groupId", "groupName", "isActive", "hasAggregatedView", "domainCount"],
    "properties": {
      "groupId": { "type": "string" },
      "groupName": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasAggregatedView": { "type": "boolean" },
      "domainCount": { "type": "integer", "minimum": 1 },
      "healthScore": { "type": "number", "minimum": 0, "maximum": 100 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Convex hull accurately wraps all child nodes with minimum 16px visual padding.
  - Zooming out automatically collapses detailed nodes into consolidated cluster pills.
  - Group creation CLI command immediately applies regex pattern matching to available nodes.

---

### Improvement 109: Real-Time Latency & Jitter Heat-Ribbon
- **ID:** `WARMUP-CANVAS-109`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Network latency spikes and TCP socket jitter during SMTP handshakes degrade warm-up ramp pace, but are invisible in summary throughput tables.
- **Concrete UI Specification:**
  - **Layout:** Continuous streaming horizontal SVG sparkline ribbon (`.jitter-heat-ribbon`) pinned beneath the primary canvas toolbar.
  - **Visual Elements:**
    - High-density color-mapped gradient stroke (cool green <30ms, cyan 30-70ms, amber 70-150ms, rose >150ms).
    - Jitter variance zone indicated by semi-transparent envelope area around the median ping line.
    - Monospace readout indicating standard deviation and packet loss percentage.
    - Interactive scrub bar allowing retrospective timeline replay of past 60 minutes.
  - **Zero-CDN Invariant:** Native SVG path points generated directly from in-memory ring buffer.
- **CLI Counterpart Command:**
  - `gitmap warmup latency ribbon --window "1h" --resolution "10s"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["medianLatencyMs", "jitterMs", "isActive", "hasPacketLoss", "samples"],
    "properties": {
      "medianLatencyMs": { "type": "number" },
      "jitterMs": { "type": "number" },
      "isActive": { "type": "boolean" },
      "hasPacketLoss": { "type": "boolean" },
      "packetLossRate": { "type": "number", "minimum": 0, "maximum": 1 },
      "samples": { "type": "array", "items": { "type": "number" } }
    }
  }
  ```
- **Acceptance Criteria:**
  - Ribbon updates at 1Hz frequency without causing redraws on the underlying node canvas.
  - Color gradient dynamically reflects latency spikes exceeding 150ms within one refresh cycle.
  - Memory consumption stays constant by recycling fixed-size 360-sample array.

---

### Improvement 110: Interactive Canvas Snapshot & State Exporter
- **ID:** `WARMUP-CANVAS-110`
- **Cluster:** 11 (Interactive Node Canvas & Delivery Path Routing)
- **Problem Statement:** Support and deliverability consultants cannot easily export visual routing topology states for asynchronous client reporting or compliance archiving.
- **Concrete UI Specification:**
  - **Layout:** Export modal dialog (`.canvas-export-modal`) accessible from top canvas menu.
  - **Visual Elements:**
    - High-resolution SVG and PNG snapshot export options with customizable scale multiplier (1x, 2x, 4x).
    - Watermark toggle embedding timestamp, commit hash, and operator initials.
    - Redaction checkbox (`input#redact-domains`) that obscures sensitive corporate domains into masked tokens (`ac***.com`).
    - One-click copy vector markup to clipboard button with animated checkmark feedback.
  - **Zero-CDN Invariant:** Utilizes native browser Canvas `toBlob` and XMLSerializer without external html2canvas scripts.
- **CLI Counterpart Command:**
  - `gitmap warmup canvas export --format svg --out "./reports/topology-snapshot.svg" --redact-sensitive`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["snapshotId", "timestampIso", "isRedacted", "hasVectorSource", "fileSizeBytes"],
    "properties": {
      "snapshotId": { "type": "string" },
      "timestampIso": { "type": "string" },
      "isRedacted": { "type": "boolean" },
      "hasVectorSource": { "type": "boolean" },
      "fileSizeBytes": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Exported SVG renders cleanly in standalone vector viewers without missing styles.
  - Redaction feature completely strips identifiable customer domains from export bundle.
  - Generated export matches exact visual layout rendered on screen at capture time.

---

# Cluster 12: AI Deliverability Co-Pilot & Subject Line Entropy (Items 111–120)

---

### Improvement 111: Shannon Token Entropy Gauge
- **ID:** `WARMUP-COPILOT-111`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Subject lines with repetitive token patterns or uniform length trigger Bayesian spam filters, yet authors lack instant mathematical entropy feedback while typing.
- **Concrete UI Specification:**
  - **Layout:** Real-time entropy dial gauge (`.entropy-dial-widget`) mounted adjacent to subject input fields.
  - **Visual Elements:**
    - Semicircular SVG gauge arc with needle indicator rotating between 0.0 and 8.0 bits/token.
    - Optimal entropy target zone highlighted in emerald green (`3.8 - 5.4 bits/token`).
    - Monospace numerical display with two decimal precision updating per keystroke (`H = 4.72 bits`).
    - Sub-label badge indicating token distribution status (`Low Variance`, `Optimal Natural Variation`, `High Chaos/Spam-Like`).
  - **Zero-CDN Invariant:** Native SVG arc geometry calculated via trigonometry formulas in vanilla JS.
- **CLI Counterpart Command:**
  - `gitmap warmup entropy calculate --subject "Quick follow up regarding our call" --format json`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["subjectLine", "entropyScore", "isOptimal", "hasHighChaos", "tokenCount"],
    "properties": {
      "subjectLine": { "type": "string" },
      "entropyScore": { "type": "number", "minimum": 0, "maximum": 8 },
      "isOptimal": { "type": "boolean" },
      "hasHighChaos": { "type": "boolean" },
      "tokenCount": { "type": "integer", "minimum": 1 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Needle rotation recalculates within 8ms of typing without input lag.
  - Shannon formula accurately computes `-Σ P(w) log2 P(w)` over white-space delimited tokens.
  - Gauge transitions into amber warning zone when repetitive tokens drop entropy below 2.5.

---

### Improvement 112: Semantic Spam-Trigger Heatmap Highlighter
- **ID:** `WARMUP-COPILOT-112`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Outbound copy frequently contains subtle financial, urgent, or promotional trigger tokens that trigger machine-learning ISP filters without author awareness.
- **Concrete UI Specification:**
  - **Layout:** Rich content-editable draft canvas (`.draft-copilot-canvas`) with inline token highlight overlays.
  - **Visual Elements:**
    - Real-time token underline highlights color-coded by severity category:
      - Urgency triggers (`Act Now`, `Urgent`, `Expires`): Amber underline with tooltip.
      - Financial/Promotional triggers (`Free`, `Guarantee`, `100% Risk Free`): Rose underline.
      - Punctuation/Case anomalies (`FREE!!!`, `CLICK HERE NOW`): Crimson wavy border.
    - Hover card displaying trigger weight and 3 instant one-click replacement alternatives.
    - Aggregate spam risk score pill in top header (`Spam Risk: 12% - LOW`).
  - **Zero-CDN Invariant:** Uses standard DOM range selections and CSS text-decoration-color styling.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot audit-text --file "./templates/outreach.txt" --highlight-triggers`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["riskScore", "isSafeToSend", "hasHighRiskTriggers", "triggerCount", "triggers"],
    "properties": {
      "riskScore": { "type": "number", "minimum": 0, "maximum": 100 },
      "isSafeToSend": { "type": "boolean" },
      "hasHighRiskTriggers": { "type": "boolean" },
      "triggerCount": { "type": "integer" },
      "triggers": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["word", "severity", "suggestedAlternative"],
          "properties": {
            "word": { "type": "string" },
            "severity": { "type": "string", "enum": ["low", "medium", "critical"] },
            "suggestedAlternative": { "type": "string" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Highlighting renders instantaneously on paste of up to 5,000 words without browser freezing.
  - Clicking any suggested alternative replaces text seamlessly while preserving cursor position.
  - CLI prints colored terminal highlights highlighting matches with replacement flags.

---

### Improvement 113: Instant One-Click Spintax Morphing Engine
- **ID:** `WARMUP-COPILOT-113`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Sending uniform body copy to warm-up seed mailboxes triggers hash-based content deduplication across Google Workspace and Microsoft 365.
- **Concrete UI Specification:**
  - **Layout:** Spintax generation drawer panel (`.spintax-engine-drawer`) with live permutation previewer.
  - **Visual Elements:**
    - Interactive syntax bracket badges highlighting Spintax groups `{Hello|Hi|Greetings}` in cyan pill chips.
    - "Generate Spintax Variations" AI button with rotating diamond sparkle icon.
    - Permutation counter showing total unique mathematical combinations (`Combinations: 1,296`).
    - Randomize sample button re-rolling preview pane in real-time with fade-in animation.
  - **Zero-CDN Invariant:** Local combinatorial parsing engine implemented in pure vanilla TypeScript/JS.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot spintax generate --template "Hi {Name}, hope you are doing well." --samples 5`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["sourceTemplate", "combinationCount", "isSyntacticallyValid", "hasSpintaxTokens"],
    "properties": {
      "sourceTemplate": { "type": "string" },
      "combinationCount": { "type": "integer", "minimum": 1 },
      "isSyntacticallyValid": { "type": "boolean" },
      "hasSpintaxTokens": { "type": "boolean" },
      "sampleOutputs": { "type": "array", "items": { "type": "string" } }
    }
  }
  ```
- **Acceptance Criteria:**
  - Nested Spintax patterns `{A|{B1|B2}}` parse accurately without syntax panics.
  - Permutation count handles large numbers up to $10^{12}$ combinations safely.
  - Preview card refreshes sample permutation in under 15ms.

---

### Improvement 114: Local LLM Deliverability Scoring Co-Pilot
- **ID:** `WARMUP-COPILOT-114`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Operators rely on cloud-hosted AI APIs for email feedback, risking confidential enterprise outreach copy leakage and introducing recurring API costs.
- **Concrete UI Specification:**
  - **Layout:** Co-pilot chat assistant card (`.copilot-assistant-dock`) anchored in the bottom-right corner of the editor.
  - **Visual Elements:**
    - Local engine status badge (`Local Ollama / Llama-3 8B` or `Embedded ONNX Runtime`) with green heartbeat indicator.
    - Circular deliverability readiness score dial (`0-100`) accompanied by three qualitative rating dimensions:
      - Conversational Authenticity (`88/100`)
      - Value Proposition Clarity (`94/100`)
      - ISP Filter Clearance Probability (`91/100`)
    - Streaming typewriter markdown response text stream rendering analysis insights in real time.
    - One-click "Apply Suggested Polishes" diff view with side-by-side green/red change highlights.
  - **Zero-CDN Invariant:** Local model inference via local HTTP endpoint (`http://127.0.0.1:11434`) or local WASM/ONNX runtime. Zero external API calls.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot score --template-file "./outreach.md" --local-model "llama3:latest"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["overallScore", "isReadyToSend", "hasLocalInference", "dimensions", "reasoning"],
    "properties": {
      "overallScore": { "type": "number", "minimum": 0, "maximum": 100 },
      "isReadyToSend": { "type": "boolean" },
      "hasLocalInference": { "type": "boolean" },
      "dimensions": {
        "type": "object",
        "required": ["authenticity", "clarity", "filterClearance"],
        "properties": {
          "authenticity": { "type": "number" },
          "clarity": { "type": "number" },
          "filterClearance": { "type": "number" }
        }
      },
      "reasoning": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Connects strictly to local inference server without contacting external cloud endpoints.
  - Side-by-side diff clearly displays additions in green and removals in strikethrough red.
  - Score updates automatically upon template changes after 800ms debounce.

---

### Improvement 115: Spam Trap Syntax & Honeypot Detector
- **ID:** `WARMUP-COPILOT-115`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Accidental inclusion of pristine or recycled spam trap addresses in warm-up seed lists can permanently burn sender domain reputation within hours.
- **Concrete UI Specification:**
  - **Layout:** Dedicated seed list hygiene inspection panel (`.spamtrap-detector-panel`).
  - **Visual Elements:**
    - Seed mailbox classification radar chart displaying breakdown between Organic Seeds, Corporate Seeds, and Suspicious Patterns.
    - Warning alert banner displaying specific risk flags:
      - Typo-squatted domains (`@gmai.com`, `@yaho.com`).
      - Role-based accounts (`abuse@`, `postmaster@`, `noc@`).
      - Dormant MX domain indicators (`MX points to known honeypot ASN`).
    - Red action button "Purge High-Risk Traps" with confirmation modal.
  - **Zero-CDN Invariant:** Internal regex lookup tables and local ASN blocklists stored in SQLite.
- **CLI Counterpart Command:**
  - `gitmap warmup seeds audit-traps --list-file "./seeds.csv" --purge-flagged`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["scannedCount", "cleanCount", "trapRiskCount", "isClean", "flaggedAddresses"],
    "properties": {
      "scannedCount": { "type": "integer" },
      "cleanCount": { "type": "integer" },
      "trapRiskCount": { "type": "integer" },
      "isClean": { "type": "boolean" },
      "flaggedAddresses": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["email", "riskReason"],
          "properties": {
            "email": { "type": "string" },
            "riskReason": { "type": "string" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately detects 100% of top 50 common typographical domain errors.
  - Flags generic role-based mailboxes before initiating test campaigns.
  - Purge button cleans local database list and records audit log trace.

---

### Improvement 116: Personalization Token Variance Analyzer
- **ID:** `WARMUP-COPILOT-116`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Excessive fallback values (e.g. `Hi there` replacing `{{firstName}}`) result in content duplication patterns that alert ISP reputation filters.
- **Concrete UI Specification:**
  - **Layout:** Variable distribution bar charts (`.token-variance-widget`) positioned below template inputs.
  - **Visual Elements:**
    - Horizontal stacked progress bars indicating variable fill rates vs fallback percentages.
    - Color ramp: Blue (populated unique value), Slate (generic fallback), Red (unresolved null token).
    - Unresolved token warning banner (`Warning: 18% of contacts will receive empty token!`).
    - Sample preview carousel flipping through 5 populated recipient records.
  - **Zero-CDN Invariant:** Native CSS bar styles and inline SVG icons.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot check-tokens --template-id "tpl-01" --dataset "./contacts.json"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["tokenId", "fillPercentage", "fallbackPercentage", "isSufficientVariance"],
    "properties": {
      "tokenId": { "type": "string" },
      "fillPercentage": { "type": "number", "minimum": 0, "maximum": 100 },
      "fallbackPercentage": { "type": "number", "minimum": 0, "maximum": 100 },
      "isSufficientVariance": { "type": "boolean" },
      "missingTokenCount": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Highlights empty token placeholders before batch sending can be approved.
  - Variance bar displays exact percentages calculated from linked recipient seed list.
  - CLI prints tabular summary with color-coded warning banners.

---

### Improvement 117: Reading Grade Level & Cognitive Load Calculator
- **ID:** `WARMUP-COPILOT-117`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Email templates written at an overly academic or convoluted reading grade level suffer lower organic reply rates, hindering warm-up reputation building.
- **Concrete UI Specification:**
  - **Layout:** Compact metrics strip (`.grade-level-strip`) nestled inside editor footer.
  - **Visual Elements:**
    - Flesch-Kincaid Grade Level pill badge (`Grade 6.2 - Conversational & Accessible`).
    - Sentence length distribution histogram with bar highlights for sentences exceeding 25 words.
    - Average syllables per word indicator with visual thermometer gauge.
    - One-click "Simplify Complex Sentences" prompt suggestion.
  - **Zero-CDN Invariant:** Pure JavaScript implementation of Flesch-Kincaid and Gunning fog index formulas.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot readability --file "./template.html"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["fleschKincaidGrade", "fleschReadingEase", "isConversational", "hasLongSentences"],
    "properties": {
      "fleschKincaidGrade": { "type": "number" },
      "fleschReadingEase": { "type": "number" },
      "isConversational": { "type": "boolean" },
      "hasLongSentences": { "type": "boolean" },
      "averageWordsPerSentence": { "type": "number" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately computes Flesch-Kincaid formula matching standard linguistic reference tests.
  - Sentences >25 words visually flagged with light amber background highlight in editor.
  - Readability score updates within 50ms of typing pause.

---

### Improvement 118: Link & Tracking Domain Reputation Shield
- **ID:** `WARMUP-COPILOT-118`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Including external hyperlinks or blacklisted custom tracking domains inside warm-up emails causes instant placement in junk folders regardless of sender domain age.
- **Concrete UI Specification:**
  - **Layout:** Hyperlink security inspection card (`.link-reputation-card`).
  - **Visual Elements:**
    - List of all extracted URLs in template with destination domain resolution.
    - Blacklist check status badge checking against Spamhaus DBL and SURBL blocklists stored locally.
    - Shortened URL warning (`bit.ly`, `tinyurl` flagged as high risk).
    - Link-to-text mismatch alert (e.g. anchor text says `company.com` but destination is `tracking-hop.net`).
  - **Zero-CDN Invariant:** Local DNSBL resolution performed via local socket queries.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot audit-links --template-file "./template.html"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["extractedUrlCount", "isClean", "hasBlacklistedLinks", "hasMismatchedAnchors", "links"],
    "properties": {
      "extractedUrlCount": { "type": "integer" },
      "isClean": { "type": "boolean" },
      "hasBlacklistedLinks": { "type": "boolean" },
      "hasMismatchedAnchors": { "type": "boolean" },
      "links": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["targetUrl", "isBlacklisted", "hasMismatch"],
          "properties": {
            "targetUrl": { "type": "string" },
            "isBlacklisted": { "type": "boolean" },
            "hasMismatch": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Extracts all URLs embedded in both HTML anchor tags and plain-text regex matches.
  - Correctly flags URL shorteners and anchor mismatches.
  - Warns operator if warm-up template contains more than 1 outbound hyperlink.

---

### Improvement 119: Dark-Mode Email Client Preview Simulator
- **ID:** `WARMUP-COPILOT-119`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Email designs that fail to adapt gracefully to inverted dark-mode rendering in Apple Mail and Outlook cause poor recipient engagement and lower reply scores.
- **Concrete UI Specification:**
  - **Layout:** Dual-viewport comparison frame (`.preview-mode-dual-deck`) displaying Light Mode and Inverted Dark Mode side by side.
  - **Visual Elements:**
    - Realistic email client chrome mockup (Outlook, Apple Mail, Gmail header bars).
    - Client-specific CSS color inversion engine simulating Outlook's aggressive background reversal.
    - Contrast ratio checker flagging unreadable dark text on dark backgrounds (<4.5:1 WCAG violation).
    - Image logo transparency checker highlighting opaque white logo bounding boxes.
  - **Zero-CDN Invariant:** Uses standard CSS `filter: invert(1) hue-rotate(180deg)` and custom dark rules.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot preview-dark --template "./template.html" --client "outlook"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["clientName", "hasContrastViolation", "isReadable", "minContrastRatio"],
    "properties": {
      "clientName": { "type": "string" },
      "hasContrastViolation": { "type": "boolean" },
      "isReadable": { "type": "boolean" },
      "minContrastRatio": { "type": "number", "minimum": 0, "maximum": 21 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Simulates both partial dark mode (color adjustment) and full dark mode (complete invert).
  - Flags text elements failing 4.5:1 WCAG contrast ratio with blinking red cursor.
  - CLI returns non-zero error code if contrast violations exceed configured threshold.

---

### Improvement 120: Subject Line Sentiment & Tone Radar
- **ID:** `WARMUP-COPILOT-120`
- **Cluster:** 12 (AI Deliverability Co-Pilot & Subject Line Entropy)
- **Problem Statement:** Aggressive or hyper-promotional tones depress open and reply rates, but authors lack objective tone analysis before launching warm-up sequences.
- **Concrete UI Specification:**
  - **Layout:** Five-axis SVG radar polygon chart (`.tone-radar-chart`) embedded in subject assistant panel.
  - **Visual Elements:**
    - Radar axes: Urgency, Curiosity, Professionalism, Warmth, Promotionality.
    - Translucent cyan polygon displaying calculated profile against benchmark baseline.
    - Tone label pill (`Neutral Professional`, `Aggressive Promotional`, `Conversational Warm`).
    - Recommendation tooltip suggesting adjustments to bring tones into peer alignment.
  - **Zero-CDN Invariant:** Native SVG polygon points calculated via standard trigonometry in vanilla JS.
- **CLI Counterpart Command:**
  - `gitmap warmup copilot tone --subject "Quick question regarding your infrastructure"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["primaryTone", "isAcceptableTone", "hasExcessiveUrgency", "scores"],
    "properties": {
      "primaryTone": { "type": "string" },
      "isAcceptableTone": { "type": "boolean" },
      "hasExcessiveUrgency": { "type": "boolean" },
      "scores": {
        "type": "object",
        "required": ["urgency", "curiosity", "professionalism", "warmth", "promotionality"],
        "properties": {
          "urgency": { "type": "number", "minimum": 0, "maximum": 1 },
          "curiosity": { "type": "number", "minimum": 0, "maximum": 1 },
          "professionalism": { "type": "number", "minimum": 0, "maximum": 1 },
          "warmth": { "type": "number", "minimum": 0, "maximum": 1 },
          "promotionality": { "type": "number", "minimum": 0, "maximum": 1 }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Radar polygon recomputes coordinates immediately upon subject string edit.
  - Flags promotionality scores exceeding 0.75 as high-risk delivery anomalies.
  - CLI emits clean Unicode ASCII radar graph representation.

---

# Cluster 13: Multi-Persona Seed Simulation & Conversational Dynamics (Items 121–130)

---

### Improvement 121: Autonomous Seed Persona Matrix Builder
- **ID:** `WARMUP-PERSONA-121`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** Static seed mailboxes with identical user agents and flat engagement patterns are recognized as bots by sophisticated machine-learning anti-spam systems.
- **Concrete UI Specification:**
  - **Layout:** Grid card gallery (`.persona-matrix-grid`) displaying active simulated seed profiles.
  - **Visual Elements:**
    - Persona identity cards displaying simulated avatar glyph, job role, company type, and locale.
    - Persona archetype badges (`Enterprise CTO`, `SaaS Operations Lead`, `Consumer Retailer`).
    - Configured engagement sliders for Open Rate Probability, Reply Probability, and Mark As Starred Probability.
    - Real-time online status beacon showing simulated local working hours for the persona's timezone.
  - **Zero-CDN Invariant:** Uses local CSS grid and inline SVG avatars.
- **CLI Counterpart Command:**
  - `gitmap warmup persona list --active-only`
  - `gitmap warmup persona create --role "CTO" --timezone "America/New_York" --reply-rate 0.45`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["personaId", "personaRole", "timezone", "isActive", "hasWorkHoursSchedule", "replyProbability"],
    "properties": {
      "personaId": { "type": "string" },
      "personaRole": { "type": "string" },
      "timezone": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasWorkHoursSchedule": { "type": "boolean" },
      "replyProbability": { "type": "number", "minimum": 0, "maximum": 1 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Seed actions only fire during realistic simulated business hours for the selected timezone.
  - Slider adjustments persist to SQLite within 10ms.
  - UI displays live time indicator for each persona's active timezone.

---

### Improvement 122: Multi-Turn Conversational Reply Tree Visualizer
- **ID:** `WARMUP-PERSONA-122`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** One-shot automated replies look synthetic; ISPs monitor thread depth and conversation longevity to evaluate sender credibility.
- **Concrete UI Specification:**
  - **Layout:** Interactive tree branch canvas (`.reply-tree-canvas`) rendering conversational depth across email threads.
  - **Visual Elements:**
    - Branching nodes showing Initial Send -> Seed Reply 1 -> Follow Up -> Seed Reply 2.
    - Node pills displaying message snippet, sentiment score, and response latency delay.
    - Branch depth indicator badge showing thread depth level (e.g. `Depth: 4 Turns`).
    - Clickable node inspection drawer displaying full RFC 822 email headers and Message-ID threading references (`In-Reply-To`, `References`).
  - **Zero-CDN Invariant:** Native SVG tree connectors with CSS curve styling.
- **CLI Counterpart Command:**
  - `gitmap warmup threads tree --thread-id "thr_4812" --depth 5`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["threadId", "depthLevel", "isThreadActive", "hasProperThreadingHeaders", "turns"],
    "properties": {
      "threadId": { "type": "string" },
      "depthLevel": { "type": "integer", "minimum": 1 },
      "isThreadActive": { "type": "boolean" },
      "hasProperThreadingHeaders": { "type": "boolean" },
      "turns": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["turnIndex", "senderType", "messageSnippet"],
          "properties": {
            "turnIndex": { "type": "integer" },
            "senderType": { "type": "string", "enum": ["system", "seed_persona"] },
            "messageSnippet": { "type": "string" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Tree properly renders multi-turn branches up to 8 conversation turns deep.
  - Verifies presence of valid `In-Reply-To` and `References` headers in simulated replies.
  - CLI displays structured ASCII indentation tree representation.

---

### Improvement 123: Authentic Dwell-Time & Scroll Velocity Simulator
- **ID:** `WARMUP-PERSONA-123`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** Instantaneous seed opens (<1 second) are flagged by webmail analytics as automated script interactions, negating reputation gains.
- **Concrete UI Specification:**
  - **Layout:** Behavioral physics dashboard card (`.dwell-physics-card`).
  - **Visual Elements:**
    - Dwell-time distribution curve showing Gaussian bell curve distribution around human reading speeds (12 to 45 seconds).
    - Real-time reading simulation gauge displaying simulated scroll velocity and viewport pauses.
    - Monospace readout indicating simulated WPM (Words Per Minute) calculation (`245 WPM`).
    - Random jitter toggle introducing non-deterministic delay variations (+/- 25%).
  - **Zero-CDN Invariant:** Pure SVG normal distribution curve and CSS progress bars.
- **CLI Counterpart Command:**
  - `gitmap warmup behavior set-dwell --mean-seconds 24 --jitter-percent 20`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["meanDwellSeconds", "jitterPercentage", "isActive", "hasGaussianNoise"],
    "properties": {
      "meanDwellSeconds": { "type": "number", "minimum": 5, "maximum": 120 },
      "jitterPercentage": { "type": "number", "minimum": 0, "maximum": 50 },
      "isActive": { "type": "boolean" },
      "hasGaussianNoise": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Seed execution loop applies Box-Muller Gaussian transform to calculate human-like dwell times.
  - Prevents sub-second automated opens across 100% of simulated seed mailboxes.
  - Interactive slider updates distribution curve graphic instantaneously.

---

### Improvement 124: "Mark as Important" & Star Interaction Physics
- **ID:** `WARMUP-PERSONA-124`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** Passive inbox delivery without positive user engagement signals (starring, flagging as important, pinning) leaves reputation score gains sluggish.
- **Concrete UI Specification:**
  - **Layout:** Positive signal orchestration panel (`.positive-signals-panel`).
  - **Visual Elements:**
    - Toggle matrix for affirmative mailbox actions:
      - Star Message (`⭐`)
      - Mark as Important / High Priority (`🏷️`)
      - Move to Category / Primary Tab (`📥`)
      - Add Sender to Contacts / Address Book (`👤`)
    - Action probability dials allowing individual threshold tuning (e.g. `Star Rate: 18%`).
    - Cumulative reputation boost multiplier readout (`Reputation Multiplier: 1.45x`).
  - **Zero-CDN Invariant:** Native SVG icons and CSS switch controls.
- **CLI Counterpart Command:**
  - `gitmap warmup signals configure --star-rate 0.15 --add-contact-rate 0.05`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["starRate", "markImportantRate", "addToContactsRate", "isActive"],
    "properties": {
      "starRate": { "type": "number", "minimum": 0, "maximum": 1 },
      "markImportantRate": { "type": "number", "minimum": 0, "maximum": 1 },
      "addToContactsRate": { "type": "number", "minimum": 0, "maximum": 1 },
      "isActive": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Seed IMAP worker dispatches corresponding IMAP `\Flagged` and `$PhishingClean` flags.
  - Confirmed positive signals update SQLite tracking journal immediately.
  - Multiplier recalculates dynamically based on assigned action frequencies.

---

### Improvement 125: Spam Rescue & "Not Spam" Evacuation Engine
- **ID:** `WARMUP-PERSONA-125`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** When test emails land in the spam folder, failing to autonomously move them to the primary inbox permanently damages reputation with that provider.
- **Concrete UI Specification:**
  - **Layout:** Spam evacuation radar card (`.spam-rescue-radar-card`) with prominent rescue counter.
  - **Visual Elements:**
    - Large digital counter showing rescued emails in the last 24 hours (`27 Rescued`).
    - Animated progress track showing 3-stage rescue workflow:
      1. Spam Folder Detected (`Junk/Spam`)
      2. Flagged as "Not Spam" (`IMAP MOVE INBOX`)
      3. Reply & Mark as Starred Sent
    - Rescued ISP breakdown pie chart (Google vs Microsoft vs Yahoo).
    - Real-time alert pill flashing green when an automated rescue completes.
  - **Zero-CDN Invariant:** Native CSS SVG circular progress meters.
- **CLI Counterpart Command:**
  - `gitmap warmup rescue run-sweep --domain "acme-corp.com"`
  - `gitmap warmup rescue status`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["rescuedLast24h", "pendingSpamCount", "isRescueActive", "hasAutoInboxMove"],
    "properties": {
      "rescuedLast24h": { "type": "integer", "minimum": 0 },
      "pendingSpamCount": { "type": "integer", "minimum": 0 },
      "isRescueActive": { "type": "boolean" },
      "hasAutoInboxMove": { "type": "boolean" },
      "lastRescueTimestamp": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - IMAP worker autonomously scans spam/junk folders every 15 minutes.
  - Dispatches `MOVE` command to `INBOX` and clears spam classification tags.
  - Rescue count increments accurately in both UI and SQLite state.

---

### Improvement 126: Seed Mailbox IMAP/SMTP Connection Vault
- **ID:** `WARMUP-PERSONA-126`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** Storing seed mailbox credentials in plaintext configuration files creates security vulnerabilities and violates SOC2 compliance requirements.
- **Concrete UI Specification:**
  - **Layout:** Encrypted credentials table (`.credential-vault-table`) with masked secret fields.
  - **Visual Elements:**
    - Password and token inputs rendered as hidden bullet points (`••••••••••••`).
    - "Test Connectivity" button with animated spinning lock glyph.
    - Security encryption badge indicating AES-256-GCM hardware key storage.
    - OAuth2 refresh token expiration indicator showing days until re-authorization.
  - **Zero-CDN Invariant:** Uses standard HTML password inputs and CSS security badges.
- **CLI Counterpart Command:**
  - `gitmap warmup vault add-seed --email "seed-01@domain.com" --provider "google-oauth" --encrypt`
  - `gitmap warmup vault test-connection --seed-id "seed-01"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["seedId", "emailAddress", "authType", "isConnected", "isEncrypted"],
    "properties": {
      "seedId": { "type": "string" },
      "emailAddress": { "type": "string" },
      "authType": { "type": "string", "enum": ["password", "app_password", "oauth2"] },
      "isConnected": { "type": "boolean" },
      "isEncrypted": { "type": "boolean" },
      "tokenExpiryIso": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - All credentials stored encrypted with AES-256-GCM in local SQLite database.
  - Plaintext passwords never logged to terminal or audit trails.
  - Connection test validates both IMAP SSL and SMTP TLS handshake successfully.

---

### Improvement 127: Synthetic Out-of-Office (OOO) Auto-Responder
- **ID:** `WARMUP-PERSONA-127`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** Real corporate mailboxes generate authentic auto-reply bouncebacks and vacation responders; warm-up pools without them appear unnatural to filtering algorithms.
- **Concrete UI Specification:**
  - **Layout:** Out-of-office simulation schedule card (`.ooo-responder-card`).
  - **Visual Elements:**
    - Calendar date picker defining simulated vacation or leave windows for personas.
    - OOO template editor pre-filled with professional corporate templates (`"I am currently out of the office until Monday..."`).
    - Probability slider controlling how frequently an OOO response is fired (e.g. `2% of initial outreaches`).
    - Header compliance tag verifying presence of `Auto-Submitted: auto-replied` header.
  - **Zero-CDN Invariant:** Native HTML5 date pickers and CSS grid layouts.
- **CLI Counterpart Command:**
  - `gitmap warmup persona set-ooo --persona-id "persona_03" --from "2026-11-01" --to "2026-11-08"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["personaId", "isOooActive", "hasAutoSubmittedHeader", "startDate", "endDate"],
    "properties": {
      "personaId": { "type": "string" },
      "isOooActive": { "type": "boolean" },
      "hasAutoSubmittedHeader": { "type": "boolean" },
      "startDate": { "type": "string" },
      "endDate": { "type": "string" },
      "responseTemplate": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Generates RFC 3834 compliant auto-reply headers including `Auto-Submitted: auto-replied`.
  - Simulates temporary vacation periods for specified personas without breaking warm-up thread continuity.
  - Tracks auto-reply events in local SQLite audit journal.

---

### Improvement 128: Seed Mailbox Storage Quota & Purge Gauge
- **ID:** `WARMUP-PERSONA-128`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** Seed mailboxes accumulate thousands of test emails over weeks, eventually filling storage quotas and causing connection failures that abort warm-up runs.
- **Concrete UI Specification:**
  - **Layout:** Storage usage gauge card (`.mailbox-quota-card`) with linear percentage bars.
  - **Visual Elements:**
    - Multi-segment storage bar showing Used vs Free space with color thresholds (Green <70%, Amber 70-90%, Rose >90%).
    - Automated lifecycle purge configuration toggle (`Delete emails older than 30 days`).
    - "Execute Immediate Purge" button with confirmation counter showing reclaimable megabytes.
    - Archive option to export mail store to compressed Mbox file before purging.
  - **Zero-CDN Invariant:** Native CSS flexbox and linear gradient progress bars.
- **CLI Counterpart Command:**
  - `gitmap warmup seeds quota-check --all-seeds`
  - `gitmap warmup seeds purge-old --older-than-days 30 --dry-run=false`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["seedId", "usedMegabytes", "quotaMegabytes", "isNearFull", "hasAutoPurge"],
    "properties": {
      "seedId": { "type": "string" },
      "usedMegabytes": { "type": "number" },
      "quotaMegabytes": { "type": "number" },
      "isNearFull": { "type": "boolean" },
      "hasAutoPurge": { "type": "boolean" },
      "oldestMessageDate": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Queries IMAP `GETQUOTAROOT` where supported, falling back to message size summation.
  - Automatically purges messages older than retention window when auto-purge is enabled.
  - Warns operator prominently when any mailbox reaches 85% capacity.

---

### Improvement 129: Multi-Language Seed Conversational Corpus
- **ID:** `WARMUP-PERSONA-129`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** Running warm-up campaigns for European or Asian markets using purely English seed conversational text triggers cross-border language classification anomalies.
- **Concrete UI Specification:**
  - **Layout:** Language corpus selector dropdown and distribution matrix (`.language-corpus-matrix`).
  - **Visual Elements:**
    - Locale selection chips for major outbound target markets (German `de-DE`, French `fr-FR`, Japanese `ja-JP`, Spanish `es-ES`, English `en-US`).
    - Corpus distribution slider bar allowing percentage allocation across languages.
    - Live snippet preview displaying simulated localized dialogue exchanges.
    - Character encoding verification badge ensuring correct UTF-8 and ISO-2022-JP charset headers.
  - **Zero-CDN Invariant:** Uses native Unicode font stacks and CSS pill components.
- **CLI Counterpart Command:**
  - `gitmap warmup corpus set-languages --distribution "en-US=50,de-DE=30,fr-FR=20"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["primaryLocale", "isMultiLanguageActive", "hasUtf8Encoding", "distributions"],
    "properties": {
      "primaryLocale": { "type": "string" },
      "isMultiLanguageActive": { "type": "boolean" },
      "hasUtf8Encoding": { "type": "boolean" },
      "distributions": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["localeCode", "percentage"],
          "properties": {
            "localeCode": { "type": "string" },
            "percentage": { "type": "number", "minimum": 0, "maximum": 100 }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Seed dialog engine generates grammatically natural replies in configured target locales.
  - Generates appropriate RFC 2047 encoded-word headers for non-ASCII subject lines.
  - Distribution percentages enforce exact 100% total sum validation.

---

### Improvement 130: Autonomous Seed Pool Health & Burn-Rate Circuit Breaker
- **ID:** `WARMUP-PERSONA-130`
- **Cluster:** 13 (Multi-Persona Seed Simulation & Conversational Dynamics)
- **Problem Statement:** If a seed pool is detected and blacklisted by an ISP, continuing to route traffic to it burns sending domain reputation across the entire fleet.
- **Concrete UI Specification:**
  - **Layout:** High-visibility circuit-breaker master card (`.circuit-breaker-master-card`) anchored at the top of the persona control center.
  - **Visual Elements:**
    - Master circuit status badge (`NORMAL - CLOSED` in emerald green, `TRIPPED - OPEN` in crimson red).
    - Threshold setting controls for automated tripping:
      - Bounce Rate > 3.5%
      - Spam Placement > 15%
      - Authentication Failure > 1%
    - Emergency manual trip switch (`button.emergency-kill-switch`) styled with industrial striped border warning.
    - Post-trip cooldown timer with step-by-step diagnostic checklist required for re-engagement.
  - **Zero-CDN Invariant:** Custom CSS hatched warning borders and native HTML button controls.
- **CLI Counterpart Command:**
  - `gitmap warmup circuit-breaker status`
  - `gitmap warmup circuit-breaker trip --reason "Manual operator intervention"`
  - `gitmap warmup circuit-breaker reset --force`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["circuitState", "isTripped", "hasEmergencyHalt", "failureMetrics", "lastTripTimestamp"],
    "properties": {
      "circuitState": { "type": "string", "enum": ["closed", "half_open", "open"] },
      "isTripped": { "type": "boolean" },
      "hasEmergencyHalt": { "type": "boolean" },
      "failureMetrics": {
        "type": "object",
        "required": ["bounceRate", "spamRate"],
        "properties": {
          "bounceRate": { "type": "number" },
          "spamRate": { "type": "number" }
        }
      },
      "lastTripTimestamp": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Automatically trips to `open` state and halts all warm-up dispatches if bounce rate exceeds 3.5%.
  - Emits high-priority alert notification in UI and audit log.
  - Reset requires explicit operator acknowledgement or CLI `--force` flag.
