# Specification 339: Improvements 201 to 230 (Web Components, Time-Travel & Curves)

- **Sequence:** 339
- **Slug:** `339-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/02-improvements-201-to-230-web-components-time-travel-curves.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## Improvement Catalog Summary (201 - 230)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 21** | Items 201–210 | Zero-CDN Shadow DOM Custom Web Components |
| **Cluster 22** | Items 211–220 | Time-Travel Deliverability Debugger & Event Playback |
| **Cluster 23** | Items 221–230 | Adaptive Warming Gradient & Reputation Acceleration Engine |

---

# Cluster 21: Zero-CDN Shadow DOM Custom Web Components (Items 201–210)

---

### Improvement 201: Encapsulated Circular Deliverability Metric Gauge Element
- **ID:** `WARMUP-WC-201`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Existing status dashboards embed monolithic HTML snippets that leak styles into parent pages and break when imported into different administrative portals or third-party webviews.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-metric-gauge>`
  - **Shadow DOM Structure:**
    - Root host with CSS variables `--gauge-fill-color`, `--gauge-bg-track`, and `--gauge-radius`.
    - SVG-based circular dial (`<svg viewBox="0 0 100 100" class="gauge-ring">`) with animated `stroke-dasharray` and `stroke-dashoffset` representing completion percentage (0–100%).
    - Central numeric readout (`<span class="gauge-value">98.4%</span>`) and subtitle label (`<span class="gauge-label">Inbox Placement</span>`).
    - Ambient radial glow pulsing on high thresholds (`filter: drop-shadow(0 0 6px var(--gauge-fill-color))`).
  - **Zero-CDN Invariant:** Native ES2022 `customElements.define('warmup-metric-gauge', MetricGaugeElement)` with inline `<template>` and encapsulated `<style>` block. Zero external assets or libraries.
- **CLI Counterpart Command:**
  - `gitmap warmup component render --tag "warmup-metric-gauge" --metric-val "98.4" --label "Inbox Placement"`
  - `gitmap warmup component inspect --name "metric-gauge"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["elementId", "metricValue", "metricLabel", "isActive", "hasOptimalScore", "strokeColorHex"],
    "properties": {
      "elementId": { "type": "string" },
      "metricValue": { "type": "number", "minimum": 0, "maximum": 100 },
      "metricLabel": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasOptimalScore": { "type": "boolean" },
      "strokeColorHex": { "type": "string", "pattern": "^#[0-9a-fA-F]{6}$" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Element mounts inside Shadow DOM without bleeding CSS styles into external DOM tree.
  - Setting attribute `data-value="75"` triggers smooth CSS transition across the SVG arc in <= 250ms.
  - `gitmap warmup component render` outputs pristine standalone HTML snippet.

---

### Improvement 202: Glassmorphic Domain Warmup Health Envelope Card Element
- **ID:** `WARMUP-WC-202`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Domain identity overviews require cluttered table rows that fail to convey composite status (SPF, DKIM, DMARC, IP warmth) in a compact, modular card format.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-envelope-card>`
  - **Shadow DOM Structure:**
    - Container card (`.envelope-card`) styled with `backdrop-filter: blur(12px)`, background `rgba(17, 24, 39, 0.85)`, and 1px border `rgba(255, 255, 255, 0.08)`.
    - Header slot (`<slot name="header">`) displaying sender domain name and active phase badge (`.badge-phase-warming`).
    - Micro-grid status row rendering 4 mini status pills: SPF, DKIM, DMARC, and MX, each showing an inline SVG checkmark or warning icon.
    - Mini sparkline curve (`<svg class="mini-sparkline">`) showing recent 7-day outbound volume trajectory.
  - **Zero-CDN Invariant:** Completely self-contained stylesheet within Shadow Root using native CSS grid and flexbox.
- **CLI Counterpart Command:**
  - `gitmap warmup component card --domain "outbound.gitmap.dev" --format web-component`
  - `gitmap warmup domain card-status --domain "outbound.gitmap.dev"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["cardId", "domainName", "currentPhase", "isActive", "hasDkimValid", "hasSpfValid", "hasDmarcValid"],
    "properties": {
      "cardId": { "type": "string" },
      "domainName": { "type": "string" },
      "currentPhase": { "type": "string", "enum": ["initial_ramp", "reputation_building", "peak_maintenance", "cooldown"] },
      "isActive": { "type": "boolean" },
      "hasDkimValid": { "type": "boolean" },
      "hasSpfValid": { "type": "boolean" },
      "hasDmarcValid": { "type": "boolean" },
      "sparklinePoints": {
        "type": "array",
        "items": { "type": "number" }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Renders all 4 authentication status indicators with responsive SVGs in under 2ms initialization time.
  - Exposes custom event `domain-selected` with domain payload when user clicks the card body.
  - Handles light and dark mode automatically via host CSS token binding.

---

### Improvement 203: Self-Contained Deliverability Time-Travel Scrubber Bar Element
- **ID:** `WARMUP-WC-203`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Historical timeline navigation tools depend on heavy external range slider libraries (e.g. noUiSlider) that increase bundle weight and introduce runtime security vulnerabilities.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-time-scrubber>`
  - **Shadow DOM Structure:**
    - Scrub track (`.scrubber-track`) with gradient color stops indicating deliverability health across historical time slices.
    - Draggable thumb element (`.scrubber-thumb`) with active grab indicator, timestamp tooltip on hover, and keyframe snap ticks.
    - Play/Pause toggle button with embedded SVG icons (`<svg class="icon-play">`, `<svg class="icon-pause">`).
    - Scrub speed selector pill (`1x`, `2x`, `5x`, `10x`) styled with hairline borders and active cyan indicator.
  - **Zero-CDN Invariant:** Implemented strictly using native pointer events (`setPointerCapture`, `onpointerdown`, `onpointermove`) and native requestAnimationFrame. Zero JS dependencies.
- **CLI Counterpart Command:**
  - `gitmap warmup timeline scrubber --start-epoch 1704067200 --end-epoch 1704672000`
  - `gitmap warmup timeline export-ticks --domain "mail.gitmap.dev"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["scrubberId", "startEpoch", "endEpoch", "currentEpoch", "isPlaying", "hasKeyframeSnapping"],
    "properties": {
      "scrubberId": { "type": "string" },
      "startEpoch": { "type": "integer" },
      "endEpoch": { "type": "integer" },
      "currentEpoch": { "type": "integer" },
      "playbackSpeedMultiplier": { "type": "number", "minimum": 0.25, "maximum": 50 },
      "isPlaying": { "type": "boolean" },
      "hasKeyframeSnapping": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Pointer dragging updates `currentEpoch` at 60 FPS without layout jitter.
  - Emits native CustomEvent `time-scrub` with payload `{ epoch: 1704123400, isLive: false }`.
  - Tabbing through controls allows keyboard scrubbing using ArrowLeft and ArrowRight keys.

---

### Improvement 204: Shadow DOM ISP Inbox vs Spam Foldering Split Meter Element
- **ID:** `WARMUP-WC-204`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Visualizing segmented placement ratios (Inbox, Spam, Promotions, Missing) across major providers (Google, Microsoft, Yahoo) currently requires rendering full-page canvas dashboards.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-isp-split-meter>`
  - **Shadow DOM Structure:**
    - Horizontal segmented progress bar (`.split-meter-bar`) divided into four color-coded segments: Inbox (`#10b981`), Spam (`#ef4444`), Promo (`#f59e0b`), and Missing (`#6b7280`).
    - Floating legend header displaying ISP brand icon (SVG) and percentage distribution tags.
    - Hover expansion drawer expanding into detailed per-seed counts (`seed_count_inbox`, `seed_count_spam`).
  - **Zero-CDN Invariant:** Native CSS flexbox with `flex-grow` proportional to segment ratios and transition timing `width 300ms ease`.
