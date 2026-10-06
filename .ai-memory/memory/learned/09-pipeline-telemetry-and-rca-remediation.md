# 09 — Pipeline Telemetry and 4-Part RCA Remediation

- **Subsystem:** CI/CD Diagnostics & Pipeline Intelligence
- **Status:** Authoritative Reference

## 1. Pipeline Error Extraction (`pe`)
- Automatically captures failing job logs from GitHub Actions and local runner pipelines.
- Isolates traceback heatmaps and failing unit test assertions.

## 2. 4-Part Root Cause Analysis (RCA) Protocol
- **Part 1: Symptoms:** Exact console output, failure exit codes, and affected environments.
- **Part 2: Underlying Mechanism:** Concrete code-level fault, race condition, or contract breach.
- **Part 3: Surgical Fix:** Precise code changes applied to resolve the defect.
- **Part 4: Verification Evidence:** Verification log confirming green passage.
