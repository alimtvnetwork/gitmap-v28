# gitmap import

Import repositories from an export file into the local database.

## Alias

im

## Usage

    gitmap import <file>

## Flags

None.

## Prerequisites

- An export file from `gitmap export` (see export.md)

## Examples

### Example 1: Import from an export file

    gitmap import gitmap-export.json

**Output:**

    Importing from gitmap-export.json...
    [1/42] my-api... added
    [2/42] web-app... added
    [3/42] billing-svc... added
    ...
    ✓ 42 repos imported
    ✓ 3 groups restored (backend, frontend, infra)
    ✓ 5 aliases restored

### Example 2: Import with duplicates

    gitmap im backup.json

**Output:**

    Importing from backup.json...
    [1/15] my-api... skipped (already exists)
    [2/15] web-app... skipped (already exists)
    [3/15] new-service... added
    [4/15] analytics-api... added
    ...
    ✓ 15 repos processed (5 added, 10 skipped)

### Example 3: Import into a fresh profile

    gitmap profile create new-machine
    gitmap profile switch new-machine
    gitmap import team-repos.json

**Output:**

    ✓ Switched to profile 'new-machine'
    Importing from team-repos.json...
    ✓ 28 repos imported
    ✓ 4 groups restored

## See Also

- [export](export.md) — Export the database to a file
- [scan](scan.md) — Scan directories as an alternative to import
- [profile](profile.md) — Manage database profiles

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter import
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
