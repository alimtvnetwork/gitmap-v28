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
| ❌ Multi-repo `git status` loops or manual directory scanning | ✅ `gitmap pc` (`gitmap pending-commits`) | High-speed multi-repo status with 90s SQLite cache, consolidated `UNCOMMITTED` metrics, short branch display, tree-view surgical fixes, and batch footer remediation. |
| ❌ Manual git log digging or release note compilation | ✅ `gitmap summary [$repo] [N]` or `gitmap fs [N]` | Split-DB cached ≤200w release gists, heated file churn tracking, 48h active filter, and TreeView master commit remediation. |
| ❌ Blind multi-repo git rebase / copy-paste merge overwrites | ✅ `gitmap merge-ai <dest-target> <sources...>` (`gitmap ma`) | Single-commit staging, chronological inspection, `01_`/`02_` collision sequencing, `merge-ai-manifest.json`, and `instruction.md`. |
| ❌ Redundant cluster SSH queries polling identical repo clones | ✅ `gitmap nodes fs` / `gitmap nodes fspe` / `gitmap nodes pe all` | Two-phase fleet handshake with local-machine precedence deduplication eliminating redundant SSH calls. |
| ❌ Guessing newly added CLI commands from raw git commit logs | ✅ `gitmap nc` (`gitmap new-commands [--limit 100]`) | Native cataloging of recently introduced commands, syntax, flags, and usage examples directly from git history. |
| ❌ `gh auth login` interactive prompts or plaintext `.env` files | ✅ `gitmap login --web`, `gitmap login --status`, `gitmap login --token <PAT>` | Token validated via GitHub API before writing; auto-resolved for all clone, pull, and push commands. |
| ❌ Colons in commit messages (`git commit -m "feat: ..."` or `gitmap cpf "feat: ..."`) | ✅ `gitmap cpf "<module> - <summary>"` (hyphen-separated only) | GitMap automatically provides `Feature: ` or `Bug: ` prefix. Colons inside the message argument cause duplicate prefixes. |
| ❌ Saving temporary scratch or test scripts into repo git tree | ✅ `gitmap rc text "<content>" --slug <slug> --ext .ps1` | Centralized script storage in `repo-cache` (`repo-storage`) for permanent cross-repo reuse without polluting git worktrees. |
| ❌ Committing `.env` or credentials to standard repositories | ✅ `gitmap rs text "<secret>" --slug <slug>` | Strict zero-secrets policy; stores credentials exclusively in `repo-secrets` vault. |
| ❌ `gh run watch` or tight polling loops (`while true; sleep 5`) | ✅ `gitmap pipeline-ai status -t <eta>` or `gitmap pe -t` | Dynamic timeout waiting driven by calculated workflow ETA without burning CPU or Actions API quotas. |
| ❌ Slow Python fleet sync (`python 03-ai-scripts/38-sync-prompts-skills-scripts.py`) | ✅ `gitmap sync [--workers 8] [--projects <path|json>]` | Native Go multi-repo synchronization across 43 repositories in <5s with 6-stage safe ceremony (backup branch, pre-pull, 5 boundaries, atomic commit). |
| ❌ Python SQLite task manager (`python 03-ai-scripts/46-agent-sqlite-task-manager.py`) | ✅ `gitmap task <init|add|claim|complete|fail|status|schema>` | Native compiled Go SQLite task manager (<1ms) with WAL mode, single-writer locking, and 1:1 identical schema for multi-agent workflows. |

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
- `gitmap pc` (alias `gitmap pending-commits`) — Fast multi-repo pending commits table with 90s SQLite cache (`pending_commits_cache.db`), persistent backup DB (`pending_commits_backup.db`), consolidated `UNCOMMITTED` column (unique untracked + modified + staged files), `VER/BRANCH` display, tree remediation hints (`├── Option 1: ...`, `└── Option 2: ...`), and batch remediation command in footer (`gitmap cpar "wip: save changes"`).
- `gitmap pc --refresh` — Force-refreshes SQLite cache by re-evaluating live git status.
- `gitmap pc --no-cache` — Disables SQLite cache completely for strict pre-flight verification.
- `gitmap pc --backup` (or `--serve-backup`) — Serves status from persistent backup DB without querying live filesystem.
- `gitmap pc --json` — Emits structured JSON telemetry of uncommitted and unpushed repositories.
- `gitmap summary [$repo] [N]` — Executive release summary of the last $N$ releases (default $N=8$; target defaults to current repository `.`, or relative path / remote URL). Constrains summaries to $\le 200$ words gist per release with raw file lists suppressed. Computes heated file churn metrics (modification frequency + line deltas across commit boundaries, surfacing top 5 heated files with functional rationale). Backed by Split-DB SQLite hot cache (`.gitmap/summary.db`) with $<15\text{ms}$ hash-matched lookups (`repo_url`, `tag_commit_hash`, `head_commit_hash`) and incremental delta calculations for newly added commits.
  - *Example (current repo, default 8 releases):* `gitmap summary`
  - *Example (current repo, 3 releases):* `gitmap summary 3`
  - *Example (relative repo path, 5 releases):* `gitmap summary ./sub-repo 5`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
