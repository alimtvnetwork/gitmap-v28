# Specification 339: Warm-Up UI and Emailing UI Improvements (Overview & Architecture)

- **Sequence:** 339
- **Slug:** `339-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## 1. Executive Summary & Core Innovations

Sequence 339 establishes the modern, declarative Web Component foundation and intelligent historical analysis suite for the GitMap Warm-Up & Outbound Emailing platform. Building upon the node topology canvas, AI co-pilot, and persona networks of Sequence 338, Sequence 339 focuses on modular zero-CDN custom elements, deterministic temporal state scrubbing, and adaptive machine learning warming curves:

1. **Zero-CDN Shadow DOM Custom Web Components** (Cluster 21, Items 201–210): Fully encapsulated, framework-agnostic native custom elements (`<warmup-metric-gauge>`, `<warmup-envelope-card>`, `<warmup-time-scrubber>`, etc.) with closed/open Shadow DOM boundaries, CSS variables theming, self-contained SVG micro-icons, and zero external runtime dependencies.
2. **Time-Travel Deliverability Debugger & Event Playback** (Cluster 22, Items 211–220): Deterministic temporal state engine allowing operators to scrub backwards and forwards through historical delivery events, inspect snapshot MTA states, identify exact timestamps of reputation degradation, isolate transient throttling spikes, and visually replay deliverability incidents.
3. **Adaptive Warming Gradient & Reputation Acceleration Engine** (Cluster 23, Items 221–230): Self-adjusting mathematical curve models (logistic, exponential sigmoid, multi-phase plateau) powered by real-time ISP feedback signals, inbox placement ratios, spam complaint rates, and automatic throttling deceleration.

```
                          ┌──────────────────────────────────────────────┐
                          │    GitMap Sequence 339 Architecture Matrix   │
                          └──────────────────────┬───────────────────────┘
                                                 │
                   ┌─────────────────────────────┼─────────────────────────────┐
                   ▼                             ▼                             ▼
       ┌───────────────────────┐     ┌───────────────────────┐     ┌───────────────────────┐
       │  Cluster 21 (201-210) │     │  Cluster 22 (211-220) │     │  Cluster 23 (221-230) │
       │  Zero-CDN Shadow DOM  │     │  Time-Travel          │     │  Adaptive Warming     │
       │  Custom Web Elements  │     │  Deliverability       │     │  Gradient Engine &    │
       │  & Card Design Tokens │     │  Temporal Debugger    │     │  Reputation Curves    │
       └───────────┬───────────┘     └───────────┬───────────┘     └───────────┬───────────┘
                   │                             │                             │
                   └─────────────────────────────┼─────────────────────────────┘
                                                 │
                                                 ▼
                         ┌──────────────────────────────────────────────┐
                         │      Core Architecture & Shared Invariants   │
                         │   - Split-DB SQLite Store (`warmup.db`)      │
                         │   - Error Return Sovereignty (Zero Exits)    │
                         │   - Pure HTML5/SVG/CSS3 (Zero External CDN)  │
                         │   - Coupled CLI Tools (`gitmap warmup ...`)  │
                         │   - Positive Booleans (`is*`, `has*` Only)   │
                         └──────────────────────────────────────────────┘
