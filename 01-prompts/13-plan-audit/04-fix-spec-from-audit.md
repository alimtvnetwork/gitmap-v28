# Specification Remediation from Audit Findings — Execution Spec (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously ingest the latest specification audit file from `spec/25-app-spec-audit/`, decompose every finding into an exhaustive 1:1 remediation checklist, spawn parallel subagents to fix the specifications, verify 100% compliance, and remove the audit gap at the final stage.

```text
N = 200
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Audit Ingestion, Finding Matrix & Subtask Decomposition)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Remediation, CI Verification & Gap Removal)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Shared Directory Contract (Non-Negotiable)

Both the auditor and fixer agents MUST operate against these exact paths:
- **Audit Reports Directory:** `spec/25-app-spec-audit/`
- **Default Target Spec Directory:** `spec/21-app/` (or user-specified subfolder)
- **Subtask Tracking:** `.lovable/plans/subtasks/xx-spec-fix/`
- **Completed Archive:** `.lovable/plans/completed/`
- **Agent State Directory:** `.lovable/temp-agents/`

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/fix-spec-from-audit/skill.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 4-Phase continuous loop. Do not skip steps.

### Phase 1: Audit Ingestion & 1:1 Finding Matrix (Steps 1 to PHASE_1_STEPS)

1. **Locate Latest Audit:** Scan `spec/25-app-spec-audit/` using `03-ai-scripts/17-fast-file-reader.py` and select the file with the highest numerical sequence prefix (`NN-audit-*.md`).
2. **Exhaustive Finding Parsing:** Parse the Markdown Summary Table at the bottom of the audit file. Extract every single row without skipping a single issue.
3. **Build 1:1 Remediation Ledger:** Create `.lovable/plans/pending/xx-spec-remediation.md` containing an explicit checkbox for every finding:
   ```markdown
   - [ ] Finding [ID]: File `<path>`, Issue: `<description>`, Remedy: `<proposed fix>`
   ```
4. **Lean Subtask Decomposition:** Group the findings by folder/file and write lean subtasks to `.lovable/plans/subtasks/xx-spec-fix/01-<slug>.md`. Each subtask MUST follow this format:
   ```markdown
   # Subtask: [Target Spec File]
   **Target File:** `spec/21-app/...`
   **Findings to Fix:**
   - [ID 1]: [Specific modification]
   - [ID 2]: [Specific modification]
   **Constraints:** No generic names, preserve line length <= 15, update cross-links.
   ```
5. **MANDATORY AUTO-LOOP (DO NOT STOP):** As soon as Phase 1 completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2.

### Phase 2: Parallel Multi-Agent Remediation (Steps PHASE_1_STEPS+1 to N)

1. **Parallel Dispatch:** Use the `invoke_subagent` tool to spawn up to 2 execution subagents concurrently (max 2 threads each), assigning disjoint subtasks from `.lovable/plans/subtasks/xx-spec-fix/`.
2. **Minimal Context Diet:** Provide subagents with minimal instructions (e.g., "Read `.lovable/plans/subtasks/xx-spec-fix/01-<slug>.md` and execute the fixes on the specified spec file").
3. **Isolated Agent State & Communication:** Each subagent MUST create `.lovable/temp-agents/xx-<task-name>/` and track its progress in `state.md`.
4. **Failure Protocol:** If a subagent fails, record the error in `.lovable/memory/issues/xx-spec-fix-failure.md`. The next subagent must read the failure log first to remediate.
5. **Mark Reconciliation Ledger:** As subtasks finish, mark the corresponding findings in `.lovable/plans/pending/xx-spec-remediation.md` as `[x]`.

### Phase 3: Quality Gate & Cross-Link Verification (Validation Gate)

1. **Cross-Link Integrity:** Run `python linter-scripts/check-spec-cross-links.py` to ensure no internal links were broken by the spec edits.
2. **Markdown Standards:** Run `python 03-ai-scripts/31-md-gap-fixer.py --fix` and verify spacing.
3. **Full CI Runner:** Run `python 03-ai-scripts/06-cicd-local-runner.py` ensuring all 36 quality gates exit with code 0 (`exit 0`).

### Phase 4: Audit Gap Removal & Final Archive (End of Loop)

> **CRITICAL (NON-NEGOTIABLE):** The audit gap MUST be officially closed on disk before concluding.

1. **Verify 100% Closure:** Confirm that every checkbox in `.lovable/plans/pending/xx-spec-remediation.md` is marked `[x]`.
2. **Remove Audit Gap on Disk:** Delete the original audit file from `spec/25-app-spec-audit/NN-audit-*.md` (or move it to `.lovable/plans/completed/NN-audit-*.md-resolved`) so no unresolved audit gaps remain in the active spec directory.
3. **Subtask Cleanup:** Consolidate completed subtasks into a single `.lovable/plans/completed/xx-spec-remediation-completed.md` file noting how many steps it took, and delete the granular `.lovable/plans/subtasks/xx-spec-fix/` files.
4. **Index Synchronization:** Update `.lovable/plans/01-index.md` and `.lovable/what-to-read.md` to reflect that the audit gap is 100% resolved.
5. **Git Commit:** Stage all modified spec files and commit with `fix(spec): remediate all audit findings and close audit gap`.

---

### Temp-Agent Isolated Task Directory & Communication Protocol (Non-Negotiable)

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.lovable/temp-agents/xx-<task-name>/`:

1. **Per-Task Isolation:** On task start, the assigned subagent creates its isolated directory `.lovable/temp-agents/xx-<task-name>/`.
2. **State & Progress Tracking:** Create `.lovable/temp-agents/xx-<task-name>/state.md` documenting:
   - `TASK_NAME`: `<task-name>`
   - `STATUS`: `IN_PROGRESS` | `DONE` | `FAILED`
   - `ASSIGNED_AGENT`: Agent identifier and thread index
   - `CURRENT_STEP`: Detailed micro-step description
3. **Inter-Agent Communication & Handoff:**
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.lovable/temp-agents/xx-<task-name>/`.
   - Sibling or successor agents MUST inspect this dedicated folder before resuming work or fixing errors.
4. **On Error/Crash:** Append the exact error, root cause, and `STATUS: FAILED` to `.lovable/temp-agents/xx-<task-name>/state.md` before exiting.
5. **On Success:** Mark `STATUS: DONE` in `.lovable/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

---

## Non-Negotiable Rules (Auto-Reject on Violation)

- [ ] Zero Skipped Findings: Every row in the audit Summary Table must have a corresponding code fix.
- [ ] No Lingering Audit Files: `spec/25-app-spec-audit/` must be clean of the resolved audit file.
- [ ] Strictly Unix LF (`\n`) line endings and UTF-8 encoding.
- [ ] No absolute file paths or `file:///` URIs.
- [ ] All 36 CI/CD gates green.

## MUST FOLLOW NON-NEGOTIABLE

Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job.
