# Command Usage Statistics

## Overview

The `stats` command aggregates data from the `CommandHistory` table to
show usage patterns, performance metrics, and failure rates for all
gitmap commands.

---

## Command

### `gitmap stats` (alias: `ss`)

Display aggregated command usage statistics.

**Synopsis:**

```
gitmap stats [--command <name>] [--json]
```

**Flags:**

| Flag               | Description                              | Default |
|--------------------|------------------------------------------|---------|
| `--command <name>` | Show stats for a specific command only   | (all)   |
| `--json`           | Output as JSON                           | false   |

**Output fields:**

| Field       | Description                            |
|-------------|----------------------------------------|
| Command     | Command name                           |
| Runs        | Total number of executions             |
| Success     | Count of successful runs (exit code 0) |
| Fail        | Count of failed runs (exit code ≠ 0)   |
| Fail%       | Failure rate as percentage             |
| Avg(ms)     | Average execution duration             |
| Min(ms)     | Fastest execution                      |
| Max(ms)     | Slowest execution                      |
| Last Used   | Timestamp of most recent execution     |

**Summary row** shows overall totals: total executions, unique commands,
success/fail counts, overall failure rate, and average duration.

**Examples:**

```bash
# Show all command stats
gitmap stats
gitmap ss

# Stats for scan command only
gitmap stats --command scan

# JSON output for scripting
gitmap stats --json
```

---

## Data Source

All statistics are computed from the `CommandHistory` table using SQL
aggregation queries (`COUNT`, `AVG`, `MIN`, `MAX`, `SUM`). No additional
tables are required.

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_stats.go` | SQL queries, messages, formatting |
| `model/stats.go` | CommandStats and OverallStats structs |
| `store/stats.go` | Stats query methods |
| `cmd/stats.go` | Stats command handler |

---

## Constraints

- Reuses existing `CommandHistory` table — no schema changes.
- All files under 200 lines, all functions 8–15 lines.
- PascalCase for all SQL column references.
