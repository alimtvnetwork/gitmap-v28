---
name: Release gitmap-only body + snapshots
description: gitmap release attaches release-version.{ps1,sh} snapshots and CHANGELOG body only to gitmap source repos; non-gitmap repos get tag-only releases (v5.16.0+)
type: feature
---

# Release: gitmap-only body + snapshot assets

As of **v5.16.0**, `gitmap release` is strict about what it puts into
GitHub releases of repositories that are NOT the gitmap source repo.

## Rule

`ShouldPrintInstallHint(getRemoteURL())` is the single gate. It is true
only when the current repo's remote matches
`alimtvnetwork/gitmap-v<N>`.

When the gate is **false** (any other repo):

- **Release body is empty.** `DetectChangelog()` is NOT called.
  `AppendPinnedInstallSnippet` is NOT called.
- **No `release-version-vX.Y.Z.{ps1,sh}` snapshot assets** are uploaded.
  Those snapshots hard-code `REPO="alimtvnetwork/gitmap-v28"` and
  `BINARY_NAME="gitmap"` and would download the gitmap binary if a user
  ran them from a non-gitmap release page.
- Every other asset still flows: cross-compiled Go binaries, zip groups,
  ad-hoc bundles, checksums, docs-site asset.

When the gate is **true** (gitmap source repo): behavior is unchanged.

## Why

The two pieces of gitmap-specific content were unconditional before
v5.16.0. A user releasing `img-pdf-v2` saw gitmap's CHANGELOG.md in
their release body and `release-version.sh` files pointing at the
wrong repo as assets.

## Code

- `gitmap/release/workflowfinalize.go::pushAndFinalize` — wraps
  `buildReleaseVersionSnapshots` in the gate.
- `gitmap/release/workflowgithub.go::uploadToGitHub` — body stays `""`
  unless the gate fires.
- `gitmap/release/releaseinstallhint.go::ShouldPrintInstallHint` — the
  matcher.

## Spec

`spec/02-app-issues/27-release-body-and-snapshots-gitmap-only.md`
