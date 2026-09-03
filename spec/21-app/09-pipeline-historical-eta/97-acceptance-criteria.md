# Acceptance Criteria: Pipeline Historical ETA & Diagnostics

## Overview

Verification contracts and automated acceptance tests for the Historical Success Baseline ETA model, dynamic `-t` flag, and diagnostic error extraction.

## Scenario 1: ETA Excludes Failed and Canceled Runs

```gherkin
Given a set of past workflow runs containing 3 successful runs (averaging 120s) and 2 failed runs (failing at 10s)
When calculateETA is invoked for an active run created 30s ago
Then the ETA baseline is computed as 120s (not 76s)
And the returned remaining ETA is 90s (120s - 30s)
```

## Scenario 2: Dynamic Timeline Flag (`-t`)

```gherkin
Given an active workflow run with an ETA of 60s
When gitmap pipeline wait-time -t is executed
Then the command computes the timeline and auto-polls with an adaptive cadence
And terminates with exit code 0 when the pipeline succeeds
```

## Scenario 3: Error Diagnostics on Pipeline Failure

```gherkin
Given a failed workflow run containing 2,000 lines of mixed logs
When gitmap pipeline status is executed
Then only the matching failure lines (e.g. "##[error]", "--- FAIL:") are extracted
And the diagnostic error block is printed directly under the status card
```
