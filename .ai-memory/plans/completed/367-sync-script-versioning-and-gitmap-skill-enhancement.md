# Plan 367: Sync Script Versioning, GitMap Skill Anti-Pattern Ban Matrix & Cat Query Window Enhancements

- **Plan Slug:** `sync-script-versioning-and-gitmap-skill-enhancement`
- **Root Slug:** `sync-script-versioning-and-gitmap-skill-enhancement`
- **Specification:** [02-spec/21-app/367-sync-script-versioning-and-gitmap-skill-enhancement/readme.md](../../../02-spec/21-app/367-sync-script-versioning-and-gitmap-skill-enhancement/readme.md)
- **Status:** `completed`
- **Created At:** 31-Oct-26
- **Architecture Lead:** V7 Master Autonomous Execution Engine

---

## 1. Goal

Eliminate ad-hoc Python file reading, filesystem crawling, sleep polling, and Node/PowerShell pipeline abuses across autonomous AI agents by upgrading `gitmap cat` with `-t` tail flag and combined query-window filtering. Establish strict dual-versioned fleet sync scripts in `coding-guidelines` (V6 unrestricted vs V7 morning-restricted & autobot governed), enshrine the Anti-Pattern Ban & Alternative Reference Matrix in `.agents/skills/gitmap/skill.md`, synchronize across repositories, and release GitMap `v6.640.0`.

---

## 2. Subtasks Breakdown

### Task 01: Architecture & Spec Authoring [completed]
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 01`
- **Subtask File:** `.ai-memory/plans/subtasks/367-sync-script-versioning-and-gitmap-skill-enhancement/01-architecture-and-spec-authoring.md`
- **Description:** Author 01-architecture-spec.md, 02-cat-enhancements-spec.md, 03-sync-script-dual-versioning-spec.md, 04-gitmap-skill-and-anti-patterns-spec.md, and register in master plan indices.

### Task 02: `gitmap cat` Flag & Query Filtering Enhancements [completed]
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 02`
- **Subtask File:** `.ai-memory/plans/subtasks/367-sync-script-versioning-and-gitmap-skill-enhancement/02-cat-enhancements.md`
- **Description:** Implement `-t` flag parsing in `cat_options.go`, combined query-window filtering in `cat_reader.go`, and inline anti-pattern detection in `py_cmd.go`.

### Task 03: Coding-Guidelines Sync Script Dual Versioning (V6 vs V7) [completed]
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 03`
- **Subtask File:** `.ai-memory/plans/subtasks/367-sync-script-versioning-and-gitmap-skill-enhancement/03-sync-script-dual-versioning.md`
- **Description:** Verify and ensure dual execution parity between `38-sync-prompts-skills-scripts.py` (V6 unrestricted) and `38-sync-prompts-skills-scripts-v7.py` (V7 restricted) across `gitmap` and `coding-guidelines`.

### Task 04: GitMap Skill Enhancement & Anti-Pattern Ban Matrix [completed]
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 04`
- **Subtask File:** `.ai-memory/plans/subtasks/367-sync-script-versioning-and-gitmap-skill-enhancement/04-gitmap-skill-and-anti-patterns.md`
- **Description:** Update `.agents/skills/gitmap/skill.md` with explicit restrict condition bans, the full 7-item anti-pattern ban table, and side-by-side executable GitMap commands.

### Task 05: Fleet Synchronization Across Connected Repositories [completed]
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 05`
- **Subtask File:** `.ai-memory/plans/subtasks/367-sync-script-versioning-and-gitmap-skill-enhancement/05-fleet-sync.md`
- **Description:** Synchronize skills and scripts across connected repositories using the fleet sync engine.

### Task 06: Verification, Minor Bump Release, and Completion Matrix [completed]
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 06`
- **Subtask File:** `.ai-memory/plans/subtasks/367-sync-script-versioning-and-gitmap-skill-enhancement/06-verification-and-release.md`
- **Description:** Run unit tests, bump version to `v6.640.0`, commit via `gitmap cpr`, push upstream, and output completion matrix.
