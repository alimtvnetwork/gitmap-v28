---
name: cfrp/cfr remote-based no-vN detection
description: clone-fix-repo / cfrp must decide no-vN from the cloned Git remote repo name, not the flattened local folder; skip only when the remote repo lacks -vN. v5.1.0+.
type: feature
---

# Feature: cfrp/cfr remote-based no-vN detection (v5.1.0+)

## Rule
`gitmap clone-fix-repo` (cfr) and `gitmap clone-fix-repo-pub` (cfrp)
MUST check the cloned Git remote repo name for the `-vN` suffix, not
the flattened local destination folder. Example: cloning
`https://github.com/alimtvnetwork/gitmap-v27` into local folder
`gitmap/` MUST still run `fix-repo --all` because the remote repo is
versioned.

## Behavior
- Default: skip with `fix-repo: skipped (repo "<name>" has no -vN suffix, nothing to rewrite)` only when the Git remote repo name lacks `-vN`. Exit 0.
- `--require-version` flag: restore strict mode → exit `ExitCloneFixRepoChainFailed` (10) with a clear "ERROR --require-version set" message.

## Why
Original report (May 2026): `gitmap cfrp https://github.com/.../img-pdf`
clone succeeded, then `fix-repo: ERROR no -vN suffix on repo name "img-pdf"`
killed the pipeline AFTER the clone + Desktop + VS Code side effects had
already run. Result: the user got the artifacts but a non-zero exit and a
scary error. The fix-repo step is meaningless on non-versioned repos —
there's nothing to rewrite — so the correct default is "skip silently
(with a notice)".

## Files
- `gitmap/cmd/clonefixrepo.go::maybeRunFixRepoStep` — gates the chained step on remote-derived repo identity, falling back to local folder only if remote lookup fails.
- `gitmap/constants/constants_clonefixrepo.go` — `FlagRequireVersion`, `MsgCloneFixRepoSkipNoVer`, `ErrCloneFixRepoNeedVersion`.

## Spec
- `spec/02-app-issues/26-cfrp-no-version-suffix-hard-error.md` (TODO follow-up)