- `gitmap full summary [N]` (aliases: `gitmap full status [N]`, `gitmap fs [N]`, default $N=3$) — Workspace activity heatmap TreeView (`├──`, `└──`). Scans repositories in `gitmap.db` and filters to those with commit activity within the last 48 hours or in an uncommitted dirty state (unstaged, staged, untracked changes); dormant clean repos are excluded. Features detailed dirty worktree breakdown per repository, provides per-repo commit command suggestions, and outputs a consolidated master sanitize commit command at the bottom (`gitmap sanitize-all --message "wip: save active progress across dirty repos"`).
  - *Example (default 3 releases, 48h active filter):* `gitmap fs`
  - *Example (show last 5 releases):* `gitmap fs 5`
  - *Example (canonical full summary command):* `gitmap full summary`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
- `gitmap full summary+pe` (aliases: `gitmap full status+pe`, `gitmap fs+pe [--json]`, `gitmap fspe`) — CI/CD error stack trace fusion. Embeds latest pipeline diagnostic telemetry from `repodb/pipeline.db` directly into the activity heatmap TreeView. Repositories with passing pipelines display a concise green checkmark (`✓ CI/CD Passing`). Repositories with failing pipelines inline the failed job/step name, exit code, and the last 25 lines of failure stack traces with file paths and line numbers for instant AI root cause analysis.
  - *Example (fused tree status):* `gitmap fspe`
  - *Example (standard alias):* `gitmap fs+pe`
  - *Example (structured JSON output):* `gitmap fs+pe --json`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`

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

### 5. High-Performance Automation (AUM), Tree & File Discovery
- `gitmap tree [path] [--file/-f <out>]` — Generates directory tree and exports directly to `.txt` (ASCII tree), `.json` (structured report), or `.yaml` with automatic format deduction. Handles wildcard patterns (e.g. `gitmap tree "*.go"`) without filesystem errors.
- `gitmap tree-search -f <file> "<pattern*ends*>"` (alias `gitmap ts`) — High-speed wildcard pattern matching with `*` and `?` across exported tree files.
- `gitmap tree-search-startsWith -f <file> "<prefix>"` (alias `gitmap tss`) — Prefix matching with folder hierarchy preservation.
- `gitmap tree-search-contains -f <file> "<substring>"` (alias `gitmap tsc`) — Case-insensitive substring searching across paths and filenames.
- `gitmap tree-search-endsWith -f <file> "<suffix>"` (alias `gitmap tse`) — Suffix and extension searching (e.g. `_test.go`, `.yaml`).
- `gitmap tree-search-grep -f <file> "<regex>"` (alias `gitmap tsg`) — Regular expression search across exported tree files.
- `gitmap tree-learn -f <file>` (alias `gitmap tl`) — Ingests exported tree manifests into Split SQLite DB (`.gitmap/data/treedb/sql.db`) with NOCASE indexes for sub-millisecond retrieval.
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
- `gitmap pe all [--json]` — High-speed fleet CI/CD pipeline error check across repositories active in the last 48 hours. Enforces the **Concise Green Rule**: passing repositories output a single concise green status line (`✓ CI/CD Passing`), whereas failing repositories expand with full error diagnostics, failed step details, exit codes, and bounded 25-line stack traces with file paths and line numbers for rapid 4-part RCA.
  - *Example (fleet CI check):* `gitmap pe all`
  - *Example (structured JSON telemetry):* `gitmap pe all --json`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
- `gitmap pe all --force-all` — Bypasses the 48-hour activity window filter, scanning all indexed repositories across the entire workspace regardless of activity recency.
  - *Example:* `gitmap pe all --force-all`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
- `gitmap pipeline purge` — Actions zero-storage purge maintaining 0.0 GB footprint.

### 7. Semantic Hyphen-Separated Commit & Push
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

### 8. Autonomous Agent Onboarding & Curriculum (LLM)
- `gitmap llm train` (alias: `gitmap llm chain`) — Full 4-stage chained curriculum, auto-generates Antigravity skill, author/sponsor attribution.
- `gitmap llm train --text-only` — Output curriculum to stdout without modifying files on disk.
- `gitmap llm-docs` (alias: `gitmap ld`) — Consolidated markdown command matrix reference for LLMs.
- `gitmap llm` — Display full LLM specification and operational guidelines.
- `gitmap new-commands` (alias `gitmap nc`) — Discovers and filters the last 100 commands added across recent git history with runnable examples.
- `gitmap nc --limit <N>` (alias `-n`) — Limits number of commands returned (default: 100).
- `gitmap nc --filter "<pattern>"` / `gitmap nc -f "<pattern>"` / `gitmap nc -q "<pattern>"` — Filters new commands by keyword across name, alias, description, and copy-pasteable example.
- `gitmap nc --category "<cat>"` (alias `-c`) — Filters commands by functional category (`commits`, `diagnostics`, `scanner`, `fleet`, `ai`, `os`, `spec`, `storage`, `sync`, `tooling`).
- `gitmap nc --json` (alias `-j`) — Emits machine-readable JSON array of discovered commands.

### 9. Multi-Repo, Cluster & Toolchain Operations
- `gitmap pae --json` — Multi-repo pull with compact JSON telemetry (use only when explicitly requested; ban routine polling).
- `gitmap cluster --help` — Orchestrate multi-node clusters and health checks.
- `gitmap sc --help` — Servers-clients topology and background task manager.
- `gitmap ssh --help` — SSH discovery, connection pooling, and remote command execution.
- Distributed fleet commands employ a **Two-Phase Handshake (`gitmap scan export --lean-manifest`) with Local Machine Precedence Deduplication**:
  - The local master catalogs locally hosted repository paths and URLs first.
  - An ultra-lean SSH discovery probe collects manifests from remote cluster nodes.
  - Repositories hosted locally on the master machine are evaluated locally with zero SSH overhead. Remote SSH queries are dispatched *only* for repositories unique to remote nodes, preventing redundant network queries across identical clones.
  - Aggregates local and remote telemetry into a single unified TreeView or JSON payload.
- `gitmap nodes fs [N]` (aliases: `gitmap nodes full status [N]`, `gitmap nodes full summary [N]`, default $N=3$) — Distributed workspace activity heatmap and release summary across cluster nodes with handshake deduplication.
  - *Example (cluster full summary):* `gitmap nodes fs`
  - *Example (last 5 releases per active repo):* `gitmap nodes fs 5`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
- `gitmap nodes summary <repo> [N]` (default $N=8$) — Queries remote cluster nodes for release summary and heated churn metrics of a specific repository.
  - *Example:* `gitmap nodes summary my-service`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
- `gitmap nodes fs+pe [--json]` (alias: `gitmap nodes fspe`) — Distributed full summary with fused CI/CD pipeline failure logs and stack traces across all fleet nodes.
  - *Example (cluster fspe):* `gitmap nodes fspe`
  - *Example (JSON output):* `gitmap nodes fs+pe --json`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
- `gitmap nodes pe all [--json]` (alias: `gitmap nodes pipe-error-all`) — Distributed CI/CD pipeline error check across all fleet nodes, respecting the 48-hour active filter and concise green line rule.
  - *Example:* `gitmap nodes pe all`
  - *Example (JSON output):* `gitmap nodes pe all --json`
  - *Specification Reference:* `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`
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
- `gitmap sync --list` — Lists all 43 registered fleet repositories and paths.
- `gitmap task init --name "<task>" --budget <N>` — Initializes SQLite task manager in `.ai-memory/temp-agents/<slug>/agent-task.db`.
- `gitmap task add --db <path> --code <code-id> --title <title> [--files <paths>] [--role <role>]` — Adds a subtask.
- `gitmap task claim --db <path> --agent <agent-name>` — Claims the next pending subtask atomically.
- `gitmap task log-action --db <path> --subtask-id <id> --agent <name> --action <action> --file <path> --details <desc>` — Logs in-flight agent action for crash forensics.
- `gitmap task complete --db <path> --subtask-id <id> --evidence <evidence>` — Marks subtask completed.
- `gitmap task fail --db <path> --subtask-id <id> --reason <reason>` — Marks subtask failed with reason.
- `gitmap task status --db <path>` — Emits machine-readable JSON summary of task progress.
- `gitmap task schema [--json|--ddl]` — Emits task database schema and DDL definitions.

### 11. Multi-Repo AI Merge Orchestration & "Repo Feature" Resolver

#### A. Universal "Repo Feature" Destination Resolver
Whenever `<dest-target>` is passed to GitMap commands (`gitmap merge-ai`, `gitmap ma`, `gitmap clone-as`, `gitmap init-remote`), the **Universal "Repo Feature" Resolver** autonomously classifies and prepares the target repository without requiring manual, multi-step interventions:
- **Branch A: Remote Git URL** (`https://github.com/org/repo.git` or `git@github.com:...`):
  - Checks `gitmap.db` to determine if a local clone exists in the active workspace.
  - If present locally: Adopts the existing local directory directly.
  - If missing locally: Probes remote availability using authenticated GitMap credentials (`gitmap login`). If the remote repo exists, clones into `<workspace>/<slug>`; if the remote repo does not exist, creates the GitHub repository under the authenticated user/org via API, initializes the local folder, and links `origin`.
