# Specification 337: Improvements 01 to 30 (Warmup Telemetry & Scheduler)

- **Sequence:** 337
- **Slug:** `337-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/02-improvements-01-to-30-warmup-telemetry-and-scheduler.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## Improvement Catalog Summary (01 - 30)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 1** | Items 01–10 | Warm-Up Ramp Curves & Visual Schedule Timeline |
| **Cluster 2** | Items 11–20 | Deliverability Health Gauges & ISP Inbox/Spam Placement |
| **Cluster 3** | Items 21–30 | Seed List Scheduler, Reply Simulation & Automated Interaction Physics |

---

# Cluster 1: Warm-Up Ramp Curves & Visual Schedule Timeline (Items 01–10)

---

### Improvement 01: Interactive Bezier S-Curve Ramp Visualizer
- **ID:** `WARMUP-UI-01`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Operators currently view warm-up progress as dry ASCII tables without visual understanding of mathematical ramp progression, making it impossible to anticipate volume inflection points or detect schedule stalls.
- **Concrete UI Specification:**
  - **Layout:** A 16:9 responsive SVG viewport (`.ramp-bezier-card`) rendering current vs. projected daily email send quotas across a 30-day timeline.
  - **Visual Elements:**
    - SVG cubic Bezier spline (`stroke: #3b82f6; stroke-width: 2.5; fill: url(#rampGradient);`).
    - Ambient gradient fill fading to `#0f172a`.
    - Interactive draggable control handles (`circle.control-knot`, `r=6`, `fill: #60a5fa`) allowing operators to adjust ramp steepness dynamically.
    - Floating zero-CDN HTML tooltip updating `Quota: X emails/day` and `Date: YYYY-MM-DD` on mouseover.
    - Smooth CSS transition `transition: d 350ms cubic-bezier(0.4, 0, 0.2, 1);`.
  - **Zero-CDN Invariant:** Pure SVG vector canvas with inline CSS math; no D3, Chart.js, or external font stylesheets.
