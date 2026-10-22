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
| ❌ `Get-ChildItem` (all variants: bare `Get-ChildItem`, `-Path`, `-Filter`, `-Recurse`, etc.), `dir`, `ls` (in shell), `find . -name "..."`, `dir /s /b` | ✅ `gitmap lf <dir>` or `gitmap list-files [dir]` or `gitmap find "<pattern>" [dir]` | Sub-millisecond directory listing and file glob discovery without shell traversal overhead. TOTAL BAN on PowerShell Get-ChildItem. |
| ❌ Inspecting plans manually via `Get-ChildItem -Path ".ai-memory/plans"` or `dir` | ✅ `gitmap plans [--json\|--yaml] [N]` (alias: `gitmap plan`) | Instant Tokyo Night styled inventory of completed and pending plans with sequence IDs, status pills, and backlog overflow navigation. |
| ❌ Directory File Inventory: `Get-ChildItem -Path <dir>`, `dir <dir>`, `ls <dir>` | ✅ `gitmap lf <dir>` (e.g. `gitmap lf src-tauri/src/modules/repo_db`) or `gitmap find "*.*" <dir>` | Instant directory file listing and inventory without PowerShell object pipeline or shell traversal latency. TOTAL BAN on PowerShell Get-ChildItem. |
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
| ❌ Python or ad-hoc SQLite scripts (`python -c "import sqlite3..."`), `sqlite3` CLI, DB browser tools | ✅ `gitmap db-read`, `gitmap db-add`, `gitmap db-edit`, `gitmap db-find`, `gitmap db-listing`, `gitmap db-view`, `gitmap db-rest`, `gitmap db-record-del`, `gitmap db-rm` | Compiled native Go SQLite exploration, mutation, view join, REST API server with Swagger UI & visualizer dashboard, and database lifecycle management. |
| ❌ Ad-hoc Python network scrapers or image link extractors (e.g. `python -c "import urllib.request, re..."`, `fetch_urls.py`, curl pipes to scrape `prnt.sc`/lightshot/screenshot URLs) | ✅ `gitmap prnt <url1> [url2...] [--extract-only] [--json]` (aliases: `gitmap download-prnt`, `gitmap asset download-prnt`) | Native compiled parallel fetching, multi-tier CDN regex resolution, User-Agent emulation, zero external dependencies, and automatic bidirectional Antigravity Brain synchronization. |
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

#### Directory Inventory & Path-Scoped Listing
- `gitmap lf <dir>` — Instant directory file listing (e.g. `gitmap lf src-tauri/src/modules/repo_db`).
- `gitmap find "*.*" <dir>` — Glob-filtered directory listing without disk traversal overhead.
- `gitmap aum search "<pattern>" <dir> -e <.ext>` — Scoped live code/regex search directly within target directory and extension.
- **Concrete Developer Workflows & Real-World Examples:**
  * **Module Directory Inventory:** `gitmap lf src-tauri/src/modules/repo_db` — Lists files within specific module path in <1ms without invoking PowerShell `Get-ChildItem`.
  * **Spec Folder Inspection:** `gitmap find "*.*" 02-spec/21-app/gitmap-commands-agent-spawning-heatmap-sync` — Discovers all spec files within a feature directory.
  * **Path-Scoped Code Search:** `gitmap aum search "HeatmapActionLog" src-tauri/src -e .rs` — Scoped search restricted to Rust source tree.
  * **Subtask Plan Inventory:** `gitmap lf .ai-memory/plans/subtasks` — Rapidly enumerates in-flight and pending subtask manifests.

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
- `gitmap llm train --skill-path <path>` — Target filesystem path for writing the generated skill markdown file (default: `.agents/skills/gitmap/skill.md`).
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
| `--skill-path <path>` | `string` | `.agents/skills/gitmap/skill.md` | Target filesystem path for writing generated skill markdown file. |
| `--text-only` | `bool` | `false` | Emits complete curriculum text to stdout without writing files to disk. |
| `--urls` | `bool` | `false` | Outputs complete 7-document Authoritative Public Documentation Registry for LLM memory ingestion. |
| `-j`, `--json` | `bool` | `false` | Emits structured machine-readable JSON curriculum and document metadata. |

