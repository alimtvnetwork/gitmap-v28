# Subtask 311-04: Installation Download Speed

## Objective

Measurably reduce the time to download and install gitmap from a release,
without breaking the installer or the self-update binary swap.

## Steps

1. Profile the current path: installer scripts (install.sh family at repo
   root) and cli/cmdupdate/ — how the release tarball and checksums.txt are
   fetched (connections, order, retries, resume), and where time goes.
2. Time a baseline full download of the current release asset.
3. Implement, in priority order and consistent with spec 242 (aria2c
   multi-connection download) and spec 286 (download speedup):
   - multi-connection segmented download (aria2c when present, built-in
     parallel range download otherwise);
   - checksums.txt fetched concurrently with the asset;
   - resume of partial downloads;
   - skip-download fast path when the target version is already cached.
4. Do not regress the known E1148 defect direction in
   cli/cmdupdate/update_binary_swap.go (the tarball ships
   gitmap-v<version>-<os>-<arch>, not a bare `gitmap` binary name).
5. Time the same download after the change; record before/after in
   02-spec/21-app/311-cicd-issues-and-end-to-end-testing/02-component-spec.md.

## Exit Criteria

- After timing is faster than before timing for the same asset and network.
- Installer Smoke CI jobs (linux + windows) green.
- No new CI failures introduced by the change.

## Profiling results + design (2026-10-23)

### Current path
- install.sh `download_asset()`: preflight HEAD -> asset download ->
  checksums.txt download, STRICTLY SERIAL. `download()` already uses aria2c
  (16 conn, 80 splits) when present; curl/wget are single-stream.
- Go `downloadReleaseStaging()` (cli/cmdupdate/update_binary_swap.go):
  archive, then checksums.txt — strictly serial. `downloadHTTP()`
  (update_ui_types.go) is a plain single-stream http.Get; the release
  candidate already carries AssetSize. Telemetry records engine +
  DurationMs, so before/after is measurable.
- Repo-root install.sh re-fetches cli/scripts/install.sh; check whether the
  two are byte-identical before editing (edit both if mirrored).

### Blocking defect found while profiling (fix in this subtask)
- release.yml packages Unix assets as
  `tar czf gitmap-v<ver>-<os>-<arch>.tar.gz gitmap-v<ver>-<os>-<arch>` — the
  tar's only entry carries the VERSIONED name.
- `extractTarGzBinary()` only accepts entries named exactly `gitmap` /
  `gitmap.exe` -> never matches -> E1148 on EVERY Unix self-update (same
  root cause as the `gitmap update` E1148 defect filed 2026-10-23; the zip
  path works because release.yml renames the Windows binary to gitmap.exe).
- Fix: accept `gitmap`, `gitmap.exe`, or any `gitmap-*` versioned entry.

### Concrete changes
1. Go: checksums.txt fetched concurrently with the archive in
   downloadReleaseStaging (goroutine + channel); E1138/E1139 preserved.
2. Go: segmented parallel archive download when AssetSize >= 4 MiB —
   4 Range segments written via WriteAt, aggregate progress, ANY range
   failure falls back to the existing single-stream downloadHTTP. New file
   cli/cmdupdate/update_download_parallel.go (< 300 lines, flat control
   flow for the nested-if/boolean linters; atomic counters only — the race
   suite does not cover cmdupdate, keep shared state minimal anyway).
3. Go: E1148 extractor fix (versioned entry names).
4. install.sh: background the checksums.txt download at the point the
   asset download starts in download_asset(); wait for it before the
   verify step; strict_fail semantics and exit codes preserved.

## Implementation status (2026-10-23)
- Items 1-4 DONE: update_download_parallel.go (probe + 4 segments +
  WriteAt + aggregated progress + single-stream fallback; go vet clean),
  concurrent staging in downloadReleaseStaging with the engine name
  threaded into UpdateResult telemetry, E1148 extractor fix
  (isReleaseBinaryName), install.sh concurrency applied to BOTH mirrored
  copies (bash -n clean, copies still byte-identical).
- install.ps1: concurrency DEFERRED deliberately. Its checksums fetch is
  the same serial pattern, but the only portable mechanism (Start-Job)
  spawns a full PowerShell process whose startup cost rivals the single
  RTT it would save for a ~1 KB file; the payload itself already uses
  aria2c multi-connection when present. Revisit only if update telemetry
  shows Windows checksum latency matters.
- Before/after timings: recorded in the 311 component spec evidence
  section at release time.
