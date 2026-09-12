# Pipeline Architecture & Cross-Platform Python Automation — Workflow (must follow)

Trigger Keywords & Aliases: `create-ci-cd`, `cicd-create`, `create cicd`, `setup cicd pipeline`, `build ci-cd pipeline`, `cicd create python`, `05-cicd-pipeline-create`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 300
```

N = total self-loop steps budget for end-to-end CI/CD creation and cross-platform automation.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase to inventory all architectural violations and anti-patterns.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py` with exit code 0 (`exit 0`).
9. [ ] /learn Ingest `spec/12-cicd-pipeline-workflows/` for domain-specific architectural specifications.
10. [ ] /learn Ingest `spec/02-coding-guidelines/06-cicd-integration/` for domain-specific architectural specifications.
11. [ ] /learn Ingest `spec/11-powershell-integration/` for domain-specific architectural specifications.
12. [ ] /learn Ingest `spec/14-update/` for domain-specific architectural specifications.
13. [ ] /learn Ingest `spec/15-distribution-and-runner/` for domain-specific architectural specifications.
14. [ ] /learn Ingest `spec/16-generic-release/` for domain-specific architectural specifications.
15. [ ] /learn Ingest `spec/17-consolidated-guidelines/18-cicd-pipeline-workflows.md` for domain-specific architectural specifications.
16. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. 150: Spec Ingestion, Python Automation Design, Local Runner Setup)
PHASE_2_STEPS = N / 2   (Steps 151 .. 300: Workflow Generation, Python Linter Verification, Green Gate Loop)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before any execution, check if this prompt is installed as a native Antigravity Skill.

1. Check if `.agents/skills/ci-cd-create/skill.md` exists in the workspace.
2. If it does NOT exist, create it now. Write the core instructions to `.agents/skills/ci-cd-create/skill.md` with frontmatter:
   ```yaml
   ---
   name: ci-cd-create
   description: >-
     Use this skill to autonomously design, create, and verify cross-platform Python CI/CD pipelines, linters, and workflows.
   ---
   ```
3. Once installed, load it on-demand via progressive disclosure for all future runs.

---

## 1. Specification Ingestion Checklist (Read Index / Overview and Subfiles)

For every specification folder below, AI agents MUST read `index.md` or `00-overview.md` (or both) first, then deeply internalize all referenced subfiles before authoring any CI/CD or automation files:

### A. CI/CD Pipeline Workflows (`spec/12-cicd-pipeline-workflows/`)

- [ ] Read `00-overview.md` or `index.md`
- [ ] `02-ci-pipeline.md`: Core CI triggers, matrix builds, artifact passing, and caching contracts.
- [ ] `03-shared-conventions.md`: Universal environment variables, exit codes, and cross-platform paths.
- [ ] `05-release-pipeline.md` & `04-github-release-standard.md`: Automated release tagging, checksums, and assets.
- [ ] `01-browser-extension-deploy/`: (`00-overview.md`, `02-ci-pipeline.md`, `05-release-pipeline.md`, `99-consistency-report.md`).
- [ ] `02-go-binary-deploy/`: (`00-overview.md`, `02-ci-pipeline.md`, `05-release-pipeline.md`, `03-complete-workflow-reference.md`).
- [ ] `03-reusable-ci-guards/`: (`00-overview.md`, `01-forbidden-name-guard.md`, `02-grandfather-baseline-naming.md`, `03-cross-file-collision-audit.md`, `04-baseline-diff-lint-gate.md`, `05-actionable-lint-suggestions.md`, `06-matrix-test-aggregator.md`, `07-shared-cli-wrapper.md`, `08-config-schema.md`, `09-workflow-templates.md`, `99-ai-implementation-guide.md`).
- [ ] `06-vulnerability-scanning.md`: Trivy, Snyk, and dependency vulnerability audits.
- [ ] `07-install-script-generation.md` & `08-installation-flow.md`: Cross-platform installer generation.
- [ ] `09-changelog-integration.md` & `14-release-body-and-changelog.md`: Automated changelog parsing and notes.
- [ ] `10-code-signing.md`: Binary signing for Windows (`signtool`/`osslsigncode`) and macOS (`codesign`).
- [ ] `11-self-update-mechanism.md` & `12-version-and-help.md`: Self-update validation in CI.
- [ ] `13-environment-variable-setup.md`: Secure secret injection and default environment fallbacks.
- [ ] `15-terminal-output-standards.md`: Clean, color-coded ANSI logging without clutter.
- [ ] `16-binary-icon-branding.md`: Resource embedding (`rsrc`, `windres`) in CI build steps.
- [ ] `17-release-pipeline-issues-rca.md`: 4-part RCA logging for failed release stages.
- [ ] `18-blue-green-deployment.md`: Zero-downtime deployment pipelines.
- [ ] `19-flaky-test-quarantine.md`: Quarantine strategies and automatic retry bounds.
- [ ] `20-contract-testing.md`: API schema validation and contract tests.
- [ ] `21-e2e-testing-pattern.md`: End-to-end integration and browser automated test flows.

