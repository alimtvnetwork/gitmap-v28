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
| ❌ `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr` | ✅ `gitmap aum search "<pattern>" [dir] [--ext <ext>]` or `gitmap search "<pattern>"` | Multi-core streaming live search (<15ms) with binary null-byte probe, 500 KB file cap, and DH2D SQLite hot cache. Prevents buffer bloat and terminal freezing. |
| ❌ `Get-ChildItem -Recurse`, `find . -name "..."`, `dir /s /b` | ✅ `gitmap find "<pattern>" [-ext <ext>]` or `gitmap list-files [dir]` | Sub-millisecond indexed glob discovery across 10,000+ files without disk traversal overhead. |
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
| ❌ Slow Python fleet sync (`python 03-ai-scripts/38-sync-prompts-skills-scripts.py`) | ✅ `gitmap sync [--workers 8] [--projects <path|json>]` | Native Go multi-repo synchronization across 43 repositories in <5s with 6-stage safe ceremony (backup branch, pre-pull, 5 boundaries, atomic commit). |
| ❌ Python SQLite task scripts or ad-hoc queries (`python -c "import sqlite3..."`) | ✅ `gitmap agent <task|subtask|log|diagnose|crashed|temp-clear>` | Native compiled Go SQLite task manager (<1ms) with WAL mode, single-writer locking, atomic in-flight telemetry, and crash autopsy. |
| ❌ Modifying shared files blindly without checking multi-agent contention | ✅ `gitmap heatmap discuss` or `gitmap agent heatmap` | Ranks write claims, active leases, and collision blast radius across agents (<5ms). Prevents file lock crashes and concurrent overwrites. |

---

## Essential Command Cheat Sheet

### 1. Authentication & GitHub Credential Management
- `gitmap login` — Interactive authentication picker (browser login or secure token paste).
- `gitmap login --web` (alias `--browser`) — Non-interactive browser login via GitHub CLI / OAuth flow.
- `gitmap login --token <PAT>` — Non-interactive token ingestion (pre-validated against `api.github.com/user` before saving to global git config).
- `gitmap login --token <PAT> --no-verify` — Offline token ingestion without network probe.
- `gitmap login --status` — Displays current masked credential state and active user.
- `gitmap token list` — Resolves active GitHub token and shows credential origin.
- `gitmap logout` — Purges stored GitHub credentials from Git global configuration.

### 2. Workspace & Multi-Repository Status
- `gitmap status` (alias `gitmap st`) — Comprehensive multi-repo status table across tracked workspace.
- `gitmap status --dirty` (`-d`) — Filters solely to repositories with uncommitted working changes.
- `gitmap status --ahead` — Isolates repositories with local commits waiting to push.
- `gitmap status --behind` — Isolates repositories with remote commits waiting to pull.
- `gitmap status --json` (`-j`) — Emits structured JSON array of repository states for automated scripts.
- `gitmap status --table` — Forces rich ANSI status table with colored columns.
- `gitmap status --compact` — Emits single-line compact summary per repository.
- `gitmap has-any-updates` (alias `gitmap hau`, `gitmap hac`) — Checks remote tracking branch for incoming commits.
- `gitmap latest-branch` (alias `gitmap lb`) — Discovers the most recently updated remote branch.
- `gitmap watch` (alias `gitmap w`) — Live-refresh terminal dashboard monitoring repository status changes.

### 3. Script Execution & Runner Engines
- `gitmap py "<code-or-script>"` — High-performance cross-platform Python script execution.
- `gitmap py -c "<code-or-expression>"` — Direct inline Python command evaluation.
- `gitmap pwsh "<cmd>"` / `gitmap ps "<cmd>"` — Cross-platform PowerShell execution with `-NoProfile`.
- `gitmap ps -c "<cmd>"` — Direct inline PowerShell command evaluation.
- `gitmap bash "<cmd>"` / `gitmap sh "<cmd>"` — Cross-platform Bash execution.
- `gitmap bash -c "<cmd>"` — Direct inline POSIX Bash evaluation.
- All runner commands support `--dry-run` and `--json` flags for pipeline automation.

