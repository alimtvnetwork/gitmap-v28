# Parent Task Plan: 71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix

**Task Code:** 71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix  
**Status:** In Progress (Phase 1 Planning & Discovery)  
**Budget:** N = 300 steps (Phase 1: 150 steps, Phase 2: 150 steps)  
**Concurrency:** A = 2 subagents, H = 2 operational hands  
**Ledger:** `02-spec/21-app/71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix/00-master-audit-ledger.md`  

---

## Objective

Deliver an end-to-end upgrade to GitMap's pull pipeline (`gitmap pa` / `gitmap pull`) and pipeline error inspector (`gitmap pe`):
1. Support `-l`, `-limit`, `--limit`, `-n`, `--lines` flags in `gitmap pe` without misidentifying numbers as target repositories.
2. Root-cause and fix the "unknown (repo-cache)" issue in `gitmap pa`.
3. Implement automated Fast-Forward merge fallback (without rebase) when pull encounters divergent branches.
4. Provide a single command (`gitmap pull-fix` / `gitmap pa --fix`) to resolve all failed repositories together, plus an interactive prompt at the end of `gitmap pa` offering "all together" or "one by one" remediation.
5. Implement CPU-aware auto-scaling for `gitmap pa` concurrency with two presets (low pressure vs high pressure) and configuration flags in `gitmap pa --help`.
6. Full verification, linters, and minor release ceremony.

---

## Execution Phases

### Phase 1: Planning, Discovery & Detailed Spec (Steps 1–150)
- **Step 1A:** Capture user request losslessly, extract discrete deliverables (Task-01 to Task-06).
- **Step 1B:** Spawn 2 parallel research subagents to discover exact code paths:
  - Research 01: `gitmap pe` flag parsing & `unknown (repo-cache)` name resolution.
  - Research 02: Pull-All FF-merge execution, failure tracking/batch fix, and CPU auto-scaling.
- **Step 1C:** Author detailed architecture and component specifications:
  - `02-spec/21-app/71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix/01-architecture-spec.md`
  - `02-spec/21-app/71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix/02-component-spec.md`
- **Step 1D:** Author modular subtask execution plans in `.ai-memory/plans/subtasks/71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix/`.

### Phase 2: Autonomous Multi-Agent Execution & Self-Looping (Steps 151–300)
- **Wave 1:**
  - Worker 01: Task-01 (`gitmap pe` flags) & Task-02 (`unknown repo-cache` fix).
  - Worker 02: Task-03 (Auto FF-merge in pull-all) & Task-04 (Batch failure fix & prompt).
- **Wave 2:**
  - Worker 01: Task-05 (CPU auto-scaling & presets).
  - Worker 02: Quality gates, unit tests, guideline autofixer, and nested if verification.
- **Wave 3:**
  - Lead Orchestrator: Version bump, changelog, and minor release ceremony.
