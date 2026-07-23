# 110 — `gitmap update` (remote-installer flow)

> **Status:** Implemented in v5.51.0. Supersedes the in-tree source-rebuild
> path that was the default in v5.50.x and earlier (now opt-in via
> `--source-rebuild`).

## Goal

`gitmap update` must work for users who installed gitmap from a released
artifact and never cloned the source repo. The legacy flow required a
local checkout, a Go toolchain, and a writable source path — none of
which a typical end user has.

## Behavior

1. `requireOnline()` — bail early if offline.
2. Unless `--source-rebuild` is present, run `runUpdateRemoteInstall()`:
   - Resolve the installer URL for the host OS:
     - **Windows** → `https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install.ps1`
     - **macOS / Linux** → `https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install.sh`
   - Download to `os.TempDir()/gitmap-update-*.{ps1|sh}` (UTF-8 BOM on
     Windows, `chmod 0755` on Unix).
   - Exec the script with the right shell (`powershell -ExecutionPolicy
     Bypass -NoProfile -NoLogo -File <path>` on Windows, `bash <path>`
     (or `sh`) elsewhere). Stdout/stderr/stdin are inherited so the
     installer's interactive prompts and progress lines flow through.
   - On non-zero exit, propagate the installer's exit code via
     `os.Exit(exitErr.ExitCode())`.
3. On any pre-exec failure (download error, etc.), print
   `MsgUpdateRemoteFallback` and fall through to the legacy
   source-rebuild path so power users with a local checkout still
   succeed.
4. `--source-rebuild` skips step 2 entirely and goes straight to the
   legacy handoff (`resolveRepoPath` + `createHandoffCopy` +
   `launchHandoff`).

## Why we don't probe sibling repos in Go

The downloaded `install.{ps1,sh}` already implements the parallel
`-v<N+i>` HEAD probe documented in
[`07-generic-release/09-generic-install-script-behavior.md`](../07-generic-release/09-generic-install-script-behavior.md):
20 parallel `HEAD` requests against `gitmap-v<current+1..current+20>`,
max-hit wins, then falls back to the latest release of the current
repo, then `main` HEAD. Reimplementing that loop inside `cmd/update.go`
would duplicate the probe contract and drift over time. The single
source of truth stays in the installer.

## Constants

Declared in `gitmap/constants/constants_update.go`:

| Constant | Purpose |
|----------|---------|
| `UpdateRemoteInstallerPwsh` | Windows installer raw URL |
| `UpdateRemoteInstallerBash` | Unix installer raw URL |
| `FlagSourceRebuild` | `--source-rebuild` opt-out into legacy flow |
| `MsgUpdateRemoteFetch` | "Fetching remote installer: %s" header |
| `MsgUpdateRemoteRun` | "Running installer: %s" header |
| `MsgUpdateRemoteDone` | Success footer |
| `MsgUpdateRemoteFallback` | Falling-back-to-source notice |
| `ErrUpdateRemoteDownload` | Download error |
| `ErrUpdateRemoteRun` | Installer exec error |

## Files

| File | Role |
|------|------|
| `gitmap/cmd/update.go` | `runUpdate()` dispatcher; chooses remote vs source-rebuild |
| `gitmap/cmd/updateremoteinstall.go` | `runUpdateRemoteInstall()`, download + exec helpers |
| `gitmap/constants/constants_update.go` | URL + flag + message constants |
| `install.ps1`, `install.sh` (repo root) | Canonical installers downloaded at runtime |

## Migration notes

- The `gitmap-updater` sidecar binary is **unchanged** and still works
  for source-rebuild deployments — it is the second-tier fallback when
  the remote install fails AND no local repo path can be resolved.
- No flag/positional changes for users: `gitmap update` with no args
  Just Works on a release-installed binary.
- CI: nothing to update. The legacy `runUpdateRunner` handoff path is
  still reachable via `--source-rebuild` and remains covered by the
  existing `updatecleanup_*_test.go` suite.
