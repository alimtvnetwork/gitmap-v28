# gitmap watch

Live-refresh repository status dashboard with configurable interval.

## Alias

w

## Usage

    gitmap watch [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --interval \<seconds\> | 30 | Refresh interval (min: 5) |
| --group \<name\> | — | Monitor only repos in a group |
| --no-fetch | false | Skip git fetch before status |
| --json | false | Output single snapshot as JSON |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Watch all repos with default interval

    gitmap watch

**Output:**

    Watching 42 repos (30s refresh) — Ctrl+C to stop
    ─────────────────────────────────────────────────
    REPO             BRANCH     STATUS  AHEAD/BEHIND
    my-api           main       clean   0/0
    web-app          develop    dirty   2/1
    billing-svc      main       clean   0/0
    auth-gateway     feature/x  dirty   5/0
    ─────────────────────────────────────────────────
    Refreshing in 28s...

### Example 2: Watch a group with fast refresh

    gitmap w --group backend --interval 10

**Output:**

    Watching 5 repos (group: backend, 10s refresh) — Ctrl+C to stop
    ─────────────────────────────────────────────────
    REPO             BRANCH   STATUS  AHEAD/BEHIND
    billing-svc      main     clean   0/0
    auth-gateway     main     clean   0/0
    payments-api     main     dirty   1/0
    user-svc         develop  clean   0/2
    notification-svc main     clean   0/0
    ─────────────────────────────────────────────────
    Refreshing in 8s...

### Example 3: Single JSON snapshot (no live mode)

    gitmap watch --json --no-fetch

**Output:**

    [
      {"name":"my-api","branch":"main","status":"clean","ahead":0,"behind":0},
      {"name":"web-app","branch":"develop","status":"dirty","ahead":2,"behind":1},
      {"name":"billing-svc","branch":"main","status":"clean","ahead":0,"behind":0}
    ]

## See Also

- [status](status.md) — One-time status snapshot
- [scan](scan.md) — Scan directories to populate the database
- [group](group.md) — Manage repo groups for filtered watching

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter watch
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
