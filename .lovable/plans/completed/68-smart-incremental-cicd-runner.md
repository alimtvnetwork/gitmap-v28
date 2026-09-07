# 68-smart-incremental-cicd-runner.md: Smart Incremental CI/CD Runner with Resumption, In-Flight Telemetry & Gate Skipping

## 1. Executive Summary

This plan upgrades `03-ai-scripts/06-cicd-local-runner.py` and its synchronized mirror `.lovable/ai-fix-scripts/06-cicd-local-runner.py` into an intelligent, crash-resilient, incremental test runner. The enhanced runner retains phase execution state, fingerprinting inputs (git commit SHA, uncommitted working tree deltas, linter script timestamps, and configuration files) to skip previously passing, unchanged gates in O(1) time (~0.5ms per gate). It features a multi-session storage architecture (`.lovable/temp/cicd/runs/<timestamp>/` and `.lovable/temp/cicd/latest/`), dual-mode parallel telemetry with in-flight worker slot tracking, zero-delay terminal failure reporting with suspect file extraction, and automatic crash resumption.

## 2. Task-Specific Rules & Constraints

1. **Rule 1 (Strict No Automatic Releases):** Under NO circumstances bump versions in `version.json`, edit changelogs, or cut release tags. Commits must remain standard development commits.
2. **Rule 2 (Strict Relative Paths):** Zero `file:///` URIs or machine-specific absolute paths in code, plans, or telemetry. All paths must be relative to repository root (`/`).
3. **Rule 3 (Atomic Disk Persistence):** State and JSON reports must be written using atomic `write-tmp-flush-rename` (`os.replace`) to prevent corruption if interrupted or killed mid-write.
4. **Rule 4 (Quiet Success by Default & Active In-Flight Telemetry):** Standard runs must not spam passing statements. Show live in-flight worker state (e.g. 4-5 concurrent running gates with elapsed timers) and report full details (command, code, stack trace, suspect files) immediately when a failure occurs.
5. **Rule 5 (Crash Resumption & Cache Controls):** If interrupted or crashed, re-running automatically resumes from the last persistent `state.json`. Provide `--force` (`--fresh`, `--clean`) to purge cache and run all 33 gates from scratch.
6. **Rule 6 (Coding Guidelines Compliance):** Functions $\le 15$ lines, blank line before return statements, positive booleans (`is_*`, `has_*`), no magic numbers, and zero swallowed exceptions (`except Exception: pass` is strictly forbidden).

## 3. Architecture & Subsystems

### A. Session Scaffolding & Storage Subsystem (`.lovable/temp/cicd/`)
- Isolated run directories under `.lovable/temp/cicd/runs/<timestamp>/`.
- Symlink / Directory Junction at `.lovable/temp/cicd/latest/` with fallback to `latest_run.txt`.
- Backward-compatible root mirror files (`errors.log`, `errors.json`, `run.log`, `summary.json`, `state.json`) so external tools and subagents continue working seamlessly.

### B. Two-Tier Incremental Fingerprinting Engine
- **Tier 1: Git Delta Resolution**: Inspects `git rev-parse HEAD` and `git status --porcelain=v1 -uall` to compute repository delta set $\Delta_{repo}$.
- **Tier 2: Gate Dependency Mapping**: 33 gates mapped to declarative `GateSpec` instances tracking tool scripts, configs, relevant path globs, exclusion globs, and upstream artifacts. Unchanged passing gates are skipped instantly.

### C. Adaptive Dual-Mode Parallel Telemetry
- Interactive TTY: Multi-slot in-flight progress dashboard showing active worker threads, running gate names, and elapsed timers.
- Non-Interactive / AI Subagent: Periodic rate-limited heartbeat lines (e.g. `[IN-FLIGHT] 4 active: [Gate A (3.2s), Gate B (1.8s)...] | 18/33 done`).
- Zero-Delay Failure Interrupts: Immediate ANSI failure banners with command, exit code, suspect files, and stack traces printed to terminal and streamed to disk.

### D. Crash Recovery & Resumption Engine
- Records gate lifecycle transitions: `PENDING` -> `RUNNING` -> `PASSED` / `FAILED` / `TIMEOUT` / `SKIPPED`.
- On start, detects prior incomplete sessions (PID dead or stale heartbeat). Automatically marks interrupted gates as `CRASHED` and queues them for re-execution while preserving `PASSED` gates.

## 4. Subtasks Breakdown

1. [01-session-scaffolding-and-state-persistence.md](subtasks/68-smart-incremental-cicd-runner/01-session-scaffolding-and-state-persistence.md): Multi-session directory hierarchy, atomic `os.replace` StateStore, PID tracking, and backward-compatible root mirrors.
2. [02-fingerprinting-and-incremental-cache-engine.md](subtasks/68-smart-incremental-cicd-runner/02-fingerprinting-and-incremental-cache-engine.md): Two-tier Git change detector, 33 GateSpec mapping, cache persistence, and skip decision logic.
3. [03-in-flight-terminal-display-and-failure-reporting.md](subtasks/68-smart-incremental-cicd-runner/03-in-flight-terminal-display-and-failure-reporting.md): Active parallel worker slot display, rate-limited non-TTY heartbeat, immediate failure banners, and AI agent docstring header.
4. [04-runner-verification-and-resumption-test.md](subtasks/68-smart-incremental-cicd-runner/04-runner-verification-and-resumption-test.md): Full execution verification, caching validation, `--force` test, and sync to `.lovable/ai-fix-scripts/06-cicd-local-runner.py`.

## 5. Acceptance Criteria

- [x] All 33 quality gates pass cleanly with `exit 0`.
- [x] Subsequent run without modifications skips previously passed gates in < 2 seconds.
- [x] Running with `--force` (or `--fresh` / `--clean`) clears cache and re-executes all 33 gates.
- [x] Terminal displays in-flight active gates during execution without flooding stdout with passing gate logs.
- [x] Any gate failure prints immediate stack trace, command, exit code, and suspect files to terminal and logs to disk.
- [x] Session runs are persisted under `.lovable/temp/cicd/runs/<timestamp>/` with updated `state.json`.
- [x] Strict coding guidelines (<= 15 line functions, blank line before returns, affirmative booleans, no swallowed exceptions) enforced.