### 4. Storage, Cache & Secrets Offloading
- `gitmap rs file <filepath> [--repo <name>]` — Copies secret file into `repo-secrets` and auto-pushes.
- `gitmap rs folder <folderpath> [--repo <name>]` — Copies secret directory into `repo-secrets` and auto-pushes.
- `gitmap rs text "<secret>" [--slug <slug>]` — Writes sensitive text into `repo-secrets/<repo>/<slug>.txt`.
- `gitmap rc file <filepath> [--repo <name>]` — Copies reusable script into `repo-cache` (`repo-storage`) and auto-pushes.
- `gitmap rc folder <folderpath> [--repo <name>]` — Copies reusable fixture directory into `repo-cache`.
- `gitmap rc text "<content>" --slug <slug> --ext <.ps1|.py>` — Writes test harnesses into `repo-cache` for permanent cross-repo reuse.
- `gitmap cd rs` / `gitmap cd rc` — Navigates directly to `repo-secrets` or `repo-cache`.

### 5. High-Performance Automation (AUM) & File Discovery
- `gitmap aum search "<pattern>" [dir] [--ext <ext>] [-r] [-i]` (alias: `gitmap aum grep`) — Multi-core streaming live search with lazy regex and binary filtering. ALWAYS scope with target `[dir]` and `--ext`. Replaces slow PowerShell `Select-String`, `Get-ChildItem -Recurse`, and `git grep`. TOTAL BAN on PowerShell `Select-String`, `rg`, `ripgrep`, and `git grep`.
- `gitmap search "<query>" [--limit <n>]` — Instant SQLite cached symbol & keyword search across scanned repositories using DH2D split-db hot cache.
- `gitmap aum guard` — Enforces 500 KB limit, large JSON exclusion, and binary null-byte probe.
- `gitmap aum sequence` — Markdown sequence gap detector and `# XX Title` autofixer.
- `gitmap aum exclude list` — Query persistent search exclusions from SQLite.
- `gitmap aum newlines --fix` — Polyglot CRLF to LF and trailing whitespace normalizer.
- `gitmap aum cache status` — Sub-millisecond in-memory cache status.
- `gitmap aum locate [tool]` — Ultra-fast tool finder (<15ms, e.g. `vcvarsall.bat`, `msbuild`, `python`; replaces slow PowerShell traversal).
- `gitmap aum benchmark all` — Side-by-side Go vs Python execution benchmarks.
- `gitmap find "<pattern>" [-ext <ext>]` — Find files matching glob pattern in <10ms across 10,000+ files.
- `gitmap find-files <name>` (alias: `gitmap ff <name>`) — Find exact filename with optional `-ext`.
- `gitmap find-files-any <str>` (alias: `gitmap ffa <str>`) — Find files matching substring.
- `gitmap find-files-startswith <prefix>` (alias: `gitmap ffs <prefix>`) — Find by filename prefix.
- `gitmap find-files-endswith <suffix>` (alias: `gitmap ffe <suffix>`) — Find by filename suffix (e.g. `_test.go`).
- `gitmap list-files [dir]` (alias: `gitmap lf [dir]`) — List relative file paths matching pattern or directory.
- `gitmap cat <filepath>` — Direct zero-disk stream of file content into process stdout. Inspect file content immediately in resource-constrained CLI sessions without buffer bloat.
- `gitmap replace <old> <new>` — Exact literal string replacement with audit trail.
- `gitmap replace-regex <pat> <subst>` — Regex replacement across repository.

### 6. Autonomous CI/CD Self-Healing (Pipeline AI)
- `gitmap pipeline-ai status --json` — Check workflow execution state, active branch, and ETA.
- `gitmap pipeline-ai status -t <eta>` — Wait dynamically for pipeline completion without tight polling.
- `gitmap pipeline error-logs` (alias: `gitmap pe`) — Extract failing step logs to file for 4-part RCA.
- `gitmap pe -t` — Telemetry mode extracting concise failure summaries.
- `gitmap pe history-ai` — Analyze CI/CD pipeline history across branches and recent runs.
- `gitmap pipeline purge` — Actions zero-storage purge maintaining 0.0 GB footprint.

### 7. Semantic & Flat Commit Suite

