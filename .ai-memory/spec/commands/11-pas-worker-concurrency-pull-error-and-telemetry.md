# Command Specification: PAS Worker Concurrency, Pull-Error Split-DB Subsystem & Machine Telemetry

Spec Reference: [02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)
Plan Reference: [.ai-memory/plans/pending/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../plans/pending/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)

## 1. Overview

This specification establishes the commands, flags, and operational contracts for:
1. Clean, dedicated machine telemetry reporting (`gitmap machine`).
2. Exact GitMap version inspection during fleet enqueue banners.
3. Concurrency worker and hand controls (`--w`, `--hand`, `--wwoh`, `paswh`).
4. Heartbeat terminal progress reporting (>30s) and structured output grouping (Updated, Dirty, Failed).
5. The Split-DB pull error inspection subsystem (`gitmap pull-error`).
6. Windows Credential Manager (`wincredman`) diagnostic detection and actionable recovery.

---

## 2. Command Index & Usage

### 2.1 Machine Telemetry (`gitmap machine`)
- **Commands:** `gitmap machine`, `gitmap machines`, `gitmap machine-name`, `gitmap hostname`
- **Description:** Renders local machine telemetry identity card (alias, hostname, IP address, OS platform, architecture, current user, GitMap version) without dumping unrelated command help menus.
- **Flags:**
  - `--json` (`-j`): Returns structured machine identity JSON.
  - `--ssh` (`-s`): Queries and tabulates machine identity across all configured SSH fleet nodes.
  - `-h`, `--help`: Explicitly displays command usage and subcommands.

### 2.2 Concurrency Flags & `paswh` Command
- **Command:** `gitmap paswh [N] [Y]`
- **Description:** Shorthand for `gitmap pull-all --ssh --w N --hand Y`.
  - Default when omitted: `N=1, Y=1` (1 worker with 1 hand).
  - Positional args: `gitmap paswh 2 2` launches 2 workers with 2 concurrent hands each.
- **Flags in `gitmap pa`, `gitmap pas`, `gitmap pull`:**
  - `--w <N>`, `--worker <N>`, `--workers <N>`: Sets concurrent worker pool count $N$.
  - `--hand <Y>`, `--h <Y>`, `--hands <Y>`: Sets concurrent tasks per worker $Y$.
  - `--worker-with-one-hand`, `--wwoh`: Locks execution to 1 hand per worker.

### 2.3 Pull-Error Inspection Subsystem (`gitmap pull-error`)
- **Commands:** `gitmap pull-error`, `gitmap pull-errors`, `gitmap pulle`, `gitmap pull-e`
- **Description:** Context-aware inspection of repository pull errors stored in the Split-DB `pull_errors` table (`data/pull-errors/sql.db`).
  - Inside a repo (`cwd` contains `.git`): automatically displays error history for current repository.
  - Outside a repo: defaults to `all` repositories from the latest pull run.
  - Positional argument `<repo-slug|path>`: filters errors specifically for that target.
- **Flags:**
  - `--ssh`: Aggregates pull error logs across remote fleet nodes with node IDs and GitMap versions.
  - `--json`: Outputs machine-readable error report for autonomous agent self-healing.
  - `--limit <N>`: Limits error entries displayed (default: 10).

### 2.4 Wincredman Diagnostic Remediation
- When `fatal: Unable to persist credentials with the 'wincredman' credential store` is detected during repository pull operations, GitMap automatically classifies the root cause and suggests:
  `↳ Next Step: Run 'gitmap fix-credential' (alias: fc) or check Windows Credential Manager service`