### 9. Multi-Repo, Cluster & Toolchain Operations
- `gitmap pae --json` — Multi-repo pull with compact JSON telemetry (use only when explicitly requested; ban routine polling).
- `gitmap cluster status` — Inspect multi-node cluster status, heartbeat connectivity, and worker nodes.
- `gitmap cluster exec workers "<cmd>"` — Execute command in parallel across online worker cluster nodes.
- `gitmap sc list` — Query servers-clients topology, active services, and background daemon tasks.
- `gitmap sc ping` — Verify server-client connectivity with zero timeout lag.
- `gitmap ssh health --all` — Probe port 22 liveness across all known remote hosts and SSH nodes.
- `gitmap ssh identity --email <email> --seq <N> --operator <name>` — Formats canonical SSH comment token `<email>-muse<muse-sequence>-<operator-name>` and sets `git config` user identity.
- `gitmap cargo status` — Inspect Rust and Cargo toolchain status.
- `gitmap install cargo --check` — Validate Cargo and Rust toolchain installation.
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

### 12. File Contention Heatmap, Cognitive Telemetry & Agent Activity Observation
- `gitmap heatmap log --task <task> --action <action> --file <file> --state <state> --expected "<goal>"` — Logs granular file touches, reads, modifications, cognitive states, and reasoning directly into `.ai-memory/heatmap.db`.
- `gitmap heatmap update --id <id> --outcome "<outcome>" --met <true|false>` — Updates logged action with final outcome and verified expectation status.
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

#### Mandatory 4-Gate Cognitive Heatmap Workflow
In multi-agent and autonomous workflows, every file modification or major task phase MUST strictly follow the 4-Gate Cognitive Heatmap Workflow:
- **Gate 0: Pre-Flight Discovery & Schema Verification:**
  `gitmap heatmap status` and `gitmap heatmap tables`
  Inspect active task telemetry, verify SQLite schema tables, and detect any active file leases or contention hotspots before touching files.
- **Gate 1: Intent, State & Lease Declaration:**
  `gitmap heatmap log --task <task> --action <action> --file <path> --state <state> --expected "<goal>"`
  Declare intent, file path, cognitive reasoning state (`EXPLORATION`, `PLANNING`, `EXECUTION`, `VERIFICATION`), and expected goal before making modifications. Records action ID for lifecycle tracking.
- **Gate 2: Outcome Recording & Telemetry Close:**
  `gitmap heatmap update --id <id> --outcome "<outcome>" --met true`
  Record actual outcome, verify if expectation was met (`--met true` or `--met false`), and close the telemetry record immediately upon completing execution.
- **Gate 3: Verification & Hierarchy Audit:**
  `gitmap heatmap tree --task <task>` and `gitmap heatmap status --task <task>`
  Inspect hierarchical touch tree, file churn counts, and expectation compliance rate across the task.

### 13. Multi-AI Companion & repo-cache Synchronization (RC Heatmap Suite)
- **Automatic RC Vault Creation:** If `repo-cache` (`rc`) or `.ai-memory/heatmap.db` does not exist, GitMap auto-provisions the required split-db structure linked to the active mission and folder.
- `gitmap rc heatmap <status|tables|query|view|tree|live|sync>` — Complete companion RC telemetry suite routed directly through the unified `rc heatmap` CLI.
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
- `gitmap rc heatmap <status|tables|query|view|tree|live|sync>`
- `gitmap heatmap rc <status|tables|query|view|tree|live|sync>`
- `gitmap heatmap rc-heatmap <status|tables|query|view|tree|live|sync>`

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

### 15. Native SQLite Database Exploration, REST API Server & Lifecycle Suite
- `gitmap db-read "<path-or-alias>" "<sql>"` — Executes read-only SQL query with request timing.
- `gitmap db-read "<path-or-alias>" ls` — Lists all tables and views with row/column counts.
- `gitmap db-read "<path>" -as <alias>` — Registers database file into persistent alias registry.
- `gitmap db-read <alias> table <t1> [t2] [--columns] [--json]` — Displays schema, data types, and primary keys.
- `gitmap db-add <alias> <table> csv <col1>, <col2>, ...` — Inserts record from CSV values.
- `gitmap db-edit <alias> -table <table> -id <id> csv <col1>, <col2>, ...` — Updates record by ID.
- `gitmap db-record-del <alias> -table <table> -id <id> [-y] [--json]` (alias `gitmap dbrd`) — Deletes single record by primary key.
- `gitmap db-find <alias> -table <t1,t2> --find "<text>" [--columns <c1,c2>] [--json]` — Full-text and wildcard search.
- `gitmap db-listing` (alias `gitmap dbl`) — Lists all registered and system databases with file sizes, table counts, and query metrics.
- `gitmap db-view <alias> -t "t1 as a, t2 as b" --viewname <name> --find "<condition>"` — Joins tables and creates persistent views on the fly.
- `gitmap db-rest "<db-path>" as <alias>` — Registers and binds database alias for REST.
- `gitmap db-rest <alias> <port1,port2>` — Launches REST API daemon probing candidate ports, auto-generates dynamic Swagger UI (`/swagger`) and custom Interactive Database Visualizer (`/ui`).
- `gitmap db-rest <alias1,alias2> <ports>` — Mounts multiple databases into a single unified REST gateway.
- `gitmap db-rm <alias> [--purge] [-y]` (alias `gitmap dbrm`, `gitmap db-delete`) — Unregisters database alias and optionally purges underlying SQLite file.