### B. Coding Guidelines CI/CD Integration (`spec/02-coding-guidelines/06-cicd-integration/`)

- [ ] Read `00-overview.md` or `index.md`
- [ ] `01-sarif-contract.md`: Standard SARIF format for static analysis output.
- [ ] `02-plugin-model.md`: Extensible plugin architectures for linting.
- [ ] `04-ci-templates.md`: Reusable GitHub Actions / GitLab CI workflow templates.
- [ ] `05-distribution.md`: Packaging linters and CLI tools for multi-architecture distribution.
- [ ] `06-rules-mapping.md`: Mapping code-red rules to linter error codes.
- [ ] `07-performance.md`: Caching dependencies, build matrices, and sub-minute execution goals.
- [ ] `08-fix-repo-and-installers/`: (`00-overview.md`, `01-fix-repo-contract.md`, `02-installer-contract.md`, `03-visibility-change-contract.md`).

### C. Cross-Platform Automation & Update Architecture

- [ ] `spec/11-powershell-integration/`: (`00-overview.md`, `02-configuration-schema.md`, `04-script-reference.md`, `05-integration-guide.md`, `06-error-codes.md`).
- [ ] `spec/13-generic-cli/`: (`11-build-deploy.md`, `12-testing.md`, `18-batch-execution.md`, `20-terminal-output-design.md`, `22-self-update-gold-standard.md`).
- [ ] `spec/14-update/`: (`05-build-scripts.md`, `17-cross-compilation.md`, `18-release-pipeline.md`, `19-install-scripts.md`, `24-update-check-mechanism/`).
- [ ] `spec/15-distribution-and-runner/`: (`02-install-contract.md`, `03-runner-contract.md`, `04-release-pipeline.md`, `05-install-config.md`).
- [ ] `spec/16-generic-release/`: (`02-cross-compilation.md`, `05-release-pipeline.md`, `04-install-scripts.md`, `05-checksums-verification.md`, `06-release-assets.md`, `09-version-pinned-release-installers.md`).
- [ ] `spec/17-consolidated-guidelines/18-cicd-pipeline-workflows.md`: Master consolidated CI/CD reference.

---

## 2. Strict In-Repository Execution & Cross-Platform Python-First Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/` and `.lovable/cicd-issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Python-First Cross-Platform Automation:** All CI/CD build scripts, test runners, validation checks, and linters MUST be written in **Python 3** (`subprocess`, `sys`, `os`, `pathlib`, `concurrent.futures`, `json`, `shutil`) ensuring identical, deterministic execution across **Windows, Linux, and macOS**. Shell scripts (`.sh`, `.ps1`) must ONLY act as lightweight one-line entrypoints invoking Python.
> 4. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.lovable/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `spec/03-error-manage/01-index.md`, `.lovable/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands](spec/13-generic-cli/01-index.md)`
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

### Standard Cross-Platform Python Scripts Architecture

```text
scripts/ or 03-ai-scripts/
├── 01-file-manipulator.py       # File renaming, sequence renumbering, CRLF normalization
├── 05-guideline-autofixer.py    # Auto-formats return new lines and cleans booleans
├── 06-cicd-local-runner.py      # Master native host runner executing all CI/CD jobs
├── build.py                     # Cross-platform compilation & bundling script
├── test.py                      # Parallel test runner with coverage aggregation
└── lint.py                      # Multi-linter orchestrator (markdownlint, eslint, tsc, go, etc.)
```

---

## 3. Phase 1: Local Automation & Python Runner Setup (Steps 1 to PHASE_1_STEPS)

> [!IMPORTANT]
> **Build the local cross-platform foundation first.** Remote CI/CD workflows merely mirror the local Python runner.

### Step 1: Detect Tech Stack & Toolchains

