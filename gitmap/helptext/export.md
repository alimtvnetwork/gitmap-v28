# gitmap export

Export the local database to a portable file.

## Alias

ex

## Usage

    gitmap export [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Export as JSON format |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Export database (default format)

    gitmap export

**Output:**

    Exporting 42 repos from profile 'default'...
    Including:
      42 repositories
       3 groups (backend, frontend, infra)
       5 aliases
    ✓ Exported to gitmap-export.json (12.4 KB)

### Example 2: Export as JSON

    gitmap ex --json

**Output:**

    Exporting 42 repos from profile 'default'...
    {
      "repos": [...],
      "groups": [...],
      "aliases": [...]
    }
    ✓ Exported to gitmap-export.json

### Example 3: Export from a different profile

    gitmap profile switch work
    gitmap export

**Output:**

    ✓ Switched to profile 'work'

    Exporting 18 repos from profile 'work'...
    ✓ Exported to gitmap-export.json (6.2 KB)

## See Also

- [import](import.md) — Import repos from an export file
- [scan](scan.md) — Scan directories to populate the database
- [profile](profile.md) — Manage database profiles

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter export
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
