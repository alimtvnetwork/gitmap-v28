# Subtask 01: Warm-Up Ramp Curves, Schedule Timelines & Deliverability Health Gauges

## Status: PLANNED
- **Subtask Slug:** `01-warmup-telemetry-and-ramp-curves`
- **Parent Plan:** `.ai-memory/plans/pending/337-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Owner:** Deliverability & Physics Subagent
- **Target Subsystem:** `pkg/mail/warmup/`, `pkg/mail/seed/`, `pkg/web/widgets/`
- **Improvements Covered:** 01 to 30 (Clusters 1, 2, and 3)
- **Exclusive Owned Files:**
  - `pkg/mail/warmup/ramp_curves.go`
  - `pkg/mail/warmup/schedule_timeline.go`
  - `pkg/mail/warmup/throttle_physics.go`
  - `pkg/mail/seed/seed_orchestrator.go`
  - `pkg/mail/seed/reply_simulator.go`
  - `pkg/web/widgets/deliverability_gauge.go`
  - `pkg/web/widgets/isp_breakdown.go`
  - `pkg/mail/warmup/types.go`

---

## 1. Executive Summary & Context

Subtask 01 delivers the foundational mathematical engine and visualization layer for email warm-up automation. Without disciplined warmup ramp curves and deliverability telemetry, newly provisioned sending domains face immediate spam filtering and reputation degradation.

This subtask implements:
1. **Mathematical Ramp Curve Profiles (Improvements 01–10):** Linear, Exponential, Sigmoid (S-Curve), Fibonacci, and Wavefront sending progressions with adaptive throttle physics and interactive milestone gates.
2. **Deliverability Health Gauges & ISP Breakdown (Improvements 11–20):** Tri-color radial SVG dial gauges, inbox vs spam placement monitors across major ISPs (Google Workspace, M365, Yahoo, iCloud), and circuit-breaker pause triggers.
3. **Seed List Scheduler & Reply Simulation (Improvements 21–30):** Coordinated seed mailboxes, synthetic reading velocity, multi-turn reply thread weaving, and automated sentiment-scored engagement.

---

## 2. File Modification & Structural Ownership

```text
pkg/mail/
├── warmup/
│   ├── ramp_curves.go          [NEW] - 5 mathematical progression curve calculators
│   ├── schedule_timeline.go    [NEW] - 30-day timeline scheduler with timezone dispersion
│   ├── throttle_physics.go     [NEW] - Jittered send interval and burst controllers
│   └── types.go                [NEW] - Domain structs, enums, and Result wrappers
├── seed/
│   ├── seed_orchestrator.go    [NEW] - Seed mailbox rotation and peer grouping
│   └── reply_simulator.go      [NEW] - Multi-turn conversational reply thread engine
pkg/web/widgets/
├── deliverability_gauge.go     [NEW] - Tri-color SVG radial score dial
└── isp_breakdown.go            [NEW] - Multi-ISP placement comparison cards
```

---

## 3. Technical Requirements & Architecture

### 3.1 Mathematical Progression Curves
- Implements `CalculateDailyQuota(day int, totalDays int, targetDailyVolume int, curveType RampCurveType) int`.
- Supports 5 distinct progression types:
  - `RampCurveLinear`: Constant daily increment.
  - `RampCurveExponential`: Conservative initial base doubling until midpoint.
  - `RampCurveSigmoid`: S-Curve providing gradual start, steep middle, and smooth plateau.
  - `RampCurveFibonacci`: Stepped milestone increases mapped to Fibonacci intervals.
  - `RampCurveWavefront`: Pulsing burst sending with resting recovery days.

### 3.2 Positive Boolean State Model
All models enforce strictly positive booleans:
- `isActive` (bool): Sending engine currently operational.
- `isPaused` (bool): Warmup temporarily paused due to threshold warning.
- `hasWarning` (bool): Bounce or complaint rate elevated.
- `isHealthy` (bool): Deliverability score >= 90.0%.
- `hasReplied` (bool): Simulated seed recipient replied to message.

### 3.3 Deliverability Health Gauges
- Zero external charting libraries; renders pure inline SVG radial dials.
- SVG path calculations for stroke-dasharray arc progression (0 to 100%).
- Real-time ISP segmentation: Google Workspace, Microsoft 365, Yahoo/AOL, Apple iCloud, and Private Corporate MX.

---

## 4. Verification & Testing Checklist

- [ ] Unit tests in `pkg/mail/warmup/ramp_curves_test.go` verifying curve monotonicity and upper volume bounds.
- [ ] Schedule timeline tests validating timezone dispersion math across UTC-8 to UTC+10.
- [ ] SVG markup generator tests verifying valid XML and zero script injections.
- [ ] Seed interaction tests ensuring replies maintain threading headers (`In-Reply-To`, `References`).
- [ ] Coding guidelines audit: <= 15 line functions, positive booleans, error return sovereignty.
