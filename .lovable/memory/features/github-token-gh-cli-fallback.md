---
name: GitHub token gh CLI fallback
description: All GitHub API calls (release upload, repo creation) resolve token via ghtoken.Resolve() — env GITHUB_TOKEN/GH_TOKEN, then `gh auth token` from the locally installed GitHub CLI. v4.43.0+.
type: feature
---

# Feature: GitHub Token gh CLI Fallback (v4.43.0+)

## Rule
Every code path that needs a GitHub API token MUST call
`gitmap/ghtoken.Resolve()` instead of reading `GITHUB_TOKEN`
directly. Resolution order is fixed and documented in the package
header.

## Resolution order
1. `GITHUB_TOKEN` env var.
2. `GH_TOKEN` env var (gh CLI's own override).
3. `gh auth token` shell-out — only when `gh` is on PATH.

On success: returns `(token, source, nil)`. On failure: returns
`("", SourceNone, ErrNoToken)`.

## Why
User ran `gitmap r v1.3.4`, release succeeded, but
`Error: GITHUB_TOKEN not set — skipping asset upload` killed the
asset-upload step even though `gh auth login` had already cached a
token. The expectation: if the official GitHub CLI is authenticated,
gitmap should use it transparently.

## Files
- `gitmap/ghtoken/ghtoken.go` — `Resolve()`, `Source`, `ErrNoToken`.
- `gitmap/release/workflowgithub.go` — release+asset upload.
- `gitmap/clonenext/github.go` — `RepoExists`, `CreateRepo`.
- `gitmap/constants/constants_assets.go` — `MsgTokenFromSource` (green log line), updated `ErrAssetNoToken` with multi-source hint.

## UX
On success the release flow prints a colored line:
`🔑 Authenticated via GitHub CLI (gh auth token)`
On failure the warning lists ALL three sources tried with a fix
hint: `Fix: run \`gh auth login\` once, or set GITHUB_TOKEN`.

## Future
Same resolver should replace any future direct `os.Getenv("GITHUB_TOKEN")` calls — search-replace audit is a one-liner.
