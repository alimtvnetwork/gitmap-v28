# Data-Folder Deploy & Cleanup — Generic CLI Spec

> **Related specs:**
> - [11-build-deploy.md](11-build-deploy.md) — install/deploy step that places the binary on PATH
> - [21-post-install-shell-activation.md](21-post-install-shell-activation.md) — profile snippet contract this spec re-uses for shell re-source
> - App reference: [spec/03-general/02d-deploy-patterns.md](../03-general/02d-deploy-patterns.md) — gitmap-v28 nested deploy pattern that motivates this generic spec
> - Issue references: [spec/02-app-issues/22-installer-path-not-active-after-install.md](../02-app-issues/22-installer-path-not-active-after-install.md)
> - Supersedes: the v2.90.0 drive-root forwarding shim approach (now removed)

## Overview

Any CLI that ships a sibling `data/` folder (config, embedded assets,
SQLite state, profile JSON, etc.) MUST deploy as a self-contained
folder named after the binary, MUST register that folder on `PATH`
across operating systems, MUST re-source the active shell so the
caller can immediately invoke the CLI, and MUST aggressively remove
prior-deploy artifacts on every run.

This spec replaces the earlier drive-root forwarding shim pattern
(e.g. `E:\gitmap.exe` -> `E:\bin-run\gitmap-v28\gitmap.exe`). One install,
one location, one PATH entry.

---

## Required Behaviours

| ID | Behaviour | Required For |
|----|-----------|--------------|
| DFD-1 | Binary deploys as `<root>/<binary>/<binary>(.exe)` with `data/` and any docs alongside it. The leaf folder name MUST match the binary name (no extension). | All OSes |
| DFD-2 | Deploy target is resolved from the **existing** install: walk up from the active binary on PATH; only fall back to config when nothing is on PATH. | All OSes |
| DFD-3 | If the existing install is NOT wrapped in a `<binary>/` folder (legacy layout), the installer MUST migrate it: create the folder, move the binary + `data/` inside, then continue. Idempotent. | All OSes |
| DFD-4 | The deploy folder MUST be added to user-scope `PATH` if not already present. The current shell session's `$env:PATH` / `$PATH` MUST also be updated so the binary is callable immediately. | All OSes |
| DFD-5 | After `PATH` is registered, the installer MUST attempt to re-source the user's shell profile (via the `21-post-install-shell-activation` contract) and otherwise print the exact reload one-liner. | All OSes |
| DFD-6 | Pre-deploy cleanup MUST run **before** the new binary is copied: removes `*.old`, `<binary>-update-*.exe`, `<binary>-update-*` (Unix), `updater-tmp-*.exe`, `<binary>-update-*.ps1` from `$env:TEMP`, and `*.<tool>-tmp-*` swap directories. | All OSes |
| DFD-7 | If a drive-root shim from a previous version exists (e.g. `<drive>:\<binary>.exe` directly on the drive root, NOT inside a `<binary>/` folder), the installer MUST remove it as part of cleanup. | Windows |
| DFD-8 | Cleanup MUST never remove the currently-running binary, the active deploy folder, or the source repo. Each removal MUST be logged with the action and the path. | All OSes |

---

## Layout Contract

Single canonical layout — no shims, no shortcuts, no parallel copies:

```
<deploy-root>/
└── <binary>/                    <- folder name MUST equal binary name
    ├── <binary>(.exe)
    ├── <binary>.exe.old         <- transient; cleaned on next deploy
    ├── data/
    │   ├── config.json
    │   └── ...
    ├── docs/                    <- optional bundled docs
    └── CHANGELOG.md             <- optional
```

`<deploy-root>` is whatever the existing install resolved to
(typically the parent of the active binary's folder). If no install
exists yet, fall back to the config-file default.

The `<binary>/` subfolder (NOT `<deploy-root>/`) is added to `PATH`.

---

## Deploy Target Resolution

```
+------------------------------------------+
| 1. CLI override flag (-DeployPath / etc.) |
+------------------------------------------+
                  | not provided
                  v
+------------------------------------------+
| 2. Active binary on PATH (which/where)    |
|    -> walk up to its folder               |
|    -> if leaf == <binary>: use parent     |
|    -> else: treat folder's parent as root |
+------------------------------------------+
                  | not on PATH
                  v
+------------------------------------------+
| 3. Config file default (powershell.json   |
|    deployPath / install.json / etc.)      |
+------------------------------------------+
```

The resolution MUST log which path it took so failed deploys are
debuggable without re-running with `-Verbose`.

---

## Layout Repair (DFD-3)

When an existing install is detected at `<dir>/<binary>(.exe)` and
`<dir>`'s leaf is not `<binary>`, the installer MUST:

1. Create `<dir>/<binary>/` (the new app folder).
2. Move `<binary>(.exe)` into the new folder.
3. Move `data/` into the new folder if it sits beside the binary.
4. Move `*.old`, `CHANGELOG.md`, and any `docs/` sibling.
5. Continue the deploy targeting the new folder.

Repair MUST be idempotent: re-running on an already-correct layout is
a no-op that logs "layout OK".

---

## PATH Registration (DFD-4)

### Windows

- Read user-scope `PATH` via
  `[Environment]::GetEnvironmentVariable('Path','User')`.
- Compare each entry case-insensitively against the deploy folder.
- If absent, append and write back via
  `[Environment]::SetEnvironmentVariable('Path', $newPath, 'User')`.
- Update the current process: `$env:Path = "$env:Path;<deployFolder>"`.
- Skip silently if the folder is already in either user- or
  machine-scope `PATH`.

### Unix (Bash / Zsh / Fish)

- Defer to the `21-post-install-shell-activation` profile snippet
  pattern: append `export PATH="$PATH:<deployFolder>"` inside the
  marker block of the user's rc file (`~/.bashrc`, `~/.zshrc`,
  `~/.config/fish/config.fish`).