- **CLI Counterpart Command:**
  - `gitmap warmup isp meter --provider "google" --domain "acme.org"`
  - `gitmap warmup isp stats --format web-component`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["meterId", "ispName", "inboxRatio", "spamRatio", "promoRatio", "isActive", "hasSpamAnomaly"],
    "properties": {
      "meterId": { "type": "string" },
      "ispName": { "type": "string", "enum": ["google", "microsoft", "yahoo", "apple", "generic"] },
      "inboxRatio": { "type": "number", "minimum": 0, "maximum": 1 },
      "spamRatio": { "type": "number", "minimum": 0, "maximum": 1 },
      "promoRatio": { "type": "number", "minimum": 0, "maximum": 1 },
      "isActive": { "type": "boolean" },
      "hasSpamAnomaly": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Segment ratios always sum to 1.0 (100%) with automatic normalization guard.
  - Hovering over spam segment displays tooltip with exact count of spam trap hits.
  - Shadow DOM encapsulation prevents host page flexbox rules from collapsing segment widths.

---

### Improvement 205: Adaptive Sigmoid Warming Curve Visualizer Web Component
- **ID:** `WARMUP-WC-205`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Operators cannot embed live interactive mathematical warming curves into ad-hoc reports or markdown previews without importing multi-megabyte charting libraries.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-curve-visualizer>`
  - **Shadow DOM Structure:**
    - High-density SVG coordinate viewport (`<svg viewBox="0 0 600 240" class="curve-viewport">`).
    - Background grid mesh with subtle dashed lines (`stroke: rgba(255, 255, 255, 0.05)`).
    - Calculated SVG `<path>` plotting mathematical sigmoid curve: $V(t) = \frac{V_{max}}{1 + e^{-k(t - t_0)}}$.
    - Secondary overlay path showing actual sent volume compared against planned curve.
    - Interactive hover crosshair vertical line displaying day number, planned volume, and actual volume.
  - **Zero-CDN Invariant:** Mathematical path coordinates calculated directly in component constructor via pure JavaScript string interpolation. Zero Chart.js or Highcharts.
- **CLI Counterpart Command:**
  - `gitmap warmup curve render --domain "outbound.ai" --days 30 --format svg-component`
  - `gitmap warmup curve calculate --base 50 --max 5000 --days 21`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["curveId", "domainName", "modelType", "plannedVolumePoints", "isActive", "hasDeviationAlert"],
    "properties": {
      "curveId": { "type": "string" },
      "domainName": { "type": "string" },
      "modelType": { "type": "string", "enum": ["logistic_sigmoid", "exponential", "stepped_linear"] },
      "plannedVolumePoints": {
        "type": "array",
        "items": { "type": "number" }
      },
      "actualVolumePoints": {
        "type": "array",
        "items": { "type": "number" }
      },
      "isActive": { "type": "boolean" },
      "hasDeviationAlert": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Computes and renders 60-point Bezier curve in < 4ms.
  - Crosshair aligns precisely with cursor pointer position and updates values in live DOM readout.
  - Automatically scales vertical axis based on maximum projected volume.

---

### Improvement 206: Real-Time SMTP Gateway Throughput Dial Element
- **ID:** `WARMUP-WC-206`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Gateway concurrency limits and MTA channel throttling indicators are buried in log text files rather than surfaced in compact HUD gauges for NOC engineers.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-gateway-dial>`
  - **Shadow DOM Structure:**
    - Semi-circular speedometer gauge with needle indicator rotated via CSS transform (`transform: rotate(calc(var(--speed-angle) * 1deg))`).
    - Needle dampening spring animation (`transition: transform 400ms cubic-bezier(0.34, 1.56, 0.64, 1)`).
    - Warning threshold arc (`#f59e0b`) and throttling danger arc (`#ef4444`) on outer perimeter.
    - Digital readout displaying current emails/minute throughput and active connection count.
  - **Zero-CDN Invariant:** CSS variables calculate needle rotation angle dynamically from element attribute `data-throughput`.
- **CLI Counterpart Command:**
  - `gitmap warmup gateway dial --gateway "smtp-eu-01.gitmap.internal"`
  - `gitmap warmup gateway live --interval 1s`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["gatewayId", "currentThroughputPerMin", "maxCapacityPerMin", "isActive", "hasThrottled", "activeConnections"],
    "properties": {
      "gatewayId": { "type": "string" },
      "currentThroughputPerMin": { "type": "integer", "minimum": 0 },
      "maxCapacityPerMin": { "type": "integer", "minimum": 1 },
      "activeConnections": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasThrottled": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Needle rotation dynamically bounds between -90deg (0 rate) and +90deg (100% saturation).
  - Emits `gateway-throttle-alert` event when throughput exceeds 95% capacity for > 5 seconds.
  - Self-contained SVG dial assets embedded directly inside element Shadow Root.

---

### Improvement 207: Multi-Tab Seed Persona Simulation Activity Feed Element
- **ID:** `WARMUP-WC-207`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Tracking virtual seed interactions across synthetic inboxes requires separate windows, making it difficult to verify persona activity in real time.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-persona-feed>`
  - **Shadow DOM Structure:**
    - Header filter bar allowing filtering by persona type: Executive, IT Admin, Sales, Consumer.
    - Virtualized scrolling stream (`.persona-event-list`) displaying simulated actions: Read Email (42s dwell time), Star/Mark Important, Move from Spam to Inbox, Reply Thread.
    - Avatar micro-icon with persona initials and colored status ring indicating online status.
    - Action duration pill showing synthetic engagement dwell time.
  - **Zero-CDN Invariant:** Pure CSS virtual scroll container using `contain: content` and native DOM element recycling.
- **CLI Counterpart Command:**
  - `gitmap warmup persona feed --persona "vp-engineering" --limit 20`
  - `gitmap warmup persona stream --format wc-embed`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["feedId", "personaId", "actionType", "dwellTimeSeconds", "isActive", "hasReplied", "timestampEpoch"],
    "properties": {
      "feedId": { "type": "string" },
      "personaId": { "type": "string" },
      "actionType": { "type": "string", "enum": ["open", "read", "star", "unspam", "reply", "archive"] },
      "dwellTimeSeconds": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasReplied": { "type": "boolean" },
      "timestampEpoch": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Renders 500 stream items without causing host window frame drops.
  - Automatically prepends new simulated events with slide-down CSS animation (`opacity: 1; transform: translateY(0)`).
  - Clicking a persona entry opens detailed simulation log drawer.

---

### Improvement 208: Cryptographic Signature & SPF/DKIM Key Status Badge Element
- **ID:** `WARMUP-WC-208`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Technical DNS security keys (DKIM 2048-bit selectors, DMARC policy enforcement) require verbose text displays that operators struggle to parse quickly.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-crypto-badge>`
  - **Shadow DOM Structure:**
    - Inline micro-pill badge (`.crypto-badge`) displaying selector tag, key bit length (e.g. `rsa2048`), and cryptographic hash algorithm.
    - Status beacon light (`.beacon-indicator`) glowing emerald for valid signatures, rose for verification failures.
    - Hover popover displaying complete public key DNS record and alignment status (Strict vs Relaxed).
  - **Zero-CDN Invariant:** Inline SVG lock vector with CSS clip-path highlighting key verification status.
- **CLI Counterpart Command:**
  - `gitmap warmup crypto verify --domain "send.domain.com" --selector "gm2026"`
  - `gitmap warmup crypto badge --domain "send.domain.com"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["badgeId", "selector", "keyBitLength", "isActive", "hasValidSignature", "alignmentMode"],
    "properties": {
      "badgeId": { "type": "string" },
      "selector": { "type": "string" },
      "keyBitLength": { "type": "integer", "enum": [1024, 2048, 4096] },
      "alignmentMode": { "type": "string", "enum": ["strict", "relaxed"] },
      "isActive": { "type": "boolean" },
      "hasValidSignature": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Displays valid green state when `hasValidSignature` is true, switching to red pulse when false.
  - Popover respects document viewport boundaries without clipping.
  - Native tooltip accessible via ARIA `aria-describedby` attribute.

---

### Improvement 209: Responsive Outbound Email Volume Heatmap Calendar Element
- **ID:** `WARMUP-WC-209`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** Weekly and monthly warming volume distribution cannot be verified at a glance, making weekend volume drops and ramp irregularities difficult to detect.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-volume-heatmap>`
  - **Shadow DOM Structure:**
    - 7x5 CSS grid layout representing days of the week and weeks of the warming phase.
    - Day cell elements (`.heatmap-cell`) styled with 5 discrete background color steps from dark navy (`rgba(17, 24, 39, 0.4)`) to vivid cyan (`#06b6d4`).
    - Tooltip displaying calendar date, total emails sent, target volume, and adherence percentage.
    - Day headers (Mon, Tue, Wed, Thu, Fri, Sat, Sun) styled with muted uppercase typography.
  - **Zero-CDN Invariant:** Pure CSS Grid with CSS custom property `--cell-heat-level: 0..4`.
- **CLI Counterpart Command:**
  - `gitmap warmup heatmap render --month "2026-10" --domain "sales.org"`
  - `gitmap warmup heatmap export --format json`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["heatmapId", "monthIdentifier", "daysData", "isActive", "hasAnomalyDay"],
    "properties": {
      "heatmapId": { "type": "string" },
      "monthIdentifier": { "type": "string", "pattern": "^[0-9]{4}-[0-9]{2}$" },
      "isActive": { "type": "boolean" },
      "hasAnomalyDay": { "type": "boolean" },
      "daysData": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["dateString", "volumeSent", "heatLevel", "hasMetTarget"],
          "properties": {
            "dateString": { "type": "string" },
            "volumeSent": { "type": "integer", "minimum": 0 },
            "heatLevel": { "type": "integer", "minimum": 0, "maximum": 4 },
            "hasMetTarget": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Automatically calculates heat levels based on relative percentage of target daily volume.
  - Cells with `hasMetTarget: false` display subtle amber corner warning pip.
  - Heatmap is fully responsive, collapsing cleanly into horizontal strip on mobile viewports.

---

### Improvement 210: Autonomous Warmup Alert Banner & Remediation Action Bar Element
- **ID:** `WARMUP-WC-210`
- **Cluster:** 21 (Zero-CDN Shadow DOM Custom Web Components)
- **Problem Statement:** When warming domains encounter sudden spam spikes or DNS invalidations, notification banners are detached from remediation actions, requiring operators to navigate across multiple screens.
- **Concrete UI Specification:**
  - **Tag Name:** `<warmup-alert-action-bar>`
  - **Shadow DOM Structure:**
    - Slide-in glassmorphic bar (`.alert-bar`) fixed to bottom or top of viewport.
    - Status icon slot with animated SVG alert beacon.
    - Problem description text with highlighted error tokens.
    - Quick-action buttons: "Decelerate Warming", "Pause IP Pool", "Run DNS Repair", and "Dismiss".
    - Action progress loader indicator embedded in button states.
  - **Zero-CDN Invariant:** Self-contained CSS keyframe animations and native click handlers emitting strongly-typed action events.
- **CLI Counterpart Command:**
  - `gitmap warmup alert bar --severity "warning" --action "decelerate"`
  - `gitmap warmup alert trigger --domain "outbound.io" --type "spam_trap"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["alertId", "severityLevel", "alertMessage", "recommendedAction", "isActive", "hasUserInteracted"],
    "properties": {
      "alertId": { "type": "string" },
      "severityLevel": { "type": "string", "enum": ["info", "warning", "critical", "emergency"] },
      "alertMessage": { "type": "string" },
      "recommendedAction": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasUserInteracted": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Mounts with slide-up transition in 200ms when triggered.
  - Clicking "Decelerate Warming" emits `warmup-remediate` CustomEvent and sets button state to `isProcessing: true`.
  - Conforms to WCAG AA contrast standards across all severity color themes.

---

# Cluster 22: Time-Travel Deliverability Debugger & Event Playback (Items 211–220)

---

### Improvement 211: Deterministic Deliverability State Keyframe Snapshot Journal
- **ID:** `WARMUP-TIME-211`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Deliverability logs record unstructured events, making it impossible to reconstruct the exact system state (active volume, IP reputation, queue depth, DNS alignment) at an arbitrary historical timestamp.
- **Concrete UI Specification:**
  - **Layout:** State Keyframe Inspector table (`.keyframe-table`) displaying temporal snapshot entries with expandable JSON diff viewer.
  - **Visual Elements:**
    - Keyframe marker badges with epoch timestamp, formatted UTC string, and state delta byte size.
    - State hash chip (`#a9f82d...`) allowing cryptographic verification of historical snapshot integrity.
    - Split diff view highlighting variable shifts between consecutive keyframes (e.g. `reputationScore: 0.94 -> 0.82`).
    - "Restore to Viewport" button enabling instantaneous projection of the historical state onto all active UI components.
  - **Zero-CDN Invariant:** Pure CSS diff styling with green (`+`) and red (`-`) background highlights.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog snapshot capture --domain "outbound.corp"`
  - `gitmap warmup timelog snapshot list --domain "outbound.corp" --since "24h"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["snapshotId", "domainName", "epochTimestamp", "stateHash", "isActive", "hasDeltaChanges"],
    "properties": {
      "snapshotId": { "type": "string" },
      "domainName": { "type": "string" },
      "epochTimestamp": { "type": "integer" },
      "stateHash": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasDeltaChanges": { "type": "boolean" },
      "stateMetrics": {
        "type": "object",
        "required": ["reputationScore", "activeVolume", "inboxPlacementRatio"],
        "properties": {
          "reputationScore": { "type": "number" },
          "activeVolume": { "type": "integer" },
          "inboxPlacementRatio": { "type": "number" }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Captures state keyframes every 15 minutes or immediately following anomaly triggers.
  - Snapshot delta compression reduces storage overhead by > 80% compared to raw full dumps.
  - Restoring a keyframe updates all child Web Components in < 10ms.

---

### Improvement 212: High-Speed Temporal Event Playback Transport Controller
- **ID:** `WARMUP-TIME-212`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Debugging cascading deliverability failures requires manually scrolling through thousands of lines of logs, losing the chronological sequence of events as they unfolded.
- **Concrete UI Specification:**
  - **Layout:** Sleek media transport toolbar (`.playback-transport-bar`) docked at the bottom of the time-travel workbench.
  - **Visual Elements:**
    - Play, Pause, Step-Forward (1 frame), Step-Backward (1 frame), and Jump-to-Incident buttons with embedded SVGs.
    - Playback speed rocker switch (`0.5x`, `1x`, `2x`, `5x`, `20x`, `60x real-time`).
    - Circular progress ring indicating playback progression through current inspection window.
    - Timecode display showing both absolute UTC time (`2026-10-27 14:32:00Z`) and relative offset (`T-02:14:30`).
  - **Zero-CDN Invariant:** Native JS timer loop using high-precision `performance.now()` clock.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog play --domain "mail.dev" --start "-6h" --speed 10x`
  - `gitmap warmup timelog step --direction forward --count 1`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["transportId", "currentOffsetSeconds", "playbackSpeed", "isPlaying", "hasReachedEnd"],
    "properties": {
      "transportId": { "type": "string" },
      "currentOffsetSeconds": { "type": "number" },
      "playbackSpeed": { "type": "number" },
      "isPlaying": { "type": "boolean" },
      "hasReachedEnd": { "type": "boolean" },
      "activeKeyframeId": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Playback engine dispatches time-tick events at precise intervals corresponding to selected multiplier.
  - Stepping 1 frame forward/backward accurately positions state at exactly the adjacent keyframe.
  - UI transport controls remain responsive without UI freezing during 60x high-speed playback.

---

### Improvement 213: ISP Reputation Degradation Point Incident Auto-Seeker
- **ID:** `WARMUP-TIME-213`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Locating the exact moment an ISP begins silently junking outbound messages requires tedious manual timeline inspection across days of data.
- **Concrete UI Specification:**
  - **Layout:** Incident bookmark strip (`.incident-bookmark-strip`) embedded above the scrubber track.
  - **Visual Elements:**
    - Diamond-shaped SVG marker pips pinned to exact timeline coordinates where reputation slope $\frac{dRep}{dt} < -\theta$.
    - Incident preview popover displaying ISP name (e.g. Microsoft 365), drop severity (-24%), and detected trigger (e.g. sudden 421 4.7.0 rate limit).
    - "Jump to Point of Failure" CTA button that immediately seeks playback to $T - 60$ seconds prior to the incident.
    - Visual shockwave pulse animation highlighting the incident point when scrubbing passes over it.
  - **Zero-CDN Invariant:** CSS keyframes with `transform: scale()` and pure SVG marker vectors.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog incident seek --domain "outbound.biz" --threshold 0.1`
  - `gitmap warmup timelog incident list --domain "outbound.biz"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["incidentId", "ispIdentifier", "incidentEpoch", "reputationDropPercent", "isActive", "hasAutoSeeked"],
    "properties": {
      "incidentId": { "type": "string" },
      "ispIdentifier": { "type": "string" },
      "incidentEpoch": { "type": "integer" },
      "reputationDropPercent": { "type": "number" },
      "triggerReason": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasAutoSeeked": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Mathematical gradient analysis detects reputation drop-off points with 100% deterministic accuracy.
  - Clicking an incident pin seeks time scrubber to 60 seconds prior to incident in < 15ms.
  - Visual shockwave pulse renders smoothly without DOM layout reflow.

---

### Improvement 214: Multi-Track Deliverability Event Sequencer View
- **ID:** `WARMUP-TIME-214`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Outbound sending involves concurrent asynchronous systems (DNS updates, MTA queue flushes, seed opens, spam bounces) that are impossible to correlate on a single timeline track.
- **Concrete UI Specification:**
  - **Layout:** Digital Audio Workstation (DAW) style multi-track event sequencer (`.sequencer-viewport`).
  - **Visual Elements:**
    - Separate horizontal tracks for:
      - Track 1: Outbound Dispatch Bursts
      - Track 2: ISP Throttle Signals (421/451 SMTP codes)
      - Track 3: Seed Inbox Placement Events
      - Track 4: Seed Spam Folder Events
      - Track 5: DNS & IP Reputation Shifts
    - Colored event block clips on each track showing start time, duration, and volume density.
    - Synchronized vertical playhead needle spanning all 5 tracks simultaneously.
  - **Zero-CDN Invariant:** Pure CSS flexbox layout with timeline tracks synchronized via single CSS custom property `--playhead-x`.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog sequencer --domain "acme.mail" --window "12h"`
  - `gitmap warmup timelog sequencer export-csv --domain "acme.mail"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["sequencerId", "trackName", "events", "isActive", "hasOverlappingSpikes"],
    "properties": {
      "sequencerId": { "type": "string" },
      "trackName": { "type": "string", "enum": ["dispatch", "throttles", "inbox", "spam", "reputation"] },
      "isActive": { "type": "boolean" },
      "hasOverlappingSpikes": { "type": "boolean" },
      "events": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["startEpoch", "durationSeconds", "densityScore"],
          "properties": {
            "startEpoch": { "type": "integer" },
            "durationSeconds": { "type": "integer" },
            "densityScore": { "type": "number" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Playhead needle position tracks smoothly across all 5 tracks without desynchronization.
  - Hovering over any event block reveals cross-track correlated events occurring within ±30 seconds.
  - Track zooming supports levels from 1 hour to 14 days without DOM rebuilding.

---

### Improvement 215: Time-Slice SMTP Response Code Heat Stream
- **ID:** `WARMUP-TIME-215`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Fleeting SMTP error codes (e.g., transient Greylisting 450 vs hard 550 mailbox unavailable) disappear into aggregate bounce rates, hiding early warnings of warming failure.
- **Concrete UI Specification:**
  - **Layout:** Continuous vertical heat stream band (`.smtp-heat-stream`) displaying color-density bands for 2xx, 4xx, and 5xx response status codes.
  - **Visual Elements:**
    - Green streams for 250 OK success handshakes.
    - Orange streams for 421/450/451 transient delays and rate limits.
    - Red streams for 550/554 permanent rejections and spam blocks.
    - Granular minute-by-minute slice inspection modal when clicking on any stream segment.
  - **Zero-CDN Invariant:** HTML5 `<canvas>` rendering stream bands with native 2D canvas context and zero external charting libraries.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog smtp-codes --domain "growth.io" --granularity 1m`
  - `gitmap warmup timelog stream-stats --code-prefix 4`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["streamId", "epochMinute", "count2xx", "count4xx", "count5xx", "isActive", "hasTransientSurge"],
    "properties": {
      "streamId": { "type": "string" },
      "epochMinute": { "type": "integer" },
      "count2xx": { "type": "integer", "minimum": 0 },
      "count4xx": { "type": "integer", "minimum": 0 },
      "count5xx": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasTransientSurge": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Canvas updates smoothly during time-scrubbing at 60 FPS.
  - Accurately classifies and separates 4xx transient throttling from 5xx fatal rejects.
  - Provides instant drill-down to raw MTA bounce strings on click.

---

### Improvement 216: Simulated What-If Warming Strategy Branching Sandbox
- **ID:** `WARMUP-TIME-216`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Operators facing a reputation dip cannot test hypothetical corrective actions (e.g. cutting send volume by 50% vs doubling seed engagement) without risking live sending domains.
- **Concrete UI Specification:**
  - **Layout:** Split-screen branch comparison workbench (`.branch-sandbox-viewport`).
  - **Visual Elements:**
    - "Original Reality" historical trajectory curve vs "What-If Branch" simulated projection curve.
    - Interactive parameter knobs to simulate: Adjust Volume Slope (-50% to +100%), Double Seed Reply Velocity, Change Sending Schedule to Business Hours Only.
    - AI-projected reputation outcome indicator displaying expected recovery date and confidence score.
    - "Commit Branch as Active Plan" button that transforms the simulated curve into live system configuration.
  - **Zero-CDN Invariant:** Client-side polynomial projection engine computed in vanilla JavaScript.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog branch create --fork-epoch 1704200000 --adjust-volume -30`
  - `gitmap warmup timelog branch compare --branch-id "recovery-sim-01"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["branchId", "forkEpoch", "volumeAdjustmentFactor", "projectedRecoveryDays", "isActive", "hasHigherReputation"],
    "properties": {
      "branchId": { "type": "string" },
      "forkEpoch": { "type": "integer" },
      "volumeAdjustmentFactor": { "type": "number" },
      "seedEngagementMultiplier": { "type": "number" },
      "projectedRecoveryDays": { "type": "number" },
      "confidenceScore": { "type": "number", "minimum": 0, "maximum": 1 },
      "isActive": { "type": "boolean" },
      "hasHigherReputation": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Computes 30-day projection curves across modified parameters in < 25ms.
  - Visual comparison overlays historical data against simulation with distinct dashed lines.
  - Committing a branch writes new parameters directly to `.gitmap/warmup.db`.

---

### Improvement 217: Time-Travel Deliverability Forensic Audit Report Exporter
- **ID:** `WARMUP-TIME-217`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Delivering proof of deliverability health or post-mortem incident reports to clients and leadership requires tedious manual screenshotting and fragmented spreadsheet exports.
- **Concrete UI Specification:**
  - **Layout:** Self-contained printable audit report preview drawer (`.report-preview-drawer`).
  - **Visual Elements:**
    - Executive summary scorecard with verified cryptographic timestamp hashes.
    - Embedded SVG charts showing volume progression, placement ratios, and incident markers.
    - Full MTA timeline transcript highlighting key corrective actions taken by autonomous warming daemon.
    - One-click "Export Standalone HTML Report" button generating a single-file, 100% offline-viewable HTML package.
  - **Zero-CDN Invariant:** Generated reports use embedded inline CSS, SVG, and base64 encoded font vectors with zero external references.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog export report --domain "brand.com" --start "2026-10-01" --output "report.html"`
  - `gitmap warmup timelog export json --domain "brand.com"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["reportId", "domainName", "startDate", "endDate", "totalEmailsSent", "avgInboxRate", "isActive", "hasCertifiedHash"],
    "properties": {
      "reportId": { "type": "string" },
      "domainName": { "type": "string" },
      "startDate": { "type": "string" },
      "endDate": { "type": "string" },
      "totalEmailsSent": { "type": "integer" },
      "avgInboxRate": { "type": "number" },
      "isActive": { "type": "boolean" },
      "hasCertifiedHash": { "type": "boolean" },
      "certificationHash": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Exported HTML report opens completely offline in any modern browser without broken styles.
  - Print stylesheet (`@media print`) formats tables and charts into high-resolution pagination.
  - CLI command generates report file in < 200ms.

---

### Improvement 218: Time-Synchronized Seed Mailbox State Reconstructor
- **ID:** `WARMUP-TIME-218`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** When investigating an inboxing failure, operators can only see current seed inbox contents, having no way to view what a specific seed inbox looked like 3 days ago.
- **Concrete UI Specification:**
  - **Layout:** Virtual seed inbox mock window (`.seed-inbox-mockup`) reflecting past mail states.
  - **Visual Elements:**
    - Traditional email client folder tree: Inbox (unread count), Spam, Trash, Archive.
    - Message list view showing past emails from target domain with subject line, sender, and received time.
    - Visual badge indicating whether the email was originally placed in Spam and then programmatically moved to Inbox by the seed persona.
    - Full email body preview pane showing rendered HTML content, headers, and authentication pass/fail tags.
  - **Zero-CDN Invariant:** Pure HTML/CSS simulated webmail client styled with Obsidian glassmorphic skin.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog seed-view --seed "ceo-seed@partner.org" --epoch 1704280000`
  - `gitmap warmup seed inspect --id "seed-042"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["reconstructionId", "seedEmail", "targetEpoch", "folderLocation", "isActive", "hasBeenRescuedFromSpam"],
    "properties": {
      "reconstructionId": { "type": "string" },
      "seedEmail": { "type": "string" },
      "targetEpoch": { "type": "integer" },
      "folderLocation": { "type": "string", "enum": ["inbox", "spam", "trash", "promotions"] },
      "isActive": { "type": "boolean" },
      "hasBeenRescuedFromSpam": { "type": "boolean" },
      "emailHeaders": {
        "type": "object",
        "required": ["subject", "from", "date"],
        "properties": {
          "subject": { "type": "string" },
          "from": { "type": "string" },
          "date": { "type": "string" }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately renders historical seed mailbox contents for any recorded epoch.
  - Displays explicit visual indicator if the seed persona engaged in synthetic rescues (moved spam to inbox).
  - Sanitizes rendered email HTML to prevent script execution vulnerabilities.

---

### Improvement 219: Delivery Latency Jitter & Connection Queue Time-Lapse
- **ID:** `WARMUP-TIME-219`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Deliverability bottlenecks often originate from SMTP handshake latency jitter and TCP queue stalls that are invisible in high-level throughput charts.
- **Concrete UI Specification:**
  - **Layout:** High-precision latency jitter histogram (`.latency-jitter-chart`).
  - **Visual Elements:**
    - Microsecond/millisecond latency buckets (0–50ms, 50–200ms, 200–500ms, 500ms–2s, >2s).
    - Color ramp shifting from cool cyan (sub-100ms) to hot crimson (>2000ms latency).
    - Time-lapse playback showing shifting latency distribution curves as network congestion develops.
    - TCP connection queue fill percentage gauge with active thread counts.
  - **Zero-CDN Invariant:** SVG path spline calculation with hardware-accelerated CSS opacity transitions.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog latency --gateway "relay.core" --window "4h"`
  - `gitmap warmup gateway latency-test --host "mta.google.com"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["latencyId", "gatewayHost", "medianLatencyMs", "p99LatencyMs", "isActive", "hasJitterSpike"],
    "properties": {
      "latencyId": { "type": "string" },
      "gatewayHost": { "type": "string" },
      "medianLatencyMs": { "type": "number" },
      "p99LatencyMs": { "type": "number" },
      "isActive": { "type": "boolean" },
      "hasJitterSpike": { "type": "boolean" },
      "latencyBuckets": {
        "type": "array",
        "items": { "type": "integer" }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Renders P50, P90, and P99 latency percentiles with precision.
  - Detects network latency jitter spikes (> 3x median) and marks corresponding time tick.
  - Renders 100 latency sample frames during playback without UI hitching.

---

### Improvement 220: Temporal Incident Bookmark & Annotation Collaboration Hub
- **ID:** `WARMUP-TIME-220`
- **Cluster:** 22 (Time-Travel Deliverability Debugger & Event Playback)
- **Problem Statement:** Deliverability engineers investigating historical anomalies cannot leave persistent comments or shared markers for teammates at specific timeline coordinates.
- **Concrete UI Specification:**
  - **Layout:** Timeline annotation sidebar drawer (`.annotation-sidebar`).
  - **Visual Elements:**
    - List of collaborative comment pins linked to specific timestamps.
    - Add Annotation modal allowing engineers to tag teammates, add Markdown explanations, and assign remediation tasks.
    - Hovering an annotation highlights the corresponding point on the scrubber bar.
    - Status badge: Open, Investigating, Resolved.
  - **Zero-CDN Invariant:** Native DOM dialog element (`<dialog>`) with glassmorphic backdrop styling.
- **CLI Counterpart Command:**
  - `gitmap warmup timelog annotate add --epoch 1704291000 --message "Throttled by Yahoo due to burst"`
  - `gitmap warmup timelog annotate list --domain "outbound.net"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["annotationId", "epochTimestamp", "authorUsername", "commentText", "isActive", "hasTaskAssigned"],
    "properties": {
      "annotationId": { "type": "string" },
      "epochTimestamp": { "type": "integer" },
      "authorUsername": { "type": "string" },
      "commentText": { "type": "string" },
      "status": { "type": "string", "enum": ["open", "investigating", "resolved"] },
      "isActive": { "type": "boolean" },
      "hasTaskAssigned": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Saving an annotation immediately updates timeline markers across all active sessions.
  - Clicking an annotation automatically jumps the scrubber playhead to that exact timestamp.
  - Persists annotations to SQLite Split-DB table `WarmupTimelineAnnotations`.

---

# Cluster 23: Adaptive Warming Gradient & Reputation Acceleration Engine (Items 221–230)

---

### Improvement 221: Dynamic Logistic Sigmoid Volume Curvature Calculator
- **ID:** `WARMUP-CURVE-221`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Fixed linear warming schedules (e.g. +50 emails every day) ignore ISP trust dynamics, either warming too aggressively in the fragile first week or unnecessarily slowly in later weeks.
- **Concrete UI Specification:**
  - **Layout:** Curve parameter tuning console (`.curve-tuner-panel`).
  - **Visual Elements:**
    - Live mathematical formula display rendering current logistic sigmoid equation: $V(t) = \frac{K}{1 + Q \cdot e^{-B(t - M)}}$.
    - Interactive slider handles for Carrying Capacity ($K$), Growth Rate ($B$), and Midpoint Day ($M$).
    - Real-time curve graph updating instantaneously as sliders move.
    - Projected time-to-peak KPI chip (e.g. `Target 50,000/day achieved on Day 24`).
  - **Zero-CDN Invariant:** Dynamic SVG path generation using standard trigonometric and exponential functions in pure JS.
- **CLI Counterpart Command:**
  - `gitmap warmup curve calculate --model sigmoid --capacity 50000 --growth 0.25 --midpoint 14`
  - `gitmap warmup curve set --domain "send.saas.com" --model sigmoid`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["curveConfigId", "modelType", "carryingCapacity", "growthRate", "isActive", "hasOptimalParameters"],
    "properties": {
      "curveConfigId": { "type": "string" },
      "modelType": { "type": "string", "enum": ["logistic_sigmoid", "exponential", "stepped"] },
      "carryingCapacity": { "type": "integer", "minimum": 100 },
      "growthRate": { "type": "number", "minimum": 0.01, "maximum": 2.0 },
      "midpointDay": { "type": "integer", "minimum": 1 },
      "isActive": { "type": "boolean" },
      "hasOptimalParameters": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Calculates daily volume targets with smooth monotonic increase and zero negative slopes.
  - Slider adjustments re-render the 30-day projection path in < 8ms.
  - Generates mathematically sound day-by-day volume schedule stored directly in `.gitmap/warmup.db`.

---

### Improvement 222: Automated ISP Throttling Deceleration Brake Engine
- **ID:** `WARMUP-CURVE-222`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** When an ISP issues temporary rate limits (421 4.7.0), static warming systems continue blasting messages, turning transient greylisting into hard IP blocklist listings.
- **Concrete UI Specification:**
  - **Layout:** Throttle brake telemetry card (`.throttle-brake-card`) with dynamic brake engagement status.
  - **Visual Elements:**
    - Brake activation meter with animated disc brake caliper SVG graphic that clamps shut when throttle threshold is exceeded.
    - Deceleration ratio readout (e.g. `Brake Engaged: Volume Scaled by -40%`).
    - Recovery cooldown timer countdown (e.g. `Cooldown: 03h 42m remaining until ramp resumption`).
    - Provider impact chips showing which specific ISPs triggered the brake (e.g. Gmail 100% OK, Outlook Brake Active).
  - **Zero-CDN Invariant:** Pure CSS transform animation with variable deceleration dampening.
- **CLI Counterpart Command:**
  - `gitmap warmup curve brake status --domain "outbound.ai"`
  - `gitmap warmup curve brake trigger --domain "outbound.ai" --isp "microsoft" --factor 0.5`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["brakeId", "domainName", "isBrakeEngaged", "decelerationFactor", "cooldownRemainingSeconds", "hasHardBlockRisk"],
    "properties": {
      "brakeId": { "type": "string" },
      "domainName": { "type": "string" },
      "isBrakeEngaged": { "type": "boolean" },
      "decelerationFactor": { "type": "number", "minimum": 0.1, "maximum": 0.9 },
      "cooldownRemainingSeconds": { "type": "integer", "minimum": 0 },
      "hasHardBlockRisk": { "type": "boolean" },
      "triggeringIsp": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Automatically engages within 500ms when 4xx error rate from a single ISP exceeds 5% in a 10-minute window.
  - Dynamically scales down that ISP's target volume by the designated `decelerationFactor`.
  - Safely resumes gradual acceleration once the cooldown window expires with zero subsequent 4xx errors.

---

### Improvement 223: Multi-Phase Plateau Warming Stabilization Model
- **ID:** `WARMUP-CURVE-223`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Continuous daily ramping without volume plateaus prevents spam filter statistical baselines from settling, increasing the risk of sudden spam folder dumping.
- **Concrete UI Specification:**
  - **Layout:** Step-and-Plateau curve graph (`.plateau-curve-graph`).
  - **Visual Elements:**
    - Stepped curve profile showing ramp phases followed by flat plateau holding periods (e.g. Ramp 3 days -> Hold 2 days -> Ramp 3 days).
    - Plateau stability gauge measuring reputation variance during the hold period.
    - "Plateau Cleared" green checkpoint badge when inbox placement remains > 95% across the entire hold period.
    - Dynamic adjustment handle allowing operators to extend plateau duration if metrics fluctuate.
  - **Zero-CDN Invariant:** SVG path generation with step-after interpolation (`d="M... L... H... V..."`).
- **CLI Counterpart Command:**
  - `gitmap warmup curve plateau configure --domain "mail.cloud" --ramp-days 3 --hold-days 2`
  - `gitmap warmup curve plateau status --domain "mail.cloud"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["plateauConfigId", "rampDays", "holdDays", "isInHoldPeriod", "hasAchievedStability", "currentPlateauIndex"],
    "properties": {
      "plateauConfigId": { "type": "string" },
      "rampDays": { "type": "integer", "minimum": 1 },
      "holdDays": { "type": "integer", "minimum": 1 },
      "isInHoldPeriod": { "type": "boolean" },
      "hasAchievedStability": { "type": "boolean" },
      "currentPlateauIndex": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Volume remains strictly clamped to constant target during designated hold days.
  - Automatically advances to next ramp phase only when stability criteria are satisfied.
  - Displays visual indicator distinguishing active ramp days from stabilization days.

---

### Improvement 224: Seed Engagement Velocity Reputation Multiplier Engine
- **ID:** `WARMUP-CURVE-224`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Static seed networks interact with emails uniformly, but modern ISP machine learning algorithms heavily weight high-velocity interactions (immediate opens, multi-turn replies, stars, address book additions).
- **Concrete UI Specification:**
  - **Layout:** Engagement Velocity Matrix (`.velocity-matrix-panel`).
  - **Visual Elements:**
    - Velocity dial displaying real-time engagement score multiplier ($1.0x$ to $3.5x$).
    - Interaction breakdown chips:
      - Rapid Open (< 60s): +0.2x multiplier
      - Two-Way Thread Reply: +0.5x multiplier
      - Starred / Marked Important: +0.4x multiplier
      - Moved from Spam to Inbox: +0.8x multiplier
    - Acceleration gauge indicating how much faster the warming curve can proceed due to elevated seed engagement.
  - **Zero-CDN Invariant:** CSS radial progress indicator with variable gradient glow.
- **CLI Counterpart Command:**
  - `gitmap warmup curve velocity status --domain "outbound.dev"`
  - `gitmap warmup curve velocity boost --multiplier 1.5`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["velocityId", "currentMultiplier", "rapidOpenCount", "threadReplyCount", "isActive", "hasBonusAcceleration"],
    "properties": {
      "velocityId": { "type": "string" },
      "currentMultiplier": { "type": "number", "minimum": 1.0, "maximum": 5.0 },
      "rapidOpenCount": { "type": "integer", "minimum": 0 },
      "threadReplyCount": { "type": "integer", "minimum": 0 },
      "unspamRescueCount": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasBonusAcceleration": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Seed engagement events contribute to multiplier calculation with verifiable mathematical weights.
  - Elevating multiplier safely shortens total warming schedule by up to 35% without triggering spam alarms.
  - Multiplier decays gradually if engagement actions taper off.

---

### Improvement 225: ISP-Specific Independent Warming Curve Partitioner
- **ID:** `WARMUP-CURVE-225`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Warming an IP pool with a single aggregate schedule fails because Google, Microsoft, and Yahoo have vastly different acceptance tolerances and warming speeds.
- **Concrete UI Specification:**
  - **Layout:** Multi-curve partitioner dashboard (`.isp-curve-partitioner`).
  - **Visual Elements:**
    - Overlay chart displaying 4 distinct, independently calculated warming curves:
      - Google Workspace Curve (Cyan)
      - Microsoft 365 Curve (Violet)
      - Yahoo/AOL Curve (Emerald)
      - Generic MTAs (Amber)
    - Independent pace sliders for each ISP allowing custom volume caps and acceleration rates.
    - Composite total volume bar at the bottom showing the sum of all individual ISP targets.
  - **Zero-CDN Invariant:** Multi-path SVG visualization rendered using native SVG `<g>` groups and scoped CSS variables.
- **CLI Counterpart Command:**
  - `gitmap warmup curve isp-split --domain "scale.org"`
  - `gitmap warmup curve isp-set --isp "google" --target 20000 --days 18`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["partitionId", "domainName", "ispCurves", "isActive", "hasDivergentSchedules"],
    "properties": {
      "partitionId": { "type": "string" },
      "domainName": { "type": "string" },
      "isActive": { "type": "boolean" },
      "hasDivergentSchedules": { "type": "boolean" },
      "ispCurves": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["ispName", "targetVolume", "currentDayVolume", "isActive"],
          "properties": {
            "ispName": { "type": "string" },
            "targetVolume": { "type": "integer" },
            "currentDayVolume": { "type": "integer" },
            "isActive": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Outbound dispatch scheduler enforces volume limits strictly per target ISP domain MX records.
  - Modifying one ISP curve does not disrupt or alter the schedules of other partitioned ISPs.
  - Aggregate total volume updates reactively when any individual ISP curve is modified.

---

### Improvement 226: Machine Learning Spam Trap Proximity Deceleration Sensor
- **ID:** `WARMUP-CURVE-226`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Hitting a pristine spam trap can instantly burn a new IP range, but standard warming curves have no early warning system to decelerate when recipient address risk increases.
- **Concrete UI Specification:**
  - **Layout:** Spam Trap Radar HUD (`.spam-trap-radar-hud`).
  - **Visual Elements:**
    - Circular radar scanner animation (`.radar-sweep`) sweeping across recipient address hygiene domains.
    - Threat proximity index gauge (0.0 to 10.0) measuring recipient MX typography, abandoned domain patterns, and role-based addresses.
    - Warning alert banner showing "High Risk Address Density Detected: Automatic Curve Deceleration Active".
    - Address filter toggle to automatically isolate and quarantine suspect recipients before dispatch.
  - **Zero-CDN Invariant:** Pure CSS radial gradient with keyframe sweep animation.
- **CLI Counterpart Command:**
  - `gitmap warmup curve trap-sensor status --domain "outbound.net"`
  - `gitmap warmup curve trap-sensor scan --recipient-list "leads.csv"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["sensorId", "threatScore", "quarantinedCount", "isDecelerationTriggered", "hasTrapProximityRisk"],
    "properties": {
      "sensorId": { "type": "string" },
      "threatScore": { "type": "number", "minimum": 0, "maximum": 10 },
      "quarantinedCount": { "type": "integer", "minimum": 0 },
      "isDecelerationTriggered": { "type": "boolean" },
      "hasTrapProximityRisk": { "type": "boolean" },
      "lastDetectionEpoch": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Automatically intercepts addresses with known spam trap syntax (e.g. role-accounts on dead domains).
  - Triggers warming curve deceleration when threat score exceeds 6.5.
  - Quarantine actions are logged with full audit trails in SQLite `WarmupTrapAlerts`.

---

### Improvement 227: Time-of-Day Natural Human Sending Distribution Curve
- **ID:** `WARMUP-CURVE-227`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Dispatching warmup emails in flat, mechanical bursts (e.g. 1 email every 30 seconds 24/7) triggers behavioral bot heuristics at major email providers.
- **Concrete UI Specification:**
  - **Layout:** Diurnal Sending Distribution bell curve chart (`.diurnal-curve-chart`).
  - **Visual Elements:**
    - 24-hour horizontal axis representing local recipient business hours.
    - Gaussian bell curve envelope peaking at 10:00 AM and 2:30 PM with natural drop-offs during lunch and night hours.
    - Poisson distribution randomness toggle ("Human Jitter Engine") adding randomized micro-delays between individual dispatches.
    - Timezone alignment selector aligning outbound send waves with recipient geographic centers (Americas, EMEA, APAC).
  - **Zero-CDN Invariant:** Calculated SVG `<path>` plotting Gaussian normal distribution: $f(x) = \frac{1}{\sigma \sqrt{2\pi}} e^{-\frac{1}{2}\left(\frac{x - \mu}{\sigma}\right)^2}$.
- **CLI Counterpart Command:**
  - `gitmap warmup curve diurnal generate --timezone "America/New_York" --jitter true`
  - `gitmap warmup curve diurnal preview --hours 24`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["diurnalConfigId", "targetTimezone", "hasHumanJitter", "peakHourMorning", "peakHourAfternoon", "isActive"],
    "properties": {
      "diurnalConfigId": { "type": "string" },
      "targetTimezone": { "type": "string" },
      "hasHumanJitter": { "type": "boolean" },
      "peakHourMorning": { "type": "number", "minimum": 6, "maximum": 12 },
      "peakHourAfternoon": { "type": "number", "minimum": 13, "maximum": 18 },
      "isActive": { "type": "boolean" },
      "hourlyVolumeDistribution": {
        "type": "array",
        "items": { "type": "number" }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Volume curves naturally drop to near zero during nighttime hours (11:00 PM - 6:00 AM).
  - Random Poisson jitter prevents static time-interval patterns between successive SMTP connections.
  - Dispatches align precisely with configured recipient timezone curves.

---

### Improvement 228: Autonomous DNS Alignment Verification & Ramp Gatekeeper
- **ID:** `WARMUP-CURVE-228`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** If a domain's SPF or DKIM DNS record fails or gets accidentally deleted mid-warmup, the ramp continues blindly, severely damaging domain reputation.
- **Concrete UI Specification:**
  - **Layout:** DNS Ramp Gatekeeper HUD (`.dns-gatekeeper-hud`).
  - **Visual Elements:**
    - Live DNS verification radar checking SPF TXT, DKIM Selector, DMARC Policy, and Reverse DNS (PTR) every 60 seconds.
    - Security Interlock status badge: "Gatekeeper Armed: All Records Aligned" (Green) vs "Gatekeeper Tripped: Ramp Halted" (Red Pulse).
    - Emergency Ramp Freeze overlay preventing any new outbound dispatch if DNS resolution fails.
    - Instant DNS re-check button with animated rotation spinner.
  - **Zero-CDN Invariant:** SVG shield status icon with pure CSS border-glow transitions.
- **CLI Counterpart Command:**
  - `gitmap warmup curve gatekeeper check --domain "outbound.pro"`
  - `gitmap warmup curve gatekeeper interlock --arm true`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["gatekeeperId", "domainName", "isInterlockArmed", "isRampHalted", "hasValidSpf", "hasValidDkim", "hasValidDmarc", "hasValidPtr"],
    "properties": {
      "gatekeeperId": { "type": "string" },
      "domainName": { "type": "string" },
      "isInterlockArmed": { "type": "boolean" },
      "isRampHalted": { "type": "boolean" },
      "hasValidSpf": { "type": "boolean" },
      "hasValidDkim": { "type": "boolean" },
      "hasValidDmarc": { "type": "boolean" },
      "hasValidPtr": { "type": "boolean" },
      "lastCheckEpoch": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Immediately freezes warming volume progression within 5 seconds of DNS validation failure.
  - Dispatches high-priority notification to operator action bar.
  - Automatically unfreezes and resumes curve ramp once DNS records return valid verification responses.

---

### Improvement 229: Multi-IP Pool Rotational Balancing & Warmth Equalizer
- **ID:** `WARMUP-CURVE-229`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Large outbound infrastructures utilize multiple dedicated IP addresses, but unbalanced sending causes one IP to warm prematurely while others remain cold and untrusted.
- **Concrete UI Specification:**
  - **Layout:** IP Pool Equalizer mixing console (`.ip-equalizer-console`).
  - **Visual Elements:**
    - Vertical channel faders for each IP address in the pool (e.g. `198.51.100.1` to `198.51.100.8`).
    - Channel LED volume meter indicating current daily email throughput.
    - "Auto-Balance Pool" master switch distributing dispatch load inversely proportional to each IP's current warmth score.
    - Pool aggregate warmth score gauge (0–100%) tracking overall fleet readiness.
  - **Zero-CDN Invariant:** CSS flexbox fader rails with pure SVG knob indicators and gradient LED strips.
- **CLI Counterpart Command:**
  - `gitmap warmup curve pool balance --pool "us-east-pool" --mode auto`
  - `gitmap warmup curve pool status --pool "us-east-pool"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["poolId", "poolName", "totalIpsCount", "isAutoBalancing", "hasImbalancedIps", "aggregateWarmthScore"],
    "properties": {
      "poolId": { "type": "string" },
      "poolName": { "type": "string" },
      "totalIpsCount": { "type": "integer", "minimum": 1 },
      "isAutoBalancing": { "type": "boolean" },
      "hasImbalancedIps": { "type": "boolean" },
      "aggregateWarmthScore": { "type": "number", "minimum": 0, "maximum": 100 },
      "ipAllocations": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["ipAddress", "currentVolume", "warmthScore", "isActive"],
          "properties": {
            "ipAddress": { "type": "string" },
            "currentVolume": { "type": "integer" },
            "warmthScore": { "type": "number" },
            "isActive": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Dynamic load balancing router shifts volume towards cooler IPs to normalize warmth across the entire pool.
  - Prevents any single IP from exceeding its daily warming curve ceiling.
  - Visual equalizer meters update in real-time as dispatches rotate across interfaces.

---

### Improvement 230: Post-Warmup Peak Maintenance & Cold-Start Re-Warming Profiler
- **ID:** `WARMUP-CURVE-230`
- **Cluster:** 23 (Adaptive Warming Gradient & Reputation Acceleration Engine)
- **Problem Statement:** Once a domain completes its warming ramp, volume typically spikes or drops abruptly during seasonal quiet periods, destroying established ISP reputation.
- **Concrete UI Specification:**
  - **Layout:** Maintenance & Re-Warming Profile manager (`.maintenance-profile-panel`).
  - **Visual Elements:**
    - Operational mode toggle: Active Maintenance vs Cold-Start Re-Warming.
    - Minimum daily volume baseline floor indicator ensuring the domain never sends below reputation threshold.
    - Autonomous "Keep-Alive" seed burst scheduler injecting synthetic engagement traffic during quiet business days.
    - Re-warming ramp accelerator designed specifically for returning dormant domains to peak volume in 50% less time.
  - **Zero-CDN Invariant:** Obsidian glassmorphic card deck with custom SVG mode badges and status switches.
- **CLI Counterpart Command:**
  - `gitmap warmup curve maintenance set --domain "corp.com" --baseline-floor 2500`
  - `gitmap warmup curve rewarm start --domain "corp.com" --dormant-days 45`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["profileId", "domainName", "operationalMode", "baselineVolumeFloor", "isActive", "hasDormancyRisk"],
    "properties": {
      "profileId": { "type": "string" },
      "domainName": { "type": "string" },
      "operationalMode": { "type": "string", "enum": ["peak_maintenance", "cold_start_rewarming", "dormant_idle"] },
      "baselineVolumeFloor": { "type": "integer", "minimum": 0 },
      "isActive": { "type": "boolean" },
      "hasDormancyRisk": { "type": "boolean" },
      "daysDormant": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Enforces minimum baseline floor volume through autonomous seed exchange during low-volume sending days.
  - Cold-start re-warming profile safely ramps dormant domains back to 100% capacity within 7 days.
  - Automatically raises alert if domain dormancy exceeds 14 consecutive days.
