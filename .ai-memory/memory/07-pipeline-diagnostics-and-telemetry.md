# 07 — Pipeline Telemetry, Diagnostics & 4-Part RCA Remediation

- **Domain:** CI/CD Telemetry, Pipeline Failure Extraction & Root Cause Analysis
- **Authoritative Specification:** [07-pipeline-and-diagnostics](../../02-spec/21-app/07-pipeline-and-diagnostics/01-architecture-spec.md)
- **Status:** Active & Ratified

---

## 1. CI/CD Pipeline Telemetry (`gitmap pe`, `gitmap pea`)

GitMap provides real-time and historical telemetry for local and remote CI/CD pipeline runs:

- **Pipeline Error Commands:**
  - `gitmap pe`: Displays compact failure summaries for the most recent pipeline run.
  - `gitmap pea`: Displays aggregated failure logs across all recorded pipeline runs.
  - `gitmap pe -t`: Renders terminal traceback mode with syntax highlighting for failed assertions.
- **Negative Commit Offsets:** Supports inspecting failures from prior commits using relative offsets (`gitmap pe -1`, `gitmap pe -2`).
- **OK-Line Suppression:** Filters out noisy passing tests (`--- PASS:`), compilation successes, and dependency download lines by default, reducing log volume by over 80%.

---

## 2. Bounded Stack Trace Extraction

When a pipeline job fails, the error extractor isolates the critical failure frame:

- **Extraction Window:** Isolates 5 leading lines of context + the failing assertion line + 20 trailing lines.
- **Exit-Code Halting:** Stops extraction upon encountering process termination boundaries (e.g., `exit status 1`), preventing irrelevant trailing teardown logs from contaminating the failure record.
- **Sanitized Paths:** Replaces local runner machine paths with relative Git paths to allow agents to find target files instantly.

---

## 3. Dynamic ETA Sleep Synchronization & Heatmap Telemetry

- **`runner-eta.json` Cache:** Records historical run durations for each test suite and CI stage.
- **Dynamic ETA Sleep:** Instead of polling the pipeline in high-frequency CPU-spinning loops, the runner calculates estimated completion time and sleeps until 90% of expected duration before checking status.
- **Heatmap Telemetry:** Summarizes failing test frequency across suites, highlighting flaky tests and regression hotspots in an ANSI terminal heatmap table.

---

## 4. 4-Part Root Cause Analysis (RCA) Framework

All automated and manual bug fixes in GitMap must follow the grounded 4-Part RCA protocol:

```markdown
### 1. Root Cause
- Concrete analysis of the bug without speculation.
- Identify the exact line of code, failing condition, or race condition.

### 2. Blast Radius
- Enumerate all affected packages, commands, or data tables.
- Identify downstream components dependent on the failing contract.

### 3. Surgical Fix
- Precise, minimal code changes addressing only the root cause.
- Adhere strictly to coding guidelines (positive booleans, <= 15 line functions).

### 4. Verification Gate
- Exact terminal commands to verify the resolution.
- Confirm zero regressions in neighboring tests and static linters.
```

---

## 5. Zero-Storage GitHub Actions Governance

- **Zero-Storage Rule:** GitHub Actions workflows must not store build artifacts or test output caches on GitHub infrastructure (maintaining 0.0 GB storage usage).
- **Automated Purge:** AI scripts monitor workflow runs and purge generated artifacts upon completion.