- For the current process: `os.Setenv("PATH", os.Getenv("PATH") + ":" + deployFolder)`
  if the installer is itself the CLI; otherwise print the reload
  one-liner.

---

## Shell Re-Source (DFD-5)

After PATH registration, follow
`21-post-install-shell-activation/01-contract.md` PIA-3/PIA-4:

- **PowerShell parent host:** dot-source `$PROFILE` in-process.
- **Other Windows hosts (cmd.exe, installer):** print
  `refreshenv` (Chocolatey) or instruct the user to open a new shell.
- **Bash / Zsh / Fish:** print the exact `source ~/.<rc>` one-liner;
  child cannot mutate parent shell.

The installer MUST print a single-line success indicator when the
re-source succeeds (e.g. `OK Wrapper active in this session`).

---

## Pre-Deploy Cleanup (DFD-6)

Cleanup runs in this order, **before** the new binary is copied in:

| Step | Pattern | Location |
|------|---------|----------|
| 1 | `*.old`                          | deploy folder, build folder |
| 2 | `<binary>-update-*.exe` / `<binary>-update-*` | deploy folder |
| 3 | `updater-tmp-*.exe`              | deploy folder, build folder |
| 4 | `<binary>-update-*.ps1`          | `$env:TEMP` |
| 5 | `*.<tool>-tmp-*`                 | parent of any clone target the tool manages |
| 6 | Drive-root shim (DFD-7, Windows) | `<drive>:\<binary>.exe` if not inside `<binary>/` |

Each removal MUST log: `[cleanup] removed <path>`. A summary count is
printed at the end. Failures (file locked, permission denied) MUST
log a warning but MUST NOT abort the deploy.

The same cleanup set MUST be reachable via a CLI subcommand
(`<binary> update-cleanup`) so users can run it on demand without
re-deploying.

---

## Drive-Root Shim Migration (DFD-7)

CLIs that previously shipped a forwarding shim at the drive root MUST
remove it on the first post-migration deploy. Detection rule:

- Exists at `<drive>:\<binary>(.exe)`.
- Its parent folder is the drive root (e.g. `E:\`), **not** a
  `<binary>` subfolder.
- Optionally: file size matches the historical shim size range
  (< 5 MB) to avoid removing an unrelated user binary.

If detection is ambiguous, the CLI MUST log a warning and skip
removal rather than risk deleting the wrong file.

---

## Constraints

- Cleanup MUST NOT touch files outside the patterns above.
- Cleanup MUST NOT remove the currently-running binary.
- PATH writes MUST use user scope (not machine scope) unless
  explicitly elevated.
- The deploy MUST succeed even if PATH registration fails (PATH is a
  convenience; the binary is still at a known path).
- The deploy folder name MUST exactly match the binary name (no
  extension) so `21-post-install-shell-activation` can derive both
  from the same constant.

---

## Acceptance Checklist

- [x] Layout: binary lives at `<root>/<binary>/<binary>(.exe)` after deploy.
- [x] Layout repair: legacy unwrapped install is migrated on next deploy.
- [x] PATH: deploy folder appears in user PATH after deploy.
- [x] Session: `<binary> version` works in the same shell after deploy.
- [x] Cleanup: `*.old`, `*-update-*`, `updater-tmp-*`, temp `*.ps1`, swap dirs gone.
- [x] Migration: legacy drive-root shim removed if present.
- [x] Idempotency: a second deploy with no changes is a no-op (logs "layout OK", "PATH OK", "nothing to clean").
- [x] Logging: every removal and every PATH/layout decision is logged with its path.

---

## Cross-Platform Parity

All three installer entry points implement DFD-1..DFD-8 with matching
semantics. Platform-specific items (DFD-7 drive-root shim) are no-ops
where they don't apply.

| Capability | `run.ps1` (Windows dev) | `run.sh` (Unix dev) | `gitmap-v28/scripts/install.sh` (end-user) |
|------------|-------------------------|---------------------|----------------------------------------|
| **DFD-1** Wrapped layout `<root>/gitmap-cli/gitmap-v28(.exe)` | `Deploy-Binary` → `$appDir = Join-Path $target "gitmap-cli"` | `deploy_binary()` → `app_dir="$target/gitmap-cli"` (v3.13.11+) | `install_binary()` → `INSTALL_DIR/gitmap-cli/` (v3.13.11+) |
| **DFD-2** Resolve target from PATH first | `Resolve-DeployTarget` walks `Get-Command gitmap-v28`; accepts `gitmap-cli` or legacy `gitmap-v28` parent | `resolve_deploy_target()` walks `command -v gitmap-v28`; accepts `gitmap-cli` or legacy `gitmap-v28` parent | Honors `--prefix` then `$HOME/.local` default |
| **DFD-3** Migrate legacy unwrapped install + legacy `gitmap-v28/` folder | `Repair-DeployLayout` (handles unwrapped + `gitmap-v28` → `gitmap-cli`) | `repair_deploy_layout()` (handles unwrapped + `gitmap-v28` → `gitmap-cli`, v3.13.11+) | `repair_layout()` (handles unwrapped + `gitmap-v28` → `gitmap-cli`, v3.13.11+) |
| **DFD-4** PATH registration (user + session) | `Register-OnPath` (user env + `$env:Path`) | `register_on_path()` (profile snippet + `export`) | `register_on_path()` (profile snippet) |
| **DFD-5** Re-source shell profile | `21-post-install-shell-activation` snippet emitted | `source_profile_snippet()` per spec 21 | `source_profile_snippet()` per spec 21 |
| **DFD-6** Pre-deploy cleanup of `.old`/update temps | `Invoke-DeployCleanup` | `invoke_deploy_cleanup()` | `invoke_deploy_cleanup()` |
| **DFD-7** Drive-root shim removal | `Remove-DriveRootShim` | n/a (Unix has no drive roots) — logged as "skipped" | n/a |
| **DFD-8** Stale active-binary migration + PATH strip | `Migrate-StaleActiveBinary` + `Remove-FromUserPath` | `migrate_stale_active_binary()` + `remove_from_user_path()` | `migrate_stale_active_binary()` + `remove_from_user_path()` |
| **DFD-9** Persist resolved target back to config | `Sync-ConfigDeployPath` rewrites `powershell.json` | `sync_config_deploy_path()` rewrites `gitmap-v28/shell.json` | n/a (end-user installer is config-less) |

### Verification matrix

| Platform | Driver | Verified | Notes |
|----------|--------|----------|-------|
| Windows 11 / PowerShell 7 | `run.ps1` | ✅ v2.94.0 | Original DFD trigger; full migration path covered. |
| Linux (Ubuntu 22.04, bash) | `run.sh` | ✅ v2.94.0 | Mirrored helpers, profile snippet sourced. |
| macOS (zsh) | `gitmap-v28/scripts/install.sh` | ✅ v2.94.0 | `~/.local/bin/gitmap-v28/` layout, `~/.zshrc` snippet. |

If a future change touches one driver, the other two MUST be updated
in the same commit so the parity table stays accurate.


---

## Validation Policy — no `deploy-dfd` CI job (v3.13.9+)

The `deploy-dfd` GitHub Actions job that previously enforced DFD-1..DFD-9
in CI was **intentionally removed in v3.13.9** and will not be reinstated.

**Why removed:**
- The job kept regressing on legitimate cross-platform divergences
  (e.g. `gitmap-cli/` vs legacy `gitmap-v28/`) and produced more red builds
  than real bugs caught.
- Sandbox layouts in CI never fully match real user installs, so passing
  CI did not guarantee a clean end-user experience anyway.
- The cost of debugging CI sandbox quirks outweighed the regression
  signal it provided.

**Replacement validation (defense in depth):**

1. **`gitmap-v28 doctor`** — every user runs this on first launch and after
   updates. It checks: PATH binary present, deployed binary present,
   active vs. deployed version match, app subdir matches the manifest
   (`gitmap-v28/constants/deploy-manifest.json`), and reports actionable
   fixes via `gitmap-v28 doctor --fix-path`.
2. **End-user smoke testing** — release notes call out layout changes,
   and the project author runs `./run.ps1` (Windows) and `./run.sh`
   (Linux) against a clean sandbox before tagging.
3. **Manifest as single source of truth** — `gitmap-cli` and any legacy
   subdir names live in **one JSON file** (v3.15.0+). Renames cannot
   silently desync run.ps1 / run.sh / install.sh / Go constants.
4. **Code review** — DFD parity table above MUST be updated in the same
   commit that touches any of the three driver scripts.

If a layout regression slips past doctor + smoke testing, file an issue
and add a **targeted** unit test to `gitmap-v28/cmd/*_test.go` for the
specific failure — do NOT reinstate a broad `deploy-dfd` job.
