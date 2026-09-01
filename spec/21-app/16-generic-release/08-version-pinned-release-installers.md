# 08 — Version-Pinned Release Installers

> **Version:** 1.0.0
> **Updated:** 2026-04-20
> **Status:** Authoritative — implementers MUST follow this spec verbatim.

> **Related specs:**
> - [03-install-scripts.md](03-install-scripts.md) — generic installer pipeline (placeholder substitution, checksum, PATH)
> - [02-release-pipeline.md](02-release-pipeline.md) — release CI stages
> - [05-release-assets.md](05-release-assets.md) — asset naming and packaging
> - [06-release-metadata.md](06-release-metadata.md) — version normalization and tagging

---

## 1. Purpose & Contract

When a Git tag (`vX.Y.Z`) is pushed and a GitHub Release is created, the
release **MUST** ship two installer scripts whose version is **baked in
at publish time**:

| Asset on Release page | Platform | Behavior |
|----------------------|----------|----------|
| `install.ps1` | Windows / PowerShell 5.1+ | Installs **exactly** `vX.Y.Z` |
| `install.sh`  | Linux / macOS / WSL Bash 4+ | Installs **exactly** `vX.Y.Z` |

### The Hard Contract (non-negotiable)

1. **Spec-first ordering.** The spec at this path (and any companion
   spec it references) **MUST** be updated and merged to `main`
   *before* the release tag is pushed. Any AI or engineer reading
   only the released artifacts must be able to reconstruct intent
   from the spec without guessing.
2. **Version is embedded.** Each released installer carries a
   non-empty constant `EMBEDDED_VERSION = "vX.Y.Z"` at the top of
   the file. This value is written by CI, never by hand.
3. **Embedded version wins, always.** When `EMBEDDED_VERSION` is
   non-empty, the installer:
   - **MUST NOT** call the GitHub `releases/latest` API.
   - **MUST NOT** read any "latest" redirect.
   - **MUST NOT** switch to a different repository.
   - **MUST** download assets exclusively from
     `https://github.com/<repo>/releases/download/vX.Y.Z/…`.
4. **Determinism.** Re-running the same released installer six months
   later — even if newer releases exist — must still install
   `vX.Y.Z`. The installer is a snapshot, not a probe.
5. **Override is opt-in.** A user may pass `-Version vA.B.C` (PS) or
   `--version vA.B.C` (sh) to install a different pinned version.
   In that case the override wins, but the installer still does
   **not** consult "latest".
6. **Checksum is mandatory.** The installer downloads
   `checksums.txt` from the **same** release tag and verifies the
   archive SHA-256 before extraction. Mismatch ⇒ exit non-zero.

---

## 2. Why This Spec Exists (Background for handoff AIs)

Historically the repo's root-level `install.sh` and `install.ps1` are
**bootstrap** scripts living on `main`. They probe the GitHub API for
the latest tag and may even hop branches. That behavior is correct for
the canonical `raw.githubusercontent.com/.../main/install.sh` URL, but
**wrong** for installers attached to a specific GitHub Release.

A user who copies the install command from
`https://github.com/<repo>/releases/tag/v3.10.0` expects v3.10.0 — not
"whatever is latest now". This spec eliminates that ambiguity by
introducing a **separate, generated, version-pinned variant** that is
uploaded as a Release asset and is the *only* installer linked from
release notes.

---

## 3. File Layout

```
repo/
├── install.sh                       # bootstrap (main branch, follows latest)
├── install.ps1                      # bootstrap (main branch, follows latest)
├── linter-scripts/
│   └── installer-templates/
│       ├── install.pinned.sh.tmpl   # template → becomes release asset
│       └── install.pinned.ps1.tmpl  # template → becomes release asset
└── .github/workflows/
    └── release.yml                  # generates pinned installers per tag
```

The two bootstrap scripts at repo root remain unchanged. The two new
template files live under `linter-scripts/installer-templates/` and
are **never executed directly** — CI renders them per release.

---

## 4. Template Contract

Each template file is plain text with exactly these placeholder tokens:

| Token | Replaced With | Example |
|-------|---------------|---------|
| `__EMBEDDED_VERSION__` | The release tag, with leading `v` | `v3.11.0` |
| `__REPO_SLUG__`        | `owner/repo` | `alimtvnetwork/coding-guidelines-v24` |
| `__BUILD_DATE_UTC__`   | ISO-8601 UTC timestamp | `2026-04-20T07:42:11Z` |
| `__COMMIT_SHA__`       | Full 40-char commit SHA | `a1b2c3…` |

Substitution is a literal `sed`/`Replace-AString` pass — no template
engine, no logic. CI must fail the release job if any token remains
unreplaced after substitution.

### 4.1 Bash template (`install.pinned.sh.tmpl`) — required header

