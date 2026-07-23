# 08 — Pinned-Version Install Snippet (Release-Page Contract)

> **Audience:** NEA and any future AI/maintainer cutting a GitHub release.
> **Status:** Active since gitmap-v27 **v3.12.0** (2026-04-20).
> **Related:** `spec/01-app/94-install-script.md`,
> `spec/01-app/95-installer-script-find-latest-repo.md`,
> `spec/07-generic-release/03-install-scripts.md`.

## 1. Goal

When a user copies an install snippet from a **specific GitHub release page**
(e.g. `…/releases/tag/v3.11.1`), running it MUST install exactly **that
tag**, regardless of:

* whether a newer tag has since been published,
* whether a newer **versioned sibling repo** exists
  (`gitmap-v27`, `gitmap-v27`, …),
* whether the user is offline from the GitHub releases API.

The snippet is the **single source of truth** for "give me v3.11.1". It is
the contract between the release page and the installer scripts.

## 2. What gets appended to the release body

`gitmap-v27/release/installsnippet.go::AppendPinnedInstallSnippet` runs inside
`uploadToGitHub` (in `workflowgithub.go`) **after** `DetectChangelog()`
and **before** `CreateGitHubRelease`. It appends a markdown block to the
release body, gated by a hidden HTML marker so re-runs are idempotent:

```html
<!-- gitmap-pinned-install-snippet:v3.11.1 -->
## Install this exact version (v3.11.1)

… powershell + bash code-fences …
```

The exact template lives in `constants_release.go` as
`ReleaseSnippetTemplate` / `ReleaseSnippetMarker`.

### 2.1 PowerShell snippet (rendered)

```powershell
$ver = 'v3.11.1'
$installer = irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap-v27/scripts/install.ps1
& ([scriptblock]::Create($installer)) -Version $ver -NoDiscovery
```

### 2.2 Bash snippet (rendered)

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap-v27/scripts/install.sh \
  | bash -s -- --version v3.11.1 --no-discovery
```

Both snippets pin the version twice over: explicit `--version` + explicit
`--no-discovery`. Either alone is sufficient (see §3); shipping both
makes the contract self-documenting.

## 3. Installer-side contract

Both `install.ps1` and `install.sh` enforce the contract symmetrically:

| Flag set | Behavior |
|---|---|
| `-Version v3.11.1` (ps1) / `--version v3.11.1` (sh) | Skip the `releases/latest` API call. Skip versioned-repo discovery (the `-v<N>` sibling probe). Download `…/releases/download/v3.11.1/…` directly. |
| `-NoDiscovery` (ps1) / `--no-discovery` (sh) | Skip the discovery probe even if no version is pinned. |
| Neither | Pre-existing behavior: probe `…-v<N+1>`, `…-v<N+2>`, … then `releases/latest`. |

The pinned-version short-circuit was added in **v3.12.0**:

* `gitmap-v27/scripts/install.ps1` lines around the `INSTALLER_DELEGATED`
  branch — new `elseif (-not [string]::IsNullOrWhiteSpace($Version))`
  arm prints `[discovery] -Version <tag> pinned; skipping repo probe`.
* `gitmap-v27/scripts/install.sh` lines around the same branch — new
  `elif [ -n "${VERSION}" ]` arm prints
  `[discovery] --version <tag> pinned; skipping repo probe`.

## 4. NEA / AI handoff checklist

When cutting a new release tag:

1. **Bump version** in `gitmap-v27/constants/constants.go` and
   `src/constants/index.ts`.
2. **Update `CHANGELOG.md`** and `src/data/changelog.ts`.
3. **Update this spec** if either the snippet template or the installer
   contract changes. If only the rendered output changes (new flags,
   new repo URL), update §2.1/§2.2.
4. **Run `gitmap-v27 release`**. The publisher auto-appends the pinned
   snippet — do not paste it manually.
5. **Verify on the release page** that the snippet renders the correct
   tag (e.g. `v3.12.0`) and points at `gitmap-v27` (not `gitmap-v27`).

## 5. Test contract

Negative case for CI (future work, not yet wired):

* Construct a release body, call `AppendPinnedInstallSnippet(body, "v9.9.9")`
  twice. Assert the marker appears exactly once.
* Run `bash install.sh --version v9.9.9 --no-discovery --help` and assert
  the discovery probe log line does NOT appear.

## 6. History

| Version | Change |
|---|---|
| v3.12.0 | Initial spec + implementation. Snippet auto-appended; both installers honor `--version` / `-Version` to skip discovery + latest lookup. Repo renamed `gitmap-v27` → `gitmap-v27` everywhere. |
