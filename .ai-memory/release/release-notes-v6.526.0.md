# GitMap v6.526.0

## What's Changed in v6.526.0

**Pending Commits Table Improvement, Dual-DB Status Cache & Backup Engine, and New Commands Discovery**

### 1. Pending Commits Table Improvement (`gitmap pc` / `gitmap pending-commits`)
- **Consolidated `UNCOMMITTED` Column**: Combined disparate status indicators (untracked, modified, and staged) into a single, unified `UNCOMMITTED` count column, saving terminal horizontal space.
- **Short Version & Branch (`VER/BRANCH`)**: Enriched each repository row with its release tag version and active git branch (e.g. `v6.525.0/main`).
- **Hierarchical Tree-View Remediation**: For every dirty repository, renders a structured tree view immediately beneath the row with concrete, runnable solution commands (`├── Option 1: git -C "<path>" add -A && git commit -m "wip: save changes" && git push`, `└── Option 2: git -C "<path>" stash -u`).
- **Overall Fleet Batch Fix Command**: Prominently featured in the table summary footer (`Fleet Remediation: gitmap cpar "wip: save changes"`), providing a one-line command to batch-fix all repositories across the workspace.

### 2. Dual-Database Status Caching & Persistent Backup Architecture
- **Ephemeral Status Cache DB (`pending_commits_cache.db`)**: Enforces a strict 90-second Time-To-Live (1.5 minutes, within the required 1–2 minute window). Stale entries past the TTL are immediately purged and never trusted. Bypass flags: `--no-cache` and `--refresh`.
- **Persistent Status Backup DB (`pending_commits_backup.db`)**: Dedicated secondary SQLite database that preserves status backups and historical snapshots. Survives ephemeral cache purges and can be served again on demand via `gitmap pc --backup` or `gitmap pc --serve-backup`.
- **Dual-Write Synchronization**: Live git scans automatically record entries to both the ephemeral cache and persistent backup databases.

### 3. New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`)
- **Git History Inspection**: Catalogs the last 100 new commands and subcommands introduced across recent GitMap milestones.
- **Filtering & Search**: Supports `--limit N` (default 100), `--category <cat>` (10 functional categories), and `--filter <query>` (aliases `-f` and `-q`).
- **Syntax, Descriptions & Runnable Examples**: Every entry provides clear purpose descriptions and copy-pasteable CLI invocation examples.
- **Machine-Readable JSON**: Supports `--json` / `-j` flag emitting structured telemetry for automated agent ingestion.

### 4. CLI Help Coverage & LLM Skills Integration
- **CLI Helpdocs**: Comprehensive documentation in `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md`.
- **CLI Constants**: Centralized constants in `cli/constants/constants_cli.go`.
- **LLM Skills Sync**: Integrated into `.agents/skills/gitmap/SKILL.md` allowing autonomous AI coding agents to discover repository changes, apply tree remediation commands, and find recent CLI capabilities.

### Specifications & Subtasks
- Architecture Specification: [01-architecture-spec.md](02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md)
- Component Specification: [02-component-spec.md](02-spec/21-app/table-improvement-and-command-enhancement/02-component-spec.md)
- Master Plan: [.ai-memory/plans/completed/81-table-improvement-and-command-enhancement.md](.ai-memory/plans/completed/81-table-improvement-and-command-enhancement.md)

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.526.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.526.0/install.sh | sh
```
