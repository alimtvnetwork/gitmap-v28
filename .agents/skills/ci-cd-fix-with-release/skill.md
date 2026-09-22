---
name: ci-cd-fix-with-release
description: >-
  Use this skill to autonomously diagnose and fix CI/CD pipelines using local runner scripts, GitMap pipeline AI with dynamic waiting, 4-part RCA, and execute full release ceremony.
---

# Instruction (must follow): Release-Triggered CI/CD Fix Loop & Release Orchestration

Trigger Keywords & Aliases: `fix and release`, `ci release`, `fix CI/CD and release`, `cicd fix release`

/goal Autonomously diagnose and repair CI/CD pipeline issues, execute full unit test suites and quality gates to verify green status, and perform complete automated release publication with version bump and changelog sync.

/learn Ingest recent RCAs from `.ai-memory/cicd-issues/`, `.ai-memory/issues/`, `02-spec/02-coding-guidelines/02-canonical-size-tier.md`, `02-spec/02-coding-guidelines/01-cross-language/01-index.md`, `02-spec/02-coding-guidelines/01-cross-language/01-index.md`, and `02-spec/03-error-manage/` before touching any code so past mistakes are never repeated.

---

## Variables — Configurable at Runtime

```text
N = 200  (Total self-loop steps budget. The user may override this when triggering the prompt.)

PHASE_1_STEPS = N / 2  (Steps 1 .. N/2: Screenshot / GitMap Pipeline Discovery, Update 06-cicd-local-runner.py, Register New JOBS)
PHASE_2_STEPS = N / 2  (Steps N/2+1 .. N: Singly-Done Self-Loop Fixing, Zero in on Errors, 4-Part RCA, Green Gate Verification)
```

Both N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them.

> [!CAUTION]
> **SMART TARGETED TEST & RELEASE AUTHORITY:**
> This skill (`ci-cd-fix-with-release`) authorizes targeted test execution strictly on failed or modified packages to achieve the fastest green exit and release. The AI MUST execute tests in the smartest way possible:
> 1. **Priority Incremental Runner:** Execute smart incremental Go tests and gates via `python 03-ai-scripts/06-cicd-local-runner.py run-smart` (or alias `smart`, `--smart`, `-s`), which inspects Git changed files, builds ONLY changed packages into OS temp, and runs the Quad Runner.
> 2. **Specific Package Targeting:** Run/build ONLY packages and test functions directly cited in the failure stack trace: `python 03-ai-scripts/06-cicd-local-runner.py --pkg <target_package_or_file>`.
> 3. **Heatmap & Fast-Path Testing:** Use `--fast` to run only hot and warm tests based on `.ai-memory/test-heatmap.json`, skipping cold tests (`python 03-ai-scripts/06-cicd-local-runner.py --fast`).
> 4. **Changed Packages from Last Git Hash:** Compare against the last known git hash (`git diff --name-only HEAD~1`) and isolate packages that actually changed using `python 03-ai-scripts/06-cicd-local-runner.py --changed-only`.
> 5. **Bounded Stack Trace Extraction (RCA 58):** Always extract bounded failure frames (strictly 5 preceding + 20 trailing lines) halting on exit code or job boundary (`gitmap pe`) to prevent clipboard and report bloat.
> 6. **File State & Hash Tracking:** Every fix must persist modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).
> 7. **In-Flight Heartbeats & ETA Wait:** The local runner emits heartbeats every 25s (`--heartbeat-interval 25.0`) and writes status to `.ai-memory/temp/runner-eta.json`. Agents must sleep for 60s or remaining ETA rather than busy-polling.
> 8. **Strict Ban on Extraneous Runs:** NEVER run the entire test suite, spellcheckers, or unrelated packages. Verify strictly using `python 03-ai-scripts/06-cicd-local-runner.py run-smart`, `--changed-only`, or `--pkg <target>`, then proceed immediately to the automated release ceremony.

---