```

---

## 2. Antigravity Design Parity & UI Foundations

Sequence 339 adheres strictly to the Google Antigravity Manager UX framework and core repository invariants:

- **Zero External CDNs:** 100% vanilla Web Components using the native browser `customElements.define()` API, encapsulated Shadow DOM stylesheets, inline SVGs, and system font stacks (`-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", sans-serif`). Zero external CDNs, Google Fonts, unpkg bundles, or NPM client-side runtime dependencies.
- **Glassmorphic Obsidian Palette:** Deep dark surfaces (`#0b0f19`, `#111827`, `#1f2937`), ultra-subtle hairline borders (`rgba(255, 255, 255, 0.08)`), vibrant telemetry accents (`#10b981` emerald, `#06b6d4` cyan, `#f59e0b` amber, `#ef4444` rose, `#8b5cf6` violet), and GPU-accelerated backdrop filters (`backdrop-filter: blur(12px)`).
- **Encapsulated Styling with CSS Custom Properties:** Each custom element exposes standardized CSS variables (`--warmup-bg-card`, `--warmup-border-color`, `--warmup-accent-cyan`, `--warmup-glow-spread`) enabling seamless runtime theme switching and token isolation across host documents.
- **Micro-State Determinism & Strict Boolean Conventions:** All element attributes, DOM properties, and JSON data contracts enforce strictly positive boolean naming (`isActive`, `hasWarning`, `isScrubbing`, `hasThrottled`, `isAdaptive`). Negative flags (`isNotReady`, `disablePlayback`) are strictly forbidden.

---

## 3. High-Level System Architecture & Component Interaction

The warm-up telemetry stream flows from local SQLite storage through the CLI and Web Component view layer deterministically:

```
   [Browser UI / Embedded Host View] ◄──────────────┐
                 │                                   │ WebSocket Frame / Local Journal Feed
                 ▼                                   │
   ┌─────────────────────────────────────────────┐   │
   │ Custom Web Component Shell (<warmup-app>)   │   │
   │ ┌───────────────────┐ ┌───────────────────┐ │   │
   │ │ Shadow DOM Cards  │ │ Time-Travel Scrub │ │   │
   │ │ & Micro Gauges    │ │ Event Playback    │ │   │
   │ └─────────┬─────────┘ └─────────┬─────────┘ │   │
   │           │                     │           │   │
   │           ▼                     ▼           │   │
   │ ┌─────────────────────────────────────────┐ │   │
   │ │ Adaptive Machine Learning Curve Renderer│ │   │
   │ └─────────────────────┬───────────────────┘ │   │
   └───────────────────────┼─────────────────────┘   │
                           │                         │
                           ▼                         │
   [CLI Command Layer: `gitmap warmup`] ─────────────┘
                 │
                 ▼
   ┌─────────────────────────────────────────────────────────┐
   │ GitMap Warm-Up Daemon & Time-Travel Event Log Core      │
   │                                                         │
   │ ┌─────────────────────────┐ ┌─────────────────────────┐ │
   │ │ Event Snapshot Journal  │ │ Polynomial Curve Engine │ │
   │ │ - Delta Compression     │ │ - Sigmoid Warming Calc  │ │
   │ │ - Keyframe Checkpoints  │ │ - ISP Throttle Brake    │ │
   │ └───────────┬─────────────┘ └───────────┬─────────────┘ │
   │             │                           │               │
   │             ▼                           ▼               │
   │ ┌─────────────────────────────────────────────────────┐ │
   │ │ Split-DB SQLite Storage (`.gitmap/warmup.db`)       │ │
   │ │ - `WarmupTimelineEvents`, `StateKeyframes`          │ │
   │ │ - `AdaptiveCurvatures`, `ISPTelemetryDeltas`        │ │
   │ └─────────────────────────────────────────────────────┘ │
   └─────────────────────────────────────────────────────────┘
```

---

## 4. Design System Tokens & Invariants (Zero-CDN)

Sequence 339 elements share a unified, self-contained CSS token palette declared at `:root` or on the `<warmup-shell>` host:

```css
:host, :root {
  /* Surfaces & Backgrounds */
  --warmup-bg-base: #0b0f19;
  --warmup-bg-surface: #111827;
  --warmup-bg-card: rgba(17, 24, 39, 0.75);
  --warmup-bg-card-hover: rgba(31, 41, 55, 0.85);
  --warmup-bg-control: #1f2937;
  
  /* Hairline Borders & Glassmorphic Highlights */
  --warmup-border-subtle: rgba(255, 255, 255, 0.08);
  --warmup-border-focus: rgba(6, 182, 212, 0.45);
  --warmup-border-warning: rgba(245, 158, 11, 0.45);
  --warmup-border-danger: rgba(239, 68, 68, 0.45);
  
  /* Semantic Accent Ramps */
  --warmup-accent-cyan: #06b6d4;
  --warmup-accent-emerald: #10b981;
  --warmup-accent-amber: #f59e0b;
  --warmup-accent-rose: #ef4444;
  --warmup-accent-violet: #8b5cf6;
  
  /* Typography & Spatial Rhythm */
  --warmup-font-mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  --warmup-font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  --warmup-radius-sm: 4px;
  --warmup-radius-md: 8px;
  --warmup-radius-lg: 12px;
  --warmup-transition-fast: 150ms cubic-bezier(0.16, 1, 0.3, 1);
  --warmup-transition-normal: 250ms cubic-bezier(0.16, 1, 0.3, 1);
}
```

---

## 5. Storage & Database Schema Architecture (`warmup.db`)

Temporal event scrubbing and dynamic curvature calculation utilize two high-performance SQLite Split-DB tables defined in `.gitmap/warmup.db`:

```sql
-- Temporal Keyframe Snapshots for Time-Travel Debugging
CREATE TABLE IF NOT EXISTS WarmupStateSnapshots (
    snapshotId TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    epochTimestamp INTEGER NOT NULL,
    activeVolume INTEGER NOT NULL,
    reputationScore REAL NOT NULL,
    inboxPlacementRatio REAL NOT NULL,
    spamTrapHits INTEGER NOT NULL,
    throttlingEventsCount INTEGER NOT NULL,
    ispReputationDeltasJson TEXT NOT NULL,
    isActive BOOLEAN NOT NULL CHECK (isActive IN (0, 1)),
    hasAnomalies BOOLEAN NOT NULL CHECK (hasAnomalies IN (0, 1)),
    createdAt TEXT NOT NULL
);

-- Adaptive Warming Curvature Models and Coefficients
CREATE TABLE IF NOT EXISTS AdaptiveWarmingCurves (
    curveId TEXT PRIMARY KEY,
    domainName TEXT NOT NULL,
    modelType TEXT NOT NULL CHECK (modelType IN ('logistic', 'exponential_sigmoid', 'linear_stepped', 'adaptive_dynamic')),
    baseVolume INTEGER NOT NULL,
    targetMaxVolume INTEGER NOT NULL,
    currentStepIndex INTEGER NOT NULL,
    accelerationFactor REAL NOT NULL,
    decelerationThreshold REAL NOT NULL,
    curvatureParametersJson TEXT NOT NULL,
    isActive BOOLEAN NOT NULL CHECK (isActive IN (0, 1)),
    hasThrottleBrakeEngaged BOOLEAN NOT NULL CHECK (hasThrottleBrakeEngaged IN (0, 1)),
    updatedAt TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_snapshots_domain_epoch ON WarmupStateSnapshots (domainName, epochTimestamp DESC);
CREATE INDEX IF NOT EXISTS idx_curves_domain ON AdaptiveWarmingCurves (domainName, isActive);
```

---

## 6. Verification & Quality Gates

All specifications authored in Sequence 339 conform to the non-negotiable GitMap quality protocol:
1. **100% Relative Links:** All intra-document links strictly avoid absolute local machine paths.
2. **Error Return Sovereignty:** Code snippets and CLI counterparts adhere to leaf error returns (`return apperror.Wrap(...)`) without rogue calls to `os.Exit`.
3. **Strict Function & Complexity Bounds:** Helper modules adhere to <= 15 lines per function with early returns and zero nested if blocks.
4. **Autonomous Testing Isolation:** Unit tests simulate time-travel buffers and curve integration math in-memory without firing live external SMTP connections or triggering OS-level disruptions.
