---
name: ci-cd-create
description: >-
  Use this skill to autonomously design, create, and verify cross-platform Python CI/CD pipelines, linters, and workflows.
---

# Instruction (must follow): CI/CD Pipeline & Cross-Platform Python Automation Creation

Trigger Keywords & Aliases: `create-ci-cd`, `cicd-create`, `create cicd`, `setup cicd pipeline`, `build ci-cd pipeline`, `cicd create python`, `05-cicd-pipeline-create`

/goal Autonomously design, create, and verify complete cross-platform Python CI/CD pipelines, local runners, and GitHub Actions workflows by following specifications in `02-spec/12-cicd-pipeline-workflows/` and `02-spec/02-coding-guidelines/06-cicd-integration/`, maintaining a continuous N-step self-loop until all jobs exit code 0.

/learn Ingest `02-spec/12-cicd-pipeline-workflows/`, `02-spec/02-coding-guidelines/06-cicd-integration/`, `02-spec/11-powershell-integration/`, `02-spec/14-update/`, `02-spec/15-distribution-and-runner/`, `02-spec/16-generic-release/`, `02-spec/17-consolidated-guidelines/18-cicd-pipeline-workflows.md`, and `.lovable/coding-guidelines.md` before writing code.

---

## Variables — Configurable at Runtime

```text
N = 300  (Total self-loop steps budget)

PHASE_1_STEPS = N / 2  (Steps 1 .. 150: Spec Ingestion, Python Automation Design, Local Runner Setup)
PHASE_2_STEPS = N / 2  (Steps 151 .. 300: Workflow Generation, Python Linter Verification, Green Gate Loop)
```

---

## Strict In-Repository Execution & Cross-Platform Python-First Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/` and `.lovable/cicd-issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **Python-First Cross-Platform Automation:** All CI/CD runners, linters, test harnesses, and build orchestrators MUST be written in **Python 3** ensuring identical, deterministic execution across **Windows, Linux, and macOS**. Shell scripts (`.sh`, `.ps1`) are only thin wrappers invoking Python.
> 6. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.lovable/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `02-spec/03-error-manage/01-index.md`, `.lovable/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands](02-spec/13-generic-cli/01-index.md)`
> 7. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.
> 8. **Zero-Storage GitHub Actions Mandate (Total Ban on CI Artifact Uploads):** Workflows MUST NOT upload test outputs, coverage files, Playwright reports, or drift summaries via `actions/upload-artifact`. Free-tier accounts have a strict 0.5 GB shared quota across all repositories. All reports, failures, and summaries MUST be emitted directly to `$GITHUB_STEP_SUMMARY`, console stdout (`cat log.txt`), or sticky PR comments with zero storage consumption. Only true GitHub release assets (binaries/tarballs on tagged releases) are permitted.

---

## STRICT AVOIDANCE: Never Disable CLI Linting or CI/CD Checks

> [!CAUTION]
> **TOTAL BAN ON DISABLING, SKIPPING, OR BYPASSING CLI LINTERS AND CI/CD GATES:**
>
> - **NEVER** disable, comment out, delete, or skip any CLI linting command (`golangci-lint`, `eslint`, `markdownlint`, `tsc`, `pytest`, `phpstan`, `mypy`, `check-*.py`), build step, or test suite.
> - **NEVER** add `|| true`, `continue-on-error: true`, `# nolint`, `// eslint-disable`, or ignore flags to fake a pipeline pass.
> - Disabling or bypassing any CI/CD or CLI lint check is an automatic and immediate rejection.

---

## End of Tunnel Checklist

- [ ] All CI/CD specs (`02-spec/12-cicd-pipeline-workflows/`, `02-spec/02-coding-guidelines/06-cicd-integration/`) read and followed.
- [ ] Python cross-platform runner `03-ai-scripts/06-cicd-local-runner.py` created and verified.
- [ ] All linters in `linter-scripts/` verified and passing without bypass.
- [ ] GitHub Actions workflows created in `.github/workflows/` (strictly adhering to Zero-Storage mandate — no `actions/upload-artifact` in CI).
- [ ] Local runner exited with code 0 on all jobs.
- [ ] All changes committed and pushed to remote repository.
