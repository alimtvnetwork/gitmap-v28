# Command Specification: GitIgnore Split-DB Cache & Concurrency Governance

Spec Reference: [02-spec/21-app/201-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md](../../../02-spec/21-app/201-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md)  
Plan Reference: [.ai-memory/plans/pending/64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md](../../plans/pending/64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md)

---

## 1. Overview

This specification establishes the CLI commands, flags, and configuration interfaces for managing the `.gitignore` audit caching engine, throttling concurrency, and inspecting Split-DB state:
1. Inspection of current `.gitignore` audit intervals and Split-DB cache storage location.
2. Configuration of audit frequency via user settings (`gitignore_check_interval`).
3. Tabular inspection of cached repository audit statuses and elapsed durations.
4. Granular cache invalidation to force fresh repository audits.
5. Command-line duration overrides for individual pull and ignore operations.

---

## 2. Command Index & Usage

### 2.1 Configuration Subcommands
- `gitmap ignore config`: Displays active check interval, database file path, cache entry count, and usage subtext examples.
- `gitmap ignore config set interval <duration>`: Updates and persists the `.gitignore` audit check interval in user settings. Supported formats: `"24h"`, `"1d"`, `"12h"`, `"30m"`, `"0"` / `"off"` (disabled).
- `gitmap ig config`: Shorthand alias for `gitmap ignore config`.

### 2.2 Cache Management Subcommands
- `gitmap ignore cache`: Renders a formatted terminal table showing all cached repositories, last verification timestamp, status (`clean`, `has_issues`, `remediated`), and duration.
- `gitmap ignore cache clear`: Invalidates all cached repository checks, prompting a fresh scan on subsequent pull operations.
- `gitmap ignore cache clear --force` / `-f`: Flushes the cache without interactive confirmation.

### 2.3 Command Flags
- `--interval <duration>`, `-i <duration>`: Overrides the configured check interval for a single `gitmap pull`, `gitmap pa`, or `gitmap ignore` execution.
- `--force`, `-f`: Bypasses Split-DB cache checks and forces immediate audit across all resolved repositories.

---

## 3. CLI Examples

```bash
# Inspect current ignore audit settings and Split-DB status:
gitmap ignore config
gitmap ig config

# Configure audit check interval to once every 24 hours (default):
gitmap ignore config set interval 24h
gitmap ignore config set interval 1d

# Configure audit check interval to 12 hours:
gitmap ignore config set interval 12h

# Disable audit caching (audit on every run):
gitmap ignore config set interval off
gitmap ignore config set interval 0

# List cached repository audit records:
gitmap ignore cache

# Flush cache to force fresh audits on next pull:
gitmap ignore cache clear --force

# Run pull-all with custom cache interval override:
gitmap pa --interval 6h
gitmap pa --force
```
