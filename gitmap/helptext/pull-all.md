# gitmap pull-all

Batch-pull every tracked repository in the catalog. This is a thin
shorthand for `gitmap pull --all` — same resolver, same parallelism,
same pending-task accounting, same exit-code semantics. It exists so
the right-click context menu and shell history can name the fan-out
intent unambiguously.

## Alias

pa

## Usage

    gitmap pull-all [flags]

All `pull` flags are forwarded verbatim. `--all` is injected
automatically and is idempotent (passing it again is a no-op).

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --verbose | false | Enable verbose logging |
| --parallel \<N\> | 1 | Run up to N pulls concurrently (worker pool) |
| --only-available | false | Skip repos whose latest probe reports no new tag |
| --stop-on-fail | false | Halt the batch after the first failure |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Plain batch pull

    gitmap pull-all

**Output:**

    Pull batch (37 repos, parallel=1)
    [01/37] my-api .................. up to date
    [02/37] frontend ................ pulled (3 commits)
    ...
    37 ok · 0 failed

### Example 2: 8-way parallel, only repos with new commits

    gitmap pull-all --parallel 8 --only-available --stop-on-fail

## See also

- `pull` — single-repo / group-scoped pull (this is its `--all` form)
- Right-click context menu — Clone ▸ Pull all (Shift+right-click on
  Windows; confirm-gated dialog on macOS/Linux)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pull-all
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