1. Inspect project files (`package.json`, `go.mod`, `pyproject.toml`, `Cargo.toml`, `composer.json`, `version.json`).
2. Identify target compilers, linters, and test engines.
3. Record exact commands needed for:
   - Dependency installation (`npm ci`, `go mod download`, `pip install -r requirements.txt`, `composer install`).
   - Linting & static analysis (`golangci-lint`, `eslint`, `markdownlint`, `tsc --noEmit`, `phpcs`, `mypy`).
   - Build & compilation (`npm run build`, `go build`, `python -m build`, `cargo build`).
   - Automated testing (`npm test`, `go test -race ./...`, `pytest`, `cargo test`).

### Step 2: Implement Cross-Platform Python Linters

If custom checks or guideline validations are required:

1. Write or update Python scripts in `linter-scripts/` (e.g., `validate-guidelines.py`, `check-markdown-header-spacing.py`).
2. Ensure linters output standard error formats: `file:line:col: [RULE-CODE] message`.
3. Support `--fix` flag wherever deterministic AST/regex automated repair is possible.

### Step 3: Check for Existing Script & Force Override

- Check if `03-ai-scripts/06-cicd-local-runner.py` already exists.
- **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
- If it EXISTS and the user did **not** specify `force`: proceed to verification.
- If it is MISSING or `force` was requested: generate the script following the worker pool architecture below.

### Step 4: Write `06-cicd-local-runner.py` (Worker Pool & Log Aggregation Architecture)

Generate `03-ai-scripts/06-cicd-local-runner.py` following these architectural requirements:

1. **Round-Robin Worker Process / Thread Pool Architecture:** Runs tasks (tests, linters, builds) concurrently using `concurrent.futures.ThreadPoolExecutor(max_workers=3)` (2–3 concurrent tasks).
2. **Enqueuing Announcement:** The script MUST announce upfront how many tasks it has enqueued across the worker pool (e.g. `[INFO] Enqueued 20 quality gates across 3 concurrent workers...`).
3. **Real-Time Progress & Timing:** Prints job completions in real time with individual runtimes (e.g. `PASS [Job Name] (X.XXs)`).
4. **Graceful Non-Cancelling Failure Handling:** If one job fails in a running batch, the runner DOES NOT abort or cancel other active workers. It lets running tasks finish gracefully, capturing all stdout and stderr.
5. **Consolidated Summary Report & Full Diagnostic Logs:** Prints a complete final summary with total executed, passed, failed, and timeouts. For every failed or timed-out job, it outputs the full command line, return code, stdout, and stderr so the agent has 100% complete RCA context.
6. **Clean Exit Code:** Exits with code 0 only when ALL jobs pass; exits non-zero if ANY job fails.

**Standard generated runner template:**

```python
#!/usr/bin/env python3
"""Auto-generated Cross-Platform CI/CD Local Runner with concurrent worker pool."""
from concurrent.futures import ThreadPoolExecutor, as_completed, TimeoutError as FuturesTimeout
import os
from pathlib import Path
import subprocess
import sys
import time

# ── Configurable Variables ──────────────────────────────────────────────────
BATCH_SIZE      = 3    # Number of jobs to run concurrently (round-robin worker pool)
JOB_TIMEOUT_SEC = 120  # Maximum seconds before a single job is timed out

os.environ.setdefault("CI", "true")

# Populate JOBS dictionary from project requirements
JOBS = {
    "lint":      [sys.executable, "scripts/lint.py"],
    "typecheck": ["npx", "tsc", "--noEmit"],
    "build":     [sys.executable, "scripts/build.py"],
    "test":      [sys.executable, "scripts/test.py"],
}

def run_job(name: str, cmd: list):
    start_time = time.time()
    try:
        res = subprocess.run(cmd, capture_output=True, text=True, timeout=JOB_TIMEOUT_SEC)
        elapsed = round(time.time() - start_time, 2)
        return name, cmd, res.returncode, res.stdout, res.stderr, elapsed
    except subprocess.TimeoutExpired as e:
        elapsed = round(time.time() - start_time, 2)
        return name, cmd, "timeout", e.stdout or "", f"Timed out after {JOB_TIMEOUT_SEC}s", elapsed
    except Exception as e:
        elapsed = round(time.time() - start_time, 2)
        return name, cmd, 1, "", str(e), elapsed

def main():
    job_items = list(JOBS.items())
    total_jobs = len(job_items)
    print(f"[INFO] Enqueued {total_jobs} quality gates across {BATCH_SIZE} concurrent workers...\n")

    all_results = {}
    total_start = time.time()

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

    total_elapsed = round(time.time() - total_start, 2)

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

## 4. Phase 2: Workflow Generation & Green Gate Verification (Steps PHASE_1_STEPS+1 to N)

### Step 1: Generate GitHub Actions Workflows (`.github/workflows/`)

1. Create `.github/workflows/ci.yml` covering:
   - Matrix builds across OS platforms (`ubuntu-latest`, `windows-latest`, `macos-latest`) where applicable.
   - Cache steps (`actions/cache`, `actions/setup-node`, `actions/setup-python`, `actions/setup-go`).
   - Linting, typechecking, building, and testing jobs.
   - Invocation of native Python runner and scripts.
2. Create `.github/workflows/release.yml` covering:
   - Tagged release triggers (`v*`).
   - Multi-arch cross-compilation.
   - Checksum generation (`sha256sum`).
   - Release creation with `--notes-file` and Quick Install one-liners.

### Step 2: Autonomous Local Verification Loop

Execute the local runner and resolve every failure singly using self-looping turns:

```text
LOOP:
    1. Execute: python 03-ai-scripts/06-cicd-local-runner.py
    2. IF exit code == 0:
           SUCCESS -> Proceed to Stage & Commit.
    3. IF exit code != 0:
           Zero in on the first failing job.
           Document 4-part RCA in .lovable/memory/issues/xx-<slug>.md.
           Apply surgical fix to source code or python script.
           Re-verify with 06-cicd-local-runner.py.
           Repeat loop.
