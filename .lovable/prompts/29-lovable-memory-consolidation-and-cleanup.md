# 29-lovable-memory-consolidation-and-cleanup.md

# Lovable Memory Consolidation, Cleanup & Audit Removal — Workflow (must follow)

Trigger Keywords & Aliases: `consolidate-plans`, `consolidate completed plans`, `clean completed plans`, `resequence completed plans`, `merge plans`, `archive completed plans`, `cleanup plans completed`, `memory consolidation`, `backup and consolidate plans`, `compact plans`, `reduce plan file count`, `compact completed plans`

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously create a timestamped backup branch, scan, analyze, cluster, aggressively consolidate, and re-sequence all completed plan files and subtasks within `.lovable/plans/` into minimal, hyper-compact milestone summaries, combining 2, 3, or more related tasks and common checklists into single files to drastically reduce total file count while strictly preserving 100% of core architectural contracts, error codes, and verified outcomes.

---

## 🛑 MANDATORY RULES FOR CONSOLIDATION & CLEANUP

### 1. Mandatory Removal of Resolved Audit Folders
When an audit has been figured out or resolved, all temporary/stale audit folders must be removed from the active workspace and purged using `03-ai-scripts/33-git-history-tracer-and-purger.py`:
- **Folder 25 (Spec Audit):** `spec/21-app/25-app-spec-audit/`
- **Folder 26 (Coding Guideline Audit):** `spec/21-app/26-coding-guideline-audit/`
- **Main Worker Service Audit:** `spec/19-main-worker-service/audit/`
- **Lovable Audits:** `.lovable/audits/`
- **Stale Audit Scripts:** `audit.py`, `audit_extra.py`, `audit_summary.py`, `audit_ts.py`

**Removal Execution Protocol:**
1. Back up all files to OS temp directory (`%TEMP%` on Windows, `/tmp` on Unix) to enable instant rollback.
2. Delete workspace files using the **Recycle Bin** via `python 03-ai-scripts/33-git-history-tracer-and-purger.py --spec-25-audit --delete` or `--purge`.
3. Report the exact backup location path immediately upon completion.

### 2. Streamlining Coding Guidelines in Milestone Files
- **Do NOT Repeat Full Coding Guidelines:** Milestone summaries exist solely to chronicle **which business logic tasks and feature contracts are done**. Never paste redundant multi-page coding guidelines into milestone files.
- **Single Source of Truth Reference:** Represent coding guidelines as a compact 5-bullet checklist referencing `.lovable/coding-guidelines.md`:
  - Functions <= 15 lines body cap.
  - Affirmative boolean naming (`is_*`, `has_*`) with implicit evaluation.
  - Guard returns and zero nested `if` statements.
  - Zero swallowed errors (`*apperror.AppError` and `cliexit`).
  - 100% relative repository paths (zero absolute paths or `file:///` URIs).
- **Pruning Pure Coding Guideline Tasks:** If a historical task or subtask was purely about "fixing coding guidelines" (e.g. style fixes, boolean renaming, function sizing, line gaps) and contains no application business logic, **prune / remove that task** from the final milestone summaries. Keep milestones focused strictly on domain features, database engines, CLI commands, installers, and architectural contracts.

---

## Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Safety Backup & Audit Preflight): Create timestamped safety backup branch `backup/plans-consolidation-YYYYMMDD-HHMMSS` and verify clean working tree.
2. [ ] /goal Phase 1 (Audit Cleanup): Scan and remove resolved audit folders (`spec/21-app/25-app-spec-audit/`, `spec/19-main-worker-service/audit/`, `.lovable/audits/`) using `33-git-history-tracer-and-purger.py` with OS temp backup and Recycle Bin deletion.
3. [ ] /goal Phase 1 (Inventory & Clustering): Group all completed plans and subtasks into maximum 10-12 high-density milestone summaries, discarding tasks that only fixed coding guidelines without business logic.
4. [ ] /goal Phase 2 (Authoring Milestones): Author consolidated milestone summaries in `.lovable/plans/completed/` embedding business logic, contracts, error envelopes, and single-file checklist references.
5. [ ] /goal Phase 2 (Pruning Superseded Files): Cleanly remove superseded micro-plan files and subtask folders via `git rm`.
6. [ ] /goal Phase 2 (Resequencing & Index Sync): Verify monotonic sequence `01-`, `02-`... with zero gaps; synchronize `.lovable/plans/01-index.md`, `.lovable/plans/index.md`, and `.lovable/what-to-read.md`.
7. [ ] /goal Phase 2 (Verification): Run `python 03-ai-scripts/06-cicd-local-runner.py` to ensure all CI/CD quality gates pass with exit code 0.
