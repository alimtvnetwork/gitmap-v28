# gitmap pending-commits

Scans local workspace repositories or fleet cluster nodes over SSH to identify uncommitted working tree modifications and unpushed commits ahead of origin.

## Alias

pc

## Usage

```bash
gitmap pending-commits [repoName|all] [flags]
gitmap pc [repoName|all] [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--sort <mode>` | `-s` | `priority` | Sort order: `priority`, `name`, or `count` |
| `--detail <mode>` | `-d` | `summary` | Detail view: `summary`, `all`, `1`, or `none` |
| `--no-cache` | | `false` | Bypass SQLite status cache and query git status directly |
| `--refresh` | | `false` | Force refresh SQLite status cache with fresh repository metrics |
| `--backup` | `--serve-backup` | `false` | Serve status from secondary persistent backup DB without live git scan |
| `--ssh` | `-ssh` | `false` | Aggregate pending commits across all registered cluster SSH nodes |
| `--json` | `-j` | `false` | Output machine-readable JSON telemetry |
| `--dirty-only` | | `true` | Display only repositories with changes |
| `--all` | | `false` | Include completely clean repositories in output |
| `-h`, `--help` | | `false` | Display command help menu |

## Dual-Database Architecture

GitMap implements a high-performance dual-database architecture for pending commit telemetry:

1. **Transient Status Cache (`pending_commits_cache.db`)**:
   - Ultra-fast query response (<5ms) backed by SQLite WAL mode.
   - Enforces a strict **90-second TTL** (1.5 minutes). Stale records are purged immediately and never trusted.
   - Dual-gate HEAD SHA verification guarantees cache invalidation whenever new commits occur.
   - Bypass or refresh using `--no-cache` or `--refresh`.

2. **Persistent Backup Store (`pending_commits_backup.db`)**:
   - Durable long-term snapshot persistence that survives cache expiration.
   - Dual-written on every live git scan.
   - Served offline or during troubleshooting via `--backup`.

## Table Columns & Tree View Output

The output renders a clean terminal table consolidating repository health metrics:
- **REPOSITORY** (22 chars): Repository slug or folder name.
- **VER/BRANCH** (16 chars): Tagged SemVer release and active local branch (e.g. `v6.523.1/main`).
- **UNCOMMITTED** (12 chars): Consolidated count of all uncommitted working files (untracked + modified + staged).
- **UNPUSHED** (10 chars): Number of local commits ahead of remote tracking branch.
- **STATUS** (9 chars): Status indicator (`● PEND` or `○ CLEAN`).

For dirty repositories, a structured tree view provides instant surgical fix commands:
```text
├── Option 1: git -C "<repo>" add -A && git commit -m "wip: save changes" && git push
└── Option 2: git -C "<repo>" stash -u
```

A global fleet remediation command is rendered in the footer when dirty repositories exist:
```text
Fleet Remediation: gitmap cpar "wip: save changes"
```

## Examples

### Example 1: Standard status overview (uses 90s SQLite cache)
```bash
gitmap pc
```

### Example 2: Force refresh cache with live git scan
```bash
gitmap pc --refresh
```

### Example 3: Completely bypass cache
```bash
gitmap pc --no-cache
```

### Example 4: Serve status from persistent backup DB
```bash
gitmap pc --backup
```

### Example 5: Output machine-readable JSON
```bash
gitmap pc --json
```

