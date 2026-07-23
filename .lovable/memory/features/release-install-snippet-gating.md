---
name: Release install-snippet gating
description: Pinned-install snippet appended to GitHub release body is gated to versioned gitmap repos only
type: feature
---

# Release Install-Snippet Gating (v5.2.0+)

`gitmap release` appends a PowerShell + bash pinned-version installer snippet
to the GitHub release body via `AppendPinnedInstallSnippet`
(`gitmap/release/installsnippet.go`). This snippet installs the gitmap binary
at a specific tag, so it is meaningless — and confusing — in releases of any
non-gitmap repository.

## Rule

`uploadToGitHub` in `gitmap/release/workflowgithub.go` MUST gate the call
with `ShouldPrintInstallHint(getRemoteURL())`:

```go
body := DetectChangelog()
if ShouldPrintInstallHint(getRemoteURL()) {
    body = AppendPinnedInstallSnippet(body, v.String())
}
```

`ShouldPrintInstallHint` (in `releaseinstallhint.go`) returns true only when
the remote matches `alimtvnetwork/gitmap-vN` (numeric suffix). The same gate
already governs the terminal install hint printed by `printInstallHint`.

## Root cause of v5.1.x bug

The snippet was originally added (spec
`spec/04-release/08-pinned-version-install-snippet.md`) without the gate, so
every release body in every repo got the gitmap installer block plus a
`## gitmap vX.Y.Z` header. Fixed in v5.2.0.

## Test guidance

Any future change to release-body composition MUST preserve the
`ShouldPrintInstallHint` gate. Add a test that calls `uploadToGitHub`-level
composition with a non-gitmap remote and asserts the body does NOT contain
the snippet marker.
