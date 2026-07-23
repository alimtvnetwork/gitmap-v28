# gitmap db-migrate

Run pending database schema migrations on the active gitmap profile.

## Aliases

`dbm`

## Usage

    gitmap db-migrate [--verbose]

## What it does

1. Opens the active-profile SQLite database
   (`.gitmap/output/data/gitmap[-profile].db`).
2. Re-runs every `CREATE TABLE IF NOT EXISTS` and column-migration step.
3. Prints a single-line summary; warnings (if any) are streamed to
   stderr with the offending **table**, **column**, and **action** so
   they can be diagnosed without trial-and-error.

The command is **idempotent** — re-running it on an up-to-date schema
performs no writes. It is safe to invoke at any time, in any working
directory, on any OS.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--verbose` | false | Print extra context about what was checked. |

## When to run it

- Immediately after `gitmap update` — the post-update worker now does
  this automatically, but `db-migrate` lets you re-trigger it manually
  if the auto-run was skipped (e.g. update aborted, DB locked).
- After restoring `.gitmap/output/data/` from a backup or another
  machine that ran an older release.
- Whenever you see a warning such as
  `⚠ Migration failed: table=TempReleases column=Commit ...`.

## Examples

    gitmap db-migrate
    gitmap dbm --verbose

## Errors

| Condition | Exit | Message |
|-----------|------|---------|
| Cannot open the active-profile DB | 1 | `Error: gitmap db-migrate failed: ...` |
| ALTER TABLE rejected by SQLite | 0 | (warning printed, schema kept consistent) |

The command **does not** modify your data; it only adds missing
columns, indexes, and seed rows.

## See also

- `gitmap db-reset --confirm` — destructive: drop and rebuild every
  table from scratch.
- `gitmap update` — pulls the latest binary and runs `db-migrate`
  automatically afterwards.
- `gitmap doctor` — broader environment + schema health check.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-migrate
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
