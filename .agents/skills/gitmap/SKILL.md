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

### 9. Multi-Repo, Cluster & Toolchain Operations
- `gitmap pae --json` — Multi-repo pull with compact JSON telemetry (use only when explicitly requested; ban routine polling).
- `gitmap cluster --help` — Orchestrate multi-node clusters and health checks.
- `gitmap sc --help` — Servers-clients topology and background task manager.
- `gitmap ssh --help` — SSH discovery, connection pooling, and remote command execution.
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
