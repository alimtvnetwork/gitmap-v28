# Errorlogs Dynamic Timeline and Real-Time Watch Loop

## Context

Developers and automated agents operating in continuous integration workflows need an immediate, zero-friction mechanism to watch active pipelines to completion and extract only actionable error logs without manual polling loops or terminal noise.

## Architectural Specification

### 1. Command Dispatching & Aliases

The command suite is registered under both `gitmap pipeline` and root-level dispatch tables:
- `gitmap pipeline errorlogs`
- `gitmap pipeline error-logs`
- `gitmap errorlogs`
- `gitmap error-logs`
- Short aliases: `error-log`, `errorlog`, `errors`, `err`

### 2. Timeline Flag Resolution

The watch behavior is activated via:
- `-t`
- `--timeline`
- `--timeout`
- `-w`
- `--watch`

### 3. Adaptive Timeline Cadence

When an active pipeline run is in progress or queued:
1. GitMap enters a polling loop reporting current ETA countdown and elapsed time:
   `⏳ [Workflow Name] in progress: ETA ~Xs (elapsed: Ys)`
2. Sleep intervals adapt dynamically to remaining time:
   - Remaining ETA > 120s: 15-second interval
   - Remaining ETA > 60s: 10-second interval
   - Remaining ETA <= 60s: 5-second interval
3. Once the workflow transitions to a terminal state (`completed`), polling immediately ceases.

### 4. Failure Log Extraction

Upon workflow termination:
- If `conclusion == "success"`: Emits green completion indicator and baseline run metrics.
- If `conclusion == "failure"`: Queries GitHub API for failed step logs, strips ANSI escape sequences, deduplicates repetitive errors, and renders the high-visibility CI/CD Error Diagnostics card.
