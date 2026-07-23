# gitmap reset

Hard reset: physically delete the active profile's SQLite database file,
recreate the schema from scratch, and reapply all JSON-based seeds.

This is more destructive than `db-reset`, which only drops and recreates
tables inside the existing file. Use `reset` when the file itself is
suspected to be corrupted, or when you want a guaranteed-fresh start.

## Alias

None

## Usage

    gitmap reset [--confirm]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --confirm | false | Required. Confirms permanent deletion of the DB file |

## What it does

1. Locates the active profile's DB file under
   `gitmap-output/data/<profile>.db`.
2. Deletes the file from disk (missing file is treated as success).
3. Opens a fresh DB at the same path, which rebuilds all tables via
   `Migrate()` and reseeds `ProjectTypes` and `TaskTypes`.
4. Reapplies optional JSON seeds (e.g. `data/seo-templates.json`) when
   present in the working directory.

## Prerequisites

- None

## Examples

### Example 1: Reset without prompt

    gitmap reset --confirm

**Output:**

    Removed database file: gitmap-output/data/gitmap.db
    Reseeded data/seo-templates.json
    Reset complete: database file deleted, schema rebuilt, seeds reapplied.

### Example 2: Reset without confirmation

    gitmap reset

**Output:**

    Error: this will permanently delete the database file and rebuild it from scratch.
    Run with --confirm to proceed: gitmap reset --confirm

### Example 3: Reset when no DB file exists

    gitmap reset --confirm

**Output:**

    Reset complete: database file deleted, schema rebuilt, seeds reapplied.

## See Also

- [db-reset](db-reset.md) — Drop & recreate tables without deleting the file
- [scan](scan.md) — Re-scan to repopulate repos after reset
- [profile](profile.md) — Manage profiles (reset only affects active profile)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter reset
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
