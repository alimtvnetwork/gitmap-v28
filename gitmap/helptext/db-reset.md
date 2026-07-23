# gitmap db-reset

Reset the local SQLite database, removing all tracked repos and metadata.

## Alias

None

## Usage

    gitmap db-reset [--confirm]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --confirm | false | Skip confirmation prompt |

## Prerequisites

- None

## Examples

### Example 1: Reset with interactive confirmation

    gitmap db-reset

**Output:**

    Current profile: default
    Data to be removed:
      42 repositories
       3 groups (backend, frontend, infra)
       5 aliases
      12 release records
      65 history entries

    This will permanently delete all data. Continue? [y/N]: y
    ✓ Database reset complete
    → Run 'gitmap scan' to rebuild

### Example 2: Reset without prompt

    gitmap db-reset --confirm

**Output:**

    ✓ Database reset (42 repos, 3 groups, 5 aliases removed)
    → Run 'gitmap scan' to rebuild

### Example 3: Reset an empty database

    gitmap db-reset

**Output:**

    Database is already empty. Nothing to reset.

## See Also

- [scan](scan.md) — Re-scan to rebuild the database
- [history-reset](history-reset.md) — Clear command history only
- [profile](profile.md) — Manage profiles (reset affects current profile)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-reset
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
