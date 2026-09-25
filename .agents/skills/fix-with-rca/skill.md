---
name: fix-with-rca
description: Diagnose, analyze, and fix bugs using grounded 4-part Root Cause Analysis (RCA) without guessing.
---

# Fix with Root Cause Analysis (RCA)

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
> 
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

Solves bugs and pipeline failures through structured 4-part Root Cause Analysis.

## Issue Destination & Root Cause Analysis (RCA) Routing Mandate

Whenever resolving an issue or performing a fix with RCA:
- **CI/CD Issues & Pipeline Failures:** Record the RCA in `.ai-memory/cicd-issues/NN-<slug>.md` and register it in `.ai-memory/cicd-index.md`.
- **Non-CI/CD Issues (Application Bugs, Feature Defects, Logic/Runtime Errors):** Document in `02-spec/22-app-issues/NN-<slug>.md` (indexed in `02-spec/22-app-issues/readme.md`, cross-referencing in `.ai-memory/memory/issues/` for institutional memory).

## GitMap Pipeline-AI & Dynamic Waiting Protocol

When diagnosing bugs or failures that involve CI/CD pipelines or remote workflows:
1. **Pipeline Inspection:** The agent MUST use `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`) to inspect pipeline status and failure diagnostics.
2. **Targeted Failure Extraction:** Extract targeted failure lines (`##[error]`, `FAIL:`, compile errors) directly using GitMap's targeted diagnostic extraction instead of dumping entire passing logs. Feed these lines into the 4-part RCA document.
3. **Adaptive Waiting / Dynamic Sleep (Anti-Credit-Waste Mandate):** Rapid tight-loop polling (`gh run view` in a loop) is STRICTLY BANNED as it burns user credits.
   - When a pipeline is running, use `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`) to wait dynamically based on historical completion ETA.
   - Adaptive sleep intervals:
     - ETA > 120s: wait 20s-30s
     - 60s < ETA <= 120s: wait 10s-20s
     - ETA <= 60s: wait 5s-10s
4. **Atomic Change Tracking:** Record modified files safely under lock via:
   ```bash
   python 03-ai-scripts/33-test-inventory-generator.py --record <files...>
   ```

## 4-Part RCA Structure

1. **Symptom:** Raw error output, failing command, and reproduction steps (or GitMap targeted failure output).
2. **Root Cause:** One-sentence root cause identifying the exact mechanism of failure.
3. **Resolution:** Direct minimal fix applied to code or configuration.
4. **Prevention & Learnings:** Specific rule or avoidance logged to prevent recurring failures.

## Verification & Commit Checklist

- [ ] **Top-Instruction Priority Mandate:** Verified that directives before this section or prompt take highest priority and are non-negotiable.
- [ ] **Issue Destination Routing:** Verified that CI/CD failures are recorded in `.ai-memory/cicd-issues/NN-<slug>.md` (and `.ai-memory/cicd-index.md`), while application bugs are recorded in `02-spec/22-app-issues/NN-<slug>.md` (and `02-spec/22-app-issues/readme.md`).
- [ ] **Atomic Change Tracking:** Recorded all modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
- [ ] **Targeted Verification:** Executed targeted tests/linters on modified packages (avoiding broad test runner loops).
- [ ] **Atomic Commit & Push:** All modifications committed in a single atomic commit and pushed to origin (ban on per-file commits).