### 16. Screenshot Asset Extraction, Lightshot/prnt.sc Downloading & Brain Sync
- `gitmap prnt <url> [dest]` — Downloads Lightshot / prnt.sc screenshot, saves to `assets/screenshots/`, and auto-syncs to active Antigravity Brain.
- `gitmap prnt <url1> <url2>...` — Concurrent batch download of multiple screenshot URLs.
- `gitmap prnt <url1> <url2>... -e` (or `--extract-only`, `-p`) — Resolves and displays direct image links without downloading.
- `gitmap prnt <url1> <url2>... --json` — Emits structured Universal JSON response envelope.
- `gitmap prnt <url1> <url2>... -o <dir>` — Saves batch downloads to designated custom directory.

---

## Practical Real-World Command Examples & Workflows

This section provides authoritative, production-tested, copy-pasteable command examples with concrete flags, arguments, and expected behavior across all 12 operational categories. AI agents must refer directly to these concrete templates rather than probing `--help` flags.

### 1. Authentication & GitHub Credential Management
Managing credentials safely without interactive prompts, plaintext leakage, or browser popups during automated turns.

```bash
# Non-interactive token ingestion (pre-validated against GitHub API before saving to global git config)
gitmap login --token "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
gitmap login token "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# Offline token ingestion without network verification (for air-gapped runners or local testing)
gitmap login --token "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" --no-verify

# Non-interactive SSH key profile configuration
gitmap login ssh work -y

# Inspect current authenticated user and masked credential origin
gitmap login --status

# List resolved token metadata, scopes, and credential origin from git configuration
gitmap token list

# Interactive OAuth browser flow (manual developer sessions only)
gitmap login --web

# Securely purge stored credentials from global git configuration
gitmap logout
```

### 2. Workspace & Multi-Repository Status
High-speed workspace auditing across 40+ repositories in a single invocation without shell traversal loops.

```bash
# Multi-repository status matrix filtered strictly to dirty working trees
gitmap status --dirty
gitmap st -d

# Check repositories with unpushed local commits waiting for remote delivery
gitmap status --ahead

# Check repositories with incoming remote commits waiting to be pulled
gitmap status --behind

# Structured machine-readable JSON status for automated agents and scripts
gitmap status --json

# Check if active repository remote tracking branch has any incoming updates
gitmap has-any-updates
gitmap hau

# Query default remote branch name (main vs master) across repositories
gitmap branch default
gitmap b def

# Discovers the most recently committed or updated remote branch
gitmap latest-branch
gitmap lb

# Terminal live-refresh status monitor without background network fetching
gitmap watch --json --no-fetch

# Single-line compact status summary across all tracked repositories
gitmap status --compact
```

### 3. Script Execution & Runner Engines
Cross-platform polyglot execution with automatic interpreter resolution from `installation.db`, `-NoProfile` safety, and zero shell quoting breakage.

```bash
# Run file manipulator script with subcommands
gitmap py 03-ai-scripts/03-file-manipulator.py lowercase

# Run Python script with verified interpreter from installation.db
gitmap py "03-ai-scripts/06-cicd-local-runner.py"

# Inline Python expression evaluation with sub-millisecond execution (<15ms)
gitmap py -c "import sys; print(sys.version_info[:2])"

# Cross-platform PowerShell execution with deterministic -NoProfile and UTF-8
gitmap pwsh "Get-Process | Where-Object WorkingSet -gt 500MB"

# Direct inline PowerShell command evaluation
gitmap ps -c "Test-Path .ai-memory/heatmap.db"

# Cross-platform POSIX Bash execution across Windows, Linux, and macOS
gitmap bash "chmod +x scripts/*.sh && ./scripts/check.sh"

# Direct inline Bash command evaluation
gitmap sh -c "uname -sm"

# Runner dry-run preview with JSON telemetry output
gitmap py "03-ai-scripts/14-version-sync-checker.py" --dry-run --json
```

