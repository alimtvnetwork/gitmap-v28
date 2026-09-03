# Historical Success Baseline ETA Algorithm

## Overview

The historical baseline ETA algorithm calculates the estimated remaining completion time for active workflow pipelines and segments by referencing only previous successful executions.

## Problem Statement

Previous ETA heuristics computed an average runtime across all completed workflow runs, regardless of conclusion. When a workflow was cancelled early (e.g. after 5 seconds) or failed fast due to a syntax error (e.g. after 12 seconds), these runs severely skewed the average downward. Consequently, a pipeline that normally requires 4 to 6 minutes was estimated to complete in 35 seconds, leading to premature timeouts and confusing status readouts.

## Core Invariants & Algorithm

1. **Strict Success Filter**:
   Only runs satisfying `status == "completed"` and `conclusion == "success"` are considered valid samples for calculating the typical runtime.

2. **Duration Computation**:
   For each qualified run \(R_i\):
   $$\text{Duration}(R_i) = \text{UpdatedAt}(R_i) - \text{CreatedAt}(R_i)$$
   Runs with duration under 10 seconds are flagged as anomalous cache hits and excluded from baseline computation.

3. **Segment & Workflow Baseline Average**:
   Given $M$ valid successful runs:
   $$\text{AvgSuccessDuration} = \frac{1}{M} \sum_{i=1}^{M} \text{Duration}(R_i)$$
   If $M = 0$, the system applies standard fallback baselines based on the workflow category:
   - Release pipelines: 95s
   - CI / Test suites: 180s
   - Other workflows: 90s

4. **Remaining ETA Calculation**:
   For the active workflow run $R_{\text{active}}$:
   $$\text{Elapsed} = \text{Now} - \text{CreatedAt}(R_{\text{active}})$$
   $$\text{ETA}_{\text{raw}} = \text{AvgSuccessDuration} - \text{Elapsed}$$
   $$\text{ETA} = \max(\text{ETA}_{\text{raw}}, \text{GraceFloor})$$
   Where $\text{GraceFloor} = 15\text{s}$ ensures the countdown does not display negative or zero values while the pipeline is still confirmed to be executing.
