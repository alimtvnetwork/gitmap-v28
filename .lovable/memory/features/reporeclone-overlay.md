---
name: Repo Reclone Overlay
description: Single-repo `gitmap reclone` / `rc` overlay that wipes + re-clones from remote.origin.url
type: feature
---

# Repo Reclone Overlay (v6.5.0+)

Single-repo flow layered on top of the existing manifest-based `reclone` pipeline.
Spec: `spec/04-generic-cli/32-repo-reclone.md`.

## Trigger shape
- `gitmap reclone` / `rc` / `rec` / `relclone` / `clone-now` from **inside** a git repo (no positional, cwd has `.git/`).
- `gitmap rc <folder>` where `<folder>` resolves to a sibling directory containing `.git/`.
- Any other shape (multiple positionals, missing path, non-git path, manifest-looking arg) falls through to the original manifest pipeline — overlay must not regress it.

## Flow
1. Read `remote.origin.url` via `git -C <target> config --get remote.origin.url`. Empty → exit 1.
2. Print plan: target abs path + origin URL.
3. Confirm: interactive `y/N` prompt, or skip with `-y` / `--yes`. Non-TTY without `-y` → refuse (protects `yes | gitmap rc` in CI).
4. `escapeCwdIfInside(target)` releases the Windows cwd handle, then `os.RemoveAll(target)`.
5. `git clone <url> <target>` into the parent. Write shell handoff so the user lands back inside the new clone.

## Exit codes
- 0 — re-clone succeeded.
- 1 — aborted at prompt, missing origin URL, clone failure, or non-TTY without `-y`.

## Code map
- Entry split: `runCloneNow` in `gitmap/cmd/clonenow.go` calls `splitRepoRecloneArgs` + `resolveRepoRecloneTarget`; on overlay match dispatches to `runRepoReclone`.
- Overlay: `gitmap/cmd/reporeclone.go` (pipeline) + `reporeclone_test.go` (helpers) + `reporeclone_e2e_test.go` (in-process round-trip against a local bare repo).
- Constants: `gitmap/constants/constants_reporeclone.go` (messages + exit codes, no magic strings).
- Help: `gitmap/helptext/reclone.md` — canonical page; `clone-now.md` is a stub redirect.

## Invariants
- Destructive `os.RemoveAll` MUST be gated by prompt-or-`-y`. Never delete on a bare positional alone.
- Manifest pipeline behavior is unchanged for any non-overlay shape — covered by `cmd/...` regression tests that don't pass single-repo args.
- `swapStdin` in tests must point at `/dev/null` (or equivalent) so e2e never blocks.
