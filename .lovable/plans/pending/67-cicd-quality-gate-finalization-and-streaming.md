# 67-cicd-quality-gate-finalization-and-streaming.md: CI/CD Quality Gate Finalization, Real-Time Streaming & AI Telemetry

## 1. Executive Summary

This plan finalizes the local CI/CD quality gate architecture across all 33 gates enqueued in `03-ai-scripts/06-cicd-local-runner.py` (and `.lovable/ai-fix-scripts/06-cicd-local-runner.py`). It implements real-time failure streaming to `.lovable/temp/cicd/`, instant terminal failure dumps with stack traces and suspect file lists, and verifies that the complete test suite passes with `exit 0`.

## 2. Task-Specific Rules & Constraints

1. **Rule 1 (Real-Time Failure Streaming):** All quality gate failures must immediately append to `.lovable/temp/cicd/errors.log`, `.lovable/temp/cicd/errors.json`, and `.lovable/temp/cicd/run.log` without waiting for the test suite to complete.
2. **Rule 2 (Immediate Terminal Visibility):** When a gate fails, print an immediate ANSI warning banner with failing command, exit code, stack trace, and suspect file paths so an AI agent can start fixing issues concurrently.
3. **Rule 3 (Quiet Success by Default):** In standard execution, passing gates remain quiet or single-line tick to prevent cluttering terminal context for the AI. Full logs are shown only with `--all-paths`.
4. **Rule 4 (No Automatic Releases):** Do NOT bump `version.json`, changelogs, or cut release tags.
5. **Rule 5 (Relative Paths Only):** Zero `file:///` absolute paths in any planning documents or code; all paths must remain strictly relative to repository root.

## 3. Architecture & Components

### A. Real-Time Telemetry Files (`.lovable/temp/cicd/`)
- `errors.log`: Human-readable Markdown log appending failures as they happen with timestamps, commands, exit codes, and extracted stack traces.
- `errors.json`: Atomic JSON array containing structured error objects.
- `run.log`: Full chronological stream logging every gate (pass or fail) as it finishes.
- `summary.json`: High-level run state (`running`, `completed`, `failed`), timestamps, gate counts (`total_gates`, `passed_gates`, `failed_gates`, `remaining_gates`), and paths to all stream files.

### B. Quality Gate Hardening
- `gitmap/logging`: `jsonlog_test.go` verifies logger defaults, no-op disabled mode, and level outputs (100% coverage).
- `gitmap/visibility`: `pattern_test.go`, `exclude_test.go`, and `fuzzy_test.go` verify pattern compilation, owner repo matching, version parsing, exclusion ranges, and fuzzy distance (91.8% coverage, beating 75% floor).
- `.github/scripts/coverage-floor.py`: Parses coverage output against `.github/coverage.floor` with module-aware working directory resolution.
- `03-ai-scripts/06-cicd-local-runner.py`: Set `GOTMPDIR` to local `.tmp/` on `D:` drive to prevent disk space exhaustion.

## 4. Subtasks Breakdown

1. `01-runner-live-streaming-and-terminal-output.md`: Refactor runner to implement live `run.log`, live `summary.json`, suspect file extraction, and JSON banner safety.
2. `02-quality-gate-fixes-and-coverage-floors.md`: Verify `logging` and `visibility` tests, `coverage-floor.py`, and remove `-coverpkg=./...` from coverage profile generator.
3. `03-suite-execution-and-verification.md`: Execute full runner suite and verify all 33 quality gates pass cleanly with `exit 0`.

## 5. Acceptance Criteria

- [ ] All 33 quality gates in `03-ai-scripts/06-cicd-local-runner.py` and `.lovable/ai-fix-scripts/06-cicd-local-runner.py` pass cleanly (`exit 0`).
- [ ] Real-time failures stream immediately to `.lovable/temp/cicd/errors.log` and `.lovable/temp/cicd/errors.json`.
- [ ] `run.log` and `summary.json` update in real time with active counters.
- [ ] No swallowed errors (`except Exception: pass` replaced with specific exceptions or warnings).
- [ ] Zero automatic releases or version bumps.
