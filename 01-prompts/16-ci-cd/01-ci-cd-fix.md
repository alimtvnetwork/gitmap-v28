# CI/CD Fix Loop with 4-Part RCA & Local Runner — Workflow (must follow)

Trigger Keywords & Aliases: `fix with RCA`, `fix`, `fix, fix`, `CI/CD fix`, `cicd fix`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal First `N/2` steps (Phase 1): Review the central CI/CD pipeline definitions (`.github/workflows`, `.gitlab-ci.yml`, etc.) and cross-reference them with the local Python runner (`03-ai-scripts/06-cicd-local-runner.py`).
   - **Condition:** If `03-ai-scripts/06-cicd-local-runner.py` does not exist, you must create it immediately.
   - **Condition:** You must ensure that **every single CI/CD case** that needs to run in the pipeline can also be run locally from this Python script (with Docker stripped for native host execution). Improve the Python script to cover all cases if any are missing.
2. [ ] /goal Second `N/2` steps (Phase 2): Run the local runner script (`python 03-ai-scripts/06-cicd-local-runner.py --all`) to catch all errors. Singly execute the script in an autonomous self-loop, zeroing in on one failing error per turn (4-part RCA -> surgical fix -> guideline autofixer -> re-verify).
3. [ ] /goal Finalize CI/CD: Your ultimate goal is to fix and finalize the CI/CD. You must loop until the Python local runner script executes flawlessly with **no errors** (exit code 0) for all registered cases. Do not stop until this goal is met.
4. [ ] /learn Ingest `.lovable/cicd-issues/` for domain-specific architectural specifications.
5. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
6. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
7. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
8. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
9. [ ] /learn Ingest `spec/03-error-manage/` for error handling architectures and AppError.
10. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Screenshot Pipeline Discovery, Update 06-cicd-local-runner.py, Register New JOBS)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Singly-Done Self-Loop Fixing, Zero in on Errors, 4-Part RCA, Green Gate Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before any execution, check if this prompt is installed as a native Antigravity Skill.

1. Check if `.agents/skills/ci-cd-fix/skill.md` exists.
2. If it does NOT exist, create it now. Write the core instructions of this prompt to `.agents/skills/ci-cd-fix/skill.md` with frontmatter:
   ```yaml
   ---
   name: ci-cd-fix
   description: >-
     Use this skill to autonomously diagnose, fix, and verify CI/CD pipelines using local runner scripts, 4-part RCA, and self-looping.
   ---
   ```
3. Once installed, load it on-demand via progressive disclosure for all future runs.

---

## Strict In-Repository Execution & `.lovable/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/` and `.lovable/cicd-issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.lovable/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `spec/03-error-manage/01-index.md`, `.lovable/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands](spec/13-generic-cli/01-index.md)`
> 4. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.
> 5. **Cross-Platform Python CI Mandate (TOTAL BAN on new `.sh` scripts in CI):** All newly created or refactored CI/CD verification tools, determinism checks, fixtures, and linter jobs MUST be implemented in pure, cross-platform Python (`.py`). Legacy `.sh` scripts must be converted to `.py` scripts so all pipelines run natively across Linux, macOS, and Windows without relying on bash emulation.
> 6. **Temp & Failure Folder Isolation:** All temporary directories, runner caches, and test artifacts MUST be strictly placed in `.lovable/temp/`. Creating `.tmp/` at the repository root or outside `.lovable/` is strictly forbidden.
>    - Dedicated Failure Directory: `.lovable/temp/failures/` is the dedicated folder where failed tests and failed quality gates write error logs (`<test-or-job-name>.log`).
>    - Passing Tests Completely Silent: Passing tests must produce ZERO filesystem artifacts (zero files written) and remain completely silent in output logs.
> 7. **Runner In-Flight ETA Wait Protocol:** When running background commands, the runner dynamically writes live status and remaining ETA to `.lovable/temp/runner-eta.json` (emitting in-flight heartbeats strictly every 25 seconds or more). If an agent inspects an active background job and it is still running, the agent MUST sleep/wait for **1 minute (60 seconds) each time**, or dynamically sleep for the remaining ETA duration read from `.lovable/temp/runner-eta.json` (or based on previous total approximate delay) instead of busy-polling or querying in loops.
> 8. **Centralized Test Inventory & Incremental Caching:** All unit tests are cataloged in `.lovable/test-inventory.json` with strictly repository-relative paths (`target_file`, `test_file`). First run executes all tests to establish baseline timings; subsequent runs execute incrementally only if target code files or test files change. Slow test threshold defaults to `4.0s` (configurable via `GITMAP_SLOW_TEST_THRESHOLD`).
> 9. **Dual-Queue Worker Pools:** Slow tests run in a dedicated 4-worker pool running at most 2 tests at a time per batch. Fast tests run in a 4-worker pool running at most 4 tests at a time, pulling in chunks of 100 tests from the test inventory queue until all are complete.
> 10. **Dynamic ETA Sleep Protocol:** The AI agent reads `.lovable/temp/runner-eta.json`, sleeps for the estimated wait time rather than looping, and if still active upon waking, re-checks remaining ETA and sleeps again to avoid burning tokens.