```bash
#!/usr/bin/env bash

# ──────────────────────────────────────────────────────────────────────

# install.sh — version-pinned release installer

#   Repo:    __REPO_SLUG__

#   Version: __EMBEDDED_VERSION__

#   Built:   __BUILD_DATE_UTC__   (commit __COMMIT_SHA__)

#

# This script installs EXACTLY __EMBEDDED_VERSION__.

# It will NEVER probe "latest" or switch repositories.

# ──────────────────────────────────────────────────────────────────────

set -euo pipefail

readonly EMBEDDED_VERSION="__EMBEDDED_VERSION__"
readonly EMBEDDED_REPO="__REPO_SLUG__"
readonly EMBEDDED_COMMIT="__COMMIT_SHA__"
```

### 4.2 PowerShell template (`install.pinned.ps1.tmpl`) — required header

```powershell

# ──────────────────────────────────────────────────────────────────────

# install.ps1 — version-pinned release installer

#   Repo:    __REPO_SLUG__

#   Version: __EMBEDDED_VERSION__

#   Built:   __BUILD_DATE_UTC__   (commit __COMMIT_SHA__)

#

# This script installs EXACTLY __EMBEDDED_VERSION__.

# It will NEVER probe "latest" or switch repositories.

# ──────────────────────────────────────────────────────────────────────

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$Script:EmbeddedVersion = '__EMBEDDED_VERSION__'
$Script:EmbeddedRepo    = '__REPO_SLUG__'
$Script:EmbeddedCommit  = '__COMMIT_SHA__'
```

---

## 5. Runtime Resolution Algorithm

Both installers MUST implement this exact resolution function. Naming
follows the project's positive-guard convention (P1).

```
resolve_target_version(cli_arg, embedded):
    if isNonEmpty(cli_arg):       return cli_arg          # explicit override
    if isNonEmpty(embedded):      return embedded         # pinned release
    return PROBE_LATEST                                   # bootstrap fallback
```

| Source             | Action                                        |
|--------------------|-----------------------------------------------|
| CLI `--version`    | Use as-is. **No** "latest" probe.             |
| `EMBEDDED_VERSION` | Use as-is. **No** "latest" probe.             |
| Both empty         | (Bootstrap-only path) call GitHub API.        |

The released installers — by definition — always have a non-empty
`EMBEDDED_VERSION`, so they can never reach the `PROBE_LATEST` branch.
That branch exists only so the same template could theoretically be
reused as a bootstrap; in practice the repo-root scripts handle
bootstrap and these templates do not.

### 5.1 Forbidden operations (released installer)

The pinned installer **MUST NOT** contain or call:

- `api.github.com/repos/.../releases/latest`
- `releases/latest/download/…` URLs
- Any `redirect`-following logic that resolves to a different tag
- Any code path that swaps `EMBEDDED_REPO` for another repository

CI MUST `grep` the rendered output and fail the release if any of the
following literal strings appear:

```
releases/latest
/releases/latest/download
api.github.com/repos
```

(Exception: these strings may appear inside a comment that is
explicitly prefixed with `# FORBIDDEN-OK:` for documentation
purposes. The grep MUST exclude that prefix.)

---

## 6. Download URL Construction

Both installers build URLs from the embedded values **only**:

```
BASE = https://github.com/${EMBEDDED_REPO}/releases/download/${VERSION}
ARCHIVE_URL  = ${BASE}/${ARCHIVE_NAME}
CHECKSUM_URL = ${BASE}/checksums.txt
```

`ARCHIVE_NAME` is computed from OS + arch detection per
[03-install-scripts.md §OS and Architecture Detection](03-install-scripts.md).
There is no fallback URL, no mirror, no CDN.

---

## 7. Release Pipeline Integration

Add a new stage between **Generate Scripts** and **Publish** in
[02-release-pipeline.md](02-release-pipeline.md):

```
… → 7. Generate Scripts → 7a. Render Pinned Installers → 8. Extract Changelog → 9. Publish
```

### 7a — Render Pinned Installers (CI step)

```yaml
- name: Render version-pinned installers
  env:
    VERSION:    ${{ needs.setup.outputs.version }}   # e.g. v3.11.0
    REPO_SLUG:  ${{ github.repository }}
    COMMIT_SHA: ${{ github.sha }}
  run: |
    mkdir -p release-assets
    BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

    for tmpl in linter-scripts/installer-templates/*.tmpl; do
      out_name="$(basename "${tmpl%.tmpl}")"            # install.pinned.sh / install.pinned.ps1
      final_name="${out_name#install.pinned.}"          # sh / ps1
      out_path="release-assets/install.${final_name}"

      sed -e "s|__EMBEDDED_VERSION__|${VERSION}|g" \
          -e "s|__REPO_SLUG__|${REPO_SLUG}|g" \
          -e "s|__BUILD_DATE_UTC__|${BUILD_DATE}|g" \
          -e "s|__COMMIT_SHA__|${COMMIT_SHA}|g" \
          "$tmpl" > "$out_path"

      # Fail if any placeholder survived
      if grep -q '__[A-Z_]\+__' "$out_path"; then
        echo "::error file=${out_path}::Unreplaced placeholder found"
        grep -n '__[A-Z_]\+__' "$out_path"
        exit 1
      fi

      # Fail if forbidden "latest" lookups slipped in
      if grep -E '(releases/latest|api\.github\.com/repos)' "$out_path" \
         | grep -v '# FORBIDDEN-OK:'; then
        echo "::error file=${out_path}::Forbidden latest-lookup detected"
        exit 1
      fi

      chmod +x "$out_path"
    done
```

