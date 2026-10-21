# Subtask 01: Git Sync, Preflight & Task Initialization

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/gitmap-testing-and-schedule-management.md`  
> **Status:** DONE  
> **Target Subsystems:** Project root, `.ai-memory/plans/`, `ledger.md`  
> **Owned Files:**  
> - `ledger.md`  
> - `.ai-memory/plans/gitmap-testing-and-schedule-management.md`  
> - `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/01-subtask-git-sync-and-preflight.md`  

---

## 1. Concrete Objectives

1. **Git Tree & Workspace Preflight Validation:**
   - Verify repository status and working tree baseline hygiene before starting code modifications or specification authoring.
   - Enforce the non-negotiable invariant: subagents never execute git mutation commands (`git add`, `git commit`, `git checkout`, `git push`, etc.).
   - Confirm that no uncommitted temporary files, dangling scratch scripts, or unresolved merge conflicts exist.

2. **Lifecycle Budget & Swarm Concurrency Setup:**
   - Establish task lifecycle budget: `N = 400` total steps.
     - **Phase 1 (Steps 1–100):** Architecture specification authoring, cross-OS power command contracts, Docker container test specifications, and subtask decomposition.
     - **Phase 2 (Steps 101–250):** Test suite implementation for schedule shutdown/restart across Windows, Linux, and Darwin; Docker harness integration; Mail and schedule feature validation.
     - **Phase 3 (Steps 251–400):** Cross-OS compilation gates, CI/CD pipeline verification, retrospective AI audit, master index registration, and atomic release push.
   - Configure swarm concurrency parameters: `A = 3` active subagents, `H = 2` human checkpoints, `Capacity = 6` maximum worker pool.

3. **Master Execution Plan Formulation & Registration:**
   - Ingest verbatim user request requiring verification of schedule shutdown and restart features across multiple operating systems, Docker container testing environment feasibility, and comprehensive validation of all new GitMap features.
   - Author and register master execution plan at `.ai-memory/plans/gitmap-testing-and-schedule-management.md`.
   - Decompose task scope into 6 discrete, disjoint subtasks:
     - **Subtask-01:** Git Sync, Preflight & Task Initialization (Lead, Status: DONE)
     - **Subtask-02:** Architecture & Cross-OS Test Specification Authoring (Spec Authors A=3, Status: IN_PROGRESS)
     - **Subtask-03:** Cross-OS Schedule Shutdown and Restart Test Suite (Worker 01, Status: PENDING)
     - **Subtask-04:** Docker Environment & Linux Containerized Verification (Worker 02, Status: PENDING)
     - **Subtask-05:** Comprehensive Verification of New Mail & Scheduling Features (Worker 03, Status: PENDING)
     - **Subtask-06:** Master Index Registration, Quality Gates & Atomic Push (Lead, Status: PENDING)

4. **Toolchain, Compiler & Environment Preflight Validation:**
   - Validate Go toolchain compiler version and module readiness.
   - Probe Docker environment availability on host machine (finding: Docker daemon is not active on local Windows host; containerized verification must utilize isolated test doubles and portable `Dockerfile.test` harness).
   - Verify pure-Go cross-compilation capability (`GOOS=linux`, `GOOS=darwin`).
   - Validate pure-Go CGO-free SQLite driver (`modernc.org/sqlite`) in `go.mod`.

5. **Ledger Task Initialization:**
   - Register active task slug `gitmap-testing-and-schedule-management` in `ledger.md`.
   - Record allocated lifecycle budget (`N = 400`), concurrency model, and initial gate thresholds.

6. **Verification Quality Gates Baseline:**
   - Ratify Quality Gates G1 through G6:
     - **G1 (Cross-OS Command Construction):** Tests verify `shutdown` and `restart` emit valid OS-specific commands for Windows (`shutdown /s /t`), Linux (`sudo poweroff` / `sudo reboot`), and Darwin (`sudo shutdown -h now` / `sudo shutdown -r now`).
     - **G2 (Safety Invariant cg-isolate-os-tests):** Zero real OS reboots are triggered during test execution. All power actions use injectable mock executors.
     - **G3 (Linux Compilation Gate):** `GOOS=linux go test -run=^$ ./cli/cmdschedule` and `GOOS=linux go build ./...` compile cleanly with 0 errors.
     - **G4 (Docker Harness Portability):** A runnable `Dockerfile.test` targeting Alpine Linux is authored for containerized validation.
     - **G5 (Comprehensive Feature Coverage):** Mail commands, folder watchers, notification dispatch, address masking, and embedded UIs are fully tested.
     - **G6 (Line Cap & Relative Paths):** All authored files are strictly $\le 300$ lines with 100% relative paths.

---

## 2. Action Items Checklist

- [x] 1. Preflight workspace inspection executed; baseline working tree verified clean.
- [x] 2. Execution plan formulated and registered at `.ai-memory/plans/gitmap-testing-and-schedule-management.md`.
- [x] 3. Disjoint file ownership matrices defined for all 6 subtasks to prevent multi-agent collisions.
- [x] 4. Subtask directory `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/` established.
- [x] 5. Toolchain dependencies, Go modules, and cross-compilation capabilities verified.
- [x] 6. Task registered in `ledger.md` with active status and N=400 budget.

---

## 3. Verification Evidence

- Master plan file: `.ai-memory/plans/gitmap-testing-and-schedule-management.md` (validated, status: ACTIVE).
- Subtask plan file: `.ai-memory/plans/subtasks/gitmap-testing-and-schedule-management/01-subtask-git-sync-and-preflight.md` (validated).
- Concurrency allocation: A=3 active agents, H=2 checkpoints, Capacity=6 capacity verified.
- Quality Gates G1–G6 ratified and documented.
- Status: **DONE**