---

## Screenshot & Pipeline Discovery Protocol (Execute First When Any Image Is Provided)

> [!IMPORTANT]
> **If the user provides any image — a CI/CD dashboard, a pipeline run view, a failure screenshot, or a log screenshot — you MUST process it FIRST before modifying application code.**
> Images are first-class diagnostic input. Treat every pixel of text in the image as ground truth.

### Bounded Single-Step Self-Loop Sequence (Singly Done — No Overloaded Steps)

Every step must be **singly done** using bounded self-looping turns. Do NOT try to do scanning, fixing, updating runner, and testing in a single turn.

- **Self-Loop Step 1 (Extract Pipeline Name from Screenshot):**
  1. Carefully read the image to extract:
     - The **pipeline or workflow name** (e.g. `"build-and-test"`, `"CI / lint"`, `"Deploy to staging"`, `"test-matrix"`).
     - The **failing job/step name** (marked with ❌ or "Failure").
     - The **error text, log snippets, or stack traces** visible in the screenshot.
  2. Scan repository CI/CD files (`.github/workflows/*.yml`, `.gitlab-ci.yml`, etc.) to locate whatever newly added pipeline jobs, steps, or linter scripts correspond to that pipeline name.

- **Self-Loop Step 2 (FIRST ACTION: Update Python Runner Script):**
  1. Open `03-ai-scripts/06-cicd-local-runner.py`.
  2. Check whether the `JOBS` dictionary already covers the newly identified pipeline/job.
  3. **If NOT covered or new steps were added:**
     - Extract all shell commands from the workflow YAML.
     - Strip all Docker wrappers (translate to native host commands).
     - Update the `JOBS` dictionary in `06-cicd-local-runner.py` to register the new job.
     - Save the runner script and verify syntax.
     - Log: `"Updated 06-cicd-local-runner.py to include newly discovered pipeline job '<name>'."`

- **Self-Loop Step 3 (Execute Runner & Establish Baseline Failures):**
  1. Run `python 03-ai-scripts/06-cicd-local-runner.py`.
  2. Capture the full output and exit code.
  3. If exit code = 0: All jobs pass! Proceed to Phase 3 (RCA & Final Completion).
  4. If exit code != 0: Zero in on the first specific failing job and its error output.

- **Self-Loop Step 4 (RCA & Zero In on the Specific Error):**
  1. For the zeroed-in failure, write a mandatory 4-part RCA file:
     - Path: `.lovable/memory/issues/xx-<slug>.md` (next sequential number)
     - Sections: **Why it happened / How it happened / Root Cause / Code Fix**
  2. Update `.lovable/01-index.md` and `.lovable/cicd-index.md`.
  3. Append any newly identified anti-pattern to `.lovable/strictly-avoid.md`.

- **Self-Loop Step 5 (Surgical Code Fix):**
  1. Open the specific offending source file and line identified in the RCA.
  2. Apply the minimal surgical fix cleanly.
  3. Run the guideline autofixer on modified files:
     ```text
     python 03-ai-scripts/05-guideline-autofixer.py <modified-files>
     ```