- **Branch B: Local Folder with Existing `.git`** (`./relative/path/to/repo`, `d:/work/existing-project`):
  - Verifies `.git` integrity and active branch. Adopts the directory directly as destination workspace, strictly preserving all existing commit history, branches, and tracked files.
- **Branch C: Local Folder without `.git`** (`./unversioned-folder`, `d:/work/legacy-app`):
  - Normalizes folder name to a kebab-case slug, executes `git init -b main` in the folder, queries active authenticated GitHub identity (`gitmap login --status`), creates a matching remote repository via GitHub API, and links `git remote add origin <url>`.
- **Branch D: Bare Repository Name / Slug** (`new-service-slug`):
  - Detects single-token input without path separators or schemes. Resolves to canonical `<workspace>/<slug>`, creates local folder, creates matching remote repository via GitMap credentials, initializes `git init -b main`, and links origin.
- **Preservation Invariant:** Target directories are NEVER wiped, cleared, or deleted. Existing files are preserved as base staging files.
- **Specification Reference:** `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/repo-feature.md`

#### B. AI Merge Orchestrator (`gitmap merge-ai` / `gitmap ma`)
- `gitmap merge-ai <dest-target> <sources...>` (alias `gitmap ma`) — Consolidates multiple repositories into a unified target (`<dest-target>` resolved via the "Repo Feature" resolver).
- **Single-Commit Staging:** Prepares the destination workspace for a single atomic commit. Multi-commit rebase stacking is intentionally disabled to eliminate silent overwrite regressions.
- **Chronological Inspection:** Inspects all source repositories in parallel, sorting them chronologically (earliest to latest) to establish baseline precedence. Base repository files are placed directly into the destination directory tree.
- **Collision Sequencing (`01_`, `02_`, `03_`):** Non-conflicting files are placed directly at their canonical paths. Conflicting files are renamed with sequential numerical prefixes (`01_<filename>` for base/earlier version, `02_<filename>` for subsequent source version, `03_<filename>` for tertiary version), ensuring zero silent file overwrites.
- **`merge-ai-manifest.json`:** Emits a comprehensive manifest at the destination root recording all source repositories, source commit hashes, active branches, and the complete collision map.
- **`instruction.md`:** Emits an authoritative checklist at the destination root instructing downstream AI agents to consolidate prefixed files, unify business logic, restore canonical filenames, and run test/build verification.
- **Supported Input Formats:**
  - Space-separated URLs or paths: `gitmap ma <dest> <src1> <src2> <src3>`
  - Comma-separated URLs or paths: `gitmap ma <dest> <src1>,<src2>,<src3>`
  - Text file (newline-delimited): `gitmap ma <dest> repo-list.txt`
  - Configuration JSON: `gitmap ma <dest> merge-config.json`
