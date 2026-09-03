# Dynamic Timeline & Timeout (`-t` / `--timeout`)

## Overview

The dynamic timeline and timeout flag (`-t`, `--timeout`, `--timeline`) eliminates the need for manual retry seconds. Instead of requiring users or scripts to guess an arbitrary delay (e.g. `retry 30`), passing `-t` instructs GitMap to orchestrate an adaptive polling schedule derived directly from the computed historical ETA.

## Operational Contracts

1. **CLI Flag Declaration**:
   - Short flag: `-t`
   - Long flags: `--timeout`, `--timeline`
   - Available on: `gitmap pipeline status`, `gitmap pipeline wait-time`, `gitmap pipeline-ai`.

2. **Adaptive Polling Cadence**:
   The polling interval dynamically adjusts depending on the remaining ETA:
   - When $\text{ETA} > 120\text{s}$: Poll every $20\text{s}$.
   - When $60\text{s} < \text{ETA} \le 120\text{s}$: Poll every $10\text{s}$.
   - When $\text{ETA} \le 60\text{s}$: Poll every $5\text{s}$.

3. **Autonomous Timeout Budget**:
   The dynamic timeout deadline is computed as:
   $$\text{Deadline} = \text{AvgSuccessDuration} \times 1.5 + 60\text{s}$$
   If a pipeline exceeds this deadline without completing, the command exits with code `1` and emits an actionable diagnostic timeout report.

4. **Terminal Output Presentation**:
   When `-t` is enabled, the CLI outputs a dynamic timeline tracking progress:
   ```text
   ● Pipeline [CI] in progress: [=====>     ] 52% (ETA: 45s, Elapsed: 49s)
   ```