- **Self-Loop Step 6 (Re-Verify & Loop):**
  1. Re-run `python 03-ai-scripts/06-cicd-local-runner.py`.
  2. If the current error is fixed and other failures remain, self-loop to Step 4 to zero in on the next error.
  3. Continue looping until exit code = 0.

---

## Phase 1: Local Runner Script Generation (Steps 1 to PHASE_1_STEPS)

> [!IMPORTANT]
> **This entire phase is dedicated to one goal ONLY: creating `03-ai-scripts/06-cicd-local-runner.py`.**
> Do NOT attempt to fix code in this phase. Read, understand, and generate the script.
> If Image Input Handling (above) already updated the runner, verify the update is complete and move on.

### Step 1: Check for Existing Script & Force Override

- Check if `03-ai-scripts/06-cicd-local-runner.py` already exists.
- **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
- If it EXISTS and the user did **not** specify `force`: skip to Phase 2 immediately.
- If it is MISSING or the user said `force`: execute Steps 2 through PHASE_1_STEPS.

### Step 2: Deep CI/CD Configuration Scan

Spend up to PHASE_1_STEPS self-loop iterations reading the codebase in this order:

1. **Locate all CI/CD runner configuration files:**
   - GitHub Actions: `.github/workflows/*.yml` (scan ALL workflow files)
   - GitLab CI: `.gitlab-ci.yml`
   - Azure Pipelines: `azure-pipelines.yml`
   - Bitbucket Pipelines: `bitbucket-pipelines.yml`
   - CircleCI: `.circleci/config.yml`
   - Any custom runner scripts: `Makefile`, `03-ai-scripts/06-cicd-local-runner.py`, `run.sh`, `run.ps1`
   - Project lockfiles: `package-lock.json`, `go.sum`, `poetry.lock`, `bun.lockb` (to detect toolchain versions)
   - Language config: `.nvmrc`, `.python-version`, `go.mod`, `pyproject.toml`, `tsconfig.json`

2. **Extract all jobs and steps:** For every CI/CD job, record:
   - The `runs-on` image (e.g., `ubuntu-latest`, `node:20-alpine`)
   - All `run:` shell commands or `uses:` action steps
   - Environment variables set with `env:`
   - Cache keys used (e.g., `actions/cache`, `cache: npm`)
   - The dependency install command (e.g., `npm ci`, `go mod download`, `pip install -r requirements.txt`)
   - The lint command (e.g., `npm run lint`, `golangci-lint run`)
   - The typecheck/static analysis command (e.g., `tsc --noEmit`, `mypy .`)
   - The build command (e.g., `npm run build`, `go build ./...`)
   - The test command (e.g., `npm test`, `go test ./...`, `pytest`)

3. **Understand the local environment:** Check which tools are already installed on the host machine:
   - `node --version`, `npm --version`, `go version`, `python3 --version`, `cargo --version`, `php --version`
   - This determines which commands need translation vs. which can run as-is.

### Step 3: Docker Translation Rule (CRITICAL)

CI/CD pipelines run inside Docker containers. Your local runner MUST treat the **host machine as if it were that Docker image**. This means:

- **Strip all Docker invocations:** Never include `docker run`, `docker build`, `docker-compose`, or `docker pull` commands in the runner script.
- **Translate container steps to native host commands:**
  - A step like `docker run --rm node:20 npm ci` becomes `npm ci` (run directly on host).
  - A step like `docker run --rm python:3.12 pytest` becomes `python3 -m pytest` (run directly on host).
  - A step like `docker run --rm golang:1.22 go test ./...` becomes `go test ./...` (run directly on host).
- **Replace Docker-specific env injection** with Python `os.environ` assignments before running subprocess commands.
- **Skip container-only operations** like `docker login`, image tagging, or container registry pushes — these are deployment steps, not CI checks.

### Step 4: Write `06-cicd-local-runner.py` (Worker Pool & Log Aggregation Architecture)

Using all information gathered, generate `03-ai-scripts/06-cicd-local-runner.py` following these architectural requirements:

1. **Round-Robin Worker Process / Thread Pool Architecture:** Runs tasks (tests, linters, builds) concurrently using `concurrent.futures.ThreadPoolExecutor(max_workers=3)` (2–3 concurrent tasks).
2. **Enqueuing Announcement:** The script MUST announce upfront how many tasks it has enqueued across the worker pool (e.g. `[INFO] Enqueued 20 quality gates across 3 concurrent workers...`).
3. **Real-Time Progress & Timing:** Prints job completions in real time with individual runtimes (e.g. `PASS [Job Name] (X.XXs)`).
4. **Graceful Non-Cancelling Failure Handling:** If one job fails in a running batch, the runner DOES NOT abort or cancel other active workers. It lets running tasks finish gracefully, capturing all stdout and stderr.
5. **Consolidated Summary Report & Full Diagnostic Logs:** Prints a complete final summary with total executed, passed, failed, and timeouts. For every failed or timed-out job, it outputs the full command line, return code, stdout, and stderr so the agent has 100% complete RCA context.
6. **Clean Exit Code:** Exits with code 0 if all quality gates pass; exits with non-zero (1) if any job fails.

**Standard generated runner template:**

```python
#!/usr/bin/env python3
"""Auto-generated CI/CD local runner with concurrent worker pool and log aggregation.
Do not edit manually. Re-generate by running:
python 03-ai-scripts/06-cicd-local-runner.py --rebuild
"""
from concurrent.futures import ThreadPoolExecutor, as_completed, TimeoutError as FuturesTimeout
import os
import subprocess
import sys
import time

# ── Configurable Variables ──────────────────────────────────────────────────
BATCH_SIZE      = 3    # Number of jobs to run concurrently (round-robin worker pool)
JOB_TIMEOUT_SEC = 300  # Maximum seconds before a single job is timed out

# ── Environment Configuration ───────────────────────────────────────────────
os.environ.setdefault("CI", "true")
os.environ.setdefault("NODE_ENV", "test")

# ── Job Definitions (extracted from CI/CD workflow steps) ───────────────────
JOBS = {
    "install":   ["npm", "ci"],
    "lint":      ["npm", "run", "lint"],
    "typecheck": ["npx", "tsc", "--noEmit"],
    "build":     ["npm", "run", "build"],
    "test":      ["npm", "test", "--", "--watchAll=false"],
}

def run_job(name, cmd):
    start = time.monotonic()
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=JOB_TIMEOUT_SEC)
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd, result.returncode, result.stdout, result.stderr, elapsed
    except subprocess.TimeoutExpired as e:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd, "timeout", e.stdout or "", f"Job timed out after {JOB_TIMEOUT_SEC}s", elapsed
    except Exception as e:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd, 1, "", str(e), elapsed

def main():
    job_items = list(JOBS.items())
    total_jobs = len(job_items)
    print(f"[INFO] Enqueued {total_jobs} quality gates across {BATCH_SIZE} concurrent workers...\n")

    all_results = {}
    total_start = time.monotonic()

    # Execute jobs using concurrent worker pool
    with ThreadPoolExecutor(max_workers=BATCH_SIZE) as executor:
        futures = {executor.submit(run_job, name, cmd): name for name, cmd in job_items}
        for future in as_completed(futures):
            try:
                name, cmd, code, out, err, elapsed = future.result()
                all_results[name] = (code, out, err, elapsed, cmd)
                if code == 0:
                    print(f"  PASS [{name}] ({elapsed}s)")
                elif code == "timeout":
                    print(f"  TIMEOUT [{name}] ({elapsed}s)")
                else:
                    print(f"  FAIL [{name}] ({elapsed}s)")
            except Exception as ex:
                job_name = futures[future]
                all_results[job_name] = (1, "", str(ex), 0, JOBS.get(job_name, []))
                print(f"  FAIL [{job_name}] (Exception: {ex})")

    total_elapsed = round(time.monotonic() - total_start, 2)

    # ── Final Consolidated Summary Report ──────────────────────────────────
    print("\n" + "=" * 60)
    print("           CI/CD EXECUTION SUMMARY REPORT")
    print("=" * 60)

    passed_jobs = []
    failed_jobs = []
    timeout_jobs = []

    for name, (code, out, err, elapsed, cmd) in all_results.items():
        if code == 0:
            passed_jobs.append((name, elapsed))
        elif code == "timeout":
            timeout_jobs.append((name, elapsed, err, cmd))
        else:
            failed_jobs.append((name, elapsed, out, err, cmd))

    print(f"Total: {total_jobs} | Passed: {len(passed_jobs)} | Failed: {len(failed_jobs)} | Timeouts: {len(timeout_jobs)} | Time: {total_elapsed}s\n")

    if failed_jobs or timeout_jobs:
        print("Detailed Failure Logs:")
        print("-" * 60)
        for name, elapsed, out, err, cmd in failed_jobs:
            print(f"\n[FAILURE LOG] Job: {name} (Duration: {elapsed}s)")
            print(f"Command: {' '.join(cmd)}")
            if out.strip():
                print(f"Stdout:\n{out.strip()}")
            if err.strip():
                print(f"Stderr:\n{err.strip()}")
            print("-" * 60)

        for name, elapsed, err, cmd in timeout_jobs:
            print(f"\n[TIMEOUT LOG] Job: {name} (Duration: {elapsed}s)")
            print(f"Command: {' '.join(cmd)}")
            print(f"Reason: {err}")
            print("-" * 60)

        print(f"\n[FAILURE] CI/CD quality gates failed with {len(failed_jobs) + len(timeout_jobs)} error(s).")
        sys.exit(1)
    else:
        print(f"\n[SUCCESS] All {total_jobs} CI/CD quality gates passed (exit 0)!")
        sys.exit(0)

if __name__ == "__main__":
    main()
```

