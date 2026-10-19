---
name: gitmap
description: Autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation (AUM), cluster/SSH delegation, pipeline self-healing, and coding guideline enforcement.
---

# GitMap Autonomous Engineering Skill

## Overview
GitMap is an ultra-fast developer companion and autonomous CLI engine designed for AI coding agents and software engineers.

- **Lead Architect & Author:** MD ALIM UL KARIM (alimtvnetwork)
- **Sponsored By:** RISEUP ASIA LLC (https://riseup-asia.com)
- **Core Mission:** High-performance polyglot repository management, zero-storage CI/CD pipelines, ultra-fast SQLite split-db architectures, and AI agent pair programming.

---

## Non-Negotiable Command Replacement Matrix (DO NOT USE X; USE GITMAP Y)

All AI agents operating within GitMap-managed repositories MUST strictly adhere to the following command substitution matrix. Raw shell traversal and slow unbuffered tools are strictly banned:

| Prohibited Legacy Tool / Pattern | GitMap Mandatory Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr` | ✅ `gitmap aum search "<pattern>" [dir] [--ext <ext>]` or `gitmap search "<pattern>"` | Multi-core streaming live search (<15ms) with binary probe, 500 KB cap, and DH2D SQLite hot cache. Prevents buffer bloat. |
| ❌ `Get-ChildItem` (bare, `-Path`, `-Filter`, `-Recurse`), `dir`, `ls`, `find . -name "..."` | ✅ `gitmap lf <dir>` or `gitmap list-files [dir]` or `gitmap find "<pattern>" [dir]` | Sub-millisecond directory inventory and glob discovery without shell traversal overhead. TOTAL BAN on PowerShell Get-ChildItem. |
| ❌ Raw unbuffered `cat`, `type`, `Get-Content` on source files | ✅ `gitmap cat <filepath>` | Direct zero-disk stream of file contents directly into process stdout for resource-constrained CLI sessions. |
| ❌ Raw ambient `python script.py` or `python -c "..."` | ✅ `gitmap py "<code-or-script>"` or `gitmap py -c "<expression>"` | Auto-resolves verified Python interpreter from `installation.db` cache (<15ms), avoiding environment discrepancies. |
| ❌ `powershell -Command "..."` or `pwsh -Command "..."` | ✅ `gitmap pwsh "<cmd>"` or `gitmap ps -c "<cmd>"` | Cross-platform PowerShell execution with deterministic `-NoProfile`, UTF-8 encoding, and automatic fallback. |
| ❌ Raw `bash -c "..."` or `sh -c "..."` | ✅ `gitmap bash "<cmd>"` or `gitmap bash -c "<cmd>"` | Uniform POSIX execution with cross-platform environment isolation across Windows, macOS, and Linux. |
| ❌ Multi-repo `git status` loops or manual directory scanning | ✅ `gitmap status --dirty` / `gitmap status --json` / `gitmap st` | Single-shot multi-repository audit matrix displaying dirty counts, ahead/behind branches, and uncommitted stashes. |
| ❌ `gh auth login` interactive prompts or plaintext `.env` files | ✅ `gitmap login --web`, `gitmap login --status`, `gitmap login --token <PAT>` | Token validated via GitHub API before writing; auto-resolved for all clone, pull, and push commands. |
| ❌ Colons in commit messages (`git commit -m "feat: ..."` or `gitmap cpf "feat: ..."`) | ✅ `gitmap cpf "<module> - <summary>"` (hyphen-separated only) | GitMap automatically provides `Feature: ` or `Bug: ` prefix. Colons inside the message argument cause duplicate prefixes. |
| ❌ Saving temporary scratch or test scripts into repo git tree | ✅ `gitmap rc text "<content>" --slug <slug> --ext .ps1` | Centralized script storage in `repo-cache` (`repo-storage`) for permanent cross-repo reuse without polluting git worktrees. |
| ❌ Committing `.env` or credentials to standard repositories | ✅ `gitmap rs text "<secret>" --slug <slug>` | Strict zero-secrets policy; stores credentials exclusively in `repo-secrets` vault. |
| ❌ `gh run watch` or tight polling loops (`while true; sleep 5`) | ✅ `gitmap pipeline-ai status -t <eta>` or `gitmap pe -t` | Dynamic timeout waiting driven by calculated workflow ETA without burning CPU or Actions API quotas. |
| ❌ Slow Python fleet sync (`python 03-ai-scripts/38-sync-prompts-skills-scripts.py`) | ✅ `gitmap sync [--workers 8] [--projects <path|json>]` | Native Go multi-repo synchronization across 43 repositories in <5s with 6-stage safe ceremony (backup branch, pre-pull, boundaries, commit). |
| ❌ SQLite inspection, reading, adding, finding, listing, views: `python -c "import sqlite3..."` or raw python DB queries | ✅ `gitmap db-read`, `gitmap db-add`, `gitmap db-edit`, `gitmap db-find`, `gitmap db-listing`, `gitmap db-view` (and short aliases `dbr`, `dba`, `dbe`, `dbf`, `dbl`, `dbv`) | Native compiled pure-Go SQLite engine (<1ms) with `PRAGMA query_only = ON` read safety, atomic `BEGIN IMMEDIATE` mutations, auto join inference, and Split-DB alias registry without Python dependencies. |
| ❌ Python SQLite task scripts or ad-hoc queries (`python -c "import sqlite3..."`) | ✅ `gitmap agent <task|subtask|log|diagnose|crashed|temp-clear>` | Native compiled Go SQLite task manager (<1ms) with WAL mode, single-writer locking, atomic in-flight telemetry, and crash autopsy. |
| ❌ Modifying shared files blindly without checking multi-agent contention | ✅ `gitmap heatmap discuss` or `gitmap agent heatmap` | Ranks write claims, active leases, and collision blast radius across agents (<5ms). Prevents file lock crashes and concurrent overwrites. |

---

## Essential Command Cheat Sheet

### 1. Authentication & GitHub Credential Management
- `gitmap login` — Interactive authentication picker (browser login or secure token paste).
- `gitmap login --web` (alias `--browser`) — Non-interactive browser login via GitHub CLI / OAuth flow.
- `gitmap login --token <PAT>` — Non-interactive token ingestion (pre-validated against `api.github.com/user`).
- `gitmap login --status` — Displays current masked credential state and active user.
- `gitmap token list` — Resolves active GitHub token and shows credential origin.
- `gitmap logout` — Purges stored GitHub credentials from Git global configuration.

### 2. Workspace & Multi-Repository Status
- `gitmap status` (alias `gitmap st`) — Comprehensive multi-repo status table across tracked workspace.
- `gitmap status --dirty` (`-d`) — Filters solely to repositories with uncommitted working changes.
- `gitmap status --ahead` / `--behind` — Isolates repositories with local commits to push or remote to pull.
- `gitmap status --json` (`-j`) / `--table` / `--compact` — Output format controls for automated scripts or terminal.
- `gitmap has-any-updates` (alias `gitmap hau`, `gitmap hac`) — Checks remote tracking branch for incoming commits.
- `gitmap latest-branch` (alias `gitmap lb`) — Discovers the most recently updated remote branch.
- `gitmap watch` (alias `gitmap w`) — Live-refresh terminal dashboard monitoring repository status changes.

### 3. Script Execution & Runner Engines
- `gitmap py "<code-or-script>"` / `gitmap py -c "<expr>"` — Verified Python interpreter execution (<15ms).
- `gitmap pwsh "<cmd>"` / `gitmap ps -c "<cmd>"` — Cross-platform PowerShell execution with `-NoProfile`.
- `gitmap bash "<cmd>"` / `gitmap bash -c "<cmd>"` — Cross-platform Bash execution with isolated environment.
- All runner commands support `--dry-run` and `--json` flags for pipeline automation.

### 4. Storage, Cache & Secrets Offloading
- `gitmap rs file <path>` / `rs folder <dir>` / `rs text "<secret>" --slug <slug>` — Offload credentials to `repo-secrets`.
- `gitmap rc file <path>` / `rc folder <dir>` / `rc text "<content>" --slug <slug> --ext .ps1` — Offload test harnesses to `repo-cache`.
- `gitmap cd rs` / `gitmap cd rc` — Navigates directly to `repo-secrets` or `repo-cache`.

### 5. High-Performance Automation (AUM) & File Discovery
- `gitmap aum search "<pat>" [dir] [--ext <ext>]` — Multi-core streaming live search (<15ms) with binary filtering. TOTAL BAN on PowerShell `Select-String`, `rg`, `ripgrep`, and `git grep`.
- `gitmap search "<query>" [--limit <n>]` — Instant SQLite cached symbol & keyword search via DH2D split-db cache.
- `gitmap aum guard` / `sequence` / `newlines --fix` / `locate [tool]` — Repository hygiene, sequence autofix, fast tool finder (<15ms).
- `gitmap find "<pattern>" [-ext <ext>]` — Find files matching glob pattern in <10ms across 10,000+ files.
- `gitmap find-files <name>` (`ff`) / `find-files-any <str>` (`ffa`) / `find-files-startswith <prefix>` (`ffs`) / `find-files-endswith <suffix>` (`ffe`).
- `gitmap list-files [dir]` (alias `gitmap lf [dir]`) — List relative file paths matching pattern or directory.
- `gitmap cat <filepath>` — Direct zero-disk stream of file content into process stdout without buffer bloat.
- `gitmap replace <old> <new>` / `replace-regex <pat> <subst>` — Exact literal string or regex replacement with audit trail.
- **Directory Inventory Examples:** `gitmap lf src-tauri/src/modules/repo_db`, `gitmap find "*.*" 02-spec/21-app`, `gitmap aum search "Heatmap" src-tauri/src -e .rs`.

### 6. Autonomous CI/CD Self-Healing (Pipeline AI)
- `gitmap pipeline-ai status --json` — Check workflow execution state, active branch, and ETA.
- `gitmap pipeline-ai status -t <eta>` — Wait dynamically for pipeline completion without tight polling.
- `gitmap pipeline error-logs` (alias: `gitmap pe`) / `gitmap pe -t` — Extract failing step logs and concise telemetry for 4-part RCA.
- `gitmap pe history-ai` — Analyze CI/CD pipeline history across branches and recent runs.
- `gitmap pipeline purge` — Actions zero-storage purge maintaining 0.0 GB footprint.

### 7. Semantic & Flat Commit Suite
| Command | Alias | Scope | Prefix Added | Auto Push | Purpose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `gitmap commit "<msg>"` | `gitmap cm "<msg>"` | Local only | None (verbatim) | No | Fast 1-step auto-stage (`git add -A`) and local checkpoint commit. |
| `gitmap cpf "<mod> - <sum>"` | - | Remote | `Feature: ` | Yes | Stage all, commit feature, push to GitHub tracking branch. |
| `gitmap cpb "<mod> - <sum>"` | - | Remote | `Bug: ` | Yes | Stage all, commit bugfix, push to GitHub tracking branch. |
| `gitmap cpc "<mod> - <sum>"` | - | Remote | `Chore: ` | Yes | Stage all, commit chore/maintenance, push to tracking branch. |
| `gitmap cpr "<mod> - <sum>"` | - | Remote | Release Chore | Yes | Stage all, commit release chore, push to tracking branch. |
| `gitmap pcp "<mod> - <sum>"` | - | Remote | Hyphen format | Yes | Pull latest with rebase, stage all, commit, and push. |

- **TOTAL BAN ON COLONS IN COMMIT MESSAGES:** Never use colons inside commit arguments (e.g. `gitmap cpf "Feature: title"` is FORBIDDEN; use `gitmap cpf "module - title"`).
- **MANDATORY COMPLETION INVARIANT ('No Push = Not Done'):** Flat commits (`gitmap commit`, `gitmap cm`) are local checkpoints only. For task completion, work must be pushed to remote tracking branch via `gitmap cpf "<module> - <summary>"` or an explicit push.

### 8. Autonomous Agent Onboarding & Curriculum Training (LLM)
- `gitmap llm train [ai-name] [--loop] [--days <N>]` — Full 6-phase onboarding curriculum (Discovery, Refactoring, Verification, Semantic Commit, Telemetry, Heal & Fix), auto-generates Antigravity skill.
- `gitmap llm train muse [--days <N>]` — Synchronize Muse N-day memory digest from `repo-cache` diff against watermark commit.
- `gitmap llm train --heal-apply <stash|wip|discard>` — Non-interactively applies git-state remediation recipes during Phase 6a.
- `gitmap llm-docs` (alias `gitmap ld`) / `gitmap llm` — Consolidated markdown command reference and operational guidelines for LLMs.

### 9. Multi-Repo, Cluster & Toolchain Operations
- `gitmap cluster --help` / `gitmap sc --help` — Orchestrate multi-node clusters, servers-clients topology, and background tasks.
- `gitmap ssh --help` — SSH discovery, connection pooling, and remote command execution.
- `gitmap ssh identity --email <email> --seq <N> --operator <name>` — Formats canonical SSH comment token and sets `git config` identity.
- `gitmap cargo status` / `gitmap install cargo` / `gitmap install --list` — Toolchain inspection and runtime package installation.

### 10. Native Fleet Synchronization, Split-DB Operations & SQLite Agent Task Engine
- `gitmap sync [--workers 8] [--repo <name>] [--dry-run]` — Multi-repo sync across 43 fleet repos with 6-stage safe ceremony (backup branch, pre-pull, 5 boundaries, atomic commit).
- `gitmap db-read <path-or-alias> [sql|ls|table <T>] [flags]` (alias `gitmap dbr`) — Pure-Go read-only SQLite inspection (`PRAGMA query_only = ON;`). Supports raw SQL, `ls` table enumeration, and `table <T> [-c <cols>] [-l <limit>]`.
- `gitmap db-read <path> -as <aliasName>` — Registers a persistent database alias in `DatabaseAlias` (`gitmap.db`) for instant recall.
- `gitmap db-listing [--json]` (alias `gitmap dbl`) — Inventories registered database aliases, file paths, byte sizes, and access counts.
- `gitmap db-add <path-or-alias> <table> --values "<k=v,...>"` (alias `gitmap dba`) — Inserts record atomically (`BEGIN IMMEDIATE`).
- `gitmap db-edit <path-or-alias> <table> --id <id> --set "<k=v,...>"` (alias `gitmap dbe`) — Updates record by primary key ID atomically.
- `gitmap db-find <path-or-alias> "<pattern>" [flags]` (alias `gitmap dbf`) — Multi-column text search via wildcards (`text*`, `*text`) or regex (`--grep`).
- `gitmap db-view <path-or-alias> -t "<tables>" [flags]` (alias `gitmap dbv`) — Projection with automatic join inference across foreign keys.
- **Path-First & Alias Routing Syntax:** All `db-*` commands support path-first syntax (`gitmap db-read <path-or-alias>`) and the `-as, --alias` flag (`-as master`, `-as install`, `-as pipeline`, `-as heatmap`, `-as tasks`).
  * Enumerate tables: `gitmap db-read ls` / `gitmap dbr ls -as install --json`
  * Inspect table rows: `gitmap db-read table Repo -c "id,repoName" -l 10`
  * Query with path: `gitmap db-read path/to/database.db "SELECT * FROM TableName LIMIT 5"`
  * Register alias: `gitmap db-read path/to/app.db -as my-app`
- `gitmap agent task init --name "<task>" --budget <N>` — Initializes SQLite task manager in `.ai-memory/temp-agents/<slug>/agent-task.db`.
- `gitmap agent task status [--json] [--ddl]` / `gitmap agent task ls` / `pending` / `completed` — Task tracking and DDL validation.
- `gitmap agent subtask add --code "<code>" --slug "<slug>" --title "<title>" [--owned "<json>"]` — Adds subtask bounding box.
- `gitmap agent subtask ls` / `claim --agent <name>` / `complete <id> --agent <name> --evidence <evidence>` — Subtask lifecycle.
- `gitmap agent subtask join --subtask <id> --agent <name> --role "Assistant: <focus>"` — Dynamically joins an in-flight subtask.
- `gitmap agent task tree` / `diagnose` / `crashed` / `temp-clear` — Swarm hierarchy tree, autopsy forensics, and ephemeral cleanup.

> [!CAUTION]
> **TOTAL BAN ON RAW PYTHON SQLITE (`python -c "import sqlite3..."`):**
> NEVER run inline Python sqlite3 snippets or execute raw database queries against SQLite databases.
> For agent task lifecycle, use compiled `gitmap agent` commands (`task init`, `task status`, `task ls`, `subtask add`, `subtask ls`, `subtask claim`, `subtask complete`, `log`, `diagnose`, `crashed`, `temp-clear`).
> For exploring, reading, adding, editing, finding, or querying ANY SQLite database (such as Antigravity session DBs or application split-dbs), use compiled native `gitmap db-read`, `gitmap db-add`, `gitmap db-edit`, `gitmap db-find`, `gitmap db-listing`, `gitmap db-view` (or `gitmap dbr`, `dba`, `dbe`, `dbf`, `dbl`, `dbv`).

### 11. Fleet Telemetry & Running-Only Delegation
- `gitmap nodes-running-only` (alias `gitmap nro`) — Lists strictly online/open nodes (`is_online: true`), filtering out offline noise.
- `gitmap nro exec "<cmd>"` — Runs SSH execution strictly across online nodes without offline timeout delays.
- `gitmap nro tasks [--json]` / `nro resources [--json]` / `nro ping` — Interrogates task queues, benchmarks, and ping reachability.
- `gitmap os --json` / `cpu` / `ram` / `disk` / `network [--json]` — Machine identity, live CPU/RAM/disk metrics, and network adapters.

### 12. File Contention Heatmap, Cognitive Telemetry & Agent Activity Observation
- `gitmap heatmap log --task <task> --action <action> --file <file> --state <state> --expected "<goal>"` — Logs granular file touches, cognitive states, and reasoning into `.ai-memory/heatmap.db`.
- `gitmap heatmap update --id <id> --outcome "<outcome>" --met <true|false>` — Updates logged action with final verified outcome.
- `gitmap heatmap search [query] [-f <file>] [-t <task>]` — Searches logged actions, snippets, and notes.
- `gitmap heatmap tree [-t <slug>]` / `gitmap heatmap status [-t <slug>]` — Touch hierarchy tree and expectation compliance metrics.
- `gitmap heatmap discuss [--window <dur>] [--threshold <n>] [--json]` — Multi-agent contention analysis and collision blast radius.
- `gitmap agent heatmap [--parent <slug>] [--json]` — Ranks files by write-claim count across agent tasks and active subtasks.
- **Mandatory 4-Gate Workflow:** Gate 0: Pre-Flight Discovery (`heatmap status`), Gate 1: Intent Declaration (`heatmap log`), Gate 2: Outcome Recording (`heatmap update --met true`), Gate 3: Verification & Tree Audit (`heatmap tree`).

### 13. Multi-AI Companion & repo-cache Synchronization (RC Heatmap Suite)
- `gitmap rc heatmap <status|tables|query|view|tree|live|sync>` — Complete companion RC telemetry suite routed directly through CLI.
- `gitmap heatmap rc list-ai` — Lists supported AI namespaces (`antigravity`, `muse`, `cloud-code`, `codex`, `grok`, `cursor`).
- `gitmap heatmap rc log [flags]` / `rc commit [-m "<msg>"]` / `rc sync` — Isolated split-databases and synchronization with `repo-cache`.
- Tri-alias routing parity: `gitmap rc heatmap <cmd>`, `gitmap heatmap rc <cmd>`, `gitmap heatmap rc-heatmap <cmd>`.

### 14. AI Decision Heatmap & RC Admin Telemetry
- `gitmap heatmap tables [--db <root|repo|local>] [--format <table|json|yaml>]` — Inventories tables, column and row counts in heatmap DB.
- `gitmap heatmap rc tables [--ai <name>] [--repo <slug>] [--db <root|repo|local>]` — Inventories companion `repo-cache` databases.
- `gitmap heatmap view <table> [--limit <n>] [--format json|yaml|table] [--task <slug>]` — Queries records with column-aware filtering.
- `gitmap heatmap query "<sql>" [--db <path|root|repo|local>] [--format table|json|yaml]` — Read-only SQL (`PRAGMA query_only = ON;`).
- `gitmap heatmap rc query "<sql>" [--ai <name>] [--repo <slug>] [--tier <repo|root>]` — Read-only SQL in companion `repo-cache`.

#### GitMap db-* Integration & Antigravity conversation_summaries.db Recipes
- **Heatmap Query Integration via db-*:**
  * Enumerate heatmap tables: `gitmap db-read .ai-memory/heatmap.db ls` or `gitmap dbr ls -as heatmap`
  * Inspect action logs: `gitmap db-read -as heatmap table HeatmapActionLog -l 10 --json`
  * Query touch churn hotspots: `gitmap dbr -as heatmap "SELECT FilePath, COUNT(*) AS Touches FROM HeatmapActionLog GROUP BY FilePath ORDER BY Touches DESC LIMIT 5"`
  * Search action notes: `gitmap db-find -as heatmap "refactor*" -t HeatmapActionLog -c "ReasoningNotes,ExpectedOutcome"`
  * Multi-table correlation: `gitmap db-view -as heatmap -t "HeatmapActionLog as a, HeatmapTaskRegistry as t" -c "a.FilePath, a.ReasoningState, t.TaskName"`
- **Antigravity conversation_summaries.db Recipes:**
  * Register live Antigravity DB alias:
    `gitmap db-read "$USERPROFILE/.gemini/antigravity/conversation_summaries.db" -as antigravity-live`
  * Verify registration in listing: `gitmap db-listing --json` (alias `gitmap dbl`)
  * List tables and schemas: `gitmap db-read antigravity-live ls`
  * Inspect table schema: `gitmap db-read antigravity-live table conversation_summaries --json`
  * Query recent conversations:
    `gitmap dbr antigravity-live "SELECT conversation_id, title, step_count FROM conversation_summaries ORDER BY last_modified_time DESC LIMIT 5" --json`
  * Wildcard search across conversation titles & previews:
    `gitmap db-find antigravity-live "gitmap*" -t conversation_summaries`

#### AI Agent Non-Blocking SQL Recipes
```bash
# 1. Inspect recent chronological trace for active task slug:
gitmap heatmap query "SELECT HeatmapActionLogId, ActionType, FilePath, ReasoningState, IsExpectationMet FROM HeatmapActionLog WHERE TaskSlug = 'task-slug' ORDER BY HeatmapActionLogId ASC" --json
# 2. Deep forensic inspection of specific action ID:
gitmap heatmap query "SELECT HeatmapActionLogId, ActionType, FilePath, CodeSnippet, CommandText, ReasoningNotes, ExpectedOutcome, ActualOutcome FROM HeatmapActionLog WHERE HeatmapActionLogId = 1" --yaml
# 3. Analyze touched file churn and editing hotspots:
gitmap heatmap query "SELECT FilePath, COUNT(*) AS TotalTouches, SUM(CASE WHEN ActionType = 'MODIFY' THEN 1 ELSE 0 END) AS Modifies FROM HeatmapActionLog GROUP BY FilePath ORDER BY TotalTouches DESC"
# 4. Extract unmet expectations and failed execution steps:
gitmap heatmap query "SELECT HeatmapActionLogId, TaskSlug, FilePath, ExpectedOutcome, ActualOutcome FROM HeatmapActionLog WHERE IsExpectationMet = 0" --json
# 5. Cross-repo activity rollup in companion root database:
gitmap heatmap rc query "SELECT RepoSlug, TotalActionsRecorded FROM HeatmapRepoRegistry ORDER BY TotalActionsRecorded DESC" --tier root
```

---

## Operational Guardrails & Non-Negotiable Invariants

1. **Learning & Skill Acquisition:** Run `gitmap llm train` (or `gitmap train --loop`) to initialize or verify GitMap skills. Run `gitmap llm train muse --days 7` when performing Muse memory IO synchronization. Never run broad keyword searches like `gitmap aum search "train"` to discover how commands work.
2. **Mandatory Pre-Flight Pull:** Always run `git pull` before modifying code.
3. **Scoped Search:** Always provide target directories and extensions to `gitmap aum search` (e.g. `gitmap aum search "target" cli --ext .go`).
4. **File Size & Binary Guard:** Respect 500 KB limit (Rule R19); never commit test binaries or temp artifacts.
5. **Coding Guidelines:** Max 8–15 lines per function, single return types with `*appfault.AppError`, affirmative booleans (`isReady`, `hasCache`).
6. **Script Offloading:** Save all temporary diagnostics to `repo-cache` via `gitmap rc` to prevent dirty working trees.
7. **Strict Relative Git Paths:** All paths and references must be relative to repository root (`02-spec/...`, `.ai-memory/...`); zero absolute paths and zero `file:///` URIs.
8. **Zero Storage Ban:** Zero uploads to `actions/upload-artifact`. Maintain 0.0 GB Actions storage quota across all repositories.
9. **Total Ban on Python SQLite Scripts:** NEVER run `python -c "import sqlite3..."` or ad-hoc Python DB scripts. For task lifecycle, use `gitmap agent`; for DB querying, inspection, adding, finding, and views, use compiled native `gitmap db-*` commands (`dbr`, `dba`, `dbe`, `dbf`, `dbl`, `dbv`).