### 4. Storage, Cache & Secrets Offloading
Preventing repository pollution by routing scratch scripts, diagnostic fixtures, and secrets to dedicated vaults (`repo-secrets` / `repo-cache`).

```bash
# Offload sensitive environment file into repo-secrets and auto-push
gitmap rs file .env --repo coding-guidelines

# Offload private directory into repo-secrets
gitmap rs folder private/keys --repo coding-guidelines

# Write secret string directly to repo-secrets vault
gitmap rs text "ghp_secret_pat_value" --slug github-token

# Enumerate stored secret files and metadata in JSON format
gitmap rs ls --json

# Offload temporary PowerShell diagnostic script to repo-cache (rc)
gitmap rc file debug-runner.ps1 --repo coding-guidelines

# Offload fixture directory to repo-cache
gitmap rc folder test/fixtures --repo coding-guidelines

# Write reusable temporary test harness directly to repo-cache
gitmap rc text "$env:DEBUG=1; go test -v ./..." --slug test-debug --ext .ps1

# Fast navigation shortcuts to vault locations
gitmap cd rs
gitmap cd rc
```

### 5. High-Performance Automation (AUM) & File Discovery
Replacing slow PowerShell `Get-ChildItem`, `Select-String`, and grep with multi-core streaming, lazy regex, and SQLite hot caches.

```bash
# Multi-core streaming live code search with extension and directory scoping
gitmap aum search "RunFleetPASCommand" cli -e .go

# Case-insensitive regex search with boundary patterns
gitmap aum search -r "(\"pas\"|\"pa\")" cli/cmd -e .go -i

# Sub-millisecond directory listing (instant replacement for Get-ChildItem)
gitmap lf "02-spec/21-app"

# Fast glob file discovery across 10,000+ files in <10ms
gitmap find "*.md" 02-spec/21-app

# Find files matching substring across entire codebase
gitmap find-files-any "heatmap"
gitmap ffa "heatmap"

# Find files by filename prefix or suffix
gitmap find-files-startswith "01-architecture"
gitmap find-files-endswith "_test.go"

# Direct zero-disk file content streaming into stdout
gitmap cat "version.json"

# Markdown sequence gap detector and header autofixer
gitmap aum sequence 02-spec/21-app --fix

# Normalize CRLF to LF line endings and strip trailing whitespace
gitmap aum newlines . --fix

# Enforce 500 KB file limit, large JSON exclusion, and binary probes
gitmap aum guard . --max-kb 500 --json

# Exact literal string replacement across files with audit trail
gitmap replace "old_function_name" "new_function_name"
```

### 6. Autonomous CI/CD Self-Healing (Pipeline AI)
Non-polling CI/CD status inspection, calculated ETA timeout waits, bounded failure extraction, and zero-storage artifact cleanup.

```bash
# Dynamic non-polling ETA wait for running GitHub Actions workflow (in seconds)
gitmap pipeline-ai status -t 45
gitmap pipeline-ai status -t 120

# Check pipeline workflow execution state with JSON output
gitmap pipeline-ai status --json

# Extract concise failing step error summary for 4-part Root Cause Analysis (RCA)
gitmap pe -t

# Extract bounded failure frame lines (5 before, 20 after) from latest failed run
gitmap pe -1

# Clear accumulated pipeline error logs non-interactively
gitmap pe clear -y

# Analyze CI/CD pipeline historical trend across recent runs
gitmap pe history-ai

# Purge Actions workflow run artifacts to maintain strict 0.0 GB storage quota
gitmap pipeline purge
```

### 7. Semantic & Flat Commit Suite
Disciplined atomic commit workflow enforcing hyphen-separated messages and remote delivery invariants.

```bash
# Intermediate local checkpoint commit (auto-stages, no remote push, verbatim message)
gitmap commit "wip - draft initial types and interfaces"
gitmap cm "wip - draft initial types and interfaces"

# Feature commit (auto-prefixes 'Feature: ', stages all, and pushes to GitHub tracking branch)
gitmap cpf "skills - enrich gitmap skill with 12 real world categories"

# Bugfix commit (auto-prefixes 'Bug: ', stages all, and pushes to GitHub tracking branch)
gitmap cpb "cloner - correct concurrency flag parsing and worker limits"

# Chore/Maintenance commit (auto-prefixes 'Chore: ', stages all, and pushes)
gitmap cpc "docs - realign companion skills and remove help flags"

# Release chore commit (auto-prefixes release chore, stages, and pushes)
gitmap cpr "release - bump minor version and update changelog"

# Pull latest with rebase, stage all, commit, and push in one verified flow
gitmap pcp "skills - sync cross-platform runner templates"

# Apply automated repository remediation (stash, discard, or auto-fix)
gitmap fix all -y

# Safe 2-step git mv file case normalization for lowercase hygiene
gitmap lowercase *.md -y
gitmap lowercase
gitmap lowercase-readme
```

