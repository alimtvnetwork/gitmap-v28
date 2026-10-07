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

## Essential Command Cheat Sheet

### 1. High-Performance Automation (AUM)
- `gitmap aum search <pattern> [dir] [--ext <ext>] [-r] [-i]` (alias: `gitmap aum grep`) — Multi-core streaming live search with lazy regex and binary filtering. ALWAYS scope with target `[dir]` and `--ext`. Replaces slow PowerShell `Select-String`, `Get-ChildItem -Recurse`, and `git grep`. TOTAL BAN on PowerShell `Select-String`, `rg`, `ripgrep`, and `git grep`.
- `gitmap search <query> [--limit <n>]` — Instant SQLite cached symbol & keyword search across scanned repositories using DH2D split-db hot cache.
- `gitmap aum guard` — Enforces 500 KB limit, large JSON exclusion, and binary null-byte probe.
- `gitmap aum sequence` — Markdown sequence gap detector and `# XX Title` autofixer.
- `gitmap aum exclude list` — Query persistent search exclusions from SQLite.
- `gitmap aum newlines --fix` — Polyglot CRLF to LF and trailing whitespace normalizer.
- `gitmap aum cache status` — Sub-millisecond in-memory cache status.
- `gitmap aum locate [tool]` — Ultra-fast tool finder (<15ms, e.g. `vcvarsall.bat`, `msbuild`, `python`; replaces slow PowerShell traversal).
- `gitmap aum benchmark all` — Side-by-side Go vs Python execution benchmarks.

### 2. Autonomous Agent Onboarding & Curriculum (LLM)
- `gitmap llm train` (alias: `gitmap llm chain`) — Full 4-stage chained curriculum, auto-generates Antigravity skill, author/sponsor attribution.
- `gitmap llm train --text-only` — Output curriculum to stdout without modifying files on disk.
- `gitmap llm-docs` (alias: `gitmap ld`) — Consolidated markdown command matrix reference for LLMs.
- `gitmap llm` — Display full LLM specification and operational guidelines.

### 3. Autonomous CI/CD Self-Healing (Pipeline AI)
- `gitmap pipeline-ai status --json` — Check workflow execution state, active branch, and ETA.
- `gitmap pipeline-ai status -t <eta>` — Wait dynamically for pipeline completion without tight polling.
- `gitmap pipeline error-logs` (alias: `gitmap pe`) — Extract failing step logs to file for 4-part RCA.
- `gitmap pe history-ai` — Analyze CI/CD pipeline history across branches and recent runs.
- `gitmap pipeline purge` — Actions zero-storage purge maintaining 0.0 GB footprint.

### 4. Fast File Discovery & Inventory
- `gitmap find "<pattern>" [-ext <ext>]` — Find files matching glob pattern in <10ms across 10,000+ files.
- `gitmap find-files <name>` (alias: `gitmap ff <name>`) — Find exact filename with optional `-ext`.
- `gitmap find-files-any <str>` (alias: `gitmap ffa <str>`) — Find files matching substring.
- `gitmap find-files-startswith <prefix>` (alias: `gitmap ffs <prefix>`) — Find by filename prefix.
- `gitmap find-files-endswith <suffix>` (alias: `gitmap ffe <suffix>`) — Find by filename suffix (e.g. `_test.go`).
- `gitmap list-files [dir]` (alias: `gitmap lf [dir]`) — List relative file paths matching pattern or directory.
- `gitmap replace <old> <new>` — Exact literal string replacement with audit trail.
- `gitmap replace-regex <pat> <subst>` — Regex replacement across repository.

### 5. Deterministic Terminal File Streaming
- `gitmap cat <filepath>` — Direct zero-disk stream of file content into process stdout. Inspect file content immediately in resource-constrained CLI sessions without buffer bloat.

### 6. Script Execution & Cache Offloading
- `gitmap py "<code-or-script>"` — Cross-platform on-the-fly Python execution.
- `gitmap pwsh "<cmd>"` / `gitmap ps "<cmd>"` — Cross-platform PowerShell execution with `-NoProfile` and automatic fallback.
- `gitmap bash "<cmd>"` / `gitmap sh "<cmd>"` — Cross-platform Bash execution.
- `gitmap rc file <filepath> [--repo <repo-name>]` — Copies script into `repo-cache` (`repo-storage`), commits and pushes automatically.
- `gitmap rc folder <folderpath> [--repo <repo-name>]` — Copies folder into `repo-cache`, commits and pushes automatically.
- `gitmap rc text "<content>" --slug <slug> --ext <.ps1|.py>` — Writes temporary script directly into `repo-cache` for permanent cross-repo reuse.

### 7. Toolchain Locator & Cache Backup
- `gitmap aum locate python` — Ultra-fast toolchain discovery (<15ms) caching exact Python executable path into SQLite `installation.db`.
- `gitmap aum cache status` — Sub-millisecond status check verifying cached toolchain paths.

### 8. Multi-Repo & Cluster Operations
- `gitmap pae --json` — Multi-repo pull with compact JSON telemetry (use only when explicitly requested; ban routine polling).
- `gitmap cluster --help` — Orchestrate multi-node clusters and health checks.
- `gitmap sc --help` — Servers-clients topology and background task manager.
- `gitmap ssh --help` — SSH discovery, connection pooling, and remote command execution.
- `gitmap cargo status` — Inspect Rust and Cargo toolchain status.
- `gitmap install cargo` — Install Rust toolchain if missing.
- `gitmap install --list` — Discover developer toolchains, profiles, and runtime packages.

### 9. Semantic Hyphen-Separated Commit & Push
- `gitmap cpf "<module> - <summary>"` — Stage, commit, and push feature branch (GitMap auto-prefixes `Feature: `).
- `gitmap cpb "<module> - <summary>"` — Stage, commit, and push bugfix branch (GitMap auto-prefixes `Bug: `).
- `gitmap cpr "<module> - <summary>"` — Stage, commit, and push release chore.
- `gitmap pcp "<module> - <summary>"` — Pull latest, commit, and push with preflight verification.
- **TOTAL BAN ON COLONS IN COMMIT MESSAGES:** Never use colons inside commit arguments (e.g. `gitmap cpf "Feature: title"` is FORBIDDEN; use `gitmap cpf "module - title"`).

## Operational Guardrails
1. **Learning & Skill Acquisition:** Run `gitmap llm train` to initialize or update GitMap skills. Never run broad keyword searches like `gitmap aum search "train"` to discover how commands work.
2. **Mandatory Pre-Flight Pull:** Always run `git pull` before modifying code.
3. **Scoped Search:** Always provide target directories and extensions to `gitmap aum search` (e.g. `gitmap aum search "target" cli --ext .go`).
4. **File Size & Binary Guard:** Respect 500 KB limit (Rule R19); never commit test binaries or temp artifacts.
5. **Coding Guidelines:** Max 8–15 lines per function, single return types with `*appfault.AppError`, affirmative booleans (`isReady`, `hasCache`).
6. **Script Offloading:** Save all temporary diagnostics to `repo-cache` via `gitmap rc` to prevent dirty working trees.
7. **Strict Relative Git Paths:** All paths and references must be relative to repository root; zero absolute paths and zero `file:///` URIs.