```

---

## STRICT AVOIDANCE: Never Disable CLI Linting or CI/CD Checks

> [!CAUTION]
> **TOTAL BAN ON DISABLING, SKIPPING, OR BYPASSING CLI LINTERS AND CI/CD GATES:**
>
> - **NEVER** disable, comment out, delete, or skip any CLI linting command (`golangci-lint`, `eslint`, `markdownlint`, `tsc`, `pytest`, `phpstan`, `mypy`, `check-*.py`), build step, or test suite.
> - **NEVER** add `|| true`, `continue-on-error: true`, `# nolint`, `// eslint-disable`, or ignore flags to "quickly win the race" or fake a pipeline pass.
> - **Your job is to legitimately fix the underlying source code and python scripts.** If resolving complex lint errors or test failures requires multiple sub-steps or nested self-looping turns, you MUST execute all necessary turns until the code is 100% clean and compliant.
> - Disabling or bypassing any CI/CD or CLI lint check is an automatic and immediate rejection.

---

## Strictly Avoid: No Automatic Releases & Run All Tests (Strict Policy)

> [!CAUTION]
> **NO AUTOMATIC RELEASES:** This prompt creates and verifies CI/CD infrastructure. You MUST NOT bump versions, update changelogs, or cut a release unless the user explicitly commands it in chat (e.g., "cut a release" or "bump the version"). All commits use `feat(ci): <description>` or `fix(ci): <description>`.
> **RUN ALL TESTS & QUALITY GATES (FULL PIPELINE FIDELITY):** All newly created CI/CD pipelines, local runners, and jobs MUST execute and verify all unit test suites, integration tests, and quality gates properly (`python 03-ai-scripts/06-cicd-local-runner.py`). NEVER skip or disable tests with `--no-tests` during pipeline creation and verification.
> **ATOMIC CHANGE TRACKING:** Record all modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.

---

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write all CI/CD pipelines in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently (up to 2-3 threads max). Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only write Python lint script for Go"). Never give a sub-agent a generic or multi-directory task.

---

## Pre-Reply / Loop Checklist (Must Verify Every Turn)

- [ ] **Specs Ingested:** Read `index.md` or `00-overview.md` and all subfiles across `spec/12-cicd-pipeline-workflows/` and `spec/02-coding-guidelines/06-cicd-integration/`.
- [ ] **Cross-Platform Python-First:** All automation, build, test, and linter scripts written in Python 3.
- [ ] **Local Runner Configured:** `03-ai-scripts/06-cicd-local-runner.py` created and tested with batching.
- [ ] **Workflows Generated:** `.github/workflows/ci.yml` and `.github/workflows/release.yml` created.
- [ ] **Zero Linting/CI/CD Bypass:** Confirmed NO linters, static analysis tools, or test scripts were disabled, commented out, or bypassed with `|| true`.
- [ ] **Local Runner 100% Green:** `python 03-ai-scripts/06-cicd-local-runner.py` exited with code 0.
- [ ] **Strict Lowercase Filenames:** All generated files use strictly lowercase naming.
- [ ] **Stage, Commit & Push:** Grouped fixes into clean commit and pushed to remote branch.
