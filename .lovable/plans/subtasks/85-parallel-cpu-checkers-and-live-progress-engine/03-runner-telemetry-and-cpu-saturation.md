# Subtask 03: Runner Telemetry and CPU Saturation

## Objective
Enhance `03-ai-scripts/06-cicd-local-runner.py` to ensure `TelemetryTracker` does not display stagnant 0% progress, sets concurrency to saturate all CPU cores, and streams child checker output so intermediate percentages are visible.

## Files to Touch
- `03-ai-scripts/06-cicd-local-runner.py`

## Implementation Steps
1. In `03-ai-scripts/06-cicd-local-runner.py`:
   - Verify `CPU_CORES = os.cpu_count() or 16` and ensure `max_workers` across CPU-bound batches (Batch 1 Linters & AST Checks) uses `CPU_CORES`.
   - Update `TelemetryTracker.tick()`:
     - When `self.completed_count == 0` but jobs are running, show active elapsed duration and in-progress status instead of raw stagnant `0% done`.
     - Display:
       `[IN-FLIGHT] 8 active: [Go Format Check (1.2s), Nested If (0.8s), ...] | running batch 1 (Linters & AST)`
   - When running subcommands in `run_job`, ensure subprocess output is unbuffered so child tools emitting progress lines flush immediately.
2. Run sanity tests on runner flags (`--filter`, `--all-paths`).
