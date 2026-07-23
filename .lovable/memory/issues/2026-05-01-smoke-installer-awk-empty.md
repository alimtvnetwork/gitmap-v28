# RCA — installer-smoke (release) failed with "Binary not found at $DEST/gitmap"

**Date:** 2026-05-01  
**Workflow:** `.github/workflows/release.yml` → job `installer-smoke` → `.github/scripts/smoke-installer.sh release`  
**Symptom:**

```
gitmap v4.3.0
Install summary
    Binary: /tmp/tmp.6o6e2hVgqU/install/gitmap-cli/gitmap
    Install dir: /tmp/tmp.6o6e2hVgqU/install/gitmap-cli
...
Done!
Error: Binary not found or not executable at /tmp/tmp.6o6e2hVgqU/install/gitmap
Error: Process completed with exit code 3.
```

The installer wrote the binary at `…/install/gitmap-cli/gitmap` (correct,
canonical layout). The smoke harness then asserted on `…/install/gitmap`
(legacy unwrapped layout) and exited 3.

## Root cause

`smoke-installer.sh#load_deploy_manifest` parsed `gitmap-cli` and `gitmap`
out of `gitmap/constants/deploy-manifest.json` using awk with `-F'"'`:

```bash
APP_SUBDIR="$(awk -F'"' '/"appSubdir"/ {print $4; exit}' "$manifest_path")"
BINARY_NAME_UNIX="$(awk -F'"' '/"unix"/ {print $4; exit}' "$manifest_path")"
```

Locally that produces `gitmap-cli` / `gitmap`. On the GitHub `ubuntu-latest`
runner this expression intermittently returned the empty string (most
likely an awk locale/quoting interaction with the embedded `"` field
separator under `set -euo pipefail` + process substitution), which made
`$APP_SUBDIR=""`. The candidate loop then collapsed:

```
"$DEST/$APP_SUBDIR/$BINARY_NAME_UNIX"  =>  "$DEST//gitmap"   (no match)
"$DEST/$BINARY_NAME_UNIX"              =>  "$DEST/gitmap"    (no match)
```

…and the legacy loop also couldn't find a match. The `find` fallback was
only reached when the first two were empty, but the failure already short-
circuited because the bash word-split made the first candidate look
"populated enough" to break the loop on the test runner. Either way, the
script reported a phantom `$BIN=$DEST/gitmap` path that the installer
never wrote.

The install was **fine** — only the test harness's manifest reader was
fragile.

## Fix

`.github/scripts/smoke-installer.sh`

1. **Replaced the awk-only manifest reader with a tiered parser:**
   `jq` (always present on `ubuntu-latest`) → `python3` (always present)
   → awk fallback. JSON is parsed by an actual JSON parser, not by
   `-F'"'` field tricks.
2. **Logged the resolved manifest values** (`APP_SUBDIR`, `BINARY_NAME_UNIX`,
   `LEGACY_APP_SUBDIRS`) before the candidate loop so any future drift is
   diagnosable from the CI log alone.
3. **Always run the `find` fallback** when direct candidates miss — it is
   now a guaranteed safety net regardless of what the manifest reader
   returned.

The PowerShell counterpart (`smoke-installer.ps1`) was already robust
because it uses `ConvertFrom-Json` plus a recursive `Get-ChildItem` fallback.
No change needed there.

## Why this didn't bite earlier

This is the first release where the smoke script was relying on the
manifest reader for the canonical path. Earlier versions hardcoded
`$DEST/gitmap-cli/gitmap` and only fell back to the manifest for legacy
migration cases, so the empty-`APP_SUBDIR` bug masked itself.

## Permanent guarantees

- JSON config files in CI scripts must be parsed with `jq` or `python3`,
  never with `awk -F'"'`. The latter is acceptable only as a last-resort
  fallback **and** only when paired with a non-empty validation that
  defaults on miss.
- Any path-resolution loop in CI must end with an exhaustive `find` (or
  `Get-ChildItem -Recurse`) fallback before declaring "not found". Direct
  candidates are an optimization, not a contract.

## Files touched

- `.github/scripts/smoke-installer.sh` — tiered JSON parser + always-on
  `find` fallback + manifest diagnostic line.
