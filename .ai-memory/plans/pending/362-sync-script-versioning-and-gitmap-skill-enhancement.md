# Plan 362: Sync Script Versioning, GitMap Cat Range Inspection, and Ad-Hoc Python Reading Ban

- **Plan Slug:** `sync-script-versioning-and-gitmap-skill-enhancement`
- **Root Slug:** `sync-script-versioning-and-gitmap-skill-enhancement`
- **Spec Reference:** `02-spec/21-app/362-sync-script-versioning-and-gitmap-skill-enhancement/01-architecture-spec.md`
- **Status:** `pending`
- **Created At:** 2026-10-31

---

## 1. Goal

Eliminate ad-hoc Python file reading and script abuse across AI agents by empowering `gitmap cat` with native line range slicing, line numbering, head/tail extraction, and query filtering. Introduce dual versioning for fleet sync scripts in `coding-guidelines` (V6 unrestricted vs V7 time-restricted & autobot governed), enshrine non-negotiable bans in `.agents/skills/gitmap/skill.md`, synchronize across repositories, and release GitMap `v6.636.0`.

---

## 2. Subtasks Breakdown

### Task 01: Architecture & Spec Authoring
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 01`
- **Subtask File:** `.ai-memory/plans/subtasks/362-sync-script-versioning-and-gitmap-skill-enhancement/01-architecture-and-spec-authoring.md`
- **Description:** Author 01-architecture-spec.md, 02-component-and-cli-spec.md, plan, and register subtasks in agent-task.db.

### Task 02: `gitmap cat` Range, Query, and Line Inspection Engine
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 02`
- **Subtask File:** `.ai-memory/plans/subtasks/362-sync-script-versioning-and-gitmap-skill-enhancement/02-cat-range-query-and-line-inspection.md`
- **Description:** Implement `cli/cmdmacro/cat_options.go` with `-L start:end`, `-n / --head`, `--tail`, `--num`, `-q / --query`, format banner stats, and update `cli/helpdoc/cat.md`.

### Task 03: Sync Script Dual Versioning in `coding-guidelines`
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 03`
- **Subtask File:** `.ai-memory/plans/subtasks/362-sync-script-versioning-and-gitmap-skill-enhancement/03-sync-script-dual-versioning.md`
- **Description:** Add `--mode v7` / `--v7` support to `38-sync-prompts-skills-scripts.py` and author dedicated `38-sync-prompts-skills-scripts-v7.py` with time-restriction and autobot checks.

### Task 04: GitMap Skill Enhancement & Python Reading Ban
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 04`
- **Subtask File:** `.ai-memory/plans/subtasks/362-sync-script-versioning-and-gitmap-skill-enhancement/04-gitmap-skill-enhancement-and-python-reading-ban.md`
- **Description:** Update `.agents/skills/gitmap/skill.md` with Rule 19 (Ban on Ad-Hoc Python Reading) and Rule 20 (Ban on Restrict Condition Bypasses), replacement matrix, and examples.

### Task 05: Fleet Synchronization Across Connected Repositories
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 05`
- **Subtask File:** `.ai-memory/plans/subtasks/362-sync-script-versioning-and-gitmap-skill-enhancement/05-fleet-synchronization.md`
- **Description:** Synchronize skills and prompts across connected repositories using `38-sync-prompts-skills-scripts.py`.

### Task 06: Verification, Minor Bump Release, and Completion Matrix
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 06`
- **Subtask File:** `.ai-memory/plans/subtasks/362-sync-script-versioning-and-gitmap-skill-enhancement/06-verification-release-and-matrix.md`
- **Description:** Compile with `-p 1`, deploy to LocalAppData, verify all commands, bump version to `v6.636.0`, commit via `gitmap cpr`, push upstream, and output completion tables.