#### GitMap Commit & Release Commands
| Command | Alias | Scope | Prefix Added | Auto Push | Purpose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `gitmap commit "<message>"` | `gitmap cm "<msg>"` | Local only | None (verbatim) | No | Fast 1-step auto-stage (`git add -A`) and local commit for scratch/micro-steps. |
| `gitmap cpf "<module> - <summary>"` | - | Remote | `Feature: ` | Yes | Stage all, commit feature, and push directly to GitHub tracking branch. |
| `gitmap cpb "<module> - <summary>"` | - | Remote | `Bug: ` | Yes | Stage all, commit bugfix, and push directly to GitHub tracking branch. |
| `gitmap cpc "<module> - <summary>"` | - | Remote | `Chore: ` | Yes | Stage all, commit chore/maintenance, and push directly to GitHub tracking branch. |
| `gitmap cpr "<module> - <summary>"` | - | Remote | Release Chore | Yes | Stage all, commit release chore, and push directly to GitHub tracking branch. |
| `gitmap pcp "<module> - <summary>"` | - | Remote | Hyphen format | Yes | Pull latest with rebase, stage all, commit, and push with preflight verification. |

#### Command Details & Behavioral Guidelines
- `gitmap commit "<message>"` (alias `gitmap cm "<message>"`) — 1-step auto-stage (`git add -A`) and local commit without semantic prefixing ('Feature: ', 'Bug: ') and without remote git push. Ideal for quick local iterations, WIP checkpoints, and subtask milestones.
- `gitmap cpc "<module> - <summary>"` — Stage, commit, and push chore/maintenance branch (GitMap auto-prefixes `Chore: `).
- `gitmap cpf "<module> - <summary>"` — Stage, commit, and push feature branch (GitMap auto-prefixes `Feature: `).
- `gitmap cpb "<module> - <summary>"` — Stage, commit, and push bugfix branch (GitMap auto-prefixes `Bug: `).
- `gitmap cpr "<module> - <summary>"` — Stage, commit, and push release chore.
- `gitmap pcp "<module> - <summary>"` — Pull latest, commit, and push with preflight verification.
- `gitmap pull [repo]` (alias `gitmap p`) — Pull targeted repository.
- `gitmap pull-all` (alias `gitmap pa`) — Pull all repositories in workspace.
- `gitmap fix [repo] [action]` — Apply remediation to repo (aliases: `stash`, `wip`, `discard`).
- `gitmap lowercase` (alias `gitmap lcf`) — Safe 2-step `git mv` file case normalization.
- `gitmap lowercase-readme` — Safe 2-step `git mv` case normalization for root `readme.md`.
- **TOTAL BAN ON COLONS IN COMMIT MESSAGES:** Never use colons inside commit arguments (e.g. `gitmap cpf "Feature: title"` is FORBIDDEN; use `gitmap cpf "module - title"`).
- **MANDATORY COMPLETION INVARIANT ('No Push = Not Done'):** Flat commits (`gitmap commit`, `gitmap cm`) are local checkpoints only. For task completion, AGENTS.md §13 strictly mandates that work must be pushed to the remote tracking branch via `gitmap cpf "<module> - <summary>"` or an explicit remote push.

### 8. Autonomous Agent Onboarding, Curriculum & Memory Training (LLM)
- `gitmap llm train [ai-name] [flags]` (aliases: `gitmap train`, `gitmap llm-train`, `gitmap llm chain`) — Full 6-phase self-looping onboarding curriculum (Discovery, Refactoring, Verification, Semantic Commit, Telemetry, Heal & Fix), auto-generates Antigravity skill, author/sponsor attribution.
- `gitmap llm train muse [--days <N>]` (aliases: `gitmap train muse`, `gitmap train-muse`, `gitmap llm-train-muse`) — Synchronize Muse N-day memory digest from `repo-cache` (fallback `repo-secrets`), diff against watermark commit, output AI traversal checklist for sequential memory reading, and reconstitute `chat-seeds`.
- `gitmap llm train --loop` / `gitmap llm train --self-loop <N>` — Autonomous self-looping flags executing continuous 6-phase AI self-looping execution cycles.
- `gitmap llm train --heal-apply <stash|wip|discard>` — Non-interactively applies git-state remediation recipes during Phase 6a.
- `gitmap llm train --skill-path <path>` — Target filesystem path for writing the generated skill markdown file (default: `.agents/skills/gitmap/SKILL.md`).
- `gitmap llm train --text-only` — Outputs curriculum text to stdout without writing skill file to disk.
- `gitmap llm train --urls` — Outputs the 7 authoritative public documentation registry links for LLM ingestion.
- `gitmap llm train --url` — Outputs raw public URL to `llm.md` specification.
- `gitmap llm train -j` / `gitmap llm train --json` — Emits structured machine-readable JSON curriculum and document metadata.
- `gitmap llm-docs` (alias: `gitmap ld`) — Consolidated markdown command matrix reference for LLMs.
- `gitmap llm` — Displays full LLM specification and operational guidelines.

