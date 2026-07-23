# gitmap history-reset

Clear the CLI command execution history.

## Alias

hr

## Usage

    gitmap history-reset [--confirm]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --confirm | false | Skip confirmation prompt |

## Prerequisites

- None

## Examples

### Example 1: Reset with interactive confirmation

    gitmap history-reset

**Output:**

    Current history: 42 entries
    This will permanently delete all command history.
    Continue? [y/N]: y
    ✓ History cleared (42 entries removed)

### Example 2: Reset without prompt (scripting)

    gitmap hr --confirm

**Output:**

    ✓ History cleared (42 entries removed)

### Example 3: Empty history

    gitmap history-reset

**Output:**

    No history to clear (0 entries).

## See Also

- [history](history.md) — View command history
- [db-reset](db-reset.md) — Reset the entire database
- [stats](stats.md) — View usage metrics (also cleared)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter history-reset
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