- *Example (Merge two repos into a new bare slug):*
  ```bash
  gitmap ma merged-core https://github.com/my-org/repo-a.git https://github.com/my-org/repo-b.git
  ```
- *Example (Merge multiple repos from text manifest into existing folder):*
  ```bash
  gitmap merge-ai ./target-app 02-spec/21-app/manifest-repos.txt
  ```
- *Specification Reference:** `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`

---

## Operational Guardrails & Non-Negotiable Invariants

1. **Learning & Skill Acquisition:** Run `gitmap llm train` to initialize or update GitMap skills. Never run broad keyword searches like `gitmap aum search "train"` to discover how commands work.
2. **Mandatory Pre-Flight Pull:** Always run `git pull` before modifying code.
3. **Scoped Search:** Always provide target directories and extensions to `gitmap aum search` (e.g. `gitmap aum search "target" cli --ext .go`).
4. **File Size & Binary Guard:** Respect 500 KB limit (Rule R19); never commit test binaries or temp artifacts.
5. **Coding Guidelines:** Max 8–15 lines per function, single return types with `*appfault.AppError`, affirmative booleans (`isReady`, `hasCache`).
6. **Script Offloading:** Save all temporary diagnostics to `repo-cache` via `gitmap rc` to prevent dirty working trees.
7. **Strict Relative Git Paths:** All paths and references must be relative to repository root (`02-spec/...`, `.ai-memory/...`); zero absolute paths and zero `file:///` URIs.
8. **Zero Storage Ban:** Zero uploads to `actions/upload-artifact`. Maintain 0.0 GB Actions storage quota across all repositories.
9. **Pending Status Caching, Backup & Remediation Invariant:** When checking multi-repo status, agents should invoke `gitmap pc`. Respect the 90-second SQLite status cache (`pending_commits_cache.db`) and persistent backup DB (`pending_commits_backup.db`). When validating state immediately after applying code modifications or git operations, pass `gitmap pc --refresh` or `gitmap pc --no-cache` to ensure live filesystem validation. If dirty repositories are reported, use the suggested tree remediation command (`├── Option 1: commit & push`, `└── Option 2: stash`) or batch footer command (`gitmap cpar "wip: save changes"`). To review historical status snapshots offline, pass `gitmap pc --backup`. To discover recently added commands and usage examples, use `gitmap nc` (with `-f` / `-q` and `-c`).
10. **Universal "Repo Feature" Invariant:** When resolving `<dest-target>`, agents must never delete or wipe pre-existing target directories. Follow `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/repo-feature.md` for zero-prompt authenticated repository creation and local adoption.
11. **Single-Commit AI Merge Invariant:** When executing `gitmap merge-ai` / `gitmap ma`, conflicting files must be sequenced as `01_<filename>`, `02_<filename>` and paired with `merge-ai-manifest.json` and `instruction.md`. Downstream agents must consolidate logic and restore canonical filenames before final commit.
12. **Fleet Handshake Deduplication Invariant:** For distributed operations (`gitmap nodes fs`, `gitmap nodes fspe`, `gitmap nodes pe all`), local-machine repositories take precedence and must be evaluated locally with zero SSH overhead. Remote nodes are queried strictly for repositories unique to those nodes.