---

## Phase 2: Autonomous Fix Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> After generating the runner script, you enter an autonomous fix loop.
> Never stop to tell the user "here are the errors". Fix them immediately.
> Never await remote CI/CD results. The local runner IS your CI/CD.

Execute this exact loop for up to PHASE_2_STEPS iterations:

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Run:  python 03-ai-scripts/06-cicd-local-runner.py
    2. Capture exit_code and full output

    IF exit_code == 0:
        BREAK  ← All checks pass. Proceed to End of Tunnel.

    ELSE:
        3. Parse the failure output: identify the exact failing job, error message, file, and line.
        4. Document 4-part RCA in .lovable/memory/issues/xx-<slug>.md
        5. Apply the minimal surgical code fix to the codebase.
        6. Run the guideline autofixer on modified files:
              python 03-ai-scripts/05-guideline-autofixer.py <modified-files>
        7. Loop back to step 1 immediately. DO NOT stop.

IF STEP >= PHASE_2_STEPS AND exit_code != 0:
    Report the remaining failures clearly in chat and ask the user for guidance.
```

### Self-Looping Agent Rules

- **Each loop iteration = one AI tool call cycle.** After each iteration, self-loop by ending your turn and immediately starting the next iteration.
- **No output clutter:** Only print progress when a job status changes (new failure or new fix). Do not reprint passing jobs every iteration.
- **Spawn sub-agents for independent failures:** If multiple unrelated jobs are failing (e.g., lint + test + typecheck all failing independently), spawn one dedicated sub-agent per failure. Give each sub-agent a single-file bounding box. The parent agent collects results and re-runs the runner.

---

## Parallel Batch Execution Rules (Mandatory)

Every time `06-cicd-local-runner.py` is executed, it MUST run jobs in parallel batches, not one at a time:

- `BATCH_SIZE = 3` means at most 3 jobs run simultaneously within a batch.
- After each batch completes, wait for all futures to resolve before starting the next batch.
- Jobs within a batch that are NOT order-dependent (e.g., lint, typecheck, and unit tests) run concurrently.
- Jobs that ARE order-dependent (e.g., install → build → test) must be placed in **sequential batches**, not the same batch.
- After every batch, inspect the results immediately. Do not wait for the entire run to finish before reading output.

---

## Timeout Detection and Auto-Increase (Mandatory)

If any job is marked ⏱ TIMEOUT in the runner output:

1. **Do NOT treat it as a code failure.** Timeout ≠ broken code; it means the job took longer than `JOB_TIMEOUT_SEC` allows.
2. **Diagnose the cause:**
   - Is the job doing work proportional to a large input set (e.g., full test suite, large build)?
   - Is a network dependency involved (e.g., `npm ci` downloading packages)?
   - Did a previous step leave a hung process or lock file?
3. **Increase the timeout intelligently:**
   - If the job's measured elapsed time (before it was killed) was close to the limit: increase `JOB_TIMEOUT_SEC` by `50%`.
   - If the job appears to hang indefinitely (no output for >30s): look for a subprocess deadlock or missing stdin; fix the subprocess call, do not just increase the timer.
   - Open `06-cicd-local-runner.py`, update `JOB_TIMEOUT_SEC` to the new value, and save the file.
4. Re-run the runner immediately after the timeout fix. The timeout adjustment counts as one Phase 2 loop step.
5. Document the timeout increase in `.lovable/cicd-issues/` as a standard CI/CD issue entry.

---

## Error Enqueuing — Plan Task & CI/CD Issue (Mandatory on Every Failure)

Every time the runner reports a ❌ FAIL or ⏱ TIMEOUT, you MUST do **both** of the following before applying the code fix:

### A. Enqueue into Plan Tasks

Create (or append to) a pending plan task file at:

```text
.lovable/plans/pending/XX-cicd-<slug>.md
```

Where `XX` is the next available sequential number and `<slug>` is a short kebab-case description of the failure (e.g., `03-cicd-lint-unused-import`).

The plan task file MUST contain:

```markdown
# CI/CD Task: <short failure description>