#### Autonomous Training Flag Matrix
| Flag | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `--days <N>` | `int` | `7` | Sets memory digest inspection window in days for agent targets (e.g. `gitmap llm train muse --days 7`). |
| `--loop` | `bool` | `false` | Triggers continuous autonomous 6-phase AI self-looping execution cycles. |
| `--self-loop <N>` | `int` | `0` | Specifies exact number of consecutive iterations for AI self-loop. |
| `--heal-apply <recipe>` | `string` | `""` | Non-interactively applies git-state remediation recipes during Phase 6a (`stash`, `wip`, `discard`). |
| `--skill-path <path>` | `string` | `.agents/skills/gitmap/SKILL.md` | Target filesystem path for writing generated skill markdown file. |
| `--text-only` | `bool` | `false` | Emits complete curriculum text to stdout without writing files to disk. |
| `--urls` | `bool` | `false` | Outputs complete 7-document Authoritative Public Documentation Registry for LLM memory ingestion. |
| `-j`, `--json` | `bool` | `false` | Emits structured machine-readable JSON curriculum and document metadata. |

### 9. Multi-Repo, Cluster & Toolchain Operations
- `gitmap pae --json` — Multi-repo pull with compact JSON telemetry (use only when explicitly requested; ban routine polling).
- `gitmap cluster --help` — Orchestrate multi-node clusters and health checks.
- `gitmap sc --help` — Servers-clients topology and background task manager.
- `gitmap ssh --help` — SSH discovery, connection pooling, and remote command execution.
- `gitmap ssh identity --email <email> --seq <N> --operator <name>` — Formats canonical SSH comment token `<email>-muse<muse-sequence>-<operator-name>` and sets `git config` user identity.
- `gitmap cargo status` — Inspect Rust and Cargo toolchain status.
- `gitmap install cargo` — Install Rust toolchain if missing.
- `gitmap install --list` — Discover developer toolchains, profiles, and runtime packages.

### 10. Native Fleet Synchronization & SQLite Agent Task Engine
- `gitmap sync` — Synchronizes canonical prompts, skills, shared specs (`02-spec/01-20`), and additive scripts across all 43 registered fleet repositories.
- `gitmap sync --projects <path|json>` — Accepts a path to JSON file or inline JSON array of repositories (e.g. `'[{"folder": "cat-my"}]'`).
- `gitmap sync --repo <name>` — Synchronizes a single target repository by name.
- `gitmap sync --workers <N>` — Sets parallel worker pool concurrency (default: 8).
- `gitmap sync --dry-run` — Previews changes across all repositories without making git or file mutations.
- `gitmap sync --no-push` — Applies changes and commits locally without pushing to remote.
- `gitmap sync --no-release` — Disables post-sync SemVer release tagging.
- `gitmap agent task init --name "<task>" --budget <N>` — Initializes SQLite task manager in `.ai-memory/temp-agents/<slug>/agent-task.db`. If task already exists, reports parent task ID, active subagents, subtasks status, and touched files. Supports prompt MD5 hash verification to resume seamlessly without redundant reads.
- `gitmap agent task status [--json] [--ddl]` — Emits structured task status, subtask progress, and schema validation.
- `gitmap agent task ls` — Lists active and completed parent tasks across the workspace.
- `gitmap agent task pending` — Lists or counts pending subtasks across agent tasks.
- `gitmap agent task completed` — Lists completed root-level parent tasks with proof.
- `gitmap agent subtask add --code "<code>" --slug "<slug>" --title "<title>" [--owned "<json>"]` — Adds a decomposed subtask bounding box to the active parent task queue.
- `gitmap agent subtask ls` — Lists all subtasks, their assignees, statuses, and file boxes.
- `gitmap agent subtask claim --agent <agent-name>` — Claims the next pending subtask atomically.
- `gitmap agent subtask claim-files` — Pre-declares file leases for collision and contention tracking.
- `gitmap agent log --subtask <id> --agent <name> --action <action> --file <path> --details <desc>` — Logs in-flight agent action for crash forensics.
- `gitmap agent subtask complete <id> --agent <name> --evidence <evidence>` — Marks subtask completed.
- `gitmap agent subtask join --subtask <id> --agent <name> --role "Assistant: <focus>"` — Dynamically joins an in-flight subtask as a collaborative assistant for parallel reading, caller analysis, or verification without file collisions.
- `gitmap agent task tree` — Renders hierarchical task swarm tree (Lead Orchestrator -> Branch Coordinators -> Leaf Workers & Dynamic Assistants).
- `gitmap agent diagnose` — Performs deep autopsy and diagnostic reporting on agent task runs.
- `gitmap agent crashed` — Detects crashed or abandoned agent tasks.
- `gitmap agent temp-clear` — Safely deletes ephemeral `.ai-memory/temp-agents/` directory.

