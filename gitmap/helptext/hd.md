# gitmap hd

One-screen operator dashboard: pending tasks, upgradable repos, last scan,
current gitmap version, and the tail of completed tasks.

`hd` is a read-only alias for `gitmap help --dashboard`. It never touches the
network unless `--refresh` is passed.

## Alias

hd, help-dashboard

## Usage

    gitmap hd [flags]
    gitmap help --dashboard [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Emit the dashboard payload as JSON |
| --refresh | false | Re-probe upgrades before rendering (equivalent to `gitmap list --update --stale-days 0`) |
| --tail \<n\> | 5 | How many recent completed tasks to show |

## Prerequisites

- Local gitmap database (any prior scan)

## Rendered View

    ╔══════════════════════════════════════════════════════════╗
    ║  gitmap dashboard                          v6.80.0       ║
    ╚══════════════════════════════════════════════════════════╝

      Version
        installed:  v6.80.0
        latest:     v6.80.0   (up to date)

      Repositories
        scanned:      142
        upgradable:    22       →  gitmap list --update
        unknown:        2

      Pending
        queued:         3       →  gitmap pending
        oldest:  #48  Upgrade   acme/api  (2h ago)

      Last scan:  2026-07-22 10:14   (13h ago)

      Recent completions
        #47  Clone   acme/api          2h ago
        #46  Scan    /home/a/work      13h ago
        #45  Pull    acme/web          1d ago

      Next steps
        • gitmap update all       Upgrade every upgradable repo
        • gitmap do-pending       Retry pending tasks
        • gitmap scan .           Re-scan the current directory

## JSON Output

    $ gitmap hd --json
    {
      "version": {"installed":"v6.80.0","latest":"v6.80.0","behind":0},
      "repos":   {"scanned":142,"upgradable":22,"unknown":2},
      "pending": {"queued":3,"oldest":{"id":48,"type":"Upgrade","target":"acme/api","ageSeconds":7200}},
      "lastScanAt": "2026-07-22T10:14:00Z",
      "recent": [
        {"id":47,"type":"Clone","target":"acme/api","completedAt":"..."},
        {"id":46,"type":"Scan","target":"/home/a/work","completedAt":"..."}
      ]
    }

## Examples

Human-readable one-screen view:

```text
$ gitmap hd
gitmap dashboard  v6.80.1  (up to date)
scanned=142  upgradable=22  pending=3
last scan: 13h ago
```

JSON payload for scripts:

```json
{"version":{"installed":"v6.80.1","latest":"v6.80.1","behind":0},"repos":{"scanned":142,"upgradable":22,"unknown":2},"pending":{"queued":3}}
```

Refresh probes before rendering:

```text
$ gitmap hd --refresh --tail 10
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Dashboard rendered |
| 2 | Local database is missing (run `gitmap scan` first) |

## See Also

- [stats](stats.md) — Deeper aggregate metrics
- [list-update](list-update.md) — Full upgradable list
- [pending](pending.md) — Full pending queue
- [history](history.md) — Complete command history
- [dashboard](dashboard.md) — Interactive HTML dashboard for a single repo

## Scripting (JSON)

    gitmap help --json --filter hd

Schema: `spec/08-json-schemas/hd.schema.json` (v6.80.0+).
