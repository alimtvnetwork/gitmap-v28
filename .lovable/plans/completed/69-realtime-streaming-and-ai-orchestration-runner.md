# 69-realtime-streaming-and-ai-orchestration-runner.md: Real-Time Streaming Telemetry, Unbuffered Event Pipeline & AI Parallel Orchestration

## 1. Executive Summary

This plan elevates `03-ai-scripts/06-cicd-local-runner.py` (and `.lovable/ai-fix-scripts/06-cicd-local-runner.py`) into a truly unbuffered, parallel-remediation CI/CD streaming engine. It guarantees that any test or quality gate failure is immediately broadcast to the terminal and streamed to `.lovable/temp/cicd/` without OS buffer delays using `os.fsync`. It introduces an append-only line-delimited `events.jsonl` stream and `changelog.log`, enhances suspect file extraction across all 33 gates (including Windows drive letters, TypeScript `(line,col)`, Python tracebacks, and semantic git delta fallbacks), displays `Working Dir` and `Env Overrides` in terminal failure banners, provides an end-of-run `AI AGENT REMEDIATION SUMMARY & LOG LOCATIONS` dashboard, and equips the runner with an Active Operational Instruction Manual in the module header for AI agents working in parallel.

## 2. Task-Specific Rules & Constraints

1. **Rule 1 (Strict No Automatic Releases):** Under NO circumstances bump versions in `version.json`, edit changelogs, or cut release tags. Commits must remain standard development commits.
2. **Rule 2 (Strict Relative Paths):** Zero `file:///` URIs or machine-specific absolute paths in code, plans, or telemetry. All paths must be relative to repository root (`/`).
3. **Rule 3 (Unbuffered fsync Disk Flushing):** All file writes to `.lovable/temp/cicd/` (`errors.log`, `run.log`, `events.jsonl`, `changelog.log`) must call `flush()` and `os.fsync(fh.fileno())` under a write lock so concurrent tailing agents see every byte instantaneously.
4. **Rule 4 (Windows Sharing Violation Resilience):** `atomic_write_text` must employ exponential backoff retries and fallback writes for `PermissionError` [WinError 32] when background agents hold read handles on JSON files.
5. **Rule 5 (Immediate ANSI Failure Interrupt & Post-Run Remediation Banner):** Failures must immediately interrupt terminal output with full command, cwd, env, exit code, suspect files, and stack traces. When execution finishes with errors, a consolidated AI remediation banner must show exact log file paths and copy-pasteable single-gate re-test commands.
6. **Rule 6 (Coding Guidelines Compliance):** Functions $\le 15$ lines, blank line before every return statement, affirmative booleans (`is_*`, `has_*`), no magic numbers, and zero swallowed exceptions (`except: pass` strictly forbidden).

## 3. Subsystems & Architecture

### A. Universal Synchronized Disk Streaming Engine (`.lovable/temp/cicd/`)
- `direct_append_sync(file_path, content)`: Thread-safe, unbuffered append with `os.fsync` synchronization.
- `events.jsonl`: Append-only NDJSON stream (`run_started`, `gate_started`, `gate_cached`, `gate_passed`, `gate_failed`, `gate_timeout`, `run_completed`).
- `changelog.log`: High-level human-readable event log.
- `summary.json`: Initialized immediately on runner startup with `"status": "running"` to eliminate race conditions with prior runs.

### B. Multi-Tier Suspect File Extraction & CWD Resolution
- Support for colon syntax (`gitmap/cmd/root.go:42:10`), Windows drive letters (`D:\...`), TypeScript parentheses (`src/App.tsx(45,12)`), Python tracebacks (`File "...", line X`), and Go panics/races.
- Automated resolution against `cwd` (e.g., prefixing `gitmap/` if run inside `cwd="gitmap"`).
- Semantic fallback: If no files in error text, extract gate's relevant patterns intersecting `repo_delta`, then tool scripts and configs.

### C. Enhanced Terminal Failure Banner & AI Remediation Summary
- Immediate failure banner displaying `Command`, `Working Dir`, `Env Overrides`, `Exit Code`, `Failing Files`, `Stream Log`, `Stream JSON`, `Stream Events`, and stack trace.
- End-of-execution `AI AGENT REMEDIATION SUMMARY & LOG LOCATIONS` banner aggregating deduplicated suspect files, exact log file paths, targeted single-gate re-test commands, and autonomous remediation instructions.

### D. Active Operational Instruction Manual Docstring
- Structured header explaining the 3-batch pipeline model, artifact locations, parallel AI remediation playbook, incremental caching mechanics, and CLI cheat sheet.

## 4. Subtasks Breakdown

1. [01-unbuffered-streaming-and-event-pipeline.md](subtasks/69-realtime-streaming-and-ai-orchestration-runner/01-unbuffered-streaming-and-event-pipeline.md): Thread-safe `direct_append_sync` with `os.fsync`, `events.jsonl` and `changelog.log` event dispatchers, Windows file lock retry in `atomic_write_text`, and startup `summary.json` initialization.
2. [02-suspect-path-resolution-and-job-context.md](subtasks/69-realtime-streaming-and-ai-orchestration-runner/02-suspect-path-resolution-and-job-context.md): Propagate `cwd` and `env_overrides` into `JobResult`, implement multi-tier `extract_suspect_files` with regex expansion and semantic git delta fallback.
3. [03-immediate-terminal-banner-and-remediation-summary.md](subtasks/69-realtime-streaming-and-ai-orchestration-runner/03-immediate-terminal-banner-and-remediation-summary.md): Update `format_failure_banner` with cwd/env/events, implement `format_ai_remediation_banner`, and integrate into runner conclusion.
4. [04-active-operational-header-and-full-suite-verification.md](subtasks/69-realtime-streaming-and-ai-orchestration-runner/04-active-operational-header-and-full-suite-verification.md): Update module docstring with operational instruction manual, verify all 33 gates and guidelines (AST function lengths $\le 15$, blank lines before returns, zero swallowed exceptions), and sync mirror.

## 5. Acceptance Criteria

- [x] All file streams (`errors.log`, `run.log`, `events.jsonl`, `changelog.log`) flushed with `os.fsync` immediately.
- [x] `events.jsonl` records lifecycle events (`run_started`, `gate_started`, `gate_passed`, `gate_failed`, `gate_cached`, `run_completed`).
- [x] `summary.json` initialized immediately with `"status": "running"`.
- [x] Terminal failure banners display failing command, `Working Dir`, `Env Overrides`, exit code, suspect files, and stack trace.
- [x] Suspect file extractor handles Windows drive letters, TypeScript formats, Python tracebacks, and falls back to git delta.
- [x] Post-run AI remediation summary outputs exact relative log paths, deduplicated suspect files, and single-gate re-test commands.
- [x] Module docstring serves as an operational instruction manual for AI agents.
- [x] All 33 quality gates pass cleanly (`exit 0`), and clean incremental run finishes in $<0.5\text{s}$.
- [x] All functions $\le 15$ lines, blank line before returns, affirmative booleans, no swallowed exceptions.
