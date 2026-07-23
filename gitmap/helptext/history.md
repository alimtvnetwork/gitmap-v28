# gitmap history

Show CLI command execution history with timestamps.

## Alias

hi

## Usage

    gitmap history [--limit N] [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --limit \<N\> | 20 | Number of entries to show |
| --json | false | Output as structured JSON |

## Prerequisites

- None (history is recorded automatically)

## Examples

### Example 1: Show recent command history

    gitmap history

**Output:**

     #   COMMAND                           TIMESTAMP
     1   scan D:\wp-work                   2025-03-10 14:30:12
     2   clone json --target-dir D:\proj   2025-03-10 14:32:45
     3   status --all                      2025-03-10 15:00:03
     4   pull --group backend              2025-03-10 15:05:18
     5   release --bump patch              2025-03-10 16:20:00
     6   cd my-api                         2025-03-10 16:25:33
     7   watch --interval 10               2025-03-10 16:30:00
    7 entries (showing last 20)

### Example 2: Show last 3 entries

    gitmap hi --limit 3

**Output:**

     #   COMMAND                           TIMESTAMP
     1   cd my-api                         2025-03-10 16:25:33
     2   watch --interval 10               2025-03-10 16:30:00
     3   release --bump minor --dry-run    2025-03-10 17:00:15
    3 entries

### Example 3: Export history as JSON

    gitmap history --json --limit 2

**Output:**

    [
      {"id":1,"command":"scan D:\\wp-work","timestamp":"2025-03-10T14:30:12Z"},
      {"id":2,"command":"clone json","timestamp":"2025-03-10T14:32:45Z"}
    ]

## See Also

- [history-reset](history-reset.md) — Clear command history
- [stats](stats.md) — View aggregated usage metrics
- [bookmark](bookmark.md) — Save commands for re-execution

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter history
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
