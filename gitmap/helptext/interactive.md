# Interactive TUI

Launch a full-screen interactive terminal UI for browsing, searching,
and managing repositories.

## Alias

i

## Usage

    gitmap interactive [--refresh <seconds>]
    gitmap i [--refresh <seconds>]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--refresh` | config or 30 | Dashboard auto-refresh interval in seconds |

## Prerequisites

- Run `gitmap scan` at least once to populate the database (see scan.md)
- Terminal must support alternate screen mode

## Examples

### Example 1: Launch the TUI

    gitmap interactive

**Output:**

    Launching interactive TUI (42 repos loaded)...
    ┌─ Repos ──────────────────────────────────────┐
    │ > my-api         main    clean   0/0         │
    │   web-app        develop dirty   2/1         │
    │   billing-svc    main    clean   0/0         │
    │ 42 repos │ Tab: switch view │ /: search      │
    └──────────────────────────────────────────────┘

### Example 2: Launch with custom refresh interval

    gitmap i --refresh 10

**Output:**

    Launching interactive TUI (42 repos loaded)...
    Auto-refresh: every 10 seconds
    ┌─ Status ─────────────────────────────────────┐
    │ REPO           BRANCH    STATUS  AHEAD/BEHIND │
    │ my-api         main      clean   0/0          │
    │ web-app        develop   dirty   2/1          │
    │ 2 dirty │ 40 clean │ Refreshing in 10s        │
    └──────────────────────────────────────────────┘

### Example 3: No repos in database

    gitmap interactive

**Output:**

    ✗ No repositories found in database.
    → Run 'gitmap scan <directory>' first to populate the database.

## Key Bindings

    q / Esc      Quit
    Tab          Switch view (Repos → Actions → Groups → Status)
    j / ↓        Move down
    k / ↑        Move up
    Space        Select/deselect repo
    /            Focus search

## Prerequisites

- Run `gitmap scan` at least once to populate the database (see scan.md)
- Terminal must support alternate screen mode

## See Also

- [list](list.md) — List tracked repositories
- [group](group.md) — Manage repo groups
- [status](status.md) — View repo statuses
- [watch](watch.md) — Live-refresh status dashboard

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter interactive
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