> [!CAUTION]
> **TOTAL BAN ON RAW PYTHON SQLITE (`python -c "import sqlite3..."`):**
> NEVER run inline Python sqlite3 snippets or execute raw database queries against `.ai-memory/temp-agents/<slug>/agent-task.db`.
> All agent task lifecycle operations, status inspections, claiming, logging, and diagnostics MUST use compiled `gitmap agent` CLI commands (`task init`, `task status`, `task ls`, `subtask add`, `subtask ls`, `subtask claim`, `subtask complete`, `subtask fail`, `log`, `diagnose`, `crashed`, `temp-clear`).

### 11. Fleet Telemetry & Running-Only Delegation
- `gitmap nodes-running-only` (alias `gitmap nro`) — Lists strictly online/open nodes (`is_online: true`), filtering out offline noise.
- `gitmap nro --json` (`-j`) — Emits clean JSON array of online cluster nodes.
- `gitmap nro exec "<cmd>"` — Runs SSH execution strictly across online nodes without offline timeout delays.
- `gitmap nro tasks [--json]` — Interrogates active and pending task queues on online nodes.
- `gitmap nro resources [--json]` — Probes CPU & RAM benchmarks across online nodes.
- `gitmap nro ping` — Pings reachable nodes with zero timeout lag.
- `gitmap nro pull` / `gitmap nro cp` / `gitmap nro scan` — High-speed running-only fleet operations.
- `gitmap os --json` — Structured machine and OS identity with kernel, build, and version telemetry.
- `gitmap cpu` / `gitmap processor` — Live & 5s cached CPU utilization, architecture, and core counts.
- `gitmap ram` — Live & 5s cached memory capacity, free/used MB, and visual utilization meter.
- `gitmap disk` — Disk volume capacity, free space, and mount telemetry.
- `gitmap network [--json]` — Network adapter configuration, IPs, MAC addresses, link speeds, and throughput metrics.
- `gitmap nodes tasks [--json] [--open-only]` — Queries local and remote node task queues with `?` role placeholder.

### 12. File Contention Heatmap & Agent Activity Observation
- `gitmap heatmap log -t <task-slug> -a <action> -f <file> -r "<reasoning>" [-s <state>]` — Logs granular file touches, reads, modifications, and reasoning directly into the centralized heatmap database (`.ai-memory/heatmap.db`).
- `gitmap heatmap search [query] [-f <file>] [-t <task>]` — Searches logged actions, code snippets, and reasoning notes across tasks and repositories.
- `gitmap heatmap tree [-t <task-slug>]` — Renders hierarchical tree view of touched files, line ranges, and action counts.
- `gitmap heatmap status [-t <task-slug>]` — Summary of actions logged, unique files touched, expectation metrics, and verification states.
- `gitmap heatmap discuss [--window <duration>] [--threshold <n>] [--json]` — Multi-agent file contention analysis and collision blast radius report before touching shared files.
- `gitmap heatmap live` — Real-time continuous stream of agent actions, file touches, and reasoning across concurrent agents.
- `gitmap agent log --subtask <id> --agent "<name>" --action "<read|modify|plan>" --file "<path>" --details "<reasoning>"` — In-flight subtask action logging with explicit rationale.
- `gitmap agent heatmap [--parent <slug>] [--json]` — Ranks files by write-claim count across agent tasks and active subtasks.
- `gitmap full-summary` (alias `gitmap fs`) — 48-hour active repository activity heatmap TreeView showing commit velocity and dirty working states.
- **MANDATORY HEATMAP & CONCURRENCY RULE:** In multi-agent autonomous sessions, agents MUST record all file read/edit decisions with reasoning via `gitmap heatmap log` or `gitmap agent log`, and inspect active file leases and contention hotspots via `gitmap agent heatmap` (or `gitmap heatmap status`) before editing shared files.
- **Roadmap-Connected Heatmap for Dynamic Join:** Newly spawned or dynamically joining subagents (`gitmap agent subtask join`) inspect `gitmap rc heatmap status` and `gitmap rc heatmap tree` to discover where prior agents read/edited and what subtasks remain in `.ai-memory/plans/<slug>.md`, eliminating redundant reads.

