# Subtask 01: Predictive Analytics, Bayesian Curves & ISP Threshold Cliffs

## Status: READY
- **Subtask Slug:** `01-predictive-analytics-and-bayesian-curves`
- **Parent Plan:** `.ai-memory/plans/pending/341-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/01-improvements-401-to-430-predictive-analytics-and-bayesian-curves.md`
- **Target Range:** Improvements 401–430 (Clusters 41–43)
- **Target Subsystem:** `pkg/mail/predictive/`, `pkg/mail/rhythm/`
- **Exclusive Owned Files:**
  - `pkg/mail/predictive/bayesian_curve.go`
  - `pkg/mail/predictive/threshold_cliff.go`
  - `pkg/mail/predictive/decay_modeler.go`
  - `pkg/mail/predictive/ramp_scheduler.go`
  - `pkg/mail/predictive/cadence_optimizer.go`
  - `cli/cmdwarmup/predictive_cmd.go`

---

## 1. Context & Objectives

This subtask delivers the predictive deliverability layer and mathematical Bayesian volume ramp curves:
1. **Bayesian Warming Curve Calculation (Cluster 41, Items 401–410):** Algorithmic regression models forecasting ISP reception ceilings and dynamically adjusting daily sending ramps based on historical inbox placement.
2. **ISP Threshold Cliff & Anomaly Detection (Cluster 42, Items 411–420):** Statistical detection of sudden ISP throttling cliffs (Google rate-limiting, Yahoo deferrals) with predictive backoff trigger curves.
3. **Seed Engagement Decay & Cadence Optimization (Cluster 43, Items 421–430):** Modeling seed engagement probability decay curves and optimizing temporal cadence to prevent inbox fatigue.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isCurveOptimal`, `hasCliffDetected`, `isRampPermitted`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Bayesian Volume Ramp & Curve Generation
- Implement Bayesian regression calculator in `pkg/mail/predictive/bayesian_curve.go`.
- Define non-linear daily sending volume curve generator in `pkg/mail/predictive/ramp_scheduler.go`.

### Step 2: ISP Threshold Cliff Modeling & Real-Time Backoff
- Construct ISP-specific reception cliff detector in `pkg/mail/predictive/threshold_cliff.go`.
- Connect predictive backoff triggers to warm-up dispatcher queues.

### Step 3: Seed Engagement Decay & CLI Surface
- Build seed interaction decay modeler in `pkg/mail/predictive/decay_modeler.go`.
- Expose predictive forecast commands via `cli/cmdwarmup/predictive_cmd.go`.