## Source
- Runner job: <job-name>
- Error type: FAIL | TIMEOUT
- Detected at: <timestamp>

## Error Summary
<paste the exact error message or timeout log here>

## Required Fix
<one-sentence description of the fix needed>

## Acceptance Criteria
- [ ] `06-cicd-local-runner.py` reports ✅ PASS for job `<job-name>`
- [ ] No regression in any other job

## Status
- [ ] pending
```

Update `.lovable/plans/01-index.md` to register the new task entry immediately.

### B. Record in CI/CD Issues

Create a CI/CD issue record at:

```text
.lovable/cicd-issues/xx-<slug>.md
```

The CI/CD issue file MUST contain:

```markdown
# CI/CD Issue: <short failure description>

- Job: <job-name>
- Type: FAIL | TIMEOUT
- Detected: <timestamp>
- Status: open | resolved

## Error
<exact error output>

## Root Cause
<one-sentence root cause>

## Fix Applied
<what was changed to fix it>

## Plan Task
Enqueued at `.lovable/plans/pending/XX-cicd-<slug>.md`
```

Update `.lovable/cicd-index.md` in the same operation. Never delete existing entries.

---

## Phase 3: 4-Part RCA Requirement (Mandatory Memory File)

For each distinct failure type encountered in Phase 2, write a memory file at `.lovable/memory/issues/xx-<slug>.md` with exactly four sections:

1. **Why it happened:** High-level architectural reason for the failure.
2. **How it happened:** Exact execution flow that triggered the error.
3. **Root Cause:** Exact file, line number, and dependency responsible.
4. **Code Fix:** Code snippet showing the before and after of the fix.

Also append any new forbidden patterns to `.lovable/strictly-avoid.md`.

---

## Pre-Flight: Past RCA Ingestion

Before entering Phase 1 or Phase 2:

- Read `.lovable/cicd-issues/` and `.lovable/strictly-avoid.md`.
- Read the provided CI/CD error log if the user supplied one.
- These provide the known failure history and ensure you never repeat a past mistake.

---

## STRICT AVOIDANCE: Never Disable CLI Linting, Static Analysis, or CI/CD Checks (No Shortcut Cheating)

> [!CAUTION]
> **TOTAL BAN ON DISABLING, SKIPPING, OR BYPASSING CLI LINTERS AND CI/CD GATES:**
>
> - **NEVER** disable, comment out, delete, or skip any CLI linting command (`golangci-lint`, `eslint`, `markdownlint`, `tsc`, `pytest`, `phpstan`, `mypy`, `check-*.py`), build step, or test suite.
> - **NEVER** add `|| true`, `continue-on-error: true`, `# nolint`, `// eslint-disable`, or ignore flags to "quickly win the race" or fake a pipeline pass.
> - **Your job is to legitimately fix the underlying source code.** If resolving complex lint errors or test failures requires multiple sub-steps, sub-agents, or nested self-looping turns, you MUST execute all necessary turns until the code is 100% clean and compliant.
> - Disabling or bypassing any CI/CD or CLI lint check is an automatic and immediate rejection.