### 13. Multi-AI Companion & repo-cache Synchronization (RC Heatmap Suite)
- **Automatic RC Vault Creation:** If `repo-cache` (`rc`) or `.ai-memory/heatmap.db` does not exist, GitMap auto-provisions the required split-db structure linked to the active mission and folder.
- `gitmap heatmap rc list-ai` — Lists all supported AI namespaces (`antigravity`, `muse`, `cloud-code`, `codex`, `grok`, `cursor`) and displays ambient detected engine.
- `gitmap heatmap rc log [flags]` — Logs cognitive actions and decisions into dedicated `repo-cache` isolated split-databases (`.heatmap/<ai-name>/<repo>/heatmap.db`).
- `gitmap heatmap rc commit [-m "<msg>"]` — Creates a telemetry checkpoint commit in `repo-cache` and marks local changes clean.
- `gitmap heatmap rc sync` — Bidirectionally reconciles in-flight decisions with `repo-cache` without workspace git conflicts.
- `gitmap heatmap rc tree` / `gitmap rc heatmap tree` — Renders hierarchical roadmap and touched file tree in RC vault.
- `gitmap heatmap rc live` / `gitmap rc heatmap live` — Real-time activity stream in RC vault.

### 14. AI Decision Heatmap & RC Admin Telemetry
- `gitmap heatmap tables [--db <root|repo|local>] [--format <table|json|yaml>]` — Inspects the SQLite master catalog and inventories all tables, views, column counts, and row counts in the active heatmap database.
- `gitmap heatmap rc tables [--ai <name>] [--repo <slug>] [--db <root|repo|local>] [--format <table|json|yaml>]` — Inventories tables across companion `repo-cache` databases (`root-heatmap.db` or `<repo-slug>/heatmap.db`).
- `gitmap heatmap view <table> [--limit <n>] [--format json|yaml|table] [--task <slug>] [--repo <slug>] [--id <id>]` — Queries records from a specific heatmap table with column-aware filtering and bounded limits.
- `gitmap heatmap rc view <table> [--limit <n>] [--format json|yaml|table] [--task <slug>] [--repo <slug>] [--id <id>] [--ai <name>]` — Views table records from the target `repo-cache` companion database.
- `gitmap heatmap query "<sql>" [--db <path|root|repo|local>] [--format table|json|yaml] [--limit <n>]` — Executes arbitrary read-only SQL queries against the local or split heatmap database with automatic `mode=ro` protection, `PRAGMA query_only = ON;`, and 5000ms busy timeout.
- `gitmap heatmap rc query "<sql>" [--ai <name>] [--repo <slug>] [--tier <repo|root>] [--format table|json|yaml] [--limit <n>]` — Direct read-only SQL execution against companion `repo-cache` split or root databases.

#### Tri-Alias Routing Parity
All companion RC commands support full tri-alias routing parity:
- `gitmap heatmap rc <status|tables|view|query|tree|live|sync>`
- `gitmap rc heatmap <status|tables|view|query|tree|live|sync>`
- `gitmap heatmap rc-heatmap <status|tables|view|query|tree|live|sync>`

#### AI Agent Non-Blocking SQL Recipes
AI agents analyzing prior session decisions, execution trees, and file hotspots should use direct read-only queries with `--json` or `--yaml` to prevent lock contention:

```bash
# 1. Inspect recent chronological trace for active task slug:
gitmap heatmap query "SELECT HeatmapActionLogId, ActionType, FilePath, ReasoningState, IsExpectationMet FROM HeatmapActionLog WHERE TaskSlug = 'heatmap-commands-testing-and-rc-admin-instructions' ORDER BY HeatmapActionLogId ASC" --json

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
9. **Total Ban on Python SQLite Task Scripts:** NEVER run `python 03-ai-scripts/46-agent-sqlite-task-manager.py` or `python -c "import sqlite3..."`. Use native compiled `gitmap agent` subcommands exclusively.
