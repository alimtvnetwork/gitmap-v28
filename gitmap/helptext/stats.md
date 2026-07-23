# gitmap stats

Show aggregated usage and performance metrics for gitmap commands.

## Alias

ss

## Usage

    gitmap stats [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- None (stats are recorded automatically)

## Examples

### Example 1: Show stats summary

    gitmap stats

**Output:**

    COMMAND          RUNS    AVG TIME    LAST RUN
    scan             15      2.3s        2025-03-10 14:30
    clone             8      12.1s       2025-03-09 10:15
    pull             32       1.8s       2025-03-10 15:05
    status           24       0.9s       2025-03-10 16:00
    release           5       8.4s       2025-03-10 16:20
    cd               41       0.1s       2025-03-10 16:25
    exec             12       3.2s       2025-03-08 09:30
    ─────────────────────────────────────────────
    Total: 137 executions across 7 commands

### Example 2: Stats as JSON

    gitmap ss --json

**Output:**

    [
      {"command":"scan","runs":15,"avg_ms":2300,"last_run":"2025-03-10T14:30:00Z"},
      {"command":"clone","runs":8,"avg_ms":12100,"last_run":"2025-03-09T10:15:00Z"},
      {"command":"pull","runs":32,"avg_ms":1800,"last_run":"2025-03-10T15:05:00Z"}
    ]

## See Also

- [history](history.md) — View command execution history
- [history-reset](history-reset.md) — Clear history data

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter stats
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