### 8. Multi-Repo, Cluster & Toolchain Operations
Managing remote cluster topologies, SSH node pooling, and developer toolchains.

```bash
# Inspect cluster node status, heartbeat connectivity, and worker nodes
gitmap cluster status

# Execute command across all online cluster workers in parallel
gitmap cluster exec workers "df -h / && free -m"

# Servers-clients topology service list and status inspection
gitmap sc list

# Run background daemon task via server-client engine
gitmap sc bash "python -m http.server 8080"

# Verify server-client daemon connectivity and response latency
gitmap sc ping

# Check remote SSH fleet health and port 22 connectivity across all nodes
gitmap ssh health --all

# Register remote node with authentication vault
gitmap ssh join user@192.168.1.50 --alias worker-node-1

# Inspect Rust and Cargo toolchain status
gitmap cargo status

# Validate Cargo toolchain installation without modifying environment
gitmap install cargo --check

# Install missing toolchain or system runtime
gitmap install cargo
```

### 9. Native Fleet Synchronization & SQLite Agent Task Engine
Compiled native Go task tracking (<1ms WAL mode) and 43-repo synchronization.

```bash
# Synchronize canonical prompts, skills, and scripts across all 43 repositories
gitmap sync --workers 8

# Synchronize a single targeted repository by name
gitmap sync --repo aukgit/core

# Preview fleet synchronization without applying changes
gitmap sync --dry-run

# Initialize agent task tracking in .ai-memory/temp-agents/<slug>/agent-task.db
gitmap agent task init --name "skills-enrichment" --budget 400

# Interrogate task status and active subtask progress in JSON format
gitmap agent task status --json

# Add decomposed subtask bounding box with exclusive file ownership
gitmap agent subtask add --code "ST-01" --slug "enrich-gitmap-skill" --title "Enrich gitmap skill with examples"

# List pending subtasks in current task ledger
gitmap agent task pending

# Claim next available subtask atomically
gitmap agent subtask claim --agent "worker-01"

# Mark claimed subtask completed with verification evidence
gitmap agent subtask complete 1 --agent "worker-01" --evidence "12 categories added and verified"

# Mark subtask failed with explicit diagnostic rationale
gitmap agent subtask fail 1 --agent "worker-01" --reason "syntax error in target spec"

# Perform post-mortem autopsy and crash diagnosis
gitmap agent diagnose
gitmap agent crashed

# Clean ephemeral agent task database after verified completion
gitmap agent temp-clear
```

### 10. File Contention Heatmap, Cognitive Telemetry & RC Sync
Sub-millisecond SQLite WAL tracking of agent decisions, active leases, and contention hotspots.

```bash
# Step 0 Pre-Flight schema and status verification
gitmap heatmap status
gitmap heatmap tables

# Pre-Action Intent Logging (mandatory prior to reading or modifying files)
gitmap heatmap log --task "skills-enrichment" --action "READ" --file ".agents/skills/gitmap/skill.md" --state "EXPLORATION" --expected "Analyze existing command matrix"

# Post-Action Outcome Closure (mandatory immediately after operation)
gitmap heatmap update --id 1 --outcome "Identified missing real-world examples" --met true

# Render hierarchical causal tree of all actions taken during the task
gitmap heatmap tree --task "skills-enrichment"

# Multi-agent file contention analysis and collision blast radius check
gitmap agent heatmap

# Safe read-only SQL query against cognitive heatmap database
gitmap heatmap query "SELECT HeatmapActionLogId, ActionType, FilePath, IsExpectationMet FROM HeatmapActionLog ORDER BY HeatmapActionLogId DESC LIMIT 5" --json

# Companion repo-cache (rc) heatmap status and synchronization
gitmap rc heatmap status
gitmap rc heatmap sync
```

### 11. Native SQLite Database Exploration, REST API Server & Lifecycle Suite
High-speed native Go SQLite exploration, mutation, join views, and dynamic REST gateways.