## Strict In-Repository Execution & `.ai-memory/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.ai-memory/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.ai-memory/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.ai-memory/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`).
>    - RCA & Issue Logs: `.ai-memory/memory/issues/` and `.ai-memory/cicd-issues/`.
>    - Execution Plans & Subtasks: `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.ai-memory/coding-guidelines.md`.
> 3. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.ai-memory/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `02-spec/03-error-manage/01-index.md`, `.ai-memory/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands]`02-spec/13-generic-cli/01-index.md)`
> 4. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.
> 5. **Cross-Platform Python CI Mandate (TOTAL BAN on new `.sh` scripts in CI):** All newly created or refactored CI/CD verification tools, determinism checks, fixtures, and linter jobs MUST be implemented in pure, cross-platform Python (`.py`). Legacy `.sh` scripts must be converted to `.py` scripts so all pipelines run natively across Linux, macOS, and Windows without relying on bash emulation.
> 6. **Temp & Failure Folder Isolation:** All temporary directories, runner caches, and test artifacts MUST be strictly placed in `.ai-memory/temp/`. Creating `.tmp/` at the repository root or outside `.ai-memory/` is strictly forbidden.
>    - Dedicated Failure Directory: `.ai-memory/temp/failures/` is the dedicated folder where failed tests and failed quality gates write error logs (`<test-or-job-name>.log`).
>    - Passing Tests Completely Silent: Passing tests must produce ZERO filesystem artifacts (zero files written) and remain completely silent in output logs.
> 7. **Runner In-Flight ETA Wait & GitMap Dynamic Waiting Protocol:**
>    - **Local Runner:** The runner dynamically writes live status and remaining ETA to `.ai-memory/temp/runner-eta.json` (emitting in-flight heartbeats strictly every 25 seconds or more). If an agent inspects an active background job and it is still running, the agent MUST sleep/wait for **1 minute (60 seconds) each time**, or dynamically sleep for the remaining ETA duration read from `.ai-memory/temp/runner-eta.json` instead of busy-polling or querying in loops.
>    - **Remote CI/CD Pipelines (GitMap Mandate):** When checking remote pipeline workflows, the agent MUST use `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`). To inspect running pipelines without burning tokens or user credits, use `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`). Tight-loop polling (e.g. `gh run view` in rapid loops) is STRICTLY BANNED. Adaptive sleep intervals: ETA > 120s wait 20s-30s; 60s < ETA <= 120s wait 10s-20s; ETA <= 60s wait 5s-10s.
>    - **Targeted Failure Extraction:** Extract targeted failure lines (`##[error]`, `FAIL:`, compile errors) from GitMap output directly into 4-part RCA files without reading noisy passing step logs.
> 8. **Centralized Test Inventory & Incremental Caching:** All unit tests are cataloged in `.ai-memory/test-inventory.json` with strictly repository-relative paths (`target_file`, `test_file`). Modified files MUST be recorded safely under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`. First run executes all tests to establish baseline timings; subsequent runs execute incrementally only if target code files or test files change. Slow test threshold defaults to `4.0s` (configurable via `GITMAP_SLOW_TEST_THRESHOLD`).
> 9. **Dual-Queue Worker Pools:** Slow tests run in a dedicated 4-worker pool running at most 2 tests at a time per batch. Fast tests run in a 4-worker pool running at most 4 tests at a time, pulling in chunks of 100 tests from the test inventory queue until all are complete.
> 10. **Dynamic ETA Sleep Protocol:** The AI agent reads `.ai-memory/temp/runner-eta.json` (or `gitmap pipeline-ai status --json`), sleeps for the estimated wait time rather than looping, and if still active upon waking, re-checks remaining ETA and sleeps again to avoid burning tokens.
> 11. **Zero-Storage GitHub Actions Mandate (Total Ban on CI Artifact Uploads):** Workflows MUST NOT upload test outputs, coverage files, Playwright reports, or drift summaries via `actions/upload-artifact`. Free-tier accounts have a strict 0.5 GB shared quota across all repositories. All reports, failures, and summaries MUST be emitted directly to `$GITHUB_STEP_SUMMARY`, console stdout (`cat log.txt`), or sticky PR comments with zero storage consumption. Only true GitHub release assets (binaries/tarballs on tagged releases) are permitted.

---

## Fast File Discovery & Diagnostic Toolchain (Mandatory Acceleration)

To rapidly locate failing pipeline definitions, broken source files, test fixtures, and error logs without hitting 50-result tool caps, the AI agent MUST utilize the diagnostic toolchain:
- **Remote Pipeline AI Status (<50ms):** `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`)
- **Remote Dynamic Timeout Wait:** `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`)
- **Extract Failing Step Error Logs:** `gitmap pipeline error-logs` (or alias `gitmap pe`, clear with `gitmap pe clear -y`)
- **Pipeline Runner Targets & Cache Table:** `gitmap pipeline details` (or alias `gitmap pd`)
- **Scan Source & Test Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts,py --limit 100 --stats`
- **Fast Cached Pattern Search (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<error-or-symbol>" --limit 50`
- **Sub-Millisecond Folder Listing & Reader:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .github/workflows --limit 20`
- **Read Workflow or Log File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file .github/workflows/ci.yml`
- **Codebase Topology Overview:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`
- **Record Modified Files Under Lock:** `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`

> [!NOTE]
> **CI/CD Fix Verification Allowance:** Unlike routine refactoring turns, CI/CD Fix workflows ARE explicitly authorized and expected to run local builds, tests, and runner scripts (`python 03-ai-scripts/06-cicd-local-runner.py`) to diagnose failures, reproduce errors, and verify that all quality gates pass (exit code 0).

---

## Screenshot & Pipeline Discovery Protocol (Execute First When Any Image Is Provided)

> [!IMPORTANT]
> **If the user provides any image or screenshot showing a CI/CD pipeline name, failing workflow, or error log:**
>
> 1. **FIRST ACTION — Update Python Runner:** Locate the pipeline/job in `.github/workflows/*.yml` (or repo CI configs) to find whatever new jobs, steps, or linters were added, and **immediately update `03-ai-scripts/06-cicd-local-runner.py`** to include them in the `JOBS` dictionary.
> 2. **SECOND ACTION — Singly-Done Self-Loop Execution:** Run the Python script iteratively, zeroing in on one failing error at a time using strictly bounded self-loop turns until all checks exit with code 0 (`exit 0`).

### Bounded Single-Step Self-Loop Sequence (Singly Done — No Overloaded Steps)

Every step must be **singly done** using bounded self-looping turns:

- **Self-Loop Step 1 (Extract Pipeline Name from Screenshot or GitMap):**
  1. Read image or query `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status -t <etaSeconds>`) to extract the pipeline name, failing job name, and targeted failure diagnostics (`##[error]`, `FAIL:`, compile errors).
  2. Extract failing step logs via `gitmap pipeline error-logs` (`gitmap pe`) using bounded stack trace extraction (5 before + 20 after). Clear stale errors with `gitmap pe clear -y`. Inspect runner jobs via `gitmap pipeline details` (`gitmap pd`).
  3. Scan `.github/workflows/*.yml` to identify the corresponding shell commands and dependencies.

- **Self-Loop Step 2 (FIRST ACTION: Update Python Runner Script):**
  1. Open `03-ai-scripts/06-cicd-local-runner.py`.
  2. If new jobs/steps are found in workflows, strip Docker wrappers and register the new commands in the `JOBS` dictionary.
  3. Save the runner script and verify syntax.

- **Self-Loop Step 3 (Execute Runner & Establish Baseline Failures):**
  1. Run `python 03-ai-scripts/06-cicd-local-runner.py run-smart` (or `python 03-ai-scripts/06-cicd-local-runner.py --changed-only` or `--pkg <target>`).
  2. Capture the full output and exit code.
  3. If exit code = 0: All jobs pass! Proceed to Phase 3 (Final Verification Gate) and Phase 4 (Release Publication).
  4. If exit code != 0: Zero in on the first specific failing job and its error output.

- **Self-Loop Step 4 (RCA & Zero In on Error):**
  1. Write 4-part RCA in `.ai-memory/memory/issues/xx-<slug>.md`.
  2. Register in `.ai-memory/01-index.md` and `.ai-memory/strictly-avoid.md`.

- **Self-Loop Step 5 (Surgical Code Fix):**
  1. Open the specific file and line, apply minimal surgical fix.
  2. Run `python 03-ai-scripts/05-guideline-autofixer.py <modified-files>`.
  3. Record modified files safely under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <modified-files>`.

- **Self-Loop Step 6 (Re-Verify & Loop):**
  1. Re-run `python 03-ai-scripts/06-cicd-local-runner.py run-smart` (or `--pkg <affected_pkg>`).
  2. If resolved and more errors remain, self-loop to Step 4 to zero in on the next error until exit code = 0.

- **Self-Loop Step 7 (Proceed to Release):**
  1. Once exit code is 0, proceed directly to Phase 3 (Final Verification Gate) and Phase 4 (Release Publication).

---

## Phase 1: Local Runner Script Generation (Steps 1 to PHASE_1_STEPS)

> [!IMPORTANT]
> **This phase is dedicated ONLY to creating `03-ai-scripts/06-cicd-local-runner.py`.**
> Do NOT fix code in this phase. Read, understand, and generate the script.

### Step 1: Check for Existing Script & Force Override

- Check if `03-ai-scripts/06-cicd-local-runner.py` already exists.
- **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
- If it EXISTS and the user did **not** say `force`: skip to Phase 2.
- If MISSING or `force` was requested: execute Steps 2–4 for up to PHASE_1_STEPS iterations.

### Step 2: Deep CI/CD Configuration Scan

Spend up to PHASE_1_STEPS self-loop iterations reading in this order:

1. **CI/CD configuration files:**
   - GitHub Actions: `.github/workflows/*.yml` (ALL files)
   - GitLab CI: `.gitlab-ci.yml`
   - Azure Pipelines: `azure-pipelines.yml`
   - Bitbucket: `bitbucket-pipelines.yml`
   - CircleCI: `.circleci/config.yml`
   - Custom runners: `Makefile`, `03-ai-scripts/06-cicd-local-runner.py`, `run.sh`, `run.ps1`
2. **Language configuration:** `.nvmrc`, `.python-version`, `go.mod`, `pyproject.toml`, `tsconfig.json`, lockfiles
3. **For every CI/CD job, record:**
   - `runs-on` image (e.g., `ubuntu-latest`, `node:20-alpine`)
   - All `run:` shell commands and `uses:` action steps
   - Environment variables from `env:` blocks
   - Dependency install commands (`npm ci`, `go mod download`, `pip install -r requirements.txt`)
   - Lint, typecheck, build, and test commands
4. **Check local toolchain:** Run `node --version`, `go version`, `python3 --version`, etc. to know what is available natively.

### Step 3: Docker Translation Rule (CRITICAL)

The host machine IS the Docker container. Strip all Docker wrappers:

- `docker run --rm node:20 npm ci` → `npm ci`
- `docker run --rm python:3.12 pytest` → `python3 -m pytest`
- `docker run --rm golang:1.22 go test ./...` → `go test ./...`
- Replace Docker `env` injection with Python `os.environ` assignments.
- **Skip entirely:** `docker login`, image tagging, container registry pushes — these are deployment steps, not CI checks.

### Step 4: Write `06-cicd-local-runner.py` (Worker Pool & Log Aggregation Architecture)

Generate `03-ai-scripts/06-cicd-local-runner.py` that:

1. **Round-Robin Worker Process / Thread Pool Architecture:** Runs tasks (tests, linters, builds) concurrently using `concurrent.futures.ThreadPoolExecutor(max_workers=3)` (2–3 concurrent tasks).
2. **Enqueuing Announcement:** The script MUST announce upfront how many tasks it has enqueued across the worker pool (e.g. `[INFO] Enqueued 20 quality gates across 3 concurrent workers...`).
3. **Real-Time Progress & Timing:** Prints job completions in real time with individual runtimes (e.g. `PASS [Job Name] (X.XXs)`).
4. **Graceful Non-canceling Failure Handling:** If one job fails in a running batch, the runner DOES NOT abort or cancel other active workers. It lets running tasks finish gracefully, capturing all stdout and stderr.
5. **Consolidated Summary Report & Full Diagnostic Logs:** Prints a complete final summary with total executed, passed, failed, and timeouts. For every failed or timed-out job, it outputs the full command line, return code, stdout, and stderr so the agent has 100% complete RCA context.
6. **Clean Exit Code:** Exits with code 0 only when ALL jobs pass; exits non-zero if ANY job fails.

**Template (adapt JOBS dict from actual CI/CD config):**

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
> Never stop to tell the user "here are the errors". Fix them. Never await remote CI/CD results.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Run: python 03-ai-scripts/06-cicd-local-runner.py
    2. Capture exit_code and full output.

    IF exit_code == 0:
        BREAK  ← All checks pass. Proceed to End of Tunnel.

    ELSE:
        3. Parse failure: identify exact failing job, error message, file, and line.
        4. Document 4-part RCA in .ai-memory/memory/issues/xx-<slug>.md
        5. Apply the minimal surgical code fix.
        6. Run: python 03-ai-scripts/05-guideline-autofixer.py <modified-files>
        7. Loop immediately to step 1. DO NOT stop.

IF STEP >= PHASE_2_STEPS AND exit_code != 0:
    Report remaining failures clearly in chat and ask the user for guidance.
    DO NOT proceed to Phase 3 or Phase 4 until all failures are resolved.
```

### Self-Looping Agent Rules

- Each loop iteration = one AI tool call cycle. End your turn and immediately start the next.
- Only print progress when a job status changes. Do not reprint passing jobs every iteration.
- Spawn one sub-agent per independent failure (lint / test / typecheck failing separately). Give each sub-agent a single-file bounding box.

---

## Parallel Batch Execution Rules (Mandatory)

- `BATCH_SIZE = 3`: at most 3 jobs run simultaneously per batch.
- After each batch completes, wait for all futures before starting the next batch.
- Non-order-dependent jobs (lint, typecheck, unit tests) run within the same batch.
- Order-dependent jobs (install → build → test) must be placed in sequential batches.
- After every batch, inspect results immediately. Do not wait for the full run.

---

## Timeout Detection and Auto-Increase (Mandatory)

If any job shows ⏱ TIMEOUT:

1. Do NOT treat it as a code failure. Timeout means the job is slow, not broken.
2. Diagnose: large input set? network dependency (`npm ci`)? hung subprocess or lock file?
3. Increase `JOB_TIMEOUT_SEC` intelligently:
   - Elapsed close to the limit → increase by 50%.
   - Job hangs with no output for >30s → fix the subprocess call, not the timer.
   - Open `06-cicd-local-runner.py`, update `JOB_TIMEOUT_SEC`, save, re-run.
4. Re-run the runner. The timeout adjustment counts as one Phase 2 loop step.
5. Document the timeout change in `.ai-memory/cicd-issues/`.

---

## Error Enqueuing — Plan Task & CI/CD Issue (Mandatory on Every Failure)

On every ❌ FAIL or ⏱ TIMEOUT, BEFORE applying any code fix, do both:

### A. Enqueue into Plan Tasks (`.ai-memory/plans/pending/XX-cicd-<slug>.md`)

```markdown
# CI/CD Task: <short failure description>

## Source
- Runner job: <job-name>
- Error type: FAIL | TIMEOUT
- Detected at: <timestamp>

## Error Summary
<exact error message>

## Required Fix
<one-sentence description>

## Acceptance Criteria
- [ ] `06-cicd-local-runner.py` reports ✅ PASS for job `<job-name>`
- [ ] No regression in any other job

## Status
- [ ] pending
```

Update `.ai-memory/plans/01-index.md` immediately.

### B. Record in CI/CD Issues (`.ai-memory/cicd-issues/xx-<slug>.md`)

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
<what was changed>

## Plan Task
Enqueued at `.ai-memory/plans/pending/XX-cicd-<slug>.md`
```

Update `.ai-memory/cicd-index.md` in the same operation. Never delete existing entries.

---

## Phase 3: Final Verification (Gate Before Release)

> [!IMPORTANT]
> Phase 3 is a hard gate. The release MUST NOT start until every item below is green.
> If any item fails, loop back to Phase 2 immediately.

- [ ] **Smart Targeted Test & CI/CD Verification (Before Release):** Run targeted verification via priority shortcuts (`python 03-ai-scripts/06-cicd-local-runner.py run-smart`, `--changed-only`, or `--pkg <affected_pkg>` with optional `--fast` heatmap filtering) covering all failing stack trace targets and packages changed since the last git hash. All modified/failing package quality gates MUST pass 100% green (`exit 0`). The release MUST NOT start if any targeted test fails.
- [ ] **Test Inventory Validation:** Check `.ai-memory/temp/recent-file-changes.json` against `.ai-memory/test-inventory.json` and verify all tests associated with modified files pass.
- [ ] **No open plan tasks from this run:** All `.ai-memory/plans/pending/XX-cicd-*.md` files created in this run are marked `resolved` or closed.
- [ ] **All RCA files written:** Every failure encountered has a `.ai-memory/memory/issues/xx-<slug>.md` with all 4 sections.
- [ ] **Coding standards pass:** Run `python 03-ai-scripts/05-guideline-autofixer.py` on all modified files. Zero violations remain.
- [ ] **Git working tree is clean:** Run `git status`. No untracked or unstaged files. Commit any remaining changes with `fix(ci): final pre-release fixes`.

---

## Phase 4: Release

> [!IMPORTANT]
> Only enter Phase 4 after Phase 3 passes 100%. No exceptions.

### Step R-1: Pre-Flight Checks

1. Run `git status`. If uncommitted changes remain: commit them now with `fix(ci): pre-release cleanup`.
2. Run `git pull` to merge any upstream changes.
3. Read the current version from `version.json` (canonical source). Print it.
4. Idempotency guard: if the canonical version already equals the computed new version, STOP. Someone half-ran a release. Detect what is done, resume from the first incomplete step. Do NOT double-bump.
5. Placeholder guard: if the previous version's changelog entry is empty or contains `TBD`/`WIP`, refuse to release until it is filled (or the user overrides).

### Step R-2: Bump the Version & Execute 5-Step Release Branching Lifecycle

**Primary path — use the release orchestrator or dedicated bump script:**

```bash
# Option A (Recommended): Full automated release orchestrator
python 03-ai-scripts/29-release-orchestrator.py --tier minor

# Option B: Dedicated bump script with full release lifecycle
python .ai-memory/release/bump_versions.py --type minor --create-release
# or: python 03-ai-scripts/37-bump-version.py --tier minor
```

All release execution MUST strictly enforce the **5-Step Release Branching Lifecycle**:
1. **Step 1: Create & Switch to Release Branch:**
   Create and switch to `release/vX.Y.Z` FIRST before modifying any version files (`git checkout -b release/vX.Y.Z`).
2. **Step 2: Bump Version on Release Branch via Python Script:**
   Execute the dedicated Python bump script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) on the release branch.
   - **CRITICAL REPOSITORY ADAPTATION & SCRIPT REPAIR:** Inspect the target repository architecture to identify where versions are defined (`version.json`, `package.json`, `Cargo.toml`, `pyproject.toml`, `go.mod`, etc.), where and how they will change (`readme.md`, `changelog.md`, install scripts), and post-bump synchronization commands (`npm run sync`, `go generate ./...`). If the bump script is missing, outdated, or lacks support for this repository's version pin sites, **the agent MUST fix or recreate the Python bump script immediately** before executing the release.
3. **Step 3: Commit in Release Branch:**
   Stage and commit all version bump and generated release files on the release branch (`release: vX.Y.Z <scope>`).
4. **Step 4: Create Annotated Git Tag:**
   Create the annotated tag on the release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
5. **Step 5: Put Commit Back to Main Branch & Push:**
   Switch to `main` (`git checkout main`), merge the release branch commit (`git merge release/vX.Y.Z`), push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`, then restore the starting branch if different from `main`.

---

### MANDATORY: Release Page Install One-Liners (FATAL IF MISSED ON GITHUB/GITLAB)

> [!CAUTION]
> **NEVER run `gh release create <tag> --generate-notes` ALONE.**
> Running `--generate-notes` without a structured `--notes-file` is a FATAL DEFECT: GitHub will only display commit hashes (as seen in broken release pages) and completely omits the installation one-liners!
>
> Every published GitHub / GitLab release page MUST have the **Quick Install One-Liners** prominently placed right at the top of the release body!

Before calling `gh release create` or `glab release create`, the release automation MUST assemble a release notes file (e.g. `.ai-memory/release/release-notes-vX.Y.Z.md` or `/tmp/release-body.md`) containing:

#### 1. Quick Install One-Liners by Project Type

**For Binary / Download Asset Repositories (e.g., Go/Rust/C CLI tools like `gitmap`):**

```markdown
## Quick Install vX.Y.Z

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/<owner>/<repo>/releases/download/vX.Y.Z/install.ps1 | iex
```

### Linux / macOS (Bash)

```bash
curl -fsSL https://github.com/<owner>/<repo>/releases/download/vX.Y.Z/install.sh | bash
```
```

**For Script / Meta-Repositories (e.g., `prompt-architect`):**

```markdown
## Quick Install vX.Y.Z

### Windows (PowerShell)
```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "vX.Y.Z"
```

### Unix / Bash

```bash
curl -sL https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.sh | bash -s -- ".ai-memory/prompts" "vX.Y.Z"
```
```

#### 2. Changelog Section for vX.Y.Z

Directly below the Quick Install block, include the extracted `[vX.Y.Z]` section from `changelog.md`.

#### 3. Platform Release Command

Pass the generated notes file via `--notes-file`:

```bash
gh release create "vX.Y.Z" --title "vX.Y.Z" --notes-file ".ai-memory/release/release-notes-vX.Y.Z.md" --generate-notes
```

*(Note: `--generate-notes` may be appended so GitHub adds commit logs below the install one-liners and changelog, but `--notes-file` is MANDATORY).*

---

**Fallback chain (if `.ai-memory/release/bump_versions.py` is missing):**

1. **Fallback 1:** Read `.ai-memory/release/release-method.md` to identify all version pin sites. Regenerate `bump_versions.py` from that documentation. Ensure it generates the release notes file with the Quick Install one-liners before running `gh release create`.
2. **Fallback 2:** If `release-method.md` is also missing, walk the repository with Python `os.walk` (ignoring `.git`, `node_modules`, `.venv`) to discover all version pin sites. Write `release-method.md` documenting them. Generate `bump_versions.py` with the correct `FILES_TO_BUMP`, release notes generator, `git checkout -b`, `git commit`, `git tag`, `git push`, and `gh release create ... --notes-file` logic. Run it.
3. **Fallback 3:** If discovery fails, stop and ask the user to specify the version pin sites explicitly.

> [!CAUTION]
> **NEVER use `rg`, `grep`, or `find` to globally search for version strings.** Follow the fallback chain above. Global searches on large repos are slow and error-prone.

### Step R-3: Pin the New Version in `readme.md`

Rewrite every occurrence of the previous version (`vX.Y.Z` and bare `X.Y.Z`) in badges, install snippets, and inline references. After this step, `grep "<previous-version>" readme.md` MUST return nothing.

### Step R-4: Write the Changelog Entry

Add the following block at the top of `changelog.md`, directly under `# Changelog`:

```markdown
## [vX.Y.Z] YYYY-MM-DD <short headline>

### Install <Project Name> vX.Y.Z

Unix/Bash:
`curl -sL https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.sh | bash -s -- ".ai-memory/prompts" "vX.Y.Z"`

PowerShell:
`Invoke-WebRequest -Uri https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "vX.Y.Z"`

### Added / Changed / Fixed / Removed

- <one bullet per real change, naming the exact file or behavior>

### Issues (only if any step failed)

- [xx-<slug>](.ai-memory/release/issues/XX-vX.Y.Z-<slug>.md) short description
```

Dynamically discover `<owner>/<repo>` by running `git config --get remote.origin.url`. Do NOT hardcode URLs.

### Step R-5: Final Verification After Bump

1. Run the version-sync check if one exists (`scripts/check-version-sync.*`, `scripts/verify-versions.*`). It MUST exit 0.
2. Verify ALL pin sites reference the new version. No previous-version strings outside the historic allow-list (`changelog.md`, `release_notes.md`, `.ai-memory/release/`, dated archives).
3. All markdown filenames in the repository MUST be strictly lowercase. Rename any uppercase files with `git mv` in the same turn.

### Step R-6: Issue Logging (If Anything Goes Wrong)

If any release step fails, write an issue file at:

```text
.ai-memory/release/issues/XX-vX.Y.Z-<slug>.md
```

Include: previous version, new version, step number and name, command run, full error output, files involved, resolution or `unresolved`. Link it from the `### Issues` bullet in the changelog entry.

---

## Phase 4 Release Checklist

- [ ] Phase 3 (Final Verification) passed with exit code 0 legitimately (no CLI linters or tests bypassed).
- [ ] Confirmed that NO CLI linting (`golangci-lint`, `eslint`, `markdownlint`, `tsc`, `pytest`), build steps, or test runs were skipped, commented out, or bypassed with `|| true`.
- [ ] Git working tree was clean before release steps.
- [ ] `git pull` completed with no conflicts.
- [ ] Previous and new versions both stated explicitly.
- [ ] Executed 5-step release branching lifecycle: `release/vX.Y.Z` created first, bumped via repository-aware Python bump script, committed on release branch, tagged `vX.Y.Z`, merged back into `main`, and pushed to remote.
- [ ] `python 03-ai-scripts/29-release-orchestrator.py` or `.ai-memory/release/bump_versions.py` or `03-ai-scripts/37-bump-version.py` ran successfully.
- [ ] All version pin sites updated to the new version.
- [ ] `readme.md` pinned to new version. No previous version strings remain.
- [ ] Changelog entry added with real bullets. No `TBD` or empty entries.
- [ ] All markdown filenames in repo are strictly lowercase.
- [ ] `### Issues` block present in changelog if any step failed, with links.
- [ ] Release notes file generated containing Quick Install One-Liners (PowerShell & Bash) and changelog.
- [ ] Release branch, release tag `vX.Y.Z`, and merged `main` branch pushed to remote.
- [ ] GitHub/GitLab release created via `gh release create --notes-file` or `glab release create --notes-file` (NEVER bare `--generate-notes`).
- [ ] Release description on GitHub/GitLab verified to contain the Quick Install one-liners, NOT just raw commit hashes.
- [ ] Verified `.agents/skills/ci-cd-fix-with-release/skill.md` and `.agents/skills/ci-cd-fix/skill.md` are present and synchronized with the latest rules.
- [ ] Report posted in chat: previous version, new version, bump tier, exact files changed.

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

## Mandatory Pre-Release Full Unit Tests & CI/CD Verification (Strict Policy)

> [!CAUTION]
> **ALL CI/CD FIXES & RELEASES MUST RUN TESTS PROPERLY:** This workflow repairs CI/CD pipelines and cuts a release. You MUST NOT skip or disable tests with `--no-tests`. All unit test suites, integration tests, AST checks, and quality gates MUST be executed and verified green (`python 03-ai-scripts/06-cicd-local-runner.py --run-tests`).
> **ATOMIC CHANGE TRACKING:** Record all modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.

---

## Non-Negotiable Coding Standards

- [ ] **GitMap Pipeline-AI & Dynamic Waiting:** If inspecting remote pipeline status, used `gitmap pipeline-ai status --json` / `gitmap pl-ai status -t <sec>` with adaptive sleep to prevent credit waste, and extracted targeted diagnostics (`##[error]`, `FAIL:`, compile errors).
- [ ] **Atomic Change Tracking:** All modified files were recorded safely under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
- [ ] **No Disabling CLI Linting (Zero Bypassing):** All CLI linters and CI/CD quality gates executed fully without `|| true`, `continue-on-error`, or suppression comments. Code was legitimately fixed.
- [ ] **Zero Actions Storage (Total Ban on CI Artifacts):** Confirmed that NO `actions/upload-artifact` steps exist in CI workflows; all diagnostic outputs stream to `$GITHUB_STEP_SUMMARY` or console logs.
- [ ] **Legitimate Multi-Step Self-Looping:** If complex errors occurred, I performed dedicated, single-step self-loop iterations to resolve each underlying failure instead of taking shortcuts.
- [ ] **Return New Line (R13-R16):** Blank line before `return`/`throw` (unless sole statement). Blank line after `}`. Never two blank lines in a row.
- [ ] **No Explicit True Checks:** Never `== true`. Write `if isReady`.
- [ ] **No Mixed Polarity:** Never `if isA && !isB`. Extract to a named boolean.
- [ ] **Boolean Prefixes & IsDefined:** All booleans start with `is` or `has` only (all other prefixes banned). MANDATORY: Use `isDefined` (or `res.IsDefined()`) instead of inverted empty checks (`!isEmpty` / `!res.IsEmpty()`). Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`.
- [ ] **Error Handling:** No swallowed errors. Wrap with `apperror.Wrap(err, "opName", ctx)`.
- [ ] **Strict Lowercase Files:** All generated/modified files use strictly lowercase naming.
- [ ] **Go Generate Sync:** If Go constants, enums, or stringers were modified, run `go generate ./...` and commit generated files.

---

## Metadata

- slug: cicd-fix-with-release
- status: active

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the conclusion of Phase 4 Release, after version manifests, changelog entries, and release notes are updated, you MUST stage everything (`git add -A`), commit with clean semantic messages, and push directly to the remote repository. Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them. All modified files across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

