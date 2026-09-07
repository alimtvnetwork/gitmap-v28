# Subtask 03: In-Flight Parallel Terminal Telemetry & Immediate Failure Reporting

## Objective
Implement dual-mode parallel worker in-flight terminal display, rate-limited heartbeats, immediate failure banners, and detailed AI instruction headers.

## Requirements
1. **In-Flight Parallel Telemetry**:
   - Active worker slot tracking: Track currently executing gates in real-time.
   - For TTY terminals: Display in-flight active gates and elapsed time dynamically.
   - For Non-TTY / AI subagent / CI: Emit periodic rate-limited progress line (e.g. `[IN-FLIGHT] 4 active: [Spell Check (2.1s), Web App Build (4.5s)...] | 12/33 done`).
   - Default run remains quiet on success (no passing log spam).
2. **Immediate Failure Output**:
   - On gate failure, immediately interrupt and print ANSI failure banner containing:
     - Gate Name and Exit Code
     - Failing Command and Working Directory
     - Suspect Files (extracted from error traces)
     - Full Stderr / Stdout Trace
   - Stream failure to `errors.log` and `errors.json` without delay.
3. **AI Agent Instructions Header**:
   - Update runner module docstring with clear instructions explaining real-time failure streaming, how to monitor `.lovable/temp/cicd/`, and parallel remediation.
4. **Coding Guidelines**:
   - All functions $\le 15$ lines.
   - Blank line before every return statement.
   - Affirmative booleans (`is_*`, `has_*`).
   - Zero swallowed exceptions.

## Target Files
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [x] Active parallel gates visible during execution.
- [x] Passing gates do not flood stdout when run without `--all-paths`.
- [x] Failures immediately print full stack trace and suspect files to stdout.
- [x] Module docstring contains comprehensive AI instructions.
