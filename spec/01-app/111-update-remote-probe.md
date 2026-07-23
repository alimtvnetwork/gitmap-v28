# 111 — `gitmap update` Native Probe + Remote Installer

> **Status:** Planned for v5.52.0. Refines [110](110-update-remote-installer.md)
> by moving the `-v<N+i>` sibling-repo probe out of the downloaded
> `install.{ps1,sh}` and into Go. The remote installer is still exec'd, but
> only after Go has picked the winning repo slug.

## Goal

Single source of truth for "which `gitmap-vN` repo is newest?" lives in Go,
not in two shell scripts. The downloaded installer is reduced to its
mechanical role: fetch + place the binary for the slug Go hands it.

## Resolution Contract

`resolveLatestRepoSlug()` returns `(slug, source, err)` where `source` is
one of `"sibling-probe"`, `"current-release"`, `"current-main"`.

1. **Sibling probe** — parse current slug (`gitmap-v27`) → base
   (`gitmap`) + `N=23`. Fire `UpdateProbeMaxSiblings` (20) parallel
   `HEAD https://github.com/<owner>/gitmap-v<N+i>` requests with a
   per-request timeout of `UpdateProbeTimeoutSec` (5s). Any 2xx counts as
   a hit. **Max-hit index wins** (e.g. if v24, v25, v27 all 200, pick
   v27). Returns `("gitmap-v27", "sibling-probe", nil)`.
2. **Current-release fallback** — no sibling hits → `GET
   api.github.com/repos/<owner>/<current>/releases/latest`. If 200,
   return current slug with source `"current-release"`.
3. **Main HEAD fallback** — releases API fails → `HEAD
   https://github.com/<owner>/<current>` against `main`. Return current
   slug with source `"current-main"`. This always succeeds for a valid
   binary.

## Install Contract

Once a slug is resolved:

1. Download `https://raw.githubusercontent.com/<owner>/<slug>/main/install.ps1`
   (Windows) or `install.sh` (Unix) to a temp file.
2. Exec with inherited stdio (`powershell -ExecutionPolicy Bypass ...` /
   `bash <path>`). The installer no longer needs to probe — Go already
   resolved the target — so it installs from its own repo directly.
3. Propagate the installer's exit code.

## Flags

| Flag | Behavior |
|------|----------|
| (default) | Probe → install winning slug |
| `--probe-only` | Print resolved slug + source, exit 0 — no install |
| `--no-probe` | Skip probe, install from current repo only |
| `--source-rebuild` | Unchanged. Bypass remote path entirely, use legacy in-tree rebuild |

## Logging

```
■ Probing sibling repos (gitmap-v27..gitmap-v43)...
  hit: gitmap-v27 (HTTP 200)
  hit: gitmap-v27 (HTTP 200)
✓ Resolved: gitmap-v27 (source: sibling-probe)
■ Fetching installer: https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install.ps1
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Install succeeded |
| 1 | Probe resolved nothing AND fallbacks all failed |
| 2 | Installer download failed |
| (installer exit) | Propagated from `install.{ps1,sh}` |

## Constants (`gitmap/constants/constants_update.go`)

| Constant | Value |
|----------|-------|
| `UpdateProbeMaxSiblings` | `20` |
| `UpdateProbeTimeoutSec` | `5` |
| `UpdateProbeRepoBase` | `"gitmap"` |
| `UpdateRepoOwner` | `"alimtvnetwork"` |
| `UpdateRepoHEADTmpl` | `"https://github.com/%s/%s"` |
| `UpdateRawInstallerTmpl` | `"https://raw.githubusercontent.com/%s/%s/main/%s"` |
| `UpdateReleasesAPITmpl` | `"https://api.github.com/repos/%s/%s/releases/latest"` |
| `FlagProbeOnly` | `"--probe-only"` |
| `FlagNoProbe` | `"--no-probe"` |

## Why Go (not the shell installer) Owns the Probe

- Single contract: install.ps1 and install.sh historically drifted. Moving
  the probe to Go eliminates two parallel implementations.
- Testable: `httptest.Server` covers every branch in CI without network.
- Observable: probe lifecycle is logged through the standard
  `MsgUpdate*` constants and can be captured by `--verbose`.

## Files

| File | Role |
|------|------|
| `gitmap/cmd/updateprobe.go` | `resolveLatestRepoSlug`, `probeSiblings`, fallback chain |
| `gitmap/cmd/updateremoteinstall.go` | Calls probe, then download+exec |
| `gitmap/cmd/update.go` | Flag dispatch (`--probe-only`, `--no-probe`) |
| `gitmap/cmd/updateprobe_test.go` | httptest coverage |
| `gitmap/constants/constants_update.go` | All probe + URL constants |
