# Issue Domain 05: Pipeline and CI/CD Telemetry

- **Domain:** GitHub Actions Polling, ETA Estimation, and Error Capture
- **Status:** Consolidated Problem & Resolution Matrix

## 1. False Positive "Clean Status" Reports
- **Symptoms:** Failed GitHub Actions runs reported as `(clean status)` with empty error logs.
- **Root Cause:** Deduplication logic keyed only on repository name, allowing older successful runs to mask newer failures.
- **Resolution:** Scoped deduplication key to `branch:sha:workflow`, correctly identifying all failure conclusions.

## 2. Negative or Zero ETA Forecasting
- **Symptoms:** Estimated time remaining (`eta`) showed 0s or negative numbers during long-running builds.
- **Root Cause:** Stage duration calculations subtracted elapsed time from hardcoded minimum thresholds.
- **Resolution:** Switched to moving-average historical stage duration with minimum 10s floor.
