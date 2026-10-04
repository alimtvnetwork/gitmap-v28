---
name: gitmap
description: Autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation (AUM), cluster/SSH delegation, pipeline self-healing, multi-tier AI agent task orchestration, and coding guideline enforcement.
---

# GitMap Autonomous Engineering Skill

## Overview
GitMap is an ultra-fast developer companion and autonomous CLI engine designed for AI coding agents and software engineers.

- **Lead Architect & Author:** MD ALIM UL KARIM (alimtvnetwork)
- **Sponsored By:** RISEUP ASIA LLC (https://riseup-asia.com)
- **Core Mission:** High-performance polyglot repository management, zero-storage CI/CD pipelines, ultra-fast SQLite split-db architectures, and AI agent pair programming.

## Essential Command Cheat Sheet

### 1. High-Performance Automation (AUM)
- `gitmap aum search <pattern> [dir] [--ext <ext>] [-r] [-i]` (alias: `gitmap aum grep`) — Multi-core streaming live search with lazy regex and binary filtering. ALWAYS scope with target `[dir]` and `--ext`. Replaces slow PowerShell `Select-String`, `Get-ChildItem -Recurse`, and `git grep`. TOTAL BAN on PowerShell `Select-String` and `git grep`.
- `gitmap search <query> [--limit <n>]` — Instant SQLite cached symbol & keyword search across scanned repositories.
- `gitmap aum guard` — Enforces 500 KB limit, large JSON exclusion, and binary null-byte probe
- `gitmap aum sequence` — Markdown sequence gap detector and # XX Title autofixer
- `gitmap aum exclude list` — Query persistent search exclusions from SQLite
- `gitmap aum newlines --fix` — Polyglot CRLF to LF and trailing whitespace normalizer
- `gitmap aum cache status` — Sub-millisecond in-memory cache status
- `gitmap aum locate [tool]` — Ultra-fast tool finder (<15ms, e.g. vcvarsall.bat, msbuild; replaces slow PowerShell Get-ChildItem)
- `gitmap aum benchmark all` — Side-by-side Go vs Python execution benchmarks

### 2. Autonomous Agent Onboarding & Curriculum (LLM)
- `gitmap llm train` (alias: `gitmap llm chain`) — Full 4-stage chained curriculum, auto-generates Antigravity skill, author/sponsor attribution
- `gitmap llm train --text-only` — Output curriculum to stdout without modifying files on disk
- `gitmap llm-docs` (alias: `gitmap ld`) — Consolidated markdown command matrix reference for LLMs
- `gitmap llm` — Display full LLM specification and operational guidelines

### 3. Autonomous CI/CD Self-Healing (Pipeline AI)
- `gitmap pipeline-ai status --json` — Check workflow execution state, active branch, and ETA
- `gitmap pipeline-ai status -t <eta>` — Wait dynamically for pipeline completion without tight polling
- `gitmap pipeline error-logs` — Extract failing step logs to file for 4-part RCA
- `gitmap pipeline purge` — Actions zero-storage purge maintaining 0.0 GB footprint (Rule R18)

### 4. Fast File Discovery & Refactoring
- `gitmap find-files <name>` (alias: `gitmap ff <name>`) — Find exact filename with optional -ext
- `gitmap find-files-any <str>` (alias: `gitmap ffa <str>`) — Find files matching substring
- `gitmap find-files-startswith <prefix>` (alias: `gitmap ffs <prefix>`) — Find by filename prefix
- `gitmap find-files-endswith <suffix>` (alias: `gitmap ffe <suffix>`) — Find by filename suffix (e.g. _test.go)
- `gitmap replace <old> <new>` — Exact literal string replacement with audit trail
- `gitmap replace-regex <pat> <subst>` — Regex replacement across repository

### 5. Semantic Commit & Push
- `gitmap cpf "<msg>"` — Stage, commit, and push feature branch
- `gitmap cpb "<msg>"` — Stage, commit, and push bugfix branch
- `gitmap cpr "<msg>"` — Stage, commit, and push release chore
- `gitmap pcp "<msg>"` — Pull latest, commit, and push with preflight verification

### 6. Script Runners & Multi-Repo Operations
- `gitmap pwsh "<cmd>"` / `gitmap ps "<cmd>"` — Cross-platform PowerShell execution with `-NoProfile` and automatic fallback
- `gitmap bash "<cmd>"` / `gitmap sh "<cmd>"` — Cross-platform Bash execution
- `gitmap pae --json` — Multi-repo pull with compact JSON telemetry (use only when explicitly requested; ban routine polling)

### 7. Multi-Node Cluster & Remote Delegation
- `gitmap cluster --help` — Orchestrate multi-node clusters and health checks
- `gitmap sc --help` — Servers-clients topology and background task manager
- `gitmap ssh --help` — SSH discovery, connection pooling, and remote command execution
- `gitmap ports [service|--common]` — Inspect local listening ports, active PIDs, and firewall rules
- `gitmap ssh enable [--port <n>]` — Cross-platform OpenSSH Server installation, auto-start, and firewall configuration
- `gitmap ssh troubleshoot <alias|ip>` — Diagnostic connectivity probe for SSH timeouts and firewall blocks

### 8. Rust & Toolchain Package Management
- `gitmap cargo status` — Inspect Rust and Cargo toolchain status
- `gitmap install cargo` — Install Rust toolchain if missing
- `gitmap install --list` — Discover developer toolchains, profiles, and runtime packages

### 9. AI Agent Task Orchestrator & Multi-Tier Split-DB (`gitmap agent` / `gitmap ai-agents`)
- `gitmap agent task init --name "<name>" [--budget <n>]` — Initialize parent task & register in master `ai_agents.db`
- `gitmap agent task ls [--limit <n>] [--all] [--status <s>]` — List active and completed parent tasks
- `gitmap agent task status [--task-id <id>]` — Show task rollup, budget consumption, and subtask progress
- `gitmap agent subtask add --parent <id> --file <file>` (or `--json '<json>'`) — Enqueue atomic subtasks
- `gitmap agent subtask claim --agent "<role>"` — Atomically claim next available pending subtask
- `gitmap agent subtask start <subtask-id> --agent "<role>"` — Mark claimed subtask as in-progress
- `gitmap agent log --agent "<role>" --subtask <id> --action "<action>" --file "<file>" --details "<details>"` — Record granular telemetry into agent split DB
- `gitmap agent subtask complete <subtask-id> --agent "<role>" --evidence "<ev>"` — Mark subtask done with evidence
- `gitmap agent subtask fail <subtask-id> --agent "<role>" --reason "<reason>"` — Mark subtask failed with RCA
- `gitmap agent crashed [--task-id <id>]` (alias: `gitmap agent diagnose`) — Detect abandoned/crashed workers and display autopsy report
- `gitmap agent ui [--port <p>] [--browse]` — Launch interactive browser dashboard with tree view and live action timeline
- `gitmap agent clear [--task-id <id>] [-y]` — Clean completed tasks or specific task runs
- `gitmap agent reset [-y]` — Reset agent database tables while preserving directory structures
- `gitmap agent temp-clear [-y]` — Purge entire agent temp directory (`.ai-memory/temp-agents/`)

## Operational Guardrails
1. **Learning & Skill Acquisition:** Run `gitmap llm train` to initialize or update GitMap skills. Never run broad keyword searches like `gitmap aum search "train"` to discover how commands work.
2. **Mandatory Pre-Flight Pull:** Always run `git pull` before modifying code.
3. **Scoped Search:** Always provide target directories and extensions to `gitmap aum search` (e.g. `gitmap aum search "target" cli --ext .go`).
4. **File Size & Binary Guard:** Respect 500 KB limit (Rule R19); never commit test binaries or temp artifacts.
5. **Coding Guidelines:** Max 8–15 lines per function, single return types with `*appfault.AppError`, affirmative booleans.
