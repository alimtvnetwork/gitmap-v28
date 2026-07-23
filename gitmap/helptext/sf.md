# sf

Manage scan folders — the absolute roots that `gitmap scan` was invoked
against.

## Synopsis

```
gitmap sf add <absolute-path> [--label <text>] [--notes <text>]
gitmap sf list                                            # alias: ls
gitmap sf rm  <absolute-path|id>                          # alias: remove
```

## What it does

Every time `gitmap scan` runs, it auto-registers the scan root as a
`ScanFolder` row and tags every discovered repo with that
`ScanFolderId`. The `sf` subcommands let you inspect, label, and
remove those registrations without re-scanning.

`add` is **idempotent**: re-running it against the same path bumps
`LastScannedAt` and only overwrites `Label` / `Notes` when the new
values are non-empty (so a manually set label survives subsequent
scans).

`rm` performs a **detach-then-delete** transaction: every repo
pointing at the folder is first set to `ScanFolderId = NULL`, then the
`ScanFolder` row is deleted. Repos themselves are never deleted —
they just lose their scan-folder pointer.

## Subcommands

### `sf add <path>`

| Flag | Effect |
|---|---|
| `--label <text>` | Friendly name shown in `sf list` (e.g. "Work projects") |
| `--notes <text>` | Free-form notes for the operator |

The path is resolved to an absolute path before insert.

### `sf list` / `sf ls`

Prints every registered scan folder, newest-scanned first, with live
repo counts.

### `sf rm <path|id>` / `sf remove <path|id>`

Removes a scan folder by **absolute path** or by **numeric id** (look
up ids via `sf list`). Reports how many repos were detached.

## Examples

```
$ gitmap sf list
Scan folders (2):
  [1] E:\src
      label: Work | repos: 12 | last scanned: 2026-04-19 06:11:42
  [2] D:\experiments
      label: (none) | repos: 3 | last scanned: 2026-04-15 21:03:08

$ gitmap sf add E:\new-roots --label "Forks"
✓ Registered scan folder: E:\new-roots (id=3)

$ gitmap sf rm 2
✓ Removed scan folder: D:\experiments (id=2, 3 repos detached)
```

## Why it matters

`ScanFolderId` powers the `--scan-folder <id>` filter on `gitmap
find-next`, letting you query "what's new in `E:\src`" without seeing
unrelated repos. It also makes future bulk operations (Phase 2.6:
`gitmap cn next all --scan-folder <id>`) scope cleanly.

## See also

- `gitmap scan` — auto-registers scan folders as a side effect
- `gitmap find-next --scan-folder <id>` — filter updates by folder
- `gitmap probe` — populate the `VersionProbe` data that find-next reads

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter sf
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