```bash
# List all registered and system SQLite databases with sizes and query metrics
gitmap db-listing
gitmap dbl

# List all tables and views in database with column and row counts
gitmap db-read master ls
gitmap db-read ".ai-memory/heatmap.db" ls

# Inspect table schema, data types, and primary key constraints
gitmap db-read tasks table Subtask --columns
gitmap db-read ".ai-memory/heatmap.db" table HeatmapActionLog --columns --json

# Execute read-only SQL query with timing telemetry
gitmap db-read ".ai-memory/heatmap.db" "SELECT COUNT(*) FROM HeatmapActionLog"

# Insert new record into table from CSV values
gitmap db-add tasks Subtask csv "SUB-01","my-task","PENDING"

# Update existing record by primary key
gitmap db-edit tasks -table Subtask -id 42 csv "SUB-01","my-task","COMPLETED"

# Search database across multiple tables using full-text wildcard
gitmap db-find master -table Project,Alias --find "enrichment" --json

# Delete a single record by primary key
gitmap db-record-del tasks -table Subtask -id 42 -y

# Create persistent joined view across tables
gitmap db-view tasks -t "Subtask as s, Task as t" --viewname ViewSubtasks --find "s.TaskId = t.TaskId"

# Launch REST API daemon with auto-generated Swagger UI (/swagger) and Visualizer (/ui)
gitmap db-rest master,tasks 8080,8081
```

### 12. Screenshot Extraction & Lightshot Downloading
Native compiled parallel screenshot asset extraction, CDN resolution, and Antigravity Brain sync.

```bash
# Download single Lightshot/prnt.sc screenshot and auto-sync to Antigravity Brain
gitmap prnt "https://prnt.sc/abcdef"

# Resolve direct image URLs without downloading (dry-run probe)
gitmap prnt "https://prnt.sc/abcdef" -e

# Emits structured Universal JSON response envelope and saves to custom directory
gitmap prnt "https://prnt.sc/abcdef" -o assets/screenshots --json

# Batch download multiple screenshot URLs in parallel
gitmap prnt "https://prnt.sc/url1" "https://prnt.sc/url2" "https://prnt.sc/url3"
```

### 13. Plan & Task Discovery (`gitmap plans`)
Automated plan discovery and inventory across `.ai-memory/plans/` (root, `pending/`, and `completed/`) with Tokyo Night Lipgloss tables and machine-readable exports.

```bash
# Default view: shows last 10 completed plans and top 10 pending plans
gitmap plans
gitmap plan

# Specify custom limit per section (e.g. 20 items per table)
gitmap plans 20
gitmap plans --limit 20
gitmap plans -n 20

# View all plans without truncation
gitmap plans --all
gitmap plans -a

# Filter exclusively to pending and in-progress plans
gitmap plans --pending
gitmap plans -p

# Filter exclusively to completed plans
gitmap plans --completed
gitmap plans -c

# Emit pure machine-readable JSON to stdout (pipeable to jq)
gitmap plans --json
gitmap plans -j 20

# Emit pure machine-readable YAML to stdout
gitmap plans --yaml
gitmap plans --yml
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
9. **Total Ban on Python SQLite Scripts:** NEVER run `python 03-ai-scripts/46-agent-sqlite-task-manager.py` or `python -c "import sqlite3..."`. For agent tasks, use native compiled `gitmap agent` subcommands exclusively. For inspecting, exploring, querying, mutating, or viewing SQLite databases (e.g. Antigravity DBs, split databases, application storage), use `gitmap db-read`, `gitmap db-add`, `gitmap db-edit`, `gitmap db-find`, `gitmap db-listing`, `gitmap db-view`, `gitmap db-rest`, `gitmap db-record-del`, `gitmap db-rm` (or aliases `dbr`, `dba`, `dbe`, `dbf`, `dbl`, `dbv`, `dbrs`, `dbrd`, `dbrm`).
10. **Total Ban on Ad-Hoc Python Web Scrapers:** NEVER run `python -c "import urllib.request..."` or write ad-hoc Python scripts (`fetch_urls.py`, etc.) to scrape web pages, extract images, or download screenshot URLs (e.g., `prnt.sc`, `image.prntscr.com`, `i.imgur.com`). You MUST use native compiled `gitmap prnt <url1> [url2...] [--extract-only] [--json]` (alias `gitmap download-prnt`), which automatically handles User-Agent headers, multi-tier CDN resolution, and bidirectional Antigravity Brain synchronization. See `@[.agents/skills/gitmap]` for full command reference.