- **CLI Counterpart Command:**
  - `gitmap warmup curve show --domain "acme-corp.com" --format svg`
  - `gitmap warmup curve set --domain "acme-corp.com" --curve-type "s-curve" --steepness 1.4`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["domain", "curveType", "isActive", "hasCompleted", "points"],
    "properties": {
      "domain": { "type": "string" },
      "curveType": { "type": "string", "enum": ["linear", "exponential", "s-curve", "fibonacci"] },
      "isActive": { "type": "boolean" },
      "hasCompleted": { "type": "boolean" },
      "steepnessFactor": { "type": "number", "minimum": 0.1, "maximum": 5.0 },
      "points": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["dayIndex", "targetQuota", "actualSent", "isMilestone"],
          "properties": {
            "dayIndex": { "type": "integer", "minimum": 1 },
            "targetQuota": { "type": "integer", "minimum": 0 },
            "actualSent": { "type": "integer", "minimum": 0 },
            "isMilestone": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Bezier path re-renders without DOM remount in `< 16ms`.
  - Path handles clamp between minimum 5 emails/day and maximum 5,000 emails/day.
  - CLI prints clean ANSI terminal curve when `--format ascii` is specified.

---

### Improvement 02: Multi-Domain Ramp Overlay Comparison Deck
- **ID:** `WARMUP-UI-02`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Managing multiple sender domains requires toggling between separate screens, hiding cross-domain synchronization issues and cluster-wide ISP reputation risks.
- **Concrete UI Specification:**
  - **Layout:** Multi-layer SVG graph displaying up to 8 sender domains simultaneously with color-coded distinct strokes (`#3b82f6`, `#10b981`, `#f59e0b`, `#8b5cf6`, `#ec4899`, etc.).
  - **Visual Elements:**
    - Domain toggle pill badges (`.domain-pill`) with active checkboxes and glow borders.
    - Synchronized vertical crosshair scanning across all domains on X-axis hover.
    - Cumulative aggregate volume line (`stroke-dasharray: 4 2; stroke: #94a3b8;`).
  - **Zero-CDN Invariant:** Inline SVG multi-path container with scoped CSS variables for stroke palettes.
- **CLI Counterpart Command:**
  - `gitmap warmup compare --domains "alpha.com,beta.com,gamma.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["clusterId", "isSynchronized", "domains"],
    "properties": {
      "clusterId": { "type": "string" },
      "isSynchronized": { "type": "boolean" },
      "domains": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["domainName", "hexColor", "isActive", "dailyCap"],
          "properties": {
            "domainName": { "type": "string" },
            "hexColor": { "type": "string" },
            "isActive": { "type": "boolean" },
            "dailyCap": { "type": "integer" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Toggling domain pills animates visibility via `opacity: 0` without layout shifts.
  - Aggregate line updates calculation immediately upon domain toggle.

---

### Improvement 03: Hourly Send Velocity Heatmap & Jitter Timeline
- **ID:** `WARMUP-UI-03`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Outbound emails sent in robotic fixed-interval bursts trigger ISP spam filters. Users need visibility into micro-jitter and diurnal hourly distribution.
- **Concrete UI Specification:**
  - **Layout:** 24-hour horizontal grid (00:00 to 23:00) with 60 micro-cells per hour representing dispatch jitter.
  - **Visual Elements:**
    - Cells shaded with 5-step monochromatic blues: `rgba(59, 130, 246, alpha)` based on velocity density.
    - Diurnal human-emulation wave overlay (bell curve matching 08:00 - 18:00 recipient local time).
    - Jitter scatter dots (`circle.jitter-dot`, `r=2`) illustrating microsecond variance.
  - **Zero-CDN Invariant:** CSS Grid matrix with native CSS custom property opacity shading.
- **CLI Counterpart Command:**
  - `gitmap warmup jitter show --domain "acme.com" --date "today"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "hasDiurnalMask", "isJitterEnabled", "hourlyBuckets"],
    "properties": {
      "domain": { "type": "string" },
      "hasDiurnalMask": { "type": "boolean" },
      "isJitterEnabled": { "type": "boolean" },
      "hourlyBuckets": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["hour", "sentCount", "hasVariance"],
          "properties": {
            "hour": { "type": "integer", "minimum": 0, "maximum": 23 },
            "sentCount": { "type": "integer" },
            "hasVariance": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Hourly heat cells highlight active hour with pulsing blue ring (`animation: pulse 2s infinite`).
  - Tooltip shows average inter-send delay and standard deviation.

---

### Improvement 04: Dynamic Backoff & Recovery Milestone Gauge
- **ID:** `WARMUP-UI-04`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** When soft bounces occur, fixed schedules plow forward blindly, damaging sender reputation permanently.
- **Concrete UI Specification:**
  - **Layout:** Stepper progress tracker with milestone checkpoints (Day 1: 20, Day 5: 100, Day 14: 500, Day 30: 2,500).
  - **Visual Elements:**
    - Nodes transition to Amber warning status when soft bounce rate exceeds 2.0%.
    - Visual step-down trajectory animation showing 48-hour cool-down curve.
    - "Recovery Mode" badge with glowing orange border and SVG shield icon.
  - **Zero-CDN Invariant:** Pure CSS flex stepper with SVG icon defs.
- **CLI Counterpart Command:**
  - `gitmap warmup backoff status --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isBackoffActive", "hasRecoveryTarget", "currentStage"],
    "properties": {
      "domain": { "type": "string" },
      "isBackoffActive": { "type": "boolean" },
      "hasRecoveryTarget": { "type": "boolean" },
      "currentStage": { "type": "integer" },
      "coolDownHoursRemaining": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Auto-triggers UI state change within 500ms of soft bounce threshold breach.
  - Displays recovery milestone target and estimated days to regain standard curve.

---

### Improvement 05: Ramp Pause & Resume Circuit-Breaker Visual Switch
- **ID:** `WARMUP-UI-05`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Pausing warm-up currently requires killing CLI background processes or editing config files manually, causing orphaned locks.
- **Concrete UI Specification:**
  - **Layout:** Prominent glassmorphism header switch with safety guard latch.
  - **Visual Elements:**
    - 2-state mechanical toggle switch (`.circuit-toggle`) with LED status indicator.
    - Green LED (`#10b981`) = Active Sending; Pulsing Red LED (`#ef4444`) = Interrupted / Safe Hold.
    - Confirmation slide dialog with zero external modal libraries.
  - **Zero-CDN Invariant:** Native CSS checkbox styling with backdrop-filter glassmorphism.
- **CLI Counterpart Command:**
  - `gitmap warmup circuit pause --domain "acme.com" --reason "DNS propagation"`
  - `gitmap warmup circuit resume --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isTripped", "hasGracefulDrain", "trippedReason"],
    "properties": {
      "domain": { "type": "string" },
      "isTripped": { "type": "boolean" },
      "hasGracefulDrain": { "type": "boolean" },
      "trippedReason": { "type": "string" },
      "trippedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Toggling pause cleanly drains in-flight SMTP connections within 5 seconds.
  - Persists circuit state into `.gitmap/warmup.db` immediately.

---

### Improvement 06: Milestone Achievement Confetti & Notification Banner
- **ID:** `WARMUP-UI-06`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Operators have no immediate celebratory feedback when sender domains successfully clear critical warmup hurdles (e.g. 100/day without spam flag).
- **Concrete UI Specification:**
  - **Layout:** Slide-down top banner with lightweight inline SVG particle burst.
  - **Visual Elements:**
    - Dismissible toast banner (`.milestone-banner`) with emerald border and celebratory icon.
    - 24 inline CSS keyframe particle elements emitting from badge center.
    - Sound toggle switch with inline synthesized Web Audio API chime (zero MP3 files).
  - **Zero-CDN Invariant:** CSS `@keyframes` particle bursts and native Web Audio API.
- **CLI Counterpart Command:**
  - `gitmap warmup milestones list --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "milestoneId", "isCelebrated", "achievedAt"],
    "properties": {
      "domain": { "type": "string" },
      "milestoneId": { "type": "string" },
      "isCelebrated": { "type": "boolean" },
      "milestoneTitle": { "type": "string" },
      "achievedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Toast banner auto-dismisses after 8 seconds or upon clicking the `x` button.
  - Milestone record is marked `isCelebrated: true` in SQLite to prevent duplicate alerts.

---

### Improvement 07: Global Timezone Distribution Arc
- **ID:** `WARMUP-UI-07`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Sending to worldwide recipients without timezone awareness delivers warmup emails at 03:00 AM recipient local time, dropping open rates.
- **Concrete UI Specification:**
  - **Layout:** Circular polar projection / radial clock wheel displaying 24 timezones.
  - **Visual Elements:**
    - Concentric SVG rings showing active sending windows for Americas, EMEA, and APAC.
    - Color arcs showing sun/business hours (09:00 - 17:00 local time).
    - Current UTC rotating indicator beam.
  - **Zero-CDN Invariant:** Pure SVG polar coordinates calculated via trigonometric math in vanilla JS.
- **CLI Counterpart Command:**
  - `gitmap warmup timezone set --domain "acme.com" --zones "US/Eastern,Europe/London,Asia/Tokyo"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isLocalWindowEnforced", "targetZones"],
    "properties": {
      "domain": { "type": "string" },
      "isLocalWindowEnforced": { "type": "boolean" },
      "targetZones": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["timeZoneName", "isSendingAllowed", "utcOffsetHours"],
          "properties": {
            "timeZoneName": { "type": "string" },
            "isSendingAllowed": { "type": "boolean" },
            "utcOffsetHours": { "type": "integer" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Restricts automated seed/ramp dispatch to recipient business hours when enabled.
  - Arc updates in real time according to current system UTC.

---

### Improvement 08: Custom Ramp Formula Sandbox & Live Preview
- **ID:** `WARMUP-UI-08`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Experienced email engineers require custom progression formulas (e.g., custom polynomial or logarithmic progression) but have no safe sandbox to model output.
- **Concrete UI Specification:**
  - **Layout:** Split-view editor: formula syntax input on the left, live curve recalculation on the right.
  - **Visual Elements:**
    - Formula input field with monospace font and inline mathematical variable chips (`n`, `day`, `base`).
    - Instantaneous SVG curve repaint without page reload.
    - Error callout box highlighting invalid expressions with early validation.
  - **Zero-CDN Invariant:** In-browser AST math parser (zero eval, zero external dependencies).
- **CLI Counterpart Command:**
  - `gitmap warmup formula test --expr "base * pow(1.15, day)" --days 30`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["expression", "isValidFormula", "projectedTotals"],
    "properties": {
      "expression": { "type": "string" },
      "isValidFormula": { "type": "boolean" },
      "projectedTotals": {
        "type": "array",
        "items": { "type": "integer" }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Formula validation executes within `< 2ms` with complete syntax safety.
  - Invalid formulas display inline error message without disrupting active schedule.

---

### Improvement 09: Schedule Overlap & Campaign Clash Detector
- **ID:** `WARMUP-UI-09`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Operators schedule cold outreach or newsletter blasts simultaneously with warmup ramps, exceeding ESP domain throughput limits.
- **Concrete UI Specification:**
  - **Layout:** Dual stacked Gantt bar chart (`.schedule-clash-card`).
  - **Visual Elements:**
    - Top bar: Scheduled warmup traffic volume.
    - Bottom bar: External marketing/transactional broadcast traffic volume.
    - Hatch-patterned diagonal red stripes (`#ef4444`) marking time intervals where combined load exceeds daily sender rate limits.
  - **Zero-CDN Invariant:** Native HTML5 `<svg>` pattern defs for striping.
- **CLI Counterpart Command:**
  - `gitmap warmup check-clash --domain "acme.com" --campaign-id "launch-v2"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "hasClash", "maxCombinedRatePerHour", "clashWindows"],
    "properties": {
      "domain": { "type": "string" },
      "hasClash": { "type": "boolean" },
      "maxCombinedRatePerHour": { "type": "integer" },
      "clashWindows": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["startHour", "endHour", "excessEmailsPerHour"],
          "properties": {
            "startHour": { "type": "string" },
            "endHour": { "type": "string" },
            "excessEmailsPerHour": { "type": "integer" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Flashes visual warning icon when combined hourly volume exceeds 80% of ISP ceiling.
  - Provides one-click "De-conflict Schedule" button that auto-adjusts warmup slots.

---

### Improvement 10: Historical Ramp Performance Replay & Scrubber
- **ID:** `WARMUP-UI-10`
- **Cluster:** 1 (Warm-Up Ramp Curves & Visual Schedule Timeline)
- **Problem Statement:** Diagnosing past deliverability drops requires manually joining historical logs across days, lacking a scrubbable temporal playback interface.
- **Concrete UI Specification:**
  - **Layout:** Media player style slider bar with Play/Pause, Rewind, and Step buttons.
  - **Visual Elements:**
    - Temporal range slider (`<input type="range" class="timeline-scrubber">`).
    - Dynamic re-drawing of gauge values and curve status as slider scrubs backwards through past days.
    - Marker flags denoting major DNS or IP configuration changes.
  - **Zero-CDN Invariant:** Pure vanilla JS requestAnimationFrame loop and HTML range slider.
- **CLI Counterpart Command:**
  - `gitmap warmup history replay --domain "acme.com" --from "2026-09-01" --to "2026-10-01"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "hasSnapshots", "totalFrames", "frames"],
    "properties": {
      "domain": { "type": "string" },
      "hasSnapshots": { "type": "boolean" },
      "totalFrames": { "type": "integer" },
      "frames": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["timestamp", "dayIndex", "sentVolume", "spamRate"],
          "properties": {
            "timestamp": { "type": "string", "format": "date-time" },
            "dayIndex": { "type": "integer" },
            "sentVolume": { "type": "integer" },
            "spamRate": { "type": "number" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Smooth 60fps scrubber dragging across 90 days of historical data.
  - Instantly exposes exact day and hour an ISP penalty was incurred.

---

# Cluster 2: Deliverability Health Gauges & ISP Inbox/Spam Placement (Items 11–20)

---

### Improvement 11: Tri-Color Radial Deliverability Score Gauge
- **ID:** `WARMUP-UI-11`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Operators glance at multiple scattered metrics (bounces, complaints, open rates) without a single consolidated authoritative health score.
- **Concrete UI Specification:**
  - **Layout:** Circular SVG dial (180-degree semi-circle gauge) centered in primary health card.
  - **Visual Elements:**
    - Background track: `stroke: #334155; stroke-width: 12;`.
    - Dynamic progress arc: `stroke-dasharray` and `stroke-dashoffset` animated via CSS transition.
    - Three status zones: Red (0-59), Amber (60-79), Emerald Green (80-100).
    - Bold center typography displaying score (e.g., `94/100`) and subtext `EXCELLENT REPUTATION`.
  - **Zero-CDN Invariant:** Mathematical SVG arc formulas with inline CSS variables.
- **CLI Counterpart Command:**
  - `gitmap warmup health score --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "healthScore", "isHealthy", "statusTier"],
    "properties": {
      "domain": { "type": "string" },
      "healthScore": { "type": "number", "minimum": 0, "maximum": 100 },
      "isHealthy": { "type": "boolean" },
      "statusTier": { "type": "string", "enum": ["poor", "warning", "optimal", "excellent"] }
    }
  }
  ```
- **Acceptance Criteria:**
  - Arc smoothly animates to target score on page load within 600ms.
  - Automatically turns Amber if score drops below 80 and Red below 60.

---

### Improvement 12: Major ISP Placement Breakdown Bar Stack
- **ID:** `WARMUP-UI-12`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Overall deliverability hides ISP-specific penalties (e.g. 100% inbox on Google Workspace but 0% on Microsoft Outlook/O365).
- **Concrete UI Specification:**
  - **Layout:** Stacked horizontal bar table segmented by primary providers: Google Workspace, Microsoft 365, Yahoo/AOL, Apple iCloud, and Private MX.
  - **Visual Elements:**
    - Tri-segment bar per ISP: Green = Primary Inbox %, Amber = Promotions/Tabs %, Red = Spam/Junk %.
    - Dedicated ISP inline monochrome SVG logo alongside provider name.
    - Percentage label rendered inside bar segments when width > 12%.
  - **Zero-CDN Invariant:** Flexbox bar layout using inline SVGs for provider emblems.
- **CLI Counterpart Command:**
  - `gitmap warmup isp breakdown --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isMultiIspEvaluated", "providers"],
    "properties": {
      "domain": { "type": "string" },
      "isMultiIspEvaluated": { "type": "boolean" },
      "providers": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["providerCode", "inboxPercentage", "spamPercentage", "isDegraded"],
          "properties": {
            "providerCode": { "type": "string", "enum": ["google", "microsoft", "yahoo", "apple", "other"] },
            "inboxPercentage": { "type": "number", "minimum": 0, "maximum": 100 },
            "spamPercentage": { "type": "number", "minimum": 0, "maximum": 100 },
            "isDegraded": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Sum of inbox, tab, and spam percentages per provider equals 100.0%.
  - Highlights degraded provider row with subtle amber border.

---

### Improvement 13: DNS Authentication Matrix (SPF / DKIM / DMARC / BIMI)
- **ID:** `WARMUP-UI-13`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Missing or misconfigured DNS security records quietly destroy deliverability without explicit error feedback.
- **Concrete UI Specification:**
  - **Layout:** 4-column card grid displaying real-time DNS authentication validation status.
  - **Visual Elements:**
    - Cards for SPF, DKIM (Selector-aware), DMARC (Enforcement policy badge: `none`, `quarantine`, `reject`), and BIMI (SVG icon status).
    - Status pill badges: Emerald checkmark (`hasValidRecord: true`) or Red warning icon with inline copyable fix snippet.
    - "Re-Verify Now" icon button with spinning loading animation.
  - **Zero-CDN Invariant:** Self-contained inline SVG status icons.
- **CLI Counterpart Command:**
  - `gitmap warmup dns check --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "hasSpf", "hasDkim", "hasDmarc", "hasBimi", "isCompliant"],
    "properties": {
      "domain": { "type": "string" },
      "hasSpf": { "type": "boolean" },
      "hasDkim": { "type": "boolean" },
      "hasDmarc": { "type": "boolean" },
      "hasBimi": { "type": "boolean" },
      "isCompliant": { "type": "boolean" },
      "dmarcPolicy": { "type": "string", "enum": ["none", "quarantine", "reject", "missing"] }
    }
  }
  ```
- **Acceptance Criteria:**
  - Queries active DNS records via local Go DNS resolver without blocking GUI.
  - Displays one-click "Copy Recommended TXT Record" to clipboard.

---

### Improvement 14: Real-Time Seed Inbox Placement Feed
- **ID:** `WARMUP-UI-14`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Operators lack a live log of synthetic seed test email arrivals, forcing them to guess whether tests are currently landing.
- **Concrete UI Specification:**
  - **Layout:** Streaming live activity log table with auto-scroll lock toggle.
  - **Visual Elements:**
    - Rows flash subtle blue highlight upon new IMAP seed detection.
    - Columns: Timestamp, Seed Address (masked: `j***@gmail.com`), Destination Folder (`INBOX`, `SPAM`, `PROMOTIONS`), Delivery Latency (ms).
    - Folder badges: Green for `INBOX`, Red for `SPAM`, Yellow for `CATEGORY_PROMO`.
  - **Zero-CDN Invariant:** Lightweight JavaScript DOM table append with capped 100-row buffer.
- **CLI Counterpart Command:**
  - `gitmap warmup feed tail --domain "acme.com" --limit 20`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isStreaming", "events"],
    "properties": {
      "domain": { "type": "string" },
      "isStreaming": { "type": "boolean" },
      "events": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["eventId", "timestamp", "isInboxPlaced", "folderName", "latencyMs"],
          "properties": {
            "eventId": { "type": "string" },
            "timestamp": { "type": "string", "format": "date-time" },
            "isInboxPlaced": { "type": "boolean" },
            "folderName": { "type": "string" },
            "latencyMs": { "type": "integer" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Activity feed limits memory usage by retaining at most 100 items in DOM.
  - Displays instant delivery latency in milliseconds.

---

### Improvement 15: Blacklist & Blocklist Radar Scanner
- **ID:** `WARMUP-UI-15`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** IP or domain listing on DNSBLs (Spamhaus, Barracuda, SORBS) causes instant silent delivery dropouts.
- **Concrete UI Specification:**
  - **Layout:** Radar scanner widget listing 50+ DNSBL authorities with health status.
  - **Visual Elements:**
    - Circular rotating radar sweep line (`@keyframes radarSweep`) over status grid.
    - High-density list with green dot for clean and flashing red alert tag for listed.
    - One-click "Delisting Request" link helper for any positive blacklist match.
  - **Zero-CDN Invariant:** CSS transform rotation radar animation.
- **CLI Counterpart Command:**
  - `gitmap warmup blacklist scan --ip "198.51.100.42" --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["target", "isClean", "hasActiveListings", "checkedAuthoritiesCount"],
    "properties": {
      "target": { "type": "string" },
      "isClean": { "type": "boolean" },
      "hasActiveListings": { "type": "boolean" },
      "checkedAuthoritiesCount": { "type": "integer" },
      "listings": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["authorityName", "isListed", "delistUrl"],
          "properties": {
            "authorityName": { "type": "string" },
            "isListed": { "type": "boolean" },
            "delistUrl": { "type": "string" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Asynchronously queries 50+ RBLs in parallel via worker pool in `< 3.5s`.
  - Highlights any listing immediately with prominent red warning badge.

---

### Improvement 16: Spam Filter Word & Content Heuristic Linter
- **ID:** `WARMUP-UI-16`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Synthetic warmup email copy containing spam trigger words (e.g., "100% FREE", "ACT NOW", excessive exclamation points) poisons seed reputation.
- **Concrete UI Specification:**
  - **Layout:** Interactive copy inspector panel with subject line and body input tabs.
  - **Visual Elements:**
    - Inline yellow wavy underlines on detected trigger phrases.
    - Real-time SpamAssassin / heuristic score counter (e.g. `Score: 1.2 / 5.0 (SAFE)`).
    - HTML-to-text ratio meter and uppercase character percentage bar.
  - **Zero-CDN Invariant:** Local regex dictionary matching (zero remote API dependencies).
- **CLI Counterpart Command:**
  - `gitmap warmup lint-content --file "./templates/warmup-01.eml"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["isSpamClean", "heuristicScore", "hasExcessiveCaps", "flaggedTokens"],
    "properties": {
      "isSpamClean": { "type": "boolean" },
      "heuristicScore": { "type": "number" },
      "hasExcessiveCaps": { "type": "boolean" },
      "flaggedTokens": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["token", "severity", "category"],
          "properties": {
            "token": { "type": "string" },
            "severity": { "type": "string", "enum": ["low", "medium", "critical"] },
            "category": { "type": "string" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Lints 2,000-word template in `< 10ms`.
  - Warns if HTML-to-text ratio exceeds 40% code overhead.

---

### Improvement 17: Bounce Classification & Diagnostic Breakdown
- **ID:** `WARMUP-UI-17`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Bounces are often lumped together without distinguishing between temporary greylisting, full mailboxes, and permanent domain blocks.
- **Concrete UI Specification:**
  - **Layout:** Donut chart breakdown categorizing bounces: Hard (5.1.1), Soft (4.2.2 Mailbox Full), Greylist (4.5.1), and Policy Rejections (5.7.1).
  - **Visual Elements:**
    - Segmented SVG ring donut with hover segment expansion.
    - Drill-down drawer showing raw SMTP server response transcripts (`550 5.7.1 Service unavailable`).
    - Recommendation pills ("Retry in 4 hours", "Prune address immediately").
  - **Zero-CDN Invariant:** Pure SVG donut ring with native hover states.
- **CLI Counterpart Command:**
  - `gitmap warmup bounces list --domain "acme.com" --type "hard"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "hasBounces", "totalBounces", "hardBounceCount", "softBounceCount"],
    "properties": {
      "domain": { "type": "string" },
      "hasBounces": { "type": "boolean" },
      "totalBounces": { "type": "integer" },
      "hardBounceCount": { "type": "integer" },
      "softBounceCount": { "type": "integer" },
      "bounceRate": { "type": "number" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately parses standard RFC 3463 status codes.
  - Hard bounce addresses are flagged for automatic seed pool removal.

---

### Improvement 18: Feedback Loop (FBL) Complaint Tracker
- **ID:** `WARMUP-UI-18`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Spam complaint rates above 0.1% (1 in 1,000) permanently destroy Google and Yahoo inbox placement.
- **Concrete UI Specification:**
  - **Layout:** High-visibility safety threshold gauge with critical warning redline at `0.10%`.
  - **Visual Elements:**
    - Needle pointer indicating current 30-day rolling complaint rate (`0.02%`).
    - Flashing alert bar when rate exceeds `0.08%` (pre-penalty buffer zone).
    - Auto-remediation link to immediately suppress complaint origins.
  - **Zero-CDN Invariant:** Inline SVG dial with mathematical needle rotation.
- **CLI Counterpart Command:**
  - `gitmap warmup fbl stats --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "complaintRate", "isCompliant", "hasExceededThreshold"],
    "properties": {
      "domain": { "type": "string" },
      "complaintRate": { "type": "number", "minimum": 0 },
      "isCompliant": { "type": "boolean" },
      "hasExceededThreshold": { "type": "boolean" },
      "rollingWindowDays": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - If complaint rate exceeds 0.08%, automatically prompts operator to throttle ramp speed by 50%.
  - Supports ARF (Abuse Reporting Format) parsing.

---

### Improvement 19: Mailbox Storage & Quota Capacity Meter
- **ID:** `WARMUP-UI-19`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Warmup sender mailboxes fill up with synthetic replies and sent items, failing silently due to exceeded disk storage quotas.
- **Concrete UI Specification:**
  - **Layout:** Storage card with horizontal multi-color capacity bar per connected sender account.
  - **Visual Elements:**
    - Progress bar: Blue (0-70%), Amber (71-89%), Red (90-100%).
    - Storage breakdown: Messages vs. Attachments vs. System Logs.
    - One-click "Safe Auto-Prune Seed Threads (> 30 days)" action button.
  - **Zero-CDN Invariant:** CSS flex bar with percentage width styling.
- **CLI Counterpart Command:**
  - `gitmap warmup mailbox quota --account "sales@acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["accountEmail", "usedBytes", "totalBytes", "isStorageCritical", "hasAutoPrune"],
    "properties": {
      "accountEmail": { "type": "string" },
      "usedBytes": { "type": "integer" },
      "totalBytes": { "type": "integer" },
      "isStorageCritical": { "type": "boolean" },
      "hasAutoPrune": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Warns operator when storage utilization exceeds 85%.
  - Auto-prune respects seed preservation rules, deleting only aged synthetic threads.

---

### Improvement 20: Comprehensive Domain Deliverability Export Card
- **ID:** `WARMUP-UI-20`
- **Cluster:** 2 (Deliverability Health Gauges & ISP Inbox/Spam Placement)
- **Problem Statement:** Sharing deliverability health with stakeholders requires screenshotting UI views instead of generating clean, zero-dependency PDF/SVG reports.
- **Concrete UI Specification:**
  - **Layout:** Export modal card offering standalone self-contained HTML, SVG vector bundle, or Markdown audit format.
  - **Visual Elements:**
    - Checkbox selector for report sections (Health score, ISP breakdown, DNS status, Timeline).
    - "Generate Standalone Report" button with progress spinner.
    - Output report contains embedded base64 SVGs and inline CSS (100% offline viewable).
  - **Zero-CDN Invariant:** Generates static self-contained HTML without external CSS or scripts.
- **CLI Counterpart Command:**
  - `gitmap warmup report export --domain "acme.com" --output "./reports/health.html"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isGenerated", "reportPath", "fileSizeBytes"],
    "properties": {
      "domain": { "type": "string" },
      "isGenerated": { "type": "boolean" },
      "reportPath": { "type": "string" },
      "fileSizeBytes": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Exported HTML opens cleanly in offline browser with zero console errors or network calls.
  - File generation takes `< 1.2s`.

---

# Cluster 3: Seed List Scheduler, Reply Simulation & Automated Interaction Physics (Items 21–30)

---

### Improvement 21: Seed Pool Diversity Distribution Radar
- **ID:** `WARMUP-UI-21`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Warmup networks with poorly distributed seed pools (e.g. 95% Gmail, 5% others) fail to train Microsoft and Yahoo algorithms, biasing reputation.
- **Concrete UI Specification:**
  - **Layout:** Pentagonal SVG radar/spider web chart displaying pool distribution across 5 provider pillars.
  - **Visual Elements:**
    - Background grid rings representing 20%, 40%, 60%, 80%, 100% quota balance.
    - Shaded polygon area (`fill: rgba(59, 130, 246, 0.3); stroke: #3b82f6; stroke-width: 2;`).
    - Ideal balance reference polygon outline (`stroke: #10b981; stroke-dasharray: 4 2;`).
  - **Zero-CDN Invariant:** Native SVG polygon rendering calculated with polar coordinate transforms.
- **CLI Counterpart Command:**
  - `gitmap warmup seeds diversity --pool-id "default"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["poolId", "isBalanced", "totalSeeds", "providerDistribution"],
    "properties": {
      "poolId": { "type": "string" },
      "isBalanced": { "type": "boolean" },
      "totalSeeds": { "type": "integer" },
      "providerDistribution": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["provider", "seedCount", "percentage", "isSufficient"],
          "properties": {
            "provider": { "type": "string" },
            "seedCount": { "type": "integer" },
            "percentage": { "type": "number" },
            "isSufficient": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Flags warning badge if any major provider comprises less than 15% of total active seed pool.
  - Updates instantaneously when new seeds are enqueued.

---

### Improvement 22: Conversational Multi-Turn Reply Thread Visualizer
- **ID:** `WARMUP-UI-22`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** ISPs detect one-way automated emails. Warmup requires multi-turn back-and-forth conversation trees, but operators cannot visualize conversation depth.
- **Concrete UI Specification:**
  - **Layout:** Interactive threaded conversation tree view with nested branch lines.
  - **Visual Elements:**
    - Message nodes displaying sender icon, timestamp, sentiment score, and turn index (`Turn 1/4`, `Turn 2/4`).
    - Connected vertical SVG tree-branch guides (`constants.TreeBranch`, `constants.TreeCorner`).
    - Expandable node body previewing context-aware AI synthetic replies.
  - **Zero-CDN Invariant:** CSS flexbox hierarchy with SVG connector line segments.
- **CLI Counterpart Command:**
  - `gitmap warmup threads show --thread-id "thr-84920"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["threadId", "isCompleted", "turnCount", "messages"],
    "properties": {
      "threadId": { "type": "string" },
      "isCompleted": { "type": "boolean" },
      "turnCount": { "type": "integer", "minimum": 1 },
      "messages": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["messageId", "senderEmail", "recipientEmail", "isSeedSender", "turnIndex"],
          "properties": {
            "messageId": { "type": "string" },
            "senderEmail": { "type": "string" },
            "recipientEmail": { "type": "string" },
            "isSeedSender": { "type": "boolean" },
            "turnIndex": { "type": "integer" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately renders deep conversation threads up to 6 turns without layout overlap.
  - Highlights active turn waiting for scheduled synthetic reply.

---

### Improvement 23: Human Interaction Velocity & Read-Time Physics Simulator
- **ID:** `WARMUP-UI-23`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Opening an email 0.05 seconds after delivery and replying instantly flags synthetic bot behavior at Gmail/Outlook.
- **Concrete UI Specification:**
  - **Layout:** Interactive physics parameter control card (`.interaction-physics-card`).
  - **Visual Elements:**
    - Dual range sliders for "Read Duration Delay" (30s to 300s) and "Reply Delay Window" (5m to 120m).
    - Gaussian bell-curve graph demonstrating simulated human behavioral latency distribution.
    - Live countdown badge on pending replies showing simulated "reading" status.
  - **Zero-CDN Invariant:** Mathematical SVG normal distribution curve updated via range slider input.
- **CLI Counterpart Command:**
  - `gitmap warmup physics set --domain "acme.com" --min-read 45s --max-read 240s`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isGaussianJitter", "minReadSeconds", "maxReadSeconds", "minReplyMinutes"],
    "properties": {
      "domain": { "type": "string" },
      "isGaussianJitter": { "type": "boolean" },
      "minReadSeconds": { "type": "integer" },
      "maxReadSeconds": { "type": "integer" },
      "minReplyMinutes": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Injected delay adheres to Box-Muller transform Gaussian distribution.
  - Minimum simulated read duration cannot be set below 15 seconds.

---

### Improvement 24: "Mark as Important" & "Unspam" Remediation Trigger Monitor
- **ID:** `WARMUP-UI-24`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** The strongest positive deliverability signal is a seed user moving an email from Spam to Inbox and marking it Starred/Important, but operators cannot monitor this remediation.
- **Concrete UI Specification:**
  - **Layout:** Action counter widget displaying positive signal multipliers: Starred, Marked Important, Moved from Spam to Inbox, and Added to Contacts.
  - **Visual Elements:**
    - Gold star SVG badge (`#f59e0b`), Yellow tag badge for Important, Green arrow for Rescue from Spam.
    - Animated counter incrementing live when a seed mailbox IMAP agent executes remediation.
    - Historical trend line showing weekly rescue actions.
  - **Zero-CDN Invariant:** Native CSS keyframe bounce animation on badge increment.
- **CLI Counterpart Command:**
  - `gitmap warmup rescue stats --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isAutoRescueEnabled", "rescuedFromSpamCount", "markedImportantCount"],
    "properties": {
      "domain": { "type": "string" },
      "isAutoRescueEnabled": { "type": "boolean" },
      "rescuedFromSpamCount": { "type": "integer" },
      "markedImportantCount": { "type": "integer" },
      "addedToContactsCount": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - When email lands in Spam on seed box, auto-rescue triggers within 2-5 minutes.
  - Visual tally updates immediately via websocket/event log.

---

### Improvement 25: Seed Mailbox Health & Authentication Credential Vault
- **ID:** `WARMUP-UI-25`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Seed mailboxes experience password expiry, OAuth revocation, or IMAP connection lockouts without immediate operator visibility.
- **Concrete UI Specification:**
  - **Layout:** High-density seed account grid showing individual mailbox health status.
  - **Visual Elements:**
    - Green status dot for authenticated active, red for credential failure, amber for rate-limited.
    - Protocol badges: `OAUTH2`, `APP-PASS`, `IMAP/SMTP`.
    - One-click "Test Connection" button per account with inline latency result.
  - **Zero-CDN Invariant:** Clean CSS table grid with inline SVG status glyphs.
- **CLI Counterpart Command:**
  - `gitmap warmup seeds test-all --pool-id "primary"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["poolId", "hasFailingSeeds", "totalActive", "seeds"],
    "properties": {
      "poolId": { "type": "string" },
      "hasFailingSeeds": { "type": "boolean" },
      "totalActive": { "type": "integer" },
      "seeds": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["seedEmail", "isHealthy", "hasImapAccess", "lastConnectedAt"],
          "properties": {
            "seedEmail": { "type": "string" },
            "isHealthy": { "type": "boolean" },
            "hasImapAccess": { "type": "boolean" },
            "lastConnectedAt": { "type": "string", "format": "date-time" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Dead seed mailboxes are automatically quarantined from active schedules within 60 seconds.
  - Masked email addresses protect privacy in shared UI viewports.

---

### Improvement 26: Reply Ratio Dynamics & Target Balancing Slider
- **ID:** `WARMUP-UI-26`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Rigid 100% reply rates look fake to machine learning spam classifiers, while 0% reply rates signal broadcast spam.
- **Concrete UI Specification:**
  - **Layout:** Dynamic target ratio slider card (`.reply-ratio-card`).
  - **Visual Elements:**
    - Dual handle slider setting target reply percentage band (e.g. 25% - 40%).
    - Real-time gauge comparing Target Ratio vs. Actual Realized Ratio over past 7 days.
    - Visual indicator of natural human benchmarks (B2B baseline: 18-32%).
  - **Zero-CDN Invariant:** Native range inputs with gradient accent track.
- **CLI Counterpart Command:**
  - `gitmap warmup reply-ratio set --domain "acme.com" --min 20 --max 35`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isTargetMet", "targetMinPercent", "targetMaxPercent", "actualPercent"],
    "properties": {
      "domain": { "type": "string" },
      "isTargetMet": { "type": "boolean" },
      "targetMinPercent": { "type": "number" },
      "targetMaxPercent": { "type": "number" },
      "actualPercent": { "type": "number" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Algorithm randomly samples target reply decisions within configured bounds.
  - UI highlights green when actual ratio falls within configured target band.

---

### Improvement 27: AI Subject Line & Content Variation Cluster Inspector
- **ID:** `WARMUP-UI-27`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Sending identical subject lines across hundreds of warmup emails triggers duplicate content fingerprinting at Gmail.
- **Concrete UI Specification:**
  - **Layout:** Topic cluster bubble visualizer and subject line entropy inspector.
  - **Visual Elements:**
    - SVG cluster circles representing topic domains (Tech, Logistics, Scheduling, Inquiries).
    - Content entropy score bar (`Entropy: 8.4/10 - HIGH DIVERSITY`).
    - Collapsible drawer listing recent subject line variations and sentiment polarity.
  - **Zero-CDN Invariant:** Pure SVG circle layout with hover scale transform.
- **CLI Counterpart Command:**
  - `gitmap warmup content entropy --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isDiverse", "entropyScore", "topicClusters"],
    "properties": {
      "domain": { "type": "string" },
      "isDiverse": { "type": "boolean" },
      "entropyScore": { "type": "number", "minimum": 0, "maximum": 10 },
      "topicClusters": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["topicName", "sampleCount", "hasAdequateDiversity"],
          "properties": {
            "topicName": { "type": "string" },
            "sampleCount": { "type": "integer" },
            "hasAdequateDiversity": { "type": "boolean" }
          }
        }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Flags warning when Levenshtein similarity between any two warmup subjects exceeds 70%.
  - Maintains minimum 5 distinct conversation topic clusters.

---

### Improvement 28: Seed Mailbox Tagging & Category Nudge Matrix
- **ID:** `WARMUP-UI-28`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Landing in the "Promotions" or "Updates" tab at Gmail prevents high open rates. Seeds must actively nudge messages into the "Primary" inbox.
- **Concrete UI Specification:**
  - **Layout:** Matrix table tracking Gmail tab classification and drag-to-primary remediation actions.
  - **Visual Elements:**
    - Tabs breakdown: Primary Inbox %, Promotions %, Updates %, Spam %.
    - "Drag-to-Primary" action log displaying automated IMAP label reclassification.
    - Trend sparkline showing migration from Promotions tab to Primary tab over 14 days.
  - **Zero-CDN Invariant:** SVG sparkline polyline with inline points.
- **CLI Counterpart Command:**
  - `gitmap warmup tabs status --domain "acme.com"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isPrimaryDominant", "primaryPercent", "promotionsPercent"],
    "properties": {
      "domain": { "type": "string" },
      "isPrimaryDominant": { "type": "boolean" },
      "primaryPercent": { "type": "number" },
      "promotionsPercent": { "type": "number" },
      "updatesPercent": { "type": "number" },
      "reclassificationCount": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Seed IMAP worker removes `\Promotions` label and assigns `\Inbox` upon detection.
  - Visual matrix updates within 3 minutes of label modification.

---

### Improvement 29: Seed Peer-to-Peer Mesh Connection Graph
- **ID:** `WARMUP-UI-29`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Circular email rings (A -> B -> C -> A) create detectable graph motifs that spam engines easily identify and penalize.
- **Concrete UI Specification:**
  - **Layout:** Force-directed SVG topology node mesh displaying seed-to-seed exchange vectors.
  - **Visual Elements:**
    - Circular nodes representing sender domain and seed mailboxes.
    - Directed arrows (`<marker id="arrow">`) indicating communication paths.
    - Graph motif detection indicator: Green for Acyclic / High Entropy Mesh, Red for Detected Closed Loops.
  - **Zero-CDN Invariant:** Native SVG node graph with CSS transform animations.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh verify --pool-id "primary"`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["poolId", "isAcyclic", "hasClosedLoops", "nodeCount", "edgeCount"],
    "properties": {
      "poolId": { "type": "string" },
      "isAcyclic": { "type": "boolean" },
      "hasClosedLoops": { "type": "boolean" },
      "nodeCount": { "type": "integer" },
      "edgeCount": { "type": "integer" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Validates that no 2-node or 3-node cyclic loops exist in the exchange schedule.
  - Graph renders cleanly with up to 100 concurrent nodes in `< 25ms`.

---

### Improvement 30: Warm-Up Completion & Production Handoff Wizard
- **ID:** `WARMUP-UI-30`
- **Cluster:** 3 (Seed List Scheduler, Reply Simulation & Automated Interaction Physics)
- **Problem Statement:** Once warmup is complete, operators struggle to safely transition domains into production sending without abrupt volume spikes.
- **Concrete UI Specification:**
  - **Layout:** Multi-step production handoff wizard card (`.handoff-wizard-card`).
  - **Visual Elements:**
    - Readiness checklist with 6 green verification badges (DNS, Health > 90, Seed Pool stable, 0% Blacklists, 30-day stability).
    - Production volume transition slider blending warmup baseline with real marketing traffic.
    - "Graduate Domain to Production" primary action button with safety unlock toggle.
  - **Zero-CDN Invariant:** Vanilla JS wizard step controller with CSS transition fades.
- **CLI Counterpart Command:**
  - `gitmap warmup graduate --domain "acme.com" --baseline-pulse 15`
- **Data Contract:**
  ```json
  {
    "type": "object",
    "required": ["domain", "isGraduated", "isProductionReady", "readinessScore", "baselinePulsePercent"],
    "properties": {
      "domain": { "type": "string" },
      "isGraduated": { "type": "boolean" },
      "isProductionReady": { "type": "boolean" },
      "readinessScore": { "type": "number", "minimum": 0, "maximum": 100 },
      "baselinePulsePercent": { "type": "integer", "default": 15 },
      "graduatedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Prevents graduation if deliverability health score is below 85.
  - Configures background baseline pulse of 10-15% volume to preserve IP warmth permanently.
