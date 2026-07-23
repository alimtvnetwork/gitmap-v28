---
name: version-probe
description: Hybrid HEAD-then-clone version probe (v3.8.0+). v3.134.0 added a capped foreground worker pool. v3.135.0 unifies flag names (--probe-workers / --probe-depth) across `probe` and `scan`, deprecating --workers and --probe-concurrency.
type: feature
---
# Version Probe (Phase 2.3, v3.8.0)

## Overview

Two changes ship together in v3.8.0:

1. **`gitmap scan` auto-tags every discovered repo** with the `ScanFolderId` of the just-registered scan root. New helper: `cmd/scan.go::tagReposWithScanFolder` calls `EnsureScanFolder(absDir, "", "")` and then `db.TagReposByScanFolder(folder.ID, paths)`. Failures log to stderr but do NOT fail the scan.
2. **New `gitmap probe [<repo-path>|--all]` command** runs the hybrid HEAD-then-clone version probe and persists results into the `VersionProbe` table.

## Probe strategy

Order matters — fall through to the next strategy only when the previous one fails:

| # | Strategy | Command | When it fails |
|---|---|---|---|
| 1 | `ls-remote` | `git ls-remote --tags --sort=-v:refname <url>` | Server rejects unauthenticated probes, returns zero tags, or git exits non-zero |
| 2 | `shallow-clone` | `git clone --depth 1 --filter=blob:none --no-checkout <url>` into `os.MkdirTemp` then `git tag --sort=-v:refname` | Network/auth failure |

The shallow-clone fallback is **treeless** (`--filter=blob:none`) and **checkout-less** (`--no-checkout`) so we only pay for the refs database — no working tree, no blobs.

## Database

`store/version_probe.go` adds three methods on `*DB`:

- `TagReposByScanFolder(scanFolderID int64, paths []string) error` — bulk `UPDATE Repo SET ScanFolderId = ? WHERE AbsolutePath IN (?,?,?)` via interpolated placeholders. No-op when paths is empty.
- `RecordVersionProbe(model.VersionProbe) error` — inserts a row, mapping `IsAvailable bool` to `INTEGER 0|1`.
- `LatestVersionProbe(repoID int64) (model.VersionProbe, error)` — returns `sql.ErrNoRows` when no probe has run yet (caller handles).

## URL preference

`pickProbeURL` prefers `HTTPSUrl` over `SSHUrl` — HTTPS has less auth friction in CI / first-time-ever clones. SSH only kicks in when HTTPS is empty.

## Semver int

`probe.parseSemverInt` packs `vMAJOR.MINOR.PATCH` into `MAJOR*1e6 + MINOR*1e3 + PATCH` for use in `ORDER BY NextVersionNum DESC` queries. Pre-release suffixes (e.g. `1.2.3-rc1`) collapse to the numeric prefix only — display logic should always use `NextVersionTag`, never `NextVersionNum`.

## CLI surface

```
gitmap probe                                   # default: 2 workers, depth 1
gitmap probe --all                             # explicit "all"
gitmap probe E:\src\my-repo                    # single repo by path
gitmap probe --all --json                      # JSON, input order preserved
gitmap probe --all --probe-workers 3           # raise foreground pool (cap = 3)
gitmap probe --all --probe-depth 25            # deeper shallow-clone fallback
gitmap scan . --probe-workers 3 --probe-depth 5
```

Per-repo line format:
- `✓ <slug> → v1.2.3 (method=ls-remote)`
- `· <slug> → no new version (method=ls-remote)`
- `✗ <slug> → <error>`

Final summary: `✓ Probe complete: <available> available, <unchanged> unchanged, <failed> failed.`

## Foreground worker pool (v3.134.0)

`gitmap probe` runs through a capped worker pool:

- Default `--probe-workers 2` — sweet spot for residential bandwidth.
- Hard cap of 3 (`constants.ProbeMaxWorkers`); higher values clamp with a
  one-line stderr notice.
- Values < 1 are rejected at parse time.
- JSON output preserves input order regardless of completion order
  (each worker writes to its own pre-allocated slot in the entries
  slice). Human progress lines print as workers complete.
- Counter updates and the per-line print share a single `counterMu` so
  totals stay coherent and stdout lines never interleave mid-print.

## Unified probe flags (v3.135.0)

Both `gitmap probe` and `gitmap scan` now accept the same two value
flags. Naming is intentional — muscle memory carries between commands:

| Flag | Default | Applies to | Notes |
|---|---|---|---|
| `--probe-workers N` | 2 (probe) / 3 (scan) | probe + scan | Pool size. Probe caps at 3; scan keeps its existing default. |
| `--probe-depth N` | 1 | probe + scan | `git clone --depth N` for the shallow-clone fallback only. No effect on the `ls-remote` fast path. Coerced to `>=1` inside `tryShallowClone`. |
| `--workers N` | — | probe (deprecated) | Alias for `--probe-workers`; emits `MsgProbeWorkersAlias` once. |
| `--probe-concurrency N` | — | scan (deprecated) | Alias for `--probe-workers`; emits `MsgScanProbeConcurrencyAlias`. New flag wins when both are set. |

Depth is plumbed through `BackgroundRunner.SetCloneDepth` (called
BEFORE the first `Start` so workers observe the value race-free) and
`probe.RunOneWithDepth` for the foreground path. The legacy
`probe.RunOne(url)` shim still exists and forwards to depth=1 so any
external integration keeps working.

## Files

- `gitmap/probe/probe.go` — `RunOne` (legacy shim), `RunOneWithDepth`, `tryLsRemote`, `parseFirstTag`, `parseSemverInt`, `Result.AsModel`
- `gitmap/probe/clone.go` — `tryShallowClone(url, depth)`, `summarize`
- `gitmap/probe/background.go` — `BackgroundRunner` + `SetCloneDepth`, `SetFailureHook`
- `gitmap/store/version_probe.go` — DB methods
- `gitmap/cmd/probe.go` — dispatcher + `runProbePool` / `probeWorker`
- `gitmap/cmd/probeflags.go` — `parseProbeArgs`, `--probe-workers` / `--workers` (deprecated) / `--probe-depth`
- `gitmap/cmd/probereport.go` — `executeOneProbe(db, repo, depth)`, JSON shaping, `tallyProbe`
- `gitmap/cmd/rootflags.go` — `parseScanFlags` + `resolveScanProbeOptions` reconciles deprecated alias
- `gitmap/cmd/scanbackgroundprobe.go` — calls `runner.SetCloneDepth(opts.Depth)` before enqueue
- `gitmap/constants/constants_probe.go` — `ProbeDefaultWorkers=2`, `ProbeMaxWorkers=3`, `ProbeDefaultDepth=1`, all flag tokens + deprecation messages
