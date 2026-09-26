# gitmap pull-all

Batch-pull every tracked repository in the catalog. By default, `gitmap pa` /
`gitmap pull all` / `gitmap pull-all` runs in **fast concise mode**, skipping slow
post-pull remote branch, PR, and tag subprocess inspections and printing an
aligned bullet summary with total duration.

To render the full post-pull status table, pass `--status` (or `--table`), or run
`gitmap pat`, `gitmap pull-all-table`, or `gitmap pull all table`.

## Alias

pa, pat (status table), pull-all-table (status table)

## Usage

    gitmap pull-all [flags]
    gitmap pa [--status | --table | --json]
    gitmap pat [flags]
    gitmap pull all table [flags]

All `pull` flags are forwarded verbatim. `--all` is injected
automatically and is idempotent (passing it again is a no-op).

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --status, --table | false | Render the full post-pull repository status table |
| --json | false | Suppress interactive progress bar and emit JSON summary to stdout |
| --verbose | false | Enable verbose logging |
| --parallel \<N\> | 0 (auto) | Run up to N pulls concurrently (worker pool) |
| --only-available | false | Skip repos whose latest probe reports no new tag |
| --stop-on-fail | false | Halt the batch after the first failure |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Default fast batch pull

```bash
gitmap pa
```

**Output:**

    • my-api                      up-to-date
    • frontend                    +12/-3 (2 files)

  ✓ Pull all complete: 37 pulled (1.4s)

### Example 2: Full post-pull status table (`--status` or `pat`)

```bash
gitmap pa --status
gitmap pat
gitmap pull all table
```

### Example 3: Machine-readable JSON summary (`--json`)

```bash
gitmap pa --json
```

### Example 4: 8-way parallel, only repos with new commits

```bash
gitmap pull-all --parallel 8 --only-available --stop-on-fail
```

## See also

- `pull` — single-repo / group-scoped pull (this is its `--all` form)
- Right-click context menu — Clone ▸ Pull all (Shift+right-click on
  Windows; confirm-gated dialog on macOS/Linux)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pull-all
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
