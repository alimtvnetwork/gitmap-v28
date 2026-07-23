---
name: Release-Version Install Script
description: Dedicated `release-version.ps1` / `.sh` used ONLY by /release/:version pages. Pinned to an exact version, never auto-upgrades, never falls back to latest silently. Ships in two forms — generic (parameterized) + per-version snapshot (baked) — both per release. Spec 105.
type: feature
---

# Feature: `release-version.ps1` / `release-version.sh` (spec 105, shipped v3.39.0; web Vitest wiring v4.24.0)

**Spec:** `spec/01-app/105-release-version-script.md`
**Replaces:** the (incorrect) practice of pointing `/release/:version` pages at the generic latest-resolving `install.ps1`.

## Behaviour summary

- **Two artefacts per release, same script body:**
  - **Generic** — lives at `/scripts/release-version.ps1` on the docs site, reads `-Version` parameter, always current.
  - **Snapshot** — uploaded as a release asset (`release-version-vX.Y.Z.ps1`), has `$Version = 'vX.Y.Z'` baked at line 1, drift-proof forever.
- **Release pages link to snapshot by default.** Generic shown as "advanced".
- **Front-page Get Started box stays on `install.ps1`** (latest-resolving) — out of scope here.
- **Steps:** validate version → detect OS/arch → GitHub API `/releases/tags/<v>` → download → SHA256 verify → extract to install dir → PATH update → chain `gitmap self-install` → run `gitmap --version` and confirm match.
- **Missing version → interactive prompt** listing 5 most recent valid releases (default N quit). Non-interactive (no TTY / `-Quiet` / closed stdin) → exit 1. `-AllowFallback` = same-minor-series patch fallback, no prompt.
- **Word "latest" never appears in the script.** No `/releases/latest`, no `self-update`.

## Flags

| `-Version` | Required for generic; ignored in snapshot |
| `-AllowFallback` | Same-minor patch fallback, no prompt |
| `-Quiet` | No prompts, no progress, non-interactive failure |
| `-InstallDir` | Override default install dir |
| `-NoPath` | Skip PATH modification |
| `-NoSelfInstall` | Skip chained self-install (download + extract only) |

## Exit codes

| 0 | Installed + verified |
| 1 | Requested version missing, user declined fallback |
| 2 | Network/download error |
| 3 | Checksum mismatch |
| 4 | OS/arch unsupported |
| 5 | PATH update failed (warning, not fatal) |
| 6 | self-install chain failed |
| 7 | Verified version mismatch |

## Why these decisions (all chosen by user)

- **Both generic + snapshot** — belt and suspenders. Snapshot survives deletion of the docs site or rewrites of the generic; generic stays maintainable.
- **Name `release-version.ps1`** — matches user's literal phrasing; clearly distinct from `install.ps1`.
- **Chain `self-install`** — matches first-run experience; `-NoSelfInstall` provides escape hatch.
- **Interactive prompt for missing version** — implementation must guard against non-TTY contexts (CI) by exiting 1 rather than blocking forever.

## Implementation footprint

| File | What |
|------|------|
| `gitmap/scripts/release-version.ps1` | New, embedded via `go:embed` |
| `gitmap/scripts/release-version.sh` | New, embedded via `go:embed` |
| `cmd/release.go` | Snapshot generation step in release pipeline |
| `constants/constants_install.go` | `ScriptReleaseVersionPS1`, `ScriptReleaseVersionSh`, snapshot filename format |
| `src/pages/Release.tsx` (docs site) | Render two install boxes per `/release/:version` |