The two rendered files (`release-assets/install.sh` and
`release-assets/install.ps1`) are then included in the release
asset list and folded into `checksums.txt` like any other asset.

---

## 8. Release Notes Wording

The release body MUST link to the **pinned** installers, not the
`raw.githubusercontent.com/.../main/` ones. Use this canonical block:

```markdown

## Install vX.Y.Z (exact version)

**Linux / macOS:**
```bash
curl -fsSL https://github.com/<repo>/releases/download/vX.Y.Z/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://github.com/<repo>/releases/download/vX.Y.Z/install.ps1 | iex
```

> These commands install **exactly vX.Y.Z**. They will not silently
> upgrade to a newer release. To install the current latest, use the
> bootstrap installer on `main` instead.
```

The wording above is generated from a CI template; do not hand-edit
release notes.

---

## 9. CLI Surface (released installer)

| Bash flag           | PowerShell param   | Effect                                          |
|---------------------|--------------------|-------------------------------------------------|
| `--version vA.B.C`  | `-Version vA.B.C`  | Override embedded pin (still no "latest" probe) |
| `--dir <path>`      | `-InstallDir <p>`  | Override install destination                    |
| `--arch <arch>`     | `-Arch <arch>`     | Override architecture detection                 |
| `--no-path`         | `-NoPath`          | Skip PATH modification                          |
| `--no-verify`       | `-NoVerify`        | **Disallowed.** Must exit 2 with explanation.   |
| `--help`            | `-Help`            | Print embedded version and usage                |

Note: `--no-verify` is intentionally **rejected** by the pinned
installer. Checksum verification is part of the contract.

`--help` output MUST include the embedded version on the first line:

```
install.sh — version-pinned installer for <repo> @ vX.Y.Z
```

---

## 10. Test Matrix (must pass before tagging)

CI runs these tests against the rendered installers in a clean
container before publishing the release:

| # | Scenario | Expected |
|---|----------|----------|
| 1 | Run `install.sh` with no args | Installs `vX.Y.Z` (the embedded one) |
| 2 | Run `install.sh --version vA.B.C` | Installs `vA.B.C` (override honored) |
| 3 | Block `api.github.com` via `/etc/hosts` and run `install.sh` | Installation still succeeds (proves no API probe) |
| 4 | Tamper with the downloaded archive | Script exits non-zero with checksum error |
| 5 | Run `install.sh --no-verify` | Script exits 2 with "verification cannot be disabled" |
| 6 | Same scenarios under PowerShell 5.1 and 7.x | All pass |

Test 3 is the **decisive** one — it proves the installer is fully
self-contained against the embedded version.

---

## 11. Spec-First Workflow (the release dance)

For every release, in order:

1. Update or extend any spec under `spec/16-generic-release/` that
   describes the change (this file, `03-install-scripts.md`,
   `06-release-metadata.md`, etc.).
2. Bump `package.json` version, run
   `node scripts/sync-version.mjs && node scripts/sync-spec-tree.mjs`.
3. Commit the spec + sync artifacts on `main`. **No tag yet.**
4. Push the commit; verify docs viewer renders the updated spec.
5. Only then push the tag `vX.Y.Z`. CI will:
   - Render the pinned installers from templates,
   - Verify no placeholder / no forbidden URL,
   - Run the §10 test matrix,
   - Publish the release with both installers as assets.
6. Confirm the release page's quick-install snippets match §8 and
   resolve to the correct pinned URLs.

If step 1 is skipped, **abort the release**. A handoff AI must be
able to read the spec at the tagged commit and reproduce the
installer behavior exactly.

---

## 12. Constraints (summary)

- Version is embedded at render time, never read from runtime
  environment, never resolved from "latest".
- Released `install.sh` / `install.ps1` are **separate files** from
  the repo-root bootstrap scripts. Never reuse the bootstrap as a
  release asset.
- `checksums.txt` from the **same** release tag is the only trust
  anchor. `--no-verify` is rejected.
- Spec at this path is updated and merged **before** the tag is
  pushed. Tag-then-spec is forbidden.
- The CI step at §7a MUST fail the release if any placeholder
  survives or any forbidden URL appears.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
