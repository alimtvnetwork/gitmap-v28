# Subtask 01: Git Sync, Preflight & Task Initialization

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`  
> **Status:** DONE  
> **Target Subsystems:** Project root, `.ai-memory/plans/`, `ledger.md`  
> **Owned Files:**  
> - `ledger.md`  
> - `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`  
> - `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/01-subtask-git-sync-and-preflight.md`  

---

## 1. Concrete Objectives

1. **Git Tree & Workspace Preflight Validation:**
   - Verify workspace cleanliness and working tree baseline hygiene before starting code modifications and specification authoring.
   - Confirm that no uncommitted temporary files, dangling scratch scripts, or unresolved merge conflicts exist.
   - Enforce the non-negotiable invariant: subagents never execute git mutation commands (`git add`, `git commit`, `git checkout`, etc.).

2. **Lifecycle Budget & Swarm Concurrency Setup:**
   - Establish task lifecycle budget: `N = 400` total steps.
     - **Phase 1 (Steps 1–100):** Architecture specification authoring, CLI command contracts, and subtask decomposition.
     - **Phase 2 (Steps 101–250):** Multi-mailbox checker implementation, schedule email CLI, folder watcher poller, privacy masking engine, and mail command tree.
     - **Phase 3 (Steps 251–400):** Embedded web dashboards (`/mail/ui` and `/schedule/ui`), unit test suites, regression prevention, and acceptance gate verification.
   - Configure swarm concurrency parameters: `A = 3` active subagents, `H = 2` human checkpoints, `Capacity = 6` maximum worker pool.

3. **Master Execution Plan Formulation & Registration:**
   - Ingest verbatim user request requiring scheduling additions (email, multi-mailbox checking, folder/file change triggers, notifier recipient dispatch), privacy masking, dedicated mail CLI verbs, and dual embedded UIs.
   - Author and register master execution plan at `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`.
   - Decompose task scope into 6 discrete, disjoint subtasks:
     - **Subtask-01:** Git Sync, Preflight & Task Initialization (Lead, Status: DONE)
     - **Subtask-02:** Architecture & Specification Authoring (Spec Authors A=3, Status: IN_PROGRESS)
     - **Subtask-03:** Multi-Mailbox Checker, Schedule Email CLI & Privacy Masker (Worker 01, Status: PENDING)
     - **Subtask-04:** Folder Watcher Event Trigger & Mail CLI Engine (Worker 02, Status: PENDING)
     - **Subtask-05:** Embedded Mailing Composer UI & Schedule Settings UI (Worker 03, Status: PENDING)
     - **Subtask-06:** Targeted Unit Tests & Verification Suite (Worker Swarm, Status: PENDING)

4. **Toolchain & Subsystem Preflight Validation:**
   - Confirm Go compiler toolchain availability and module health.
   - Verify presence of pure-Go CGO-free SQLite driver (`modernc.org/sqlite`) in `go.mod`.
   - Enforce standard library `net/http` utilization for embedded UIs without third-party web frameworks or CDN dependencies.
   - Verify existing `cli/cmdschedule` package structure and map integration touchpoints for `cli/cmdmail`.

5. **Ledger Task Initialization:**
   - Register active task slug `gitmap-mailing-and-scheduling-enhancements` in `ledger.md`.
   - Record allocated lifecycle budget (`N = 400`), concurrency model, and initial gate thresholds.

6. **Verification Quality Gates Baseline:**
   - Ratify Quality Gates G1 through G6:
     - **G1 (Multi-Mailbox Polling):** `gitmap schedule mail-check --all` evaluates configured mailboxes and unread counts.
     - **G2 (Scheduled Email Registration):** `gitmap schedule email` registers jobs in Split-DB table `ScheduledEmail`.
     - **G3 (Folder Event Triggering):** Creating or modifying monitored files dispatches debounced email notifications to configured recipients.
     - **G4 (Privacy Masking):** Sender and recipient addresses display masked format (`u***r@d***n.com`) by default across CLI and logs.
     - **G5 (Mail CLI Hierarchy):** `gitmap mail send`, `check`, `test`, `template`, `notify`, `ui` dispatch correctly.
     - **G6 (Embedded Web UIs):** Both `/mail/ui` and `/schedule/ui` serve responsive HTML/CSS/JS without external CDN calls.

---

## 2. Action Items Checklist

- [x] 1. Preflight workspace inspection executed; baseline working tree verified clean.
- [x] 2. Execution plan formulated and registered at `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`.
- [x] 3. Disjoint file ownership matrices defined for all 6 subtasks to prevent multi-agent collisions.
- [x] 4. Subtask directory `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/` created.
- [x] 5. Toolchain dependencies, Go modules, and pure-Go SQLite driver verified.
- [x] 6. Task registered in `ledger.md` with active status and N=400 budget.

---

## 3. Verification Evidence

- Master plan file: `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md` (validated, status: ACTIVE).
- Subtask plan file: `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/01-subtask-git-sync-and-preflight.md` (validated).
- Concurrency allocation: A=3 active agents, H=2 checkpoints, Capacity=6 capacity verified.
- Quality Gates G1–G6 ratified and documented.
- Status: **DONE**
