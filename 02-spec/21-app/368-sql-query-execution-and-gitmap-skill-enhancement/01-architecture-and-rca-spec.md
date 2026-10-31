# 368: SQL Query Execution, AGY Database Commands, AI Code Execution Restrictions & Heatmap Integration

## Executive Summary
Autonomous AI agents frequently revert to writing ad-hoc Python scripts (e.g. `check_summaries.py`, `python -c "import sqlite3..."`) to query local SQLite databases such as Antigravity's `conversation_summaries.db`, and fail to notify GitMap's Heatmap system when inspecting or analyzing codebase files. This specification details the root causes, architecture, native CLI additions, restriction governance, and fleet synchronization required to permanently eliminate these anti-patterns.

---

## 1. Grounded Root Cause Analysis (4-Part RCA)

### 1.1 First Principles Discovery
- **Observation:** Agents in workspaces such as `Antigravity-Manager` execute ad-hoc Python scripts (`gitmap py C:\...\scratch\check_summaries.py` or inline `python -c "import sqlite3..."`) targeting `~/.gemini/antigravity/conversation_summaries.db` to inspect conversation states (`SELECT conversation_id, title, status, not_fully_idle, last_modified_time...`).
- **Observation:** Agents read and analyze critical source files across repositories without registering their touch or read intent with GitMap Heatmap, causing multi-agent blind spots.

### 1.2 Underlying Root Causes
1. **Lack of Dedicated AGY Conversation Database Command:** Although `gitmap db-read` natively supports custom DB paths, agents had no high-level command like `gitmap agy summaries` (or `gitmap agy convs`) that automatically resolves `~/.gemini/antigravity/conversation_summaries.db` and renders Tokyo Night formatted tables.
2. **Missing `agy` Alias in DB Subsystem:** `gitmap db-read` supported aliases (`master`, `tasks`, `pipeline`, `heatmap`), but did not include an auto-registered `agy` alias pointing to `~/.gemini/antigravity/conversation_summaries.db`.
3. **No Centralized AI Code Execution Restriction File:** Guidelines were dispersed across multiple markdown files without a single sovereign rule file (`.agents/rules/ai_code_execution_restrictions.md`) explicitly banning ad-hoc Python SQLite scripts and detailing mandatory native replacements.
4. **Letterly Prompt Missing Strict SQL & Heatmap Mandates:** Letterly prompt generators still permitted ambiguity around `gitmap py` and lacked explicit mandates to log file analysis into GitMap Heatmap (`gitmap heatmap touch <path> --reason "<why>"`).

### 1.3 Architectural Remedy
- **Native AGY Command:** Add `gitmap agy summaries [N]` (aliases: `gitmap agy convs`, `gitmap agy db`) to directly read and format `conversation_summaries.db` in sub-millisecond time.
- **Built-in `agy` Alias:** Register `agy` in `db-read` so `gitmap db-read agy "SELECT ..."` works out of the box.
- **Python Interceptor:** Expand `warnAdHocAntiPatterns` in `cli/cmdpy/py_cmd.go` to intercept `conversation_summaries.db` or ad-hoc `sqlite3.connect` calls and display the exact `gitmap agy summaries` replacement.
- **Heatmap Touch Command:** Ensure `gitmap heatmap touch <path> [--reason <text>]` (or `gitmap cat -t / --touch`) registers file analysis in `heatmap.db` and displays live telemetry.
- **Centralized Restrictions Rule:** Author `.agents/rules/ai_code_execution_restrictions.md` and reference it across all skills and prompt templates.
- **Fleet Sync:** Propagate updated skills and restriction rules across all connected repositories.
