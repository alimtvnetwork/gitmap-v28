# Subtask 01: Zero-CDN Web Components, Time-Travel Debugger & Adaptive Curves

## Status: READY
- **Subtask Slug:** `01-web-components-and-time-travel`
- **Parent Plan:** `.ai-memory/plans/pending/339-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Owner:** Spec Author Subagent 1 / Front-End Systems Engineer
- **Target Subsystem:** `pkg/web/components/`, `pkg/mail/timetravel/`, `pkg/mail/warmup/`
- **Exclusive Owned Files:**
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/02-improvements-201-to-230-web-components-and-time-travel.md`
  - `pkg/web/components/custom_elements.go`
  - `pkg/mail/timetravel/state_scrubber.go`
  - `pkg/mail/warmup/curve_modeler.go`
  - `web/static/js/components/warmup_elements.js`

---

## 1. Context & Objectives

This subtask delivers the core front-end Web Component framework and temporal debugging engine for Sequence 339:
1. **Zero-CDN Shadow DOM Elements (Cluster 21, Items 201–210):** Native encapsulated Web Components (`<warmup-metric-gauge>`, `<warmup-envelope-card>`, `<warmup-time-scrubber>`) with zero external CDN dependencies, self-contained SVG micro-icons, and CSS custom property theming.
2. **Time-Travel Deliverability Debugger (Cluster 22, Items 211–220):** Temporal state engine enabling bidirectional timeline scrubbing, MTA event snapshot inspection, reputation degradation root-cause isolation, and incident replay.
3. **Adaptive Warming Gradient Engine (Cluster 23, Items 221–230):** Non-linear mathematical warming curves (logistic, sigmoid) dynamically adjusted based on real-time ISP feedback and throttling telemetry.
4. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isActive`, `hasWarning`, `isScrubbing`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Shadow DOM Custom Element Definitions
- Implement `<warmup-metric-gauge>` with SVG circular progress visualization.
- Implement `<warmup-envelope-card>` with status badges, retry indicators, and SPF/DKIM verification tags.
- Implement `<warmup-time-scrubber>` with scrubbing handle, milestone markers, and playback controls.

### Step 2: Time-Travel Deliverability Debugger
- Architect `pkg/mail/timetravel/state_scrubber.go` to compute deterministic MTA snapshot deltas from SQLite event logs.
- Provide step-by-step playback API and binary state diffing algorithms.

### Step 3: Adaptive Warming Gradient Engine
- Implement mathematical curve calculation (logistic expansion, plateau stabilization, decay dampening) in `pkg/mail/warmup/curve_modeler.go`.
- Connect ISP bounce signals to dynamically throttle or accelerate daily sending limits.

---

## 3. Verification & Acceptance Criteria
- [ ] Custom Web Components register cleanly in vanilla DOM without external runtime libraries.
- [ ] Time-travel scrubber smoothly replays past 10,000 delivery events with sub-16ms frame rates.
- [ ] Adaptive warming curves conform to ISP-specific daily maximum sending bounds.
