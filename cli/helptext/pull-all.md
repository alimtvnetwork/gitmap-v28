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
    gitmap pa [--probe] [--status | --table | --json]
    gitmap pa --probe [-y]
    gitmap pat [flags]
    gitmap pull all table [flags]

All `pull` flags are forwarded verbatim. `--all` is injected
automatically and is idempotent (passing it again is a no-op).

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --probe | false | Probe GitHub account & workspace for companion repositories (`repo-secrets` & `repo-cache`) before pulling |
| --status, --table | false | Render the full post-pull repository status table |
| --json | false | Suppress interactive progress bar and emit JSON summary to stdout |
| --ssh, -s | false | Dispatch pull-all across SSH cluster fleet and Local VM concurrently |
| --verbose | false | Enable verbose logging |
| --parallel \<N\> | 0 (auto) | Run up to N pulls concurrently (worker pool) |
| --only-available | false | Skip repos whose latest probe reports no new tag |
| --stop-on-fail | false | Halt the batch after the first failure |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Default fast batch pull (shows only active/updated repos)

```bash
gitmap pa
```

**Output:**

    • frontend                    +12/-3 (2 files)
    • backend-api                 dirty

  ✓ Pull all complete: 64 pulled (2 active, 62 up-to-date) (1.4s)

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

### Example 4: SSH fleet batch pull (`--ssh`)

```bash
gitmap pa --ssh
```

**Output:**

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [alpha-win] (10.20.0.11): Enqueued (async)
    • Remote Node [beta-linux] (10.20.0.12): Enqueued (async)
    • Local VM (127.0.0.1 - localhost): Running locally

  ▶ Local VM (127.0.0.1 - localhost): 64 pulled (1 active, 63 up-to-date)
      • wp-exam-v2    +1266/-131 (10)

  ▶ Node [alpha-win] (IP: 10.20.0.11): 64 pulled (all up-to-date)
      (all repositories are up-to-date)

### Example 5: 8-way parallel, only repos with new commits

```bash
gitmap pull-all --parallel 8 --only-available --stop-on-fail
```

### Example 6: Pull all with companion repository probing (`--probe`)

Probe your remote GitHub account and local workspace for companion repositories (`repo-secrets` and `repo-cache`), prompting to clone any missing repositories before pulling all repositories:

```bash
gitmap pa --probe
gitmap pa --probe -y
gitmap pull all --probe
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
