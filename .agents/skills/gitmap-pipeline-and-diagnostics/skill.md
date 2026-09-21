---
name: gitmap-pipeline-and-diagnostics
description: Autonomously manage, optimize, and diagnose CI/CD pipeline runs, error extraction, bounded stack traces, and 4-part RCA generation.
---

# GitMap Pipeline & Diagnostics Skill (`gitmap-pipeline-and-diagnostics`)

## Mission & Purpose
This skill provides authoritative architectural guidance, code navigation, and execution rules for managing GitMap's CI/CD pipelines, bounded stack trace extraction, dynamic ETA waiting, and Root Cause Analysis (RCA) records.

---

## 1. Key Architectural Components & Code Map

| Component | Primary Location | Key Responsibilities |
|---|---|---|
| **Pipeline CLI Engine** | `cli/cmdpipeline/` | Subcommands for `status`, `error-logs` (`pe`), `details` (`pd`), `purge`, and live error streaming. |
| **Bounded Stack Extractor** | `cli/cmdpipeline/pipeline_stacktrace.go` | Extracts strictly 5 preceding and 20 trailing failure frame lines. Halts on exit code or job boundary to prevent clipboard explosion (RCA 58). |
| **Local CI/CD Runner** | `03-ai-scripts/06-cicd-local-runner.py` | Multi-threaded runner executing linters and gates via `ThreadPoolExecutor` with `--failed` log isolation. |
| **RCA Knowledge Base** | `.ai-memory/cicd-issues/` | Numbered post-mortem documents (`xx-<slug>-rca.md`) capturing root causes, reproduction steps, and applied fixes. |
| **Actions Zero-Storage** | `cli/cmdpipeline/purge.go` | Automated pruning of old GitHub Actions run artifacts maintaining 0.0 GB storage footprint. |

---

## 2. Essential Commands

```bash
# Check remote or local pipeline status with JSON output
gitmap pipeline-ai status --json

# Dynamic non-polling ETA wait for running pipelines
gitmap pipeline-ai status -t <eta>

# Extract failing step error logs (alias: gitmap pe)
gitmap pipeline error-logs
gitmap pe

# Clear past error logs
gitmap pe clear -y

# Display pipeline execution details with target table (alias: gitmap pd)
gitmap pipeline details
gitmap pd

# Zero-storage artifact purge
gitmap pipeline purge
```

---

## 3. Core Invariants & Engineering Guardrails

1. **Bounded Failure Context (RCA 58):** Never dump entire raw log files or unbounded stack traces to chat, clipboard, or reports. Always bound extractions to the specific failing frame (5 before, 20 after) using `extractBoundedStackLines`.
2. **Terminal Boundary Halting:** Stack trace extraction must halt immediately upon encountering error delimiters (`[ERROR]`, `exit code`, or job transitions).
3. **No Polling Loops:** When waiting for CI pipelines or asynchronous commands, agents must never busy-poll in tight loops. Use dynamic sleep protocols based on `.ai-memory/temp/runner-eta.json` or in-flight 25-second heartbeats.
4. **4-Part Root Cause Analysis:** Every recurring or significant CI failure must be documented in `.ai-memory/cicd-issues/` with:
   - Part 1: Symptoms & Error Signature
   - Part 2: Root Cause (mechanics of why it failed)
   - Part 3: Applied Solution & Code Changes
   - Part 4: Prevention Invariant (why it will never happen again)
5. **No Disabling of CI/CD:** Strictly avoid commenting out, bypassing, or removing CI/CD checks to force a pass. Fix the underlying code.
