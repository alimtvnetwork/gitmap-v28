---
name: find-next
description: gitmap find-next (alias fn) reads latest VersionProbe rows where IsAvailable=1 and prints repos with available updates. Read-only, supports --scan-folder filter and --json output. v3.9.0.
type: feature
---
# Find-Next (Phase 2.4, v3.9.0)

## Overview

`gitmap find-next` (alias `fn`) is the read-only consumer of the data that Phase 2.3's `gitmap probe` writes into the `VersionProbe` table. It surfaces every repo whose **latest** probe row reports `IsAvailable = 1`, sorted by `NextVersionNum DESC` so the freshest tags float to the top.

The command is intentionally read-only. To refresh stale results, the user runs `gitmap probe --all` first; `find-next` never triggers a probe itself. This separation keeps the latency model predictable: probe is the slow, network-bound command; find-next is an instant DB read.

## SQL

The query joins `Repo` against `VersionProbe` and uses a correlated subquery on `MAX(ProbedAt)` to pick only the newest probe per repo:

```sql
SELECT r.*, p.NextVersionTag, p.NextVersionNum, p.Method, p.ProbedAt
FROM Repo r
JOIN VersionProbe p ON p.RepoId = r.RepoId
WHERE p.IsAvailable = 1
  AND p.ProbedAt = (SELECT MAX(ProbedAt) FROM VersionProbe WHERE RepoId = r.RepoId)
ORDER BY p.NextVersionNum DESC, r.Slug ASC
```

A second variant (`SQLSelectFindNextByScanFolder`) adds `AND r.ScanFolderId = ?` for the `--scan-folder <id>` filter.

## CLI surface

| Flag | Effect |
|---|---|
| (none) | List every repo with an available update across the whole DB |
| `--scan-folder <id>` | Restrict to one ScanFolder (look up id via `gitmap sf list`) |
| `--json` | Emit `[]model.FindNextRow` as indented JSON for CI consumption |

Aliases: `fn` is registered alongside `find-next` in the dispatcher.

## Files

- `gitmap/cmd/findnext.go` — `runFindNext`, `parseFindNextFlags`, `emitFindNextJSON`, `emitFindNextText`
- `gitmap/store/find_next.go` — `(*DB).FindNext(scanFolderID int64) ([]FindNextRow, error)`
- `gitmap/model/find_next.go` — `FindNextRow` struct (embeds `ScanRecord`, adds tag/method/probedAt)
- `gitmap/constants/constants_find_next.go` — SQL, error/message strings, CLI tokens
- `gitmap/helptext/find-next.md` — `gitmap help find-next` content

## Phase 2.5 readiness

Phase 2.5 (parallel `gitmap pull`) will reuse `(*DB).FindNext` to pre-filter the worker queue: repos with no available probe row get skipped without spawning a pull worker at all.
