# Command History & Audit Trail

## Overview

Every CLI command execution is automatically logged to the `CommandHistory`
SQLite table, creating a queryable audit trail of all gitmap operations.
Users can review, filter, and reset the history via dedicated commands.

---

## Table Schema

### CommandHistory Table

| Column     | Type    | Constraints               | Notes                              |
|------------|---------|---------------------------|------------------------------------|
| Id         | TEXT    | PRIMARY KEY               | Timestamp-based unique ID          |
| Command    | TEXT    | NOT NULL                  | Command name (e.g. `scan`)         |
| Alias      | TEXT    | DEFAULT ''                | Alias used (e.g. `s`)             |
| Args       | TEXT    | DEFAULT ''                | Positional arguments               |
| Flags      | TEXT    | DEFAULT ''                | Flags passed (e.g. `--mode ssh`)   |
| StartedAt  | TEXT    | NOT NULL                  | RFC3339 start timestamp            |
| FinishedAt | TEXT    | DEFAULT ''                | RFC3339 end timestamp              |
| DurationMs | INTEGER | DEFAULT 0                 | Execution time in milliseconds     |
| ExitCode   | INTEGER | DEFAULT 0                 | 0 = success, non-zero = failure    |
| Summary    | TEXT    | DEFAULT ''                | Human-readable result summary      |
| RepoCount  | INTEGER | DEFAULT 0                 | Number of repos affected           |
| CreatedAt  | TEXT    | DEFAULT CURRENT_TIMESTAMP |                                    |

**Insert strategy:** Every command execution inserts a new row at start,
then updates it with completion details (duration, exit code, summary).

---

## Commands

### `gitmap history` (alias: `hi`)

Display the command execution audit log.

**Synopsis:**

```
gitmap history [--detail basic|standard|detailed] [--command <name>] [--limit N] [--json]
```

**Flags:**

| Flag                 | Description                                | Default    |
|----------------------|--------------------------------------------|------------|
| `--detail <level>`   | Output detail: basic, standard, detailed   | standard   |
| `--command <name>`   | Filter by command name                     | (all)      |
| `--limit N`          | Show only the last N entries               | 0 (all)    |
| `--json`             | Output as JSON                             | false      |

**Detail levels:**

- **basic** — Command name, timestamp, status (OK/FAIL)
- **standard** — Command, timestamp, flags, status, duration
- **detailed** — Command, timestamp, args, flags, status, duration, repo count, summary

**Examples:**

```bash
# Show recent history (standard detail)
gitmap history
gitmap hi

# Show basic view
gitmap history --detail basic

# Show detailed view for scan commands
gitmap history --detail detailed --command scan

# Last 10 entries as JSON
gitmap history --json --limit 10
```

### `gitmap history-reset` (alias: `hr`)

Clear all command history from the database.

**Synopsis:**

```
gitmap history-reset --confirm
```

Requires `--confirm` to prevent accidental deletion.

**Examples:**

```bash
gitmap history-reset --confirm
gitmap hr --confirm
```

---

## Audit Hook

The audit system works via a two-phase approach:

1. **Start phase** — Before the command runs, a record is inserted with
   command name, args, flags, and start timestamp.
2. **End phase** — After the command completes, the record is updated
   with finish timestamp, duration, exit code, summary, and repo count.

The audit hook is non-blocking — if the database is unavailable, the
command still executes normally. Audit failures are silently ignored.

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_history.go` | SQL, command names, messages |
| `model/history.go` | CommandHistoryRecord struct |
| `store/history.go` | History CRUD operations |
| `cmd/history.go` | History command (display) |
| `cmd/historyreset.go` | History-reset command |
| `cmd/audit.go` | Audit hook (start/end recording) |

---

## Constraints

- Non-blocking: audit failures never prevent command execution.
- PascalCase table and column names.
- `--confirm` required for history-reset (consistent with db-reset).
- `db-reset --confirm` also clears CommandHistory.
- All files under 200 lines, all functions 8–15 lines.
