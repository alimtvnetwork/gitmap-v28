# Plan 368: SQL Query Execution, AGY Database Commands, AI Code Execution Restrictions & Heatmap Integration

- **Plan Slug:** `sql-query-execution-and-gitmap-skill-enhancement`
- **Root Slug:** `sql-query-execution-and-gitmap-skill-enhancement`
- **Specification:** [02-spec/21-app/368-sql-query-execution-and-gitmap-skill-enhancement/readme.md](../../../02-spec/21-app/368-sql-query-execution-and-gitmap-skill-enhancement/readme.md)
- **Status:** `pending`
- **Created At:** 31-Oct-26
- **Architecture Lead:** V7 Master Autonomous Execution Engine

---

## 1. Goal

Permanently eliminate ad-hoc Python SQLite queries (e.g. `check_summaries.py`, `python -c "import sqlite3..."`), scratch script creation, and unnotified file analysis by implementing native `gitmap agy summaries` (and `gitmap agy convs`), registering built-in `agy` database alias in `gitmap db-read`, expanding `gitmap py` interceptor warnings, creating the centralized `.agents/rules/ai_code_execution_restrictions.md`, mandating Heatmap touch notifications in skills, updating Letterly prompt standards, synchronizing across 43 fleet repositories, and releasing GitMap `v6.642.0`.

---

## 2. Subtasks Breakdown

### Task 01: Architecture, Spec Authoring & 4-Part RCA [in_progress]
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 01`
- **Subtask File:** `.ai-memory/plans/subtasks/368-sql-query-execution-and-gitmap-skill-enhancement/01-architecture-and-rca.md`
- **Description:** Author 01-architecture-and-rca-spec.md, 02-agy-db-commands-and-alias-spec.md, 03-restrictions-and-letterly-prompts-spec.md, 04-heatmap-file-touch-spec.md, and register in master plan indices.

### Task 02: Dedicated `gitmap agy summaries` CLI & Built-in `agy` DB Alias [pending]
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 02`
- **Subtask File:** `.ai-memory/plans/subtasks/368-sql-query-execution-and-gitmap-skill-enhancement/02-agy-db-commands-and-alias.md`
- **Description:** Implement `cli/cmdagy/agy_summaries.go` to parse and render `conversation_summaries.db` in Tokyo Night table or JSON, wire `gitmap agy summaries` / `convs` / `db`, and register the `agy` alias in `gitmap db-read`.

### Task 03: Centralized AI Code Execution Restrictions Rule & Letterly Standards [pending]
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 03`
- **Subtask File:** `.ai-memory/plans/subtasks/368-sql-query-execution-and-gitmap-skill-enhancement/03-restrictions-and-letterly-standards.md`
- **Description:** Create `.agents/rules/ai_code_execution_restrictions.md`, update `linter-scripts/check-letterly-format.py`, and eliminate all ad-hoc script loopholes from prompt templates.

### Task 04: GitMap Py Interceptor Expansion & Heatmap File Touch Logging [pending]
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 04`
- **Subtask File:** `.ai-memory/plans/subtasks/368-sql-query-execution-and-gitmap-skill-enhancement/04-py-interceptor-and-heatmap-touch.md`
- **Description:** Expand `warnAdHocAntiPatterns` in `cli/cmdpy/py_cmd.go` to detect `conversation_summaries.db`, `sqlite3.connect`, and verify `gitmap heatmap touch` command for file analysis logging.

### Task 05: GitMap Skill Updates & Cross-Repository Fleet Synchronization [pending]
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 05`
- **Subtask File:** `.ai-memory/plans/subtasks/368-sql-query-execution-and-gitmap-skill-enhancement/05-skill-updates-and-fleet-sync.md`
- **Description:** Update `.agents/skills/gitmap/skill.md` with explicit database query substitutions and heatmap analysis rules, and execute fleet sync across connected repositories.

### Task 06: Verification, Minor Bump Release `v6.642.0` & Completion Matrix [pending]
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 06`
- **Subtask File:** `.ai-memory/plans/subtasks/368-sql-query-execution-and-gitmap-skill-enhancement/06-verification-and-release.md`
- **Description:** Run unit tests, compile `gitmap.exe` with memory guards, deploy to `%LOCALAPPDATA%\gitmap-cli\`, commit via `gitmap cpr`, push upstream, and present completion report.
