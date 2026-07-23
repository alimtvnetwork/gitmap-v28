---
name: cn folder-arg dispatch
description: gitmap cn accepts `cn vX <folder>` and `cn <folder>` (defaults to v++) — chdirs into folder, runs in-place pipeline, chdirs back. Disambiguation via path-sep + os.Stat heuristic so release-alias form keeps working.
type: feature
---

# `gitmap cn` folder-arg forms (v3.117.0+)

## New forms

- `gitmap cn vX <folder>` — explicit version, explicit folder.
- `gitmap cn v+1 <folder>` / `cn v++ <folder>` — version-bump shortcuts.
- `gitmap cn <folder>` — single positional, defaults to `v++`.

`<folder>` accepts absolute, relative, `~`-prefixed, and bare-name
paths. Spaces honored when shell-quoted.

## Disambiguation

`looksLikeVersion` extended to match `v++` and `v+N` in addition to
`v?N.N.N`. `isFolderShaped` returns true when the token contains
`/`, `\`, or starts with `~`, OR `os.Stat` succeeds as a directory.

Two-positional truth table:
- version + folder-shaped → NEW form
- folder-shaped + version → EXISTING release-alias form
- both versions or both folders → exit 1 with canonical message

## Pipeline

Reuses `performCrossDirCloneNext` (chdir → runCloneNext → chdir back).
All flags forwarded via `extractFlagArgs`. No new flags introduced.

## References

- Spec: `spec/01-app/111-cn-folder-arg.md`
- Plan: `.lovable/memory/plans/08-cn-folder-arg-plan.md`
- Dispatcher: `gitmap/cmd/clonenextfolderdispatch.go`
- Cross-dir reuse: `gitmap/cmd/clonenextcrossdir.go`
