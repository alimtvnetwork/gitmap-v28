# gitmap cluster history

Inspect the audit trail of past delegated cluster executions.

## Usage

```bash
gitmap cluster history [RunRef] [--limit N] [--json]
```

## Description

Displays execution records including run ID, command kind, target selector, timestamp, and per-node outcomes (OK, Failed, Skipped).

## Flags

- `[RunRef]`: Optional run ID (e.g. `RUN-20260817-001`) to inspect detailed node logs.
- `--limit <n>`: Limit number of historical runs displayed (default 20).
- `--json`: Output execution history in JSON format.

## Examples

```bash
gitmap cluster history
gitmap cluster history RUN-20260817-001
gitmap cluster history --limit 5 --json
```

See also: `gitmap cluster`, `gitmap servers-clients`