---

## Strictly Avoid: No Automatic Releases & Run All Tests (Strict Policy)

> [!CAUTION]
> **NO AUTOMATIC RELEASES:** This is a development fix workflow. You MUST NOT bump versions, update changelogs, or cut a release unless the user explicitly commands it in chat (e.g., "cut a release" or "bump the version"). All commits use `fix(ci): <description>`.
> **RUN ALL TESTS & QUALITY GATES (FULL PIPELINE FIDELITY):** `ci-cd-fix` is an explicit CI pipeline repair workflow. It MUST run all CI/CD quality gates, linters, static analysis, and test suites properly (`python 03-ai-scripts/06-cicd-local-runner.py`). NEVER skip or disable tests with `--no-tests` in CI/CD fix workflows; all unit test suites, AST checks, and verification jobs must be executed and diagnosed until they legitimately pass (`exit 0`).
> **ATOMIC CHANGE TRACKING:** Record all modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.

---

## Non-Negotiable Coding Standards

Run `03-ai-scripts/05-guideline-autofixer.py` on all modified files, then verify:

- [ ] **No Disabling CLI Linting (Zero Bypassing):** All CLI linters and CI/CD quality gates executed fully without `|| true`, `continue-on-error`, or suppression comments. Code was legitimately fixed.
- [ ] **Legitimate Multi-Step Self-Looping:** If complex errors occurred, I performed dedicated, single-step self-loop iterations to resolve each underlying failure instead of taking shortcuts.
- [ ] **Return New Line (R13-R16):** One blank line before every `return`/`throw` (unless sole statement). One blank line after closing `}`. Never two blank lines in a row.
- [ ] **No Explicit True Checks:** NEVER write `== true` or `=== true`. Write `if isReady`, not `if isReady == true`.
- [ ] **No Mixed Polarity:** NEVER write `if isA && !isB`. Extract to a named boolean.
- [ ] **Boolean Prefixes:** All boolean variables start with is and has only (can, should, was, etc. are banned), `was`, `will`, `did`, or `must`.
- [ ] **Error Handling:** No swallowed errors. Every propagated error is wrapped with `apperror.Wrap(err, "opName", ctx)`.
- [ ] **Strict Lowercase Files:** All generated/modified files use strictly lowercase naming.
- [ ] **Go Generate Sync:** If Go constants, enums, or stringers were modified, run `go generate ./...` and commit generated files.

---

## End of Tunnel Checklist

When `06-cicd-local-runner.py` exits with code 0:

- [ ] **Zero Linting/CI/CD Bypass:** Confirmed that NO CLI linters, static analysis tools, or test scripts were disabled, commented out, skipped, or bypassed with `|| true`.
- [ ] **Local CI Runner 100% Green:** All jobs in `06-cicd-local-runner.py` (including all unit tests, linters, and quality gates) passed legitimately (exit code = 0).
- [ ] **RCA Documented:** All encountered failures have memory files in `.lovable/memory/issues/`.
- [ ] **Antigravity Skill Updated:** Verified `.agents/skills/ci-cd-fix/skill.md` is present and synchronized with the latest rules.
- [ ] **Stage & Commit:** Group all related fixes into a single descriptive commit: `fix(ci): resolve <summary>`.
- [ ] **Push to Remote:** Push the commit to the current branch.
- [ ] **File Change Summary (MANDATORY):** In chat, list every file changed, what specifically changed inside it, and why. This summary is critical.

---

## Metadata

- slug: cicd-fix
- status: active
