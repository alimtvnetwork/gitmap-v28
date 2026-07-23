# `gitmap-v28 release` must not leak gitmap-specific content into other repos' releases

## Symptom

Running `gitmap-v28 release` inside a non-gitmap repository (e.g. `img-pdf-v2`,
`some-tool-v3`, etc.) produced GitHub releases that:

1. **Attached `release-version-vX.Y.Z.ps1` / `release-version-vX.Y.Z.sh`
   snapshot assets** whose body hard-coded
   `REPO="alimtvnetwork/gitmap-v28"` and `BINARY_NAME="gitmap-v28"`.
   A user who copy-pasted that one-liner would download and install the
   gitmap-v28 binary into the wrong repo's release flow.

2. **Dumped gitmap-v28's `CHANGELOG.md` notes into the release body** because
   `DetectChangelog()` was always called and `AppendPinnedInstallSnippet`
   was the only gated step. The body therefore contained release notes
   for a totally unrelated project.

## Root cause

Both pieces of behavior were unconditional in `pushAndFinalize` /
`uploadToGitHub`. The existing `ShouldPrintInstallHint` gate (which
matches `alimtvnetwork/gitmap-v<N>` remotes only) was applied only to
the install snippet, not to the snapshot asset builder or to the
changelog-body resolution.

## Fix (v5.16.0)

- `gitmap-v28/release/workflowfinalize.go::pushAndFinalize` —
  `buildReleaseVersionSnapshots` is now wrapped in
  `if ShouldPrintInstallHint(getRemoteURL())`. Non-gitmap repos no
  longer receive the `.ps1` / `.sh` snapshot assets.
- `gitmap-v28/release/workflowgithub.go::uploadToGitHub` — the release body
  is empty by default. Only when `ShouldPrintInstallHint` returns true
  (i.e. the current repo is a `alimtvnetwork/gitmap-v<N>` source repo)
  do we call `DetectChangelog()` and append the pinned install snippet.

## Contract

For any repository whose origin remote does NOT match
`alimtvnetwork/gitmap-v<N>`:

- The GitHub release is created with the **tag + name only**; the body
  is an empty string.
- No `release-version-vX.Y.Z.{ps1,sh}` snapshots are uploaded as assets.
- Every other asset (Go cross-compiled binaries, zip groups, ad-hoc
  bundles, checksums, docs-site) is unaffected.

For the canonical gitmap-v28 source repos (`gitmap-v28`, `gitmap-v28`, …):

- Behavior is unchanged. Body = `DetectChangelog()` +
  `AppendPinnedInstallSnippet`. Snapshot assets are still attached.

## Related

- `gitmap-v28/release/releaseinstallhint.go::ShouldPrintInstallHint` — the
  single source of truth for "is this the gitmap-v28 source repo?".
- `spec/01-app/105-release-version-script.md` — the snapshot script spec
  (already gitmap-specific; this gate enforces that scoping at the
  publishing layer).
