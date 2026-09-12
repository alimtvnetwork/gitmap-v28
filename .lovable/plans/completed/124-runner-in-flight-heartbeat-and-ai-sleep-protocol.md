# 124-runner-in-flight-heartbeat-and-ai-sleep-protocol.md: Runner 25s In-Flight Heartbeat & 1-Minute AI Agent Sleep/Wait Protocol

**Status: completed**

> **Task Initiation & Execution:**
> Initiated in response to user directive:
> 1. In-flight heartbeat interval increased to **25 seconds** (previously 10.0s) in `03-ai-scripts/06-cicd-local-runner.py` because fast checking produces excessive terminal noise.
> 2. AI agents and runner prompts updated with a mandatory **1 minute (60 seconds) sleep/wait protocol** each time, or dynamically sleeping for the remaining duration read from `.lovable/temp/runner-eta.json` (or based on previous known approximate total delay) instead of busy-polling.
> 3. Synchronized across `01-prompts/`, `.agents/skills/`, `.lovable/memory/`, and connected workspace `coding-guidelines`.
> 4. Fixed remaining `gitmap` folder/cwd references in runner `JOB_BATCHES` to `cli`.
> Completed autonomously in a single consolidated iteration without triggering releases or uncommanded unit tests.

## 1. Problem Diagnosis & Architectural Goal
1. **Terminal Noise from Rapid Heartbeats**:
   - `TelemetryTracker` in `03-ai-scripts/06-cicd-local-runner.py` was previously configured with a 10.0-second heartbeat interval, resulting in repetitive terminal lines and excessive log volume.
   - The user provided terminal artifact `.lovable/assets/terminals/04-runner-in-flight-interval.png` showing rapid emissions.
   - Solution: Set `DEFAULT_HEARTBEAT_INTERVAL = float(os.environ.get("RUNNER_HEARTBEAT_INTERVAL", 25.0))`, add CLI option `--heartbeat-interval`, and ensure `tick()` rate limits emissions strictly to 25 seconds or more unless forced by the background heartbeat thread.
2. **AI Agent Busy-Polling**:
   - When AI agents execute background test runs or long-running commands, polling every few seconds wastes token budget and generates excessive tool calls.
   - Solution: Prompts and skills now strictly mandate sleeping for **1 minute (60 seconds)** per iteration, or dynamically sleeping for `remaining_eta_sec` from `.lovable/temp/runner-eta.json` / previous total approximate delay.
3. **Runner Job Batch CWD Parity**:
   - Following the Plan 123 migration (`gitmap/` -> `cli/`), `JOB_BATCHES` entries in `03-ai-scripts/06-cicd-local-runner.py` still contained `cwd: "gitmap"` and `"-C", "gitmap"`. These were updated to `"cli"`.

## 2. Task-Specific Rule Set (Enforced)
1. **Rule H1 (25-Second Heartbeat Cadence)**: The runner's in-flight telemetry emits strictly every 25 seconds or more.
2. **Rule H2 (1-Minute AI Sleep Discipline)**: Agents checking background runners MUST sleep for 1 minute (60s) or wait for remaining ETA before inspecting task status.
3. **Rule H3 (Zero-Poll Protocol)**: AI agents must never busy-loop or query task status in rapid succession.
4. **Rule H4 (Strict Relative Paths)**: Zero absolute OS paths (`D:/...`, `C:/...`) or `file:///` URIs in any markdown or committed document.
5. **Rule H5 (Zero Releases & No Uncommanded Tests)**: No version bumping, changelog editing, or test running during routine turns.

## 3. Subtask Consolidation & Deliverables

### Subtask 1: CI/CD Local Runner In-Flight Telemetry Update (Completed)
- In `03-ai-scripts/06-cicd-local-runner.py`:
  - Added `DEFAULT_HEARTBEAT_INTERVAL = float(os.environ.get("RUNNER_HEARTBEAT_INTERVAL", 25.0))`.
  - Updated `TelemetryTracker.__init__` with default `heartbeat_interval=DEFAULT_HEARTBEAT_INTERVAL`.
  - Added `--heartbeat-interval` CLI option to `add_reporting_arguments`.
  - Passed `heartbeat_interval` dynamically in `prepare_runner_context`.
  - Updated `tick(force=False)` condition: rate-limits output strictly every 25 seconds or more unless `force=True` from the heartbeat timer.
  - Replaced remaining `gitmap` folder and `-C gitmap` references in `JOB_BATCHES` with `cli`.
  - Verified compilation via `python -m py_compile 03-ai-scripts/06-cicd-local-runner.py`.

### Subtask 2: Primary Repository Prompts & Skills Synchronization (Completed)
- Updated Runner In-Flight ETA Wait Protocol in `gitmap`:
  - `01-prompts/16-ci-cd/01-ci-cd-fix.md` (Rule 7)
  - `01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md` (Rule 7)
  - `01-prompts/14-execute/02-execute-parent-task-with-n-steps.md` (Rule 8 & Phase 2)
  - `.agents/skills/ci-cd-fix/skill.md` (Rule 7)
  - `.agents/skills/autonomous-qa-and-testing/skill.md` (Rule 12)
  - `.agents/skills/execute-parent-task/skill.md` (Checklist item)
  - `.agents/skills/execute-parent-task-with-n-steps/skill.md` (Checklist item)
  - `.lovable/memory/learned/14-smart-test-runner-and-temp-isolation.md` (Section D)

### Subtask 3: Connected Workspace (`coding-guidelines`) Synchronization (Completed)
- Synchronized matching prompts and skills in connected workspace:
  - `01-prompts/16-ci-cd/01-ci-cd-fix.md`
  - `01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md`
  - `01-prompts/14-execute/02-execute-parent-task-with-n-steps.md`
  - `.agents/skills/ci-cd-fix/skill.md`
  - `.agents/skills/autonomous-qa-and-testing/skill.md`
  - `.agents/skills/execute-parent-task-with-n-steps/skill.md`

## 4. Verification Results
- `python -m py_compile 03-ai-scripts/06-cicd-local-runner.py` passed with exit code 0.
- `python linter-scripts/check-nested-ifs.py` and `python linter-scripts/check-enum-and-boolean.py` passing clean.
- Image asset `.lovable/assets/terminals/04-runner-in-flight-interval.png` tracked and persisted.
- Zero absolute paths and zero `file:///` URIs in committed markdown.
- Pushed clean to remote.
