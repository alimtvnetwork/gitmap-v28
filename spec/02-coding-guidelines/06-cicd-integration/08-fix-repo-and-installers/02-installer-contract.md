# `install.sh` / `install.ps1` & `release-install.*` — Normative Contract

> **Version:** 1.0.0 · **Updated:** 2026-04-28
> Implements: `install.sh`, `install.ps1`, `release-install.sh`, `release-install.ps1`
> Companion to: [`spec/14-update/27-generic-installer-behavior.md`](../../../14-update/27-generic-installer-behavior.md)

---

## 1. Two installers, two modes

| Script | Mode | Source | Resolves "latest"? |
|---|---|---|---|
| `install.sh` / `install.ps1` | **Implicit** (default) or **Pinned** (with `--version`) | `raw.githubusercontent.com/<repo>/<branch>` | Yes (when implicit) |
| `release-install.sh` / `release-install.ps1` | **Pinned only** (always) | GitHub Release asset | NEVER |

`release-install.*` is a thin pinning wrapper that hands off to
`install.*` with `--no-latest --version <tag> --pinned-by-release-install <tag>`.

---

## 2. Resolution precedence (release-install only)

Highest precedence first:

1. `--version <tag>` / `-Version <tag>` (CLI flag)
2. `$INSTALLER_VERSION` / `$env:INSTALLER_VERSION` (env var)
3. `__VERSION_PLACEHOLDER__` baked at release-asset build time

If two sources disagree, emit a `WARN` line and the higher-precedence
value wins. None present → exit `1`.

Version string MUST match `^v?\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$`.
Mismatch → exit `2`.

---

## 3. `install.sh` / `install.ps1` flag surface

| Bash | PowerShell | Purpose |
|---|---|---|
| `--repo owner/repo` | `-Repo owner/repo` | Override source repo |
| `--branch <name>` | `-Branch <name>` | Override branch (ignored if `--version`) |
| `--version vX.Y.Z` | `-Version vX.Y.Z` | Install a specific release tag (PINNED) |
| `--folders a,b,c` | `-Folders a,b,c` | Comma-separated folder list (subpaths OK) |
| `--dest <path>` | `-Dest <path>` | Install destination (default: cwd) |
| `--config <file>` | `-Config <file>` | Custom `install-config.json` |
| `--log-dir <path>` | `-LogDir <path>` | fix-repo log directory (default: `<dest>/.install-logs`) |
| `--show-fix-repo-log` | `-ShowFixRepoLog` | Print latest fix-repo log after run |
| `--max-fix-repo-logs N` | `-MaxFixRepoLogs N` | Keep only newest N logs (0 = keep all) |
| `--prompt` | `-Prompt` | Ask before overwriting each file (`y/n/a/s`) |
| `--force` | `-Force` | Overwrite all without prompting |
| `--dry-run` | `-DryRun` | Show what would change; write nothing |
| `--list-versions` | `-ListVersions` | List available tags and exit |
| `--list-folders` | `-ListFolders` | List top-level folders for ref and exit |
| `-n` / `--no-latest` / `--no-probe` | `-NoLatest` / `-NoProbe` | Skip latest-version probe |
| `--no-discovery` | `-NoDiscovery` | Skip V→V+N parallel discovery |
| `--no-main-fallback` | `-NoMainFallback` | Skip main-branch fallback |
| `--offline` / `--use-local-archive` | `-Offline` | Skip all network ops |
| `--run-fix-repo` | `-RunFixRepo` | Execute `fix-repo.{sh,ps1}` after verify |
| `-h` / `--help` | `-Help` / `-?` | Show help and exit 0 |

Unknown flag → exit `1` (Bash) / exit `1` (PowerShell). Both MUST print
the offending flag to stderr.

---

## 4. Banner (mandatory, both scripts)

Printed before any network call:

```
    📦 coding-guidelines installer
       mode:    <implicit|pinned>
       repo:    <owner/repo>
       version: <latest|vX.Y.Z>
       source:  <raw-branch|release-asset (tag)>
       dest:    <path>
```

`release-install.*` substitutes `release-install (pinned)` for the title.

---

## 5. Verification (mandatory unless `-n` passed to `linters-cicd/install.sh`)

1. Download `checksums.txt` from the same release/branch.
2. SHA-256 the downloaded archive.
3. Compare against the matching line in `checksums.txt`.
4. Mismatch → exit `4` (`E_VERIFY_FAILED`). Missing `checksums.txt` →
   `WARN` and continue (do **not** fail). With `-n` no exit `4` is ever
   raised.

---

## 6. Pinned-mode forbidden behaviors (release-install + `--version`)

Pinned mode MUST NOT:

- query `/releases/latest`
- fall back to the `main` branch
- cross repo boundaries (no V→V+N discovery)
- pick a "compatible" or "nearest" version
- silently downgrade to implicit mode

The dual-endpoint probe (`/releases/download/<tag>/...` then
`/archive/refs/tags/<tag>`) is allowed because both URLs resolve to the
**same pinned tag** — this is not a `main` fallback.

---

## 7. Exit codes (unified)

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Generic failure (missing tool, unknown flag, network exhausted, no version resolvable) |
| `2` | Offline mode required network OR invalid version string OR handshake mismatch |
| `3` | Pinned release / asset not found (PINNED MODE only) |
| `4` | Verification failed (checksum / required-paths) — only when verification is ON |
| `5` | Inner installer / handoff rejected |

---

## 8. fix-repo integration (`--run-fix-repo`)

When the flag is set (or `INSTALL_RUN_FIX_REPO=1`):

1. After verify+extract succeeds, invoke `fix-repo.sh` (POSIX) or
   `fix-repo.ps1` (Windows) at the install destination.
2. Stream output into `<log-dir>/fix-repo-<UTC-ISO8601>.log`.
3. If `--show-fix-repo-log` is set, `tail`/`Get-Content` the file to
   stdout after completion.
4. If `--max-fix-repo-logs N` is set with `N > 0`, prune oldest logs.
5. fix-repo's exit code is propagated 1:1 (its codes 2–7 surface as the
   installer's final exit code).

---

## 9. CODE RED compliance

- Banner-to-exit path ≤ 300 lines per script (excluding embedded help).
- Functions 8–15 effective lines; helpers split into `scripts/install/*`
  if needed.
- Zero nested conditionals; guard-and-return.
- Every error logged to stderr AND reflected in the exit code.
- No magic numbers in branching — every exit code is a named constant.
